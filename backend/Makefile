.PHONY: docs docs-fmt docs-clean

# Swagger документация
docs:
	go tool github.com/swaggo/swag/cmd/swag init --parseDependency --parseInternal -g cmd/api/main.go -o docs --exclude ./internal/domain/models/custom_component.go

docs-fmt:
	go tool github.com/swaggo/swag/cmd/swag fmt

docs-clean:
	rm -rf docs