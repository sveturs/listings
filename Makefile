.PHONY: help build run test clean docker-build docker-up docker-down migrate-up migrate-down proto lint format tidy deps

# Variables
APP_NAME := listings-service
BUILD_DIR := bin
MIGRATIONS_DIR := migrations
PROTO_DIR := api/proto/listings/v1

# Go variables
GO := go
GOFLAGS := -v
LDFLAGS := -w -s

# Docker variables
DOCKER_IMAGE := $(APP_NAME):latest
DOCKER_COMPOSE := docker-compose

# Migration variables
MIGRATE := migrate
DATABASE_URL := postgres://listings_user:listings_password@localhost:35433/listings_db?sslmode=disable

# Color output
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

help: ## Show this help message
	@echo "$(GREEN)$(APP_NAME) - Available commands:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

## Build commands

build: ## Build the application binary
	@echo "$(GREEN)Building $(APP_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@$(GO) build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME) ./cmd/server
	@echo "$(GREEN)Build complete: $(BUILD_DIR)/$(APP_NAME)$(NC)"

build-all: ## Build for all platforms (linux, darwin, windows)
	@echo "$(GREEN)Building for all platforms...$(NC)"
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 ./cmd/server
	GOOS=darwin GOARCH=amd64 $(GO) build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 ./cmd/server
	GOOS=darwin GOARCH=arm64 $(GO) build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 ./cmd/server
	GOOS=windows GOARCH=amd64 $(GO) build -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe ./cmd/server
	@echo "$(GREEN)Multi-platform build complete$(NC)"

run: build ## Build and run the application locally
	@echo "$(GREEN)Running $(APP_NAME)...$(NC)"
	@./$(BUILD_DIR)/$(APP_NAME)

clean: ## Clean build artifacts
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -rf coverage.txt coverage.html
	@echo "$(GREEN)Clean complete$(NC)"

## Testing commands

test: ## Run all tests
	@echo "$(GREEN)Running tests...$(NC)"
	@$(GO) test -v -race -cover ./...

test-coverage: ## Run tests with coverage report
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	@$(GO) test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	@$(GO) tool cover -html=coverage.txt -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

test-integration: ## Run integration tests
	@echo "$(GREEN)Running integration tests...$(NC)"
	@$(GO) test -v -tags=integration ./...

bench: ## Run benchmarks
	@echo "$(GREEN)Running benchmarks...$(NC)"
	@$(GO) test -bench=. -benchmem ./...

## Code quality commands

lint: ## Run linter (golangci-lint)
	@echo "$(GREEN)Running linter...$(NC)"
	@which golangci-lint > /dev/null || (echo "$(RED)golangci-lint not installed. Run: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin$(NC)" && exit 1)
	@golangci-lint run --timeout=5m

format: ## Format code with gofmt and goimports
	@echo "$(GREEN)Formatting code...$(NC)"
	@$(GO) fmt ./...
	@which goimports > /dev/null && goimports -w . || echo "$(YELLOW)goimports not found, skipping. Install: go install golang.org/x/tools/cmd/goimports@latest$(NC)"

tidy: ## Tidy go modules
	@echo "$(GREEN)Tidying go modules...$(NC)"
	@$(GO) mod tidy
	@$(GO) mod verify

deps: ## Download dependencies
	@echo "$(GREEN)Downloading dependencies...$(NC)"
	@$(GO) mod download

## Docker commands

docker-build: ## Build Docker image
	@echo "$(GREEN)Building Docker image...$(NC)"
	@docker build -t $(DOCKER_IMAGE) .
	@echo "$(GREEN)Docker image built: $(DOCKER_IMAGE)$(NC)"

docker-up: ## Start Docker Compose services
	@echo "$(GREEN)Starting Docker Compose services...$(NC)"
	@$(DOCKER_COMPOSE) up -d
	@echo "$(GREEN)Services started. Use 'make docker-logs' to view logs$(NC)"

docker-down: ## Stop Docker Compose services
	@echo "$(YELLOW)Stopping Docker Compose services...$(NC)"
	@$(DOCKER_COMPOSE) down
	@echo "$(GREEN)Services stopped$(NC)"

docker-restart: docker-down docker-up ## Restart Docker Compose services

docker-logs: ## View Docker Compose logs
	@$(DOCKER_COMPOSE) logs -f

docker-ps: ## Show running containers
	@$(DOCKER_COMPOSE) ps

docker-clean: docker-down ## Remove Docker volumes and images
	@echo "$(YELLOW)Cleaning Docker resources...$(NC)"
	@$(DOCKER_COMPOSE) down -v --rmi local
	@echo "$(GREEN)Docker cleanup complete$(NC)"

## Database migration commands

migrate-install: ## Install golang-migrate tool
	@echo "$(GREEN)Installing golang-migrate...$(NC)"
	@which migrate > /dev/null || go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "$(GREEN)golang-migrate installed$(NC)"

migrate-up: ## Run database migrations up
	@echo "$(GREEN)Running migrations up...$(NC)"
	@$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up
	@echo "$(GREEN)Migrations applied$(NC)"

migrate-down: ## Rollback last migration
	@echo "$(YELLOW)Rolling back last migration...$(NC)"
	@$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down 1
	@echo "$(GREEN)Rollback complete$(NC)"

migrate-reset: ## Reset database (down all + up all)
	@echo "$(RED)Resetting database...$(NC)"
	@$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down -all || true
	@$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up
	@echo "$(GREEN)Database reset complete$(NC)"

migrate-version: ## Show current migration version
	@$(MIGRATE) -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" version

migrate-create: ## Create a new migration file (usage: make migrate-create NAME=add_users_table)
	@if [ -z "$(NAME)" ]; then echo "$(RED)NAME is required. Usage: make migrate-create NAME=add_users_table$(NC)"; exit 1; fi
	@echo "$(GREEN)Creating migration: $(NAME)...$(NC)"
	@$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) -seq $(NAME)
	@echo "$(GREEN)Migration created in $(MIGRATIONS_DIR)$(NC)"

## Protobuf commands

proto: ## Generate Go code from protobuf files
	@echo "$(GREEN)Generating protobuf code...$(NC)"
	@which protoc > /dev/null || (echo "$(RED)protoc not installed. Visit: https://grpc.io/docs/protoc-installation/$(NC)" && exit 1)
	@which protoc-gen-go > /dev/null || go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@which protoc-gen-go-grpc > /dev/null || go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/*.proto
	@echo "$(GREEN)Protobuf code generated$(NC)"

## Development commands

dev: docker-up migrate-up ## Setup development environment
	@echo "$(GREEN)Development environment ready!$(NC)"
	@echo "$(YELLOW)Run 'make run' to start the application$(NC)"

dev-reset: docker-down docker-clean docker-up migrate-up ## Reset development environment
	@echo "$(GREEN)Development environment reset complete$(NC)"

## CI/CD commands

ci: deps lint test build ## Run CI pipeline locally

pre-commit: format lint test ## Run pre-commit checks

## Info commands

version: ## Show application version
	@$(GO) run ./cmd/server version 2>/dev/null || echo "Build the app first: make build"

env: ## Show current environment configuration
	@echo "$(GREEN)Current environment:$(NC)"
	@echo "  APP_NAME: $(APP_NAME)"
	@echo "  BUILD_DIR: $(BUILD_DIR)"
	@echo "  GO_VERSION: $(shell $(GO) version)"
	@echo "  DATABASE_URL: $(DATABASE_URL)"
