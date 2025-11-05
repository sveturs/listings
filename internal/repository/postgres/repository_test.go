package postgres

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sveturs/listings/internal/domain"
	"github.com/sveturs/listings/tests"
)

func setupTestRepo(t *testing.T) (*Repository, *tests.TestDB) {
	t.Helper()

	// Skip if running in short mode
	tests.SkipIfShort(t)
	tests.SkipIfNoDocker(t)

	// Setup test database
	testDB := tests.SetupTestPostgres(t)

	// Run migrations
	tests.RunMigrations(t, testDB.DB, "../../../migrations")

	// Create sqlx.DB wrapper
	db := sqlx.NewDb(testDB.DB, "postgres")

	// Create logger
	logger := zerolog.New(zerolog.NewTestWriter(t)).With().Timestamp().Logger()

	// Create repository
	repo := NewRepository(db, logger)

	// Create test categories (fixtures)
	setupTestCategories(t, db)

	return repo, testDB
}

// setupTestCategories creates test category fixtures
func setupTestCategories(t *testing.T, db *sqlx.DB) {
	t.Helper()

	categories := []struct {
		id          int
		name        string
		slug        string
		description string
	}{
		{100, "Test Electronics", "test-electronics", "Test category for electronics"},
		{200, "Test Fashion", "test-fashion", "Test category for fashion items"},
		{300, "Test Home & Garden", "test-home-garden", "Test category for home and garden"},
	}

	for _, cat := range categories {
		_, err := db.Exec(`
			INSERT INTO c2c_categories (id, name, slug, description, is_active, level, sort_order)
			VALUES ($1, $2, $3, $4, true, 0, 0)
			ON CONFLICT (id) DO NOTHING
		`, cat.id, cat.name, cat.slug, cat.description)
		if err != nil {
			t.Fatalf("failed to create test category: %v", err)
		}
	}
}

func TestNewRepository(t *testing.T) {
	db := &sqlx.DB{}
	logger := zerolog.Nop()

	repo := NewRepository(db, logger)

	assert.NotNil(t, repo)
	assert.NotNil(t, repo.db)
}

func TestCreateListing(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	testCases := []struct {
		name      string
		input     *domain.CreateListingInput
		wantErr   bool
		errString string
	}{
		{
			name: "valid listing",
			input: &domain.CreateListingInput{
				UserID:       1,
				StorefrontID: nil,
				Title:        "Test Product",
				Description:  stringPtr("Test Description"),
				Price:        99.99,
				Currency:     "USD",
				CategoryID:   100,
				Quantity:     10,
				SKU:          stringPtr("TEST-SKU-001"),
			},
			wantErr: false,
		},
		{
			name: "valid listing with storefront",
			input: &domain.CreateListingInput{
				UserID:       2,
				StorefrontID: int64Ptr(123),
				Title:        "Storefront Product",
				Description:  stringPtr("Storefront Description"),
				Price:        199.99,
				Currency:     "EUR",
				CategoryID:   200,
				Quantity:     5,
				SKU:          stringPtr("STORE-SKU-001"),
			},
			wantErr: false,
		},
		{
			name: "missing required fields - invalid user",
			input: &domain.CreateListingInput{
				UserID:   0, // Invalid
				Title:    "",
				Price:    0,
				Currency: "USD",
			},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tests.TestContext(t)

			listing, err := repo.CreateListing(ctx, tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, listing)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, listing)
			assert.Greater(t, listing.ID, int64(0))
			assert.NotEmpty(t, listing.UUID)
			assert.Equal(t, tt.input.UserID, listing.UserID)
			assert.Equal(t, tt.input.Title, listing.Title)
			assert.Equal(t, tt.input.Description, listing.Description)
			assert.Equal(t, tt.input.Price, listing.Price)
			assert.Equal(t, tt.input.Currency, listing.Currency)
			assert.Equal(t, "draft", listing.Status) // Default status
			assert.NotZero(t, listing.CreatedAt)
			assert.NotZero(t, listing.UpdatedAt)
		})
	}
}

func TestGetListingByID(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	// Create a test listing first
	input := &domain.CreateListingInput{
		UserID:      1,
		Title:       "Test Listing",
		Description: stringPtr("Description"),
		Price:       50.00,
		Currency:    "USD",
		CategoryID:  100,
		Quantity:    1,
		SKU:         stringPtr("TEST-001"),
	}
	created, err := repo.CreateListing(ctx, input)
	require.NoError(t, err)

	testCases := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{
			name:    "existing listing",
			id:      created.ID,
			wantErr: false,
		},
		{
			name:    "non-existent listing",
			id:      99999,
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			listing, err := repo.GetListingByID(ctx, tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, listing)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, listing)
			assert.Equal(t, tt.id, listing.ID)
			assert.Equal(t, created.Title, listing.Title)
		})
	}
}

