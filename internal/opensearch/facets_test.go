package opensearch

import (
	"encoding/json"
	"testing"
)

func TestBuildFacetsAggregation_AttributeOnly(t *testing.T) {
	cfg := FacetsConfig{
		AttributeCodes: []string{"brand", "color"},
	}

	aggs := BuildFacetsAggregation(cfg)

	// Should have 2 attribute facets
	if len(aggs) != 2 {
		t.Fatalf("Expected 2 aggregations, got %d", len(aggs))
	}

	// Check brand facet
	if _, ok := aggs["facet_attr_brand"]; !ok {
		t.Error("Missing facet_attr_brand aggregation")
	}

	// Check color facet
	if _, ok := aggs["facet_attr_color"]; !ok {
		t.Error("Missing facet_attr_color aggregation")
	}
}

func TestBuildFacetsAggregation_AllFacets(t *testing.T) {
	cfg := FacetsConfig{
		AttributeCodes:     []string{"brand"},
		IncludeCategories:  true,
		IncludePriceRanges: true,
		IncludeSourceTypes: true,
		IncludeStockStatus: true,
	}

	aggs := BuildFacetsAggregation(cfg)

	// Should have: 1 attribute + categories + price + source + stock = 5
	if len(aggs) != 5 {
		t.Fatalf("Expected 5 aggregations, got %d", len(aggs))
	}

	// Verify all expected aggregations exist
	expectedKeys := []string{
		"facet_attr_brand",
		"facet_categories",
		"facet_price_ranges",
		"facet_source_types",
		"facet_stock_status",
	}

	for _, key := range expectedKeys {
		if _, ok := aggs[key]; !ok {
			t.Errorf("Missing expected aggregation: %s", key)
		}
	}
}

func TestBuildAttributeFacetAggregation(t *testing.T) {
	agg := buildAttributeFacetAggregation("brand")

	// Check nested structure
	nested, ok := agg["nested"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected nested aggregation")
	}

	if nested["path"] != "attributes" {
		t.Errorf("Expected path 'attributes', got %v", nested["path"])
	}

	// Check sub-aggregations
	subAggs, ok := agg["aggs"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected sub-aggregations")
	}

	filterByCode, ok := subAggs["filter_by_code"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected filter_by_code aggregation")
	}

	// Check filter term
	filter := filterByCode["filter"].(map[string]interface{})
	term := filter["term"].(map[string]interface{})
	if term["attributes.code"] != "brand" {
		t.Errorf("Expected filter on 'brand', got %v", term["attributes.code"])
	}

	// Check value aggregations exist
	valueAggs := filterByCode["aggs"].(map[string]interface{})
	expectedValueAggs := []string{
		"values_select",
		"values_multiselect",
		"values_number_stats",
		"values_boolean",
	}

	for _, aggName := range expectedValueAggs {
		if _, ok := valueAggs[aggName]; !ok {
			t.Errorf("Missing value aggregation: %s", aggName)
		}
	}
}

func TestBuildPriceRangeAggregation_Default(t *testing.T) {
	agg := buildPriceRangeAggregation(nil)

	rangeAgg := agg["range"].(map[string]interface{})
	if rangeAgg["field"] != "price" {
		t.Errorf("Expected field 'price', got %v", rangeAgg["field"])
	}

	ranges := rangeAgg["ranges"].([]map[string]interface{})
	if len(ranges) != 6 {
		t.Errorf("Expected 6 default price ranges, got %d", len(ranges))
	}

	// Check first range
	firstRange := ranges[0]
	if firstRange["from"] != float64(0) || firstRange["to"] != float64(50) {
		t.Errorf("Unexpected first range: %v - %v", firstRange["from"], firstRange["to"])
	}
}

func TestBuildPriceRangeAggregation_Custom(t *testing.T) {
	customRanges := []PriceRange{
		{From: 0, To: 100},
		{From: 100, To: 500},
		{From: 500, To: 1000},
	}

	agg := buildPriceRangeAggregation(customRanges)

	rangeAgg := agg["range"].(map[string]interface{})
	ranges := rangeAgg["ranges"].([]map[string]interface{})

	if len(ranges) != 3 {
		t.Fatalf("Expected 3 custom ranges, got %d", len(ranges))
	}

	// Verify custom ranges
	if ranges[0]["from"] != float64(0) || ranges[0]["to"] != float64(100) {
		t.Errorf("First range mismatch")
	}
	if ranges[2]["from"] != float64(500) || ranges[2]["to"] != float64(1000) {
		t.Errorf("Third range mismatch")
	}
}

