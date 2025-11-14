# Performance Testing Guide - Orders Microservice

**Version:** 1.0.0
**Last Updated:** 2025-11-14

This guide provides comprehensive instructions for running, analyzing, and interpreting performance tests for the Orders Microservice.

---

## Table of Contents

1. [Overview](#overview)
2. [Test Types](#test-types)
3. [Setup](#setup)
4. [Running Tests](#running-tests)
5. [Analyzing Results](#analyzing-results)
6. [Performance Targets](#performance-targets)
7. [Troubleshooting](#troubleshooting)
8. [CI/CD Integration](#cicd-integration)

---

## Overview

The Orders Microservice performance testing suite includes:

- **Unit Benchmarks**: Go benchmarks for individual functions and methods
- **Integration Benchmarks**: End-to-end benchmarks with database and Redis
- **Load Tests**: gRPC load testing simulating production traffic
- **Stress Tests**: Peak load and failure scenario testing

### Goals

1. **Validate Performance**: Ensure all operations meet latency and throughput targets
2. **Detect Regressions**: Identify performance degradation before deployment
3. **Capacity Planning**: Understand system limits and scaling needs
4. **Optimization**: Identify bottlenecks and optimization opportunities

---

## Test Types

### 1. Unit Benchmarks

**Location:** `tests/performance/orders_benchmarks_test.go`

**Purpose:** Measure performance of individual operations in isolation

**Tests:**
- `BenchmarkAddToCart` - Cart item addition
- `BenchmarkGetCart` - Cart retrieval
- `BenchmarkGetCartWithItems` - Cart with items retrieval
- `BenchmarkCreateOrder` - Order creation
- `BenchmarkListOrders` - Order listing with pagination

**Concurrent Tests:**
- `BenchmarkConcurrentAddToCart` - Concurrent writes
- `BenchmarkConcurrentGetCart` - Concurrent reads
- `BenchmarkConcurrentCreateOrder` - Concurrent order creation
- `BenchmarkMixedOrderOperations` - Mixed workload (70% reads, 30% writes)

### 2. Load Tests

**Location:** `load-tests/ghz-orders.sh`

**Purpose:** Simulate production load and measure system behavior

**Scenarios:**
1. **GetCart** - High read load (200 RPS, 60s)
2. **AddToCart** - Concurrent writes (100 RPS, 60s)
3. **CreateOrder** - Transaction load (50 RPS, 60s)
4. **ListOrders** - Pagination load (100 RPS, 60s)
5. **MixedWorkload** - Peak load (300 RPS, 60s)

### 3. Integration Tests

**Location:** `tests/integration/*_test.go`

**Purpose:** Validate end-to-end functionality with real dependencies

---

## Setup

### Prerequisites

```bash
# Install required tools
sudo apt-get update
sudo apt-get install -y jq bc

# Install k6 (HTTP load testing)
snap install k6

# Install ghz (gRPC load testing)
go install github.com/bojand/ghz/cmd/ghz@latest

# Verify installations
ghz --version
k6 version
```

### Start Services

```bash
# 1. Start PostgreSQL (if not running)
docker start listings_postgres || \
docker run -d --name listings_postgres \
  -e POSTGRES_PASSWORD=listings_password \
  -e POSTGRES_USER=listings_user \
  -e POSTGRES_DB=listings_db \
  -p 5432:5432 \
  postgres:15-alpine

# 2. Start Redis (if not running)
docker start listings_redis || \
docker run -d --name listings_redis \
  -p 6379:6379 \
  redis:7-alpine

# 3. Start Orders Microservice
/home/dim/.local/bin/start-listings-microservice.sh

# 4. Verify services are running
netstat -tlnp | grep -E '50052|5432|6379'
```

### Setup Test Data

```bash
# Run migrations
cd /p/github.com/sveturs/listings
make migrate

# Optional: Load test fixtures
# psql "postgres://listings_user:listings_password@localhost:5432/listings_db" < test_fixtures.sql
```

---

## Running Tests

### Quick Start (All Tests)

```bash
cd /p/github.com/sveturs/listings

# Run all performance tests
make perf-test

# Or manually:
go test -bench=. -benchmem -benchtime=10s ./tests/performance/...
```

### Unit Benchmarks

```bash
# Run all benchmarks
go test -bench=. -benchmem -benchtime=10s \
  -run=^$ \
  ./tests/performance/... \
  | tee benchmark_results.txt

# Run specific benchmark
go test -bench=BenchmarkAddToCart -benchmem -benchtime=10s \
  ./tests/performance/

# With CPU profiling
go test -bench=BenchmarkCreateOrder -cpuprofile=cpu.prof \
  ./tests/performance/

# Analyze CPU profile
go tool pprof cpu.prof
```

### Load Tests (gRPC)

```bash
cd /p/github.com/sveturs/listings/load-tests

# Run all load tests
./ghz-orders.sh

# Run specific test scenario
# (edit ghz-orders.sh and comment out unwanted tests)

# Custom configuration
export GRPC_HOST="localhost:50052"
export TEST_USER_ID="1"
export TEST_STOREFRONT_ID="1"
export TEST_LISTING_ID="1"
./ghz-orders.sh
```

### Integration Tests

```bash
# Run all integration tests with testcontainers
go test -v -timeout=10m -tags=integration \
  ./test/integration/...

# Run specific test
go test -v -timeout=10m -tags=integration \
  -run TestOrdersE2E \
  ./test/integration/
```

---

## Analyzing Results

### Benchmark Results

**Understanding Output:**
```
BenchmarkAddToCart-8    50000    25473 ns/op    1024 B/op    12 allocs/op
                   ^      ^        ^              ^           ^
                   |      |        |              |           └─ Allocations per op
                   |      |        |              └─ Bytes allocated per op
                   |      |        └─ Nanoseconds per operation
                   |      └─ Number of iterations
                   └─ Benchmark name with GOMAXPROCS
```

**Performance Evaluation:**

| Operation | Target ns/op | Good | Acceptable | Poor |
|-----------|--------------|------|------------|------|
| AddToCart | <50,000 | <30ms | 30-50ms | >50ms |
| GetCart | <20,000 | <10ms | 10-20ms | >20ms |
| CreateOrder | <200,000 | <100ms | 100-200ms | >200ms |

**Memory Analysis:**
- Low allocations per op: <1000 B/op
- High allocations per op: >10,000 B/op
- Watch for increasing allocations in concurrent tests (potential leak)

### Load Test Results

**Metrics Explained:**

```json
{
  "count": 6000,              // Total requests
  "rps": 100.45,              // Requests per second
  "average": 9847123,         // Average latency (nanoseconds)
  "fastest": 5234567,         // Fastest request
  "slowest": 45678901,        // Slowest request
  "latencyDistribution": [
    {"percentage": 50, "latency": "8.5ms"},   // p50 (median)
    {"percentage": 95, "latency": "15.2ms"},  // p95 (SLA metric)
    {"percentage": 99, "latency": "23.7ms"}   // p99 (tail latency)
  ]
}
```

**Latency Percentiles:**
- **p50 (median)**: Half of requests complete under this time
- **p95**: 95% of requests complete under this time (main SLA)
- **p99**: 99% of requests complete under this time (tail latency)

**Success Criteria:**

| Test | p50 | p95 | p99 | Error Rate |
|------|-----|-----|-----|------------|
| GetCart | <10ms | <20ms | <50ms | <1% |
| AddToCart | <20ms | <50ms | <100ms | <1% |
| CreateOrder | <100ms | <200ms | <500ms | <1% |
| ListOrders | <50ms | <100ms | <200ms | <1% |

### Generate Report

```bash
cd /p/github.com/sveturs/listings

# Generate comprehensive report
./scripts/generate_performance_report.sh \
  ./load-tests/results \
  ./PERFORMANCE_REPORT.md

# View report
cat PERFORMANCE_REPORT.md
```

---

## Performance Targets

### Latency Targets

| Operation | Target P95 | Justification |
|-----------|------------|---------------|
| **GetCart** | <20ms | Frequent operation, should be cached |
| **AddToCart** | <50ms | Write operation with validation |
| **CreateOrder** | <200ms | Complex transaction with stock check |
| **ListOrders** | <100ms | Paginated read with JOINs |
| **UpdateCart** | <50ms | Write operation |
| **GetOrder** | <50ms | Read operation with relations |

### Throughput Targets

| Operation | Target RPS | Notes |
|-----------|------------|-------|
| **GetCart** | >500 | High read volume |
| **AddToCart** | >100 | Moderate write volume |
| **CreateOrder** | >50 | Transaction-heavy |
| **ListOrders** | >100 | Read-heavy |

### Resource Limits

| Resource | Limit | Alert Threshold |
|----------|-------|-----------------|
| **CPU** | <70% | >60% for 5min |
| **Memory** | <512MB | >400MB |
| **DB Connections** | <50 | >40 |
| **Goroutines** | <5000 | >3000 |

---

## Troubleshooting

### High Latency

**Symptoms:**
- P95 latency exceeds targets
- Slow response times

**Diagnosis:**
```bash
# 1. Check database query performance
psql "postgres://..." -c "SELECT query, mean_exec_time FROM pg_stat_statements ORDER BY mean_exec_time DESC LIMIT 10;"

# 2. Profile the application
go test -bench=BenchmarkCreateOrder -cpuprofile=cpu.prof
go tool pprof cpu.prof
(pprof) top10
(pprof) list CreateOrder

# 3. Check Redis latency
redis-cli --latency
```

**Solutions:**
- Add database indexes
- Optimize queries (reduce JOINs)
- Implement caching
- Use connection pooling
- Optimize serialization

### High Error Rate

**Symptoms:**
- Error rate > 1%
- Timeouts or connection failures

**Diagnosis:**
```bash
# 1. Check error logs
tail -f /tmp/listings-microservice.log | grep ERROR

# 2. Analyze error distribution
cat load-tests/results/*.json | jq '.errorDistribution'

# 3. Check database connections
psql "postgres://..." -c "SELECT count(*) FROM pg_stat_activity;"
```

**Solutions:**
- Increase connection pool size
- Implement retry logic
- Add circuit breakers
- Optimize timeout values
- Fix race conditions

### Memory Leaks

**Symptoms:**
- Memory usage continuously increasing
- OOM kills

**Diagnosis:**
```bash
# 1. Generate heap profile
curl http://localhost:6060/debug/pprof/heap > heap.prof
go tool pprof heap.prof
(pprof) top10

# 2. Check goroutine leaks
curl http://localhost:6060/debug/pprof/goroutine > goroutine.prof
go tool pprof goroutine.prof
```

**Solutions:**
- Fix goroutine leaks (use context cancellation)
- Close database connections properly
- Release resources in defer statements
- Use object pools for frequently allocated objects

### Database Connection Pool Exhaustion

**Symptoms:**
- "too many connections" errors
- High connection wait times

**Diagnosis:**
```bash
# Check current connections
psql "postgres://..." -c "SELECT count(*), state FROM pg_stat_activity GROUP BY state;"

# Check pool stats
curl http://localhost:6060/debug/vars | jq '.db'
```

**Solutions:**
- Increase max connections in PostgreSQL
- Reduce max open connections in application
- Implement connection pooling
- Close connections promptly
- Use read replicas for read-heavy workloads

---

## CI/CD Integration

### GitHub Actions

**Workflow:** `.github/workflows/orders-service-ci.yml`

**Stages:**
1. **Lint** - Code quality checks
2. **Unit Tests** - Fast unit tests with coverage
3. **Integration Tests** - E2E tests with testcontainers
4. **Performance Tests** - Benchmark regression detection
5. **Build** - Docker image creation
6. **Deploy** - Automatic deployment to dev environment

**Performance Test Stage:**
```yaml
- name: Run benchmark tests
  run: |
    go test -bench=. -benchmem -benchtime=10s \
      -run=^$ \
      ./tests/performance/... \
      | tee benchmark_results.txt

- name: Compare with baseline
  run: |
    # Download previous benchmark
    # Compare and fail if regression > 20%
```

### Pre-commit Hooks

**Setup:**
```bash
# Install pre-commit
pip install pre-commit

# Install hooks
pre-commit install

# Run manually
pre-commit run --all-files
```

**Configuration:** `.pre-commit-config.yaml`

---

## Best Practices

### Writing Benchmarks

1. **Reset Timer**: Always call `b.ResetTimer()` after setup
2. **Report Allocations**: Use `b.ReportAllocs()` to track memory
3. **Avoid Optimization**: Use `b.N` loop variable correctly
4. **Parallel Tests**: Use `b.RunParallel()` for concurrent tests
5. **Consistent Environment**: Run on same hardware for comparisons

**Example:**
```go
func BenchmarkMyOperation(b *testing.B) {
    // Setup (not measured)
    setup := prepareTestData()

    // Reset timer before actual test
    b.ResetTimer()
    b.ReportAllocs()

    // Actual benchmark
    for i := 0; i < b.N; i++ {
        result := myOperation(setup)
        _ = result // Prevent compiler optimization
    }
}
```

### Running Load Tests

1. **Warm-up Phase**: Always include warm-up to populate caches
2. **Realistic Data**: Use production-like data volumes
3. **Gradual Ramp-up**: Increase load gradually to avoid false positives
4. **Sustained Load**: Run sustained load for sufficient duration (>60s)
5. **Cool-down**: Allow services to stabilize between tests

### Analyzing Results

1. **Look at Distributions**: Don't rely on averages alone
2. **Check Tail Latencies**: P99 matters for user experience
3. **Monitor System Resources**: CPU, memory, disk I/O
4. **Compare with Baseline**: Track performance trends
5. **Investigate Regressions**: Any increase >20% requires investigation

---

## References

- [Go Benchmarking Tutorial](https://golang.org/pkg/testing/#hdr-Benchmarks)
- [ghz Documentation](https://ghz.sh/)
- [Performance Testing Best Practices](https://www.oreilly.com/library/view/web-performance-testing/9781491948781/)
- [Prometheus Best Practices](https://prometheus.io/docs/practices/)

---

## Appendix: Quick Reference

### Common Commands

```bash
# Run all benchmarks
make perf-test

# Run specific benchmark
go test -bench=BenchmarkAddToCart -benchtime=10s ./tests/performance/

# Run load tests
cd load-tests && ./ghz-orders.sh

# Generate report
./scripts/generate_performance_report.sh

# View metrics
curl http://localhost:50052/metrics

# Check service health
grpcurl -plaintext localhost:50052 grpc.health.v1.Health/Check
```

### Performance Checklist

- [ ] All benchmarks pass (p95 within targets)
- [ ] No memory leaks detected
- [ ] Error rate < 1%
- [ ] CPU usage < 70% under normal load
- [ ] Database connection pool not exhausted
- [ ] Redis cache hit rate > 80%
- [ ] No goroutine leaks
- [ ] Performance report generated
- [ ] Results compared with baseline
- [ ] No regressions > 20%

---

**Document Owner:** Development Team
**Last Review:** 2025-11-14
**Next Review:** 2025-12-14
