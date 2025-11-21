package domain

import (
	"fmt"
	"time"
)

// MessageStatus represents delivery/read status
type MessageStatus string

const (
	MessageStatusSent      MessageStatus = "sent"
	MessageStatusDelivered MessageStatus = "delivered"
	MessageStatusRead      MessageStatus = "read"
	MessageStatusFailed    MessageStatus = "failed"
)

// Message represents a single message in a chat
type Message struct {
	// Identification
	ID         int64 `json:"id"`
	ChatID     int64 `json:"chat_id"`
	SenderID   int64 `json:"sender_id"`
	ReceiverID int64 `json:"receiver_id"`

	// Content
	Content          string `json:"content"`
	OriginalLanguage string `json:"original_language"`

	// Context (optional - inherited from chat if not provided)
	ListingID           *int64 `json:"listing_id,omitempty"`
	StorefrontProductID *int64 `json:"storefront_product_id,omitempty"`

	// Status
	Status MessageStatus `json:"status"`
	IsRead bool          `json:"is_read"`

	// Attachments
	HasAttachments   bool                `json:"has_attachments"`
	AttachmentsCount int32               `json:"attachments_count"`
	Attachments      []*ChatAttachment   `json:"attachments,omitempty"`

	// Timestamps
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	ReadAt    *time.Time `json:"read_at,omitempty"`

	// Denormalized for UI
	SenderName *string `json:"sender_name,omitempty"`
}

// Validate validates the message fields
func (m *Message) Validate() error {
	if m.ChatID == 0 {
		return fmt.Errorf("chat_id is required")
	}
	if m.SenderID == 0 {
		return fmt.Errorf("sender_id is required")
	}
	if m.ReceiverID == 0 {
		return fmt.Errorf("receiver_id is required")
	}
	if m.SenderID == m.ReceiverID {
		return fmt.Errorf("sender_id and receiver_id cannot be the same")
	}

	// Validate content
	if len(m.Content) == 0 {
		return fmt.Errorf("content is required")
	}
	if len(m.Content) > 10000 {
		return fmt.Errorf("content exceeds maximum length of 10000 characters")
	}

	// Validate language code
	if m.OriginalLanguage != "" && len(m.OriginalLanguage) != 2 {
		return fmt.Errorf("original_language must be a 2-letter ISO 639-1 code")
	}

	// Validate status
	if m.Status != "" {
		if m.Status != MessageStatusSent && m.Status != MessageStatusDelivered &&
			m.Status != MessageStatusRead && m.Status != MessageStatusFailed {
			return fmt.Errorf("invalid message status: %s", m.Status)
		}
	}

	// Ensure only one context is set
	contextCount := 0
	if m.ListingID != nil && *m.ListingID > 0 {
		contextCount++
	}
	if m.StorefrontProductID != nil && *m.StorefrontProductID > 0 {
		contextCount++
	}

	if contextCount > 1 {
		return fmt.Errorf("message can only have one context: listing OR storefront_product OR none")
	}

	return nil
}

// MarkAsRead marks the message as read with the current timestamp
func (m *Message) MarkAsRead() {
	m.IsRead = true
	m.Status = MessageStatusRead
	now := time.Now()
	m.ReadAt = &now
}
