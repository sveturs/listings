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

    // Public routes
    s.app.Get("/rooms", s.handlers.Rooms.List)
    s.app.Get("/rooms/:id", s.handlers.Rooms.Get)
    s.app.Get("/rooms/:id/images", s.handlers.Rooms.ListImages)
    s.app.Get("/rooms/:id/available-beds", s.handlers.Rooms.GetAvailableBeds)
    s.app.Get("/beds/:id/images", s.handlers.Rooms.ListBedImages)


    // Auth routes
    auth := s.app.Group("/auth")
    auth.Get("/session", s.handlers.Auth.GetSession)
    auth.Get("/google", s.handlers.Auth.GoogleAuth)
    auth.Get("/google/callback", s.handlers.Auth.GoogleCallback)
    auth.Get("/logout", s.handlers.Auth.Logout)

    // Protected API routes
    api := s.app.Group("/api/v1", s.middleware.AuthRequired)
    cars := api.Group("/cars")
    cars.Post("/", s.handlers.Cars.AddCar)
    cars.Get("/available", s.handlers.Cars.GetAvailableCars)
    cars.Post("/:id/images", s.handlers.Cars.UploadImages)
    cars.Get("/:id/images", s.handlers.Cars.GetImages)


    // Protected room routes
    rooms := api.Group("/rooms")
    rooms.Post("/", s.handlers.Rooms.Create)
    rooms.Post("/:id/images", s.handlers.Rooms.UploadImages)
    rooms.Delete("/:id/images/:imageId", s.handlers.Rooms.DeleteImage)
    rooms.Post("/:id/beds", s.handlers.Rooms.AddBed)
    rooms.Post("/:roomId/beds/:bedId/images", s.handlers.Rooms.UploadBedImages)

    // Protected booking routes
    bookings := api.Group("/bookings")
    bookings.Post("/", s.handlers.Bookings.Create)
    bookings.Get("/", s.handlers.Bookings.List)
    bookings.Delete("/:id", s.handlers.Bookings.Delete)

    // Protected user routes
    users := api.Group("/users")
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