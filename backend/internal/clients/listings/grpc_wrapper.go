package listings

import (
	"context"

	"backend/internal/domain/models"

	pb "github.com/sveturs/listings/api/proto/listings/v1"
)

// GRPCWrapper обертка над gRPC клиентом для упрощения использования
type GRPCWrapper struct {
	client *Client
}

// NewGRPCWrapper создает новую обертку над gRPC клиентом
func NewGRPCWrapper(client *Client) *GRPCWrapper {
	return &GRPCWrapper{
		client: client,
	}
}

// GetListing получает listing по ID и конвертирует в UnifiedListing
func (w *GRPCWrapper) GetListing(ctx context.Context, id int64) (*models.UnifiedListing, error) {
	req := &pb.GetListingRequest{
		Id: id,
	}

	resp, err := w.client.GetListing(ctx, req)
	if err != nil {
		return nil, err
	}

	// Конвертируем proto → UnifiedListing
	unified, err := ProtoToUnifiedListing(resp.Listing)
	if err != nil {
		return nil, err
	}

	return unified, nil
}

// CreateListing создает listing и конвертирует результат в UnifiedListing
func (w *GRPCWrapper) CreateListing(ctx context.Context, unified *models.UnifiedListing) (*models.UnifiedListing, error) {
	// Конвертируем UnifiedListing → proto request
	req, err := UnifiedToProtoCreateRequest(unified)
	if err != nil {
		return nil, err
	}

	resp, err := w.client.CreateListing(ctx, req)
	if err != nil {
		return nil, err
	}

	// Конвертируем proto → UnifiedListing
	created, err := ProtoToUnifiedListing(resp.Listing)
	if err != nil {
		return nil, err
	}

	return created, nil
}

// UpdateListing обновляет listing и конвертирует результат в UnifiedListing
func (w *GRPCWrapper) UpdateListing(ctx context.Context, unified *models.UnifiedListing) (*models.UnifiedListing, error) {
	// Конвертируем UnifiedListing → proto request
	req, err := UnifiedToProtoUpdateRequest(unified)
	if err != nil {
		return nil, err
	}

	resp, err := w.client.UpdateListing(ctx, req)
	if err != nil {
		return nil, err
	}

	// Конвертируем proto → UnifiedListing
	updated, err := ProtoToUnifiedListing(resp.Listing)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

// DeleteListing удаляет listing
func (w *GRPCWrapper) DeleteListing(ctx context.Context, id int64, userID int64) error {
	req := &pb.DeleteListingRequest{
		Id:     id,
		UserId: userID,
	}

	_, err := w.client.DeleteListing(ctx, req)
	return err
}
