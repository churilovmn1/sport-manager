package repository

import (
	"context"
	"database/sql"
	"fmt"
)

// SimpleModel представляет структуру для справочников (Sport и Rank)
type SimpleModel struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	// Description только для Rank
	Description sql.NullString `json:"description,omitempty"` 
}

// SimpleRepository общий репозиторий для справочников
type SimpleRepository struct {
	db        *sql.DB
	tableName string
}

func NewSimpleRepository(db *sql.DB, tableName string) *SimpleRepository {
	return &SimpleRepository{db: db, tableName: tableName}
}

// ----------------------------------------------------
// 1. CREATE
// ----------------------------------------------------
func (r *SimpleRepository) Create(ctx context.Context, model *SimpleModel) error {
	var query string
	var args []interface{}
	
	if r.tableName == "ranks" {
		query = fmt.Sprintf("INSERT INTO %s (name, description) VALUES ($1, $2) RETURNING id", r.tableName)
		args = []interface{}{model.Name, model.Description}
	} else {
		query = fmt.Sprintf("INSERT INTO %s (name) VALUES ($1) RETURNING id", r.tableName)
		args = []interface{}{model.Name}
	}
	
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&model.ID)
	
	if err != nil {
		return fmt.Errorf("repository: failed to create %s: %w", r.tableName, err)
	}
	return nil
}

// ----------------------------------------------------
// 2. READ ALL (List)
// ----------------------------------------------------
func (r *SimpleRepository) GetAll(ctx context.Context) ([]*SimpleModel, error) {
	selectFields := "id, name"
	if r.tableName == "ranks" {
		selectFields += ", description"
	}
	
	query := fmt.Sprintf("SELECT %s FROM %s ORDER BY name", selectFields, r.tableName)
	
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to get all %s: %w", r.tableName, err)
	}
	defer rows.Close()

	var models []*SimpleModel
	for rows.Next() {
		model := &SimpleModel{}
		var scanArgs []interface{} = []interface{}{&model.ID, &model.Name}
		
		if r.tableName == "ranks" {
			scanArgs = append(scanArgs, &model.Description)
		}
		
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, fmt.Errorf("repository: failed to scan %s row: %w", r.tableName, err)
		}
		models = append(models, model)
	}
	return models, rows.Err()
}

// ----------------------------------------------------
// 3. DELETE
// ----------------------------------------------------
func (r *SimpleRepository) Delete(ctx context.Context, id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", r.tableName)
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("repository: failed to delete %s: %w", r.tableName, err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("%s not found", r.tableName)
	}
	
	return nil
}