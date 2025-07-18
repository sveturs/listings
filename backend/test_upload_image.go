package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// Создаем multipart форму
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Добавляем файл
	file, err := os.Open("test_product_image.jpg")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", filepath.Base("test_product_image.jpg"))
	if err != nil {
		fmt.Printf("Error creating form file: %v\n", err)
		return
	}
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Printf("Error copying file: %v\n", err)
		return
	}

	// Добавляем другие поля
	writer.WriteField("alt_text", "Test Product Image")
	writer.WriteField("is_primary", "true")

	err = writer.Close()
	if err != nil {
		fmt.Printf("Error closing writer: %v\n", err)
		return
	}

	// Создаем запрос
	req, err := http.NewRequest("POST", "http://localhost:3000/api/v1/storefronts/slug/test_vitrina-1751210133/products/4109/images", &requestBody)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTI5NTY1NTQsImlkIjoxfQ.e1jLhe3UFNIqH0470KkecbbBDh0dHmrkUisr82i8h4A")

	// Отправляем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Response: %s\n", string(body))
}