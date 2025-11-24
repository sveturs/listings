# Test Infrastructure Fix Checklist

**Date**: 2025-11-14
**Estimated Time**: 4 hours
**Priority**: CRITICAL (blocks production deployment)

---

## Quick Status

- ✅ **Compilation Fixed** - All code compiles
- ❌ **Tests Broken** - Design issues prevent execution
- ⚠️ **CI/CD Ready** - Will fail due to test issues

---

## Fix #1: Update TestEnvironment Interface (30 minutes)

**Priority**: 🔴 CRITICAL
**File**: `/internal/testing/environment.go`

### Current Problem:
```go
func NewTestEnvironment(t *testing.T) *TestEnvironment  // ❌ Only works with tests
```

### Fix:
```go
func NewTestEnvironment(t testing.TB) *TestEnvironment  // ✅ Works with tests AND benchmarks
func NewTestEnvironmentWithConfig(t testing.TB, config TestEnvironmentConfig) *TestEnvironment
```

### Changes Required:
1. Import `testing` package (already imported)
2. Change parameter type from `*testing.T` to `testing.TB`
3. Update all internal functions that use `t` (they should already work since TB is interface)

### Test:
```bash
go test -v ./internal/testing/ -run TestNewTestEnvironment_Docker -timeout 5m
```

---

## Fix #2: Rewrite Performance Benchmarks (2 hours)

**Priority**: 🔴 CRITICAL
**File**: `/tests/performance/orders_benchmarks_test.go`

### Current Problem:
```go
func BenchmarkAddToCart(b *testing.B) {
    tests.SkipIfNoDocker(&testing.T{})           // ❌ WRONG! Creates empty struct
    testDB := tests.SetupTestPostgres(&testing.T{}) // ❌ WRONG!

    for i := 0; i < b.N; i++ {
        // benchmark code
    }
}
```

### Fix Pattern:
```go
func BenchmarkAddToCart(b *testing.B) {
    if testing.Short() {
        b.Skip("Skipping benchmark in short mode")
    }

    // Setup ONCE (outside loop)
    env := testing.NewTestEnvironment(b)  // ✅ Now works because TB interface
    defer env.Cleanup()

    ctx := context.Background()

    // Create test data ONCE
    cart := createTestCart(b, ctx, env.Repo, 1)
    listing := createTestListing(b, ctx, env.Repo)

    // Benchmark actual operation
    b.ResetTimer()
    b.ReportAllocs()

    for i := 0; i < b.N; i++ {
        item := &domain.CartItem{
            CartID:        cart.ID,
            ListingID:     listing.ID,
            Quantity:      1,
            PriceSnapshot: listing.Price,
        }

        _, err := env.Repo.AddItemToCart(ctx, item)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

### Functions to Fix:

1. ✅ Update helper functions to accept `testing.TB`:
   ```go
   func createTestListing(t testing.TB, ctx context.Context, repo *postgres.Repository) *domain.Listing
   func createTestCart(t testing.TB, ctx context.Context, repo *postgres.Repository, storefrontID int64) *domain.Cart
   func addTestItemToCart(t testing.TB, ctx context.Context, repo *postgres.Repository, cartID, listingID int64, price float64)
   ```

2. ❌ Remove all `&testing.T{}` constructs

3. ✅ Move setup outside `b.N` loop

4. ✅ Add `b.ResetTimer()` before benchmark loop

### Benchmarks to Fix:
- [ ] BenchmarkAddToCart
- [ ] BenchmarkGetCart
- [ ] BenchmarkUpdateCartItem
- [ ] BenchmarkRemoveCartItem
- [ ] BenchmarkClearCart
- [ ] BenchmarkCreateOrder
- [ ] BenchmarkGetOrder
- [ ] BenchmarkListOrders
- [ ] BenchmarkCancelOrder
- [ ] BenchmarkGetOrderWithItems
- [ ] BenchmarkBulkAddToCart
- [ ] BenchmarkConcurrentCartOperations

### Test Each Benchmark:
```bash
# Test one benchmark
go test -bench=BenchmarkAddToCart -benchmem -benchtime=3s -run=^$ ./tests/performance/

