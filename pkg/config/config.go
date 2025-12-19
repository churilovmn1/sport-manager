package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config объединяет все настройки приложения в одной структуре.
// Использование тегов mapstructure позволяет Viper автоматически сопоставлять
// ключи из файлов конфигурации (.env) с полями структуры.
type Config struct {
	// Настройки базы данных (PostgreSQL)
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     int    `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`

	// Настройки HTTP-сервера
	HTTPServerPort int `mapstructure:"HTTP_SERVER_PORT"`

	// Параметры безопасности JWT (JSON Web Token)
	JWTSecret             string        `mapstructure:"JWT_SECRET"`
	JWTExpirationDuration time.Duration `mapstructure:"JWT_EXPIRATION_DURATION"`
}

// LoadConfig инициализирует конфигурацию, соблюдая строгую иерархию приоритетов:
// 1. Переменные окружения (Environment Variables)
// 2. Файл конфигурации (app.env)
// 3. Значения по умолчанию (Defaults)
func LoadConfig() (*Config, error) {
	// Настройка поиска файла конфигурации
	viper.AddConfigPath(".")   // Искать в корневой папке проекта
	viper.SetConfigName("app") // Имя файла: app.env
	viper.SetConfigType("env") // Формат файла

	// Включаем автоматическое чтение переменных окружения.
	// Полезно при развертывании в Docker или на сервере.
	viper.AutomaticEnv()

	// Установка базовых значений (Defaults)
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5432)
	viper.SetDefault("HTTP_SERVER_PORT", 8080)
	viper.SetDefault("JWT_EXPIRATION_DURATION", time.Hour*24)

	// Попытка чтения файла app.env
	if err := viper.ReadInConfig(); err != nil {
		// Проверяем, является ли ошибка отсутствием файла
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}

		// Если файла нет, устанавливаем жесткие значения для учебной среды,
		// чтобы избежать конфликтов с системными переменными (например, DB_USER в Windows)
		viper.Set("DB_USER", "postgres")
		viper.Set("DB_PASSWORD", "student")
		viper.Set("DB_NAME", "sport_manager")
	}

	// Декодируем накопленные данные в структуру Config
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
