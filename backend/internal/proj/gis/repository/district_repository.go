package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"backend/internal/logger"
	"backend/internal/proj/gis/types"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

// ErrDistrictNotFound возвращается когда район не найден
var ErrDistrictNotFound = errors.New("district not found")

// ErrMunicipalityNotFound возвращается когда муниципалитет не найден
var ErrMunicipalityNotFound = errors.New("municipality not found")

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

	if len(params.CityIDs) > 0 {
		placeholders := make([]string, len(params.CityIDs))
		for i, cityID := range params.CityIDs {
			argCount++
			placeholders[i] = fmt.Sprintf("$%d", argCount)
			args = append(args, cityID)
		}
		query += fmt.Sprintf(" AND city_id IN (%s)", strings.Join(placeholders, ","))
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
	}

	query += " ORDER BY name"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query districts with query: %s, args: %v", query, args)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error().Err(err).Msg("Failed to close rows")
		}
	}()

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
		if boundaryJSON.Valid && boundaryJSON.String != "" {
			boundary, err := parsePolygonFromGeoJSON(boundaryJSON.String)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to parse boundary for district %s", d.ID)
			}
			d.Boundary = boundary
		}
		if centerJSON.Valid && centerJSON.String != "" {
			centerPoint, err := parsePointFromGeoJSON(centerJSON.String)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to parse center point for district %s", d.ID)
			}
			d.CenterPoint = centerPoint
		}

		districts = append(districts, d)
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error iterating over districts")
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
			return nil, ErrDistrictNotFound
		}
		return nil, errors.Wrap(err, "failed to get district by ID")
	}

	// Parse geometry JSON
	if boundaryJSON.Valid && boundaryJSON.String != "" {
		boundary, err := parsePolygonFromGeoJSON(boundaryJSON.String)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse boundary for district %s", d.ID)
		}
		d.Boundary = boundary
	}
	if centerJSON.Valid && centerJSON.String != "" {
		centerPoint, err := parsePointFromGeoJSON(centerJSON.String)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse center point for district %s", d.ID)
		}
		d.CenterPoint = centerPoint
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
	}

	query += " ORDER BY name"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to query municipalities")
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error().Err(err).Msg("Failed to close rows")
		}
	}()

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
		if boundaryJSON.Valid && boundaryJSON.String != "" {
			boundary, err := parsePolygonFromGeoJSON(boundaryJSON.String)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to parse boundary for municipality %s", m.ID)
			}
			m.Boundary = boundary
		}
		if centerJSON.Valid && centerJSON.String != "" {
			centerPoint, err := parsePointFromGeoJSON(centerJSON.String)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to parse center point for municipality %s", m.ID)
			}
			m.CenterPoint = centerPoint
		}

		municipalities = append(municipalities, m)
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error iterating over municipalities")
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
			return nil, ErrMunicipalityNotFound
		}
		return nil, errors.Wrap(err, "failed to get municipality by ID")
	}

	// Parse geometry JSON
	if boundaryJSON.Valid && boundaryJSON.String != "" {
		boundary, err := parsePolygonFromGeoJSON(boundaryJSON.String)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse boundary for municipality %s", m.ID)
		}
		m.Boundary = boundary
	}
	if centerJSON.Valid && centerJSON.String != "" {
		centerPoint, err := parsePointFromGeoJSON(centerJSON.String)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse center point for municipality %s", m.ID)
		}
		m.CenterPoint = centerPoint
	}

	return &m, nil
}

