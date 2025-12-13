package handler

import (
	"encoding/json"
	"log" // Добавлен для логгирования
	"net/http"
	
	"sport-manager/internal/service"
)

// LoginRequest определяет структуру для входящего запроса логина
type LoginRequest struct {
	Username string `json:"username"` // <-- ИСПРАВЛЕНО: Добавлен JSON-тег
	Password string `json:"password"` // <-- ИСПРАВЛЕНО: Добавлен JSON-тег
}

// AuthHandler содержит методы для обработки HTTP-запросов аутентификации
type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

// Login обрабатывает запрос POST /api/v1/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	
	// 1. Декодирование JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// 2. ДИАГНОСТИЧЕСКИЙ ЛОГ
	log.Printf("DEBUG: Login attempt for user: %s, password length: %d", req.Username, len(req.Password))
	if req.Username == "" || req.Password == "" {
        log.Println("DEBUG: Username or Password field is empty after decoding. Check JSON tags in LoginRequest.")
    }
	// ----------------------

	// 3. Аутентификация
	token, err := h.service.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		// В идеале: вернуть JSON-ошибку, а не простой текст
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid credentials"})
		return
	}

	// 4. Успешный ответ
	// Возвращаем access_token (соответствует ожиданиям JS-фронтенда)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"access_token": token})
}