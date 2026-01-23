package search

import (
	"fmt"
)

// ============================================================================
// PHASE 21.2: OpenSearch DSL Query Builders
// ============================================================================

// BuildFacetsQuery builds an aggregations-only query for GetSearchFacets
func BuildFacetsQuery(req *FacetsRequest) map[string]interface{} {
	query := map[string]interface{}{
		"size": 0, // No documents, only aggregations
		"aggs": buildAggregations(),
	}

	// Add filters if provided (pre-filter before aggregating)
	boolQuery := map[string]interface{}{
		"must": []map[string]interface{}{
			{"term": map[string]interface{}{"status": "active"}},
		},
	}

	if req.Query != "" {
		boolQuery["must"] = append(
			boolQuery["must"].([]map[string]interface{}),
			map[string]interface{}{
				"multi_match": map[string]interface{}{
					"query":  req.Query,
					"fields": []string{"title^3", "description"},
					"type":   "best_fields",
				},
			},
		)
	}

	if req.CategoryID != nil {
		boolQuery["must"] = append(
			boolQuery["must"].([]map[string]interface{}),
			map[string]interface{}{
				"term": map[string]interface{}{"category_id": *req.CategoryID},
			},
		)
	}

	if req.Filters != nil {
		filterClauses := buildFilterClauses(req.Filters)
		boolQuery["must"] = append(boolQuery["must"].([]map[string]interface{}), filterClauses...)
	}

	query["query"] = map[string]interface{}{"bool": boolQuery}

	return query
}

// buildAggregations constructs all facet aggregations
func buildAggregations() map[string]interface{} {
	return map[string]interface{}{
		// Category distribution
		"categories": map[string]interface{}{
			"terms": map[string]interface{}{
				"field": "category_id",
				"size":  50,
			},
		},

		// Price histogram (buckets of 100)
		"price_ranges": map[string]interface{}{
			"histogram": map[string]interface{}{
				"field":    "price",
				"interval": 100,
			},
		},

		// Source types (c2c, b2c)
		"source_types": map[string]interface{}{
			"terms": map[string]interface{}{
				"field": "source_type",
			},
		},

		// Stock statuses
		"stock_statuses": map[string]interface{}{
			"terms": map[string]interface{}{
				"field": "stock_status",
			},
		},

		// Attributes (nested aggregation)
		"attributes": map[string]interface{}{
			"nested": map[string]interface{}{
				"path": "attributes",
			},
			"aggs": map[string]interface{}{
				"attribute_keys": map[string]interface{}{
					"terms": map[string]interface{}{
						"field": "attributes.code",
						"size":  50,
					},
					"aggs": map[string]interface{}{
						"attribute_values": map[string]interface{}{
							"terms": map[string]interface{}{
								"field": "attributes.value_text.keyword",
								"size":  20,
							},
						},
					},
				},
			},
		},
	}
}

// BuildFilteredSearchQuery builds a complex query with filters, sorting, and optional facets
func BuildFilteredSearchQuery(req *SearchFiltersRequest) map[string]interface{} {
	mustClauses := []map[string]interface{}{
		{"term": map[string]interface{}{"status": "active"}},
	}

	// Text search
	if req.Query != "" {
		mustClauses = append(mustClauses, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  req.Query,
				"fields": []string{"title^3", "description"},
				"type":   "best_fields",
			},
		})
	}

	// Category filter
	if req.CategoryID != nil {
		mustClauses = append(mustClauses, map[string]interface{}{
			"term": map[string]interface{}{"category_id": *req.CategoryID},
		})
	}

	// Advanced filters
	if req.Filters != nil {
		filterClauses := buildFilterClauses(req.Filters)
		mustClauses = append(mustClauses, filterClauses...)
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": mustClauses,
			},
		},
		"size": req.Limit,
		"from": req.Offset,
	}

	// Sorting
	if req.Sort != nil {
		query["sort"] = buildSort(req.Sort)
	} else {
		// Default sort
		if req.Query != "" {
			// With query: sort by relevance (score)
			query["sort"] = []map[string]interface{}{
				{"_score": map[string]interface{}{"order": "desc"}},
			}
		} else {
			// Without query: sort by created_at
			query["sort"] = []map[string]interface{}{
				{"created_at": map[string]interface{}{"order": "desc"}},
			}
		}
	}

	// Include facets if requested
	if req.IncludeFacets {
		query["aggs"] = buildAggregations()
	}

	return query
}

