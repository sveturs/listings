package performance

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"

	"github.com/sveturs/listings/internal/domain"
	"github.com/sveturs/listings/internal/repository/postgres"
	"github.com/sveturs/listings/tests"
)

// BenchmarkCreateListing measures listing creation performance
func BenchmarkCreateListing(b *testing.B) {
	tests.SkipIfNoDocker(&testing.T{})

	testDB := tests.SetupTestPostgres(&testing.T{})
	defer testDB.TeardownTestPostgres(&testing.T{})

	tests.RunMigrations(&testing.T{}, testDB.DB, "../../migrations")

	db := sqlx.NewDb(testDB.DB, "postgres")
	logger := zerolog.Nop()
	repo := postgres.NewRepository(db, logger)

	ctx := context.Background()

	input := &domain.CreateListingInput{
		UserID:      1,
		Title:       "Benchmark Product",
		Description: stringPtr("Benchmark Description"),
		Price:       99.99,
		Currency:    "USD",
		CategoryID:  100,
		Quantity:    10,
		SKU:         stringPtr("BENCH-001"),
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := repo.CreateListing(ctx, input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGetListingByID measures listing retrieval performance
func BenchmarkGetListingByID(b *testing.B) {
	tests.SkipIfNoDocker(&testing.T{})

	testDB := tests.SetupTestPostgres(&testing.T{})
	defer testDB.TeardownTestPostgres(&testing.T{})

	tests.RunMigrations(&testing.T{}, testDB.DB, "../../migrations")

	db := sqlx.NewDb(testDB.DB, "postgres")
	logger := zerolog.Nop()
	repo := postgres.NewRepository(db, logger)

	ctx := context.Background()

	// Create a listing first
	input := &domain.CreateListingInput{
		UserID:      1,
		Title:       "Benchmark Product",
		Description: stringPtr("Benchmark Description"),
		Price:       99.99,
		Currency:    "USD",
		CategoryID:  100,
		Quantity:    10,
		SKU:         stringPtr("BENCH-GET-001"),
	}

	listing, err := repo.CreateListing(ctx, input)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := repo.GetListingByID(ctx, listing.ID)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkUpdateListing measures listing update performance
func BenchmarkUpdateListing(b *testing.B) {
	tests.SkipIfNoDocker(&testing.T{})

	testDB := tests.SetupTestPostgres(&testing.T{})
	defer testDB.TeardownTestPostgres(&testing.T{})

	tests.RunMigrations(&testing.T{}, testDB.DB, "../../migrations")

	db := sqlx.NewDb(testDB.DB, "postgres")
	logger := zerolog.Nop()
	repo := postgres.NewRepository(db, logger)

	ctx := context.Background()

	// Create a listing first
	input := &domain.CreateListingInput{
		UserID:      1,
		Title:       "Benchmark Product",
		Description: stringPtr("Benchmark Description"),
		Price:       99.99,
		Currency:    "USD",
		CategoryID:  100,
		Quantity:    10,
		SKU:         stringPtr("BENCH-UPDATE-001"),
	}

	listing, err := repo.CreateListing(ctx, input)
	if err != nil {
		b.Fatal(err)
	}

	update := &domain.UpdateListingInput{
		Title: stringPtr("Updated Title"),
		Price: float64Ptr(149.99),
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := repo.UpdateListing(ctx, listing.ID, update)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkListListings measures listing pagination performance
func BenchmarkListListings(b *testing.B) {
	tests.SkipIfNoDocker(&testing.T{})

	testDB := tests.SetupTestPostgres(&testing.T{})
	defer testDB.TeardownTestPostgres(&testing.T{})

	tests.RunMigrations(&testing.T{}, testDB.DB, "../../migrations")

	db := sqlx.NewDb(testDB.DB, "postgres")
	logger := zerolog.Nop()
	repo := postgres.NewRepository(db, logger)

	ctx := context.Background()

	// Create 100 listings for realistic benchmark
	for i := 0; i < 100; i++ {
		input := &domain.CreateListingInput{
			UserID:      int64(i % 10),
			Title:       "Benchmark Product",
			Description: stringPtr("Benchmark Description"),
			Price:       99.99,
			Currency:    "USD",
			CategoryID:  100,
			Quantity:    10,
			SKU:         stringPtr("BENCH-LIST-001"),
		}
		_, err := repo.CreateListing(ctx, input)
		if err != nil {
			b.Fatal(err)
		}
	}

	filter := &domain.ListListingsFilter{
		Limit:  10,
		Offset: 0,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, err := repo.ListListings(ctx, filter)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkParallelGetListing measures concurrent read performance
func BenchmarkParallelGetListing(b *testing.B) {
	tests.SkipIfNoDocker(&testing.T{})

	testDB := tests.SetupTestPostgres(&testing.T{})
	defer testDB.TeardownTestPostgres(&testing.T{})

	tests.RunMigrations(&testing.T{}, testDB.DB, "../../migrations")

	db := sqlx.NewDb(testDB.DB, "postgres")
	logger := zerolog.Nop()
	repo := postgres.NewRepository(db, logger)

	ctx := context.Background()

	// Create a listing
	input := &domain.CreateListingInput{
		UserID:      1,
		Title:       "Parallel Benchmark Product",
		Description: stringPtr("Parallel Benchmark Description"),
		Price:       99.99,
		Currency:    "USD",
		CategoryID:  100,
		Quantity:    10,
		SKU:         stringPtr("BENCH-PARALLEL-001"),
	}

	listing, err := repo.CreateListing(ctx, input)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := repo.GetListingByID(ctx, listing.ID)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}
