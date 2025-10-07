package service

import (
	"context"
	"fmt"

	"backend/internal/domain/models"
	"backend/internal/logger"
	"backend/internal/storage/postgres"
)

// CategoryRepository defines interface for category operations
type CategoryRepository interface {
	GetCategoryByID(ctx context.Context, id int) (*models.MarketplaceCategory, error)
}

// CategoryProposalService handles category proposal operations
type CategoryProposalService struct {
	proposalsRepo    postgres.CategoryProposalsRepositoryInterface
	categoriesRepo   CategoryRepository
	aiCategoryMapper *AICategoryMapper
}

// NewCategoryProposalService creates a new category proposal service
func NewCategoryProposalService(
	proposalsRepo postgres.CategoryProposalsRepositoryInterface,
	categoriesRepo CategoryRepository,
	aiCategoryMapper *AICategoryMapper,
) *CategoryProposalService {
	return &CategoryProposalService{
		proposalsRepo:    proposalsRepo,
		categoriesRepo:   categoriesRepo,
		aiCategoryMapper: aiCategoryMapper,
	}
}

// CreateProposal creates a new category proposal
func (s *CategoryProposalService) CreateProposal(ctx context.Context, userID int, req *models.CreateCategoryProposalRequest) (*models.CategoryProposal, error) {
	proposal := &models.CategoryProposal{
		ProposedByUserID:       userID,
		StorefrontID:           req.StorefrontID,
		Name:                   req.Name,
		NameTranslations:       req.NameTranslations,
		ParentCategoryID:       req.ParentCategoryID,
		Description:            req.Description,
		Reasoning:              req.Reasoning,
		ExpectedProducts:       req.ExpectedProducts,
		ExternalCategorySource: req.ExternalCategorySource,
		SimilarCategories:      req.SimilarCategories,
		Tags:                   req.Tags,
		Status:                 models.CategoryProposalStatusPending,
	}

	// Validate parent category exists if specified
	if proposal.ParentCategoryID != nil {
		_, err := s.categoriesRepo.GetCategoryByID(ctx, *proposal.ParentCategoryID)
		if err != nil {
			return nil, fmt.Errorf("parent category not found: %w", err)
		}
	}

	err := s.proposalsRepo.Create(ctx, proposal)
	if err != nil {
		return nil, fmt.Errorf("failed to create proposal: %w", err)
	}

	logger.Info().
		Int("proposal_id", proposal.ID).
		Int("user_id", userID).
		Str("name", proposal.Name).
		Msg("Category proposal created")

	return proposal, nil
}

// GetProposal retrieves a proposal by ID
func (s *CategoryProposalService) GetProposal(ctx context.Context, id int) (*models.CategoryProposal, error) {
	proposal, err := s.proposalsRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get proposal: %w", err)
	}
	return proposal, nil
}

// ListProposals retrieves proposals with filters and pagination
func (s *CategoryProposalService) ListProposals(ctx context.Context, filter *postgres.CategoryProposalFilter) (*models.CategoryProposalListResponse, error) {
	proposals, total, err := s.proposalsRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list proposals: %w", err)
	}

	totalPages := (total + filter.Limit - 1) / filter.Limit
	page := (filter.Offset / filter.Limit) + 1

	response := &models.CategoryProposalListResponse{
		Proposals:  make([]models.CategoryProposal, 0, len(proposals)),
		Total:      total,
		Page:       page,
		PageSize:   filter.Limit,
		TotalPages: totalPages,
	}

	for _, p := range proposals {
		response.Proposals = append(response.Proposals, *p)
	}

	return response, nil
}

