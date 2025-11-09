# Listings Microservice Troubleshooting Guide

**Last Updated:** 2025-11-05
**Version:** 1.0.0
**Maintainer:** Platform Team

## Table of Contents

- [Decision Tree](#decision-tree)
- [Debugging Tools](#debugging-tools)
  - [Logs](#logs)
  - [Metrics](#metrics)
  - [Profiling](#profiling)
  - [Tracing](#tracing)
- [Common Error Messages](#common-error-messages)
- [Performance Investigation](#performance-investigation)
  - [Slow Requests](#slow-requests)
  - [Database Query Analysis](#database-query-analysis)
  - [Memory Profiling](#memory-profiling)
  - [CPU Profiling](#cpu-profiling)
- [Network Issues](#network-issues)
  - [Connectivity Testing](#connectivity-testing)
  - [DNS Resolution](#dns-resolution)
  - [Certificate Issues](#certificate-issues)
- [Component-Specific Troubleshooting](#component-specific-troubleshooting)

---

## Decision Tree

```
Service Issue Detected
│
├─ Service not responding
│  │
│  ├─ Health check failing?
│  │  ├─ YES → Check service status
│  │  │        └─ Not running? → Check logs for crash
│  │  │                         └─ See: Service Down (RUNBOOK.md)
│  │  │
│  │  └─ NO → Check port bindings
│  │           └─ Not listening? → Check firewall/iptables
│  │
│  └─ Timeouts?
│     ├─ YES → Check latency metrics
│     │        └─ High P99? → See: Performance Investigation
│     │
│     └─ NO → Check error rate
│              └─ High errors? → See: Common Error Messages
│
├─ High error rate (>1%)
│  │
│  ├─ Database errors?
│  │  ├─ Connection refused → Check PostgreSQL status
│  │  ├─ Too many connections → See: DB Connection Pool (RUNBOOK.md)
│  │  └─ Slow queries → See: Slow Queries (RUNBOOK.md)
│  │
│  ├─ Redis errors?
│  │  ├─ Connection refused → Check Redis status
│  │  ├─ OOM → Clear cache keys
│  │  └─ Timeout → Check network latency
│  │
│  └─ OpenSearch errors?
│     ├─ Cluster red → See: OpenSearch Cluster Red (RUNBOOK.md)
│     └─ Index not found → Reindex from PostgreSQL
│
├─ High latency (P99 >2s)
│  │
│  ├─ Which endpoint?
│  │  ├─ Check metrics by method
│  │  └─ Identify slow operations
│  │
│  ├─ Database slow?
│  │  ├─ Check active queries
│  │  ├─ Review query plans
│  │  └─ Check for lock contention
│  │
│  ├─ CPU high?
│  │  └─ Generate CPU profile
│  │
│  └─ Memory high?
│     └─ Generate heap profile
│
└─ Service instability
   │
   ├─ Frequent restarts?
   │  ├─ OOM killed? → Check memory usage
   │  ├─ Panic? → Check logs for stack traces
   │  └─ Config error? → Validate environment variables
   │
   └─ Intermittent failures?
      ├─ Race conditions? → Run with -race flag
      └─ Resource exhaustion? → Check file descriptors
```

---

## Debugging Tools

### Logs

#### Structured Logging

The service uses structured JSON logging with zerolog. All logs include:
- `timestamp`: ISO8601 format
- `level`: debug, info, warn, error, fatal
- `message`: Human-readable message
- `component`: Source component (service, repository, transport)
- Additional context fields

#### Accessing Logs

**Systemd Journal (Production):**
```bash
# Tail logs in real-time
sudo journalctl -u listings-service -f

# Last 100 lines
sudo journalctl -u listings-service -n 100

# Logs since timestamp
sudo journalctl -u listings-service --since "2025-11-05 10:00:00"

# Logs in time range
sudo journalctl -u listings-service --since "10:00" --until "11:00"

# Filter by log level
sudo journalctl -u listings-service | grep '"level":"error"'

# Export to file
sudo journalctl -u listings-service --since "1 hour ago" > /tmp/listings-logs.txt
```

**Parse JSON Logs:**
```bash
# Pretty-print JSON logs
sudo journalctl -u listings-service -n 100 | grep -v "^--" | jq .

# Filter by field
sudo journalctl -u listings-service --since "10 minutes ago" | \
  jq -r 'select(.component == "repository")'

# Count errors by message
sudo journalctl -u listings-service --since "1 hour ago" | \
  jq -r 'select(.level == "error") | .message' | \
  sort | uniq -c | sort -rn

# Extract slow queries
sudo journalctl -u listings-service --since "30 minutes ago" | \
  jq -r 'select(.query_duration > 1000) | "\(.timestamp) \(.query)"'
```

**Docker Logs (Development):**
```bash
# Tail service logs
docker logs -f listings_app

# Last 100 lines
docker logs listings_app --tail 100

# Since timestamp
docker logs listings_app --since 2025-11-05T10:00:00
```

#### Log Levels

| Level | When to Use | Example |
|-------|-------------|---------|
| **debug** | Development debugging, verbose output | "Rate limit check passed" |
| **info** | Normal operations, business events | "Listing created successfully" |
| **warn** | Recoverable errors, degraded performance | "Rate limit exceeded" |
| **error** | Request failures, requires investigation | "Database connection failed" |
| **fatal** | Service cannot continue, immediate shutdown | "Failed to load configuration" |

#### Useful Log Queries

**Find errors in last hour:**
```bash
sudo journalctl -u listings-service --since "1 hour ago" | \
  jq -r 'select(.level == "error") | "\(.timestamp) [\(.component)] \(.message)"'
```

**Track request by ID:**
```bash
sudo journalctl -u listings-service | \
  jq -r 'select(.request_id == "abc123")'
```

**Find slow database queries:**
```bash
sudo journalctl -u listings-service --since "30 minutes ago" | \
  jq -r 'select(.query_duration_ms > 1000) | "\(.query_duration_ms)ms - \(.query)"' | \
  sort -rn | head -20
```

**Count requests by endpoint:**
```bash
sudo journalctl -u listings-service --since "1 hour ago" | \
  jq -r 'select(.method != null) | .method' | \
  sort | uniq -c | sort -rn
```

### Metrics

#### Prometheus Metrics Endpoint

All metrics are exposed at: `http://localhost:8086/metrics`

**Accessing Metrics:**
```bash
# All metrics
curl -s http://localhost:8086/metrics

# Filter by prefix
curl -s http://localhost:8086/metrics | grep listings_

# Specific metric
curl -s http://localhost:8086/metrics | grep listings_grpc_requests_total

# Export to file
curl -s http://localhost:8086/metrics > /tmp/metrics-$(date +%s).txt
```

#### Key Metric Categories

**1. gRPC Handler Metrics**
```bash
# Request count by method and status
curl -s http://localhost:8086/metrics | grep listings_grpc_requests_total

# Request duration (histogram)
curl -s http://localhost:8086/metrics | grep listings_grpc_request_duration_seconds

# Active requests
curl -s http://localhost:8086/metrics | grep listings_grpc_handler_requests_active
```

**2. Database Metrics**
```bash
# Connection pool status
curl -s http://localhost:8086/metrics | grep listings_db_connections

# Query duration by operation
curl -s http://localhost:8086/metrics | grep listings_db_query_duration_seconds
```

**3. Rate Limiting Metrics**
```bash
# Rate limit evaluations
curl -s http://localhost:8086/metrics | grep listings_rate_limit_hits_total

# Rejected requests
curl -s http://localhost:8086/metrics | grep listings_rate_limit_rejected_total

# Rejection rate (requires calculation)
curl -s http://localhost:8086/metrics | \
  awk '/listings_rate_limit_rejected_total/{rejected=$2}
       /listings_rate_limit_hits_total/{hits=$2}
       END{print "Rejection rate: " (rejected/hits)*100 "%"}'
```

**4. Timeout Metrics**
```bash
# Timeout occurrences
curl -s http://localhost:8086/metrics | grep listings_timeouts_total

# Near-timeout warnings
curl -s http://localhost:8086/metrics | grep listings_near_timeouts_total
```

**5. Business Metrics**
```bash
# Listings operations
curl -s http://localhost:8086/metrics | grep -E 'listings_listings_(created|updated|deleted)_total'

# Search queries
curl -s http://localhost:8086/metrics | grep listings_listings_searched_total

# Product views
curl -s http://localhost:8086/metrics | grep listings_inventory_product_views_total
```

#### Metrics Analysis

**Calculate Error Rate:**
```bash
# Error rate = (errors / total requests) * 100
curl -s http://localhost:8086/metrics | awk '
  /listings_grpc_requests_total.*status="[45]/ {errors+=$2}
  /listings_grpc_requests_total/ {total+=$2}
  END {print "Error Rate: " (errors/total)*100 "%"}'
```

**Calculate P95 Latency (approximate):**
```bash
# Extract histogram buckets and calculate
curl -s http://localhost:8086/metrics | \
  grep 'listings_grpc_request_duration_seconds_bucket' | \
  sort -t= -k2 -n
```

**Top 5 Slowest Endpoints:**
```bash
curl -s http://localhost:8086/metrics | \
  grep 'listings_grpc_request_duration_seconds_sum' | \
  sed 's/.*method="\([^"]*\)".* \(.*\)/\2 \1/' | \
  sort -rn | head -5
```

**Database Connection Pool Utilization:**
```bash
curl -s http://localhost:8086/metrics | awk '
  /listings_db_connections_open/ {open=$2}
  /listings_db_connections_idle/ {idle=$2}
  END {print "Active: " (open-idle) ", Idle: " idle ", Total: " open}'
```

### Profiling

#### CPU Profiling

**Capture CPU Profile:**
```bash
# 30-second CPU profile
curl http://localhost:8086/debug/pprof/profile?seconds=30 > /tmp/cpu-$(date +%s).prof

# During high CPU, capture immediately
top -b -n 1 -p $(pgrep listings-service)
curl http://localhost:8086/debug/pprof/profile?seconds=10 > /tmp/cpu-urgent.prof
```

**Analyze CPU Profile:**
```bash
# Top functions by CPU time
go tool pprof -top /tmp/cpu-*.prof

# Interactive mode
go tool pprof /tmp/cpu-*.prof
# Commands: top, list <function>, web

# Generate visualization (requires graphviz)
go tool pprof -pdf /tmp/cpu-*.prof > /tmp/cpu-profile.pdf

# Find hot path
go tool pprof -list=main /tmp/cpu-*.prof
```

**CPU Profile Output Example:**
```
Type: cpu
Showing nodes accounting for 2.5s, 83.33% of 3s total
      flat  flat%   sum%        cum   cum%
     1.2s 40.00% 40.00%      1.8s 60.00%  runtime.scanobject
     0.5s 16.67% 56.67%      0.5s 16.67%  runtime.memmove
     0.4s 13.33% 70.00%      0.6s 20.00%  database/sql.(*DB).query
     0.4s 13.33% 83.33%      0.4s 13.33%  encoding/json.(*encodeState).string
```

#### Memory Profiling

**Capture Heap Profile:**
```bash
# Current heap snapshot
curl http://localhost:8086/debug/pprof/heap > /tmp/heap-$(date +%s).prof

# Capture series for comparison
for i in {1..5}; do
  curl http://localhost:8086/debug/pprof/heap > /tmp/heap-$i.prof
  echo "Captured profile $i"
  sleep 60
done
```

**Analyze Memory Profile:**
```bash
# Top memory consumers
go tool pprof -top /tmp/heap-*.prof

# Compare profiles (identify leaks)
go tool pprof -base /tmp/heap-1.prof /tmp/heap-5.prof

# Show allocations
go tool pprof -alloc_space /tmp/heap-*.prof

# Interactive analysis
go tool pprof /tmp/heap-*.prof
# Commands: top, list <function>, web
```

**Memory Profile Output Example:**
```
Type: inuse_space
Showing nodes accounting for 512MB, 89.47% of 572MB total
      flat  flat%   sum%        cum   cum%
   256MB 44.76% 44.76%    256MB 44.76%  database/sql.(*DB).connectionOpener
   128MB 22.38% 67.14%    128MB 22.38%  net/http.(*persistConn).readLoop
    64MB 11.19% 78.32%     64MB 11.19%  github.com/redis/go-redis/v9.(*Client).newConn
    64MB 11.19% 89.51%     64MB 11.19%  bufio.NewReaderSize
```

#### Goroutine Profiling

**Check Goroutine Count:**
```bash
# Current goroutine count
curl -s http://localhost:8086/debug/pprof/goroutine?debug=1 | grep "goroutine profile" | awk '{print $4}'

# Detailed goroutine dump
curl http://localhost:8086/debug/pprof/goroutine?debug=2 > /tmp/goroutines.txt

# Count by state
curl -s http://localhost:8086/debug/pprof/goroutine?debug=1 | \
  grep -oP '\[.*?\]' | sort | uniq -c | sort -rn
```

**Analyze Goroutines:**
```bash
# Find goroutine leaks
curl http://localhost:8086/debug/pprof/goroutine?debug=1 | \
  grep -oP '# 0x[0-9a-f]+ \K.*' | \
  sort | uniq -c | sort -rn | head -20

# Interactive analysis
go tool pprof http://localhost:8086/debug/pprof/goroutine
```

**Healthy Goroutine Count:**
- **Baseline:** 10-30 (workers, metrics collector, server)
- **Under load:** 50-200
- **Warning:** >500 (possible leak)
- **Critical:** >1000 (investigate immediately)

#### Block Profiling

**Enable and Capture:**
```bash
# Block profile (contention profiling)
curl http://localhost:8086/debug/pprof/block > /tmp/block-$(date +%s).prof

# Analyze
go tool pprof -top /tmp/block-*.prof
```

#### Mutex Profiling

**Enable and Capture:**
```bash
# Mutex profile
curl http://localhost:8086/debug/pprof/mutex > /tmp/mutex-$(date +%s).prof

# Analyze
go tool pprof -top /tmp/mutex-*.prof
```

### Tracing

**Note:** Distributed tracing with Jaeger is currently disabled by default (`SVETULISTINGS_TRACING_ENABLED=false`).

#### Enable Tracing (Development)

```bash
# In .env file
SVETULISTINGS_TRACING_ENABLED=true
SVETULISTINGS_JAEGER_ENDPOINT=http://localhost:14268/api/traces

# Restart service
sudo systemctl restart listings-service
```

#### Jaeger UI

Access: `http://localhost:16686`

**Useful Queries:**
- Find slow traces: Duration > 2s
- Error traces: Tags: error=true
- Specific operation: Operation: `/listings.v1.ListingsService/GetListing`

---

## Common Error Messages

### Database Errors

#### "pq: sorry, too many clients already"

**Cause:** PostgreSQL connection limit reached.

**Investigation:**
```bash
# Check active connections
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT count(*), state FROM pg_stat_activity WHERE datname='listings_db' GROUP BY state;"

# Check connection pool
curl -s http://localhost:8086/metrics | grep listings_db_connections
```

**Solution:**
```bash
# Kill idle connections
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT pg_terminate_backend(pid) FROM pg_stat_activity
   WHERE datname='listings_db' AND state='idle'
   AND now() - state_change > interval '5 minutes';"

# Restart service
sudo systemctl restart listings-service
```

**Prevention:** See RUNBOOK.md → Database Connection Pool Exhausted

---

#### "pq: canceling statement due to user request"

**Cause:** Query cancelled due to timeout or context cancellation.

**Investigation:**
```bash
# Check timeout metrics
curl -s http://localhost:8086/metrics | grep listings_timeouts_total

# Check slow queries
sudo journalctl -u listings-service --since "10 minutes ago" | \
  jq -r 'select(.query_duration_ms > 1000)'
```

**Solution:**
- If legitimate slow query: Optimize query or add index
- If timeout too aggressive: Review timeout configuration in `internal/timeout/config.go`

---

#### "pq: deadlock detected"

**Cause:** Multiple transactions waiting for each other.

**Investigation:**
```bash
# Check for lock contention
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT pid, state, query FROM pg_stat_activity
   WHERE datname='listings_db' AND wait_event_type='Lock';"
```

**Solution:**
```bash
# Kill blocking transaction
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT pg_terminate_backend(pid) FROM pg_stat_activity
   WHERE datname='listings_db' AND state='active'
   AND now() - query_start > interval '30 seconds';"
```

**Prevention:** Review transaction isolation levels and locking patterns

---

### Redis Errors

#### "redis: connection refused"

**Cause:** Redis service not running or not accessible.

**Investigation:**
```bash
# Check Redis status
sudo systemctl status redis

# Test connectivity
redis-cli -h localhost -p 36380 -a redis_password ping

# Check port binding
netstat -tlnp | grep 36380
```

**Solution:**
```bash
# Start Redis
sudo systemctl start redis

# Verify connection
redis-cli -h localhost -p 36380 -a redis_password ping
# Expected: PONG
```

**Impact:** Service continues (fail-open mode), but rate limiting and caching disabled.

---

#### "redis: OOM command not allowed when used memory > 'maxmemory'"

**Cause:** Redis reached memory limit.

**Investigation:**
```bash
# Check memory usage
redis-cli -h localhost -p 36380 -a redis_password info memory | grep used_memory_human

# Check maxmemory setting
redis-cli -h localhost -p 36380 -a redis_password config get maxmemory
```

**Solution:**
```bash
# Clear cache keys
redis-cli -h localhost -p 36380 -a redis_password --scan --pattern 'cache:*' | \
  xargs redis-cli -h localhost -p 36380 -a redis_password DEL

# Clear rate limit keys
redis-cli -h localhost -p 36380 -a redis_password --scan --pattern 'rate_limit:*' | \
  xargs redis-cli -h localhost -p 36380 -a redis_password DEL

# Increase maxmemory (if system has capacity)
redis-cli -h localhost -p 36380 -a redis_password config set maxmemory 2gb
```

---

### OpenSearch Errors

#### "cluster_block_exception: blocked by: [FORBIDDEN/12/index read-only / allow delete (api)]"

**Cause:** OpenSearch index set to read-only due to disk space.

**Investigation:**
```bash
# Check disk usage
curl -u admin:admin http://localhost:9200/_cat/allocation?v

# Check cluster settings
curl -u admin:admin http://localhost:9200/_cluster/settings?pretty
```

**Solution:**
```bash
# Remove read-only block
curl -u admin:admin -X PUT "http://localhost:9200/_all/_settings" \
  -H 'Content-Type: application/json' -d'
{
  "index.blocks.read_only_allow_delete": null
}'

# Free disk space (see RUNBOOK.md → Disk Space Critical)
```

---

#### "no such index [listings_microservice]"

**Cause:** OpenSearch index not created or deleted.

**Solution:**
```bash
# Create index
cd /p/github.com/sveturs/listings
python3 scripts/create_opensearch_index.py

# Reindex data
python3 scripts/reindex_via_docker.py --target-password admin

# Validate
python3 scripts/validate_opensearch.py --target-password admin
```

---

### gRPC Errors

#### "rpc error: code = ResourceExhausted desc = rate limit exceeded"

**Cause:** Client exceeded rate limit for endpoint.

**Investigation:**
```bash
# Check rate limit rejections
curl -s http://localhost:8086/metrics | grep listings_rate_limit_rejected_total

# Identify rate-limited IPs
sudo journalctl -u listings-service --since "10 minutes ago" | \
  grep "rate limit exceeded" | \
  jq -r '.identifier' | sort | uniq -c | sort -rn
```

**Solution:**
- If legitimate traffic: Increase rate limit in `internal/ratelimit/config.go`
- If abuse: Block IP (see RUNBOOK.md → Rate Limit Abuse)

---

#### "rpc error: code = DeadlineExceeded desc = context deadline exceeded"

**Cause:** Request exceeded timeout limit.

**Investigation:**
```bash
# Check timeout metrics
curl -s http://localhost:8086/metrics | grep listings_timeouts_total

# Check endpoint latency
curl -s http://localhost:8086/metrics | grep listings_grpc_request_duration_seconds_sum
```

**Solution:**
- Optimize slow operation
- Increase timeout if legitimate (edit `internal/timeout/config.go`)

---

#### "rpc error: code = Unavailable desc = connection refused"

**Cause:** Service not listening on gRPC port.

**Investigation:**
```bash
# Check service status
sudo systemctl status listings-service

# Check port binding
netstat -tlnp | grep 50053
```

**Solution:**
```bash
# Restart service
sudo systemctl restart listings-service

# Verify gRPC health
grpcurl -plaintext localhost:50053 grpc.health.v1.Health/Check
```

---

### Application Errors

#### "failed to load configuration"

**Cause:** Missing or invalid environment variables.

**Investigation:**
```bash
# Check service logs
sudo journalctl -u listings-service -n 100 | grep configuration

# Verify environment file
cat /opt/listings-dev/.env

# Check required variables
env | grep SVETULISTINGS_
```

**Solution:**
```bash
# Validate .env file
cd /opt/listings-dev
cp .env.example .env.new
# Edit .env.new with correct values
sudo mv .env.new .env

# Restart service
sudo systemctl restart listings-service
```

---

#### "panic: runtime error: invalid memory address or nil pointer dereference"

**Cause:** Nil pointer access (bug in code).

**Investigation:**
```bash
# Get full stack trace
sudo journalctl -u listings-service | grep -A 50 "panic:"

# Identify source line
sudo journalctl -u listings-service | grep "panic:" | grep -oP '\.go:\d+'
```

**Solution:**
- Emergency: Restart service immediately
- Long-term: File bug report with stack trace
- Hotfix: Deploy fix as soon as possible

---

## Performance Investigation

### Slow Requests

#### Step 1: Identify Slow Endpoint

```bash
# Check P99 latency by method
curl -s http://localhost:8086/metrics | \
  grep 'listings_grpc_request_duration_seconds' | \
  grep 'quantile="0.99"' | \
  sort -t= -k2 -rn

# Check request counts
curl -s http://localhost:8086/metrics | \
  grep 'listings_grpc_requests_total' | \
  grep -v "#"
```

#### Step 2: Analyze Request Flow

**Check Component Latency:**
```bash
# Database query duration
curl -s http://localhost:8086/metrics | grep listings_db_query_duration_seconds

# Cache hit rate
curl -s http://localhost:8086/metrics | awk '
  /listings_cache_hits_total/ {hits+=$2}
  /listings_cache_misses_total/ {misses+=$2}
  END {print "Hit Rate: " (hits/(hits+misses))*100 "%"}'
```

**Check for Bottlenecks:**
```bash
# Active database connections
curl -s http://localhost:8086/metrics | awk '
  /listings_db_connections_open/ {open=$2}
  /listings_db_connections_idle/ {idle=$2}
  END {print "Active: " (open-idle) " / " open}'

# Near-timeout requests (>80% of limit)
curl -s http://localhost:8086/metrics | grep listings_near_timeouts_total
```

#### Step 3: Profile During Load

```bash
# Capture CPU profile during slow period
curl http://localhost:8086/debug/pprof/profile?seconds=30 > /tmp/cpu-slow.prof

# Analyze
go tool pprof -top /tmp/cpu-slow.prof
```

### Database Query Analysis

#### Identify Slow Queries

**From Service Logs:**
```bash
# Find queries >1 second
sudo journalctl -u listings-service --since "30 minutes ago" | \
  jq -r 'select(.query_duration_ms > 1000) |
         "\(.query_duration_ms)ms - \(.query)"' | \
  sort -rn | head -20
```

**From PostgreSQL:**
```bash
# Active slow queries
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "SELECT pid, now() - pg_stat_activity.query_start AS duration, query
   FROM pg_stat_activity
   WHERE state = 'active' AND datname = 'listings_db'
   AND now() - pg_stat_activity.query_start > interval '1 second'
   ORDER BY duration DESC;"
```

#### Analyze Query Plan

```bash
# Get query plan
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c \
  "EXPLAIN ANALYZE SELECT * FROM listings WHERE category_id = 1301 LIMIT 10;"
```

**Interpreting EXPLAIN Output:**
- **Seq Scan:** Full table scan (BAD for large tables)
- **Index Scan:** Using index (GOOD)
- **Bitmap Heap Scan:** Index + heap (OK)
- **Actual time:** Real execution time
- **Rows:** Estimated vs actual rows

#### Common Query Optimizations

**Missing Index:**
```sql
-- Problem: Seq Scan on listings (cost=0.00..431.24)
-- Solution: Create index
CREATE INDEX CONCURRENTLY idx_listings_category_id ON listings(category_id);
```

**Inefficient WHERE Clause:**
```sql
-- Problem: WHERE LOWER(title) LIKE '%search%'
-- Solution: Use trigram index for partial matching
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX idx_listings_title_trgm ON listings USING gin(title gin_trgm_ops);
```

**N+1 Query Problem:**
```sql
-- Problem: Loading images in loop
-- Solution: Use JOIN or batch loading
SELECT l.*, array_agg(li.*) AS images
FROM listings l
LEFT JOIN listing_images li ON li.listing_id = l.id
WHERE l.id IN (1, 2, 3, 4, 5)
GROUP BY l.id;
```

### Memory Profiling

#### Detect Memory Leak

**Monitor Memory Growth:**
```bash
# Capture heap profiles every 5 minutes
for i in {1..12}; do
  curl http://localhost:8086/debug/pprof/heap > /tmp/heap-$i.prof
  SIZE=$(ps aux | grep listings-service | awk '{print $6/1024}')
  echo "$i: ${SIZE}MB at $(date)"
  sleep 300
done
```

**Compare Profiles:**
```bash
# Identify growing allocations
go tool pprof -base /tmp/heap-1.prof /tmp/heap-12.prof -top

# Interactive comparison
go tool pprof -base /tmp/heap-1.prof /tmp/heap-12.prof
# Commands: top, list <function>, web
```

#### Common Memory Leaks

**1. Unclosed HTTP Response Bodies**
```bash
# Check for http.Response leaks
go tool pprof -list=http /tmp/heap-*.prof
```

**2. Growing Maps/Slices**
```bash
# Check for map/slice growth
go tool pprof -alloc_space /tmp/heap-*.prof | grep -E 'map|slice'
```

**3. Goroutine Leaks**
```bash
# Correlate goroutines with memory
GOROUTINES=$(curl -s http://localhost:8086/debug/pprof/goroutine?debug=1 | grep "goroutine profile" | awk '{print $4}')
MEMORY=$(ps aux | grep listings-service | awk '{print $6/1024}')
echo "Goroutines: $GOROUTINES, Memory: ${MEMORY}MB"
```

### CPU Profiling

#### Capture CPU Profile

```bash
# During high CPU usage
CPU_USAGE=$(top -b -n 1 -p $(pgrep listings-service) | tail -1 | awk '{print $9}')
echo "Current CPU: ${CPU_USAGE}%"

if (( $(echo "$CPU_USAGE > 50" | bc -l) )); then
  echo "High CPU detected, capturing profile..."
  curl http://localhost:8086/debug/pprof/profile?seconds=30 > /tmp/cpu-high-$(date +%s).prof
fi
```

#### Analyze CPU Hotspots

```bash
# Top CPU consumers
go tool pprof -top /tmp/cpu-high-*.prof

# Cumulative view (includes callees)
go tool pprof -cum -top /tmp/cpu-high-*.prof

# Flamegraph (requires go-torch)
go-torch /tmp/cpu-high-*.prof
```

**Interpreting CPU Profile:**
- **flat:** Time spent in function itself
- **cum:** Time spent in function + callees
- Focus on high `flat%` for optimization

#### Common CPU Issues

**1. JSON Encoding/Decoding:**
```bash
# Check for excessive JSON operations
go tool pprof -list=encoding/json /tmp/cpu-high-*.prof
```

**2. Database Query Overhead:**
```bash
# Check for database/sql CPU usage
go tool pprof -list=database/sql /tmp/cpu-high-*.prof
```

**3. Regular Expression Matching:**
```bash
# Check for regexp CPU usage
go tool pprof -list=regexp /tmp/cpu-high-*.prof
```

---

## Network Issues

### Connectivity Testing

#### PostgreSQL

```bash
# TCP connectivity
telnet localhost 35433

# psql connection test
psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c "SELECT 1;"

# Connection timing
time psql "postgres://listings_user:listings_password@localhost:35433/listings_db" -c "SELECT 1;"

# Check network stats
netstat -an | grep 35433 | wc -l
```

#### Redis

```bash
# TCP connectivity
telnet localhost 36380

# Redis ping
redis-cli -h localhost -p 36380 -a redis_password ping

# Connection timing
time redis-cli -h localhost -p 36380 -a redis_password ping

# Check latency
redis-cli -h localhost -p 36380 -a redis_password --latency
```

#### OpenSearch

```bash
# TCP connectivity
telnet localhost 9200

# HTTP request
curl -u admin:admin http://localhost:9200/_cluster/health

# Connection timing
time curl -u admin:admin http://localhost:9200/_cluster/health

# Check response time
curl -u admin:admin -w "@-" -o /dev/null -s http://localhost:9200/_cluster/health <<'EOF'
\ntime_total: %{time_total}s\n
EOF
```

#### gRPC

```bash
# Health check
grpcurl -plaintext localhost:50053 grpc.health.v1.Health/Check

# Connection timing
time grpcurl -plaintext localhost:50053 grpc.health.v1.Health/Check

# Test endpoint
grpcurl -plaintext -d '{"id": 1}' localhost:50053 listings.v1.ListingsService/GetListing
```

### DNS Resolution

```bash
# Resolve hostname
nslookup localhost
dig localhost

# Check /etc/hosts
cat /etc/hosts | grep localhost

# Test resolution timing
time nslookup localhost
```

### Certificate Issues

**Note:** Service uses plaintext gRPC in current configuration. If TLS is enabled:

```bash
# Check certificate validity
openssl s_client -connect localhost:50053 -servername listings.local

# View certificate
openssl s_client -connect localhost:50053 < /dev/null | openssl x509 -noout -text

# Check certificate expiry
openssl s_client -connect localhost:50053 < /dev/null 2>/dev/null | \
  openssl x509 -noout -enddate
```

---

## Component-Specific Troubleshooting

### gRPC Server

**Issue:** gRPC server not starting

```bash
# Check port binding
netstat -tlnp | grep 50053

# Check for port conflicts
sudo lsof -i :50053

# Verify gRPC port in config
env | grep SVETULISTINGS_GRPC_PORT

# Test gRPC manually
grpcurl -plaintext localhost:50053 list
```

### HTTP Server

**Issue:** HTTP endpoints not responding

```bash
# Check port binding
netstat -tlnp | grep 8086

# Test health endpoint
curl -v http://localhost:8086/health

# Check HTTP metrics
curl http://localhost:8086/metrics | grep listings_http_

# Test API endpoint
curl http://localhost:8086/api/v1/listings?limit=5
```

### Rate Limiter

**Issue:** Rate limiting not working

```bash
# Check Redis connectivity
redis-cli -h localhost -p 36380 -a redis_password ping

# Check rate limit keys
redis-cli -h localhost -p 36380 -a redis_password --scan --pattern 'rate_limit:*'

# Check rate limit metrics
curl -s http://localhost:8086/metrics | grep listings_rate_limit

# Test rate limit manually
for i in {1..250}; do
  grpcurl -plaintext -d '{"id": 1}' localhost:50053 listings.v1.ListingsService/GetListing 2>&1 | \
    grep -q "ResourceExhausted" && echo "BLOCKED at $i" || echo "OK $i"
done
```

### Worker (Async Indexing)

**Issue:** OpenSearch indexing not happening

```bash
# Check worker status in logs
sudo journalctl -u listings-service | grep worker

# Check indexing queue size
curl -s http://localhost:8086/metrics | grep listings_indexing_queue_size

# Check indexing job metrics
curl -s http://localhost:8086/metrics | grep listings_indexing_jobs_processed_total

# Verify worker enabled
env | grep SVETULISTINGS_WORKER_ENABLED
```

---

## Additional Resources

- **RUNBOOK.md:** Common incidents and resolution procedures
- **DISASTER_RECOVERY.md:** Emergency recovery procedures
- **MONITORING_GUIDE.md:** Grafana dashboard usage
- **ON_CALL_GUIDE.md:** On-call engineer handbook

---

**Document Version:** 1.0.0
**Last Reviewed:** 2025-11-05
**Next Review:** 2025-12-05
**Owner:** Platform Team
