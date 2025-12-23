package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	opensearch "github.com/opensearch-project/opensearch-go/v2"
	"github.com/rs/zerolog"

	"github.com/vondi-global/listings/internal/domain"
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

// NewClientSimple creates a simple OpenSearch client with just a URL (for CLI tools)
func NewClientSimple(url string, logger zerolog.Logger) (*Client, error) {
	return NewClient([]string{url}, "", "", "marketplace_listings", logger)
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
		"source_type":     listing.SourceType,  // c2c or b2c
		"document_type":   "listing",           // Required for unified search filtering
		"stock_status":    listing.StockStatus, // in_stock, out_of_stock, etc
		"views_count":     listing.ViewsCount,
		"favorites_count": listing.FavoritesCount,
		"created_at":      listing.CreatedAt,
		"updated_at":      listing.UpdatedAt,
		"published_at":    listing.PublishedAt,
	}

	// Add location if available (for geo_distance similarity scoring)
	if listing.Location != nil {
		loc := listing.Location
		if loc.Latitude != nil && loc.Longitude != nil && (*loc.Latitude != 0 || *loc.Longitude != 0) {
			doc["location"] = map[string]interface{}{
				"lat": *loc.Latitude,
				"lon": *loc.Longitude,
			}
			doc["has_location"] = true
		}
		if loc.City != nil {
			doc["city"] = *loc.City
		}
		if loc.Country != nil {
			doc["country"] = *loc.Country
		}
		if loc.AddressLine1 != nil {
			doc["address"] = *loc.AddressLine1
		}
	}

	// Add attributes from cache if available
	attributes, searchableText, err := c.getAttributesFromCache(ctx, int32(listing.ID))
	if err != nil {
		// Log but don't fail - attributes are optional
		c.logger.Debug().Err(err).Int64("listing_id", listing.ID).Msg("no attributes cache found")
	} else {
		doc["attributes"] = attributes
		doc["attributes_searchable_text"] = searchableText
	}

	// Add listing attributes if present (for similarity matching)
	if len(listing.Attributes) > 0 {
		attrMap := make(map[string]string)
		for _, attr := range listing.Attributes {
			attrMap[attr.AttributeKey] = attr.AttributeValue
		}
		doc["listing_attributes"] = attrMap
	}

	// Add images if present
	if len(listing.Images) > 0 {
		images := make([]map[string]interface{}, 0, len(listing.Images))
		for _, img := range listing.Images {
			images = append(images, map[string]interface{}{
				"public_url": img.URL,
				"is_main":    img.IsPrimary,
			})
		}
		doc["images"] = images
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
		body, _ := io.ReadAll(res.Body)
		c.logger.Error().Int("status", res.StatusCode).Int64("listing_id", listing.ID).Str("response_body", string(body)).Msg("OpenSearch index error")
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

		// CategoryID can be string (UUID) or float64 (legacy int)
		var categoryID string
		if catStr, ok := source["category_id"].(string); ok {
			categoryID = catStr
		} else if catNum, ok := source["category_id"].(float64); ok {
			// Legacy: convert int to string
			categoryID = fmt.Sprintf("%d", int64(catNum))
		}

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
			CategoryID:     categoryID,
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

// GetSimilarListings finds listings similar to the given listing using multi-factor similarity.
//
// Similarity factors (in order of importance):
// 1. Title similarity (more_like_this on title field)
// 2. Description similarity (more_like_this on description field)
// 3. Category match (same category = higher score)
// 4. Price range (±30% for broader results)
// 5. Geographic proximity (if location available)
// 6. Attribute matching (if attributes available)
//
// Results are ranked by combined similarity score.
func (c *Client) GetSimilarListings(ctx context.Context, listingID int64, limit int32) ([]*domain.Listing, int32, error) {
	// First, fetch the source listing to get its properties
	res, err := c.client.Get(
		c.index,
		fmt.Sprintf("%d", listingID),
		c.client.Get.WithContext(ctx),
	)
	if err != nil {
		c.logger.Error().Err(err).Int64("listing_id", listingID).Msg("failed to fetch source listing")
		return nil, 0, fmt.Errorf("failed to fetch source listing: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		if res.StatusCode == 404 {
			c.logger.Warn().Int64("listing_id", listingID).Msg("source listing not found in index")
			return []*domain.Listing{}, 0, nil
		}
		return nil, 0, fmt.Errorf("failed to get listing: %s", res.Status())
	}

	var getResult map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&getResult); err != nil {
		return nil, 0, fmt.Errorf("failed to parse get response: %w", err)
	}

	source, ok := getResult["_source"].(map[string]interface{})
	if !ok {
		return nil, 0, fmt.Errorf("invalid source structure in get response")
	}

	// Extract source listing properties
	// category_id is stored as string (UUID)
	categoryID, _ := source["category_id"].(string)
	title, _ := source["title"].(string)
	description, _ := source["description"].(string)

	// Extract all title translations for cross-language matching
	titleSr, _ := source["title_sr"].(string)
	titleEn, _ := source["title_en"].(string)
	titleRu, _ := source["title_ru"].(string)

	// Combine all titles for cross-language search
	// This allows "Gamepad" in Serbian title to match "gamepad" in English translation
	allTitles := title
	if titleSr != "" {
		allTitles += " " + titleSr
	}
	if titleEn != "" {
		allTitles += " " + titleEn
	}
	if titleRu != "" {
		allTitles += " " + titleRu
	}

	// Price can be float64 or int (OpenSearch returns long as float64)
	var price float64
	if p, ok := source["price"].(float64); ok {
		price = p
	}

	// Location for geo_distance scoring (if available)
	var lat, lon float64
	var hasLocation bool
	if loc, ok := source["location"].(map[string]interface{}); ok {
		if latVal, ok := loc["lat"].(float64); ok {
			lat = latVal
		}
		if lonVal, ok := loc["lon"].(float64); ok {
			lon = lonVal
		}
		hasLocation = lat != 0 || lon != 0
	}

	// Build multi-factor similarity query using function_score
	// Base filters: must be active, public, not the same listing
	// Note: status, visibility, category_id are already keyword type in mapping, no .keyword suffix needed
	mustFilters := []interface{}{
		map[string]interface{}{"term": map[string]interface{}{"status": "active"}},
		map[string]interface{}{"term": map[string]interface{}{"visibility": "public"}},
	}

	mustNotFilters := []interface{}{
		map[string]interface{}{"term": map[string]interface{}{"id": listingID}},
	}

	// Build should clauses for similarity scoring
	shouldClauses := []interface{}{}

	// 1. Title similarity using multi_match with cross_fields
	// This enables cross-language matching where "Gamepad" in Serbian title
	// matches "gamepad" in English translation of another listing
	if allTitles != "" {
		shouldClauses = append(shouldClauses, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query": allTitles,
				"fields": []string{
					"title^3", "title_sr^2", "title_en^2", "title_ru^2",
				},
				"type":                 "best_fields",
				"tie_breaker":          0.3,
				"minimum_should_match": "1", // At least 1 term must match
				"boost":                10.0, // High boost for title matching
			},
		})
	}

	// 2. Description similarity (LOWER priority)
	// Description often contains common words that may create false positives
	if description != "" {
		shouldClauses = append(shouldClauses, map[string]interface{}{
			"more_like_this": map[string]interface{}{
				"fields": []string{
					"description", "description_sr", "description_en", "description_ru",
				},
				"like": []interface{}{
					map[string]interface{}{
						"_index": c.index,
						"_id":    fmt.Sprintf("%d", listingID),
					},
				},
				"min_term_freq":        1,
				"min_doc_freq":         1,
				"max_query_terms":      50,
				"minimum_should_match": "1",
				"boost":                2.0, // Lower boost for description
			},
		})
	}

	// 3. Same category (high weight)
	if categoryID != "" {
		shouldClauses = append(shouldClauses, map[string]interface{}{
			"term": map[string]interface{}{
				"category_id": map[string]interface{}{
					"value": categoryID,
					"boost": 4.0,
				},
			},
		})
	}

	// 4. Price range filter (±30% for broader results, but not a must)
	if price > 0 {
		minPrice := price * 0.7
		maxPrice := price * 1.3
		shouldClauses = append(shouldClauses, map[string]interface{}{
			"range": map[string]interface{}{
				"price": map[string]interface{}{
					"gte":   minPrice,
					"lte":   maxPrice,
					"boost": 2.0,
				},
			},
		})
	}

	// Build the main query
	boolQuery := map[string]interface{}{
		"must":     mustFilters,
		"must_not": mustNotFilters,
		"should":   shouldClauses,
	}

	// Require at least one should clause to match
	if len(shouldClauses) > 0 {
		boolQuery["minimum_should_match"] = 1
	}

	// Build function_score query for additional scoring factors
	functions := []interface{}{}

	// 5. Geographic proximity scoring (if location available)
	if hasLocation {
		functions = append(functions, map[string]interface{}{
			"gauss": map[string]interface{}{
				"location": map[string]interface{}{
					"origin": map[string]interface{}{
						"lat": lat,
						"lon": lon,
					},
					"scale":  "10km", // Distance at which score drops to ~0.5
					"offset": "1km",  // No decay within this distance
					"decay":  0.5,
				},
			},
			"weight": 2.0, // Location proximity weight
		})
	}

	// 6. Price proximity scoring (closer price = higher score)
	if price > 0 {
		functions = append(functions, map[string]interface{}{
			"gauss": map[string]interface{}{
				"price": map[string]interface{}{
					"origin": price,
					"scale":  price * 0.2, // 20% price difference = ~0.5 score
					"offset": price * 0.05,
					"decay":  0.5,
				},
			},
			"weight": 1.5,
		})
	}

	// 7. Recency boost (newer listings slightly preferred)
	functions = append(functions, map[string]interface{}{
		"gauss": map[string]interface{}{
			"created_at": map[string]interface{}{
				"origin": "now",
				"scale":  "30d",
				"offset": "7d",
				"decay":  0.5,
			},
		},
		"weight": 0.5,
	})

	// Build final query
	var searchQuery map[string]interface{}
	if len(functions) > 0 {
		searchQuery = map[string]interface{}{
			"query": map[string]interface{}{
				"function_score": map[string]interface{}{
					"query":      map[string]interface{}{"bool": boolQuery},
					"functions":  functions,
					"score_mode": "sum",      // Sum all function scores
					"boost_mode": "multiply", // Multiply with query score
				},
			},
			"size": limit,
		}
	} else {
		searchQuery = map[string]interface{}{
			"query": map[string]interface{}{"bool": boolQuery},
			"size":  limit,
			"sort": []interface{}{
				map[string]interface{}{"_score": "desc"},
				map[string]interface{}{"created_at": "desc"},
			},
		}
	}

	body, err := json.Marshal(searchQuery)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to marshal similarity query: %w", err)
	}

	c.logger.Debug().
		Int64("listing_id", listingID).
		Str("category_id", categoryID).
		Str("title", title).
		Float64("price", price).
		Bool("has_location", hasLocation).
		RawJSON("query", body).
		Msg("executing similarity search")

	searchRes, err := c.client.Search(
		c.client.Search.WithContext(ctx),
		c.client.Search.WithIndex(c.index),
		c.client.Search.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		c.logger.Error().Err(err).Msg("similarity search failed")
		return nil, 0, fmt.Errorf("similarity search failed: %w", err)
	}
	defer searchRes.Body.Close()

	if searchRes.IsError() {
		bodyBytes, _ := io.ReadAll(searchRes.Body)
		c.logger.Error().
			Int("status", searchRes.StatusCode).
			Str("response", string(bodyBytes)).
			Msg("similarity search error")
		return nil, 0, fmt.Errorf("similarity search error: %s", searchRes.Status())
	}

	// Parse search results
	var result map[string]interface{}
	if err := json.NewDecoder(searchRes.Body).Decode(&result); err != nil {
		return nil, 0, fmt.Errorf("failed to parse similarity search response: %w", err)
	}

	hits, ok := result["hits"].(map[string]interface{})
	if !ok {
		return nil, 0, fmt.Errorf("invalid hits structure in similarity response")
	}

	total, ok := hits["total"].(map[string]interface{})
	if !ok {
		return nil, 0, fmt.Errorf("invalid total structure in similarity response")
	}

	totalValue, ok := total["value"].(float64)
	if !ok {
		return nil, 0, fmt.Errorf("invalid total value in similarity response")
	}
	totalHits := int64(totalValue)

	hitsArray, ok := hits["hits"].([]interface{})
	if !ok {
		return nil, 0, fmt.Errorf("invalid hits array in similarity response")
	}

	listings := make([]*domain.Listing, 0, len(hitsArray))
	for _, hit := range hitsArray {
		hitMap, ok := hit.(map[string]interface{})
		if !ok {
			continue
		}

		hitSource, ok := hitMap["_source"].(map[string]interface{})
		if !ok {
			continue
		}

		listing := c.parseListingFromSource(hitSource)
		if listing != nil {
			listings = append(listings, listing)
		}
	}

	c.logger.Info().
		Int64("listing_id", listingID).
		Str("category_id", categoryID).
		Float64("price", price).
		Int("results", len(listings)).
		Int64("total", totalHits).
		Msg("similarity search completed")

	return listings, int32(totalHits), nil
}

