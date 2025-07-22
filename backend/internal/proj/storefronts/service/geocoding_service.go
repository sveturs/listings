package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"backend/internal/domain/models"
)

// GeocodingService интерфейс сервиса геокодирования
type GeocodingService interface {
	// GeocodeAddress преобразует адрес в координаты
	GeocodeAddress(ctx context.Context, address string) (*models.Location, error)

	// ReverseGeocode преобразует координаты в адрес
	ReverseGeocode(ctx context.Context, lat, lng float64) (*models.Location, error)

	// SmartGeocode умный геокодинг с определением ближайшего здания
	SmartGeocode(ctx context.Context, userLat, userLng float64) (*models.Location, error)

	// ValidateAddress проверяет корректность адреса
	ValidateAddress(ctx context.Context, address string) (bool, error)
}

// nominatimResponse ответ от Nominatim API
type nominatimResponse struct {
	PlaceID     string   `json:"place_id"`
	Lat         string   `json:"lat"`
	Lon         string   `json:"lon"`
	DisplayName string   `json:"display_name"`
	Address     address  `json:"address"`
	BoundingBox []string `json:"boundingbox"`
}

type address struct {
	HouseNumber  string `json:"house_number"`
	Road         string `json:"road"`
	Suburb       string `json:"suburb"`
	City         string `json:"city"`
	Municipality string `json:"municipality"`
	State        string `json:"state"`
	Postcode     string `json:"postcode"`
	Country      string `json:"country"`
	CountryCode  string `json:"country_code"`
}

// geocodingService реализация сервиса геокодирования
type geocodingService struct {
	httpClient *http.Client
	userAgent  string
	cache      map[string]*cacheEntry // Простой in-memory кеш
}

type cacheEntry struct {
	data      *models.Location
	expiresAt time.Time
}

// NewGeocodingService создает новый сервис геокодирования
func NewGeocodingService() GeocodingService {
	return &geocodingService{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		userAgent: "SVE-TU-Platform/1.0",
		cache:     make(map[string]*cacheEntry),
	}
}

// GeocodeAddress преобразует адрес в координаты
func (g *geocodingService) GeocodeAddress(ctx context.Context, address string) (*models.Location, error) {
	// Проверяем кеш
	cacheKey := "geocode:" + address
	if cached, ok := g.cache[cacheKey]; ok && cached.expiresAt.After(time.Now()) {
		return cached.data, nil
	}

	// Формируем запрос к Nominatim
	params := url.Values{}
	params.Set("q", address)
	params.Set("format", "json")
	params.Set("limit", "1")
	params.Set("addressdetails", "1")
	params.Set("accept-language", "sr,en") // Сербский приоритет

	reqURL := fmt.Sprintf("https://nominatim.openstreetmap.org/search?%s", params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", g.userAgent)

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to geocode: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			// Логирование ошибки закрытия Body
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geocoding failed with status %d", resp.StatusCode)
	}

	var results []nominatimResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results found for address: %s", address)
	}

	location := g.parseNominatimResponse(&results[0])

	// Сохраняем в кеш
	g.cache[cacheKey] = &cacheEntry{
		data:      location,
		expiresAt: time.Now().Add(24 * time.Hour),
	}

	// Задержка для соблюдения лимитов Nominatim (1 запрос в секунду)
	time.Sleep(time.Second)

	return location, nil
}

// ReverseGeocode преобразует координаты в адрес
func (g *geocodingService) ReverseGeocode(ctx context.Context, lat, lng float64) (*models.Location, error) {
	// Проверяем кеш
	cacheKey := fmt.Sprintf("reverse:%.6f,%.6f", lat, lng)
	if cached, ok := g.cache[cacheKey]; ok && cached.expiresAt.After(time.Now()) {
		return cached.data, nil
	}

	// Формируем запрос к Nominatim
	params := url.Values{}
	params.Set("lat", fmt.Sprintf("%.8f", lat))
	params.Set("lon", fmt.Sprintf("%.8f", lng))
	params.Set("format", "json")
	params.Set("addressdetails", "1")
	params.Set("accept-language", "sr,en")

	reqURL := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?%s", params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", g.userAgent)

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to reverse geocode: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			// Логирование ошибки закрытия Body
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("reverse geocoding failed with status %d", resp.StatusCode)
	}

	var result nominatimResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	location := g.parseNominatimResponse(&result)
	location.UserLat = lat
	location.UserLng = lng

	// Сохраняем в кеш
	g.cache[cacheKey] = &cacheEntry{
		data:      location,
		expiresAt: time.Now().Add(24 * time.Hour),
	}

	// Задержка для соблюдения лимитов
	time.Sleep(time.Second)

	return location, nil
}

