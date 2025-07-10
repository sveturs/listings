// backend/internal/proj/global/handler/unified_search.go
package handler

import (
	"context"
	"math"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"backend/internal/domain/behavior"
	"backend/internal/domain/models"
	"backend/internal/domain/search"
	"backend/internal/logger"
	globalService "backend/internal/proj/global/service"
	storefrontOpenSearch "backend/internal/proj/storefronts/storage/opensearch"
	"backend/pkg/utils"
)

// UnifiedSearchHandler обрабатывает унифицированные поисковые запросы
type UnifiedSearchHandler struct {
	services globalService.ServicesInterface
}

// NewUnifiedSearchHandler создает новый обработчик унифицированного поиска
func NewUnifiedSearchHandler(services globalService.ServicesInterface) *UnifiedSearchHandler {
	return &UnifiedSearchHandler{
		services: services,
	}
}

// UnifiedSearchParams параметры унифицированного поиска
type UnifiedSearchParams struct {
	Query            string                 `json:"query" form:"query"`
	ProductTypes     []string               `json:"product_types" form:"product_types"` // ["marketplace", "storefront"]
	Page             int                    `json:"page" form:"page"`
	Limit            int                    `json:"limit" form:"limit"`
	CategoryID       string                 `json:"category_id" form:"category_id"`
	PriceMin         float64                `json:"price_min" form:"price_min"`
	PriceMax         float64                `json:"price_max" form:"price_max"`
	SortBy           string                 `json:"sort_by" form:"sort_by"`
	SortOrder        string                 `json:"sort_order" form:"sort_order"`
	AttributeFilters map[string]interface{} `json:"attribute_filters" form:"attribute_filters"`
	StorefrontID     int                    `json:"storefront_id" form:"storefront_id"`
	City             string                 `json:"city" form:"city"`
	Language         string                 `json:"language" form:"language"`
}

// UnifiedSearchResult результат унифицированного поиска
type UnifiedSearchResult struct {
	Items      []UnifiedSearchItem    `json:"items"`
	Total      int                    `json:"total"`
	Page       int                    `json:"page"`
	Limit      int                    `json:"limit"`
	TotalPages int                    `json:"total_pages"`
	HasMore    bool                   `json:"has_more"`
	TookMs     int64                  `json:"took_ms"`
	Facets     map[string]interface{} `json:"facets,omitempty"`
}

// UnifiedSearchItem унифицированный элемент поиска
type UnifiedSearchItem struct {
	ID          string                 `json:"id"`           // Уникальный ID (ml_123 или sp_456)
	ProductType string                 `json:"product_type"` // "marketplace" или "storefront"
	ProductID   int                    `json:"product_id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Price       float64                `json:"price"`
	Currency    string                 `json:"currency"`
	Images      []UnifiedProductImage  `json:"images"`
	Category    UnifiedCategoryInfo    `json:"category"`
	Location    *UnifiedLocationInfo   `json:"location,omitempty"`
	Storefront  *UnifiedStorefrontInfo `json:"storefront,omitempty"` // Только для storefront товаров
	Score       float64                `json:"score"`
	Highlights  map[string][]string    `json:"highlights,omitempty"`
	ViewsCount  int                    `json:"views_count,omitempty"` // Для расчета популярности
	CreatedAt   *time.Time             `json:"created_at,omitempty"`  // Для расчета свежести
}

// UnifiedProductImage изображение товара
type UnifiedProductImage struct {
	URL     string `json:"url"`
	AltText string `json:"alt_text,omitempty"`
	IsMain  bool   `json:"is_main"`
}

// UnifiedCategoryInfo информация о категории
type UnifiedCategoryInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug,omitempty"`
}

// UnifiedLocationInfo информация о местоположении
type UnifiedLocationInfo struct {
	City    string  `json:"city,omitempty"`
	Country string  `json:"country,omitempty"`
	Lat     float64 `json:"lat,omitempty"`
	Lng     float64 `json:"lng,omitempty"`
}

