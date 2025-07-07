package main

import (
	"context"
	"fmt"
	"log"
	"strings"

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
	fileStorage, err := filestorage.NewFileStorage(cfg.FileStorage)
	if err != nil {
		log.Fatal("Failed to create file storage:", err)
	}

	// Подключаемся к базе данных
	db, err := postgres.NewDatabase(cfg.DatabaseURL, nil, "", fileStorage, cfg.SearchWeights)
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
	targetEmail := "voroshilovdo@gmail.com"

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

	// Генерируем refresh токен
	refreshToken, _, err := authService.GenerateTokensForOAuth(ctx, user.ID, user.Email, "127.0.0.1", "CLI Tool")
	if err != nil {
		log.Fatal("Failed to generate refresh token:", err)
	}

	// Выводим результат
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Printf("🎉 JWT ТОКЕНЫ АВТОРИЗАЦИИ УСПЕШНО СОЗДАНЫ!\n")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("\n👤 Пользователь: %s (ID: %d)\n", user.Email, user.ID)
	fmt.Printf("📧 Имя: %s\n", user.Name)
	fmt.Printf("🖼️  Фото: %s\n", user.PictureURL)

	fmt.Println("\n🔑 ACCESS TOKEN (JWT):")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("%s\n", jwtToken)

	fmt.Println("\n🔄 REFRESH TOKEN:")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("%s\n", refreshToken)

	fmt.Println("\n📝 СПОСОБЫ ИСПОЛЬЗОВАНИЯ:")
	fmt.Println(strings.Repeat("-", 80))

	fmt.Println("\n1️⃣  Authorization header (рекомендуется):")
	fmt.Printf("   curl -H \"Authorization: Bearer %s\" \\\n", jwtToken)
	fmt.Println("        http://localhost:3000/api/v1/user/profile")

	fmt.Println("\n2️⃣  В JavaScript (axios):")
	fmt.Printf("   axios.defaults.headers.common['Authorization'] = 'Bearer %s';\n", jwtToken)

	fmt.Println("\n3️⃣  В JavaScript (fetch):")
	fmt.Println("   fetch('http://localhost:3000/api/v1/user/profile', {")
	fmt.Println("     headers: {")
	fmt.Printf("       'Authorization': 'Bearer %s'\n", jwtToken)
	fmt.Println("     }")
	fmt.Println("   })")

	fmt.Println("\n4️⃣  Проверка токена:")
	fmt.Printf("   curl -H \"Authorization: Bearer %s\" \\\n", jwtToken)
	fmt.Println("        http://localhost:3000/api/v1/auth/me")

	fmt.Println("\n5️⃣  Обновление токена (когда истечет):")
	fmt.Println("   curl -X POST http://localhost:3000/api/v1/auth/refresh \\")
	fmt.Println("        -H \"Content-Type: application/json\" \\")
	fmt.Printf("        -d '{\"refresh_token\": \"%s\"}'\n", refreshToken)

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Printf("⏰ Access токен действителен: %d часов\n", cfg.JWTExpirationHours)
	fmt.Println("⏰ Refresh токен действителен: 30 дней")
	fmt.Println("🔒 Тип авторизации: JWT Bearer")
	fmt.Println(strings.Repeat("=", 80))
}

// Вспомогательная функция для форматирования
func repeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