// SearchListingsByDistrict searches for listings within a district
func (r *DistrictRepository) SearchListingsByDistrict(ctx context.Context, params types.DistrictListingSearchParams) ([]types.GeoListing, error) {
	query := `
		SELECT
			ml.id,
			ml.title,
			COALESCE(ml.description, '') as description,
			COALESCE(ml.price, 0) as price,
			COALESCE(ml.category_id, 0) as category_id,
			COALESCE(mc.name, '') as category_name,
			COALESCE(ml.user_id, 0) as user_id,
			COALESCE(u.email, '') as user_email,
			COALESCE(u.name, '') as user_name,
			COALESCE(mlg.formatted_address, '') as address,
			COALESCE(ml.address_city, '') as city,
			COALESCE(ml.address_country, '') as country,
			ST_Y(mlg.location::geometry) as latitude,
			ST_X(mlg.location::geometry) as longitude,
			ml.created_at,
			ml.updated_at,
			COALESCE(
				(SELECT mi.public_url
				 FROM marketplace_images mi
				 WHERE mi.listing_id = ml.id
				 ORDER BY mi.created_at
				 LIMIT 1),
				''
			) as first_image_url
		FROM c2c_listings ml
		JOIN listings_geo mlg ON ml.id = mlg.listing_id
		LEFT JOIN c2c_categories mc ON ml.category_id = mc.id
		LEFT JOIN users u ON ml.user_id = u.id
		JOIN districts d ON d.id = $1
		WHERE ST_Contains(d.boundary, mlg.location::geometry)
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
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error().Err(err).Msg("Failed to close rows")
		}
	}()

	var results []types.GeoListing
	for rows.Next() {
		var result types.GeoListing
		var categoryName, userEmail, userName, firstImageURL sql.NullString
		var lat, lng float64
		var categoryID sql.NullString
		var address, city, country string
		err := rows.Scan(
			&result.ID,
			&result.Title,
			&result.Description,
			&result.Price,
			&categoryID,
			&categoryName,
			&result.UserID,
			&userEmail,
			&userName,
			&address,
			&city,
			&country,
			&lat,
			&lng,
			&result.CreatedAt,
			&result.UpdatedAt,
			&firstImageURL,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan listing result")
		}

		result.Address = address
		result.Location = types.Point{Lat: lat, Lng: lng}
		if categoryName.Valid {
			result.Category = categoryName.String
		}

		// Добавляем первое изображение в массив, если есть
		if firstImageURL.Valid && firstImageURL.String != "" {
			result.Images = []string{firstImageURL.String}
		}
		results = append(results, result)
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error iterating over district listings")
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
			ml.category_id,
			mc.name as category_name,
			ml.user_id,
			u.email as user_email,
			u.name as user_name,
			COALESCE(mlg.formatted_address, '') as address,
			COALESCE(ml.address_city, '') as city,
			COALESCE(ml.address_country, '') as country,
			ST_Y(mlg.location::geometry) as latitude,
			ST_X(mlg.location::geometry) as longitude,
			ml.created_at,
			ml.updated_at,
			COALESCE(
				(SELECT mi.public_url
				 FROM marketplace_images mi
				 WHERE mi.listing_id = ml.id
				 ORDER BY mi.created_at
				 LIMIT 1),
				''
			) as first_image_url
		FROM c2c_listings ml
		JOIN listings_geo mlg ON ml.id = mlg.listing_id
		LEFT JOIN c2c_categories mc ON ml.category_id = mc.id
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
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error().Err(err).Msg("Failed to close rows")
		}
	}()

	var results []types.GeoListing
	for rows.Next() {
		var result types.GeoListing
		var categoryName, userEmail, userName, firstImageURL sql.NullString
		var lat, lng float64
		var categoryID sql.NullString
		var address, city, country string
		err := rows.Scan(
			&result.ID,
			&result.Title,
			&result.Description,
			&result.Price,
			&categoryID,
			&categoryName,
			&result.UserID,
			&userEmail,
			&userName,
			&address,
			&city,
			&country,
			&lat,
			&lng,
			&result.CreatedAt,
			&result.UpdatedAt,
			&firstImageURL,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan listing result")
		}

		result.Address = address
		result.Location = types.Point{Lat: lat, Lng: lng}
		if categoryName.Valid {
			result.Category = categoryName.String
		}

		// Добавляем первое изображение в массив, если есть
		if firstImageURL.Valid && firstImageURL.String != "" {
			result.Images = []string{firstImageURL.String}
		}
		results = append(results, result)
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error iterating over municipality listings")
	}

	return results, nil
}

// GetDistrictBoundaryGeoJSON returns district boundary as GeoJSON string
func (r *DistrictRepository) GetDistrictBoundaryGeoJSON(ctx context.Context, districtID uuid.UUID) (string, error) {
	var boundaryGeoJSON sql.NullString

	query := `
		SELECT ST_AsGeoJSON(boundary) as boundary_geojson
		FROM districts
		WHERE id = $1
	`

	err := r.db.QueryRowContext(ctx, query, districtID).Scan(&boundaryGeoJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", types.ErrDistrictNotFound
		}
		return "", errors.Wrapf(err, "failed to get district boundary for ID: %s", districtID)
	}

	if !boundaryGeoJSON.Valid || boundaryGeoJSON.String == "" {
		return "", errors.New("district boundary is null or empty")
	}

	return boundaryGeoJSON.String, nil
}

// GetCities returns cities with optional filtering
func (r *DistrictRepository) GetCities(ctx context.Context, params types.CitySearchParams) ([]types.City, error) {
	query := `
		SELECT
			id, name, slug, country_code,
			ST_AsGeoJSON(center_point) as center_json,
			ST_AsGeoJSON(boundary) as boundary_json,
			population, area_km2, postal_codes,
			has_districts, priority,
			created_at, updated_at
		FROM cities
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 0

	// Filter by country code
	if params.CountryCode != "" {
		argCount++
		query += fmt.Sprintf(" AND country_code = $%d", argCount)
		args = append(args, params.CountryCode)
	}

	// Filter by districts availability
	if params.HasDistricts != nil {
		argCount++
		query += fmt.Sprintf(" AND has_districts = $%d", argCount)
		args = append(args, *params.HasDistricts)
	}

	// Text search by name
	if params.SearchQuery != "" {
		argCount++
		query += fmt.Sprintf(" AND (name ILIKE $%d OR slug ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+params.SearchQuery+"%")
	}

	// Bounds filtering - use boundary if available, otherwise use center_point
	if params.Bounds != nil {
		argCount++
		query += fmt.Sprintf(` AND (
			(boundary IS NOT NULL AND ST_Intersects(boundary, ST_MakeEnvelope($%d, $%d, $%d, $%d, 4326)))
			OR
			(boundary IS NULL AND ST_Intersects(center_point, ST_MakeEnvelope($%d, $%d, $%d, $%d, 4326)))
		)`, argCount, argCount+1, argCount+2, argCount+3, argCount, argCount+1, argCount+2, argCount+3)
		args = append(args, params.Bounds.West, params.Bounds.South, params.Bounds.East, params.Bounds.North)
		argCount += 3
	}

	// Order by priority and name
	query += " ORDER BY priority DESC, name"

	// Apply limit and offset
	if params.Limit > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, params.Limit)
	}

	if params.Offset > 0 {
		argCount++
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, params.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to query cities with query: %s, args: %v", query, args)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Error().Err(err).Msg("Failed to close rows")
		}
	}()

	var cities []types.City
	for rows.Next() {
		var c types.City
		var centerJSON, boundaryJSON sql.NullString

		err := rows.Scan(
			&c.ID, &c.Name, &c.Slug, &c.CountryCode,
			&centerJSON, &boundaryJSON,
			&c.Population, &c.AreaKm2, pq.Array(&c.PostalCodes),
			&c.HasDistricts, &c.Priority,
			&c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan city")
		}

		// Parse center point geometry from JSON
		if centerJSON.Valid && centerJSON.String != "" {
			centerPoint, err := parsePointFromGeoJSON(centerJSON.String)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to parse center point for city %s", c.ID)
			}
			c.CenterPoint = centerPoint
		}

		// Parse boundary geometry from JSON
		if boundaryJSON.Valid && boundaryJSON.String != "" {
			boundary, err := parsePolygonFromGeoJSON(boundaryJSON.String)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to parse boundary for city %s", c.ID)
			}
			c.Boundary = boundary
		}

		cities = append(cities, c)
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error iterating over cities")
	}

	return cities, nil
}

