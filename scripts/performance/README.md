# Performance Testing Scripts

This directory contains performance baseline measurement and testing scripts for the Listings Microservice.

## Scripts Overview

### 1. baseline.sh
**Purpose:** Comprehensive baseline measurement of 13 critical gRPC endpoints

**Features:**
- Measures P50/P95/P99 latencies
- Calculates throughput (RPS)
- Tracks error rates
- Generates JSON + text reports
- Configurable test parameters

**Usage:**
```bash
# Default configuration (30s, 10 concurrency, 100 RPS)
./baseline.sh

# Custom configuration
./baseline.sh --duration 60 --concurrency 20 --rate 200

# Save to specific location
./baseline.sh --output /tmp/baseline_$(date +%Y%m%d).json

# Custom gRPC server
./baseline.sh --grpc-addr dev.vondi.rs:8086
```

**Options:**
- `-d, --duration SECONDS` - Duration for each test (default: 30)
- `-c, --concurrency NUM` - Number of concurrent connections (default: 10)
- `-r, --rate NUM` - Target requests per second (default: 100)
- `-o, --output FILE` - Output file path (default: baseline_results.json)
- `--grpc-addr ADDR` - gRPC server address (default: localhost:8086)
- `--metrics-url URL` - Prometheus metrics URL (default: http://localhost:8086/metrics)
- `-h, --help` - Show help message

**Prerequisites:**
```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
go install github.com/bojand/ghz/cmd/ghz@latest
```

**Output:**
- `baseline_results.json` - Machine-readable JSON with all metrics
- `baseline_results.txt` - Human-readable text report with summary

---

### 2. quick-check.sh
**Purpose:** Rapid 10-second performance validation of 4 most critical endpoints

**Features:**
- Fast execution (40 seconds total)
- Color-coded results (Fast/OK/Slow)
- Minimal resource usage
- Quick health check

**Usage:**
```bash
# Local server
./quick-check.sh

# Remote server
./quick-check.sh dev.vondi.rs:8086
```

**Tested Endpoints:**
1. GetProduct - Single product lookup
2. CheckStockAvailability - Stock check (critical for orders)
3. GetAllCategories - Category tree
4. ListProducts - Paginated product list

**Output Example:**
```
==========================================
  Quick Performance Check
==========================================
gRPC Server: localhost:8086
Duration: 10 seconds per endpoint

Testing GetProduct... FAST (P95: 12.5ms, RPS: 102)
Testing CheckStock... FAST (P95: 18.3ms, RPS: 98)
Testing GetCategories... FAST (P95: 5.2ms, RPS: 205)
Testing ListProducts... OK (P95: 32.1ms, RPS: 95)

Quick check complete!
```

**Color Coding:**
- 🟢 FAST: P95 < 50ms
- 🟡 OK: P95 50-100ms
- 🔴 SLOW: P95 > 100ms

---

## Critical Endpoints Tested

### CRITICAL Priority (Order Processing)
- **CheckStockAvailability** (single & batch) - P95 target: <30ms
- **DecrementStock** - P95 target: <30ms

### HIGH Priority (Core CRUD)
- **GetProduct** - P95 target: <50ms
- **ListProducts** - P95 target: <50ms
- **GetListing** - P95 target: <50ms
- **SearchListings** - P95 target: <70ms
- **GetAllCategories** - P95 target: <20ms
- **GetRootCategories** - P95 target: <20ms
- **GetStorefront** - P95 target: <50ms
- **ListStorefronts** - P95 target: <50ms

### MEDIUM Priority (Batch/Stats)
- **GetProductsByIDs** - P95 target: <100ms
- **GetProductStats** - P95 target: <100ms

---

## Installation

### Install Prerequisites

```bash
# Install ghz (gRPC benchmarking tool)
go install github.com/bojand/ghz/cmd/ghz@latest

# Install grpcurl (gRPC debugging tool)
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Install jq (JSON processor)
sudo apt-get install jq  # Ubuntu/Debian
brew install jq          # macOS

# Install bc (calculator)
sudo apt-get install bc  # Ubuntu/Debian
brew install bc          # macOS

# Verify installation
ghz --version
grpcurl --version
jq --version
bc --version
```

### Verify Service Running

```bash
# Check if listings service is accessible
netstat -tlnp | grep :8086

# Test gRPC connectivity
grpcurl -plaintext localhost:8086 list

# Verify metrics endpoint
curl -s http://localhost:8086/metrics | head -20
```

---

## Examples

### Example 1: Daily Baseline Test
```bash
#!/bin/bash
# Run daily baseline and archive results

DATE=$(date +%Y%m%d)
RESULTS_DIR=/var/log/performance/listings

mkdir -p "$RESULTS_DIR"

./baseline.sh \
  --duration 60 \
  --concurrency 10 \
  --rate 100 \
  --output "$RESULTS_DIR/baseline_$DATE.json"

# Archive older results (keep 30 days)
find "$RESULTS_DIR" -name "baseline_*.json" -mtime +30 -delete
```

### Example 2: Pre-Deployment Validation
```bash
#!/bin/bash
# Run quick check before deploying

./quick-check.sh localhost:8086

if [ $? -ne 0 ]; then
    echo "Performance check failed! Aborting deployment."
    exit 1
fi

echo "Performance check passed. Proceeding with deployment."
```

### Example 3: Continuous Monitoring
```bash
#!/bin/bash
# Run baseline every hour and alert on degradation

while true; do
    ./baseline.sh --output /tmp/baseline_current.json

    # Extract average P95
    AVG_P95=$(jq -r '[.results[].latency_ms.p95] | add / length' /tmp/baseline_current.json)

    if (( $(echo "$AVG_P95 > 50" | bc -l) )); then
        echo "ALERT: Average P95 latency is ${AVG_P95}ms (threshold: 50ms)"
        # Send alert notification here
    fi

    sleep 3600  # Wait 1 hour
done
```

### Example 4: Compare Before/After
```bash
#!/bin/bash
# Compare performance before and after changes

# Before changes
./baseline.sh --output before.json

# Apply changes...
# ...

# After changes
./baseline.sh --output after.json

# Compare results
echo "Performance Comparison:"
jq -r '.results[] | "\(.method): \(.latency_ms.p95)ms"' before.json > before.txt
jq -r '.results[] | "\(.method): \(.latency_ms.p95)ms"' after.json > after.txt

echo "=== BEFORE ==="
cat before.txt
echo ""
echo "=== AFTER ==="
cat after.txt
echo ""

# Show difference
paste before.txt after.txt | awk '{print $1, $2, "->", $4, "(" ($4-$2) "ms)"}'
```

---

## Interpreting Results

### Latency Metrics

**P50 (Median):**
- 50% of requests complete in this time or less
- Represents typical user experience
- Target: <20ms for critical operations

**P95:**
- 95% of requests complete in this time or less
- Primary SLO metric
- Target: <50ms for core operations

**P99:**
- 99% of requests complete in this time or less
- Worst-case monitoring
- Target: <100ms

### Error Rate

- **0% errors:** Perfect (rare in production)
- **<0.1% errors:** Excellent
- **0.1-0.5% errors:** Good
- **0.5-1% errors:** Warning - investigate
- **>1% errors:** Critical - immediate action required

### Throughput (RPS)

- **Achieved RPS ≈ Target RPS:** Service handling load well
- **Achieved RPS < Target RPS:** Service saturated or limited
- **High RPS + High latency:** Resource contention
- **Low RPS + Low latency:** Underutilized

---

## Troubleshooting

### Error: "ghz: command not found"
```bash
# Install ghz
go install github.com/bojand/ghz/cmd/ghz@latest

# Add Go bin to PATH
export PATH=$PATH:$(go env GOPATH)/bin

# Or permanently add to ~/.bashrc
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc
```

### Error: "Cannot connect to gRPC server"
```bash
# Check service is running
netstat -tlnp | grep :8086

# Test basic connectivity
grpcurl -plaintext localhost:8086 list

# Check service logs
docker logs listings_app

# Verify firewall
sudo iptables -L | grep 8086
```

### Error: "Benchmark times out"
```bash
# Reduce load
./baseline.sh --concurrency 5 --rate 50

# Increase timeout in script (edit baseline.sh)
# Change: --duration 30s
# To:     --duration 60s
```

### Error: "No metrics available"
```bash
# Check metrics endpoint
curl -s http://localhost:8086/metrics | head -20

# Verify Prometheus instrumentation
curl -s http://localhost:8086/metrics | grep listings_grpc

# Check if metrics are being recorded
watch -n 1 'curl -s http://localhost:8086/metrics | grep -c listings_grpc_requests_total'
```

---

## Integration with CI/CD

### GitHub Actions Example
```yaml
name: Performance Baseline

on:
  pull_request:
    branches: [main]

jobs:
  performance:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2

      - name: Install prerequisites
        run: |
          go install github.com/bojand/ghz/cmd/ghz@latest
          sudo apt-get install -y jq bc

      - name: Start services
        run: docker-compose up -d

      - name: Wait for service
        run: |
          timeout 60 bash -c 'until grpcurl -plaintext localhost:8086 list; do sleep 2; done'

      - name: Run baseline test
        run: |
          cd scripts/performance
          ./baseline.sh --duration 30 --output baseline.json

      - name: Check performance regression
        run: |
          AVG_P95=$(jq -r '[.results[].latency_ms.p95] | add / length' baseline.json)
          if (( $(echo "$AVG_P95 > 50" | bc -l) )); then
            echo "Performance regression detected: P95=${AVG_P95}ms"
            exit 1
          fi

      - name: Upload results
        uses: actions/upload-artifact@v2
        with:
          name: performance-baseline
          path: baseline.json
```

---

## References

- **Full Documentation:** `/p/github.com/sveturs/listings/PHASE_13_4_1_BASELINE.md`
- **Grafana Dashboard:** `/p/github.com/sveturs/listings/monitoring/grafana/listings-performance-dashboard.json`
- **Prometheus Alerts:** `/p/github.com/sveturs/listings/monitoring/prometheus/performance_alerts.yml`

---

## Support

For issues or questions:
1. Check the troubleshooting section above
2. Review main documentation in `PHASE_13_4_1_BASELINE.md`
3. Check Grafana dashboard for real-time metrics
4. Review Prometheus alerts for active issues
