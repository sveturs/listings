// backend/internal/proj/marketplace/handler/search.go
package handler

import (
	"backend/internal/domain/models"
	"backend/internal/domain/search"
	globalService "backend/internal/proj/global/service"
	"backend/internal/proj/marketplace/service"
	"backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"log"
	"math"
	"strconv"
)

// SearchHandler обрабатывает запросы, связанные с поиском
type SearchHandler struct {
	services           globalService.ServicesInterface
	marketplaceService service.MarketplaceServiceInterface
}

// NewSearchHandler создает новый обработчик поиска
func NewSearchHandler(services globalService.ServicesInterface) *SearchHandler {
	return &SearchHandler{
		services:           services,
		marketplaceService: services.Marketplace(),
	}
}

// SearchListingsAdvanced выполняет расширенный поиск объявлений
func (h *SearchHandler) SearchListingsAdvanced(c *fiber.Ctx) error {
	// Парсим параметры поиска из запроса
	var params search.ServiceParams
	if err := c.BodyParser(&params); err != nil {
		log.Printf("Failed to parse search params: %v", err)

		// Попробуем разобрать запрос как form-data
		params = search.ServiceParams{
			Query:         c.FormValue("query"),
			Page:          parseIntOrDefault(c.FormValue("page"), 1),
			Size:          parseIntOrDefault(c.FormValue("limit"), 20),
			Sort:          c.FormValue("sort_by"),
			SortDirection: c.FormValue("sort_order"),
		}

		// Разбор фильтров из form-data
		categoryID := c.FormValue("category_id")
		if categoryID != "" {
			params.CategoryID = categoryID
		}

		minPrice := c.FormValue("min_price")
		if minPrice != "" {
			price, err := strconv.ParseFloat(minPrice, 64)
			if err == nil && price > 0 {
				params.PriceMin = price
			}
		}

		maxPrice := c.FormValue("max_price")
		if maxPrice != "" {
			price, err := strconv.ParseFloat(maxPrice, 64)
			if err == nil && price > 0 {
				params.PriceMax = price
			}
		}

		// Обработка фильтров атрибутов
		// Собираем все параметры, начинающиеся с "attr_"
		attributeFilters := make(map[string]string)
		c.Context().QueryArgs().VisitAll(func(key, value []byte) {
			keyStr := string(key)
			if len(keyStr) > 5 && keyStr[:5] == "attr_" {
				attrName := keyStr[5:]
				attributeFilters[attrName] = string(value)
			}
		})

		if len(attributeFilters) > 0 {
			params.AttributeFilters = attributeFilters
		}
	}

	// Получаем язык из контекста или из запроса
	lang := c.Query("lang")
	if lang == "" {
		if ctxLang, ok := c.Locals("language").(string); ok && ctxLang != "" {
			lang = ctxLang
		} else {
			lang = "ru" // Язык по умолчанию
		}
	}

	// Устанавливаем значения по умолчанию
	if params.Size <= 0 {
		params.Size = 20
	}
	if params.Page <= 0 {
		params.Page = 1
	}

	// Ограничиваем размер страницы
	if params.Size > 100 {
		params.Size = 100
	}

	// Устанавливаем язык
	params.Language = lang

	// Проверяем режим просмотра карты
	viewMode := c.Query("view_mode")
	if viewMode == "map" {
		// Для режима карты увеличиваем лимит
		params.Size = 5000
	}

	// Создаем контекст с языком
	ctx := c.Context()
	ctx.SetUserValue("language", lang)

	// Выполняем поиск
	results, err := h.marketplaceService.SearchListingsAdvanced(ctx, &params)
	if err != nil {
		log.Printf("Failed to perform advanced search: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Ошибка при выполнении поиска")
	}

	// Проверяем, что results.Items не nil
	items := results.Items
	if items == nil {
		items = []*models.MarketplaceListing{}
	}

	// Вычисляем метаданные пагинации
	total := results.Total
	totalPages := int(math.Ceil(float64(total) / float64(params.Size)))
	hasMore := params.Page < totalPages

	// Получаем текущую страницу из параметров (или значение по умолчанию)
	page := params.Page
	if page <= 0 {
		page = 1
	}

	// Получаем размер страницы из параметров (или значение по умолчанию)
	size := params.Size
	if size <= 0 {
		size = 20
	}

	// Логируем метаданные пагинации для отладки
	log.Printf("Метаданные пагинации: total=%d, total_pages=%d, page=%d, size=%d, has_more=%v",
		total, totalPages, page, size, hasMore)

	// ВАЖНОЕ ИЗМЕНЕНИЕ: структура, соответствующая ожиданиям фронтенда
	response := fiber.Map{
		"data": items,
		"meta": fiber.Map{
			"total":               total,
			"page":                page,
			"size":                size,
			"total_pages":         totalPages,
			"has_more":            hasMore,
			"facets":              results.Facets,
			"suggestions":         results.Suggestions,
			"spelling_suggestion": results.SpellingSuggestion,
			"took_ms":             results.Took,
		},
	}

	// ИЗМЕНЕНИЕ: теперь прямой возврат response вместо utils.SuccessResponse
	return c.JSON(response)
}

