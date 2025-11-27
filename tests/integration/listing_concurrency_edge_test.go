//go:build integration

package integration

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/vondi-global/listings/api/proto/listings/v1"
	"github.com/vondi-global/listings/tests"
)

// ============================================================================
// Concurrency Edge Cases - Race Conditions, Deadlocks, Parallel Operations
// ============================================================================

// TestConcurrency_ParallelCreates tests creating multiple listings concurrently
func TestConcurrency_ParallelCreates(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	concurrency := 50
	var wg sync.WaitGroup
	var successCount int32
	var errorCount int32
	var mu sync.Mutex
	errors := make([]error, 0)

	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func(index int) {
			defer wg.Done()

			req := &pb.CreateListingRequest{
				UserId:     100,
				Title:      "Concurrent Test Product",
				Price:      99.99,
				Currency:   "USD",
				CategoryId: 1,
				Quantity:   1,
			}

			resp, err := client.CreateListing(ctx, req)
			if err != nil {
				atomic.AddInt32(&errorCount, 1)
				mu.Lock()
				errors = append(errors, err)
				mu.Unlock()
			} else {
				atomic.AddInt32(&successCount, 1)
				require.NotNil(t, resp)
			}
		}(i)
	}

	wg.Wait()

	t.Logf("Concurrent creates: %d success, %d errors", successCount, errorCount)
	assert.Greater(t, successCount, int32(0), "At least some creates should succeed")

	// Most should succeed
	assert.Greater(t, successCount, int32(concurrency/2), "Most concurrent creates should succeed")

	if errorCount > 0 {
		t.Logf("Sample errors: %v", errors[:min(5, len(errors))])
	}
}

// TestConcurrency_ParallelUpdates tests updating same listing concurrently
func TestConcurrency_ParallelUpdates(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	// Create initial listing
	createReq := &pb.CreateListingRequest{
		UserId:     100,
		Title:      "Original Title",
		Price:      100.00,
		Currency:   "USD",
		CategoryId: 1,
		Quantity:   100,
	}

	createResp, err := client.CreateListing(ctx, createReq)
	require.NoError(t, err)
	require.NotNil(t, createResp)

	listingID := createResp.Listing.Id

	// Update concurrently
	concurrency := 20
	var wg sync.WaitGroup
	var successCount int32
	var errorCount int32

	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func(index int) {
			defer wg.Done()

			updateReq := &pb.UpdateListingRequest{
				Id:       listingID,
				Quantity: int32Ptr(int32(200 + index)),
			}

			_, err := client.UpdateListing(ctx, updateReq)
			if err != nil {
				atomic.AddInt32(&errorCount, 1)
			} else {
				atomic.AddInt32(&successCount, 1)
			}
		}(i)
	}

	wg.Wait()

	t.Logf("Concurrent updates: %d success, %d errors", successCount, errorCount)

	// Verify final state (should be consistent)
	finalResp, err := client.GetListing(ctx, &pb.GetListingRequest{Id: listingID})
	require.NoError(t, err)
	require.NotNil(t, finalResp)

	// Quantity should be one of the update values (200-219)
	assert.GreaterOrEqual(t, finalResp.Listing.Quantity, int32(200))
	assert.LessOrEqual(t, finalResp.Listing.Quantity, int32(219))
	t.Logf("Final quantity: %d", finalResp.Listing.Quantity)
}

