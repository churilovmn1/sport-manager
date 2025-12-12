package service

import (
	"context"
	"errors"
	"time"

	"sport-manager/internal/repository"
	"sport-manager/pkg/config"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// AuthService содержит методы для аутентификации и работы с токенами
type AuthService struct {
	repo      *repository.AuthRepository
	jwtSecret string
}

func NewAuthService(repo *repository.AuthRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		repo:      repo,
		jwtSecret: cfg.JWTSecret,
	}
}

// Claims определяет структуру полезной нагрузки JWT
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Login проверяет учетные данные и генерирует JWT-токен.
func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		// Скрываем детали об ошибке "user not found" для безопасности
		return "", errors.New("invalid credentials") 
	}

	// Сравниваем хеш пароля
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Генерация JWT
	expirationTime := time.Now().Add(24 * time.Hour) // Токен действителен 24 часа
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}