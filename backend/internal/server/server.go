// backend/internal/server/server.go
package server

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/gofiber/websocket/v2"
	"github.com/pkg/errors"

	_ "backend/docs"
	"backend/internal/config"
	"backend/internal/logger"
	"backend/internal/middleware"
	balanceHandler "backend/internal/proj/balance/handler"
	geocodeHandler "backend/internal/proj/geocode/handler"
	globalService "backend/internal/proj/global/service"
	marketplaceHandler "backend/internal/proj/marketplace/handler"
	marketplaceService "backend/internal/proj/marketplace/service"
	notificationHandler "backend/internal/proj/notifications/handler"
	paymentService "backend/internal/proj/payments/service"
	reviewHandler "backend/internal/proj/reviews/handler"
	storefrontHandler "backend/internal/proj/storefront/handler"
	userHandler "backend/internal/proj/users/handler"
	"backend/internal/storage/filestorage"
	"backend/internal/storage/opensearch"
	"backend/internal/storage/postgres"
)

type Server struct {
	app           *fiber.App
	cfg           *config.Config
	users         *userHandler.Handler
	middleware    *middleware.Middleware
	review        *reviewHandler.Handler
	marketplace   *marketplaceHandler.Handler
	notifications *notificationHandler.Handler
	balance       *balanceHandler.Handler
	payments      paymentService.PaymentServiceInterface
	storefront    *storefrontHandler.Handler
	geocode       *geocodeHandler.GeocodeHandler
	contacts      *marketplaceHandler.ContactsHandler
	fileStorage   filestorage.FileStorageInterface
}

func NewServer(cfg *config.Config) (*Server, error) {
	fileStorage, err := filestorage.NewFileStorage(cfg.FileStorage)
	if err != nil {
		return nil, errors.Wrap(err, "Ошибка инициализации файлового хранилища")
	}

	osClient, err := initializeOpenSearch(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "OpenSearch initialization failed")
	} else {
		logger.Info().Msg("Успешное подключение к OpenSearch")
	}
	db, err := postgres.NewDatabase(cfg.DatabaseURL, osClient, cfg.OpenSearch.MarketplaceIndex, fileStorage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize database")
	}

	translationService, err := initializeTranslationService(cfg, db)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize translation service: %w", err)
	}

	services := globalService.NewService(db, cfg, translationService)

	usersHandler := userHandler.NewHandler(services)
	reviewHandler := reviewHandler.NewHandler(services)
	notificationsHandler := notificationHandler.NewHandler(services.Notification())
	marketplaceHandlerInstance := marketplaceHandler.NewHandler(services)
	balanceHandler := balanceHandler.NewHandler(services)
	storefrontHandler := storefrontHandler.NewHandler(services)
	contactsHandler := marketplaceHandler.NewContactsHandler(services)
	middleware := middleware.NewMiddleware(cfg, services)
	geocodeHandler := geocodeHandler.NewGeocodeHandler(services.Geocode())

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Детальное логирование ошибки
			logger.Error().
				Err(err).
				Str("error_type", fmt.Sprintf("%T", err)).
				Str("path", c.Path()).
				Msg("Error in handler")

			// Стандартная обработка ошибки
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			return c.Status(code).JSON(fiber.Map{
				"error":  err.Error(),
				"status": code,
				"path":   c.Path(),
			})
		},
		BodyLimit:               50 * 1024 * 1024,
		EnableTrustedProxyCheck: true,
		TrustedProxies:          []string{"127.0.0.1", "::1"},
	})

	server := &Server{
		app:           app,
		cfg:           cfg,
		users:         usersHandler,
		middleware:    middleware,
		review:        reviewHandler,
		marketplace:   marketplaceHandlerInstance,
		notifications: notificationsHandler,
		balance:       balanceHandler,
		storefront:    storefrontHandler,
		payments:      services.Payment(),
		geocode:       geocodeHandler,
		contacts:      contactsHandler,
		fileStorage:   fileStorage,
	}

	notificationsHandler.ConnectTelegramWebhook()
	server.setupMiddleware()
	server.setupRoutes()

	return server, nil
}

