package opensearch

import (
    "backend/internal/domain/models"
)

// SearchParams параметры поиска
type SearchParams struct {
    Query         string            // Текст поиска
    CategoryID    *int              // ID категории
    PriceMin      *float64          // Минимальная цена
    PriceMax      *float64          // Максимальная цена
    Condition     string            // Состояние (новое, б/у)
    City          string            // Город
    Country       string            // Страна
    StorefrontID  *int              // ID витрины
    Sort          string            // Поле сортировки
    SortDirection string            // Направление сортировки (asc, desc)
    Location      *GeoLocation      // Координаты для геопоиска
    Distance      string            // Радиус поиска (5km, 10km, ...)
    Page          int               // Номер страницы
    Size          int               // Размер страницы
    Aggregations  []string          // Запрашиваемые агрегации
    Language      string            // Язык для поиска
}

// GeoLocation координаты для геопоиска
type GeoLocation struct {
    Lat float64
    Lon float64
}

// SearchResult результаты поиска
type SearchResult struct {
    Listings     []*models.MarketplaceListing // Найденные объявления
    Total        int                         // Общее количество найденных объявлений
    Took         int64                       // Время выполнения запроса в мс
    Aggregations map[string][]Bucket        // Фасеты для фильтров
    Suggestions  []string                    // Подсказки (для исправления опечаток)
}

// Bucket для агрегаций
type Bucket struct {
    Key   string // Ключ бакета
    Count int    // Количество документов
}