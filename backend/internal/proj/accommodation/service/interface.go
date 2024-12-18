package service

import (
    "context"
    "mime/multipart"
    "backend/internal/domain/models"
)
type RoomServiceInterface interface {
    CreateRoom(ctx context.Context, room *models.Room) (int, error)
    GetRooms(ctx context.Context, filters map[string]string, sortBy string, sortDirection string, limit int, offset int) ([]models.Room, int64, error)
    GetRoomByID(ctx context.Context, id int) (*models.Room, error)
    ProcessImage(file *multipart.FileHeader) (string, error)
    AddRoomImage(ctx context.Context, image *models.RoomImage) (int, error)
    GetRoomImages(ctx context.Context, roomID string) ([]models.RoomImage, error)
    DeleteRoomImage(ctx context.Context, imageID string) (string, error)
    AddBed(ctx context.Context, roomID int, bedNumber string, pricePerNight float64, hasOutlet bool, hasLight bool, hasShelf bool, bedType string) (int, error)
    GetAvailableBeds(ctx context.Context, roomID string, startDate string, endDate string) ([]models.Bed, error)
    GetBedImages(ctx context.Context, bedID string) ([]models.RoomImage, error)
    AddBedImage(ctx context.Context, image *models.RoomImage) (int, error)
}
type BookingServiceInterface interface {
    CreateBooking(ctx context.Context, userID int, booking *models.BookingRequest) error
    GetAllBookings(ctx context.Context) ([]models.Booking, error)
    DeleteBooking(ctx context.Context, bookingID string, bookingType string) error
    GetAvailableBeds(ctx context.Context, roomID string, startDate string, endDate string) ([]models.Bed, error)
    AddBed(ctx context.Context, roomID int, bedNumber string, pricePerNight float64, hasOutlet bool, hasLight bool, hasShelf bool, bedType string) (int, error)
}