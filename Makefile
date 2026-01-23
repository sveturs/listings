.PHONY: help build run test clean docker-build docker-up docker-down migrate-up migrate-down proto lint format tidy deps dump2mig-dump dump2mig

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
DATABASE_URL := postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable

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

test: ## Run unit tests (excludes integration)
	@echo "$(GREEN)Running unit tests...$(NC)"
	@$(GO) test -v -race -cover -short ./...

test-all: ## Run all tests (unit + integration)
	@echo "$(GREEN)Running all tests (unit + integration)...$(NC)"
	@$(GO) test -v -race -cover ./...

test-coverage: ## Run tests with coverage report
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	@$(GO) test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	@$(GO) tool cover -html=coverage.txt -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"
	@echo "$(YELLOW)Coverage summary:$(NC)"
	@$(GO) tool cover -func=coverage.txt | tail -1

test-coverage-check: test-coverage ## Check if coverage meets threshold (70%)
	@echo "$(GREEN)Checking coverage threshold...$(NC)"
	@COVERAGE=$$($(GO) tool cover -func=coverage.txt | grep total | awk '{print $$3}' | sed 's/%//'); \
	if [ "$$(echo "$$COVERAGE < 70" | bc -l)" -eq 1 ]; then \
		echo "$(RED)Coverage $$COVERAGE% is below 70% threshold$(NC)"; \
		exit 1; \
	else \
		echo "$(GREEN)Coverage $$COVERAGE% meets threshold$(NC)"; \
	fi

test-unit: ## Run only unit tests
	@echo "$(GREEN)Running unit tests...$(NC)"
	@$(GO) test -v -race -cover -short ./internal/...

test-integration: ## Run integration tests
	@echo "$(GREEN)Running integration tests...$(NC)"
	@$(GO) test -v -tags=integration ./tests/integration/...

test-e2e: ## Run end-to-end tests
	@echo "$(GREEN)Running E2E tests...$(NC)"
	@$(GO) test -v ./tests/e2e/...

bench: ## Run benchmarks
	@echo "$(GREEN)Running benchmarks...$(NC)"
	@$(GO) test -bench=. -benchmem ./tests/performance/...

bench-cpu: ## Run benchmarks with CPU profiling
	@echo "$(GREEN)Running benchmarks with CPU profiling...$(NC)"
	@$(GO) test -bench=. -benchmem -cpuprofile=cpu.prof ./tests/performance/...
	@echo "$(YELLOW)View profile: go tool pprof cpu.prof$(NC)"

bench-mem: ## Run benchmarks with memory profiling
	@echo "$(GREEN)Running benchmarks with memory profiling...$(NC)"
	@$(GO) test -bench=. -benchmem -memprofile=mem.prof ./tests/performance/...
	@echo "$(YELLOW)View profile: go tool pprof mem.prof$(NC)"

test-verbose: ## Run tests with verbose output
	@echo "$(GREEN)Running tests (verbose)...$(NC)"
	@$(GO) test -v -race -cover ./... 2>&1 | tee test-output.log

test-watch: ## Run tests in watch mode (requires entr)
	@echo "$(GREEN)Running tests in watch mode...$(NC)"
	@which entr > /dev/null || (echo "$(RED)entr not installed$(NC)" && exit 1)
	@find . -name '*.go' | entr -c make test

generate-mocks: ## Generate mocks using mockgen
	@echo "$(GREEN)Generating mocks...$(NC)"
	@which mockgen > /dev/null || go install go.uber.org/mock/mockgen@latest
	@echo "$(YELLOW)Note: Mock generation will be implemented in Sprint 4.2$(NC)"

## Code quality commands

lint: ## Run linter (golangci-lint)
	@echo "$(GREEN)Running linter...$(NC)"
	@which golangci-lint > /dev/null || (echo "$(RED)golangci-lint not installed. Run: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin latest$(NC)" && exit 1)
	@golangci-lint run --timeout=5m

