# AI instruction when works with OpenApi, Swagger etc

## IMPORTANT RULES

- Когда изменяются Swagger аннотации на backend ВСЕГДА нужно запускать `make generate-types`

## API Documentation (Swagger)

Backend использует **goswag** для автоматической генерации OpenAPI/Swagger документации. См [@backend/Makefile](backend/Makefile)

### Сгенерированные файлы

После выполнения `make generate-types` создаются:
- `docs/docs.go` - Go код для встраивания документации
- `docs/openapi3.json` - OpenAPI v3 спецификация в JSON
- `docs/openapi3.yaml` - OpenAPI v3 спецификация в YAML
- `docs/swagger.json` - Swagger (Open API v2) спецификация в JSON
- `docs/swagger.yaml` - Swagger (Open API v2) спецификация в YAML

### Просмотр документации

Swagger UI доступен по адресу: http://localhost:3000/swagger/index.html (в режиме разработки)