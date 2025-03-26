package repository

import (
	"context"
	"ia-online-golang/internal/domain/models"
)

type TokenRepository interface {
	RefreshTokenByUserId(ctx context.Context, userID int64) (models.Token, error)
	RefreshTokenByToken(ctx context.Context, refresh_token string) (models.Token, error)
	SaveRefreshToken(ctx context.Context, userID int64, token string) error
	DeleteRefreshToken(ctx context.Context, refreshToken string) error
	UpdateRefreshToken(ctx context.Context, userID int64, token string) error
}
