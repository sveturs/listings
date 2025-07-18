package types

import (
	"encoding/json"
	"time"
)

// ProductVariantAttribute represents a variant attribute (color, size, etc.)
type ProductVariantAttribute struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	DisplayName string    `json:"display_name" db:"display_name"`
	Type        string    `json:"type" db:"type"` // text, color, image, number
	IsRequired  bool      `json:"is_required" db:"is_required"`
	SortOrder   int       `json:"sort_order" db:"sort_order"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// ProductVariantAttributeValue represents possible values for an attribute
type ProductVariantAttributeValue struct {
	ID          int       `json:"id" db:"id"`
	AttributeID int       `json:"attribute_id" db:"attribute_id"`
	Value       string    `json:"value" db:"value"`
	DisplayName string    `json:"display_name" db:"display_name"`
	ColorHex    *string   `json:"color_hex,omitempty" db:"color_hex"`
	ImageURL    *string   `json:"image_url,omitempty" db:"image_url"`
	SortOrder   int       `json:"sort_order" db:"sort_order"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// ProductVariant represents a specific variant of a product
type ProductVariant struct {
	ID                int                    `json:"id" db:"id"`
	ProductID         int                    `json:"product_id" db:"product_id"`
	SKU               *string                `json:"sku,omitempty" db:"sku"`
	Barcode           *string                `json:"barcode,omitempty" db:"barcode"`
	Price             *float64               `json:"price,omitempty" db:"price"`
	CompareAtPrice    *float64               `json:"compare_at_price,omitempty" db:"compare_at_price"`
	CostPrice         *float64               `json:"cost_price,omitempty" db:"cost_price"`
	StockQuantity     int                    `json:"stock_quantity" db:"stock_quantity"`
	StockStatus       string                 `json:"stock_status" db:"stock_status"`
	LowStockThreshold *int                   `json:"low_stock_threshold,omitempty" db:"low_stock_threshold"`
	VariantAttributes map[string]interface{} `json:"variant_attributes" db:"variant_attributes"`
	Weight            *float64               `json:"weight,omitempty" db:"weight"`
	Dimensions        map[string]interface{} `json:"dimensions,omitempty" db:"dimensions"`
	IsActive          bool                   `json:"is_active" db:"is_active"`
	IsDefault         bool                   `json:"is_default" db:"is_default"`
	ViewCount         int                    `json:"view_count" db:"view_count"`
	SoldCount         int                    `json:"sold_count" db:"sold_count"`
	CreatedAt         time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at" db:"updated_at"`

	// Related data
	Images []ProductVariantImage `json:"images,omitempty"`
}

// ProductVariantImage represents an image for a specific variant
type ProductVariantImage struct {
	ID           int       `json:"id" db:"id"`
	VariantID    int       `json:"variant_id" db:"variant_id"`
	ImageURL     string    `json:"image_url" db:"image_url"`
	ThumbnailURL *string   `json:"thumbnail_url,omitempty" db:"thumbnail_url"`
	AltText      *string   `json:"alt_text,omitempty" db:"alt_text"`
	DisplayOrder int       `json:"display_order" db:"display_order"`
	IsMain       bool      `json:"is_main" db:"is_main"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// CreateVariantRequest represents request to create a new variant
type CreateVariantRequest struct {
	ProductID         int                         `json:"product_id" validate:"required"`
	SKU               *string                     `json:"sku,omitempty"`
	Barcode           *string                     `json:"barcode,omitempty"`
	Price             *float64                    `json:"price,omitempty"`
	CompareAtPrice    *float64                    `json:"compare_at_price,omitempty"`
	CostPrice         *float64                    `json:"cost_price,omitempty"`
	StockQuantity     int                         `json:"stock_quantity" validate:"min=0"`
	LowStockThreshold *int                        `json:"low_stock_threshold,omitempty"`
	VariantAttributes map[string]interface{}      `json:"variant_attributes" validate:"required"`
	Weight            *float64                    `json:"weight,omitempty"`
	Dimensions        map[string]interface{}      `json:"dimensions,omitempty"`
	IsDefault         bool                        `json:"is_default"`
	Images            []CreateVariantImageRequest `json:"images,omitempty"`
}

// CreateVariantImageRequest represents request to create variant image
type CreateVariantImageRequest struct {
	ImageURL     string  `json:"image_url" validate:"required"`
	ThumbnailURL *string `json:"thumbnail_url,omitempty"`
	AltText      *string `json:"alt_text,omitempty"`
	DisplayOrder int     `json:"display_order"`
	IsMain       bool    `json:"is_main"`
}

// UpdateVariantRequest represents request to update a variant
type UpdateVariantRequest struct {
	SKU               *string                `json:"sku,omitempty"`
	Barcode           *string                `json:"barcode,omitempty"`
	Price             *float64               `json:"price,omitempty"`
	CompareAtPrice    *float64               `json:"compare_at_price,omitempty"`
	CostPrice         *float64               `json:"cost_price,omitempty"`
	StockQuantity     *int                   `json:"stock_quantity,omitempty" validate:"omitempty,min=0"`
	LowStockThreshold *int                   `json:"low_stock_threshold,omitempty"`
	VariantAttributes map[string]interface{} `json:"variant_attributes,omitempty"`
	Weight            *float64               `json:"weight,omitempty"`
	Dimensions        map[string]interface{} `json:"dimensions,omitempty"`
	IsActive          *bool                  `json:"is_active,omitempty"`
	IsDefault         *bool                  `json:"is_default,omitempty"`
}

// VariantSearchRequest represents search parameters for variants
type VariantSearchRequest struct {
	ProductID         *int                   `json:"product_id,omitempty"`
	SKU               *string                `json:"sku,omitempty"`
	StockStatus       *string                `json:"stock_status,omitempty"`
	IsActive          *bool                  `json:"is_active,omitempty"`
	VariantAttributes map[string]interface{} `json:"variant_attributes,omitempty"`
	Page              int                    `json:"page" validate:"min=1"`
	Limit             int                    `json:"limit" validate:"min=1,max=100"`
}

// VariantSearchResponse represents search results
type VariantSearchResponse struct {
	Variants   []ProductVariant `json:"variants"`
	Total      int              `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"total_pages"`
}

// BulkCreateVariantsRequest represents request to create multiple variants
type BulkCreateVariantsRequest struct {
	ProductID int                    `json:"product_id" validate:"required"`
	Variants  []CreateVariantRequest `json:"variants" validate:"required,min=1"`
}

// StorefrontProductAttribute represents seller's custom attribute configuration
type StorefrontProductAttribute struct {
	ID           int              `json:"id" db:"id"`
	ProductID    int              `json:"product_id" db:"product_id"`
	AttributeID  int              `json:"attribute_id" db:"attribute_id"`
	IsEnabled    bool             `json:"is_enabled" db:"is_enabled"`
	IsRequired   bool             `json:"is_required" db:"is_required"`
	CustomValues []AttributeValue `json:"custom_values" db:"custom_values"`
	CreatedAt    time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at" db:"updated_at"`

	// Related data
	Attribute    *ProductVariantAttribute       `json:"attribute,omitempty"`
	GlobalValues []ProductVariantAttributeValue `json:"global_values,omitempty"`
}

