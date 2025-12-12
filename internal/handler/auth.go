package handler

import (
	"encoding/json"
	"net/http"
	
	"sport-manager/internal/service"
)

// LoginRequest определяет структуру для входящего запроса логина
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthHandler содержит методы для обработки HTTP-запросов аутентификации
type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

// Login обрабатывает запрос POST /login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		// Возвращаем 401 Unauthorized для неверных учетных данных
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}