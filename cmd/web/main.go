package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "time"
    
    "sport-manager/internal/config"
    "sport-manager/internal/handler"
    "sport-manager/internal/middleware"
    "sport-manager/internal/pkg/database"
    "sport-manager/internal/repository/postgres"
    "sport-manager/internal/service"
    "github.com/gorilla/mux"
    "github.com/rs/cors"
)

func main() {
    // Загрузка конфигурации
    cfg := config.Load()
    
    // Подключение к БД
    db, err := database.NewPostgres(cfg.DatabaseURL)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()
    
    // Инициализация репозиториев
    userRepo := postgres.NewUserRepository(db)
    athleteRepo := postgres.NewAthleteRepository(db)
    competitionRepo := postgres.NewCompetitionRepository(db)
    // ... другие репозитории
    
    // Инициализация сервисов
    authService := service.NewAuthService(userRepo, cfg.JWTSecret)
    athleteService := service.NewAthleteService(athleteRepo)
    competitionService := service.NewCompetitionService(competitionRepo)
    // ... другие сервисы
    
    // Инициализация обработчиков
    authHandler := handler.NewAuthHandler(authService)
    athleteHandler := handler.NewAthleteHandler(athleteService)
    competitionHandler := handler.NewCompetitionHandler(competitionService)
    // ... другие обработчики
    
    // Создание роутера
    router := mux.NewRouter()
    
    // Публичные маршруты
    router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })
    
    // Аутентификация
    authHandler.RegisterRoutes(router)
    
    // Защищенные маршруты (требуется аутентификация)
    apiRouter := router.PathPrefix("/api").Subrouter()
    apiRouter.Use(middleware.AuthMiddleware(authService))
    
    athleteHandler.RegisterRoutes(apiRouter)
    competitionHandler.RegisterRoutes(apiRouter)
    // ... другие защищенные маршруты
    
    // Настройка CORS
    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:3000"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Authorization", "Content-Type"},
        AllowCredentials: true,
    })
    
    // Настройка сервера
    srv := &http.Server{
        Addr:         ":" + cfg.Port,
        Handler:      c.Handler(router),
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }
    
    // Graceful shutdown
    go func() {
        log.Printf("Server starting on port %s", cfg.Port)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal("Server error:", err)
        }
    }()
    
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt)
    <-quit
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Server shutdown error:", err)
    }
    log.Println("Server stopped")
}