package factory

import (
	"context"
	"fmt"
	"time"

	"backend/internal/config"
	"backend/internal/proj/delivery/interfaces"
	"backend/internal/proj/postexpress"
)

// PostExpressAdapter - адаптер для Post Express
type PostExpressAdapter struct {
	service *postexpress.Service
}

// NewPostExpressAdapter создает новый адаптер для Post Express
func NewPostExpressAdapter(service *postexpress.Service) *PostExpressAdapter {
	return &PostExpressAdapter{
		service: service,
	}
}

// GetCode возвращает код провайдера
func (a *PostExpressAdapter) GetCode() string {
	return "post_express"
}

// GetName возвращает название провайдера
func (a *PostExpressAdapter) GetName() string {
	return "Post Express"
}

// IsActive проверяет, активен ли провайдер
func (a *PostExpressAdapter) IsActive() bool {
	return a.service != nil
}

// GetCapabilities возвращает возможности провайдера
func (a *PostExpressAdapter) GetCapabilities() *interfaces.ProviderCapabilities {
	return &interfaces.ProviderCapabilities{
		MaxWeightKg:       30.0,
		MaxVolumeM3:       0.5,
		MaxLengthCm:       120.0,
		DeliveryZones:     []string{"local", "national"},
		DeliveryTypes:     []string{"standard", "express"},
		SupportsCOD:       true,
		SupportsInsurance: true,
		SupportsTracking:  true,
		SupportsPickup:    true,
		SupportsReturn:    false,
		Services:          []string{"tracking", "insurance", "cod", "sms"},
	}
}

// CalculateRate рассчитывает стоимость доставки
func (a *PostExpressAdapter) CalculateRate(ctx context.Context, req *interfaces.RateRequest) (*interfaces.RateResponse, error) {
	if a.service == nil {
		return nil, fmt.Errorf("post Express service not initialized")
	}

	// Маппинг из универсального запроса в Post Express формат
	peReq := &postexpress.RateRequest{
		FromCity:  req.FromAddress.City,
		ToCity:    req.ToAddress.City,
		Weight:    calculateTotalWeight(req.Packages),
		Value:     req.InsuranceValue,
		CODAmount: req.CODAmount,
		Services:  req.Services,
	}

	peResp, err := a.service.CalculateRate(ctx, peReq)
	if err != nil {
		return nil, fmt.Errorf("post Express rate calculation failed: %w", err)
	}

	// Маппинг из Post Express ответа в универсальный формат
	deliveryOptions := make([]interfaces.DeliveryOption, 0, len(peResp.DeliveryOptions))
	for _, option := range peResp.DeliveryOptions {
		deliveryOptions = append(deliveryOptions, interfaces.DeliveryOption{
			Type:          option.Type,
			Name:          option.Name,
			TotalCost:     option.TotalPrice,
			EstimatedDays: option.EstimatedDays,
			CostBreakdown: &interfaces.CostBreakdown{
				BasePrice:     option.BasePrice,
				CODFee:        option.CODFee,
				InsuranceFee:  option.InsuranceFee,
				FuelSurcharge: option.FuelSurcharge,
			},
		})
	}

	return &interfaces.RateResponse{
		ProviderCode:    "post_express",
		ProviderName:    "Post Express",
		Currency:        config.GetGlobalDefaultCurrency(),
		ValidUntil:      time.Now().Add(24 * time.Hour),
		DeliveryOptions: deliveryOptions,
	}, nil
}

