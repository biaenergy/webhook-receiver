# Webhook Receiver Makefile

.PHONY: help build run test clean install deps

# Variables
BINARY_NAME=webhook-receiver
BUILD_DIR=build
GO_VERSION=1.21

# Colores para output
GREEN=\033[0;32m
YELLOW=\033[1;33m
RED=\033[0;31m
NC=\033[0m # No Color

help: ## Mostrar ayuda
	@echo "$(GREEN)Webhook Receiver - Comandos disponibles:$(NC)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-15s$(NC) %s\n", $$1, $$2}'

install: ## Instalar dependencias
	@echo "$(GREEN)Instalando dependencias...$(NC)"
	go mod download
	go mod tidy

deps: install ## Alias para install

build: ## Compilar el proyecto
	@echo "$(GREEN)Compilando $(BINARY_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) main.go
	@echo "$(GREEN)✅ Compilación exitosa: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

run: ## Ejecutar el servidor en modo desarrollo
	@echo "$(GREEN)Iniciando servidor webhook receiver...$(NC)"
	@echo "$(YELLOW)Variables de entorno:$(NC)"
	@echo "  PORT=$(shell echo $${PORT:-8080})"
	@echo "  WEBHOOK_SECRET_KEY=$(shell echo $${WEBHOOK_SECRET_KEY:-default-secret-key})"
	@echo "  GIN_MODE=$(shell echo $${GIN_MODE:-debug})"
	@echo ""
	go run main.go

run-prod: ## Ejecutar en modo producción
	@echo "$(GREEN)Iniciando servidor en modo producción...$(NC)"
	GIN_MODE=release go run main.go

test: ## Ejecutar tests
	@echo "$(GREEN)Ejecutando tests...$(NC)"
	go test -v ./...

test-webhook: ## Probar webhook con script
	@echo "$(GREEN)Probando webhook...$(NC)"
	@chmod +x scripts/test_webhook.sh
	./scripts/test_webhook.sh

test-client: ## Ejecutar ejemplo de cliente
	@echo "$(GREEN)Ejecutando ejemplo de cliente...$(NC)"
	go run examples/client_example.go

clean: ## Limpiar archivos generados
	@echo "$(GREEN)Limpiando archivos...$(NC)"
	rm -rf $(BUILD_DIR)
	go clean

docker-build: ## Construir imagen Docker
	@echo "$(GREEN)Construyendo imagen Docker...$(NC)"
	docker build -t webhook-receiver:latest .

docker-run: ## Ejecutar con Docker
	@echo "$(GREEN)Ejecutando con Docker...$(NC)"
	docker run -p 8080:8080 \
		-e PORT=8080 \
		-e WEBHOOK_SECRET_KEY=default-secret-key \
		-e GIN_MODE=release \
		webhook-receiver:latest

dev: ## Modo desarrollo completo (instalar + ejecutar)
	@echo "$(GREEN)Modo desarrollo - instalando dependencias y ejecutando...$(NC)"
	$(MAKE) install
	$(MAKE) run

check: ## Verificar código
	@echo "$(GREEN)Verificando código...$(NC)"
	go vet ./...
	go fmt ./...

lint: check ## Alias para check

# Comando por defecto
.DEFAULT_GOAL := help
