package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config содержит все настройки приложения.
type Config struct {
	// Настройки сервера
	HTTPServerPort int

	// Настройки базы данных (PostgreSQL)
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string

	// Настройки безопасности
	JWTSecret string
}

// LoadConfig загружает конфигурацию из переменных окружения (или .env файла).
func LoadConfig() (*Config, error) {
	// Загружаем .env файл для локальной разработки
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables.")
	}

	cfg := &Config{
		// Значения по умолчанию для порта, если не указаны
		HTTPServerPort: getIntEnv("HTTP_PORT", 8080),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getIntEnv("DB_PORT", 5432),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "sport_manager_db"),

		// Секретный ключ для JWT. Крайне важно, чтобы он был установлен.
		JWTSecret: getEnv("JWT_SECRET", "super_secret_key_change_me_in_production"),
	}

	return cfg, nil
}

// Вспомогательная функция для получения строковой переменной окружения
func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Вспомогательная функция для получения целочисленной переменной окружения
func getIntEnv(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Warning: Invalid integer value for %s. Using default: %d", key, defaultValue)
		return defaultValue
	}
	return value
}