// CreateShipment создает отправление
func (a *PostExpressAdapter) CreateShipment(ctx context.Context, req *interfaces.ShipmentRequest) (*interfaces.ShipmentResponse, error) {
	if a.service == nil {
		return nil, fmt.Errorf("post Express service not initialized")
	}

	// Маппинг в Post Express B2B формат
	config := a.service.GetConfig()

	// Определяем способ оплаты
	nacinPlacanja := "POF" // По умолчанию postanska uplatnica
	if req.CODAmount > 0 {
		nacinPlacanja = "N" // Наличные при получении
	}

	// Конвертируем вес из кг в граммы
	weightGrams := int(calculateTotalWeight(req.Packages) * 1000)

	// Конвертируем денежные суммы в para (1 RSD = 100 para)
	codPara := int(req.CODAmount * 100)
	insurancePara := int(req.InsuranceValue * 100)

	peReq := &postexpress.ShipmentRequest{
		// Обязательные B2B поля
		ExtBrend:          config.Brand,
		ExtMagacin:        "WAREHOUSE1",
		ExtReferenca:      fmt.Sprintf("SVETU-%d", time.Now().Unix()),
		NacinPrijema:      "K", // K - курьер
		ImaPrijemniBrojDN: "N", // "N" = нет приёмного номера, "D" = есть приёмный номер
		NacinPlacanja:     nacinPlacanja,

		// Отправитель внутри отправления
		Posiljalac: postexpress.SenderInfo{
			Naziv: config.Brand,
			Adresa: &postexpress.AddressInfo{
				Ulica:         req.FromAddress.Street,
				Mesto:         req.FromAddress.City,
				PostanskiBroj: req.FromAddress.PostalCode,
				OznakaZemlje:  "RS",
			},
			Mesto:         req.FromAddress.City,
			PostanskiBroj: req.FromAddress.PostalCode,
			Telefon:       req.FromAddress.Phone,
			Email:         req.FromAddress.Email,
			OznakaZemlje:  "RS",
		},

		// Место забора для курьера
		MestoPreuzimanja: &postexpress.SenderInfo{
			Naziv: config.Brand,
			Adresa: &postexpress.AddressInfo{
				Ulica:         req.FromAddress.Street,
				Mesto:         req.FromAddress.City,
				PostanskiBroj: req.FromAddress.PostalCode,
				OznakaZemlje:  "RS",
			},
			Mesto:         req.FromAddress.City,
			PostanskiBroj: req.FromAddress.PostalCode,
			Telefon:       req.FromAddress.Phone,
			Email:         req.FromAddress.Email,
			OznakaZemlje:  "RS",
		},

		// Основные данные
		BrojPosiljke: fmt.Sprintf("SVETU-%d", time.Now().Unix()),
		IDRukovanje:  71, // PE_Danas_za_sutra_isporuka (ID 29 отменена с 01.01.2022)
		Masa:         weightGrams,

		// Получатель
		Primalac: postexpress.ReceiverInfo{
			TipAdrese: "S", // Стандартный адрес
			Naziv:     req.ToAddress.Name,
			Telefon:   req.ToAddress.Phone,
			Email:     req.ToAddress.Email,
			Adresa: &postexpress.AddressInfo{
				Ulica:         req.ToAddress.Street,
				Mesto:         req.ToAddress.City,
				PostanskiBroj: req.ToAddress.PostalCode,
				OznakaZemlje:  "RS",
			},
			Mesto:         req.ToAddress.City,
			PostanskiBroj: req.ToAddress.PostalCode,
			OznakaZemlje:  "RS",
		},

		// COD и ценность
		Otkupnina: codPara,
		Vrednost:  insurancePara,

		// Дополнительные услуги через запятую
		PosebneUsluge: buildPosebneUsluge(req.Services), // PNA + другие услуги
	}

	// Валидация перед отправкой
	if err := a.service.ValidateShipment(peReq); err != nil {
		return nil, fmt.Errorf("shipment validation failed: %w", err)
	}

	// Создание отправления
	peResp, err := a.service.CreateShipment(ctx, peReq)
	if err != nil {
		return nil, fmt.Errorf("post Express shipment creation failed: %w", err)
	}

	// Маппинг ответа
	resp := &interfaces.ShipmentResponse{
		ShipmentID:     fmt.Sprintf("%d", peResp.IDPosiljke),
		TrackingNumber: peResp.TrackingNumber,
		ExternalID:     peResp.BrojPosiljke,
		Status:         mapPostExpressStatus(peResp.Status),
		TotalCost:      0, // Будет заполнено позже из расчета
		CreatedAt:      time.Now(),
	}

	// Добавляем этикетку если есть
	if peResp.LabelURL != "" {
		resp.Labels = []interfaces.LabelInfo{
			{
				Type:   "shipping",
				Format: "pdf",
				URL:    peResp.LabelURL,
			},
		}
	}

	return resp, nil
}

// TrackShipment отслеживает отправление
func (a *PostExpressAdapter) TrackShipment(ctx context.Context, trackingNumber string) (*interfaces.TrackingResponse, error) {
	if a.service == nil {
		return nil, fmt.Errorf("post Express service not initialized")
	}

	peResp, err := a.service.TrackShipment(ctx, trackingNumber)
	if err != nil {
		return nil, fmt.Errorf("post Express tracking failed: %w", err)
	}

	// Маппинг событий
	events := make([]interfaces.TrackingEvent, 0, len(peResp.Events))
	for _, event := range peResp.Events {
		events = append(events, interfaces.TrackingEvent{
			Timestamp:   event.Timestamp,
			Status:      mapPostExpressStatus(event.Status),
			Description: event.Description,
			Location:    event.Location,
			Details:     event.Details,
		})
	}

	resp := &interfaces.TrackingResponse{
		TrackingNumber:  trackingNumber,
		Status:          mapPostExpressStatus(peResp.Status),
		StatusText:      peResp.StatusText,
		CurrentLocation: peResp.CurrentLocation,
		EstimatedDate:   peResp.EstimatedDate,
		DeliveredDate:   peResp.DeliveredDate,
		Events:          events,
	}

	// Добавляем подтверждение доставки если есть
	if peResp.ProofOfDelivery != nil {
		resp.ProofOfDelivery = &interfaces.ProofOfDelivery{
			RecipientName: peResp.ProofOfDelivery.RecipientName,
			SignatureURL:  peResp.ProofOfDelivery.SignatureURL,
			PhotoURL:      peResp.ProofOfDelivery.PhotoURL,
			DeliveredAt:   peResp.ProofOfDelivery.DeliveredAt,
			Notes:         peResp.ProofOfDelivery.Notes,
		}
	}

	return resp, nil
}

