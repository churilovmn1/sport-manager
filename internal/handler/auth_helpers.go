package handler

import (
	"net/http"
	"sport-manager/internal/auth"
)

// checkAdminRole извлекает роль из контекста и проверяет, является ли пользователь админом.
// Возвращает false, если роль не найдена или не 'admin'.
func checkAdminRole(r *http.Request) bool {
	// Извлекаем роль из контекста, используя ключ из пакета auth
	role, ok := r.Context().Value(auth.ContextKeyRole).(string)
	if !ok {
		return false // Роль не найдена
	}
	return role == "admin"
}

// requireAdmin проверяет права администратора.
// Если пользователь не администратор, отправляет ответ 403 Forbidden и возвращает false.
func requireAdmin(w http.ResponseWriter, r *http.Request) bool {
	if !checkAdminRole(r) {
		http.Error(w, "Access Denied. Only administrators can perform this action.", http.StatusForbidden)
		return false
	}
	return true
}
