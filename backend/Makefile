# CI/CD Status Notifier Bot Makefile

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=cicd-notifier-bot
BINARY_UNIX=$(BINARY_NAME)_unix

# Docker parameters
DOCKER_IMAGE=cicd-notifier-bot
DOCKER_TAG=latest

# Database parameters
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=jQlwjVKoQw
DB_NAME=cicd_notifier
DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

.PHONY: all build clean test coverage deps dev docker-build docker-run docker-stop help

## Build the application
build:
	@echo "Building application..."
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/main.go

## Build for linux
build-linux:
	@echo "Building for Linux..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v ./cmd/main.go

## Clean build files
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

## Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v -race -timeout 30s ./...

## Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

## Run tests and display coverage in terminal
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -func=coverage.out

## Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

## Run the application in development mode
dev:
	@echo "Starting application in development mode..."
	$(GOCMD) run cmd/main.go

## Watch and restart on file changes (requires air)
watch:
	@echo "Starting application with hot reload..."
	@if [ ! -f .air.toml ]; then air init; fi
	air

## Install air for hot reloading
install-air:
	@echo "Installing air..."
	$(GOGET) -u github.com/cosmtrek/air

## Install migrate tool (cross-platform)
install-migrate:
	@echo "Installing migrate tool..."
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

## Run database migrations up (with better error handling)
migrate-up:
	@echo "Running database migrations up..."
	@MIGRATE_BIN=$$(go env GOPATH)/bin/migrate; \
	if [ ! -f "$$MIGRATE_BIN" ]; then \
		echo "migrate tool not found. Installing..."; \
		$(MAKE) install-migrate; \
	fi; \
	$$MIGRATE_BIN -path scripts/migrations -database "$(DB_URL)" up

## Run database migrations down (with better error handling)
migrate-down:
	@echo "Running database migrations down..."
	@MIGRATE_BIN=$$(go env GOPATH)/bin/migrate; \
	if [ ! -f "$$MIGRATE_BIN" ]; then \
		echo "migrate tool not found. Installing..."; \
		$(MAKE) install-migrate; \
	fi; \
	$$MIGRATE_BIN -path scripts/migrations -database "$(DB_URL)" down

## Force migrate down all migrations
migrate-force-down:
	@echo "Force running all migrations down..."
	@MIGRATE_BIN=$$(go env GOPATH)/bin/migrate; \
	if [ ! -f "$$MIGRATE_BIN" ]; then \
		echo "migrate tool not found. Installing..."; \
		$(MAKE) install-migrate; \
	fi; \
	$$MIGRATE_BIN -path scripts/migrations -database "$(DB_URL)" down -all

## Reset database (drop and recreate)
db-reset:
	@echo "Resetting database..."
	@psql -U $(DB_USER) -h $(DB_HOST) -p $(DB_PORT) -c "DROP DATABASE IF EXISTS $(DB_NAME);" || true
	@psql -U $(DB_USER) -h $(DB_HOST) -p $(DB_PORT) -c "CREATE DATABASE $(DB_NAME);"
	@echo "Database $(DB_NAME) has been reset"

## Setup database (create database if not exists)
db-setup:
	@echo "Setting up database..."
	@psql -U $(DB_USER) -h $(DB_HOST) -p $(DB_PORT) -c "CREATE DATABASE $(DB_NAME);" || echo "Database $(DB_NAME) already exists"

## Complete database reset and migrate
db-fresh: db-reset migrate-up
	@echo "Database fresh migration completed!"

## Check migration status
migrate-status:
	@echo "Checking migration status..."
	@MIGRATE_BIN=$$(go env GOPATH)/bin/migrate; \
	if [ ! -f "$$MIGRATE_BIN" ]; then \
		echo "migrate tool not found. Installing..."; \
		$(MAKE) install-migrate; \
	fi; \
	$$MIGRATE_BIN -path scripts/migrations -database "$(DB_URL)" version

## Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

## Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run -d --name $(DOCKER_IMAGE) \
		-p 8080:8080 \
		-e DB_HOST=host.docker.internal \
		$(DOCKER_IMAGE):$(DOCKER_TAG)

## Stop Docker container
docker-stop:
	@echo "Stopping Docker container..."
	docker stop $(DOCKER_IMAGE) || true
	docker rm $(DOCKER_IMAGE) || true

## Run with Docker Compose
compose-up:
	@echo "Starting application with Docker Compose..."
	docker-compose up -d

## Stop Docker Compose
compose-down:
	@echo "Stopping Docker Compose..."
	docker-compose down

## Development environment setup
setup-dev: deps install-air install-migrate db-setup
	@echo "Development environment setup complete!"
	@echo "Run 'make dev' to start the application"

## Setup environment with database fresh migration
setup-fresh: deps install-air install-migrate db-fresh
	@echo "Fresh development environment setup complete!"
	@echo "Run 'make dev' to start the application"

## Lint the code
lint:
	@echo "Running linter..."
	golangci-lint run

## Install linter
install-lint:
	@echo "Installing golangci-lint..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.54.2

## Format code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

## Vet code
vet:
	@echo "Vetting code..."
	$(GOCMD) vet ./...

## Run security check
sec:
	@echo "Running security check..."
	gosec ./...

## Install security scanner
install-sec:
	@echo "Installing gosec..."
	$(GOGET) github.com/securecodewarrior/gosec/v2/cmd/gosec

## Generate mock files
generate:
	@echo "Generating mock files..."
	$(GOCMD) generate ./...

## All quality checks
quality: fmt vet lint test-coverage sec

## Show help
help:
	@echo "Available commands:"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "Examples:"
	@echo "  make setup-dev     # Setup development environment"
	@echo "  make dev           # Run application in development mode"
	@echo "  make test          # Run tests"
	@echo "  make build         # Build the application"

## Default target
all: deps fmt vet lint test build
