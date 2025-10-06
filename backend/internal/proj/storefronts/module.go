package storefronts

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"backend/internal/domain/models"
	"backend/internal/logger"
	"backend/internal/middleware"
	aiHandler "backend/internal/proj/ai/handler"
	"backend/internal/proj/global/service"
	marketplaceServices "backend/internal/proj/marketplace/services"
	storefrontHandler "backend/internal/proj/storefront/handler"
	"backend/internal/proj/storefront/repository"
	"backend/internal/proj/storefronts/handler"
	storefrontService "backend/internal/proj/storefronts/service"
	internalServices "backend/internal/services"
	"backend/internal/storage/postgres"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	authMiddleware "github.com/sveturs/auth/pkg/http/fiber/middleware"
	"go.uber.org/zap"
)

// Module представляет модуль витрин с продуктами
type Module struct {
	services             service.ServicesInterface
	storefrontHandler    *handler.StorefrontHandler
	productHandler       *handler.ProductHandler
	importHandler        *handler.ImportHandler
	imageHandler         *handler.ImageHandler
	dashboardHandler     *handler.DashboardHandler
	variantHandler       *storefrontHandler.VariantHandler
	publicVariantHandler *storefrontHandler.PublicVariantHandler
	aiProductHandler     *handler.AIProductHandler
	importQueueManager   *storefrontService.ImportQueueManager
}

// NewModule создает новый модуль витрин
func NewModule(ctx context.Context, services service.ServicesInterface) *Module {
	storefrontSvc := services.Storefront()
	// Создаем адаптер для storage
	db := services.Storage().(*postgres.Database)
	productStorage := &storageAdapter{db: db}
	// Получаем OpenSearch репозиторий для товаров витрин
	var searchRepo storefrontService.ProductSearchRepository
	if osProductRepo := services.Storage().StorefrontProductSearch(); osProductRepo != nil {
		searchRepo = osProductRepo.(storefrontService.ProductSearchRepository)
	}

	// Создаем variant repository и service
	variantRepo := repository.NewVariantRepository(db.GetSQLXDB())
	variantService := storefrontService.NewVariantService(variantRepo)

	// Создаем ProductService с VariantService
	productSvc := storefrontService.NewProductService(productStorage, searchRepo, variantService)

	// Создаем единый ImageService с конфигурацией buckets
	imageRepo := postgres.NewImageRepository(db.GetSQLXDB())
	cfg := services.Config()
	imageCfg := internalServices.ImageServiceConfig{
		BucketListings:    cfg.FileStorage.MinioBucketName,
		BucketStorefront:  cfg.FileStorage.MinioStorefrontBucket,
		BucketChatFiles:   cfg.FileStorage.MinioChatBucket,
		BucketReviewPhoto: cfg.FileStorage.MinioReviewPhotosBucket,
	}
	imageService := services.NewImageService(services.FileStorage(), imageRepo, imageCfg)

	// Создаем repository для import jobs
	importJobsRepo := postgres.NewImportJobsRepository(db.GetPool())

	// Создаем repository для category mappings
	categoryMappingsRepo := postgres.NewCategoryMappingsRepository(db.GetPool())

	// Создаем сервис импорта (categoryMappingService будет установлен позже)
	importSvc := storefrontService.NewImportService(productSvc, importJobsRepo, imageService, nil)

	// Создаем Import Queue Manager для асинхронной обработки импорта
	// Параметры: workerCount (количество воркеров), queueSize (размер очереди)
	importQueueManager := storefrontService.NewImportQueueManager(4, 100, importSvc)

	// Устанавливаем queue manager в ImportService
	importSvc.SetQueueManager(importQueueManager)

	// Запускаем queue manager
	if err := importQueueManager.Start(); err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to start import queue manager")
		// Продолжаем работу без queue manager (fallback на синхронный импорт)
	} else {
		logger.Info().Msg("Import queue manager started successfully")
	}

	// Получаем storefront repository
	storefrontRepo := services.Storage().Storefront().(postgres.StorefrontRepository)

	// Создаем variant handlers (используем уже созданный variantRepo)
	variantHandler := storefrontHandler.NewVariantHandler(variantRepo)
	publicVariantHandler := storefrontHandler.NewPublicVariantHandler(variantRepo)

	// Инициализируем AI Product Handler (переиспользуем marketplace AI infrastructure)
	// Создаем AICategoryDetector (он легковесный, использует тот же DB и Redis cache)
	log := logger.Info() // Используем logger из internal/logger
	_ = log              // Избегаем warning о неиспользуемой переменной

	aiDetector := marketplaceServices.NewAICategoryDetector(ctx, db.GetSQLXDB(), zap.L())

	// Создаем CategoryMappingService с AI detector (после создания aiDetector)
	categoryMappingSvc := storefrontService.NewCategoryMappingService(categoryMappingsRepo, aiDetector, zap.L())

	// Устанавливаем CategoryMappingService в ImportService
	importSvc.SetCategoryMappingService(categoryMappingSvc)

	// Создаем AI Category Mapper и Analyzer для умного импорта
	aiCategoryMapper := storefrontService.NewAICategoryMapper(aiDetector, categoryMappingSvc)
	aiCategoryAnalyzer := storefrontService.NewAICategoryAnalyzer(aiCategoryMapper, categoryMappingSvc)

	// Создаем AI Handler для общих AI операций (анализ изображений, переводы и т.д.)
	aiHandlerInstance := aiHandler.NewHandler(cfg, services)

	// Создаем AI Product Handler для витрин
	aiProductHandlerInstance := handler.NewAIProductHandler(
		aiDetector,
		aiHandlerInstance,
		zap.L(),
	)

	return &Module{
		services:             services,
		storefrontHandler:    handler.NewStorefrontHandler(storefrontSvc),
		productHandler:       handler.NewProductHandler(productSvc),
		importHandler:        handler.NewImportHandler(importSvc, aiCategoryMapper, aiCategoryAnalyzer),
		imageHandler:         handler.NewImageHandler(imageService, productSvc),
		dashboardHandler:     handler.NewDashboardHandler(storefrontSvc, productSvc, storefrontRepo),
		variantHandler:       variantHandler,
		publicVariantHandler: publicVariantHandler,
		aiProductHandler:     aiProductHandlerInstance,
		importQueueManager:   importQueueManager,
	}
}

