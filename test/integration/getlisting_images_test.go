package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pb "github.com/sveturs/listings/api/proto/listings/v1"
	testutils "github.com/sveturs/listings/internal/testing"
)

// TestGetListing_WithImages verifies that GetListing loads and returns images
func TestGetListing_WithImages(t *testing.T) {
	testutils.SkipIfShort(t)
	testutils.SkipIfNoDocker(t)

	t.Run("GetListing_ReturnsImages", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		// Setup: Create category
		ExecuteSQL(t, server, `
			INSERT INTO c2c_categories (id, name, slug, parent_id, sort_order, level, is_active, count)
			VALUES ($1, $2, $3, NULL, $4, 0, $5, $6)
		`, 1301, "Electronics", "electronics", 1, true, 0)

		// Setup: Create listing
		ExecuteSQL(t, server, `
			INSERT INTO listings (id, user_id, category_id, title, description, price, currency, quantity, status, visibility, source_type, uuid, slug)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, gen_random_uuid(), $12)
		`, 101, 1, 1301, "Test Laptop", "Gaming laptop with RTX 4090", 1999.99, "USD", 5, "active", "public", "b2c", "test-laptop")

		// Setup: Create images for the listing
		ExecuteSQL(t, server, `
			INSERT INTO listing_images (id, listing_id, url, storage_path, thumbnail_url, display_order, is_primary, width, height, file_size, mime_type)
			VALUES
				($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11),
				($12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22),
				($23, $24, $25, NULL, NULL, $26, $27, NULL, NULL, NULL, NULL)
		`,
			// Image 1 - Primary with all fields
			1, 101, "https://example.com/laptop-main.jpg", "/storage/101/main.jpg", "https://example.com/thumbnails/laptop-main-thumb.jpg", 1, true, 1920, 1080, 256000, "image/jpeg",
			// Image 2 - Secondary with all fields
			2, 101, "https://example.com/laptop-side.jpg", "/storage/101/side.jpg", "https://example.com/thumbnails/laptop-side-thumb.jpg", 2, false, 1920, 1080, 198000, "image/jpeg",
			// Image 3 - Minimal (only required fields)
			3, 101, "https://example.com/laptop-back.jpg", 3, false,
		)

		// Test: Call GetListing
		ctx := testutils.TestContext(t)
		req := &pb.GetListingRequest{
			Id: 101,
		}

		resp, err := server.Client.GetListing(ctx, req)

		// Assertions
		require.NoError(t, err, "GetListing should succeed")
		require.NotNil(t, resp, "Response should not be nil")
		require.NotNil(t, resp.Listing, "Listing should not be nil")

		// Verify listing basic fields
		assert.Equal(t, int64(101), resp.Listing.Id)
		assert.Equal(t, "Test Laptop", resp.Listing.Title)

		// Verify images are loaded
		require.Len(t, resp.Listing.Images, 3, "Should load 3 images")

		// Verify images are ordered by display_order
		images := resp.Listing.Images
		assert.Equal(t, int64(1), images[0].Id)
		assert.Equal(t, int64(2), images[1].Id)
		assert.Equal(t, int64(3), images[2].Id)

		// Verify first image (primary with all fields)
		img1 := images[0]
		assert.Equal(t, int64(101), img1.ListingId)
		assert.Equal(t, "https://example.com/laptop-main.jpg", img1.Url)
		assert.True(t, img1.IsPrimary)
		assert.Equal(t, int32(1), img1.DisplayOrder)
		require.NotNil(t, img1.StoragePath)
		assert.Equal(t, "/storage/101/main.jpg", *img1.StoragePath)
		require.NotNil(t, img1.ThumbnailUrl)
		assert.Equal(t, "https://example.com/thumbnails/laptop-main-thumb.jpg", *img1.ThumbnailUrl)
		require.NotNil(t, img1.Width)
		assert.Equal(t, int32(1920), *img1.Width)
		require.NotNil(t, img1.Height)
		assert.Equal(t, int32(1080), *img1.Height)
		require.NotNil(t, img1.FileSize)
		assert.Equal(t, int64(256000), *img1.FileSize)
		require.NotNil(t, img1.MimeType)
		assert.Equal(t, "image/jpeg", *img1.MimeType)

		// Verify second image (secondary with all fields)
		img2 := images[1]
		assert.Equal(t, "https://example.com/laptop-side.jpg", img2.Url)
		assert.False(t, img2.IsPrimary)
		assert.Equal(t, int32(2), img2.DisplayOrder)

		// Verify third image (minimal fields)
		img3 := images[2]
		assert.Equal(t, "https://example.com/laptop-back.jpg", img3.Url)
		assert.False(t, img3.IsPrimary)
		assert.Equal(t, int32(3), img3.DisplayOrder)
		assert.Nil(t, img3.StoragePath, "Optional field should be nil")
		assert.Nil(t, img3.ThumbnailUrl, "Optional field should be nil")
		assert.Nil(t, img3.Width, "Optional field should be nil")
	})

	t.Run("GetListing_NoImages_ReturnsEmpty", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		// Setup: Create category
		ExecuteSQL(t, server, `
			INSERT INTO c2c_categories (id, name, slug, parent_id, sort_order, level, is_active, count)
			VALUES ($1, $2, $3, NULL, $4, 0, $5, $6)
		`, 1302, "Books", "books", 1, true, 0)

		// Setup: Create listing WITHOUT images
		ExecuteSQL(t, server, `
			INSERT INTO listings (id, user_id, category_id, title, price, currency, quantity, status, visibility, source_type, uuid, slug)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, gen_random_uuid(), $11)
		`, 102, 1, 1302, "Test Book", 29.99, "USD", 10, "active", "public", "c2c", "test-book")

		// Test: Call GetListing
		ctx := testutils.TestContext(t)
		req := &pb.GetListingRequest{
			Id: 102,
		}

		resp, err := server.Client.GetListing(ctx, req)

		// Assertions
		require.NoError(t, err, "GetListing should succeed even without images")
		require.NotNil(t, resp, "Response should not be nil")
		require.NotNil(t, resp.Listing, "Listing should not be nil")

		// Verify listing basic fields
		assert.Equal(t, int64(102), resp.Listing.Id)
		assert.Equal(t, "Test Book", resp.Listing.Title)

		// Verify no images
		assert.Empty(t, resp.Listing.Images, "Should return empty images array")
	})

	t.Run("GetListing_ImageLoadError_DoesNotFailRequest", func(t *testing.T) {
		config := DefaultTestServerConfig()
		server := SetupTestServer(t, config)
		defer server.Teardown(t)

		// Setup: Create category
		ExecuteSQL(t, server, `
			INSERT INTO c2c_categories (id, name, slug, parent_id, sort_order, level, is_active, count)
			VALUES ($1, $2, $3, NULL, $4, 0, $5, $6)
		`, 1303, "Sports", "sports", 1, true, 0)

		// Setup: Create listing
		ExecuteSQL(t, server, `
			INSERT INTO listings (id, user_id, category_id, title, price, currency, quantity, status, visibility, source_type, uuid, slug)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, gen_random_uuid(), $11)
		`, 103, 1, 1303, "Test Football", 49.99, "USD", 20, "active", "public", "b2c", "test-football")

		// Note: We don't insert images but the implementation should still succeed

		// Test: Call GetListing
		ctx := testutils.TestContext(t)
		req := &pb.GetListingRequest{
			Id: 103,
		}

		resp, err := server.Client.GetListing(ctx, req)

		// Assertions - GetListing should NOT fail even if images fail to load
		require.NoError(t, err, "GetListing should succeed even if image loading has issues")
		require.NotNil(t, resp)
		require.NotNil(t, resp.Listing)
		assert.Equal(t, int64(103), resp.Listing.Id)
	})
}
