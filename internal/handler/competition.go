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

// CompetitionHandler реализует интерфейс обработки запросов для соревнований
type CompetitionHandler struct {
	service *service.CompetitionService
}

// NewCompetitionHandler создает новый экземпляр хендлера соревнований
func NewCompetitionHandler(s *service.CompetitionService) *CompetitionHandler {
	return &CompetitionHandler{service: s}
}

// CreateCompetition создает новое соревнование на основе JSON-данных
func (h *CompetitionHandler) CreateCompetition(w http.ResponseWriter, r *http.Request) {
	var competition repository.Competition

	// Декодируем входящий JSON
	if err := json.NewDecoder(r.Body).Decode(&competition); err != nil {
		log.Printf("ERROR: Competition JSON Decode failed: %v", err)
		writeErrorResponse(w, http.StatusBadRequest, "Неверный формат тела запроса или даты. Ожидается RFC3339 (например, 2025-12-14T00:00:00Z)")
		return
	}

	// Вызываем сервис для сохранения в БД
	if err := h.service.Create(r.Context(), &competition); err != nil {
		log.Printf("ERROR: Competition creation failed: %v", err)

		if strings.Contains(err.Error(), "validation error") {
			writeErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		writeErrorResponse(w, http.StatusInternalServerError, "Ошибка при создании соревнования на сервере")
		return
	}

	writeJSONResponse(w, http.StatusCreated, competition)
}

// ListCompetitions возвращает список всех доступных соревнований
func (h *CompetitionHandler) ListCompetitions(w http.ResponseWriter, r *http.Request) {
	competitions, err := h.service.ListAll(r.Context())
	if err != nil {
		log.Printf("ERROR: Failed to list competitions: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Не удалось получить список соревнований")
		return
	}

	// Формируем ответ в формате, который ожидает фронтенд
	writeJSONResponse(w, http.StatusOK, map[string]interface{}{
		"competitions": competitions,
	})
}

// GetCompetition возвращает детальную информацию об одном соревновании по ID
func (h *CompetitionHandler) GetCompetition(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Некорректный ID соревнования")
		return
	}

	competition, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeErrorResponse(w, http.StatusNotFound, "Соревнование не найдено")
			return
		}
		log.Printf("ERROR: Failed to get competition: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Ошибка при поиске соревнования")
		return
	}

	writeJSONResponse(w, http.StatusOK, competition)
}

// UpdateCompetition обновляет данные существующего соревнования
func (h *CompetitionHandler) UpdateCompetition(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Некорректный ID соревнования")
		return
	}

	var competition repository.Competition
	if err := json.NewDecoder(r.Body).Decode(&competition); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Неверный формат данных")
		return
	}

	competition.ID = id // Принудительно ставим ID из URL

	if err := h.service.Update(r.Context(), &competition); err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeErrorResponse(w, http.StatusNotFound, "Соревнование не найдено")
			return
		}
		log.Printf("ERROR: Competition update failed: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Ошибка при обновлении данных")
		return
	}

	writeJSONResponse(w, http.StatusOK, competition)
}

// DeleteCompetition удаляет соревнование из системы
func (h *CompetitionHandler) DeleteCompetition(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Некорректный ID соревнования")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			writeErrorResponse(w, http.StatusNotFound, "Соревнование не найдено")
			return
		}
		log.Printf("ERROR: Competition deletion failed: %v", err)
		writeErrorResponse(w, http.StatusInternalServerError, "Не удалось удалить соревнование")
		return
	}

	// Возвращаем статус 204 (успешно, без тела ответа)
	w.WriteHeader(http.StatusNoContent)
}
