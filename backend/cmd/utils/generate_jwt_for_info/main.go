package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"backend/internal/config"
	"backend/internal/proj/users/service"
	"backend/internal/storage/filestorage"
	"backend/internal/storage/postgres"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Создаем файловое хранилище
	fileStorage, err := filestorage.NewFileStorage(context.Background(), cfg.FileStorage)
	if err != nil {
		log.Fatal("Failed to create file storage:", err)
	}

	// Подключаемся к базе данных
	db, err := postgres.NewDatabase(context.Background(), cfg.DatabaseURL, nil, "", fileStorage, cfg.SearchWeights)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Создаем сервис авторизации
	authService := service.NewAuthService(
		cfg.GoogleClientID,
		cfg.GoogleClientSecret,
		cfg.GoogleRedirectURL,
		db,
		cfg.JWTSecret,
		cfg.JWTExpirationHours,
	)

	// Целевой email
	ctx := context.Background()
	targetEmail := "info@svetu.rs"

	// Получаем пользователя
	user, err := db.GetUserByEmail(ctx, targetEmail)
	if err != nil || user == nil {
		log.Fatal("User not found:", err)
	}

	// Генерируем JWT токен
	jwtToken, err := authService.GenerateJWT(user.ID, user.Email)
	if err != nil {
		log.Fatal("Failed to generate JWT token:", err)
	}

	fmt.Printf("Generated JWT token for %s (ID: %d)\n", user.Email, user.ID)
	fmt.Printf("Token: %s\n\n", jwtToken)

	// Выполняем GET запрос к /api/v1/auth/session
	req, err := http.NewRequest("GET", "http://localhost:3000/api/v1/auth/session", nil)
	if err != nil {
		log.Fatal("Failed to create request:", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwtToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to execute request:", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Failed to read response:", err)
	}

	fmt.Printf("Response status: %d\n", resp.StatusCode)
	fmt.Printf("Response body: %s\n", string(body))
}
