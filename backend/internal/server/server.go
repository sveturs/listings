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

	// CSRF токен
	s.app.Get("/api/v1/csrf-token", s.middleware.GetCSRFToken())

	// Применяем rate limiting для authentication endpoints
	// New JWT endpoints with consistent /api/auth prefix

	// API routes с общим rate limiting по пользователю (300 запросов в минуту)
	// Используем JWT как основной метод аутентификации
	authedAPIGroup := s.app.Group("/api/v1", s.middleware.AuthRequiredJWT, s.middleware.CSRFProtection(), s.middleware.RateLimitByUser(300, time.Minute))

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
	contacts := authedAPIGroup.Group("/contacts")
	contacts.Get("/", s.contacts.GetContacts)
	contacts.Post("/", s.contacts.AddContact)
	contacts.Put("/:contact_user_id", s.contacts.UpdateContactStatus)
	contacts.Delete("/:contact_user_id", s.contacts.RemoveContact)
	contacts.Get("/privacy", s.contacts.GetPrivacySettings)
	contacts.Put("/privacy", s.contacts.UpdatePrivacySettings)
	contacts.Get("/status/:contact_user_id", s.contacts.GetContactStatus)
}

// registerProjectRoutes регистрирует роуты проектов через новую систему
func (s *Server) registerProjectRoutes() {
	// Создаем слайс всех проектов, которые реализуют RouteRegistrar
	var registrars []RouteRegistrar

	// Добавляем notifications проект
	registrars = append(registrars, s.notifications, s.users, s.review, s.marketplace)

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
