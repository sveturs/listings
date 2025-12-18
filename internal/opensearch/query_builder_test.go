package opensearch

import (
	"encoding/json"
	"testing"
)

func TestBuildFilterQuery_CategoryOnly(t *testing.T) {
	cfg := FilterQuery{
		CategoryID: 123,
		Limit:      20,
		Offset:     0,
	}

	query := BuildFilterQuery(cfg)

	// Check structure
	if query["from"] != 0 {
		t.Errorf("Expected offset 0, got %v", query["from"])
	}
	if query["size"] != 20 {
		t.Errorf("Expected limit 20, got %v", query["size"])
	}

	// Check category filter
	queryPart := query["query"].(map[string]interface{})
	boolPart := queryPart["bool"].(map[string]interface{})
	mustClauses := boolPart["must"].([]map[string]interface{})

	if len(mustClauses) != 1 {
		t.Fatalf("Expected 1 must clause, got %d", len(mustClauses))
	}

	termQuery := mustClauses[0]["term"].(map[string]interface{})
	if termQuery["category_id"] != int64(123) {
		t.Errorf("Expected category_id 123, got %v", termQuery["category_id"])
	}
}

func TestBuildFilterQuery_WithPriceRange(t *testing.T) {
	minPrice := 100.0
	maxPrice := 500.0

	cfg := FilterQuery{
		CategoryID: 123,
		PriceMin:   &minPrice,
		PriceMax:   &maxPrice,
		Limit:      20,
		Offset:     0,
	}

	query := BuildFilterQuery(cfg)

	queryPart := query["query"].(map[string]interface{})
	boolPart := queryPart["bool"].(map[string]interface{})
	mustClauses := boolPart["must"].([]map[string]interface{})

	// Should have 2 clauses: category + price range
	if len(mustClauses) != 2 {
		t.Fatalf("Expected 2 must clauses, got %d", len(mustClauses))
	}

	// Find price range clause
	var priceClause map[string]interface{}
	for _, clause := range mustClauses {
		if _, ok := clause["range"]; ok {
			priceClause = clause
			break
		}
	}

	if priceClause == nil {
		t.Fatal("Price range clause not found")
	}

	rangeQuery := priceClause["range"].(map[string]interface{})
	priceRange := rangeQuery["price"].(map[string]interface{})

	if priceRange["gte"] != 100.0 {
		t.Errorf("Expected min price 100.0, got %v", priceRange["gte"])
	}
	if priceRange["lte"] != 500.0 {
		t.Errorf("Expected max price 500.0, got %v", priceRange["lte"])
	}
}

func TestBuildFilterQuery_WithAttributeFilters(t *testing.T) {
	cfg := FilterQuery{
		CategoryID: 123,
		Attributes: []AttributeFilter{
			{Code: "brand", Type: "select", Values: []string{"apple"}},
			{Code: "color", Type: "select", Values: []string{"black", "white"}},
		},
		Limit:  20,
		Offset: 0,
	}

	query := BuildFilterQuery(cfg)

	queryPart := query["query"].(map[string]interface{})
	boolPart := queryPart["bool"].(map[string]interface{})
	mustClauses := boolPart["must"].([]map[string]interface{})

	// Should have 3 clauses: category + 2 attribute filters
	if len(mustClauses) != 3 {
		t.Fatalf("Expected 3 must clauses, got %d", len(mustClauses))
	}

	// Count nested queries (attribute filters)
	nestedCount := 0
	for _, clause := range mustClauses {
		if _, ok := clause["nested"]; ok {
			nestedCount++
		}
	}

	if nestedCount != 2 {
		t.Errorf("Expected 2 nested queries, got %d", nestedCount)
	}
}

