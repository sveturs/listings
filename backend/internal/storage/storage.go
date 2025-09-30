// backend/internal/storage/storage.go
package storage

import (
	"context"
	"database/sql"

	"backend/internal/domain/models"
	"backend/internal/domain/search"
	"backend/internal/proj/storefronts/storage/opensearch"
	"backend/internal/storage/filestorage"
)

type Storage interface {
	// User methods
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	GetFavoritedUsers(ctx context.Context, listingID int) ([]int, error)

	// Методы для работы с администраторами
	IsUserAdmin(ctx context.Context, email string) (bool, error)
	GetAllAdmins(ctx context.Context) ([]*models.AdminUser, error)
	AddAdmin(ctx context.Context, admin *models.AdminUser) error
	RemoveAdmin(ctx context.Context, email string) error

	// Reviews
	CreateReview(ctx context.Context, review *models.Review) (*models.Review, error)
	GetReviews(ctx context.Context, filter models.ReviewsFilter) ([]models.Review, int64, error)
	GetReviewByID(ctx context.Context, id int) (*models.Review, error)
	UpdateReview(ctx context.Context, review *models.Review) error
	UpdateReviewStatus(ctx context.Context, reviewId int, status string) error
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

	// Marketplace listing variants methods
	CreateListingVariants(ctx context.Context, listingID int, variants []models.MarketplaceListingVariant) error
	GetListingVariants(ctx context.Context, listingID int) ([]models.MarketplaceListingVariant, error)
	UpdateListingVariant(ctx context.Context, variant *models.MarketplaceListingVariant) error
	DeleteListingVariant(ctx context.Context, variantID int) error
	GetListingBySlug(ctx context.Context, slug string) (*models.MarketplaceListing, error)
	IsSlugUnique(ctx context.Context, slug string, excludeID int) (bool, error)
	GenerateUniqueSlug(ctx context.Context, baseSlug string, excludeID int) (string, error)
	IncrementViewsCount(ctx context.Context, id int) error
	UpdateListing(ctx context.Context, listing *models.MarketplaceListing) error
	DeleteListing(ctx context.Context, id int, userID int) error
	DeleteListingAdmin(ctx context.Context, id int) error
	GetCategoryTree(ctx context.Context) ([]models.CategoryTreeNode, error)
	AddListingImage(ctx context.Context, image *models.MarketplaceImage) (int, error)
	GetListingImages(ctx context.Context, listingID string) ([]models.MarketplaceImage, error)

	// FileStorage возвращает интерфейс для работы с файловым хранилищем
	FileStorage() filestorage.FileStorageInterface

	// GetListingImageByID возвращает информацию об изображении по ID
	GetListingImageByID(ctx context.Context, imageID int) (*models.MarketplaceImage, error)

	// DeleteListingImage удаляет информацию об изображении из базы данных
	DeleteListingImage(ctx context.Context, imageID int) error

	GetAttributeOptionTranslations(ctx context.Context, attributeName, optionValue string) (map[string]string, error)
	GetAttributeRanges(ctx context.Context, categoryID int) (map[string]map[string]interface{}, error)

	GetCategoryAttributes(ctx context.Context, categoryID int) ([]models.CategoryAttribute, error)
	GetCategories(ctx context.Context) ([]models.MarketplaceCategory, error)
	GetAllCategories(ctx context.Context) ([]models.MarketplaceCategory, error)
	GetPopularCategories(ctx context.Context, limit int) ([]models.MarketplaceCategory, error)
	GetCategoryByID(ctx context.Context, id int) (*models.MarketplaceCategory, error)

	AddToFavorites(ctx context.Context, userID int, listingID int) error
	RemoveFromFavorites(ctx context.Context, userID int, listingID int) error
	GetUserFavorites(ctx context.Context, userID int) ([]models.MarketplaceListing, error)

	// Storefront favorites
	AddStorefrontToFavorites(ctx context.Context, userID int, productID int) error
	RemoveStorefrontFromFavorites(ctx context.Context, userID int, productID int) error
	GetUserStorefrontFavorites(ctx context.Context, userID int) ([]models.MarketplaceListing, error)

	GetPriceHistory(ctx context.Context, listingID int) ([]models.PriceHistoryEntry, error)
	AddPriceHistoryEntry(ctx context.Context, entry *models.PriceHistoryEntry) error
	ClosePriceHistoryEntry(ctx context.Context, listingID int) error
	CheckPriceManipulation(ctx context.Context, listingID int) (bool, error)

	SaveListingAttributes(ctx context.Context, listingID int, attributes []models.ListingAttributeValue) error
	GetListingAttributes(ctx context.Context, listingID int) ([]models.ListingAttributeValue, error)
	SynchronizeDiscountMetadata(ctx context.Context) error // Добавьте эту строку

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

	// Chat attachments methods
	CreateChatAttachment(ctx context.Context, attachment *models.ChatAttachment) error
	GetChatAttachment(ctx context.Context, attachmentID int) (*models.ChatAttachment, error)
	GetMessageAttachments(ctx context.Context, messageID int) ([]*models.ChatAttachment, error)
	DeleteChatAttachment(ctx context.Context, attachmentID int) error
	UpdateMessageAttachmentsCount(ctx context.Context, messageID int, count int) error
	GetMessageByID(ctx context.Context, messageID int) (*models.MarketplaceMessage, error)

	// Chat verification methods
	GetChatActivityStats(ctx context.Context, buyerID int, sellerID int, listingID int) (*models.ChatActivityStats, error)

	// Rating aggregation methods
	GetUserAggregatedRating(ctx context.Context, userID int) (*models.UserAggregatedRating, error)
	GetStorefrontAggregatedRating(ctx context.Context, storefrontID int) (*models.StorefrontAggregatedRating, error)
	RefreshRatingViews(ctx context.Context) error

	// Review confirmation and dispute methods
	CreateReviewConfirmation(ctx context.Context, confirmation *models.ReviewConfirmation) error
	GetReviewConfirmation(ctx context.Context, reviewID int) (*models.ReviewConfirmation, error)
	CreateReviewDispute(ctx context.Context, dispute *models.ReviewDispute) error
	GetReviewDispute(ctx context.Context, reviewID int) (*models.ReviewDispute, error)
	UpdateReviewDispute(ctx context.Context, dispute *models.ReviewDispute) error

	// Review permission check
	CanUserReviewEntity(ctx context.Context, userID int, entityType string, entityID int) (*models.CanReviewResponse, error)

	Exec(ctx context.Context, sql string, args ...interface{}) (sql.Result, error)

	// Storefront methods
	CreateStorefront(ctx context.Context, userID int, dto *models.StorefrontCreateDTO) (*models.Storefront, error)
	GetUserStorefronts(ctx context.Context, userID int) ([]models.Storefront, error)
	GetStorefrontByID(ctx context.Context, id int) (*models.Storefront, error)
	UpdateStorefront(ctx context.Context, storefront *models.Storefront) error
	DeleteStorefront(ctx context.Context, id int) error
	GetStorefrontOwnerByProductID(ctx context.Context, productID int) (int, error)

	// Storefront repository access
	Storefront() interface{}

	// Orders system repository access
	Cart() interface{}
	Order() interface{}
	Inventory() interface{}

	// Marketplace orders repository access
	MarketplaceOrder() interface{}

	// Storefront product search repository access
	StorefrontProductSearch() interface{}

	// OpenSearch методы
	SearchListings(ctx context.Context, params *search.SearchParams) (*search.SearchResult, error)

	SearchListingsOpenSearch(ctx context.Context, params *search.SearchParams) (*search.SearchResult, error)
	SuggestListings(ctx context.Context, prefix string, size int) ([]string, error)
	ReindexAllListings(ctx context.Context) error
	IndexListing(ctx context.Context, listing *models.MarketplaceListing) error
	DeleteListingIndex(ctx context.Context, id string) error
	PrepareIndex(ctx context.Context) error

	// Storefront OpenSearch методы
	SearchStorefrontsOpenSearch(ctx context.Context, params *opensearch.StorefrontSearchParams) (*opensearch.StorefrontSearchResult, error)
	IndexStorefront(ctx context.Context, storefront *models.Storefront) error
	DeleteStorefrontIndex(ctx context.Context, storefrontID int) error
	ReindexAllStorefronts(ctx context.Context) error

	// Translation methods
	GetTranslationsForEntity(ctx context.Context, entityType string, entityID int) ([]models.Translation, error)

	// Search query methods
	GetPopularSearchQueries(ctx context.Context, query string, limit int) ([]interface{}, error)
	SaveSearchQuery(ctx context.Context, query, normalizedQuery string, resultsCount int, language string) error
	SearchCategories(ctx context.Context, query string, limit int) ([]models.MarketplaceCategory, error)

	// Fuzzy search methods
	ExpandSearchQuery(ctx context.Context, query string, language string) (string, error)
	SearchCategoriesFuzzy(ctx context.Context, searchTerm string, language string, similarityThreshold float64) ([]interface{}, error)

	// User Contacts methods
	AddContact(ctx context.Context, contact *models.UserContact) error
	UpdateContactStatus(ctx context.Context, userID, contactUserID int, status, notes string) error
	GetContact(ctx context.Context, userID, contactUserID int) (*models.UserContact, error)
	GetUserContacts(ctx context.Context, userID int, status string, page, limit int) ([]models.UserContact, int, error)
	GetIncomingContactRequests(ctx context.Context, userID int, page, limit int) ([]models.UserContact, int, error)
	RemoveContact(ctx context.Context, userID, contactUserID int) error
	GetUserPrivacySettings(ctx context.Context, userID int) (*models.UserPrivacySettings, error)
	UpdateUserPrivacySettings(ctx context.Context, userID int, settings *models.UpdatePrivacySettingsRequest) error
	CanAddContact(ctx context.Context, userID, targetUserID int) (bool, error)

	// Privacy Settings methods
	GetPrivacySettings(ctx context.Context, userID int) (*models.UserPrivacySettings, error)
	UpdatePrivacySettings(ctx context.Context, userID int, settings *models.UpdatePrivacySettingsRequest) error

	// Car Makes and Models methods
	GetCarMakes(ctx context.Context, country string, isDomestic bool, isMotorcycle bool, activeOnly bool) ([]models.CarMake, error)
	GetCarMakeBySlug(ctx context.Context, slug string) (*models.CarMake, error)
	GetCarModelsByMake(ctx context.Context, makeSlug string, activeOnly bool) ([]models.CarModel, error)
	GetCarGenerationsByModel(ctx context.Context, modelID int, activeOnly bool) ([]models.CarGeneration, error)
	SearchCarMakes(ctx context.Context, query string, limit int) ([]models.CarMake, error)
	GetCarListingsCount(ctx context.Context) (int, error)
	GetTotalCarModelsCount(ctx context.Context) (int, error)

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
	Err() error // Добавляем метод Err() для проверки ошибок после итерации
}
