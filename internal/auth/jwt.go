package auth

import (
	"fmt"
	"time"

	"sport-manager/internal/repository"
	"sport-manager/pkg/config"

	"github.com/golang-jwt/jwt/v5"
)

// Claims расширяет стандартные Claims, добавляя Username и Role
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken создает новый JWT для пользователя
func GenerateToken(user *repository.User, cfg *config.Config) (string, error) {
	// ИСПРАВЛЕНО: Доступ к публичному полю JWTExpirationDuration
	expirationTime := time.Now().Add(cfg.JWTExpirationDuration)

	claims := &Claims{
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))

	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken проверяет валидность токена
func ValidateToken(tokenString string, cfg *config.Config) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	return claims, nil
}
