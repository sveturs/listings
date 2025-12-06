package domain

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// StorefrontInvitationType represents the type of invitation
type StorefrontInvitationType string

const (
	InvitationTypeEmail StorefrontInvitationType = "email" // One-time email invitation
	InvitationTypeLink  StorefrontInvitationType = "link"  // Shareable link invitation
)

// StorefrontInvitationStatus represents the invitation status
type StorefrontInvitationStatus string

const (
	InvitationStatusPending  StorefrontInvitationStatus = "pending"
	InvitationStatusAccepted StorefrontInvitationStatus = "accepted"
	InvitationStatusDeclined StorefrontInvitationStatus = "declined"
	InvitationStatusExpired  StorefrontInvitationStatus = "expired"
	InvitationStatusRevoked  StorefrontInvitationStatus = "revoked"
)

// InviteCodePrefix is the prefix for shareable link codes
const InviteCodePrefix = "sf-"

// StorefrontInvitation represents an invitation to join a storefront staff
type StorefrontInvitation struct {
	// Identification
	ID           int64                        `db:"id" json:"id"`
	StorefrontID int64                        `db:"storefront_id" json:"storefront_id"`
	Role         string                       `db:"role" json:"role"`
	Type         StorefrontInvitationType     `db:"type" json:"type"`

	// Email invitation fields
	InvitedEmail  *string `db:"invited_email" json:"invited_email,omitempty"`
	InvitedUserID *int64  `db:"invited_user_id" json:"invited_user_id,omitempty"`

	// Link invitation fields
	InviteCode   *string    `db:"invite_code" json:"invite_code,omitempty"`
	ExpiresAt    *time.Time `db:"expires_at" json:"expires_at,omitempty"`
	MaxUses      *int32     `db:"max_uses" json:"max_uses,omitempty"`
	CurrentUses  int32      `db:"current_uses" json:"current_uses"`

	// Metadata
	InvitedByID int64  `db:"invited_by_id" json:"invited_by_id"`
	Comment     string `db:"comment" json:"comment,omitempty"`

	// Status
	Status      StorefrontInvitationStatus `db:"status" json:"status"`

	// Timestamps
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
	AcceptedAt  *time.Time `db:"accepted_at" json:"accepted_at,omitempty"`
	DeclinedAt  *time.Time `db:"declined_at" json:"declined_at,omitempty"`
}

// IsExpired checks if the invitation has expired
func (i *StorefrontInvitation) IsExpired() bool {
	if i.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*i.ExpiresAt)
}

// CanAccept checks if the invitation can be accepted
func (i *StorefrontInvitation) CanAccept() bool {
	// Check status
	if i.Status != InvitationStatusPending {
		return false
	}

	// Check if expired
	if i.IsExpired() {
		return false
	}

	// Check max uses for link invitations
	if i.Type == InvitationTypeLink && i.MaxUses != nil {
		if i.CurrentUses >= *i.MaxUses {
			return false
		}
	}

	return true
}

// IncrementUses increments the current uses counter for link invitations
func (i *StorefrontInvitation) IncrementUses() error {
	if i.Type != InvitationTypeLink {
		return fmt.Errorf("cannot increment uses for email invitation")
	}

	if i.MaxUses != nil && i.CurrentUses >= *i.MaxUses {
		return fmt.Errorf("invitation has reached max uses")
	}

	i.CurrentUses++
	return nil
}

// IsLinkInvitation checks if this is a shareable link invitation
func (i *StorefrontInvitation) IsLinkInvitation() bool {
	return i.Type == InvitationTypeLink
}

// IsEmailInvitation checks if this is an email invitation
func (i *StorefrontInvitation) IsEmailInvitation() bool {
	return i.Type == InvitationTypeEmail
}

// IsPending checks if the invitation is still pending
func (i *StorefrontInvitation) IsPending() bool {
	return i.Status == InvitationStatusPending
}

// IsAccepted checks if the invitation has been accepted
func (i *StorefrontInvitation) IsAccepted() bool {
	return i.Status == InvitationStatusAccepted
}

// IsDeclined checks if the invitation has been declined
func (i *StorefrontInvitation) IsDeclined() bool {
	return i.Status == InvitationStatusDeclined
}

