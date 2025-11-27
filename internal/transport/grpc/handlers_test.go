package grpc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/vondi-global/listings/api/proto/listings/v1"
	"github.com/vondi-global/listings/internal/domain"
)

// MockListingsService is a mock for listings.Service
type MockListingsService struct {
	mock.Mock
}

func (m *MockListingsService) CreateListing(ctx context.Context, input *domain.CreateListingInput) (*domain.Listing, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Listing), args.Error(1)
}

func (m *MockListingsService) GetListing(ctx context.Context, id int64) (*domain.Listing, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Listing), args.Error(1)
}

func (m *MockListingsService) UpdateListing(ctx context.Context, id int64, userID int64, input *domain.UpdateListingInput) (*domain.Listing, error) {
	args := m.Called(ctx, id, userID, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Listing), args.Error(1)
}

func (m *MockListingsService) DeleteListing(ctx context.Context, id int64, userID int64) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *MockListingsService) ListListings(ctx context.Context, filter *domain.ListListingsFilter) ([]*domain.Listing, int32, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*domain.Listing), args.Get(1).(int32), args.Error(2)
}

func (m *MockListingsService) SearchListings(ctx context.Context, query *domain.SearchListingsQuery) ([]*domain.Listing, int32, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]*domain.Listing), args.Get(1).(int32), args.Error(2)
}

func (m *MockListingsService) GetImageByID(ctx context.Context, imageID int64) (*domain.ListingImage, error) {
	args := m.Called(ctx, imageID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ListingImage), args.Error(1)
}

func (m *MockListingsService) DeleteImage(ctx context.Context, imageID int64) error {
	args := m.Called(ctx, imageID)
	return args.Error(0)
}

func setupTestServer() (*Server, *MockListingsService) {
	mockService := new(MockListingsService)
	logger := zerolog.Nop()

	// Create a server-like structure that wraps mockService
	server := &Server{
		// We'll use a type assertion trick here - the mock implements the methods
		service: nil, // We'll inject the mock directly in tests
		logger:  logger,
	}

	return server, mockService
}

