# Health Check System Documentation

## Overview

The Listings microservice provides a comprehensive health check system with multiple endpoints for different use cases:

- **Overall Health** (`/health`) - Full dependency checks with detailed status
- **Liveness Probe** (`/health/live`) - Minimal check for Kubernetes liveness
- **Readiness Probe** (`/health/ready`) - Service readiness for traffic
- **Startup Probe** (`/health/startup`) - Initial startup verification
- **Deep Diagnostics** (`/health/deep`) - Extended diagnostics for debugging

All health checks are production-ready with:
- ✅ Timeout handling (prevents blocking)
- ✅ Result caching (reduces load on dependencies)
- ✅ Thread-safe concurrent access
- ✅ Graceful degradation (partial outages)
- ✅ Detailed error reporting
- ✅ Prometheus-compatible metrics

---

## Endpoints

### 1. Overall Health Check

**Endpoint:** `GET /health`

**Purpose:** Comprehensive health check of all dependencies

**Response Format:**
```json
{
  "status": "healthy|degraded|unhealthy",
  "version": "0.1.0",
  "uptime": "15h23m",
  "checks": {
    "database": {
      "status": "healthy",
      "response_time_ms": 5,
      "details": "5 connections active, 3 idle, 2 in use"
    },
    "redis": {
      "status": "healthy",
      "response_time_ms": 2,
      "details": "120 hits, 15 misses, 10 total conns"
    },
    "opensearch": {
      "status": "healthy",
      "response_time_ms": 15,
      "details": "cluster status: 200 OK"
    },
    "minio": {
      "status": "healthy",
      "response_time_ms": 8,
      "details": "bucket accessible"
    }
  },
  "timestamp": "2025-11-05T18:00:00Z"
}
```

**HTTP Status Codes:**
- `200 OK` - Service is healthy or degraded (accepting traffic)
- `503 Service Unavailable` - Service is unhealthy (critical dependencies down)

**Status Definitions:**
- **healthy**: All dependencies operational
- **degraded**: Critical dependencies OK, optional dependencies down (OpenSearch, MinIO)
- **unhealthy**: Critical dependencies down (PostgreSQL, Redis)

**Use Cases:**
- Monitoring dashboards
- Load balancer health checks
- Manual service verification
- CI/CD health validation

---

### 2. Liveness Probe

**Endpoint:** `GET /health/live`

**Purpose:** Kubernetes liveness probe - checks if service is running

**Response Format:**
```json
{
  "status": "healthy",
  "timestamp": 1730826000
}
```

**HTTP Status Codes:**
- `200 OK` - Service is alive
- `503 Service Unavailable` - Service is dead (should be restarted)

**Characteristics:**
- ⚡ Ultra-fast (no dependency checks)
- 🎯 Minimal overhead
- 🔄 Never fails unless service is completely dead

**Kubernetes Configuration:**
```yaml
livenessProbe:
  httpGet:
    path: /health/live
    port: 8086
  initialDelaySeconds: 10
  periodSeconds: 30
  timeoutSeconds: 5
  failureThreshold: 3
```

---

### 3. Readiness Probe

**Endpoint:** `GET /health/ready`

**Purpose:** Kubernetes readiness probe - checks if service can accept traffic

**Response Format:**
```json
{
  "status": "ready",
  "timestamp": 1730826000
}
```

**HTTP Status Codes:**
- `200 OK` - Service is ready
- `503 Service Unavailable` - Service is not ready

**Dependencies Checked:**
- ✅ PostgreSQL (database connectivity)
- ✅ Redis (cache connectivity)

**Kubernetes Configuration:**
```yaml
readinessProbe:
  httpGet:
    path: /health/ready
    port: 8086
  initialDelaySeconds: 5
  periodSeconds: 10
  timeoutSeconds: 5
  failureThreshold: 3
```

**Behavior:**
- Service is removed from load balancer if readiness fails
- Does not restart the service (unlike liveness probe)
- Checks only critical dependencies required for serving requests

---

### 4. Startup Probe

**Endpoint:** `GET /health/startup`

**Purpose:** Kubernetes startup probe - checks service initialization

**Response Format:**
```json
{
  "status": "started",
  "timestamp": 1730826000
}
```

