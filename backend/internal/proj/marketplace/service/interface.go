package service

import (
    "backend/internal/domain/models"
    "backend/internal/domain/search"
    "context"
    "mime/multipart"
    "backend/internal/storage"  
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
    GetUserFavorites(ctx context.Context, userID int) ([]models.MarketplaceListing, error)
    GetFavoritedUsers(ctx context.Context, listingID int) ([]int, error)
    UpdateTranslation(ctx context.Context, translation *models.Translation) error
    GetSubcategories(ctx context.Context, parentID string, limit int, offset int) ([]models.CategoryTreeNode, error)
    RefreshCategoryListingCounts(ctx context.Context) error
    
    // OpenSearch методы
    SearchListingsAdvanced(ctx context.Context, params *search.ServiceParams) (*search.ServiceResult, error)
    GetSuggestions(ctx context.Context, prefix string, size int) ([]string, error)
    ReindexAllListings(ctx context.Context) error
    GetCategorySuggestions(ctx context.Context, query string, size int) ([]models.CategorySuggestion, error)
    Storage() storage.Storage
    //атрибуты
    GetCategoryAttributes(ctx context.Context, categoryID int) ([]models.CategoryAttribute, error)
    SaveListingAttributes(ctx context.Context, listingID int, attributes []models.ListingAttributeValue) error


}