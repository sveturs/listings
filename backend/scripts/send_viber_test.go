//go:build ignore

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run send_viber_test.go <your_phone_number>")
		fmt.Println("Example: go run send_viber_test.go 381604485063")
		os.Exit(1)
	}

	phoneNumber := os.Args[1]

	// Отправляем простое тестовое сообщение
	payload := map[string]interface{}{
		"viber_id": phoneNumber,
		"text":     "🎉 Привет! Это тестовое сообщение от SveTu бота!\n\nЕсли ты видишь это сообщение - значит всё работает! 🚀",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Отправляем запрос к локальному API
	resp, err := http.Post(
		"http://localhost:3000/api/viber/send",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Response: %s\n", string(body))

	if resp.StatusCode == 200 {
		fmt.Println("\n✅ Сообщение отправлено! Проверь Viber на номере", phoneNumber)
		fmt.Println("\n📱 Что ты должен увидеть в Viber:")
		fmt.Println("   1. Сообщение от SveTu бота")
		fmt.Println("   2. Если нет - возможно нужно найти бота в Viber и подписаться")
		fmt.Println("   3. Или настроить webhook в Infobip портале")
	} else {
		fmt.Println("\n❌ Ошибка отправки! Проверь:")
		fmt.Println("   1. Backend запущен на порту 3000")
		fmt.Println("   2. Конфигурация Infobip в .env.dev правильная")
		fmt.Println("   3. Номер телефона в правильном формате (381XXXXXXXXX)")
	}
}
