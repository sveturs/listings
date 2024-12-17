//backend/internal/server/server.go
package server

import (

	userHandler 			"backend/internal/proj/users/handler"
	accommodationHandler 	"backend/internal/proj/accommodation/handler"


	"backend/internal/config"
	"backend/internal/handlers"
	"backend/internal/middleware"
	"backend/internal/services"
	"backend/internal/storage/postgres"
	"context"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
)

// Определяем структуру Server перед использованием
type Server struct {
	app        		*fiber.App
	cfg        		*config.Config
	handlers   		*handlers.Handler
	users      		*userHandler.Handler
	accommodation 	*accommodationHandler.Handler
	middleware 		*middleware.Middleware
}

// NewServer создает новый сервер
func NewServer(cfg *config.Config) (*Server, error) {
	// Инициализация базы данных
	db, err := postgres.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Инициализация сервисов
	services := services.NewServices(db, cfg)
	usersHandler := userHandler.NewHandler(services)
	accommodationHandler := accommodationHandler.NewHandler(services)


	// Инициализация обработчиков
	handlers := handlers.NewHandler(services)

	// Инициализация middleware
	middleware := middleware.NewMiddleware(cfg, services)

	// Инициализация Fiber
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// Применение middleware
	middleware.Setup(app)

	// Инициализация сервера
	server := &Server{
		app:        app,
		cfg:        cfg,
		handlers:   handlers,
		users:      usersHandler, // Исправлено здесь
		middleware: middleware,
		accommodation: accommodationHandler,
	}

 

	// Настройка маршрутов
	server.setupRoutes()

	return server, nil
}

// setupRoutes настраивает маршруты сервера
func (s *Server) setupRoutes() {
	// Root path
	s.app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hostel Booking System API")
	})

	// Static files
	s.app.Static("/uploads", "./uploads")
	os.MkdirAll("./uploads", os.ModePerm)

	// Пример маршрутов
	s.app.Get("/rooms", s.accommodation.Room.List)
	s.app.Get("/rooms/:id", s.accommodation.Room.Get)
	s.app.Get("/rooms/:id/images", s.accommodation.Room.ListImages)
	s.app.Get("/rooms/:id/available-beds", s.accommodation.Room.GetAvailableBeds)
	s.app.Get("/beds/:id/images", s.accommodation.Room.ListBedImages)

	// Публичные маршруты для автомобилей
	s.app.Get("/api/v1/cars/available", s.handlers.Cars.GetAvailableCars)
	s.app.Get("/api/v1/cars/:id/images", s.handlers.Cars.GetImages)
	s.app.Post("/api/v1/car-bookings", s.handlers.Cars.CreateBooking)

	// Публичные маршруты маркетплейса
	marketplace := s.app.Group("/api/v1/marketplace")
	marketplace.Get("/listings", s.handlers.Marketplace.GetListings)
	marketplace.Get("/listings/:id", s.handlers.Marketplace.GetListing)
	marketplace.Get("/categories", s.handlers.Marketplace.GetCategories)
	marketplace.Get("/category-tree", s.handlers.Marketplace.GetCategoryTree)
 
	// Публичные маршруты для отзывов
    reviews := s.app.Group("/api/v1/reviews")
    reviews.Get("/", s.handlers.Reviews.GetReviews)  // Получение списка отзывов
    reviews.Get("/:id", s.handlers.Reviews.GetReviewByID)  // Получение отдельного отзыва
    reviews.Get("/stats", s.handlers.Reviews.GetStats)  // Статистика по отзывам

	// Auth routes
	auth := s.app.Group("/auth")
	auth.Get("/session", s.users.Auth.GetSession) // Исправлено
	auth.Get("/google", s.users.Auth.GoogleAuth)
	auth.Get("/google/callback", s.users.Auth.GoogleCallback)
	auth.Get("/logout", s.users.Auth.Logout)

	// Protected API routes
	api := s.app.Group("/api/v1", s.middleware.AuthRequired)
	cars := api.Group("/cars")
	cars.Post("/", s.handlers.Cars.AddCar)
	cars.Post("/:id/images", s.handlers.Cars.UploadImages)

    // Маршруты для отзывов, требующие авторизации
    protectedReviews := s.app.Group("/api/v1/reviews", s.middleware.AuthRequired)
    protectedReviews.Post("/", s.handlers.Reviews.CreateReview)  // Создание отзыва
    protectedReviews.Put("/:id", s.handlers.Reviews.UpdateReview)  // Обновление отзыва
    protectedReviews.Delete("/:id", s.handlers.Reviews.DeleteReview)  // Удаление отзыва
    protectedReviews.Post("/:id/vote", s.handlers.Reviews.VoteForReview)  // Голосование за отзыв
    protectedReviews.Post("/:id/response", s.handlers.Reviews.AddResponse)  // Добавление ответа на отзыв
    protectedReviews.Post("/:id/photos", s.handlers.Reviews.UploadPhotos)  // Загрузка фотографий к отзыву

    // Маршруты для статистики по сущностям
    entityStats := s.app.Group("/api/v1/entity")
    entityStats.Get("/:type/:id/rating", s.handlers.Reviews.GetEntityRating)  // Получение рейтинга сущности
    entityStats.Get("/:type/:id/stats", s.handlers.Reviews.GetEntityStats)  // Получение статистики по отзывам сущности


	// Protected room routes
	rooms := api.Group("/rooms")
	rooms.Post("/", s.accommodation.Room.Create)
	rooms.Post("/:id/images", s.accommodation.Room.UploadImages)
	rooms.Delete("/:id/images/:imageId", s.accommodation.Room.DeleteImage)
	rooms.Post("/:id/beds", s.accommodation.Room.AddBed)
	rooms.Post("/:roomId/beds/:bedId/images", s.accommodation.Room.UploadBedImages)

	// Protected booking routes
	bookings := api.Group("/bookings")
	bookings.Post("/", s.accommodation.Booking.Create)
	bookings.Get("/", s.accommodation.Booking.List)
	bookings.Delete("/:id", s.accommodation.Booking.Delete)
	// Protected user routes
	users := s.app.Group("/api/v1/users")
	users.Post("/register", s.users.User.Register)      // Исправлено
	users.Get("/me", s.users.User.GetProfile)           // Исправлено
	users.Put("/me", s.users.User.UpdateProfile)        // Исправлено
	users.Get("/profile", s.users.User.GetProfile)      // Исправлено
	users.Put("/profile", s.users.User.UpdateProfile)   // Исправлено

	// Защищенные маршруты маркетплейса
	marketplaceProtected := api.Group("/marketplace")
	marketplaceProtected.Post("/listings", s.handlers.Marketplace.CreateListing)
	marketplaceProtected.Put("/listings/:id", s.handlers.Marketplace.UpdateListing)
	marketplaceProtected.Delete("/listings/:id", s.handlers.Marketplace.DeleteListing)
	marketplaceProtected.Post("/listings/:id/images", s.handlers.Marketplace.UploadImages)
	marketplaceProtected.Post("/listings/:id/favorite", s.handlers.Marketplace.AddToFavorites)
	marketplaceProtected.Delete("/listings/:id/favorite", s.handlers.Marketplace.RemoveFromFavorites)
}
func (s *Server) Start() error {
	return s.app.Listen(fmt.Sprintf(":%s", s.cfg.Port))
}

// Shutdown останавливает сервер
func (s *Server) Shutdown(ctx context.Context) error {
	return s.app.Shutdown()
}
