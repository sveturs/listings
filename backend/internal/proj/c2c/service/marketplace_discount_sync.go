// backend/internal/proj/c2c/service/marketplace_discount_sync.go
package service

import (
	"context"
	"fmt"

	"backend/internal/domain/models"
	"backend/internal/logger"
)

// SynchronizeDiscountData координирует процесс синхронизации скидок для объявления
//
// Этапы работы:
// 1. Получение данных (listing, price history)
// 2. Детектирование манипуляций с ценами
// 3. Расчёт скидки из истории цен (если манипуляций нет)
// 4. Парсинг скидки из description (если нет в истории)
// 5. Применение скидки к БД
// 6. Переиндексация в OpenSearch
//
// Функция следует принципу Single Responsibility - координирует процесс,
// делегируя конкретные задачи специализированным функциям
func (s *MarketplaceService) SynchronizeDiscountData(ctx context.Context, listingID int) error {
	logger.Debug().
		Int("listing_id", listingID).
		Msg("Starting discount synchronization")

	// ========================================
	// ЭТАП 1: Получение данных
	// ========================================
	listing, err := s.storage.GetListingByID(ctx, listingID)
	if err != nil {
		logger.Error().
			Err(err).
			Int("listing_id", listingID).
			Msg("Failed to get listing")
		return fmt.Errorf("ошибка получения объявления: %w", err)
	}

	priceHistory, err := s.storage.GetPriceHistory(ctx, listingID)
	if err != nil {
		logger.Warn().
			Err(err).
			Int("listing_id", listingID).
			Msg("Failed to get price history, continuing without it")
		priceHistory = []models.PriceHistoryEntry{}
	}

	// ========================================
	// ЭТАП 2: Детектирование манипуляций
	// ========================================
	if len(priceHistory) > 1 {
		manipulation := s.DetectPriceManipulation(ctx, listing, priceHistory)

		if manipulation.IsManipulated {
			// Удаляем скидку и завершаем
			if err := s.RemoveDiscountDueToManipulation(ctx, listing); err != nil {
				return err
			}

			// Переиндексируем без скидки и возвращаемся
			logger.Info().
				Int("listing_id", listingID).
				Msg("Discount removed due to manipulation, reindexing")
			return s.storage.IndexListing(ctx, listing)
		}
	}

	// ========================================
	// ЭТАП 3: Расчёт скидки из истории цен
	// ========================================
	var discountApplied bool

	if len(priceHistory) > 0 {
		calculatedDiscount := s.CalculateDiscountFromHistory(ctx, listing, priceHistory)

		if err := s.ApplyCalculatedDiscount(ctx, listing, calculatedDiscount); err != nil {
			logger.Error().
				Err(err).
				Int("listing_id", listingID).
				Msg("Failed to apply calculated discount")
			return err
		}

		// Если скидка была успешно рассчитана и применена - помечаем флаг
		if calculatedDiscount.IsValid && IsDiscountSignificant(calculatedDiscount.DiscountPercent) {
			discountApplied = true
			logger.Info().
				Int("listing_id", listingID).
				Int("discount_percent", calculatedDiscount.DiscountPercent).
				Msg("Applied discount from price history")
		}
	}

	// ========================================
	// ЭТАП 4: Парсинг скидки из description
	// ========================================
	// Парсим только если не было скидки из истории
	if !discountApplied {
		// Проверяем что метаданные не содержат скидку
		hasExistingDiscount := listing.Metadata != nil && listing.Metadata["discount"] != nil

		if !hasExistingDiscount {
			parsedDiscount := s.ParseDiscountFromDescription(ctx, listing)

			if parsedDiscount.IsValid {
				if err := s.ApplyParsedDiscount(ctx, listing, parsedDiscount); err != nil {
					logger.Error().
						Err(err).
						Int("listing_id", listingID).
						Msg("Failed to apply parsed discount")
					return err
				}

				logger.Info().
					Int("listing_id", listingID).
					Int("discount_percent", parsedDiscount.DiscountPercent).
					Msg("Applied discount from description")
			}
		}
	}

	// ========================================
	// ЭТАП 5: Переиндексация
	// ========================================
	if err := s.storage.IndexListing(ctx, listing); err != nil {
		logger.Error().
			Err(err).
			Int("listing_id", listingID).
			Msg("Failed to reindex listing")
		return err
	}

	logger.Debug().
		Int("listing_id", listingID).
		Msg("Discount synchronization completed successfully")

	return nil
}
