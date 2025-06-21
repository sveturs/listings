package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"backend/internal/config"
	"backend/internal/logger"
	"backend/internal/proj/storefronts/storage/opensearch"
	osClient "backend/internal/storage/opensearch"
	"backend/internal/storage/postgres"
	
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Флаги командной строки
	var (
		createIndex = flag.Bool("create", false, "Create new index before reindexing")
		batchSize   = flag.Int("batch", 100, "Batch size for indexing")
		indexName   = flag.String("index", "storefront_products", "OpenSearch index name")
	)
	flag.Parse()

	// Загрузка конфигурации
	cfg := config.MustLoad()

	// Подключение к PostgreSQL
	dbConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to parse database URL: %v", err)
	}

	dbPool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Подключение к OpenSearch
	osClientInstance := osClient.NewClient(cfg.OpenSearch)
	if err := osClientInstance.Ping(); err != nil {
		log.Fatalf("Failed to connect to OpenSearch: %v", err)
	}

	// Создание репозитория для товаров витрин
	productRepo := opensearch.NewProductRepository(osClientInstance, *indexName)

	// Создание индекса если указано
	if *createIndex {
		logger.Info().Msg("Creating/updating OpenSearch index...")
		if err := productRepo.PrepareIndex(context.Background()); err != nil {
			log.Fatalf("Failed to prepare index: %v", err)
		}
		logger.Info().Msg("Index prepared successfully")
	}

	// Получение общего количества товаров
	var totalCount int
	err = dbPool.QueryRow(context.Background(), `
		SELECT COUNT(*) 
		FROM storefront_products 
		WHERE status = 'active'
	`).Scan(&totalCount)
	if err != nil {
		log.Fatalf("Failed to count products: %v", err)
	}

	logger.Info().Msgf("Found %d active products to index", totalCount)

	// Индексация товаров батчами
	offset := 0
	indexed := 0
	failed := 0
	startTime := time.Now()

	for offset < totalCount {
		// Получение батча товаров
		rows, err := dbPool.Query(context.Background(), `
			SELECT 
				p.id,
				p.storefront_id,
				p.category_id,
				p.name,
				p.description,
				p.price,
				p.currency,
				p.sku,
				p.barcode,
				p.status,
				p.metadata,
				p.attributes,
				p.inventory_count,
				p.reserved_count,
				p.low_stock_threshold,
				p.track_inventory,
				p.sales_count,
				p.views_count,
				p.created_at,
				p.updated_at,
				-- Storefront info
				s.id as sf_id,
				s.name as sf_name,
				s.slug as sf_slug,
				s.city as sf_city,
				s.country as sf_country,
				s.latitude as sf_latitude,
				s.longitude as sf_longitude,
				s.rating as sf_rating,
				s.is_verified as sf_is_verified,
				-- Category info
				c.id as cat_id,
				c.name as cat_name,
				c.slug as cat_slug,
				c.path as cat_path
			FROM storefront_products p
			LEFT JOIN user_storefronts s ON p.storefront_id = s.id
			LEFT JOIN marketplace_categories c ON p.category_id = c.id
			WHERE p.status = 'active'
			ORDER BY p.id
			LIMIT $1 OFFSET $2
		`, *batchSize, offset)
		
		if err != nil {
			log.Printf("Failed to fetch products batch at offset %d: %v", offset, err)
			break
		}
		defer rows.Close()

		// Парсинг и индексация батча
		batch := make([]*models.StorefrontProduct, 0, *batchSize)
		
		for rows.Next() {
			product := &models.StorefrontProduct{}
			storefront := &models.Storefront{}
			category := &models.MarketplaceCategory{}
			
			var sfLatitude, sfLongitude *float64
			var catPath *string
			
			err := rows.Scan(
				&product.ID,
				&product.StorefrontID,
				&product.CategoryID,
				&product.Name,
				&product.Description,
				&product.Price,
				&product.Currency,
				&product.SKU,
				&product.Barcode,
				&product.Status,
				&product.Metadata,
				&product.Attributes,
				&product.InventoryCount,
				&product.ReservedCount,
				&product.LowStockThreshold,
				&product.TrackInventory,
				&product.SalesCount,
				&product.ViewsCount,
				&product.CreatedAt,
				&product.UpdatedAt,
				// Storefront
				&storefront.ID,
				&storefront.Name,
				&storefront.Slug,
				&storefront.City,
				&storefront.Country,
				&sfLatitude,
				&sfLongitude,
				&storefront.Rating,
				&storefront.IsVerified,
				// Category
				&category.ID,
				&category.Name,
				&category.Slug,
				&catPath,
			)
			
			if err != nil {
				log.Printf("Failed to scan product: %v", err)
				failed++
				continue
			}
			
			// Установка связей
			if storefront.ID > 0 {
				storefront.Latitude = sfLatitude
				storefront.Longitude = sfLongitude
				product.Storefront = storefront
			}
			
			if category.ID > 0 {
				if catPath != nil {
					category.Path = *catPath
				}
				product.Category = category
			}
			
			// Получение изображений
			imageRows, err := dbPool.Query(context.Background(), `
				SELECT id, url, alt_text, is_main, position
				FROM storefront_product_images
				WHERE product_id = $1
				ORDER BY position, id
			`, product.ID)
			
			if err == nil {
				defer imageRows.Close()
				
				for imageRows.Next() {
					img := &models.StorefrontProductImage{}
					err := imageRows.Scan(
						&img.ID,
						&img.URL,
						&img.AltText,
						&img.IsMain,
						&img.Position,
					)
					if err == nil {
						product.Images = append(product.Images, img)
					}
				}
			}
			
			// Получение вариантов
			variantRows, err := dbPool.Query(context.Background(), `
				SELECT id, name, sku, price, attributes
				FROM storefront_product_variants
				WHERE product_id = $1
				ORDER BY id
			`, product.ID)
			
			if err == nil {
				defer variantRows.Close()
				
				for variantRows.Next() {
					variant := &models.StorefrontProductVariant{}
					err := variantRows.Scan(
						&variant.ID,
						&variant.Name,
						&variant.SKU,
						&variant.Price,
						&variant.Attributes,
					)
					if err == nil {
						product.Variants = append(product.Variants, variant)
					}
				}
			}
			
			batch = append(batch, product)
		}
		
		if err := rows.Err(); err != nil {
			log.Printf("Error iterating rows: %v", err)
		}

		// Индексация батча
		if len(batch) > 0 {
			if err := productRepo.BulkIndexProducts(context.Background(), batch); err != nil {
				log.Printf("Failed to index batch at offset %d: %v", offset, err)
				failed += len(batch)
			} else {
				indexed += len(batch)
				logger.Info().Msgf("Indexed batch: %d products (total: %d/%d)", len(batch), indexed, totalCount)
			}
		}

		offset += *batchSize
		
		// Прогресс
		progress := float64(offset) / float64(totalCount) * 100
		elapsed := time.Since(startTime)
		eta := time.Duration(float64(elapsed) / float64(offset) * float64(totalCount-offset))
		
		fmt.Printf("\rProgress: %.1f%% (%d/%d) | Indexed: %d | Failed: %d | ETA: %s", 
			progress, offset, totalCount, indexed, failed, eta.Round(time.Second))
	}

	fmt.Println() // Новая строка после прогресса
	
	// Финальная статистика
	duration := time.Since(startTime)
	logger.Info().Msgf("Reindexing completed in %s", duration)
	logger.Info().Msgf("Total products: %d", totalCount)
	logger.Info().Msgf("Successfully indexed: %d", indexed)
	logger.Info().Msgf("Failed: %d", failed)
	logger.Info().Msgf("Average speed: %.2f products/second", float64(indexed)/duration.Seconds())
}

