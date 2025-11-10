# Performance Testing - Quick Start Guide

## Installation (One-time)

```bash
# Install Go tools
go install github.com/bojand/ghz/cmd/ghz@latest
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Install utilities (Ubuntu/Debian)
sudo apt-get install -y jq bc

# Verify installation
ghz --version && grpcurl --version && jq --version
```

## Quick Check (40 seconds)

```bash
cd /p/github.com/sveturs/listings
./scripts/performance/quick-check.sh
```

**Output Example:**
```
Testing GetProduct... FAST (P95: 12.5ms, RPS: 102)
Testing CheckStock... FAST (P95: 18.3ms, RPS: 98)
Testing GetCategories... FAST (P95: 5.2ms, RPS: 205)
Testing ListProducts... OK (P95: 32.1ms, RPS: 95)
```

## Full Baseline (7 minutes)

```bash
cd /p/github.com/sveturs/listings
./scripts/performance/baseline.sh
```

**Results:**
- `baseline_results.json` - Machine-readable
- `baseline_results.txt` - Human-readable

## View Results

```bash
# Summary statistics
cat baseline_results.txt | grep -A 20 "SUMMARY"

# P95 latencies
jq -r '.results[] | "\(.method): \(.latency_ms.p95)ms"' baseline_results.json

# Find slowest endpoints
jq -r '.results | sort_by(.latency_ms.p99) | reverse | .[0:5] | .[] |
  "\(.method): P99=\(.latency_ms.p99)ms"' baseline_results.json
```

## Grafana Dashboard

```bash
# 1. Open Grafana
firefox http://localhost:3030

# 2. Login: admin / admin123

# 3. Import dashboard
#    - Go to Dashboards → Import
#    - Upload: monitoring/grafana/listings-performance-dashboard.json
#    - Select datasource: Prometheus
#    - Click Import

# 4. Access dashboard
firefox http://localhost:3030/d/listings-performance
```

## Prometheus Alerts

```bash
# Copy alerts to Prometheus
cp monitoring/prometheus/performance_alerts.yml \
   deployment/prometheus/rules/

# Reload Prometheus
docker exec listings_prometheus kill -HUP 1

# Verify alerts loaded
curl -s http://localhost:9090/api/v1/rules | \
  jq '.data.groups[] | select(.name == "performance_critical_alerts") | .rules[] | .name'
```

## Common Issues

### "ghz: command not found"
```bash
go install github.com/bojand/ghz/cmd/ghz@latest
export PATH=$PATH:$(go env GOPATH)/bin
```

### "Cannot connect to gRPC server"
```bash
# Check service running
netstat -tlnp | grep :8086

# Test connectivity
grpcurl -plaintext localhost:8086 list
```

### Performance degradation detected
```bash
# Check current metrics
curl -s http://localhost:8086/metrics | grep listings_grpc

# View Grafana dashboard
firefox http://localhost:3030/d/listings-performance

# Check Prometheus alerts
curl -s http://localhost:9090/api/v1/alerts
```

## Targets

| Metric | Warning | Critical |
|--------|---------|----------|
| P95 (core) | >50ms | >100ms |
| P99 (core) | >100ms | >200ms |
| Stock P95 | >30ms | >50ms |
| Error rate | >0.5% | >1% |

## Full Documentation

- Main: `/p/github.com/sveturs/listings/PHASE_13_4_1_BASELINE.md`
- Scripts: `/p/github.com/sveturs/listings/scripts/performance/README.md`
- Summary: `/p/github.com/sveturs/listings/PHASE_13_4_1_SUMMARY.md`
