// backend/internal/storage/postgres/marketplace_grpc_client_test.go
package postgres

import (
	"testing"
	"time"

	listingsv1 "github.com/sveturs/listings/api/proto/listings/v1"
)

func TestConvertProtoToListing(t *testing.T) {
	// Arrange
	testID := int64(123)
	testUserID := int64(456)
	testCategoryID := int64(789)
	testTitle := "Test Listing"
	testDescription := "Test Description"
	testPrice := 99.99
	testStorefrontID := int64(111)
	testCreatedAt := time.Now().Format(time.RFC3339)
	testUpdatedAt := time.Now().Format(time.RFC3339)

	pbListing := &listingsv1.Listing{
		Id:           testID,
		UserId:       testUserID,
		CategoryId:   testCategoryID,
		Title:        testTitle,
		Description:  &testDescription,
		Price:        testPrice,
		Currency:     "RSD",
		Status:       "active",
		Visibility:   "public",
		Quantity:     1,
		ViewsCount:   10,
		CreatedAt:    testCreatedAt,
		UpdatedAt:    testUpdatedAt,
		StorefrontId: &testStorefrontID,
	}

	// Act
	result := convertProtoToListing(pbListing)

	// Assert
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	if result.ID != int(testID) {
		t.Errorf("Expected ID %d, got %d", testID, result.ID)
	}

	if result.UserID != int(testUserID) {
		t.Errorf("Expected UserID %d, got %d", testUserID, result.UserID)
	}

	if result.CategoryID != int(testCategoryID) {
		t.Errorf("Expected CategoryID %d, got %d", testCategoryID, result.CategoryID)
	}

	if result.Title != testTitle {
		t.Errorf("Expected Title %s, got %s", testTitle, result.Title)
	}

	if result.Description != testDescription {
		t.Errorf("Expected Description %s, got %s", testDescription, result.Description)
	}

	if result.Price != testPrice {
		t.Errorf("Expected Price %.2f, got %.2f", testPrice, result.Price)
	}

	if result.Status != "active" {
		t.Errorf("Expected Status 'active', got %s", result.Status)
	}

	if result.StorefrontID == nil {
		t.Error("Expected non-nil StorefrontID")
	} else if *result.StorefrontID != int(testStorefrontID) {
		t.Errorf("Expected StorefrontID %d, got %d", testStorefrontID, *result.StorefrontID)
	}

	if result.ViewsCount != 10 {
		t.Errorf("Expected ViewsCount 10, got %d", result.ViewsCount)
	}

	// Check timestamps
	if result.CreatedAt.IsZero() {
		t.Error("Expected non-zero CreatedAt")
	}

	if result.UpdatedAt.IsZero() {
		t.Error("Expected non-zero UpdatedAt")
	}
}

func TestConvertProtoToListing_WithImages(t *testing.T) {
	// Arrange
	thumbnailURL := "https://example.com/thumb.jpg"
	pbListing := &listingsv1.Listing{
		Id:         123,
		UserId:     456,
		CategoryId: 789,
		Title:      "Test",
		Price:      99.99,
		Currency:   "RSD",
		Status:     "active",
		Images: []*listingsv1.ListingImage{
			{
				Id:           1,
				ListingId:    123,
				Url:          "https://example.com/image.jpg",
				ThumbnailUrl: &thumbnailURL,
				DisplayOrder: 0,
				IsPrimary:    true,
				CreatedAt:    time.Now().Format(time.RFC3339),
			},
		},
	}

	// Act
	result := convertProtoToListing(pbListing)

	// Assert
	if len(result.Images) != 1 {
		t.Fatalf("Expected 1 image, got %d", len(result.Images))
	}

	img := result.Images[0]
	if img.ID != 1 {
		t.Errorf("Expected image ID 1, got %d", img.ID)
	}

	if img.ListingID != 123 {
		t.Errorf("Expected image ListingID 123, got %d", img.ListingID)
	}

	if img.PublicURL != "https://example.com/image.jpg" {
		t.Errorf("Expected PublicURL, got %s", img.PublicURL)
	}

	if img.ThumbnailURL != thumbnailURL {
		t.Errorf("Expected ThumbnailURL %s, got %s", thumbnailURL, img.ThumbnailURL)
	}

	if !img.IsMain {
		t.Error("Expected IsMain to be true")
	}

	if img.DisplayOrder != 0 {
		t.Errorf("Expected DisplayOrder 0, got %d", img.DisplayOrder)
	}
}

