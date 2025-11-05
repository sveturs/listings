# Mocks for Listings Service Testing

This package contains mock implementations of all interfaces required by the listings service layer for unit testing.

## Mock Types

### MockRepository
Mock implementation of `listings.Repository` interface.
- Implements all repository methods (listings, products, variants, images, categories, favorites)
- Uses `testify/mock.Mock` for expectation setup and verification
- Supports all CRUD operations, bulk operations, and transaction management

### MockCacheRepository
Mock implementation of `listings.CacheRepository` interface.
- Implements Get, Set, Delete cache operations
- Useful for testing caching behavior and cache invalidation logic

### MockIndexingService
Mock implementation of `listings.IndexingService` interface.
- Implements IndexListing, UpdateListing, DeleteListing operations
- Useful for testing search indexing integration

## Usage Example

```go
func TestCreateListing_Success(t *testing.T) {
    // Setup
    service, mockRepo, mockCache, mockIndexer := SetupServiceTest(t)
    ctx := TestContext()

    input := NewCreateListingInput(100, "Test Listing")
    expectedListing := NewTestListing(1, 100, "Test Listing")

    // Setup expectations
    mockRepo.On("CreateListing", ctx, input).Return(expectedListing, nil)
    mockRepo.On("EnqueueIndexing", ctx, int64(1), domain.IndexOpIndex).Return(nil)

    // Execute
    result, err := service.CreateListing(ctx, input)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, expectedListing.ID, result.ID)

    // Verify expectations
    mockRepo.AssertExpectations(t)
}
```

## Available Test Helpers

See `../test_helpers.go` for a comprehensive list of helper functions:

### Setup Functions
- `SetupServiceTest(t)` - Creates service with all mocks
- `TestContext()` - Creates test context with timeout

### Listing Helpers
- `NewTestListing(id, userID, title)` - Create test listing
- `NewCreateListingInput(userID, title)` - Create listing input
- `NewUpdateListingInput(title, price)` - Update listing input
- `NewListListingsFilter(limit)` - List filter
- `NewSearchListingsQuery(query, limit)` - Search query

### Product Helpers
- `NewTestProduct(id, storefrontID, name)` - Create test product
- `NewCreateProductInput(storefrontID, name)` - Create product input
- `NewUpdateProductInput(name, price)` - Update product input
- `NewTestProductStats()` - Product statistics

### Variant Helpers
- `NewTestProductVariant(id, productID)` - Create test variant
- `NewCreateVariantInput(productID)` - Create variant input
- `NewUpdateVariantInput(price, stockQuantity)` - Update variant input

### Stock/Inventory Helpers
- `NewStockUpdateItem(productID, quantity)` - Stock update item
- `NewStockUpdateResult(productID, before, after, success)` - Stock update result

### Other Helpers
- `NewTestCategory(id, name)` - Create test category
- `NewTestListingImage(id, listingID)` - Create test image
- `NewBulkUpdateProductInput(productID, name, price)` - Bulk update input
- `NewBulkUpdateProductsResult(products, errors)` - Bulk update result

## Best Practices

1. **Always use SetupServiceTest()** - Creates properly initialized service with all dependencies
2. **Use TestContext()** - Provides timeout to prevent hanging tests
3. **Use helper functions** - Creates valid test data with proper defaults
4. **Setup expectations before execution** - Use `On()` and `Return()` methods
5. **Verify expectations** - Call `AssertExpectations(t)` at end of test
6. **Clean test data** - Each test should be independent

## Running Tests

```bash
# Run all service tests
go test -v ./internal/service/listings/

# Run specific test
go test -v ./internal/service/listings/ -run TestCreateListing_Success

# Run with coverage
go test -v ./internal/service/listings/ -cover

# Run with race detection
go test -v ./internal/service/listings/ -race
```

## Package Dependencies

- `github.com/stretchr/testify/mock` - Mock framework
- `github.com/rs/zerolog` - Logging (test writer)
- Project internal packages:
  - `github.com/sveturs/listings/internal/domain`
  - `github.com/sveturs/listings/internal/service/listings`
