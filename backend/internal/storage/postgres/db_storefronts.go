// backend/internal/storage/postgres/db_storefronts.go
package postgres

import "backend/internal/domain/search"
import (
	"context"
	"fmt"
	"backend/internal/domain/models"
)

// Storefront methods - TODO: temporarily disabled during refactoring

func (db *Database) CreateStorefront(ctx context.Context, userID int, dto *models.StorefrontCreateDTO) (*models.Storefront, error) {
	return nil, fmt.Errorf("storefront service temporarily disabled")
}

func (db *Database) GetStorefrontByID(ctx context.Context, id int) (*models.Storefront, error) {
	return nil, fmt.Errorf("storefront service temporarily disabled")
}

func (db *Database) IncrementListingViewCount(ctx context.Context, id int, userIdentifier string) error {
	return nil // Silent no-op
}

// B2C Product methods

func (db *Database) GetB2CProductImages(ctx context.Context, productID int) ([]models.MarketplaceImage, error) {
	return []models.MarketplaceImage{}, nil
}

// Favorites

func (db *Database) AddToFavorites(ctx context.Context, userID, listingID int) error {
	return fmt.Errorf("favorites service temporarily disabled")
}

func (db *Database) RemoveFromFavorites(ctx context.Context, userID, listingID int) error {
	return fmt.Errorf("favorites service temporarily disabled")
}

func (db *Database) GetUserFavorites(ctx context.Context, userID int) ([]models.MarketplaceListing, error) {
	return []models.MarketplaceListing{}, nil
}

func (db *Database) AddStorefrontToFavorites(ctx context.Context, userID, productID int) error {
	return fmt.Errorf("storefront favorites temporarily disabled")
}

func (db *Database) RemoveStorefrontFromFavorites(ctx context.Context, userID, productID int) error {
	return fmt.Errorf("storefront favorites temporarily disabled")
}

func (db *Database) GetUserStorefrontFavorites(ctx context.Context, userID int) ([]models.MarketplaceListing, error) {
	return []models.MarketplaceListing{}, nil
}


// DeleteStorefront - TODO: disabled during refactoring
func (db *Database) DeleteStorefront(ctx context.Context, id int) error {
	return fmt.Errorf("storefront delete temporarily disabled")
}
func (db *Database) DeleteStorefrontIndex(ctx context.Context, id int) error { return nil }
func (db *Database) GetStorefrontOwnerByProductID(ctx context.Context, productID int) (int, error) { return 0, fmt.Errorf("disabled") }
func (db *Database) GetUserStorefronts(ctx context.Context, userID int) ([]models.Storefront, error) { return nil, fmt.Errorf("disabled") }
func (db *Database) IncrementViewsCount(ctx context.Context, id int) error { return nil }
func (db *Database) IndexStorefront(ctx context.Context, store *models.Storefront) error { return nil }
func (db *Database) GetStorefrontBySlug(ctx context.Context, slug string) (*models.Storefront, error) { return nil, fmt.Errorf("disabled") }
func (db *Database) UpdateStorefront(ctx context.Context, store *models.Storefront) error { return fmt.Errorf("disabled") }
func (db *Database) ReindexAllStorefronts(ctx context.Context) error { return nil }
func (db *Database) SearchStorefronts(ctx context.Context, params *search.SearchParams) (*search.SearchResult, error) { return &search.SearchResult{Listings: []*models.MarketplaceListing{}}, nil }
func (db *Database) SuggestStorefronts(ctx context.Context, prefix string, size int) ([]string, error) { return []string{}, nil }
func (db *Database) PrepareSearchIndex(ctx context.Context) error { return nil }
func (db *Database) Storefront() interface{} { return db.storefrontRepo }
