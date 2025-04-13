package token

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"ia-online-golang/internal/dto"
	"ia-online-golang/internal/services/lead"
	"ia-online-golang/internal/services/referral"
	"ia-online-golang/internal/services/user"
	"ia-online-golang/internal/storage"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

type TokenService struct {
	log                   *logrus.Logger
	SecretKeyAccess       string
	SecretKeyRefresh      string
	ExpirationTimeAccess  int64
	ExpirationTimeRefresh int64
	TokenRepository       storage.TokenRepositoryI
	UserService           user.UserServiceI
	LeadService           lead.LeadServiceI
	ReferralService       referral.ReferralServiceI
}

type TokenServiceI interface {
	CreateUserTokens(ctx context.Context, userID int64) (dto.AuthTokensDTO, error)
	GenerateTokens(ctx context.Context, payloadAccess any, payloadRefresh any) (dto.AuthTokensDTO, error)
	GenerateAccessToken(ctx context.Context, payloadAccess any) (string, error)
	SaveToken(ctx context.Context, userId int64, refreshToken string) error
	ValidateRefreshToken(ctx context.Context, refresh_token string, payloadStruct any) (any, error)
	ValidateAccessToken(ctx context.Context, token string, payloadStruct any) (any, error)
}

var (
	ErrSaveRefreshToken          = errors.New("error saving refresh token")
	ErrRefreshTokenAlreadyExists = errors.New("refreshing token already exists")
	ErrRefreshTokenNotExists     = errors.New("refreshing token not exists")
	ErrInvalidToken              = errors.New("token is invalid")
	ErrInvalidRefreshToken       = errors.New("refresh token is invalid")
	ErrInvalidAccessToken        = errors.New("access token is invalid")
	ErrExpiredToken              = errors.New("expired token")
	ErrExpiredRefreshToken       = errors.New("expired access token")
	ErrExpiredAccessToken        = errors.New("expired refresh token")
)

func New(log *logrus.Logger,
	secretKeyAccess string,
	secretKeyRefresh string,
	expiryTimeAccess int64,
	expiryTimeRefresh int64,
	tokenRepository storage.TokenRepositoryI,
	userService user.UserServiceI,
	leadService lead.LeadServiceI,
	referralService referral.ReferralServiceI) *TokenService {
	return &TokenService{
		log:                   log,
		SecretKeyAccess:       secretKeyAccess,
		SecretKeyRefresh:      secretKeyRefresh,
		ExpirationTimeAccess:  expiryTimeAccess,
		ExpirationTimeRefresh: expiryTimeRefresh,
		TokenRepository:       tokenRepository,
		UserService:           userService,
		LeadService:           leadService,
		ReferralService:       referralService,
	}
}

