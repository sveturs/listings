package calculator

import (
	"time"

	"backend/internal/proj/delivery/models"
)

// CalculationResponse - ответ с расчетом стоимости для фронтенда
type CalculationResponse struct {
	Success bool             `json:"success"`
	Data    *CalculationData `json:"data,omitempty"`
	Message string           `json:"message,omitempty"`
}

// CalculationData - данные расчета
type CalculationData struct {
	Providers   []ProviderQuote `json:"providers"`
	Cheapest    *ProviderQuote  `json:"cheapest,omitempty"`
	Fastest     *ProviderQuote  `json:"fastest,omitempty"`
	Recommended *ProviderQuote  `json:"recommended,omitempty"`
}

// ProviderQuote - предложение от провайдера для фронтенда
type ProviderQuote struct {
	ProviderID            int                      `json:"provider_id"`
	ProviderCode          string                   `json:"provider_code"`
	ProviderName          string                   `json:"provider_name"`
	DeliveryType          string                   `json:"delivery_type"`
	TotalPrice            float64                  `json:"total_price"`
	DeliveryCost          float64                  `json:"delivery_cost"`
	InsuranceCost         float64                  `json:"insurance_cost,omitempty"`
	CODFee                float64                  `json:"cod_fee,omitempty"`
	CostBreakdown         models.CostBreakdown     `json:"cost_breakdown"`
	EstimatedDays         int                      `json:"estimated_days"`
	EstimatedDeliveryDate *time.Time               `json:"estimated_delivery_date,omitempty"`
	Services              []models.DeliveryService `json:"services,omitempty"`
	IsAvailable           bool                     `json:"is_available"`
	UnavailableReason     string                   `json:"unavailable_reason,omitempty"`
}
