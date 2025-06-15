# Шаг 4: Docker с runtime поддержкой

## Цель
Создать универсальный Docker образ, который поддерживает изменение конфигурации через переменные окружения без пересборки.

## Задачи

### 4.1 Обновление Dockerfile

Файл: `/frontend/svetu/Dockerfile`

```dockerfile
# syntax=docker/dockerfile:1

# Build arguments
ARG NODE_VERSION=22
ARG ALPINE_VERSION=3.19

# ==========================================
# Base stage - common dependencies
# ==========================================
FROM node:${NODE_VERSION}-alpine${ALPINE_VERSION} AS base
RUN apk add --no-cache libc6-compat dumb-init
WORKDIR /app

# ==========================================
# Dependencies stage
# ==========================================
FROM base AS deps
COPY package.json yarn.lock ./
# Install production dependencies
RUN yarn install --frozen-lockfile --production && \
    yarn cache clean

# ==========================================
# Dev dependencies stage
# ==========================================
FROM base AS dev-deps
COPY package.json yarn.lock ./
# Install all dependencies (including dev)
RUN yarn install --frozen-lockfile && \
    yarn cache clean

# ==========================================
# Builder stage
# ==========================================
FROM base AS builder
WORKDIR /app

# Copy dependencies
COPY --from=dev-deps /app/node_modules ./node_modules
COPY . .

# Build arguments for compile-time config
ARG NEXT_PUBLIC_API_URL="__NEXT_PUBLIC_API_URL__"
ARG NEXT_PUBLIC_MINIO_URL="__NEXT_PUBLIC_MINIO_URL__"
ARG NEXT_PUBLIC_IMAGE_HOSTS="__NEXT_PUBLIC_IMAGE_HOSTS__"
ARG NEXT_PUBLIC_WEBSOCKET_URL="__NEXT_PUBLIC_WEBSOCKET_URL__"
ARG NEXT_PUBLIC_ENABLE_PAYMENTS="__NEXT_PUBLIC_ENABLE_PAYMENTS__"

# Set build-time env vars with placeholders
ENV NEXT_PUBLIC_API_URL=${NEXT_PUBLIC_API_URL}
ENV NEXT_PUBLIC_MINIO_URL=${NEXT_PUBLIC_MINIO_URL}
ENV NEXT_PUBLIC_IMAGE_HOSTS=${NEXT_PUBLIC_IMAGE_HOSTS}
ENV NEXT_PUBLIC_WEBSOCKET_URL=${NEXT_PUBLIC_WEBSOCKET_URL}
ENV NEXT_PUBLIC_ENABLE_PAYMENTS=${NEXT_PUBLIC_ENABLE_PAYMENTS}

# Disable telemetry during build
ENV NEXT_TELEMETRY_DISABLED=1

# Build the application
RUN yarn build

# ==========================================
# Production stage
# ==========================================
FROM base AS runner
WORKDIR /app

# Production environment
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

# Create non-root user
RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs

# Copy production files
COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

# Copy and setup entrypoint script
COPY --chown=nextjs:nodejs docker-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:3000/api/health || exit 1

# Switch to non-root user
USER nextjs

# Expose port
EXPOSE 3000

# Use dumb-init to handle signals properly
ENTRYPOINT ["dumb-init", "--"]
CMD ["/usr/local/bin/docker-entrypoint.sh"]
```

### 4.2 Создание entrypoint скрипта

Файл: `/frontend/svetu/docker-entrypoint.sh`

