// backend/internal/proj/marketplace/handler/listings.go
package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"

	"backend/internal/domain/models"
	"backend/internal/logger"
	globalService "backend/internal/proj/global/service"
	"backend/internal/proj/marketplace/cache"
	"backend/internal/proj/marketplace/service"
	searchlogsTypes "backend/internal/proj/searchlogs/types"
	"backend/internal/storage/postgres"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// Define typed context keys
type contextKey string

const (
	contextKeyUserID    contextKey = "user_id"
	contextKeyIPAddress contextKey = "ip_address"
	contextKeyLocale    contextKey = "locale"
)

// ListingsHandler обрабатывает запросы, связанные с объявлениями
type ListingsHandler struct {
	services           globalService.ServicesInterface
	marketplaceService service.MarketplaceServiceInterface
	cache              *cache.UniversalCache
}

// NewListingsHandler создает новый обработчик объявлений
func NewListingsHandler(services globalService.ServicesInterface, cache *cache.UniversalCache) *ListingsHandler {
	return &ListingsHandler{
		services:           services,
		marketplaceService: services.Marketplace(),
		cache:              cache,
	}
}

// loadUserInfoForListings загружает информацию о пользователях из auth-service для списка объявлений
func (h *ListingsHandler) loadUserInfoForListings(ctx context.Context, listings []models.MarketplaceListing) {
	// Собираем все уникальные user IDs
	userIDs := make(map[int]bool)
	for i := range listings {
		if listings[i].UserID > 0 {
			userIDs[listings[i].UserID] = true
		}
	}

	// Загружаем пользователей параллельно
	userCache := make(map[int]*models.User)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for userID := range userIDs {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			user, err := h.services.User().GetUserByID(ctx, id)
			if err != nil {
				logger.Warn().Err(err).Int("userId", id).Msg("Failed to load user from auth-service")
				return
			}
			mu.Lock()
			userCache[id] = user
			mu.Unlock()
		}(userID)
	}
	wg.Wait()

	// Заполняем информацию о пользователях в объявлениях
	for i := range listings {
		if user, ok := userCache[listings[i].UserID]; ok {
			listings[i].User = user
		}
	}
}

// CreateListing создает новое объявление
// @Summary Create new listing
// @Description Creates a new marketplace listing with attributes
// @Tags marketplace-listings
// @Accept json
// @Produce json
// @Param body body models.MarketplaceListing true "Listing data"
// @Success 200 {object} utils.SuccessResponseSwag{data=IDMessageResponse} "Listing created successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidData"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 409 {object} utils.ErrorResponseSwag "marketplace.duplicateTitle"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.createError"
// @Security BearerAuth
// @Router /api/v1/marketplace/listings [post]
func (h *ListingsHandler) CreateListing(c *fiber.Ctx) error {
	// Получаем ID пользователя из контекста
	userID, ok := authMiddleware.GetUserID(c)
	if !ok {
		logger.Warn().Msg("User ID not found in context")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.required")
	}

	var listing models.MarketplaceListing
	if err := c.BodyParser(&listing); err != nil {
		logger.Error().Err(err).Msg("Failed to parse request body")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	// Дополнительная обработка для атрибутов
	var requestBody map[string]interface{}
	if err := json.Unmarshal(c.Body(), &requestBody); err == nil {
		processAttributesFromRequest(requestBody, &listing)
	}

	// Проверяем, была ли категория определена автоматически
	var categoryDetectionStatsID *int32
	var detectedKeywords []string
	detectionLanguage := "ru" // значение по умолчанию
	if requestBody != nil {
		if statsID, ok := requestBody["category_detection_stats_id"].(float64); ok {
			statsIDInt := int32(statsID)
			categoryDetectionStatsID = &statsIDInt
		}
		if keywords, ok := requestBody["detected_keywords"].([]interface{}); ok {
			for _, kw := range keywords {
				if kwStr, ok := kw.(string); ok {
					detectedKeywords = append(detectedKeywords, kwStr)
				}
			}
		}
		// Получаем язык оригинала для правильного обновления счетчиков
		if lang, ok := requestBody["original_language"].(string); ok {
			detectionLanguage = lang
		}
	}

	listing.UserID = userID
	listing.Status = "active"

	// Санитизация полей для защиты от XSS
	listing.Title = utils.SanitizeText(listing.Title)
	listing.Description = utils.SanitizeText(listing.Description)
	if listing.Location != "" {
		listing.Location = utils.SanitizeText(listing.Location)
	}

	// Создаем объявление
	id, err := h.marketplaceService.CreateListing(c.Context(), &listing)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create listing")
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return utils.ErrorResponse(c, fiber.StatusConflict, "marketplace.duplicateTitle")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.createError")
	}

	// Если категория была определена автоматически, обновляем статистику
	if categoryDetectionStatsID != nil {
		go h.updateCategoryDetectionStats(c.Context(), *categoryDetectionStatsID, listing.CategoryID, detectedKeywords, detectionLanguage)
	}

	// Возвращаем ID созданного объявления
	return utils.SuccessResponse(c, IDMessageResponse{
		ID:      id,
		Message: "marketplace.createSuccess",
	})
}