lint-fix: ## Run linter with auto-fix
	@echo "$(GREEN)Running linter with auto-fix...$(NC)"
	@which golangci-lint > /dev/null || (echo "$(RED)golangci-lint not installed. Run: make lint$(NC)" && exit 1)
	@golangci-lint run --fix --timeout=5m

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

## Docker dependency commands (PostgreSQL + Redis only)

deps-up: ## Start PostgreSQL and Redis
	@echo "$(GREEN)Starting PostgreSQL and Redis...$(NC)"
	@$(DOCKER_COMPOSE) up -d
	@echo "$(YELLOW)Waiting for PostgreSQL to be ready...$(NC)"
	@timeout=30; counter=0; \
	until $(DOCKER_COMPOSE) exec -T postgres pg_isready -U listings_user -d listings_dev_db > /dev/null 2>&1; do \
		sleep 1; \
		counter=$$((counter + 1)); \
		if [ $$counter -ge $$timeout ]; then \
			echo "$(RED)Timeout waiting for PostgreSQL$(NC)"; \
			exit 1; \
		fi; \
	done
	@echo "$(GREEN)PostgreSQL is ready$(NC)"

deps-down: ## Stop PostgreSQL and Redis (keep data)
	@echo "$(YELLOW)Stopping PostgreSQL and Redis (keeping data)...$(NC)"
	@$(DOCKER_COMPOSE) down
	@echo "$(GREEN)Services stopped$(NC)"

deps-clean: ## Remove PostgreSQL and Redis with all data
	@echo "$(RED)Stopping and removing PostgreSQL and Redis with all data...$(NC)"
	@$(DOCKER_COMPOSE) down -v --remove-orphans
	@echo "$(GREEN)All containers and volumes removed$(NC)"

deps-reset: deps-clean deps-up ## Clean restart: remove all data + start fresh
	@echo ""
	@echo "=========================================="
	@echo "$(GREEN)✅ Dependencies reset complete!$(NC)"
	@echo "=========================================="
	@echo "PostgreSQL and Redis are running with fresh data."
	@echo "Run 'make migrate-up' to apply migrations."

## Service start/stop commands

start: stop ## Start the service in background with logs to file
	@mkdir -p logs
	@echo "$(GREEN)Building $(APP_NAME)...$(NC)"
	@$(GO) build $(GOFLAGS) -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(APP_NAME) ./cmd/server
	@echo "$(GREEN)Starting $(APP_NAME) in background...$(NC)"
	@nohup ./$(BUILD_DIR)/$(APP_NAME) > logs/$(APP_NAME).log 2>&1 &
	@sleep 2
	@if lsof -ti:8086 > /dev/null 2>&1; then \
		echo "$(GREEN)$(APP_NAME) started successfully$(NC)"; \
		echo "  HTTP: http://localhost:8086"; \
		echo "  gRPC: localhost:50053"; \
		echo "  Logs: tail -f logs/$(APP_NAME).log"; \
	else \
		echo "$(RED)Failed to start $(APP_NAME). Check logs/$(APP_NAME).log$(NC)"; \
		exit 1; \
	fi

stop: ## Stop the service
	@if lsof -ti:8086 > /dev/null 2>&1; then \
		echo "$(YELLOW)Stopping $(APP_NAME)...$(NC)"; \
		lsof -ti:8086 | xargs kill 2>/dev/null || true; \
		sleep 1; \
		if lsof -ti:8086 > /dev/null 2>&1; then \
			echo "$(YELLOW)Force killing $(APP_NAME)...$(NC)"; \
			lsof -ti:8086 | xargs kill -9 2>/dev/null || true; \
		fi; \
		echo "$(GREEN)$(APP_NAME) stopped$(NC)"; \
	fi

## Docker commands (legacy aliases)

docker-build: ## Build Docker image
	@echo "$(GREEN)Building Docker image...$(NC)"
	@docker build -t $(DOCKER_IMAGE) .
	@echo "$(GREEN)Docker image built: $(DOCKER_IMAGE)$(NC)"

docker-up: deps-up ## Alias for deps-up

docker-down: deps-down ## Alias for deps-down

docker-restart: deps-down deps-up ## Restart Docker Compose services

