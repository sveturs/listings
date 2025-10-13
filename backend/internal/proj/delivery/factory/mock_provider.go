package factory

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	"backend/internal/config"
	"backend/internal/proj/delivery/interfaces"
)

// cryptoRandIntn возвращает случайное число от 0 до n-1 используя crypto/rand
func cryptoRandIntn(n int) int {
	if n <= 0 {
		return 0
	}
	var bytes [4]byte
	_, err := rand.Read(bytes[:])
	if err != nil {
		// Fallback в случае ошибки - возвращаем фиксированное значение
		return 0
	}
	// Используем 32-битное значение для безопасности
	randomValue := binary.BigEndian.Uint32(bytes[:])
	return int(randomValue) % n
}

// MockProvider - мок-провайдер для тестирования
type MockProvider struct {
	code         string
	name         string
	isActive     bool
	capabilities *interfaces.ProviderCapabilities
}

// NewMockProvider создает новый мок-провайдер
func NewMockProvider(code, name string) *MockProvider {
	return &MockProvider{
		code:     code,
		name:     name,
		isActive: true,
		capabilities: &interfaces.ProviderCapabilities{
			MaxWeightKg:       50,
			MaxVolumeM3:       1.0,
			MaxLengthCm:       200,
			DeliveryZones:     []string{"serbia", "region"},
			DeliveryTypes:     []string{"standard", "express"},
			SupportsCOD:       true,
			SupportsInsurance: true,
			SupportsTracking:  true,
			SupportsPickup:    true,
			SupportsReturn:    false,
			Services:          []string{"sms_notification", "email_notification"},
		},
	}
}

// GetCode возвращает код провайдера
func (m *MockProvider) GetCode() string {
	return m.code
}

// GetName возвращает название провайдера
func (m *MockProvider) GetName() string {
	return m.name
}

// IsActive проверяет, активен ли провайдер
func (m *MockProvider) IsActive() bool {
	return m.isActive
}

// CalculateRate рассчитывает стоимость доставки
func (m *MockProvider) CalculateRate(ctx context.Context, req *interfaces.RateRequest) (*interfaces.RateResponse, error) {
	// Симулируем задержку API
	time.Sleep(time.Millisecond * time.Duration(cryptoRandIntn(500)))

	// Рассчитываем общий вес
	totalWeight := 0.0
	totalVolume := 0.0
	hasFragile := false

	for _, pkg := range req.Packages {
		totalWeight += pkg.Weight
		if pkg.Dimensions != nil {
			volume := (pkg.Dimensions.Length * pkg.Dimensions.Width * pkg.Dimensions.Height) / 1000000
			totalVolume += volume
		}
		if pkg.IsFragile {
			hasFragile = true
		}
	}

	// Базовая стоимость
	basePrice := 300.0
	weightPrice := totalWeight * 50    // 50 RSD за кг
	volumePrice := totalVolume * 10000 // 10000 RSD за м³

	// Берем максимум из весовой и объемной стоимости
	shippingPrice := basePrice + max(weightPrice, volumePrice)

	// Доплаты
	if hasFragile {
		shippingPrice += 100
	}

	insuranceFee := 0.0
	if req.InsuranceValue > 0 {
		insuranceFee = req.InsuranceValue * 0.01
	}

	codFee := 0.0
	if req.CODAmount > 0 {
		codFee = max(50, req.CODAmount*0.02)
	}

	// Создаем опции доставки
	options := []interfaces.DeliveryOption{}

	// Стандартная доставка
	standardCost := shippingPrice + insuranceFee + codFee
	options = append(options, interfaces.DeliveryOption{
		Type:          "standard",
		Name:          "Стандартная доставка",
		TotalCost:     standardCost,
		EstimatedDays: 3 + cryptoRandIntn(3),
		CostBreakdown: &interfaces.CostBreakdown{
			BasePrice:       basePrice,
			WeightSurcharge: weightPrice,
			FragileSurcharge: func() float64 {
				if hasFragile {
					return 100
				}
				return 0
			}(),
			InsuranceFee: insuranceFee,
			CODFee:       codFee,
		},
	})

	// Экспресс доставка
	expressCost := (shippingPrice * 1.5) + insuranceFee + codFee
	options = append(options, interfaces.DeliveryOption{
		Type:          "express",
		Name:          "Экспресс доставка",
		TotalCost:     expressCost,
		EstimatedDays: 1 + cryptoRandIntn(2),
		CostBreakdown: &interfaces.CostBreakdown{
			BasePrice:       basePrice * 1.5,
			WeightSurcharge: weightPrice * 1.5,
			FragileSurcharge: func() float64 {
				if hasFragile {
					return 150
				}
				return 0
			}(),
			InsuranceFee: insuranceFee,
			CODFee:       codFee,
		},
	})

	return &interfaces.RateResponse{
		ProviderCode:    m.code,
		ProviderName:    m.name,
		DeliveryOptions: options,
		Currency:        config.GetGlobalDefaultCurrency(),
		ValidUntil:      time.Now().Add(24 * time.Hour),
	}, nil
}

