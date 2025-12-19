// Package domain defines core business entities for the listings microservice.
// This file contains Stock Reservation domain models for inventory management.
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// StockReservation represents a temporary stock reservation for order processing
type StockReservation struct {
	ID        uuid.UUID `json:"id" db:"id"`
	VariantID uuid.UUID `json:"variant_id" db:"variant_id"`
	OrderID   uuid.UUID `json:"order_id" db:"order_id"`
	Quantity  int32     `json:"quantity" db:"quantity"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	Status    string    `json:"status" db:"status"` // active, confirmed, cancelled, expired
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	// Relations (loaded on demand)
	Variant *ProductVariantV2 `json:"variant,omitempty" db:"-"`
}

// CreateReservationInput represents input for creating a stock reservation
type CreateReservationInput struct {
	VariantID  uuid.UUID `json:"variant_id" validate:"required"`
	OrderID    uuid.UUID `json:"order_id" validate:"required"`
	Quantity   int32     `json:"quantity" validate:"required,gte=1"`
	TTLMinutes int32     `json:"ttl_minutes" validate:"required,gte=1,lte=1440"` // Max 24 hours
}

// UpdateReservationInput represents input for updating a reservation
type UpdateReservationInput struct {
	Quantity  *int32  `json:"quantity,omitempty" validate:"omitempty,gte=1"`
	Status    *string `json:"status,omitempty" validate:"omitempty,oneof=active confirmed cancelled expired"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// ListReservationsFilter represents filters for listing reservations
type ListReservationsFilter struct {
	VariantID  *uuid.UUID `json:"variant_id,omitempty"`
	OrderID    *uuid.UUID `json:"order_id,omitempty"`
	Status     *string    `json:"status,omitempty" validate:"omitempty,oneof=active confirmed cancelled expired"`
	ActiveOnly bool       `json:"active_only"`
	ExpiredOnly bool      `json:"expired_only"`
}

// Stock Reservation status constants (prefixed to avoid conflict with existing ReservationStatus)
const (
	StockReservationStatusActive    = "active"
	StockReservationStatusConfirmed = "confirmed"
	StockReservationStatusCancelled = "cancelled"
	StockReservationStatusExpired   = "expired"
)

// Default reservation TTL in minutes
const (
	DefaultReservationTTL = 30 // 30 minutes
	MaxReservationTTL     = 1440 // 24 hours
)

// IsExpired returns true if reservation has expired
func (r *StockReservation) IsExpired() bool {
	return time.Now().After(r.ExpiresAt)
}

// IsActive returns true if reservation is active and not expired
func (r *StockReservation) IsActive() bool {
	return r.Status == StockReservationStatusActive && !r.IsExpired()
}

// CanBeReleased returns true if reservation can be released (active or expired)
func (r *StockReservation) CanBeReleased() bool {
	return r.Status == StockReservationStatusActive || r.Status == StockReservationStatusExpired
}

// CanBeConfirmed returns true if reservation can be confirmed (active and not expired)
func (r *StockReservation) CanBeConfirmed() bool {
	return r.Status == StockReservationStatusActive && !r.IsExpired()
}

// GetRemainingTime returns duration until expiration
func (r *StockReservation) GetRemainingTime() time.Duration {
	if r.IsExpired() {
		return 0
	}
	return time.Until(r.ExpiresAt)
}

// Validate performs basic validation on the reservation
func (r *StockReservation) Validate() error {
	if r.Quantity <= 0 {
		return ErrInvalidReservationQuantity
	}

	if r.ExpiresAt.Before(time.Now()) {
		return ErrReservationAlreadyExpired
	}

	if r.Status == "" {
		return ErrInvalidReservationStatus
	}

	return nil
}

// Domain errors for StockReservation
var (
	ErrStockReservationNotFound         = errors.New("stock reservation not found")
	ErrInvalidReservationQuantity  = errors.New("reservation quantity must be positive")
	ErrReservationAlreadyExpired   = errors.New("reservation has already expired")
	ErrInvalidReservationStatus    = errors.New("invalid reservation status")
	ErrCannotReleaseReservation    = errors.New("reservation cannot be released in current state")
	ErrCannotConfirmReservation    = errors.New("reservation cannot be confirmed (expired or invalid state)")
	ErrReservationExpired          = errors.New("reservation has expired")
)
