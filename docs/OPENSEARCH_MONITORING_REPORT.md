# OpenSearch Monitoring Implementation Report

## Реализация: Фаза 6 - Мониторинг OpenSearch

**Дата:** 2025-12-19
**Статус:** Успешно реализовано
**Директория:** `/p/github.com/vondi-global/listings/`

---

## Реализованные компоненты

### 1. Prometheus метрики (monitoring.go)

Созданы следующие метрики для OpenSearch:

#### Запросы и производительность
- `vondi_listings_opensearch_request_duration_seconds` - длительность запросов
- `vondi_listings_opensearch_requests_total` - общее количество запросов
- `vondi_listings_opensearch_search_latency_seconds` - latency поисковых запросов

#### Статус кластера и индекса
- `vondi_listings_opensearch_cluster_status` - здоровье кластера (0=red, 1=yellow, 2=green)
- `vondi_listings_opensearch_indexed_documents` - количество документов в индексе
- `vondi_listings_opensearch_index_size_mb` - размер индекса в MB
- `vondi_listings_opensearch_shard_count` - количество шардов (primary/replica)

#### Переиндексация
- `vondi_listings_opensearch_reindex_progress` - прогресс переиндексации (0-100%)

---

### 2. Health Check (HealthCheckDetailed)

Метод `HealthCheckDetailed(ctx) (*HealthStatus, error)` возвращает:

```go
type HealthStatus struct {
    Status        string        // healthy, degraded, unhealthy
    ClusterHealth string        // green, yellow, red
    IndexExists   bool
    DocsCount     int64
    LastCheck     time.Time
    Latency       time.Duration
    IndexSizeMB   float64
    Shards        ShardInfo
}
```

**Логика определения статуса:**
- `unhealthy`: cluster_health=red ИЛИ index не существует
- `degraded`: cluster_health=yellow ИЛИ latency > 500ms
- `healthy`: cluster_health=green И latency < 500ms

---

### 3. Index Statistics (GetIndexStats)

Детальная статистика индекса:

```go
type IndexStats struct {
    Index        string
    DocsCount    int64
    StoreSize    string        // человекочитаемый формат
    StoreSizeMB  float64
    SegmentCount int
    Shards       ShardInfo
    RefreshTime  time.Duration
    FlushTime    time.Duration
}
```

---

### 4. Periodic Stats Collector

Автоматический сборщик статистики:

- **Класс:** `StatsCollector`
- **Интервал:** 60 секунд (по умолчанию)
- **Действия:**
  - Вызывает `HealthCheckDetailed()`
  - Обновляет Prometheus метрики
  - Логирует статус в zerolog

**Запуск в main.go:**
```go
if searchClient != nil {
    statsCollector := opensearchRepo.NewStatsCollector(searchClient, 60*time.Second, zerologLogger)
    go statsCollector.Start(context.Background())
    defer statsCollector.Stop()
}
```

---

### 5. HTTP Endpoints

#### GET /health/opensearch
Детальный health check OpenSearch:

**Response (200 OK):**
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
  },
  "last_check": "2025-12-19T23:00:00Z",
  "timestamp": 1734649200
}
```

**Response (503 Service Unavailable):**
```json
{
  "status": "unhealthy",
  "cluster_health": "red",
  "index_exists": true,
  "docs_count": 1234,
  "error": "index marketplace_listings does not exist",
  "timestamp": 1734649200
}
```

#### GET /metrics/opensearch
Детальная статистика индекса:

**Response (200 OK):**
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
  "flush_time_ms": 567,
  "timestamp": 1734649200
}
```

#### GET /metrics/opensearch/cluster
Краткий статус кластера:

**Response (200 OK):**
```json
{
  "cluster_status": "green",
  "index_exists": true,
  "overall_status": "healthy",
  "timestamp": 1734649200
}
```

---

## Интеграция с существующей системой

### HealthHandler (обновлён)
- Добавлен опциональный `searchClient *opensearchRepo.Client`
- Новый конструктор: `NewHealthHandlerWithOpenSearch(checker, searchClient, logger)`
- Автоматическая регистрация OpenSearch endpoints при наличии клиента

### main.go (обновлён)
```go
// Initialize OpenSearch stats collector
if searchClient != nil {
    statsCollector := opensearchRepo.NewStatsCollector(searchClient, 60*time.Second, zerologLogger)
    go statsCollector.Start(context.Background())
    defer statsCollector.Stop()
}

// Use enhanced health handler
if searchClient != nil {
    healthHandler = httpTransport.NewHealthHandlerWithOpenSearch(healthChecker, searchClient, zerologLogger)
} else {
    healthHandler = httpTransport.NewHealthHandler(healthChecker, zerologLogger)
}
```

---

## Prometheus Integration

Все метрики автоматически доступны на стандартном endpoint:

