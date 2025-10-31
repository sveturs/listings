package service

// Client provides a Go client for accessing the listings service
// This is the PUBLIC API that other services will use
// TODO: Sprint 4.2 - Implement client methods
type Client struct {
	// httpClient *http.Client
	// grpcConn *grpc.ClientConn
	// baseURL string
}

// NewClient creates a new listings service client
func NewClient( /* config */ ) *Client {
	return &Client{}
}

// Placeholder client methods - will be implemented in Sprint 4.2
// GetListing(ctx context.Context, id int64) (*Listing, error)
// CreateListing(ctx context.Context, req *CreateListingRequest) (*Listing, error)
// UpdateListing(ctx context.Context, id int64, req *UpdateListingRequest) (*Listing, error)
// DeleteListing(ctx context.Context, id int64) error
// SearchListings(ctx context.Context, query string) ([]*Listing, error)
