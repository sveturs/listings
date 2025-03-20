// backend/internal/domain/models/price_history.go
package models

import (
	"time"
)

// PriceHistoryEntry представляет запись в истории изменения цен
type PriceHistoryEntry struct {
	ID               int       `json:"id"`
	ListingID        int       `json:"listing_id"`
	Price            float64   `json:"price"`
	EffectiveFrom    time.Time `json:"effective_from"`
	EffectiveTo      *time.Time `json:"effective_to,omitempty"`
	ChangeSource     string    `json:"change_source"` // manual, import, system, etc.
	ChangePercentage float64   `json:"change_percentage,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

// DiscountInfo содержит информацию о скидке для объявления
type DiscountInfo struct {
	CurrentPrice         float64   `json:"current_price"`
	PreviousPrice        float64   `json:"previous_price"`
	MaxPrice             float64   `json:"max_price"`
	DiscountPercent      int       `json:"discount_percent"`
	MaxDiscountPercent   int       `json:"max_discount_percent"`
	EffectiveFrom        time.Time `json:"effective_from"`
	PreviousEffectiveFrom time.Time `json:"previous_effective_from"`
	MaxPriceEffectiveFrom time.Time `json:"max_price_effective_from"`
	IsRealDiscount       bool      `json:"is_real_discount"`
	IsSuspicious         bool      `json:"is_suspicious"`
}

// DiscountData содержит упрощенную информацию о скидке для отображения на фронтенде
type DiscountData struct {
	DiscountPercent    int       `json:"discount_percent"`    // Процент скидки
	PreviousPrice      float64   `json:"previous_price"`      // Предыдущая цена
	EffectiveFrom      time.Time `json:"effective_from"`      // Когда скидка начала действовать
	HasPriceHistory    bool      `json:"has_price_history"`   // Есть ли история цен для отображения графика
}