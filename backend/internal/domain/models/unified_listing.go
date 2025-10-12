// backend/internal/domain/models/unified_listing.go
package models

import (
	"encoding/json"
	"time"
)

// UnifiedListing объединяет C2C listings и B2C products без дублирования данных
type UnifiedListing struct {
	ID           int                    `json:"id" db:"id"`
	SourceType   string                 `json:"source_type" db:"source_type"` // "c2c" или "b2c"
	UserID       int                    `json:"user_id" db:"user_id"`
	CategoryID   int                    `json:"category_id" db:"category_id"`
	Title        string                 `json:"title" db:"title"`
	Description  string                 `json:"description" db:"description"`
	Price        float64                `json:"price" db:"price"`
	Condition    string                 `json:"condition" db:"condition"`
	Status       string                 `json:"status" db:"status"`
	Location     string                 `json:"location" db:"location"`
	Latitude     *float64               `json:"latitude,omitempty" db:"latitude"`
	Longitude    *float64               `json:"longitude,omitempty" db:"longitude"`
	City         string                 `json:"city" db:"address_city"`
	Country      string                 `json:"country" db:"address_country"`
	ViewsCount   int                    `json:"views_count" db:"views_count"`
	CreatedAt    time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at" db:"updated_at"`
	ShowOnMap    bool                   `json:"show_on_map" db:"show_on_map"`
	OriginalLang string                 `json:"original_language" db:"original_language"`
	StorefrontID *int                   `json:"storefront_id,omitempty" db:"storefront_id"` // только для B2C
	Metadata     map[string]interface{} `json:"metadata,omitempty" db:"metadata"`

	// JSONB поле из VIEW (массив изображений)
	ImagesJSON json.RawMessage `json:"-" db:"images"`

	// Связанные данные (загружаются отдельно)
	Images       []UnifiedImage         `json:"images,omitempty" db:"-"`
	User         *User                  `json:"user,omitempty" db:"-"`
	Category     *MarketplaceCategory   `json:"category,omitempty" db:"-"`
	Storefront   *Storefront            `json:"storefront,omitempty" db:"-"` // только для B2C
	Translations map[string]interface{} `json:"translations,omitempty" db:"-"`

	// Флаги и дополнительные поля
	IsFavorite      bool     `json:"is_favorite" db:"-"`
	HasDiscount     bool     `json:"has_discount" db:"-"`
	OldPrice        *float64 `json:"old_price,omitempty" db:"-"`
	DiscountPercent *int     `json:"discount_percentage,omitempty" db:"-"`
}

// UnifiedImage унифицированная структура изображения
type UnifiedImage struct {
	ID           int    `json:"id"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnail_url,omitempty"`
	IsMain       bool   `json:"is_main"`
	DisplayOrder int    `json:"display_order"`
}

// UnifiedListingsFilters фильтры для unified listings
type UnifiedListingsFilters struct {
	SourceType   string  `json:"source_type"` // "all", "c2c", "b2c"
	CategoryID   int     `json:"category_id"`
	MinPrice     float64 `json:"min_price"`
	MaxPrice     float64 `json:"max_price"`
	Condition    string  `json:"condition"`
	Query        string  `json:"query"`
	UserID       int     `json:"user_id"`
	StorefrontID int     `json:"storefront_id"`
	Limit        int     `json:"limit"`
	Offset       int     `json:"offset"`
}

// UnifiedListingsResponse ответ для unified listings
type UnifiedListingsResponse struct {
	Data       []UnifiedListing `json:"data"`
	Total      int64            `json:"total"`
	Limit      int              `json:"limit"`
	Offset     int              `json:"offset"`
	SourceType string           `json:"source_type,omitempty"`
}

// ParseImages парсит JSONB images в структуру UnifiedImage
func (u *UnifiedListing) ParseImages() error {
	if len(u.ImagesJSON) == 0 {
		u.Images = []UnifiedImage{}
		return nil
	}

	return json.Unmarshal(u.ImagesJSON, &u.Images)
}
