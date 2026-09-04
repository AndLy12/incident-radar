# Короткие команды для установки, запуска и проверки проекта.
GO ?= go
ARGS ?=
PORT ?= 8000
BINARY ?= bin/incident-radar

.PHONY: install build run serve format check test test-race

install: ## Скачать зависимости приложения и локальный линтер
	$(GO) mod download

build: ## Собрать исполняемый файл
	$(GO) build -trimpath -o $(BINARY) ./cmd/incident-radar

run: ## Запустить приложение, дополнительные параметры передаются через ARGS
	$(GO) run ./cmd/incident-radar $(ARGS)

serve: ## Запустить локальный HTTP-сервер с тестовыми NGINX-логами
	$(GO) run ./cmd/log-server --port $(PORT)

format: ## Применить gofumpt и goimports
	$(GO) tool golangci-lint fmt

check: ## Проверить go.mod, форматирование и код линтерами
	$(GO) mod tidy -diff
	$(GO) tool golangci-lint config verify
	$(GO) tool golangci-lint fmt --diff
	$(GO) tool golangci-lint run

test: ## Запустить тесты и сформировать HTML-отчёт покрытия
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -func=coverage.out
	$(GO) tool cover -html=coverage.out -o coverage.html

test-race: ## Проверить тесты с race detector (требуется CGO)
	$(GO) test -race ./...
