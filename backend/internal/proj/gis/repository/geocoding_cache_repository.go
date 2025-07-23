package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"backend/internal/logger"
	"backend/internal/proj/gis/types"

	"github.com/jmoiron/sqlx"
)

// ErrGeocodingCacheNotFound возвращается когда запись в кэше не найдена
var ErrGeocodingCacheNotFound = errors.New("geocoding cache entry not found")

// GeocodingCacheRepository репозиторий для кэша геокодирования
type GeocodingCacheRepository struct {
	db *sqlx.DB
}

// NewGeocodingCacheRepository создает новый репозиторий
func NewGeocodingCacheRepository(db *sqlx.DB) *GeocodingCacheRepository {
	return &GeocodingCacheRepository{db: db}
}

// GetByAddress получение из кэша по адресу
func (r *GeocodingCacheRepository) GetByAddress(ctx context.Context, normalizedAddress, language, countryCode string) (*types.GeocodingCacheEntry, error) {
	query := `
		SELECT 
			id, input_address, normalized_address,
			ST_Y(location::geometry) as lat,
			ST_X(location::geometry) as lng,
			address_components, formatted_address, confidence,
			provider, language, country_code, cache_hits,
			created_at, updated_at, expires_at
		FROM geocoding_cache 
		WHERE normalized_address = $1 
			AND language = $2 
			AND (country_code = $3 OR country_code IS NULL OR $3 = '')
			AND expires_at > CURRENT_TIMESTAMP
		ORDER BY cache_hits DESC, confidence DESC
		LIMIT 1`

	var entry types.GeocodingCacheEntry
	var lat, lng float64
	var componentsJSON []byte

	err := r.db.QueryRowContext(ctx, query, normalizedAddress, language, countryCode).Scan(
		&entry.ID,
		&entry.InputAddress,
		&entry.NormalizedAddress,
		&lat,
		&lng,
		&componentsJSON,
		&entry.FormattedAddress,
		&entry.Confidence,
		&entry.Provider,
		&entry.Language,
		&entry.CountryCode,
		&entry.CacheHits,
		&entry.CreatedAt,
		&entry.UpdatedAt,
		&entry.ExpiresAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrGeocodingCacheNotFound // Не найдено в кэше - это нормально
		}
		return nil, fmt.Errorf("failed to get geocoding cache entry: %w", err)
	}

	entry.Location = types.Point{Lat: lat, Lng: lng}

	if err := json.Unmarshal(componentsJSON, &entry.AddressComponents); err != nil {
		return nil, fmt.Errorf("failed to unmarshal address components: %w", err)
	}

	return &entry, nil
}

// Save сохранение в кэш
func (r *GeocodingCacheRepository) Save(ctx context.Context, entry *types.GeocodingCacheEntry) error {
	componentsJSON, err := json.Marshal(entry.AddressComponents)
	if err != nil {
		return fmt.Errorf("failed to marshal address components: %w", err)
	}

	query := `
		INSERT INTO geocoding_cache (
			input_address, normalized_address, location, 
			address_components, formatted_address, confidence,
			provider, language, country_code, expires_at
		) VALUES (
			$1, $2, ST_SetSRID(ST_MakePoint($3, $4), 4326),
			$5, $6, $7, $8, $9, $10, $11
		)
		ON CONFLICT (normalized_address, language, country_code) 
		DO UPDATE SET
			cache_hits = geocoding_cache.cache_hits + 1,
			updated_at = CURRENT_TIMESTAMP,
			confidence = CASE 
				WHEN EXCLUDED.confidence > geocoding_cache.confidence 
				THEN EXCLUDED.confidence 
				ELSE geocoding_cache.confidence 
			END,
			formatted_address = CASE 
				WHEN EXCLUDED.confidence > geocoding_cache.confidence 
				THEN EXCLUDED.formatted_address 
				ELSE geocoding_cache.formatted_address 
			END
		RETURNING id`

	err = r.db.QueryRowContext(ctx, query,
		entry.InputAddress,
		entry.NormalizedAddress,
		entry.Location.Lng,
		entry.Location.Lat,
		componentsJSON,
		entry.FormattedAddress,
		entry.Confidence,
		entry.Provider,
		entry.Language,
		entry.CountryCode,
		entry.ExpiresAt,
	).Scan(&entry.ID)
	if err != nil {
		return fmt.Errorf("failed to save geocoding cache entry: %w", err)
	}

	return nil
}

