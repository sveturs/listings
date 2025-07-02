package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"backend/internal/storage/filestorage"
	"backend/internal/storage/opensearch"
	"backend/internal/storage/postgres"
	"backend/internal/types"
	"backend/pkg/utils"
)

func main() {
	// Simple hardcoded config
	cfg := &struct {
		Database    struct{ ConnectionString string }
		FileStorage filestorage.Config
		OpenSearch  opensearch.Config
	}{
		Database: struct{ ConnectionString string }{
			ConnectionString: "postgres://postgres:password@localhost:5432/hostel_db?sslmode=disable",
		},
		FileStorage: filestorage.Config{
			Provider: "minio",
			MinIO: &filestorage.MinIOConfig{
				Endpoint:      "localhost:9000",
				AccessKey:     "minioadmin",
				SecretKey:     "1321321321321",
				BucketName:    "listings",
				Region:        "eu-central-1",
				UseSSL:        false,
				PublicBaseURL: "http://localhost:3000",
			},
		},
		OpenSearch: opensearch.Config{
			URL:      "http://localhost:9200",
			Username: "admin",
			Password: "admin",
			Index:    "marketplace",
		},
	}

	// Connect to database
	osClient, err := opensearch.NewClient(cfg.OpenSearch)
	if err != nil {
		log.Fatal("Failed to connect to OpenSearch:", err)
	}

	// Create file storage
	fileStorage, err := filestorage.NewFactory(cfg.FileStorage).CreateStorage()
	if err != nil {
		log.Fatal("Failed to create file storage:", err)
	}

	db, err := postgres.NewDatabase(cfg.Database.ConnectionString, osClient, "marketplace", fileStorage)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close(context.Background())

	// Get admin user
	ctx := context.Background()
	adminEmail := "voroshilovdo@gmail.com" // Используем email администратора из логов

	user, err := db.GetUserByEmail(ctx, adminEmail)
	if err != nil || user == nil {
		log.Fatal("Failed to get user or user not found:", err)
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

	// Save session to database
	sessionJSON, _ := json.Marshal(sessionData)
	expiry := time.Now().Add(24 * time.Hour)

	if err := db.SaveSession(ctx, sessionToken, string(sessionJSON), expiry); err != nil {
		log.Fatal("Failed to save session:", err)
	}

	fmt.Printf("Generated test token for user %s (ID: %d)\n", user.Email, user.ID)
	fmt.Printf("Token: %s\n", sessionToken)
	fmt.Println("\nHow to use with curl:")
	fmt.Printf("curl -H \"Cookie: session_token=%s\" http://localhost:3000/api/admin/attribute-groups\n", sessionToken)
}
