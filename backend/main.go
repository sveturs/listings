package main

import (
    "context"
    "fmt"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/disintegration/imaging"
    "log"
    "mime/multipart"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)
type RoomImage struct {
    ID          int    `json:"id"`
    RoomID      int    `json:"room_id"`
    FilePath    string `json:"file_path"`
    FileName    string `json:"file_name"`
    FileSize    int    `json:"file_size"`
    ContentType string `json:"content_type"`
    IsMain      bool   `json:"is_main"`
    CreatedAt   time.Time `json:"created_at"` // Изменено с string на time.Time
}
func main() {
	app := fiber.New()

	// Настройка CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3001,http://localhost:3000",
		AllowMethods: "GET,POST,DELETE,PUT",
		AllowHeaders: "Origin, Content-Type, Accept",
		ExposeHeaders: "Content-Length",
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

        // Генерируем уникальное имя файла
        ext := filepath.Ext(file.Filename)
        fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
        filePath := filepath.Join("uploads", fileName)

        // Создаем файл для сохранения
        dst, err := os.Create(filePath)
        if err != nil {
            return "", err
        }
        defer dst.Close()

        // Открываем изображение для обработки
        img, err := imaging.Decode(src)
        if err != nil {
            return "", err
        }

        // Изменяем размер изображения (например, максимальная ширина 1200px)
        resized := imaging.Resize(img, 1200, 0, imaging.Lanczos)

        // Сохраняем обработанное изображение
        err = imaging.Save(resized, filePath)
        if err != nil {
            return "", err
        }

        return fileName, nil
    }

    // Добавляем эндпоинт для загрузки изображений
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
                IsMain:     isMain,
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
    Name              string  `json:"name"`
    Capacity          int     `json:"capacity"`
    PricePerNight     *float64 `json:"price_per_night"`
    AddressStreet     string  `json:"address_street"`
    AddressCity       string  `json:"address_city"`
    AddressState      string  `json:"address_state"`
    AddressCountry    string  `json:"address_country"`
    AddressPostalCode string  `json:"address_postal_code"`
    AccommodationType string  `json:"accommodation_type"`
    IsShared          bool    `json:"is_shared"`
    TotalBeds         *int     `json:"total_beds,omitempty"`
    AvailableBeds     *int     `json:"available_beds,omitempty"`
    HasPrivateBathroom bool   `json:"has_private_bathroom"`
}

type Bed struct {
    ID            int     `json:"id"`
    RoomID        int     `json:"room_id"`
    BedNumber     string  `json:"bed_number"`  // изменено с int на string
    IsAvailable   bool    `json:"is_available"`
    PricePerNight float64 `json:"price_per_night"`
}

type BedBooking struct {
    ID        int       `json:"id"`
    BedID     int       `json:"bed_id"`
    UserID    int       `json:"user_id"`
    StartDate string    `json:"start_date"`
    EndDate   string    `json:"end_date"`
}
app.Post("/rooms", func(c *fiber.Ctx) error {
    var room Room
    if err := c.BodyParser(&room); err != nil {
        return c.Status(400).SendString("Неверный формат данных")
    }

    var roomId int
    err := pool.QueryRow(context.Background(), `
        INSERT INTO rooms (
            name, capacity, price_per_night,
            address_street, address_city, address_state,
            address_country, address_postal_code,
            accommodation_type, is_shared,
            total_beds, available_beds, has_private_bathroom
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
        RETURNING id`,
        room.Name, room.Capacity, room.PricePerNight,
        room.AddressStreet, room.AddressCity, room.AddressState,
        room.AddressCountry, room.AddressPostalCode,
        room.AccommodationType, room.IsShared,
        room.TotalBeds, room.AvailableBeds, room.HasPrivateBathroom,
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
        BedNumber    string  `json:"bed_number"`
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
        "id": bedID,
        "room_id": roomID,
        "bed_number": bedReq.BedNumber,
        "price_per_night": bedReq.PricePerNight,
        "is_available": true    })
})

