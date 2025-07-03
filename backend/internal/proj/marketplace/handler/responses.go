// Package handler содержит структуры ответов для Swagger документации
package handler

import (
	"time"

	"backend/internal/domain/models"
	"backend/internal/domain/search"
)

// ErrorResponse базовая структура ответа с ошибкой
type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Error   string `json:"error" example:"marketplace.error"`
}

// MessageResponse ответ с сообщением
type MessageResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"marketplace.operationSuccess"`
}

// IDMessageResponse ответ с ID и сообщением
type IDMessageResponse struct {
	ID      int    `json:"id" example:"123"`
	Message string `json:"message,omitempty" example:"marketplace.created"`
}

// PaginationMeta метаданные пагинации
type PaginationMeta struct {
	Total      int  `json:"total" example:"100"`
	Page       int  `json:"page" example:"1"`
	Limit      int  `json:"limit" example:"20"`
	TotalPages int  `json:"total_pages,omitempty" example:"5"`
	HasMore    bool `json:"has_more,omitempty" example:"true"`
}

// ListingsResponse ответ со списком объявлений
type ListingsResponse struct {
	Success bool                        `json:"success" example:"true"`
	Data    []models.MarketplaceListing `json:"data"`
	Meta    PaginationMeta              `json:"meta"`
}

// ListingResponse ответ с одним объявлением
type ListingResponse struct {
	Success bool                       `json:"success" example:"true"`
	Data    *models.MarketplaceListing `json:"data"`
}

// PriceHistoryResponse ответ с историей цен
type PriceHistoryResponse struct {
	Success bool                       `json:"success" example:"true"`
	Data    []models.PriceHistoryEntry `json:"data"`
}

// TranslationUpdateRequest запрос обновления переводов
type TranslationUpdateRequest struct {
	Language     string            `json:"language" example:"en"`
	Translations map[string]string `json:"translations"`
	IsVerified   bool              `json:"is_verified" example:"false"`
	Provider     string            `json:"provider" example:"google"`
}

// TranslationsResponse ответ с переводами
type TranslationsResponse struct {
	Success bool                 `json:"success" example:"true"`
	Data    []models.Translation `json:"data"`
}

// TranslateTextRequest запрос на перевод текста
type TranslateTextRequest struct {
	Text       string `json:"text" example:"Hello world"`
	SourceLang string `json:"source_lang" example:"en"`
	TargetLang string `json:"target_lang" example:"ru"`
	Provider   string `json:"provider" example:"google"`
}

// TranslateTextResponse ответ перевода текста
type TranslateTextResponse struct {
	Success bool               `json:"success" example:"true"`
	Data    TranslatedTextData `json:"data"`
}

// TranslatedTextData данные переведенного текста
type TranslatedTextData struct {
	TranslatedText string `json:"translated_text" example:"Привет мир"`
	SourceLang     string `json:"source_lang" example:"en"`
	TargetLang     string `json:"target_lang" example:"ru"`
	Provider       string `json:"provider" example:"google"`
}

// DetectLanguageRequest запрос определения языка
type DetectLanguageRequest struct {
	Text string `json:"text" example:"Hello world"`
}

// DetectLanguageResponse ответ определения языка
type DetectLanguageResponse struct {
	Success bool                 `json:"success" example:"true"`
	Data    DetectedLanguageData `json:"data"`
}

// DetectedLanguageData данные определенного языка
type DetectedLanguageData struct {
	Language   string  `json:"language" example:"en"`
	Confidence float64 `json:"confidence" example:"0.99"`
}

// TranslationLimitsResponse ответ с лимитами переводов
type TranslationLimitsResponse struct {
	Success bool                  `json:"success" example:"true"`
	Data    TranslationLimitsData `json:"data"`
}

// TranslationLimitsData данные лимитов переводов
type TranslationLimitsData struct {
	DailyLimit int    `json:"daily_limit" example:"10000"`
	UsedToday  int    `json:"used_today" example:"3450"`
	Remaining  int    `json:"remaining" example:"6550"`
	Provider   string `json:"provider" example:"google"`
}

// SetProviderRequest запрос установки провайдера
type SetProviderRequest struct {
	Provider string `json:"provider" example:"google"`
}

// SetProviderResponse ответ установки провайдера
type SetProviderResponse struct {
	Success bool         `json:"success" example:"true"`
	Message string       `json:"message" example:"marketplace.providerSet"`
	Data    ProviderData `json:"data"`
}

// ProviderData данные провайдера
type ProviderData struct {
	Provider string `json:"provider" example:"google"`
}

// BatchTranslateRequest запрос пакетного перевода
type BatchTranslateRequest struct {
	ListingIDs []int  `json:"listing_ids" example:"[1,2,3]"`
	TargetLang string `json:"target_lang" example:"en"`
	Provider   string `json:"provider" example:"google"`
}

// BatchTranslateResponse ответ пакетного перевода
type BatchTranslateResponse struct {
	Success bool               `json:"success" example:"true"`
	Message string             `json:"message" example:"marketplace.translationProcessStarted"`
	Data    BatchTranslateData `json:"data"`
}

// BatchTranslateData данные пакетного перевода
type BatchTranslateData struct {
	ListingCount int    `json:"listing_count" example:"3"`
	TargetLang   string `json:"target_lang" example:"en"`
	Provider     string `json:"provider" example:"google"`
}

// AttributeGroupsResponse ответ со списком групп атрибутов
type AttributeGroupsResponse struct {
	Groups []*models.AttributeGroup `json:"groups"`
}

// AttributeGroupResponse ответ с группой атрибутов
type AttributeGroupResponse struct {
	Group *models.AttributeGroup `json:"group"`
}

// AttributeGroupWithItemsResponse ответ с группой и её элементами
type AttributeGroupWithItemsResponse struct {
	Success bool                        `json:"success" example:"true"`
	Data    AttributeGroupWithItemsData `json:"data"`
}

// AttributeGroupWithItemsData данные группы с элементами
type AttributeGroupWithItemsData struct {
	Group *models.AttributeGroup       `json:"group"`
	Items []*models.AttributeGroupItem `json:"items"`
}

// SearchResponse ответ поиска
type SearchResponse struct {
	Data []models.MarketplaceListing `json:"data"`
	Meta SearchMetadata              `json:"meta"`
}

// SearchMetadata метаданные поиска
type SearchMetadata struct {
	Total              int                        `json:"total" example:"100"`
	Page               int                        `json:"page" example:"1"`
	Size               int                        `json:"size" example:"20"`
	TotalPages         int                        `json:"total_pages" example:"5"`
	HasMore            bool                       `json:"has_more" example:"true"`
	Facets             map[string][]search.Bucket `json:"facets,omitempty"`
	Suggestions        []string                   `json:"suggestions,omitempty"`
	SpellingSuggestion string                     `json:"spelling_suggestion,omitempty"`
	TookMs             int64                      `json:"took_ms,omitempty" example:"25"`
}

// SuggestionItem элемент подсказки
type SuggestionItem struct {
	Text     string           `json:"text" example:"iPhone 13"`
	Type     string           `json:"type" example:"product"`
	Category *CategorySummary `json:"category,omitempty"`
}

// CategorySummary краткая информация о категории
type CategorySummary struct {
	ID   int    `json:"id" example:"10"`
	Name string `json:"name" example:"Smartphones"`
	Slug string `json:"slug" example:"smartphones"`
}

// SuggestionsResponse ответ с подсказками
type SuggestionsResponse struct {
	Success bool             `json:"success" example:"true"`
	Data    []SuggestionItem `json:"data"`
}

// MapBounds границы карты
type MapBounds struct {
	NE Coordinates `json:"ne"`
	SW Coordinates `json:"sw"`
}

// Coordinates координаты
type Coordinates struct {
	Lat float64 `json:"lat" example:"55.7558"`
	Lng float64 `json:"lng" example:"37.6173"`
}

// MapListingsResponse ответ с объявлениями на карте
type MapListingsResponse struct {
	Listings []models.MarketplaceListing `json:"listings"`
	Bounds   MapBounds                   `json:"bounds"`
	Zoom     int                         `json:"zoom" example:"10"`
	Count    int                         `json:"count" example:"50"`
}

// MapClustersResponse ответ с кластерами на карте
type MapClustersResponse struct {
	Success bool            `json:"success" example:"true"`
	Data    MapClustersData `json:"data"`
}

// AttributeRangesResponse ответ с диапазонами атрибутов
type AttributeRangesResponse struct {
	Success bool                      `json:"success" example:"true"`
	Data    map[string]AttributeRange `json:"data"`
}

// AttributeRange диапазон значений атрибута
type AttributeRange struct {
	Min float64 `json:"min" example:"0"`
	Max float64 `json:"max" example:"1000000"`
}

// ChatMessage сообщение чата
type ChatMessage struct {
	ID           int       `json:"id"`
	ChatID       int       `json:"chat_id"`
	SenderID     int       `json:"sender_id"`
	RecipientID  int       `json:"recipient_id"`
	Message      string    `json:"message"`
	IsRead       bool      `json:"is_read"`
	CreatedAt    time.Time `json:"created_at"`
	SenderName   string    `json:"sender_name,omitempty"`
	SenderAvatar string    `json:"sender_avatar,omitempty"`
}

// ChatResponse ответ с чатом
type ChatResponse struct {
	Success bool                   `json:"success" example:"true"`
	Data    models.MarketplaceChat `json:"data"`
}

// ChatsResponse ответ со списком чатов
type ChatsResponse struct {
	Success bool                     `json:"success" example:"true"`
	Data    []models.MarketplaceChat `json:"data"`
}

// MessagesResponse ответ со списком сообщений
type MessagesResponse struct {
	Success bool          `json:"success" example:"true"`
	Data    []ChatMessage `json:"data"`
}

// UnreadCountResponse ответ с количеством непрочитанных
type UnreadCountResponse struct {
	Success bool `json:"success" example:"true"`
	Count   int  `json:"count" example:"5"`
}

// ImageUploadResponse ответ загрузки изображения
type ImageUploadResponse struct {
	Success  bool   `json:"success" example:"true"`
	ImageURL string `json:"image_url" example:"https://svetu.rs/listings/1234_image_0.jpg"`
	Message  string `json:"message,omitempty" example:"marketplace.imageUploaded"`
}

// ReindexResponse ответ переиндексации
type ReindexResponse struct {
	Message string `json:"message" example:"Reindexing started"`
}

// AdminStatsResponse ответ со статистикой администратора
type AdminStatsResponse struct {
	Success bool           `json:"success" example:"true"`
	Data    AdminStatsData `json:"data"`
}

// AdminStatsData данные статистики администратора
type AdminStatsData struct {
	TotalListings   int `json:"total_listings" example:"1000"`
	ActiveListings  int `json:"active_listings" example:"850"`
	TotalUsers      int `json:"total_users" example:"500"`
	TotalCategories int `json:"total_categories" example:"50"`
	TotalAttributes int `json:"total_attributes" example:"200"`
	ListingsToday   int `json:"listings_today" example:"25"`
}

// CustomComponentResponse ответ с пользовательским компонентом
type CustomComponentResponse struct {
	Success bool                     `json:"success" example:"true"`
	Data    models.CustomUIComponent `json:"data"`
}

// CustomComponentsResponse ответ со списком пользовательских компонентов
type CustomComponentsResponse struct {
	Success bool                       `json:"success" example:"true"`
	Data    []models.CustomUIComponent `json:"data"`
}

// MapBoundsData represents data for map bounds response
type MapBoundsData struct {
	Listings []models.MapMarker `json:"listings"`
	Bounds   MapBounds          `json:"bounds"`
	Zoom     int                `json:"zoom" example:"10"`
	Count    int                `json:"count" example:"25"`
}

// MapBoundsResponse represents the response for map bounds
type MapBoundsResponse struct {
	Success bool          `json:"success" example:"true"`
	Data    MapBoundsData `json:"data"`
}

// MapClustersData represents data for map clusters response
type MapClustersData struct {
	Type  string      `json:"type" example:"clusters"`
	Data  interface{} `json:"data"`
	Zoom  int         `json:"zoom" example:"10"`
	Count int         `json:"count" example:"50"`
}

// ImagesUploadResponse represents the response for image upload
type ImagesUploadResponse struct {
	Success bool                      `json:"success" example:"true"`
	Message string                    `json:"message" example:"marketplace.imagesUploaded"`
	Images  []models.MarketplaceImage `json:"images"`
	Count   int                       `json:"count" example:"3"`
}

// ModerationData represents image moderation data
type ModerationData struct {
	Labels           []string `json:"labels" example:"image,photo"`
	ProhibitedLabels []string `json:"prohibited_labels"`
	HasProhibited    bool     `json:"has_prohibited" example:"false"`
}

// ImageModerationResponse represents the response for image moderation
type ImageModerationResponse struct {
	Success bool           `json:"success" example:"true"`
	Data    ModerationData `json:"data"`
}

// EnhancePreviewData represents enhancement preview data
type EnhancePreviewData struct {
	PreviewURL string `json:"preview_url" example:"https://example.com/preview/123/quality"`
}

// EnhancePreviewResponse represents the response for enhancement preview
type EnhancePreviewResponse struct {
	Success bool               `json:"success" example:"true"`
	Data    EnhancePreviewData `json:"data"`
}

// EnhanceImagesData represents enhancement job data
type EnhanceImagesData struct {
	Message string `json:"message" example:"marketplace.imageEnhancementStarted"`
	JobID   string `json:"job_id" example:"enhance_123_1234567890"`
}

// EnhanceImagesResponse represents the response for image enhancement
type EnhanceImagesResponse struct {
	Success bool              `json:"success" example:"true"`
	Data    EnhanceImagesData `json:"data"`
}

// ChatMessagesResponse ответ со списком сообщений чата
type ChatMessagesResponse struct {
	Messages []models.MarketplaceMessage `json:"messages"`
	Total    int                         `json:"total"`
	Page     int                         `json:"page"`
	Limit    int                         `json:"limit"`
}

// UnreadCountData данные о количестве непрочитанных сообщений
type UnreadCountData struct {
	Count int `json:"count" example:"5"`
}

// FavoritesResponse ответ со списком избранных объявлений
type FavoritesResponse struct {
	Success bool                        `json:"success" example:"true"`
	Data    []models.MarketplaceListing `json:"data"`
}

// FavoriteStatusResponse ответ о статусе избранного
type FavoriteStatusResponse struct {
	Success bool               `json:"success" example:"true"`
	Data    FavoriteStatusData `json:"data"`
}

// FavoriteStatusData данные о статусе избранного
type FavoriteStatusData struct {
	IsInFavorites bool `json:"is_in_favorites" example:"true"`
}

// FavoritesCountResponse ответ с количеством избранных
type FavoritesCountResponse struct {
	Success bool               `json:"success" example:"true"`
	Data    FavoritesCountData `json:"data"`
}

// FavoritesCountData данные о количестве избранных
type FavoritesCountData struct {
	Count int `json:"count" example:"5"`
}

// ReindexStartedResponse ответ о начале переиндексации
type ReindexStartedResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"marketplace.reindexStarted"`
}

// AttributeCreateResponse ответ создания атрибута
type AttributeCreateResponse struct {
	ID      int    `json:"id" example:"123"`
	Message string `json:"message" example:"marketplace.attributeCreated"`
}

// BulkUpdateResult результат массового обновления
type BulkUpdateResult struct {
	SuccessCount int      `json:"success_count" example:"5"`
	TotalCount   int      `json:"total_count" example:"10"`
	Errors       []string `json:"errors,omitempty"`
}

// ImportAttributesResult результат импорта атрибутов
type ImportAttributesResult struct {
	SuccessCount int      `json:"success_count" example:"5"`
	TotalCount   int      `json:"total_count" example:"10"`
	Errors       []string `json:"errors,omitempty"`
}

// PartialOperationResponse ответ при частичном выполнении операции
type PartialOperationResponse struct {
	Success bool        `json:"success" example:"false"`
	Error   string      `json:"error" example:"marketplace.partialOperationCompleted"`
	Data    interface{} `json:"data"`
}

// TranslationResult результат автоматического перевода атрибута
type TranslationResult struct {
	AttributeID  int                    `json:"attribute_id" example:"123"`
	Translations map[string]interface{} `json:"translations"`
	Errors       []string               `json:"errors,omitempty"`
}
