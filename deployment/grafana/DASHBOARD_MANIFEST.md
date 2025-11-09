# Grafana Dashboard Manifest

## Dashboard Inventory

| Dashboard | UID | Panels | Purpose | Priority |
|-----------|-----|--------|---------|----------|
| Service Health | `listings-service-health` | 10 | Real-time monitoring | 🔴 Critical |
| Infrastructure | `listings-infrastructure` | 14 | Resource monitoring | 🟡 High |
| Business Metrics | `listings-business-metrics` | 17 | Analytics | 🟢 Medium |
| Alerting & SLO | `listings-alerting-slo` | 16 | Reliability | 🔴 Critical |

**Total**: 57 panels across 4 dashboards

---

## Metrics Coverage

### Service Layer (12 metrics)
- ✅ `grpc_requests_total` - Request counter by method/status
- ✅ `grpc_request_duration_seconds` - Latency histogram
- ✅ `grpc_handler_requests_active` - In-flight requests
- ✅ `http_requests_total` - HTTP request counter
- ✅ `http_request_duration_seconds` - HTTP latency
- ✅ `http_requests_in_flight` - Active HTTP requests
- ✅ `rate_limit_hits_total` - Rate limit evaluations
- ✅ `rate_limit_allowed_total` - Allowed requests
- ✅ `rate_limit_rejected_total` - Rejected requests
- ✅ `timeouts_total` - Timeout counter
- ✅ `near_timeouts_total` - Near-timeout counter
- ✅ `timeout_duration_seconds` - Timeout duration histogram

### Business Layer (8 metrics)
- ✅ `listings_created_total` - Listings created
- ✅ `listings_updated_total` - Listings updated
- ✅ `listings_deleted_total` - Listings deleted
- ✅ `listings_searched_total` - Search queries
- ✅ `inventory_product_views_total` - Product view increments
- ✅ `inventory_stock_operations_total` - Stock operations
- ✅ `inventory_movements_recorded_total` - Movement tracking
- ✅ `inventory_stock_low_threshold_reached_total` - Low stock alerts

### Infrastructure Layer (9 metrics)
- ✅ `db_connections_open` - Open DB connections
- ✅ `db_connections_idle` - Idle DB connections
- ✅ `db_query_duration_seconds` - Query duration histogram
- ✅ `cache_hits_total` - Cache hits
- ✅ `cache_misses_total` - Cache misses
- ✅ `indexing_queue_size` - Queue size
- ✅ `indexing_jobs_processed_total` - Jobs processed
- ✅ `indexing_job_duration_seconds` - Job duration
- ✅ `errors_total` - Error counter by component/type

### Inventory Layer (5 metrics)
- ✅ `inventory_stock_value` - Current stock value
- ✅ `inventory_out_of_stock_products` - Out-of-stock count
- ✅ `inventory_product_views_errors_total` - View errors
- ✅ `inventory_movements_errors_total` - Movement errors
- ✅ `inventory_stock_low_threshold_reached_total` - Low stock alerts

### System Metrics (Go runtime, Process)
- ✅ `go_goroutines` - Active goroutines
- ✅ `go_memstats_*` - Memory statistics
- ✅ `process_cpu_seconds_total` - CPU usage
- ✅ `process_resident_memory_bytes` - Memory usage
- ✅ `process_open_fds` - Open file descriptors
- ✅ `process_start_time_seconds` - Process start time

**Total Metrics Tracked**: 67+

---

## SLO Definitions

### Availability SLO
- **Target**: 99.9% (30 days)
- **Error Budget**: 43.2 minutes/month
- **Measurement**: `(1 - error_rate) * 100`
- **Dashboard**: Alerting & SLO

### Latency SLO (P95)
- **Target**: < 200ms
- **Dashboard**: Service Health, Alerting & SLO
- **Measurement**: `histogram_quantile(0.95, grpc_request_duration_seconds_bucket)`

### Latency SLO (P99)
- **Target**: < 1000ms
- **Dashboard**: Service Health, Alerting & SLO
- **Measurement**: `histogram_quantile(0.99, grpc_request_duration_seconds_bucket)`

