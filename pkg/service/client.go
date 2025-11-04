// Package service provides a unified client for accessing the listings microservice.
// It supports both gRPC (primary) and HTTP (fallback) communication methods.
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	// TODO: Uncomment when proto files are generated in Sprint 4.4
	// pb "github.com/sveturs/listings/api/proto/listings/v1"
)

// ErrNotFound is returned when a listing is not found.
var ErrNotFound = errors.New("listing not found")

// ErrUnavailable is returned when the service is unavailable.
var ErrUnavailable = errors.New("service unavailable")

// ErrInvalidInput is returned when input validation fails.
var ErrInvalidInput = errors.New("invalid input")

// ClientConfig holds configuration for the listings service client.
type ClientConfig struct {
	// GRPCAddr is the address of the gRPC server (e.g., "localhost:50053")
	GRPCAddr string

	// HTTPBaseURL is the base URL for HTTP fallback (e.g., "http://localhost:8086")
	HTTPBaseURL string

	// AuthToken is the authentication token for service-to-service communication
	AuthToken string

	// Timeout is the default timeout for requests
	Timeout time.Duration

	// EnableFallback enables HTTP fallback when gRPC fails
	EnableFallback bool

	// Logger for structured logging
	Logger zerolog.Logger
}

// Client provides a unified interface for accessing the listings microservice.
// It attempts gRPC first, and falls back to HTTP REST if gRPC is unavailable and fallback is enabled.
type Client struct {
	// TODO: Uncomment when proto files are generated in Sprint 4.4
	// grpcClient pb.ListingsServiceClient
	grpcClient interface{} // Temporary placeholder
	grpcConn   *grpc.ClientConn
	httpClient *HTTPClient
	config     ClientConfig
	logger     zerolog.Logger
}

// NewClient creates a new listings service client with the given configuration.
// It establishes both gRPC and HTTP connections (if fallback is enabled).
//
// Example:
//
//	client, err := service.NewClient(service.ClientConfig{
//	    GRPCAddr:       "localhost:50053",
//	    HTTPBaseURL:    "http://localhost:8086",
//	    AuthToken:      os.Getenv("SERVICE_TOKEN"),
//	    Timeout:        5 * time.Second,
//	    EnableFallback: true,
//	    Logger:         logger,
//	})
func NewClient(config ClientConfig) (*Client, error) {
	if config.Timeout == 0 {
		config.Timeout = 5 * time.Second
	}

	client := &Client{
		config: config,
		logger: config.Logger,
	}

	// Establish gRPC connection
	if config.GRPCAddr != "" {
		conn, err := grpc.NewClient(
			config.GRPCAddr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			client.logger.Warn().Err(err).Str("addr", config.GRPCAddr).Msg("Failed to connect to gRPC server")
			if !config.EnableFallback {
				return nil, fmt.Errorf("failed to connect to gRPC: %w", err)
			}
		} else {
			client.grpcConn = conn
			// TODO: Uncomment when proto files are generated in Sprint 4.4
			// client.grpcClient = pb.NewListingsServiceClient(conn)
			client.logger.Info().Str("addr", config.GRPCAddr).Msg("Connected to gRPC server")
		}
	}

	// Setup HTTP fallback client
	if config.EnableFallback && config.HTTPBaseURL != "" {
		httpClient, err := NewHTTPClient(config.HTTPBaseURL, config.AuthToken, config.Timeout, config.Logger)
		if err != nil {
			client.logger.Warn().Err(err).Msg("Failed to create HTTP client")
		} else {
			client.httpClient = httpClient
		}
	}

	// Ensure at least one transport is available
	// TODO: Change to grpcClient once proto is generated
	if client.grpcConn == nil && client.httpClient == nil {
		return nil, errors.New("no transport available: both gRPC and HTTP failed")
	}

	return client, nil
}

// Close closes all connections (gRPC and HTTP).
func (c *Client) Close() error {
	if c.grpcConn != nil {
		return c.grpcConn.Close()
	}
	return nil
}

// GetListing retrieves a single listing by ID.
// It tries gRPC first, then falls back to HTTP if enabled.
func (c *Client) GetListing(ctx context.Context, id int64) (*Listing, error) {
	if c.grpcClient != nil {
		listing, err := c.getListingGRPC(ctx, id)
		if err == nil {
			return listing, nil
		}

		c.logger.Warn().Err(err).Int64("id", id).Msg("gRPC GetListing failed")

		// If not a permanent error, don't fallback
		if !c.shouldFallback(err) {
			return nil, err
		}
	}

	// Fallback to HTTP
	if c.httpClient != nil {
		return c.httpClient.GetListing(ctx, id)
	}

	return nil, ErrUnavailable
}

