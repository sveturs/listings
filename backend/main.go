package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"hostel-backend/auth"
	"hostel-backend/database"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"context"
	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

type Server struct {
	app  *fiber.App
	db   *database.Database
	auth *auth.AuthManager
}

// Вспомогательная функция для обработки изображений
func processImage(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	ext := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join("uploads", fileName)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	img, err := imaging.Decode(src)
	if err != nil {
		return "", err
	}

	resized := imaging.Resize(img, 1200, 0, imaging.Lanczos)
	err = imaging.Save(resized, filePath)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

// Вспомогательная функция для генерации токена сессии
func generateSessionToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func (s *Server) setupRoutes() {
	// Основные маршруты
	s.app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hostel Booking System API")
	})

	// Добавление пользователя
	s.app.Post("/users", func(c *fiber.Ctx) error {
		var user struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		if err := c.BodyParser(&user); err != nil {
			return c.Status(400).SendString("Неверный формат данных")
		}

		err := s.db.AddUser(c.Context(), user.Name, user.Email)
		if err != nil {
			if strings.Contains(err.Error(), "unique constraint") {
				return c.Status(400).SendString("Email уже используется")
			}
			return c.Status(500).SendString("Ошибка добавления пользователя")
		}
		return c.SendString("Пользователь добавлен успешно")
	})

	// Маршруты для комнат
	s.app.Post("/rooms", func(c *fiber.Ctx) error {
		var room database.Room
		if err := c.BodyParser(&room); err != nil {
			return c.Status(400).SendString("Неверный формат данных")
		}

		roomID, err := s.db.AddRoom(c.Context(), room)
		if err != nil {
			if strings.Contains(err.Error(), "unique constraint") {
				return c.Status(400).SendString("Комната с таким названием уже существует")
			}
			return c.Status(500).SendString("Ошибка добавления комнаты")
		}

		return c.JSON(fiber.Map{"id": roomID})
	})

	s.app.Get("/rooms", func(c *fiber.Ctx) error {
		filters := map[string]string{
			"capacity":           c.Query("capacity"),
			"start_date":         c.Query("start_date"),
			"end_date":           c.Query("end_date"),
			"min_price":          c.Query("min_price"),
			"max_price":          c.Query("max_price"),
			"city":               c.Query("city"),
			"country":            c.Query("country"),
			"accommodation_type": c.Query("accommodation_type"),
			"has_private_rooms":  c.Query("has_private_rooms"),
		}
		log.Printf("Handling /rooms request with filters: %+v", filters) // Добавить

		rooms, err := s.db.GetRooms(c.Context(), filters)
		if err != nil {
			log.Printf("Error getting rooms: %v", err) // Добавить
			return c.Status(500).SendString("Ошибка получения списка комнат")
		}
	
		log.Printf("Successfully retrieved %d rooms", len(rooms)) // Добавить
		log.Printf("Returning rooms: %+v", rooms)
		return c.JSON(rooms)
	})

	// Изображения комнат
	s.app.Post("/rooms/:id/images", func(c *fiber.Ctx) error {
		roomID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(400).SendString("Неверный ID комнаты")
		}

		form, err := c.MultipartForm()
		if err != nil {
			return c.Status(400).SendString("Ошибка получения файлов")
		}

		files := form.File["images"]
		isMain := len(files) > 0

		var uploadedImages []database.RoomImage
		for _, file := range files {
			if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
				return c.Status(400).SendString("Допустимы только изображения")
			}

			if file.Size > 5*1024*1024 {
				return c.Status(400).SendString("Размер файла не должен превышать 5MB")
			}

			fileName, err := processImage(file)
			if err != nil {
				return c.Status(500).SendString("Ошибка обработки изображения")
			}

			image := database.RoomImage{
				RoomID:      roomID,
				FilePath:    fileName,
				FileName:    file.Filename,
				FileSize:    int(file.Size),
				ContentType: file.Header.Get("Content-Type"),
				IsMain:      isMain,
			}

			imageID, err := s.db.AddRoomImage(c.Context(), roomID, image)
			if err != nil {
				return c.Status(500).SendString("Ошибка сохранения информации об изображении")
			}

			image.ID = imageID
			uploadedImages = append(uploadedImages, image)
			isMain = false
		}

		return c.JSON(uploadedImages)
	})

	s.app.Get("/rooms/:id/images", func(c *fiber.Ctx) error {
		images, err := s.db.GetRoomImages(c.Context(), c.Params("id"))
		if err != nil {
			return c.Status(500).SendString("Ошибка получения изображений")
		}
		return c.JSON(images)
	})

	s.app.Delete("/rooms/:roomId/images/:imageId", func(c *fiber.Ctx) error {
		filePath, err := s.db.DeleteRoomImage(c.Context(), c.Params("imageId"))
		if err != nil {
			return c.Status(404).SendString("Изображение не найдено")
		}

		os.Remove(filepath.Join("uploads", filePath))
		return c.SendString("Изображение удалено")
	})

	// Маршруты для кроватей
	s.app.Post("/rooms/:id/beds", func(c *fiber.Ctx) error {
		type BedRequest struct {
			BedNumber     string  `json:"bed_number"`
			PricePerNight float64 `json:"price_per_night"`
		}

		roomID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(400).SendString("Неверный ID комнаты")
		}

		var bedReq BedRequest
		if err := c.BodyParser(&bedReq); err != nil {
			return c.Status(400).SendString("Неверный формат данных")
		}

		bedID, err := s.db.AddBed(c.Context(), roomID, bedReq.BedNumber, bedReq.PricePerNight)
		if err != nil {
			if err.Error() == "room not found" {
				return c.Status(404).SendString("Комната не найдена")
			}
			return c.Status(500).SendString("Ошибка добавления кровати")
		}

		return c.JSON(fiber.Map{
			"id":              bedID,
			"room_id":         roomID,
			"bed_number":      bedReq.BedNumber,
			"price_per_night": bedReq.PricePerNight,
			"is_available":    true,
		})
	})

	s.app.Get("/rooms/:id/available-beds", func(c *fiber.Ctx) error {
		startDate := c.Query("start_date")
		endDate := c.Query("end_date")

		if startDate == "" || endDate == "" {
			return c.Status(400).SendString("Необходимо указать даты")
		}

		beds, err := s.db.GetAvailableBeds(c.Context(), c.Params("id"), startDate, endDate)
		if err != nil {
			return c.Status(500).SendString("Ошибка получения списка кроватей")
		}

		return c.JSON(beds)
	})

	// Маршруты для бронирований
	s.app.Post("/bookings", func(c *fiber.Ctx) error {
		// Получаем сессию пользователя
		sessionToken := c.Cookies("session_token")
		if sessionToken == "" {
			return c.Status(401).SendString("Необходима авторизация")
		}
	
		sessionData, ok := s.auth.GetSession(sessionToken)
		if !ok {
			return c.Status(401).SendString("Необходима авторизация")
		}
	
		var booking database.BookingRequest
		if err := c.BodyParser(&booking); err != nil {
			return c.Status(400).SendString("Неверный формат данных")
		}
	
		// Получаем id пользователя по email из сессии
		userId, err := s.db.GetUserIDByEmail(c.Context(), sessionData.Email)
		if err != nil {
			return c.Status(500).SendString("Ошибка получения данных пользователя")
		}
	
		// Устанавливаем ID пользователя из сессии
		booking.UserID = userId
	
		err = s.db.CreateBooking(c.Context(), booking)
		if err != nil {
			switch err.Error() {
			case "bed ID is required for bed booking":
				return c.Status(400).SendString("Для койко-места необходимо указать ID кровати")
			case "bed is not available":
				return c.Status(400).SendString("Койко-место недоступно")
			case "bed is already booked for these dates":
				return c.Status(400).SendString("Койко-место уже забронировано на эти даты")
			case "room is already booked for these dates":
				return c.Status(400).SendString("Помещение занято на указанные даты")
			case "check-out date must be after check-in date":
				return c.Status(400).SendString("Дата выезда должна быть позже даты заезда")
			default:
				return c.Status(500).SendString("Ошибка создания бронирования")
			}
		}
	
		return c.SendString("Бронирование создано успешно")
	})


	s.app.Get("/bookings", func(c *fiber.Ctx) error {
		bookings, err := s.db.GetAllBookings(c.Context())
		if err != nil {
			return c.Status(500).SendString("Ошибка получения списка бронирований")
		}
		return c.JSON(bookings)
	})

	s.app.Delete("/rooms/:roomId/bookings/:bookingId", func(c *fiber.Ctx) error {
		err := s.db.DeleteBooking(c.Context(), c.Params("bookingId"), "room")
		if err != nil {
			return c.Status(500).SendString("Ошибка удаления бронирования")
		}
		return c.SendString("Бронирование удалено")
	})

	s.app.Delete("/beds/:bedId/bookings/:bookingId", func(c *fiber.Ctx) error {
		err := s.db.DeleteBooking(c.Context(), c.Params("bookingId"), "bed")
		if err != nil {
			return c.Status(500).SendString("Ошибка удаления бронирования")
		}
		return c.SendString("Бронирование удалено")
	})

	// Статическая раздача изображений
	s.app.Static("/uploads", "./uploads")
}

