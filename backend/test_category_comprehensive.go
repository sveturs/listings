package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type TestCase struct {
	Name        string
	RequestData map[string]interface{}
	Expected    int // Expected category ID
}

func main() {
	testCases := []TestCase{
		{
			Name: "Volkswagen Touran - Минивэн",
			RequestData: map[string]interface{}{
				"title":       "Volkswagen Touran 2017",
				"description": "Минивэн в отличном состоянии",
				"aiHints": map[string]interface{}{
					"domain":      "automotive",
					"productType": "minivan",
					"keywords":    []string{"volkswagen", "touran", "минивэн"},
				},
			},
			Expected: 1301, // Cars category
		},
		{
			Name: "Mercedes-Benz E-Class - Седан",
			RequestData: map[string]interface{}{
				"title":       "Mercedes-Benz E220d 2020",
				"description": "Бизнес седан, полный привод",
				"aiHints": map[string]interface{}{
					"domain":      "automotive",
					"productType": "sedan",
					"keywords":    []string{"mercedes", "седан", "e-class"},
				},
			},
			Expected: 1301,
		},
		{
			Name: "BMW X5 - Внедорожник",
			RequestData: map[string]interface{}{
				"title":       "BMW X5 xDrive 2019",
				"description": "Внедорожник премиум класса",
				"aiHints": map[string]interface{}{
					"domain":      "automotive",
					"productType": "suv",
					"keywords":    []string{"bmw", "x5", "внедорожник", "кроссовер"},
				},
			},
			Expected: 1301,
		},
		{
			Name: "Yamaha R1 - Мотоцикл",
			RequestData: map[string]interface{}{
				"title":       "Yamaha R1 2021",
				"description": "Спортивный мотоцикл",
				"aiHints": map[string]interface{}{
					"domain":      "automotive",
					"productType": "motorcycle",
					"keywords":    []string{"yamaha", "мотоцикл", "r1"},
				},
			},
			Expected: 1302, // Motorcycles category
		},
		{
			Name: "Vespa Primavera - Скутер",
			RequestData: map[string]interface{}{
				"title":       "Vespa Primavera 125",
				"description": "Итальянский скутер",
				"aiHints": map[string]interface{}{
					"domain":      "automotive",
					"productType": "scooter",
					"keywords":    []string{"vespa", "скутер", "primavera"},
				},
			},
			Expected: 1302, // Motorcycles category
		},
		{
			Name: "Michelin Pilot Sport - Шины",
			RequestData: map[string]interface{}{
				"title":       "Michelin Pilot Sport 4 225/45 R17",
				"description": "Летние шины для легкового автомобиля",
				"aiHints": map[string]interface{}{
					"domain":      "automotive",
					"productType": "tires",
					"keywords":    []string{"michelin", "шины", "резина"},
				},
			},
			Expected: 1303, // Auto parts category
		},
		{
			Name: "BBS Диски - Колеса",
			RequestData: map[string]interface{}{
				"title":       "BBS CH-R 19 дюймов",
				"description": "Литые диски для BMW",
				"aiHints": map[string]interface{}{
					"domain":      "automotive",
					"productType": "wheels",
					"keywords":    []string{"bbs", "диски", "колеса"},
				},
			},
			Expected: 1303, // Auto parts category
		},
	}

	fmt.Println("=== TESTING AUTOMOTIVE CATEGORY DETECTION ===\n")

	successCount := 0
	failureCount := 0

	for _, tc := range testCases {
		fmt.Printf("Testing: %s\n", tc.Name)

		jsonData, _ := json.Marshal(tc.RequestData)

		req, err := http.NewRequest("POST", "http://localhost:3000/api/v1/marketplace/categories/detect", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("  ❌ Error creating request: %v\n\n", err)
			failureCount++
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("  ❌ Error making request: %v\n\n", err)
			failureCount++
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("  ❌ Error reading response: %v\n\n", err)
			failureCount++
			continue
		}

		var result struct {
			Data struct {
				CategoryID     int     `json:"category_id"`
				CategoryName   string  `json:"category_name"`
				CategorySlug   string  `json:"category_slug"`
				ConfidenceScore float64 `json:"confidence_score"`
				Method         string  `json:"method"`
				ProcessingTime int     `json:"processing_time_ms"`
			} `json:"data"`
			Success bool `json:"success"`
		}

		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Printf("  ❌ Error parsing response: %v\n", err)
			fmt.Printf("  Raw response: %s\n\n", string(body))
			failureCount++
			continue
		}

		if result.Data.CategoryID == tc.Expected {
			fmt.Printf("  ✅ SUCCESS: Detected category %d (%s)\n", result.Data.CategoryID, result.Data.CategoryName)
			fmt.Printf("     Method: %s, Confidence: %.2f, Time: %dms\n\n",
				result.Data.Method, result.Data.ConfidenceScore, result.Data.ProcessingTime)
			successCount++
		} else {
			fmt.Printf("  ❌ FAILURE: Expected %d, got %d (%s)\n", tc.Expected, result.Data.CategoryID, result.Data.CategoryName)
			fmt.Printf("     Method: %s, Confidence: %.2f\n\n", result.Data.Method, result.Data.ConfidenceScore)
			failureCount++
		}
	}

	fmt.Println("=== TEST SUMMARY ===")
	fmt.Printf("✅ Passed: %d\n", successCount)
	fmt.Printf("❌ Failed: %d\n", failureCount)
	fmt.Printf("Total: %d\n", len(testCases))

	if failureCount == 0 {
		fmt.Println("\n🎉 All tests passed! Category detection is working correctly!")
	} else {
		fmt.Printf("\n⚠️  %d tests failed. Review the results above.\n", failureCount)
	}
}