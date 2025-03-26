package repository

import (
	"context"
	"ia-online-golang/internal/domain/models"
)

type ActivationLinkRepository interface {
	ActivationLinkByActivationId(ctx context.Context, activationID string) (models.ActivationLink, error)
	ActivationLinkByUserId(ctx context.Context, userID int64) (models.ActivationLink, error)
	SaveActivationLink(ctx context.Context, activation models.ActivationLink) error
	DeleteActivationLink(ctx context.Context, activation models.ActivationLink) error
	UpdateActivationLink(ctx context.Context, activation models.ActivationLink) error
}
