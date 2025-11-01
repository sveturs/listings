// backend/internal/proj/c2c/service/marketplace_discount_parser.go
package service

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"

	"backend/internal/domain/models"
	"backend/internal/logger"
)

// ParsedDiscount содержит информацию о скидке, распарсенной из описания
type ParsedDiscount struct {
	DiscountPercent int
	OldPrice        float64
	IsValid         bool
	ValidationError string
}

// ParseDiscountFromDescription парсит информацию о скидке из текста описания
// Ищет паттерны вида "30% СКИДКА" и "Старая цена: 150 RSD"
func (s *MarketplaceService) ParseDiscountFromDescription(
	ctx context.Context,
	listing *models.MarketplaceListing,
) *ParsedDiscount {
	// Проверяем наличие ключевых слов о скидке
	if !strings.Contains(listing.Description, "СКИДКА") &&
		!strings.Contains(listing.Description, "СКИДКА!") {
		return &ParsedDiscount{
			IsValid: false,
		}
	}

	// Regex для поиска процента скидки: "30% СКИДКА"
	discountRegex := regexp.MustCompile(`(\d+)%\s*СКИДКА`)
	discountMatches := discountRegex.FindStringSubmatch(listing.Description)

	// Regex для поиска старой цены: "Старая цена: 150 RSD" или "Старая цена: 150.50 RSD"
	priceRegex := regexp.MustCompile(`Старая цена:\s*(\d+[\.,]?\d*)\s*RSD`)
	priceMatches := priceRegex.FindStringSubmatch(listing.Description)

	// Если не нашли оба значения - скидка невалидна
	if len(discountMatches) < 2 || len(priceMatches) < 2 {
		return &ParsedDiscount{
			IsValid:         false,
			ValidationError: "Не найдена полная информация о скидке в описании",
		}
	}

	// Парсим процент скидки
	discountPercent, err := strconv.Atoi(discountMatches[1])
	if err != nil {
		logger.Warn().
			Err(err).
			Int("listing_id", listing.ID).
			Str("discount_str", discountMatches[1]).
			Msg("Failed to parse discount percent")
		return &ParsedDiscount{
			IsValid:         false,
			ValidationError: "Неверный формат процента скидки",
		}
	}

	// Парсим старую цену (заменяем запятую на точку для float парсинга)
	oldPriceStr := strings.ReplaceAll(priceMatches[1], ",", ".")
	oldPrice, err := strconv.ParseFloat(oldPriceStr, 64)
	if err != nil {
		logger.Warn().
			Err(err).
			Int("listing_id", listing.ID).
			Str("price_str", oldPriceStr).
			Msg("Failed to parse old price")
		return &ParsedDiscount{
			IsValid:         false,
			ValidationError: "Неверный формат старой цены",
		}
	}

	// Валидируем реальность скидки
	// Рассчитываем фактический процент скидки
	calculatedDiscount := int((oldPrice - listing.Price) / oldPrice * 100)

	// Проверяем что заявленная скидка примерно соответствует реальной (±5% tolerance)
	discountDifference := abs(calculatedDiscount - discountPercent)
	if calculatedDiscount < 0 || discountDifference > 5 {
		logger.Warn().
			Int("listing_id", listing.ID).
			Int("claimed_discount", discountPercent).
			Int("calculated_discount", calculatedDiscount).
			Float64("old_price", oldPrice).
			Float64("current_price", listing.Price).
			Msg("Discount validation failed: mismatch between claimed and calculated")

		return &ParsedDiscount{
			IsValid:         false,
			ValidationError: "Заявленная скидка не соответствует реальной",
		}
	}

	logger.Info().
		Int("listing_id", listing.ID).
		Int("discount_percent", discountPercent).
		Float64("old_price", oldPrice).
		Float64("current_price", listing.Price).
		Msg("Successfully parsed discount from description")

	return &ParsedDiscount{
		DiscountPercent: discountPercent,
		OldPrice:        oldPrice,
		IsValid:         true,
	}
}

// CreatePriceHistoryFromParsedDiscount создаёт записи в истории цен
// на основе распарсенной информации о скидке
func (s *MarketplaceService) CreatePriceHistoryFromParsedDiscount(
	ctx context.Context,
	listing *models.MarketplaceListing,
	discount *ParsedDiscount,
) error {
	// 1. Закрываем все предыдущие открытые записи истории цен
	if err := s.storage.ClosePriceHistoryEntry(ctx, listing.ID); err != nil {
		logger.Error().
			Err(err).
			Int("listing_id", listing.ID).
			Msg("Failed to close previous price history entries")
		// Продолжаем выполнение, т.к. это не критично
	}

	// 2. Создаём запись со старой ценой, датированную неделю назад
	// (предполагаем что старая цена действовала как минимум неделю)
	effectiveFrom := time.Now().AddDate(0, 0, -7)

	oldPriceEntry := &models.PriceHistoryEntry{
		ListingID:     listing.ID,
		Price:         discount.OldPrice,
		EffectiveFrom: effectiveFrom,
		ChangeSource:  "parsed_from_description",
	}

	if err := s.storage.AddPriceHistoryEntry(ctx, oldPriceEntry); err != nil {
		logger.Error().
			Err(err).
			Int("listing_id", listing.ID).
			Float64("price", discount.OldPrice).
			Msg("Failed to add old price to history")
		return err
	}

	// 3. Создаём запись с текущей (сниженной) ценой
	currentTime := time.Now()
	newPriceEntry := &models.PriceHistoryEntry{
		ListingID:     listing.ID,
		Price:         listing.Price,
		EffectiveFrom: currentTime,
		ChangeSource:  "parsed_from_description",
	}

	if err := s.storage.AddPriceHistoryEntry(ctx, newPriceEntry); err != nil {
		logger.Error().
			Err(err).
			Int("listing_id", listing.ID).
			Float64("price", listing.Price).
			Msg("Failed to add new price to history")
		return err
	}

	logger.Info().
		Int("listing_id", listing.ID).
		Float64("old_price", discount.OldPrice).
		Float64("new_price", listing.Price).
		Msg("Created price history from parsed discount")

	return nil
}
