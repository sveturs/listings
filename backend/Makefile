.PHONY: docs docs-fmt docs-clean format format-check lint pre-commit generate-grpc test-production-readiness test-integration test-e2e test-load mock-start mock-stop test-microservice-connectivity test-traffic-router test-data-consistency test-performance-reliability test-microservice-smoke test-microservice-integration test-microservice-coverage

# Generate gRPC code from proto files (delivery service)
generate-grpc:
	@echo "Generating gRPC code from proto files..."
	@mkdir -p pkg/grpc/delivery/v1
	@PATH=$$PATH:$$HOME/go/bin protoc \
		--go_out=pkg/grpc --go_opt=paths=source_relative \
		--go-grpc_out=pkg/grpc --go-grpc_opt=paths=source_relative \
		--proto_path=proto \
		proto/delivery/v1/delivery.proto
	@echo "✅ gRPC code generated in pkg/grpc/delivery/v1/"

# Generate proto code for Listings service
generate-proto:
	@echo "Generating proto files for Listings service..."
	@mkdir -p pkg/proto/listings/v1
	@PATH=$$PATH:$$HOME/go/bin protoc \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		pkg/proto/listings/v1/listings.proto
	@echo "✅ Proto generation complete"

build_migrator:
	go build -ldflags "-X main.gitCommit=$$(git rev-parse --short HEAD) -X main.buildTime=$$(date +%Y%m%d_%H%M%S)" \
		-o bin/migrator ./cmd/migrator

build_api:
	go build -ldflags "-X main.gitCommit=$$(git rev-parse --short HEAD) -X main.buildTime=$$(date +%Y%m%d_%H%M%S)" \
		-o bin/api ./cmd/api

build_reindex:
	go build -ldflags "-X main.gitCommit=$$(git rev-parse --short HEAD) -X main.buildTime=$$(date +%Y%m%d_%H%M%S)" \
		-o bin/reindex ./cmd/reindex

build_all: build_api build_migrator build_reindex
	@echo "Build DONE"

# Swagger документация
docs:
	go tool github.com/swaggo/swag/cmd/swag init --parseDependency --parseInternal -g cmd/api/main.go -o docs
	@echo "Swagger 2.0 documentation generated in docs/"

# Generate OpenAPI 3.0 documentation
docs-openapi3: docs
	@if ! command -v swagger2openapi >/dev/null 2>&1; then \
		echo "Installing swagger2openapi..."; \
		npm install -g swagger2openapi; \
	fi
	@echo "Converting Swagger 2.0 to OpenAPI 3.0..."
	swagger2openapi docs/swagger.yaml -o docs/openapi3.yaml
	swagger2openapi docs/swagger.json -o docs/openapi3.json
	@echo "OpenAPI 3.0 documentation generated in docs/openapi3.yaml and docs/openapi3.json"

# Generate TypeScript types from OpenAPI spec
generate-types: docs-openapi3
	@echo "Generating TypeScript types..."
	cd ../frontend/svetu && npx openapi-typescript ../../backend/docs/openapi3.yaml -o src/types/generated/api.ts && yarn format
	@echo "TypeScript types generated in frontend/svetu/src/types/generated/api.ts"

docs-fmt:
	go tool github.com/swaggo/swag/cmd/swag fmt

docs-clean:
	rm -rf docs

# Go formatting and linting
format:
	@echo "Formatting Go code with gofumpt..."
	go tool golang.org/x/tools/cmd/goimports -w .
	go tool mvdan.cc/gofumpt -w .
	@echo "Go code formatting completed!"

format-check:
	@echo "Checking Go code formatting..."
	@if [ -n "$$(go tool mvdan.cc/gofumpt -d .)" ]; then \
		echo "❌ Code is not formatted. Run 'make format' to fix."; \
		go tool mvdan.cc/gofumpt -d .; \
		exit 1; \
	fi
	@if [ -n "$$(go tool golang.org/x/tools/cmd/goimports -d .)" ]; then \
		echo "❌ Imports are not formatted. Run 'make format' to fix."; \
		go tool golang.org/x/tools/cmd/goimports -d .; \
		exit 1; \
	fi
	@echo "✅ Go code formatting is correct!"

