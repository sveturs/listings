// backend/internal/proj/users/handler/admin.go
package handler

import (
	"backend/internal/domain/models"
	"backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"log"
	"strconv"
)

// GetAllUsers возвращает список всех пользователей с пагинацией
// GetAllUsers возвращает список всех пользователей с пагинацией
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	log.Printf("GetAllUsers handler called with path: %s", c.Path())

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
		log.Printf("GetAllUsers: error getting users: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка получения списка пользователей",
		})
	}

	// Возвращаем данные напрямую
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":  users,
		"total": total,
		"page":  page,
		"limit": limit,
		"pages": (total + limit - 1) / limit, // Округление вверх
	})
}

// GetUserByID возвращает информацию о пользователе по ID
func (h *UserHandler) GetUserByIDAdmin(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный ID пользователя")
	}

	profile, err := h.userService.GetUserProfile(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Пользователь не найден")
	}

	return utils.SuccessResponse(c, profile)
}

// UpdateUserAdmin обновляет информацию о пользователе (админ)
func (h *UserHandler) UpdateUserAdmin(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный ID пользователя")
	}

	var update models.UserProfileUpdate
	if err := c.BodyParser(&update); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный формат данных")
	}

	if err := update.Validate(); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	err = h.userService.UpdateUserProfile(c.Context(), id, &update)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка обновления профиля пользователя")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Профиль пользователя успешно обновлен",
	})
}

// UpdateUserStatus обновляет статус пользователя (блокировка/разблокировка)
func (h *UserHandler) UpdateUserStatus(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный ID пользователя")
	}

	var data struct {
		Status string `json:"status"`
	}

	if err := c.BodyParser(&data); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный формат данных")
	}

	// Проверяем допустимость статуса
	if data.Status != "active" && data.Status != "blocked" && data.Status != "pending" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Недопустимый статус. Допустимые значения: active, blocked, pending")
	}

	err = h.userService.UpdateUserStatus(c.Context(), id, data.Status)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка обновления статуса пользователя")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Статус пользователя успешно обновлен",
	})
}

// DeleteUser удаляет пользователя
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный ID пользователя")
	}

	err = h.userService.DeleteUser(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка удаления пользователя")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Пользователь успешно удален",
	})
}
