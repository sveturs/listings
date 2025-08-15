package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"backend/internal/proj/postexpress/models"
	"backend/internal/proj/postexpress/storage"
	"backend/pkg/logger"
)

// ServiceImpl представляет реализацию сервиса Post Express
type ServiceImpl struct {
	repo      storage.Repository
	wspClient WSPClient
	logger    logger.Logger
	config    *ServiceConfig
}

// ServiceConfig представляет конфигурацию сервиса
type ServiceConfig struct {
	DefaultWarehouseCode string
	PickupExpiryDays     int
	EnableAutoTracking   bool
	TrackingInterval     time.Duration
	MaxRetries           int
}

// NewService создает новый экземпляр сервиса
func NewService(repo storage.Repository, wspClient WSPClient, logger logger.Logger, config *ServiceConfig) Service {
	return &ServiceImpl{
		repo:      repo,
		wspClient: wspClient,
		logger:    logger,
		config:    config,
	}
}

// =============================================================================
// НАСТРОЙКИ
// =============================================================================

func (s *ServiceImpl) GetSettings(ctx context.Context) (*models.PostExpressSettings, error) {
	settings, err := s.repo.GetSettings(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}

	if settings == nil {
		// Возвращаем настройки по умолчанию
		return &models.PostExpressSettings{
			APIEndpoint:            "https://wsp.postexpress.rs/api/Transakcija",
			SenderName:             "Sve Tu Platform",
			SenderAddress:          "Улица Микија Манојловића 53",
			SenderCity:             "Нови Сад",
			SenderPostalCode:       "21000",
			SenderPhone:            "+381 21 XXX-XXXX",
			SenderEmail:            "shipping@svetu.rs",
			Enabled:                false,
			TestMode:               true,
			AutoPrintLabels:        false,
			AutoTrackShipments:     true,
			NotifyOnPickup:         true,
			NotifyOnDelivery:       true,
			NotifyOnFailedDelivery: true,
		}, nil
	}

	return settings, nil
}

func (s *ServiceImpl) UpdateSettings(ctx context.Context, settings *models.PostExpressSettings) error {
	if settings.ID == 0 {
		return s.repo.CreateSettings(ctx, settings)
	}
	return s.repo.UpdateSettings(ctx, settings)
}

// =============================================================================
// ЛОКАЦИИ
// =============================================================================

func (s *ServiceImpl) SearchLocations(ctx context.Context, query string) ([]*models.PostExpressLocation, error) {
	// Сначала ищем в локальной базе
	locations, err := s.repo.SearchLocations(ctx, query, 50)
	if err != nil {
		return nil, fmt.Errorf("failed to search local locations: %w", err)
	}

	// Если найдено достаточно результатов, возвращаем их
	if len(locations) >= 5 {
		return locations, nil
	}

	// Иначе запрашиваем у Post Express API
	s.logger.Debug("Searching locations via WSP API", "query", query)

	wspLocations, err := s.wspClient.GetLocations(ctx, query)
	if err != nil {
		s.logger.Error("Failed to get locations from WSP API: %v", err)
		// Возвращаем локальные результаты если API недоступен
		return locations, nil
	}

	// Преобразуем и сохраняем новые локации
	newLocations := make([]*models.PostExpressLocation, 0, len(wspLocations))
	for _, wspLoc := range wspLocations {
		// Проверяем, есть ли уже такая локация
		existing, err := s.repo.GetLocationByPostExpressID(ctx, wspLoc.ID)
		if err != nil {
			s.logger.Error("Failed to check existing location", "id", wspLoc.ID, "error", err.Error())
			continue
		}

		if existing != nil {
			newLocations = append(newLocations, existing)
			continue
		}

		// Создаем новую локацию
		location := &models.PostExpressLocation{
			PostExpressID:   wspLoc.ID,
			Name:            wspLoc.Name,
			PostalCode:      stringPtr(wspLoc.PostalCode),
			Municipality:    stringPtr(wspLoc.Municipality),
			IsActive:        true,
			SupportsCOD:     true,
			SupportsExpress: true,
		}

		err = s.repo.CreateLocation(ctx, location)
		if err != nil {
			s.logger.Error("Failed to create location", "location", wspLoc.Name, "error", err.Error())
			continue
		}

		newLocations = append(newLocations, location)
	}

	// Объединяем локальные и новые результаты
	allLocations := append(locations, newLocations...)

	// Удаляем дубликаты и ограничиваем количество
	seen := make(map[int]bool)
	result := make([]*models.PostExpressLocation, 0, 50)

	for _, loc := range allLocations {
		if !seen[loc.PostExpressID] && len(result) < 50 {
			seen[loc.PostExpressID] = true
			result = append(result, loc)
		}
	}

	return result, nil
}

