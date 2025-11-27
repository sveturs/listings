package grpc

import (
	"fmt"

	listingspb "github.com/vondi-global/listings/api/proto/listings/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ============================================================================
// ANALYTICS SERVICE HANDLERS
// ============================================================================

// GetOverviewStats retrieves platform-wide analytics statistics (admin only)
// Handler methods are defined in handlers.go under "ANALYTICS SERVICE HANDLERS"
// This file contains handler method wrappers and helper functions for analytics

// validateOverviewStatsRequest validates GetOverviewStatsRequest
func (s *Server) validateOverviewStatsRequestHelper(req *listingspb.GetOverviewStatsRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	// Validate date_from and date_to are present
	if req.DateFrom == nil {
		return fmt.Errorf("date_from is required")
	}
	if req.DateTo == nil {
		return fmt.Errorf("date_to is required")
	}

	// Validate date range is valid
	dateFrom := req.DateFrom.AsTime()
	dateTo := req.DateTo.AsTime()

	if dateFrom.After(dateTo) {
		return fmt.Errorf("date_from must be before or equal to date_to")
	}

	// Validate listing_type if provided
	if req.ListingType != nil && *req.ListingType != "" {
		listingType := *req.ListingType
		if listingType != "b2c" && listingType != "c2c" {
			return fmt.Errorf("listing_type must be either 'b2c' or 'c2c'")
		}
	}

	return nil
}

// validateListingStatsRequest validates GetListingStatsRequest
func (s *Server) validateListingStatsRequestHelper(req *listingspb.GetListingStatsRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	// Must have listing_id or product_id
	if req.GetListingId() == 0 && req.GetProductId() == 0 {
		return fmt.Errorf("either listing_id or product_id must be provided")
	}

	// Validate date_from and date_to are present
	if req.DateFrom == nil {
		return fmt.Errorf("date_from is required")
	}
	if req.DateTo == nil {
		return fmt.Errorf("date_to is required")
	}

	// Validate date range is valid
	dateFrom := req.DateFrom.AsTime()
	dateTo := req.DateTo.AsTime()

	if dateFrom.After(dateTo) {
		return fmt.Errorf("date_from must be before or equal to date_to")
	}

	return nil
}

// ============================================================================
// ERROR MAPPING
// ============================================================================

// mapAnalyticsError maps service errors to gRPC status codes
func mapAnalyticsError(err error, operation string) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()

	// Map authorization errors
	if contains(errMsg, "admin") || contains(errMsg, "unauthorized") || contains(errMsg, "permission") {
		return status.Errorf(codes.PermissionDenied, "%s: permission denied", operation)
	}

	// Map validation errors
	if contains(errMsg, "invalid") || contains(errMsg, "required") {
		return status.Errorf(codes.InvalidArgument, "%s: %v", operation, err)
	}

	// Map not found errors
	if contains(errMsg, "not found") {
		return status.Errorf(codes.NotFound, "%s: %v", operation, err)
	}

	// Default to internal error
	return status.Errorf(codes.Internal, "%s: %v", operation, err)
}
