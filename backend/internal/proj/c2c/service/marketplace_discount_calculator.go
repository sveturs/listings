// backend/internal/proj/c2c/service/marketplace_discount_calculator.go
package service

import (
	"context"
	"time"

	"backend/internal/domain/models"
	"backend/internal/logger"
)

// CalculatedDiscount содержит результат расчёта скидки из истории цен
type CalculatedDiscount struct {
	DiscountPercent int
	PreviousPrice   float64
	EffectiveFrom   time.Time
	IsValid         bool
}

// CalculateDiscountFromHistory рассчитывает скидку на основе истории цен
// Находит максимальную цену в истории (с учётом минимальной длительности)
// и рассчитывает процент скидки относительно текущей цены
func (s *MarketplaceService) CalculateDiscountFromHistory(
	ctx context.Context,
	listing *models.MarketplaceListing,
	priceHistory []models.PriceHistoryEntry,
) *CalculatedDiscount {
	// Если истории нет - скидку рассчитать невозможно
	if len(priceHistory) == 0 {
		return &CalculatedDiscount{
			IsValid: false,
		}
	}

	// Минимальная длительность действия цены - 1 день
	// Это предотвращает учёт кратковременных манипуляций
	minDuration := 24 * time.Hour

	// Ищем максимальную цену среди тех, что действовали достаточно долго
	var maxPrice float64
	var maxPriceDate time.Time

	for _, entry := range priceHistory {
		// Рассчитываем длительность действия цены
		var duration time.Duration
		if entry.EffectiveTo == nil {
			// Если цена всё ещё действует - считаем от начала до сейчас
			duration = time.Since(entry.EffectiveFrom)
		} else {
			// Если цена уже не действует - считаем полную длительность
			duration = entry.EffectiveTo.Sub(entry.EffectiveFrom)
		}

		// Учитываем только цены, которые действовали достаточно долго
		if duration >= minDuration && entry.Price > maxPrice {
			maxPrice = entry.Price
			maxPriceDate = entry.EffectiveFrom
		}
	}

	// Если максимальная цена не найдена или не превышает текущую - скидки нет
	if maxPrice == 0 || maxPrice <= listing.Price {
		return &CalculatedDiscount{
			IsValid: false,
		}
	}

	// Рассчитываем процент скидки
	discountPercent := int((maxPrice - listing.Price) / maxPrice * 100)

	logger.Debug().
		Int("listing_id", listing.ID).
		Float64("max_price", maxPrice).
		Float64("current_price", listing.Price).
		Int("discount_percent", discountPercent).
		Msg("Calculated discount from price history")

	return &CalculatedDiscount{
		DiscountPercent: discountPercent,
		PreviousPrice:   maxPrice,
		EffectiveFrom:   maxPriceDate,
		IsValid:         true,
	}
}

// IsDiscountSignificant проверяет что скидка достаточно большая чтобы показывать её
// Минимальный порог - 5%
func IsDiscountSignificant(discountPercent int) bool {
	const minDiscountPercent = 5
	return discountPercent >= minDiscountPercent
}

// CreateDiscountMetadata создаёт map с метаданными о скидке
// для сохранения в поле metadata объявления
func CreateDiscountMetadata(discount *CalculatedDiscount) map[string]interface{} {
	return map[string]interface{}{
		"discount_percent":  discount.DiscountPercent,
		"previous_price":    discount.PreviousPrice,
		"effective_from":    discount.EffectiveFrom.Format(time.RFC3339),
		"has_price_history": true,
	}
}

// CreateDiscountMetadataFromParsed создаёт map с метаданными о скидке
// из распарсенных данных (для скидок из description)
func CreateDiscountMetadataFromParsed(discount *ParsedDiscount, effectiveFrom time.Time) map[string]interface{} {
	return map[string]interface{}{
		"discount_percent":  discount.DiscountPercent,
		"previous_price":    discount.OldPrice,
		"effective_from":    effectiveFrom.Format(time.RFC3339),
		"has_price_history": true,
	}
}
