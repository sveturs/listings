.PHONY: deploy docs docs-fmt docs-clean dev_backend_rebuild dev_frontend_rebuild dev_all_rebuild dev_migration_up

# Default target
deploy:
	@./scripts/blue_green_deploy_on_svetu.rs.sh $(filter-out $@,$(MAKECMDGOALS))

# Swagger документация
docs:
	cd backend && make docs && cd ..


docs-fmt:
	cd backend && make docs-fmt

docs-clean:
	cd backend && make docs-clean

# Docker Compose Dev commands
dev_backend_rebuild:
	@echo "🔨 Rebuilding and restarting backend container..."
	docker-compose -f docker-compose.dev.yml stop backend
	docker-compose -f docker-compose.dev.yml rm -f backend
	docker-compose -f docker-compose.dev.yml up -d --build backend
	@echo "✅ Backend container rebuilt and restarted"

dev_frontend_rebuild:
	@echo "🔨 Rebuilding and restarting frontend container..."
	docker-compose -f docker-compose.dev.yml stop frontend
	docker-compose -f docker-compose.dev.yml rm -f frontend
	docker-compose -f docker-compose.dev.yml up -d --build frontend
	@echo "✅ Frontend container rebuilt and restarted"

dev_all_rebuild: dev_backend_rebuild dev_frontend_rebuild
	@echo "✅ All containers rebuilt and restarted"

dev_migration_up:
	@echo "🔄 Running migrations in backend container..."
	docker-compose -f docker-compose.dev.yml exec backend /app/migrator migrate
	@echo "✅ Migrations completed"

# This allows passing arguments to the deploy target
%:
	@:
