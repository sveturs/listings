package grpc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	listingssvcv1 "github.com/sveturs/listings/api/proto/listings/v1"
)

// ============================================================================
// Mock Analytics Service
// ============================================================================

type mockAnalyticsService struct {
	getOverviewStatsFunc func(ctx context.Context, req *listingssvcv1.GetOverviewStatsRequest, userID int64, isAdmin bool) (*listingssvcv1.GetOverviewStatsResponse, error)
	getListingStatsFunc  func(ctx context.Context, req *listingssvcv1.GetListingStatsRequest, userID int64, isAdmin bool) (*listingssvcv1.GetListingStatsResponse, error)
	getTrendingStatsFunc func(ctx context.Context, req *listingssvcv1.GetTrendingStatsRequest) (*listingssvcv1.GetTrendingStatsResponse, error)
}

func (m *mockAnalyticsService) GetOverviewStats(
	ctx context.Context,
	req *listingssvcv1.GetOverviewStatsRequest,
	userID int64,
	isAdmin bool,
) (*listingssvcv1.GetOverviewStatsResponse, error) {
	if m.getOverviewStatsFunc != nil {
		return m.getOverviewStatsFunc(ctx, req, userID, isAdmin)
	}
	return &listingssvcv1.GetOverviewStatsResponse{}, nil
}

func (m *mockAnalyticsService) GetListingStats(
	ctx context.Context,
	req *listingssvcv1.GetListingStatsRequest,
	userID int64,
	isAdmin bool,
) (*listingssvcv1.GetListingStatsResponse, error) {
	if m.getListingStatsFunc != nil {
		return m.getListingStatsFunc(ctx, req, userID, isAdmin)
	}
	return &listingssvcv1.GetListingStatsResponse{}, nil
}

func (m *mockAnalyticsService) GetTrendingStats(
	ctx context.Context,
	req *listingssvcv1.GetTrendingStatsRequest,
) (*listingssvcv1.GetTrendingStatsResponse, error) {
	if m.getTrendingStatsFunc != nil {
		return m.getTrendingStatsFunc(ctx, req)
	}
	return &listingssvcv1.GetTrendingStatsResponse{}, nil
}

// ============================================================================
// Helper Functions for Tests
// ============================================================================

func datePtr(t time.Time) *timestamppb.Timestamp {
	return timestamppb.New(t)
}

func analyticsTestTimeNow() time.Time {
	return time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)
}

func analyticsTestTimeDaysAgo(days int) time.Time {
	return analyticsTestTimeNow().AddDate(0, 0, -days)
}

// analyticsContainsError checks if error message contains a substring (case-insensitive)
func analyticsContainsError(err error, substr string) bool {
	if err == nil {
		return false
	}
	return contains(err.Error(), substr)
}

// ============================================================================
// Validation Tests for GetOverviewStatsRequest
// ============================================================================

func TestValidateOverviewStatsRequest_Success(t *testing.T) {
	// Arrange
	now := analyticsTestTimeNow()
	dateFrom := analyticsTestTimeDaysAgo(7)

	server := &Server{
		logger: zerolog.Nop(),
	}

	req := &listingssvcv1.GetOverviewStatsRequest{
		DateFrom: datePtr(dateFrom),
		DateTo:   datePtr(now),
	}

	// Act
	err := server.validateOverviewStatsRequestHelper(req)

	// Assert
	require.NoError(t, err)
}

func TestValidateOverviewStatsRequest_WithFilters(t *testing.T) {
	// Arrange
	now := analyticsTestTimeNow()
	dateFrom := analyticsTestTimeDaysAgo(30)
	categoryID := int64(1001)
	storefrontID := int64(42)
	listingType := "b2c"

	server := &Server{
		logger: zerolog.Nop(),
	}

	req := &listingssvcv1.GetOverviewStatsRequest{
		DateFrom:     datePtr(dateFrom),
		DateTo:       datePtr(now),
		CategoryId:   &categoryID,
		StorefrontId: &storefrontID,
		ListingType:  &listingType,
	}

	// Act
	err := server.validateOverviewStatsRequestHelper(req)

	// Assert
	require.NoError(t, err)
}

func TestValidateOverviewStatsRequest_NilRequest(t *testing.T) {
	// Arrange
	server := &Server{
		logger: zerolog.Nop(),
	}

	// Act
	err := server.validateOverviewStatsRequestHelper(nil)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be nil")
}

