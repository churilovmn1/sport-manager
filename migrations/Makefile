# ============================================
# СПОРТИВНЫЙ МЕНЕДЖЕР - МАКЕФАЙЛ
# ============================================

# Объявляем фиктивные цели (не являются файлами)
.PHONY: help run build test migrate clean docker-up docker-down lint format deps

# Цель по умолчанию (выполняется при вызове просто `make`)
.DEFAULT_GOAL := help

# ============================================
# ПОМОЩЬ
# ============================================

help: ## Показать это сообщение помощи
	@echo "Доступные команды:"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# ============================================
# РАЗРАБОТКА
# ============================================

run: deps ## Запустить приложение в режиме разработки
	@echo "Запуск приложения..."
	@go run cmd/web/main.go

dev: deps ## Запуск с горячей перезагрузкой (нужен air)
	@echo "Запуск с hot reload..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Установите air: go install github.com/cosmtrek/air@latest"; \
		exit 1; \
	fi

# ============================================
# СБОРКА
# ============================================

build: deps ## Собрать бинарный файл
	@echo "Сборка приложения..."
	@mkdir -p bin
	@go build -ldflags="-s -w" -o bin/sport-manager cmd/web/main.go
	@echo "Собран файл: bin/sport-manager"

build-linux: deps ## Собрать для Linux
	@echo "Сборка для Linux..."
	@mkdir -p bin
	@GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/sport-manager-linux cmd/web/main.go

build-windows: deps ## Собрать для Windows
	@echo "Сборка для Windows..."
	@mkdir -p bin
	@GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o bin/sport-manager.exe cmd/web/main.go

# ============================================
# ТЕСТИРОВАНИЕ
# ============================================

test: ## Запустить unit-тесты
	@echo "Запуск тестов..."
	@go test -v ./... -cover

test-cover: ## Запустить тесты с покрытием
	@echo "Запуск тестов с покрытием..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Отчет покрытия: coverage.html"

# ============================================
# БАЗА ДАННЫХ
# ============================================

migrate: ## Применить миграции
	@echo "Применение миграций..."
	@psql -h localhost -U admin -d sport_manager -f migrations/001_init_schema.sql

migrate-create: ## Создать новую миграцию
	@echo "Создание миграции..."
	@read -p "Введите название миграции: " name; \
	timestamp=$$(date +%Y%m%d%H%M%S); \
	file="migrations/$${timestamp}_$${name}.sql"; \
	echo "-- Миграция: $${name}" > $$file; \
	echo "-- Автор: $$(whoami)" >> $$file; \
	echo "-- Дата: $$(date)" >> $$file; \
	echo "" >> $$file; \
	echo "BEGIN;" >> $$file; \
	echo "" >> $$file; \
	echo "-- Впишите SQL команды здесь" >> $$file; \
	echo "" >> $$file; \
	echo "COMMIT;" >> $$file; \
	echo "Создан файл: $$file"

# ============================================
# DOCKER
# ============================================

docker-up: ## Запустить контейнеры
	@echo "Запуск Docker контейнеров..."
	@docker-compose up -d --build
	@echo "Приложение доступно по адресу: http://localhost:8080"

docker-down: ## Остановить контейнеры
	@echo "Остановка Docker контейнеров..."
	@docker-compose down

docker-logs: ## Показать логи контейнеров
	@docker-compose logs -f

docker-clean: ## Полная очистка Docker
	@echo "Полная очистка Docker..."
	@docker-compose down -v
	@docker system prune -f

# ============================================
# КАЧЕСТВО КОДА
# ============================================

lint: ## Проверка стиля кода
	@echo "Проверка кода..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "Установите golangci-lint: https://golangci-lint.run/"; \
	fi

format: ## Форматирование кода
	@echo "Форматирование кода..."
	@go fmt ./...
	@if command -v goimports > /dev/null; then \
		goimports -w .; \
	fi

# ============================================
# ЗАВИСИМОСТИ
# ============================================

deps: ## Установить/обновить зависимости
	@echo "Проверка зависимостей..."
	@go mod tidy
	@go mod download

deps-upgrade: ## Обновить все зависимости
	@echo "Обновление зависимостей..."
	@go get -u ./...
	@go mod tidy

# ============================================
# ОЧИСТКА
# ============================================

clean: ## Очистить проект
	@echo "Очистка проекта..."
	@rm -rf bin/
	@rm -rf coverage.out coverage.html
	@go clean
	@echo "Очистка завершена"

# ============================================
# ДОКУМЕНТАЦИЯ
# ============================================

swagger: ## Генерация Swagger документации
	@echo "Генерация Swagger документации..."
	@if command -v swag > /dev/null; then \
		swag init -g cmd/web/main.go -o api/docs; \
	else \
		echo "Установите swag: go install github.com/swaggo/swag/cmd/swag@latest"; \
	fi

docs: swagger ## Запустить документацию
	@echo "Запуск документации..."
	@if [ -f "bin/sport-manager" ]; then \
		./bin/sport-manager & \
		sleep 3; \
		open http://localhost:8080/swagger/index.html; \
	else \
		echo "Сначала соберите проект: make build"; \
	fi