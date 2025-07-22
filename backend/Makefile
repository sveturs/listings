.PHONY: docs docs-fmt docs-clean format format-check lint pre-commit

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

lint-errcheck:
	@echo "Running errcheck linter only..."
	go tool github.com/golangci/golangci-lint/v2/cmd/golangci-lint run --enable-only=errcheck ./...

lint-goconst:
	@echo "Running goconst linter only..."
	go tool github.com/golangci/golangci-lint/v2/cmd/golangci-lint run --enable-only=goconst ./...

lint-staticcheck:
	@echo "Running staticcheck linter only..."
	go tool github.com/golangci/golangci-lint/v2/cmd/golangci-lint run --enable-only=staticcheck ./...

pre-commit: format-check lint
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
