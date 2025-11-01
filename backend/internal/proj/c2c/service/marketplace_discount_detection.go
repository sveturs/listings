// backend/internal/proj/c2c/service/marketplace_discount_detection.go
package service

import (
	"context"
	"sort"
	"time"

	"backend/internal/domain/models"
	"backend/internal/logger"
)

// PriceManipulationResult содержит результат детектирования манипуляций с ценами
type PriceManipulationResult struct {
	IsManipulated bool
	Reason        string
}

// DetectPriceManipulation проверяет историю цен на наличие манипуляций
// Манипуляция определяется как:
// - Повышение цены более чем на 50%
// - Длительность повышенной цены менее 3 дней
// - Последующее снижение цены
//
// Эта практика используется для искусственного создания "скидок"
func (s *MarketplaceService) DetectPriceManipulation(
	ctx context.Context,
	listing *models.MarketplaceListing,
	priceHistory []models.PriceHistoryEntry,
) *PriceManipulationResult {
	// Если истории нет или записей мало - манипуляции быть не может
	if len(priceHistory) < 2 {
		return &PriceManipulationResult{
			IsManipulated: false,
		}
	}

	// Сортируем историю цен по дате (от старых к новым)
	sortedHistory := make([]models.PriceHistoryEntry, len(priceHistory))
	copy(sortedHistory, priceHistory)
	sort.Slice(sortedHistory, func(i, j int) bool {
		return sortedHistory[i].EffectiveFrom.Before(sortedHistory[j].EffectiveFrom)
	})

	// Проверяем каждый переход цены
	for i := 0; i < len(sortedHistory)-1; i++ {
		currentEntry := sortedHistory[i]
		nextEntry := sortedHistory[i+1]

		// Определяем когда закончилось действие следующей цены
		var nextEffectiveTo time.Time
		if nextEntry.EffectiveTo == nil {
			nextEffectiveTo = time.Now()
		} else {
			nextEffectiveTo = *nextEntry.EffectiveTo
		}

		// Рассчитываем длительность действия повышенной цены
		duration := nextEffectiveTo.Sub(nextEntry.EffectiveFrom)

		// Проверяем условия манипуляции:
		// 1. Цена выросла более чем на 50%
		priceIncrease := currentEntry.Price * 1.5
		if nextEntry.Price <= priceIncrease {
			continue
		}

		// 2. Повышенная цена действовала менее 3 дней
		if duration >= 3*24*time.Hour {
			continue
		}

		// 3. После повышения была записана более низкая цена
		if i+2 >= len(sortedHistory) {
			continue
		}
		subsequentEntry := sortedHistory[i+2]
		if subsequentEntry.Price >= nextEntry.Price {
			continue
		}

		// Все условия выполнены - это манипуляция!
		logger.Warn().
			Int("listing_id", listing.ID).
			Float64("price_before", currentEntry.Price).
			Float64("price_inflated", nextEntry.Price).
			Float64("price_after", subsequentEntry.Price).
			Float64("duration_hours", duration.Hours()).
			Msg("Detected price manipulation")

		return &PriceManipulationResult{
			IsManipulated: true,
			Reason:        "Обнаружена манипуляция с ценой: резкое повышение с последующим быстрым снижением",
		}
	}

	return &PriceManipulationResult{
		IsManipulated: false,
	}
}

// RemoveDiscountDueToManipulation удаляет информацию о скидке из метаданных
// и обновляет объявление в БД
func (s *MarketplaceService) RemoveDiscountDueToManipulation(
	ctx context.Context,
	listing *models.MarketplaceListing,
) error {
	// Если скидки нет - нечего удалять
	if listing.Metadata == nil || listing.Metadata["discount"] == nil {
		return nil
	}

	// Удаляем скидку из метаданных
	delete(listing.Metadata, "discount")

	// Обновляем объявление в БД
	_, err := s.storage.Exec(ctx, `
		UPDATE c2c_listings
		SET metadata = $1
		WHERE id = $2
	`, listing.Metadata, listing.ID)
	if err != nil {
		logger.Error().
			Err(err).
			Int("listing_id", listing.ID).
			Msg("Failed to remove discount metadata due to manipulation")
		return err
	}

	logger.Info().
		Int("listing_id", listing.ID).
		Msg("Removed discount due to price manipulation detection")

	return nil
}
