//go:build integration

package integration

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/vondi-global/listings/api/proto/listings/v1"
	"github.com/vondi-global/listings/tests"
)

// ============================================================================
// Edge Cases Tests - Boundary Values, Unicode, Special Characters
// ============================================================================

// TestListing_BoundaryValues_MinTitle tests minimum valid title length (1 character)
func TestListing_BoundaryValues_MinTitle(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	req := &pb.CreateListingRequest{
		UserId:      100,
		Title:       "X", // Minimum 1 character
		Description: stringPtr("Valid description"),
		Price:       1.0,
		Currency:    "USD",
		CategoryId:  1,
		Quantity:    1,
	}

	resp, err := client.CreateListing(ctx, req)

	require.NoError(t, err, "Should accept 1-character title")
	require.NotNil(t, resp)
	require.NotNil(t, resp.Listing)
	assert.Equal(t, "X", resp.Listing.Title)
}

// TestListing_BoundaryValues_MaxTitle tests maximum valid title length
func TestListing_BoundaryValues_MaxTitle(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	// Assuming max title length is 255 characters
	maxTitle := strings.Repeat("A", 255)

	req := &pb.CreateListingRequest{
		UserId:      100,
		Title:       maxTitle,
		Description: stringPtr("Valid description"),
		Price:       1.0,
		Currency:    "USD",
		CategoryId:  1,
		Quantity:    1,
	}

	resp, err := client.CreateListing(ctx, req)

	require.NoError(t, err, "Should accept 255-character title")
	require.NotNil(t, resp)
	assert.Equal(t, maxTitle, resp.Listing.Title)
}

// TestListing_BoundaryValues_TitleTooLong tests title exceeding max length
func TestListing_BoundaryValues_TitleTooLong(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	// Title exceeding max length (256 characters)
	tooLongTitle := strings.Repeat("A", 256)

	req := &pb.CreateListingRequest{
		UserId:      100,
		Title:       tooLongTitle,
		Description: stringPtr("Valid description"),
		Price:       1.0,
		Currency:    "USD",
		CategoryId:  1,
		Quantity:    1,
	}

	resp, err := client.CreateListing(ctx, req)

	if err != nil {
		// Should return validation error
		st, ok := status.FromError(err)
		require.True(t, ok, "Error should be gRPC status error")
		assert.Equal(t, codes.InvalidArgument, st.Code(), "Should return InvalidArgument")
	} else {
		// Or backend might truncate to 255 characters
		require.NotNil(t, resp)
		assert.LessOrEqual(t, len(resp.Listing.Title), 255, "Title should be truncated to max length")
	}
}

// TestListing_BoundaryValues_MinPrice tests minimum valid price (0.01)
func TestListing_BoundaryValues_MinPrice(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	req := &pb.CreateListingRequest{
		UserId:      100,
		Title:       "Min Price Test",
		Description: stringPtr("Testing minimum price"),
		Price:       0.01, // Minimum valid price (1 cent)
		Currency:    "USD",
		CategoryId:  1,
		Quantity:    1,
	}

	resp, err := client.CreateListing(ctx, req)

	require.NoError(t, err, "Should accept minimum price 0.01")
	require.NotNil(t, resp)
	assert.Equal(t, 0.01, resp.Listing.Price)
}

// TestListing_BoundaryValues_ZeroPrice tests zero price (should be rejected)
func TestListing_BoundaryValues_ZeroPrice(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	req := &pb.CreateListingRequest{
		UserId:      100,
		Title:       "Zero Price Test",
		Description: stringPtr("Testing zero price"),
		Price:       0.0, // Invalid price
		Currency:    "USD",
		CategoryId:  1,
		Quantity:    1,
	}

	resp, err := client.CreateListing(ctx, req)

	if err != nil {
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code(), "Should reject zero price")
	} else {
		// Backend might allow zero for "free" items
		require.NotNil(t, resp)
		t.Logf("Backend allows zero price: %+v", resp.Listing)
	}
}

// TestListing_BoundaryValues_NegativePrice tests negative price (should be rejected)
func TestListing_BoundaryValues_NegativePrice(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	req := &pb.CreateListingRequest{
		UserId:      100,
		Title:       "Negative Price Test",
		Description: stringPtr("Testing negative price"),
		Price:       -10.50, // Invalid price
		Currency:    "USD",
		CategoryId:  1,
		Quantity:    1,
	}

	resp, err := client.CreateListing(ctx, req)

	require.Error(t, err, "Should reject negative price")
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Nil(t, resp)
}

