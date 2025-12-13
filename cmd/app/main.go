package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"sport-manager/internal/auth"
	"sport-manager/internal/handler"
	"sport-manager/internal/repository"
	"sport-manager/internal/service"
	"sport-manager/pkg/config"
	"sport-manager/pkg/utils"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// 1. Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("FATAL: Failed to load config: %v", err)
	}

	// --- 1.1. Запуск миграций ---
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	log.Println("Attempting to run database migrations...")
	// NOTE: Убедитесь, что миграции, добавляющие "place" и удаляющие "sport_id", выполнены!
	if err := utils.RunMigrations(dbURL, "migrations"); err != nil {
		log.Fatalf("FATAL: Database migration failed: %v", err)
	}
	log.Println("Migrations executed successfully.")

	// 2. Инициализация подключения к БД
	db, err := initDB(cfg)
	if err != nil {
		log.Fatalf("FATAL: Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("SUCCESS: Successfully connected to the PostgreSQL database.")

	// 3. Инициализация Repositories, Services, Handlers

	// --- Repositories ---
	authRepo := repository.NewAuthRepository(db)
	athleteRepo := repository.NewAthleteRepository(db)
	competitionRepo := repository.NewCompetitionRepository(db)
	participationRepo := repository.NewParticipationRepository(db)

	// --- Services ---
	authService := service.NewAuthService(authRepo, cfg)
	athleteService := service.NewAthleteService(athleteRepo)
	competitionService := service.NewCompetitionService(competitionRepo)
	participationService := service.NewParticipationService(participationRepo, athleteRepo, competitionRepo)

	// --- Handlers ---
	authHandler := handler.NewAuthHandler(authService)
	athleteHandler := handler.NewAthleteHandler(athleteService)
	competitionHandler := handler.NewCompetitionHandler(competitionService)
	participationHandler := handler.NewParticipationHandler(participationService)

	// 4. Настройка роутера (mux)
	router := mux.NewRouter()

	// --- 4.1. Публичные маршруты ---
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Service is up and running!"))
	}).Methods("GET")
	router.HandleFunc("/api/v1/auth/login", authHandler.Login).Methods("POST")

	// NOTE: Удалены явные PathPrefix для .html файлов (participants.html, athletes.html, competitions.html).
	// Они будут обслуживаться универсальным FileServer ниже.

	// --- 4.2. Защищенные маршруты (API CRUD) ---
	apiRouter := router.PathPrefix("/api/v1").Subrouter()
	// NOTE: Если вы не залогинены, API запросы будут возвращать 401.
	// Убедитесь, что AuthMiddleware (auth.AuthMiddleware(cfg)) работает корректно.
	apiRouter.Use(auth.AuthMiddleware(cfg))

	// CRUD: Спортсмены (Athletes)
	apiRouter.HandleFunc("/athletes", athleteHandler.CreateAthlete).Methods("POST")
	apiRouter.HandleFunc("/athletes", athleteHandler.ListAllAthletes).Methods("GET")
	apiRouter.HandleFunc("/athletes/{id}", athleteHandler.GetAthleteByID).Methods("GET")
	apiRouter.HandleFunc("/athletes/{id}", athleteHandler.UpdateAthlete).Methods("PUT")
	apiRouter.HandleFunc("/athletes/{id}", athleteHandler.DeleteAthlete).Methods("DELETE")

	// CRUD: Соревнования (Competitions)
	apiRouter.HandleFunc("/competitions", competitionHandler.CreateCompetition).Methods("POST")
	apiRouter.HandleFunc("/competitions", competitionHandler.ListCompetitions).Methods("GET")
	apiRouter.HandleFunc("/competitions/{id}", competitionHandler.GetCompetition).Methods("GET")
	apiRouter.HandleFunc("/competitions/{id}", competitionHandler.UpdateCompetition).Methods("PUT")
	apiRouter.HandleFunc("/competitions/{id}", competitionHandler.DeleteCompetition).Methods("DELETE")

	// CRUD: Участие/Результаты (Participations)
	apiRouter.HandleFunc("/participations", participationHandler.CreateParticipation).Methods("POST")
	apiRouter.HandleFunc("/participations", participationHandler.ListParticipations).Methods("GET")
	apiRouter.HandleFunc("/participations/{id}/place", participationHandler.UpdatePlace).Methods("PUT")
	apiRouter.HandleFunc("/participations/{id}", participationHandler.DeleteParticipation).Methods("DELETE")

	// CRUD: Справочники (Sports)
	
	

	// CRUD: Справочники (Ranks)
	
	

	// -----------------------------------------------------------------
	// 4.3. ОБСЛУЖИВАНИЕ СТАТИЧЕСКИХ ФАЙЛОВ (ОДИН УНИВЕРСАЛЬНЫЙ МАРШРУТ)
	// Этот маршрут должен быть ПОСЛЕДНИМ, чтобы не перехватывать /api/v1.
	// -----------------------------------------------------------------
	fileServer := http.FileServer(http.Dir("./web/"))
	router.PathPrefix("/").Handler(fileServer)

	// 5. Запуск сервера
	log.Printf("Starting HTTP server on :%d", cfg.HTTPServerPort)
	serverAddr := fmt.Sprintf(":%d", cfg.HTTPServerPort)
	if err := http.ListenAndServe(serverAddr, router); err != nil && err != http.ErrServerClosed {
		log.Fatalf("FATAL: Could not listen on %s: %v", serverAddr, err)
	}
}

// initDB настраивает и проверяет подключение к PostgreSQL.
func initDB(cfg *config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("open database connection: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