docker-logs: ## View Docker Compose logs
	@$(DOCKER_COMPOSE) logs -f

docker-ps: ## Show running containers
	@$(DOCKER_COMPOSE) ps

docker-clean: deps-clean ## Alias for deps-clean

## Database migration commands

# Build the migrator
build-migrator: ## Build migrator binary
	@echo "$(GREEN)Building migrator...$(NC)"
	$(GO) build -v -o bin/migrator cmd/migrator/main.go
	@echo "$(GREEN)Migrator built: bin/migrator$(NC)"

# Build all binaries
build-all: build build-migrator ## Build all binaries

migrate-up: ## Run database migrations (schema only)
	@echo "$(GREEN)Running database migrations...$(NC)"
	$(GO) run cmd/migrator/main.go -command up

migrate-up-fixtures: ## Run migrations + fixtures
	@echo "$(GREEN)Running database migrations with fixtures...$(NC)"
	$(GO) run cmd/migrator/main.go -command up -with-fixtures

migrate-only-fixtures: ## Run only fixtures (no schema)
	@echo "$(GREEN)Running only fixtures...$(NC)"
	$(GO) run cmd/migrator/main.go -command up -only-fixtures

migrate-down: ## Rollback last migration
	@echo "$(YELLOW)Rolling back last migration...$(NC)"
	$(GO) run cmd/migrator/main.go -command down

migrate-down-all: ## Rollback all migrations
	@echo "$(RED)Rolling back all migrations...$(NC)"
	$(GO) run cmd/migrator/main.go -command down-all

migrate-status: ## Show migration status
	@echo "$(GREEN)Checking migration status...$(NC)"
	$(GO) run cmd/migrator/main.go -command status

migrate-force: ## Force set migration version (usage: make migrate-force VERSION=5)
	@if [ -z "$(VERSION)" ]; then echo "$(RED)VERSION is required. Usage: make migrate-force VERSION=5$(NC)"; exit 1; fi
	@echo "$(YELLOW)Forcing migration version to $(VERSION)...$(NC)"
	$(GO) run cmd/migrator/main.go -command force -version $(VERSION)

migrate-create: ## Create a new migration file (usage: make migrate-create NAME=add_users_table)
	@if [ -z "$(NAME)" ]; then echo "$(RED)NAME is required. Usage: make migrate-create NAME=add_users_table$(NC)"; exit 1; fi
	@echo "$(GREEN)Creating migration: $(NAME)...$(NC)"
	$(GO) run cmd/migrator/main.go -command create -name $(NAME)

migrate-reset: migrate-down-all migrate-up ## Reset database (down all + up all)

migrate-version: migrate-status ## Alias for migrate-status (backward compatibility)

## Protobuf commands

