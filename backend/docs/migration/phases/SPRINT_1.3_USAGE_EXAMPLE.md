# Sprint 1.3: Async Indexing - Usage Examples

## Quick Start

### 1. Enable Async Indexing in Main App

**File:** `backend/internal/storage/postgres/db.go`

Find where OpenSearch repository is initialized and add:

```go
// After creating OpenSearch repository
osMarketplaceRepo := opensearch.NewRepository(osClient, indexName, db, searchWeights)

// Enable async indexing
if sqlxDB != nil {
    err := osMarketplaceRepo.EnableAsyncIndexing(
        sqlxDB,   // *sqlx.DB for DLQ
        5,        // 5 workers
        1000,     // queue size 1000
    )
    if err != nil {
        log.Printf("WARNING: Failed to enable async indexing: %v", err)
        log.Printf("Falling back to synchronous indexing")
    } else {
        log.Printf("✅ Async indexing enabled (workers: 5, queue: 1000)")
    }
}
```

### 2. Graceful Shutdown

**File:** `backend/cmd/api/main.go`

Add shutdown hook:

```go
func main() {
    // ... existing setup code ...

    // Get storage instance
    storage := getStorage() // your method to get storage

    // Shutdown hook
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)

    go func() {
        <-c
        log.Println("Shutting down async indexer...")

        // Get OpenSearch repository
        if osRepo := storage.GetOpenSearchRepository(); osRepo != nil {
            err := osRepo.ShutdownAsyncIndexer(30 * time.Second)
            if err != nil {
                log.Printf("Error during shutdown: %v", err)
            }
        }

        os.Exit(0)
    }()

    // ... start server ...
}
```

---

## Usage in Code

### Creating/Updating Listing

```go
// backend/internal/proj/c2c/handler/listing_handler.go

func (h *Handler) CreateListing(c *fiber.Ctx) error {
    var listing models.MarketplaceListing
    // ... parse request, validate ...

    // Save to PostgreSQL
    listingID, err := h.storage.CreateListing(ctx, &listing)
    if err != nil {
        return err
    }

    // Index to OpenSearch (now async!)
    // This returns immediately, indexing happens in background
    err = h.storage.IndexListing(ctx, &listing)
    if err != nil {
        // Log error but don't fail the request
        // Indexing will be retried automatically
        log.Printf("Indexing enqueue failed: %v", err)
    }

    return c.JSON(listing)
}
```

**No code changes needed!** The IndexListing method automatically uses async indexing if enabled.

---

## Monitoring

### Health Check Endpoint

Create an admin endpoint to check indexer health:

```go
// backend/internal/proj/admin/handler/monitoring.go

func (h *Handler) GetIndexerStatus(c *fiber.Ctx) error {
    osRepo := h.storage.GetOpenSearchRepository()
    if osRepo == nil {
        return c.Status(404).JSON(fiber.Map{
            "error": "OpenSearch not configured",
        })
    }

    indexer := osRepo.GetAsyncIndexer()
    if indexer == nil {
        return c.JSON(fiber.Map{
            "status": "disabled",
            "mode":   "synchronous",
        })
    }

    return c.JSON(fiber.Map{
        "status":     "enabled",
        "mode":       "asynchronous",
        "healthy":    indexer.IsHealthy(),
        "queue_size": indexer.GetQueueSize(),
    })
}
```

**Response example:**

```json
{
  "status": "enabled",
  "mode": "asynchronous",
  "healthy": true,
  "queue_size": 23
}
```

### DLQ Management Endpoint

```go
// backend/internal/proj/admin/handler/dlq_handler.go

func (h *Handler) GetDLQStats(c *fiber.Ctx) error {
    var stats struct {
        Total      int    `db:"total"`
        OldestTask string `db:"oldest_task"`
    }

    err := h.db.Get(&stats, `
        SELECT
            COUNT(*) as total,
            MIN(created_at) as oldest_task
        FROM opensearch_indexing_dlq
    `)
    if err != nil {
        return err
    }

    return c.JSON(stats)
}

func (h *Handler) RetryDLQ(c *fiber.Ctx) error {
    limit := c.QueryInt("limit", 100)

    osRepo := h.storage.GetOpenSearchRepository()
    indexer := osRepo.GetAsyncIndexer()

    err := indexer.RetryDLQ(c.Context(), limit)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.JSON(fiber.Map{
        "status":  "success",
        "retried": limit,
    })
}
```

