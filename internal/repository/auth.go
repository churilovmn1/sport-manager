package repository

import (
	"context"
	"database/sql"
	"fmt"
)

// User представляет данные пользователя для аутентификации
type User struct {
	ID           int
	Username     string
	PasswordHash string
}

// AuthRepository отвечает за работу с таблицей users
type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// GetUserByUsername находит пользователя по имени и возвращает его данные, включая хеш пароля.
func (r *AuthRepository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	query := `SELECT id, username, password_hash FROM users WHERE username = $1`
	
	user := &User{}
	err := r.db.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &user.PasswordHash)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("database query failed: %w", err)
	}
	
	return user, nil
}