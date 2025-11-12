# Makefile for Prometheus Monitoring Stack
# Provides convenient commands for managing the monitoring infrastructure

.PHONY: help validate start stop restart logs status clean test-alerts reload-prometheus reload-alertmanager backup

# Default target
.DEFAULT_GOAL := help

##@ General

help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Configuration

validate: ## Validate all configuration files
	@echo "Validating configuration..."
	@./validate-config.sh

install-tools: ## Install required tools (promtool, amtool, yamllint)
	@echo "Installing monitoring tools..."
	@echo "Note: This requires Go to be installed"
	@command -v go >/dev/null 2>&1 || { echo "Error: Go is not installed"; exit 1; }
	@echo "Installing promtool..."
	@go install github.com/prometheus/prometheus/cmd/promtool@latest
	@echo "Installing amtool..."
	@go install github.com/prometheus/alertmanager/cmd/amtool@latest
	@echo "Installing yamllint (requires pip)..."
	@pip3 install yamllint 2>/dev/null || echo "Warning: Could not install yamllint (pip3 required)"
	@echo "Done! Tools installed to $$GOPATH/bin"

##@ Lifecycle

start: ## Start the monitoring stack
	@echo "Starting monitoring stack..."
	@docker compose up -d
	@echo "Waiting for services to be healthy..."
	@sleep 5
	@make status

stop: ## Stop the monitoring stack
	@echo "Stopping monitoring stack..."
	@docker compose down

restart: ## Restart the monitoring stack
	@echo "Restarting monitoring stack..."
	@docker compose restart

start-dev: validate ## Validate config and start stack (development)
	@echo "Starting in development mode..."
	@make validate
	@make start
	@sleep 10
	@make test-alerts-quick

##@ Monitoring

status: ## Show status of all services
	@echo "Service Status:"
	@docker compose ps
	@echo ""
	@echo "Health Checks:"
	@echo -n "Prometheus:   "
	@curl -sf http://localhost:9090/-/healthy 2>/dev/null && echo "✓ Healthy" || echo "✗ Unhealthy"
	@echo -n "Grafana:      "
	@curl -sf http://localhost:3030/api/health 2>/dev/null && echo "✓ Healthy" || echo "✗ Unhealthy"
	@echo -n "Alertmanager: "
	@curl -sf http://localhost:9093/-/healthy 2>/dev/null && echo "✓ Healthy" || echo "✗ Unhealthy"

logs: ## Show logs from all services
	@docker compose logs -f

logs-prometheus: ## Show Prometheus logs
	@docker compose logs -f prometheus

logs-grafana: ## Show Grafana logs
	@docker compose logs -f grafana

logs-alertmanager: ## Show Alertmanager logs
	@docker compose logs -f alertmanager

##@ Testing

test-alerts: ## Run full alert test suite
	@./test-alerts.sh --full

test-alerts-quick: ## Quick alert test (active alerts only)
	@echo "Checking active alerts..."
	@curl -sf http://localhost:9090/api/v1/alerts | jq -r '.data.alerts[] | "\(.labels.alertname): \(.state)"' 2>/dev/null || echo "No alerts or jq not installed"

test-alerts-interactive: ## Run interactive alert testing
	@./test-alerts.sh

simulate-alert: ## Send test alert to Alertmanager
	@./test-alerts.sh --simulate

##@ Operations

reload-prometheus: ## Reload Prometheus configuration (no restart)
	@echo "Reloading Prometheus configuration..."
	@curl -sf -X POST http://localhost:9090/-/reload && echo "✓ Reloaded" || echo "✗ Failed"
	@echo "Note: Ensure Prometheus was started with --web.enable-lifecycle"

reload-alertmanager: ## Reload Alertmanager configuration (no restart)
	@echo "Reloading Alertmanager configuration..."
	@curl -sf -X POST http://localhost:9093/-/reload && echo "✓ Reloaded" || echo "✗ Failed"

check-targets: ## Show Prometheus scrape targets status
	@echo "Scrape Targets:"
	@curl -sf http://localhost:9090/api/v1/targets | jq -r '.data.activeTargets[] | "\(.labels.job): \(.health) (last scrape: \(.lastScrape))"' 2>/dev/null || echo "Failed to fetch targets"

check-rules: ## Show loaded alert and recording rules
	@echo "Loaded Rules:"
	@curl -sf http://localhost:9090/api/v1/rules | jq -r '.data.groups[] | "Group: \(.name) (\(.rules | length) rules)"' 2>/dev/null || echo "Failed to fetch rules"