func (s *ServiceImpl) GetLocationByID(ctx context.Context, id int) (*models.PostExpressLocation, error) {
	return s.repo.GetLocationByID(ctx, id)
}

func (s *ServiceImpl) SyncLocations(ctx context.Context) error {
	s.logger.Info("Starting locations sync")

	// Получаем популярные города для синхронизации
	cities := []string{"Београд", "Нови Сад", "Ниш", "Крагујевац", "Суботица", "Панчево", "Чачак", "Нова Пазова"}

	allLocations := make([]*models.PostExpressLocation, 0)

	for _, city := range cities {
		s.logger.Debug("Syncing city", "city", city)

		wspLocations, err := s.wspClient.GetLocations(ctx, city)
		if err != nil {
			s.logger.Error("Failed to get locations for city", "city", city, "error", err.Error())
			continue
		}

		for _, wspLoc := range wspLocations {
			location := &models.PostExpressLocation{
				PostExpressID:   wspLoc.ID,
				Name:            wspLoc.Name,
				PostalCode:      stringPtr(wspLoc.PostalCode),
				Municipality:    stringPtr(wspLoc.Municipality),
				IsActive:        true,
				SupportsCOD:     true,
				SupportsExpress: true,
			}
			allLocations = append(allLocations, location)
		}

		// Небольшая пауза между запросами
		time.Sleep(100 * time.Millisecond)
	}

	if len(allLocations) > 0 {
		err := s.repo.BulkUpsertLocations(ctx, allLocations)
		if err != nil {
			return fmt.Errorf("failed to bulk upsert locations: %w", err)
		}
	}

	s.logger.Info("Locations sync completed", "count", len(allLocations))
	return nil
}

// =============================================================================
// ОТДЕЛЕНИЯ
// =============================================================================

func (s *ServiceImpl) GetOfficesByLocation(ctx context.Context, locationID int) ([]*models.PostExpressOffice, error) {
	return s.repo.GetOfficesByLocationID(ctx, locationID)
}

func (s *ServiceImpl) GetOfficeByCode(ctx context.Context, code string) (*models.PostExpressOffice, error) {
	return s.repo.GetOfficeByCode(ctx, code)
}

func (s *ServiceImpl) SyncOffices(ctx context.Context) error {
	s.logger.Info("Starting offices sync")

	// Получаем все активные локации
	locations, err := s.repo.SearchLocations(ctx, "", 1000)
	if err != nil {
		return fmt.Errorf("failed to get locations for office sync: %w", err)
	}

	allOffices := make([]*models.PostExpressOffice, 0)

	for _, location := range locations {
		s.logger.Debug("Syncing offices for location", "location", location.Name, "id", location.PostExpressID)

		wspOffices, err := s.wspClient.GetOffices(ctx, location.PostExpressID)
		if err != nil {
			s.logger.Error("Failed to get offices for location", "location", location.Name, "error", err.Error())
			continue
		}

		for _, wspOffice := range wspOffices {
			office := &models.PostExpressOffice{
				OfficeCode:      wspOffice.Code,
				LocationID:      &location.ID,
				Name:            wspOffice.Name,
				Address:         wspOffice.Address,
				Phone:           stringPtr(wspOffice.Phone),
				AcceptsPackages: true,
				IssuesPackages:  true,
				IsActive:        true,
			}

			// Парсим рабочие часы если есть
			if wspOffice.WorkingHours != "" {
				// TODO: Парсинг рабочих часов в JSON формат
			}

			allOffices = append(allOffices, office)
		}

		// Пауза между запросами
		time.Sleep(200 * time.Millisecond)
	}

	if len(allOffices) > 0 {
		err := s.repo.BulkUpsertOffices(ctx, allOffices)
		if err != nil {
			return fmt.Errorf("failed to bulk upsert offices: %w", err)
		}
	}

	s.logger.Info("Offices sync completed", "count", len(allOffices))
	return nil
}

