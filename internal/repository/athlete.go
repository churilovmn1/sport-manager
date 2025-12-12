// internal/repository/athlete.go
package repository

import (
    "context"
    "database/sql"
)

type Athlete struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    // ... другие поля
}

type AthleteRepository struct {
    db *sql.DB
}

func NewAthleteRepository(db *sql.DB) *AthleteRepository {
    return &AthleteRepository{db: db}
}

func (r *AthleteRepository) Create(ctx context.Context, athlete *Athlete) error {
    query := `INSERT INTO athletes (name) VALUES ($1) RETURNING id`
    return r.db.QueryRowContext(ctx, query, athlete.Name).Scan(&athlete.ID)
}

// ... другие CRUD-методы (Get, Update, Delete)