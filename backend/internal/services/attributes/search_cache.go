package attributes

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// SearchCache provides caching for search results with attributes
type SearchCache struct {
	redis  *redis.Client
	logger *zap.Logger
	ttl    time.Duration
	prefix string
}

// SearchCacheConfig contains configuration for search cache
type SearchCacheConfig struct {
	TTL    time.Duration
	Prefix string
}

// NewSearchCache creates a new search cache service
func NewSearchCache(redis *redis.Client, logger *zap.Logger, config SearchCacheConfig) *SearchCache {
	if config.TTL == 0 {
		config.TTL = 5 * time.Minute // Default 5 minutes cache
	}
	if config.Prefix == "" {
		config.Prefix = "search_attr:"
	}

	return &SearchCache{
		redis:  redis,
		logger: logger,
		ttl:    config.TTL,
		prefix: config.Prefix,
	}
}

// SearchParams represents search parameters including attributes
type SearchParams struct {
	Query      string              `json:"query"`
	CategoryID *int                `json:"category_id,omitempty"`
	MinPrice   *float64            `json:"min_price,omitempty"`
	MaxPrice   *float64            `json:"max_price,omitempty"`
	Attributes map[string][]string `json:"attributes,omitempty"`
	Sort       string              `json:"sort"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
}

// SearchResult represents cached search results
type SearchResult struct {
	Items    []interface{} `json:"items"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	Limit    int           `json:"limit"`
	HasMore  bool          `json:"has_more"`
	CachedAt time.Time     `json:"cached_at"`
	CacheHit bool          `json:"cache_hit,omitempty"`
}

// GenerateCacheKey creates a cache key from search parameters
func (sc *SearchCache) GenerateCacheKey(params SearchParams) string {
	// Create a deterministic key from parameters
	keyData, _ := json.Marshal(params)
	return fmt.Sprintf("%s%x", sc.prefix, keyData)
}

// Get retrieves cached search results
func (sc *SearchCache) Get(ctx context.Context, params SearchParams) (*SearchResult, error) {
	key := sc.GenerateCacheKey(params)

	data, err := sc.redis.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			// Cache miss
			return nil, nil
		}
		sc.logger.Error("Failed to get from cache", zap.Error(err), zap.String("key", key))
		return nil, err
	}

	var result SearchResult
	if err := json.Unmarshal(data, &result); err != nil {
		sc.logger.Error("Failed to unmarshal cached data", zap.Error(err))
		return nil, err
	}

	result.CacheHit = true
	return &result, nil
}

// Set stores search results in cache
func (sc *SearchCache) Set(ctx context.Context, params SearchParams, result *SearchResult) error {
	key := sc.GenerateCacheKey(params)
	result.CachedAt = time.Now()

	data, err := json.Marshal(result)
	if err != nil {
		sc.logger.Error("Failed to marshal search result", zap.Error(err))
		return err
	}

	if err := sc.redis.Set(ctx, key, data, sc.ttl).Err(); err != nil {
		sc.logger.Error("Failed to set cache", zap.Error(err), zap.String("key", key))
		return err
	}

	// Update cache metrics
	sc.incrementCacheMetric(ctx, "set")

	return nil
}

// InvalidateCategory removes all cached searches for a category
func (sc *SearchCache) InvalidateCategory(ctx context.Context, categoryID int) error {
	pattern := fmt.Sprintf("%s*category_id*%d*", sc.prefix, categoryID)
	return sc.invalidatePattern(ctx, pattern)
}

// InvalidateAll removes all cached search results
func (sc *SearchCache) InvalidateAll(ctx context.Context) error {
	pattern := fmt.Sprintf("%s*", sc.prefix)
	return sc.invalidatePattern(ctx, pattern)
}

// invalidatePattern removes all keys matching pattern
func (sc *SearchCache) invalidatePattern(ctx context.Context, pattern string) error {
	var cursor uint64
	var keys []string

	for {
		var err error
		var scanKeys []string
		scanKeys, cursor, err = sc.redis.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			sc.logger.Error("Failed to scan keys", zap.Error(err), zap.String("pattern", pattern))
			return err
		}

		keys = append(keys, scanKeys...)

		if cursor == 0 {
			break
		}
	}

	if len(keys) > 0 {
		if err := sc.redis.Del(ctx, keys...).Err(); err != nil {
			sc.logger.Error("Failed to delete keys", zap.Error(err), zap.Int("count", len(keys)))
			return err
		}
		sc.logger.Info("Invalidated cache entries", zap.Int("count", len(keys)))
	}

	return nil
}

