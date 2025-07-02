package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"backend/internal/config"
	"backend/internal/domain/models"
	"backend/internal/logger"
	"backend/internal/proj/storefronts/storage/opensearch"
	osClientPkg "backend/internal/storage/opensearch"
	postgresStorage "backend/internal/storage/postgres"

	"github.com/joho/godotenv"
)

func main() {
	// Parse command line flags
	var (
		batchSize = flag.Int("batch", 100, "Batch size for indexing")
		recreate  = flag.Bool("recreate", false, "Recreate index before reindexing")
	)
	flag.Parse()

	// Load environment variables
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}
	if err := godotenv.Load(envFile); err != nil {
		fmt.Printf("Warning: Could not load .env file: %s\n", err)
	}

	// Initialize logger
	if err := logger.Init(os.Getenv("APP_MODE"), os.Getenv("LOG_LEVEL")); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	// Initialize configuration
	cfg, err := config.NewConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load config")
	}

	// Initialize OpenSearch client
	osClient, err := osClientPkg.NewOpenSearchClient(osClientPkg.Config{
		URL:      cfg.OpenSearch.URL,
		Username: cfg.OpenSearch.Username,
		Password: cfg.OpenSearch.Password,
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create OpenSearch client")
	}

	// Initialize storage with all necessary dependencies
	storage, err := postgresStorage.NewDatabase(
		cfg.DatabaseURL,
		osClient,
		"",  // minioEndpoint - not needed for this task
		nil, // fileStorage - not needed for this task
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize storage")
	}
	defer storage.Close()

	// Get storefront product search repository
	searchRepo := storage.StorefrontProductSearch()
	if searchRepo == nil {
		logger.Fatal().Msg("Storefront product search repository not configured")
	}

	productSearchRepo, ok := searchRepo.(opensearch.ProductSearchRepository)
	if !ok {
		logger.Fatal().Msg("Invalid storefront product search repository type")
	}

	ctx := context.Background()

	// If recreate flag is set, recreate the index
	if *recreate {
		logger.Info().Msg("Recreating index...")
		if err := productSearchRepo.PrepareIndex(ctx); err != nil {
			logger.Fatal().Err(err).Msg("Failed to prepare index")
		}
	}

	// Get all storefront products
	logger.Info().Msg("Fetching storefront products...")

	// We need to access products directly from the database
	query := `
		SELECT 
			p.id, p.storefront_id, p.category_id, p.name, p.description,
			p.price, p.currency, p.sku, p.barcode, p.stock_quantity,
			p.stock_status, p.is_active, p.attributes, p.created_at, p.updated_at,
			p.sold_count, p.view_count,
			c.id as cat_id, c.name as cat_name, c.slug as cat_slug
		FROM storefront_products p
		LEFT JOIN marketplace_categories c ON c.id = p.category_id
		WHERE p.is_active = true
		ORDER BY p.id
	`

	rows, err := storage.Query(ctx, query)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to fetch products")
	}
	defer rows.Close()

	// Process products in batches
	var products []*models.StorefrontProduct
	totalCount := 0
	indexedCount := 0

	for rows.Next() {
		var p models.StorefrontProduct
		var categoryID, catID *int
		var catName, catSlug *string

		err := rows.Scan(
			&p.ID, &p.StorefrontID, &categoryID, &p.Name, &p.Description,
			&p.Price, &p.Currency, &p.SKU, &p.Barcode, &p.StockQuantity,
			&p.StockStatus, &p.IsActive, &p.Attributes, &p.CreatedAt, &p.UpdatedAt,
			&p.SoldCount, &p.ViewCount,
			&catID, &catName, &catSlug,
		)
		if err != nil {
			logger.Error().Err(err).Msg("Error scanning product")
			continue
		}

		if categoryID != nil {
			p.CategoryID = *categoryID
		}

		// Set category info if available
		if catID != nil && catName != nil {
			p.Category = &models.MarketplaceCategory{
				ID:   *catID,
				Name: *catName,
			}
			if catSlug != nil {
				p.Category.Slug = *catSlug
			}
		}

		// Get product images
		imgQuery := `
			SELECT id, image_url, thumbnail_url, alt_text, is_default, display_order
			FROM storefront_product_images
			WHERE product_id = $1
			ORDER BY display_order, id
		`
		imgRows, err := storage.Query(ctx, imgQuery, p.ID)
		if err == nil {
			for imgRows.Next() {
				var img models.StorefrontProductImage
				var altText *string
				err := imgRows.Scan(&img.ID, &img.ImageURL, &img.ThumbnailURL, &altText, &img.IsDefault, &img.DisplayOrder)
				if err == nil {
					p.Images = append(p.Images, img)
				}
			}
			imgRows.Close()
		}

		// Get product variants
		varQuery := `
			SELECT id, name, sku, price, stock_quantity, attributes
			FROM storefront_product_variants
			WHERE product_id = $1
			ORDER BY id
		`
		varRows, err := storage.Query(ctx, varQuery, p.ID)
		if err == nil {
			for varRows.Next() {
				var v models.StorefrontProductVariant
				err := varRows.Scan(&v.ID, &v.Name, &v.SKU, &v.Price, &v.StockQuantity, &v.Attributes)
				if err == nil {
					p.Variants = append(p.Variants, v)
				}
			}
			varRows.Close()
		}

		products = append(products, &p)
		totalCount++

		// Index batch when reaching batch size
		if len(products) >= *batchSize {
			logger.Info().Msgf("Indexing batch of %d products...", len(products))
			if err := productSearchRepo.BulkIndexProducts(ctx, products); err != nil {
				logger.Error().Err(err).Msg("Error indexing batch")
			} else {
				indexedCount += len(products)
			}
			products = nil
		}
	}

	// Index remaining products
	if len(products) > 0 {
		logger.Info().Msgf("Indexing final batch of %d products...", len(products))
		if err := productSearchRepo.BulkIndexProducts(ctx, products); err != nil {
			logger.Error().Err(err).Msg("Error indexing final batch")
		} else {
			indexedCount += len(products)
		}
	}

	logger.Info().
		Int("total", totalCount).
		Int("indexed", indexedCount).
		Msg("Reindexing completed!")

	if indexedCount < totalCount {
		logger.Error().Msgf("WARNING: %d products failed to index", totalCount-indexedCount)
		os.Exit(1)
	}
}
