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
		ID                int                     `json:"id"`
		UserID            int                     `json:"user_id"`
		CategoryID        int                     `json:"category_id"`
		Title             string                  `json:"title"`
		Description       string                  `json:"description"`
		Price             float64                 `json:"price"`
		Condition         string                  `json:"condition"`
		Status            string                  `json:"status"`
		Location          string                  `json:"location"`
		Latitude          *float64                `json:"latitude,omitempty"`
		Longitude         *float64                `json:"longitude,omitempty"`
		City              string                  `json:"city"`
		Country           string                  `json:"country"`
		ViewsCount        int                     `json:"views_count"`
		CreatedAt         time.Time               `json:"created_at"`
		UpdatedAt         time.Time               `json:"updated_at"`
		Images            []MarketplaceImage      `json:"images,omitempty"`
		User              *User                   `json:"user,omitempty"`
		Category          *MarketplaceCategory    `json:"category,omitempty"`
		HelpfulVotes      int                     `json:"helpful_votes"`
		NotHelpfulVotes   int                     `json:"not_helpful_votes"`
		IsFavorite        bool                    `json:"is_favorite"`
		ShowOnMap         bool                    `json:"show_on_map"`
		OriginalLanguage  string                  `json:"original_language,omitempty"`
		RawTranslations   interface{}             `json:"-"` // Для хранения "сырых" данных
		Translations      TranslationMap          `json:"translations,omitempty"`
		CategoryPathNames []string                `json:"category_path_names,omitempty"`
		CategoryPathIds   []int                   `json:"category_path_ids,omitempty"`
		CategoryPathSlugs []string                `json:"category_path_slugs,omitempty"`
		CategoryPath      []string                `json:"category_path,omitempty"`
		StorefrontID      *int                    `json:"storefront_id,omitempty"` // связь с витриной
		ExternalID        string                  `json:"external_id,omitempty"`
		Attributes        []ListingAttributeValue `json:"attributes,omitempty"`
		OldPrice          float64                 `json:"old_price,omitempty"`
		HasDiscount       bool                    `json:"has_discount"`
		Metadata          map[string]interface{}  `json:"metadata,omitempty"` // Для хранения дополнительной информации, включая данные о скидке

		AverageRating float64 `json:"average_rating,omitempty"`
		ReviewCount   int     `json:"review_count,omitempty"`
	}
)

type MarketplaceCategory struct {
	ID                int               `json:"id"`
	Name              string            `json:"name"`
	Slug              string            `json:"slug"`
	ParentID          *int              `json:"parent_id,omitempty"`
	Icon              string            `json:"icon,omitempty"`
	Description       string            `json:"description,omitempty"`
	IsActive          bool              `json:"is_active"`
	CreatedAt         time.Time         `json:"created_at"`
	Translations      map[string]string `json:"translations,omitempty"`
	ListingCount      int               `json:"listing_count"`
	HasCustomUI       bool              `json:"has_custom_ui,omitempty"`
	CustomUIComponent string            `json:"custom_ui_component,omitempty"`
	SortOrder         int               `json:"sort_order"`
	Level             int               `json:"level"`
	Count             int               `json:"count"`
	ExternalID        string            `json:"external_id,omitempty"`
	SEOTitle          string            `json:"seo_title,omitempty"`
	SEODescription    string            `json:"seo_description,omitempty"`
	SEOKeywords       string            `json:"seo_keywords,omitempty"`
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