func TestParseTermsBuckets(t *testing.T) {
	// Mock OpenSearch response for terms aggregation
	mockAgg := map[string]interface{}{
		"buckets": []interface{}{
			map[string]interface{}{
				"key":       "apple",
				"doc_count": float64(150),
			},
			map[string]interface{}{
				"key":       "samsung",
				"doc_count": float64(200),
			},
		},
	}

	values := parseTermsBuckets(mockAgg)

	if len(values) != 2 {
		t.Fatalf("Expected 2 values, got %d", len(values))
	}

	// Check first value
	if values[0].Value != "apple" || values[0].Count != 150 {
		t.Errorf("First value mismatch: %v", values[0])
	}

	// Check second value
	if values[1].Value != "samsung" || values[1].Count != 200 {
		t.Errorf("Second value mismatch: %v", values[1])
	}
}

func TestParseCategoryFacets(t *testing.T) {
	mockAgg := map[string]interface{}{
		"buckets": []interface{}{
			map[string]interface{}{
				"key":       float64(10),
				"doc_count": float64(50),
			},
			map[string]interface{}{
				"key":       float64(20),
				"doc_count": float64(75),
			},
		},
	}

	categories := parseCategoryFacets(mockAgg)

	if len(categories) != 2 {
		t.Fatalf("Expected 2 categories, got %d", len(categories))
	}

	if categories[0].CategoryID != 10 || categories[0].Count != 50 {
		t.Errorf("First category mismatch: %v", categories[0])
	}

	if categories[1].CategoryID != 20 || categories[1].Count != 75 {
		t.Errorf("Second category mismatch: %v", categories[1])
	}
}

func TestParsePriceRangeFacets(t *testing.T) {
	mockAgg := map[string]interface{}{
		"buckets": []interface{}{
			map[string]interface{}{
				"from":      float64(0),
				"to":        float64(50),
				"doc_count": float64(30),
			},
			map[string]interface{}{
				"from":      float64(50),
				"to":        float64(100),
				"doc_count": float64(45),
			},
			map[string]interface{}{
				"from":      float64(100),
				"to":        float64(250),
				"doc_count": float64(0), // Empty bucket - should be skipped
			},
		},
	}

	ranges := parsePriceRangeFacets(mockAgg)

	// Should have 2 ranges (empty bucket skipped)
	if len(ranges) != 2 {
		t.Fatalf("Expected 2 price ranges (empty skipped), got %d", len(ranges))
	}

	if ranges[0].Min != 0 || ranges[0].Max != 50 || ranges[0].Count != 30 {
		t.Errorf("First range mismatch: %v", ranges[0])
	}

	if ranges[1].Min != 50 || ranges[1].Max != 100 || ranges[1].Count != 45 {
		t.Errorf("Second range mismatch: %v", ranges[1])
	}
}

func TestParseSimpleFacets(t *testing.T) {
	mockAgg := map[string]interface{}{
		"buckets": []interface{}{
			map[string]interface{}{
				"key":       "c2c",
				"doc_count": float64(120),
			},
			map[string]interface{}{
				"key":       "b2c",
				"doc_count": float64(80),
			},
		},
	}

	facets := parseSimpleFacets(mockAgg)

	if len(facets) != 2 {
		t.Fatalf("Expected 2 facets, got %d", len(facets))
	}

	if facets[0].Key != "c2c" || facets[0].Count != 120 {
		t.Errorf("First facet mismatch: %v", facets[0])
	}

	if facets[1].Key != "b2c" || facets[1].Count != 80 {
		t.Errorf("Second facet mismatch: %v", facets[1])
	}
}

