// run_tests_manually.go - Temporary script to run tests manually
//go:build ignore

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"

	"backend/internal/config"
	testingService "backend/internal/proj/admin/testing/service"
	testingStorage "backend/internal/proj/admin/testing/storage/postgres"
	"backend/internal/storage/postgres"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Initialize database
	ctx := context.Background()
	db, err := postgres.NewDatabase(ctx, cfg.DatabaseURL, nil, "", nil, cfg.SearchWeights, cfg)
	if err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		os.Exit(1)
	}

	// Initialize testing components
	testStorage := testingStorage.NewStorage(db.GetSQLXDB(), logger)
	testAuthMgr := testingService.NewTestAuthManager(cfg.BackendURL, "admin@admin.rs", "P@$S4@dmi№", logger)
	testRunner := testingService.NewTestRunner(testStorage, testAuthMgr, cfg.BackendURL, logger)

	// Run API endpoint tests
	fmt.Println("Running API endpoint tests...")
	testRun, err := testRunner.RunTestSuite(ctx, "api-endpoints", 1, true)
	if err != nil {
		fmt.Printf("Failed to run test suite: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Test run created: ID=%d, UUID=%s, Status=%s\n", testRun.ID, testRun.RunUUID, testRun.Status)
	fmt.Println("Tests are running in background. Check database for results.")
}
