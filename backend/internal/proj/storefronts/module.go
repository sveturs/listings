package storefronts

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"backend/internal/domain/models"
	"backend/internal/logger"
	"backend/internal/middleware"
	"backend/internal/proj/global/service"
	storefrontHandler "backend/internal/proj/storefront/handler"
	"backend/internal/proj/storefront/repository"
	"backend/internal/proj/storefronts/handler"
	storefrontService "backend/internal/proj/storefronts/service"
	"backend/internal/storage/postgres"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
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
}

// NewModule создает новый модуль витрин
func NewModule(services service.ServicesInterface) *Module {
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

	// Создаем сервис импорта
	importSvc := storefrontService.NewImportService(productSvc)

	// Создаем единый ImageService
	imageRepo := postgres.NewImageRepository(db.GetSQLXDB())
	imageService := services.NewImageService(services.FileStorage(), imageRepo)

	// Получаем storefront repository
	storefrontRepo := services.Storage().Storefront().(postgres.StorefrontRepository)

	// Создаем variant handlers (используем уже созданный variantRepo)
	variantHandler := storefrontHandler.NewVariantHandler(variantRepo)
	publicVariantHandler := storefrontHandler.NewPublicVariantHandler(variantRepo)

	return &Module{
		services:             services,
		storefrontHandler:    handler.NewStorefrontHandler(storefrontSvc),
		productHandler:       handler.NewProductHandler(productSvc),
		importHandler:        handler.NewImportHandler(importSvc),
		imageHandler:         handler.NewImageHandler(imageService, productSvc),
		dashboardHandler:     handler.NewDashboardHandler(storefrontSvc, productSvc, storefrontRepo),
		variantHandler:       variantHandler,
		publicVariantHandler: publicVariantHandler,
	}
}

// GetPrefix возвращает префикс для маршрутов
func (m *Module) GetPrefix() string {
	return "/api/v1/storefronts"
}

