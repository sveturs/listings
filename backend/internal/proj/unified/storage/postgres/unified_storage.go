// backend/internal/proj/unified/storage/postgres/unified_storage.go
package postgres

import (
	"context"
	"fmt"
	"strings"

	"backend/internal/domain/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type UnifiedStorage struct {
	pool *pgxpool.Pool
	log  *zerolog.Logger
}

func NewUnifiedStorage(pool *pgxpool.Pool, log *zerolog.Logger) *UnifiedStorage {
	return &UnifiedStorage{
		pool: pool,
		log:  log,
	}
}

// GetUnifiedListings получает объединенный список C2C + B2C из VIEW
func (s *UnifiedStorage) GetUnifiedListings(
	ctx context.Context,
	filters models.UnifiedListingsFilters,
) ([]models.UnifiedListing, int64, error) {
	// Базовый запрос через unified_listings VIEW
	query := `
	WITH filtered_listings AS (
		SELECT
			ul.*,
			COUNT(*) OVER() as total_count
		FROM unified_listings ul
		WHERE 1=1
	`

	args := []interface{}{}
	argCount := 0

	// Фильтр по типу источника
	if filters.SourceType != "" && filters.SourceType != "all" {
		argCount++
		query += fmt.Sprintf(" AND ul.source_type = $%d", argCount)
		args = append(args, filters.SourceType)
	}

	// Фильтр по категории
	if filters.CategoryID > 0 {
		argCount++
		query += fmt.Sprintf(" AND ul.category_id = $%d", argCount)
		args = append(args, filters.CategoryID)
	}

	// Фильтр по цене
	if filters.MinPrice > 0 {
		argCount++
		query += fmt.Sprintf(" AND ul.price >= $%d", argCount)
		args = append(args, filters.MinPrice)
	}

	if filters.MaxPrice > 0 {
		argCount++
		query += fmt.Sprintf(" AND ul.price <= $%d", argCount)
		args = append(args, filters.MaxPrice)
	}

	// Фильтр по условию
	if filters.Condition != "" {
		argCount++
		query += fmt.Sprintf(" AND ul.condition = $%d", argCount)
		args = append(args, filters.Condition)
	}

	// Фильтр по витрине (только для B2C)
	if filters.StorefrontID > 0 {
		argCount++
		query += fmt.Sprintf(" AND ul.storefront_id = $%d", argCount)
		args = append(args, filters.StorefrontID)
	}

	// Фильтр по пользователю
	if filters.UserID > 0 {
		argCount++
		query += fmt.Sprintf(" AND ul.user_id = $%d", argCount)
		args = append(args, filters.UserID)
	}

	// Текстовый поиск
	if filters.Query != "" {
		argCount++
		query += fmt.Sprintf(`
			AND (
				LOWER(ul.title) LIKE LOWER($%d)
				OR LOWER(ul.description) LIKE LOWER($%d)
			)
		`, argCount, argCount)
		args = append(args, "%"+filters.Query+"%")
	}

	// Сортировка и пагинация
	query += `
		ORDER BY ul.created_at DESC
	`

	argCount++
	query += fmt.Sprintf(" LIMIT $%d", argCount)
	args = append(args, filters.Limit)

	argCount++
	query += fmt.Sprintf(" OFFSET $%d", argCount)
	args = append(args, filters.Offset)

	query += `
	)
	SELECT * FROM filtered_listings
	`

	// Выполнить запрос
	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to query unified listings")
		return nil, 0, fmt.Errorf("failed to query unified listings: %w", err)
	}
	defer rows.Close()

	listings := []models.UnifiedListing{}
	var totalCount int64

	for rows.Next() {
		var listing models.UnifiedListing
		var tc int64
		var externalID, needsReindex, addressMultilingual interface{} // Игнорируем лишние поля

		err := rows.Scan(
			&listing.ID,
			&listing.SourceType,
			&listing.UserID,
			&listing.CategoryID,
			&listing.Title,
			&listing.Description,
			&listing.Price,
			&listing.Condition,
			&listing.Status,
			&listing.Location,
			&listing.Latitude,
			&listing.Longitude,
			&listing.City,
			&listing.Country,
			&listing.ViewsCount,
			&listing.ShowOnMap,
			&listing.OriginalLang,
			&listing.CreatedAt,
			&listing.UpdatedAt,
			&listing.StorefrontID,
			&externalID, // Игнорируем
			&listing.Metadata,
			&needsReindex,        // Игнорируем
			&addressMultilingual, // Игнорируем
			&listing.ImagesJSON,
			&tc,
		)
		if err != nil {
			s.log.Error().Err(err).Msg("Failed to scan unified listing")
			continue
		}

		// Парсить изображения из JSONB
		if err := listing.ParseImages(); err != nil {
			s.log.Error().Err(err).Int("listing_id", listing.ID).Msg("Failed to parse images")
		}

		listings = append(listings, listing)
		totalCount = tc
	}

	if err := rows.Err(); err != nil {
		s.log.Error().Err(err).Msg("Error iterating unified listings rows")
		return nil, 0, fmt.Errorf("error iterating rows: %w", err)
	}

	s.log.Info().
		Int("count", len(listings)).
		Int64("total", totalCount).
		Str("source_type", filters.SourceType).
		Msg("Fetched unified listings")

	return listings, totalCount, nil
}