func TestValidateOverviewStatsRequest_MissingDateFrom(t *testing.T) {
	// Arrange
	now := analyticsTestTimeNow()
	server := &Server{
		logger: zerolog.Nop(),
	}

	req := &listingssvcv1.GetOverviewStatsRequest{
		DateFrom: nil,
		DateTo:   datePtr(now),
	}

	// Act
	err := server.validateOverviewStatsRequestHelper(req)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "date_from")
}

func TestValidateOverviewStatsRequest_MissingDateTo(t *testing.T) {
	// Arrange
	dateFrom := analyticsTestTimeDaysAgo(7)
	server := &Server{
		logger: zerolog.Nop(),
	}

	req := &listingssvcv1.GetOverviewStatsRequest{
		DateFrom: datePtr(dateFrom),
		DateTo:   nil,
	}

	// Act
	err := server.validateOverviewStatsRequestHelper(req)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "date_to")
}

func TestValidateOverviewStatsRequest_InvalidDateRange(t *testing.T) {
	// Arrange
	now := analyticsTestTimeNow()
	dateAfter := now.AddDate(0, 0, 1)

	server := &Server{
		logger: zerolog.Nop(),
	}

	req := &listingssvcv1.GetOverviewStatsRequest{
		DateFrom: datePtr(now),
		DateTo:   datePtr(dateAfter.AddDate(0, 0, -2)),
	}

	// Act
	err := server.validateOverviewStatsRequestHelper(req)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "date_from must be before")
}

func TestValidateOverviewStatsRequest_InvalidListingType(t *testing.T) {
	// Arrange
	now := analyticsTestTimeNow()
	dateFrom := analyticsTestTimeDaysAgo(7)
	invalidType := "invalid_type"

	server := &Server{
		logger: zerolog.Nop(),
	}

	req := &listingssvcv1.GetOverviewStatsRequest{
		DateFrom:    datePtr(dateFrom),
		DateTo:      datePtr(now),
		ListingType: &invalidType,
	}

	// Act
	err := server.validateOverviewStatsRequestHelper(req)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "listing_type")
}

func TestValidateOverviewStatsRequest_ValidListingTypes(t *testing.T) {
	tests := []struct {
		name        string
		listingType string
		expectError bool
	}{
		{
			name:        "b2c type is valid",
			listingType: "b2c",
			expectError: false,
		},
		{
			name:        "c2c type is valid",
			listingType: "c2c",
			expectError: false,
		},
		{
			name:        "invalid type fails",
			listingType: "b3c",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			now := analyticsTestTimeNow()
			dateFrom := analyticsTestTimeDaysAgo(7)
			listingType := tt.listingType

			server := &Server{
				logger: zerolog.Nop(),
			}

			req := &listingssvcv1.GetOverviewStatsRequest{
				DateFrom:    datePtr(dateFrom),
				DateTo:      datePtr(now),
				ListingType: &listingType,
			}

			// Act
			err := server.validateOverviewStatsRequestHelper(req)

			// Assert
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// ============================================================================
// Validation Tests for GetListingStatsRequest
// ============================================================================

func TestValidateListingStatsRequest_Success_WithListingID(t *testing.T) {
	// Arrange
	now := analyticsTestTimeNow()
	dateFrom := analyticsTestTimeDaysAgo(7)
	listingID := int64(281)

	server := &Server{
		logger: zerolog.Nop(),
	}

	req := &listingssvcv1.GetListingStatsRequest{
		Identifier: &listingssvcv1.GetListingStatsRequest_ListingId{ListingId: listingID},
		DateFrom:   datePtr(dateFrom),
		DateTo:     datePtr(now),
	}

	// Act
	err := server.validateListingStatsRequestHelper(req)

	// Assert
	require.NoError(t, err)
}

func TestValidateListingStatsRequest_Success_WithProductID(t *testing.T) {
	// Arrange
	now := analyticsTestTimeNow()
	dateFrom := analyticsTestTimeDaysAgo(7)
	productID := int64(501)

	server := &Server{
		logger: zerolog.Nop(),
	}

	req := &listingssvcv1.GetListingStatsRequest{
		Identifier: &listingssvcv1.GetListingStatsRequest_ProductId{ProductId: productID},
		DateFrom:   datePtr(dateFrom),
		DateTo:     datePtr(now),
	}

	// Act
	err := server.validateListingStatsRequestHelper(req)

	// Assert
	require.NoError(t, err)
}

func TestValidateListingStatsRequest_NilRequest(t *testing.T) {
	// Arrange
	server := &Server{
		logger: zerolog.Nop(),
	}

	// Act
	err := server.validateListingStatsRequestHelper(nil)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be nil")
}

func TestValidateListingStatsRequest_MissingIdentifier(t *testing.T) {
	// Arrange
	now := analyticsTestTimeNow()
	dateFrom := analyticsTestTimeDaysAgo(7)

	server := &Server{
		logger: zerolog.Nop(),
	}

	req := &listingssvcv1.GetListingStatsRequest{
		DateFrom: datePtr(dateFrom),
		DateTo:   datePtr(now),
	}

	// Act
	err := server.validateListingStatsRequestHelper(req)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "listing_id or product_id")
}

func TestValidateListingStatsRequest_MissingDateFrom(t *testing.T) {
	// Arrange
	now := analyticsTestTimeNow()
	listingID := int64(281)

	server := &Server{
		logger: zerolog.Nop(),
	}

	req := &listingssvcv1.GetListingStatsRequest{
		Identifier: &listingssvcv1.GetListingStatsRequest_ListingId{ListingId: listingID},
		DateFrom:   nil,
		DateTo:     datePtr(now),
	}

	// Act
	err := server.validateListingStatsRequestHelper(req)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "date_from")
}

func TestValidateListingStatsRequest_MissingDateTo(t *testing.T) {
	// Arrange
	dateFrom := analyticsTestTimeDaysAgo(7)
	listingID := int64(281)

	server := &Server{
		logger: zerolog.Nop(),
	}

	req := &listingssvcv1.GetListingStatsRequest{
		Identifier: &listingssvcv1.GetListingStatsRequest_ListingId{ListingId: listingID},
		DateFrom:   datePtr(dateFrom),
		DateTo:     nil,
	}

	// Act
	err := server.validateListingStatsRequestHelper(req)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "date_to")
}

