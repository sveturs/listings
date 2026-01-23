# Prometheus Monitoring Stack - Technical Overview

Production-ready monitoring infrastructure for Listings Microservice.

## Stack Summary

| Component | Version | Port | Purpose |
|-----------|---------|------|---------|
| **Prometheus** | v2.48.0 | 9090 | Metrics storage, alerting engine |
| **Grafana** | v10.2.2 | 3030 | Visualization, dashboards |
| **Alertmanager** | v0.26.0 | 9093 | Alert routing, notifications |
| **Node Exporter** | v1.7.0 | 9100 | System metrics (CPU, RAM, disk) |
| **Postgres Exporter** | v0.15.0 | 9187 | Database metrics |
| **Redis Exporter** | v1.55.0 | 9121 | Cache metrics |
| **Blackbox Exporter** | v0.24.0 | 9115 | HTTP probing, uptime |

## Metrics Coverage

### Application Metrics (67+ metrics)

**HTTP Metrics:**
- `http_requests_total` - Total requests by method, status
- `http_request_duration_seconds` - Request latency histogram
- `http_request_size_bytes` - Request payload sizes
- `http_response_size_bytes` - Response payload sizes

**Database Metrics:**
- `listings_db_connections_open` - Active connections
- `listings_db_connections_idle` - Idle connections
- `listings_db_connections_max` - Max pool size
- `listings_db_queries_total` - Query count by operation
- `listings_db_query_duration_seconds` - Query latency
- `listings_db_errors_total` - Database errors
- `listings_db_connection_wait_duration_seconds` - Wait time for connection

**Cache Metrics:**
- `listings_cache_hits_total` - Cache hits
- `listings_cache_misses_total` - Cache misses
- `listings_cache_operations_total` - Operations (get, set, delete)
- `listings_cache_errors_total` - Cache errors
- `listings_cache_evictions_total` - Evicted keys

**Business Metrics:**
- `listings_created_total` - Listings created
- `listings_updated_total` - Listings updated
- `listings_deleted_total` - Listings deleted
- `listings_active_total` - Active listings count
- `listings_by_status` - Distribution by status
- `listings_by_category` - Distribution by category

**Rate Limiting:**
- `listings_rate_limit_rejections_total` - Rejected requests

**Go Runtime:**
- `go_goroutines` - Goroutine count
- `go_threads` - OS threads
- `go_memstats_*` - Memory statistics
- `go_gc_duration_seconds` - GC duration

### Infrastructure Metrics

**System (via Node Exporter):**
- CPU usage, load average
- Memory usage, swap
- Disk I/O, space usage
- Network throughput, errors

**Database (via Postgres Exporter):**
- Connection pool stats
- Query performance
- Table sizes, bloat
- Replication lag
- Lock stats, deadlocks

**Cache (via Redis Exporter):**
- Memory usage
- Hit/miss rates
- Eviction stats
- Command statistics

## Alert Configuration

### Alert Distribution

| Severity | Count | Destination | Repeat Interval |
|----------|-------|-------------|-----------------|
| **Critical** | 6 | PagerDuty | 1 hour |
| **Warning** | 8 | Slack | 12 hours |
| **SLO** | 6 | PagerDuty + Slack | 2 hours |

### Alert Categories

**Service Health (Critical):**
1. `ServiceDown` - Health check failing (1m)
2. `HighErrorRateCritical` - >1% errors (5m)
3. `HighLatencyP99Critical` - P99 >2s (5m)

**Resource Health (Critical):**
4. `DatabaseConnectionsHigh` - >90% pool (5m)
5. `DiskSpaceCritical` - <15% free (5m)
6. `MemoryUsageCritical` - >95% used (5m)

**Performance (Warning):**
7. `HighErrorRateWarning` - >0.5% errors (10m)
8. `HighLatencyP95Warning` - P95 >500ms (10m)
9. `CacheHitRatioLow` - <70% hit rate (15m)

**Resource Trends (Warning):**
10. `CPUUsageHigh` - >70% (15m)
11. `MemoryUsageHigh` - >80% (15m)
12. `RateLimitRejectionsHigh` - >10/s (10m)

**Database Health (Warning):**
13. `DatabaseQueryDurationHigh` - P95 >1s (10m)
14. `GoroutinesHigh` - >1000 (15m)

