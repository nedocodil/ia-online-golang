package middleware

import (
	"context"
	"ia-online-golang/internal/http/context_keys"
	"ia-online-golang/internal/http/responses"
	"ia-online-golang/internal/services/token"
	"net/http"
	"strings"
)

func JWTMiddleware(ctx context.Context, tokenService token.TokenServiceI) func(http.Handler) http.Handler {
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

			tokenAccess := parts[1]

			var payload token.PayloadUserAccess
			claims, err := tokenService.ValidateAccessToken(ctx, tokenAccess, &payload)
			if err != nil {
				responses.InvalidAccessToken(w)
				return
			}

			userClaims, ok := claims.(*token.PayloadUserAccess)
			if !ok {
				responses.InvalidAccessToken(w)
				return
			}

			// Добавляем userID и роли в контекст
			ctx := context.WithValue(r.Context(), context_keys.UserIDKey, userClaims.UserID)
			ctx = context.WithValue(ctx, context_keys.UserRoleKey, userClaims.Roles)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RoleMiddleware(requiredRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			userRolesInterface := r.Context().Value(context_keys.UserRoleKey)
			if userRolesInterface == nil {
				responses.Forbidden(w)
				return
			}

			userRoles, ok := userRolesInterface.([]string)
			if !ok {
				responses.Forbidden(w)
				return
			}

			// Проверяем, есть ли у пользователя хотя бы одна из требуемых ролей
			for _, userRole := range userRoles {
				for _, requiredRole := range requiredRoles {
					if userRole == requiredRole {
						next.ServeHTTP(w, r)
						return
					}
				}
			}

			responses.Forbidden(w)
		})
	}
}