// GetListing получает детали объявления
// @Summary Get listing details
// @Description Returns detailed information about a specific listing including attributes and images
// @Tags marketplace-listings
// @Accept json
// @Produce json
// @Param id path int true "Listing ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.MarketplaceListing} "Listing details"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidId"
// @Failure 404 {object} utils.ErrorResponseSwag "marketplace.notFound"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.getError"
// @Router /api/v1/marketplace/listings/{id} [get]
func (h *ListingsHandler) GetListing(c *fiber.Ctx) error {
	// Получаем ID объявления из параметров URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidId")
	}

	// Получаем язык из query параметра
	lang := c.Query("lang", "en")

	// Создаем контекст с языком
	ctx := context.WithValue(c.Context(), contextKeyLocale, lang)

	var listing *models.MarketplaceListing
	fromCache := false

	// Пробуем получить из кеша, если включен
	if h.cache != nil {
		if cachedData, err := h.cache.GetListingDetails(ctx, id); err == nil {
			// Преобразуем данные из кеша
			if jsonData, err := json.Marshal(cachedData); err == nil {
				var cachedListing models.MarketplaceListing
				if err := json.Unmarshal(jsonData, &cachedListing); err == nil {
					listing = &cachedListing
					fromCache = true
					logger.Debug().Int("listingId", id).Msg("Listing details retrieved from cache")
				}
			}
		}
	}

	// Если не нашли в кеше, получаем из базы
	if !fromCache {
		var err error
		listing, err = h.marketplaceService.GetListingByID(ctx, id)
		if err != nil {
			logger.Error().Err(err).Int("listingId", id).Msg("Failed to get listing")
			if err.Error() == "listing not found" {
				return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.notFound")
			}
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.getError")
		}

		// Сохраняем в кеш для будущих запросов
		if h.cache != nil && listing != nil {
			if err := h.cache.SetListingDetails(ctx, id, listing); err != nil {
				logger.Warn().Err(err).Int("listingId", id).Msg("Failed to cache listing details")
			} else {
				logger.Debug().Int("listingId", id).Msg("Listing details cached successfully")
			}
		}
	}

	// Делаем запрос на увеличение счетчика просмотров в горутине, чтобы не задерживать ответ
	// Создаем новый контекст с данными из текущего запроса
	currentUserID, _ := authMiddleware.GetUserID(c)
	viewCtx := context.WithValue(context.Background(), contextKeyUserID, currentUserID)

	// Получаем IP адрес клиента
	clientIP := c.IP()
	if clientIP == "" {
		// Если c.IP() пустой, пробуем получить из заголовков
		clientIP = c.Get("X-Forwarded-For", "")
		if clientIP == "" {
			clientIP = c.Get("X-Real-IP", "")
		}
		if clientIP == "" {
			// В крайнем случае используем localhost
			clientIP = "127.0.0.1"
		}
	}
	viewCtx = context.WithValue(viewCtx, contextKeyIPAddress, clientIP)

	logger.Debug().Str("clientIP", clientIP).Int("listingId", id).Msg("Incrementing views count")

	go func(listingID int, ctx context.Context) {
		err := h.services.Storage().IncrementViewsCount(ctx, listingID)
		if err != nil {
			logger.Error().Err(err).Int("listingId", listingID).Msg("Failed to increment views count")
		}
	}(id, viewCtx)

	// Загружаем информацию о пользователе из auth-service
	if listing.UserID > 0 {
		userInfo, err := h.services.User().GetUserByID(c.Context(), listing.UserID)
		if err != nil {
			logger.Warn().Err(err).Int("userId", listing.UserID).Msg("Failed to load user info from auth-service")
			// Не прерываем выполнение, просто оставляем пустой User
		} else {
			listing.User = userInfo
		}
	}

	// Получаем ID пользователя из контекста для проверки избранного
	userID, ok := authMiddleware.GetUserID(c)
	if ok && userID > 0 {
		// Проверяем, находится ли объявление в избранном у пользователя
		var favorites []models.MarketplaceListing
		favorites, err = h.marketplaceService.GetUserFavorites(c.Context(), userID)
		if err == nil {
			for _, fav := range favorites {
				if fav.ID == listing.ID {
					listing.IsFavorite = true
					break
				}
			}
		}
	}

	// Возвращаем детали объявления
	return utils.SuccessResponse(c, listing)
}