// UnifiedStorefrontInfo информация о витрине
type UnifiedStorefrontInfo struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Slug       string  `json:"slug,omitempty"`
	Rating     float64 `json:"rating,omitempty"`
	IsVerified bool    `json:"is_verified"`
}

// UnifiedSearch выполняет унифицированный поиск по marketplace и storefront товарам
// @Summary Unified search across all product types
// @Description Searches both marketplace listings and storefront products
// @Tags search
// @Accept json
// @Produce json
// @Param body body UnifiedSearchParams false "Search parameters"
// @Param query query string false "Search query"
// @Param product_types query []string false "Product types to search" Enums(marketplace,storefront)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param category_id query string false "Category ID"
// @Param price_min query number false "Minimum price"
// @Param price_max query number false "Maximum price"
// @Param sort_by query string false "Sort field" Enums(relevance,price,date,popularity)
// @Param sort_order query string false "Sort order" Enums(asc,desc)
// @Param storefront_id query int false "Storefront ID filter"
// @Param city query string false "City filter"
// @Param language query string false "Language" default(ru)
// @Success 200 {object} UnifiedSearchResult "Search results"
// @Failure 400 {object} utils.ErrorResponseSwag "search.invalidParams"
// @Failure 500 {object} utils.ErrorResponseSwag "search.searchError"
// @Router /api/v1/search [get]
func (h *UnifiedSearchHandler) UnifiedSearch(c *fiber.Ctx) error {
	ctx := c.Context()
	startTime := time.Now()

	// Парсим параметры поиска
	var params UnifiedSearchParams

	// Сначала пытаемся получить из JSON body
	if c.Get("Content-Type") == "application/json" {
		if err := c.BodyParser(&params); err != nil {
			logger.Debug().Err(err).Msg("Failed to parse JSON body, trying query params")
		}
	}

	// Получаем параметры из query string (перезаписывают JSON если есть)
	// Поддерживаем оба варианта: "q" и "query"
	if query := c.Query("q"); query != "" {
		params.Query = query
	} else if query := c.Query("query"); query != "" {
		params.Query = query
	}
	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			params.Page = p
		}
	}
	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			params.Limit = l
		}
	}
	if categoryID := c.Query("category_id"); categoryID != "" {
		params.CategoryID = categoryID
	}
	if priceMin := c.Query("price_min"); priceMin != "" {
		if p, err := strconv.ParseFloat(priceMin, 64); err == nil && p >= 0 {
			params.PriceMin = p
		}
	}
	if priceMax := c.Query("price_max"); priceMax != "" {
		if p, err := strconv.ParseFloat(priceMax, 64); err == nil && p >= 0 {
			params.PriceMax = p
		}
	}
	if sortBy := c.Query("sort_by"); sortBy != "" {
		params.SortBy = sortBy
	}
	if sortOrder := c.Query("sort_order"); sortOrder != "" {
		params.SortOrder = sortOrder
	}
	if storefrontID := c.Query("storefront_id"); storefrontID != "" {
		if s, err := strconv.Atoi(storefrontID); err == nil && s > 0 {
			params.StorefrontID = s
		}
	}
	if city := c.Query("city"); city != "" {
		params.City = city
	}
	if lang := c.Query("language"); lang != "" {
		params.Language = lang
	}

	// Обработка product_types
	if productTypes := c.Query("product_types"); productTypes != "" {
		// Парсинг comma-separated значений
		params.ProductTypes = strings.Split(productTypes, ",")
		// Очистка пробелов
		for i, pt := range params.ProductTypes {
			params.ProductTypes[i] = strings.TrimSpace(pt)
		}
	}

	// Установка значений по умолчанию
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 20
	}
	if params.Limit > 100 {
		params.Limit = 100
	}
	if params.Language == "" {
		params.Language = "ru"
	}
	if len(params.ProductTypes) == 0 {
		params.ProductTypes = []string{"marketplace", "storefront"}
	}

	// Выполняем поиск
	result, err := h.performUnifiedSearch(ctx, &params)
	if err != nil {
		logger.Error().Err(err).Msg("Unified search failed")
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "search.searchError")
	}

	// Сохраняем поисковый запрос при успешном поиске
	if params.Query != "" && result.Total > 0 {
		if err := h.services.Marketplace().SaveSearchQuery(ctx, params.Query, result.Total, params.Language); err != nil {
			logger.Error().Err(err).Msg("Failed to save search query")
			// Не возвращаем ошибку пользователю, так как основной поиск прошел успешно
		}
	}

	// Трекинг поискового события (только для первой страницы)
	if params.Query != "" && params.Page == 1 {
		// Извлекаем все необходимые данные из fiber.Ctx перед запуском горутины
		trackCtx := &trackingContext{
			userAgent: c.Get("User-Agent"),
			referer:   c.Get("Referer"),
			ipAddress: c.IP(),
		}

		// Получаем userID если пользователь авторизован
		if c.Locals("userID") != nil {
			if uid, ok := c.Locals("userID").(int); ok && uid > 0 {
				trackCtx.userID = &uid
			}
		}

		// Получаем session_id из заголовков или query параметров
		sessionID := c.Get("X-Session-ID")
		if sessionID == "" {
			sessionID = c.Query("session_id")
		}
		trackCtx.sessionID = sessionID

		go h.trackSearchEvent(trackCtx, &params, result, time.Since(startTime))
	}

	return c.JSON(result)
}

