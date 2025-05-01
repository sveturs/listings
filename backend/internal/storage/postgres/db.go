// backend/internal/storage/postgres/db.go
package postgres

import (
	"backend/internal/domain/models"
	"backend/internal/domain/search"
	marketplaceService "backend/internal/proj/marketplace/service"
	"backend/internal/storage"
	"backend/internal/storage/filestorage"
	"backend/internal/types"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"

	marketplaceStorage "backend/internal/proj/marketplace/storage/postgres"
	notificationStorage "backend/internal/proj/notifications/storage/postgres"
	reviewStorage "backend/internal/proj/reviews/storage/postgres"
	userStorage "backend/internal/proj/users/storage/postgres"

	"backend/internal/proj/marketplace/storage/opensearch"
	osClient "backend/internal/storage/opensearch"
)

type Database struct {
	pool          *pgxpool.Pool
	marketplaceDB *marketplaceStorage.Storage

	reviewDB          *reviewStorage.Storage
	usersDB           *userStorage.Storage
	notificationsDB   *notificationStorage.Storage
	osMarketplaceRepo opensearch.MarketplaceSearchRepository
	db                *sql.DB
	marketplaceIndex  string
	fsStorage         filestorage.FileStorageInterface
}

func NewDatabase(dbURL string, osClient *osClient.OpenSearchClient, indexName string, fileStorage filestorage.FileStorageInterface) (*Database, error) {
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("error creating connection pool: %w", err)
	}

	// Создаем сервис переводов
	translationService, err := marketplaceService.NewTranslationService(os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		return nil, fmt.Errorf("error creating translation service: %w", err)
	}

	db := &Database{
		pool:            pool,
		marketplaceDB:   marketplaceStorage.NewStorage(pool, translationService),
		reviewDB:        reviewStorage.NewStorage(pool, translationService),
		usersDB:         userStorage.NewStorage(pool),
		notificationsDB: notificationStorage.NewNotificationStorage(pool),
		fsStorage:       fileStorage, // Используем переданный параметр
	}

	// Инициализируем репозиторий OpenSearch, если клиент передан
	if osClient != nil {
		db.osMarketplaceRepo = opensearch.NewRepository(osClient, indexName, db)
		// Подготавливаем индекс
		if err := db.osMarketplaceRepo.PrepareIndex(context.Background()); err != nil {
			log.Printf("Ошибка подготовки индекса OpenSearch: %v", err)
		}
	}

	return db, nil
}

var _ storage.Storage = (*Database)(nil)

func (db *Database) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
}
func (db *Database) FileStorage() filestorage.FileStorageInterface {
	return db.fsStorage
}
func (db *Database) SearchListingsOpenSearch(ctx context.Context, params *search.SearchParams) (*search.SearchResult, error) {
	if db.osMarketplaceRepo == nil {
		return nil, fmt.Errorf("OpenSearch не настроен")
	}
	return db.osMarketplaceRepo.SearchListings(ctx, params)
}
func (db *Database) GetListingImageByID(ctx context.Context, imageID int) (*models.MarketplaceImage, error) {
	var image models.MarketplaceImage

	err := db.pool.QueryRow(ctx, `
		SELECT id, listing_id, file_path, file_name, file_size, content_type, is_main,
		       storage_type, storage_bucket, public_url, created_at
		FROM marketplace_images
		WHERE id = $1
	`, imageID).Scan(
		&image.ID, &image.ListingID, &image.FilePath, &image.FileName, &image.FileSize,
		&image.ContentType, &image.IsMain, &image.StorageType, &image.StorageBucket,
		&image.PublicURL, &image.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("image not found")
		}
		return nil, err
	}

	return &image, nil
}

func (db *Database) DeleteListingImage(ctx context.Context, imageID int) error {
	_, err := db.pool.Exec(ctx, `
		DELETE FROM marketplace_images
		WHERE id = $1
	`, imageID)

	return err
}

func (db *Database) IndexListing(ctx context.Context, listing *models.MarketplaceListing) error {
	if db.osMarketplaceRepo == nil {
		return fmt.Errorf("OpenSearch не настроен")
	}

	return db.osMarketplaceRepo.IndexListing(ctx, listing)
}
func (db *Database) PrepareIndex(ctx context.Context) error {
	if db.osMarketplaceRepo == nil {
		// Если репозиторий OpenSearch не инициализирован, просто возвращаем nil
		// Поиск будет работать без OpenSearch
		return nil
	}

	// Используем уже инициализированный репозиторий для проверки индекса
	return db.osMarketplaceRepo.PrepareIndex(ctx)
}
func (db *Database) DeleteListingIndex(ctx context.Context, id string) error {
	if db.osMarketplaceRepo == nil {
		return fmt.Errorf("OpenSearch не настроен")
	}

	return db.osMarketplaceRepo.DeleteListing(ctx, id)
}

