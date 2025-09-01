// backend/internal/domain/models/models.go
package models

import (
	"time"
)

const (
	NotificationTypeNewMessage     = "new_message"
	NotificationTypeNewReview      = "new_review"
	NotificationTypeReviewVote     = "review_vote"
	NotificationTypeReviewResponse = "review_response"
	NotificationTypeListingStatus  = "listing_status"
	NotificationTypeFavoritePrice  = "favorite_price"
	NotificationTypeDeliveryStatus = "delivery_status"
)

// PaginatedResponse представляет пагинированный ответ
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	Total      int         `json:"total"`
	TotalPages int         `json:"total_pages"`
}

type User struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	GoogleID   string    `json:"google_id"`
	PictureURL string    `json:"picture_url"`
	Phone      *string   `json:"phone,omitempty"`
	Password   *string   `json:"-"` // скрываем пароль в JSON ответах, может быть NULL для OAuth пользователей
	Provider   string    `json:"provider"`
	CreatedAt  time.Time `json:"created_at"`
}
type (
	TranslationMap     map[string]map[string]string
	MarketplaceListing struct {
		ID                 int                     `json:"id"`
		UserID             int                     `json:"user_id"`
		CategoryID         int                     `json:"category_id"`
		Title              string                  `json:"title"`
		Description        string                  `json:"description"`
		Price              float64                 `json:"price"`
		Condition          string                  `json:"condition"`
		Status             string                  `json:"status"`
		Location           string                  `json:"location"`
		Latitude           *float64                `json:"latitude,omitempty"`
		Longitude          *float64                `json:"longitude,omitempty"`
		City               string                  `json:"address_city"`
		Country            string                  `json:"address_country"`
		ViewsCount         int                     `json:"views_count"`
		CreatedAt          time.Time               `json:"created_at"`
		UpdatedAt          time.Time               `json:"updated_at"`
		Images             []MarketplaceImage      `json:"images,omitempty"`
		User               *User                   `json:"user,omitempty"`
		Category           *MarketplaceCategory    `json:"category,omitempty"`
		HelpfulVotes       int                     `json:"helpful_votes"`
		NotHelpfulVotes    int                     `json:"not_helpful_votes"`
		IsFavorite         bool                    `json:"is_favorite"`
		ShowOnMap          bool                    `json:"show_on_map"`
		LocationPrivacy    string                  `json:"location_privacy,omitempty"`
		OriginalLanguage   string                  `json:"original_language,omitempty"`
		RawTranslations    interface{}             `json:"-"` // Для хранения "сырых" данных
		Translations       TranslationMap          `json:"translations,omitempty"`
		CategoryPathNames  []string                `json:"category_path_names,omitempty"`
		CategoryPathIds    []int                   `json:"category_path_ids,omitempty"`
		CategoryPathSlugs  []string                `json:"category_path_slugs,omitempty"`
		CategoryPath       []string                `json:"category_path,omitempty"`
		StorefrontID       *int                    `json:"storefront_id,omitempty"` // связь с витриной
		Storefront         *Storefront             `json:"storefront,omitempty"`    // данные витрины
		ExternalID         string                  `json:"external_id,omitempty"`
		Attributes         []ListingAttributeValue `json:"attributes,omitempty"`
		OldPrice           *float64                `json:"old_price,omitempty"`
		HasDiscount        bool                    `json:"has_discount"`
		DiscountPercentage *int                    `json:"discount_percentage,omitempty"`
		Metadata           map[string]interface{}  `json:"metadata,omitempty"` // Для хранения дополнительной информации, включая данные о скидке

		AverageRating float64 `json:"average_rating,omitempty"`
		ReviewCount   int     `json:"review_count,omitempty"`

		// Поля для товаров витрин
		StockQuantity *int    `json:"stock_quantity,omitempty"`
		StockStatus   *string `json:"stock_status,omitempty"`
	}
)

type MarketplaceCategory struct {
	ID                int               `json:"id" db:"id"`
	Name              string            `json:"name" db:"name"`
	Slug              string            `json:"slug" db:"slug"`
	ParentID          *int              `json:"parent_id,omitempty" db:"parent_id"`
	Icon              *string           `json:"icon,omitempty" db:"icon"`
	Description       string            `json:"description,omitempty" db:"description"`
	IsActive          bool              `json:"is_active" db:"is_active"`
	CreatedAt         time.Time         `json:"created_at" db:"created_at"`
	Translations      map[string]string `json:"translations,omitempty"`
	ListingCount      int               `json:"listing_count" db:"listing_count"`
	HasCustomUI       bool              `json:"has_custom_ui,omitempty" db:"has_custom_ui"`
	CustomUIComponent string            `json:"custom_ui_component,omitempty" db:"custom_ui_component"`
	SortOrder         int               `json:"sort_order" db:"sort_order"`
	Level             int               `json:"level" db:"level"`
	Count             int               `json:"count" db:"count"`
	ExternalID        string            `json:"external_id,omitempty" db:"external_id"`
	SEOTitle          string            `json:"seo_title,omitempty" db:"seo_title"`
	SEODescription    string            `json:"seo_description,omitempty" db:"seo_description"`
	SEOKeywords       string            `json:"seo_keywords,omitempty" db:"seo_keywords"`
}

