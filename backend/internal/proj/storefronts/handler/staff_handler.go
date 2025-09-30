package handler

import (
	"errors"
	"strconv"

	"backend/internal/domain/models"
	"backend/internal/proj/storefronts/service"
	"backend/internal/storage/postgres"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// AddStaff добавляет сотрудника в витрину
// @Summary Add staff member
// @Description Adds a new staff member to the storefront
// @Tags storefronts,staff
// @Accept json
// @Produce json
// @Param id path int true "Storefront ID"
// @Param staff body AddStaffRequest true "Staff data"
// @Success 200 {object} map[string]string "Staff added"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Bad request"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "Unauthorized"
// @Failure 403 {object} backend_pkg_utils.ErrorResponseSwag "Insufficient permissions or staff limit reached"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{id}/staff [post]
func (h *StorefrontHandler) AddStaff(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	storefrontID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.invalid_id")
	}

	var req AddStaffRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.invalid_request_body")
	}

	// Валидация роли
	validRoles := map[models.StaffRole]bool{
		models.StaffRoleManager:   true,
		models.StaffRoleSupport:   true,
		models.StaffRoleModerator: true,
	}
	if !validRoles[req.Role] {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.invalid_staff_role")
	}

	err = h.service.AddStaff(c.Context(), userID, storefrontID, req.UserID, req.Role)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInsufficientPermissions):
			return utils.ErrorResponse(c, fiber.StatusForbidden, "storefronts.error.insufficient_permissions")
		case errors.Is(err, service.ErrStaffLimitReached):
			return utils.ErrorResponse(c, fiber.StatusForbidden, "storefronts.error.staff_limit_reached")
		default:
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefronts.error.add_staff_failed")
		}
	}

	return c.JSON(fiber.Map{
		"message": "Staff member added successfully",
	})
}

// UpdateStaffPermissions обновляет права сотрудника
// @Summary Update staff permissions
// @Description Updates permissions for a staff member
// @Tags storefronts,staff
// @Accept json
// @Produce json
// @Param id path int true "Storefront ID"
// @Param staffId path int true "Staff ID"
// @Param permissions body backend_internal_domain_models.JSONB true "Permissions map"
// @Success 200 {object} map[string]string "Permissions updated"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Bad request"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "Unauthorized"
// @Failure 403 {object} backend_pkg_utils.ErrorResponseSwag "Insufficient permissions"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "Staff not found"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{id}/staff/{staffId}/permissions [put]
func (h *StorefrontHandler) UpdateStaffPermissions(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	staffID, err := strconv.Atoi(c.Params("staffId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.invalid_staff_id")
	}

	var permissions models.JSONB
	if err := c.BodyParser(&permissions); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.invalid_request_body")
	}

	err = h.service.UpdateStaffPermissions(c.Context(), userID, staffID, permissions)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInsufficientPermissions):
			return utils.ErrorResponse(c, fiber.StatusForbidden, "storefronts.error.insufficient_permissions")
		case errors.Is(err, postgres.ErrNotFound):
			return utils.ErrorResponse(c, fiber.StatusNotFound, "storefronts.error.staff_not_found")
		default:
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefronts.error.update_permissions_failed")
		}
	}

	return c.JSON(fiber.Map{
		"message": "Permissions updated successfully",
	})
}

// RemoveStaff удаляет сотрудника из витрины
// @Summary Remove staff member
// @Description Removes a staff member from the storefront
// @Tags storefronts,staff
// @Accept json
// @Produce json
// @Param id path int true "Storefront ID"
// @Param userId path int true "User ID of staff member"
// @Success 200 {object} map[string]string "Staff removed"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Bad request"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "Unauthorized"
// @Failure 403 {object} backend_pkg_utils.ErrorResponseSwag "Insufficient permissions"
// @Failure 404 {object} backend_pkg_utils.ErrorResponseSwag "Staff not found"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{id}/staff/{userId} [delete]
func (h *StorefrontHandler) RemoveStaff(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	storefrontID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.invalid_id")
	}

	staffUserID, err := strconv.Atoi(c.Params("userId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.invalid_user_id")
	}

	err = h.service.RemoveStaff(c.Context(), userID, storefrontID, staffUserID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInsufficientPermissions):
			return utils.ErrorResponse(c, fiber.StatusForbidden, "storefronts.error.insufficient_permissions")
		case errors.Is(err, postgres.ErrNotFound):
			return utils.ErrorResponse(c, fiber.StatusNotFound, "storefronts.error.staff_not_found")
		default:
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefronts.error.remove_staff_failed")
		}
	}

	return c.JSON(fiber.Map{
		"message": "Staff member removed successfully",
	})
}

// GetStaff получает список сотрудников витрины
// @Summary Get storefront staff
// @Description Returns list of staff members with their roles and permissions
// @Tags storefronts,staff
// @Accept json
// @Produce json
// @Param id path int true "Storefront ID"
// @Success 200 {object} StaffListResponse "Staff list"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "Bad request"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/storefronts/{id}/staff [get]
func (h *StorefrontHandler) GetStaff(c *fiber.Ctx) error {
	storefrontID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "storefronts.error.invalid_id")
	}

	staff, err := h.service.GetStaff(c.Context(), storefrontID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "storefronts.error.get_staff_failed")
	}

	return c.JSON(StaffListResponse{
		Staff: staff,
		Total: len(staff),
	})
}

// Request/Response types for staff

type AddStaffRequest struct {
	UserID int              `json:"user_id" validate:"required"`
	Role   models.StaffRole `json:"role" validate:"required"`
}

type StaffListResponse struct {
	Staff []*models.StorefrontStaff `json:"staff"`
	Total int                       `json:"total"`
}
