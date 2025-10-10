//go:build ignore

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Загружаем .env
	if err := godotenv.Load("../.env.dev"); err != nil {
		log.Printf("Warning: Error loading .env.dev file: %v\n", err)
	}

	apiKey := os.Getenv("INFOBIP_API_KEY")
	baseURL := os.Getenv("INFOBIP_BASE_URL")
	senderID := os.Getenv("INFOBIP_SENDER_ID")

	fmt.Println("=== Infobip Configuration ===")
	fmt.Printf("API Key: %s\n", apiKey)
	fmt.Printf("Base URL: %s\n", baseURL)
	fmt.Printf("Sender ID: %s\n", senderID)
	fmt.Println()

	if apiKey == "" || baseURL == "" || senderID == "" {
		log.Fatal("❌ Missing configuration! Check INFOBIP_API_KEY, INFOBIP_BASE_URL, INFOBIP_SENDER_ID in .env.dev")
	}

	// Отправляем тестовое сообщение
	phoneNumber := "381604485063"

	// Формируем запрос по спецификации Infobip API v2
	payload := map[string]interface{}{
		"messages": []map[string]interface{}{
			{
				"sender": senderID,
				"destinations": []map[string]string{
					{"to": phoneNumber},
				},
				"content": map[string]interface{}{
					"type": "TEXT",
					"text": "🎉 Тест от SveTu! Если ты видишь это - интеграция работает!",
				},
			},
		},
	}

	jsonData, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	fmt.Println("=== Request Payload ===")
	fmt.Println(string(jsonData))
	fmt.Println()

	// Создаём HTTP запрос
	url := fmt.Sprintf("https://%s/viber/2/messages", baseURL)
	fmt.Printf("=== Sending to: %s ===\n", url)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "App "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	fmt.Println("=== Request Headers ===")
	for key, values := range req.Header {
		for _, value := range values {
			if key == "Authorization" {
				// Скрываем полный токен
				fmt.Printf("%s: App %s...%s\n", key, apiKey[:10], apiKey[len(apiKey)-10:])
			} else {
				fmt.Printf("%s: %s\n", key, value)
			}
		}
	}
	fmt.Println()

	// Отправляем запрос
	fmt.Println("=== Sending Request ===")
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("❌ Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Читаем ответ
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	fmt.Printf("=== Response Status: %d ===\n", resp.StatusCode)
	fmt.Println("=== Response Headers ===")
	for key, values := range resp.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", key, value)
		}
	}
	fmt.Println()

	fmt.Println("=== Response Body ===")

	// Пытаемся распарсить как JSON для красивого вывода
	var prettyJSON map[string]interface{}
	if err := json.Unmarshal(respBody, &prettyJSON); err == nil {
		formatted, _ := json.MarshalIndent(prettyJSON, "", "  ")
		fmt.Println(string(formatted))
	} else {
		// Если не JSON, выводим как есть
		fmt.Println(string(respBody))
	}
	fmt.Println()

	// Анализ результата
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		fmt.Println("✅ SUCCESS! Message sent successfully!")
		fmt.Printf("\n📱 Check Viber on phone number: %s\n", phoneNumber)
	} else {
		fmt.Println("❌ FAILED! Check error details above.")
		fmt.Println("\nPossible issues:")
		fmt.Println("1. Invalid Sender ID - check in Infobip portal")
		fmt.Println("2. Phone number not subscribed to bot")
		fmt.Println("3. Bot not approved yet")
		fmt.Println("4. API key not authorized for this sender")
		fmt.Println("5. Incorrect Base URL")
	}
}
