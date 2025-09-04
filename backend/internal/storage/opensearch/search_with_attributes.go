package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/opensearch-project/opensearch-go/v2"
	"go.uber.org/zap"
)

// AttributeSearchService provides optimized search with attributes
type AttributeSearchService struct {
	client *opensearch.Client
	logger *zap.Logger
	index  string
}

// NewAttributeSearchService creates a new attribute search service
func NewAttributeSearchService(client *opensearch.Client, logger *zap.Logger, index string) *AttributeSearchService {
	return &AttributeSearchService{
		client: client,
		logger: logger,
		index:  index,
	}
}

// SearchRequest represents an optimized search request with attributes
type SearchRequest struct {
	Query      string              `json:"query"`
	CategoryID *int                `json:"category_id,omitempty"`
	MinPrice   *float64            `json:"min_price,omitempty"`
	MaxPrice   *float64            `json:"max_price,omitempty"`
	Attributes map[string][]string `json:"attributes,omitempty"`
	Sort       string              `json:"sort"`
	From       int                 `json:"from"`
	Size       int                 `json:"size"`

	// Optimization flags
	IncludeAggregations bool `json:"include_aggregations"`
	OptimizeForCache    bool `json:"optimize_for_cache"`
}

// SearchResponse represents the search response
type SearchResponse struct {
	Hits         []json.RawMessage      `json:"hits"`
	Total        int64                  `json:"total"`
	Aggregations map[string]interface{} `json:"aggregations,omitempty"`
	TookMs       int64                  `json:"took_ms"`
}

// Search performs an optimized search with attributes
func (s *AttributeSearchService) Search(ctx context.Context, req SearchRequest) (*SearchResponse, error) {
	// Build optimized query
	query := s.buildOptimizedQuery(req)

	// Execute search
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("failed to encode query: %w", err)
	}

	res, err := s.client.Search(
		s.client.Search.WithContext(ctx),
		s.client.Search.WithIndex(s.index),
		s.client.Search.WithBody(&buf),
		s.client.Search.WithTrackTotalHits(true),
		s.client.Search.WithFrom(req.From),
		s.client.Search.WithSize(req.Size),
		s.client.Search.WithRequestCache(req.OptimizeForCache), // Enable request cache
	)
	if err != nil {
		return nil, fmt.Errorf("search request failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.String())
	}

	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return s.parseSearchResponse(result), nil
}

