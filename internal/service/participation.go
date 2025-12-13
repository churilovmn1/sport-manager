package service

import (
	"context"
	"fmt"

	"sport-manager/internal/repository"
)

// ParticipationService содержит бизнес-логику для участия
type ParticipationService struct {
	repo            *repository.ParticipationRepository
	athleteRepo     *repository.AthleteRepository
	competitionRepo *repository.CompetitionRepository
}

func NewParticipationService(
	repo *repository.ParticipationRepository,
	athleteRepo *repository.AthleteRepository,
	competitionRepo *repository.CompetitionRepository,
) *ParticipationService {
	return &ParticipationService{
		repo:            repo,
		athleteRepo:     athleteRepo,
		competitionRepo: competitionRepo,
	}
}

// Create создает новую запись об участии
func (s *ParticipationService) Create(ctx context.Context, p *repository.Participation) error {
	if p.AthleteID == 0 || p.CompetitionID == 0 {
		return fmt.Errorf("validation error: athlete ID and competition ID are required")
	}

	// Дополнительная проверка существования спортсмена и соревнования
	if _, err := s.athleteRepo.GetByID(ctx, p.AthleteID); err != nil {
		return fmt.Errorf("service: athlete not found: %w", err)
	}
	if _, err := s.competitionRepo.GetByID(ctx, p.CompetitionID); err != nil {
		return fmt.Errorf("service: competition not found: %w", err)
	}

	return s.repo.Create(ctx, p)
}

// ListAll получает список всех участий.
// NOTE: После обновления репозитория, этот метод теперь возвращает
// структуру Participation с заполненными полями AthleteName и CompetitionName.
// Дополнительная логика на уровне Service здесь не требуется.
func (s *ParticipationService) ListAll(ctx context.Context) ([]repository.Participation, error) {
	participations, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("service: failed to list participations: %w", err)
	}
	return participations, nil
}

// UpdatePlace обновляет место в соревновании
func (s *ParticipationService) UpdatePlace(ctx context.Context, id int, place int) error {
	if id == 0 || place <= 0 {
		return fmt.Errorf("validation error: valid ID and place number are required")
	}
	return s.repo.UpdatePlace(ctx, id, place)
}

// Delete удаляет запись об участии
func (s *ParticipationService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
