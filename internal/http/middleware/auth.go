package middleware

import (
	"context"
	"ia-online-golang/internal/http/responses"
	payloads "ia-online-golang/internal/interfaces/payload"
	"ia-online-golang/internal/services/token"
	"net/http"
	"strings"
)

// Определяем собственный тип ключа для контекста
type contextKey string

const userIDKey contextKey = "userID"

func JWTMiddleware(ctx context.Context, tokenService *token.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				responses.AccessTokenNotFound(w)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				responses.InvalidBearerFormat(w)
				return
			}

			token := parts[1]

			var payload payloads.PayloadUserAccess
			claims, err := tokenService.ValidateAccessToken(ctx, token, &payload)
			if err != nil {
				responses.InvalidAccessToken(w)
				return
			}
			// Проверяем, содержит ли claims нужный UserID
			userClaims, ok := claims.(*payloads.PayloadUserAccess)
			if !ok {
				responses.InvalidAccessToken(w)
				return
			}

			// Добавляем userID в контекст (если есть)
			ctx := context.WithValue(r.Context(), userIDKey, userClaims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
