// backend/internal/storage/postgres/db.go
package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"backend/internal/config"
	"backend/internal/domain/models"
	"backend/internal/domain/search"
	marketplaceService "backend/internal/proj/marketplace/service"
	"backend/internal/storage"
	"backend/internal/storage/filestorage"
	"backend/internal/types"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	authservice "github.com/sveturs/auth/pkg/http/service"

	marketplaceStorage "backend/internal/proj/marketplace/storage/postgres"
	notificationStorage "backend/internal/proj/notifications/storage/postgres"
	reviewStorage "backend/internal/proj/reviews/storage/postgres"

	"backend/internal/proj/marketplace/storage/opensearch"
	storefrontOpenSearch "backend/internal/proj/storefronts/storage/opensearch"
	osClient "backend/internal/storage/opensearch"
)

// ErrSessionNotFound возвращается когда сессия не найдена
var ErrSessionNotFound = errors.New("session not found")

// ErrAttributeTranslationNotFound возвращается когда перевод атрибута не найден
var ErrAttributeTranslationNotFound = errors.New("attribute translation not found")

// ErrReviewConfirmationNotFound возвращается когда подтверждение отзыва не найдено
var ErrReviewConfirmationNotFound = errors.New("review confirmation not found")

// ErrReviewDisputeNotFound возвращается когда спор по отзыву не найден
var ErrReviewDisputeNotFound = errors.New("review dispute not found")

// ErrStorefrontNotFound возвращается когда витрина не найдена
var ErrStorefrontNotFound = errors.New("storefront not found")

type Database struct {
	pool          *pgxpool.Pool
	marketplaceDB *marketplaceStorage.Storage

	reviewDB             *reviewStorage.Storage
	notificationsDB      *notificationStorage.Storage
	osMarketplaceRepo    opensearch.MarketplaceSearchRepository
	osStorefrontRepo     storefrontOpenSearch.StorefrontSearchRepository
	osClient             *osClient.OpenSearchClient // Клиент OpenSearch для прямых запросов
	db                   *sql.DB
	sqlxDB               *sqlx.DB // sqlx.DB для работы с sqlx библиотекой
	marketplaceIndex     string
	storefrontIndex      string
	attributeGroups      AttributeGroupStorage
	fsStorage            filestorage.FileStorageInterface
	storefrontRepo       StorefrontRepository                         // Репозиторий для витрин
	cartRepo             CartRepositoryInterface                      // Репозиторий для корзин
	orderRepo            OrderRepositoryInterface                     // Репозиторий для заказов
	searchWeights        *config.SearchWeights                        // Веса для поиска
	inventoryRepo        InventoryRepositoryInterface                 // Репозиторий для инвентаря
	marketplaceOrderRepo *MarketplaceOrderRepository                  // Репозиторий для заказов маркетплейса
	productSearchRepo    storefrontOpenSearch.ProductSearchRepository // Репозиторий для поиска товаров витрин
}

func NewDatabase(ctx context.Context, dbURL string, osClient *osClient.OpenSearchClient, indexName string, fileStorage filestorage.FileStorageInterface, searchWeights *config.SearchWeights) (*Database, error) {
	// Настраиваем конфигурацию пула
	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing connection string: %w", err)
	}

	// Увеличиваем количество соединений для behavior tracking
	poolConfig.MaxConns = 50 // Увеличиваем максимальное количество соединений
	poolConfig.MinConns = 10 // Минимальное количество соединений

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating connection pool: %w", err)
	}

	// Создаем sql.DB из pgxpool для совместимости с sqlx
	stdDB := stdlib.OpenDBFromPool(pool)
	sqlxDB := sqlx.NewDb(stdDB, "pgx")

	// Создаем сервис переводов
	translationService, err := marketplaceService.NewTranslationService(os.Getenv("OPENAI_API_KEY"))
	if err != nil {
		return nil, fmt.Errorf("error creating translation service: %w", err)
	}

	db := &Database{
		pool:             pool,
		db:               stdDB,
		sqlxDB:           sqlxDB,
		marketplaceDB:    marketplaceStorage.NewStorage(pool, translationService, nil), // userService будет установлен позже
		reviewDB:         reviewStorage.NewStorage(pool, translationService),
		notificationsDB:  notificationStorage.NewNotificationStorage(pool),
		osClient:         osClient,      // Сохраняем клиент OpenSearch
		marketplaceIndex: indexName,     // Сохраняем имя индекса
		storefrontIndex:  "storefronts", // Индекс для витрин
		fsStorage:        fileStorage,   // Используем переданный параметр
		attributeGroups:  NewAttributeGroupStorage(pool),
	}

	// Инициализируем репозиторий витрин
	db.storefrontRepo = NewStorefrontRepository(db)

	// Инициализируем репозиторий корзин
	db.cartRepo = NewCartRepository(pool)

	// Инициализируем репозиторий заказов
	db.orderRepo = NewOrderRepository(pool)

	// Инициализируем репозиторий инвентаря
	db.inventoryRepo = NewInventoryRepository(pool)

	// Инициализируем репозиторий заказов маркетплейса
	db.marketplaceOrderRepo = NewMarketplaceOrderRepository(pool)

	// Инициализируем репозиторий OpenSearch, если клиент передан
	if osClient != nil {
		db.osMarketplaceRepo = opensearch.NewRepository(osClient, indexName, db, searchWeights)
		// Подготавливаем индекс
		if err := db.osMarketplaceRepo.PrepareIndex(ctx); err != nil {
			log.Printf("Ошибка подготовки индекса OpenSearch: %v", err)
		}

		// Инициализируем репозиторий витрин в OpenSearch
		db.osStorefrontRepo = storefrontOpenSearch.NewStorefrontRepository(osClient, db.storefrontIndex)
		// Подготавливаем индекс витрин
		if err := db.osStorefrontRepo.PrepareIndex(ctx); err != nil {
			log.Printf("Ошибка подготовки индекса витрин в OpenSearch: %v", err)
		}

		// Инициализируем репозиторий товаров витрин в OpenSearch
		// ВАЖНО: Используем тот же индекс что и для marketplace (унифицированный поиск)
		db.productSearchRepo = storefrontOpenSearch.NewProductRepository(osClient, "marketplace_listings")
		// Подготавливаем индекс товаров витрин
		if err := db.productSearchRepo.PrepareIndex(ctx); err != nil {
			log.Printf("Ошибка подготовки индекса товаров витрин в OpenSearch: %v", err)
		}
	}

	// Сохраняем search weights
	db.searchWeights = searchWeights

	return db, nil
}

var _ storage.Storage = (*Database)(nil)

