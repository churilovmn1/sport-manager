package service

import (
	"context"
	"errors"
	"fmt"
	"sport-manager/internal/repository"
)

// AthleteService реализует бизнес-логику для работы со спортсменами.
type AthleteService struct {
	repo *repository.AthleteRepository
}

// NewAthleteService создает новый экземпляр сервиса.
func NewAthleteService(repo *repository.AthleteRepository) *AthleteService {
	return &AthleteService{repo: repo}
}

// --- БИЗНЕС-ЛОГИКА ---

// CreateAthlete проверяет данные и создает нового спортсмена.
func (s *AthleteService) CreateAthlete(ctx context.Context, athlete *repository.Athlete) error {
	// Валидация: имя не может быть пустым
	if athlete.FullName == "" {
		return errors.New("валидация: ФИО спортсмена обязательно для заполнения")
	}

	// Вызов слоя данных
	if err := s.repo.Create(ctx, athlete); err != nil {
		return fmt.Errorf("service: не удалось создать атлета: %w", err)
	}
	return nil
}

// ListAll возвращает список всех спортсменов.
func (s *AthleteService) ListAll(ctx context.Context) ([]repository.Athlete, error) {
	athletes, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("service: ошибка получения списка: %w", err)
	}
	return athletes, nil
}

// GetByID возвращает данные конкретного спортсмена по ID.
func (s *AthleteService) GetByID(ctx context.Context, id int) (*repository.Athlete, error) {
	if id <= 0 {
		return nil, errors.New("валидация: некорректный ID спортсмена")
	}

	athlete, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service: ошибка при поиске атлета: %w", err)
	}
	return athlete, nil
}

// Update обновляет информацию о спортсмене, предварительно проверяя данные.
func (s *AthleteService) Update(ctx context.Context, athlete *repository.Athlete) error {
	if athlete.ID <= 0 {
		return errors.New("валидация: ID спортсмена обязателен для обновления")
	}
	if athlete.FullName == "" {
		return errors.New("валидация: ФИО не может быть пустым")
	}

	if err := s.repo.Update(ctx, athlete); err != nil {
		return fmt.Errorf("service: не удалось обновить данные: %w", err)
	}
	return nil
}

// Delete удаляет спортсмена из системы.
func (s *AthleteService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("валидация: некорректный ID для удаления")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("service: ошибка при удалении: %w", err)
	}
	return nil
}