**SLO Compliance:**
15. `ErrorBudgetBurnRateFast` - 5x burn (5m)
16. `ErrorBudgetBurnRateSlow` - 2x burn (30m)
17. `AvailabilitySLOBreach` - <99.9% (5m)
18. `SuccessRateSLOBreach` - <99.5% (5m)
19. `LatencyP95SLOBreach` - >200ms (10m)
20. `LatencyP99SLOBreach` - >1s (10m)

## Recording Rules

### Performance Optimization

Recording rules pre-compute expensive queries, improving dashboard load times and reducing Prometheus CPU usage.

**Request Rate Rules (6 rules):**
- `job:http_requests:rate1m/5m/1h`
- `method:http_requests:rate1m/5m`
- `service:http_requests_per_second:sum`

**Error Rate Rules (4 rules):**
- `job:http_requests_error_rate:ratio1m/5m/1h`
- `job:http_requests_4xx_rate:ratio5m`

**Latency Rules (7 rules):**
- `job:http_request_duration:p50/p95/p99/avg_5m`
- `method:http_request_duration:p50/p95/p99_5m`

**SLO Rules (7 rules):**
- `service:availability:ratio1h/24h`
- `service:success_rate:ratio1h/24h`
- `service:error_budget_remaining:ratio30d`
- `service:latency_slo_p95/p99:ratio1h`

**Database Rules (7 rules):**
- `job:db_connection_pool:utilization`
- `job:db_queries:rate1m/5m`
- `job:db_query_duration:p95/p99_5m`

**Cache Rules (5 rules):**
- `job:cache_hit_ratio:rate5m/1h`
- `job:cache_operations:rate5m`
- `job:cache_evictions:rate5m`

**System Rules (5 rules):**
- `instance:memory_usage:ratio`
- `instance:cpu_usage:ratio`
- `instance:disk_usage:ratio`
- `instance:disk_io:rate5m`
- `instance:network_throughput:rate5m`

**Total: 48 recording rules**

## SLO Targets

### Service Level Objectives

| Metric | Target | Measurement Window | Error Budget |
|--------|--------|-------------------|--------------|
| **Availability** | 99.9% | 30 days | 43.2 min/month |
| **Success Rate** | 99.5% | 30 days | 3.6 hours/month |
| **P95 Latency** | <200ms | Real-time | N/A |
| **P99 Latency** | <1s | Real-time | N/A |

### Error Budget Burn Rates

| Rate | Window | Alert Threshold | Response Time |
|------|--------|-----------------|---------------|
| **Fast (5x)** | 1 hour | 5 minutes | Immediate (PagerDuty) |
| **Slow (2x)** | 6 hours | 30 minutes | High priority (Slack) |

### SLO Calculation

```promql
# Availability (uptime)
(
  sum(rate(http_requests_total{status!~"5.."}[1h]))
  /
  sum(rate(http_requests_total[1h]))
) * 100

# Success Rate (successful requests)
(
  sum(rate(http_requests_total{status=~"2.."}[1h]))
  /
  sum(rate(http_requests_total[1h]))
) * 100

# Error Budget Remaining (monthly)
1 - (
  sum(increase(http_requests_total{status=~"5.."}[30d]))
  /
  sum(increase(http_requests_total[30d]))
) / 0.001
```

## Data Retention

| Data Type | Retention | Storage Size | Compaction |
|-----------|-----------|--------------|------------|
| **Raw metrics** | 15 days | ~50GB | 2 hour blocks |
| **Recording rules** | 15 days | ~5GB | Pre-aggregated |
| **Grafana dashboards** | Unlimited | ~100MB | N/A |
| **Alert history** | 120 hours | ~1GB | N/A |

### Storage Calculation

```
Metrics per second: ~500 samples
Storage per sample: ~2 bytes (compressed)
Daily storage: 500 * 86400 * 2 / 1024^3 ≈ 0.08GB
15-day storage: 0.08 * 15 ≈ 1.2GB

With overhead and indices: ~50GB total
```

## Network Architecture

### Service Discovery

```yaml
# Static configuration (development)
static_configs:
  - targets: ['localhost:9093', 'localhost:3000']

# Docker DNS (production)
static_configs:
  - targets: ['listings:9093', 'backend:3000']
```

### Relabeling

All targets get labeled with:
- `job`: Service name
- `instance`: Host:port
- `service`: Logical service name
- `service_type`: microservice/monolith
- `layer`: application/database/cache/monitoring

