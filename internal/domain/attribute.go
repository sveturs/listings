// Package domain defines core business entities and domain models for the listings microservice.
// This file contains attribute-related domain models.
package domain

import (
	"encoding/json"
	"fmt"
	"time"
)

// AttributeType represents the type of an attribute
type AttributeType string

const (
	AttributeTypeText        AttributeType = "text"
	AttributeTypeTextarea    AttributeType = "textarea"
	AttributeTypeNumber      AttributeType = "number"
	AttributeTypeBoolean     AttributeType = "boolean"
	AttributeTypeSelect      AttributeType = "select"
	AttributeTypeMultiselect AttributeType = "multiselect"
	AttributeTypeDate        AttributeType = "date"
	AttributeTypeColor       AttributeType = "color"
	AttributeTypeSize        AttributeType = "size"
)

// AttributePurpose represents the purpose of an attribute
type AttributePurpose string

const (
	AttributePurposeRegular AttributePurpose = "regular"
	AttributePurposeVariant AttributePurpose = "variant"
	AttributePurposeBoth    AttributePurpose = "both"
)

// AttributeOption represents an option for select/multiselect attributes
type AttributeOption struct {
	Value string            `json:"value"`
	Label map[string]string `json:"label"` // i18n: {"en": "...", "ru": "...", "sr": "..."}
}

// Attribute represents an attribute metadata entity
type Attribute struct {
	ID                  int32                  `json:"id" db:"id"`
	Code                string                 `json:"code" db:"code"`
	Name                map[string]string      `json:"name" db:"name"`                 // i18n JSONB
	DisplayName         map[string]string      `json:"display_name" db:"display_name"` // i18n JSONB
	AttributeType       AttributeType          `json:"attribute_type" db:"attribute_type"`
	Purpose             AttributePurpose       `json:"purpose" db:"purpose"`
	Options             []AttributeOption      `json:"options,omitempty" db:"options"`                   // JSONB
	ValidationRules     map[string]interface{} `json:"validation_rules,omitempty" db:"validation_rules"` // JSONB
	UISettings          map[string]interface{} `json:"ui_settings,omitempty" db:"ui_settings"`           // JSONB
	IsSearchable        bool                   `json:"is_searchable" db:"is_searchable"`
	IsFilterable        bool                   `json:"is_filterable" db:"is_filterable"`
	IsRequired          bool                   `json:"is_required" db:"is_required"`
	IsVariantCompatible bool                   `json:"is_variant_compatible" db:"is_variant_compatible"`
	AffectsStock        bool                   `json:"affects_stock" db:"affects_stock"`
	AffectsPrice        bool                   `json:"affects_price" db:"affects_price"`
	ShowInCard          bool                   `json:"show_in_card" db:"show_in_card"`
	IsActive            bool                   `json:"is_active" db:"is_active"`
	SortOrder           int32                  `json:"sort_order" db:"sort_order"`
	Icon                string                 `json:"icon,omitempty" db:"icon"`
	CreatedAt           time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at" db:"updated_at"`
}

// CategoryAttribute represents a category-attribute relationship with overrides
type CategoryAttribute struct {
	ID                      int32                  `json:"id" db:"id"`
	CategoryID              int32                  `json:"category_id" db:"category_id"`
	AttributeID             int32                  `json:"attribute_id" db:"attribute_id"`
	Attribute               *Attribute             `json:"attribute,omitempty" db:"-"` // Loaded on demand
	IsEnabled               bool                   `json:"is_enabled" db:"is_enabled"`
	IsRequired              *bool                  `json:"is_required,omitempty" db:"is_required"`     // Nullable override
	IsSearchable            *bool                  `json:"is_searchable,omitempty" db:"is_searchable"` // Nullable override
	IsFilterable            *bool                  `json:"is_filterable,omitempty" db:"is_filterable"` // Nullable override
	SortOrder               int32                  `json:"sort_order" db:"sort_order"`
	CategorySpecificOptions []AttributeOption      `json:"category_specific_options,omitempty" db:"category_specific_options"` // JSONB
	CustomValidationRules   map[string]interface{} `json:"custom_validation_rules,omitempty" db:"custom_validation_rules"`     // JSONB
	CustomUISettings        map[string]interface{} `json:"custom_ui_settings,omitempty" db:"custom_ui_settings"`               // JSONB
	IsActive                bool                   `json:"is_active" db:"is_active"`
	CreatedAt               time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt               time.Time              `json:"updated_at" db:"updated_at"`
}