func TestValidateListingStatsRequest_InvalidDateRange(t *testing.T) {
	// Arrange
	now := analyticsTestTimeNow()
	listingID := int64(281)

	server := &Server{
		logger: zerolog.Nop(),
	}

	req := &listingssvcv1.GetListingStatsRequest{
		Identifier: &listingssvcv1.GetListingStatsRequest_ListingId{ListingId: listingID},
		DateFrom:   datePtr(now),
		DateTo:     datePtr(now.AddDate(0, 0, -1)),
	}

	// Act
	err := server.validateListingStatsRequestHelper(req)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "date_from must be before")
}

// ============================================================================
// Error Mapping Tests
// ============================================================================

func TestMapAnalyticsError_PermissionErrors(t *testing.T) {
	tests := []struct {
		name             string
		err              error
		operation        string
		expectedCode     codes.Code
		expectedContains string
	}{
		{
			name:             "admin authorization error",
			err:              errors.New("ErrUnauthorized: admin access required"),
			operation:        "GetOverviewStats",
			expectedCode:     codes.PermissionDenied,
			expectedContains: "permission denied",
		},
		{
			name:             "unauthorized error",
			err:              errors.New("ErrUnauthorized: user is unauthorized"),
			operation:        "GetListingStats",
			expectedCode:     codes.PermissionDenied,
			expectedContains: "permission denied",
		},
		{
			name:             "permission error",
			err:              errors.New("permission denied for this resource"),
			operation:        "GetListingStats",
			expectedCode:     codes.PermissionDenied,
			expectedContains: "permission denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := mapAnalyticsError(tt.err, tt.operation)

			// Assert
			assert.Error(t, result)
			st, ok := status.FromError(result)
			require.True(t, ok)
			assert.Equal(t, tt.expectedCode, st.Code())
			assert.Contains(t, st.Message(), tt.expectedContains)
		})
	}
}

func TestMapAnalyticsError_ValidationErrors(t *testing.T) {
	tests := []struct {
		name             string
		err              error
		operation        string
		expectedCode     codes.Code
		expectedContains string
	}{
		{
			name:             "invalid input",
			err:              errors.New("invalid: date range too large"),
			operation:        "GetOverviewStats",
			expectedCode:     codes.InvalidArgument,
			expectedContains: "invalid",
		},
		{
			name:             "required field missing",
			err:              errors.New("required: listing_id is required"),
			operation:        "GetListingStats",
			expectedCode:     codes.InvalidArgument,
			expectedContains: "required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := mapAnalyticsError(tt.err, tt.operation)

			// Assert
			assert.Error(t, result)
			st, ok := status.FromError(result)
			require.True(t, ok)
			assert.Equal(t, tt.expectedCode, st.Code())
			assert.Contains(t, st.Message(), tt.expectedContains)
		})
	}
}

