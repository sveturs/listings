package handler

import (
	"strconv"

	"backend/internal/domain/models"
	"backend/internal/logger"
	"backend/internal/proj/storefronts/service"
	"backend/internal/storage/postgres"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"
)

// CategoryProposalHandler handles category proposal endpoints
type CategoryProposalHandler struct {
	proposalService *service.CategoryProposalService
}

// NewCategoryProposalHandler creates a new category proposal handler
func NewCategoryProposalHandler(proposalService *service.CategoryProposalService) *CategoryProposalHandler {
	return &CategoryProposalHandler{
		proposalService: proposalService,
	}
}

// CreateProposal creates a new category proposal
// @Summary Create category proposal
// @Description Create a new category proposal (requires authentication)
// @Tags admin,categories,proposals
// @Accept json
// @Produce json
// @Param request body models.CreateCategoryProposalRequest true "Proposal data"
// @Success 201 {object} models.CategoryProposal
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /api/v1/admin/category-proposals [post]
func (h *CategoryProposalHandler) CreateProposal(c *fiber.Ctx) error {
	userID, ok := authMiddleware.GetUserID(c)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "categories.unauthorized")
	}

	var req models.CreateCategoryProposalRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "categories.invalid_request")
	}

	proposal, err := h.proposalService.CreateProposal(c.Context(), userID, &req)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create category proposal")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "categories.create_failed")
	}

	return c.Status(fiber.StatusCreated).JSON(proposal)
}

// GetProposal retrieves a proposal by ID
// @Summary Get category proposal
// @Description Get a category proposal by ID (admin only)
// @Tags admin,categories,proposals
// @Produce json
// @Param id path int true "Proposal ID"
// @Success 200 {object} models.CategoryProposal
// @Failure 404 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /api/v1/admin/category-proposals/{id} [get]
func (h *CategoryProposalHandler) GetProposal(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "categories.invalid_id")
	}

	proposal, err := h.proposalService.GetProposal(c.Context(), id)
	if err != nil {
		logger.Error().Err(err).Int("id", id).Msg("Failed to get category proposal")
		return utils.ErrorResponse(c, fiber.StatusNotFound, "categories.not_found")
	}

	return c.JSON(proposal)
}

// ListProposals lists category proposals with filters
// @Summary List category proposals
// @Description List category proposals with pagination and filters (admin only)
// @Tags admin,categories,proposals
// @Produce json
// @Param status query string false "Filter by status (pending/approved/rejected)"
// @Param storefront_id query int false "Filter by storefront ID"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 20, max: 100)"
// @Success 200 {object} models.CategoryProposalListResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /api/v1/admin/category-proposals [get]
func (h *CategoryProposalHandler) ListProposals(c *fiber.Ctx) error {
	// Parse filters
	filter := &postgres.CategoryProposalFilter{
		Limit:  20,
		Offset: 0,
	}

	// Status filter
	if statusStr := c.Query("status"); statusStr != "" {
		status := models.CategoryProposalStatus(statusStr)
		filter.Status = &status
	}

	// Storefront filter
	if storefrontIDStr := c.Query("storefront_id"); storefrontIDStr != "" {
		storefrontID, err := strconv.Atoi(storefrontIDStr)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "categories.invalid_storefront_id")
		}
		filter.StorefrontID = &storefrontID
	}

	// Pagination
	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	pageSize := 20
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	filter.Limit = pageSize
	filter.Offset = (page - 1) * pageSize

	response, err := h.proposalService.ListProposals(c.Context(), filter)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to list category proposals")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "categories.list_failed")
	}

	return c.JSON(response)
}

// UpdateProposal updates a category proposal
// @Summary Update category proposal
// @Description Update a pending category proposal (admin only)
// @Tags admin,categories,proposals
// @Accept json
// @Produce json
// @Param id path int true "Proposal ID"
// @Param request body models.UpdateCategoryProposalRequest true "Update data"
// @Success 200 {object} models.CategoryProposal
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /api/v1/admin/category-proposals/{id} [put]
func (h *CategoryProposalHandler) UpdateProposal(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "categories.invalid_id")
	}

	var req models.UpdateCategoryProposalRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "categories.invalid_request")
	}

	proposal, err := h.proposalService.UpdateProposal(c.Context(), id, &req)
	if err != nil {
		logger.Error().Err(err).Int("id", id).Msg("Failed to update category proposal")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "categories.update_failed")
	}

	return c.JSON(proposal)
}

