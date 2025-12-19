package service

import (
	"context"
	"fmt"

	"sport-manager/internal/repository"
)

// ParticipationService управляет логикой регистрации атлетов на соревнования.
// Он координирует работу нескольких репозиториев для обеспечения целостности связей.
type ParticipationService struct {
	repo            *repository.ParticipationRepository
	athleteRepo     *repository.AthleteRepository
	competitionRepo *repository.CompetitionRepository
}

// NewParticipationService инициализирует сервис со всеми необходимыми зависимостями.
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

// --- БИЗНЕС-ЛОГИКА ---

// Create регистрирует спортсмена на соревнование, предварительно проверяя их существование.
func (s *ParticipationService) Create(ctx context.Context, p *repository.Participation) error {
	// 1. Первичная валидация входных данных
	if p.AthleteID == 0 || p.CompetitionID == 0 {
		return fmt.Errorf("ошибка валидации: ID атлета и ID соревнования обязательны")
	}

	// 2. Проверка существования атлета
	// Это предотвращает создание "битых" связей в базе данных
	if _, err := s.athleteRepo.GetByID(ctx, p.AthleteID); err != nil {
		return fmt.Errorf("service: указанный атлет не найден: %w", err)
	}

	// 3. Проверка существования соревнования
	if _, err := s.competitionRepo.GetByID(ctx, p.CompetitionID); err != nil {
		return fmt.Errorf("service: указанное соревнование не найдено: %w", err)
	}

	// 4. Сохранение записи в БД
	return s.repo.Create(ctx, p)
}

// ListAll возвращает расширенный список участий (с именами атлетов и названиями турниров).
func (s *ParticipationService) ListAll(ctx context.Context) ([]repository.Participation, error) {
	participations, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("service: не удалось получить список регистраций: %w", err)
	}
	return participations, nil
}

// UpdatePlace фиксирует результат (место), занятое атлетом.
func (s *ParticipationService) UpdatePlace(ctx context.Context, id int, place int) error {
	if id <= 0 {
		return fmt.Errorf("ошибка: некорректный ID записи")
	}
	if place <= 0 {
		return fmt.Errorf("ошибка валидации: занятое место должно быть положительным числом")
	}

	return s.repo.UpdatePlace(ctx, id, place)
}

// Delete аннулирует участие атлета в соревновании.
func (s *ParticipationService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("ошибка: некорректный ID для удаления")
	}
	return s.repo.Delete(ctx, id)
}
