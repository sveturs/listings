// backend/internal/proj/accommodation/service/bed.go
package service

import (
    "context"
    "backend/internal/domain/models"
    "backend/internal/storage"
	"fmt"
)

type BedService struct {
    storage storage.Storage
}

func NewBedService(storage storage.Storage) BedServiceInterface {
    return &BedService{
        storage: storage,
    }
}

func (s *BedService) AddBed(ctx context.Context, roomID int, bedNumber string, pricePerNight float64, hasOutlet bool, hasLight bool, hasShelf bool, bedType string) (int, error) {
    return s.storage.AddBed(ctx, roomID, bedNumber, pricePerNight, hasOutlet, hasLight, hasShelf, bedType)
}

func (s *BedService) GetAvailableBeds(ctx context.Context, roomID string, startDate string, endDate string) ([]models.Bed, error) {
    return s.storage.GetAvailableBeds(ctx, roomID, startDate, endDate)
}

func (s *BedService) GetBedImages(ctx context.Context, bedID string) ([]models.RoomImage, error) {
    return s.storage.GetBedImages(ctx, bedID)
}

func (s *BedService) AddBedImage(ctx context.Context, image *models.RoomImage) (int, error) {
    return s.storage.AddBedImage(ctx, image)
}

func (s *BedService) UpdateBedAttributes(ctx context.Context, bedID int, bedReq *models.BedRequest) error {
    // Проверяем существование кровати перед обновлением
    bed, err := s.storage.GetBedByID(ctx, bedID)
    if err != nil {
        return err
    }

    // Проверяем валидность bedType
    if bedReq.BedType != "top" && bedReq.BedType != "bottom" && bedReq.BedType != "single" {
        bedReq.BedType = "single" // Устанавливаем значение по умолчанию
    }

    // Проверяем, что цена не отрицательная
    if bedReq.PricePerNight < 0 {
        bedReq.PricePerNight = bed.PricePerNight // Оставляем текущую цену
    }

    return s.storage.UpdateBedAttributes(ctx, bedID, bedReq)
}

func (s *BedService) GetBedByID(ctx context.Context, id int) (*models.Bed, error) {
    return s.storage.GetBedByID(ctx, id)
}

func (s *BedService) GetBedsByRoom(ctx context.Context, roomID int) ([]models.Bed, error) {
    return s.storage.GetBedsByRoomID(ctx, roomID)
}

func (s *BedService) UpdateBedAvailability(ctx context.Context, bedID int, isAvailable bool) error {
    // Проверяем существование кровати
    if _, err := s.storage.GetBedByID(ctx, bedID); err != nil {
        return err
    }
    return s.storage.UpdateBedAvailability(ctx, bedID, isAvailable)
}

func (s *BedService) UpdateBedPrice(ctx context.Context, bedID int, price float64) error {
    if price < 0 {
        return fmt.Errorf("price cannot be negative")
    }
    
    // Проверяем существование кровати
    if _, err := s.storage.GetBedByID(ctx, bedID); err != nil {
        return err
    }
    return s.storage.UpdateBedPrice(ctx, bedID, price)
}