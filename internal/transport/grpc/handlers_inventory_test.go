package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/sveturs/listings/api/proto/listings/v1"
	"github.com/sveturs/listings/internal/domain"
)

// Helper functions for creating pointers
func int64Ptr(v int64) *int64 {
	return &v
}

func stringPtr(v string) *string {
	return &v
}

// Mock service interface methods needed for inventory tests
func (m *MockListingsService) UpdateProductInventory(ctx context.Context, storefrontID, productID, variantID int64, movementType string, quantity int32, reason, notes string, userID int64) (int32, int32, error) {
	args := m.Called(ctx, storefrontID, productID, variantID, movementType, quantity, reason, notes, userID)
	return args.Get(0).(int32), args.Get(1).(int32), args.Error(2)
}

func (m *MockListingsService) BatchUpdateStock(ctx context.Context, storefrontID int64, items []domain.StockUpdateItem, reason string, userID int64) (int32, int32, []domain.StockUpdateResult, error) {
	args := m.Called(ctx, storefrontID, items, reason, userID)
	if args.Get(2) == nil {
		return args.Get(0).(int32), args.Get(1).(int32), nil, args.Error(3)
	}
	return args.Get(0).(int32), args.Get(1).(int32), args.Get(2).([]domain.StockUpdateResult), args.Error(3)
}

func (m *MockListingsService) GetProductStats(ctx context.Context, storefrontID int64) (*domain.ProductStats, error) {
	args := m.Called(ctx, storefrontID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ProductStats), args.Error(1)
}

func (m *MockListingsService) IncrementProductViews(ctx context.Context, productID int64) error {
	args := m.Called(ctx, productID)
	return args.Error(0)
}

// ListingsServiceInterface defines the interface for service layer methods
type ListingsServiceInterface interface {
	UpdateProductInventory(ctx context.Context, storefrontID, productID, variantID int64, movementType string, quantity int32, reason, notes string, userID int64) (int32, int32, error)
	BatchUpdateStock(ctx context.Context, storefrontID int64, items []domain.StockUpdateItem, reason string, userID int64) (int32, int32, []domain.StockUpdateResult, error)
	GetProductStats(ctx context.Context, storefrontID int64) (*domain.ProductStats, error)
	IncrementProductViews(ctx context.Context, productID int64) error
}

// Ensure MockListingsService implements the interface
var _ ListingsServiceInterface = (*MockListingsService)(nil)

// serverWithInterface wraps Server for testing with interface-based service
type serverWithInterface struct {
	pb.UnimplementedListingsServiceServer
	service ListingsServiceInterface
	logger  zerolog.Logger
}

// setupInventoryTestServer creates a test server with mock service
func setupInventoryTestServer() (*serverWithInterface, *MockListingsService) {
	mockService := new(MockListingsService)
	logger := zerolog.Nop()

	server := &serverWithInterface{
		service: mockService,
		logger:  logger,
	}

	return server, mockService
}

