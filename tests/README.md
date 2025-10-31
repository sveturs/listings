# Testing Guide - Listings Microservice

Comprehensive testing documentation for the Listings microservice.

## Table of Contents

- [Overview](#overview)
- [Test Structure](#test-structure)
- [Running Tests](#running-tests)
- [Test Categories](#test-categories)
- [Writing Tests](#writing-tests)
- [Coverage Goals](#coverage-goals)
- [Troubleshooting](#troubleshooting)

---

## Overview

The Listings microservice uses a comprehensive testing strategy with multiple layers:

- **Unit Tests**: Test individual components in isolation
- **Integration Tests**: Test component interactions with real dependencies
- **E2E Tests**: Test complete user workflows
- **Performance Tests**: Benchmark critical operations

### Testing Framework

- **Framework**: [testify](https://github.com/stretchr/testify) - assertions and test suites
- **Mocking**: [gomock](https://github.com/uber-go/mock) - mock generation and verification
- **Integration**: [dockertest](https://github.com/ory/dockertest) - Docker containers for tests
- **Coverage Target**: **70%+**

---

## Test Structure

```
tests/
├── testing.go              # Test helpers and fixtures
├── integration/            # Integration tests
│   ├── database_test.go   # PostgreSQL integration
│   ├── redis_test.go      # Redis integration
│   ├── opensearch_test.go # OpenSearch integration
│   └── auth_test.go       # Auth Service integration
├── e2e/                   # End-to-end tests
│   ├── listing_flow_test.go
│   └── search_flow_test.go
└── performance/           # Performance benchmarks
    └── benchmarks_test.go

internal/
├── repository/postgres/
│   └── repository_test.go # Repository unit tests
├── service/listings/
│   └── service_test.go    # Service unit tests
└── transport/
    ├── http/server_test.go   # HTTP handler tests
    └── grpc/server_test.go   # gRPC handler tests
```

---

## Running Tests

### Quick Start

```bash
# Run all unit tests (fast, no Docker required)
make test

# Run all tests (unit + integration)
make test-all

# Run tests with coverage report
make test-coverage

# Check coverage threshold (70%)
make test-coverage-check
```

### Specific Test Categories

```bash
# Unit tests only
make test-unit

# Integration tests (requires Docker)
make test-integration

# E2E tests
make test-e2e

# Performance benchmarks
make bench
```

### Advanced Testing

```bash
# Verbose output
make test-verbose

# Watch mode (auto-run on file change)
make test-watch

# CPU profiling
make bench-cpu

# Memory profiling
make bench-mem
```

### Manual Test Execution

```bash
# Run specific test
go test -v ./internal/repository/postgres -run TestCreateListing

# Run with race detector
go test -race ./...

# Run in short mode (skips integration tests)
go test -short ./...

# Integration tests only
go test -tags=integration ./tests/integration/...
```

---

## Test Categories

### 1. Unit Tests

**Location**: `internal/*/\*_test.go`

**Purpose**: Test individual functions and methods in isolation

**Characteristics**:
- Fast execution (<100ms per test)
- No external dependencies
- Uses mocks for dependencies
- High code coverage target

**Example**:

```go
func TestCreateListing(t *testing.T) {
    repo, testDB := setupTestRepo(t)
    defer testDB.TeardownTestPostgres(t)

    input := &domain.CreateListingInput{
        UserID: 1,
        Title:  "Test Product",
        // ...
    }

    listing, err := repo.CreateListing(ctx, input)

    require.NoError(t, err)
    assert.Equal(t, input.Title, listing.Title)
}
```

### 2. Integration Tests

**Location**: `tests/integration/`

**Build Tag**: `// +build integration`

**Purpose**: Test component interactions with real services

**Characteristics**:
- Uses Docker containers (PostgreSQL, Redis, etc.)
- Tests real database operations
- Slower than unit tests
- Validates integration points

**Example**:

```go
// +build integration

func TestDatabaseIntegration(t *testing.T) {
    tests.SkipIfNoDocker(t)

    testDB := tests.SetupTestPostgres(t)
    defer testDB.TeardownTestPostgres(t)

    // Test full workflow
    // ...
}
```

### 3. E2E Tests

**Location**: `tests/e2e/`

**Purpose**: Test complete user workflows

**Characteristics**:
- Full system running
- Tests HTTP/gRPC endpoints
- Validates business logic
- Most comprehensive

### 4. Performance Tests

**Location**: `tests/performance/`

**Purpose**: Benchmark critical operations

**Characteristics**:
- Measures execution time
- Memory allocation tracking
- Concurrency testing
- Performance regression detection

**Example**:

```go
func BenchmarkGetListingByID(b *testing.B) {
    // Setup

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        repo.GetListingByID(ctx, id)
    }
}
```

---

## Writing Tests

### Test Naming Convention

```go
// Unit test
func TestFunctionName(t *testing.T)

// Table-driven test
func TestFunctionName(t *testing.T) {
    testCases := []struct {
        name    string
        input   Input
        want    Output
        wantErr bool
    }{
        {"valid input", validInput, expectedOutput, false},
        {"invalid input", invalidInput, nil, true},
    }

    for _, tt := range testCases {
        t.Run(tt.name, func(t *testing.T) {
            // Test logic
        })
    }
}

// Benchmark
func BenchmarkFunctionName(b *testing.B)
```

### Using Test Helpers

```go
import "github.com/sveturs/listings/tests"

// Setup test database
testDB := tests.SetupTestPostgres(t)
defer testDB.TeardownTestPostgres(t)

// Run migrations
tests.RunMigrations(t, testDB.DB, "../../migrations")

// Load fixtures
tests.LoadTestFixtures(t, testDB.DB, "fixtures.sql")

// Generate test data
listing := tests.GenerateTestListing(1, "Test Product")

// Test context with timeout
ctx := tests.TestContext(t)

// Skip conditions
tests.SkipIfShort(t)
tests.SkipIfNoDocker(t)
```

### Assertions

```go
// testify/require - fails test immediately
require.NoError(t, err)
require.NotNil(t, listing)
require.Equal(t, expected, actual)

// testify/assert - continues test execution
assert.Equal(t, expected, actual)
assert.Greater(t, listing.ID, int64(0))
assert.Contains(t, list, item)
```

### Mocking (gomock)

```go
// Generate mocks (to be implemented in Sprint 4.2)
//go:generate mockgen -destination=mocks/mock_repository.go -package=mocks . Repository

// Use in tests
mockRepo := mocks.NewMockRepository(ctrl)
mockRepo.EXPECT().GetListingByID(ctx, id).Return(listing, nil)
```

---

## Coverage Goals

### Coverage Targets

- **Overall**: ≥70%
- **Critical Paths**: ≥90% (repository, service layers)
- **Handlers**: ≥80%
- **Utils**: ≥60%

### Checking Coverage

```bash
# Generate coverage report
make test-coverage

# Check coverage threshold
make test-coverage-check

# View HTML report
open coverage.html

# View coverage by function
go tool cover -func=coverage.txt
```

### Coverage Exclusions

Some code is excluded from coverage requirements:
- Main entry points
- Generated code (protobuf, mocks)
- Third-party integrations (tested via integration tests)

---

## Performance Benchmarks

### Running Benchmarks

```bash
# Run all benchmarks
make bench

# Run specific benchmark
go test -bench=BenchmarkCreateListing ./tests/performance/

# With memory stats
go test -bench=. -benchmem ./tests/performance/

# Compare results
go test -bench=. -benchmem ./tests/performance/ > old.txt
# Make changes
go test -bench=. -benchmem ./tests/performance/ > new.txt
benchcmp old.txt new.txt
```

### Performance Targets

| Operation | Target (p95) | Current |
|-----------|--------------|---------|
| CreateListing | <50ms | TBD |
| GetListing | <20ms | TBD |
| Search | <100ms | TBD |
| UpdateListing | <50ms | TBD |

---

## Troubleshooting

### Common Issues

#### 1. "Docker not available"

**Problem**: Integration tests fail with Docker error

**Solution**:
```bash
# Check Docker is running
docker ps

# Install Docker if needed
# https://docs.docker.com/get-docker/

# Skip integration tests
make test-unit
```

#### 2. "Port already in use"

**Problem**: Test containers fail to start

**Solution**:
```bash
# Find and kill process using port
lsof -ti:35433 | xargs kill -9

# Clean up Docker containers
docker ps -a | grep listings | awk '{print $1}' | xargs docker rm -f
```

#### 3. Tests hang or timeout

**Problem**: Tests don't complete

**Solution**:
- Check for missing `defer cleanup()` calls
- Verify context timeouts are set
- Look for goroutine leaks

#### 4. Flaky tests

**Problem**: Tests fail intermittently

**Solution**:
- Avoid time-based assertions
- Use proper test cleanup
- Ensure test independence
- Check for race conditions: `go test -race`

### Debug Tips

```bash
# Verbose test output
go test -v ./...

# Show test coverage
go test -cover ./...

# Race detector
go test -race ./...

# CPU profiling
go test -cpuprofile=cpu.prof
go tool pprof cpu.prof

# Memory profiling
go test -memprofile=mem.prof
go tool pprof mem.prof
```

---

## Best Practices

1. **Test Independence**: Each test should clean up after itself
2. **Table-Driven Tests**: Use for multiple scenarios
3. **Clear Test Names**: Describe what is being tested
4. **AAA Pattern**: Arrange, Act, Assert
5. **Mock External Dependencies**: Keep unit tests fast
6. **Integration Tests**: Test real scenarios
7. **Performance Benchmarks**: Track regressions
8. **Coverage**: Aim for meaningful coverage, not just numbers

---

## CI/CD Integration

Tests run automatically on:
- Every push
- Pull requests
- Before deployment

**CI Pipeline**:
1. Lint code
2. Run unit tests
3. Check coverage threshold
4. Run integration tests
5. Run E2E tests
6. Upload coverage report

---

## Additional Resources

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Testify Documentation](https://github.com/stretchr/testify)
- [gomock Documentation](https://github.com/uber-go/mock)
- [dockertest Documentation](https://github.com/ory/dockertest)

---

**Last Updated**: 2025-10-31
**Sprint**: 4.2 - Testing Implementation
