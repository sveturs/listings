// TEMPORARY: Will be moved to microservice
package storage

import (
	"context"

	"backend/internal/domain/models"
)

type MarketplaceStorage interface {
	// Categories
	GetCategories(ctx context.Context, lang string) ([]models.MarketplaceCategory, error)
	GetPopularCategories(ctx context.Context, lang string, limit int) ([]models.MarketplaceCategory, error)
	GetCategoryByID(ctx context.Context, id int) (*models.MarketplaceCategory, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*models.MarketplaceCategory, error)

	// Favorites
	GetUserFavorites(ctx context.Context, userID int) ([]int, error)
	AddToFavorites(ctx context.Context, userID, listingID int) error
	RemoveFromFavorites(ctx context.Context, userID, listingID int) error
	IsFavorite(ctx context.Context, userID, listingID int) (bool, error)

	// Attributes
	GetCategoryAttributes(ctx context.Context, categoryID int) ([]models.CategoryAttribute, error)
	GetVariantAttributes(ctx context.Context, categorySlug string) ([]models.CategoryVariantAttribute, error)

	// Storefronts (B2C)
	GetStorefronts(ctx context.Context, filters StorefrontFilters) ([]models.Storefront, int, error)
	GetStorefrontByID(ctx context.Context, id int) (*models.Storefront, error)
	CreateStorefront(ctx context.Context, userID int, dto *models.StorefrontCreateDTO) (*models.Storefront, error)

	// Listings (TEMPORARY: proxies to microservice)
	CreateListing(ctx context.Context, userID int, categoryID int, title string, description *string, price float64, currency string, quantity int32, sku *string, storefrontID *int) (*models.MarketplaceListing, error)
	GetListing(ctx context.Context, listingID int) (*models.MarketplaceListing, error)
}

// StorefrontFilters параметры фильтрации витрин
type StorefrontFilters struct {
	IsActive  *bool
	Page      int
	Limit     int
	SortBy    string // "products_count", "rating", "created_at"
	SortOrder string // "asc", "desc"
}
