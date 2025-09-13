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

# This allows passing arguments to the deploy target
%:
	@:

restart-backend:
	docker-compose ps
	docker-compose stop backend
	docker-compose rm -f backend
	docker-compose up -d --build backend
	@echo "Run docker-compose logs -f backend"

restart-frontend:
	docker-compose ps
	docker-compose stop frontend
	docker-compose rm -f frontend
	docker-compose up -d --build frontend
	@echo "Run docker-compose logs -f frontend"

