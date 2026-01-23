//go:build integration

package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vondi-global/listings/internal/domain"
	"github.com/vondi-global/listings/internal/repository/postgres"
	"github.com/vondi-global/listings/tests"
)

// TestDatabaseIntegration tests full database workflow
func TestDatabaseIntegration(t *testing.T) {
	tests.SkipIfNoDocker(t)

	// Setup test database
	testDB := tests.SetupTestPostgres(t)
	defer testDB.TeardownTestPostgres(t)

	// Run migrations
	tests.RunMigrations(t, testDB.DB, "../../migrations")

	// Create repository
	db := sqlx.NewDb(testDB.DB, "postgres")
	logger := zerolog.New(zerolog.NewTestWriter(t)).With().Timestamp().Logger()
	repo := postgres.NewRepository(db, logger)

	ctx := context.Background()

	t.Run("Create and Retrieve Listing", func(t *testing.T) {
		input := &domain.CreateListingInput{
			UserID:      1,
			Title:       "Integration Test Product",
			Description: stringPtr("Integration Test Description"),
			Price:       199.99,
			Currency:    "USD",
			CategoryID:  "200",
			Quantity:    5,
			SKU:         stringPtr("INT-TEST-001"),
		}

		// Create
		created, err := repo.CreateListing(ctx, input)
		require.NoError(t, err)
		require.NotNil(t, created)
		assert.Greater(t, created.ID, int64(0))

		// Retrieve by ID
		retrieved, err := repo.GetListingByID(ctx, created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, retrieved.ID)
		assert.Equal(t, input.Title, retrieved.Title)

		// Retrieve by UUID
		retrievedByUUID, err := repo.GetListingByUUID(ctx, created.UUID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, retrievedByUUID.ID)
	})

	t.Run("Update and Delete Workflow", func(t *testing.T) {
		// Create
		input := &domain.CreateListingInput{
			UserID:      2,
			Title:       "Update Test Product",
			Description: stringPtr("Update Test Description"),
			Price:       299.99,
			Currency:    "EUR",
			CategoryID:  "300",
			Quantity:    10,
			SKU:         stringPtr("UPD-TEST-001"),
		}

		created, err := repo.CreateListing(ctx, input)
		require.NoError(t, err)

		// Update
		update := &domain.UpdateListingInput{
			Title:    stringPtr("Updated Product"),
			Price:    float64Ptr(349.99),
			Quantity: int32Ptr(15),
		}

		updated, err := repo.UpdateListing(ctx, created.ID, update)
		require.NoError(t, err)
		assert.Equal(t, "Updated Product", updated.Title)
		assert.Equal(t, 349.99, updated.Price)
		assert.Equal(t, int32(15), updated.Quantity)

		// Delete (soft delete)
		err = repo.DeleteListing(ctx, created.ID)
		require.NoError(t, err)

		// Verify deletion
		_, err = repo.GetListingByID(ctx, created.ID)
		assert.Error(t, err)
	})

	t.Run("List with Filters", func(t *testing.T) {
		// Clean database
		tests.CleanupTestDB(t, testDB.DB)

		// Create multiple listings
		for i := 0; i < 10; i++ {
			input := &domain.CreateListingInput{
				UserID:      int64(i%3 + 1),
				Title:       "Filter Test Product",
				Description: stringPtr("Filter Test Description"),
				Price:       float64(100 * (i + 1)),
				Currency:    "USD",
				CategoryID:  fmt.Sprintf("%d", i%2+1),
				Quantity:    int32(i + 1),
				SKU:         stringPtr("FILT-TEST"),
			}
			_, err := repo.CreateListing(ctx, input)
			require.NoError(t, err)
		}

		// Test pagination
		filter := &domain.ListListingsFilter{
			Limit:  5,
			Offset: 0,
		}

		listings, total, err := repo.ListListings(ctx, filter)
		require.NoError(t, err)
		assert.Len(t, listings, 5)
		assert.Equal(t, int32(10), total)

		// Test user filter
		userID := int64(1)
		filterByUser := &domain.ListListingsFilter{
			UserID: &userID,
			Limit:  10,
			Offset: 0,
		}

		userListings, _, err := repo.ListListings(ctx, filterByUser)
		require.NoError(t, err)
		for _, listing := range userListings {
			assert.Equal(t, int64(1), listing.UserID)
		}
	})

	t.Run("Concurrent Operations", func(t *testing.T) {
		// Test concurrent reads
		input := &domain.CreateListingInput{
			UserID:      5,
			Title:       "Concurrent Test Product",
			Description: stringPtr("Concurrent Test Description"),
			Price:       499.99,
			Currency:    "USD",
			CategoryID:  "500",
			Quantity:    20,
			SKU:         stringPtr("CONC-TEST-001"),
		}

		created, err := repo.CreateListing(ctx, input)
		require.NoError(t, err)

		// Run 10 concurrent reads
		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func() {
				_, err := repo.GetListingByID(ctx, created.ID)
				assert.NoError(t, err)
				done <- true
			}()
		}

		// Wait for all goroutines
		for i := 0; i < 10; i++ {
			<-done
		}
	})
}

// TestHealthCheck verifies database health check
func TestHealthCheck(t *testing.T) {
	tests.SkipIfNoDocker(t)

	testDB := tests.SetupTestPostgres(t)
	defer testDB.TeardownTestPostgres(t)

	tests.RunMigrations(t, testDB.DB, "../../migrations")

	db := sqlx.NewDb(testDB.DB, "postgres")
	logger := zerolog.Nop()
	repo := postgres.NewRepository(db, logger)

	ctx := context.Background()

	err := repo.HealthCheck(ctx)
	assert.NoError(t, err)
}
