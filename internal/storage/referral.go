package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ia-online-golang/internal/models"
)

type ReferralRepositoryI interface {
	ReferralByReferralId(ctx context.Context, referral_id string) (models.Referral, error)
	SaveReferral(ctx context.Context, userId int64, referral_id string) error
	ReferralsUser(ctx context.Context, referral_id string) ([]models.ReferralAndUser, error)
	Referrals(ctx context.Context) ([]models.Referral, error)
	GetInactiveReferralsWithReadyLeads(ctx context.Context) ([]models.Referral, error)
	UpdateActive(ctx context.Context, referral_id int64, active bool) error
	ActiveReferralsByReferralId(ctx context.Context, referral_id string) ([]models.Referral, error)
}

var (
	ErrReferralNotFound  = errors.New("referral not found")
	ErrReferralsNotFound = errors.New("referrals not found")
)

func (s *Storage) ReferralByReferralId(ctx context.Context, referral_id string) (models.Referral, error) {
	const op = "storage.auth.ReferralByReferralId"
	var referral models.Referral

	// Запрос для получения refresh-токена по user_id
	query := "SELECT id, user_id, referral_id, cost FROM referrals WHERE referral_id = $1"
	err := s.db.QueryRowContext(ctx, query, referral_id).Scan(
		&referral.ID,
		&referral.UserID,
		&referral.ReferralCode,
		&referral.Cost,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Referral{}, ErrReferralNotFound
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

func (s *Storage) ReferralsUser(ctx context.Context, referral_id string) ([]models.ReferralAndUser, error) {
	const op = "storage.auth.ReferralsUser"

	query := `
		SELECT u.id, u.phone_number, u.email, u.name, u.telegram, u.city,
		       u.password_hash, u.referral_code, u.created_at, u.is_active,
		       u.roles, r.id, r.user_id, r.referral_id, r.created_at, r.active
		FROM users u
		JOIN referrals r ON u.id = r.user_id
		WHERE r.referral_id = $1
	`

	rows, err := s.db.QueryContext(ctx, query, referral_id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var referralsAndUsers []models.ReferralAndUser
	for rows.Next() {
		var ru models.ReferralAndUser
		err := rows.Scan(
			&ru.User.ID, &ru.User.PhoneNumber, &ru.User.Email, &ru.User.Name, &ru.User.Telegram, &ru.User.City,
			&ru.User.PasswordHash, &ru.User.ReferralCode, &ru.User.CreatedAt, &ru.User.IsActive,
			&ru.User.Roles, &ru.Referral.ID, &ru.Referral.UserID, &ru.Referral.ReferralCode, &ru.Referral.CreatedAt, &ru.Referral.Active,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		referralsAndUsers = append(referralsAndUsers, ru)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(referralsAndUsers) == 0 {
		return nil, ErrReferralsNotFound
	}

	return referralsAndUsers, nil
}

func (s *Storage) Referrals(ctx context.Context) ([]models.Referral, error) {
	const op = "storage.Referrals"

	query := "SELECT id, user_id, referral_id, created_at, active, cost FROM referrals"

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var referrals []models.Referral
	for rows.Next() {
		var r models.Referral
		err := rows.Scan(&r.ID, &r.UserID, &r.ReferralCode, &r.CreatedAt, &r.Active, &r.Cost)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		referrals = append(referrals, r)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(referrals) == 0 {
		return nil, ErrReferralsNotFound
	}

	return referrals, nil
}

func (s *Storage) ActiveReferralsByReferralId(ctx context.Context, referralID string) ([]models.Referral, error) {
	const op = "storage.Referrals"

	query := "SELECT id, user_id, referral_id, created_at, active, cost FROM referrals WHERE referral_id = $1 AND active = true"

	rows, err := s.db.QueryContext(ctx, query, referralID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var referrals []models.Referral
	for rows.Next() {
		var r models.Referral
		err := rows.Scan(&r.ID, &r.UserID, &r.ReferralCode, &r.CreatedAt, &r.Active, &r.Cost)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		referrals = append(referrals, r)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(referrals) == 0 {
		return nil, ErrReferralsNotFound
	}

	return referrals, nil
}

func (s *Storage) GetInactiveReferralsWithReadyLeads(ctx context.Context) ([]models.Referral, error) {
	const op = "storage.referral.GetInactiveReferralsWithReadyLeads"

	query := `
				SELECT r.*
		FROM referrals r
		JOIN users u ON u.id = r.user_id
		JOIN leads l ON l.user_id = u.id
		WHERE r.active = false
		  AND l.status_id = 4
		GROUP BY r.id
		HAVING COUNT(l.id) > 2;

	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var referrals []models.Referral
	for rows.Next() {
		var ref models.Referral
		if err := rows.Scan(&ref.ID, &ref.UserID, &ref.ReferralCode, &ref.CreatedAt, &ref.Active, &ref.Cost); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		referrals = append(referrals, ref)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(referrals) == 0 {
		return nil, ErrReferralsNotFound
	}

	return referrals, nil
}

func (s *Storage) UpdateActive(ctx context.Context, referral_id int64, active bool) error {
	const op = "storage.referral.UpdateReferral"

	query := "UPDATE referrals SET active = $2 WHERE id = $1"
	result, err := s.db.ExecContext(ctx, query, referral_id, active)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Проверяем, сколько строк было обновлено
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Если обновлено 0 строк, значит пользователь не найден
	if rowsAffected == 0 {
		return ErrReferralNotFound
	}

	return nil
}