Example:
```yaml
labels:
  job: "listings-microservice"
  instance: "listings:9093"
  service: "listings"
  service_type: "microservice"
  layer: "application"
```

## Security Considerations

### Current State (Development)

- No authentication required
- HTTP only (no TLS)
- All ports exposed on localhost
- Default passwords (Grafana)

### Production Hardening Required

1. **Authentication:**
   - Enable Prometheus basic auth
   - Configure OAuth for Grafana
   - Restrict Alertmanager access

2. **Encryption:**
   - TLS for all HTTP endpoints
   - Certificate management (Let's Encrypt)

3. **Network:**
   - Internal network only
   - Reverse proxy (Nginx/Traefik)
   - Firewall rules

4. **Secrets:**
   - Move credentials to env vars
   - Use Docker secrets
   - Rotate API keys regularly

5. **Access Control:**
   - RBAC in Grafana
   - Service accounts
   - Audit logging

## Resource Requirements

### Minimum (Development)

- **CPU**: 2 cores
- **RAM**: 4GB
- **Disk**: 50GB
- **Network**: 10 Mbps

### Recommended (Production)

- **CPU**: 4 cores
- **RAM**: 8GB
- **Disk**: 100GB SSD
- **Network**: 100 Mbps

### Per Service

| Service | CPU | RAM | Disk |
|---------|-----|-----|------|
| Prometheus | 1-2 cores | 2-4GB | 50GB |
| Grafana | 0.5 cores | 512MB-1GB | 1GB |
| Alertmanager | 0.25 cores | 256MB | 1GB |
| Exporters | 0.25 cores | 128MB | - |

## Scalability

### Current Limits

- **Metrics/second**: 1,000
- **Active series**: 100,000
- **Queries/second**: 100
- **Concurrent dashboards**: 50

### Scaling Options

**Vertical:**
- Increase CPU/RAM for Prometheus
- Add SSDs for better I/O

**Horizontal:**
- Prometheus federation (multiple instances)
- Thanos for long-term storage
- Cortex for multi-tenancy
- Victoria Metrics for higher scale

## Integration Points

### Current

- Listings microservice (:9093/metrics)
- Monolith backend (:3000/metrics)
- PostgreSQL database
- Redis cache

### Future

- OpenSearch cluster
- Message queues (RabbitMQ/Kafka)
- Object storage (MinIO)
- CDN metrics
- Load balancer metrics

## Operational Runbooks

Quick reference for common operations:

| Task | Command | Frequency |
|------|---------|-----------|
| **Validate config** | `make validate` | Before deploy |
| **Start stack** | `make start` | - |
| **Check status** | `make status` | Daily |
| **View logs** | `make logs` | As needed |
| **Test alerts** | `make test-alerts` | Weekly |
| **Reload config** | `make reload-prometheus` | After changes |
| **Backup data** | `make backup` | Weekly |
| **Update images** | `make update` | Monthly |
| **Clean data** | `make clean` | Rarely |

## Files Overview

```
prometheus/
├── prometheus.yml              # Main config (scraping)
├── alerts.yml                  # Alert rules (20 alerts)
├── recording_rules.yml         # Recording rules (48 rules)
├── alertmanager.yml           # Alert routing
├── blackbox.yml               # HTTP probing config
├── postgres-exporter-queries.yml  # Custom DB queries
├── docker-compose.yml         # Stack definition
├── .env.example               # Environment template
├── .gitignore                # Git exclusions
├── Makefile                  # Convenience commands
│
├── validate-config.sh        # Config validator
├── test-alerts.sh           # Alert tester
│
├── README.md                # Full documentation
├── QUICK_START.md          # 5-minute setup
├── OVERVIEW.md             # This file
│
└── grafana/
    ├── provisioning/
    │   ├── datasources/    # Auto-configure Prometheus
    │   └── dashboards/     # Auto-import dashboards
    └── dashboards/         # Custom dashboards
```

## Next Steps

1. **Quick Start**: Follow [QUICK_START.md](QUICK_START.md)
2. **Full Setup**: Read [README.md](README.md)
3. **Validate**: Run `./validate-config.sh`
4. **Deploy**: Run `make start-dev`
5. **Test**: Run `make test-alerts`
6. **Monitor**: Open Grafana at http://localhost:3030

## Support

- **Documentation**: [README.md](README.md)
- **Issues**: Check logs with `make logs`
- **Testing**: Use `test-alerts.sh`
- **Team**: platform-team@vondi.rs
