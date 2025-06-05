// backend/internal/proj/marketplace/handler/handler.go
package handler

import (
	"backend/internal/domain/models"
	globalService "backend/internal/proj/global/service"
	postgres "backend/internal/storage/postgres"
	"backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"sync"
	"time"
)

// Глобальные переменные для кеширования категорий
var (
	categoryTreeCache      []models.CategoryTreeNode
	categoryTreeLastUpdate time.Time
	categoryTreeMutex      sync.RWMutex
)

// Handler объединяет все обработчики маркетплейса
type Handler struct {
	Listings           *ListingsHandler
	Images             *ImagesHandler
	Categories         *CategoriesHandler
	Search             *SearchHandler
	Translations       *TranslationsHandler
	Favorites          *FavoritesHandler
	Indexing           *IndexingHandler
	Chat               *ChatHandler
	AdminCategories    *AdminCategoriesHandler
	AdminAttributes    *AdminAttributesHandler
	CustomComponents   *CustomComponentHandler
	MarketplaceHandler *MarketplaceHandler
	service            globalService.ServicesInterface
}

// NewHandler создает новый обработчик маркетплейса
func NewHandler(services globalService.ServicesInterface) *Handler {
	// Сначала создаем базовые обработчики
	categoriesHandler := NewCategoriesHandler(services)
	// Получаем storage из services и создаем хранилище для кастомных компонентов
	marketplaceService := services.Marketplace()

	// Приводим storage к postgres.Database для доступа к pool
	if postgresDB, ok := marketplaceService.Storage().(*postgres.Database); ok {
		// Создаем Storage с AttributeGroups
		storage := postgres.NewStorage(postgresDB.GetPool(), services.Translation())

		// Создаем MarketplaceHandler
		marketplaceHandler := NewMarketplaceHandler(storage)

		customComponentStorage := postgres.NewCustomComponentStorage(postgresDB)
		customComponentHandler := NewCustomComponentHandler(customComponentStorage)
		return &Handler{
			Listings:           NewListingsHandler(services),
			Images:             NewImagesHandler(services),
			Categories:         categoriesHandler,
			Search:             NewSearchHandler(services),
			Translations:       NewTranslationsHandler(services),
			Favorites:          NewFavoritesHandler(services),
			Indexing:           NewIndexingHandler(services),
			Chat:               NewChatHandler(services, services.Config()),
			AdminCategories:    NewAdminCategoriesHandler(categoriesHandler),
			AdminAttributes:    NewAdminAttributesHandler(services),
			CustomComponents:   customComponentHandler,
			MarketplaceHandler: marketplaceHandler,
			service:            services,
		}
	}

	// Возвращаем handler без CustomComponents, если приведение не удалось
	return &Handler{
		Listings:           NewListingsHandler(services),
		Images:             NewImagesHandler(services),
		Categories:         categoriesHandler,
		Search:             NewSearchHandler(services),
		Translations:       NewTranslationsHandler(services),
		Favorites:          NewFavoritesHandler(services),
		Indexing:           NewIndexingHandler(services),
		Chat:               NewChatHandler(services, services.Config()),
		AdminCategories:    NewAdminCategoriesHandler(categoriesHandler),
		AdminAttributes:    NewAdminAttributesHandler(services),
		CustomComponents:   nil,
		MarketplaceHandler: nil,
		service:            services,
	}
}

