package viber

import (
	"github.com/gofiber/fiber/v2"

	"backend/internal/middleware"
	"backend/internal/proj/global/service"
	"backend/internal/storage/postgres"
)

// Module представляет модуль Viber Bot
type Module struct {
	viberHandler *ViberHandler
}

// NewModule создает новый модуль Viber
func NewModule(services service.ServicesInterface) *Module {
	// Получаем database из services
	db := services.Storage().(*postgres.Database)

	// Создаем handler с сервисами
	viberHandler := NewViberHandler(db, services)

	return &Module{
		viberHandler: viberHandler,
	}
}

// RegisterRoutes регистрирует роуты модуля Viber
func (m *Module) RegisterRoutes(app *fiber.App, middleware *middleware.Middleware) error {
	api := app.Group("/api")

	// Роуты для Viber webhook (публичные)
	viber := api.Group("/viber")
	viber.Post("/webhook", m.viberHandler.HandleViberWebhook)
	viber.Post("/infobip-webhook", m.viberHandler.HandleInfobipWebhook)

	// Роуты для отправки сообщений (требуют авторизации)
	viberAPI := api.Group("/viber", middleware.AuthRequiredJWT)
	viberAPI.Post("/send", m.viberHandler.SendMessage)
	viberAPI.Post("/send-tracking", m.viberHandler.SendTrackingNotification)
	viberAPI.Get("/stats", m.viberHandler.GetSessionStats)
	viberAPI.Post("/estimate-cost", m.viberHandler.EstimateMessageCost)

	return nil
}

// GetPrefix возвращает префикс модуля для логирования
func (m *Module) GetPrefix() string {
	return "viber"
}