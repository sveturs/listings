// Package utils
// backend/pkg/utils/utils.go
package utils

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"time"

	"backend/pkg/jwt"

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

// GenerateJWTToken генерирует JWT токен для пользователя с дефолтным временем жизни
func GenerateJWTToken(userID int, email string, secret string) (string, error) {
	return jwt.GenerateToken(userID, email, secret)
}

// GenerateJWTTokenWithDuration генерирует JWT токен для пользователя с заданным временем жизни
func GenerateJWTTokenWithDuration(userID int, email string, secret string, duration time.Duration) (string, error) {
	return jwt.GenerateTokenWithDuration(userID, email, secret, duration)
}

// ValidateJWTToken валидирует JWT токен
func ValidateJWTToken(tokenString string, secret string) (*jwt.Claims, error) {
	return jwt.ValidateToken(tokenString, secret)
}
