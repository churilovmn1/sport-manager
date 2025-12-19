package service

import (
	"context"
	"fmt"
	"log"
	"sport-manager/internal/auth"
	"sport-manager/internal/repository"
	"sport-manager/pkg/config"

	"golang.org/x/crypto/bcrypt"
)

// AuthService содержит бизнес-логику управления пользователями и сессиями.
type AuthService struct {
	repo *repository.AuthRepository
	cfg  *config.Config
}

// NewAuthService создает новый экземпляр сервиса аутентификации.
func NewAuthService(repo *repository.AuthRepository, cfg *config.Config) *AuthService {
	return &AuthService{repo: repo, cfg: cfg}
}

// Login проверяет учетные данные пользователя и возвращает JWT-токен при успехе.
func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	log.Printf("Попытка входа пользователя: %s", username)

	// 1. Ищем пользователя в базе данных
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		log.Printf("Ошибка входа: пользователь %s не найден", username)
		return "", fmt.Errorf("неверные учетные данные")
	}

	// 2. Проверяем соответствие введенного пароля сохраненному хешу
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		log.Printf("Ошибка входа: неверный пароль для %s", username)
		return "", fmt.Errorf("неверные учетные данные")
	}

	// 3. Генерируем токен доступа на основе данных пользователя
	token, err := auth.GenerateToken(user, s.cfg)
	if err != nil {
		log.Printf("Ошибка генерации токена: %v", err)
		return "", fmt.Errorf("ошибка сервера при авторизации")
	}

	log.Printf("Успешный вход: %s (Роль: %s)", username, user.Role)
	return token, nil
}

// Register выполняет безопасную регистрацию нового пользователя.
func (s *AuthService) Register(ctx context.Context, username, email, password string) error {
	if len(password) < 6 {
		return fmt.Errorf("пароль слишком короткий (минимум 6 символов)")
	}

	// Хешируем пароль перед сохранением в базу данных.
	// Никогда не храните пароли в открытом виде!
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("ошибка при обработке пароля: %w", err)
	}

	user := &repository.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         "user", // Роль по умолчанию
	}

	return s.repo.CreateUser(ctx, user)
}
