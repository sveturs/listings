// backend/internal/server/server.go
package server

import (
	marketplaceHandler "backend/internal/proj/marketplace/handler"
	reviewHandler "backend/internal/proj/reviews/handler"
	userHandler "backend/internal/proj/users/handler"
//	service "backend/internal/proj/marketplace/service"  

	"github.com/gofiber/websocket/v2"

	globalService "backend/internal/proj/global/service"
	notificationHandler "backend/internal/proj/notifications/handler"

	"backend/internal/config"
	"backend/internal/middleware"
	"backend/internal/storage/postgres"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

// Определяем структуру Server перед использованием
type Server struct {
	app           *fiber.App
	cfg           *config.Config
	users         *userHandler.Handler
	middleware    *middleware.Middleware
	review        *reviewHandler.Handler
	marketplace   *marketplaceHandler.Handler
	notifications *notificationHandler.Handler
}

func NewServer(cfg *config.Config) (*Server, error) {
    db, err := postgres.NewDatabase(cfg.DatabaseURL)
    if err != nil {
        return nil, fmt.Errorf("failed to initialize database: %w", err)
    }

    // Create global services
    services := globalService.NewService(db, cfg)

    usersHandler := userHandler.NewHandler(services)
    reviewHandler := reviewHandler.NewHandler(services)
    marketplaceHandler := marketplaceHandler.NewHandler(services)
    notificationsHandler := notificationHandler.NewHandler(services)

    middleware := middleware.NewMiddleware(cfg, services)

    // Initialize Fiber with file upload size limit
    app := fiber.New(fiber.Config{
        ErrorHandler: middleware.ErrorHandler,
        BodyLimit:    20 * 1024 * 1024, // 20MB
    })

    // Apply middleware
    middleware.Setup(app)

    // Initialize server
    server := &Server{
        app: app,
        cfg: cfg,
        users:         usersHandler,
        middleware:    middleware,
        review:        reviewHandler,
        marketplace:   marketplaceHandler,
        notifications: notificationsHandler,
    }
    notificationsHandler.Notification.ConnectTelegramWebhook()

    // Setup routes
    server.setupRoutes()

    return server, nil
}
// setupRoutes настраивает маршруты сервера
func (s *Server) setupRoutes() {
	// Root path
	s.app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hostel Booking System API")
	})
	s.app.Post("/api/v1/notifications/telegram/webhook", func(c *fiber.Ctx) error {
		log.Printf("Received webhook request: %s", string(c.Body()))
		return s.notifications.Notification.HandleTelegramWebhook(c)
	})
	// Static files
	s.app.Static("/uploads", "./uploads")
	s.app.Static("/public", "./public")
	s.app.Get("/service-worker.js", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/javascript")
		return c.SendFile("./public/service-worker.js")
	})
	os.MkdirAll("./uploads", os.ModePerm)
	os.MkdirAll("./public", os.ModePerm)

	// Публичные маршруты маркетплейса
	marketplace := s.app.Group("/api/v1/marketplace")
	marketplace.Get("/listings", s.marketplace.Marketplace.GetListings)
	marketplace.Get("/categories", s.marketplace.Marketplace.GetCategories)
	marketplace.Get("/category-tree", s.marketplace.Marketplace.GetCategoryTree)

	marketplace.Get("/listings/:id", s.marketplace.Marketplace.GetListing) // Детали товара

	// Публичные маршруты для отзывов
	review := s.app.Group("/api/v1/reviews")
	review.Get("/", s.review.Review.GetReviews)       // Получение списка отзывов
	review.Get("/:id", s.review.Review.GetReviewByID) // Получение отдельного отзыва
	review.Get("/stats", s.review.Review.GetStats)    // Статистика по отзывам

	// Auth routes
	auth := s.app.Group("/auth")
	auth.Get("/session", s.users.Auth.GetSession) // Исправлено
	auth.Get("/google", s.users.Auth.GoogleAuth)
	auth.Get("/google/callback", s.users.Auth.GoogleCallback)
	auth.Get("/logout", s.users.Auth.Logout)

	// Protected API routes
	api := s.app.Group("/api/v1", s.middleware.AuthRequired)

	// Маршруты для отзывов, требующие авторизации
	protectedReviews := s.app.Group("/api/v1/reviews", s.middleware.AuthRequired)
	protectedReviews.Post("/", s.review.Review.CreateReview)            // Создание отзыва
	protectedReviews.Put("/:id", s.review.Review.UpdateReview)          // Обновление отзыва
	protectedReviews.Delete("/:id", s.review.Review.DeleteReview)       // Удаление отзыва
	protectedReviews.Post("/:id/vote", s.review.Review.VoteForReview)   // Голосование за отзыв
	protectedReviews.Post("/:id/response", s.review.Review.AddResponse) // Добавление ответа на отзыв
	protectedReviews.Post("/:id/photos", s.review.Review.UploadPhotos)  // Загрузка фотографий к отзыву

	// Маршруты для статистики по сущностям
	entityStats := s.app.Group("/api/v1/entity")
	entityStats.Get("/:type/:id/rating", s.review.Review.GetEntityRating) // Получение рейтинга сущности
	entityStats.Get("/:type/:id/stats", s.review.Review.GetEntityStats)   // Получение статистики по отзывам сущности

	// Protected user routes
	users := s.app.Group("/api/v1/users")
	users.Post("/register", s.users.User.Register)    // Исправлено
	users.Get("/me", s.users.User.GetProfile)         // Исправлено
	users.Put("/me", s.users.User.UpdateProfile)      // Исправлено
	users.Get("/profile", s.users.User.GetProfile)    // Исправлено
	users.Put("/profile", s.users.User.UpdateProfile) // Исправлено

	// Защищенные маршруты маркетплейса
	marketplaceProtected := api.Group("/marketplace")
	marketplaceProtected.Post("/listings", s.marketplace.Marketplace.CreateListing)
	marketplaceProtected.Get("/listings/:id", s.marketplace.Marketplace.GetListing)
	marketplaceProtected.Put("/listings/:id", s.marketplace.Marketplace.UpdateListing)
	marketplaceProtected.Delete("/listings/:id", s.marketplace.Marketplace.DeleteListing)
	marketplaceProtected.Post("/listings/:id/images", s.marketplace.Marketplace.UploadImages)
	marketplaceProtected.Post("/listings/:id/favorite", s.marketplace.Marketplace.AddToFavorites)
	marketplaceProtected.Delete("/listings/:id/favorite", s.marketplace.Marketplace.RemoveFromFavorites)
	marketplaceProtected.Get("/favorites", s.marketplace.Marketplace.GetFavorites)
	// Чат для маркетплейса

	chat := api.Group("/marketplace/chat")
	chat.Get("/", s.marketplace.Chat.GetChats)
	chat.Get("/:listing_id/messages", s.marketplace.Chat.GetMessages)

	chat.Post("/messages", s.marketplace.Chat.SendMessage)
	chat.Put("/messages/read", s.marketplace.Chat.MarkAsRead)
	chat.Post("/:chat_id/archive", s.marketplace.Chat.ArchiveChat)
	chat.Get("/unread-count", s.marketplace.Chat.GetUnreadCount)

	notifications := api.Group("/notifications")
	notifications.Get("/", s.notifications.Notification.GetNotifications)
	notifications.Get("/settings", s.notifications.Notification.GetSettings)
	notifications.Put("/settings", s.notifications.Notification.UpdateSettings)
	notifications.Get("/telegram", s.notifications.Notification.GetTelegramStatus)

	notifications.Put("/:id/read", s.notifications.Notification.MarkAsRead)
 	notifications.Post("/telegram/token", func(c *fiber.Ctx) error {
        // Добавляем логирование для отладки
        log.Printf("Handling telegram token request for user: %v", c.Locals("user_id"))
        response := s.notifications.Notification.GetTelegramToken(c)
        log.Printf("Token generation response: %v", response)
        return response
    })
	notifications.Post("/test", s.notifications.Notification.SendTestNotification)

	// WebSocket эндпоинт
	s.app.Use("/ws/chat", s.middleware.AuthRequired) // Защищаем WebSocket
	s.app.Get("/ws/chat", websocket.New(s.marketplace.Chat.HandleWebSocket))

}
func (s *Server) Start() error {
	return s.app.Listen(fmt.Sprintf(":%s", s.cfg.Port))
}

// Shutdown останавливает сервер
func (s *Server) Shutdown(ctx context.Context) error {
	return s.app.Shutdown()
}
