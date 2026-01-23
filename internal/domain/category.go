package domain

import (
	"time"

	"github.com/google/uuid"
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

// ============================================================================
// V2 Types - UUID-based categories with JSONB support
// ============================================================================

// CategoryV2 represents a category with UUID and JSONB multilingual fields
type CategoryV2 struct {
	ID              uuid.UUID         `json:"id" db:"id"`
	Slug            string            `json:"slug" db:"slug"`
	ParentID        *uuid.UUID        `json:"parent_id,omitempty" db:"parent_id"`
	Level           int32             `json:"level" db:"level"`
	Path            string            `json:"path" db:"path"`
	SortOrder       int32             `json:"sort_order" db:"sort_order"`
	Name            map[string]string `json:"name" db:"name"`
	Description     map[string]string `json:"description,omitempty" db:"description"`
	MetaTitle       map[string]string `json:"meta_title,omitempty" db:"meta_title"`
	MetaDescription map[string]string `json:"meta_description,omitempty" db:"meta_description"`
	MetaKeywords    map[string]string `json:"meta_keywords,omitempty" db:"meta_keywords"`
	Icon            *string           `json:"icon,omitempty" db:"icon"`
	ImageURL        *string           `json:"image_url,omitempty" db:"image_url"`
	IsActive        bool              `json:"is_active" db:"is_active"`
	CreatedAt       time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at" db:"updated_at"`
}

// LocalizedCategory represents a category with localized (single language) fields
type LocalizedCategory struct {
	ID              uuid.UUID  `json:"id"`
	Slug            string     `json:"slug"`
	ParentID        *uuid.UUID `json:"parent_id,omitempty"`
	Level           int32      `json:"level"`
	Path            string     `json:"path"`
	SortOrder       int32      `json:"sort_order"`
	Name            string     `json:"name"`
	Description     string     `json:"description,omitempty"`
	MetaTitle       string     `json:"meta_title,omitempty"`
	MetaDescription string     `json:"meta_description,omitempty"`
	MetaKeywords    string     `json:"meta_keywords,omitempty"`
	Icon            *string    `json:"icon,omitempty"`
	ImageURL        *string    `json:"image_url,omitempty"`
	IsActive        bool       `json:"is_active"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// CategoryTreeV2 represents a category tree with V2 categories
type CategoryTreeV2 struct {
	Category      *LocalizedCategory `json:"category"`
	Subcategories []*CategoryTreeV2  `json:"subcategories,omitempty"`
}

// CategoryBreadcrumb represents a breadcrumb item for navigation
type CategoryBreadcrumb struct {
	ID    uuid.UUID `json:"id"`
	Slug  string    `json:"slug"`
	Name  string    `json:"name"` // Localized name
	Level int32     `json:"level"`
}

// GetCategoryTreeFilterV2 contains filter options for V2 category tree
type GetCategoryTreeFilterV2 struct {
	RootID     *uuid.UUID `json:"root_id,omitempty"`
	Locale     string     `json:"locale"` // sr, en, ru
	ActiveOnly bool       `json:"active_only"`
	MaxDepth   *int32     `json:"max_depth,omitempty"`
}

// Localize extracts localized values from JSONB fields
func (c *CategoryV2) Localize(locale string) *LocalizedCategory {
	return &LocalizedCategory{
		ID:              c.ID,
		Slug:            c.Slug,
		ParentID:        c.ParentID,
		Level:           c.Level,
		Path:            c.Path,
		SortOrder:       c.SortOrder,
		Name:            getLocalized(c.Name, locale),
		Description:     getLocalized(c.Description, locale),
		MetaTitle:       getLocalized(c.MetaTitle, locale),
		MetaDescription: getLocalized(c.MetaDescription, locale),
		MetaKeywords:    getLocalized(c.MetaKeywords, locale),
		Icon:            c.Icon,
		ImageURL:        c.ImageURL,
		IsActive:        c.IsActive,
		CreatedAt:       c.CreatedAt,
		UpdatedAt:       c.UpdatedAt,
	}
}

// getLocalized extracts a localized string from a map with fallback logic
// Priority: requested locale -> "sr" (default) -> "en" -> first available -> empty string
func getLocalized(m map[string]string, locale string) string {
	if m == nil {
		return ""
	}

	// Try requested locale
	if val, ok := m[locale]; ok && val != "" {
		return val
	}

	// Fallback to Serbian (default)
	if locale != "sr" {
		if val, ok := m["sr"]; ok && val != "" {
			return val
		}
	}

	// Fallback to English
	if locale != "en" {
		if val, ok := m["en"]; ok && val != "" {
			return val
		}
	}

	// Return first available value
	for _, val := range m {
		if val != "" {
			return val
		}
	}

	return ""
}