lint:
	@echo "Running Go linters..."
	go tool github.com/golangci/golangci-lint/v2/cmd/golangci-lint run ./...
	@echo "✅ Linting completed!"

lint_ext: lint-govet lint-lll

lint-govet:
	@echo "Running govet linter only..."
	go tool github.com/golangci/golangci-lint/v2/cmd/golangci-lint run --enable-only=govet ./...

lint-lll:
	@echo "Running lll linter only..."
	go tool github.com/golangci/golangci-lint/v2/cmd/golangci-lint run --enable-only=lll ./...

test:
	@echo "Running Go tests..."
	go test ./... -v -count=1
	@echo "✅ Tests completed!"

test-short:
	@echo "Running Go tests (short mode)..."
	go test ./... -short -run "^(Test|Benchmark).*" -skip "Integration|OpenSearch"
	@echo "✅ Short tests completed!"

test-unit:
	@echo "Running unit tests only..."
	go test ./... -run "^(Test|Benchmark).*" -skip "Integration|OpenSearch|RealWorld"
	@echo "✅ Unit tests completed!"

test-coverage:
	@echo "Running Go tests with coverage..."
	@PKG_LIST=$$(go list ./... | grep -v '/cmd/'); \
	go test $$PKG_LIST -cover -coverprofile=coverage.out 2>&1 | grep -v "go: no such tool"; \
	if [ -f coverage.out ]; then \
		go tool cover -html=coverage.out -o coverage.html; \
		echo "✅ Coverage report generated: coverage.html"; \
		echo "📊 Coverage summary:"; \
		go tool cover -func=coverage.out | tail -1; \
	else \
		echo "❌ Failed to generate coverage.out"; \
		exit 1; \
	fi

test-coverage-detailed:
	@echo "Running Go tests with detailed coverage report..."
	@PKG_LIST=$$(go list ./... | grep -v '/cmd/'); \
	go test $$PKG_LIST -cover -coverprofile=coverage.out 2>&1 | grep -v "go: no such tool"; \
	if [ -f coverage.out ]; then \
		go tool cover -html=coverage.out -o coverage.html; \
		echo "✅ Coverage report generated: coverage.html"; \
		echo ""; \
		echo "📊 Top 20 files by coverage:"; \
		go tool cover -func=coverage.out | grep -v "total:" | sort -k3 -rn | head -20; \
		echo ""; \
		echo "📉 Bottom 20 files by coverage:"; \
		go tool cover -func=coverage.out | grep -v "total:" | sort -k3 -n | head -20; \
		echo ""; \
		echo "📊 Overall summary:"; \
		go tool cover -func=coverage.out | tail -1; \
	else \
		echo "❌ Failed to generate coverage.out"; \
		exit 1; \
	fi

test-coverage-json:
	@echo "Running Go tests and generating JSON coverage report..."
	@PKG_LIST=$$(go list ./... | grep -v '/cmd/'); \
	go test $$PKG_LIST -cover -coverprofile=coverage.out 2>&1 | grep -v "go: no such tool" > /dev/null; \
	if [ -f coverage.out ]; then \
		TOTAL=$$(go tool cover -func=coverage.out | tail -1 | awk '{print $$NF}' | sed 's/%//'); \
		echo "{\"total\": $$TOTAL, \"html\": \"coverage.html\"}"; \
	else \
		echo "{\"error\": \"Failed to generate coverage report\"}"; \
		exit 1; \
	fi

test-postexpress:
	@echo "Running Post Express integration test..."
	cd scripts && go run test_postexpress.go
	@echo "✅ Post Express test completed!"

test-redis:
	@echo "Running Redis cache integration test..."
	cd scripts && go run test_redis.go
	@echo "✅ Redis test completed!"

test-opensearch:
	@echo "Running OpenSearch integration test..."
	cd scripts && go run test_opensearch.go
	@echo "✅ OpenSearch test completed!"

test-postgres:
	@echo "Running PostgreSQL integration test..."
	cd scripts && go run test_postgres.go
	@echo "✅ PostgreSQL test completed!"

pre-commit: format-check lint test-short
	@echo "✅ Pre-commit checks passed!"

migrate_up:
	bin/migrator up

migrate_down:
	bin/migrator down

reindex_products:
	./reindex-products.sh

