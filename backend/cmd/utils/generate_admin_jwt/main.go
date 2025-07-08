package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type Claims struct {
	ID   int    `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func main() {
	// Load .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Get JWT secret
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	// Create claims for admin user (ID 25 based on previous sessions)
	claims := &Claims{
		ID:   25,
		Role: "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour * 30)), // 30 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Fatalf("Failed to sign token: %v", err)
	}

	fmt.Printf("Generated JWT token for admin user (ID: %d):\n", claims.ID)
	fmt.Println(tokenString)
	fmt.Printf("\nExpires at: %s\n", claims.ExpiresAt.Time.Format(time.RFC3339))
}
