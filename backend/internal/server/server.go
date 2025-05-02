// backend/internal/server/server.go
package server

import (
	"backend/internal/config"
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
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"log"
	"os"
	"time"
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
	translation   *marketplaceHandler.TranslationHandler
	fileStorage   filestorage.FileStorageInterface
}

func NewServer(cfg *config.Config) (*Server, error) {
	fileStorage, err := filestorage.NewFileStorage(cfg.FileStorage)
	if err != nil {
		log.Printf("Ошибка инициализации файлового хранилища: %v. Функции загрузки файлов могут быть недоступны.", err)
	}

	var osClient *opensearch.OpenSearchClient
	if cfg.OpenSearch.URL != "" {
		var err error
		osClient, err = opensearch.NewOpenSearchClient(opensearch.Config{
			URL:      cfg.OpenSearch.URL,
			Username: cfg.OpenSearch.Username,
			Password: cfg.OpenSearch.Password,
		})
		if err != nil {
			log.Printf("Ошибка подключения к OpenSearch: %v", err)
		} else {
			log.Println("Успешное подключение к OpenSearch")
		}
	} else {
		log.Println("OpenSearch URL не указан, поиск будет отключен")
	}

	db, err := postgres.NewDatabase(cfg.DatabaseURL, osClient, cfg.OpenSearch.MarketplaceIndex, fileStorage)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	var translationService marketplaceService.TranslationServiceInterface
	if cfg.GoogleTranslateAPIKey != "" && cfg.OpenAIAPIKey != "" {
		translationFactory, err := marketplaceService.NewTranslationServiceFactory(cfg.GoogleTranslateAPIKey, cfg.OpenAIAPIKey, db)
		if err != nil {
			log.Printf("Ошибка создания фабрики перевода: %v, будет использован только OpenAI", err)
			translationService, err = marketplaceService.NewTranslationService(cfg.OpenAIAPIKey)
			if err != nil {
				return nil, fmt.Errorf("failed to create translation service: %w", err)
			}
		} else {
			translationService = translationFactory
			log.Printf("Создана фабрика сервисов перевода с поддержкой Google Translate и OpenAI")
		}
	} else if cfg.OpenAIAPIKey != "" {
		var err error
		translationService, err = marketplaceService.NewTranslationService(cfg.OpenAIAPIKey)
		if err != nil {
			return nil, fmt.Errorf("failed to create translation service: %w", err)
		}
		log.Printf("Создан сервис перевода на базе OpenAI")
	} else {
		return nil, fmt.Errorf("не указан ни один API ключ для перевода")
	}

	services := globalService.NewService(db, cfg, translationService)

	usersHandler := userHandler.NewHandler(services)
	reviewHandler := reviewHandler.NewHandler(services)
	marketplaceHandler := marketplaceHandler.NewHandler(services)
	notificationsHandler := notificationHandler.NewHandler(services)
	balanceHandler := balanceHandler.NewHandler(services)
	storefrontHandler := storefrontHandler.NewHandler(services)
	middleware := middleware.NewMiddleware(cfg, services)
	geocodeHandler := geocodeHandler.NewGeocodeHandler(services.Geocode())

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Детальное логирование ошибки
			log.Printf("Error type: %T", err)
			log.Printf("Error details: %+v", err)
			log.Printf("Error handler called with path: %s", c.Path())

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
		marketplace:   marketplaceHandler,
		notifications: notificationsHandler,
		balance:       balanceHandler,
		storefront:    storefrontHandler,
		payments:      services.Payment(),
		geocode:       geocodeHandler,
		translation:   marketplaceHandler.Translation,
		fileStorage:   fileStorage,
	}

	notificationsHandler.Notification.ConnectTelegramWebhook()
	server.setupMiddleware()
	server.setupRoutes()

	return server, nil
}

func (s *Server) setupMiddleware() {
	s.app.Use(s.middleware.CORS())
	s.app.Use(s.middleware.Logger())
	s.app.Use("/ws", s.middleware.AuthRequired)
	os.MkdirAll("./uploads", os.ModePerm)
	os.MkdirAll("./public", os.ModePerm)
}

