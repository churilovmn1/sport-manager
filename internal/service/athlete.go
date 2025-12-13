package service

import (
	"context"
	"fmt"
	"sport-manager/internal/repository"
)

type AthleteService struct {
	repo *repository.AthleteRepository
}

func NewAthleteService(repo *repository.AthleteRepository) *AthleteService {
	return &AthleteService{repo: repo}
}

// CreateAthlete
func (s *AthleteService) CreateAthlete(ctx context.Context, athlete *repository.Athlete) error {
	if err := s.repo.Create(ctx, athlete); err != nil {
		return fmt.Errorf("service: failed to create athlete: %w", err)
	}
	return nil
}

// ListAll
func (s *AthleteService) ListAll(ctx context.Context) ([]repository.Athlete, error) {
	athletes, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("service: failed to list all athletes: %w", err)
	}
	return athletes, nil
}

// GetByID
func (s *AthleteService) GetByID(ctx context.Context, id int) (*repository.Athlete, error) {
	athlete, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("service: failed to get athlete by ID: %w", err)
	}
	return athlete, nil
}

// Update
func (s *AthleteService) Update(ctx context.Context, athlete *repository.Athlete) error {
	if err := s.repo.Update(ctx, athlete); err != nil {
		return fmt.Errorf("service: failed to update athlete: %w", err)
	}
	return nil
}

// Delete
func (s *AthleteService) Delete(ctx context.Context, id int) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("service: failed to delete athlete: %w", err)
	}
	return nil
}