// GetListingBySlug получает детали объявления по slug
// @Summary Get listing details by slug
// @Description Returns detailed information about a specific listing by URL slug including attributes and images
// @Tags marketplace-listings
// @Accept json
// @Produce json
// @Param slug path string true "Listing slug"
// @Success 200 {object} utils.SuccessResponseSwag{data=models.MarketplaceListing} "Listing details"
// @Failure 404 {object} utils.ErrorResponseSwag "marketplace.notFound"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.getError"
// @Router /api/v1/marketplace/listings/slug/{slug} [get]
func (h *ListingsHandler) GetListingBySlug(c *fiber.Ctx) error {
	// Получаем slug из параметров URL
	slug := c.Params("slug")
	if slug == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidSlug")
	}

	// Получаем язык из query параметра
	lang := c.Query("lang", "en")

	// Создаем контекст с языком
	ctx := context.WithValue(c.Context(), contextKeyLocale, lang)

	// Получаем детали объявления по slug
	listing, err := h.marketplaceService.GetListingBySlug(ctx, slug)
	if err != nil {
		logger.Error().Err(err).Str("slug", slug).Msg("Failed to get listing by slug")
		if strings.Contains(err.Error(), "not found") {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.notFound")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.getError")
	}

	// Делаем запрос на увеличение счетчика просмотров в горутине, чтобы не задерживать ответ
	// Создаем новый контекст с данными из текущего запроса
	currentUserID, _ := authMiddleware.GetUserID(c)
	viewCtx := context.WithValue(context.Background(), contextKeyUserID, currentUserID)

	// Получаем IP адрес клиента
	clientIP := c.IP()
	if clientIP == "" {
		// Если c.IP() пустой, пробуем получить из заголовков
		clientIP = c.Get("X-Forwarded-For", "")
		if clientIP == "" {
			clientIP = c.Get("X-Real-IP", "")
		}
		if clientIP == "" {
			// В крайнем случае используем localhost
			clientIP = "127.0.0.1"
		}
	}
	viewCtx = context.WithValue(viewCtx, contextKeyIPAddress, clientIP)

	logger.Debug().Str("clientIP", clientIP).Int("listingId", listing.ID).Msg("Incrementing views count")

	go func(listingID int, ctx context.Context) {
		err := h.services.Storage().IncrementViewsCount(ctx, listingID)
		if err != nil {
			logger.Error().Err(err).Int("listingId", listingID).Msg("Failed to increment views count")
		}
	}(listing.ID, viewCtx)

	// Получаем ID пользователя из контекста для проверки избранного
	userID, ok := authMiddleware.GetUserID(c)
	if ok && userID > 0 {
		// Проверяем, находится ли объявление в избранном у пользователя
		var favorites []models.MarketplaceListing
		favorites, err = h.marketplaceService.GetUserFavorites(c.Context(), userID)
		if err == nil {
			for _, fav := range favorites {
				if fav.ID == listing.ID {
					listing.IsFavorite = true
					break
				}
			}
		}
	}

	// Возвращаем детали объявления
	return utils.SuccessResponse(c, listing)
}

// GetListings получает список объявлений
// @Summary Get listings list
// @Description Returns paginated list of listings with optional filters
// @Tags marketplace-listings
// @Accept json
// @Produce json
// @Param query query string false "Search query"
// @Param category_id query int false "Category ID filter"
// @Param condition query string false "Condition filter (new, used, etc.)"
// @Param min_price query number false "Minimum price filter"
// @Param max_price query number false "Maximum price filter"
// @Param sort_by query string false "Sort order (price_asc, price_desc, date_desc, etc.)"
// @Param user_id query int false "User ID filter"
// @Param storefront_id query int false "Storefront ID filter"
// @Param exclude_storefronts query boolean false "Exclude storefront products (for admin P2P listings)"
// @Param limit query int false "Number of items per page" default(20)
// @Param offset query int false "Number of items to skip" default(0)
// @Success 200 {object} utils.SuccessResponseSwag{data=ListingsResponse} "Listings list with pagination"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.listError"
// @Router /api/v1/marketplace/listings [get]
func (h *ListingsHandler) GetListings(c *fiber.Ctx) error {
	// Засекаем время начала для измерения производительности
	startTime := time.Now()
	// Получаем параметры фильтрации из запроса
	query := c.Query("q")
	if query == "" {
		query = c.Query("query") // fallback для обратной совместимости
	}
	category := c.Query("category_id")
	condition := c.Query("condition")
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")
	sortBy := c.Query("sort_by")
	userIDStr := c.Query("user_id")
	storefrontIDStr := c.Query("storefront_id")
	excludeStorefronts := c.Query("exclude_storefronts") // Параметр для исключения товаров витрин

	// Значения по умолчанию для пагинации
	limit := 20
	offset := 0

	// Получаем лимит и смещение из запроса
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Формируем фильтры
	filters := make(map[string]string)
	if query != "" {
		filters["q"] = query
	}
	if category != "" {
		filters["category_id"] = category
	}
	if condition != "" {
		filters["condition"] = condition
	}
	if minPrice != "" {
		filters["min_price"] = minPrice
	}
	if maxPrice != "" {
		filters["max_price"] = maxPrice
	}
	if sortBy != "" {
		filters["sort_by"] = sortBy
	}
	if userIDStr != "" {
		filters["user_id"] = userIDStr
	}
	if storefrontIDStr != "" {
		filters["storefront_id"] = storefrontIDStr
	}
	if excludeStorefronts != "" {
		filters["exclude_storefronts"] = excludeStorefronts
	}

	// Получаем список объявлений
	listings, total, err := h.marketplaceService.GetListings(c.Context(), filters, limit, offset)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get listings")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.listError")
	}

	// Проверяем, что listings не nil
	if listings == nil {
		listings = []models.MarketplaceListing{}
	}

	// Загружаем информацию о пользователях из auth-service
	h.loadUserInfoForListings(c.Context(), listings)

	// Асинхронное логирование поискового запроса (если есть query)
	if query != "" && h.services.SearchLogs() != nil {
		// Извлекаем данные из контекста Fiber ДО запуска горутины
		var userID *int
		if uid, ok := authMiddleware.GetUserID(c); ok && uid > 0 {
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
			deviceType := detectDeviceType(userAgent)

			// Парсим цены для логирования
			var priceMin, priceMax *float64
			if minPrice != "" {
				if val, err := strconv.ParseFloat(minPrice, 64); err == nil {
					priceMin = &val
				}
			}
			if maxPrice != "" {
				if val, err := strconv.ParseFloat(maxPrice, 64); err == nil {
					priceMax = &val
				}
			}

			// Преобразуем filters из map[string]string в map[string]interface{}
			filtersInterface := make(map[string]interface{})
			for k, v := range filters {
				filtersInterface[k] = v
			}

			// Преобразуем category в *int
			var categoryID *int
			if category != "" {
				if catID, err := strconv.Atoi(category); err == nil {
					categoryID = &catID
				}
			}

			// Создаем запись лога
			logEntry := &searchlogsTypes.SearchLogEntry{
				Query:           query,
				UserID:          userID,
				SessionID:       sessionID, // Убрали указатель, так как ожидается string
				ResultCount:     int(total),
				ResponseTimeMS:  int64(responseTime),
				Filters:         filtersInterface, // Преобразуем map[string]string в map[string]interface{}
				CategoryID:      categoryID,       // Преобразуем в *int
				PriceMin:        priceMin,
				PriceMax:        priceMax,
				Location:        nil,  // TODO: добавить поддержку локации
				Language:        "ru", // TODO: получать из контекста
				DeviceType:      deviceType,
				UserAgent:       userAgent,
				IP:              ipAddress,
				SearchType:      "listings",
				HasSpellCorrect: false,
				ClickedItems:    []int{},
				Timestamp:       time.Now(),
			}

			// Логируем асинхронно
			if err := h.services.SearchLogs().LogSearch(context.Background(), logEntry); err != nil {
				logger.Error().Err(err).Msg("Failed to log search query")
			}
		}()
	}

	// Возвращаем список объявлений с пагинацией
	return utils.SuccessResponse(c, ListingsResponse{
		Success: true,
		Data:    listings,
		Meta: PaginationMeta{
			Total: int(total),
			Page:  offset/limit + 1,
			Limit: limit,
		},
	})
}

