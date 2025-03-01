// backend/internal/storage/storage.go
package storage

import (
	"backend/internal/domain/models"
	"backend/internal/types"
	"context"
	"database/sql"
)

type Storage interface {
	// User methods
	GetOrCreateGoogleUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, user *models.User) error
	GetUserProfile(ctx context.Context, id int) (*models.UserProfile, error)
	UpdateUserProfile(ctx context.Context, id int, update *models.UserProfileUpdate) error
	UpdateLastSeen(ctx context.Context, id int) error
	GetFavoritedUsers(ctx context.Context, listingID int) ([]int, error)
	GetSession(ctx context.Context, token string) (*types.SessionData, error)

	// Reviews
	CreateReview(ctx context.Context, review *models.Review) (*models.Review, error)
	GetReviews(ctx context.Context, filter models.ReviewsFilter) ([]models.Review, int64, error)
	GetReviewByID(ctx context.Context, id int) (*models.Review, error)
	UpdateReview(ctx context.Context, review *models.Review) error
	DeleteReview(ctx context.Context, id int) error
	AddReviewResponse(ctx context.Context, response *models.ReviewResponse) error
	AddReviewVote(ctx context.Context, vote *models.ReviewVote) error
	GetReviewVotes(ctx context.Context, reviewId int) (helpful int, notHelpful int, err error)
	GetUserReviewVote(ctx context.Context, userId int, reviewId int) (string, error)
	GetEntityRating(ctx context.Context, entityType string, entityId int) (float64, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) Row
	Query(ctx context.Context, sql string, args ...interface{}) (Rows, error)

    GetUserReviews(ctx context.Context, userID int, filter models.ReviewsFilter) ([]models.Review, error)
    GetStorefrontReviews(ctx context.Context, storefrontID int, filter models.ReviewsFilter) ([]models.Review, error)
    GetUserRatingSummary(ctx context.Context, userID int) (*models.UserRatingSummary, error)
    GetStorefrontRatingSummary(ctx context.Context, storefrontID int) (*models.StorefrontRatingSummary, error)


	// Notification methods
	GetNotificationSettings(ctx context.Context, userID int) ([]models.NotificationSettings, error)
	UpdateNotificationSettings(ctx context.Context, settings *models.NotificationSettings) error
	SaveTelegramConnection(ctx context.Context, userID int, chatID string, username string) error
	GetTelegramConnection(ctx context.Context, userID int) (*models.TelegramConnection, error)
	DeleteTelegramConnection(ctx context.Context, userID int) error
	CreateNotification(ctx context.Context, notification *models.Notification) error
	GetUserNotifications(ctx context.Context, userID int, limit, offset int) ([]models.Notification, error)
	MarkNotificationAsRead(ctx context.Context, userID int, notificationID int) error
	DeleteNotification(ctx context.Context, userID int, notificationID int) error

	// Marketplace methods
	CreateListing(ctx context.Context, listing *models.MarketplaceListing) (int, error)
	GetListings(ctx context.Context, filters map[string]string, limit int, offset int) ([]models.MarketplaceListing, int64, error)
	GetListingByID(ctx context.Context, id int) (*models.MarketplaceListing, error)
	UpdateListing(ctx context.Context, listing *models.MarketplaceListing) error
	DeleteListing(ctx context.Context, id int, userID int) error
	GetCategoryTree(ctx context.Context) ([]models.CategoryTreeNode, error)
	GetSubcategories(ctx context.Context, parentID *int, limit int, offset int) ([]models.CategoryTreeNode, error)
	AddListingImage(ctx context.Context, image *models.MarketplaceImage) (int, error)
	GetListingImages(ctx context.Context, listingID string) ([]models.MarketplaceImage, error)
	DeleteListingImage(ctx context.Context, imageID string) (string, error)

	GetCategories(ctx context.Context) ([]models.MarketplaceCategory, error)
	GetCategoryByID(ctx context.Context, id int) (*models.MarketplaceCategory, error)

	AddToFavorites(ctx context.Context, userID int, listingID int) error
	RemoveFromFavorites(ctx context.Context, userID int, listingID int) error
	GetUserFavorites(ctx context.Context, userID int) ([]models.MarketplaceListing, error)

	// Balance methods
	GetUserBalance(ctx context.Context, userID int) (*models.UserBalance, error)
	GetUserTransactions(ctx context.Context, userID int, limit, offset int) ([]models.BalanceTransaction, error)
	CreateTransaction(ctx context.Context, transaction *models.BalanceTransaction) (int, error)
	GetActivePaymentMethods(ctx context.Context) ([]models.PaymentMethod, error)
	UpdateBalance(ctx context.Context, userID int, amount float64) error
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Transaction, error)
	// Marketplace Chat methods
	CreateMessage(ctx context.Context, msg *models.MarketplaceMessage) error
	GetMessages(ctx context.Context, listingID int, userID int, offset int, limit int) ([]models.MarketplaceMessage, error)
	GetChats(ctx context.Context, userID int) ([]models.MarketplaceChat, error)
	GetChat(ctx context.Context, chatID int, userID int) (*models.MarketplaceChat, error)
	MarkMessagesAsRead(ctx context.Context, messageIDs []int, userID int) error
	ArchiveChat(ctx context.Context, chatID int, userID int) error
	GetUnreadMessagesCount(ctx context.Context, userID int) (int, error)

	Exec(ctx context.Context, sql string, args ...interface{}) (sql.Result, error)

	// Storefront methods
	CreateStorefront(ctx context.Context, storefront *models.Storefront) (int, error)
	GetUserStorefronts(ctx context.Context, userID int) ([]models.Storefront, error)
	GetStorefrontByID(ctx context.Context, id int) (*models.Storefront, error)
	UpdateStorefront(ctx context.Context, storefront *models.Storefront) error
	DeleteStorefront(ctx context.Context, id int) error

	// Import Source methods
	CreateImportSource(ctx context.Context, source *models.ImportSource) (int, error)
	GetImportSourceByID(ctx context.Context, id int) (*models.ImportSource, error)
	GetImportSources(ctx context.Context, storefrontID int) ([]models.ImportSource, error)
	UpdateImportSource(ctx context.Context, source *models.ImportSource) error
	DeleteImportSource(ctx context.Context, id int) error

	// Import History methods
	CreateImportHistory(ctx context.Context, history *models.ImportHistory) (int, error)
	GetImportHistory(ctx context.Context, sourceID int, limit, offset int) ([]models.ImportHistory, error)
	UpdateImportHistory(ctx context.Context, history *models.ImportHistory) error
	
	// Database connection
	Close()
	Ping(ctx context.Context) error
}
type Transaction interface {
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Query(ctx context.Context, query string, args ...interface{}) (Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) Row
	Commit() error
	Rollback() error
}
type Row interface {
	Scan(dest ...interface{}) error
}

type Rows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close() error
}
