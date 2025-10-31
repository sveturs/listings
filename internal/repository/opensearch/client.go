package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	opensearch "github.com/opensearch-project/opensearch-go/v2"
	"github.com/rs/zerolog"

	"github.com/sveturs/listings/internal/domain"
)

// Client handles OpenSearch operations for listings
type Client struct {
	client *opensearch.Client
	index  string
	logger zerolog.Logger
}

// NewClient creates a new OpenSearch client
func NewClient(addresses []string, username, password, index string, logger zerolog.Logger) (*Client, error) {
	cfg := opensearch.Config{
		Addresses: addresses,
		Username:  username,
		Password:  password,
	}

	client, err := opensearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenSearch client: %w", err)
	}

	// Test connection
	res, err := client.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to OpenSearch: %w", err)
	}
	defer res.Body.Close()

	logger.Info().
		Strs("addresses", addresses).
		Str("index", index).
		Msg("OpenSearch client initialized")

	return &Client{
		client: client,
		index:  index,
		logger: logger.With().Str("component", "opensearch_client").Logger(),
	}, nil
}

// IndexListing indexes a listing document in OpenSearch
func (c *Client) IndexListing(ctx context.Context, listing *domain.Listing) error {
	// Prepare document for indexing
	doc := map[string]interface{}{
		"id":              listing.ID,
		"uuid":            listing.UUID,
		"user_id":         listing.UserID,
		"storefront_id":   listing.StorefrontID,
		"title":           listing.Title,
		"description":     listing.Description,
		"price":           listing.Price,
		"currency":        listing.Currency,
		"category_id":     listing.CategoryID,
		"status":          listing.Status,
		"visibility":      listing.Visibility,
		"quantity":        listing.Quantity,
		"sku":             listing.SKU,
		"views_count":     listing.ViewsCount,
		"favorites_count": listing.FavoritesCount,
		"created_at":      listing.CreatedAt,
		"updated_at":      listing.UpdatedAt,
		"published_at":    listing.PublishedAt,
	}

	body, err := json.Marshal(doc)
	if err != nil {
		c.logger.Error().Err(err).Int64("listing_id", listing.ID).Msg("failed to marshal document")
		return fmt.Errorf("failed to marshal document: %w", err)
	}

	// Index document
	res, err := c.client.Index(
		c.index,
		bytes.NewReader(body),
		c.client.Index.WithContext(ctx),
		c.client.Index.WithDocumentID(fmt.Sprintf("%d", listing.ID)),
		c.client.Index.WithRefresh("false"), // Async refresh
	)

	if err != nil {
		c.logger.Error().Err(err).Int64("listing_id", listing.ID).Msg("failed to index listing")
		return fmt.Errorf("failed to index listing: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		c.logger.Error().Int("status", res.StatusCode).Int64("listing_id", listing.ID).Msg("OpenSearch index error")
		return fmt.Errorf("OpenSearch index error: %s", res.Status())
	}

	c.logger.Debug().Int64("listing_id", listing.ID).Msg("listing indexed successfully")
	return nil
}

// UpdateListing updates a listing document in OpenSearch
func (c *Client) UpdateListing(ctx context.Context, listing *domain.Listing) error {
	// For simplicity, re-index the entire document
	return c.IndexListing(ctx, listing)
}

// DeleteListing removes a listing from OpenSearch index
func (c *Client) DeleteListing(ctx context.Context, listingID int64) error {
	res, err := c.client.Delete(
		c.index,
		fmt.Sprintf("%d", listingID),
		c.client.Delete.WithContext(ctx),
	)

	if err != nil {
		c.logger.Error().Err(err).Int64("listing_id", listingID).Msg("failed to delete listing from index")
		return fmt.Errorf("failed to delete listing: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		c.logger.Error().Int("status", res.StatusCode).Int64("listing_id", listingID).Msg("OpenSearch delete error")
		return fmt.Errorf("OpenSearch delete error: %s", res.Status())
	}

	c.logger.Debug().Int64("listing_id", listingID).Msg("listing deleted from index")
	return nil
}

// HealthCheck performs a health check on OpenSearch
func (c *Client) HealthCheck(ctx context.Context) error {
	res, err := c.client.Cluster.Health()
	if err != nil {
		return fmt.Errorf("OpenSearch health check failed: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("OpenSearch health check error: %s", res.Status())
	}

	return nil
}

// Close closes the OpenSearch client (no-op for opensearch-go client)
func (c *Client) Close() error {
	// opensearch-go client doesn't require explicit closing
	return nil
}
