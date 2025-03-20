// backend/internal/proj/marketplace/service/price_history.go
package service

import (
	"backend/internal/domain/models"
	"backend/internal/storage"
	"context"
	"log"
	"math"
	"time"
)

// PriceHistoryServiceInterface определяет интерфейс для работы с историей цен
type PriceHistoryServiceInterface interface {
	GetPriceHistory(ctx context.Context, listingID int) ([]models.PriceHistoryEntry, error)
	AnalyzeDiscount(ctx context.Context, listingID int) (*models.DiscountInfo, error)
	RecordPriceChange(ctx context.Context, listingID int, oldPrice, newPrice float64, source string) error
}

// PriceHistoryService реализует интерфейс PriceHistoryServiceInterface
type PriceHistoryService struct {
	storage storage.Storage
}

// NewPriceHistoryService создает новый сервис истории цен
func NewPriceHistoryService(storage storage.Storage) PriceHistoryServiceInterface {
	return &PriceHistoryService{
		storage: storage,
	}
}

 // RecordPriceChange записывает изменение цены в историю
 func (s *PriceHistoryService) RecordPriceChange(ctx context.Context, listingID int, oldPrice, newPrice float64, source string) error {
    // Если цена не изменилась, ничего не делаем
    if oldPrice == newPrice {
        return nil
    }

    // Вычисляем процент изменения цены
    var changePercentage float64
    if oldPrice > 0 {
        changePercentage = ((newPrice - oldPrice) / oldPrice) * 100
    }

    // Закрываем текущую активную запись
    if err := s.storage.ClosePriceHistoryEntry(ctx, listingID); err != nil {
        log.Printf("Ошибка при закрытии активной записи истории цен: %v", err)
        // Продолжаем выполнение, не прерываем операцию
    }

    // Создаем новую запись
    entry := &models.PriceHistoryEntry{
        ListingID:        listingID,
        Price:            newPrice,
        EffectiveFrom:    time.Now().UTC(),
        ChangeSource:     source,
        ChangePercentage: changePercentage,
    }

    // Сохраняем запись в базу
    return s.storage.AddPriceHistoryEntry(ctx, entry)
}
// GetPriceHistory возвращает историю изменения цен объявления
func (s *PriceHistoryService) GetPriceHistory(ctx context.Context, listingID int) ([]models.PriceHistoryEntry, error) {
	// Сначала мы должны расширить интерфейс Storage
	// Временно возвращаем пустой массив
	return []models.PriceHistoryEntry{}, nil
}

// AnalyzeDiscount анализирует историю цен и определяет, является ли текущая скидка настоящей
// Возвращает информацию о скидке, если она обнаружена
func (s *PriceHistoryService) AnalyzeDiscount(ctx context.Context, listingID int) (*models.DiscountInfo, error) {
	// Получаем историю цен для данного объявления
	history, err := s.GetPriceHistory(ctx, listingID)
	if err != nil {
		return nil, err
	}

	// Если истории нет или она слишком короткая, скидки нет
	if len(history) < 2 {
		return nil, nil
	}

	// Получаем текущую цену (последняя активная запись)
	var currentPrice float64
	var currentEffectiveFrom time.Time
	var previousPrice float64
	var previousEffectiveFrom time.Time
	var maxPrice float64
	var maxPriceEffectiveFrom time.Time

	// Инициализируем значения
	currentPrice = history[0].Price
	currentEffectiveFrom = history[0].EffectiveFrom
	maxPrice = currentPrice
	maxPriceEffectiveFrom = currentEffectiveFrom

	// Проходим по истории, ищем предыдущую и максимальную цены
	for i, entry := range history {
		// Пропускаем текущую запись
		if i == 0 {
			continue
		}

		// Находим предыдущую цену, если еще не нашли
		if previousPrice == 0 {
			previousPrice = entry.Price
			previousEffectiveFrom = entry.EffectiveFrom
		}

		// Обновляем максимальную цену, если нашли большую
		if entry.Price > maxPrice {
			maxPrice = entry.Price
			maxPriceEffectiveFrom = entry.EffectiveFrom
		}
	}

	// Временно пропускаем проверку манипуляций с ценой, пока не расширим интерфейс storage
	isSuspicious := false

	// Если текущая цена ниже предыдущей, это может быть скидка
	if currentPrice < previousPrice {
		// Вычисляем процент скидки от предыдущей цены
		discountPercent := math.Round((1 - currentPrice/previousPrice) * 100)

		// Вычисляем процент скидки от максимальной цены
		maxDiscountPercent := math.Round((1 - currentPrice/maxPrice) * 100)

		// Определяем, насколько давно была установлена предыдущая цена
		daysSincePrevious := time.Since(previousEffectiveFrom).Hours() / 24
		
		// Проверяем, является ли скидка реальной:
		// 1. Скидка должна быть значительной (более 5%)
		// 2. Предыдущая цена должна быть установлена не менее 7 дней назад
		// 3. Скидка не должна быть подозрительной
		isRealDiscount := discountPercent >= 5 && daysSincePrevious >= 7 && !isSuspicious

		// Создаем объект с информацией о скидке
		discountInfo := &models.DiscountInfo{
			CurrentPrice:       currentPrice,
			PreviousPrice:      previousPrice,
			MaxPrice:           maxPrice,
			DiscountPercent:    int(discountPercent),
			MaxDiscountPercent: int(maxDiscountPercent),
			EffectiveFrom:      currentEffectiveFrom,
			PreviousEffectiveFrom: previousEffectiveFrom,
			MaxPriceEffectiveFrom: maxPriceEffectiveFrom,
			IsRealDiscount:     isRealDiscount,
			IsSuspicious:       isSuspicious,
		}

		return discountInfo, nil
	}

	// Если текущая цена не ниже предыдущей, скидки нет
	return nil, nil
}