// buildOptimizedQuery builds an optimized OpenSearch query
func (s *AttributeSearchService) buildOptimizedQuery(req SearchRequest) map[string]interface{} {
	query := map[string]interface{}{
		"track_total_hits": true,
		"_source": []string{
			"id", "title", "description", "price", "currency",
			"category_id", "user_id", "location", "images",
			"created_at", "updated_at", "attributes",
		},
	}

	// Build bool query
	boolQuery := map[string]interface{}{
		"must":   []interface{}{},
		"filter": []interface{}{},
	}

	// Add text search if query provided
	if req.Query != "" {
		boolQuery["must"] = append(boolQuery["must"].([]interface{}), map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query": req.Query,
				"fields": []string{
					"title^3",       // Title has highest boost
					"description^2", // Description medium boost
					"attributes.*",  // Search in all attributes
				},
				"type":      "best_fields",
				"operator":  "or",
				"fuzziness": "AUTO",
			},
		})
	}

	// Add category filter
	if req.CategoryID != nil {
		boolQuery["filter"] = append(boolQuery["filter"].([]interface{}), map[string]interface{}{
			"term": map[string]interface{}{
				"category_id": *req.CategoryID,
			},
		})
	}

	// Add price range filter
	if req.MinPrice != nil || req.MaxPrice != nil {
		rangeQuery := map[string]interface{}{}
		if req.MinPrice != nil {
			rangeQuery["gte"] = *req.MinPrice
		}
		if req.MaxPrice != nil {
			rangeQuery["lte"] = *req.MaxPrice
		}

		boolQuery["filter"] = append(boolQuery["filter"].([]interface{}), map[string]interface{}{
			"range": map[string]interface{}{
				"price": rangeQuery,
			},
		})
	}

	// Add attribute filters with optimization
	if len(req.Attributes) > 0 {
		for attrKey, attrValues := range req.Attributes {
			if len(attrValues) == 0 {
				continue
			}

			// Use nested query for attributes
			nestedQuery := map[string]interface{}{
				"nested": map[string]interface{}{
					"path": "attributes",
					"query": map[string]interface{}{
						"bool": map[string]interface{}{
							"must": []interface{}{
								map[string]interface{}{
									"term": map[string]interface{}{
										"attributes.key": attrKey,
									},
								},
							},
						},
					},
				},
			}

			// Handle multiple values with should clause
			if len(attrValues) == 1 {
				nestedQuery["nested"].(map[string]interface{})["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = append(
					nestedQuery["nested"].(map[string]interface{})["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]interface{}),
					map[string]interface{}{
						"term": map[string]interface{}{
							"attributes.value.keyword": attrValues[0],
						},
					},
				)
			} else {
				shouldClauses := make([]interface{}, len(attrValues))
				for i, val := range attrValues {
					shouldClauses[i] = map[string]interface{}{
						"term": map[string]interface{}{
							"attributes.value.keyword": val,
						},
					}
				}

				nestedQuery["nested"].(map[string]interface{})["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = append(
					nestedQuery["nested"].(map[string]interface{})["query"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]interface{}),
					map[string]interface{}{
						"bool": map[string]interface{}{
							"should":               shouldClauses,
							"minimum_should_match": 1,
						},
					},
				)
			}

			boolQuery["filter"] = append(boolQuery["filter"].([]interface{}), nestedQuery)
		}
	}

	// Set the query
	if len(boolQuery["must"].([]interface{})) > 0 || len(boolQuery["filter"].([]interface{})) > 0 {
		query["query"] = map[string]interface{}{
			"bool": boolQuery,
		}
	} else {
		query["query"] = map[string]interface{}{
			"match_all": map[string]interface{}{},
		}
	}

	// Add sorting
	query["sort"] = s.buildSortCriteria(req.Sort)

	// Add aggregations if requested
	if req.IncludeAggregations {
		query["aggs"] = s.buildAggregations(req)
	}

	// Add performance optimizations
	query["timeout"] = "5s"
	query["terminate_after"] = 10000 // Limit documents examined

	return query
}

// buildSortCriteria builds sort criteria for the query
func (s *AttributeSearchService) buildSortCriteria(sortBy string) []interface{} {
	switch sortBy {
	case "price_asc":
		return []interface{}{
			map[string]interface{}{
				"price": map[string]interface{}{
					"order":   "asc",
					"missing": "_last",
				},
			},
		}
	case "price_desc":
		return []interface{}{
			map[string]interface{}{
				"price": map[string]interface{}{
					"order":   "desc",
					"missing": "_last",
				},
			},
		}
	case "date_desc":
		return []interface{}{
			map[string]interface{}{
				"created_at": map[string]interface{}{
					"order": "desc",
				},
			},
		}
	case "relevance":
		fallthrough
	default:
		return []interface{}{
			"_score",
			map[string]interface{}{
				"created_at": map[string]interface{}{
					"order": "desc",
				},
			},
		}
	}
}

// buildAggregations builds aggregations for getting available filter values
func (s *AttributeSearchService) buildAggregations(req SearchRequest) map[string]interface{} {
	aggs := map[string]interface{}{
		"categories": map[string]interface{}{
			"terms": map[string]interface{}{
				"field": "category_id",
				"size":  20,
			},
		},
		"price_range": map[string]interface{}{
			"stats": map[string]interface{}{
				"field": "price",
			},
		},
		"attributes": map[string]interface{}{
			"nested": map[string]interface{}{
				"path": "attributes",
			},
			"aggs": map[string]interface{}{
				"attribute_keys": map[string]interface{}{
					"terms": map[string]interface{}{
						"field": "attributes.key.keyword",
						"size":  50,
					},
					"aggs": map[string]interface{}{
						"attribute_values": map[string]interface{}{
							"terms": map[string]interface{}{
								"field": "attributes.value.keyword",
								"size":  20,
							},
						},
					},
				},
			},
		},
	}

	return aggs
}

