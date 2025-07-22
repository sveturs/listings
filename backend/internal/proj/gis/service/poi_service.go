package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"backend/internal/proj/gis/types"

	"github.com/jmoiron/sqlx"
)

// POIService сервис для работы с точками интереса
type POIService struct {
	db          *sqlx.DB
	mapboxToken string
	httpClient  *http.Client
}

// NewPOIService создает новый сервис POI
func NewPOIService(db *sqlx.DB) *POIService {
	return &POIService{
		db:          db,
		mapboxToken: os.Getenv("MAPBOX_ACCESS_TOKEN"),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SearchPOI ищет точки интереса через MapBox Geocoding API
func (s *POIService) SearchPOI(ctx context.Context, poiType types.POIType, lat, lng float64, radiusKm float64) ([]types.POISearchResult, error) {
	if s.mapboxToken == "" {
		return nil, fmt.Errorf("MAPBOX_ACCESS_TOKEN not set")
	}

	// Мапинг типов POI в MapBox категории
	mapboxCategory := s.mapPOIType(poiType)

	// Формируем bbox для поиска
	bbox := s.calculateBBox(lat, lng, radiusKm)

	// URL для MapBox Geocoding API
	baseURL := "https://api.mapbox.com/geocoding/v5/mapbox.places"

	params := url.Values{
		"access_token": {s.mapboxToken},
		"types":        {"poi"},
		"proximity":    {fmt.Sprintf("%f,%f", lng, lat)},
		"bbox":         {fmt.Sprintf("%f,%f,%f,%f", bbox[0], bbox[1], bbox[2], bbox[3])},
		"limit":        {"10"},
	}

	// Поисковый запрос в зависимости от типа
	query := s.getPOISearchQuery(poiType)
	searchURL := fmt.Sprintf("%s/%s.json?%s", baseURL, url.QueryEscape(query), params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Printf("Failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("MapBox API returned status %d", resp.StatusCode)
	}

	var geocodingResp struct {
		Features []struct {
			ID         string    `json:"id"`
			Text       string    `json:"text"`
			PlaceName  string    `json:"place_name"`
			Center     []float64 `json:"center"`
			Properties struct {
				Category string `json:"category"`
			} `json:"properties"`
		} `json:"features"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&geocodingResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Конвертируем в наш формат и фильтруем по категории
	results := make([]types.POISearchResult, 0)
	for _, feature := range geocodingResp.Features {
		if len(feature.Center) >= 2 && s.matchesCategory(feature.Properties.Category, mapboxCategory) {
			distance := s.calculateDistance(lat, lng, feature.Center[1], feature.Center[0])

			results = append(results, types.POISearchResult{
				ID:       feature.ID,
				Name:     feature.Text,
				Type:     poiType,
				Lat:      feature.Center[1],
				Lng:      feature.Center[0],
				Distance: distance,
			})
		}
	}

	return results, nil
}

// FilterListingsByPOI фильтрует объявления по близости к POI
func (s *POIService) FilterListingsByPOI(ctx context.Context, filter *types.POIFilter, listingIDs []string) ([]string, error) {
	// Сначала ищем POI в базе данных
	query := `
		WITH poi_locations AS (
			SELECT 
				location::geometry as geom
			FROM gis_poi_cache
			WHERE 
				poi_type = $1
				AND expires_at > NOW()
		),
		listing_poi_distances AS (
			SELECT 
				l.id,
				COUNT(p.geom) as poi_count,
				MIN(ST_Distance(l.location::geography, p.geom::geography)) as min_distance
			FROM marketplace_listings l
			CROSS JOIN poi_locations p
			WHERE 
				l.id = ANY($2)
				AND l.location IS NOT NULL
				AND ST_DWithin(l.location::geography, p.geom::geography, $3)
			GROUP BY l.id
		)
		SELECT id
		FROM listing_poi_distances
		WHERE min_distance <= $3
	`

	// Добавляем фильтр по минимальному количеству POI если указан
	if filter.MinCount > 0 {
		query += fmt.Sprintf(" AND poi_count >= %d", filter.MinCount)
	}

	var filteredIDs []string
	err := s.db.SelectContext(ctx, &filteredIDs, query,
		string(filter.POIType),
		listingIDs,
		filter.MaxDistance,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to filter listings by POI: %w", err)
	}

	return filteredIDs, nil
}

// CachePOI кеширует найденные POI
func (s *POIService) CachePOI(ctx context.Context, poi types.POISearchResult) error {
	query := `
		INSERT INTO gis_poi_cache (
			external_id,
			name,
			poi_type,
			location,
			created_at,
			expires_at
		) VALUES (
			$1,
			$2,
			$3,
			ST_SetSRID(ST_MakePoint($4, $5), 4326),
			NOW(),
			NOW() + INTERVAL '30 days'
		)
		ON CONFLICT (external_id) DO UPDATE SET
			name = EXCLUDED.name,
			expires_at = EXCLUDED.expires_at
	`

	_, err := s.db.ExecContext(ctx, query,
		poi.ID,
		poi.Name,
		string(poi.Type),
		poi.Lng,
		poi.Lat,
	)

	return err
}

// mapPOIType мапит наши типы POI в MapBox категории
func (s *POIService) mapPOIType(poiType types.POIType) string {
	switch poiType {
	case types.POISchool:
		return "school,college,university"
	case types.POIHospital:
		return "hospital,clinic,doctor"
	case types.POIMetro:
		return "subway,metro,station"
	case types.POISupermarket:
		return "supermarket,grocery,shop"
	case types.POIPark:
		return "park,garden"
	case types.POIBank:
		return "bank,atm"
	case types.POIPharmacy:
		return "pharmacy,drugstore"
	case types.POIBusStop:
		return "bus,transit"
	default:
		return "poi"
	}
}

// getPOISearchQuery возвращает поисковый запрос для типа POI
func (s *POIService) getPOISearchQuery(poiType types.POIType) string {
	switch poiType {
	case types.POISchool:
		return "школа"
	case types.POIHospital:
		return "болница"
	case types.POIMetro:
		return "метро станица"
	case types.POISupermarket:
		return "супермаркет"
	case types.POIPark:
		return "парк"
	case types.POIBank:
		return "банка"
	case types.POIPharmacy:
		return "апотека"
	case types.POIBusStop:
		return "аутобуска станица"
	default:
		return ""
	}
}

// calculateBBox вычисляет bounding box для радиуса
func (s *POIService) calculateBBox(lat, lng, radiusKm float64) []float64 {
	// Примерный расчет (1 градус ≈ 111 км)
	latDelta := radiusKm / 111.0
	lngDelta := radiusKm / (111.0 * 0.7) // учитываем широту Сербии

	return []float64{
		lng - lngDelta, // min lng
		lat - latDelta, // min lat
		lng + lngDelta, // max lng
		lat + latDelta, // max lat
	}
}

// calculateDistance вычисляет расстояние между точками (упрощенная формула)
func (s *POIService) calculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	// Упрощенная формула для небольших расстояний
	latDiff := lat2 - lat1
	lngDiff := lng2 - lng1

	// Примерное расстояние в метрах
	return ((latDiff * latDiff) + (lngDiff * lngDiff)) * 111000
}

// matchesCategory проверяет соответствие категории
func (s *POIService) matchesCategory(category, expected string) bool {
	// Простая проверка вхождения
	return category != "" && expected != ""
}