// ApproveProposal approves a category proposal
// @Summary Approve category proposal
// @Description Approve a pending category proposal and optionally create the category (admin only)
// @Tags admin,categories,proposals
// @Accept json
// @Produce json
// @Param id path int true "Proposal ID"
// @Param request body models.CategoryProposalApproveRequest true "Approve options"
// @Success 200 {object} utils.SuccessResponseSwag
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /api/v1/admin/category-proposals/{id}/approve [post]
func (h *CategoryProposalHandler) ApproveProposal(c *fiber.Ctx) error {
	userID, ok := authMiddleware.GetUserID(c)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "categories.unauthorized")
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "categories.invalid_id")
	}

	var req models.CategoryProposalApproveRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "categories.invalid_request")
	}

	proposal, category, err := h.proposalService.ApproveProposal(c.Context(), id, userID, req.CreateCategory)
	if err != nil {
		logger.Error().Err(err).Int("id", id).Msg("Failed to approve category proposal")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "categories.approve_failed")
	}

	return c.JSON(fiber.Map{
		"proposal": proposal,
		"category": category,
	})
}

// RejectProposal rejects a category proposal
// @Summary Reject category proposal
// @Description Reject a pending category proposal (admin only)
// @Tags admin,categories,proposals
// @Accept json
// @Produce json
// @Param id path int true "Proposal ID"
// @Param request body models.CategoryProposalRejectRequest true "Reject reason"
// @Success 200 {object} models.CategoryProposal
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /api/v1/admin/category-proposals/{id}/reject [post]
func (h *CategoryProposalHandler) RejectProposal(c *fiber.Ctx) error {
	userID, ok := authMiddleware.GetUserID(c)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "categories.unauthorized")
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "categories.invalid_id")
	}

	var req models.CategoryProposalRejectRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "categories.invalid_request")
	}

	proposal, err := h.proposalService.RejectProposal(c.Context(), id, userID, req.Reason)
	if err != nil {
		logger.Error().Err(err).Int("id", id).Msg("Failed to reject category proposal")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "categories.reject_failed")
	}

	return c.JSON(proposal)
}

// DeleteProposal deletes a category proposal
// @Summary Delete category proposal
// @Description Delete a category proposal (admin only)
// @Tags admin,categories,proposals
// @Produce json
// @Param id path int true "Proposal ID"
// @Success 204 "No content"
// @Failure 404 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /api/v1/admin/category-proposals/{id} [delete]
func (h *CategoryProposalHandler) DeleteProposal(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "categories.invalid_id")
	}

	err = h.proposalService.DeleteProposal(c.Context(), id)
	if err != nil {
		logger.Error().Err(err).Int("id", id).Msg("Failed to delete category proposal")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "categories.delete_failed")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GetPendingCount returns count of pending proposals
// @Summary Get pending count
// @Description Get count of pending category proposals (admin only)
// @Tags admin,categories,proposals
// @Produce json
// @Param storefront_id query int false "Filter by storefront ID"
// @Success 200 {object} object{count=int}
// @Failure 401 {object} models.ErrorResponse
// @Router /api/v1/admin/category-proposals/pending/count [get]
func (h *CategoryProposalHandler) GetPendingCount(c *fiber.Ctx) error {
	var storefrontID *int
	if storefrontIDStr := c.Query("storefront_id"); storefrontIDStr != "" {
		id, err := strconv.Atoi(storefrontIDStr)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "categories.invalid_storefront_id")
		}
		storefrontID = &id
	}

	count, err := h.proposalService.GetPendingCount(c.Context(), storefrontID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get pending count")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "categories.count_failed")
	}

	return c.JSON(fiber.Map{
		"count": count,
	})
}
