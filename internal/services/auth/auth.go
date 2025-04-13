package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"ia-online-golang/internal/services/email"
	"ia-online-golang/internal/services/passwordcode"
	"ia-online-golang/internal/services/token"
	UserService "ia-online-golang/internal/services/user"
	"ia-online-golang/internal/utils"

	"ia-online-golang/internal/dto"
	"ia-online-golang/internal/models"
	"ia-online-golang/internal/storage"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	log                      *logrus.Logger
	Address                  string
	UserRepository           storage.UserRepositoryI
	TokenRepository          storage.TokenRepositoryI
	ReferralRepository       storage.ReferralRepositoryI
	ActivationLinkRepository storage.ActivationLinkRepositoryI
	TokenService             token.TokenServiceI
	EmailService             email.EmailServiceI
	UserService              UserService.UserServiceI
	PasswordCodeService      passwordcode.PasswordCodeServiceI
}

type AuthServiceI interface {
	RegistrationUser(ctx context.Context, registerDTO dto.RegisterUserDTO) (dto.AuthTokensDTO, error)
	ActivationUser(ctx context.Context, activation_id string) error
	LoginUser(ctx context.Context, loginDTO dto.LoginUserDTO) (dto.AuthTokensDTO, error)
	LogoutUser(ctx context.Context, refreshToken string) error
	RefreshUserTokens(ctx context.Context, refresh_token string) (dto.AuthTokensDTO, error)
	SendActivationLink(ctx context.Context, userID int64, email string) error
	ChangingPassword(ctx context.Context, newPasswordDTO dto.NewPasswordDTO, userID int64) error
	RecoverPassword(ctx context.Context, email string) error
	NewPassword(ctx context.Context, dto dto.RecoverPasswordDTO) error
}

var (
	ErrIncorrectPassword    = errors.New("incorrect password")
	ErrIncorrectOldPassword = errors.New("incorrect old password")
)

func New(
	log *logrus.Logger,
	address string,
	userRepo storage.UserRepositoryI,
	tokenRepo storage.TokenRepositoryI,
	referralRepo storage.ReferralRepositoryI,
	activationLinkRepo storage.ActivationLinkRepositoryI,
	tokenService token.TokenServiceI,
	emailService email.EmailServiceI,
	userService UserService.UserServiceI,
	passwordCodeService passwordcode.PasswordCodeServiceI,

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
		UserService:              userService,
	}
}

