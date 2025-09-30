package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"backend/internal/logger"
	"backend/internal/proj/vin/models"

	"github.com/jmoiron/sqlx"
)

// VINStorage предоставляет методы для работы с БД для VIN данных
type VINStorage struct {
	db *sqlx.DB
}

// NewVINStorage создает новый экземпляр storage
func NewVINStorage(db *sqlx.DB) *VINStorage {
	return &VINStorage{db: db}
}

// GetCachedVIN получает закэшированные данные VIN
func (s *VINStorage) GetCachedVIN(ctx context.Context, vin string) (*models.VINDecodeCache, error) {
	query := `
		SELECT
			id, vin, make, model, year,
			engine_type, engine_displacement, transmission_type,
			drivetrain, body_type, fuel_type,
			doors, seats, color_exterior, color_interior,
			manufacturer, country_of_origin, assembly_plant,
			vehicle_class, vehicle_type, gross_vehicle_weight,
			decode_status, error_message, raw_response,
			created_at, updated_at
		FROM vin_decode_cache
		WHERE vin = $1
		LIMIT 1
	`

	var cached models.VINDecodeCache
	err := s.db.GetContext(ctx, &cached, query, vin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("VIN not found in cache: %w", sql.ErrNoRows)
		}
		return nil, fmt.Errorf("failed to get cached VIN: %w", err)
	}

	return &cached, nil
}

// SaveToCache сохраняет декодированные данные VIN в кэш
func (s *VINStorage) SaveToCache(ctx context.Context, data *models.VINDecodeCache) (int64, error) {
	query := `
		INSERT INTO vin_decode_cache (
			vin, make, model, year,
			engine_type, engine_displacement, transmission_type,
			drivetrain, body_type, fuel_type,
			doors, seats, color_exterior, color_interior,
			manufacturer, country_of_origin, assembly_plant,
			vehicle_class, vehicle_type, gross_vehicle_weight,
			decode_status, error_message, raw_response
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
			$21, $22, $23
		)
		ON CONFLICT (vin) DO UPDATE SET
			make = EXCLUDED.make,
			model = EXCLUDED.model,
			year = EXCLUDED.year,
			engine_type = EXCLUDED.engine_type,
			engine_displacement = EXCLUDED.engine_displacement,
			transmission_type = EXCLUDED.transmission_type,
			drivetrain = EXCLUDED.drivetrain,
			body_type = EXCLUDED.body_type,
			fuel_type = EXCLUDED.fuel_type,
			doors = EXCLUDED.doors,
			seats = EXCLUDED.seats,
			color_exterior = EXCLUDED.color_exterior,
			color_interior = EXCLUDED.color_interior,
			manufacturer = EXCLUDED.manufacturer,
			country_of_origin = EXCLUDED.country_of_origin,
			assembly_plant = EXCLUDED.assembly_plant,
			vehicle_class = EXCLUDED.vehicle_class,
			vehicle_type = EXCLUDED.vehicle_type,
			gross_vehicle_weight = EXCLUDED.gross_vehicle_weight,
			decode_status = EXCLUDED.decode_status,
			error_message = EXCLUDED.error_message,
			raw_response = EXCLUDED.raw_response,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id
	`

	var id int64
	err := s.db.GetContext(ctx, &id, query,
		data.VIN, data.Make, data.Model, data.Year,
		data.EngineType, data.EngineDisplacement, data.TransmissionType,
		data.Drivetrain, data.BodyType, data.FuelType,
		data.Doors, data.Seats, data.ColorExterior, data.ColorInterior,
		data.Manufacturer, data.CountryOfOrigin, data.AssemblyPlant,
		data.VehicleClass, data.VehicleType, data.GrossVehicleWeight,
		data.DecodeStatus, data.ErrorMessage, data.RawResponse,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to save VIN to cache: %w", err)
	}

	return id, nil
}

// SaveCheckHistory сохраняет историю проверки VIN
func (s *VINStorage) SaveCheckHistory(ctx context.Context, history *models.VINCheckHistory) error {
	query := `
		INSERT INTO vin_check_history (
			user_id, vin, listing_id,
			decode_success, decode_cache_id,
			check_type, ip_address, user_agent
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
	`

	_, err := s.db.ExecContext(ctx, query,
		history.UserID, history.VIN, history.ListingID,
		history.DecodeSuccess, history.DecodeCacheID,
		history.CheckType, history.IPAddress, history.UserAgent,
	)
	if err != nil {
		return fmt.Errorf("failed to save check history: %w", err)
	}

	return nil
}