// =============================================================================
// РАСЧЕТ СТОИМОСТИ
// =============================================================================

func (s *ServiceImpl) CalculateRate(ctx context.Context, req *models.CalculateRateRequest) (*models.CalculateRateResponse, error) {
	// Получаем тариф для указанного веса
	rate, err := s.repo.GetRateForWeight(ctx, req.WeightKg)
	if err != nil {
		return nil, fmt.Errorf("failed to get rate for weight: %w", err)
	}

	if rate == nil {
		return &models.CalculateRateResponse{
			ServiceAvailable: false,
		}, nil
	}

	// Базовая стоимость
	basePrice := rate.BasePrice

	// Расчет страхования
	var insuranceFee float64
	if req.DeclaredValue != nil && *req.DeclaredValue > rate.InsuranceIncludedUpTo {
		excess := *req.DeclaredValue - rate.InsuranceIncludedUpTo
		insuranceFee = excess * (rate.InsuranceRatePercent / 100)
	}

	// Комиссия за наложенный платеж
	var codFee float64
	if req.CODAmount != nil && *req.CODAmount > 0 {
		codFee = rate.CODFee
	}

	// Проверка размеров
	serviceAvailable := true
	if req.LengthCm != nil && *req.LengthCm > rate.MaxLengthCm {
		serviceAvailable = false
	}
	if req.WidthCm != nil && *req.WidthCm > rate.MaxWidthCm {
		serviceAvailable = false
	}
	if req.HeightCm != nil && *req.HeightCm > rate.MaxHeightCm {
		serviceAvailable = false
	}
	if req.LengthCm != nil && req.WidthCm != nil && req.HeightCm != nil {
		totalDimensions := *req.LengthCm + *req.WidthCm + *req.HeightCm
		if totalDimensions > rate.MaxDimensionsSumCm {
			serviceAvailable = false
		}
	}

	totalPrice := basePrice + insuranceFee + codFee

	return &models.CalculateRateResponse{
		BasePrice:        basePrice,
		InsuranceFee:     insuranceFee,
		CODFee:           codFee,
		TotalPrice:       totalPrice,
		DeliveryDaysMin:  rate.DeliveryDaysMin,
		DeliveryDaysMax:  rate.DeliveryDaysMax,
		ServiceAvailable: serviceAvailable,
	}, nil
}

func (s *ServiceImpl) GetRates(ctx context.Context) ([]*models.PostExpressRate, error) {
	return s.repo.GetRates(ctx)
}

// =============================================================================
// ОТПРАВЛЕНИЯ
// =============================================================================

