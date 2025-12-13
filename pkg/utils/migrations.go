package utils

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // PostgreSQL драйвер
	_ "github.com/golang-migrate/migrate/v4/source/file"       // Файловый источник миграций
)

// RunMigrations запускает миграции в указанной базе данных.
func RunMigrations(dbURL, migrationsPath string) error {
	// Формируем полный URL для мигратора.
	// DB URL должен содержать sslmode=disable.
	m, err := migrate.New(
		"file://"+migrationsPath,
		dbURL,
	)
	if err != nil {
		return fmt.Errorf("could not create migration instance: %w", err)
	}

	// Запускаем все доступные миграции
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	// Проверяем версию миграций и логируем результат
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNoChange {
		log.Printf("Warning: Could not get migration version: %v", err)
	} else if dirty {
		log.Printf("ATTENTION: Database is in a dirty state at version %d", version)
	} else {
		log.Printf("Database migration successful. Current version: %d", version)
	}

	return nil
}
