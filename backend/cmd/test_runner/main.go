// Простой тестовый runner для запуска functional tests напрямую
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"

	testingService "backend/internal/proj/admin/testing/service"
	testingStorage "backend/internal/proj/admin/testing/storage/postgres"
)

func main() {
	// Setup logger
	zerologLogger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Get database URL
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5432/svetubd?sslmode=disable"
	}

	// Connect to database
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create storage
	storage := testingStorage.NewStorage(db, zerologLogger)

	// Get backend URL from env
	backendURL := os.Getenv("BACKEND_URL")
	if backendURL == "" {
		backendURL = "http://localhost:3000"
	}

	// Check if we should use mock auth (for local testing without auth-service)
	useMockAuth := os.Getenv("USE_MOCK_AUTH")
	if useMockAuth == "" {
		useMockAuth = "true" // Default to mock auth for standalone runner
	}

	var testRunner *testingService.TestRunner

	if useMockAuth == "true" {
		fmt.Println("⚠️  Using MOCK authentication (no auth-service required)")
		mockAuthMgr, err := testingService.NewMockAuthManager(zerologLogger)
		if err != nil {
			log.Fatalf("Failed to create mock auth manager: %v", err)
		}
		testRunner = testingService.NewTestRunner(storage, mockAuthMgr, backendURL, zerologLogger)
	} else {
		fmt.Println("✅ Using REAL authentication (auth-service required)")
		testAdminEmail := os.Getenv("TEST_ADMIN_EMAIL")
		if testAdminEmail == "" {
			testAdminEmail = "admin@admin.rs"
		}

		testAdminPassword := os.Getenv("TEST_ADMIN_PASSWORD")
		if testAdminPassword == "" {
			log.Fatal("TEST_ADMIN_PASSWORD environment variable is required when USE_MOCK_AUTH=false")
		}

		realAuthMgr := testingService.NewTestAuthManager(backendURL, testAdminEmail, testAdminPassword, zerologLogger)
		testRunner = testingService.NewTestRunner(storage, realAuthMgr, backendURL, zerologLogger)
	}

	// Get test suite from env (default: functional-api)
	testSuite := os.Getenv("TEST_SUITE")
	if testSuite == "" {
		testSuite = "functional-api"
	}

	// Run tests
	ctx := context.Background()
	fmt.Printf("=== Starting %s Tests ===\n", testSuite)

	// Empty string for testName means run all tests in suite
	testRun, err := testRunner.RunTestSuite(ctx, testSuite, "", 11, false)
	if err != nil {
		log.Fatalf("Failed to run test suite: %v", err)
	}

	fmt.Printf("\n✅ Test run started successfully!\n")
	fmt.Printf("Run ID: %d\n", testRun.ID)
	fmt.Printf("Run UUID: %s\n", testRun.RunUUID)
	fmt.Printf("Status: %s\n", testRun.Status)

	// Wait a bit for tests to complete (async execution)
	fmt.Println("\n⏳ Waiting for tests to complete (60 seconds)...")
	fmt.Println("Check database for results: SELECT * FROM test_runs WHERE id = ", testRun.ID)

	// Get latest test run details
	fmt.Println("\n=== Fetching test results ===")

	// Simple polling loop
	for i := 0; i < 30; i++ {
		detail, err := testRunner.GetTestRunDetail(ctx, testRun.ID)
		if err != nil {
			log.Printf("Failed to get test run detail: %v", err)
			continue
		}

		if detail != nil && detail.TestRun != nil {
			fmt.Printf("\nStatus: %s\n", detail.TestRun.Status)
			fmt.Printf("Total: %d, Passed: %d, Failed: %d\n",
				detail.TestRun.TotalTests,
				detail.TestRun.PassedTests,
				detail.TestRun.FailedTests)

			if detail.TestRun.Status == "completed" || detail.TestRun.Status == "failed" {
				fmt.Println("\n=== Test Results ===")
				for _, result := range detail.Results {
					status := "✅"
					if result.Status == "failed" {
						status = "❌"
					} else if result.Status == "skipped" {
						status = "⏭️"
					}

					fmt.Printf("%s %s (%dms)\n", status, result.TestName, result.DurationMs)
					if result.ErrorMsg != nil && *result.ErrorMsg != "" {
						fmt.Printf("   Error: %s\n", *result.ErrorMsg)
					}
				}
				break
			}
		}

		// Wait 2 seconds before next check
		fmt.Print(".")
		time.Sleep(2 * time.Second)
	}

	fmt.Println("\n\n=== Test run completed ===")
}