// TestConcurrency_ReadWriteRace tests concurrent reads and writes
func TestConcurrency_ReadWriteRace(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	// Create initial listing
	createReq := &pb.CreateListingRequest{
		UserId:     100,
		Title:      "Race Test Listing",
		Price:      100.00,
		Currency:   "USD",
		CategoryId: 1,
		Quantity:   100,
	}

	createResp, err := client.CreateListing(ctx, createReq)
	require.NoError(t, err)
	listingID := createResp.Listing.Id

	// Start concurrent operations
	duration := 2 * time.Second
	stopChan := make(chan struct{})
	var readCount int32
	var writeCount int32
	var wg sync.WaitGroup

	// Readers (10 goroutines)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stopChan:
					return
				default:
					_, err := client.GetListing(ctx, &pb.GetListingRequest{Id: listingID})
					if err == nil {
						atomic.AddInt32(&readCount, 1)
					}
					time.Sleep(10 * time.Millisecond)
				}
			}
		}()
	}

	// Writers (3 goroutines)
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			counter := 0
			for {
				select {
				case <-stopChan:
					return
				default:
					updateReq := &pb.UpdateListingRequest{
						Id:    listingID,
						Price: float64Ptr(100.0 + float64(counter)),
					}
					_, err := client.UpdateListing(ctx, updateReq)
					if err == nil {
						atomic.AddInt32(&writeCount, 1)
					}
					counter++
					time.Sleep(50 * time.Millisecond)
				}
			}
		}(i)
	}

	// Let it run for duration
	time.Sleep(duration)
	close(stopChan)
	wg.Wait()

	t.Logf("Read/Write race: %d reads, %d writes", readCount, writeCount)
	assert.Greater(t, readCount, int32(0), "Should have successful reads")
	assert.Greater(t, writeCount, int32(0), "Should have successful writes")

	// Verify final state is consistent
	finalResp, err := client.GetListing(ctx, &pb.GetListingRequest{Id: listingID})
	require.NoError(t, err)
	require.NotNil(t, finalResp)
	assert.Equal(t, listingID, finalResp.Listing.Id)
}

// TestConcurrency_CreateAndDelete tests concurrent create and delete operations
func TestConcurrency_CreateAndDelete(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	concurrency := 30
	var wg sync.WaitGroup
	var createCount int32
	var deleteCount int32
	var createdIDs sync.Map

	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func(index int) {
			defer wg.Done()

			// Create
			createReq := &pb.CreateListingRequest{
				UserId:     100,
				Title:      "Create/Delete Test",
				Price:      99.99,
				Currency:   "USD",
				CategoryId: 1,
				Quantity:   1,
			}

			createResp, err := client.CreateListing(ctx, createReq)
			if err == nil && createResp != nil {
				atomic.AddInt32(&createCount, 1)
				createdIDs.Store(index, createResp.Listing.Id)

				// Immediately try to delete
				deleteReq := &pb.DeleteListingRequest{
					Id: createResp.Listing.Id,
				}

				_, err := client.DeleteListing(ctx, deleteReq)
				if err == nil {
					atomic.AddInt32(&deleteCount, 1)
				}
			}
		}(i)
	}

	wg.Wait()

	t.Logf("Create/Delete concurrency: %d created, %d deleted", createCount, deleteCount)
	assert.Greater(t, createCount, int32(0), "Should create some listings")
	assert.Greater(t, deleteCount, int32(0), "Should delete some listings")
}

// TestConcurrency_BulkOperations tests concurrent bulk operations
func TestConcurrency_BulkOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping bulk concurrency test in short mode")
	}

	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	// Create initial listings
	listingIDs := make([]int64, 100)
	for i := 0; i < 100; i++ {
		createReq := &pb.CreateListingRequest{
			UserId:     100,
			Title:      "Bulk Test Listing",
			Price:      99.99,
			Currency:   "USD",
			CategoryId: 1,
			Quantity:   10,
		}

		createResp, err := client.CreateListing(ctx, createReq)
		require.NoError(t, err)
		listingIDs[i] = createResp.Listing.Id
	}

	// Perform concurrent bulk operations
	var wg sync.WaitGroup
	var successCount int32

	// Concurrent bulk gets
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()

			listReq := &pb.ListListingsRequest{
				Limit:  50,
				Offset: 0,
			}

			_, err := client.ListListings(ctx, listReq)
			if err == nil {
				atomic.AddInt32(&successCount, 1)
			}
		}()
	}

	wg.Wait()

	t.Logf("Bulk operations: %d successful list requests", successCount)
	assert.Greater(t, successCount, int32(0), "Should complete some bulk operations")
}

// ============================================================================
// Deadlock Detection Tests
// ============================================================================