func (s *ServiceImpl) CreateShipment(ctx context.Context, req *models.CreateShipmentRequest) (*models.PostExpressShipment, error) {
	// Получаем настройки
	settings, err := s.GetSettings(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}

	if !settings.Enabled {
		return nil, fmt.Errorf("Post Express integration is disabled")
	}

	// Рассчитываем стоимость
	rateReq := &models.CalculateRateRequest{
		WeightKg:      req.WeightKg,
		LengthCm:      req.LengthCm,
		WidthCm:       req.WidthCm,
		HeightCm:      req.HeightCm,
		DeclaredValue: req.DeclaredValue,
		CODAmount:     req.CODAmount,
		ServiceType:   req.ServiceType,
	}

	rateResp, err := s.CalculateRate(ctx, rateReq)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate rate: %w", err)
	}

	if !rateResp.ServiceAvailable {
		return nil, fmt.Errorf("service not available for specified parameters")
	}

	// Создаем отправление в БД
	shipment := &models.PostExpressShipment{
		MarketplaceOrderID: req.OrderID,
		StorefrontOrderID:  req.StorefrontOrderID,

		// Отправитель (из настроек)
		SenderName:       settings.SenderName,
		SenderAddress:    settings.SenderAddress,
		SenderCity:       settings.SenderCity,
		SenderPostalCode: settings.SenderPostalCode,
		SenderPhone:      settings.SenderPhone,
		SenderEmail:      &settings.SenderEmail,

		// Получатель
		RecipientName:       req.RecipientName,
		RecipientAddress:    req.RecipientAddress,
		RecipientCity:       req.RecipientCity,
		RecipientPostalCode: req.RecipientPostalCode,
		RecipientPhone:      req.RecipientPhone,
		RecipientEmail:      req.RecipientEmail,

		// Параметры посылки
		WeightKg:      req.WeightKg,
		LengthCm:      req.LengthCm,
		WidthCm:       req.WidthCm,
		HeightCm:      req.HeightCm,
		DeclaredValue: req.DeclaredValue,

		// Услуги
		ServiceType:     getServiceType(req.ServiceType),
		CODAmount:       req.CODAmount,
		InsuranceAmount: req.InsuranceAmount,

		// Стоимость
		BasePrice:    rateResp.BasePrice,
		InsuranceFee: rateResp.InsuranceFee,
		CODFee:       rateResp.CODFee,
		TotalPrice:   rateResp.TotalPrice,

		// Статус
		Status: models.ShipmentStatusCreated,

		// Дополнительно
		DeliveryInstructions: req.DeliveryInstructions,
		Notes:                req.Notes,
	}

	createdShipment, err := s.repo.CreateShipment(ctx, shipment)
	if err != nil {
		return nil, fmt.Errorf("failed to create shipment in database: %w", err)
	}

	// Регистрируем отправление в Post Express
	err = s.registerShipmentWithPostExpress(ctx, createdShipment)
	if err != nil {
		s.logger.Error("Failed to register shipment with Post Express",
			"shipment_id", createdShipment.ID, "error", err.Error())
		// Не возвращаем ошибку, так как отправление уже создано в БД
		// Можно попробовать зарегистрировать позже
	}

	return createdShipment, nil
}

