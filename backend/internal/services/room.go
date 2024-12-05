package services

import (
	"backend/internal/domain/models"
    "backend/internal/storage"
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
	"log"
	"github.com/disintegration/imaging"
)
var _ RoomServiceInterface = (*RoomService)(nil) 
type RoomService struct {
    storage storage.Storage
}

func NewRoomService(storage storage.Storage) *RoomService {
    return &RoomService{
        storage: storage,
    }
}
func (s *RoomService) AddBedImage(ctx context.Context, image *models.RoomImage) (int, error) {
    return s.storage.AddBedImage(ctx, image)
}
func (s *RoomService) GetBedImages(ctx context.Context, bedID string) ([]models.RoomImage, error) {
    return s.storage.GetBedImages(ctx, bedID)
}
func (s *RoomService) CreateRoom(ctx context.Context, room *models.Room) (int, error) {
    return s.storage.AddRoom(ctx, room)
}

func (s *RoomService) GetRooms(ctx context.Context, filters map[string]string) ([]models.Room, error) {
    log.Printf("RoomService.GetRooms called with filters: %+v", filters)
    
    rooms, err := s.storage.GetRooms(ctx, filters)
    if err != nil {
        log.Printf("Error in RoomService.GetRooms: %v", err)
        return nil, fmt.Errorf("error getting rooms from storage: %w", err)
    }
    
    return rooms, nil
}

func (s *RoomService) ProcessImage(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	ext := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join("uploads", fileName)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	img, err := imaging.Decode(src)
	if err != nil {
		return "", err
	}

	resized := imaging.Resize(img, 1200, 0, imaging.Lanczos)
	err = imaging.Save(resized, filePath)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func (s *RoomService) AddRoomImage(ctx context.Context, image *models.RoomImage) (int, error) {
    return s.storage.AddRoomImage(ctx, image)
}

func (s *RoomService) GetRoomImages(ctx context.Context, roomID string) ([]models.RoomImage, error) {
	return s.storage.GetRoomImages(ctx, roomID)
}

func (s *RoomService) DeleteRoomImage(ctx context.Context, imageID string) (string, error) {
	return s.storage.DeleteRoomImage(ctx, imageID)
}

func (s *RoomService) GetRoomByID(ctx context.Context, id int) (*models.Room, error) {
    return s.storage.GetRoomByID(ctx, id)
}

func (s *RoomService) AddBed(ctx context.Context, roomID int, bedNumber string, pricePerNight float64, hasOutlet bool, hasLight bool, hasShelf bool, bedType string) (int, error) {
    return s.storage.AddBed(ctx, roomID, bedNumber, pricePerNight, hasOutlet, hasLight, hasShelf, bedType)
}

func (s *RoomService) GetAvailableBeds(ctx context.Context, roomID string, startDate string, endDate string) ([]models.Bed, error) {
    return s.storage.GetAvailableBeds(ctx, roomID, startDate, endDate)
}