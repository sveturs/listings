package grpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	categoriesv2 "github.com/vondi-global/listings/api/proto/categories/v2"
	"github.com/vondi-global/listings/internal/cache"
	"github.com/vondi-global/listings/internal/domain"
)

// GetCategoryTreeV2 retrieves category tree with localization and caching
func (s *Server) GetCategoryTreeV2(ctx context.Context, req *categoriesv2.GetCategoryTreeV2Request) (*categoriesv2.GetCategoryTreeV2Response, error) {
	s.logger.Debug().
		Interface("root_id", req.RootId).
		Str("locale", req.Locale).
		Bool("active_only", req.ActiveOnly).
		Interface("max_depth", req.MaxDepth).
		Msg("GetCategoryTreeV2 called")

	// Validate locale
	locale := strings.ToLower(req.Locale)
	if locale == "" || !isValidLocale(locale) {
		return nil, status.Error(codes.InvalidArgument, "locale must be one of: sr, en, ru")
	}

	// Build cache key
	cacheKey := buildTreeCacheKey(req.RootId, locale, req.ActiveOnly, req.MaxDepth)

	// Try cache first
	if s.categoryCache != nil {
		cachedTree, err := s.categoryCache.GetCategoryTree(ctx, cacheKey)
		if err == nil && cachedTree != nil {
			s.logger.Debug().Str("cache_key", cacheKey).Msg("category tree cache hit")
			return &categoriesv2.GetCategoryTreeV2Response{
				Tree: DomainToProtoCategoryTreeV2(cachedTree),
			}, nil
		}
	}

	// Build filter
	filter := &domain.GetCategoryTreeFilterV2{
		Locale:     locale,
		ActiveOnly: req.ActiveOnly,
		MaxDepth:   nil,
	}

	// Parse root_id if provided
	if req.RootId != nil && *req.RootId != "" {
		rootUUID, err := uuid.Parse(*req.RootId)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid root_id UUID: %v", err))
		}
		filter.RootID = &rootUUID
	}

	// Parse max_depth if provided
	if req.MaxDepth != nil {
		filter.MaxDepth = req.MaxDepth
	}

	// Fetch from repository
	tree, err := s.categoryRepoV2.GetTreeV2(ctx, filter)
	if err != nil {
		s.logger.Error().Err(err).Interface("filter", filter).Msg("failed to get category tree V2")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get category tree: %v", err))
	}

	// Cache result (ignore cache errors)
	if s.categoryCache != nil {
		if err := s.categoryCache.SetCategoryTree(ctx, cacheKey, tree, cache.DefaultCategoryTTL); err != nil {
			s.logger.Warn().Err(err).Str("cache_key", cacheKey).Msg("failed to cache category tree")
		}
	}

	s.logger.Debug().Int("tree_count", len(tree)).Msg("category tree V2 retrieved")

	return &categoriesv2.GetCategoryTreeV2Response{
		Tree: DomainToProtoCategoryTreeV2(tree),
	}, nil
}

// GetCategoryBySlugV2 retrieves a single category by slug with localization
func (s *Server) GetCategoryBySlugV2(ctx context.Context, req *categoriesv2.GetCategoryBySlugV2Request) (*categoriesv2.GetCategoryBySlugV2Response, error) {
	s.logger.Debug().
		Str("slug", req.Slug).
		Str("locale", req.Locale).
		Msg("GetCategoryBySlugV2 called")

	// Validate slug
	if req.Slug == "" {
		return nil, status.Error(codes.InvalidArgument, "slug is required")
	}

	// Validate locale
	locale := strings.ToLower(req.Locale)
	if locale == "" || !isValidLocale(locale) {
		return nil, status.Error(codes.InvalidArgument, "locale must be one of: sr, en, ru")
	}

	// Try cache first
	if s.categoryCache != nil {
		cachedCat, err := s.categoryCache.GetCategoryBySlug(ctx, req.Slug)
		if err == nil && cachedCat != nil {
			s.logger.Debug().Str("slug", req.Slug).Msg("category slug cache hit")
			localized := cachedCat.Localize(locale)
			return &categoriesv2.GetCategoryBySlugV2Response{
				Category: DomainToProtoCategoryV2(localized),
			}, nil
		}
	}

	// Fetch from repository
	category, err := s.categoryRepoV2.GetBySlugV2(ctx, req.Slug)
	if err != nil {
		s.logger.Error().Err(err).Str("slug", req.Slug).Msg("failed to get category by slug V2")
		if contains(err.Error(), "not found") {
			return nil, status.Error(codes.NotFound, "category not found")
		}
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get category: %v", err))
	}

	// Cache result (ignore cache errors)
	if s.categoryCache != nil {
		if err := s.categoryCache.SetCategoryBySlug(ctx, req.Slug, category, cache.DefaultCategoryTTL); err != nil {
			s.logger.Warn().Err(err).Str("slug", req.Slug).Msg("failed to cache category by slug")
		}
	}

	// Localize
	localized := category.Localize(locale)

	s.logger.Debug().Str("slug", req.Slug).Str("id", category.ID.String()).Msg("category by slug V2 retrieved")

	return &categoriesv2.GetCategoryBySlugV2Response{
		Category: DomainToProtoCategoryV2(localized),
	}, nil
}

