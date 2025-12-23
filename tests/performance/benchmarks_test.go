package performance

import (
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"

	"github.com/vondi-global/listings/internal/domain"
	"github.com/vondi-global/listings/internal/repository/postgres"
	"github.com/vondi-global/listings/tests"
)

// BenchmarkCreateListing measures listing creation performance
func BenchmarkCreateListing(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping Docker benchmark in short mode")
	}

	tests.SkipIfNoDocker(b)

	testDB := tests.SetupTestPostgres(b)
	defer testDB.TeardownTestPostgres(b)

	tests.RunMigrations(b, testDB.DB, "../../migrations")

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
		CategoryID:  "3b4246cc-9970-403c-af01-c142a4178dc6",
		Quantity:    10,
		SKU:         stringPtr("BENCH-001"),
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := repo.CreateListing(ctx, input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkGetListingByID measures listing retrieval performance
func BenchmarkGetListingByID(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping Docker benchmark in short mode")
	}

	tests.SkipIfNoDocker(b)

	testDB := tests.SetupTestPostgres(b)
	defer testDB.TeardownTestPostgres(b)

	tests.RunMigrations(b, testDB.DB, "../../migrations")

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
		CategoryID:  "3b4246cc-9970-403c-af01-c142a4178dc6",
		Quantity:    10,
		SKU:         stringPtr("BENCH-GET-001"),
	}

	listing, err := repo.CreateListing(ctx, input)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := repo.GetListingByID(ctx, listing.ID)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkUpdateListing measures listing update performance
func BenchmarkUpdateListing(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping Docker benchmark in short mode")
	}

	tests.SkipIfNoDocker(b)

	testDB := tests.SetupTestPostgres(b)
	defer testDB.TeardownTestPostgres(b)

	tests.RunMigrations(b, testDB.DB, "../../migrations")

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
		CategoryID:  "3b4246cc-9970-403c-af01-c142a4178dc6",
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
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := repo.UpdateListing(ctx, listing.ID, update)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkListListings measures listing pagination performance
func BenchmarkListListings(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping Docker benchmark in short mode")
	}

	tests.SkipIfNoDocker(b)

	testDB := tests.SetupTestPostgres(b)
	defer testDB.TeardownTestPostgres(b)

	tests.RunMigrations(b, testDB.DB, "../../migrations")

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
			CategoryID:  "3b4246cc-9970-403c-af01-c142a4178dc6",
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
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _, err := repo.ListListings(ctx, filter)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkParallelGetListing measures concurrent read performance
func BenchmarkParallelGetListing(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping Docker benchmark in short mode")
	}

	tests.SkipIfNoDocker(b)

	testDB := tests.SetupTestPostgres(b)
	defer testDB.TeardownTestPostgres(b)

	tests.RunMigrations(b, testDB.DB, "../../migrations")

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
		CategoryID:  "3b4246cc-9970-403c-af01-c142a4178dc6",
		Quantity:    10,
		SKU:         stringPtr("BENCH-PARALLEL-001"),
	}

	listing, err := repo.CreateListing(ctx, input)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.ReportAllocs()

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