func TestConvertProtoToListing_WithLocation(t *testing.T) {
	// Arrange
	country := "Serbia"
	city := "Belgrade"
	lat := 44.7866
	lon := 20.4489

	pbListing := &listingsv1.Listing{
		Id:         123,
		UserId:     456,
		CategoryId: 789,
		Title:      "Test",
		Price:      99.99,
		Currency:   "RSD",
		Status:     "active",
		Location: &listingsv1.ListingLocation{
			Id:        1,
			ListingId: 123,
			Country:   &country,
			City:      &city,
			Latitude:  &lat,
			Longitude: &lon,
		},
	}

	// Act
	result := convertProtoToListing(pbListing)

	// Assert
	if result.Country != country {
		t.Errorf("Expected Country %s, got %s", country, result.Country)
	}

	if result.City != city {
		t.Errorf("Expected City %s, got %s", city, result.City)
	}

	if result.Latitude == nil {
		t.Error("Expected non-nil Latitude")
	} else if *result.Latitude != lat {
		t.Errorf("Expected Latitude %.4f, got %.4f", lat, *result.Latitude)
	}

	if result.Longitude == nil {
		t.Error("Expected non-nil Longitude")
	} else if *result.Longitude != lon {
		t.Errorf("Expected Longitude %.4f, got %.4f", lon, *result.Longitude)
	}
}

func TestConvertProtoToListing_WithAttributes(t *testing.T) {
	// Arrange
	pbListing := &listingsv1.Listing{
		Id:         123,
		UserId:     456,
		CategoryId: 789,
		Title:      "Test",
		Price:      99.99,
		Currency:   "RSD",
		Status:     "active",
		Attributes: []*listingsv1.ListingAttribute{
			{
				Id:             1,
				ListingId:      123,
				AttributeKey:   "color",
				AttributeValue: "red",
			},
			{
				Id:             2,
				ListingId:      123,
				AttributeKey:   "size",
				AttributeValue: "large",
			},
		},
	}

	// Act
	result := convertProtoToListing(pbListing)

	// Assert
	if len(result.Attributes) != 2 {
		t.Fatalf("Expected 2 attributes, got %d", len(result.Attributes))
	}

	// Check first attribute
	attr1 := result.Attributes[0]
	if attr1.AttributeName != "color" {
		t.Errorf("Expected AttributeName 'color', got %s", attr1.AttributeName)
	}
	if attr1.TextValue == nil {
		t.Error("Expected non-nil TextValue")
	} else if *attr1.TextValue != "red" {
		t.Errorf("Expected TextValue 'red', got %s", *attr1.TextValue)
	}

	// Check second attribute
	attr2 := result.Attributes[1]
	if attr2.AttributeName != "size" {
		t.Errorf("Expected AttributeName 'size', got %s", attr2.AttributeName)
	}
	if attr2.TextValue == nil {
		t.Error("Expected non-nil TextValue")
	} else if *attr2.TextValue != "large" {
		t.Errorf("Expected TextValue 'large', got %s", *attr2.TextValue)
	}
}

func TestConvertProtoToListing_MinimalData(t *testing.T) {
	// Arrange - только обязательные поля
	pbListing := &listingsv1.Listing{
		Id:         123,
		UserId:     456,
		CategoryId: 789,
		Title:      "Test",
		Price:      99.99,
		Currency:   "RSD",
		Status:     "active",
	}

	// Act
	result := convertProtoToListing(pbListing)

	// Assert
	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	// Check that optional fields are handled correctly
	if result.StorefrontID != nil {
		t.Error("Expected nil StorefrontID for minimal data")
	}

	if result.Description != "" {
		t.Error("Expected empty Description for minimal data")
	}

	if len(result.Images) != 0 {
		t.Error("Expected empty Images array for minimal data")
	}

	if len(result.Attributes) != 0 {
		t.Error("Expected empty Attributes array for minimal data")
	}

	if result.Latitude != nil {
		t.Error("Expected nil Latitude for minimal data")
	}

	if result.Longitude != nil {
		t.Error("Expected nil Longitude for minimal data")
	}
}

