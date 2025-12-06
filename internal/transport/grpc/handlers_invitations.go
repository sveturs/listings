package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	listingspb "github.com/vondi-global/listings/api/proto/listings/v1"
	"github.com/vondi-global/listings/internal/domain"
)

// CreateStorefrontEmailInvitation creates an email invitation
func (s *Server) CreateStorefrontEmailInvitation(ctx context.Context, req *listingspb.CreateStorefrontEmailInvitationRequest) (*listingspb.StorefrontInvitation, error) {
	s.logger.Info().
		Int64("storefront_id", req.StorefrontId).
		Str("email", req.Email).
		Str("role", req.Role).
		Msg("CreateStorefrontEmailInvitation called")

	// Validate request
	if req.StorefrontId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "storefront_id is required")
	}
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}
	if req.Role == "" {
		return nil, status.Error(codes.InvalidArgument, "role is required")
	}
	if req.InvitedById <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invited_by_id is required")
	}

	// Convert to domain request
	domainReq := &domain.CreateEmailInvitationRequest{
		StorefrontID: req.StorefrontId,
		InvitedEmail: req.Email,
		Role:         req.Role,
		InvitedByID:  req.InvitedById,
		Comment:      stringPtrFromOptional(req.Comment),
	}

	// Create invitation
	inv, err := s.invitationService.CreateEmailInvitation(ctx, domainReq)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to create email invitation")
		return nil, status.Errorf(codes.Internal, "failed to create invitation: %v", err)
	}

	s.logger.Info().Int64("invitation_id", inv.ID).Msg("Email invitation created successfully")

	return domainInvitationToProto(inv), nil
}

// CreateStorefrontLinkInvitation creates a shareable link invitation
func (s *Server) CreateStorefrontLinkInvitation(ctx context.Context, req *listingspb.CreateStorefrontLinkInvitationRequest) (*listingspb.StorefrontInvitation, error) {
	s.logger.Info().
		Int64("storefront_id", req.StorefrontId).
		Str("role", req.Role).
		Int32("expires_in_days", req.ExpiresInDays).
		Int32("max_uses", req.MaxUses).
		Msg("CreateStorefrontLinkInvitation called")

	// Validate request
	if req.StorefrontId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "storefront_id is required")
	}
	if req.Role == "" {
		return nil, status.Error(codes.InvalidArgument, "role is required")
	}
	if req.InvitedById <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invited_by_id is required")
	}
	if req.ExpiresInDays <= 0 {
		return nil, status.Error(codes.InvalidArgument, "expires_in_days must be greater than 0")
	}

	// Calculate expiration time
	expiresAt := time.Now().Add(time.Duration(req.ExpiresInDays) * 24 * time.Hour)

	// Handle max_uses (0 means unlimited)
	var maxUses *int32
	if req.MaxUses > 0 {
		maxUses = &req.MaxUses
	}

	// Convert to domain request
	domainReq := &domain.CreateLinkInvitationRequest{
		StorefrontID: req.StorefrontId,
		Role:         req.Role,
		InvitedByID:  req.InvitedById,
		Comment:      stringPtrFromOptional(req.Comment),
		ExpiresAt:    &expiresAt,
		MaxUses:      maxUses,
	}

	// Create invitation
	inv, err := s.invitationService.CreateLinkInvitation(ctx, domainReq)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to create link invitation")
		return nil, status.Errorf(codes.Internal, "failed to create invitation: %v", err)
	}

	s.logger.Info().Int64("invitation_id", inv.ID).Msg("Link invitation created successfully")

	return domainInvitationToProto(inv), nil
}

