# Production Readiness Test Suite

This directory contains comprehensive tests for Sprint 6.2 resilience patterns: timeout, circuit breaker, and fallback mechanisms.

## Quick Start

```bash
# Run all production readiness tests (recommended)
cd /p/github.com/sveturs/svetu/backend
make test-production-readiness

# Quick smoke test (2 minutes)
make test-smoke

# Individual test suites
make test-integration          # Timeout, circuit breaker tests
make test-integration-race     # Race condition detection
make test-load                 # Performance benchmarks
make test-e2e                  # End-to-end (requires running backend)
```

## Test Structure

```
tests/
├── integration/           # Integration tests
│   ├── timeout_test.go           # Timeout pattern (7 tests)
│   └── circuit_breaker_test.go   # Circuit breaker (8 tests)
├── e2e/                   # End-to-end tests
│   └── resilience_test.go        # Full system tests (6 tests)
├── load/                  # Load/performance tests
│   └── resilience_load_test.go   # Benchmarks (5 tests)
└── mocks/                 # Mock services
    └── microservice_mock.go      # gRPC mock server
```

## Test Categories

### 1. Timeout Integration Tests (7 tests)

**File:** `integration/timeout_test.go`

Tests the 500ms timeout pattern with fallback to monolith:

- ✅ Timeout triggers at configured duration (500ms)
- ✅ Fallback to monolith works on timeout
- ✅ Context cancellation propagates correctly
- ✅ Multiple concurrent timeouts handled
- ✅ No goroutine leaks on timeout
- ✅ Timeout works with retries

**Run:**
```bash
go test -v ./tests/integration/timeout_test.go
```

### 2. Circuit Breaker Tests (8 tests)

**File:** `integration/circuit_breaker_test.go`

Tests circuit breaker opening, closing, and state transitions:

- ✅ Circuit opens after 5 consecutive failures
- ✅ Circuit rejects requests in OPEN state
- ✅ Circuit transitions to HALF_OPEN after 30s timeout
- ✅ Circuit closes after 2 successful requests
- ✅ Concurrent requests handled in HALF_OPEN
- ✅ No race conditions under load

**Run:**
```bash
go test -v ./tests/integration/circuit_breaker_test.go
go test -v -race ./tests/integration/circuit_breaker_test.go  # With race detector
```

### 3. End-to-End Resilience Tests (6 tests)

**File:** `e2e/resilience_test.go`

Full system tests with real backend and mock microservice:

- ✅ Slow microservice → timeout → fallback
- ✅ Failing microservice → circuit opens
- ✅ Microservice recovery → circuit closes
- ✅ Mixed load with partial degradation
- ✅ Cascading failure prevention
- ✅ End-to-end latency <200ms

**Run:**
```bash
# Requires backend running on :3000
go test -v -tags=e2e ./tests/e2e/resilience_test.go
```

### 4. Load Tests (5 benchmarks)

**File:** `load/resilience_load_test.go`

Performance and stress tests:

- ✅ 1000 RPS with 10% timeouts
- ✅ 500 RPS with circuit breaker cycles
- ✅ 2000 RPS with mixed success/failure
- ✅ Memory stability under sustained load
- ✅ No goroutine leaks

**Run:**
```bash
go test -v -bench=. -benchtime=30s ./tests/load/resilience_load_test.go
```

## Mock Microservice

**File:** `mocks/microservice_mock.go`

Provides a controllable gRPC service for testing:

### Starting the Mock

```bash
# Option 1: Direct run
go run tests/mocks/microservice_mock.go

# Option 2: Via Makefile
make mock-start

# Stop
make mock-stop
```

### Ports

- **gRPC Service:** `:50051`
- **Control API:** `:50052` (HTTP)

### Control API

Configure mock behavior via HTTP:

```bash
# Check health
curl http://localhost:50052/health

# Get current config
curl http://localhost:50052/control/status

# Normal mode (fast responses)
curl -X POST http://localhost:50052/control/config \
  -H "Content-Type: application/json" \
  -d '{"mode":"normal"}'

# Slow mode (1s delay)
curl -X POST http://localhost:50052/control/config \
  -H "Content-Type: application/json" \
  -d '{"mode":"slow","delay":"1s"}'

# Error mode (always fail)
curl -X POST http://localhost:50052/control/config \
  -H "Content-Type: application/json" \
  -d '{"mode":"error"}'

# Partial mode (50% failure rate)
curl -X POST http://localhost:50052/control/config \
  -H "Content-Type: application/json" \
  -d '{"mode":"partial","failure_rate":50}'
```

### Special Test IDs

The mock recognizes special listing IDs for targeted testing:

- **ID 1:** Normal success response
- **ID 777:** Always returns internal error (triggers circuit breaker)
- **ID 888:** Returns unavailable error
- **ID 999:** Delays 1s then times out (triggers timeout)

## Success Criteria

All tests must meet these criteria to pass:

### Performance
- ✅ P99 latency <200ms
- ✅ Throughput >1000 RPS
- ✅ Error rate <1% under normal load

### Reliability
- ✅ Circuit breaker opens after 5 failures
- ✅ Circuit breaker recovers after 30s
- ✅ Timeout triggers at 500ms
- ✅ Fallback to monolith works

### Stability
- ✅ Memory increase <10% under sustained load
- ✅ No goroutine leaks
- ✅ No race conditions

## Test Results

See detailed test report:
- **Report:** `docs/migration/PRODUCTION_READINESS_TEST_REPORT.md`
- **Status:** ✅ ALL TESTS PASSING (26/26)

## Common Issues

### Mock service not running

**Error:** `connection refused`

**Solution:**
```bash
make mock-start
# Wait 2 seconds for startup
```

### Tests timeout

**Error:** `timeout after 60s`

**Solution:**
```bash
# Increase timeout
go test -v -timeout=120s ./tests/integration/...
```

### Race detector warnings

**Warning:** Not an error, but indicates potential issue

**Solution:**
```bash
# Fix race conditions before proceeding
go test -v -race ./tests/integration/...
```

## CI/CD Integration

### GitHub Actions

```yaml
- name: Run production readiness tests
  run: |
    cd backend
    make test-production-readiness
```

### Pre-deployment Check

```bash
# Must pass before deploying to staging
make test-production-readiness

# Generate coverage report
make test-coverage
```

## Debugging Tests

### Verbose output

```bash
go test -v ./tests/integration/timeout_test.go
```

### Run single test

```bash
go test -v -run TestTimeoutTriggersAtConfiguredDuration ./tests/integration/timeout_test.go
```

### Enable race detector

```bash
go test -v -race ./tests/integration/circuit_breaker_test.go
```

### Profile memory

```bash
go test -v -memprofile=mem.prof ./tests/load/resilience_load_test.go
go tool pprof mem.prof
```

### Profile CPU

```bash
go test -v -cpuprofile=cpu.prof -bench=. ./tests/load/resilience_load_test.go
go tool pprof cpu.prof
```

## Contributing

When adding new tests:

1. Follow existing test structure
2. Use descriptive test names (`TestFeature_Scenario_ExpectedResult`)
3. Add detailed comments explaining test purpose
4. Update this README with new tests
5. Ensure all tests pass before PR
6. Run race detector (`-race` flag)

## Related Documentation

- [Production Readiness Report](../../docs/migration/PRODUCTION_READINESS_TEST_REPORT.md)
- [Sprint 6.2 Plan](../../docs/migration/SPRINT_6_2_RESILIENCE_PATTERNS.md)
- [Circuit Breaker Design](../../backend/internal/clients/listings/client.go)
- [Metrics Documentation](../../backend/internal/metrics/migration_metrics.go)

## Contact

For questions or issues with tests:
- Create GitHub issue with `testing` label
- Ping @sveturs in Slack #backend channel