// performUnifiedSearch выполняет поиск по всем типам товаров
func (h *UnifiedSearchHandler) performUnifiedSearch(ctx context.Context, params *UnifiedSearchParams) (*UnifiedSearchResult, error) {
	var allItems []UnifiedSearchItem
	var totalCount int
	var tookMs int64

	// Поиск в marketplace (если включен)
	if h.containsProductType(params.ProductTypes, "marketplace") {
		marketplaceItems, count, took, err := h.searchMarketplace(ctx, params)
		if err != nil {
			logger.Error().Err(err).Msg("Marketplace search failed")
		} else {
			allItems = append(allItems, marketplaceItems...)
			totalCount += count
			tookMs += took
		}
	}

	// Поиск в storefront (если включен)
	if h.containsProductType(params.ProductTypes, "storefront") {
		storefrontItems, count, took, err := h.searchStorefront(ctx, params)
		if err != nil {
			logger.Error().Err(err).Msg("Storefront search failed")
		} else {
			allItems = append(allItems, storefrontItems...)
			totalCount += count
			tookMs += took
		}
	}

	// Объединяем и ранжируем результаты
	rankedItems := h.mergeAndRankResults(allItems, params)

	// Применяем пагинацию к объединенным результатам
	offset := (params.Page - 1) * params.Limit
	end := offset + params.Limit

	var pagedItems []UnifiedSearchItem
	if offset >= len(rankedItems) {
		pagedItems = []UnifiedSearchItem{}
	} else if end > len(rankedItems) {
		pagedItems = rankedItems[offset:]
	} else {
		pagedItems = rankedItems[offset:end]
	}

	// Вычисляем метаданные
	totalPages := int(math.Ceil(float64(totalCount) / float64(params.Limit)))
	hasMore := params.Page < totalPages

	return &UnifiedSearchResult{
		Items:      pagedItems,
		Total:      totalCount,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
		HasMore:    hasMore,
		TookMs:     tookMs,
	}, nil
}

