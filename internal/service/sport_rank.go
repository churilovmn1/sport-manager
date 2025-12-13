package service

import (
	"context"
	"fmt"
	"sport-manager/internal/repository"
)

// SimpleService содержит логику для справочников
type SimpleService struct {
	repo *repository.SimpleRepository
}

func NewSimpleService(repo *repository.SimpleRepository) *SimpleService {
	return &SimpleService{repo: repo}
}

func (s *SimpleService) Create(ctx context.Context, model *repository.SimpleModel) error {
	if len(model.Name) < 2 {
		return fmt.Errorf("validation error: name must be at least 2 characters")
	}
	return s.repo.Create(ctx, model)
}

func (s *SimpleService) GetAll(ctx context.Context) ([]*repository.SimpleModel, error) {
	return s.repo.GetAll(ctx)
}

func (s *SimpleService) Delete(ctx context.Context, id int) error {
	// В реальном проекте, перед удалением нужно проверить,
	// не используется ли этот Sport/Rank каким-либо Athlete.
	// Если используется, то удаление не пройдет из-за FOREIGN KEY.
	return s.repo.Delete(ctx, id)
}