// parseSearchResponse parses the OpenSearch response
func (s *AttributeSearchService) parseSearchResponse(result map[string]interface{}) *SearchResponse {
	response := &SearchResponse{
		Hits:         []json.RawMessage{},
		Aggregations: map[string]interface{}{},
	}

	// Parse hits
	if hits, ok := result["hits"].(map[string]interface{}); ok {
		// Get total count
		if total, ok := hits["total"].(map[string]interface{}); ok {
			if value, ok := total["value"].(float64); ok {
				response.Total = int64(value)
			}
		}

		// Get hit documents
		if docs, ok := hits["hits"].([]interface{}); ok {
			for _, doc := range docs {
				if docMap, ok := doc.(map[string]interface{}); ok {
					if source, ok := docMap["_source"]; ok {
						if jsonBytes, err := json.Marshal(source); err == nil {
							response.Hits = append(response.Hits, jsonBytes)
						}
					}
				}
			}
		}
	}

	// Parse took time
	if took, ok := result["took"].(float64); ok {
		response.TookMs = int64(took)
	}

	// Parse aggregations
	if aggs, ok := result["aggregations"].(map[string]interface{}); ok {
		response.Aggregations = aggs
	}

	return response
}

// BulkIndexAttributes performs bulk indexing of listings with attributes
func (s *AttributeSearchService) BulkIndexAttributes(ctx context.Context, documents []map[string]interface{}) error {
	if len(documents) == 0 {
		return nil
	}

	var buf bytes.Buffer
	for _, doc := range documents {
		// Add index action
		meta := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": s.index,
				"_id":    doc["id"],
			},
		}

		if err := json.NewEncoder(&buf).Encode(meta); err != nil {
			return fmt.Errorf("failed to encode meta: %w", err)
		}

		if err := json.NewEncoder(&buf).Encode(doc); err != nil {
			return fmt.Errorf("failed to encode document: %w", err)
		}
	}

	// Execute bulk request
	res, err := s.client.Bulk(
		bytes.NewReader(buf.Bytes()),
		s.client.Bulk.WithContext(ctx),
		s.client.Bulk.WithRefresh("wait_for"), // Wait for refresh
	)
	if err != nil {
		return fmt.Errorf("bulk index failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("bulk index error: %s", res.String())
	}

	// Check for errors in response
	var bulkRes map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&bulkRes); err != nil {
		return fmt.Errorf("failed to decode bulk response: %w", err)
	}

	if errors, ok := bulkRes["errors"].(bool); ok && errors {
		return fmt.Errorf("bulk index had errors")
	}

	s.logger.Info("Bulk indexed documents",
		zap.Int("count", len(documents)),
		zap.Int64("took_ms", int64(bulkRes["took"].(float64))),
	)

	return nil
}

// UpdateMapping updates the index mapping for better attribute search
func (s *AttributeSearchService) UpdateMapping(ctx context.Context) error {
	mapping := map[string]interface{}{
		"properties": map[string]interface{}{
			"attributes": map[string]interface{}{
				"type": "nested",
				"properties": map[string]interface{}{
					"key": map[string]interface{}{
						"type": "text",
						"fields": map[string]interface{}{
							"keyword": map[string]interface{}{
								"type":         "keyword",
								"ignore_above": 256,
							},
						},
					},
					"value": map[string]interface{}{
						"type": "text",
						"fields": map[string]interface{}{
							"keyword": map[string]interface{}{
								"type":         "keyword",
								"ignore_above": 256,
							},
						},
					},
					"type": map[string]interface{}{
						"type": "keyword",
					},
				},
			},
			"category_id": map[string]interface{}{
				"type": "integer",
			},
			"price": map[string]interface{}{
				"type": "float",
			},
			"created_at": map[string]interface{}{
				"type": "date",
			},
			"title": map[string]interface{}{
				"type":     "text",
				"analyzer": "standard",
			},
			"description": map[string]interface{}{
				"type":     "text",
				"analyzer": "standard",
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(mapping); err != nil {
		return fmt.Errorf("failed to encode mapping: %w", err)
	}

	res, err := s.client.Indices.PutMapping(
		strings.NewReader(buf.String()),
		s.client.Indices.PutMapping.WithIndex(s.index),
		s.client.Indices.PutMapping.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("failed to update mapping: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("mapping update error: %s", res.String())
	}

	s.logger.Info("Updated index mapping for optimized attribute search")
	return nil
}
