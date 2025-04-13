package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ia-online-golang/internal/models"
)

type TokenRepositoryI interface {
	RefreshTokenByUserId(ctx context.Context, userID int64) (models.Token, error)
	RefreshTokenByToken(ctx context.Context, refresh_token string) (models.Token, error)
	SaveRefreshToken(ctx context.Context, userID int64, token string) error
	DeleteRefreshToken(ctx context.Context, refreshToken string) error
	UpdateRefreshToken(ctx context.Context, userID int64, token string) error
}

var (
	ErrTokenNotFound = errors.New("token not found")
)

func (s *Storage) RefreshTokenByToken(ctx context.Context, refresh_token string) (models.Token, error) {
	const op = "storage.auth.GetRefreshTokenByToken"
	var token models.Token

	// Запрос для получения refresh-токена по user_id
	query := "SELECT id, user_id, refresh_token FROM tokens WHERE refresh_token = $1"
	err := s.db.QueryRowContext(ctx, query, refresh_token).Scan(
		&token.ID,
		&token.UserID,
		&token.RefreshToken,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Token{}, ErrTokenNotFound
		}
		return models.Token{}, fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (s *Storage) RefreshTokenByUserId(ctx context.Context, userID int64) (models.Token, error) {
	const op = "storage.auth.GetRefreshToken"
	var token models.Token

	// Запрос для получения refresh-токена по user_id
	query := "SELECT id, user_id, refresh_token FROM tokens WHERE user_id = $1"
	err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&token.ID,
		&token.UserID,
		&token.RefreshToken,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Token{}, ErrTokenNotFound
		}
		return models.Token{}, fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (s *Storage) SaveRefreshToken(ctx context.Context, userID int64, token string) error {
	const op = "storage.auth.SaveRefreshToken"

	// Запрос на сохранение нового refresh-токена
	query := "INSERT INTO tokens (user_id, refresh_token) VALUES ($1, $2)"
	_, err := s.db.ExecContext(ctx, query, userID, token)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	const op = "storage.auth.DeleteRefreshToken"

	// Запрос на удаление refresh-токена
	query := "DELETE FROM tokens WHERE refresh_token = $1"
	result, err := s.db.ExecContext(ctx, query, refreshToken)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Проверяем, сколько строк было удалено
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return ErrTokenNotFound
	}

	return nil
}

func (s *Storage) UpdateRefreshToken(ctx context.Context, userID int64, token string) error {
	const op = "storage.auth.UpdateRefreshToken"

	// Запрос на обновление refresh-токена
	query := "UPDATE tokens SET refresh_token = $2 WHERE user_id = $1"
	result, err := s.db.ExecContext(ctx, query, userID, token)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Проверка, был ли обновлен хотя бы один пользователь
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return ErrTokenNotFound
	}

	return nil
}
