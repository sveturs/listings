// backend/internal/domain/models/marketplace_chat.go
package models

import "time"

type MarketplaceMessage struct {
	ID                  int       `json:"id"`
	ChatID              int       `json:"chat_id"`
	ListingID           int       `json:"listing_id"`
	StorefrontProductID int       `json:"storefront_product_id"` // Новое поле для товаров витрин
	SenderID            int       `json:"sender_id"`
	ReceiverID          int       `json:"receiver_id"`
	Content             string    `json:"content"`
	IsRead              bool      `json:"is_read"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`

	// Дополнительные поля для отображения
	Sender            *User               `json:"sender,omitempty"`
	Receiver          *User               `json:"receiver,omitempty"`
	Listing           *MarketplaceListing `json:"listing,omitempty"`
	StorefrontProduct *StorefrontProduct  `json:"storefront_product,omitempty"` // Новое поле для товара витрины
	// Мультиязычность
	OriginalLanguage        string                   `json:"original_language"`
	Translations            map[string]string        `json:"translations,omitempty"` // {"en": "Hello", "ru": "Привет"}
	ChatTranslationMetadata *ChatTranslationMetadata `json:"translation_metadata,omitempty"`

	// Поля для поддержки вложений
	HasAttachments   bool             `json:"has_attachments"`
	AttachmentsCount int              `json:"attachments_count"`
	Attachments      []ChatAttachment `json:"attachments,omitempty"`
}

type MarketplaceChat struct {
	ID                  int       `json:"id"`
	ListingID           int       `json:"listing_id"`
	StorefrontProductID int       `json:"storefront_product_id"` // Новое поле для товаров витрин
	BuyerID             int       `json:"buyer_id"`
	SellerID            int       `json:"seller_id"`
	LastMessageAt       time.Time `json:"last_message_at"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	IsArchived          bool      `json:"is_archived"`

	// Дополнительные поля для отображения
	Buyer       *User               `json:"buyer,omitempty"`
	Seller      *User               `json:"seller,omitempty"`
	OtherUser   *User               `json:"other_user,omitempty"`
	Listing     *MarketplaceListing `json:"listing,omitempty"`
	LastMessage *MarketplaceMessage `json:"last_message,omitempty"`
	UnreadCount int                 `json:"unread_count"`
}

// Структуры для запросов
type CreateMessageRequest struct {
	ListingID           int    `json:"listing_id"`
	StorefrontProductID int    `json:"storefront_product_id"` // Новое поле для товаров витрин
	ChatID              int    `json:"chat_id"`
	ReceiverID          int    `json:"receiver_id" validate:"required"`
	Content             string `json:"content" validate:"required"`
}

type GetMessagesRequest struct {
	ListingID int `query:"listing_id"`
	ChatID    int `query:"chat_id"`
	Page      int `query:"page"`
	Limit     int `query:"limit"`
}

type MarkAsReadRequest struct {
	MessageIDs []int `json:"message_ids" validate:"required"`
}

// ChatActivityStats структура для статистики активности в чате
type ChatActivityStats struct {
	ChatExists        bool      `json:"chat_exists"`
	TotalMessages     int       `json:"total_messages"`
	BuyerMessages     int       `json:"buyer_messages"`
	SellerMessages    int       `json:"seller_messages"`
	DaysSinceFirstMsg int       `json:"days_since_first_msg"`
	DaysSinceLastMsg  int       `json:"days_since_last_msg"`
	FirstMessageDate  time.Time `json:"first_message_date"`
	LastMessageDate   time.Time `json:"last_message_date"`
}

// ChatTranslationMetadata содержит метаинформацию о переводе сообщения
type ChatTranslationMetadata struct {
	TranslatedFrom string    `json:"translated_from"` // "ru"
	TranslatedTo   string    `json:"translated_to"`   // "en"
	TranslatedAt   time.Time `json:"translated_at"`   // Timestamp
	CacheHit       bool      `json:"cache_hit"`       // From Redis cache?
	Provider       string    `json:"provider"`        // "claude-haiku"
}

// ChatUserSettings содержит настройки чата пользователя
type ChatUserSettings struct {
	AutoTranslate     bool   `json:"auto_translate_chat"`
	PreferredLanguage string `json:"preferred_language"` // "ru", "en", "sr"
	ShowLanguageBadge bool   `json:"show_original_language_badge"`
	ModerateTone      bool   `json:"chat_tone_moderation"` // Модерация тона сообщений
}
