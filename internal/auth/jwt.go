package auth

import (
	"fmt"
	"time"

	"sport-manager/internal/repository"
	"sport-manager/pkg/config"

	"github.com/golang-jwt/jwt/v5"
)

// Claims — структура полезной нагрузки токена.
// Мы расширяем стандартные поля JWT (RegisteredClaims), добавляя имя и роль,
// чтобы не лезть в базу данных при каждом запросе.
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken создает новый JWT-токен для пользователя.
// Токен — это зашифрованная строка, которая подтверждает личность пользователя.
func GenerateToken(user *repository.User, cfg *config.Config) (string, error) {
	// Устанавливаем время жизни токена из конфигурации
	expirationTime := time.Now().Add(cfg.JWTExpirationDuration)

	// Заполняем данные (Payload)
	claims := &Claims{
		Username: user.Username,
		Role:     user.Role, // Роль записывается в токен для работы Middleware
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Subject:   user.Username,
			IssuedAt:  jwt.NewNumericDate(time.Now()), // Время выдачи токена
		},
	}

	// Создаем токен с использованием алгоритма шифрования HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен секретным ключом из конфига
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("ошибка подписи токена: %w", err)
	}

	return tokenString, nil
}

// ValidateToken проверяет подлинность токена и его срок годности.
func ValidateToken(tokenString string, cfg *config.Config) (*Claims, error) {
	claims := &Claims{}

	// Парсим строку токена и проверяем подпись
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Проверка метода шифрования (защита от подмены алгоритма)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		// Возвращаем секретный ключ для проверки подписи
		return []byte(cfg.JWTSecret), nil
	})

	// Если токен поврежден или подпись не совпала
	if err != nil {
		return nil, fmt.Errorf("невалидный токен: %w", err)
	}

	// Если срок годности истек
	if !token.Valid {
		return nil, fmt.Errorf("токен просрочен или недействителен")
	}

	return claims, nil
}
