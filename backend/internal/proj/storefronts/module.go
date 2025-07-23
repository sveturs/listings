package storefronts

import (
	"backend/internal/logger"
	"backend/internal/middleware"
	"backend/internal/proj/global/service"
	"backend/internal/proj/storefronts/handler"
	storefrontService "backend/internal/proj/storefronts/service"
	"backend/internal/storage/postgres"

	"github.com/gofiber/fiber/v2"
)

// Module представляет модуль витрин с продуктами
type Module struct {
	services          service.ServicesInterface
	storefrontHandler *handler.StorefrontHandler
	productHandler    *handler.ProductHandler
	importHandler     *handler.ImportHandler
	imageHandler      *handler.ImageHandler
}

// NewModule создает новый модуль витрин
func NewModule(services service.ServicesInterface) *Module {
	storefrontSvc := services.Storefront()
	// ProductService использует тот же storage, что и все остальные сервисы
	// Предполагаем, что Storage реализует нужные методы
	productStorage, ok := services.Storage().(storefrontService.Storage)
	if !ok {
		// Если storage не реализует нужный интерфейс, создаем заглушку
		panic("Storage does not implement product storage interface")
	}
	// Получаем OpenSearch репозиторий для товаров витрин
	var searchRepo storefrontService.ProductSearchRepository
	if osProductRepo := services.Storage().StorefrontProductSearch(); osProductRepo != nil {
		searchRepo = osProductRepo.(storefrontService.ProductSearchRepository)
	}
	productSvc := storefrontService.NewProductService(productStorage, searchRepo)

	// Создаем сервис импорта
	importSvc := storefrontService.NewImportService(productSvc)

	// Создаем единый ImageService
	db := services.Storage().(*postgres.Database)
	imageRepo := postgres.NewImageRepository(db.GetSQLXDB())
	imageService := services.NewImageService(services.FileStorage(), imageRepo)

	return &Module{
		services:          services,
		storefrontHandler: handler.NewStorefrontHandler(storefrontSvc),
		productHandler:    handler.NewProductHandler(productSvc),
		importHandler:     handler.NewImportHandler(importSvc),
		imageHandler:      handler.NewImageHandler(imageService, productSvc),
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

		if err == postgres.ErrNotFound {
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

// setStorefrontIDBySlug добавляет storefront ID в контекст по slug
func (m *Module) setStorefrontIDBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Storefront slug is required",
		})
	}

	// Получаем витрину через сервис по slug
	storefront, err := m.services.Storefront().GetBySlug(c.Context(), slug)
	if err != nil {
		if err == postgres.ErrNotFound {
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
