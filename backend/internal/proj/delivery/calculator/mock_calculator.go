package calculator

import (
	"context"
	"math"
	"time"

	"backend/internal/proj/delivery/models"
)

// MockCalculator - временный калькулятор для демонстрации
type MockCalculator struct{}

// NewMockCalculator создает новый mock калькулятор
func NewMockCalculator() *MockCalculator {
	return &MockCalculator{}
}

// CalculateMock - простой расчет доставки для демонстрации
func (m *MockCalculator) CalculateMock(ctx context.Context, req *CalculationRequest) (*CalculationResponse, error) {
	providers := []struct {
		id   int
		code string
		name string
		logo string
	}{
		{1, "post-express", "Post Express", "📮"},
		{2, "bex-express", "BEX Express", "📦"},
		{3, "aks-express", "AKS Express", "🚚"},
		{4, "d-express", "D-Express", "🚛"},
		{5, "city-express", "City Express", "🏙️"},
	}

	var quotes []ProviderQuote

	// Рассчитываем общий вес и объем
	totalWeight := 0.0
	totalValue := 0.0
	hasFragile := false

	for _, item := range req.Items {
		if item.Attributes != nil {
			totalWeight += item.Attributes.WeightKg * float64(item.Quantity)
			if item.Attributes.IsFragile {
				hasFragile = true
			}
		} else {
			// Дефолтный вес 1кг на товар
			totalWeight += float64(item.Quantity)
		}
		// Предполагаем цену товара 1000 RSD (для расчета страховки)
		totalValue += 1000 * float64(item.Quantity)
	}

	// Генерируем предложения от каждого провайдера
	for _, provider := range providers {
		// Базовая цена зависит от провайдера
		basePrice := 300.0 + float64(provider.id)*40

		// Надбавки
		weightSurcharge := 0.0
		if totalWeight > 2 {
			weightSurcharge = math.Ceil(totalWeight-2) * 50
		}

		fragileSurcharge := 0.0
		if hasFragile {
			fragileSurcharge = 100
		}

		// Разные типы доставки для каждого провайдера
		deliveryTypes := []struct {
			name          string
			priceModifier float64
			days          int
			services      []models.DeliveryService
		}{
			{
				name:          "standard",
				priceModifier: 1.0,
				days:          3,
				services: []models.DeliveryService{
					{Name: "Отслеживание", Code: "tracking", IsIncluded: true},
					{Name: "SMS уведомления", Code: "sms", IsIncluded: false, Price: 50, IsAvailable: true},
				},
			},
			{
				name:          "express",
				priceModifier: 1.5,
				days:          1,
				services: []models.DeliveryService{
					{Name: "Отслеживание", Code: "tracking", IsIncluded: true},
					{Name: "SMS уведомления", Code: "sms", IsIncluded: true, IsAvailable: true},
					{Name: "Приоритетная обработка", Code: "priority", IsIncluded: true, IsAvailable: true},
				},
			},
		}

		// Добавляем самовывоз для некоторых провайдеров
		if provider.id%2 == 0 {
			deliveryTypes = append(deliveryTypes, struct {
				name          string
				priceModifier float64
				days          int
				services      []models.DeliveryService
			}{
				name:          "pickup",
				priceModifier: 0,
				days:          1,
				services: []models.DeliveryService{
					{Name: "Самовывоз со склада", Code: "self_pickup", IsIncluded: true, IsAvailable: true},
					{Name: "Хранение 7 дней", Code: "storage", IsIncluded: true, IsAvailable: true},
				},
			})
		}

		for _, dt := range deliveryTypes {
			deliveryCost := basePrice * dt.priceModifier
			if dt.name == "pickup" {
				deliveryCost = 0 // Самовывоз бесплатный
				weightSurcharge = 0
				fragileSurcharge = 0
			}

			totalCost := deliveryCost + weightSurcharge + fragileSurcharge

			// Добавляем страховку если запрошена
			insuranceCost := 0.0
			if req.InsuranceValue > 0 {
				insuranceCost = req.InsuranceValue * 0.03
				totalCost += insuranceCost
			}

			// COD комиссия
			codFee := 0.0
			if req.CODAmount > 0 {
				codFee = req.CODAmount * 0.02
				totalCost += codFee
			}

			// Дата доставки
			estimatedDate := time.Now().AddDate(0, 0, dt.days)

			quote := ProviderQuote{
				ProviderID:    provider.id,
				ProviderCode:  provider.code,
				ProviderName:  provider.name,
				DeliveryType:  dt.name,
				TotalPrice:    totalCost,
				DeliveryCost:  deliveryCost,
				InsuranceCost: insuranceCost,
				CODFee:        codFee,
				CostBreakdown: models.CostBreakdown{
					BasePrice:        basePrice * dt.priceModifier,
					WeightSurcharge:  weightSurcharge,
					FragileSurcharge: fragileSurcharge,
				},
				EstimatedDays:         dt.days,
				EstimatedDeliveryDate: &estimatedDate,
				Services:              dt.services,
				IsAvailable:           true,
			}

			quotes = append(quotes, quote)
		}
	}

	// Выбираем лучшие предложения
	var cheapest, fastest, recommended *ProviderQuote

	for i := range quotes {
		quote := &quotes[i]

		// Самый дешевый
		if cheapest == nil || quote.TotalPrice < cheapest.TotalPrice {
			cheapest = quote
		}

		// Самый быстрый
		if fastest == nil || quote.EstimatedDays < fastest.EstimatedDays {
			fastest = quote
		}
	}

	// Рекомендуемый - баланс цены и скорости
	// Выбираем стандартную доставку от Post Express как рекомендуемую
	for i := range quotes {
		quote := &quotes[i]
		if quote.ProviderCode == "post-express" && quote.DeliveryType == "standard" {
			recommended = quote
			break
		}
	}

	// Если не нашли Post Express, берем первый доступный
	if recommended == nil && len(quotes) > 0 {
		recommended = &quotes[0]
	}

	return &CalculationResponse{
		Success: true,
		Data: &CalculationData{
			Providers:   quotes,
			Cheapest:    cheapest,
			Fastest:     fastest,
			Recommended: recommended,
		},
	}, nil
}
