package service

import (
	"context"
	"encoding/json"
	"errors"
	"sort"

	"backend/internal/proj/gis/constants"
	"backend/internal/proj/gis/repository"
	"backend/internal/proj/gis/types"

	"github.com/google/uuid"
	pkgErrors "github.com/pkg/errors"
)

// ErrDistrictNotFoundByPoint возвращается когда район не найден по координатам
var ErrDistrictNotFoundByPoint = errors.New("district not found by point")

// ErrMunicipalityNotFoundByPoint возвращается когда муниципалитет не найден по координатам
var ErrMunicipalityNotFoundByPoint = errors.New("municipality not found by point")

// DistrictService handles business logic for districts and municipalities
type DistrictService struct {
	repo *repository.DistrictRepository
}

// NewDistrictService creates a new district service
func NewDistrictService(repo *repository.DistrictRepository) *DistrictService {
	return &DistrictService{repo: repo}
}

// GetDistricts returns all districts with optional filtering
func (s *DistrictService) GetDistricts(ctx context.Context, params types.DistrictSearchParams) ([]types.District, error) {
	districts, err := s.repo.GetDistricts(ctx, params)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "failed to get districts")
	}
	return districts, nil
}

// GetDistrictByID returns a district by ID
func (s *DistrictService) GetDistrictByID(ctx context.Context, id uuid.UUID) (*types.District, error) {
	district, err := s.repo.GetDistrictByID(ctx, id)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "failed to get district by ID")
	}
	if district == nil {
		return nil, errors.New("district not found")
	}
	return district, nil
}

// GetMunicipalities returns all municipalities with optional filtering
func (s *DistrictService) GetMunicipalities(ctx context.Context, params types.MunicipalitySearchParams) ([]types.Municipality, error) {
	municipalities, err := s.repo.GetMunicipalities(ctx, params)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "failed to get municipalities")
	}
	return municipalities, nil
}

// GetMunicipalityByID returns a municipality by ID
func (s *DistrictService) GetMunicipalityByID(ctx context.Context, id uuid.UUID) (*types.Municipality, error) {
	municipality, err := s.repo.GetMunicipalityByID(ctx, id)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "failed to get municipality by ID")
	}
	if municipality == nil {
		return nil, errors.New("municipality not found")
	}
	return municipality, nil
}

// SearchListingsByDistrict searches for listings within a district
func (s *DistrictService) SearchListingsByDistrict(ctx context.Context, districtID uuid.UUID, params types.DistrictListingSearchParams) ([]types.GeoListing, error) {
	// Ensure district exists
	district, err := s.GetDistrictByID(ctx, districtID)
	if err != nil {
		return nil, err
	}
	if district == nil {
		return nil, errors.New("district not found")
	}

	params.DistrictID = districtID

	// Set default limit if not provided
	if params.Limit <= 0 {
		params.Limit = constants.DEFAULT_LIMIT
	}
	if params.Limit > constants.MAX_LIMIT {
		params.Limit = constants.MAX_LIMIT
	}

	results, err := s.repo.SearchListingsByDistrict(ctx, params)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "failed to search listings by district")
	}

	return results, nil
}

// SearchListingsByMunicipality searches for listings within a municipality
func (s *DistrictService) SearchListingsByMunicipality(ctx context.Context, municipalityID uuid.UUID, params types.MunicipalityListingSearchParams) ([]types.GeoListing, error) {
	// Ensure municipality exists
	municipality, err := s.GetMunicipalityByID(ctx, municipalityID)
	if err != nil {
		return nil, err
	}
	if municipality == nil {
		return nil, errors.New("municipality not found")
	}

	params.MunicipalityID = municipalityID

	// Set default limit if not provided
	if params.Limit <= 0 {
		params.Limit = constants.DEFAULT_LIMIT
	}
	if params.Limit > constants.MAX_LIMIT {
		params.Limit = constants.MAX_LIMIT
	}

	results, err := s.repo.SearchListingsByMunicipality(ctx, params)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "failed to search listings by municipality")
	}

	return results, nil
}

// GetDistrictByPoint finds the district containing a specific point
func (s *DistrictService) GetDistrictByPoint(ctx context.Context, lat, lng float64) (*types.District, error) {
	params := types.DistrictSearchParams{
		Point: &types.Point{
			Lat: lat,
			Lng: lng,
		},
	}

	districts, err := s.repo.GetDistricts(ctx, params)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "failed to find district by point")
	}

	if len(districts) == 0 {
		return nil, ErrDistrictNotFoundByPoint
	}

	return &districts[0], nil
}

// GetMunicipalityByPoint finds the municipality containing a specific point
func (s *DistrictService) GetMunicipalityByPoint(ctx context.Context, lat, lng float64) (*types.Municipality, error) {
	params := types.MunicipalitySearchParams{
		Point: &types.Point{
			Lat: lat,
			Lng: lng,
		},
	}

	municipalities, err := s.repo.GetMunicipalities(ctx, params)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "failed to find municipality by point")
	}

	if len(municipalities) == 0 {
		return nil, ErrMunicipalityNotFoundByPoint
	}

	return &municipalities[0], nil
}

// GetDistrictBoundary returns district boundary as GeoJSON
func (s *DistrictService) GetDistrictBoundary(ctx context.Context, districtID uuid.UUID) (*types.DistrictBoundaryResponse, error) {
	district, err := s.repo.GetDistrictByID(ctx, districtID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "failed to get district for boundary")
	}
	if district == nil {
		return nil, errors.New("district not found")
	}

	// Get boundary from repository as GeoJSON string
	boundaryGeoJSON, err := s.repo.GetDistrictBoundaryGeoJSON(ctx, districtID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "failed to get district boundary")
	}

	response := &types.DistrictBoundaryResponse{
		ID:       district.ID.String(),
		Name:     district.Name,
		Boundary: json.RawMessage(boundaryGeoJSON),
	}

	if district.CityID != nil {
		cityIDStr := district.CityID.String()
		response.CityID = &cityIDStr
	}

	return response, nil
}

// GetCities returns all cities with optional filtering
func (s *DistrictService) GetCities(ctx context.Context, params types.CitySearchParams) ([]types.City, error) {
	cities, err := s.repo.GetCities(ctx, params)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "failed to get cities")
	}
	return cities, nil
}

// GetVisibleCities returns cities visible in viewport with distance calculation
func (s *DistrictService) GetVisibleCities(ctx context.Context, req types.VisibleCitiesRequest) (*types.VisibleCitiesResponse, error) {
	// Get cities that intersect with the viewport bounds
	params := types.CitySearchParams{
		Bounds:       req.Bounds,
		HasDistricts: nil, // Get all cities, not just those with districts
		Limit:        100,
		Offset:       0,
	}

	allCities, err := s.repo.GetCities(ctx, params)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "failed to get cities in viewport")
	}

	// Calculate distances and sort
	var citiesWithDistance []types.CityWithDistance
	for _, city := range allCities {
		if city.CenterPoint != nil {
			distance := city.CenterPoint.Distance(*req.Center)
			citiesWithDistance = append(citiesWithDistance, types.CityWithDistance{
				City:     city,
				Distance: distance,
			})
		}
	}

	// Sort by distance
	sort.Slice(citiesWithDistance, func(i, j int) bool {
		return citiesWithDistance[i].Distance < citiesWithDistance[j].Distance
	})

	response := &types.VisibleCitiesResponse{
		VisibleCities: citiesWithDistance,
	}

	// Set closest city (first in sorted list)
	if len(citiesWithDistance) > 0 {
		response.ClosestCity = &citiesWithDistance[0]
	}

	return response, nil
}
