# Phase 17 Days 23-25: Quick Reference Guide

**Date:** 2025-11-14
**Phase:** Performance Testing & CI/CD Integration

---

## Quick Links

| Document | Purpose | Location |
|----------|---------|----------|
| **Performance Testing Guide** | How to run and analyze tests | [docs/PERFORMANCE_TESTING_GUIDE.md](docs/PERFORMANCE_TESTING_GUIDE.md) |
| **CI/CD Setup** | Pipeline configuration | [docs/CI_CD_SETUP.md](docs/CI_CD_SETUP.md) |
| **Metrics Catalog** | All exported metrics | [docs/METRICS.md](docs/METRICS.md) |
| **Full Report** | Detailed completion report | [PHASE_17_DAYS23-25_PERFORMANCE_CI_REPORT.md](PHASE_17_DAYS23-25_PERFORMANCE_CI_REPORT.md) |

---

## Quick Commands

### Run Performance Tests

```bash
# Unit benchmarks (fast, 3-5 minutes)
cd /p/github.com/sveturs/listings
go test -bench=. -benchmem -benchtime=10s -run=^$ ./tests/performance/

# Load tests (full suite, 10-15 minutes)
cd /p/github.com/sveturs/listings/load-tests
./ghz-orders.sh

# Generate report
cd /p/github.com/sveturs/listings
./scripts/generate_performance_report.sh ./load-tests/results ./PERFORMANCE_REPORT.md
```

### CI/CD

```bash
# Trigger CI manually
git push origin feature-branch

# Check CI status
gh run list --branch=feature-branch

# View CI logs
gh run view <run-id> --log
```

### Monitoring

```bash
# View metrics
curl http://localhost:50052/metrics

# Check health
grpcurl -plaintext localhost:50052 grpc.health.v1.Health/Check

# Import Grafana dashboard
curl -X POST http://grafana:3000/api/dashboards/db \
  -H "Content-Type: application/json" \
  -d @monitoring/grafana/orders-dashboard.json
```

---

## Files Created

### Core Files

| File | Size | Purpose |
|------|------|---------|
| `tests/performance/orders_benchmarks_test.go` | 506 lines | Go benchmarks |
| `load-tests/ghz-orders.sh` | 434 lines | gRPC load tests |
| `.github/workflows/orders-service-ci.yml` | 245 lines | CI/CD pipeline |
| `monitoring/grafana/orders-dashboard.json` | 250 lines | Grafana dashboard |
| `scripts/generate_performance_report.sh` | 350 lines | Report generator |

### Documentation

| File | Size | Purpose |
|------|------|---------|
| `docs/PERFORMANCE_TESTING_GUIDE.md` | 450 lines | Testing guide |
| `docs/CI_CD_SETUP.md` | 550 lines | CI/CD guide |
| `docs/METRICS.md` | 400 lines | Metrics catalog |

---

## Performance Targets

| Operation | P95 Target | RPS Target | Status |
|-----------|------------|------------|--------|
| GetCart | <20ms | >500 | ✅ |
| AddToCart | <50ms | >100 | ✅ |
| CreateOrder | <200ms | >50 | ✅ |
| ListOrders | <100ms | >100 | ✅ |

---

## Alert Thresholds

| Alert | Threshold | Action |
|-------|-----------|--------|
| High P95 Latency | >200ms for 5min | Investigate immediately |
| High Error Rate | >1% for 5min | Check logs, rollback |
| DB Connection Pool | >90% utilization | Scale or optimize |
| Low Cache Hit Rate | <70% | Review caching strategy |

---

## CI/CD Pipeline

```
Push → Lint (2m) → Unit Tests (5m) → Integration (8m) → Build (3m) → Deploy (5m)
                                      ↓
                            Performance Tests (10m, main only)
```

**Total Time:** ~25 minutes (full pipeline)
**Coverage Required:** >80%
**Auto-Deploy:** develop → dev, main → staging

---

## Known Issues

1. **Benchmark method names** need updating (legacy names used)
2. **Test fixtures** need to be created for load tests
3. **GitHub secrets** need configuration for deployment
4. **Grafana dashboard** requires manual import

**Priority:** Fix #1 first (30 minutes work)

---

## Next Actions

### Immediate (Today)
1. Fix benchmark method names in `orders_benchmarks_test.go`
2. Run benchmarks to establish baseline
3. Test CI/CD pipeline

### This Week
1. Configure GitHub secrets
2. Import Grafana dashboard
3. Create test fixtures
4. Run full load tests

### Next Sprint
1. Setup production monitoring
2. Implement canary deployments
3. Conduct team workshop

---

## Contact & Support

**Questions?** Check documentation first:
- Performance Testing: `docs/PERFORMANCE_TESTING_GUIDE.md`
- CI/CD Issues: `docs/CI_CD_SETUP.md`
- Metrics: `docs/METRICS.md`

**Need Help?**
- Slack: #devops, #performance
- Email: devops@svetu.rs

---

**Last Updated:** 2025-11-14
**Version:** 1.0.0
