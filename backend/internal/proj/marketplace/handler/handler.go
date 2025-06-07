// backend/internal/proj/marketplace/handler/handler.go
package handler

import (
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"

	"backend/internal/domain/models"
	"backend/internal/middleware"
	globalService "backend/internal/proj/global/service"
	postgres "backend/internal/storage/postgres"
	"backend/pkg/utils"
)

// Global variables for caching categories
var (
	categoryTreeCache      []models.CategoryTreeNode
	categoryTreeLastUpdate time.Time
	categoryTreeMutex      sync.RWMutex
)

// Handler combines all marketplace handlers
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

func (h *Handler) GetPrefix() string {
	return "/api/v1/marketplace"
}

// NewHandler creates a new marketplace handler
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

func (h *Handler) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {

	marketplace := app.Group("/api/v1/marketplace")
	marketplace.Get("/listings", h.Listings.GetListings)
	marketplace.Get("/categories", h.Categories.GetCategories)
	marketplace.Get("/category-tree", h.Categories.GetCategoryTree)
	marketplace.Get("/listings/:id", h.Listings.GetListing)
	marketplace.Get("/search", h.Search.SearchListingsAdvanced) // маршрут поиска
	marketplace.Get("/suggestions", h.Search.GetSuggestions)    // маршрут автодополнения
	marketplace.Get("/category-suggestions", h.Search.GetCategorySuggestions)

	// Временный публичный маршрут для проверки
	app.Get("/admin-categories-test", h.AdminCategories.GetCategories)
	marketplace.Get("/categories/:id/attributes", h.Categories.GetCategoryAttributes)
	marketplace.Get("/listings/:id/price-history", h.Listings.GetPriceHistory)
	marketplace.Get("/listings/:id/similar", h.Search.GetSimilarListings)
	marketplace.Get("/categories/:id/attribute-ranges", h.Categories.GetAttributeRanges)
	marketplace.Get("/enhanced-suggestions", h.Search.GetEnhancedSuggestions)

	// Карта - геопространственные маршруты
	marketplace.Get("/map/bounds", h.GetListingsInBounds)
	marketplace.Get("/map/clusters", h.GetMapClusters)

	// Обновлено: маршруты API переводов используют обработчик переводов
	translation := app.Group("/api/v1/translation")
	translation.Get("/limits", h.Translations.GetTranslationLimits)
	translation.Post("/provider", h.Translations.SetTranslationProvider)

	authedAPIGroup := app.Group("/api/v1", mw.AuthRequiredJWT, mw.CSRFProtection())

	marketplaceProtected := authedAPIGroup.Group("/marketplace")
	marketplaceProtected.Post("/listings", h.Listings.CreateListing)
	marketplaceProtected.Put("/listings/:id", h.Listings.UpdateListing)
	marketplaceProtected.Delete("/listings/:id", h.Listings.DeleteListing)
	marketplaceProtected.Post("/listings/:id/images", h.Images.UploadImages)
	marketplaceProtected.Post("/listings/:id/favorite", h.Favorites.AddToFavorites)
	marketplaceProtected.Delete("/listings/:id/favorite", h.Favorites.RemoveFromFavorites)
	marketplaceProtected.Get("/favorites", h.Favorites.GetFavorites)
	marketplaceProtected.Put("/translations/:id", h.Translations.UpdateTranslations)
	marketplaceProtected.Post("/translations/batch", h.Translations.TranslateText) // Предполагается, что этот метод переименован
	marketplaceProtected.Post("/moderate-image", h.Images.ModerateImage)
	marketplaceProtected.Post("/enhance-preview", h.Images.EnhancePreview)
	marketplaceProtected.Post("/enhance-images", h.Images.EnhanceImages)

	// маршруты для новых методов в TranslationsHandler
	marketplaceProtected.Post("/translations/batch-translate", h.Translations.BatchTranslateListings)
	marketplaceProtected.Post("/translations/translate", h.Translations.TranslateText)
	marketplaceProtected.Post("/translations/detect-language", h.Translations.DetectLanguage)
	marketplaceProtected.Get("/translations/:id", h.Translations.GetTranslations)

	adminRoutes := app.Group("/api/v1/admin", mw.AuthRequiredJWT, mw.AdminRequired, mw.CSRFProtection())

	// Регистрируем маршруты администрирования категорий
	adminRoutes.Post("/categories", h.AdminCategories.CreateCategory)
	adminRoutes.Get("/categories", h.AdminCategories.GetCategories)
	adminRoutes.Get("/categories/:id", h.AdminCategories.GetCategoryByID)
	adminRoutes.Put("/categories/:id", h.AdminCategories.UpdateCategory)
	adminRoutes.Delete("/categories/:id", h.AdminCategories.DeleteCategory)
	adminRoutes.Post("/categories/:id/reorder", h.AdminCategories.ReorderCategories)
	adminRoutes.Put("/categories/:id/move", h.AdminCategories.MoveCategory)
	adminRoutes.Post("/categories/:id/attributes", h.AdminCategories.AddAttributeToCategory)
	adminRoutes.Delete("/categories/:id/attributes/:attr_id", h.AdminCategories.RemoveAttributeFromCategory)
	adminRoutes.Put("/categories/:id/attributes/:attr_id", h.AdminCategories.UpdateAttributeCategory)

	// Регистрируем маршруты администрирования атрибутов
	adminRoutes.Post("/attributes", h.AdminAttributes.CreateAttribute)
	adminRoutes.Get("/attributes", h.AdminAttributes.GetAttributes)
	adminRoutes.Get("/attributes/:id", h.AdminAttributes.GetAttributeByID)
	adminRoutes.Put("/attributes/:id", h.AdminAttributes.UpdateAttribute)
	adminRoutes.Delete("/attributes/:id", h.AdminAttributes.DeleteAttribute)
	adminRoutes.Post("/attributes/bulk-update", h.AdminAttributes.BulkUpdateAttributes)

	// Маршруты для экспорта/импорта настроек атрибутов
	adminRoutes.Get("/categories/:categoryId/attributes/export", h.AdminAttributes.ExportCategoryAttributes)
	adminRoutes.Post("/categories/:categoryId/attributes/import", h.AdminAttributes.ImportCategoryAttributes)
	adminRoutes.Post("/categories/:targetCategoryId/attributes/copy", h.AdminAttributes.CopyAttributesSettings)

	// Для обратной совместимости добавим маршруты без v1
	legacyAdmin := app.Group("/api/admin", mw.AuthRequiredJWT, mw.AdminRequired, mw.CSRFProtection())

	// Все маршруты для категорий
	legacyAdmin.Get("/categories", h.AdminCategories.GetCategories)
	legacyAdmin.Post("/categories", h.AdminCategories.CreateCategory)
	legacyAdmin.Get("/categories/:id", h.AdminCategories.GetCategoryByID)
	legacyAdmin.Put("/categories/:id", h.AdminCategories.UpdateCategory)
	legacyAdmin.Delete("/categories/:id", h.AdminCategories.DeleteCategory)
	legacyAdmin.Post("/categories/:id/reorder", h.AdminCategories.ReorderCategories)
	legacyAdmin.Put("/categories/:id/move", h.AdminCategories.MoveCategory)
	legacyAdmin.Post("/categories/:id/attributes", h.AdminCategories.AddAttributeToCategory)
	legacyAdmin.Delete("/categories/:id/attributes/:attr_id", h.AdminCategories.RemoveAttributeFromCategory)
	legacyAdmin.Put("/categories/:id/attributes/:attr_id", h.AdminCategories.UpdateAttributeCategory)

	// Маршруты для атрибутов
	legacyAdmin.Post("/attributes", h.AdminAttributes.CreateAttribute)
	legacyAdmin.Get("/attributes", h.AdminAttributes.GetAttributes)
	legacyAdmin.Get("/attributes/:id", h.AdminAttributes.GetAttributeByID)
	legacyAdmin.Put("/attributes/:id", h.AdminAttributes.UpdateAttribute)
	legacyAdmin.Delete("/attributes/:id", h.AdminAttributes.DeleteAttribute)
	legacyAdmin.Post("/attributes/bulk-update", h.AdminAttributes.BulkUpdateAttributes)

	// Добавляем маршруты для экспорта/импорта настроек атрибутов
	legacyAdmin.Get("/categories/:categoryId/attributes/export", h.AdminAttributes.ExportCategoryAttributes)
	legacyAdmin.Post("/categories/:categoryId/attributes/import", h.AdminAttributes.ImportCategoryAttributes)
	legacyAdmin.Post("/categories/:targetCategoryId/attributes/copy", h.AdminAttributes.CopyAttributesSettings)

	// Маршруты для групп атрибутов
	legacyAdmin.Get("/attribute-groups", h.MarketplaceHandler.ListAttributeGroups)
	legacyAdmin.Post("/attribute-groups", h.MarketplaceHandler.CreateAttributeGroup)
	legacyAdmin.Get("/attribute-groups/:id", h.MarketplaceHandler.GetAttributeGroup)
	legacyAdmin.Put("/attribute-groups/:id", h.MarketplaceHandler.UpdateAttributeGroup)
	legacyAdmin.Delete("/attribute-groups/:id", h.MarketplaceHandler.DeleteAttributeGroup)
	legacyAdmin.Get("/attribute-groups/:id/items", h.MarketplaceHandler.GetAttributeGroupWithItems)
	legacyAdmin.Post("/attribute-groups/:id/items", h.MarketplaceHandler.AddItemToGroup)
	legacyAdmin.Delete("/attribute-groups/:id/items/:attributeId", h.MarketplaceHandler.RemoveItemFromGroup)

	// Маршруты для привязки групп к категориям
	legacyAdmin.Get("/categories/:id/attribute-groups", h.MarketplaceHandler.GetCategoryGroups)
	legacyAdmin.Post("/categories/:id/attribute-groups", h.MarketplaceHandler.AttachGroupToCategory)
	legacyAdmin.Delete("/categories/:id/attribute-groups/:groupId", h.MarketplaceHandler.DetachGroupFromCategory)

	// Маршруты для кастомных UI компонентов
	// ВАЖНО: Более специфичные роуты должны идти раньше параметризованных

	// Маршруты для шаблонов (должны быть перед :id, чтобы не конфликтовать)
	adminRoutes.Get("/custom-components/templates", h.CustomComponents.ListTemplates)
	adminRoutes.Post("/custom-components/templates", h.CustomComponents.CreateTemplate)

	// Маршруты для использования компонентов
	adminRoutes.Get("/custom-components/usage", h.CustomComponents.GetComponentUsages)
	adminRoutes.Post("/custom-components/usage", h.CustomComponents.AddComponentUsage)
	adminRoutes.Delete("/custom-components/usage/:id", h.CustomComponents.RemoveComponentUsage)

	// Основные маршруты компонентов (параметризованные идут последними)
	adminRoutes.Post("/custom-components", h.CustomComponents.CreateComponent)
	adminRoutes.Get("/custom-components", h.CustomComponents.ListComponents)
	adminRoutes.Get("/custom-components/:id", h.CustomComponents.GetComponent)
	adminRoutes.Put("/custom-components/:id", h.CustomComponents.UpdateComponent)
	adminRoutes.Delete("/custom-components/:id", h.CustomComponents.DeleteComponent)

	adminRoutes.Get("/categories/:category_id/components", h.CustomComponents.GetCategoryComponents)

	// Маршруты для групп атрибутов
	adminRoutes.Get("/attribute-groups", h.MarketplaceHandler.ListAttributeGroups)
	adminRoutes.Post("/attribute-groups", h.MarketplaceHandler.CreateAttributeGroup)
	adminRoutes.Get("/attribute-groups/:id", h.MarketplaceHandler.GetAttributeGroup)
	adminRoutes.Put("/attribute-groups/:id", h.MarketplaceHandler.UpdateAttributeGroup)
	adminRoutes.Delete("/attribute-groups/:id", h.MarketplaceHandler.DeleteAttributeGroup)
	adminRoutes.Get("/attribute-groups/:id/items", h.MarketplaceHandler.GetAttributeGroupWithItems)
	adminRoutes.Post("/attribute-groups/:id/items", h.MarketplaceHandler.AddItemToGroup)
	adminRoutes.Delete("/attribute-groups/:id/items/:attributeId", h.MarketplaceHandler.RemoveItemFromGroup)

	// Маршруты для привязки групп к категориям
	adminRoutes.Get("/categories/:id/attribute-groups", h.MarketplaceHandler.GetCategoryGroups)
	adminRoutes.Post("/categories/:id/attribute-groups", h.MarketplaceHandler.AttachGroupToCategory)
	adminRoutes.Delete("/categories/:id/attribute-groups/:groupId", h.MarketplaceHandler.DetachGroupFromCategory)

	// Использовать реальный обработчик из UserHandler

	// Управление администраторами

	// Обновлено: маршруты админских функций используют обработчик индексации
	adminRoutes.Post("/reindex-listings", h.Indexing.ReindexAll)
	adminRoutes.Post("/reindex-listings-with-translations", h.Indexing.ReindexAllWithTranslations)
	adminRoutes.Post("/sync-discounts", h.Listings.SynchronizeDiscounts) // Оставляем в Listings, т.к. это работа с объявлениями
	adminRoutes.Post("/reindex-ratings", h.Indexing.ReindexRatings)

	chat := authedAPIGroup.Group("/marketplace/chat")
	chat.Get("/", h.Chat.GetChats)
	chat.Get("/messages", h.Chat.GetMessages)

	// Применяем rate limiting для отправки сообщений и загрузки файлов
	chat.Post("/messages", mw.RateLimitMessages(), h.Chat.SendMessage)
	chat.Put("/messages/read", h.Chat.MarkAsRead)
	chat.Post("/:chat_id/archive", h.Chat.ArchiveChat)

	// Роуты для работы с вложениями с rate limiting
	chat.Post("/messages/:id/attachments", mw.RateLimitMessages(), h.Chat.UploadAttachments)
	chat.Get("/attachments/:id", h.Chat.GetAttachment)
	chat.Delete("/attachments/:id", h.Chat.DeleteAttachment)
	chat.Get("/unread-count", h.Chat.GetUnreadCount)

	return nil
}

