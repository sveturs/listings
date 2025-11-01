// backend/internal/proj/c2c/service/marketplace_discount_sync_test.go
package service

import (
	"context"
	"testing"
	"time"

	"backend/internal/domain/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOpenSearch мок для OpenSearch (IndexListing)
type MockOpenSearch struct {
	mock.Mock
}

func (m *MockOpenSearch) IndexListing(ctx context.Context, listing *models.MarketplaceListing) error {
	args := m.Called(ctx, listing)
	return args.Error(0)
}

func TestSynchronizeDiscountData_NoHistory(t *testing.T) {
	ctx := context.Background()
	mockStorage := new(MockStorage)

	listing := &models.MarketplaceListing{
		ID:          1,
		Price:       100,
		Description: "Обычный товар",
		Metadata:    make(map[string]interface{}),
	}

	// Mock GetListingByID
	mockStorage.On("GetListingByID", ctx, 1).Return(listing, nil)
	// Mock GetPriceHistory - пустая история
	mockStorage.On("GetPriceHistory", ctx, 1).Return([]models.PriceHistoryEntry{}, nil)
	// Mock IndexListing
	mockStorage.On("IndexListing", ctx, listing).Return(nil)

	service := &MarketplaceService{
		storage: mockStorage,
	}

	err := service.SynchronizeDiscountData(ctx, 1)

	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
}

func TestSynchronizeDiscountData_WithManipulation(t *testing.T) {
	ctx := context.Background()
	mockStorage := new(MockStorage)

	listing := &models.MarketplaceListing{
		ID:    2,
		Price: 90,
		Metadata: map[string]interface{}{
			"discount": map[string]interface{}{
				"discount_percent": 50,
			},
		},
	}

	// История с манипуляцией
	priceHistory := []models.PriceHistoryEntry{
		{
			ID:            1,
			ListingID:     2,
			Price:         100,
			EffectiveFrom: time.Now().AddDate(0, 0, -10),
			EffectiveTo:   ptrTime(time.Now().AddDate(0, 0, -5)),
		},
		{
			ID:            2,
			ListingID:     2,
			Price:         200, // Повышение на 100%
			EffectiveFrom: time.Now().AddDate(0, 0, -5),
			EffectiveTo:   ptrTime(time.Now().AddDate(0, 0, -3)), // 2 дня
		},
		{
			ID:            3,
			ListingID:     2,
			Price:         90,
			EffectiveFrom: time.Now().AddDate(0, 0, -3),
			EffectiveTo:   nil,
		},
	}

	// Mock GetListingByID
	mockStorage.On("GetListingByID", ctx, 2).Return(listing, nil)
	// Mock GetPriceHistory
	mockStorage.On("GetPriceHistory", ctx, 2).Return(priceHistory, nil)
	// Mock Exec для удаления скидки
	mockStorage.On("Exec", ctx, mock.Anything, mock.Anything).Return(nil, nil)
	// Mock IndexListing
	mockStorage.On("IndexListing", ctx, listing).Return(nil)

	service := &MarketplaceService{
		storage: mockStorage,
	}

	err := service.SynchronizeDiscountData(ctx, 2)

	assert.NoError(t, err)
	// Проверяем что скидка была удалена
	assert.Nil(t, listing.Metadata["discount"])
	mockStorage.AssertExpectations(t)
}

func TestSynchronizeDiscountData_WithCalculatedDiscount(t *testing.T) {
	ctx := context.Background()
	mockStorage := new(MockStorage)

	listing := &models.MarketplaceListing{
		ID:          3,
		Price:       70,
		Description: "Товар со скидкой",
		Metadata:    make(map[string]interface{}),
	}

	// Легитимная история со скидкой
	priceHistory := []models.PriceHistoryEntry{
		{
			ID:            1,
			ListingID:     3,
			Price:         100,
			EffectiveFrom: time.Now().AddDate(0, 0, -30),
			EffectiveTo:   ptrTime(time.Now().AddDate(0, 0, -5)),
		},
		{
			ID:            2,
			ListingID:     3,
			Price:         70,
			EffectiveFrom: time.Now().AddDate(0, 0, -5),
			EffectiveTo:   nil,
		},
	}

	// Mock GetListingByID
	mockStorage.On("GetListingByID", ctx, 3).Return(listing, nil)
	// Mock GetPriceHistory
	mockStorage.On("GetPriceHistory", ctx, 3).Return(priceHistory, nil)
	// Mock Exec для установки скидки
	mockStorage.On("Exec", ctx, mock.Anything, mock.Anything).Return(nil, nil)
	// Mock IndexListing
	mockStorage.On("IndexListing", ctx, listing).Return(nil)

	service := &MarketplaceService{
		storage: mockStorage,
	}

	err := service.SynchronizeDiscountData(ctx, 3)

	assert.NoError(t, err)
	// Проверяем что скидка была установлена
	assert.NotNil(t, listing.Metadata["discount"])
	mockStorage.AssertExpectations(t)
}

func TestSynchronizeDiscountData_WithParsedDiscount(t *testing.T) {
	ctx := context.Background()
	mockStorage := new(MockStorage)

	listing := &models.MarketplaceListing{
		ID:          4,
		Price:       70,
		Description: "Отличный товар! 30% СКИДКА Старая цена: 100 RSD",
		Metadata:    make(map[string]interface{}),
		CreatedAt:   time.Now().AddDate(0, 0, -10),
	}

	// Пустая история цен
	priceHistory := []models.PriceHistoryEntry{}

	// Mock GetListingByID
	mockStorage.On("GetListingByID", ctx, 4).Return(listing, nil)
	// Mock GetPriceHistory
	mockStorage.On("GetPriceHistory", ctx, 4).Return(priceHistory, nil)
	// Mock ClosePriceHistoryEntry
	mockStorage.On("ClosePriceHistoryEntry", ctx, 4).Return(nil)
	// Mock AddPriceHistoryEntry для старой цены
	mockStorage.On("AddPriceHistoryEntry", ctx, mock.MatchedBy(func(entry *models.PriceHistoryEntry) bool {
		return entry.Price == 100 && entry.ListingID == 4
	})).Return(nil)
	// Mock AddPriceHistoryEntry для новой цены
	mockStorage.On("AddPriceHistoryEntry", ctx, mock.MatchedBy(func(entry *models.PriceHistoryEntry) bool {
		return entry.Price == 70 && entry.ListingID == 4
	})).Return(nil)
	// Mock Exec для установки metadata
	mockStorage.On("Exec", ctx, mock.Anything, mock.Anything).Return(nil, nil)
	// Mock IndexListing
	mockStorage.On("IndexListing", ctx, listing).Return(nil)

	service := &MarketplaceService{
		storage: mockStorage,
	}

	err := service.SynchronizeDiscountData(ctx, 4)

	assert.NoError(t, err)
	// Проверяем что скидка была установлена
	assert.NotNil(t, listing.Metadata["discount"])
	mockStorage.AssertExpectations(t)
}

func TestSynchronizeDiscountData_NoDiscount(t *testing.T) {
	ctx := context.Background()
	mockStorage := new(MockStorage)

	listing := &models.MarketplaceListing{
		ID:          5,
		Price:       100,
		Description: "Обычный товар без скидки",
		Metadata:    make(map[string]interface{}),
	}

	// История без скидки (цена не менялась значительно)
	priceHistory := []models.PriceHistoryEntry{
		{
			ID:            1,
			ListingID:     5,
			Price:         102,
			EffectiveFrom: time.Now().AddDate(0, 0, -30),
			EffectiveTo:   ptrTime(time.Now().AddDate(0, 0, -5)),
		},
		{
			ID:            2,
			ListingID:     5,
			Price:         100,
			EffectiveFrom: time.Now().AddDate(0, 0, -5),
			EffectiveTo:   nil,
		},
	}

	// Mock GetListingByID
	mockStorage.On("GetListingByID", ctx, 5).Return(listing, nil)
	// Mock GetPriceHistory
	mockStorage.On("GetPriceHistory", ctx, 5).Return(priceHistory, nil)
	// Mock Exec для удаления маленькой скидки (если будет)
	mockStorage.On("Exec", ctx, mock.Anything, mock.Anything).Return(nil, nil)
	// Mock IndexListing
	mockStorage.On("IndexListing", ctx, listing).Return(nil)

	service := &MarketplaceService{
		storage: mockStorage,
	}

	err := service.SynchronizeDiscountData(ctx, 5)

	assert.NoError(t, err)
	// Проверяем что скидки нет
	assert.Nil(t, listing.Metadata["discount"])
	mockStorage.AssertExpectations(t)
}
