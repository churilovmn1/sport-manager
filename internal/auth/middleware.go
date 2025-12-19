package auth

import (
	"context"
	"net/http"
	"strings"

	"sport-manager/pkg/config"
)

// ContextKey — специальный тип для ключей контекста.
// Это защищает нас от случайных совпадений (коллизий), если другие пакеты тоже используют контекст.
type ContextKey string

const (
	ContextKeyRole     ContextKey = "userRole"
	ContextKeyUsername ContextKey = "userUsername"
)

// AuthMiddleware — основной фильтр (посредник), который проверяет JWT токен.
// Он выполняется ДО того, как запрос попадет в хендлер.
func AuthMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Извлекаем заголовок Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
				return
			}

			// 2. Проверяем формат заголовка (должен быть: Bearer <token>)
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Неверный формат токена", http.StatusUnauthorized)
				return
			}
			tokenString := parts[1]

			// 3. Валидируем токен (проверка подписи и срока годности)
			claims, err := ValidateToken(tokenString, cfg)
			if err != nil {
				http.Error(w, "Токен недействителен или просрочен", http.StatusUnauthorized)
				return
			}

			// 4. Передаем данные пользователя дальше по цепочке через Context.
			// Контекст позволяет хендлерам узнать, кто делает запрос, не перечитывая токен.
			ctx := context.WithValue(r.Context(), ContextKeyRole, claims.Role)
			ctx = context.WithValue(ctx, ContextKeyUsername, claims.Username)

			// Передаем управление следующему обработчику с обновленным контекстом
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// AdminOnly — Middleware для разграничения прав доступа (RBAC).
// Он пропускает запрос дальше только в том случае, если у пользователя роль "admin".
func AdminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем роль, которую ранее сохранил AuthMiddleware
		role, ok := r.Context().Value(ContextKeyRole).(string)

		// Если роли нет или она не админская — возвращаем 403 Forbidden
		if !ok || role != "admin" {
			http.Error(w, "Доступ запрещен: требуются права администратора", http.StatusForbidden)
			return
		}

		// Если всё в порядке, выполняем основной код хендлера
		next.ServeHTTP(w, r)
	}
}
