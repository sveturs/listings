package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"backend/internal/domain/models"
	globalService "backend/internal/proj/global/service"
	"backend/internal/proj/users/service"
	"backend/pkg/utils"
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

// Response structures for Swagger documentation

// UserProfileResponse представляет ответ с профилем пользователя
type UserProfileResponse struct {
	Success bool                `json:"success" example:"true"`
	Data    *models.UserProfile `json:"data"`
}

// MessageResponse представляет ответ с сообщением
type MessageResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Операция выполнена успешно"`
}

// RegisterResponse представляет ответ после регистрации
type RegisterResponse struct {
	Success bool         `json:"success" example:"true"`
	Message string       `json:"message" example:"Пользователь успешно зарегистрирован"`
	User    *models.User `json:"user"`
}

// PublicUserResponse представляет публичные данные пользователя
type PublicUserResponse struct {
	ID         int    `json:"id" example:"1"`
	Name       string `json:"name" example:"Иван Иванов"`
	Email      string `json:"email" example:"user@example.com"`
	PictureURL string `json:"picture_url" example:"https://example.com/avatar.jpg"`
	CreatedAt  string `json:"created_at" example:"2023-01-01T12:00:00Z"`
}

// PublicUserResponseWrapper обертка для публичного профиля
type PublicUserResponseWrapper struct {
	Success bool               `json:"success" example:"true"`
	Data    PublicUserResponse `json:"data"`
}

// AdminCheckResponse представляет ответ проверки администратора
type AdminCheckResponse struct {
	Success bool `json:"success" example:"true"`
	IsAdmin bool `json:"is_admin" example:"false"`
}

// AdminCheckResponseWrapper обертка для проверки администратора
type AdminCheckResponseWrapper struct {
	Success bool               `json:"success" example:"true"`
	Data    AdminCheckResponse `json:"data"`
}

// RegisterRequest представляет запрос на регистрацию
type RegisterRequest struct {
	Name  string `json:"name" validate:"required" example:"Иван Иванов"`
	Email string `json:"email" validate:"required,email" example:"user@example.com"`
}

// GetProfile получает профиль текущего пользователя
// @Summary Получить профиль пользователя
// @Description Возвращает полный профиль авторизованного пользователя
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} UserProfileResponse
// @Failure 401 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Security BearerAuth
// @Router /api/v1/users/profile [get]
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	profile, err := h.services.User().GetUserProfile(c.Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "users.profile.error.fetch")
	}

	return c.JSON(UserProfileResponse{
		Success: true,
		Data:    profile,
	})
}

// UpdateProfile обновляет профиль текущего пользователя
// @Summary Обновить профиль пользователя
// @Description Обновляет профиль авторизованного пользователя
// @Tags Users
// @Accept json
// @Produce json
// @Param profile body models.UserProfileUpdate true "Данные для обновления профиля"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} utils.ErrorResponseSwag
// @Failure 401 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Security BearerAuth
// @Router /api/v1/users/profile [put]
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)

	var update models.UserProfileUpdate
	if err := c.BodyParser(&update); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "users.profile.error.invalid_data")
	}

	// Валидация данных
	if err := update.Validate(); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "users.profile.error.validation")
	}

	err := h.services.User().UpdateUserProfile(c.Context(), userID, &update)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "users.profile.error.update")
	}

	return c.JSON(MessageResponse{
		Success: true,
		Message: "users.profile.success.updated",
	})
}

// Register регистрирует нового пользователя
// @Summary Регистрация пользователя
// @Description Создает нового пользователя в системе
// @Tags Users
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "Данные для регистрации"
// @Success 201 {object} RegisterResponse
// @Failure 400 {object} utils.ErrorResponseSwag
// @Failure 401 {object} utils.ErrorResponseSwag
// @Failure 500 {object} utils.ErrorResponseSwag
// @Security BearerAuth
// @Router /api/v1/users/register [post]
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var registerData RegisterRequest

	if err := c.BodyParser(&registerData); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "users.register.error.invalid_data")
	}

	// Базовая валидация
	if registerData.Name == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "users.register.error.name_required")
	}
	if registerData.Email == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "users.register.error.email_required")
	}

	user := &models.User{
		Name:  registerData.Name,
		Email: registerData.Email,
	}

	err := h.services.User().CreateUser(c.Context(), user)
	if err != nil {
		if err.Error() == "email already exists" {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "users.register.error.email_exists")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "users.register.error.create_failed")
	}

	return c.Status(fiber.StatusCreated).JSON(RegisterResponse{
		Success: true,
		Message: "users.register.success.created",
		User:    user,
	})
}

// GetProfileByID получает публичный профиль пользователя
// @Summary Получить публичный профиль пользователя
// @Description Возвращает публичную информацию о пользователе по ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} PublicUserResponseWrapper
// @Failure 400 {object} utils.ErrorResponseSwag
// @Failure 401 {object} utils.ErrorResponseSwag
// @Failure 404 {object} utils.ErrorResponseSwag
// @Security BearerAuth
// @Router /api/v1/users/{id}/profile [get]
func (h *UserHandler) GetProfileByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "users.profile.error.invalid_id")
	}

	user, err := h.userService.GetUserByID(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "users.profile.error.not_found")
	}

	publicUser := PublicUserResponse{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		PictureURL: user.PictureURL,
		CreatedAt:  user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	return c.JSON(PublicUserResponseWrapper{
		Success: true,
		Data:    publicUser,
	})
}

// IsAdminSimple проверяет, является ли пользователь администратором (простая реализация)
// @Summary Проверка статуса администратора (упрощенная)
// @Description Проверяет, является ли пользователь с указанным email администратором (упрощенная проверка по ID)
// @Tags Users
// @Accept json
// @Produce json
// @Param email path string true "Email пользователя"
// @Success 200 {object} AdminCheckResponseWrapper
// @Failure 400 {object} utils.ErrorResponseSwag
// @Failure 401 {object} utils.ErrorResponseSwag
// @Failure 404 {object} utils.ErrorResponseSwag
// @Security BearerAuth
// @Router /api/v1/users/admin-check/{email} [get]
func (h *UserHandler) IsAdminSimple(c *fiber.Ctx) error {
	email := c.Params("email")
	if email == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "users.admin_check.error.email_required")
	}

	// Получаем пользователя по email
	user, err := h.userService.GetUserByEmail(c.Context(), email)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "users.admin_check.error.user_not_found")
	}

	// Проверяем, является ли пользователь администратором
	// В этой простой версии считаем администраторами только пользователей с ID 1, 2, 3
	isAdmin := user.ID == 1 || user.ID == 2 || user.ID == 3

	return c.JSON(AdminCheckResponseWrapper{
		Success: true,
		Data: AdminCheckResponse{
			IsAdmin: isAdmin,
		},
	})
}
