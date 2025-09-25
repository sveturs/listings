package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	// Тестовые данные для Volkswagen Touran
	requestData := map[string]interface{}{
		"title":       "Volkswagen Touran 2017",
		"description": "Минивэн в отличном состоянии",
		"aiHints": map[string]interface{}{
			"domain":      "automotive",
			"productType": "minivan",
			"keywords":    []string{"volkswagen", "touran", "минивэн"},
		},
	}

	jsonData, _ := json.Marshal(requestData)

	// Попробуем простой эндпоинт без авторизации
	req, err := http.NewRequest("POST", "http://localhost:3000/api/v1/marketplace/categories/detect", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Response: %s\n", string(body))

	// Парсим JSON ответ для удобного вывода
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err == nil {
		if categoryId, ok := result["categoryId"]; ok {
			fmt.Printf("\n✅ Определена категория: %v\n", categoryId)
			if categoryName, ok := result["categoryName"]; ok {
				fmt.Printf("   Название: %v\n", categoryName)
			}
			if algorithm, ok := result["algorithm"]; ok {
				fmt.Printf("   Алгоритм: %v\n", algorithm)
			}
		}
	}
}