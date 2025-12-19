package handler

import (
	"encoding/json"
	"net/http"
	"sport-manager/internal/service"
)

// --- СТРУКТУРЫ ДАННЫХ ДЛЯ ЗАПРОСОВ ---

// LoginRequest описывает входящие данные для входа
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterRequest описывает входящие данные для создания аккаунта
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthHandler отвечает за обработку HTTP-запросов, связанных с аутентификацией
type AuthHandler struct {
	service *service.AuthService
}

// NewAuthHandler создает новый экземпляр хендлера (Dependency Injection)
func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

// --- МЕТОДЫ ОБРАБОТКИ ---

// Login выполняет аутентификацию пользователя и выдает JWT токен
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	// Декодируем тело JSON-запроса в структуру
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Некорректный формат запроса")
		return
	}

	// Вызываем бизнес-логику из сервиса
	token, err := h.service.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		// Если пароль неверный или пользователь не найден
		h.respondWithError(w, http.StatusUnauthorized, "Неверный логин или пароль")
		return
	}

	// Возвращаем токен в случае успеха
	h.respondWithJSON(w, http.StatusOK, map[string]string{
		"access_token": token,
		"status":       "success",
	})
}

// Register создает нового пользователя в системе
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Некорректный формат запроса")
		return
	}

	// Передаем данные в слой сервиса для регистрации
	err := h.service.Register(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Ошибка при создании пользователя")
		return
	}

	// Успешный ответ (201 Created)
	h.respondWithJSON(w, http.StatusCreated, map[string]string{"status": "success"})
}

// --- ВСПОМОГАТЕЛЬНЫЕ МЕТОДЫ (Утилиты для чистоты кода) ---

// respondWithError упрощает отправку сообщений об ошибках в формате JSON
func (h *AuthHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]string{"error": message})
}

// respondWithJSON унифицирует отправку успешных ответов и установку заголовков
func (h *AuthHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
