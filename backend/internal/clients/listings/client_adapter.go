package listings

import (
	"context"
	"fmt"

	"backend/internal/domain/models"

	pb "github.com/sveturs/listings/api/proto/listings/v1"
)

// ClientAdapter адаптирует gRPC Client к интерфейсу ListingsGRPCClient
// который ожидает service layer (работает с models.UnifiedListing, не protobuf)
type ClientAdapter struct {
	client *Client
}

// NewClientAdapter создаёт новый ClientAdapter
func NewClientAdapter(client *Client) *ClientAdapter {
	return &ClientAdapter{
		client: client,
	}
}

// GetListing получает listing по ID
func (a *ClientAdapter) GetListing(ctx context.Context, id int64) (*models.UnifiedListing, error) {
	req := &pb.GetListingRequest{
		Id: id,
	}

	resp, err := a.client.GetListing(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get listing from microservice: %w", err)
	}

	if resp.Listing == nil {
		return nil, fmt.Errorf("listing not found: id=%d", id)
	}

	unified, err := ProtoToUnifiedListing(resp.Listing)
	if err != nil {
		return nil, fmt.Errorf("failed to convert proto to unified listing: %w", err)
	}

	return unified, nil
}

// CreateListing создаёт новый listing
func (a *ClientAdapter) CreateListing(ctx context.Context, unified *models.UnifiedListing) (*models.UnifiedListing, error) {
	req, err := UnifiedToProtoCreateRequest(unified)
	if err != nil {
		return nil, fmt.Errorf("failed to convert unified to proto create request: %w", err)
	}

	resp, err := a.client.CreateListing(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create listing in microservice: %w", err)
	}

	if resp.Listing == nil {
		return nil, fmt.Errorf("microservice returned nil listing")
	}

	created, err := ProtoToUnifiedListing(resp.Listing)
	if err != nil {
		return nil, fmt.Errorf("failed to convert proto to unified listing: %w", err)
	}

	return created, nil
}

// UpdateListing обновляет существующий listing
func (a *ClientAdapter) UpdateListing(ctx context.Context, unified *models.UnifiedListing) (*models.UnifiedListing, error) {
	req, err := UnifiedToProtoUpdateRequest(unified)
	if err != nil {
		return nil, fmt.Errorf("failed to convert unified to proto update request: %w", err)
	}

	resp, err := a.client.UpdateListing(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update listing in microservice: %w", err)
	}

	if resp.Listing == nil {
		return nil, fmt.Errorf("microservice returned nil listing")
	}

	updated, err := ProtoToUnifiedListing(resp.Listing)
	if err != nil {
		return nil, fmt.Errorf("failed to convert proto to unified listing: %w", err)
	}

	return updated, nil
}

// DeleteListing удаляет listing (soft delete)
func (a *ClientAdapter) DeleteListing(ctx context.Context, id int64, userID int64) error {
	req := &pb.DeleteListingRequest{
		Id:     id,
		UserId: userID,
	}

	_, err := a.client.DeleteListing(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to delete listing in microservice: %w", err)
	}

	return nil
}
