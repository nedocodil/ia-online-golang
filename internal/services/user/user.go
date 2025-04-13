package user

import (
	"context"
	"errors"
	"fmt"
	"ia-online-golang/internal/dto"
	"ia-online-golang/internal/http/context_keys"
	"ia-online-golang/internal/models"
	"ia-online-golang/internal/storage"
	"ia-online-golang/internal/utils"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type UserService struct {
	log            *logrus.Logger
	UserRepository storage.UserRepositoryI
}

type UserServiceI interface {
	UserById(ctx context.Context, id int64) (dto.UserDTO, error)
	UserByEmail(ctx context.Context, email string) (dto.UserDTO, error)
	SaveUser(ctx context.Context, userRegisterDTO dto.RegisterUserDTO, passHash string) (dto.UserDTO, error)
	Users(ctx context.Context) ([]dto.UserDTO, error)
	EditUser(ctx context.Context, userDTO dto.UserDTO) error
}

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotActivated  = errors.New("user not activated")
	ErrUserNotFound      = errors.New("user not found")
)

func New(
	log *logrus.Logger,
	userRepo storage.UserRepositoryI,
) *UserService {
	return &UserService{
		log:            log,
		UserRepository: userRepo,
	}
}

func (u *UserService) UserById(ctx context.Context, userID int64) (dto.UserDTO, error) {
	op := "UserService.User"

	user, err := u.UserRepository.UserById(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return dto.UserDTO{}, ErrUserNotFound
		}
		return dto.UserDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	userDTO := utils.UserToDTO(user)

	return userDTO, nil
}

func (u *UserService) UserByEmail(ctx context.Context, email string) (dto.UserDTO, error) {
	op := "UserService.UserByEmail"

	user, err := u.UserRepository.UserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return dto.UserDTO{}, ErrUserNotFound
		}
		return dto.UserDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	userDTO := utils.UserToDTO(user)

	return userDTO, nil
}

func (u *UserService) SaveUser(ctx context.Context, userRegisterDTO dto.RegisterUserDTO, passHash string) (dto.UserDTO, error) {
	op := "UserService.SaveUser"

	err := u.UserRepository.ValidationUser(ctx, userRegisterDTO.Email, userRegisterDTO.PhoneNumber, userRegisterDTO.Telegram)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return dto.UserDTO{}, fmt.Errorf("%s: %w", op, ErrUserAlreadyExists)
		}

		return dto.UserDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	referralCode := uuid.New()

	userNew := models.User{
		PhoneNumber:  userRegisterDTO.PhoneNumber,
		Email:        userRegisterDTO.Email,
		Name:         userRegisterDTO.Name,
		Telegram:     userRegisterDTO.Telegram,
		City:         userRegisterDTO.City,
		PasswordHash: passHash,
		ReferralCode: referralCode.String(),
		IsActive:     false,
	}

	user, err := u.UserRepository.CreateUser(ctx, userNew)
	if err != nil {
		u.log.Error(err)

		return dto.UserDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	userDTO := utils.UserToDTO(user)

	return userDTO, nil
}

func (u *UserService) Users(ctx context.Context) ([]dto.UserDTO, error) {
	op := "UserService.Users"

	users, err := u.UserRepository.Users(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var userDTOs []dto.UserDTO
	for _, user := range users {
		userDTOs = append(userDTOs, utils.UserToDTO(user))
	}

	return userDTOs, nil
}

func (u *UserService) EditUser(ctx context.Context, userDTO dto.UserDTO) error {
	op := "UserService.EditUser"

	user := utils.DtoToUser(userDTO)

	if userDTO.ID == nil {
		userIDValue := ctx.Value(context_keys.UserIDKey)
		userID, ok := userIDValue.(int64)
		if !ok {
			return fmt.Errorf("%s: error receiving userID ", op)
		}
		user.ID = userID
	}

	fmt.Println(user)

	err := u.UserRepository.UpdateUser(ctx, user)
	if err != nil {
		if errors.Is(err, storage.ErrUserIsNotUpdated) {
			return ErrUserNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
