# Listings Service Client - Usage Examples

This document provides detailed usage examples for the Listings Service Client Library.

## Table of Contents

1. [Basic Client Setup](#basic-client-setup)
2. [CRUD Operations](#crud-operations)
3. [Search and Filtering](#search-and-filtering)
4. [Fiber Middleware Integration](#fiber-middleware-integration)
5. [gRPC Connection Pooling](#grpc-connection-pooling)
6. [Error Handling](#error-handling)
7. [Advanced Patterns](#advanced-patterns)

---

## Basic Client Setup

### Minimal Configuration (gRPC Only)

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

    client, err := service.NewClient(service.ClientConfig{
        GRPCAddr: "localhost:50053",
        Timeout:  5 * time.Second,
        Logger:   logger,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Use client...
}
```

### Production Configuration (with Fallback)

```go
client, err := service.NewClient(service.ClientConfig{
    GRPCAddr:       os.Getenv("LISTINGS_GRPC_ADDR"),
    HTTPBaseURL:    os.Getenv("LISTINGS_HTTP_URL"),
    AuthToken:      os.Getenv("SERVICE_TOKEN"),
    Timeout:        10 * time.Second,
    EnableFallback: true,  // Auto-fallback to HTTP if gRPC fails
    Logger:         logger,
})
if err != nil {
    return fmt.Errorf("failed to create listings client: %w", err)
}
```

---

## CRUD Operations

### Get Listing by ID

```go
func getListing(client *service.Client, id int64) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    listing, err := client.GetListing(ctx, id)
    if err != nil {
        if err == service.ErrNotFound {
            log.Printf("Listing %d not found", id)
            return nil
        }
        return fmt.Errorf("failed to get listing: %w", err)
    }

    log.Printf("Found listing: %s (Price: %.2f %s)",
        listing.Title,
        listing.Price,
        listing.Currency,
    )

    return nil
}
```

### Create Listing

```go
func createListing(client *service.Client, userID int64) (*service.Listing, error) {
    ctx := context.Background()

    description := "High-quality product in excellent condition"
    req := &service.CreateListingRequest{
        UserID:      userID,
        Title:       "Vintage Camera",
        Description: &description,
        Price:       299.99,
        Currency:    "USD",
        CategoryID:  1001,
        Quantity:    1,
    }

    listing, err := client.CreateListing(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("failed to create listing: %w", err)
    }

    log.Printf("Created listing ID: %d, UUID: %s", listing.ID, listing.UUID)
    return listing, nil
}
```

### Update Listing

```go
func updateListing(client *service.Client, id int64) error {
    ctx := context.Background()

    newTitle := "Vintage Camera - Price Reduced!"
    newPrice := 249.99
    status := service.StatusActive

    req := &service.UpdateListingRequest{
        Title:  &newTitle,
        Price:  &newPrice,
        Status: &status,
    }

    listing, err := client.UpdateListing(ctx, id, req)
    if err != nil {
        return fmt.Errorf("failed to update listing: %w", err)
    }

    log.Printf("Updated listing: %s (New price: %.2f)", listing.Title, listing.Price)
    return nil
}
```

### Delete Listing

```go
func deleteListing(client *service.Client, id int64) error {
    ctx := context.Background()

    if err := client.DeleteListing(ctx, id); err != nil {
        return fmt.Errorf("failed to delete listing: %w", err)
    }

    log.Printf("Deleted listing ID: %d", id)
    return nil
}
```

---

## Search and Filtering

### Full-Text Search

```go
func searchListings(client *service.Client, query string) error {
    ctx := context.Background()

    req := &service.SearchListingsRequest{
        Query:  query,
        Limit:  20,
        Offset: 0,
    }

    resp, err := client.SearchListings(ctx, req)
    if err != nil {
        return fmt.Errorf("search failed: %w", err)
    }

    log.Printf("Found %d listings matching '%s'", resp.Total, query)
    for _, listing := range resp.Listings {
        log.Printf("- %s (%.2f %s)", listing.Title, listing.Price, listing.Currency)
    }

    return nil
}
```

### Advanced Search with Filters

```go
func advancedSearch(client *service.Client) error {
    ctx := context.Background()

    categoryID := int64(1001)
    minPrice := 100.0
    maxPrice := 500.0

    req := &service.SearchListingsRequest{
        Query:      "camera",
        CategoryID: &categoryID,
        MinPrice:   &minPrice,
        MaxPrice:   &maxPrice,
        Limit:      10,
        Offset:     0,
    }

    resp, err := client.SearchListings(ctx, req)
    if err != nil {
        return fmt.Errorf("search failed: %w", err)
    }

    log.Printf("Found %d cameras in price range $%.2f-$%.2f",
        resp.Total, minPrice, maxPrice)

    return nil
}
```

### List User's Listings

```go
func listUserListings(client *service.Client, userID int64) error {
    ctx := context.Background()

    status := service.StatusActive
    req := &service.ListListingsRequest{
        UserID: &userID,
        Status: &status,
        Limit:  50,
        Offset: 0,
    }

    resp, err := client.ListListings(ctx, req)
    if err != nil {
        return fmt.Errorf("failed to list listings: %w", err)
    }

    log.Printf("User %d has %d active listings", userID, resp.Total)
    return nil
}
```

### Pagination Example

```go
func paginateListings(client *service.Client) error {
    ctx := context.Background()

    const pageSize = 20
    page := 0

    for {
        req := &service.ListListingsRequest{
            Limit:  pageSize,
            Offset: int32(page * pageSize),
        }

        resp, err := client.ListListings(ctx, req)
        if err != nil {
            return err
        }

        if len(resp.Listings) == 0 {
            break  // No more results
        }

        log.Printf("Page %d: %d listings", page+1, len(resp.Listings))

        // Process listings...
        for _, listing := range resp.Listings {
            log.Printf("  - %s", listing.Title)
        }

        page++
    }

    return nil
}
```

---

## Fiber Middleware Integration

### Complete Handler Setup

```go
package handlers

import (
    "github.com/gofiber/fiber/v2"
    authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"
    listingsMiddleware "github.com/sveturs/listings/pkg/http/fiber/middleware"
    "github.com/sveturs/listings/pkg/service"
)

type Handler struct {
    listingsClient *service.Client
    logger         zerolog.Logger
}

func NewHandler(client *service.Client, logger zerolog.Logger) *Handler {
    return &Handler{
        listingsClient: client,
        logger:         logger,
    }
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
    // Inject client into all routes
    app.Use("/listings", listingsMiddleware.InjectListingsClient(h.listingsClient))

    // Log all operations
    app.Use("/listings", listingsMiddleware.LogListingOperations(h.logger))

    // Public routes
    listings := app.Group("/listings")
    listings.Get("/:id", h.GetListing)
    listings.Get("/search", h.SearchListings)

    // Protected routes (authentication required)
    protected := listings.Group("")
    protected.Use(authMiddleware.RequireAuth())
    protected.Post("/", h.CreateListing)

    // Owner-only routes (must own the listing)
    protected.Put("/:id",
        listingsMiddleware.RequireListingOwnership(h.listingsClient),
        h.UpdateListing,
    )
    protected.Delete("/:id",
        listingsMiddleware.RequireListingOwnership(h.listingsClient),
        h.DeleteListing,
    )
}
```

### Handler Implementation

```go
// GetListing retrieves a single listing
func (h *Handler) GetListing(c *fiber.Ctx) error {
    id, err := c.ParamsInt("id")
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "invalid_id",
        })
    }

    client := listingsMiddleware.GetListingsClient(c)
    listing, err := client.GetListing(c.Context(), int64(id))
    if err != nil {
        if err == service.ErrNotFound {
            return c.Status(404).JSON(fiber.Map{
                "error": "listing_not_found",
            })
        }
        return c.Status(500).JSON(fiber.Map{
            "error": "internal_error",
        })
    }

    return c.JSON(fiber.Map{
        "data": listing,
    })
}

// UpdateListing updates an existing listing (owner only)
func (h *Handler) UpdateListing(c *fiber.Ctx) error {
    // Listing is already verified by RequireListingOwnership middleware
    listing := listingsMiddleware.GetListing(c)

    var req service.UpdateListingRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "invalid_request_body",
        })
    }

    client := listingsMiddleware.GetListingsClient(c)
    updated, err := client.UpdateListing(c.Context(), listing.ID, &req)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "update_failed",
        })
    }

    return c.JSON(fiber.Map{
        "data": updated,
    })
}
```

### Rate Limiting Example

```go
func setupRateLimiting(app *fiber.App) {
    // Rate limit by user for write operations
    app.Post("/listings",
        authMiddleware.RequireAuth(),
        listingsMiddleware.RateLimitByUserID(10, time.Minute),  // 10 req/min
        handlers.CreateListing,
    )
}
```

### Caching Example

```go
func setupCaching(app *fiber.App, logger zerolog.Logger) {
    cacheConfig := listingsMiddleware.CacheListingsConfig{
        TTL:    5 * time.Minute,
        Logger: logger,
        KeyGenerator: func(c *fiber.Ctx) string {
            return fmt.Sprintf("listing:%s", c.Params("id"))
        },
    }

    app.Get("/listings/:id",
        listingsMiddleware.CacheListings(cacheConfig),
        handlers.GetListing,
    )
}
```

---

## gRPC Connection Pooling

### Basic Pool Setup

```go
package main

