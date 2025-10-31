package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

// RedisCache implements caching using Redis
type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
	logger zerolog.Logger
}

// NewRedisCache creates a new Redis cache client
func NewRedisCache(addr, password string, db int, poolSize, minIdleConns int, ttl time.Duration, logger zerolog.Logger) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		PoolSize:     poolSize,
		MinIdleConns: minIdleConns,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info().
		Str("addr", addr).
		Int("pool_size", poolSize).
		Int("min_idle_conns", minIdleConns).
		Dur("default_ttl", ttl).
		Msg("Redis cache initialized")

	return &RedisCache{
		client: client,
		ttl:    ttl,
		logger: logger.With().Str("component", "redis_cache").Logger(),
	}, nil
}

// Get retrieves a value from cache
func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("cache miss: key not found")
		}
		c.logger.Error().Err(err).Str("key", key).Msg("failed to get from cache")
		return fmt.Errorf("cache get failed: %w", err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		c.logger.Error().Err(err).Str("key", key).Msg("failed to unmarshal cached data")
		return fmt.Errorf("cache unmarshal failed: %w", err)
	}

	c.logger.Debug().Str("key", key).Msg("cache hit")
	return nil
}

// Set stores a value in cache with default TTL
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}) error {
	return c.SetWithTTL(ctx, key, value, c.ttl)
}

// SetWithTTL stores a value in cache with custom TTL
func (c *RedisCache) SetWithTTL(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		c.logger.Error().Err(err).Str("key", key).Msg("failed to marshal value")
		return fmt.Errorf("cache marshal failed: %w", err)
	}

	if err := c.client.Set(ctx, key, data, ttl).Err(); err != nil {
		c.logger.Error().Err(err).Str("key", key).Msg("failed to set cache")
		return fmt.Errorf("cache set failed: %w", err)
	}

	c.logger.Debug().Str("key", key).Dur("ttl", ttl).Msg("cache set")
	return nil
}

// Delete removes a key from cache
func (c *RedisCache) Delete(ctx context.Context, key string) error {
	if err := c.client.Del(ctx, key).Err(); err != nil {
		c.logger.Error().Err(err).Str("key", key).Msg("failed to delete from cache")
		return fmt.Errorf("cache delete failed: %w", err)
	}

	c.logger.Debug().Str("key", key).Msg("cache deleted")
	return nil
}

// DeletePattern deletes all keys matching a pattern
func (c *RedisCache) DeletePattern(ctx context.Context, pattern string) error {
	iter := c.client.Scan(ctx, 0, pattern, 0).Iterator()
	pipe := c.client.Pipeline()

	count := 0
	for iter.Next(ctx) {
		pipe.Del(ctx, iter.Val())
		count++
	}

	if err := iter.Err(); err != nil {
		c.logger.Error().Err(err).Str("pattern", pattern).Msg("failed to scan keys")
		return fmt.Errorf("cache scan failed: %w", err)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		c.logger.Error().Err(err).Str("pattern", pattern).Msg("failed to delete keys")
		return fmt.Errorf("cache delete pattern failed: %w", err)
	}

	c.logger.Debug().Str("pattern", pattern).Int("count", count).Msg("deleted keys by pattern")
	return nil
}

// Exists checks if a key exists in cache
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		c.logger.Error().Err(err).Str("key", key).Msg("failed to check key existence")
		return false, fmt.Errorf("cache exists check failed: %w", err)
	}

	return n > 0, nil
}

// Increment increments a counter key
func (c *RedisCache) Increment(ctx context.Context, key string) (int64, error) {
	val, err := c.client.Incr(ctx, key).Result()
	if err != nil {
		c.logger.Error().Err(err).Str("key", key).Msg("failed to increment counter")
		return 0, fmt.Errorf("cache increment failed: %w", err)
	}

	return val, nil
}

// HealthCheck performs a health check on Redis connection
func (c *RedisCache) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := c.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis health check failed: %w", err)
	}

	return nil
}

// GetPoolStats returns connection pool statistics
func (c *RedisCache) GetPoolStats() *redis.PoolStats {
	return c.client.PoolStats()
}

// Close closes the Redis client connection
func (c *RedisCache) Close() error {
	return c.client.Close()
}

// FlushAll flushes all keys from Redis (use with caution!)
func (c *RedisCache) FlushAll(ctx context.Context) error {
	if err := c.client.FlushAll(ctx).Err(); err != nil {
		c.logger.Error().Err(err).Msg("failed to flush all keys")
		return fmt.Errorf("cache flush all failed: %w", err)
	}

	c.logger.Warn().Msg("flushed all cache keys")
	return nil
}
