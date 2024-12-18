// backend/internal/storage/postgres/db.go
package postgres

import (
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/jackc/pgx/v5"
    "backend/internal/storage"
    "backend/internal/domain/models" 
    "fmt"
    //"log"
    //"strconv"

    carStorage "backend/internal/proj/car/storage/postgres"
    marketplaceStorage "backend/internal/proj/marketplace/storage/postgres"
    reviewStorage "backend/internal/proj/reviews/storage/postgres"
    accommodationStorage "backend/internal/proj/accommodation/storage/postgres"
)

type Database struct {
    pool *pgxpool.Pool
    carDB *carStorage.Storage
    marketplaceDB *marketplaceStorage.Storage
    reviewDB *reviewStorage.Storage
    accommodationDB *accommodationStorage.Storage
    
}

func NewDatabase(dbURL string) (*Database, error) {
    pool, err := pgxpool.New(context.Background(), dbURL)
    if err != nil {
        return nil, fmt.Errorf("error creating connection pool: %w", err)
    }

    return &Database{
        pool:          pool,
        carDB:         carStorage.NewStorage(pool),
        marketplaceDB: marketplaceStorage.NewStorage(pool),
        reviewDB:      reviewStorage.NewStorage(pool),
        accommodationDB: accommodationStorage.NewStorage(pool),
    }, nil
}

var _ storage.Storage = (*Database)(nil) 

func (db *Database) Close() {
    if db.pool != nil {
        db.pool.Close()
    }
}

func (db *Database) Ping(ctx context.Context) error {
    return db.pool.Ping(ctx)
}

type RowsWrapper struct {
    rows pgx.Rows
}

func (r *RowsWrapper) Next() bool {
    return r.rows.Next()
}

func (r *RowsWrapper) Scan(dest ...interface{}) error {
    return r.rows.Scan(dest...)
}

func (r *RowsWrapper) Close() error {
    r.rows.Close()
    return nil
}

func (db *Database) Query(ctx context.Context, sql string, args ...interface{}) (storage.Rows, error) {
    rows, err := db.pool.Query(ctx, sql, args...)
    if err != nil {
        return nil, err
    }
    return &RowsWrapper{rows: rows}, nil
}

func (db *Database) QueryRow(ctx context.Context, sql string, args ...interface{}) storage.Row {
    return db.pool.QueryRow(ctx, sql, args...)
}

// Car methods
func (db *Database) AddCar(ctx context.Context, car *models.Car) (int, error) {
    return db.carDB.AddCar(ctx, car)
}

func (db *Database) GetAvailableCars(ctx context.Context, filters map[string]string) ([]models.Car, error) {
    return db.carDB.GetAvailableCars(ctx, filters)
}

func (db *Database) GetCarWithFeatures(ctx context.Context, carID int) (*models.Car, error) {
    return db.carDB.GetCarWithFeatures(ctx, carID)
}

func (db *Database) GetCarFeatures(ctx context.Context) ([]models.CarFeature, error) {
    return db.carDB.GetCarFeatures(ctx)
}

func (db *Database) GetCarCategories(ctx context.Context) ([]models.CarCategory, error) {
    return db.carDB.GetCarCategories(ctx)
}

func (db *Database) CreateCarBooking(ctx context.Context, booking *models.CarBooking) error {
    return db.carDB.CreateCarBooking(ctx, booking)
}

// Car image methods
func (db *Database) AddCarImage(ctx context.Context, image *models.CarImage) (int, error) {
    return db.carDB.AddCarImage(ctx, image)
}

func (db *Database) GetCarImages(ctx context.Context, carID string) ([]models.CarImage, error) {
    return db.carDB.GetCarImages(ctx, carID)
}

func (db *Database) DeleteCarImage(ctx context.Context, imageID string) (string, error) {
    return db.carDB.DeleteCarImage(ctx, imageID)
}
// Marketplace methods
func (db *Database) CreateListing(ctx context.Context, listing *models.MarketplaceListing) (int, error) {
    return db.marketplaceDB.CreateListing(ctx, listing)
}

