// internal/handler/athlete.go
package handler

import (
	"encoding/json"
	"net/http"
	"sport-manager/internal/service"
)

type AthleteHandler struct {
	service service.AthleteService
}

func NewAthleteHandler(s service.AthleteService) *AthleteHandler {
	return &AthleteHandler{service: s}
}

func (h *AthleteHandler) CreateAthlete(w http.ResponseWriter, r *http.Request) {
	var athlete repository.Athlete // Используйте модель DTO или напрямую репозиторий
	if err := json.NewDecoder(r.Body).Decode(&athlete); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.RegisterAthlete(r.Context(), &athlete); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(athlete)
}
