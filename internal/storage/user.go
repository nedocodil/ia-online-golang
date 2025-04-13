package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ia-online-golang/internal/models"

	"github.com/lib/pq"
)

type UserRepositoryI interface {
	UserByReferralCode(ctx context.Context, referral_code string) (models.User, error)
	Users(ctx context.Context) ([]models.User, error)
	UserByEmail(ctx context.Context, email string) (models.User, error)
	UserById(ctx context.Context, id int64) (models.User, error)
	UserIdByEmail(ctx context.Context, email string) (int64, error)
	UserIdByPhone(ctx context.Context, phone string) (int64, error)
	ValidationUser(ctx context.Context, email string, phone string, telegram string) error
	UserIdByTelegram(ctx context.Context, telegram string) (int64, error)
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	UpdateActiveUser(ctx context.Context, userID int64, isActive bool) error
	UpdatePasswordUser(ctx context.Context, password_hash string, userID int64) error
	UpdateUser(ctx context.Context, user models.User) error
	DeleteUser(ctx context.Context, id int) error
}

var (
	ErrUserExists       = errors.New("user already exists")
	ErrUserNotFound     = errors.New("user not found")
	ErrUserIsNotUpdated = errors.New("user is not updated")
)

// Получение пользователя по email
func (s *Storage) UserByEmail(ctx context.Context, email string) (models.User, error) {
	const op = "storage.user.UserByEmail"
	var user models.User

	query := "SELECT id, email, name, phone_number, telegram, is_active, created_at, city, password_hash, referral_code, roles FROM users WHERE email = $1"
	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PhoneNumber,
		&user.Telegram,
		&user.IsActive,
		&user.CreatedAt,
		&user.City,
		&user.PasswordHash,
		&user.ReferralCode,
		&user.Roles,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, ErrUserNotFound
		}
		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) Users(ctx context.Context) ([]models.User, error) {
	const op = "storage.user.Users"
	var users []models.User

	query := "SELECT id, email, name, phone_number, telegram, is_active, created_at, city, password_hash, referral_code, roles FROM users"
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID, &user.Email, &user.Name, &user.PhoneNumber, &user.Telegram,
			&user.IsActive, &user.CreatedAt, &user.City, &user.PasswordHash,
			&user.ReferralCode, &user.Roles,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}

