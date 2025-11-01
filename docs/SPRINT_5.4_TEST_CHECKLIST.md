# Sprint 5.4 Testing Checklist

**Based on Sprint 5.3 Verification Report**
**Target Coverage:** 65%+ (from current 29.8%)

---

## Phase 1: Quick Wins (Day 1-2) ⚡

### Code Quality Fixes
- [ ] Run `gofmt -w internal/transport/grpc/converters.go`
- [ ] Run `gofmt -w internal/transport/grpc/handlers_test.go`
- [ ] Extract magic numbers to constants
  - [ ] Create `internal/transport/grpc/constants.go`
  - [ ] Move MaxLimit = 100
  - [ ] Move MinTitleLength = 3
  - [ ] Move MaxTitleLength = 255
  - [ ] Move MinQueryLength = 2
  - [ ] Move CurrencyCodeLength = 3

### Typed Errors (30 min)
- [ ] Create `internal/domain/errors.go`
- [ ] Add `ErrUnauthorized = errors.New("unauthorized")`
- [ ] Add `ErrNotFound = errors.New("not found")`
- [ ] Update handlers.go to use `errors.Is()` instead of string comparison

### Missing Converter Tests (2 hours)
- [ ] Add `TestDomainToProtoAttribute()` in handlers_test.go
- [ ] Add `TestDomainToProtoLocation()` in handlers_test.go
- [ ] Add `TestProtoToUpdateListingInput()` in handlers_test.go
- [ ] Add `TestProtoToListListingsFilter()` in handlers_test.go
- [ ] Add `TestProtoToSearchListingsQuery()` in handlers_test.go

**Expected Impact:** Coverage 29.8% → 40%

---

## Phase 2: Handler Success Path Tests (Day 3-4) 🎯

### Complete Mock Service Integration
- [ ] Fix mock service wiring in `setupTestServer()` (handlers_test.go:69-81)
- [ ] Ensure all tests use mock service correctly

### Test GetListing Success Path
- [ ] Test GetListing with valid ID and existing listing
- [ ] Test GetListing with images attached
- [ ] Test GetListing with attributes attached
- [ ] Test GetListing with location attached

### Test CreateListing Success Path
- [ ] Test CreateListing with minimal required fields
- [ ] Test CreateListing with all optional fields
- [ ] Test CreateListing with storefront_id
- [ ] Test CreateListing with description and SKU

### Test UpdateListing Success Path
- [ ] Test UpdateListing with title only
- [ ] Test UpdateListing with multiple fields
- [ ] Test UpdateListing with status change
- [ ] Test UpdateListing authorization check

### Test DeleteListing Success Path
- [ ] Test DeleteListing with valid ID and owner
- [ ] Test DeleteListing authorization check
- [ ] Test DeleteListing with non-existent listing

### Test SearchListings Success Path
- [ ] Test SearchListings with query only
- [ ] Test SearchListings with category filter
- [ ] Test SearchListings with price range
- [ ] Test SearchListings with pagination
- [ ] Test SearchListings with no results

### Test ListListings Success Path
- [ ] Test ListListings with default parameters
- [ ] Test ListListings with user_id filter
- [ ] Test ListListings with storefront_id filter
- [ ] Test ListListings with status filter
- [ ] Test ListListings with price range
- [ ] Test ListListings pagination (offset/limit)

**Expected Impact:** Coverage 40% → 60%

---

## Phase 3: Service Error Scenarios (Day 5) ⚠️

### Create `TestHandlerErrors` Test Suite
- [ ] Test GetListing when service returns error
- [ ] Test CreateListing when service returns error
- [ ] Test UpdateListing when service returns error
- [ ] Test UpdateListing when user doesn't own listing (PermissionDenied)
- [ ] Test DeleteListing when service returns error
- [ ] Test DeleteListing when user doesn't own listing (PermissionDenied)
- [ ] Test SearchListings when service returns error
- [ ] Test ListListings when service returns error

### Verify gRPC Status Codes
- [ ] Validate InvalidArgument code for validation errors
- [ ] Validate NotFound code for missing resources
- [ ] Validate PermissionDenied code for authorization failures
- [ ] Validate Internal code for service errors

**Expected Impact:** Coverage 60% → 65%

---

## Phase 4: Integration Tests Setup (Day 6-8) 🗄️

### Docker Test Infrastructure
- [ ] Create `test/docker-compose.test.yml`
  - [ ] PostgreSQL test service (port 5433)
  - [ ] OpenSearch test service (port 9201)
- [ ] Create `internal/transport/grpc/testutil/` package
- [ ] Add `SetupTestDB(t *testing.T) *sql.DB` helper
- [ ] Add `TeardownTestDB(t *testing.T, db *sql.DB)` helper
- [ ] Add `CreateTestListing()` fixture helper
- [ ] Add `CleanupTestData()` helper

