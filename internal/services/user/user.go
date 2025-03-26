package user

import (
	"context"
	"fmt"
	"ia-online-golang/internal/domain/repository"
	"ia-online-golang/internal/interfaces/dto"

	"github.com/sirupsen/logrus"
)

type UserService struct {
	log            *logrus.Logger
	UserRepository repository.UserRepository
}

func New(
	log *logrus.Logger,
	userRepo repository.UserRepository,
) *UserService {
	return &UserService{
		log:            log,
		UserRepository: userRepo,
	}
}

func (u *UserService) Users(ctx context.Context) ([]dto.UserDTO, error) {
	op := "UserService.Users"

	users, err := u.UserRepository.Users(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var userDTOs []dto.UserDTO
	for _, user := range users {
		userDTOs = append(userDTOs, dto.UserDTO{
			ID:             user.ID,
			Email:          user.Email,
			Name:           user.Name,
			PhoneNumber:    user.PhoneNumber,
			City:           user.City,
			RewardInternet: user.RewardInternet.Float64,
			RewardCleaning: user.RewardCleaning.Float64,
			RewardShipping: user.RewardShipping.Float64,
		})
	}

	return userDTOs, nil
}
