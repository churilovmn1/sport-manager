package auth

import (
	"context"
	"net/http"
	"strings"

	"sport-manager/pkg/config"
)

// ContextKey - собственный тип для ключей контекста
type ContextKey string

// ContextKeyRole - ключ для хранения роли пользователя в контексте
const ContextKeyRole ContextKey = "userRole" // <-- Добавлено

// AuthMiddleware проверяет JWT и сохраняет роль в контексте
func AuthMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Получение токена из заголовка Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			// Ожидаем формат "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid token format", http.StatusUnauthorized)
				return
			}
			tokenString := parts[1]

			// 2. Валидация токена
			claims, err := ValidateToken(tokenString, cfg)
			if err != nil {
				http.Error(w, "Invalid or expired token: "+err.Error(), http.StatusUnauthorized)
				return
			}

			// 3. Сохранение роли в контексте запроса
			ctx := context.WithValue(r.Context(), ContextKeyRole, claims.Role)
			r = r.WithContext(ctx)

			// Продолжение обработки запроса
			next.ServeHTTP(w, r)
		})
	}
}
