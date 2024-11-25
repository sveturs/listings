package main

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoomImage struct {
	ID          int       `json:"id"`
	RoomID      int       `json:"room_id"`
	FilePath    string    `json:"file_path"`
	FileName    string    `json:"file_name"`
	FileSize    int       `json:"file_size"`
	ContentType string    `json:"content_type"`
	IsMain      bool      `json:"is_main"`
	CreatedAt   time.Time `json:"created_at"`
}

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("CORS_ORIGINS"),
		AllowMethods:     "GET,POST,DELETE,PUT",
		AllowHeaders:     "Origin, Content-Type, Accept",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: true,
	}))

	os.MkdirAll("./uploads", os.ModePerm)
	// Подключение к базе данных
	dbURL := os.Getenv("DATABASE_URL")
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer pool.Close()
	processImage := func(file *multipart.FileHeader) (string, error) {
		src, err := file.Open()
		if err != nil {
			return "", err
		}
		defer src.Close()

		//  уникальное имя файла
		ext := filepath.Ext(file.Filename)
		fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		filePath := filepath.Join("uploads", fileName)

		//  файл для сохранения
		dst, err := os.Create(filePath)
		if err != nil {
			return "", err
		}
		defer dst.Close()

		//  изображение для обработки
		img, err := imaging.Decode(src)
		if err != nil {
			return "", err
		}

		// Изменяем размер изображения (например, максимальная ширина 1200px)
		resized := imaging.Resize(img, 1200, 0, imaging.Lanczos)

		//    обработанное изображение
		err = imaging.Save(resized, filePath)
		if err != nil {
			return "", err
		}

		return fileName, nil
	}
	// Удаление бронирования комнаты
	app.Delete("/rooms/:roomId/bookings/:bookingId", func(c *fiber.Ctx) error {
		bookingID := c.Params("bookingId")
		_, err := pool.Exec(context.Background(), "DELETE FROM bookings WHERE id = $1", bookingID)
		if err != nil {
			return c.Status(500).SendString("Ошибка удаления бронирования")
		}
		return c.SendString("Бронирование удалено")
	})

	// Удаление бронирования койко-места
	app.Delete("/beds/:bedId/bookings/:bookingId", func(c *fiber.Ctx) error {
		bookingID := c.Params("bookingId")
		_, err := pool.Exec(context.Background(), "DELETE FROM bed_bookings WHERE id = $1", bookingID)
		if err != nil {
			return c.Status(500).SendString("Ошибка удаления бронирования")
		}
		return c.SendString("Бронирование удалено")
	})
	//  эндпоинт для загрузки изображений
	app.Post("/rooms/:id/images", func(c *fiber.Ctx) error {
		log.Printf("Начало загрузки изображений")
		roomID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			log.Printf("Ошибка преобразования ID комнаты: %v", err)
			return c.Status(400).SendString("Неверный ID комнаты")
		}

		form, err := c.MultipartForm()
		if err != nil {
			log.Printf("Ошибка получения формы: %v", err)
			return c.Status(400).SendString("Ошибка получения файлов")
		}

		files := form.File["images"]
		log.Printf("Получено %d файлов", len(files))

		isMain := len(files) > 0 // Первое изображение будет главным

		var uploadedImages []RoomImage
		for _, file := range files {
			// Проверяем тип файла
			if !strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
				return c.Status(400).SendString("Допустимы только изображения")
			}

			// Проверяем размер файла (например, максимум 5MB)
			if file.Size > 5*1024*1024 {
				return c.Status(400).SendString("Размер файла не должен превышать 5MB")
			}

			// Обрабатываем и сохраняем изображение
			fileName, err := processImage(file)
			if err != nil {
				log.Printf("Ошибка обработки изображения: %v", err)
				return c.Status(500).SendString("Ошибка обработки изображения")
			}

			// Сохраняем информацию в базу данных
			var imageID int
			err = pool.QueryRow(context.Background(), `
                INSERT INTO room_images (room_id, file_path, file_name, file_size, content_type, is_main)
                VALUES ($1, $2, $3, $4, $5, $6)
                RETURNING id
            `, roomID, fileName, file.Filename, file.Size, file.Header.Get("Content-Type"), isMain).Scan(&imageID)

			if err != nil {
				log.Printf("Ошибка сохранения информации об изображении: %v", err)
				return c.Status(500).SendString("Ошибка сохранения информации об изображении")
			}

			uploadedImages = append(uploadedImages, RoomImage{
				ID:          imageID,
				RoomID:      roomID,
				FilePath:    fileName,
				FileName:    file.Filename,
				FileSize:    int(file.Size),
				ContentType: file.Header.Get("Content-Type"),
				IsMain:      isMain,
			})

			isMain = false // Только первое изображение главное
		}

		return c.JSON(uploadedImages)
	})

	// Получение изображений комнаты
	app.Get("/rooms/:id/images", func(c *fiber.Ctx) error {
		roomID := c.Params("id")
		log.Printf("Получение изображений для комнаты: %s", roomID)

		rows, err := pool.Query(context.Background(), `
			SELECT id, room_id, file_path, file_name, file_size, content_type, is_main, created_at
			FROM room_images
			WHERE room_id = $1
			ORDER BY is_main DESC, created_at DESC
		`, roomID)
		if err != nil {
			log.Printf("Ошибка запроса изображений: %v", err)
			return c.Status(500).SendString("Ошибка получения изображений")
		}
		defer rows.Close()

		var images []RoomImage
		for rows.Next() {
			var img RoomImage
			err := rows.Scan(
				&img.ID,
				&img.RoomID,
				&img.FilePath,
				&img.FileName,
				&img.FileSize,
				&img.ContentType,
				&img.IsMain,
				&img.CreatedAt,
			)
			if err != nil {
				log.Printf("Ошибка сканирования изображения: %v", err)
				continue
			}
			images = append(images, img)
		}

		if len(images) == 0 {
			log.Printf("Изображения не найдены для комнаты: %s", roomID)
		} else {
			log.Printf("Найдено %d изображений для комнаты: %s", len(images), roomID)
		}

		return c.JSON(images)
	})

	// Удаление изображения
	app.Delete("/rooms/:roomId/images/:imageId", func(c *fiber.Ctx) error {
		imageID := c.Params("imageId")
		var filePath string
		err := pool.QueryRow(context.Background(), "SELECT file_path FROM room_images WHERE id = $1", imageID).Scan(&filePath)
		if err != nil {
			return c.Status(404).SendString("Изображение не найдено")
		}

		// Удаляем файл
		os.Remove(filepath.Join("uploads", filePath))

		// Удаляем запись из базы
		_, err = pool.Exec(context.Background(), "DELETE FROM room_images WHERE id = $1", imageID)
		if err != nil {
			return c.Status(500).SendString("Ошибка удаления изображения")
		}

		return c.SendString("Изображение удалено")
	})

	// Статическая раздача изображений
	app.Static("/uploads", "./uploads")

	// Главный маршрут
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hostel Booking System API с PostgreSQL")
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

		_, err := pool.Exec(context.Background(), "INSERT INTO users (name, email) VALUES ($1, $2)", user.Name, user.Email)
		if err != nil {
			if strings.Contains(err.Error(), "unique constraint") {
				return c.Status(400).SendString("Email уже используется")
			}
			log.Printf("Ошибка добавления пользователя: %v", err)
			return c.Status(500).SendString("Ошибка добавления пользователя")
		}

		return c.SendString("Пользователь добавлен успешно")
	})

	type Room struct {
		Name               string   `json:"name"`
		Capacity           int      `json:"capacity"`
		PricePerNight      *float64 `json:"price_per_night"`
		AddressStreet      string   `json:"address_street"`
		AddressCity        string   `json:"address_city"`
		AddressState       string   `json:"address_state"`
		AddressCountry     string   `json:"address_country"`
		AddressPostalCode  string   `json:"address_postal_code"`
		AccommodationType  string   `json:"accommodation_type"`
		IsShared           bool     `json:"is_shared"`
		TotalBeds          *int     `json:"total_beds"`
		AvailableBeds      *int     `json:"available_beds"`
		HasPrivateBathroom bool     `json:"has_private_bathroom"`
		Latitude           *float64 `json:"latitude"`
		Longitude          *float64 `json:"longitude"`
		FormattedAddress   string   `json:"formatted_address"`
	}

	type Bed struct {
		ID            int     `json:"id"`
		RoomID        int     `json:"room_id"`
		BedNumber     string  `json:"bed_number"` // изменено с int на string
		IsAvailable   bool    `json:"is_available"`
		PricePerNight float64 `json:"price_per_night"`
	}

	type BedBooking struct {
		ID        int    `json:"id"`
		BedID     int    `json:"bed_id"`
		UserID    int    `json:"user_id"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}
	app.Post("/rooms", func(c *fiber.Ctx) error {
		var room Room
		if err := c.BodyParser(&room); err != nil {
			return c.Status(400).SendString("Неверный формат данных")
		}

		// Устанавливаем значения по умолчанию для total_beds и available_beds
		totalBeds := 0
		availableBeds := 0
		if room.TotalBeds != nil {
			totalBeds = *room.TotalBeds
			// Для available_beds используем либо переданное значение, либо total_beds
			if room.AvailableBeds != nil {
				availableBeds = *room.AvailableBeds
			} else {
				availableBeds = totalBeds
			}
		}

		var roomId int
		err := pool.QueryRow(context.Background(), `
			INSERT INTO rooms (
				name, capacity, price_per_night,
				address_street, address_city, address_state,
				address_country, address_postal_code,
				accommodation_type, is_shared,
				total_beds, available_beds, has_private_bathroom,
				latitude, longitude, formatted_address
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
				$11, $12, $13, $14, $15, $16)
			RETURNING id
		`, room.Name, room.Capacity, room.PricePerNight,
			room.AddressStreet, room.AddressCity, room.AddressState,
			room.AddressCountry, room.AddressPostalCode,
			room.AccommodationType, room.IsShared,
			totalBeds, availableBeds, room.HasPrivateBathroom,
			room.Latitude, room.Longitude, room.FormattedAddress,
		).Scan(&roomId)

		if err != nil {
			if strings.Contains(err.Error(), "unique constraint") {
				return c.Status(400).SendString("Комната с таким названием уже существует")
			}
			log.Printf("Ошибка добавления комнаты: %v", err)
			return c.Status(500).SendString("Ошибка добавления комнаты")
		}

		return c.JSON(fiber.Map{"id": roomId})
	})

	// Добавление кровати
	// Обработчик создания кровати
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
			log.Printf("Ошибка парсинга данных: %v", err)
			return c.Status(400).SendString("Неверный формат данных")
		}

		// Проверяем, существует ли комната
		var roomExists bool
		err = pool.QueryRow(context.Background(),
			"SELECT EXISTS(SELECT 1 FROM rooms WHERE id = $1)",
			roomID).Scan(&roomExists)
		if err != nil || !roomExists {
			return c.Status(404).SendString("Комната не найдена")
		}

		// Добавляем кровать
		var bedID int
		err = pool.QueryRow(context.Background(), `
        INSERT INTO beds (room_id, bed_number, price_per_night, is_available) 
        VALUES ($1, $2, $3, true)
        RETURNING id`,
			roomID, bedReq.BedNumber, bedReq.PricePerNight).Scan(&bedID)

		if err != nil {
			log.Printf("Ошибка добавления кровати: %v", err)
			return c.Status(500).SendString("Ошибка добавления кровати")
		}

		return c.JSON(fiber.Map{
			"id":              bedID,
			"room_id":         roomID,
			"bed_number":      bedReq.BedNumber,
			"price_per_night": bedReq.PricePerNight,
			"is_available":    true})
	})

	// Получение доступных кроватей
	app.Get("/rooms/:id/available-beds", func(c *fiber.Ctx) error {
		roomID := c.Params("id")
		startDate := c.Query("start_date")
		endDate := c.Query("end_date")

		// Проверяем корректность дат
		if startDate == "" || endDate == "" {
			return c.Status(400).SendString("Необходимо указать даты")
		}

		// проверяет пересечение периодов бронирования
		query := `
    SELECT b.id, b.bed_number, b.price_per_night
    FROM beds b
    WHERE b.room_id = $1
    AND b.is_available = true
    AND NOT EXISTS (
        SELECT 1
        FROM bed_bookings bb
        WHERE bb.bed_id = b.id
        AND bb.status = 'confirmed'
        AND (
            (bb.start_date <= $3 AND bb.end_date >= $2) -- Проверяем пересечение периодов
        )
    )
    ORDER BY b.bed_number
`

		rows, err := pool.Query(context.Background(), query, roomID, startDate, endDate)
		if err != nil {
			log.Printf("Ошибка запроса доступных кроватей: %v", err)
			return c.Status(500).SendString("Ошибка получения списка кроватей")
		}
		defer rows.Close()

		var beds []Bed
		for rows.Next() {
			var bed Bed
			if err := rows.Scan(&bed.ID, &bed.BedNumber, &bed.PricePerNight); err != nil {
				log.Printf("Ошибка сканирования кровати: %v", err)
				continue
			}
			bed.RoomID, _ = strconv.Atoi(roomID)
			bed.IsAvailable = true
			beds = append(beds, bed)
		}

		if err := rows.Err(); err != nil {
			log.Printf("Ошибка при итерации по кроватям: %v", err)
			return c.Status(500).SendString("Ошибка обработки данных")
		}

		// Обновляем количество доступных кроватей в комнате
		_, err = pool.Exec(context.Background(), `
			UPDATE rooms 
			SET available_beds = $1
			WHERE id = $2
		`, len(beds), roomID)

		if err != nil {
			log.Printf("Ошибка обновления количества доступных кроватей: %v", err)
		}

		return c.JSON(beds)
	})

	app.Get("/rooms", func(c *fiber.Ctx) error {
		capacity := c.Query("capacity")
		startDate := c.Query("start_date")
		endDate := c.Query("end_date")
		minPrice := c.Query("min_price")
		maxPrice := c.Query("max_price")
		city := c.Query("city")
		country := c.Query("country")

		baseQuery := `
WITH room_availability AS (
    SELECT 
        r.id,
        COALESCE(r.total_beds, 0) as total_beds,
        CASE 
            WHEN r.accommodation_type = 'bed' THEN 
                COALESCE(
                    r.total_beds - COALESCE((
                        SELECT COUNT(DISTINCT bb.bed_id)
                        FROM beds b2
                        LEFT JOIN bed_bookings bb ON b2.id = bb.bed_id
                        WHERE b2.room_id = r.id
                        AND bb.status = 'confirmed'
                        AND bb.start_date <= $2
                        AND bb.end_date >= $1
                    ), 0),
                    r.total_beds
                )
            ELSE 
                CASE WHEN EXISTS (
                    SELECT 1 FROM bookings b
                    WHERE b.room_id = r.id
                    AND b.status = 'confirmed'
                    AND b.start_date <= $2
                    AND b.end_date >= $1
                ) THEN 0 ELSE 1 END
        END as available_count,
        CASE 
            WHEN r.accommodation_type = 'bed' THEN
                COALESCE(
                    (SELECT MIN(b3.price_per_night) 
                     FROM beds b3 
                     WHERE b3.room_id = r.id 
                     AND b3.is_available = true
                     AND b3.id NOT IN (
                        SELECT bb.bed_id
                        FROM bed_bookings bb
                        WHERE bb.status = 'confirmed'
                        AND bb.start_date <= $2
                        AND bb.end_date >= $1
                     )),
                    r.price_per_night
                )
            ELSE r.price_per_night
        END as actual_price
    FROM rooms r
)
SELECT 
    r.id, 
    r.name, 
    r.capacity,
    r.latitude,
    r.longitude,
    ra.actual_price as price_per_night,
    r.address_street, 
    r.address_city, 
    r.address_state,
    r.address_country, 
    r.address_postal_code,
    r.accommodation_type, 
    r.is_shared,
    COALESCE(r.total_beds, 0) as total_beds,
    COALESCE(ra.available_count, 0) as available_beds,
    r.has_private_bathroom,
    r.created_at
FROM rooms r
JOIN room_availability ra ON r.id = ra.id
WHERE 1=1
    AND (
        CASE 
            WHEN r.accommodation_type = 'bed' 
            THEN COALESCE(ra.available_count, 0) > 0 
            ELSE COALESCE(ra.available_count, 1) = 1
        END
    )`

		var conditions []string
		args := []interface{}{startDate, endDate}
		// Добавляем фильтр по типу размещения
		if accommodationType := c.Query("accommodation_type"); accommodationType != "" {
			conditions = append(conditions, fmt.Sprintf("r.accommodation_type = $%d", len(args)+1))
			args = append(args, accommodationType)

			// Дополнительная проверка для приватных комнат
			if accommodationType == "room" && c.Query("has_private_rooms") == "true" {
				conditions = append(conditions, "r.is_shared = false")
			}
		}
		// Фильтр по вместимости
		if capacity != "" {
			conditions = append(conditions, fmt.Sprintf("r.capacity >= $%d", len(args)+1))
			args = append(args, capacity)
		}

		// Фильтр по минимальной цене
		if minPrice != "" {
			minPriceFloat, err := strconv.ParseFloat(minPrice, 64)
			if err == nil {
				conditions = append(conditions, fmt.Sprintf("ra.actual_price >= $%d", len(args)+1))
				args = append(args, minPriceFloat)
			}
		}

		// Фильтр по максимальной цене
		if maxPrice != "" {
			maxPriceFloat, err := strconv.ParseFloat(maxPrice, 64)
			if err == nil {
				conditions = append(conditions, fmt.Sprintf("ra.actual_price <= $%d", len(args)+1))
				args = append(args, maxPriceFloat)
			}
		}

		// Фильтры по городу и стране
		if city != "" {
			conditions = append(conditions, fmt.Sprintf("r.address_city ILIKE $%d", len(args)+1))
			args = append(args, "%"+city+"%")
		}
		if country != "" {
			conditions = append(conditions, fmt.Sprintf("r.address_country ILIKE $%d", len(args)+1))
			args = append(args, "%"+country+"%")
		}

		// Добавляем условия WHERE, если они есть
		if len(conditions) > 0 {
			baseQuery += " AND " + strings.Join(conditions, " AND ")
		}

		// Добавляем сортировку
		baseQuery += " ORDER BY ra.available_count DESC, ra.actual_price ASC, r.created_at DESC"
		// Логирование запроса
		log.Printf("SQL Query: %s, Args: %v", baseQuery, args)

		// Выполнение запроса
		rows, err := pool.Query(context.Background(), baseQuery, args...)

		if err != nil {
			log.Printf("Ошибка выполнения запроса: %v", err)
			return c.Status(500).SendString("Ошибка получения списка комнат")
		}
		defer rows.Close()

		// Обработка результатов
		var rooms []map[string]interface{}
		for rows.Next() {
			var id, capacity int
			var totalBeds, availableBeds int
			var name, addressStreet, addressCity, addressState,
				addressCountry, addressPostalCode, accommodationType string
			var latitude, longitude float64
			var pricePerNight float64
			var isShared, hasPrivateBathroom bool
			var createdAt time.Time

			if err := rows.Scan(
				&id, &name, &capacity, &latitude, &longitude, &pricePerNight,
				&addressStreet, &addressCity, &addressState,
				&addressCountry, &addressPostalCode,
				&accommodationType, &isShared,
				&totalBeds, &availableBeds, // теперь это не может быть NULL
				&hasPrivateBathroom,
				&createdAt,
			); err != nil {
				log.Printf("Ошибка сканирования строки: %v", err)
				continue
			}

			rooms = append(rooms, map[string]interface{}{
				"id":                   id,
				"name":                 name,
				"capacity":             capacity,
				"latitude":             latitude,
				"longitude":            longitude,
				"price_per_night":      pricePerNight,
				"address_street":       addressStreet,
				"address_city":         addressCity,
				"address_state":        addressState,
				"address_country":      addressCountry,
				"address_postal_code":  addressPostalCode,
				"accommodation_type":   accommodationType,
				"is_shared":            isShared,
				"total_beds":           totalBeds,
				"available_beds":       availableBeds,
				"has_private_bathroom": hasPrivateBathroom,
				"created_at":           createdAt.Format("2006-01-02 15:04:05"),
			})
		}
		return c.JSON(rooms)
	})

	// Добавление изображений койко-места
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
		var uploadedImages []map[string]interface{}

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

			// Сохраняем информацию в базу данных
			var imageID int
			err = pool.QueryRow(context.Background(), `
            INSERT INTO bed_images (bed_id, file_path, file_name, file_size, content_type)
            VALUES ($1, $2, $3, $4, $5)
            RETURNING id
        `, bedID, fileName, file.Filename, file.Size, file.Header.Get("Content-Type")).Scan(&imageID)

			if err != nil {
				return c.Status(500).SendString("Ошибка сохранения информации об изображении")
			}

			uploadedImages = append(uploadedImages, map[string]interface{}{
				"id":           imageID,
				"bed_id":       bedID,
				"file_path":    fileName,
				"file_name":    file.Filename,
				"file_size":    file.Size,
				"content_type": file.Header.Get("Content-Type"),
			})
		}

		return c.JSON(uploadedImages)
	})

	// Получение изображений койко-места
	app.Get("/beds/:id/images", func(c *fiber.Ctx) error {
		bedID := c.Params("id")

		rows, err := pool.Query(context.Background(), `
        SELECT id, bed_id, file_path, file_name, file_size, content_type, created_at
        FROM bed_images
        WHERE bed_id = $1
        ORDER BY created_at DESC
    `, bedID)
		if err != nil {
			return c.Status(500).SendString("Ошибка получения изображений")
		}
		defer rows.Close()

		var images []map[string]interface{}
		for rows.Next() {
			var (
				id          int
				bedID       int
				filePath    string
				fileName    string
				fileSize    int
				contentType string
				createdAt   time.Time
			)

			err := rows.Scan(&id, &bedID, &filePath, &fileName, &fileSize, &contentType, &createdAt)
			if err != nil {
				continue
			}

			images = append(images, map[string]interface{}{
				"id":           id,
				"bed_id":       bedID,
				"file_path":    filePath,
				"file_name":    fileName,
				"file_size":    fileSize,
				"content_type": contentType,
				"created_at":   createdAt,
			})
		}

		return c.JSON(images)
	})
	// Добавление бронирования
	app.Post("/bookings", func(c *fiber.Ctx) error {
		type BookingRequest struct {
			UserID    int    `json:"user_id"`
			RoomID    int    `json:"room_id"`
			BedID     *int   `json:"bed_id,omitempty"`
			StartDate string `json:"start_date"`
			EndDate   string `json:"end_date"`
		}

		var booking BookingRequest
		if err := c.BodyParser(&booking); err != nil {
			return c.Status(400).SendString("Неверный формат данных")
		}
		if booking.StartDate == booking.EndDate {
			return c.Status(400).SendString("Дата выезда должна быть позже даты заезда")
		}
		tx, err := pool.Begin(context.Background())
		if err != nil {
			return c.Status(500).SendString("Ошибка начала транзакции")
		}
		defer tx.Rollback(context.Background())

		// Проверяем существование пользователя
		var userExists bool
		err = tx.QueryRow(context.Background(),
			"SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)",
			booking.UserID).Scan(&userExists)
		if err != nil || !userExists {
			return c.Status(400).SendString("Пользователь не найден")
		}

		// Получаем информацию о комнате
		var roomType string
		var isShared bool
		err = tx.QueryRow(context.Background(), `
            SELECT accommodation_type, is_shared 
            FROM rooms 
            WHERE id = $1`,
			booking.RoomID).Scan(&roomType, &isShared)
		if err != nil {
			return c.Status(500).SendString("Ошибка получения информации о комнате")
		}

		if roomType == "bed" {
			if booking.BedID == nil {
				return c.Status(400).SendString("Для койко-места необходимо указать ID кровати")
			}

			// Проверяем доступность койко-места
			var isAvailable bool
			err = tx.QueryRow(context.Background(), `
                SELECT is_available 
                FROM beds 
                WHERE id = $1 AND room_id = $2`,
				*booking.BedID, booking.RoomID).Scan(&isAvailable)
			if err != nil || !isAvailable {
				return c.Status(400).SendString("Койко-место недоступно")
			}

			// Проверяем, не забронировано ли койко-место на эти даты
			var conflictCount int
			err = tx.QueryRow(context.Background(), `
				SELECT COUNT(*) 
				FROM bed_bookings 
				WHERE bed_id = $1 
				AND status = 'confirmed'
				AND (
					(start_date <= $2 AND end_date >= $2) OR
					(start_date <= $3 AND end_date >= $3) OR
					(start_date >= $2 AND end_date <= $3)
				)`,
				*booking.BedID, booking.StartDate, booking.EndDate).Scan(&conflictCount)

			if err != nil || conflictCount > 0 {
				return c.Status(400).SendString("Койко-место уже забронировано на эти даты")
			}

			// Создаем бронирование койко-места
			_, err = tx.Exec(context.Background(), `
                INSERT INTO bed_bookings (bed_id, user_id, start_date, end_date, status)
                VALUES ($1, $2, $3, $4, 'confirmed')`,
				*booking.BedID, booking.UserID, booking.StartDate, booking.EndDate)

			// Обновляем количество доступных мест в комнате
			_, err = tx.Exec(context.Background(), `
                UPDATE rooms 
                SET available_beds = available_beds - 1
                WHERE id = $1`,
				booking.RoomID)
		} else {
			// Для комнат и квартир проверяем общую доступность
			var count int
			err = tx.QueryRow(context.Background(), `
                SELECT COUNT(*) 
                FROM bookings 
                WHERE room_id = $1 
                    AND start_date <= $3 
                    AND end_date >= $2
                    AND status = 'confirmed'`,
				booking.RoomID, booking.StartDate, booking.EndDate).Scan(&count)
			if err != nil || count > 0 {
				return c.Status(400).SendString("Помещение занято на указанные даты")
			}

			// Создаем обычное бронирование
			_, err = tx.Exec(context.Background(), `
                INSERT INTO bookings (user_id, room_id, start_date, end_date, status)
                VALUES ($1, $2, $3, $4, 'confirmed')`,
				booking.UserID, booking.RoomID, booking.StartDate, booking.EndDate)
		}

		if err != nil {
			log.Printf("Ошибка создания бронирования: %v", err)
			return c.Status(500).SendString("Ошибка создания бронирования")
		}

		if err = tx.Commit(context.Background()); err != nil {
			return c.Status(500).SendString("Ошибка фиксации транзакции")
		}

		return c.SendString("Бронирование создано успешно")
	})

	// Получение списка всех бронирований
	app.Get("/bookings", func(c *fiber.Ctx) error {
		// Получаем бронирования комнат
		roomBookingsQuery := `
            SELECT b.id, b.user_id, b.room_id, NULL as bed_id, 
                   b.start_date, b.end_date, b.status,
                   r.name as room_name, r.accommodation_type,
                   u.name as user_name, u.email as user_email
            FROM bookings b
            JOIN rooms r ON b.room_id = r.id
            JOIN users u ON b.user_id = u.id
        `
		roomRows, err := pool.Query(context.Background(), roomBookingsQuery)
		if err != nil {
			log.Printf("Ошибка получения бронирований комнат: %v", err)
			return c.Status(500).SendString("Ошибка получения списка бронирований")
		}
		defer roomRows.Close()

		// Получаем бронирования койко-мест
		bedBookingsQuery := `
            SELECT bb.id, bb.user_id, b.room_id, bb.bed_id,
                   bb.start_date, bb.end_date, bb.status,
                   r.name as room_name, r.accommodation_type,
                   u.name as user_name, u.email as user_email
            FROM bed_bookings bb
            JOIN beds b ON bb.bed_id = b.id
            JOIN rooms r ON b.room_id = r.id
            JOIN users u ON bb.user_id = u.id
        `
		bedRows, err := pool.Query(context.Background(), bedBookingsQuery)
		if err != nil {
			log.Printf("Ошибка получения бронирований койко-мест: %v", err)
			return c.Status(500).SendString("Ошибка получения списка бронирований")
		}
		defer bedRows.Close()

		var bookings []map[string]interface{}

		// Обработка бронирований комнат
		for roomRows.Next() {
			var (
				id, userID, roomID  int
				bedID               *int
				startDate, endDate  time.Time
				status, roomName    string
				accommodationType   string
				userName, userEmail string
			)

			if err := roomRows.Scan(
				&id, &userID, &roomID, &bedID, &startDate, &endDate, &status,
				&roomName, &accommodationType, &userName, &userEmail,
			); err != nil {
				log.Printf("Ошибка сканирования бронирования комнаты: %v", err)
				continue
			}

			bookings = append(bookings, map[string]interface{}{
				"id":           id,
				"user_id":      userID,
				"room_id":      roomID,
				"bed_id":       bedID,
				"start_date":   startDate.Format("2006-01-02"),
				"end_date":     endDate.Format("2006-01-02"),
				"status":       status,
				"room_name":    roomName,
				"type":         accommodationType,
				"user_name":    userName,
				"user_email":   userEmail,
				"booking_type": "room",
			})
		}

		// Обработка бронирований койко-мест
		for bedRows.Next() {
			var (
				id, userID, roomID  int
				bedID               int
				startDate, endDate  time.Time
				status, roomName    string
				accommodationType   string
				userName, userEmail string
			)

			if err := bedRows.Scan(
				&id, &userID, &roomID, &bedID, &startDate, &endDate, &status,
				&roomName, &accommodationType, &userName, &userEmail,
			); err != nil {
				log.Printf("Ошибка сканирования бронирования койко-места: %v", err)
				continue
			}

			bookings = append(bookings, map[string]interface{}{
				"id":           id,
				"user_id":      userID,
				"room_id":      roomID,
				"bed_id":       bedID,
				"start_date":   startDate.Format("2006-01-02"),
				"end_date":     endDate.Format("2006-01-02"),
				"status":       status,
				"room_name":    roomName,
				"type":         accommodationType,
				"user_name":    userName,
				"user_email":   userEmail,
				"booking_type": "bed",
			})
		}

		return c.JSON(bookings)
	})

	// Удаление комнаты
	app.Delete("/rooms/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		result, err := pool.Exec(context.Background(), "DELETE FROM rooms WHERE id=$1", id)
		if err != nil {
			log.Printf("Ошибка удаления комнаты: %v", err)
			return c.Status(500).SendString("Ошибка удаления комнаты")
		}
		if result.RowsAffected() == 0 {
			return c.Status(404).SendString("Комната не найдена")
		}
		return c.SendString("Комната успешно удалена")
	})

	// Запуск приложения
	log.Fatal(app.Listen("0.0.0.0:3000"))
}