// parseListingFromSource extracts a Listing from OpenSearch _source document
func (c *Client) parseListingFromSource(source map[string]interface{}) *domain.Listing {
	// Extract ID (required field)
	id, ok := source["id"].(float64)
	if !ok {
		return nil
	}

	listing := &domain.Listing{
		ID: int64(id),
	}

	// String fields
	if v, ok := source["uuid"].(string); ok {
		listing.UUID = v
	}
	if v, ok := source["title"].(string); ok {
		listing.Title = v
	}
	if v, ok := source["currency"].(string); ok {
		listing.Currency = v
	}
	if v, ok := source["status"].(string); ok {
		listing.Status = v
	}
	if v, ok := source["visibility"].(string); ok {
		listing.Visibility = v
	}
	if v, ok := source["source_type"].(string); ok {
		listing.SourceType = v
	}

	// CategoryID can be string (UUID) or float64 (legacy int)
	if catStr, ok := source["category_id"].(string); ok {
		listing.CategoryID = catStr
	} else if catNum, ok := source["category_id"].(float64); ok {
		listing.CategoryID = fmt.Sprintf("%d", int64(catNum))
	}

	// Numeric fields
	if v, ok := source["user_id"].(float64); ok {
		listing.UserID = int64(v)
	}
	if v, ok := source["price"].(float64); ok {
		listing.Price = v
	}
	if v, ok := source["quantity"].(float64); ok {
		listing.Quantity = int32(v)
	}
	if v, ok := source["views_count"].(float64); ok {
		listing.ViewsCount = int32(v)
	}
	if v, ok := source["favorites_count"].(float64); ok {
		listing.FavoritesCount = int32(v)
	}

	// Optional string fields (pointers)
	if v, ok := source["description"].(string); ok {
		listing.Description = &v
	}
	if v, ok := source["sku"].(string); ok {
		listing.SKU = &v
	}
	if v, ok := source["stock_status"].(string); ok {
		listing.StockStatus = &v
	}
	if v, ok := source["attributes"].(string); ok {
		listing.AttributesJSON = &v
	}

	// Optional numeric fields (pointers)
	if v, ok := source["storefront_id"].(float64); ok {
		sfID := int64(v)
		listing.StorefrontID = &sfID
	}

	return listing
}

