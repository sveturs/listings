// Package domain defines core business entities and domain models for the listings microservice.
// It contains data structures, validation rules, and business logic types used across the application.
package domain

import (
	"time"
)

// Listing represents a marketplace listing entity
type Listing struct {
	ID             int64      `json:"id" db:"id"`
	UUID           string     `json:"uuid" db:"uuid"`
	Slug           string     `json:"slug" db:"slug"`
	UserID         int64      `json:"user_id" db:"user_id"`
	StorefrontID   *int64     `json:"storefront_id,omitempty" db:"storefront_id"`
	Title          string     `json:"title" db:"title"`
	Description    *string    `json:"description,omitempty" db:"description"`
	Price          float64    `json:"price" db:"price"`
	Currency       string     `json:"currency" db:"currency"`
	CategoryID     int64      `json:"category_id" db:"category_id"`
	Status         string     `json:"status" db:"status"`
	Visibility     string     `json:"visibility" db:"visibility"`
	Quantity       int32      `json:"quantity" db:"quantity"`
	SKU            *string    `json:"sku,omitempty" db:"sku"`
	SourceType     string     `json:"source_type" db:"source_type"`              // c2c or b2c
	StockStatus    *string    `json:"stock_status,omitempty" db:"stock_status"`  // in_stock, out_of_stock, low_stock, discontinued (enum from DB)
	AttributesJSON *string    `json:"attributes_json,omitempty" db:"attributes"` // JSONB column stored as string for flexibility
	ViewsCount     int32      `json:"views_count" db:"view_count"`
	FavoritesCount int32      `json:"favorites_count" db:"favorites_count"`
	ExpiresAt      *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
	PublishedAt    *time.Time `json:"published_at,omitempty" db:"published_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	IsDeleted      bool       `json:"is_deleted" db:"is_deleted"`

	// Relations (loaded on demand)
	Attributes []*ListingAttribute `json:"attributes,omitempty" db:"-"`
	Images     []*ListingImage     `json:"images,omitempty" db:"-"`
	Tags       []string            `json:"tags,omitempty" db:"-"`
	Location   *ListingLocation    `json:"location,omitempty" db:"-"`
}

// ListingAttribute represents flexible key-value attributes
type ListingAttribute struct {
	ID             int64     `json:"id" db:"id"`
	ListingID      int64     `json:"listing_id" db:"listing_id"`
	AttributeKey   string    `json:"attribute_key" db:"attribute_key"`
	AttributeValue string    `json:"attribute_value" db:"attribute_value"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// ListingImage represents an image associated with a listing
type ListingImage struct {
	ID           int64     `json:"id" db:"id"`
	ListingID    int64     `json:"listing_id" db:"listing_id"`
	URL          string    `json:"url" db:"url"`
	StoragePath  *string   `json:"storage_path,omitempty" db:"storage_path"`
	ThumbnailURL *string   `json:"thumbnail_url,omitempty" db:"thumbnail_url"`
	DisplayOrder int32     `json:"display_order" db:"display_order"`
	IsPrimary    bool      `json:"is_primary" db:"is_primary"`
	Width        *int32    `json:"width,omitempty" db:"width"`
	Height       *int32    `json:"height,omitempty" db:"height"`
	FileSize     *int64    `json:"file_size,omitempty" db:"file_size"`
	MimeType     *string   `json:"mime_type,omitempty" db:"mime_type"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// ListingLocation represents the geographic location of a listing
type ListingLocation struct {
	ID           int64     `json:"id" db:"id"`
	ListingID    int64     `json:"listing_id" db:"listing_id"`
	Country      *string   `json:"country,omitempty" db:"country"`
	City         *string   `json:"city,omitempty" db:"city"`
	PostalCode   *string   `json:"postal_code,omitempty" db:"postal_code"`
	AddressLine1 *string   `json:"address_line1,omitempty" db:"address_line1"`
	AddressLine2 *string   `json:"address_line2,omitempty" db:"address_line2"`
	Latitude     *float64  `json:"latitude,omitempty" db:"latitude"`
	Longitude    *float64  `json:"longitude,omitempty" db:"longitude"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// ListingStats represents cached statistics for a listing
type ListingStats struct {
	ListingID      int64      `json:"listing_id" db:"listing_id"`
	ViewsCount     int32      `json:"views_count" db:"view_count"`
	FavoritesCount int32      `json:"favorites_count" db:"favorites_count"`
	InquiriesCount int32      `json:"inquiries_count" db:"inquiries_count"`
	LastViewedAt   *time.Time `json:"last_viewed_at,omitempty" db:"last_viewed_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

// IndexingQueueItem represents a pending indexing operation
type IndexingQueueItem struct {
	ID           int64      `json:"id" db:"id"`
	ListingID    int64      `json:"listing_id" db:"listing_id"`
	Operation    string     `json:"operation" db:"operation"` // index, update, delete
	Status       string     `json:"status" db:"status"`       // pending, processing, completed, failed
	RetryCount   int32      `json:"retry_count" db:"retry_count"`
	MaxRetries   int32      `json:"max_retries" db:"max_retries"`
	ErrorMessage *string    `json:"error_message,omitempty" db:"error_message"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	ProcessedAt  *time.Time `json:"processed_at,omitempty" db:"processed_at"`
}

// CreateListingInput represents input for creating a new listing
type CreateListingInput struct {
	UserID       int64   `json:"user_id" validate:"required"`
	StorefrontID *int64  `json:"storefront_id,omitempty"`
	Title        string  `json:"title" validate:"required,min=3,max=255"`
	Description  *string `json:"description,omitempty"`
	Price        float64 `json:"price" validate:"required,gte=0"`
	Currency     string  `json:"currency" validate:"required,len=3"`
	CategoryID   int64   `json:"category_id" validate:"required"`
	Quantity     int32   `json:"quantity" validate:"required,gte=0"`
	SKU          *string `json:"sku,omitempty"`
	SourceType   string  `json:"source_type" validate:"required,oneof=c2c b2c"`
}

// UpdateListingInput represents input for updating an existing listing
type UpdateListingInput struct {
	Title       *string  `json:"title,omitempty" validate:"omitempty,min=3,max=255"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty" validate:"omitempty,gte=0"`
	Quantity    *int32   `json:"quantity,omitempty" validate:"omitempty,gte=0"`
	Status      *string  `json:"status,omitempty" validate:"omitempty,oneof=draft active inactive sold archived"`
}

// ListListingsFilter represents filters for listing queries
type ListListingsFilter struct {
	UserID       *int64   `json:"user_id,omitempty"`
	StorefrontID *int64   `json:"storefront_id,omitempty"`
	CategoryID   *int64   `json:"category_id,omitempty"`
	Status       *string  `json:"status,omitempty"`
	SourceType   *string  `json:"source_type,omitempty" validate:"omitempty,oneof=c2c b2c"`
	MinPrice     *float64 `json:"min_price,omitempty"`
	MaxPrice     *float64 `json:"max_price,omitempty"`
	Limit        int32    `json:"limit" validate:"required,gte=1,lte=100"`
	Offset       int32    `json:"offset" validate:"gte=0"`
}

// SearchListingsQuery represents a search query for listings
type SearchListingsQuery struct {
	Query      string   `json:"query" validate:"required,min=2"`
	CategoryID *int64   `json:"category_id,omitempty"`
	SourceType *string  `json:"source_type,omitempty" validate:"omitempty,oneof=c2c b2c"` // Filter by c2c or b2c listings
	MinPrice   *float64 `json:"min_price,omitempty"`
	MaxPrice   *float64 `json:"max_price,omitempty"`
	Limit      int32    `json:"limit" validate:"required,gte=1,lte=100"`
	Offset     int32    `json:"offset" validate:"gte=0"`
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

// Constants for listing source type
const (
	SourceTypeC2C = "c2c" // Consumer-to-Consumer
	SourceTypeB2C = "b2c" // Business-to-Consumer
)

// Constants for indexing operations
const (
	IndexOpIndex  = "index"
	IndexOpUpdate = "update"
	IndexOpDelete = "delete"
)

// Constants for indexing queue status
const (
	IndexStatusPending    = "pending"
	IndexStatusProcessing = "processing"
	IndexStatusCompleted  = "completed"
	IndexStatusFailed     = "failed"
)

// Category represents a marketplace category
type Category struct {
	ID                int64     `json:"id" db:"id"`
	Name              string    `json:"name" db:"name"`
	Slug              string    `json:"slug" db:"slug"`
	ParentID          *int64    `json:"parent_id,omitempty" db:"parent_id"`
	Icon              *string   `json:"icon,omitempty" db:"icon"`
	Description       *string   `json:"description,omitempty" db:"description"`
	IsActive          bool      `json:"is_active" db:"is_active"`
	ListingCount      int32     `json:"listing_count" db:"listing_count"`
	SortOrder         int32     `json:"sort_order" db:"sort_order"`
	Level             int32     `json:"level" db:"level"`
	HasCustomUI       bool      `json:"has_custom_ui" db:"has_custom_ui"`
	CustomUIComponent *string   `json:"custom_ui_component,omitempty" db:"custom_ui_component"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
}

// CategoryTreeNode represents a category with its children in a tree structure
type CategoryTreeNode struct {
	ID                int64              `json:"id"`
	Name              string             `json:"name"`
	Slug              string             `json:"slug"`
	Icon              *string            `json:"icon,omitempty"`
	ParentID          *int64             `json:"parent_id,omitempty"`
	Level             int32              `json:"level"`
	Path              string             `json:"path"`
	ListingCount      int32              `json:"listing_count"`
	ChildrenCount     int32              `json:"children_count"`
	Children          []CategoryTreeNode `json:"children"`
	HasCustomUI       bool               `json:"has_custom_ui"`
	CustomUIComponent *string            `json:"custom_ui_component,omitempty"`
	CreatedAt         string             `json:"created_at"`
}

// ListingVariant represents a product variant
type ListingVariant struct {
	ID         int64             `json:"id" db:"id"`
	ListingID  int64             `json:"listing_id" db:"listing_id"`
	SKU        string            `json:"sku" db:"sku"`
	Price      *float64          `json:"price,omitempty" db:"price"`
	Stock      *int32            `json:"stock,omitempty" db:"stock"`
	Attributes map[string]string `json:"attributes,omitempty" db:"attributes"`
	ImageURL   *string           `json:"image_url,omitempty" db:"image_url"`
	IsActive   bool              `json:"is_active" db:"is_active"`
	CreatedAt  *time.Time        `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt  *time.Time        `json:"updated_at,omitempty" db:"updated_at"`
}