// GetStats returns cache statistics
func (sc *SearchCache) GetStats(ctx context.Context) (map[string]int64, error) {
	stats := make(map[string]int64)

	// Get hit count
	hitCount, _ := sc.redis.Get(ctx, sc.prefix+"stats:hits").Int64()
	stats["hits"] = hitCount

	// Get miss count
	missCount, _ := sc.redis.Get(ctx, sc.prefix+"stats:misses").Int64()
	stats["misses"] = missCount

	// Get set count
	setCount, _ := sc.redis.Get(ctx, sc.prefix+"stats:sets").Int64()
	stats["sets"] = setCount

	// Calculate hit rate
	total := hitCount + missCount
	if total > 0 {
		stats["hit_rate_percent"] = (hitCount * 100) / total
	}

	// Count cached entries
	var count int64
	cursor := uint64(0)
	pattern := fmt.Sprintf("%s*", sc.prefix)

	for {
		keys, newCursor, err := sc.redis.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			break
		}

		// Exclude stats keys
		for _, key := range keys {
			if !containsString(key, "stats:") {
				count++
			}
		}

		cursor = newCursor
		if cursor == 0 {
			break
		}
	}

	stats["cached_entries"] = count

	return stats, nil
}

// incrementCacheMetric increments a cache metric counter
func (sc *SearchCache) incrementCacheMetric(ctx context.Context, metric string) {
	key := fmt.Sprintf("%sstats:%ss", sc.prefix, metric)
	sc.redis.Incr(ctx, key)
	// Set expiry for stats to prevent indefinite growth
	sc.redis.Expire(ctx, key, 24*time.Hour)
}

// RecordHit records a cache hit
func (sc *SearchCache) RecordHit(ctx context.Context) {
	sc.incrementCacheMetric(ctx, "hit")
}

// RecordMiss records a cache miss
func (sc *SearchCache) RecordMiss(ctx context.Context) {
	sc.incrementCacheMetric(ctx, "miss")
}

// BatchGetSearchResults retrieves multiple search results in batch
func (sc *SearchCache) BatchGetSearchResults(ctx context.Context, paramsList []SearchParams) ([]*SearchResult, error) {
	if len(paramsList) == 0 {
		return nil, nil
	}

	// Generate keys for all parameters
	keys := make([]string, len(paramsList))
	for i, params := range paramsList {
		keys[i] = sc.GenerateCacheKey(params)
	}

	// Get all values in a single pipeline
	pipe := sc.redis.Pipeline()
	cmds := make([]*redis.StringCmd, len(keys))

	for i, key := range keys {
		cmds[i] = pipe.Get(ctx, key)
	}

	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		sc.logger.Error("Failed to execute pipeline", zap.Error(err))
		return nil, err
	}

	// Process results
	results := make([]*SearchResult, len(cmds))
	for i, cmd := range cmds {
		data, err := cmd.Bytes()
		if err != nil {
			if err == redis.Nil {
				// Cache miss for this key
				results[i] = nil
				continue
			}
			sc.logger.Error("Failed to get result from pipeline", zap.Error(err), zap.Int("index", i))
			continue
		}

		var result SearchResult
		if err := json.Unmarshal(data, &result); err != nil {
			sc.logger.Error("Failed to unmarshal cached data", zap.Error(err), zap.Int("index", i))
			results[i] = nil
			continue
		}

		result.CacheHit = true
		results[i] = &result
	}

	return results, nil
}

// WarmupCache pre-populates cache with common searches
func (sc *SearchCache) WarmupCache(ctx context.Context, commonSearches []SearchParams,
	searchFunc func(context.Context, SearchParams) (*SearchResult, error),
) error {
	for _, params := range commonSearches {
		// Check if already cached
		existing, _ := sc.Get(ctx, params)
		if existing != nil {
			continue
		}

		// Perform search and cache result
		result, err := searchFunc(ctx, params)
		if err != nil {
			sc.logger.Error("Failed to warmup cache entry", zap.Error(err))
			continue
		}

		if err := sc.Set(ctx, params, result); err != nil {
			sc.logger.Error("Failed to cache warmup result", zap.Error(err))
		}
	}

	sc.logger.Info("Cache warmup completed", zap.Int("entries", len(commonSearches)))
	return nil
}

// Helper function to check if string contains substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr
}