reindex_products_build:
	go build -o bin/reindex-products ./cmd/reindex-products/main.go

reindex_products_direct:
	go run ./cmd/reindex-products/main.go --batch 100

# Docker commands for migrator
docker_build_migrator:
	docker build -f Dockerfile.migrator -t migrator:local \
		--build-arg GIT_COMMIT=$$(git rev-parse --short HEAD) \
		--build-arg BUILD_TIME=$$(date +%Y%m%d_%H%M%S) \
		.

# Build with buildx (if installed)
docker_buildx_migrator:
	#docker buildx build -f Dockerfile.migrator -t migrator:local --load \
#		--build-arg GIT_COMMIT=$$(git rev-parse --short HEAD) \
#		--build-arg BUILD_TIME=$$(date +%Y%m%d_%H%M%S) \
#		.

.PHONY: watch
watch:
	go tool github.com/air-verse/air

docker_run_migrator_up:
	docker run --rm --network host \
		-e DATABASE_URL="${DATABASE_URL}" \
		-e MIGRATIONS_PATH="/migrations" \
		-e MIGRATION_DIRECTION="up" \
		migrator:local

docker_run_migrator_down:
	docker run --rm --network host \
		-e DATABASE_URL="${DATABASE_URL}" \
		-e MIGRATIONS_PATH="/migrations" \
		-e MIGRATION_DIRECTION="down" \
		migrator:local

# Development restart command - build, kill old processes and restart with logging
dev-restart:
	@echo "🔨 Building all binaries..."
	@$(MAKE) build_all
	@echo "📦 Renaming api to api_dev..."
	@mv bin/api bin/api_dev
	@echo "🔪 Killing old api_dev processes..."
	@for pid in $$(pgrep -f "^\./bin/api_dev" || echo ""); do \
		if [ -n "$$pid" ]; then \
			echo "  Killing PID: $$pid"; \
			kill $$pid 2>/dev/null || true; \
		fi; \
	done
	@echo "🚀 Starting api_dev in background with logging..."
	@nohup ./bin/api_dev > api_dev.log 2>&1 &
	@echo "✅ api_dev started! Check logs with: tail -f api_dev.log"
	@sleep 1
	@echo "📍 PID: $$(pgrep -f "^\./bin/api_dev" | head -n1)"

# ==============================================================================
# Production Readiness Testing
# ==============================================================================

# Start mock microservice for testing
mock-start:
	@echo "Starting mock gRPC microservice..."
	@go run tests/mocks/microservice_mock.go &
	@sleep 2
	@echo "✅ Mock microservice started on :50051 (gRPC) and :50052 (HTTP control)"
	@echo "   Control API: curl http://localhost:50052/health"

# Stop mock microservice
mock-stop:
	@echo "Stopping mock microservice..."
	@pkill -f "tests/mocks/microservice_mock" || true
	@echo "✅ Mock microservice stopped"

# Integration tests (timeout, circuit breaker)
test-integration:
	@echo "Running integration tests..."
	@go test -v -timeout=60s ./tests/integration/...
	@echo "✅ Integration tests completed"

# Integration tests with race detector
test-integration-race:
	@echo "Running integration tests with race detector..."
	@go test -v -race -timeout=60s ./tests/integration/...
	@echo "✅ Integration tests with race detector completed"

# E2E tests (requires running backend and mock service)
test-e2e:
	@echo "Running end-to-end tests..."
	@go test -v -tags=e2e -timeout=120s ./tests/e2e/...
	@echo "✅ E2E tests completed"

# Load tests (benchmarks)
test-load:
	@echo "Running load tests..."
	@go test -v -bench=. -benchtime=30s -timeout=5m ./tests/load/...
	@echo "✅ Load tests completed"

