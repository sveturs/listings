// backend/internal/proj/c2c/service/marketplace_discount_calculator_test.go
package service

import (
	"context"
	"testing"
	"time"

	"backend/internal/domain/models"

	"github.com/stretchr/testify/assert"
)

func TestCalculateDiscountFromHistory(t *testing.T) {
	ctx := context.Background()
	service := &MarketplaceService{}

	tests := []struct {
		name                string
		listing             *models.MarketplaceListing
		priceHistory        []models.PriceHistoryEntry
		wantValid           bool
		wantDiscountPercent int
		wantPreviousPrice   float64
	}{
		{
			name: "Valid discount - price decreased from max",
			listing: &models.MarketplaceListing{
				ID:    1,
				Price: 70,
			},
			priceHistory: []models.PriceHistoryEntry{
				{
					ID:            1,
					ListingID:     1,
					Price:         100,
					EffectiveFrom: time.Now().AddDate(0, 0, -30),
					EffectiveTo:   ptrTime(time.Now().AddDate(0, 0, -5)),
				},
				{
					ID:            2,
					ListingID:     1,
					Price:         70,
					EffectiveFrom: time.Now().AddDate(0, 0, -5),
					EffectiveTo:   nil,
				},
			},
			wantValid:           true,
			wantDiscountPercent: 30, // (100-70)/100 = 30%
			wantPreviousPrice:   100,
		},
		{
			name: "Valid discount - multiple price changes",
			listing: &models.MarketplaceListing{
				ID:    2,
				Price: 60,
			},
			priceHistory: []models.PriceHistoryEntry{
				{
					ID:            1,
					ListingID:     2,
					Price:         80,
					EffectiveFrom: time.Now().AddDate(0, 0, -50),
					EffectiveTo:   ptrTime(time.Now().AddDate(0, 0, -30)),
				},
				{
					ID:            2,
					ListingID:     2,
					Price:         120, // Max price
					EffectiveFrom: time.Now().AddDate(0, 0, -30),
					EffectiveTo:   ptrTime(time.Now().AddDate(0, 0, -10)),
				},
				{
					ID:            3,
					ListingID:     2,
					Price:         60,
					EffectiveFrom: time.Now().AddDate(0, 0, -10),
					EffectiveTo:   nil,
				},
			},
			wantValid:           true,
			wantDiscountPercent: 50, // (120-60)/120 = 50%
			wantPreviousPrice:   120,
		},
		{
			name: "Invalid - no price decrease",
			listing: &models.MarketplaceListing{
				ID:    3,
				Price: 100,
			},
			priceHistory: []models.PriceHistoryEntry{
				{
					ID:            1,
					ListingID:     3,
					Price:         80,
					EffectiveFrom: time.Now().AddDate(0, 0, -30),
					EffectiveTo:   ptrTime(time.Now().AddDate(0, 0, -5)),
				},
				{
					ID:            2,
					ListingID:     3,
					Price:         100,
					EffectiveFrom: time.Now().AddDate(0, 0, -5),
					EffectiveTo:   nil,
				},
			},
			wantValid: false, // Current price >= max price
		},
		{
			name: "Invalid - empty price history",
			listing: &models.MarketplaceListing{
				ID:    4,
				Price: 100,
			},
			priceHistory: []models.PriceHistoryEntry{},
			wantValid:    false,
		},
		{
			name: "Invalid - price too short duration (ignored)",
			listing: &models.MarketplaceListing{
				ID:    5,
				Price: 50,
			},
			priceHistory: []models.PriceHistoryEntry{
				{
					ID:            1,
					ListingID:     5,
					Price:         200, // High price but lasted < 24h
					EffectiveFrom: time.Now().Add(-12 * time.Hour),
					EffectiveTo:   ptrTime(time.Now().Add(-6 * time.Hour)),
				},
				{
					ID:            2,
					ListingID:     5,
					Price:         100, // This is max (lasted > 24h)
					EffectiveFrom: time.Now().AddDate(0, 0, -30),
					EffectiveTo:   ptrTime(time.Now().Add(-13 * time.Hour)),
				},
				{
					ID:            3,
					ListingID:     5,
					Price:         50,
					EffectiveFrom: time.Now().Add(-5 * time.Hour),
					EffectiveTo:   nil,
				},
			},
			wantValid:           true,
			wantDiscountPercent: 50, // (100-50)/100 = 50%
			wantPreviousPrice:   100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.CalculateDiscountFromHistory(ctx, tt.listing, tt.priceHistory)

			assert.NotNil(t, result)
			assert.Equal(t, tt.wantValid, result.IsValid, "IsValid mismatch")

			if tt.wantValid {
				assert.Equal(t, tt.wantDiscountPercent, result.DiscountPercent, "DiscountPercent mismatch")
				assert.Equal(t, tt.wantPreviousPrice, result.PreviousPrice, "PreviousPrice mismatch")
			}
		})
	}
}

func TestIsDiscountSignificant(t *testing.T) {
	tests := []struct {
		name            string
		discountPercent int
		wantSignificant bool
	}{
		{
			name:            "Significant - exactly 5%",
			discountPercent: 5,
			wantSignificant: true,
		},
		{
			name:            "Significant - 30%",
			discountPercent: 30,
			wantSignificant: true,
		},
		{
			name:            "Not significant - 4%",
			discountPercent: 4,
			wantSignificant: false,
		},
		{
			name:            "Not significant - 0%",
			discountPercent: 0,
			wantSignificant: false,
		},
		{
			name:            "Significant - 100%",
			discountPercent: 100,
			wantSignificant: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsDiscountSignificant(tt.discountPercent)
			assert.Equal(t, tt.wantSignificant, result)
		})
	}
}

func TestCreateDiscountMetadata(t *testing.T) {
	effectiveFrom := time.Date(2024, 10, 1, 12, 0, 0, 0, time.UTC)

	discount := &CalculatedDiscount{
		DiscountPercent: 30,
		PreviousPrice:   100,
		EffectiveFrom:   effectiveFrom,
		IsValid:         true,
	}

	metadata := CreateDiscountMetadata(discount)

	assert.NotNil(t, metadata)
	assert.Equal(t, 30, metadata["discount_percent"])
	assert.Equal(t, 100.0, metadata["previous_price"])
	assert.Equal(t, effectiveFrom.Format(time.RFC3339), metadata["effective_from"])
	assert.Equal(t, true, metadata["has_price_history"])
}