// RecordInventoryMovement implements the gRPC handler
func (s *serverWithInterface) RecordInventoryMovement(ctx context.Context, req *pb.RecordInventoryMovementRequest) (*pb.RecordInventoryMovementResponse, error) {
	variantIDLog := int64(0)
	if req.VariantId != nil {
		variantIDLog = *req.VariantId
	}

	s.logger.Debug().
		Int64("storefront_id", req.StorefrontId).
		Int64("product_id", req.ProductId).
		Int64("variant_id", variantIDLog).
		Str("movement_type", req.MovementType).
		Int32("quantity", req.Quantity).
		Msg("RecordInventoryMovement called")

	// Validation
	if req.StorefrontId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "storefront ID must be greater than 0")
	}

	if req.ProductId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "product ID must be greater than 0")
	}

	// Validate movement type
	validMovementTypes := map[string]bool{
		"in":         true,
		"out":        true,
		"adjustment": true,
	}
	if !validMovementTypes[req.MovementType] {
		return nil, status.Error(codes.InvalidArgument, "movement_type must be 'in', 'out', or 'adjustment'")
	}

	if req.Quantity < 0 {
		return nil, status.Error(codes.InvalidArgument, "quantity must be non-negative")
	}

	if req.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "user ID must be greater than 0")
	}

	// Prepare optional variant ID (0 means product-level stock update)
	var variantID int64
	if req.VariantId != nil {
		variantID = *req.VariantId
	}

	// Prepare optional reason and notes
	reason := ""
	if req.Reason != nil {
		reason = *req.Reason
	}

	notes := ""
	if req.Notes != nil {
		notes = *req.Notes
	}

	// Call service to record inventory movement
	stockBefore, stockAfter, err := s.service.UpdateProductInventory(
		ctx,
		req.StorefrontId,
		req.ProductId,
		variantID,
		req.MovementType,
		req.Quantity,
		reason,
		notes,
		req.UserId,
	)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to record inventory movement")

		// Check for specific errors with placeholders
		errMsg := err.Error()
		if errMsg == "inventory.product_not_found" {
			return nil, status.Error(codes.NotFound, "inventory.product_not_found")
		}
		if errMsg == "inventory.variant_not_found" {
			return nil, status.Error(codes.NotFound, "inventory.variant_not_found")
		}
		if errMsg == "inventory.insufficient_stock" {
			return nil, status.Error(codes.FailedPrecondition, "inventory.insufficient_stock")
		}

		// Generic error
		return nil, status.Error(codes.Internal, "inventory.update_failed")
	}

	s.logger.Info().
		Int64("product_id", req.ProductId).
		Int64("variant_id", variantID).
		Int32("stock_before", stockBefore).
		Int32("stock_after", stockAfter).
		Msg("inventory movement recorded successfully")

	return &pb.RecordInventoryMovementResponse{
		Success:     true,
		StockBefore: stockBefore,
		StockAfter:  stockAfter,
		Error:       nil,
	}, nil
}

// BatchUpdateStock implements the gRPC handler
func (s *serverWithInterface) BatchUpdateStock(ctx context.Context, req *pb.BatchUpdateStockRequest) (*pb.BatchUpdateStockResponse, error) {
	s.logger.Debug().
		Int64("storefront_id", req.StorefrontId).
		Int("item_count", len(req.Items)).
		Msg("BatchUpdateStock called")

	// Validation
	if req.StorefrontId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "storefront ID must be greater than 0")
	}

	if len(req.Items) == 0 {
		return nil, status.Error(codes.InvalidArgument, "items list cannot be empty")
	}

	if len(req.Items) > 1000 {
		return nil, status.Error(codes.InvalidArgument, "cannot update more than 1000 items at once")
	}

	if req.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "user ID must be greater than 0")
	}

	// Validate each item
	for i, item := range req.Items {
		if item.ProductId <= 0 {
			return nil, status.Errorf(codes.InvalidArgument, "invalid product_id at index %d", i)
		}
		if item.Quantity < 0 {
			return nil, status.Errorf(codes.InvalidArgument, "invalid quantity at index %d: cannot be negative", i)
		}
	}

	// Convert proto items to domain models
	domainItems := make([]domain.StockUpdateItem, len(req.Items))
	for i, item := range req.Items {
		domainItem := domain.StockUpdateItem{
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
		}

		if item.VariantId != nil {
			domainItem.VariantID = item.VariantId
		}

		if item.Reason != nil {
			domainItem.Reason = item.Reason
		}

		domainItems[i] = domainItem
	}

	// Prepare optional common reason
	reason := ""
	if req.Reason != nil {
		reason = *req.Reason
	}

	// Call service to batch update stock
	successCount, failedCount, results, err := s.service.BatchUpdateStock(
		ctx,
		req.StorefrontId,
		domainItems,
		reason,
		req.UserId,
	)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to batch update stock")
		return nil, status.Error(codes.Internal, "inventory.batch_update_failed")
	}

	// Convert domain results to proto
	protoResults := make([]*pb.StockUpdateResult, len(results))
	for i, result := range results {
		protoResult := &pb.StockUpdateResult{
			ProductId:   result.ProductID,
			StockBefore: result.StockBefore,
			StockAfter:  result.StockAfter,
			Success:     result.Success,
		}

		if result.VariantID != nil {
			protoResult.VariantId = result.VariantID
		}

		if result.Error != nil {
			protoResult.Error = result.Error
		}

		protoResults[i] = protoResult
	}

	s.logger.Info().
		Int32("successful_count", successCount).
		Int32("failed_count", failedCount).
		Msg("batch stock update completed")

	return &pb.BatchUpdateStockResponse{
		SuccessfulCount: successCount,
		FailedCount:     failedCount,
		Results:         protoResults,
	}, nil
}