func (db *Database) GetListings(ctx context.Context, filters map[string]string, limit int, offset int) ([]models.MarketplaceListing, int64, error) {
    return db.marketplaceDB.GetListings(ctx, filters, limit, offset)
}

func (db *Database) GetListingByID(ctx context.Context, id int) (*models.MarketplaceListing, error) {
    return db.marketplaceDB.GetListingByID(ctx, id)
}

func (db *Database) UpdateListing(ctx context.Context, listing *models.MarketplaceListing) error {
    return db.marketplaceDB.UpdateListing(ctx, listing)
}

func (db *Database) DeleteListing(ctx context.Context, id int, userID int) error {
    return db.marketplaceDB.DeleteListing(ctx, id, userID)
}

func (db *Database) GetCategories(ctx context.Context) ([]models.MarketplaceCategory, error) {
    return db.marketplaceDB.GetCategories(ctx)
}

func (db *Database) GetCategoryByID(ctx context.Context, id int) (*models.MarketplaceCategory, error) {
    return db.marketplaceDB.GetCategoryByID(ctx, id)
}

func (db *Database) GetCategoryTree(ctx context.Context) ([]models.CategoryTreeNode, error) {
    return db.marketplaceDB.GetCategoryTree(ctx)
}

func (db *Database) AddListingImage(ctx context.Context, image *models.MarketplaceImage) (int, error) {
    return db.marketplaceDB.AddListingImage(ctx, image)
}

func (db *Database) GetListingImages(ctx context.Context, listingID string) ([]models.MarketplaceImage, error) {
    return db.marketplaceDB.GetListingImages(ctx, listingID)
}

func (db *Database) DeleteListingImage(ctx context.Context, imageID string) (string, error) {
    return db.marketplaceDB.DeleteListingImage(ctx, imageID)
}

func (db *Database) AddToFavorites(ctx context.Context, userID int, listingID int) error {
    return db.marketplaceDB.AddToFavorites(ctx, userID, listingID)
}

func (db *Database) RemoveFromFavorites(ctx context.Context, userID int, listingID int) error {
    return db.marketplaceDB.RemoveFromFavorites(ctx, userID, listingID)
}

func (db *Database) GetUserFavorites(ctx context.Context, userID int) ([]models.MarketplaceListing, error) {
    return db.marketplaceDB.GetUserFavorites(ctx, userID)
}
// Добавляем делегирующие методы
func (db *Database) CreateReview(ctx context.Context, review *models.Review) (*models.Review, error) {
    return db.reviewDB.CreateReview(ctx, review)
}

func (db *Database) GetReviews(ctx context.Context, filter models.ReviewsFilter) ([]models.Review, int64, error) {
    return db.reviewDB.GetReviews(ctx, filter)
}

func (db *Database) GetReviewByID(ctx context.Context, id int) (*models.Review, error) {
    return db.reviewDB.GetReviewByID(ctx, id)
}

func (db *Database) UpdateReview(ctx context.Context, review *models.Review) error {
    return db.reviewDB.UpdateReview(ctx, review)
}

func (db *Database) DeleteReview(ctx context.Context, id int) error {
    return db.reviewDB.DeleteReview(ctx, id)
}

func (db *Database) AddReviewResponse(ctx context.Context, response *models.ReviewResponse) error {
    return db.reviewDB.AddReviewResponse(ctx, response)
}

func (db *Database) AddReviewVote(ctx context.Context, vote *models.ReviewVote) error {
    return db.reviewDB.AddReviewVote(ctx, vote)
}

func (db *Database) GetReviewVotes(ctx context.Context, reviewId int) (helpful int, notHelpful int, err error) {
    return db.reviewDB.GetReviewVotes(ctx, reviewId)
}

