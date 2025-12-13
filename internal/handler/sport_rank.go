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

// SimpleHandler обрабатывает HTTP-запросы для справочников (Sport, Rank)
type SimpleHandler struct {
	service *service.SimpleService
}

func NewSimpleHandler(s *service.SimpleService) *SimpleHandler {
	return &SimpleHandler{service: s}
}

// ----------------------------------------------------
// 1. POST /api/v1/sports или /ranks
// ----------------------------------------------------
func (h *SimpleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var model repository.SimpleModel
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.Create(r.Context(), &model); err != nil {
		if strings.Contains(err.Error(), "validation error") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal server error creating item", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusCreated, model)
}

// ----------------------------------------------------
// 2. GET /api/v1/sports или /ranks
// ----------------------------------------------------
func (h *SimpleHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	models, err := h.service.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Internal server error listing items", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusOK, models)
}

// ----------------------------------------------------
// 3. DELETE /api/v1/sports/{id} или /ranks/{id}
// ----------------------------------------------------
func (h *SimpleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}
		// Ошибка FOREIGN KEY:
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			http.Error(w, "Cannot delete item: it is currently linked to one or more athletes.", http.StatusConflict)
			return
		}
		http.Error(w, "Internal server error deleting item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}