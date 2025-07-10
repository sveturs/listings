package behavior

import (
	"time"
)

// EventType представляет тип поведенческого события
type EventType string

const (
	EventTypeSearchPerformed     EventType = "search_performed"
	EventTypeResultClicked       EventType = "result_clicked"
	EventTypeItemViewed          EventType = "item_viewed"
	EventTypeItemPurchased       EventType = "item_purchased"
	EventTypeSearchFilterApplied EventType = "search_filter_applied"
	EventTypeSearchSortChanged   EventType = "search_sort_changed"
	EventTypeItemAddedToCart     EventType = "item_added_to_cart"
)

// ItemType представляет тип элемента
type ItemType string

const (
	ItemTypeMarketplace ItemType = "marketplace"
	ItemTypeStorefront  ItemType = "storefront"
)

// BehaviorEvent представляет поведенческое событие пользователя
type BehaviorEvent struct {
	ID          int64                  `json:"id,omitempty"`
	EventType   EventType              `json:"event_type" validate:"required"`
	UserID      *int                   `json:"user_id,omitempty"`
	SessionID   string                 `json:"session_id" validate:"required"`
	SearchQuery string                 `json:"search_query,omitempty"`
	ItemID      string                 `json:"item_id,omitempty"`
	ItemType    ItemType               `json:"item_type,omitempty"`
	Position    *int                   `json:"position,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"created_at,omitempty"`
}

// SearchMetrics представляет агрегированные метрики поиска
type SearchMetrics struct {
	ID               int64     `json:"id,omitempty"`
	SearchQuery      string    `json:"search_query"`
	TotalSearches    int       `json:"total_searches"`
	TotalClicks      int       `json:"total_clicks"`
	CTR              float64   `json:"ctr"`
	AvgClickPosition float64   `json:"avg_click_position"`
	Conversions      int       `json:"conversions"`
	ConversionRate   float64   `json:"conversion_rate"`
	PeriodStart      time.Time `json:"period_start"`
	PeriodEnd        time.Time `json:"period_end"`
	CreatedAt        time.Time `json:"created_at,omitempty"`
	UpdatedAt        time.Time `json:"updated_at,omitempty"`
}

// ItemMetrics представляет метрики для конкретного товара
type ItemMetrics struct {
	ItemID         string    `json:"item_id"`
	ItemType       ItemType  `json:"item_type"`
	Views          int       `json:"views"`
	Clicks         int       `json:"clicks"`
	Purchases      int       `json:"purchases"`
	CTR            float64   `json:"ctr"`
	ConversionRate float64   `json:"conversion_rate"`
	AvgPosition    float64   `json:"avg_position"`
	PeriodStart    time.Time `json:"period_start"`
	PeriodEnd      time.Time `json:"period_end"`
}

// TrackEventRequest представляет запрос на отслеживание события
type TrackEventRequest struct {
	EventType   EventType              `json:"event_type" validate:"required,oneof=search_performed result_clicked item_viewed item_purchased search_filter_applied search_sort_changed item_added_to_cart"`
	SessionID   string                 `json:"session_id,omitempty"` // Если не передан, будет сгенерирован
	SearchQuery string                 `json:"search_query,omitempty"`
	ItemID      string                 `json:"item_id,omitempty"`
	ItemType    ItemType               `json:"item_type,omitempty" validate:"omitempty,oneof=marketplace storefront"`
	Position    *int                   `json:"position,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// SearchMetricsQuery представляет параметры запроса метрик поиска
type SearchMetricsQuery struct {
	Query       string    `query:"query"`
	PeriodStart time.Time `query:"period_start"`
	PeriodEnd   time.Time `query:"period_end"`
	Limit       int       `query:"limit"`
	Offset      int       `query:"offset"`
	SortBy      string    `query:"sort_by" validate:"omitempty,oneof=ctr conversions total_searches"`
	OrderBy     string    `query:"order_by" validate:"omitempty,oneof=asc desc"`
}

// ItemMetricsQuery представляет параметры запроса метрик товаров
type ItemMetricsQuery struct {
	ItemType    ItemType  `query:"item_type" validate:"omitempty,oneof=marketplace storefront"`
	PeriodStart time.Time `query:"period_start"`
	PeriodEnd   time.Time `query:"period_end"`
	Limit       int       `query:"limit"`
	Offset      int       `query:"offset"`
	SortBy      string    `query:"sort_by" validate:"omitempty,oneof=views clicks purchases ctr conversion_rate"`
	OrderBy     string    `query:"order_by" validate:"omitempty,oneof=asc desc"`
}

// TrackEventBatch представляет пакет событий для отслеживания
type TrackEventBatch struct {
	Events    []TrackEventRequest `json:"events" validate:"required,dive"`
	BatchID   string              `json:"batch_id" validate:"required"`
	CreatedAt string              `json:"created_at" validate:"required"`
}

// AggregatedSearchMetrics представляет агрегированные метрики поиска
type AggregatedSearchMetrics struct {
	TotalSearches           int64                  `json:"total_searches"`
	UniqueSearches          int64                  `json:"unique_searches"`
	AverageSearchDurationMs float64                `json:"average_search_duration_ms"`
	SearchTrends            []SearchTrend          `json:"search_trends"`
	ClickMetrics            AggregatedClickMetrics `json:"click_metrics"`
}

// SearchTrend представляет тренд поиска по дням
type SearchTrend struct {
	Date          string  `json:"date"`
	SearchesCount int     `json:"searches_count"`
	ClicksCount   int     `json:"clicks_count"`
	CTR           float64 `json:"ctr"`
}

// AggregatedClickMetrics представляет агрегированные метрики кликов
type AggregatedClickMetrics struct {
	TotalClicks          int     `json:"total_clicks"`
	AverageClickPosition float64 `json:"average_click_position"`
	CTR                  float64 `json:"ctr"`
	ConversionRate       float64 `json:"conversion_rate"`
}

// TopSearchQuery представляет топ поисковый запрос с полной статистикой
type TopSearchQuery struct {
	Query       string  `json:"query"`
	Count       int     `json:"count"`
	CTR         float64 `json:"ctr"`
	AvgPosition float64 `json:"avg_position"`
	AvgResults  float64 `json:"avg_results"`
}
