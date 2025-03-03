// backend/internal/server/server.go

package server

import (
	"backend/internal/config"
	"backend/internal/middleware"
	balanceHandler "backend/internal/proj/balance/handler"
	globalService "backend/internal/proj/global/service"
	marketplaceHandler "backend/internal/proj/marketplace/handler"
	marketplaceService "backend/internal/proj/marketplace/service"
	notificationHandler "backend/internal/proj/notifications/handler"
	paymentService "backend/internal/proj/payments/service"
	reviewHandler "backend/internal/proj/reviews/handler"
	storefrontHandler "backend/internal/proj/storefront/handler"
	userHandler "backend/internal/proj/users/handler"
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
}

// Обновить функцию NewServer:
func NewServer(cfg *config.Config) (*Server, error) {
	// Инициализируем клиент OpenSearch
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

	// Инициализируем базу данных с OpenSearch
	db, err := postgres.NewDatabase(cfg.DatabaseURL, osClient, cfg.OpenSearch.MarketplaceIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	translationService, err := marketplaceService.NewTranslationService(cfg.OpenAIAPIKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create translation service: %w", err)
	}

	services := globalService.NewService(db, cfg, translationService)

	usersHandler := userHandler.NewHandler(services)
	reviewHandler := reviewHandler.NewHandler(services)
	marketplaceHandler := marketplaceHandler.NewHandler(services)
	notificationsHandler := notificationHandler.NewHandler(services)
	balanceHandler := balanceHandler.NewHandler(services)
	storefrontHandler := storefrontHandler.NewHandler(services) // Вот эта строка вызывает ошибку
	middleware := middleware.NewMiddleware(cfg, services)

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
		BodyLimit:    20 * 1024 * 1024, // 20MB
		// Добавляем конфигурацию для WebSocket
		EnableTrustedProxyCheck: true,
		TrustedProxies:          []string{"127.0.0.1", "::1"},
	})

	// Initialize server
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
	}

	// Инициализируем webhooks для телеграма
	notificationsHandler.Notification.ConnectTelegramWebhook()

	// Устанавливаем глобальные middleware до настройки роутов
	server.setupMiddleware()

	// Настраиваем роуты
	server.setupRoutes()

	return server, nil
}

func (s *Server) setupMiddleware() {
	// Глобальные middleware
	s.app.Use(s.middleware.CORS())
	s.app.Use(s.middleware.Logger())

	// WebSocket middleware должен быть настроен до роутов WebSocket
	s.app.Use("/ws", s.middleware.AuthRequired)

	// Создаем необходимые директории
	os.MkdirAll("./uploads", os.ModePerm)
	os.MkdirAll("./public", os.ModePerm)
}

