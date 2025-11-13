package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	listingspb "github.com/sveturs/listings/api/proto/listings/v1"
	"github.com/sveturs/listings/internal/service/listings"
)

// DecrementStock atomically decrements stock for multiple products/variants
// Used during order creation to reserve inventory
func (s *Server) DecrementStock(ctx context.Context, req *listingspb.DecrementStockRequest) (*listingspb.DecrementStockResponse, error) {
	logger := s.logger.With().
		Str("method", "DecrementStock").
		Int("items_count", len(req.Items)).
		Logger()

	if req.OrderId != nil {
		logger = logger.With().Str("order_id", *req.OrderId).Logger()
	}

	logger.Debug().Msg("DecrementStock called")

	// Validate request
	if len(req.Items) == 0 {
		return &listingspb.DecrementStockResponse{
			Success: false,
			Error:   strPtr("no items provided"),
		}, nil
	}

	// Convert proto items to domain
	items := make([]listings.StockItem, 0, len(req.Items))
	for _, item := range req.Items {
		if item.ProductId <= 0 {
			return &listingspb.DecrementStockResponse{
				Success: false,
				Error:   strPtr("invalid product_id"),
			}, nil
		}

		if item.Quantity <= 0 {
			return &listingspb.DecrementStockResponse{
				Success: false,
				Error:   strPtr("quantity must be positive"),
			}, nil
		}

		stockItem := listings.StockItem{
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
		}

		if item.VariantId != nil {
			stockItem.VariantID = item.VariantId
		}

		items = append(items, stockItem)
	}

	// Execute stock decrement
	results, err := s.service.DecrementStock(ctx, items, req.OrderId)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to decrement stock")

		// Convert results to proto
		pbResults := make([]*listingspb.StockResult, 0, len(results))
		for _, result := range results {
			pbResults = append(pbResults, &listingspb.StockResult{
				ProductId:   result.ProductID,
				VariantId:   result.VariantID,
				StockBefore: result.StockBefore,
				StockAfter:  result.StockAfter,
				Success:     result.Success,
				Error:       result.Error,
			})
		}

		return &listingspb.DecrementStockResponse{
			Success: false,
			Error:   strPtr(err.Error()),
			Results: pbResults,
		}, nil
	}

	// Convert successful results to proto
	pbResults := make([]*listingspb.StockResult, 0, len(results))
	for _, result := range results {
		pbResults = append(pbResults, &listingspb.StockResult{
			ProductId:   result.ProductID,
			VariantId:   result.VariantID,
			StockBefore: result.StockBefore,
			StockAfter:  result.StockAfter,
			Success:     result.Success,
			Error:       result.Error,
		})
	}

	logger.Info().
		Int("items_processed", len(results)).
		Msg("Stock decremented successfully")

	return &listingspb.DecrementStockResponse{
		Success: true,
		Results: pbResults,
	}, nil
}

// RollbackStock restores stock if order creation fails
// Compensating transaction for failed orders
func (s *Server) RollbackStock(ctx context.Context, req *listingspb.RollbackStockRequest) (*listingspb.RollbackStockResponse, error) {
	logger := s.logger.With().
		Str("method", "RollbackStock").
		Int("items_count", len(req.Items)).
		Logger()

	if req.OrderId != nil {
		logger = logger.With().Str("order_id", *req.OrderId).Logger()
	}

	logger.Debug().Msg("RollbackStock called")

	// Validate request
	if len(req.Items) == 0 {
		return &listingspb.RollbackStockResponse{
			Success: false,
			Error:   strPtr("no items provided"),
		}, nil
	}

	// Convert proto items to domain
	items := make([]listings.StockItem, 0, len(req.Items))
	for _, item := range req.Items {
		if item.ProductId <= 0 {
			return &listingspb.RollbackStockResponse{
				Success: false,
				Error:   strPtr("invalid product_id"),
			}, nil
		}

		if item.Quantity <= 0 {
			return &listingspb.RollbackStockResponse{
				Success: false,
				Error:   strPtr("quantity must be positive"),
			}, nil
		}

		stockItem := listings.StockItem{
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
		}

		if item.VariantId != nil {
			stockItem.VariantID = item.VariantId
		}

		items = append(items, stockItem)
	}

	// Execute stock rollback
	results, err := s.service.RollbackStock(ctx, items, req.OrderId)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to rollback stock")
		return &listingspb.RollbackStockResponse{
			Success: false,
			Error:   strPtr(err.Error()),
		}, nil
	}

	// Convert results to proto and check if all succeeded
	pbResults := make([]*listingspb.StockResult, 0, len(results))
	allSucceeded := true
	for _, result := range results {
		pbResults = append(pbResults, &listingspb.StockResult{
			ProductId:   result.ProductID,
			VariantId:   result.VariantID,
			StockBefore: result.StockBefore,
			StockAfter:  result.StockAfter,
			Success:     result.Success,
			Error:       result.Error,
		})

		if !result.Success {
			allSucceeded = false
		}
	}

	logger.Info().
		Int("items_processed", len(results)).
		Bool("all_succeeded", allSucceeded).
		Msg("Stock rollback completed")

	return &listingspb.RollbackStockResponse{
		Success: allSucceeded,
		Results: pbResults,
	}, nil
}

// CheckStockAvailability verifies if requested quantities are available
// Used for validation before order creation
func (s *Server) CheckStockAvailability(ctx context.Context, req *listingspb.CheckStockAvailabilityRequest) (*listingspb.CheckStockAvailabilityResponse, error) {
	logger := s.logger.With().
		Str("method", "CheckStockAvailability").
		Int("items_count", len(req.Items)).
		Logger()

	logger.Debug().Msg("CheckStockAvailability called")

	// Validate request
	if len(req.Items) == 0 {
		return nil, status.Error(codes.InvalidArgument, "no items provided")
	}

	// Convert proto items to domain
	items := make([]listings.StockItem, 0, len(req.Items))
	for _, item := range req.Items {
		if item.ProductId <= 0 {
			return nil, status.Error(codes.InvalidArgument, "invalid product_id")
		}

		if item.Quantity <= 0 {
			return nil, status.Error(codes.InvalidArgument, "quantity must be positive")
		}

		stockItem := listings.StockItem{
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
		}

		if item.VariantId != nil {
			stockItem.VariantID = item.VariantId
		}

		items = append(items, stockItem)
	}

	// Check availability
	allAvailable, availabilities, err := s.service.CheckStockAvailability(ctx, items)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to check stock availability")
		return nil, status.Error(codes.Internal, "failed to check stock availability")
	}

	// Convert results to proto
	pbAvailabilities := make([]*listingspb.StockAvailability, 0, len(availabilities))
	for _, avail := range availabilities {
		pbAvailabilities = append(pbAvailabilities, &listingspb.StockAvailability{
			ProductId:         avail.ProductID,
			VariantId:         avail.VariantID,
			RequestedQuantity: avail.RequestedQuantity,
			AvailableQuantity: avail.AvailableQuantity,
			IsAvailable:       avail.IsAvailable,
		})
	}

	logger.Debug().
		Bool("all_available", allAvailable).
		Msg("Stock availability checked")

	return &listingspb.CheckStockAvailabilityResponse{
		AllAvailable: allAvailable,
		Items:        pbAvailabilities,
	}, nil
}

// strPtr returns a pointer to a string (helper for optional fields)
func strPtr(s string) *string {
	return &s
}