// SmartGeocode умный геокодинг
func (g *geocodingService) SmartGeocode(ctx context.Context, userLat, userLng float64) (*models.Location, error) {
	// 1. Получаем адрес по координатам клика
	location, err := g.ReverseGeocode(ctx, userLat, userLng)
	if err != nil {
		return nil, fmt.Errorf("reverse geocode failed: %w", err)
	}

	// 2. Если есть номер дома, получаем точные координаты здания
	if location.HouseNumber != "" && location.Street != "" {
		buildingAddress := fmt.Sprintf("%s %s, %s, %s",
			location.Street, location.HouseNumber, location.City, location.Country)

		buildingLocation, err := g.GeocodeAddress(ctx, buildingAddress)
		if err == nil {
			// Сохраняем оригинальные координаты клика
			buildingLocation.UserLat = userLat
			buildingLocation.UserLng = userLng
			return buildingLocation, nil
		}
	}

	// 3. Если не удалось найти точное здание, возвращаем оригинальные данные
	return location, nil
}

// ValidateAddress проверяет адрес
func (g *geocodingService) ValidateAddress(ctx context.Context, address string) (bool, error) {
	location, err := g.GeocodeAddress(ctx, address)
	if err != nil {
		return false, nil
	}

	// Проверяем что нашли хотя бы город
	return location.City != "", nil
}

// parseNominatimResponse парсит ответ от Nominatim
func (g *geocodingService) parseNominatimResponse(resp *nominatimResponse) *models.Location {
	lat := 0.0
	lng := 0.0
	if _, err := fmt.Sscanf(resp.Lat, "%f", &lat); err != nil {
		// Не удалось распарсить широту, оставляем 0.0
	}
	if _, err := fmt.Sscanf(resp.Lon, "%f", &lng); err != nil {
		// Не удалось распарсить долготу, оставляем 0.0
	}

	// Определяем город (может быть в разных полях)
	city := resp.Address.City
	if city == "" {
		city = resp.Address.Municipality
	}

	// Формируем полный адрес
	addressParts := []string{}
	if resp.Address.Road != "" {
		addressParts = append(addressParts, resp.Address.Road)
		if resp.Address.HouseNumber != "" {
			addressParts = append(addressParts, resp.Address.HouseNumber)
		}
	}

	fullAddress := strings.Join(addressParts, " ")
	if city != "" {
		if fullAddress != "" {
			fullAddress += ", "
		}
		fullAddress += city
	}

	if resp.Address.Postcode != "" {
		fullAddress += " " + resp.Address.Postcode
	}

	// Определяем код страны для Балкан
	countryCode := strings.ToUpper(resp.Address.CountryCode)
	if countryCode == "" {
		// Определяем по названию страны
		switch strings.ToLower(resp.Address.Country) {
		case "србија", "serbia":
			countryCode = "RS"
		case "hrvatska", "croatia":
			countryCode = "HR"
		case "bosna i hercegovina", "bosnia and herzegovina":
			countryCode = "BA"
		case "crna gora", "montenegro":
			countryCode = "ME"
		case "slovenija", "slovenia":
			countryCode = "SI"
		case "северна македонија", "north macedonia":
			countryCode = "MK"
		default:
			countryCode = "RS" // По умолчанию Сербия
		}
	}

	return &models.Location{
		UserLat:     lat,
		UserLng:     lng,
		BuildingLat: lat,
		BuildingLng: lng,
		FullAddress: fullAddress,
		Street:      resp.Address.Road,
		HouseNumber: resp.Address.HouseNumber,
		PostalCode:  resp.Address.Postcode,
		City:        city,
		Country:     countryCode,
		BuildingInfo: models.JSONB{
			"display_name": resp.DisplayName,
			"place_id":     resp.PlaceID,
			"suburb":       resp.Address.Suburb,
			"state":        resp.Address.State,
		},
	}
}

// Вспомогательные функции для работы с сербскими адресами

// normalizeAddress нормализует адрес для Балкан
func normalizeAddress(address string) string {
	// Заменяем распространенные сокращения
	replacements := map[string]string{
		"ул.":   "улица",
		"бул.":  "булевар",
		"тргпј": "трг",
		"бр.":   "",
		"br.":   "",
		"ul.":   "ulica",
		"bul.":  "bulevar",
	}

	normalized := strings.ToLower(address)
	for old, new := range replacements {
		normalized = strings.ReplaceAll(normalized, old, new)
	}

	return normalized
}
