package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func main() {
	// Load .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Get database URL
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// Get JWT secret
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	// Connect to database
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Get admin user info
	var id int
	var email string
	err = pool.QueryRow(ctx, `
		SELECT id, email 
		FROM users 
		WHERE is_admin = true
		ORDER BY id
		LIMIT 1
	`).Scan(&id, &email)
	if err != nil {
		log.Printf("Failed to get admin user: %v", err)
		return
	}

	fmt.Printf("Found admin user: ID=%d, Email=%s\n", id, email)

	// Create claims with proper structure
	now := time.Now()
	claims := &Claims{
		UserID: id,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour * 30)), // 30 days
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "svetu-backend",
			Subject:   fmt.Sprintf("user:%d", id),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Printf("Failed to sign token: %v", err)
		return
	}

	fmt.Printf("\n=== Generated JWT token ===\n")
	fmt.Println(tokenString)
	fmt.Printf("\nExpires at: %s\n", claims.ExpiresAt.Format(time.RFC3339))
}
