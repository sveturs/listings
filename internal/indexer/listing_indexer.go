// Package indexer provides indexing services for OpenSearch integration.
package indexer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"

	"github.com/sveturs/listings/internal/domain"
	"github.com/sveturs/listings/internal/repository/opensearch"
)

// ListingIndexer handles listing indexing operations combining data from DB and OpenSearch
type ListingIndexer struct {
	db               *sqlx.DB
	osClient         *opensearch.Client
	attributeIndexer *AttributeIndexer
	logger           zerolog.Logger
}

// NewListingIndexer creates a new ListingIndexer instance
func NewListingIndexer(db *sqlx.DB, osClient *opensearch.Client, logger zerolog.Logger) *ListingIndexer {
	return &ListingIndexer{
		db:               db,
		osClient:         osClient,
		attributeIndexer: NewAttributeIndexer(db, logger),
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
	}

	// Optional fields
	if listing.Description != nil {
		doc["description"] = *listing.Description
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
	}

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
		batchSize = 100
	}

	idx.logger.Info().Int("batch_size", batchSize).Msg("starting full reindex with attributes")

	// Get all active listings
	query := `
		SELECT id, uuid, user_id, storefront_id, title, description,
		       price, currency, category_id, status, visibility, quantity, sku,
		       source_type, stock_status, view_count, favorites_count,
		       created_at, updated_at, published_at
		FROM listings
		WHERE status = 'active' AND visibility = 'public' AND is_deleted = false
		ORDER BY id ASC
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
		)
		if err != nil {
			idx.logger.Error().Err(err).Msg("failed to scan listing")
			continue
		}

		listings = append(listings, &listing)

		// Process batch when full
		if len(listings) >= batchSize {
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
