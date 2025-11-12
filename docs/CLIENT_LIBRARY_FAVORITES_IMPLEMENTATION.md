# Client Library - Favorites Implementation

## Summary

Successfully implemented client library methods for Favorites functionality in the listings microservice. All methods support both gRPC (primary) and HTTP fallback (secondary) transports.

## Implementation Date
2025-11-11

## Files Modified

### 1. `/pkg/service/client.go` (432 lines)
**Changes:**
- Added proto imports: `pb "github.com/sveturs/listings/api/proto/listings/v1"`
- Updated `Client` struct to use `pb.ListingsServiceClient` instead of `interface{}`
- Updated `NewClient()` to initialize `pb.NewListingsServiceClient(conn)`
- Added 5 public methods for Favorites functionality:
  - `AddToFavorites(ctx, userID, listingID) error`
  - `RemoveFromFavorites(ctx, userID, listingID) error`
  - `GetUserFavorites(ctx, userID) ([]int64, int, error)`
  - `IsFavorite(ctx, userID, listingID) (bool, error)`
  - `GetFavoritedUsers(ctx, listingID) ([]int64, error)`

Each method follows the existing pattern:
1. Try gRPC first
2. Log errors and check if fallback is needed
3. Fallback to HTTP if enabled
4. Return `ErrUnavailable` if no transport available

### 2. `/pkg/service/grpc_client.go` (172 lines)
**Changes:**
- Uncommented proto imports
- Added 5 gRPC implementation methods:
  - `addToFavoritesGRPC(ctx, userID, listingID) error`
  - `removeFromFavoritesGRPC(ctx, userID, listingID) error`
  - `getUserFavoritesGRPC(ctx, userID) ([]int64, int, error)`
  - `isFavoriteGRPC(ctx, userID, listingID) (bool, error)`
  - `getFavoritedUsersGRPC(ctx, listingID) ([]int64, error)`

Each method:
- Uses context with timeout
- Creates appropriate proto request
- Calls gRPC client method
- Converts gRPC errors using `convertGRPCError()`
- Logs debug information

### 3. `/pkg/service/http_client.go` (492 lines)
**Changes:**
- Added 5 HTTP fallback implementation methods:
  - `AddToFavorites(ctx, userID, listingID) error`
  - `RemoveFromFavorites(ctx, userID, listingID) error`
  - `GetUserFavorites(ctx, userID) ([]int64, int, error)`
  - `IsFavorite(ctx, userID, listingID) (bool, error)`
  - `GetFavoritedUsers(ctx, listingID) ([]int64, error)`

HTTP endpoints used:
- `POST /api/v1/favorites/{listingId}` - AddToFavorites
- `DELETE /api/v1/favorites/{listingId}?user_id={userId}` - RemoveFromFavorites
- `GET /api/v1/users/{userId}/favorites` - GetUserFavorites
- `GET /api/v1/favorites/{listingId}/is-favorite?user_id={userId}` - IsFavorite
- `GET /api/v1/favorites/{listingId}/users` - GetFavoritedUsers

## Proto Messages Used

From `/api/proto/listings/v1/listings.proto`:

```protobuf
// Requests
message AddToFavoritesRequest {
  int64 user_id = 1;
  int64 listing_id = 2;
}

message RemoveFromFavoritesRequest {
  int64 user_id = 1;
  int64 listing_id = 2;
}

message GetUserFavoritesRequest {
  int64 user_id = 1;
  int32 limit = 2;
  int32 offset = 3;
}

message IsFavoriteRequest {
  int64 user_id = 1;
  int64 listing_id = 2;
}

message ListingIDRequest {
  int64 listing_id = 1;
}

// Responses
message GetUserFavoritesResponse {
  repeated int64 listing_ids = 1;
  int32 total = 2;
}

message IsFavoriteResponse {
  bool is_favorite = 1;
}

message UserIDsResponse {
  repeated int64 user_ids = 1;
}
```

## Testing

### Compilation Test
```bash
cd /p/github.com/sveturs/listings
go build ./pkg/service/...
# ✅ PASS
```

### Unit Tests
```bash
cd /p/github.com/sveturs/listings
go test ./pkg/service/... -v
# ✅ PASS - All 10 tests passed
```

### Go Vet
```bash
cd /p/github.com/sveturs/listings
go vet ./pkg/service/...
# ✅ PASS - No issues found
```

## Usage Example

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/rs/zerolog"
    "github.com/sveturs/listings/pkg/service"
)

