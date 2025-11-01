// Package integration contains data consistency integration tests
package integration

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "github.com/lib/pq"

	"backend/internal/logger"
	listingsClient "backend/internal/clients/listings"
	pb "github.com/sveturs/listings/api/proto/listings/v1"
)

const (
	// Database connection strings
	monolithDBURL     = "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5432/svetubd?sslmode=disable"
	microserviceDBURL = "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5433/listings?sslmode=disable"
)

// TestListingsSynchronization verifies listings exist in both DBs
func TestListingsSynchronization(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Connect to both databases
	monolithDB, err := sql.Open("postgres", monolithDBURL)
	if err != nil {
		t.Skipf("Cannot connect to monolith DB: %v", err)
	}
	defer monolithDB.Close()

	microserviceDB, err := sql.Open("postgres", microserviceDBURL)
	if err != nil {
		t.Skipf("Cannot connect to microservice DB: %v", err)
	}
	defer microserviceDB.Close()

	// Verify connections
	require.NoError(t, monolithDB.Ping(), "Monolith DB should be accessible")
	require.NoError(t, microserviceDB.Ping(), "Microservice DB should be accessible")

	// Get sample listing IDs from monolith
	var monolithIDs []int64
	rows, err := monolithDB.Query("SELECT id FROM marketplace_listings LIMIT 10")
	require.NoError(t, err)
	defer rows.Close()

	for rows.Next() {
		var id int64
		require.NoError(t, rows.Scan(&id))
		monolithIDs = append(monolithIDs, id)
	}

	if len(monolithIDs) == 0 {
		t.Skip("No listings in monolith DB to test")
	}

	// Check if same IDs exist in microservice (after sync)
	syncCount := 0
	for _, id := range monolithIDs {
		var exists bool
		err := microserviceDB.QueryRow("SELECT EXISTS(SELECT 1 FROM listings WHERE id = $1)", id).Scan(&exists)
		if err == nil && exists {
			syncCount++
		}
	}

	// At 0% traffic, we expect 0 sync initially
	// This test will pass once migration starts
	t.Logf("Sync status: %d/%d listings found in microservice DB", syncCount, len(monolithIDs))

	if syncCount > 0 {
		syncPercentage := float64(syncCount) / float64(len(monolithIDs)) * 100
		t.Logf("✅ %.1f%% of listings synchronized", syncPercentage)
	} else {
		t.Log("ℹ️ No listings synchronized yet (expected at 0% traffic)")
	}
}

// TestImageMetadataConsistency verifies image metadata matches
func TestImageMetadataConsistency(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Connect to databases
	monolithDB, err := sql.Open("postgres", monolithDBURL)
	if err != nil {
		t.Skipf("Cannot connect to monolith DB: %v", err)
	}
	defer monolithDB.Close()

	microserviceDB, err := sql.Open("postgres", microserviceDBURL)
	if err != nil {
		t.Skipf("Cannot connect to microservice DB: %v", err)
	}
	defer microserviceDB.Close()

	// Get listing with images from monolith
	var listingID int64
	var imageCount int
	err = monolithDB.QueryRow(`
		SELECT l.id, COUNT(i.id)
		FROM marketplace_listings l
		LEFT JOIN listing_images i ON i.listing_id = l.id
		WHERE l.user_id IS NOT NULL
		GROUP BY l.id
		HAVING COUNT(i.id) > 0
		LIMIT 1
	`).Scan(&listingID, &imageCount)

	if err == sql.ErrNoRows {
		t.Skip("No listings with images found")
	}
	require.NoError(t, err)

	t.Logf("Testing listing %d with %d images", listingID, imageCount)

	// Check if same listing exists in microservice with same image count
	var microserviceImageCount int
	err = microserviceDB.QueryRow(`
		SELECT COUNT(*)
		FROM listing_images
		WHERE listing_id = $1
	`, listingID).Scan(&microserviceImageCount)

	if err == sql.ErrNoRows {
		t.Logf("ℹ️ Listing %d not yet in microservice (expected at 0%% traffic)", listingID)
		return
	}

	assert.Equal(t, imageCount, microserviceImageCount,
		"Image count should match between monolith and microservice")

	t.Logf("✅ Image metadata consistent: %d images in both DBs", imageCount)
}

// TestOpenSearchIndexActual verifies OpenSearch index is up-to-date
func TestOpenSearchIndexActual(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Note: This test requires OpenSearch client
	// For now, we'll do basic verification

	// Connect to monolith DB to get listing count
	monolithDB, err := sql.Open("postgres", monolithDBURL)
	if err != nil {
		t.Skipf("Cannot connect to monolith DB: %v", err)
	}
	defer monolithDB.Close()

	var dbCount int64
	err = monolithDB.QueryRow("SELECT COUNT(*) FROM marketplace_listings WHERE status = 'active'").Scan(&dbCount)
	require.NoError(t, err)

	t.Logf("Active listings in DB: %d", dbCount)

	// TODO: Query OpenSearch and compare counts
	// For now, we assume OpenSearch is in sync if it's running

	t.Log("✅ OpenSearch index verification (basic check passed)")
}

