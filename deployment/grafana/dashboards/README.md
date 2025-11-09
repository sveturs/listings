# Grafana Dashboards for Listings Microservice

Production-ready Grafana dashboards for comprehensive monitoring of the listings microservice.

## Dashboards Overview

### 1. Service Health Dashboard (`service-health.json`)
**Purpose**: Real-time service health and performance monitoring

**Key Metrics**:
- **Service Uptime**: 24-hour availability percentage
- **Current RPS**: Real-time requests per second
- **Total Requests**: 24-hour request volume
- **Error Rate**: Percentage of failed requests
- **Request Latency**: P50, P95, P99 percentiles over time
- **Throughput by Method**: RPS breakdown by gRPC method
- **Error Analysis**: Error rates and status code distribution
- **Request Volume Heatmap**: Hourly request patterns
- **Active Requests**: In-flight request gauge

**Use Cases**:
- Monitor service health at a glance
- Identify performance degradation
- Track request patterns and peak times
- Quick incident detection

---

### 2. Infrastructure Dashboard (`infrastructure.json`)
**Purpose**: System resources and infrastructure monitoring

**Key Metrics**:
- **Resource Usage**: CPU, Memory, Goroutines (gauges + timeseries)
- **Database Connections**: Open/idle connection tracking
- **Database Query Performance**: Duration histogram and slow query detection
- **Redis Performance**: Commands/sec, hit ratio, cache operations
- **Rate Limiting**: Evaluations, allowed, and rejected requests
- **Network I/O**: File descriptors and request rates

**Use Cases**:
- Capacity planning and resource optimization
- Database connection pool tuning
- Cache performance optimization
- Rate limiting effectiveness monitoring
- Infrastructure troubleshooting

---

### 3. Business Metrics Dashboard (`business-metrics.json`)
**Purpose**: Business operations and feature usage tracking

**Key Metrics**:
- **Listing Operations**: Created/updated/deleted counts and rates
- **Inventory Management**: Stock operations, out-of-stock products
- **Inventory Movements**: Movement types distribution
- **Low Stock Alerts**: Products approaching threshold
- **Product Views**: Total views and top 10 most viewed products
- **Search Analytics**: Query volume and search rate
- **Cache Efficiency**: Hit ratio and hit/miss breakdown
- **Indexing Queue**: Queue size, job processing, and duration

**Use Cases**:
- Business intelligence and analytics
- Feature adoption tracking
- Inventory health monitoring
- Search performance analysis
- Worker queue health monitoring

---

### 4. Alerting & SLO Dashboard (`alerting-slo.json`)
**Purpose**: SLO compliance tracking and alert management

**Key Metrics**:
- **Active Alerts**: Real-time alert table with severity
- **Service Availability SLO**: 30-day availability (target: 99.9%)
- **Success Rate SLO**: 7-day success rate (target: 99.5%)
- **Error Budget**: Remaining error budget percentage
- **Latency SLOs**: P95 < 200ms, P99 < 1000ms compliance
- **Alert History**: Timeline of fired alerts
- **Error Budget Burn Rate**: Current vs historical burn rate
- **Error Analysis**: Errors by component and type
- **Rate Limit Rejections**: Tracking rejected requests
- **On-Call Status**: Current on-call information
- **Incidents**: Critical incident count (7 days)

**Use Cases**:
- SLO compliance monitoring
- Error budget tracking
- Incident management
- Alert triage and investigation
- Reliability engineering

---

## Installation

### Option 1: Grafana UI Import