// searchMarketplace поиск в marketplace
func (h *UnifiedSearchHandler) searchMarketplace(ctx context.Context, params *UnifiedSearchParams) ([]UnifiedSearchItem, int, int64, error) {
	// Конвертируем параметры в формат для marketplace поиска
	searchParams := &search.ServiceParams{
		Query:         params.Query,
		Page:          params.Page,
		Size:          params.Limit * 2, // Берем больше для лучшего смешивания результатов
		CategoryID:    params.CategoryID,
		PriceMin:      params.PriceMin,
		PriceMax:      params.PriceMax,
		Sort:          params.SortBy,
		SortDirection: params.SortOrder,
		// AttributeFilters: params.AttributeFilters, // TODO: конвертировать типы
		Language: params.Language,
	}

	results, err := h.services.Marketplace().SearchListingsAdvanced(ctx, searchParams)
	if err != nil {
		return nil, 0, 0, err
	}

	// Конвертируем результаты в унифицированный формат
	items := make([]UnifiedSearchItem, 0, len(results.Items))
	for _, listing := range results.Items {
		if listing == nil {
			continue
		}

		item := UnifiedSearchItem{
			ID:          "ml_" + strconv.Itoa(listing.ID),
			ProductType: "marketplace",
			ProductID:   listing.ID,
			Name:        listing.Title,
			Description: listing.Description,
			Price:       listing.Price,
			Currency:    "RSD", // TODO: получить реальную валюту из конфигурации
			Images:      h.convertMarketplaceImages(listing.Images),
			Category:    h.convertMarketplaceCategory(listing.Category),
			Location:    h.convertMarketplaceLocation(listing),
			Score:       1.0, // TODO: получить реальный score из OpenSearch
			ViewsCount:  listing.ViewsCount,
			CreatedAt:   &listing.CreatedAt,
		}

		items = append(items, item)
	}

	return items, results.Total, results.Took, nil
}

// searchStorefront поиск в storefront
func (h *UnifiedSearchHandler) searchStorefront(ctx context.Context, params *UnifiedSearchParams) ([]UnifiedSearchItem, int, int64, error) {
	// Получаем репозиторий поиска товаров витрин
	searchRepo := h.services.Storage().StorefrontProductSearch()
	if searchRepo == nil {
		logger.Debug().Msg("Storefront product search repository not configured")
		return []UnifiedSearchItem{}, 0, 0, nil
	}

	productSearchRepo, ok := searchRepo.(storefrontOpenSearch.ProductSearchRepository)
	if !ok {
		logger.Error().Msg("Invalid storefront product search repository type")
		return []UnifiedSearchItem{}, 0, 0, nil
	}

	// Конвертируем параметры в формат для поиска товаров витрин
	categoryID := 0
	if params.CategoryID != "" {
		// Попробуем конвертировать в int
		if id, err := strconv.Atoi(params.CategoryID); err == nil {
			categoryID = id
		}
	}

	searchParams := &storefrontOpenSearch.ProductSearchParams{
		Query:        params.Query,
		StorefrontID: params.StorefrontID,
		CategoryID:   categoryID,
		PriceMin:     params.PriceMin,
		PriceMax:     params.PriceMax,
		City:         params.City,
		Limit:        params.Limit * 2, // Берем больше для лучшего смешивания результатов
		Offset:       (params.Page - 1) * params.Limit * 2,
		SortBy:       params.SortBy,
		SortOrder:    params.SortOrder,
	}

	// Выполняем поиск
	results, err := productSearchRepo.SearchProducts(ctx, searchParams)
	if err != nil {
		logger.Error().Err(err).Msg("Storefront product search failed")
		return nil, 0, 0, err
	}

	// Конвертируем результаты в унифицированный формат
	items := make([]UnifiedSearchItem, 0, len(results.Products))
	for _, product := range results.Products {
		if product == nil {
			continue
		}

		// Преобразуем изображения
		images := make([]UnifiedProductImage, 0, len(product.Images))
		for _, img := range product.Images {
			images = append(images, UnifiedProductImage{
				URL:     img.URL,
				AltText: img.AltText,
				IsMain:  img.IsMain,
			})
		}

		// Создаем унифицированный элемент
		item := UnifiedSearchItem{
			ID:          "sp_" + strconv.Itoa(product.ProductID),
			ProductType: "storefront",
			ProductID:   product.ProductID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Currency:    product.Currency,
			Images:      images,
			Category: UnifiedCategoryInfo{
				ID:   product.Category.ID,
				Name: product.Category.Name,
				Slug: product.Category.Slug,
			},
			Score: product.Score,
		}

		// Добавляем информацию о локации, если есть
		if product.Storefront.City != "" || product.Storefront.Country != "" {
			item.Location = &UnifiedLocationInfo{
				City:    product.Storefront.City,
				Country: product.Storefront.Country,
			}
		}

		// Добавляем информацию о витрине
		if product.StorefrontID > 0 {
			item.Storefront = &UnifiedStorefrontInfo{
				ID:         product.StorefrontID,
				Name:       product.Storefront.Name,
				Slug:       product.Storefront.Slug,
				Rating:     product.Storefront.Rating,
				IsVerified: product.Storefront.IsVerified,
			}
		}

		// Добавляем highlights, если есть
		if len(product.Highlights) > 0 {
			item.Highlights = product.Highlights
		}

		items = append(items, item)
	}

	return items, results.Total, results.TookMs, nil
}

