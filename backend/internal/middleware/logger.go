package middleware

import (
    "github.com/gofiber/fiber/v2"
    "time"
    "log"
)

func (m *Middleware) Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		
		// Продолжаем обработку запроса
		err := c.Next()
		
		// Логируем информацию о запросе
		msg := "%s - %s %s - %d - %v"
		if err != nil {
			msg = "%s - %s %s - %d - %v - Error: %v"
		}

		duration := time.Since(start)
		
		if err != nil {
			log.Printf(msg,
				c.IP(),
				c.Method(),
				c.Path(),
				c.Response().StatusCode(),
				duration,
				err,
			)
		} else {
			log.Printf(msg,
				c.IP(),
				c.Method(),
				c.Path(),
				c.Response().StatusCode(),
				duration,
			)
		}

		return err
	}
}