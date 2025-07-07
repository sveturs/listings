package handler

import (
	"log"

	"backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes регистрирует все маршруты проекта
func (h *Handler) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
	// API v1 group
	api := app.Group("/api/v1")

	// Публичные маршруты для поиска
	searchGroup := api.Group("/search")

	// Публичные GET маршруты для конфигурации
	configGroup := searchGroup.Group("/config")
	configGroup.Get("/", h.GetConfig)

	// Публичные GET маршруты для весов
	weightsGroup := configGroup.Group("/weights")
	weightsGroup.Get("/", h.GetWeights)
	weightsGroup.Get("/:field", h.GetWeightByField)

	// Публичные GET маршруты для синонимов
	synonymsGroup := configGroup.Group("/synonyms")
	synonymsGroup.Get("/", h.GetSynonyms)

	// Публичные GET маршруты для транслитерации
	translitGroup := configGroup.Group("/transliteration")
	translitGroup.Get("/", h.GetTransliterationRules)

	// Публичные маршруты для статистики
	statsGroup := searchGroup.Group("/statistics")
	statsGroup.Get("/", h.GetSearchStatistics)
	statsGroup.Get("/popular", h.GetPopularSearches)

	// Защищенные admin маршруты
	adminRoutes := app.Group("/api/v1/admin", mw.AuthRequiredJWT, mw.AdminRequired)

	// Аналитика поиска
	adminRoutes.Get("/search/analytics", h.GetSearchAnalytics)

	// Конфигурация поиска
	adminRoutes.Put("/search/config", h.UpdateConfig)

	// Веса
	adminRoutes.Get("/search/weights", h.GetWeights)
	adminRoutes.Post("/search/weights", h.CreateWeight)
	adminRoutes.Put("/search/weights/:id", h.UpdateWeight)
	adminRoutes.Delete("/search/weights/:id", h.DeleteWeight)

	// Синонимы
	adminRoutes.Get("/search/synonyms", h.GetSynonyms)
	adminRoutes.Post("/search/synonyms", func(c *fiber.Ctx) error {
		log.Printf("POST /search/synonyms - middleware reached")
		return h.CreateSynonym(c)
	})
	adminRoutes.Put("/search/synonyms/:id", h.UpdateSynonym)
	adminRoutes.Delete("/search/synonyms/:id", h.DeleteSynonym)

	// Транслитерация
	adminRoutes.Post("/search/transliteration", h.CreateTransliterationRule)
	adminRoutes.Put("/search/transliteration/:id", h.UpdateTransliterationRule)
	adminRoutes.Delete("/search/transliteration/:id", h.DeleteTransliterationRule)

	return nil
}

// GetPrefix возвращает префикс проекта для логирования
func (h *Handler) GetPrefix() string {
	return "search_admin"
}