// GetListingsInBounds returns listings within specified map bounds
// @Summary Get listings in bounds
// @Description Returns all listings within the specified geographical bounds
// @Tags marketplace-map
// @Accept json
// @Produce json
// @Param ne_lat query number true "Northeast latitude"
// @Param ne_lng query number true "Northeast longitude"
// @Param sw_lat query number true "Southwest latitude"
// @Param sw_lng query number true "Southwest longitude"
// @Param zoom query int false "Map zoom level" default(10)
// @Param categories query string false "Comma-separated category IDs"
// @Param condition query string false "Item condition filter"
// @Param min_price query number false "Minimum price filter"
// @Param max_price query number false "Maximum price filter"
// @Success 200 {object} object{listings=[]models.MarketplaceListing,bounds=object{ne=object{lat=number,lng=number},sw=object{lat=number,lng=number}},zoom=int,count=int} "Listings within bounds"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidBounds"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.mapError"
// @Router /api/v1/marketplace/map/bounds [get]
func (h *Handler) GetListingsInBounds(c *fiber.Ctx) error {
	// Получаем параметры bounds
	neLat := c.Query("ne_lat")
	neLng := c.Query("ne_lng")
	swLat := c.Query("sw_lat")
	swLng := c.Query("sw_lng")
	zoomStr := c.Query("zoom", "10")

	// Валидируем параметры
	if neLat == "" || neLng == "" || swLat == "" || swLng == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.missingBounds")
	}

	// Парсим координаты
	neLat64, err := strconv.ParseFloat(neLat, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidLatitude")
	}

	neLng64, err := strconv.ParseFloat(neLng, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidLongitude")
	}

	swLat64, err := strconv.ParseFloat(swLat, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidLatitude")
	}

	swLng64, err := strconv.ParseFloat(swLng, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidLongitude")
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
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.mapError")
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

// GetMapClusters returns clustered data for map view
// @Summary Get map clusters
// @Description Returns clustered listings data for efficient map rendering
// @Tags marketplace-map
// @Accept json
// @Produce json
// @Param ne_lat query number true "Northeast latitude"
// @Param ne_lng query number true "Northeast longitude"
// @Param sw_lat query number true "Southwest latitude"
// @Param sw_lng query number true "Southwest longitude"
// @Param zoom query int false "Map zoom level" default(10)
// @Param categories query string false "Comma-separated category IDs"
// @Param condition query string false "Item condition filter"
// @Param min_price query number false "Minimum price filter"
// @Param max_price query number false "Maximum price filter"
// @Success 200 {object} object{clusters=[]models.MarketplaceListing,bounds=object{ne=object{lat=number,lng=number},sw=object{lat=number,lng=number}},zoom=int,count=int} "Map clusters data"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.invalidBounds"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.mapError"
// @Router /api/v1/marketplace/map/clusters [get]
func (h *Handler) GetMapClusters(c *fiber.Ctx) error {
	// Получаем параметры bounds
	neLat := c.Query("ne_lat")
	neLng := c.Query("ne_lng")
	swLat := c.Query("sw_lat")
	swLng := c.Query("sw_lng")
	zoomStr := c.Query("zoom", "10")

	// Валидируем параметры
	if neLat == "" || neLng == "" || swLat == "" || swLng == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.missingBounds")
	}

	// Парсим координаты
	neLat64, err := strconv.ParseFloat(neLat, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidLatitude")
	}

	neLng64, err := strconv.ParseFloat(neLng, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidLongitude")
	}

	swLat64, err := strconv.ParseFloat(swLat, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidLatitude")
	}

	swLng64, err := strconv.ParseFloat(swLng, 64)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidLongitude")
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

	clusters, err := h.service.Marketplace().GetListingsInBounds(c.Context(),
		neLat64, neLng64, swLat64, swLng64, zoom,
		categoryIDs, condition, minPriceFloat, maxPriceFloat)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.mapError")
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
