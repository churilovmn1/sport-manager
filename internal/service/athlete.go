package service

import (
	"context"
	"fmt"
	"time"

	"sport-manager/internal/repository"
)

// AthleteService содержит методы бизнес-логики для спортсменов
type AthleteService struct {
	repo *repository.AthleteRepository
}

func NewAthleteService(repo *repository.AthleteRepository) *AthleteService {
	return &AthleteService{repo: repo}
}

// Create проверяет данные и вызывает репозиторий для создания спортсмена
func (s *AthleteService) Create(ctx context.Context, a *repository.Athlete) error {
	// ------------------------------------------------
	// Бизнес-логика 1: Проверка валидности данных
	// ------------------------------------------------
	if len(a.FullName) == 0 {
		return fmt.Errorf("validation error: full name is required")
	}
	if a.BirthDate.IsZero() || a.BirthDate.After(time.Now()) {
		return fmt.Errorf("validation error: birth date is invalid")
	}

	// В Go-проекте такого уровня часто можно добавить дополнительную проверку:
	// 1. Проверить, что SportID и RankID существуют в БД (через отдельные репозитории Sports/Ranks)

	return s.repo.Create(ctx, a)
}

// GetByID получает спортсмена
func (s *AthleteService) GetByID(ctx context.Context, id int) (*repository.Athlete, error) {
	athlete, err := s.repo.GetByID(ctx, id)
	if err != nil {
		// Здесь можно обработать ошибку, если она специфична для сервиса
		return nil, err
	}
	return athlete, nil
}

// GetAll получает список всех спортсменов
func (s *AthleteService) GetAll(ctx context.Context) ([]*repository.Athlete, error) {
	return s.repo.GetAll(ctx)
}

// Update обновляет данные спортсмена
func (s *AthleteService) Update(ctx context.Context, a *repository.Athlete) error {
	// ------------------------------------------------
	// Бизнес-логика 2: Проверка перед обновлением
	// ------------------------------------------------
	if a.ID <= 0 {
		return fmt.Errorf("validation error: athlete ID is required for update")
	}
	// Можно добавить повторную проверку валидности полей, как в Create

	return s.repo.Update(ctx, a)
}

// Delete удаляет спортсмена
func (s *AthleteService) Delete(ctx context.Context, id int) error {
	// ------------------------------------------------
	// Бизнес-логика 3: Проверка перед удалением
	// ------------------------------------------------
	// В реальном проекте здесь можно проверить,
	// участвует ли спортсмен в будущих соревнованиях,
	// и запретить удаление, если да.

	return s.repo.Delete(ctx, id)
}
