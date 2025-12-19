package opensearch

import (
	"fmt"
	"strings"
)

// Facet represents a facet with its possible values and counts
type Facet struct {
	Code   string        `json:"code"`
	Name   string        `json:"name"`
	Type   string        `json:"type"`
	Values []FacetValue  `json:"values"`
}

// FacetValue represents a single value option in a facet
type FacetValue struct {
	Value string `json:"value"`
	Label string `json:"label"`
	Count int64  `json:"count"`
}

// CategoryFacet represents category distribution
type CategoryFacet struct {
	CategoryID int64 `json:"category_id"`
	Count      int64 `json:"count"`
}

// PriceRangeFacet represents price range bucket
type PriceRangeFacet struct {
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
	Count int64   `json:"count"`
}

// FacetsConfig contains configuration for facet aggregations
type FacetsConfig struct {
	// Attribute codes to generate facets for
	AttributeCodes []string

	// Include category distribution
	IncludeCategories bool

	// Include price ranges
	IncludePriceRanges bool
	PriceRanges        []PriceRange // Custom price ranges, if empty uses default

	// Include source type distribution (c2c, b2c)
	IncludeSourceTypes bool

	// Include stock status distribution
	IncludeStockStatus bool
}

// PriceRange represents a price range for aggregation
type PriceRange struct {
	From float64
	To   float64
}

// BuildFacetsAggregation constructs OpenSearch aggregations for facets
func BuildFacetsAggregation(cfg FacetsConfig) map[string]interface{} {
	aggs := map[string]interface{}{}

	// Attribute facets (nested aggregations)
	for _, code := range cfg.AttributeCodes {
		aggs[fmt.Sprintf("facet_attr_%s", code)] = buildAttributeFacetAggregation(code)
	}

	// Category distribution
	if cfg.IncludeCategories {
		aggs["facet_categories"] = map[string]interface{}{
			"terms": map[string]interface{}{
				"field": "category_id",
				"size":  100, // Max categories to return
			},
		}
	}

	// Price ranges
	if cfg.IncludePriceRanges {
		aggs["facet_price_ranges"] = buildPriceRangeAggregation(cfg.PriceRanges)
	}

	// Source types (c2c, b2c)
	if cfg.IncludeSourceTypes {
		aggs["facet_source_types"] = map[string]interface{}{
			"terms": map[string]interface{}{
				"field": "source_type",
				"size":  10,
			},
		}
	}

	// Stock status
	if cfg.IncludeStockStatus {
		aggs["facet_stock_status"] = map[string]interface{}{
			"terms": map[string]interface{}{
				"field": "stock_status",
				"size":  10,
			},
		}
	}

	return aggs
}

// buildAttributeFacetAggregation creates nested aggregation for a single attribute
func buildAttributeFacetAggregation(code string) map[string]interface{} {
	return map[string]interface{}{
		"nested": map[string]interface{}{
			"path": "attributes",
		},
		"aggs": map[string]interface{}{
			"filter_by_code": map[string]interface{}{
				"filter": map[string]interface{}{
					"term": map[string]interface{}{
						"attributes.code": code,
					},
				},
				"aggs": map[string]interface{}{
					// Aggregate select values
					"values_select": map[string]interface{}{
						"terms": map[string]interface{}{
							"field": "attributes.value_select",
							"size":  100,
						},
					},
					// Aggregate multiselect values
					"values_multiselect": map[string]interface{}{
						"terms": map[string]interface{}{
							"field": "attributes.value_multiselect",
							"size":  100,
						},
					},
					// Stats for numeric values
					"values_number_stats": map[string]interface{}{
						"stats": map[string]interface{}{
							"field": "attributes.value_number",
						},
					},
					// Boolean value counts
					"values_boolean": map[string]interface{}{
						"terms": map[string]interface{}{
							"field": "attributes.value_boolean",
							"size":  2,
						},
					},
				},
			},
		},
	}
}