// parseIntOrDefault преобразует строку в число, возвращая значение по умолчанию в случае ошибки
func parseIntOrDefault(str string, defaultValue int) int {
	if str == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(str)
	if err != nil {
		return defaultValue
	}

	return value
}

// GetSuggestions возвращает предложения автодополнения
func (h *SearchHandler) GetSuggestions(c *fiber.Ctx) error {
	// Получаем префикс для автодополнения из параметров
	prefix := c.Query("prefix")
	if prefix == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Необходимо указать префикс")
	}

	// Получаем размер выборки
	size := 10
	if sizeStr := c.Query("size"); sizeStr != "" {
		if parsedSize, err := strconv.Atoi(sizeStr); err == nil && parsedSize > 0 {
			size = parsedSize
		}
	}

	// Получаем предложения
	suggestions, err := h.marketplaceService.GetSuggestions(c.Context(), prefix, size)
	if err != nil {
		log.Printf("Failed to get suggestions for prefix '%s': %v", prefix, err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось получить предложения")
	}

	// Возвращаем предложения
	return utils.SuccessResponse(c, suggestions)
}

// GetEnhancedSuggestions возвращает расширенные предложения автодополнения
func (h *SearchHandler) GetEnhancedSuggestions(c *fiber.Ctx) error {
	// Получаем префикс для автодополнения из параметров
	prefix := c.Query("prefix")
	if prefix == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Необходимо указать префикс")
	}

	// Получаем размер выборки
	size := 10
	if sizeStr := c.Query("size"); sizeStr != "" {
		if parsedSize, err := strconv.Atoi(sizeStr); err == nil && parsedSize > 0 {
			size = parsedSize
		}
	}

	// Получаем язык из контекста или из запроса
	lang := c.Query("lang")
	if lang == "" {
		if ctxLang, ok := c.Locals("language").(string); ok && ctxLang != "" {
			lang = ctxLang
		} else {
			lang = "ru" // Язык по умолчанию
		}
	}

	// Получаем предложения для текущего языка
	suggestions, err := h.marketplaceService.GetSuggestions(c.Context(), prefix, size)
	if err != nil {
		log.Printf("Failed to get suggestions for prefix '%s': %v", prefix, err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось получить предложения")
	}

	// Для более короткого запроса (менее 3 символов) добавляем категории
	var categorySuggestions []models.CategorySuggestion
	if len(prefix) < 5 {
		categorySuggestions, err = h.marketplaceService.GetCategorySuggestions(c.Context(), prefix, size)
		if err != nil {
			log.Printf("Failed to get category suggestions for prefix '%s': %v", prefix, err)
			// Продолжаем даже при ошибке
		}
	}

	// Создаем расширенный ответ
	type enhancedSuggestion struct {
		Text     string `json:"text"`
		Type     string `json:"type"`
		Category *struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Slug string `json:"slug"`
		} `json:"category,omitempty"`
	}

	var result []enhancedSuggestion

	// Добавляем текстовые предложения
	for _, text := range suggestions {
		result = append(result, enhancedSuggestion{
			Text: text,
			Type: "text",
		})
	}

	// Добавляем категории
	for _, cat := range categorySuggestions {
		result = append(result, enhancedSuggestion{
			Text: cat.Name,
			Type: "category",
			Category: &struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
				Slug string `json:"slug"`
			}{
				ID:   cat.ID,
				Name: cat.Name,
				Slug: "", // CategorySuggestion не имеет поля Slug
			},
		})
	}

	// Возвращаем расширенные предложения
	return utils.SuccessResponse(c, result)
}

// GetCategorySuggestions возвращает предложения категорий
func (h *SearchHandler) GetCategorySuggestions(c *fiber.Ctx) error {
	// Получаем строку запроса
	query := c.Query("query")
	if query == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Необходимо указать запрос")
	}

	// Получаем размер выборки
	size := 10
	if sizeStr := c.Query("size"); sizeStr != "" {
		if parsedSize, err := strconv.Atoi(sizeStr); err == nil && parsedSize > 0 {
			size = parsedSize
		}
	}

	// Получаем предложения категорий
	suggestions, err := h.marketplaceService.GetCategorySuggestions(c.Context(), query, size)
	if err != nil {
		log.Printf("Failed to get category suggestions for query '%s': %v", query, err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось получить предложения категорий")
	}

	// Возвращаем предложения категорий
	return utils.SuccessResponse(c, suggestions)
}

// GetSimilarListings возвращает похожие объявления
func (h *SearchHandler) GetSimilarListings(c *fiber.Ctx) error {
	// Получаем ID объявления из параметров URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Некорректный ID объявления")
	}

	// Получаем лимит
	limit := 5
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Получаем похожие объявления
	listings, err := h.marketplaceService.GetSimilarListings(c.Context(), id, limit)
	if err != nil {
		log.Printf("Failed to get similar listings for listing %d: %v", id, err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Не удалось получить похожие объявления")
	}

	// Проверяем, что listings не nil
	if listings == nil {
		listings = []*models.MarketplaceListing{}
	}

	// Возвращаем похожие объявления
	return utils.SuccessResponse(c, listings)
}
