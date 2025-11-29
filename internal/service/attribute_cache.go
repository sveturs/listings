// Package service implements business logic for the listings microservice.
package service

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
	cacheKeyAttributeID     = "attr:id:%d"
	cacheKeyAttributeCode   = "attr:code:%s"
	cacheKeyCategoryAttrs   = "cat_attrs:%d"
	cacheKeyCategoryVariant = "cat_variant_attrs:%d"
	cacheKeyListingAttrs    = "listing_attrs:%d"

	// Cache TTL (30 minutes as per architecture)
	cacheTTL = 30 * time.Minute
)

// AttributeCache provides caching functionality for attributes
type AttributeCache struct {
	client redis.UniversalClient
	logger zerolog.Logger
}

// NewAttributeCache creates a new attribute cache service
func NewAttributeCache(client redis.UniversalClient, logger zerolog.Logger) *AttributeCache {
	return &AttributeCache{
		client: client,
		logger: logger.With().Str("component", "attribute_cache").Logger(),
	}
}

// GetAttribute retrieves a cached attribute by ID
func (c *AttributeCache) GetAttribute(ctx context.Context, id int32) (*domain.Attribute, error) {
	key := fmt.Sprintf(cacheKeyAttributeID, id)
	return c.getAttribute(ctx, key)
}

// GetAttributeByCode retrieves a cached attribute by code
func (c *AttributeCache) GetAttributeByCode(ctx context.Context, code string) (*domain.Attribute, error) {
	key := fmt.Sprintf(cacheKeyAttributeCode, code)
	return c.getAttribute(ctx, key)
}

// SetAttribute caches an attribute by ID and code
func (c *AttributeCache) SetAttribute(ctx context.Context, attr *domain.Attribute) error {
	if attr == nil {
		return fmt.Errorf("attribute cannot be nil")
	}

	data, err := json.Marshal(attr)
	if err != nil {
		return fmt.Errorf("failed to marshal attribute: %w", err)
	}

	// Cache by ID
	keyID := fmt.Sprintf(cacheKeyAttributeID, attr.ID)
	if err := c.client.Set(ctx, keyID, data, cacheTTL).Err(); err != nil {
		c.logger.Warn().Err(err).Int32("id", attr.ID).Msg("failed to cache attribute by ID")
		return err
	}

	// Cache by code
	keyCode := fmt.Sprintf(cacheKeyAttributeCode, attr.Code)
	if err := c.client.Set(ctx, keyCode, data, cacheTTL).Err(); err != nil {
		c.logger.Warn().Err(err).Str("code", attr.Code).Msg("failed to cache attribute by code")
		return err
	}

	c.logger.Debug().Int32("id", attr.ID).Str("code", attr.Code).Msg("attribute cached")
	return nil
}

// GetCategoryAttributes retrieves cached category attributes
func (c *AttributeCache) GetCategoryAttributes(ctx context.Context, categoryID int32) ([]*domain.CategoryAttribute, error) {
	key := fmt.Sprintf(cacheKeyCategoryAttrs, categoryID)

	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		c.logger.Warn().Err(err).Int32("category_id", categoryID).Msg("failed to get category attributes from cache")
		return nil, err
	}

	var attrs []*domain.CategoryAttribute
	if err := json.Unmarshal(data, &attrs); err != nil {
		c.logger.Error().Err(err).Int32("category_id", categoryID).Msg("failed to unmarshal category attributes")
		// Delete corrupted cache entry
		_ = c.client.Del(ctx, key).Err()
		return nil, err
	}

	c.logger.Debug().Int32("category_id", categoryID).Int("count", len(attrs)).Msg("category attributes cache hit")
	return attrs, nil
}

// SetCategoryAttributes caches category attributes
func (c *AttributeCache) SetCategoryAttributes(ctx context.Context, categoryID int32, attrs []*domain.CategoryAttribute) error {
	key := fmt.Sprintf(cacheKeyCategoryAttrs, categoryID)

	data, err := json.Marshal(attrs)
	if err != nil {
		return fmt.Errorf("failed to marshal category attributes: %w", err)
	}

	if err := c.client.Set(ctx, key, data, cacheTTL).Err(); err != nil {
		c.logger.Warn().Err(err).Int32("category_id", categoryID).Msg("failed to cache category attributes")
		return err
	}

	c.logger.Debug().Int32("category_id", categoryID).Int("count", len(attrs)).Msg("category attributes cached")
	return nil
}

// GetCategoryVariantAttributes retrieves cached variant attributes for a category
func (c *AttributeCache) GetCategoryVariantAttributes(ctx context.Context, categoryID int32) ([]*domain.VariantAttribute, error) {
	key := fmt.Sprintf(cacheKeyCategoryVariant, categoryID)

	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		c.logger.Warn().Err(err).Int32("category_id", categoryID).Msg("failed to get variant attributes from cache")
		return nil, err
	}

	var attrs []*domain.VariantAttribute
	if err := json.Unmarshal(data, &attrs); err != nil {
		c.logger.Error().Err(err).Int32("category_id", categoryID).Msg("failed to unmarshal variant attributes")
		// Delete corrupted cache entry
		_ = c.client.Del(ctx, key).Err()
		return nil, err
	}

	c.logger.Debug().Int32("category_id", categoryID).Int("count", len(attrs)).Msg("variant attributes cache hit")
	return attrs, nil
}

