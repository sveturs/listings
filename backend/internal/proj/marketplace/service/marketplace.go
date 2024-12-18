// backend/internal/proj/marketplace/service/marketplace.go
package service

import (
    "backend/internal/domain/models"
    "backend/internal/storage"
    "context"
    "mime/multipart"
    "path/filepath"
    "fmt"
    "time"
)

type MarketplaceService struct {
    storage storage.Storage
}

func NewMarketplaceService(storage storage.Storage) MarketplaceServiceInterface {
    return &MarketplaceService{
        storage: storage,
    }
}

func (s *MarketplaceService) CreateListing(ctx context.Context, listing *models.MarketplaceListing) (int, error) {
    // Устанавливаем начальные значения
    listing.Status = "active"
    listing.ViewsCount = 0
    
    return s.storage.CreateListing(ctx, listing)
}

func (s *MarketplaceService) GetListings(ctx context.Context, filters map[string]string, limit int, offset int) ([]models.MarketplaceListing, int64, error) {
    return s.storage.GetListings(ctx, filters, limit, offset)
}

func (s *MarketplaceService) GetListingByID(ctx context.Context, id int) (*models.MarketplaceListing, error) {
    return s.storage.GetListingByID(ctx, id)
}

func (s *MarketplaceService) UpdateListing(ctx context.Context, listing *models.MarketplaceListing) error {
    return s.storage.UpdateListing(ctx, listing)
}
func (s *MarketplaceService) GetCategoryTree(ctx context.Context) ([]models.CategoryTreeNode, error) {
    return s.storage.GetCategoryTree(ctx)
}
func (s *MarketplaceService) DeleteListing(ctx context.Context, id int, userID int) error {
    // Проверяем, что пользователь является владельцем объявления
    listing, err := s.storage.GetListingByID(ctx, id)
    if err != nil {
        return err
    }
    
    if listing.UserID != userID {
        return fmt.Errorf("не хватает прав для удаления объявления")
    }
    
    return s.storage.DeleteListing(ctx, id, userID)
}

func (s *MarketplaceService) ProcessImage(file *multipart.FileHeader) (string, error) {
    ext := filepath.Ext(file.Filename)
    fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

    // Используем переменную в возвращаемом значении
    return fileName, nil
}

func (s *MarketplaceService) AddListingImage(ctx context.Context, image *models.MarketplaceImage) (int, error) {
    return s.storage.AddListingImage(ctx, image)
}

func (s *MarketplaceService) GetCategories(ctx context.Context) ([]models.MarketplaceCategory, error) {
    return s.storage.GetCategories(ctx)
}

func (s *MarketplaceService) AddToFavorites(ctx context.Context, userID int, listingID int) error {
    return s.storage.AddToFavorites(ctx, userID, listingID)
}

func (s *MarketplaceService) RemoveFromFavorites(ctx context.Context, userID int, listingID int) error {
    return s.storage.RemoveFromFavorites(ctx, userID, listingID)
}