# TaskFlow — Makefile
# ==========================================

BINARY_NAME=taskflow
MAIN_PATH=./cmd/taskflow
DB_PATH=taskflow.db

# Default target
.PHONY: all
all: docker-build

# Run the application locally (requires GCC)
.PHONY: run
run:
	go run $(MAIN_PATH)

# Download dependencies
.PHONY: deps
deps:
	@echo "📦 Descargando dependencias..."
	go mod tidy
	go mod download
	@echo "✓ Dependencias listas"

# Run tests
.PHONY: test
test:
	go test ./...

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Clean
.PHONY: clean
clean:
	@echo "🧹 Limpiando..."
	-del /f $(DB_PATH) 2>nul
	@echo "✓ Limpieza completada"

# Docker
.PHONY: docker-build
docker-build:
	@echo "🐳 Construyendo imagen Docker..."
	docker build -t taskflow .
	@echo "✓ Imagen construida: taskflow"

.PHONY: docker-run
docker-run:
	@echo "🐳 Iniciando contenedor..."
	docker compose up -d
	@echo "TaskFlow disponible en http://localhost:3000"

.PHONY: docker-up
docker-up:
	@echo "🐳 Construyendo e iniciando..."
	docker compose up -d --build
	@echo "TaskFlow disponible en http://localhost:3000"

.PHONY: docker-stop
docker-stop:
	docker compose down

.PHONY: docker-logs
docker-logs:
	docker compose logs -f

# Help
.PHONY: help
help:
	@echo "TaskFlow - Comandos disponibles:"
	@echo ""
	@echo "  make run           Ejecuta localmente (requiere GCC)"
	@echo "  make deps          Descarga dependencias"
	@echo "  make test          Ejecuta tests"
	@echo "  make fmt           Formatea el codigo"
	@echo "  make clean         Limpia archivos generados"
	@echo ""
	@echo "  make docker-build  Construye imagen Docker"
	@echo "  make docker-up     Construye e inicia (todo en uno)"
	@echo "  make docker-run    Inicia con Docker Compose"
	@echo "  make docker-stop   Detiene el contenedor"
	@echo "  make docker-logs   Ver logs del contenedor"
	@echo ""
	@echo "  make help          Muestra esta ayuda"