// Shutdown gracefully shuts down the module
func (m *Module) Shutdown() error {
	logger.Info().Msg("Shutting down storefronts module...")

	// Stop import queue manager
	if m.importQueueManager != nil && m.importQueueManager.IsRunning() {
		if err := m.importQueueManager.Stop(); err != nil {
			logger.Error().
				Err(err).
				Msg("Failed to stop import queue manager")
			return err
		}
		logger.Info().Msg("Import queue manager stopped successfully")
	}

	logger.Info().Msg("Storefronts module shutdown complete")
	return nil
}

// GetPrefix возвращает префикс для маршрутов
func (m *Module) GetPrefix() string {
	return "/api/v1/storefronts"
}

// RegisterRoutes регистрирует все маршруты модуля
func (m *Module) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
	api := app.Group("/api/v1")

	// Публичные маршруты витрин (без авторизации) - РЕГИСТРИРУЕМ ПЕРВЫМИ
	// ВАЖНО: Конкретные маршруты должны быть ПЕРЕД параметризованными

	// Защищенный маршрут /my - специфичный, регистрируем рано
	api.Get("/storefronts/my", mw.JWTParser(), authMiddleware.RequireAuth(), m.storefrontHandler.GetMyStorefronts)

	// Публичный endpoint для получения товара по ID (для чата)
	api.Get("/storefronts/products/:id", m.productHandler.GetProductByID)

	// Списки и поиск (публичные)
	api.Get("/storefronts", m.storefrontHandler.ListStorefronts)
	api.Get("/storefronts/search", m.storefrontHandler.SearchOpenSearch)
	api.Get("/storefronts/nearby", m.storefrontHandler.GetNearbyStorefronts)

	// Картографические данные (публичные)
	api.Get("/storefronts/map", m.storefrontHandler.GetMapData)
	api.Get("/storefronts/building", m.storefrontHandler.GetBusinessesInBuilding)

	// Маршруты с slug (публичные)
	api.Get("/storefronts/slug/:slug", m.storefrontHandler.GetStorefrontBySlug)
	api.Get("/storefronts/slug/:slug/products", m.getProductsBySlug)
	api.Get("/storefronts/slug/:slug/products/:id", m.getProductBySlug)

	// Публичные маршруты для вариантов товаров
	api.Get("/storefronts/slug/:slug/products/:product_id/variants", m.publicVariantHandler.GetProductVariantsPublic)
	api.Get("/storefronts/:slug/products/:product_id", m.publicVariantHandler.GetProductPublic)
	api.Get("/storefronts/variants/:variant_id", m.publicVariantHandler.GetVariantByIDPublic)

	// Параметризованные маршруты (должны быть последними)
	// Получение витрины (публичное)
	api.Get("/storefronts/:id", m.storefrontHandler.GetStorefront)

	// Персонал (просмотр, публичное)
	api.Get("/storefronts/:id/staff", m.storefrontHandler.GetStaff)

	// Аналитика (запись просмотра, публичное)
	api.Post("/storefronts/:id/view", m.storefrontHandler.RecordView)

	// Защищенные маршруты витрин (требуют авторизации)
	// ВАЖНО: Регистрируем напрямую через api вместо группы, чтобы избежать конфликта с public группой
	{
		// AI endpoints для витрин (защищенные)
		if m.aiProductHandler != nil {
			api.Post("/storefronts/ai/analyze-product-image", mw.JWTParser(), authMiddleware.RequireAuth(), m.aiProductHandler.AnalyzeProductImage)
			api.Post("/storefronts/ai/detect-category", mw.JWTParser(), authMiddleware.RequireAuth(), m.aiProductHandler.DetectCategory)
			api.Post("/storefronts/ai/ab-test-titles", mw.JWTParser(), authMiddleware.RequireAuth(), m.aiProductHandler.ABTestTitles)
			api.Post("/storefronts/ai/translate-content", mw.JWTParser(), authMiddleware.RequireAuth(), m.aiProductHandler.TranslateContent)
			api.Get("/storefronts/ai/metrics", mw.JWTParser(), authMiddleware.RequireAuth(), m.aiProductHandler.GetMetrics)
		}

		// Управление витринами
		// ВАЖНО: Используем /storefronts/create вместо /storefronts, чтобы избежать конфликта с GET /storefronts
		api.Post("/storefronts/create", mw.JWTParser(), authMiddleware.RequireAuth(), m.storefrontHandler.CreateStorefront)
		api.Put("/storefronts/:id", mw.JWTParser(), authMiddleware.RequireAuth(), m.storefrontHandler.UpdateStorefront)
		api.Delete("/storefronts/:id", mw.JWTParser(), authMiddleware.RequireAuth(), m.storefrontHandler.DeleteStorefront)
		api.Post("/storefronts/:id/restore", mw.JWTParser(), authMiddleware.RequireAuth(), m.storefrontHandler.RestoreStorefront)

		// Настройки витрины
		api.Put("/storefronts/:id/hours", mw.JWTParser(), authMiddleware.RequireAuth(), m.storefrontHandler.UpdateWorkingHours)
		api.Put("/storefronts/:id/payment-methods", mw.JWTParser(), authMiddleware.RequireAuth(), m.storefrontHandler.UpdatePaymentMethods)
		api.Put("/storefronts/:id/delivery-options", mw.JWTParser(), authMiddleware.RequireAuth(), m.storefrontHandler.UpdateDeliveryOptions)

		// Управление персоналом
		api.Post("/storefronts/:id/staff", mw.JWTParser(), authMiddleware.RequireAuth(), m.storefrontHandler.AddStaff)
		api.Put("/storefronts/:id/staff/:staffId/permissions", mw.JWTParser(), authMiddleware.RequireAuth(), m.storefrontHandler.UpdateStaffPermissions)
		api.Delete("/storefronts/:id/staff/:userId", mw.JWTParser(), authMiddleware.RequireAuth(), m.storefrontHandler.RemoveStaff)

		// Загрузка изображений
		api.Post("/storefronts/:id/logo", mw.JWTParser(), authMiddleware.RequireAuth(), m.storefrontHandler.UploadLogo)
		api.Post("/storefronts/:id/banner", mw.JWTParser(), authMiddleware.RequireAuth(), m.storefrontHandler.UploadBanner)

		// Аналитика (просмотр)
		api.Get("/storefronts/:id/analytics", mw.JWTParser(), authMiddleware.RequireAuth(), m.storefrontHandler.GetAnalytics)

		// Защищенные маршруты товаров с использованием slug
		api.Post("/storefronts/slug/:slug/products", mw.JWTParser(), authMiddleware.RequireAuth(), m.createProductBySlug)
		api.Put("/storefronts/slug/:slug/products/:id", mw.JWTParser(), authMiddleware.RequireAuth(), m.updateProductBySlug)
		api.Delete("/storefronts/slug/:slug/products/:id", mw.JWTParser(), authMiddleware.RequireAuth(), m.deleteProductBySlug)
		api.Post("/storefronts/slug/:slug/products/:id/inventory", mw.JWTParser(), authMiddleware.RequireAuth(), m.updateInventoryBySlug)
		api.Get("/storefronts/slug/:slug/products/stats", mw.JWTParser(), authMiddleware.RequireAuth(), m.getProductStatsBySlug)

		// Маршруты для изображений товаров
		api.Post("/storefronts/slug/:slug/products/:product_id/images", mw.JWTParser(), authMiddleware.RequireAuth(), m.uploadProductImageBySlug)
		api.Get("/storefronts/slug/:slug/products/:product_id/images", mw.JWTParser(), authMiddleware.RequireAuth(), m.getProductImagesBySlug)
		api.Delete("/storefronts/slug/:slug/products/:product_id/images/:image_id", mw.JWTParser(), authMiddleware.RequireAuth(), m.deleteProductImageBySlug)
		api.Post("/storefronts/slug/:slug/products/:product_id/images/:image_id/main", mw.JWTParser(), authMiddleware.RequireAuth(), m.setMainProductImageBySlug)
		api.Put("/storefronts/slug/:slug/products/:product_id/images/order", mw.JWTParser(), authMiddleware.RequireAuth(), m.updateImageOrderBySlug)

		// Bulk операции с товарами
		api.Post("/storefronts/slug/:slug/products/bulk/create", mw.JWTParser(), authMiddleware.RequireAuth(), m.bulkCreateProductsBySlug)
		api.Put("/storefronts/slug/:slug/products/bulk/update", mw.JWTParser(), authMiddleware.RequireAuth(), m.bulkUpdateProductsBySlug)
		api.Delete("/storefronts/slug/:slug/products/bulk/delete", mw.JWTParser(), authMiddleware.RequireAuth(), m.bulkDeleteProductsBySlug)
		api.Put("/storefronts/slug/:slug/products/bulk/status", mw.JWTParser(), authMiddleware.RequireAuth(), m.bulkUpdateStatusBySlug)

		// Маршруты импорта товаров
		api.Post("/storefronts/:storefront_id/import/url", mw.JWTParser(), authMiddleware.RequireAuth(), m.importHandler.ImportFromURL)
		api.Post("/storefronts/:storefront_id/import/file", mw.JWTParser(), authMiddleware.RequireAuth(), m.importHandler.ImportFromFile)
		api.Post("/storefronts/:storefront_id/import/validate", mw.JWTParser(), authMiddleware.RequireAuth(), m.importHandler.ValidateImportFile)
		api.Post("/storefronts/:storefront_id/import/preview", mw.JWTParser(), authMiddleware.RequireAuth(), m.importHandler.PreviewImportFile)

		// Новые маршруты для AI анализа импорта
		api.Post("/storefronts/:storefront_id/import/analyze-categories", mw.JWTParser(), authMiddleware.RequireAuth(), m.importHandler.AnalyzeCategories)
		api.Post("/storefronts/:storefront_id/import/analyze-attributes", mw.JWTParser(), authMiddleware.RequireAuth(), m.importHandler.AnalyzeAttributes)
		api.Post("/storefronts/:storefront_id/import/detect-variants", mw.JWTParser(), authMiddleware.RequireAuth(), m.importHandler.DetectVariants)
		api.Post("/storefronts/:storefront_id/import/analyze-client-categories", mw.JWTParser(), authMiddleware.RequireAuth(), m.importHandler.AnalyzeClientCategories)

		api.Get("/storefronts/:storefront_id/import/jobs", mw.JWTParser(), authMiddleware.RequireAuth(), m.importHandler.GetJobs)
		api.Get("/storefronts/:storefront_id/import/jobs/:jobId", mw.JWTParser(), authMiddleware.RequireAuth(), m.importHandler.GetJobDetails)
		api.Get("/storefronts/:storefront_id/import/jobs/:jobId/status", mw.JWTParser(), authMiddleware.RequireAuth(), m.importHandler.GetJobStatus)
		api.Post("/storefronts/:storefront_id/import/jobs/:jobId/cancel", mw.JWTParser(), authMiddleware.RequireAuth(), m.importHandler.CancelJob)
		api.Post("/storefronts/:storefront_id/import/jobs/:jobId/retry", mw.JWTParser(), authMiddleware.RequireAuth(), m.importHandler.RetryJob)

		// Маршруты импорта через slug
		api.Post("/storefronts/slug/:slug/import/url", mw.JWTParser(), authMiddleware.RequireAuth(), m.importFromURLBySlug)
		api.Post("/storefronts/slug/:slug/import/file", mw.JWTParser(), authMiddleware.RequireAuth(), m.importFromFileBySlug)
		api.Post("/storefronts/slug/:slug/import/validate", mw.JWTParser(), authMiddleware.RequireAuth(), m.validateImportBySlug)
		api.Post("/storefronts/slug/:slug/import/preview", mw.JWTParser(), authMiddleware.RequireAuth(), m.previewImportBySlug)
		api.Get("/storefronts/slug/:slug/import/jobs", mw.JWTParser(), authMiddleware.RequireAuth(), m.getJobsBySlug)
		api.Get("/storefronts/slug/:slug/import/jobs/:jobId", mw.JWTParser(), authMiddleware.RequireAuth(), m.getJobDetailsBySlug)
		api.Get("/storefronts/slug/:slug/import/jobs/:jobId/status", mw.JWTParser(), authMiddleware.RequireAuth(), m.getJobStatusBySlug)
		api.Post("/storefronts/slug/:slug/import/jobs/:jobId/cancel", mw.JWTParser(), authMiddleware.RequireAuth(), m.cancelJobBySlug)
		api.Post("/storefronts/slug/:slug/import/jobs/:jobId/retry", mw.JWTParser(), authMiddleware.RequireAuth(), m.retryJobBySlug)

		// Маршруты товаров по slug без префикса (для frontend совместимости)
		api.Put("/storefronts/:slug/products/:id", mw.JWTParser(), authMiddleware.RequireAuth(), m.updateProductBySlugDirect)
		api.Post("/storefronts/:slug/products/:product_id/images", mw.JWTParser(), authMiddleware.RequireAuth(), m.uploadProductImageBySlugDirect)
		api.Delete("/storefronts/:slug/products/:product_id/images/:image_id", mw.JWTParser(), authMiddleware.RequireAuth(), m.deleteProductImageBySlugDirect)
		api.Post("/storefronts/:slug/products/:product_id/images/:image_id/main", mw.JWTParser(), authMiddleware.RequireAuth(), m.setMainProductImageBySlugDirect)

		// Dashboard маршруты
		api.Get("/storefronts/:slug/dashboard/stats", mw.JWTParser(), authMiddleware.RequireAuth(), m.dashboardHandler.GetDashboardStats)
		api.Get("/storefronts/:slug/dashboard/recent-orders", mw.JWTParser(), authMiddleware.RequireAuth(), m.dashboardHandler.GetRecentOrders)
		api.Get("/storefronts/:slug/dashboard/low-stock", mw.JWTParser(), authMiddleware.RequireAuth(), m.dashboardHandler.GetLowStockProducts)
		api.Get("/storefronts/:slug/dashboard/notifications", mw.JWTParser(), authMiddleware.RequireAuth(), m.dashboardHandler.GetDashboardNotifications)

		// Приватные маршруты для управления вариантами товаров (требуют авторизации)
		// Глобальные атрибуты вариантов
		api.Get("/storefronts/storefront/variants/attributes", mw.JWTParser(), authMiddleware.RequireAuth(), m.variantHandler.GetVariantAttributes)
		api.Get("/storefronts/storefront/variants/attributes/:attribute_id/values", mw.JWTParser(), authMiddleware.RequireAuth(), m.variantHandler.GetVariantAttributeValues)

		// Управление вариантами товара
		api.Get("/storefronts/storefront/products/:product_id/variants", mw.JWTParser(), authMiddleware.RequireAuth(), m.variantHandler.GetVariantsByProductID)
		api.Post("/storefronts/storefront/variants", mw.JWTParser(), authMiddleware.RequireAuth(), m.variantHandler.CreateVariant)
		api.Put("/storefronts/storefront/variants/:variant_id", mw.JWTParser(), authMiddleware.RequireAuth(), m.variantHandler.UpdateVariant)
		api.Delete("/storefronts/storefront/variants/:variant_id", mw.JWTParser(), authMiddleware.RequireAuth(), m.variantHandler.DeleteVariant)
		api.Get("/storefronts/storefront/variants/:variant_id", mw.JWTParser(), authMiddleware.RequireAuth(), m.variantHandler.GetVariantByID)

		// Генерация и массовые операции
		api.Post("/storefronts/storefront/variants/generate", mw.JWTParser(), authMiddleware.RequireAuth(), m.variantHandler.GenerateVariants)
		api.Post("/storefronts/storefront/variants/bulk", mw.JWTParser(), authMiddleware.RequireAuth(), m.variantHandler.BulkCreateVariants)
		api.Get("/storefronts/storefront/products/:product_id/variant-matrix", mw.JWTParser(), authMiddleware.RequireAuth(), m.variantHandler.GetVariantMatrix)
		api.Post("/storefronts/storefront/products/:product_id/variants/bulk-update-stock", mw.JWTParser(), authMiddleware.RequireAuth(), m.variantHandler.BulkUpdateStock)
		api.Get("/storefronts/storefront/products/:product_id/variants/analytics", mw.JWTParser(), authMiddleware.RequireAuth(), m.variantHandler.GetVariantAnalytics)

		// CSV импорт/экспорт
		api.Post("/storefronts/storefront/products/:product_id/variants/import", mw.JWTParser(), authMiddleware.RequireAuth(), m.variantHandler.ImportVariants)
		api.Get("/storefronts/storefront/products/:product_id/variants/export", mw.JWTParser(), authMiddleware.RequireAuth(), m.variantHandler.ExportVariants)

		// Настройка атрибутов товара
		api.Post("/storefronts/storefront/products/attributes/setup", mw.JWTParser(), authMiddleware.RequireAuth(), m.variantHandler.SetupProductAttributes)
		api.Get("/storefronts/storefront/products/:product_id/attributes", mw.JWTParser(), authMiddleware.RequireAuth(), m.variantHandler.GetProductAttributes)
		api.Get("/storefronts/storefront/categories/:category_id/attributes", mw.JWTParser(), authMiddleware.RequireAuth(), m.variantHandler.GetAvailableAttributesForCategory)
	}

	// Публичные маршруты для вариантов (БЕЗ АВТОРИЗАЦИИ)
	publicVariants := api.Group("/public")
	{
		// Публичные endpoints для покупателей
		publicVariants.Get("/storefronts/:slug/products/:product_id", m.publicVariantHandler.GetProductPublic)
		publicVariants.Get("/storefronts/:slug/products/:product_id/variants", m.publicVariantHandler.GetProductVariantsPublic)
		publicVariants.Get("/variants/attributes", m.publicVariantHandler.GetVariantAttributesPublic)
		publicVariants.Get("/variants/attributes/:attribute_id/values", m.publicVariantHandler.GetVariantAttributeValuesPublic)
		publicVariants.Get("/variants/:variant_id", m.publicVariantHandler.GetVariantByIDPublic)
		publicVariants.Get("/categories/:category_id/attributes", m.publicVariantHandler.GetAvailableAttributesForCategoryPublic)
	}

	// Публичные маршруты импорта (для получения шаблонов и документации)
	api.Get("/storefronts/import/csv-template", m.importHandler.GetCSVTemplate)
	api.Get("/storefronts/import/formats", m.importHandler.GetImportFormats)

	return nil
}

