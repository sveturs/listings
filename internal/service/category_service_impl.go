// Package service implements business logic for the listings microservice.
package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"

	"github.com/vondi-global/listings/internal/domain"
)

// CategoryRepository defines minimal repository interface needed by the service
type CategoryRepository interface {
	GetCategoriesWithPagination(ctx context.Context, parentID *string, isActive *bool, limit, offset int32) ([]*domain.Category, int32, error)
	GetCategoryByID(ctx context.Context, id string) (*domain.Category, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*domain.Category, error)
	GetCategoryTree(ctx context.Context, categoryID string) (*domain.CategoryTreeNode, error)
	CreateCategory(ctx context.Context, cat *domain.Category) (*domain.Category, error)
	UpdateCategory(ctx context.Context, cat *domain.Category) (*domain.Category, error)
	DeleteCategory(ctx context.Context, categoryID string) error
}

const (
	// Category cache keys and TTL
	categoryCacheKeyByID   = "category:id:%s"
	categoryCacheKeyBySlug = "category:slug:%s"
	categoryTreeCacheKey   = "category:tree:%s"
	categoryCacheTTL       = 1 * time.Hour
)

// CategoryServiceImpl implements CategoryService interface
type CategoryServiceImpl struct {
	repo   CategoryRepository
	cache  *CategoryCache
	logger zerolog.Logger
}

// CategoryCache provides caching functionality for categories
type CategoryCache struct {
	client redis.UniversalClient
	logger zerolog.Logger
}

// NewCategoryCache creates a new category cache service
func NewCategoryCache(client redis.UniversalClient, logger zerolog.Logger) *CategoryCache {
	return &CategoryCache{
		client: client,
		logger: logger.With().Str("component", "category_cache").Logger(),
	}
}

// NewCategoryService creates a new category service
func NewCategoryService(
	repo CategoryRepository,
	cacheClient redis.UniversalClient,
	logger zerolog.Logger,
) CategoryService {
	return &CategoryServiceImpl{
		repo:   repo,
		cache:  NewCategoryCache(cacheClient, logger),
		logger: logger.With().Str("component", "category_service").Logger(),
	}
}

// =============================================================================
// Public Read Operations
// =============================================================================

// GetCategories retrieves a list of categories with optional filtering and pagination
func (s *CategoryServiceImpl) GetCategories(ctx context.Context, parentID *string, isActive *bool, limit, offset int32) ([]*domain.Category, int32, error) {
	categories, total, err := s.repo.GetCategoriesWithPagination(ctx, parentID, isActive, limit, offset)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get categories")
		return nil, 0, fmt.Errorf("failed to get categories: %w", err)
	}

	s.logger.Debug().
		Interface("parent_id", parentID).
		Interface("is_active", isActive).
		Int32("limit", limit).
		Int32("offset", offset).
		Int32("total", total).
		Int("count", len(categories)).
		Msg("categories retrieved successfully")

	return categories, total, nil
}

// GetCategory retrieves a single category by ID (with caching)
func (s *CategoryServiceImpl) GetCategory(ctx context.Context, id string) (*domain.Category, error) {
	// Try cache first
	cacheKey := fmt.Sprintf(categoryCacheKeyByID, id)
	cached, err := s.cache.Get(ctx, cacheKey)
	if err == nil && cached != nil {
		s.logger.Debug().Str("id", id).Msg("category retrieved from cache")
		return cached, nil
	}

	// Cache miss - fetch from repository
	s.logger.Debug().Str("id", id).Msg("cache miss, fetching from repository")
	category, err := s.repo.GetCategoryByID(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Str("id", id).Msg("failed to get category")
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	// Cache the result
	if err := s.cache.Set(ctx, cacheKey, category, categoryCacheTTL); err != nil {
		s.logger.Warn().Err(err).Str("id", id).Msg("failed to cache category")
		// Don't fail on cache error
	}

	return category, nil
}

// GetCategoryBySlug retrieves a single category by slug (with caching)
func (s *CategoryServiceImpl) GetCategoryBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	// Try cache first
	cacheKey := fmt.Sprintf(categoryCacheKeyBySlug, slug)
	cached, err := s.cache.Get(ctx, cacheKey)
	if err == nil && cached != nil {
		s.logger.Debug().Str("slug", slug).Msg("category retrieved from cache")
		return cached, nil
	}

	// Cache miss - fetch from repository
	s.logger.Debug().Str("slug", slug).Msg("cache miss, fetching from repository")
	category, err := s.repo.GetCategoryBySlug(ctx, slug)
	if err != nil {
		s.logger.Error().Err(err).Str("slug", slug).Msg("failed to get category by slug")
		return nil, fmt.Errorf("failed to get category by slug: %w", err)
	}

	// Cache the result (both by slug and by ID)
	if err := s.cache.Set(ctx, cacheKey, category, categoryCacheTTL); err != nil {
		s.logger.Warn().Err(err).Str("slug", slug).Msg("failed to cache category by slug")
	}
	idKey := fmt.Sprintf(categoryCacheKeyByID, category.ID)
	if err := s.cache.Set(ctx, idKey, category, categoryCacheTTL); err != nil {
		s.logger.Warn().Err(err).Str("id", category.ID).Msg("failed to cache category by id")
	}

	return category, nil
}

// GetCategoryTree retrieves the full category tree starting from a specific category (with caching)
func (s *CategoryServiceImpl) GetCategoryTree(ctx context.Context, categoryID string) (*domain.CategoryTreeNode, error) {
	// Try cache first
	cacheKey := fmt.Sprintf(categoryTreeCacheKey, categoryID)
	cached, err := s.cache.GetTree(ctx, cacheKey)
	if err == nil && cached != nil {
		s.logger.Debug().Str("category_id", categoryID).Msg("category tree retrieved from cache")
		return cached, nil
	}

	// Cache miss - fetch from repository
	s.logger.Debug().Str("category_id", categoryID).Msg("cache miss, fetching tree from repository")
	tree, err := s.repo.GetCategoryTree(ctx, categoryID)
	if err != nil {
		s.logger.Error().Err(err).Str("category_id", categoryID).Msg("failed to get category tree")
		return nil, fmt.Errorf("failed to get category tree: %w", err)
	}

	// Cache the result
	if err := s.cache.SetTree(ctx, cacheKey, tree, categoryCacheTTL); err != nil {
		s.logger.Warn().Err(err).Str("category_id", categoryID).Msg("failed to cache category tree")
		// Don't fail on cache error
	}

	return tree, nil
}

// =============================================================================
// Admin Write Operations
// =============================================================================

// CreateCategory creates a new category with validation
func (s *CategoryServiceImpl) CreateCategory(ctx context.Context, cat *domain.Category) (*domain.Category, error) {
	// Validate input
	if err := s.validateCategory(cat); err != nil {
		s.logger.Warn().Err(err).Str("name", cat.Name).Msg("category validation failed")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check slug uniqueness
	if cat.Slug != "" {
		existing, err := s.repo.GetCategoryBySlug(ctx, cat.Slug)
		if err == nil && existing != nil {
			return nil, fmt.Errorf("category with slug '%s' already exists", cat.Slug)
		}
	}

	// Verify parent exists if specified
	if cat.ParentID != nil {
		parent, err := s.repo.GetCategoryByID(ctx, *cat.ParentID)
		if err != nil || parent == nil {
			return nil, fmt.Errorf("parent category with id %s not found", *cat.ParentID)
		}
	}

	// Create category
	created, err := s.repo.CreateCategory(ctx, cat)
	if err != nil {
		s.logger.Error().Err(err).Str("name", cat.Name).Msg("failed to create category")
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	// Invalidate parent's tree cache if this is a subcategory
	if cat.ParentID != nil {
		s.invalidateTreeCache(ctx, *cat.ParentID)
	}

	s.logger.Info().Str("id", created.ID).Str("name", created.Name).Msg("category created successfully")
	return created, nil
}

// UpdateCategory updates an existing category with validation
func (s *CategoryServiceImpl) UpdateCategory(ctx context.Context, cat *domain.Category) (*domain.Category, error) {
	// Validate input
	if cat.ID == "" {
		return nil, fmt.Errorf("category ID is required")
	}

	// Get existing category
	existing, err := s.repo.GetCategoryByID(ctx, cat.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing category: %w", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("category with id %s not found", cat.ID)
	}

	// Validate updated category
	if err := s.validateCategory(cat); err != nil {
		s.logger.Warn().Err(err).Str("id", cat.ID).Msg("category validation failed")
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check slug uniqueness (if changed)
	if cat.Slug != existing.Slug {
		slugExists, err := s.repo.GetCategoryBySlug(ctx, cat.Slug)
		if err == nil && slugExists != nil && slugExists.ID != cat.ID {
			return nil, fmt.Errorf("category with slug '%s' already exists", cat.Slug)
		}
	}

	// Verify parent exists if changed
	if cat.ParentID != nil {
		if existing.ParentID == nil || *cat.ParentID != *existing.ParentID {
			parent, err := s.repo.GetCategoryByID(ctx, *cat.ParentID)
			if err != nil || parent == nil {
				return nil, fmt.Errorf("parent category with id %s not found", *cat.ParentID)
			}

			// Prevent circular dependency
			if *cat.ParentID == cat.ID {
				return nil, fmt.Errorf("category cannot be its own parent")
			}
		}
	}

	// Update category
	updated, err := s.repo.UpdateCategory(ctx, cat)
	if err != nil {
		s.logger.Error().Err(err).Str("id", cat.ID).Msg("failed to update category")
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	// Invalidate cache
	if err := s.InvalidateCache(ctx, cat.ID); err != nil {
		s.logger.Warn().Err(err).Str("id", cat.ID).Msg("failed to invalidate cache after update")
	}

	// Invalidate parent's tree cache if parent changed
	if existing.ParentID != nil {
		s.invalidateTreeCache(ctx, *existing.ParentID)
	}
	if cat.ParentID != nil && (existing.ParentID == nil || *cat.ParentID != *existing.ParentID) {
		s.invalidateTreeCache(ctx, *cat.ParentID)
	}

	s.logger.Info().Str("id", updated.ID).Msg("category updated successfully")
	return updated, nil
}

// DeleteCategory soft-deletes a category (sets is_active=false)
func (s *CategoryServiceImpl) DeleteCategory(ctx context.Context, categoryID string) error {
	// Get existing category
	existing, err := s.repo.GetCategoryByID(ctx, categoryID)
	if err != nil {
		return fmt.Errorf("failed to get category: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("category with id %s not found", categoryID)
	}

	// Delete category
	if err := s.repo.DeleteCategory(ctx, categoryID); err != nil {
		s.logger.Error().Err(err).Str("id", categoryID).Msg("failed to delete category")
		return fmt.Errorf("failed to delete category: %w", err)
	}

	// Invalidate cache
	if err := s.InvalidateCache(ctx, categoryID); err != nil {
		s.logger.Warn().Err(err).Str("id", categoryID).Msg("failed to invalidate cache after delete")
	}

	// Invalidate parent's tree cache
	if existing.ParentID != nil {
		s.invalidateTreeCache(ctx, *existing.ParentID)
	}

	s.logger.Info().Str("id", categoryID).Msg("category deleted successfully")
	return nil
}

// =============================================================================
// Cache Management
// =============================================================================

// InvalidateCache invalidates all cache entries for a category
func (s *CategoryServiceImpl) InvalidateCache(ctx context.Context, categoryID string) error {
	// Get category to get slug for cache invalidation
	category, err := s.repo.GetCategoryByID(ctx, categoryID)
	if err != nil {
		s.logger.Warn().Err(err).Str("id", categoryID).Msg("failed to get category for cache invalidation")
		// Continue with ID-based invalidation even if we can't get the slug
	}

	// Invalidate by ID
	idKey := fmt.Sprintf(categoryCacheKeyByID, categoryID)
	if err := s.cache.client.Del(ctx, idKey).Err(); err != nil {
		s.logger.Warn().Err(err).Str("key", idKey).Msg("failed to delete cache key")
	}

	// Invalidate by slug if we have it
	if category != nil {
		slugKey := fmt.Sprintf(categoryCacheKeyBySlug, category.Slug)
		if err := s.cache.client.Del(ctx, slugKey).Err(); err != nil {
			s.logger.Warn().Err(err).Str("key", slugKey).Msg("failed to delete cache key")
		}
	}

	// Invalidate tree cache
	s.invalidateTreeCache(ctx, categoryID)

	s.logger.Debug().Str("category_id", categoryID).Msg("category cache invalidated")
	return nil
}

// invalidateTreeCache invalidates the tree cache for a category
func (s *CategoryServiceImpl) invalidateTreeCache(ctx context.Context, categoryID string) {
	treeKey := fmt.Sprintf(categoryTreeCacheKey, categoryID)
	if err := s.cache.client.Del(ctx, treeKey).Err(); err != nil {
		s.logger.Warn().Err(err).Str("key", treeKey).Msg("failed to delete tree cache key")
	}
}

// =============================================================================
// Validation Helpers
// =============================================================================

// validateCategory validates category fields
func (s *CategoryServiceImpl) validateCategory(cat *domain.Category) error {
	if cat == nil {
		return fmt.Errorf("category cannot be nil")
	}

	// Validate name
	name := strings.TrimSpace(cat.Name)
	if len(name) < 2 {
		return fmt.Errorf("category name must be at least 2 characters long")
	}
	if len(name) > 100 {
		return fmt.Errorf("category name must be at most 100 characters long")
	}

	// Validate slug if provided
	if cat.Slug != "" {
		slug := strings.TrimSpace(cat.Slug)
		if len(slug) < 2 || len(slug) > 100 {
			return fmt.Errorf("category slug must be between 2 and 100 characters long")
		}
	}

	return nil
}

// =============================================================================
// Cache Operations (CategoryCache)
// =============================================================================

// Get retrieves a category from cache
func (c *CategoryCache) Get(ctx context.Context, key string) (*domain.Category, error) {
	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil // Cache miss
	}
	if err != nil {
		c.logger.Warn().Err(err).Str("key", key).Msg("cache get error")
		return nil, err
	}

	var category domain.Category
	if err := json.Unmarshal(data, &category); err != nil {
		c.logger.Error().Err(err).Str("key", key).Msg("failed to unmarshal cached category")
		return nil, err
	}

	return &category, nil
}

// Set stores a category in cache
func (c *CategoryCache) Set(ctx context.Context, key string, category *domain.Category, ttl time.Duration) error {
	data, err := json.Marshal(category)
	if err != nil {
		c.logger.Error().Err(err).Str("key", key).Msg("failed to marshal category for cache")
		return err
	}

	if err := c.client.Set(ctx, key, data, ttl).Err(); err != nil {
		c.logger.Warn().Err(err).Str("key", key).Msg("cache set error")
		return err
	}

	return nil
}

// GetTree retrieves a category tree from cache
func (c *CategoryCache) GetTree(ctx context.Context, key string) (*domain.CategoryTreeNode, error) {
	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil // Cache miss
	}
	if err != nil {
		c.logger.Warn().Err(err).Str("key", key).Msg("cache get error")
		return nil, err
	}

	var tree domain.CategoryTreeNode
	if err := json.Unmarshal(data, &tree); err != nil {
		c.logger.Error().Err(err).Str("key", key).Msg("failed to unmarshal cached category tree")
		return nil, err
	}

	return &tree, nil
}

// SetTree stores a category tree in cache
func (c *CategoryCache) SetTree(ctx context.Context, key string, tree *domain.CategoryTreeNode, ttl time.Duration) error {
	data, err := json.Marshal(tree)
	if err != nil {
		c.logger.Error().Err(err).Str("key", key).Msg("failed to marshal category tree for cache")
		return err
	}

	if err := c.client.Set(ctx, key, data, ttl).Err(); err != nil {
		c.logger.Warn().Err(err).Str("key", key).Msg("cache set error")
		return err
	}

	return nil
}
