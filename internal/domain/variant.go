package domain

import "time"

// Variant represents a B2C product variant with complete inventory and pricing information
type Variant struct {
	ID                int64     `json:"id" db:"id"`
	ProductID         int64     `json:"product_id" db:"product_id"`
	SKU               *string   `json:"sku,omitempty" db:"sku"`
	Barcode           *string   `json:"barcode,omitempty" db:"barcode"`
	Price             *float64  `json:"price,omitempty" db:"price"`
	CompareAtPrice    *float64  `json:"compare_at_price,omitempty" db:"compare_at_price"`
	CostPrice         *float64  `json:"cost_price,omitempty" db:"cost_price"`
	StockQuantity     int32     `json:"stock_quantity" db:"stock_quantity"`
	StockStatus       string    `json:"stock_status" db:"stock_status"` // in_stock, low_stock, out_of_stock
	LowStockThreshold *int32    `json:"low_stock_threshold,omitempty" db:"low_stock_threshold"`
	VariantAttributes string    `json:"variant_attributes" db:"variant_attributes"` // JSONB
	Weight            *float64  `json:"weight,omitempty" db:"weight"`
	Dimensions        *string   `json:"dimensions,omitempty" db:"dimensions"` // JSONB
	IsActive          bool      `json:"is_active" db:"is_active"`
	IsDefault         bool      `json:"is_default" db:"is_default"`
	ViewCount         int32     `json:"view_count" db:"view_count"`
	SoldCount         int32     `json:"sold_count" db:"sold_count"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// VariantUpdate represents fields that can be updated for a variant
// Uses **T pattern for optional nullable fields (nil = no update, *nil = set to NULL, *value = update to value)
type VariantUpdate struct {
	SKU               **string  `json:"sku,omitempty"`
	Barcode           **string  `json:"barcode,omitempty"`
	Price             **float64 `json:"price,omitempty"`
	CompareAtPrice    **float64 `json:"compare_at_price,omitempty"`
	CostPrice         **float64 `json:"cost_price,omitempty"`
	StockQuantity     *int32    `json:"stock_quantity,omitempty"`
	StockStatus       *string   `json:"stock_status,omitempty"`
	LowStockThreshold **int32   `json:"low_stock_threshold,omitempty"`
	VariantAttributes *string   `json:"variant_attributes,omitempty"` // JSON string
	Weight            **float64 `json:"weight,omitempty"`
	Dimensions        **string  `json:"dimensions,omitempty"` // JSON string
	IsActive          *bool     `json:"is_active,omitempty"`
	IsDefault         *bool     `json:"is_default,omitempty"`
}

// VariantFilters represents filters for querying variants
type VariantFilters struct {
	ProductID   int64   `json:"product_id"`
	ActiveOnly  *bool   `json:"active_only,omitempty"`
	StockStatus *string `json:"stock_status,omitempty"`
}
