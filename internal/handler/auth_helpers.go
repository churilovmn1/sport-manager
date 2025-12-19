package handler

import (
	"net/http"
	"sport-manager/internal/auth"
)

// checkAdminRole — внутренняя вспомогательная функция.
// Она достает роль пользователя из контекста запроса (куда её ранее положил AuthMiddleware).
func checkAdminRole(r *http.Request) bool {
	// r.Context().Value позволяет получить данные, привязанные к текущему запросу.
	// Приводим интерфейс к строке через .(string)
	role, ok := r.Context().Value(auth.ContextKeyRole).(string)
	if !ok {
		return false // Если токен был пустой или роль не указана
	}
	return role == "admin"
}

// requireAdmin — функция-страж.
// Используется внутри хендлеров, чтобы пресечь действия не-админов.
func requireAdmin(w http.ResponseWriter, r *http.Request) bool {
	if !checkAdminRole(r) {
		// Если не админ — сразу отдаем статус 403 (Доступ запрещен)
		http.Error(w, "Доступ запрещен. Это действие доступно только администраторам.", http.StatusForbidden)
		return false
	}
	// Если админ — возвращаем true, и хендлер продолжает работу
	return true
}