// Получение доступных кроватей
app.Get("/rooms/:id/available-beds", func(c *fiber.Ctx) error {
    roomID := c.Params("id")
    startDate := c.Query("start_date")
    endDate := c.Query("end_date")

    query := `
        SELECT b.id, b.bed_number, b.price_per_night
        FROM beds b
        WHERE b.room_id = $1
        AND b.is_available = true
        AND b.id NOT IN (
            SELECT bed_id
            FROM bed_bookings
            WHERE $2 < end_date AND $3 > start_date
            AND status != 'cancelled'
        )
    `

    rows, err := pool.Query(context.Background(), query, roomID, startDate, endDate)
    if err != nil {
        return c.Status(500).SendString("Ошибка получения списка кроватей")
    }
    defer rows.Close()

    var beds []Bed
    for rows.Next() {
        var bed Bed
        if err := rows.Scan(&bed.ID, &bed.BedNumber, &bed.PricePerNight); err != nil {
            continue
        }
        beds = append(beds, bed)
    }

    return c.JSON(beds)
})
	// В endpoint получения списка комнат обновляем запрос
	app.Get("/rooms", func(c *fiber.Ctx) error {
		capacity := c.Query("capacity")
		startDate := c.Query("start_date")
		endDate := c.Query("end_date")
		minPrice := c.Query("min_price")
		maxPrice := c.Query("max_price")
		city := c.Query("city")      // Добавляем параметр города
		country := c.Query("country") // Добавляем параметр страны
	
		query := `
        SELECT id, name, capacity, price_per_night, address_street, address_city, 
       		address_state, address_country, address_postal_code,
       		accommodation_type, is_shared, COALESCE(total_beds, 0) AS total_beds, 
       		COALESCE(available_beds, 0) AS available_beds, 
       		has_private_bathroom, created_at 
		FROM rooms`
		args := []interface{}{}
		conditions := []string{}
	
		// Фильтр по вместимости
		if capacity != "" {
			conditions = append(conditions, "capacity >= $"+strconv.Itoa(len(args)+1))
			args = append(args, capacity)
		}
	
		// Фильтр по минимальной цене
		if minPrice != "" {
			minPriceFloat, err := strconv.ParseFloat(minPrice, 64)
			if err != nil {
				return c.Status(400).SendString("Некорректное значение минимальной цены")
			}
			conditions = append(conditions, "price_per_night >= $"+strconv.Itoa(len(args)+1))
			args = append(args, minPriceFloat)
		}
	
		// Фильтр по максимальной цене
		if maxPrice != "" {
			maxPriceFloat, err := strconv.ParseFloat(maxPrice, 64)
			if err != nil {
				return c.Status(400).SendString("Некорректное значение максимальной цены")
			}
			conditions = append(conditions, "price_per_night <= $"+strconv.Itoa(len(args)+1))
			args = append(args, maxPriceFloat)
		}
	
		// Проверка корректности диапазона цен
		if minPrice != "" && maxPrice != "" {
			minPriceFloat, _ := strconv.ParseFloat(minPrice, 64)
			maxPriceFloat, _ := strconv.ParseFloat(maxPrice, 64)
			if minPriceFloat > maxPriceFloat {
				return c.Status(400).SendString("Минимальная цена не может быть больше максимальной")
			}
		}
	
		// Фильтр по доступности дат
		if startDate != "" && endDate != "" {
			conditions = append(conditions, `
				id NOT IN (
					SELECT room_id FROM bookings 
					WHERE $`+strconv.Itoa(len(args)+1)+` < end_date 
					  AND $`+strconv.Itoa(len(args)+2)+` > start_date
				)
			`)
			args = append(args, startDate, endDate)
		}
	    if city != "" {
			conditions = append(conditions, "address_city ILIKE $"+strconv.Itoa(len(args)+1))
			args = append(args, "%"+city+"%")
		}
		if country != "" {
			conditions = append(conditions, "address_country ILIKE $"+strconv.Itoa(len(args)+1))
			args = append(args, "%"+country+"%")
		}




		// Добавление условий в запрос
		if len(conditions) > 0 {
			query += " WHERE " + strings.Join(conditions, " AND ")
		}
	
		// Логирование запроса
		log.Printf("SQL Query: %s, Args: %v", query, args)
	
		// Выполнение запроса
		rows, err := pool.Query(context.Background(), query, args...)
		if err != nil {
			log.Printf("Ошибка выполнения запроса: %v", err)
			return c.Status(500).SendString("Ошибка получения списка комнат")
		}
		defer rows.Close()
	
		// Обработка результатов
    var rooms []map[string]interface{}
    for rows.Next() {
        var id, capacity, totalBeds, availableBeds int
        var name, addressStreet, addressCity, addressState, 
            addressCountry, addressPostalCode, accommodationType string
        var pricePerNight float64
        var isShared, hasPrivateBathroom bool
        var createdAt time.Time
        
        if err := rows.Scan(
            &id, &name, &capacity, &pricePerNight,
            &addressStreet, &addressCity, &addressState,
            &addressCountry, &addressPostalCode,
            &accommodationType, &isShared,
            &totalBeds, &availableBeds, &hasPrivateBathroom,
            &createdAt,
        ); err != nil {
            log.Printf("Ошибка сканирования строки: %v", err)
            return c.Status(500).SendString("Ошибка обработки данных")
        }
        
        rooms = append(rooms, map[string]interface{}{
            "id":                 id,
            "name":               name,
            "capacity":           capacity,
            "price_per_night":    pricePerNight,
            "address_street":     addressStreet,
            "address_city":       addressCity,
            "address_state":      addressState,
            "address_country":    addressCountry,
            "address_postal_code": addressPostalCode,
            "accommodation_type":  accommodationType,
            "is_shared":          isShared,
            "total_beds":         totalBeds,
            "available_beds":     availableBeds,
            "has_private_bathroom": hasPrivateBathroom,
            "created_at":         createdAt.Format("2006-01-02 15:04:05"),
        })
    }

    return c.JSON(rooms)
})
	
	
	// Добавление бронирования
	app.Post("/bookings", func(c *fiber.Ctx) error {
		type Booking struct {
			UserID    int    `json:"user_id"`
			RoomID    int    `json:"room_id"`
			StartDate string `json:"start_date"`
			EndDate   string `json:"end_date"`
		}
	
		var booking Booking
		if err := c.BodyParser(&booking); err != nil {
			return c.Status(400).SendString("Неверный формат данных")
		}
	
		// Проверяем, существует ли пользователь
		var userCount int
		if err := pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM users WHERE id = $1", booking.UserID).Scan(&userCount); err != nil || userCount == 0 {
			return c.Status(400).SendString("Пользователь не найден")
		}
	
		// Проверяем, существует ли комната
		var roomCount int
		if err := pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM rooms WHERE id = $1", booking.RoomID).Scan(&roomCount); err != nil || roomCount == 0 {
			return c.Status(400).SendString("Комната не найдена")
		}
	
		// Проверяем доступность комнаты
		var count int
		err := pool.QueryRow(context.Background(), `
			SELECT COUNT(*) FROM bookings 
			WHERE room_id = $1 AND $2 < end_date AND $3 > start_date
		`, booking.RoomID, booking.StartDate, booking.EndDate).Scan(&count)
		if err != nil {
			log.Printf("Ошибка проверки доступности комнаты: %v", err)
			return c.Status(500).SendString("Ошибка проверки доступности комнаты")
		}
		if count > 0 {
			return c.Status(400).SendString("Комната занята на указанные даты")
		}
	
		_, err = pool.Exec(context.Background(), `
			INSERT INTO bookings (user_id, room_id, start_date, end_date) 
			VALUES ($1, $2, $3, $4)
		`, booking.UserID, booking.RoomID, booking.StartDate, booking.EndDate)
		if err != nil {
			log.Printf("Ошибка добавления бронирования: %v", err)
			return c.Status(500).SendString("Ошибка добавления бронирования")
		}
	
		return c.SendString("Бронирование добавлено успешно")
	})
	

	// Получение списка всех бронирований
	app.Get("/bookings", func(c *fiber.Ctx) error {
		query := "SELECT id, user_id, room_id, start_date, end_date FROM bookings"
		rows, err := pool.Query(context.Background(), query)
		if err != nil {
			log.Printf("Ошибка выполнения запроса: %v", err)
			return c.Status(500).SendString("Ошибка получения списка бронирований")
		}
		defer rows.Close()
	
		var bookings []map[string]interface{}
		for rows.Next() {
			var id, userID, roomID int
			var startDate, endDate time.Time
			if err := rows.Scan(&id, &userID, &roomID, &startDate, &endDate); err != nil {
				log.Printf("Ошибка сканирования строки: %v", err)
				return c.Status(500).SendString("Ошибка обработки данных бронирования")
			}
			bookings = append(bookings, map[string]interface{}{
				"id":         id,
				"user_id":    userID,
				"room_id":    roomID,
				"start_date": startDate.Format("2006-01-02"),
				"end_date":   endDate.Format("2006-01-02"),
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

	// Добавление бронирования
	// (оставлено без изменений для краткости)

	// Запуск приложения
	log.Fatal(app.Listen(":3000"))
}
