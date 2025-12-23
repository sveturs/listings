// Package indexer provides indexing services for OpenSearch integration.
package indexer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/rs/zerolog"

	"github.com/vondi-global/listings/internal/domain"
	"github.com/vondi-global/listings/internal/repository/opensearch"
	"github.com/vondi-global/listings/internal/repository/postgres"
)

// ListingIndexer handles listing indexing operations combining data from DB and OpenSearch
type ListingIndexer struct {
	db               *sqlx.DB
	osClient         *opensearch.Client
	attributeIndexer *AttributeIndexer
	variantRepo      *postgres.VariantRepository
	logger           zerolog.Logger
}

// NewListingIndexer creates a new ListingIndexer instance
func NewListingIndexer(db *sqlx.DB, osClient *opensearch.Client, logger zerolog.Logger) *ListingIndexer {
	return &ListingIndexer{
		db:               db,
		osClient:         osClient,
		attributeIndexer: NewAttributeIndexer(db, logger),
		variantRepo:      postgres.NewVariantRepository(db, logger),
		logger:           logger.With().Str("component", "listing_indexer").Logger(),
	}
}

// IndexListing indexes a listing with its attributes in OpenSearch
func (idx *ListingIndexer) IndexListing(ctx context.Context, listing *domain.Listing) error {
	if listing == nil {
		return fmt.Errorf("listing cannot be nil")
	}

	// Get attributes from cache
	attributes, searchableText, err := idx.getAttributesFromCache(ctx, int32(listing.ID))
	if err != nil {
		// Log but continue - attributes are optional
		idx.logger.Debug().Err(err).Int64("listing_id", listing.ID).Msg("no attributes cache found")
	}

	// Build document with attributes
	doc := idx.buildListingDocument(listing, attributes, searchableText)

	// Index in OpenSearch using bulk operation structure
	body, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	// Use underlying OpenSearch client
	osClient := idx.osClient.GetClient()
	res, err := osClient.Index(
		"marketplace_listings", // Index name
		bytes.NewReader(body),
		osClient.Index.WithContext(ctx),
		osClient.Index.WithDocumentID(fmt.Sprintf("%d", listing.ID)),
		osClient.Index.WithRefresh("false"),
	)

	if err != nil {
		idx.logger.Error().Err(err).Int64("listing_id", listing.ID).Msg("failed to index listing")
		return fmt.Errorf("failed to index listing: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		idx.logger.Error().Int("status", res.StatusCode).Int64("listing_id", listing.ID).Msg("OpenSearch index error")
		return fmt.Errorf("OpenSearch index error: %s", res.Status())
	}

	idx.logger.Debug().Int64("listing_id", listing.ID).Msg("listing indexed successfully")
	return nil
}

// BulkIndexListings indexes multiple listings with attributes in a single bulk request
func (idx *ListingIndexer) BulkIndexListings(ctx context.Context, listings []*domain.Listing) error {
	if len(listings) == 0 {
		return nil
	}

	idx.logger.Info().Int("count", len(listings)).Msg("bulk indexing listings with attributes")

	// Build bulk request body
	var bulkBody bytes.Buffer
	successCount := 0

	for _, listing := range listings {
		// Get attributes from cache
		attributes, searchableText, err := idx.getAttributesFromCache(ctx, int32(listing.ID))
		if err != nil {
			idx.logger.Debug().Err(err).Int64("listing_id", listing.ID).Msg("no attributes cache for listing")
			// Continue without attributes
		}

		// Build document
		doc := idx.buildListingDocument(listing, attributes, searchableText)

		// Action line
		action := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": "marketplace_listings",
				"_id":    fmt.Sprintf("%d", listing.ID),
			},
		}

		actionJSON, err := json.Marshal(action)
		if err != nil {
			idx.logger.Error().Err(err).Int64("listing_id", listing.ID).Msg("failed to marshal action")
			continue
		}
		bulkBody.Write(actionJSON)
		bulkBody.WriteByte('\n')

		// Document line
		docJSON, err := json.Marshal(doc)
		if err != nil {
			idx.logger.Error().Err(err).Int64("listing_id", listing.ID).Msg("failed to marshal document")
			continue
		}
		bulkBody.Write(docJSON)
		bulkBody.WriteByte('\n')

		successCount++
	}

	if successCount == 0 {
		return fmt.Errorf("no listings to index")
	}

	// Execute bulk request
	osClient := idx.osClient.GetClient()
	res, err := osClient.Bulk(
		bytes.NewReader(bulkBody.Bytes()),
		osClient.Bulk.WithContext(ctx),
	)

	if err != nil {
		idx.logger.Error().Err(err).Int("listing_count", successCount).Msg("bulk index request failed")
		return fmt.Errorf("bulk index request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		idx.logger.Error().Int("status", res.StatusCode).Str("body", string(body)).Msg("bulk index response error")
		return fmt.Errorf("bulk index response error: %s", res.Status())
	}

	// Parse response
	var bulkResp map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&bulkResp); err != nil {
		idx.logger.Error().Err(err).Msg("failed to parse bulk response")
		return fmt.Errorf("failed to parse bulk response: %w", err)
	}

	// Check for errors
	if errors, ok := bulkResp["errors"].(bool); ok && errors {
		idx.logger.Warn().Msg("some items in bulk request failed")
		return fmt.Errorf("some items in bulk request failed (see logs)")
	}

	idx.logger.Info().Int("count", successCount).Msg("bulk index completed successfully")
	return nil
}

