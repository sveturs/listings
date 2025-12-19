package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/vondi-global/listings/internal/domain"
)

const (
	// Cache key prefixes
	CategoryTreePrefix       = "category:tree:"
	CategoryBySlugPrefix     = "category:slug:"
	CategoryByUUIDPrefix     = "category:uuid:"
	CategoryBreadcrumbPrefix = "category:breadcrumb:"

	// Default TTL
	DefaultCategoryTTL = 1 * time.Hour
)

// CategoryCache handles category caching operations
type CategoryCache struct {
	redis  *redis.Client
	logger zerolog.Logger
}

// NewCategoryCache creates a new category cache instance
func NewCategoryCache(redisClient *redis.Client, logger zerolog.Logger) *CategoryCache {
	return &CategoryCache{
		redis:  redisClient,
		logger: logger.With().Str("component", "category_cache").Logger(),
	}
}

// GetCategoryTree retrieves cached category tree
func (c *CategoryCache) GetCategoryTree(ctx context.Context, key string) ([]*domain.CategoryTreeV2, error) {
	fullKey := CategoryTreePrefix + key
	val, err := c.redis.Get(ctx, fullKey).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	}
	if err != nil {
		c.logger.Warn().Err(err).Str("key", fullKey).Msg("failed to get category tree from cache")
		return nil, err
	}

	var tree []*domain.CategoryTreeV2
	if err := json.Unmarshal([]byte(val), &tree); err != nil {
		c.logger.Error().Err(err).Str("key", fullKey).Msg("failed to unmarshal category tree")
		return nil, err
	}

	c.logger.Debug().Str("key", fullKey).Msg("category tree cache hit")
	return tree, nil
}

// SetCategoryTree caches category tree
func (c *CategoryCache) SetCategoryTree(ctx context.Context, key string, tree []*domain.CategoryTreeV2, ttl time.Duration) error {
	fullKey := CategoryTreePrefix + key
	data, err := json.Marshal(tree)
	if err != nil {
		c.logger.Error().Err(err).Msg("failed to marshal category tree")
		return err
	}

	if ttl == 0 {
		ttl = DefaultCategoryTTL
	}

	if err := c.redis.Set(ctx, fullKey, data, ttl).Err(); err != nil {
		c.logger.Warn().Err(err).Str("key", fullKey).Msg("failed to cache category tree")
		return err
	}

	c.logger.Debug().Str("key", fullKey).Dur("ttl", ttl).Msg("cached category tree")
	return nil
}

// GetCategoryBySlug retrieves cached category by slug
func (c *CategoryCache) GetCategoryBySlug(ctx context.Context, slug string) (*domain.CategoryV2, error) {
	fullKey := CategoryBySlugPrefix + slug
	val, err := c.redis.Get(ctx, fullKey).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	}
	if err != nil {
		c.logger.Warn().Err(err).Str("key", fullKey).Msg("failed to get category from cache")
		return nil, err
	}

	var cat domain.CategoryV2
	if err := json.Unmarshal([]byte(val), &cat); err != nil {
		c.logger.Error().Err(err).Str("key", fullKey).Msg("failed to unmarshal category")
		return nil, err
	}

	c.logger.Debug().Str("key", fullKey).Msg("category cache hit")
	return &cat, nil
}

// SetCategoryBySlug caches category by slug
func (c *CategoryCache) SetCategoryBySlug(ctx context.Context, slug string, cat *domain.CategoryV2, ttl time.Duration) error {
	fullKey := CategoryBySlugPrefix + slug
	data, err := json.Marshal(cat)
	if err != nil {
		c.logger.Error().Err(err).Msg("failed to marshal category")
		return err
	}

	if ttl == 0 {
		ttl = DefaultCategoryTTL
	}

	if err := c.redis.Set(ctx, fullKey, data, ttl).Err(); err != nil {
		c.logger.Warn().Err(err).Str("key", fullKey).Msg("failed to cache category")
		return err
	}

	c.logger.Debug().Str("key", fullKey).Dur("ttl", ttl).Msg("cached category")
	return nil
}

// GetCategoryByUUID retrieves cached category by UUID
func (c *CategoryCache) GetCategoryByUUID(ctx context.Context, uuid string) (*domain.CategoryV2, error) {
	fullKey := CategoryByUUIDPrefix + uuid
	val, err := c.redis.Get(ctx, fullKey).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	}
	if err != nil {
		c.logger.Warn().Err(err).Str("key", fullKey).Msg("failed to get category from cache")
		return nil, err
	}

	var cat domain.CategoryV2
	if err := json.Unmarshal([]byte(val), &cat); err != nil {
		c.logger.Error().Err(err).Str("key", fullKey).Msg("failed to unmarshal category")
		return nil, err
	}

	c.logger.Debug().Str("key", fullKey).Msg("category cache hit")
	return &cat, nil
}