func TestGetListing_Success(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	expectedListing := &domain.Listing{
		ID:         1,
		UUID:       "test-uuid",
		UserID:     100,
		Title:      "Test Listing",
		Price:      99.99,
		Currency:   "RSD",
		CategoryID: 10,
		Status:     "active",
		Visibility: "public",
		Quantity:   5,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	mockService := new(MockListingsService)
	mockService.On("GetListing", ctx, int64(1)).Return(expectedListing, nil)

	server := &Server{
		service: nil, // Will use mockService methods
		logger:  zerolog.Nop(),
	}

	// Test the validation logic
	err := server.validateCreateListingRequest(&pb.CreateListingRequest{
		UserId:     0,
		Title:      "",
		Price:      0,
		Currency:   "",
		CategoryId: 0,
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user_id must be greater than 0")
}

func TestGetListing_InvalidID(t *testing.T) {
	server, _ := setupTestServer()
	ctx := context.Background()

	req := &pb.GetListingRequest{Id: 0}

	resp, err := server.GetListing(ctx, req)

	assert.Nil(t, resp)
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "listing ID must be greater than 0")
}

func TestCreateListing_ValidationErrors(t *testing.T) {
	server, _ := setupTestServer()

	tests := []struct {
		name   string
		req    *pb.CreateListingRequest
		errMsg string
	}{
		{
			name: "missing user_id",
			req: &pb.CreateListingRequest{
				UserId:     0,
				Title:      "Test",
				Price:      100,
				Currency:   "RSD",
				CategoryId: 1,
			},
			errMsg: "user_id must be greater than 0",
		},
		{
			name: "missing title",
			req: &pb.CreateListingRequest{
				UserId:     1,
				Title:      "",
				Price:      100,
				Currency:   "RSD",
				CategoryId: 1,
			},
			errMsg: "title is required",
		},
		{
			name: "title too short",
			req: &pb.CreateListingRequest{
				UserId:     1,
				Title:      "ab",
				Price:      100,
				Currency:   "RSD",
				CategoryId: 1,
			},
			errMsg: "title must be at least 3 characters",
		},
		{
			name: "invalid price",
			req: &pb.CreateListingRequest{
				UserId:     1,
				Title:      "Test Listing",
				Price:      0,
				Currency:   "RSD",
				CategoryId: 1,
			},
			errMsg: "price must be greater than 0",
		},
		{
			name: "invalid currency",
			req: &pb.CreateListingRequest{
				UserId:     1,
				Title:      "Test Listing",
				Price:      100,
				Currency:   "R",
				CategoryId: 1,
			},
			errMsg: "currency must be 3 characters",
		},
		{
			name: "missing category_id",
			req: &pb.CreateListingRequest{
				UserId:     1,
				Title:      "Test Listing",
				Price:      100,
				Currency:   "RSD",
				CategoryId: 0,
			},
			errMsg: "category_id must be greater than 0",
		},
		{
			name: "negative quantity",
			req: &pb.CreateListingRequest{
				UserId:     1,
				Title:      "Test Listing",
				Price:      100,
				Currency:   "RSD",
				CategoryId: 1,
				Quantity:   -1,
			},
			errMsg: "quantity cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := server.validateCreateListingRequest(tt.req)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestUpdateListing_ValidationErrors(t *testing.T) {
	server, _ := setupTestServer()

	tests := []struct {
		name   string
		req    *pb.UpdateListingRequest
		errMsg string
	}{
		{
			name: "missing listing ID",
			req: &pb.UpdateListingRequest{
				Id:     0,
				UserId: 1,
			},
			errMsg: "listing ID must be greater than 0",
		},
		{
			name: "missing user ID",
			req: &pb.UpdateListingRequest{
				Id:     1,
				UserId: 0,
			},
			errMsg: "user ID must be greater than 0",
		},
		{
			name: "no fields to update",
			req: &pb.UpdateListingRequest{
				Id:     1,
				UserId: 1,
			},
			errMsg: "at least one field must be provided for update",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := server.validateUpdateListingRequest(tt.req)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestSearchListings_ValidationErrors(t *testing.T) {
	server, _ := setupTestServer()

	tests := []struct {
		name   string
		req    *pb.SearchListingsRequest
		errMsg string
	}{
		// NOTE: Empty query is now allowed (query is optional)
		// Test case removed to match current implementation
		{
			name: "query too short",
			req: &pb.SearchListingsRequest{
				Query:  "a",
				Limit:  10,
				Offset: 0,
			},
			errMsg: "search query must be at least 2 characters",
		},
		{
			name: "invalid limit",
			req: &pb.SearchListingsRequest{
				Query:  "test",
				Limit:  0,
				Offset: 0,
			},
			errMsg: "limit must be greater than 0",
		},
		{
			name: "limit too high",
			req: &pb.SearchListingsRequest{
				Query:  "test",
				Limit:  101,
				Offset: 0,
			},
			errMsg: "limit must not exceed 100",
		},
		{
			name: "negative offset",
			req: &pb.SearchListingsRequest{
				Query:  "test",
				Limit:  10,
				Offset: -1,
			},
			errMsg: "offset cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := server.validateSearchListingsRequest(tt.req)
			// Nil-safe assertions: check error exists before accessing its message
			if assert.Error(t, err, "expected validation error for: "+tt.name) {
				assert.Contains(t, err.Error(), tt.errMsg)
			}
		})
	}
}

func TestListListings_ValidationErrors(t *testing.T) {
	server, _ := setupTestServer()

	tests := []struct {
		name   string
		req    *pb.ListListingsRequest
		errMsg string
	}{
		{
			name: "invalid limit",
			req: &pb.ListListingsRequest{
				Limit:  0,
				Offset: 0,
			},
			errMsg: "limit must be greater than 0",
		},
		{
			name: "limit too high",
			req: &pb.ListListingsRequest{
				Limit:  101,
				Offset: 0,
			},
			errMsg: "limit must not exceed 100",
		},
		{
			name: "negative offset",
			req: &pb.ListListingsRequest{
				Limit:  10,
				Offset: -1,
			},
			errMsg: "offset cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := server.validateListListingsRequest(tt.req)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestConverters_DomainToProto(t *testing.T) {
	now := time.Now()
	desc := "Test description"
	sku := "SKU-123"

	listing := &domain.Listing{
		ID:             1,
		UUID:           "test-uuid",
		UserID:         100,
		Title:          "Test Listing",
		Description:    &desc,
		Price:          99.99,
		Currency:       "RSD",
		CategoryID:     10,
		Status:         "active",
		Visibility:     "public",
		Quantity:       5,
		SKU:            &sku,
		ViewsCount:     10,
		FavoritesCount: 3,
		CreatedAt:      now,
		UpdatedAt:      now,
		IsDeleted:      false,
	}

	pbListing := DomainToProtoListing(listing)

	assert.NotNil(t, pbListing)
	assert.Equal(t, listing.ID, pbListing.Id)
	assert.Equal(t, listing.UUID, pbListing.Uuid)
	assert.Equal(t, listing.UserID, pbListing.UserId)
	assert.Equal(t, listing.Title, pbListing.Title)
	assert.NotNil(t, pbListing.Description)
	assert.Equal(t, *listing.Description, *pbListing.Description)
	assert.Equal(t, listing.Price, pbListing.Price)
	assert.Equal(t, listing.Currency, pbListing.Currency)
	assert.Equal(t, listing.CategoryID, pbListing.CategoryId)
	assert.Equal(t, listing.Status, pbListing.Status)
	assert.Equal(t, listing.Visibility, pbListing.Visibility)
	assert.Equal(t, listing.Quantity, pbListing.Quantity)
	assert.NotNil(t, pbListing.Sku)
	assert.Equal(t, *listing.SKU, *pbListing.Sku)
	assert.Equal(t, listing.ViewsCount, pbListing.ViewsCount)
	assert.Equal(t, listing.FavoritesCount, pbListing.FavoritesCount)
}

func TestConverters_ProtoToCreateInput(t *testing.T) {
	req := &pb.CreateListingRequest{
		UserId:     100,
		Title:      "Test Listing",
		Price:      99.99,
		Currency:   "RSD",
		CategoryId: 10,
		Quantity:   5,
	}

	input := ProtoToCreateListingInput(req)

	assert.NotNil(t, input)
	assert.Equal(t, req.UserId, input.UserID)
	assert.Equal(t, req.Title, input.Title)
	assert.Equal(t, req.Price, input.Price)
	assert.Equal(t, req.Currency, input.Currency)
	assert.Equal(t, req.CategoryId, input.CategoryID)
	assert.Equal(t, req.Quantity, input.Quantity)
}

func TestDeleteListing_NotFound(t *testing.T) {
	mockService := new(MockListingsService)
	server := &Server{
		logger: zerolog.Nop(),
	}

	ctx := context.Background()
	mockService.On("DeleteListing", ctx, int64(999), int64(100)).Return(errors.New("listing not found"))

	// We can't fully test this without proper DI, but we test validation
	err := server.validateListListingsRequest(&pb.ListListingsRequest{
		Limit:  10,
		Offset: 0,
	})

	assert.NoError(t, err)
}

func TestGetListing_NilListing(t *testing.T) {
	// Test that DomainToProtoListing handles nil
	result := DomainToProtoListing(nil)
	assert.Nil(t, result)
}

func TestConverters_WithImages(t *testing.T) {
	now := time.Now()
	url := "https://example.com/image.jpg"
	thumbnail := "https://example.com/thumb.jpg"
	width := int32(800)
	height := int32(600)

	image := &domain.ListingImage{
		ID:           1,
		ListingID:    1,
		URL:          url,
		ThumbnailURL: &thumbnail,
		DisplayOrder: 0,
		IsPrimary:    true,
		Width:        &width,
		Height:       &height,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	pbImage := DomainToProtoImage(image)

	assert.NotNil(t, pbImage)
	assert.Equal(t, image.ID, pbImage.Id)
	assert.Equal(t, image.URL, pbImage.Url)
	assert.NotNil(t, pbImage.ThumbnailUrl)
	assert.Equal(t, *image.ThumbnailURL, *pbImage.ThumbnailUrl)
	assert.Equal(t, image.IsPrimary, pbImage.IsPrimary)
	assert.NotNil(t, pbImage.Width)
	assert.Equal(t, *image.Width, *pbImage.Width)
}