// UpdateProposal updates an existing proposal
func (s *CategoryProposalService) UpdateProposal(ctx context.Context, id int, req *models.UpdateCategoryProposalRequest) (*models.CategoryProposal, error) {
	// Get existing proposal
	proposal, err := s.proposalsRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get proposal: %w", err)
	}

	// Check if proposal is still pending
	if proposal.Status != models.CategoryProposalStatusPending {
		return nil, fmt.Errorf("cannot update non-pending proposal")
	}

	// Update fields if provided
	if req.Name != nil {
		proposal.Name = *req.Name
	}
	if req.NameTranslations != nil {
		proposal.NameTranslations = req.NameTranslations
	}
	if req.ParentCategoryID != nil {
		// Validate parent category exists
		_, err := s.categoriesRepo.GetCategoryByID(ctx, *req.ParentCategoryID)
		if err != nil {
			return nil, fmt.Errorf("parent category not found: %w", err)
		}
		proposal.ParentCategoryID = req.ParentCategoryID
	}
	if req.Description != nil {
		proposal.Description = req.Description
	}
	if req.Tags != nil {
		proposal.Tags = req.Tags
	}

	err = s.proposalsRepo.Update(ctx, proposal)
	if err != nil {
		return nil, fmt.Errorf("failed to update proposal: %w", err)
	}

	logger.Info().
		Int("proposal_id", id).
		Str("name", proposal.Name).
		Msg("Category proposal updated")

	return proposal, nil
}

// ApproveProposal approves a proposal and optionally creates the category
func (s *CategoryProposalService) ApproveProposal(ctx context.Context, id int, reviewedByUserID int, createCategory bool) (*models.CategoryProposal, *models.MarketplaceCategory, error) {
	// Get proposal
	proposal, err := s.proposalsRepo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get proposal: %w", err)
	}

	// Check if proposal is pending
	if proposal.Status != models.CategoryProposalStatusPending {
		return nil, nil, fmt.Errorf("proposal is not pending")
	}

	// Approve proposal
	err = s.proposalsRepo.Approve(ctx, id, reviewedByUserID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to approve proposal: %w", err)
	}

	// Get updated proposal
	proposal, err = s.proposalsRepo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get updated proposal: %w", err)
	}

	logger.Info().
		Int("proposal_id", id).
		Int("reviewed_by", reviewedByUserID).
		Bool("create_category", createCategory).
		Msg("Category proposal approved")

	// Create category if requested
	var newCategory *models.MarketplaceCategory
	if createCategory {
		// TODO: Implement category creation when admin category management API is ready
		// For now, just return proposal without creating category
		logger.Warn().
			Int("proposal_id", id).
			Msg("Category creation requested but not yet implemented - manual creation required")
	}

	return proposal, newCategory, nil
}

// RejectProposal rejects a proposal
func (s *CategoryProposalService) RejectProposal(ctx context.Context, id int, reviewedByUserID int, reason *string) (*models.CategoryProposal, error) {
	// Get proposal
	proposal, err := s.proposalsRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get proposal: %w", err)
	}

	// Check if proposal is pending
	if proposal.Status != models.CategoryProposalStatusPending {
		return nil, fmt.Errorf("proposal is not pending")
	}

	// Reject proposal
	err = s.proposalsRepo.Reject(ctx, id, reviewedByUserID, reason)
	if err != nil {
		return nil, fmt.Errorf("failed to reject proposal: %w", err)
	}

	// Get updated proposal
	proposal, err = s.proposalsRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated proposal: %w", err)
	}

	logger.Info().
		Int("proposal_id", id).
		Int("reviewed_by", reviewedByUserID).
		Msg("Category proposal rejected")

	return proposal, nil
}

// DeleteProposal deletes a proposal
func (s *CategoryProposalService) DeleteProposal(ctx context.Context, id int) error {
	err := s.proposalsRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete proposal: %w", err)
	}

	logger.Info().
		Int("proposal_id", id).
		Msg("Category proposal deleted")

	return nil
}

// GetPendingCount returns count of pending proposals
func (s *CategoryProposalService) GetPendingCount(ctx context.Context, storefrontID *int) (int, error) {
	count, err := s.proposalsRepo.GetPendingCount(ctx, storefrontID)
	if err != nil {
		return 0, fmt.Errorf("failed to get pending count: %w", err)
	}
	return count, nil
}
