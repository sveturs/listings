package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/sveturs/listings/internal/domain"
)

// GetStorefront retrieves a single storefront by ID
func (r *Repository) GetStorefront(ctx context.Context, storefrontID int64) (*domain.Storefront, error) {
	query := `
		SELECT
			id, user_id, slug, name, description,
			logo_url, banner_url,
			phone, email, website,
			address, city, postal_code, country, latitude, longitude,
			is_active, is_verified,
			rating, reviews_count, products_count, sales_count, views_count, followers_count,
			created_at, updated_at
		FROM storefronts
		WHERE id = $1 AND deleted_at IS NULL
	`

	var sf domain.Storefront
	err := r.db.QueryRowContext(ctx, query, storefrontID).Scan(
		&sf.ID, &sf.UserID, &sf.Slug, &sf.Name, &sf.Description,
		&sf.LogoURL, &sf.BannerURL,
		&sf.Phone, &sf.Email, &sf.Website,
		&sf.Address, &sf.City, &sf.PostalCode, &sf.Country, &sf.Latitude, &sf.Longitude,
		&sf.IsActive, &sf.IsVerified,
		&sf.Rating, &sf.ReviewsCount, &sf.ProductsCount, &sf.SalesCount, &sf.ViewsCount, &sf.FollowersCount,
		&sf.CreatedAt, &sf.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Warn().Int64("storefront_id", storefrontID).Msg("storefront not found")
			return nil, fmt.Errorf("storefront not found: %w", err)
		}
		r.logger.Error().Err(err).Int64("storefront_id", storefrontID).Msg("failed to get storefront")
		return nil, fmt.Errorf("failed to get storefront: %w", err)
	}

	r.logger.Debug().Int64("storefront_id", storefrontID).Str("slug", sf.Slug).Msg("storefront retrieved successfully")
	return &sf, nil
}

// GetStorefrontBySlug retrieves a storefront by slug
func (r *Repository) GetStorefrontBySlug(ctx context.Context, slug string) (*domain.Storefront, error) {
	query := `
		SELECT
			id, user_id, slug, name, description,
			logo_url, banner_url,
			phone, email, website,
			address, city, postal_code, country, latitude, longitude,
			is_active, is_verified,
			rating, reviews_count, products_count, sales_count, views_count, followers_count,
			created_at, updated_at
		FROM storefronts
		WHERE slug = $1 AND deleted_at IS NULL
	`

	var sf domain.Storefront
	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&sf.ID, &sf.UserID, &sf.Slug, &sf.Name, &sf.Description,
		&sf.LogoURL, &sf.BannerURL,
		&sf.Phone, &sf.Email, &sf.Website,
		&sf.Address, &sf.City, &sf.PostalCode, &sf.Country, &sf.Latitude, &sf.Longitude,
		&sf.IsActive, &sf.IsVerified,
		&sf.Rating, &sf.ReviewsCount, &sf.ProductsCount, &sf.SalesCount, &sf.ViewsCount, &sf.FollowersCount,
		&sf.CreatedAt, &sf.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Warn().Str("slug", slug).Msg("storefront not found by slug")
			return nil, fmt.Errorf("storefront not found: %w", err)
		}
		r.logger.Error().Err(err).Str("slug", slug).Msg("failed to get storefront by slug")
		return nil, fmt.Errorf("failed to get storefront by slug: %w", err)
	}

	r.logger.Debug().Int64("storefront_id", sf.ID).Str("slug", slug).Msg("storefront retrieved by slug successfully")
	return &sf, nil
}

// ListStorefronts returns a paginated list of storefronts
func (r *Repository) ListStorefronts(ctx context.Context, limit, offset int) ([]*domain.Storefront, int64, error) {
	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM storefronts WHERE deleted_at IS NULL`
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		r.logger.Error().Err(err).Msg("failed to get storefronts count")
		return nil, 0, fmt.Errorf("failed to get storefronts count: %w", err)
	}

	// Get paginated results
	query := `
		SELECT
			id, user_id, slug, name, description,
			logo_url, banner_url,
			phone, email, website,
			address, city, postal_code, country, latitude, longitude,
			is_active, is_verified,
			rating, reviews_count, products_count, sales_count, views_count, followers_count,
			created_at, updated_at
		FROM storefronts
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		r.logger.Error().Err(err).Int("limit", limit).Int("offset", offset).Msg("failed to list storefronts")
		return nil, 0, fmt.Errorf("failed to list storefronts: %w", err)
	}
	defer rows.Close()

	var storefronts []*domain.Storefront
	for rows.Next() {
		var sf domain.Storefront
		if err := rows.Scan(
			&sf.ID, &sf.UserID, &sf.Slug, &sf.Name, &sf.Description,
			&sf.LogoURL, &sf.BannerURL,
			&sf.Phone, &sf.Email, &sf.Website,
			&sf.Address, &sf.City, &sf.PostalCode, &sf.Country, &sf.Latitude, &sf.Longitude,
			&sf.IsActive, &sf.IsVerified,
			&sf.Rating, &sf.ReviewsCount, &sf.ProductsCount, &sf.SalesCount, &sf.ViewsCount, &sf.FollowersCount,
			&sf.CreatedAt, &sf.UpdatedAt,
		); err != nil {
			r.logger.Error().Err(err).Msg("failed to scan storefront")
			return nil, 0, fmt.Errorf("failed to scan storefront: %w", err)
		}
		storefronts = append(storefronts, &sf)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error().Err(err).Msg("rows iteration error")
		return nil, 0, fmt.Errorf("rows iteration error: %w", err)
	}

	r.logger.Debug().
		Int("count", len(storefronts)).
		Int64("total", total).
		Int("limit", limit).
		Int("offset", offset).
		Msg("storefronts listed successfully")

	return storefronts, total, nil
}
