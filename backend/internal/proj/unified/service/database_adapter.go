// backend/internal/proj/unified/service/database_adapter.go
package service

import (
	"context"
	"fmt"

	"backend/internal/domain/models"
	"backend/internal/domain/search"
)

// DatabaseC2CAdapter адаптер для Database к интерфейсу C2CRepository
type DatabaseC2CAdapter struct {
	db DatabaseInterface
}

// DatabaseInterface определяет методы Database которые нужны для адаптера
type DatabaseInterface interface {
	CreateListing(ctx context.Context, listing *models.MarketplaceListing) (int, error)
	GetListingByID(ctx context.Context, id int) (*models.MarketplaceListing, error)
	UpdateListing(ctx context.Context, listing *models.MarketplaceListing) error
	DeleteListingAdmin(ctx context.Context, id int) error
	GetStorefrontProductByID(ctx context.Context, productID int) (*models.StorefrontProduct, error)
	CreateStorefrontProduct(ctx context.Context, storefrontID int, req *models.CreateProductRequest) (*models.StorefrontProduct, error)
	UpdateStorefrontProduct(ctx context.Context, storefrontID, productID int, req *models.UpdateProductRequest) error
	DeleteStorefrontProduct(ctx context.Context, storefrontID, productID int) error
	GetStorefrontProducts(ctx context.Context, filter models.ProductFilter) ([]*models.StorefrontProduct, error)
	GetStorefrontByID(ctx context.Context, id int) (*models.Storefront, error)
	SearchListings(ctx context.Context, params *search.SearchParams) (*search.SearchResult, error)
	IndexListing(ctx context.Context, listing *models.MarketplaceListing) error
	DeleteListingIndex(ctx context.Context, id string) error
}

// NewDatabaseC2CAdapter создает новый адаптер для C2C
func NewDatabaseC2CAdapter(db DatabaseInterface) *DatabaseC2CAdapter {
	return &DatabaseC2CAdapter{db: db}
}

// CreateListing создает C2C listing
func (a *DatabaseC2CAdapter) CreateListing(ctx context.Context, listing *models.MarketplaceListing) (int, error) {
	return a.db.CreateListing(ctx, listing)
}

// GetListing получает C2C listing по ID
func (a *DatabaseC2CAdapter) GetListing(ctx context.Context, id int) (*models.MarketplaceListing, error) {
	return a.db.GetListingByID(ctx, id)
}

// UpdateListing обновляет C2C listing
func (a *DatabaseC2CAdapter) UpdateListing(ctx context.Context, listing *models.MarketplaceListing) error {
	return a.db.UpdateListing(ctx, listing)
}

// DeleteListing удаляет C2C listing (admin version - без проверки userID)
func (a *DatabaseC2CAdapter) DeleteListing(ctx context.Context, id int) error {
	return a.db.DeleteListingAdmin(ctx, id)
}

// DatabaseB2CAdapter адаптер для Database к интерфейсу B2CRepository
type DatabaseB2CAdapter struct {
	db DatabaseInterface
}

// NewDatabaseB2CAdapter создает новый адаптер для B2C
func NewDatabaseB2CAdapter(db DatabaseInterface) *DatabaseB2CAdapter {
	return &DatabaseB2CAdapter{db: db}
}

// GetStorefrontProductByID получает B2C product по ID
func (a *DatabaseB2CAdapter) GetStorefrontProductByID(ctx context.Context, productID int) (*models.StorefrontProduct, error) {
	return a.db.GetStorefrontProductByID(ctx, productID)
}

// CreateStorefrontProduct создает B2C product
func (a *DatabaseB2CAdapter) CreateStorefrontProduct(ctx context.Context, storefrontID int, req *models.CreateProductRequest) (*models.StorefrontProduct, error) {
	return a.db.CreateStorefrontProduct(ctx, storefrontID, req)
}

