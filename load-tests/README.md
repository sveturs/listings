# Load Testing Suite - Listings Microservice

Production-ready load testing setup for validating performance and scalability of the Listings microservice under various load scenarios.

## 📋 Overview

This suite provides comprehensive load testing for both HTTP and gRPC endpoints with:

- **HTTP Load Testing** - k6-based tests for REST API endpoints
- **gRPC Load Testing** - ghz-based tests for gRPC service methods
- **System Monitoring** - Real-time CPU, memory, and load tracking
- **Automated Reporting** - Detailed performance metrics and analysis

## 🎯 Success Criteria

All tests must meet these criteria:

| Metric | Target | Description |
|--------|--------|-------------|
| **p95 Latency** | < 100ms | 95% of requests complete under 100ms |
| **Error Rate** | < 1% | Less than 1% of requests fail |
| **Throughput** | 100 RPS | Sustain 100 requests/second without degradation |
| **Memory** | No leaks | Memory usage remains stable |

## 🚀 Quick Start

### Prerequisites

```bash
# Install k6 (for HTTP tests)
snap install k6

# Install ghz (for gRPC tests)
go install github.com/bojand/ghz/cmd/ghz@latest

# Install jq (for result parsing)
sudo apt-get install jq bc
```

### Running Tests

```bash
cd /p/github.com/sveturs/listings/load-tests

# Run all tests (HTTP + gRPC)
./run-all-tests.sh

# Run only HTTP tests
./run-all-tests.sh --http-only

# Run only gRPC tests
./run-all-tests.sh --grpc-only

# Skip monitoring
./run-all-tests.sh --no-monitor
```

## 📊 Test Scenarios

### HTTP Load Test (k6-http.js)

Tests REST API endpoints with realistic traffic patterns:

| Phase | Duration | Target RPS | Purpose |
|-------|----------|------------|---------|
| **Warmup** | 30s | 10 | Initialize connections, warm caches |
| **Ramp-up** | 1m | 10→100 | Gradually increase load |
| **Sustained** | 2m | 100 | Test stability under normal load |
| **Peak** | 1m | 200 | Test maximum capacity |
| **Cool-down** | 30s | 200→0 | Graceful shutdown |

**Total Duration:** ~5 minutes

**Endpoints Tested:**
- `GET /health` - Health check endpoint
- `GET /api/v1/storefronts` - List storefronts (paginated)
- `GET /api/v1/storefronts/{id}` - Get specific storefront
- `GET /api/v1/listings` - List listings (filtered)
- `GET /api/v1/listings/{id}` - Get specific listing

**Traffic Distribution:**
- 30% Health checks
- 35% Storefront API calls
- 35% Listings API calls

### gRPC Load Test (ghz-grpc.sh)

Tests gRPC service methods with multiple scenarios:

#### Scenario 1: GetAllCategories
- **Purpose:** Test cached read performance
- **Load:** 50 RPS for 60s
- **Expected:** Very fast (< 50ms p95) due to caching

#### Scenario 2: ListStorefronts
- **Purpose:** Test paginated list queries
- **Load:** 50 RPS for 60s
- **Expected:** Moderate latency (< 150ms p95)

#### Scenario 3: GetListing
- **Purpose:** Test single item retrieval
- **Load:** 100 RPS for 60s
- **Expected:** Fast (< 100ms p95)

#### Scenario 4: Mixed Workload (Stress Test)
- **Purpose:** Test peak capacity
- **Load:** 200 RPS for 60s
- **Expected:** Higher latency but < 1% errors

**Total Duration:** ~6 minutes (4 scenarios + warmup)

## 📁 File Structure

```
load-tests/
├── k6-http.js              # HTTP load test (k6)
├── ghz-grpc.sh             # gRPC load test (ghz)
├── run-all-tests.sh        # Orchestration script
├── README.md               # This file
└── results/                # Test results (generated)
    ├── k6_results_*.json
    ├── get_all_categories_*.json
    ├── list_storefronts_*.json
    ├── get_listing_*.json
    ├── mixed_workload_*.json
    ├── system_metrics_*.log
    └── summary_*.txt
```

## 🔧 Configuration

### Environment Variables

```bash
# HTTP endpoint
export HTTP_BASE_URL="http://localhost:8086"

# gRPC endpoint
export GRPC_HOST="localhost:50051"

# Results directory
export RESULTS_DIR="/p/github.com/sveturs/listings/load-tests/results"

# Test data
export STOREFRONT_ID="1"
export PRODUCT_ID="1"
export CATEGORY_ID="1"
```

### Service Configuration

Before running tests, ensure the service is running:

```bash
cd /p/github.com/sveturs/listings

# Start the service
make run

# Or with Docker
docker-compose up -d

# Verify services are available
curl http://localhost:8086/health
grpcurl -plaintext localhost:50051 list
```

## 📈 Analyzing Results

### View HTTP Results

```bash
# Pretty-print JSON results
cat results/k6_results_*.json | jq

# View specific metrics
cat results/k6_results_*.json | jq '.metrics'

# Check thresholds
cat results/k6_results_*.json | jq '.root_group.checks'
```

### View gRPC Results

```bash
# View all gRPC results
ls -lh results/*_grpc_*.json

# Analyze specific scenario
cat results/get_all_categories_*.json | jq '.latencyDistribution'

# Compare scenarios
for f in results/*_2025*.json; do
    echo "$f: $(jq -r '.average' $f)ms avg"
done
```

### View System Metrics

```bash
# View CPU/Memory during test
cat results/system_metrics_*.log | column -t -s,

# Plot CPU usage (requires gnuplot)
gnuplot -e "set terminal dumb; set datafile separator ','; plot 'results/system_metrics_*.log' using 1:2 with lines title 'CPU %'"
```