func (db *Database) GetUserReviewVote(ctx context.Context, userId int, reviewId int) (string, error) {
    return db.reviewDB.GetUserReviewVote(ctx, userId, reviewId)
}

func (db *Database) GetEntityRating(ctx context.Context, entityType string, entityId int) (float64, error) {
    return db.reviewDB.GetEntityRating(ctx, entityType, entityId)
}
// Методы для работы с images
func (db *Database) AddBedImage(ctx context.Context, image *models.RoomImage) (int, error) {
    return db.accommodationDB.AddBedImage(ctx, image)
}

func (db *Database) GetBedImages(ctx context.Context, bedID string) ([]models.RoomImage, error) {
    return db.accommodationDB.GetBedImages(ctx, bedID)
}

func (db *Database) AddRoomImage(ctx context.Context, image *models.RoomImage) (int, error) {
    return db.accommodationDB.AddRoomImage(ctx, image)
}

func (db *Database) GetRoomImages(ctx context.Context, roomID string) ([]models.RoomImage, error) {
    return db.accommodationDB.GetRoomImages(ctx, roomID)
}

func (db *Database) DeleteRoomImage(ctx context.Context, imageID string) (string, error) {
    return db.accommodationDB.DeleteRoomImage(ctx, imageID)
}

// Методы для работы с rooms
func (db *Database) AddRoom(ctx context.Context, room *models.Room) (int, error) {
    return db.accommodationDB.AddRoom(ctx, room)
}

func (db *Database) GetRooms(ctx context.Context, filters map[string]string, sortBy string, sortDirection string, limit int, offset int) ([]models.Room, int64, error) {
    return db.accommodationDB.GetRooms(ctx, filters, sortBy, sortDirection, limit, offset)
}

func (db *Database) GetRoomByID(ctx context.Context, id int) (*models.Room, error) {
    return db.accommodationDB.GetRoomByID(ctx, id)
}

// Методы для работы с beds
func (db *Database) AddBed(ctx context.Context, roomID int, bedNumber string, pricePerNight float64, hasOutlet bool, hasLight bool, hasShelf bool, bedType string) (int, error) {
    return db.accommodationDB.AddBed(ctx, roomID, bedNumber, pricePerNight, hasOutlet, hasLight, hasShelf, bedType)
}

func (db *Database) GetBedByID(ctx context.Context, id int) (*models.Bed, error) {
    return db.accommodationDB.GetBedByID(ctx, id)
}

func (db *Database) GetBedsByRoomID(ctx context.Context, roomID int) ([]models.Bed, error) {
    return db.accommodationDB.GetBedsByRoomID(ctx, roomID)
}

func (db *Database) GetAvailableBeds(ctx context.Context, roomID string, startDate string, endDate string) ([]models.Bed, error) {
    return db.accommodationDB.GetAvailableBeds(ctx, roomID, startDate, endDate)
}

func (db *Database) UpdateBedAvailability(ctx context.Context, bedID int, isAvailable bool) error {
    return db.accommodationDB.UpdateBedAvailability(ctx, bedID, isAvailable)
}

func (db *Database) UpdateBedPrice(ctx context.Context, bedID int, price float64) error {
    return db.accommodationDB.UpdateBedPrice(ctx, bedID, price)
}

func (db *Database) UpdateBedAttributes(ctx context.Context, bedID int, bedReq *models.BedRequest) error {
    return db.accommodationDB.UpdateBedAttributes(ctx, bedID, bedReq)
}

// Методы для работы с bookings
func (db *Database) CreateBooking(ctx context.Context, booking *models.BookingRequest) error {
    return db.accommodationDB.CreateBooking(ctx, booking)
}

func (db *Database) GetAllBookings(ctx context.Context) ([]models.Booking, error) {
    return db.accommodationDB.GetAllBookings(ctx)
}

func (db *Database) DeleteBooking(ctx context.Context, bookingID string, bookingType string) error {
    return db.accommodationDB.DeleteBooking(ctx, bookingID, bookingType)
}