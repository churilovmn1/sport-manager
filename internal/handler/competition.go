package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"sport-manager/internal/repository"
	"sport-manager/internal/service"

	"github.com/gorilla/mux"
)

// NOTE: Предполагается, что функции writeErrorResponse, writeJSONResponse
// и requireAdmin доступны в этом пакете handler.

// CompetitionHandler обрабатывает HTTP-запросы для соревнований
type CompetitionHandler struct {
	service *service.CompetitionService
}

func NewCompetitionHandler(s *service.CompetitionService) *CompetitionHandler {
	return &CompetitionHandler{service: s}
}

// POST /api/v1/competitions (Create)
func (h *CompetitionHandler) CreateCompetition(w http.ResponseWriter, r *http.Request) {
	// --- RBAC: Только для Администратора ---
	// if !requireAdmin(w, r) {
	// 	return
	// }
	// ---------------------------------------

	var competition repository.Competition
	if err := json.NewDecoder(r.Body).Decode(&competition); err != nil {
		log.Printf("ERROR: Competition JSON Decode failed: %v", err)
		// ИСПРАВЛЕНО: writeErrorResponse вместо http.Error
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body or date format. Expected RFC3339 (e.g., 2025-12-14T00:00:00Z)")
		return
	}

	if err := h.service.Create(r.Context(), &competition); err != nil {
		log.Printf("ERROR: Competition creation failed. Input: %+v. Server Error: %v", competition, err)

		if strings.Contains(err.Error(), "validation error") {
			// ИСПРАВЛЕНО: writeErrorResponse вместо http.Error
			writeErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		// ИСПРАВЛЕНО: writeErrorResponse вместо http.Error
		writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Internal server error creating competition: %v", err))
		return
	}

	writeJSONResponse(w, http.StatusCreated, competition)
}

// GET /api/v1/competitions (List) - Чтение разрешено всем
func (h *CompetitionHandler) ListCompetitions(w http.ResponseWriter, r *http.Request) {
	competitions, err := h.service.ListAll(r.Context())
	if err != nil {
		log.Printf("ERROR: Failed to list competitions: %v", err)
		// ИСПРАВЛЕНО: writeErrorResponse вместо http.Error
		writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Internal server error fetching list: %v", err))
		return
	}

	// Структура ответа, которую ожидает ваш JS: {"competitions": [...]}
	writeJSONResponse(w, http.StatusOK, map[string]interface{}{"competitions": competitions})
}

// GET /api/v1/competitions/{id} (GetByID) - Чтение разрешено всем
func (h *CompetitionHandler) GetCompetition(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		// ИСПРАВЛЕНО: writeErrorResponse вместо http.Error
		writeErrorResponse(w, http.StatusBadRequest, "Invalid competition ID")
		return
	}

	competition, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			// ИСПРАВЛЕНО: writeErrorResponse вместо http.Error
			writeErrorResponse(w, http.StatusNotFound, "Competition not found")
			return
		}
		log.Printf("ERROR: Failed to get competition: %v", err)
		// ИСПРАВЛЕНО: writeErrorResponse вместо http.Error
		writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeJSONResponse(w, http.StatusOK, competition)
}

// PUT /api/v1/competitions/{id} (Update)
func (h *CompetitionHandler) UpdateCompetition(w http.ResponseWriter, r *http.Request) {
	// --- RBAC: Только для Администратора ---
	// if !requireAdmin(w, r) {
	// 	return
	// }
	// ---------------------------------------

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		// ИСПРАВЛЕНО: writeErrorResponse вместо http.Error
		writeErrorResponse(w, http.StatusBadRequest, "Invalid competition ID")
		return
	}

	var competition repository.Competition
	if err := json.NewDecoder(r.Body).Decode(&competition); err != nil {
		// ИСПРАВЛЕНО: writeErrorResponse вместо http.Error
		writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	competition.ID = id // Устанавливаем ID из URL

	if err := h.service.Update(r.Context(), &competition); err != nil {
		if strings.Contains(err.Error(), "not found") {
			// ИСПРАВЛЕНО: writeErrorResponse вместо http.Error
			writeErrorResponse(w, http.StatusNotFound, "Competition not found")
			return
		}
		log.Printf("ERROR: Competition update failed: %v", err)
		// ИСПРАВЛЕНО: writeErrorResponse вместо http.Error
		writeErrorResponse(w, http.StatusInternalServerError, "Internal server error updating competition")
		return
	}

	writeJSONResponse(w, http.StatusOK, competition)
}

// DELETE /api/v1/competitions/{id} (Delete)
func (h *CompetitionHandler) DeleteCompetition(w http.ResponseWriter, r *http.Request) {
	// --- RBAC: Только для Администратора ---
	// if !requireAdmin(w, r) {
	// 	return
	// }
	// ---------------------------------------

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		// ИСПРАВЛЕНО: writeErrorResponse вместо http.Error
		writeErrorResponse(w, http.StatusBadRequest, "Invalid competition ID")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			// ИСПРАВЛЕНО: writeErrorResponse вместо http.Error
			writeErrorResponse(w, http.StatusNotFound, "Competition not found")
			return
		}
		log.Printf("ERROR: Competition deletion failed: %v", err)
		// ИСПРАВЛЕНО: writeErrorResponse вместо http.Error
		writeErrorResponse(w, http.StatusInternalServerError, "Internal server error deleting competition")
		return
	}

	// 204 No Content - не возвращает тело
	w.WriteHeader(http.StatusNoContent)
}