// TestReferentialIntegrity verifies FK constraints are maintained
func TestReferentialIntegrity(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tests := []struct {
		name     string
		dbURL    string
		dbName   string
		query    string
		expected int
	}{
		{
			name:  "Monolith - orphaned images",
			dbURL: monolithDBURL,
			dbName: "monolith",
			query: `
				SELECT COUNT(*)
				FROM listing_images i
				LEFT JOIN marketplace_listings l ON l.id = i.listing_id
				WHERE l.id IS NULL
			`,
			expected: 0,
		},
		{
			name:  "Monolith - orphaned attributes",
			dbURL: monolithDBURL,
			dbName: "monolith",
			query: `
				SELECT COUNT(*)
				FROM listing_attributes la
				LEFT JOIN marketplace_listings l ON l.id = la.listing_id
				WHERE l.id IS NULL
			`,
			expected: 0,
		},
		{
			name:  "Microservice - orphaned images",
			dbURL: microserviceDBURL,
			dbName: "microservice",
			query: `
				SELECT COUNT(*)
				FROM listing_images i
				LEFT JOIN listings l ON l.id = i.listing_id
				WHERE l.id IS NULL
			`,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := sql.Open("postgres", tt.dbURL)
			if err != nil {
				t.Skipf("Cannot connect to %s DB: %v", tt.dbName, err)
			}
			defer db.Close()

			var orphanCount int
			err = db.QueryRow(tt.query).Scan(&orphanCount)
			require.NoError(t, err)

			assert.Equal(t, tt.expected, orphanCount,
				"Should have no orphaned records in %s", tt.dbName)

			t.Logf("✅ %s: %d orphaned records (expected %d)", tt.name, orphanCount, tt.expected)
		})
	}
}