// Helper methods

func (h *UnifiedSearchHandler) containsProductType(types []string, target string) bool {
	for _, t := range types {
		if t == target {
			return true
		}
	}
	return false
}

// mergeAndRankResults объединяет и ранжирует результаты поиска
func (h *UnifiedSearchHandler) mergeAndRankResults(items []UnifiedSearchItem, params *UnifiedSearchParams) []UnifiedSearchItem {
	// Проверяем входные параметры
	if items == nil || params == nil {
		return []UnifiedSearchItem{}
	}

	// Если нет поискового запроса, просто сортируем по указанному критерию
	if params.Query == "" {
		h.sortScoredItems(items, params.SortBy, params.SortOrder)
		return items
	}

	// Вычисляем оценку релевантности для каждого элемента
	for i := range items {
		items[i].Score = h.calculateRelevanceScore(&items[i], params.Query)
	}

	// Сортируем результаты
	h.sortScoredItems(items, params.SortBy, params.SortOrder)

	return items
}

// calculateRelevanceScore вычисляет оценку релевантности для элемента
func (h *UnifiedSearchHandler) calculateRelevanceScore(item *UnifiedSearchItem, query string) float64 {
	// Проверяем валидность входных параметров
	if item == nil || query == "" {
		return 0.0
	}

	score := item.Score // Начинаем с базового score из OpenSearch

	// Приводим к нижнему регистру для сравнения
	queryLower := strings.ToLower(query)
	titleLower := strings.ToLower(item.Name)
	descLower := strings.ToLower(item.Description)

	// Точное совпадение в заголовке (вес 5.0)
	if titleLower == queryLower {
		score += 5.0
	} else if strings.Contains(titleLower, queryLower) {
		// Частичное совпадение в заголовке (вес 3.0)
		score += 3.0
	}

	// Совпадение в описании (вес 2.0)
	if strings.Contains(descLower, queryLower) {
		score += 2.0
	}

	// Учитываем популярность (до 1.0 балла)
	if item.ViewsCount > 0 {
		popularityScore := math.Log10(float64(item.ViewsCount+1)) / 3.0 // нормализуем до ~1.0
		if popularityScore > 1.0 {
			popularityScore = 1.0
		}
		score += popularityScore
	}

	// Учитываем свежесть объявления (до 0.5 балла)
	if item.CreatedAt != nil {
		daysSinceCreated := time.Since(*item.CreatedAt).Hours() / 24
		if daysSinceCreated < 7 {
			score += 0.5 // Новые объявления получают бонус
		} else if daysSinceCreated < 30 {
			score += 0.3
		} else if daysSinceCreated < 90 {
			score += 0.1
		}
	}

	return score
}

