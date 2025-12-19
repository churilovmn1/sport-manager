package repository

import (
	"context"
	"database/sql"
	"fmt"
)

// Participation представляет собой связующую сущность (Many-to-Many)
// между атлетами и соревнованиями, включая результат (место).
type Participation struct {
	ID            int `json:"id"`
	AthleteID     int `json:"athlete_id"`
	CompetitionID int `json:"competition_id"`
	Place         int `json:"place"` // 0 означает, что соревнование еще не завершено

	// Поля, заполняемые через JOIN для удобства отображения на фронтенде
	AthleteName     string `json:"athlete_name"`
	CompetitionName string `json:"competition_name"`
}

// ParticipationRepository управляет записями об участии спортсменов в турнирах.
type ParticipationRepository struct {
	db *sql.DB
}

// NewParticipationRepository создает новый экземпляр репозитория участий.
func NewParticipationRepository(db *sql.DB) *ParticipationRepository {
	return &ParticipationRepository{db: db}
}

// --- МЕТОДЫ РАБОТЫ С ДАННЫМИ ---

// Create регистрирует атлета на соревнование.
func (r *ParticipationRepository) Create(ctx context.Context, p *Participation) error {
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
		return fmt.Errorf("repo: не удалось создать запись об участии: %w", err)
	}
	return nil
}

// ListAll извлекает все записи об участии, объединяя данные из таблиц athletes и competitions.
func (r *ParticipationRepository) ListAll(ctx context.Context) ([]Participation, error) {
	// Используем JOIN, чтобы сразу получить читаемые названия вместо простых ID
	query := `
		SELECT 
			p.id, p.athlete_id, p.competition_id, p.place,
			a.full_name AS athlete_name,
			c.name      AS competition_name
		FROM participations p
		JOIN athletes a ON a.id = p.athlete_id
		JOIN competitions c ON c.id = p.competition_id
		ORDER BY p.id ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repo: ошибка при выполнении JOIN-запроса: %w", err)
	}
	defer rows.Close()

	participations := make([]Participation, 0)
	for rows.Next() {
		var p Participation
		// Сканируем все поля, включая полученные через JOIN
		err := rows.Scan(
			&p.ID, &p.AthleteID, &p.CompetitionID, &p.Place,
			&p.AthleteName, &p.CompetitionName,
		)
		if err != nil {
			return nil, fmt.Errorf("repo: ошибка сканирования строки участия: %w", err)
		}
		participations = append(participations, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repo: ошибка после завершения чтения строк: %w", err)
	}

	return participations, nil
}

// UpdatePlace обновляет только результат (место) атлета в соревновании.
func (r *ParticipationRepository) UpdatePlace(ctx context.Context, id int, place int) error {
	result, err := r.db.ExecContext(ctx,
		"UPDATE participations SET place = $2 WHERE id = $1",
		id, place,
	)
	if err != nil {
		return fmt.Errorf("repo: не удалось обновить результат: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("запись об участии с ID %d не найдена", id)
	}
	return nil
}

// Delete удаляет запись об участии (аннулирует регистрацию атлета).
func (r *ParticipationRepository) Delete(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM participations WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("repo: ошибка при удалении записи: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("запись об участии с ID %d не найдена", id)
	}
	return nil
}