// Wrapper функции для добавления storefrontID в контекст по slug

func (m *Module) getProductsBySlug(c *fiber.Ctx) error {
	if err := m.setStorefrontIDBySlug(c); err != nil {
		return err
	}
	return m.productHandler.GetProducts(c)
}

func (m *Module) getProductBySlug(c *fiber.Ctx) error {
	if err := m.setStorefrontIDBySlug(c); err != nil {
		return err
	}
	return m.productHandler.GetProduct(c)
}

func (m *Module) createProductBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")

	// Декодируем slug на случай если он URL-encoded
	decodedSlug, err := url.QueryUnescape(slug)
	if err == nil {
		slug = decodedSlug
	}

	logger.Info().
		Str("slug", slug).
		Str("method", "createProductBySlug").
		Msg("Starting product creation")

	if slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Storefront slug is required",
		})
	}

	// Получаем витрину через сервис по slug
	storefront, err := m.services.Storefront().GetBySlug(c.Context(), slug)
	if err != nil {
		logger.Error().
			Err(err).
			Str("slug", slug).
			Msg("Failed to get storefront")

		if errors.Is(err, postgres.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Storefront not found",
				"slug":  slug,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get storefront: " + err.Error(),
			"slug":  slug,
		})
	}

	logger.Info().
		Int("storefront_id", storefront.ID).
		Int("owner_id", storefront.UserID).
		Msg("Found storefront")

	// Сохраняем ID в контекст
	c.Locals("storefrontID", storefront.ID)

	// Проверяем доступ
	userID, ok := authMiddleware.GetUserID(c)
	if !ok {
		logger.Error().
			Msg("User ID not found in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User ID not found in context",
		})
	}

	logger.Info().
		Int("user_id", userID).
		Int("storefront_owner_id", storefront.UserID).
		Bool("is_owner", storefront.UserID == userID).
		Msg("Checking access")

	// Проверяем, является ли пользователь владельцем или персоналом
	if storefront.UserID != userID {
		// Проверяем, есть ли пользователь в персонале
		staff, err := m.services.Storefront().GetStaff(c.Context(), storefront.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to check permissions: " + err.Error(),
			})
		}

		// Ищем пользователя среди персонала с правами на продукты
		hasAccess := false
		for _, s := range staff {
			if s.UserID == userID {
				// Проверяем права на управление продуктами
				permissions, ok := s.Permissions["products"].(bool)
				if ok && permissions {
					hasAccess = true
					break
				}
			}
		}

		if !hasAccess {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":             "Access denied: user is not owner or staff with product permissions",
				"storefrontOwnerID": storefront.UserID,
				"currentUserID":     userID,
			})
		}
	}

	return m.productHandler.CreateProduct(c)
}

