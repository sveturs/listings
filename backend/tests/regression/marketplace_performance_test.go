// backend/tests/regression/marketplace_performance_test.go
package regression

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"backend/internal/domain/search"

	"backend/internal/domain/models"
	"backend/internal/proj/unified/service"

	"github.com/rs/zerolog"
)

// BenchmarkGetListing_Monolith бенчмарк GetListing через monolith (local DB)
// defaultRoutingContext создаёт дефолтный routing context для тестов
func defaultRoutingContext() *service.RoutingContext {
	return &service.RoutingContext{
		UserID:  100,
		IsAdmin: false,
	}
}

func BenchmarkGetListing_Monolith(b *testing.B) {
	logger := zerolog.New(nil).Level(zerolog.Disabled)
	ctx := context.Background()

	mockC2C := &mockC2CRepository{}
	svc := service.NewMarketplaceService(mockC2C, nil, &mockOpenSearchRepository{}, logger)
	// Microservice disabled (default)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = svc.GetListing(ctx, int64(i%1000), service.SourceTypeC2C, defaultRoutingContext())
	}
}

// BenchmarkGetListing_Microservice бенчмарк GetListing через microservice
func BenchmarkGetListing_Microservice(b *testing.B) {
	logger := zerolog.New(nil).Level(zerolog.Disabled)
	ctx := context.Background()

	mockGRPC := service.NewMockListingsGRPCClient()
	svc := service.NewMarketplaceService(nil, nil, &mockOpenSearchRepository{}, logger)
	svc.SetListingsGRPCClient(mockGRPC, true)

	// Pre-populate some listings
	for i := 0; i < 1000; i++ {
		_, _ = mockGRPC.CreateListing(ctx, &models.UnifiedListing{
			UserID:     100,
			SourceType: service.SourceTypeC2C,
			Title:      fmt.Sprintf("Test %d", i),
			Price:      float64(i),
			CategoryID: 1,
		})
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = svc.GetListing(ctx, int64(i%1000+1), service.SourceTypeC2C, defaultRoutingContext())
	}
}

// BenchmarkCreateListing_Monolith бенчмарк CreateListing через monolith
func BenchmarkCreateListing_Monolith(b *testing.B) {
	logger := zerolog.New(nil).Level(zerolog.Disabled)
	ctx := context.Background()

	mockC2C := &mockC2CRepository{}
	svc := service.NewMarketplaceService(mockC2C, nil, &mockOpenSearchRepository{}, logger)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		unified := &models.UnifiedListing{
			UserID:     int(rand.Int31n(1000)),
			SourceType: service.SourceTypeC2C,
			Title:      fmt.Sprintf("Benchmark %d", i),
			Price:      float64(rand.Intn(1000)),
			CategoryID: rand.Intn(10) + 1,
		}
		_, _ = svc.CreateListing(ctx, unified, defaultRoutingContext())
	}
}

// BenchmarkCreateListing_Microservice бенчмарк CreateListing через microservice
func BenchmarkCreateListing_Microservice(b *testing.B) {
	logger := zerolog.New(nil).Level(zerolog.Disabled)
	ctx := context.Background()

	mockGRPC := service.NewMockListingsGRPCClient()
	svc := service.NewMarketplaceService(nil, nil, &mockOpenSearchRepository{}, logger)
	svc.SetListingsGRPCClient(mockGRPC, true)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		unified := &models.UnifiedListing{
			UserID:     int(rand.Int31n(1000)),
			SourceType: service.SourceTypeC2C,
			Title:      fmt.Sprintf("Benchmark %d", i),
			Price:      float64(rand.Intn(1000)),
			CategoryID: rand.Intn(10) + 1,
		}
		_, _ = svc.CreateListing(ctx, unified, defaultRoutingContext())
	}
}

// BenchmarkUpdateListing_Monolith бенчмарк UpdateListing через monolith
func BenchmarkUpdateListing_Monolith(b *testing.B) {
	logger := zerolog.New(nil).Level(zerolog.Disabled)
	ctx := context.Background()

	mockC2C := &mockC2CRepository{}
	svc := service.NewMarketplaceService(mockC2C, nil, &mockOpenSearchRepository{}, logger)

	// Pre-create a listing
	unified := &models.UnifiedListing{
		ID:         100,
		UserID:     100,
		SourceType: service.SourceTypeC2C,
		Title:      "Original",
		Price:      50.00,
		CategoryID: 1,
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		unified.Title = fmt.Sprintf("Updated %d", i)
		unified.Price = float64(i)
		_ = svc.UpdateListing(ctx, unified, defaultRoutingContext())
	}
}