// buildFilters constructs filter clauses for search query
// Fields like status, visibility, source_type, category_id are already keyword type in index mapping
func buildFilters(query *domain.SearchListingsQuery) []interface{} {
	// Note: status, visibility, source_type, category_id are already keyword type in mapping
	filters := []interface{}{
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

	// Add attribute filters
	if len(query.AttributeFilters) > 0 {
		for _, attrFilter := range query.AttributeFilters {
			// Build nested query for each attribute filter
			if attrFilter.MinNumber != nil || attrFilter.MaxNumber != nil {
				// Range query for numeric attributes
				filters = append(filters, GetAttributeRangeQuery(
					attrFilter.Code,
					attrFilter.MinNumber,
					attrFilter.MaxNumber,
				))
			} else {
				// Exact match query
				filters = append(filters, GetAttributeNestedQuery(
					attrFilter.Code,
					attrFilter.ValueText,
					attrFilter.ValueNumber,
					attrFilter.ValueBool,
				))
			}
		}
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

	// CategoryID can be string (UUID) or float64 (legacy int)
	var categoryID string
	if catStr, ok := source["category_id"].(string); ok {
		categoryID = catStr
	} else if catNum, ok := source["category_id"].(float64); ok {
		// Legacy: convert int to string
		categoryID = fmt.Sprintf("%d", int64(catNum))
	}

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
		CategoryID: categoryID,
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

// IndexProduct indexes a single B2C product in OpenSearch
// Used for real-time indexing when product is created
func (c *Client) IndexProduct(ctx context.Context, product *domain.Listing) error {
	if product == nil {
		return fmt.Errorf("product cannot be nil")
	}

	// Prepare document for indexing (similar format to C2C listings for unified index)
	doc := c.buildProductDocument(product)

	body, err := json.Marshal(doc)
	if err != nil {
		c.logger.Error().Err(err).Int64("product_id", product.ID).Msg("failed to marshal product document")
		return fmt.Errorf("failed to marshal product document: %w", err)
	}

	// Index document
	res, err := c.client.Index(
		c.index,
		bytes.NewReader(body),
		c.client.Index.WithContext(ctx),
		c.client.Index.WithDocumentID(fmt.Sprintf("%d", product.ID)),
		c.client.Index.WithRefresh("false"), // Async refresh for performance
	)

	if err != nil {
		c.logger.Error().Err(err).Int64("product_id", product.ID).Msg("failed to index product")
		return fmt.Errorf("failed to index product: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		c.logger.Error().Int("status", res.StatusCode).Int64("product_id", product.ID).Msg("OpenSearch index error")
		return fmt.Errorf("OpenSearch index error: %s", res.Status())
	}

	c.logger.Debug().Int64("product_id", product.ID).Msg("product indexed successfully")
	return nil
}

// UpdateProduct updates a B2C product document in OpenSearch
// Used for real-time updates when product is modified
func (c *Client) UpdateProduct(ctx context.Context, product *domain.Listing) error {
	// For simplicity, re-index the entire document
	// In production, you might want partial updates for efficiency
	return c.IndexProduct(ctx, product)
}

// DeleteProduct removes a B2C product from OpenSearch index
// Used for real-time deletion when product is deleted
func (c *Client) DeleteProduct(ctx context.Context, productID int64) error {
	res, err := c.client.Delete(
		c.index,
		fmt.Sprintf("%d", productID),
		c.client.Delete.WithContext(ctx),
	)

	if err != nil {
		c.logger.Error().Err(err).Int64("product_id", productID).Msg("failed to delete product from index")
		return fmt.Errorf("failed to delete product: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		c.logger.Error().Int("status", res.StatusCode).Int64("product_id", productID).Msg("OpenSearch delete error")
		return fmt.Errorf("OpenSearch delete error: %s", res.Status())
	}

	c.logger.Debug().Int64("product_id", productID).Msg("product deleted from index")
	return nil
}

// BulkIndexProducts indexes multiple products in a single bulk request
// Used for batch operations and reindexing
// Recommended batch size: 100-1000 products
func (c *Client) BulkIndexProducts(ctx context.Context, products []*domain.Listing) error {
	if len(products) == 0 {
		return nil
	}

	c.logger.Debug().Int("count", len(products)).Msg("bulk indexing products")

	// Build bulk request body
	var bulkBody bytes.Buffer
	for _, product := range products {
		// Action line (index operation)
		action := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": c.index,
				"_id":    fmt.Sprintf("%d", product.ID),
			},
		}

		actionJSON, err := json.Marshal(action)
		if err != nil {
			c.logger.Error().Err(err).Int64("product_id", product.ID).Msg("failed to marshal action")
			return fmt.Errorf("failed to marshal action for product %d: %w", product.ID, err)
		}
		bulkBody.Write(actionJSON)
		bulkBody.WriteByte('\n')

		// Document line
		doc := c.buildProductDocument(product)
		docJSON, err := json.Marshal(doc)
		if err != nil {
			c.logger.Error().Err(err).Int64("product_id", product.ID).Msg("failed to marshal document")
			return fmt.Errorf("failed to marshal document for product %d: %w", product.ID, err)
		}
		bulkBody.Write(docJSON)
		bulkBody.WriteByte('\n')
	}

	// Execute bulk request
	res, err := c.client.Bulk(
		bytes.NewReader(bulkBody.Bytes()),
		c.client.Bulk.WithContext(ctx),
	)

	if err != nil {
		c.logger.Error().Err(err).Int("product_count", len(products)).Msg("bulk index request failed")
		return fmt.Errorf("bulk index request failed: %w", err)
	}
	defer res.Body.Close()

	// Check response
	if res.IsError() {
		body, _ := io.ReadAll(res.Body)
		c.logger.Error().Int("status", res.StatusCode).Str("body", string(body)).Msg("bulk index response error")
		return fmt.Errorf("bulk index response error: %s", res.Status())
	}

	// Parse response to check for individual errors
	var bulkResp map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&bulkResp); err != nil {
		c.logger.Error().Err(err).Msg("failed to parse bulk response")
		return fmt.Errorf("failed to parse bulk response: %w", err)
	}

	// Check for errors in bulk response
	if errors, ok := bulkResp["errors"].(bool); ok && errors {
		c.logger.Warn().Msg("some items in bulk request failed")
		// Log individual errors but don't fail the entire operation
		if items, ok := bulkResp["items"].([]interface{}); ok {
			for i, item := range items {
				if itemMap, ok := item.(map[string]interface{}); ok {
					for action, details := range itemMap {
						if detailsMap, ok := details.(map[string]interface{}); ok {
							if errorInfo, hasError := detailsMap["error"]; hasError {
								c.logger.Warn().
									Int("index", i).
									Str("action", action).
									Interface("error", errorInfo).
									Msg("bulk item error")
							}
						}
					}
				}
			}
		}
		return fmt.Errorf("some items in bulk request failed (see logs)")
	}

	c.logger.Info().Int("count", len(products)).Msg("bulk index completed successfully")
	return nil
}

// buildProductDocument builds an OpenSearch document from a product/listing
// Formats document to match the unified index structure used by monolith
func (c *Client) buildProductDocument(product *domain.Listing) map[string]interface{} {
	doc := map[string]interface{}{
		"id":              product.ID,
		"uuid":            product.UUID,
		"source_type":     product.SourceType, // "b2c" or "c2c"
		"document_type":   "listing",          // Required for unified search
		"title":           product.Title,
		"price":           product.Price,
		"currency":        product.Currency,
		"category_id":     product.CategoryID,
		"status":          product.Status,
		"visibility":      product.Visibility,
		"views_count":     product.ViewsCount,
		"favorites_count": product.FavoritesCount,
		"created_at":      product.CreatedAt,
		"updated_at":      product.UpdatedAt,
	}

	// Optional fields
	if product.Description != nil {
		doc["description"] = *product.Description
	}
	if product.StorefrontID != nil {
		doc["storefront_id"] = *product.StorefrontID
	}
	if product.SKU != nil {
		doc["sku"] = *product.SKU
	}
	if product.StockStatus != nil {
		doc["stock_status"] = *product.StockStatus
	}
	if product.AttributesJSON != nil {
		doc["attributes"] = *product.AttributesJSON
	}
	if product.PublishedAt != nil {
		doc["published_at"] = *product.PublishedAt
	}

	// Add images if present
	if len(product.Images) > 0 {
		images := make([]map[string]interface{}, 0, len(product.Images))
		for _, img := range product.Images {
			images = append(images, map[string]interface{}{
				"id":        img.ID,
				"file_path": img.URL,
				"is_main":   img.IsPrimary,
			})
		}
		doc["images"] = images
	}

	// Add location if present (B2C products may have individual locations)
	if product.Location != nil {
		loc := product.Location
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

	return doc
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

// getAttributesFromCache retrieves cached attributes for a listing from attribute_search_cache
// Returns: attributes array, searchable text, error
func (c *Client) getAttributesFromCache(ctx context.Context, listingID int32) ([]interface{}, string, error) {
	// This method requires database connection which Client doesn't have
	// It should be called from a higher level service that has both OpenSearch and DB access
	// For now, return empty to avoid breaking changes
	// TODO: Refactor to pass attributes from service layer
	return nil, "", fmt.Errorf("attributes cache not implemented in client")
}

// DeleteIndex deletes an OpenSearch index
func (c *Client) DeleteIndex(ctx context.Context, indexName string) error {
	res, err := c.client.Indices.Delete(
		[]string{indexName},
		c.client.Indices.Delete.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("failed to delete index %s: %w", indexName, err)
	}
	defer res.Body.Close()

	if res.IsError() {
		if res.StatusCode == 404 {
			// Index doesn't exist - not an error
			c.logger.Debug().Str("index", indexName).Msg("index does not exist")
			return nil
		}
		return fmt.Errorf("failed to delete index %s: %s", indexName, res.Status())
	}

	c.logger.Info().Str("index", indexName).Msg("index deleted successfully")
	return nil
}

// CreateIndex creates an OpenSearch index with the given mapping
func (c *Client) CreateIndex(ctx context.Context, indexName string, mapping map[string]interface{}) error {
	body, err := json.Marshal(mapping)
	if err != nil {
		return fmt.Errorf("failed to marshal mapping: %w", err)
	}

	res, err := c.client.Indices.Create(
		indexName,
		c.client.Indices.Create.WithContext(ctx),
		c.client.Indices.Create.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		return fmt.Errorf("failed to create index %s: %w", indexName, err)
	}
	defer res.Body.Close()

	if res.IsError() {
		bodyBytes, _ := io.ReadAll(res.Body)
		return fmt.Errorf("failed to create index %s: %s - %s", indexName, res.Status(), string(bodyBytes))
	}

	c.logger.Info().Str("index", indexName).Msg("index created successfully")
	return nil
}

// CountDocuments returns the number of documents in an index
func (c *Client) CountDocuments(ctx context.Context, indexName string) (int, error) {
	res, err := c.client.Count(
		c.client.Count.WithContext(ctx),
		c.client.Count.WithIndex(indexName),
	)
	if err != nil {
		return 0, fmt.Errorf("failed to count documents in %s: %w", indexName, err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return 0, fmt.Errorf("failed to count documents in %s: %s", indexName, res.Status())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to parse count response: %w", err)
	}

	count, ok := result["count"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid count response format")
	}

	return int(count), nil
}