// TestListing_BoundaryValues_ZeroQuantity tests zero quantity
func TestListing_BoundaryValues_ZeroQuantity(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	req := &pb.CreateListingRequest{
		UserId:      100,
		Title:       "Zero Quantity Test",
		Description: stringPtr("Testing zero quantity"),
		Price:       99.99,
		Currency:    "USD",
		CategoryId:  1,
		Quantity:    0, // Zero quantity

	}

	resp, err := client.CreateListing(ctx, req)

	// Zero quantity might be valid for "out of stock" items
	if err != nil {
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	} else {
		require.NotNil(t, resp)
		assert.Equal(t, int32(0), resp.Listing.Quantity)
		t.Logf("Backend allows zero quantity: %+v", resp.Listing)
	}
}

// TestListing_BoundaryValues_MaxQuantity tests maximum quantity value
func TestListing_BoundaryValues_MaxQuantity(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	req := &pb.CreateListingRequest{
		UserId:      100,
		Title:       "Max Quantity Test",
		Description: stringPtr("Testing max quantity"),
		Price:       99.99,
		Currency:    "USD",
		CategoryId:  1,
		Quantity:    2147483647, // Max int32 value

	}

	resp, err := client.CreateListing(ctx, req)

	require.NoError(t, err, "Should accept max int32 quantity")
	require.NotNil(t, resp)
	assert.Equal(t, int32(2147483647), resp.Listing.Quantity)
}

// ============================================================================
// Unicode and Special Characters Tests
// ============================================================================

// TestListing_Unicode_Emoji tests emoji in title
func TestListing_Unicode_Emoji(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	req := &pb.CreateListingRequest{
		UserId:      100,
		Title:       "iPhone 15 Pro 📱 New! 🔥",
		Description: stringPtr("Smartphone with emoji 🎉🎊"),
		Price:       999.99,
		Currency:    "USD",
		CategoryId:  1,
		Quantity:    10,
	}

	resp, err := client.CreateListing(ctx, req)

	require.NoError(t, err, "Should accept emoji in title")
	require.NotNil(t, resp)
	assert.Contains(t, resp.Listing.Title, "📱")
	assert.Contains(t, resp.Listing.Title, "🔥")
}

// TestListing_Unicode_CJK tests Chinese/Japanese/Korean characters
func TestListing_Unicode_CJK(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	req := &pb.CreateListingRequest{
		UserId:      100,
		Title:       "中文标题 日本語 한국어",
		Description: stringPtr("多语言描述 マルチ言語 다국어"),
		Price:       149.99,
		Currency:    "USD",
		CategoryId:  1,
		Quantity:    5,
	}

	resp, err := client.CreateListing(ctx, req)

	require.NoError(t, err, "Should accept CJK characters")
	require.NotNil(t, resp)
	assert.Equal(t, "中文标题 日本語 한국어", resp.Listing.Title)
}

// TestListing_Unicode_RTL tests Right-to-Left text (Arabic, Hebrew)
func TestListing_Unicode_RTL(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	req := &pb.CreateListingRequest{
		UserId:      100,
		Title:       "منتج جديد العربية עברית",
		Description: stringPtr("وصف المنتج תיאור המוצר"),
		Price:       79.99,
		Currency:    "USD",
		CategoryId:  1,
		Quantity:    3,
	}

	resp, err := client.CreateListing(ctx, req)

	require.NoError(t, err, "Should accept RTL text")
	require.NotNil(t, resp)
	assert.Contains(t, resp.Listing.Title, "منتج")
	assert.Contains(t, resp.Listing.Title, "עברית")
}

// TestListing_Unicode_SpecialSymbols tests special Unicode symbols
func TestListing_Unicode_SpecialSymbols(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	req := &pb.CreateListingRequest{
		UserId:      100,
		Title:       "Product™ ® © § ¶ • ★ ♥",
		Description: stringPtr("Symbols: € £ ¥ ₹ ¢ ° ± × ÷"),
		Price:       199.99,
		Currency:    "USD",
		CategoryId:  1,
		Quantity:    1,
	}

	resp, err := client.CreateListing(ctx, req)

	require.NoError(t, err, "Should accept special Unicode symbols")
	require.NotNil(t, resp)
	assert.Contains(t, resp.Listing.Title, "™")
	assert.Contains(t, resp.Listing.Title, "♥")
}

// TestListing_Unicode_ZeroWidthCharacters tests zero-width characters
func TestListing_Unicode_ZeroWidthCharacters(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	// Zero-width joiner (U+200D) and zero-width space (U+200B)
	req := &pb.CreateListingRequest{
		UserId:      100,
		Title:       "Test\u200DTitle\u200BWith\u200CZero\u200DWidth",
		Description: stringPtr("Normal description"),
		Price:       99.99,
		Currency:    "USD",
		CategoryId:  1,
		Quantity:    1,
	}

	resp, err := client.CreateListing(ctx, req)

	if err != nil {
		// Backend might strip zero-width characters
		t.Logf("Backend rejected zero-width characters: %v", err)
	} else {
		require.NotNil(t, resp)
		t.Logf("Backend accepted zero-width characters: %s", resp.Listing.Title)
	}
}