// BenchmarkUpdateListing_Microservice бенчмарк UpdateListing через microservice
func BenchmarkUpdateListing_Microservice(b *testing.B) {
	logger := zerolog.New(nil).Level(zerolog.Disabled)
	ctx := context.Background()

	mockGRPC := service.NewMockListingsGRPCClient()
	svc := service.NewMarketplaceService(nil, nil, &mockOpenSearchRepository{}, logger)
	svc.SetListingsGRPCClient(mockGRPC, true)

	// Pre-create a listing
	unified, _ := mockGRPC.CreateListing(ctx, &models.UnifiedListing{
		UserID:     100,
		SourceType: service.SourceTypeC2C,
		Title:      "Original",
		Price:      50.00,
		CategoryID: 1,
	})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		unified.Title = fmt.Sprintf("Updated %d", i)
		unified.Price = float64(i)
		_ = svc.UpdateListing(ctx, unified, defaultRoutingContext())
	}
}

// BenchmarkDeleteListing_Monolith бенчмарк DeleteListing через monolith
func BenchmarkDeleteListing_Monolith(b *testing.B) {
	logger := zerolog.New(nil).Level(zerolog.Disabled)
	ctx := context.Background()

	mockC2C := &mockC2CRepository{}
	svc := service.NewMarketplaceService(mockC2C, nil, &mockOpenSearchRepository{}, logger)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = svc.DeleteListing(ctx, int64(i), service.SourceTypeC2C, defaultRoutingContext())
	}
}

// BenchmarkDeleteListing_Microservice бенчмарк DeleteListing через microservice
func BenchmarkDeleteListing_Microservice(b *testing.B) {
	logger := zerolog.New(nil).Level(zerolog.Disabled)
	ctx := context.Background()

	mockGRPC := service.NewMockListingsGRPCClient()
	svc := service.NewMarketplaceService(nil, nil, &mockOpenSearchRepository{}, logger)
	svc.SetListingsGRPCClient(mockGRPC, true)

	// Pre-create listings
	for i := 0; i < b.N; i++ {
		_, _ = mockGRPC.CreateListing(ctx, &models.UnifiedListing{
			UserID:     100,
			SourceType: service.SourceTypeC2C,
			Title:      fmt.Sprintf("Test %d", i),
			Price:      10.00,
			CategoryID: 1,
		})
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = svc.DeleteListing(ctx, int64(i+1), service.SourceTypeC2C, defaultRoutingContext())
	}
}

// BenchmarkConcurrentReads_Monolith бенчмарк конкурентного чтения (monolith)
func BenchmarkConcurrentReads_Monolith(b *testing.B) {
	logger := zerolog.New(nil).Level(zerolog.Disabled)
	ctx := context.Background()

	mockC2C := &mockC2CRepository{}
	svc := service.NewMarketplaceService(mockC2C, nil, &mockOpenSearchRepository{}, logger)

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			_, _ = svc.GetListing(ctx, int64(i%1000), service.SourceTypeC2C, defaultRoutingContext())
			i++
		}
	})
}

// BenchmarkConcurrentReads_Microservice бенчмарк конкурентного чтения (microservice)
func BenchmarkConcurrentReads_Microservice(b *testing.B) {
	logger := zerolog.New(nil).Level(zerolog.Disabled)
	ctx := context.Background()

	mockGRPC := service.NewMockListingsGRPCClient()
	svc := service.NewMarketplaceService(nil, nil, &mockOpenSearchRepository{}, logger)
	svc.SetListingsGRPCClient(mockGRPC, true)

	// Pre-populate
	for i := 0; i < 1000; i++ {
		_, _ = mockGRPC.CreateListing(ctx, &models.UnifiedListing{
			UserID:     100,
			SourceType: service.SourceTypeC2C,
			Title:      fmt.Sprintf("Test %d", i),
			Price:      float64(i),
			CategoryID: 1,
		})
	}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			_, _ = svc.GetListing(ctx, int64(i%1000+1), service.SourceTypeC2C, defaultRoutingContext())
			i++
		}
	})
}

// BenchmarkConcurrentWrites_Monolith бенчмарк конкурентной записи (monolith)
func BenchmarkConcurrentWrites_Monolith(b *testing.B) {
	logger := zerolog.New(nil).Level(zerolog.Disabled)
	ctx := context.Background()

	mockC2C := &mockC2CRepository{}
	svc := service.NewMarketplaceService(mockC2C, nil, &mockOpenSearchRepository{}, logger)

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			unified := &models.UnifiedListing{
				UserID:     100,
				SourceType: service.SourceTypeC2C,
				Title:      fmt.Sprintf("Concurrent %d", i),
				Price:      float64(i),
				CategoryID: 1,
			}
			_, _ = svc.CreateListing(ctx, unified, defaultRoutingContext())
			i++
		}
	})
}

