package service

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/jmoiron/sqlx"

	"backend/internal/proj/gis/repository"
	"backend/internal/proj/gis/types"
)

// SpatialService сервис для работы с пространственными данными
type SpatialService struct {
	repo *repository.PostGISRepository
	db   *sqlx.DB
}

// NewSpatialService создает новый сервис
func NewSpatialService(db *sqlx.DB) *SpatialService {
	return &SpatialService{
		repo: repository.NewPostGISRepository(db),
		db:   db,
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

// SearchByRadius специализированный метод радиусного поиска
func (s *SpatialService) SearchByRadius(ctx context.Context, req types.RadiusSearchRequest) (*types.RadiusSearchResponse, error) {
	// Валидация запроса
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Преобразуем в стандартные параметры поиска
	searchParams := req.ToSearchParams()

	// Выполняем поиск через существующий метод
	searchResponse, err := s.SearchListings(ctx, searchParams)
	if err != nil {
		return nil, fmt.Errorf("failed to perform radius search: %w", err)
	}

	// Создаем специализированный ответ для радиусного поиска
	radiusResponse := &types.RadiusSearchResponse{
		Listings:     searchResponse.Listings,
		TotalCount:   searchResponse.TotalCount,
		HasMore:      searchResponse.HasMore,
		SearchRadius: req.Radius,
		SearchCenter: types.Point{Lat: req.Latitude, Lng: req.Longitude},
	}

	// Дополнительно обогащаем результаты расстояниями если еще не рассчитаны
	s.enrichListingsWithDistances(radiusResponse.Listings, radiusResponse.SearchCenter)

	return radiusResponse, nil
}

// enrichListingsWithDistances обогащает объявления точными расстояниями
func (s *SpatialService) enrichListingsWithDistances(listings []types.GeoListing, center types.Point) {
	for i := range listings {
		if listings[i].Distance == nil {
			distance := s.CalculateDistance(center, listings[i].Location)
			listings[i].Distance = &distance
		}
	}
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

// ========== PHASE 2: Новые методы для умного ввода адресов ==========

// UpdateListingAddress обновление адреса объявления
func (s *SpatialService) UpdateListingAddress(
	ctx context.Context,
	listingID, userID int64,
	req types.UpdateAddressRequest,
	ipAddress, userAgent string,
) (*types.EnhancedListingGeo, error) {
	// Проверяем права доступа
	if err := s.checkListingAccess(ctx, listingID, userID); err != nil {
		return nil, err
	}

	// Получаем текущие данные для логирования
	currentGeo, err := s.getEnhancedListingGeo(ctx, listingID)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get current geo data: %w", err)
	}

	// Вычисляем размытую локацию в зависимости от настроек приватности
	blurredLocation := s.calculateBlurredLocation(req.Location, req.LocationPrivacy)

	// Создаем обновленную запись
	updatedGeo := &types.EnhancedListingGeo{
		ListingID:           listingID,
		Location:            req.Location,
		BlurredLocation:     blurredLocation,
		LocationPrivacy:     req.LocationPrivacy,
		AddressComponents:   req.AddressComponents,
		FormattedAddress:    req.Address,
		AddressVerified:     req.Verified,
		InputMethod:         req.InputMethod,
		GeocodingConfidence: s.calculateConfidenceScore(req),
		UpdatedAt:           time.Now(),
	}

	// Обновляем в базе данных
	if err := s.updateEnhancedListingGeo(ctx, updatedGeo); err != nil {
		return nil, fmt.Errorf("failed to update listing geo: %w", err)
	}

	// Логируем изменение если есть предыдущие данные
	if currentGeo != nil {
		changeLog := &types.AddressChangeLog{
			ListingID:        listingID,
			UserID:           userID,
			OldAddress:       currentGeo.FormattedAddress,
			NewAddress:       req.Address,
			OldLocation:      &currentGeo.Location,
			NewLocation:      &req.Location,
			ChangeReason:     string(req.InputMethod),
			ConfidenceBefore: currentGeo.GeocodingConfidence,
			ConfidenceAfter:  updatedGeo.GeocodingConfidence,
			IPAddress:        ipAddress,
			UserAgent:        userAgent,
			CreatedAt:        time.Now(),
		}

		if err := s.logAddressChange(ctx, changeLog); err != nil {
			// Логируем ошибку, но не прерываем выполнение
			fmt.Printf("Warning: failed to log address change: %v\n", err)
		}
	}

	return updatedGeo, nil
}

// checkListingAccess проверка прав доступа к объявлению
func (s *SpatialService) checkListingAccess(ctx context.Context, listingID, userID int64) error {
	query := `SELECT user_id FROM marketplace_listings WHERE id = $1`
	var ownerID int64

	err := s.db.QueryRowContext(ctx, query, listingID).Scan(&ownerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.ErrListingNotFound
		}
		return err
	}

	if ownerID != userID {
		return types.ErrAccessDenied
	}

	return nil
}

// getEnhancedListingGeo получение расширенных геоданных объявления
func (s *SpatialService) getEnhancedListingGeo(ctx context.Context, listingID int64) (*types.EnhancedListingGeo, error) {
	query := `
		SELECT 
			id, listing_id, 
			ST_Y(location::geometry) as lat,
			ST_X(location::geometry) as lng,
			COALESCE(ST_Y(blurred_location::geometry), 0) as blurred_lat,
			COALESCE(ST_X(blurred_location::geometry), 0) as blurred_lng,
			geohash, is_precise, blur_radius,
			address_components, formatted_address, geocoding_confidence,
			address_verified, input_method, location_privacy,
			created_at, updated_at
		FROM listings_geo 
		WHERE listing_id = $1`

	var geo types.EnhancedListingGeo
	var lat, lng, blurredLat, blurredLng float64
	var addressComponentsJSON sql.NullString

	err := s.db.QueryRowContext(ctx, query, listingID).Scan(
		&geo.ID,
		&geo.ListingID,
		&lat,
		&lng,
		&blurredLat,
		&blurredLng,
		&geo.Geohash,
		&geo.IsPrecise,
		&geo.BlurRadius,
		&addressComponentsJSON,
		&geo.FormattedAddress,
		&geo.GeocodingConfidence,
		&geo.AddressVerified,
		&geo.InputMethod,
		&geo.LocationPrivacy,
		&geo.CreatedAt,
		&geo.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	geo.Location = types.Point{Lat: lat, Lng: lng}

	if blurredLat != 0 || blurredLng != 0 {
		geo.BlurredLocation = &types.Point{Lat: blurredLat, Lng: blurredLng}
	}

	// Парсим компоненты адреса если есть
	if addressComponentsJSON.Valid && addressComponentsJSON.String != "" {
		// TODO: реализовать парсинг JSON
	}

	return &geo, nil
}

// updateEnhancedListingGeo обновление расширенных геоданных
func (s *SpatialService) updateEnhancedListingGeo(ctx context.Context, geo *types.EnhancedListingGeo) error {
	// Для простоты используем UPSERT
	query := `
		INSERT INTO listings_geo (
			listing_id, location, blurred_location, geohash, is_precise, blur_radius,
			address_components, formatted_address, geocoding_confidence,
			address_verified, input_method, location_privacy, updated_at
		) VALUES (
			$1, ST_SetSRID(ST_MakePoint($2, $3), 4326), 
			CASE WHEN $4 IS NOT NULL AND $5 IS NOT NULL 
				THEN ST_SetSRID(ST_MakePoint($4, $5), 4326) 
				ELSE NULL END,
			$6, $7, $8, $9, $10, $11, $12, $13, $14, $15
		)
		ON CONFLICT (listing_id) 
		DO UPDATE SET
			location = EXCLUDED.location,
			blurred_location = EXCLUDED.blurred_location,
			formatted_address = EXCLUDED.formatted_address,
			geocoding_confidence = EXCLUDED.geocoding_confidence,
			address_verified = EXCLUDED.address_verified,
			input_method = EXCLUDED.input_method,
			location_privacy = EXCLUDED.location_privacy,
			updated_at = EXCLUDED.updated_at
		RETURNING id`

	var blurredLng, blurredLat sql.NullFloat64
	if geo.BlurredLocation != nil {
		blurredLng = sql.NullFloat64{Float64: geo.BlurredLocation.Lng, Valid: true}
		blurredLat = sql.NullFloat64{Float64: geo.BlurredLocation.Lat, Valid: true}
	}

	// Вычисляем geohash - для простоты используем базовую реализацию
	geohash := s.calculateGeohash(geo.Location)

	err := s.db.QueryRowContext(ctx, query,
		geo.ListingID,
		geo.Location.Lng,
		geo.Location.Lat,
		blurredLng,
		blurredLat,
		geohash,
		geo.LocationPrivacy != types.PrivacyExact,
		geo.LocationPrivacy.CalculateBlurRadius(),
		nil, // address_components JSON - пока nil
		geo.FormattedAddress,
		geo.GeocodingConfidence,
		geo.AddressVerified,
		string(geo.InputMethod),
		string(geo.LocationPrivacy),
		time.Now(),
	).Scan(&geo.ID)

	return err
}

// logAddressChange логирование изменения адреса
func (s *SpatialService) logAddressChange(ctx context.Context, log *types.AddressChangeLog) error {
	query := `
		INSERT INTO address_change_log (
			listing_id, user_id, old_address, new_address,
			old_location, new_location, change_reason,
			confidence_before, confidence_after, ip_address, user_agent, created_at
		) VALUES (
			$1, $2, $3, $4,
			CASE WHEN $5 IS NOT NULL AND $6 IS NOT NULL 
				THEN ST_SetSRID(ST_MakePoint($5, $6), 4326) 
				ELSE NULL END,
			CASE WHEN $7 IS NOT NULL AND $8 IS NOT NULL 
				THEN ST_SetSRID(ST_MakePoint($7, $8), 4326) 
				ELSE NULL END,
			$9, $10, $11, $12, $13, $14
		)`

	var oldLng, oldLat, newLng, newLat sql.NullFloat64

	if log.OldLocation != nil {
		oldLng = sql.NullFloat64{Float64: log.OldLocation.Lng, Valid: true}
		oldLat = sql.NullFloat64{Float64: log.OldLocation.Lat, Valid: true}
	}

	if log.NewLocation != nil {
		newLng = sql.NullFloat64{Float64: log.NewLocation.Lng, Valid: true}
		newLat = sql.NullFloat64{Float64: log.NewLocation.Lat, Valid: true}
	}

	_, err := s.db.ExecContext(ctx, query,
		log.ListingID,
		log.UserID,
		log.OldAddress,
		log.NewAddress,
		oldLng,
		oldLat,
		newLng,
		newLat,
		log.ChangeReason,
		log.ConfidenceBefore,
		log.ConfidenceAfter,
		log.IPAddress,
		log.UserAgent,
		log.CreatedAt,
	)

	return err
}

// calculateBlurredLocation вычисление размытой локации
func (s *SpatialService) calculateBlurredLocation(location types.Point, privacy types.LocationPrivacyLevel) *types.Point {
	if privacy == types.PrivacyExact {
		return nil // Возвращаем точную локацию
	}

	// Получаем радиус размытия
	blurRadius := privacy.CalculateBlurRadius()

	// Генерируем случайное смещение в пределах радиуса
	// Используем равномерное распределение по площади круга
	angle := rand.Float64() * 2 * math.Pi
	distance := math.Sqrt(rand.Float64()) * blurRadius

	// Преобразуем метры в градусы (приблизительно)
	offsetLat := (distance * math.Cos(angle)) / 111000 // ~111км на градус широты
	offsetLng := (distance * math.Sin(angle)) / (111000 * math.Cos(location.Lat*math.Pi/180))

	blurredLocation := &types.Point{
		Lat: location.Lat + offsetLat,
		Lng: location.Lng + offsetLng,
	}

	return blurredLocation
}

// calculateConfidenceScore вычисление показателя доверия
func (s *SpatialService) calculateConfidenceScore(req types.UpdateAddressRequest) float64 {
	baseScore := 0.5

	// Увеличиваем доверие в зависимости от метода ввода
	switch req.InputMethod {
	case types.InputGeocoded:
		baseScore += 0.3
	case types.InputMapClick:
		baseScore += 0.2
	case types.InputCurrentLocation:
		baseScore += 0.4
	case types.InputManual:
		baseScore += 0.1
	}

	// Увеличиваем если адрес подтвержден пользователем
	if req.Verified {
		baseScore += 0.2
	}

	// Увеличиваем если есть структурированные компоненты адреса
	if req.AddressComponents != nil && req.AddressComponents.HouseNumber != "" {
		baseScore += 0.1
	}

	// Ограничиваем значение
	if baseScore > 1.0 {
		baseScore = 1.0
	}
	if baseScore < 0.0 {
		baseScore = 0.0
	}

	return baseScore
}

// calculateGeohash базовая реализация geohash
func (s *SpatialService) calculateGeohash(point types.Point) string {
	// Упрощенная версия - в реальности стоит использовать библиотеку
	return fmt.Sprintf("gh%.6f%.6f", point.Lat, point.Lng)
}

// Helper функции

func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}