// SetCategoryByUUID caches category by UUID
func (c *CategoryCache) SetCategoryByUUID(ctx context.Context, uuid string, cat *domain.CategoryV2, ttl time.Duration) error {
	fullKey := CategoryByUUIDPrefix + uuid
	data, err := json.Marshal(cat)
	if err != nil {
		c.logger.Error().Err(err).Msg("failed to marshal category")
		return err
	}

	if ttl == 0 {
		ttl = DefaultCategoryTTL
	}

	if err := c.redis.Set(ctx, fullKey, data, ttl).Err(); err != nil {
		c.logger.Warn().Err(err).Str("key", fullKey).Msg("failed to cache category")
		return err
	}

	c.logger.Debug().Str("key", fullKey).Dur("ttl", ttl).Msg("cached category")
	return nil
}

// GetBreadcrumb retrieves cached breadcrumb
func (c *CategoryCache) GetBreadcrumb(ctx context.Context, categoryID, locale string) ([]*domain.CategoryBreadcrumb, error) {
	fullKey := fmt.Sprintf("%s%s:%s", CategoryBreadcrumbPrefix, categoryID, locale)
	val, err := c.redis.Get(ctx, fullKey).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	}
	if err != nil {
		c.logger.Warn().Err(err).Str("key", fullKey).Msg("failed to get breadcrumb from cache")
		return nil, err
	}

	var breadcrumbs []*domain.CategoryBreadcrumb
	if err := json.Unmarshal([]byte(val), &breadcrumbs); err != nil {
		c.logger.Error().Err(err).Str("key", fullKey).Msg("failed to unmarshal breadcrumb")
		return nil, err
	}

	c.logger.Debug().Str("key", fullKey).Msg("breadcrumb cache hit")
	return breadcrumbs, nil
}

// SetBreadcrumb caches breadcrumb
func (c *CategoryCache) SetBreadcrumb(ctx context.Context, categoryID, locale string, breadcrumbs []*domain.CategoryBreadcrumb, ttl time.Duration) error {
	fullKey := fmt.Sprintf("%s%s:%s", CategoryBreadcrumbPrefix, categoryID, locale)
	data, err := json.Marshal(breadcrumbs)
	if err != nil {
		c.logger.Error().Err(err).Msg("failed to marshal breadcrumb")
		return err
	}

	if ttl == 0 {
		ttl = DefaultCategoryTTL
	}

	if err := c.redis.Set(ctx, fullKey, data, ttl).Err(); err != nil {
		c.logger.Warn().Err(err).Str("key", fullKey).Msg("failed to cache breadcrumb")
		return err
	}

	c.logger.Debug().Str("key", fullKey).Dur("ttl", ttl).Msg("cached breadcrumb")
	return nil
}

// InvalidateCategoryCache invalidates all category caches matching pattern
func (c *CategoryCache) InvalidateCategoryCache(ctx context.Context, pattern string) error {
	// Use SCAN to find matching keys (safer than KEYS *)
	var cursor uint64
	deletedCount := 0

	for {
		keys, nextCursor, err := c.redis.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			c.logger.Error().Err(err).Str("pattern", pattern).Msg("failed to scan keys")
			return err
		}

		if len(keys) > 0 {
			if err := c.redis.Del(ctx, keys...).Err(); err != nil {
				c.logger.Warn().Err(err).Strs("keys", keys).Msg("failed to delete keys")
			} else {
				deletedCount += len(keys)
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	c.logger.Info().Str("pattern", pattern).Int("deleted", deletedCount).Msg("invalidated category cache")
	return nil
}

// InvalidateAll invalidates all category caches
func (c *CategoryCache) InvalidateAll(ctx context.Context) error {
	patterns := []string{
		CategoryTreePrefix + "*",
		CategoryBySlugPrefix + "*",
		CategoryByUUIDPrefix + "*",
		CategoryBreadcrumbPrefix + "*",
	}

	for _, pattern := range patterns {
		if err := c.InvalidateCategoryCache(ctx, pattern); err != nil {
			return err
		}
	}

	return nil
}

// InvalidateCategory invalidates all caches for a specific category
func (c *CategoryCache) InvalidateCategory(ctx context.Context, categoryID, slug string) error {
	keys := []string{
		CategoryByUUIDPrefix + categoryID,
		CategoryBySlugPrefix + slug,
	}

	// Also invalidate breadcrumbs for all locales
	breadcrumbPattern := CategoryBreadcrumbPrefix + categoryID + ":*"
	if err := c.InvalidateCategoryCache(ctx, breadcrumbPattern); err != nil {
		c.logger.Warn().Err(err).Msg("failed to invalidate breadcrumbs")
	}

	// Delete category-specific keys
	if err := c.redis.Del(ctx, keys...).Err(); err != nil {
		c.logger.Warn().Err(err).Strs("keys", keys).Msg("failed to delete category keys")
		return err
	}

	c.logger.Debug().Str("category_id", categoryID).Str("slug", slug).Msg("invalidated category cache")
	return nil
}
