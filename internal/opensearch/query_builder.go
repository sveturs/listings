package opensearch

import (
	"fmt"
)

// AttributeFilter represents a filter for a specific attribute
type AttributeFilter struct {
	Code   string   // Attribute code (e.g., "brand", "color", "ram")
	Type   string   // Attribute type (e.g., "select", "number", "boolean")
	Values []string // Filter values (can be multiple for OR logic)
}

// FilterQuery represents query configuration for search with filters
type FilterQuery struct {
	// Required
	CategoryID int64

	// Optional filters
	PriceMin     *float64
	PriceMax     *float64
	SourceType   string // "c2c" or "b2c"
	StockStatus  string // "in_stock", "out_of_stock", "low_stock"
	LatLon       *GeoLocation
	RadiusKM     *float64
	Attributes   []AttributeFilter
	SearchQuery  string // Text search query

	// Pagination and sorting
	Limit  int
	Offset int
	SortBy string // "relevance", "price", "created_at", "views_count", "favorites_count"
	SortOrder string // "asc", "desc"
}

// GeoLocation represents latitude and longitude
type GeoLocation struct {
	Lat float64
	Lon float64
}

// BuildFilterQuery constructs OpenSearch DSL query from filter configuration
func BuildFilterQuery(cfg FilterQuery) map[string]interface{} {
	// Build bool query with must clauses
	mustClauses := []map[string]interface{}{}

	// Category filter (required if specified)
	if cfg.CategoryID > 0 {
		mustClauses = append(mustClauses, map[string]interface{}{
			"term": map[string]interface{}{
				"category_id": cfg.CategoryID,
			},
		})
	}

	// Source type filter
	if cfg.SourceType != "" {
		mustClauses = append(mustClauses, map[string]interface{}{
			"term": map[string]interface{}{
				"source_type": cfg.SourceType,
			},
		})
	}

	// Stock status filter
	if cfg.StockStatus != "" {
		mustClauses = append(mustClauses, map[string]interface{}{
			"term": map[string]interface{}{
				"stock_status": cfg.StockStatus,
			},
		})
	}

	// Text search query (full-text search in title and description)
	if cfg.SearchQuery != "" {
		mustClauses = append(mustClauses, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  cfg.SearchQuery,
				"fields": []string{"title^3", "description", "title.autocomplete"},
				"type":   "best_fields",
				"operator": "and",
			},
		})
	}

	// Price range filter
	if cfg.PriceMin != nil || cfg.PriceMax != nil {
		priceRange := buildRangeFilter("price", cfg.PriceMin, cfg.PriceMax)
		mustClauses = append(mustClauses, priceRange)
	}

	// Geo distance filter
	if cfg.LatLon != nil && cfg.RadiusKM != nil {
		mustClauses = append(mustClauses, map[string]interface{}{
			"geo_distance": map[string]interface{}{
				"distance": fmt.Sprintf("%.2fkm", *cfg.RadiusKM),
				"location": map[string]interface{}{
					"lat": cfg.LatLon.Lat,
					"lon": cfg.LatLon.Lon,
				},
			},
		})
	}

	// Attribute filters (nested queries)
	for _, attrFilter := range cfg.Attributes {
		nestedQuery := buildAttributeNestedQuery(attrFilter)
		mustClauses = append(mustClauses, nestedQuery)
	}

	// Construct final query
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": mustClauses,
			},
		},
		"from": cfg.Offset,
		"size": cfg.Limit,
	}

	// Add sorting
	if sort := buildSort(cfg.SortBy, cfg.SortOrder, cfg.SearchQuery != ""); sort != nil {
		query["sort"] = sort
	}

	return query
}

// buildAttributeNestedQuery constructs nested query for attribute filter
func buildAttributeNestedQuery(filter AttributeFilter) map[string]interface{} {
	// Determine value field based on attribute type
	valueField := determineValueField(filter.Type)

	// Build value query (term for single value, terms for multiple)
	var valueQuery map[string]interface{}
	if len(filter.Values) == 1 {
		// Single value - term query
		valueQuery = map[string]interface{}{
			"term": map[string]interface{}{
				valueField: filter.Values[0],
			},
		}
	} else {
		// Multiple values - terms query (OR logic)
		valueQuery = map[string]interface{}{
			"terms": map[string]interface{}{
				valueField: filter.Values,
			},
		}
	}

	// Nested query: match attribute code AND one of the values
	return map[string]interface{}{
		"nested": map[string]interface{}{
			"path": "attributes",
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": []map[string]interface{}{
						{
							"term": map[string]interface{}{
								"attributes.code": filter.Code,
							},
						},
						valueQuery,
					},
				},
			},
		},
	}
}

// determineValueField returns the correct field name based on attribute type
func determineValueField(attrType string) string {
	switch attrType {
	case "number":
		return "attributes.value_number"
	case "boolean":
		return "attributes.value_boolean"
	case "date":
		return "attributes.value_date"
	case "text":
		return "attributes.value_text.keyword"
	case "multiselect":
		return "attributes.value_multiselect"
	case "select":
		fallthrough
	default:
		return "attributes.value_select"
	}
}

// buildRangeFilter constructs range filter for numeric fields
func buildRangeFilter(field string, min, max *float64) map[string]interface{} {
	rangeParams := map[string]interface{}{}

	if min != nil {
		rangeParams["gte"] = *min
	}
	if max != nil {
		rangeParams["lte"] = *max
	}

	return map[string]interface{}{
		"range": map[string]interface{}{
			field: rangeParams,
		},
	}
}

// buildSort constructs sort configuration
func buildSort(sortBy, sortOrder string, hasTextSearch bool) interface{} {
	// Default sort order
	if sortOrder == "" {
		sortOrder = "desc"
	}

	// Default sort field
	if sortBy == "" {
		if hasTextSearch {
			sortBy = "relevance"
		} else {
			sortBy = "created_at"
		}
	}

	// Relevance sorting (by _score)
	if sortBy == "relevance" {
		return []map[string]interface{}{
			{
				"_score": map[string]interface{}{
					"order": "desc",
				},
			},
		}
	}

	// Other field sorting
	var sortField string
	switch sortBy {
	case "price":
		sortField = "price"
	case "created_at":
		sortField = "created_at"
	case "views_count":
		sortField = "views_count"
	case "favorites_count":
		sortField = "favorites_count"
	default:
		sortField = "created_at"
	}

	return []map[string]interface{}{
		{
			sortField: map[string]interface{}{
				"order": sortOrder,
			},
		},
	}
}
