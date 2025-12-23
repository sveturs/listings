package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

// ============================================================================
// PHASE 4: Advanced Search Queries - Function Score, Did You Mean, Autocomplete
// ============================================================================

// UniversalSearchParams represents parameters for universal search with ranking
type UniversalSearchParams struct {
	Query          string
	CategoryIDs    []string            // UUID strings
	PriceMin       *float64
	PriceMax       *float64
	SourceTypes    []string            // "c2c", "b2c"
	City           string
	Country        string
	Attributes     map[string][]string // attribute_code -> values
	SortBy         string              // "price", "created_at", "relevance", "popularity"
	SortOrder      string              // "asc", "desc"
	Page           int
	Limit          int
	IncludeVariants bool
	EnableFuzzy    bool
	Highlighting   bool // Enable highlighting
}

// UniversalSearchResult represents search result with ranking signals
type UniversalSearchResult struct {
	Listings    []*ListingResult
	Total       int64
	Page        int
	Limit       int
	Suggestions []string
	TookMs      int
}

// ListingResult represents a single listing in search results
type ListingResult struct {
	ID          int64
	UUID        string
	Title       string
	Price       float64
	Image       string
	Score       float64
	Highlights  map[string][]string
	IsPromoted  bool
	IsFeatured  bool
}

// UniversalSearch performs advanced search with function_score ranking
func (c *Client) UniversalSearch(ctx context.Context, params *UniversalSearchParams) (*UniversalSearchResult, error) {
	// Build multi_match query with boost
	multiMatchQuery := map[string]interface{}{
		"multi_match": map[string]interface{}{
			"query": params.Query,
			"fields": []string{
				"title^3",
				"title.autocomplete^2",
				"description^1",
				"brand^2",
				"attributes_searchable_text^0.5",
				"tags^1.5",
				"storefront_name^1",
			},
			"type":      "best_fields",
			"fuzziness": "AUTO",
		},
	}

	if !params.EnableFuzzy {
		delete(multiMatchQuery["multi_match"].(map[string]interface{}), "fuzziness")
	}

	// Build filters
	filters := []map[string]interface{}{
		{"term": map[string]interface{}{"status": "active"}},
		{"term": map[string]interface{}{"visibility": "public"}},
	}

	// Category filter
	if len(params.CategoryIDs) > 0 {
		filters = append(filters, map[string]interface{}{
			"terms": map[string]interface{}{"category_id": params.CategoryIDs},
		})
	}

	// Price range filter
	if params.PriceMin != nil || params.PriceMax != nil {
		priceRange := map[string]interface{}{}
		if params.PriceMin != nil {
			priceRange["gte"] = *params.PriceMin
		}
		if params.PriceMax != nil {
			priceRange["lte"] = *params.PriceMax
		}
		filters = append(filters, map[string]interface{}{
			"range": map[string]interface{}{"price": priceRange},
		})
	}

	// Source type filter (c2c, b2c)
	if len(params.SourceTypes) > 0 {
		filters = append(filters, map[string]interface{}{
			"terms": map[string]interface{}{"source_type": params.SourceTypes},
		})
	}

	// Location filters
	if params.City != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{"city.keyword": params.City},
		})
	}
	if params.Country != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{"country.keyword": params.Country},
		})
	}

	// Attribute filters (nested)
	for attrCode, attrValues := range params.Attributes {
		if len(attrValues) > 0 {
			filters = append(filters, map[string]interface{}{
				"nested": map[string]interface{}{
					"path": "attributes",
					"query": map[string]interface{}{
						"bool": map[string]interface{}{
							"must": []map[string]interface{}{
								{"term": map[string]interface{}{"attributes.code": attrCode}},
								{"terms": map[string]interface{}{"attributes.value_text.keyword": attrValues}},
							},
						},
					},
				},
			})
		}
	}

	// Build function_score query for ranking
	functionScoreQuery := map[string]interface{}{
		"function_score": map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must":   multiMatchQuery,
					"filter": filters,
				},
			},
			"functions": []map[string]interface{}{
				// Promoted listings get 1.5x boost
				{
					"filter": map[string]interface{}{"term": map[string]interface{}{"is_promoted": true}},
					"weight": 1.5,
				},
				// Featured listings get 1.3x boost
				{
					"filter": map[string]interface{}{"term": map[string]interface{}{"is_featured": true}},
					"weight": 1.3,
				},
				// Verified sellers get 1.2x boost
				{
					"filter": map[string]interface{}{"term": map[string]interface{}{"seller_verified": true}},
					"weight": 1.2,
				},
				// High-rated listings (≥4.5 stars) get 1.1x boost
				{
					"filter": map[string]interface{}{
						"range": map[string]interface{}{"rating": map[string]interface{}{"gte": 4.5}},
					},
					"weight": 1.1,
				},
				// Popularity score (views, favorites) - logarithmic decay
				{
					"field_value_factor": map[string]interface{}{
						"field":    "popularity_score",
						"factor":   1.2,
						"modifier": "log1p",
						"missing":  0,
					},
				},
				// Recency boost (newer listings preferred within 7 days)
				{
					"gauss": map[string]interface{}{
						"created_at": map[string]interface{}{
							"origin": "now",
							"scale":  "7d",
							"decay":  0.5,
						},
					},
				},
			},
			"score_mode": "sum",      // Sum all function scores
			"boost_mode": "multiply", // Multiply with query score
		},
	}

	// Calculate offset
	offset := (params.Page - 1) * params.Limit
	if offset < 0 {
		offset = 0
	}

	// Build final query
	query := map[string]interface{}{
		"query": functionScoreQuery,
		"size":  params.Limit,
		"from":  offset,
	}

	// Add highlighting if requested
	if params.Highlighting {
		query["highlight"] = map[string]interface{}{
			"fields": map[string]interface{}{
				"title": map[string]interface{}{
					"number_of_fragments": 1,
				},
				"description": map[string]interface{}{
					"number_of_fragments": 3,
					"fragment_size":       150,
				},
			},
			"pre_tags":  []string{"<mark>"},
			"post_tags": []string{"</mark>"},
		}
	}

	// Add sorting
	query["sort"] = c.buildSort(params.SortBy, params.SortOrder, params.Query != "")

	// Execute search via OpenSearch client
	queryJSON, err := json.Marshal(query)
	if err != nil {
		c.logger.Error().Err(err).Msg("failed to marshal query")
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	res, err := c.client.Search(
		c.client.Search.WithContext(ctx),
		c.client.Search.WithIndex(c.index),
		c.client.Search.WithBody(bytes.NewReader(queryJSON)),
		c.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		c.logger.Error().Err(err).Msg("universal search failed")
		return nil, fmt.Errorf("universal search failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		c.logger.Error().Str("status", res.Status()).Msg("search returned error")
		return nil, fmt.Errorf("search returned error: %s", res.Status())
	}

	// Decode response
	var result AdvancedSearchResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		c.logger.Error().Err(err).Msg("failed to decode response")
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Parse results
	listings := make([]*ListingResult, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		listing := c.parseListingResultWithHighlight(hit)
		listings = append(listings, listing)
	}

	return &UniversalSearchResult{
		Listings: listings,
		Total:    result.Hits.Total.Value,
		Page:     params.Page,
		Limit:    params.Limit,
		TookMs:   result.Took,
	}, nil
}

// DidYouMean provides phrase suggestions using phrase suggester
func (c *Client) DidYouMean(ctx context.Context, query string) ([]string, error) {
	// Build phrase suggester query
	suggestQuery := map[string]interface{}{
		"suggest": map[string]interface{}{
			"did_you_mean": map[string]interface{}{
				"text": query,
				"phrase": map[string]interface{}{
					"field":     "title.trigram",
					"size":      3,
					"gram_size": 3,
					"direct_generator": []map[string]interface{}{
						{
							"field":        "title.trigram",
							"suggest_mode": "always",
						},
					},
					"highlight": map[string]interface{}{
						"pre_tag":  "<em>",
						"post_tag": "</em>",
					},
				},
			},
		},
	}

	// Execute search
	queryJSON, err := json.Marshal(suggestQuery)
	if err != nil {
		c.logger.Error().Err(err).Msg("failed to marshal suggest query")
		return nil, fmt.Errorf("failed to marshal suggest query: %w", err)
	}

	res, err := c.client.Search(
		c.client.Search.WithContext(ctx),
		c.client.Search.WithIndex(c.index),
		c.client.Search.WithBody(bytes.NewReader(queryJSON)),
	)
	if err != nil {
		c.logger.Error().Err(err).Msg("did you mean query failed")
		return nil, fmt.Errorf("did you mean query failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		c.logger.Error().Str("status", res.Status()).Msg("suggest returned error")
		return nil, fmt.Errorf("suggest returned error: %s", res.Status())
	}

	// Decode response
	var result AdvancedSearchResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		c.logger.Error().Err(err).Msg("failed to decode suggest response")
		return nil, fmt.Errorf("failed to decode suggest response: %w", err)
	}

	// Parse suggestions
	suggestions := []string{}
	if result.Suggest != nil {
		if didYouMean, ok := result.Suggest["did_you_mean"].([]interface{}); ok {
			for _, item := range didYouMean {
				if itemMap, ok := item.(map[string]interface{}); ok {
					if options, ok := itemMap["options"].([]interface{}); ok {
						for _, option := range options {
							if optMap, ok := option.(map[string]interface{}); ok {
								if text, ok := optMap["text"].(string); ok {
									suggestions = append(suggestions, text)
								}
							}
						}
					}
				}
			}
		}
	}

	c.logger.Debug().
		Str("query", query).
		Int("suggestions", len(suggestions)).
		Msg("did you mean completed")

	return suggestions, nil
}

// AutocompleteSuggestion represents a single autocomplete suggestion
type AutocompleteSuggestion struct {
	Text      string
	ID        int64
	Highlight string
}

// Autocomplete provides autocomplete suggestions with highlighting
func (c *Client) Autocomplete(ctx context.Context, prefix string, limit int) ([]AutocompleteSuggestion, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	// Build match_phrase_prefix query on autocomplete field
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match_phrase_prefix": map[string]interface{}{
							"title.autocomplete": map[string]interface{}{
								"query":          prefix,
								"max_expansions": 50,
							},
						},
					},
				},
				"filter": []map[string]interface{}{
					{"term": map[string]interface{}{"status": "active"}},
					{"term": map[string]interface{}{"visibility": "public"}},
				},
			},
		},
		"size":    limit,
		"_source": []string{"id", "title", "images"},
		"highlight": map[string]interface{}{
			"fields": map[string]interface{}{
				"title": map[string]interface{}{},
			},
			"pre_tags":  []string{"<mark>"},
			"post_tags": []string{"</mark>"},
		},
	}

	// Execute search
	queryJSON, err := json.Marshal(query)
	if err != nil {
		c.logger.Error().Err(err).Msg("failed to marshal autocomplete query")
		return nil, fmt.Errorf("failed to marshal autocomplete query: %w", err)
	}

	res, err := c.client.Search(
		c.client.Search.WithContext(ctx),
		c.client.Search.WithIndex(c.index),
		c.client.Search.WithBody(bytes.NewReader(queryJSON)),
		c.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		c.logger.Error().Err(err).Msg("autocomplete query failed")
		return nil, fmt.Errorf("autocomplete query failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		c.logger.Error().Str("status", res.Status()).Msg("autocomplete returned error")
		return nil, fmt.Errorf("autocomplete returned error: %s", res.Status())
	}

	// Decode response
	var result AdvancedSearchResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		c.logger.Error().Err(err).Msg("failed to decode autocomplete response")
		return nil, fmt.Errorf("failed to decode autocomplete response: %w", err)
	}

	// Parse results
	suggestions := []AutocompleteSuggestion{}
	for _, hit := range result.Hits.Hits {
		suggestion := AutocompleteSuggestion{}

		// Parse ID
		if id, ok := hit.Source["id"].(float64); ok {
			suggestion.ID = int64(id)
		}

		// Parse title
		if title, ok := hit.Source["title"].(string); ok {
			suggestion.Text = title
		}

		// Parse highlight
		if hit.Highlight != nil {
			if titleHighlight, ok := hit.Highlight["title"].([]interface{}); ok {
				if len(titleHighlight) > 0 {
					if hl, ok := titleHighlight[0].(string); ok {
						suggestion.Highlight = hl
					}
				}
			}
		}

		suggestions = append(suggestions, suggestion)
	}

	c.logger.Debug().
		Str("prefix", prefix).
		Int("suggestions", len(suggestions)).
		Msg("autocomplete completed")

	return suggestions, nil
}