// DeleteListing removes a listing from OpenSearch index
func (idx *ListingIndexer) DeleteListing(ctx context.Context, listingID int64) error {
	return idx.osClient.DeleteListing(ctx, listingID)
}

// buildListingDocument builds an OpenSearch document from listing and attributes
func (idx *ListingIndexer) buildListingDocument(listing *domain.Listing, attributes []AttributeForIndex, searchableText string) map[string]interface{} {
	doc := map[string]interface{}{
		"id":              listing.ID,
		"uuid":            listing.UUID,
		"user_id":         listing.UserID,
		"title":           listing.Title,
		"price":           listing.Price,
		"currency":        listing.Currency,
		"category_id":     listing.CategoryID,
		"status":          listing.Status,
		"visibility":      listing.Visibility,
		"quantity":        listing.Quantity,
		"source_type":     listing.SourceType,
		"document_type":   "listing",
		"views_count":     listing.ViewsCount,
		"favorites_count": listing.FavoritesCount,
		"created_at":      listing.CreatedAt,
		"updated_at":      listing.UpdatedAt,

		// New fields - popularity and flags
		"popularity_score": calculatePopularityScore(listing),
		"is_new_arrival":   isNewArrival(listing.CreatedAt),
		"is_promoted":      false, // TODO: implement promotion logic
		"is_featured":      false, // TODO: implement featured logic
	}

	// Add category_slug if available
	if listing.CategorySlug != nil {
		doc["category_slug"] = *listing.CategorySlug
	}

	// Optional fields
	if listing.Description != nil {
		doc["description"] = *listing.Description
	}

	// Add translations for multilingual search
	// Supported languages: sr (Serbian), en (English), ru (Russian)
	if listing.OriginalLanguage != "" {
		doc["original_language"] = listing.OriginalLanguage
	}
	if len(listing.TitleTranslations) > 0 {
		for lang, translation := range listing.TitleTranslations {
			if translation != "" {
				doc["title_"+lang] = translation
			}
		}
	}
	if len(listing.DescriptionTranslations) > 0 {
		for lang, translation := range listing.DescriptionTranslations {
			if translation != "" {
				doc["description_"+lang] = translation
			}
		}
	}
	if listing.StorefrontID != nil {
		doc["storefront_id"] = *listing.StorefrontID
	}
	if listing.SKU != nil {
		doc["sku"] = *listing.SKU
	}
	if listing.StockStatus != nil {
		doc["stock_status"] = *listing.StockStatus
	}
	if listing.PublishedAt != nil {
		doc["published_at"] = *listing.PublishedAt
	}

	// Add images
	if len(listing.Images) > 0 {
		images := make([]map[string]interface{}, 0, len(listing.Images))
		for _, img := range listing.Images {
			images = append(images, map[string]interface{}{
				"id":         img.ID,
				"public_url": img.URL,
				"file_path":  img.URL,
				"is_main":    img.IsPrimary,
			})
		}
		doc["images"] = images
	}

	// Add location
	if listing.Location != nil {
		loc := listing.Location
		if loc.Latitude != nil && loc.Longitude != nil {
			doc["has_individual_location"] = true
			doc["individual_latitude"] = *loc.Latitude
			doc["individual_longitude"] = *loc.Longitude
			doc["location"] = map[string]interface{}{
				"lat": *loc.Latitude,
				"lon": *loc.Longitude,
			}
		}
		if loc.Country != nil {
			doc["country"] = *loc.Country
		}
		if loc.City != nil {
			doc["city"] = *loc.City
		}
	}

	// Add attributes
	if len(attributes) > 0 {
		doc["attributes"] = attributes
		doc["attributes_searchable_text"] = searchableText

		// Extract brand from attributes
		brand := extractBrand(attributes)
		if brand != "" {
			doc["brand"] = brand
		}
	}

	// Add tags if available
	if len(listing.Tags) > 0 {
		doc["tags"] = unique(listing.Tags)
	}

	// Add location translations (if available)
	if listing.Location != nil {
		if len(listing.CountryTranslations) > 0 {
			for lang, translation := range listing.CountryTranslations {
				if translation != "" {
					doc["country_"+lang] = translation
				}
			}
		}
		if len(listing.CityTranslations) > 0 {
			for lang, translation := range listing.CityTranslations {
				if translation != "" {
					doc["city_"+lang] = translation
				}
			}
		}
		if len(listing.LocationTranslations) > 0 {
			for lang, translation := range listing.LocationTranslations {
				if translation != "" {
					doc["address_"+lang] = translation
				}
			}
		}
	}

	// Load and add storefront data for B2C listings
	if listing.StorefrontID != nil && *listing.StorefrontID > 0 {
		storefrontData, err := idx.loadStorefrontData(context.Background(), *listing.StorefrontID)
		if err == nil {
			doc["storefront_name"] = storefrontData.Name
			doc["storefront_slug"] = storefrontData.Slug
			doc["storefront_rating"] = storefrontData.Rating
			doc["seller_verified"] = storefrontData.IsVerified
		} else {
			idx.logger.Debug().Err(err).Int64("storefront_id", *listing.StorefrontID).Msg("failed to load storefront data")
		}
	}

	// Load and add variants for B2C products
	if listing.SourceType == domain.SourceTypeB2C {
		variants, err := idx.loadVariantsForListing(context.Background(), listing)
		if err == nil && len(variants) > 0 {
			doc["has_variants"] = true
			doc["variants_count"] = len(variants)

			// Collect variant data
			variantSKUs := make([]string, 0, len(variants))
			variantColors := make([]string, 0)
			variantSizes := make([]string, 0)
			prices := make([]float64, 0, len(variants))
			totalStock := int32(0)

			variantsNested := make([]map[string]interface{}, 0, len(variants))

			for _, v := range variants {
				if v.SKU != "" {
					variantSKUs = append(variantSKUs, v.SKU)
				}

				// Extract color and size from attributes
				for _, attr := range v.Attributes {
					if attr.ValueText != nil {
						// Assuming attribute_id mappings: 1=color, 2=size (adjust if needed)
						switch attr.AttributeID {
						case 1: // color
							variantColors = append(variantColors, *attr.ValueText)
						case 2: // size
							variantSizes = append(variantSizes, *attr.ValueText)
						}
					}
				}

				// Calculate available stock
				availableStock := v.StockQuantity - v.ReservedQuantity
				totalStock += availableStock

				if v.Price != nil {
					prices = append(prices, *v.Price)
				}

				// Build nested variant object
				variantDoc := map[string]interface{}{
					"sku":   v.SKU,
					"stock": availableStock,
				}
				if v.Price != nil {
					variantDoc["price"] = *v.Price
				}
				if v.CompareAtPrice != nil {
					variantDoc["compare_at_price"] = *v.CompareAtPrice
				}

				// Add attributes to variant
				variantAttrs := make(map[string]interface{})
				for _, attr := range v.Attributes {
					if attr.ValueText != nil {
						// Use attribute code from DB (requires join) or use ID as key
						variantAttrs[fmt.Sprintf("attr_%d", attr.AttributeID)] = *attr.ValueText
					} else if attr.ValueNumber != nil {
						variantAttrs[fmt.Sprintf("attr_%d", attr.AttributeID)] = *attr.ValueNumber
					}
				}
				if len(variantAttrs) > 0 {
					variantDoc["attributes"] = variantAttrs
				}

				variantsNested = append(variantsNested, variantDoc)
			}

			doc["variant_skus"] = unique(variantSKUs)
			if len(variantColors) > 0 {
				doc["variant_colors"] = unique(variantColors)
			}
			if len(variantSizes) > 0 {
				doc["variant_sizes"] = unique(variantSizes)
			}
			doc["total_stock"] = totalStock

			// Calculate min/max price from variants
			if len(prices) > 0 {
				minPrice := prices[0]
				maxPrice := prices[0]
				for _, p := range prices {
					if p < minPrice {
						minPrice = p
					}
					if p > maxPrice {
						maxPrice = p
					}
				}
				doc["min_price"] = minPrice
				doc["max_price"] = maxPrice
			}

			// Add nested variants array
			doc["variants"] = variantsNested
		} else if err != nil {
			idx.logger.Debug().Err(err).Int64("listing_id", listing.ID).Msg("failed to load variants")
		}
	}

	// TODO: Add shipping, discount, rating, condition, etc. when available in domain.Listing

	return doc
}

