# Quick Start Guide - Load Testing

Get up and running with load tests in 5 minutes.

## 🚀 Step 1: Install Dependencies

```bash
# Install k6 (HTTP load testing)
sudo snap install k6

# Install ghz (gRPC load testing)
go install github.com/bojand/ghz/cmd/ghz@latest

# Install utilities
sudo apt-get install jq bc netcat
```

## 🏗️ Step 2: Start the Service

### Option A: Local Development

```bash
cd /p/github.com/sveturs/listings

# Start dependencies
docker-compose up -d postgres redis

# Run migrations
make migrate-up

# Start the service
make run
```

### Option B: Full Docker Stack

```bash
cd /p/github.com/sveturs/listings/load-tests

# Start all services (app + monitoring)
docker-compose -f docker-compose.load-test.yml up -d

# Wait for services to be ready (~30 seconds)
docker-compose -f docker-compose.load-test.yml ps

# Check health
curl http://localhost:8086/health
```

## ✅ Step 3: Verify Service is Ready

```bash
# Check HTTP endpoint
curl http://localhost:8086/health

# Check gRPC endpoint
nc -zv localhost 50051

# Optional: Check gRPC methods
grpcurl -plaintext localhost:50051 list listings.v1.ListingsService
```

## 🧪 Step 4: Run Load Tests

### Quick Test (Recommended for first run)

```bash
cd /p/github.com/sveturs/listings/load-tests

# Run all tests
./run-all-tests.sh
```

This will:
- ✅ Check all dependencies
- ✅ Verify service availability
- ✅ Run HTTP load tests (~5 minutes)
- ✅ Run gRPC load tests (~6 minutes)
- ✅ Monitor system metrics
- ✅ Generate summary report

### Run Individual Tests

```bash
# Only HTTP tests
./run-all-tests.sh --http-only

# Only gRPC tests
./run-all-tests.sh --grpc-only

# Without monitoring
./run-all-tests.sh --no-monitor
```

### Manual Test Execution

```bash
# HTTP test (k6)
k6 run --env BASE_URL=http://localhost:8086 k6-http.js

# gRPC test (ghz)
./ghz-grpc.sh
```

## 📊 Step 5: View Results

### Quick Summary

```bash
# View latest test summary
cat results/summary_*.txt | tail -100

# List all results
ls -lh results/
```

### Detailed Analysis

```bash
# Analyze latest results
./analyze-results.sh

# Analyze specific test run
./analyze-results.sh 20251110_143000

# View HTTP results (JSON)
cat results/k6_results_*.json | jq '.metrics'

# View gRPC results
cat results/get_all_categories_*.json | jq '.latencyDistribution'

# View system metrics
cat results/system_metrics_*.log | column -t -s,
```

### Grafana Dashboard (Docker setup only)

```bash
# Open Grafana
open http://localhost:3000

# Login: admin / admin
# Navigate to Dashboards → Browse
```

## ✅ Success Criteria Check

Your tests pass if you see:

```
✅ Success Criteria Evaluation:
--------------------------------
✓ p95 latency < 100ms: 87.32ms
✓ No errors (0% error rate)
✓ Throughput >= 100 RPS: 125.45 RPS

🎉 All success criteria passed!
```

## 🐛 Troubleshooting

### Service not available

```bash
# Check if service is running
ps aux | grep listings

# Check ports
netstat -tlnp | grep -E '8086|50051'

# View service logs
docker-compose -f docker-compose.load-test.yml logs -f listings
```

### High error rates during test

```bash
# Check database connections
docker exec listings_load_test_db psql -U listings_user -d listings_db -c "SELECT count(*) FROM pg_stat_activity;"

# Check Redis
docker exec listings_load_test_redis redis-cli ping

# Reduce load
# Edit k6-http.js: lower target RPS values
# Edit ghz-grpc.sh: reduce --rps parameters
```

### Dependencies not found

```bash
# Verify PATH includes Go bin
echo $PATH | grep -q "$HOME/go/bin" || echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc

# Reload shell
source ~/.bashrc

# Check installations
which k6 ghz jq bc
```

## 📈 Understanding Results

### Key Metrics

| Metric | Good | Warning | Critical |
|--------|------|---------|----------|
| **p95 Latency** | < 100ms | 100-200ms | > 200ms |
| **Error Rate** | < 0.1% | 0.1-1% | > 1% |
| **RPS** | > 100 | 50-100 | < 50 |
| **CPU** | < 70% | 70-85% | > 85% |
| **Memory** | < 70% | 70-85% | > 85% |

### What to Look For

✅ **Good Signs:**
- Latency stays flat during sustained load
- Error rate near 0%
- CPU/Memory stable
- No connection errors

❌ **Warning Signs:**
- Latency increases over time (leak?)
- Error rate > 1%
- CPU constantly high
- Database connection errors

## 🎯 Next Steps

1. **Baseline Established**: Save these results as baseline for future comparisons
2. **Continuous Testing**: Run before each deployment
3. **Performance Regression**: Alert if metrics degrade > 10%
4. **Capacity Planning**: Use results to plan scaling strategy

## 📚 Full Documentation

See [README.md](README.md) for complete documentation including:
- Detailed test scenarios
- CI/CD integration
- Advanced analysis
- Performance tuning tips

---

**Quick Links:**
- [Main README](README.md)
- [Test Results](results/)
- [Grafana](http://localhost:3000) (Docker setup)
- [Prometheus](http://localhost:9090) (Docker setup)