// sortScoredItems сортирует элементы по указанному критерию
func (h *UnifiedSearchHandler) sortScoredItems(items []UnifiedSearchItem, sortBy, sortOrder string) {
	sort.Slice(items, func(i, j int) bool {
		switch sortBy {
		case "price":
			if sortOrder == "asc" {
				return items[i].Price < items[j].Price
			}
			return items[i].Price > items[j].Price

		case "date":
			if items[i].CreatedAt == nil || items[j].CreatedAt == nil {
				return false
			}
			if sortOrder == "asc" {
				return items[i].CreatedAt.Before(*items[j].CreatedAt)
			}
			return items[i].CreatedAt.After(*items[j].CreatedAt)

		case "popularity":
			if sortOrder == "asc" {
				return items[i].ViewsCount < items[j].ViewsCount
			}
			return items[i].ViewsCount > items[j].ViewsCount

		case "relevance", "":
			// По умолчанию сортируем по релевантности (score) по убыванию
			return items[i].Score > items[j].Score

		default:
			// Для неизвестных критериев сортируем по score
			return items[i].Score > items[j].Score
		}
	})
}

func (h *UnifiedSearchHandler) convertMarketplaceImages(images []models.MarketplaceImage) []UnifiedProductImage {
	result := make([]UnifiedProductImage, 0, len(images))
	for _, img := range images {
		result = append(result, UnifiedProductImage{
			URL:     img.PublicURL,
			AltText: img.FileName, // используем имя файла как alt text
			IsMain:  img.IsMain,
		})
	}
	return result
}

func (h *UnifiedSearchHandler) convertMarketplaceCategory(category *models.MarketplaceCategory) UnifiedCategoryInfo {
	if category == nil {
		return UnifiedCategoryInfo{}
	}
	return UnifiedCategoryInfo{
		ID:   category.ID,
		Name: category.Name,
		Slug: category.Slug,
	}
}

func (h *UnifiedSearchHandler) convertMarketplaceLocation(listing *models.MarketplaceListing) *UnifiedLocationInfo {
	if listing == nil {
		return nil
	}

	location := &UnifiedLocationInfo{
		City:    listing.City,
		Country: listing.Country,
	}

	if listing.Latitude != nil {
		location.Lat = *listing.Latitude
	}
	if listing.Longitude != nil {
		location.Lng = *listing.Longitude
	}

	return location
}

// trackingContext содержит все данные из fiber.Ctx для безопасной передачи в горутину
type trackingContext struct {
	userID    *int
	sessionID string
	userAgent string
	referer   string
	ipAddress string
}

