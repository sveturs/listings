.PHONY: deploy docs docs-fmt docs-clean

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