// ListStorefrontInvitations lists invitations for a storefront
func (s *Server) ListStorefrontInvitations(ctx context.Context, req *listingspb.ListStorefrontInvitationsRequest) (*listingspb.ListStorefrontInvitationsResponse, error) {
	s.logger.Debug().
		Int64("storefront_id", req.StorefrontId).
		Msg("ListStorefrontInvitations called")

	// Validate request
	if req.StorefrontId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "storefront_id is required")
	}

	// Build filter
	filter := &domain.ListInvitationsFilter{
		Limit:  req.Limit,
		Page:   req.Offset/req.Limit + 1, // Convert offset to page
	}

	if req.StatusFilter != nil {
		statusStr := protoStatusToDomainStatus(*req.StatusFilter)
		filter.Status = &statusStr
	}

	// Get invitations
	invitations, total, err := s.invitationService.ListInvitations(ctx, req.StorefrontId, filter)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to list invitations")
		return nil, status.Errorf(codes.Internal, "failed to list invitations: %v", err)
	}

	// Convert to proto
	protoInvitations := make([]*listingspb.StorefrontInvitation, len(invitations))
	for i, inv := range invitations {
		protoInvitations[i] = domainInvitationToProto(inv)
	}

	return &listingspb.ListStorefrontInvitationsResponse{
		Invitations: protoInvitations,
		Total:       total,
	}, nil
}

// RevokeStorefrontInvitation revokes an invitation
func (s *Server) RevokeStorefrontInvitation(ctx context.Context, req *listingspb.RevokeStorefrontInvitationRequest) (*emptypb.Empty, error) {
	s.logger.Info().
		Int64("invitation_id", req.InvitationId).
		Int64("revoked_by_id", req.RevokedById).
		Msg("RevokeStorefrontInvitation called")

	// Validate request
	if req.InvitationId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invitation_id is required")
	}
	if req.RevokedById <= 0 {
		return nil, status.Error(codes.InvalidArgument, "revoked_by_id is required")
	}

	// Revoke invitation
	if err := s.invitationService.RevokeInvitation(ctx, req.InvitationId, req.RevokedById); err != nil {
		s.logger.Error().Err(err).Msg("Failed to revoke invitation")
		return nil, status.Errorf(codes.Internal, "failed to revoke invitation: %v", err)
	}

	s.logger.Info().Int64("invitation_id", req.InvitationId).Msg("Invitation revoked successfully")

	return &emptypb.Empty{}, nil
}

// GetMyStorefrontInvitations gets pending invitations for a user
func (s *Server) GetMyStorefrontInvitations(ctx context.Context, req *listingspb.GetMyStorefrontInvitationsRequest) (*listingspb.ListStorefrontInvitationsResponse, error) {
	s.logger.Debug().
		Int64("user_id", req.UserId).
		Msg("GetMyStorefrontInvitations called")

	// Validate request
	if req.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	email := ""
	if req.Email != nil {
		email = *req.Email
	}

	// Get user invitations
	invitations, err := s.invitationService.GetMyInvitations(ctx, req.UserId, email)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to get user invitations")
		return nil, status.Errorf(codes.Internal, "failed to get invitations: %v", err)
	}

	// Convert to proto
	protoInvitations := make([]*listingspb.StorefrontInvitation, len(invitations))
	for i, inv := range invitations {
		protoInvitations[i] = domainInvitationToProto(inv)
	}

	return &listingspb.ListStorefrontInvitationsResponse{
		Invitations: protoInvitations,
		Total:       int32(len(invitations)),
	}, nil
}

// AcceptStorefrontInvitation accepts an invitation
func (s *Server) AcceptStorefrontInvitation(ctx context.Context, req *listingspb.AcceptStorefrontInvitationRequest) (*emptypb.Empty, error) {
	s.logger.Info().
		Int64("invitation_id", req.InvitationId).
		Int64("user_id", req.UserId).
		Msg("AcceptStorefrontInvitation called")

	// Validate request
	if req.InvitationId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invitation_id is required")
	}
	if req.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	// Accept invitation
	if err := s.invitationService.AcceptInvitation(ctx, req.InvitationId, req.UserId); err != nil {
		s.logger.Error().Err(err).Msg("Failed to accept invitation")
		return nil, status.Errorf(codes.Internal, "failed to accept invitation: %v", err)
	}

	s.logger.Info().Int64("invitation_id", req.InvitationId).Msg("Invitation accepted successfully")

	return &emptypb.Empty{}, nil
}