func TestUpdateListing(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	// Create a test listing
	input := &domain.CreateListingInput{
		UserID:      1,
		Title:       "Original Title",
		Description: stringPtr("Original Description"),
		Price:       100.00,
		Currency:    "USD",
		CategoryID:  100,
		Quantity:    5,
		SKU:         stringPtr("ORIG-001"),
	}
	created, err := repo.CreateListing(ctx, input)
	require.NoError(t, err)

	testCases := []struct {
		name    string
		id      int64
		update  *domain.UpdateListingInput
		wantErr bool
	}{
		{
			name: "update title and price",
			id:   created.ID,
			update: &domain.UpdateListingInput{
				Title:    stringPtr("Updated Title"),
				Price:    float64Ptr(150.00),
				Quantity: int32Ptr(10),
			},
			wantErr: false,
		},
		{
			name: "update non-existent listing",
			id:   99999,
			update: &domain.UpdateListingInput{
				Title: stringPtr("Should Fail"),
			},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			listing, err := repo.UpdateListing(ctx, tt.id, tt.update)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, listing)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, listing)

			if tt.update.Title != nil {
				assert.Equal(t, *tt.update.Title, listing.Title)
			}
			if tt.update.Price != nil {
				assert.Equal(t, *tt.update.Price, listing.Price)
			}
			if tt.update.Quantity != nil {
				assert.Equal(t, *tt.update.Quantity, listing.Quantity)
			}
		})
	}
}

func TestDeleteListing(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	// Create a test listing
	input := &domain.CreateListingInput{
		UserID:      1,
		Title:       "To Be Deleted",
		Description: stringPtr("Description"),
		Price:       75.00,
		Currency:    "USD",
		CategoryID:  100,
		Quantity:    3,
		SKU:         stringPtr("DEL-001"),
	}
	created, err := repo.CreateListing(ctx, input)
	require.NoError(t, err)

	testCases := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{
			name:    "delete existing listing",
			id:      created.ID,
			wantErr: false,
		},
		{
			name:    "delete non-existent listing",
			id:      99999,
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.DeleteListing(ctx, tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Verify listing is soft-deleted
			deleted, err := repo.GetListingByID(ctx, tt.id)
			assert.Error(t, err) // Should not find deleted listing
			assert.Nil(t, deleted)
		})
	}
}

func TestListListings(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	// Create multiple test listings
	for i := 0; i < 5; i++ {
		desc := stringPtr("Test Description")
		sku := stringPtr("LIST-SKU-" + string(rune('A'+i)))
		input := &domain.CreateListingInput{
			UserID:      1,
			Title:       "Test Listing " + string(rune('A'+i)),
			Description: desc,
			Price:       float64(100 * (i + 1)),
			Currency:    "USD",
			CategoryID:  100,
			Quantity:    int32(i + 1),
			SKU:         sku,
		}
		_, err := repo.CreateListing(ctx, input)
		require.NoError(t, err)
	}

	testCases := []struct {
		name      string
		filter    *domain.ListListingsFilter
		wantCount int
		wantErr   bool
	}{
		{
			name: "get all listings",
			filter: &domain.ListListingsFilter{
				Limit:  10,
				Offset: 0,
			},
			wantCount: 5,
			wantErr:   false,
		},
		{
			name: "get with pagination",
			filter: &domain.ListListingsFilter{
				Limit:  2,
				Offset: 0,
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "get specific user listings",
			filter: &domain.ListListingsFilter{
				UserID: int64Ptr(1),
				Limit:  10,
				Offset: 0,
			},
			wantCount: 5,
			wantErr:   false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			listings, total, err := repo.ListListings(ctx, tt.filter)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Len(t, listings, tt.wantCount)
			assert.GreaterOrEqual(t, total, int32(tt.wantCount))
		})
	}
}

func TestHealthCheck(t *testing.T) {
	repo, testDB := setupTestRepo(t)
	defer testDB.TeardownTestPostgres(t)

	ctx := tests.TestContext(t)

	err := repo.HealthCheck(ctx)
	assert.NoError(t, err)
}

// Helper functions
func int64Ptr(v int64) *int64 {
	return &v
}

func stringPtr(v string) *string {
	return &v
}

func float64Ptr(v float64) *float64 {
	return &v
}

func int32Ptr(v int32) *int32 {
	return &v
}
