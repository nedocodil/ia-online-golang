package services

import (
	"context"
	"ia-online-golang/internal/interfaces/dto"
)

type UserService interface {
	Users(ctx context.Context) ([]dto.UserDTO, error)
}
