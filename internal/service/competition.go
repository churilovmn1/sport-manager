package service

import (
	"context"
	"fmt"
	"time"

	"sport-manager/internal/repository"
)

// CompetitionService содержит методы бизнес-логики для соревнований
type CompetitionService struct {
	repo *repository.CompetitionRepository
}

func NewCompetitionService(repo *repository.CompetitionRepository) *CompetitionService {
	return &CompetitionService{repo: repo}
}

// Create проверяет данные и вызывает репозиторий для создания соревнования
func (s *CompetitionService) Create(ctx context.Context, c *repository.Competition) error {
	if len(c.Name) == 0 {
		return fmt.Errorf("validation error: competition name is required")
	}
	if len(c.Location) == 0 {
		return fmt.Errorf("validation error: competition location is required")
	}
	if c.StartDate.IsZero() {
		return fmt.Errorf("validation error: start date is required")
	}

	// Можно добавить проверку, что дата не в прошлом
	if c.StartDate.Before(time.Now().AddDate(0, 0, -1)) {
		// return fmt.Errorf("validation error: start date must be in the future")
	}

	return s.repo.Create(ctx, c)
}

// GetByID вызывает репозиторий для получения одного соревнования
func (s *CompetitionService) GetByID(ctx context.Context, id int) (*repository.Competition, error) {
	return s.repo.GetByID(ctx, id)
}

// ListAll вызывает репозиторий для получения списка соревнований
func (s *CompetitionService) ListAll(ctx context.Context) ([]repository.Competition, error) {
	return s.repo.ListAll(ctx)
}

// Update проверяет данные и вызывает репозиторий для обновления
func (s *CompetitionService) Update(ctx context.Context, c *repository.Competition) error {
	if len(c.Name) == 0 {
		return fmt.Errorf("validation error: competition name is required")
	}
	return s.repo.Update(ctx, c)
}

// Delete вызывает репозиторий для удаления
func (s *CompetitionService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
