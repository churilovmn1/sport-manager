package service

import (
	"context"
	"fmt"
	"time"

	"sport-manager/internal/repository"
)

// CompetitionService реализует бизнес-логику управления спортивными мероприятиями.
type CompetitionService struct {
	repo *repository.CompetitionRepository
}

// NewCompetitionService создает новый экземпляр сервиса соревнований.
func NewCompetitionService(repo *repository.CompetitionRepository) *CompetitionService {
	return &CompetitionService{repo: repo}
}

// --- МЕТОДЫ БИЗНЕС-ЛОГИКИ ---

// Create выполняет валидацию данных перед регистрацией нового соревнования.
func (s *CompetitionService) Create(ctx context.Context, c *repository.Competition) error {
	// Валидация обязательных полей
	if len(c.Name) == 0 {
		return fmt.Errorf("ошибка валидации: название соревнования обязательно")
	}
	if len(c.Location) == 0 {
		return fmt.Errorf("ошибка валидации: место проведения обязательно")
	}
	if c.StartDate.IsZero() {
		return fmt.Errorf("ошибка валидации: дата начала обязательна")
	}

	// Бизнес-правило: соревнование не может быть создано на прошедшую дату
	if c.StartDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return fmt.Errorf("ошибка валидации: дата начала не может быть в прошлом")
	}

	return s.repo.Create(ctx, c)
}

// GetByID возвращает детальную информацию о соревновании.
func (s *CompetitionService) GetByID(ctx context.Context, id int) (*repository.Competition, error) {
	if id <= 0 {
		return nil, fmt.Errorf("ошибка: некорректный ID соревнования")
	}
	return s.repo.GetByID(ctx, id)
}

// ListAll возвращает полный перечень запланированных мероприятий.
func (s *CompetitionService) ListAll(ctx context.Context) ([]repository.Competition, error) {
	return s.repo.ListAll(ctx)
}

// Update проверяет обновленные данные перед сохранением в базу.
func (s *CompetitionService) Update(ctx context.Context, c *repository.Competition) error {
	if c.ID <= 0 {
		return fmt.Errorf("ошибка: для обновления необходим ID записи")
	}
	if len(c.Name) == 0 {
		return fmt.Errorf("ошибка валидации: название не может быть пустым")
	}

	return s.repo.Update(ctx, c)
}

// Delete удаляет мероприятие из системы.
func (s *CompetitionService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("ошибка: некорректный ID для удаления")
	}
	return s.repo.Delete(ctx, id)
}
