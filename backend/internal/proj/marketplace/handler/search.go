// backend/internal/proj/marketplace/handler/search.go
package handler

import (
	"context"
	"math"
	"strconv"
	"strings"
	"time"

	"backend/internal/domain/models"
	"backend/internal/domain/search"
	"backend/internal/logger"
	globalService "backend/internal/proj/global/service"
	"backend/internal/proj/marketplace/service"
	searchlogsTypes "backend/internal/proj/searchlogs/types"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
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
// @Summary Advanced search for listings
// @Description Performs advanced search with filters, facets, and suggestions
// @Tags marketplace-search
// @Accept json
// @Produce json
// @Param body body search.ServiceParams true "Search parameters"
// @Success 200 {object} handler.SearchResponse "Search results with metadata"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.searchError"
// @Router /api/v1/marketplace/search [get]
// @Router /api/v1/marketplace/search [post]
func (h *SearchHandler) SearchListingsAdvanced(c *fiber.Ctx) error {
	// Засекаем время начала для измерения производительности
	startTime := time.Now()
	// Парсим параметры поиска из запроса
	var params search.ServiceParams

	// Если это POST запрос, пробуем распарсить JSON body
	if c.Method() == "POST" {
		// Структура для POST запроса
		var postRequest struct {
			search.ServiceParams
			AdvancedGeoFilters *search.AdvancedGeoFilters `json:"advanced_geo_filters"`
		}

		if err := c.BodyParser(&postRequest); err != nil {
			logger.Error().Err(err).Msg("Failed to parse POST search params")
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "validation.failed")
		}

		params = postRequest.ServiceParams
		params.AdvancedGeoFilters = postRequest.AdvancedGeoFilters

		// Также парсим query параметры для POST запроса
		if query := c.Query("query"); query != "" {
			params.Query = query
		}
		if page := c.QueryInt("page", 0); page > 0 {
			params.Page = page
		}
		if limit := c.QueryInt("limit", 0); limit > 0 {
			params.Size = limit
		}
	} else if err := c.BodyParser(&params); err != nil {
		logger.Error().Err(err).Msg("Failed to parse search params")

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

	// Устанавливаем язык из контекста или из запроса
	if params.Language == "" {
		lang := c.Query("lang")
		if lang == "" {
			if ctxLang, ok := c.Locals("language").(string); ok && ctxLang != "" {
				lang = ctxLang
			} else {
				lang = "ru" // Язык по умолчанию
			}
		}
		params.Language = lang
	}

	// Проверяем, нужно ли использовать нечеткий поиск
	useFuzzy := c.Query("fuzzy", "true") // По умолчанию включен
	if useFuzzy == "true" || useFuzzy == "1" {
		params.UseSynonyms = true
		if params.Fuzziness == "" {
			params.Fuzziness = "AUTO"
		}
	}

	// Проверяем режим просмотра карты
	viewMode := c.Query("view_mode")
	if viewMode == "map" {
		// Для режима карты увеличиваем лимит
		params.Size = 5000
	}

	// Создаем контекст с языком
	ctx := c.Context()
	ctx.SetUserValue("language", params.Language)

	// Выполняем поиск
	results, err := h.marketplaceService.SearchListingsAdvanced(ctx, &params)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to perform advanced search")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.searchError")
	}

	// Проверяем, что results.Items не nil
	items := results.Items
	if items == nil {
		items = []*models.MarketplaceListing{}
	}

	// Преобразуем []*models.MarketplaceListing в []models.MarketplaceListing
	// ВАЖНО: Используем глубокую копию для сохранения срезов (Images)
	listings := make([]models.MarketplaceListing, 0, len(items))
	for _, item := range items {
		if item != nil {
			// Создаем глубокую копию структуры
			listingCopy := *item
			// Если изображения есть, создаем копию среза
			if item.Images != nil {
				listingCopy.Images = make([]models.MarketplaceImage, len(item.Images))
				copy(listingCopy.Images, item.Images)
			}
			listings = append(listings, listingCopy)
		}
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
	logger.Info().Int("total", total).Int("totalPages", totalPages).Int("page", page).Int("size", size).Bool("hasMore", hasMore).Msg("Pagination metadata")

	// ВАЖНОЕ ИЗМЕНЕНИЕ: структура, соответствующая ожиданиям фронтенда
	response := SearchResponse{
		Data: listings,
		Meta: SearchMetadata{
			Total:              total,
			Page:               page,
			Size:               size,
			TotalPages:         totalPages,
			HasMore:            hasMore,
			Facets:             results.Facets,
			Suggestions:        results.Suggestions,
			SpellingSuggestion: results.SpellingSuggestion,
			TookMs:             results.Took,
		},
	}

	// Асинхронное логирование поискового запроса
	if searchLogsSvc := h.services.SearchLogs(); searchLogsSvc != nil {
		logger.Info().Msg("SearchLogs service is available, logging search query")

		// Извлекаем данные из контекста Fiber ДО запуска горутины
		var userID *int
		if uid, ok := c.Locals("user_id").(int); ok && uid > 0 {
			userID = &uid
		}

		// Получаем session ID из cookie или заголовков
		sessionID := c.Cookies("session_id")
		if sessionID == "" {
			sessionID = c.Get("X-Session-ID")
		}

		// Определяем тип устройства из User-Agent
		userAgent := c.Get("User-Agent")
		ipAddress := c.IP()

		go func() {
			// Вычисляем время ответа
			responseTime := time.Since(startTime).Milliseconds()

			// Определяем тип устройства из User-Agent
			deviceType := detectDeviceTypeSearch(userAgent)

			// Преобразуем filters из map[string]string в map[string]interface{}
			filtersInterface := make(map[string]interface{})
			for k, v := range params.AttributeFilters {
				filtersInterface[k] = v
			}

			// Преобразуем CategoryID в *int
			var categoryIDInt *int
			if params.CategoryID != "" {
				if catID, err := strconv.Atoi(params.CategoryID); err == nil {
					categoryIDInt = &catID
				}
			}

			// Создаем запись лога
			logEntry := &searchlogsTypes.SearchLogEntry{
				Query:           params.Query,
				UserID:          userID,
				SessionID:       sessionID, // Убрали указатель
				ResultCount:     total,
				ResponseTimeMS:  int64(responseTime),
				Filters:         filtersInterface,
				CategoryID:      categoryIDInt,
				PriceMin:        &params.PriceMin,
				PriceMax:        &params.PriceMax,
				Location:        nil, // TODO: добавить поддержку локации
				Language:        params.Language,
				DeviceType:      deviceType,
				UserAgent:       userAgent,
				IP:              ipAddress,
				SearchType:      "advanced",
				HasSpellCorrect: results.SpellingSuggestion != "",
				ClickedItems:    []int{}, // Будет заполняться позже при кликах
				Timestamp:       time.Now(),
			}

			// Логируем асинхронно
			logger.Info().
				Str("query", logEntry.Query).
				Int("results", logEntry.ResultCount).
				Int64("response_ms", logEntry.ResponseTimeMS).
				Msg("Logging search query")

			if err := searchLogsSvc.LogSearch(context.Background(), logEntry); err != nil {
				logger.Error().Err(err).Msg("Failed to log search query")
			} else {
				logger.Info().Msg("Search query logged successfully")
			}
		}()
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
// @Summary Get search suggestions
// @Description Returns autocomplete suggestions based on prefix
// @Tags marketplace-search
// @Accept json
// @Produce json
// @Param prefix query string true "Search prefix"
// @Param size query int false "Number of suggestions" default(10)
// @Success 200 {object} utils.SuccessResponseSwag{data=[]string} "Suggestions list"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.prefixRequired"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.suggestionsError"
// @Router /api/v1/marketplace/suggestions [get]
func (h *SearchHandler) GetSuggestions(c *fiber.Ctx) error {
	// Засекаем время начала для измерения производительности
	startTime := time.Now()
	// Получаем префикс для автодополнения из параметров
	prefix := c.Query("prefix")
	if prefix == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.prefixRequired")
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
		logger.Error().Err(err).Str("prefix", prefix).Msg("Failed to get suggestions")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.suggestionsError")
	}

	// Асинхронное логирование запроса на автодополнение
	if searchLogsSvc := h.services.SearchLogs(); searchLogsSvc != nil {
		// Извлекаем данные из контекста Fiber ДО запуска горутины
		var userID *int
		if uid, ok := c.Locals("user_id").(int); ok && uid > 0 {
			userID = &uid
		}

		// Получаем session ID из cookie или заголовков
		sessionID := c.Cookies("session_id")
		if sessionID == "" {
			sessionID = c.Get("X-Session-ID")
		}

		// Определяем тип устройства из User-Agent
		userAgent := c.Get("User-Agent")
		ipAddress := c.IP()

		go func() {
			// Вычисляем время ответа
			responseTime := time.Since(startTime).Milliseconds()

			// Определяем тип устройства из User-Agent
			deviceType := detectDeviceTypeSearch(userAgent)

			// Создаем запись лога
			logEntry := &searchlogsTypes.SearchLogEntry{
				Query:           prefix,
				UserID:          userID,
				SessionID:       sessionID, // Убрали указатель
				ResultCount:     len(suggestions),
				ResponseTimeMS:  responseTime,
				Filters:         nil,
				CategoryID:      nil,
				PriceMin:        nil,
				PriceMax:        nil,
				Location:        nil,
				Language:        "ru", // TODO: получать из контекста
				DeviceType:      deviceType,
				UserAgent:       userAgent,
				IP:              ipAddress,
				SearchType:      "suggestions",
				HasSpellCorrect: false,
				ClickedItems:    []int{},
				Timestamp:       time.Now(),
			}

			// Логируем асинхронно
			if err := searchLogsSvc.LogSearch(context.Background(), logEntry); err != nil {
				logger.Error().Err(err).Msg("Failed to log suggestions query")
			}
		}()
	}

	// Возвращаем предложения
	return utils.SuccessResponse(c, suggestions)
}

// GetEnhancedSuggestions возвращает расширенные предложения автодополнения
// @Summary Get enhanced search suggestions
// @Description Returns autocomplete suggestions with categories
// @Tags marketplace-search
// @Accept json
// @Produce json
// @Param prefix query string true "Search prefix"
// @Param size query int false "Number of suggestions" default(10)
// @Param lang query string false "Language" default(ru)
// @Success 200 {object} utils.SuccessResponseSwag{data=[]handler.SuggestionItem} "Enhanced suggestions list"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.prefixRequired"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.suggestionsError"
// @Router /api/v1/marketplace/enhanced-suggestions [get]
func (h *SearchHandler) GetEnhancedSuggestions(c *fiber.Ctx) error {
	// Получаем префикс для автодополнения из параметров
	prefix := c.Query("prefix")
	if prefix == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.prefixRequired")
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
		logger.Error().Err(err).Str("prefix", prefix).Msg("Failed to get suggestions")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.suggestionsError")
	}

	// Для более короткого запроса (менее 3 символов) добавляем категории
	var categorySuggestions []models.CategorySuggestion
	if len(prefix) < 5 {
		categorySuggestions, err = h.marketplaceService.GetCategorySuggestions(c.Context(), prefix, size)
		if err != nil {
			logger.Error().Err(err).Str("prefix", prefix).Msg("Failed to get category suggestions")
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
// @Summary Get category suggestions
// @Description Returns category suggestions based on query
// @Tags marketplace-search
// @Accept json
// @Produce json
// @Param query query string true "Search query"
// @Param size query int false "Number of suggestions" default(10)
// @Success 200 {object} utils.SuccessResponseSwag{data=[]models.CategorySuggestion} "Category suggestions list"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.queryRequired"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.categorySuggestionsError"
// @Router /api/v1/marketplace/category-suggestions [get]
func (h *SearchHandler) GetCategorySuggestions(c *fiber.Ctx) error {
	// Получаем строку запроса
	query := c.Query("query")
	if query == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.queryRequired")
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
		logger.Error().Err(err).Str("query", query).Msg("Failed to get category suggestions")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.categorySuggestionsError")
	}

	// Возвращаем предложения категорий
	return utils.SuccessResponse(c, suggestions)
}

// GetSimilarListings возвращает похожие объявления
// @Summary Get similar listings
// @Description Returns listings similar to a specific listing
// @Tags marketplace-search
// @Accept json
// @Produce json
// @Param id path int true "Listing ID"
// @Param limit query int false "Number of similar listings" default(5)
// @Success 200 {object} utils.SuccessResponseSwag{data=[]models.MarketplaceListing} "Similar listings list"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidId"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.similarListingsError"
// @Router /api/v1/marketplace/listings/{id}/similar [get]
func (h *SearchHandler) GetSimilarListings(c *fiber.Ctx) error {
	// Получаем ID объявления из параметров URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidId")
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
		logger.Error().Err(err).Int("listingId", id).Msg("Failed to get similar listings")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.similarListingsError")
	}

	// Проверяем, что listings не nil
	if listings == nil {
		listings = []*models.MarketplaceListing{}
	}

	// Возвращаем похожие объявления
	return utils.SuccessResponse(c, listings)
}

// detectDeviceTypeSearch определяет тип устройства по User-Agent для поиска
func detectDeviceTypeSearch(userAgent string) string {
	ua := strings.ToLower(userAgent)

	// Проверка на мобильные устройства
	mobileKeywords := []string{
		"mobile", "android", "iphone", "ipad", "ipod",
		"blackberry", "windows phone", "opera mini", "iemobile",
	}

	for _, keyword := range mobileKeywords {
		if strings.Contains(ua, keyword) {
			// Планшеты
			if strings.Contains(ua, "ipad") || strings.Contains(ua, "tablet") {
				return "tablet"
			}
			return "mobile"
		}
	}

	// Проверка на боты
	botKeywords := []string{
		"bot", "crawl", "spider", "scraper", "curl", "wget",
	}

	for _, keyword := range botKeywords {
		if strings.Contains(ua, keyword) {
			return "bot"
		}
	}

	// По умолчанию - десктоп
	return "desktop"
}
