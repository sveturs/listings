package models

import "time"

// VariantAttributeMapping представляет связь между вариативными атрибутами и категориями
type VariantAttributeMapping struct {
	ID                 int       `json:"id" db:"id"`
	VariantAttributeID int       `json:"variant_attribute_id" db:"variant_attribute_id"`
	CategoryID         int       `json:"category_id" db:"category_id"`
	SortOrder          int       `json:"sort_order" db:"sort_order"`
	IsRequired         bool      `json:"is_required" db:"is_required"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`

	// Связанные объекты
	Attribute *UnifiedAttribute    `json:"attribute,omitempty"`
	Category  *MarketplaceCategory `json:"category,omitempty"`
}

// VariantAttributeMappingFilter фильтр для поиска связей
type VariantAttributeMappingFilter struct {
	CategoryID  *int  `json:"category_id,omitempty"`
	AttributeID *int  `json:"attribute_id,omitempty"`
	IsRequired  *bool `json:"is_required,omitempty"`
	Limit       int   `json:"limit,omitempty"`
	Offset      int   `json:"offset,omitempty"`
}

// VariantAttributeMappingCreateRequest запрос на создание/обновление связи
type VariantAttributeMappingCreateRequest struct {
	VariantAttributeID int  `json:"variant_attribute_id" validate:"required"`
	CategoryID         int  `json:"category_id" validate:"required"`
	SortOrder          int  `json:"sort_order"`
	IsRequired         bool `json:"is_required"`
}

// VariantAttributeMappingUpdateRequest запрос на обновление связи
type VariantAttributeMappingUpdateRequest struct {
	SortOrder  *int  `json:"sort_order,omitempty"`
	IsRequired *bool `json:"is_required,omitempty"`
}

// CategoryVariantAttributesUpdateRequest запрос на обновление вариативных атрибутов категории
type CategoryVariantAttributesUpdateRequest struct {
	CategoryID int                                  `json:"category_id" validate:"required"`
	Attributes []CategoryVariantAttributeAssignment `json:"attributes" validate:"required"`
}

// CategoryVariantAttributeAssignment назначение атрибута категории
type CategoryVariantAttributeAssignment struct {
	AttributeID int  `json:"attribute_id" validate:"required"`
	SortOrder   int  `json:"sort_order"`
	IsRequired  bool `json:"is_required"`
}
