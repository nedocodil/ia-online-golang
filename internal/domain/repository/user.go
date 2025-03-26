package repository

import (
	"context"
	"ia-online-golang/internal/domain/models"
)

type UserRepository interface {
	UserByReferralCode(ctx context.Context, referral_code string) (models.User, error)
	Users(ctx context.Context) ([]models.User, error)
	UserByEmail(ctx context.Context, email string) (models.User, error)
	UserById(ctx context.Context, id int64) (models.User, error)
	UserIdByEmail(ctx context.Context, email string) (int64, error)
	UserIdByPhone(ctx context.Context, phone string) (int64, error)
	ValidationUser(ctx context.Context, email string, phone string) error
	UserIdByTelegram(ctx context.Context, telegram string) (int64, error)
	CreateUser(ctx context.Context, user models.User) (int64, error)
	UpdateActiveUser(ctx context.Context, userID int64, isActive bool) error
	UpdateUser(ctx context.Context, user models.User) error
	DeleteUser(ctx context.Context, id int) error
}
