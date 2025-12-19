package repository

import (
	"context"
	"database/sql"
	"fmt"
)

// User описывает модель пользователя системы.
// Поле PasswordHash помечено тегом `json:"-"`, чтобы оно никогда не попало в API ответы.
type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"` // Скрываем хеш пароля при сериализации
	Role         string `json:"role"`
}

// AuthRepository управляет хранением и поиском учетных данных пользователей.
type AuthRepository struct {
	db *sql.DB
}

// NewAuthRepository инициализирует репозиторий для работы с аутентификацией.
func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// --- МЕТОДЫ РАБОТЫ С ПОЛЬЗОВАТЕЛЯМИ ---

// GetByUsername ищет пользователя по имени для процесса логина.
func (r *AuthRepository) GetByUsername(ctx context.Context, username string) (*User, error) {
	user := &User{}
	query := `
		SELECT id, username, email, password_hash, role 
		FROM users 
		WHERE username = $1`

	// Выполняем запрос с использованием контекста для контроля времени выполнения
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("пользователь '%s' не найден", username)
		}
		return nil, fmt.Errorf("repo: ошибка при поиске пользователя: %w", err)
	}

	return user, nil
}

// CreateUser регистрирует нового пользователя в системе.
func (r *AuthRepository) CreateUser(ctx context.Context, user *User) error {
	// Устанавливаем роль по умолчанию, если она не задана
	if user.Role == "" {
		user.Role = "user"
	}

	query := `
		INSERT INTO users (username, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	// RETURNING id позволяет нам сразу обновить структуру user после вставки
	err := r.db.QueryRowContext(ctx, query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Role,
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("repo: не удалось создать пользователя: %w", err)
	}

	return nil
}
