package main

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func main() {
	// Test admin user ID 1, email test@example.com
	userID := 1
	email := "test@example.com"

	// Get secret from environment or use default
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-here"
	}

	// Create token claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"nbf":     time.Now().Unix(),
		"iat":     time.Now().Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating token: %v\n", err)
		os.Exit(1)
	}

	// Output token only
	fmt.Print(tokenString)
}
