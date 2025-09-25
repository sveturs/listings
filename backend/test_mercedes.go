package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	// Test with different variations to understand the issue
	tests := []map[string]interface{}{
		{
			"title":       "Mercedes-Benz E220d 2020",
			"description": "Бизнес седан, полный привод",
			"aiHints": map[string]interface{}{
				"domain":      "automotive",
				"productType": "sedan",
				"keywords":    []string{"mercedes", "седан", "e-class"},
			},
		},
		{
			"title":       "Mercedes E220d 2020",
			"description": "Седан бизнес класса",
			"aiHints": map[string]interface{}{
				"domain":      "automotive",
				"productType": "sedan",
				"keywords":    []string{"mercedes", "седан"},
			},
		},
		{
			"title":       "Автомобиль Mercedes-Benz E220d",
			"description": "Легковой автомобиль седан",
			"aiHints": map[string]interface{}{
				"domain":      "automotive",
				"productType": "car",
				"keywords":    []string{"автомобиль", "mercedes", "седан"},
			},
		},
	}

	for i, testData := range tests {
		fmt.Printf("\n=== Test %d ===\n", i+1)
		fmt.Printf("Title: %s\n", testData["title"])

		jsonData, _ := json.Marshal(testData)

		req, err := http.NewRequest("POST", "http://localhost:3000/api/v1/marketplace/categories/detect", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error making request: %v\n", err)
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response: %v\n", err)
			continue
		}

		// Parse and display result
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err == nil {
			if data, ok := result["data"].(map[string]interface{}); ok {
				fmt.Printf("Result: Category %v (%v)\n", data["category_id"], data["category_name"])
				fmt.Printf("Method: %v, Confidence: %v\n", data["method"], data["confidence_score"])

				// Show matched keywords if available
				if matched, ok := data["matched_keywords"].([]interface{}); ok {
					fmt.Printf("Matched Keywords: %v\n", matched)
				}
			}
		} else {
			fmt.Printf("Raw response: %s\n", string(body))
		}
	}
}