type MarketplaceImage struct {
	ID            int       `json:"id"`
	ListingID     int       `json:"listing_id"`
	FilePath      string    `json:"file_path"`                // Путь к файлу в хранилище
	FileName      string    `json:"file_name"`                // Оригинальное имя файла
	FileSize      int       `json:"file_size"`                // Размер файла в байтах
	ContentType   string    `json:"content_type"`             // MIME-тип файла
	IsMain        bool      `json:"is_main"`                  // Является ли изображение основным
	StorageType   string    `json:"storage_type"`             // Тип хранилища: "local" или "minio"
	StorageBucket string    `json:"storage_bucket,omitempty"` // Имя бакета для MinIO
	PublicURL     string    `json:"public_url,omitempty"`     // Публичный URL для доступа к файлу
	ImageURL      string    `json:"image_url,omitempty"`      // URL изображения для API
	ThumbnailURL  string    `json:"thumbnail_url,omitempty"`  // URL миниатюры для API
	DisplayOrder  int       `json:"display_order"`            // Порядок отображения
	CreatedAt     time.Time `json:"created_at"`
}

// Реализация интерфейса ImageInterface для MarketplaceImage
func (m *MarketplaceImage) GetID() int {
	return m.ID
}

func (m *MarketplaceImage) GetEntityType() string {
	return "marketplace"
}

func (m *MarketplaceImage) GetEntityID() int {
	return m.ListingID
}

func (m *MarketplaceImage) GetFilePath() string {
	return m.FilePath
}

func (m *MarketplaceImage) GetFileName() string {
	return m.FileName
}

func (m *MarketplaceImage) GetFileSize() int {
	return m.FileSize
}

func (m *MarketplaceImage) GetContentType() string {
	return m.ContentType
}

func (m *MarketplaceImage) GetIsMain() bool {
	return m.IsMain
}

func (m *MarketplaceImage) GetStorageType() string {
	return m.StorageType
}

func (m *MarketplaceImage) GetStorageBucket() string {
	return m.StorageBucket
}

func (m *MarketplaceImage) GetPublicURL() string {
	return m.PublicURL
}

func (m *MarketplaceImage) GetImageURL() string {
	return m.ImageURL
}

func (m *MarketplaceImage) GetThumbnailURL() string {
	return m.ThumbnailURL
}

func (m *MarketplaceImage) GetDisplayOrder() int {
	return m.DisplayOrder
}

func (m *MarketplaceImage) GetCreatedAt() time.Time {
	return m.CreatedAt
}

func (m *MarketplaceImage) SetID(id int) {
	m.ID = id
}

func (m *MarketplaceImage) SetEntityID(entityID int) {
	m.ListingID = entityID
}

func (m *MarketplaceImage) SetFilePath(filePath string) {
	m.FilePath = filePath
}

func (m *MarketplaceImage) SetFileName(fileName string) {
	m.FileName = fileName
}

func (m *MarketplaceImage) SetFileSize(fileSize int) {
	m.FileSize = fileSize
}

func (m *MarketplaceImage) SetContentType(contentType string) {
	m.ContentType = contentType
}

func (m *MarketplaceImage) SetIsMain(isMain bool) {
	m.IsMain = isMain
}

func (m *MarketplaceImage) SetStorageType(storageType string) {
	m.StorageType = storageType
}

func (m *MarketplaceImage) SetStorageBucket(bucket string) {
	m.StorageBucket = bucket
}

func (m *MarketplaceImage) SetPublicURL(url string) {
	m.PublicURL = url
}

func (m *MarketplaceImage) SetImageURL(url string) {
	m.ImageURL = url
}

func (m *MarketplaceImage) SetThumbnailURL(url string) {
	m.ThumbnailURL = url
}

func (m *MarketplaceImage) SetDisplayOrder(order int) {
	m.DisplayOrder = order
}

func (m *MarketplaceImage) SetCreatedAt(createdAt time.Time) {
	m.CreatedAt = createdAt
}

func (m *MarketplaceImage) IsMainImage() bool {
	return m.IsMain
}

func (m *MarketplaceImage) SetMainImage(isMain bool) {
	m.IsMain = isMain
}

type CategoryTreeNode struct {
	ID                int                `json:"id"`
	Name              string             `json:"name"`
	Slug              string             `json:"slug"`
	Icon              string             `json:"icon,omitempty"`
	ParentID          *int               `json:"parent_id,omitempty"`
	CreatedAt         string             `json:"created_at"`
	Level             int                `json:"level"`
	Path              string             `json:"path"`
	ListingCount      int                `json:"listing_count"`
	Children          []CategoryTreeNode `json:"children,omitempty"`
	ChildrenCount     int                `json:"children_count"` // Новое поле
	Translations      map[string]string  `json:"translations,omitempty"`
	HasCustomUI       bool               `json:"has_custom_ui,omitempty"`
	CustomUIComponent string             `json:"custom_ui_component,omitempty"`
}