func main() {
    logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

    // Create client with gRPC and HTTP fallback
    client, err := service.NewClient(service.ClientConfig{
        GRPCAddr:       "localhost:50053",
        HTTPBaseURL:    "http://localhost:8086",
        AuthToken:      "your-service-token",
        Timeout:        5 * time.Second,
        EnableFallback: true,
        Logger:         logger,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    ctx := context.Background()
    userID := int64(1)
    listingID := int64(328)

    // Add to favorites
    err = client.AddToFavorites(ctx, userID, listingID)
    if err != nil {
        log.Printf("Failed to add to favorites: %v", err)
    }

    // Check if favorited
    isFav, err := client.IsFavorite(ctx, userID, listingID)
    if err != nil {
        log.Printf("Failed to check favorite: %v", err)
    }
    log.Printf("Is favorite: %v", isFav)

    // Get user's favorites
    listingIDs, total, err := client.GetUserFavorites(ctx, userID)
    if err != nil {
        log.Printf("Failed to get favorites: %v", err)
    }
    log.Printf("User has %d favorites: %v", total, listingIDs)

    // Get users who favorited a listing
    userIDs, err := client.GetFavoritedUsers(ctx, listingID)
    if err != nil {
        log.Printf("Failed to get favorited users: %v", err)
    }
    log.Printf("Listing favorited by %d users: %v", len(userIDs), userIDs)

    // Remove from favorites
    err = client.RemoveFromFavorites(ctx, userID, listingID)
    if err != nil {
        log.Printf("Failed to remove from favorites: %v", err)
    }
}
```

## Integration with Monolith

To integrate these methods into the monolith (`/p/github.com/sveturs/svetu`):

### 1. Import the Client Library
```bash
cd /p/github.com/sveturs/svetu/backend
go get github.com/sveturs/listings/pkg/service@latest
```

### 2. Initialize Client in Server Startup
```go
// In backend/internal/server/server.go or similar
listingsClient, err := service.NewClient(service.ClientConfig{
    GRPCAddr:       os.Getenv("LISTINGS_GRPC_ADDR"),      // e.g., "localhost:50053"
    HTTPBaseURL:    os.Getenv("LISTINGS_HTTP_BASE_URL"),  // e.g., "http://localhost:8086"
    AuthToken:      os.Getenv("SERVICE_TO_SERVICE_TOKEN"),
    Timeout:        5 * time.Second,
    EnableFallback: true,
    Logger:         logger,
})
```

### 3. Use in Handlers
```go
// In marketplace/favorites handler
func (h *Handler) AddToFavorites(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(int64)
    listingID, _ := c.ParamsInt("listing_id")

    err := h.listingsClient.AddToFavorites(c.Context(), userID, int64(listingID))
    if err != nil {
        if errors.Is(err, service.ErrNotFound) {
            return c.Status(404).JSON(fiber.Map{"error": "listing not found"})
        }
        return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
    }

    return c.Status(200).JSON(fiber.Map{"success": true})
}
```

### 4. Environment Variables
Add to `.env`:
```bash
LISTINGS_GRPC_ADDR=localhost:50053
LISTINGS_HTTP_BASE_URL=http://localhost:8086
SERVICE_TO_SERVICE_TOKEN=your-secure-token
```

## Notes

1. **gRPC vs HTTP**: gRPC is always tried first. HTTP is used as fallback only when:
   - gRPC returns transient errors (Unavailable, DeadlineExceeded, Canceled, Unknown)
   - `EnableFallback` is set to `true`
   - `httpClient` is initialized

2. **Error Handling**: All errors are properly converted:
   - `codes.NotFound` → `service.ErrNotFound`
   - `codes.InvalidArgument` → `service.ErrInvalidInput`
   - `codes.Unavailable` → `service.ErrUnavailable`

3. **Logging**: All methods include debug logging with:
   - Request parameters (userID, listingID)
   - Result counts
   - Error details

4. **Timeout**: Each gRPC call uses context with timeout from `ClientConfig.Timeout` (default 5s)

5. **Thread Safety**: Client is thread-safe and can be shared across goroutines

## What's Ready

✅ Proto definitions exist and are compiled
✅ gRPC handlers implemented in listings microservice
✅ Service layer implemented in listings microservice
✅ Client library methods implemented (gRPC + HTTP)
✅ All tests passing
✅ Code compiles without errors
✅ Proper error handling and logging

## What's Needed for Integration

1. **Deploy listings microservice** with gRPC server running on port 50053
2. **Configure monolith** with environment variables for listings service connection
3. **Update monolith handlers** to use client library instead of direct DB access
4. **Add service-to-service authentication token** for secure communication
5. **Test end-to-end flow** with both gRPC and HTTP fallback

## Status

**Status:** ✅ READY FOR INTEGRATION

The client library is fully implemented, tested, and ready to be integrated into the monolith. All Favorites functionality can now be accessed through the unified client interface with automatic gRPC/HTTP fallback.
