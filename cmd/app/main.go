package main

import (
	// ... (старые импорты)
	"net/http"
	"sport-manager/internal/auth" // Добавляем импорт пакета auth

	"github.com/gorilla/mux"
	"github.com/yourusername/sport-manager/internal/handler"
	"github.com/yourusername/sport-manager/internal/repository"
	"github.com/yourusername/sport-manager/internal/service"
	// ...
)

func main() {
	// ... (загрузка конфига, инициализация БД, инициализация Auth Repo/Service/Handler)

	// ...
	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo, cfg)
	authHandler := handler.NewAuthHandler(authService)

	// 4. Настройка роутера (mux)
	router := mux.NewRouter()

	// Публичные маршруты (Auth)
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Service is up and running!"))
	}).Methods("GET")
	router.HandleFunc("/api/v1/auth/login", authHandler.Login).Methods("POST")

	// --- 5. Защищенные маршруты (CRUD) ---
	// Создаем под-роутер, к которому применяем мидлвар
	apiRouter := router.PathPrefix("/api/v1").Subrouter()
	apiRouter.Use(auth.AuthMiddleware(cfg)) // Применяем мидлвар ко всем маршрутам apiRouter

	// Теперь, все что мы добавим в apiRouter, будет требовать JWT.
	// Пример (заглушка):
	apiRouter.HandleFunc("/athletes", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, authorized admin!"))
	}).Methods("GET")

	// ... (Запуск сервера)
	// ...
}

// initDB (остается без изменений)
// ...
