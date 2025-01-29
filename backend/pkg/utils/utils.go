//backend/pkg/utils/utils.go
package utils

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GenerateSessionToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func ErrorResponse(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{
		"error": message,
	})
}

func SuccessResponse(c *fiber.Ctx, data interface{}) error {
    return c.JSON(fiber.Map{
        "success": true,
        "data":    data,
    })
}
func StringToInt(str string, defaultValue int) int {
    if str == "" {
        return defaultValue
    }
    val, err := strconv.Atoi(str)
    if err != nil {
        return defaultValue
    }
    return val
}