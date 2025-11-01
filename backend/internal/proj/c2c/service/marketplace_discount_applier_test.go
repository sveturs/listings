// backend/internal/proj/c2c/service/marketplace_discount_applier_test.go
package service

import (
	"context"
	"testing"
	"time"

	"backend/internal/domain/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestApplyDiscountMetadata_Set(t *testing.T) {
	ctx := context.Background()
	mockStorage := new(MockStorage)

	service := &MarketplaceService{
		storage: mockStorage,
	}

	listing := &models.MarketplaceListing{
		ID:       1,
		Price:    70,
		Metadata: make(map[string]interface{}),
	}

	discountData := map[string]interface{}{
		"discount_percent": 30,
		"previous_price":   100.0,
		"effective_from":   time.Now().Format(time.RFC3339),
	}

	// Mock успешного обновления БД
	mockStorage.On("Exec", ctx, mock.Anything, mock.Anything).Return(nil, nil)

	err := service.ApplyDiscountMetadata(ctx, listing, ApplyDiscountActionSet, discountData)

	assert.NoError(t, err)
	assert.NotNil(t, listing.Metadata["discount"])
	assert.Equal(t, discountData, listing.Metadata["discount"])
	mockStorage.AssertExpectations(t)
}

func TestApplyDiscountMetadata_Remove(t *testing.T) {
	ctx := context.Background()
	mockStorage := new(MockStorage)

	service := &MarketplaceService{
		storage: mockStorage,
	}

	listing := &models.MarketplaceListing{
		ID:    1,
		Price: 100,
		Metadata: map[string]interface{}{
			"discount": map[string]interface{}{
				"discount_percent": 30,
			},
		},
	}

	// Mock успешного обновления БД
	mockStorage.On("Exec", ctx, mock.Anything, mock.Anything).Return(nil, nil)

	err := service.ApplyDiscountMetadata(ctx, listing, ApplyDiscountActionRemove, nil)

	assert.NoError(t, err)
	assert.Nil(t, listing.Metadata["discount"])
	mockStorage.AssertExpectations(t)
}

func TestApplyCalculatedDiscount_ValidSignificantDiscount(t *testing.T) {
	ctx := context.Background()
	mockStorage := new(MockStorage)

	service := &MarketplaceService{
		storage: mockStorage,
	}

	listing := &models.MarketplaceListing{
		ID:       1,
		Price:    70,
		Metadata: make(map[string]interface{}),
	}

	discount := &CalculatedDiscount{
		DiscountPercent: 30,
		PreviousPrice:   100,
		EffectiveFrom:   time.Now(),
		IsValid:         true,
	}

	// Mock успешного обновления БД
	mockStorage.On("Exec", ctx, mock.Anything, mock.Anything).Return(nil, nil)

	err := service.ApplyCalculatedDiscount(ctx, listing, discount)

	assert.NoError(t, err)
	assert.NotNil(t, listing.Metadata["discount"])
	mockStorage.AssertExpectations(t)
}

func TestApplyCalculatedDiscount_InvalidDiscount(t *testing.T) {
	ctx := context.Background()
	mockStorage := new(MockStorage)

	service := &MarketplaceService{
		storage: mockStorage,
	}

	listing := &models.MarketplaceListing{
		ID:    1,
		Price: 100,
		Metadata: map[string]interface{}{
			"discount": map[string]interface{}{
				"discount_percent": 10,
			},
		},
	}

	discount := &CalculatedDiscount{
		IsValid: false,
	}

	// Mock успешного обновления БД (для удаления скидки)
	mockStorage.On("Exec", ctx, mock.Anything, mock.Anything).Return(nil, nil)

	err := service.ApplyCalculatedDiscount(ctx, listing, discount)

	assert.NoError(t, err)
	assert.Nil(t, listing.Metadata["discount"])
	mockStorage.AssertExpectations(t)
}

func TestApplyCalculatedDiscount_TooSmallDiscount(t *testing.T) {
	ctx := context.Background()
	mockStorage := new(MockStorage)

	service := &MarketplaceService{
		storage: mockStorage,
	}

	listing := &models.MarketplaceListing{
		ID:       1,
		Price:    97,
		Metadata: make(map[string]interface{}),
	}

	discount := &CalculatedDiscount{
		DiscountPercent: 3, // < 5%
		PreviousPrice:   100,
		EffectiveFrom:   time.Now(),
		IsValid:         true,
	}

	// Mock успешного обновления БД (для удаления скидки)
	mockStorage.On("Exec", ctx, mock.Anything, mock.Anything).Return(nil, nil)

	err := service.ApplyCalculatedDiscount(ctx, listing, discount)

	assert.NoError(t, err)
	assert.Nil(t, listing.Metadata["discount"])
	mockStorage.AssertExpectations(t)
}

func TestApplyParsedDiscount_ValidDiscount(t *testing.T) {
	ctx := context.Background()
	mockStorage := new(MockStorage)

	service := &MarketplaceService{
		storage: mockStorage,
	}

	createdAt := time.Now().AddDate(0, 0, -10)
	listing := &models.MarketplaceListing{
		ID:        1,
		Price:     70,
		Metadata:  make(map[string]interface{}),
		CreatedAt: createdAt,
	}

	discount := &ParsedDiscount{
		DiscountPercent: 30,
		OldPrice:        100,
		IsValid:         true,
	}

	// Mock успешного закрытия истории
	mockStorage.On("ClosePriceHistoryEntry", ctx, 1).Return(nil)
	// Mock успешного добавления старой цены
	mockStorage.On("AddPriceHistoryEntry", ctx, mock.MatchedBy(func(entry *models.PriceHistoryEntry) bool {
		return entry.Price == 100 && entry.ListingID == 1
	})).Return(nil)
	// Mock успешного добавления новой цены
	mockStorage.On("AddPriceHistoryEntry", ctx, mock.MatchedBy(func(entry *models.PriceHistoryEntry) bool {
		return entry.Price == 70 && entry.ListingID == 1
	})).Return(nil)
	// Mock успешного обновления metadata
	mockStorage.On("Exec", ctx, mock.Anything, mock.Anything).Return(nil, nil)

	err := service.ApplyParsedDiscount(ctx, listing, discount)

	assert.NoError(t, err)
	assert.NotNil(t, listing.Metadata["discount"])
	mockStorage.AssertExpectations(t)
}

func TestApplyParsedDiscount_InvalidDiscount(t *testing.T) {
	ctx := context.Background()
	service := &MarketplaceService{}

	listing := &models.MarketplaceListing{
		ID:       1,
		Price:    100,
		Metadata: make(map[string]interface{}),
	}

	discount := &ParsedDiscount{
		IsValid:         false,
		ValidationError: "Invalid format",
	}

	err := service.ApplyParsedDiscount(ctx, listing, discount)

	assert.NoError(t, err) // Функция не возвращает ошибку для invalid discount
	assert.Empty(t, listing.Metadata)
}
