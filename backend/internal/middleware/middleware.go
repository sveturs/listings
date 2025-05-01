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
	//	log.Printf("AdminRequired check for user ID: %d, isOk: %v", userID, ok)

	if !ok || userID == 0 {
		//		log.Printf("AdminRequired: User ID not found or invalid")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Требуется авторизация")
	}

	// В этой простой версии считаем администраторами только пользователей с ID 1, 2, 3
	if userID != 1 && userID != 2 && userID != 3 {
		//		log.Printf("AdminRequired: User ID %d is not an admin", userID)
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Отказано в доступе")
	}

	//	log.Printf("AdminRequired: Access granted for user ID %d", userID)
	return c.Next()
}
