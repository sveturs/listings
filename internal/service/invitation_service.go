package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/vondi-global/listings/internal/domain"
	"github.com/vondi-global/listings/internal/repository/postgres"
)

// InvitationService handles storefront invitation business logic
type InvitationService struct {
	invitationRepo      *postgres.StorefrontInvitationRepository
	repo                *postgres.Repository
	storefrontStaffRepo *postgres.StorefrontStaffRepository
	logger              zerolog.Logger
}

// NewInvitationService creates a new invitation service
func NewInvitationService(
	db *sql.DB,
	dbx *sqlx.DB,
	logger zerolog.Logger,
) *InvitationService {
	return &InvitationService{
		invitationRepo:      postgres.NewStorefrontInvitationRepository(db),
		repo:                postgres.NewRepository(dbx, logger),
		storefrontStaffRepo: postgres.NewStorefrontStaffRepository(db),
		logger:              logger.With().Str("service", "invitation").Logger(),
	}
}

// CreateEmailInvitation creates an email invitation
func (s *InvitationService) CreateEmailInvitation(ctx context.Context, req *domain.CreateEmailInvitationRequest) (*domain.StorefrontInvitation, error) {
	s.logger.Info().
		Int64("storefront_id", req.StorefrontID).
		Str("email", req.InvitedEmail).
		Str("role", req.Role).
		Msg("Creating email invitation")

	// Validate storefront exists
	storefront, err := s.repo.GetStorefrontByID(ctx, req.StorefrontID, nil)
	if err != nil {
		return nil, fmt.Errorf("storefront not found: %w", err)
	}
	if storefront == nil {
		return nil, fmt.Errorf("storefront not found")
	}

	// Check if user already has pending invitation
	existingInv, err := s.invitationRepo.GetByEmailAndStorefront(ctx, req.InvitedEmail, req.StorefrontID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing invitation: %w", err)
	}
	if existingInv != nil {
		return nil, fmt.Errorf("pending invitation already exists for this email")
	}

	// Create invitation
	inv := &domain.StorefrontInvitation{
		StorefrontID: req.StorefrontID,
		Type:         domain.InvitationTypeEmail,
		Role:         req.Role,
		InvitedEmail: &req.InvitedEmail,
		InvitedByID:  req.InvitedByID,
		Comment:      req.Comment,
		Status:       domain.InvitationStatusPending,
		CurrentUses:  0,
	}

	// Validate
	if err := inv.Validate(); err != nil {
		return nil, fmt.Errorf("invalid invitation: %w", err)
	}

	// Save to database
	if err := s.invitationRepo.Create(ctx, inv); err != nil {
		return nil, fmt.Errorf("failed to create invitation: %w", err)
	}

	s.logger.Info().
		Int64("invitation_id", inv.ID).
		Str("email", req.InvitedEmail).
		Msg("Email invitation created successfully")

	return inv, nil
}

// CreateLinkInvitation creates a shareable link invitation
func (s *InvitationService) CreateLinkInvitation(ctx context.Context, req *domain.CreateLinkInvitationRequest) (*domain.StorefrontInvitation, error) {
	s.logger.Info().
		Int64("storefront_id", req.StorefrontID).
		Str("role", req.Role).
		Int32("max_uses", *req.MaxUses).
		Msg("Creating link invitation")

	// Validate storefront exists
	storefront, err := s.repo.GetStorefrontByID(ctx, req.StorefrontID, nil)
	if err != nil {
		return nil, fmt.Errorf("storefront not found: %w", err)
	}
	if storefront == nil {
		return nil, fmt.Errorf("storefront not found")
	}

	// Generate unique invite code
	var inviteCode string
	maxAttempts := 10
	for i := 0; i < maxAttempts; i++ {
		code, err := domain.GenerateInviteCode()
		if err != nil {
			return nil, fmt.Errorf("failed to generate invite code: %w", err)
		}

		// Check if code already exists
		exists, err := s.invitationRepo.CheckInviteCodeExists(ctx, code)
		if err != nil {
			return nil, fmt.Errorf("failed to check invite code: %w", err)
		}

		if !exists {
			inviteCode = code
			break
		}
	}

	if inviteCode == "" {
		return nil, fmt.Errorf("failed to generate unique invite code after %d attempts", maxAttempts)
	}

	// Calculate expiration
	expiresAt := req.ExpiresAt
	if expiresAt == nil {
		// Default expiration if not provided
		expires := time.Now().Add(7 * 24 * time.Hour) // 7 days
		expiresAt = &expires
	}

	// Create invitation
	inv := &domain.StorefrontInvitation{
		StorefrontID: req.StorefrontID,
		Type:         domain.InvitationTypeLink,
		Role:         req.Role,
		InviteCode:   &inviteCode,
		ExpiresAt:    expiresAt,
		MaxUses:      req.MaxUses,
		CurrentUses:  0,
		InvitedByID:  req.InvitedByID,
		Comment:      req.Comment,
		Status:       domain.InvitationStatusPending,
	}

	// Validate
	if err := inv.Validate(); err != nil {
		return nil, fmt.Errorf("invalid invitation: %w", err)
	}

	// Save to database
	if err := s.invitationRepo.Create(ctx, inv); err != nil {
		return nil, fmt.Errorf("failed to create invitation: %w", err)
	}

	s.logger.Info().
		Int64("invitation_id", inv.ID).
		Str("code", inviteCode).
		Msg("Link invitation created successfully")

	return inv, nil
}

