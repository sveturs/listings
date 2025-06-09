// Package handler
// backend/internal/proj/users/handler/admin.go
package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"backend/internal/domain/models"
	"backend/internal/logger"
	"backend/pkg/utils"
)

// GetAllUsers returns paginated list of all users
// @Summary Get all users (Admin)
// @Description Returns paginated list of all users in the system
// @Tags admin-users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10) minimum(1) maximum(100)
// @Success 200 {object} utils.SuccessResponseSwag{data=AdminUserListResponse} "List of users"
// @Failure 500 {object} utils.ErrorResponseSwag "admin.users.error.fetch_failed"
// @Security BearerAuth
// @Router /api/v1/admin/users [get]
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	logger.Debug().Str("path", c.Path()).Msg("GetAllUsers handler called")

	// Получаем параметры пагинации
	page := utils.StringToInt(c.Query("page", "1"), 1)
	limit := utils.StringToInt(c.Query("limit", "10"), 10)

	// Проверяем корректность параметров
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Вычисляем смещение
	offset := (page - 1) * limit

	// Получаем пользователей из базы данных
	users, total, err := h.userService.GetAllUsers(c.Context(), limit, offset)
	if err != nil {
		logger.Error().Err(err).Msg("GetAllUsers: error getting users")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "admin.users.error.fetch_failed")
	}

	// Возвращаем данные напрямую
	response := AdminUserListResponse{
		Data:  users,
		Total: total,
		Page:  page,
		Limit: limit,
		Pages: (total + limit - 1) / limit, // Округление вверх
	}
	return utils.SuccessResponse(c, response)
}

// GetUserByIDAdmin returns user information by ID
// @Summary Get user by ID (Admin)
// @Description Returns detailed user information by user ID
// @Tags admin-users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.UserProfile} "User profile"
// @Failure 400 {object} utils.ErrorResponseSwag "admin.users.error.invalid_user_id"
// @Failure 404 {object} utils.ErrorResponseSwag "admin.users.error.user_not_found"
// @Security BearerAuth
// @Router /api/v1/admin/users/{id} [get]
func (h *UserHandler) GetUserByIDAdmin(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "admin.users.error.invalid_user_id")
	}

	profile, err := h.userService.GetUserProfile(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "admin.users.error.user_not_found")
	}

	return utils.SuccessResponse(c, profile)
}

// UpdateUserAdmin updates user information
// @Summary Update user (Admin)
// @Description Updates user profile information by administrator
// @Tags admin-users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param body body models.UserProfileUpdate true "User update data"
// @Success 200 {object} utils.SuccessResponseSwag{data=AdminMessageResponse} "Profile updated"
// @Failure 400 {object} utils.ErrorResponseSwag "admin.users.error.invalid_user_id or admin.users.error.invalid_format or admin.users.error.validation_failed"
// @Failure 500 {object} utils.ErrorResponseSwag "admin.users.error.update_failed"
// @Security BearerAuth
// @Router /api/v1/admin/users/{id} [put]
func (h *UserHandler) UpdateUserAdmin(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "admin.users.error.invalid_user_id")
	}

	var update models.UserProfileUpdate
	if err := c.BodyParser(&update); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "admin.users.error.invalid_format")
	}

	if err := update.Validate(); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "admin.users.error.validation_failed")
	}

	err = h.userService.UpdateUserProfile(c.Context(), id, &update)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "admin.users.error.update_failed")
	}

	response := AdminMessageResponse{
		Message: "admin.users.success.profile_updated",
	}
	return utils.SuccessResponse(c, response)
}

// UpdateUserStatus updates user status (block/unblock)
// @Summary Update user status (Admin)
// @Description Updates user status (active, blocked, pending)
// @Tags admin-users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param body body AdminStatusUpdateRequest true "Status update (active, blocked, pending)"
// @Success 200 {object} utils.SuccessResponseSwag{data=AdminMessageResponse} "Status updated"
// @Failure 400 {object} utils.ErrorResponseSwag "admin.users.error.invalid_user_id or admin.users.error.invalid_format or admin.users.error.invalid_status"
// @Failure 500 {object} utils.ErrorResponseSwag "admin.users.error.status_update_failed"
// @Security BearerAuth
// @Router /api/v1/admin/users/{id}/status [put]
func (h *UserHandler) UpdateUserStatus(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "admin.users.error.invalid_user_id")
	}

	var data AdminStatusUpdateRequest

	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "admin.users.error.invalid_format")
	}

	// Проверяем допустимость статуса
	if data.Status != "active" && data.Status != "blocked" && data.Status != "pending" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "admin.users.error.invalid_status")
	}

	err = h.userService.UpdateUserStatus(c.Context(), id, data.Status)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "admin.users.error.status_update_failed")
	}

	response := AdminMessageResponse{
		Message: "admin.users.success.status_updated",
	}
	return utils.SuccessResponse(c, response)
}

// DeleteUser deletes a user from the system
// @Summary Delete user (Admin)
// @Description Permanently deletes a user from the system
// @Tags admin-users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=AdminMessageResponse} "User deleted"
// @Failure 400 {object} utils.ErrorResponseSwag "admin.users.error.invalid_user_id"
// @Failure 500 {object} utils.ErrorResponseSwag "admin.users.error.delete_failed"
// @Security BearerAuth
// @Router /api/v1/admin/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "admin.users.error.invalid_user_id")
	}

	err = h.userService.DeleteUser(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "admin.users.error.delete_failed")
	}

	response := AdminMessageResponse{
		Message: "admin.users.success.deleted",
	}
	return utils.SuccessResponse(c, response)
}