// IncrementHits увеличение счетчика попаданий
func (r *GeocodingCacheRepository) IncrementHits(ctx context.Context, id int64) error {
	query := `UPDATE geocoding_cache SET cache_hits = cache_hits + 1, updated_at = CURRENT_TIMESTAMP WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to increment cache hits: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no cache entry found with id %d", id)
	}

	return nil
}

// CleanupExpired очистка устаревших записей
func (r *GeocodingCacheRepository) CleanupExpired(ctx context.Context) (int64, error) {
	query := `DELETE FROM geocoding_cache WHERE expires_at < CURRENT_TIMESTAMP`

	result, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup expired cache entries: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

// GetStats получение статистики кэша
func (r *GeocodingCacheRepository) GetStats(ctx context.Context) (map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(*) as total_entries,
			COUNT(*) FILTER (WHERE expires_at > CURRENT_TIMESTAMP) as active_entries,
			COUNT(*) FILTER (WHERE expires_at <= CURRENT_TIMESTAMP) as expired_entries,
			COALESCE(SUM(cache_hits), 0) as total_cache_hits,
			ROUND(AVG(confidence), 3) as avg_confidence,
			COUNT(DISTINCT provider) as provider_count
		FROM geocoding_cache`

	var stats struct {
		TotalEntries   int64   `db:"total_entries"`
		ActiveEntries  int64   `db:"active_entries"`
		ExpiredEntries int64   `db:"expired_entries"`
		TotalCacheHits int64   `db:"total_cache_hits"`
		AvgConfidence  float64 `db:"avg_confidence"`
		ProviderCount  int64   `db:"provider_count"`
	}

	err := r.db.GetContext(ctx, &stats, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get cache stats: %w", err)
	}

	// Получаем топ провайдеров
	providersQuery := `
		SELECT provider, COUNT(*) as count 
		FROM geocoding_cache 
		WHERE expires_at > CURRENT_TIMESTAMP 
		GROUP BY provider 
		ORDER BY count DESC 
		LIMIT 5`

	type providerStat struct {
		Provider string `db:"provider"`
		Count    int64  `db:"count"`
	}

	var providers []providerStat
	err = r.db.SelectContext(ctx, &providers, providersQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider stats: %w", err)
	}

	result := map[string]interface{}{
		"total_entries":    stats.TotalEntries,
		"active_entries":   stats.ActiveEntries,
		"expired_entries":  stats.ExpiredEntries,
		"total_cache_hits": stats.TotalCacheHits,
		"avg_confidence":   stats.AvgConfidence,
		"provider_count":   stats.ProviderCount,
		"top_providers":    providers,
		"cache_hit_rate":   0.0,
	}

	// Вычисляем hit rate если есть данные
	if stats.TotalEntries > 0 {
		hitRate := float64(stats.TotalCacheHits) / float64(stats.TotalEntries)
		result["cache_hit_rate"] = hitRate
	}

	return result, nil
}

// SearchSimilar поиск похожих адресов в кэше для предложений
func (r *GeocodingCacheRepository) SearchSimilar(ctx context.Context, query string, language string, limit int) ([]types.GeocodingCacheEntry, error) {
	searchQuery := `
		SELECT 
			id, input_address, normalized_address,
			ST_Y(location::geometry) as lat,
			ST_X(location::geometry) as lng,
			address_components, formatted_address, confidence,
			provider, language, country_code, cache_hits,
			created_at, updated_at, expires_at
		FROM geocoding_cache 
		WHERE expires_at > CURRENT_TIMESTAMP
			AND language = $1
			AND (
				normalized_address ILIKE $2 
				OR formatted_address ILIKE $2
				OR input_address ILIKE $2
			)
		ORDER BY 
			cache_hits DESC, 
			confidence DESC,
			CASE 
				WHEN normalized_address ILIKE $3 THEN 1
				WHEN formatted_address ILIKE $3 THEN 2
				ELSE 3
			END
		LIMIT $4`

	likePattern := "%" + query + "%"
	exactPattern := query + "%"

	rows, err := r.db.QueryContext(ctx, searchQuery, language, likePattern, exactPattern, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search similar addresses: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error().Err(err).Msg("Failed to close rows")
		}
	}()

	var entries []types.GeocodingCacheEntry

	for rows.Next() {
		var entry types.GeocodingCacheEntry
		var lat, lng float64
		var componentsJSON []byte

		err := rows.Scan(
			&entry.ID,
			&entry.InputAddress,
			&entry.NormalizedAddress,
			&lat,
			&lng,
			&componentsJSON,
			&entry.FormattedAddress,
			&entry.Confidence,
			&entry.Provider,
			&entry.Language,
			&entry.CountryCode,
			&entry.CacheHits,
			&entry.CreatedAt,
			&entry.UpdatedAt,
			&entry.ExpiresAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cache entry: %w", err)
		}

		entry.Location = types.Point{Lat: lat, Lng: lng}

		if err := json.Unmarshal(componentsJSON, &entry.AddressComponents); err != nil {
			continue // Пропускаем запись с некорректными данными
		}

		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating cache entries: %w", err)
	}

	return entries, nil
}

// UpdateExpiry обновление времени истечения для активных записей
func (r *GeocodingCacheRepository) UpdateExpiry(ctx context.Context, id int64, newExpiry time.Time) error {
	query := `UPDATE geocoding_cache SET expires_at = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, newExpiry, id)
	if err != nil {
		return fmt.Errorf("failed to update cache expiry: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no cache entry found with id %d", id)
	}

	return nil
}
