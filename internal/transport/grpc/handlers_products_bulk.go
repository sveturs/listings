package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/sveturs/listings/api/proto/listings/v1"
)

// BulkDeleteProducts deletes multiple products in a single atomic operation
func (s *Server) BulkDeleteProducts(ctx context.Context, req *pb.BulkDeleteProductsRequest) (*pb.BulkDeleteProductsResponse, error) {
	s.logger.Debug().
		Int64("storefront_id", req.StorefrontId).
		Int("product_count", len(req.ProductIds)).
		Bool("hard_delete", req.HardDelete).
		Msg("BulkDeleteProducts called")

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
	for _, id := range req.ProductIds {
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
	protoErrors := make([]*pb.BulkOperationError, 0, len(errorMap))
	for productID, errorCode := range errorMap {
		id := productID
		protoErrors = append(protoErrors, &pb.BulkOperationError{
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

	return &pb.BulkDeleteProductsResponse{
		SuccessfulCount: successCount,
		FailedCount:     failedCount,
		VariantsDeleted: variantsDeleted,
		Errors:          protoErrors,
	}, nil
}