// SetCategoryVariantAttributes caches variant attributes for a category
func (c *AttributeCache) SetCategoryVariantAttributes(ctx context.Context, categoryID int32, attrs []*domain.VariantAttribute) error {
	key := fmt.Sprintf(cacheKeyCategoryVariant, categoryID)

	data, err := json.Marshal(attrs)
	if err != nil {
		return fmt.Errorf("failed to marshal variant attributes: %w", err)
	}

	if err := c.client.Set(ctx, key, data, cacheTTL).Err(); err != nil {
		c.logger.Warn().Err(err).Int32("category_id", categoryID).Msg("failed to cache variant attributes")
		return err
	}

	c.logger.Debug().Int32("category_id", categoryID).Int("count", len(attrs)).Msg("variant attributes cached")
	return nil
}

// GetListingAttributes retrieves cached listing attribute values
func (c *AttributeCache) GetListingAttributes(ctx context.Context, listingID int32) ([]*domain.ListingAttributeValue, error) {
	key := fmt.Sprintf(cacheKeyListingAttrs, listingID)

	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		c.logger.Warn().Err(err).Int32("listing_id", listingID).Msg("failed to get listing attributes from cache")
		return nil, err
	}

	var attrs []*domain.ListingAttributeValue
	if err := json.Unmarshal(data, &attrs); err != nil {
		c.logger.Error().Err(err).Int32("listing_id", listingID).Msg("failed to unmarshal listing attributes")
		// Delete corrupted cache entry
		_ = c.client.Del(ctx, key).Err()
		return nil, err
	}

	c.logger.Debug().Int32("listing_id", listingID).Int("count", len(attrs)).Msg("listing attributes cache hit")
	return attrs, nil
}

// SetListingAttributes caches listing attribute values
func (c *AttributeCache) SetListingAttributes(ctx context.Context, listingID int32, attrs []*domain.ListingAttributeValue) error {
	key := fmt.Sprintf(cacheKeyListingAttrs, listingID)

	data, err := json.Marshal(attrs)
	if err != nil {
		return fmt.Errorf("failed to marshal listing attributes: %w", err)
	}

	if err := c.client.Set(ctx, key, data, cacheTTL).Err(); err != nil {
		c.logger.Warn().Err(err).Int32("listing_id", listingID).Msg("failed to cache listing attributes")
		return err
	}

	c.logger.Debug().Int32("listing_id", listingID).Int("count", len(attrs)).Msg("listing attributes cached")
	return nil
}

// InvalidateAttribute invalidates cache for a specific attribute (by ID and code)
func (c *AttributeCache) InvalidateAttribute(ctx context.Context, id int32, code string) error {
	keys := []string{
		fmt.Sprintf(cacheKeyAttributeID, id),
		fmt.Sprintf(cacheKeyAttributeCode, code),
	}

	for _, key := range keys {
		if err := c.client.Del(ctx, key).Err(); err != nil {
			c.logger.Warn().Err(err).Str("key", key).Msg("failed to invalidate cache key")
		}
	}

	c.logger.Debug().Int32("id", id).Str("code", code).Msg("attribute cache invalidated")
	return nil
}

// InvalidateCategory invalidates all cache entries for a category
func (c *AttributeCache) InvalidateCategory(ctx context.Context, categoryID int32) error {
	keys := []string{
		fmt.Sprintf(cacheKeyCategoryAttrs, categoryID),
		fmt.Sprintf(cacheKeyCategoryVariant, categoryID),
	}

	for _, key := range keys {
		if err := c.client.Del(ctx, key).Err(); err != nil {
			c.logger.Warn().Err(err).Str("key", key).Msg("failed to invalidate cache key")
		}
	}

	c.logger.Debug().Int32("category_id", categoryID).Msg("category cache invalidated")
	return nil
}

// InvalidateListing invalidates cache for listing attributes
func (c *AttributeCache) InvalidateListing(ctx context.Context, listingID int32) error {
	key := fmt.Sprintf(cacheKeyListingAttrs, listingID)

	if err := c.client.Del(ctx, key).Err(); err != nil {
		c.logger.Warn().Err(err).Str("key", key).Msg("failed to invalidate cache key")
		return err
	}

	c.logger.Debug().Int32("listing_id", listingID).Msg("listing cache invalidated")
	return nil
}

// getAttribute is a helper method to retrieve and unmarshal an attribute from cache
func (c *AttributeCache) getAttribute(ctx context.Context, key string) (*domain.Attribute, error) {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		c.logger.Warn().Err(err).Str("key", key).Msg("failed to get attribute from cache")
		return nil, err
	}

	var attr domain.Attribute
	if err := json.Unmarshal(data, &attr); err != nil {
		c.logger.Error().Err(err).Str("key", key).Msg("failed to unmarshal attribute")
		// Delete corrupted cache entry
		_ = c.client.Del(ctx, key).Err()
		return nil, err
	}

	c.logger.Debug().Str("key", key).Msg("attribute cache hit")
	return &attr, nil
}
