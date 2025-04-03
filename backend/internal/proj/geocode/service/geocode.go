package service

import (
    "backend/internal/domain/models"
    "backend/internal/storage"
    "context"
    "encoding/json"
    "fmt"
    "io/ioutil"
	"net/http"
    "net/url"
    "strconv"
    "strings"
)

type GeocodeService struct {
    storage storage.Storage
}

func NewGeocodeService(storage storage.Storage) GeocodeServiceInterface {
    return &GeocodeService{
        storage: storage,
    }
}

// ReverseGeocode получает адрес по координатам
// Исправленная функция ReverseGeocode в файле backend/internal/service/geocode.go

// ReverseGeocode получает адрес по координатам
func (s *GeocodeService) ReverseGeocode(ctx context.Context, lat, lon float64) (*models.GeoLocation, error) {
    // Используем Nominatim API (OpenStreetMap)
    nominatimURL := fmt.Sprintf(
        "https://nominatim.openstreetmap.org/reverse?format=json&lat=%f&lon=%f&zoom=10&addressdetails=1",
        lat, lon,
    )

    client := &http.Client{}
    req, err := http.NewRequest("GET", nominatimURL, nil)
    if err != nil {
        return nil, err
    }

    // Добавляем User-Agent, как требует OSM
    req.Header.Add("User-Agent", "HostelBookingApp/1.0")
    
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var result map[string]interface{}
    err = json.Unmarshal(body, &result)
    if err != nil {
        return nil, err
    }

    // Извлекаем данные из ответа
    address, ok := result["address"].(map[string]interface{})
    if !ok {
        // Вместо возврата ошибки используем значения по умолчанию
        return &models.GeoLocation{
            City:    "Unknown City",
            Country: "Unknown Country",
            Lat:     lat,
            Lon:     lon,
        }, nil
    }

    // Пытаемся найти город
    var city, country string
    
    // Проверяем разные поля, так как структура зависит от региона
    for _, field := range []string{"city", "town", "village", "hamlet", "county", "state", "region"} {
        if value, ok := address[field].(string); ok && value != "" {
            city = value
            break
        }
    }

    // Если город не найден, используем display_name
    if city == "" && result["display_name"] != nil {
        displayName, ok := result["display_name"].(string)
        if ok {
            parts := strings.Split(displayName, ",")
            if len(parts) > 0 {
                city = strings.TrimSpace(parts[0])
            }
        }
    }

    // Если город всё еще не найден, используем значение по умолчанию
    if city == "" {
        city = "Unknown City"
    }

    // Определяем страну
    if value, ok := address["country"].(string); ok {
        country = value
    } else {
        country = "Unknown Country"
    }

    // Создаем модель локации
    location := &models.GeoLocation{
        City:    city,
        Country: country,
        Lat:     lat,
        Lon:     lon,
    }

    return location, nil
}

// GetCitySuggestions возвращает предложения городов по частичному названию
func (s *GeocodeService) GetCitySuggestions(ctx context.Context, query string, limit int) ([]models.GeoLocation, error) {
    // Очищаем строку запроса
    query = strings.TrimSpace(query)
    query = url.QueryEscape(query)

    // Запрос в Nominatim API для поиска городов
    nominatimURL := fmt.Sprintf(
        "https://nominatim.openstreetmap.org/search?format=json&q=%s&limit=%d&addressdetails=1&accept-language=ru,en",
        query, limit,
    )

    client := &http.Client{}
    req, err := http.NewRequest("GET", nominatimURL, nil)
    if err != nil {
        return nil, err
    }

    // Добавляем User-Agent, как требует OSM
    req.Header.Add("User-Agent", "HostelBookingApp/1.0")
    
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var results []map[string]interface{}
    err = json.Unmarshal(body, &results)
    if err != nil {
        return nil, err
    }

    // Парсим результаты
    var suggestions []models.GeoLocation
    for _, result := range results {
        // Проверяем тип места - нам нужны только города
        placeType, _ := result["type"].(string)
        if !isCity(placeType) {
            continue
        }

        var city, country string
        lat, lon := 0.0, 0.0

        // Извлекаем широту и долготу
        if latStr, ok := result["lat"].(string); ok {
            lat, _ = strconv.ParseFloat(latStr, 64)
        }
        if lonStr, ok := result["lon"].(string); ok {
            lon, _ = strconv.ParseFloat(lonStr, 64)
        }

        // Получаем адрес
        address, _ := result["address"].(map[string]interface{})
        if address != nil {
            // Ищем город
            for _, field := range []string{"city", "town", "village"} {
                if value, ok := address[field].(string); ok && value != "" {
                    city = value
                    break
                }
            }

            // Если город не найден, пробуем использовать имя места
            if city == "" {
                if name, ok := result["display_name"].(string); ok {
                    parts := strings.Split(name, ",")
                    if len(parts) > 0 {
                        city = strings.TrimSpace(parts[0])
                    }
                }
            }

            // Определяем страну
            if value, ok := address["country"].(string); ok {
                country = value
            }
        }

        // Если город все еще не найден, пропускаем
        if city == "" {
            continue
        }

        suggestion := models.GeoLocation{
            City:    city,
            Country: country,
            Lat:     lat,
            Lon:     lon,
        }

        suggestions = append(suggestions, suggestion)
    }

    return suggestions, nil
}

// isCity проверяет, является ли результат городом или похожим населенным пунктом
func isCity(placeType string) bool {
    cityTypes := []string{"city", "town", "village", "hamlet", "administrative"}
    for _, t := range cityTypes {
        if placeType == t {
            return true
        }
    }
    return false
}