package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	listingspb "github.com/vondi-global/listings/api/proto/listings/v1"
	"github.com/vondi-global/listings/internal/domain"
	"github.com/vondi-global/listings/internal/service"
	"github.com/vondi-global/listings/internal/timeout"
)

// RecordInventoryMovement records a stock change with movement tracking
func (s *Server) RecordInventoryMovement(ctx context.Context, req *listingspb.RecordInventoryMovementRequest) (*listingspb.RecordInventoryMovementResponse, error) {
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

		// Record error metric
		errMsg := err.Error()
		if s.metrics != nil {
			errorReason := "unknown"
			switch errMsg {
			case "products.not_found":
				errorReason = "product_not_found"
			case "inventory.variant_not_found":
				errorReason = "variant_not_found"
			case "inventory.insufficient_stock":
				errorReason = "insufficient_stock"
			}
			s.metrics.RecordInventoryMovementError(errorReason)
		}

		// Check for specific errors with placeholders
		if errMsg == "products.not_found" {
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

	// Record success metric
	if s.metrics != nil {
		s.metrics.RecordInventoryMovement(req.MovementType)
	}

	s.logger.Info().
		Int64("product_id", req.ProductId).
		Int64("variant_id", variantID).
		Int32("stock_before", stockBefore).
		Int32("stock_after", stockAfter).
		Msg("inventory movement recorded successfully")

	return &listingspb.RecordInventoryMovementResponse{
		Success:     true,
		StockBefore: stockBefore,
		StockAfter:  stockAfter,
		Error:       nil,
	}, nil
}

// BatchUpdateStock updates stock for multiple products/variants atomically
func (s *Server) BatchUpdateStock(ctx context.Context, req *listingspb.BatchUpdateStockRequest) (*listingspb.BatchUpdateStockResponse, error) {
	s.logger.Debug().
		Int64("storefront_id", req.StorefrontId).
		Int("item_count", len(req.Items)).
		Msg("BatchUpdateStock called")

	// Check remaining time before starting (batch operation requires at least 5s)
	if !timeout.HasSufficientTime(ctx, 5*time.Second) {
		s.logger.Warn().
			Dur("remaining", timeout.RemainingTime(ctx)).
			Msg("insufficient time for batch update operation")
		return nil, status.Error(codes.DeadlineExceeded, "insufficient time remaining for batch operation")
	}

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
		// Check context periodically during validation
		if i%100 == 0 {
			if err := timeout.CheckDeadline(ctx); err != nil {
				s.logger.Warn().Int("items_validated", i).Msg("timeout during validation")
				return nil, status.Error(codes.DeadlineExceeded, "operation cancelled during validation")
			}
		}

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

		// Record error metric
		if s.metrics != nil {
			s.metrics.RecordInventoryStockOperation("batch_update", "error")
		}

		return nil, status.Error(codes.Internal, "inventory.batch_update_failed")
	}

	// Record success metric
	if s.metrics != nil {
		s.metrics.RecordInventoryStockOperation("batch_update", "success")
	}

	// Convert domain results to proto
	protoResults := make([]*listingspb.StockUpdateResult, len(results))
	for i, result := range results {
		protoResult := &listingspb.StockUpdateResult{
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

	return &listingspb.BatchUpdateStockResponse{
		SuccessfulCount: successCount,
		FailedCount:     failedCount,
		Results:         protoResults,
	}, nil
}

// GetProductStats retrieves statistics for storefront products
func (s *Server) GetProductStats(ctx context.Context, req *listingspb.GetProductStatsRequest) (*listingspb.GetProductStatsResponse, error) {
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
	protoStats := &listingspb.ProductStats{
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

	return &listingspb.GetProductStatsResponse{Stats: protoStats}, nil
}

// IncrementProductViews increments view counter for analytics
func (s *Server) IncrementProductViews(ctx context.Context, req *listingspb.IncrementProductViewsRequest) (*emptypb.Empty, error) {
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

		// Record error metric
		if s.metrics != nil {
			s.metrics.RecordInventoryProductViewError()
		}

		return nil, status.Error(codes.Internal, "products.increment_views_failed")
	}

	// Record success metric
	if s.metrics != nil {
		s.metrics.RecordInventoryProductView(fmt.Sprintf("%d", req.ProductId))
	}

	s.logger.Debug().Int64("product_id", req.ProductId).Msg("product views incremented successfully")
	return &emptypb.Empty{}, nil
}

// ============================================================================
// Inventory Reservation Handlers (for Transfers and other use cases)
// ============================================================================

// CreateReservation creates a temporary stock hold for transfers
func (s *Server) CreateReservation(ctx context.Context, req *listingspb.CreateReservationRequest) (*listingspb.CreateReservationResponse, error) {
	s.logger.Debug().
		Str("reference_type", req.ReferenceType).
		Int64("reference_id", req.ReferenceId).
		Msg("CreateReservation called")

	// Validation
	if req.Item == nil {
		return nil, status.Error(codes.InvalidArgument, "item is required")
	}

	if req.Item.ListingId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "listing_id must be greater than 0")
	}

	if req.Item.Quantity <= 0 {
		return nil, status.Error(codes.InvalidArgument, "quantity must be greater than 0")
	}

	if req.ReferenceType == "" {
		return nil, status.Error(codes.InvalidArgument, "reference_type is required")
	}

	if req.ReferenceId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "reference_id must be greater than 0")
	}

	// Validate reference type
	refType := domain.ReferenceType(req.ReferenceType)
	if refType != domain.ReferenceTypeOrder && refType != domain.ReferenceTypeTransfer {
		return nil, status.Error(codes.InvalidArgument, "reference_type must be 'order' or 'transfer'")
	}

	// Prepare TTL
	ttlMinutes := 30 // default
	if req.TtlMinutes != nil {
		ttlMinutes = int(*req.TtlMinutes)
	}

	// Create reservation request
	createReq := &service.CreateReservationRequest{
		ListingID:     req.Item.ListingId,
		VariantID:     req.Item.VariantId,
		ReferenceType: refType,
		ReferenceID:   req.ReferenceId,
		Quantity:      req.Item.Quantity,
		TTLMinutes:    ttlMinutes,
	}

	// Call inventory service
	reservation, err := s.inventoryService.CreateReservation(ctx, createReq)
	if err != nil {
		s.logger.Error().Err(err).
			Int64("listing_id", req.Item.ListingId).
			Str("reference_type", req.ReferenceType).
			Int64("reference_id", req.ReferenceId).
			Msg("failed to create reservation")

		errMsg := err.Error()
		errStr := errMsg

		// Check for specific errors
		switch {
		case contains(errMsg, "not found"):
			return &listingspb.CreateReservationResponse{
				Success: false,
				Error:   &errStr,
			}, nil
		case contains(errMsg, "stock") || contains(errMsg, "available"):
			return &listingspb.CreateReservationResponse{
				Success: false,
				Error:   &errStr,
			}, nil
		default:
			return nil, status.Error(codes.Internal, "reservation.create_failed")
		}
	}

	s.logger.Info().
		Int64("reservation_id", reservation.ID).
		Int64("listing_id", req.Item.ListingId).
		Str("reference_type", req.ReferenceType).
		Int64("reference_id", req.ReferenceId).
		Msg("reservation created successfully")

	return &listingspb.CreateReservationResponse{
		Success:       true,
		ReservationId: reservation.ID,
		ExpiresAt:     timestamppb.New(reservation.ExpiresAt),
	}, nil
}

