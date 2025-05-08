.PHONY: deploy

# Default target
deploy:
	@./scripts/blue_green_deploy_on_svetu.rs.sh $(filter-out $@,$(MAKECMDGOALS))

# This allows passing arguments to the deploy target
%:
	@:
