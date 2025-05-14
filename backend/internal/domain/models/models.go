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

type User struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	GoogleID   string    `json:"google_id"`
	PictureURL string    `json:"picture_url"`
	Phone      *string   `json:"phone,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}
type TranslationMap map[string]map[string]string
type MarketplaceListing struct {
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

type MarketplaceCategory struct {
	ID                int               `json:"id"`
	Name              string            `json:"name"`
	Slug              string            `json:"slug"`
	ParentID          *int              `json:"parent_id,omitempty"`
	Icon              string            `json:"icon,omitempty"`
	CreatedAt         time.Time         `json:"created_at"`
	Translations      map[string]string `json:"translations,omitempty"`
	ListingCount      int               `json:"listing_count"`
	HasCustomUI       bool              `json:"has_custom_ui,omitempty"`
	CustomUIComponent string            `json:"custom_ui_component,omitempty"`
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
	CreatedAt     time.Time `json:"created_at"`
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