// CreateShipment создает отправление
func (m *MockProvider) CreateShipment(ctx context.Context, req *interfaces.ShipmentRequest) (*interfaces.ShipmentResponse, error) {
	// Симулируем задержку API
	time.Sleep(time.Millisecond * time.Duration(cryptoRandIntn(1000)))

	// Генерируем случайный tracking number
	trackingNumber := fmt.Sprintf("%s-%d-%d", m.code, time.Now().Unix(), cryptoRandIntn(10000))

	// Рассчитываем стоимость
	rateReq := &interfaces.RateRequest{
		FromAddress:    req.FromAddress,
		ToAddress:      req.ToAddress,
		Packages:       req.Packages,
		DeliveryType:   req.DeliveryType,
		InsuranceValue: req.InsuranceValue,
		CODAmount:      req.CODAmount,
	}

	rateResp, err := m.CalculateRate(ctx, rateReq)
	if err != nil {
		return nil, err
	}

	// Находим выбранную опцию доставки
	var selectedOption *interfaces.DeliveryOption
	for _, option := range rateResp.DeliveryOptions {
		if option.Type == req.DeliveryType {
			selectedOption = &option
			break
		}
	}

	if selectedOption == nil && len(rateResp.DeliveryOptions) > 0 {
		selectedOption = &rateResp.DeliveryOptions[0]
	}

	estimatedDate := time.Now().AddDate(0, 0, selectedOption.EstimatedDays)

	return &interfaces.ShipmentResponse{
		ShipmentID:     fmt.Sprintf("%d", cryptoRandIntn(1000000)),
		TrackingNumber: trackingNumber,
		ExternalID:     fmt.Sprintf("EXT-%s-%d", m.code, cryptoRandIntn(100000)),
		Status:         interfaces.StatusConfirmed,
		TotalCost:      selectedOption.TotalCost,
		CostBreakdown:  selectedOption.CostBreakdown,
		EstimatedDate:  &estimatedDate,
		Labels: []interfaces.LabelInfo{
			{
				Type:   "shipping",
				Format: "pdf",
				URL:    fmt.Sprintf("https://mock-labels.example.com/%s.pdf", trackingNumber),
			},
		},
		CreatedAt: time.Now(),
	}, nil
}

// TrackShipment отслеживает отправление
func (m *MockProvider) TrackShipment(ctx context.Context, trackingNumber string) (*interfaces.TrackingResponse, error) {
	// Симулируем задержку API
	time.Sleep(time.Millisecond * time.Duration(cryptoRandIntn(500)))

	// Генерируем случайную историю событий
	events := []interfaces.TrackingEvent{
		{
			Timestamp:   time.Now().Add(-72 * time.Hour),
			Status:      interfaces.StatusPending,
			Description: "Заказ создан",
			Location:    "Белград",
		},
		{
			Timestamp:   time.Now().Add(-48 * time.Hour),
			Status:      interfaces.StatusConfirmed,
			Description: "Заказ подтвержден",
			Location:    "Белград",
		},
		{
			Timestamp:   time.Now().Add(-36 * time.Hour),
			Status:      interfaces.StatusPickedUp,
			Description: "Посылка забрана курьером",
			Location:    "Белград, Склад",
		},
		{
			Timestamp:   time.Now().Add(-24 * time.Hour),
			Status:      interfaces.StatusInTransit,
			Description: "Посылка в пути",
			Location:    "Нови-Сад",
		},
	}

	// Случайный текущий статус
	statuses := []string{
		interfaces.StatusInTransit,
		interfaces.StatusOutForDelivery,
		interfaces.StatusDelivered,
	}
	currentStatus := statuses[cryptoRandIntn(len(statuses))]

	response := &interfaces.TrackingResponse{
		TrackingNumber:  trackingNumber,
		Status:          currentStatus,
		StatusText:      getStatusDescription(currentStatus),
		CurrentLocation: "Нови-Сад",
		Events:          events,
	}

	// Если доставлено, добавляем информацию
	if currentStatus == interfaces.StatusDelivered {
		deliveredTime := time.Now().Add(-2 * time.Hour)
		response.DeliveredDate = &deliveredTime
		response.ProofOfDelivery = &interfaces.ProofOfDelivery{
			RecipientName: "Иван Иванович",
			DeliveredAt:   deliveredTime,
			SignatureURL:  fmt.Sprintf("https://mock-signatures.example.com/%s.png", trackingNumber),
		}

		// Добавляем событие доставки
		response.Events = append(response.Events, interfaces.TrackingEvent{
			Timestamp:   deliveredTime,
			Status:      interfaces.StatusDelivered,
			Description: "Посылка доставлена получателю",
			Location:    "Нови-Сад",
			Details:     "Подписано: Иван Иванович",
		})
	} else {
		estimatedDate := time.Now().Add(24 * time.Hour)
		response.EstimatedDate = &estimatedDate
	}

	return response, nil
}

// CancelShipment отменяет отправление
func (m *MockProvider) CancelShipment(ctx context.Context, shipmentID string) error {
	// Симулируем задержку API
	time.Sleep(time.Millisecond * time.Duration(cryptoRandIntn(500)))

	// Случайно возвращаем ошибку (10% вероятность)
	if cryptoRandIntn(10) == 0 {
		return fmt.Errorf("shipment %s cannot be canceled: already in transit", shipmentID)
	}

	return nil
}

// GetLabel получает этикетку
func (m *MockProvider) GetLabel(ctx context.Context, shipmentID string) (*interfaces.LabelResponse, error) {
	// Симулируем задержку API
	time.Sleep(time.Millisecond * time.Duration(cryptoRandIntn(500)))

	return &interfaces.LabelResponse{
		Labels: []interfaces.LabelInfo{
			{
				Type:   "shipping",
				Format: "pdf",
				URL:    fmt.Sprintf("https://mock-labels.example.com/shipment-%s.pdf", shipmentID),
				Data:   []byte("Mock PDF content"),
			},
			{
				Type:   "return",
				Format: "pdf",
				URL:    fmt.Sprintf("https://mock-labels.example.com/return-%s.pdf", shipmentID),
			},
		},
	}, nil
}