proto: ## Generate Go code from protobuf files
	@echo "$(GREEN)Generating protobuf code...$(NC)"
	@which protoc > /dev/null || (echo "$(RED)protoc not installed. Visit: https://grpc.io/docs/protoc-installation/$(NC)" && exit 1)
	@which protoc-gen-go > /dev/null || go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@which protoc-gen-go-grpc > /dev/null || go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	PATH=/home/dim/go/bin:$$PATH protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/*.proto
	PATH=/home/dim/go/bin:$$PATH protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/proto/categories/v1/*.proto
	PATH=/home/dim/go/bin:$$PATH protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/proto/categories/v2/*.proto
	PATH=/home/dim/go/bin:$$PATH protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/proto/attributes/v1/*.proto
	PATH=/home/dim/go/bin:$$PATH protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/proto/search/v1/*.proto
	PATH=/home/dim/go/bin:$$PATH protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/proto/chat/v1/*.proto
	PATH=/home/dim/go/bin:$$PATH protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/proto/variants/v1/*.proto
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

## Load Testing commands

load-test: ## Run all load tests (HTTP + gRPC)
	@echo "$(GREEN)Running load tests...$(NC)"
	@cd load-tests && ./run-all-tests.sh

load-test-http: ## Run only HTTP load tests
	@echo "$(GREEN)Running HTTP load tests...$(NC)"
	@cd load-tests && ./run-all-tests.sh --http-only

load-test-grpc: ## Run only gRPC load tests
	@echo "$(GREEN)Running gRPC load tests...$(NC)"
	@cd load-tests && ./run-all-tests.sh --grpc-only

load-test-analyze: ## Analyze latest load test results
	@echo "$(GREEN)Analyzing load test results...$(NC)"
	@cd load-tests && ./analyze-results.sh

load-test-setup: ## Setup load testing environment with Docker
	@echo "$(GREEN)Setting up load testing environment...$(NC)"
	@cd load-tests && docker-compose -f docker-compose.load-test.yml up -d
	@echo "$(YELLOW)Waiting for services to be ready...$(NC)"
	@sleep 10
	@echo "$(GREEN)Load testing environment ready!$(NC)"
	@echo "$(YELLOW)Grafana: http://localhost:3000 (admin/admin)$(NC)"
	@echo "$(YELLOW)Prometheus: http://localhost:9090$(NC)"

load-test-teardown: ## Stop load testing environment
	@echo "$(YELLOW)Stopping load testing environment...$(NC)"
	@cd load-tests && docker-compose -f docker-compose.load-test.yml down
	@echo "$(GREEN)Load testing environment stopped$(NC)"

load-test-clean: ## Clean load test results
	@echo "$(YELLOW)Cleaning load test results...$(NC)"
	@rm -rf load-tests/results/*.json load-tests/results/*.log load-tests/results/*.txt
	@echo "$(GREEN)Load test results cleaned$(NC)"

## Dump2Mig commands (database schema migration generation)

dump2mig-dump: ## Create PostgreSQL dump for migration generation
	@echo "$(GREEN)Creating database dump for dump2mig...$(NC)"
	@echo "$(YELLOW)Connecting to: localhost:35434/listings_dev_db$(NC)"
	@echo "$(YELLOW)Note: Using schema-only dump (--schema-only) to avoid PostGIS type issues$(NC)"
	@PGPASSWORD=listings_secret pg_dump \
		-h localhost \
		-p 35434 \
		-U listings_user \
		-d listings_dev_db \
		--no-owner \
		--no-acl \
		--schema-only \
		-f dump2mig_dump.sql
	@echo "$(GREEN)Dump created: dump2mig_dump.sql ($$(wc -l < dump2mig_dump.sql) lines)$(NC)"

dump2mig: dump2mig-dump ## Full dump2mig cycle: dump -> split -> combine -> replace
	@echo ""
	@echo "$(GREEN)Splitting dump into migrations and fixtures...$(NC)"
	@rm -rf migrations_new fixtures_new migrations_big migrations_combined
	~/go/bin/dump2mig-split dump2mig_dump.sql
	@mv migrations migrations_new
	@mv fixtures fixtures_new
	@echo "$(GREEN)✅ Split complete: $$(ls migrations_new/*.sql | wc -l) migrations, $$(ls fixtures_new/*.sql | wc -l) fixtures$(NC)"
	@echo ""
	@echo "$(GREEN)Combining migration files...$(NC)"
	@mv migrations_new migrations
	~/go/bin/dump2mig-combine
	@rm -rf migrations
	@mv migrations_combined migrations_new
	@echo "$(GREEN)✅ Combine complete: $$(ls migrations_new/*.sql | wc -l) combined migrations$(NC)"
	@echo ""
	@echo "$(YELLOW)Replacing old migrations and fixtures...$(NC)"
	@rm -rf migrations fixtures migrations_big
	@mv migrations_new migrations
	@mv fixtures_new fixtures
	@rm -f dump2mig_dump.sql
	@echo ""
	@echo "=========================================="
	@echo "$(GREEN)✅ dump2mig complete!$(NC)"
	@echo "=========================================="
	@echo "Generated files:"
	@echo "  migrations/ - $$(ls migrations/*.sql 2>/dev/null | wc -l | tr -d ' ') files"
	@echo "  fixtures/   - $$(ls fixtures/*.sql 2>/dev/null | wc -l | tr -d ' ') files"
	@echo ""
