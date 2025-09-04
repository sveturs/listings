package main

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	// Секретный ключ из переменной окружения или дефолтный
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-here"
	}

	// Создаем claims для пользователя 7 (владелец витрин)
	claims := jwt.MapClaims{
		"user_id":  7,
		"email":    "testuser7@example.com", // Test user
		"is_admin": false,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	// Создаем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating token: %v\n", err)
		os.Exit(1)
	}

	// Выводим токен
	fmt.Print(tokenString)
}