func TestBuildFilterQuery_WithTextSearch(t *testing.T) {
	cfg := FilterQuery{
		CategoryID:  123,
		SearchQuery: "iphone 15",
		Limit:       20,
		Offset:      0,
	}

	query := BuildFilterQuery(cfg)

	queryPart := query["query"].(map[string]interface{})
	boolPart := queryPart["bool"].(map[string]interface{})
	mustClauses := boolPart["must"].([]map[string]interface{})

	// Should have 2 clauses: category + text search
	if len(mustClauses) != 2 {
		t.Fatalf("Expected 2 must clauses, got %d", len(mustClauses))
	}

	// Find multi_match clause
	var multiMatchClause map[string]interface{}
	for _, clause := range mustClauses {
		if _, ok := clause["multi_match"]; ok {
			multiMatchClause = clause
			break
		}
	}

	if multiMatchClause == nil {
		t.Fatal("multi_match clause not found")
	}

	multiMatch := multiMatchClause["multi_match"].(map[string]interface{})
	if multiMatch["query"] != "iphone 15" {
		t.Errorf("Expected query 'iphone 15', got %v", multiMatch["query"])
	}
}

func TestBuildFilterQuery_WithGeoDistance(t *testing.T) {
	lat := 44.7866
	lon := 20.4489
	radius := 10.0

	cfg := FilterQuery{
		CategoryID: 123,
		LatLon:     &GeoLocation{Lat: lat, Lon: lon},
		RadiusKM:   &radius,
		Limit:      20,
		Offset:     0,
	}

	query := BuildFilterQuery(cfg)

	queryPart := query["query"].(map[string]interface{})
	boolPart := queryPart["bool"].(map[string]interface{})
	mustClauses := boolPart["must"].([]map[string]interface{})

	// Find geo_distance clause
	var geoClause map[string]interface{}
	for _, clause := range mustClauses {
		if _, ok := clause["geo_distance"]; ok {
			geoClause = clause
			break
		}
	}

	if geoClause == nil {
		t.Fatal("geo_distance clause not found")
	}

	geoDist := geoClause["geo_distance"].(map[string]interface{})
	if geoDist["distance"] != "10.00km" {
		t.Errorf("Expected distance '10.00km', got %v", geoDist["distance"])
	}

	location := geoDist["location"].(map[string]interface{})
	if location["lat"] != lat {
		t.Errorf("Expected lat %f, got %v", lat, location["lat"])
	}
	if location["lon"] != lon {
		t.Errorf("Expected lon %f, got %v", lon, location["lon"])
	}
}

func TestBuildFilterQuery_WithSourceTypeAndStockStatus(t *testing.T) {
	cfg := FilterQuery{
		CategoryID:  123,
		SourceType:  "b2c",
		StockStatus: "in_stock",
		Limit:       20,
		Offset:      0,
	}

	query := BuildFilterQuery(cfg)

	queryPart := query["query"].(map[string]interface{})
	boolPart := queryPart["bool"].(map[string]interface{})
	mustClauses := boolPart["must"].([]map[string]interface{})

	// Should have 3 clauses: category + source_type + stock_status
	if len(mustClauses) != 3 {
		t.Fatalf("Expected 3 must clauses, got %d", len(mustClauses))
	}
}

func TestBuildAttributeNestedQuery_SingleValue(t *testing.T) {
	filter := AttributeFilter{
		Code:   "brand",
		Type:   "select",
		Values: []string{"apple"},
	}

	nestedQuery := buildAttributeNestedQuery(filter)

	// Check nested structure
	nested := nestedQuery["nested"].(map[string]interface{})
	if nested["path"] != "attributes" {
		t.Errorf("Expected path 'attributes', got %v", nested["path"])
	}

	// Check query structure
	query := nested["query"].(map[string]interface{})
	boolPart := query["bool"].(map[string]interface{})
	mustClauses := boolPart["must"].([]map[string]interface{})

	if len(mustClauses) != 2 {
		t.Fatalf("Expected 2 must clauses in nested query, got %d", len(mustClauses))
	}

	// Check code match
	codeClause := mustClauses[0]["term"].(map[string]interface{})
	if codeClause["attributes.code"] != "brand" {
		t.Errorf("Expected code 'brand', got %v", codeClause["attributes.code"])
	}

	// Check value match (term for single value)
	valueClause := mustClauses[1]["term"].(map[string]interface{})
	if valueClause["attributes.value_select"] != "apple" {
		t.Errorf("Expected value 'apple', got %v", valueClause["attributes.value_select"])
	}
}