func TestMapAnalyticsError_NotFoundErrors(t *testing.T) {
	// Arrange
	err := errors.New("not found: listing 999 does not exist")
	operation := "GetListingStats"

	// Act
	result := mapAnalyticsError(err, operation)

	// Assert
	assert.Error(t, result)
	st, ok := status.FromError(result)
	require.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
}

func TestMapAnalyticsError_InternalErrors(t *testing.T) {
	tests := []struct {
		name             string
		err              error
		operation        string
		expectedCode     codes.Code
		expectedContains string
	}{
		{
			name:             "database connection error",
			err:              errors.New("database connection failed"),
			operation:        "GetOverviewStats",
			expectedCode:     codes.Internal,
			expectedContains: "Internal",
		},
		{
			name:             "service error",
			err:              errors.New("service error: something went wrong"),
			operation:        "GetListingStats",
			expectedCode:     codes.Internal,
			expectedContains: "Internal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := mapAnalyticsError(tt.err, tt.operation)

			// Assert
			assert.Error(t, result)
			st, ok := status.FromError(result)
			require.True(t, ok)
			assert.Equal(t, tt.expectedCode, st.Code())
		})
	}
}

func TestMapAnalyticsError_NilError(t *testing.T) {
	// Act
	result := mapAnalyticsError(nil, "GetOverviewStats")

	// Assert
	assert.NoError(t, result)
	assert.Nil(t, result)
}

// ============================================================================
// Edge Cases and Boundary Conditions
// ============================================================================

func TestValidateOverviewStatsRequest_BoundaryDateRange_ZeroDays(t *testing.T) {
	// Arrange - same day
	now := analyticsTestTimeNow()
	server := &Server{
		logger: zerolog.Nop(),
	}

	req := &listingssvcv1.GetOverviewStatsRequest{
		DateFrom: datePtr(now),
		DateTo:   datePtr(now),
	}

	// Act
	err := server.validateOverviewStatsRequestHelper(req)

	// Assert
	require.NoError(t, err, "Same day date range should be valid")
}

func TestValidateOverviewStatsRequest_FullYearRange(t *testing.T) {
	// Arrange - 365 days exactly
	now := analyticsTestTimeNow()
	yearAgo := now.AddDate(-1, 0, 0)

	server := &Server{
		logger: zerolog.Nop(),
	}

	req := &listingssvcv1.GetOverviewStatsRequest{
		DateFrom: datePtr(yearAgo),
		DateTo:   datePtr(now),
	}

	// Act
	err := server.validateOverviewStatsRequestHelper(req)

	// Assert
	require.NoError(t, err, "One year (365 days) should be valid range")
}

func TestValidateListingStatsRequest_WithVariantOptions(t *testing.T) {
	// Arrange
	now := analyticsTestTimeNow()
	dateFrom := analyticsTestTimeDaysAgo(7)
	productID := int64(501)
	includeVariants := true
	includeGeo := true

	server := &Server{
		logger: zerolog.Nop(),
	}

	req := &listingssvcv1.GetListingStatsRequest{
		Identifier:      &listingssvcv1.GetListingStatsRequest_ProductId{ProductId: productID},
		DateFrom:        datePtr(dateFrom),
		DateTo:          datePtr(now),
		IncludeVariants: &includeVariants,
		IncludeGeo:      &includeGeo,
	}

	// Act
	err := server.validateListingStatsRequestHelper(req)

	// Assert
	require.NoError(t, err)
}

// ============================================================================
// Comprehensive Service Integration Tests
// ============================================================================

func TestMapAnalyticsError_CaseSensitivity(t *testing.T) {
	// Test that error mapping is case-insensitive
	errors_to_test := []struct {
		name        string
		err         error
		expectedErr string
	}{
		{
			name:        "Unauthorized with capital U",
			err:         errors.New("UNAUTHORIZED: admin access required"),
			expectedErr: "permission",
		},
		{
			name:        "Invalid with capital I",
			err:         errors.New("INVALID: bad input"),
			expectedErr: "InvalidArgument",
		},
	}

	for _, tt := range errors_to_test {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := mapAnalyticsError(tt.err, "TestOp")

			// Assert
			assert.Error(t, result)
		})
	}
}

