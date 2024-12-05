package handlers

import (
    "github.com/gofiber/fiber/v2"
    "backend/internal/services"
    "backend/internal/types"
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
    sessionData := c.Locals("user").(*types.SessionData)
    if sessionData == nil {
        return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Необходима авторизация")
    }

    user, err := h.services.User().GetUserByID(c.Context(), sessionData.UserID)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка получения профиля")
    }

    return utils.SuccessResponse(c, user)
}

func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
    sessionData := c.Locals("user").(*types.SessionData)
    if sessionData == nil {
        return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Необходима авторизация")
    }

    var updateData struct {
        Name string `json:"name"`
    }

    if err := c.BodyParser(&updateData); err != nil {
        return utils.ErrorResponse(c, fiber.StatusBadRequest, "Неверный формат данных")
    }

    user, err := h.services.User().GetUserByID(c.Context(), sessionData.UserID)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка получения профиля")
    }

    user.Name = updateData.Name
    err = h.services.User().UpdateUser(c.Context(), user)
    if err != nil {
        return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка обновления профиля")
    }

    return utils.SuccessResponse(c, user)
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