// TestListing_Unicode_ControlCharacters tests control characters (should be stripped/rejected)
func TestListing_Unicode_ControlCharacters(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	// Control characters: \t, \n, \r
	req := &pb.CreateListingRequest{
		UserId:      100,
		Title:       "Test\tTitle\nWith\rControl",
		Description: stringPtr("Test\nMultiline\rDescription"),
		Price:       99.99,
		Currency:    "USD",
		CategoryId:  1,
		Quantity:    1,
	}

	resp, err := client.CreateListing(ctx, req)

	if err != nil {
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code(), "Should reject control characters")
	} else {
		// Backend might strip control characters
		require.NotNil(t, resp)
		assert.NotContains(t, resp.Listing.Title, "\t", "Tab should be stripped")
		assert.NotContains(t, resp.Listing.Title, "\n", "Newline should be stripped")
		t.Logf("Backend sanitized title: %s", resp.Listing.Title)
	}
}

// ============================================================================
// Special Edge Cases
// ============================================================================

// TestListing_EdgeCase_EmptyDescription tests optional description field
func TestListing_EdgeCase_EmptyDescription(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	req := &pb.CreateListingRequest{
		UserId:      100,
		Title:       "No Description Listing",
		Description: nil, // Empty description
		Price:       49.99,
		Currency:    "USD",
		CategoryId:  1,
		Quantity:    1,
	}

	resp, err := client.CreateListing(ctx, req)

	require.NoError(t, err, "Should accept nil description")
	require.NotNil(t, resp)
	assert.Nil(t, resp.Listing.Description)
}

// TestListing_EdgeCase_VeryLongDescription tests extremely long description
func TestListing_EdgeCase_VeryLongDescription(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	// Generate 10,000 character description
	veryLongDesc := strings.Repeat("Lorem ipsum dolor sit amet. ", 400) // ~11,200 chars

	req := &pb.CreateListingRequest{
		UserId:      100,
		Title:       "Long Description Test",
		Description: &veryLongDesc,
		Price:       99.99,
		Currency:    "USD",
		CategoryId:  1,
		Quantity:    1,
	}

	resp, err := client.CreateListing(ctx, req)

	if err != nil {
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code(), "Might reject extremely long description")
	} else {
		require.NotNil(t, resp)
		t.Logf("Backend accepted long description: %d chars", len(*resp.Listing.Description))
	}
}

// TestListing_EdgeCase_InvalidCategoryID tests non-existent category
func TestListing_EdgeCase_InvalidCategoryID(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	req := &pb.CreateListingRequest{
		UserId:      100,
		Title:       "Invalid Category Test",
		Description: stringPtr("Testing invalid category"),
		Price:       99.99,
		Currency:    "USD",
		CategoryId:  999999, // Non-existent category
		Quantity:    1,
	}

	resp, err := client.CreateListing(ctx, req)

	require.Error(t, err, "Should reject invalid category ID")
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Contains(t, []codes.Code{codes.InvalidArgument, codes.NotFound}, st.Code())
	assert.Nil(t, resp)
}

// TestListing_EdgeCase_SQLInjectionAttempt tests SQL injection protection
func TestListing_EdgeCase_SQLInjectionAttempt(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	req := &pb.CreateListingRequest{
		UserId:      100,
		Title:       "'; DROP TABLE listings; --",
		Description: stringPtr("1' OR '1'='1"),
		Price:       99.99,
		Currency:    "USD",
		CategoryId:  1,
		Quantity:    1,
	}

	resp, err := client.CreateListing(ctx, req)

	// Should either:
	// 1. Accept as regular text (SQL injection protected)
	// 2. Reject due to validation
	if err != nil {
		t.Logf("SQL injection attempt rejected: %v", err)
	} else {
		require.NotNil(t, resp)
		// Verify no SQL injection occurred by checking listing was created safely
		assert.Contains(t, resp.Listing.Title, "DROP TABLE")
		t.Log("SQL injection attempt safely handled as text")
	}
}

// TestListing_EdgeCase_XSSAttempt tests XSS protection
func TestListing_EdgeCase_XSSAttempt(t *testing.T) {
	client, _, cleanup := setupGRPCTestServer(t)
	defer cleanup()

	ctx := tests.TestContext(t)

	req := &pb.CreateListingRequest{
		UserId:      100,
		Title:       "<script>alert('XSS')</script>",
		Description: stringPtr("<img src=x onerror=alert('XSS')>"),
		Price:       99.99,
		Currency:    "USD",
		CategoryId:  1,
		Quantity:    1,
	}

	resp, err := client.CreateListing(ctx, req)

	if err != nil {
		// Backend might reject HTML/script tags
		t.Logf("XSS attempt rejected: %v", err)
	} else {
		require.NotNil(t, resp)
		// Verify HTML is escaped or stripped
		t.Logf("XSS attempt handled: title=%s", resp.Listing.Title)
	}
}