// UpdateListing обновляет существующее объявление
// @Summary Update listing
// @Description Updates an existing marketplace listing. Only the owner can update
// @Tags marketplace-listings
// @Accept json
// @Produce json
// @Param id path int true "Listing ID"
// @Param body body models.MarketplaceListing true "Updated listing data"
// @Success 200 {object} utils.SuccessResponseSwag{data=MessageResponse} "Listing updated successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidData"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 403 {object} utils.ErrorResponseSwag "marketplace.forbidden"
// @Failure 404 {object} utils.ErrorResponseSwag "marketplace.notFound"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.updateError"
// @Security BearerAuth
// @Router /api/v1/marketplace/listings/{id} [put]
func (h *ListingsHandler) UpdateListing(c *fiber.Ctx) error {
	// Получаем ID пользователя из контекста
	userID, ok := authMiddleware.GetUserID(c)
	if !ok {
		logger.Warn().Msg("User ID not found in context")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.required")
	}

	// Получаем ID объявления из параметров URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidId")
	}

	// Получаем текущие данные объявления для проверки владельца
	currentListing, err := h.marketplaceService.GetListingByID(c.Context(), id)
	if err != nil {
		logger.Error().Err(err).Int("listingId", id).Msg("Failed to get listing")
		return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.notFound")
	}

	// Проверяем, является ли пользователь владельцем объявления
	if currentListing.UserID != userID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "marketplace.forbidden")
	}

	// Парсим данные из запроса
	var listing models.MarketplaceListing
	if err := c.BodyParser(&listing); err != nil {
		logger.Error().Err(err).Msg("Failed to parse request body")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	// Дополнительная обработка для атрибутов
	var requestBody map[string]interface{}
	if err := json.Unmarshal(c.Body(), &requestBody); err == nil {
		processAttributesFromRequest(requestBody, &listing)
	}

	// Устанавливаем ID объявления и пользователя
	listing.ID = id
	listing.UserID = userID

	// Санитизация полей для защиты от XSS
	listing.Title = utils.SanitizeText(listing.Title)
	listing.Description = utils.SanitizeText(listing.Description)
	if listing.Location != "" {
		listing.Location = utils.SanitizeText(listing.Location)
	}

	// Обрабатываем изменение цены - если она отличается, сохраняем в историю
	if currentListing.Price != listing.Price {
		// Создаем запись в истории цен
		priceHistory := models.PriceHistoryEntry{
			ListingID:     id,
			Price:         listing.Price,
			EffectiveFrom: time.Now(),
			ChangeSource:  "manual",
		}

		err = h.services.Storage().ClosePriceHistoryEntry(c.Context(), id)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to close previous price history entry")
		}

		err = h.services.Storage().AddPriceHistoryEntry(c.Context(), &priceHistory)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to add price history entry")
		}

		// Проверяем, не является ли изменение цены манипуляцией
		isManipulation, err := h.services.Storage().CheckPriceManipulation(c.Context(), id)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to check price manipulation")
		}

		if isManipulation {
			logger.Warn().Int("listingId", id).Msg("Detected price manipulation")
			// Здесь можно добавить логику для обработки манипуляций с ценой
		}
	}

	// Обновляем объявление
	err = h.marketplaceService.UpdateListing(c.Context(), &listing)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to update listing")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.updateError")
	}

	// Возвращаем успешный результат
	return utils.SuccessResponse(c, MessageResponse{
		Message: "marketplace.updateSuccess",
	})
}

