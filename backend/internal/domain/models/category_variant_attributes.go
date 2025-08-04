package models

import "time"

// CategoryVariantAttribute представляет связь между категорией и доступными для нее вариативными атрибутами
type CategoryVariantAttribute struct {
	ID                   int       `json:"id" db:"id"`
	CategoryID           int       `json:"category_id" db:"category_id"`
	VariantAttributeName string    `json:"variant_attribute_name" db:"variant_attribute_name"`
	SortOrder            int       `json:"sort_order" db:"sort_order"`
	IsRequired           bool      `json:"is_required" db:"is_required"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`

	// Вложенные данные для удобства
	VariantAttribute *ProductVariantAttribute `json:"variant_attribute,omitempty" db:"-"`
}

// CategoryVariantAttributesRequest представляет запрос на обновление вариативных атрибутов категории
type CategoryVariantAttributesRequest struct {
	VariantAttributes []CategoryVariantAttributeUpdate `json:"variant_attributes"`
}

// CategoryVariantAttributeUpdate представляет обновление одного вариативного атрибута
type CategoryVariantAttributeUpdate struct {
	VariantAttributeName string `json:"variant_attribute_name"`
	SortOrder            int    `json:"sort_order"`
	IsRequired           bool   `json:"is_required"`
}
