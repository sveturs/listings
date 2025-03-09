// backend/internal/storage/postgres/db.go
package postgres

import (
	"backend/internal/domain/models"
	"backend/internal/domain/search"
	marketplaceService "backend/internal/proj/marketplace/service"
	"backend/internal/storage"
	osClient "backend/internal/storage/opensearch"
	"backend/internal/types"
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"time"

	marketplaceStorage "backend/internal/proj/marketplace/storage/postgres"
	notificationStorage "backend/internal/proj/notifications/storage/postgres"
	reviewStorage "backend/internal/proj/reviews/storage/postgres"
	userStorage "backend/internal/proj/users/storage/postgres"

	"backend/internal/proj/marketplace/storage/opensearch"
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
}

func NewDatabase(dbURL string, osClient *osClient.OpenSearchClient, indexName string) (*Database, error) {
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
func (db *Database) SearchListingsOpenSearch(ctx context.Context, params *search.SearchParams) (*search.SearchResult, error) {
	if db.osMarketplaceRepo == nil {
		return nil, fmt.Errorf("OpenSearch не настроен")
	}

	return db.osMarketplaceRepo.SearchListings(ctx, params)
}
func (db *Database) IndexListing(ctx context.Context, listing *models.MarketplaceListing) error {
    if db.osMarketplaceRepo == nil {
        return fmt.Errorf("OpenSearch не настроен")
    }

    // Проверяем, является ли это автомобильным объявлением
    isAuto := false
    
    // Список ID автомобильных категорий
    autoCategories := []int{2000, 2100, 2200, 2210, 2220, 2230, 2240, 2300, 2310, 2315, 2320, 2325, 2330, 2335, 2340, 2345, 2350, 2355, 2360, 2365}
    
    // Проверяем напрямую по списку
    for _, id := range autoCategories {
        if id == listing.CategoryID {
            isAuto = true
            break
        }
    }
    
    // Если не нашли в прямом списке, проверяем по дереву категорий
    if !isAuto {
        row := db.pool.QueryRow(ctx, `
            WITH RECURSIVE category_tree AS (
                -- Базовый случай: категория 2000 (Автомобили)
                SELECT id, parent_id FROM marketplace_categories WHERE id = 2000
                UNION ALL
                -- Рекурсивный случай: все дочерние категории
                SELECT c.id, c.parent_id
                FROM marketplace_categories c
                JOIN category_tree ct ON c.parent_id = ct.id
            )
            SELECT EXISTS(SELECT 1 FROM category_tree WHERE id = $1)
        `, listing.CategoryID)

        err := row.Scan(&isAuto)
        if err != nil {
            log.Printf("Ошибка при проверке автомобильной категории: %v", err)
            // Продолжаем стандартную индексацию
        }
    }

    // Если это автомобильное объявление, добавляем автомобильные свойства
    if isAuto {
        autoProps, err := db.GetAutoPropertiesByListingID(ctx, listing.ID)
        if err == nil && autoProps != nil {
            // Создаем расширенный документ с автосвойствами
            doc := map[string]interface{}{
                "id":                listing.ID,
                "title":             listing.Title,
                "description":       listing.Description,
                "price":             listing.Price,
                "condition":         listing.Condition,
                "status":            listing.Status,
                "category_id":       listing.CategoryID,
                "user_id":           listing.UserID,
                "original_language": listing.OriginalLanguage,
                "created_at":        listing.CreatedAt.Format(time.RFC3339),
                "updated_at":        listing.UpdatedAt.Format(time.RFC3339),

                // Автомобильные свойства
                "auto_brand":           autoProps.Brand,
                "auto_model":           autoProps.Model,
                "auto_year":            autoProps.Year,
                "auto_mileage":         autoProps.Mileage,
                "auto_fuel_type":       autoProps.FuelType,
                "auto_transmission":    autoProps.Transmission,
                "auto_body_type":       autoProps.BodyType,
                "auto_drive_type":      autoProps.DriveType,
                "auto_engine_capacity": autoProps.EngineCapacity,
                "auto_power":           autoProps.Power,
                "auto_color":           autoProps.Color,
                "auto_number_of_doors": autoProps.NumberOfDoors,
                "auto_number_of_seats": autoProps.NumberOfSeats,
            }

            // Индексируем обогащенный документ напрямую через клиент OpenSearch
            if db.osMarketplaceRepo != nil {
                if client, err := db.getOpenSearchClient(); err == nil && client != nil {
                    // Если можем получить клиент напрямую
                    err = client.IndexDocument(db.marketplaceIndex, fmt.Sprintf("%d", listing.ID), doc)
                    if err != nil {
                        log.Printf("Ошибка при индексации автомобильного объявления напрямую: %v", err)
                    } else {
                        // Успешно индексировали напрямую
                        return nil
                    }
                }
            }
        }
    }

    // Стандартная индексация через репозиторий
    return db.osMarketplaceRepo.IndexListing(ctx, listing)
}

// Вспомогательный метод для получения клиента OpenSearch
func (db *Database) getOpenSearchClient() (*osClient.OpenSearchClient, error) {
	// В данном случае у нас нет прямого доступа к клиенту
	return nil, fmt.Errorf("not implemented")
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