// UpdateListingStatus обновляет статус объявления
// @Summary Update listing status
// @Description Updates the status of a marketplace listing (active/inactive). Only the owner can update
// @Tags marketplace-listings
// @Accept json
// @Produce json
// @Param id path int true "Listing ID"
// @Param body body map[string]string true "Status update data"
// @Success 200 {object} utils.SuccessResponseSwag{data=MessageResponse} "Status updated successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidData"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 403 {object} utils.ErrorResponseSwag "marketplace.forbidden"
// @Failure 404 {object} utils.ErrorResponseSwag "marketplace.notFound"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.updateError"
// @Security BearerAuth
// @Router /api/v1/marketplace/listings/{id}/status [patch]
func (h *ListingsHandler) UpdateListingStatus(c *fiber.Ctx) error {
	// Получаем ID пользователя из контекста
	userID, ok := authMiddleware.GetUserID(c)
	if !ok {
		logger.Warn().Msg("User ID not found in context")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.required")
	}

	// Получаем ID объявления из параметров URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidId")
	}

	// Парсим тело запроса
	var request struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&request); err != nil {
		logger.Error().Err(err).Msg("Failed to parse request body")
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	// Проверяем валидность статуса
	if request.Status != "active" && request.Status != "inactive" && request.Status != "draft" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidStatus")
	}

	// Получаем существующее объявление для проверки прав
	existingListing, err := h.marketplaceService.GetListingByID(c.Context(), id)
	if err != nil {
		logger.Error().Err(err).Int("listingId", id).Msg("Failed to get existing listing")
		return utils.ErrorResponse(c, fiber.StatusNotFound, "marketplace.notFound")
	}

	// Проверяем, что пользователь является владельцем объявления
	if existingListing.UserID != userID {
		logger.Warn().
			Int("userID", userID).
			Int("ownerID", existingListing.UserID).
			Int("listingID", id).
			Msg("User tried to update listing status they don't own")
		return utils.ErrorResponse(c, fiber.StatusForbidden, "marketplace.forbidden")
	}

	// Обновляем только статус
	existingListing.Status = request.Status

	// Сохраняем изменения
	err = h.marketplaceService.UpdateListing(c.Context(), existingListing)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to update listing status")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.updateError")
	}

	logger.Info().
		Int("listingId", id).
		Int("userId", userID).
		Str("newStatus", request.Status).
		Msg("Listing status updated successfully")

	// Возвращаем успешный результат
	return utils.SuccessResponse(c, MessageResponse{
		Message: "marketplace.statusUpdateSuccess",
	})
}

// CheckSlugAvailability проверяет доступность slug
// @Summary Check slug availability
// @Description Checks if a slug is available for use and suggests alternatives if not
// @Tags marketplace-listings
// @Accept json
// @Produce json
// @Param body body map[string]interface{} true "Slug to check"
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "Slug availability status"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidData"
// @Router /api/v1/marketplace/listings/check-slug [post]
func (h *ListingsHandler) CheckSlugAvailability(c *fiber.Ctx) error {
	var request struct {
		Slug      string `json:"slug"`
		ExcludeID int    `json:"exclude_id,omitempty"`
	}

	if err := c.BodyParser(&request); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidData")
	}

	if request.Slug == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.slugRequired")
	}

	// Проверяем доступность slug
	isAvailable, err := h.marketplaceService.IsSlugAvailable(c.Context(), request.Slug, request.ExcludeID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to check slug availability")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.checkError")
	}

	response := map[string]interface{}{
		"available": isAvailable,
		"slug":      request.Slug,
	}

	// Если slug занят, генерируем альтернативы
	if !isAvailable {
		suggestions, err := h.marketplaceService.GenerateUniqueSlug(c.Context(), request.Slug, request.ExcludeID)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to generate slug suggestions")
			// Не прерываем выполнение, просто не возвращаем предложения
		} else {
			response["suggestion"] = suggestions
		}
	}

	return utils.SuccessResponse(c, response)
}

// DeleteListing удаляет объявление
// @Summary Delete listing
// @Description Deletes a marketplace listing. Only the owner can delete
// @Tags marketplace-listings
// @Accept json
// @Produce json
// @Param id path int true "Listing ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=MessageResponse} "Listing deleted successfully"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidId"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 403 {object} utils.ErrorResponseSwag "marketplace.forbidden"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.deleteError"
// @Security BearerAuth
// @Router /api/v1/marketplace/listings/{id} [delete]
func (h *ListingsHandler) DeleteListing(c *fiber.Ctx) error {
	// Получаем ID пользователя из контекста
	userID, ok := authMiddleware.GetUserID(c)
	if !ok {
		logger.Warn().Msg("User ID not found in context")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.required")
	}

	// Получаем полный профиль пользователя с ролью
	isAdmin := false
	userProfile, err := h.services.User().GetUserProfile(c.Context(), userID)
	if err == nil && userProfile != nil {
		// Проверяем роль пользователя - админ или супер админ
		// Флаг IsAdmin уже установлен в GetUserProfile на основе роли
		if userProfile.IsAdmin {
			isAdmin = true
			roleName := ""
			if userProfile.Role != nil {
				roleName = userProfile.Role.Name
			}
			logger.Info().
				Int("userID", userID).
				Str("email", userProfile.Email).
				Bool("isAdmin", userProfile.IsAdmin).
				Str("role", roleName).
				Msg("Admin access for deletion")
		}
	} else {
		logger.Error().Err(err).Int("userID", userID).Msg("Failed to get user profile")
	}

	// Получаем ID объявления из параметров URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidId")
	}

	// Удаляем объявление (передаем флаг администратора)
	err = h.marketplaceService.DeleteListingWithAdmin(c.Context(), id, userID, isAdmin)
	if err != nil {
		logger.Error().Err(err).Int("listingId", id).Msg("Failed to delete listing")
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "permission") {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "marketplace.forbidden")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.deleteError")
	}

	// Удаляем документ из OpenSearch
	go func() {
		err := h.services.Storage().DeleteListingIndex(context.Background(), fmt.Sprintf("%d", id))
		if err != nil {
			logger.Error().Err(err).Int("listingId", id).Msg("Failed to delete listing index")
		}
	}()

	// Возвращаем успешный результат
	return utils.SuccessResponse(c, MessageResponse{
		Message: "marketplace.deleteSuccess",
	})
}

