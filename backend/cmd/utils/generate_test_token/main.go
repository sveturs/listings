package main

import (
	"backend/internal/config"
	"backend/internal/proj/global/service"
	"backend/internal/storage/postgres"
	"backend/internal/types"
	"backend/pkg/utils"
	"context"
	"fmt"
	"log"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect to database
	db, err := postgres.New(cfg.Database.ConnectionString)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Create global service
	globalService := service.NewService(cfg, db)

	// Check if admin user exists
	ctx := context.Background()
	adminEmail := "bevzenko.sergey@gmail.com"

	// Get user by email
	user, err := db.GetUserByEmail(ctx, adminEmail)
	if err != nil {
		log.Fatal("Failed to get user:", err)
	}

	if user == nil {
		log.Fatal("User not found")
	}

	// Generate session token
	sessionToken := utils.GenerateSessionToken()

	// Create session data
	sessionData := &types.SessionData{
		UserID:     user.ID,
		Name:       user.Name,
		Email:      user.Email,
		GoogleID:   user.GoogleID,
		PictureURL: user.PictureURL,
		Provider:   "google",
	}

	// Save session
	globalService.Auth().SaveSession(sessionToken, sessionData)

	fmt.Printf("Generated test token for user %s (ID: %d)\n", user.Email, user.ID)
	fmt.Printf("Token: %s\n", sessionToken)
	fmt.Println("\nHow to use:")
	fmt.Printf("1. Set cookie: session_token=%s\n", sessionToken)
	fmt.Printf("2. Or add to localStorage: localStorage.setItem('user_session_token', '%s')\n", sessionToken)
	fmt.Printf("3. Or add query parameter: ?session_token=%s\n", sessionToken)
}