func TestBuildAttributeNestedQuery_MultipleValues(t *testing.T) {
	filter := AttributeFilter{
		Code:   "color",
		Type:   "select",
		Values: []string{"black", "white"},
	}

	nestedQuery := buildAttributeNestedQuery(filter)

	nested := nestedQuery["nested"].(map[string]interface{})
	query := nested["query"].(map[string]interface{})
	boolPart := query["bool"].(map[string]interface{})
	mustClauses := boolPart["must"].([]map[string]interface{})

	// Check value match (terms for multiple values)
	valueClause := mustClauses[1]["terms"].(map[string]interface{})
	values := valueClause["attributes.value_select"].([]string)

	if len(values) != 2 {
		t.Fatalf("Expected 2 values, got %d", len(values))
	}
	if values[0] != "black" || values[1] != "white" {
		t.Errorf("Expected values [black, white], got %v", values)
	}
}

func TestDetermineValueField(t *testing.T) {
	tests := []struct {
		attrType  string
		expected  string
	}{
		{"select", "attributes.value_select"},
		{"multiselect", "attributes.value_multiselect"},
		{"number", "attributes.value_number"},
		{"boolean", "attributes.value_boolean"},
		{"date", "attributes.value_date"},
		{"text", "attributes.value_text.keyword"},
		{"unknown", "attributes.value_select"}, // default
	}

	for _, tt := range tests {
		t.Run(tt.attrType, func(t *testing.T) {
			result := determineValueField(tt.attrType)
			if result != tt.expected {
				t.Errorf("For type %s, expected %s, got %s", tt.attrType, tt.expected, result)
			}
		})
	}
}

func TestBuildSort_DefaultSorting(t *testing.T) {
	// Without text search - should sort by created_at desc
	sort := buildSort("", "", false)
	sortArr := sort.([]map[string]interface{})

	if len(sortArr) != 1 {
		t.Fatalf("Expected 1 sort clause, got %d", len(sortArr))
	}

	if _, ok := sortArr[0]["created_at"]; !ok {
		t.Error("Expected sort by created_at")
	}
}

func TestBuildSort_RelevanceSorting(t *testing.T) {
	// With text search - should sort by relevance (score)
	sort := buildSort("", "", true)
	sortArr := sort.([]map[string]interface{})

	if len(sortArr) != 1 {
		t.Fatalf("Expected 1 sort clause, got %d", len(sortArr))
	}

	if _, ok := sortArr[0]["_score"]; !ok {
		t.Error("Expected sort by _score")
	}
}

func TestBuildSort_PriceSorting(t *testing.T) {
	sort := buildSort("price", "asc", false)
	sortArr := sort.([]map[string]interface{})

	priceSort := sortArr[0]["price"].(map[string]interface{})
	if priceSort["order"] != "asc" {
		t.Errorf("Expected order 'asc', got %v", priceSort["order"])
	}
}

func TestBuildFilterQuery_JSONSerialization(t *testing.T) {
	// Test that generated query can be serialized to JSON (important for OpenSearch)
	minPrice := 100.0
	cfg := FilterQuery{
		CategoryID: 123,
		PriceMin:   &minPrice,
		SearchQuery: "test",
		Attributes: []AttributeFilter{
			{Code: "brand", Type: "select", Values: []string{"apple"}},
		},
		Limit:  20,
		Offset: 0,
	}

	query := BuildFilterQuery(cfg)

	// Try to serialize to JSON
	jsonBytes, err := json.Marshal(query)
	if err != nil {
		t.Fatalf("Failed to serialize query to JSON: %v", err)
	}

	// Try to deserialize back
	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("Failed to deserialize JSON: %v", err)
	}

	// Basic validation
	if result["from"] != float64(0) {
		t.Errorf("JSON round-trip failed for 'from' field")
	}
	if result["size"] != float64(20) {
		t.Errorf("JSON round-trip failed for 'size' field")
	}
}
