package tracking

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"backend/internal/middleware"
	"backend/internal/storage/postgres"
)

// Module представляет модуль трекинга доставок
type Module struct {
	trackingHandler *TrackingHandler
	DeliveryService *DeliveryService
	Hub             *Hub
}

// NewModule создает новый модуль трекинга
func NewModule(db *postgres.Database) *Module {
	// Создаем сервисы для трекинга
	deliveryService := NewDeliveryService(db)
	courierService := NewCourierService(db)
	hub := NewHub()

	// Запускаем Hub в отдельной горутине
	go hub.Run(context.Background())

	// Создаем handler
	trackingHandler := NewTrackingHandler(deliveryService, courierService, hub)

	return &Module{
		trackingHandler: trackingHandler,
		DeliveryService: deliveryService,
		Hub:             hub,
	}
}

// RegisterRoutes регистрирует роуты модуля трекинга
func (m *Module) RegisterRoutes(app *fiber.App, middleware *middleware.Middleware) error {
	api := app.Group("/api/v1")

	// Публичные роуты для трекинга
	tracking := api.Group("/tracking")
	tracking.Get("/:token", m.trackingHandler.GetDeliveryByToken)
	tracking.Get("/connections", m.trackingHandler.GetWebSocketConnections)

	// Роуты для курьеров (требуют авторизации)
	courier := api.Group("/courier", middleware.AuthRequiredJWT)
	courier.Post("/:courier_id/location", m.trackingHandler.UpdateCourierLocation)
	courier.Get("/:courier_id/deliveries", m.trackingHandler.GetActiveDeliveries)

	// Роуты для управления доставками
	delivery := api.Group("/delivery", middleware.AuthRequiredJWT)
	delivery.Post("/", m.trackingHandler.CreateDelivery)
	delivery.Put("/:delivery_id/status", m.trackingHandler.UpdateDeliveryStatus)

	return nil
}

// GetPrefix возвращает префикс модуля для логирования
func (m *Module) GetPrefix() string {
	return "tracking"
}
