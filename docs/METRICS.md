# Orders Microservice - Metrics Catalog

**Version:** 1.0.0
**Last Updated:** 2025-11-14

This document describes all metrics exported by the Orders Microservice for monitoring, alerting, and performance analysis.

---

## Table of Contents

1. [gRPC Metrics](#grpc-metrics)
2. [Orders Domain Metrics](#orders-domain-metrics)
3. [Cart Domain Metrics](#cart-domain-metrics)
4. [Database Metrics](#database-metrics)
5. [Redis Cache Metrics](#redis-cache-metrics)
6. [System Metrics](#system-metrics)
7. [Alert Thresholds](#alert-thresholds)

---

## gRPC Metrics

### `grpc_server_handled_total`
**Type:** Counter
**Description:** Total number of RPCs completed on the server, regardless of success or failure
**Labels:**
- `grpc_service`: Service name (e.g., `listings.v1.OrdersService`)
- `grpc_method`: Method name (e.g., `CreateOrder`, `GetCart`)
- `grpc_code`: gRPC status code (e.g., `OK`, `NotFound`, `Internal`)

**Example Query:**
```promql
rate(grpc_server_handled_total{grpc_service="listings.v1.OrdersService"}[5m])
```

**Alert Threshold:** N/A (informational)

---

### `grpc_server_handling_seconds`
**Type:** Histogram
**Description:** Histogram of response latency (seconds) of gRPC that had been application-level handled by the server
**Labels:**
- `grpc_service`: Service name
- `grpc_method`: Method name
- `grpc_code`: gRPC status code

**Example Query (P95 latency):**
```promql
histogram_quantile(0.95, rate(grpc_server_handling_seconds_bucket{grpc_service="listings.v1.OrdersService"}[5m]))
```

**Alert Thresholds:**
- **GetCart:** P95 > 20ms
- **AddToCart:** P95 > 50ms
- **CreateOrder:** P95 > 200ms
- **ListOrders:** P95 > 100ms

---

### `grpc_server_started_total`
**Type:** Counter
**Description:** Total number of RPCs started on the server
**Labels:**
- `grpc_service`: Service name
- `grpc_method`: Method name

**Example Query:**
```promql
rate(grpc_server_started_total{grpc_service="listings.v1.OrdersService"}[5m])
```

**Alert Threshold:** N/A (informational)

---

### `grpc_server_msg_received_total`
**Type:** Counter
**Description:** Total number of RPC stream messages received on the server
**Labels:**
- `grpc_service`: Service name
- `grpc_method`: Method name

**Example Query:**
```promql
rate(grpc_server_msg_received_total{grpc_service="listings.v1.OrdersService"}[5m])
```

**Alert Threshold:** N/A (informational)

---

### `grpc_server_msg_sent_total`
**Type:** Counter
**Description:** Total number of gRPC stream messages sent by the server
**Labels:**
- `grpc_service`: Service name
- `grpc_method`: Method name

**Example Query:**
```promql
rate(grpc_server_msg_sent_total{grpc_service="listings.v1.OrdersService"}[5m])
```

**Alert Threshold:** N/A (informational)

---

## Orders Domain Metrics

### `orders_total`
**Type:** Counter
**Description:** Total number of orders created
**Labels:**
- `status`: Order status (e.g., `pending`, `confirmed`, `shipped`)
- `storefront_id`: Storefront ID

**Example Query:**
```promql
sum(rate(orders_total[5m])) by (status)
```

**Alert Threshold:** Sudden drop in order creation rate (>20% decrease over 15min)

---

### `orders_value_total`
**Type:** Counter
**Description:** Total value of all orders (in base currency)
**Labels:**
- `currency`: Currency code (e.g., `USD`, `EUR`, `RSD`)
- `storefront_id`: Storefront ID

**Example Query:**
```promql
sum(rate(orders_value_total{currency="USD"}[1h])) by (storefront_id)
```

**Alert Threshold:** N/A (business metric)

---

### `order_processing_duration_seconds`
**Type:** Histogram
**Description:** Time taken to process an order from creation to confirmation
**Labels:**
- `status`: Final order status

**Example Query (P95):**
```promql
histogram_quantile(0.95, rate(order_processing_duration_seconds_bucket[5m]))
```

**Alert Threshold:** P95 > 5s (order processing too slow)

---

### `order_failures_total`
**Type:** Counter
**Description:** Total number of failed order operations
**Labels:**
- `operation`: Operation type (e.g., `create`, `update`, `cancel`)
- `reason`: Failure reason (e.g., `insufficient_stock`, `payment_failed`)

**Example Query:**
```promql
rate(order_failures_total[5m])
```

**Alert Threshold:** Error rate > 1%

---

## Cart Domain Metrics

### `cart_items_added_total`
**Type:** Counter
**Description:** Total number of items added to carts
**Labels:**
- `storefront_id`: Storefront ID

**Example Query:**
```promql
sum(rate(cart_items_added_total[5m])) by (storefront_id)
```

**Alert Threshold:** N/A (informational)

---

### `cart_items_removed_total`
**Type:** Counter
**Description:** Total number of items removed from carts
**Labels:**
- `storefront_id`: Storefront ID

**Example Query:**
```promql
sum(rate(cart_items_removed_total[5m])) by (storefront_id)
```

**Alert Threshold:** N/A (informational)

---

### `cart_abandonment_rate`
**Type:** Gauge
**Description:** Percentage of carts abandoned (not converted to orders) in the last hour
**Labels:** None

**Example Query:**
```promql
cart_abandonment_rate
```

**Alert Threshold:** > 80% (high abandonment)

---

### `active_carts`
**Type:** Gauge
**Description:** Number of active carts (carts with items, not converted to orders)
**Labels:**
- `storefront_id`: Storefront ID

**Example Query:**
```promql
sum(active_carts) by (storefront_id)
```

**Alert Threshold:** N/A (informational)

---

## Database Metrics

### `db_connections_open`
**Type:** Gauge
**Description:** Number of open database connections
**Labels:**
- `database`: Database name

**Example Query:**
```promql
db_connections_open{database="listings_db"}
```

**Alert Threshold:** > 80% of max_connections (connection pool exhaustion)

---

### `db_connections_in_use`
**Type:** Gauge
**Description:** Number of database connections currently in use
**Labels:**
- `database`: Database name

**Example Query:**
```promql
db_connections_in_use{database="listings_db"}
```

**Alert Threshold:** N/A (informational)

---

### `db_connections_idle`
**Type:** Gauge
**Description:** Number of idle database connections
**Labels:**
- `database`: Database name

**Example Query:**
```promql
db_connections_idle{database="listings_db"}
```

**Alert Threshold:** N/A (informational)

---

### `db_query_duration_seconds`
**Type:** Histogram
**Description:** Duration of database queries in seconds
**Labels:**
- `query_type`: Type of query (e.g., `select`, `insert`, `update`)
- `table`: Table name

**Example Query (P95):**
```promql
histogram_quantile(0.95, rate(db_query_duration_seconds_bucket[5m]))
```

**Alert Threshold:** P95 > 100ms (slow queries)

---

### `db_query_errors_total`
**Type:** Counter
**Description:** Total number of database query errors
**Labels:**
- `error_type`: Error type (e.g., `timeout`, `connection_lost`, `constraint_violation`)

**Example Query:**
```promql
rate(db_query_errors_total[5m])
```

**Alert Threshold:** > 10 errors/min

---

## Redis Cache Metrics

### `redis_cache_hits`
**Type:** Counter
**Description:** Total number of cache hits
**Labels:**
- `cache_key_prefix`: Cache key prefix (e.g., `cart:`, `order:`)

**Example Query:**
```promql
rate(redis_cache_hits[5m])
```

**Alert Threshold:** N/A (informational)

---

### `redis_cache_misses`
**Type:** Counter
**Description:** Total number of cache misses
**Labels:**
- `cache_key_prefix`: Cache key prefix

**Example Query:**
```promql
rate(redis_cache_misses[5m])
```

**Alert Threshold:** Hit rate < 70%

---

### `redis_cache_hit_rate`
**Type:** Gauge
**Description:** Cache hit rate (0.0 to 1.0)
**Labels:** None

**Example Query:**
```promql
redis_cache_hit_rate
```

**Calculation:**
```promql
rate(redis_cache_hits[5m]) / (rate(redis_cache_hits[5m]) + rate(redis_cache_misses[5m]))
```

**Alert Threshold:** < 70%

---

### `redis_operations_duration_seconds`
**Type:** Histogram
**Description:** Duration of Redis operations in seconds
**Labels:**
- `operation`: Operation type (e.g., `get`, `set`, `del`)

**Example Query (P95):**
```promql
histogram_quantile(0.95, rate(redis_operations_duration_seconds_bucket[5m]))
```

**Alert Threshold:** P95 > 10ms (Redis latency spike)

---

## System Metrics

### `process_cpu_seconds_total`
**Type:** Counter
**Description:** Total user and system CPU time spent in seconds
**Labels:** None

**Example Query (CPU usage %):**
```promql
rate(process_cpu_seconds_total{job="orders-service"}[5m]) * 100
```

**Alert Threshold:** > 80% sustained for 5 minutes

---

### `process_resident_memory_bytes`
**Type:** Gauge
**Description:** Resident memory size in bytes
**Labels:** None

**Example Query:**
```promql
process_resident_memory_bytes{job="orders-service"}
```

**Alert Threshold:** > 1GB (memory leak detection)

---

### `go_goroutines`
**Type:** Gauge
**Description:** Number of goroutines that currently exist
**Labels:** None

**Example Query:**
```promql
go_goroutines{job="orders-service"}
```

**Alert Threshold:** > 10,000 (potential goroutine leak)

---

### `go_gc_duration_seconds`
**Type:** Summary
**Description:** A summary of the pause duration of garbage collection cycles
**Labels:** None

**Example Query:**
```promql
rate(go_gc_duration_seconds_sum{job="orders-service"}[5m])
```

**Alert Threshold:** P99 > 100ms (GC pauses affecting latency)

---

### `go_memstats_alloc_bytes`
**Type:** Gauge
**Description:** Number of bytes allocated and still in use
**Labels:** None

**Example Query:**
```promql
go_memstats_alloc_bytes{job="orders-service"}
```

**Alert Threshold:** Continuously increasing (memory leak)

---

## Alert Thresholds

### Critical Alerts (Immediate Action Required)

| Metric | Threshold | Impact |
|--------|-----------|--------|
| `grpc_server_error_rate` | > 5% | High error rate, service degradation |
| `db_connections_open` | > 90% of max | Connection pool exhaustion |
| `order_failures_total` | > 50 errors/min | Order creation failing |
| `process_cpu_seconds_total` | > 90% for 5min | CPU exhaustion |
| `redis_cache_hit_rate` | < 30% | Cache ineffective |

### Warning Alerts (Investigation Needed)

| Metric | Threshold | Impact |
|--------|-----------|--------|
| `grpc_server_error_rate` | > 1% | Elevated error rate |
| `p95_latency` | > target + 20% | Performance degradation |
| `cart_abandonment_rate` | > 80% | High cart abandonment |
| `db_query_duration_seconds_p95` | > 100ms | Slow database queries |
| `go_goroutines` | > 5,000 | Potential goroutine leak |

### Info Alerts (Monitoring Only)

| Metric | Threshold | Impact |
|--------|-----------|--------|
| `orders_total` | Trend analysis | Business KPI |
| `orders_value_total` | Trend analysis | Revenue tracking |
| `active_carts` | Trend analysis | User engagement |

---

## Grafana Dashboard

A comprehensive Grafana dashboard is available at:
`/p/github.com/sveturs/listings/monitoring/grafana/orders-dashboard.json`

**Dashboard Includes:**
- RPC request rate and latency
- Error rate by method
- Order status distribution
- Cart operations throughput
- Database connection pool
- Redis cache hit rate
- CPU and memory usage
- Goroutine count
- GC pause duration

---

## Prometheus Configuration

### Scrape Configuration

```yaml
scrape_configs:
  - job_name: 'orders-service'
    static_configs:
      - targets: ['localhost:50052']
    metrics_path: '/metrics'
    scrape_interval: 15s
```

### Recording Rules

```yaml
groups:
  - name: orders
    interval: 30s
    rules:
      - record: grpc_server_error_rate
        expr: |
          rate(grpc_server_handled_total{grpc_code!="OK"}[5m])
          /
          rate(grpc_server_handled_total[5m])

      - record: redis_cache_hit_rate
        expr: |
          rate(redis_cache_hits[5m])
          /
          (rate(redis_cache_hits[5m]) + rate(redis_cache_misses[5m]))
```

---

## Querying Tips

### Top 5 Slowest gRPC Methods
```promql
topk(5, histogram_quantile(0.95, rate(grpc_server_handling_seconds_bucket[5m])))
```

### Total Requests Per Second
```promql
sum(rate(grpc_server_handled_total{grpc_service="listings.v1.OrdersService"}[5m]))
```

### Error Rate by Method
```promql
sum(rate(grpc_server_handled_total{grpc_code!="OK"}[5m])) by (grpc_method)
```

### Database Connection Utilization %
```promql
(db_connections_in_use / db_connections_open) * 100
```

### Memory Growth Rate
```promql
deriv(process_resident_memory_bytes[10m])
```

---

## References

- [Prometheus Best Practices](https://prometheus.io/docs/practices/naming/)
- [gRPC Prometheus Integration](https://github.com/grpc-ecosystem/go-grpc-prometheus)
- [Grafana Dashboard Best Practices](https://grafana.com/docs/grafana/latest/best-practices/)

---

**Document Owner:** Development Team
**Contact:** devops@vondi.rs
**Version History:**
- v1.0.0 (2025-11-14): Initial version