func TestParseAttributeFacet(t *testing.T) {
	// Mock nested aggregation response for brand attribute
	mockAgg := map[string]interface{}{
		"filter_by_code": map[string]interface{}{
			"values_select": map[string]interface{}{
				"buckets": []interface{}{
					map[string]interface{}{
						"key":       "apple",
						"doc_count": float64(100),
					},
					map[string]interface{}{
						"key":       "samsung",
						"doc_count": float64(150),
					},
				},
			},
			"values_multiselect": map[string]interface{}{
				"buckets": []interface{}{},
			},
			"values_boolean": map[string]interface{}{
				"buckets": []interface{}{},
			},
			"values_number_stats": map[string]interface{}{
				"count": float64(0),
			},
		},
	}

	metadata := map[string]AttributeMetadata{
		"brand": {Code: "brand", Name: "Brand", Type: "select"},
	}

	facet := parseAttributeFacet("brand", mockAgg, metadata)

	if facet == nil {
		t.Fatal("Expected facet, got nil")
	}

	if facet.Code != "brand" {
		t.Errorf("Expected code 'brand', got %s", facet.Code)
	}

	if facet.Name != "Brand" {
		t.Errorf("Expected name 'Brand', got %s", facet.Name)
	}

	if len(facet.Values) != 2 {
		t.Fatalf("Expected 2 values, got %d", len(facet.Values))
	}

	if facet.Values[0].Value != "apple" || facet.Values[0].Count != 100 {
		t.Errorf("First value mismatch: %v", facet.Values[0])
	}
}

func TestParseFacetsResponse_Complete(t *testing.T) {
	// Complete mock aggregations response
	mockAggs := map[string]interface{}{
		"facet_attr_brand": map[string]interface{}{
			"filter_by_code": map[string]interface{}{
				"values_select": map[string]interface{}{
					"buckets": []interface{}{
						map[string]interface{}{
							"key":       "apple",
							"doc_count": float64(50),
						},
					},
				},
				"values_multiselect": map[string]interface{}{"buckets": []interface{}{}},
				"values_boolean":     map[string]interface{}{"buckets": []interface{}{}},
				"values_number_stats": map[string]interface{}{"count": float64(0)},
			},
		},
		"facet_categories": map[string]interface{}{
			"buckets": []interface{}{
				map[string]interface{}{
					"key":       float64(10),
					"doc_count": float64(25),
				},
			},
		},
		"facet_price_ranges": map[string]interface{}{
			"buckets": []interface{}{
				map[string]interface{}{
					"from":      float64(0),
					"to":        float64(100),
					"doc_count": float64(15),
				},
			},
		},
		"facet_source_types": map[string]interface{}{
			"buckets": []interface{}{
				map[string]interface{}{
					"key":       "c2c",
					"doc_count": float64(30),
				},
			},
		},
		"facet_stock_status": map[string]interface{}{
			"buckets": []interface{}{
				map[string]interface{}{
					"key":       "in_stock",
					"doc_count": float64(40),
				},
			},
		},
	}

	metadata := map[string]AttributeMetadata{
		"brand": {Code: "brand", Name: "Brand", Type: "select"},
	}

	result := ParseFacetsResponse(mockAggs, metadata)

	// Verify all facet types were parsed
	if len(result.Attributes) != 1 {
		t.Errorf("Expected 1 attribute facet, got %d", len(result.Attributes))
	}

	if len(result.Categories) != 1 {
		t.Errorf("Expected 1 category facet, got %d", len(result.Categories))
	}

	if len(result.PriceRanges) != 1 {
		t.Errorf("Expected 1 price range, got %d", len(result.PriceRanges))
	}

	if len(result.SourceTypes) != 1 {
		t.Errorf("Expected 1 source type facet, got %d", len(result.SourceTypes))
	}

	if len(result.StockStatuses) != 1 {
		t.Errorf("Expected 1 stock status facet, got %d", len(result.StockStatuses))
	}
}

func TestFacetsAggregation_JSONSerialization(t *testing.T) {
	// Test that aggregation can be serialized to JSON
	cfg := FacetsConfig{
		AttributeCodes:     []string{"brand", "color"},
		IncludeCategories:  true,
		IncludePriceRanges: true,
	}

	aggs := BuildFacetsAggregation(cfg)

	jsonBytes, err := json.Marshal(aggs)
	if err != nil {
		t.Fatalf("Failed to serialize aggregations to JSON: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("Failed to deserialize JSON: %v", err)
	}

	if len(result) != 4 {
		t.Errorf("Expected 4 aggregations after JSON round-trip, got %d", len(result))
	}
}