// GetCategoryByUUID retrieves a single category by UUID with localization
func (s *Server) GetCategoryByUUID(ctx context.Context, req *categoriesv2.GetCategoryByUUIDRequest) (*categoriesv2.GetCategoryByUUIDResponse, error) {
	s.logger.Debug().
		Str("id", req.Id).
		Str("locale", req.Locale).
		Msg("GetCategoryByUUID called")

	// Validate UUID
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if _, err := uuid.Parse(req.Id); err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid UUID: %v", err))
	}

	// Validate locale
	locale := strings.ToLower(req.Locale)
	if locale == "" || !isValidLocale(locale) {
		return nil, status.Error(codes.InvalidArgument, "locale must be one of: sr, en, ru")
	}

	// Try cache first
	if s.categoryCache != nil {
		cachedCat, err := s.categoryCache.GetCategoryByUUID(ctx, req.Id)
		if err == nil && cachedCat != nil {
			s.logger.Debug().Str("id", req.Id).Msg("category UUID cache hit")
			localized := cachedCat.Localize(locale)
			return &categoriesv2.GetCategoryByUUIDResponse{
				Category: DomainToProtoCategoryV2(localized),
			}, nil
		}
	}

	// Fetch from repository
	category, err := s.categoryRepoV2.GetByUUID(ctx, req.Id)
	if err != nil {
		s.logger.Error().Err(err).Str("id", req.Id).Msg("failed to get category by UUID")
		if contains(err.Error(), "not found") {
			return nil, status.Error(codes.NotFound, "category not found")
		}
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get category: %v", err))
	}

	// Cache result (ignore cache errors)
	if s.categoryCache != nil {
		if err := s.categoryCache.SetCategoryByUUID(ctx, req.Id, category, cache.DefaultCategoryTTL); err != nil {
			s.logger.Warn().Err(err).Str("id", req.Id).Msg("failed to cache category by UUID")
		}
	}

	// Localize
	localized := category.Localize(locale)

	s.logger.Debug().Str("id", req.Id).Msg("category by UUID retrieved")

	return &categoriesv2.GetCategoryByUUIDResponse{
		Category: DomainToProtoCategoryV2(localized),
	}, nil
}

// GetBreadcrumb retrieves breadcrumb trail for a category
func (s *Server) GetBreadcrumb(ctx context.Context, req *categoriesv2.GetBreadcrumbRequest) (*categoriesv2.GetBreadcrumbResponse, error) {
	s.logger.Debug().
		Str("category_id", req.CategoryId).
		Str("locale", req.Locale).
		Msg("GetBreadcrumb called")

	// Validate category_id
	if req.CategoryId == "" {
		return nil, status.Error(codes.InvalidArgument, "category_id is required")
	}

	if _, err := uuid.Parse(req.CategoryId); err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("invalid category_id UUID: %v", err))
	}

	// Validate locale
	locale := strings.ToLower(req.Locale)
	if locale == "" || !isValidLocale(locale) {
		return nil, status.Error(codes.InvalidArgument, "locale must be one of: sr, en, ru")
	}

	// Try cache first
	if s.categoryCache != nil {
		cachedBreadcrumbs, err := s.categoryCache.GetBreadcrumb(ctx, req.CategoryId, locale)
		if err == nil && cachedBreadcrumbs != nil {
			s.logger.Debug().Str("category_id", req.CategoryId).Str("locale", locale).Msg("breadcrumb cache hit")
			return &categoriesv2.GetBreadcrumbResponse{
				Breadcrumbs: DomainToProtoBreadcrumb(cachedBreadcrumbs),
			}, nil
		}
	}

	// Fetch from repository
	breadcrumbs, err := s.categoryRepoV2.GetBreadcrumb(ctx, req.CategoryId, locale)
	if err != nil {
		s.logger.Error().Err(err).Str("category_id", req.CategoryId).Msg("failed to get breadcrumb")
		if contains(err.Error(), "not found") {
			return nil, status.Error(codes.NotFound, "category not found")
		}
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get breadcrumb: %v", err))
	}

	// Cache result (ignore cache errors)
	if s.categoryCache != nil {
		if err := s.categoryCache.SetBreadcrumb(ctx, req.CategoryId, locale, breadcrumbs, cache.DefaultCategoryTTL); err != nil {
			s.logger.Warn().Err(err).Str("category_id", req.CategoryId).Msg("failed to cache breadcrumb")
		}
	}

	s.logger.Debug().Str("category_id", req.CategoryId).Int("count", len(breadcrumbs)).Msg("breadcrumb retrieved")

	return &categoriesv2.GetBreadcrumbResponse{
		Breadcrumbs: DomainToProtoBreadcrumb(breadcrumbs),
	}, nil
}

// Helper functions

// isValidLocale checks if locale is supported
func isValidLocale(locale string) bool {
	validLocales := map[string]bool{"sr": true, "en": true, "ru": true}
	return validLocales[locale]
}

// buildTreeCacheKey generates a cache key for category tree
func buildTreeCacheKey(rootID *string, locale string, activeOnly bool, maxDepth *int32) string {
	key := locale
	if rootID != nil && *rootID != "" {
		key += ":" + *rootID
	} else {
		key += ":root"
	}
	if activeOnly {
		key += ":active"
	}
	if maxDepth != nil {
		key += fmt.Sprintf(":depth%d", *maxDepth)
	}
	return key
}
