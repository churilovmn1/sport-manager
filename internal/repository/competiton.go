package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Competition описывает модель турнира или спортивного события.
type Competition struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	StartDate time.Time `json:"start_date"`
}

// CompetitionRepository инкапсулирует логику работы с таблицей соревнований.
type CompetitionRepository struct {
	db *sql.DB
}

// NewCompetitionRepository создает новый экземпляр репозитория соревнований.
func NewCompetitionRepository(db *sql.DB) *CompetitionRepository {
	return &CompetitionRepository{db: db}
}

// --- МЕТОДЫ РАБОТЫ С БАЗОЙ ДАННЫХ ---

// Create сохраняет новое соревнование и возвращает сгенерированный базой ID.
func (r *CompetitionRepository) Create(ctx context.Context, c *Competition) error {
	query := `
		INSERT INTO competitions (name, location, start_date)
		VALUES ($1, $2, $3)
		RETURNING id`

	// Используем QueryRowContext для безопасного выполнения в рамках контекста запроса
	err := r.db.QueryRowContext(ctx, query,
		c.Name, c.Location, c.StartDate,
	).Scan(&c.ID)

	if err != nil {
		return fmt.Errorf("repo: ошибка при создании соревнования: %w", err)
	}
	return nil
}

// GetByID возвращает данные конкретного соревнования по его первичному ключу.
func (r *CompetitionRepository) GetByID(ctx context.Context, id int) (*Competition, error) {
	query := `
		SELECT id, name, location, start_date
		FROM competitions WHERE id = $1`

	c := &Competition{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&c.ID, &c.Name, &c.Location, &c.StartDate,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("соревнование с ID %d не найдено", id)
		}
		return nil, fmt.Errorf("repo: ошибка при получении данных: %w", err)
	}
	return c, nil
}

// ListAll возвращает полный список соревнований, отсортированный по дате начала (сначала новые).
func (r *CompetitionRepository) ListAll(ctx context.Context) ([]Competition, error) {
	query := `
		SELECT id, name, location, start_date
		FROM competitions ORDER BY start_date DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repo: ошибка запроса списка: %w", err)
	}
	defer rows.Close() // Закрываем соединение после обработки всех строк

	var competitions []Competition
	for rows.Next() {
		c := Competition{}
		if err := rows.Scan(&c.ID, &c.Name, &c.Location, &c.StartDate); err != nil {
			return nil, fmt.Errorf("repo: ошибка сканирования данных: %w", err)
		}
		competitions = append(competitions, c)
	}

	// Проверяем финальное состояние итератора на наличие скрытых ошибок
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repo: ошибка обработки строк: %w", err)
	}
	return competitions, nil
}

// Update изменяет информацию о соревновании. Возвращает ошибку, если запись не найдена.
func (r *CompetitionRepository) Update(ctx context.Context, c *Competition) error {
	query := `
		UPDATE competitions 
		SET name = $1, location = $2, start_date = $3
		WHERE id = $4`

	result, err := r.db.ExecContext(ctx, query, c.Name, c.Location, c.StartDate, c.ID)
	if err != nil {
		return fmt.Errorf("repo: ошибка обновления данных: %w", err)
	}

	// Важный шаг: проверяем, затронул ли запрос хоть одну строку
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("соревнование с ID %d не найдено", c.ID)
	}
	return nil
}

// Delete удаляет соревнование по его ID.
func (r *CompetitionRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM competitions WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("repo: ошибка при удалении записи: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("соревнование с ID %d не найдено для удаления", id)
	}
	return nil
}