// ReleaseReservation releases a previously created reservation
func (s *Server) ReleaseReservation(ctx context.Context, req *listingspb.ReleaseReservationRequest) (*listingspb.ReleaseReservationResponse, error) {
	s.logger.Debug().
		Int64("reservation_id", req.ReservationId).
		Msg("ReleaseReservation called")

	// Validation
	if req.ReservationId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "reservation_id must be greater than 0")
	}

	// Get reservation first to get quantity for response
	reservation, err := s.inventoryService.GetReservationByID(ctx, req.ReservationId)
	if err != nil {
		s.logger.Error().Err(err).Int64("reservation_id", req.ReservationId).Msg("failed to get reservation")
		errStr := err.Error()
		return &listingspb.ReleaseReservationResponse{
			Success: false,
			Error:   &errStr,
		}, nil
	}

	quantityToRelease := reservation.Quantity

	// Call inventory service to release
	if err := s.inventoryService.ReleaseReservation(ctx, req.ReservationId); err != nil {
		s.logger.Error().Err(err).Int64("reservation_id", req.ReservationId).Msg("failed to release reservation")
		errStr := err.Error()
		return &listingspb.ReleaseReservationResponse{
			Success: false,
			Error:   &errStr,
		}, nil
	}

	s.logger.Info().
		Int64("reservation_id", req.ReservationId).
		Int32("quantity_released", quantityToRelease).
		Msg("reservation released successfully")

	return &listingspb.ReleaseReservationResponse{
		Success:          true,
		QuantityReleased: quantityToRelease,
	}, nil
}

