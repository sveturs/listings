package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type SessionData struct {
	UserID     int    `json:"user_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	GoogleID   string `json:"google_id"`
	PictureURL string `json:"picture_url"`
	Provider   string `json:"provider"`
}

func generateSessionToken() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

func main() {
	// Connect to database
	db, err := sql.Open("postgres", "postgres://postgres:password@localhost:5432/hostel_db?sslmode=disable")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Get admin user
	adminEmail := "voroshilovdo@gmail.com"

	var userID int
	var name, email, googleID, pictureURL string

	row := db.QueryRow(`
		SELECT id, name, email, google_id, COALESCE(picture_url, '') 
		FROM users 
		WHERE email = $1
	`, adminEmail)

	if err := row.Scan(&userID, &name, &email, &googleID, &pictureURL); err != nil {
		log.Fatal("Failed to get user:", err)
	}

	// Generate session token
	sessionToken := generateSessionToken()

	// Create session data
	sessionData := &SessionData{
		UserID:     userID,
		Name:       name,
		Email:      email,
		GoogleID:   googleID,
		PictureURL: pictureURL,
		Provider:   "google",
	}

	// Save session to database
	sessionJSON, _ := json.Marshal(sessionData)
	expiry := time.Now().Add(24 * time.Hour)

	_, err = db.Exec(`
		INSERT INTO user_sessions (id, data, expiry) 
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE SET data = $2, expiry = $3
	`, sessionToken, string(sessionJSON), expiry)
	if err != nil {
		log.Fatal("Failed to save session:", err)
	}

	fmt.Printf("Generated test token for user %s (ID: %d)\n", email, userID)
	fmt.Printf("Token: %s\n", sessionToken)
	fmt.Println("\nHow to use with curl:")
	fmt.Printf("curl -H \"Cookie: session_token=%s\" http://localhost:3000/api/admin/attribute-groups\n", sessionToken)
}