### Database Integration Tests
- [ ] Create `internal/transport/grpc/integration_test.go`
- [ ] Add build tag: `// +build integration`
- [ ] Test GetListing with real DB
- [ ] Test CreateListing with real DB (verify inserted data)
- [ ] Test UpdateListing with real DB (verify updated data)
- [ ] Test DeleteListing with real DB (verify soft delete)
- [ ] Test ListListings with real DB (verify pagination)

### CRUD Flow Integration Test
- [ ] Test full flow: Create → Get → Update → Get → Delete → Get (NotFound)
- [ ] Verify data consistency at each step
- [ ] Verify transactions work correctly

**Expected Impact:** High confidence in production readiness

---

## Phase 5: Advanced Tests (Optional - Day 9-10) 🚀

### Search Integration Tests (if OpenSearch ready)
- [ ] Setup OpenSearch in test environment
- [ ] Test full-text search with indexed data
- [ ] Test category filtering
- [ ] Test price range filtering
- [ ] Test pagination in search results

### E2E Tests with gRPC Client
- [ ] Create `test/e2e/` directory
- [ ] Setup real gRPC client in tests
- [ ] Start/stop server automatically
- [ ] Test full request/response cycle
- [ ] Test concurrent requests (10-100 simultaneous)

### Performance Benchmarks
- [ ] Create `internal/transport/grpc/handlers_bench_test.go`
- [ ] Add `BenchmarkGetListing(b *testing.B)`
- [ ] Add `BenchmarkCreateListing(b *testing.B)`
- [ ] Add `BenchmarkDomainToProtoListing(b *testing.B)`
- [ ] Add `BenchmarkSearchListings(b *testing.B)`
- [ ] Document baseline performance (README.md)

**Expected Impact:** Production confidence + performance baseline

---

## Acceptance Criteria

### Must Have (Sprint 5.4 Success)
- ✅ All converter functions have tests (100% coverage)
- ✅ All handler functions have success path tests (>50% coverage)
- ✅ Service error scenarios tested (>80% coverage)
- ✅ Overall code coverage ≥65%
- ✅ All existing tests still pass
- ✅ Code formatting fixed
- ✅ Typed errors implemented

### Should Have
- ✅ Database integration tests running
- ✅ CRUD flow integration test passing
- ✅ Docker test infrastructure working

### Nice to Have
- ⚠️ Search integration tests (depends on OpenSearch availability)
- ⚠️ E2E tests with real gRPC client
- ⚠️ Performance benchmarks documented

---

## Running Tests

### Unit Tests Only
```bash
go test -v -short ./internal/transport/grpc/...
```

### With Coverage
```bash
go test -v -short -coverprofile=coverage.out ./internal/transport/grpc/...
go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Integration Tests (with Docker)
```bash
# Start test infrastructure
docker-compose -f test/docker-compose.test.yml up -d

# Run integration tests
go test -v -tags=integration ./internal/transport/grpc/...

# Cleanup
docker-compose -f test/docker-compose.test.yml down -v
```

### All Tests
```bash
go test -v ./internal/transport/grpc/...
```

---

## Validation Commands

Before marking Sprint 5.4 as complete:

```bash
# 1. Code quality
gofmt -l internal/transport/grpc/  # Should be empty
go vet ./internal/transport/grpc/...  # Should pass

# 2. Test execution
go test -v -short ./internal/transport/grpc/...  # All pass

# 3. Coverage check
go test -short -coverprofile=/tmp/coverage.out ./internal/transport/grpc/...
go tool cover -func=/tmp/coverage.out | tail -1
# Should show: total: (statements) ≥65.0%

# 4. Build verification
go build -v ./cmd/server  # Should succeed

# 5. Binary execution
./server --help  # Should start and fail gracefully (no DB)
```

---

## Notes for Test Engineer

1. **Prioritize based on risk**: Handler tests are more important than converter tests.

2. **Mock service properly**: Fix the mock integration first - many tests will be easier afterwards.

3. **Use table-driven tests**: Keep the pattern from existing tests - it's clean and maintainable.

4. **Test one thing at a time**: Each test should verify one specific behavior.

5. **Don't skip cleanup**: Always cleanup test data, even if tests fail.

6. **Document edge cases**: If you find an edge case during testing, document it.

7. **Measure, don't guess**: Use coverage tools to verify actual coverage, don't estimate.

8. **Integration tests are expensive**: They're slower and more fragile. Balance unit vs integration.

9. **Performance baselines**: Document baseline performance for future comparison.

10. **Keep tests fast**: Unit tests should run in <1s, integration in <10s.

---

**Checklist Version:** 1.0
**Created:** 2025-11-01
**Based on:** Sprint 5.3 Verification Report