// getAttributesFromCache retrieves cached attributes for a listing
func (idx *ListingIndexer) getAttributesFromCache(ctx context.Context, listingID int32) ([]AttributeForIndex, string, error) {
	query := `
		SELECT attributes_flat, attributes_searchable
		FROM attribute_search_cache
		WHERE listing_id = $1
	`

	var attributesJSON []byte
	var searchableText *string

	err := idx.db.QueryRowContext(ctx, query, listingID).Scan(&attributesJSON, &searchableText)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get cache: %w", err)
	}

	var attributes []AttributeForIndex
	if len(attributesJSON) > 0 {
		if err := json.Unmarshal(attributesJSON, &attributes); err != nil {
			return nil, "", fmt.Errorf("failed to unmarshal attributes: %w", err)
		}
	}

	searchText := ""
	if searchableText != nil {
		searchText = *searchableText
	}

	return attributes, searchText, nil
}

// ReindexAllWithAttributes reindexes all listings with their attributes
func (idx *ListingIndexer) ReindexAllWithAttributes(ctx context.Context, batchSize int) error {
	if batchSize <= 0 {
		batchSize = 500 // Increased from 100 to 500
	}

	idx.logger.Info().Int("batch_size", batchSize).Msg("starting full reindex with attributes")

	// Get all active listings with category slug and translations
	// Note: category_id is UUID in DB but we don't use it for indexing (use category_slug instead)
	query := `
		SELECT l.id, l.uuid, l.user_id, l.storefront_id, l.title, l.description,
		       l.price, l.currency, l.category_id::text, c.slug AS category_slug,
		       l.status, l.visibility, l.quantity, l.sku,
		       l.source_type, l.stock_status, l.view_count, l.favorites_count,
		       l.created_at, l.updated_at, l.published_at,
		       COALESCE(l.title_translations, '{}'::jsonb) AS title_translations,
		       COALESCE(l.description_translations, '{}'::jsonb) AS description_translations,
		       COALESCE(l.original_language, 'sr') AS original_language
		FROM listings l
		LEFT JOIN categories c ON l.category_id = c.id
		WHERE l.status = 'active' AND l.visibility = 'public' AND l.is_deleted = false
		ORDER BY l.id ASC
	`

	rows, err := idx.db.QueryContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to query listings: %w", err)
	}
	defer rows.Close()

	var listings []*domain.Listing
	totalProcessed := 0

	for rows.Next() {
		var listing domain.Listing
		var titleTransJSON, descTransJSON []byte

		err := rows.Scan(
			&listing.ID,
			&listing.UUID,
			&listing.UserID,
			&listing.StorefrontID,
			&listing.Title,
			&listing.Description,
			&listing.Price,
			&listing.Currency,
			&listing.CategoryID,
			&listing.CategorySlug,
			&listing.Status,
			&listing.Visibility,
			&listing.Quantity,
			&listing.SKU,
			&listing.SourceType,
			&listing.StockStatus,
			&listing.ViewsCount,
			&listing.FavoritesCount,
			&listing.CreatedAt,
			&listing.UpdatedAt,
			&listing.PublishedAt,
			&titleTransJSON,
			&descTransJSON,
			&listing.OriginalLanguage,
		)
		if err != nil {
			idx.logger.Error().Err(err).Msg("failed to scan listing")
			continue
		}

		// Parse translations from JSONB
		if len(titleTransJSON) > 0 {
			if err := json.Unmarshal(titleTransJSON, &listing.TitleTranslations); err != nil {
				idx.logger.Debug().Err(err).Int64("listing_id", listing.ID).Msg("failed to parse title translations")
			}
		}
		if len(descTransJSON) > 0 {
			if err := json.Unmarshal(descTransJSON, &listing.DescriptionTranslations); err != nil {
				idx.logger.Debug().Err(err).Int64("listing_id", listing.ID).Msg("failed to parse description translations")
			}
		}

		listings = append(listings, &listing)

		// Process batch when full
		if len(listings) >= batchSize {
			// Load images for this batch
			if err := idx.loadImagesForListings(ctx, listings); err != nil {
				idx.logger.Warn().Err(err).Msg("failed to load images for batch")
			}

			if err := idx.BulkIndexListings(ctx, listings); err != nil {
				idx.logger.Error().Err(err).Msg("failed to index batch")
			} else {
				totalProcessed += len(listings)
			}
			listings = nil // Reset batch
		}
	}

	// Process remaining listings
	if len(listings) > 0 {
		// Load images for final batch
		if err := idx.loadImagesForListings(ctx, listings); err != nil {
			idx.logger.Warn().Err(err).Msg("failed to load images for final batch")
		}

		if err := idx.BulkIndexListings(ctx, listings); err != nil {
			idx.logger.Error().Err(err).Msg("failed to index final batch")
		} else {
			totalProcessed += len(listings)
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating listings: %w", err)
	}

	idx.logger.Info().Int("total_processed", totalProcessed).Msg("reindex completed")
	return nil
}

