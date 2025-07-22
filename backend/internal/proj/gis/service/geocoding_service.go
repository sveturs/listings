package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"backend/internal/logger"
	"backend/internal/proj/gis/repository"
	"backend/internal/proj/gis/types"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

const (
	mapboxGeocodingAPIURL = "https://api.mapbox.com/geocoding/v5/mapbox.places"
)

// GeocodingService сервис для геокодирования
type GeocodingService struct {
	db             *sqlx.DB
	cacheRepo      *repository.GeocodingCacheRepository
	mapboxToken    string
	httpClient     *http.Client
	maxSuggestions int
	defaultTTL     time.Duration
	logger         zerolog.Logger
}

// NewGeocodingService создает новый сервис геокодирования
func NewGeocodingService(db *sqlx.DB) *GeocodingService {
	mapboxToken := os.Getenv("MAPBOX_ACCESS_TOKEN")
	if mapboxToken == "" {
		// Логируем предупреждение, но не останавливаем сервис
		fmt.Println("WARNING: MAPBOX_ACCESS_TOKEN not set, geocoding will use fallback methods")
	}

	maxSuggestions := 5
	if val := os.Getenv("GEOCODING_MAX_SUGGESTIONS"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil && parsed > 0 {
			maxSuggestions = parsed
		}
	}

	ttlHours := 720 // 30 дней по умолчанию
	if val := os.Getenv("GEOCODING_CACHE_TTL_HOURS"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil && parsed > 0 {
			ttlHours = parsed
		}
	}

	return &GeocodingService{
		db:             db,
		cacheRepo:      repository.NewGeocodingCacheRepository(db),
		mapboxToken:    mapboxToken,
		httpClient:     &http.Client{Timeout: 10 * time.Second},
		maxSuggestions: maxSuggestions,
		defaultTTL:     time.Duration(ttlHours) * time.Hour,
		logger:         logger.Get().With().Str("service", "geocoding").Logger(),
	}
}

// ValidateAndGeocode валидация и геокодирование адреса
func (s *GeocodingService) ValidateAndGeocode(ctx context.Context, req types.GeocodeValidateRequest) (*types.GeocodeValidateResponse, error) {
	// Нормализуем адрес для поиска в кэше
	normalizedAddress := s.normalizeAddress(req.Address)

	// Проверяем кэш
	if cached, err := s.cacheRepo.GetByAddress(ctx, normalizedAddress, req.Language, req.CountryCode); err == nil && cached != nil {
		// Обновляем счетчик попаданий в кэш
		if err := s.cacheRepo.IncrementHits(ctx, cached.ID); err != nil {
			// Логируем ошибку, но не прерываем выполнение
			fmt.Printf("Warning: failed to increment cache hits: %v\n", err)
		}

		return &types.GeocodeValidateResponse{
			Success:           true,
			Location:          &cached.Location,
			AddressComponents: &cached.AddressComponents,
			FormattedAddress:  cached.FormattedAddress,
			Confidence:        cached.Confidence,
		}, nil
	}

	// Выполняем геокодирование через MapBox API
	result, err := s.geocodeWithMapbox(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("geocoding failed: %w", err)
	}

	// Сохраняем в кэш если результат успешный
	if result.Success && result.Location != nil {
		cacheEntry := &types.GeocodingCacheEntry{
			InputAddress:      req.Address,
			NormalizedAddress: normalizedAddress,
			Location:          *result.Location,
			AddressComponents: *result.AddressComponents,
			FormattedAddress:  result.FormattedAddress,
			Confidence:        result.Confidence,
			Provider:          "mapbox",
			Language:          req.Language,
			CountryCode:       req.CountryCode,
			ExpiresAt:         time.Now().Add(s.defaultTTL),
		}

		if err := s.cacheRepo.Save(ctx, cacheEntry); err != nil {
			// Логируем ошибку, но не прерываем выполнение
			fmt.Printf("Warning: failed to save geocoding cache: %v\n", err)
		}
	}

	return result, nil
}