// CancelShipment отменяет отправление
func (a *PostExpressAdapter) CancelShipment(ctx context.Context, externalID string) error {
	if a.service == nil {
		return fmt.Errorf("post Express service not initialized")
	}

	return a.service.CancelShipment(ctx, externalID, "Отменено пользователем")
}

// GetLabel получает этикетку отправления
func (a *PostExpressAdapter) GetLabel(ctx context.Context, shipmentID string) (*interfaces.LabelResponse, error) {
	// Post Express возвращает URL этикетки в response CreateShipment
	// Здесь можно реализовать отдельный запрос если API поддерживает
	return &interfaces.LabelResponse{
		Labels: []interfaces.LabelInfo{
			{
				Type:   "shipping",
				Format: "pdf",
				URL:    "", // TODO: реализовать если API поддерживает отдельный endpoint
			},
		},
	}, nil
}

// ValidateAddress проверяет корректность адреса
func (a *PostExpressAdapter) ValidateAddress(ctx context.Context, address *interfaces.Address) (*interfaces.AddressValidationResponse, error) {
	// Базовая валидация
	errors := []string{}

	if address.City == "" {
		errors = append(errors, "City is required")
	}
	if address.PostalCode == "" {
		errors = append(errors, "Postal code is required")
	}
	if address.Street == "" {
		errors = append(errors, "Street address is required")
	}

	isValid := len(errors) == 0

	// Проверяем доступность доставки в город через список офисов
	deliveryAvailable := false
	zone := "national"

	if a.service != nil && address.City != "" {
		officeReq := &postexpress.OfficeListRequest{
			City: address.City,
		}
		officesResp, err := a.service.GetOffices(ctx, officeReq)
		if err == nil && len(officesResp.Offices) > 0 {
			deliveryAvailable = true

			// Определяем зону (локальная или национальная)
			config := a.service.GetConfig()
			if address.City == config.Brand { // Если город совпадает с базовым - локальная зона
				zone = "local"
			}
		}
	}

	return &interfaces.AddressValidationResponse{
		IsValid:           isValid,
		ValidationErrors:  errors,
		DeliveryAvailable: deliveryAvailable,
		Zone:              zone,
	}, nil
}

// HandleWebhook обрабатывает webhook от Post Express
func (a *PostExpressAdapter) HandleWebhook(ctx context.Context, payload []byte, headers map[string]string) (*interfaces.WebhookResponse, error) {
	// TODO: реализовать парсинг webhook от Post Express
	// Формат webhook нужно уточнить в документации API

	response := &interfaces.WebhookResponse{
		Processed:     true,
		Timestamp:     time.Now(),
		StatusDetails: "Post Express webhook received and processed",
	}

	return response, nil
}

// Вспомогательные функции

// calculateTotalWeight рассчитывает общий вес всех посылок
func calculateTotalWeight(packages []interfaces.Package) float64 {
	total := 0.0
	for _, pkg := range packages {
		total += pkg.Weight
	}
	return total
}

// mapPostExpressStatus маппит статус Post Express в универсальный статус
func mapPostExpressStatus(peStatus string) string {
	mapping := map[string]string{
		postexpress.StatusCreated:           interfaces.StatusPending,
		postexpress.StatusPickupScheduled:   interfaces.StatusConfirmed,
		postexpress.StatusPickedUp:          interfaces.StatusPickedUp,
		postexpress.StatusInTransit:         interfaces.StatusInTransit,
		postexpress.StatusArrived:           interfaces.StatusInTransit,
		postexpress.StatusOutForDelivery:    interfaces.StatusOutForDelivery,
		postexpress.StatusDelivered:         interfaces.StatusDelivered,
		postexpress.StatusDeliveryAttempted: interfaces.StatusDeliveryAttempted,
		postexpress.StatusReturning:         interfaces.StatusReturning,
		postexpress.StatusReturned:          interfaces.StatusReturned,
		postexpress.StatusCancelled:         interfaces.StatusCancelled,
		postexpress.StatusLost:              interfaces.StatusLost,
		postexpress.StatusDamaged:           interfaces.StatusDamaged,
	}

	if mapped, ok := mapping[peStatus]; ok {
		return mapped
	}

	return interfaces.StatusPending // default
}

// contains проверяет наличие строки в слайсе
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// buildPosebneUsluge строит строку дополнительных услуг через запятую
func buildPosebneUsluge(services []string) string {
	usluge := []string{"PNA"} // Приём на адресе для курьера - всегда включён

	// Добавляем SMS если запрошен
	if contains(services, "sms") {
		usluge = append(usluge, "SMS")
	}

	// Добавляем OTK (откупнина) если есть COD
	if contains(services, "cod") {
		usluge = append(usluge, "OTK")
	}

	// Добавляем VD (ценная посылка) если есть insurance
	if contains(services, "insurance") {
		usluge = append(usluge, "VD")
	}

	// Соединяем через запятую
	result := ""
	for i, u := range usluge {
		if i > 0 {
			result += ","
		}
		result += u
	}
	return result
}
