package handler

import (
	"encoding/json"
	"log" // <-- Добавлен для логгирования ошибок
	"net/http"
	"strconv"
	"strings"

	"sport-manager/internal/repository"
	"sport-manager/internal/service"

	"github.com/gorilla/mux"
)

// AthleteHandler обрабатывает HTTP-запросы для спортсменов
type AthleteHandler struct {
	service *service.AthleteService
}

func NewAthleteHandler(s *service.AthleteService) *AthleteHandler {
	return &AthleteHandler{service: s}
}

// helper
func writeJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// ----------------------------------------------------
// 1. POST /api/v1/athletes (Create)
// ----------------------------------------------------
func (h *AthleteHandler) CreateAthlete(w http.ResponseWriter, r *http.Request) {
	var athlete repository.Athlete
	if err := json.NewDecoder(r.Body).Decode(&athlete); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// --- ИСПРАВЛЕНИЕ 1: Инициализация обязательных полей ---
	// Фронтенд не присылает Gender. Если это обязательное поле в БД, нужно дать заглушку.
	if athlete.Gender == "" {
		athlete.Gender = "Не указан" 
	}
	if !athlete.IsActive {
        athlete.IsActive = true // Предполагаем, что новый спортсмен активен
    }
	// -------------------------------------------------------

	if err := h.service.Create(r.Context(), &athlete); err != nil {
		// --- ИСПРАВЛЕНИЕ 2: ДИАГНОСТИЧЕСКИЙ ЛОГ ---
		log.Printf("ERROR: Athlete creation failed. Input: %+v. Server Error: %v", athlete, err)
		// ------------------------------------------

		if strings.Contains(err.Error(), "validation error") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Улучшенный ответ для фронтенда, чтобы намекнуть на проблему с внешним ключом
		http.Error(w, "Ошибка создания спортсмена на сервере. Возможно, неверный ID спорта/разряда.", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusCreated, athlete)
}

// ----------------------------------------------------
// 2. GET /api/v1/athletes/{id} (Read One)
// ----------------------------------------------------
func (h *AthleteHandler) GetAthlete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid athlete ID", http.StatusBadRequest)
		return
	}

	athlete, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Athlete not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error retrieving athlete", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusOK, athlete)
}

// ----------------------------------------------------
// 3. GET /api/v1/athletes (Read All)
// ----------------------------------------------------
func (h *AthleteHandler) ListAthletes(w http.ResponseWriter, r *http.Request) {
	athletes, err := h.service.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Internal server error listing athletes", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusOK, athletes)
}

// ----------------------------------------------------
// 4. PUT /api/v1/athletes/{id} (Update)
// ----------------------------------------------------
func (h *AthleteHandler) UpdateAthlete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid athlete ID", http.StatusBadRequest)
		return
	}

	var athlete repository.Athlete
	if err := json.NewDecoder(r.Body).Decode(&athlete); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	athlete.ID = id // Убеждаемся, что обновляем правильный ID

	if err := h.service.Update(r.Context(), &athlete); err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Athlete not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error updating athlete", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusOK, athlete)
}

// ----------------------------------------------------
// 5. DELETE /api/v1/athletes/{id} (Delete)
// ----------------------------------------------------
func (h *AthleteHandler) DeleteAthlete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid athlete ID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Athlete not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error deleting athlete", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusNoContent, nil) // 204 No Content
}