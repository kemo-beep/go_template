# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=server
BINARY_UNIX=$(BINARY_NAME)_unix

# Database parameters
DB_HOST=localhost
DB_PORT=5433
DB_USER=app
DB_PASS=secret
DB_NAME=myapp
DB_SSLMODE=disable
DB_DSN=postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)

.PHONY: all build clean test coverage deps run dev docker-build docker-run migrate-up migrate-down migrate-create setup lint format

# Default target
all: test build

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/server

# Build for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v ./cmd/server

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Run the application
run: build
	./$(BINARY_NAME)

# Run in development mode with hot reload
dev:
	air

# Docker build
docker-build:
	docker build -t go-mobile-backend .

# Docker run
docker-run:
	docker-compose up --build

# Docker run production
docker-prod:
	docker-compose -f deploy/docker-compose.yml up --build

# Run migrations up
migrate-up:
	goose -dir ./internal/db/migrations postgres "$(DB_DSN)" up

# Run migrations down
migrate-down:
	goose -dir ./internal/db/migrations postgres "$(DB_DSN)" down

# Create new migration
migrate-create:
	@read -p "Enter migration name: " name; \
	goose -dir ./internal/db/migrations create $$name sql

# Setup development environment
setup: deps
	@echo "Setting up development environment..."
	@if [ ! -f .env ]; then cp env.example .env; fi
	@echo "Please update .env file with your configuration"
	@echo "Run 'make docker-run' to start the application"

# Lint code
lint:
	golangci-lint run

# Format code
format:
	$(GOCMD) fmt ./...
	goimports -w .

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/cosmtrek/air@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger docs
swagger:
	swag init -g cmd/server/main.go -o docs/api

# Database operations
db-create:
	createdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) $(DB_NAME)

db-drop:
	dropdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) $(DB_NAME)

# Help
help:
	@echo "Available commands:"
	@echo "  build          - Build the application"
	@echo "  build-linux    - Build for Linux"
	@echo "  clean          - Clean build artifacts"
	@echo "  test           - Run tests"
	@echo "  coverage       - Run tests with coverage"
	@echo "  deps           - Download dependencies"
	@echo "  run            - Run the application"
	@echo "  dev            - Run in development mode"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run with Docker Compose"
	@echo "  docker-prod    - Run production with Docker Compose"
	@echo "  migrate-up     - Run database migrations"
	@echo "  migrate-down   - Rollback database migrations"
	@echo "  migrate-create - Create new migration"
	@echo "  setup          - Setup development environment"
	@echo "  lint           - Lint code"
	@echo "  format         - Format code"
	@echo "  install-tools  - Install development tools"
	@echo "  swagger        - Generate Swagger docs"
	@echo "  db-create      - Create database"
	@echo "  db-drop        - Drop database"
