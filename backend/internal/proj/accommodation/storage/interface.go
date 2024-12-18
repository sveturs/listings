// internal/proj/accommodation/storage/interface.go
package storage

import (
    "context"
    "backend/internal/domain/models"
)

type Repository interface {
    // Room methods
    AddRoom(ctx context.Context, room *models.Room) (int, error)
    GetRooms(ctx context.Context, filters map[string]string, sortBy string, sortDirection string, limit int, offset int) ([]models.Room, int64, error)
    GetRoomByID(ctx context.Context, id int) (*models.Room, error)

    // Bed methods
    AddBed(ctx context.Context, roomID int, bedNumber string, pricePerNight float64, hasOutlet bool, hasLight bool, hasShelf bool, bedType string) (int, error)
    GetAvailableBeds(ctx context.Context, roomID string, startDate string, endDate string) ([]models.Bed, error)
    GetBedByID(ctx context.Context, id int) (*models.Bed, error)
    GetBedsByRoomID(ctx context.Context, roomID int) ([]models.Bed, error)
    UpdateBedAttributes(ctx context.Context, bedID int, bedReq *models.BedRequest) error
    UpdateBedAvailability(ctx context.Context, bedID int, isAvailable bool) error
    UpdateBedPrice(ctx context.Context, bedID int, price float64) error

    // Image methods
    AddRoomImage(ctx context.Context, image *models.RoomImage) (int, error)
    GetRoomImages(ctx context.Context, roomID string) ([]models.RoomImage, error)
    DeleteRoomImage(ctx context.Context, imageID string) (string, error)
    AddBedImage(ctx context.Context, image *models.RoomImage) (int, error)
    GetBedImages(ctx context.Context, bedID string) ([]models.RoomImage, error)

    // Booking methods
    CreateBooking(ctx context.Context, booking *models.BookingRequest) error
    GetAllBookings(ctx context.Context) ([]models.Booking, error)
    DeleteBooking(ctx context.Context, bookingID string, bookingType string) error
}