// RegisterRoutes регистрирует все маршруты модуля
func (m *Module) RegisterRoutes(app *fiber.App, mw *middleware.Middleware) error {
	api := app.Group("/api/v1")

	// Регистрируем защищенный маршрут /my первым, чтобы он имел приоритет
	api.Get("/storefronts/my", mw.AuthRequiredJWT, m.storefrontHandler.GetMyStorefronts)

	// Публичный endpoint для получения товара по ID (для чата) - ВАЖНО: должен быть ПЕРЕД группой /storefronts
	api.Get("/storefronts/products/:id", m.productHandler.GetProductByID)

	// Публичные маршруты витрин (без авторизации)
	public := api.Group("/storefronts")
	{
		// Конкретные маршруты (должны быть перед параметризованными)
		// Списки и поиск
		public.Get("/", m.storefrontHandler.ListStorefronts)
		public.Get("/search", m.storefrontHandler.SearchOpenSearch)
		public.Get("/nearby", m.storefrontHandler.GetNearbyStorefronts)

		// Картографические данные
		public.Get("/map", m.storefrontHandler.GetMapData)
		public.Get("/building", m.storefrontHandler.GetBusinessesInBuilding)

		// Маршруты с slug
		public.Get("/slug/:slug", m.storefrontHandler.GetStorefrontBySlug)

		// Публичные маршруты товаров с использованием slug
		public.Get("/slug/:slug/products", m.getProductsBySlug)
		public.Get("/slug/:slug/products/:id", m.getProductBySlug)

		// Публичные маршруты для вариантов товаров (БЕЗ АВТОРИЗАЦИИ)
		public.Get("/slug/:slug/products/:product_id/variants", m.publicVariantHandler.GetProductVariantsPublic)
		public.Get("/:slug/products/:product_id", m.publicVariantHandler.GetProductPublic)
		public.Get("/variants/:variant_id", m.publicVariantHandler.GetVariantByIDPublic)

		// Параметризованные маршруты (должны быть последними)
		// Получение витрины
		public.Get("/:id", m.storefrontHandler.GetStorefront)

		// Персонал (просмотр)
		public.Get("/:id/staff", m.storefrontHandler.GetStaff)

		// Аналитика (запись просмотра)
		public.Post("/:id/view", m.storefrontHandler.RecordView)
	}

	// Защищенные маршруты витрин (требуют авторизации)
	protected := api.Group("/storefronts")
	protected.Use(mw.AuthRequiredJWT)
	{
		// Управление витринами
		protected.Post("/", m.storefrontHandler.CreateStorefront)
		protected.Put("/:id", m.storefrontHandler.UpdateStorefront)
		protected.Delete("/:id", m.storefrontHandler.DeleteStorefront)

		// Настройки витрины
		protected.Put("/:id/hours", m.storefrontHandler.UpdateWorkingHours)
		protected.Put("/:id/payment-methods", m.storefrontHandler.UpdatePaymentMethods)
		protected.Put("/:id/delivery-options", m.storefrontHandler.UpdateDeliveryOptions)

		// Управление персоналом
		protected.Post("/:id/staff", m.storefrontHandler.AddStaff)
		protected.Put("/:id/staff/:staffId/permissions", m.storefrontHandler.UpdateStaffPermissions)
		protected.Delete("/:id/staff/:userId", m.storefrontHandler.RemoveStaff)

		// Загрузка изображений
		protected.Post("/:id/logo", m.storefrontHandler.UploadLogo)
		protected.Post("/:id/banner", m.storefrontHandler.UploadBanner)

		// Аналитика (просмотр)
		protected.Get("/:id/analytics", m.storefrontHandler.GetAnalytics)

		// Защищенные маршруты товаров с использованием slug
		protected.Post("/slug/:slug/products", m.createProductBySlug)
		protected.Put("/slug/:slug/products/:id", m.updateProductBySlug)
		protected.Delete("/slug/:slug/products/:id", m.deleteProductBySlug)
		protected.Post("/slug/:slug/products/:id/inventory", m.updateInventoryBySlug)
		protected.Get("/slug/:slug/products/stats", m.getProductStatsBySlug)

		// Маршруты для изображений товаров
		protected.Post("/slug/:slug/products/:product_id/images", m.uploadProductImageBySlug)
		protected.Get("/slug/:slug/products/:product_id/images", m.getProductImagesBySlug)
		protected.Delete("/slug/:slug/products/:product_id/images/:image_id", m.deleteProductImageBySlug)
		protected.Post("/slug/:slug/products/:product_id/images/:image_id/main", m.setMainProductImageBySlug)
		protected.Put("/slug/:slug/products/:product_id/images/order", m.updateImageOrderBySlug)

		// Bulk операции с товарами
		protected.Post("/slug/:slug/products/bulk/create", m.bulkCreateProductsBySlug)
		protected.Put("/slug/:slug/products/bulk/update", m.bulkUpdateProductsBySlug)
		protected.Delete("/slug/:slug/products/bulk/delete", m.bulkDeleteProductsBySlug)
		protected.Put("/slug/:slug/products/bulk/status", m.bulkUpdateStatusBySlug)

		// Маршруты импорта товаров
		protected.Post("/:id/import/url", m.importHandler.ImportFromURL)
		protected.Post("/:id/import/file", m.importHandler.ImportFromFile)
		protected.Post("/:id/import/validate", m.importHandler.ValidateImportFile)
		protected.Get("/:id/import/jobs", m.importHandler.GetJobs)
		protected.Get("/:id/import/jobs/:jobId", m.importHandler.GetJobDetails)
		protected.Get("/:id/import/jobs/:jobId/status", m.importHandler.GetJobStatus)
		protected.Post("/:id/import/jobs/:jobId/cancel", m.importHandler.CancelJob)
		protected.Post("/:id/import/jobs/:jobId/retry", m.importHandler.RetryJob)

		// Маршруты импорта через slug
		protected.Post("/slug/:slug/import/url", m.importFromURLBySlug)
		protected.Post("/slug/:slug/import/file", m.importFromFileBySlug)
		protected.Post("/slug/:slug/import/validate", m.validateImportBySlug)
		protected.Get("/slug/:slug/import/jobs", m.getJobsBySlug)
		protected.Get("/slug/:slug/import/jobs/:jobId", m.getJobDetailsBySlug)
		protected.Get("/slug/:slug/import/jobs/:jobId/status", m.getJobStatusBySlug)
		protected.Post("/slug/:slug/import/jobs/:jobId/cancel", m.cancelJobBySlug)
		protected.Post("/slug/:slug/import/jobs/:jobId/retry", m.retryJobBySlug)

		// Маршруты товаров по slug без префикса (для frontend совместимости)
		protected.Put("/:slug/products/:id", m.updateProductBySlugDirect)
		protected.Post("/:slug/products/:product_id/images", m.uploadProductImageBySlugDirect)
		protected.Delete("/:slug/products/:product_id/images/:image_id", m.deleteProductImageBySlugDirect)
		protected.Post("/:slug/products/:product_id/images/:image_id/main", m.setMainProductImageBySlugDirect)

		// Dashboard маршруты
		protected.Get("/:slug/dashboard/stats", m.dashboardHandler.GetDashboardStats)
		protected.Get("/:slug/dashboard/recent-orders", m.dashboardHandler.GetRecentOrders)
		protected.Get("/:slug/dashboard/low-stock", m.dashboardHandler.GetLowStockProducts)
		protected.Get("/:slug/dashboard/notifications", m.dashboardHandler.GetDashboardNotifications)

		// Приватные маршруты для управления вариантами товаров (требуют авторизации)
		variants := protected.Group("/storefront")
		{
			// Глобальные атрибуты вариантов
			variants.Get("/variants/attributes", m.variantHandler.GetVariantAttributes)
			variants.Get("/variants/attributes/:attribute_id/values", m.variantHandler.GetVariantAttributeValues)

			// Управление вариантами товара
			variants.Get("/products/:product_id/variants", m.variantHandler.GetVariantsByProductID)
			variants.Post("/variants", m.variantHandler.CreateVariant)
			variants.Put("/variants/:variant_id", m.variantHandler.UpdateVariant)
			variants.Delete("/variants/:variant_id", m.variantHandler.DeleteVariant)
			variants.Get("/variants/:variant_id", m.variantHandler.GetVariantByID)

			// Генерация и массовые операции
			variants.Post("/variants/generate", m.variantHandler.GenerateVariants)
			variants.Post("/variants/bulk", m.variantHandler.BulkCreateVariants)
			variants.Get("/products/:product_id/variant-matrix", m.variantHandler.GetVariantMatrix)
			variants.Post("/products/:product_id/variants/bulk-update-stock", m.variantHandler.BulkUpdateStock)
			variants.Get("/products/:product_id/variants/analytics", m.variantHandler.GetVariantAnalytics)

			// CSV импорт/экспорт
			variants.Post("/products/:product_id/variants/import", m.variantHandler.ImportVariants)
			variants.Get("/products/:product_id/variants/export", m.variantHandler.ExportVariants)

			// Настройка атрибутов товара
			variants.Post("/products/attributes/setup", m.variantHandler.SetupProductAttributes)
			variants.Get("/products/:product_id/attributes", m.variantHandler.GetProductAttributes)
			variants.Get("/categories/:category_id/attributes", m.variantHandler.GetAvailableAttributesForCategory)
		}
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
	userID, ok := c.Locals("user_id").(int)
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

	var storefront *models.Storefront
	var err error

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

	userID, ok := c.Locals("user_id").(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User ID not found in context",
		})
	}

	// Проверяем, является ли пользователь администратором
	isAdmin, _ := c.Locals("is_admin").(bool)
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