func (m *Module) updateProductBySlug(c *fiber.Ctx) error {
	if err := m.setStorefrontIDBySlug(c); err != nil {
		return err
	}

	// Проверяем доступ после установки storefrontID
	if err := m.checkStorefrontAccess(c); err != nil {
		return err
	}

	return m.productHandler.UpdateProduct(c)
}

func (m *Module) deleteProductBySlug(c *fiber.Ctx) error {
	if err := m.setStorefrontIDBySlug(c); err != nil {
		return err
	}

	// Проверяем доступ после установки storefrontID
	if err := m.checkStorefrontAccess(c); err != nil {
		return err
	}

	// Проверяем параметр hard для выбора типа удаления
	if c.Query("hard") == "true" {
		return m.productHandler.HardDeleteProduct(c)
	}

	return m.productHandler.DeleteProduct(c)
}

func (m *Module) updateInventoryBySlug(c *fiber.Ctx) error {
	if err := m.setStorefrontIDBySlug(c); err != nil {
		return err
	}

	// Проверяем доступ после установки storefrontID
	if err := m.checkStorefrontAccess(c); err != nil {
		return err
	}

	return m.productHandler.UpdateInventory(c)
}

func (m *Module) getProductStatsBySlug(c *fiber.Ctx) error {
	if err := m.setStorefrontIDBySlug(c); err != nil {
		return err
	}

	// Проверяем доступ после установки storefrontID
	if err := m.checkStorefrontAccess(c); err != nil {
		return err
	}

	return m.productHandler.GetProductStats(c)
}

