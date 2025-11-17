package cache

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

// SearchCache provides caching for search results
type SearchCache struct {
	client *redis.Client
	ttl    time.Duration
	logger zerolog.Logger
}

// SearchRequest represents search parameters for cache key generation
type SearchRequest struct {
	Query      string
	CategoryID *int64
	Limit      int32
	Offset     int32
}

// SearchResult represents cached search result
type SearchResult struct {
	Listings   []map[string]interface{} `json:"listings"`
	Total      int64                    `json:"total"`
	TookMs     int32                    `json:"took_ms"`
	CachedAt   time.Time                `json:"cached_at"`
	CacheKeyID string                   `json:"cache_key_id"` // For debugging
}

// NewSearchCache creates a new search cache client
func NewSearchCache(redisURL string, ttl time.Duration, logger zerolog.Logger) (*SearchCache, error) {
	if ttl <= 0 {
		ttl = 5 * time.Minute // Default TTL
	}

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opts)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info().
		Dur("ttl", ttl).
		Msg("Search cache initialized")

	return &SearchCache{
		client: client,
		ttl:    ttl,
		logger: logger.With().Str("component", "search_cache").Logger(),
	}, nil
}

// Get retrieves cached search result
func (sc *SearchCache) Get(ctx context.Context, req *SearchRequest) (*SearchResult, error) {
	key := sc.generateKey(req)

	data, err := sc.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			sc.logger.Debug().
				Str("key", key).
				Msg("cache miss")
			return nil, fmt.Errorf("cache miss")
		}
		sc.logger.Error().
			Err(err).
			Str("key", key).
			Msg("failed to get from cache")
		return nil, fmt.Errorf("cache get failed: %w", err)
	}

	var result SearchResult
	if err := json.Unmarshal(data, &result); err != nil {
		sc.logger.Error().
			Err(err).
			Str("key", key).
			Msg("failed to unmarshal cached data")
		return nil, fmt.Errorf("cache unmarshal failed: %w", err)
	}

	sc.logger.Debug().
		Str("key", key).
		Int64("total", result.Total).
		Time("cached_at", result.CachedAt).
		Msg("cache hit")

	return &result, nil
}

// Set stores search result in cache
func (sc *SearchCache) Set(ctx context.Context, req *SearchRequest, result *SearchResult) error {
	key := sc.generateKey(req)

	// Add cache metadata
	result.CachedAt = time.Now()
	result.CacheKeyID = key

	data, err := json.Marshal(result)
	if err != nil {
		sc.logger.Error().
			Err(err).
			Str("key", key).
			Msg("failed to marshal search result")
		return fmt.Errorf("cache marshal failed: %w", err)
	}

	if err := sc.client.Set(ctx, key, data, sc.ttl).Err(); err != nil {
		sc.logger.Error().
			Err(err).
			Str("key", key).
			Msg("failed to set cache")
		return fmt.Errorf("cache set failed: %w", err)
	}

	sc.logger.Debug().
		Str("key", key).
		Dur("ttl", sc.ttl).
		Int64("total", result.Total).
		Msg("cache set")

	return nil
}

// generateKey creates a unique cache key from search parameters
func (sc *SearchCache) generateKey(req *SearchRequest) string {
	// Create deterministic string representation
	keyParts := fmt.Sprintf("q:%s|cat:%v|lim:%d|off:%d",
		req.Query,
		req.CategoryID,
		req.Limit,
		req.Offset,
	)

	// Generate MD5 hash for compact key
	hash := md5.Sum([]byte(keyParts))
	hashStr := hex.EncodeToString(hash[:])

	// Use prefix for easy identification and cleanup
	return fmt.Sprintf("search:v1:%s", hashStr)
}

// Invalidate removes a specific cached search result
func (sc *SearchCache) Invalidate(ctx context.Context, req *SearchRequest) error {
	key := sc.generateKey(req)

	if err := sc.client.Del(ctx, key).Err(); err != nil {
		sc.logger.Error().
			Err(err).
			Str("key", key).
			Msg("failed to invalidate cache")
		return fmt.Errorf("cache invalidate failed: %w", err)
	}

	sc.logger.Debug().
		Str("key", key).
		Msg("cache invalidated")

	return nil
}

// InvalidateAll removes all search cache entries (use with caution!)
func (sc *SearchCache) InvalidateAll(ctx context.Context) error {
	pattern := "search:v1:*"

	iter := sc.client.Scan(ctx, 0, pattern, 0).Iterator()
	pipe := sc.client.Pipeline()

	count := 0
	for iter.Next(ctx) {
		pipe.Del(ctx, iter.Val())
		count++
	}

	if err := iter.Err(); err != nil {
		sc.logger.Error().
			Err(err).
			Str("pattern", pattern).
			Msg("failed to scan keys")
		return fmt.Errorf("cache scan failed: %w", err)
	}

	if count > 0 {
		if _, err := pipe.Exec(ctx); err != nil {
			sc.logger.Error().
				Err(err).
				Str("pattern", pattern).
				Int("count", count).
				Msg("failed to delete keys")
			return fmt.Errorf("cache delete failed: %w", err)
		}

		sc.logger.Info().
			Str("pattern", pattern).
			Int("count", count).
			Msg("invalidated all search cache")
	}

	return nil
}

// GetStats returns cache statistics
func (sc *SearchCache) GetStats(ctx context.Context) (map[string]interface{}, error) {
	pattern := "search:v1:*"

	var count int64
	iter := sc.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		count++
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("failed to count keys: %w", err)
	}

	return map[string]interface{}{
		"total_keys": count,
		"pattern":    pattern,
		"ttl":        sc.ttl.String(),
	}, nil
}

// HealthCheck performs a health check on Redis connection
func (sc *SearchCache) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := sc.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis health check failed: %w", err)
	}

	return nil
}

// Close closes the Redis client connection
func (sc *SearchCache) Close() error {
	return sc.client.Close()
}