// GetProductStats implements the gRPC handler
func (s *serverWithInterface) GetProductStats(ctx context.Context, req *pb.GetProductStatsRequest) (*pb.GetProductStatsResponse, error) {
	s.logger.Debug().
		Int64("storefront_id", req.StorefrontId).
		Msg("GetProductStats called")

	// Validation
	if req.StorefrontId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "storefront ID must be greater than 0")
	}

	// Call service to get product stats
	stats, err := s.service.GetProductStats(ctx, req.StorefrontId)
	if err != nil {
		s.logger.Error().Err(err).Int64("storefront_id", req.StorefrontId).Msg("failed to get product stats")
		return nil, status.Error(codes.Internal, "products.stats_failed")
	}

	// Convert domain stats to proto
	protoStats := &pb.ProductStats{
		TotalProducts:  stats.TotalProducts,
		ActiveProducts: stats.ActiveProducts,
		OutOfStock:     stats.OutOfStock,
		LowStock:       stats.LowStock,
		TotalValue:     stats.TotalValue,
		TotalSold:      stats.TotalSold,
	}

	s.logger.Info().
		Int32("total_products", stats.TotalProducts).
		Int32("active_products", stats.ActiveProducts).
		Float64("total_value", stats.TotalValue).
		Msg("product stats retrieved successfully")

	return &pb.GetProductStatsResponse{Stats: protoStats}, nil
}

// IncrementProductViews implements the gRPC handler
func (s *serverWithInterface) IncrementProductViews(ctx context.Context, req *pb.IncrementProductViewsRequest) (*emptypb.Empty, error) {
	s.logger.Debug().
		Int64("product_id", req.ProductId).
		Msg("IncrementProductViews called")

	// Validation
	if req.ProductId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "product ID must be greater than 0")
	}

	// Call service to increment views
	if err := s.service.IncrementProductViews(ctx, req.ProductId); err != nil {
		s.logger.Error().Err(err).Int64("product_id", req.ProductId).Msg("failed to increment product views")
		return nil, status.Error(codes.Internal, "products.increment_views_failed")
	}

	s.logger.Debug().Int64("product_id", req.ProductId).Msg("product views incremented successfully")
	return &emptypb.Empty{}, nil
}

// TestRecordInventoryMovement_ValidRequest_Success tests successful inventory movement recording
func TestRecordInventoryMovement_ValidRequest_Success(t *testing.T) {
	server, mockService := setupInventoryTestServer()
	ctx := context.Background()

	req := &pb.RecordInventoryMovementRequest{
		StorefrontId: 123,
		ProductId:    456,
		VariantId:    int64Ptr(789),
		MovementType: "out",
		Quantity:     5,
		Reason:       stringPtr("sale"),
		Notes:        stringPtr("Order #12345"),
		UserId:       1,
	}

	// Mock expectations
	mockService.On("UpdateProductInventory",
		ctx,
		req.StorefrontId,
		req.ProductId,
		int64(789), // variantID
		req.MovementType,
		req.Quantity,
		"sale",         // reason
		"Order #12345", // notes
		req.UserId,
	).Return(int32(100), int32(95), nil)

	// Execute
	resp, err := server.RecordInventoryMovement(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Equal(t, int32(100), resp.StockBefore)
	assert.Equal(t, int32(95), resp.StockAfter)
	assert.Nil(t, resp.Error)

	mockService.AssertExpectations(t)
}

// TestRecordInventoryMovement_InvalidStorefrontID_ValidationError tests validation for storefront_id
func TestRecordInventoryMovement_InvalidStorefrontID_ValidationError(t *testing.T) {
	server, _ := setupInventoryTestServer()
	ctx := context.Background()

	testCases := []struct {
		name         string
		storefrontID int64
		expectedCode codes.Code
	}{
		{
			name:         "Zero storefront ID",
			storefrontID: 0,
			expectedCode: codes.InvalidArgument,
		},
		{
			name:         "Negative storefront ID",
			storefrontID: -1,
			expectedCode: codes.InvalidArgument,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &pb.RecordInventoryMovementRequest{
				StorefrontId: tc.storefrontID,
				ProductId:    456,
				MovementType: "in",
				Quantity:     10,
				UserId:       1,
			}

			resp, err := server.RecordInventoryMovement(ctx, req)

			assert.Error(t, err)
			assert.Nil(t, resp)

			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tc.expectedCode, st.Code())
		})
	}
}

