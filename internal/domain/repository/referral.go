package repository

import (
	"context"
	"ia-online-golang/internal/domain/models"
)

type ReferralRepository interface {
	ReferralByReferralId(ctx context.Context, referral_id string) (models.Referral, error)
	SaveReferral(ctx context.Context, userId int64, referral_id string) error
}