// GetUnifiedListingByID получает unified listing по ID и типу
func (s *UnifiedStorage) GetUnifiedListingByID(
	ctx context.Context,
	id int,
	sourceType string,
) (*models.UnifiedListing, error) {
	if sourceType != "c2c" && sourceType != "b2c" {
		return nil, fmt.Errorf("invalid source_type: %s", sourceType)
	}

	query := `
	SELECT *
	FROM unified_listings
	WHERE id = $1 AND source_type = $2
	`

	var listing models.UnifiedListing
	var externalID, needsReindex, addressMultilingual interface{}
	err := s.pool.QueryRow(ctx, query, id, sourceType).Scan(
		&listing.ID,
		&listing.SourceType,
		&listing.UserID,
		&listing.CategoryID,
		&listing.Title,
		&listing.Description,
		&listing.Price,
		&listing.Condition,
		&listing.Status,
		&listing.Location,
		&listing.Latitude,
		&listing.Longitude,
		&listing.City,
		&listing.Country,
		&listing.ViewsCount,
		&listing.ShowOnMap,
		&listing.OriginalLang,
		&listing.CreatedAt,
		&listing.UpdatedAt,
		&listing.StorefrontID,
		&externalID,
		&listing.Metadata,
		&needsReindex,
		&addressMultilingual,
		&listing.ImagesJSON,
	)
	if err != nil {
		s.log.Error().Err(err).Int("id", id).Str("source_type", sourceType).Msg("Failed to get unified listing by ID")
		return nil, fmt.Errorf("failed to get unified listing: %w", err)
	}

	// Парсить изображения
	if err := listing.ParseImages(); err != nil {
		s.log.Error().Err(err).Int("listing_id", listing.ID).Msg("Failed to parse images")
	}

	return &listing, nil
}

// GetUnifiedListingsByIDs получает unified listings по списку IDs
func (s *UnifiedStorage) GetUnifiedListingsByIDs(
	ctx context.Context,
	ids []int,
	sourceType string,
) ([]models.UnifiedListing, error) {
	if len(ids) == 0 {
		return []models.UnifiedListing{}, nil
	}

	// Построить IN clause
	placeholders := []string{}
	args := []interface{}{}
	for i, id := range ids {
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		args = append(args, id)
	}

	query := `
	SELECT *
	FROM unified_listings
	WHERE id IN (` + strings.Join(placeholders, ",") + `)
	`

	if sourceType != "" && sourceType != "all" {
		query += fmt.Sprintf(" AND source_type = $%d", len(args)+1)
		args = append(args, sourceType)
	}

	query += " ORDER BY created_at DESC"

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to get unified listings by IDs")
		return nil, fmt.Errorf("failed to get unified listings: %w", err)
	}
	defer rows.Close()

	listings := []models.UnifiedListing{}
	for rows.Next() {
		var listing models.UnifiedListing
		var externalID, needsReindex, addressMultilingual interface{}
		err := rows.Scan(
			&listing.ID,
			&listing.SourceType,
			&listing.UserID,
			&listing.CategoryID,
			&listing.Title,
			&listing.Description,
			&listing.Price,
			&listing.Condition,
			&listing.Status,
			&listing.Location,
			&listing.Latitude,
			&listing.Longitude,
			&listing.City,
			&listing.Country,
			&listing.ViewsCount,
			&listing.ShowOnMap,
			&listing.OriginalLang,
			&listing.CreatedAt,
			&listing.UpdatedAt,
			&listing.StorefrontID,
			&externalID,
			&listing.Metadata,
			&needsReindex,
			&addressMultilingual,
			&listing.ImagesJSON,
		)
		if err != nil {
			s.log.Error().Err(err).Msg("Failed to scan unified listing")
			continue
		}

		// Парсить изображения
		if err := listing.ParseImages(); err != nil {
			s.log.Error().Err(err).Int("listing_id", listing.ID).Msg("Failed to parse images")
		}

		listings = append(listings, listing)
	}

	return listings, nil
}
