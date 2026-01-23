package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	listingspb "github.com/vondi-global/listings/api/proto/listings/v1"
	"github.com/vondi-global/listings/internal/timeout"
)

// BulkDeleteProducts deletes multiple products in a single atomic operation
func (s *Server) BulkDeleteProducts(ctx context.Context, req *listingspb.BulkDeleteProductsRequest) (*listingspb.BulkDeleteProductsResponse, error) {
	s.logger.Debug().
		Int64("storefront_id", req.StorefrontId).
		Int("product_count", len(req.ProductIds)).
		Bool("hard_delete", req.HardDelete).
		Msg("BulkDeleteProducts called")

	// Check remaining time before starting (bulk delete requires at least 5s for cascade operations)
	if !timeout.HasSufficientTime(ctx, 5*time.Second) {
		s.logger.Warn().
			Dur("remaining", timeout.RemainingTime(ctx)).
			Msg("insufficient time for bulk delete operation")
		return nil, status.Error(codes.DeadlineExceeded, "insufficient time remaining for bulk delete operation")
	}

	// Validation
	if req.StorefrontId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "storefront ID must be greater than 0")
	}

	if len(req.ProductIds) == 0 {
		return nil, status.Error(codes.InvalidArgument, "product IDs list cannot be empty")
	}

	if len(req.ProductIds) > 1000 {
		return nil, status.Error(codes.InvalidArgument, "cannot delete more than 1000 products at once")
	}

	// Validate all product IDs are positive
	for i, id := range req.ProductIds {
		// Check context periodically during validation
		if i%100 == 0 {
			if err := timeout.CheckDeadline(ctx); err != nil {
				s.logger.Warn().Int("ids_validated", i).Msg("timeout during validation")
				return nil, status.Error(codes.DeadlineExceeded, "operation cancelled during validation")
			}
		}

		if id <= 0 {
			return nil, status.Error(codes.InvalidArgument, "all product IDs must be greater than 0")
		}
	}

	// Call service method
	successCount, failedCount, variantsDeleted, errorMap, err := s.service.BulkDeleteProducts(ctx, req.StorefrontId, req.ProductIds, req.HardDelete)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to bulk delete products")
		return nil, status.Error(codes.Internal, "products.bulk_delete_failed")
	}

	// Convert error map to proto errors
	protoErrors := make([]*listingspb.BulkOperationError, 0, len(errorMap))
	for productID, errorCode := range errorMap {
		id := productID
		protoErrors = append(protoErrors, &listingspb.BulkOperationError{
			ProductId:    &id,
			ErrorCode:    errorCode,
			ErrorMessage: errorCode, // Use error code as message (frontend will translate)
		})
	}

	s.logger.Info().
		Int32("success_count", successCount).
		Int32("failed_count", failedCount).
		Int32("variants_deleted", variantsDeleted).
		Msg("bulk delete products completed")

	return &listingspb.BulkDeleteProductsResponse{
		SuccessfulCount: successCount,
		FailedCount:     failedCount,
		VariantsDeleted: variantsDeleted,
		Errors:          protoErrors,
	}, nil
}