// buildPriceRangeAggregation creates price range aggregation
func buildPriceRangeAggregation(customRanges []PriceRange) map[string]interface{} {
	ranges := customRanges

	// Default ranges if none provided
	if len(ranges) == 0 {
		ranges = []PriceRange{
			{From: 0, To: 50},
			{From: 50, To: 100},
			{From: 100, To: 250},
			{From: 250, To: 500},
			{From: 500, To: 1000},
			{From: 1000, To: 999999},
		}
	}

	rangeAggs := []map[string]interface{}{}
	for _, r := range ranges {
		rangeAggs = append(rangeAggs, map[string]interface{}{
			"from": r.From,
			"to":   r.To,
		})
	}

	return map[string]interface{}{
		"range": map[string]interface{}{
			"field":  "price",
			"ranges": rangeAggs,
		},
	}
}

// ParseFacetsResponse parses OpenSearch aggregation response into facets
func ParseFacetsResponse(aggs map[string]interface{}, attributeMetadata map[string]AttributeMetadata) FacetsResult {
	result := FacetsResult{
		Attributes:   []Facet{},
		Categories:   []CategoryFacet{},
		PriceRanges:  []PriceRangeFacet{},
		SourceTypes:  []SimpleFacet{},
		StockStatuses: []SimpleFacet{},
	}

	// Parse attribute facets
	for key, value := range aggs {
		if strings.HasPrefix(key, "facet_attr_") {
			code := strings.TrimPrefix(key, "facet_attr_")
			facet := parseAttributeFacet(code, value, attributeMetadata)
			if facet != nil {
				result.Attributes = append(result.Attributes, *facet)
			}
		}
	}

	// Parse category distribution
	if catAgg, ok := aggs["facet_categories"]; ok {
		result.Categories = parseCategoryFacets(catAgg)
	}

	// Parse price ranges
	if priceAgg, ok := aggs["facet_price_ranges"]; ok {
		result.PriceRanges = parsePriceRangeFacets(priceAgg)
	}

	// Parse source types
	if sourceAgg, ok := aggs["facet_source_types"]; ok {
		result.SourceTypes = parseSimpleFacets(sourceAgg)
	}

	// Parse stock status
	if stockAgg, ok := aggs["facet_stock_status"]; ok {
		result.StockStatuses = parseSimpleFacets(stockAgg)
	}

	return result
}

// AttributeMetadata contains metadata about an attribute (name, type, etc.)
type AttributeMetadata struct {
	Code string
	Name string
	Type string
}

// FacetsResult contains all parsed facets
type FacetsResult struct {
	Attributes    []Facet           `json:"attributes"`
	Categories    []CategoryFacet   `json:"categories"`
	PriceRanges   []PriceRangeFacet `json:"price_ranges"`
	SourceTypes   []SimpleFacet     `json:"source_types"`
	StockStatuses []SimpleFacet     `json:"stock_statuses"`
}

// SimpleFacet represents a simple key-count facet
type SimpleFacet struct {
	Key   string `json:"key"`
	Count int64  `json:"count"`
}

// parseAttributeFacet parses nested attribute aggregation
func parseAttributeFacet(code string, aggData interface{}, metadata map[string]AttributeMetadata) *Facet {
	aggMap, ok := aggData.(map[string]interface{})
	if !ok {
		return nil
	}

	// Navigate nested structure: nested -> filter_by_code -> values
	filterByCode, ok := aggMap["filter_by_code"].(map[string]interface{})
	if !ok {
		return nil
	}

	facet := &Facet{
		Code:   code,
		Values: []FacetValue{},
	}

	// Get metadata if available
	if meta, exists := metadata[code]; exists {
		facet.Name = meta.Name
		facet.Type = meta.Type
	}

	// Parse select values
	if selectAgg, ok := filterByCode["values_select"].(map[string]interface{}); ok {
		facet.Values = append(facet.Values, parseTermsBuckets(selectAgg)...)
	}

	// Parse multiselect values
	if multiselectAgg, ok := filterByCode["values_multiselect"].(map[string]interface{}); ok {
		facet.Values = append(facet.Values, parseTermsBuckets(multiselectAgg)...)
	}

	// Parse boolean values
	if boolAgg, ok := filterByCode["values_boolean"].(map[string]interface{}); ok {
		facet.Values = append(facet.Values, parseTermsBuckets(boolAgg)...)
	}

	// Parse number stats (for range facets)
	if statsAgg, ok := filterByCode["values_number_stats"].(map[string]interface{}); ok {
		if count, ok := statsAgg["count"].(float64); ok && count > 0 {
			min := statsAgg["min"].(float64)
			max := statsAgg["max"].(float64)
			facet.Values = append(facet.Values, FacetValue{
				Value: fmt.Sprintf("%.0f-%.0f", min, max),
				Label: fmt.Sprintf("%.0f - %.0f", min, max),
				Count: int64(count),
			})
		}
	}

	return facet
}