// parsePointFromGeoJSON parses a Point from PostGIS GeoJSON output
func parsePointFromGeoJSON(geoJSON string) (*types.Point, error) {
	var geom struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	}

	if err := json.Unmarshal([]byte(geoJSON), &geom); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal point GeoJSON")
	}

	if geom.Type != "Point" {
		return nil, fmt.Errorf("expected Point geometry, got %s", geom.Type)
	}

	if len(geom.Coordinates) != 2 {
		return nil, fmt.Errorf("point coordinates must have exactly 2 elements, got %d", len(geom.Coordinates))
	}

	return &types.Point{
		Lng: geom.Coordinates[0], // GeoJSON uses [longitude, latitude]
		Lat: geom.Coordinates[1],
	}, nil
}

// parsePolygonFromGeoJSON parses a Polygon from PostGIS GeoJSON output
func parsePolygonFromGeoJSON(geoJSON string) (*types.Polygon, error) {
	var geom struct {
		Type        string        `json:"type"`
		Coordinates [][][]float64 `json:"coordinates"`
	}

	if err := json.Unmarshal([]byte(geoJSON), &geom); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal polygon GeoJSON")
	}

	if geom.Type != "Polygon" {
		return nil, fmt.Errorf("expected Polygon geometry, got %s", geom.Type)
	}

	return &types.Polygon{
		Type:        geom.Type,
		Coordinates: geom.Coordinates,
	}, nil
}
