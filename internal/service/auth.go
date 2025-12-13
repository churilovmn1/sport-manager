package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
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
		log.Printf("DEBUG: User '%s' not found in repository or DB error: %v", username, err)
		return "", errors.New("invalid credentials")
	}

	// 1. ОЧИСТКА ХЕША
	cleanHash := strings.TrimSpace(user.PasswordHash)

	// --- ФИНАЛЬНЫЙ ДИАГНОСТИЧЕСКИЙ ЛОГ ---
	log.Printf("DEBUG: Found User: %s. Original Hash length: %d. Clean Hash length: %d. Clean Hash: [%s]",
		user.Username, len(user.PasswordHash), len(cleanHash), cleanHash)
	// ------------------------------------

	// 2. Сравнение пароля с ЧИСТЫМ хешем
	err = bcrypt.CompareHashAndPassword([]byte(cleanHash), []byte(password))
	if err != nil {
		log.Printf("DEBUG: bcrypt comparison FAILED for user %s. Error: %v", user.Username, err)
		return "", errors.New("invalid credentials")
	}

	// 3. Генерация JWT (Успешный вход)
	expirationTime := time.Now().Add(24 * time.Hour)
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