// GetPriceHistory получает историю цен для объявления
// @Summary Get listing price history
// @Description Returns price history for a specific listing
// @Tags marketplace-listings
// @Accept json
// @Produce json
// @Param id path int true "Listing ID"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]models.PriceHistoryEntry} "Price history entries"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidId"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.priceHistoryError"
// @Router /api/v1/marketplace/listings/{id}/price-history [get]
func (h *ListingsHandler) GetPriceHistory(c *fiber.Ctx) error {
	// Получаем ID объявления из параметров URL
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidId")
	}

	// Получаем историю цен
	priceHistory, err := h.marketplaceService.GetPriceHistory(c.Context(), id)
	if err != nil {
		logger.Error().Err(err).Int("listingId", id).Msg("Failed to get price history")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.priceHistoryError")
	}

	// Проверяем, что priceHistory не nil
	if priceHistory == nil {
		priceHistory = []models.PriceHistoryEntry{}
	}

	// Возвращаем историю цен
	return utils.SuccessResponse(c, priceHistory)
}

// GetMyListings returns current user's listings
// @Summary Get my listings
// @Description Returns all listings created by the authenticated user
// @Tags marketplace
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(20)
// @Param offset query int false "Offset" default(0)
// @Param status query string false "Filter by status (active, sold, draft)"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]models.MarketplaceListing} "User's listings"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.listError"
// @Security BearerAuth
// @Router /api/v1/marketplace/my-listings [get]
func (h *ListingsHandler) GetMyListings(c *fiber.Ctx) error {
	// Получаем ID текущего пользователя из контекста
	userID, ok := authMiddleware.GetUserID(c)
	if !ok {
		logger.Warn().Msg("User ID not found in context for my-listings")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.required")
	}

	// Значения по умолчанию для пагинации
	limit := 20
	offset := 0

	// Получаем лимит и смещение из запроса
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Формируем фильтры
	filters := map[string]string{
		"user_id": strconv.Itoa(userID),
	}

	// Добавляем фильтр по статусу если указан
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}

	// Получаем список объявлений пользователя
	listings, total, err := h.marketplaceService.GetListings(c.Context(), filters, limit, offset)
	if err != nil {
		logger.Error().Err(err).Int("userId", userID).Msg("Failed to get user listings")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.listError")
	}

	// Проверяем, что listings не nil
	if listings == nil {
		listings = []models.MarketplaceListing{}
	}

	// Загружаем информацию о пользователе из auth-service
	h.loadUserInfoForListings(c.Context(), listings)

	// Возвращаем объявления с общим количеством
	return c.JSON(fiber.Map{
		"success": true,
		"data":    listings,
		"total":   total,
	})
}

// SynchronizeDiscounts синхронизирует данные о скидках
// @Summary Synchronize discount data
// @Description Synchronizes discount data for all listings (Admin only)
// @Tags marketplace-admin
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=MessageResponse} "Discounts synchronized successfully"
// @Failure 401 {object} utils.ErrorResponseSwag "auth.required"
// @Failure 403 {object} utils.ErrorResponseSwag "admin.required"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.syncError"
// @Security BearerAuth
// @Router /api/v1/admin/sync-discounts [post]
func (h *ListingsHandler) SynchronizeDiscounts(c *fiber.Ctx) error {
	// Проверяем, является ли пользователь администратором
	userID, ok := authMiddleware.GetUserID(c)
	if !ok {
		logger.Warn().Msg("User ID not found in context")
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.required")
	}

	// Получаем пользователя для проверки email
	user, err := h.services.User().GetUserByID(c.Context(), userID)
	if err != nil {
		logger.Error().Err(err).Int("userId", userID).Msg("Failed to get user")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "admin.checkError")
	}

	// Проверяем права администратора
	isAdmin, err := h.services.User().IsUserAdmin(c.Context(), user.Email)
	if err != nil || !isAdmin {
		logger.Error().Err(err).Int("userId", userID).Msg("User is not admin")
		return utils.ErrorResponse(c, fiber.StatusForbidden, "admin.required")
	}

	// Запускаем синхронизацию
	err = h.marketplaceService.SynchronizeDiscountData(c.Context(), 0)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to synchronize discount data")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.syncError")
	}

	// Возвращаем успешный результат
	return utils.SuccessResponse(c, MessageResponse{
		Message: "marketplace.syncSuccess",
	})
}

