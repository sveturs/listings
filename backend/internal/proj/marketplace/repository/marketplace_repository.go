package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"backend/internal/domain/models"
)

// MarketplaceRepository предоставляет методы для работы с объявлениями
type MarketplaceRepository struct {
	db *sqlx.DB
}

// NewMarketplaceRepository создает новый репозиторий маркетплейса
func NewMarketplaceRepository(db *sqlx.DB) *MarketplaceRepository {
	return &MarketplaceRepository{db: db}
}

// GetListingByID получает объявление по ID
func (r *MarketplaceRepository) GetListingByID(ctx context.Context, id int) (*models.MarketplaceListing, error) {
	var listing models.MarketplaceListing
	query := `
		SELECT 
			id, user_id, category_id, status, type, title, description,
			condition, warranty_type, return_policy, location_name, latitude, longitude,
			price, currency_code, price_type, is_negotiable,
			is_featured, promotion_expires_at, views_count, favorites_count,
			average_rating, review_count, last_activity_at,
			created_at, updated_at
		FROM marketplace_listings
		WHERE id = $1`

	err := r.db.GetContext(ctx, &listing, query, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get listing by id")
	}

	return &listing, nil
}
