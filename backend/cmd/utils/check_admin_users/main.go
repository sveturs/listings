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
	ID   int    `json:"id"`
	Role string `json:"role"`
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

	// Connect to database
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// First check table structure
	fmt.Println("=== Users table structure ===")
	rows, err := pool.Query(ctx, `
		SELECT column_name, data_type 
		FROM information_schema.columns 
		WHERE table_name = 'users' 
		ORDER BY ordinal_position
	`)
	if err != nil {
		log.Printf("Failed to get table structure: %v", err)
		return
	}

	for rows.Next() {
		var colName, dataType string
		if err := rows.Scan(&colName, &dataType); err != nil {
			log.Printf("Warning: failed to scan row: %v", err)
			continue
		}
		fmt.Printf("- %s: %s\n", colName, dataType)
	}
	rows.Close()

	// Check admin users
	fmt.Println("\n=== Admin users in database ===")
	rows, err = pool.Query(ctx, `
		SELECT id, name, email, is_admin 
		FROM users 
		WHERE is_admin = true
		ORDER BY id
	`)
	if err != nil {
		log.Fatalf("Failed to query admin users: %v", err)
	}
	defer rows.Close()

	var adminFound bool
	var firstAdminID int
	for rows.Next() {
		var id int
		var name, email string
		var isAdmin bool
		if err := rows.Scan(&id, &name, &email, &isAdmin); err != nil {
			log.Printf("Failed to scan row: %v", err)
			continue
		}
		fmt.Printf("ID: %d, Name: %s, Email: %s, IsAdmin: %v\n", id, name, email, isAdmin)
		if !adminFound {
			adminFound = true
			firstAdminID = id
		}
	}

	if !adminFound {
		fmt.Println("No admin users found!")
		return
	}

	// Generate JWT for the first admin
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	// Create claims for admin user
	claims := &Claims{
		ID:   firstAdminID,
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

	fmt.Printf("\n=== Generated JWT token for admin user (ID: %d) ===\n", claims.ID)
	fmt.Println(tokenString)
	fmt.Printf("\nExpires at: %s\n", claims.ExpiresAt.Format(time.RFC3339))
}