// DeclineStorefrontInvitation declines an invitation
func (s *Server) DeclineStorefrontInvitation(ctx context.Context, req *listingspb.DeclineStorefrontInvitationRequest) (*emptypb.Empty, error) {
	s.logger.Info().
		Int64("invitation_id", req.InvitationId).
		Int64("user_id", req.UserId).
		Msg("DeclineStorefrontInvitation called")

	// Validate request
	if req.InvitationId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invitation_id is required")
	}
	if req.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	// Decline invitation
	if err := s.invitationService.DeclineInvitation(ctx, req.InvitationId, req.UserId); err != nil {
		s.logger.Error().Err(err).Msg("Failed to decline invitation")
		return nil, status.Errorf(codes.Internal, "failed to decline invitation: %v", err)
	}

	s.logger.Info().Int64("invitation_id", req.InvitationId).Msg("Invitation declined successfully")

	return &emptypb.Empty{}, nil
}

// ValidateStorefrontInviteCode validates an invite code
func (s *Server) ValidateStorefrontInviteCode(ctx context.Context, req *listingspb.ValidateStorefrontInviteCodeRequest) (*listingspb.ValidateStorefrontInviteCodeResponse, error) {
	s.logger.Debug().
		Str("code", req.Code).
		Msg("ValidateStorefrontInviteCode called")

	// Validate request
	if req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "code is required")
	}

	// Validate invite code
	inv, storefrontName, err := s.invitationService.ValidateInviteCode(ctx, req.Code)
	if err != nil {
		// Return validation result with error details
		errorCode := "invalid"
		if inv != nil {
			switch {
			case inv.IsExpired():
				errorCode = "expired"
			case inv.Status == domain.InvitationStatusRevoked:
				errorCode = "revoked"
			case inv.MaxUses != nil && inv.CurrentUses >= *inv.MaxUses:
				errorCode = "max_uses_reached"
			}
		} else {
			errorCode = "not_found"
		}

		resp := &listingspb.ValidateStorefrontInviteCodeResponse{
			IsValid:   false,
			ErrorCode: &errorCode,
		}

		if inv != nil {
			resp.Invitation = domainInvitationToProto(inv)
			resp.StorefrontName = storefrontName
		}

		return resp, nil
	}

	// Valid invitation
	return &listingspb.ValidateStorefrontInviteCodeResponse{
		Invitation:     domainInvitationToProto(inv),
		StorefrontName: storefrontName,
		IsValid:        true,
	}, nil
}

// AcceptStorefrontInviteCode accepts an invitation via code
func (s *Server) AcceptStorefrontInviteCode(ctx context.Context, req *listingspb.AcceptStorefrontInviteCodeRequest) (*emptypb.Empty, error) {
	s.logger.Info().
		Str("code", req.Code).
		Int64("user_id", req.UserId).
		Msg("AcceptStorefrontInviteCode called")

	// Validate request
	if req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "code is required")
	}
	if req.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	// Accept invitation via code
	if err := s.invitationService.AcceptInviteCode(ctx, req.Code, req.UserId); err != nil {
		s.logger.Error().Err(err).Msg("Failed to accept invitation via code")
		return nil, status.Errorf(codes.Internal, "failed to accept invitation: %v", err)
	}

	s.logger.Info().Str("code", req.Code).Msg("Invitation accepted via code successfully")

	return &emptypb.Empty{}, nil
}

// ============================================================================
// Helper functions for conversion
// ============================================================================