// ============================================================================
// Helper Methods
// ============================================================================

// buildSort constructs sort configuration based on parameters
func (c *Client) buildSort(sortBy, sortOrder string, hasQuery bool) []map[string]interface{} {
	// Default sort order
	if sortOrder == "" {
		sortOrder = "desc"
	}

	// Build sort based on field
	switch sortBy {
	case "price":
		return []map[string]interface{}{
			{"price": map[string]interface{}{"order": sortOrder}},
			{"_score": map[string]interface{}{"order": "desc"}},
		}
	case "created_at":
		return []map[string]interface{}{
			{"created_at": map[string]interface{}{"order": sortOrder}},
			{"_score": map[string]interface{}{"order": "desc"}},
		}
	case "popularity":
		return []map[string]interface{}{
			{"popularity_score": map[string]interface{}{"order": sortOrder}},
			{"_score": map[string]interface{}{"order": "desc"}},
		}
	case "relevance":
		return []map[string]interface{}{
			{"_score": map[string]interface{}{"order": "desc"}},
			{"created_at": map[string]interface{}{"order": "desc"}},
		}
	default:
		// Default: if query provided, sort by relevance, else by created_at
		if hasQuery {
			return []map[string]interface{}{
				{"_score": map[string]interface{}{"order": "desc"}},
				{"created_at": map[string]interface{}{"order": "desc"}},
			}
		}
		return []map[string]interface{}{
			{"created_at": map[string]interface{}{"order": "desc"}},
		}
	}
}