func (db *Database) SuggestListings(ctx context.Context, prefix string, size int) ([]string, error) {
	if db.osMarketplaceRepo == nil {
		return nil, fmt.Errorf("OpenSearch не настроен")
	}

	return db.osMarketplaceRepo.SuggestListings(ctx, prefix, size)
}

func (db *Database) ReindexAllListings(ctx context.Context) error {
	if db.osMarketplaceRepo == nil {
		return fmt.Errorf("OpenSearch не настроен")
	}

	return db.osMarketplaceRepo.ReindexAll(ctx)
}

// GetCategoryAttributes получает атрибуты для указанной категории
func (db *Database) GetCategoryAttributes(ctx context.Context, categoryID int) ([]models.CategoryAttribute, error) {
	return db.marketplaceDB.GetCategoryAttributes(ctx, categoryID)
}

// SaveListingAttributes сохраняет значения атрибутов для объявления
func (db *Database) SaveListingAttributes(ctx context.Context, listingID int, attributes []models.ListingAttributeValue) error {
	return db.marketplaceDB.SaveListingAttributes(ctx, listingID, attributes)
}
func (db *Database) GetAttributeRanges(ctx context.Context, categoryID int) (map[string]map[string]interface{}, error) {
	return db.marketplaceDB.GetAttributeRanges(ctx, categoryID)
}

// GetListingAttributes получает значения атрибутов для объявления
func (db *Database) GetListingAttributes(ctx context.Context, listingID int) ([]models.ListingAttributeValue, error) {
	return db.marketplaceDB.GetListingAttributes(ctx, listingID)
}
func (db *Database) GetSession(ctx context.Context, token string) (*types.SessionData, error) {
	var session types.SessionData
	err := db.pool.QueryRow(ctx, `
        SELECT user_id, name, email, google_id, picture_url, provider
        FROM sessions s
        JOIN users u ON s.user_id = u.id
        WHERE token = $1 AND expires_at > NOW()
    `, token).Scan(
		&session.UserID,
		&session.Name,
		&session.Email,
		&session.GoogleID,
		&session.PictureURL,
		&session.Provider,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

func (db *Database) GetFavoritedUsers(ctx context.Context, listingID int) ([]int, error) {
	query := `
        SELECT user_id
        FROM marketplace_favorites
        WHERE listing_id = $1
    `
	rows, err := db.pool.Query(ctx, query, listingID)
	if err != nil {
		return nil, fmt.Errorf("error querying favorited users: %w", err)
	}
	defer rows.Close()

	var userIDs []int
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return nil, fmt.Errorf("error scanning user ID: %w", err)
		}
		userIDs = append(userIDs, userID)
	}

	return userIDs, nil
}

type pgxResult struct {
	ct pgconn.CommandTag
}

func (db *Database) GetSubcategories(ctx context.Context, parentID *int, limit int, offset int) ([]models.CategoryTreeNode, error) {
	return db.marketplaceDB.GetSubcategories(ctx, parentID, limit, offset)
}
func (r pgxResult) LastInsertId() (int64, error) {
	return 0, fmt.Errorf("LastInsertId is not supported by PostgreSQL")
}

func (r pgxResult) RowsAffected() (int64, error) {
	return r.ct.RowsAffected(), nil
}

