package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"backend/internal/config"
	"backend/internal/storage/filestorage"
	"backend/internal/storage/postgres"

	"github.com/golang-jwt/jwt/v5"
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

	// Целевой email администратора
	ctx := context.Background()
	targetEmail := "admin@test.com"

	// Получаем пользователя
	user, err := db.GetUserByEmail(ctx, targetEmail)
	if err != nil || user == nil {
		log.Fatal("Admin user not found:", err)
	}

	// Получаем профиль пользователя с полем is_admin
	userProfile, err := db.GetUserProfile(ctx, user.ID)
	if err != nil {
		log.Fatal("Failed to get user profile:", err)
	}

	// Проверяем, что пользователь является администратором
	if !userProfile.IsAdmin {
		log.Fatal("User is not an admin")
	}

	// Генерируем JWT токен
	jwtToken, err := generateJWT(int64(user.ID), user.Email, cfg.JWTSecret, cfg.JWTExpirationHours)
	if err != nil {
		log.Fatal("Failed to generate JWT token:", err)
	}

	// Выводим результат
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Printf("🎉 ADMIN JWT ТОКЕН УСПЕШНО СОЗДАН!\n")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("\n👤 Администратор: %s (ID: %d)\n", user.Email, user.ID)
	fmt.Printf("📧 Имя: %s\n", user.Name)
	fmt.Printf("🔑 Admin права: %v\n", userProfile.IsAdmin)

	fmt.Println("\n🔑 ACCESS TOKEN (JWT):")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("%s\n", jwtToken)

	fmt.Println("\n📝 СПОСОБЫ ИСПОЛЬЗОВАНИЯ:")
	fmt.Println(strings.Repeat("-", 80))

	fmt.Println("\n1️⃣  Тест API синонимов:")
	fmt.Printf("   curl -H \"Authorization: Bearer %s\" \\\n", jwtToken)
	fmt.Println("        'http://localhost:3000/api/v1/admin/search/synonyms?page=1&limit=20&language=ru'")

	fmt.Println("\n2️⃣  Добавление синонима:")
	fmt.Printf("   curl -X POST -H \"Authorization: Bearer %s\" \\\n", jwtToken)
	fmt.Println("        -H \"Content-Type: application/json\" \\")
	fmt.Println("        -d '{\"word\": \"телефон\", \"synonyms\": [\"смартфон\", \"мобильный\"], \"language\": \"ru\"}' \\")
	fmt.Println("        'http://localhost:3000/api/v1/admin/search/synonyms'")

	fmt.Println("\n3️⃣  Для использования в браузере (установка в localStorage):")
	fmt.Printf("   localStorage.setItem('access_token', '%s');\n", jwtToken)

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Printf("⏰ Access токен действителен: %d часов\n", cfg.JWTExpirationHours)
	fmt.Println("🔒 Тип авторизации: JWT Bearer")
	fmt.Println(strings.Repeat("=", 80))
}

func generateJWT(userID int64, email string, secret string, expirationHours int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"iss":     "svetu-backend",
		"sub":     fmt.Sprintf("user:%d", userID),
		"exp":     time.Now().Add(time.Hour * time.Duration(expirationHours)).Unix(),
		"nbf":     time.Now().Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
