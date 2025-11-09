# Listings Microservice Monitoring Guide

**Last Updated:** 2025-11-05
**Version:** 1.0.0
**Maintainer:** Platform Team

## Table of Contents

- [Overview](#overview)
- [Monitoring Stack](#monitoring-stack)
- [Grafana Dashboards](#grafana-dashboards)
- [Prometheus Metrics](#prometheus-metrics)
- [Alert Definitions](#alert-definitions)
- [Log Analysis](#log-analysis)
- [Metric Interpretation](#metric-interpretation)
- [Creating Custom Dashboards](#creating-custom-dashboards)
- [Best Practices](#best-practices)

---

## Overview

### Monitoring Philosophy

**The Three Pillars of Observability:**
1. **Metrics** - Numerical measurements over time (Prometheus)
2. **Logs** - Event records with context (Journald + structured logging)
3. **Traces** - Request flow through system (Jaeger - optional)

**Our Approach:**
```
Metrics → WHAT is happening (service is slow)
Logs → WHY it's happening (database query timeout)
Traces → WHERE in the flow (which service/function)
```

### Access Information

| Tool | URL | Purpose | Access Level |
|------|-----|---------|--------------|
| **Grafana** | https://grafana.svetu.rs | Dashboards and visualization | All engineers |
| **Prometheus** | http://prometheus.svetu.rs:9090 | Metrics and queries | All engineers |
| **AlertManager** | http://alertmanager.svetu.rs:9093 | Alert management | On-call engineers |
| **Jaeger** | http://jaeger.svetu.rs:16686 | Distributed tracing | All engineers |

**Credentials:** Stored in company password manager

---

## Monitoring Stack

### Architecture

```
┌─────────────────────┐
│ Listings Service    │
│  :8086/metrics      │──┐
└─────────────────────┘  │
                         │
┌─────────────────────┐  │  ┌──────────────────┐
│ PostgreSQL Exporter │──┼─→│   Prometheus     │
│  :9187/metrics      │  │  │   (Storage +     │
└─────────────────────┘  │  │    Queries)      │
                         │  └────────┬─────────┘
┌─────────────────────┐  │           │
│ Redis Exporter      │──┘           │
│  :9121/metrics      │              │
└─────────────────────┘              │
                                     ↓
                            ┌─────────────────┐
                            │    Grafana      │
                            │  (Visualization)│
                            └─────────────────┘
                                     ↑
                                     │
                            ┌─────────────────┐
                            │  AlertManager   │
                            │  (Alerting)     │
                            └─────────────────┘
                                     ↓
                            ┌─────────────────┐
                            │   PagerDuty     │
                            │  Slack, Email   │
                            └─────────────────┘
```

### Components

**Prometheus:**
- Metrics collection (15s scrape interval)
- Time-series database
- Query language (PromQL)
- Alert evaluation

**Grafana:**
- Dashboard visualization
- Multi-datasource support
- Alerting (supplementary)
- User management

**Exporters:**
- Listings service: Native `/metrics` endpoint
- PostgreSQL: `postgres_exporter`
- Redis: `redis_exporter`
- Node: `node_exporter` (system metrics)

---

## Grafana Dashboards

### Primary Dashboards

#### 1. Listings Service Overview

**Location:** Dashboards → Listings → Overview
**URL:** https://grafana.svetu.rs/d/listings-overview

**Purpose:** High-level service health and SLO tracking

**Panels:**
- **Service Status** - Health check status
- **Request Rate** - Requests per second (by method)
- **Error Rate** - % of failed requests
- **Latency** - P50, P95, P99 latencies
- **SLO Compliance** - Availability vs target
- **Error Budget** - Remaining error budget
- **Top Errors** - Most common errors
- **Active Alerts** - Current firing alerts

**When to Use:**
- Daily health check
- Incident investigation (first stop)
- SLO review
- Executive reporting

**Key Metrics to Watch:**
```
✅ Green: All metrics within SLO
⚠️ Yellow: Approaching thresholds
🔴 Red: SLO breach or critical issue
```

---

#### 2. Listings Service Details

**Location:** Dashboards → Listings → Details
**URL:** https://grafana.svetu.rs/d/listings-details

**Purpose:** Deep-dive into service internals

**Panels:**

**Request Metrics:**
- Requests by endpoint
- Request duration histogram
- Active requests (in-flight)
- Request size distribution

**Error Metrics:**
- Errors by type
- Errors by endpoint
- Error rate trend (7d)
- Top error messages

**Timeout Metrics:**
- Timeouts by endpoint
- Near-timeouts (>80% of limit)
- Timeout duration distribution

**Rate Limiting:**
- Rate limit hits
- Rate limit rejections
- Rejection rate by endpoint
- Top rate-limited IPs

**When to Use:**
- Performance investigation
- Capacity planning
- Rate limit tuning
- Debugging specific endpoints

---

#### 3. Database Performance

**Location:** Dashboards → Listings → Database
**URL:** https://grafana.svetu.rs/d/listings-database

**Purpose:** PostgreSQL monitoring

**Panels:**

**Connection Pool:**
- Open connections
- Idle connections
- Active connections
- Connection wait time

**Query Performance:**
- Query duration (by operation)
- Slow queries (>1s)
- Transactions per second
- Rows returned per query

**Database Health:**
- Table sizes
- Index usage
- Cache hit ratio
- Replication lag (if applicable)

**Locks and Blocking:**
- Lock wait time
- Deadlocks
- Blocking queries

**When to Use:**
- Database slow query investigation
- Connection pool tuning
- Capacity planning
- Query optimization

**Key Queries:**

```promql
# Connection pool utilization
(listings_db_connections_open - listings_db_connections_idle) /
listings_db_connections_open * 100

# Cache hit ratio (should be >95%)
sum(rate(pg_stat_database_blks_hit[5m])) /
(sum(rate(pg_stat_database_blks_hit[5m])) +
 sum(rate(pg_stat_database_blks_read[5m]))) * 100

# Slow queries (>1s)
sum(rate(listings_db_query_duration_seconds_bucket{le="1.0"}[5m]))
```

---

#### 4. Redis Performance

**Location:** Dashboards → Listings → Redis
**URL:** https://grafana.svetu.rs/d/listings-redis

**Purpose:** Cache and rate limiter monitoring

**Panels:**

**Connection Metrics:**
- Connected clients
- Blocked clients
- Client longest output list

**Memory:**
- Used memory
- Memory fragmentation ratio
- Evicted keys
- Expired keys

**Performance:**
- Commands per second
- Hit rate
- Keyspace hits/misses
- Network I/O

**Cache Effectiveness:**
- Cache hit rate (listings service)
- Cache TTL distribution
- Top cache keys by size

**When to Use:**
- Cache performance tuning
- Memory issues investigation
- Rate limiter troubleshooting
- Cache eviction analysis

---

#### 5. SLO Dashboard

**Location:** Dashboards → Listings → SLO
**URL:** https://grafana.svetu.rs/d/listings-slo

**Purpose:** SLO tracking and error budget management

**Panels:**

**Availability:**
- Current availability (30d rolling)
- SLO target line (99.9%)
- Error budget remaining
- Error budget burn rate

**Latency:**
- P95 latency trend
- P99 latency trend
- SLO target lines
- Latency by endpoint

**Error Rate:**
- Current error rate
- Error rate trend
- Error budget impact
- Errors by category

**Incident Impact:**
- Downtime per incident
- Cumulative downtime
- Projected month-end status

**When to Use:**
- Monthly SLO reviews
- Error budget tracking
- Incident impact assessment
- Reporting to leadership

---

### Dashboard Navigation Tips

**Time Range Selection:**
```
Last 5 minutes   → Real-time incident investigation
Last 1 hour      → Recent performance analysis
Last 24 hours    → Daily patterns
Last 7 days      → Weekly trends
Last 30 days     → SLO tracking, monthly review
Custom range     → Historical analysis
```

**Refresh Rate:**
```
Off             → Historical analysis
5s              → Active incident response
30s             → Real-time monitoring
1m              → Dashboard displays
```

**Variables (Dropdown Filters):**
- **Method:** Filter by specific gRPC method
- **Status:** Filter by response status code
- **Instance:** Filter by service instance (if multiple)

**Zoom and Pan:**
- Click and drag to zoom
- Shift + drag to pan
- Double-click to reset zoom

**Panel Actions:**
- **View** → See raw query and data
- **Edit** → Modify panel (requires edit permission)
- **Share** → Get shareable link
- **Explore** → Open in Prometheus query interface

---

## Prometheus Metrics

### Metric Categories

#### 1. gRPC Handler Metrics

**Request Counters:**
```promql
# Total requests by method and status
listings_grpc_requests_total{method="/listings.v1.ListingsService/GetListing", status="200"}

# Requests per second
rate(listings_grpc_requests_total[5m])

# Requests by status code family
sum by (status) (rate(listings_grpc_requests_total[5m]))
```

**Request Duration:**
```promql
# P50 latency
histogram_quantile(0.50, rate(listings_grpc_request_duration_seconds_bucket[5m]))

# P95 latency
histogram_quantile(0.95, rate(listings_grpc_request_duration_seconds_bucket[5m]))

# P99 latency
histogram_quantile(0.99, rate(listings_grpc_request_duration_seconds_bucket[5m]))

# Average latency by method
rate(listings_grpc_request_duration_seconds_sum[5m]) /
rate(listings_grpc_request_duration_seconds_count[5m])
```

**Active Requests:**
```promql
# Current in-flight requests
listings_grpc_handler_requests_active

# By method
listings_grpc_handler_requests_active{method="/listings.v1.ListingsService/GetListing"}
```

---

#### 2. Database Metrics

**Connection Pool:**
```promql
# Open connections
listings_db_connections_open

# Idle connections
listings_db_connections_idle

# Active connections
listings_db_connections_open - listings_db_connections_idle

# Pool utilization %
(listings_db_connections_open - listings_db_connections_idle) /
listings_db_connections_open * 100
```

**Query Performance:**
```promql
# Query duration by operation
rate(listings_db_query_duration_seconds_sum{operation="SELECT"}[5m]) /
rate(listings_db_query_duration_seconds_count{operation="SELECT"}[5m])

# Slow query rate (>1s)
sum(rate(listings_db_query_duration_seconds_bucket{le="1.0", operation="SELECT"}[5m]))
```

---

#### 3. Rate Limiting Metrics

```promql
# Rate limit evaluations
rate(listings_rate_limit_hits_total[5m])

# Rejection rate
rate(listings_rate_limit_rejected_total[5m]) /
rate(listings_rate_limit_hits_total[5m]) * 100

# Top rate-limited methods
topk(5, sum by (method) (rate(listings_rate_limit_rejected_total[5m])))
```

---

#### 4. Timeout Metrics

```promql
# Timeouts by method
sum by (method) (rate(listings_timeouts_total[5m]))

# Near-timeout warnings (>80% of limit)
sum by (method) (rate(listings_near_timeouts_total[5m]))

# Timeout duration histogram
histogram_quantile(0.99, rate(listings_timeout_duration_seconds_bucket[5m]))
```

---

#### 5. Business Metrics

**Inventory Operations:**
```promql
# Product views
rate(listings_inventory_product_views_total[5m])

# Stock updates
rate(listings_inventory_stock_operations_total{operation="update"}[5m])

# Inventory movements
rate(listings_inventory_movements_recorded_total[5m])

# Out of stock products
listings_inventory_out_of_stock_products
```

**Listing Operations:**
```promql
# Listings created
rate(listings_listings_created_total[5m])

# Listings updated
rate(listings_listings_updated_total[5m])

# Listings deleted
rate(listings_listings_deleted_total[5m])

# Search queries
rate(listings_listings_searched_total[5m])
```

---

#### 6. Cache Metrics

```promql
# Cache hit rate
rate(listings_cache_hits_total[5m]) /
(rate(listings_cache_hits_total[5m]) + rate(listings_cache_misses_total[5m])) * 100

# Cache operations by type
sum by (cache_type) (rate(listings_cache_hits_total[5m]))
```

---

#### 7. Worker Metrics

```promql
# Indexing queue size
listings_indexing_queue_size

# Jobs processed (success/failure)
sum by (status) (rate(listings_indexing_jobs_processed_total[5m]))

# Job duration
rate(listings_indexing_job_duration_seconds_sum[5m]) /
rate(listings_indexing_job_duration_seconds_count[5m])
```

---

### Useful PromQL Queries

**Error Rate Calculation:**
```promql
# Overall error rate
sum(rate(listings_grpc_requests_total{status=~"5.."}[5m])) /
sum(rate(listings_grpc_requests_total[5m])) * 100

# Error rate by method
sum by (method) (rate(listings_grpc_requests_total{status=~"5.."}[5m])) /
sum by (method) (rate(listings_grpc_requests_total[5m])) * 100
```

**Top Endpoints by Traffic:**
```promql
topk(10, sum by (method) (rate(listings_grpc_requests_total[5m])))
```

**Top Endpoints by Latency:**
```promql
topk(10, sum by (method) (
  rate(listings_grpc_request_duration_seconds_sum[5m]) /
  rate(listings_grpc_request_duration_seconds_count[5m])
))
```

**Requests vs Capacity:**
```promql
# Current RPS
sum(rate(listings_grpc_requests_total[5m]))

# vs historical max RPS (last 7 days)
max_over_time(sum(rate(listings_grpc_requests_total[5m]))[7d:5m])
```

---

## Alert Definitions

### Alert Configuration

**Location:** `/p/github.com/sveturs/listings/deployment/prometheus/alerts.yml`

### Critical Alerts (P1)

#### ListingsServiceDown

```yaml
- alert: ListingsServiceDown
  expr: up{job="listings"} == 0
  for: 1m
  labels:
    severity: critical
    service: listings
  annotations:
    summary: "Listings service is down"
    description: "Listings service has been unavailable for > 1 minute"
    runbook: "docs/operations/RUNBOOK.md#3-service-down"
```

**Meaning:** Service process not running or not responding to Prometheus scrapes.

**Action:** Immediate investigation and restart. See RUNBOOK.md → Service Down.

---

#### ListingsCriticalErrorRate

```yaml
- alert: ListingsCriticalErrorRate
  expr: |
    (sum(rate(listings_grpc_requests_total{status=~"5.."}[5m])) /
     sum(rate(listings_grpc_requests_total[5m]))) * 100 > 10
  for: 5m
  labels:
    severity: critical
    service: listings
  annotations:
    summary: "Critical error rate: {{ $value | humanize }}%"
    description: "Error rate above 10% for 5 minutes"
```

**Meaning:** More than 10% of requests failing.

**Action:** Immediate investigation. Likely database or dependency issue.

---

### High Severity Alerts (P2)

#### ListingsHighErrorRate

```yaml
- alert: ListingsHighErrorRate
  expr: |
    (sum(rate(listings_grpc_requests_total{status=~"5.."}[5m])) /
     sum(rate(listings_grpc_requests_total[5m]))) * 100 > 1
  for: 10m
  labels:
    severity: high
    service: listings
  annotations:
    summary: "High error rate: {{ $value | humanize }}%"
    description: "Error rate above 1% for 10 minutes (SLO breach)"
    runbook: "docs/operations/RUNBOOK.md#1-high-error-rate-1"
```

---

#### ListingsHighLatency

```yaml
- alert: ListingsHighLatency
  expr: |
    histogram_quantile(0.99,
      rate(listings_grpc_request_duration_seconds_bucket[5m])
    ) > 2
  for: 10m
  labels:
    severity: high
    service: listings
  annotations:
    summary: "High latency: P99 = {{ $value | humanize }}s"
    description: "P99 latency above 2s for 10 minutes (SLO breach)"
    runbook: "docs/operations/RUNBOOK.md#2-high-latency-p99-2s"
```

---

#### ListingsDBPoolExhausted

```yaml
- alert: ListingsDBPoolExhausted
  expr: |
    (listings_db_connections_open - listings_db_connections_idle) /
    listings_db_connections_open * 100 > 90
  for: 5m
  labels:
    severity: high
    service: listings
  annotations:
    summary: "Database connection pool at {{ $value | humanize }}%"
    description: "Connection pool utilization above 90%"
    runbook: "docs/operations/RUNBOOK.md#4-database-connection-pool-exhausted"
```

---

### Medium Severity Alerts (P3)

#### ListingsElevatedLatency

```yaml
- alert: ListingsElevatedLatency
  expr: |
    histogram_quantile(0.95,
      rate(listings_grpc_request_duration_seconds_bucket[5m])
    ) > 1
  for: 15m
  labels:
    severity: medium
    service: listings
  annotations:
    summary: "Elevated latency: P95 = {{ $value | humanize }}s"
    description: "P95 latency above 1s for 15 minutes"
```

---

#### ListingsRedisDown

```yaml
- alert: ListingsRedisDown
  expr: redis_up{job="redis-listings"} == 0
  for: 5m
  labels:
    severity: medium
    service: listings
  annotations:
    summary: "Redis is down"
    description: "Redis unavailable - rate limiting and caching disabled"
    runbook: "docs/operations/RUNBOOK.md#5-redis-connection-issues"
```

---

#### ListingsHighMemory

```yaml
- alert: ListingsHighMemory
  expr: |
    process_resident_memory_bytes{job="listings"} /
    1024 / 1024 / 1024 > 4
  for: 10m
  labels:
    severity: medium
    service: listings
  annotations:
    summary: "High memory usage: {{ $value | humanize }}GB"
    description: "Memory usage above 4GB for 10 minutes"
    runbook: "docs/operations/RUNBOOK.md#7-memory-leak--oom-killed"
```

---

### SLO Alerts

#### SLOAvailabilityAtRisk

```yaml
- alert: SLOAvailabilityAtRisk
  expr: |
    (sum(rate(listings_grpc_requests_total{status!~"5.."}[30d])) /
     sum(rate(listings_grpc_requests_total[30d]))) * 100 < 99.95
  labels:
    severity: high
    service: listings
    slo: availability
  annotations:
    summary: "SLO at risk: {{ $value | humanize }}%"
    description: "30-day availability below 99.95% (buffer zone)"
```

---

## Log Analysis

### Structured Logging Format

**All logs are JSON-formatted:**
```json
{
  "timestamp": "2025-11-05T14:30:00Z",
  "level": "error",
  "message": "Database query failed",
  "component": "repository",
  "operation": "GetListing",
  "listing_id": 328,
  "error": "pq: connection refused",
  "duration_ms": 5000,
  "user_id": "user-123",
  "request_id": "req-abc-456"
}
```

### Log Query Examples

**Find errors in last hour:**
```bash
sudo journalctl -u listings-service --since "1 hour ago" | \
  grep '"level":"error"' | \
  jq -r '"\(.timestamp) [\(.component)] \(.message)"'
```

**Count errors by type:**
```bash
sudo journalctl -u listings-service --since "1 hour ago" | \
  grep '"level":"error"' | \
  jq -r '.error' | \
  sort | uniq -c | sort -rn
```

**Find slow queries:**
```bash
sudo journalctl -u listings-service --since "30 minutes ago" | \
  jq -r 'select(.duration_ms > 1000) | "\(.duration_ms)ms - \(.operation) - \(.query)"'
```

**Track specific request:**
```bash
REQUEST_ID="req-abc-456"
sudo journalctl -u listings-service | \
  jq -r "select(.request_id == \"$REQUEST_ID\")"
```

**Find requests by user:**
```bash
USER_ID="user-123"
sudo journalctl -u listings-service --since "1 day ago" | \
  jq -r "select(.user_id == \"$USER_ID\") | \"\(.timestamp) \(.operation) \(.duration_ms)ms\""
```

---

## Metric Interpretation

### Understanding Histograms

**Histogram metrics have three components:**
1. `_bucket{le="X"}` - Count of observations ≤ X
2. `_sum` - Sum of all observations
3. `_count` - Total count of observations

**Example:**
```promql
listings_grpc_request_duration_seconds_bucket{le="0.5"} 1000  # 1000 requests ≤ 500ms
listings_grpc_request_duration_seconds_bucket{le="1.0"} 1800  # 1800 requests ≤ 1s
listings_grpc_request_duration_seconds_sum 1500              # Total 1500s
listings_grpc_request_duration_seconds_count 2000            # 2000 requests

Average latency = 1500 / 2000 = 0.75s
```

### Understanding Percentiles

**What percentiles mean:**
- **P50 (median):** 50% of requests faster than this
- **P95:** 95% of requests faster than this
- **P99:** 99% of requests faster than this

**Why P99 matters:**
```
If P50 = 100ms and P99 = 5s:
- Most users have good experience (100ms)
- But 1% have terrible experience (5s)
- At 1000 req/s, that's 10 unhappy users per second!
```

**Rule of thumb:**
- **P50** - Optimize for typical user
- **P95** - SLO target
- **P99** - Debugging outliers
- **P99.9** - Identifying edge cases

### Understanding Rate Functions

**`rate()` calculates per-second rate:**
```promql
# Raw counter
listings_grpc_requests_total = 10000

# Rate over 5 minutes
rate(listings_grpc_requests_total[5m]) = 33.3  # 33.3 requests/second
```

**`increase()` calculates total increase:**
```promql
# Total increase over 5 minutes
increase(listings_grpc_requests_total[5m]) = 10000  # 10k requests in 5 min
```

**Range selection best practices:**
- **[1m]** - Very short-term spikes
- **[5m]** - Standard for most alerts and dashboards
- **[1h]** - Hourly aggregations
- **[24h]** - Daily patterns
- **[30d]** - Monthly SLO tracking

---

## Creating Custom Dashboards

### Dashboard Creation Workflow

**1. Plan Your Dashboard:**
- What question am I answering?
- Who is the audience?
- What metrics are needed?
- What time ranges?

**2. Create Dashboard:**
```
Grafana → + → Dashboard → Add new panel
```

**3. Configure Panel:**
- **Query:** Write PromQL query
- **Visualization:** Choose graph type
- **Legend:** Configure labels
- **Thresholds:** Set warning/error levels
- **Unit:** Select appropriate unit (seconds, bytes, etc.)

**4. Save and Share:**
- Name descriptively
- Add to appropriate folder
- Set permissions
- Document in this guide

### Example Custom Panels

**Request Rate by Endpoint:**
```promql
Query: sum by (method) (rate(listings_grpc_requests_total[5m]))
Visualization: Time series
Legend: {{method}}
Unit: requests/sec
```

**Error Rate Gauge:**
```promql
Query: (sum(rate(listings_grpc_requests_total{status=~"5.."}[5m])) /
        sum(rate(listings_grpc_requests_total[5m]))) * 100
Visualization: Gauge
Unit: percent (0-100)
Thresholds: 0-1% green, 1-5% yellow, 5-100% red
```

**Database Connection Heatmap:**
```promql
Query: listings_db_connections_open - listings_db_connections_idle
Visualization: Heatmap
Time range: Last 24 hours
Bucket data: 5 minute intervals
```

---

## Best Practices

### Dashboard Design

**Do's ✅**
- Use consistent time ranges across panels
- Group related metrics together
- Use appropriate visualizations (line for trends, gauge for current value)
- Set meaningful thresholds and colors
- Include SLO target lines
- Add panel descriptions
- Use variables for filtering

**Don'ts ❌**
- Don't overcrowd dashboards (max 12 panels per page)
- Don't use fancy visualizations without purpose
- Don't set artificial thresholds
- Don't forget units
- Don't mix incompatible metrics on same axis

### Alert Design

**Good Alerts:**
- Actionable (responder knows what to do)
- Meaningful (indicates real problem)
- Contextual (includes runbook link)
- Tuned (low false positive rate)

**Bad Alerts:**
- "Server CPU high" (not actionable - which server? why?)
- Fires every day at same time (predictable, tune it!)
- No context (what should I do?)
- Duplicate alerts (consolidate!)

### Query Optimization

**Efficient Queries:**
```promql
# Good: Specific labels
sum(rate(listings_grpc_requests_total{job="listings", status="200"}[5m]))

# Bad: Too broad, slow
sum(rate(listings_grpc_requests_total[5m]))
```

**Use recording rules for expensive queries:**
```yaml
# Pre-calculate frequently used metrics
- record: listings:request_rate:5m
  expr: sum(rate(listings_grpc_requests_total[5m]))
```

---

## Monitoring Checklist

### Daily

- [ ] Check "Listings Overview" dashboard
- [ ] Verify all panels showing data
- [ ] Check for active alerts
- [ ] Review error rate (should be <1%)
- [ ] Check SLO compliance

### Weekly

- [ ] Review "Listings Details" dashboard
- [ ] Analyze performance trends
- [ ] Check database performance
- [ ] Review Redis metrics
- [ ] Update dashboard documentation

### Monthly

- [ ] Full SLO review with dashboard data
- [ ] Alert tuning (reduce false positives)
- [ ] Dashboard cleanup (remove unused)
- [ ] Create new panels for new features
- [ ] Training for new team members

---

## Troubleshooting Monitoring

### Dashboard Not Loading

**Symptoms:** Blank panels, "No data" errors

**Checks:**
```bash
# Check Prometheus
curl http://prometheus.svetu.rs:9090/-/healthy

# Check Grafana datasource
Grafana → Configuration → Data Sources → Prometheus → Test

# Check if metrics exist
curl http://localhost:8086/metrics | grep listings_
```

### Missing Metrics

**Symptoms:** Panel shows "No data" for specific metric

**Checks:**
```bash
# Check if service is exporting metric
curl -s http://localhost:8086/metrics | grep <metric_name>

# Check Prometheus targets
curl http://prometheus.svetu.rs:9090/api/v1/targets | jq '.data.activeTargets[] | select(.labels.job == "listings")'

# Check for metric typo in query
```

### Alerts Not Firing

**Symptoms:** Expected alert not received

**Checks:**
```bash
# Check alert rule
curl http://prometheus.svetu.rs:9090/api/v1/rules | jq '.data.groups[].rules[] | select(.name == "ListingsHighErrorRate")'

# Check alert state
curl http://prometheus.svetu.rs:9090/api/v1/alerts | jq '.data.alerts[] | select(.labels.alertname == "ListingsHighErrorRate")'

# Check AlertManager
curl http://alertmanager.svetu.rs:9093/api/v2/alerts
```

---

## Resources

### Documentation
- **Prometheus Docs:** https://prometheus.io/docs/
- **PromQL Tutorial:** https://prometheus.io/docs/prometheus/latest/querying/basics/
- **Grafana Docs:** https://grafana.com/docs/
- **Best Practices:** https://prometheus.io/docs/practices/

### Internal
- **RUNBOOK.md:** Incident response procedures
- **TROUBLESHOOTING.md:** Debugging guide
- **SLO_GUIDE.md:** SLO management

### Training
- **Prometheus Basics:** https://training.promlabs.com/
- **Grafana Fundamentals:** https://grafana.com/tutorials/

---

**Document Version:** 1.0.0
**Last Reviewed:** 2025-11-05
**Next Review:** 2025-12-05
**Owner:** Platform Team

**Remember: Monitoring is not a goal, it's a tool. The goal is reliable service. Monitor what matters, alert on what's actionable, and iterate continuously.**
