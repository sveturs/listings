// backend/internal/domain/search/types.go
package search

import (
	"backend/internal/domain/models"
)

// SearchParams параметры поиска для OpenSearch
type SearchParams struct {
	Query              string            // Текст поиска
	CategoryID         *int              // ID категории
	PriceMin           *float64          // Минимальная цена
	PriceMax           *float64          // Максимальная цена
	Condition          string            // Состояние (новое, б/у)
	City               string            // Город
	Country            string            // Страна
	StorefrontID       *int              // ID витрины
	Sort               string            // Поле сортировки
	SortDirection      string            // Направление сортировки (asc, desc)
	Location           *GeoLocation      // Координаты для геопоиска
	Distance           string            // Радиус поиска (5km, 10km, ...)
	Page               int               // Номер страницы
	Size               int               // Размер страницы
	Aggregations       []string          // Запрашиваемые агрегации
	Language           string            // Язык для поиска
	MinimumShouldMatch string            // Минимальное количество совпадений (70%, 50% и т.д.)
	Fuzziness          string            // Уровень нечеткости (AUTO, 1, 2, ...)
	AttributeFilters   map[string]string //  поле для фильтров атрибутов
}

// SearchParams параметры поиска для OpenSearch
type ServiceParams struct {
	Query              string            // Текстовый запрос
	CategoryID         string            // ID категории
	PriceMin           float64           // Минимальная цена
	PriceMax           float64           // Максимальная цена
	Condition          string            // Состояние (новое, б/у)
	City               string            // Город
	Country            string            // Страна
	StorefrontID       string            // ID витрины
	Sort               string            // Поле сортировки
	SortDirection      string            // Направление сортировки (asc, desc)
	Latitude           float64           // Широта для геопоиска
	Longitude          float64           // Долгота для геопоиска
	Distance           string            // Радиус поиска
	Page               int               // Номер страницы
	Size               int               // Размер страницы
	Aggregations       []string          // Запрашиваемые агрегации
	Language           string            // Язык для поиска
	MinimumShouldMatch string            // Минимальное количество совпадений (70%, 50% и т.д.)
	Fuzziness          string            // Уровень нечеткости (AUTO, 1, 2, ...)
	AttributeFilters   map[string]string //  поле для фильтров атрибутов
}

// GeoLocation координаты для геопоиска
type GeoLocation struct {
	Lat float64
	Lon float64
}

// Bucket для агрегаций
type Bucket struct {
	Key   string // Ключ бакета
	Count int    // Количество документов
}

// SearchResult результаты поиска
type SearchResult struct {
	Listings     []*models.MarketplaceListing // Найденные объявления
	Total        int                          // Общее количество найденных объявлений
	Took         int64                        // Время выполнения запроса в мс
	Aggregations map[string][]Bucket          // Фасеты для фильтров
	Suggestions  []string                     // Подсказки (для исправления опечаток)
}

// ServiceResult результаты для сервисного слоя
type ServiceResult struct {
	Items              []*models.MarketplaceListing `json:"items"`
	Total              int                          `json:"total"`
	Page               int                          `json:"page"`
	Size               int                          `json:"size"`
	TotalPages         int                          `json:"total_pages"`
	Took               int64                        `json:"took_ms"`
	Facets             map[string][]Bucket          `json:"facets,omitempty"`
	Suggestions        []string                     `json:"suggestions,omitempty"`
	SpellingSuggestion string                       `json:"spelling_suggestion,omitempty"` // Добавленное поле

}
