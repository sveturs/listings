package repository

import (
	"context"
	"database/sql"
	"fmt"

	"backend/internal/proj/gis/types"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

// DistrictRepository handles district and municipality data access
type DistrictRepository struct {
	db *sqlx.DB
}

// NewDistrictRepository creates a new district repository
func NewDistrictRepository(db *sqlx.DB) *DistrictRepository {
	return &DistrictRepository{db: db}
}

// GetDistricts returns all districts with optional filtering
func (r *DistrictRepository) GetDistricts(ctx context.Context, params types.DistrictSearchParams) ([]types.District, error) {
	query := `
		SELECT 
			id, name, city_id, country_code,
			ST_AsGeoJSON(boundary) as boundary_json,
			ST_AsGeoJSON(center_point) as center_json,
			population, area_km2, postal_codes,
			created_at, updated_at
		FROM districts
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 0

	if params.CountryCode != "" {
		argCount++
		query += fmt.Sprintf(" AND country_code = $%d", argCount)
		args = append(args, params.CountryCode)
	}

	if params.CityID != nil {
		argCount++
		query += fmt.Sprintf(" AND city_id = $%d", argCount)
		args = append(args, *params.CityID)
	}

	if params.Name != "" {
		argCount++
		query += fmt.Sprintf(" AND name ILIKE $%d", argCount)
		args = append(args, "%"+params.Name+"%")
	}

	if params.Point != nil {
		argCount++
		query += fmt.Sprintf(" AND ST_Contains(boundary, ST_SetSRID(ST_MakePoint($%d, $%d), 4326))", argCount, argCount+1)
		args = append(args, params.Point.Lng, params.Point.Lat)
		argCount++
	}

	query += " ORDER BY name"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query districts with query: %s, args: %v", query, args)
	}
	defer rows.Close()

	var districts []types.District
	for rows.Next() {
		var d types.District
		var boundaryJSON, centerJSON sql.NullString

		err := rows.Scan(
			&d.ID, &d.Name, &d.CityID, &d.CountryCode,
			&boundaryJSON, &centerJSON,
			&d.Population, &d.AreaKm2, pq.Array(&d.PostalCodes),
			&d.CreatedAt, &d.UpdatedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan district")
		}

		// Parse geometry JSON
		if boundaryJSON.Valid {
			// TODO: Parse GeoJSON to Polygon type
		}
		if centerJSON.Valid {
			// TODO: Parse GeoJSON to Point type
		}

		districts = append(districts, d)
	}

	return districts, nil
}

// GetDistrictByID returns a district by ID
func (r *DistrictRepository) GetDistrictByID(ctx context.Context, id uuid.UUID) (*types.District, error) {
	var d types.District
	var boundaryJSON, centerJSON sql.NullString

	query := `
		SELECT 
			id, name, city_id, country_code,
			ST_AsGeoJSON(boundary) as boundary_json,
			ST_AsGeoJSON(center_point) as center_json,
			population, area_km2, postal_codes,
			created_at, updated_at
		FROM districts
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&d.ID, &d.Name, &d.CityID, &d.CountryCode,
		&boundaryJSON, &centerJSON,
		&d.Population, &d.AreaKm2, pq.Array(&d.PostalCodes),
		&d.CreatedAt, &d.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to get district by ID")
	}

	// Parse geometry JSON
	if boundaryJSON.Valid {
		// TODO: Parse GeoJSON to Polygon type
	}
	if centerJSON.Valid {
		// TODO: Parse GeoJSON to Point type
	}

	return &d, nil
}

