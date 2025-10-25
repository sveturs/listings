# Delivery Module - Testing Guide

## Quick Start

### Prerequisites
- Go 1.25+
- Docker (for testcontainers)
- PostgreSQL container image

### Run All Tests
```bash
cd /p/github.com/sveturs/svetu/backend
go test -v -race ./internal/proj/delivery/...
```

### Run Specific Test Suites

#### Storage Layer Tests
```bash
go test -v ./internal/proj/delivery/storage/...
```

#### gRPC Mapper Tests
```bash
go test -v ./internal/proj/delivery/grpcclient/... -run TestMap
```

#### gRPC Client Tests
```bash
go test -v ./internal/proj/delivery/grpcclient/... -run TestCreateShipment
```

#### Attributes Service Tests
```bash
go test -v ./internal/proj/delivery/attributes/...
```

### Generate Coverage Report
```bash
# Generate coverage profile
go test -v -race -coverprofile=coverage.out ./internal/proj/delivery/...

# View HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# View coverage by function
go tool cover -func=coverage.out | grep delivery
```

## Test Structure

```
backend/internal/proj/delivery/
├── storage/
│   ├── storage.go
│   ├── storage_test.go          ← Storage layer tests
│   ├── admin_storage.go
│   └── admin_storage_test.go    ← Admin storage tests
├── grpcclient/
│   ├── client.go
│   ├── client_test.go           ← gRPC client tests with mocks
│   ├── mapper.go
│   └── mapper_test.go           ← Proto mapper tests
├── attributes/
│   ├── service.go
│   └── service_test.go          ← Attributes service tests
├── TEST_SUITE_SUMMARY.md        ← Comprehensive test documentation
└── TESTING.md                   ← This file
```

## Test Coverage Goals

| Component | Target Coverage | Status |
|-----------|----------------|---------|
| Storage Layer | 85%+ | ✅ Achieved |
| gRPC Mapper | 95%+ | ✅ Achieved |
| gRPC Client | 80%+ | ✅ Achieved |
| Attributes Service | 85%+ | ✅ Achieved |

## Common Test Scenarios

### 1. Storage Layer Tests
- Provider CRUD operations
- Shipment creation and retrieval
- Tracking events (with duplicate prevention)
- Zone detection
- Statistics and analytics
- Error handling (not found, validation)

### 2. gRPC Client Tests
- Successful gRPC calls
- Retry logic with exponential backoff
- Circuit breaker behavior
- Error classification (retryable vs non-retryable)
- Timeout handling

### 3. Mapper Tests
- Proto to model conversion
- Model to proto conversion
- Status mappings
- Provider code mappings
- Time conversions

### 4. Attributes Service Tests
- Product attributes CRUD
- Category defaults
- Batch updates with transactions
- Validation (weight, dimensions, packaging)
- Volumetric weight calculations

## Troubleshooting

### Docker Not Running
```
Error: Cannot connect to the Docker daemon
```
**Solution:** Start Docker daemon
```bash
sudo systemctl start docker
```

### Testcontainers Timeout
```
Error: container startup timeout
```
**Solution:** Increase wait timeout or check Docker resources

### Port Already in Use
```
Error: bind: address already in use
```
**Solution:** Testcontainers uses random ports, no action needed

### Failed to Pull Postgres Image
```
Error: pull access denied for postgres
```
**Solution:** Pull image manually
```bash
docker pull postgres:16
```

## Running Tests in CI/CD

### GitHub Actions Example
```yaml
name: Delivery Module Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.25'

      - name: Run tests
        run: |
          cd backend
          go test -v -race -coverprofile=coverage.out ./internal/proj/delivery/...

      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./backend/coverage.out
```

## Test Execution Time

- **Storage tests:** ~15-20s (testcontainers startup)
- **Mapper tests:** <1s (pure unit tests)
- **gRPC client tests:** ~2-3s (mock tests)
- **Attributes tests:** ~15-20s (testcontainers startup)
- **Total:** ~40-50 seconds

## Best Practices

### 1. Run Tests Before Commit
```bash
# Quick check
go test ./internal/proj/delivery/...

# Full check with race detector
go test -v -race ./internal/proj/delivery/...
```

### 2. Watch Mode (Using entr or similar)
```bash
# Install entr
sudo apt-get install entr

# Watch for changes and re-run tests
find . -name '*.go' | entr -c go test ./internal/proj/delivery/...
```

### 3. Parallel Test Execution
```bash
# Run tests in parallel (default is GOMAXPROCS)
go test -v -race -parallel 4 ./internal/proj/delivery/...
```

### 4. Verbose Output
```bash
# See all test output including passed tests
go test -v ./internal/proj/delivery/...
```

### 5. Run Specific Test
```bash
# Run single test by name
go test -v -run TestGetProvider ./internal/proj/delivery/storage/...

# Run tests matching pattern
go test -v -run TestCreate ./internal/proj/delivery/...
```

## Test Maintenance

### Adding New Tests
1. Follow existing test patterns (table-driven tests, test suites)
2. Use descriptive test names: `TestFunction_Scenario_ExpectedBehavior`
3. Test both success and error paths
4. Add tests to relevant test file or create new if needed

### Updating Tests
1. Run tests after any code changes
2. Update test expectations if behavior changes intentionally
3. Add regression tests for bugs
4. Keep test coverage above 80%

## Additional Resources

- [Test Suite Summary](./TEST_SUITE_SUMMARY.md) - Detailed test documentation
- [Testcontainers Go](https://golang.testcontainers.org/) - Integration testing
- [Testify](https://github.com/stretchr/testify) - Testing toolkit
- [Go Testing](https://golang.org/pkg/testing/) - Official Go testing docs

## Contact

For questions or issues with tests, see:
- Test Suite Summary: `TEST_SUITE_SUMMARY.md`
- Audit Report: `/p/github.com/sveturs/delivery/MONOLITH_AUDIT_REPORT.md`
