// backend/internal/storage/postgres/db.go
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"backend/internal/config"
	"backend/internal/storage"
	"backend/internal/storage/filestorage"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"

	notificationStorage "backend/internal/proj/notifications/storage/postgres"
	reviewStorage "backend/internal/proj/reviews/storage/postgres"

	marketplaceStorage "backend/internal/proj/marketplace/storage"
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

// Database представляет главную структуру для работы с базой данных
type Database struct {
	pool                 *pgxpool.Pool
	reviewDB             *reviewStorage.Storage
	notificationsDB      *notificationStorage.Storage
	osClient             *osClient.OpenSearchClient // Клиент OpenSearch для прямых запросов
	db                   *sql.DB
	sqlxDB               *sqlx.DB // sqlx.DB для работы с sqlx библиотекой
	marketplaceIndex     string
	storefrontIndex      string
	attributeGroups      AttributeGroupStorage
	fsStorage            filestorage.FileStorageInterface
	storefrontRepo       StorefrontRepository                  // Репозиторий для витрин
	cartRepo             CartRepositoryInterface               // Репозиторий для корзин
	orderRepo            OrderRepositoryInterface              // Репозиторий для заказов
	searchWeights        *config.SearchWeights                 // Веса для поиска
	inventoryRepo        InventoryRepositoryInterface          // Репозиторий для инвентаря
	marketplaceOrderRepo *MarketplaceOrderRepository           // Репозиторий для заказов маркетплейса
	productSearchRepo    ProductSearchRepositoryStub           // Заглушка для поиска товаров витрин (TODO: восстановить после рефакторинга OpenSearch)
	marketplaceStorage   marketplaceStorage.MarketplaceStorage // Marketplace storage
	grpcClient           *MarketplaceGRPCClient                // gRPC клиент для listings микросервиса
	config               *config.Config                        // Конфигурация приложения (для feature flags)
}

// NewDatabase создает новый экземпляр Database
func NewDatabase(ctx context.Context, dbURL string, osClient *osClient.OpenSearchClient, indexName string, fileStorage filestorage.FileStorageInterface, searchWeights *config.SearchWeights, cfg *config.Config) (*Database, error) {
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

	db := &Database{
		pool:             pool,
		db:               stdDB,
		sqlxDB:           sqlxDB,
		reviewDB:         reviewStorage.NewStorage(pool, nil), // translation service теперь не нужен для reviews
		notificationsDB:  notificationStorage.NewNotificationStorage(pool),
		osClient:         osClient,     // Сохраняем клиент OpenSearch
		marketplaceIndex: indexName,    // Сохраняем имя индекса
		storefrontIndex:  "b2c_stores", // Индекс для витрин
		fsStorage:        fileStorage,  // Используем переданный параметр
		attributeGroups:  NewAttributeGroupStorage(pool),
		config:           cfg, // Сохраняем конфигурацию для feature flags
	}

	// Инициализируем репозиторий витрин
	db.storefrontRepo = NewStorefrontRepository(db)

	// Инициализируем репозиторий корзин (с логгером и listingsClient будет передан позже)
	cartLogger := zerolog.New(log.Writer()).With().Timestamp().Str("component", "cart_repository").Logger()
	db.cartRepo = NewCartRepository(pool, nil, cartLogger)

	// Инициализируем репозиторий заказов
	db.orderRepo = NewOrderRepository(pool)

	// Инициализируем репозиторий инвентаря
	db.inventoryRepo = NewInventoryRepository(pool)

	// Инициализируем репозиторий заказов маркетплейса
	db.marketplaceOrderRepo = NewMarketplaceOrderRepository(pool)

	// TODO: Инициализация OpenSearch репозиториев временно отключена
	// Необходимо рефакторинг OpenSearch модуля
	if osClient != nil {
		log.Println("OpenSearch client available, but repositories initialization is disabled during refactoring")
		// Используем заглушку для productSearchRepo
		db.productSearchRepo = &productSearchStub{}
	}

	// Сохраняем search weights
	db.searchWeights = searchWeights

	// Инициализируем marketplace storage (без listings client здесь - будет передан из server.go)
	logger := zerolog.New(log.Writer()).With().Timestamp().Str("component", "marketplace_storage").Logger()
	db.marketplaceStorage = marketplaceStorage.NewPostgresMarketplaceStorage(sqlxDB, logger)

	// Инициализируем gRPC клиент для listings микросервиса
	// DEV MODE: Microservice is REQUIRED - fail fast if URL is empty or connection fails
	grpcURL := cfg.ListingsGRPCURL
	if grpcURL == "" {
		return nil, fmt.Errorf("LISTINGS_GRPC_URL is required - listings microservice must be running. Backend cannot start without it")
	}

	grpcClient, err := NewMarketplaceGRPCClient(ctx, grpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to listings microservice at %s: %w. Make sure microservice is running on this address", grpcURL, err)
	}

	db.grpcClient = grpcClient
	logger.Info().
		Str("grpc_url", grpcURL).
		Msg("Connected to listings microservice (REQUIRED for operation)")

	return db, nil
}

var _ storage.Storage = (*Database)(nil)

// ProductSearchRepositoryStub - временная заглушка для ProductSearchRepository
// TODO: Восстановить после завершения рефакторинга OpenSearch
type ProductSearchRepositoryStub interface{}

// productSearchStub - реализация заглушки
type productSearchStub struct{}

// Ensure stub implements the interface (if needed)
var _ ProductSearchRepositoryStub = (*productSearchStub)(nil)
