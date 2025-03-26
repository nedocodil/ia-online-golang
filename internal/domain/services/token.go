package services

import (
	"context"
)

type TokenService interface {
	GenerateTokens(ctx context.Context, payloadAccess any, payloadRefresh any) (string, string, error)
	GenerateAccessToken(ctx context.Context, payloadAccess any) (string, error)
	SaveToken(ctx context.Context, userId int64, refreshToken string) error
	ValidateRefreshToken(ctx context.Context, refresh_token string, payloadStruct any) (any, error)
	ValidateAccessToken(ctx context.Context, token string, payloadStruct any) (any, error)
}
