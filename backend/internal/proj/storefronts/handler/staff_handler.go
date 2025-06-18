package handler

import (
	"backend/internal/domain/models"
	"backend/internal/proj/storefronts/service"
	"backend/internal/storage/postgres"
	"strconv"

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
// @Success 200 {object} SuccessResponse "Staff added"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Insufficient permissions or staff limit reached"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{id}/staff [post]
func (h *StorefrontHandler) AddStaff(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	
	storefrontID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid storefront ID",
		})
	}

	var req AddStaffRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid request body",
		})
	}

	// Валидация роли
	validRoles := map[models.StaffRole]bool{
		models.StaffRoleManager:   true,
		models.StaffRoleSupport:   true,
		models.StaffRoleModerator: true,
	}
	if !validRoles[req.Role] {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid staff role",
		})
	}

	err = h.service.AddStaff(c.Context(), userID, storefrontID, req.UserID, req.Role)
	if err != nil {
		if err == service.ErrInsufficientPermissions {
			return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{
				Error: "Insufficient permissions",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Failed to add staff",
		})
	}

	return c.JSON(SuccessResponse{
		Message: "Staff member added successfully",
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
// @Param permissions body models.JSONB true "Permissions map"
// @Success 200 {object} SuccessResponse "Permissions updated"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Insufficient permissions"
// @Failure 404 {object} ErrorResponse "Staff not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{id}/staff/{staffId}/permissions [put]
func (h *StorefrontHandler) UpdateStaffPermissions(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	
	staffID, err := strconv.Atoi(c.Params("staffId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid staff ID",
		})
	}

	var permissions models.JSONB
	if err := c.BodyParser(&permissions); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid request body",
		})
	}

	err = h.service.UpdateStaffPermissions(c.Context(), userID, staffID, permissions)
	if err != nil {
		switch err {
		case service.ErrInsufficientPermissions:
			return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{
				Error: "Insufficient permissions",
			})
		case postgres.ErrNotFound:
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error: "Staff member not found",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error: "Failed to update permissions",
			})
		}
	}

	return c.JSON(SuccessResponse{
		Message: "Permissions updated successfully",
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
// @Success 200 {object} SuccessResponse "Staff removed"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Insufficient permissions"
// @Failure 404 {object} ErrorResponse "Staff not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/storefronts/{id}/staff/{userId} [delete]
func (h *StorefrontHandler) RemoveStaff(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	
	storefrontID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid storefront ID",
		})
	}

	staffUserID, err := strconv.Atoi(c.Params("userId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid user ID",
		})
	}

	err = h.service.RemoveStaff(c.Context(), userID, storefrontID, staffUserID)
	if err != nil {
		switch err {
		case service.ErrInsufficientPermissions:
			return c.Status(fiber.StatusForbidden).JSON(ErrorResponse{
				Error: "Insufficient permissions",
			})
		case postgres.ErrNotFound:
			return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
				Error: "Staff member not found",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error: "Failed to remove staff",
			})
		}
	}

	return c.JSON(SuccessResponse{
		Message: "Staff member removed successfully",
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
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/v1/storefronts/{id}/staff [get]
func (h *StorefrontHandler) GetStaff(c *fiber.Ctx) error {
	storefrontID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error: "Invalid storefront ID",
		})
	}

	staff, err := h.service.GetStaff(c.Context(), storefrontID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error: "Failed to get staff",
		})
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