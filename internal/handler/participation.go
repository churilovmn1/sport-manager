package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"sport-manager/internal/repository"
	"sport-manager/internal/service"

	"github.com/gorilla/mux"
)

// ParticipationHandler обрабатывает HTTP-запросы для участия/результатов
type ParticipationHandler struct {
	service *service.ParticipationService
}

func NewParticipationHandler(s *service.ParticipationService) *ParticipationHandler {
	return &ParticipationHandler{service: s}
}

// ----------------------------------------------------
// 1. POST /api/v1/participations (Регистрация)
// ----------------------------------------------------
func (h *ParticipationHandler) CreateParticipation(w http.ResponseWriter, r *http.Request) {
	var p repository.Participation
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.Create(r.Context(), &p); err != nil {
		if strings.Contains(err.Error(), "validation error") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error registering athlete", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusCreated, p)
}

// ----------------------------------------------------
// 2. GET /api/v1/participations (Список всех участий)
// ----------------------------------------------------
func (h *ParticipationHandler) ListParticipations(w http.ResponseWriter, r *http.Request) {
	details, err := h.service.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Internal server error listing participations", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusOK, details)
}

// ----------------------------------------------------
// 3. PUT /api/v1/participations/{id}/place (Проставление результата)
// ----------------------------------------------------
// Ожидаемый JSON: {"place": 1}
func (h *ParticipationHandler) UpdatePlace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid participation ID", http.StatusBadRequest)
		return
	}

	var data struct {
		Place int `json:"place"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body: must contain 'place'", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdatePlace(r.Context(), id, data.Place); err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Participation not found", http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "validation error") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error updating place", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ----------------------------------------------------
// 4. DELETE /api/v1/participations/{id} (Снятие с соревнований)
// ----------------------------------------------------
func (h *ParticipationHandler) DeleteParticipation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid participation ID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Participation not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error deleting participation", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}