### Success Rate SLO
- **Target**: 99.5% (7 days)
- **Dashboard**: Alerting & SLO
- **Measurement**: `(status=OK / total_requests) * 100`

---

## Panel Types Used

| Type | Count | Usage |
|------|-------|-------|
| Timeseries | 22 | Trend visualization |
| Stat | 16 | Current values |
| Gauge | 8 | Threshold monitoring |
| Table | 4 | Detailed breakdowns |
| Histogram | 3 | Distribution analysis |
| Piechart | 1 | Composition |
| Barchart | 1 | Comparison |
| Bargauge | 1 | Ranked metrics |
| Heatmap | 1 | Pattern detection |

---

## Query Complexity

| Dashboard | Simple | Medium | Complex |
|-----------|--------|--------|---------|
| Service Health | 3 | 5 | 2 |
| Infrastructure | 6 | 6 | 2 |
| Business Metrics | 8 | 7 | 2 |
| Alerting & SLO | 4 | 6 | 6 |

**Simple**: Direct metric queries (e.g., `listings_db_connections_open`)
**Medium**: Rate calculations (e.g., `rate(listings_grpc_requests_total[5m])`)
**Complex**: Histogram quantiles, multi-metric calculations

---

## Alert Integration Points

Dashboards visualize metrics that trigger these alert categories:

1. **Service Availability** (Critical)
   - Service down
   - High error rate (>5%)
   - SLO breach

2. **Performance** (Warning)
   - High latency (P95 > 200ms, P99 > 1s)
   - Slow queries (P99 > 100ms)

3. **Capacity** (Warning)
   - High CPU/memory usage
   - DB connection saturation
   - Queue backup

4. **Business** (Info)
   - Low cache hit ratio (<70%)
   - High rate limit rejections
   - Inventory alerts

---

## Deployment Checklist

- [ ] Grafana installed and accessible
- [ ] Prometheus datasource configured
- [ ] Metrics endpoint reachable (http://localhost:9093/metrics)
- [ ] Prometheus scraping listings service
- [ ] Import all 4 dashboards
- [ ] Verify data appears in dashboards
- [ ] Set up Prometheus alert rules
- [ ] Configure Alertmanager routing
- [ ] Test alert flow (fire test alert)
- [ ] Share dashboard URLs with team
- [ ] Add to runbooks and documentation

---

## Maintenance Schedule

### Daily
- Check Service Health dashboard for anomalies

### Weekly
- Review Alerting & SLO dashboard for SLO compliance
- Check for new slow queries in Infrastructure
- Review business metrics trends

### Monthly
- Update alert thresholds based on capacity changes
- Review and optimize dashboard queries
- Update documentation with lessons learned
- Check for new metrics to add

### Quarterly
- Full dashboard review and optimization
- Update SLO targets if needed
- Refine panel layouts based on usage
- Add new features as service evolves

---

## Version History

### v1.0.0 (2025-01-05)
- Initial production release
- 4 comprehensive dashboards
- 57 panels covering 67+ metrics
- Full SLO tracking
- Alert integration points
- Complete documentation

---

## Files Generated

```
deployment/grafana/
├── dashboards/
│   ├── service-health.json       10.7 KB  (10 panels)
│   ├── infrastructure.json       13.2 KB  (14 panels)
│   ├── business-metrics.json     16.4 KB  (17 panels)
│   ├── alerting-slo.json         18.1 KB  (16 panels)
│   └── README.md                 13.5 KB  (Full docs)
├── import-dashboards.sh          2.8 KB   (Import script)
├── validate-dashboards.sh        2.3 KB   (Validation)
├── QUICKSTART.md                 4.2 KB   (Quick guide)
└── DASHBOARD_MANIFEST.md         [This file]
```

**Total**: 81.2 KB of dashboard assets

---

## Contact & Support

For issues with dashboards:
1. Check README.md for troubleshooting
2. Validate metrics endpoint is working
3. Verify Prometheus scrape config
4. Check Grafana datasource connection
5. Review dashboard JSON for query errors

**Metrics Source Code**: `/internal/metrics/metrics.go`
**Documentation**: Complete in `dashboards/README.md`