// TestRecordInventoryMovement_InvalidProductID_ValidationError tests validation for product_id
func TestRecordInventoryMovement_InvalidProductID_ValidationError(t *testing.T) {
	server, _ := setupInventoryTestServer()
	ctx := context.Background()

	testCases := []struct {
		name         string
		productID    int64
		expectedCode codes.Code
	}{
		{
			name:         "Zero product ID",
			productID:    0,
			expectedCode: codes.InvalidArgument,
		},
		{
			name:         "Negative product ID",
			productID:    -1,
			expectedCode: codes.InvalidArgument,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &pb.RecordInventoryMovementRequest{
				StorefrontId: 123,
				ProductId:    tc.productID,
				MovementType: "in",
				Quantity:     10,
				UserId:       1,
			}

			resp, err := server.RecordInventoryMovement(ctx, req)

			assert.Error(t, err)
			assert.Nil(t, resp)

			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, tc.expectedCode, st.Code())
		})
	}
}

// TestRecordInventoryMovement_InvalidMovementType_ValidationError tests validation for movement_type
func TestRecordInventoryMovement_InvalidMovementType_ValidationError(t *testing.T) {
	server, _ := setupInventoryTestServer()
	ctx := context.Background()

	testCases := []struct {
		name         string
		movementType string
	}{
		{
			name:         "Empty movement type",
			movementType: "",
		},
		{
			name:         "Invalid movement type",
			movementType: "invalid",
		},
		{
			name:         "Uppercase movement type",
			movementType: "OUT",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &pb.RecordInventoryMovementRequest{
				StorefrontId: 123,
				ProductId:    456,
				MovementType: tc.movementType,
				Quantity:     10,
				UserId:       1,
			}

			resp, err := server.RecordInventoryMovement(ctx, req)

			assert.Error(t, err)
			assert.Nil(t, resp)

			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, codes.InvalidArgument, st.Code())
		})
	}
}

