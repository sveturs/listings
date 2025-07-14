package types

import (
	"time"

	"github.com/google/uuid"
)

// District represents a city district (e.g., Belgrade districts)
type District struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	CityID      *uuid.UUID `json:"city_id,omitempty" db:"city_id"`
	CountryCode string     `json:"country_code" db:"country_code"`
	Boundary    *Polygon   `json:"boundary,omitempty" db:"-"`
	CenterPoint *Point     `json:"center_point,omitempty" db:"-"`
	Population  *int       `json:"population,omitempty" db:"population"`
	AreaKm2     *float64   `json:"area_km2,omitempty" db:"area_km2"`
	PostalCodes []string   `json:"postal_codes,omitempty" db:"postal_codes"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// Municipality represents a municipality
type Municipality struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	DistrictID  *uuid.UUID `json:"district_id,omitempty" db:"district_id"`
	CountryCode string     `json:"country_code" db:"country_code"`
	Boundary    *Polygon   `json:"boundary,omitempty" db:"-"`
	CenterPoint *Point     `json:"center_point,omitempty" db:"-"`
	Population  *int       `json:"population,omitempty" db:"population"`
	AreaKm2     *float64   `json:"area_km2,omitempty" db:"area_km2"`
	PostalCodes []string   `json:"postal_codes,omitempty" db:"postal_codes"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// Polygon represents a geographic polygon
type Polygon struct {
	Type        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
}

// DistrictSearchParams represents search parameters for districts
type DistrictSearchParams struct {
	CountryCode string      `json:"country_code,omitempty"`
	CityID      *uuid.UUID  `json:"city_id,omitempty"`
	CityIDs     []uuid.UUID `json:"city_ids,omitempty"` // Filter by multiple cities
	Name        string      `json:"name,omitempty"`
	Point       *Point      `json:"point,omitempty"` // Find district containing this point
}

// MunicipalitySearchParams represents search parameters for municipalities
type MunicipalitySearchParams struct {
	CountryCode string     `json:"country_code,omitempty"`
	DistrictID  *uuid.UUID `json:"district_id,omitempty"`
	Name        string     `json:"name,omitempty"`
	Point       *Point     `json:"point,omitempty"` // Find municipality containing this point
}

// DistrictListingSearchParams represents parameters for searching listings by district
type DistrictListingSearchParams struct {
	DistrictID uuid.UUID  `json:"district_id"`
	CategoryID *uuid.UUID `json:"category_id,omitempty"`
	MinPrice   *float64   `json:"min_price,omitempty"`
	MaxPrice   *float64   `json:"max_price,omitempty"`
	Limit      int        `json:"limit"`
	Offset     int        `json:"offset"`
}

// MunicipalityListingSearchParams represents parameters for searching listings by municipality
type MunicipalityListingSearchParams struct {
	MunicipalityID uuid.UUID  `json:"municipality_id"`
	CategoryID     *uuid.UUID `json:"category_id,omitempty"`
	MinPrice       *float64   `json:"min_price,omitempty"`
	MaxPrice       *float64   `json:"max_price,omitempty"`
	Limit          int        `json:"limit"`
	Offset         int        `json:"offset"`
}