func (a *AuthService) RegistrationUser(ctx context.Context, registerDTO dto.RegisterUserDTO) (dto.AuthTokensDTO, error) {
	const op = "AuthService.RegistrationUser"

	if registerDTO.ReferralCode != "" {
		_, err := a.UserRepository.UserByReferralCode(ctx, registerDTO.ReferralCode)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				return dto.AuthTokensDTO{}, fmt.Errorf("%s: %w", op, ErrReferralIdNotFound)
			}

			return dto.AuthTokensDTO{}, fmt.Errorf("%s: %w", op, err)
		}
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(registerDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		a.log.Error(err)

		return dto.AuthTokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	userDTO, err := a.UserService.SaveUser(ctx, registerDTO, string(passHash))
	if err != nil {
		if errors.Is(err, UserService.ErrUserAlreadyExists) {
			return dto.AuthTokensDTO{}, UserService.ErrUserAlreadyExists
		}

		a.log.Error(err)

		return dto.AuthTokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	if registerDTO.ReferralCode != "" {
		err = a.ReferralRepository.SaveReferral(ctx, *userDTO.ID, registerDTO.ReferralCode)
		if err != nil {
			a.log.Error(err)

			return dto.AuthTokensDTO{}, fmt.Errorf("%s: %v", op, err)
		}
	}

	err = a.SendActivationLink(ctx, *userDTO.ID, registerDTO.Email)
	if err != nil {
		a.log.Error(err)

		return dto.AuthTokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	tokens, err := a.TokenService.CreateUserTokens(ctx, *userDTO.ID)
	if err != nil {
		a.log.Error(err)

		return dto.AuthTokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	return tokens, nil
}

func (a *AuthService) LoginUser(ctx context.Context, loginDTO dto.LoginUserDTO) (dto.AuthTokensDTO, error) {
	const op = "AuthService.LoginUser"

	user, err := a.UserRepository.UserByEmail(ctx, loginDTO.Email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Error(err)

			return dto.AuthTokensDTO{}, fmt.Errorf("%s: %w", op, UserService.ErrUserNotFound)
		}

		a.log.Error(err)

		return dto.AuthTokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginDTO.Password))
	if err != nil {
		a.log.Error(err)

		return dto.AuthTokensDTO{}, ErrIncorrectPassword
	}

	if !user.IsActive {
		a.log.Error(err)

		err = a.SendActivationLink(ctx, user.ID, user.Email)
		if err != nil {
			a.log.Error(err)

			return dto.AuthTokensDTO{}, fmt.Errorf("%s: %w", op, err)
		}

		return dto.AuthTokensDTO{}, UserService.ErrUserNotActivated
	}

	userDTO := utils.UserToDTO(user)

	tokens, err := a.TokenService.CreateUserTokens(ctx, *userDTO.ID)
	if err != nil {
		a.log.Error(err)

		return dto.AuthTokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	return tokens, nil
}

func (a *AuthService) SendActivationLink(ctx context.Context, userID int64, email string) error {
	op := "AuthService.SendActivationLink"

	var activationObj models.ActivationLink

	activationObj, err := a.ActivationLinkRepository.ActivationLinkByUserId(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrActivationLinkIsNotFound) {
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
		if errors.Is(err, storage.ErrTokenNotFound) {
			return token.ErrRefreshTokenNotExists
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *AuthService) ActivationUser(ctx context.Context, activationID string) error {
	const op = "AuthService.ActivationUser"

	activation, err := a.ActivationLinkRepository.ActivationLinkByActivationId(ctx, activationID)
	if err != nil {
		if errors.Is(err, storage.ErrActivationLinkIsNotFound) {
			return ErrActiveLinkNotExists
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	// Проверяем, не истек ли срок действия активационной ссылки
	if activation.ExpiresAt.Before(time.Now()) {
		return ErrActiveLinkExpired
	}

	err = a.UserRepository.UpdateActiveUser(ctx, activation.UserID, true)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *AuthService) RefreshUserTokens(ctx context.Context, refresh_token string) (dto.AuthTokensDTO, error) {
	op := "AuthService.RefreshToken"

	// Создаём переменную payload и передаём её в `ValidateRefreshToken`
	var payload token.PayloadUserRefresh
	_, err := a.TokenService.ValidateRefreshToken(ctx, refresh_token, &payload)
	if err != nil {
		if errors.Is(err, token.ErrInvalidRefreshToken) {
			return dto.AuthTokensDTO{}, token.ErrInvalidRefreshToken
		}
		if errors.Is(err, token.ErrRefreshTokenNotExists) {
			return dto.AuthTokensDTO{}, token.ErrRefreshTokenNotExists
		}
		return dto.AuthTokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	tokens, err := a.TokenService.CreateUserTokens(ctx, payload.UserID)
	if err != nil {
		a.log.Error(err)

		return dto.AuthTokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	return tokens, nil
}

func (a *AuthService) ChangingPassword(ctx context.Context, newPasswordDTO dto.NewPasswordDTO, userID int64) error {
	op := "AuthService.NewPassword"

	user, err := a.UserRepository.UserById(ctx, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(newPasswordDTO.OldPassword)); err != nil {
		return ErrIncorrectOldPassword
	}

	newPassHash, err := bcrypt.GenerateFromPassword([]byte(newPasswordDTO.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = a.UserRepository.UpdatePasswordUser(ctx, string(newPassHash), userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *AuthService) RecoverPassword(ctx context.Context, email string) error {
	const op = "AuthService.RecoverPassword"

	user, err := a.UserRepository.UserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return UserService.ErrUserNotFound
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	new_password, err := utils.GeneratePasswordCode(8)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	new_password_hash, err := bcrypt.GenerateFromPassword([]byte(new_password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = a.UserRepository.UpdatePasswordUser(ctx, string(new_password_hash), user.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = a.EmailService.SendNewPassword(ctx, email, new_password)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *AuthService) NewPassword(ctx context.Context, dto dto.RecoverPasswordDTO) error {
	op := "AuthService.NewPassword"

	password_code, err := a.PasswordCodeService.PasswordCode(ctx, dto.Code)
	if err != nil {
		if errors.Is(err, passwordcode.ErrPasswordCodeIsNotFound) {
			return passwordcode.ErrPasswordCodeIsNotFound
		}

		if errors.Is(err, passwordcode.ErrPasswordCodeIncorrect) {
			return passwordcode.ErrPasswordCodeIncorrect
		}

		if errors.Is(err, passwordcode.ErrPasswordCodeHasExpired) {
			return passwordcode.ErrPasswordCodeIncorrect
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	newPassHash, err := bcrypt.GenerateFromPassword([]byte(dto.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = a.UserRepository.UpdatePasswordUser(ctx, string(newPassHash), password_code.UserID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = a.PasswordCodeService.DeleteRefreshToken(ctx, password_code.Code)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
