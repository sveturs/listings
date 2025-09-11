package main

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	// Получаем секретный ключ из переменной окружения или используем дефолтный
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production"
	}

	// Создаем claims для админа
	claims := jwt.MapClaims{
		"user_id":  1,
		"email":    "admin@example.com",
		"is_admin": true,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	// Создаем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating token: %v\n", err)
		os.Exit(1)
	}

	// Выводим токен
	fmt.Println(tokenString)
}