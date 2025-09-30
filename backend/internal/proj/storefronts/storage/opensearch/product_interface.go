package opensearch

import (
	"context"
	"time"

	"backend/internal/domain/models"
)

// ProductSearchRepository интерфейс для работы с поиском товаров витрин
type ProductSearchRepository interface {
	// PrepareIndex подготавливает индекс (создает если не существует)
	PrepareIndex(ctx context.Context) error

	// IndexProduct индексирует один товар
	IndexProduct(ctx context.Context, product *models.StorefrontProduct) error

	// BulkIndexProducts индексирует несколько товаров
	BulkIndexProducts(ctx context.Context, products []*models.StorefrontProduct) error

	// UpdateProduct обновляет товар в индексе
	UpdateProduct(ctx context.Context, product *models.StorefrontProduct) error

	// DeleteProduct удаляет товар из индекса
	DeleteProduct(ctx context.Context, productID int) error

	// UpdateProductStock частично обновляет поля склада товара
	UpdateProductStock(ctx context.Context, productID int, stockData map[string]interface{}) error

	// SearchProducts выполняет поиск товаров
	SearchProducts(ctx context.Context, params *ProductSearchParams) (*ProductSearchResult, error)

	// ReindexAll переиндексирует все товары витрин
	ReindexAll(ctx context.Context) error

	// SearchSimilarProducts выполняет поиск похожих товаров
	SearchSimilarProducts(ctx context.Context, productID int, limit int) ([]*models.MarketplaceListing, error)
}

// ProductSearchParams параметры поиска товаров витрин
type ProductSearchParams struct {
	Query           string                 // Текстовый поиск
	StorefrontID    int                    // Фильтр по витрине
	CategoryID      int                    // Фильтр по категории
	CategoryIDs     []int                  // Фильтр по множественным категориям
	CategoryPath    string                 // Фильтр по пути категории
	Brand           string                 // Фильтр по бренду
	PriceMin        float64                // Минимальная цена
	PriceMax        float64                // Максимальная цена
	InStock         *bool                  // Только в наличии
	IsActive        *bool                  // Фильтр по активности (nil = не фильтровать, для админки)
	Attributes      map[string]interface{} // Фильтры по атрибутам
	City            string                 // Город витрины
	Latitude        float64                // Широта для геопоиска
	Longitude       float64                // Долгота для геопоиска
	RadiusKm        int                    // Радиус поиска в километрах
	SortBy          string                 // Поле сортировки
	SortOrder       string                 // Порядок сортировки (asc/desc)
	IncludeVariants bool                   // Включить варианты в результаты
	OnlyVerified    bool                   // Только от верифицированных витрин
	MinQualityScore float64                // Минимальный показатель качества
	Aggregations    []string               // Запрашиваемые агрегации
	Limit           int                    // Количество результатов
	Offset          int                    // Смещение
	Language        string                 // Язык для анализатора
}

// ProductSearchResult результат поиска товаров
type ProductSearchResult struct {
	Products     []*ProductSearchItem
	Total        int
	Aggregations map[string]interface{}
	TookMs       int64
}

// ProductSearchItem элемент результата поиска
type ProductSearchItem struct {
	ID                string              // Уникальный ID документа (sp_123)
	ProductID         int                 // ID товара
	StorefrontID      int                 // ID витрины
	Name              string              // Название товара
	Description       string              // Описание
	Price             float64             // Цена
	PriceMin          float64             // Минимальная цена (с учетом вариантов)
	PriceMax          float64             // Максимальная цена (с учетом вариантов)
	Currency          string              // Валюта
	SKU               string              // Артикул
	Brand             string              // Бренд
	InStock           bool                // В наличии
	AvailableQuantity int                 // Доступное количество
	Images            []ProductImage      // Изображения
	Storefront        StorefrontInfo      // Информация о витрине
	Category          CategoryInfo        // Информация о категории
	Attributes        []ProductAttribute  // Атрибуты
	Variants          []ProductVariant    // Варианты
	Score             float64             // Релевантность
	Distance          *float64            // Расстояние в км (если есть)
	Highlights        map[string][]string // Подсвеченные фрагменты
	PopularityScore   float64             // Показатель популярности
	QualityScore      float64             // Показатель качества
	ViewsCount        int                 // Количество просмотров
	SoldCount         int                 // Количество продаж
	CreatedAt         *time.Time          // Дата создания товара
	UpdatedAt         *time.Time          // Дата обновления товара
}

// ProductImage изображение товара
type ProductImage struct {
	ID       int    `json:"id"`
	URL      string `json:"url"`
	AltText  string `json:"alt_text"`
	IsMain   bool   `json:"is_main"`
	Position int    `json:"position"`
}

// StorefrontInfo информация о витрине
type StorefrontInfo struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Slug       string  `json:"slug"`
	City       string  `json:"city"`
	Country    string  `json:"country"`
	Rating     float64 `json:"rating"`
	IsVerified bool    `json:"is_verified"`
}

// CategoryInfo информация о категории
type CategoryInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	Path string `json:"path,omitempty"`
}

// ProductAttribute атрибут товара
type ProductAttribute struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Type        string      `json:"type"`
	Value       interface{} `json:"value"`
	DisplayName string      `json:"display_name"`
}

// ProductVariant вариант товара
type ProductVariant struct {
	ID                int                    `json:"id"`
	Name              string                 `json:"name"`
	SKU               string                 `json:"sku"`
	Price             float64                `json:"price"`
	StockQuantity     int                    `json:"stock_quantity"`
	AvailableQuantity int                    `json:"available_quantity"`
	StockStatus       string                 `json:"stock_status"`
	IsActive          bool                   `json:"is_active"`
	Attributes        map[string]interface{} `json:"attributes"`
}
