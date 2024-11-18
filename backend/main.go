package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
		"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3001", // Укажите фронтенд URL
		AllowMethods: "GET,POST,DELETE,PUT",
	}))
	// Подключение к базе данных
	dbURL := os.Getenv("DATABASE_URL")
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer pool.Close()

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

	// Добавление комнаты
	app.Post("/rooms", func(c *fiber.Ctx) error {
		type Room struct {
			Name          string  `json:"name"`
			Capacity      int     `json:"capacity"`
			PricePerNight float64 `json:"price_per_night"`
		}

		var room Room
		if err := c.BodyParser(&room); err != nil {
			return c.Status(400).SendString("Неверный формат данных")
		}

		_, err := pool.Exec(context.Background(), "INSERT INTO rooms (name, capacity, price_per_night) VALUES ($1, $2, $3)", room.Name, room.Capacity, room.PricePerNight)
		if err != nil {
			log.Printf("Ошибка добавления комнаты: %v", err)
			return c.Status(500).SendString("Ошибка добавления комнаты")
		}

		return c.SendString("Комната добавлена успешно")
	})

	// Получение списка комнат с фильтрами
	app.Get("/rooms", func(c *fiber.Ctx) error {
		capacity := c.Query("capacity")
		startDate := c.Query("start_date")
		endDate := c.Query("end_date")
		minPrice := c.Query("min_price")
		maxPrice := c.Query("max_price")

		query := "SELECT id, name, capacity, price_per_night, created_at FROM rooms"
		args := []interface{}{}
		conditions := []string{}

		if capacity != "" {
			conditions = append(conditions, "capacity >= $"+strconv.Itoa(len(args)+1))
			args = append(args, capacity)
		}

		if minPrice != "" {
			conditions = append(conditions, "price_per_night >= $"+strconv.Itoa(len(args)+1))
			args = append(args, minPrice)
		}

		if maxPrice != "" {
			conditions = append(conditions, "price_per_night <= $"+strconv.Itoa(len(args)+1))
			args = append(args, maxPrice)
		}

		if minPrice != "" && maxPrice != "" {
			min, err1 := strconv.ParseFloat(minPrice, 64)
			max, err2 := strconv.ParseFloat(maxPrice, 64)
			if err1 != nil || err2 != nil || min > max {
				return c.Status(400).SendString("Некорректный диапазон цен")
			}
		}

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

		if len(conditions) > 0 {
			query += " WHERE " + strings.Join(conditions, " AND ")
		}

		rows, err := pool.Query(context.Background(), query, args...)
		if err != nil {
			log.Printf("Ошибка выполнения запроса: %v", err)
			return c.Status(500).SendString("Ошибка получения списка комнат")
		}
		defer rows.Close()

		var rooms []map[string]interface{}
		for rows.Next() {
			var id, capacity int
			var name string
			var pricePerNight float64
			var createdAt time.Time
			if err := rows.Scan(&id, &name, &capacity, &pricePerNight, &createdAt); err != nil {
				log.Printf("Ошибка сканирования строки: %v", err)
				return c.Status(500).SendString("Ошибка обработки данных")
			}
			rooms = append(rooms, map[string]interface{}{
				"id":              id,
				"name":            name,
				"capacity":        capacity,
				"price_per_night": pricePerNight,
				"created_at":      createdAt.Format("2006-01-02 15:04:05"),
			})
		}

		return c.JSON(rooms)
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

		var pricePerNight float64
		err = pool.QueryRow(context.Background(), "SELECT price_per_night FROM rooms WHERE id=$1 AND price_per_night IS NOT NULL", booking.RoomID).Scan(&pricePerNight)
		if err != nil {
			log.Printf("Ошибка получения стоимости комнаты: %v", err)
			return c.Status(400).SendString("Комната не найдена или цена не задана")
		}

		layout := "2006-01-02"
		startDate, _ := time.Parse(layout, booking.StartDate)
		endDate, _ := time.Parse(layout, booking.EndDate)
		if startDate.After(endDate) || startDate.Equal(endDate) {
			return c.Status(400).SendString("Некорректный диапазон дат")
		}
		totalDays := int(endDate.Sub(startDate).Hours() / 24)
		totalCost := pricePerNight * float64(totalDays)

		_, err = pool.Exec(context.Background(), "INSERT INTO bookings (user_id, room_id, start_date, end_date) VALUES ($1, $2, $3, $4)",
			booking.UserID, booking.RoomID, booking.StartDate, booking.EndDate)
		if err != nil {
			log.Printf("Ошибка добавления бронирования: %v", err)
			return c.Status(500).SendString("Ошибка добавления бронирования")
		}

		return c.JSON(fiber.Map{
			"message":    "Бронирование добавлено успешно",
			"total_cost": totalCost,
		})
	})

	// Удаление бронирования
	app.Delete("/bookings/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		result, err := pool.Exec(context.Background(), "DELETE FROM bookings WHERE id=$1", id)
		if err != nil {
			log.Printf("Ошибка удаления бронирования: %v", err)
			return c.Status(500).SendString("Ошибка удаления бронирования")
		 }
		 if result.RowsAffected() == 0 {
			 return c.Status(404).SendString("Бронирование не найдено")
		 }
		return c.SendString("Бронирование успешно удалено")
	})

	// Запуск приложения
	app.Listen(":3000")
}
