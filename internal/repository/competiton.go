package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Competition представляет модель данных для таблицы competitions
type Competition struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	StartDate time.Time `json:"start_date"`
}

// CompetitionRepository отвечает за работу с таблицей competitions
type CompetitionRepository struct {
	db *sql.DB
}

func NewCompetitionRepository(db *sql.DB) *CompetitionRepository {
	return &CompetitionRepository{db: db}
}

// Create создает новое соревнование
func (r *CompetitionRepository) Create(ctx context.Context, c *Competition) error {
	query := `
		INSERT INTO competitions (name, location, start_date)
		VALUES ($1, $2, $3)
		RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		c.Name, c.Location, c.StartDate,
	).Scan(&c.ID)

	if err != nil {
		return fmt.Errorf("repository: failed to create competition: %w", err)
	}
	return nil
}

// GetByID получает соревнование по его ID
func (r *CompetitionRepository) GetByID(ctx context.Context, id int) (*Competition, error) {
	query := `
		SELECT id, name, location, start_date
		FROM competitions WHERE id = $1`

	competition := &Competition{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&competition.ID, &competition.Name, &competition.Location, &competition.StartDate,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("competition not found")
		}
		return nil, fmt.Errorf("repository: failed to get competition by ID: %w", err)
	}
	return competition, nil
}

// ListAll получает список всех соревнований
func (r *CompetitionRepository) ListAll(ctx context.Context) ([]Competition, error) {
	query := `
		SELECT id, name, location, start_date
		FROM competitions ORDER BY start_date DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to query competitions: %w", err)
	}
	defer rows.Close()

	var competitions []Competition
	for rows.Next() {
		c := Competition{}
		if err := rows.Scan(&c.ID, &c.Name, &c.Location, &c.StartDate); err != nil {
			return nil, fmt.Errorf("repository: failed to scan competition row: %w", err)
		}
		competitions = append(competitions, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return competitions, nil
}

// Update обновляет существующее соревнование
func (r *CompetitionRepository) Update(ctx context.Context, c *Competition) error {
	query := `
		UPDATE competitions SET name=$1, location=$2, start_date=$3
		WHERE id=$4`

	result, err := r.db.ExecContext(ctx, query, c.Name, c.Location, c.StartDate, c.ID)
	if err != nil {
		return fmt.Errorf("repository: failed to update competition: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("competition not found")
	}
	return nil
}

// Delete удаляет соревнование по ID
func (r *CompetitionRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM competitions WHERE id=$1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("repository: failed to delete competition: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("competition not found")
	}
	return nil
}