// setStorefrontIDBySlug добавляет storefront ID в контекст по slug или ID
func (m *Module) setStorefrontIDBySlug(c *fiber.Ctx) error {
	slugOrID := c.Params("slug")
	if slugOrID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Storefront slug is required",
		})
	}

	// Декодируем slug на случай если он URL-encoded (для кириллических slug)
	decodedSlug, err := url.QueryUnescape(slugOrID)
	if err == nil {
		slugOrID = decodedSlug
	}

	var storefront *models.Storefront

	// Пробуем сначала как ID
	if id, parseErr := strconv.Atoi(slugOrID); parseErr == nil {
		storefront, err = m.services.Storefront().GetByID(c.Context(), id)
		if err == nil {
			c.Locals("storefrontID", storefront.ID)
			return nil
		}
		// Если не нашли по ID, пробуем как slug
	}

	// Получаем витрину через сервис по slug
	storefront, err = m.services.Storefront().GetBySlug(c.Context(), slugOrID)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Storefront not found",
				"slug":  slugOrID,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get storefront: " + err.Error(),
			"slug":  slugOrID,
		})
	}

	c.Locals("storefrontID", storefront.ID)
	return nil
}

// checkStorefrontAccess проверяет доступ к витрине
func (m *Module) checkStorefrontAccess(c *fiber.Ctx) error {
	storefrontID, ok := c.Locals("storefrontID").(int)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Storefront ID not found in context",
		})
	}

	userID, ok := authMiddleware.GetUserID(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User ID not found in context",
		})
	}

	// Проверяем, является ли пользователь администратором
	isAdmin := authMiddleware.IsAdmin(c)
	if isAdmin {
		// Администраторы имеют полный доступ ко всем витринам
		return nil
	}

	// Получаем витрину для проверки владельца
	storefront, err := m.services.Storefront().GetByID(c.Context(), storefrontID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get storefront",
		})
	}

	// Проверяем, является ли пользователь владельцем или персоналом
	if storefront.UserID != userID {
		// Проверяем, есть ли пользователь в персонале
		staff, err := m.services.Storefront().GetStaff(c.Context(), storefrontID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to check permissions: " + err.Error(),
			})
		}

		// Ищем пользователя среди персонала с правами на продукты
		hasAccess := false
		for _, s := range staff {
			if s.UserID == userID {
				// Проверяем права на управление продуктами
				permissions, ok := s.Permissions["products"].(bool)
				if ok && permissions {
					hasAccess = true
					break
				}
			}
		}

		if !hasAccess {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":             "Access denied: user is not owner or staff with product permissions",
				"storefrontOwnerID": storefront.UserID,
				"currentUserID":     userID,
			})
		}
	}

	return nil
}

