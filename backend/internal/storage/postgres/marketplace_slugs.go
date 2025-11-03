// backend/internal/storage/postgres/marketplace_slugs.go
package postgres

import (
	"context"
	"fmt"

	"backend/internal/logger"

	"github.com/jackc/pgx/v5"
)

// IsSlugUnique checks if slug is unique in c2c_listings table
// excludeID - ID листинга, который нужно исключить из проверки (для update)
// Возвращает true если slug уникален (не занят другими листингами)
func (db *Database) IsSlugUnique(ctx context.Context, slug string, excludeID int) (bool, error) {
	// NOTE: This implementation assumes that 'slug' column exists in c2c_listings table
	// If column doesn't exist yet, this will return an error
	// Migration to add slug column should be created separately

	var query string
	var args []interface{}

	if excludeID > 0 {
		// Проверяем уникальность, исключая указанный ID (для update)
		query = `SELECT EXISTS(SELECT 1 FROM c2c_listings WHERE slug = $1 AND id != $2)`
		args = []interface{}{slug, excludeID}
	} else {
		// Проверяем уникальность для нового листинга
		query = `SELECT EXISTS(SELECT 1 FROM c2c_listings WHERE slug = $1)`
		args = []interface{}{slug}
	}

	var exists bool
	err := db.pool.QueryRow(ctx, query, args...).Scan(&exists)
	if err != nil {
		// Check if error is due to missing column
		if err.Error() == "ERROR: column \"slug\" does not exist (SQLSTATE 42703)" {
			return false, fmt.Errorf("slug column does not exist in c2c_listings table, migration required")
		}
		return false, fmt.Errorf("failed to check slug uniqueness: %w", err)
	}

	logger.Debug().
		Str("slug", slug).
		Int("exclude_id", excludeID).
		Bool("is_unique", !exists).
		Msg("Checked slug uniqueness")

	return !exists, nil
}

// GenerateUniqueSlug generates a unique slug by appending numeric suffix if needed
// baseSlug - базовый slug (например, "prodayu-iphone-15")
// excludeID - ID листинга для update (0 для create)
// Возвращает уникальный slug (добавляя суффикс -1, -2, ... если нужно)
func (db *Database) GenerateUniqueSlug(ctx context.Context, baseSlug string, excludeID int) (string, error) {
	// Check base slug first
	isUnique, err := db.IsSlugUnique(ctx, baseSlug, excludeID)
	if err != nil {
		return "", err
	}

	if isUnique {
		logger.Debug().
			Str("slug", baseSlug).
			Msg("Base slug is unique, using as-is")
		return baseSlug, nil
	}

	// Try with numeric suffix
	for i := 1; i <= 100; i++ {
		candidateSlug := fmt.Sprintf("%s-%d", baseSlug, i)
		isUnique, err := db.IsSlugUnique(ctx, candidateSlug, excludeID)
		if err != nil {
			return "", err
		}

		if isUnique {
			logger.Debug().
				Str("base_slug", baseSlug).
				Str("generated_slug", candidateSlug).
				Int("suffix", i).
				Msg("Generated unique slug with suffix")
			return candidateSlug, nil
		}
	}

	// If all 100 attempts failed, return error
	return "", fmt.Errorf("failed to generate unique slug after 100 attempts for base: %s", baseSlug)
}

// GetListingBySlug retrieves a single listing by slug from database
// NOTE: This is a direct database query, bypassing gRPC microservice
// Use this method when you need to fetch by slug specifically
func (db *Database) GetListingBySlugDirect(ctx context.Context, slug string) (*pgx.Row, error) {
	// NOTE: This implementation assumes that 'slug' column exists in c2c_listings table
	// If column doesn't exist yet, this will return an error

	query := `
		SELECT
			id, user_id, category_id, title, description, price,
			condition, status, location, latitude, longitude,
			address_city, address_country, views_count, show_on_map,
			original_language, created_at, updated_at,
			storefront_id, external_id, metadata,
			address_multilingual
		FROM c2c_listings
		WHERE slug = $1 AND status != 'deleted'
		LIMIT 1
	`

	row := db.pool.QueryRow(ctx, query, slug)
	return &row, nil
}
