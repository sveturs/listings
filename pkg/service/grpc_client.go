package service

import (
	"context"
	// TODO: Uncomment when proto files are generated in Sprint 4.4
	// pb "github.com/sveturs/listings/api/proto/listings/v1"
	// "google.golang.org/protobuf/types/known/timestamppb"
)

// TODO: Sprint 4.4 - These methods will be implemented when proto files are generated
// For now, they are stubs that will be replaced with actual gRPC calls

// getListingGRPC retrieves a listing via gRPC.
// TODO: Sprint 4.4 - Implement when proto files are generated
func (c *Client) getListingGRPC(ctx context.Context, id int64) (*Listing, error) {
	// Temporary stub - will be implemented in Sprint 4.4
	return nil, ErrUnavailable

	/*  // TODO: Uncomment in Sprint 4.4
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	req := &pb.GetListingRequest{
		Id: id,
	}

	resp, err := c.grpcClient.GetListing(ctx, req)
	if err != nil {
		return nil, convertGRPCError(err)
	}

	return protoToListing(resp.Listing), nil
	*/
}

// createListingGRPC creates a listing via gRPC.
// TODO: Sprint 4.4 - Implement when proto files are generated
func (c *Client) createListingGRPC(ctx context.Context, req *CreateListingRequest) (*Listing, error) {
	// Temporary stub - will be implemented in Sprint 4.4
	return nil, ErrUnavailable
}

// updateListingGRPC updates a listing via gRPC.
// TODO: Sprint 4.4 - Implement when proto files are generated
func (c *Client) updateListingGRPC(ctx context.Context, id int64, req *UpdateListingRequest) (*Listing, error) {
	// Temporary stub - will be implemented in Sprint 4.4
	return nil, ErrUnavailable
}

// deleteListingGRPC deletes a listing via gRPC.
// TODO: Sprint 4.4 - Implement when proto files are generated
func (c *Client) deleteListingGRPC(ctx context.Context, id int64) error {
	// Temporary stub - will be implemented in Sprint 4.4
	return ErrUnavailable
}

// searchListingsGRPC searches listings via gRPC.
// TODO: Sprint 4.4 - Implement when proto files are generated
func (c *Client) searchListingsGRPC(ctx context.Context, req *SearchListingsRequest) (*SearchListingsResponse, error) {
	// Temporary stub - will be implemented in Sprint 4.4
	return nil, ErrUnavailable
}

// listListingsGRPC lists listings via gRPC.
// TODO: Sprint 4.4 - Implement when proto files are generated
func (c *Client) listListingsGRPC(ctx context.Context, req *ListListingsRequest) (*ListListingsResponse, error) {
	// Temporary stub - will be implemented in Sprint 4.4
	return nil, ErrUnavailable
}

// TODO: Sprint 4.4 - Implement proto conversion functions when proto files are generated
// protoToListing converts a protobuf Listing to a domain Listing.
// listingToProto converts a domain Listing to a protobuf Listing.