// parseTermsBuckets parses terms aggregation buckets
func parseTermsBuckets(aggData interface{}) []FacetValue {
	aggMap, ok := aggData.(map[string]interface{})
	if !ok {
		return nil
	}

	buckets, ok := aggMap["buckets"].([]interface{})
	if !ok {
		return nil
	}

	values := []FacetValue{}
	for _, bucket := range buckets {
		b, ok := bucket.(map[string]interface{})
		if !ok {
			continue
		}

		key := fmt.Sprintf("%v", b["key"])
		count := int64(b["doc_count"].(float64))

		values = append(values, FacetValue{
			Value: key,
			Label: key, // Can be enhanced with i18n later
			Count: count,
		})
	}

	return values
}

// parseCategoryFacets parses category distribution
func parseCategoryFacets(aggData interface{}) []CategoryFacet {
	aggMap, ok := aggData.(map[string]interface{})
	if !ok {
		return nil
	}

	buckets, ok := aggMap["buckets"].([]interface{})
	if !ok {
		return nil
	}

	categories := []CategoryFacet{}
	for _, bucket := range buckets {
		b, ok := bucket.(map[string]interface{})
		if !ok {
			continue
		}

		categoryID := int64(b["key"].(float64))
		count := int64(b["doc_count"].(float64))

		categories = append(categories, CategoryFacet{
			CategoryID: categoryID,
			Count:      count,
		})
	}

	return categories
}

// parsePriceRangeFacets parses price range buckets
func parsePriceRangeFacets(aggData interface{}) []PriceRangeFacet {
	aggMap, ok := aggData.(map[string]interface{})
	if !ok {
		return nil
	}

	buckets, ok := aggMap["buckets"].([]interface{})
	if !ok {
		return nil
	}

	ranges := []PriceRangeFacet{}
	for _, bucket := range buckets {
		b, ok := bucket.(map[string]interface{})
		if !ok {
			continue
		}

		min := 0.0
		max := 0.0

		if from, ok := b["from"].(float64); ok {
			min = from
		}
		if to, ok := b["to"].(float64); ok {
			max = to
		}

		count := int64(b["doc_count"].(float64))

		// Skip empty buckets
		if count > 0 {
			ranges = append(ranges, PriceRangeFacet{
				Min:   min,
				Max:   max,
				Count: count,
			})
		}
	}

	return ranges
}

// parseSimpleFacets parses simple terms aggregation (source_type, stock_status)
func parseSimpleFacets(aggData interface{}) []SimpleFacet {
	aggMap, ok := aggData.(map[string]interface{})
	if !ok {
		return nil
	}

	buckets, ok := aggMap["buckets"].([]interface{})
	if !ok {
		return nil
	}

	facets := []SimpleFacet{}
	for _, bucket := range buckets {
		b, ok := bucket.(map[string]interface{})
		if !ok {
			continue
		}

		key := fmt.Sprintf("%v", b["key"])
		count := int64(b["doc_count"].(float64))

		facets = append(facets, SimpleFacet{
			Key:   key,
			Count: count,
		})
	}

	return facets
}
