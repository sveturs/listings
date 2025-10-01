package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"backend/internal/config"

	"github.com/rs/zerolog"
	authClient "github.com/sveturs/auth/pkg/http/client"
	authEntity "github.com/sveturs/auth/pkg/http/entity"
	authService "github.com/sveturs/auth/pkg/http/service"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Инициализируем логгер
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	fmt.Println("🔐 Генерация токена через Auth Service\n")
	fmt.Println("⚠️  ВАЖНО: Этот скрипт требует email и пароль пользователя.")
	fmt.Println("    Убедитесь, что пользователь существует в Auth Service.")
	fmt.Println()

	var email, password string

	// Проверяем env переменные
	email = os.Getenv("ADMIN_EMAIL")
	password = os.Getenv("ADMIN_PASSWORD")

	// Если не заданы в env, запрашиваем у пользователя
	if email == "" {
		fmt.Print("📧 Введите email пользователя: ")
		fmt.Scanln(&email)
	}

	if password == "" {
		fmt.Print("🔑 Введите пароль: ")
		fmt.Scanln(&password)
	}

	if email == "" || password == "" {
		log.Fatal("❌ Email и пароль обязательны!")
	}

	// Создаем Auth Service клиент
	client, err := authClient.NewClientWithResponses(cfg.AuthServiceURL)
	if err != nil {
		log.Fatalf("❌ Не удалось создать Auth Service клиент: %v", err)
	}

	authSvc := authService.NewAuthService(client, logger)

	// Выполняем логин
	ctx := context.Background()
	loginReq := authEntity.UserLoginRequest{
		Email:      email,
		Password:   password,
		DeviceID:   "generate_token_script",
		DeviceName: "Token Generator Script",
	}

	fmt.Printf("\n🔄 Выполняем логин через Auth Service (%s)...\n", cfg.AuthServiceURL)

	resp, err := authSvc.Login(ctx, loginReq)
	if err != nil {
		log.Fatalf("❌ Ошибка при логине: %v", err)
	}

	// Проверяем статус
	if resp.StatusCode() != 200 {
		log.Fatalf("❌ Ошибка логина: статус %d", resp.StatusCode())
	}

	if resp.JSON200 == nil || resp.JSON200.AccessToken == nil || *resp.JSON200.AccessToken == "" {
		log.Fatal("❌ Не получен access token от Auth Service")
	}

	accessToken := *resp.JSON200.AccessToken
	refreshToken := ""
	if resp.JSON200.RefreshToken != nil && *resp.JSON200.RefreshToken != "" {
		refreshToken = *resp.JSON200.RefreshToken
	}

	// Получаем информацию о пользователе
	var userName string
	var isAdmin bool

	if resp.JSON200.User != nil {
		if resp.JSON200.User.Name != nil {
			userName = *resp.JSON200.User.Name
		}
		if resp.JSON200.User.IsAdmin != nil {
			isAdmin = *resp.JSON200.User.IsAdmin
		}
	}

	// Выводим результат
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Printf("🎉 ТОКЕНЫ УСПЕШНО ПОЛУЧЕНЫ!\n")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("\n👤 Пользователь: %s\n", email)
	if userName != "" {
		fmt.Printf("📧 Имя: %s\n", userName)
	}
	fmt.Printf("🔑 Admin права: %v\n", isAdmin)

	fmt.Println("\n🔑 ACCESS TOKEN (JWT):")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("%s\n", accessToken)

	if refreshToken != "" {
		fmt.Println("\n🔄 REFRESH TOKEN:")
		fmt.Println(strings.Repeat("-", 80))
		fmt.Printf("%s\n", refreshToken)
	}

	fmt.Println("\n📝 СПОСОБЫ ИСПОЛЬЗОВАНИЯ:")
	fmt.Println(strings.Repeat("-", 80))

	fmt.Println("\n1️⃣  Тест API синонимов:")
	fmt.Printf("   curl -H \"Authorization: Bearer %s\" \\\n", accessToken)
	fmt.Println("        'http://localhost:3000/api/v1/admin/search/synonyms?page=1&limit=20&language=ru'")

	fmt.Println("\n2️⃣  Добавление синонима:")
	fmt.Printf("   curl -X POST -H \"Authorization: Bearer %s\" \\\n", accessToken)
	fmt.Println("        -H \"Content-Type: application/json\" \\")
	fmt.Println("        -d '{\"word\": \"телефон\", \"synonyms\": [\"смартфон\", \"мобильный\"], " +
		"\"language\": \"ru\"}' \\")
	fmt.Println("        'http://localhost:3000/api/v1/admin/search/synonyms'")

	fmt.Println("\n3️⃣  Проверка профиля:")
	fmt.Printf("   curl -H \"Authorization: Bearer %s\" \\\n", accessToken)
	fmt.Println("        'http://localhost:3000/api/v1/auth/me'")

	fmt.Println("\n4️⃣  Сохранить в переменную для последующего использования:")
	fmt.Printf("   export TOKEN='%s'\n", accessToken)
	fmt.Println("   curl -H \"Authorization: Bearer $TOKEN\" http://localhost:3000/api/v1/users/me")

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("⏰ Access токен действителен: обычно 15 минут")
	fmt.Println("🔄 Refresh токен действителен: обычно 30 дней")
	fmt.Println("🔒 Алгоритм: RS256 (Auth Service)")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()
}
