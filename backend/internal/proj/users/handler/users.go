//backend/internal/proj/users/handler/users.go
package handler

import (
    "github.com/gofiber/fiber/v2"
    "backend/internal/services"
//    "backend/internal/types"
    "backend/pkg/utils"
    "backend/internal/domain/models"
)

type UserHandler struct {
    services services.ServicesInterface  // Изменено с *services.ServicesInterface
}

func NewUserHandler(services services.ServicesInterface) *UserHandler {  // Изменено с *services.ServicesInterface
    return &UserHandler{
        services: services,
    }
}

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