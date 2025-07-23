package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	fmt.Println("=== Testing Category Cache Performance ===")
	fmt.Println()

	// Тестируем производительность с кешем и без
	locales := []string{"en", "ru", "sr"}

	for _, locale := range locales {
		fmt.Printf("\nTesting locale: %s\n", locale)
		fmt.Println(strings.Repeat("-", 40))

		// Первый запрос (может быть из кеша или из БД)
		duration1, count1 := makeRequest(locale, "First")

		// Второй запрос (должен быть из кеша)
		duration2, count2 := makeRequest(locale, "Second")

		// Третий запрос (точно из кеша)
		duration3, count3 := makeRequest(locale, "Third")

		fmt.Printf("\nResults for locale %s:\n", locale)
		fmt.Printf("  First request:  %v (%d categories)\n", duration1, count1)
		fmt.Printf("  Second request: %v (%d categories)\n", duration2, count2)
		fmt.Printf("  Third request:  %v (%d categories)\n", duration3, count3)

		// Если второй и третий запросы значительно быстрее первого, кеш работает
		switch {
		case duration2 < duration1/2 && duration3 < duration1/2:
			fmt.Println("  ✓ Cache appears to be working!")
		case duration2 < time.Duration(float64(duration1)*0.8) || duration3 < time.Duration(float64(duration1)*0.8):
			fmt.Println("  ~ Cache might be working (small improvement)")
		default:
			fmt.Println("  ✗ Cache doesn't seem to be working")
		}
	}

	// Тест массовых запросов
	fmt.Println("\n\nBulk request test (100 requests):")
	fmt.Println(strings.Repeat("=", 50))

	start := time.Now()
	totalCategories := 0

	for i := 0; i < 100; i++ {
		locale := locales[i%3]
		_, count := makeRequest(locale, "")
		totalCategories += count
	}

	elapsed := time.Since(start)
	avgTime := elapsed / 100

	fmt.Printf("Total time: %v\n", elapsed)
	fmt.Printf("Average time per request: %v\n", avgTime)
	fmt.Printf("Total categories fetched: %d\n", totalCategories)

	switch {
	case avgTime < 10*time.Millisecond:
		fmt.Println("✓ Excellent performance - cache is definitely working!")
	case avgTime < 50*time.Millisecond:
		fmt.Println("✓ Good performance - cache is likely working")
	default:
		fmt.Println("✗ Poor performance - cache might not be working")
	}
}

func makeRequest(locale string, label string) (time.Duration, int) {
	start := time.Now()

	url := "http://localhost:3000/api/v1/marketplace/categories"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return 0, 0
	}

	if locale != "" {
		req.Header.Set("Accept-Language", locale)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return 0, 0
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("Failed to close response body: %v", closeErr)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		return 0, 0
	}

	var result struct {
		Data []interface{} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		return 0, 0
	}

	elapsed := time.Since(start)

	if label != "" {
		fmt.Printf("  %s request completed in %v\n", label, elapsed)
	}

	return elapsed, len(result.Data)
}

var strings = struct {
	Repeat func(s string, count int) string
}{
	Repeat: func(s string, count int) string {
		result := ""
		for i := 0; i < count; i++ {
			result += s
		}
		return result
	},
}
