// Package domain defines core business entities and domain models for the listings microservice.
package domain

import (
	"errors"
	"time"

	pb "github.com/vondi-global/listings/api/proto/listings/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ReservationStatus represents the state of inventory reservation
type ReservationStatus string

const (
	ReservationStatusUnspecified ReservationStatus = "unspecified"
	ReservationStatusActive      ReservationStatus = "active"    // Reservation active, stock held
	ReservationStatusCommitted   ReservationStatus = "committed" // Reservation committed (order confirmed)
	ReservationStatusReleased    ReservationStatus = "released"  // Reservation released, stock restored
	ReservationStatusExpired     ReservationStatus = "expired"   // Reservation expired (TTL exceeded)
)

// ReferenceType defines the type of entity that owns the reservation
type ReferenceType string

const (
	ReferenceTypeOrder    ReferenceType = "order"
	ReferenceTypeTransfer ReferenceType = "transfer"
)

// InventoryReservation represents a temporary stock hold
type InventoryReservation struct {
	ID            int64             `json:"id" db:"id"`
	ListingID     int64             `json:"listing_id" db:"listing_id"`           // FK to listing or product
	VariantID     *int64            `json:"variant_id,omitempty" db:"variant_id"` // FK to variant (if applicable)
	ReferenceType ReferenceType     `json:"reference_type" db:"reference_type"`   // Type: "order", "transfer"
	ReferenceID   int64             `json:"reference_id" db:"reference_id"`       // FK to order or transfer
	Quantity      int32             `json:"quantity" db:"quantity"`               // Quantity reserved
	Status        ReservationStatus `json:"status" db:"status"`                   // Reservation state
	ExpiresAt     time.Time         `json:"expires_at" db:"expires_at"`           // TTL for reservation
	CreatedAt     time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at" db:"updated_at"`
	CommittedAt   *time.Time        `json:"committed_at,omitempty" db:"committed_at"` // When reservation committed
	ReleasedAt    *time.Time        `json:"released_at,omitempty" db:"released_at"`   // When reservation released
}

// Validate validates the InventoryReservation entity
func (r *InventoryReservation) Validate() error {
	if r == nil {
		return errors.New("reservation cannot be nil")
	}

	if r.ListingID <= 0 {
		return errors.New("listing_id must be greater than 0")
	}

	if r.ReferenceType != ReferenceTypeOrder && r.ReferenceType != ReferenceTypeTransfer {
		return errors.New("reference_type must be 'order' or 'transfer'")
	}

	if r.ReferenceID <= 0 {
		return errors.New("reference_id must be greater than 0")
	}

	if r.Quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	if r.ExpiresAt.IsZero() {
		return errors.New("expires_at is required")
	}

	return nil
}

// IsExpired checks if the reservation has expired
func (r *InventoryReservation) IsExpired() bool {
	if r == nil {
		return true
	}

	return time.Now().After(r.ExpiresAt) && r.Status == ReservationStatusActive
}

// CanCommit checks if reservation can be committed
func (r *InventoryReservation) CanCommit() bool {
	if r == nil {
		return false
	}

	// Can only commit active reservations that haven't expired
	return r.Status == ReservationStatusActive && !r.IsExpired()
}

// CanRelease checks if reservation can be released
func (r *InventoryReservation) CanRelease() bool {
	if r == nil {
		return false
	}

	// Can release active reservations (expired or not)
	return r.Status == ReservationStatusActive
}

// Commit marks the reservation as committed
func (r *InventoryReservation) Commit() error {
	if !r.CanCommit() {
		return errors.New("reservation cannot be committed")
	}

	now := time.Now()
	r.Status = ReservationStatusCommitted
	r.CommittedAt = &now
	r.UpdatedAt = now

	return nil
}

// Release marks the reservation as released
func (r *InventoryReservation) Release() error {
	if !r.CanRelease() {
		return errors.New("reservation cannot be released")
	}

	now := time.Now()
	r.Status = ReservationStatusReleased
	r.ReleasedAt = &now
	r.UpdatedAt = now

	return nil
}

// Expire marks the reservation as expired
func (r *InventoryReservation) Expire() error {
	if r == nil {
		return errors.New("reservation cannot be nil")
	}

	if r.Status != ReservationStatusActive {
		return errors.New("only active reservations can be expired")
	}

	now := time.Now()
	r.Status = ReservationStatusExpired
	r.UpdatedAt = now

	return nil
}

// CalculateTTL calculates the time-to-live for the reservation in minutes
func (r *InventoryReservation) CalculateTTL() int64 {
	if r == nil || r.ExpiresAt.IsZero() {
		return 0
	}

	ttl := time.Until(r.ExpiresAt)
	if ttl < 0 {
		return 0
	}

	return int64(ttl.Minutes())
}

// NewInventoryReservation creates a new inventory reservation with default TTL (30 minutes)
func NewInventoryReservation(listingID int64, variantID *int64, refType ReferenceType, refID int64, quantity int32) *InventoryReservation {
	now := time.Now()
	expiresAt := now.Add(30 * time.Minute) // Default TTL: 30 minutes

	return &InventoryReservation{
		ListingID:     listingID,
		VariantID:     variantID,
		ReferenceType: refType,
		ReferenceID:   refID,
		Quantity:      quantity,
		Status:        ReservationStatusActive,
		ExpiresAt:     expiresAt,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// NewInventoryReservationWithTTL creates a new inventory reservation with custom TTL
func NewInventoryReservationWithTTL(listingID int64, variantID *int64, refType ReferenceType, refID int64, quantity int32, ttlMinutes int) *InventoryReservation {
	now := time.Now()
	expiresAt := now.Add(time.Duration(ttlMinutes) * time.Minute)

	return &InventoryReservation{
		ListingID:     listingID,
		VariantID:     variantID,
		ReferenceType: refType,
		ReferenceID:   refID,
		Quantity:      quantity,
		Status:        ReservationStatusActive,
		ExpiresAt:     expiresAt,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// ReservationStatusFromProto converts proto ReservationStatus to domain ReservationStatus
func ReservationStatusFromProto(pbStatus pb.ReservationStatus) ReservationStatus {
	switch pbStatus {
	case pb.ReservationStatus_RESERVATION_STATUS_ACTIVE:
		return ReservationStatusActive
	case pb.ReservationStatus_RESERVATION_STATUS_COMMITTED:
		return ReservationStatusCommitted
	case pb.ReservationStatus_RESERVATION_STATUS_RELEASED:
		return ReservationStatusReleased
	case pb.ReservationStatus_RESERVATION_STATUS_EXPIRED:
		return ReservationStatusExpired
	default:
		return ReservationStatusUnspecified
	}
}

// ToProtoReservationStatus converts domain ReservationStatus to proto ReservationStatus
func (s ReservationStatus) ToProtoReservationStatus() pb.ReservationStatus {
	switch s {
	case ReservationStatusActive:
		return pb.ReservationStatus_RESERVATION_STATUS_ACTIVE
	case ReservationStatusCommitted:
		return pb.ReservationStatus_RESERVATION_STATUS_COMMITTED
	case ReservationStatusReleased:
		return pb.ReservationStatus_RESERVATION_STATUS_RELEASED
	case ReservationStatusExpired:
		return pb.ReservationStatus_RESERVATION_STATUS_EXPIRED
	default:
		return pb.ReservationStatus_RESERVATION_STATUS_UNSPECIFIED
	}
}

// InventoryReservationFromProto converts proto ReservationDetails to domain InventoryReservation
func InventoryReservationFromProto(details *pb.ReservationDetails) *InventoryReservation {
	if details == nil {
		return nil
	}

	reservation := &InventoryReservation{
		ID:            details.Id,
		ListingID:     details.ListingId,
		ReferenceType: ReferenceType(details.ReferenceType),
		ReferenceID:   details.ReferenceId,
		Quantity:      details.Quantity,
		Status:        ReservationStatus(details.Status),
	}

	if details.VariantId != nil {
		variantID := *details.VariantId
		reservation.VariantID = &variantID
	}

	if details.ExpiresAt != nil {
		reservation.ExpiresAt = details.ExpiresAt.AsTime()
	}

	if details.CreatedAt != nil {
		reservation.CreatedAt = details.CreatedAt.AsTime()
	}

	if details.UpdatedAt != nil {
		reservation.UpdatedAt = details.UpdatedAt.AsTime()
	}

	return reservation
}

// ToProto converts domain InventoryReservation to proto ReservationDetails
func (r *InventoryReservation) ToProto() *pb.ReservationDetails {
	if r == nil {
		return nil
	}

	details := &pb.ReservationDetails{
		Id:            r.ID,
		ListingId:     r.ListingID,
		Quantity:      r.Quantity,
		Status:        string(r.Status),
		ReferenceType: string(r.ReferenceType),
		ReferenceId:   r.ReferenceID,
		ExpiresAt:     timestamppb.New(r.ExpiresAt),
		CreatedAt:     timestamppb.New(r.CreatedAt),
		UpdatedAt:     timestamppb.New(r.UpdatedAt),
	}

	if r.VariantID != nil {
		details.VariantId = r.VariantID
	}

	return details
}
