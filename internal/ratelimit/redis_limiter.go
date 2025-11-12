package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

// RedisLimiter implements rate limiting using Redis
type RedisLimiter struct {
	client *redis.Client
	logger zerolog.Logger
}

// NewRedisLimiter creates a new Redis-backed rate limiter
func NewRedisLimiter(client *redis.Client, logger zerolog.Logger) *RedisLimiter {
	return &RedisLimiter{
		client: client,
		logger: logger.With().Str("component", "redis_limiter").Logger(),
	}
}

// Allow checks if a request is allowed under the rate limit using token bucket algorithm
// Key pattern: rate_limit:{method}:{identifier}
func (r *RedisLimiter) Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	rateLimitKey := fmt.Sprintf("rate_limit:%s", key)

	// Use Lua script for atomic operations to avoid race conditions
	// This implements a token bucket algorithm
	script := redis.NewScript(`
		local key = KEYS[1]
		local limit = tonumber(ARGV[1])
		local window = tonumber(ARGV[2])
		local current_time = tonumber(ARGV[3])

		-- Get current value
		local current = redis.call('GET', key)

		if current == false then
			-- First request in this window
			redis.call('SET', key, limit - 1, 'EX', window)
			return {1, limit - 1}  -- allowed, remaining
		end

		current = tonumber(current)

		if current > 0 then
			-- Decrement the counter
			local remaining = redis.call('DECR', key)
			return {1, remaining}  -- allowed, remaining
		else
			-- Rate limit exceeded
			local ttl = redis.call('TTL', key)
			return {0, 0, ttl}  -- not allowed, no remaining, time until reset
		end
	`)

	// Execute the Lua script
	result, err := script.Run(
		ctx,
		r.client,
		[]string{rateLimitKey},
		limit,
		int(window.Seconds()),
		time.Now().Unix(),
	).Result()

	if err != nil {
		// If Redis fails, fail open (allow the request) to avoid cascading failures
		r.logger.Error().
			Err(err).
			Str("key", key).
			Msg("redis error, failing open (allowing request)")
		return true, nil
	}

	// Parse result
	resultSlice, ok := result.([]interface{})
	if !ok || len(resultSlice) < 2 {
		r.logger.Error().
			Str("key", key).
			Msg("unexpected result format from lua script")
		return true, nil // Fail open
	}

	allowed := resultSlice[0].(int64) == 1
	remaining := resultSlice[1].(int64)

	if !allowed {
		ttl := int64(-1)
		if len(resultSlice) > 2 {
			ttl = resultSlice[2].(int64)
		}
		r.logger.Debug().
			Str("key", key).
			Int("limit", limit).
			Int64("ttl", ttl).
			Msg("rate limit exceeded")
	} else {
		r.logger.Debug().
			Str("key", key).
			Int64("remaining", remaining).
			Msg("request allowed")
	}

	return allowed, nil
}

// Remaining returns the number of remaining requests in the current window
func (r *RedisLimiter) Remaining(ctx context.Context, key string) (int, error) {
	rateLimitKey := fmt.Sprintf("rate_limit:%s", key)

	val, err := r.client.Get(ctx, rateLimitKey).Int()
	if err != nil {
		if err == redis.Nil {
			// Key doesn't exist, all requests available
			return -1, nil
		}
		return 0, fmt.Errorf("failed to get remaining count: %w", err)
	}

	return val, nil
}

// Reset clears the rate limit for a specific key
func (r *RedisLimiter) Reset(ctx context.Context, key string) error {
	rateLimitKey := fmt.Sprintf("rate_limit:%s", key)

	err := r.client.Del(ctx, rateLimitKey).Err()
	if err != nil {
		return fmt.Errorf("failed to reset rate limit: %w", err)
	}

	r.logger.Debug().Str("key", key).Msg("rate limit reset")
	return nil
}

// HealthCheck verifies Redis connectivity
func (r *RedisLimiter) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := r.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis health check failed: %w", err)
	}

	return nil
}
