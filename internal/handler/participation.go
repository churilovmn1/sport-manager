package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"sport-manager/internal/repository"
	"sport-manager/internal/service"

	"github.com/gorilla/mux"
)

// ParticipationHandler обрабатывает HTTP-запросы для участия и результатов
type ParticipationHandler struct {
	service *service.ParticipationService
}

func NewParticipationHandler(s *service.ParticipationService) *ParticipationHandler {
	return &ParticipationHandler{service: s}
}

// POST /api/v1/participations (CreateParticipation)
func (h *ParticipationHandler) CreateParticipation(w http.ResponseWriter, r *http.Request) {
	// --- RBAC: Только для Администратора ---
	// NOTE: Раскомментировать, когда функция requireAdmin будет доступна
	// if !requireAdmin(w, r) {
	// 	return
	// }
	// ---------------------------------------

	var participation repository.Participation
	if err := json.NewDecoder(r.Body).Decode(&participation); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Place будет инициализирован как 0 при декодировании, что соответствует
	// тому, что результат пока неизвестен.

	if err := h.service.Create(r.Context(), &participation); err != nil {
		// Используем writeErrorResponse для гарантированного возврата JSON
		writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create participation: %v", err))
		return
	}

	writeJSONResponse(w, http.StatusCreated, participation)
}

// GET /api/v1/participations (ListParticipations) - Чтение разрешено всем
func (h *ParticipationHandler) ListParticipations(w http.ResponseWriter, r *http.Request) {
	participations, err := h.service.ListAll(r.Context())
	if err != nil {
		// Используем writeErrorResponse для гарантированного возврата JSON
		writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to list participations: %v", err))
		return
	}

	writeJSONResponse(w, http.StatusOK, participations)
}

// PUT /api/v1/participations/{id}/place (UpdatePlace)
func (h *ParticipationHandler) UpdatePlace(w http.ResponseWriter, r *http.Request) {
	// --- RBAC: Только для Администратора ---
	// NOTE: Раскомментировать, когда функция requireAdmin будет доступна
	// if !requireAdmin(w, r) {
	// 	return
	// }
	// ---------------------------------------

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid participation ID")
		return
	}

	var requestBody struct {
		Place int `json:"place"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request payload or missing 'place'")
		return
	}

	if err := h.service.UpdatePlace(r.Context(), id, requestBody.Place); err != nil {
		// Используем writeErrorResponse для гарантированного возврата JSON
		writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to update place: %v", err))
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"message": "Place updated successfully"})
}

// DELETE /api/v1/participations/{id} (DeleteParticipation)
func (h *ParticipationHandler) DeleteParticipation(w http.ResponseWriter, r *http.Request) {
	// --- RBAC: Только для Администратора ---
	// NOTE: Раскомментировать, когда функция requireAdmin будет доступна
	// if !requireAdmin(w, r) {
	// 	return
	// }
	// ---------------------------------------

	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid participation ID")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		// Используем writeErrorResponse для гарантированного возврата JSON
		writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to delete participation: %v", err))
		return
	}

	// HTTP Status 204 No Content обычно не возвращает тело
	w.WriteHeader(http.StatusNoContent)
}
