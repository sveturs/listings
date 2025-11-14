# E2E Testing Guide for Orders Microservice

## Overview

This guide explains how to run end-to-end (E2E) tests for the Orders Microservice using Docker-based testcontainers for complete isolation and reproducibility.

## Test Infrastructure

### Testcontainers Setup

The test infrastructure uses **testcontainers-go** to provide:
- PostgreSQL 15 (ephemeral container per test environment)
- Redis 7 (ephemeral container per test environment)
- Automatic schema migrations
- Test data seeding
- Complete cleanup after tests

**Location**: `/internal/testing/environment.go`

### TestEnvironment

The `TestEnvironment` struct provides a complete testing environment:

```go
type TestEnvironment struct {
    // Infrastructure
    DB       *sqlx.DB
    Redis    *redis.Client
    Logger   zerolog.Logger

    // Repository
    Repo *postgres.Repository

    // Services
    CartService      service.CartService
    OrderService     service.OrderService
    InventoryService service.InventoryService
}
```

## Running Tests

### Prerequisites

- Docker installed and running
- Go 1.25+
- Network access for pulling Docker images (postgres:15-alpine, redis:7-alpine)

### Quick Start

```bash
# Run all E2E tests
cd /p/github.com/sveturs/listings
go test -v ./internal/transport/grpc/... -run Test

# Run specific test
go test -v ./internal/transport/grpc/... -run TestAddToCart_Success

# Skip long-running tests (for CI)
go test -short ./internal/transport/grpc/...

# Run with timeout
go test -v ./internal/transport/grpc/... -timeout 10m
```

### Test Execution Time

- **Per test environment creation**: ~5-10 seconds
- **Single test execution**: ~1-3 seconds
- **All 23 E2E tests**: ~60-120 seconds (depending on hardware)

## Test Categories

### 1. Cart Operations (6 tests)

Tests for shopping cart functionality:

1. **TestAddToCart_Success** - Add items to cart successfully
2. **TestAddToCart_InvalidInput** - Validation error handling (5 subtests)
3. **TestGetCart_Success** - Retrieve cart with items
4. **TestUpdateCartItem_Success** - Update item quantity
5. **TestRemoveFromCart_Success** - Remove item from cart
6. **TestClearCart_Success** - Clear all cart items
7. **TestGetUserCarts_Success** - Get all user carts across storefronts

### 2. Order Operations (6 tests)

Tests for order management:

1. **TestCreateOrder_Success** - Create order from cart
2. **TestCreateOrder_EmptyCart** - Fail when cart is empty
3. **TestGetOrder_Success** - Retrieve order by ID
4. **TestGetOrder_Unauthorized** - Access control for orders
5. **TestListOrders_Success** - List user orders with pagination
6. **TestCancelOrder_Success** - Cancel order and initiate refund

### 3. Admin Operations (2 tests)

Tests for administrative functions:

1. **TestUpdateOrderStatus_Success** - Change order status (admin)
2. **TestGetOrderStats_Success** - Retrieve order statistics

### 4. Error Handling (9 tests embedded)

Tests embedded in validation scenarios:
- No user_id or session_id
- Both user_id and session_id provided
- Invalid storefront_id
- Invalid listing_id
- Invalid quantity
- Unauthorized access
- Empty cart operations
- Missing required fields

## Test Structure

### Standard Test Pattern

```go
func TestExampleOperation(t *testing.T) {
    // Skip in short mode
    if testing.Short() {
        t.Skip("Skipping E2E test in short mode")
    }

    // Setup test environment (accepts testing.TB interface - works with *testing.T and *testing.B)
    env := testingpkg.NewTestEnvironment(t)
    defer env.Cleanup()

    // Seed test data
    env.SeedTestData(t)

    // Create gRPC handler
    orderHandler := grpcTransport.NewOrderServiceServer(
        env.CartService,
        env.OrderService,
        env.InventoryService,
        env.Logger,
    )

    // Execute test
    ctx := context.Background()
    req := &ordersspb.SomeRequest{...}
    resp, err := orderHandler.SomeMethod(ctx, req)

    // Assertions
    require.NoError(t, err)
    assert.NotNil(t, resp)
}
```

### Test Data Seeding

The `SeedTestData` function creates:

- **Categories**: Test Electronics (ID: 1000), Test Phones (ID: 1001)
- **Storefronts**: Test Store 1 (ID: 1), Test Store 2 (ID: 2)
- **Products**: IDs 100, 101, 200 with prices €99.99, €149.99, €199.99
- **Inventory**: 100, 50, 200 units respectively

## Writing New Tests

### Step 1: Define Test Function

```go
func TestYourNewTest(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping E2E test in short mode")
    }

    env := testingpkg.NewTestEnvironment(t)
    defer env.Cleanup()
    env.SeedTestData(t)

    // Your test code here
}
```