# Test all benchmarks
go test -bench=. -benchmem -benchtime=3s -run=^$ ./tests/performance/
```

---

## Fix #3: Fix Environment Benchmark (15 minutes)

**Priority**: 🔴 CRITICAL
**File**: `/internal/testing/environment_test.go`

### Current Problem:
```go
func BenchmarkEnvironmentCreation(b *testing.B) {
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        env := testingpkg.NewTestEnvironment(b)  // ❌ Creates Docker containers in loop!
        env.Cleanup()
    }
}
```

**Why This is Wrong**:
- Benchmarks infrastructure setup (Docker containers), not code
- Each iteration takes ~30 seconds (Docker startup)
- Will timeout in CI/CD
- Not measuring anything useful

### Option A: Remove It (RECOMMENDED)
```go
// REMOVED: BenchmarkEnvironmentCreation
// Reason: Benchmarking infrastructure setup is not useful
// Use: Integration tests to verify environment works
```

### Option B: Disable It
```go
func BenchmarkEnvironmentCreation(b *testing.B) {
    b.Skip("TODO: Benchmark design needs rethinking - measures infra setup, not code")
}
```

### Option C: Benchmark Something Meaningful
```go
func BenchmarkDatabaseQuery(b *testing.B) {
    // Setup ONCE
    env := testingpkg.NewTestEnvironment(b)
    defer env.Cleanup()

    ctx := context.Background()

    // Benchmark actual database operation
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := env.Repo.GetListing(ctx, 1)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

### Test:
```bash
go test -bench=BenchmarkEnvironmentCreation -run=^$ ./internal/testing/
```

---

## Fix #4: Update Documentation (30 minutes)

**Priority**: 🟡 HIGH
**File**: `/internal/testing/E2E_TESTING_GUIDE.md`

### Updates Required:

1. **Add section on testing.TB interface**:
   ```markdown
   ## Using TestEnvironment with Benchmarks

   The `NewTestEnvironment()` function now accepts `testing.TB` interface,
   which means it works with both tests (`*testing.T`) and benchmarks (`*testing.B`):

   ```go
   // In tests
   func TestAddToCart(t *testing.T) {
       env := testing.NewTestEnvironment(t)  // ✅ Works
       defer env.Cleanup()
   }

   // In benchmarks
   func BenchmarkAddToCart(b *testing.B) {
       env := testing.NewTestEnvironment(b)  // ✅ Also works!
       defer env.Cleanup()
   }
   ```
   ```

2. **Add section on pgxpool requirement**:
   ```markdown
   ## Database Connections

   TestEnvironment creates TWO database connections:

   1. `env.DB` - *sqlx.DB connection (for main Repository)
   2. `env.PgPool` - *pgxpool.Pool connection (for Order/Reservation repos)

   This dual-connection approach supports both legacy (sqlx) and new (pgxpool)
   repository patterns.
   ```

3. **Add section on benchmark best practices**:
   ```markdown
   ## Writing Benchmarks

   ✅ **DO**:
   - Setup infrastructure ONCE (outside b.N loop)
   - Use `b.ResetTimer()` before benchmark loop
   - Benchmark actual operations, not infrastructure
   - Use `testing.Short()` to skip slow benchmarks

   ❌ **DON'T**:
   - Create Docker containers in benchmark loop
   - Use `&testing.T{}` constructs
   - Benchmark infrastructure setup
   - Forget to call `defer env.Cleanup()`
   ```

4. **Update examples** to show pgxpool usage

---

## Fix #5: Test CI/CD Locally (30 minutes)

**Priority**: 🟡 HIGH

### Run What CI Will Run:

```bash
# 1. Linting
cd /p/github.com/sveturs/listings
go vet ./...

# 2. Unit tests with coverage (what CI runs)
go test -v -race -coverprofile=coverage.txt -covermode=atomic -short \
  ./internal/... \
  ./pkg/...

# 3. Check coverage threshold
COVERAGE=$(go tool cover -func=coverage.txt | grep total | awk '{print $3}' | sed 's/%//')
echo "Coverage: $COVERAGE%"
if (( $(echo "$COVERAGE < 80" | bc -l) )); then
  echo "❌ Coverage below 80%"
else
  echo "✅ Coverage meets threshold"
fi

# 4. Build
make build

# 5. Integration tests (if you have Docker and time)
go test -v -timeout=10m -tags=integration \
  ./internal/transport/grpc/... \
  ./test/integration/...

# 6. Performance benchmarks
go test -bench=. -benchmem -benchtime=10s \
  -run=^$ \
  ./tests/performance/...
```

### Expected Results:
- ✅ All tests pass
- ✅ Coverage > 80%
- ✅ Build succeeds
- ✅ Benchmarks run without errors

---

## Verification Checklist

After completing all fixes, verify:

### Compilation
- [ ] `go build ./internal/testing/...` - SUCCESS
- [ ] `go build ./tests/performance/...` - SUCCESS
- [ ] `go build ./internal/transport/grpc/...` - SUCCESS
- [ ] No compilation warnings

### Unit Tests
- [ ] `go test -short ./internal/testing/` - PASS
- [ ] `go test ./internal/testing/ -timeout 15m` - PASS (with Docker)
- [ ] All 11 environment tests pass

### Performance Benchmarks
- [ ] `go test -bench=BenchmarkAddToCart -benchtime=3s ./tests/performance/` - RUNS
- [ ] All 12 benchmarks run successfully
- [ ] No panics or fatal errors
- [ ] Reasonable performance (< 100ms per op for most)

### E2E Tests
- [ ] `go test -v ./internal/transport/grpc/... -timeout 15m` - PASS
- [ ] All 23 E2E tests pass
- [ ] No Docker container cleanup issues

### CI/CD Simulation
- [ ] Lint passes
- [ ] Unit tests pass with `-short` flag
- [ ] Coverage > 80%
- [ ] Build succeeds
- [ ] Integration tests pass (optional, if Docker available)

### Documentation
- [ ] E2E_TESTING_GUIDE.md updated
- [ ] Examples updated
- [ ] Changelog/commit message describes changes

---

## Estimated Time Breakdown

| Task | Time | Priority |
|------|------|----------|
| Fix #1: TB Interface | 30 min | 🔴 CRITICAL |
| Fix #2: Benchmarks | 2 hours | 🔴 CRITICAL |
| Fix #3: Environment Benchmark | 15 min | 🔴 CRITICAL |
| Fix #4: Documentation | 30 min | 🟡 HIGH |
| Fix #5: CI/CD Testing | 30 min | 🟡 HIGH |
| **TOTAL** | **4 hours** | |

---

## Quick Win Option (30 minutes)

If you need a QUICK fix to unblock CI/CD:

1. **Disable problematic tests** (5 min):
   ```go
   func BenchmarkAddToCart(b *testing.B) {
       b.Skip("TODO: Fix testing patterns - issue #XXX")
   }
   ```

2. **Fix TestEnvironment interface** (25 min):
   ```go
   func NewTestEnvironment(t testing.TB) *TestEnvironment
   ```

3. **Test compilation**:
   ```bash
   go build ./...
   ```

This gets CI/CD green, then you can fix benchmarks properly later.

---

## Getting Help

- **Compilation errors**: Check `/p/github.com/sveturs/listings/PHASE_17_DAYS23-25_VALIDATION_REPORT.md` section 1
- **Test design questions**: See validation report section 2
- **CI/CD issues**: See validation report section 3
- **Dependencies**: See validation report section 4

---

## After Fixes Complete

1. Create commit:
   ```bash
   git add -A
   git commit -m "fix: resolve test infrastructure issues

   - Update TestEnvironment to use testing.TB interface
   - Rewrite performance benchmarks with correct patterns
   - Remove/fix problematic environment benchmark
   - Update documentation with pgxpool requirements

   Fixes compilation errors and test design issues identified
   in PHASE_17_DAYS23-25_VALIDATION_REPORT.md"
   ```

2. Push and verify CI/CD:
   ```bash
   git push origin <branch>
   ```

3. Monitor GitHub Actions to ensure all checks pass

4. Create PR with link to validation report

---

**Last Updated**: 2025-11-14
**Validator**: Test Engineer Agent
**Status**: Ready for implementation
