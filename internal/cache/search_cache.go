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
	config SearchCacheConfig
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

// NewSearchCache creates a new search cache client with default config
func NewSearchCache(redisURL string, ttl time.Duration, logger zerolog.Logger) (*SearchCache, error) {
	config := DefaultSearchCacheConfig()
	if ttl > 0 {
		// Override SearchTTL if provided (backward compatibility)
		config.SearchTTL = ttl
	}
	return NewSearchCacheWithConfig(redisURL, config, logger)
}

// NewSearchCacheWithConfig creates a new search cache client with custom config
func NewSearchCacheWithConfig(redisURL string, config SearchCacheConfig, logger zerolog.Logger) (*SearchCache, error) {
	// Validate and fix config
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid cache config: %w", err)
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
		Dur("search_ttl", config.SearchTTL).
		Dur("facets_ttl", config.FacetsTTL).
		Dur("suggestions_ttl", config.SuggestionsTTL).
		Dur("popular_ttl", config.PopularTTL).
		Dur("filtered_search_ttl", config.FilteredSearchTTL).
		Msg("Search cache initialized")

	return &SearchCache{
		client: client,
		config: config,
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

	if err := sc.client.Set(ctx, key, data, sc.config.SearchTTL).Err(); err != nil {
		sc.logger.Error().
			Err(err).
			Str("key", key).
			Msg("failed to set cache")
		return fmt.Errorf("cache set failed: %w", err)
	}

	sc.logger.Debug().
		Str("key", key).
		Dur("ttl", sc.config.SearchTTL).
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
		"total_keys":          count,
		"pattern":             pattern,
		"search_ttl":          sc.config.SearchTTL.String(),
		"facets_ttl":          sc.config.FacetsTTL.String(),
		"suggestions_ttl":     sc.config.SuggestionsTTL.String(),
		"popular_ttl":         sc.config.PopularTTL.String(),
		"filtered_search_ttl": sc.config.FilteredSearchTTL.String(),
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

// ============================================================================
// PHASE 21.2: Advanced Search Cache Methods
// ============================================================================

// Note: We need to import the search service types to avoid circular dependency.
// The cache layer should work with generic map[string]interface{} for flexibility.

// GetFacets retrieves cached facets
func (sc *SearchCache) GetFacets(ctx context.Context, key string) (map[string]interface{}, error) {
	data, err := sc.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			sc.logger.Debug().Str("key", key).Msg("facets cache miss")
			return nil, fmt.Errorf("cache miss")
		}
		sc.logger.Error().Err(err).Str("key", key).Msg("failed to get facets from cache")
		return nil, fmt.Errorf("cache get failed: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		sc.logger.Error().Err(err).Str("key", key).Msg("failed to unmarshal cached facets")
		return nil, fmt.Errorf("cache unmarshal failed: %w", err)
	}

	sc.logger.Debug().Str("key", key).Msg("facets cache hit")
	return result, nil
}

// SetFacets stores facets in cache
func (sc *SearchCache) SetFacets(ctx context.Context, key string, facets map[string]interface{}) error {
	data, err := json.Marshal(facets)
	if err != nil {
		sc.logger.Error().Err(err).Str("key", key).Msg("failed to marshal facets")
		return fmt.Errorf("cache marshal failed: %w", err)
	}

	if err := sc.client.Set(ctx, key, data, sc.config.FacetsTTL).Err(); err != nil {
		sc.logger.Error().Err(err).Str("key", key).Msg("failed to set facets cache")
		return fmt.Errorf("cache set failed: %w", err)
	}

	sc.logger.Debug().Str("key", key).Dur("ttl", sc.config.FacetsTTL).Msg("facets cache set")
	return nil
}

// GetSuggestions retrieves cached suggestions
func (sc *SearchCache) GetSuggestions(ctx context.Context, key string) (map[string]interface{}, error) {
	data, err := sc.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			sc.logger.Debug().Str("key", key).Msg("suggestions cache miss")
			return nil, fmt.Errorf("cache miss")
		}
		sc.logger.Error().Err(err).Str("key", key).Msg("failed to get suggestions from cache")
		return nil, fmt.Errorf("cache get failed: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		sc.logger.Error().Err(err).Str("key", key).Msg("failed to unmarshal cached suggestions")
		return nil, fmt.Errorf("cache unmarshal failed: %w", err)
	}

	sc.logger.Debug().Str("key", key).Msg("suggestions cache hit")
	return result, nil
}

// SetSuggestions stores suggestions in cache
func (sc *SearchCache) SetSuggestions(ctx context.Context, key string, suggestions map[string]interface{}) error {
	data, err := json.Marshal(suggestions)
	if err != nil {
		sc.logger.Error().Err(err).Str("key", key).Msg("failed to marshal suggestions")
		return fmt.Errorf("cache marshal failed: %w", err)
	}

	if err := sc.client.Set(ctx, key, data, sc.config.SuggestionsTTL).Err(); err != nil {
		sc.logger.Error().Err(err).Str("key", key).Msg("failed to set suggestions cache")
		return fmt.Errorf("cache set failed: %w", err)
	}

	sc.logger.Debug().Str("key", key).Dur("ttl", sc.config.SuggestionsTTL).Msg("suggestions cache set")
	return nil
}

