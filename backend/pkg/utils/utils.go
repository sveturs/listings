// Package utils
// backend/pkg/utils/utils.go
package utils

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"

	"backend/internal/logger"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// ErrorResponseSwag представляет структуру ответа об ошибке для swagger
type ErrorResponseSwag struct {
	Error string `json:"error" example:"Описание ошибки"`
}

// SuccessResponseSwag представляет структуру успешного ответа для swagger
type SuccessResponseSwag struct {
	Success bool        `json:"success" example:"true"`
	Data    interface{} `json:"data"`
}

func GenerateSessionToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func ErrorResponse(c *fiber.Ctx, status int, message string) error {
	// Логируем ошибку для отладки
	if status == 500 {
		logger.Error().
			Str("path", c.Path()).
			Str("error_message", message).
			Str("method", c.Method()).
			Msg("ErrorResponse 500 called")
	}
	return c.Status(status).JSON(fiber.Map{
		"error": message,
	})
}

// GetUserIDFromContext извлекает userID из контекста Fiber
func GetUserIDFromContext(c *fiber.Ctx) int {
	if userID := c.Locals("userID"); userID != nil {
		if uid, ok := userID.(int); ok {
			return uid
		}
	}
	return 0
}

// IsAdmin проверяет, является ли пользователь администратором
func IsAdmin(c *fiber.Ctx) bool {
	// Проверяем роль пользователя из контекста
	if role := c.Locals("role"); role != nil {
		if roleStr, ok := role.(string); ok && roleStr == "admin" {
			return true
		}
	}

	// Также проверяем флаг isAdmin
	if isAdmin := c.Locals("isAdmin"); isAdmin != nil {
		if admin, ok := isAdmin.(bool); ok && admin {
			return true
		}
	}

	return false
}

func SuccessResponse(c *fiber.Ctx, data interface{}) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

// SendErrorResponse отправляет ответ об ошибке с дополнительными данными
func SendErrorResponse(c *fiber.Ctx, status int, message string, data fiber.Map) error {
	// Логируем ошибку для отладки
	if status == 500 {
		logger.Error().
			Str("path", c.Path()).
			Str("error_message", message).
			Str("method", c.Method()).
			Any("data", data).
			Msg("SendErrorResponse 500 called")
	}

	response := fiber.Map{
		"error": message,
	}

	// Добавляем дополнительные данные если они есть
	if data != nil {
		for k, v := range data {
			response[k] = v
		}
	}

	return c.Status(status).JSON(response)
}

// SendSuccessResponse отправляет успешный ответ с данными и сообщением
func SendSuccessResponse(c *fiber.Ctx, data interface{}, message string) error {
	return c.JSON(fiber.Map{
		"success": true,
		"data":    data,
		"message": message,
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

// HashPassword хеширует пароль с использованием bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash проверяет пароль против хеша
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// JWT функции удалены - используйте Auth Service для генерации и валидации токенов
