// backend/internal/proj/global/handler/handler.go
package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"backend/internal/config"
	"backend/internal/middleware"
	globalService "backend/internal/proj/global/service"
)

// Handler объединяет все глобальные обработчики
type Handler struct {
	UnifiedSearch *UnifiedSearchHandler
	service       globalService.ServicesInterface
	searchWeights *config.SearchWeights
}

// NewHandler создает новый глобальный обработчик
func NewHandler(services globalService.ServicesInterface, searchWeights *config.SearchWeights) *Handler {
	return &Handler{
		UnifiedSearch: NewUnifiedSearchHandler(services),
		service:       services,
		searchWeights: searchWeights,
	}
}

// GetPrefix возвращает префикс для глобальных API
func (h *Handler) GetPrefix() string {
	return "/api/v1"
}

// RegisterRoutes регистрирует все глобальные маршруты
func (h *Handler) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
	// Регистрируем унифицированный поиск напрямую в app,
	// чтобы избежать конфликтов с другими middleware
	app.Get("/api/v1/search", h.UnifiedSearch.UnifiedSearch)

	// Добавляем алиас для suggestions
	app.Get("/api/v1/search/suggestions", h.GetSuggestions)

	return nil
}

// GetSuggestions обрабатывает запросы автодополнения для поиска
// @Summary Get search suggestions
// @Description Returns autocomplete suggestions with types (search, product, category)
// @Tags search
// @Accept json
// @Produce json
// @Param q query string false "Search query (alias for prefix)"
// @Param prefix query string false "Search prefix"
// @Param limit query int false "Number of suggestions (alias for size)" default(10)
// @Param size query int false "Number of suggestions" default(10)
// @Param types query string false "Comma-separated types: queries,categories,products" default("queries,categories,products")
// @Success 200 {object} backend_pkg_utils.SuccessResponseSwag{data=[]backend_internal_proj_marketplace_service.SuggestionItem} "Enhanced suggestions list"
// @Failure 400 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.prefixRequired"
// @Failure 500 {object} backend_pkg_utils.ErrorResponseSwag "marketplace.suggestionsError"
// @Router /api/v1/search/suggestions [get]
func (h *Handler) GetSuggestions(c *fiber.Ctx) error {
	// Получаем префикс - поддерживаем оба параметра для совместимости
	prefix := c.Query("q")
	if prefix == "" {
		prefix = c.Query("prefix")
	}

	if prefix == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Query parameter is required")
	}

	// Получаем размер выборки
	sizeStr := c.Query("limit")
	if sizeStr == "" {
		sizeStr = c.Query("size", "10")
	}
	size := 10
	if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 {
		size = s
	}

	// Получаем типы подсказок
	types := c.Query("types", "queries,categories,products")

	// Используем GetEnhancedSuggestions для получения расширенных подсказок
	suggestions, err := h.service.Marketplace().GetEnhancedSuggestions(c.Context(), prefix, size, types)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to get suggestions")
	}

	// Возвращаем в формате совместимом с API
	return c.JSON(fiber.Map{
		"data": suggestions,
	})
}
