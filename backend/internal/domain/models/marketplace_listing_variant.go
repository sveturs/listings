// backend/internal/domain/models/marketplace_listing_variant.go
package models

// MarketplaceListingVariant представляет вариант товара в маркетплейсе
type MarketplaceListingVariant struct {
	ID         int               `json:"id" db:"id"`
	ListingID  int               `json:"listing_id" db:"listing_id"`
	SKU        string            `json:"sku" db:"sku"`
	Price      *float64          `json:"price" db:"price"`
	Stock      *int              `json:"stock" db:"stock"`
	Attributes map[string]string `json:"attributes" db:"attributes"` // JSON поле для хранения атрибутов варианта
	ImageURL   *string           `json:"image_url,omitempty" db:"image_url"`
	IsActive   bool              `json:"is_active" db:"is_active"`
	CreatedAt  *string           `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt  *string           `json:"updated_at,omitempty" db:"updated_at"`
}
