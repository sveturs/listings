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
notificationStorage "backend/internal/proj/notifications/storage/postgres"
    marketplaceStorage "backend/internal/proj/marketplace/storage/postgres"
    reviewStorage "backend/internal/proj/reviews/storage/postgres"
    userStorage             "backend/internal/proj/users/storage/postgres"
)

type Database struct {
    pool *pgxpool.Pool
    marketplaceDB *marketplaceStorage.Storage
    reviewDB *reviewStorage.Storage
    usersDB *userStorage.Storage 
    notificationsDB *notificationStorage.Storage
    
}

func NewDatabase(dbURL string) (*Database, error) {
    pool, err := pgxpool.New(context.Background(), dbURL)
    if err != nil {
        return nil, fmt.Errorf("error creating connection pool: %w", err)
    }

    return &Database{
        pool:          pool,
        marketplaceDB: marketplaceStorage.NewStorage(pool),
        reviewDB:      reviewStorage.NewStorage(pool),
        usersDB:        userStorage.NewStorage(pool),
        notificationsDB: notificationStorage.NewNotificationStorage(pool),
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

// User methods
func (db *Database) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
    return db.usersDB.GetUserByEmail(ctx, email)
}

func (db *Database) GetUserByID(ctx context.Context, id int) (*models.User, error) {
    return db.usersDB.GetUserByID(ctx, id)
}

func (db *Database) CreateUser(ctx context.Context, user *models.User) error {
    return db.usersDB.CreateUser(ctx, user)
}

func (db *Database) UpdateUser(ctx context.Context, user *models.User) error {
    return db.usersDB.UpdateUser(ctx, user)
}

func (db *Database) GetOrCreateGoogleUser(ctx context.Context, user *models.User) (*models.User, error) {
    return db.usersDB.GetOrCreateGoogleUser(ctx, user)
}

func (db *Database) GetUserProfile(ctx context.Context, id int) (*models.UserProfile, error) {
    return db.usersDB.GetUserProfile(ctx, id)
}

func (db *Database) UpdateUserProfile(ctx context.Context, id int, update *models.UserProfileUpdate) error {
    return db.usersDB.UpdateUserProfile(ctx, id, update)
}

func (db *Database) UpdateLastSeen(ctx context.Context, id int) error {
    return db.usersDB.UpdateLastSeen(ctx, id)
}

// Добавить следующие методы в структуру Database:

func (db *Database) ArchiveChat(ctx context.Context, chatID int, userID int) error {
    return db.marketplaceDB.ArchiveChat(ctx, chatID, userID)
}

func (db *Database) CreateMessage(ctx context.Context, msg *models.MarketplaceMessage) error {
    return db.marketplaceDB.CreateMessage(ctx, msg)
}

func (db *Database) GetMessages(ctx context.Context, listingID int, userID int, offset int, limit int) ([]models.MarketplaceMessage, error) {
    return db.marketplaceDB.GetMessages(ctx, listingID, userID, offset, limit)
}

func (db *Database) GetChats(ctx context.Context, userID int) ([]models.MarketplaceChat, error) {
    return db.marketplaceDB.GetChats(ctx, userID)
}

func (db *Database) GetChat(ctx context.Context, chatID int, userID int) (*models.MarketplaceChat, error) {
    return db.marketplaceDB.GetChat(ctx, chatID, userID)
}

func (db *Database) MarkMessagesAsRead(ctx context.Context, messageIDs []int, userID int) error {
    return db.marketplaceDB.MarkMessagesAsRead(ctx, messageIDs, userID)
}
func (db *Database) GetUnreadMessagesCount(ctx context.Context, userID int) (int, error) {
    var count int
    err := db.pool.QueryRow(ctx, `
        SELECT COUNT(*) 
        FROM marketplace_messages 
        WHERE receiver_id = $1 AND is_read = false
    `, userID).Scan(&count)
    
    if err != nil {
        return 0, err
    }
    
    return count, nil
}
func (db *Database) CreateNotification(ctx context.Context, n *models.Notification) error {
    return db.notificationsDB.CreateNotification(ctx, n)
}

func (db *Database) GetNotificationSettings(ctx context.Context, userID int) ([]models.NotificationSettings, error) {
    return db.notificationsDB.GetNotificationSettings(ctx, userID)
}

func (db *Database) UpdateNotificationSettings(ctx context.Context, s *models.NotificationSettings) error {
    return db.notificationsDB.UpdateNotificationSettings(ctx, s)
}

func (db *Database) SaveTelegramConnection(ctx context.Context, userID int, chatID string, username string) error {
    return db.notificationsDB.SaveTelegramConnection(ctx, userID, chatID, username)
}

func (db *Database) GetTelegramConnection(ctx context.Context, userID int) (*models.TelegramConnection, error) {
    return db.notificationsDB.GetTelegramConnection(ctx, userID)
}

func (db *Database) DeleteTelegramConnection(ctx context.Context, userID int) error {
    return db.notificationsDB.DeleteTelegramConnection(ctx, userID)
}

func (db *Database) SavePushSubscription(ctx context.Context, sub *models.PushSubscription) error {
    return db.notificationsDB.SavePushSubscription(ctx, sub)
}

func (db *Database) GetPushSubscriptions(ctx context.Context, userID int) ([]models.PushSubscription, error) {
    return db.notificationsDB.GetPushSubscriptions(ctx, userID)
}

func (db *Database) DeletePushSubscription(ctx context.Context, userID int, endpoint string) error {
    return db.notificationsDB.DeletePushSubscription(ctx, userID, endpoint)
}

func (db *Database) GetUserNotifications(ctx context.Context, userID int, limit, offset int) ([]models.Notification, error) {
    return db.notificationsDB.GetUserNotifications(ctx, userID, limit, offset)
}

func (db *Database) MarkNotificationAsRead(ctx context.Context, userID int, notificationID int) error {
    return db.notificationsDB.MarkNotificationAsRead(ctx, userID, notificationID)
}

func (db *Database) DeleteNotification(ctx context.Context, userID int, notificationID int) error {
    return db.notificationsDB.DeleteNotification(ctx, userID, notificationID)
}