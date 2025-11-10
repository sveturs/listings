package integration

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/sveturs/listings/api/proto/listings/v1"
	testutils "github.com/sveturs/listings/internal/testing"
)

// =============================================================================
// 1. CreateListing Tests (8 scenarios)
// =============================================================================

func TestCreateListing(t *testing.T) {
	testutils.SkipIfShort(t)
	testutils.SkipIfNoDocker(t)

	t.Run("ValidC2CListing_Success", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		// Setup: Insert test category
		ExecuteSQL(t, server, `
			INSERT INTO c2c_categories (id, name, slug, parent_id, sort_order, level, is_active, count)
			VALUES ($1, $2, $3, NULL, $4, 0, $5, $6)
		`, 1, "Electronics", "electronics", 1, true, 0)

		ctx := testutils.TestContext(t)
		req := &pb.CreateListingRequest{
			UserId:      100,
			Title:       "Test Electronics Listing",
			Description: testutils.StringPtr("A test electronic device"),
			Price:       99.99,
			Currency:    "USD",
			CategoryId:  1,
			Quantity:    1,
		}

		resp, err := server.Client.CreateListing(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.NotZero(t, resp.Listing.Id)
		assert.Equal(t, "Test Electronics Listing", resp.Listing.Title)
		assert.Equal(t, 99.99, resp.Listing.Price)
		assert.Equal(t, "USD", resp.Listing.Currency)
		assert.Equal(t, int64(1), resp.Listing.CategoryId)
		assert.Equal(t, "draft", resp.Listing.Status) // Default status is "draft"
		assert.NotEmpty(t, resp.Listing.Uuid)
		assert.False(t, resp.Listing.IsDeleted)
	})

	t.Run("ValidB2CListing_WithStorefront_Success", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		// Setup: Insert category and storefront
		ExecuteSQL(t, server, `
			INSERT INTO c2c_categories (id, name, slug, parent_id, sort_order, level, is_active, count)
			VALUES ($1, $2, $3, NULL, $4, 0, $5, $6)
		`, 2, "Fashion", "fashion", 1, true, 0)

		ExecuteSQL(t, server, `
			INSERT INTO storefronts (id, user_id, slug, name, country, is_active, is_verified)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, 5001, 200, "test-store", "Test Store", "US", true, true)

		ctx := testutils.TestContext(t)
		storefrontID := int64(5001)
		req := &pb.CreateListingRequest{
			UserId:       200,
			StorefrontId: &storefrontID,
			Title:        "Fashion Item from Store",
			Description:  testutils.StringPtr("A fashion item"),
			Price:        49.99,
			Currency:     "USD",
			CategoryId:   2,
			Quantity:     5,
		}

		resp, err := server.Client.CreateListing(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.NotZero(t, resp.Listing.Id)
		assert.Equal(t, "Fashion Item from Store", resp.Listing.Title)
		require.NotNil(t, resp.Listing.StorefrontId)
		assert.Equal(t, int64(5001), *resp.Listing.StorefrontId)
	})

	t.Run("AllOptionalFields_Success", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		// Setup: category
		ExecuteSQL(t, server, `
			INSERT INTO c2c_categories (id, name, slug, parent_id, sort_order, level, is_active, count)
			VALUES ($1, $2, $3, NULL, $4, 0, $5, $6)
		`, 3, "Vehicles", "vehicles", 1, true, 0)

		ctx := testutils.TestContext(t)
		sku := "TEST-SKU-001"
		req := &pb.CreateListingRequest{
			UserId:      300,
			Title:       "Vehicle with All Fields",
			Description: testutils.StringPtr("Complete vehicle listing with all optional fields"),
			Price:       25000.00,
			Currency:    "USD",
			CategoryId:  3,
			Quantity:    1,
			Sku:         &sku,
		}

		resp, err := server.Client.CreateListing(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.NotZero(t, resp.Listing.Id)
		require.NotNil(t, resp.Listing.Sku)
		assert.Equal(t, "TEST-SKU-001", *resp.Listing.Sku)
	})

	t.Run("MinimalRequiredFields_Success", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		// Setup: category
		ExecuteSQL(t, server, `
			INSERT INTO c2c_categories (id, name, slug, parent_id, sort_order, level, is_active, count)
			VALUES ($1, $2, $3, NULL, $4, 0, $5, $6)
		`, 4, "Books", "books", 1, true, 0)

		ctx := testutils.TestContext(t)
		req := &pb.CreateListingRequest{
			UserId:     400,
			Title:      "Minimal Listing",
			Price:      9.99,
			Currency:   "USD",
			CategoryId: 4,
			Quantity:   1,
		}

		resp, err := server.Client.CreateListing(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.NotZero(t, resp.Listing.Id)
		assert.Equal(t, "Minimal Listing", resp.Listing.Title)
		assert.Nil(t, resp.Listing.Description) // Optional field not provided
	})

	t.Run("InvalidInput_MissingTitle_Error", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		ctx := testutils.TestContext(t)
		req := &pb.CreateListingRequest{
			UserId:     500,
			Title:      "", // Empty title (invalid)
			Price:      10.00,
			Currency:   "USD",
			CategoryId: 1,
			Quantity:   1,
		}

		resp, err := server.Client.CreateListing(ctx, req)

		require.Error(t, err)
		assert.Nil(t, resp)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("InvalidInput_NegativePrice_Error", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		ctx := testutils.TestContext(t)
		req := &pb.CreateListingRequest{
			UserId:     500,
			Title:      "Invalid Price Listing",
			Price:      -10.00, // Negative price (invalid)
			Currency:   "USD",
			CategoryId: 1,
			Quantity:   1,
		}

		resp, err := server.Client.CreateListing(ctx, req)

		require.Error(t, err)
		assert.Nil(t, resp)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("InvalidInput_NonExistentCategory_Error", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		ctx := testutils.TestContext(t)
		req := &pb.CreateListingRequest{
			UserId:     500,
			Title:      "Listing with Invalid Category",
			Price:      10.00,
			Currency:   "USD",
			CategoryId: 99999, // Non-existent category
			Quantity:   1,
		}

		resp, err := server.Client.CreateListing(ctx, req)

		require.Error(t, err)
		assert.Nil(t, resp)
		st, ok := status.FromError(err)
		require.True(t, ok)
		// Could be InvalidArgument or FailedPrecondition depending on implementation
		assert.Contains(t, []codes.Code{codes.InvalidArgument, codes.FailedPrecondition, codes.NotFound}, st.Code())
	})

	t.Run("ConcurrentCreation_MultipleListings_Success", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		// Setup: category
		ExecuteSQL(t, server, `
			INSERT INTO c2c_categories (id, name, slug, parent_id, sort_order, level, is_active, count)
			VALUES ($1, $2, $3, NULL, $4, 0, $5, $6)
		`, 5, "Concurrent Test", "concurrent-test", 1, true, 0)

		ctx := testutils.TestContext(t)
		concurrency := 5
		var wg sync.WaitGroup
		errors := make(chan error, concurrency)
		successes := make(chan *pb.CreateListingResponse, concurrency)

		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()

				req := &pb.CreateListingRequest{
					UserId:     int64(600 + index),
					Title:      fmt.Sprintf("Concurrent Listing %d", index),
					Price:      float64(10 + index),
					Currency:   "USD",
					CategoryId: 5,
					Quantity:   1,
				}

				resp, err := server.Client.CreateListing(ctx, req)
				if err != nil {
					errors <- err
				} else {
					successes <- resp
				}
			}(i)
		}

		wg.Wait()
		close(errors)
		close(successes)

		// Verify all succeeded
		var successCount int
		for resp := range successes {
			assert.NotNil(t, resp)
			assert.NotZero(t, resp.Listing.Id)
			successCount++
		}

		var errorCount int
		for err := range errors {
			t.Logf("Concurrent creation error: %v", err)
			errorCount++
		}

		assert.Equal(t, concurrency, successCount, "All concurrent creates should succeed")
		assert.Equal(t, 0, errorCount, "No errors expected")
	})
}

// =============================================================================
// 2. UpdateListing Tests (7 scenarios)
// =============================================================================

func TestUpdateListing(t *testing.T) {
	testutils.SkipIfShort(t)
	testutils.SkipIfNoDocker(t)

	t.Run("UpdateTitleDescriptionPrice_Success", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		// Setup: category and listing
		ExecuteSQL(t, server, `
			INSERT INTO c2c_categories (id, name, slug, parent_id, sort_order, level, is_active, count)
			VALUES ($1, $2, $3, NULL, $4, 0, $5, $6)
		`, 10, "Update Test", "update-test", 1, true, 0)

		ExecuteSQL(t, server, `
			INSERT INTO listings (
				id, user_id, title, description, price, currency, category_id,
				status, visibility, quantity, view_count, favorites_count
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
			)
		`, 2001, 1000, "Original Title", "Original Description", 50.00, "USD", 10,
			"active", "public", 1, 0, 0)

		ctx := testutils.TestContext(t)
		newDesc := "Updated Description"
		req := &pb.UpdateListingRequest{
			Id:          2001,
			UserId:      1000,
			Title:       testutils.StringPtr("Updated Title"),
			Description: &newDesc,
			Price:       testutils.Float64Ptr(75.00),
		}

		resp, err := server.Client.UpdateListing(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, int64(2001), resp.Listing.Id)
		assert.Equal(t, "Updated Title", resp.Listing.Title)
		assert.NotNil(t, resp.Listing.Description)
		assert.Equal(t, "Updated Description", *resp.Listing.Description)
		assert.Equal(t, 75.00, resp.Listing.Price)
	})

	t.Run("PartialUpdate_OnlyQuantity_Success", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		// Setup
		ExecuteSQL(t, server, `
			INSERT INTO c2c_categories (id, name, slug, parent_id, sort_order, level, is_active, count)
			VALUES ($1, $2, $3, NULL, $4, 0, $5, $6)
		`, 11, "Partial Update", "partial-update", 1, true, 0)

		ExecuteSQL(t, server, `
			INSERT INTO listings (
				id, user_id, title, description, price, currency, category_id,
				status, visibility, quantity, view_count, favorites_count
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
			)
		`, 2002, 1001, "Partial Title", "Partial Description", 30.00, "USD", 11,
			"active", "public", 5, 0, 0)

		ctx := testutils.TestContext(t)
		req := &pb.UpdateListingRequest{
			Id:       2002,
			UserId:   1001,
			Quantity: testutils.Int32Ptr(10), // Only update quantity
		}

		resp, err := server.Client.UpdateListing(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, int32(10), resp.Listing.Quantity)
		// Other fields should remain unchanged
		assert.Equal(t, "Partial Title", resp.Listing.Title)
		assert.Equal(t, 30.00, resp.Listing.Price)
	})

	t.Run("UpdateWithValidationError_NegativeQuantity_Error", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		// Setup
		ExecuteSQL(t, server, `
			INSERT INTO c2c_categories (id, name, slug, parent_id, sort_order, level, is_active, count)
			VALUES ($1, $2, $3, NULL, $4, 0, $5, $6)
		`, 12, "Validation", "validation", 1, true, 0)

		ExecuteSQL(t, server, `
			INSERT INTO listings (
				id, user_id, title, description, price, currency, category_id,
				status, visibility, quantity, view_count, favorites_count
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
			)
		`, 2003, 1002, "Valid Listing", "Description", 20.00, "USD", 12,
			"active", "public", 5, 0, 0)

		ctx := testutils.TestContext(t)
		req := &pb.UpdateListingRequest{
			Id:       2003,
			UserId:   1002,
			Quantity: testutils.Int32Ptr(-5), // Invalid negative quantity
		}

		resp, err := server.Client.UpdateListing(ctx, req)

		require.Error(t, err)
		assert.Nil(t, resp)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})

	t.Run("UpdateNonExistentListing_NotFound", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		ctx := testutils.TestContext(t)
		req := &pb.UpdateListingRequest{
			Id:     99999, // Non-existent listing
			UserId: 1003,
			Title:  testutils.StringPtr("Should Fail"),
		}

		resp, err := server.Client.UpdateListing(ctx, req)

		require.Error(t, err)
		assert.Nil(t, resp)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
	})

	t.Run("UpdateByWrongUser_PermissionDenied", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		// Setup
		ExecuteSQL(t, server, `
			INSERT INTO c2c_categories (id, name, slug, parent_id, sort_order, level, is_active, count)
			VALUES ($1, $2, $3, NULL, $4, 0, $5, $6)
		`, 13, "Permission", "permission", 1, true, 0)

		ExecuteSQL(t, server, `
			INSERT INTO listings (
				id, user_id, title, description, price, currency, category_id,
				status, visibility, quantity, view_count, favorites_count
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
			)
		`, 2004, 1004, "Owner's Listing", "Description", 15.00, "USD", 13,
			"active", "public", 1, 0, 0)

		ctx := testutils.TestContext(t)
		req := &pb.UpdateListingRequest{
			Id:     2004,
			UserId: 9999, // Different user (not owner)
			Title:  testutils.StringPtr("Unauthorized Update"),
		}

		resp, err := server.Client.UpdateListing(ctx, req)

		require.Error(t, err)
		assert.Nil(t, resp)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.PermissionDenied, st.Code())
	})

	t.Run("ConcurrentUpdate_SameListing_Success", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		// Setup
		ExecuteSQL(t, server, `
			INSERT INTO c2c_categories (id, name, slug, parent_id, sort_order, level, is_active, count)
			VALUES ($1, $2, $3, NULL, $4, 0, $5, $6)
		`, 14, "Concurrent Update", "concurrent-update", 1, true, 0)

		ExecuteSQL(t, server, `
			INSERT INTO listings (
				id, user_id, title, description, price, currency, category_id,
				status, visibility, quantity, view_count, favorites_count
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
			)
		`, 2005, 1005, "Concurrent Test", "Description", 100.00, "USD", 14,
			"active", "public", 10, 0, 0)

		ctx := testutils.TestContext(t)
		concurrency := 5
		var wg sync.WaitGroup
		errors := make(chan error, concurrency)

		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()

				req := &pb.UpdateListingRequest{
					Id:     2005,
					UserId: 1005,
					Price:  testutils.Float64Ptr(100.00 + float64(index)),
				}

				_, err := server.Client.UpdateListing(ctx, req)
				if err != nil {
					errors <- err
				}
			}(i)
		}

		wg.Wait()
		close(errors)

		// At least some updates should succeed (last-write-wins)
		errorCount := 0
		for err := range errors {
			t.Logf("Concurrent update error: %v", err)
			errorCount++
		}

		// All updates should succeed or handle optimistic locking
		assert.LessOrEqual(t, errorCount, concurrency/2, "Most concurrent updates should succeed")
	})

	t.Run("StatusTransition_DraftToActive_Success", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		// Setup
		ExecuteSQL(t, server, `
			INSERT INTO c2c_categories (id, name, slug, parent_id, sort_order, level, is_active, count)
			VALUES ($1, $2, $3, NULL, $4, 0, $5, $6)
		`, 15, "Status", "status", 1, true, 0)

		ExecuteSQL(t, server, `
			INSERT INTO listings (
				id, user_id, title, description, price, currency, category_id,
				status, visibility, quantity, view_count, favorites_count
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
			)
		`, 2006, 1006, "Draft Listing", "Description", 50.00, "USD", 15,
			"draft", "private", 1, 0, 0)

		ctx := testutils.TestContext(t)
		req := &pb.UpdateListingRequest{
			Id:     2006,
			UserId: 1006,
			Status: testutils.StringPtr("active"),
		}

		resp, err := server.Client.UpdateListing(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, "active", resp.Listing.Status)
	})
}

// =============================================================================
// 3. GetListing Tests (6 scenarios)
// =============================================================================

func TestGetListing(t *testing.T) {
	testutils.SkipIfShort(t)
	testutils.SkipIfNoDocker(t)

	t.Run("GetByID_Success", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		// Setup
		ExecuteSQL(t, server, `
			INSERT INTO c2c_categories (id, name, slug, parent_id, sort_order, level, is_active, count)
			VALUES ($1, $2, $3, NULL, $4, 0, $5, $6)
		`, 20, "Get Test", "get-test", 1, true, 0)

		ExecuteSQL(t, server, `
			INSERT INTO listings (
				id, user_id, title, description, price, currency, category_id,
				status, visibility, quantity, view_count, favorites_count
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
			)
		`, 3001, 2000, "Test Listing for Get", "Description", 99.99, "USD", 20,
			"active", "public", 1, 0, 0)

		ctx := testutils.TestContext(t)
		req := &pb.GetListingRequest{Id: 3001}

		resp, err := server.Client.GetListing(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, int64(3001), resp.Listing.Id)
		assert.Equal(t, "Test Listing for Get", resp.Listing.Title)
		assert.Equal(t, 99.99, resp.Listing.Price)
		assert.Equal(t, "active", resp.Listing.Status)
	})

	t.Run("GetNonExistent_NotFound", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		ctx := testutils.TestContext(t)
		req := &pb.GetListingRequest{Id: 99999}

		resp, err := server.Client.GetListing(ctx, req)

		require.Error(t, err)
		assert.Nil(t, resp)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
	})

	t.Run("GetDeleted_SoftDelete_NotFound", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		// Setup
		ExecuteSQL(t, server, `
			INSERT INTO c2c_categories (id, name, slug, parent_id, sort_order, level, is_active, count)
			VALUES ($1, $2, $3, NULL, $4, 0, $5, $6)
		`, 21, "Deleted Test", "deleted-test", 1, true, 0)

		ExecuteSQL(t, server, `
			INSERT INTO listings (
				id, user_id, title, description, price, currency, category_id,
				status, visibility, quantity, view_count, favorites_count, deleted_at, is_deleted
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW(), true
			)
		`, 3002, 2001, "Deleted Listing", "Description", 50.00, "USD", 21,
			"archived", "private", 0, 0, 0)

		ctx := testutils.TestContext(t)
		req := &pb.GetListingRequest{Id: 3002}

		resp, err := server.Client.GetListing(ctx, req)

		require.Error(t, err)
		assert.Nil(t, resp)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
	})

	t.Run("GetWithRelatedData_ImagesAndLocation_Success", func(t *testing.T) {
		t.Skip("Images and location loading not yet implemented in GetListingByID")
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		// Setup
		ExecuteSQL(t, server, `
			INSERT INTO c2c_categories (id, name, slug, parent_id, sort_order, level, is_active, count)
			VALUES ($1, $2, $3, NULL, $4, 0, $5, $6)
		`, 22, "Related Data", "related-data", 1, true, 0)

		ExecuteSQL(t, server, `
			INSERT INTO listings (
				id, user_id, title, description, price, currency, category_id,
				status, visibility, quantity, view_count, favorites_count
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
			)
		`, 3003, 2002, "Listing with Relations", "Description", 150.00, "USD", 22,
			"active", "public", 2, 0, 0)

		ExecuteSQL(t, server, `
			INSERT INTO listing_images (
				id, listing_id, url, display_order, is_primary
			) VALUES (
				$1, $2, $3, $4, $5
			), (
				$2, $2, $6, $7, $8
			)
		`, 4001, 3003, "https://example.com/img1.jpg", 1, true,
			"https://example.com/img2.jpg", 2, false)

		ExecuteSQL(t, server, `
			INSERT INTO listing_locations (
				id, listing_id, country, city, postal_code
			) VALUES (
				$1, $2, $3, $4, $5
			)
		`, 5001, 3003, "US", "New York", "10001")

		ctx := testutils.TestContext(t)
		req := &pb.GetListingRequest{Id: 3003}

		resp, err := server.Client.GetListing(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, int64(3003), resp.Listing.Id)

		// Verify images
		assert.Len(t, resp.Listing.Images, 2)
		assert.True(t, resp.Listing.Images[0].IsPrimary)
		assert.Equal(t, "https://example.com/img1.jpg", resp.Listing.Images[0].Url)

		// Verify location
		require.NotNil(t, resp.Listing.Location)
		assert.Equal(t, "US", *resp.Listing.Location.Country)
		assert.Equal(t, "New York", *resp.Listing.Location.City)
	})

	t.Run("GetWithAttributes_Success", func(t *testing.T) {
		t.Skip("Attributes loading not yet implemented in GetListingByID")
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		// Setup
		ExecuteSQL(t, server, `
			INSERT INTO c2c_categories (id, name, slug, parent_id, sort_order, level, is_active, count)
			VALUES ($1, $2, $3, NULL, $4, 0, $5, $6)
		`, 23, "Attributes", "attributes", 1, true, 0)

		ExecuteSQL(t, server, `
			INSERT INTO listings (
				id, user_id, title, description, price, currency, category_id,
				status, visibility, quantity, view_count, favorites_count
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
			)
		`, 3004, 2003, "Listing with Attributes", "Description", 75.00, "USD", 23,
			"active", "public", 1, 0, 0)

		ExecuteSQL(t, server, `
			INSERT INTO listing_attributes (
				id, listing_id, attribute_key, attribute_value
			) VALUES (
				$1, $2, $3, $4
			), (
				$5, $2, $6, $7
			)
		`, 6001, 3004, "brand", "TestBrand",
			6002, "color", "Red")

		ctx := testutils.TestContext(t)
		req := &pb.GetListingRequest{Id: 3004}

		resp, err := server.Client.GetListing(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, int64(3004), resp.Listing.Id)
		assert.Len(t, resp.Listing.Attributes, 2)
	})

	t.Run("MultiLanguageSupport_Cyrillic_Success", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		// Setup
		ExecuteSQL(t, server, `
			INSERT INTO c2c_categories (id, name, slug, parent_id, sort_order, level, is_active, count)
			VALUES ($1, $2, $3, NULL, $4, 0, $5, $6)
		`, 24, "Multilang", "multilang", 1, true, 0)

		ExecuteSQL(t, server, `
			INSERT INTO listings (
				id, user_id, title, description, price, currency, category_id,
				status, visibility, quantity, view_count, favorites_count
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
			)
		`, 3005, 2004, "Товар на русском языке", "Описание на русском", 120.00, "USD", 24,
			"active", "public", 1, 0, 0)

		ctx := testutils.TestContext(t)
		req := &pb.GetListingRequest{Id: 3005}

		resp, err := server.Client.GetListing(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, int64(3005), resp.Listing.Id)
		assert.Equal(t, "Товар на русском языке", resp.Listing.Title)
		assert.Equal(t, "Описание на русском", *resp.Listing.Description)
	})
}

// =============================================================================
// 4. DeleteListing Tests (4 scenarios)
// =============================================================================

func TestDeleteListing(t *testing.T) {
	testutils.SkipIfShort(t)
	testutils.SkipIfNoDocker(t)

	t.Run("SoftDelete_Success", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		// Setup
		ExecuteSQL(t, server, `
			INSERT INTO c2c_categories (id, name, slug, parent_id, sort_order, level, is_active, count)
			VALUES ($1, $2, $3, NULL, $4, 0, $5, $6)
		`, 30, "Delete Test", "delete-test", 1, true, 0)

		ExecuteSQL(t, server, `
			INSERT INTO listings (
				id, user_id, title, description, price, currency, category_id,
				status, visibility, quantity, view_count, favorites_count
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
			)
		`, 4001, 3000, "Listing to Delete", "Description", 45.00, "USD", 30,
			"active", "public", 1, 0, 0)

		ctx := testutils.TestContext(t)
		req := &pb.DeleteListingRequest{
			Id:     4001,
			UserId: 3000,
		}

		resp, err := server.Client.DeleteListing(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.True(t, resp.Success)

		// Verify soft delete: deleted_at should be set
		exists := RowExists(t, server, "listings", "id = $1 AND deleted_at IS NOT NULL AND is_deleted = true", 4001)
		assert.True(t, exists, "Listing should be soft-deleted")
	})

	t.Run("DeleteNonExistent_NotFound", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		ctx := testutils.TestContext(t)
		req := &pb.DeleteListingRequest{
			Id:     99999,
			UserId: 3001,
		}

		resp, err := server.Client.DeleteListing(ctx, req)

		require.Error(t, err)
		assert.Nil(t, resp)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
	})

	t.Run("DeleteByWrongUser_PermissionDenied", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		// Setup
		ExecuteSQL(t, server, `
			INSERT INTO c2c_categories (id, name, slug, parent_id, sort_order, level, is_active, count)
			VALUES ($1, $2, $3, NULL, $4, 0, $5, $6)
		`, 31, "Permission Delete", "permission-delete", 1, true, 0)

		ExecuteSQL(t, server, `
			INSERT INTO listings (
				id, user_id, title, description, price, currency, category_id,
				status, visibility, quantity, view_count, favorites_count
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
			)
		`, 4002, 3002, "Protected Listing", "Description", 30.00, "USD", 31,
			"active", "public", 1, 0, 0)

		ctx := testutils.TestContext(t)
		req := &pb.DeleteListingRequest{
			Id:     4002,
			UserId: 9999, // Different user
		}

		resp, err := server.Client.DeleteListing(ctx, req)

		require.Error(t, err)
		assert.Nil(t, resp)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.PermissionDenied, st.Code())
	})

	t.Run("DeleteAlreadyDeleted_NotFound", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		// Setup
		ExecuteSQL(t, server, `
			INSERT INTO c2c_categories (id, name, slug, parent_id, sort_order, level, is_active, count)
			VALUES ($1, $2, $3, NULL, $4, 0, $5, $6)
		`, 32, "Already Deleted", "already-deleted", 1, true, 0)

		ExecuteSQL(t, server, `
			INSERT INTO listings (
				id, user_id, title, description, price, currency, category_id,
				status, visibility, quantity, view_count, favorites_count, deleted_at, is_deleted
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW(), true
			)
		`, 4003, 3003, "Already Deleted", "Description", 25.00, "USD", 32,
			"archived", "private", 0, 0, 0)

		ctx := testutils.TestContext(t)
		req := &pb.DeleteListingRequest{
			Id:     4003,
			UserId: 3003,
		}

		resp, err := server.Client.DeleteListing(ctx, req)

		require.Error(t, err)
		assert.Nil(t, resp)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.NotFound, st.Code())
	})
}

// =============================================================================
// 5. SearchListings Tests (10 scenarios)
// =============================================================================

func TestSearchListings(t *testing.T) {
	testutils.SkipIfShort(t)
	testutils.SkipIfNoDocker(t)

	t.Run("SearchByCategory_Success", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		// Setup
		ExecuteSQL(t, server, `
			INSERT INTO c2c_categories (id, name, slug, parent_id, sort_order, level, is_active, count)
			VALUES ($1, $2, $3, NULL, $4, 0, $5, $6)
		`, 40, "Search Category", "search-category", 1, true, 0)

		for i := 1; i <= 5; i++ {
			ExecuteSQL(t, server, `
				INSERT INTO listings (
					id, user_id, title, description, price, currency, category_id,
					status, visibility, quantity, view_count, favorites_count
				) VALUES (
					$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
				)
			`, 5000+i, 4000, fmt.Sprintf("Search Listing %d", i),
				"Description", float64(10+i), "USD", 40, "active", "public", 1, 0, 0)
		}

		ctx := testutils.TestContext(t)
		categoryID := int64(40)
		req := &pb.SearchListingsRequest{
			Query:      "",
			CategoryId: &categoryID,
			Limit:      10,
			Offset:     0,
		}

		resp, err := server.Client.SearchListings(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.GreaterOrEqual(t, len(resp.Listings), 5)
		assert.GreaterOrEqual(t, resp.Total, int32(5))
	})

	t.Run("SearchByPriceRange_Success", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		// Setup
		ExecuteSQL(t, server, `
			INSERT INTO c2c_categories (id, name, slug, parent_id, sort_order, level, is_active, count)
			VALUES ($1, $2, $3, NULL, $4, 0, $5, $6)
		`, 41, "Price Range", "price-range", 1, true, 0)

		prices := []float64{10.00, 25.00, 50.00, 75.00, 100.00}
		for i, price := range prices {
			ExecuteSQL(t, server, `
				INSERT INTO listings (
					id, user_id, title, description, price, currency, category_id,
					status, visibility, quantity, view_count, favorites_count
				) VALUES (
					$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
				)
			`, 5010+i, 4001, fmt.Sprintf("Price Listing %d", i),
				"Description", price, "USD", 41, "active", "public", 1, 0, 0)
		}

		ctx := testutils.TestContext(t)
		minPrice := 20.0
		maxPrice := 80.0
		req := &pb.SearchListingsRequest{
			Query:    "",
			MinPrice: &minPrice,
			MaxPrice: &maxPrice,
			Limit:    10,
			Offset:   0,
		}

		resp, err := server.Client.SearchListings(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		// Should find listings with prices 25.00, 50.00, 75.00
		assert.GreaterOrEqual(t, len(resp.Listings), 3)

		for _, listing := range resp.Listings {
			assert.GreaterOrEqual(t, listing.Price, 20.0)
			assert.LessOrEqual(t, listing.Price, 80.0)
		}
	})

	t.Run("SearchByTitle_TextSearch_Success", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		// Setup
		ExecuteSQL(t, server, `
			INSERT INTO c2c_categories (id, name, slug, parent_id, sort_order, level, is_active, count)
			VALUES ($1, $2, $3, NULL, $4, 0, $5, $6)
		`, 42, "Text Search", "text-search", 1, true, 0)

		ExecuteSQL(t, server, `
			INSERT INTO listings (
				id, user_id, title, description, price, currency, category_id,
				status, visibility, quantity, view_count, favorites_count
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
			)
		`, 5020, 4002, "Unique Laptop Device", "Electronic device", 500.00, "USD", 42,
			"active", "public", 1, 0, 0)

		ctx := testutils.TestContext(t)
		req := &pb.SearchListingsRequest{
			Query:  "Laptop",
			Limit:  10,
			Offset: 0,
		}

		resp, err := server.Client.SearchListings(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		// Should find the "Unique Laptop Device"
		assert.GreaterOrEqual(t, len(resp.Listings), 1)
	})

	t.Run("Pagination_OffsetAndLimit_Success", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		// Setup
		ExecuteSQL(t, server, `
			INSERT INTO c2c_categories (id, name, slug, parent_id, sort_order, level, is_active, count)
			VALUES ($1, $2, $3, NULL, $4, 0, $5, $6)
		`, 43, "Pagination", "pagination", 1, true, 0)

		for i := 1; i <= 15; i++ {
			ExecuteSQL(t, server, `
				INSERT INTO listings (
					id, user_id, title, description, price, currency, category_id,
					status, visibility, quantity, view_count, favorites_count
				) VALUES (
					$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
				)
			`, 5030+i, 4003, fmt.Sprintf("Pagination Listing %d", i),
				"Description", float64(20+i), "USD", 43, "active", "public", 1, 0, 0)
		}

		ctx := testutils.TestContext(t)

		// Page 1
		req1 := &pb.SearchListingsRequest{
			Query:  "",
			Limit:  5,
			Offset: 0,
		}
		resp1, err := server.Client.SearchListings(ctx, req1)
		require.NoError(t, err)
		assert.Len(t, resp1.Listings, 5)

		// Page 2
		req2 := &pb.SearchListingsRequest{
			Query:  "",
			Limit:  5,
			Offset: 5,
		}
		resp2, err := server.Client.SearchListings(ctx, req2)
		require.NoError(t, err)
		assert.Len(t, resp2.Listings, 5)

		// Ensure different results
		assert.NotEqual(t, resp1.Listings[0].Id, resp2.Listings[0].Id)
	})

	t.Run("CombinedFilters_CategoryAndPrice_Success", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		// Setup
		ExecuteSQL(t, server, `
			INSERT INTO c2c_categories (id, name, slug, parent_id, sort_order, level, is_active, count)
			VALUES ($1, $2, $3, NULL, $4, 0, $5, $6)
		`, 44, "Combined", "combined", 1, true, 0)

		for i := 1; i <= 10; i++ {
			ExecuteSQL(t, server, `
				INSERT INTO listings (
					id, user_id, title, description, price, currency, category_id,
					status, visibility, quantity, view_count, favorites_count
				) VALUES (
					$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
				)
			`, 5040+i, 4004, fmt.Sprintf("Combined Listing %d", i),
				"Description", float64(30+i*5), "USD", 44, "active", "public", 1, 0, 0)
		}

		ctx := testutils.TestContext(t)
		categoryID := int64(44)
		minPrice := 40.0
		maxPrice := 70.0
		req := &pb.SearchListingsRequest{
			Query:      "",
			CategoryId: &categoryID,
			MinPrice:   &minPrice,
			MaxPrice:   &maxPrice,
			Limit:      20,
			Offset:     0,
		}

		resp, err := server.Client.SearchListings(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Greater(t, len(resp.Listings), 0)

		for _, listing := range resp.Listings {
			assert.Equal(t, int64(44), listing.CategoryId)
			assert.GreaterOrEqual(t, listing.Price, 40.0)
			assert.LessOrEqual(t, listing.Price, 70.0)
		}
	})

	t.Run("EmptyResults_NoMatch_Success", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		ctx := testutils.TestContext(t)
		req := &pb.SearchListingsRequest{
			Query:  "NonExistentSearchTerm12345",
			Limit:  10,
			Offset: 0,
		}

		resp, err := server.Client.SearchListings(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Empty(t, resp.Listings)
		assert.Equal(t, int32(0), resp.Total)
	})

	t.Run("SortByPrice_Ascending_Success", func(t *testing.T) {
		t.Skip("Sorting implementation pending in SearchListings")
		// TODO: Implement when sorting is added to SearchListings
	})

	t.Run("SortByDate_Descending_Success", func(t *testing.T) {
		t.Skip("Sorting implementation pending in SearchListings")
		// TODO: Implement when sorting is added to SearchListings
	})

	t.Run("FilterBySourceType_C2CVsB2C_Success", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		// Setup
		ExecuteSQL(t, server, `
			INSERT INTO c2c_categories (id, name, slug, parent_id, sort_order, level, is_active, count)
			VALUES ($1, $2, $3, NULL, $4, 0, $5, $6)
		`, 45, "Source Type", "source-type", 1, true, 0)

		ExecuteSQL(t, server, `
			INSERT INTO storefronts (id, user_id, slug, name, country, is_active, is_verified)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, 5002, 4005, "source-store", "Source Store", "US", true, true)

		// C2C listing (no storefront_id)
		ExecuteSQL(t, server, `
			INSERT INTO listings (
				id, user_id, storefront_id, title, description, price, currency, category_id,
				status, visibility, quantity, view_count, favorites_count
			) VALUES (
				$1, $2, $3, NULL, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
			)
		`, 5050, 4005, "C2C Listing", "Description", 30.00, "USD", 45,
			"active", "public", 1, 0, 0)

		// B2C listing (with storefront_id)
		ExecuteSQL(t, server, `
			INSERT INTO listings (
				id, user_id, storefront_id, title, description, price, currency, category_id,
				status, visibility, quantity, view_count, favorites_count
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
			)
		`, 5051, 4005, 5002, "B2C Listing", "Description", 40.00, "USD", 45,
			"active", "public", 2, 0, 0)

		ctx := testutils.TestContext(t)
		categoryID := int64(45)
		req := &pb.SearchListingsRequest{
			Query:      "",
			CategoryId: &categoryID,
			Limit:      10,
			Offset:     0,
		}

		resp, err := server.Client.SearchListings(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.GreaterOrEqual(t, len(resp.Listings), 2)

		// Verify both C2C and B2C are returned
		hasC2C := false
		hasB2C := false
		for _, listing := range resp.Listings {
			if listing.StorefrontId == nil {
				hasC2C = true
			} else {
				hasB2C = true
			}
		}
		assert.True(t, hasC2C, "Should have C2C listings")
		assert.True(t, hasB2C, "Should have B2C listings")
	})

	t.Run("Performance_LargeDataset_Success", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		// Setup
		ExecuteSQL(t, server, `
			INSERT INTO c2c_categories (id, name, slug, parent_id, sort_order, level, is_active, count)
			VALUES ($1, $2, $3, NULL, $4, 0, $5, $6)
		`, 46, "Performance", "performance", 1, true, 0)

		// Insert 100 listings for performance test
		for i := 1; i <= 100; i++ {
			ExecuteSQL(t, server, `
				INSERT INTO listings (
					id, user_id, title, description, price, currency, category_id,
					status, visibility, quantity, view_count, favorites_count
				) VALUES (
					$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
				)
			`, 6000+i, 5000, fmt.Sprintf("Performance Listing %d", i),
				"Performance test description", float64(50+i), "USD", 46, "active", "public", 1, 0, 0)
		}

		ctx := testutils.TestContext(t)
		categoryID := int64(46)
		req := &pb.SearchListingsRequest{
			Query:      "",
			CategoryId: &categoryID,
			Limit:      50,
			Offset:     0,
		}

		start := time.Now()
		resp, err := server.Client.SearchListings(ctx, req)
		duration := time.Since(start)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Len(t, resp.Listings, 50)
		assert.Equal(t, int32(100), resp.Total)

		// Performance check: should complete in reasonable time
		assert.Less(t, duration, 1*time.Second, "Search should complete within 1 second")
	})
}

