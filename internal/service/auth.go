package service

import (
	"context"
	"errors"
	"fmt"

	"sport-manager/internal/auth"
	"sport-manager/internal/repository"
	"sport-manager/pkg/config"

	"golang.org/x/crypto/bcrypt" // <-- НОВЫЙ ИМПОРТ ДЛЯ ХЕШИРОВАНИЯ
)

// AuthService содержит бизнес-логику для аутентификации
type AuthService struct {
	repo *repository.AuthRepository
	cfg  *config.Config
}

func NewAuthService(repo *repository.AuthRepository, cfg *config.Config) *AuthService {
	return &AuthService{repo: repo, cfg: cfg}
}

// Login проверяет учетные данные и генерирует JWT-токен
func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	if username == "" || password == "" {
		return "", fmt.Errorf("username and password cannot be empty")
	}

	// 1. Получение пользователя из репозитория
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		// Не раскрываем, что пользователь не найден.
		return "", fmt.Errorf("invalid credentials")
	}

	// 2. КРИТИЧЕСКОЕ ИСПРАВЛЕНИЕ: Проверка хеша bcrypt
	// Сравниваем введенный пароль с хешем, хранящимся в базе данных.
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))

	if err != nil {
		// Если ошибка - это несовпадение хеша, возвращаем общую ошибку
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", fmt.Errorf("invalid credentials")
		}
		// Для других ошибок (например, неверный формат хеша)
		return "", fmt.Errorf("authentication failed: %w", err)
	}

	// 3. Генерация токена
	token, err := auth.GenerateToken(user, s.cfg)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}
