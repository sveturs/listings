package middleware

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
)

func (m *Middleware) Logger() fiber.Handler {
	masker := NewSensitiveDataMasker()

	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Логируем запрос с маскированием sensitive данных
		path := c.Path()
		method := c.Method()

		// Маскируем query параметры если есть sensitive данные
		queryParams := c.Request().URI().QueryString()
		if len(queryParams) > 0 {
			maskedQuery := masker.Mask(string(queryParams))
			log.Printf("REQUEST: %s %s?%s", method, path, maskedQuery)
		} else {
			log.Printf("REQUEST: %s %s", method, path)
		}

		// Логируем тело запроса только для POST/PUT/PATCH с маскированием
		if method == "POST" || method == "PUT" || method == "PATCH" {
			body := c.Body()
			if len(body) > 0 && len(body) < 10000 { // Не логируем большие тела запросов
				// Пытаемся распарсить как JSON для красивого вывода
				var jsonBody interface{}
				if err := json.Unmarshal(body, &jsonBody); err == nil {
					maskedBody := masker.Mask(string(body))
					log.Printf("REQUEST BODY: %s", maskedBody)
				}
			}
		}

		err := c.Next()

		// Логируем результат
		log.Printf("RESPONSE: %s %s - %d - %v",
			method,
			path,
			c.Response().StatusCode(),
			time.Since(start),
		)

		return err
	}
}
