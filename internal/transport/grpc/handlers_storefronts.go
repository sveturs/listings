package grpc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	pb "github.com/sveturs/listings/api/proto/listings/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetStorefront retrieves a single storefront by ID
func (s *Server) GetStorefront(ctx context.Context, req *pb.GetStorefrontRequest) (*pb.StorefrontResponse, error) {
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
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Warn().
				Int64("storefront_id", req.Id).
				Msg("Storefront not found")
			return nil, status.Error(codes.NotFound, "storefront not found")
		}

		s.logger.Error().
			Err(err).
			Int64("storefront_id", req.Id).
			Msg("Failed to get storefront")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get storefront: %v", err))
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

// GetStorefrontBySlug retrieves a storefront by slug
func (s *Server) GetStorefrontBySlug(ctx context.Context, req *pb.GetStorefrontBySlugRequest) (*pb.StorefrontResponse, error) {
	s.logger.Debug().Str("slug", req.Slug).Msg("GetStorefrontBySlug called")

	// Validation
	if req.Slug == "" {
		s.logger.Warn().Msg("Empty slug provided")
		return nil, status.Error(codes.InvalidArgument, "slug cannot be empty")
	}

	// Get storefront from service
	storefront, err := s.service.GetStorefrontBySlug(ctx, req.Slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Warn().
				Str("slug", req.Slug).
				Msg("Storefront not found by slug")
			return nil, status.Error(codes.NotFound, "storefront not found")
		}

		s.logger.Error().
			Err(err).
			Str("slug", req.Slug).
			Msg("Failed to get storefront by slug")
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get storefront: %v", err))
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

// ListStorefronts returns a paginated list of storefronts
func (s *Server) ListStorefronts(ctx context.Context, req *pb.ListStorefrontsRequest) (*pb.ListStorefrontsResponse, error) {
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
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to list storefronts: %v", err))
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