**GET /metrics** (Prometheus format)
```prometheus
# HELP vondi_listings_opensearch_cluster_status Cluster health status (0=red, 1=yellow, 2=green)
# TYPE vondi_listings_opensearch_cluster_status gauge
vondi_listings_opensearch_cluster_status{cluster="marketplace_listings"} 2

# HELP vondi_listings_opensearch_indexed_documents Number of documents in index
# TYPE vondi_listings_opensearch_indexed_documents gauge
vondi_listings_opensearch_indexed_documents{index="marketplace_listings"} 1234

# HELP vondi_listings_opensearch_request_duration_seconds Duration of OpenSearch requests
# TYPE vondi_listings_opensearch_request_duration_seconds histogram
vondi_listings_opensearch_request_duration_seconds_bucket{operation="search",status="success",le="0.005"} 100
...
```

---

## Использование

### Проверка здоровья OpenSearch
```bash
curl http://localhost:8086/health/opensearch | jq '.'
```

### Получение статистики индекса
```bash
curl http://localhost:8086/metrics/opensearch | jq '.'
```

### Prometheus scraping
```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'listings-service'
    static_configs:
      - targets: ['localhost:8086']
    metrics_path: '/metrics'
    scrape_interval: 15s
```

### Grafana Dashboard
Создайте dashboard с панелями:

1. **Cluster Health Gauge**
   - Metric: `vondi_listings_opensearch_cluster_status`
   - Thresholds: red=0, yellow=1, green=2

2. **Document Count**
   - Metric: `vondi_listings_opensearch_indexed_documents`

3. **Request Latency (p95, p99)**
   - Metric: `histogram_quantile(0.95, vondi_listings_opensearch_request_duration_seconds)`

4. **Index Size**
   - Metric: `vondi_listings_opensearch_index_size_mb`

5. **Search Latency by Query Type**
   - Metric: `vondi_listings_opensearch_search_latency_seconds`

---

## Компиляция

Проект успешно компилируется:

```bash
cd /p/github.com/vondi-global/listings
go build -o bin/listings-server ./cmd/server/main.go
# Success! Binary created: bin/listings-server
```

---

## Файлы

### Созданные файлы:
1. `/internal/repository/opensearch/monitoring.go` (459 строк)
   - Prometheus метрики
   - HealthCheckDetailed
   - GetIndexStats
   - StatsCollector

2. `/internal/transport/http/opensearch_monitoring_handler.go` (160 строк)
   - Отдельный handler для OpenSearch endpoints

### Изменённые файлы:
1. `/internal/transport/http/health_handler.go`
   - Добавлены OpenSearch endpoints
   - NewHealthHandlerWithOpenSearch конструктор

2. `/cmd/server/main.go`
   - Инициализация StatsCollector
   - Использование enhanced health handler

---

## Метрики и Health Checks Summary

| Тип | Компонент | Описание |
|-----|-----------|----------|
| **Метрики** | cluster_status | Здоровье кластера (gauge) |
| **Метрики** | indexed_documents | Количество документов (gauge) |
| **Метрики** | index_size_mb | Размер индекса (gauge) |
| **Метрики** | request_duration_seconds | Latency запросов (histogram) |
| **Метрики** | search_latency_seconds | Latency поиска (histogram) |
| **Метрики** | shard_count | Количество шардов (gauge) |
| **Health** | GET /health/opensearch | Детальный health check |
| **Health** | GET /metrics/opensearch | Статистика индекса |
| **Health** | GET /metrics/opensearch/cluster | Статус кластера |
| **Worker** | StatsCollector | Периодический сбор (60s) |

---

## Рекомендации для продакшна

### Алерты в Prometheus/Grafana
```yaml
# prometheus-alerts.yml
groups:
  - name: opensearch
    rules:
      # Cluster RED
      - alert: OpenSearchClusterRed
        expr: vondi_listings_opensearch_cluster_status{cluster="marketplace_listings"} == 0
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "OpenSearch cluster is RED"

      # High latency
      - alert: OpenSearchHighLatency
        expr: histogram_quantile(0.95, vondi_listings_opensearch_search_latency_seconds) > 1
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "OpenSearch p95 latency > 1s"

      # Low disk space (size > 80% of expected max)
      - alert: OpenSearchIndexSizeHigh
        expr: vondi_listings_opensearch_index_size_mb > 8000
        for: 1h
        labels:
          severity: warning
        annotations:
          summary: "OpenSearch index size exceeds 8GB"
```

### Logging
Все операции логируются через zerolog:
- Health check failures → ERROR
- Slow queries (>500ms) → WARN
- Stats collection → DEBUG

---

## Выполнено 100%

- ✅ Prometheus метрики для OpenSearch
- ✅ Детальный health check (HealthCheckDetailed)
- ✅ Index statistics (GetIndexStats)
- ✅ Periodic stats collector (StatsCollector)
- ✅ HTTP endpoints (/health/opensearch, /metrics/opensearch)
- ✅ Интеграция с main.go
- ✅ Компиляция без ошибок
- ✅ Совместимость с существующими health checks

---

**Автор:** Claude Sonnet 4.5
**Дата:** 2025-12-19