// GetPopular retrieves cached popular searches
func (sc *SearchCache) GetPopular(ctx context.Context, key string) (map[string]interface{}, error) {
	data, err := sc.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			sc.logger.Debug().Str("key", key).Msg("popular cache miss")
			return nil, fmt.Errorf("cache miss")
		}
		sc.logger.Error().Err(err).Str("key", key).Msg("failed to get popular from cache")
		return nil, fmt.Errorf("cache get failed: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		sc.logger.Error().Err(err).Str("key", key).Msg("failed to unmarshal cached popular")
		return nil, fmt.Errorf("cache unmarshal failed: %w", err)
	}

	sc.logger.Debug().Str("key", key).Msg("popular cache hit")
	return result, nil
}

// SetPopular stores popular searches in cache
func (sc *SearchCache) SetPopular(ctx context.Context, key string, popular map[string]interface{}) error {
	data, err := json.Marshal(popular)
	if err != nil {
		sc.logger.Error().Err(err).Str("key", key).Msg("failed to marshal popular")
		return fmt.Errorf("cache marshal failed: %w", err)
	}

	if err := sc.client.Set(ctx, key, data, sc.config.PopularTTL).Err(); err != nil {
		sc.logger.Error().Err(err).Str("key", key).Msg("failed to set popular cache")
		return fmt.Errorf("cache set failed: %w", err)
	}

	sc.logger.Debug().Str("key", key).Dur("ttl", sc.config.PopularTTL).Msg("popular cache set")
	return nil
}

// GetFiltered retrieves cached filtered search results
func (sc *SearchCache) GetFiltered(ctx context.Context, key string) (map[string]interface{}, error) {
	data, err := sc.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			sc.logger.Debug().Str("key", key).Msg("filtered search cache miss")
			return nil, fmt.Errorf("cache miss")
		}
		sc.logger.Error().Err(err).Str("key", key).Msg("failed to get filtered search from cache")
		return nil, fmt.Errorf("cache get failed: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		sc.logger.Error().Err(err).Str("key", key).Msg("failed to unmarshal cached filtered search")
		return nil, fmt.Errorf("cache unmarshal failed: %w", err)
	}

	sc.logger.Debug().Str("key", key).Msg("filtered search cache hit")
	return result, nil
}

// SetFiltered stores filtered search results in cache
func (sc *SearchCache) SetFiltered(ctx context.Context, key string, filtered map[string]interface{}) error {
	data, err := json.Marshal(filtered)
	if err != nil {
		sc.logger.Error().Err(err).Str("key", key).Msg("failed to marshal filtered search")
		return fmt.Errorf("cache marshal failed: %w", err)
	}

	if err := sc.client.Set(ctx, key, data, sc.config.FilteredSearchTTL).Err(); err != nil {
		sc.logger.Error().Err(err).Str("key", key).Msg("failed to set filtered search cache")
		return fmt.Errorf("cache set failed: %w", err)
	}

	sc.logger.Debug().Str("key", key).Dur("ttl", sc.config.FilteredSearchTTL).Msg("filtered search cache set")
	return nil
}

// GenerateFacetsKey creates a unique cache key for facets request
func (sc *SearchCache) GenerateFacetsKey(query string, categoryID *int64, filters map[string]interface{}) string {
	parts := []string{
		"q:" + query,
		fmt.Sprintf("cat:%v", categoryID),
	}

	// Add filters if present
	if filters != nil {
		filtersJSON, _ := json.Marshal(filters)
		parts = append(parts, fmt.Sprintf("filters:%s", string(filtersJSON)))
	}

	hash := md5.Sum([]byte(fmt.Sprintf("%v", parts)))
	hashStr := hex.EncodeToString(hash[:])
	return fmt.Sprintf("search:facets:v1:%s", hashStr)
}

// GenerateSuggestionsKey creates a unique cache key for suggestions request
func (sc *SearchCache) GenerateSuggestionsKey(prefix string, categoryID *int64) string {
	parts := []string{
		"prefix:" + prefix,
		fmt.Sprintf("cat:%v", categoryID),
	}

	hash := md5.Sum([]byte(fmt.Sprintf("%v", parts)))
	hashStr := hex.EncodeToString(hash[:])
	return fmt.Sprintf("search:suggestions:v1:%s", hashStr)
}

// GeneratePopularKey creates a unique cache key for popular searches request
func (sc *SearchCache) GeneratePopularKey(categoryID *int64, timeRange string) string {
	parts := []string{
		fmt.Sprintf("cat:%v", categoryID),
		"range:" + timeRange,
	}

	hash := md5.Sum([]byte(fmt.Sprintf("%v", parts)))
	hashStr := hex.EncodeToString(hash[:])
	return fmt.Sprintf("search:popular:v1:%s", hashStr)
}

// GenerateFilteredKey creates a unique cache key for filtered search request
func (sc *SearchCache) GenerateFilteredKey(query string, categoryID *int64, filters map[string]interface{}, sort map[string]string, limit, offset int32) string {
	parts := []string{
		"q:" + query,
		fmt.Sprintf("cat:%v", categoryID),
		fmt.Sprintf("lim:%d", limit),
		fmt.Sprintf("off:%d", offset),
	}

	// Add filters if present
	if filters != nil {
		filtersJSON, _ := json.Marshal(filters)
		parts = append(parts, fmt.Sprintf("filters:%s", string(filtersJSON)))
	}

	// Add sort if present
	if sort != nil {
		sortJSON, _ := json.Marshal(sort)
		parts = append(parts, fmt.Sprintf("sort:%s", string(sortJSON)))
	}

	hash := md5.Sum([]byte(fmt.Sprintf("%v", parts)))
	hashStr := hex.EncodeToString(hash[:])
	return fmt.Sprintf("search:filtered:v1:%s", hashStr)
}
