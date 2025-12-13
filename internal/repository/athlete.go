package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// Athlete представляет модель данных для таблицы athletes
type Athlete struct {
	ID        int       `json:"id"`
	FullName  string    `json:"full_name"`
	BirthDate time.Time `json:"birth_date"`
	Gender    string    `json:"gender"`
	IsActive  bool      `json:"is_active"`
	Address   string    `json:"address"`
}

type AthleteRepository struct {
	db *sql.DB
}

func NewAthleteRepository(db *sql.DB) *AthleteRepository {
	return &AthleteRepository{db: db}
}

// Create создает нового спортсмена
func (r *AthleteRepository) Create(ctx context.Context, athlete *Athlete) error {
	// SQL-запрос без sport_id и rank_id
	query := `
		INSERT INTO athletes (full_name, birth_date, gender, is_active, address)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		athlete.FullName,
		athlete.BirthDate,
		athlete.Gender,
		athlete.IsActive,
		athlete.Address,
	).Scan(&athlete.ID)

	if err != nil {
		return fmt.Errorf("repository: failed to create athlete: %w", err)
	}
	return nil
}

// ListAll получает список всех спортсменов
func (r *AthleteRepository) ListAll(ctx context.Context) ([]Athlete, error) {
	query := `
		SELECT 
			id, full_name, birth_date, gender, is_active, address
		FROM athletes
		ORDER BY id ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to query athletes: %w", err)
	}
	defer rows.Close()

	athletes := make([]Athlete, 0)
	for rows.Next() {
		var a Athlete

		err := rows.Scan(
			&a.ID, &a.FullName, &a.BirthDate, &a.Gender, &a.IsActive, &a.Address,
		)
		if err != nil {
			return nil, fmt.Errorf("repository: failed to scan athlete row: %w", err)
		}

		athletes = append(athletes, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return athletes, nil
}

// GetByID получает спортсмена по ID
func (r *AthleteRepository) GetByID(ctx context.Context, id int) (*Athlete, error) {
	athlete := &Athlete{}
	query := `
		SELECT 
			id, full_name, birth_date, gender, is_active, address 
		FROM athletes
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&athlete.ID, &athlete.FullName, &athlete.BirthDate, &athlete.Gender,
		&athlete.IsActive, &athlete.Address,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("athlete not found")
		}
		return nil, fmt.Errorf("repository: failed to get athlete by ID: %w", err)
	}

	return athlete, nil
}

// Update обновляет данные спортсмена
func (r *AthleteRepository) Update(ctx context.Context, athlete *Athlete) error {
	// SQL-запрос без rank_id и sport_id
	query := `
		UPDATE athletes SET full_name = $2, birth_date = $3, gender = $4,
		is_active = $5, address = $6
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		athlete.ID,
		athlete.FullName,
		athlete.BirthDate,
		athlete.Gender,
		athlete.IsActive,
		athlete.Address,
	)
	if err != nil {
		return fmt.Errorf("repository: failed to update athlete: %w", err)
	}
	return nil
}

// Delete удаляет спортсмена
func (r *AthleteRepository) Delete(ctx context.Context, id int) error {
	log.Printf("Executing DELETE for athlete ID: %d", id)
	result, err := r.db.ExecContext(ctx, "DELETE FROM athletes WHERE id = $1", id)

	if err != nil {
		log.Printf("SQL Error during delete for ID %d: %v", id, err)
		return fmt.Errorf("repository: failed to delete athlete: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		log.Printf("Athlete ID %d not found for deletion.", id)
		return fmt.Errorf("athlete not found")
	}

	log.Printf("Delete successful, rows affected: %d", rowsAffected)
	return nil
}
