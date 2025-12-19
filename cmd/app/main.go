package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"sport-manager/internal/auth"
	"sport-manager/internal/handler"
	"sport-manager/internal/repository"
	"sport-manager/internal/service"
	"sport-manager/pkg/config"
	"sport-manager/pkg/utils"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // Драйвер PostgreSQL
)

func main() {
	// 1. ЗАГРУЗКА КОНФИГУРАЦИИ
	// Читаем переменные окружения или конфиг-файл
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Критическая ошибка: конфигурация не загружена: %v", err)
	}

	// 2. БАЗА ДАННЫХ И МИГРАЦИИ
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	log.Println("Выполнение миграций (создание таблиц)...")
	if err := utils.RunMigrations(dbURL, "migrations"); err != nil {
		log.Fatalf("Ошибка миграций: %v", err)
	}

	db, err := initDB(cfg)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	// 3. ИНИЦИАЛИЗАЦИЯ СЛОЕВ (Dependency Injection)
	// Инициализируем репозитории (работа с БД)
	authRepo := repository.NewAuthRepository(db)
	athleteRepo := repository.NewAthleteRepository(db)
	competitionRepo := repository.NewCompetitionRepository(db)
	participationRepo := repository.NewParticipationRepository(db)

	// Инициализируем сервисы (бизнес-логика)
	authService := service.NewAuthService(authRepo, cfg)
	athleteService := service.NewAthleteService(athleteRepo)
	competitionService := service.NewCompetitionService(competitionRepo)
	// ParticipationService зависит от нескольких репозиториев для проверки существования записей
	participationService := service.NewParticipationService(participationRepo, athleteRepo, competitionRepo)

	// Инициализируем хендлеры (обработка HTTP запросов)
	authHandler := handler.NewAuthHandler(authService)
	athleteHandler := handler.NewAthleteHandler(athleteService)
	competitionHandler := handler.NewCompetitionHandler(competitionService)
	participationHandler := handler.NewParticipationHandler(participationService)

	// 4. НАСТРОЙКА РОУТИНГА
	router := mux.NewRouter()

	// Группа для API v1
	api := router.PathPrefix("/api/v1").Subrouter()

	// --- Публичные маршруты ---
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	api.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// --- Защищенные маршруты (JWT Middleware) ---
	protected := api.PathPrefix("").Subrouter()
	protected.Use(auth.AuthMiddleware(cfg))

	// Спортсмены
	protected.HandleFunc("/athletes", athleteHandler.ListAllAthletes).Methods("GET")
	protected.HandleFunc("/athletes/{id}", athleteHandler.GetAthleteByID).Methods("GET")
	protected.HandleFunc("/athletes", auth.AdminOnly(athleteHandler.CreateAthlete)).Methods("POST")
	protected.HandleFunc("/athletes/{id}", auth.AdminOnly(athleteHandler.UpdateAthlete)).Methods("PUT")
	protected.HandleFunc("/athletes/{id}", auth.AdminOnly(athleteHandler.DeleteAthlete)).Methods("DELETE")

	// Соревнования
	protected.HandleFunc("/competitions", competitionHandler.ListCompetitions).Methods("GET")
	protected.HandleFunc("/competitions/{id}", competitionHandler.GetCompetition).Methods("GET")
	protected.HandleFunc("/competitions", auth.AdminOnly(competitionHandler.CreateCompetition)).Methods("POST")
	protected.HandleFunc("/competitions/{id}", auth.AdminOnly(competitionHandler.UpdateCompetition)).Methods("PUT")
	protected.HandleFunc("/competitions/{id}", auth.AdminOnly(competitionHandler.DeleteCompetition)).Methods("DELETE")

	// Участие (регистрация атлетов на турниры)
	protected.HandleFunc("/participations", participationHandler.ListParticipations).Methods("GET")
	protected.HandleFunc("/participations", auth.AdminOnly(participationHandler.CreateParticipation)).Methods("POST")
	protected.HandleFunc("/participations/{id}/place", auth.AdminOnly(participationHandler.UpdatePlace)).Methods("PUT")
	protected.HandleFunc("/participations/{id}", auth.AdminOnly(participationHandler.DeleteParticipation)).Methods("DELETE")

	// 5. РАЗДАЧА ФРОНТЕНДА
	// Важно: Static файлы регистрируются ПОСЛЕ API, чтобы не перекрывать маршруты
	fileServer := http.FileServer(http.Dir("./web/"))
	router.PathPrefix("/").Handler(fileServer)

	// 6. ЗАПУСК СЕРВЕРА
	serverAddr := fmt.Sprintf(":%d", cfg.HTTPServerPort)
	log.Printf("=== Sport Manager Backend запущен на %s ===", serverAddr)

	if err := http.ListenAndServe(serverAddr, router); err != nil {
		log.Fatalf("Ошибка остановки сервера: %v", err)
	}
}

// initDB открывает соединение и проверяет доступность базы
func initDB(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия базы: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("база недоступна: %w", err)
	}

	return db, nil
}