// Функции-обертки для импорта через slug
func (m *Module) importFromURLBySlug(c *fiber.Ctx) error {
	if err := m.setStorefrontIDBySlug(c); err != nil {
		return err
	}

	// Проверяем доступ после установки storefrontID
	if err := m.checkStorefrontAccess(c); err != nil {
		return err
	}

	// Storefront ID уже в locals, ImportHandler его оттуда возьмет
	return m.importHandler.ImportFromURL(c)
}

func (m *Module) importFromFileBySlug(c *fiber.Ctx) error {
	if err := m.setStorefrontIDBySlug(c); err != nil {
		return err
	}

	// Проверяем доступ после установки storefrontID
	if err := m.checkStorefrontAccess(c); err != nil {
		return err
	}

	// Storefront ID уже в locals, ImportHandler его оттуда возьмет
	return m.importHandler.ImportFromFile(c)
}

func (m *Module) validateImportBySlug(c *fiber.Ctx) error {
	if err := m.setStorefrontIDBySlug(c); err != nil {
		return err
	}

	// Проверяем доступ после установки storefrontID
	if err := m.checkStorefrontAccess(c); err != nil {
		return err
	}

	// Storefront ID уже в locals, ImportHandler его оттуда возьмет
	return m.importHandler.ValidateImportFile(c)
}

func (m *Module) previewImportBySlug(c *fiber.Ctx) error {
	if err := m.setStorefrontIDBySlug(c); err != nil {
		return err
	}

	// Проверяем доступ после установки storefrontID
	if err := m.checkStorefrontAccess(c); err != nil {
		return err
	}

	// Storefront ID уже в locals, ImportHandler его оттуда возьмет
	return m.importHandler.PreviewImportFile(c)
}

// Функции-обертки для работы с jobs через slug
func (m *Module) getJobsBySlug(c *fiber.Ctx) error {
	if err := m.setStorefrontIDBySlug(c); err != nil {
		return err
	}

	// Проверяем доступ после установки storefrontID
	if err := m.checkStorefrontAccess(c); err != nil {
		return err
	}

	// Storefront ID уже в locals, ImportHandler его оттуда возьмет
	return m.importHandler.GetJobs(c)
}

func (m *Module) getJobDetailsBySlug(c *fiber.Ctx) error {
	if err := m.setStorefrontIDBySlug(c); err != nil {
		return err
	}

	// Проверяем доступ после установки storefrontID
	if err := m.checkStorefrontAccess(c); err != nil {
		return err
	}

	return m.importHandler.GetJobDetails(c)
}

func (m *Module) getJobStatusBySlug(c *fiber.Ctx) error {
	if err := m.setStorefrontIDBySlug(c); err != nil {
		return err
	}

	// Проверяем доступ после установки storefrontID
	if err := m.checkStorefrontAccess(c); err != nil {
		return err
	}

	return m.importHandler.GetJobStatus(c)
}

func (m *Module) cancelJobBySlug(c *fiber.Ctx) error {
	if err := m.setStorefrontIDBySlug(c); err != nil {
		return err
	}

	// Проверяем доступ после установки storefrontID
	if err := m.checkStorefrontAccess(c); err != nil {
		return err
	}

	return m.importHandler.CancelJob(c)
}

func (m *Module) retryJobBySlug(c *fiber.Ctx) error {
	if err := m.setStorefrontIDBySlug(c); err != nil {
		return err
	}

	// Проверяем доступ после установки storefrontID
	if err := m.checkStorefrontAccess(c); err != nil {
		return err
	}

	return m.importHandler.RetryJob(c)
}

