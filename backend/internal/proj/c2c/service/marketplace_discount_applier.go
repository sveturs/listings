// backend/internal/proj/c2c/service/marketplace_discount_applier.go
package service

import (
	"context"

	"backend/internal/domain/models"
	"backend/internal/logger"
)

// ApplyDiscountAction определяет тип действия с скидкой
type ApplyDiscountAction string

const (
	ApplyDiscountActionSet    ApplyDiscountAction = "set"    // Установить/обновить скидку
	ApplyDiscountActionRemove ApplyDiscountAction = "remove" // Удалить скидку
)

// ApplyDiscountMetadata применяет изменения метаданных о скидке к объявлению в БД
func (s *MarketplaceService) ApplyDiscountMetadata(
	ctx context.Context,
	listing *models.MarketplaceListing,
	action ApplyDiscountAction,
	discountData map[string]interface{},
) error {
	// Инициализируем metadata если его нет
	if listing.Metadata == nil {
		listing.Metadata = make(map[string]interface{})
	}

	// Выполняем действие
	switch action {
	case ApplyDiscountActionSet:
		// Устанавливаем/обновляем скидку
		listing.Metadata["discount"] = discountData

		logger.Info().
			Int("listing_id", listing.ID).
			Interface("discount", discountData).
			Msg("Setting discount metadata")

	case ApplyDiscountActionRemove:
		// Удаляем скидку
		if listing.Metadata["discount"] != nil {
			delete(listing.Metadata, "discount")

			logger.Info().
				Int("listing_id", listing.ID).
				Msg("Removing discount metadata")
		}
	}

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
			Str("action", string(action)).
			Msg("Failed to update discount metadata in database")
		return err
	}

	return nil
}

// ApplyCalculatedDiscount применяет рассчитанную скидку к объявлению
// Проверяет что скидка значимая (>=5%) и обновляет БД
func (s *MarketplaceService) ApplyCalculatedDiscount(
	ctx context.Context,
	listing *models.MarketplaceListing,
	discount *CalculatedDiscount,
) error {
	// Если скидка невалидна - удаляем информацию о скидке
	if !discount.IsValid {
		return s.ApplyDiscountMetadata(ctx, listing, ApplyDiscountActionRemove, nil)
	}

	// Проверяем что скидка достаточно большая
	if !IsDiscountSignificant(discount.DiscountPercent) {
		logger.Debug().
			Int("listing_id", listing.ID).
			Int("discount_percent", discount.DiscountPercent).
			Msg("Discount is too small, removing")

		return s.ApplyDiscountMetadata(ctx, listing, ApplyDiscountActionRemove, nil)
	}

	// Создаём метаданные и применяем
	metadata := CreateDiscountMetadata(discount)
	return s.ApplyDiscountMetadata(ctx, listing, ApplyDiscountActionSet, metadata)
}

// ApplyParsedDiscount применяет распарсенную скидку к объявлению
func (s *MarketplaceService) ApplyParsedDiscount(
	ctx context.Context,
	listing *models.MarketplaceListing,
	discount *ParsedDiscount,
) error {
	// Если скидка невалидна - ничего не делаем
	if !discount.IsValid {
		logger.Debug().
			Int("listing_id", listing.ID).
			Str("reason", discount.ValidationError).
			Msg("Parsed discount is invalid, skipping")
		return nil
	}

	// Создаём историю цен
	if err := s.CreatePriceHistoryFromParsedDiscount(ctx, listing, discount); err != nil {
		return err
	}

	// Создаём метаданные и применяем
	// Используем время неделю назад как effective_from (как в оригинале)
	effectiveFrom := listing.CreatedAt.AddDate(0, 0, -7)
	metadata := CreateDiscountMetadataFromParsed(discount, effectiveFrom)

	return s.ApplyDiscountMetadata(ctx, listing, ApplyDiscountActionSet, metadata)
}