// SearchSuggestions поиск предложений адресов
func (s *GeocodingService) SearchSuggestions(ctx context.Context, query string, limit int, language, countryCode string) ([]types.AddressSuggestion, error) {
	if limit <= 0 || limit > s.maxSuggestions {
		limit = s.maxSuggestions
	}

	// Сначала проверяем кэш на похожие адреса
	cachedSuggestions := s.getCachedSuggestions(ctx, query, language, limit/2)

	// Затем запрашиваем fresh данные из MapBox
	freshSuggestions, err := s.getFreshSuggestions(ctx, query, limit, language, countryCode)
	if err != nil {
		// Если MapBox недоступен, возвращаем только кэшированные результаты
		return cachedSuggestions, nil
	}

	// Объединяем результаты, убирая дубликаты
	suggestions := s.mergeSuggestions(cachedSuggestions, freshSuggestions, limit)

	return suggestions, nil
}

// ReverseGeocode обратное геокодирование
func (s *GeocodingService) ReverseGeocode(ctx context.Context, point types.Point, language string) (*types.AddressSuggestion, error) {
	// Формируем URL для обратного геокодирования
	baseURL := mapboxGeocodingAPIURL
	params := url.Values{}
	params.Set("access_token", s.mapboxToken)
	params.Set("language", language)
	params.Set("limit", "1")

	requestURL := fmt.Sprintf("%s/%.6f,%.6f.json?%s", baseURL, point.Lng, point.Lat, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.logger.Warn().Err(err).Msg("Failed to close response body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("mapbox API returned status %d", resp.StatusCode)
	}

	var mapboxResp MapboxGeocodingResponse
	if err := json.NewDecoder(resp.Body).Decode(&mapboxResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(mapboxResp.Features) == 0 {
		return nil, fmt.Errorf("no address found for coordinates")
	}

	result := s.convertMapboxFeatureToSuggestion(mapboxResp.Features[0])
	return &result, nil
}

// GetCacheStats получение статистики кэша
func (s *GeocodingService) GetCacheStats(ctx context.Context) (map[string]interface{}, error) {
	return s.cacheRepo.GetStats(ctx)
}

// CleanupExpiredCache очистка устаревшего кэша
func (s *GeocodingService) CleanupExpiredCache(ctx context.Context) (int64, error) {
	return s.cacheRepo.CleanupExpired(ctx)
}

// Приватные методы

func (s *GeocodingService) normalizeAddress(address string) string {
	// Приводим адрес к нормализованному виду для кэширования
	normalized := strings.TrimSpace(address)
	normalized = strings.ToLower(normalized)

	// Убираем лишние пробелы и знаки препинания
	normalized = strings.Join(strings.Fields(normalized), " ")

	// Заменяем некоторые сокращения
	replacements := map[string]string{
		" str. ":  " street ",
		" st. ":   " street ",
		" ave. ":  " avenue ",
		" blvd. ": " boulevard ",
		" ul. ":   " ulica ",
		" бул. ":  " булевар ",
		" ул. ":   " улица ",
	}

	for old, new := range replacements {
		normalized = strings.ReplaceAll(normalized, old, new)
	}

	return normalized
}

func (s *GeocodingService) geocodeWithMapbox(ctx context.Context, req types.GeocodeValidateRequest) (*types.GeocodeValidateResponse, error) {
	if s.mapboxToken == "" {
		return &types.GeocodeValidateResponse{
			Success:  false,
			Warnings: []string{"MapBox API token not configured"},
		}, nil
	}

	// Формируем URL для MapBox Geocoding API
	baseURL := mapboxGeocodingAPIURL
	params := url.Values{}
	params.Set("access_token", s.mapboxToken)
	params.Set("limit", "1")

	if req.Language != "" {
		params.Set("language", req.Language)
	}

	if req.CountryCode != "" {
		params.Set("country", strings.ToLower(req.CountryCode))
	}

	// Добавляем типы для более точного поиска
	params.Set("types", "address,poi")

	// Добавляем proximity если указан в контексте
	if req.Context != nil && req.Context.ProximityTo != nil {
		proximity := fmt.Sprintf("%.6f,%.6f", req.Context.ProximityTo.Lng, req.Context.ProximityTo.Lat)
		params.Set("proximity", proximity)
	}

	requestURL := fmt.Sprintf("%s/%s.json?%s", baseURL, url.QueryEscape(req.Address), params.Encode())

	httpReq, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.logger.Warn().Err(err).Msg("Failed to close response body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("mapbox API returned status %d", resp.StatusCode)
	}

	var mapboxResp MapboxGeocodingResponse
	if err := json.NewDecoder(resp.Body).Decode(&mapboxResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(mapboxResp.Features) == 0 {
		return &types.GeocodeValidateResponse{
			Success:  false,
			Warnings: []string{"No results found for the given address"},
		}, nil
	}

	feature := mapboxResp.Features[0]
	suggestion := s.convertMapboxFeatureToSuggestion(feature)

	// Определяем уровень доверия на основе типа места и точности
	confidence := s.calculateConfidence(feature)

	result := &types.GeocodeValidateResponse{
		Success:           true,
		Location:          &suggestion.Location,
		AddressComponents: &suggestion.AddressComponents,
		FormattedAddress:  suggestion.PlaceName,
		Confidence:        confidence,
	}

	// Добавляем предупреждения если нужно
	if confidence < 0.7 {
		result.Warnings = append(result.Warnings, "Low confidence geocoding result")
	}

	return result, nil
}

func (s *GeocodingService) getCachedSuggestions(ctx context.Context, query string, language string, limit int) []types.AddressSuggestion {
	cached, err := s.cacheRepo.SearchSimilar(ctx, query, language, limit)
	if err != nil {
		return nil
	}

	suggestions := make([]types.AddressSuggestion, 0, len(cached))
	for _, entry := range cached {
		suggestion := types.AddressSuggestion{
			ID:                fmt.Sprintf("cache_%d", entry.ID),
			Text:              entry.InputAddress,
			PlaceName:         entry.FormattedAddress,
			Location:          entry.Location,
			AddressComponents: entry.AddressComponents,
			Confidence:        entry.Confidence,
			PlaceTypes:        []string{"cached"},
		}
		suggestions = append(suggestions, suggestion)
	}

	return suggestions
}

func (s *GeocodingService) getFreshSuggestions(ctx context.Context, query string, limit int, language, countryCode string) ([]types.AddressSuggestion, error) {
	if s.mapboxToken == "" {
		return nil, fmt.Errorf("mapbox token not configured")
	}

	// Формируем URL для MapBox Geocoding API
	baseURL := mapboxGeocodingAPIURL
	params := url.Values{}
	params.Set("access_token", s.mapboxToken)
	params.Set("autocomplete", "true")
	params.Set("limit", fmt.Sprintf("%d", limit))
	params.Set("language", language)

	if countryCode != "" {
		params.Set("country", strings.ToLower(countryCode))
	}

	// Типы мест для более точного поиска адресов
	params.Set("types", "address,poi")

	requestURL := fmt.Sprintf("%s/%s.json?%s", baseURL, url.QueryEscape(query), params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.logger.Warn().Err(err).Msg("Failed to close response body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("mapbox API returned status %d", resp.StatusCode)
	}

	var mapboxResp MapboxGeocodingResponse
	if err := json.NewDecoder(resp.Body).Decode(&mapboxResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Конвертируем ответ MapBox в наш формат
	suggestions := make([]types.AddressSuggestion, 0, len(mapboxResp.Features))
	for _, feature := range mapboxResp.Features {
		suggestion := s.convertMapboxFeatureToSuggestion(feature)
		suggestions = append(suggestions, suggestion)
	}

	return suggestions, nil
}

func (s *GeocodingService) mergeSuggestions(cached, fresh []types.AddressSuggestion, limit int) []types.AddressSuggestion {
	// Создаем мапу для избежания дубликатов
	seen := make(map[string]bool)
	var result []types.AddressSuggestion

	// Сначала добавляем кэшированные результаты (они приоритетнее)
	for _, suggestion := range cached {
		key := fmt.Sprintf("%.6f,%.6f", suggestion.Location.Lat, suggestion.Location.Lng)
		if !seen[key] && len(result) < limit {
			seen[key] = true
			result = append(result, suggestion)
		}
	}

	// Затем добавляем fresh результаты
	for _, suggestion := range fresh {
		key := fmt.Sprintf("%.6f,%.6f", suggestion.Location.Lat, suggestion.Location.Lng)
		if !seen[key] && len(result) < limit {
			seen[key] = true
			result = append(result, suggestion)
		}
	}

	return result
}

func (s *GeocodingService) convertMapboxFeatureToSuggestion(feature MapboxFeature) types.AddressSuggestion {
	// Извлекаем координаты
	coords := feature.Geometry.Coordinates
	location := types.Point{
		Lng: coords[0],
		Lat: coords[1],
	}

	// Извлекаем компоненты адреса из контекста
	components := types.AddressComponents{
		Formatted: feature.PlaceName,
	}

	// Парсим контекст для извлечения компонентов
	for _, ctx := range feature.Context {
		switch {
		case strings.HasPrefix(ctx.ID, "country"):
			components.Country = ctx.Text
			if ctx.ShortCode != "" {
				components.CountryCode = strings.ToUpper(ctx.ShortCode)
			}
		case strings.HasPrefix(ctx.ID, "place"):
			components.City = ctx.Text
		case strings.HasPrefix(ctx.ID, "district"):
			components.District = ctx.Text
		case strings.HasPrefix(ctx.ID, "postcode"):
			components.PostalCode = ctx.Text
		case strings.HasPrefix(ctx.ID, "address"):
			components.Street = ctx.Text
		}
	}

	// Извлекаем дополнительную информацию из properties
	if feature.Properties.Address != "" {
		components.HouseNumber = feature.Properties.Address
	}

	return types.AddressSuggestion{
		ID:                feature.ID,
		Text:              feature.Text,
		PlaceName:         feature.PlaceName,
		Location:          location,
		AddressComponents: components,
		Confidence:        s.calculateConfidence(feature),
		PlaceTypes:        feature.PlaceType,
	}
}

func (s *GeocodingService) calculateConfidence(feature MapboxFeature) float64 {
	// Базовый уровень доверия
	confidence := 0.5

	// Увеличиваем доверие в зависимости от типа места
	for _, placeType := range feature.PlaceType {
		switch placeType {
		case "address":
			confidence += 0.4
		case "poi":
			confidence += 0.3
		case "place":
			confidence += 0.2
		case "locality":
			confidence += 0.1
		}
	}

	// Увеличиваем доверие если есть номер дома
	if feature.Properties.Address != "" {
		confidence += 0.1
	}

	// Увеличиваем доверие на основе relevance от MapBox
	if feature.Relevance > 0 {
		confidence += feature.Relevance * 0.2
	}

	// Ограничиваем значение от 0.0 до 1.0
	if confidence > 1.0 {
		confidence = 1.0
	}
	if confidence < 0.0 {
		confidence = 0.0
	}

	return confidence
}

// Структуры для работы с MapBox API
type MapboxGeocodingResponse struct {
	Type     string          `json:"type"`
	Query    []interface{}   `json:"query"`
	Features []MapboxFeature `json:"features"`
}

type MapboxFeature struct {
	ID         string           `json:"id"`
	Type       string           `json:"type"`
	PlaceType  []string         `json:"place_type"`
	Relevance  float64          `json:"relevance"`
	Properties MapboxProperties `json:"properties"`
	Text       string           `json:"text"`
	PlaceName  string           `json:"place_name"`
	Center     []float64        `json:"center"`
	Geometry   MapboxGeometry   `json:"geometry"`
	Context    []MapboxContext  `json:"context"`
}

type MapboxProperties struct {
	Accuracy string `json:"accuracy,omitempty"`
	Address  string `json:"address,omitempty"`
	Category string `json:"category,omitempty"`
	Maki     string `json:"maki,omitempty"`
}

type MapboxGeometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type MapboxContext struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	Wikidata  string `json:"wikidata,omitempty"`
	ShortCode string `json:"short_code,omitempty"`
}
