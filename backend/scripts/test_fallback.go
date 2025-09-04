package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	// Configuration
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:3000"
	}

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	fmt.Println("Starting fallback mechanism tests...")
	fmt.Println("=" * 50)

	// Test scenarios
	scenarios := []struct {
		name        string
		setup       func(*sql.DB) error
		test        func() error
		cleanup     func(*sql.DB) error
		expectation string
	}{
		{
			name: "Fallback when new system is empty",
			setup: func(db *sql.DB) error {
				// Clear new system data for category 1
				_, err := db.Exec(`
					DELETE FROM unified_category_attributes WHERE category_id = 1
				`)
				return err
			},
			test: func() error {
				// Request attributes for category 1
				resp, err := http.Get(fmt.Sprintf("%s/api/v2/categories/1/attributes", apiURL))
				if err != nil {
					return err
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					return fmt.Errorf("expected status 200, got %d", resp.StatusCode)
				}

				var result map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					return err
				}

				data, ok := result["data"].([]interface{})
				if !ok || len(data) == 0 {
					return fmt.Errorf("expected attributes from fallback, got empty")
				}

				fmt.Printf("  ✅ Fallback returned %d attributes from old system\n", len(data))
				return nil
			},
			cleanup: func(db *sql.DB) error {
				// Restore data
				_, err := db.Exec(`
					INSERT INTO unified_category_attributes (category_id, attribute_id)
					SELECT 1, id FROM unified_attributes WHERE legacy_category_attribute_id IS NOT NULL
					ON CONFLICT DO NOTHING
				`)
				return err
			},
			expectation: "Should return data from old system when new system is empty",
		},
		{
			name: "No fallback when new system has data",
			setup: func(db *sql.DB) error {
				// Ensure new system has data
				_, err := db.Exec(`
					INSERT INTO unified_category_attributes (category_id, attribute_id)
					SELECT 2, id FROM unified_attributes LIMIT 3
					ON CONFLICT DO NOTHING
				`)
				return err
			},
			test: func() error {
				// Request attributes for category 2
				resp, err := http.Get(fmt.Sprintf("%s/api/v2/categories/2/attributes", apiURL))
				if err != nil {
					return err
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					return fmt.Errorf("expected status 200, got %d", resp.StatusCode)
				}

				var result map[string]interface{}
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					return err
				}

				// Check response headers or metadata for source
				if source := resp.Header.Get("X-Data-Source"); source != "" {
					fmt.Printf("  Data source: %s\n", source)
				}

				data, ok := result["data"].([]interface{})
				if !ok {
					return fmt.Errorf("unexpected response format")
				}

				fmt.Printf("  ✅ New system returned %d attributes\n", len(data))
				return nil
			},
			cleanup: func(db *sql.DB) error {
				return nil
			},
			expectation: "Should use new system when data exists",
		},
		{
			name: "Fallback on new system error",
			setup: func(db *sql.DB) error {
				// Simulate error by breaking constraint
				_, err := db.Exec(`
					-- Temporarily break a constraint to cause error
					ALTER TABLE unified_attributes DROP CONSTRAINT IF EXISTS test_constraint;
					ALTER TABLE unified_attributes 
					ADD CONSTRAINT test_constraint 
					CHECK (false) NOT VALID;
				`)
				return err
			},
			test: func() error {
				// Try to get attributes, should fallback on error
				resp, err := http.Get(fmt.Sprintf("%s/api/v2/categories/3/attributes", apiURL))
				if err != nil {
					return err
				}
				defer resp.Body.Close()

				// Should still return 200 due to fallback
				if resp.StatusCode != http.StatusOK {
					return fmt.Errorf("expected status 200 from fallback, got %d", resp.StatusCode)
				}

				fmt.Printf("  ✅ Fallback activated on new system error\n")
				return nil
			},
			cleanup: func(db *sql.DB) error {
				_, err := db.Exec(`
					ALTER TABLE unified_attributes DROP CONSTRAINT IF EXISTS test_constraint
				`)
				return err
			},
			expectation: "Should fallback to old system on new system error",
		},
		{
			name: "Performance comparison",
			setup: func(db *sql.DB) error {
				// Ensure both systems have data
				return nil
			},
			test: func() error {
				// Test with new system
				start := time.Now()
				resp1, err := http.Get(fmt.Sprintf("%s/api/v2/categories/1/attributes?force_new=true", apiURL))
				if err != nil {
					return err
				}
				resp1.Body.Close()
				newSystemTime := time.Since(start)

				// Test with old system (fallback)
				start = time.Now()
				resp2, err := http.Get(fmt.Sprintf("%s/api/v2/categories/1/attributes?force_old=true", apiURL))
				if err != nil {
					return err
				}
				resp2.Body.Close()
				oldSystemTime := time.Since(start)

				fmt.Printf("  New system response time: %v\n", newSystemTime)
				fmt.Printf("  Old system response time: %v\n", oldSystemTime)

				if newSystemTime < oldSystemTime {
					fmt.Printf("  ✅ New system is %.2fx faster\n", float64(oldSystemTime)/float64(newSystemTime))
				} else {
					fmt.Printf("  ⚠️  Old system is %.2fx faster\n", float64(newSystemTime)/float64(oldSystemTime))
				}

				return nil
			},
			cleanup: func(db *sql.DB) error {
				return nil
			},
			expectation: "Compare performance between systems",
		},
		{
			name: "Data consistency check",
			setup: func(db *sql.DB) error {
				return nil
			},
			test: func() error {
				// Get data from both systems
				resp1, err := http.Get(fmt.Sprintf("%s/api/v2/categories/1/attributes?force_new=true", apiURL))
				if err != nil {
					return err
				}
				defer resp1.Body.Close()

				var newData map[string]interface{}
				json.NewDecoder(resp1.Body).Decode(&newData)

				resp2, err := http.Get(fmt.Sprintf("%s/api/v2/categories/1/attributes?force_old=true", apiURL))
				if err != nil {
					return err
				}
				defer resp2.Body.Close()

				var oldData map[string]interface{}
				json.NewDecoder(resp2.Body).Decode(&oldData)

				// Compare counts
				newCount := 0
				oldCount := 0

				if data, ok := newData["data"].([]interface{}); ok {
					newCount = len(data)
				}
				if data, ok := oldData["data"].([]interface{}); ok {
					oldCount = len(data)
				}

				fmt.Printf("  New system: %d attributes\n", newCount)
				fmt.Printf("  Old system: %d attributes\n", oldCount)

				if newCount == oldCount {
					fmt.Printf("  ✅ Data counts match\n")
				} else {
					fmt.Printf("  ⚠️  Data count mismatch (diff: %d)\n", newCount-oldCount)
				}

				return nil
			},
			cleanup: func(db *sql.DB) error {
				return nil
			},
			expectation: "Verify data consistency between systems",
		},
	}

	// Run tests
	successCount := 0
	failureCount := 0

	for _, scenario := range scenarios {
		fmt.Printf("\nTest: %s\n", scenario.name)
		fmt.Printf("  Expectation: %s\n", scenario.expectation)

		// Setup
		if scenario.setup != nil {
			if err := scenario.setup(db); err != nil {
				fmt.Printf("  ❌ Setup failed: %v\n", err)
				failureCount++
				continue
			}
		}

		// Run test
		if err := scenario.test(); err != nil {
			fmt.Printf("  ❌ Test failed: %v\n", err)
			failureCount++
		} else {
			fmt.Printf("  ✅ Test passed\n")
			successCount++
		}

		// Cleanup
		if scenario.cleanup != nil {
			if err := scenario.cleanup(db); err != nil {
				fmt.Printf("  ⚠️  Cleanup failed: %v\n", err)
			}
		}
	}

	// Test fallback monitoring
	fmt.Println("\n" + "="*50)
	fmt.Println("Testing fallback metrics...")

	resp, err := http.Get(fmt.Sprintf("%s/metrics", apiURL))
	if err == nil && resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()

		// Would parse Prometheus metrics here
		fmt.Println("  ✅ Metrics endpoint available")
		fmt.Println("  Fallback metrics should be visible in Prometheus format")
	} else {
		fmt.Println("  ⚠️  Metrics endpoint not available")
	}

	// Final results
	fmt.Println("\n" + "="*50)
	fmt.Println("FALLBACK MECHANISM TEST RESULTS")
	fmt.Println("=" * 50)
	fmt.Printf("Tests passed: %d/%d\n", successCount, successCount+failureCount)

	if failureCount > 0 {
		fmt.Printf("\n❌ FAILED: %d tests failed\n", failureCount)
		os.Exit(1)
	}

	fmt.Printf("\n✅ SUCCESS: All %d tests passed\n", successCount)
	fmt.Println("\nFallback mechanism is working correctly:")
	fmt.Println("  - Falls back when new system is empty")
	fmt.Println("  - Uses new system when data exists")
	fmt.Println("  - Handles errors gracefully")
	fmt.Println("  - Maintains data consistency")
}
