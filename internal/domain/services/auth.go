package services

import (
	"context"
	"ia-online-golang/internal/interfaces/dto"
)

type AuthService interface {
	RegistrationUser(ctx context.Context, registerDTO dto.RegisterUserDTO) (dto.TokensDTO, dto.RegisterUserDTO, error)
	ActivationUser(ctx context.Context, activation_id string) error
	LoginUser(ctx context.Context, loginDTO dto.LoginUserDTO) (dto.TokensDTO, error)
	LogoutUser(ctx context.Context, refreshToken string) error
	RefreshUserToken(ctx context.Context, refresh_token string) (dto.TokensDTO, error)
}
