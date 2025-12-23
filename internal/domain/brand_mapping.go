package domain

import (
	"time"

	"github.com/google/uuid"
)

// BrandCategoryMapping - маппинг бренда на категорию
type BrandCategoryMapping struct {
	ID           uuid.UUID `json:"id" db:"id"`
	BrandName    string    `json:"brand_name" db:"brand_name"`
	BrandAliases []string  `json:"brand_aliases" db:"brand_aliases"`
	CategorySlug string    `json:"category_slug" db:"category_slug"`
	Confidence   float64   `json:"confidence" db:"confidence"`
	IsVerified   bool      `json:"is_verified" db:"is_verified"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// BrandMatch - результат детекции по бренду
type BrandMatch struct {
	Brand         string
	MatchedAlias  string
	CategorySlug  string
	ConfidenceScore float64
}
