package grpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/sveturs/listings/api/proto/listings/v1"
	"github.com/sveturs/listings/internal/metrics"
	"github.com/sveturs/listings/internal/service/listings"
)

// contains checks if substring is in string (case-insensitive helper)
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// Server implements gRPC ListingsServiceServer
type Server struct {
	pb.UnimplementedListingsServiceServer
	service *listings.Service
	metrics *metrics.Metrics
	logger  zerolog.Logger
}

// NewServer creates a new gRPC server instance
func NewServer(service *listings.Service, m *metrics.Metrics, logger zerolog.Logger) *Server {
	return &Server{
		service: service,
		metrics: m,
		logger:  logger.With().Str("component", "grpc_handler").Logger(),
	}
}

// GetListing retrieves a single listing by ID
func (s *Server) GetListing(ctx context.Context, req *pb.GetListingRequest) (*pb.GetListingResponse, error) {
	s.logger.Debug().Int64("listing_id", req.Id).Msg("GetListing called")

	// Validate request
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "listing ID must be greater than 0")
	}

	// Get listing from service
	listing, err := s.service.GetListing(ctx, req.Id)
	if err != nil {
		s.logger.Error().Err(err).Int64("listing_id", req.Id).Msg("failed to get listing")
		return nil, status.Error(codes.NotFound, fmt.Sprintf("listing not found: %v", err))
	}

	// Check for nil listing (defensive programming)
	if listing == nil {
		s.logger.Error().Int64("listing_id", req.Id).Msg("listing returned nil without error")
		return nil, status.Error(codes.NotFound, "listing not found")
	}

	// Convert to proto
	pbListing := DomainToProtoListing(listing)

	return &pb.GetListingResponse{
		Listing: pbListing,
	}, nil
}

// CreateListing creates a new listing
func (s *Server) CreateListing(ctx context.Context, req *pb.CreateListingRequest) (*pb.CreateListingResponse, error) {
	s.logger.Debug().Int64("user_id", req.UserId).Str("title", req.Title).Msg("CreateListing called")

	// Validate request
	if err := s.validateCreateListingRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Convert proto to domain input
	input := ProtoToCreateListingInput(req)

	// Create listing via service
	listing, err := s.service.CreateListing(ctx, input)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to create listing")

		// Map errors to appropriate gRPC codes (use contains for wrapped errors)
		errMsg := err.Error()
		if contains(errMsg, "category not found") || contains(errMsg, "fk_listings_category_id") {
			return nil, status.Error(codes.InvalidArgument, "category not found")
		}
		if contains(errMsg, "validation failed") {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create listing: %v", err))
	}

	// Convert to proto
	pbListing := DomainToProtoListing(listing)

	s.logger.Info().Int64("listing_id", listing.ID).Msg("listing created successfully")
	return &pb.CreateListingResponse{
		Listing: pbListing,
	}, nil
}

// UpdateListing updates an existing listing
func (s *Server) UpdateListing(ctx context.Context, req *pb.UpdateListingRequest) (*pb.UpdateListingResponse, error) {
	s.logger.Debug().Int64("listing_id", req.Id).Int64("user_id", req.UserId).Msg("UpdateListing called")

	// Validate request
	if err := s.validateUpdateListingRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Convert proto to domain input
	input := ProtoToUpdateListingInput(req)

	// Update listing via service (with ownership check)
	listing, err := s.service.UpdateListing(ctx, req.Id, req.UserId, input)
	if err != nil {
		s.logger.Error().Err(err).Int64("listing_id", req.Id).Msg("failed to update listing")

		// Map errors to appropriate gRPC codes (use contains for wrapped errors)
		errMsg := err.Error()
		if errMsg == "unauthorized: user does not own this listing" {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		if contains(errMsg, "listing not found") || contains(errMsg, "sql: no rows in result set") {
			return nil, status.Error(codes.NotFound, "listing not found")
		}

		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to update listing: %v", err))
	}

	// Convert to proto
	pbListing := DomainToProtoListing(listing)

	s.logger.Info().Int64("listing_id", listing.ID).Msg("listing updated successfully")
	return &pb.UpdateListingResponse{
		Listing: pbListing,
	}, nil
}

// DeleteListing soft-deletes a listing
func (s *Server) DeleteListing(ctx context.Context, req *pb.DeleteListingRequest) (*pb.DeleteListingResponse, error) {
	s.logger.Debug().Int64("listing_id", req.Id).Int64("user_id", req.UserId).Msg("DeleteListing called")

	// Validate request
	if req.Id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "listing ID must be greater than 0")
	}

	if req.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "user ID must be greater than 0")
	}

	// Delete listing via service (with ownership check)
	err := s.service.DeleteListing(ctx, req.Id, req.UserId)
	if err != nil {
		s.logger.Error().Err(err).Int64("listing_id", req.Id).Msg("failed to delete listing")

		// Map errors to appropriate gRPC codes (use contains for wrapped errors)
		errMsg := err.Error()
		if errMsg == "unauthorized: user does not own this listing" {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		if contains(errMsg, "listing not found") || contains(errMsg, "sql: no rows in result set") {
			return nil, status.Error(codes.NotFound, "listing not found")
		}

		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to delete listing: %v", err))
	}

	s.logger.Info().Int64("listing_id", req.Id).Msg("listing deleted successfully")
	return &pb.DeleteListingResponse{
		Success: true,
	}, nil
}