// IsRevoked checks if the invitation has been revoked
func (i *StorefrontInvitation) IsRevoked() bool {
	return i.Status == InvitationStatusRevoked
}

// MarkAsAccepted marks the invitation as accepted
func (i *StorefrontInvitation) MarkAsAccepted() {
	i.Status = InvitationStatusAccepted
	now := time.Now()
	i.AcceptedAt = &now
}

// MarkAsDeclined marks the invitation as declined
func (i *StorefrontInvitation) MarkAsDeclined() {
	i.Status = InvitationStatusDeclined
	now := time.Now()
	i.DeclinedAt = &now
}

// MarkAsExpired marks the invitation as expired
func (i *StorefrontInvitation) MarkAsExpired() {
	i.Status = InvitationStatusExpired
}

// MarkAsRevoked marks the invitation as revoked
func (i *StorefrontInvitation) MarkAsRevoked() {
	i.Status = InvitationStatusRevoked
}

// GetRemainingUses returns the remaining uses for link invitations
func (i *StorefrontInvitation) GetRemainingUses() *int32 {
	if i.Type != InvitationTypeLink || i.MaxUses == nil {
		return nil
	}

	remaining := *i.MaxUses - i.CurrentUses
	return &remaining
}

// IsUnlimitedUses checks if the link invitation has unlimited uses
func (i *StorefrontInvitation) IsUnlimitedUses() bool {
	return i.Type == InvitationTypeLink && i.MaxUses == nil
}

// Validate checks if the invitation data is valid
func (i *StorefrontInvitation) Validate() error {
	if i.StorefrontID <= 0 {
		return fmt.Errorf("storefront_id is required")
	}

	if i.InvitedByID <= 0 {
		return fmt.Errorf("invited_by_id is required")
	}

	if i.Role == "" {
		return fmt.Errorf("role is required")
	}

	// Validate role
	validRoles := map[string]bool{
		"owner": true, "manager": true, "staff": true, "cashier": true,
	}
	if !validRoles[i.Role] {
		return fmt.Errorf("invalid role: %s", i.Role)
	}

	// Type-specific validation
	switch i.Type {
	case InvitationTypeEmail:
		if i.InvitedEmail == nil || *i.InvitedEmail == "" {
			return fmt.Errorf("invited_email is required for email invitations")
		}
	case InvitationTypeLink:
		if i.InviteCode == nil || *i.InviteCode == "" {
			return fmt.Errorf("invite_code is required for link invitations")
		}
	default:
		return fmt.Errorf("invalid invitation type: %s", i.Type)
	}

	return nil
}

// GenerateInviteCode generates a unique invite code for shareable link
func GenerateInviteCode() (string, error) {
	// Generate 8 random bytes (will be 16 hex characters)
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random code: %w", err)
	}

	// Convert to hex and add prefix
	code := InviteCodePrefix + hex.EncodeToString(bytes)
	return code, nil
}

// CreateEmailInvitationRequest represents request to create email invitation
type CreateEmailInvitationRequest struct {
	StorefrontID  int64  `json:"storefront_id" validate:"required"`
	InvitedEmail  string `json:"invited_email" validate:"required,email"`
	Role          string `json:"role" validate:"required,oneof=owner manager staff cashier"`
	InvitedByID   int64  `json:"invited_by_id" validate:"required"`
	Comment       string `json:"comment"`
}

// CreateLinkInvitationRequest represents request to create link invitation
type CreateLinkInvitationRequest struct {
	StorefrontID int64      `json:"storefront_id" validate:"required"`
	Role         string     `json:"role" validate:"required,oneof=owner manager staff cashier"`
	InvitedByID  int64      `json:"invited_by_id" validate:"required"`
	Comment      string     `json:"comment"`
	ExpiresAt    *time.Time `json:"expires_at"`
	MaxUses      *int32     `json:"max_uses" validate:"omitempty,min=1"`
}

// ListInvitationsFilter represents filter for listing invitations
type ListInvitationsFilter struct {
	StorefrontID *int64                      `json:"storefront_id"`
	Status       *StorefrontInvitationStatus `json:"status"`
	Type         *StorefrontInvitationType   `json:"type"`
	InvitedByID  *int64                      `json:"invited_by_id"`
	Page         int32                       `json:"page"`
	Limit        int32                       `json:"limit"`
}