// ============================================================================
// Request Validation Coverage
// ============================================================================

func TestValidateOverviewStatsRequest_EmptyListingType(t *testing.T) {
	// Arrange
	now := analyticsTestTimeNow()
	dateFrom := analyticsTestTimeDaysAgo(7)
	emptyType := ""

	server := &Server{
		logger: zerolog.Nop(),
	}

	req := &listingssvcv1.GetOverviewStatsRequest{
		DateFrom:    datePtr(dateFrom),
		DateTo:      datePtr(now),
		ListingType: &emptyType,
	}

	// Act
	err := server.validateOverviewStatsRequestHelper(req)

	// Assert
	require.NoError(t, err, "Empty listing type should be valid (means all types)")
}

func TestValidateListingStatsRequest_BothIdentifiersZero(t *testing.T) {
	// Arrange
	now := analyticsTestTimeNow()
	dateFrom := analyticsTestTimeDaysAgo(7)

	server := &Server{
		logger: zerolog.Nop(),
	}

	// Create request with zero values for both
	req := &listingssvcv1.GetListingStatsRequest{
		Identifier: &listingssvcv1.GetListingStatsRequest_ListingId{ListingId: 0},
		DateFrom:   datePtr(dateFrom),
		DateTo:     datePtr(now),
	}

	// Act
	err := server.validateListingStatsRequestHelper(req)

	// Assert
	require.Error(t, err, "Should fail when listing_id is 0")
}

// ============================================================================
// Test Service Methods (Mocked)
// ============================================================================

func TestGetOverviewStats_SuccessWithMockedService(t *testing.T) {
	// Arrange
	now := analyticsTestTimeNow()
	dateFrom := analyticsTestTimeDaysAgo(7)

	mockSvc := &mockAnalyticsService{
		getOverviewStatsFunc: func(ctx context.Context, req *listingssvcv1.GetOverviewStatsRequest, userID int64, isAdmin bool) (*listingssvcv1.GetOverviewStatsResponse, error) {
			assert.True(t, isAdmin, "Should be admin")
			return &listingssvcv1.GetOverviewStatsResponse{
				Listings: &listingssvcv1.ListingsStats{
					TotalListings:  int32(500),
					ActiveListings: int32(480),
				},
				Revenue: &listingssvcv1.RevenueStats{
					TotalRevenue:  15000.00,
					AvgOrderValue: 150.00,
					Transactions:  int32(100),
				},
				Orders: &listingssvcv1.OrdersStats{
					TotalOrders:     int32(100),
					CompletedOrders: int32(95),
				},
				GeneratedAt: datePtr(now),
				DataFrom:    datePtr(dateFrom),
				DataTo:      datePtr(now),
			}, nil
		},
	}

	// Act
	resp, err := mockSvc.GetOverviewStats(context.Background(), &listingssvcv1.GetOverviewStatsRequest{
		DateFrom: datePtr(dateFrom),
		DateTo:   datePtr(now),
	}, 123, true)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int32(500), resp.Listings.TotalListings)
	assert.Equal(t, 15000.00, resp.Revenue.TotalRevenue)
}

func TestGetListingStats_SuccessWithMockedService(t *testing.T) {
	// Arrange
	now := analyticsTestTimeNow()
	dateFrom := analyticsTestTimeDaysAgo(7)
	listingID := int64(281)

	mockSvc := &mockAnalyticsService{
		getListingStatsFunc: func(ctx context.Context, req *listingssvcv1.GetListingStatsRequest, userID int64, isAdmin bool) (*listingssvcv1.GetListingStatsResponse, error) {
			return &listingssvcv1.GetListingStatsResponse{
				ListingId:      listingID,
				ListingName:    "Test Product",
				ListingType:    "b2c",
				TotalViews:     int32(5000),
				FavoriteCount:  int32(120),
				TotalSales:     int32(45),
				TotalRevenue:   45000.00,
				AvgOrderValue:  1000.00,
				ConversionRate: 0.009,
				GeneratedAt:    datePtr(now),
				DataFrom:       datePtr(dateFrom),
				DataTo:         datePtr(now),
			}, nil
		},
	}

	// Act
	resp, err := mockSvc.GetListingStats(context.Background(), &listingssvcv1.GetListingStatsRequest{
		Identifier: &listingssvcv1.GetListingStatsRequest_ListingId{ListingId: listingID},
		DateFrom:   datePtr(dateFrom),
		DateTo:     datePtr(now),
	}, 123, false)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, listingID, resp.ListingId)
	assert.Equal(t, "Test Product", resp.ListingName)
	assert.Equal(t, int32(5000), resp.TotalViews)
}

