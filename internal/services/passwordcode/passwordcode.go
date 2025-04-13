package passwordcode

import (
	"errors"
	"fmt"
	"ia-online-golang/internal/models"
	"ia-online-golang/internal/storage"
	"ia-online-golang/internal/utils"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type PasswordCodeService struct {
	log                    *logrus.Logger
	PasswordCodeRepository storage.PasswordCodeRepositoryI
}

type PasswordCodeServiceI interface {
	PasswordCode(ctx context.Context, password_code string) (models.PasswordCode, error)
	GeneratePasswordCode(ctx context.Context, userID int64) (string, error)
	DeleteRefreshToken(ctx context.Context, code string) error
}

var (
	ErrPasswordCodeIsNotFound = errors.New("password code is not found")
	ErrPasswordCodeIncorrect  = errors.New("password code incorrect")
	ErrPasswordCodeHasExpired = errors.New("password code has expired")
)

func New(log *logrus.Logger, passwordCodeRepository storage.PasswordCodeRepositoryI) *PasswordCodeService {
	return &PasswordCodeService{
		log:                    log,
		PasswordCodeRepository: passwordCodeRepository,
	}
}

func (p *PasswordCodeService) PasswordCode(ctx context.Context, password_code string) (models.PasswordCode, error) {
	op := "PasswordCodeService.PasswordCode"

	code, err := p.PasswordCodeRepository.PasswordCode(ctx, password_code)
	if err != nil {
		if errors.Is(err, storage.ErrPasswordCodeIsNotFound) {
			return models.PasswordCode{}, ErrPasswordCodeIsNotFound
		}

		return models.PasswordCode{}, fmt.Errorf("%s:%w", op, err)
	}

	return code, nil
}
func (p *PasswordCodeService) GeneratePasswordCode(ctx context.Context, userID int64) (string, error) {
	op := "PasswordCodeService.GeneratePasswordCode"

	now := time.Now()
	tomorrow := now.Add(24 * time.Hour)

	newCode, err := utils.GeneratePasswordCode(6)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	result := models.PasswordCode{
		UserID:    userID,
		Code:      newCode,
		ExpiresAt: tomorrow,
	}

	password_code, err := p.PasswordCodeRepository.PasswordCodeByUserId(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrPasswordCodeIsNotFound) {

			err := p.PasswordCodeRepository.SavePasswordCode(ctx, result)
			if err != nil {
				return "", fmt.Errorf("%s: %w", op, err)
			}

			return result.Code, nil
		}

		return "", fmt.Errorf("%s:%w", op, err)
	}

	if now.After(password_code.ExpiresAt) {
		err := p.PasswordCodeRepository.UpdatePasswordCode(ctx, result)
		if err != nil {
			return "", fmt.Errorf("%s:%w", op, err)
		}

		return result.Code, nil
	}

	return password_code.Code, nil
}
func (p *PasswordCodeService) DeleteRefreshToken(ctx context.Context, code string) error {
	op := "PasswordCodeService.GeneratePasswordCode"

	err := p.PasswordCodeRepository.DeletePasswordCode(ctx, code)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