func (s *Storage) UserByReferralCode(ctx context.Context, referral_code string) (models.User, error) {
	const op = "storage.user.UserByReferralCode"
	var user models.User

	query := "SELECT id, email, name, phone_number, telegram, is_active, created_at, city, password_hash, referral_code, roles FROM users WHERE referral_code = $1"
	err := s.db.QueryRowContext(ctx, query, referral_code).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PhoneNumber,
		&user.Telegram,
		&user.IsActive,
		&user.CreatedAt,
		&user.City,
		&user.PasswordHash,
		&user.ReferralCode,
		&user.Roles,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, ErrUserNotFound
		}
		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) UserById(ctx context.Context, id int64) (models.User, error) {
	const op = "storage.user.UserById"
	var user models.User

	query := "SELECT id, email, name, phone_number, telegram, is_active, created_at, city, password_hash, referral_code, roles FROM users WHERE id = $1"
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PhoneNumber,
		&user.Telegram,
		&user.IsActive,
		&user.CreatedAt,
		&user.City,
		&user.PasswordHash,
		&user.ReferralCode,
		&user.Roles,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, ErrUserNotFound
		}
		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) UserIdByEmail(ctx context.Context, email string) (int64, error) {
	const op = "storage.user.UserIdByEmail"
	var userId int64

	query := "SELECT id FROM users WHERE email = $1"
	err := s.db.QueryRowContext(ctx, query, email).Scan(&userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrUserNotFound
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userId, nil
}

func (s *Storage) UserIdByPhone(ctx context.Context, phone string) (int64, error) {
	const op = "storage.user.UserIdByPhone"
	var userId int64

	query := "SELECT id FROM users WHERE phone_number = $1"
	err := s.db.QueryRowContext(ctx, query, phone).Scan(&userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrUserNotFound
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userId, nil
}

func (s *Storage) UserIdByTelegram(ctx context.Context, telegram string) (int64, error) {
	const op = "storage.user.UserIdByTelegram"
	var userId int64

	query := "SELECT id FROM users WHERE telegram = $1"
	err := s.db.QueryRowContext(ctx, query, telegram).Scan(&userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrUserNotFound
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userId, nil
}

func (s *Storage) ValidationUser(ctx context.Context, email string, phone string, telegram string) error {
	const op = "storage.user.ValidationUser"

	var count int
	query := `SELECT COUNT(*) FROM users WHERE phone_number = $1 OR email = $2 OR telegram = $3`

	err := s.db.QueryRowContext(ctx, query, phone, email, telegram).Scan(&count)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if count > 0 {
		return ErrUserExists
	}

	return nil
}

func (s *Storage) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	const op = "storage.user.CreateUser"

	// Правильный SQL-запрос для PostgreSQL, который возвращает все поля пользователя
	query := `INSERT INTO users (email, password_hash, phone_number, name, telegram, city, referral_code, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, roles, email, password_hash, phone_number, name, telegram, city, referral_code, is_active`

	var newUser models.User
	// Извлекаем все данные о пользователе
	err := s.db.QueryRowContext(ctx, query, user.Email, user.PasswordHash, user.PhoneNumber, user.Name, user.Telegram, user.City, user.ReferralCode, user.IsActive).
		Scan(&newUser.ID, &newUser.Roles, &newUser.Email, &newUser.PasswordHash, &newUser.PhoneNumber, &newUser.Name, &newUser.Telegram, &newUser.City, &newUser.ReferralCode, &newUser.IsActive)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return newUser, nil
}

func (s *Storage) UpdateActiveUser(ctx context.Context, userID int64, isActive bool) error {
	const op = "storage.user.UpdateActiveUser"

	// Пытаемся обновить пользователя
	query := "UPDATE users SET is_active = $1 WHERE id = $2"
	result, err := s.db.ExecContext(ctx, query, isActive, userID)
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
		return ErrUserNotFound
	}

	return nil
}

func (s *Storage) UpdateUser(ctx context.Context, user models.User) error {
	const op = "storage.user.UpdateUser"

	// Карта для хранения полей и их значений
	updateFields := make(map[string]interface{})
	if user.Email != "" {
		updateFields["email"] = user.Email
	}
	if user.PasswordHash != "" {
		updateFields["password"] = user.PasswordHash
	}
	if user.Name != "" {
		updateFields["name"] = user.Name
	}
	if user.City != "" {
		updateFields["city"] = user.City
	}
	if user.Telegram != "" {
		updateFields["telegram"] = user.Telegram
	}
	if user.PhoneNumber != "" {
		updateFields["phone_number"] = user.PhoneNumber
	}
	if len(user.Roles) > 0 {
		updateFields["roles"] = pq.Array(user.Roles)
	}

	// Строим динамический запрос
	query := "UPDATE users SET "
	var args []interface{}
	argPos := 1

	for field, value := range updateFields {
		query += fmt.Sprintf("%s = $%d, ", field, argPos)
		args = append(args, value)
		argPos++
	}

	// Убираем последнюю запятую
	query = query[:len(query)-2]

	// Добавляем условие WHERE
	query += fmt.Sprintf(" WHERE id = $%d", argPos)
	args = append(args, user.ID)

	// Выполняем запрос
	res, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// Проверяем, обновилась ли хоть одна строка
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s: %w", op, ErrUserIsNotUpdated)
	}

	return nil
}

func (s *Storage) UpdatePasswordUser(ctx context.Context, password_hash string, userID int64) error {
	const op = "storage.user.UpdatePasswordUser"

	// Пытаемся обновить пользователя
	query := "UPDATE users SET password_hash = $1 WHERE id = $2"
	result, err := s.db.ExecContext(ctx, query, password_hash, userID)
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
		return ErrUserNotFound
	}

	return nil
}

func (s *Storage) DeleteUser(ctx context.Context, id int) error {
	const op = "storage.user.DeleteUser"

	query := "DELETE FROM users WHERE id = $1"
	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