func (s *ServiceImpl) registerShipmentWithPostExpress(ctx context.Context, shipment *models.PostExpressShipment) error {
	// Подготавливаем запрос для WSP API
	wspReq := &WSPShipmentRequest{
		SenderName:          shipment.SenderName,
		SenderAddress:       shipment.SenderAddress,
		SenderCity:          shipment.SenderCity,
		SenderPostalCode:    shipment.SenderPostalCode,
		SenderPhone:         shipment.SenderPhone,
		RecipientName:       shipment.RecipientName,
		RecipientAddress:    shipment.RecipientAddress,
		RecipientCity:       shipment.RecipientCity,
		RecipientPostalCode: shipment.RecipientPostalCode,
		RecipientPhone:      shipment.RecipientPhone,
		Weight:              shipment.WeightKg,
		ServiceType:         shipment.ServiceType,
		Content:             "Товары из интернет-магазина",
	}

	if shipment.CODAmount != nil {
		wspReq.CODAmount = *shipment.CODAmount
	}
	if shipment.DeclaredValue != nil {
		wspReq.InsuranceAmount = *shipment.DeclaredValue
	}
	if shipment.Notes != nil {
		wspReq.Note = *shipment.Notes
	}

	// Отправляем запрос
	resp, err := s.wspClient.CreateShipment(ctx, wspReq)
	if err != nil {
		return fmt.Errorf("failed to create shipment via WSP API: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("Post Express API error: %s", resp.ErrorMessage)
	}

	// Обновляем отправление данными от Post Express
	shipment.TrackingNumber = &resp.TrackingNumber
	shipment.Barcode = &resp.Barcode
	shipment.PostExpressID = &resp.TrackingNumber // Используем tracking number как ID
	shipment.LabelURL = &resp.LabelURL
	shipment.Status = models.ShipmentStatusRegistered
	now := time.Now()
	shipment.RegisteredAt = &now

	err = s.repo.UpdateShipment(ctx, shipment)
	if err != nil {
		s.logger.Error("Failed to update shipment after registration",
			"shipment_id", shipment.ID, "error", err.Error())
		// Не критичная ошибка
	}

	s.logger.Info("Shipment registered with Post Express",
		"shipment_id", shipment.ID, "tracking_number", resp.TrackingNumber)

	return nil
}

func (s *ServiceImpl) GetShipment(ctx context.Context, id int) (*models.PostExpressShipment, error) {
	return s.repo.GetShipmentByID(ctx, id)
}

func (s *ServiceImpl) GetShipmentByTrackingNumber(ctx context.Context, trackingNumber string) (*models.PostExpressShipment, error) {
	return s.repo.GetShipmentByTrackingNumber(ctx, trackingNumber)
}

func (s *ServiceImpl) ListShipments(ctx context.Context, filters storage.ShipmentFilters) ([]*models.PostExpressShipment, int, error) {
	return s.repo.ListShipments(ctx, filters)
}

func (s *ServiceImpl) UpdateShipmentStatus(ctx context.Context, id int, status models.ShipmentStatus) error {
	return s.repo.UpdateShipmentStatus(ctx, id, status)
}

func (s *ServiceImpl) CancelShipment(ctx context.Context, id int, reason string) error {
	shipment, err := s.repo.GetShipmentByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get shipment: %w", err)
	}

	if shipment == nil {
		return fmt.Errorf("shipment not found")
	}

	// Отменяем в Post Express если есть ID
	if shipment.PostExpressID != nil {
		err = s.wspClient.CancelShipment(ctx, *shipment.PostExpressID)
		if err != nil {
			s.logger.Error("Failed to cancel shipment in Post Express",
				"shipment_id", id, "error", err.Error())
			// Продолжаем отмену в локальной БД
		}
	}

	// Обновляем статус
	shipment.Status = models.ShipmentStatusFailed
	shipment.FailedReason = &reason
	now := time.Now()
	shipment.FailedAt = &now

	return s.repo.UpdateShipment(ctx, shipment)
}

// =============================================================================
// ДОКУМЕНТЫ
// =============================================================================

func (s *ServiceImpl) GetShipmentLabel(ctx context.Context, id int) ([]byte, error) {
	shipment, err := s.repo.GetShipmentByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipment: %w", err)
	}

	if shipment == nil {
		return nil, fmt.Errorf("shipment not found")
	}

	if shipment.PostExpressID == nil {
		return nil, fmt.Errorf("shipment not registered with Post Express")
	}

	return s.wspClient.PrintLabel(ctx, *shipment.PostExpressID)
}

func (s *ServiceImpl) GetShipmentInvoice(ctx context.Context, id int) ([]byte, error) {
	// TODO: Реализовать получение накладной
	return nil, fmt.Errorf("invoice generation not implemented")
}

func (s *ServiceImpl) PrintShipmentDocuments(ctx context.Context, id int) error {
	// TODO: Реализовать печать документов
	return fmt.Errorf("document printing not implemented")
}

// =============================================================================
// ОТСЛЕЖИВАНИЕ
// =============================================================================