```bash
#!/bin/sh
set -e

echo "🚀 Starting Next.js application..."
echo "📋 Environment: ${NODE_ENV:-production}"
echo "🔧 Runtime configuration:"
echo "   API URL: ${NEXT_PUBLIC_API_URL:-not set}"
echo "   MinIO URL: ${NEXT_PUBLIC_MINIO_URL:-not set}"
echo "   WebSocket: ${NEXT_PUBLIC_WEBSOCKET_URL:-not set}"

# Function to replace placeholders in static files
replace_env_vars() {
  local search="$1"
  local replace="$2"
  
  echo "   Replacing $search with $replace"
  
  # Find and replace in all JS files
  find /app/.next/static -name "*.js" -type f | while read -r file; do
    # Use temp file to avoid issues with file permissions
    if grep -q "$search" "$file" 2>/dev/null; then
      sed "s|$search|$replace|g" "$file" > "$file.tmp" && mv "$file.tmp" "$file"
    fi
  done
}

# Replace placeholders with actual values
echo "🔄 Injecting runtime environment variables..."

# Required variables
if [ -n "$NEXT_PUBLIC_API_URL" ]; then
  replace_env_vars "__NEXT_PUBLIC_API_URL__" "$NEXT_PUBLIC_API_URL"
else
  echo "⚠️  Warning: NEXT_PUBLIC_API_URL not set, using build default"
fi

if [ -n "$NEXT_PUBLIC_MINIO_URL" ]; then
  replace_env_vars "__NEXT_PUBLIC_MINIO_URL__" "$NEXT_PUBLIC_MINIO_URL"
else
  echo "⚠️  Warning: NEXT_PUBLIC_MINIO_URL not set, using build default"
fi

# Optional variables
if [ -n "$NEXT_PUBLIC_IMAGE_HOSTS" ]; then
  replace_env_vars "__NEXT_PUBLIC_IMAGE_HOSTS__" "$NEXT_PUBLIC_IMAGE_HOSTS"
fi

if [ -n "$NEXT_PUBLIC_WEBSOCKET_URL" ]; then
  replace_env_vars "__NEXT_PUBLIC_WEBSOCKET_URL__" "$NEXT_PUBLIC_WEBSOCKET_URL"
fi

if [ -n "$NEXT_PUBLIC_ENABLE_PAYMENTS" ]; then
  replace_env_vars "__NEXT_PUBLIC_ENABLE_PAYMENTS__" "$NEXT_PUBLIC_ENABLE_PAYMENTS"
fi

echo "✅ Environment variables injected successfully"

# Create health check endpoint if it doesn't exist
if [ ! -f "/app/api/health/route.js" ]; then
  echo "📌 Creating health check endpoint..."
  mkdir -p /app/api/health
  cat > /app/api/health/route.js << 'EOF'
exports.GET = async function GET() {
  return new Response(JSON.stringify({ status: 'ok', timestamp: new Date().toISOString() }), {
    status: 200,
    headers: { 'Content-Type': 'application/json' },
  });
};
EOF
fi

# Start the application
echo "🎯 Starting Node.js server..."
exec node server.js
```

### 4.3 Docker Compose для разработки

Файл: `/frontend/svetu/docker-compose.yml`

```yaml
version: '3.8'

services:
  frontend:
    build:
      context: .
      dockerfile: Dockerfile
      target: runner
      args:
        NODE_VERSION: "22"
    ports:
      - "3001:3000"
    environment:
      # Runtime variables
      - NEXT_PUBLIC_API_URL=${NEXT_PUBLIC_API_URL:-http://localhost:3000}
      - NEXT_PUBLIC_MINIO_URL=${NEXT_PUBLIC_MINIO_URL:-http://localhost:9000}
      - NEXT_PUBLIC_IMAGE_HOSTS=${NEXT_PUBLIC_IMAGE_HOSTS:-http:localhost:9000,https:svetu.rs:443}
      - NEXT_PUBLIC_WEBSOCKET_URL=${NEXT_PUBLIC_WEBSOCKET_URL:-ws://localhost:3000}
      - NEXT_PUBLIC_ENABLE_PAYMENTS=${NEXT_PUBLIC_ENABLE_PAYMENTS:-false}
      # Server variables
      - INTERNAL_API_URL=http://backend:3000
      - NODE_ENV=production
    networks:
      - app-network
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:3000/api/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
    restart: unless-stopped

  # Development mode with hot reload
  frontend-dev:
    build:
      context: .
      dockerfile: Dockerfile.dev
    ports:
      - "3001:3000"
    volumes:
      - ./src:/app/src
      - ./public:/app/public
      - ./.env.local:/app/.env.local
    environment:
      - NODE_ENV=development
    networks:
      - app-network
    profiles:
      - dev

networks:
  app-network:
    driver: bridge
```

### 4.4 Dockerfile для разработки

Файл: `/frontend/svetu/Dockerfile.dev`

```dockerfile
FROM node:22-alpine

WORKDIR /app

# Install dependencies
COPY package.json yarn.lock ./
RUN yarn install && yarn cache clean

# Copy application code
COPY . .

# Expose port
EXPOSE 3000

# Start development server
CMD ["yarn", "dev"]
```

### 4.5 Makefile для удобства

Файл: `/frontend/svetu/Makefile`

