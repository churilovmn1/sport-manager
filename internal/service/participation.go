package service

import (
	"context"
	"fmt"
	"strings"

	"sport-manager/internal/repository"
)

// ParticipationService содержит методы бизнес-логики для участия и результатов
type ParticipationService struct {
	repo         *repository.ParticipationRepository
	athleteRepo  *repository.AthleteRepository // Для проверки AthleteID
	compRepo     *repository.CompetitionRepository // Для проверки CompetitionID
}

func NewParticipationService(
	repo *repository.ParticipationRepository,
	athleteRepo *repository.AthleteRepository,
	compRepo *repository.CompetitionRepository,
) *ParticipationService {
	return &ParticipationService{
		repo: repo,
		athleteRepo: athleteRepo,
		compRepo: compRepo,
	}
}

// ----------------------------------------------------
// 1. CREATE (Регистрация)
// ----------------------------------------------------
func (s *ParticipationService) Create(ctx context.Context, p *repository.Participation) error {
	// Бизнес-логика 1: Проверка существования спортсмена
	if _, err := s.athleteRepo.GetByID(ctx, p.AthleteID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return fmt.Errorf("validation error: athlete with ID %d not found", p.AthleteID)
		}
		return fmt.Errorf("failed to check athlete existence: %w", err)
	}

	// Бизнес-логика 2: Проверка существования соревнования
	if _, err := s.compRepo.GetByID(ctx, p.CompetitionID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return fmt.Errorf("validation error: competition with ID %d not found", p.CompetitionID)
		}
		return fmt.Errorf("failed to check competition existence: %w", err)
	}

	// Бизнес-логика 3: Проверка на дублирование (опционально, но важно)
	// В реальном проекте здесь нужно убедиться, что спортсмен еще не зарегистрирован
	// (например, написав отдельный метод в репозитории: GetByAthleteAndCompetitionID)

	return s.repo.Create(ctx, p)
}

// ----------------------------------------------------
// 2. READ
// ----------------------------------------------------
func (s *ParticipationService) GetByID(ctx context.Context, id int) (*repository.ParticipationDetails, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ParticipationService) GetAll(ctx context.Context) ([]*repository.ParticipationDetails, error) {
	return s.repo.GetAll(ctx)
}

// ----------------------------------------------------
// 3. UPDATE (Результаты)
// ----------------------------------------------------
func (s *ParticipationService) UpdatePlace(ctx context.Context, id int, place int) error {
	// Бизнес-логика: Проверка валидности места
	if place <= 0 {
		return fmt.Errorf("validation error: place must be a positive number")
	}

	return s.repo.UpdatePlace(ctx, id, place)
}

// ----------------------------------------------------
// 4. DELETE
// ----------------------------------------------------
func (s *ParticipationService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}