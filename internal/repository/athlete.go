package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// Athlete описывает структуру таблицы спортсменов в базе данных
type Athlete struct {
	ID        int       `json:"id"`
	FullName  string    `json:"full_name"`
	BirthDate time.Time `json:"birth_date"`
	Gender    string    `json:"gender"`
	IsActive  bool      `json:"is_active"`
	Address   string    `json:"address"`
}

// AthleteRepository предоставляет методы для взаимодействия с таблицей athletes
type AthleteRepository struct {
	db *sql.DB
}

// NewAthleteRepository создает новый экземпляр репозитория
func NewAthleteRepository(db *sql.DB) *AthleteRepository {
	return &AthleteRepository{db: db}
}

// --- МЕТОДЫ ДОСТУПА К ДАННЫМ ---

// Create добавляет нового спортсмена и возвращает присвоенный ID
func (r *AthleteRepository) Create(ctx context.Context, athlete *Athlete) error {
	query := `
		INSERT INTO athletes (full_name, birth_date, gender, is_active, address)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	// Используем QueryRowContext для получения сгенерированного ID через RETURNING
	err := r.db.QueryRowContext(ctx, query,
		athlete.FullName,
		athlete.BirthDate,
		athlete.Gender,
		athlete.IsActive,
		athlete.Address,
	).Scan(&athlete.ID)

	if err != nil {
		return fmt.Errorf("repo: не удалось создать атлета: %w", err)
	}
	return nil
}

// ListAll возвращает список всех спортсменов из базы данных
func (r *AthleteRepository) ListAll(ctx context.Context) ([]Athlete, error) {
	query := `
		SELECT id, full_name, birth_date, gender, is_active, address
		FROM athletes
		ORDER BY id ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repo: ошибка при получении списка: %w", err)
	}
	defer rows.Close() // Обязательно освобождаем ресурсы

	athletes := make([]Athlete, 0)
	for rows.Next() {
		var a Athlete
		if err := rows.Scan(&a.ID, &a.FullName, &a.BirthDate, &a.Gender, &a.IsActive, &a.Address); err != nil {
			return nil, fmt.Errorf("repo: ошибка сканирования строки: %w", err)
		}
		athletes = append(athletes, a)
	}

	// Проверяем, не возникло ли ошибок в процессе итерации
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repo: ошибка итерации строк: %w", err)
	}

	return athletes, nil
}

// GetByID находит одного спортсмена по его уникальному идентификатору
func (r *AthleteRepository) GetByID(ctx context.Context, id int) (*Athlete, error) {
	a := &Athlete{}
	query := `
		SELECT id, full_name, birth_date, gender, is_active, address 
		FROM athletes
		WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&a.ID, &a.FullName, &a.BirthDate, &a.Gender, &a.IsActive, &a.Address,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("атлет с ID %d не найден", id)
		}
		return nil, fmt.Errorf("repo: ошибка при поиске по ID: %w", err)
	}

	return a, nil
}

// Update обновляет все поля существующего спортсмена по его ID
func (r *AthleteRepository) Update(ctx context.Context, a *Athlete) error {
	query := `
		UPDATE athletes 
		SET full_name = $2, birth_date = $3, gender = $4, is_active = $5, address = $6
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		a.ID, a.FullName, a.BirthDate, a.Gender, a.IsActive, a.Address,
	)

	if err != nil {
		return fmt.Errorf("repo: не удалось обновить данные: %w", err)
	}
	return nil
}

// Delete безвозвратно удаляет запись о спортсмене из базы данных
func (r *AthleteRepository) Delete(ctx context.Context, id int) error {
	log.Printf("Запрос на удаление атлета ID: %d", id)

	result, err := r.db.ExecContext(ctx, "DELETE FROM athletes WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("repo: критическая ошибка при удалении: %w", err)
	}

	// Проверяем, была ли удалена хоть одна строка
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("атлет с ID %d не найден для удаления", id)
	}

	log.Printf("Успешно удалено строк: %d", rowsAffected)
	return nil
}
