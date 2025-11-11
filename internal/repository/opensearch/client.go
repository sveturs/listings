package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	opensearch "github.com/opensearch-project/opensearch-go/v2"
	"github.com/rs/zerolog"

	"github.com/sveturs/listings/internal/domain"
)

// Client handles OpenSearch operations for listings
type Client struct {
	client *opensearch.Client
	index  string
	logger zerolog.Logger
}

// NewClient creates a new OpenSearch client
func NewClient(addresses []string, username, password, index string, logger zerolog.Logger) (*Client, error) {
	cfg := opensearch.Config{
		Addresses: addresses,
		Username:  username,
		Password:  password,
	}

	client, err := opensearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenSearch client: %w", err)
	}

	// Test connection
	res, err := client.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to OpenSearch: %w", err)
	}
	defer res.Body.Close()

	logger.Info().
		Strs("addresses", addresses).
		Str("index", index).
		Msg("OpenSearch client initialized")

	return &Client{
		client: client,
		index:  index,
		logger: logger.With().Str("component", "opensearch_client").Logger(),
	}, nil
}

// IndexListing indexes a listing document in OpenSearch
func (c *Client) IndexListing(ctx context.Context, listing *domain.Listing) error {
	// Prepare document for indexing
	doc := map[string]interface{}{
		"id":              listing.ID,
		"uuid":            listing.UUID,
		"user_id":         listing.UserID,
		"storefront_id":   listing.StorefrontID,
		"title":           listing.Title,
		"description":     listing.Description,
		"price":           listing.Price,
		"currency":        listing.Currency,
		"category_id":     listing.CategoryID,
		"status":          listing.Status,
		"visibility":      listing.Visibility,
		"quantity":        listing.Quantity,
		"sku":             listing.SKU,
		"source_type":     listing.SourceType,     // c2c or b2c
		"document_type":   "listing",              // Required for unified search filtering
		"stock_status":    listing.StockStatus,    // in_stock, out_of_stock, etc
		"attributes":      listing.AttributesJSON, // JSONB attributes as string
		"views_count":     listing.ViewsCount,
		"favorites_count": listing.FavoritesCount,
		"created_at":      listing.CreatedAt,
		"updated_at":      listing.UpdatedAt,
		"published_at":    listing.PublishedAt,
	}

	body, err := json.Marshal(doc)
	if err != nil {
		c.logger.Error().Err(err).Int64("listing_id", listing.ID).Msg("failed to marshal document")
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	// Index document
	res, err := c.client.Index(
		c.index,
		bytes.NewReader(body),
		c.client.Index.WithContext(ctx),
		c.client.Index.WithDocumentID(fmt.Sprintf("%d", listing.ID)),
		c.client.Index.WithRefresh("false"), // Async refresh
	)

	if err != nil {
		c.logger.Error().Err(err).Int64("listing_id", listing.ID).Msg("failed to index listing")
		return fmt.Errorf("failed to index listing: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		c.logger.Error().Int("status", res.StatusCode).Int64("listing_id", listing.ID).Msg("OpenSearch index error")
		return fmt.Errorf("OpenSearch index error: %s", res.Status())
	}

	c.logger.Debug().Int64("listing_id", listing.ID).Msg("listing indexed successfully")
	return nil
}

// UpdateListing updates a listing document in OpenSearch
func (c *Client) UpdateListing(ctx context.Context, listing *domain.Listing) error {
	// For simplicity, re-index the entire document
	return c.IndexListing(ctx, listing)
}

// DeleteListing removes a listing from OpenSearch index
func (c *Client) DeleteListing(ctx context.Context, listingID int64) error {
	res, err := c.client.Delete(
		c.index,
		fmt.Sprintf("%d", listingID),
		c.client.Delete.WithContext(ctx),
	)

	if err != nil {
		c.logger.Error().Err(err).Int64("listing_id", listingID).Msg("failed to delete listing from index")
		return fmt.Errorf("failed to delete listing: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		c.logger.Error().Int("status", res.StatusCode).Int64("listing_id", listingID).Msg("OpenSearch delete error")
		return fmt.Errorf("OpenSearch delete error: %s", res.Status())
	}

	c.logger.Debug().Int64("listing_id", listingID).Msg("listing deleted from index")
	return nil
}

// HealthCheck performs a health check on OpenSearch
func (c *Client) HealthCheck(ctx context.Context) error {
	res, err := c.client.Cluster.Health()
	if err != nil {
		return fmt.Errorf("OpenSearch health check failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("OpenSearch health check error: %s", res.Status())
	}

	return nil
}

// SearchListings performs a full-text search on listings
func (c *Client) SearchListings(ctx context.Context, query *domain.SearchListingsQuery) ([]*domain.Listing, int64, error) {
	// Build OpenSearch query
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []interface{}{
					map[string]interface{}{
						"multi_match": map[string]interface{}{
							"query":  query.Query,
							"fields": []string{"title^3", "description"},
							"type":   "best_fields",
						},
					},
				},
				"filter": buildFilters(query),
			},
		},
		"from": query.Offset,
		"size": query.Limit,
		"sort": []interface{}{
			map[string]interface{}{"_score": "desc"},
			map[string]interface{}{"created_at": "desc"},
		},
	}

	body, err := json.Marshal(searchQuery)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to marshal search query: %w", err)
	}

	res, err := c.client.Search(
		c.client.Search.WithContext(ctx),
		c.client.Search.WithIndex(c.index),
		c.client.Search.WithBody(bytes.NewReader(body)),
	)

	if err != nil {
		c.logger.Error().Err(err).Msg("search query failed")
		return nil, 0, fmt.Errorf("search query failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, 0, fmt.Errorf("search query error: %s", res.Status())
	}

	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, 0, fmt.Errorf("failed to parse search response: %w", err)
	}

	hits, ok := result["hits"].(map[string]interface{})
	if !ok {
		return nil, 0, fmt.Errorf("invalid hits structure in search response")
	}

	total, ok := hits["total"].(map[string]interface{})
	if !ok {
		return nil, 0, fmt.Errorf("invalid total structure in search response")
	}

	totalValue, ok := total["value"].(float64)
	if !ok {
		return nil, 0, fmt.Errorf("invalid total value in search response")
	}
	totalHits := int64(totalValue)

	hitsArray, ok := hits["hits"].([]interface{})
	if !ok {
		return nil, 0, fmt.Errorf("invalid hits array in search response")
	}

	var listings []*domain.Listing
	for _, hit := range hitsArray {
		hitMap, ok := hit.(map[string]interface{})
		if !ok {
			continue
		}

		source, ok := hitMap["_source"].(map[string]interface{})
		if !ok {
			continue
		}

		// Extract fields with safe type assertions (already validated structure above)
		//nolint:errcheck // Type assertions are safe after validating structure
		id, _ := source["id"].(float64)
		//nolint:errcheck
		uuid, _ := source["uuid"].(string)
		//nolint:errcheck
		userID, _ := source["user_id"].(float64)
		//nolint:errcheck
		title, _ := source["title"].(string)
		//nolint:errcheck
		price, _ := source["price"].(float64)
		//nolint:errcheck
		currency, _ := source["currency"].(string)
		//nolint:errcheck
		categoryID, _ := source["category_id"].(float64)
		//nolint:errcheck
		status, _ := source["status"].(string)
		//nolint:errcheck
		visibility, _ := source["visibility"].(string)
		//nolint:errcheck
		quantity, _ := source["quantity"].(float64)
		//nolint:errcheck
		sourceType, _ := source["source_type"].(string)
		//nolint:errcheck
		viewsCount, _ := source["views_count"].(float64)
		//nolint:errcheck
		favoritesCount, _ := source["favorites_count"].(float64)

		listing := &domain.Listing{
			ID:             int64(id),
			UUID:           uuid,
			UserID:         int64(userID),
			Title:          title,
			Price:          price,
			Currency:       currency,
			CategoryID:     int64(categoryID),
			Status:         status,
			Visibility:     visibility,
			Quantity:       int32(quantity),
			SourceType:     sourceType,
			ViewsCount:     int32(viewsCount),
			FavoritesCount: int32(favoritesCount),
		}

		// Optional fields
		if desc, ok := source["description"].(string); ok {
			listing.Description = &desc
		}
		if sfID, ok := source["storefront_id"].(float64); ok {
			id := int64(sfID)
			listing.StorefrontID = &id
		}
		if sku, ok := source["sku"].(string); ok {
			listing.SKU = &sku
		}
		if stockStatus, ok := source["stock_status"].(string); ok {
			listing.StockStatus = &stockStatus
		}
		if attributes, ok := source["attributes"].(string); ok {
			listing.AttributesJSON = &attributes
		}

		listings = append(listings, listing)
	}

	c.logger.Debug().
		Str("query", query.Query).
		Int("results", len(listings)).
		Int64("total", totalHits).
		Msg("search completed")

	return listings, totalHits, nil
}

