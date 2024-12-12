//backend/internal/services/interfaces.go
package services

import (
    "context"
    "mime/multipart"
    "backend/internal/domain/models"
    "backend/internal/types"
)

type RoomServiceInterface interface {
    CreateRoom(ctx context.Context, room *models.Room) (int, error)
    GetRooms(ctx context.Context, filters map[string]string, sortBy string, sortDirection string, limit int, offset int) ([]models.Room, int64, error)
    GetRoomByID(ctx context.Context, id int) (*models.Room, error)
    ProcessImage(file *multipart.FileHeader) (string, error)
    AddRoomImage(ctx context.Context, image *models.RoomImage) (int, error)
    GetRoomImages(ctx context.Context, roomID string) ([]models.RoomImage, error)
    AddBedImage(ctx context.Context, image *models.RoomImage) (int, error)
    GetBedImages(ctx context.Context, bedID string) ([]models.RoomImage, error)
    DeleteRoomImage(ctx context.Context, imageID string) (string, error)
    AddBed(ctx context.Context, roomID int, bedNumber string, pricePerNight float64, hasOutlet bool, hasLight bool, hasShelf bool, bedType string) (int, error)
    GetAvailableBeds(ctx context.Context, roomID string, startDate string, endDate string) ([]models.Bed, error)
}
type CarServiceInterface interface {
    AddCar(ctx context.Context, car *models.Car) (int, error)
    GetAvailableCars(ctx context.Context, filters map[string]string) ([]models.Car, error)
    ProcessImage(file *multipart.FileHeader) (string, error)
    AddCarImage(ctx context.Context, image *models.CarImage) (int, error)
    GetCarImages(ctx context.Context, carID string) ([]models.CarImage, error)
    CreateBooking(ctx context.Context, booking *models.CarBooking) error
}

type AuthServiceInterface interface {
    GetGoogleAuthURL() string
    HandleGoogleCallback(ctx context.Context, code string) (*types.SessionData, error)
    SaveSession(token string, data *types.SessionData)
    GetSession(token string) (*types.SessionData, bool)
    DeleteSession(token string)
}

type BookingServiceInterface interface {
    CreateBooking(ctx context.Context, userID int, booking *models.BookingRequest) error
    GetAllBookings(ctx context.Context) ([]models.Booking, error)
    DeleteBooking(ctx context.Context, bookingID string, bookingType string) error
    GetAvailableBeds(ctx context.Context, roomID string, startDate string, endDate string) ([]models.Bed, error)
    AddBed(ctx context.Context, roomID int, bedNumber string, pricePerNight float64, hasOutlet bool, hasLight bool, hasShelf bool, bedType string) (int, error)
}

type UserServiceInterface interface {
    GetUserByID(ctx context.Context, id int) (*models.User, error)
    GetUserByEmail(ctx context.Context, email string) (*models.User, error)
    CreateUser(ctx context.Context, user *models.User) error
    UpdateUser(ctx context.Context, user *models.User) error
}

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