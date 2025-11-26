# Prometheus Monitoring Stack for Listings Microservice

Production-ready Prometheus configuration with comprehensive alerting, SLO tracking, and Grafana dashboards.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Quick Start](#quick-start)
- [Configuration Files](#configuration-files)
- [SLO Targets](#slo-targets)
- [Alert Rules](#alert-rules)
- [Recording Rules](#recording-rules)
- [Grafana Dashboards](#grafana-dashboards)
- [Testing](#testing)
- [Production Deployment](#production-deployment)
- [Troubleshooting](#troubleshooting)
- [Maintenance](#maintenance)

---

## Overview

This monitoring stack provides:

- **Metrics Collection**: 67+ custom metrics from listings microservice
- **Alert Management**: Critical, warning, and SLO-based alerts
- **SLO Tracking**: 99.9% availability, 99.5% success rate, P95 < 200ms, P99 < 1s
- **Performance Monitoring**: Request rates, error rates, latency percentiles
- **Resource Monitoring**: CPU, memory, disk, database, cache
- **Visualization**: Grafana dashboards for all key metrics

### Components

- **Prometheus** (v2.48.0): Metrics storage and alerting engine
- **Grafana** (v10.2.2): Visualization and dashboarding
- **Alertmanager** (v0.26.0): Alert routing and notifications
- **Node Exporter** (v1.7.0): System metrics (CPU, memory, disk)
- **Postgres Exporter** (v0.15.0): Database metrics
- **Redis Exporter** (v1.55.0): Cache metrics
- **Blackbox Exporter** (v0.24.0): HTTP endpoint probing

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Listings Microservice                    │
│                    :9093/metrics endpoint                    │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           │ scrape (15s interval)
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                       Prometheus                             │
│  - Metrics storage (15d retention)                          │
│  - Alert evaluation (every 15s)                             │
│  - Recording rules (every 30s)                              │
└──────────┬────────────────────────┬─────────────────────────┘
           │                        │
           │ alerts                 │ query
           ▼                        ▼
┌──────────────────────┐   ┌──────────────────────┐
│   Alertmanager       │   │      Grafana         │
│  - PagerDuty         │   │  - Dashboards        │
│  - Slack             │   │  - Visualization     │
│  - Email             │   │  - Alerting UI       │
└──────────────────────┘   └──────────────────────┘
```

---

## Quick Start

### Prerequisites

- Docker and Docker Compose installed
- Listings microservice running on port 9093
- At least 4GB RAM available
- 100GB disk space for metrics storage

### 1. Validate Configuration

Before starting, validate all configuration files:

```bash
cd /p/github.com/sveturs/listings/deployment/prometheus
./validate-config.sh
```

This will check:
- YAML syntax
- Prometheus configuration
- Alert rules syntax
- Recording rules syntax
- Alertmanager configuration
- Docker Compose configuration

### 2. Update Configuration

**IMPORTANT**: Replace placeholder values before production deployment!

#### Alertmanager Configuration

Edit `alertmanager.yml`:

```yaml
# Replace these placeholders:
pagerduty_configs:
  - routing_key: 'YOUR_PAGERDUTY_INTEGRATION_KEY'  # ← Replace this

slack_configs:
  - api_url: 'YOUR_SLACK_WEBHOOK_URL'  # ← Replace this
```

#### Docker Compose

Edit `docker-compose.yml`:

```yaml
# Change default Grafana password
environment:
  - GF_SECURITY_ADMIN_PASSWORD=admin123  # ← Change this!

# Update PostgreSQL connection
environment:
  - DATA_SOURCE_NAME=postgresql://user:pass@host:5432/db  # ← Update this
```

### 3. Start Monitoring Stack

```bash
# Start all services
docker compose up -d

# Check service health
docker compose ps

# View logs
docker compose logs -f prometheus
docker compose logs -f grafana
docker compose logs -f alertmanager
```

### 4. Verify Setup

**Prometheus:**
- URL: http://localhost:9090
- Health: http://localhost:9090/-/healthy
- Targets: http://localhost:9090/targets
- Alerts: http://localhost:9090/alerts

**Grafana:**
- URL: http://localhost:3030
- Login: admin / admin123 (change this!)
- Datasource should be auto-configured

**Alertmanager:**
- URL: http://localhost:9093
- Health: http://localhost:9093/-/healthy

### 5. Test Alerts

Run the alert testing suite:

```bash
./test-alerts.sh --full
```

Or use interactive mode:

```bash
./test-alerts.sh
```

---

## Configuration Files

### prometheus.yml

Main Prometheus configuration with scrape targets:

- **Scrape Intervals**:
  - Application services: 15s
  - Infrastructure: 30s

- **Targets**:
  - `listings-microservice`: Port 9093
  - `svetu-monolith`: Port 3000
  - `node-exporter`: Port 9100
  - `postgres-exporter`: Port 9187
  - `redis-exporter`: Port 9121
  - `blackbox-exporter`: Port 9115

### alerts.yml

Alert rules organized by severity:

**Critical Alerts (PagerDuty)**:
- `ServiceDown`: Health check failing (1m)
- `HighErrorRateCritical`: Error rate > 1% (5m)
- `HighLatencyP99Critical`: P99 > 2s (5m)
- `DatabaseConnectionsHigh`: Pool > 90% (5m)
- `DiskSpaceCritical`: Disk < 15% (5m)
- `MemoryUsageCritical`: Memory > 95% (5m)

**Warning Alerts (Slack)**:
- `HighErrorRateWarning`: Error rate > 0.5% (10m)
- `HighLatencyP95Warning`: P95 > 500ms (10m)
- `CacheHitRatioLow`: Hit ratio < 70% (15m)
- `CPUUsageHigh`: CPU > 70% (15m)
- `MemoryUsageHigh`: Memory > 80% (15m)
- `RateLimitRejectionsHigh`: > 10/s (10m)

**SLO Alerts**:
- `ErrorBudgetBurnRateFast`: 5x burn rate (5m)
- `ErrorBudgetBurnRateSlow`: 2x burn rate (30m)
- `AvailabilitySLOBreach`: < 99.9% (5m)
- `SuccessRateSLOBreach`: < 99.5% (5m)
- `LatencyP95SLOBreach`: > 200ms (10m)
- `LatencyP99SLOBreach`: > 1s (10m)

### recording_rules.yml

Pre-aggregated metrics for performance:

**Request Rate Metrics**:
- `job:http_requests:rate1m`, `rate5m`, `rate1h`
- `method:http_requests:rate1m`, `rate5m`
- `service:http_requests_per_second:sum`

**Error Rate Metrics**:
- `job:http_requests_error_rate:ratio1m`, `ratio5m`, `ratio1h`
- `job:http_requests_4xx_rate:ratio5m`
- `service:http_requests_5xx:sum`

**Latency Percentiles**:
- `job:http_request_duration:p50_5m`, `p95_5m`, `p99_5m`, `avg_5m`
- `method:http_request_duration:p50_5m`, `p95_5m`, `p99_5m`

**SLO Metrics**:
- `service:availability:ratio1h`, `ratio24h`
- `service:success_rate:ratio1h`, `ratio24h`
- `service:error_budget_remaining:ratio30d`
- `service:latency_slo_p95:ratio1h`, `latency_slo_p99:ratio1h`

**Database Metrics**:
- `job:db_connection_pool:utilization`
- `job:db_queries:rate1m`, `rate5m`
- `job:db_query_duration:p95_5m`, `p99_5m`

**Cache Metrics**:
- `job:cache_hit_ratio:rate5m`, `rate1h`
- `job:cache_operations:rate5m`
- `job:cache_evictions:rate5m`

---

## SLO Targets

### Service Level Objectives

| Metric | Target | Error Budget |
|--------|--------|--------------|
| **Availability** | 99.9% | 43.2 min/month |
| **Success Rate** | 99.5% | 3.6 hours/month |
| **P95 Latency** | < 200ms | - |
| **P99 Latency** | < 1s | - |

### Monitoring Windows

- **Fast Burn (5x)**: 1 hour window, alert after 5 minutes
- **Slow Burn (2x)**: 6 hour window, alert after 30 minutes
- **SLO Compliance**: 1 hour window for real-time tracking
- **SLO Reporting**: 24 hour and 30 day windows

---

## Alert Rules

### Alert Severity Levels

**Critical (PagerDuty - Immediate Response)**:
- Service completely unavailable
- High error rates affecting users
- Critical resource exhaustion
- Database connection failures

**Warning (Slack - Proactive Monitoring)**:
- Elevated error rates
- Performance degradation
- Resource usage trends
- Cache inefficiency

**Info (Email/Slack - Awareness)**:
- Deployment notifications
- Configuration changes
- Scheduled maintenance

### Alert Routing

Alerts are routed based on labels:

```yaml
severity: critical → PagerDuty (repeat every 1h)
severity: warning  → Slack (repeat every 12h)
severity: info     → Email (repeat every 24h)
slo: true          → PagerDuty + Slack
```

### Inhibition Rules

- Warning alerts suppressed when critical alert active
- All alerts suppressed when `ServiceDown` firing

---

## Recording Rules

Recording rules improve query performance by pre-computing expensive queries.

### Benefits

1. **Faster Queries**: Pre-aggregated data loads instantly
2. **Reduced Load**: Less CPU/memory for repeated queries
3. **Consistent Metrics**: Same calculation across all dashboards
4. **Historical Data**: Preserved for full retention period

### Usage in Queries

Instead of:
```promql
# Slow - calculated every time
sum(rate(http_requests_total[5m])) by (method)
```

Use:
```promql
# Fast - pre-calculated
method:http_requests:rate5m
```

---

## Grafana Dashboards

### Auto-Provisioning

Datasource and dashboards are automatically configured on startup.

### Recommended Dashboards

Import these community dashboards:

1. **Node Exporter Full** (ID: 1860)
   - Comprehensive system metrics
   - CPU, memory, disk, network

2. **PostgreSQL Database** (ID: 6417)
   - Database performance metrics
   - Query statistics, connections

3. **Go Processes** (ID: 6671)
   - Golang runtime metrics
   - Goroutines, GC, memory

### Custom Dashboards

Place custom dashboard JSON files in:
```
grafana/dashboards/
```

They will be auto-imported on Grafana startup.

### Creating Dashboards

1. Login to Grafana: http://localhost:3030
2. Create dashboard manually
3. Export JSON
4. Save to `grafana/dashboards/` directory
5. Commit to version control

---

## Testing

### Validation Script

```bash
./validate-config.sh
```

Checks:
- YAML syntax (yamllint)
- Prometheus config (promtool)
- Alert rules (promtool)
- Recording rules (promtool)
- Alertmanager config (amtool)
- Docker Compose config
- Placeholder values

### Alert Testing Script

```bash
# Full test suite
./test-alerts.sh --full

# Interactive mode
./test-alerts.sh

# Specific tests
./test-alerts.sh
  1. Full test suite
  2. List alert rules
  3. Show active alerts
  4. Test critical alerts
  5. Test SLO alerts
  6. Test recording rules
  7. Show scrape targets
  8. Check metrics availability
  9. Simulate test alert
  10. Reload Prometheus config
```

### Manual Testing

**Test Alert Query:**
```bash
# Check if alert would fire
curl -s 'http://localhost:9090/api/v1/query?query=up{job="listings-microservice"}==0' | jq .
```

**Send Test Alert:**
```bash
curl -X POST -H "Content-Type: application/json" \
  -d '[{"labels":{"alertname":"Test","severity":"warning"},"annotations":{"summary":"Test"}}]' \
  http://localhost:9093/api/v1/alerts
```

**Reload Config:**
```bash
# Reload Prometheus without restart
curl -X POST http://localhost:9090/-/reload

# Reload Alertmanager
curl -X POST http://localhost:9093/-/reload
```

---

## Production Deployment

### Pre-Deployment Checklist

- [ ] Validate all configurations
- [ ] Replace placeholder values (PagerDuty, Slack)
- [ ] Update database credentials
- [ ] Change default Grafana password
- [ ] Configure persistent volumes
- [ ] Set up backup strategy
- [ ] Configure firewall rules
- [ ] Set up TLS/SSL certificates
- [ ] Configure retention policies
- [ ] Test alert routing

### Security Hardening

**1. Change Default Passwords**

```yaml
# Grafana
environment:
  - GF_SECURITY_ADMIN_PASSWORD=${GRAFANA_ADMIN_PASSWORD}

# Use environment variables or secrets
```

**2. Enable Authentication**

```yaml
# Prometheus
command:
  - '--web.external-url=https://prometheus.example.com'
  - '--web.route-prefix=/'

# Add reverse proxy with auth (Nginx, Traefik)
```

**3. Restrict Network Access**

```yaml
# Only expose on localhost
ports:
  - "127.0.0.1:9090:9090"

# Or use internal network only
networks:
  - monitoring
```

**4. Enable TLS**

```yaml
# Use reverse proxy for TLS termination
# Or configure Prometheus with TLS
```

### Resource Limits

Add resource limits to prevent OOM:

```yaml
services:
  prometheus:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 4G
        reservations:
          cpus: '1'
          memory: 2G
```

### Backup Strategy

**Prometheus Data:**
```bash
# Create snapshot
curl -X POST http://localhost:9090/api/v1/admin/tsdb/snapshot

# Backup snapshot directory
tar -czf prometheus-backup-$(date +%Y%m%d).tar.gz \
  /var/lib/docker/volumes/listings-prometheus-data/_data/snapshots/
```

**Grafana Dashboards:**
```bash
# Export all dashboards
for dashboard in $(curl -s http://admin:admin123@localhost:3030/api/search | jq -r '.[].uid'); do
  curl -s http://admin:admin123@localhost:3030/api/dashboards/uid/$dashboard | \
    jq '.dashboard' > "dashboard-${dashboard}.json"
done
```

### High Availability

For production HA setup:

1. **Prometheus Federation**: Multiple Prometheus instances
2. **Thanos/Cortex**: Long-term storage and global view
3. **Alertmanager Cluster**: HA alert routing
4. **Grafana HA**: Multiple Grafana instances with shared DB

---

## Troubleshooting

### Prometheus Issues

**Prometheus not scraping targets:**
```bash
# Check targets status
curl http://localhost:9090/api/v1/targets | jq .

# Check Prometheus logs
docker compose logs prometheus

# Common issues:
# - Service not exposing metrics endpoint
# - Firewall blocking connection
# - DNS resolution failure
```

**High memory usage:**
```bash
# Check cardinality
curl http://localhost:9090/api/v1/status/tsdb | jq .

# Reduce retention or increase memory
# Add resource limits
```

**Alerts not firing:**
```bash
# Check alert rules loaded
curl http://localhost:9090/api/v1/rules | jq .

# Check alert state
curl http://localhost:9090/api/v1/alerts | jq .

# Validate alert query manually
curl 'http://localhost:9090/api/v1/query?query=up==0' | jq .
```

### Alertmanager Issues

**Alerts not being sent:**
```bash
# Check Alertmanager logs
docker compose logs alertmanager

# Check Alertmanager config
docker compose exec alertmanager amtool check-config /etc/alertmanager/alertmanager.yml

# Test receiver
docker compose exec alertmanager amtool config routes test \
  --config.file=/etc/alertmanager/alertmanager.yml \
  severity=critical service=listings
```

### Grafana Issues

**Datasource not working:**
```bash
# Check datasource health
curl http://admin:admin123@localhost:3030/api/datasources/1/health

# Check Grafana logs
docker compose logs grafana

# Re-provision datasources
docker compose restart grafana
```

**Dashboards not loading:**
```bash
# Check dashboard provisioning
docker compose exec grafana ls -la /var/lib/grafana/dashboards/

# Check provisioning config
docker compose exec grafana cat /etc/grafana/provisioning/dashboards/default.yml
```

### Common Issues

**"Too many open files":**
```bash
# Increase file descriptor limit
ulimit -n 65536

# Or add to docker-compose.yml
ulimits:
  nofile:
    soft: 65536
    hard: 65536
```

**Disk space issues:**
```bash
# Check disk usage
docker system df

# Clean old data
docker volume prune

# Reduce Prometheus retention
# Edit prometheus.yml: --storage.tsdb.retention.time=7d
```

---

## Maintenance

### Regular Tasks

**Daily:**
- Check active alerts
- Monitor disk usage
- Review error logs

**Weekly:**
- Review SLO compliance
- Check alert noise (too many/few alerts)
- Update dashboards based on feedback

**Monthly:**
- Review and update alert thresholds
- Clean up old data
- Update exporters and Prometheus
- Review security configurations

### Updating Components

```bash
# Pull latest images
docker compose pull

# Restart with new images
docker compose up -d

# Verify health
docker compose ps
./test-alerts.sh --full
```

### Monitoring the Monitors

Monitor Prometheus itself:

- Scrape success rate
- Rule evaluation time
- Query performance
- TSDB compaction time
- Disk usage trends

Add alerts for Prometheus health:

```yaml
- alert: PrometheusDown
  expr: up{job="prometheus"} == 0
  for: 5m

- alert: PrometheusDiskSpaceHigh
  expr: (prometheus_tsdb_storage_blocks_bytes / prometheus_tsdb_storage_available_bytes) > 0.85
  for: 10m
```

---

## Additional Resources

### Documentation

- [Prometheus Docs](https://prometheus.io/docs/)
- [Alertmanager Docs](https://prometheus.io/docs/alerting/latest/alertmanager/)
- [Grafana Docs](https://grafana.com/docs/)
- [PromQL Tutorial](https://prometheus.io/docs/prometheus/latest/querying/basics/)

### Tools

- [promtool](https://github.com/prometheus/prometheus): Prometheus CLI tool
- [amtool](https://github.com/prometheus/alertmanager): Alertmanager CLI tool
- [pint](https://github.com/cloudflare/pint): Prometheus rule linter

### Community Dashboards

- [Grafana Dashboards](https://grafana.com/grafana/dashboards/)
- [Awesome Prometheus](https://github.com/roaldnefs/awesome-prometheus)

---

## Support

For issues or questions:

1. Check [Troubleshooting](#troubleshooting) section
2. Review Prometheus/Grafana logs
3. Validate configuration with `validate-config.sh`
4. Contact platform team: platform-team@vondi.rs

---

## License

Internal use only - Svetu Platform Team
