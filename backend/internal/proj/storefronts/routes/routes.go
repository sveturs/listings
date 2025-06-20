package routes

import (
	"backend/internal/middleware"
	"backend/internal/proj/storefronts/handler"

	"github.com/gofiber/fiber/v2"
)

// RegisterStorefrontRoutes регистрирует маршруты для витрин
func RegisterStorefrontRoutes(app *fiber.App, h *handler.StorefrontHandler, authMiddleware *middleware.Middleware) {
	api := app.Group("/api/v1")

	// Публичные маршруты (без авторизации)
	public := api.Group("/storefronts")
	{
		// Получение витрины
		public.Get("/:id", h.GetStorefront)
		public.Get("/slug/:slug", h.GetStorefrontBySlug)
		
		// Списки и поиск
		public.Get("/", h.ListStorefronts)
		public.Get("/nearby", h.GetNearbyStorefronts)
		
		// Картографические данные
		public.Get("/map", h.GetMapData)
		public.Get("/building", h.GetBusinessesInBuilding)
		
		// Персонал (просмотр)
		public.Get("/:id/staff", h.GetStaff)
		
		// Аналитика (запись просмотра)
		public.Post("/:id/view", h.RecordView)
	}

	// Защищенные маршруты (требуют авторизации)
	protected := api.Group("/storefronts")
	protected.Use(authMiddleware.AuthRequiredJWT)
	{
		// Мои витрины (должен быть перед /:id чтобы не конфликтовать)
		protected.Get("/my", h.GetMyStorefronts)
		
		// Управление витринами
		protected.Post("/", h.CreateStorefront)
		protected.Put("/:id", h.UpdateStorefront)
		protected.Delete("/:id", h.DeleteStorefront)
		
		// Настройки витрины
		protected.Put("/:id/hours", h.UpdateWorkingHours)
		protected.Put("/:id/payment-methods", h.UpdatePaymentMethods)
		protected.Put("/:id/delivery-options", h.UpdateDeliveryOptions)
		
		// Управление персоналом
		protected.Post("/:id/staff", h.AddStaff)
		protected.Put("/:id/staff/:staffId/permissions", h.UpdateStaffPermissions)
		protected.Delete("/:id/staff/:userId", h.RemoveStaff)
		
		// Аналитика (просмотр)
		protected.Get("/:id/analytics", h.GetStorefrontAnalytics)
	}
}

// RegisterStorefrontWebhooks регистрирует вебхуки для интеграций
func RegisterStorefrontWebhooks(app *fiber.App) {
	webhooks := app.Group("/webhooks/storefronts")
	
	// Платежные системы
	webhooks.Post("/payment/postanska", handlePostanskaWebhook)
	webhooks.Post("/payment/kekspay", handleKeksPayWebhook)
	
	// Службы доставки
	webhooks.Post("/delivery/postasrbije", handlePostaSrbijeWebhook)
	webhooks.Post("/delivery/aks", handleAKSWebhook)
	webhooks.Post("/delivery/bex", handleBEXWebhook)
}

// Заглушки для вебхуков (будут реализованы позже)
func handlePostanskaWebhook(c *fiber.Ctx) error {
	// TODO: Implement Poštanska štedionica webhook
	return c.SendStatus(fiber.StatusOK)
}

func handleKeksPayWebhook(c *fiber.Ctx) error {
	// TODO: Implement Keks Pay webhook
	return c.SendStatus(fiber.StatusOK)
}

func handlePostaSrbijeWebhook(c *fiber.Ctx) error {
	// TODO: Implement Pošta Srbije tracking webhook
	return c.SendStatus(fiber.StatusOK)
}

func handleAKSWebhook(c *fiber.Ctx) error {
	// TODO: Implement AKS tracking webhook
	return c.SendStatus(fiber.StatusOK)
}

func handleBEXWebhook(c *fiber.Ctx) error {
	// TODO: Implement BEX tracking webhook
	return c.SendStatus(fiber.StatusOK)
}