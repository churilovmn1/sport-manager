package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Participation представляет запись об участии спортсмена в соревновании
type Participation struct {
	ID            int          `json:"id"`
	AthleteID     int          `json:"athlete_id"`
	CompetitionID int          `json:"competition_id"`
	Place         sql.NullInt32 `json:"place"` // Место (результат). sql.NullInt32, так как может быть NULL до завершения.
	RegisteredAt  time.Time    `json:"registered_at"`
}

// ParticipationDetails для вывода полной информации (для List/Get)
type ParticipationDetails struct {
	ID              int       `json:"id"`
	AthleteID       int       `json:"athlete_id"`
	AthleteFullName string    `json:"athlete_full_name"`
	CompetitionID   int       `json:"competition_id"`
	CompetitionName string    `json:"competition_name"`
	Place           int32     `json:"place,omitempty"`
	RegisteredAt    time.Time `json:"registered_at"`
}

// ParticipationRepository отвечает за работу с таблицей participations
type ParticipationRepository struct {
	db *sql.DB
}

func NewParticipationRepository(db *sql.DB) *ParticipationRepository {
	return &ParticipationRepository{db: db}
}

// ----------------------------------------------------
// 1. CREATE (Регистрация спортсмена на соревнование)
// ----------------------------------------------------
func (r *ParticipationRepository) Create(ctx context.Context, p *Participation) error {
	query := `
		INSERT INTO participations (athlete_id, competition_id, registered_at)
		VALUES ($1, $2, $3)
		RETURNING id`
	
	// Вставляем AthleteID, CompetitionID и текущее время регистрации
	p.RegisteredAt = time.Now()
	err := r.db.QueryRowContext(ctx, query,
		p.AthleteID, p.CompetitionID, p.RegisteredAt,
	).Scan(&p.ID)
	
	if err != nil {
		return fmt.Errorf("repository: failed to create participation: %w", err)
	}
	return nil
}

// ----------------------------------------------------
// 2. GET (Получение одной записи об участии)
// ----------------------------------------------------
func (r *ParticipationRepository) GetByID(ctx context.Context, id int) (*ParticipationDetails, error) {
	// Используем JOIN для получения имени спортсмена и названия соревнования
	query := `
		SELECT 
			p.id, p.athlete_id, a.full_name, p.competition_id, c.title, p.place, p.registered_at
		FROM participations p
		JOIN athletes a ON p.athlete_id = a.id
		JOIN competitions c ON p.competition_id = c.id
		WHERE p.id = $1`
	
	detail := &ParticipationDetails{}
	var place sql.NullInt32

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&detail.ID, &detail.AthleteID, &detail.AthleteFullName,
		&detail.CompetitionID, &detail.CompetitionName, &place, &detail.RegisteredAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("participation not found")
		}
		return nil, fmt.Errorf("repository: failed to get participation by ID: %w", err)
	}

	if place.Valid {
		detail.Place = place.Int32
	}
	
	return detail, nil
}

// ----------------------------------------------------
// 3. GET ALL (Получение списка всех участий)
// ----------------------------------------------------
func (r *ParticipationRepository) GetAll(ctx context.Context) ([]*ParticipationDetails, error) {
	// Используем JOIN для получения полной информации
	query := `
		SELECT 
			p.id, p.athlete_id, a.full_name, p.competition_id, c.title, p.place, p.registered_at
		FROM participations p
		JOIN athletes a ON p.athlete_id = a.id
		JOIN competitions c ON p.competition_id = c.id
		ORDER BY c.start_date DESC, a.full_name`
	
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to get all participations: %w", err)
	}
	defer rows.Close()

	var details []*ParticipationDetails
	for rows.Next() {
		detail := &ParticipationDetails{}
		var place sql.NullInt32
		
		err := rows.Scan(
			&detail.ID, &detail.AthleteID, &detail.AthleteFullName,
			&detail.CompetitionID, &detail.CompetitionName, &place, &detail.RegisteredAt,
		)
		
		if err != nil {
			return nil, fmt.Errorf("repository: failed to scan participation row: %w", err)
		}

		if place.Valid {
			detail.Place = place.Int32
		}
		details = append(details, detail)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: rows error: %w", err)
	}
	return details, nil
}

// ----------------------------------------------------
// 4. UPDATE (Проставление результата/места)
// ----------------------------------------------------
func (r *ParticipationRepository) UpdatePlace(ctx context.Context, id int, place int) error {
	query := `
		UPDATE participations 
		SET place = $1
		WHERE id = $2`
	
	result, err := r.db.ExecContext(ctx, query, place, id)
	
	if err != nil {
		return fmt.Errorf("repository: failed to update participation place: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository: failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("participation not found")
	}
	
	return nil
}

// ----------------------------------------------------
// 5. DELETE (Снятие спортсмена с соревнования)
// ----------------------------------------------------
func (r *ParticipationRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM participations WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("repository: failed to delete participation: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository: failed to get rows affected after delete: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("participation not found")
	}
	
	return nil
}