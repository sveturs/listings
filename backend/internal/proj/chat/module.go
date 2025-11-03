package chat

import (
	"backend/internal/config"
	"backend/internal/proj/global/service"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Module struct{}

func New(services service.ServicesInterface, cfg *config.Config, jwtParserMW fiber.Handler) *Module {
	return &Module{}
}

func (m *Module) RegisterRoutes(app *fiber.App) error {
	// WebSocket роут для чата - временно возвращает сообщение об отключении
	app.Get("/ws/chat", func(c *fiber.Ctx) error {
		// Проверяем, что это WebSocket запрос
		if websocket.IsWebSocketUpgrade(c) {
			// WebSocket upgrade - отправим сообщение об отключении
			return websocket.New(func(conn *websocket.Conn) {
				// Отправляем сообщение об отключении функционала
				_ = conn.WriteJSON(fiber.Map{
					"type":    "error",
					"message": "Chat functionality is temporarily disabled during migration to microservice architecture",
					"code":    "CHAT_DISABLED",
				})
				_ = conn.Close()
			})(c)
		}

		// Обычный HTTP запрос - возвращаем JSON ошибку
		return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "chat.temporarily_disabled")
	})

	return nil
}
