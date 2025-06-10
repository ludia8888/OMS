.PHONY: help build run test lint clean migrate docker-build docker-run

# Variables
APP_NAME = oms-service
VERSION := $(shell git describe --tags --always --dirty)
BUILD_DIR = ./bin
MAIN_PATH = ./cmd/server
DOCKER_IMAGE = openfoundry/oms-service

# Default target
help:
	@echo "Available commands:"
	@echo "  make build       - Build the application"
	@echo "  make run         - Run the application"
	@echo "  make test        - Run tests"
	@echo "  make lint        - Run linter"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make migrate-up  - Run database migrations"
	@echo "  make migrate-down - Rollback database migrations"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run  - Run Docker container"

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@go build -ldflags "-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

# Run the application
run:
	@echo "Running $(APP_NAME)..."
	@go run $(MAIN_PATH)/main.go

# Run with environment variables
run-dev:
	@echo "Running $(APP_NAME) in development mode..."
	@DB_PASSWORD=testpass JWT_SECRET=testsecret go run $(MAIN_PATH)/main.go

# Run tests
test:
	@echo "Running tests..."
	@go test -v -cover -race ./...

# Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linter
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: brew install golangci-lint"; \
		exit 1; \
	fi

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Code formatted"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

# Generate mocks
generate:
	@echo "Generating mocks..."
	@go generate ./...

# Database migrations
migrate-up:
	@echo "Running migrations..."
	@migrate -path ./internal/infrastructure/database/migrations -database "postgresql://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=disable" up

migrate-down:
	@echo "Rolling back migrations..."
	@migrate -path ./internal/infrastructure/database/migrations -database "postgresql://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=disable" down

migrate-create:
	@echo "Creating new migration..."
	@migrate create -ext sql -dir ./internal/infrastructure/database/migrations -seq $(name)

# Docker commands
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE):$(VERSION) -t $(DOCKER_IMAGE):latest .

docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE):latest

# Development helpers
dev-deps:
	@echo "Installing development dependencies..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install github.com/99designs/gqlgen@latest
	@echo "Development dependencies installed"

# GraphQL generation
gqlgen:
	@echo "Generating GraphQL code..."
	@go run github.com/99designs/gqlgen generate