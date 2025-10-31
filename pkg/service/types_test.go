package service

import (
	"testing"
	"time"
)

func TestListingConstantsExist(t *testing.T) {
	// Test that status constants exist
	statuses := []string{
		StatusDraft,
		StatusActive,
		StatusInactive,
		StatusSold,
		StatusArchived,
	}

	for _, status := range statuses {
		if status == "" {
			t.Errorf("status constant is empty")
		}
	}

	// Test that visibility constants exist
	visibilities := []string{
		VisibilityPublic,
		VisibilityPrivate,
		VisibilityUnlisted,
	}

	for _, visibility := range visibilities {
		if visibility == "" {
			t.Errorf("visibility constant is empty")
		}
	}
}

func TestListingStructure(t *testing.T) {
	// Test that Listing struct can be created
	now := time.Now()
	listing := &Listing{
		ID:         1,
		UUID:       "test-uuid",
		UserID:     100,
		Title:      "Test Listing",
		Price:      99.99,
		Currency:   "USD",
		CategoryID: 1,
		Status:     StatusActive,
		Visibility: VisibilityPublic,
		Quantity:   10,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if listing.ID != 1 {
		t.Errorf("expected ID 1, got %d", listing.ID)
	}

	if listing.Title != "Test Listing" {
		t.Errorf("expected title 'Test Listing', got %s", listing.Title)
	}

	if listing.Price != 99.99 {
		t.Errorf("expected price 99.99, got %f", listing.Price)
	}
}

func TestCreateListingRequest(t *testing.T) {
	req := &CreateListingRequest{
		UserID:     100,
		Title:      "Test",
		Price:      10.0,
		Currency:   "USD",
		CategoryID: 1,
		Quantity:   5,
	}

	if req.UserID != 100 {
		t.Errorf("expected UserID 100, got %d", req.UserID)
	}

	if req.Title != "Test" {
		t.Errorf("expected Title 'Test', got %s", req.Title)
	}
}

func TestUpdateListingRequest(t *testing.T) {
	title := "Updated Title"
	price := 20.0
	quantity := int32(10)
	status := StatusActive

	req := &UpdateListingRequest{
		Title:    &title,
		Price:    &price,
		Quantity: &quantity,
		Status:   &status,
	}

	if req.Title == nil || *req.Title != "Updated Title" {
		t.Errorf("expected Title 'Updated Title', got %v", req.Title)
	}

	if req.Price == nil || *req.Price != 20.0 {
		t.Errorf("expected Price 20.0, got %v", req.Price)
	}
}

func TestSearchListingsRequest(t *testing.T) {
	categoryID := int64(1)
	minPrice := 10.0
	maxPrice := 100.0

	req := &SearchListingsRequest{
		Query:      "test",
		CategoryID: &categoryID,
		MinPrice:   &minPrice,
		MaxPrice:   &maxPrice,
		Limit:      10,
		Offset:     0,
	}

	if req.Query != "test" {
		t.Errorf("expected Query 'test', got %s", req.Query)
	}

	if req.CategoryID == nil || *req.CategoryID != 1 {
		t.Errorf("expected CategoryID 1, got %v", req.CategoryID)
	}
}

func TestListListingsRequest(t *testing.T) {
	userID := int64(100)
	storefrontID := int64(50)
	status := StatusActive

	req := &ListListingsRequest{
		UserID:       &userID,
		StorefrontID: &storefrontID,
		Status:       &status,
		Limit:        20,
		Offset:       0,
	}

	if req.UserID == nil || *req.UserID != 100 {
		t.Errorf("expected UserID 100, got %v", req.UserID)
	}

	if req.StorefrontID == nil || *req.StorefrontID != 50 {
		t.Errorf("expected StorefrontID 50, got %v", req.StorefrontID)
	}

	if req.Limit != 20 {
		t.Errorf("expected Limit 20, got %d", req.Limit)
	}
}