func TestGetOverviewStats_UnauthorizedWithMockedService(t *testing.T) {
	// Arrange
	now := analyticsTestTimeNow()
	dateFrom := analyticsTestTimeDaysAgo(7)

	mockSvc := &mockAnalyticsService{
		getOverviewStatsFunc: func(ctx context.Context, req *listingssvcv1.GetOverviewStatsRequest, userID int64, isAdmin bool) (*listingssvcv1.GetOverviewStatsResponse, error) {
			return nil, errors.New("ErrUnauthorized: admin access required")
		},
	}

	// Act
	resp, err := mockSvc.GetOverviewStats(context.Background(), &listingssvcv1.GetOverviewStatsRequest{
		DateFrom: datePtr(dateFrom),
		DateTo:   datePtr(now),
	}, 456, false)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "admin")
}

func TestGetListingStats_UnauthorizedWithMockedService(t *testing.T) {
	// Arrange
	now := analyticsTestTimeNow()
	dateFrom := analyticsTestTimeDaysAgo(7)
	listingID := int64(281)

	mockSvc := &mockAnalyticsService{
		getListingStatsFunc: func(ctx context.Context, req *listingssvcv1.GetListingStatsRequest, userID int64, isAdmin bool) (*listingssvcv1.GetListingStatsResponse, error) {
			return nil, errors.New("ErrUnauthorized: user does not own this listing")
		},
	}

	// Act
	resp, err := mockSvc.GetListingStats(context.Background(), &listingssvcv1.GetListingStatsRequest{
		Identifier: &listingssvcv1.GetListingStatsRequest_ListingId{ListingId: listingID},
		DateFrom:   datePtr(dateFrom),
		DateTo:     datePtr(now),
	}, 456, false)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "user")
}

func TestGetOverviewStats_ServiceErrorWithMockedService(t *testing.T) {
	// Arrange
	now := analyticsTestTimeNow()
	dateFrom := analyticsTestTimeDaysAgo(7)

	mockSvc := &mockAnalyticsService{
		getOverviewStatsFunc: func(ctx context.Context, req *listingssvcv1.GetOverviewStatsRequest, userID int64, isAdmin bool) (*listingssvcv1.GetOverviewStatsResponse, error) {
			return nil, errors.New("database connection failed")
		},
	}

	// Act
	resp, err := mockSvc.GetOverviewStats(context.Background(), &listingssvcv1.GetOverviewStatsRequest{
		DateFrom: datePtr(dateFrom),
		DateTo:   datePtr(now),
	}, 123, true)

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
}

// ============================================================================
// GetTrendingStats Handler Tests
// ============================================================================

func TestGetTrendingStats_Success(t *testing.T) {
	// Arrange
	now := analyticsTestTimeNow()
	mockSvc := &mockAnalyticsService{
		getTrendingStatsFunc: func(ctx context.Context, req *listingssvcv1.GetTrendingStatsRequest) (*listingssvcv1.GetTrendingStatsResponse, error) {
			return &listingssvcv1.GetTrendingStatsResponse{
				TrendingCategories: []*listingssvcv1.TrendingCategory{
					{
						CategoryId:     1001,
						CategoryName:   "Electronics",
						OrderCount_30D: 450,
						OrderCount_7D:  150,
						GrowthRate:     42.8,
						TrendScore:     1950.0,
					},
					{
						CategoryId:     1002,
						CategoryName:   "Fashion",
						OrderCount_30D: 320,
						OrderCount_7D:  110,
						GrowthRate:     35.2,
						TrendScore:     1420.0,
					},
				},
				HotListings: []*listingssvcv1.HotListing{
					{
						ListingId:        281,
						Title:            "iPhone 15 Pro",
						Orders_24H:       25,
						Orders_7D:        80,
						OrdersGrowth:     2.19,
						QuantitySold_24H: 30,
						Price:            1299.99,
					},
				},
				PopularSearches: []*listingssvcv1.PopularSearch{
					{
						Query:       "iphone",
						SearchCount: 1250,
					},
					{
						Query:       "laptop",
						SearchCount: 890,
					},
				},
				GeneratedAt: datePtr(now),
			}, nil
		},
	}

	// Act
	resp, err := mockSvc.GetTrendingStats(context.Background(), &listingssvcv1.GetTrendingStatsRequest{})

	// Assert
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Len(t, resp.TrendingCategories, 2)
	assert.Len(t, resp.HotListings, 1)
	assert.Len(t, resp.PopularSearches, 2)
	assert.Equal(t, "Electronics", resp.TrendingCategories[0].CategoryName)
	assert.Equal(t, "iPhone 15 Pro", resp.HotListings[0].Title)
	assert.Equal(t, "iphone", resp.PopularSearches[0].Query)
}

