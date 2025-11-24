package domain

import (
	"time"
)

// CategoryDetail represents a complete category with all fields (for gRPC)
type CategoryDetail struct {
	ID                int32     `db:"id" json:"id"`
	Name              string    `db:"name" json:"name"`
	Slug              string    `db:"slug" json:"slug"`
	ParentID          *int32    `db:"parent_id" json:"parent_id,omitempty"`
	Icon              *string   `db:"icon" json:"icon,omitempty"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	HasCustomUI       bool      `db:"has_custom_ui" json:"has_custom_ui"`
	CustomUIComponent *string   `db:"custom_ui_component" json:"custom_ui_component,omitempty"`
	SortOrder         int32     `db:"sort_order" json:"sort_order"`
	Level             int32     `db:"level" json:"level"`
	Count             int32     `db:"count" json:"count"`
	ExternalID        *string   `db:"external_id" json:"external_id,omitempty"`
	Description       *string   `db:"description" json:"description,omitempty"`
	IsActive          bool      `db:"is_active" json:"is_active"`
	SEOTitle          *string   `db:"seo_title" json:"seo_title,omitempty"`
	SEODescription    *string   `db:"seo_description" json:"seo_description,omitempty"`
	SEOKeywords       *string   `db:"seo_keywords" json:"seo_keywords,omitempty"`
	TitleEn           *string   `db:"title_en" json:"title_en,omitempty"`
	TitleRu           *string   `db:"title_ru" json:"title_ru,omitempty"`
	TitleSr           *string   `db:"title_sr" json:"title_sr,omitempty"`
}

// CategoryTree represents a category with its subcategories
type CategoryTree struct {
	Category      *Category       `json:"category"`
	Subcategories []*CategoryTree `json:"subcategories"`
}

// CreateCategoryInput contains fields for creating a new category
type CreateCategoryInput struct {
	Name              string  `json:"name" validate:"required,min=1,max=100"`
	Slug              string  `json:"slug" validate:"required,min=1,max=100"`
	ParentID          *int32  `json:"parent_id,omitempty"`
	Icon              *string `json:"icon,omitempty"`
	Description       *string `json:"description,omitempty"`
	CustomUIComponent *string `json:"custom_ui_component,omitempty"`
	SortOrder         int32   `json:"sort_order"`
	SEOTitle          *string `json:"seo_title,omitempty"`
	SEODescription    *string `json:"seo_description,omitempty"`
	SEOKeywords       *string `json:"seo_keywords,omitempty"`
	TitleEn           *string `json:"title_en,omitempty"`
	TitleRu           *string `json:"title_ru,omitempty"`
	TitleSr           *string `json:"title_sr,omitempty"`
}

// UpdateCategoryInput contains fields for updating an existing category
type UpdateCategoryInput struct {
	Name              *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Slug              *string `json:"slug,omitempty" validate:"omitempty,min=1,max=100"`
	ParentID          *int32  `json:"parent_id,omitempty"`
	Icon              *string `json:"icon,omitempty"`
	Description       *string `json:"description,omitempty"`
	CustomUIComponent *string `json:"custom_ui_component,omitempty"`
	SortOrder         *int32  `json:"sort_order,omitempty"`
	IsActive          *bool   `json:"is_active,omitempty"`
	SEOTitle          *string `json:"seo_title,omitempty"`
	SEODescription    *string `json:"seo_description,omitempty"`
	SEOKeywords       *string `json:"seo_keywords,omitempty"`
	TitleEn           *string `json:"title_en,omitempty"`
	TitleRu           *string `json:"title_ru,omitempty"`
	TitleSr           *string `json:"title_sr,omitempty"`
}

// GetCategoriesFilter contains filter options for listing categories
type GetCategoriesFilter struct {
	ParentID        *int32 `json:"parent_id,omitempty"`
	IsActive        *bool  `json:"is_active,omitempty"`
	IncludeChildren bool   `json:"include_children"`
	Page            int32  `json:"page"`
	PageSize        int32  `json:"page_size"`
}

// GetCategoryTreeFilter contains filter options for getting category tree
type GetCategoryTreeFilter struct {
	RootID     *int32 `json:"root_id,omitempty"`
	ActiveOnly bool   `json:"active_only"`
}