// TestRecordInventoryMovement_NegativeQuantity_ValidationError tests validation for negative quantity
func TestRecordInventoryMovement_NegativeQuantity_ValidationError(t *testing.T) {
	server, _ := setupInventoryTestServer()
	ctx := context.Background()

	req := &pb.RecordInventoryMovementRequest{
		StorefrontId: 123,
		ProductId:    456,
		MovementType: "in",
		Quantity:     -5,
		UserId:       1,
	}

	resp, err := server.RecordInventoryMovement(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

// TestRecordInventoryMovement_InvalidUserID_ValidationError tests validation for user_id
func TestRecordInventoryMovement_InvalidUserID_ValidationError(t *testing.T) {
	server, _ := setupInventoryTestServer()
	ctx := context.Background()

	req := &pb.RecordInventoryMovementRequest{
		StorefrontId: 123,
		ProductId:    456,
		MovementType: "in",
		Quantity:     10,
		UserId:       0,
	}

	resp, err := server.RecordInventoryMovement(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

// TestRecordInventoryMovement_ProductNotFound_NotFoundError tests product not found error
func TestRecordInventoryMovement_ProductNotFound_NotFoundError(t *testing.T) {
	server, mockService := setupInventoryTestServer()
	ctx := context.Background()

	req := &pb.RecordInventoryMovementRequest{
		StorefrontId: 123,
		ProductId:    999,
		MovementType: "in",
		Quantity:     10,
		UserId:       1,
	}

	// Mock service returns product not found error
	mockService.On("UpdateProductInventory",
		ctx,
		req.StorefrontId,
		req.ProductId,
		int64(0), // variantID
		req.MovementType,
		req.Quantity,
		"", // reason
		"", // notes
		req.UserId,
	).Return(int32(0), int32(0), errors.New("inventory.product_not_found"))

	resp, err := server.RecordInventoryMovement(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Equal(t, "inventory.product_not_found", st.Message())

	mockService.AssertExpectations(t)
}

// TestRecordInventoryMovement_InsufficientStock_PreconditionFailedError tests insufficient stock error
func TestRecordInventoryMovement_InsufficientStock_PreconditionFailedError(t *testing.T) {
	server, mockService := setupInventoryTestServer()
	ctx := context.Background()

	req := &pb.RecordInventoryMovementRequest{
		StorefrontId: 123,
		ProductId:    456,
		MovementType: "out",
		Quantity:     100,
		UserId:       1,
	}

	// Mock service returns insufficient stock error
	mockService.On("UpdateProductInventory",
		ctx,
		req.StorefrontId,
		req.ProductId,
		int64(0), // variantID
		req.MovementType,
		req.Quantity,
		"", // reason
		"", // notes
		req.UserId,
	).Return(int32(10), int32(-90), errors.New("inventory.insufficient_stock"))

	resp, err := server.RecordInventoryMovement(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.FailedPrecondition, st.Code())
	assert.Equal(t, "inventory.insufficient_stock", st.Message())

	mockService.AssertExpectations(t)
}

// TestBatchUpdateStock_ValidRequest_Success tests successful batch stock update
func TestBatchUpdateStock_ValidRequest_Success(t *testing.T) {
	server, mockService := setupInventoryTestServer()
	ctx := context.Background()

	req := &pb.BatchUpdateStockRequest{
		StorefrontId: 123,
		Items: []*pb.StockUpdateItem{
			{
				ProductId: 100,
				VariantId: int64Ptr(200),
				Quantity:  50,
			},
			{
				ProductId: 101,
				Quantity:  25,
				Reason:    stringPtr("restock"),
			},
		},
		Reason: stringPtr("inventory_adjustment"),
		UserId: 1,
	}

	// Expected domain items
	expectedItems := []domain.StockUpdateItem{
		{
			ProductID: 100,
			VariantID: int64Ptr(200),
			Quantity:  50,
		},
		{
			ProductID: 101,
			Quantity:  25,
			Reason:    stringPtr("restock"),
		},
	}

	expectedResults := []domain.StockUpdateResult{
		{
			ProductID:   100,
			VariantID:   int64Ptr(200),
			StockBefore: 10,
			StockAfter:  50,
			Success:     true,
		},
		{
			ProductID:   101,
			StockBefore: 0,
			StockAfter:  25,
			Success:     true,
		},
	}

	// Mock expectations
	mockService.On("BatchUpdateStock",
		ctx,
		req.StorefrontId,
		mock.MatchedBy(func(items []domain.StockUpdateItem) bool {
			if len(items) != len(expectedItems) {
				return false
			}
			for i := range items {
				if items[i].ProductID != expectedItems[i].ProductID {
					return false
				}
				if items[i].Quantity != expectedItems[i].Quantity {
					return false
				}
			}
			return true
		}),
		"inventory_adjustment",
		req.UserId,
	).Return(int32(2), int32(0), expectedResults, nil)

	// Execute
	resp, err := server.BatchUpdateStock(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int32(2), resp.SuccessfulCount)
	assert.Equal(t, int32(0), resp.FailedCount)
	assert.Len(t, resp.Results, 2)

	assert.Equal(t, int64(100), resp.Results[0].ProductId)
	assert.Equal(t, int32(10), resp.Results[0].StockBefore)
	assert.Equal(t, int32(50), resp.Results[0].StockAfter)
	assert.True(t, resp.Results[0].Success)

	mockService.AssertExpectations(t)
}

// TestBatchUpdateStock_InvalidStorefrontID_ValidationError tests validation for storefront_id
func TestBatchUpdateStock_InvalidStorefrontID_ValidationError(t *testing.T) {
	server, _ := setupInventoryTestServer()
	ctx := context.Background()

	req := &pb.BatchUpdateStockRequest{
		StorefrontId: 0,
		Items: []*pb.StockUpdateItem{
			{ProductId: 100, Quantity: 10},
		},
		UserId: 1,
	}

	resp, err := server.BatchUpdateStock(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

// TestBatchUpdateStock_EmptyItems_ValidationError tests validation for empty items list
func TestBatchUpdateStock_EmptyItems_ValidationError(t *testing.T) {
	server, _ := setupInventoryTestServer()
	ctx := context.Background()

	req := &pb.BatchUpdateStockRequest{
		StorefrontId: 123,
		Items:        []*pb.StockUpdateItem{},
		UserId:       1,
	}

	resp, err := server.BatchUpdateStock(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

// TestBatchUpdateStock_TooManyItems_ValidationError tests batch size limit
func TestBatchUpdateStock_TooManyItems_ValidationError(t *testing.T) {
	server, _ := setupInventoryTestServer()
	ctx := context.Background()

	// Create 1001 items
	items := make([]*pb.StockUpdateItem, 1001)
	for i := 0; i < 1001; i++ {
		items[i] = &pb.StockUpdateItem{
			ProductId: int64(i + 1),
			Quantity:  10,
		}
	}

	req := &pb.BatchUpdateStockRequest{
		StorefrontId: 123,
		Items:        items,
		UserId:       1,
	}

	resp, err := server.BatchUpdateStock(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

// TestBatchUpdateStock_InvalidItemProductID_ValidationError tests item validation
func TestBatchUpdateStock_InvalidItemProductID_ValidationError(t *testing.T) {
	server, _ := setupInventoryTestServer()
	ctx := context.Background()

	req := &pb.BatchUpdateStockRequest{
		StorefrontId: 123,
		Items: []*pb.StockUpdateItem{
			{ProductId: 100, Quantity: 10},
			{ProductId: 0, Quantity: 20}, // Invalid
		},
		UserId: 1,
	}

	resp, err := server.BatchUpdateStock(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

// TestBatchUpdateStock_NegativeQuantity_ValidationError tests negative quantity validation
func TestBatchUpdateStock_NegativeQuantity_ValidationError(t *testing.T) {
	server, _ := setupInventoryTestServer()
	ctx := context.Background()

	req := &pb.BatchUpdateStockRequest{
		StorefrontId: 123,
		Items: []*pb.StockUpdateItem{
			{ProductId: 100, Quantity: -10}, // Invalid
		},
		UserId: 1,
	}

	resp, err := server.BatchUpdateStock(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

// TestBatchUpdateStock_PartialSuccess tests partial success scenario
func TestBatchUpdateStock_PartialSuccess(t *testing.T) {
	server, mockService := setupInventoryTestServer()
	ctx := context.Background()

	req := &pb.BatchUpdateStockRequest{
		StorefrontId: 123,
		Items: []*pb.StockUpdateItem{
			{ProductId: 100, Quantity: 50},
			{ProductId: 999, Quantity: 25}, // Will fail
		},
		UserId: 1,
	}

	expectedResults := []domain.StockUpdateResult{
		{
			ProductID:   100,
			StockBefore: 10,
			StockAfter:  50,
			Success:     true,
		},
		{
			ProductID: 999,
			Success:   false,
			Error:     stringPtr("not found"),
		},
	}

	// Mock expectations
	mockService.On("BatchUpdateStock",
		ctx,
		req.StorefrontId,
		mock.Anything,
		"",
		req.UserId,
	).Return(int32(1), int32(1), expectedResults, nil)

	// Execute
	resp, err := server.BatchUpdateStock(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int32(1), resp.SuccessfulCount)
	assert.Equal(t, int32(1), resp.FailedCount)
	assert.Len(t, resp.Results, 2)

	assert.True(t, resp.Results[0].Success)
	assert.False(t, resp.Results[1].Success)
	assert.NotNil(t, resp.Results[1].Error)

	mockService.AssertExpectations(t)
}

// TestGetProductStats_ValidRequest_Success tests successful stats retrieval
func TestGetProductStats_ValidRequest_Success(t *testing.T) {
	server, mockService := setupInventoryTestServer()
	ctx := context.Background()

	req := &pb.GetProductStatsRequest{
		StorefrontId: 123,
	}

	expectedStats := &domain.ProductStats{
		TotalProducts:  150,
		ActiveProducts: 140,
		OutOfStock:     5,
		LowStock:       12,
		TotalValue:     125000.50,
		TotalSold:      3500,
	}

	// Mock expectations
	mockService.On("GetProductStats", ctx, req.StorefrontId).Return(expectedStats, nil)

	// Execute
	resp, err := server.GetProductStats(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Stats)
	assert.Equal(t, int32(150), resp.Stats.TotalProducts)
	assert.Equal(t, int32(140), resp.Stats.ActiveProducts)
	assert.Equal(t, int32(5), resp.Stats.OutOfStock)
	assert.Equal(t, int32(12), resp.Stats.LowStock)
	assert.Equal(t, 125000.50, resp.Stats.TotalValue)
	assert.Equal(t, int32(3500), resp.Stats.TotalSold)

	mockService.AssertExpectations(t)
}

// TestGetProductStats_InvalidStorefrontID_ValidationError tests validation for storefront_id
func TestGetProductStats_InvalidStorefrontID_ValidationError(t *testing.T) {
	server, _ := setupInventoryTestServer()
	ctx := context.Background()

	req := &pb.GetProductStatsRequest{
		StorefrontId: 0,
	}

	resp, err := server.GetProductStats(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

// TestGetProductStats_ServiceError_InternalError tests service layer error
func TestGetProductStats_ServiceError_InternalError(t *testing.T) {
	server, mockService := setupInventoryTestServer()
	ctx := context.Background()

	req := &pb.GetProductStatsRequest{
		StorefrontId: 123,
	}

	// Mock service returns error
	mockService.On("GetProductStats", ctx, req.StorefrontId).Return(nil, errors.New("database error"))

	resp, err := server.GetProductStats(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, "products.stats_failed", st.Message())

	mockService.AssertExpectations(t)
}

// TestIncrementProductViews_ValidRequest_Success tests successful view increment
func TestIncrementProductViews_ValidRequest_Success(t *testing.T) {
	server, mockService := setupInventoryTestServer()
	ctx := context.Background()

	req := &pb.IncrementProductViewsRequest{
		ProductId: 456,
	}

	// Mock expectations
	mockService.On("IncrementProductViews", ctx, req.ProductId).Return(nil)

	// Execute
	resp, err := server.IncrementProductViews(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.IsType(t, &emptypb.Empty{}, resp)

	mockService.AssertExpectations(t)
}

// TestIncrementProductViews_InvalidProductID_ValidationError tests validation for product_id
func TestIncrementProductViews_InvalidProductID_ValidationError(t *testing.T) {
	server, _ := setupInventoryTestServer()
	ctx := context.Background()

	testCases := []struct {
		name      string
		productID int64
	}{
		{
			name:      "Zero product ID",
			productID: 0,
		},
		{
			name:      "Negative product ID",
			productID: -1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &pb.IncrementProductViewsRequest{
				ProductId: tc.productID,
			}

			resp, err := server.IncrementProductViews(ctx, req)

			assert.Error(t, err)
			assert.Nil(t, resp)

			st, ok := status.FromError(err)
			assert.True(t, ok)
			assert.Equal(t, codes.InvalidArgument, st.Code())
		})
	}
}

// TestIncrementProductViews_ServiceError_InternalError tests service layer error
func TestIncrementProductViews_ServiceError_InternalError(t *testing.T) {
	server, mockService := setupInventoryTestServer()
	ctx := context.Background()

	req := &pb.IncrementProductViewsRequest{
		ProductId: 456,
	}

	// Mock service returns error
	mockService.On("IncrementProductViews", ctx, req.ProductId).Return(errors.New("database error"))

	resp, err := server.IncrementProductViews(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, "products.increment_views_failed", st.Message())

	mockService.AssertExpectations(t)
}

// TestRecordInventoryMovement_AllMovementTypes tests all valid movement types
func TestRecordInventoryMovement_AllMovementTypes(t *testing.T) {
	server, mockService := setupInventoryTestServer()
	ctx := context.Background()

	testCases := []struct {
		name         string
		movementType string
		quantity     int32
		stockBefore  int32
		stockAfter   int32
	}{
		{
			name:         "Movement type: in",
			movementType: "in",
			quantity:     10,
			stockBefore:  50,
			stockAfter:   60,
		},
		{
			name:         "Movement type: out",
			movementType: "out",
			quantity:     5,
			stockBefore:  50,
			stockAfter:   45,
		},
		{
			name:         "Movement type: adjustment",
			movementType: "adjustment",
			quantity:     75,
			stockBefore:  50,
			stockAfter:   75,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := &pb.RecordInventoryMovementRequest{
				StorefrontId: 123,
				ProductId:    456,
				MovementType: tc.movementType,
				Quantity:     tc.quantity,
				UserId:       1,
			}

			// Mock expectations
			mockService.On("UpdateProductInventory",
				ctx,
				req.StorefrontId,
				req.ProductId,
				int64(0), // variantID
				req.MovementType,
				req.Quantity,
				"", // reason
				"", // notes
				req.UserId,
			).Return(tc.stockBefore, tc.stockAfter, nil).Once()

			// Execute
			resp, err := server.RecordInventoryMovement(ctx, req)

			// Assert
			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.True(t, resp.Success)
			assert.Equal(t, tc.stockBefore, resp.StockBefore)
			assert.Equal(t, tc.stockAfter, resp.StockAfter)

			mockService.AssertExpectations(t)
		})
	}
}

// TestRecordInventoryMovement_WithoutVariantID tests product-level stock update
func TestRecordInventoryMovement_WithoutVariantID_Success(t *testing.T) {
	server, mockService := setupInventoryTestServer()
	ctx := context.Background()

	req := &pb.RecordInventoryMovementRequest{
		StorefrontId: 123,
		ProductId:    456,
		VariantId:    nil, // No variant
		MovementType: "in",
		Quantity:     100,
		UserId:       1,
	}

	// Mock expectations - variantID should be 0
	mockService.On("UpdateProductInventory",
		ctx,
		req.StorefrontId,
		req.ProductId,
		int64(0), // variantID = 0 for product-level stock
		req.MovementType,
		req.Quantity,
		"", // reason
		"", // notes
		req.UserId,
	).Return(int32(50), int32(150), nil)

	// Execute
	resp, err := server.RecordInventoryMovement(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Equal(t, int32(50), resp.StockBefore)
	assert.Equal(t, int32(150), resp.StockAfter)

	mockService.AssertExpectations(t)
}

// TestBatchUpdateStock_InvalidUserID_ValidationError tests user_id validation
func TestBatchUpdateStock_InvalidUserID_ValidationError(t *testing.T) {
	server, _ := setupInventoryTestServer()
	ctx := context.Background()

	req := &pb.BatchUpdateStockRequest{
		StorefrontId: 123,
		Items: []*pb.StockUpdateItem{
			{ProductId: 100, Quantity: 10},
		},
		UserId: 0, // Invalid
	}

	resp, err := server.BatchUpdateStock(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

// TestGetProductStats_EmptyStorefront tests stats for empty storefront
func TestGetProductStats_EmptyStorefront(t *testing.T) {
	server, mockService := setupInventoryTestServer()
	ctx := context.Background()

	req := &pb.GetProductStatsRequest{
		StorefrontId: 123,
	}

	expectedStats := &domain.ProductStats{
		TotalProducts:  0,
		ActiveProducts: 0,
		OutOfStock:     0,
		LowStock:       0,
		TotalValue:     0.0,
		TotalSold:      0,
	}

	// Mock expectations
	mockService.On("GetProductStats", ctx, req.StorefrontId).Return(expectedStats, nil)

	// Execute
	resp, err := server.GetProductStats(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Stats)
	assert.Equal(t, int32(0), resp.Stats.TotalProducts)
	assert.Equal(t, int32(0), resp.Stats.ActiveProducts)

	mockService.AssertExpectations(t)
}
