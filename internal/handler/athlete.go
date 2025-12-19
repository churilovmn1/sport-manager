package handler

import (
	"encoding/json"
	"net/http"
	"sport-manager/internal/repository"
	"sport-manager/internal/service"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// AthleteHandler — слой контроллеров для сущности "Спортсмен".
// Он не знает о базе данных напрямую, а работает только через AthleteService.
type AthleteHandler struct {
	service *service.AthleteService
}

// NewAthleteHandler создает новый экземпляр хендлера (Dependency Injection)
func NewAthleteHandler(service *service.AthleteService) *AthleteHandler {
	return &AthleteHandler{service: service}
}

// CreateAthlete обрабатывает POST /api/v1/athletes
// Создает нового спортсмена на основе данных из тела запроса (JSON).
func (h *AthleteHandler) CreateAthlete(w http.ResponseWriter, r *http.Request) {
	// Анонимная структура для парсинга входных данных
	var input struct {
		FullName  string `json:"full_name"`
		BirthDate string `json:"birth_date"`
		Gender    string `json:"gender"`
		Address   string `json:"address"`
	}

	// Декодируем JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Неверный формат JSON")
		return
	}

	// Преобразуем строковую дату в формат time.Time
	birthDate, err := time.Parse(time.RFC3339Nano, input.BirthDate)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Неверный формат даты рождения")
		return
	}

	// Создаем объект модели для передачи в сервис
	athlete := &repository.Athlete{
		FullName:  input.FullName,
		BirthDate: birthDate,
		Gender:    input.Gender,
		IsActive:  true, // По умолчанию спортсмен активен
		Address:   input.Address,
	}

	// Вызываем бизнес-логику создания
	if err := h.service.CreateAthlete(r.Context(), athlete); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Ошибка при сохранении спортсмена")
		return
	}

	// Возвращаем успешный ответ с ID нового спортсмена
	writeJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"status": "success",
		"id":     athlete.ID,
	})
}

// ListAllAthletes обрабатывает GET /api/v1/athletes
// Возвращает полный список спортсменов.
func (h *AthleteHandler) ListAllAthletes(w http.ResponseWriter, r *http.Request) {
	athletes, err := h.service.ListAll(r.Context())
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Ошибка получения списка")
		return
	}

	writeJSONResponse(w, http.StatusOK, athletes)
}

// GetAthleteByID обрабатывает GET /api/v1/athletes/{id}
func (h *AthleteHandler) GetAthleteByID(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID из параметров пути (URL)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Некорректный ID")
		return
	}

	athlete, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeErrorResponse(w, http.StatusNotFound, "Спортсмен не найден")
		return
	}

	writeJSONResponse(w, http.StatusOK, athlete)
}

// UpdateAthlete обрабатывает PUT /api/v1/athletes/{id}
func (h *AthleteHandler) UpdateAthlete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Некорректный ID")
		return
	}

	var input repository.Athlete
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Ошибка парсинга данных")
		return
	}

	input.ID = id

	if err := h.service.Update(r.Context(), &input); err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Ошибка при обновлении")
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"status": "success"})
}

// DeleteAthlete обрабатывает DELETE /api/v1/athletes/{id}
func (h *AthleteHandler) DeleteAthlete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Некорректный ID")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		writeErrorResponse(w, http.StatusNotFound, "Не удалось удалить спортсмена")
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"status": "success"})
}