// AttributeValue represents a single attribute value (global or custom)
type AttributeValue struct {
	Value       string  `json:"value"`
	DisplayName string  `json:"display_name"`
	ColorHex    *string `json:"color_hex,omitempty"`
	ImageURL    *string `json:"image_url,omitempty"`
	IsCustom    bool    `json:"is_custom"` // true if added by seller
}

// GenerateVariantsRequest represents request to auto-generate variants
type GenerateVariantsRequest struct {
	ProductID         int                                    `json:"product_id" validate:"required"`
	AttributeMatrix   map[string][]string                    `json:"attribute_matrix" validate:"required"` // {"color": ["red", "blue"], "brand": ["Samsung", "Apple"]}
	PriceModifiers    map[string]float64                     `json:"price_modifiers,omitempty"`            // {"256gb": 50.0, "Apple": 200.0}
	StockQuantities   map[string]int                         `json:"stock_quantities,omitempty"`           // {"Samsung-Galaxy_A54-red": 10}
	DefaultAttributes map[string]string                      `json:"default_attributes,omitempty"`         // {"color": "black", "brand": "Samsung"}
	ImageMappings     map[string][]CreateVariantImageRequest `json:"image_mappings,omitempty"`             // {"red": [images], "Samsung": [images]}
}

// SetupProductAttributesRequest represents request to configure product attributes
type SetupProductAttributesRequest struct {
	ProductID  int                     `json:"product_id" validate:"required"`
	Attributes []ProductAttributeSetup `json:"attributes" validate:"required"`
}

// ProductAttributeSetup represents attribute configuration for a product
type ProductAttributeSetup struct {
	AttributeID          int              `json:"attribute_id" validate:"required"`
	IsEnabled            bool             `json:"is_enabled"`
	IsRequired           bool             `json:"is_required"`
	CustomValues         []AttributeValue `json:"custom_values,omitempty"`          // seller's custom values
	SelectedGlobalValues []string         `json:"selected_global_values,omitempty"` // selected from global values
}

// Custom JSON marshaling for JSONB fields
func (pv *ProductVariant) MarshalJSON() ([]byte, error) {
	type Alias ProductVariant
	return json.Marshal(&struct {
		*Alias
		VariantAttributes json.RawMessage `json:"variant_attributes"`
		Dimensions        json.RawMessage `json:"dimensions,omitempty"`
	}{
		Alias:             (*Alias)(pv),
		VariantAttributes: mustMarshalJSON(pv.VariantAttributes),
		Dimensions:        mustMarshalJSON(pv.Dimensions),
	})
}

func mustMarshalJSON(v interface{}) json.RawMessage {
	if v == nil {
		return json.RawMessage("{}")
	}
	data, _ := json.Marshal(v)
	return data
}