---

## Cron Job for DLQ Retry

Create a cron job to automatically retry failed tasks:

**File:** `backend/cmd/dlq-retry/main.go`

```go
package main

import (
    "context"
    "log"
    "time"

    "backend/internal/storage/postgres"
)

func main() {
    // Initialize storage (same as main app)
    storage, err := postgres.NewDatabase(/* config */)
    if err != nil {
        log.Fatalf("Failed to init storage: %v", err)
    }
    defer storage.Close()

    // Get OpenSearch repository
    osRepo := storage.GetOpenSearchRepository()
    if osRepo == nil {
        log.Fatal("OpenSearch not configured")
    }

    indexer := osRepo.GetAsyncIndexer()
    if indexer == nil {
        log.Fatal("Async indexer not enabled")
    }

    // Retry failed tasks
    ctx := context.Background()
    err = indexer.RetryDLQ(ctx, 100)
    if err != nil {
        log.Fatalf("DLQ retry failed: %v", err)
    }

    log.Println("✅ DLQ retry completed successfully")
}
```

**Crontab entry (run every 5 minutes):**

```cron
*/5 * * * * /path/to/dlq-retry >> /var/log/dlq-retry.log 2>&1
```

---

## SQL Queries for DLQ

### View Failed Tasks

```sql
-- Recent failures
SELECT
    listing_id,
    action,
    attempts,
    last_error,
    created_at,
    last_attempt_at
FROM opensearch_indexing_dlq
ORDER BY last_attempt_at DESC
LIMIT 20;
```

### Failed Tasks by Error Type

```sql
SELECT
    last_error,
    COUNT(*) as count,
    MIN(created_at) as first_seen,
    MAX(last_attempt_at) as last_seen
FROM opensearch_indexing_dlq
GROUP BY last_error
ORDER BY count DESC;
```

### Clear Old Failed Tasks

```sql
-- Delete tasks older than 7 days
DELETE FROM opensearch_indexing_dlq
WHERE created_at < NOW() - INTERVAL '7 days';
```

### Manual Retry Specific Listing

```sql
-- Delete from DLQ to trigger fresh indexing
DELETE FROM opensearch_indexing_dlq
WHERE listing_id = 12345;

-- Then manually trigger indexing via API or code
```

---

## Prometheus Metrics Setup

### Scrape Configuration

**File:** `prometheus.yml`

```yaml
scrape_configs:
  - job_name: 'svetu-backend'
    static_configs:
      - targets: ['localhost:3000']
    metrics_path: '/metrics'
```

### Grafana Dashboard Queries

**Queue Size Panel:**

```promql
opensearch_indexing_queue_size
```

**Success Rate Panel:**

```promql
rate(opensearch_indexing_success_total[5m]) /
  (rate(opensearch_indexing_success_total[5m]) +
   rate(opensearch_indexing_failure_total[5m]))
```

**P95 Latency Panel:**

```promql
histogram_quantile(0.95,
  rate(opensearch_indexing_latency_seconds_bucket[5m])
)
```

**Retry Rate Panel:**

```promql
rate(opensearch_indexing_retry_total[5m])
```

---

## Alerting Rules

**File:** `prometheus/alerts.yml`

```yaml
groups:
  - name: opensearch_indexing
    interval: 30s
    rules:
      - alert: OpenSearchDLQGrowing
        expr: |
          (SELECT COUNT(*) FROM opensearch_indexing_dlq) > 100
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "OpenSearch DLQ has {{ $value }} failed tasks"
          description: "Investigate indexing failures"

      - alert: OpenSearchQueueSaturated
        expr: opensearch_indexing_queue_size > 900
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Queue {{ $value }}/1000 full"
          description: "Consider increasing workers or queue size"

      - alert: OpenSearchHighFailureRate
        expr: |
          rate(opensearch_indexing_failure_total[5m]) > 10
        for: 2m
        labels:
          severity: critical
        annotations:
          summary: "{{ $value }} indexing failures per second"
          description: "OpenSearch may be down or overloaded"
```

---

## Troubleshooting

### Issue: Queue is full, falling back to sync indexing

