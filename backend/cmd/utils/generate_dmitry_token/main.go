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

	// Target user email
	ctx := context.Background()
	targetEmail := "voroshilovdo@gmail.com"

	// Get user by email
	user, err := db.GetUserByEmail(ctx, targetEmail)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		// Если пользователь не найден, попробуем создать его
		fmt.Println("\nПользователь не найден. Создаем нового пользователя...")
		
		// Создаем нового пользователя
		newUser := &types.User{
			Name:       "Dmitry Voroshilov",
			Email:      targetEmail,
			GoogleID:   "google_" + utils.GenerateSessionToken()[:20], // Фиктивный Google ID
			PictureURL: "https://lh3.googleusercontent.com/a/default-user=s96-c", // Дефолтное изображение
			Provider:   "google",
		}
		
		// Сохраняем пользователя в базе данных
		userID, err := db.CreateUser(ctx, newUser)
		if err != nil {
			log.Fatal("Failed to create user:", err)
		}
		
		// Получаем созданного пользователя
		user, err = db.GetUserByID(ctx, userID)
		if err != nil {
			log.Fatal("Failed to get created user:", err)
		}
		
		fmt.Printf("Создан новый пользователь: %s (ID: %d)\n", user.Email, user.ID)
	}

	if user == nil {
		log.Fatal("User not found and could not be created")
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

	fmt.Printf("\n✅ Generated authentication token for user %s (ID: %d)\n", user.Email, user.ID)
	fmt.Printf("\n🔑 Token: %s\n", sessionToken)
	fmt.Println("\n📝 Как использовать токен:")
	fmt.Println("====================================")
	fmt.Printf("1. Cookie (для браузера):\n   document.cookie = 'session_token=%s; path=/'\n\n", sessionToken)
	fmt.Printf("2. LocalStorage (для SPA):\n   localStorage.setItem('user_session_token', '%s')\n\n", sessionToken)
	fmt.Printf("3. Query parameter:\n   http://localhost:3001/?session_token=%s\n\n", sessionToken)
	fmt.Printf("4. cURL запрос:\n   curl -H \"Cookie: session_token=%s\" http://localhost:3000/api/v1/auth/session\n\n", sessionToken)
	fmt.Println("====================================")
	fmt.Println("⏰ Токен действителен 24 часа")
}