// loadImagesForListings loads images for a batch of listings
func (idx *ListingIndexer) loadImagesForListings(ctx context.Context, listings []*domain.Listing) error {
	if len(listings) == 0 {
		return nil
	}

	// Collect all listing IDs
	ids := make([]int64, len(listings))
	listingMap := make(map[int64]*domain.Listing, len(listings))
	for i, l := range listings {
		ids[i] = l.ID
		listingMap[l.ID] = l
	}

	// Query all images for these listings in one batch
	query := `
		SELECT id, listing_id, url, storage_path, thumbnail_url,
		       display_order, is_primary, width, height, file_size, mime_type,
		       created_at, updated_at
		FROM listing_images
		WHERE listing_id = ANY($1)
		ORDER BY listing_id, is_primary DESC, display_order ASC, id ASC
	`

	rows, err := idx.db.QueryContext(ctx, query, pq.Array(ids))
	if err != nil {
		return fmt.Errorf("failed to query images: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var img domain.ListingImage
		err := rows.Scan(
			&img.ID,
			&img.ListingID,
			&img.URL,
			&img.StoragePath,
			&img.ThumbnailURL,
			&img.DisplayOrder,
			&img.IsPrimary,
			&img.Width,
			&img.Height,
			&img.FileSize,
			&img.MimeType,
			&img.CreatedAt,
			&img.UpdatedAt,
		)
		if err != nil {
			idx.logger.Warn().Err(err).Msg("failed to scan image")
			continue
		}

		// Add image to corresponding listing
		if listing, ok := listingMap[img.ListingID]; ok {
			listing.Images = append(listing.Images, &img)
		}
	}

	return rows.Err()
}

// loadVariantsForListing loads product variants for a B2C listing
func (idx *ListingIndexer) loadVariantsForListing(ctx context.Context, listing *domain.Listing) ([]*domain.ProductVariantV2, error) {
	// Parse listing UUID to use as product_id
	productID, err := uuid.Parse(listing.UUID)
	if err != nil {
		return nil, fmt.Errorf("invalid listing UUID: %w", err)
	}

	filter := &domain.ListVariantsFilter{
		ProductID:         productID,
		ActiveOnly:        true,
		IncludeAttributes: true,
	}

	variants, err := idx.variantRepo.ListByProduct(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to load variants: %w", err)
	}

	return variants, nil
}

// extractBrand extracts brand from listing attributes
func extractBrand(attrs []AttributeForIndex) string {
	for _, attr := range attrs {
		if strings.ToLower(attr.Code) == "brand" && attr.ValueText != nil {
			return *attr.ValueText
		}
	}
	return ""
}

// calculatePopularityScore calculates popularity score using weighted formula
func calculatePopularityScore(listing *domain.Listing) float64 {
	// popularity_score = views*0.3 + favorites*0.5 + orders*0.2
	// Note: orders count не доступен в listing entity, используем 0
	views := float64(listing.ViewsCount)
	favorites := float64(listing.FavoritesCount)
	return views*0.3 + favorites*0.5
}

// isNewArrival checks if listing was created within last 7 days
func isNewArrival(createdAt time.Time) bool {
	return time.Since(createdAt) < 7*24*time.Hour
}

// unique removes duplicates from string slice
func unique(slice []string) []string {
	seen := make(map[string]bool)
	result := make([]string, 0, len(slice))
	for _, s := range slice {
		if s != "" && !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	return result
}

// loadStorefrontData loads storefront information for B2C listings
func (idx *ListingIndexer) loadStorefrontData(ctx context.Context, storefrontID int64) (*StorefrontData, error) {
	query := `
		SELECT id, slug, name, rating, is_verified
		FROM storefronts
		WHERE id = $1 AND deleted_at IS NULL
	`

	var data StorefrontData
	err := idx.db.QueryRowContext(ctx, query, storefrontID).Scan(
		&data.ID,
		&data.Slug,
		&data.Name,
		&data.Rating,
		&data.IsVerified,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load storefront data: %w", err)
	}

	return &data, nil
}

// StorefrontData holds minimal storefront info for indexing
type StorefrontData struct {
	ID         int64
	Slug       string
	Name       string
	Rating     float64
	IsVerified bool
}

// CountActiveListings returns the count of active public listings in database
func (idx *ListingIndexer) CountActiveListings(ctx context.Context) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM listings
		WHERE status = 'active' AND visibility = 'public' AND is_deleted = false
	`

	var count int64
	err := idx.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count active listings: %w", err)
	}

	return count, nil
}

// GetDB returns the database connection (needed for reindexer)
func (idx *ListingIndexer) GetDB() *sqlx.DB {
	return idx.db
}
