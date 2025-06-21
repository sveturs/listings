// backend/internal/proj/global/handler/unified_search.go
package handler

import (
	"context"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"backend/internal/domain/models"
	"backend/internal/domain/search"
	"backend/internal/logger"
	globalService "backend/internal/proj/global/service"
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
	ProductTypes     []string               `json:"product_types" form:"product_types"`     // ["marketplace", "storefront"]
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
	Items      []UnifiedSearchItem `json:"items"`
	Total      int                 `json:"total"`
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	TotalPages int                 `json:"total_pages"`
	HasMore    bool                `json:"has_more"`
	TookMs     int64               `json:"took_ms"`
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
	
	// Парсим параметры поиска
	var params UnifiedSearchParams
	
	// Сначала пытаемся получить из JSON body
	if c.Get("Content-Type") == "application/json" {
		if err := c.BodyParser(&params); err != nil {
			logger.Debug().Err(err).Msg("Failed to parse JSON body, trying query params")
		}
	}
	
	// Получаем параметры из query string (перезаписывают JSON если есть)
	if query := c.Query("query"); query != "" {
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
		// Простая обработка comma-separated значений
		params.ProductTypes = []string{productTypes} // TODO: реализовать парсинг CSV
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
	
	// Сортируем результаты по релевантности (score)
	if len(allItems) > 1 {
		h.sortUnifiedResults(allItems, params.SortBy, params.SortOrder)
	}
	
	// Применяем пагинацию к объединенным результатам
	offset := (params.Page - 1) * params.Limit
	end := offset + params.Limit
	
	if offset >= len(allItems) {
		allItems = []UnifiedSearchItem{}
	} else if end > len(allItems) {
		allItems = allItems[offset:]
	} else {
		allItems = allItems[offset:end]
	}
	
	// Вычисляем метаданные
	totalPages := int(math.Ceil(float64(totalCount) / float64(params.Limit)))
	hasMore := params.Page < totalPages
	
	return &UnifiedSearchResult{
		Items:      allItems,
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
		Query:            params.Query,
		Page:             params.Page,
		Size:             params.Limit * 2, // Берем больше для лучшего смешивания результатов
		CategoryID:       params.CategoryID,
		PriceMin:         params.PriceMin,
		PriceMax:         params.PriceMax,
		Sort:             params.SortBy,
		SortDirection:    params.SortOrder,
		// AttributeFilters: params.AttributeFilters, // TODO: конвертировать типы
		Language:         params.Language,
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
		}
		
		items = append(items, item)
	}
	
	return items, results.Total, results.Took, nil
}

// searchStorefront поиск в storefront
func (h *UnifiedSearchHandler) searchStorefront(ctx context.Context, params *UnifiedSearchParams) ([]UnifiedSearchItem, int, int64, error) {
	// TODO: Реализовать поиск в storefront через OpenSearch
	// Пока возвращаем пустой результат
	logger.Debug().Msg("Storefront search not yet implemented")
	return []UnifiedSearchItem{}, 0, 0, nil
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

func (h *UnifiedSearchHandler) sortUnifiedResults(items []UnifiedSearchItem, sortBy, sortOrder string) {
	// TODO: Реализовать сортировку результатов
	logger.Debug().Str("sortBy", sortBy).Str("sortOrder", sortOrder).Msg("Sorting not yet implemented")
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