func initializeTranslationService(cfg *config.Config, db *postgres.Database) (marketplaceService.TranslationServiceInterface, error) {
	if cfg.GoogleTranslateAPIKey != "" && cfg.OpenAIAPIKey != "" {
		translationFactory, err := marketplaceService.NewTranslationServiceFactory(cfg.GoogleTranslateAPIKey, cfg.OpenAIAPIKey, db)
		if err == nil {
			logger.Info().Msg("Создана фабрика сервисов перевода с поддержкой Google Translate и OpenAI")
			return translationFactory, nil
		} else {
			logger.Error().Err(err).Msg("Ошибка создания фабрики перевода, будет использован только OpenAI")
		}
	}

	if cfg.OpenAIAPIKey != "" {
		translationService, err := marketplaceService.NewTranslationService(cfg.OpenAIAPIKey)
		if err != nil {
			return nil, err
		}
		logger.Info().Msg("Создан сервис перевода на базе OpenAI")
		return translationService, nil
	}

	return nil, fmt.Errorf("не указан ни один API ключ для перевода")
}

func initializeOpenSearch(cfg *config.Config) (*opensearch.OpenSearchClient, error) {
	if cfg.OpenSearch.URL == "" {
		return nil, errors.New("OpenSearch URL не указан, поиск будет отключен")
	}

	osClient, err := opensearch.NewOpenSearchClient(opensearch.Config{
		URL:      cfg.OpenSearch.URL,
		Username: cfg.OpenSearch.Username,
		Password: cfg.OpenSearch.Password,
	})

	if err != nil {
		return nil, errors.New("Ошибка подключения к OpenSearch")
	}

	return osClient, nil
}

func (s *Server) setupMiddleware() {
	// Общие middleware для observability
	// Security headers должны быть первыми
	s.app.Use(s.middleware.SecurityHeaders())
	s.app.Use(s.middleware.CORS())
	s.app.Use(s.middleware.Logger())
	// TODO: Добавить middleware для метрик

	s.app.Use("/ws", s.middleware.AuthRequiredJWT)
}

