package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/sveturs/listings/api/proto/listings/v1"
)

// ReindexAll performs full reindexing of all products to OpenSearch
// This is an administrative operation used for rebuilding search index after schema changes or data migration
func (s *Server) ReindexAll(ctx context.Context, req *pb.ReindexAllRequest) (*pb.ReindexAllResponse, error) {
	// Get source type and batch size from request (with defaults)
	sourceType := ""
	if req.SourceType != nil {
		sourceType = *req.SourceType
	}

	batchSize := 1000 // Default batch size
	if req.BatchSize != nil && *req.BatchSize > 0 && *req.BatchSize <= 10000 {
		batchSize = int(*req.BatchSize)
	}

	s.logger.Info().
		Str("source_type", sourceType).
		Int("batch_size", batchSize).
		Msg("ReindexAll gRPC called")

	// Check if OpenSearch indexer is configured
	// We validate this by attempting to call the service method
	// which will return an error if indexer is nil
	startTime := time.Now()

	// Call service layer to perform reindexing
	totalIndexed, totalFailed, durationSeconds, errors, err := s.service.ReindexAll(
		ctx,
		sourceType,
		batchSize,
	)

	if err != nil {
		s.logger.Error().Err(err).Msg("ReindexAll operation failed")

		// Map specific errors to gRPC codes
		errMsg := err.Error()
		switch errMsg {
		case "opensearch_not_configured":
			return nil, status.Error(codes.FailedPrecondition, "OpenSearch indexer is not configured")
		case "invalid source_type: must be 'b2c', 'c2c', or empty":
			return nil, status.Error(codes.InvalidArgument, "invalid source_type: must be 'b2c', 'c2c', or empty")
		default:
			// Check for context cancellation
			if ctx.Err() == context.Canceled {
				return nil, status.Error(codes.Canceled, "reindexing operation was cancelled")
			}
			if ctx.Err() == context.DeadlineExceeded {
				return nil, status.Error(codes.DeadlineExceeded, "reindexing operation timed out")
			}

			// Generic internal error
			return nil, status.Error(codes.Internal, "reindexing operation failed")
		}
	}

	// Calculate actual duration (in case service didn't measure it correctly)
	actualDuration := int32(time.Since(startTime).Seconds())
	if durationSeconds > 0 {
		actualDuration = int32(durationSeconds)
	}

	// Truncate errors list to max 10 items (as per proto spec)
	if len(errors) > 10 {
		errors = errors[:10]
	}

	s.logger.Info().
		Int32("total_indexed", totalIndexed).
		Int32("total_failed", totalFailed).
		Int32("duration_seconds", actualDuration).
		Int("error_count", len(errors)).
		Msg("ReindexAll operation completed successfully")

	return &pb.ReindexAllResponse{
		TotalIndexed:    totalIndexed,
		TotalFailed:     totalFailed,
		DurationSeconds: actualDuration,
		Errors:          errors,
	}, nil
}
