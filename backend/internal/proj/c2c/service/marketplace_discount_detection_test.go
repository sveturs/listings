// backend/internal/proj/c2c/service/marketplace_discount_detection_test.go
package service

import (
	"context"
	"testing"
	"time"

	"backend/internal/domain/models"

	"github.com/stretchr/testify/assert"
)

func TestDetectPriceManipulation(t *testing.T) {
	ctx := context.Background()
	service := &MarketplaceService{}

	tests := []struct {
		name             string
		listing          *models.MarketplaceListing
		priceHistory     []models.PriceHistoryEntry
		wantManipulation bool
		wantReason       string
	}{
		{
			name: "No manipulation - stable price history",
			listing: &models.MarketplaceListing{
				ID:    1,
				Price: 100,
			},
			priceHistory: []models.PriceHistoryEntry{
				{
					ID:            1,
					ListingID:     1,
					Price:         100,
					EffectiveFrom: time.Now().AddDate(0, 0, -30),
					EffectiveTo:   ptrTime(time.Now().AddDate(0, 0, -10)),
				},
				{
					ID:            2,
					ListingID:     1,
					Price:         95,
					EffectiveFrom: time.Now().AddDate(0, 0, -10),
					EffectiveTo:   nil,
				},
			},
			wantManipulation: false,
		},
		{
			name: "Manipulation detected - price spike with quick drop",
			listing: &models.MarketplaceListing{
				ID:    2,
				Price: 100,
			},
			priceHistory: []models.PriceHistoryEntry{
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
					EffectiveTo:   ptrTime(time.Now().AddDate(0, 0, -3)), // 2 дня (< 3 дней)
				},
				{
					ID:            3,
					ListingID:     2,
					Price:         90, // Снижение ниже первоначальной
					EffectiveFrom: time.Now().AddDate(0, 0, -3),
					EffectiveTo:   nil,
				},
			},
			wantManipulation: true,
			wantReason:       "Обнаружена манипуляция с ценой",
		},
		{
			name: "No manipulation - price spike lasted long enough",
			listing: &models.MarketplaceListing{
				ID:    3,
				Price: 100,
			},
			priceHistory: []models.PriceHistoryEntry{
				{
					ID:            1,
					ListingID:     3,
					Price:         100,
					EffectiveFrom: time.Now().AddDate(0, 0, -20),
					EffectiveTo:   ptrTime(time.Now().AddDate(0, 0, -10)),
				},
				{
					ID:            2,
					ListingID:     3,
					Price:         200,
					EffectiveFrom: time.Now().AddDate(0, 0, -10),
					EffectiveTo:   ptrTime(time.Now().AddDate(0, 0, -3)), // 7 дней (> 3 дней)
				},
				{
					ID:            3,
					ListingID:     3,
					Price:         90,
					EffectiveFrom: time.Now().AddDate(0, 0, -3),
					EffectiveTo:   nil,
				},
			},
			wantManipulation: false,
		},
		{
			name: "No manipulation - price increase < 50%",
			listing: &models.MarketplaceListing{
				ID:    4,
				Price: 100,
			},
			priceHistory: []models.PriceHistoryEntry{
				{
					ID:            1,
					ListingID:     4,
					Price:         100,
					EffectiveFrom: time.Now().AddDate(0, 0, -10),
					EffectiveTo:   ptrTime(time.Now().AddDate(0, 0, -5)),
				},
				{
					ID:            2,
					ListingID:     4,
					Price:         120, // Повышение на 20% (< 50%)
					EffectiveFrom: time.Now().AddDate(0, 0, -5),
					EffectiveTo:   ptrTime(time.Now().AddDate(0, 0, -3)),
				},
				{
					ID:            3,
					ListingID:     4,
					Price:         90,
					EffectiveFrom: time.Now().AddDate(0, 0, -3),
					EffectiveTo:   nil,
				},
			},
			wantManipulation: false,
		},
		{
			name: "No manipulation - insufficient history",
			listing: &models.MarketplaceListing{
				ID:    5,
				Price: 100,
			},
			priceHistory: []models.PriceHistoryEntry{
				{
					ID:            1,
					ListingID:     5,
					Price:         100,
					EffectiveFrom: time.Now().AddDate(0, 0, -5),
					EffectiveTo:   nil,
				},
			},
			wantManipulation: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.DetectPriceManipulation(ctx, tt.listing, tt.priceHistory)

			assert.NotNil(t, result)
			assert.Equal(t, tt.wantManipulation, result.IsManipulated)

			if tt.wantManipulation {
				assert.Contains(t, result.Reason, tt.wantReason)
			}
		})
	}
}

// Helper function to create time pointers
func ptrTime(t time.Time) *time.Time {
	return &t
}