func (s *Server) setupRoutes() {
	s.app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hostel Booking System API")
	})

	s.app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	s.app.Get("/listings/*", func(c *fiber.Ctx) error {
		path := c.Params("*")
		log.Printf("Serving MinIO file: %s", path)
		minioUrl := fmt.Sprintf("http://localhost:9000/listings/%s", path)
		log.Printf("Redirecting to public MinIO URL: %s", minioUrl)
		return c.Redirect(minioUrl, 302)
	})

	s.app.Static("/uploads", "./uploads")
	s.app.Static("/public", "./public")

	s.app.Get("/service-worker.js", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/javascript")
		return c.SendFile("./public/service-worker.js")
	})

	s.app.Get("/ws/chat", websocket.New(s.marketplace.Chat.HandleWebSocket, websocket.Config{
		HandshakeTimeout: 10 * time.Second,
	}))

	s.app.Post("/reindex-ratings-public", s.marketplace.Marketplace.ReindexRatings)

	s.app.Post("/api/v1/notifications/telegram/webhook", func(c *fiber.Ctx) error {
		log.Printf("Received webhook request: %s", string(c.Body()))
		return s.notifications.Notification.HandleTelegramWebhook(c)
	})

	s.app.Post("/api/v1/public/send-email", s.notifications.Notification.SendPublicEmail)
	s.app.Get("/api/v1/public/storefronts/:id", s.storefront.Storefront.GetPublicStorefront)
	s.app.Get("/v1/notifications/telegram", s.notifications.Notification.GetTelegramStatus)

	balanceRoutes := s.app.Group("/api/v1/balance", s.middleware.AuthRequired)
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

	marketplace := s.app.Group("/api/v1/marketplace")
	marketplace.Get("/listings", s.marketplace.Marketplace.GetListings)
	marketplace.Get("/categories", s.marketplace.Marketplace.GetCategories)
	marketplace.Get("/category-tree", s.marketplace.Marketplace.GetCategoryTree)
	marketplace.Get("/listings/:id", s.marketplace.Marketplace.GetListing)
	marketplace.Get("/search", s.marketplace.Marketplace.SearchListingsAdvanced) // маршрут поиска
	marketplace.Get("/suggestions", s.marketplace.Marketplace.GetSuggestions)    // маршрут автодополнения
	marketplace.Get("/category-suggestions", s.marketplace.Marketplace.GetCategorySuggestions)
	marketplace.Get("/categories/:id/attributes", s.marketplace.Marketplace.GetCategoryAttributes)
	marketplace.Get("/listings/:id/price-history", s.marketplace.Marketplace.GetPriceHistory)
	marketplace.Get("/listings/:id/similar", s.marketplace.Marketplace.GetSimilarListings)
	marketplace.Get("/categories/:id/attribute-ranges", s.marketplace.Marketplace.GetAttributeRanges)

	translation := s.app.Group("/api/v1/translation")
	translation.Get("/limits", s.translation.GetTranslationLimits)
	translation.Post("/provider", s.translation.SetTranslationProvider)

	review := s.app.Group("/api/v1/reviews")
	review.Get("/", s.review.Review.GetReviews)
	review.Get("/:id", s.review.Review.GetReviewByID)
	review.Get("/stats", s.review.Review.GetStats)

	entityStats := s.app.Group("/api/v1/entity")
	entityStats.Get("/:type/:id/rating", s.review.Review.GetEntityRating)
	entityStats.Get("/:type/:id/stats", s.review.Review.GetEntityStats)

	auth := s.app.Group("/auth")
	auth.Get("/session", s.users.Auth.GetSession)
	auth.Get("/google", s.users.Auth.GoogleAuth)
	auth.Get("/google/callback", s.users.Auth.GoogleCallback)
	auth.Get("/logout", s.users.Auth.Logout)

	api := s.app.Group("/api/v1", s.middleware.AuthRequired)

	protectedReviews := api.Group("/reviews")
	protectedReviews.Post("/", s.review.Review.CreateReview)
	protectedReviews.Put("/:id", s.review.Review.UpdateReview)
	protectedReviews.Delete("/:id", s.review.Review.DeleteReview)
	protectedReviews.Post("/:id/vote", s.review.Review.VoteForReview)
	protectedReviews.Post("/:id/response", s.review.Review.AddResponse)
	protectedReviews.Post("/:id/photos", s.review.Review.UploadPhotos)

	storefronts := api.Group("/storefronts")
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

	api.Get("/users/:id/reviews", s.review.Review.GetUserReviews)
	api.Get("/users/:id/rating", s.review.Review.GetUserRatingSummary)
	api.Get("/storefronts/:id/reviews", s.review.Review.GetStorefrontReviews)
	api.Get("/storefronts/:id/rating", s.review.Review.GetStorefrontRatingSummary)

	geocodeApi := s.app.Group("/api/v1/geocode")
	geocodeApi.Get("/reverse", s.geocode.ReverseGeocode)

	citiesApi := s.app.Group("/api/v1/cities")
	citiesApi.Get("/suggest", s.geocode.GetCitySuggestions)

	users := api.Group("/users")
	users.Post("/register", s.users.User.Register)
	users.Get("/me", s.users.User.GetProfile)
	users.Put("/me", s.users.User.UpdateProfile)
	users.Get("/profile", s.users.User.GetProfile)
	users.Put("/profile", s.users.User.UpdateProfile)
	users.Get("/:id/profile", s.users.User.GetProfileByID)

	marketplaceProtected := api.Group("/marketplace")
	marketplaceProtected.Post("/listings", s.marketplace.Marketplace.CreateListing)
	marketplaceProtected.Put("/listings/:id", s.marketplace.Marketplace.UpdateListing)
	marketplaceProtected.Delete("/listings/:id", s.marketplace.Marketplace.DeleteListing)
	marketplaceProtected.Post("/listings/:id/images", s.marketplace.Marketplace.UploadImages)
	marketplaceProtected.Post("/listings/:id/favorite", s.marketplace.Marketplace.AddToFavorites)
	marketplaceProtected.Delete("/listings/:id/favorite", s.marketplace.Marketplace.RemoveFromFavorites)
	marketplaceProtected.Get("/favorites", s.marketplace.Marketplace.GetFavorites)
	marketplaceProtected.Put("/translations/:id", s.marketplace.Marketplace.UpdateTranslations)
	marketplaceProtected.Post("/translations/batch", s.marketplace.Marketplace.BatchTranslateListings)
	marketplaceProtected.Post("/moderate-image", s.marketplace.Marketplace.ModerateImage)
	marketplaceProtected.Post("/enhance-preview", s.marketplace.Marketplace.EnhancePreview)
	marketplaceProtected.Post("/enhance-images", s.marketplace.Marketplace.EnhanceImages)

	// Тестовый маршрут
	//	s.app.Get("/api/test-admin", s.middleware.AuthRequired, s.middleware.AdminRequired, func(c *fiber.Ctx) error {
	//		log.Println("Test admin endpoint called")
	//		return c.Status(200).JSON(fiber.Map{
	//			"message": "Admin test successful",
	//		})
	//	})

	// Административные маршруты
	adminRoutes := s.app.Group("/api/v1/admin")
	adminRoutes.Use(s.middleware.AuthRequired)
	adminRoutes.Use(s.middleware.AdminRequired)

	// Использовать реальный обработчик из UserHandler
	adminRoutes.Get("/users", s.users.User.GetAllUsers)
	adminRoutes.Get("/users/:id", s.users.User.GetUserByIDAdmin)
	adminRoutes.Put("/users/:id", s.users.User.UpdateUserAdmin)
	adminRoutes.Put("/users/:id/status", s.users.User.UpdateUserStatus)
	adminRoutes.Delete("/users/:id", s.users.User.DeleteUser)
	adminRoutes.Get("/users/:id/balance", s.users.User.GetUserBalance)
	adminRoutes.Get("/users/:id/transactions", s.users.User.GetUserTransactions)
	adminRoutes.Post("/reindex-listings", s.marketplace.Marketplace.ReindexAll)
	adminRoutes.Post("/reindex-listings-with-translations", s.marketplace.Marketplace.ReindexAllWithTranslations)
	adminRoutes.Post("/sync-discounts", s.marketplace.Marketplace.SynchronizeDiscounts)
	adminRoutes.Post("/reindex-ratings", s.marketplace.Marketplace.ReindexRatings)

	chat := api.Group("/marketplace/chat")
	chat.Get("/", s.marketplace.Chat.GetChats)
	chat.Get("/messages", s.marketplace.Chat.GetMessages)
	chat.Post("/messages", s.marketplace.Chat.SendMessage)
	chat.Put("/messages/read", s.marketplace.Chat.MarkAsRead)
	chat.Post("/:chat_id/archive", s.marketplace.Chat.ArchiveChat)
	chat.Get("/unread-count", s.marketplace.Chat.GetUnreadCount)

	notifications := api.Group("/notifications")
	notifications.Get("/", s.notifications.Notification.GetNotifications)
	notifications.Get("/settings", s.notifications.Notification.GetSettings)
	notifications.Put("/settings", s.notifications.Notification.UpdateSettings)
	notifications.Get("/telegram", s.notifications.Notification.GetTelegramStatus)
	notifications.Get("/telegram/token", s.notifications.Notification.GetTelegramToken)
	notifications.Put("/:id/read", s.notifications.Notification.MarkAsRead)
	notifications.Post("/telegram/token", s.notifications.Notification.GetTelegramToken)
}

func (s *Server) Start() error {
	return s.app.Listen(fmt.Sprintf(":%s", s.cfg.Port))
}

func (s *Server) Shutdown(ctx context.Context) error {
	services := globalService.NewService(nil, nil, nil)
	services.Shutdown()
	return s.app.Shutdown()
}