// GetMunicipalities returns all municipalities with optional filtering
func (r *DistrictRepository) GetMunicipalities(ctx context.Context, params types.MunicipalitySearchParams) ([]types.Municipality, error) {
	query := `
		SELECT 
			id, name, district_id, country_code,
			ST_AsGeoJSON(boundary) as boundary_json,
			ST_AsGeoJSON(center_point) as center_json,
			population, area_km2, postal_codes,
			created_at, updated_at
		FROM municipalities
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 0

	if params.CountryCode != "" {
		argCount++
		query += fmt.Sprintf(" AND country_code = $%d", argCount)
		args = append(args, params.CountryCode)
	}

	if params.DistrictID != nil {
		argCount++
		query += fmt.Sprintf(" AND district_id = $%d", argCount)
		args = append(args, *params.DistrictID)
	}

	if params.Name != "" {
		argCount++
		query += fmt.Sprintf(" AND name ILIKE $%d", argCount)
		args = append(args, "%"+params.Name+"%")
	}

	if params.Point != nil {
		argCount++
		query += fmt.Sprintf(" AND ST_Contains(boundary, ST_SetSRID(ST_MakePoint($%d, $%d), 4326))", argCount, argCount+1)
		args = append(args, params.Point.Lng, params.Point.Lat)
		argCount++
	}

	query += " ORDER BY name"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query municipalities")
	}
	defer rows.Close()

	var municipalities []types.Municipality
	for rows.Next() {
		var m types.Municipality
		var boundaryJSON, centerJSON sql.NullString

		err := rows.Scan(
			&m.ID, &m.Name, &m.DistrictID, &m.CountryCode,
			&boundaryJSON, &centerJSON,
			&m.Population, &m.AreaKm2, pq.Array(&m.PostalCodes),
			&m.CreatedAt, &m.UpdatedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan municipality")
		}

		// Parse geometry JSON
		if boundaryJSON.Valid {
			// TODO: Parse GeoJSON to Polygon type
		}
		if centerJSON.Valid {
			// TODO: Parse GeoJSON to Point type
		}

		municipalities = append(municipalities, m)
	}

	return municipalities, nil
}

// GetMunicipalityByID returns a municipality by ID
func (r *DistrictRepository) GetMunicipalityByID(ctx context.Context, id uuid.UUID) (*types.Municipality, error) {
	var m types.Municipality
	var boundaryJSON, centerJSON sql.NullString

	query := `
		SELECT 
			id, name, district_id, country_code,
			ST_AsGeoJSON(boundary) as boundary_json,
			ST_AsGeoJSON(center_point) as center_json,
			population, area_km2, postal_codes,
			created_at, updated_at
		FROM municipalities
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.ID, &m.Name, &m.DistrictID, &m.CountryCode,
		&boundaryJSON, &centerJSON,
		&m.Population, &m.AreaKm2, pq.Array(&m.PostalCodes),
		&m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to get municipality by ID")
	}

	// Parse geometry JSON
	if boundaryJSON.Valid {
		// TODO: Parse GeoJSON to Polygon type
	}
	if centerJSON.Valid {
		// TODO: Parse GeoJSON to Point type
	}

	return &m, nil
}

// SearchListingsByDistrict searches for listings within a district
func (r *DistrictRepository) SearchListingsByDistrict(ctx context.Context, params types.DistrictListingSearchParams) ([]types.GeoListing, error) {
	query := `
		SELECT 
			ml.id,
			ml.title,
			ml.description,
			ml.price,
			ml.currency,
			ml.category_id,
			mc.name as category_name,
			ml.user_id,
			u.email as user_email,
			u.name as user_name,
			mlg.address,
			mlg.city,
			mlg.country,
			ST_Y(mlg.location) as latitude,
			ST_X(mlg.location) as longitude,
			ml.created_at,
			ml.updated_at,
			COALESCE(
				(SELECT mi.file_url 
				 FROM marketplace_images mi 
				 WHERE mi.listing_id = ml.id 
				 ORDER BY mi.order_index, mi.created_at 
				 LIMIT 1), 
				''
			) as first_image_url
		FROM marketplace_listings ml
		JOIN listings_geo mlg ON ml.id = mlg.listing_id
		LEFT JOIN marketplace_categories mc ON ml.category_id = mc.id
		LEFT JOIN users u ON ml.user_id = u.id
		WHERE mlg.district_id = $1
			AND ml.status = 'active'
	`

	args := []interface{}{params.DistrictID}
	argCount := 1

	if params.CategoryID != nil {
		argCount++
		query += fmt.Sprintf(" AND ml.category_id = $%d", argCount)
		args = append(args, *params.CategoryID)
	}

	if params.MinPrice != nil {
		argCount++
		query += fmt.Sprintf(" AND ml.price >= $%d", argCount)
		args = append(args, *params.MinPrice)
	}

	if params.MaxPrice != nil {
		argCount++
		query += fmt.Sprintf(" AND ml.price <= $%d", argCount)
		args = append(args, *params.MaxPrice)
	}

	query += fmt.Sprintf(" ORDER BY ml.created_at DESC LIMIT $%d OFFSET $%d", argCount+1, argCount+2)
	args = append(args, params.Limit, params.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to search listings by district")
	}
	defer rows.Close()

	var results []types.GeoListing
	for rows.Next() {
		var result types.GeoListing
		var categoryName, userEmail, userName, firstImageURL sql.NullString
		var lat, lng float64
		var categoryID sql.NullString
		err := rows.Scan(
			&result.ID,
			&result.Title,
			&result.Description,
			&result.Price,
			&result.Currency,
			&categoryID,
			&categoryName,
			&result.UserID,
			&userEmail,
			&userName,
			&result.Address,
			&result.Address, // city - временно дублируем address
			&result.Address, // country - временно дублируем address
			&lat,
			&lng,
			&result.CreatedAt,
			&result.UpdatedAt,
			&firstImageURL,
		)

		result.Location = types.Point{Lat: lat, Lng: lng}
		if categoryName.Valid {
			result.Category = categoryName.String
		}

		// Добавляем первое изображение в массив, если есть
		if firstImageURL.Valid && firstImageURL.String != "" {
			result.Images = []string{firstImageURL.String}
		}
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan listing result")
		}
		results = append(results, result)
	}

	return results, nil
}