// processAttributesFromRequest обрабатывает атрибуты из запроса
func processAttributesFromRequest(requestBody map[string]interface{}, listing *models.MarketplaceListing) {
	// Обработка переводов
	if translationsRaw, ok := requestBody["translations"]; ok {
		if translationsMap, ok := translationsRaw.(map[string]interface{}); ok {
			translations := make(models.TranslationMap)
			for lang, fieldsRaw := range translationsMap {
				if fields, ok := fieldsRaw.(map[string]interface{}); ok {
					langTranslations := make(map[string]string)
					for field, value := range fields {
						if strValue, ok := value.(string); ok {
							langTranslations[field] = strValue
						}
					}
					if len(langTranslations) > 0 {
						translations[lang] = langTranslations
					}
				}
			}
			if len(translations) > 0 {
				listing.Translations = translations
			}
		}
	}

	// Проверяем наличие атрибутов в запросе
	if attributesRaw, ok := requestBody["attributes"]; ok {
		if attributesSlice, ok := attributesRaw.([]interface{}); ok {
			var attributes []models.ListingAttributeValue

			for _, attrRaw := range attributesSlice {
				if attrMap, ok := attrRaw.(map[string]interface{}); ok {
					var attr models.ListingAttributeValue

					// Перенос всех полей из JSON-объекта
					if id, ok := attrMap["attribute_id"].(float64); ok {
						attr.AttributeID = int(id)
					}
					if name, ok := attrMap["attribute_name"].(string); ok {
						attr.AttributeName = name
					}
					if displayName, ok := attrMap["display_name"].(string); ok {
						attr.DisplayName = displayName
					}
					if attrType, ok := attrMap["attribute_type"].(string); ok {
						attr.AttributeType = attrType
					}
					if unit, ok := attrMap["unit"].(string); ok {
						attr.Unit = unit
					}
					if displayValue, ok := attrMap["display_value"].(string); ok {
						attr.DisplayValue = displayValue
					}

					// Обрабатываем значение в зависимости от типа атрибута
					switch attr.AttributeType {
					case "text", "select":
						if textValue, ok := attrMap["text_value"].(string); ok && textValue != "" {
							attr.TextValue = &textValue
						} else if textValue, ok := attrMap["value"].(string); ok && textValue != "" {
							attr.TextValue = &textValue
						}
					case "number":
						if numValue, ok := attrMap["numeric_value"].(float64); ok {
							attr.NumericValue = &numValue
						} else if numValue, ok := attrMap["value"].(float64); ok {
							attr.NumericValue = &numValue
						} else if textValue, ok := attrMap["value"].(string); ok && textValue != "" {
							// Иногда числа приходят как строки, преобразуем
							if numVal, err := strconv.ParseFloat(textValue, 64); err == nil {
								attr.NumericValue = &numVal
							}
						}
					case "boolean":
						if boolValue, ok := attrMap["boolean_value"].(bool); ok {
							attr.BooleanValue = &boolValue
						} else if boolValue, ok := attrMap["value"].(bool); ok {
							attr.BooleanValue = &boolValue
						}
					case "multiselect":
						// Для multiselect значение хранится в JSON
						if jsonValues, ok := attrMap["json_value"]; ok {
							jsonBytes, err := json.Marshal(jsonValues)
							if err == nil {
								attr.JSONValue = jsonBytes
							}
						} else if jsonValues, ok := attrMap["value"]; ok {
							jsonBytes, err := json.Marshal(jsonValues)
							if err == nil {
								attr.JSONValue = jsonBytes
							}
						}
					}

					attributes = append(attributes, attr)
				}
			}

			listing.Attributes = attributes
		}
	}

	// Обработка вариантов товара
	if variantsRaw, ok := requestBody["productVariants"]; ok {
		if variantsSlice, ok := variantsRaw.([]interface{}); ok {
			var variants []models.MarketplaceListingVariant

			for _, variantRaw := range variantsSlice {
				if variantMap, ok := variantRaw.(map[string]interface{}); ok {
					var variant models.MarketplaceListingVariant

					// Обработка основных полей варианта
					if sku, ok := variantMap["sku"].(string); ok {
						variant.SKU = sku
					}
					if price, ok := variantMap["price"].(float64); ok {
						variant.Price = &price
					}
					if stock, ok := variantMap["stock"].(float64); ok {
						stockInt := int(stock)
						variant.Stock = &stockInt
					}
					if imageURL, ok := variantMap["image"].(string); ok && imageURL != "" {
						variant.ImageURL = &imageURL
					}

					// Обработка атрибутов варианта
					if attributesRaw, ok := variantMap["attributes"]; ok {
						if attributesMap, ok := attributesRaw.(map[string]interface{}); ok {
							variant.Attributes = make(map[string]string)
							for key, value := range attributesMap {
								if strValue, ok := value.(string); ok {
									variant.Attributes[key] = strValue
								}
							}
						}
					}

					variant.IsActive = true
					variants = append(variants, variant)
				}
			}

			if len(variants) > 0 {
				listing.Variants = variants
			}
		}
	}
}

