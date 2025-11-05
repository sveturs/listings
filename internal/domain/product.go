// Package domain defines core business entities and domain models for the listings microservice.
// This file contains product-related domain models for B2C storefront products.
package domain

import (
	"time"
)

// Product represents a B2C storefront product entity
type Product struct {
	ID                    int64                  `json:"id" db:"id"`
	StorefrontID          int64                  `json:"storefront_id" db:"storefront_id"`
	Name                  string                 `json:"name" db:"name"`
	Description           string                 `json:"description" db:"description"`
	Price                 float64                `json:"price" db:"price"`
	Currency              string                 `json:"currency" db:"currency"`
	CategoryID            int64                  `json:"category_id" db:"category_id"`
	SKU                   *string                `json:"sku,omitempty" db:"sku"`
	Barcode               *string                `json:"barcode,omitempty" db:"barcode"`
	StockQuantity         int32                  `json:"stock_quantity" db:"stock_quantity"`
	StockStatus           string                 `json:"stock_status" db:"stock_status"`
	IsActive              bool                   `json:"is_active" db:"is_active"`
	Attributes            map[string]interface{} `json:"attributes,omitempty" db:"attributes"` // JSONB
	ViewCount             int32                  `json:"view_count" db:"view_count"`
	SoldCount             int32                  `json:"sold_count" db:"sold_count"`
	CreatedAt             time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time              `json:"updated_at" db:"updated_at"`
	HasIndividualLocation bool                   `json:"has_individual_location" db:"has_individual_location"`
	IndividualAddress     *string                `json:"individual_address,omitempty" db:"individual_address"`
	IndividualLatitude    *float64               `json:"individual_latitude,omitempty" db:"individual_latitude"`
	IndividualLongitude   *float64               `json:"individual_longitude,omitempty" db:"individual_longitude"`
	LocationPrivacy       *string                `json:"location_privacy,omitempty" db:"location_privacy"`
	ShowOnMap             bool                   `json:"show_on_map" db:"show_on_map"`
	HasVariants           bool                   `json:"has_variants" db:"has_variants"`

	// Relations (loaded on demand)
	Variants []ProductVariant `json:"variants,omitempty" db:"-"`
	Images   []*ProductImage  `json:"images,omitempty" db:"-"`
}

// ProductVariant represents a product variant (size, color, etc.)
type ProductVariant struct {
	ID                int64                  `json:"id" db:"id"`
	ProductID         int64                  `json:"product_id" db:"product_id"`
	SKU               *string                `json:"sku,omitempty" db:"sku"`
	Barcode           *string                `json:"barcode,omitempty" db:"barcode"`
	Price             *float64               `json:"price,omitempty" db:"price"`
	CompareAtPrice    *float64               `json:"compare_at_price,omitempty" db:"compare_at_price"`
	CostPrice         *float64               `json:"cost_price,omitempty" db:"cost_price"`
	StockQuantity     int32                  `json:"stock_quantity" db:"stock_quantity"`
	StockStatus       string                 `json:"stock_status" db:"stock_status"`
	LowStockThreshold *int32                 `json:"low_stock_threshold,omitempty" db:"low_stock_threshold"`
	VariantAttributes map[string]interface{} `json:"variant_attributes,omitempty" db:"variant_attributes"` // JSONB
	Weight            *float64               `json:"weight,omitempty" db:"weight"`
	Dimensions        map[string]interface{} `json:"dimensions,omitempty" db:"dimensions"` // JSONB
	IsActive          bool                   `json:"is_active" db:"is_active"`
	IsDefault         bool                   `json:"is_default" db:"is_default"`
	ViewCount         int32                  `json:"view_count" db:"view_count"`
	SoldCount         int32                  `json:"sold_count" db:"sold_count"`
	CreatedAt         time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at" db:"updated_at"`

	// Relations (loaded on demand)
	Images []*ProductImage `json:"images,omitempty" db:"-"`
}

