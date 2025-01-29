package middleware

import (
    "github.com/gofiber/fiber/v2"
    "time"
    "log"
)

func (m *Middleware) Logger() fiber.Handler {
    return func(c *fiber.Ctx) error {
        start := time.Now()
        err := c.Next()
        
        // Логируем только ошибки
        if err != nil || c.Response().StatusCode() >= 400 {
            log.Printf("ERROR: %s - %s %s - %d - %v - Error: %v",
                c.IP(),
                c.Method(),
                c.Path(),
                c.Response().StatusCode(),
                time.Since(start),
                err,
            )
        }
        
        return err
    }
}