// ListInvitations lists invitations for a storefront with filtering
func (s *InvitationService) ListInvitations(ctx context.Context, storefrontID int64, filter *domain.ListInvitationsFilter) ([]*domain.StorefrontInvitation, int32, error) {
	s.logger.Debug().
		Int64("storefront_id", storefrontID).
		Msg("Listing invitations")

	// Ensure storefront_id is set in filter
	if filter == nil {
		filter = &domain.ListInvitationsFilter{}
	}
	filter.StorefrontID = &storefrontID

	return s.invitationRepo.ListByStorefront(ctx, storefrontID, filter)
}

// GetMyInvitations gets pending invitations for a user
func (s *InvitationService) GetMyInvitations(ctx context.Context, userID int64, email string) ([]*domain.StorefrontInvitation, error) {
	s.logger.Debug().
		Int64("user_id", userID).
		Str("email", email).
		Msg("Getting user invitations")

	return s.invitationRepo.ListByUser(ctx, userID, email)
}

// AcceptInvitation accepts an invitation by ID
func (s *InvitationService) AcceptInvitation(ctx context.Context, invitationID int64, userID int64) error {
	s.logger.Info().
		Int64("invitation_id", invitationID).
		Int64("user_id", userID).
		Msg("Accepting invitation")

	// Get invitation
	inv, err := s.invitationRepo.GetByID(ctx, invitationID)
	if err != nil {
		return fmt.Errorf("invitation not found: %w", err)
	}

	// Check if invitation can be accepted
	if !inv.CanAccept() {
		return fmt.Errorf("invitation cannot be accepted: status=%s, expired=%v", inv.Status, inv.IsExpired())
	}

	// For email invitations, verify user email matches
	if inv.IsEmailInvitation() {
		if inv.InvitedEmail == nil {
			return fmt.Errorf("email invitation has no invited email")
		}
		// TODO: Verify userID email matches invited email (need user service)
	}

	// Add user to storefront staff
	staff := &domain.StorefrontStaff{
		StorefrontID: inv.StorefrontID,
		UserID:       userID,
		Role:         inv.Role,
	}

	if err := s.storefrontStaffRepo.Create(ctx, staff); err != nil {
		return fmt.Errorf("failed to add user to storefront staff: %w", err)
	}

	// Mark invitation as accepted
	inv.MarkAsAccepted()
	if inv.IsLinkInvitation() {
		if err := inv.IncrementUses(); err != nil {
			return fmt.Errorf("failed to increment uses: %w", err)
		}
	}

	// Update invitation
	if err := s.invitationRepo.Update(ctx, inv); err != nil {
		return fmt.Errorf("failed to update invitation: %w", err)
	}

	s.logger.Info().
		Int64("invitation_id", invitationID).
		Int64("user_id", userID).
		Int64("storefront_id", inv.StorefrontID).
		Msg("Invitation accepted successfully")

	return nil
}