# All production readiness tests (full suite)
test-production-readiness:
	@echo "=========================================="
	@echo "Production Readiness Test Suite"
	@echo "=========================================="
	@echo ""
	@echo "Step 1: Starting mock microservice..."
	@$(MAKE) mock-start
	@echo ""
	@echo "Step 2: Running integration tests..."
	@$(MAKE) test-integration
	@echo ""
	@echo "Step 3: Running integration tests with race detector..."
	@$(MAKE) test-integration-race
	@echo ""
	@echo "Step 4: Running load tests..."
	@$(MAKE) test-load
	@echo ""
	@echo "Step 5: Stopping mock microservice..."
	@$(MAKE) mock-stop
	@echo ""
	@echo "=========================================="
	@echo "✅ All production readiness tests PASSED"
	@echo "=========================================="
	@echo ""
	@echo "Summary:"
	@echo "  - Timeout tests: PASSED"
	@echo "  - Circuit breaker tests: PASSED"
	@echo "  - Race condition tests: PASSED"
	@echo "  - Load tests: PASSED"
	@echo "  - Memory stability: PASSED"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Review test report: docs/migration/PRODUCTION_READINESS_TEST_REPORT.md"
	@echo "  2. Check Prometheus metrics during load tests"
	@echo "  3. Run E2E tests with real backend: make test-e2e"

# Quick smoke test (fast subset)
test-smoke:
	@echo "Running smoke tests..."
	@$(MAKE) mock-start
	@go test -v -short -timeout=30s ./tests/integration/timeout_test.go
	@$(MAKE) mock-stop
	@echo "✅ Smoke tests completed"

# Coverage report for production tests
test-coverage:
	@echo "Generating test coverage report..."
	@go test -v -coverprofile=coverage_production.out ./tests/integration/... ./tests/load/...
	@go tool cover -html=coverage_production.out -o coverage_production.html
	@echo "✅ Coverage report generated: coverage_production.html"
	@echo "   Open with: firefox coverage_production.html"

# ==============================================================================
# Sprint 6.3: Integration Testing for Microservice Migration
# ==============================================================================

# Run microservice connectivity tests
test-microservice-connectivity:
	@echo "Running microservice connectivity tests..."
	@go test -v -timeout=60s ./tests/integration/microservice_connectivity_test.go

# Run traffic router integration tests
test-traffic-router:
	@echo "Running traffic router integration tests..."
	@go test -v -timeout=60s ./tests/integration/traffic_router_integration_test.go

# Run data consistency tests
test-data-consistency:
	@echo "Running data consistency tests..."
	@go test -v -timeout=120s ./tests/integration/data_consistency_test.go

# Run performance & reliability tests
test-performance-reliability:
	@echo "Running performance and reliability tests..."
	@go test -v -timeout=180s ./tests/integration/performance_reliability_test.go

# Run smoke tests (quick validation <10 seconds)
test-microservice-smoke:
	@echo "Running microservice smoke tests..."
	@go test -v -timeout=30s ./tests/smoke/microservice_smoke_test.go

# Run all microservice integration tests (Sprint 6.3)
test-microservice-integration:
	@echo "=========================================="
	@echo "Sprint 6.3: Microservice Integration Tests"
	@echo "=========================================="
	@echo ""
	@echo "Step 1: Smoke tests (quick validation)..."
	@$(MAKE) test-microservice-smoke
	@echo ""
	@echo "Step 2: Connectivity tests..."
	@$(MAKE) test-microservice-connectivity
	@echo ""
	@echo "Step 3: Traffic router tests..."
	@$(MAKE) test-traffic-router
	@echo ""
	@echo "Step 4: Data consistency tests..."
	@$(MAKE) test-data-consistency
	@echo ""
	@echo "Step 5: Performance & reliability tests..."
	@$(MAKE) test-performance-reliability
	@echo ""
	@echo "=========================================="
	@echo "✅ ALL MICROSERVICE INTEGRATION TESTS COMPLETE"
	@echo "=========================================="

# Generate integration test coverage report
test-microservice-coverage:
	@echo "Generating microservice integration test coverage..."
	@go test -v -coverprofile=coverage_microservice_integration.out \
		./tests/integration/microservice_connectivity_test.go \
		./tests/integration/traffic_router_integration_test.go \
		./tests/integration/data_consistency_test.go \
		./tests/integration/performance_reliability_test.go \
		./tests/smoke/microservice_smoke_test.go
	@if [ -f coverage_microservice_integration.out ]; then \
		go tool cover -html=coverage_microservice_integration.out -o coverage_microservice_integration.html; \
		echo "✅ Integration coverage report: coverage_microservice_integration.html"; \
		echo "📊 Coverage summary:"; \
		go tool cover -func=coverage_microservice_integration.out | tail -1; \
	fi