// ListingAttributeValue represents an attribute value for a listing
type ListingAttributeValue struct {
	ID           int32                  `json:"id" db:"id"`
	ListingID    int32                  `json:"listing_id" db:"listing_id"`
	AttributeID  int32                  `json:"attribute_id" db:"attribute_id"`
	Attribute    *Attribute             `json:"attribute,omitempty" db:"-"` // Loaded on demand
	ValueText    *string                `json:"value_text,omitempty" db:"value_text"`
	ValueNumber  *float64               `json:"value_number,omitempty" db:"value_number"`
	ValueBoolean *bool                  `json:"value_boolean,omitempty" db:"value_boolean"`
	ValueDate    *time.Time             `json:"value_date,omitempty" db:"value_date"`
	ValueJSON    map[string]interface{} `json:"value_json,omitempty" db:"value_json"` // For multiselect, complex objects
	CreatedAt    time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at" db:"updated_at"`
}

// VariantAttribute represents a variant attribute definition for a category
type VariantAttribute struct {
	ID           int32      `json:"id" db:"id"`
	CategoryID   int32      `json:"category_id" db:"category_id"`
	AttributeID  int32      `json:"attribute_id" db:"attribute_id"`
	Attribute    *Attribute `json:"attribute,omitempty" db:"-"` // Loaded on demand
	IsRequired   bool       `json:"is_required" db:"is_required"`
	AffectsPrice bool       `json:"affects_price" db:"affects_price"`
	AffectsStock bool       `json:"affects_stock" db:"affects_stock"`
	SortOrder    int32      `json:"sort_order" db:"sort_order"`
	DisplayAs    string     `json:"display_as" db:"display_as"` // dropdown, buttons, swatches, radio
	IsActive     bool       `json:"is_active" db:"is_active"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// VariantAttributeValue represents an attribute value for a product variant
type VariantAttributeValue struct {
	ID                int32                  `json:"id" db:"id"`
	VariantID         int32                  `json:"variant_id" db:"variant_id"`
	AttributeID       int32                  `json:"attribute_id" db:"attribute_id"`
	Attribute         *Attribute             `json:"attribute,omitempty" db:"-"` // Loaded on demand
	ValueText         *string                `json:"value_text,omitempty" db:"value_text"`
	ValueNumber       *float64               `json:"value_number,omitempty" db:"value_number"`
	ValueBoolean      *bool                  `json:"value_boolean,omitempty" db:"value_boolean"`
	ValueDate         *time.Time             `json:"value_date,omitempty" db:"value_date"`
	ValueJSON         map[string]interface{} `json:"value_json,omitempty" db:"value_json"`
	PriceModifier     float64                `json:"price_modifier" db:"price_modifier"`           // For price adjustments
	PriceModifierType string                 `json:"price_modifier_type" db:"price_modifier_type"` // fixed, percent
	CreatedAt         time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at" db:"updated_at"`
}

// AttributeOption represents a predefined option for select/multiselect attributes
type AttributeOptionEntity struct {
	ID          int32             `json:"id" db:"id"`
	AttributeID int32             `json:"attribute_id" db:"attribute_id"`
	OptionValue string            `json:"option_value" db:"option_value"`
	OptionLabel map[string]string `json:"option_label" db:"option_label"`     // i18n JSONB
	ColorHex    *string           `json:"color_hex,omitempty" db:"color_hex"` // For color swatches
	ImageURL    *string           `json:"image_url,omitempty" db:"image_url"` // For image-based options
	Icon        *string           `json:"icon,omitempty" db:"icon"`           // For icon-based options
	IsDefault   bool              `json:"is_default" db:"is_default"`
	IsActive    bool              `json:"is_active" db:"is_active"`
	SortOrder   int32             `json:"sort_order" db:"sort_order"`
	CreatedAt   time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" db:"updated_at"`
}

