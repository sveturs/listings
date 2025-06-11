.PHONY: docs docs-fmt docs-clean

build_migrator:
	go build -o bin/migrator ./cmd/migrator

build_api:
	go build -o bin/api ./cmd/api

build_reindex:
	go build -o bin/reindex ./cmd/reindex

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

migrate_up:
	bin/migrator up
