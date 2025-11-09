package grpc

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/sveturs/listings/api/proto/listings/v1"
	"github.com/sveturs/listings/internal/domain"
)

// MockStorefrontService extends MockListingsService with storefront methods
type MockStorefrontService struct {
	mock.Mock
}

func (m *MockStorefrontService) GetStorefront(ctx context.Context, storefrontID int64) (*domain.Storefront, error) {
	args := m.Called(ctx, storefrontID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Storefront), args.Error(1)
}

func (m *MockStorefrontService) GetStorefrontBySlug(ctx context.Context, slug string) (*domain.Storefront, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Storefront), args.Error(1)
}

func (m *MockStorefrontService) ListStorefronts(ctx context.Context, limit, offset int) ([]*domain.Storefront, int64, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Storefront), args.Get(1).(int64), args.Error(2)
}

// StorefrontServiceInterface defines the interface for storefront operations
type StorefrontServiceInterface interface {
	GetStorefront(ctx context.Context, storefrontID int64) (*domain.Storefront, error)
	GetStorefrontBySlug(ctx context.Context, slug string) (*domain.Storefront, error)
	ListStorefronts(ctx context.Context, limit, offset int) ([]*domain.Storefront, int64, error)
}

// Ensure MockStorefrontService implements the interface
var _ StorefrontServiceInterface = (*MockStorefrontService)(nil)

// serverWithStorefronts wraps Server for storefront testing
type serverWithStorefronts struct {
	pb.UnimplementedListingsServiceServer
	service StorefrontServiceInterface
	logger  zerolog.Logger
}

// setupStorefrontsTestServer creates a test server with mock service
func setupStorefrontsTestServer() (*serverWithStorefronts, *MockStorefrontService) {
	mockService := new(MockStorefrontService)
	logger := zerolog.Nop()

	server := &serverWithStorefronts{
		service: mockService,
		logger:  logger,
	}

	return server, mockService
}

// GetStorefront implements the gRPC handler for testing
func (s *serverWithStorefronts) GetStorefront(ctx context.Context, req *pb.GetStorefrontRequest) (*pb.StorefrontResponse, error) {
	s.logger.Debug().Int64("storefront_id", req.Id).Msg("GetStorefront called")

	// Validation
	if req.Id <= 0 {
		s.logger.Warn().
			Int64("storefront_id", req.Id).
			Msg("Invalid storefront ID")
		return nil, status.Error(codes.InvalidArgument, "storefront ID must be greater than 0")
	}

	// Get storefront from service
	storefront, err := s.service.GetStorefront(ctx, req.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Warn().
				Int64("storefront_id", req.Id).
				Msg("Storefront not found")
			return nil, status.Error(codes.NotFound, "storefront not found")
		}

		s.logger.Error().
			Err(err).
			Int64("storefront_id", req.Id).
			Msg("Failed to get storefront")
		return nil, status.Error(codes.Internal, "failed to get storefront: "+err.Error())
	}

	// Convert to proto
	protoStorefront := StorefrontToProto(storefront)

	s.logger.Info().
		Int64("storefront_id", req.Id).
		Str("slug", storefront.Slug).
		Msg("Storefront retrieved successfully")

	return &pb.StorefrontResponse{
		Storefront: protoStorefront,
	}, nil
}

// GetStorefrontBySlug implements the gRPC handler for testing
func (s *serverWithStorefronts) GetStorefrontBySlug(ctx context.Context, req *pb.GetStorefrontBySlugRequest) (*pb.StorefrontResponse, error) {
	s.logger.Debug().Str("slug", req.Slug).Msg("GetStorefrontBySlug called")

	// Validation
	if req.Slug == "" {
		s.logger.Warn().Msg("Empty slug provided")
		return nil, status.Error(codes.InvalidArgument, "slug cannot be empty")
	}

	// Get storefront from service
	storefront, err := s.service.GetStorefrontBySlug(ctx, req.Slug)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Warn().
				Str("slug", req.Slug).
				Msg("Storefront not found by slug")
			return nil, status.Error(codes.NotFound, "storefront not found")
		}

		s.logger.Error().
			Err(err).
			Str("slug", req.Slug).
			Msg("Failed to get storefront by slug")
		return nil, status.Error(codes.Internal, "failed to get storefront: "+err.Error())
	}

	// Convert to proto
	protoStorefront := StorefrontToProto(storefront)

	s.logger.Info().
		Int64("storefront_id", storefront.ID).
		Str("slug", req.Slug).
		Msg("Storefront retrieved by slug successfully")

	return &pb.StorefrontResponse{
		Storefront: protoStorefront,
	}, nil
}

// ListStorefronts implements the gRPC handler for testing
func (s *serverWithStorefronts) ListStorefronts(ctx context.Context, req *pb.ListStorefrontsRequest) (*pb.ListStorefrontsResponse, error) {
	s.logger.Debug().
		Int32("limit", req.Limit).
		Int32("offset", req.Offset).
		Msg("ListStorefronts called")

	// Validation and defaults
	limit := req.Limit
	if limit <= 0 {
		limit = 20 // Default limit
	}
	if limit > 100 {
		limit = 100 // Max limit
	}

	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	// Get storefronts from service
	storefronts, total, err := s.service.ListStorefronts(ctx, int(limit), int(offset))
	if err != nil {
		s.logger.Error().
			Err(err).
			Int32("limit", limit).
			Int32("offset", offset).
			Msg("Failed to list storefronts")
		return nil, status.Error(codes.Internal, "failed to list storefronts: "+err.Error())
	}

	// Convert to proto
	protoStorefronts := make([]*pb.Storefront, 0, len(storefronts))
	for _, sf := range storefronts {
		protoStorefronts = append(protoStorefronts, StorefrontToProto(sf))
	}

	s.logger.Info().
		Int("count", len(storefronts)).
		Int64("total", total).
		Int32("limit", limit).
		Int32("offset", offset).
		Msg("Storefronts listed successfully")

	return &pb.ListStorefrontsResponse{
		Storefronts: protoStorefronts,
		Total:       int32(total),
	}, nil
}

