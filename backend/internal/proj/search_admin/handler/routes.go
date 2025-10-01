package handler

import (
	"log"

	"backend/internal/middleware"

	"github.com/gofiber/fiber/v2"
	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"
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

	// Защищенные admin маршруты для поиска (специфичная группа)
	adminSearchRoutes := app.Group("/api/v1/admin/search", mw.JWTParser(), authMiddleware.RequireAuth(), mw.AdminRequired)

	// Аналитика поиска - REMOVED (deprecated, use /api/v1/analytics/metrics/search instead)

	// Конфигурация поиска
	adminSearchRoutes.Put("/config", h.UpdateConfig)

	// Веса
	adminSearchRoutes.Get("/weights", h.GetWeights)
	adminSearchRoutes.Post("/weights", h.CreateWeight)
	adminSearchRoutes.Put("/weights/:id", h.UpdateWeight)
	adminSearchRoutes.Delete("/weights/:id", h.DeleteWeight)

	// Синонимы
	adminSearchRoutes.Get("/synonyms", h.GetSynonyms)
	adminSearchRoutes.Post("/synonyms", func(c *fiber.Ctx) error {
		log.Printf("POST /search/synonyms - middleware reached")
		return h.CreateSynonym(c)
	})
	adminSearchRoutes.Put("/synonyms/:id", h.UpdateSynonym)
	adminSearchRoutes.Delete("/synonyms/:id", h.DeleteSynonym)

	// Транслитерация
	adminSearchRoutes.Post("/transliteration", h.CreateTransliterationRule)
	adminSearchRoutes.Put("/transliteration/:id", h.UpdateTransliterationRule)
	adminSearchRoutes.Delete("/transliteration/:id", h.DeleteTransliterationRule)

	// Управление индексом
	indexGroup := adminSearchRoutes.Group("/index")
	indexGroup.Get("/info", h.GetIndexInfo)
	indexGroup.Get("/statistics", h.GetIndexStatistics)
	indexGroup.Get("/documents", h.SearchIndexedDocuments)
	indexGroup.Post("/reindex", h.ReindexDocuments)

	return nil
}

// GetPrefix возвращает префикс проекта для логирования
func (h *Handler) GetPrefix() string {
	return "search_admin"
}