// =============================================================================
// 6. Error Cases Tests (4 scenarios)
// =============================================================================

func TestListingErrorCases(t *testing.T) {
	testutils.SkipIfShort(t)
	testutils.SkipIfNoDocker(t)

	t.Run("Timeout_LongRunningOperation_Error", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		// Create context with very short timeout
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()

		req := &pb.GetListingRequest{Id: 1}
		_, err := server.Client.GetListing(ctx, req)

		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.DeadlineExceeded, st.Code())
	})

	t.Run("InvalidProtoMessage_MalformedRequest_Error", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		ctx := testutils.TestContext(t)

		// Invalid request: negative ID
		req := &pb.GetListingRequest{Id: -1}
		resp, err := server.Client.GetListing(ctx, req)

		require.Error(t, err)
		assert.Nil(t, resp)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Contains(t, []codes.Code{codes.InvalidArgument, codes.NotFound}, st.Code())
	})

	t.Run("DatabaseError_SimulatedFailure", func(t *testing.T) {
		t.Skip("Database failure simulation requires special setup")
		// TODO: Implement with test database that can simulate failures
	})

	t.Run("RateLimiting_TooManyRequests", func(t *testing.T) {
		t.Skip("Rate limiting not implemented in test infrastructure")
		// TODO: Implement when rate limiting is added to gRPC service
	})
}

