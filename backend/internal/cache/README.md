# Redis Cache Package

This package provides a Redis-based caching implementation with graceful degradation support.

## Features

- ✅ Basic cache operations (Get, Set, Delete)
- ✅ Pattern-based deletion (DeletePattern)
- ✅ Key existence check (Exists)
- ✅ GetOrSet pattern for cache-aside loading
- ✅ TTL support for automatic expiration
- ✅ Graceful degradation when Redis is unavailable
- ✅ JSON serialization for complex data types
- ✅ Concurrent-safe operations
- ✅ Comprehensive logging with logrus

## Usage

```go
// Create a new Redis cache instance
logger := logrus.New()
cache, err := NewRedisCache("localhost:6379", "", 0, 10, logger)
if err != nil {
    log.Fatal(err)
}
defer cache.Close()

// Set a value with TTL
err = cache.Set(ctx, "user:123", user, 5*time.Minute)

// Get a value
var user User
err = cache.Get(ctx, "user:123", &user)
if err == ErrCacheMiss {
    // Key not found
}

// GetOrSet pattern - get from cache or load and cache
err = cache.GetOrSet(ctx, "user:123", &user, 5*time.Minute, func() (interface{}, error) {
    // Load from database
    return db.GetUser(123)
})

// Delete keys
err = cache.Delete(ctx, "user:123", "user:124")

// Delete by pattern
err = cache.DeletePattern(ctx, "user:*")

// Check if key exists
exists, err := cache.Exists(ctx, "user:123")
```

## Testing

The package includes comprehensive unit tests using miniredis for mocking Redis:

```bash
# Run all tests
go test ./internal/cache -v

# Run benchmarks
go test ./internal/cache -bench=. -benchtime=10s
```

## Graceful Degradation

The cache implements graceful degradation:
- When Redis is unavailable, `GetOrSet` will still return data from the loader function
- All operations log errors but don't panic
- The application can continue working without cache

## Performance

Based on benchmarks with miniredis:
- Set operation: ~5.5μs per operation
- Get operation: ~5.3μs per operation
- GetOrSet operation: ~4.4μs per operation

Real Redis performance will vary based on network latency and server load.