```makefile
# Variables
DOCKER_IMAGE_NAME = svetu-frontend
DOCKER_TAG = latest
DOCKER_REGISTRY = harbor.svetu.rs/svetu

# Colors
GREEN = \033[0;32m
YELLOW = \033[0;33m
RED = \033[0;31m
NC = \033[0m

# Build Docker image
.PHONY: docker-build
docker-build:
	@echo "$(GREEN)Building Docker image...$(NC)"
	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_TAG) .

# Run Docker container locally
.PHONY: docker-run
docker-run:
	@echo "$(GREEN)Running Docker container...$(NC)"
	docker run -p 3001:3000 \
		-e NEXT_PUBLIC_API_URL=http://localhost:3000 \
		-e NEXT_PUBLIC_MINIO_URL=http://localhost:9000 \
		$(DOCKER_IMAGE_NAME):$(DOCKER_TAG)

# Run with docker-compose
.PHONY: docker-up
docker-up:
	@echo "$(GREEN)Starting services with docker-compose...$(NC)"
	docker-compose up -d

# Stop docker-compose services
.PHONY: docker-down
docker-down:
	@echo "$(YELLOW)Stopping services...$(NC)"
	docker-compose down

# Development mode with hot reload
.PHONY: docker-dev
docker-dev:
	@echo "$(GREEN)Starting development mode...$(NC)"
	docker-compose --profile dev up frontend-dev

# Push to registry
.PHONY: docker-push
docker-push: docker-build
	@echo "$(GREEN)Pushing image to registry...$(NC)"
	docker tag $(DOCKER_IMAGE_NAME):$(DOCKER_TAG) $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(DOCKER_TAG)
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE_NAME):$(DOCKER_TAG)

# Test runtime config
.PHONY: docker-test-runtime
docker-test-runtime: docker-build
	@echo "$(GREEN)Testing runtime configuration...$(NC)"
	@echo "$(YELLOW)Starting with default config...$(NC)"
	docker run -d --name test-default -p 3001:3000 $(DOCKER_IMAGE_NAME):$(DOCKER_TAG)
	@sleep 5
	@echo "$(YELLOW)Testing with custom config...$(NC)"
	docker run -d --name test-custom -p 3002:3000 \
		-e NEXT_PUBLIC_API_URL=https://api.example.com \
		-e NEXT_PUBLIC_ENABLE_PAYMENTS=true \
		$(DOCKER_IMAGE_NAME):$(DOCKER_TAG)
	@sleep 5
	@echo "$(GREEN)Check http://localhost:3001 and http://localhost:3002$(NC)"
	@echo "$(YELLOW)Press any key to cleanup...$(NC)"
	@read -n 1
	docker stop test-default test-custom
	docker rm test-default test-custom

# Clean up
.PHONY: docker-clean
docker-clean:
	@echo "$(RED)Cleaning up Docker resources...$(NC)"
	docker-compose down -v
	docker rmi $(DOCKER_IMAGE_NAME):$(DOCKER_TAG) || true

# Help
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  $(GREEN)make docker-build$(NC)       - Build Docker image"
	@echo "  $(GREEN)make docker-run$(NC)         - Run container locally"
	@echo "  $(GREEN)make docker-up$(NC)          - Start with docker-compose"
	@echo "  $(GREEN)make docker-down$(NC)        - Stop docker-compose"
	@echo "  $(GREEN)make docker-dev$(NC)         - Run in development mode"
	@echo "  $(GREEN)make docker-push$(NC)        - Push to registry"
	@echo "  $(GREEN)make docker-test-runtime$(NC) - Test runtime config"
	@echo "  $(GREEN)make docker-clean$(NC)       - Clean up resources"
```

### 4.6 GitHub Actions для сборки

Файл: `/frontend/svetu/.github/workflows/docker-build.yml`

```yaml
name: Build and Push Docker Image

on:
  push:
    branches: [main, develop]
    paths:
      - 'frontend/svetu/**'
      - '.github/workflows/docker-build.yml'
  pull_request:
    paths:
      - 'frontend/svetu/**'

env:
  REGISTRY: harbor.svetu.rs
  IMAGE_NAME: svetu/frontend

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Log in to Harbor
        if: github.event_name != 'pull_request'
        uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ secrets.HARBOR_USERNAME }}
          password: ${{ secrets.HARBOR_PASSWORD }}

      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
          tags: |
            type=ref,event=branch
            type=ref,event=pr
            type=sha,prefix={{branch}}-
            type=raw,value=latest,enable={{is_default_branch}}

      - name: Build and push Docker image
        uses: docker/build-push-action@v5
        with:
          context: ./frontend/svetu
          push: ${{ github.event_name != 'pull_request' }}
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=gha
          cache-to: type=gha,mode=max
          build-args: |
            NODE_VERSION=22
```

## Тестирование

### Локальное тестирование
```bash
# Собрать образ
make docker-build

# Запустить с дефолтной конфигурацией
make docker-run

# Запустить с кастомной конфигурацией
docker run -p 3001:3000 \
  -e NEXT_PUBLIC_API_URL=https://api.production.com \
  -e NEXT_PUBLIC_ENABLE_PAYMENTS=true \
  svetu-frontend:latest

# Проверить runtime замену
make docker-test-runtime
```

### Проверка в браузере
1. Открыть http://localhost:3001
2. Открыть DevTools → Console
3. Выполнить: `window.__NEXT_DATA__.runtimeConfig`
4. Убедиться что переменные соответствуют заданным

## Результат
После выполнения этого шага:
1. Docker образ будет поддерживать runtime конфигурацию
2. Один образ можно использовать для всех окружений
3. Переменные можно менять без пересборки
4. Образ оптимизирован по размеру и безопасности