// ValidateAddress проверяет адрес
func (m *MockProvider) ValidateAddress(ctx context.Context, address *interfaces.Address) (*interfaces.AddressValidationResponse, error) {
	// Симулируем задержку API
	time.Sleep(time.Millisecond * time.Duration(cryptoRandIntn(200)))

	// Список поддерживаемых городов
	supportedCities := map[string]bool{
		"Белград":    true,
		"Belgrade":   true,
		"Нови Сад":   true,
		"Novi Sad":   true,
		"Ниш":        true,
		"Nis":        true,
		"Крагујевац": true,
		"Kragujevac": true,
		"Суботица":   true,
		"Subotica":   true,
	}

	isValid := supportedCities[address.City]

	response := &interfaces.AddressValidationResponse{
		IsValid:           isValid,
		DeliveryAvailable: isValid,
	}

	if !isValid {
		response.ValidationErrors = []string{
			fmt.Sprintf("City '%s' is not in our delivery zone", address.City),
		}

		// Предлагаем ближайший город
		response.SuggestedAddress = &interfaces.Address{
			Name:       address.Name,
			Phone:      address.Phone,
			Email:      address.Email,
			Street:     address.Street,
			City:       "Белград", // По умолчанию предлагаем Белград
			PostalCode: address.PostalCode,
			Country:    address.Country,
		}
	}

	// Определяем зону
	switch {
	case address.City == "Белград" || address.City == "Belgrade":
		response.Zone = "local"
	case isValid:
		response.Zone = "national"
	default:
		response.Zone = "unavailable"
	}

	return response, nil
}

// GetCapabilities возвращает возможности провайдера
func (m *MockProvider) GetCapabilities() *interfaces.ProviderCapabilities {
	return m.capabilities
}

// Вспомогательные функции

func getStatusDescription(status string) string {
	descriptions := map[string]string{
		interfaces.StatusPending:           "Ожидает обработки",
		interfaces.StatusConfirmed:         "Заказ подтвержден",
		interfaces.StatusPickedUp:          "Забрано курьером",
		interfaces.StatusInTransit:         "В пути",
		interfaces.StatusOutForDelivery:    "Передано курьеру для доставки",
		interfaces.StatusDelivered:         "Доставлено",
		interfaces.StatusDeliveryAttempted: "Попытка доставки не удалась",
		interfaces.StatusReturning:         "Возвращается отправителю",
		interfaces.StatusReturned:          "Возвращено отправителю",
		interfaces.StatusCancelled:         "Отменено",
		interfaces.StatusLost:              "Утеряно",
		interfaces.StatusDamaged:           "Повреждено при транспортировке",
	}

	if desc, ok := descriptions[status]; ok {
		return desc
	}
	return status
}

// HandleWebhook обрабатывает webhook от мок-провайдера
func (m *MockProvider) HandleWebhook(ctx context.Context, payload []byte, headers map[string]string) (*interfaces.WebhookResponse, error) {
	// Для мок-провайдера создаем симуляцию webhook обработки
	response := &interfaces.WebhookResponse{
		Processed: true,
		Timestamp: time.Now(),
	}

	// Попробуем распарсить payload как JSON для извлечения трек-номера и статуса
	var webhookData map[string]interface{}
	if err := json.Unmarshal(payload, &webhookData); err != nil {
		// Если не удается распарсить, используем заглушку
		response.TrackingNumber = "MOCK_" + fmt.Sprintf("%d", time.Now().Unix())
		response.Status = interfaces.StatusInTransit
		response.StatusDetails = "Mock status update"
	} else {
		// Извлекаем данные из webhook payload
		if trackingNum, ok := webhookData["tracking_number"].(string); ok {
			response.TrackingNumber = trackingNum
		}

		if status, ok := webhookData["status"].(string); ok {
			response.Status = status
		} else {
			response.Status = interfaces.StatusInTransit
		}

		if location, ok := webhookData["location"].(string); ok {
			response.Location = location
		}

		if details, ok := webhookData["status_details"].(string); ok {
			response.StatusDetails = details
		}
	}

	// Создаем симуляцию события отслеживания
	event := interfaces.TrackingEvent{
		Timestamp:   response.Timestamp,
		Status:      response.Status,
		Location:    response.Location,
		Description: response.StatusDetails,
	}

	// Если статус "доставлено", добавляем детали доставки
	if response.Status == interfaces.StatusDelivered {
		response.DeliveryDetails = &interfaces.ProofOfDelivery{
			RecipientName: "Mock Recipient",
			DeliveredAt:   response.Timestamp,
			Notes:         "Delivered by mock provider",
		}
		event.Description = "Package delivered successfully"
	}

	response.Events = []interfaces.TrackingEvent{event}

	return response, nil
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