func TestGetTrendingStats_EmptyResults(t *testing.T) {
	// Arrange
	now := analyticsTestTimeNow()
	mockSvc := &mockAnalyticsService{
		getTrendingStatsFunc: func(ctx context.Context, req *listingssvcv1.GetTrendingStatsRequest) (*listingssvcv1.GetTrendingStatsResponse, error) {
			return &listingssvcv1.GetTrendingStatsResponse{
				TrendingCategories: []*listingssvcv1.TrendingCategory{},
				HotListings:        []*listingssvcv1.HotListing{},
				PopularSearches:    []*listingssvcv1.PopularSearch{},
				GeneratedAt:        datePtr(now),
			}, nil
		},
	}

	// Act
	resp, err := mockSvc.GetTrendingStats(context.Background(), &listingssvcv1.GetTrendingStatsRequest{})

	// Assert
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Empty(t, resp.TrendingCategories)
	assert.Empty(t, resp.HotListings)
	assert.Empty(t, resp.PopularSearches)
}

func TestGetTrendingStats_ServiceError(t *testing.T) {
	// Arrange
	mockSvc := &mockAnalyticsService{
		getTrendingStatsFunc: func(ctx context.Context, req *listingssvcv1.GetTrendingStatsRequest) (*listingssvcv1.GetTrendingStatsResponse, error) {
			return nil, errors.New("database connection failed")
		},
	}

	// Act
	resp, err := mockSvc.GetTrendingStats(context.Background(), &listingssvcv1.GetTrendingStatsRequest{})

	// Assert
	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "database connection failed")
}

func TestGetTrendingStats_CacheEmpty(t *testing.T) {
	// Arrange
	now := analyticsTestTimeNow()
	mockSvc := &mockAnalyticsService{
		getTrendingStatsFunc: func(ctx context.Context, req *listingssvcv1.GetTrendingStatsRequest) (*listingssvcv1.GetTrendingStatsResponse, error) {
			// Simulating empty cache - returning empty data but no error
			return &listingssvcv1.GetTrendingStatsResponse{
				TrendingCategories: []*listingssvcv1.TrendingCategory{},
				HotListings:        []*listingssvcv1.HotListing{},
				PopularSearches:    []*listingssvcv1.PopularSearch{},
				GeneratedAt:        datePtr(now),
			}, nil
		},
	}

	// Act
	resp, err := mockSvc.GetTrendingStats(context.Background(), &listingssvcv1.GetTrendingStatsRequest{})

	// Assert
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotNil(t, resp.GeneratedAt)
}

func TestGetTrendingStats_TrendingCategories_Ordering(t *testing.T) {
	// Arrange
	mockSvc := &mockAnalyticsService{
		getTrendingStatsFunc: func(ctx context.Context, req *listingssvcv1.GetTrendingStatsRequest) (*listingssvcv1.GetTrendingStatsResponse, error) {
			return &listingssvcv1.GetTrendingStatsResponse{
				TrendingCategories: []*listingssvcv1.TrendingCategory{
					{CategoryId: 1001, CategoryName: "Electronics", TrendScore: 1950.0},
					{CategoryId: 1002, CategoryName: "Fashion", TrendScore: 1420.0},
					{CategoryId: 1003, CategoryName: "Home", TrendScore: 1100.0},
				},
				HotListings:     []*listingssvcv1.HotListing{},
				PopularSearches: []*listingssvcv1.PopularSearch{},
				GeneratedAt:     datePtr(analyticsTestTimeNow()),
			}, nil
		},
	}

	// Act
	resp, err := mockSvc.GetTrendingStats(context.Background(), &listingssvcv1.GetTrendingStatsRequest{})

	// Assert
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.TrendingCategories, 3)

	// Verify descending order by trend_score
	assert.Greater(t, resp.TrendingCategories[0].TrendScore, resp.TrendingCategories[1].TrendScore)
	assert.Greater(t, resp.TrendingCategories[1].TrendScore, resp.TrendingCategories[2].TrendScore)
}

