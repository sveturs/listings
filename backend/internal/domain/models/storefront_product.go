package models

import (
	"time"
)

// StorefrontProduct represents a product in a storefront
type StorefrontProduct struct {
	ID            int                    `json:"id" db:"id"`
	StorefrontID  int                    `json:"storefront_id" db:"storefront_id"`
	Name          string                 `json:"name" db:"name"`
	Description   string                 `json:"description" db:"description"`
	Price         float64                `json:"price" db:"price"`
	Currency      string                 `json:"currency" db:"currency"`
	CategoryID    int                    `json:"category_id" db:"category_id"`
	SKU           *string                `json:"sku,omitempty" db:"sku"`
	Barcode       *string                `json:"barcode,omitempty" db:"barcode"`
	StockQuantity int                    `json:"stock_quantity" db:"stock_quantity"`
	StockStatus   string                 `json:"stock_status" db:"stock_status"` // in_stock, low_stock, out_of_stock
	IsActive      bool                   `json:"is_active" db:"is_active"`
	Attributes    map[string]interface{} `json:"attributes,omitempty" db:"attributes"`
	ViewCount     int                    `json:"view_count" db:"view_count"`
	SoldCount     int                    `json:"sold_count" db:"sold_count"`
	CreatedAt     time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at" db:"updated_at"`

	// Relations
	Images   []StorefrontProductImage   `json:"images,omitempty" db:"-"`
	Category *MarketplaceCategory       `json:"category,omitempty" db:"-"`
	Variants []StorefrontProductVariant `json:"variants,omitempty" db:"-"`
}

// StorefrontProductImage represents an image of a storefront product
type StorefrontProductImage struct {
	ID               int       `json:"id" db:"id"`
	StorefrontProductID int    `json:"storefront_product_id" db:"storefront_product_id"`
	ImageURL         string    `json:"image_url" db:"image_url"`
	ThumbnailURL     string    `json:"thumbnail_url" db:"thumbnail_url"`
	DisplayOrder     int       `json:"display_order" db:"display_order"`
	IsDefault        bool      `json:"is_default" db:"is_default"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// StorefrontProductVariant represents a variant of a product (e.g., size, color)
type StorefrontProductVariant struct {
	ID                  int                    `json:"id" db:"id"`
	StorefrontProductID int                    `json:"storefront_product_id" db:"storefront_product_id"`
	Name                string                 `json:"name" db:"name"` // e.g., "Red - Large"
	SKU                 *string                `json:"sku,omitempty" db:"sku"`
	Price               float64                `json:"price" db:"price"`
	StockQuantity       int                    `json:"stock_quantity" db:"stock_quantity"`
	Attributes          map[string]interface{} `json:"attributes,omitempty" db:"attributes"` // e.g., {"color": "red", "size": "L"}
	IsActive            bool                   `json:"is_active" db:"is_active"`
	CreatedAt           time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at" db:"updated_at"`
}