// SearchListingsByMunicipality searches for listings within a municipality
func (r *DistrictRepository) SearchListingsByMunicipality(ctx context.Context, params types.MunicipalityListingSearchParams) ([]types.GeoListing, error) {
	query := `
		SELECT 
			ml.id,
			ml.title,
			ml.description,
			ml.price,
			ml.currency,
			ml.category_id,
			mc.name as category_name,
			ml.user_id,
			u.email as user_email,
			u.name as user_name,
			mlg.address,
			mlg.city,
			mlg.country,
			ST_Y(mlg.location) as latitude,
			ST_X(mlg.location) as longitude,
			ml.created_at,
			ml.updated_at,
			COALESCE(
				(SELECT mi.file_url 
				 FROM marketplace_images mi 
				 WHERE mi.listing_id = ml.id 
				 ORDER BY mi.order_index, mi.created_at 
				 LIMIT 1), 
				''
			) as first_image_url
		FROM marketplace_listings ml
		JOIN listings_geo mlg ON ml.id = mlg.listing_id
		LEFT JOIN marketplace_categories mc ON ml.category_id = mc.id
		LEFT JOIN users u ON ml.user_id = u.id
		WHERE mlg.municipality_id = $1
			AND ml.status = 'active'
	`

	args := []interface{}{params.MunicipalityID}
	argCount := 1

	if params.CategoryID != nil {
		argCount++
		query += fmt.Sprintf(" AND ml.category_id = $%d", argCount)
		args = append(args, *params.CategoryID)
	}

	if params.MinPrice != nil {
		argCount++
		query += fmt.Sprintf(" AND ml.price >= $%d", argCount)
		args = append(args, *params.MinPrice)
	}

	if params.MaxPrice != nil {
		argCount++
		query += fmt.Sprintf(" AND ml.price <= $%d", argCount)
		args = append(args, *params.MaxPrice)
	}

	query += fmt.Sprintf(" ORDER BY ml.created_at DESC LIMIT $%d OFFSET $%d", argCount+1, argCount+2)
	args = append(args, params.Limit, params.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to search listings by municipality")
	}
	defer rows.Close()

	var results []types.GeoListing
	for rows.Next() {
		var result types.GeoListing
		var categoryName, userEmail, userName, firstImageURL sql.NullString
		var lat, lng float64
		var categoryID sql.NullString
		err := rows.Scan(
			&result.ID,
			&result.Title,
			&result.Description,
			&result.Price,
			&result.Currency,
			&categoryID,
			&categoryName,
			&result.UserID,
			&userEmail,
			&userName,
			&result.Address,
			&result.Address, // city - временно дублируем address
			&result.Address, // country - временно дублируем address
			&lat,
			&lng,
			&result.CreatedAt,
			&result.UpdatedAt,
			&firstImageURL,
		)

		result.Location = types.Point{Lat: lat, Lng: lng}
		if categoryName.Valid {
			result.Category = categoryName.String
		}

		// Добавляем первое изображение в массив, если есть
		if firstImageURL.Valid && firstImageURL.String != "" {
			result.Images = []string{firstImageURL.String}
		}
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan listing result")
		}
		results = append(results, result)
	}

	return results, nil
}