func (s *TokenService) CreateUserTokens(ctx context.Context, userID int64) (dto.AuthTokensDTO, error) {
	op := "TokenService.CreateUserTokens"

	// Получаем пользователя по `UserID`
	user, err := s.UserService.UserById(ctx, userID)
	if err != nil {
		return dto.AuthTokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	now := time.Now()
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	statistic, err := s.LeadService.GetUserPaymentStatistic(ctx, userID, &firstOfMonth, nil)
	if err != nil {
		return dto.AuthTokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	referrals, err := s.ReferralService.ReferralsUser(ctx, user.ReferralCode)
	if err != nil {
		return dto.AuthTokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	payloadAccess := PayloadUserAccess{
		UserID:       *user.ID,
		Roles:        user.Roles,
		Name:         user.Name,
		Email:        user.Email,
		PhoneNumber:  user.PhoneNumber,
		City:         user.City,
		Telegram:     user.Telegram,
		ReferralCode: user.ReferralCode,
		Statistic:    statistic,
		Referrals:    referrals,
	}
	payloadRefresh := PayloadUserRefresh{
		UserID: *user.ID,
	}

	tokens, err := s.GenerateTokens(ctx, payloadAccess, payloadRefresh)
	if err != nil {
		s.log.Error(err)

		return dto.AuthTokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	err = s.SaveToken(ctx, *user.ID, tokens.RefreshToken)
	if err != nil {
		s.log.Error(err)

		return dto.AuthTokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	return tokens, nil
}

func (t *TokenService) GenerateTokens(ctx context.Context, payloadAccess any, payloadRefresh any) (dto.AuthTokensDTO, error) {
	op := "TokenService.GenerateTokens"

	payloadMapAccess, err := t.structToMap(payloadAccess)
	if err != nil {
		return dto.AuthTokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	payloadMapRefresh, err := t.structToMap(payloadRefresh)
	if err != nil {
		return dto.AuthTokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Создаем access-токен
	accessToken, err := t.createToken(payloadMapAccess, t.ExpirationTimeAccess, t.SecretKeyAccess)
	if err != nil {
		return dto.AuthTokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	// Создаем refresh-токен
	refreshToken, err := t.createToken(payloadMapRefresh, t.ExpirationTimeRefresh, t.SecretKeyRefresh)
	if err != nil {
		return dto.AuthTokensDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	return dto.AuthTokensDTO{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (t *TokenService) GenerateAccessToken(ctx context.Context, payloadAccess any) (string, error) {
	op := "TokenService.GenerateAccessToken"

	payloadMapAccess, err := t.structToMap(payloadAccess)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// Создаем access-токен
	accessToken, err := t.createToken(payloadMapAccess, t.ExpirationTimeAccess, t.SecretKeyAccess)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return accessToken, nil
}

func (t *TokenService) SaveToken(ctx context.Context, userId int64, refreshToken string) error {
	op := "TokenService.SaveToken"

	_, err := t.TokenRepository.RefreshTokenByUserId(ctx, userId)
	if err != nil {
		if errors.Is(err, storage.ErrTokenNotFound) {
			// Токена нет — сохраняем новый
			if err := t.TokenRepository.SaveRefreshToken(ctx, userId, refreshToken); err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
			return nil
		}
		// Любая другая ошибка — возвращаем её сразу
		return fmt.Errorf("%s: %w", op, err)
	}

	// Токен уже есть — обновляем
	if err := t.TokenRepository.UpdateRefreshToken(ctx, userId, refreshToken); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (t *TokenService) ValidateRefreshToken(ctx context.Context, refresh_token string, payloadStruct any) (any, error) {
	op := "TokenService.ValidateRefreshToken"

	payload, err := t.validateToken(refresh_token, t.SecretKeyRefresh, payloadStruct)
	if err != nil {
		if errors.Is(err, ErrInvalidToken) {
			return nil, ErrInvalidRefreshToken
		}

		if errors.Is(err, ErrExpiredToken) {
			return nil, ErrExpiredRefreshToken
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = t.TokenRepository.RefreshTokenByToken(ctx, refresh_token)
	if err != nil {
		if errors.Is(err, storage.ErrTokenNotFound) {
			return nil, storage.ErrTokenNotFound
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return payload, nil
}

func (t *TokenService) ValidateAccessToken(ctx context.Context, token string, payloadStruct any) (any, error) {
	op := "TokenService.ValidateAccessToken"

	payload, err := t.validateToken(token, t.SecretKeyAccess, payloadStruct)
	if err != nil {
		if errors.Is(err, ErrInvalidToken) {
			return nil, ErrInvalidAccessToken
		}

		if errors.Is(err, ErrExpiredToken) {
			return nil, ErrExpiredAccessToken
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return payload, nil
}

func (t *TokenService) createToken(payload map[string]interface{}, expirationTime int64, secretKey string) (string, error) {
	op := "TokenService.createToken"

	claims := jwt.MapClaims{}
	for key, value := range payload {
		claims[key] = value
	}
	claims["exp"] = time.Now().Add(time.Duration(expirationTime) * time.Second).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return tokenString, nil
}

func (t *TokenService) validateToken(tokenString string, secretKey string, payloadStruct any) (any, error) {
	op := "TokenService.validateToken"

	// Парсим токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, fmt.Errorf("%s: error method %s", op, token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Проверяем валидность токена
	if !token.Valid {
		return nil, ErrInvalidToken
	}

	// Извлекаем claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("%s: invalid claims format", op)
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("%s: missing exp field", op)
	}

	// Проверяем, не истёк ли срок действия
	if time.Now().Unix() > int64(exp) {
		return nil, ErrExpiredToken
	}

	// Конвертируем `claims` в JSON и затем в структуру
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = json.Unmarshal(claimsJSON, payloadStruct)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return payloadStruct, nil
}

func (t *TokenService) structToMap(s interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
