// Package domain defines core business entities and domain models for the listings microservice.
package domain

import (
	"encoding/json"
	"errors"
	"time"
)

// ============================================================================
// STOREFRONT EVENT TRACKING
// ============================================================================

// StorefrontEvent represents an analytics event for storefront tracking
type StorefrontEvent struct {
	ID           int64           `json:"id" db:"id"`
	StorefrontID int64           `json:"storefront_id" db:"storefront_id"`
	EventType    string          `json:"event_type" db:"event_type"`
	EventData    json.RawMessage `json:"event_data" db:"event_data"` // JSONB additional data
	UserID       *int64          `json:"user_id,omitempty" db:"user_id"`
	SessionID    string          `json:"session_id" db:"session_id"`
	IPAddress    string          `json:"ip_address" db:"ip_address"`
	UserAgent    string          `json:"user_agent" db:"user_agent"`
	Referrer     string          `json:"referrer" db:"referrer"`
	CreatedAt    time.Time       `json:"created_at" db:"created_at"`
}

// StorefrontEventType defines valid event types
type StorefrontEventType string

const (
	EventTypePageView     StorefrontEventType = "page_view"
	EventTypeProductView  StorefrontEventType = "product_view"
	EventTypeAddToCart    StorefrontEventType = "add_to_cart"
	EventTypeCheckout     StorefrontEventType = "checkout"
	EventTypeOrder        StorefrontEventType = "order"
	EventTypeUnspecified  StorefrontEventType = "unspecified"
)

// Valid returns true if the event type is valid
func (e StorefrontEventType) Valid() bool {
	switch e {
	case EventTypePageView, EventTypeProductView, EventTypeAddToCart, EventTypeCheckout, EventTypeOrder:
		return true
	default:
		return false
	}
}

// String returns string representation of event type
func (e StorefrontEventType) String() string {
	return string(e)
}

// ============================================================================
// VALIDATION
// ============================================================================

// Validate validates the storefront event
func (e *StorefrontEvent) Validate() error {
	if e.StorefrontID <= 0 {
		return errors.New("storefront_id must be positive")
	}

	if e.SessionID == "" {
		return errors.New("session_id is required")
	}

	eventType := StorefrontEventType(e.EventType)
	if !eventType.Valid() {
		return errors.New("invalid event_type")
	}

	return nil
}

// SetDefaults sets default values for optional fields
func (e *StorefrontEvent) SetDefaults() {
	if e.EventData == nil {
		e.EventData = json.RawMessage("{}")
	}

	if e.CreatedAt.IsZero() {
		e.CreatedAt = time.Now()
	}
}
