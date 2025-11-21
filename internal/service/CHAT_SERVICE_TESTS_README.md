# Chat Service Test Suite

## Overview

This test suite provides comprehensive unit tests for the Chat functionality in the listings microservice. Tests are implemented using the testify/mock framework and follow Go testing best practices.

## Test Coverage

### Service Layer Tests (`chat_service_test.go`)

**Total Tests**: 13 passing tests

**Coverage Breakdown** (for chat_service.go methods):
- GetChat: 80.0%
- GetUnreadCount: 80.0%
- GetUserChats: 73.3%
- MarkMessagesAsRead: 70.6%
- validateGetMessagesRequest: 57.1%
- GetMessages: 64.0%
- ArchiveChat: 63.6%
- DeleteChat: 54.5%

### Tested Methods

#### ✅ GetChat
- Success scenario with unread count
- Not participant error
- Chat not found error

#### ✅ GetUserChats
- Success with multiple chats and unread counts
- Pagination limit validation (max 100)

#### ✅ GetMessages
- Success with message retrieval
- Not participant error

#### ✅ MarkMessagesAsRead
- Mark specific messages
- Mark all messages

#### ✅ GetUnreadCount
- Specific chat count
- All chats count

#### ✅ ArchiveChat
- Success scenario

#### ✅ DeleteChat
- Success scenario

### Not Tested (Limitations)

The following methods are NOT tested due to architectural constraints:

#### ❌ CreateChat / GetOrCreateChat
**Reason**: Requires mocking `*postgres.Repository` (concrete type, not interface)
**Workaround**: Would need either:
1. Refactoring productsRepo to use an interface
2. Using reflection/unsafe to inject mocks
3. Integration tests with real database

#### ❌ SendMessage
**Reason**: Uses database transactions (pgx.Tx) which are difficult to mock in unit tests
**Workaround**: Would benefit from integration tests

#### ❌ Attachment Operations
**Reason**: Depends on file upload/storage logic and SendMessage (transactions)
**Workaround**: Integration or E2E tests

## Running Tests

### Run all chat service tests:
```bash
cd /p/github.com/sveturs/listings
go test ./internal/service -v -run "TestChatService_"
```

### Run with coverage:
```bash
go test ./internal/service -cover -run "TestChatService_"
```

### Generate coverage report:
```bash
go test ./internal/service -coverprofile=/tmp/coverage.out -run "TestChatService_"
go tool cover -html=/tmp/coverage.out -o /tmp/coverage.html
```

## Test Structure

### Mock Repositories

All repository dependencies are mocked using testify/mock:
- `MockChatRepository` - Chat CRUD operations
- `MockMessageRepository` - Message operations
- `MockAttachmentRepository` - Attachment operations
- `MockProductsRepository` - Product/listing validation (unused due to limitations)

### Test Pattern

Tests follow the AAA (Arrange-Act-Assert) pattern:

```go
func TestChatService_MethodName_Scenario(t *testing.T) {
    // Arrange
    service, chatRepo, messageRepo, _, _ := setupTestChatService(t)
    ctx := context.Background()

    // Setup mocks
    chatRepo.On("Method", args...).Return(result, nil)

    // Act
    result, err := service.Method(ctx, request)

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expected, result)

    chatRepo.AssertExpectations(t)
}
```

## Known Issues

1. **Limited Coverage of CreateChat**: Cannot test listing validation due to concrete type dependency
2. **No Transaction Tests**: Methods using transactions (SendMessage) are skipped
3. **Attachment Tests Missing**: Would require file upload mocking and transaction support

## Recommendations

1. **Refactor ProductsRepo**: Change from `*postgres.Repository` to interface for better testability
2. **Integration Tests**: Add tests with real PostgreSQL database using testcontainers
3. **E2E Tests**: Test full chat flow including file uploads via gRPC API

## Architecture Note

The service has a limitation where `productsRepo` is typed as `*postgres.Repository` (concrete type) rather than an interface. This prevents easy mocking and limits unit test coverage for methods that validate listings/products.

### Future Improvement

Extract interface:
```go
type ProductRepository interface {
    GetListingByID(ctx context.Context, listingID int64) (*domain.Listing, error)
}
```

Then inject in chatService:
```go
type chatService struct {
    // ...
    productsRepo   ProductRepository // Changed from *postgres.Repository
    // ...
}
```

This would allow full mock coverage of all service methods.

## Test File Location

- **Service Tests**: `/p/github.com/sveturs/listings/internal/service/chat_service_test.go`
- **Handler Tests**: `/p/github.com/sveturs/listings/internal/transport/grpc/handlers_chat_test.go` (NOT YET CREATED)

## Dependencies

- `github.com/stretchr/testify` - Assertions and mocking
- `github.com/rs/zerolog` - Logging (disabled in tests)
- `github.com/jackc/pgx/v5` - PostgreSQL driver types

---

**Last Updated**: 2025-11-21
**Test Coverage**: 13 tests, 63.6-80% coverage on tested methods
**Status**: ✅ All service tests passing