// Модели для индексации (упрощенные версии)
package models

type StorefrontProduct struct {
	ID                int                          `json:"id"`
	StorefrontID      int                          `json:"storefront_id"`
	CategoryID        int                          `json:"category_id"`
	Name              string                       `json:"name"`
	Description       string                       `json:"description"`
	Price             float64                      `json:"price"`
	Currency          string                       `json:"currency"`
	SKU               string                       `json:"sku"`
	Barcode           string                       `json:"barcode"`
	Status            string                       `json:"status"`
	Metadata          map[string]interface{}       `json:"metadata"`
	Attributes        []map[string]interface{}     `json:"attributes"`
	InventoryCount    int                          `json:"inventory_count"`
	ReservedCount     int                          `json:"reserved_count"`
	LowStockThreshold int                          `json:"low_stock_threshold"`
	TrackInventory    bool                         `json:"track_inventory"`
	SalesCount        int                          `json:"sales_count"`
	ViewsCount        int                          `json:"views_count"`
	CreatedAt         time.Time                    `json:"created_at"`
	UpdatedAt         time.Time                    `json:"updated_at"`
	
	// Relations
	Storefront        *Storefront                  `json:"storefront,omitempty"`
	Category          *MarketplaceCategory         `json:"category,omitempty"`
	Images            []*StorefrontProductImage    `json:"images,omitempty"`
	Variants          []*StorefrontProductVariant  `json:"variants,omitempty"`
}

type Storefront struct {
	ID         int      `json:"id"`
	Name       string   `json:"name"`
	Slug       string   `json:"slug"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	Latitude   *float64 `json:"latitude,omitempty"`
	Longitude  *float64 `json:"longitude,omitempty"`
	Rating     float64  `json:"rating"`
	IsVerified bool     `json:"is_verified"`
}

type MarketplaceCategory struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	Path string `json:"path"`
}

type StorefrontProductImage struct {
	ID       int    `json:"id"`
	URL      string `json:"url"`
	AltText  string `json:"alt_text"`
	IsMain   bool   `json:"is_main"`
	Position int    `json:"position"`
}

type StorefrontProductVariant struct {
	ID         int                    `json:"id"`
	Name       string                 `json:"name"`
	SKU        string                 `json:"sku"`
	Price      float64                `json:"price"`
	Attributes map[string]interface{} `json:"attributes"`
}