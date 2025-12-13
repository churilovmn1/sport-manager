package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"sport-manager/internal/service" // Используем Claims из сервиса
	"sport-manager/pkg/config"

	"github.com/golang-jwt/jwt/v5"
)

// Контекстный ключ для UserID
type ContextKey string

const ContextUserID ContextKey = "userID"

// Middleware для проверки JWT-токена
func AuthMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// 1. Получение токена из заголовка Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// Ожидаемый формат: "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid token format", http.StatusUnauthorized)
				return
			}
			tokenString := parts[1]

			// 2. Парсинг и валидация токена
			claims := &service.Claims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				// Проверка метода подписи
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(cfg.JWTSecret), nil
			})

			if err != nil || !token.Valid {
				// Проверка на истечение срока действия токена
				if errors.Is(err, jwt.ErrTokenExpired) || (claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now())) {
					http.Error(w, "Token expired", http.StatusUnauthorized)
					return
				}
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// 3. Сохранение данных пользователя в контексте
			// Передаем UserID, чтобы Handler знал, кто делает запрос
			ctx := context.WithValue(r.Context(), ContextUserID, claims.UserID)

			// 4. Передача управления следующему обработчику
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