check-metrics: ## Verify key metrics are available
	@echo "Checking key metrics..."
	@echo -n "http_requests_total: "
	@curl -sf 'http://localhost:9090/api/v1/query?query=http_requests_total' | jq -r '.data.result | length' 2>/dev/null || echo "0"
	@echo -n "go_goroutines: "
	@curl -sf 'http://localhost:9090/api/v1/query?query=go_goroutines' | jq -r '.data.result | length' 2>/dev/null || echo "0"
	@echo -n "up: "
	@curl -sf 'http://localhost:9090/api/v1/query?query=up' | jq -r '.data.result | length' 2>/dev/null || echo "0"

##@ Maintenance

clean: ## Remove all data (WARNING: deletes metrics!)
	@echo "WARNING: This will delete all metrics data!"
	@echo "Press Ctrl+C to cancel, or wait 5 seconds to continue..."
	@sleep 5
	@docker compose down -v
	@echo "All data removed"

backup: ## Create backup of metrics data
	@echo "Creating backup..."
	@mkdir -p backups
	@BACKUP_FILE="backups/prometheus-backup-$$(date +%Y%m%d-%H%M%S).tar.gz"; \
	echo "Backing up to $$BACKUP_FILE..."; \
	docker compose exec -T prometheus tar czf - /prometheus > $$BACKUP_FILE && \
	echo "✓ Backup created: $$BACKUP_FILE" || \
	echo "✗ Backup failed"

backup-grafana: ## Export all Grafana dashboards
	@echo "Exporting Grafana dashboards..."
	@mkdir -p backups/grafana-dashboards
	@for uid in $$(curl -sf http://admin:admin123@localhost:3030/api/search | jq -r '.[].uid'); do \
		curl -sf http://admin:admin123@localhost:3030/api/dashboards/uid/$$uid | \
		jq '.dashboard' > backups/grafana-dashboards/dashboard-$$uid.json; \
		echo "Exported dashboard: $$uid"; \
	done
	@echo "✓ Dashboards exported to backups/grafana-dashboards/"

restore: ## Restore from backup (requires BACKUP_FILE variable)
	@if [ -z "$(BACKUP_FILE)" ]; then \
		echo "Error: Please specify BACKUP_FILE"; \
		echo "Usage: make restore BACKUP_FILE=backups/prometheus-backup-20250105-120000.tar.gz"; \
		exit 1; \
	fi
	@echo "Restoring from $(BACKUP_FILE)..."
	@docker compose exec -T prometheus tar xzf - -C / < $(BACKUP_FILE)
	@echo "✓ Restored. Restarting Prometheus..."
	@docker compose restart prometheus

update: ## Update all container images
	@echo "Pulling latest images..."
	@docker compose pull
	@echo "Recreating containers with new images..."
	@docker compose up -d
	@make status

##@ Development

shell-prometheus: ## Open shell in Prometheus container
	@docker compose exec prometheus sh

shell-grafana: ## Open shell in Grafana container
	@docker compose exec grafana sh

shell-alertmanager: ## Open shell in Alertmanager container
	@docker compose exec alertmanager sh

query: ## Execute PromQL query (requires QUERY variable)
	@if [ -z "$(QUERY)" ]; then \
		echo "Error: Please specify QUERY"; \
		echo "Usage: make query QUERY='up'"; \
		exit 1; \
	fi
	@curl -sf "http://localhost:9090/api/v1/query?query=$(QUERY)" | jq .

##@ URLs

urls: ## Show all service URLs
	@echo "Service URLs:"
	@echo "  Prometheus:   http://localhost:9090"
	@echo "  Grafana:      http://localhost:3030 (admin/admin123)"
	@echo "  Alertmanager: http://localhost:9093"
	@echo ""
	@echo "Metrics Endpoints:"
	@echo "  Listings:     http://localhost:9093/metrics"
	@echo "  Monolith:     http://localhost:3000/metrics"

open: ## Open all UIs in browser
	@echo "Opening Prometheus..."
	@xdg-open http://localhost:9090 2>/dev/null || open http://localhost:9090 2>/dev/null || echo "Please open http://localhost:9090"
	@sleep 1
	@echo "Opening Grafana..."
	@xdg-open http://localhost:3030 2>/dev/null || open http://localhost:3030 2>/dev/null || echo "Please open http://localhost:3030"
	@sleep 1
	@echo "Opening Alertmanager..."
	@xdg-open http://localhost:9093 2>/dev/null || open http://localhost:9093 2>/dev/null || echo "Please open http://localhost:9093"

##@ Documentation

docs: ## Show documentation files
	@echo "Documentation:"
	@echo "  README.md       - Full documentation"
	@echo "  QUICK_START.md  - Quick start guide"
	@echo ""
	@echo "Configuration:"
	@echo "  prometheus.yml           - Scraping configuration"
	@echo "  alerts.yml              - Alert rules"
	@echo "  recording_rules.yml     - Recording rules"
	@echo "  alertmanager.yml        - Alertmanager routing"
	@echo ""
	@echo "Scripts:"
	@echo "  validate-config.sh  - Validate configurations"
	@echo "  test-alerts.sh      - Test alert rules"
