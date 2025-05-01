package handler

import (
	"backend/internal/domain/models"
	globalService "backend/internal/proj/global/service"
	"backend/internal/proj/users/service"
	"backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"strconv" // Добавляем импорт для strconv
	// "encoding/json" // Добавляем импорт для json
)

type UserHandler struct {
	services    globalService.ServicesInterface
	userService service.UserServiceInterface
}

func NewUserHandler(services globalService.ServicesInterface) *UserHandler {
	return &UserHandler{
		services:    services,
		userService: services.User(),
	}
}

// Добавляем методы, связанные с администрированием
// GetAllUsers экспортирован в admin.go
// GetUserByIDAdmin экспортирован в admin.go
// UpdateUserAdmin экспортирован в admin.go
// UpdateUserStatus экспортирован в admin.go
// DeleteUser экспортирован в admin.go
// GetUserBalance экспортирован в admin_balance.go
// GetUserTransactions экспортирован в admin_balance.go

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	profile, err := h.services.User().GetUserProfile(c.Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка получения профиля")
	}

	return utils.SuccessResponse(c, profile)
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	var update models.UserProfileUpdate
	if err := c.BodyParser(&update); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный формат данных")
	}

	err := h.services.User().UpdateUserProfile(c.Context(), userID, &update)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка обновления профиля")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Профиль успешно обновлен",
	})
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var registerData struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := c.BodyParser(&registerData); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный формат данных")
	}

	user := &models.User{
		Name:  registerData.Name,
		Email: registerData.Email,
	}

	err := h.services.User().CreateUser(c.Context(), user)
	if err != nil {
		if err.Error() == "email already exists" {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Email уже используется")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка регистрации пользователя")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"message": "Пользователь успешно зарегистрирован",
		"user":    user,
	})
}

func (h *UserHandler) GetProfileByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid user ID")
	}

	user, err := h.userService.GetUserByID(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "User not found")
	}

	// Получаем только данные пользователя из User и возвращаем их
	return utils.SuccessResponse(c, fiber.Map{
		"id":          user.ID,
		"name":        user.Name,
		"email":       user.Email,
		"picture_url": user.PictureURL,
		"created_at":  user.CreatedAt,
	})
}
