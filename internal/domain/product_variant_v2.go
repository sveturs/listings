// Package domain defines core business entities for the listings microservice.
// This file contains Product Variant V2 domain models with UUID support.
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ProductVariantV2 represents a product variant with UUID-based architecture (Phase 3)
// This replaces the old int64-based ProductVariant for new variant system
type ProductVariantV2 struct {
	ID              uuid.UUID `json:"id" db:"id"`
	ProductID       uuid.UUID `json:"product_id" db:"product_id"`
	SKU             string    `json:"sku" db:"sku"`
	Price           *float64  `json:"price,omitempty" db:"price"` // NULL = use product base_price
	CompareAtPrice  *float64  `json:"compare_at_price,omitempty" db:"compare_at_price"`
	StockQuantity   int32     `json:"stock_quantity" db:"stock_quantity"`
	ReservedQuantity int32    `json:"reserved_quantity" db:"reserved_quantity"`
	LowStockAlert   int32     `json:"low_stock_alert" db:"low_stock_alert"`
	WeightGrams     *float64  `json:"weight_grams,omitempty" db:"weight_grams"`
	Barcode         *string   `json:"barcode,omitempty" db:"barcode"`
	IsDefault       bool      `json:"is_default" db:"is_default"`
	Position        int32     `json:"position" db:"position"`
	Status          string    `json:"status" db:"status"` // active, out_of_stock, discontinued
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`

	// Relations (loaded on demand)
	Attributes []*VariantAttributeValueV2 `json:"attributes,omitempty" db:"-"`
}

// VariantAttributeValueV2 represents an attribute value for a product variant
type VariantAttributeValueV2 struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	VariantID   uuid.UUID  `json:"variant_id" db:"variant_id"`
	AttributeID int32      `json:"attribute_id" db:"attribute_id"`
	ValueText   *string    `json:"value_text,omitempty" db:"value_text"`
	ValueNumber *float64   `json:"value_number,omitempty" db:"value_number"`
	ValueBoolean *bool     `json:"value_boolean,omitempty" db:"value_boolean"`
	ValueDate   *time.Time `json:"value_date,omitempty" db:"value_date"`
	ValueJSON   []byte     `json:"value_json,omitempty" db:"value_json"` // JSONB for multiselect
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`

	// Loaded attribute metadata (optional)
	Attribute *Attribute `json:"attribute,omitempty" db:"-"`
}

// CreateVariantInputV2 represents input for creating a new variant (Phase 3)
type CreateVariantInputV2 struct {
	ProductID      uuid.UUID                    `json:"product_id" validate:"required"`
	SKU            string                       `json:"sku" validate:"required,max=100"`
	Price          *float64                     `json:"price,omitempty" validate:"omitempty,gte=0"`
	CompareAtPrice *float64                     `json:"compare_at_price,omitempty" validate:"omitempty,gte=0"`
	StockQuantity  int32                        `json:"stock_quantity" validate:"gte=0"`
	LowStockAlert  int32                        `json:"low_stock_alert" validate:"gte=0"`
	WeightGrams    *float64                     `json:"weight_grams,omitempty" validate:"omitempty,gte=0"`
	Barcode        *string                      `json:"barcode,omitempty" validate:"omitempty,max=50"`
	IsDefault      bool                         `json:"is_default"`
	Position       int32                        `json:"position" validate:"gte=0"`
	Attributes     []CreateVariantAttributeValue `json:"attributes" validate:"required,min=1"`
}

// CreateVariantAttributeValue represents an attribute value when creating a variant
type CreateVariantAttributeValue struct {
	AttributeID  int32      `json:"attribute_id" validate:"required"`
	ValueText    *string    `json:"value_text,omitempty"`
	ValueNumber  *float64   `json:"value_number,omitempty"`
	ValueBoolean *bool      `json:"value_boolean,omitempty"`
	ValueDate    *time.Time `json:"value_date,omitempty"`
	ValueJSON    []byte     `json:"value_json,omitempty"`
}

// UpdateVariantInputV2 represents input for updating an existing variant (Phase 3)
type UpdateVariantInputV2 struct {
	SKU            *string  `json:"sku,omitempty" validate:"omitempty,max=100"`
	Price          *float64 `json:"price,omitempty" validate:"omitempty,gte=0"`
	CompareAtPrice *float64 `json:"compare_at_price,omitempty" validate:"omitempty,gte=0"`
	StockQuantity  *int32   `json:"stock_quantity,omitempty" validate:"omitempty,gte=0"`
	LowStockAlert  *int32   `json:"low_stock_alert,omitempty" validate:"omitempty,gte=0"`
	WeightGrams    *float64 `json:"weight_grams,omitempty" validate:"omitempty,gte=0"`
	Barcode        *string  `json:"barcode,omitempty" validate:"omitempty,max=50"`
	IsDefault      *bool    `json:"is_default,omitempty"`
	Position       *int32   `json:"position,omitempty" validate:"omitempty,gte=0"`
	Status         *string  `json:"status,omitempty" validate:"omitempty,oneof=active out_of_stock discontinued"`
}

