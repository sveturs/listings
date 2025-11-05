package ratelimit

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRedis(t *testing.T) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15, // Use a separate DB for tests
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := client.Ping(ctx).Err()
	if err != nil {
		t.Skip("Redis not available, skipping integration tests")
	}

	// Clean up test DB
	err = client.FlushDB(ctx).Err()
	require.NoError(t, err)

	return client
}

func TestRedisLimiter_Allow(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	logger := zerolog.New(zerolog.NewTestWriter(t))
	limiter := NewRedisLimiter(client, logger)

	ctx := context.Background()
	key := "test:user:123"
	limit := 5
	window := 10 * time.Second

	t.Run("allows requests under limit", func(t *testing.T) {
		// Clean up
		_ = limiter.Reset(ctx, key)

		// First 5 requests should be allowed
		for i := 0; i < limit; i++ {
			allowed, err := limiter.Allow(ctx, key, limit, window)
			require.NoError(t, err)
			assert.True(t, allowed, "request %d should be allowed", i+1)
		}
	})

	t.Run("blocks requests over limit", func(t *testing.T) {
		// Clean up
		_ = limiter.Reset(ctx, key)

		// First 5 requests allowed
		for i := 0; i < limit; i++ {
			allowed, err := limiter.Allow(ctx, key, limit, window)
			require.NoError(t, err)
			require.True(t, allowed)
		}

		// 6th request should be blocked
		allowed, err := limiter.Allow(ctx, key, limit, window)
		require.NoError(t, err)
		assert.False(t, allowed, "request over limit should be blocked")
	})

	t.Run("resets after window expires", func(t *testing.T) {
		// Use a very short window for this test
		shortWindow := 2 * time.Second
		key := "test:user:short-window"

		// Clean up
		_ = limiter.Reset(ctx, key)

		// Use up the limit
		for i := 0; i < limit; i++ {
			allowed, err := limiter.Allow(ctx, key, limit, shortWindow)
			require.NoError(t, err)
			require.True(t, allowed)
		}

		// Next request should be blocked
		allowed, err := limiter.Allow(ctx, key, limit, shortWindow)
		require.NoError(t, err)
		assert.False(t, allowed)

		// Wait for window to expire
		time.Sleep(shortWindow + 500*time.Millisecond)

		// Request should be allowed again
		allowed, err = limiter.Allow(ctx, key, limit, shortWindow)
		require.NoError(t, err)
		assert.True(t, allowed, "request should be allowed after window expires")
	})

	t.Run("separate keys are independent", func(t *testing.T) {
		key1 := "test:user:key1"
		key2 := "test:user:key2"

		// Clean up
		_ = limiter.Reset(ctx, key1)
		_ = limiter.Reset(ctx, key2)

		// Use up limit for key1
		for i := 0; i < limit; i++ {
			allowed, err := limiter.Allow(ctx, key1, limit, window)
			require.NoError(t, err)
			require.True(t, allowed)
		}

		// key1 should be blocked
		allowed, err := limiter.Allow(ctx, key1, limit, window)
		require.NoError(t, err)
		assert.False(t, allowed)

		// But key2 should still be allowed
		allowed, err = limiter.Allow(ctx, key2, limit, window)
		require.NoError(t, err)
		assert.True(t, allowed, "different key should have independent limit")
	})
}

func TestRedisLimiter_Remaining(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	logger := zerolog.New(zerolog.NewTestWriter(t))
	limiter := NewRedisLimiter(client, logger)

	ctx := context.Background()
	key := "test:user:remaining"
	limit := 10
	window := 10 * time.Second

	// Clean up
	_ = limiter.Reset(ctx, key)

	// Use 3 requests
	for i := 0; i < 3; i++ {
		_, err := limiter.Allow(ctx, key, limit, window)
		require.NoError(t, err)
	}

	// Check remaining
	remaining, err := limiter.Remaining(ctx, key)
	require.NoError(t, err)
	assert.Equal(t, 7, remaining, "should have 7 remaining after using 3")
}

func TestRedisLimiter_Reset(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	logger := zerolog.New(zerolog.NewTestWriter(t))
	limiter := NewRedisLimiter(client, logger)

	ctx := context.Background()
	key := "test:user:reset"
	limit := 5
	window := 10 * time.Second

	// Use up the limit
	for i := 0; i < limit; i++ {
		allowed, err := limiter.Allow(ctx, key, limit, window)
		require.NoError(t, err)
		require.True(t, allowed)
	}

	// Should be blocked
	allowed, err := limiter.Allow(ctx, key, limit, window)
	require.NoError(t, err)
	assert.False(t, allowed)

	// Reset the limit
	err = limiter.Reset(ctx, key)
	require.NoError(t, err)

	// Should be allowed again
	allowed, err = limiter.Allow(ctx, key, limit, window)
	require.NoError(t, err)
	assert.True(t, allowed, "should be allowed after reset")
}

func TestRedisLimiter_HealthCheck(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	logger := zerolog.New(zerolog.NewTestWriter(t))
	limiter := NewRedisLimiter(client, logger)

	ctx := context.Background()

	err := limiter.HealthCheck(ctx)
	assert.NoError(t, err, "health check should pass with valid connection")
}

func TestRedisLimiter_ConcurrentRequests(t *testing.T) {
	client := setupTestRedis(t)
	defer client.Close()

	logger := zerolog.New(zerolog.NewTestWriter(t))
	limiter := NewRedisLimiter(client, logger)

	ctx := context.Background()
	key := "test:user:concurrent"
	limit := 10
	window := 10 * time.Second

	// Clean up
	_ = limiter.Reset(ctx, key)

	// Make concurrent requests
	results := make(chan bool, 20)
	for i := 0; i < 20; i++ {
		go func() {
			allowed, err := limiter.Allow(ctx, key, limit, window)
			if err != nil {
				results <- false
				return
			}
			results <- allowed
		}()
	}

	// Collect results
	allowedCount := 0
	blockedCount := 0
	for i := 0; i < 20; i++ {
		if <-results {
			allowedCount++
		} else {
			blockedCount++
		}
	}

	// Exactly 10 should be allowed, 10 should be blocked
	assert.Equal(t, limit, allowedCount, "exactly %d requests should be allowed", limit)
	assert.Equal(t, 10, blockedCount, "exactly 10 requests should be blocked")
}