func (db *Database) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
	if db.db != nil {
		if err := db.db.Close(); err != nil {
			// Логируем ошибку, но не прерываем выполнение закрытия
			_ = err // Explicitly ignore error
		}
	}
}

func (db *Database) GetSearchWeights() *config.SearchWeights {
	return db.searchWeights
}

// GetSQLXDB возвращает sqlx.DB для использования в модулях, которые требуют sqlx
func (db *Database) GetSQLXDB() *sqlx.DB {
	// Если sqlxDB уже инициализирован, возвращаем его
	if db.sqlxDB != nil {
		return db.sqlxDB
	}
	// Создаем новый sqlx.DB из пула и СОХРАНЯЕМ его
	stdDB := stdlib.OpenDBFromPool(db.pool)
	db.sqlxDB = sqlx.NewDb(stdDB, "pgx")
	return db.sqlxDB
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

// GetOpenSearchClient возвращает клиент OpenSearch для прямого выполнения запросов
func (db *Database) GetOpenSearchClient() (interface {
	Execute(ctx context.Context, method, path string, body []byte) ([]byte, error)
}, error,
) {
	if db.osClient == nil {
		return nil, fmt.Errorf("OpenSearch клиент не настроен")
	}
	return db.osClient, nil
}

// GetOpenSearchRepository возвращает репозиторий OpenSearch для маркетплейса
func (db *Database) GetOpenSearchRepository() opensearch.MarketplaceSearchRepository {
	return db.osMarketplaceRepo
}

// SearchStorefrontsOpenSearch выполняет поиск витрин через OpenSearch
func (db *Database) SearchStorefrontsOpenSearch(ctx context.Context, params *storefrontOpenSearch.StorefrontSearchParams) (*storefrontOpenSearch.StorefrontSearchResult, error) {
	if db.osStorefrontRepo == nil {
		return nil, fmt.Errorf("OpenSearch для витрин не настроен")
	}
	return db.osStorefrontRepo.Search(ctx, params)
}

// IndexStorefront индексирует витрину в OpenSearch
func (db *Database) IndexStorefront(ctx context.Context, storefront *models.Storefront) error {
	if db.osStorefrontRepo == nil {
		return fmt.Errorf("OpenSearch для витрин не настроен")
	}
	return db.osStorefrontRepo.Index(ctx, storefront)
}

// DeleteStorefrontIndex удаляет витрину из индекса OpenSearch
func (db *Database) DeleteStorefrontIndex(ctx context.Context, storefrontID int) error {
	if db.osStorefrontRepo == nil {
		return fmt.Errorf("OpenSearch для витрин не настроен")
	}
	return db.osStorefrontRepo.Delete(ctx, storefrontID)
}

// ReindexAllStorefronts переиндексирует все витрины
func (db *Database) ReindexAllStorefronts(ctx context.Context) error {
	if db.osStorefrontRepo == nil {
		return fmt.Errorf("OpenSearch для витрин не настроен")
	}
	return db.osStorefrontRepo.ReindexAll(ctx)
}

func (db *Database) GetListingImageByID(ctx context.Context, imageID int) (*models.MarketplaceImage, error) {
	var image models.MarketplaceImage
	var storageBucket, publicURL sql.NullString

	err := db.pool.QueryRow(ctx, `
		SELECT id, listing_id, file_path, file_name, file_size, content_type, is_main,
		       storage_type, storage_bucket, public_url, created_at
		FROM marketplace_images
		WHERE id = $1
	`, imageID).Scan(
		&image.ID, &image.ListingID, &image.FilePath, &image.FileName, &image.FileSize,
		&image.ContentType, &image.IsMain, &image.StorageType, &storageBucket,
		&publicURL, &image.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("image not found")
		}
		return nil, err
	}

	// Handle nullable fields
	if storageBucket.Valid {
		image.StorageBucket = storageBucket.String
	}
	if publicURL.Valid {
		image.PublicURL = publicURL.String
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

// GetSession - DEPRECATED: Sessions are now managed via JWT tokens in auth-service
// This method is kept for backward compatibility but should not be used in new code
func (db *Database) GetSession(ctx context.Context, token string) (*types.SessionData, error) {
	return nil, fmt.Errorf("GetSession: moved to JWT-based auth, sessions table no longer used")
}

// Refresh Token methods - DEPRECATED: moved to auth-service
func (db *Database) CreateRefreshToken(ctx context.Context, token *models.RefreshToken) error {
	return fmt.Errorf("CreateRefreshToken: moved to auth-service")
}

func (db *Database) GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	return nil, fmt.Errorf("GetRefreshToken: moved to auth-service")
}

func (db *Database) GetRefreshTokenByID(ctx context.Context, id int) (*models.RefreshToken, error) {
	return nil, fmt.Errorf("GetRefreshTokenByID: moved to auth-service")
}

func (db *Database) GetUserRefreshTokens(ctx context.Context, userID int) ([]*models.RefreshToken, error) {
	return nil, fmt.Errorf("GetUserRefreshTokens: moved to auth-service")
}

func (db *Database) UpdateRefreshToken(ctx context.Context, token *models.RefreshToken) error {
	return fmt.Errorf("UpdateRefreshToken: moved to auth-service")
}

func (db *Database) RevokeRefreshToken(ctx context.Context, tokenID int) error {
	return fmt.Errorf("RevokeRefreshToken: moved to auth-service")
}

func (db *Database) RevokeRefreshTokenByValue(ctx context.Context, tokenValue string) error {
	return fmt.Errorf("RevokeRefreshTokenByValue: moved to auth-service")
}

func (db *Database) RevokeUserRefreshTokens(ctx context.Context, userID int) error {
	return fmt.Errorf("RevokeUserRefreshTokens: moved to auth-service")
}

func (db *Database) DeleteExpiredRefreshTokens(ctx context.Context) (int64, error) {
	return 0, fmt.Errorf("DeleteExpiredRefreshTokens: moved to auth-service")
}

// GetSQLDB returns the raw sql.DB connection
func (db *Database) GetSQLDB() *sql.DB {
	return db.db
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

func (db *Database) GetListingBySlug(ctx context.Context, slug string) (*models.MarketplaceListing, error) {
	return db.marketplaceDB.GetListingBySlug(ctx, slug)
}

func (db *Database) IsSlugUnique(ctx context.Context, slug string, excludeID int) (bool, error) {
	return db.marketplaceDB.IsSlugUnique(ctx, slug, excludeID)
}

func (db *Database) GenerateUniqueSlug(ctx context.Context, baseSlug string, excludeID int) (string, error) {
	return db.marketplaceDB.GenerateUniqueSlug(ctx, baseSlug, excludeID)
}

func (db *Database) UpdateListing(ctx context.Context, listing *models.MarketplaceListing) error {
	return db.marketplaceDB.UpdateListing(ctx, listing)
}

func (db *Database) DeleteListing(ctx context.Context, id int, userID int) error {
	return db.marketplaceDB.DeleteListing(ctx, id, userID)
}

func (db *Database) DeleteListingAdmin(ctx context.Context, id int) error {
	return db.marketplaceDB.DeleteListingAdmin(ctx, id)
}

func (db *Database) GetCategories(ctx context.Context) ([]models.MarketplaceCategory, error) {
	return db.marketplaceDB.GetCategories(ctx)
}

func (db *Database) GetAllCategories(ctx context.Context) ([]models.MarketplaceCategory, error) {
	return db.marketplaceDB.GetAllCategories(ctx)
}

func (db *Database) GetPopularCategories(ctx context.Context, limit int) ([]models.MarketplaceCategory, error) {
	return db.marketplaceDB.GetPopularCategories(ctx, limit)
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

// Storefront favorites
func (db *Database) AddStorefrontToFavorites(ctx context.Context, userID int, productID int) error {
	return db.marketplaceDB.AddStorefrontToFavorites(ctx, userID, productID)
}

func (db *Database) RemoveStorefrontFromFavorites(ctx context.Context, userID int, productID int) error {
	return db.marketplaceDB.RemoveStorefrontFromFavorites(ctx, userID, productID)
}

func (db *Database) GetUserStorefrontFavorites(ctx context.Context, userID int) ([]models.MarketplaceListing, error) {
	return db.marketplaceDB.GetUserStorefrontFavorites(ctx, userID)
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

func (db *Database) UpdateReviewStatus(ctx context.Context, reviewId int, status string) error {
	return db.reviewDB.UpdateReviewStatus(ctx, reviewId, status)
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

func (db *Database) GetAttributeGroups() AttributeGroupStorage {
	return db.attributeGroups
}

// User methods - DEPRECATED: moved to auth-service
func (db *Database) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return nil, fmt.Errorf("GetUserByEmail: moved to auth-service")
}

func (db *Database) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	return nil, fmt.Errorf("GetUserByID: moved to auth-service")
}

func (db *Database) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	return nil, fmt.Errorf("CreateUser: moved to auth-service")
}

func (db *Database) UpdateUser(ctx context.Context, user *models.User) error {
	return fmt.Errorf("UpdateUser: moved to auth-service")
}

func (db *Database) GetOrCreateGoogleUser(ctx context.Context, user *models.User) (*models.User, error) {
	return nil, fmt.Errorf("GetOrCreateGoogleUser: moved to auth-service")
}

func (db *Database) GetUserProfile(ctx context.Context, id int) (*models.UserProfile, error) {
	return nil, fmt.Errorf("GetUserProfile: moved to auth-service")
}

func (db *Database) UpdateUserProfile(ctx context.Context, id int, update *models.UserProfileUpdate) error {
	return fmt.Errorf("UpdateUserProfile: moved to auth-service")
}

func (db *Database) UpdateLastSeen(ctx context.Context, id int) error {
	return fmt.Errorf("UpdateLastSeen: moved to auth-service")
}

// Административные методы для управления пользователями - DEPRECATED: moved to auth-service
func (db *Database) GetAllUsers(ctx context.Context, limit, offset int) ([]*models.UserProfile, int, error) {
	return nil, 0, fmt.Errorf("GetAllUsers: moved to auth-service")
}

func (db *Database) GetAllUsersWithSort(ctx context.Context, limit, offset int, sortBy, sortOrder, statusFilter string) ([]*models.UserProfile, int, error) {
	return nil, 0, fmt.Errorf("GetAllUsersWithSort: moved to auth-service")
}

func (db *Database) UpdateUserStatus(ctx context.Context, id int, status string) error {
	return fmt.Errorf("UpdateUserStatus: moved to auth-service")
}

func (db *Database) UpdateUserRole(ctx context.Context, id int, roleID int) error {
	return fmt.Errorf("UpdateUserRole: moved to auth-service")
}

func (db *Database) GetAllRoles(ctx context.Context) ([]*models.Role, error) {
	return nil, fmt.Errorf("GetAllRoles: moved to auth-service")
}

func (db *Database) DeleteUser(ctx context.Context, id int) error {
	return fmt.Errorf("DeleteUser: moved to auth-service")
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
        FROM marketplace_messages m
        JOIN marketplace_chats c ON m.chat_id = c.id
        WHERE m.receiver_id = $1
        AND NOT m.is_read
        AND NOT c.is_archived
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAttributeTranslationNotFound
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

// IncrementViewsCount увеличивает счетчик просмотров объявления на 1
// IncrementViewsCount увеличивает счетчик просмотров объявления на 1,
// но только если данный пользователь ещё не просматривал это объявление
func (db *Database) IncrementViewsCount(ctx context.Context, id int) error {
	// Логируем вызов функции
	fmt.Printf("IncrementViewsCount called for listing %d\n", id)

	// Получаем ID пользователя из контекста
	var userID int
	if uid := ctx.Value("user_id"); uid != nil {
		if uidInt, ok := uid.(int); ok {
			userID = uidInt
		}
		fmt.Printf("User ID from context: %v (type: %T)\n", uid, uid)
	} else {
		fmt.Printf("No user_id in context\n")
	}

	// Для неавторизованных пользователей используем IP-адрес как идентификатор
	var userIdentifier string
	if ip := ctx.Value("ip_address"); ip != nil {
		if ipStr, ok := ip.(string); ok {
			userIdentifier = ipStr
		}
		fmt.Printf("IP address from context: %v (type: %T)\n", ip, ip)
	} else {
		fmt.Printf("No ip_address in context\n")
	}

	// Проверяем, есть ли уже запись о просмотре этого объявления данным пользователем
	// Если userID > 0, проверяем по user_id, иначе по ip_hash
	var viewExists bool
	var err error

	switch {
	case userID > 0:
		// Для авторизованных пользователей проверяем по ID
		err = db.pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM listing_views
				WHERE listing_id = $1 AND user_id = $2
			)
		`, id, userID).Scan(&viewExists)
	case userIdentifier != "":
		// Для неавторизованных пользователей проверяем строго по IP-адресу,
		// убедившись, что user_id IS NULL (чтобы не конфликтовать с ограничением уникальности)
		err = db.pool.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1 FROM listing_views
				WHERE listing_id = $1 AND ip_hash = $2 AND user_id IS NULL
			)
		`, id, userIdentifier).Scan(&viewExists)
	default:
		// Если нет ни ID пользователя, ни IP - считаем, что просмотр уже был (перестраховка)
		return nil
	}

	if err != nil {
		return err
	}

	// Если просмотра ещё не было, увеличиваем счетчик и добавляем запись о просмотре
	fmt.Printf("View exists check for listing %d: %v\n", id, viewExists)
	if !viewExists {
		fmt.Printf("View does not exist, incrementing counter for listing %d\n", id)
		// Начинаем транзакцию
		tx, err := db.pool.Begin(ctx)
		if err != nil {
			return err
		}
		defer func() {
			if err := tx.Rollback(ctx); err != nil {
				// Игнорируем ошибку если транзакция уже была завершена
				_ = err // Explicitly ignore error
			}
		}()

		// Увеличиваем счетчик просмотров
		// Сначала пробуем обновить в marketplace_listings
		commandTag, err := tx.Exec(ctx, `
			UPDATE marketplace_listings
			SET views_count = views_count + 1
			WHERE id = $1
		`, id)
		if err != nil {
			fmt.Printf("Error updating marketplace_listings: %v\n", err)
			return err
		}

		rowsAffected := commandTag.RowsAffected()
		fmt.Printf("marketplace_listings rows affected: %d\n", rowsAffected)

		// Если не нашли в marketplace_listings, пробуем в storefront_products
		if rowsAffected == 0 {
			fmt.Printf("Trying to update storefront_products for id %d\n", id)
			commandTag, err = tx.Exec(ctx, `
				UPDATE storefront_products
				SET view_count = view_count + 1
				WHERE id = $1
			`, id)
			if err != nil {
				fmt.Printf("Error updating storefront_products: %v\n", err)
				return err
			}

			rowsAffected = commandTag.RowsAffected()
			fmt.Printf("storefront_products rows affected: %d\n", rowsAffected)
			if rowsAffected == 0 {
				// Если ни в одной таблице не нашли товар
				return fmt.Errorf("listing or product with id %d not found", id)
			}
		}

		// Добавляем запись о просмотре
		if userID > 0 {
			_, err = tx.Exec(ctx, `
				INSERT INTO listing_views (listing_id, user_id, view_time)
				VALUES ($1, $2, NOW())
			`, id, userID)
		} else {
			_, err = tx.Exec(ctx, `
				INSERT INTO listing_views (listing_id, ip_hash, view_time, user_id)
				VALUES ($1, $2, NOW(), NULL)
			`, id, userIdentifier)
		}
		if err != nil {
			return err
		}

		// Фиксируем транзакцию
		fmt.Printf("Committing transaction for listing %d view increment\n", id)
		err = tx.Commit(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("View increment committed successfully for listing %d\n", id)

		// После успешного обновления в PostgreSQL синхронизируем данные с OpenSearch
		if db.osMarketplaceRepo != nil && db.osClient != nil {
			// Получаем обновленное значение счетчика просмотров
			var viewsCount int
			// Сначала пробуем получить из marketplace_listings
			err = db.pool.QueryRow(ctx, "SELECT views_count FROM marketplace_listings WHERE id = $1", id).Scan(&viewsCount)
			if err != nil {
				// Если не нашли в marketplace_listings, пробуем из storefront_products
				err = db.pool.QueryRow(ctx, "SELECT view_count FROM storefront_products WHERE id = $1", id).Scan(&viewsCount)
				if err != nil {
					log.Printf("Ошибка при получении обновленного счетчика просмотров: %v", err)
					// Не прерываем выполнение, так как главное - обновить в PostgreSQL
				} else {
					// Обновляем данные в OpenSearch
					go db.updateViewCountInOpenSearch(id, viewsCount) //nolint:contextcheck // фоновое обновление
				}
			} else {
				// Обновляем данные в OpenSearch
				go db.updateViewCountInOpenSearch(id, viewsCount) //nolint:contextcheck // фоновое обновление
			}
		}
	}

	return nil
}

// updateViewCountInOpenSearch обновляет счетчик просмотров в индексе OpenSearch
func (db *Database) updateViewCountInOpenSearch(id int, viewsCount int) {
	ctx := context.Background()

	// Получаем объявление из PostgreSQL
	listing, err := db.GetListingByID(ctx, id)
	if err != nil {
		log.Printf("Ошибка при получении листинга для обновления OpenSearch: %v", err)
		return
	}

	// Устанавливаем значение счетчика просмотров
	listing.ViewsCount = viewsCount

	// Индексируем объявление в OpenSearch
	if db.osMarketplaceRepo != nil {
		err = db.osMarketplaceRepo.IndexListing(ctx, listing)
		if err != nil {
			log.Printf("Ошибка при обновлении счетчика просмотров в OpenSearch: %v", err)
		} else {
			log.Printf("Успешно обновлен счетчик просмотров в OpenSearch для объявления %d: %d", id, viewsCount)
		}
	}
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

// User Contacts methods - delegating to marketplace storage
func (db *Database) AddContact(ctx context.Context, contact *models.UserContact) error {
	return db.marketplaceDB.AddContact(ctx, contact)
}

func (db *Database) UpdateContactStatus(ctx context.Context, userID, contactUserID int, status, notes string) error {
	return db.marketplaceDB.UpdateContactStatus(ctx, userID, contactUserID, status, notes)
}

func (db *Database) GetContact(ctx context.Context, userID, contactUserID int) (*models.UserContact, error) {
	return db.marketplaceDB.GetContact(ctx, userID, contactUserID)
}

func (db *Database) GetUserContacts(ctx context.Context, userID int, status string, page, limit int) ([]models.UserContact, int, error) {
	return db.marketplaceDB.GetUserContacts(ctx, userID, status, page, limit)
}

func (db *Database) GetIncomingContactRequests(ctx context.Context, userID int, page, limit int) ([]models.UserContact, int, error) {
	return db.marketplaceDB.GetIncomingContactRequests(ctx, userID, page, limit)
}

func (db *Database) RemoveContact(ctx context.Context, userID, contactUserID int) error {
	return db.marketplaceDB.RemoveContact(ctx, userID, contactUserID)
}

func (db *Database) GetUserPrivacySettings(ctx context.Context, userID int) (*models.UserPrivacySettings, error) {
	return db.marketplaceDB.GetUserPrivacySettings(ctx, userID)
}

func (db *Database) UpdateUserPrivacySettings(ctx context.Context, userID int, settings *models.UpdatePrivacySettingsRequest) error {
	return db.marketplaceDB.UpdateUserPrivacySettings(ctx, userID, settings)
}

func (db *Database) GetPrivacySettings(ctx context.Context, userID int) (*models.UserPrivacySettings, error) {
	return db.marketplaceDB.GetUserPrivacySettings(ctx, userID)
}

func (db *Database) UpdatePrivacySettings(ctx context.Context, userID int, settings *models.UpdatePrivacySettingsRequest) error {
	return db.marketplaceDB.UpdateUserPrivacySettings(ctx, userID, settings)
}

func (db *Database) CanAddContact(ctx context.Context, userID, targetUserID int) (bool, error) {
	return db.marketplaceDB.CanAddContact(ctx, userID, targetUserID)
}

func (db *Database) AreContacts(ctx context.Context, userID1, userID2 int) (bool, error) {
	return db.marketplaceDB.AreContacts(ctx, userID1, userID2)
}

// GetChatActivityStats получает статистику активности в чате между покупателем и продавцом
func (db *Database) GetChatActivityStats(ctx context.Context, buyerID int, sellerID int, listingID int) (*models.ChatActivityStats, error) {
	stats := &models.ChatActivityStats{}

	// Проверяем наличие чата и получаем статистику
	query := `
		WITH chat_info AS (
			SELECT
				c.id as chat_id,
				c.created_at as chat_created
			FROM marketplace_chats c
			WHERE c.buyer_id = $1
				AND c.seller_id = $2
				AND c.listing_id = $3
			LIMIT 1
		),
		message_stats AS (
			SELECT
				COUNT(*) as total_messages,
				COUNT(*) FILTER (WHERE m.sender_id = $1) as buyer_messages,
				COUNT(*) FILTER (WHERE m.sender_id = $2) as seller_messages,
				MIN(m.created_at) as first_message_date,
				MAX(m.created_at) as last_message_date
			FROM marketplace_messages m
			INNER JOIN chat_info ci ON m.chat_id = ci.chat_id
		)
		SELECT
			CASE WHEN ci.chat_id IS NOT NULL THEN true ELSE false END as chat_exists,
			COALESCE(ms.total_messages, 0) as total_messages,
			COALESCE(ms.buyer_messages, 0) as buyer_messages,
			COALESCE(ms.seller_messages, 0) as seller_messages,
			ms.first_message_date,
			ms.last_message_date,
			CASE
				WHEN ms.first_message_date IS NOT NULL
				THEN EXTRACT(DAY FROM NOW() - ms.first_message_date)::int
				ELSE 0
			END as days_since_first_msg,
			CASE
				WHEN ms.last_message_date IS NOT NULL
				THEN EXTRACT(DAY FROM NOW() - ms.last_message_date)::int
				ELSE 0
			END as days_since_last_msg
		FROM chat_info ci
		LEFT JOIN message_stats ms ON true
	`

	row := db.pool.QueryRow(ctx, query, buyerID, sellerID, listingID)

	err := row.Scan(
		&stats.ChatExists,
		&stats.TotalMessages,
		&stats.BuyerMessages,
		&stats.SellerMessages,
		&stats.FirstMessageDate,
		&stats.LastMessageDate,
		&stats.DaysSinceFirstMsg,
		&stats.DaysSinceLastMsg,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		// Чат не существует, возвращаем пустую статистику
		return stats, nil
	}

	return stats, err
}

// GetUserAggregatedRating получает агрегированный рейтинг пользователя
func (db *Database) GetUserAggregatedRating(ctx context.Context, userID int) (*models.UserAggregatedRating, error) {
	rating := &models.UserAggregatedRating{UserID: userID}

	query := `
		SELECT
			total_reviews, average_rating, direct_reviews, listing_reviews,
			storefront_reviews, verified_reviews, rating_1, rating_2, rating_3,
			rating_4, rating_5, recent_rating, recent_reviews, last_review_at
		FROM user_ratings
		WHERE user_id = $1
	`

	row := db.pool.QueryRow(ctx, query, userID)
	err := row.Scan(
		&rating.TotalReviews, &rating.AverageRating, &rating.DirectReviews,
		&rating.ListingReviews, &rating.StorefrontReviews, &rating.VerifiedReviews,
		&rating.Rating1, &rating.Rating2, &rating.Rating3,
		&rating.Rating4, &rating.Rating5, &rating.RecentRating,
		&rating.RecentReviews, &rating.LastReviewAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		// Если нет данных в материализованном представлении, возвращаем пустой рейтинг
		return rating, nil
	}

	return rating, err
}

// GetStorefrontAggregatedRating получает агрегированный рейтинг магазина
func (db *Database) GetStorefrontAggregatedRating(ctx context.Context, storefrontID int) (*models.StorefrontAggregatedRating, error) {
	rating := &models.StorefrontAggregatedRating{StorefrontID: storefrontID}

	query := `
		SELECT
			total_reviews, average_rating, direct_reviews, listing_reviews,
			verified_reviews, rating_1, rating_2, rating_3, rating_4, rating_5,
			recent_rating, recent_reviews, last_review_at, owner_id
		FROM storefront_ratings
		WHERE storefront_id = $1
	`

	row := db.pool.QueryRow(ctx, query, storefrontID)
	err := row.Scan(
		&rating.TotalReviews, &rating.AverageRating, &rating.DirectReviews,
		&rating.ListingReviews, &rating.VerifiedReviews,
		&rating.Rating1, &rating.Rating2, &rating.Rating3,
		&rating.Rating4, &rating.Rating5, &rating.RecentRating,
		&rating.RecentReviews, &rating.LastReviewAt, &rating.OwnerID,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return rating, nil
	}

	return rating, err
}

// RefreshRatingViews обновляет материализованные представления
func (db *Database) RefreshRatingViews(ctx context.Context) error {
	_, err := db.pool.Exec(ctx, "SELECT rebuild_all_ratings()")
	return err
}

// CreateReviewConfirmation создает подтверждение отзыва
func (db *Database) CreateReviewConfirmation(ctx context.Context, confirmation *models.ReviewConfirmation) error {
	query := `
		INSERT INTO review_confirmations
		(review_id, confirmed_by, confirmation_status, notes)
		VALUES ($1, $2, $3, $4)
		RETURNING id, confirmed_at
	`

	row := db.pool.QueryRow(ctx, query,
		confirmation.ReviewID, confirmation.ConfirmedBy,
		confirmation.ConfirmationStatus, confirmation.Notes,
	)

	return row.Scan(&confirmation.ID, &confirmation.ConfirmedAt)
}

// GetReviewConfirmation получает подтверждение отзыва
func (db *Database) GetReviewConfirmation(ctx context.Context, reviewID int) (*models.ReviewConfirmation, error) {
	confirmation := &models.ReviewConfirmation{}

	query := `
		SELECT id, review_id, confirmed_by, confirmation_status, confirmed_at, notes
		FROM review_confirmations
		WHERE review_id = $1
	`

	row := db.pool.QueryRow(ctx, query, reviewID)
	err := row.Scan(
		&confirmation.ID, &confirmation.ReviewID, &confirmation.ConfirmedBy,
		&confirmation.ConfirmationStatus, &confirmation.ConfirmedAt, &confirmation.Notes,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrReviewConfirmationNotFound
	}

	return confirmation, err
}

// CreateReviewDispute создает спор по отзыву
func (db *Database) CreateReviewDispute(ctx context.Context, dispute *models.ReviewDispute) error {
	query := `
		INSERT INTO review_disputes
		(review_id, disputed_by, dispute_reason, dispute_description, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	row := db.pool.QueryRow(ctx, query,
		dispute.ReviewID, dispute.DisputedBy, dispute.DisputeReason,
		dispute.DisputeDescription, dispute.Status,
	)

	err := row.Scan(&dispute.ID, &dispute.CreatedAt, &dispute.UpdatedAt)
	if err != nil {
		return err
	}

	// Обновляем флаг в таблице reviews
	_, err = db.pool.Exec(ctx,
		"UPDATE reviews SET has_active_dispute = true WHERE id = $1",
		dispute.ReviewID,
	)

	return err
}

// GetReviewDispute получает спор по отзыву
func (db *Database) GetReviewDispute(ctx context.Context, reviewID int) (*models.ReviewDispute, error) {
	dispute := &models.ReviewDispute{}

	query := `
		SELECT id, review_id, disputed_by, dispute_reason, dispute_description,
			   status, admin_id, admin_notes, created_at, updated_at, resolved_at
		FROM review_disputes
		WHERE review_id = $1 AND status NOT IN ('resolved_keep_review', 'resolved_remove_review', 'cancelled')
		ORDER BY created_at DESC
		LIMIT 1
	`

	row := db.pool.QueryRow(ctx, query, reviewID)
	err := row.Scan(
		&dispute.ID, &dispute.ReviewID, &dispute.DisputedBy,
		&dispute.DisputeReason, &dispute.DisputeDescription, &dispute.Status,
		&dispute.AdminID, &dispute.AdminNotes, &dispute.CreatedAt,
		&dispute.UpdatedAt, &dispute.ResolvedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrReviewDisputeNotFound
	}

	return dispute, err
}

// UpdateReviewDispute обновляет спор
func (db *Database) UpdateReviewDispute(ctx context.Context, dispute *models.ReviewDispute) error {
	query := `
		UPDATE review_disputes
		SET status = $2, admin_id = $3, admin_notes = $4,
			updated_at = NOW(), resolved_at = $5
		WHERE id = $1
	`

	_, err := db.pool.Exec(ctx, query,
		dispute.ID, dispute.Status, dispute.AdminID,
		dispute.AdminNotes, dispute.ResolvedAt,
	)
	if err != nil {
		return err
	}

	// Если спор разрешен, обновляем флаг в reviews
	if dispute.Status == "resolved_keep_review" ||
		dispute.Status == "resolved_remove_review" ||
		dispute.Status == "canceled" {
		_, err = db.pool.Exec(ctx,
			"UPDATE reviews SET has_active_dispute = false WHERE id = $1",
			dispute.ReviewID,
		)
	}

	return err
}

// CanUserReviewEntity проверяет может ли пользователь оставить отзыв
func (db *Database) CanUserReviewEntity(ctx context.Context, userID int, entityType string, entityID int) (*models.CanReviewResponse, error) {
	response := &models.CanReviewResponse{
		CanReview: true,
	}

	// Проверяем существующий отзыв
	var existingReviewID *int
	query := `
		SELECT id FROM reviews
		WHERE user_id = $1 AND entity_type = $2 AND entity_id = $3
		AND status != 'deleted'
		LIMIT 1
	`

	err := db.pool.QueryRow(ctx, query, userID, entityType, entityID).Scan(&existingReviewID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	if existingReviewID != nil {
		response.CanReview = false
		response.HasExistingReview = true
		response.ExistingReviewID = existingReviewID
		response.Reason = "Вы уже оставили отзыв на этот объект"
		return response, nil
	}

	// Для отзывов на товары проверяем владельца
	if entityType == "listing" {
		var ownerID int
		err := db.pool.QueryRow(ctx,
			"SELECT user_id FROM marketplace_listings WHERE id = $1",
			entityID,
		).Scan(&ownerID)

		if err == nil && ownerID == userID {
			response.CanReview = false
			response.Reason = "Вы не можете оставить отзыв на свой товар"
			return response, nil
		}
	}

	return response, nil
}

// stringPtr creates a pointer to a string
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// Storefront methods
func (db *Database) CreateStorefront(ctx context.Context, userID int, dto *models.StorefrontCreateDTO) (*models.Storefront, error) {
	storefront := &models.Storefront{
		UserID:           userID,
		Slug:             dto.Slug,
		Name:             dto.Name,
		Description:      stringPtr(dto.Description),
		LogoURL:          stringPtr(""), // Будет заполнено после загрузки
		BannerURL:        stringPtr(""), // Будет заполнено после загрузки
		Theme:            dto.Theme,
		Phone:            stringPtr(dto.Phone),
		Email:            stringPtr(dto.Email),
		Website:          stringPtr(dto.Website),
		Address:          stringPtr(dto.Location.FullAddress),
		City:             stringPtr(dto.Location.City),
		PostalCode:       stringPtr(dto.Location.PostalCode),
		Country:          stringPtr(dto.Location.Country),
		Latitude:         &dto.Location.BuildingLat,
		Longitude:        &dto.Location.BuildingLng,
		Settings:         dto.Settings,
		SEOMeta:          dto.SEOMeta,
		IsActive:         true,
		SubscriptionPlan: "basic",
		CommissionRate:   0.05, // 5% по умолчанию
	}

	err := db.pool.QueryRow(ctx, `
		INSERT INTO storefronts (user_id, slug, name, description, logo_url, banner_url, theme,
			phone, email, website, address, city, postal_code, country, latitude, longitude,
			settings, seo_meta, is_active, subscription_plan, commission_rate)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
		RETURNING id, created_at, updated_at
	`, storefront.UserID, storefront.Slug, storefront.Name, storefront.Description,
		storefront.LogoURL, storefront.BannerURL, storefront.Theme, storefront.Phone,
		storefront.Email, storefront.Website, storefront.Address, storefront.City,
		storefront.PostalCode, storefront.Country, storefront.Latitude, storefront.Longitude,
		storefront.Settings, storefront.SEOMeta, storefront.IsActive, storefront.SubscriptionPlan,
		storefront.CommissionRate).Scan(&storefront.ID, &storefront.CreatedAt, &storefront.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Создаем запись в unified_geo для поддержки GIS поиска витрин
	if storefront.Latitude != nil && storefront.Longitude != nil {
		// Calculate geohash from coordinates
		geohashStr := fmt.Sprintf("%.6f,%.6f", *storefront.Latitude, *storefront.Longitude)

		addressComponents := map[string]interface{}{
			"city":        storefront.City,
			"postal_code": storefront.PostalCode,
			"country":     storefront.Country,
		}
		addressComponentsJSON, _ := json.Marshal(addressComponents)

		_, err = db.pool.Exec(ctx, `
			INSERT INTO unified_geo (
				source_type, source_id, location, geohash,
				formatted_address, address_components,
				geocoding_confidence, address_verified,
				input_method, location_privacy, blur_radius,
				is_precise
			) VALUES (
				'storefront', $1, ST_SetSRID(ST_MakePoint($2, $3), 4326), $4,
				$5, $6,
				0.9, true,
				'manual', 'exact', 0,
				true
			)
			ON CONFLICT (source_type, source_id)
			DO UPDATE SET
				location = EXCLUDED.location,
				geohash = EXCLUDED.geohash,
				formatted_address = EXCLUDED.formatted_address,
				address_components = EXCLUDED.address_components,
				updated_at = CURRENT_TIMESTAMP
		`, storefront.ID, *storefront.Longitude, *storefront.Latitude, geohashStr,
			storefront.Address, addressComponentsJSON)

		if err != nil {
			log.Printf("Error creating unified_geo entry for storefront %d: %v", storefront.ID, err)
			// Не прерываем создание витрины из-за ошибки с geo
		} else {
			// Обновляем materialized view после успешного создания geo записи
			_, err = db.pool.Exec(ctx, "SELECT refresh_map_items_cache()")
			if err != nil {
				log.Printf("Error refreshing map_items_cache: %v", err)
			}
		}
	}

	return storefront, nil
}

func (db *Database) GetUserStorefronts(ctx context.Context, userID int) ([]models.Storefront, error) {
	rows, err := db.pool.Query(ctx, `
		SELECT id, user_id, slug, name, description, logo_url, banner_url, theme,
			phone, email, website, address, city, postal_code, country, latitude, longitude,
			settings, seo_meta, is_active, is_verified, verification_date, rating, reviews_count,
			products_count, sales_count, views_count, subscription_plan, subscription_expires_at,
			commission_rate, ai_agent_enabled, ai_agent_config, live_shopping_enabled,
			group_buying_enabled, created_at, updated_at
		FROM storefronts
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var storefronts []models.Storefront
	for rows.Next() {
		var s models.Storefront
		err := rows.Scan(
			&s.ID, &s.UserID, &s.Slug, &s.Name, &s.Description, &s.LogoURL, &s.BannerURL, &s.Theme,
			&s.Phone, &s.Email, &s.Website, &s.Address, &s.City, &s.PostalCode, &s.Country,
			&s.Latitude, &s.Longitude, &s.Settings, &s.SEOMeta, &s.IsActive, &s.IsVerified,
			&s.VerificationDate, &s.Rating, &s.ReviewsCount, &s.ProductsCount, &s.SalesCount,
			&s.ViewsCount, &s.SubscriptionPlan, &s.SubscriptionExpiresAt, &s.CommissionRate,
			&s.AIAgentEnabled, &s.AIAgentConfig, &s.LiveShoppingEnabled, &s.GroupBuyingEnabled,
			&s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		storefronts = append(storefronts, s)
	}

	return storefronts, nil
}

func (db *Database) GetStorefrontByID(ctx context.Context, id int) (*models.Storefront, error) {
	var s models.Storefront
	var theme, settings, seoMeta, aiConfig json.RawMessage

	err := db.pool.QueryRow(ctx, `
		SELECT id, user_id, slug, name, description, logo_url, banner_url,
			COALESCE(theme, '{}')::jsonb,
			phone, email, website, address, city, postal_code, country, latitude, longitude,
			COALESCE(settings, '{}')::jsonb, COALESCE(seo_meta, '{}')::jsonb,
			is_active, is_verified, verification_date, rating, reviews_count,
			products_count, sales_count, views_count, subscription_plan, subscription_expires_at,
			commission_rate, ai_agent_enabled, COALESCE(ai_agent_config, '{}')::jsonb,
			live_shopping_enabled, group_buying_enabled, created_at, updated_at
		FROM storefronts
		WHERE id = $1
	`, id).Scan(
		&s.ID, &s.UserID, &s.Slug, &s.Name, &s.Description, &s.LogoURL, &s.BannerURL, &theme,
		&s.Phone, &s.Email, &s.Website, &s.Address, &s.City, &s.PostalCode, &s.Country,
		&s.Latitude, &s.Longitude, &settings, &seoMeta, &s.IsActive, &s.IsVerified,
		&s.VerificationDate, &s.Rating, &s.ReviewsCount, &s.ProductsCount, &s.SalesCount,
		&s.ViewsCount, &s.SubscriptionPlan, &s.SubscriptionExpiresAt, &s.CommissionRate,
		&s.AIAgentEnabled, &aiConfig, &s.LiveShoppingEnabled, &s.GroupBuyingEnabled,
		&s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrStorefrontNotFound
		}
		return nil, err
	}

	// Конвертируем json.RawMessage в JSONB
	if theme != nil {
		if err := json.Unmarshal(theme, &s.Theme); err != nil {
			// Логируем ошибку, но не прерываем выполнение
			_ = err // Explicitly ignore error
		}
	}
	if settings != nil {
		if err := json.Unmarshal(settings, &s.Settings); err != nil {
			// Логируем ошибку, но не прерываем выполнение
			_ = err // Explicitly ignore error
		}
	}
	if seoMeta != nil {
		if err := json.Unmarshal(seoMeta, &s.SEOMeta); err != nil {
			// Логируем ошибку, но не прерываем выполнение
			_ = err // Explicitly ignore error
		}
	}
	if aiConfig != nil {
		if err := json.Unmarshal(aiConfig, &s.AIAgentConfig); err != nil {
			// Логируем ошибку, но не прерываем выполнение
			_ = err // Explicitly ignore error
		}
	}

	return &s, nil
}

// GetStorefrontOwnerByProductID возвращает ID владельца витрины по ID товара
func (db *Database) GetStorefrontOwnerByProductID(ctx context.Context, productID int) (int, error) {
	var userID int
	err := db.pool.QueryRow(ctx, `
		SELECT s.user_id
		FROM storefronts s
		INNER JOIN storefront_products sp ON sp.storefront_id = s.id
		WHERE sp.id = $1
	`, productID).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("storefront product not found: %d", productID)
		}
		return 0, err
	}

	return userID, nil
}

func (db *Database) UpdateStorefront(ctx context.Context, storefront *models.Storefront) error {
	_, err := db.pool.Exec(ctx, `
		UPDATE storefronts
		SET name = $2, description = $3, logo_url = $4, banner_url = $5, theme = $6,
			phone = $7, email = $8, website = $9, address = $10, city = $11,
			postal_code = $12, country = $13, latitude = $14, longitude = $15,
			settings = $16, seo_meta = $17, is_active = $18, updated_at = NOW()
		WHERE id = $1
	`, storefront.ID, storefront.Name, storefront.Description, storefront.LogoURL,
		storefront.BannerURL, storefront.Theme, storefront.Phone, storefront.Email,
		storefront.Website, storefront.Address, storefront.City, storefront.PostalCode,
		storefront.Country, storefront.Latitude, storefront.Longitude, storefront.Settings,
		storefront.SEOMeta, storefront.IsActive)
	return err
}

func (db *Database) DeleteStorefront(ctx context.Context, id int) error {
	_, err := db.pool.Exec(ctx, "DELETE FROM storefronts WHERE id = $1", id)
	return err
}

// Storefront возвращает репозиторий витрин
func (db *Database) Storefront() interface{} {
	if db.storefrontRepo != nil {
		return db.storefrontRepo
	}
	// Возвращаем новый репозиторий используя текущий экземпляр db
	return NewStorefrontRepository(db)
}

// Cart возвращает репозиторий корзин
func (db *Database) Cart() interface{} {
	if db.cartRepo != nil {
		return db.cartRepo
	}
	// Возвращаем новый репозиторий используя пул соединений
	return NewCartRepository(db.pool)
}

// Order возвращает репозиторий заказов
func (db *Database) Order() interface{} {
	if db.orderRepo != nil {
		return db.orderRepo
	}
	// Возвращаем новый репозиторий используя пул соединений
	return NewOrderRepository(db.pool)
}

// Inventory возвращает репозиторий инвентаря
func (db *Database) Inventory() interface{} {
	if db.inventoryRepo != nil {
		return db.inventoryRepo
	}
	// Возвращаем новый репозиторий используя пул соединений
	return NewInventoryRepository(db.pool)
}

// MarketplaceOrder возвращает репозиторий заказов маркетплейса
func (db *Database) MarketplaceOrder() interface{} {
	if db.marketplaceOrderRepo != nil {
		return db.marketplaceOrderRepo
	}
	// Возвращаем новый репозиторий используя пул соединений
	return NewMarketplaceOrderRepository(db.pool)
}

// StorefrontProductSearch возвращает репозиторий для поиска товаров витрин
func (db *Database) StorefrontProductSearch() interface{} {
	if db.productSearchRepo != nil {
		return db.productSearchRepo
	}
	// Возвращаем новый репозиторий если есть клиент OpenSearch
	if db.osClient != nil {
		// ВАЖНО: Используем тот же индекс что и для marketplace (унифицированный поиск)
		return storefrontOpenSearch.NewProductRepository(db.osClient, "marketplace_listings")
	}
	return nil
}

// QueryContext выполняет SQL запрос с контекстом и возвращает строки результата
func (db *Database) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return db.db.QueryContext(ctx, query, args...)
}