// ListVariantsFilter represents filters for listing variants
type ListVariantsFilter struct {
	ProductID       uuid.UUID `json:"product_id" validate:"required"`
	ActiveOnly      bool      `json:"active_only"`
	InStockOnly     bool      `json:"in_stock_only"`
	IncludeAttributes bool    `json:"include_attributes"`
}

// FindVariantByAttributesFilter represents filters for finding variant by attribute combination
type FindVariantByAttributesFilter struct {
	ProductID  uuid.UUID              `json:"product_id" validate:"required"`
	Attributes map[int32]interface{}  `json:"attributes" validate:"required,min=1"`
}

// Variant status constants
const (
	VariantStatusActive       = "active"
	VariantStatusOutOfStock   = "out_of_stock"
	VariantStatusDiscontinued = "discontinued"
)

// GetAvailableQuantity returns stock_quantity - reserved_quantity
func (v *ProductVariantV2) GetAvailableQuantity() int32 {
	return v.StockQuantity - v.ReservedQuantity
}

// IsAvailable returns true if variant has available stock
func (v *ProductVariantV2) IsAvailable() bool {
	return v.Status == VariantStatusActive && v.GetAvailableQuantity() > 0
}

// IsLowStock returns true if available quantity is at or below low stock alert threshold
func (v *ProductVariantV2) IsLowStock() bool {
	return v.GetAvailableQuantity() <= v.LowStockAlert && v.GetAvailableQuantity() > 0
}

// GetEffectivePrice returns variant price or nil if using product base price
func (v *ProductVariantV2) GetEffectivePrice() *float64 {
	return v.Price
}

// GetAttributeValue returns the value of a specific attribute
func (v *ProductVariantV2) GetAttributeValue(attributeID int32) *VariantAttributeValueV2 {
	if v.Attributes == nil {
		return nil
	}

	for _, attr := range v.Attributes {
		if attr.AttributeID == attributeID {
			return attr
		}
	}

	return nil
}

// HasAttribute checks if variant has a specific attribute
func (v *ProductVariantV2) HasAttribute(attributeID int32) bool {
	return v.GetAttributeValue(attributeID) != nil
}

// MatchesAttributes checks if variant matches all provided attribute values
func (v *ProductVariantV2) MatchesAttributes(filters map[int32]interface{}) bool {
	if v.Attributes == nil || len(v.Attributes) == 0 {
		return false
	}

	for attrID, expectedValue := range filters {
		attr := v.GetAttributeValue(attrID)
		if attr == nil {
			return false
		}

		// Compare based on value type
		if attr.ValueText != nil && expectedValue != nil {
			if strVal, ok := expectedValue.(string); ok {
				if *attr.ValueText != strVal {
					return false
				}
			} else {
				return false
			}
		} else if attr.ValueNumber != nil && expectedValue != nil {
			if numVal, ok := expectedValue.(float64); ok {
				if *attr.ValueNumber != numVal {
					return false
				}
			} else {
				return false
			}
		} else if attr.ValueBoolean != nil && expectedValue != nil {
			if boolVal, ok := expectedValue.(bool); ok {
				if *attr.ValueBoolean != boolVal {
					return false
				}
			} else {
				return false
			}
		}
	}

	return true
}

// Validate performs basic validation on the variant
func (v *ProductVariantV2) Validate() error {
	if v.SKU == "" {
		return ErrInvalidSKU
	}

	if v.StockQuantity < 0 {
		return ErrInvalidStockQuantity
	}

	if v.ReservedQuantity < 0 {
		return ErrInvalidReservedQuantity
	}

	if v.ReservedQuantity > v.StockQuantity {
		return ErrReservedExceedsStock
	}

	if v.Price != nil && *v.Price < 0 {
		return ErrInvalidPrice
	}

	return nil
}

// Domain errors for ProductVariantV2
var (
	ErrVariantNotFound          = errors.New("variant not found")
	ErrInvalidSKU              = errors.New("invalid or empty SKU")
	ErrDuplicateSKU            = errors.New("SKU already exists")
	ErrInvalidStockQuantity    = errors.New("stock quantity cannot be negative")
	ErrInvalidReservedQuantity = errors.New("reserved quantity cannot be negative")
	ErrReservedExceedsStock    = errors.New("reserved quantity cannot exceed stock quantity")
	ErrInvalidPrice            = errors.New("price cannot be negative")
	ErrInsufficientStock       = errors.New("insufficient stock for operation")
	ErrVariantAttributeNotFound = errors.New("variant attribute value not found")
)