// detectDeviceType определяет тип устройства по User-Agent
func detectDeviceType(userAgent string) string {
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

// updateCategoryDetectionStats обновляет статистику определения категории при создании объявления
func (h *ListingsHandler) updateCategoryDetectionStats(ctx context.Context, statsID int32, categoryID int, detectedKeywords []string, language string) {
	// Получаем storage для работы с БД
	storage := h.services.Storage()

	// Преобразуем storage к конкретному типу для доступа к методам
	if db, ok := storage.(*postgres.Database); ok {
		// Обновляем статистику как подтвержденную
		statsRepo := postgres.NewCategoryDetectionStatsRepository(db.GetSQLXDB())

		// Помечаем, что пользователь подтвердил категорию
		confirmed := true
		finalCategoryID := int32(categoryID) //nolint:gosec // CategoryID проверяется выше

		err := statsRepo.UpdateUserFeedback(ctx, statsID, confirmed, &finalCategoryID)
		if err != nil {
			logger.Error().Err(err).Int32("statsID", statsID).Msg("Failed to update category detection stats")
			return
		}

		// Обновляем success_rate для использованных ключевых слов
		if len(detectedKeywords) > 0 {
			keywordRepo := postgres.NewCategoryKeywordRepository(db.GetSQLXDB())

			// Увеличиваем счетчик использования для найденных ключевых слов
			err := keywordRepo.IncrementUsageCount(ctx, finalCategoryID, detectedKeywords, language)
			if err != nil {
				logger.Error().Err(err).
					Int32("categoryID", finalCategoryID).
					Interface("keywords", detectedKeywords).
					Msg("Failed to increment keyword usage count")
			} else {
				logger.Info().
					Int32("categoryID", finalCategoryID).
					Interface("keywords", detectedKeywords).
					Msg("Successfully incremented keyword usage count")
			}

			// Получаем все статистики для пересчета success_rate
			stats, err := statsRepo.GetRecentStats(ctx, 30) // за последние 30 дней
			if err != nil {
				logger.Error().Err(err).Msg("Failed to get recent stats for success rate update")
				return
			}

			// Подсчитываем успешность для каждого ключевого слова
			keywordSuccess := make(map[string]int)
			keywordTotal := make(map[string]int)

			logger.Info().Int("totalStats", len(stats)).Msg("Processing stats for success rate calculation")

			for _, stat := range stats {
				for _, keyword := range stat.MatchedKeywords {
					keywordTotal[keyword]++
					if stat.UserConfirmed != nil && *stat.UserConfirmed {
						keywordSuccess[keyword]++
					}
				}
			}

			logger.Info().
				Interface("keywordTotal", keywordTotal).
				Interface("keywordSuccess", keywordSuccess).
				Msg("Keyword statistics calculated")

			// Обновляем success_rate для каждого ключевого слова
			for keyword, total := range keywordTotal {
				if total > 0 {
					successRate := float64(keywordSuccess[keyword]) / float64(total)
					logger.Info().
						Str("keyword", keyword).
						Int("success", keywordSuccess[keyword]).
						Int("total", total).
						Float64("successRate", successRate).
						Msg("Updating keyword success rate")
					err := keywordRepo.UpdateSuccessRate(ctx, keyword, successRate)
					if err != nil {
						logger.Error().Err(err).Str("keyword", keyword).Msg("Failed to update keyword success rate")
					}
				}
			}
		}

		logger.Info().
			Int32("statsID", statsID).
			Int("categoryID", categoryID).
			Strs("keywords", detectedKeywords).
			Msg("Category detection stats updated successfully")
	}
}

// AdminStatisticsResponse represents statistics for admin dashboard
type AdminStatisticsResponse struct {
	Total   int `json:"total"`
	Active  int `json:"active"`
	Pending int `json:"pending"`
	Views   int `json:"views"`
}

// GetAdminStatistics возвращает статистику для админ панели
// @Summary Get admin statistics
// @Description Returns listing statistics for admin dashboard
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.SuccessResponseSwag{data=handler.AdminStatisticsResponse} "Statistics"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/admin/listings/statistics [get]
func (h *ListingsHandler) GetAdminStatistics(c *fiber.Ctx) error {
	ctx := context.Background()

	// Проверяем, что пользователь администратор
	isAdmin := authMiddleware.IsAdmin(c)
	if !isAdmin {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "admin.unauthorized")
	}

	// Получаем общую статистику
	db, ok := h.services.Storage().(*postgres.Database)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.internal")
	}

	var stats AdminStatisticsResponse

	// Получаем общее количество объявлений (только P2P, исключаем товары витрин)
	err := db.GetPool().QueryRow(ctx, `
		SELECT COUNT(*) FROM marketplace_listings WHERE storefront_id IS NULL
	`).Scan(&stats.Total)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get total listings count")
	}

	// Получаем количество активных объявлений (только P2P)
	err = db.GetPool().QueryRow(ctx, `
		SELECT COUNT(*) FROM marketplace_listings WHERE status = 'active' AND storefront_id IS NULL
	`).Scan(&stats.Active)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get active listings count")
	}

	// Получаем количество ожидающих модерации (только P2P)
	err = db.GetPool().QueryRow(ctx, `
		SELECT COUNT(*) FROM marketplace_listings WHERE status = 'pending' AND storefront_id IS NULL
	`).Scan(&stats.Pending)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get pending listings count")
	}

	// Получаем количество просмотров за последние 30 дней
	err = db.GetPool().QueryRow(ctx, `
		SELECT COUNT(*) FROM listing_views
		WHERE view_time > NOW() - INTERVAL '30 days'
	`).Scan(&stats.Views)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get views count")
	}

	return utils.SuccessResponse(c, stats)
}