// ProductImage represents an image associated with a product or variant
type ProductImage struct {
	ID           int64      `json:"id" db:"id"`
	ProductID    *int64     `json:"product_id,omitempty" db:"product_id"`
	VariantID    *int64     `json:"variant_id,omitempty" db:"variant_id"`
	URL          string     `json:"url" db:"url"`
	StoragePath  *string    `json:"storage_path,omitempty" db:"storage_path"`
	ThumbnailURL *string    `json:"thumbnail_url,omitempty" db:"thumbnail_url"`
	DisplayOrder int32      `json:"display_order" db:"display_order"`
	IsPrimary    bool       `json:"is_primary" db:"is_primary"`
	Width        *int32     `json:"width,omitempty" db:"width"`
	Height       *int32     `json:"height,omitempty" db:"height"`
	FileSize     *int64     `json:"file_size,omitempty" db:"file_size"`
	MimeType     *string    `json:"mime_type,omitempty" db:"mime_type"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

// Constants for stock status
const (
	StockStatusInStock    = "in_stock"
	StockStatusLowStock   = "low_stock"
	StockStatusOutOfStock = "out_of_stock"
	StockStatusPreOrder   = "pre_order"
)

// Constants for location privacy
const (
	LocationPrivacyExact       = "exact"
	LocationPrivacyApproximate = "approximate"
	LocationPrivacyHidden      = "hidden"
)

// IsInStock returns true if the product or variant is in stock
func (p *Product) IsInStock() bool {
	return p.StockQuantity > 0 && p.StockStatus == StockStatusInStock
}

// IsLowStock returns true if the product stock is low (less than 10 units)
func (p *Product) IsLowStock() bool {
	return p.StockQuantity > 0 && p.StockQuantity < 10
}

// GetEffectiveStock returns the actual available stock quantity
// (can be extended to consider reservations in the future)
func (p *Product) GetEffectiveStock() int32 {
	if !p.IsActive {
		return 0
	}
	return p.StockQuantity
}

// GetEffectivePrice returns the effective price considering variants
func (p *Product) GetEffectivePrice() float64 {
	// If product has variants, return the lowest variant price
	if p.HasVariants && len(p.Variants) > 0 {
		minPrice := p.Price
		for _, v := range p.Variants {
			if v.IsActive && v.Price != nil && *v.Price < minPrice {
				minPrice = *v.Price
			}
		}
		return minPrice
	}
	return p.Price
}

// IsInStock returns true if the variant is in stock
func (pv *ProductVariant) IsInStock() bool {
	return pv.StockQuantity > 0 && pv.StockStatus == StockStatusInStock
}

// IsLowStock returns true if the variant stock is below the threshold
func (pv *ProductVariant) IsLowStock() bool {
	if pv.LowStockThreshold != nil {
		return pv.StockQuantity > 0 && pv.StockQuantity <= *pv.LowStockThreshold
	}
	// Default threshold is 10
	return pv.StockQuantity > 0 && pv.StockQuantity < 10
}

// GetEffectiveStock returns the actual available stock quantity for the variant
func (pv *ProductVariant) GetEffectiveStock() int32 {
	if !pv.IsActive {
		return 0
	}
	return pv.StockQuantity
}

// GetEffectivePrice returns the effective price (uses variant price if set, otherwise falls back to parent price)
func (pv *ProductVariant) GetEffectivePrice(parentPrice float64) float64 {
	if pv.Price != nil {
		return *pv.Price
	}
	return parentPrice
}

// HasDiscount returns true if the variant has a compare-at price higher than the current price
func (pv *ProductVariant) HasDiscount() bool {
	if pv.CompareAtPrice == nil || pv.Price == nil {
		return false
	}
	return *pv.CompareAtPrice > *pv.Price
}

// GetDiscountPercentage returns the discount percentage if applicable
func (pv *ProductVariant) GetDiscountPercentage() float64 {
	if !pv.HasDiscount() {
		return 0.0
	}
	return ((*pv.CompareAtPrice - *pv.Price) / *pv.CompareAtPrice) * 100.0
}

// CreateProductInput represents input for creating a new product
type CreateProductInput struct {
	StorefrontID          int64                  `json:"storefront_id" validate:"required"`
	Name                  string                 `json:"name" validate:"required,min=3,max=255"`
	Description           string                 `json:"description" validate:"required"`
	Price                 float64                `json:"price" validate:"required,gte=0"`
	Currency              string                 `json:"currency" validate:"required,len=3"`
	CategoryID            int64                  `json:"category_id" validate:"required"`
	SKU                   *string                `json:"sku,omitempty"`
	Barcode               *string                `json:"barcode,omitempty"`
	StockQuantity         int32                  `json:"stock_quantity" validate:"required,gte=0"`
	Attributes            map[string]interface{} `json:"attributes,omitempty"`
	HasIndividualLocation bool                   `json:"has_individual_location"`
	IndividualAddress     *string                `json:"individual_address,omitempty"`
	IndividualLatitude    *float64               `json:"individual_latitude,omitempty"`
	IndividualLongitude   *float64               `json:"individual_longitude,omitempty"`
	LocationPrivacy       *string                `json:"location_privacy,omitempty" validate:"omitempty,oneof=exact approximate hidden"`
	ShowOnMap             bool                   `json:"show_on_map"`
}

// UpdateProductInput represents input for updating an existing product
type UpdateProductInput struct {
	Name                  *string                 `json:"name,omitempty" validate:"omitempty,min=3,max=255"`
	Description           *string                 `json:"description,omitempty"`
	Price                 *float64                `json:"price,omitempty" validate:"omitempty,gte=0"`
	StockQuantity         *int32                  `json:"stock_quantity,omitempty" validate:"omitempty,gte=0"`
	StockStatus           *string                 `json:"stock_status,omitempty" validate:"omitempty,oneof=in_stock low_stock out_of_stock pre_order"`
	IsActive              *bool                   `json:"is_active,omitempty"`
	Attributes            map[string]interface{}  `json:"attributes,omitempty"`
	HasIndividualLocation *bool                   `json:"has_individual_location,omitempty"`
	IndividualAddress     *string                 `json:"individual_address,omitempty"`
	IndividualLatitude    *float64                `json:"individual_latitude,omitempty"`
	IndividualLongitude   *float64                `json:"individual_longitude,omitempty"`
	LocationPrivacy       *string                 `json:"location_privacy,omitempty" validate:"omitempty,oneof=exact approximate hidden"`
	ShowOnMap             *bool                   `json:"show_on_map,omitempty"`
}

// CreateVariantInput represents input for creating a new product variant
type CreateVariantInput struct {
	ProductID         int64                  `json:"product_id" validate:"required"`
	SKU               *string                `json:"sku,omitempty"`
	Barcode           *string                `json:"barcode,omitempty"`
	Price             *float64               `json:"price,omitempty" validate:"omitempty,gte=0"`
	CompareAtPrice    *float64               `json:"compare_at_price,omitempty" validate:"omitempty,gte=0"`
	CostPrice         *float64               `json:"cost_price,omitempty" validate:"omitempty,gte=0"`
	StockQuantity     int32                  `json:"stock_quantity" validate:"required,gte=0"`
	LowStockThreshold *int32                 `json:"low_stock_threshold,omitempty" validate:"omitempty,gte=0"`
	VariantAttributes map[string]interface{} `json:"variant_attributes,omitempty"`
	Weight            *float64               `json:"weight,omitempty" validate:"omitempty,gte=0"`
	Dimensions        map[string]interface{} `json:"dimensions,omitempty"`
	IsDefault         bool                   `json:"is_default"`
}

// UpdateVariantInput represents input for updating an existing variant
type UpdateVariantInput struct {
	SKU               *string                `json:"sku,omitempty"`
	Barcode           *string                `json:"barcode,omitempty"`
	Price             *float64               `json:"price,omitempty" validate:"omitempty,gte=0"`
	CompareAtPrice    *float64               `json:"compare_at_price,omitempty" validate:"omitempty,gte=0"`
	CostPrice         *float64               `json:"cost_price,omitempty" validate:"omitempty,gte=0"`
	StockQuantity     *int32                 `json:"stock_quantity,omitempty" validate:"omitempty,gte=0"`
	StockStatus       *string                `json:"stock_status,omitempty" validate:"omitempty,oneof=in_stock low_stock out_of_stock pre_order"`
	LowStockThreshold *int32                 `json:"low_stock_threshold,omitempty" validate:"omitempty,gte=0"`
	VariantAttributes map[string]interface{} `json:"variant_attributes,omitempty"`
	Weight            *float64               `json:"weight,omitempty" validate:"omitempty,gte=0"`
	Dimensions        map[string]interface{} `json:"dimensions,omitempty"`
	IsActive          *bool                  `json:"is_active,omitempty"`
	IsDefault         *bool                  `json:"is_default,omitempty"`
}

// ListProductsFilter represents filters for product queries
type ListProductsFilter struct {
	StorefrontID *int64   `json:"storefront_id,omitempty"`
	CategoryID   *int64   `json:"category_id,omitempty"`
	IsActive     *bool    `json:"is_active,omitempty"`
	InStock      *bool    `json:"in_stock,omitempty"`
	MinPrice     *float64 `json:"min_price,omitempty"`
	MaxPrice     *float64 `json:"max_price,omitempty"`
	Limit        int32    `json:"limit" validate:"required,gte=1,lte=100"`
	Offset       int32    `json:"offset" validate:"gte=0"`
}

// SearchProductsQuery represents a search query for products
type SearchProductsQuery struct {
	Query        string   `json:"query" validate:"required,min=2"`
	StorefrontID *int64   `json:"storefront_id,omitempty"`
	CategoryID   *int64   `json:"category_id,omitempty"`
	MinPrice     *float64 `json:"min_price,omitempty"`
	MaxPrice     *float64 `json:"max_price,omitempty"`
	InStock      *bool    `json:"in_stock,omitempty"`
	Limit        int32    `json:"limit" validate:"required,gte=1,lte=100"`
	Offset       int32    `json:"offset" validate:"gte=0"`
}

// BulkProductError represents an error for a single product in bulk operation
type BulkProductError struct {
	Index        int32   `json:"index"`
	ProductID    *int64  `json:"product_id,omitempty"`
	ErrorCode    string  `json:"error_code"`
	ErrorMessage string  `json:"error_message"`
}

// BulkUpdateProductInput represents input for bulk updating a single product
type BulkUpdateProductInput struct {
	ProductID             int64                  `json:"product_id" validate:"required,gt=0"`
	Name                  *string                `json:"name,omitempty" validate:"omitempty,min=3,max=255"`
	Description           *string                `json:"description,omitempty"`
	Price                 *float64               `json:"price,omitempty" validate:"omitempty,gte=0"`
	SKU                   *string                `json:"sku,omitempty"`
	Barcode               *string                `json:"barcode,omitempty"`
	StockQuantity         *int32                 `json:"stock_quantity,omitempty" validate:"omitempty,gte=0"`
	StockStatus           *string                `json:"stock_status,omitempty" validate:"omitempty,oneof=in_stock low_stock out_of_stock pre_order"`
	IsActive              *bool                  `json:"is_active,omitempty"`
	Attributes            map[string]interface{} `json:"attributes,omitempty"`
	HasIndividualLocation *bool                  `json:"has_individual_location,omitempty"`
	IndividualAddress     *string                `json:"individual_address,omitempty"`
	IndividualLatitude    *float64               `json:"individual_latitude,omitempty"`
	IndividualLongitude   *float64               `json:"individual_longitude,omitempty"`
	LocationPrivacy       *string                `json:"location_privacy,omitempty" validate:"omitempty,oneof=exact approximate hidden"`
	ShowOnMap             *bool                  `json:"show_on_map,omitempty"`
	UpdateMask            []string               `json:"update_mask,omitempty"` // Field mask for partial updates
}

// BulkUpdateProductsResult represents the result of a bulk update operation
type BulkUpdateProductsResult struct {
	SuccessfulProducts []*Product
	FailedUpdates      []BulkUpdateError
}

// BulkUpdateError represents an error for a single product in bulk operation
type BulkUpdateError struct {
	ProductID    int64  `json:"product_id"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

// ProductStats represents product statistics for a storefront
type ProductStats struct {
	TotalProducts  int32   `json:"total_products"`
	ActiveProducts int32   `json:"active_products"`
	OutOfStock     int32   `json:"out_of_stock"`
	LowStock       int32   `json:"low_stock"`
	TotalValue     float64 `json:"total_value"`
	TotalSold      int32   `json:"total_sold"`
}

// StockUpdateItem represents a single stock update in batch operation
type StockUpdateItem struct {
	ProductID int64   `json:"product_id"`
	VariantID *int64  `json:"variant_id,omitempty"`
	Quantity  int32   `json:"quantity"`
	Reason    *string `json:"reason,omitempty"`
}

// StockUpdateResult represents the result of a single stock update
type StockUpdateResult struct {
	ProductID   int64   `json:"product_id"`
	VariantID   *int64  `json:"variant_id,omitempty"`
	StockBefore int32   `json:"stock_before"`
	StockAfter  int32   `json:"stock_after"`
	Success     bool    `json:"success"`
	Error       *string `json:"error,omitempty"`
}
