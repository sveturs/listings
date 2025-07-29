package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// GeocodingService интерфейс для сервиса геокодирования
type GeocodingService interface {
	GetStreetCoordinates(ctx context.Context, address string) (lat, lng float64, err error)
	GetDistrictCoordinates(ctx context.Context, city, district string) (lat, lng float64, err error)
	GetCityCoordinates(ctx context.Context, city, country string) (lat, lng float64, err error)
}

// NominatimGeocoding реализация геокодирования через Nominatim (OpenStreetMap)
type NominatimGeocoding struct{}

// NewNominatimGeocoding создает новый сервис геокодирования
func NewNominatimGeocoding() GeocodingService {
	return &NominatimGeocoding{}
}

// geocode выполняет геокодирование адреса
func (n *NominatimGeocoding) geocode(address string) (lat, lng float64, err error) {
	// Кодируем адрес для URL
	encodedAddress := url.QueryEscape(address)
	
	// Формируем URL запроса к Nominatim
	nominatimURL := fmt.Sprintf(
		"https://nominatim.openstreetmap.org/search?format=json&q=%s&limit=1",
		encodedAddress,
	)

	client := &http.Client{}
	req, err := http.NewRequest("GET", nominatimURL, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create request: %w", err)
	}

	// Добавляем User-Agent, как требует OSM
	req.Header.Add("User-Agent", "SveTu-Platform/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read response: %w", err)
	}

	var results []map[string]interface{}
	err = json.Unmarshal(body, &results)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(results) == 0 {
		return 0, 0, fmt.Errorf("no results found for address: %s", address)
	}

	// Получаем координаты из первого результата
	result := results[0]
	
	latStr, ok := result["lat"].(string)
	if !ok {
		return 0, 0, fmt.Errorf("invalid latitude in response")
	}
	
	lngStr, ok := result["lon"].(string)
	if !ok {
		return 0, 0, fmt.Errorf("invalid longitude in response")
	}

	// Преобразуем строки в числа
	fmt.Sscanf(latStr, "%f", &lat)
	fmt.Sscanf(lngStr, "%f", &lng)

	return lat, lng, nil
}

// GetStreetCoordinates получает координаты улицы (без номера дома)
func (n *NominatimGeocoding) GetStreetCoordinates(ctx context.Context, address string) (lat, lng float64, err error) {
	// Удаляем номер дома из адреса
	parts := strings.Split(address, ",")
	if len(parts) > 0 {
		// Убираем числа из первой части адреса
		streetPart := removeHouseNumber(parts[0])
		if streetPart != "" && len(parts) > 1 {
			// Собираем адрес обратно без номера дома
			parts[0] = streetPart
			address = strings.Join(parts, ",")
		}
	}

	return n.geocode(address)
}

// GetDistrictCoordinates получает координаты центра района
func (n *NominatimGeocoding) GetDistrictCoordinates(ctx context.Context, city, district string) (lat, lng float64, err error) {
	// Формируем запрос для поиска района
	query := fmt.Sprintf("%s, %s", district, city)
	return n.geocode(query)
}

// GetCityCoordinates получает координаты центра города
func (n *NominatimGeocoding) GetCityCoordinates(ctx context.Context, city, country string) (lat, lng float64, err error) {
	// Формируем запрос для поиска города
	query := fmt.Sprintf("%s, %s", city, country)
	return n.geocode(query)
}

// removeHouseNumber удаляет номер дома из строки адреса
func removeHouseNumber(street string) string {
	// Простая реализация - удаляем числа в начале или конце строки
	parts := strings.Fields(street)
	var result []string
	
	for _, part := range parts {
		// Проверяем, не является ли это числом или числом с буквой
		isNumber := true
		for _, r := range part {
			if (r < '0' || r > '9') && r != 'a' && r != 'b' && r != 'а' && r != 'б' {
				isNumber = false
				break
			}
		}
		
		if !isNumber {
			result = append(result, part)
		}
	}
	
	return strings.Join(result, " ")
}

// GetCoordinatesWithGeocoding получает координаты с учетом уровня приватности используя геокодирование
func GetCoordinatesWithGeocoding(ctx context.Context, lat, lng float64, address string, privacyLevel string, geocoder GeocodingService) (float64, float64, error) {
	switch privacyLevel {
	case "exact":
		// Возвращаем точные координаты
		return lat, lng, nil

	case "street", "approximate":
		// Получаем координаты улицы через геокодирование
		if address != "" && geocoder != nil {
			streetLat, streetLng, err := geocoder.GetStreetCoordinates(ctx, address)
			if err == nil {
				return streetLat, streetLng, nil
			}
			// Если геокодирование не удалось, возвращаем округленные координаты
			fmt.Printf("Failed to geocode street address: %v, using rounded coordinates\n", err)
		}
		// Fallback - округляем координаты
		return roundToDecimalPlaces(lat, 3), roundToDecimalPlaces(lng, 3), nil

	case "district":
		// Извлекаем район и город из адреса
		if address != "" && geocoder != nil {
			parts := strings.Split(address, ",")
			if len(parts) >= 2 {
				// Предполагаем формат: улица, район/город, страна
				city := strings.TrimSpace(parts[len(parts)-2])
				district := city // В простом случае используем город как район
				
				districtLat, districtLng, err := geocoder.GetDistrictCoordinates(ctx, city, district)
				if err == nil {
					return districtLat, districtLng, nil
				}
				fmt.Printf("Failed to geocode district: %v, using rounded coordinates\n", err)
			}
		}
		// Fallback - округляем координаты
		return roundToDecimalPlaces(lat, 2), roundToDecimalPlaces(lng, 2), nil

	case "city", "city_only":
		// Извлекаем город из адреса
		if address != "" && geocoder != nil {
			parts := strings.Split(address, ",")
			if len(parts) >= 2 {
				city := strings.TrimSpace(parts[len(parts)-2])
				country := "Serbia" // По умолчанию Сербия
				if len(parts) >= 3 {
					country = strings.TrimSpace(parts[len(parts)-1])
				}
				
				cityLat, cityLng, err := geocoder.GetCityCoordinates(ctx, city, country)
				if err == nil {
					return cityLat, cityLng, nil
				}
				fmt.Printf("Failed to geocode city: %v, using rounded coordinates\n", err)
			}
		}
		// Fallback - округляем координаты
		return roundToDecimalPlaces(lat, 1), roundToDecimalPlaces(lng, 1), nil

	case "hidden":
		// Возвращаем нулевые координаты
		return 0, 0, nil

	default:
		return lat, lng, nil
	}
}

// Функция roundToDecimalPlaces уже определена в address_privacy.go