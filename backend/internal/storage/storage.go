package storage

import (
    "context"
    "backend/internal/domain/models"
)

type Storage interface {
    // User methods
    GetOrCreateGoogleUser(ctx context.Context, user *models.User) (*models.User, error)
    GetUserByEmail(ctx context.Context, email string) (*models.User, error)
    GetUserByID(ctx context.Context, id int) (*models.User, error)
    CreateUser(ctx context.Context, user *models.User) error
    UpdateUser(ctx context.Context, user *models.User) error

    // Room methods
    AddRoom(ctx context.Context, room *models.Room) (int, error)
    AddBed(ctx context.Context, roomID int, bedNumber string, pricePerNight float64, hasOutlet bool, hasLight bool, hasShelf bool, bedType string) (int, error)
    GetRooms(ctx context.Context, filters map[string]string) ([]models.Room, error)
    GetRoomByID(ctx context.Context, id int) (*models.Room, error)
    AddRoomImage(ctx context.Context, image *models.RoomImage) (int, error)
    GetRoomImages(ctx context.Context, roomID string) ([]models.RoomImage, error)
    DeleteRoomImage(ctx context.Context, imageID string) (string, error)
    GetBedImages(ctx context.Context, bedID string) ([]models.RoomImage, error)
    AddBedImage(ctx context.Context, image *models.RoomImage) (int, error)
    
    // Booking methods
    CreateBooking(ctx context.Context, booking *models.BookingRequest) error
    GetAllBookings(ctx context.Context) ([]models.Booking, error)
    DeleteBooking(ctx context.Context, bookingID string, bookingType string) error
    GetAvailableBeds(ctx context.Context, roomID string, startDate string, endDate string) ([]models.Bed, error)

    // Database connection
    Close()
    Ping(ctx context.Context) error
}