func (s *ServiceImpl) TrackShipment(ctx context.Context, trackingNumber string) ([]*models.TrackingEvent, error) {
	// Получаем отправление
	shipment, err := s.repo.GetShipmentByTrackingNumber(ctx, trackingNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipment: %w", err)
	}

	if shipment == nil {
		return nil, fmt.Errorf("shipment not found")
	}

	// Обновляем статус из Post Express
	err = s.UpdateTrackingStatus(ctx, trackingNumber)
	if err != nil {
		s.logger.Error("Failed to update tracking status for %s: %v", trackingNumber, err)
	}

	// Возвращаем локальные события отслеживания
	return s.repo.GetTrackingEventsByShipmentID(ctx, shipment.ID)
}

func (s *ServiceImpl) UpdateTrackingStatus(ctx context.Context, trackingNumber string) error {
	// Получаем статус от Post Express
	tracking, err := s.wspClient.GetShipmentStatus(ctx, trackingNumber)
	if err != nil {
		return fmt.Errorf("failed to get tracking status: %w", err)
	}

	// Получаем отправление
	shipment, err := s.repo.GetShipmentByTrackingNumber(ctx, trackingNumber)
	if err != nil {
		return fmt.Errorf("failed to get shipment: %w", err)
	}

	if shipment == nil {
		return fmt.Errorf("shipment not found")
	}

	// Обновляем статус отправления
	newStatus := mapWSPStatusToShipmentStatus(tracking.Status)
	if newStatus != shipment.Status {
		err = s.repo.UpdateShipmentStatus(ctx, shipment.ID, newStatus)
		if err != nil {
			return fmt.Errorf("failed to update shipment status: %w", err)
		}
	}

	// Создаем события отслеживания
	for _, wspEvent := range tracking.Events {
		event := &models.TrackingEvent{
			ShipmentID:       shipment.ID,
			EventCode:        wspEvent.Code,
			EventDescription: wspEvent.Description,
			EventLocation:    &wspEvent.Location,
			EventTimestamp:   parseWSPDateTime(wspEvent.Date, wspEvent.Time),
		}

		err = s.repo.CreateTrackingEvent(ctx, event)
		if err != nil {
			s.logger.Error("Failed to create tracking event",
				"shipment_id", shipment.ID, "event_code", wspEvent.Code, "error", err.Error())
		}
	}

	return nil
}

func (s *ServiceImpl) SyncAllActiveShipments(ctx context.Context) error {
	// Получаем все активные отправления
	filters := storage.ShipmentFilters{
		Page:     1,
		PageSize: 1000,
	}

	shipments, _, err := s.repo.ListShipments(ctx, filters)
	if err != nil {
		return fmt.Errorf("failed to get active shipments: %w", err)
	}

	var successCount, errorCount int

	for _, shipment := range shipments {
		if shipment.TrackingNumber == nil {
			continue
		}

		// Пропускаем завершенные отправления
		if shipment.Status == models.ShipmentStatusDelivered ||
			shipment.Status == models.ShipmentStatusReturned {
			continue
		}

		err = s.UpdateTrackingStatus(ctx, *shipment.TrackingNumber)
		if err != nil {
			s.logger.Error("Failed to update tracking for shipment",
				"shipment_id", shipment.ID, "tracking", *shipment.TrackingNumber, "error", err.Error())
			errorCount++
		} else {
			successCount++
		}

		// Пауза между запросами
		time.Sleep(100 * time.Millisecond)
	}

	s.logger.Info("Tracking sync completed",
		"success_count", successCount, "error_count", errorCount)

	return nil
}

// =============================================================================
// СКЛАД И САМОВЫВОЗ
// =============================================================================

func (s *ServiceImpl) GetWarehouses(ctx context.Context) ([]*models.Warehouse, error) {
	return s.repo.GetWarehouses(ctx)
}

func (s *ServiceImpl) GetWarehouseByCode(ctx context.Context, code string) (*models.Warehouse, error) {
	return s.repo.GetWarehouseByCode(ctx, code)
}

