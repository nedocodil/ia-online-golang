package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"ia-online-golang/internal/domain/models"
	"ia-online-golang/internal/domain/repository"
	"ia-online-golang/internal/domain/services"

	"ia-online-golang/internal/interfaces/dto"
	"ia-online-golang/internal/interfaces/errors/repositories"
	ServiceAuthErrors "ia-online-golang/internal/interfaces/errors/services"
	payloads "ia-online-golang/internal/interfaces/payload"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	log                      *logrus.Logger
	Address                  string
	UserRepository           repository.UserRepository
	TokenRepository          repository.TokenRepository
	ReferralRepository       repository.ReferralRepository
	ActivationLinkRepository repository.ActivationLinkRepository
	TokenService             services.TokenService
	EmailService             services.EmailService
}

func New(
	log *logrus.Logger,
	address string,
	userRepo repository.UserRepository,
	tokenRepo repository.TokenRepository,
	referralRepo repository.ReferralRepository,
	activationLinkRepo repository.ActivationLinkRepository,
	tokenService services.TokenService,
	emailService services.EmailService,
) *AuthService {
	return &AuthService{
		log:                      log,
		Address:                  address,
		UserRepository:           userRepo,
		TokenRepository:          tokenRepo,
		ReferralRepository:       referralRepo,
		ActivationLinkRepository: activationLinkRepo,
		TokenService:             tokenService,
		EmailService:             emailService,
	}
}