// Функции-обертки для bulk операций с товарами
func (m *Module) bulkCreateProductsBySlug(c *fiber.Ctx) error {
	if err := m.setStorefrontIDBySlug(c); err != nil {
		return err
	}

	// Проверяем доступ после установки storefrontID
	if err := m.checkStorefrontAccess(c); err != nil {
		return err
	}

	// Storefront ID уже в locals, ProductHandler его оттуда возьмет
	return m.productHandler.BulkCreateProducts(c)
}

func (m *Module) bulkUpdateProductsBySlug(c *fiber.Ctx) error {
	if err := m.setStorefrontIDBySlug(c); err != nil {
		return err
	}

	// Проверяем доступ после установки storefrontID
	if err := m.checkStorefrontAccess(c); err != nil {
		return err
	}

	// Storefront ID уже в locals, ProductHandler его оттуда возьмет
	return m.productHandler.BulkUpdateProducts(c)
}

func (m *Module) bulkDeleteProductsBySlug(c *fiber.Ctx) error {
	if err := m.setStorefrontIDBySlug(c); err != nil {
		return err
	}

	// Проверяем доступ после установки storefrontID
	if err := m.checkStorefrontAccess(c); err != nil {
		return err
	}

	// Storefront ID уже в locals, ProductHandler его оттуда возьмет
	return m.productHandler.BulkDeleteProducts(c)
}

func (m *Module) bulkUpdateStatusBySlug(c *fiber.Ctx) error {
	if err := m.setStorefrontIDBySlug(c); err != nil {
		return err
	}

	// Проверяем доступ после установки storefrontID
	if err := m.checkStorefrontAccess(c); err != nil {
		return err
	}

	// Storefront ID уже в locals, ProductHandler его оттуда возьмет
	return m.productHandler.BulkUpdateStatus(c)
}

// Функции-обертки для изображений товаров
func (m *Module) uploadProductImageBySlug(c *fiber.Ctx) error {
	if err := m.setStorefrontIDBySlug(c); err != nil {
		return err
	}

	// Проверяем доступ после установки storefrontID
	if err := m.checkStorefrontAccess(c); err != nil {
		return err
	}

	return m.imageHandler.UploadProductImage(c)
}

func (m *Module) getProductImagesBySlug(c *fiber.Ctx) error {
	if err := m.setStorefrontIDBySlug(c); err != nil {
		return err
	}

	// Для получения изображений не требуется авторизация - это публичный API
	return m.imageHandler.GetProductImages(c)
}

func (m *Module) deleteProductImageBySlug(c *fiber.Ctx) error {
	if err := m.setStorefrontIDBySlug(c); err != nil {
		return err
	}

	// Проверяем доступ после установки storefrontID
	if err := m.checkStorefrontAccess(c); err != nil {
		return err
	}

	return m.imageHandler.DeleteProductImage(c)
}

func (m *Module) setMainProductImageBySlug(c *fiber.Ctx) error {
	if err := m.setStorefrontIDBySlug(c); err != nil {
		return err
	}

	// Проверяем доступ после установки storefrontID
	if err := m.checkStorefrontAccess(c); err != nil {
		return err
	}

	return m.imageHandler.SetMainProductImage(c)
}

func (m *Module) updateImageOrderBySlug(c *fiber.Ctx) error {
	if err := m.setStorefrontIDBySlug(c); err != nil {
		return err
	}

	// Проверяем доступ после установки storefrontID
	if err := m.checkStorefrontAccess(c); err != nil {
		return err
	}

	return m.imageHandler.UpdateImageOrder(c)
}

// Direct slug wrapper functions for frontend compatibility (без префикса /slug/)
func (m *Module) updateProductBySlugDirect(c *fiber.Ctx) error {
	if err := m.setStorefrontIDBySlug(c); err != nil {
		return err
	}

	// Проверяем доступ после установки storefrontID
	if err := m.checkStorefrontAccess(c); err != nil {
		return err
	}

	return m.productHandler.UpdateProduct(c)
}

func (m *Module) uploadProductImageBySlugDirect(c *fiber.Ctx) error {
	if err := m.setStorefrontIDBySlug(c); err != nil {
		return err
	}

	// Проверяем доступ после установки storefrontID
	if err := m.checkStorefrontAccess(c); err != nil {
		return err
	}

	return m.imageHandler.UploadProductImage(c)
}

func (m *Module) deleteProductImageBySlugDirect(c *fiber.Ctx) error {
	if err := m.setStorefrontIDBySlug(c); err != nil {
		return err
	}

	// Проверяем доступ после установки storefrontID
	if err := m.checkStorefrontAccess(c); err != nil {
		return err
	}

	return m.imageHandler.DeleteProductImage(c)
}

func (m *Module) setMainProductImageBySlugDirect(c *fiber.Ctx) error {
	if err := m.setStorefrontIDBySlug(c); err != nil {
		return err
	}

	// Проверяем доступ после установки storefrontID
	if err := m.checkStorefrontAccess(c); err != nil {
		return err
	}

	return m.imageHandler.SetMainProductImage(c)
}

// storageAdapter adapts postgres.Database to storefrontService.Storage interface
type storageAdapter struct {
	db *postgres.Database
}

// sqlxTransactionWrapper wraps sqlx.Tx to implement Transaction interface
type sqlxTransactionWrapper struct {
	tx *sqlx.Tx
}

func (t *sqlxTransactionWrapper) Rollback() error {
	return t.tx.Rollback()
}

func (t *sqlxTransactionWrapper) Commit() error {
	return t.tx.Commit()
}

func (t *sqlxTransactionWrapper) GetPgxTx() interface{} {
	return t.tx
}

func (t *sqlxTransactionWrapper) GetSqlxTx() *sqlx.Tx {
	return t.tx
}

// GetStorefrontProducts delegates to database
func (s *storageAdapter) GetStorefrontProducts(ctx context.Context, filter models.ProductFilter) ([]*models.StorefrontProduct, error) {
	return s.db.GetStorefrontProducts(ctx, filter)
}