// CreateStorefrontProduct delegates to database
func (s *storageAdapter) CreateStorefrontProduct(ctx context.Context, storefrontID int, req *models.CreateProductRequest) (*models.StorefrontProduct, error) {
	return s.db.CreateStorefrontProduct(ctx, storefrontID, req)
}

// UpdateStorefrontProduct delegates to database
func (s *storageAdapter) UpdateStorefrontProduct(ctx context.Context, storefrontID, productID int, req *models.UpdateProductRequest) error {
	return s.db.UpdateStorefrontProduct(ctx, storefrontID, productID, req)
}

// DeleteStorefrontProduct delegates to database
func (s *storageAdapter) DeleteStorefrontProduct(ctx context.Context, storefrontID, productID int) error {
	return s.db.DeleteStorefrontProduct(ctx, storefrontID, productID)
}

// UpdateProductInventory delegates to database
func (s *storageAdapter) UpdateProductInventory(ctx context.Context, storefrontID, productID int, userID int, req *models.UpdateInventoryRequest) error {
	return s.db.UpdateProductInventory(ctx, storefrontID, productID, userID, req)
}

// GetProductStats delegates to database
func (s *storageAdapter) GetProductStats(ctx context.Context, storefrontID int) (*models.ProductStats, error) {
	return s.db.GetProductStats(ctx, storefrontID)
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
