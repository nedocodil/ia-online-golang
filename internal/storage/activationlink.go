package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ia-online-golang/internal/models"
)

type ActivationLinkRepositoryI interface {
	ActivationLinkByActivationId(ctx context.Context, activationID string) (models.ActivationLink, error)
	ActivationLinkByUserId(ctx context.Context, userID int64) (models.ActivationLink, error)
	SaveActivationLink(ctx context.Context, activation models.ActivationLink) error
	DeleteActivationLink(ctx context.Context, activation models.ActivationLink) error
	UpdateActivationLink(ctx context.Context, activation models.ActivationLink) error
}

var (
	ErrActivationLinkIsNotFound   = errors.New("activation link is not found")
	ErrActivationLinkIsNotUpdated = errors.New("activation link is not updated")
	ErrGetActivationLink          = errors.New("error getting activation link")
	ErrSaveActivationLink         = errors.New("error saving activation link")
)

func (s *Storage) ActivationLinkByActivationId(ctx context.Context, activationID string) (models.ActivationLink, error) {
	const op = "storage.auth.GetActivationLink"

	var activationLink models.ActivationLink
	query := "SELECT id, user_id, activation_id, expires_at FROM activation_links WHERE activation_id = $1"
	err := s.db.QueryRowContext(ctx, query, activationID).Scan(
		&activationLink.ID,
		&activationLink.UserID,
		&activationLink.ActivationID,
		&activationLink.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ActivationLink{}, ErrActivationLinkIsNotFound
		}
		return models.ActivationLink{}, fmt.Errorf("%s: %w", op, err)
	}

	return activationLink, nil
}

func (s *Storage) ActivationLinkByUserId(ctx context.Context, userID int64) (models.ActivationLink, error) {
	const op = "storage.auth.ActivationLinkByUserId"

	var activationLink models.ActivationLink
	query := "SELECT id, user_id, activation_id, expires_at FROM activation_links WHERE user_id = $1"
	err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&activationLink.ID,
		&activationLink.UserID,
		&activationLink.ActivationID,
		&activationLink.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ActivationLink{}, ErrActivationLinkIsNotFound
		}
		return models.ActivationLink{}, fmt.Errorf("%s: %w", op, err)
	}

	return activationLink, nil
}

func (s *Storage) SaveActivationLink(ctx context.Context, activation models.ActivationLink) error {
	const op = "storage.auth.SaveActivationLink"

	query := "INSERT INTO activation_links (user_id, activation_id, expires_at) VALUES ($1, $2, $3)"
	result, err := s.db.ExecContext(ctx, query, activation.UserID, activation.ActivationID, activation.ExpiresAt)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return ErrSaveActivationLink
	}

	return nil
}

func (s *Storage) UpdateActivationLink(ctx context.Context, activation models.ActivationLink) error {
	const op = "storage.auth.UpdateActivationLink"

	query := `
		UPDATE activation_links
		SET activation_id = $1, expires_at = $2
		WHERE user_id = $3`
	result, err := s.db.ExecContext(ctx, query, activation.ActivationID, activation.ExpiresAt, activation.UserID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return ErrActivationLinkIsNotUpdated
	}

	return nil
}

func (s *Storage) DeleteActivationLink(ctx context.Context, activation models.ActivationLink) error {
	const op = "storage.auth.DeleteActivationLink"

	query := "DELETE FROM activation_links WHERE activation_id = $1"
	result, err := s.db.ExecContext(ctx, query, activation.ActivationID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return ErrActivationLinkIsNotFound
	}

	return nil
}
