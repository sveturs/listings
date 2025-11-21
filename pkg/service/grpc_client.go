package service

import (
	"context"

	pb "github.com/sveturs/listings/api/proto/listings/v1"
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

// ============================================================================
// Favorites gRPC Methods
// ============================================================================

// addToFavoritesGRPC adds a listing to user's favorites via gRPC.
func (c *Client) addToFavoritesGRPC(ctx context.Context, userID, listingID int64) error {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	req := &pb.AddToFavoritesRequest{
		UserId:    userID,
		ListingId: listingID,
	}

	_, err := c.grpcClient.AddToFavorites(ctx, req)
	if err != nil {
		return convertGRPCError(err)
	}

	c.logger.Debug().Int64("user_id", userID).Int64("listing_id", listingID).Msg("added to favorites via gRPC")
	return nil
}

// removeFromFavoritesGRPC removes a listing from user's favorites via gRPC.
func (c *Client) removeFromFavoritesGRPC(ctx context.Context, userID, listingID int64) error {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	req := &pb.RemoveFromFavoritesRequest{
		UserId:    userID,
		ListingId: listingID,
	}

	_, err := c.grpcClient.RemoveFromFavorites(ctx, req)
	if err != nil {
		return convertGRPCError(err)
	}

	c.logger.Debug().Int64("user_id", userID).Int64("listing_id", listingID).Msg("removed from favorites via gRPC")
	return nil
}

// getUserFavoritesGRPC retrieves list of listing IDs favorited by a user via gRPC.
func (c *Client) getUserFavoritesGRPC(ctx context.Context, userID int64) ([]int64, int, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	req := &pb.GetUserFavoritesRequest{
		UserId: userID,
		// Note: Current proto doesn't use limit/offset for this method
		// It returns all favorites. If pagination is needed, update proto first.
	}

	resp, err := c.grpcClient.GetUserFavorites(ctx, req)
	if err != nil {
		return nil, 0, convertGRPCError(err)
	}

	c.logger.Debug().Int64("user_id", userID).Int("count", len(resp.ListingIds)).Msg("retrieved user favorites via gRPC")
	return resp.ListingIds, int(resp.Total), nil
}

// isFavoriteGRPC checks if a listing is in user's favorites via gRPC.
func (c *Client) isFavoriteGRPC(ctx context.Context, userID, listingID int64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	req := &pb.IsFavoriteRequest{
		UserId:    userID,
		ListingId: listingID,
	}

	resp, err := c.grpcClient.IsFavorite(ctx, req)
	if err != nil {
		return false, convertGRPCError(err)
	}

	c.logger.Debug().Int64("user_id", userID).Int64("listing_id", listingID).Bool("is_favorite", resp.IsFavorite).Msg("checked favorite status via gRPC")
	return resp.IsFavorite, nil
}

// getFavoritedUsersGRPC retrieves list of user IDs who favorited a listing via gRPC.
func (c *Client) getFavoritedUsersGRPC(ctx context.Context, listingID int64) ([]int64, error) {
	ctx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	req := &pb.ListingIDRequest{
		ListingId: listingID,
	}

	resp, err := c.grpcClient.GetFavoritedUsers(ctx, req)
	if err != nil {
		return nil, convertGRPCError(err)
	}

	c.logger.Debug().Int64("listing_id", listingID).Int("count", len(resp.UserIds)).Msg("retrieved favorited users via gRPC")
	return resp.UserIds, nil
}

// ============================================================================
// Image Management gRPC Methods
// ============================================================================

// deleteListingImageGRPC removes an image from a listing via gRPC.
// NOTE: This method requires listing_id and user_id which are not available
// in the current Client API signature. Returns ErrUnavailable until the
// public DeleteListingImage method is updated to include these parameters.
// TODO: Update Client.DeleteListingImage signature to include listingID and userID
func (c *Client) deleteListingImageGRPC(ctx context.Context, imageID int64) error {
	// The new proto requires DeleteListingImageRequest with listing_id, image_id, user_id
	// but the current public API only provides imageID.
	// Return ErrUnavailable to fall back to HTTP transport.
	c.logger.Debug().Int64("image_id", imageID).Msg("deleteListingImageGRPC requires listing_id and user_id - falling back to HTTP")
	return ErrUnavailable
}

// reorderListingImagesGRPC updates display order for multiple images via gRPC.
// NOTE: This method requires user_id for authorization which is not available
// in the current Client API signature. Returns ErrUnavailable until the
// public ReorderListingImages method is updated to include userID.
// TODO: Update Client.ReorderListingImages signature to include userID
func (c *Client) reorderListingImagesGRPC(ctx context.Context, listingID int64, imageOrders []ImageOrder) error {
	// The new proto requires ReorderImagesRequest with listing_id, user_id, image_ids
	// but the current public API doesn't provide user_id.
	// Return ErrUnavailable to fall back to HTTP transport.
	c.logger.Debug().Int64("listing_id", listingID).Int("count", len(imageOrders)).Msg("reorderListingImagesGRPC requires user_id - falling back to HTTP")
	return ErrUnavailable
}
