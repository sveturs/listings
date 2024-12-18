// backend/internal/proj/marketplace/service/interface.go
package service

import (
    "context"
    "backend/internal/domain/models"
    "mime/multipart"
    
)

type MarketplaceServiceInterface interface {
    CreateListing(ctx context.Context, listing *models.MarketplaceListing) (int, error)
    GetListings(ctx context.Context, filters map[string]string, limit int, offset int) ([]models.MarketplaceListing, int64, error)
    GetListingByID(ctx context.Context, id int) (*models.MarketplaceListing, error)
    UpdateListing(ctx context.Context, listing *models.MarketplaceListing) error
    DeleteListing(ctx context.Context, id int, userID int) error
    ProcessImage(file *multipart.FileHeader) (string, error)
    AddListingImage(ctx context.Context, image *models.MarketplaceImage) (int, error)
    GetCategories(ctx context.Context) ([]models.MarketplaceCategory, error)
    AddToFavorites(ctx context.Context, userID int, listingID int) error
    RemoveFromFavorites(ctx context.Context, userID int, listingID int) error
    GetCategoryTree(ctx context.Context) ([]models.CategoryTreeNode, error)
}