func TestGetTrendingStats_HotListings_GrowthCalculation(t *testing.T) {
	// Arrange
	mockSvc := &mockAnalyticsService{
		getTrendingStatsFunc: func(ctx context.Context, req *listingssvcv1.GetTrendingStatsRequest) (*listingssvcv1.GetTrendingStatsResponse, error) {
			return &listingssvcv1.GetTrendingStatsResponse{
				TrendingCategories: []*listingssvcv1.TrendingCategory{},
				HotListings: []*listingssvcv1.HotListing{
					{
						ListingId:    281,
						Title:        "Viral Product",
						Orders_24H:   50,
						Orders_7D:    80,
						OrdersGrowth: 4.375, // (50 / (80/7)) = 4.375x growth
					},
					{
						ListingId:    282,
						Title:        "Steady Product",
						Orders_24H:   10,
						Orders_7D:    70,
						OrdersGrowth: 1.0, // (10 / (70/7)) = 1.0x (same rate)
					},
				},
				PopularSearches: []*listingssvcv1.PopularSearch{},
				GeneratedAt:     datePtr(analyticsTestTimeNow()),
			}, nil
		},
	}

	// Act
	resp, err := mockSvc.GetTrendingStats(context.Background(), &listingssvcv1.GetTrendingStatsRequest{})

	// Assert
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.HotListings, 2)

	// Verify growth rates
	assert.Greater(t, resp.HotListings[0].OrdersGrowth, float64(4.0))     // Viral spike
	assert.LessOrEqual(t, resp.HotListings[1].OrdersGrowth, float64(1.1)) // Steady rate
}

func TestGetTrendingStats_PopularSearches_Sorted(t *testing.T) {
	// Arrange
	mockSvc := &mockAnalyticsService{
		getTrendingStatsFunc: func(ctx context.Context, req *listingssvcv1.GetTrendingStatsRequest) (*listingssvcv1.GetTrendingStatsResponse, error) {
			return &listingssvcv1.GetTrendingStatsResponse{
				TrendingCategories: []*listingssvcv1.TrendingCategory{},
				HotListings:        []*listingssvcv1.HotListing{},
				PopularSearches: []*listingssvcv1.PopularSearch{
					{Query: "iphone", SearchCount: 1250},
					{Query: "laptop", SearchCount: 890},
					{Query: "headphones", SearchCount: 670},
				},
				GeneratedAt: datePtr(analyticsTestTimeNow()),
			}, nil
		},
	}

	// Act
	resp, err := mockSvc.GetTrendingStats(context.Background(), &listingssvcv1.GetTrendingStatsRequest{})

	// Assert
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.PopularSearches, 3)

	// Verify descending order by search_count
	assert.Greater(t, resp.PopularSearches[0].SearchCount, resp.PopularSearches[1].SearchCount)
	assert.Greater(t, resp.PopularSearches[1].SearchCount, resp.PopularSearches[2].SearchCount)
}

func TestGetTrendingStats_GeneratedAtTimestamp(t *testing.T) {
	// Arrange
	expectedTime := analyticsTestTimeNow()
	mockSvc := &mockAnalyticsService{
		getTrendingStatsFunc: func(ctx context.Context, req *listingssvcv1.GetTrendingStatsRequest) (*listingssvcv1.GetTrendingStatsResponse, error) {
			return &listingssvcv1.GetTrendingStatsResponse{
				TrendingCategories: []*listingssvcv1.TrendingCategory{},
				HotListings:        []*listingssvcv1.HotListing{},
				PopularSearches:    []*listingssvcv1.PopularSearch{},
				GeneratedAt:        datePtr(expectedTime),
			}, nil
		},
	}

	// Act
	resp, err := mockSvc.GetTrendingStats(context.Background(), &listingssvcv1.GetTrendingStatsRequest{})

	// Assert
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.GeneratedAt)
	assert.Equal(t, expectedTime.Unix(), resp.GeneratedAt.AsTime().Unix())
}