// CategoryDetectionStats статистика определения категорий
type CategoryDetectionStats struct {
	ID                      int32                  `db:"id" json:"id"`
	UserID                  *int32                 `db:"user_id" json:"user_id,omitempty"`
	SessionID               string                 `db:"session_id" json:"session_id,omitempty"`
	Method                  string                 `db:"method" json:"method"`
	AIKeywords              []string               `db:"ai_keywords" json:"ai_keywords,omitempty"`
	AIAttributes            map[string]interface{} `db:"ai_attributes" json:"ai_attributes,omitempty"`
	AIDomain                string                 `db:"ai_domain" json:"ai_domain,omitempty"`
	AIProductType           string                 `db:"ai_product_type" json:"ai_product_type,omitempty"`
	AISuggestedCategoryID   *int32                 `db:"ai_suggested_category_id" json:"ai_suggested_category_id,omitempty"`
	FinalCategoryID         *int32                 `db:"final_category_id" json:"final_category_id,omitempty"`
	AlternativeCategories   []byte                 `db:"alternative_categories" json:"-"`
	ConfidenceScore         *float64               `db:"confidence_score" json:"confidence_score,omitempty"`
	SimilarityScore         *float64               `db:"similarity_score" json:"similarity_score,omitempty"`
	KeywordScore            *float64               `db:"keyword_score" json:"keyword_score,omitempty"`
	SimilarListingsFound    *int32                 `db:"similar_listings_found" json:"similar_listings_found,omitempty"`
	TopSimilarListingID     *int32                 `db:"top_similar_listing_id" json:"top_similar_listing_id,omitempty"`
	TopSimilarityScore      *float64               `db:"top_similarity_score" json:"top_similarity_score,omitempty"`
	MatchedKeywords         []string               `db:"matched_keywords" json:"matched_keywords,omitempty"`
	MatchedNegativeKeywords []string               `db:"matched_negative_keywords" json:"matched_negative_keywords,omitempty"`
	ProcessingTimeMs        *int64                 `db:"processing_time_ms" json:"processing_time_ms,omitempty"`
	UserConfirmed           *bool                  `db:"user_confirmed" json:"user_confirmed,omitempty"`
	UserSelectedCategoryID  *int32                 `db:"user_selected_category_id" json:"user_selected_category_id,omitempty"`
	CreatedAt               time.Time              `db:"created_at" json:"created_at"`
}

// CategoryKeyword модель ключевого слова категории
type CategoryKeyword struct {
	ID          int32     `json:"id"`
	CategoryID  int32     `json:"category_id"`
	Keyword     string    `json:"keyword"`
	Language    string    `json:"language"`
	Weight      float64   `json:"weight"`
	KeywordType string    `json:"keyword_type"`
	IsNegative  bool      `json:"is_negative"`
	Source      string    `json:"source"`
	UsageCount  int32     `json:"usage_count"`
	SuccessRate float64   `json:"success_rate"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CategoryKeywordRequest запрос на создание ключевого слова
type CategoryKeywordRequest struct {
	Keyword     string  `json:"keyword" validate:"required"`
	Language    string  `json:"language" validate:"required"`
	Weight      float64 `json:"weight" validate:"gte=0,lte=10"`
	KeywordType string  `json:"keyword_type" validate:"required"`
	IsNegative  bool    `json:"is_negative"`
}

// CategoryKeywordUpdateRequest запрос на обновление ключевого слова
type CategoryKeywordUpdateRequest struct {
	Weight float64 `json:"weight" validate:"gte=0,lte=10"`
}

// CarMake представляет марку автомобиля
type CarMake struct {
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Slug         string    `json:"slug" db:"slug"`
	LogoURL      *string   `json:"logo_url,omitempty" db:"logo_url"`
	Country      *string   `json:"country,omitempty" db:"country"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	SortOrder    int       `json:"sort_order" db:"sort_order"`
	IsDomestic   bool      `json:"is_domestic" db:"is_domestic"`
	PopularityRS int       `json:"popularity_rs" db:"popularity_rs"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// CarModel представляет модель автомобиля
type CarModel struct {
	ID        int       `json:"id" db:"id"`
	MakeID    int       `json:"make_id" db:"make_id"`
	Name      string    `json:"name" db:"name"`
	Slug      string    `json:"slug" db:"slug"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	SortOrder int       `json:"sort_order" db:"sort_order"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	// Поля для JOIN
	Make *CarMake `json:"make,omitempty"`
}

// CarGeneration представляет поколение модели автомобиля
type CarGeneration struct {
	ID        int       `json:"id" db:"id"`
	ModelID   int       `json:"model_id" db:"model_id"`
	Name      string    `json:"name" db:"name"`
	Slug      string    `json:"slug" db:"slug"`
	YearStart int       `json:"year_start" db:"year_start"`
	YearEnd   *int      `json:"year_end,omitempty" db:"year_end"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	SortOrder int       `json:"sort_order" db:"sort_order"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	// Поля для JOIN
	Model *CarModel `json:"model,omitempty"`
}
