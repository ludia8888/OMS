.PHONY: help build run test clean lint docker migrate generate

# Variables
APP_NAME := oms-service
GO := go
GOFLAGS := -v
DOCKER_IMAGE := openfoundry/oms-service
DOCKER_TAG := latest

# Default target
help:
	@echo "Available targets:"
	@echo "  build      - Build the application"
	@echo "  run        - Run the application"
	@echo "  test       - Run tests"
	@echo "  lint       - Run linters"
	@echo "  clean      - Clean build artifacts"
	@echo "  docker     - Build Docker image"
	@echo "  migrate    - Run database migrations"
	@echo "  generate   - Generate code (GraphQL, protobuf)"

# Build the application
build:
	$(GO) build $(GOFLAGS) -o $(APP_NAME) cmd/server/main.go

# Run the application
run:
	$(GO) run cmd/server/main.go

# Run tests
test:
	$(GO) test -v -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

# Run unit tests only
test-unit:
	$(GO) test -v -race -short ./...

# Run integration tests
test-integration:
	$(GO) test -v -race -run Integration ./...

# Run linters
lint:
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	golangci-lint run ./...

# Clean build artifacts
clean:
	rm -f $(APP_NAME)
	rm -f coverage.out coverage.html
	rm -rf dist/

# Build Docker image
docker:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Run database migrations
migrate-up:
	migrate -path internal/infrastructure/database/migrations -database "postgresql://localhost:5432/oms?sslmode=disable" up

migrate-down:
	migrate -path internal/infrastructure/database/migrations -database "postgresql://localhost:5432/oms?sslmode=disable" down

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir internal/infrastructure/database/migrations -seq $$name

# Generate code
generate: generate-graphql generate-proto

# Generate GraphQL code
generate-graphql:
	@if ! command -v gqlgen &> /dev/null; then \
		echo "Installing gqlgen..."; \
		go install github.com/99designs/gqlgen@latest; \
	fi
	cd internal/interfaces/graphql && gqlgen generate

# Generate protobuf code
generate-proto:
	@if ! command -v protoc &> /dev/null; then \
		echo "Please install protoc first"; \
		exit 1; \
	fi
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		pkg/proto/*.proto

# Development setup
dev-setup:
	@echo "Setting up development environment..."
	@if ! command -v migrate &> /dev/null; then \
		echo "Installing golang-migrate..."; \
		go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
	fi
	@if ! command -v gqlgen &> /dev/null; then \
		echo "Installing gqlgen..."; \
		go install github.com/99designs/gqlgen@latest; \
	fi
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@echo "Development environment setup complete!"

# Start local dependencies
deps-up:
	docker-compose -f deployments/docker-compose.dev.yml up -d

# Stop local dependencies
deps-down:
	docker-compose -f deployments/docker-compose.dev.yml down

# View logs
logs:
	docker-compose -f deployments/docker-compose.dev.yml logs -f

# Format code
fmt:
	$(GO) fmt ./...
	gofmt -s -w .

# Run security check
security:
	@if ! command -v gosec &> /dev/null; then \
		echo "Installing gosec..."; \
		go install github.com/securego/gosec/v2/cmd/gosec@latest; \
	fi
	gosec ./...