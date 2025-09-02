package models

// UnifiedSuggestion представляет унифицированную подсказку для автодополнения
type UnifiedSuggestion struct {
	Type       string                 `json:"type"`                  // "query", "product", "category"
	Value      string                 `json:"value"`                 // Основное значение для поиска
	Label      string                 `json:"label"`                 // Отображаемый текст
	Count      *int                   `json:"count,omitempty"`       // Количество результатов
	CategoryID *int                   `json:"category_id,omitempty"` // ID категории для type="category"
	ProductID  *int                   `json:"product_id,omitempty"`  // ID товара для type="product"
	Icon       *string                `json:"icon,omitempty"`        // Иконка
	Metadata   *UnifiedSuggestionMeta `json:"metadata,omitempty"`    // Дополнительные данные
}

// UnifiedSuggestionMeta дополнительные метаданные для подсказок
type UnifiedSuggestionMeta struct {
	Price          *float64 `json:"price,omitempty"`           // Цена товара
	Image          *string  `json:"image,omitempty"`           // URL изображения
	Category       *string  `json:"category,omitempty"`        // Название категории
	SourceType     *string  `json:"source_type,omitempty"`     // "marketplace" или "storefront"
	StorefrontID   *int     `json:"storefront_id,omitempty"`   // ID витрины
	Storefront     *string  `json:"storefront,omitempty"`      // Название витрины
	StorefrontSlug *string  `json:"storefront_slug,omitempty"` // Slug витрины
	ParentID       *int     `json:"parent_id,omitempty"`       // ID родительской категории
	LastSearched   *string  `json:"last_searched,omitempty"`   // Дата последнего поиска
}

// SuggestionRequestParams параметры запроса подсказок
type SuggestionRequestParams struct {
	Query    string   `json:"query"`              // Поисковый запрос
	Types    []string `json:"types"`              // Типы подсказок: queries, categories, products
	Limit    int      `json:"limit"`              // Максимальное количество
	Category *string  `json:"category,omitempty"` // Фильтр по категории
	Language string   `json:"language"`           // Язык ответа
}
