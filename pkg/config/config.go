package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config содержит все настройки приложения.
type Config struct {
	// Database settings
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     int    `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`

	// HTTP Server settings
	HTTPServerPort int `mapstructure:"HTTP_SERVER_PORT"`

	// JWT settings
	JWTSecret             string        `mapstructure:"JWT_SECRET"`
	JWTExpirationDuration time.Duration `mapstructure:"JWT_EXPIRATION_DURATION"`
}

// LoadConfig загружает настройки из файла или переменных окружения.
func LoadConfig() (*Config, error) {
	// 1. Устанавливаем пути поиска и тип файла
	viper.AddConfigPath(".")
	viper.SetConfigName("app") // Ищет файл с именем 'app.env'
	viper.SetConfigType("env")

	// 2. Устанавливаем настройки для чтения переменных окружения (более высокий приоритет, чем значения по умолчанию)
	viper.AutomaticEnv()

	// 3. Устанавливаем значения по умолчанию (самый низкий приоритет)
	// Эти значения будут использоваться, только если ни файл, ни переменные окружения не найдены.
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5432)
	viper.SetDefault("HTTP_SERVER_PORT", 8080)
	viper.SetDefault("JWT_EXPIRATION_DURATION", time.Hour*24)

	// 4. Читаем файл конфигурации
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
		// Файл не найден.

		// === КРИТИЧЕСКОЕ ИСПРАВЛЕНИЕ ПРИОРИТЕТА ===
		// Если файл не найден, мы явно устанавливаем нужные DB-значения.
		// Это гарантирует, что они перекроют любое системное имя пользователя Windows ("Makar"),
		// которое могло быть подобрано через viper.AutomaticEnv().
		viper.Set("DB_USER", "postgres")
		viper.Set("DB_PASSWORD", "student")
		viper.Set("DB_NAME", "sport_manager") // Также убедимся, что имя БД задано

		// Логирование, что файл не найден, но мы используем значения по умолчанию
		// log.Println("Config file not found, using environment variables and hardcoded defaults.")
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
