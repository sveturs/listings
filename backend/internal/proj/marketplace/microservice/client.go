package microservice

import (
	"context"
	"fmt"
	"time"

	listingsv1 "backend/pkg/proto/listings/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client wraps gRPC client for listings microservice
type Client struct {
	conn   *grpc.ClientConn
	client listingsv1.ListingsServiceClient
	addr   string
}

// Config holds client configuration
type Config struct {
	Address        string        // gRPC server address (e.g., "localhost:50053")
	ConnectTimeout time.Duration // Connection timeout
}

// NewClient creates new listings microservice client
func NewClient(cfg Config) (*Client, error) {
	if cfg.ConnectTimeout == 0 {
		cfg.ConnectTimeout = 5 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectTimeout)
	defer cancel()

	// Connect to gRPC server
	conn, err := grpc.DialContext(ctx,
		cfg.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to listings microservice at %s: %w", cfg.Address, err)
	}

	return &Client{
		conn:   conn,
		client: listingsv1.NewListingsServiceClient(conn),
		addr:   cfg.Address,
	}, nil
}

// Close closes the gRPC connection
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetListing retrieves a listing by ID
func (c *Client) GetListing(ctx context.Context, id int64) (*listingsv1.Listing, error) {
	resp, err := c.client.GetListing(ctx, &listingsv1.GetListingRequest{
		Id: id,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get listing %d: %w", id, err)
	}
	return resp.Listing, nil
}

// CreateListing creates a new listing
func (c *Client) CreateListing(ctx context.Context, listing *listingsv1.Listing) (*listingsv1.Listing, error) {
	// Map Listing to CreateListingRequest
	req := &listingsv1.CreateListingRequest{
		UserId:     listing.UserId,
		Title:      listing.Title,
		Price:      listing.Price,
		Currency:   listing.Currency,
		CategoryId: listing.CategoryId,
		Quantity:   listing.Quantity,
	}

	// Optional fields
	if listing.StorefrontId != nil {
		req.StorefrontId = listing.StorefrontId
	}
	if listing.Description != nil {
		req.Description = listing.Description
	}
	if listing.Sku != nil {
		req.Sku = listing.Sku
	}

	resp, err := c.client.CreateListing(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create listing: %w", err)
	}
	return resp.Listing, nil
}

// UpdateListing updates an existing listing
func (c *Client) UpdateListing(ctx context.Context, listing *listingsv1.Listing) (*listingsv1.Listing, error) {
	// Map Listing to UpdateListingRequest
	req := &listingsv1.UpdateListingRequest{
		Id:     listing.Id,
		UserId: listing.UserId, // Required for ownership check
	}

	// Optional fields (only set if not nil/empty)
	if listing.Title != "" {
		req.Title = &listing.Title
	}
	if listing.Description != nil && *listing.Description != "" {
		req.Description = listing.Description
	}
	if listing.Price != 0 {
		req.Price = &listing.Price
	}
	if listing.Quantity != 0 {
		req.Quantity = &listing.Quantity
	}
	if listing.Status != "" {
		req.Status = &listing.Status
	}

	resp, err := c.client.UpdateListing(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update listing %d: %w", listing.Id, err)
	}
	return resp.Listing, nil
}

// DeleteListing soft-deletes a listing
func (c *Client) DeleteListing(ctx context.Context, id int64) error {
	_, err := c.client.DeleteListing(ctx, &listingsv1.DeleteListingRequest{
		Id: id,
	})
	if err != nil {
		return fmt.Errorf("failed to delete listing %d: %w", id, err)
	}
	return nil
}

// ListListings returns a paginated list of listings
func (c *Client) ListListings(ctx context.Context, req *listingsv1.ListListingsRequest) ([]*listingsv1.Listing, int32, error) {
	resp, err := c.client.ListListings(ctx, req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list listings: %w", err)
	}
	return resp.Listings, resp.Total, nil
}

// SearchListings performs full-text search
func (c *Client) SearchListings(ctx context.Context, req *listingsv1.SearchListingsRequest) ([]*listingsv1.Listing, int32, error) {
	resp, err := c.client.SearchListings(ctx, req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search listings: %w", err)
	}
	return resp.Listings, resp.Total, nil
}

// Health checks if the microservice is healthy
func (c *Client) Health(ctx context.Context) error {
	// Use ListListings with limit=0 as health check
	_, _, err := c.ListListings(ctx, &listingsv1.ListListingsRequest{
		Limit:  0,
		Offset: 0,
	})
	return err
}