// BenchmarkConcurrentWrites_Microservice бенчмарк конкурентной записи (microservice)
func BenchmarkConcurrentWrites_Microservice(b *testing.B) {
	logger := zerolog.New(nil).Level(zerolog.Disabled)
	ctx := context.Background()

	mockGRPC := service.NewMockListingsGRPCClient()
	svc := service.NewMarketplaceService(nil, nil, &mockOpenSearchRepository{}, logger)
	svc.SetListingsGRPCClient(mockGRPC, true)

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			unified := &models.UnifiedListing{
				UserID:     100,
				SourceType: service.SourceTypeC2C,
				Title:      fmt.Sprintf("Concurrent %d", i),
				Price:      float64(i),
				CategoryID: 1,
			}
			_, _ = svc.CreateListing(ctx, unified, defaultRoutingContext())
			i++
		}
	})
}

// BenchmarkFallback_Performance бенчмарк fallback производительности
func BenchmarkFallback_Performance(b *testing.B) {
	logger := zerolog.New(nil).Level(zerolog.Disabled)
	ctx := context.Background()

	mockC2C := &mockC2CRepository{}
	mockGRPC := service.NewMockListingsGRPCClient()
	mockGRPC.ShouldFailCreate = true // Force fallback

	svc := service.NewMarketplaceService(mockC2C, nil, &mockOpenSearchRepository{}, logger)
	svc.SetListingsGRPCClient(mockGRPC, true)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		unified := &models.UnifiedListing{
			UserID:     100,
			SourceType: service.SourceTypeC2C,
			Title:      fmt.Sprintf("Fallback %d", i),
			Price:      float64(i),
			CategoryID: 1,
		}
		_, _ = svc.CreateListing(ctx, unified, defaultRoutingContext())
	}
}

// BenchmarkMemoryAllocation_Monolith проверяет memory allocation (monolith)
func BenchmarkMemoryAllocation_Monolith(b *testing.B) {
	logger := zerolog.New(nil).Level(zerolog.Disabled)
	ctx := context.Background()

	mockC2C := &mockC2CRepository{}
	svc := service.NewMarketplaceService(mockC2C, nil, &mockOpenSearchRepository{}, logger)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		unified := &models.UnifiedListing{
			UserID:      100,
			SourceType:  service.SourceTypeC2C,
			Title:       "Memory Test",
			Description: "Testing memory allocation patterns",
			Price:       99.99,
			CategoryID:  1,
		}
		_, _ = svc.CreateListing(ctx, unified, defaultRoutingContext())
		_, _ = svc.GetListing(ctx, 1, service.SourceTypeC2C, defaultRoutingContext())
	}
}

// BenchmarkMemoryAllocation_Microservice проверяет memory allocation (microservice)
func BenchmarkMemoryAllocation_Microservice(b *testing.B) {
	logger := zerolog.New(nil).Level(zerolog.Disabled)
	ctx := context.Background()

	mockGRPC := service.NewMockListingsGRPCClient()
	svc := service.NewMarketplaceService(nil, nil, &mockOpenSearchRepository{}, logger)
	svc.SetListingsGRPCClient(mockGRPC, true)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		unified := &models.UnifiedListing{
			UserID:      100,
			SourceType:  service.SourceTypeC2C,
			Title:       "Memory Test",
			Description: "Testing memory allocation patterns",
			Price:       99.99,
			CategoryID:  1,
		}
		_, _ = svc.CreateListing(ctx, unified, defaultRoutingContext())
		_, _ = svc.GetListing(ctx, int64(i+1), service.SourceTypeC2C, defaultRoutingContext())
	}
}

// Mock repositories

type mockC2CRepository struct{}

func (m *mockC2CRepository) CreateListing(ctx context.Context, listing *models.MarketplaceListing) (int, error) {
	return 123, nil
}

func (m *mockC2CRepository) GetListing(ctx context.Context, id int) (*models.MarketplaceListing, error) {
	return &models.MarketplaceListing{ID: id, Title: "Mock Listing"}, nil
}

func (m *mockC2CRepository) UpdateListing(ctx context.Context, listing *models.MarketplaceListing) error {
	return nil
}

func (m *mockC2CRepository) DeleteListing(ctx context.Context, id int) error {
	return nil
}

type mockOpenSearchRepository struct{}

func (m *mockOpenSearchRepository) SearchListings(ctx context.Context, params *search.ServiceParams) (*search.ServiceResult, error) {
	return &search.ServiceResult{Items: []*models.MarketplaceListing{}, Total: 0, Page: 0, Size: 10}, nil
}

func (m *mockOpenSearchRepository) Index(ctx context.Context, listing *models.MarketplaceListing) error {
	return nil
}

func (m *mockOpenSearchRepository) Delete(ctx context.Context, listingID int) error {
	return nil
}