**HTTP Status Codes:**
- `200 OK` - Service has fully started
- `503 Service Unavailable` - Service is still starting

**Startup Checks:**
1. Grace period elapsed (5 seconds)
2. PostgreSQL connectivity verified
3. Redis connectivity verified
4. Service initialization complete

**Kubernetes Configuration:**
```yaml
startupProbe:
  httpGet:
    path: /health/startup
    port: 8086
  initialDelaySeconds: 0
  periodSeconds: 5
  timeoutSeconds: 10
  failureThreshold: 30  # 150 seconds total (5s * 30)
```

**Use Cases:**
- Slow-starting applications
- Database migration wait
- Cache warming period
- Initial index building

---

### 5. Deep Health Check

**Endpoint:** `GET /health/deep`

**Purpose:** Extended diagnostics for debugging and monitoring

**Response Format:**
```json
{
  "status": "healthy",
  "version": "0.1.0",
  "uptime": "15h23m",
  "checks": {
    "database": { ... },
    "redis": { ... },
    "opensearch": { ... },
    "minio": { ... }
  },
  "timestamp": "2025-11-05T18:00:00Z",
  "diagnostics": {
    "goroutines": 42,
    "memory_alloc_mb": 128,
    "memory_sys_mb": 256,
    "num_gc": 15,
    "recent_errors": [
      "2025-11-05T17:45:00Z: redis ping failed: connection refused",
      "2025-11-05T17:50:00Z: opensearch cluster health failed: timeout"
    ],
    "connection_pools": {
      "database": {
        "active": 5,
        "idle": 10,
        "max_open": 25,
        "wait_time_ms": 2
      },
      "redis": {
        "active": 3,
        "idle": 7,
        "max_open": 10,
        "wait_time_ms": 0
      }
    }
  }
}
```

**HTTP Status Codes:**
- `200 OK` - Service is healthy or degraded
- `503 Service Unavailable` - Service is unhealthy

**Additional Information:**
- **goroutines**: Number of active goroutines (memory leak detection)
- **memory_alloc_mb**: Current memory allocation (MB)
- **memory_sys_mb**: Total memory obtained from OS (MB)
- **num_gc**: Number of GC cycles (GC pressure)
- **recent_errors**: Last 10 health check errors
- **connection_pools**: Pool metrics for all connections

**Use Cases:**
- Performance debugging
- Memory leak investigation
- Connection pool tuning
- Error pattern analysis

**⚠️ Warning:** This endpoint is more resource-intensive. Use sparingly in production.

---

## Configuration

### Environment Variables

```bash
# Health check timeout for individual checks (default: 5s)
VONDILISTINGS_HEALTH_CHECK_TIMEOUT=5s

# Interval between cached checks (default: 30s)
VONDILISTINGS_HEALTH_CHECK_INTERVAL=30s

# Timeout for startup checks (default: 60s)
VONDILISTINGS_HEALTH_STARTUP_TIMEOUT=60s

# Duration to cache check results (default: 10s)
VONDILISTINGS_HEALTH_CACHE_DURATION=10s

# Enable deep diagnostics endpoint (default: true)
VONDILISTINGS_HEALTH_ENABLE_DEEP_CHECKS=true
```

### Default Configuration

```go
&health.Config{
    CheckTimeout:     5 * time.Second,
    CheckInterval:    30 * time.Second,
    StartupTimeout:   60 * time.Second,
    CacheDuration:    10 * time.Second,
    EnableDeepChecks: true,
}
```

---

## Dependency Checks

### 1. PostgreSQL (Critical)

**Checks Performed:**
1. Connection ping
2. Simple query execution (`SELECT COUNT(*) FROM listings`)
3. Connection pool statistics

**Timeout:** 5 seconds (configurable)

**Failure Impact:** Service marked as **unhealthy**

**Details Example:**
```
"5 connections active, 3 idle, 2 in use"
```

---

### 2. Redis (Critical)

**Checks Performed:**
1. Connection ping
2. Set/Get test operation
3. Pool statistics

**Timeout:** 5 seconds (configurable)

**Failure Impact:** Service marked as **unhealthy**

**Details Example:**
```
"120 hits, 15 misses, 10 total conns"
```

---

### 3. OpenSearch (Optional)

**Checks Performed:**
1. Cluster health API call

**Timeout:** 5 seconds (configurable)