// CreateListing creates a new listing.
// It tries gRPC first, then falls back to HTTP if enabled.
func (c *Client) CreateListing(ctx context.Context, req *CreateListingRequest) (*Listing, error) {
	if c.grpcClient != nil {
		listing, err := c.createListingGRPC(ctx, req)
		if err == nil {
			return listing, nil
		}

		c.logger.Warn().Err(err).Msg("gRPC CreateListing failed")

		if !c.shouldFallback(err) {
			return nil, err
		}
	}

	// Fallback to HTTP
	if c.httpClient != nil {
		return c.httpClient.CreateListing(ctx, req)
	}

	return nil, ErrUnavailable
}

// UpdateListing updates an existing listing.
// It tries gRPC first, then falls back to HTTP if enabled.
func (c *Client) UpdateListing(ctx context.Context, id int64, req *UpdateListingRequest) (*Listing, error) {
	if c.grpcClient != nil {
		listing, err := c.updateListingGRPC(ctx, id, req)
		if err == nil {
			return listing, nil
		}

		c.logger.Warn().Err(err).Int64("id", id).Msg("gRPC UpdateListing failed")

		if !c.shouldFallback(err) {
			return nil, err
		}
	}

	// Fallback to HTTP
	if c.httpClient != nil {
		return c.httpClient.UpdateListing(ctx, id, req)
	}

	return nil, ErrUnavailable
}

// DeleteListing soft-deletes a listing.
// It tries gRPC first, then falls back to HTTP if enabled.
func (c *Client) DeleteListing(ctx context.Context, id int64) error {
	if c.grpcClient != nil {
		err := c.deleteListingGRPC(ctx, id)
		if err == nil {
			return nil
		}

		c.logger.Warn().Err(err).Int64("id", id).Msg("gRPC DeleteListing failed")

		if !c.shouldFallback(err) {
			return err
		}
	}

	// Fallback to HTTP
	if c.httpClient != nil {
		return c.httpClient.DeleteListing(ctx, id)
	}

	return ErrUnavailable
}

// SearchListings performs full-text search on listings.
// It tries gRPC first, then falls back to HTTP if enabled.
func (c *Client) SearchListings(ctx context.Context, req *SearchListingsRequest) (*SearchListingsResponse, error) {
	if c.grpcClient != nil {
		resp, err := c.searchListingsGRPC(ctx, req)
		if err == nil {
			return resp, nil
		}

		c.logger.Warn().Err(err).Msg("gRPC SearchListings failed")

		if !c.shouldFallback(err) {
			return nil, err
		}
	}

	// Fallback to HTTP
	if c.httpClient != nil {
		return c.httpClient.SearchListings(ctx, req)
	}

	return nil, ErrUnavailable
}

// ListListings returns a paginated list of listings with filters.
// It tries gRPC first, then falls back to HTTP if enabled.
func (c *Client) ListListings(ctx context.Context, req *ListListingsRequest) (*ListListingsResponse, error) {
	if c.grpcClient != nil {
		resp, err := c.listListingsGRPC(ctx, req)
		if err == nil {
			return resp, nil
		}

		c.logger.Warn().Err(err).Msg("gRPC ListListings failed")

		if !c.shouldFallback(err) {
			return nil, err
		}
	}

	// Fallback to HTTP
	if c.httpClient != nil {
		return c.httpClient.ListListings(ctx, req)
	}

	return nil, ErrUnavailable
}

// shouldFallback determines if we should fallback to HTTP based on the error type.
// Only fallback for transient errors (unavailable, timeout, etc.)
func (c *Client) shouldFallback(err error) bool {
	if !c.config.EnableFallback || c.httpClient == nil {
		return false
	}

	// Check gRPC status codes
	st, ok := status.FromError(err)
	if !ok {
		// Not a gRPC error, might be transport error - fallback
		return true
	}

	// Fallback only for transient errors
	switch st.Code() {
	case codes.Unavailable, codes.DeadlineExceeded, codes.Canceled, codes.Unknown:
		return true
	default:
		return false
	}
}

// Helper method to convert gRPC error to public error
func convertGRPCError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.NotFound:
		return ErrNotFound
	case codes.InvalidArgument:
		return ErrInvalidInput
	case codes.Unavailable:
		return ErrUnavailable
	default:
		return err
	}
}
