// backend/internal/proj/storefront/handler/routes.go
package handler

import (
	"github.com/gofiber/fiber/v2"

	"backend/internal/middleware"
)

// RegisterRoutes регистрирует все маршруты для проекта storefront
func (h *Handler) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
	// Публичные маршруты
	app.Get("/api/v1/public/storefronts/:id", h.Storefront.GetPublicStorefront)

	// Защищенные маршруты с CSRF защитой
	authedAPIGroup := app.Group("/api/v1", mw.AuthRequiredJWT, mw.CSRFProtection())
	
	storefronts := authedAPIGroup.Group("/storefronts")
	storefronts.Get("/", h.Storefront.GetUserStorefronts)
	storefronts.Post("/", h.Storefront.CreateStorefront)
	storefronts.Get("/:id", h.Storefront.GetStorefront)
	storefronts.Put("/:id", h.Storefront.UpdateStorefront)
	storefronts.Delete("/:id", h.Storefront.DeleteStorefront)
	storefronts.Get("/:id/import-sources", h.Storefront.GetImportSources)
	storefronts.Post("/import-sources", h.Storefront.CreateImportSource)
	storefronts.Put("/import-sources/:id", h.Storefront.UpdateImportSource)
	storefronts.Delete("/import-sources/:id", h.Storefront.DeleteImportSource)
	storefronts.Post("/import-sources/:id/run", h.Storefront.RunImport)
	storefronts.Get("/import-sources/:id/history", h.Storefront.GetImportHistory)
	storefronts.Get("/import-sources/:id/category-mappings", h.Storefront.GetCategoryMappings)
	storefronts.Put("/import-sources/:id/category-mappings", h.Storefront.UpdateCategoryMappings)
	storefronts.Get("/import-sources/:id/imported-categories", h.Storefront.GetImportedCategories)
	storefronts.Post("/import-sources/:id/apply-category-mappings", h.Storefront.ApplyCategoryMappings)

	return nil
}

// GetPrefix возвращает префикс проекта для логирования
func (h *Handler) GetPrefix() string {
	return "/api/v1/storefronts"
}