package models

import "time"

type MarketplaceMessage struct {
    ID         int       `json:"id"`
    ChatID     int       `json:"chat_id"`
    ListingID  int       `json:"listing_id"`
    SenderID   int       `json:"sender_id"`
    ReceiverID int       `json:"receiver_id"`
    Content    string    `json:"content"`
    IsRead     bool      `json:"is_read"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
    
    // Дополнительные поля для отображения
    Sender   *User             `json:"sender,omitempty"`
    Receiver *User             `json:"receiver,omitempty"`
    Listing  *MarketplaceListing `json:"listing,omitempty"`
}

type MarketplaceChat struct {
	ID            int       `json:"id"`
	ListingID     int       `json:"listing_id"`
	BuyerID       int       `json:"buyer_id"`
	SellerID      int       `json:"seller_id"`
	LastMessageAt time.Time `json:"last_message_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	IsArchived    bool      `json:"is_archived"`

	// Дополнительные поля для отображения
	Buyer       *User               `json:"buyer,omitempty"`
	Seller      *User               `json:"seller,omitempty"`
	Listing     *MarketplaceListing `json:"listing,omitempty"`
	LastMessage *MarketplaceMessage `json:"last_message,omitempty"`
	UnreadCount int                 `json:"unread_count"`
}

// Структуры для запросов
type CreateMessageRequest struct {
	ListingID  int    `json:"listing_id" validate:"required"`
	ReceiverID int    `json:"receiver_id" validate:"required"`
	Content    string `json:"content" validate:"required"`
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