### Step 2: Use Test Environment

The TestEnvironment provides access to all infrastructure via `testing.TB` interface:

```go
// Access database
var count int
err := env.DB.Get(&count, "SELECT COUNT(*) FROM orders")

// Access Redis
ctx := context.Background()
err := env.Redis.Set(ctx, "key", "value", 0).Err()

// Use services (initialized automatically)
cart, err := env.CartService.GetCart(ctx, &userID, nil, storefrontID)
order, err := env.OrderService.CreateOrder(ctx, &userID, storefrontID, items)
```

### Step 3: Create Test Data

```go
// Create custom product
productID := env.CreateTestProduct(t, storefrontID, 299.99)

// Insert data directly
_, err := env.DB.Exec(`
    INSERT INTO custom_table (column1, column2)
    VALUES ($1, $2)
`, value1, value2)
```

### Step 4: Cleanup (Automatic)

The `defer env.Cleanup()` handles:
- Rolling back transactions
- Stopping Docker containers
- Closing connections

## Troubleshooting

### Issue: "Docker not available"

**Solution**: Ensure Docker is running:
```bash
docker ps
# Should show running containers
```

### Issue: "Port already in use"

**Solution**: Testcontainers uses random ports. If this happens, old containers may be lingering:
```bash
docker ps -a | grep test
docker rm -f $(docker ps -aq --filter "name=test")
```

### Issue: "Tests timeout"

**Solution**: Increase timeout:
```bash
go test -v ./internal/transport/grpc/... -timeout 15m
```

### Issue: "Too many open connections"

**Solution**: Reduce parallel tests or increase PostgreSQL max_connections:
```bash
go test -p 1 ./internal/transport/grpc/...
```

### Issue: "Migration failures"

**Solution**: Check migration files:
```bash
ls -la /p/github.com/sveturs/listings/migrations/*.up.sql
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: E2E Tests
on: [push, pull_request]

jobs:
  e2e-tests:
    runs-on: ubuntu-latest
    services:
      docker:
        image: docker:dind
        options: --privileged

    steps:
      - uses: actions/checkout@v3

      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.25'

      - name: Run E2E Tests
        run: |
          go test -v ./internal/transport/grpc/... -timeout 15m
```

### GitLab CI Example

```yaml
e2e-tests:
  image: golang:1.25
  services:
    - docker:dind
  variables:
    DOCKER_HOST: tcp://docker:2375
  script:
    - go test -v ./internal/transport/grpc/... -timeout 15m
```

## Performance Tips

### 1. Parallel Test Execution

```bash
# Run tests in parallel (default is GOMAXPROCS)
go test -v ./internal/transport/grpc/... -parallel 4
```

### 2. Skip Short Tests in CI

```bash
# In CI, skip short tests
go test -v ./internal/transport/grpc/... -short=false
```

### 3. Reuse Test Database

For faster local development, use existing database:

```go
config := testingpkg.DefaultTestEnvironmentConfig()
config.UseDocker = false
config.PostgresDSN = "postgres://test:test@localhost:5432/testdb?sslmode=disable"
config.RedisAddr = "localhost:6379"

env := testingpkg.NewTestEnvironmentWithConfig(t, config)
```

### 4. Cache Docker Images

Pre-pull images before running tests:

```bash
docker pull postgres:15-alpine
docker pull redis:7-alpine
```

## Performance Benchmarking

The test infrastructure supports **both unit tests and performance benchmarks** using the same `TestEnvironment`.

### Running Benchmarks

```bash
# Run all benchmarks with short duration
go test -bench=. -benchtime=1s -run=^$ ./tests/performance/

# Run specific benchmark
go test -bench=BenchmarkAddToCart -benchtime=5s -run=^$ ./tests/performance/

# Run benchmarks with memory profiling
go test -bench=. -benchmem -benchtime=5s -run=^$ ./tests/performance/

# Generate CPU profile
go test -bench=BenchmarkConcurrentCreateOrder -cpuprofile=cpu.prof -run=^$ ./tests/performance/
go tool pprof cpu.prof
```

### Benchmark Test Pattern

Benchmarks use the **same TestEnvironment** as E2E tests via the `testing.TB` interface:

```go
func BenchmarkYourOperation(b *testing.B) {
    if testing.Short() {
        b.Skip("Skipping benchmark in short mode")
    }

    // Setup test environment ONCE (outside b.N loop)
    env := testingpkg.NewTestEnvironment(b)
    defer env.Cleanup()

    env.SeedTestData(b)

    // Prepare test data ONCE
    ctx := context.Background()
    listing := createTestListing(b, ctx, env.Repo)
    cart := createTestCart(b, ctx, env.Repo, 1)

    // Reset timer AFTER setup
    b.ResetTimer()
    b.ReportAllocs()

    // Benchmark loop
    for i := 0; i < b.N; i++ {
        // Code to benchmark
        _, err := env.Repo.AddItemToCart(ctx, &domain.CartItem{
            CartID:   cart.ID,
            ListingID: listing.ID,
            Quantity: 1,
        })
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

### Available Benchmarks

**Location**: `/tests/performance/orders_benchmarks_test.go`

1. **BenchmarkAddToCart** - Cart item addition (Target: <50ms P95)
2. **BenchmarkGetCart** - Cart retrieval (Target: <20ms P95)
3. **BenchmarkGetCartWithItems** - Cart with items (Target: <50ms P95)
4. **BenchmarkCreateOrder** - Order creation (Target: <200ms P95)
5. **BenchmarkListOrders** - Order listing (Target: <100ms P95)
6. **BenchmarkConcurrentAddToCart** - Concurrent writes (Target: >100 ops/sec)
7. **BenchmarkConcurrentGetCart** - Concurrent reads (Target: >500 ops/sec)
8. **BenchmarkConcurrentCreateOrder** - Concurrent orders (Target: >50 ops/sec)
9. **BenchmarkMixedOrderOperations** - Realistic workload (70% reads, 20% writes, 10% orders)

### Important Benchmark Rules

✅ **DO**:
- Use `b.ResetTimer()` AFTER setup to exclude container startup time
- Setup TestEnvironment ONCE per benchmark
- Use `testing.Short()` to skip in CI
- Benchmark actual business logic, NOT Docker operations
- Use `b.ReportAllocs()` to track memory allocations

❌ **DON'T**:
- Create/destroy Docker containers inside `b.N` loop (will timeout CI)
- Construct `testing.T{}` manually (invalid, will panic)
- Benchmark environment setup (meaningless metrics)
- Run benchmarks without `-run=^$` (will also run unit tests)

### Benchmark Output Example

```
BenchmarkAddToCart-8                1000  1234567 ns/op  2048 B/op  32 allocs/op
BenchmarkGetCart-8                  2000   567890 ns/op  1024 B/op  16 allocs/op
BenchmarkConcurrentAddToCart-8      5000   234567 ns/op  4096 B/op  64 allocs/op
```

- **1000**: Number of iterations (b.N)
- **1234567 ns/op**: Average time per operation (nanoseconds)
- **2048 B/op**: Bytes allocated per operation
- **32 allocs/op**: Number of allocations per operation

## Test Coverage

Generate coverage report:

```bash
go test ./internal/transport/grpc/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## Best Practices

### 1. Test Isolation

✅ **DO**: Use separate test environment per test
```go
func TestExample(t *testing.T) {
    env := testingpkg.NewTestEnvironment(t)
    defer env.Cleanup()
}
```

❌ **DON'T**: Share environment across tests
```go
var globalEnv *testingpkg.TestEnvironment // Bad!
```

### 2. Data Cleanup

✅ **DO**: Rely on automatic cleanup
```go
defer env.Cleanup() // Handles everything
```

❌ **DON'T**: Manual cleanup (unnecessary)
```go
env.DB.Exec("DELETE FROM ...") // Not needed
```

### 3. Test Naming

✅ **DO**: Use descriptive names
```go
func TestAddToCart_Success(t *testing.T) {}
func TestAddToCart_InvalidQuantity(t *testing.T) {}
```

❌ **DON'T**: Generic names
```go
func TestCart1(t *testing.T) {}
func TestCart2(t *testing.T) {}
```

### 4. Assertions

✅ **DO**: Use testify assertions
```go
assert.NoError(t, err)
assert.Equal(t, expected, actual)
```

❌ **DON'T**: Manual error checking
```go
if err != nil {
    t.Errorf("error: %v", err)
}
```

## File Structure

```
/p/github.com/sveturs/listings/
├── internal/
│   ├── testing/
│   │   ├── environment.go          # TestEnvironment setup
│   │   ├── environment_test.go     # Environment unit tests
│   │   ├── database.go             # Database helpers
│   │   ├── fixtures.go             # Test data fixtures
│   │   ├── helpers.go              # Utility functions
│   │   └── E2E_TESTING_GUIDE.md    # This file
│   └── transport/
│       └── grpc/
│           └── handlers_orders_test.go  # 23 E2E tests
├── migrations/                      # Database migrations
└── go.mod
```

## Support & Contact

For questions or issues:
- Check this guide first
- Review test logs: `go test -v ...`
- Check Docker logs: `docker logs <container_id>`
- Create issue in project repository

---

**Last Updated**: 2025-11-14
**Version**: 1.0.0
**Author**: Phase 17 Implementation Team
