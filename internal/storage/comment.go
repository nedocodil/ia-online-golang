package storage

import (
	"context"
	"fmt"
	"ia-online-golang/internal/models"
)

type CommentsRepositoryI interface {
	SaveComment(ctx context.Context, comment models.Comment) error
}

// var (
// 	ErrLeadNotFound = errors.New("comment not found")
// )

// func (s *Storage) Comment(ctx context.Context) (models.Comment, error) {
// 	const op = "storage.user.UserByEmail"
// 	var user models.User

// 	query := "SELECT id, email, name, phone_number, telegram, is_active, created_at, city, password_hash, referral_code, roles, reward_internet, reward_cleaning, reward_shipping FROM users WHERE email = $1"
// 	err := s.db.QueryRowContext(ctx, query, email).Scan(
// 		&user.ID,
// 		&user.Email,
// 		&user.Name,
// 		&user.PhoneNumber,
// 		&user.Telegram,
// 		&user.IsActive,
// 		&user.CreatedAt,
// 		&user.City,
// 		&user.PasswordHash,
// 		&user.ReferralCode,
// 		&user.Roles,
// 		&user.RewardInternet,
// 		&user.RewardCleaning,
// 		&user.RewardShipping,
// 	)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return user, ErrUserNotFound
// 		}
// 		return user, fmt.Errorf("%s: %w", op, err)
// 	}

// 	return user, nil
// }

func (s *Storage) SaveComment(ctx context.Context, comment models.Comment) error {
	const op = "CommentRepository.SaveComment"

	query := `INSERT INTO comments (lead_id, user_id, text) VALUES ($1, $2, $3)`
	_, err := s.db.ExecContext(ctx, query, comment.LeadID, comment.UserID, comment.Text)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
