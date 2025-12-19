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

// ParticipationHandler управляет процессом регистрации спортсменов на соревнования
type ParticipationHandler struct {
	service *service.ParticipationService
}

// NewParticipationHandler создает новый экземпляр хендлера для записей участия
func NewParticipationHandler(s *service.ParticipationService) *ParticipationHandler {
	return &ParticipationHandler{service: s}
}

// CreateParticipation регистрирует спортсмена на конкретное соревнование
func (h *ParticipationHandler) CreateParticipation(w http.ResponseWriter, r *http.Request) {
	var participation repository.Participation

	// Декодируем ID атлета и соревнования из JSON
	if err := json.NewDecoder(r.Body).Decode(&participation); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Некорректный формат данных в запросе")
		return
	}

	// Поле Place по умолчанию будет 0 (результат еще не определен)
	if err := h.service.Create(r.Context(), &participation); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, fmt.Sprintf("Ошибка при создании регистрации: %v", err))
		return
	}

	writeJSONResponse(w, http.StatusCreated, participation)
}

// ListParticipations возвращает список всех регистраций с именами атлетов и турниров
func (h *ParticipationHandler) ListParticipations(w http.ResponseWriter, r *http.Request) {
	participations, err := h.service.ListAll(r.Context())
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Не удалось получить список регистраций")
		return
	}

	writeJSONResponse(w, http.StatusOK, participations)
}

// UpdatePlace обновляет занятое спортсменом место в рамках соревнования
func (h *ParticipationHandler) UpdatePlace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Некорректный ID записи участия")
		return
	}

	// Структура для приема только поля 'place'
	var requestBody struct {
		Place int `json:"place"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Необходимо указать корректное число в поле 'place'")
		return
	}

	// Обновляем только результат (место)
	if err := h.service.UpdatePlace(r.Context(), id, requestBody.Place); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Не удалось обновить результат")
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{
		"message": "Результат успешно обновлен",
	})
}

// DeleteParticipation удаляет запись о регистрации (например, при отказе атлета от участия)
func (h *ParticipationHandler) DeleteParticipation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Некорректный ID записи")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Ошибка при удалении записи из системы")
		return
	}

	// Возвращаем 204 No Content (успех без тела ответа)
	w.WriteHeader(http.StatusNoContent)
}