// buildFilters constructs filter clauses for search query
func buildFilters(query *domain.SearchListingsQuery) []interface{} {
	filters := []interface{}{
		// Only active, visible listings
		map[string]interface{}{"term": map[string]interface{}{"status": "active"}},
		map[string]interface{}{"term": map[string]interface{}{"visibility": "public"}},
	}

	// Filter by source_type (c2c vs b2c) if specified
	if query.SourceType != nil && *query.SourceType != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{"source_type": *query.SourceType},
		})
	}

	if query.CategoryID != nil {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{"category_id": *query.CategoryID},
		})
	}

	if query.MinPrice != nil || query.MaxPrice != nil {
		priceRange := map[string]interface{}{}
		if query.MinPrice != nil {
			priceRange["gte"] = *query.MinPrice
		}
		if query.MaxPrice != nil {
			priceRange["lte"] = *query.MaxPrice
		}
		filters = append(filters, map[string]interface{}{
			"range": map[string]interface{}{"price": priceRange},
		})
	}

	return filters
}

// GetListingByID retrieves a single listing by ID from OpenSearch
func (c *Client) GetListingByID(ctx context.Context, listingID int64) (*domain.Listing, error) {
	res, err := c.client.Get(
		c.index,
		fmt.Sprintf("%d", listingID),
		c.client.Get.WithContext(ctx),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get listing: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return nil, fmt.Errorf("listing not found")
	}

	if res.IsError() {
		return nil, fmt.Errorf("OpenSearch get error: %s", res.Status())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	source, ok := result["_source"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid source structure in response")
	}

	// Extract fields with safe type assertions (already validated structure above)
	//nolint:errcheck // Type assertions are safe after validating structure
	id, _ := source["id"].(float64)
	//nolint:errcheck
	uuid, _ := source["uuid"].(string)
	//nolint:errcheck
	userID, _ := source["user_id"].(float64)
	//nolint:errcheck
	title, _ := source["title"].(string)
	//nolint:errcheck
	price, _ := source["price"].(float64)
	//nolint:errcheck
	currency, _ := source["currency"].(string)
	//nolint:errcheck
	categoryID, _ := source["category_id"].(float64)
	//nolint:errcheck
	status, _ := source["status"].(string)
	//nolint:errcheck
	sourceType, _ := source["source_type"].(string)

	listing := &domain.Listing{
		ID:         int64(id),
		UUID:       uuid,
		UserID:     int64(userID),
		Title:      title,
		Price:      price,
		Currency:   currency,
		CategoryID: int64(categoryID),
		Status:     status,
		SourceType: sourceType,
	}

	// Optional fields
	if stockStatus, ok := source["stock_status"].(string); ok {
		listing.StockStatus = &stockStatus
	}
	if attributes, ok := source["attributes"].(string); ok {
		listing.AttributesJSON = &attributes
	}

	return listing, nil
}

// GetClient returns the underlying OpenSearch client for advanced usage
func (c *Client) GetClient() *opensearch.Client {
	return c.client
}

// Close closes the OpenSearch client (no-op for opensearch-go client)
func (c *Client) Close() error {
	// opensearch-go client doesn't require explicit closing
	return nil
}