// CommitReservation commits reservation (decrements actual stock)
func (s *Server) CommitReservation(ctx context.Context, req *listingspb.CommitReservationRequest) (*listingspb.CommitReservationResponse, error) {
	s.logger.Debug().
		Int64("reservation_id", req.ReservationId).
		Msg("CommitReservation called")

	// Validation
	if req.ReservationId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "reservation_id must be greater than 0")
	}

	// Get reservation first to get quantity for response
	reservation, err := s.inventoryService.GetReservationByID(ctx, req.ReservationId)
	if err != nil {
		s.logger.Error().Err(err).Int64("reservation_id", req.ReservationId).Msg("failed to get reservation")
		errStr := err.Error()
		return &listingspb.CommitReservationResponse{
			Success: false,
			Error:   &errStr,
		}, nil
	}

	quantityToCommit := reservation.Quantity

	// Call inventory service to commit
	if err := s.inventoryService.CommitReservation(ctx, req.ReservationId); err != nil {
		s.logger.Error().Err(err).Int64("reservation_id", req.ReservationId).Msg("failed to commit reservation")
		errStr := err.Error()
		return &listingspb.CommitReservationResponse{
			Success: false,
			Error:   &errStr,
		}, nil
	}

	s.logger.Info().
		Int64("reservation_id", req.ReservationId).
		Int32("quantity_committed", quantityToCommit).
		Msg("reservation committed successfully")

	return &listingspb.CommitReservationResponse{
		Success:           true,
		QuantityCommitted: quantityToCommit,
	}, nil
}

// GetReservation retrieves reservation details by ID
func (s *Server) GetReservation(ctx context.Context, req *listingspb.GetReservationRequest) (*listingspb.GetReservationResponse, error) {
	s.logger.Debug().
		Int64("reservation_id", req.ReservationId).
		Msg("GetReservation called")

	// Validation
	if req.ReservationId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "reservation_id must be greater than 0")
	}

	// Call inventory service
	reservation, err := s.inventoryService.GetReservationByID(ctx, req.ReservationId)
	if err != nil {
		s.logger.Error().Err(err).Int64("reservation_id", req.ReservationId).Msg("failed to get reservation")
		errStr := err.Error()
		return &listingspb.GetReservationResponse{
			Success: false,
			Error:   &errStr,
		}, nil
	}

	// Convert to proto ReservationDetails
	details := &listingspb.ReservationDetails{
		Id:            reservation.ID,
		ListingId:     reservation.ListingID,
		VariantId:     reservation.VariantID,
		Quantity:      reservation.Quantity,
		Status:        string(reservation.Status),
		ReferenceType: string(reservation.ReferenceType),
		ReferenceId:   reservation.ReferenceID,
		ExpiresAt:     timestamppb.New(reservation.ExpiresAt),
		CreatedAt:     timestamppb.New(reservation.CreatedAt),
		UpdatedAt:     timestamppb.New(reservation.UpdatedAt),
	}

	s.logger.Debug().
		Int64("reservation_id", req.ReservationId).
		Str("status", string(reservation.Status)).
		Msg("reservation retrieved successfully")

	return &listingspb.GetReservationResponse{
		Success:     true,
		Reservation: details,
	}, nil
}