// StorefrontInventoryMovement tracks inventory changes
type StorefrontInventoryMovement struct {
	ID                  int       `json:"id" db:"id"`
	StorefrontProductID int       `json:"storefront_product_id" db:"storefront_product_id"`
	VariantID           *int      `json:"variant_id,omitempty" db:"variant_id"`
	Type                string    `json:"type" db:"type"` // in, out, adjustment
	Quantity            int       `json:"quantity" db:"quantity"`
	Reason              string    `json:"reason" db:"reason"` // sale, return, damage, restock, adjustment
	OrderID             *int      `json:"order_id,omitempty" db:"order_id"`
	Notes               *string   `json:"notes,omitempty" db:"notes"`
	UserID              int       `json:"user_id" db:"user_id"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
}

// ProductFilter represents filter options for products
type ProductFilter struct {
	StorefrontID  int      `json:"storefront_id"`
	CategoryID    *int     `json:"category_id,omitempty"`
	Search        *string  `json:"search,omitempty"`
	MinPrice      *float64 `json:"min_price,omitempty"`
	MaxPrice      *float64 `json:"max_price,omitempty"`
	StockStatus   *string  `json:"stock_status,omitempty"`
	IsActive      *bool    `json:"is_active,omitempty"`
	SKU           *string  `json:"sku,omitempty"`
	Barcode       *string  `json:"barcode,omitempty"`
	SortBy        string   `json:"sort_by,omitempty"` // name, price, created_at, stock_quantity
	SortOrder     string   `json:"sort_order,omitempty"` // asc, desc
	Limit         int      `json:"limit,omitempty"`
	Offset        int      `json:"offset,omitempty"`
}

// CreateProductRequest represents a request to create a product
type CreateProductRequest struct {
	Name          string                 `json:"name" validate:"required,min=3,max=255"`
	Description   string                 `json:"description" validate:"required,min=10"`
	Price         float64                `json:"price" validate:"required,min=0"`
	Currency      string                 `json:"currency" validate:"required,len=3"`
	CategoryID    int                    `json:"category_id" validate:"required"`
	SKU           *string                `json:"sku,omitempty"`
	Barcode       *string                `json:"barcode,omitempty"`
	StockQuantity int                    `json:"stock_quantity" validate:"min=0"`
	IsActive      bool                   `json:"is_active"`
	Attributes    map[string]interface{} `json:"attributes,omitempty"`
}

// UpdateProductRequest represents a request to update a product
type UpdateProductRequest struct {
	Name          *string                `json:"name,omitempty" validate:"omitempty,min=3,max=255"`
	Description   *string                `json:"description,omitempty" validate:"omitempty,min=10"`
	Price         *float64               `json:"price,omitempty" validate:"omitempty,min=0"`
	CategoryID    *int                   `json:"category_id,omitempty"`
	SKU           *string                `json:"sku,omitempty"`
	Barcode       *string                `json:"barcode,omitempty"`
	StockQuantity *int                   `json:"stock_quantity,omitempty" validate:"omitempty,min=0"`
	IsActive      *bool                  `json:"is_active,omitempty"`
	Attributes    map[string]interface{} `json:"attributes,omitempty"`
}

// UpdateInventoryRequest represents a request to update product inventory
type UpdateInventoryRequest struct {
	Quantity int     `json:"quantity" validate:"required"`
	Type     string  `json:"type" validate:"required,oneof=in out adjustment"`
	Reason   string  `json:"reason" validate:"required"`
	Notes    *string `json:"notes,omitempty"`
}

// BulkInventoryUpdate represents a bulk inventory update request
type BulkInventoryUpdate struct {
	Updates []struct {
		ProductID     int  `json:"product_id" validate:"required"`
		VariantID     *int `json:"variant_id,omitempty"`
		StockQuantity int  `json:"stock_quantity" validate:"min=0"`
	} `json:"updates" validate:"required,dive"`
}

// ProductStats represents product statistics
type ProductStats struct {
	TotalProducts   int     `json:"total_products"`
	ActiveProducts  int     `json:"active_products"`
	OutOfStock      int     `json:"out_of_stock"`
	LowStock        int     `json:"low_stock"`
	TotalValue      float64 `json:"total_value"`
	TotalSold       int     `json:"total_sold"`
}

// GetStockStatus calculates the stock status based on quantity
func (p *StorefrontProduct) GetStockStatus() string {
	switch {
	case p.StockQuantity == 0:
		return "out_of_stock"
	case p.StockQuantity <= 5: // TODO: Make this configurable
		return "low_stock"
	default:
		return "in_stock"
	}
}

// CalculateTotalStock calculates total stock including variants
func (p *StorefrontProduct) CalculateTotalStock() int {
	total := p.StockQuantity
	for _, variant := range p.Variants {
		if variant.IsActive {
			total += variant.StockQuantity
		}
	}
	return total
}

// Bulk operation models

// BulkCreateProductsRequest represents a request to create multiple products
type BulkCreateProductsRequest struct {
	Products []CreateProductRequest `json:"products" validate:"required,min=1,max=100,dive"`
}

// BulkCreateProductsResponse represents the response for bulk product creation
type BulkCreateProductsResponse struct {
	Created []int                `json:"created"`    // IDs of successfully created products
	Failed  []BulkOperationError `json:"failed"`     // Errors for failed operations
}

// BulkUpdateProductsRequest represents a request to update multiple products
type BulkUpdateProductsRequest struct {
	Updates []BulkUpdateItem `json:"updates" validate:"required,min=1,max=100,dive"`
}

// BulkUpdateItem represents a single product update in bulk operation
type BulkUpdateItem struct {
	ProductID int                  `json:"product_id" validate:"required"`
	Updates   UpdateProductRequest `json:"updates" validate:"required"`
}

// BulkUpdateProductsResponse represents the response for bulk product updates
type BulkUpdateProductsResponse struct {
	Updated []int                `json:"updated"`    // IDs of successfully updated products
	Failed  []BulkOperationError `json:"failed"`     // Errors for failed operations
}

// BulkDeleteProductsRequest represents a request to delete multiple products
type BulkDeleteProductsRequest struct {
	ProductIDs []int `json:"product_ids" validate:"required,min=1,max=100"`
}

// BulkDeleteProductsResponse represents the response for bulk product deletion
type BulkDeleteProductsResponse struct {
	Deleted []int                `json:"deleted"`    // IDs of successfully deleted products
	Failed  []BulkOperationError `json:"failed"`     // Errors for failed operations
}

// BulkUpdateStatusRequest represents a request to update status of multiple products
type BulkUpdateStatusRequest struct {
	ProductIDs []int `json:"product_ids" validate:"required,min=1,max=100"`
	IsActive   bool  `json:"is_active"`
}

// BulkUpdateStatusResponse represents the response for bulk status update
type BulkUpdateStatusResponse struct {
	Updated []int                `json:"updated"`    // IDs of successfully updated products
	Failed  []BulkOperationError `json:"failed"`     // Errors for failed operations
}

// BulkOperationError represents an error for a single item in bulk operation
type BulkOperationError struct {
	Index     int    `json:"index,omitempty"`     // Index in the request array
	ProductID int    `json:"product_id,omitempty"` // Product ID if available
	Error     string `json:"error"`                // Error message
}