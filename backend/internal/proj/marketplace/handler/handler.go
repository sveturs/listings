// Package handler
// backend/internal/proj/marketplace/handler/handler.go
package handler

import (
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"backend/internal/domain/models"
	"backend/internal/logger"
	"backend/internal/middleware"
	globalService "backend/internal/proj/global/service"
	marketplaceServices "backend/internal/proj/marketplace/services"
	"backend/internal/proj/marketplace/storage/opensearch"
	"backend/internal/storage/postgres"
	"backend/pkg/utils"
)

// Global variables for caching categories
var (
	categoryTreeCache      []models.CategoryTreeNode
	categoryTreeLastUpdate time.Time
	categoryTreeMutex      sync.RWMutex
)

// SuggestionItem представляет элемент автодополнения
type SuggestionItem struct {
	Type       string                 `json:"type"`
	Value      string                 `json:"value"`
	Label      string                 `json:"label"`
	Count      int                    `json:"count,omitempty"`
	CategoryID int                    `json:"category_id,omitempty"`
	ProductID  int                    `json:"product_id,omitempty"`
	Icon       string                 `json:"icon,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// Handler combines all marketplace handlers
type Handler struct {
	Listings               *ListingsHandler
	Images                 *ImagesHandler
	Categories             *CategoriesHandler
	Search                 *SearchHandler
	Translations           *TranslationsHandler
	Favorites              *FavoritesHandler
	Indexing               *IndexingHandler
	Chat                   *ChatHandler
	AdminCategories        *AdminCategoriesHandler
	AdminAttributes        *AdminAttributesHandler
	AdminVariantAttributes *AdminVariantAttributesHandler
	AdminTranslations      *AdminTranslationsHandler
	CustomComponents       *CustomComponentHandler
	MarketplaceHandler     *MarketplaceHandler
	Orders                 *OrderHandler
	CategoryDetector       *CategoryDetectorHandler
	VariantAttributes      *VariantAttributesHandler
	Cars                   *CarsHandler
	service                globalService.ServicesInterface
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

		// Создаем OrderService если есть Orders сервис
		var orderHandler *OrderHandler
		if orderService := services.Orders(); orderService != nil {
			orderHandler = NewOrderHandler(orderService)
		}

		// Создаем репозиторий для keywords
		keywordRepo := postgres.NewCategoryKeywordRepository(postgresDB.GetSQLXDB())

		adminCategoriesHandler := NewAdminCategoriesHandler(categoriesHandler, keywordRepo)
		logger.Info().Interface("adminCategoriesHandler", adminCategoriesHandler).Msg("Created AdminCategoriesHandler")

		// Создаём CategoryDetector и CategoryDetectorHandler
		var categoryDetectorHandler *CategoryDetectorHandler
		if storage := services.Storage(); storage != nil {
			logger.Info().Msg("Storage is available, checking for OpenSearch...")
			// Пытаемся получить OpenSearch репозиторий
			if db, ok := storage.(*postgres.Database); ok {
				logger.Info().Msg("Storage is postgres.Database")
				if osRepo := db.GetOpenSearchRepository(); osRepo != nil {
					logger.Info().Msg("OpenSearch repository exists")
					// Проверяем, что это именно *opensearch.Repository
					if concreteRepo, ok := osRepo.(*opensearch.Repository); ok {
						logger.Info().Msg("OpenSearch repository is correct type")
						// Создаём сервис определения категорий
						detector, err := marketplaceServices.NewCategoryDetectorFromStorage(db, concreteRepo)
						if err != nil {
							logger.Error().Err(err).Msg("Failed to create CategoryDetector")
						} else {
							// Создаём handler
							categoryDetectorHandler = NewCategoryDetectorHandler(detector, zap.L())
							logger.Info().Msg("Created CategoryDetectorHandler successfully")
						}
					} else {
						logger.Error().Msgf("OpenSearch repository is not of expected type *opensearch.Repository, got %T", osRepo)
					}
				} else {
					logger.Error().Msg("OpenSearch repository is nil")
				}
			} else {
				logger.Error().Msg("Storage is not postgres.Database")
			}
		} else {
			logger.Error().Msg("Storage is nil")
		}

		return &Handler{
			Listings:               NewListingsHandler(services),
			Images:                 NewImagesHandler(services),
			Categories:             categoriesHandler,
			Search:                 NewSearchHandler(services),
			Translations:           NewTranslationsHandler(services),
			Favorites:              NewFavoritesHandler(services),
			Indexing:               NewIndexingHandler(services),
			Chat:                   NewChatHandler(services, services.Config()),
			AdminCategories:        adminCategoriesHandler,
			AdminAttributes:        NewAdminAttributesHandler(services),
			AdminVariantAttributes: NewAdminVariantAttributesHandler(services),
			AdminTranslations:      NewAdminTranslationsHandler(services),
			CustomComponents:       customComponentHandler,
			MarketplaceHandler:     marketplaceHandler,
			Orders:                 orderHandler,
			CategoryDetector:       categoryDetectorHandler,
			VariantAttributes:      NewVariantAttributesHandler(services),
			Cars:                   NewCarsHandler(services.Marketplace(), services.UnifiedCar()),
			service:                services,
		}
	}

	// Возвращаем handler без CustomComponents, если приведение не удалось
	// В fallback случае создаем nil keywordRepo - это временное решение
	adminCategoriesHandler := NewAdminCategoriesHandler(categoriesHandler, nil)
	logger.Info().Interface("adminCategoriesHandler", adminCategoriesHandler).Msg("Created AdminCategoriesHandler (fallback)")

	return &Handler{
		Listings:               NewListingsHandler(services),
		Images:                 NewImagesHandler(services),
		Categories:             categoriesHandler,
		Search:                 NewSearchHandler(services),
		Translations:           NewTranslationsHandler(services),
		Favorites:              NewFavoritesHandler(services),
		Indexing:               NewIndexingHandler(services),
		Chat:                   NewChatHandler(services, services.Config()),
		AdminCategories:        adminCategoriesHandler,
		AdminAttributes:        NewAdminAttributesHandler(services),
		AdminVariantAttributes: NewAdminVariantAttributesHandler(services),
		AdminTranslations:      NewAdminTranslationsHandler(services),
		CustomComponents:       nil,
		MarketplaceHandler:     nil,
		Orders:                 nil,
		CategoryDetector:       nil,
		Cars:                   NewCarsHandler(services.Marketplace(), services.UnifiedCar()),
		service:                services,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
	marketplace := app.Group("/api/v1/marketplace")
	marketplace.Get("/listings", h.Listings.GetListings)
	marketplace.Get("/categories", h.Categories.GetCategories)
	marketplace.Get("/popular-categories", h.Categories.GetPopularCategories)
	marketplace.Get("/category-tree", h.Categories.GetCategoryTree)
	marketplace.Get("/listings/slug/:slug", h.Listings.GetListingBySlug)
	marketplace.Get("/listings/:id", h.Listings.GetListing)
	marketplace.Get("/search", h.Search.SearchListingsAdvanced)  // маршрут поиска GET
	marketplace.Post("/search", h.Search.SearchListingsAdvanced) // маршрут поиска POST для расширенных фильтров
	marketplace.Get("/suggestions", h.Search.GetSuggestions)     // маршрут автодополнения
	marketplace.Get("/category-suggestions", h.Search.GetCategorySuggestions)
	marketplace.Get("/categories/:id/attributes", h.Categories.GetCategoryAttributes)
	marketplace.Get("/listings/:id/price-history", h.Listings.GetPriceHistory)
	marketplace.Get("/listings/:id/similar", h.Search.GetSimilarListings)
	marketplace.Get("/categories/:id/attribute-ranges", h.Categories.GetAttributeRanges)
	marketplace.Get("/enhanced-suggestions", h.GetEnhancedSuggestions)

	// Cars routes (public endpoints)
	if h.Cars != nil {
		cars := app.Group("/api/v1/cars") // Отдельная группа для автомобилей
		cars.Get("/makes", h.Cars.GetCarMakes)
		cars.Get("/makes/search", h.Cars.SearchCarMakes)
		cars.Get("/makes/:make_slug/models", h.Cars.GetCarModels)
		cars.Get("/models/:model_id/generations", h.Cars.GetCarGenerations)
		cars.Get("/vin/:vin/decode", h.Cars.DecodeVIN)

		logger.Info().Msg("Registered cars routes")
	}

	// Fuzzy search routes
	marketplace.Get("/test-fuzzy-search", h.Search.TestFuzzySearch)
	marketplace.Get("/fuzzy-search", h.Search.SearchWithFuzzyParams)

	// Category detection routes
	if h.CategoryDetector != nil {
		logger.Info().Msg("Registering category detection routes")
		// Добавляем тестовый эндпоинт
		marketplace.Get("/categories/detect/test", func(c *fiber.Ctx) error {
			logger.Info().Msg("Test endpoint called")
			return c.JSON(fiber.Map{"status": "ok", "message": "CategoryDetector is available"})
		})
		// Создаем wrapper функцию для вызова метода
		detectCategoryFunc := func(c *fiber.Ctx) error {
			logger.Info().Msg("=== DetectCategory route called ===")
			if h.CategoryDetector == nil {
				logger.Error().Msg("CategoryDetector is nil in route")
				return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.marketplace.categoryDetectionFailed")
			}
			logger.Info().Msg("Calling CategoryDetector.DetectCategory method...")
			return h.CategoryDetector.DetectCategory(c)
		}
		marketplace.Post("/categories/detect", detectCategoryFunc)
		marketplace.Put("/categories/detect/:stats_id/confirm", h.CategoryDetector.UpdateCategoryConfirmation)
		marketplace.Get("/categories/:category_id/keywords", h.CategoryDetector.GetCategoryKeywords)
	} else {
		logger.Error().Msg("CategoryDetector is nil, routes not registered")
	}

	// Карта - геопространственные маршруты
	marketplace.Get("/map/bounds", h.GetListingsInBounds)
	marketplace.Get("/map/clusters", h.GetMapClusters)

	// Neighborhood statistics
	marketplace.Get("/neighborhood-stats", h.MarketplaceHandler.GetNeighborhoodStats)

	// Автомобильные марки и модели
	if h.Cars != nil {
		h.Cars.RegisterRoutes(marketplace)
	}

	// Вариативные атрибуты
	marketplace.Get("/product-variant-attributes", h.VariantAttributes.GetProductVariantAttributes)
	marketplace.Get("/categories/:slug/variant-attributes", h.VariantAttributes.GetCategoryVariantAttributes)

	// Обновлено: маршруты API переводов используют обработчик переводов
	translation := app.Group("/api/v1/translation")
	translation.Get("/limits", h.Translations.GetTranslationLimits)
	translation.Post("/provider", h.Translations.SetTranslationProvider)

	authedAPIGroup := app.Group("/api/v1", mw.AuthRequiredJWT)

	marketplaceProtected := authedAPIGroup.Group("/marketplace")
	marketplaceProtected.Post("/listings", h.Listings.CreateListing)
	marketplaceProtected.Put("/listings/:id", h.Listings.UpdateListing)
	marketplaceProtected.Delete("/listings/:id", h.Listings.DeleteListing)
	marketplaceProtected.Post("/listings/check-slug", h.Listings.CheckSlugAvailability)
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

	// Регистрируем маршруты для заказов маркетплейса под marketplace префиксом
	if h.Orders != nil {
		ordersGroup := marketplaceProtected.Group("/orders")
		h.Orders.RegisterRoutes(ordersGroup)
	}

	adminRoutes := app.Group("/api/v1/admin", mw.AuthRequiredJWT, mw.AdminRequired)

	// Регистрируем маршруты администрирования категорий
	logger.Info().Msg("Registering admin categories routes")
	logger.Info().Interface("AdminCategories", h.AdminCategories).Msg("AdminCategories handler")
	if h.AdminCategories == nil {
		logger.Error().Msg("🚨🚨🚨 AdminCategories is NIL! 🚨🚨🚨")
	} else {
		logger.Info().Msg("✅ AdminCategories is NOT nil")
	}
	logger.Info().Str("route", "POST /categories").Msg("Registering CreateCategory route")

	adminRoutes.Post("/categories", h.AdminCategories.CreateCategory)
	adminRoutes.Get("/categories", h.AdminCategories.GetCategories)
	adminRoutes.Get("/categories/all", h.AdminCategories.GetAllCategories)
	adminRoutes.Get("/categories/:id", h.AdminCategories.GetCategoryByID)
	adminRoutes.Put("/categories/:id", h.AdminCategories.UpdateCategory)
	adminRoutes.Delete("/categories/:id", h.AdminCategories.DeleteCategory)
	adminRoutes.Post("/categories/:id/reorder", h.AdminCategories.ReorderCategories)
	adminRoutes.Put("/categories/:id/move", h.AdminCategories.MoveCategory)
	adminRoutes.Post("/categories/:id/attributes", h.AdminCategories.AddAttributeToCategory)
	adminRoutes.Delete("/categories/:id/attributes/:attr_id", h.AdminCategories.RemoveAttributeFromCategory)
	adminRoutes.Put("/categories/:id/attributes/:attr_id", h.AdminCategories.UpdateAttributeCategory)
	adminRoutes.Get("/categories/:id/groups", h.AdminCategories.GetCategoryAttributeGroups)
	adminRoutes.Post("/categories/:id/groups", h.AdminCategories.AttachAttributeGroupToCategory)
	adminRoutes.Delete("/categories/:id/groups/:group_id", h.AdminCategories.DetachAttributeGroupFromCategory)
	adminRoutes.Post("/categories/:id/translate", h.AdminCategories.TranslateCategory)

	// Маршруты для управления ключевыми словами категорий
	adminRoutes.Get("/categories/:category_id/keywords", h.AdminCategories.GetCategoryKeywords)
	adminRoutes.Post("/categories/:category_id/keywords", h.AdminCategories.AddCategoryKeyword)
	adminRoutes.Put("/categories/keywords/:keyword_id", h.AdminCategories.UpdateCategoryKeyword)
	adminRoutes.Delete("/categories/keywords/:keyword_id", h.AdminCategories.DeleteCategoryKeyword)

	// Маршруты для управления вариативными атрибутами категорий
	adminRoutes.Get("/categories/:id/variant-attributes", h.AdminCategories.GetCategoryVariantAttributes)
	adminRoutes.Put("/categories/:id/variant-attributes", h.AdminCategories.UpdateCategoryVariantAttributes)

	// Регистрируем маршруты администрирования атрибутов
	adminRoutes.Post("/attributes", h.AdminAttributes.CreateAttribute)
	adminRoutes.Get("/attributes", h.AdminAttributes.GetAttributes)
	adminRoutes.Get("/attributes/:id", h.AdminAttributes.GetAttributeByID)
	adminRoutes.Put("/attributes/:id", h.AdminAttributes.UpdateAttribute)
	adminRoutes.Delete("/attributes/:id", h.AdminAttributes.DeleteAttribute)
	adminRoutes.Post("/attributes/:id/translate", h.AdminAttributes.TranslateAttribute)
	adminRoutes.Post("/attributes/bulk-update", h.AdminAttributes.BulkUpdateAttributes)

	// Маршруты для экспорта/импорта настроек атрибутов
	adminRoutes.Get("/categories/:categoryId/attributes/export", h.AdminAttributes.ExportCategoryAttributes)
	adminRoutes.Post("/categories/:categoryId/attributes/import", h.AdminAttributes.ImportCategoryAttributes)
	adminRoutes.Post("/categories/:targetCategoryId/attributes/copy", h.AdminAttributes.CopyAttributesSettings)

	// Регистрируем маршруты администрирования вариативных атрибутов
	adminRoutes.Get("/variant-attributes", h.AdminVariantAttributes.GetVariantAttributes)
	adminRoutes.Post("/variant-attributes", h.AdminVariantAttributes.CreateVariantAttribute)
	adminRoutes.Get("/variant-attributes/:id", h.AdminVariantAttributes.GetVariantAttributeByID)
	adminRoutes.Put("/variant-attributes/:id", h.AdminVariantAttributes.UpdateVariantAttribute)
	adminRoutes.Delete("/variant-attributes/:id", h.AdminVariantAttributes.DeleteVariantAttribute)
	// Маршруты для управления связями вариативных атрибутов
	adminRoutes.Get("/variant-attributes/:id/mappings", h.AdminVariantAttributes.GetVariantAttributeMappings)
	adminRoutes.Put("/variant-attributes/:id/mappings", h.AdminVariantAttributes.UpdateVariantAttributeMappings)

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

	// Маршруты для админских переводов marketplace
	// Изменен путь для избежания конфликта с translation_admin модулем
	adminRoutes.Post("/marketplace-translations/batch-categories", h.AdminTranslations.BatchTranslateCategories)
	adminRoutes.Post("/marketplace-translations/batch-attributes", h.AdminTranslations.BatchTranslateAttributes)
	adminRoutes.Get("/marketplace-translations/status", h.AdminTranslations.GetTranslationStatus)
	adminRoutes.Put("/marketplace-translations/:entity_type/:entity_id/:field_name", h.AdminTranslations.UpdateFieldTranslation)

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
	chat.Get("/attachments/:id/download", h.Chat.GetAttachmentFile) // Новый защищенный роут для скачивания файлов
	chat.Delete("/attachments/:id", h.Chat.DeleteAttachment)
	chat.Get("/unread-count", h.Chat.GetUnreadCount)

	return nil
}

// GetEnhancedSuggestions returns enhanced search suggestions
// @Summary Get enhanced search suggestions
// @Description Returns enhanced autocomplete suggestions including queries, categories, and products
// @Tags marketplace-search
// @Accept json
// @Produce json
// @Param query query string true "Search query"
// @Param limit query int false "Number of suggestions" default(10)
// @Param types query string false "Comma-separated types (queries,categories,products)" default(queries,categories,products)
// @Success 200 {object} utils.SuccessResponseSwag{data=[]SuggestionItem} "Enhanced suggestions list"
// @Failure 400 {object} utils.ErrorResponseSwag "marketplace.queryRequired"
// @Failure 500 {object} utils.ErrorResponseSwag "marketplace.suggestionsError"
// @Router /api/v1/marketplace/enhanced-suggestions [get]
func (h *Handler) GetEnhancedSuggestions(c *fiber.Ctx) error {
	// Получаем параметры запроса
	query := c.Query("query")
	if query == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "marketplace.queryRequired")
	}

	// Получаем лимит
	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Получаем типы подсказок
	types := c.Query("types", "queries,categories,products")

	// Вызываем сервисный метод
	suggestions, err := h.service.Marketplace().GetEnhancedSuggestions(c.Context(), query, limit, types)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.suggestionsError")
	}

	// Возвращаем результат
	return utils.SuccessResponse(c, suggestions)
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
// @Success 200 {object} utils.SuccessResponseSwag{data=MapBoundsData} "Listings within bounds"
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

	// Получаем атрибуты из запроса
	attributes := c.Query("attributes", "")

	// Получаем объявления в указанных границах
	listings, err := h.service.Marketplace().GetListingsInBounds(c.Context(),
		neLat64, neLng64, swLat64, swLng64, zoom,
		categoryIDs, condition, minPriceFloat, maxPriceFloat, attributes)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.mapError")
	}

	response := MapBoundsData{
		Listings: listings,
		Bounds: MapBounds{
			NE: Coordinates{Lat: neLat64, Lng: neLng64},
			SW: Coordinates{Lat: swLat64, Lng: swLng64},
		},
		Zoom:  zoom,
		Count: len(listings),
	}
	return utils.SuccessResponse(c, response)
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
// @Success 200 {object} utils.SuccessResponseSwag{data=MapBoundsData} "Map clusters data"
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

	// Получаем атрибуты из запроса
	attributes := c.Query("attributes", "")

	clusters, err := h.service.Marketplace().GetListingsInBounds(c.Context(),
		neLat64, neLng64, swLat64, swLng64, zoom,
		categoryIDs, condition, minPriceFloat, maxPriceFloat, attributes)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "marketplace.mapError")
	}

	// Этот метод возвращает listings, а не clusters, поэтому используем MapBoundsData
	response := MapBoundsData{
		Listings: clusters, // clusters здесь на самом деле listings
		Bounds: MapBounds{
			NE: Coordinates{Lat: neLat64, Lng: neLng64},
			SW: Coordinates{Lat: swLat64, Lng: swLng64},
		},
		Zoom:  zoom,
		Count: len(clusters),
	}
	return utils.SuccessResponse(c, response)
}
