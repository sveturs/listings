// internal/proj/marketplace/storage/interface.go
package storage

import (
    "context"
    "backend/internal/domain/models"
)

type MarketplaceRepository interface {
    // Listings
    CreateListing(ctx context.Context, listing *models.MarketplaceListing) (int, error)
    GetListings(ctx context.Context, filters map[string]string, limit int, offset int) ([]models.MarketplaceListing, int64, error)
    GetListingByID(ctx context.Context, id int) (*models.MarketplaceListing, error)
    UpdateListing(ctx context.Context, listing *models.MarketplaceListing) error
    DeleteListing(ctx context.Context, id int, userID int) error
    
    // Categories
    GetCategories(ctx context.Context) ([]models.MarketplaceCategory, error)
    GetCategoryByID(ctx context.Context, id int) (*models.MarketplaceCategory, error)
    GetCategoryTree(ctx context.Context) ([]models.CategoryTreeNode, error)
    
    // Images
    AddListingImage(ctx context.Context, image *models.MarketplaceImage) (int, error)
    GetListingImages(ctx context.Context, listingID string) ([]models.MarketplaceImage, error)
    DeleteListingImage(ctx context.Context, imageID string) (string, error)
    
    // Favorites
    AddToFavorites(ctx context.Context, userID int, listingID int) error
    RemoveFromFavorites(ctx context.Context, userID int, listingID int) error
    GetUserFavorites(ctx context.Context, userID int) ([]models.MarketplaceListing, error)
}