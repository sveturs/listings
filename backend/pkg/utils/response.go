package utils

import (
	"github.com/gofiber/fiber/v2"
)

// SendSuccess sends a successful response with custom status code
func SendSuccess(c *fiber.Ctx, statusCode int, message string, data interface{}) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"success": true,
		"message": message,
		"data":    data,
	})
}

// SendError sends an error response with custom status code
func SendError(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"success": false,
		"error":   message,
	})
}
