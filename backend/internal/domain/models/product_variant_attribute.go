package models

import "time"

// ProductVariantAttribute представляет атрибут для вариантов товаров
type ProductVariantAttribute struct {
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	DisplayName  string    `json:"display_name" db:"display_name"`
	Type         string    `json:"type" db:"type"`
	IsRequired   bool      `json:"is_required" db:"is_required"`
	SortOrder    int       `json:"sort_order" db:"sort_order"`
	AffectsStock bool      `json:"affects_stock" db:"affects_stock"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
