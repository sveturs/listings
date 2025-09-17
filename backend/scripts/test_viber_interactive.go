package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"svetu/internal/proj/viber/infobip"
)

func main() {
	// Загружаем .env
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("INFOBIP_API_KEY")
	baseURL := os.Getenv("INFOBIP_BASE_URL")
	sender := os.Getenv("INFOBIP_SENDER_ID")

	if apiKey == "" || baseURL == "" {
		log.Fatal("INFOBIP_API_KEY and INFOBIP_BASE_URL must be set")
	}

	// Создаём клиент Infobip
	client := infobip.NewClient(apiKey, baseURL)

	ctx := context.Background()

	// Телефон получателя
	to := "381604485063" // Ваш номер телефона

	fmt.Println("Отправляем интерактивное сообщение с ссылками...")

	// 1. Отправляем интерактивное меню
	resp, err := client.SendInteractiveMessage(ctx, sender, to, nil)
	if err != nil {
		log.Fatalf("Failed to send interactive message: %v", err)
	}

	fmt.Printf("Интерактивное меню отправлено: %+v\n", resp)

	// 2. Отправляем ссылку на трекинг
	orderID := "TEST-ORDER-123"
	resp2, err := client.SendTrackingLink(ctx, sender, to, orderID, nil)
	if err != nil {
		log.Fatalf("Failed to send tracking link: %v", err)
	}

	fmt.Printf("Ссылка на трекинг отправлена: %+v\n", resp2)

	fmt.Println("\n✅ Проверьте Viber! Сообщения должны прийти с кликабельными ссылками, которые открываются прямо в приложении.")
}