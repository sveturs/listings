# Product CRUD Tests - Usage Guide

## Quick Start

### Run All Product Tests
```bash
cd /p/github.com/sveturs/listings
go test -v ./internal/repository/postgres/ -run "TestCreateProduct|TestUpdateProduct|TestDeleteProduct|TestBulkCreateProducts|TestBulkUpdateProducts|TestBulkDeleteProducts"
```

### Run Specific Test Category

**Create Tests:**
```bash
go test -v ./internal/repository/postgres/ -run TestCreateProduct
```

**Update Tests:**
```bash
go test -v ./internal/repository/postgres/ -run TestUpdateProduct
```

**Delete Tests:**
```bash
go test -v ./internal/repository/postgres/ -run TestDeleteProduct
```

**Bulk Operations:**
```bash
go test -v ./internal/repository/postgres/ -run TestBulk
```

### Run Single Test
```bash
go test -v ./internal/repository/postgres/ -run TestCreateProduct_Success
```

### Run with Coverage
```bash
go test -coverprofile=coverage.out ./internal/repository/postgres/ -run TestCreateProduct
go tool cover -html=coverage.out -o coverage.html
```

### Run in Short Mode (Skip Integration Tests)
```bash
go test -short ./internal/repository/postgres/
```

---

## Test Structure

Each test follows the AAA pattern:

```go
func TestCreateProduct_Success(t *testing.T) {
    // Arrange - Setup test data
    repo, testDB := setupTestRepo(t)
    defer testDB.TeardownTestPostgres(t)

    storefrontID := createTestStorefront(t, repo)
    product := &domain.CreateProductInput{ /* ... */ }

    // Act - Execute the operation
    createdProduct, err := repo.CreateProduct(ctx, product)

    // Assert - Verify results
    require.NoError(t, err)
    assert.Equal(t, expected, actual)
}
```

---

## Helper Functions

### Test Fixtures

**Create Test Storefront:**
```go
storefrontID := createTestStorefront(t, repo)
```

**Create Test Product:**
```go
product := createTestProduct(t, repo, storefrontID)
```

**Create Product with Custom Options:**
```go
product := createTestProductWithOptions(t, repo, storefrontID, "SKU-001", 99.99, 100)
```

**Create Product Variant:**
```go
variant := createTestVariant(t, repo, productID)
```

### Pointer Helpers

```go
sku := stringPtr("TEST-SKU")
price := float64Ptr(99.99)
quantity := int32Ptr(100)
```

---

## Test Data Patterns

### Valid Product Creation
```go
product := &domain.CreateProductInput{
    StorefrontID:  1,
    CategoryID:    100,
    Name:          "Test Product",
    Description:   "Test description",
    Price:         99.99,
    Currency:      "USD",
    SKU:           stringPtr("TEST-SKU-001"),
    StockQuantity: 100,
    Attributes:    map[string]interface{}{}, // Always include to avoid JSONB errors
}
```

### Testing Validation Errors
```go
// Empty name
product := &domain.CreateProductInput{
    Name: "",  // Should fail validation
    // ... other fields
}

// Negative price
product := &domain.CreateProductInput{
    Price: -10.00,  // Should fail validation
    // ... other fields
}
```

### Testing Bulk Operations
```go
inputs := []*domain.CreateProductInput{
    {Name: "Product 1", SKU: stringPtr("SKU-001"), /* ... */},
    {Name: "Product 2", SKU: stringPtr("SKU-002"), /* ... */},
    {Name: "Product 3", SKU: stringPtr("SKU-003"), /* ... */},
}

products, errors, err := repo.BulkCreateProducts(ctx, storefrontID, inputs)
```

---

## Common Test Scenarios

### Testing Success Cases
```go
func TestOperationName_Success(t *testing.T) {
    // Arrange
    repo, testDB := setupTestRepo(t)
    defer testDB.TeardownTestPostgres(t)

    // Act
    result, err := repo.Operation(ctx, input)

    // Assert
    require.NoError(t, err)
    assert.NotNil(t, result)
}
```

### Testing Error Cases
```go
func TestOperationName_NotFound(t *testing.T) {
    // Arrange
    repo, testDB := setupTestRepo(t)
    defer testDB.TeardownTestPostgres(t)

    // Act
    result, err := repo.Operation(ctx, invalidInput)

    // Assert
    assert.Error(t, err)
    assert.Nil(t, result)
    assert.Contains(t, err.Error(), "not_found")
}
```

