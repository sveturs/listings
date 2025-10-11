// Package handler
// backend/internal/proj/c2c/handler/handler.go
package handler

import (
	"context"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"
	"go.uber.org/zap"

	"backend/internal/config"
	"backend/internal/domain/models"
	"backend/internal/logger"
	"backend/internal/middleware"
	"backend/internal/proj/c2c/cache"
	"backend/internal/proj/c2c/repository"
	marketplaceServices "backend/internal/proj/c2c/services"
	"backend/internal/proj/c2c/storage/opensearch"
	globalService "backend/internal/proj/global/service"
	"backend/internal/storage/postgres"
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
	Listings               *ListingsHandler
	Images                 *ImagesHandler
	Categories             *CategoriesHandler
	Search                 *SearchHandler
	Translations           *TranslationsHandler
	Favorites              *FavoritesHandler
	SavedSearches          *SavedSearchesHandler
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
	VariantMappings        *VariantMappingsHandler
	Cars                   *CarsHandler
	UnifiedAttributes      *UnifiedAttributesHandler
	AICategoryHandler      *AICategoryHandler
	service                globalService.ServicesInterface
	jwtParserMW            fiber.Handler
}

func (h *Handler) GetPrefix() string {
	return "/api/v1/c2c"
}

// NewHandler creates a new marketplace handler
func NewHandler(ctx context.Context, services globalService.ServicesInterface, jwtParserMW fiber.Handler) *Handler {
	// Сначала создаем базовые обработчики
	categoriesHandler := NewCategoriesHandler(services)
	// Получаем storage из services и создаем хранилище для кастомных компонентов
	marketplaceService := services.Marketplace()

	// Приводим storage к postgres.Database для доступа к pool
	if postgresDB, ok := marketplaceService.Storage().(*postgres.Database); ok {
		// Получаем auth UserService и устанавливаем в marketplace storage
		authUserSvc := services.AuthUserService()
		postgresDB.SetMarketplaceUserService(authUserSvc)

		// Создаем Storage с AttributeGroups
		storage := postgres.NewStorage(postgresDB.GetPool(), services.Translation())

		// Создаем MarketplaceHandler
		marketplaceHandler := NewMarketplaceHandler(storage, marketplaceService)

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

		// Создаем UnifiedAttributesHandler
		// Получаем feature flags из конфигурации
		featureFlags := config.LoadFeatureFlags()
		unifiedAttrStorage := postgres.NewUnifiedAttributeStorage(postgresDB.GetPool(), featureFlags.UnifiedAttributesFallback)
		unifiedAttributesHandler := NewUnifiedAttributesHandler(unifiedAttrStorage, featureFlags)

		// Создаем универсальный кеш для маркетплейса
		var universalCache *cache.UniversalCache
		redisAddr := "localhost:6379" // TODO: взять из конфига
		if cfg := services.Config(); cfg != nil && cfg.Redis.URL != "" {
			redisAddr = cfg.Redis.URL
		}

		universalCache, err := cache.NewUniversalCache(ctx, redisAddr, zap.L(), cache.DefaultCacheConfig())
		if err != nil {
			logger.Warn().Err(err).Msg("Failed to create universal cache, continuing without cache")
			universalCache = nil
		} else {
			logger.Info().Msg("Universal cache created successfully")
		}

		// Создаём CategoryDetector и CategoryDetectorHandler
		var categoryDetectorHandler *CategoryDetectorHandler
		var aiCategoryHandler *AICategoryHandler
		if storage := services.Storage(); storage != nil {
			logger.Info().Msg("Storage is available, checking for OpenSearch...")
			// Пытаемся получить OpenSearch репозиторий
			if db, ok := storage.(*postgres.Database); ok {
				logger.Info().Msg("Storage is postgres.Database")

				// Создаём AI Category Detector независимо от OpenSearch
				// так как он использует только PostgreSQL
				aiDetector := marketplaceServices.NewAICategoryDetector(ctx, db.GetSQLXDB(), zap.L())

				// Создаём остальные AI сервисы для полной интеграции
				redisClient := redis.NewClient(&redis.Options{
					Addr: "localhost:6379",
					DB:   0,
				})

				// Создаём все необходимые сервисы
				validator := marketplaceServices.NewAICategoryValidator(zap.L(), redisClient)
				keywordRepo := repository.NewKeywordRepository(db.GetSQLXDB(), zap.L())
				keywordGenerator := marketplaceServices.NewAIKeywordGenerator(zap.L(), redisClient, validator)

				// TODO: Создать FeedbackRepository - пока используем nil
				learningSystem := marketplaceServices.NewAILearningSystem(zap.L(), redisClient, keywordRepo, validator, keywordGenerator, nil)

				aiCategoryHandler = NewAICategoryHandler(aiDetector, validator, keywordGenerator, keywordRepo, learningSystem, zap.L())
				logger.Info().Msg("Created AICategoryHandler successfully")

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
			Listings:               NewListingsHandler(services, universalCache),
			Images:                 NewImagesHandler(services),
			Categories:             categoriesHandler,
			Search:                 NewSearchHandler(services, universalCache),
			Translations:           NewTranslationsHandler(services),
			Favorites:              NewFavoritesHandler(services),
			SavedSearches:          NewSavedSearchesHandler(services),
			Indexing:               NewIndexingHandler(services),
			Chat:                   NewChatHandler(services, services.Config(), jwtParserMW),
			AdminCategories:        adminCategoriesHandler,
			AdminAttributes:        NewAdminAttributesHandler(services),
			AdminVariantAttributes: NewAdminVariantAttributesHandler(services),
			AdminTranslations:      NewAdminTranslationsHandler(services),
			CustomComponents:       customComponentHandler,
			MarketplaceHandler:     marketplaceHandler,
			Orders:                 orderHandler,
			CategoryDetector:       categoryDetectorHandler,
			VariantAttributes:      NewVariantAttributesHandler(services),
			VariantMappings:        NewVariantMappingsHandler(services, unifiedAttrStorage, featureFlags),
			Cars:                   NewCarsHandler(services.Marketplace(), services.UnifiedCar()),
			UnifiedAttributes:      unifiedAttributesHandler,
			AICategoryHandler:      aiCategoryHandler,
			service:                services,
			jwtParserMW:            jwtParserMW,
		}
	}

	// Возвращаем handler без CustomComponents, если приведение не удалось
	// В fallback случае создаем nil keywordRepo - это временное решение
	adminCategoriesHandler := NewAdminCategoriesHandler(categoriesHandler, nil)
	logger.Info().Interface("adminCategoriesHandler", adminCategoriesHandler).Msg("Created AdminCategoriesHandler (fallback)")

	// В fallback случае все равно создаем UnifiedAttributesHandler
	// (используем nil для storage - будет работать только fallback)

	return &Handler{
		Listings:               NewListingsHandler(services, nil), // В fallback случае используем nil для кеша
		Images:                 NewImagesHandler(services),
		Categories:             categoriesHandler,
		Search:                 NewSearchHandler(services, nil), // В fallback случае используем nil для кеша
		Translations:           NewTranslationsHandler(services),
		Favorites:              NewFavoritesHandler(services),
		SavedSearches:          NewSavedSearchesHandler(services),
		Indexing:               NewIndexingHandler(services),
		Chat:                   NewChatHandler(services, services.Config(), jwtParserMW),
		AdminCategories:        adminCategoriesHandler,
		AdminAttributes:        NewAdminAttributesHandler(services),
		AdminVariantAttributes: NewAdminVariantAttributesHandler(services),
		AdminTranslations:      NewAdminTranslationsHandler(services),
		CustomComponents:       nil,
		MarketplaceHandler:     nil,
		Orders:                 nil,
		CategoryDetector:       nil,
		Cars:                   NewCarsHandler(services.Marketplace(), services.UnifiedCar()),
		UnifiedAttributes:      nil, // В fallback случае не создаем
		AICategoryHandler:      nil, // В fallback случае нет AI handler
		service:                services,
		jwtParserMW:            jwtParserMW,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
	c2c := app.Group("/api/v1/c2c")
	c2c.Get("/listings", h.Listings.GetListings)
	c2c.Get("/categories", h.Categories.GetCategories)
	c2c.Get("/popular-categories", h.Categories.GetPopularCategories)
	c2c.Get("/category-tree", h.Categories.GetCategoryTree)
	c2c.Get("/listings/slug/:slug", h.Listings.GetListingBySlug)
	c2c.Get("/listings/:id", h.Listings.GetListing)
	c2c.Get("/search", h.Search.SearchListingsAdvanced)      // маршрут поиска GET
	c2c.Post("/search", h.Search.SearchListingsAdvanced)     // маршрут поиска POST для расширенных фильтров
	c2c.Get("/suggestions", h.Search.GetSuggestions)         // маршрут автодополнения
	c2c.Get("/search/autocomplete", h.Search.GetSuggestions) // алиас для совместимости с фронтендом
	c2c.Get("/category-suggestions", h.Search.GetCategorySuggestions)
	c2c.Get("/enhanced-suggestions", h.Search.GetEnhancedSuggestions) // улучшенные предложения
	c2c.Get("/categories/:id/attributes", h.Categories.GetCategoryAttributes)
	c2c.Get("/listings/:id/price-history", h.Listings.GetPriceHistory)
	c2c.Get("/listings/:id/similar", h.Search.GetSimilarListings)
	c2c.Get("/categories/:id/attribute-ranges", h.Categories.GetAttributeRanges)

	// Public recommendations endpoint
	c2c.Get("/recommendations", h.MarketplaceHandler.GetPublicRecommendations)

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
	c2c.Get("/test-fuzzy-search", h.Search.TestFuzzySearch)
	c2c.Get("/fuzzy-search", h.Search.SearchWithFuzzyParams)

	// Category detection routes
	if h.CategoryDetector != nil {
		logger.Info().Msg("Registering category detection routes")
		// Добавляем тестовый эндпоинт
		c2c.Get("/categories/detect/test", func(c *fiber.Ctx) error {
			logger.Info().Msg("Test endpoint called")
			return c.JSON(fiber.Map{"status": "ok", "message": "CategoryDetector is available"})
		})
		// Создаем wrapper функцию для вызова метода
		detectCategoryFunc := func(c *fiber.Ctx) error {
			logger.Info().Msg("=== DetectCategory route called ===")
			if h.CategoryDetector == nil {
				logger.Error().Msg("CategoryDetector is nil in route")
				return utils.ErrorResponse(c, fiber.StatusInternalServerError, "errors.c2c.categoryDetectionFailed")
			}
			logger.Info().Msg("Calling CategoryDetector.DetectCategory method...")
			return h.CategoryDetector.DetectCategory(c)
		}
		c2c.Post("/categories/detect", detectCategoryFunc)
		c2c.Put("/categories/detect/:stats_id/confirm", h.CategoryDetector.UpdateCategoryConfirmation)
		c2c.Get("/categories/:category_id/keywords", h.CategoryDetector.GetCategoryKeywords)
	} else {
		logger.Error().Msg("CategoryDetector is nil, routes not registered")
	}

	// AI Category Detection routes (enhanced)
	if h.AICategoryHandler != nil {
		logger.Info().Msg("Registering AI category detection routes")
		aiGroup := c2c.Group("/ai")
		aiGroup.Post("/detect-category", h.AICategoryHandler.DetectCategory)
		aiGroup.Post("/select-category", h.AICategoryHandler.SelectCategory)     // НОВЫЙ МЕТОД: прямой выбор через AI
		aiGroup.Post("/validate-category", h.AICategoryHandler.ValidateCategory) // ДОБАВЛЕН НЕДОСТАЮЩИЙ РОУТ
		aiGroup.Post("/confirm/:feedbackId", h.AICategoryHandler.ConfirmDetection)
		aiGroup.Get("/metrics", h.AICategoryHandler.GetAccuracyMetrics)
		aiGroup.Post("/learn", h.jwtParserMW, authMiddleware.RequireAuth(), h.AICategoryHandler.TriggerLearning) // Защищено для админов
	}

	// Карта - геопространственные маршруты
	c2c.Get("/map/bounds", h.MarketplaceHandler.GetListingsInBounds)
	c2c.Get("/map/clusters", h.MarketplaceHandler.GetMapClusters)

	// Neighborhood statistics
	c2c.Get("/neighborhood-stats", h.MarketplaceHandler.GetNeighborhoodStats)

	// Автомобильные марки и модели
	if h.Cars != nil {
		h.Cars.RegisterRoutes(c2c)
	}

	// Вариативные атрибуты
	c2c.Get("/product-variant-attributes", h.VariantAttributes.GetProductVariantAttributes)
	c2c.Get("/categories/:slug/variant-attributes", h.VariantAttributes.GetCategoryVariantAttributes)

	// V2 API с унифицированными атрибутами (если включен feature flag)
	if h.UnifiedAttributes != nil && h.service.Config().FeatureFlags != nil && h.service.Config().FeatureFlags.UseUnifiedAttributes {
		logger.Info().Msg("Registering v2 unified attributes routes")

		// Создаем middleware для проверки feature flags
		featureFlagsMiddleware := middleware.NewFeatureFlagsMiddleware(h.service.Config().FeatureFlags)

		// V2 API группа с проверкой feature flags
		v2 := app.Group("/api/v2")
		v2Marketplace := v2.Group("/marketplace", featureFlagsMiddleware.CheckUnifiedAttributes())

		// Публичные эндпоинты v2
		v2Marketplace.Get("/categories/:category_id/attributes", h.UnifiedAttributes.GetCategoryAttributes)
		v2Marketplace.Get("/listings/:listing_id/attributes", h.UnifiedAttributes.GetListingAttributeValues)
		v2Marketplace.Get("/categories/:category_id/attribute-ranges", h.UnifiedAttributes.GetAttributeRanges)

		// Защищенные эндпоинты v2 (требуют авторизации)
		v2Protected := v2.Group("/marketplace", h.jwtParserMW, authMiddleware.RequireAuth(), featureFlagsMiddleware.CheckUnifiedAttributes())
		v2Protected.Post("/listings/:listing_id/attributes", h.UnifiedAttributes.SaveListingAttributeValues)
		v2Protected.Put("/listings/:listing_id/attributes", h.UnifiedAttributes.UpdateListingAttributeValues)

		// Административные эндпоинты v2
		v2Admin := app.Group("/api/v2/admin", h.jwtParserMW, authMiddleware.RequireAuthString("admin"), featureFlagsMiddleware.CheckUnifiedAttributes())
		v2Admin.Post("/attributes", h.UnifiedAttributes.CreateAttribute)
		v2Admin.Put("/attributes/:attribute_id", h.UnifiedAttributes.UpdateAttribute)
		v2Admin.Delete("/attributes/:attribute_id", h.UnifiedAttributes.DeleteAttribute)
		v2Admin.Post("/categories/:category_id/attributes", h.UnifiedAttributes.AttachAttributeToCategory)
		v2Admin.Delete("/categories/:category_id/attributes/:attribute_id", h.UnifiedAttributes.DetachAttributeFromCategory)
		v2Admin.Put("/categories/:category_id/attributes/:attribute_id", h.UnifiedAttributes.UpdateCategoryAttribute)

		// Миграция данных (только для админов)
		v2Admin.Post("/attributes/migrate", h.UnifiedAttributes.MigrateFromLegacy)
		v2Admin.Get("/attributes/migration-status", h.UnifiedAttributes.GetMigrationStatus)

		logger.Info().Msg("V2 unified attributes routes registered successfully")
	} else {
		logger.Info().Msg("V2 unified attributes routes not registered (feature disabled or handler nil)")
	}

	// Обновлено: маршруты API переводов используют обработчик переводов
	translation := app.Group("/api/v1/translation")
	translation.Get("/limits", h.Translations.GetTranslationLimits)
	translation.Post("/provider", h.Translations.SetTranslationProvider)

	// ВАЖНО: НЕ используем Group("/api/v1") с middleware - это вызывает middleware leak!
	// Все защищенные маршруты регистрируем с inline middleware

	// Marketplace protected routes - используем прямую регистрацию с h.jwtParserMW
	authMW := []fiber.Handler{h.jwtParserMW, authMiddleware.RequireAuth()}

	app.Post("/api/v1/c2c/listings", append(authMW, h.Listings.CreateListing)...)
	app.Put("/api/v1/c2c/listings/:id", append(authMW, h.Listings.UpdateListing)...)
	app.Patch("/api/v1/c2c/listings/:id/status", append(authMW, h.Listings.UpdateListingStatus)...)
	app.Delete("/api/v1/c2c/listings/:id", append(authMW, h.Listings.DeleteListing)...)
	app.Post("/api/v1/c2c/listings/check-slug", append(authMW, h.Listings.CheckSlugAvailability)...)
	app.Post("/api/v1/c2c/listings/:id/images", append(authMW, h.Images.UploadImages)...)
	app.Delete("/api/v1/c2c/listings/:id/images/:image_id", append(authMW, h.Images.DeleteImage)...)

	// My listings route - protected
	app.Get("/api/v1/c2c/my-listings", append(authMW, h.Listings.GetMyListings)...)

	// Favorites routes - поддерживаем оба варианта для совместимости
	// Старый формат через listings
	app.Post("/api/v1/c2c/listings/:id/favorite", append(authMW, h.Favorites.AddToFavorites)...)
	app.Delete("/api/v1/c2c/listings/:id/favorite", append(authMW, h.Favorites.RemoveFromFavorites)...)

	// Новый формат - основной
	app.Get("/api/v1/c2c/favorites", append(authMW, h.Favorites.GetFavorites)...)
	app.Get("/api/v1/c2c/favorites/count", append(authMW, h.Favorites.GetFavoritesCount)...)
	app.Post("/api/v1/c2c/favorites/:id", append(authMW, h.Favorites.AddToFavorites)...)
	app.Delete("/api/v1/c2c/favorites/:id", append(authMW, h.Favorites.RemoveFromFavorites)...)

	// Saved searches routes
	app.Post("/api/v1/c2c/saved-searches", append(authMW, h.SavedSearches.CreateSavedSearch)...)
	app.Get("/api/v1/c2c/saved-searches", append(authMW, h.SavedSearches.GetSavedSearches)...)
	app.Get("/api/v1/c2c/saved-searches/:id", append(authMW, h.SavedSearches.GetSavedSearch)...)
	app.Put("/api/v1/c2c/saved-searches/:id", append(authMW, h.SavedSearches.UpdateSavedSearch)...)
	app.Delete("/api/v1/c2c/saved-searches/:id", append(authMW, h.SavedSearches.DeleteSavedSearch)...)
	app.Get("/api/v1/c2c/saved-searches/:id/execute", append(authMW, h.SavedSearches.ExecuteSavedSearch)...)
	app.Get("/api/v1/c2c/favorites/:id/check", append(authMW, h.Favorites.IsInFavorites)...)
	app.Put("/api/v1/c2c/translations/:id", append(authMW, h.Translations.UpdateTranslations)...)
	app.Post("/api/v1/c2c/translations/batch", append(authMW, h.Translations.TranslateText)...)
	app.Post("/api/v1/c2c/moderate-image", append(authMW, h.Images.ModerateImage)...)
	app.Post("/api/v1/c2c/enhance-preview", append(authMW, h.Images.EnhancePreview)...)
	app.Post("/api/v1/c2c/enhance-images", append(authMW, h.Images.EnhanceImages)...)

	// маршруты для новых методов в TranslationsHandler
	app.Post("/api/v1/c2c/translations/batch-translate", append(authMW, h.Translations.BatchTranslateListings)...)
	app.Post("/api/v1/c2c/translations/translate", append(authMW, h.Translations.TranslateText)...)
	app.Post("/api/v1/c2c/translations/detect-language", append(authMW, h.Translations.DetectLanguage)...)
	app.Get("/api/v1/c2c/translations/:id", append(authMW, h.Translations.GetTranslations)...)

	// Регистрируем маршруты для заказов маркетплейса под marketplace префиксом
	if h.Orders != nil {
		// Создаем защищенную группу ТОЛЬКО для orders - узкий префикс!
		ordersGroup := app.Group("/api/v1/c2c/orders", h.jwtParserMW, authMiddleware.RequireAuth())
		h.Orders.RegisterRoutes(ordersGroup)
	}

	// Используем h.jwtParserMW + библиотечный RequireAuthString("admin") для защиты admin роутов
	adminRoutes := app.Group("/api/v1/admin", h.jwtParserMW, authMiddleware.RequireAuthString("admin"))

	// Статистика для админ панели
	adminRoutes.Get("/listings/statistics", h.Listings.GetAdminStatistics)

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
	// Используем отдельный путь чтобы избежать конфликта с :id параметром
	adminRoutes.Get("/categories-all", h.AdminCategories.GetAllCategories)
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
	// ВАЖНО: сначала регистрируем более специфичные маршруты, потом параметризованные
	adminRoutes.Post("/attributes/bulk-update", h.AdminAttributes.BulkUpdateAttributes)
	// Регистрируем variant-compatible до :id маршрута
	if h.VariantMappings != nil {
		adminRoutes.Get("/attributes/variant-compatible", h.VariantMappings.GetVariantCompatibleAttributes)
	}
	adminRoutes.Post("/attributes", h.AdminAttributes.CreateAttribute)
	adminRoutes.Get("/attributes", h.AdminAttributes.GetAttributes)
	adminRoutes.Get("/attributes/:id", h.AdminAttributes.GetAttributeByID)
	adminRoutes.Put("/attributes/:id", h.AdminAttributes.UpdateAttribute)
	adminRoutes.Delete("/attributes/:id", h.AdminAttributes.DeleteAttribute)
	adminRoutes.Post("/attributes/:id/translate", h.AdminAttributes.TranslateAttribute)

	// Маршруты для экспорта/импорта настроек атрибутов
	adminRoutes.Get("/categories/:categoryId/attributes/export", h.AdminAttributes.ExportCategoryAttributes)
	adminRoutes.Post("/categories/:categoryId/attributes/import", h.AdminAttributes.ImportCategoryAttributes)
	adminRoutes.Post("/categories/:targetCategoryId/attributes/copy", h.AdminAttributes.CopyAttributesSettings)

	// Регистрируем маршруты администрирования вариативных атрибутов
	adminRoutes.Get("/variant-attributes", h.AdminVariantAttributes.GetVariantAttributes)
	adminRoutes.Post("/variant-attributes", h.AdminVariantAttributes.CreateVariantAttribute)

	// Новые маршруты для управления вариативными атрибутами через единый интерфейс
	// ВАЖНО: регистрируем ДО :id маршрутов, чтобы избежать конфликтов
	if h.VariantMappings != nil {
		adminRoutes.Get("/variant-attributes/mappings", h.VariantMappings.GetCategoryVariantMappings)
		adminRoutes.Post("/variant-attributes/mappings", h.VariantMappings.CreateVariantMapping)
		adminRoutes.Patch("/variant-attributes/mappings/:id", h.VariantMappings.UpdateVariantMapping)
		adminRoutes.Delete("/variant-attributes/mappings/:id", h.VariantMappings.DeleteVariantMapping)
		adminRoutes.Put("/categories/variant-attributes", h.VariantMappings.UpdateCategoryVariantAttributes)
	}

	// Маршруты с параметрами - регистрируем ПОСЛЕ статичных путей
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

	// Алиасы для обратной совместимости с frontend (c2c-translations → marketplace-translations)
	adminRoutes.Post("/c2c-translations/batch-categories", h.AdminTranslations.BatchTranslateCategories)
	adminRoutes.Post("/c2c-translations/batch-attributes", h.AdminTranslations.BatchTranslateAttributes)
	adminRoutes.Get("/c2c-translations/status", h.AdminTranslations.GetTranslationStatus)
	adminRoutes.Put("/c2c-translations/:entity_type/:entity_id/:field_name", h.AdminTranslations.UpdateFieldTranslation)

	// Управление администраторами

	// Обновлено: маршруты админских функций используют обработчик индексации
	adminRoutes.Post("/reindex-listings", h.Indexing.ReindexAll)
	adminRoutes.Post("/reindex-listings-with-translations", h.Indexing.ReindexAllWithTranslations)
	adminRoutes.Post("/sync-discounts", h.Listings.SynchronizeDiscounts) // Оставляем в Listings, т.к. это работа с объявлениями
	adminRoutes.Post("/reindex-ratings", h.Indexing.ReindexRatings)

	// Chat routes - используем узкий префикс для группы
	chat := app.Group("/api/v1/c2c/chat", h.jwtParserMW, authMiddleware.RequireAuth())
	chat.Get("/", h.Chat.GetChats)
	chat.Get("/messages", h.Chat.GetMessages)
	chat.Get("/messages/:id/translation", h.Chat.TranslateMessage) // NEW: Translation endpoint

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
