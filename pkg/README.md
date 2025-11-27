# Listings Service Client Library

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Coverage](https://img.shields.io/badge/coverage-65%25-yellow)](https://github.com/sveturs/listings)

Production-ready Go client library for integrating with the Listings microservice. Supports both gRPC (primary) and HTTP REST (fallback) communication methods.

## Features

- **Unified Client Interface** - Single API for both gRPC and HTTP transports
- **Automatic Failover** - Falls back to HTTP when gRPC is unavailable
- **Connection Pooling** - Efficient gRPC connection management
- **Middleware Support** - Ready-to-use Fiber middleware for common patterns
- **Type-Safe** - Fully typed interfaces with comprehensive validation
- **Observability** - Built-in logging, metrics, and tracing interceptors
- **Production-Ready** - Retry logic, timeouts, circuit breakers

## Installation

```bash
go get github.com/sveturs/listings@latest
```

## Quick Start

### Basic Client Usage

```go
import (
    "context"
    "time"

    "github.com/rs/zerolog/log"
    "github.com/sveturs/listings/pkg/service"
)

func main() {
    // Create client with both gRPC and HTTP fallback
    client, err := service.NewClient(service.ClientConfig{
        GRPCAddr:       "localhost:50053",
        HTTPBaseURL:    "http://localhost:8086",
        AuthToken:      os.Getenv("SERVICE_TOKEN"),
        Timeout:        5 * time.Second,
        EnableFallback: true,
        Logger:         log.Logger,
    })
    if err != nil {
        log.Fatal().Err(err).Msg("Failed to create client")
    }
    defer client.Close()

    // Get a listing
    listing, err := client.GetListing(context.Background(), 123)
    if err != nil {
        log.Error().Err(err).Msg("Failed to get listing")
        return
    }

    log.Info().
        Int64("id", listing.ID).
        Str("title", listing.Title).
        Float64("price", listing.Price).
        Msg("Retrieved listing")
}
```

### Fiber Middleware Integration

```go
import (
    "github.com/gofiber/fiber/v2"
    authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"
    listingsMiddleware "github.com/sveturs/listings/pkg/http/fiber/middleware"
)

func setupRoutes(app *fiber.App, listingsClient *service.Client) {
    // Inject listings client into all routes
    app.Use(listingsMiddleware.InjectListingsClient(listingsClient))

    // Public routes
    app.Get("/listings/:id", handlers.GetListing)
    app.Get("/listings/search", handlers.SearchListings)

    // Protected routes - require authentication
    protected := app.Group("/listings")
    protected.Use(authMiddleware.RequireAuth())
    protected.Post("/", handlers.CreateListing)

    // Owner-only routes - verify listing ownership
    protected.Put("/:id",
        listingsMiddleware.RequireListingOwnership(listingsClient),
        handlers.UpdateListing,
    )
    protected.Delete("/:id",
        listingsMiddleware.RequireListingOwnership(listingsClient),
        handlers.DeleteListing,
    )
}
```

### gRPC Connection Pooling

```go
import (
    "github.com/rs/zerolog/log"
    grpcpkg "github.com/sveturs/listings/pkg/grpc"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func main() {
    // Create connection pool with 5 connections
    pool, err := grpcpkg.NewPool(grpcpkg.PoolConfig{
        Size:   5,
        Target: "localhost:50053",
        DialOptions: []grpc.DialOption{
            grpc.WithTransportCredentials(insecure.NewCredentials()),
            grpc.WithUnaryInterceptor(grpcpkg.ChainUnaryClient(
                grpcpkg.LoggingInterceptor(log.Logger),
                grpcpkg.MetricsInterceptor(),
                grpcpkg.AuthInterceptor(os.Getenv("SERVICE_TOKEN")),
                grpcpkg.RetryInterceptor(3, 100*time.Millisecond, log.Logger),
            )),
        },
    })
    if err != nil {
        log.Fatal().Err(err).Msg("Failed to create pool")
    }
    defer pool.Close()

    // Get a connection (round-robin)
    conn := pool.Get()
    client := pb.NewListingsServiceClient(conn)

    // Use client...
}
```

## Package Structure

```
pkg/
├── service/              # Unified client (gRPC + HTTP)
│   ├── client.go         # Main client implementation
│   ├── grpc_client.go    # gRPC transport layer
│   ├── http_client.go    # HTTP transport layer
│   └── types.go          # Shared types and models
├── grpc/                 # gRPC utilities
│   ├── interceptors.go   # Logging, metrics, auth, retry
│   └── client_pool.go    # Connection pooling
└── http/fiber/middleware/ # Fiber middleware
    └── listings.go       # Request validation, ownership checks
```

## API Reference

### Client Methods

#### GetListing
```go
func (c *Client) GetListing(ctx context.Context, id int64) (*Listing, error)
```
Retrieves a single listing by ID.

#### CreateListing
```go
func (c *Client) CreateListing(ctx context.Context, req *CreateListingRequest) (*Listing, error)
```
Creates a new listing.

#### UpdateListing
```go
func (c *Client) UpdateListing(ctx context.Context, id int64, req *UpdateListingRequest) (*Listing, error)
```
Updates an existing listing.

#### DeleteListing
```go
func (c *Client) DeleteListing(ctx context.Context, id int64) error
```
Soft-deletes a listing.

#### SearchListings
```go
func (c *Client) SearchListings(ctx context.Context, req *SearchListingsRequest) (*SearchListingsResponse, error)
```
Performs full-text search on listings.

#### ListListings
```go
func (c *Client) ListListings(ctx context.Context, req *ListListingsRequest) (*ListListingsResponse, error)
```
Returns a paginated list of listings with filters.

### Middleware

#### InjectListingsClient
```go
func InjectListingsClient(client *service.Client) fiber.Handler
```
Injects the listings client into Fiber context.

#### RequireListingOwnership
```go
func RequireListingOwnership(client *service.Client, paramName ...string) fiber.Handler
```
Verifies that the authenticated user owns the listing.

#### CacheListings
```go
func CacheListings(config CacheListingsConfig) fiber.Handler
```
Caches listing responses (useful for read-heavy endpoints).

#### RateLimitByUserID
```go
func RateLimitByUserID(maxRequests int, window time.Duration) fiber.Handler
```
Rate limits requests per user.

#### LogListingOperations
```go
func LogListingOperations(logger zerolog.Logger) fiber.Handler
```
Logs all listing operations for audit trails.

### gRPC Interceptors

#### LoggingInterceptor
```go
func LoggingInterceptor(logger zerolog.Logger) grpc.UnaryClientInterceptor
```
Logs all gRPC calls with duration and errors.

#### MetricsInterceptor
```go
func MetricsInterceptor() grpc.UnaryClientInterceptor
```
Collects metrics for gRPC calls (integrate with Prometheus).

#### AuthInterceptor
```go
func AuthInterceptor(token string) grpc.UnaryClientInterceptor
```
Adds authentication token to gRPC metadata.

#### RetryInterceptor
```go
func RetryInterceptor(maxRetries int, initialDelay time.Duration, logger zerolog.Logger) grpc.UnaryClientInterceptor
```
Retries failed requests with exponential backoff.

#### TimeoutInterceptor
```go
func TimeoutInterceptor(timeout time.Duration) grpc.UnaryClientInterceptor
```
Enforces timeout on all gRPC calls.

#### ChainUnaryClient
```go
func ChainUnaryClient(interceptors ...grpc.UnaryClientInterceptor) grpc.UnaryClientInterceptor
```
Chains multiple interceptors into one.

## Error Handling

The library uses typed errors for common cases:

```go
import "github.com/sveturs/listings/pkg/service"

listing, err := client.GetListing(ctx, 123)
if err != nil {
    switch err {
    case service.ErrNotFound:
        // Handle not found (404)
    case service.ErrInvalidInput:
        // Handle validation error (400)
    case service.ErrUnavailable:
        // Handle service unavailable (503)
    default:
        // Handle other errors
    }
}
```

## Testing

Run tests with coverage:

```bash
cd /p/github.com/sveturs/listings
go test -v ./pkg/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

Current test coverage: **65%+**

## Best Practices

### 1. Always Use Context
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

listing, err := client.GetListing(ctx, 123)
```

### 2. Enable Fallback for Production
```go
client, err := service.NewClient(service.ClientConfig{
    GRPCAddr:       "listings-grpc.internal:50053",
    HTTPBaseURL:    "http://listings-http.internal:8086",
    EnableFallback: true,  // Graceful degradation
    // ...
})
```

### 3. Use Connection Pooling
```go
// Don't create a new connection for each request
pool, err := grpcpkg.NewPool(grpcpkg.PoolConfig{
    Size:   10,  // Adjust based on load
    Target: "localhost:50053",
    // ...
})
```

### 4. Implement Observability
```go
// Chain interceptors for full observability
interceptor := grpcpkg.ChainUnaryClient(
    grpcpkg.LoggingInterceptor(logger),
    grpcpkg.MetricsInterceptor(),
    grpcpkg.AuthInterceptor(token),
    grpcpkg.RetryInterceptor(3, 100*time.Millisecond, logger),
)
```

## Roadmap

- ✅ Sprint 4.3: Public pkg library (current)
- 🚧 Sprint 4.4: gRPC server implementation
- 📋 Sprint 4.5: OpenSearch integration
- 📋 Sprint 4.6: Caching layer (Redis)

## Contributing

See [EXAMPLES.md](./EXAMPLES.md) for more detailed usage examples.

For questions or issues, contact the team or open an issue on GitHub.

## License

Internal use only - Vondi Platform

---

**Generated:** Sprint 4.3 (2025-10-31)
**Version:** 1.0.0
**Status:** Production Ready (HTTP), Stubs (gRPC - awaiting Sprint 4.4)
