package main

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"hostel-backend/database"

	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

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

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,DELETE,PUT",
		AllowHeaders:     "Origin, Content-Type, Accept",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: false,
		MaxAge:           300,
	}))

	os.MkdirAll("./uploads", os.ModePerm)

	// Инициализация базы данных
	db, err := database.New(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Главный маршрут
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hostel Booking System API")
	})

	// Добавление пользователя
	app.Post("/users", func(c *fiber.Ctx) error {
		type User struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}

		var user User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(400).SendString("Неверный формат данных")
		}

		err := db.AddUser(c.Context(), user.Name, user.Email)
		if err != nil {
			if strings.Contains(err.Error(), "unique constraint") {
				return c.Status(400).SendString("Email уже используется")
			}
			return c.Status(500).SendString("Ошибка добавления пользователя")
		}

		return c.SendString("Пользователь добавлен успешно")
	})

	// Добавление комнаты
	app.Post("/rooms", func(c *fiber.Ctx) error {
		var room database.Room
		if err := c.BodyParser(&room); err != nil {
			return c.Status(400).SendString("Неверный формат данных")
		}

		roomID, err := db.AddRoom(c.Context(), room)
		if err != nil {
			if strings.Contains(err.Error(), "unique constraint") {
				return c.Status(400).SendString("Комната с таким названием уже существует")
			}
			return c.Status(500).SendString("Ошибка добавления комнаты")
		}

		return c.JSON(fiber.Map{"id": roomID})
	})

	// Получение списка комнат
	app.Get("/rooms", func(c *fiber.Ctx) error {
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

		rooms, err := db.GetRooms(c.Context(), filters)
		if err != nil {
			return c.Status(500).SendString("Ошибка получения списка комнат")
		}

		return c.JSON(rooms)
	})

	// Удаление комнаты
	app.Delete("/rooms/:id", func(c *fiber.Ctx) error {
		err := db.DeleteRoom(c.Context(), c.Params("id"))
		if err != nil {
			if err.Error() == "room not found" {
				return c.Status(404).SendString("Комната не найдена")
			}
			return c.Status(500).SendString("Ошибка удаления комнаты")
		}
		return c.SendString("Комната успешно удалена")
	})

	// Загрузка изображений комнаты
	app.Post("/rooms/:id/images", func(c *fiber.Ctx) error {
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

			imageID, err := db.AddRoomImage(c.Context(), roomID, image)
			if err != nil {
				return c.Status(500).SendString("Ошибка сохранения информации об изображении")
			}

			image.ID = imageID
			uploadedImages = append(uploadedImages, image)
			isMain = false
		}

		return c.JSON(uploadedImages)
	})

	// Получение изображений комнаты
	app.Get("/rooms/:id/images", func(c *fiber.Ctx) error {
		images, err := db.GetRoomImages(c.Context(), c.Params("id"))
		if err != nil {
			return c.Status(500).SendString("Ошибка получения изображений")
		}
		return c.JSON(images)
	})

	// Удаление изображения комнаты
	app.Delete("/rooms/:roomId/images/:imageId", func(c *fiber.Ctx) error {
		filePath, err := db.DeleteRoomImage(c.Context(), c.Params("imageId"))
		if err != nil {
			return c.Status(404).SendString("Изображение не найдено")
		}

		os.Remove(filepath.Join("uploads", filePath))
		return c.SendString("Изображение удалено")
	})

	// Добавление кровати
	app.Post("/rooms/:id/beds", func(c *fiber.Ctx) error {
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

		bedID, err := db.AddBed(c.Context(), roomID, bedReq.BedNumber, bedReq.PricePerNight)
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

	// Получение доступных кроватей
	app.Get("/rooms/:id/available-beds", func(c *fiber.Ctx) error {
		startDate := c.Query("start_date")
		endDate := c.Query("end_date")

		if startDate == "" || endDate == "" {
			return c.Status(400).SendString("Необходимо указать даты")
		}

		beds, err := db.GetAvailableBeds(c.Context(), c.Params("id"), startDate, endDate)
		if err != nil {
			return c.Status(500).SendString("Ошибка получения списка кроватей")
		}

		return c.JSON(beds)
	})

	// Загрузка изображений койко-места
	app.Post("/beds/:id/images", func(c *fiber.Ctx) error {
		bedID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(400).SendString("Неверный ID койко-места")
		}

		form, err := c.MultipartForm()
		if err != nil {
			return c.Status(400).SendString("Ошибка получения файлов")
		}

		files := form.File["images"]
		var uploadedImages []database.BedImage

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

			image := database.BedImage{
				BedID:       bedID,
				FilePath:    fileName,
				FileName:    file.Filename,
				FileSize:    int(file.Size),
				ContentType: file.Header.Get("Content-Type"),
			}

			imageID, err := db.AddBedImage(c.Context(), bedID, image)
			if err != nil {
				return c.Status(500).SendString("Ошибка сохранения информации об изображении")
			}

			image.ID = imageID
			uploadedImages = append(uploadedImages, image)
		}

		return c.JSON(uploadedImages)
	})

	// Получение изображений койко-места
	app.Get("/beds/:id/images", func(c *fiber.Ctx) error {
		images, err := db.GetBedImages(c.Context(), c.Params("id"))
		if err != nil {
			return c.Status(500).SendString("Ошибка получения изображений")
		}
		return c.JSON(images)
	})

	// Создание бронирования
	app.Post("/bookings", func(c *fiber.Ctx) error {
		var booking database.BookingRequest
		if err := c.BodyParser(&booking); err != nil {
			return c.Status(400).SendString("Неверный формат данных")
		}

		err := db.CreateBooking(c.Context(), booking)
		if err != nil {
			switch err.Error() {
			case "user not found":
				return c.Status(400).SendString("Пользователь не найден")
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

	// Получение всех бронирований
	app.Get("/bookings", func(c *fiber.Ctx) error {
		bookings, err := db.GetAllBookings(c.Context())
		if err != nil {
			return c.Status(500).SendString("Ошибка получения списка бронирований")
		}
		return c.JSON(bookings)
	})

	// Удаление бронирования комнаты
	app.Delete("/rooms/:roomId/bookings/:bookingId", func(c *fiber.Ctx) error {
		err := db.DeleteBooking(c.Context(), c.Params("bookingId"), "room")
		if err != nil {
			return c.Status(500).SendString("Ошибка удаления бронирования")
		}
		return c.SendString("Бронирование удалено")
	})

	// Удаление бронирования койко-места
	app.Delete("/beds/:bedId/bookings/:bookingId", func(c *fiber.Ctx) error {
		err := db.DeleteBooking(c.Context(), c.Params("bookingId"), "bed")
		if err != nil {
			return c.Status(500).SendString("Ошибка удаления бронирования")
		}
		return c.SendString("Бронирование удалено")
	})

	// Статическая раздача изображений
	app.Static("/uploads", "./uploads")

	// Запуск приложения
	log.Fatal(app.Listen("0.0.0.0:3000"))
}
