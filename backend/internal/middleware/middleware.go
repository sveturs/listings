// backend/internal/middleware/middleware.go
package middleware

import (
	"backend/internal/config"
	"backend/pkg/utils"
	"log"

	"github.com/gofiber/fiber/v2"

	globalService "backend/internal/proj/global/service"
)

type Middleware struct {
	config   *config.Config
	services globalService.ServicesInterface
}

func NewMiddleware(cfg *config.Config, services globalService.ServicesInterface) *Middleware {
	return &Middleware{
		config:   cfg,
		services: services,
	}
}

func (m *Middleware) Setup(app *fiber.App) {
	app.Use(m.Logger())
	app.Use(m.CORS())
}

// ErrorHandler обрабатывает все ошибки приложения
func (m *Middleware) ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	log.Printf("Error handling request: %v", err)

	return utils.ErrorResponse(c, code, message)
}
func (m *Middleware) AdminRequired(c *fiber.Ctx) error {
	log.Printf("AdminRequired middleware called for path: %s", c.Path())

	// Проверяем, является ли пользователь администратором
	userID, ok := c.Locals("user_id").(int)
	log.Printf("AdminRequired check for user ID: %d, isOk: %v", userID, ok)

	if !ok || userID == 0 {
		log.Printf("AdminRequired: User ID not found or invalid")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	// Сначала проверяем ID пользователя (для обратной совместимости)
	if userID == 1 || userID == 2 || userID == 3 {
		log.Printf("AdminRequired: Access granted for user ID %d", userID)
		return c.Next()
	}

	log.Printf("AdminRequired: User ID %d not in admin list, checking email", userID)

	// Если ID не подходит, проверяем email пользователя
	ctx := c.Context()

	if m.services == nil {
		log.Printf("AdminRequired: Error: services is nil")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	userService := m.services.User()
	if userService == nil {
		log.Printf("AdminRequired: Error: User service is nil")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	user, err := userService.GetUserByID(ctx, userID)
	if err != nil {
		log.Printf("AdminRequired: Error getting user with ID %d: %v", userID, err)
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Отказано в доступе")
	}

	if user == nil {
		log.Printf("AdminRequired: User with ID %d not found", userID)
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Отказано в доступе")
	}

	log.Printf("AdminRequired: Found user with ID %d, email: %s", userID, user.Email)

	// Проверка администратора по БД (нужно учитывать, что таблица может еще не существовать)
	isAdmin, err := userService.IsUserAdmin(ctx, user.Email)
	if err != nil {
		// Если таблица еще не создана, проверяем по жесткому списку
		log.Printf("AdminRequired: Error checking admin status in DB: %v, falling back to hardcoded list", err)
	} else if isAdmin {
		log.Printf("AdminRequired: Access granted for admin email %s (from DB)", user.Email)
		return c.Next()
	}

	// Резервный список email-адресов администраторов (на случай если миграция еще не применена)
	adminEmails := []string{
		"bevzenko.sergey@gmail.com",
		"voroshilovdo@gmail.com",
		// Здесь можно добавить дополнительные email-адреса администраторов
	}

	log.Printf("AdminRequired: Checking user email %s against hardcoded admin emails: %v", user.Email, adminEmails)

	// Проверяем, является ли email пользователя админским
	for _, adminEmail := range adminEmails {
		if user.Email == adminEmail {
			log.Printf("AdminRequired: Access granted for admin email %s (from hardcoded list)", user.Email)
			return c.Next()
		}
	}

	log.Printf("AdminRequired: User ID %d with email %s is not an admin", userID, user.Email)
	return utils.ErrorResponse(c, fiber.StatusForbidden, "Отказано в доступе")
}