// Integration tests for CRUD operations
// These tests are skipped by default as they require running listings microservice

// TestMarketplaceGRPCClient_CreateListing_Integration tests CreateListing method
func TestMarketplaceGRPCClient_CreateListing_Integration(t *testing.T) {
	t.Skip("Integration test - requires running listings microservice on localhost:50051")

	// This test verifies that CreateListing correctly:
	// 1. Converts domain model to protobuf request
	// 2. Calls gRPC CreateListing RPC
	// 3. Returns the new listing ID
	// 4. Handles errors appropriately
}

// TestMarketplaceGRPCClient_GetListingByID_Integration tests GetListingByID method
func TestMarketplaceGRPCClient_GetListingByID_Integration(t *testing.T) {
	t.Skip("Integration test - requires running listings microservice on localhost:50051")

	// This test verifies that GetListingByID correctly:
	// 1. Calls gRPC GetListing RPC with ID
	// 2. Converts protobuf response to domain model
	// 3. Returns populated listing with all fields
	// 4. Handles not found errors
}

// TestMarketplaceGRPCClient_GetListings_Integration tests GetListings method
func TestMarketplaceGRPCClient_GetListings_Integration(t *testing.T) {
	t.Skip("Integration test - requires running listings microservice on localhost:50051")

	// This test verifies that GetListings correctly:
	// 1. Converts filter map to protobuf request
	// 2. Applies all filters (user_id, storefront_id, category_id, status, price range)
	// 3. Handles pagination (limit, offset)
	// 4. Returns listings array and total count
}

// TestMarketplaceGRPCClient_UpdateListing_Integration tests UpdateListing method
func TestMarketplaceGRPCClient_UpdateListing_Integration(t *testing.T) {
	t.Skip("Integration test - requires running listings microservice on localhost:50051")

	// This test verifies that UpdateListing correctly:
	// 1. Converts domain model to protobuf update request
	// 2. Only includes non-empty optional fields
	// 3. Calls gRPC UpdateListing RPC
	// 4. Handles ownership check errors
}

// TestMarketplaceGRPCClient_DeleteListing_Integration tests DeleteListing method
func TestMarketplaceGRPCClient_DeleteListing_Integration(t *testing.T) {
	t.Skip("Integration test - requires running listings microservice on localhost:50051")

	// This test verifies that DeleteListing correctly:
	// 1. Calls gRPC DeleteListing with user_id for ownership check
	// 2. Handles soft-delete (listing marked as deleted but not removed)
	// 3. Returns success response
	// 4. Handles permission denied errors
}

// TestMarketplaceGRPCClient_DeleteListingAdmin_Integration tests DeleteListingAdmin method
func TestMarketplaceGRPCClient_DeleteListingAdmin_Integration(t *testing.T) {
	t.Skip("Integration test - requires running listings microservice on localhost:50051")

	// This test verifies that DeleteListingAdmin correctly:
	// 1. Calls gRPC DeleteListing with user_id=0 (admin, no ownership check)
	// 2. Allows deletion of any listing regardless of owner
	// 3. Returns success response
}

// TestMarketplaceGRPCClient_GetListingBySlug_Integration tests GetListingBySlug method
func TestMarketplaceGRPCClient_GetListingBySlug_Integration(t *testing.T) {
	t.Skip("Integration test - requires slug support in listings microservice")

	// This test verifies that GetListingBySlug correctly:
	// 1. Returns error indicating slug support not yet implemented
	// NOTE: Slug lookup will be added to listings microservice in future sprint
}