// CreateAttributeInput represents input for creating a new attribute
type CreateAttributeInput struct {
	Code                string                 `json:"code" validate:"required,max=100"`
	Name                map[string]string      `json:"name" validate:"required"`
	DisplayName         map[string]string      `json:"display_name" validate:"required"`
	AttributeType       AttributeType          `json:"attribute_type" validate:"required"`
	Purpose             AttributePurpose       `json:"purpose"`
	Options             []AttributeOption      `json:"options,omitempty"`
	ValidationRules     map[string]interface{} `json:"validation_rules,omitempty"`
	UISettings          map[string]interface{} `json:"ui_settings,omitempty"`
	IsSearchable        bool                   `json:"is_searchable"`
	IsFilterable        bool                   `json:"is_filterable"`
	IsRequired          bool                   `json:"is_required"`
	IsVariantCompatible bool                   `json:"is_variant_compatible"`
	AffectsStock        bool                   `json:"affects_stock"`
	AffectsPrice        bool                   `json:"affects_price"`
	ShowInCard          bool                   `json:"show_in_card"`
	SortOrder           int32                  `json:"sort_order"`
	Icon                string                 `json:"icon,omitempty"`
}

// UpdateAttributeInput represents input for updating an existing attribute
type UpdateAttributeInput struct {
	Name                *map[string]string      `json:"name,omitempty"`
	DisplayName         *map[string]string      `json:"display_name,omitempty"`
	AttributeType       *AttributeType          `json:"attribute_type,omitempty"`
	Purpose             *AttributePurpose       `json:"purpose,omitempty"`
	Options             *[]AttributeOption      `json:"options,omitempty"`
	ValidationRules     *map[string]interface{} `json:"validation_rules,omitempty"`
	UISettings          *map[string]interface{} `json:"ui_settings,omitempty"`
	IsSearchable        *bool                   `json:"is_searchable,omitempty"`
	IsFilterable        *bool                   `json:"is_filterable,omitempty"`
	IsRequired          *bool                   `json:"is_required,omitempty"`
	IsVariantCompatible *bool                   `json:"is_variant_compatible,omitempty"`
	AffectsStock        *bool                   `json:"affects_stock,omitempty"`
	AffectsPrice        *bool                   `json:"affects_price,omitempty"`
	ShowInCard          *bool                   `json:"show_in_card,omitempty"`
	IsActive            *bool                   `json:"is_active,omitempty"`
	SortOrder           *int32                  `json:"sort_order,omitempty"`
	Icon                *string                 `json:"icon,omitempty"`
}

// ListAttributesFilter represents filters for listing attributes
type ListAttributesFilter struct {
	AttributeType       *AttributeType    `json:"attribute_type,omitempty"`
	Purpose             *AttributePurpose `json:"purpose,omitempty"`
	IsSearchable        *bool             `json:"is_searchable,omitempty"`
	IsFilterable        *bool             `json:"is_filterable,omitempty"`
	IsVariantCompatible *bool             `json:"is_variant_compatible,omitempty"`
	IsActive            *bool             `json:"is_active,omitempty"`
	Limit               int32             `json:"limit" validate:"gte=1,lte=100"`
	Offset              int32             `json:"offset" validate:"gte=0"`
}

// CategoryAttributeSettings represents settings for linking an attribute to a category
type CategoryAttributeSettings struct {
	IsEnabled               bool                    `json:"is_enabled"`
	IsRequired              *bool                   `json:"is_required,omitempty"`
	IsSearchable            *bool                   `json:"is_searchable,omitempty"`
	IsFilterable            *bool                   `json:"is_filterable,omitempty"`
	SortOrder               int32                   `json:"sort_order"`
	CategorySpecificOptions *[]AttributeOption      `json:"category_specific_options,omitempty"`
	CustomValidationRules   *map[string]interface{} `json:"custom_validation_rules,omitempty"`
	CustomUISettings        *map[string]interface{} `json:"custom_ui_settings,omitempty"`
}

