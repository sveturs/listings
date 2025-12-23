# OpenSearch Monitoring - Быстрый старт

## Запуск сервера

```bash
cd /p/github.com/vondi-global/listings
./bin/listings-server
```

## Проверка мониторинга

### 1. Health Check OpenSearch
```bash
curl http://localhost:8086/health/opensearch | jq '.'
```

**Ожидаемый результат (healthy):**
```json
{
  "status": "healthy",
  "cluster_health": "green",
  "index_exists": true,
  "docs_count": 1234,
  "index_size_mb": 56.78,
  "latency_ms": 12,
  "shards": {
    "total": 5,
    "successful": 5,
    "failed": 0,
    "primary": 1,
    "replica": 4
  }
}
```

### 2. Index Statistics
```bash
curl http://localhost:8086/metrics/opensearch | jq '.'
```

**Ожидаемый результат:**
```json
{
  "index": "marketplace_listings",
  "docs_count": 1234,
  "store_size": "56.8 MB",
  "store_size_mb": 56.78,
  "segment_count": 12,
  "shards": {
    "total": 5,
    "successful": 5,
    "failed": 0,
    "primary": 1,
    "replica": 4
  },
  "refresh_time_ms": 1234,
  "flush_time_ms": 567
}
```

### 3. Prometheus Metrics
```bash
curl http://localhost:8086/metrics | grep opensearch
```

**Ожидаемый результат:**
```prometheus
vondi_listings_opensearch_cluster_status{cluster="marketplace_listings"} 2
vondi_listings_opensearch_indexed_documents{index="marketplace_listings"} 1234
vondi_listings_opensearch_index_size_mb{index="marketplace_listings"} 56.78
vondi_listings_opensearch_request_duration_seconds_bucket{operation="search",status="success",le="0.005"} 100
```

## Проверка логов

Сборщик статистики логирует каждые 60 секунд:

```bash
tail -f /tmp/listings-microservice.log | grep opensearch
```

**Ожидаемый лог:**
```json
{"level":"info","component":"opensearch_stats_collector","status":"healthy","cluster":"green","docs":1234,"size_mb":56.78,"latency":"12ms","time":"2025-12-19T23:00:00Z","message":"OpenSearch stats collected"}
```

## Grafana Dashboard (рекомендуемые панели)

### Panel 1: Cluster Health Status
```
Query: vondi_listings_opensearch_cluster_status
Visualization: Gauge
Thresholds: 0 (red), 1 (yellow), 2 (green)
```

### Panel 2: Document Count
```
Query: vondi_listings_opensearch_indexed_documents
Visualization: Stat
```

### Panel 3: Search Latency p95
```
Query: histogram_quantile(0.95, sum(rate(vondi_listings_opensearch_search_latency_seconds_bucket[5m])) by (le))
Visualization: Graph
Unit: seconds
```

### Panel 4: Index Size
```
Query: vondi_listings_opensearch_index_size_mb
Visualization: Graph
Unit: MB
```

### Panel 5: Request Success Rate
```
Query: sum(rate(vondi_listings_opensearch_requests_total{status="success"}[5m])) / sum(rate(vondi_listings_opensearch_requests_total[5m]))
Visualization: Graph
Unit: percentunit (0.0-1.0)
```

## Troubleshooting

### OpenSearch недоступен
```bash
curl http://localhost:8086/health/opensearch
# Response: {"status":"unavailable","error":"OpenSearch client not initialized"}
```

**Решение:** Проверить конфигурацию в `.env`:
```bash
VONDILISTINGS_SEARCH_ADDRESSES=http://localhost:9200
VONDILISTINGS_SEARCH_INDEX=marketplace_listings
```

### Cluster status = yellow
```json
{"cluster_health": "yellow", "status": "degraded"}
```

**Причина:** Нет replica шардов (одиночная нода OpenSearch).

**Решение:** Это нормально для dev окружения. В продакшне используйте cluster с несколькими нодами.

### High latency
```json
{"latency_ms": 1500, "status": "degraded"}
```

**Причины:**
1. Большой индекс (много документов)
2. Медленные диски
3. CPU перегружен

**Решение:**
```bash
# Проверить размер индекса
curl http://localhost:8086/metrics/opensearch | jq '.store_size_mb'

# Проверить количество сегментов (оптимально < 50)
curl http://localhost:8086/metrics/opensearch | jq '.segment_count'

# Если сегментов много - выполнить force merge
curl -X POST "http://localhost:9200/marketplace_listings/_forcemerge?max_num_segments=1"
```

## Мониторинг в реальном времени

### Watch mode (каждые 5 секунд)
```bash
watch -n 5 'curl -s http://localhost:8086/health/opensearch | jq "."'
```

### Непрерывный лог метрик
```bash
while true; do
  echo "=== $(date) ==="
  curl -s http://localhost:8086/metrics/opensearch | jq '{docs: .docs_count, size_mb: .store_size_mb, segments: .segment_count}'
  sleep 10
done
```

## Ссылки

- **Полный отчёт:** [OPENSEARCH_MONITORING_REPORT.md](./OPENSEARCH_MONITORING_REPORT.md)
- **Исходный код мониторинга:** [/internal/repository/opensearch/monitoring.go](../internal/repository/opensearch/monitoring.go)
- **HTTP endpoints:** [/internal/transport/http/health_handler.go](../internal/transport/http/health_handler.go)
