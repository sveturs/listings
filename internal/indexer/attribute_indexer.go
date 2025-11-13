// Package indexer provides indexing services for OpenSearch integration.
package indexer

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

// AttributeIndexer handles attribute indexing operations for OpenSearch
type AttributeIndexer struct {
	db     *sqlx.DB
	logger zerolog.Logger
}

// NewAttributeIndexer creates a new AttributeIndexer instance
func NewAttributeIndexer(db *sqlx.DB, logger zerolog.Logger) *AttributeIndexer {
	return &AttributeIndexer{
		db:     db,
		logger: logger.With().Str("component", "attribute_indexer").Logger(),
	}
}

// AttributeForIndex represents an attribute prepared for OpenSearch indexing
type AttributeForIndex struct {
	ID           int32    `json:"id"`
	Code         string   `json:"code"`
	Name         string   `json:"name"` // Denormalized for language
	ValueText    *string  `json:"value_text,omitempty"`
	ValueNumber  *float64 `json:"value_number,omitempty"`
	ValueBoolean *bool    `json:"value_boolean,omitempty"`
	IsSearchable bool     `json:"is_searchable"`
	IsFilterable bool     `json:"is_filterable"`
}

// PopulateAttributeSearchCache fills attribute_search_cache table for all listings with attributes
// This function processes listings in batches for better performance
func (idx *AttributeIndexer) PopulateAttributeSearchCache(ctx context.Context, batchSize int) error {
	if batchSize <= 0 {
		batchSize = 100 // Default batch size
	}

	idx.logger.Info().Int("batch_size", batchSize).Msg("starting attribute search cache population")

	// Get all listing IDs that have attributes
	listingIDs, err := idx.getListingIDsWithAttributes(ctx)
	if err != nil {
		return fmt.Errorf("failed to get listing IDs: %w", err)
	}

	idx.logger.Info().Int("total_listings", len(listingIDs)).Msg("found listings with attributes")

	// Process in batches
	totalProcessed := 0
	totalErrors := 0

	for i := 0; i < len(listingIDs); i += batchSize {
		end := i + batchSize
		if end > len(listingIDs) {
			end = len(listingIDs)
		}

		batch := listingIDs[i:end]
		idx.logger.Debug().Int("batch_start", i).Int("batch_size", len(batch)).Msg("processing batch")

		// Process batch
		for _, listingID := range batch {
			if err := idx.UpdateListingAttributeCache(ctx, listingID); err != nil {
				idx.logger.Error().Err(err).Int32("listing_id", listingID).Msg("failed to update cache for listing")
				totalErrors++
				continue
			}
			totalProcessed++
		}

		idx.logger.Info().
			Int("processed", totalProcessed).
			Int("errors", totalErrors).
			Int("total", len(listingIDs)).
			Msg("batch completed")
	}

	idx.logger.Info().
		Int("total_processed", totalProcessed).
		Int("total_errors", totalErrors).
		Msg("attribute search cache population completed")

	if totalErrors > 0 {
		return fmt.Errorf("completed with %d errors out of %d listings", totalErrors, len(listingIDs))
	}

	return nil
}

// UpdateListingAttributeCache updates the attribute_search_cache for a single listing
func (idx *AttributeIndexer) UpdateListingAttributeCache(ctx context.Context, listingID int32) error {
	// Build attributes data
	attributes, searchableText, filterableData, err := idx.BuildAttributesForIndex(ctx, listingID)
	if err != nil {
		return fmt.Errorf("failed to build attributes: %w", err)
	}

	// Marshal to JSON
	attributesJSON, err := json.Marshal(attributes)
	if err != nil {
		return fmt.Errorf("failed to marshal attributes: %w", err)
	}

	filterableJSON, err := json.Marshal(filterableData)
	if err != nil {
		return fmt.Errorf("failed to marshal filterable data: %w", err)
	}

	// Upsert into cache table
	query := `
		INSERT INTO attribute_search_cache (
			listing_id, attributes_flat, attributes_searchable, attributes_filterable, cache_version
		)
		VALUES ($1, $2, $3, $4, 1)
		ON CONFLICT (listing_id)
		DO UPDATE SET
			attributes_flat = EXCLUDED.attributes_flat,
			attributes_searchable = EXCLUDED.attributes_searchable,
			attributes_filterable = EXCLUDED.attributes_filterable,
			cache_version = EXCLUDED.cache_version,
			last_updated = CURRENT_TIMESTAMP
	`

	_, err = idx.db.ExecContext(ctx, query, listingID, attributesJSON, searchableText, filterableJSON)
	if err != nil {
		return fmt.Errorf("failed to upsert cache: %w", err)
	}

	idx.logger.Debug().Int32("listing_id", listingID).Msg("attribute cache updated")
	return nil
}