// SearchListings performs full-text search on listings
func (s *Server) SearchListings(ctx context.Context, req *pb.SearchListingsRequest) (*pb.SearchListingsResponse, error) {
	s.logger.Debug().Str("query", req.Query).Msg("SearchListings called")

	// Validate request
	if err := s.validateSearchListingsRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Convert proto to domain query
	query := ProtoToSearchListingsQuery(req)

	// Search listings via service
	listings, total, err := s.service.SearchListings(ctx, query)
	if err != nil {
		s.logger.Error().Err(err).Str("query", req.Query).Msg("failed to search listings")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to search listings: %v", err))
	}

	// Convert to proto
	pbListings := make([]*pb.Listing, len(listings))
	for i, listing := range listings {
		pbListings[i] = DomainToProtoListing(listing)
	}

	s.logger.Debug().Int("count", len(listings)).Int32("total", total).Msg("search completed")
	return &pb.SearchListingsResponse{
		Listings: pbListings,
		Total:    total,
	}, nil
}

// ListListings returns a paginated list of listings
func (s *Server) ListListings(ctx context.Context, req *pb.ListListingsRequest) (*pb.ListListingsResponse, error) {
	s.logger.Debug().Int32("limit", req.Limit).Int32("offset", req.Offset).Msg("ListListings called")

	// Validate request
	if err := s.validateListListingsRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Convert proto to domain filter
	filter := ProtoToListListingsFilter(req)

	// List listings via service
	listings, total, err := s.service.ListListings(ctx, filter)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to list listings")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to list listings: %v", err))
	}

	// Convert to proto
	pbListings := make([]*pb.Listing, len(listings))
	for i, listing := range listings {
		pbListings[i] = DomainToProtoListing(listing)
	}

	s.logger.Debug().Int("count", len(listings)).Int32("total", total).Msg("listings retrieved")
	return &pb.ListListingsResponse{
		Listings: pbListings,
		Total:    total,
	}, nil
}

// Validation helpers

func (s *Server) validateCreateListingRequest(req *pb.CreateListingRequest) error {
	if req.UserId <= 0 {
		return fmt.Errorf("user_id must be greater than 0")
	}

	if req.Title == "" {
		return fmt.Errorf("title is required")
	}

	if len(req.Title) < 3 {
		return fmt.Errorf("title must be at least 3 characters")
	}

	if len(req.Title) > 255 {
		return fmt.Errorf("title must not exceed 255 characters")
	}

	if req.Price <= 0 {
		return fmt.Errorf("price must be greater than 0")
	}

	if req.Currency == "" {
		return fmt.Errorf("currency is required")
	}

	if len(req.Currency) != 3 {
		return fmt.Errorf("currency must be 3 characters (ISO 4217)")
	}

	if req.CategoryId <= 0 {
		return fmt.Errorf("category_id must be greater than 0")
	}

	if req.Quantity < 0 {
		return fmt.Errorf("quantity cannot be negative")
	}

	return nil
}

func (s *Server) validateUpdateListingRequest(req *pb.UpdateListingRequest) error {
	if req.Id <= 0 {
		return fmt.Errorf("listing ID must be greater than 0")
	}

	if req.UserId <= 0 {
		return fmt.Errorf("user ID must be greater than 0")
	}

	// At least one field must be set
	if req.Title == nil && req.Description == nil && req.Price == nil && req.Quantity == nil && req.Status == nil {
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Validate individual fields if present
	if req.Title != nil {
		title := *req.Title
		if title == "" {
			return fmt.Errorf("title cannot be empty")
		}
		if len(title) < 3 {
			return fmt.Errorf("title must be at least 3 characters")
		}
		if len(title) > 255 {
			return fmt.Errorf("title must not exceed 255 characters")
		}
	}

	if req.Price != nil && *req.Price <= 0 {
		return fmt.Errorf("price must be greater than 0")
	}

	if req.Quantity != nil && *req.Quantity < 0 {
		return fmt.Errorf("quantity cannot be negative")
	}

	if req.Status != nil {
		validStatuses := map[string]bool{
			"draft":    true,
			"active":   true,
			"inactive": true,
			"sold":     true,
			"archived": true,
		}
		status := *req.Status
		if !validStatuses[status] {
			return fmt.Errorf("invalid status: %s", status)
		}
	}

	return nil
}

func (s *Server) validateSearchListingsRequest(req *pb.SearchListingsRequest) error {
	// Query is optional - if provided, validate it
	if req.Query != "" && len(req.Query) < 2 {
		return fmt.Errorf("search query must be at least 2 characters")
	}

	if req.Limit <= 0 {
		return fmt.Errorf("limit must be greater than 0")
	}

	if req.Limit > 100 {
		return fmt.Errorf("limit must not exceed 100")
	}

	if req.Offset < 0 {
		return fmt.Errorf("offset cannot be negative")
	}

	if req.MinPrice != nil && req.MaxPrice != nil {
		if *req.MinPrice > *req.MaxPrice {
			return fmt.Errorf("min_price cannot be greater than max_price")
		}
	}

	return nil
}

func (s *Server) validateListListingsRequest(req *pb.ListListingsRequest) error {
	if req.Limit <= 0 {
		return fmt.Errorf("limit must be greater than 0")
	}

	if req.Limit > 100 {
		return fmt.Errorf("limit must not exceed 100")
	}

	if req.Offset < 0 {
		return fmt.Errorf("offset cannot be negative")
	}

	if req.MinPrice != nil && req.MaxPrice != nil {
		if *req.MinPrice > *req.MaxPrice {
			return fmt.Errorf("min_price cannot be greater than max_price")
		}
	}

	return nil
}