// ============================================================================
// Test Cases
// ============================================================================

// TestGetStorefront_Success tests successful storefront retrieval
func TestGetStorefront_Success(t *testing.T) {
	server, mockService := setupStorefrontsTestServer()
	ctx := context.Background()

	name := "Test Storefront"
	description := "Test description"
	expectedStorefront := &domain.Storefront{
		ID:          1,
		UserID:      123,
		Name:        name,
		Slug:        "test-storefront",
		Description: &description,
		Country:     "RS",
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockService.On("GetStorefront", ctx, int64(1)).Return(expectedStorefront, nil)

	req := &pb.GetStorefrontRequest{Id: 1}
	resp, err := server.GetStorefront(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedStorefront.ID, resp.Storefront.Id)
	assert.Equal(t, expectedStorefront.Name, resp.Storefront.Name)
	assert.Equal(t, expectedStorefront.Slug, resp.Storefront.Slug)
	mockService.AssertExpectations(t)
}

// TestGetStorefront_NotFound tests storefront not found scenario
func TestGetStorefront_NotFound(t *testing.T) {
	server, mockService := setupStorefrontsTestServer()
	ctx := context.Background()

	mockService.On("GetStorefront", ctx, int64(999)).Return(nil, sql.ErrNoRows)

	req := &pb.GetStorefrontRequest{Id: 999}
	resp, err := server.GetStorefront(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	mockService.AssertExpectations(t)
}

// TestGetStorefront_InvalidID tests invalid storefront ID
func TestGetStorefront_InvalidID(t *testing.T) {
	server, _ := setupStorefrontsTestServer()
	ctx := context.Background()

	req := &pb.GetStorefrontRequest{Id: 0}
	resp, err := server.GetStorefront(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

// TestGetStorefrontBySlug_Success tests successful slug-based retrieval
func TestGetStorefrontBySlug_Success(t *testing.T) {
	server, mockService := setupStorefrontsTestServer()
	ctx := context.Background()

	expectedStorefront := &domain.Storefront{
		ID:       1,
		UserID:   123,
		Name:     "Test Storefront",
		Slug:     "test-storefront",
		Country:  "RS",
		IsActive: true,
	}

	mockService.On("GetStorefrontBySlug", ctx, "test-storefront").Return(expectedStorefront, nil)

	req := &pb.GetStorefrontBySlugRequest{Slug: "test-storefront"}
	resp, err := server.GetStorefrontBySlug(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, expectedStorefront.Slug, resp.Storefront.Slug)
	mockService.AssertExpectations(t)
}

// TestGetStorefrontBySlug_EmptySlug tests empty slug validation
func TestGetStorefrontBySlug_EmptySlug(t *testing.T) {
	server, _ := setupStorefrontsTestServer()
	ctx := context.Background()

	req := &pb.GetStorefrontBySlugRequest{Slug: ""}
	resp, err := server.GetStorefrontBySlug(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

// TestListStorefronts_Success tests successful listing
func TestListStorefronts_Success(t *testing.T) {
	server, mockService := setupStorefrontsTestServer()
	ctx := context.Background()

	expectedStorefronts := []*domain.Storefront{
		{ID: 1, Name: "Store 1", Slug: "store-1", Country: "RS", IsActive: true},
		{ID: 2, Name: "Store 2", Slug: "store-2", Country: "RS", IsActive: true},
	}

	mockService.On("ListStorefronts", ctx, 20, 0).Return(expectedStorefronts, int64(2), nil)

	req := &pb.ListStorefrontsRequest{Limit: 20, Offset: 0}
	resp, err := server.ListStorefronts(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Storefronts, 2)
	assert.Equal(t, int32(2), resp.Total)
	mockService.AssertExpectations(t)
}

// TestListStorefronts_DefaultLimit tests default limit behavior
func TestListStorefronts_DefaultLimit(t *testing.T) {
	server, mockService := setupStorefrontsTestServer()
	ctx := context.Background()

	mockService.On("ListStorefronts", ctx, 20, 0).Return([]*domain.Storefront{}, int64(0), nil)

	req := &pb.ListStorefrontsRequest{Limit: 0, Offset: 0} // Invalid limit, should use default 20
	resp, err := server.ListStorefronts(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int32(0), resp.Total)
	mockService.AssertExpectations(t)
}

// TestListStorefronts_MaxLimit tests max limit enforcement
func TestListStorefronts_MaxLimit(t *testing.T) {
	server, mockService := setupStorefrontsTestServer()
	ctx := context.Background()

	// Should enforce max limit of 100
	mockService.On("ListStorefronts", ctx, 100, 0).Return([]*domain.Storefront{}, int64(0), nil)

	req := &pb.ListStorefrontsRequest{Limit: 200, Offset: 0} // Exceeds max, should cap at 100
	resp, err := server.ListStorefronts(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	mockService.AssertExpectations(t)
}

// TestListStorefronts_NegativeOffset tests negative offset handling
func TestListStorefronts_NegativeOffset(t *testing.T) {
	server, mockService := setupStorefrontsTestServer()
	ctx := context.Background()

	// Negative offset should be normalized to 0
	mockService.On("ListStorefronts", ctx, 20, 0).Return([]*domain.Storefront{}, int64(0), nil)

	req := &pb.ListStorefrontsRequest{Limit: 20, Offset: -10}
	resp, err := server.ListStorefronts(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	mockService.AssertExpectations(t)
}