// DeclineInvitation declines an invitation
func (s *InvitationService) DeclineInvitation(ctx context.Context, invitationID int64, userID int64) error {
	s.logger.Info().
		Int64("invitation_id", invitationID).
		Int64("user_id", userID).
		Msg("Declining invitation")

	// Get invitation
	inv, err := s.invitationRepo.GetByID(ctx, invitationID)
	if err != nil {
		return fmt.Errorf("invitation not found: %w", err)
	}

	// Check status
	if inv.Status != domain.InvitationStatusPending {
		return fmt.Errorf("invitation cannot be declined: status=%s", inv.Status)
	}

	// Mark as declined
	inv.MarkAsDeclined()

	// Update invitation
	if err := s.invitationRepo.Update(ctx, inv); err != nil {
		return fmt.Errorf("failed to update invitation: %w", err)
	}

	s.logger.Info().
		Int64("invitation_id", invitationID).
		Msg("Invitation declined successfully")

	return nil
}

// RevokeInvitation revokes an invitation
func (s *InvitationService) RevokeInvitation(ctx context.Context, invitationID int64, revokedByID int64) error {
	s.logger.Info().
		Int64("invitation_id", invitationID).
		Int64("revoked_by_id", revokedByID).
		Msg("Revoking invitation")

	// Get invitation
	inv, err := s.invitationRepo.GetByID(ctx, invitationID)
	if err != nil {
		return fmt.Errorf("invitation not found: %w", err)
	}

	// Check status
	if inv.Status != domain.InvitationStatusPending {
		return fmt.Errorf("invitation cannot be revoked: status=%s", inv.Status)
	}

	// Mark as revoked
	if err := s.invitationRepo.MarkAsRevoked(ctx, invitationID); err != nil {
		return fmt.Errorf("failed to revoke invitation: %w", err)
	}

	s.logger.Info().
		Int64("invitation_id", invitationID).
		Msg("Invitation revoked successfully")

	return nil
}

// ValidateInviteCode validates an invite code and returns the invitation
func (s *InvitationService) ValidateInviteCode(ctx context.Context, code string) (*domain.StorefrontInvitation, string, error) {
	s.logger.Debug().
		Str("code", code).
		Msg("Validating invite code")

	// Get invitation by code
	inv, err := s.invitationRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, "", fmt.Errorf("invitation not found: %w", err)
	}

	// Get storefront name
	storefront, err := s.repo.GetStorefrontByID(ctx, inv.StorefrontID, nil)
	if err != nil {
		return nil, "", fmt.Errorf("storefront not found: %w", err)
	}

	storefrontName := ""
	if storefront != nil {
		storefrontName = storefront.Name
	}

	// Check if invitation can be accepted
	if !inv.CanAccept() {
		var errorCode string
		switch {
		case inv.IsExpired():
			errorCode = "expired"
		case inv.Status == domain.InvitationStatusRevoked:
			errorCode = "revoked"
		case inv.MaxUses != nil && inv.CurrentUses >= *inv.MaxUses:
			errorCode = "max_uses_reached"
		default:
			errorCode = "invalid"
		}
		return inv, storefrontName, fmt.Errorf("invitation cannot be accepted: %s", errorCode)
	}

	return inv, storefrontName, nil
}

// AcceptInviteCode accepts an invitation via invite code
func (s *InvitationService) AcceptInviteCode(ctx context.Context, code string, userID int64) error {
	s.logger.Info().
		Str("code", code).
		Int64("user_id", userID).
		Msg("Accepting invitation via code")

	// Get invitation by code
	inv, err := s.invitationRepo.GetByCode(ctx, code)
	if err != nil {
		return fmt.Errorf("invitation not found: %w", err)
	}

	// Use AcceptInvitation
	return s.AcceptInvitation(ctx, inv.ID, userID)
}

// ExpirePendingInvitations is a background job to expire pending invitations
func (s *InvitationService) ExpirePendingInvitations(ctx context.Context) (int64, error) {
	s.logger.Info().Msg("Expiring pending invitations")

	count, err := s.invitationRepo.ExpirePendingInvitations(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to expire invitations: %w", err)
	}

	s.logger.Info().Int64("count", count).Msg("Expired pending invitations")
	return count, nil
}