// GetListingsInBounds возвращает объявления в указанных границах карты
func (h *Handler) GetListingsInBounds(c *fiber.Ctx) error {
	// Получаем параметры bounds
	neLat := c.Query("ne_lat")
	neLng := c.Query("ne_lng")
	swLat := c.Query("sw_lat")
	swLng := c.Query("sw_lng")
	zoomStr := c.Query("zoom", "10")

	// Валидируем параметры
	if neLat == "" || neLng == "" || swLat == "" || swLng == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Missing bounding box parameters")
	}

	// Парсим координаты
	neLat64, err := strconv.ParseFloat(neLat, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ne_lat")
	}

	neLng64, err := strconv.ParseFloat(neLng, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ne_lng")
	}

	swLat64, err := strconv.ParseFloat(swLat, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid sw_lat")
	}

	swLng64, err := strconv.ParseFloat(swLng, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid sw_lng")
	}

	zoom, err := strconv.Atoi(zoomStr)
	if err != nil {
		zoom = 10
	}

	// Получаем фильтры
	categoryIDs := c.Query("categories", "")
	condition := c.Query("condition", "")
	minPrice := c.Query("min_price", "")
	maxPrice := c.Query("max_price", "")

	// Парсим цены
	var minPriceFloat, maxPriceFloat *float64
	if minPrice != "" {
		if parsed, err := strconv.ParseFloat(minPrice, 64); err == nil {
			minPriceFloat = &parsed
		}
	}
	if maxPrice != "" {
		if parsed, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			maxPriceFloat = &parsed
		}
	}

	// Получаем объявления в указанных границах
	listings, err := h.service.Marketplace().GetListingsInBounds(c.Context(),
		neLat64, neLng64, swLat64, swLng64, zoom,
		categoryIDs, condition, minPriceFloat, maxPriceFloat)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get listings")
	}

	return utils.SuccessResponse(c, map[string]interface{}{
		"listings": listings,
		"bounds": map[string]interface{}{
			"ne": map[string]float64{"lat": neLat64, "lng": neLng64},
			"sw": map[string]float64{"lat": swLat64, "lng": swLng64},
		},
		"zoom":  zoom,
		"count": len(listings),
	})
}

// GetMapClusters возвращает кластеризованные данные для карты
func (h *Handler) GetMapClusters(c *fiber.Ctx) error {
	// Получаем параметры bounds
	neLat := c.Query("ne_lat")
	neLng := c.Query("ne_lng")
	swLat := c.Query("sw_lat")
	swLng := c.Query("sw_lng")
	zoomStr := c.Query("zoom", "10")

	// Валидируем параметры
	if neLat == "" || neLng == "" || swLat == "" || swLng == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Missing bounding box parameters")
	}

	// Парсим координаты
	neLat64, err := strconv.ParseFloat(neLat, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ne_lat")
	}

	neLng64, err := strconv.ParseFloat(neLng, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ne_lng")
	}

	swLat64, err := strconv.ParseFloat(swLat, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid sw_lat")
	}

	swLng64, err := strconv.ParseFloat(swLng, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid sw_lng")
	}

	zoom, err := strconv.Atoi(zoomStr)
	if err != nil {
		zoom = 10
	}

	// Получаем фильтры
	categoryIDs := c.Query("categories", "")
	condition := c.Query("condition", "")
	minPrice := c.Query("min_price", "")
	maxPrice := c.Query("max_price", "")

	// Парсим цены
	var minPriceFloat, maxPriceFloat *float64
	if minPrice != "" {
		if parsed, err := strconv.ParseFloat(minPrice, 64); err == nil {
			minPriceFloat = &parsed
		}
	}
	if maxPrice != "" {
		if parsed, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			maxPriceFloat = &parsed
		}
	}

	// Получаем кластеры в указанных границах
	_ = neLat64       // используем переменную
	_ = neLng64       // используем переменную
	_ = swLat64       // используем переменную
	_ = swLng64       // используем переменную
	_ = zoom          // используем переменную
	_ = categoryIDs   // используем переменную
	_ = condition     // используем переменную
	_ = minPriceFloat // используем переменную

	clusters, err := h.service.Marketplace().GetListingsInBounds(c.Context(),
		neLat64, neLng64, swLat64, swLng64, zoom,
		categoryIDs, condition, minPriceFloat, maxPriceFloat)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get clusters")
	}

	return utils.SuccessResponse(c, map[string]interface{}{
		"clusters": clusters,
		"bounds": map[string]interface{}{
			"ne": map[string]float64{"lat": neLat64, "lng": neLng64},
			"sw": map[string]float64{"lat": swLat64, "lng": swLng64},
		},
		"zoom":  zoom,
		"count": len(clusters),
	})
}

