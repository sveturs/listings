package domain

import (
	"time"
)

// ChatStatus represents the state of a chat
type ChatStatus string

const (
	ChatStatusActive   ChatStatus = "active"   // Chat is active
	ChatStatusArchived ChatStatus = "archived" // Chat archived by user
	ChatStatusBlocked  ChatStatus = "blocked"  // Chat blocked (spam/abuse)
)

// Chat represents a conversation between buyer and seller
type Chat struct {
	// Identification
	ID       int64 `json:"id" db:"id"`
	BuyerID  int64 `json:"buyer_id" db:"buyer_id"`   // User who initiated chat
	SellerID int64 `json:"seller_id" db:"seller_id"` // User who receives chat

	// Context (what is being discussed)
	ListingID           *int64 `json:"listing_id,omitempty" db:"listing_id"`                       // If discussing marketplace listing
	StorefrontProductID *int64 `json:"storefront_product_id,omitempty" db:"storefront_product_id"` // If discussing B2C product

	// Status
	Status     ChatStatus `json:"status" db:"status"`           // Active, archived, blocked
	IsArchived bool       `json:"is_archived" db:"is_archived"` // Archived by current user

	// Metadata
	LastMessageAt time.Time `json:"last_message_at" db:"last_message_at"` // Last message timestamp
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`

	// Relations (loaded on demand)
	LastMessage  *Message `json:"last_message,omitempty" db:"-"`  // Most recent message (for list view)
	UnreadCount  int32    `json:"unread_count,omitempty" db:"-"`  // Unread messages count for current user
	BuyerName    *string  `json:"buyer_name,omitempty" db:"-"`    // Denormalized for UI
	SellerName   *string  `json:"seller_name,omitempty" db:"-"`   // Denormalized for UI
	ListingTitle *string  `json:"listing_title,omitempty" db:"-"` // Denormalized for UI
}

// IsParticipant checks if user is a participant in the chat
func (c *Chat) IsParticipant(userID int64) bool {
	return c.BuyerID == userID || c.SellerID == userID
}

// GetOtherParticipantID returns the ID of the other participant
func (c *Chat) GetOtherParticipantID(userID int64) int64 {
	if c.BuyerID == userID {
		return c.SellerID
	}
	return c.BuyerID
}

// CanArchive checks if chat can be archived
func (c *Chat) CanArchive() bool {
	return c.Status == ChatStatusActive
}
