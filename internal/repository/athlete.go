package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Athlete представляет модель данных для таблицы athletes
type Athlete struct {
	ID        int       `json:"id"`
	FullName  string    `json:"full_name"`
	BirthDate time.Time `json:"birth_date"`
	Gender    string    `json:"gender"`
	SportID   int       `json:"sport_id"`
	RankID    int       `json:"rank_id"`
	Address   string    `json:"address"`
	IsActive  bool      `json:"is_active"`
}

// AthleteRepository отвечает за работу с таблицей athletes
type AthleteRepository struct {
	db *sql.DB
}

func NewAthleteRepository(db *sql.DB) *AthleteRepository {
	return &AthleteRepository{db: db}
}

// ----------------------------------------------------
// 1. CREATE
// ----------------------------------------------------
func (r *AthleteRepository) Create(ctx context.Context, a *Athlete) error {
	query := `
		INSERT INTO athletes (full_name, birth_date, gender, sport_id, rank_id, address, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		a.FullName, a.BirthDate, a.Gender, a.SportID, a.RankID, a.Address, a.IsActive,
	).Scan(&a.ID)

	if err != nil {
		return fmt.Errorf("repository: failed to create athlete: %w", err)
	}
	return nil
}

// ----------------------------------------------------
// 2. READ ONE
// ----------------------------------------------------
func (r *AthleteRepository) GetByID(ctx context.Context, id int) (*Athlete, error) {
	query := `
		SELECT id, full_name, birth_date, gender, sport_id, rank_id, address, is_active
		FROM athletes WHERE id = $1`

	athlete := &Athlete{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&athlete.ID, &athlete.FullName, &athlete.BirthDate, &athlete.Gender,
		&athlete.SportID, &athlete.RankID, &athlete.Address, &athlete.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("athlete not found")
		}
		return nil, fmt.Errorf("repository: failed to get athlete by ID: %w", err)
	}
	return athlete, nil
}

// ----------------------------------------------------
// 3. READ ALL (List)
// ----------------------------------------------------
func (r *AthleteRepository) GetAll(ctx context.Context) ([]*Athlete, error) {
	query := `
		SELECT id, full_name, birth_date, gender, sport_id, rank_id, address, is_active
		FROM athletes ORDER BY full_name`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to get all athletes: %w", err)
	}
	defer rows.Close()

	var athletes []*Athlete
	for rows.Next() {
		athlete := &Athlete{}
		err := rows.Scan(
			&athlete.ID, &athlete.FullName, &athlete.BirthDate, &athlete.Gender,
			&athlete.SportID, &athlete.RankID, &athlete.Address, &athlete.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("repository: failed to scan athlete row: %w", err)
		}
		athletes = append(athletes, athlete)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: rows error: %w", err)
	}
	return athletes, nil
}

// ----------------------------------------------------
// 4. UPDATE
// ----------------------------------------------------
func (r *AthleteRepository) Update(ctx context.Context, a *Athlete) error {
	query := `
		UPDATE athletes 
		SET full_name = $1, birth_date = $2, gender = $3, sport_id = $4, rank_id = $5, address = $6, is_active = $7
		WHERE id = $8`

	result, err := r.db.ExecContext(ctx, query,
		a.FullName, a.BirthDate, a.Gender, a.SportID, a.RankID, a.Address, a.IsActive, a.ID)

	if err != nil {
		return fmt.Errorf("repository: failed to update athlete: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository: failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("athlete not found")
	}

	return nil
}

// ----------------------------------------------------
// 5. DELETE
// ----------------------------------------------------
func (r *AthleteRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM athletes WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("repository: failed to delete athlete: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository: failed to get rows affected after delete: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("athlete not found")
	}

	return nil
}