// trackSearchEvent отправляет событие поиска в behavior tracking
func (h *UnifiedSearchHandler) trackSearchEvent(trackCtx *trackingContext, params *UnifiedSearchParams, result *UnifiedSearchResult, duration time.Duration) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error().
				Interface("panic", r).
				Interface("stack", string(debug.Stack())).
				Msg("Panic in trackSearchEvent")
		}
	}()

	// Проверяем валидность всех параметров
	if h == nil || h.services == nil {
		logger.Error().Msg("Handler or services is nil in trackSearchEvent")
		return
	}

	if trackCtx == nil || params == nil || result == nil {
		logger.Error().
			Bool("trackCtx_nil", trackCtx == nil).
			Bool("params_nil", params == nil).
			Bool("result_nil", result == nil).
			Msg("Invalid parameters in trackSearchEvent")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Используем данные из trackingContext
	userID := trackCtx.userID
	sessionID := trackCtx.sessionID
	if sessionID == "" {
		// Генерируем временный session_id если не предоставлен
		sessionID = "backend_" + strconv.FormatInt(time.Now().UnixNano(), 36)
	}

	// Подготавливаем фильтры для трекинга
	searchFilters := make(map[string]interface{})
	if params.ProductTypes != nil && len(params.ProductTypes) > 0 {
		searchFilters["product_types"] = params.ProductTypes
	}
	if params.CategoryID != "" {
		searchFilters["category_id"] = params.CategoryID
	}
	if params.PriceMin > 0 {
		searchFilters["price_min"] = params.PriceMin
	}
	if params.PriceMax > 0 {
		searchFilters["price_max"] = params.PriceMax
	}
	if params.City != "" {
		searchFilters["city"] = params.City
	}
	if params.StorefrontID > 0 {
		searchFilters["storefront_id"] = params.StorefrontID
	}

	// Определяем тип элемента на основе типов товаров в поиске
	var itemType behavior.ItemType

	// Проверяем, что у нас есть типы товаров
	if params.ProductTypes == nil || len(params.ProductTypes) == 0 {
		// По умолчанию используем marketplace
		itemType = behavior.ItemTypeMarketplace
	} else if len(params.ProductTypes) == 1 {
		// Если ищем только один тип, устанавливаем его
		switch params.ProductTypes[0] {
		case "marketplace":
			itemType = behavior.ItemTypeMarketplace
		case "storefront":
			itemType = behavior.ItemTypeStorefront
		default:
			// Если неизвестный тип, используем marketplace по умолчанию
			itemType = behavior.ItemTypeMarketplace
		}
	} else {
		// Если ищем несколько типов или все типы, выбираем тип с большим количеством результатов
		marketplaceCount := 0
		storefrontCount := 0

		// Проверяем, что Items не nil перед итерацией
		if result.Items != nil {
			for _, item := range result.Items {
				switch item.ProductType {
				case "marketplace":
					marketplaceCount++
				case "storefront":
					storefrontCount++
				}
			}
		}

		if marketplaceCount >= storefrontCount {
			itemType = behavior.ItemTypeMarketplace
		} else {
			itemType = behavior.ItemTypeStorefront
		}
	}

	// Создаем метаданные для трекинга
	metadata := map[string]interface{}{
		"search_query":       params.Query,
		"search_filters":     searchFilters,
		"search_sort":        params.SortBy,
		"results_count":      result.Total,
		"search_duration_ms": duration.Milliseconds(),
		"page":               params.Page,
		"limit":              params.Limit,
		"language":           params.Language,
		"device_type":        getDeviceTypeFromUserAgent(trackCtx.userAgent),
		"browser":            trackCtx.userAgent,
		"referer":            trackCtx.referer,
		"user_agent":         trackCtx.userAgent,
		"ip_address":         trackCtx.ipAddress,
	}

	// Добавляем product_types только если они не nil
	if params.ProductTypes != nil {
		metadata["product_types"] = params.ProductTypes
	}

	// Создаем запрос для трекинга
	trackingReq := &behavior.TrackEventRequest{
		SessionID:   sessionID,
		EventType:   behavior.EventTypeSearchPerformed,
		SearchQuery: params.Query,
		ItemType:    itemType,
		Metadata:    metadata,
	}

	// Отправляем событие в behavior tracking сервис (если доступен)
	behaviorSvc := h.services.BehaviorTracking()
	logger.Debug().
		Bool("behavior_svc_nil", behaviorSvc == nil).
		Str("query", params.Query).
		Msg("Checking behavior tracking service")

	if behaviorSvc != nil {
		if err := behaviorSvc.TrackEvent(ctx, userID, trackingReq); err != nil {
			logger.Error().Err(err).
				Str("session_id", sessionID).
				Str("query", params.Query).
				Msg("Failed to track search event")
		} else {
			logger.Debug().
				Str("session_id", sessionID).
				Str("query", params.Query).
				Msg("Successfully sent event to behavior tracking")
		}
	} else {
		logger.Warn().Msg("Behavior tracking service is not available")
	}
}

// getDeviceTypeFromUserAgent определяет тип устройства по User-Agent
func getDeviceTypeFromUserAgent(userAgent string) string {
	userAgent = strings.ToLower(userAgent)

	if strings.Contains(userAgent, "mobile") ||
		strings.Contains(userAgent, "android") ||
		strings.Contains(userAgent, "iphone") {
		return "mobile"
	}

	if strings.Contains(userAgent, "tablet") ||
		strings.Contains(userAgent, "ipad") {
		return "tablet"
	}

	return "desktop"
}