**Symptom:**
```
WARN Queue full, falling back to sync indexing listingID=123
```

**Solution:**
```go
// Increase queue size
osRepo.EnableAsyncIndexing(db, 5, 5000) // from 1000 to 5000
```

---

### Issue: DLQ accumulating tasks

**Symptom:**
```sql
SELECT COUNT(*) FROM opensearch_indexing_dlq;
-- Result: 500+
```

**Solution 1:** Retry from DLQ
```bash
# Via API
curl -X POST http://localhost:3000/api/v1/admin/dlq/retry?limit=100

# Via SQL
# Investigate errors first
SELECT last_error, COUNT(*) FROM opensearch_indexing_dlq GROUP BY last_error;

# Then retry via code or delete if issue is fixed
```

**Solution 2:** Check OpenSearch health
```bash
curl http://localhost:9200/_cluster/health
```

---

### Issue: Slow shutdown

**Symptom:**
```
Shutting down async indexer...
(waits 30 seconds)
```

**Cause:** Many tasks in queue waiting to complete

**Solution:**
```go
// Reduce shutdown timeout or drain queue faster
osRepo.ShutdownAsyncIndexer(10 * time.Second) // faster shutdown
```

---

## Performance Tuning

### Low Traffic (< 10 req/s):

```go
Workers:    2     // Minimal resources
QueueSize:  100   // Small buffer
```

### Medium Traffic (10-100 req/s):

```go
Workers:    5     // Recommended default
QueueSize:  1000  // 10s buffer
```

### High Traffic (> 100 req/s):

```go
Workers:    10    // High parallelism
QueueSize:  5000  // 50s buffer
```

### Burst Traffic:

```go
Workers:    15    // Very high parallelism
QueueSize:  10000 // Large buffer
```

**Note:** Monitor CPU and memory usage when increasing workers.

---

## Testing Locally

### 1. Apply Migration

```bash
cd /p/github.com/sveturs/svetu/backend
./migrator up
```

### 2. Run Tests

```bash
go test -v ./internal/proj/c2c/storage/opensearch/...
```

### 3. Manual Testing

```bash
# Create a listing via API
curl -X POST http://localhost:3000/api/v1/marketplace/listings \
  -H "Content-Type: application/json" \
  -d '{"title": "Test Async", "price": 100, "category_id": 1}'

# Check indexer status
curl http://localhost:3000/api/v1/admin/indexer/status

# Check DLQ
psql "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5433/svetubd" \
  -c "SELECT * FROM opensearch_indexing_dlq;"
```

---

## Migration from Sync to Async

### Step 1: Deploy with Async Disabled

```go
// Initially deploy without enabling async
osRepo := opensearch.NewRepository(osClient, indexName, db, searchWeights)
// DO NOT call EnableAsyncIndexing yet
```

### Step 2: Monitor Baseline

- Collect sync indexing metrics
- Establish baseline latency/throughput

### Step 3: Enable Async (Canary)

```go
// Enable on one instance only
if os.Getenv("ENABLE_ASYNC_INDEXING") == "true" {
    osRepo.EnableAsyncIndexing(db, 5, 1000)
}
```

### Step 4: Monitor & Compare

- Compare async vs sync latencies
- Check DLQ accumulation
- Monitor queue saturation

### Step 5: Full Rollout

```go
// Enable on all instances
osRepo.EnableAsyncIndexing(db, 5, 1000)
```

---

## Rollback Plan

### If Issues Occur:

1. **Disable async indexing:**

```go
osRepo.DisableAsyncIndexing()
// Falls back to sync immediately
```

2. **Restart application:**
```bash
systemctl restart svetu-backend
```

3. **Clear DLQ if needed:**
```sql
TRUNCATE opensearch_indexing_dlq;
```

4. **Re-index manually:**
```bash
python3 /p/github.com/sveturs/svetu/backend/reindex_unified.py
```

---

## Questions & Support

- **Documentation:** See `SPRINT_1.3_ASYNC_INDEXING_REPORT.md`
- **Code:** `backend/internal/proj/c2c/storage/opensearch/async_indexer.go`
- **Tests:** `backend/internal/proj/c2c/storage/opensearch/async_indexer_test.go`

---

**Last Updated:** 2025-10-31
