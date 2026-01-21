package opensearch

import (
	"context"

	"github.com/vondi-global/listings/internal/domain"
)

// IndexingAdapter adapts OpenSearch Client to IndexingService interface
// This allows the service layer to use OpenSearch for indexing operations
type IndexingAdapter struct {
	client *Client
}

// NewIndexingAdapter creates a new IndexingAdapter
func NewIndexingAdapter(client *Client) *IndexingAdapter {
	return &IndexingAdapter{
		client: client,
	}
}

// IndexListing indexes a listing document in OpenSearch
func (a *IndexingAdapter) IndexListing(ctx context.Context, listing *domain.Listing) error {
	return a.client.IndexProduct(ctx, listing)
}

// UpdateListing updates a listing document in OpenSearch
func (a *IndexingAdapter) UpdateListing(ctx context.Context, listing *domain.Listing) error {
	return a.client.UpdateProduct(ctx, listing)
}

// DeleteListing removes a listing document from OpenSearch
func (a *IndexingAdapter) DeleteListing(ctx context.Context, listingID int64) error {
	return a.client.DeleteProduct(ctx, listingID)
}

// GetSimilarListings returns similar listings for a given listing
func (a *IndexingAdapter) GetSimilarListings(ctx context.Context, listingID int64, limit int32) ([]*domain.Listing, int32, error) {
	return a.client.GetSimilarListings(ctx, listingID, limit)
}

// DeleteAllDocuments deletes all documents from the index
// Used before full reindexing to remove stale/deleted documents
func (a *IndexingAdapter) DeleteAllDocuments(ctx context.Context) error {
	return a.client.DeleteAllDocuments(ctx)
}

// DeleteDocumentsBySourceType deletes documents with a specific source_type
func (a *IndexingAdapter) DeleteDocumentsBySourceType(ctx context.Context, sourceType string) error {
	return a.client.DeleteDocumentsBySourceType(ctx, sourceType)
}