func (s *Server) setupAuthRoutes() {
	s.app.Get("/auth/google", func(c *fiber.Ctx) error {
		url := s.auth.GetGoogleAuthURL()
		log.Printf("Generated Google Auth URL: %s", url) 
		return c.Redirect(url)
	})

	s.app.Get("/auth/google/callback", func(c *fiber.Ctx) error {
		code := c.Query("code")
		if code == "" {
			return c.Status(400).SendString("Missing code")
		}

		sessionData, err := s.auth.HandleGoogleCallback(c.Context(), code)
		if err != nil {
			return c.Status(500).SendString("Authentication failed")
		}

		// Получаем или создаём пользователя
		userID, err := s.db.GetOrCreateGoogleUser(
			c.Context(),
			sessionData.Name,
			sessionData.Email,
			sessionData.GoogleID,   // Должно приходить из auth.SessionData
			sessionData.PictureURL, // Должно приходить из auth.SessionData
		)
		if err != nil {
			log.Printf("Error managing user: %v", err)
			return c.Status(500).SendString("Error managing user")
		}

		sessionToken := generateSessionToken()
		sessionData.UserID = userID // Добавляем ID пользователя в данные сессии
		s.auth.SaveSession(sessionToken, sessionData)

		c.Cookie(&fiber.Cookie{
			Name:     "session_token",
			Value:    sessionToken,
			Path:     "/",
			MaxAge:   3600 * 24,
			Secure:   true,
			HTTPOnly: true,
			SameSite: "Lax",
		})

		return c.Redirect(os.Getenv("FRONTEND_URL"))
	})

	s.app.Get("/auth/session", func(c *fiber.Ctx) error {
		sessionToken := c.Cookies("session_token")
		if sessionToken == "" {
			return c.JSON(fiber.Map{
				"authenticated": false,
			})
		}

		sessionData, ok := s.auth.GetSession(sessionToken)
		if !ok {
			return c.JSON(fiber.Map{
				"authenticated": false,
			})
		}

		return c.JSON(fiber.Map{
			"authenticated": true,
			"user": fiber.Map{
				"name":     sessionData.Name,
				"email":    sessionData.Email,
				"provider": sessionData.Provider,
			},
		})
	})

	s.app.Get("/auth/logout", func(c *fiber.Ctx) error {
		sessionToken := c.Cookies("session_token")
		if sessionToken != "" {
			s.auth.DeleteSession(sessionToken)
			c.Cookie(&fiber.Cookie{
				Name:     "session_token",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				Secure:   true,
				HTTPOnly: true,
				SameSite: "Lax",
			})
		}
		return c.SendStatus(200)
	})
}

