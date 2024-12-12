//backend/internal/storage/storage.go
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

    // Методы для работы с автомобилями
    CreateCarBooking(ctx context.Context, booking *models.CarBooking) error
    AddCar(ctx context.Context, car *models.Car) (int, error)
    GetAvailableCars(ctx context.Context, filters map[string]string) ([]models.Car, error)
    AddCarImage(ctx context.Context, image *models.CarImage) (int, error)
    GetCarImages(ctx context.Context, carID string) ([]models.CarImage, error)
    DeleteCarImage(ctx context.Context, imageID string) (string, error)
    GetCarFeatures(ctx context.Context) ([]models.CarFeature, error)
    GetCarCategories(ctx context.Context) ([]models.CarCategory, error)
    GetCarWithFeatures(ctx context.Context, carID int) (*models.Car, error)

    // Room methods
    AddRoom(ctx context.Context, room *models.Room) (int, error)
    AddBed(ctx context.Context, roomID int, bedNumber string, pricePerNight float64, hasOutlet bool, hasLight bool, hasShelf bool, bedType string) (int, error)
    GetRooms(ctx context.Context, filters map[string]string, sortBy string, sortDirection string, limit int, offset int) ([]models.Room, int64, error)
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

    
     // Marketplace methods
     CreateListing(ctx context.Context, listing *models.MarketplaceListing) (int, error)
     GetListings(ctx context.Context, filters map[string]string, limit int, offset int) ([]models.MarketplaceListing, int64, error)
     GetListingByID(ctx context.Context, id int) (*models.MarketplaceListing, error)
     UpdateListing(ctx context.Context, listing *models.MarketplaceListing) error
     DeleteListing(ctx context.Context, id int, userID int) error
     GetCategoryTree(ctx context.Context) ([]models.CategoryTreeNode, error)
     
     AddListingImage(ctx context.Context, image *models.MarketplaceImage) (int, error)
     GetListingImages(ctx context.Context, listingID string) ([]models.MarketplaceImage, error)
     DeleteListingImage(ctx context.Context, imageID string) (string, error)
     
     GetCategories(ctx context.Context) ([]models.MarketplaceCategory, error)
     GetCategoryByID(ctx context.Context, id int) (*models.MarketplaceCategory, error)
     
     AddToFavorites(ctx context.Context, userID int, listingID int) error
     RemoveFromFavorites(ctx context.Context, userID int, listingID int) error
     GetUserFavorites(ctx context.Context, userID int) ([]models.MarketplaceListing, error)

    // Database connection
    Close()
    Ping(ctx context.Context) error
}