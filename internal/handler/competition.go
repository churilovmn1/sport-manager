package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"sport-manager/internal/repository"
	"sport-manager/internal/service"

	"github.com/gorilla/mux"
)

// CompetitionHandler обрабатывает HTTP-запросы для соревнований
type CompetitionHandler struct {
	service *service.CompetitionService
}

func NewCompetitionHandler(s *service.CompetitionService) *CompetitionHandler {
	return &CompetitionHandler{service: s}
}

// POST /api/v1/competitions (Create)
func (h *CompetitionHandler) CreateCompetition(w http.ResponseWriter, r *http.Request) {
	var competition repository.Competition
	if err := json.NewDecoder(r.Body).Decode(&competition); err != nil {
		log.Printf("ERROR: Competition JSON Decode failed: %v", err)
		http.Error(w, "Invalid request body or date format. Expected RFC3339 (e.g., 2025-12-14T00:00:00Z)", http.StatusBadRequest)
		return
	}

	if err := h.service.Create(r.Context(), &competition); err != nil {
		log.Printf("ERROR: Competition creation failed. Input: %+v. Server Error: %v", competition, err)

		if strings.Contains(err.Error(), "validation error") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error creating competition", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusCreated, competition)
}

// GET /api/v1/competitions (List)
func (h *CompetitionHandler) ListCompetitions(w http.ResponseWriter, r *http.Request) {
	competitions, err := h.service.ListAll(r.Context())
	if err != nil {
		log.Printf("ERROR: Failed to list competitions: %v", err)
		http.Error(w, "Internal server error fetching list", http.StatusInternalServerError)
		return
	}
	// Оборачиваем список в объект для лучшей совместимости с API
	writeJSONResponse(w, http.StatusOK, map[string]interface{}{"competitions": competitions})
}

// GET /api/v1/competitions/{id} (GetByID)
func (h *CompetitionHandler) GetCompetition(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid competition ID", http.StatusBadRequest)
		return
	}

	competition, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if err.Error() == "competition not found" {
			http.Error(w, "Competition not found", http.StatusNotFound)
			return
		}
		log.Printf("ERROR: Failed to get competition: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	writeJSONResponse(w, http.StatusOK, competition)
}

// PUT /api/v1/competitions/{id} (Update)
func (h *CompetitionHandler) UpdateCompetition(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid competition ID", http.StatusBadRequest)
		return
	}

	var competition repository.Competition
	if err := json.NewDecoder(r.Body).Decode(&competition); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	competition.ID = id // Устанавливаем ID из URL

	if err := h.service.Update(r.Context(), &competition); err != nil {
		if err.Error() == "competition not found" {
			http.Error(w, "Competition not found", http.StatusNotFound)
			return
		}
		log.Printf("ERROR: Competition update failed: %v", err)
		http.Error(w, "Internal server error updating competition", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusOK, competition)
}

// DELETE /api/v1/competitions/{id} (Delete)
func (h *CompetitionHandler) DeleteCompetition(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid competition ID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if err.Error() == "competition not found" {
			http.Error(w, "Competition not found", http.StatusNotFound)
			return
		}
		log.Printf("ERROR: Competition deletion failed: %v", err)
		http.Error(w, "Internal server error deleting competition", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusNoContent, nil)
}
