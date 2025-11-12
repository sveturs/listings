# Grafana Dashboards - Quick Start Guide

## What's Included

✅ **4 Production-Ready Dashboards** (57 total panels)
- `service-health.json` - 10 panels for real-time service monitoring
- `infrastructure.json` - 14 panels for resource and infrastructure
- `business-metrics.json` - 17 panels for business operations
- `alerting-slo.json` - 16 panels for SLO tracking and alerts

✅ **67+ Metrics** from listings microservice
✅ **Complete Documentation** (README.md with examples and alert rules)
✅ **Import/Validation Scripts** for easy deployment

---

## Quick Import (5 minutes)

### Option 1: Using Import Script (Recommended)

```bash
# 1. Navigate to directory
cd /p/github.com/sveturs/listings/deployment/grafana

# 2. Import all dashboards
./import-dashboards.sh http://your-grafana:3000 YOUR_API_KEY

# Or without authentication (if Grafana allows)
./import-dashboards.sh http://localhost:3000
```

### Option 2: Manual Import via UI

1. Open Grafana → **Dashboards** → **Import**
2. Upload each JSON file from `dashboards/` directory
3. Select **Prometheus** as datasource
4. Click **Import**

### Option 3: Kubernetes ConfigMap

```bash
kubectl create configmap listings-dashboards \
  --from-file=dashboards/ \
  -n monitoring
```

---

## Verification

Check that metrics are being exported:

```bash
# Test metrics endpoint
curl http://localhost:9093/metrics | grep listings_

# Expected output:
# listings_grpc_requests_total{...}
# listings_http_requests_total{...}
# listings_db_connections_open
# ...
```

---

## Dashboard Overview

### 🟢 Service Health (`service-health.json`)
**Use for**: Day-to-day operations, incident response
- Uptime %, RPS, Error Rate
- Latency P50/P95/P99
- Request volume and active requests

### 🔧 Infrastructure (`infrastructure.json`)
**Use for**: Capacity planning, resource optimization
- CPU, Memory, Goroutines
- Database connection pool
- Redis cache performance
- Rate limiting activity

### 📊 Business Metrics (`business-metrics.json`)
**Use for**: Product analytics, feature tracking
- Listing operations (create/update/delete)
- Inventory management and stock alerts
- Product views and search analytics
- Worker queue health

### 🚨 Alerting & SLO (`alerting-slo.json`)
**Use for**: Reliability engineering, SLO compliance
- Service availability (99.9% target)
- Success rate (99.5% target)
- Latency SLO (P95 < 200ms, P99 < 1s)
- Error budget tracking
- Active alerts and incident history

---

## Common Queries (Copy-Paste Ready)

### Current RPS
```promql
sum(rate(listings_grpc_requests_total[1m]))
```

### Error Rate
```promql
(sum(rate(listings_grpc_requests_total{status!="OK"}[5m]))
/ sum(rate(listings_grpc_requests_total[5m]))) * 100
```

### P95 Latency
```promql
histogram_quantile(0.95,
  sum(rate(listings_grpc_request_duration_seconds_bucket[5m])) by (le))
```

### Service Uptime (24h)
```promql
(1 - (sum(rate(listings_errors_total[24h]))
/ sum(rate(listings_grpc_requests_total[24h])))) * 100
```

### Cache Hit Ratio
```promql
(sum(rate(listings_cache_hits_total[5m]))
/ (sum(rate(listings_cache_hits_total[5m]))
+ sum(rate(listings_cache_misses_total[5m])))) * 100
```

---

## Alert Rules (Quick Setup)

Create `/etc/prometheus/rules/listings.yml`:

```yaml
groups:
  - name: listings_critical
    interval: 30s
    rules:
      - alert: ListingsServiceDown
        expr: up{job="listings"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Listings service is down"

      - alert: HighErrorRate
        expr: |
          (sum(rate(listings_grpc_requests_total{status!="OK"}[5m]))
          / sum(rate(listings_grpc_requests_total[5m]))) * 100 > 5
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate: {{ $value }}%"

      - alert: HighLatencyP95
        expr: |
          histogram_quantile(0.95,
          sum(rate(listings_grpc_request_duration_seconds_bucket[5m])) by (le)) > 0.2
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "P95 latency exceeds 200ms"
```

Reload Prometheus:
```bash
curl -X POST http://localhost:9090/-/reload
```

---

## Troubleshooting

### "No data" in dashboards

1. **Check service is running**:
   ```bash
   netstat -tlnp | grep 9093
   ```

2. **Verify metrics endpoint**:
   ```bash
   curl http://localhost:9093/metrics | head -20
   ```

3. **Check Prometheus scrape config** (`/etc/prometheus/prometheus.yml`):
   ```yaml
   scrape_configs:
     - job_name: 'listings'
       static_configs:
         - targets: ['localhost:9093']
   ```

4. **Verify Prometheus is scraping**:
   - Open http://localhost:9090/targets
   - Check `listings` job is UP

### Dashboards show partial data

- Increase time range (top-right corner)
- Check metric namespace matches: `listings_*`
- Verify label filters match your setup

---

## Next Steps

1. ✅ Import dashboards
2. ✅ Verify metrics are visible
3. ✅ Set up Prometheus alert rules
4. ✅ Configure alert routing (Alertmanager, PagerDuty, Slack)
5. ✅ Customize thresholds for your capacity
6. ✅ Share dashboard URLs with team

---

## Files in This Directory

```
deployment/grafana/
├── dashboards/
│   ├── service-health.json       (10 panels)
│   ├── infrastructure.json       (14 panels)
│   ├── business-metrics.json     (17 panels)
│   ├── alerting-slo.json         (16 panels)
│   └── README.md                 (Full documentation)
├── import-dashboards.sh          (Automated import)
├── validate-dashboards.sh        (JSON validation)
└── QUICKSTART.md                 (This file)
```

---

## Support & Documentation

- **Full README**: `dashboards/README.md` (detailed metrics, SLOs, examples)
- **Metrics Code**: `/internal/metrics/metrics.go`
- **Prometheus Config**: `/deployment/prometheus/prometheus.yml`
- **Alert Rules**: `/deployment/prometheus/rules/`

---

## Tips

💡 **Start with Service Health dashboard** - It's your operational command center
💡 **Set up alerts ASAP** - Copy alert rules from README.md
💡 **Bookmark dashboards** - Add to your browser favorites
💡 **Use annotations** - Mark deployments and incidents
💡 **Weekly SLO review** - Check Alerting & SLO dashboard every Monday
💡 **Customize thresholds** - Adjust gauges based on your capacity

---

**Dashboard Version**: 1.0.0
**Created**: 2025-01-05
**Metrics**: 67+ from listings microservice
**Compatibility**: Grafana 9.x+, Prometheus 2.x+