// GetCategoryAttributesFilter represents filters for getting category attributes
type GetCategoryAttributesFilter struct {
	IsEnabled    *bool `json:"is_enabled,omitempty"`
	IsRequired   *bool `json:"is_required,omitempty"`
	IsSearchable *bool `json:"is_searchable,omitempty"`
	IsFilterable *bool `json:"is_filterable,omitempty"`
}

// SetListingAttributeValue represents a single attribute value to set for a listing
type SetListingAttributeValue struct {
	AttributeID  int32                  `json:"attribute_id" validate:"required"`
	ValueText    *string                `json:"value_text,omitempty"`
	ValueNumber  *float64               `json:"value_number,omitempty"`
	ValueBoolean *bool                  `json:"value_boolean,omitempty"`
	ValueDate    *time.Time             `json:"value_date,omitempty"`
	ValueJSON    map[string]interface{} `json:"value_json,omitempty"`
}

// GetEffectiveIsRequired returns the effective is_required value (override or attribute default)
func (ca *CategoryAttribute) GetEffectiveIsRequired() bool {
	if ca.IsRequired != nil {
		return *ca.IsRequired
	}
	if ca.Attribute != nil {
		return ca.Attribute.IsRequired
	}
	return false
}

// GetEffectiveIsSearchable returns the effective is_searchable value
func (ca *CategoryAttribute) GetEffectiveIsSearchable() bool {
	if ca.IsSearchable != nil {
		return *ca.IsSearchable
	}
	if ca.Attribute != nil {
		return ca.Attribute.IsSearchable
	}
	return false
}

// GetEffectiveIsFilterable returns the effective is_filterable value
func (ca *CategoryAttribute) GetEffectiveIsFilterable() bool {
	if ca.IsFilterable != nil {
		return *ca.IsFilterable
	}
	if ca.Attribute != nil {
		return ca.Attribute.IsFilterable
	}
	return false
}

// GetEffectiveOptions returns the effective options (category-specific or attribute default)
func (ca *CategoryAttribute) GetEffectiveOptions() []AttributeOption {
	if len(ca.CategorySpecificOptions) > 0 {
		return ca.CategorySpecificOptions
	}
	if ca.Attribute != nil {
		return ca.Attribute.Options
	}
	return nil
}

// GetEffectiveValidationRules returns the effective validation rules (custom or attribute default)
func (ca *CategoryAttribute) GetEffectiveValidationRules() map[string]interface{} {
	if len(ca.CustomValidationRules) > 0 {
		return ca.CustomValidationRules
	}
	if ca.Attribute != nil {
		return ca.Attribute.ValidationRules
	}
	return nil
}

// GetEffectiveUISettings returns the effective UI settings (custom or attribute default)
func (ca *CategoryAttribute) GetEffectiveUISettings() map[string]interface{} {
	if len(ca.CustomUISettings) > 0 {
		return ca.CustomUISettings
	}
	if ca.Attribute != nil {
		return ca.Attribute.UISettings
	}
	return nil
}

// GetValueAsString returns the attribute value as a string for display
func (lav *ListingAttributeValue) GetValueAsString() string {
	if lav.ValueText != nil {
		return *lav.ValueText
	}
	if lav.ValueNumber != nil {
		return fmt.Sprintf("%.2f", *lav.ValueNumber)
	}
	if lav.ValueBoolean != nil {
		if *lav.ValueBoolean {
			return "true"
		}
		return "false"
	}
	if lav.ValueDate != nil {
		return lav.ValueDate.Format("2006-01-02")
	}
	if len(lav.ValueJSON) > 0 {
		// For multiselect or complex objects
		jsonBytes, _ := json.Marshal(lav.ValueJSON)
		return string(jsonBytes)
	}
	return ""
}