### Generate Reports

```bash
# View test summary
cat results/summary_*.txt

# Compare multiple test runs
for f in results/summary_*.txt; do
    echo "=== $(basename $f) ==="
    grep -A 5 "HTTP Tests:" $f
    echo ""
done
```

## 🎯 Success Criteria Evaluation

Each test automatically evaluates against success criteria:

```
✅ Success Criteria Evaluation:
--------------------------------
✓ p95 latency < 100ms: 87.32ms
✓ No errors (0% error rate)

🎉 All success criteria passed for GetAllCategories!
```

### Understanding Metrics

**Latency Percentiles:**
- **p50 (median):** 50% of requests complete under this time
- **p95:** 95% of requests complete under this time (main SLA metric)
- **p99:** 99% of requests complete under this time

**Response Time Targets:**
- Cached reads (categories): < 50ms p95
- Database reads (listings): < 100ms p95
- List queries: < 150ms p95

**RPS (Requests Per Second):**
- Normal load: 50-100 RPS
- Peak load: 200 RPS
- Emergency capacity: 300+ RPS (short bursts)

## 🐛 Troubleshooting

### Tests Fail to Start

```bash
# Check if service is running
curl http://localhost:8086/health
nc -zv localhost 50051

# Check for port conflicts
netstat -tlnp | grep -E '8086|50051'

# View service logs
docker-compose logs -f listings
```

### High Error Rates

```bash
# Check database connections
psql "postgres://user:pass@localhost:5432/listings_db" -c "SELECT count(*) FROM pg_stat_activity;"

# Check for connection pool exhaustion
docker stats listings

# Review application logs
tail -f /var/log/listings/app.log
```

### High Latency

```bash
# Check database query performance
# Enable slow query logging in PostgreSQL

# Check Redis cache hit rates
redis-cli info stats | grep keyspace

# Monitor system resources
htop
iotop
```

### Memory Leaks

```bash
# Monitor memory over time
watch -n 5 'ps aux | grep listings | grep -v grep'

# Generate heap profile (if Go service)
curl http://localhost:6060/debug/pprof/heap > heap.prof
go tool pprof heap.prof

# Check for connection leaks
lsof -p $(pgrep listings) | wc -l
```

## 🔄 CI/CD Integration

### GitLab CI Example

```yaml
load_test:
  stage: test
  image: grafana/k6:latest
  services:
    - postgres:14
    - redis:7
  script:
    - cd load-tests
    - ./run-all-tests.sh --no-monitor
  artifacts:
    paths:
      - load-tests/results/
    expire_in: 30 days
  only:
    - merge_requests
    - main
```

### GitHub Actions Example

```yaml
name: Load Tests

on:
  pull_request:
  push:
    branches: [main]

jobs:
  load-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Start services
        run: docker-compose up -d

      - name: Install k6
        run: |
          sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
          echo "deb https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
          sudo apt-get update
          sudo apt-get install k6

      - name: Run load tests
        run: cd load-tests && ./run-all-tests.sh

      - name: Upload results
        uses: actions/upload-artifact@v3
        with:
          name: load-test-results
          path: load-tests/results/
```

## 📊 Performance Benchmarks

### Target Performance (Single Instance)

| Metric | Value | Notes |
|--------|-------|-------|
| **Max RPS** | 200+ | With caching enabled |
| **Avg Latency** | 50-80ms | For read operations |
| **p95 Latency** | < 100ms | 95th percentile |
| **p99 Latency** | < 200ms | 99th percentile |
| **CPU Usage** | < 70% | At 100 RPS |
| **Memory Usage** | < 512MB | Base + working set |
| **DB Connections** | 10-20 | From pool of 50 |

### Scaling Recommendations

**Horizontal Scaling:**
- Add more instances behind load balancer
- Each instance handles 100-150 RPS comfortably
- 3 instances can handle 300-450 RPS with redundancy

**Vertical Scaling:**
- 2 CPU cores minimum
- 1GB RAM recommended
- SSD storage for database

**Database Optimization:**
- Connection pool: 50 connections per instance
- Read replicas for heavy read workloads
- Redis cache for frequently accessed data

## 🔍 Monitoring in Production

### Key Metrics to Monitor

```bash
# Prometheus metrics
rate(http_request_duration_seconds_bucket[5m])
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))
rate(http_requests_total{status=~"5.."}[5m])

# gRPC metrics
rate(grpc_server_handled_total[5m])
histogram_quantile(0.95, rate(grpc_server_handling_seconds_bucket[5m]))
```

### Alert Thresholds

```yaml
alerts:
  - name: HighLatency
    condition: p95_latency > 100ms for 5m
    severity: warning

  - name: HighErrorRate
    condition: error_rate > 1% for 5m
    severity: critical

  - name: LowThroughput
    condition: rps < 50 for 10m
    severity: warning
```

## 📚 Additional Resources

- [k6 Documentation](https://k6.io/docs/)
- [ghz Documentation](https://ghz.sh/)
- [gRPC Performance Best Practices](https://grpc.io/docs/guides/performance/)
- [HTTP/2 Performance Tuning](https://hpbn.co/http2/)

## 🤝 Contributing

When adding new tests:

1. Follow existing naming conventions
2. Add appropriate checks and thresholds
3. Update this README with new scenarios
4. Ensure tests pass locally before committing

## 📝 License

MIT License - See parent project for details

---

**Last Updated:** 2025-11-10
**Maintainer:** Development Team
**Version:** 1.0.0
