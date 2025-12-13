package repository

import (
	"context"
	"database/sql"
	"fmt"
	// Возможно, нужен для других полей, оставлено для безопасности
)

// Participation представляет модель данных для таблицы participations
type Participation struct {
	ID            int `json:"id"`
	AthleteID     int `json:"athlete_id"`
	CompetitionID int `json:"competition_id"`
	Place         int `json:"place"` // 0 до определения результата

	// НОВЫЕ ПОЛЯ для отображения в HTML (получаются через JOIN)
	AthleteName     string `json:"athlete_name"`
	CompetitionName string `json:"competition_name"`
}

// ParticipationRepository отвечает за работу с таблицей participations
type ParticipationRepository struct {
	db *sql.DB
}

func NewParticipationRepository(db *sql.DB) *ParticipationRepository {
	return &ParticipationRepository{db: db}
}

// Create создает новую запись об участии
func (r *ParticipationRepository) Create(ctx context.Context, p *Participation) error {
	// PLACE здесь должен быть обнуляемым (NULL) в БД, если результат еще не известен.
	// Для простоты оставим Place=0, если его нет.
	if p.Place == 0 {
		// PostgreSQL умеет работать с NULL, но Go int не может быть NULL.
		// Если Place в БД nullable, используйте sql.NullInt32.
		// Предполагая, что Place в БД - NOT NULL и временно 0:
	}

	query := `
		INSERT INTO participations (athlete_id, competition_id, place)
		VALUES ($1, $2, $3)
		RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		p.AthleteID,
		p.CompetitionID,
		p.Place,
	).Scan(&p.ID)

	if err != nil {
		return fmt.Errorf("repository: failed to create participation: %w", err)
	}
	return nil
}

// ListAll получает список всех участий с именами спортсменов и соревнований
func (r *ParticipationRepository) ListAll(ctx context.Context) ([]Participation, error) {
	query := `
		SELECT 
			p.id, p.athlete_id, p.competition_id, p.place,
			a.full_name AS athlete_name,
			c.name      AS competition_name
		FROM participations p
		JOIN athletes a ON a.id = p.athlete_id
		JOIN competitions c ON c.id = p.competition_id
		ORDER BY p.id ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to query participations: %w", err)
	}
	defer rows.Close()

	participations := make([]Participation, 0)
	for rows.Next() {
		var p Participation

		// ВАЖНО: Добавить новые поля в Scan!
		err := rows.Scan(
			&p.ID, &p.AthleteID, &p.CompetitionID, &p.Place,
			&p.AthleteName, &p.CompetitionName,
		)
		if err != nil {
			return nil, fmt.Errorf("repository: failed to scan participation row: %w", err)
		}
		participations = append(participations, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return participations, nil
}

// UpdatePlace обновляет место в соревновании
func (r *ParticipationRepository) UpdatePlace(ctx context.Context, id int, place int) error {
	result, err := r.db.ExecContext(ctx, "UPDATE participations SET place = $2 WHERE id = $1", id, place)
	if err != nil {
		return fmt.Errorf("repository: failed to update place: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("participation not found")
	}
	return nil
}

// Delete удаляет запись об участии
func (r *ParticipationRepository) Delete(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM participations WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("repository: failed to delete participation: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("participation not found")
	}
	return nil
}