// domainInvitationToProto converts domain invitation to proto
func domainInvitationToProto(inv *domain.StorefrontInvitation) *listingspb.StorefrontInvitation {
	if inv == nil {
		return nil
	}

	proto := &listingspb.StorefrontInvitation{
		Id:           inv.ID,
		StorefrontId: inv.StorefrontID,
		Type:         domainTypeToProtoType(inv.Type),
		Role:         inv.Role,
		Status:       domainStatusToProtoStatus(inv.Status),
		CurrentUses:  inv.CurrentUses,
		InvitedById:  inv.InvitedByID,
		CreatedAt:    timestamppb.New(inv.CreatedAt),
		UpdatedAt:    timestamppb.New(inv.UpdatedAt),
	}

	// Optional fields
	if inv.InvitedEmail != nil {
		proto.InvitedEmail = inv.InvitedEmail
	}
	if inv.InvitedUserID != nil {
		proto.InvitedUserId = inv.InvitedUserID
	}
	if inv.InviteCode != nil {
		proto.InviteCode = inv.InviteCode
		// Optionally generate full URL
		// proto.InviteUrl = &("https://vondi.rs/invite/" + *inv.InviteCode)
	}
	if inv.ExpiresAt != nil {
		proto.ExpiresAt = timestamppb.New(*inv.ExpiresAt)
	}
	if inv.MaxUses != nil {
		proto.MaxUses = inv.MaxUses
	}
	if inv.Comment != "" {
		proto.Comment = &inv.Comment
	}
	if inv.AcceptedAt != nil {
		proto.AcceptedAt = timestamppb.New(*inv.AcceptedAt)
	}
	if inv.DeclinedAt != nil {
		proto.DeclinedAt = timestamppb.New(*inv.DeclinedAt)
	}

	return proto
}

// domainTypeToProtoType converts domain invitation type to proto
func domainTypeToProtoType(t domain.StorefrontInvitationType) listingspb.StorefrontInvitationType {
	switch t {
	case domain.InvitationTypeEmail:
		return listingspb.StorefrontInvitationType_STOREFRONT_INVITATION_TYPE_EMAIL
	case domain.InvitationTypeLink:
		return listingspb.StorefrontInvitationType_STOREFRONT_INVITATION_TYPE_LINK
	default:
		return listingspb.StorefrontInvitationType_STOREFRONT_INVITATION_TYPE_UNSPECIFIED
	}
}

// domainStatusToProtoStatus converts domain invitation status to proto
func domainStatusToProtoStatus(s domain.StorefrontInvitationStatus) listingspb.StorefrontInvitationStatus {
	switch s {
	case domain.InvitationStatusPending:
		return listingspb.StorefrontInvitationStatus_STOREFRONT_INVITATION_STATUS_PENDING
	case domain.InvitationStatusAccepted:
		return listingspb.StorefrontInvitationStatus_STOREFRONT_INVITATION_STATUS_ACCEPTED
	case domain.InvitationStatusDeclined:
		return listingspb.StorefrontInvitationStatus_STOREFRONT_INVITATION_STATUS_DECLINED
	case domain.InvitationStatusExpired:
		return listingspb.StorefrontInvitationStatus_STOREFRONT_INVITATION_STATUS_EXPIRED
	case domain.InvitationStatusRevoked:
		return listingspb.StorefrontInvitationStatus_STOREFRONT_INVITATION_STATUS_REVOKED
	default:
		return listingspb.StorefrontInvitationStatus_STOREFRONT_INVITATION_STATUS_UNSPECIFIED
	}
}

// protoStatusToDomainStatus converts proto invitation status to domain
func protoStatusToDomainStatus(s listingspb.StorefrontInvitationStatus) domain.StorefrontInvitationStatus {
	switch s {
	case listingspb.StorefrontInvitationStatus_STOREFRONT_INVITATION_STATUS_PENDING:
		return domain.InvitationStatusPending
	case listingspb.StorefrontInvitationStatus_STOREFRONT_INVITATION_STATUS_ACCEPTED:
		return domain.InvitationStatusAccepted
	case listingspb.StorefrontInvitationStatus_STOREFRONT_INVITATION_STATUS_DECLINED:
		return domain.InvitationStatusDeclined
	case listingspb.StorefrontInvitationStatus_STOREFRONT_INVITATION_STATUS_EXPIRED:
		return domain.InvitationStatusExpired
	case listingspb.StorefrontInvitationStatus_STOREFRONT_INVITATION_STATUS_REVOKED:
		return domain.InvitationStatusRevoked
	default:
		return domain.InvitationStatusPending
	}
}

// stringPtrFromOptional converts optional string to pointer
func stringPtrFromOptional(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
