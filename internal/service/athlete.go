// internal/service/athlete.go
package service

import (
    "context"
    "sport-manager/internal/repository"
)

type AthleteService struct {
    repo repository.AthleteRepository
}

func NewAthleteService(repo repository.AthleteRepository) *AthleteService {
    return &AthleteService{repo: repo}
}

func (s *AthleteService) RegisterAthlete(ctx context.Context, athlete *repository.Athlete) error {
    // Здесь бизнес-логика, например, валидация данных
    if len(athlete.Name) < 3 {
        return errors.New("invalid athlete name")
    }
    return s.repo.Create(ctx, athlete)
}