// UpdateStorefrontProduct обновляет B2C product
func (a *DatabaseB2CAdapter) UpdateStorefrontProduct(ctx context.Context, storefrontID, productID int, req *models.UpdateProductRequest) error {
	return a.db.UpdateStorefrontProduct(ctx, storefrontID, productID, req)
}

// DeleteStorefrontProduct удаляет B2C product
func (a *DatabaseB2CAdapter) DeleteStorefrontProduct(ctx context.Context, storefrontID, productID int) error {
	return a.db.DeleteStorefrontProduct(ctx, storefrontID, productID)
}

// GetStorefrontProducts получает список B2C products
func (a *DatabaseB2CAdapter) GetStorefrontProducts(ctx context.Context, filter models.ProductFilter) ([]*models.StorefrontProduct, error) {
	return a.db.GetStorefrontProducts(ctx, filter)
}

// GetStorefrontByID получает storefront по ID
func (a *DatabaseB2CAdapter) GetStorefrontByID(ctx context.Context, id int) (*models.Storefront, error) {
	return a.db.GetStorefrontByID(ctx, id)
}

// DatabaseOpenSearchAdapter адаптер для Database к интерфейсу OpenSearchRepository
type DatabaseOpenSearchAdapter struct {
	db DatabaseInterface
}

// NewDatabaseOpenSearchAdapter создает новый адаптер для OpenSearch
func NewDatabaseOpenSearchAdapter(db DatabaseInterface) *DatabaseOpenSearchAdapter {
	return &DatabaseOpenSearchAdapter{db: db}
}

// SearchListings выполняет поиск через OpenSearch
// Конвертирует ServiceParams в SearchParams для Database
func (a *DatabaseOpenSearchAdapter) SearchListings(ctx context.Context, params *search.ServiceParams) (*search.ServiceResult, error) {
	// Конвертируем ServiceParams → SearchParams
	// CategoryID: string → *int
	var categoryID *int
	if params.CategoryID != "" {
		var id int
		fmt.Sscanf(params.CategoryID, "%d", &id)
		categoryID = &id
	}

	// PriceMin/Max: float64 → *float64
	var priceMin, priceMax *float64
	if params.PriceMin > 0 {
		priceMin = &params.PriceMin
	}
	if params.PriceMax > 0 {
		priceMax = &params.PriceMax
	}

	searchParams := &search.SearchParams{
		Query:         params.Query,
		CategoryID:    categoryID,
		PriceMin:      priceMin,
		PriceMax:      priceMax,
		Condition:     params.Condition,
		Page:          params.Page,
		Size:          params.Size,
		Sort:          params.Sort,
		SortDirection: params.SortDirection,
		City:          params.City,
		Country:       params.Country,
		Language:      params.Language,
		Fuzziness:     params.Fuzziness,
	}

	// Геолокация
	if params.Latitude != 0 && params.Longitude != 0 {
		searchParams.Location = &search.GeoLocation{
			Lat: params.Latitude,
			Lon: params.Longitude,
		}
		searchParams.Distance = params.Distance
	}

	result, err := a.db.SearchListings(ctx, searchParams)
	if err != nil {
		return nil, err
	}

	// Конвертируем SearchResult → ServiceResult
	totalPages := result.Total / params.Size
	if result.Total%params.Size > 0 {
		totalPages++
	}

	serviceResult := &search.ServiceResult{
		Items:       result.Listings,
		Total:       result.Total,
		Page:        params.Page,
		Size:        params.Size,
		TotalPages:  totalPages,
		Took:        result.Took,
		Facets:      result.Aggregations,
		Suggestions: result.Suggestions,
	}

	return serviceResult, nil
}

// Index индексирует listing в OpenSearch
func (a *DatabaseOpenSearchAdapter) Index(ctx context.Context, listing *models.MarketplaceListing) error {
	return a.db.IndexListing(ctx, listing)
}

// Delete удаляет listing из OpenSearch индекса
func (a *DatabaseOpenSearchAdapter) Delete(ctx context.Context, listingID int) error {
	return a.db.DeleteListingIndex(ctx, fmt.Sprintf("%d", listingID))
}
