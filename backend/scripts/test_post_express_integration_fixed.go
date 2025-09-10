//go:build ignore
// +build ignore

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"backend/internal/proj/postexpress/models"
	"backend/internal/proj/postexpress/service"
	"backend/pkg/logger"
)

func main() {
	// Инициализация логгера
	appLogger := logger.New()

	// Получение учетных данных из переменных окружения
	username := os.Getenv("POST_EXPRESS_USERNAME")
	password := os.Getenv("POST_EXPRESS_PASSWORD")
	endpoint := os.Getenv("POST_EXPRESS_API_URL")

	// Используем реальные учетные данные если не указаны переменные окружения
	if username == "" {
		username = "b2b@svetu.rs"
	}
	if password == "" {
		password = "Vu3@#45$%67&*89()"
	}
	if endpoint == "" {
		endpoint = "http://212.62.32.201/WspWebApi/transakcija" // Правильный URL из документации
	}

	// Конфигурация WSP клиента с исправленными параметрами
	config := &service.WSPConfig{
		Endpoint:        endpoint,
		Username:        username,
		Password:        password,
		Language:        "sr",
		DeviceType:      "2", // ИСПРАВЛЕНО: строка "2" для веб-приложения
		Timeout:         30 * time.Second,
		MaxRetries:      2,
		RetryDelay:      1 * time.Second,
		TestMode:        false,
		DeviceName:      "SveTu-Server",
		ApplicationName: "SveTu-Platform",
		Version:         "1.0.0",
		PartnerID:       10109, // ДОБАВЛЕНО: Partner ID для svetu.rs
	}

	// Создание клиента
	client := service.NewWSPClient(config, *appLogger)
	ctx := context.Background()

	fmt.Println("=====================================")
	fmt.Println("Post Express Integration Test - FIXED")
	fmt.Println("=====================================")
	fmt.Printf("Endpoint: %s\n", endpoint)
	fmt.Printf("Username: %s\n", username)
	fmt.Printf("Partner ID: %d\n", config.PartnerID)
	fmt.Println("=====================================\n")

	// Тест 1: Поиск населенных пунктов
	fmt.Println("Test 1: Searching for locations (Novi Sad)...")
	locations, err := client.GetLocations(ctx, "Novi Sad")
	if err != nil {
		log.Printf("❌ Failed to get locations: %v", err)
	} else {
		fmt.Printf("✅ Found %d locations\n", len(locations))
		for i, loc := range locations {
			if i < 3 { // Показываем первые 3 результата
				fmt.Printf("  - %s (ID: %d, Postal: %s)\n", loc.Name, loc.ID, loc.PostalCode)
			}
		}
		fmt.Println()
	}

	// Тест 2: Создание посылки через манифест (транзакция 73)
	fmt.Println("Test 2: Creating shipment via manifest (transaction 73)...")

	shipmentRequest := &service.WSPShipmentRequest{
		SenderName:          "Sve Tu Platform",
		SenderAddress:       "Улица Микија Манојловића 53",
		SenderCity:          "Нови Сад",
		SenderPostalCode:    "21000",
		SenderPhone:         "+381 60 1234567",
		RecipientName:       "Тест Получатель",
		RecipientAddress:    "Булевар Ослобођења 100",
		RecipientCity:       "Београд",
		RecipientPostalCode: "11000",
		RecipientPhone:      "+381 60 7654321",
		Weight:              1.5,
		CODAmount:           0,
		InsuranceAmount:     1000,
		ServiceType:         "PE_Danas_za_sutra_12",
		Content:             "Тестовая посылка - электроника",
		Note:                "Осторожно! Хрупкое!",
	}

	// Создание манифеста с исправленной транзакцией
	manifestResp, err := client.CreateShipmentViaManifest(ctx, shipmentRequest)
	if err != nil {
		log.Printf("❌ Failed to create shipment via manifest: %v", err)
	} else {
		if manifestResp.Success {
			fmt.Println("✅ Manifest created successfully!")
			if manifestResp.Manifest != nil {
				fmt.Printf("  Manifest Number: %s\n", manifestResp.Manifest.BrojManifesta)
				fmt.Printf("  Status: %s\n", manifestResp.Manifest.Status)
			}
			if len(manifestResp.Posiljke) > 0 {
				for _, p := range manifestResp.Posiljke {
					fmt.Printf("  Shipment %s:\n", p.BrojPosiljke)
					fmt.Printf("    - Post Express Number: %s\n", p.PostExpressBroj)
					fmt.Printf("    - Barcode: %s\n", p.Barkod)
					fmt.Printf("    - Status: %s\n", p.Status)
					if p.Greska != "" {
						fmt.Printf("    - Error: %s\n", p.Greska)
					}
				}
			}
		} else {
			fmt.Printf("❌ Manifest creation failed: %s\n", manifestResp.ErrorMessage)
		}
		fmt.Println()
	}

	// Тест 3: Отслеживание посылки
	if manifestResp != nil && len(manifestResp.Posiljke) > 0 && manifestResp.Posiljke[0].PostExpressBroj != "" {
		trackingNumber := manifestResp.Posiljke[0].PostExpressBroj
		fmt.Printf("Test 3: Tracking shipment %s...\n", trackingNumber)

		tracking, err := client.GetShipmentStatus(ctx, trackingNumber)
		if err != nil {
			log.Printf("❌ Failed to track shipment: %v", err)
		} else {
			fmt.Printf("✅ Tracking info received:\n")
			fmt.Printf("  Status: %s\n", tracking.Status)
			fmt.Printf("  Events: %d\n", len(tracking.Events))
			for _, event := range tracking.Events {
				fmt.Printf("    - %s %s: %s (%s)\n", event.Date, event.Time, event.Description, event.Location)
			}
			fmt.Println()
		}
	}

	// Тест 4: Тестирование транзакции напрямую
	fmt.Println("Test 4: Testing direct transaction with correct parameters...")

	// Подготовка тестового запроса
	testInput := map[string]interface{}{
		"TestField": "TestValue",
	}

	inputJSON, _ := json.Marshal(testInput)

	req := &models.TransactionRequest{
		TransactionType: 1, // Тестовая транзакция
		InputData:       string(inputJSON),
	}

	resp, err := client.Transaction(ctx, req)
	if err != nil {
		log.Printf("❌ Transaction failed: %v", err)
	} else {
		fmt.Printf("✅ Transaction completed:\n")
		fmt.Printf("  Success: %v\n", resp.Success)
		if resp.ErrorMessage != nil {
			fmt.Printf("  Error: %s\n", *resp.ErrorMessage)
		}
		if len(resp.OutputData) > 0 {
			fmt.Printf("  Output data received: %d bytes\n", len(resp.OutputData))
		}
	}

	fmt.Println("\n=====================================")
	fmt.Println("Integration test completed!")
	fmt.Println("=====================================")
}