1. Open Grafana (http://your-grafana-url)
2. Navigate to **Dashboards** → **Import**
3. Click **Upload JSON file**
4. Select one of the dashboard JSON files
5. Choose your Prometheus datasource
6. Click **Import**

### Option 2: Grafana Provisioning

1. Copy dashboard files to Grafana provisioning directory:
   ```bash
   cp *.json /etc/grafana/provisioning/dashboards/
   ```

2. Create/update provisioning config (`/etc/grafana/provisioning/dashboards/listings.yaml`):
   ```yaml
   apiVersion: 1

   providers:
     - name: 'Listings Dashboards'
       orgId: 1
       folder: 'Listings Service'
       type: file
       disableDeletion: false
       updateIntervalSeconds: 30
       allowUiUpdates: true
       options:
         path: /etc/grafana/provisioning/dashboards
         foldersFromFilesStructure: false
   ```

3. Restart Grafana:
   ```bash
   sudo systemctl restart grafana-server
   ```

### Option 3: Kubernetes ConfigMap

```bash
kubectl create configmap listings-grafana-dashboards \
  --from-file=service-health.json \
  --from-file=infrastructure.json \
  --from-file=business-metrics.json \
  --from-file=alerting-slo.json \
  -n monitoring
```

---

## Configuration

### Required Variables

All dashboards use these template variables:

- **`$datasource`**: Prometheus datasource (default: "Prometheus")
- **`$interval`**: Time interval for aggregation (auto, 1m, 5m, 10m, 30m, 1h)

### Annotations

Dashboards include annotations for:

- **Deployments**: Service restarts/deployments (blue markers)
- **Incidents**: Critical alerts (red markers)

### Time Range

- **Default**: Last 24 hours
- **Refresh**: 30 seconds auto-refresh
- **Configurable**: Time picker available for custom ranges

---

## Alert Rules (Example)

The dashboards visualize metrics that should trigger alerts. Example Prometheus alert rules:

```yaml
groups:
  - name: listings_service
    rules:
      # High Error Rate
      - alert: HighErrorRate
        expr: |
          (sum(rate(listings_grpc_requests_total{status!="OK"}[5m]))
          / sum(rate(listings_grpc_requests_total[5m]))) * 100 > 5
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"
          description: "Error rate is {{ $value }}% (threshold: 5%)"

      # High Latency P95
      - alert: HighLatencyP95
        expr: |
          histogram_quantile(0.95,
          sum(rate(listings_grpc_request_duration_seconds_bucket[5m])) by (le)) > 0.2
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "P95 latency exceeds SLO"
          description: "P95 latency is {{ $value }}s (threshold: 200ms)"

      # High Latency P99
      - alert: HighLatencyP99
        expr: |
          histogram_quantile(0.99,
          sum(rate(listings_grpc_request_duration_seconds_bucket[5m])) by (le)) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "P99 latency exceeds SLO"
          description: "P99 latency is {{ $value }}s (threshold: 1s)"

      # Service Down
      - alert: ServiceDown
        expr: up{job="listings"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Listings service is down"
          description: "Service has been down for 1 minute"

      # High Memory Usage
      - alert: HighMemoryUsage
        expr: process_resident_memory_bytes{job="listings"} > 2147483648
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High memory usage"
          description: "Memory usage is {{ $value | humanize }}B (threshold: 2GB)"

      # Database Connection Pool Saturation
      - alert: DBConnectionPoolSaturation
        expr: listings_db_connections_open > 90
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Database connection pool near saturation"
          description: "Open connections: {{ $value }} (max: 100)"

      # High Cache Miss Rate
      - alert: HighCacheMissRate
        expr: |
          (sum(rate(listings_cache_misses_total[5m]))
          / (sum(rate(listings_cache_hits_total[5m]))
          + sum(rate(listings_cache_misses_total[5m])))) * 100 > 50
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "High cache miss rate"
          description: "Cache miss rate is {{ $value }}% (threshold: 50%)"

      # Indexing Queue Backup
      - alert: IndexingQueueBackup
        expr: listings_indexing_queue_size > 500
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "Indexing queue is backing up"
          description: "Queue size: {{ $value }} (threshold: 500)"

      # High Rate Limit Rejections
      - alert: HighRateLimitRejections
        expr: sum(rate(listings_rate_limit_rejected_total[5m])) > 10
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High rate of rate limit rejections"
          description: "Rejection rate: {{ $value }}/s (threshold: 10/s)"
```

---

## Metrics Reference

### Namespace: `listings`

All metrics are prefixed with `listings_` (e.g., `listings_grpc_requests_total`).

#### gRPC Metrics
- `grpc_requests_total{method, status}` - Counter
- `grpc_request_duration_seconds{method}` - Histogram
- `grpc_handler_requests_active{method}` - Gauge

#### HTTP Metrics
- `http_requests_total{method, path, status}` - Counter
- `http_request_duration_seconds{method, path}` - Histogram
- `http_requests_in_flight` - Gauge

#### Business Metrics
- `listings_created_total` - Counter
- `listings_updated_total` - Counter
- `listings_deleted_total` - Counter
- `listings_searched_total` - Counter

#### Database Metrics
- `db_connections_open` - Gauge
- `db_connections_idle` - Gauge
- `db_query_duration_seconds{operation}` - Histogram

#### Cache Metrics
- `cache_hits_total{cache_type}` - Counter
- `cache_misses_total{cache_type}` - Counter

#### Inventory Metrics
- `inventory_product_views_total{product_id}` - Counter
- `inventory_product_views_errors_total` - Counter
- `inventory_stock_operations_total{operation, status}` - Counter
- `inventory_stock_low_threshold_reached_total{product_id, storefront_id}` - Counter
- `inventory_movements_recorded_total{movement_type}` - Counter
- `inventory_movements_errors_total{reason}` - Counter
- `inventory_stock_value{storefront_id, product_id}` - Gauge
- `inventory_out_of_stock_products` - Gauge

#### Indexing Queue Metrics
- `indexing_queue_size` - Gauge
- `indexing_jobs_processed_total{operation, status}` - Counter
- `indexing_job_duration_seconds` - Histogram

#### Rate Limiting Metrics
- `rate_limit_hits_total{method, identifier_type}` - Counter
- `rate_limit_allowed_total{method, identifier_type}` - Counter
- `rate_limit_rejected_total{method, identifier_type}` - Counter

#### Timeout Metrics
- `timeouts_total{method}` - Counter
- `near_timeouts_total{method}` - Counter
- `timeout_duration_seconds{method}` - Histogram

#### Error Metrics
- `errors_total{component, error_type}` - Counter

---

## SLO Targets

### Availability SLO
- **Target**: 99.9% (30 days)
- **Error Budget**: 0.1% = ~43 minutes downtime per month
- **Measurement**: `(successful_requests / total_requests) * 100`

### Latency SLO
- **P95**: < 200ms
- **P99**: < 1000ms
- **Measurement**: `histogram_quantile(0.95|0.99, grpc_request_duration_seconds)`

### Success Rate SLO
- **Target**: 99.5% (7 days)
- **Error Budget**: 0.5% of requests
- **Measurement**: `(status=OK requests / total_requests) * 100`

---

## Troubleshooting

### Dashboard shows "No data"

1. **Check Prometheus datasource**:
   - Verify datasource is configured in Grafana
   - Test connection in **Configuration** → **Data Sources**

2. **Verify metrics are exported**:
   ```bash
   curl http://localhost:9093/metrics | grep listings_
   ```

3. **Check Prometheus is scraping**:
   - Open Prometheus UI (http://localhost:9090)
   - Go to **Status** → **Targets**
   - Verify `listings` job is UP

### Queries returning empty results

1. **Check metric names**: Ensure they match exactly (case-sensitive)
2. **Verify label filters**: Labels like `{job="listings"}` must match your config
3. **Check time range**: Some metrics may not have historical data

### High memory usage in Grafana

- Reduce time range
- Increase query interval (`$interval`)
- Limit number of series in queries

---

## Best Practices

1. **Start with Service Health**: Use it as your primary operational dashboard
2. **Set up alerts**: Configure Prometheus alerts based on dashboard thresholds
3. **Regular review**: Weekly review of SLO dashboard to track trends
4. **Customize thresholds**: Adjust gauge thresholds based on your capacity
5. **Use annotations**: Mark deployments and incidents for correlation
6. **Dashboard maintenance**: Update queries when metrics change
7. **Performance**: Use `$interval` variable for high-cardinality queries

---

## Version History

- **v1.0.0** (2025-01-XX): Initial production release
  - 4 comprehensive dashboards
  - 67+ metrics visualization
  - SLO tracking and alerting
  - Production-ready with tested queries

---

## Support

For issues or questions:
- **Metrics**: Check `/p/github.com/sveturs/listings/internal/metrics/metrics.go`
- **Documentation**: See project README and ADRs
- **Monitoring**: Prometheus endpoint at `http://localhost:9093/metrics`

---

## License

Part of the Listings microservice. Same license applies.