// GetHistory получает историю проверок VIN
func (s *VINStorage) GetHistory(ctx context.Context, req *models.VINHistoryRequest) ([]*models.VINCheckHistory, error) {
	query := `
		SELECT
			h.id, h.user_id, h.vin, h.listing_id,
			h.decode_success, h.decode_cache_id,
			h.check_type, h.ip_address, h.user_agent, h.checked_at,
			c.make, c.model, c.year, c.manufacturer, c.body_type, c.fuel_type
		FROM vin_check_history h
		LEFT JOIN vin_decode_cache c ON h.decode_cache_id = c.id
		WHERE 1=1
	`

	args := []interface{}{}
	argCount := 0

	if req.UserID != nil {
		argCount++
		query += fmt.Sprintf(" AND h.user_id = $%d", argCount)
		args = append(args, *req.UserID)
	}

	if req.VIN != "" {
		argCount++
		query += fmt.Sprintf(" AND h.vin = $%d", argCount)
		args = append(args, req.VIN)
	}

	query += " ORDER BY h.checked_at DESC"

	if req.Limit > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, req.Limit)
	} else {
		query += " LIMIT 100" // Default limit
	}

	if req.Offset > 0 {
		argCount++
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, req.Offset)
	}

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get history: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error().Err(err).Msg("Failed to close rows")
		}
	}()

	var history []*models.VINCheckHistory
	for rows.Next() {
		var h models.VINCheckHistory
		var cache models.VINDecodeCache

		err := rows.Scan(
			&h.ID, &h.UserID, &h.VIN, &h.ListingID,
			&h.DecodeSuccess, &h.DecodeCacheID,
			&h.CheckType, &h.IPAddress, &h.UserAgent, &h.CheckedAt,
			&cache.Make, &cache.Model, &cache.Year,
			&cache.Manufacturer, &cache.BodyType, &cache.FuelType,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan history row: %w", err)
		}

		// Если есть данные из кэша, добавляем их
		if h.DecodeCacheID != nil {
			cache.ID = *h.DecodeCacheID
			cache.VIN = h.VIN
			h.DecodeCache = &cache
		}

		history = append(history, &h)
	}

	return history, nil
}

// GetTotalChecks получает общее количество проверок
func (s *VINStorage) GetTotalChecks(ctx context.Context, userID *int64) (int64, error) {
	query := `SELECT COUNT(*) FROM vin_check_history WHERE 1=1`
	args := []interface{}{}

	if userID != nil {
		query += " AND user_id = $1"
		args = append(args, *userID)
	}

	var count int64
	err := s.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to get total checks: %w", err)
	}

	return count, nil
}

// GetUniqueVINCount получает количество уникальных проверенных VIN
func (s *VINStorage) GetUniqueVINCount(ctx context.Context, userID *int64) (int64, error) {
	query := `SELECT COUNT(DISTINCT vin) FROM vin_check_history WHERE 1=1`
	args := []interface{}{}

	if userID != nil {
		query += " AND user_id = $1"
		args = append(args, *userID)
	}

	var count int64
	err := s.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to get unique VIN count: %w", err)
	}

	return count, nil
}

// GetManufacturerStats получает статистику по производителям
func (s *VINStorage) GetManufacturerStats(ctx context.Context, userID *int64) ([]map[string]interface{}, error) {
	query := `
		SELECT
			c.manufacturer,
			COUNT(*) as count
		FROM vin_check_history h
		JOIN vin_decode_cache c ON h.decode_cache_id = c.id
		WHERE c.manufacturer IS NOT NULL
	`

	args := []interface{}{}
	if userID != nil {
		query += " AND h.user_id = $1"
		args = append(args, *userID)
	}

	query += `
		GROUP BY c.manufacturer
		ORDER BY count DESC
		LIMIT 10
	`

	rows, err := s.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get manufacturer stats: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error().Err(err).Msg("Failed to close rows")
		}
	}()

	var stats []map[string]interface{}
	for rows.Next() {
		var manufacturer string
		var count int64

		if err := rows.Scan(&manufacturer, &count); err != nil {
			return nil, fmt.Errorf("failed to scan manufacturer stats: %w", err)
		}

		stats = append(stats, map[string]interface{}{
			"manufacturer": manufacturer,
			"count":        count,
		})
	}

	return stats, nil
}

// GetVINByListingID получает VIN по ID объявления
func (s *VINStorage) GetVINByListingID(ctx context.Context, listingID int64) (string, error) {
	// Сначала проверяем в истории проверок
	query := `
		SELECT DISTINCT vin
		FROM vin_check_history
		WHERE listing_id = $1
		ORDER BY checked_at DESC
		LIMIT 1
	`

	var vin string
	err := s.db.GetContext(ctx, &vin, query, listingID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("failed to get VIN by listing ID: %w", err)
	}

	return vin, nil
}

// UpdateVINCache обновляет кэшированные данные VIN
func (s *VINStorage) UpdateVINCache(ctx context.Context, vin string, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	// Строим динамический UPDATE запрос
	setClause := "updated_at = CURRENT_TIMESTAMP"
	args := []interface{}{}
	argCount := 0

	for field, value := range updates {
		argCount++
		setClause += fmt.Sprintf(", %s = $%d", field, argCount)
		args = append(args, value)
	}

	argCount++
	args = append(args, vin)

	query := fmt.Sprintf(`
		UPDATE vin_decode_cache
		SET %s
		WHERE vin = $%d
	`, setClause, argCount)

	_, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update VIN cache: %w", err)
	}

	return nil
}

// CleanupOldHistory удаляет старые записи истории
func (s *VINStorage) CleanupOldHistory(ctx context.Context, daysToKeep int) (int64, error) {
	query := `
		DELETE FROM vin_check_history
		WHERE checked_at < CURRENT_TIMESTAMP - INTERVAL '%d days'
	`

	result, err := s.db.ExecContext(ctx, fmt.Sprintf(query, daysToKeep))
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup old history: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}