// BuildAttributesForIndex builds attribute data structures for OpenSearch indexing
// Returns: attributes array, searchable text, filterable data map, error
func (idx *AttributeIndexer) BuildAttributesForIndex(ctx context.Context, listingID int32) ([]AttributeForIndex, string, map[string]interface{}, error) {
	// Fetch listing attribute values with attribute metadata
	query := `
		SELECT
			a.id, a.code, a.name, a.is_searchable, a.is_filterable,
			lav.value_text, lav.value_number, lav.value_boolean
		FROM listing_attribute_values lav
		INNER JOIN attributes a ON lav.attribute_id = a.id
		WHERE lav.listing_id = $1 AND a.is_active = true
		ORDER BY a.sort_order ASC, (a.name->>'en') ASC
	`

	rows, err := idx.db.QueryContext(ctx, query, listingID)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to query attributes: %w", err)
	}
	defer rows.Close()

	var attributes []AttributeForIndex
	var searchableTexts []string
	filterableData := make(map[string]interface{})

	for rows.Next() {
		var attr AttributeForIndex
		var nameJSON []byte

		err := rows.Scan(
			&attr.ID,
			&attr.Code,
			&nameJSON,
			&attr.IsSearchable,
			&attr.IsFilterable,
			&attr.ValueText,
			&attr.ValueNumber,
			&attr.ValueBoolean,
		)
		if err != nil {
			return nil, "", nil, fmt.Errorf("failed to scan attribute: %w", err)
		}

		// Parse name JSON and get English name (or first available)
		var nameMap map[string]string
		if err := json.Unmarshal(nameJSON, &nameMap); err == nil {
			if enName, ok := nameMap["en"]; ok {
				attr.Name = enName
			} else if ruName, ok := nameMap["ru"]; ok {
				attr.Name = ruName
			} else {
				// Get first available name
				for _, name := range nameMap {
					attr.Name = name
					break
				}
			}
		}

		attributes = append(attributes, attr)

		// Build searchable text from searchable attributes
		if attr.IsSearchable {
			if attr.ValueText != nil && *attr.ValueText != "" {
				searchableTexts = append(searchableTexts, *attr.ValueText)
			} else if attr.ValueNumber != nil {
				searchableTexts = append(searchableTexts, fmt.Sprintf("%.2f", *attr.ValueNumber))
			} else if attr.ValueBoolean != nil {
				if *attr.ValueBoolean {
					searchableTexts = append(searchableTexts, "yes")
				} else {
					searchableTexts = append(searchableTexts, "no")
				}
			}
		}

		// Build filterable data structure
		if attr.IsFilterable {
			if attr.ValueText != nil {
				filterableData[attr.Code] = *attr.ValueText
			} else if attr.ValueNumber != nil {
				filterableData[attr.Code] = *attr.ValueNumber
			} else if attr.ValueBoolean != nil {
				filterableData[attr.Code] = *attr.ValueBoolean
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, "", nil, fmt.Errorf("error iterating attributes: %w", err)
	}

	searchableText := strings.Join(searchableTexts, " ")

	return attributes, searchableText, filterableData, nil
}

// GetListingAttributeCache retrieves cached attribute data for a listing
func (idx *AttributeIndexer) GetListingAttributeCache(ctx context.Context, listingID int32) ([]AttributeForIndex, error) {
	query := `
		SELECT attributes_flat
		FROM attribute_search_cache
		WHERE listing_id = $1
	`

	var attributesJSON []byte
	err := idx.db.QueryRowContext(ctx, query, listingID).Scan(&attributesJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to get cache: %w", err)
	}

	var attributes []AttributeForIndex
	if err := json.Unmarshal(attributesJSON, &attributes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal attributes: %w", err)
	}

	return attributes, nil
}

// DeleteListingAttributeCache removes cache entry for a listing
func (idx *AttributeIndexer) DeleteListingAttributeCache(ctx context.Context, listingID int32) error {
	query := `DELETE FROM attribute_search_cache WHERE listing_id = $1`

	_, err := idx.db.ExecContext(ctx, query, listingID)
	if err != nil {
		return fmt.Errorf("failed to delete cache: %w", err)
	}

	idx.logger.Debug().Int32("listing_id", listingID).Msg("attribute cache deleted")
	return nil
}

// getListingIDsWithAttributes retrieves all listing IDs that have attribute values
// Only returns IDs that exist in listings table (no orphans)
func (idx *AttributeIndexer) getListingIDsWithAttributes(ctx context.Context) ([]int32, error) {
	query := `
		SELECT DISTINCT lav.listing_id
		FROM listing_attribute_values lav
		INNER JOIN listings l ON lav.listing_id = l.id
		ORDER BY lav.listing_id ASC
	`

	rows, err := idx.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query listing IDs: %w", err)
	}
	defer rows.Close()

	var listingIDs []int32
	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan listing ID: %w", err)
		}
		listingIDs = append(listingIDs, id)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating listing IDs: %w", err)
	}

	return listingIDs, nil
}

// InvalidateCache invalidates cache for a specific listing (marks for reindex)
func (idx *AttributeIndexer) InvalidateCache(ctx context.Context, listingID int32) error {
	// For now, just delete the cache entry
	// In production, you might want to mark it as stale instead
	return idx.DeleteListingAttributeCache(ctx, listingID)
}

// BulkUpdateCache updates cache for multiple listings in a single transaction
func (idx *AttributeIndexer) BulkUpdateCache(ctx context.Context, listingIDs []int32) error {
	if len(listingIDs) == 0 {
		return nil
	}

	idx.logger.Info().Int("count", len(listingIDs)).Msg("bulk updating attribute cache")

	successCount := 0
	errorCount := 0

	for _, listingID := range listingIDs {
		if err := idx.UpdateListingAttributeCache(ctx, listingID); err != nil {
			idx.logger.Error().Err(err).Int32("listing_id", listingID).Msg("failed to update cache")
			errorCount++
			continue
		}
		successCount++
	}

	idx.logger.Info().
		Int("success", successCount).
		Int("errors", errorCount).
		Msg("bulk cache update completed")

	if errorCount > 0 {
		return fmt.Errorf("completed with %d errors out of %d listings", errorCount, len(listingIDs))
	}

	return nil
}