**Failure Impact:** Service marked as **degraded** (not unhealthy)

**Details Example:**
```
"cluster status: 200 OK"
```

**Note:** Search functionality may be unavailable, but service continues operating.

---

### 4. MinIO (Optional)

**Checks Performed:**
1. Bucket existence check

**Timeout:** 5 seconds (configurable)

**Failure Impact:** Service marked as **degraded** (not unhealthy)

**Details Example:**
```
"bucket accessible"
```

**Note:** Image storage may be unavailable, but core service continues operating.

---

## Caching Strategy

To prevent overloading dependencies with health checks, results are cached:

- **Cache Duration:** 10 seconds (configurable via `VONDILISTINGS_HEALTH_CACHE_DURATION`)
- **Cache Invalidation:** Automatic expiry after duration
- **Thread Safety:** Concurrent-safe using read-write locks

**Example Behavior:**
```
T+0s:  /health called → Check DB/Redis/etc → Cache result
T+2s:  /health called → Return cached result (no dependency check)
T+5s:  /health called → Return cached result (no dependency check)
T+11s: /health called → Cache expired → Check DB/Redis/etc → Cache new result
```

**Benefits:**
- 🚀 Fast response times (< 1ms for cached results)
- 💰 Reduced load on dependencies
- 🛡️ Protection against health check storms

---

## Monitoring Integration

### Prometheus Metrics

Health check duration and failures are automatically exposed via Prometheus metrics:

```prometheus
# Health check duration histogram
http_request_duration_seconds{method="GET",path="/health",status="200"}

# Health check failure counter
http_requests_total{method="GET",path="/health",status="503"}

# Connection pool metrics
listings_db_open_connections
listings_db_in_use_connections
listings_db_idle_connections
```

### Grafana Dashboard Example

```json
{
  "title": "Listings Service Health",
  "panels": [
    {
      "title": "Health Check Status",
      "targets": [
        {
          "expr": "http_requests_total{path=\"/health\"}"
        }
      ]
    },
    {
      "title": "Database Connections",
      "targets": [
        {
          "expr": "listings_db_open_connections"
        }
      ]
    }
  ]
}
```

---

## Complete Kubernetes Deployment Example

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: listings-service
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: listings
        image: listings-service:0.1.0
        ports:
        - containerPort: 8086
          name: http

        # Startup probe - give service time to initialize
        startupProbe:
          httpGet:
            path: /health/startup
            port: http
          initialDelaySeconds: 0
          periodSeconds: 5
          timeoutSeconds: 10
          failureThreshold: 30  # 150 seconds max startup time

        # Liveness probe - restart if service is dead
        livenessProbe:
          httpGet:
            path: /health/live
            port: http
          initialDelaySeconds: 10
          periodSeconds: 30
          timeoutSeconds: 5
          failureThreshold: 3

        # Readiness probe - remove from load balancer if not ready
        readinessProbe:
          httpGet:
            path: /health/ready
            port: http
          initialDelaySeconds: 5
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 3
          successThreshold: 1

        env:
        - name: VONDILISTINGS_HEALTH_CHECK_TIMEOUT
          value: "5s"
        - name: VONDILISTINGS_HEALTH_CACHE_DURATION
          value: "10s"
        - name: VONDILISTINGS_HEALTH_ENABLE_DEEP_CHECKS
          value: "true"

        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
```

---

## Troubleshooting

### Issue: Health checks timing out

**Symptoms:**
- `/health` endpoint returns 503
- Logs show "context deadline exceeded"

**Solutions:**
1. Increase timeout: `VONDILISTINGS_HEALTH_CHECK_TIMEOUT=10s`
2. Check dependency connectivity (DB, Redis, OpenSearch, MinIO)
3. Verify network latency between service and dependencies
4. Check if dependencies are overloaded

---

### Issue: Service marked as unhealthy but working fine

**Symptoms:**
- Service processes requests successfully
- `/health` returns 503

**Solutions:**
1. Check Redis connectivity: `redis-cli -h <host> -p <port> PING`
2. Check PostgreSQL connectivity: `psql -h <host> -p <port> -U <user> -c "SELECT 1"`
3. Review recent errors: `curl http://localhost:8086/health/deep | jq .diagnostics.recent_errors`
4. Verify configuration: check `VONDILISTINGS_DB_*` and `VONDILISTINGS_REDIS_*` env vars

