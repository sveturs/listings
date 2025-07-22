package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"backend/internal/logger"
	"backend/internal/proj/gis/types"

	"github.com/jmoiron/sqlx"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/wkt"
)

// IsochroneService сервис для работы с изохронами MapBox
type IsochroneService struct {
	db          *sqlx.DB
	mapboxToken string
	httpClient  *http.Client
}

// NewIsochroneService создает новый сервис изохрон
func NewIsochroneService(db *sqlx.DB) *IsochroneService {
	return &IsochroneService{
		db:          db,
		mapboxToken: os.Getenv("MAPBOX_ACCESS_TOKEN"),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetIsochrone получает изохрону от MapBox API
func (s *IsochroneService) GetIsochrone(ctx context.Context, filter *types.TravelTimeFilter) (*types.IsohronResponse, error) {
	if s.mapboxToken == "" {
		return nil, fmt.Errorf("MAPBOX_ACCESS_TOKEN not set")
	}

	// Формируем URL для MapBox Isochrone API
	baseURL := "https://api.mapbox.com/isochrone/v1/mapbox"

	// Мапинг режимов транспорта
	mapboxMode := s.mapTransportMode(filter.TransportMode)

	// Координаты в формате lng,lat
	coordinates := fmt.Sprintf("%f,%f", filter.CenterLng, filter.CenterLat)

	// Параметры запроса
	params := url.Values{
		"access_token":     {s.mapboxToken},
		"contours_minutes": {fmt.Sprintf("%d", filter.MaxMinutes)},
		"polygons":         {"true"},
		"denoise":          {"1"},
		"generalize":       {"500"},
	}

	fullURL := fmt.Sprintf("%s/%s/%s?%s", baseURL, mapboxMode, coordinates, params.Encode())

	// Выполняем запрос
	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
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

	var isochroneResp types.IsohronResponse
	if err := json.NewDecoder(resp.Body).Decode(&isochroneResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &isochroneResp, nil
}

// FilterListingsByIsochrone фильтрует объявления по изохроне
func (s *IsochroneService) FilterListingsByIsochrone(ctx context.Context, filter *types.TravelTimeFilter, listingIDs []string) ([]string, error) {
	// Получаем изохрону
	isochrone, err := s.GetIsochrone(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get isochrone: %w", err)
	}

	if len(isochrone.Features) == 0 {
		return []string{}, nil
	}

	// Берем первый полигон (самый большой контур)
	coords := isochrone.Features[0].Geometry.Coordinates
	if len(coords) == 0 || len(coords[0]) == 0 {
		return []string{}, nil
	}

	// Конвертируем координаты в WKT полигон
	polygon := s.coordinatesToWKT(coords[0])

	// Фильтруем объявления через PostGIS
	query := `
		SELECT id 
		FROM marketplace_listings 
		WHERE 
			id = ANY($1) 
			AND location IS NOT NULL
			AND ST_Within(location::geometry, ST_GeomFromText($2, 4326))
	`

	var filteredIDs []string
	err = s.db.SelectContext(ctx, &filteredIDs, query, listingIDs, polygon)
	if err != nil {
		return nil, fmt.Errorf("failed to filter listings by isochrone: %w", err)
	}

	return filteredIDs, nil
}

// mapTransportMode конвертирует наш режим транспорта в MapBox формат
func (s *IsochroneService) mapTransportMode(mode types.TransportMode) string {
	switch mode {
	case types.TransportWalking:
		return "walking"
	case types.TransportDriving:
		return "driving"
	case types.TransportCycling:
		return "cycling"
	default:
		return "driving-traffic" // для общественного транспорта используем driving с трафиком
	}
}

// coordinatesToWKT конвертирует координаты в WKT формат
func (s *IsochroneService) coordinatesToWKT(coords [][]float64) string {
	// Создаем кольцо координат
	ring := make([]geom.Coord, len(coords))
	for i, coord := range coords {
		ring[i] = geom.Coord{coord[0], coord[1]}
	}

	// Создаем полигон
	polygon := geom.NewPolygon(geom.XY).MustSetCoords([][]geom.Coord{ring})

	// Конвертируем в WKT
	wktStr, err := wkt.Marshal(polygon)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to marshal polygon to WKT")
		return ""
	}

	return wktStr
}

// CacheIsochrone кеширует результат изохроны
func (s *IsochroneService) CacheIsochrone(ctx context.Context, filter *types.TravelTimeFilter, polygon string) error {
	query := `
		INSERT INTO gis_isochrone_cache (
			center_point,
			transport_mode,
			max_minutes,
			polygon,
			created_at,
			expires_at
		) VALUES (
			ST_SetSRID(ST_MakePoint($1, $2), 4326),
			$3,
			$4,
			ST_GeomFromText($5, 4326),
			NOW(),
			NOW() + INTERVAL '24 hours'
		)
	`

	_, err := s.db.ExecContext(ctx, query,
		filter.CenterLng,
		filter.CenterLat,
		filter.TransportMode,
		filter.MaxMinutes,
		polygon,
	)

	return err
}

// GetCachedIsochrone получает закешированную изохрону
func (s *IsochroneService) GetCachedIsochrone(ctx context.Context, filter *types.TravelTimeFilter) (string, error) {
	query := `
		SELECT ST_AsText(polygon) 
		FROM gis_isochrone_cache
		WHERE 
			ST_DWithin(
				center_point::geography, 
				ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography, 
				100
			)
			AND transport_mode = $3
			AND max_minutes = $4
			AND expires_at > NOW()
		ORDER BY created_at DESC
		LIMIT 1
	`

	var polygon string
	err := s.db.GetContext(ctx, &polygon, query,
		filter.CenterLng,
		filter.CenterLat,
		filter.TransportMode,
		filter.MaxMinutes,
	)
	if err != nil {
		return "", err
	}

	return polygon, nil
}