func (a *AuthService) RegistrationUser(ctx context.Context, registerDTO dto.RegisterUserDTO) (dto.TokensDTO, dto.RegisterUserDTO, error) {
	const op = "AuthService.RegistrationUser"

	err := a.UserRepository.ValidationUser(ctx, registerDTO.Email, registerDTO.PhoneNumber)
	if err != nil {
		if errors.Is(err, repositories.ErrUserExists) {
			return dto.TokensDTO{}, dto.RegisterUserDTO{}, fmt.Errorf("%s: %w", op, ServiceAuthErrors.ErrUserAlreadyExists)
		}

		return dto.TokensDTO{}, dto.RegisterUserDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	if registerDTO.Telegram != "" {
		_, err = a.UserRepository.UserIdByTelegram(ctx, registerDTO.Telegram)
		if err == nil {
			a.log.Error(ServiceAuthErrors.ErrUserAlreadyExists)

			return dto.TokensDTO{}, dto.RegisterUserDTO{}, fmt.Errorf("%s: %w", op, ServiceAuthErrors.ErrUserAlreadyExists)
		}
	}

	if registerDTO.ReferralCode != "" {
		_, err := a.UserRepository.UserByReferralCode(ctx, registerDTO.ReferralCode)
		if err != nil {
			if errors.Is(err, repositories.ErrUserNotFound) {
				return dto.TokensDTO{}, dto.RegisterUserDTO{}, fmt.Errorf("%s: %w", op, ServiceAuthErrors.ErrReferralAlreadyUsed)
			}
		}
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(registerDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		a.log.Error(err)

		return dto.TokensDTO{}, dto.RegisterUserDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	referralCode := uuid.New()

	user := models.User{
		PhoneNumber:  registerDTO.PhoneNumber,
		Email:        registerDTO.Email,
		Name:         registerDTO.Name,
		Telegram:     registerDTO.Telegram,
		City:         registerDTO.City,
		PasswordHash: string(passHash),
		ReferralCode: referralCode.String(),
		Role:         "user",
		IsActive:     false,
	}

	idUser, err := a.UserRepository.CreateUser(ctx, user)
	if err != nil {
		a.log.Error(err)

		return dto.TokensDTO{}, dto.RegisterUserDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	err = a.ReferralRepository.SaveReferral(ctx, idUser, registerDTO.ReferralCode)
	if err != nil {
		a.log.Error(err)

		return dto.TokensDTO{}, dto.RegisterUserDTO{}, fmt.Errorf("%s: %v", op, err)
	}

	payloadAccess := payloads.PayloadUserAccess{
		UserID:       idUser,
		Name:         user.Name,
		Email:        user.Email,
		PhoneNumber:  user.PhoneNumber,
		City:         user.City,
		Telegram:     user.Telegram,
		ReferralCode: user.ReferralCode,
	}
	payloadRefresh := payloads.PayloadUserRefresh{
		UserID: idUser,
	}

	access_token, refresh_token, err := a.TokenService.GenerateTokens(ctx, payloadAccess, payloadRefresh)
	if err != nil {
		a.log.Error(err)

		return dto.TokensDTO{}, dto.RegisterUserDTO{}, err
	}

	err = a.TokenService.SaveToken(ctx, idUser, refresh_token)
	if err != nil {
		a.log.Error(err)

		return dto.TokensDTO{}, dto.RegisterUserDTO{}, err
	}

	err = a.SendActivationLink(ctx, idUser, registerDTO.Email)
	if err != nil {
		a.log.Error(err)

		return dto.TokensDTO{}, dto.RegisterUserDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	tokens := dto.TokensDTO{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}

	return tokens, registerDTO, nil
}

func (a *AuthService) LoginUser(ctx context.Context, loginDTO dto.LoginUserDTO) (dto.TokensDTO, error) {
	const op = "AuthService.LoginUser"

	user, err := a.UserRepository.UserByEmail(ctx, loginDTO.Email)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			a.log.Error(err)

			return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, ServiceAuthErrors.ErrUserNotFound)
		}

		a.log.Error(err)

		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginDTO.Password))
	if err != nil {
		a.log.Error(err)

		return dto.TokensDTO{}, ServiceAuthErrors.ErrIncorrectPassword
	}

	if !user.IsActive {
		a.log.Error(err)

		err = a.SendActivationLink(ctx, user.ID, user.Email)
		if err != nil {
			a.log.Error(err)

			return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
		}

		return dto.TokensDTO{}, ServiceAuthErrors.ErrUserNotActivated
	}

	payloadAccess := payloads.PayloadUserAccess{
		UserID:       user.ID,
		Name:         user.Name,
		Email:        user.Email,
		PhoneNumber:  user.PhoneNumber,
		City:         user.City,
		Telegram:     user.Telegram,
		ReferralCode: user.ReferralCode,
	}
	payloadRefresh := payloads.PayloadUserRefresh{
		UserID: user.ID,
	}

	access_token, refresh_token, err := a.TokenService.GenerateTokens(ctx, payloadAccess, payloadRefresh)
	if err != nil {
		a.log.Error(err)

		return dto.TokensDTO{}, err
	}

	err = a.TokenService.SaveToken(ctx, user.ID, refresh_token)
	if err != nil {
		a.log.Error(err)

		return dto.TokensDTO{}, err
	}

	tokens := dto.TokensDTO{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}

	return tokens, nil
}

func (a *AuthService) SendActivationLink(ctx context.Context, userID int64, email string) error {
	op := "AuthService.SendActivationLink"

	var activationObj models.ActivationLink

	activationObj, err := a.ActivationLinkRepository.ActivationLinkByUserId(ctx, userID)
	if err != nil {
		if errors.Is(err, repositories.ErrActivationLinkIsNotFound) {
			activationObj = models.ActivationLink{
				UserID:       userID,
				ActivationID: uuid.New().String(),
				ExpiresAt:    time.Now().Add(24 * time.Hour),
			}
			err = a.ActivationLinkRepository.SaveActivationLink(ctx, activationObj)
			if err != nil {
				a.log.Error(err)

				return fmt.Errorf("%s: %w", op, err)
			}
		} else {
			a.log.Error(err)

			return fmt.Errorf("%s: %w", op, err)
		}
	}

	if activationObj.ExpiresAt.Before(time.Now()) {
		activationObj = models.ActivationLink{
			UserID:       userID,
			ActivationID: uuid.New().String(),
			ExpiresAt:    time.Now().Add(24 * time.Hour),
		}
		err = a.ActivationLinkRepository.UpdateActivationLink(ctx, activationObj)
		if err != nil {
			a.log.Error(err)

			return fmt.Errorf("%s: %w", op, err)
		}
	}

	activationLink := "https://" + a.Address + "/api/v1/auth/activation/" + activationObj.ActivationID
	err = a.EmailService.SendActivationLink(ctx, email, activationLink)
	if err != nil {
		a.log.Error(err)

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *AuthService) LogoutUser(ctx context.Context, refreshToken string) error {
	op := "TokenService.DeleteToken"

	err := a.TokenRepository.DeleteRefreshToken(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, repositories.ErrTokenNotFound) {
			return ServiceAuthErrors.ErrRefreshTokenNotExists
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *AuthService) ActivationUser(ctx context.Context, activationID string) error {
	const op = "AuthService.ActivationUser"

	activation, err := a.ActivationLinkRepository.ActivationLinkByActivationId(ctx, activationID)
	if err != nil {
		if errors.Is(err, repositories.ErrActivationLinkIsNotFound) {
			return ServiceAuthErrors.ErrActiveLinkNotExists
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	// Проверяем, не истек ли срок действия активационной ссылки
	if activation.ExpiresAt.Before(time.Now()) {
		return ServiceAuthErrors.ErrActiveLinkExpired
	}

	err = a.UserRepository.UpdateActiveUser(ctx, activation.UserID, true)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *AuthService) RefreshUserToken(ctx context.Context, refresh_token string) (dto.TokensDTO, error) {
	op := "AuthService.RefreshToken"

	// Создаём переменную payload и передаём её в `ValidateRefreshToken`
	var payload payloads.PayloadUserRefresh
	_, err := a.TokenService.ValidateRefreshToken(ctx, refresh_token, &payload)
	if err != nil {
		if errors.Is(err, ServiceAuthErrors.ErrInvalidRefreshToken) {
			return dto.TokensDTO{}, ServiceAuthErrors.ErrInvalidRefreshToken
		}
		if errors.Is(err, ServiceAuthErrors.ErrRefreshTokenNotExists) {
			return dto.TokensDTO{}, ServiceAuthErrors.ErrRefreshTokenNotExists
		}
		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Получаем пользователя по `UserID`
	user, err := a.UserRepository.UserById(ctx, payload.UserID)
	if err != nil {
		return dto.TokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Создаём payload'ы для токенов
	payloadAccess := payloads.PayloadUserAccess{
		UserID:       user.ID,
		Name:         user.Name,
		Email:        user.Email,
		PhoneNumber:  user.PhoneNumber,
		City:         user.City,
		Telegram:     user.Telegram,
		ReferralCode: user.ReferralCode,
	}

	// Генерируем новые токены
	access_token, err := a.TokenService.GenerateAccessToken(ctx, payloadAccess)
	if err != nil {
		a.log.Errorf("%s: %v", op, err)
		return dto.TokensDTO{}, err
	}

	// Возвращаем токены
	return dto.TokensDTO{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}, nil
}