func (s *Server) setupRoutes() {

	s.app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Svetu API")
	})

	s.app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Swagger документация
	s.app.Get("/swagger/*", swagger.HandlerDefault)
	s.app.Get("/docs/*", swagger.New(swagger.Config{
		URL:         "/swagger/doc.json",
		DeepLinking: false,
	}))

	// WebSocket с проверкой аутентификации и rate limiting
	s.app.Get("/ws/chat", s.middleware.AuthRequiredJWT, s.middleware.RateLimitByUser(5, time.Minute), func(c *fiber.Ctx) error {
		// Проверяем, что это WebSocket запрос
		if websocket.IsWebSocketUpgrade(c) {
			// Сохраняем userID для использования в WebSocket handler
			userID := c.Locals("user_id").(int)

			return websocket.New(func(conn *websocket.Conn) {
				// Передаем userID через контекст соединения
				// В Fiber WebSocket, Locals доступен только для чтения
				// Поэтому создаем обертку с сохраненным userID
				s.marketplace.Chat.HandleWebSocketWithAuth(conn, userID)
			}, websocket.Config{
				HandshakeTimeout:  10 * time.Second,
				ReadBufferSize:    1024,
				WriteBufferSize:   1024,
				EnableCompression: false,
			})(c)
		}
		return fiber.ErrUpgradeRequired
	})

	// Регистрируем роуты через новую систему
	s.registerProjectRoutes()

	// Убираем старые роуты notifications, они теперь в registerProjectRoutes
	// s.app.Post("/api/v1/notifications/telegram/webhook", s.notifications.Notification.HandleTelegramWebhook)
	// s.app.Get("/api/v1/notifications/telegram", s.notifications.Notification.GetTelegramStatus)

	// s.app.Post("/api/v1/public/send-email", s.notifications.Notification.SendPublicEmail)
	s.app.Get("/api/v1/public/storefronts/:id", s.storefront.Storefront.GetPublicStorefront)

	balanceRoutes := s.app.Group("/api/v1/balance", s.middleware.AuthRequiredJWT)
	balanceRoutes.Get("/", s.balance.Balance.GetBalance)
	balanceRoutes.Get("/transactions", s.balance.Balance.GetTransactions)
	balanceRoutes.Get("/payment-methods", s.balance.Balance.GetPaymentMethods)
	balanceRoutes.Post("/deposit", s.balance.Balance.CreateDeposit)

	s.app.Post("/webhook/stripe", func(c *fiber.Ctx) error {
		payload := c.Body()
		signature := c.Get("Stripe-Signature")

		err := s.payments.HandleWebhook(c.Context(), payload, signature)
		if err != nil {
			log.Printf("Webhook error: %v", err)
			return c.SendStatus(fiber.StatusBadRequest)
		}

		return c.SendStatus(fiber.StatusOK)
	})

	// Обновлено: маршруты публичного API маркетплейса используют соответствующие обработчики
	marketplace := s.app.Group("/api/v1/marketplace")
	marketplace.Get("/listings", s.marketplace.Listings.GetListings)
	marketplace.Get("/categories", s.marketplace.Categories.GetCategories)
	marketplace.Get("/category-tree", s.marketplace.Categories.GetCategoryTree)
	marketplace.Get("/listings/:id", s.marketplace.Listings.GetListing)
	marketplace.Get("/search", s.marketplace.Search.SearchListingsAdvanced) // маршрут поиска
	marketplace.Get("/suggestions", s.marketplace.Search.GetSuggestions)    // маршрут автодополнения
	marketplace.Get("/category-suggestions", s.marketplace.Search.GetCategorySuggestions)

	// Временный публичный маршрут для проверки
	s.app.Get("/admin-categories-test", s.marketplace.AdminCategories.GetCategories)
	marketplace.Get("/categories/:id/attributes", s.marketplace.Categories.GetCategoryAttributes)
	marketplace.Get("/listings/:id/price-history", s.marketplace.Listings.GetPriceHistory)
	marketplace.Get("/listings/:id/similar", s.marketplace.Search.GetSimilarListings)
	marketplace.Get("/categories/:id/attribute-ranges", s.marketplace.Categories.GetAttributeRanges)
	marketplace.Get("/enhanced-suggestions", s.marketplace.Search.GetEnhancedSuggestions)

	// Карта - геопространственные маршруты
	marketplace.Get("/map/bounds", s.marketplace.GetListingsInBounds)
	marketplace.Get("/map/clusters", s.marketplace.GetMapClusters)

	// Обновлено: маршруты API переводов используют обработчик переводов
	translation := s.app.Group("/api/v1/translation")
	translation.Get("/limits", s.marketplace.Translations.GetTranslationLimits)
	translation.Post("/provider", s.marketplace.Translations.SetTranslationProvider)

	// CSRF токен
	s.app.Get("/api/v1/csrf-token", s.middleware.GetCSRFToken())

	// Применяем rate limiting для authentication endpoints
	// New JWT endpoints with consistent /api/auth prefix

	// API routes с общим rate limiting по пользователю (300 запросов в минуту)
	// Используем JWT как основной метод аутентификации
	api := s.app.Group("/api/v1", s.middleware.AuthRequiredJWT, s.middleware.RateLimitByUser(300, time.Minute))

	authedAPIGroup := s.app.Group("/api/v1", s.middleware.AuthRequiredJWT, s.middleware.CSRFProtection())

	storefronts := authedAPIGroup.Group("/storefronts")
	storefronts.Get("/", s.storefront.Storefront.GetUserStorefronts)
	storefronts.Post("/", s.storefront.Storefront.CreateStorefront)
	storefronts.Get("/:id", s.storefront.Storefront.GetStorefront)
	storefronts.Put("/:id", s.storefront.Storefront.UpdateStorefront)
	storefronts.Delete("/:id", s.storefront.Storefront.DeleteStorefront)
	storefronts.Get("/:id/import-sources", s.storefront.Storefront.GetImportSources)
	storefronts.Post("/import-sources", s.storefront.Storefront.CreateImportSource)
	storefronts.Put("/import-sources/:id", s.storefront.Storefront.UpdateImportSource)
	storefronts.Delete("/import-sources/:id", s.storefront.Storefront.DeleteImportSource)
	storefronts.Post("/import-sources/:id/run", s.storefront.Storefront.RunImport)
	storefronts.Get("/import-sources/:id/history", s.storefront.Storefront.GetImportHistory)
	storefronts.Get("/import-sources/:id/category-mappings", s.storefront.Storefront.GetCategoryMappings)
	storefronts.Put("/import-sources/:id/category-mappings", s.storefront.Storefront.UpdateCategoryMappings)
	storefronts.Get("/import-sources/:id/imported-categories", s.storefront.Storefront.GetImportedCategories)
	storefronts.Post("/import-sources/:id/apply-category-mappings", s.storefront.Storefront.ApplyCategoryMappings)

	geocodeApi := s.app.Group("/api/v1/geocode")
	geocodeApi.Get("/reverse", s.geocode.ReverseGeocode)

	citiesApi := s.app.Group("/api/v1/cities")
	citiesApi.Get("/suggest", s.geocode.GetCitySuggestions)

	// Contacts routes
	contacts := api.Group("/contacts")
	contacts.Get("/", s.contacts.GetContacts)
	contacts.Post("/", s.contacts.AddContact)
	contacts.Put("/:contact_user_id", s.contacts.UpdateContactStatus)
	contacts.Delete("/:contact_user_id", s.contacts.RemoveContact)
	contacts.Get("/privacy", s.contacts.GetPrivacySettings)
	contacts.Put("/privacy", s.contacts.UpdatePrivacySettings)
	contacts.Get("/status/:contact_user_id", s.contacts.GetContactStatus)

	// Обновлено: маршруты защищенного API маркетплейса используют соответствующие обработчики
	marketplaceProtected := authedAPIGroup.Group("/marketplace")
	marketplaceProtected.Post("/listings", s.marketplace.Listings.CreateListing)
	marketplaceProtected.Put("/listings/:id", s.marketplace.Listings.UpdateListing)
	marketplaceProtected.Delete("/listings/:id", s.marketplace.Listings.DeleteListing)
	marketplaceProtected.Post("/listings/:id/images", s.marketplace.Images.UploadImages)
	marketplaceProtected.Post("/listings/:id/favorite", s.marketplace.Favorites.AddToFavorites)
	marketplaceProtected.Delete("/listings/:id/favorite", s.marketplace.Favorites.RemoveFromFavorites)
	marketplaceProtected.Get("/favorites", s.marketplace.Favorites.GetFavorites)
	marketplaceProtected.Put("/translations/:id", s.marketplace.Translations.UpdateTranslations)
	marketplaceProtected.Post("/translations/batch", s.marketplace.Translations.TranslateText) // Предполагается, что этот метод переименован
	marketplaceProtected.Post("/moderate-image", s.marketplace.Images.ModerateImage)
	marketplaceProtected.Post("/enhance-preview", s.marketplace.Images.EnhancePreview)
	marketplaceProtected.Post("/enhance-images", s.marketplace.Images.EnhanceImages)

	// маршруты для новых методов в TranslationsHandler
	marketplaceProtected.Post("/translations/batch-translate", s.marketplace.Translations.BatchTranslateListings)
	marketplaceProtected.Post("/translations/translate", s.marketplace.Translations.TranslateText)
	marketplaceProtected.Post("/translations/detect-language", s.marketplace.Translations.DetectLanguage)
	marketplaceProtected.Get("/translations/:id", s.marketplace.Translations.GetTranslations)
	// Административные маршруты
	adminRoutes := s.app.Group("/api/v1/admin", s.middleware.AuthRequiredJWT, s.middleware.AdminRequired, s.middleware.CSRFProtection())

	// Регистрируем маршруты администрирования категорий
	adminRoutes.Post("/categories", s.marketplace.AdminCategories.CreateCategory)
	adminRoutes.Get("/categories", s.marketplace.AdminCategories.GetCategories)
	adminRoutes.Get("/categories/:id", s.marketplace.AdminCategories.GetCategoryByID)
	adminRoutes.Put("/categories/:id", s.marketplace.AdminCategories.UpdateCategory)
	adminRoutes.Delete("/categories/:id", s.marketplace.AdminCategories.DeleteCategory)
	adminRoutes.Post("/categories/:id/reorder", s.marketplace.AdminCategories.ReorderCategories)
	adminRoutes.Put("/categories/:id/move", s.marketplace.AdminCategories.MoveCategory)
	adminRoutes.Post("/categories/:id/attributes", s.marketplace.AdminCategories.AddAttributeToCategory)
	adminRoutes.Delete("/categories/:id/attributes/:attr_id", s.marketplace.AdminCategories.RemoveAttributeFromCategory)
	adminRoutes.Put("/categories/:id/attributes/:attr_id", s.marketplace.AdminCategories.UpdateAttributeCategory)

	// Регистрируем маршруты администрирования атрибутов
	adminRoutes.Post("/attributes", s.marketplace.AdminAttributes.CreateAttribute)
	adminRoutes.Get("/attributes", s.marketplace.AdminAttributes.GetAttributes)
	adminRoutes.Get("/attributes/:id", s.marketplace.AdminAttributes.GetAttributeByID)
	adminRoutes.Put("/attributes/:id", s.marketplace.AdminAttributes.UpdateAttribute)
	adminRoutes.Delete("/attributes/:id", s.marketplace.AdminAttributes.DeleteAttribute)
	adminRoutes.Post("/attributes/bulk-update", s.marketplace.AdminAttributes.BulkUpdateAttributes)

	// Маршруты для экспорта/импорта настроек атрибутов
	adminRoutes.Get("/categories/:categoryId/attributes/export", s.marketplace.AdminAttributes.ExportCategoryAttributes)
	adminRoutes.Post("/categories/:categoryId/attributes/import", s.marketplace.AdminAttributes.ImportCategoryAttributes)
	adminRoutes.Post("/categories/:targetCategoryId/attributes/copy", s.marketplace.AdminAttributes.CopyAttributesSettings)

	// Для обратной совместимости добавим маршруты без v1
	legacyAdmin := s.app.Group("/api/admin")
	legacyAdmin.Use(s.middleware.AuthRequiredJWT)
	legacyAdmin.Use(s.middleware.AdminRequired)
	legacyAdmin.Use(s.middleware.CSRFProtection())

	// Все маршруты для категорий
	legacyAdmin.Get("/categories", s.marketplace.AdminCategories.GetCategories)
	legacyAdmin.Post("/categories", s.marketplace.AdminCategories.CreateCategory)
	legacyAdmin.Get("/categories/:id", s.marketplace.AdminCategories.GetCategoryByID)
	legacyAdmin.Put("/categories/:id", s.marketplace.AdminCategories.UpdateCategory)
	legacyAdmin.Delete("/categories/:id", s.marketplace.AdminCategories.DeleteCategory)
	legacyAdmin.Post("/categories/:id/reorder", s.marketplace.AdminCategories.ReorderCategories)
	legacyAdmin.Put("/categories/:id/move", s.marketplace.AdminCategories.MoveCategory)
	legacyAdmin.Post("/categories/:id/attributes", s.marketplace.AdminCategories.AddAttributeToCategory)
	legacyAdmin.Delete("/categories/:id/attributes/:attr_id", s.marketplace.AdminCategories.RemoveAttributeFromCategory)
	legacyAdmin.Put("/categories/:id/attributes/:attr_id", s.marketplace.AdminCategories.UpdateAttributeCategory)

	// Маршруты для атрибутов
	legacyAdmin.Post("/attributes", s.marketplace.AdminAttributes.CreateAttribute)
	legacyAdmin.Get("/attributes", s.marketplace.AdminAttributes.GetAttributes)
	legacyAdmin.Get("/attributes/:id", s.marketplace.AdminAttributes.GetAttributeByID)
	legacyAdmin.Put("/attributes/:id", s.marketplace.AdminAttributes.UpdateAttribute)
	legacyAdmin.Delete("/attributes/:id", s.marketplace.AdminAttributes.DeleteAttribute)
	legacyAdmin.Post("/attributes/bulk-update", s.marketplace.AdminAttributes.BulkUpdateAttributes)

	// Добавляем маршруты для экспорта/импорта настроек атрибутов
	legacyAdmin.Get("/categories/:categoryId/attributes/export", s.marketplace.AdminAttributes.ExportCategoryAttributes)
	legacyAdmin.Post("/categories/:categoryId/attributes/import", s.marketplace.AdminAttributes.ImportCategoryAttributes)
	legacyAdmin.Post("/categories/:targetCategoryId/attributes/copy", s.marketplace.AdminAttributes.CopyAttributesSettings)

	// Маршруты для групп атрибутов
	legacyAdmin.Get("/attribute-groups", s.marketplace.MarketplaceHandler.ListAttributeGroups)
	legacyAdmin.Post("/attribute-groups", s.marketplace.MarketplaceHandler.CreateAttributeGroup)
	legacyAdmin.Get("/attribute-groups/:id", s.marketplace.MarketplaceHandler.GetAttributeGroup)
	legacyAdmin.Put("/attribute-groups/:id", s.marketplace.MarketplaceHandler.UpdateAttributeGroup)
	legacyAdmin.Delete("/attribute-groups/:id", s.marketplace.MarketplaceHandler.DeleteAttributeGroup)
	legacyAdmin.Get("/attribute-groups/:id/items", s.marketplace.MarketplaceHandler.GetAttributeGroupWithItems)
	legacyAdmin.Post("/attribute-groups/:id/items", s.marketplace.MarketplaceHandler.AddItemToGroup)
	legacyAdmin.Delete("/attribute-groups/:id/items/:attributeId", s.marketplace.MarketplaceHandler.RemoveItemFromGroup)

	// Маршруты для привязки групп к категориям
	legacyAdmin.Get("/categories/:id/attribute-groups", s.marketplace.MarketplaceHandler.GetCategoryGroups)
	legacyAdmin.Post("/categories/:id/attribute-groups", s.marketplace.MarketplaceHandler.AttachGroupToCategory)
	legacyAdmin.Delete("/categories/:id/attribute-groups/:groupId", s.marketplace.MarketplaceHandler.DetachGroupFromCategory)

	// Маршруты для кастомных UI компонентов
	// ВАЖНО: Более специфичные роуты должны идти раньше параметризованных

	// Маршруты для шаблонов (должны быть перед :id, чтобы не конфликтовать)
	adminRoutes.Get("/custom-components/templates", s.marketplace.CustomComponents.ListTemplates)
	adminRoutes.Post("/custom-components/templates", s.marketplace.CustomComponents.CreateTemplate)

	// Маршруты для использования компонентов
	adminRoutes.Get("/custom-components/usage", s.marketplace.CustomComponents.GetComponentUsages)
	adminRoutes.Post("/custom-components/usage", s.marketplace.CustomComponents.AddComponentUsage)
	adminRoutes.Delete("/custom-components/usage/:id", s.marketplace.CustomComponents.RemoveComponentUsage)

	// Основные маршруты компонентов (параметризованные идут последними)
	adminRoutes.Post("/custom-components", s.marketplace.CustomComponents.CreateComponent)
	adminRoutes.Get("/custom-components", s.marketplace.CustomComponents.ListComponents)
	adminRoutes.Get("/custom-components/:id", s.marketplace.CustomComponents.GetComponent)
	adminRoutes.Put("/custom-components/:id", s.marketplace.CustomComponents.UpdateComponent)
	adminRoutes.Delete("/custom-components/:id", s.marketplace.CustomComponents.DeleteComponent)

	adminRoutes.Get("/categories/:category_id/components", s.marketplace.CustomComponents.GetCategoryComponents)

	// Маршруты для групп атрибутов
	adminRoutes.Get("/attribute-groups", s.marketplace.MarketplaceHandler.ListAttributeGroups)
	adminRoutes.Post("/attribute-groups", s.marketplace.MarketplaceHandler.CreateAttributeGroup)
	adminRoutes.Get("/attribute-groups/:id", s.marketplace.MarketplaceHandler.GetAttributeGroup)
	adminRoutes.Put("/attribute-groups/:id", s.marketplace.MarketplaceHandler.UpdateAttributeGroup)
	adminRoutes.Delete("/attribute-groups/:id", s.marketplace.MarketplaceHandler.DeleteAttributeGroup)
	adminRoutes.Get("/attribute-groups/:id/items", s.marketplace.MarketplaceHandler.GetAttributeGroupWithItems)
	adminRoutes.Post("/attribute-groups/:id/items", s.marketplace.MarketplaceHandler.AddItemToGroup)
	adminRoutes.Delete("/attribute-groups/:id/items/:attributeId", s.marketplace.MarketplaceHandler.RemoveItemFromGroup)

	// Маршруты для привязки групп к категориям
	adminRoutes.Get("/categories/:id/attribute-groups", s.marketplace.MarketplaceHandler.GetCategoryGroups)
	adminRoutes.Post("/categories/:id/attribute-groups", s.marketplace.MarketplaceHandler.AttachGroupToCategory)
	adminRoutes.Delete("/categories/:id/attribute-groups/:groupId", s.marketplace.MarketplaceHandler.DetachGroupFromCategory)

	// Использовать реальный обработчик из UserHandler

	// Управление администраторами

	// Обновлено: маршруты админских функций используют обработчик индексации
	adminRoutes.Post("/reindex-listings", s.marketplace.Indexing.ReindexAll)
	adminRoutes.Post("/reindex-listings-with-translations", s.marketplace.Indexing.ReindexAllWithTranslations)
	adminRoutes.Post("/sync-discounts", s.marketplace.Listings.SynchronizeDiscounts) // Оставляем в Listings, т.к. это работа с объявлениями
	adminRoutes.Post("/reindex-ratings", s.marketplace.Indexing.ReindexRatings)

	chat := authedAPIGroup.Group("/marketplace/chat")
	chat.Get("/", s.marketplace.Chat.GetChats)
	chat.Get("/messages", s.marketplace.Chat.GetMessages)

	// Применяем rate limiting для отправки сообщений и загрузки файлов
	chat.Post("/messages", s.middleware.RateLimitMessages(), s.marketplace.Chat.SendMessage)
	chat.Put("/messages/read", s.marketplace.Chat.MarkAsRead)
	chat.Post("/:chat_id/archive", s.marketplace.Chat.ArchiveChat)

	// Роуты для работы с вложениями с rate limiting
	chat.Post("/messages/:id/attachments", s.middleware.RateLimitMessages(), s.marketplace.Chat.UploadAttachments)
	chat.Get("/attachments/:id", s.marketplace.Chat.GetAttachment)
	chat.Delete("/attachments/:id", s.marketplace.Chat.DeleteAttachment)
	chat.Get("/unread-count", s.marketplace.Chat.GetUnreadCount)

}

// registerProjectRoutes регистрирует роуты проектов через новую систему
func (s *Server) registerProjectRoutes() {
	// Создаем слайс всех проектов, которые реализуют RouteRegistrar
	var registrars []RouteRegistrar

	// Добавляем notifications проект
	registrars = append(registrars, s.notifications, s.users, s.review)

	// Регистрируем роуты каждого проекта
	for _, registrar := range registrars {
		err := registrar.RegisterRoutes(s.app, s.middleware)
		if err != nil {
			logger.Error().
				Err(err).
				Str("prefix", registrar.GetPrefix()).
				Msg("Ошибка регистрации роутов проекта")
		} else {
			logger.Info().
				Str("prefix", registrar.GetPrefix()).
				Msg("Роуты проекта успешно зарегистрированы")
		}
	}
}

func (s *Server) Start() error {
	return s.app.Listen(fmt.Sprintf(":%s", s.cfg.Port))
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}
