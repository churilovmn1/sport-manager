package repository

import (
	"context"
	"database/sql"
	"fmt"
)

// User представляет модель данных для таблицы users
type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"` // <-- ИСПРАВЛЕНО: Теперь соответствует password_hash в БД
	Role         string `json:"role"`
}

// AuthRepository отвечает за работу с таблицей users
type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// GetByUsername получает пользователя по его имени
func (r *AuthRepository) GetByUsername(ctx context.Context, username string) (*User, error) {
	user := &User{}

	// Обновленный запрос: используем правильное имя столбца 'password_hash'
	query := `SELECT id, username, password_hash, role FROM users WHERE username = $1` // <-- ИСПРАВЛЕНО

	// Обновленный Scan: сканируем в &user.PasswordHash
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash, // <-- ИСПРАВЛЕНО
		&user.Role,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("repository: failed to get user by username: %w", err)
	}

	return user, nil
}

// CreateUser создает нового пользователя (включая роль)
// ВНИМАНИЕ: Перед вызовом этого метода пароль должен быть хеширован в сервисе!
func (r *AuthRepository) CreateUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (username, password_hash, role)
		VALUES ($1, $2, $3)
		RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		user.Username,
		user.PasswordHash, // <-- ИСПРАВЛЕНО
		user.Role,
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("repository: failed to create user: %w", err)
	}
	return nil
}
