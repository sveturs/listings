package handler

import (
	"backend/internal/middleware"
	"backend/internal/proj/global/service"
	"github.com/gofiber/fiber/v2"
)

// Handler основная структура handler для витрин
type Handler struct {
	services        service.ServicesInterface
	storefrontHandler *StorefrontHandler
}

// NewHandler создает новый handler
func NewHandler(services service.ServicesInterface) *Handler {
	storefrontService := services.Storefront()
	
	return &Handler{
		services:          services,
		storefrontHandler: NewStorefrontHandler(storefrontService),
	}
}

// GetPrefix возвращает префикс для маршрутов
func (h *Handler) GetPrefix() string {
	return "/storefronts"
}

// RegisterRoutes регистрирует маршруты
func (h *Handler) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
	api := app.Group("/api/v1")
	
	// Публичные маршруты (без авторизации)
	public := api.Group("/storefronts")
	{
		// Получение витрины
		public.Get("/:id", h.storefrontHandler.GetStorefront)
		public.Get("/slug/:slug", h.storefrontHandler.GetStorefrontBySlug)
		
		// Списки и поиск
		public.Get("/", h.storefrontHandler.ListStorefronts)
		public.Get("/nearby", h.storefrontHandler.GetNearbyStorefronts)
		
		// Картографические данные
		public.Get("/map", h.storefrontHandler.GetMapData)
		public.Get("/building", h.storefrontHandler.GetBusinessesInBuilding)
		
		// Персонал (просмотр)
		public.Get("/:id/staff", h.storefrontHandler.GetStaff)
		
		// Аналитика (запись просмотра)
		public.Post("/:id/view", h.storefrontHandler.RecordView)
	}

	// Защищенные маршруты (требуют авторизации)
	protected := api.Group("/storefronts")
	protected.Use(mw.AuthRequiredJWT)
	{
		// Управление витринами
		protected.Post("/", h.storefrontHandler.CreateStorefront)
		protected.Put("/:id", h.storefrontHandler.UpdateStorefront)
		protected.Delete("/:id", h.storefrontHandler.DeleteStorefront)
		
		// Мои витрины
		protected.Get("/my", h.storefrontHandler.GetMyStorefronts)
		
		// Настройки витрины
		protected.Put("/:id/hours", h.storefrontHandler.UpdateWorkingHours)
		protected.Put("/:id/payment-methods", h.storefrontHandler.UpdatePaymentMethods)
		protected.Put("/:id/delivery-options", h.storefrontHandler.UpdateDeliveryOptions)
		
		// Управление персоналом
		protected.Post("/:id/staff", h.storefrontHandler.AddStaff)
		protected.Put("/:id/staff/:staffId/permissions", h.storefrontHandler.UpdateStaffPermissions)
		protected.Delete("/:id/staff/:userId", h.storefrontHandler.RemoveStaff)
		
		// Загрузка изображений
		protected.Post("/:id/logo", h.storefrontHandler.UploadLogo)
		protected.Post("/:id/banner", h.storefrontHandler.UploadBanner)
		
		// Аналитика (просмотр)
		protected.Get("/:id/analytics", h.storefrontHandler.GetAnalytics)
	}
	
	return nil
}