// backend/internal/proj/c2c/service/marketplace_discount_parser_test.go
package service

import (
	"context"
	"testing"

	"backend/internal/domain/models"

	"github.com/stretchr/testify/assert"
)

func TestParseDiscountFromDescription(t *testing.T) {
	ctx := context.Background()
	service := &MarketplaceService{}

	tests := []struct {
		name                string
		listing             *models.MarketplaceListing
		wantValid           bool
		wantDiscountPercent int
		wantOldPrice        float64
		wantErrorContains   string
	}{
		{
			name: "Valid discount - 30% off",
			listing: &models.MarketplaceListing{
				ID:          1,
				Price:       70,
				Description: "Отличная куртка! 30% СКИДКА Старая цена: 100 RSD",
			},
			wantValid:           true,
			wantDiscountPercent: 30,
			wantOldPrice:        100,
		},
		{
			name: "Valid discount - with comma in price",
			listing: &models.MarketplaceListing{
				ID:          2,
				Price:       850,
				Description: "Супер товар! 15% СКИДКА! Старая цена: 1000,00 RSD",
			},
			wantValid:           true,
			wantDiscountPercent: 15,
			wantOldPrice:        1000,
		},
		{
			name: "Valid discount - with dot in price",
			listing: &models.MarketplaceListing{
				ID:          3,
				Price:       425,
				Description: "Акция! 50% СКИДКА Старая цена: 850.50 RSD",
			},
			wantValid:           true,
			wantDiscountPercent: 50,
			wantOldPrice:        850.50,
		},
		{
			name: "Invalid - no discount keyword",
			listing: &models.MarketplaceListing{
				ID:          4,
				Price:       100,
				Description: "Обычный товар без скидки",
			},
			wantValid: false,
		},
		{
			name: "Invalid - missing discount percent",
			listing: &models.MarketplaceListing{
				ID:          5,
				Price:       100,
				Description: "СКИДКА! Старая цена: 150 RSD",
			},
			wantValid:         false,
			wantErrorContains: "Не найдена полная информация",
		},
		{
			name: "Invalid - missing old price",
			listing: &models.MarketplaceListing{
				ID:          6,
				Price:       100,
				Description: "25% СКИДКА на товар!",
			},
			wantValid:         false,
			wantErrorContains: "Не найдена полная информация",
		},
		{
			name: "Invalid - discount doesn't match calculation",
			listing: &models.MarketplaceListing{
				ID:          7,
				Price:       90,
				Description: "50% СКИДКА Старая цена: 100 RSD", // На самом деле 10% скидка
			},
			wantValid:         false,
			wantErrorContains: "не соответствует реальной",
		},
		{
			name: "Invalid - negative discount (price increased)",
			listing: &models.MarketplaceListing{
				ID:          8,
				Price:       150,
				Description: "30% СКИДКА Старая цена: 100 RSD",
			},
			wantValid:         false,
			wantErrorContains: "не соответствует реальной",
		},
		{
			name: "Valid - discount within 5% tolerance",
			listing: &models.MarketplaceListing{
				ID:          9,
				Price:       67,
				Description: "30% СКИДКА Старая цена: 100 RSD", // Фактически 33% скидка
			},
			wantValid:           true,
			wantDiscountPercent: 30,
			wantOldPrice:        100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.ParseDiscountFromDescription(ctx, tt.listing)

			assert.NotNil(t, result)
			assert.Equal(t, tt.wantValid, result.IsValid, "IsValid mismatch")

			if tt.wantValid {
				assert.Equal(t, tt.wantDiscountPercent, result.DiscountPercent, "DiscountPercent mismatch")
				assert.Equal(t, tt.wantOldPrice, result.OldPrice, "OldPrice mismatch")
				assert.Empty(t, result.ValidationError)
			} else if tt.wantErrorContains != "" {
				assert.Contains(t, result.ValidationError, tt.wantErrorContains)
			}
		})
	}
}