// =============================================================================
// 7. Edge Cases Tests (2 scenarios)
// =============================================================================

func TestListingEdgeCases(t *testing.T) {
	testutils.SkipIfShort(t)
	testutils.SkipIfNoDocker(t)

	t.Run("Unicode_CyrillicAndEmoji_Success", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		// Setup
		ExecuteSQL(t, server, `
			INSERT INTO c2c_categories (id, name, slug, parent_id, sort_order, level, is_active, count)
			VALUES ($1, $2, $3, NULL, $4, 0, $5, $6)
		`, 50, "Unicode", "unicode", 1, true, 0)

		ctx := testutils.TestContext(t)
		req := &pb.CreateListingRequest{
			UserId:      6000,
			Title:       "Товар 🎉 Product with Emoji 😊",
			Description: testutils.StringPtr("Описание с кириллицей и эмодзи 🚀"),
			Price:       99.99,
			Currency:    "USD",
			CategoryId:  50,
			Quantity:    1,
		}

		resp, err := server.Client.CreateListing(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, "Товар 🎉 Product with Emoji 😊", resp.Listing.Title)
		assert.Contains(t, *resp.Listing.Description, "кириллицей")
		assert.Contains(t, *resp.Listing.Description, "🚀")

		// Verify retrieval
		getResp, err := server.Client.GetListing(ctx, &pb.GetListingRequest{Id: resp.Listing.Id})
		require.NoError(t, err)
		assert.Equal(t, "Товар 🎉 Product with Emoji 😊", getResp.Listing.Title)
	})

	t.Run("BoundaryValues_MaxPriceAndTitle_Success", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		// Setup
		ExecuteSQL(t, server, `
			INSERT INTO c2c_categories (id, name, slug, parent_id, sort_order, level, is_active, count)
			VALUES ($1, $2, $3, NULL, $4, 0, $5, $6)
		`, 51, "Boundary", "boundary", 1, true, 0)

		ctx := testutils.TestContext(t)

		// Very long title (200 characters - maximum allowed)
		longTitle := strings.Repeat("A", 200)

		req := &pb.CreateListingRequest{
			UserId:     6001,
			Title:      longTitle,
			Price:      999999.99, // Very high price
			Currency:   "USD",
			CategoryId: 51,
			Quantity:   999999, // Maximum quantity
		}

		resp, err := server.Client.CreateListing(ctx, req)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.NotZero(t, resp.Listing.Id)
		assert.Equal(t, 999999.99, resp.Listing.Price)
		assert.Equal(t, int32(999999), resp.Listing.Quantity)

		// Verify minimum values
		minReq := &pb.CreateListingRequest{
			UserId:     6002,
			Title:      "Min",  // Minimum 3 characters (per validation rules)
			Price:      0.01, // Minimum price
			Currency:   "USD",
			CategoryId: 51,
			Quantity:   1, // Minimum quantity
		}

		minResp, err := server.Client.CreateListing(ctx, minReq)
		require.NoError(t, err)
		require.NotNil(t, minResp)
		assert.Equal(t, "Min", minResp.Listing.Title)
		assert.Equal(t, 0.01, minResp.Listing.Price)
		assert.Equal(t, int32(1), minResp.Listing.Quantity)
	})
}