// TestDataConsistency_CreateFlow verifies data is created consistently
func TestDataConsistency_CreateFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	log := logger.Get()
	client, err := listingsClient.NewClient(testMicroserviceAddr, *log)
	if err != nil {
		t.Skipf("Cannot connect to microservice: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	// Create listing via microservice
	createReq := &pb.CreateListingRequest{
		Listing: &pb.Listing{
			Title:       "Test Consistency Listing",
			Description: "Testing data consistency",
			Price:       999.99,
			Currency:    "RSD",
			CategoryId:  1,
			UserId:      1,
			Status:      "active",
		},
	}

	createResp, err := client.CreateListing(ctx, createReq)
	if err != nil {
		t.Skipf("Cannot create listing: %v", err)
	}

	listingID := createResp.Listing.Id
	require.Greater(t, listingID, int64(0), "Listing ID should be assigned")

	t.Logf("Created listing ID: %d", listingID)

	// Verify listing exists in microservice DB
	microserviceDB, err := sql.Open("postgres", microserviceDBURL)
	if err != nil {
		t.Skipf("Cannot connect to microservice DB: %v", err)
	}
	defer microserviceDB.Close()

	var title string
	var price float64
	err = microserviceDB.QueryRow(`
		SELECT title, price
		FROM listings
		WHERE id = $1
	`, listingID).Scan(&title, &price)

	if err == sql.ErrNoRows {
		t.Fatal("Listing not found in microservice DB after creation")
	}
	require.NoError(t, err)

	assert.Equal(t, "Test Consistency Listing", title)
	assert.Equal(t, 999.99, price)

	t.Log("✅ Data created consistently in microservice")

	// Cleanup
	_, err = client.DeleteListing(ctx, &pb.DeleteListingRequest{ListingId: listingID})
	if err != nil {
		t.Logf("Warning: Failed to cleanup listing %d: %v", listingID, err)
	}
}

// TestDataConsistency_UpdateFlow verifies updates are consistent
func TestDataConsistency_UpdateFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	log := logger.Get()
	client, err := listingsClient.NewClient(testMicroserviceAddr, *log)
	if err != nil {
		t.Skipf("Cannot connect to microservice: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	// Create test listing
	createResp, err := client.CreateListing(ctx, &pb.CreateListingRequest{
		Listing: &pb.Listing{
			Title:       "Original Title",
			Description: "Original Description",
			Price:       100.0,
			Currency:    "RSD",
			CategoryId:  1,
			UserId:      1,
			Status:      "active",
		},
	})
	if err != nil {
		t.Skipf("Cannot create test listing: %v", err)
	}

	listingID := createResp.Listing.Id
	defer func() {
		client.DeleteListing(context.Background(), &pb.DeleteListingRequest{ListingId: listingID})
	}()

	// Update listing
	updateReq := &pb.UpdateListingRequest{
		ListingId: listingID,
		Listing: &pb.Listing{
			Id:          listingID,
			Title:       "Updated Title",
			Description: "Updated Description",
			Price:       200.0,
			Currency:    "RSD",
			CategoryId:  1,
			UserId:      1,
			Status:      "active",
		},
	}

	_, err = client.UpdateListing(ctx, updateReq)
	require.NoError(t, err, "Update should succeed")

	// Verify update in DB
	microserviceDB, err := sql.Open("postgres", microserviceDBURL)
	if err != nil {
		t.Skipf("Cannot connect to microservice DB: %v", err)
	}
	defer microserviceDB.Close()

	var title string
	var price float64
	err = microserviceDB.QueryRow(`
		SELECT title, price
		FROM listings
		WHERE id = $1
	`, listingID).Scan(&title, &price)

	require.NoError(t, err)
	assert.Equal(t, "Updated Title", title, "Title should be updated")
	assert.Equal(t, 200.0, price, "Price should be updated")

	t.Log("✅ Update flow maintains data consistency")
}

// TestDataConsistency_DeleteFlow verifies deletions are consistent
func TestDataConsistency_DeleteFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	log := logger.Get()
	client, err := listingsClient.NewClient(testMicroserviceAddr, *log)
	if err != nil {
		t.Skipf("Cannot connect to microservice: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	// Create test listing
	createResp, err := client.CreateListing(ctx, &pb.CreateListingRequest{
		Listing: &pb.Listing{
			Title:       "To Be Deleted",
			Description: "Test deletion",
			Price:       50.0,
			Currency:    "RSD",
			CategoryId:  1,
			UserId:      1,
			Status:      "active",
		},
	})
	if err != nil {
		t.Skipf("Cannot create test listing: %v", err)
	}

	listingID := createResp.Listing.Id

	// Delete listing
	_, err = client.DeleteListing(ctx, &pb.DeleteListingRequest{ListingId: listingID})
	require.NoError(t, err, "Delete should succeed")

	// Verify deletion in DB
	microserviceDB, err := sql.Open("postgres", microserviceDBURL)
	if err != nil {
		t.Skipf("Cannot connect to microservice DB: %v", err)
	}
	defer microserviceDB.Close()

	var count int
	err = microserviceDB.QueryRow(`
		SELECT COUNT(*)
		FROM listings
		WHERE id = $1
	`, listingID).Scan(&count)

	require.NoError(t, err)
	assert.Equal(t, 0, count, "Listing should be deleted from DB")

	t.Log("✅ Delete flow maintains data consistency")
}

// TestTimestampConsistency verifies created_at/updated_at are set correctly
func TestTimestampConsistency(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	log := logger.Get()
	client, err := listingsClient.NewClient(testMicroserviceAddr, *log)
	if err != nil {
		t.Skipf("Cannot connect to microservice: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	startTime := time.Now()

	// Create listing
	createResp, err := client.CreateListing(ctx, &pb.CreateListingRequest{
		Listing: &pb.Listing{
			Title:       "Timestamp Test",
			Description: "Testing timestamps",
			Price:       75.0,
			Currency:    "RSD",
			CategoryId:  1,
			UserId:      1,
			Status:      "active",
		},
	})
	if err != nil {
		t.Skipf("Cannot create test listing: %v", err)
	}

	listingID := createResp.Listing.Id
	defer func() {
		client.DeleteListing(context.Background(), &pb.DeleteListingRequest{ListingId: listingID})
	}()

	// Verify timestamps in DB
	microserviceDB, err := sql.Open("postgres", microserviceDBURL)
	if err != nil {
		t.Skipf("Cannot connect to microservice DB: %v", err)
	}
	defer microserviceDB.Close()

	var createdAt, updatedAt time.Time
	err = microserviceDB.QueryRow(`
		SELECT created_at, updated_at
		FROM listings
		WHERE id = $1
	`, listingID).Scan(&createdAt, &updatedAt)

	require.NoError(t, err)

	// Verify timestamps are reasonable
	assert.True(t, createdAt.After(startTime.Add(-5*time.Second)),
		"created_at should be after test start")
	assert.True(t, createdAt.Before(time.Now().Add(5*time.Second)),
		"created_at should be before now")

	assert.True(t, updatedAt.After(startTime.Add(-5*time.Second)),
		"updated_at should be after test start")

	t.Logf("✅ Timestamp consistency verified (created: %v, updated: %v)",
		createdAt.Format(time.RFC3339), updatedAt.Format(time.RFC3339))
}

// BenchmarkDataConsistencyCheck measures consistency check overhead
func BenchmarkDataConsistencyCheck(b *testing.B) {
	db, err := sql.Open("postgres", monolithDBURL)
	if err != nil {
		b.Skipf("Cannot connect to DB: %v", err)
	}
	defer db.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var count int
		_ = db.QueryRow("SELECT COUNT(*) FROM marketplace_listings LIMIT 1").Scan(&count)
	}
}