import (
    "log"

    grpcpkg "github.com/sveturs/listings/pkg/grpc"
    pb "github.com/sveturs/listings/api/proto/listings/v1"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func main() {
    pool, err := grpcpkg.NewPool(grpcpkg.PoolConfig{
        Size:   5,
        Target: "localhost:50053",
        DialOptions: []grpc.DialOption{
            grpc.WithTransportCredentials(insecure.NewCredentials()),
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    defer pool.Close()

    // Get connection (round-robin)
    conn := pool.Get()
    client := pb.NewListingsServiceClient(conn)

    // Use client...
}
```

### Pool with Interceptors

```go
func setupPoolWithInterceptors(logger zerolog.Logger, token string) (*grpcpkg.Pool, error) {
    interceptor := grpcpkg.ChainUnaryClient(
        grpcpkg.LoggingInterceptor(logger),
        grpcpkg.MetricsInterceptor(),
        grpcpkg.AuthInterceptor(token),
        grpcpkg.RetryInterceptor(3, 100*time.Millisecond, logger),
        grpcpkg.TimeoutInterceptor(5*time.Second),
    )

    return grpcpkg.NewPool(grpcpkg.PoolConfig{
        Size:   10,
        Target: "listings.internal:50053",
        DialOptions: []grpc.DialOption{
            grpc.WithTransportCredentials(insecure.NewCredentials()),
            grpc.WithUnaryInterceptor(interceptor),
        },
    })
}
```

### Pool Health Monitoring

```go
func monitorPoolHealth(pool *grpcpkg.Pool, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for range ticker.C {
        stats := pool.GetStats()
        log.Printf("Pool stats: Size=%d, Healthy=%d, Requests=%d, Closed=%v",
            stats.Size,
            stats.HealthyConns,
            stats.RequestCounter,
            stats.Closed,
        )

        if stats.HealthyConns < stats.Size/2 {
            log.Warn().Msg("More than half of connections are unhealthy")
        }
    }
}
```

---

## Error Handling

### Typed Error Handling

```go
func handleTypedErrors(client *service.Client, id int64) {
    listing, err := client.GetListing(context.Background(), id)
    if err != nil {
        switch err {
        case service.ErrNotFound:
            log.Printf("Listing %d does not exist", id)
            // Maybe create it?

        case service.ErrInvalidInput:
            log.Printf("Invalid input: %v", err)
            // Fix validation

        case service.ErrUnavailable:
            log.Printf("Service temporarily unavailable")
            // Retry later

        default:
            log.Printf("Unexpected error: %v", err)
            // Log and alert
        }
        return
    }

    // Use listing...
}
```

### Retry with Exponential Backoff

```go
func retryWithBackoff(client *service.Client, id int64) (*service.Listing, error) {
    var listing *service.Listing
    var err error

    backoff := 100 * time.Millisecond
    maxRetries := 5

    for i := 0; i < maxRetries; i++ {
        listing, err = client.GetListing(context.Background(), id)
        if err == nil {
            return listing, nil
        }

        if err != service.ErrUnavailable {
            return nil, err  // Don't retry non-transient errors
        }

        log.Printf("Attempt %d failed, retrying in %v", i+1, backoff)
        time.Sleep(backoff)
        backoff *= 2  // Exponential backoff
    }

    return nil, fmt.Errorf("max retries exceeded: %w", err)
}
```

---

## Advanced Patterns

### Batch Operations

```go
func batchUpdatePrices(client *service.Client, updates map[int64]float64) error {
    ctx := context.Background()

    for id, newPrice := range updates {
        req := &service.UpdateListingRequest{
            Price: &newPrice,
        }

        if _, err := client.UpdateListing(ctx, id, req); err != nil {
            log.Printf("Failed to update listing %d: %v", id, err)
            continue  // Continue with other updates
        }

        log.Printf("Updated listing %d to %.2f", id, newPrice)
    }

    return nil
}
```

### Concurrent Requests with Worker Pool

```go
func fetchListingsConcurrently(client *service.Client, ids []int64, workers int) []*service.Listing {
    jobs := make(chan int64, len(ids))
    results := make(chan *service.Listing, len(ids))

    // Start workers
    for w := 0; w < workers; w++ {
        go func() {
            for id := range jobs {
                listing, err := client.GetListing(context.Background(), id)
                if err != nil {
                    log.Printf("Failed to fetch listing %d: %v", id, err)
                    continue
                }
                results <- listing
            }
        }()
    }

    // Send jobs
    for _, id := range ids {
        jobs <- id
    }
    close(jobs)

    // Collect results
    var listings []*service.Listing
    for i := 0; i < len(ids); i++ {
        select {
        case listing := <-results:
            listings = append(listings, listing)
        case <-time.After(10 * time.Second):
            log.Println("Timeout waiting for results")
            break
        }
    }

    return listings
}
```

### Circuit Breaker Pattern

```go
type CircuitBreaker struct {
    client       *service.Client
    failures     int
    maxFailures  int
    state        string  // "closed", "open", "half-open"
    lastFailTime time.Time
    timeout      time.Duration
    mu           sync.RWMutex
}

func (cb *CircuitBreaker) GetListing(ctx context.Context, id int64) (*service.Listing, error) {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    if cb.state == "open" {
        if time.Since(cb.lastFailTime) > cb.timeout {
            cb.state = "half-open"
        } else {
            return nil, fmt.Errorf("circuit breaker is open")
        }
    }

    listing, err := cb.client.GetListing(ctx, id)
    if err != nil {
        cb.failures++
        cb.lastFailTime = time.Now()

        if cb.failures >= cb.maxFailures {
            cb.state = "open"
            log.Println("Circuit breaker opened")
        }

        return nil, err
    }

    // Success - reset
    cb.failures = 0
    cb.state = "closed"

    return listing, nil
}
```

---

## Complete Integration Example

```go
package main

import (
    "context"
    "log"
    "os"
    "time"

    "github.com/gofiber/fiber/v2"
    "github.com/rs/zerolog"
    authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"
    grpcpkg "github.com/sveturs/listings/pkg/grpc"
    listingsMiddleware "github.com/sveturs/listings/pkg/http/fiber/middleware"
    "github.com/sveturs/listings/pkg/service"
)

func main() {
    logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

    // Setup listings client
    listingsClient, err := service.NewClient(service.ClientConfig{
        GRPCAddr:       os.Getenv("LISTINGS_GRPC_ADDR"),
        HTTPBaseURL:    os.Getenv("LISTINGS_HTTP_URL"),
        AuthToken:      os.Getenv("SERVICE_TOKEN"),
        Timeout:        10 * time.Second,
        EnableFallback: true,
        Logger:         logger,
    })
    if err != nil {
        log.Fatal("Failed to create listings client:", err)
    }
    defer listingsClient.Close()

    // Setup Fiber app
    app := fiber.New()

    // Global middleware
    app.Use(listingsMiddleware.InjectListingsClient(listingsClient))
    app.Use(listingsMiddleware.LogListingOperations(logger))

    // Routes
    setupRoutes(app, listingsClient)

    // Start server
    log.Fatal(app.Listen(":8080"))
}

func setupRoutes(app *fiber.App, client *service.Client) {
    // Public routes
    app.Get("/listings/:id", getListingHandler)
    app.Get("/listings/search", searchListingsHandler)

    // Protected routes
    protected := app.Group("/listings", authMiddleware.RequireAuth())
    protected.Post("/", createListingHandler)
    protected.Put("/:id",
        listingsMiddleware.RequireListingOwnership(client),
        updateListingHandler,
    )
}
```

---

For more information, see [README.md](./README.md).