func (db *Database) Exec(ctx context.Context, sql string, args ...interface{}) (sql.Result, error) {
	ct, err := db.pool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return &pgxResult{ct: ct}, nil
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
func (r *RowsWrapper) Err() error {
	return r.rows.Err()
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

// Административные методы для управления пользователями
func (db *Database) GetAllUsers(ctx context.Context, limit, offset int) ([]*models.UserProfile, int, error) {
	return db.usersDB.GetAllUsers(ctx, limit, offset)
}

func (db *Database) UpdateUserStatus(ctx context.Context, id int, status string) error {
	return db.usersDB.UpdateUserStatus(ctx, id, status)
}

func (db *Database) DeleteUser(ctx context.Context, id int) error {
	return db.usersDB.DeleteUser(ctx, id)
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

// SearchListings выполняет поиск объявлений с пользовательским запросом
func (db *Database) SearchListings(ctx context.Context, params *search.SearchParams) (*search.SearchResult, error) {
	if db.osMarketplaceRepo == nil {
		return nil, fmt.Errorf("OpenSearch не настроен")
	}
	return db.osMarketplaceRepo.SearchListings(ctx, params)
}

// Добавить этот метод в структуру Database
func (db *Database) GetAttributeOptionTranslations(ctx context.Context, attributeName, optionValue string) (map[string]string, error) {
	query := `
        SELECT option_value, en_translation, sr_translation
        FROM attribute_option_translations
        WHERE attribute_name = $1 AND option_value = $2
    `

	var optValue, enTrans, srTrans string
	err := db.pool.QueryRow(ctx, query, attributeName, optionValue).Scan(
		&optValue, &enTrans, &srTrans,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting attribute translations: %w", err)
	}

	translations := map[string]string{
		"en": enTrans,
		"sr": srTrans,
	}

	return translations, nil
}

func (db *Database) GetTelegramConnection(ctx context.Context, userID int) (*models.TelegramConnection, error) {
	return db.notificationsDB.GetTelegramConnection(ctx, userID)
}

func (db *Database) DeleteTelegramConnection(ctx context.Context, userID int) error {
	return db.notificationsDB.DeleteTelegramConnection(ctx, userID)
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

type pgxTransaction struct {
	tx pgx.Tx
}

func (db *Database) BeginTx(ctx context.Context, opts *sql.TxOptions) (storage.Transaction, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &pgxTransaction{tx: tx}, nil
}

// Реализация методов интерфейса Transaction
func (t *pgxTransaction) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	ct, err := t.tx.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &pgxResult{ct: ct}, nil
}

func (t *pgxTransaction) Query(ctx context.Context, query string, args ...interface{}) (storage.Rows, error) {
	rows, err := t.tx.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &RowsWrapper{rows: rows}, nil
}

func (t *pgxTransaction) QueryRow(ctx context.Context, query string, args ...interface{}) storage.Row {
	return t.tx.QueryRow(ctx, query, args...)
}

func (t *pgxTransaction) Commit() error {
	return t.tx.Commit(context.Background())
}

func (t *pgxTransaction) Rollback() error {
	return t.tx.Rollback(context.Background())
}

func (db *Database) GetUserReviews(ctx context.Context, userID int, filter models.ReviewsFilter) ([]models.Review, error) {
	return db.reviewDB.GetUserReviews(ctx, userID, filter)
}

func (db *Database) GetStorefrontReviews(ctx context.Context, storefrontID int, filter models.ReviewsFilter) ([]models.Review, error) {
	return db.reviewDB.GetStorefrontReviews(ctx, storefrontID, filter)
}

func (db *Database) GetUserRatingSummary(ctx context.Context, userID int) (*models.UserRatingSummary, error) {
	return db.reviewDB.GetUserRatingSummary(ctx, userID)
}

func (db *Database) GetStorefrontRatingSummary(ctx context.Context, storefrontID int) (*models.StorefrontRatingSummary, error) {
	return db.reviewDB.GetStorefrontRatingSummary(ctx, storefrontID)
}
func (db *Database) SynchronizeDiscountMetadata(ctx context.Context) error {
	// Получаем все объявления с информацией о скидке
	query := `
        SELECT id, price, metadata
        FROM marketplace_listings
        WHERE metadata->>'discount' IS NOT NULL
    `

	rows, err := db.pool.Query(ctx, query)
	if err != nil {
		return fmt.Errorf("error querying listings with discounts: %w", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id int
		var price float64
		var metadataJSON []byte

		if err := rows.Scan(&id, &price, &metadataJSON); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		var metadata map[string]interface{}
		if err := json.Unmarshal(metadataJSON, &metadata); err != nil {
			log.Printf("Error unmarshaling metadata: %v", err)
			continue
		}

		// Проверяем и обновляем информацию о скидке
		if discount, ok := metadata["discount"].(map[string]interface{}); ok {
			if prevPrice, ok := discount["previous_price"].(float64); ok && prevPrice > 0 {
				// Пересчитываем актуальный процент скидки
				if prevPrice > price {
					discountPercent := int((prevPrice - price) / prevPrice * 100)
					discount["discount_percent"] = float64(discountPercent)

					// Обновляем метаданные в БД
					metadata["discount"] = discount
					updatedMetadataJSON, err := json.Marshal(metadata)
					if err != nil {
						log.Printf("Error marshaling updated metadata: %v", err)
						continue
					}

					_, err = db.pool.Exec(ctx, `
                        UPDATE marketplace_listings
                        SET metadata = $1
                        WHERE id = $2
                    `, updatedMetadataJSON, id)

					if err != nil {
						log.Printf("Error updating metadata for listing %d: %v", id, err)
						continue
					}

					count++
					log.Printf("Updated discount percentage for listing %d: %d%%", id, discountPercent)
				}
			}
		}
	}

	log.Printf("Synchronized discount metadata for %d listings", count)
	return nil
}
