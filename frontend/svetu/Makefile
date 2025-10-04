generate_types:
	node scripts

# Docker commands
DOCKER_IMAGE_NAME := svetu-frontend
DOCKER_REGISTRY := harbor.svetu.rs/svetu/frontend
GIT_HASH := $(shell git rev-parse --short HEAD)

.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t svetu-frontend:latest .
	#docker build -t $(DOCKER_IMAGE_NAME):$(GIT_HASH) .

.PHONY: docker-build-prod
docker-build-prod:
	@echo "Building Docker image for production..."
	docker build --platform linux/amd64 -t $(DOCKER_REGISTRY):latest .
	docker build --platform linux/amd64 -t $(DOCKER_REGISTRY):$(GIT_HASH) .

.PHONY: docker-push
docker-push:
	@echo "Pushing Docker image to Harbor..."
	docker push $(DOCKER_REGISTRY):latest
	docker push $(DOCKER_REGISTRY):$(GIT_HASH)

.PHONY: docker-run
docker-run:
	@echo "Running Docker container locally..."
	docker run -p 3001:3000 \
		-e NEXT_PUBLIC_API_URL=http://localhost:3000 \
		-e NEXT_PUBLIC_MINIO_URL=http://localhost:9000 \
		-e NEXT_PUBLIC_IMAGE_HOSTS="http:localhost:9000,http:localhost:3000" \
		-e NEXT_PUBLIC_WEBSOCKET_URL=ws://localhost:3000 \
		-e NEXT_PUBLIC_ENABLE_PAYMENTS=false \
		$(DOCKER_IMAGE_NAME):latest

.PHONY: docker-run-prod
docker-run-prod:
	@echo "Running Docker container with production settings..."
	docker run -p 3001:3000 \
		-e NODE_ENV=production \
		-e NEXT_PUBLIC_API_URL=https://api.svetu.rs \
		-e NEXT_PUBLIC_MINIO_URL=https://svetu.rs \
		-e NEXT_PUBLIC_IMAGE_HOSTS="https:svetu.rs:443" \
		-e NEXT_PUBLIC_WEBSOCKET_URL=wss://ws.svetu.rs \
		-e NEXT_PUBLIC_ENABLE_PAYMENTS=true \
		$(DOCKER_REGISTRY):latest

.PHONY: docker-health
docker-health:
	@echo "Checking Docker container health..."
	@container_id=$$(docker ps -q -f ancestor=$(DOCKER_IMAGE_NAME):latest); \
	if [ -n "$$container_id" ]; then \
		docker inspect --format='{{.State.Health.Status}}' $$container_id; \
		echo "Checking health endpoint..."; \
		curl -f http://localhost:3001/api/health | jq .; \
	else \
		echo "No running container found"; \
		exit 1; \
	fi

.PHONY: docker-logs
docker-logs:
	@echo "Showing Docker container logs..."
	@container_id=$$(docker ps -q -f ancestor=$(DOCKER_IMAGE_NAME):latest); \
	if [ -n "$$container_id" ]; then \
		docker logs -f $$container_id; \
	else \
		echo "No running container found"; \
		exit 1; \
	fi

.PHONY: docker-stop
docker-stop:
	@echo "Stopping Docker container..."
	@container_id=$$(docker ps -q -f ancestor=$(DOCKER_IMAGE_NAME):latest); \
	if [ -n "$$container_id" ]; then \
		docker stop $$container_id; \
		echo "Container stopped"; \
	else \
		echo "No running container found"; \
	fi

.PHONY: docker-clean
docker-clean:
	@echo "Cleaning up Docker images..."
	docker rmi $(DOCKER_IMAGE_NAME):latest || true
	docker rmi $(DOCKER_IMAGE_NAME):$(GIT_HASH) || true

.PHONY: docker-compose-up
docker-compose-up:
	@echo "Starting frontend with docker-compose..."
	docker-compose up -d
	@echo "Frontend is running on http://localhost:3000"
	@echo "Checking health status..."
	@sleep 5
	@docker-compose ps

.PHONY: docker-compose-down
docker-compose-down:
	@echo "Stopping docker-compose services..."
	docker-compose down

.PHONY: docker-compose-logs
docker-compose-logs:
	@echo "Showing docker-compose logs..."
	docker-compose logs -f frontend

# Development restart command - stop old process and start new one on port 3003
.PHONY: dev-restart
dev-restart:
	@echo "🔪 Stopping frontend on port 3003..."
	@lsof -ti:3003 | xargs -r kill 2>/dev/null || true
	@sleep 1
	@echo "🚀 Starting frontend on port 3003..."
	@nohup yarn dev -p 3003 > frontend-dev.log 2>&1 &
	@sleep 2
	@echo "✅ Frontend running! Logs: tail -f frontend-dev.log"

# Production build and restart for dev server
.PHONY: dev-build-restart
dev-build-restart:
	@echo "🔪 Stopping frontend on port 3003..."
	@lsof -ti:3003 | xargs -r kill 2>/dev/null || true
	@fuser -k 3003/tcp 2>/dev/null || true
	@sleep 2
	@echo "🏗️  Building production version..."
	@yarn build
	@echo "🚀 Starting production server on port 3003..."
	@nohup yarn start -p 3003 > frontend-dev.log 2>&1 &
	@sleep 3
	@echo "✅ Production frontend running! Logs: tail -f frontend-dev.log"

# Start production server (without build)
.PHONY: dev-start
dev-start:
	@echo "🚀 Starting production server on port 3003..."
	@nohup yarn start -p 3003 > frontend-dev.log 2>&1 &
	@sleep 3
	@echo "✅ Production frontend running! Logs: tail -f frontend-dev.log"