// GetStorefrontProduct delegates to database
func (s *storageAdapter) GetStorefrontProduct(ctx context.Context, storefrontID, productID int) (*models.StorefrontProduct, error) {
	return s.db.GetStorefrontProduct(ctx, storefrontID, productID)
}

// GetStorefrontProductByID delegates to database
func (s *storageAdapter) GetStorefrontProductByID(ctx context.Context, productID int) (*models.StorefrontProduct, error) {
	return s.db.GetStorefrontProductByID(ctx, productID)
}

// GetStorefrontProductBySKU delegates to database
func (s *storageAdapter) GetStorefrontProductBySKU(ctx context.Context, storefrontID int, sku string) (*models.StorefrontProduct, error) {
	return s.db.GetStorefrontProductBySKU(ctx, storefrontID, sku)
}

// GetStorefrontProductsBySKUs delegates to database
func (s *storageAdapter) GetStorefrontProductsBySKUs(ctx context.Context, storefrontID int, skus []string) (map[string]*models.StorefrontProduct, error) {
	return s.db.GetStorefrontProductsBySKUs(ctx, storefrontID, skus)
}

// CreateStorefrontProduct delegates to database
func (s *storageAdapter) CreateStorefrontProduct(ctx context.Context, storefrontID int, req *models.CreateProductRequest) (*models.StorefrontProduct, error) {
	return s.db.CreateStorefrontProduct(ctx, storefrontID, req)
}

// BatchCreateStorefrontProducts delegates to database
func (s *storageAdapter) BatchCreateStorefrontProducts(ctx context.Context, storefrontID int, requests []*models.CreateProductRequest) ([]*models.StorefrontProduct, error) {
	return s.db.BatchCreateStorefrontProducts(ctx, storefrontID, requests)
}

// UpdateStorefrontProduct delegates to database
func (s *storageAdapter) UpdateStorefrontProduct(ctx context.Context, storefrontID, productID int, req *models.UpdateProductRequest) error {
	return s.db.UpdateStorefrontProduct(ctx, storefrontID, productID, req)
}

// DeleteStorefrontProduct delegates to database
func (s *storageAdapter) DeleteStorefrontProduct(ctx context.Context, storefrontID, productID int) error {
	return s.db.DeleteStorefrontProduct(ctx, storefrontID, productID)
}

// HardDeleteStorefrontProduct delegates to database
func (s *storageAdapter) HardDeleteStorefrontProduct(ctx context.Context, storefrontID, productID int) error {
	return s.db.HardDeleteStorefrontProduct(ctx, storefrontID, productID)
}

// UpdateProductInventory delegates to database
func (s *storageAdapter) UpdateProductInventory(ctx context.Context, storefrontID, productID int, userID int, req *models.UpdateInventoryRequest) error {
	return s.db.UpdateProductInventory(ctx, storefrontID, productID, userID, req)
}

// GetProductStats delegates to database
func (s *storageAdapter) GetProductStats(ctx context.Context, storefrontID int) (*models.ProductStats, error) {
	return s.db.GetProductStats(ctx, storefrontID)
}

// IncrementProductViews delegates to database
func (s *storageAdapter) IncrementProductViews(ctx context.Context, productID int) error {
	return s.db.IncrementProductViews(ctx, productID)
}

// BulkCreateProducts delegates to database
func (s *storageAdapter) BulkCreateProducts(ctx context.Context, storefrontID int, products []models.CreateProductRequest) ([]int, []error) {
	return s.db.BulkCreateProducts(ctx, storefrontID, products)
}

// BulkUpdateProducts delegates to database
func (s *storageAdapter) BulkUpdateProducts(ctx context.Context, storefrontID int, updates []models.BulkUpdateItem) ([]int, []error) {
	return s.db.BulkUpdateProducts(ctx, storefrontID, updates)
}

// BulkDeleteProducts delegates to database
func (s *storageAdapter) BulkDeleteProducts(ctx context.Context, storefrontID int, productIDs []int) ([]int, []error) {
	return s.db.BulkDeleteProducts(ctx, storefrontID, productIDs)
}

// BulkUpdateStatus delegates to database
func (s *storageAdapter) BulkUpdateStatus(ctx context.Context, storefrontID int, productIDs []int, isActive bool) ([]int, []error) {
	return s.db.BulkUpdateStatus(ctx, storefrontID, productIDs, isActive)
}

// GetStorefrontByID delegates to database
func (s *storageAdapter) GetStorefrontByID(ctx context.Context, id int) (*models.Storefront, error) {
	storefrontRepo := s.db.Storefront().(postgres.StorefrontRepository)
	return storefrontRepo.GetByID(ctx, id)
}

// SaveTranslation delegates to database
func (s *storageAdapter) SaveTranslation(ctx context.Context, translation *models.Translation) error {
	return s.db.SaveTranslation(ctx, translation)
}

// BeginTx starts a new transaction
func (s *storageAdapter) BeginTx(ctx context.Context) (storefrontService.Transaction, error) {
	// Get sqlx.DB from database
	sqlxDB := s.db.GetSQLXDB()

	// Start transaction
	tx, err := sqlxDB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Wrap the transaction
	return &sqlxTransactionWrapper{tx: tx}, nil
}

// CreateStorefrontProductTx creates a product within a transaction
func (s *storageAdapter) CreateStorefrontProductTx(ctx context.Context, tx storefrontService.Transaction, storefrontID int, req *models.CreateProductRequest) (*models.StorefrontProduct, error) {
	// Get the underlying sqlx.Tx from interface
	sqlxWrapper, ok := tx.(*sqlxTransactionWrapper)
	if !ok {
		return nil, fmt.Errorf("invalid transaction type")
	}

	// Call the database method with the transaction
	return s.db.CreateStorefrontProductTx(ctx, sqlxWrapper.tx, storefrontID, req)
}
