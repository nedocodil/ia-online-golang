package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ia-online-golang/internal/domain/models"
	"ia-online-golang/internal/interfaces/errors/repositories"
)

func (s *Storage) ReferralByReferralId(ctx context.Context, referral_id string) (models.Referral, error) {
	const op = "storage.auth.ReferralByReferralId"
	var referral models.Referral

	// Запрос для получения refresh-токена по user_id
	query := "SELECT id, user_id, referral_id FROM referrals WHERE referral_id = $1"
	err := s.db.QueryRowContext(ctx, query, referral_id).Scan(
		&referral.ID,
		&referral.UserID,
		&referral.ReferralCode,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Referral{}, repositories.ErrReferralNotFound
		}
		return models.Referral{}, fmt.Errorf("%s: %w", op, err)
	}

	return referral, nil
}

func (s *Storage) SaveReferral(ctx context.Context, userId int64, referral_id string) error {
	const op = "storage.auth.SaveReferral"

	query := "INSERT INTO referrals (user_id, referral_id) VALUES ($1, $2)"
	_, err := s.db.ExecContext(ctx, query, userId, referral_id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