func (s *ServiceImpl) CreatePickupOrder(ctx context.Context, req *models.CreatePickupOrderRequest) (*models.WarehousePickupOrder, error) {
	// Получаем склад по умолчанию
	warehouse, err := s.repo.GetWarehouseByCode(ctx, s.config.DefaultWarehouseCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get default warehouse: %w", err)
	}

	if warehouse == nil {
		return nil, fmt.Errorf("default warehouse not found")
	}

	// Создаем заказ на самовывоз
	order := &models.WarehousePickupOrder{
		WarehouseID:        warehouse.ID,
		MarketplaceOrderID: req.OrderID,
		StorefrontOrderID:  req.StorefrontOrderID,
		CustomerName:       req.CustomerName,
		CustomerPhone:      req.CustomerPhone,
		CustomerEmail:      req.CustomerEmail,
		Notes:              req.Notes,
	}

	return s.repo.CreatePickupOrder(ctx, order)
}

func (s *ServiceImpl) GetPickupOrder(ctx context.Context, id int) (*models.WarehousePickupOrder, error) {
	return s.repo.GetPickupOrderByID(ctx, id)
}

func (s *ServiceImpl) GetPickupOrderByCode(ctx context.Context, code string) (*models.WarehousePickupOrder, error) {
	return s.repo.GetPickupOrderByCode(ctx, code)
}

func (s *ServiceImpl) ConfirmPickup(ctx context.Context, id int, confirmedBy string, documentType string, documentNumber string) error {
	order, err := s.repo.GetPickupOrderByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get pickup order: %w", err)
	}

	if order == nil {
		return fmt.Errorf("pickup order not found")
	}

	if order.Status != models.PickupOrderStatusReady {
		return fmt.Errorf("pickup order is not ready for pickup")
	}

	// Обновляем статус
	now := time.Now()
	order.Status = models.PickupOrderStatusPickedUp
	order.PickedUpAt = &now
	order.PickupConfirmedBy = &confirmedBy
	order.IDDocumentType = &documentType
	order.IDDocumentNumber = &documentNumber

	return s.repo.UpdatePickupOrder(ctx, order)
}

func (s *ServiceImpl) CancelPickupOrder(ctx context.Context, id int, reason string) error {
	order, err := s.repo.GetPickupOrderByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get pickup order: %w", err)
	}

	if order == nil {
		return fmt.Errorf("pickup order not found")
	}

	order.Status = models.PickupOrderStatusCancelled
	order.Notes = &reason

	return s.repo.UpdatePickupOrder(ctx, order)
}

// =============================================================================
// СТАТИСТИКА
// =============================================================================

func (s *ServiceImpl) GetShipmentStatistics(ctx context.Context, filters storage.StatisticsFilters) (*storage.ShipmentStatistics, error) {
	return s.repo.GetShipmentStatistics(ctx, filters)
}

func (s *ServiceImpl) GetWarehouseStatistics(ctx context.Context, warehouseID int) (*storage.WarehouseStatistics, error) {
	return s.repo.GetWarehouseStatistics(ctx, warehouseID)
}

// =============================================================================
// ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ
// =============================================================================

func getServiceType(serviceType string) string {
	if serviceType == "" {
		return "danas_za_sutra"
	}
	return serviceType
}

func mapWSPStatusToShipmentStatus(wspStatus string) models.ShipmentStatus {
	switch strings.ToLower(wspStatus) {
	case "preuzeto", "picked_up":
		return models.ShipmentStatusPickedUp
	case "u_transportu", "in_transit":
		return models.ShipmentStatusInTransit
	case "dostavljeno", "delivered":
		return models.ShipmentStatusDelivered
	case "neuspesno", "failed":
		return models.ShipmentStatusFailed
	case "vraceno", "returned":
		return models.ShipmentStatusReturned
	default:
		return models.ShipmentStatusRegistered
	}
}

func parseWSPDateTime(date, timeStr string) time.Time {
	// TODO: Реализовать корректный парсинг даты и времени из WSP формата
	dateTimeStr := date + " " + timeStr
	parsed, err := time.Parse("02.01.2006 15:04", dateTimeStr)
	if err != nil {
		return time.Now()
	}
	return parsed
}
