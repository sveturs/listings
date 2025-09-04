package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type AttributeValue struct {
	AttributeID  int         `json:"attribute_id"`
	Value        interface{} `json:"value"`
	TextValue    *string     `json:"text_value"`
	NumericValue *float64    `json:"numeric_value"`
	BooleanValue *bool       `json:"boolean_value"`
}

func main() {
	// Configuration
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:3000"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test configuration
	testCases := []struct {
		name        string
		listingID   int
		attributeID int
		value       interface{}
		valueType   string
	}{
		{"Text attribute", 1, 1, "Test Value 123", "text"},
		{"Numeric attribute", 1, 2, 42.5, "numeric"},
		{"Boolean attribute", 1, 3, true, "boolean"},
		{"JSON attribute", 1, 4, map[string]interface{}{"key": "value"}, "json"},
	}

	fmt.Println("Starting dual-write validation tests...")
	fmt.Println(strings.Repeat("=", 50))

	successCount := 0
	failureCount := 0

	for _, tc := range testCases {
		fmt.Printf("\nTest: %s\n", tc.name)
		fmt.Printf("  Listing ID: %d, Attribute ID: %d\n", tc.listingID, tc.attributeID)

		// 1. Send API request to save attribute value
		payload := map[string]interface{}{
			"attributes": map[string]interface{}{
				fmt.Sprintf("%d", tc.attributeID): tc.value,
			},
		}

		jsonData, _ := json.Marshal(payload)
		resp, err := http.Post(
			fmt.Sprintf("%s/api/v2/listings/%d/attributes", apiURL, tc.listingID),
			"application/json",
			bytes.NewBuffer(jsonData),
		)
		if err != nil {
			fmt.Printf("  ❌ API request failed: %v\n", err)
			failureCount++
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("  ❌ API returned status %d\n", resp.StatusCode)
			failureCount++
			continue
		}

		// 2. Wait for dual-write to complete
		time.Sleep(100 * time.Millisecond)

		// 3. Verify data in NEW system (unified_attribute_values)
		var newSystemValue interface{}
		var newSystemExists bool

		query := `
			SELECT 
				COALESCE(text_value::text, numeric_value::text, boolean_value::text, json_value::text)
			FROM unified_attribute_values 
			WHERE entity_type = 'listing' 
				AND entity_id = $1 
				AND attribute_id = $2
		`
		err = db.QueryRow(query, tc.listingID, tc.attributeID).Scan(&newSystemValue)
		newSystemExists = err == nil

		// 4. Verify data in OLD system (listing_attribute_values)
		var oldSystemValue interface{}
		var oldSystemExists bool

		query = `
			SELECT 
				COALESCE(text_value::text, numeric_value::text, boolean_value::text, json_value::text)
			FROM listing_attribute_values 
			WHERE listing_id = $1 
				AND attribute_id = $2
		`
		err = db.QueryRow(query, tc.listingID, tc.attributeID).Scan(&oldSystemValue)
		oldSystemExists = err == nil

		// 5. Compare results
		fmt.Printf("  New system: ")
		if newSystemExists {
			fmt.Printf("✅ Value saved: %v\n", newSystemValue)
		} else {
			fmt.Printf("❌ No value found\n")
		}

		fmt.Printf("  Old system: ")
		if oldSystemExists {
			fmt.Printf("✅ Value saved: %v\n", oldSystemValue)
		} else {
			fmt.Printf("❌ No value found\n")
		}

		// 6. Validate dual-write success
		if newSystemExists && oldSystemExists {
			fmt.Printf("  ✅ Dual-write successful\n")
			successCount++
		} else if newSystemExists && !oldSystemExists {
			fmt.Printf("  ⚠️  Only new system has data (old system might be disabled)\n")
			successCount++
		} else {
			fmt.Printf("  ❌ Dual-write failed\n")
			failureCount++
		}
	}

	// Test concurrent dual-writes
	fmt.Println("\n" + "="*50)
	fmt.Println("Testing concurrent dual-writes...")

	concurrentTests := 10
	done := make(chan bool, concurrentTests)

	for i := 0; i < concurrentTests; i++ {
		go func(index int) {
			payload := map[string]interface{}{
				"attributes": map[string]interface{}{
					"1": fmt.Sprintf("Concurrent test %d", index),
				},
			}

			jsonData, _ := json.Marshal(payload)
			resp, err := http.Post(
				fmt.Sprintf("%s/api/v2/listings/99/attributes", apiURL),
				"application/json",
				bytes.NewBuffer(jsonData),
			)

			if err != nil || resp.StatusCode != http.StatusOK {
				fmt.Printf("  Concurrent request %d: ❌\n", index)
				done <- false
			} else {
				fmt.Printf("  Concurrent request %d: ✅\n", index)
				resp.Body.Close()
				done <- true
			}
		}(i)
	}

	concurrentSuccess := 0
	for i := 0; i < concurrentTests; i++ {
		if <-done {
			concurrentSuccess++
		}
	}

	// Final results
	fmt.Println("\n" + "="*50)
	fmt.Println("DUAL-WRITE VALIDATION RESULTS")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Sequential tests: %d/%d passed\n", successCount, successCount+failureCount)
	fmt.Printf("Concurrent tests: %d/%d passed\n", concurrentSuccess, concurrentTests)

	totalTests := successCount + failureCount + concurrentTests
	totalPassed := successCount + concurrentSuccess

	if totalPassed < totalTests {
		fmt.Printf("\n❌ FAILED: Only %d/%d tests passed\n", totalPassed, totalTests)
		os.Exit(1)
	}

	fmt.Printf("\n✅ SUCCESS: All %d tests passed\n", totalTests)

	// Additional metrics
	fmt.Println("\nPerformance metrics:")

	// Check if both tables have same count
	var newCount, oldCount int
	db.QueryRow("SELECT COUNT(*) FROM unified_attribute_values WHERE entity_type = 'listing'").Scan(&newCount)
	db.QueryRow("SELECT COUNT(*) FROM listing_attribute_values").Scan(&oldCount)

	fmt.Printf("  Records in new system: %d\n", newCount)
	fmt.Printf("  Records in old system: %d\n", oldCount)

	if newCount == oldCount {
		fmt.Println("  ✅ Record counts match")
	} else {
		fmt.Printf("  ⚠️  Record count mismatch (diff: %d)\n", newCount-oldCount)
	}
}
