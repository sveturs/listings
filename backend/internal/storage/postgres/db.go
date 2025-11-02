// backend/internal/storage/postgres/db.go
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"backend/internal/config"
	"backend/internal/storage"
	"backend/internal/storage/filestorage"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	notificationStorage "backend/internal/proj/notifications/storage/postgres"
	reviewStorage "backend/internal/proj/reviews/storage/postgres"

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
	b2cProductIndex      string // Индекс для B2C товаров
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

// NewDatabase создает новый экземпляр Database
func NewDatabase(ctx context.Context, dbURL string, osClient *osClient.OpenSearchClient, indexName string, b2cIndexName string, fileStorage filestorage.FileStorageInterface, searchWeights *config.SearchWeights) (*Database, error) {
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
		b2cProductIndex:  b2cIndexName, // Индекс для B2C товаров из конфигурации
		fsStorage:        fileStorage,  // Используем переданный параметр
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
		// Используем отдельный индекс b2c_products для B2C товаров
		db.productSearchRepo = storefrontOpenSearch.NewProductRepository(osClient, db.b2cProductIndex)
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