// TestConcurrency_NoDeadlock tests that concurrent operations don't cause deadlock
func TestConcurrency_NoDeadlock(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create listings
	listingIDs := make([]int64, 10)
	for i := 0; i < 10; i++ {
		createReq := &pb.CreateListingRequest{
			UserId:     int64(100 + i),
			Title:      "Deadlock Test",
			Price:      99.99,
			Currency:   "USD",
			CategoryId: 1,
			Quantity:   10,
		}

		createResp, err := client.CreateListing(ctx, createReq)
		require.NoError(t, err)
		listingIDs[i] = createResp.Listing.Id
	}

	// Perform mixed operations concurrently
	var wg sync.WaitGroup
	operations := []func(int64){
		func(id int64) {
			client.GetListing(ctx, &pb.GetListingRequest{Id: id})
		},
		func(id int64) {
			client.UpdateListing(ctx, &pb.UpdateListingRequest{
				Id:    id,
				Price: float64Ptr(199.99),
			})
		},
		func(id int64) {
			client.ListListings(ctx, &pb.ListListingsRequest{Limit: 10})
		},
	}

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(index int) {
			defer wg.Done()

			// Randomly pick listing and operation
			listingID := listingIDs[index%len(listingIDs)]
			operation := operations[index%len(operations)]

			operation(listingID)
		}(i)
	}

	// Wait with timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		t.Log("All operations completed without deadlock")
	case <-ctx.Done():
		t.Fatal("Deadlock detected: operations did not complete within timeout")
	}
}

// TestConcurrency_DoubleDelete tests concurrent deletes of same listing
func TestConcurrency_DoubleDelete(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	// Create listing
	createReq := &pb.CreateListingRequest{
		UserId:     100,
		Title:      "Double Delete Test",
		Price:      99.99,
		Currency:   "USD",
		CategoryId: 1,
		Quantity:   1,
	}

	createResp, err := client.CreateListing(ctx, createReq)
	require.NoError(t, err)
	listingID := createResp.Listing.Id

	// Try to delete concurrently
	concurrency := 10
	var wg sync.WaitGroup
	var successCount int32
	var notFoundCount int32
	var otherErrorCount int32

	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()

			deleteReq := &pb.DeleteListingRequest{
				Id: listingID,
			}

			_, err := client.DeleteListing(ctx, deleteReq)
			if err == nil {
				atomic.AddInt32(&successCount, 1)
			} else {
				st, ok := status.FromError(err)
				if ok && st.Code() == codes.NotFound {
					atomic.AddInt32(&notFoundCount, 1)
				} else {
					atomic.AddInt32(&otherErrorCount, 1)
				}
			}
		}()
	}

	wg.Wait()

	t.Logf("Double delete: %d success, %d not_found, %d other_errors",
		successCount, notFoundCount, otherErrorCount)

	// Only one should succeed (or all might fail if soft delete is idempotent)
	assert.LessOrEqual(t, successCount, int32(1), "At most one delete should succeed")
	assert.Greater(t, notFoundCount+successCount, int32(0), "Should have some responses")
}

// TestConcurrency_ParallelSearches tests concurrent search operations
func TestConcurrency_ParallelSearches(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	// Create test listings
	for i := 0; i < 50; i++ {
		createReq := &pb.CreateListingRequest{
			UserId:     100,
			Title:      "Search Test Product",
			Price:      99.99,
			Currency:   "USD",
			CategoryId: int64(1 + (i % 3)),
			Quantity:   1,
		}

		_, err := client.CreateListing(ctx, createReq)
		require.NoError(t, err)
	}

	// Perform concurrent searches
	concurrency := 20
	var wg sync.WaitGroup
	var successCount int32
	var errorCount int32

	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func(index int) {
			defer wg.Done()

			listReq := &pb.ListListingsRequest{
				Limit:  10,
				Offset: int32(index % 5 * 10),
			}

			resp, err := client.ListListings(ctx, listReq)
			if err != nil {
				atomic.AddInt32(&errorCount, 1)
			} else {
				atomic.AddInt32(&successCount, 1)
				require.NotNil(t, resp)
			}
		}(i)
	}

	wg.Wait()

	t.Logf("Parallel searches: %d success, %d errors", successCount, errorCount)
	assert.Equal(t, int32(concurrency), successCount, "All searches should succeed")
	assert.Equal(t, int32(0), errorCount, "No search errors expected")
}

// ============================================================================
// Helper Functions
// ============================================================================

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
