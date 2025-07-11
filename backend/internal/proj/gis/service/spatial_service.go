package service

import (
	"context"
	"fmt"
	"math"

	"github.com/jmoiron/sqlx"

	"backend/internal/proj/gis/repository"
	"backend/internal/proj/gis/types"
)

// SpatialService сервис для работы с пространственными данными
type SpatialService struct {
	repo *repository.PostGISRepository
}

// NewSpatialService создает новый сервис
func NewSpatialService(db *sqlx.DB) *SpatialService {
	return &SpatialService{
		repo: repository.NewPostGISRepository(db),
	}
}

// SearchListings поиск объявлений с учетом геопозиции
func (s *SpatialService) SearchListings(ctx context.Context, params types.SearchParams) (*types.SearchResponse, error) {
	// Валидация параметров
	if err := s.validateSearchParams(&params); err != nil {
		return nil, err
	}

	// Устанавливаем дефолтные значения
	if params.Limit <= 0 {
		params.Limit = 50
	}
	if params.Limit > 1000 {
		params.Limit = 1000
	}

	// Выполняем поиск
	listings, totalCount, err := s.repo.SearchListings(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to search listings: %w", err)
	}

	// Формируем ответ
	response := &types.SearchResponse{
		Listings:   listings,
		TotalCount: totalCount,
		HasMore:    int64(params.Offset+len(listings)) < totalCount,
	}

	return response, nil
}

// GetListingLocation получение геоданных объявления
func (s *SpatialService) GetListingLocation(ctx context.Context, listingID int) (*types.GeoListing, error) {
	listing, err := s.repo.GetListingByID(ctx, listingID)
	if err != nil {
		return nil, fmt.Errorf("failed to get listing location: %w", err)
	}

	return listing, nil
}

// UpdateListingLocation обновление геолокации объявления
func (s *SpatialService) UpdateListingLocation(ctx context.Context, listingID int, location types.Point, address string) error {
	// Валидация координат
	if location.Lat < -90 || location.Lat > 90 {
		return types.ErrInvalidLatitude
	}
	if location.Lng < -180 || location.Lng > 180 {
		return types.ErrInvalidLongitude
	}

	err := s.repo.UpdateListingLocation(ctx, listingID, location, address)
	if err != nil {
		return fmt.Errorf("failed to update listing location: %w", err)
	}

	return nil
}

// GetNearbyListings получение ближайших объявлений
func (s *SpatialService) GetNearbyListings(ctx context.Context, center types.Point, radiusKm float64, limit int) (*types.SearchResponse, error) {
	if radiusKm <= 0 {
		return nil, types.ErrInvalidRadius
	}

	params := types.SearchParams{
		Center:    &center,
		RadiusKm:  radiusKm,
		Limit:     limit,
		SortBy:    "distance",
		SortOrder: "asc",
	}

	return s.SearchListings(ctx, params)
}

// GetListingsInBounds получение объявлений в заданных границах
func (s *SpatialService) GetListingsInBounds(ctx context.Context, bounds types.Bounds, categories []string, limit int) (*types.SearchResponse, error) {
	params := types.SearchParams{
		Bounds:     &bounds,
		Categories: categories,
		Limit:      limit,
		SortBy:     "created_at",
		SortOrder:  "desc",
	}

	return s.SearchListings(ctx, params)
}

// CalculateDistance вычисление расстояния между двумя точками (в километрах)
func (s *SpatialService) CalculateDistance(from, to types.Point) float64 {
	// Формула гаверсинусов для расчета расстояния на сфере
	const earthRadiusKm = 6371.0

	lat1Rad := degreesToRadians(from.Lat)
	lat2Rad := degreesToRadians(to.Lat)
	deltaLat := degreesToRadians(to.Lat - from.Lat)
	deltaLng := degreesToRadians(to.Lng - from.Lng)

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLng/2)*math.Sin(deltaLng/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}

// GetBoundsForRadius вычисление границ для заданного радиуса от центра
func (s *SpatialService) GetBoundsForRadius(center types.Point, radiusKm float64) types.Bounds {
	// Примерные вычисления (1 градус широты ≈ 111 км)
	latDelta := radiusKm / 111.0

	// Для долготы зависит от широты
	lngDelta := radiusKm / (111.0 * math.Cos(degreesToRadians(center.Lat)))

	return types.Bounds{
		North: math.Min(90, center.Lat+latDelta),
		South: math.Max(-90, center.Lat-latDelta),
		East:  math.Min(180, center.Lng+lngDelta),
		West:  math.Max(-180, center.Lng-lngDelta),
	}
}

// validateSearchParams валидация параметров поиска
func (s *SpatialService) validateSearchParams(params *types.SearchParams) error {
	// Проверка границ
	if params.Bounds != nil {
		if err := params.Bounds.Validate(); err != nil {
			return err
		}
	}

	// Проверка центра и радиуса
	if params.Center != nil {
		if params.Center.Lat < -90 || params.Center.Lat > 90 {
			return types.ErrInvalidLatitude
		}
		if params.Center.Lng < -180 || params.Center.Lng > 180 {
			return types.ErrInvalidLongitude
		}

		if params.RadiusKm < 0 {
			return types.ErrInvalidRadius
		}
	}

	// Проверка сортировки
	validSortFields := map[string]bool{
		"distance":   true,
		"price":      true,
		"created_at": true,
	}

	if params.SortBy != "" && !validSortFields[params.SortBy] {
		params.SortBy = "created_at"
	}

	if params.SortOrder != "asc" && params.SortOrder != "desc" {
		params.SortOrder = "desc"
	}

	return nil
}

// Helper функции

func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}
