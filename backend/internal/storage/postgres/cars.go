package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"backend/internal/domain/models"
)

// GetCarMakes возвращает список марок автомобилей с фильтрацией
func (d *Database) GetCarMakes(ctx context.Context, country string, isDomestic bool, isMotorcycle bool, activeOnly bool) ([]models.CarMake, error) {
	var makes []models.CarMake
	var err error

	query := `
		SELECT id, name, slug, logo_url, country, is_active, sort_order, 
		       is_domestic, popularity_rs, created_at, updated_at
		FROM car_makes
		WHERE 1=1
	`

	var conditions []string
	var args []interface{}
	argIndex := 1

	if activeOnly {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, true)
		argIndex++
	}

	if country != "" {
		conditions = append(conditions, fmt.Sprintf("country = $%d", argIndex))
		args = append(args, country)
		argIndex++
	}

	if isDomestic {
		conditions = append(conditions, fmt.Sprintf("is_domestic = $%d", argIndex))
		args = append(args, true)
		// argIndex++ // not needed anymore as we don't have more conditions
	}

	// Пока не поддерживаем фильтр по мотоциклам, так как нет колонки is_motorcycle

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	// Сортировка: сначала по приоритету для сербского рынка, затем по общему порядку
	query += " ORDER BY popularity_rs DESC, sort_order ASC, name ASC"

	err = d.sqlxDB.SelectContext(ctx, &makes, query, args...)
	if err != nil {
		return nil, err
	}

	return makes, nil
}

// GetCarModelsByMake возвращает модели автомобилей для конкретной марки
func (d *Database) GetCarModelsByMake(ctx context.Context, makeSlug string, activeOnly bool) ([]models.CarModel, error) {
	var carModels []models.CarModel

	query := `
		SELECT m.id, m.make_id, m.name, m.slug, m.is_active, m.sort_order, 
		       m.created_at, m.updated_at
		FROM car_models m
		JOIN car_makes mk ON mk.id = m.make_id
		WHERE mk.slug = $1
	`

	args := []interface{}{makeSlug}

	if activeOnly {
		query += " AND m.is_active = $2"
		args = append(args, true)
	}

	query += " ORDER BY m.sort_order ASC, m.name ASC"

	err := d.sqlxDB.SelectContext(ctx, &carModels, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.CarModel{}, nil
		}
		return nil, err
	}

	return carModels, nil
}

// GetCarGenerationsByModel возвращает поколения для конкретной модели
func (d *Database) GetCarGenerationsByModel(ctx context.Context, modelID int, activeOnly bool) ([]models.CarGeneration, error) {
	var generations []models.CarGeneration

	query := `
		SELECT id, model_id, name, slug, year_start, year_end, is_active, 
		       sort_order, created_at, updated_at
		FROM car_generations
		WHERE model_id = $1
	`

	args := []interface{}{modelID}

	if activeOnly {
		query += " AND is_active = $2"
		args = append(args, true)
	}

	query += " ORDER BY year_start DESC, sort_order ASC"

	err := d.sqlxDB.SelectContext(ctx, &generations, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []models.CarGeneration{}, nil
		}
		return nil, err
	}

	return generations, nil
}

// GetCarMakeBySlug возвращает марку автомобиля по slug
func (d *Database) GetCarMakeBySlug(ctx context.Context, slug string) (*models.CarMake, error) {
	var make models.CarMake

	query := `
		SELECT id, name, slug, logo_url, country, is_active, sort_order, 
		       is_domestic, popularity_rs, created_at, updated_at
		FROM car_makes
		WHERE slug = $1
	`

	err := d.sqlxDB.GetContext(ctx, &make, query, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("car make not found: %s", slug)
		}
		return nil, err
	}

	return &make, nil
}

// SearchCarMakes выполняет поиск марок автомобилей по названию
func (d *Database) SearchCarMakes(ctx context.Context, query string, limit int) ([]models.CarMake, error) {
	var makes []models.CarMake

	// Приводим запрос к нижнему регистру для нечувствительного к регистру поиска
	query = strings.ToLower(query)

	sqlQuery := `
		SELECT id, name, slug, logo_url, country, is_active, sort_order, 
		       is_domestic, popularity_rs, created_at, updated_at
		FROM car_makes
		WHERE is_active = true 
		  AND (LOWER(name) LIKE $1 OR LOWER(slug) LIKE $1)
		ORDER BY 
			CASE 
				WHEN LOWER(name) = $2 THEN 1
				WHEN LOWER(name) LIKE $3 THEN 2
				ELSE 3
			END,
			popularity_rs DESC,
			name ASC
		LIMIT $4
	`

	searchPattern := "%" + query + "%"
	exactMatch := query
	startPattern := query + "%"

	err := d.sqlxDB.SelectContext(ctx, &makes, sqlQuery, searchPattern, exactMatch, startPattern, limit)
	if err != nil {
		return nil, err
	}

	return makes, nil
}
