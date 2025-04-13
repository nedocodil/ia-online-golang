package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ia-online-golang/internal/models"
)

type PasswordCodeRepositoryI interface {
	PasswordCode(ctx context.Context, password_code string) (models.PasswordCode, error)
	PasswordCodeByUserId(ctx context.Context, userID int64) (models.PasswordCode, error)
	SavePasswordCode(ctx context.Context, password_code models.PasswordCode) error
	DeletePasswordCode(ctx context.Context, password_code string) error
	UpdatePasswordCode(ctx context.Context, password_code models.PasswordCode) error
}

var (
	ErrPasswordCodeIsNotFound   = errors.New("password code is not found")
	ErrPasswordCodeInNotUpdated = errors.New("password code is not updated")
)

func (s *Storage) PasswordCode(ctx context.Context, password_code string) (models.PasswordCode, error) {
	const op = "storage.passwordcode.PasswordCode"

	var PasswordCode models.PasswordCode
	query := "SELECT id, user_id, code, expires_at FROM password_codes WHERE code = $1"
	err := s.db.QueryRowContext(ctx, query, password_code).Scan(
		&PasswordCode.ID,
		&PasswordCode.UserID,
		&PasswordCode.Code,
		&PasswordCode.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.PasswordCode{}, ErrPasswordCodeIsNotFound
		}
		return models.PasswordCode{}, fmt.Errorf("%s: %w", op, err)
	}

	return PasswordCode, nil
}

func (s *Storage) PasswordCodeByUserId(ctx context.Context, userID int64) (models.PasswordCode, error) {
	const op = "storage.passwordcode.PasswordCodeByUserId"

	var PasswordCode models.PasswordCode
	query := "SELECT id, user_id, code, expires_at FROM password_codes WHERE user_id = $1"
	err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&PasswordCode.ID,
		&PasswordCode.UserID,
		&PasswordCode.Code,
		&PasswordCode.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.PasswordCode{}, ErrPasswordCodeIsNotFound
		}
		return models.PasswordCode{}, fmt.Errorf("%s: %w", op, err)
	}

	return PasswordCode, nil
}

func (s *Storage) SavePasswordCode(ctx context.Context, password_code models.PasswordCode) error {
	const op = "storage.passwordcode.SavePasswordCode"

	query := "INSERT INTO password_codes (user_id, code, expires_at) VALUES ($1, $2, $3)"
	_, err := s.db.ExecContext(ctx, query, password_code.UserID, password_code.Code, password_code.ExpiresAt)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdatePasswordCode(ctx context.Context, password_code models.PasswordCode) error {
	const op = "storage.passwordcode.UpdatePasswordCode"

	query := `
		UPDATE password_codes
		SET code = $1, expires_at = $2
		WHERE user_id = $3`
	_, err := s.db.ExecContext(ctx, query, password_code.Code, password_code.ExpiresAt, password_code.UserID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeletePasswordCode(ctx context.Context, password_code string) error {
	const op = "storage.passwordcode.DeletePasswordCode"

	query := "DELETE FROM password_codes WHERE code = $1"
	_, err := s.db.ExecContext(ctx, query, password_code)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