// QueryRowContext выполняет SQL запрос с контекстом и возвращает одну строку результата
func (db *Database) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return db.db.QueryRowContext(ctx, query, args...)
}

// ExecContext выполняет SQL запрос с контекстом без возврата результата
func (db *Database) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.db.ExecContext(ctx, query, args...)
}

// Marketplace listing variants methods
func (db *Database) CreateListingVariants(ctx context.Context, listingID int, variants []models.MarketplaceListingVariant) error {
	if len(variants) == 0 {
		return nil
	}

	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			log.Printf("Failed to rollback transaction: %v", rollbackErr)
		}
	}()

	for _, variant := range variants {
		attributesJSON, err := json.Marshal(variant.Attributes)
		if err != nil {
			return fmt.Errorf("failed to marshal variant attributes: %w", err)
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO marketplace_listing_variants (listing_id, sku, price, stock, attributes, image_url, is_active)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, listingID, variant.SKU, variant.Price, variant.Stock, attributesJSON, variant.ImageURL, variant.IsActive)
		if err != nil {
			return fmt.Errorf("failed to insert variant: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (db *Database) GetListingVariants(ctx context.Context, listingID int) ([]models.MarketplaceListingVariant, error) {
	query := `
		SELECT id, listing_id, sku, price, stock, attributes, image_url, is_active,
		       created_at::text, updated_at::text
		FROM marketplace_listing_variants
		WHERE listing_id = $1 AND is_active = true
		ORDER BY id
	`

	rows, err := db.pool.Query(ctx, query, listingID)
	if err != nil {
		return nil, fmt.Errorf("failed to query variants: %w", err)
	}
	defer rows.Close()

	var variants []models.MarketplaceListingVariant
	for rows.Next() {
		var variant models.MarketplaceListingVariant
		var attributesJSON []byte

		err := rows.Scan(
			&variant.ID, &variant.ListingID, &variant.SKU, &variant.Price, &variant.Stock,
			&attributesJSON, &variant.ImageURL, &variant.IsActive,
			&variant.CreatedAt, &variant.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan variant: %w", err)
		}

		if len(attributesJSON) > 0 {
			err = json.Unmarshal(attributesJSON, &variant.Attributes)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal variant attributes: %w", err)
			}
		}

		variants = append(variants, variant)
	}

	return variants, rows.Err()
}

func (db *Database) UpdateListingVariant(ctx context.Context, variant *models.MarketplaceListingVariant) error {
	attributesJSON, err := json.Marshal(variant.Attributes)
	if err != nil {
		return fmt.Errorf("failed to marshal variant attributes: %w", err)
	}

	query := `
		UPDATE marketplace_listing_variants
		SET sku = $1, price = $2, stock = $3, attributes = $4, image_url = $5, is_active = $6
		WHERE id = $7
	`

	result, err := db.pool.Exec(ctx, query,
		variant.SKU, variant.Price, variant.Stock, attributesJSON,
		variant.ImageURL, variant.IsActive, variant.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update variant: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("variant not found")
	}

	return nil
}

func (db *Database) DeleteListingVariant(ctx context.Context, variantID int) error {
	// Soft delete - просто помечаем как неактивный
	query := `UPDATE marketplace_listing_variants SET is_active = false WHERE id = $1`

	result, err := db.pool.Exec(ctx, query, variantID)
	if err != nil {
		return fmt.Errorf("failed to delete variant: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("variant not found")
	}

	return nil
}

// SetMarketplaceUserService устанавливает UserService для marketplace storage
func (db *Database) SetMarketplaceUserService(userService *authservice.UserService) {
	if db.marketplaceDB != nil {
		db.marketplaceDB.SetUserService(userService)
	}
}