// buildFilterClauses constructs filter clauses from SearchFilters
func buildFilterClauses(filters *SearchFilters) []map[string]interface{} {
	clauses := []map[string]interface{}{}

	// Price range filter
	if filters.Price != nil {
		rangeClause := map[string]interface{}{}
		if filters.Price.Min != nil {
			rangeClause["gte"] = *filters.Price.Min
		}
		if filters.Price.Max != nil {
			rangeClause["lte"] = *filters.Price.Max
		}
		if len(rangeClause) > 0 {
			clauses = append(clauses, map[string]interface{}{
				"range": map[string]interface{}{
					"price": rangeClause,
				},
			})
		}
	}

	// Attribute filters (nested queries)
	for key, values := range filters.Attributes {
		if len(values) > 0 {
			clauses = append(clauses, map[string]interface{}{
				"nested": map[string]interface{}{
					"path": "attributes",
					"query": map[string]interface{}{
						"bool": map[string]interface{}{
							"must": []map[string]interface{}{
								{"term": map[string]interface{}{"attributes.code": key}},
								{"terms": map[string]interface{}{"attributes.value_text.keyword": values}},
							},
						},
					},
				},
			})
		}
	}

	// Location filter (geo_distance)
	if filters.Location != nil {
		clauses = append(clauses, map[string]interface{}{
			"geo_distance": map[string]interface{}{
				"distance": fmt.Sprintf("%.1fkm", filters.Location.RadiusKm),
				"location": map[string]interface{}{
					"lat": filters.Location.Lat,
					"lon": filters.Location.Lon,
				},
			},
		})
	}

	// Source type filter
	if filters.SourceType != nil {
		clauses = append(clauses, map[string]interface{}{
			"term": map[string]interface{}{"source_type": *filters.SourceType},
		})
	}

	// Stock status filter
	if filters.StockStatus != nil {
		clauses = append(clauses, map[string]interface{}{
			"term": map[string]interface{}{"stock_status": *filters.StockStatus},
		})
	}

	return clauses
}

// buildSort constructs sort configuration
func buildSort(sort *SortConfig) []map[string]interface{} {
	return []map[string]interface{}{
		{sort.Field: map[string]interface{}{"order": sort.Order}},
	}
}

// BuildSuggestionsQuery builds a completion suggester query
func BuildSuggestionsQuery(req *SuggestionsRequest) map[string]interface{} {
	query := map[string]interface{}{
		"suggest": map[string]interface{}{
			"listing-suggest": map[string]interface{}{
				"prefix": req.Prefix,
				"completion": map[string]interface{}{
					"field":           "suggest",
					"size":            req.Limit,
					"skip_duplicates": true,
					"fuzzy": map[string]interface{}{
						"fuzziness": "AUTO",
					},
				},
			},
		},
	}

	// Add category context filter if provided
	if req.CategoryID != nil {
		completion := query["suggest"].(map[string]interface{})["listing-suggest"].(map[string]interface{})["completion"].(map[string]interface{})
		completion["contexts"] = map[string]interface{}{
			"category": []string{*req.CategoryID}, // UUID string
		}
	}

	return query
}

// BuildPopularSearchesQuery builds a query for trending searches
// Note: This would typically query a separate index tracking search events
// For now, we return a placeholder structure that can be implemented with actual analytics
func BuildPopularSearchesQuery(req *PopularSearchesRequest) map[string]interface{} {
	// This would query a "search_analytics" index with documents like:
	// {
	//   "query": "laptop",
	//   "category_id": 1001,
	//   "timestamp": "2025-11-17T10:00:00Z",
	//   "user_id": 123
	// }

	mustClauses := []map[string]interface{}{}

	// Category filter
	if req.CategoryID != nil {
		mustClauses = append(mustClauses, map[string]interface{}{
			"term": map[string]interface{}{"category_id": *req.CategoryID},
		})
	}

	// Time range filter
	var timeRangeStr string
	switch req.TimeRange {
	case "24h":
		timeRangeStr = "now-24h"
	case "7d":
		timeRangeStr = "now-7d"
	case "30d":
		timeRangeStr = "now-30d"
	default:
		timeRangeStr = "now-24h"
	}

	mustClauses = append(mustClauses, map[string]interface{}{
		"range": map[string]interface{}{
			"timestamp": map[string]interface{}{
				"gte": timeRangeStr,
			},
		},
	})

	query := map[string]interface{}{
		"size": 0, // No documents, only aggregations
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": mustClauses,
			},
		},
		"aggs": map[string]interface{}{
			"popular_queries": map[string]interface{}{
				"terms": map[string]interface{}{
					"field": "query.keyword",
					"size":  req.Limit,
					"order": map[string]interface{}{
						"_count": "desc",
					},
				},
			},
		},
	}

	return query
}

// ValidateQuery performs basic validation on generated query
func ValidateQuery(query map[string]interface{}) error {
	// Check if query has required structure
	if query == nil {
		return fmt.Errorf("query cannot be nil")
	}

	// Check for either "query" or "suggest" field
	_, hasQuery := query["query"]
	_, hasSuggest := query["suggest"]
	_, hasAggs := query["aggs"]

	if !hasQuery && !hasSuggest && !hasAggs {
		return fmt.Errorf("query must have at least one of: query, suggest, or aggs")
	}

	return nil
}
