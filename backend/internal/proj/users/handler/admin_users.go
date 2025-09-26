// backend/internal/proj/users/handler/admin_users.go
package handler

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"backend/internal/domain/models"
	"backend/internal/logger"
	"backend/pkg/utils"
)

// GetAllAdmins returns list of all administrators
// @Summary Get all administrators
// @Description Returns list of all system administrators
// @Tags admin-management
// @Accept json
// @Produce json
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]models.AdminUser} "List of administrators"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "admin.admins.error.unauthorized"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "admin.admins.error.fetch_failed"
// @Security BearerAuth
// @Router /api/v1/admin/admins [get]
func (h *UserHandler) GetAllAdmins(c *fiber.Ctx) error {
	ctx := context.Background()

	userID, ok := c.Locals("user_id").(int)
	if !ok || userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "admin.admins.error.unauthorized")
	}

	// Получаем список администраторов
	admins, err := h.userService.GetAllAdmins(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Error getting admin users")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "admin.admins.error.fetch_failed")
	}

	return utils.SuccessResponse(c, admins)
}

// AddAdmin adds a new administrator
// @Summary Add new administrator
// @Description Adds a new administrator by email
// @Tags admin-management
// @Accept json
// @Produce json
// @Param admin body models.AdminUser true "Administrator data"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=models.AdminUser} "Created administrator"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "admin.admins.error.invalid_format or admin.admins.error.email_required"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "admin.admins.error.unauthorized"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "admin.admins.error.add_failed"
// @Security BearerAuth
// @Router /api/v1/admin/admins [post]
func (h *UserHandler) AddAdmin(c *fiber.Ctx) error {
	ctx := context.Background()

	userID, ok := c.Locals("user_id").(int)
	if !ok || userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "admin.admins.error.unauthorized")
	}

	// Получаем данные из запроса
	admin := &models.AdminUser{}
	if err := c.BodyParser(admin); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "admin.admins.error.invalid_format")
	}

	// Проверяем email
	if admin.Email == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "admin.admins.error.email_required")
	}

	// Устанавливаем пользователя, который создает администратора
	admin.CreatedBy = &userID

	// Добавляем администратора
	err := h.userService.AddAdmin(ctx, admin)
	if err != nil {
		logger.Error().Err(err).Msg("Error adding admin user")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "admin.admins.error.add_failed")
	}

	return utils.SuccessResponse(c, admin)
}

// RemoveAdmin removes an administrator
// @Summary Remove administrator
// @Description Removes administrator privileges from user by email
// @Tags admin-management
// @Accept json
// @Produce json
// @Param email path string true "Administrator email"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=AdminMessageResponse} "Administrator removed"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "admin.admins.error.email_required"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "admin.admins.error.unauthorized"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "admin.admins.error.remove_failed"
// @Security BearerAuth
// @Router /api/v1/admin/admins/{email} [delete]
func (h *UserHandler) RemoveAdmin(c *fiber.Ctx) error {
	ctx := context.Background()

	userID, ok := c.Locals("user_id").(int)
	if !ok || userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "admin.admins.error.unauthorized")
	}

	// Получаем email из параметров пути
	email := c.Params("email")
	if email == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "admin.admins.error.email_required")
	}

	// Удаляем администратора
	err := h.userService.RemoveAdmin(ctx, email)
	if err != nil {
		logger.Error().Err(err).Msg("Error removing admin user")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "admin.admins.error.remove_failed")
	}

	response := AdminMessageResponse{
		Message: "admin.admins.success.removed",
	}
	return utils.SuccessResponse(c, response)
}

// IsAdmin checks if user is an administrator
// @Summary Check administrator status
// @Description Checks if user with specified email is an administrator
// @Tags admin-management
// @Accept json
// @Produce json
// @Param email path string true "User email"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=AdminAdminsResponse} "Admin status"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "admin.admins.error.email_required"
// @Failure 401 {object} backend_pkg_utils.ErrorResponseSwag "admin.admins.error.unauthorized"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "admin.admins.error.check_failed"
// @Security BearerAuth
// @Router /api/v1/admin/admins/check/{email} [get]
func (h *UserHandler) IsAdmin(c *fiber.Ctx) error {
	ctx := context.Background()

	userID, ok := c.Locals("user_id").(int)
	if !ok || userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "admin.admins.error.unauthorized")
	}

	// Получаем email из параметров пути
	email := c.Params("email")
	if email == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "admin.admins.error.email_required")
	}

	// Проверяем, является ли пользователь администратором
	isAdmin, err := h.userService.IsUserAdmin(ctx, email)
	if err != nil {
		logger.Error().Err(err).Msg("Error checking admin status")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "admin.admins.error.check_failed")
	}

	response := AdminAdminsResponse{
		Email:   email,
		IsAdmin: isAdmin,
	}
	return utils.SuccessResponse(c, response)
}

// IsAdminPublic checks if user is an administrator (public method)
// @Summary Public administrator check
// @Description Checks if user with specified email is an administrator (no authorization required)
// @Tags public
// @Accept json
// @Produce json
// @Param email path string true "User email"
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=AdminAdminsResponse} "Admin status"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "admin.admins.error.email_required"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "admin.admins.error.check_failed"
// @Router /api/v1/admin-check/{email} [get]
func (h *UserHandler) IsAdminPublic(c *fiber.Ctx) error {
	ctx := context.Background()

	// Получаем email из параметров пути
	email := c.Params("email")
	if email == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "admin.admins.error.email_required")
	}

	// Проверяем, является ли пользователь администратором
	isAdmin, err := h.userService.IsUserAdmin(ctx, email)
	if err != nil {
		logger.Error().Err(err).Msg("Error checking admin status")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "admin.admins.error.check_failed")
	}

	response := AdminAdminsResponse{
		Email:   email,
		IsAdmin: isAdmin,
	}
	return utils.SuccessResponse(c, response)
}
