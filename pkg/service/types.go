// Package service provides a unified client for accessing the listings microservice.
// It supports both gRPC (primary) and HTTP (fallback) communication methods.
package service

import (
	"time"
)

// Listing represents a marketplace listing entity.
// This is the public API type that will be used by consumers of this library.
type Listing struct {
	ID             int64      `json:"id"`
	UUID           string     `json:"uuid"`
	UserID         int64      `json:"user_id"`
	StorefrontID   *int64     `json:"storefront_id,omitempty"`
	Title          string     `json:"title"`
	Description    *string    `json:"description,omitempty"`
	Price          float64    `json:"price"`
	Currency       string     `json:"currency"`
	CategoryID     int64      `json:"category_id"`
	Status         string     `json:"status"`
	Visibility     string     `json:"visibility"`
	Quantity       int32      `json:"quantity"`
	SKU            *string    `json:"sku,omitempty"`
	ViewsCount     int32      `json:"views_count"`
	FavoritesCount int32      `json:"favorites_count"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	PublishedAt    *time.Time `json:"published_at,omitempty"`

	// Relations (loaded on demand)
	Attributes []*ListingAttribute `json:"attributes,omitempty"`
	Images     []*ListingImage     `json:"images,omitempty"`
	Tags       []string            `json:"tags,omitempty"`
	Location   *ListingLocation    `json:"location,omitempty"`
}

// ListingAttribute represents flexible key-value attributes.
type ListingAttribute struct {
	ID             int64     `json:"id"`
	ListingID      int64     `json:"listing_id"`
	AttributeKey   string    `json:"attribute_key"`
	AttributeValue string    `json:"attribute_value"`
	CreatedAt      time.Time `json:"created_at"`
}

// ListingImage represents an image associated with a listing.
type ListingImage struct {
	ID           int64     `json:"id"`
	ListingID    int64     `json:"listing_id"`
	URL          string    `json:"url"`
	StoragePath  *string   `json:"storage_path,omitempty"`
	ThumbnailURL *string   `json:"thumbnail_url,omitempty"`
	DisplayOrder int32     `json:"display_order"`
	IsPrimary    bool      `json:"is_primary"`
	Width        *int32    `json:"width,omitempty"`
	Height       *int32    `json:"height,omitempty"`
	FileSize     *int64    `json:"file_size,omitempty"`
	MimeType     *string   `json:"mime_type,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ListingLocation represents the geographic location of a listing.
type ListingLocation struct {
	ID           int64     `json:"id"`
	ListingID    int64     `json:"listing_id"`
	Country      *string   `json:"country,omitempty"`
	City         *string   `json:"city,omitempty"`
	PostalCode   *string   `json:"postal_code,omitempty"`
	AddressLine1 *string   `json:"address_line1,omitempty"`
	AddressLine2 *string   `json:"address_line2,omitempty"`
	Latitude     *float64  `json:"latitude,omitempty"`
	Longitude    *float64  `json:"longitude,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateListingRequest represents input for creating a new listing.
type CreateListingRequest struct {
	UserID       int64   `json:"user_id" validate:"required"`
	StorefrontID *int64  `json:"storefront_id,omitempty"`
	Title        string  `json:"title" validate:"required,min=3,max=255"`
	Description  *string `json:"description,omitempty"`
	Price        float64 `json:"price" validate:"required,gte=0"`
	Currency     string  `json:"currency" validate:"required,len=3"`
	CategoryID   int64   `json:"category_id" validate:"required"`
	Quantity     int32   `json:"quantity" validate:"required,gte=0"`
	SKU          *string `json:"sku,omitempty"`
}

// UpdateListingRequest represents input for updating an existing listing.
type UpdateListingRequest struct {
	Title       *string  `json:"title,omitempty" validate:"omitempty,min=3,max=255"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty" validate:"omitempty,gte=0"`
	Quantity    *int32   `json:"quantity,omitempty" validate:"omitempty,gte=0"`
	Status      *string  `json:"status,omitempty" validate:"omitempty,oneof=draft active inactive sold archived"`
}

// ListListingsRequest represents filters for listing queries.
type ListListingsRequest struct {
	UserID       *int64   `json:"user_id,omitempty"`
	StorefrontID *int64   `json:"storefront_id,omitempty"`
	CategoryID   *int64   `json:"category_id,omitempty"`
	Status       *string  `json:"status,omitempty"`
	MinPrice     *float64 `json:"min_price,omitempty"`
	MaxPrice     *float64 `json:"max_price,omitempty"`
	Limit        int32    `json:"limit" validate:"required,gte=1,lte=100"`
	Offset       int32    `json:"offset" validate:"gte=0"`
}

// SearchListingsRequest represents a search query for listings.
type SearchListingsRequest struct {
	Query      string   `json:"query" validate:"required,min=2"`
	CategoryID *int64   `json:"category_id,omitempty"`
	MinPrice   *float64 `json:"min_price,omitempty"`
	MaxPrice   *float64 `json:"max_price,omitempty"`
	Limit      int32    `json:"limit" validate:"required,gte=1,lte=100"`
	Offset     int32    `json:"offset" validate:"gte=0"`
}

// ListListingsResponse represents a paginated list of listings.
type ListListingsResponse struct {
	Listings []*Listing `json:"listings"`
	Total    int32      `json:"total"`
}

// SearchListingsResponse represents search results.
type SearchListingsResponse struct {
	Listings []*Listing `json:"listings"`
	Total    int32      `json:"total"`
}

// Constants for listing statuses
const (
	StatusDraft    = "draft"
	StatusActive   = "active"
	StatusInactive = "inactive"
	StatusSold     = "sold"
	StatusArchived = "archived"
)

// Constants for listing visibility
const (
	VisibilityPublic   = "public"
	VisibilityPrivate  = "private"
	VisibilityUnlisted = "unlisted"
)