func NewServer() (*Server, error) {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }

	db, err := database.New(os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
    if err := db.Ping(context.Background()); err != nil {
        return nil, fmt.Errorf("cannot ping database: %v", err)
    }
    log.Println("Successfully connected to database")

	authManager := auth.NewAuthManager(
		os.Getenv("GOOGLE_CLIENT_ID"),
		os.Getenv("GOOGLE_CLIENT_SECRET"),
		os.Getenv("GOOGLE_OAUTH_REDIRECT_URL"),
	)

	app := fiber.New()

	// Настройка CORS
	//app.Use(cors.New(cors.Config{
	//	AllowOrigins:     os.Getenv("FRONTEND_URL"),
	//	AllowMethods:     "GET,POST,DELETE,PUT,OPTIONS",
	//	AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
	//	ExposeHeaders:    "Content-Length",
	//	AllowCredentials: true,
	//	MaxAge:           300,
	//}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000, http://localhost:3001, https://landhub.rs, http://landhub.rs", // Укажите домены, которые разрешены
		AllowCredentials: true, // Включить передачу cookie
		AllowMethods: "GET,POST,DELETE,PUT,OPTIONS", // Разрешить основные HTTP методы
		AllowHeaders: "Origin, Content-Type, Accept, Authorization", // Разрешить заголовки
		ExposeHeaders: "Content-Length", // Экспонировать определенные заголовки
	}))
	

	os.MkdirAll("./uploads", os.ModePerm)

	server := &Server{
		app:  app,
		db:   db,
		auth: authManager,
	}

	// Настройка маршрутов
	server.setupAuthRoutes()
	server.setupRoutes()

	return server, nil
}

func main() {
	server, err := NewServer()
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server running on port %s", port)
	log.Fatal(server.app.Listen(fmt.Sprintf(":%s", port)))
}
