// backend/internal/proj/users/handler/admin_users.go
package handler

import (
	"backend/internal/domain/models"
	"backend/pkg/utils"
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
)

// GetAllAdmins возвращает список всех администраторов
// @Summary Получить список администраторов
// @Description Получает список всех администраторов системы
// @Tags Admins
// @Accept json
// @Produce json
// @Success 200 {array} models.AdminUser
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/admin/admins [get]
func (h *UserHandler) GetAllAdmins(c *fiber.Ctx) error {
	ctx := context.Background()

	userID, ok := c.Locals("user_id").(int)
	if !ok || userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	// Получаем список администраторов
	admins, err := h.userService.GetAllAdmins(ctx)
	if err != nil {
		log.Printf("Error getting admin users: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка при получении списка администраторов")
	}

	return c.JSON(admins)
}

// AddAdmin добавляет нового администратора
// @Summary Добавить администратора
// @Description Добавляет нового администратора по email
// @Tags Admins
// @Accept json
// @Produce json
// @Param admin body models.AdminUser true "Данные администратора"
// @Success 200 {object} models.AdminUser
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/admin/admins [post]
func (h *UserHandler) AddAdmin(c *fiber.Ctx) error {
	ctx := context.Background()

	userID, ok := c.Locals("user_id").(int)
	if !ok || userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	// Получаем данные из запроса
	admin := &models.AdminUser{}
	if err := c.BodyParser(admin); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный формат данных")
	}

	// Проверяем email
	if admin.Email == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Email обязателен")
	}

	// Устанавливаем пользователя, который создает администратора
	admin.CreatedBy = &userID

	// Добавляем администратора
	err := h.userService.AddAdmin(ctx, admin)
	if err != nil {
		log.Printf("Error adding admin user: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка при добавлении администратора")
	}

	return c.JSON(admin)
}

// RemoveAdmin удаляет администратора
// @Summary Удалить администратора
// @Description Удаляет администратора по email
// @Tags Admins
// @Accept json
// @Produce json
// @Param email path string true "Email администратора"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/admin/admins/{email} [delete]
func (h *UserHandler) RemoveAdmin(c *fiber.Ctx) error {
	ctx := context.Background()

	userID, ok := c.Locals("user_id").(int)
	if !ok || userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	// Получаем email из параметров пути
	email := c.Params("email")
	if email == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Email обязателен")
	}

	// Удаляем администратора
	err := h.userService.RemoveAdmin(ctx, email)
	if err != nil {
		log.Printf("Error removing admin user: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка при удалении администратора")
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Администратор успешно удален",
	})
}

// IsAdmin проверяет, является ли пользователь администратором
// @Summary Проверить статус администратора
// @Description Проверяет, является ли пользователь с указанным email администратором
// @Tags Admins
// @Accept json
// @Produce json
// @Param email path string true "Email пользователя"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/admin/admins/check/{email} [get]
func (h *UserHandler) IsAdmin(c *fiber.Ctx) error {
	ctx := context.Background()

	userID, ok := c.Locals("user_id").(int)
	if !ok || userID == 0 {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	// Получаем email из параметров пути
	email := c.Params("email")
	if email == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Email обязателен")
	}

	// Проверяем, является ли пользователь администратором
	isAdmin, err := h.userService.IsUserAdmin(ctx, email)
	if err != nil {
		log.Printf("Error checking admin status: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка при проверке статуса администратора")
	}

	return c.JSON(fiber.Map{
		"email":    email,
		"is_admin": isAdmin,
	})
}
