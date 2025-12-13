package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sport-manager/internal/repository"
	"sport-manager/internal/service"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// NOTE: Предполагается, что функции writeErrorResponse и writeJSONResponse
// определены в этом или другом файле в пакете handler и доступны.

// AthleteHandler представляет HTTP-хендлеры для работы со спортсменами
type AthleteHandler struct {
	service *service.AthleteService
}

func NewAthleteHandler(service *service.AthleteService) *AthleteHandler {
	return &AthleteHandler{service: service}
}

// CreateAthlete обрабатывает POST-запрос на создание нового спортсмена
func (h *AthleteHandler) CreateAthlete(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FullName  string `json:"full_name"`
		BirthDate string `json:"birth_date"`
		Gender    string `json:"gender"`
		// SportID и RankID УДАЛЕНЫ
		IsActive bool   `json:"is_active"`
		Address  string `json:"address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON input")
		return
	}

	var birthDate time.Time
	var err error

	// Корректный парсинг даты
	birthDate, err = time.Parse(time.RFC3339Nano, input.BirthDate)

	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("Invalid birth date format. Error: %v", err))
		return
	}

	if input.Gender == "" {
		input.Gender = "M"
	}

	athlete := &repository.Athlete{
		FullName:  input.FullName,
		BirthDate: birthDate,
		Gender:    input.Gender,
		// SportID и RankID УДАЛЕНЫ
		IsActive: true,
		Address:  input.Address,
	}

	if err := h.service.CreateAthlete(r.Context(), athlete); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create athlete: %v", err))
		return
	}

	writeJSONResponse(w, http.StatusCreated, map[string]interface{}{"status": "success", "id": athlete.ID})
}

// ListAllAthletes обрабатывает GET-запрос на получение всех спортсменов
func (h *AthleteHandler) ListAllAthletes(w http.ResponseWriter, r *http.Request) {
	athletes, err := h.service.ListAll(r.Context())
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to list athletes: %v", err))
		return
	}

	writeJSONResponse(w, http.StatusOK, athletes)
}

// GetAthleteByID обрабатывает GET-запрос на получение спортсмена по ID
func (h *AthleteHandler) GetAthleteByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid athlete ID")
		return
	}

	athlete, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Athlete not found: %v", err))
		return
	}

	writeJSONResponse(w, http.StatusOK, athlete)
}

// UpdateAthlete обрабатывает PUT-запрос на обновление спортсмена
func (h *AthleteHandler) UpdateAthlete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid athlete ID")
		return
	}

	var input repository.Athlete
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON input")
		return
	}

	input.ID = id

	if err := h.service.Update(r.Context(), &input); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to update athlete: %v", err))
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"status": "success"})
}

// DeleteAthlete обрабатывает DELETE-запрос на удаление спортсмена
func (h *AthleteHandler) DeleteAthlete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid athlete ID")
		return
	}

	log.Printf("Handler: Attempting to delete athlete with ID: %d", id)

	if err := h.service.Delete(r.Context(), id); err != nil {
		log.Printf("Handler: Deletion failed for ID %d: %v", id, err)
		writeErrorResponse(w, http.StatusNotFound, fmt.Sprintf("Failed to delete athlete: %v", err))
		return
	}

	log.Printf("Handler: Successfully deleted athlete with ID: %d", id)

	writeJSONResponse(w, http.StatusOK, map[string]string{"status": "success"})
}
