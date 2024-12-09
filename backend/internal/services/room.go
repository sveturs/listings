//backend/internal/services/room.go
package services

import (
    "backend/internal/domain/models"
    "backend/internal/storage"
    "context"
    "errors"
    "fmt"
    "mime/multipart"
    "os"
	"path/filepath"
    "strconv"
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

func (s *RoomService) GetRooms(ctx context.Context, filters map[string]string, sortBy string, sortDirection string, limit int, offset int) ([]models.Room, int64, error) {
    // Проверяем и нормализуем параметры сортировки
    if sortDirection != "asc" && sortDirection != "desc" {
        sortDirection = "asc"
    }
    
    validSortFields := map[string]bool{
        "price_per_night": true,
        "rating": true,
        "created_at": true,
    }
    if !validSortFields[sortBy] {
        sortBy = "created_at"
    }

    // Валидация фильтров
    if filters["capacity"] != "" {
        if _, err := strconv.Atoi(filters["capacity"]); err != nil {
            return nil, 0, errors.New("invalid capacity value")
        }
    }

    if filters["start_date"] != "" {
        if _, err := time.Parse("2006-01-02", filters["start_date"]); err != nil {
            return nil, 0, errors.New("invalid start_date format")
        }
    }

    if filters["end_date"] != "" {
        if _, err := time.Parse("2006-01-02", filters["end_date"]); err != nil {
            return nil, 0, errors.New("invalid end_date format")
        }
    }

    // Получаем комнаты и общее количество
    rooms, total, err := s.storage.GetRooms(ctx, filters, sortBy, sortDirection, limit, offset)
    if err != nil {
        return nil, 0, fmt.Errorf("error getting rooms: %w", err)
    }

    // Для каждой комнаты получаем изображения
    for i := range rooms {
        images, err := s.storage.GetRoomImages(ctx, strconv.Itoa(rooms[i].ID))
        if err != nil {
            log.Printf("Error getting images for room %d: %v", rooms[i].ID, err)
            continue
        }
        rooms[i].Images = images
    }

    return rooms, total, nil
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