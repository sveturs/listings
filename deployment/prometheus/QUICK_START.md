# Prometheus Quick Start Guide

5-minute setup guide for local development.

## Prerequisites

```bash
# Ensure listings service is running
curl http://localhost:9093/metrics

# Expected: Metrics output starting with "# HELP"
```

## Setup in 3 Steps

### 1. Validate & Configure

```bash
cd /p/github.com/sveturs/listings/deployment/prometheus

# Validate configuration
./validate-config.sh

# (Optional) Update alertmanager.yml with your Slack/PagerDuty keys
# For development, can skip this step
```

### 2. Start Stack

```bash
# Start all services
docker compose up -d

# Wait 30 seconds for startup
sleep 30

# Check status
docker compose ps
```

Expected output:
```
NAME                          STATUS
listings-prometheus           Up (healthy)
listings-grafana              Up (healthy)
listings-alertmanager         Up (healthy)
listings-node-exporter        Up
listings-blackbox-exporter    Up
```

### 3. Verify

```bash
# Test Prometheus
curl -s http://localhost:9090/-/healthy
# Expected: "Prometheus is Healthy."

# Test Grafana
curl -s http://localhost:3030/api/health
# Expected: {"database": "ok", ...}

# Check metrics are being scraped
curl -s http://localhost:9090/api/v1/targets | jq '.data.activeTargets[] | {job: .labels.job, health: .health}'
```

## Access UIs

| Service | URL | Credentials |
|---------|-----|-------------|
| **Prometheus** | http://localhost:9090 | - |
| **Grafana** | http://localhost:3030 | admin / admin123 |
| **Alertmanager** | http://localhost:9093 | - |

## Quick Tests

### Test Alerts

```bash
# Run full test suite
./test-alerts.sh --full

# Or interactive mode
./test-alerts.sh
```

### View Metrics

```bash
# Request rate
curl -s 'http://localhost:9090/api/v1/query?query=rate(http_requests_total[5m])' | jq .

# P95 latency
curl -s 'http://localhost:9090/api/v1/query?query=histogram_quantile(0.95,rate(http_request_duration_seconds_bucket[5m]))' | jq .

# Error rate
curl -s 'http://localhost:9090/api/v1/query?query=sum(rate(http_requests_total{status=~"5.."}[5m]))/sum(rate(http_requests_total[5m]))' | jq .
```

### Simulate Test Alert

```bash
./test-alerts.sh --simulate

# Check in Alertmanager UI
open http://localhost:9093
```

## Common Commands

```bash
# View logs
docker compose logs -f prometheus
docker compose logs -f grafana

# Restart service
docker compose restart prometheus

# Reload config (without restart)
curl -X POST http://localhost:9090/-/reload

# Stop stack
docker compose down

# Stop and remove volumes
docker compose down -v
```

## Grafana Setup

1. Open http://localhost:3030
2. Login: admin / admin123
3. Datasource "Prometheus" is auto-configured
4. Import dashboards:
   - Go to Dashboards → Import
   - Enter dashboard ID:
     - **1860** - Node Exporter Full
     - **6417** - PostgreSQL Database
     - **6671** - Go Processes
   - Select "Prometheus" datasource
   - Click Import

## Troubleshooting

### No Metrics Showing

```bash
# Check if listings service is exposing metrics
curl http://localhost:9093/metrics | head -20

# Check Prometheus targets
curl http://localhost:9090/api/v1/targets | jq '.data.activeTargets[] | {job: .labels.job, health: .health, lastError: .lastError}'

# Check Prometheus logs
docker compose logs prometheus | tail -50
```

### Alerts Not Working

```bash
# Check alert rules loaded
curl http://localhost:9090/api/v1/rules | jq '.data.groups[].name'

# Check active alerts
curl http://localhost:9090/api/v1/alerts | jq '.data.alerts[] | {name: .labels.alertname, state: .state}'

# Validate alert rules
docker compose exec prometheus promtool check rules /etc/prometheus/alerts.yml
```

### Grafana Can't Connect to Prometheus

```bash
# Test datasource from Grafana container
docker compose exec grafana wget -O- http://prometheus:9090/api/v1/status/config

# Check network
docker compose exec grafana ping prometheus

# Restart Grafana
docker compose restart grafana
```

## Next Steps

After basic setup is working:

1. Review and customize alert thresholds in `alerts.yml`
2. Add Slack webhook to `alertmanager.yml`
3. Create custom Grafana dashboards
4. Set up long-term storage (Thanos/Cortex)
5. Configure backup automation
6. Review [Full Documentation](README.md)

## Quick Reference

### PromQL Examples

```promql
# Request rate (per second)
rate(http_requests_total[5m])

# Error rate (percentage)
(sum(rate(http_requests_total{status=~"5.."}[5m])) / sum(rate(http_requests_total[5m]))) * 100

# P95 latency
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# Database connections usage
(listings_db_connections_open / listings_db_connections_max) * 100

# Cache hit ratio
(rate(listings_cache_hits_total[5m]) / (rate(listings_cache_hits_total[5m]) + rate(listings_cache_misses_total[5m]))) * 100

# Goroutines count
go_goroutines

# Memory usage
process_resident_memory_bytes / (1024 * 1024)  # in MB
```

### API Endpoints

```bash
# Prometheus
http://localhost:9090/api/v1/query?query=up          # Instant query
http://localhost:9090/api/v1/query_range?query=...   # Range query
http://localhost:9090/api/v1/targets                  # Scrape targets
http://localhost:9090/api/v1/alerts                   # Active alerts
http://localhost:9090/api/v1/rules                    # Alert rules
http://localhost:9090/api/v1/status/config            # Configuration

# Alertmanager
http://localhost:9093/api/v1/alerts                   # Active alerts
http://localhost:9093/api/v1/silences                 # Silences
http://localhost:9093/api/v1/status                   # Status

# Grafana
http://localhost:3030/api/health                      # Health check
http://localhost:3030/api/datasources                 # Datasources
http://localhost:3030/api/dashboards/home             # Home dashboard
```

## Need Help?

- **Full docs**: [README.md](README.md)
- **Validation**: `./validate-config.sh`
- **Testing**: `./test-alerts.sh`
- **Logs**: `docker compose logs <service>`
