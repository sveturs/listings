package server

import (
    "context"
    "fmt"
    "os"
    "backend/internal/config"
    "backend/internal/handlers"
    "backend/internal/middleware"
    "backend/internal/services"
    "backend/internal/storage/postgres"
    "github.com/gofiber/fiber/v2"
)

type Server struct {
    app        *fiber.App
    cfg        *config.Config
    handlers   *handlers.Handler
    middleware *middleware.Middleware
}

func NewServer(cfg *config.Config) (*Server, error) {
    // Инициализация базы данных
    db, err := postgres.NewDatabase(cfg.DatabaseURL)
    if err != nil {
        return nil, fmt.Errorf("failed to initialize database: %w", err)
    }

    // Инициализация сервисов
    services := services.NewServices(db, cfg)

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

    server := &Server{
        app:        app,
        cfg:        cfg,
        handlers:   handlers,
        middleware: middleware,
    }

    server.setupRoutes()

    return server, nil
}

func (s *Server) setupRoutes() {
    // Root path
    s.app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hostel Booking System API")
    })

    // Static files
    s.app.Static("/uploads", "./uploads")
    os.MkdirAll("./uploads", os.ModePerm)

    // Public routes (without /api/v1 prefix)
    s.app.Get("/rooms", s.handlers.Rooms.List)
    s.app.Get("/rooms/:id", s.handlers.Rooms.Get)
    s.app.Get("/rooms/:id/images", s.handlers.Rooms.ListImages)
	s.app.Get("/rooms/:id/available-beds", s.handlers.Rooms.GetAvailableBeds) 
    s.app.Post("/rooms", s.handlers.Rooms.Create)
    
    // Auth routes
    s.app.Get("/auth/session", s.handlers.Auth.GetSession)
    s.app.Get("/auth/google", s.handlers.Auth.GoogleAuth)
    s.app.Get("/auth/google/callback", s.handlers.Auth.GoogleCallback)
    s.app.Get("/auth/logout", s.handlers.Auth.Logout)

    // API v1 routes
    api := s.app.Group("/api/v1")

    // Protected API routes
    protected := api.Use(s.middleware.AuthRequired)
    
    // Protected room routes
    rooms := protected.Group("/rooms")
    rooms.Get("/", s.handlers.Rooms.List)
    rooms.Post("/", s.handlers.Rooms.Create)
    rooms.Get("/:id", s.handlers.Rooms.Get)
    rooms.Post("/:id/images", s.handlers.Rooms.UploadImages)
    rooms.Get("/:id/images", s.handlers.Rooms.ListImages)
    rooms.Delete("/:id/images/:imageId", s.handlers.Rooms.DeleteImage)
    rooms.Post("/:id/beds", s.handlers.Rooms.AddBed)
    rooms.Get("/:id/available-beds", s.handlers.Rooms.GetAvailableBeds)
	s.app.Get("/beds/:id/images", s.handlers.Rooms.ListBedImages)  // Публичный маршрут
	rooms.Post("/:roomId/beds/:bedId/images", s.handlers.Rooms.UploadBedImages)

    // Protected booking routes
    bookings := protected.Group("/bookings")
    bookings.Post("/", s.handlers.Bookings.Create)
    bookings.Get("/", s.handlers.Bookings.List)
    bookings.Delete("/:id", s.handlers.Bookings.Delete)

    // Protected user routes
    users := protected.Group("/users")
    users.Post("/register", s.handlers.Users.Register)
    users.Get("/me", s.handlers.Users.GetProfile)
    users.Put("/me", s.handlers.Users.UpdateProfile)
}

func (s *Server) Start() error {
    return s.app.Listen(fmt.Sprintf(":%s", s.cfg.Port))
}

func (s *Server) Shutdown(ctx context.Context) error {
    return s.app.Shutdown()
}