---

### Issue: Too many health check requests

**Symptoms:**
- High CPU usage during health checks
- Dependencies (DB/Redis) experiencing load spikes
- Slow response times

**Solutions:**
1. Increase cache duration: `VONDILISTINGS_HEALTH_CACHE_DURATION=30s`
2. Reduce health check frequency in Kubernetes probes
3. Use separate health check endpoints for different purposes:
   - Liveness: `/health/live` (lightweight)
   - Readiness: `/health/ready` (medium weight)
   - Monitoring: `/health` (full checks, less frequent)

---

### Issue: Startup probe keeps failing

**Symptoms:**
- Service never becomes ready
- Kubernetes keeps restarting the pod

**Solutions:**
1. Increase startup timeout: `VONDILISTINGS_HEALTH_STARTUP_TIMEOUT=120s`
2. Increase Kubernetes probe `failureThreshold`: 30 → 60
3. Check startup logs for initialization errors
4. Verify database migrations are completing
5. Check if dependencies are available during startup

---

### Issue: Memory leaks detected in deep health check

**Symptoms:**
- `memory_alloc_mb` continuously increasing
- `num_gc` very high

**Solutions:**
1. Enable pprof profiling: access `:6060/debug/pprof/heap`
2. Check goroutine leaks: `curl http://localhost:8086/health/deep | jq .diagnostics.goroutines`
3. Review connection pool settings (may have leaks)
4. Check for unclosed database connections or Redis clients

---

## Testing

### Manual Testing

```bash
# Test overall health
curl -i http://localhost:8086/health | jq

# Test liveness
curl -i http://localhost:8086/health/live

# Test readiness
curl -i http://localhost:8086/health/ready

# Test startup
curl -i http://localhost:8086/health/startup

# Test deep diagnostics
curl -i http://localhost:8086/health/deep | jq
```

### Automated Testing

```bash
# Run unit tests
go test -v ./internal/health/...

# Run with coverage
go test -v -cover ./internal/health/...

# Run with race detector
go test -v -race ./internal/health/...
```

### Integration Testing

```bash
# Start dependencies
docker-compose up -d postgres redis opensearch minio

# Start service
go run ./cmd/server

# Test health endpoints
./scripts/test-health-checks.sh
```

---

## Performance Characteristics

| Endpoint | Avg Response Time | Dependencies Checked | Cache Enabled |
|----------|-------------------|----------------------|---------------|
| `/health` | 10-50ms (uncached) / < 1ms (cached) | DB, Redis, OpenSearch, MinIO | Yes |
| `/health/live` | < 1ms | None | No |
| `/health/ready` | 5-20ms | DB, Redis | Yes |
| `/health/startup` | 10-30ms | DB, Redis | No |
| `/health/deep` | 20-100ms | All + Diagnostics | Partial |

**Recommendations:**
- Use `/health/live` for frequent checks (every 10s)
- Use `/health/ready` for readiness (every 10-30s)
- Use `/health` for monitoring (every 30-60s)
- Use `/health/deep` for debugging only (manual)

---

## Best Practices

1. **Use appropriate endpoints for each use case**
   - Liveness: Fast, no dependencies
   - Readiness: Critical dependencies only
   - Health: Full checks for monitoring

2. **Configure reasonable timeouts**
   - Too short: False negatives
   - Too long: Slow failure detection
   - Recommended: 5-10 seconds

3. **Enable caching in production**
   - Reduces dependency load
   - Faster response times
   - Recommended: 10-30 seconds

4. **Monitor health check failures**
   - Set up alerts for persistent failures
   - Track failure patterns
   - Investigate sudden spikes

5. **Test failure scenarios**
   - Simulate DB outages
   - Simulate Redis failures
   - Verify graceful degradation

6. **Document dependency relationships**
   - Critical vs optional dependencies
   - Impact of failures
   - Recovery procedures

---

## References

- [Kubernetes Probes Documentation](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/)
- [Health Check Patterns](https://microservices.io/patterns/observability/health-check-api.html)
- [12-Factor App - Health Checks](https://12factor.net/admin-processes)

---

**Last Updated:** 2025-11-05
**Version:** 0.1.0