### Testing Constraints
```go
func TestOperationName_DuplicateKey(t *testing.T) {
    // Arrange
    repo, testDB := setupTestRepo(t)
    defer testDB.TeardownTestPostgres(t)

    // Create first item
    _, err := repo.Create(ctx, input1)
    require.NoError(t, err)

    // Act - Try to create duplicate
    _, err = repo.Create(ctx, input2WithSameKey)

    // Assert
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "duplicate")
}
```

---

## Debugging Failed Tests

### View Detailed Logs
```bash
go test -v ./internal/repository/postgres/ -run TestCreateProduct_Success 2>&1 | less
```

### Check Database State
Tests use Docker containers. To inspect:
```bash
# Find running container
docker ps | grep postgres

# Connect to test DB (during test execution)
docker exec -it <container_id> psql -U test_user -d test_db

# View products
SELECT * FROM b2c_products;
```

### Enable Verbose Logging
The repository already logs at debug/info/error levels. Check test output for structured logs:
```json
{"level":"debug","component":"postgres_repository","product_id":1,"message":"creating product"}
```

---

## CI/CD Integration

### GitHub Actions Example
```yaml
- name: Run Product CRUD Tests
  run: |
    cd listings
    go test -v ./internal/repository/postgres/ \
      -run "TestCreateProduct|TestUpdateProduct|TestDeleteProduct|TestBulk" \
      -coverprofile=coverage.out

    go tool cover -func=coverage.out
```

### Pre-commit Hook
```bash
#!/bin/bash
# .git/hooks/pre-commit

cd /p/github.com/sveturs/listings
go test ./internal/repository/postgres/ -run TestCreateProduct -short
if [ $? -ne 0 ]; then
    echo "Tests failed! Commit aborted."
    exit 1
fi
```

---

## Known Issues & Workarounds

### Issue 1: JSONB Marshaling Error
**Error:** `pq: invalid input syntax for type json`

**Cause:** `Attributes` field is `nil` instead of empty map.

**Fix:** Always initialize Attributes:
```go
Attributes: map[string]interface{}{},  // Not nil!
```

### Issue 2: Soft Delete Not Working
**Error:** `GetProductByID` returns soft-deleted products

**Status:** Known issue, fix pending in repository code

**Workaround:** Use hard delete for now:
```go
repo.DeleteProduct(ctx, productID, storefrontID, true)  // hardDelete=true
```

### Issue 3: Transaction Rollback in BulkUpdate
**Error:** `failed to commit transaction: pq: Could not complete operation in a failed transaction`

**Status:** Known issue, fix pending

**Workaround:** Avoid testing duplicate SKU in bulk updates

---

## Performance Benchmarks

### Run Benchmarks
```bash
go test -bench=. -benchmem ./internal/repository/postgres/ -run=^$
```

### Example Benchmark
```go
func BenchmarkCreateProduct(b *testing.B) {
    repo, testDB := setupTestRepo(&testing.T{})
    defer testDB.TeardownTestPostgres(&testing.T{})

    storefrontID := createTestStorefront(&testing.T{}, repo)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        input := &domain.CreateProductInput{
            StorefrontID: storefrontID,
            Name: fmt.Sprintf("Product %d", i),
            SKU: stringPtr(fmt.Sprintf("SKU-%d", i)),
            Price: 99.99,
            // ...
        }
        repo.CreateProduct(context.Background(), input)
    }
}
```

---

## Best Practices

1. **Always cleanup:** Use `defer testDB.TeardownTestPostgres(t)`
2. **Use helpers:** Don't repeat fixture creation code
3. **Test isolation:** Each test should be independent
4. **Clear naming:** `TestOperation_Scenario` format
5. **Document intent:** Add comments for non-obvious assertions
6. **Verify DB state:** Don't just check return values
7. **Test errors:** Error paths are as important as success paths
8. **Use require vs assert:**
   - `require`: Test stops on failure (setup phase)
   - `assert`: Test continues on failure (assertions phase)

---

## Resources

- **Test File:** `internal/repository/postgres/products_test.go`
- **Summary:** `PRODUCT_CRUD_TESTS_SUMMARY.md`
- **Repository Code:** `internal/repository/postgres/products_repository.go`
- **Domain Models:** `internal/domain/product.go`
- **Test Infrastructure:** `tests/testing.go`

---

**Last Updated:** 2025-11-05
**Coverage:** 63.1% average for Product CRUD operations
**Test Count:** 38 tests (35 passing, 3 failing)
