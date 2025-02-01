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
func (s *MarketplaceService) GetUserFavorites(ctx context.Context, userID int) ([]models.MarketplaceListing, error) {
    return s.storage.GetUserFavorites(ctx, userID)
}
func (s *MarketplaceService) CreateListing(ctx context.Context, listing *models.MarketplaceListing) (int, error) {
    // Устанавливаем начальные значения
    if listing.OriginalLanguage == "" {
        detectedLang, _, err := s.translationService.DetectLanguage(ctx, listing.Title)
        if err != nil {
            return 0, fmt.Errorf("failed to detect language: %w", err)
        }
        listing.OriginalLanguage = detectedLang
    }
    listing.Status = "active"
    listing.ViewsCount = 0
    
    return s.storage.CreateListing(ctx, listing)
}

func (s *MarketplaceService) GetListings(ctx context.Context, filters map[string]string, limit int, offset int) ([]models.MarketplaceListing, int64, error) {
    return s.storage.GetListings(ctx, filters, limit, offset)
}
func (s *MarketplaceService) GetFavoritedUsers(ctx context.Context, listingID int) ([]int, error) {
    return s.storage.GetFavoritedUsers(ctx, listingID)
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
    // Получаем расширение файла
    ext := filepath.Ext(file.Filename)
    if ext == "" {
        // Если расширение отсутствует, определяем его по MIME-типу
        switch file.Header.Get("Content-Type") {
        case "image/jpeg", "image/jpg":
            ext = ".jpg"
        case "image/png":
            ext = ".png"
        case "image/gif":
            ext = ".gif"
        case "image/webp":
            ext = ".webp"
        default:
            ext = ".jpg" // По умолчанию используем .jpg
        }
    }

    // Генерируем уникальное имя файла с расширением
    fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

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