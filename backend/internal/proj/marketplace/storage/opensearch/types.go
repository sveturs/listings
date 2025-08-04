package opensearch

import (
	"backend/internal/domain/models"
)

// GeoLocation координаты для геопоиска
type GeoLocation struct {
	Lat float64
	Lon float64
}

// DBTranslation представляет перевод из таблицы translations
type DBTranslation struct {
	Language       string `json:"language"`
	FieldName      string `json:"field_name"`
	TranslatedText string `json:"translated_text"`
}

// SearchResult результаты поиска
type SearchResult struct {
	Listings     []*models.MarketplaceListing // Найденные объявления
	Total        int                          // Общее количество найденных объявлений
	Took         int64                        // Время выполнения запроса в мс
	Aggregations map[string][]Bucket          // Фасеты для фильтров
	Suggestions  []string                     // Подсказки (для исправления опечаток)
}

// Bucket для агрегаций
type Bucket struct {
	Key   string // Ключ бакета
	Count int    // Количество документов
}