// parseListingResultWithHighlight parses a single hit with highlights into ListingResult
func (c *Client) parseListingResultWithHighlight(hit SearchHitWithHighlight) *ListingResult {
	listing := &ListingResult{}

	// Parse ID
	if id, ok := hit.Source["id"].(float64); ok {
		listing.ID = int64(id)
	}

	// Parse UUID
	if uuid, ok := hit.Source["uuid"].(string); ok {
		listing.UUID = uuid
	}

	// Parse title
	if title, ok := hit.Source["title"].(string); ok {
		listing.Title = title
	}

	// Parse price
	if price, ok := hit.Source["price"].(float64); ok {
		listing.Price = price
	}

	// Parse score
	if hit.Score != nil {
		listing.Score = *hit.Score
	}

	// Parse promoted flag
	if isPromoted, ok := hit.Source["is_promoted"].(bool); ok {
		listing.IsPromoted = isPromoted
	}

	// Parse featured flag
	if isFeatured, ok := hit.Source["is_featured"].(bool); ok {
		listing.IsFeatured = isFeatured
	}

	// Parse image (first image)
	if images, ok := hit.Source["images"].([]interface{}); ok {
		if len(images) > 0 {
			if img, ok := images[0].(map[string]interface{}); ok {
				if url, ok := img["public_url"].(string); ok {
					listing.Image = url
				}
			}
		}
	}

	// Parse highlights
	if hit.Highlight != nil {
		listing.Highlights = make(map[string][]string)
		for field, highlights := range hit.Highlight {
			if hlArray, ok := highlights.([]interface{}); ok {
				strings := make([]string, 0, len(hlArray))
				for _, hl := range hlArray {
					if hlStr, ok := hl.(string); ok {
						strings = append(strings, hlStr)
					}
				}
				listing.Highlights[field] = strings
			}
		}
	}

	return listing
}

// SearchHitWithHighlight represents a single search result hit with highlights
type SearchHitWithHighlight struct {
	ID        string                 `json:"_id"`
	Score     *float64               `json:"_score"`
	Source    map[string]interface{} `json:"_source"`
	Highlight map[string]interface{} `json:"highlight,omitempty"`
}

// AdvancedSearchResponse extends SearchResponse with highlights
type AdvancedSearchResponse struct {
	Took int `json:"took"`
	Hits struct {
		Total struct {
			Value    int64  `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		Hits []SearchHitWithHighlight `json:"hits"`
	} `json:"hits"`
	Aggregations map[string]interface{} `json:"aggregations,omitempty"`
	Suggest      map[string]interface{} `json:"suggest,omitempty"`
}