func (s *Server) setupRoutes() {
	// Root path
	s.app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hostel Booking System API")
	})
	s.app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
	// Static files
	s.app.Static("/uploads", "./uploads")
	s.app.Static("/public", "./public")

	// Service worker
	s.app.Get("/service-worker.js", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/javascript")
		return c.SendFile("./public/service-worker.js")
	})

	// WebSocket endpoint должен быть настроен до других роутов
	s.app.Get("/ws/chat", websocket.New(s.marketplace.Chat.HandleWebSocket, websocket.Config{
		HandshakeTimeout: 10 * time.Second,
	}))

	// Telegram webhook
	s.app.Post("/api/v1/notifications/telegram/webhook", func(c *fiber.Ctx) error {
		log.Printf("Received webhook request: %s", string(c.Body()))
		return s.notifications.Notification.HandleTelegramWebhook(c)
	})
	// маршрут для витрин
	s.app.Get("/api/v1/public/storefronts/:id", s.storefront.Storefront.GetPublicStorefront)

	// Balance routes
	balanceRoutes := s.app.Group("/api/v1/balance", s.middleware.AuthRequired)
	balanceRoutes.Get("/", s.balance.Balance.GetBalance) // Добавляем .Balance
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

	// Public marketplace routes
	marketplace := s.app.Group("/api/v1/marketplace")
	marketplace.Get("/listings", s.marketplace.Marketplace.GetListings)
	marketplace.Get("/categories", s.marketplace.Marketplace.GetCategories)
	marketplace.Get("/category-tree", s.marketplace.Marketplace.GetCategoryTree)
	marketplace.Get("/listings/:id", s.marketplace.Marketplace.GetListing)
	marketplace.Get("/subcategories", s.marketplace.Marketplace.GetSubcategories)

	marketplace.Get("/search", s.marketplace.Marketplace.SearchListingsAdvanced) // маршрут поиска
	marketplace.Get("/suggestions", s.marketplace.Marketplace.GetSuggestions)    // маршрут автодополнения
 	marketplace.Get("/category-suggestions", s.marketplace.Marketplace.GetCategorySuggestions)


	// Public review routes
	review := s.app.Group("/api/v1/reviews")
	review.Get("/", s.review.Review.GetReviews)
	review.Get("/:id", s.review.Review.GetReviewByID)
	review.Get("/stats", s.review.Review.GetStats)

	// Auth routes
	auth := s.app.Group("/auth")
	auth.Get("/session", s.users.Auth.GetSession)
	auth.Get("/google", s.users.Auth.GoogleAuth)
	auth.Get("/google/callback", s.users.Auth.GoogleCallback)
	auth.Get("/logout", s.users.Auth.Logout)

	// Protected routes
	api := s.app.Group("/api/v1", s.middleware.AuthRequired)

	// Protected reviews routes
	protectedReviews := api.Group("/reviews")
	protectedReviews.Post("/", s.review.Review.CreateReview)
	protectedReviews.Put("/:id", s.review.Review.UpdateReview)
	protectedReviews.Delete("/:id", s.review.Review.DeleteReview)
	protectedReviews.Post("/:id/vote", s.review.Review.VoteForReview)
	protectedReviews.Post("/:id/response", s.review.Review.AddResponse)
	protectedReviews.Post("/:id/photos", s.review.Review.UploadPhotos)
	// маршруты для витрин
	storefronts := api.Group("/storefronts")
	storefronts.Get("/", s.storefront.Storefront.GetUserStorefronts)
	storefronts.Post("/", s.storefront.Storefront.CreateStorefront)
	storefronts.Get("/:id", s.storefront.Storefront.GetStorefront)
	storefronts.Put("/:id", s.storefront.Storefront.UpdateStorefront)
	storefronts.Delete("/:id", s.storefront.Storefront.DeleteStorefront)

	// Маршруты для источников импорта
	storefronts.Get("/:id/import-sources", s.storefront.Storefront.GetImportSources)
	storefronts.Post("/import-sources", s.storefront.Storefront.CreateImportSource)
	storefronts.Put("/import-sources/:id", s.storefront.Storefront.UpdateImportSource)
	storefronts.Delete("/import-sources/:id", s.storefront.Storefront.DeleteImportSource)

	// Маршруты для импорта данных
	storefronts.Post("/import-sources/:id/run", s.storefront.Storefront.RunImport)
	storefronts.Get("/import-sources/:id/history", s.storefront.Storefront.GetImportHistory)

	// Entity stats routes
	entityStats := api.Group("/entity")
	entityStats.Get("/:type/:id/rating", s.review.Review.GetEntityRating)
	entityStats.Get("/:type/:id/stats", s.review.Review.GetEntityStats)

	// Маршруты для отзывов пользователей и витрин
	api.Get("/users/:id/reviews", s.review.Review.GetUserReviews)
	api.Get("/users/:id/rating", s.review.Review.GetUserRatingSummary)
	api.Get("/storefronts/:id/reviews", s.review.Review.GetStorefrontReviews)
	api.Get("/storefronts/:id/rating", s.review.Review.GetStorefrontRatingSummary)

	// Protected user routes
	users := api.Group("/users")
	users.Post("/register", s.users.User.Register)
	users.Get("/me", s.users.User.GetProfile)
	users.Put("/me", s.users.User.UpdateProfile)
	users.Get("/profile", s.users.User.GetProfile)
	users.Put("/profile", s.users.User.UpdateProfile)
	users.Get("/:id/profile", s.users.User.GetProfileByID)

	// Protected marketplace routes
	marketplaceProtected := api.Group("/marketplace")
	marketplaceProtected.Post("/listings", s.marketplace.Marketplace.CreateListing)
	marketplaceProtected.Put("/listings/:id", s.marketplace.Marketplace.UpdateListing)
	marketplaceProtected.Delete("/listings/:id", s.marketplace.Marketplace.DeleteListing)
	marketplaceProtected.Post("/listings/:id/images", s.marketplace.Marketplace.UploadImages)
	marketplaceProtected.Post("/listings/:id/favorite", s.marketplace.Marketplace.AddToFavorites)
	marketplaceProtected.Delete("/listings/:id/favorite", s.marketplace.Marketplace.RemoveFromFavorites)
	marketplaceProtected.Get("/favorites", s.marketplace.Marketplace.GetFavorites)
	marketplaceProtected.Put("/translations/:id", s.marketplace.Marketplace.UpdateTranslations)

	// Административный маршрут для переиндексации

	api.Post("/admin/reindex-listings", s.middleware.AdminRequired, s.marketplace.Marketplace.ReindexAll)
	// Chat routes
	chat := api.Group("/marketplace/chat")
	chat.Get("/", s.marketplace.Chat.GetChats)
	chat.Get("/messages", s.marketplace.Chat.GetMessages)
	chat.Post("/messages", s.marketplace.Chat.SendMessage)
	chat.Put("/messages/read", s.marketplace.Chat.MarkAsRead)
	chat.Post("/:chat_id/archive", s.marketplace.Chat.ArchiveChat)
	chat.Get("/unread-count", s.marketplace.Chat.GetUnreadCount)

	// Notification routes
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
	return s.app.Shutdown()
}
