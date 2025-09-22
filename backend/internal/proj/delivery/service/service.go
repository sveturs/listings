package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"

	"backend/internal/proj/delivery/attributes"
	"backend/internal/proj/delivery/calculator"
	"backend/internal/proj/delivery/factory"
	"backend/internal/proj/delivery/interfaces"
	"backend/internal/proj/delivery/models"
	"backend/internal/proj/delivery/notifications"
	"backend/internal/proj/delivery/storage"
	notifService "backend/internal/proj/notifications/service"
)

// Service - основной сервис управления доставкой
type Service struct {
	db              *sqlx.DB
	storage         *storage.Storage
	calculator      *calculator.Service
	attributes      *attributes.Service
	providerFactory *factory.ProviderFactory
	notifications   *notifications.DeliveryNotificationIntegration
	cache           map[string]interfaces.DeliveryProvider // кеш провайдеров
}

// NewService создает новый экземпляр сервиса доставки
func NewService(db *sqlx.DB, providerFactory *factory.ProviderFactory) *Service {
	return &Service{
		db:              db,
		storage:         storage.NewStorage(db),
		calculator:      calculator.NewService(db),
		attributes:      attributes.NewService(db),
		providerFactory: providerFactory,
		cache:           make(map[string]interfaces.DeliveryProvider),
	}
}

// SetNotificationService устанавливает сервис уведомлений
func (s *Service) SetNotificationService(notifService notifService.NotificationServiceInterface) {
	if notifService != nil {
		s.notifications = notifications.NewDeliveryNotificationIntegration(notifService)
		log.Info().Msg("Notification service integrated with delivery module")
	}
}

// CalculateDelivery - рассчитывает стоимость доставки для всех доступных провайдеров
func (s *Service) CalculateDelivery(ctx context.Context, req *calculator.CalculationRequest) (*calculator.CalculationResponse, error) {
	return s.calculator.Calculate(ctx, req)
}

// GetProductAttributes - получает атрибуты доставки товара
func (s *Service) GetProductAttributes(ctx context.Context, productID int, productType string) (*models.DeliveryAttributes, error) {
	return s.attributes.GetProductAttributes(ctx, productID, productType)
}

// UpdateProductAttributes - обновляет атрибуты доставки товара
func (s *Service) UpdateProductAttributes(ctx context.Context, productID int, productType string, attrs *models.DeliveryAttributes) error {
	return s.attributes.UpdateProductAttributes(ctx, productID, productType, attrs)
}

// GetCategoryDefaults - получает дефолтные атрибуты категории
func (s *Service) GetCategoryDefaults(ctx context.Context, categoryID int) (*models.CategoryDefaults, error) {
	return s.attributes.GetCategoryDefaults(ctx, categoryID)
}

// UpdateCategoryDefaults - обновляет дефолтные атрибуты категории
func (s *Service) UpdateCategoryDefaults(ctx context.Context, defaults *models.CategoryDefaults) error {
	return s.attributes.UpdateCategoryDefaults(ctx, defaults)
}

// CreateShipment - создает отправление через выбранного провайдера
func (s *Service) CreateShipment(ctx context.Context, req *CreateShipmentRequest) (*models.Shipment, error) {
	// Получаем провайдера
	provider, err := s.getProvider(ctx, req.ProviderCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	// Подготавливаем запрос для провайдера
	providerReq := &interfaces.ShipmentRequest{
		OrderID:        req.OrderID,
		FromAddress:    req.FromAddress,
		ToAddress:      req.ToAddress,
		Packages:       req.Packages,
		DeliveryType:   req.DeliveryType,
		PickupDate:     req.PickupDate,
		InsuranceValue: req.InsuranceValue,
		CODAmount:      req.CODAmount,
		Services:       req.Services,
		Reference:      req.Reference,
		Notes:          req.Notes,
	}

	// Создаем отправление через провайдера
	providerResp, err := provider.CreateShipment(ctx, providerReq)
	if err != nil {
		return nil, fmt.Errorf("provider failed to create shipment: %w", err)
	}

	// Сохраняем в БД
	shipment := &models.Shipment{
		ProviderID:     req.ProviderID,
		OrderID:        &req.OrderID,
		ExternalID:     &providerResp.ExternalID,
		TrackingNumber: &providerResp.TrackingNumber,
		Status:         providerResp.Status,
		DeliveryCost:   &providerResp.TotalCost,
	}

	// Сохраняем адреса
	senderInfo, _ := json.Marshal(req.FromAddress)
	recipientInfo, _ := json.Marshal(req.ToAddress)
	packageInfo, _ := json.Marshal(req.Packages)
	costBreakdown, _ := json.Marshal(providerResp.CostBreakdown)
	labels, _ := json.Marshal(providerResp.Labels)

	shipment.SenderInfo = senderInfo
	shipment.RecipientInfo = recipientInfo
	shipment.PackageInfo = packageInfo
	shipment.CostBreakdown = costBreakdown
	shipment.Labels = labels

	if req.InsuranceValue > 0 {
		shipment.InsuranceCost = &req.InsuranceValue
	}
	if req.CODAmount > 0 {
		shipment.CODAmount = &req.CODAmount
	}
	if req.PickupDate != nil {
		shipment.PickupDate = req.PickupDate
	}
	if providerResp.EstimatedDate != nil {
		shipment.EstimatedDelivery = providerResp.EstimatedDate
	}

	// Сохраняем в БД
	if err := s.storage.CreateShipment(ctx, shipment); err != nil {
		return nil, fmt.Errorf("failed to save shipment: %w", err)
	}

	// Обновляем заказ если нужно
	if req.OrderID > 0 {
		if err := s.storage.UpdateOrderShipment(ctx, req.OrderID, shipment.ID); err != nil {
			// Логируем ошибку но не прерываем
			fmt.Printf("Failed to update order shipment: %v\n", err)
		}
	}

	return shipment, nil
}

// TrackShipment - отслеживает отправление
func (s *Service) TrackShipment(ctx context.Context, trackingNumber string) (*TrackingInfo, error) {
	// Получаем отправление из БД
	shipment, err := s.storage.GetShipmentByTracking(ctx, trackingNumber)
	if err != nil {
		return nil, fmt.Errorf("shipment not found: %w", err)
	}

	// Получаем провайдера
	providerInfo, err := s.storage.GetProvider(ctx, shipment.ProviderID)
	if err != nil {
		return nil, fmt.Errorf("provider not found: %w", err)
	}

	provider, err := s.getProvider(ctx, providerInfo.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	// Запрашиваем статус у провайдера
	trackingResp, err := provider.TrackShipment(ctx, trackingNumber)
	if err != nil {
		// Если провайдер недоступен, возвращаем последний известный статус
		events, _ := s.storage.GetTrackingEvents(ctx, shipment.ID)

		return &TrackingInfo{
			ShipmentID:      shipment.ID,
			TrackingNumber:  trackingNumber,
			Status:          shipment.Status,
			CurrentLocation: "", // неизвестно
			Events:          convertEvents(events),
			LastUpdated:     shipment.UpdatedAt,
		}, nil
	}

	// Обновляем статус отправления и отправляем уведомление
	if shipment.Status != trackingResp.Status {
		oldStatus := shipment.Status
		shipment.Status = trackingResp.Status
		if trackingResp.DeliveredDate != nil {
			shipment.ActualDeliveryDate = trackingResp.DeliveredDate
		}
		if err := s.storage.UpdateShipmentStatus(ctx, shipment.ID, shipment.Status, shipment.ActualDeliveryDate); err != nil {
			log.Error().Err(err).Int("shipment_id", shipment.ID).Msg("Failed to update shipment status")
		}

		// Отправляем уведомление об изменении статуса
		s.sendStatusNotification(ctx, shipment, oldStatus, trackingResp.Status, trackingResp.CurrentLocation, trackingResp.StatusText)
	}

	// Сохраняем новые события
	for _, event := range trackingResp.Events {
		trackEvent := &models.TrackingEvent{
			ShipmentID:  shipment.ID,
			ProviderID:  shipment.ProviderID,
			EventTime:   event.Timestamp,
			Status:      event.Status,
			Location:    &event.Location,
			Description: &event.Description,
		}
		if err := s.storage.CreateTrackingEvent(ctx, trackEvent); err != nil {
			log.Error().Err(err).Int("shipment_id", shipment.ID).Msg("Failed to create tracking event")
		}
	}

	return &TrackingInfo{
		ShipmentID:      shipment.ID,
		TrackingNumber:  trackingNumber,
		Status:          trackingResp.Status,
		StatusText:      trackingResp.StatusText,
		CurrentLocation: trackingResp.CurrentLocation,
		EstimatedDate:   trackingResp.EstimatedDate,
		DeliveredDate:   trackingResp.DeliveredDate,
		Events:          trackingResp.Events,
		ProofOfDelivery: trackingResp.ProofOfDelivery,
		LastUpdated:     shipment.UpdatedAt,
	}, nil
}

// CancelShipment - отменяет отправление
func (s *Service) CancelShipment(ctx context.Context, shipmentID int, reason string) error {
	// Получаем отправление
	shipment, err := s.storage.GetShipment(ctx, shipmentID)
	if err != nil {
		return fmt.Errorf("shipment not found: %w", err)
	}

	// Проверяем статус
	if shipment.Status == models.ShipmentStatusDelivered {
		return fmt.Errorf("cannot cancel delivered shipment")
	}
	if shipment.Status == models.ShipmentStatusCancelled {
		return fmt.Errorf("shipment already canceled")
	}

	// Получаем провайдера
	providerInfo, err := s.storage.GetProvider(ctx, shipment.ProviderID)
	if err != nil {
		return fmt.Errorf("provider not found: %w", err)
	}

	provider, err := s.getProvider(ctx, providerInfo.Code)
	if err != nil {
		return fmt.Errorf("failed to get provider: %w", err)
	}

	// Отменяем у провайдера
	if shipment.ExternalID != nil {
		if err := provider.CancelShipment(ctx, *shipment.ExternalID); err != nil {
			return fmt.Errorf("provider failed to cancel: %w", err)
		}
	}

	// Обновляем статус в БД
	return s.storage.UpdateShipmentStatus(ctx, shipmentID, models.ShipmentStatusCancelled, nil)
}

// GetProviders - получает список доступных провайдеров
func (s *Service) GetProviders(ctx context.Context, activeOnly bool) ([]models.Provider, error) {
	return s.storage.GetProviders(ctx, activeOnly)
}

// getProvider - получает экземпляр провайдера
func (s *Service) getProvider(ctx context.Context, code string) (interfaces.DeliveryProvider, error) {
	// Проверяем кеш
	if provider, ok := s.cache[code]; ok {
		return provider, nil
	}

	// Создаем через фабрику
	provider, err := s.providerFactory.CreateProvider(code)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кеш
	s.cache[code] = provider

	return provider, nil
}

// ApplyCategoryDefaults - применяет дефолтные атрибуты к товарам категории
func (s *Service) ApplyCategoryDefaults(ctx context.Context, categoryID int) (int, error) {
	return s.attributes.ApplyCategoryDefaultsToProducts(ctx, categoryID)
}

// convertEvents - конвертирует события из БД в интерфейс
func convertEvents(dbEvents []models.TrackingEvent) []interfaces.TrackingEvent {
	events := make([]interfaces.TrackingEvent, len(dbEvents))
	for i, e := range dbEvents {
		events[i] = interfaces.TrackingEvent{
			Timestamp: e.EventTime,
			Status:    e.Status,
			Description: func() string {
				if e.Description != nil {
					return *e.Description
				}
				return ""
			}(),
			Location: func() string {
				if e.Location != nil {
					return *e.Location
				}
				return ""
			}(),
		}
	}
	return events
}

// sendStatusNotification отправляет уведомление об изменении статуса
func (s *Service) sendStatusNotification(ctx context.Context, shipment *models.Shipment, oldStatus, newStatus, location, description string) {
	if s.notifications == nil {
		log.Debug().Msg("Notification service not configured, skipping notification")
		return
	}

	// Получаем информацию о пользователе (получателе)
	var userID int
	if shipment.OrderID != nil {
		// TODO: Получить user_id из заказа
		// order, err := s.storage.GetOrder(ctx, *shipment.OrderID)
		// if err == nil && order.UserID != nil {
		//     userID = *order.UserID
		// }
		userID = 0 // Временно устанавливаем 0 до реализации
	}

	if userID == 0 {
		log.Debug().Msg("No user ID found for shipment, skipping notification")
		return
	}

	// Получаем трек-номер
	trackingNumber := ""
	if shipment.TrackingNumber != nil {
		trackingNumber = *shipment.TrackingNumber
	}

	// Создаем событие изменения статуса
	event := &notifications.StatusChangeEvent{
		ShipmentID:     shipment.ID,
		TrackingNumber: trackingNumber,
		OldStatus:      oldStatus,
		NewStatus:      newStatus,
		Location:       location,
		Description:    description,
		EventTime:      time.Now(),
	}

	// Отправляем уведомление асинхронно
	go func(ctx context.Context) {
		if err := s.notifications.SendDeliveryStatusUpdate(ctx, userID, event); err != nil {
			log.Error().Err(err).
				Int("shipment_id", shipment.ID).
				Str("tracking_number", trackingNumber).
				Msg("Failed to send delivery notification")
		}
	}(ctx)
}

// Структуры запросов и ответов

// CreateShipmentRequest - запрос создания отправления
type CreateShipmentRequest struct {
	ProviderID     int                  `json:"provider_id"`
	ProviderCode   string               `json:"provider_code"`
	OrderID        int                  `json:"order_id"`
	FromAddress    *interfaces.Address  `json:"from_address"`
	ToAddress      *interfaces.Address  `json:"to_address"`
	Packages       []interfaces.Package `json:"packages"`
	DeliveryType   string               `json:"delivery_type"`
	PickupDate     *time.Time           `json:"pickup_date,omitempty"`
	InsuranceValue float64              `json:"insurance_value,omitempty"`
	CODAmount      float64              `json:"cod_amount,omitempty"`
	Services       []string             `json:"services,omitempty"`
	Reference      string               `json:"reference,omitempty"`
	Notes          string               `json:"notes,omitempty"`
}

// TrackingInfo - информация об отслеживании
type TrackingInfo struct {
	ShipmentID      int                         `json:"shipment_id"`
	TrackingNumber  string                      `json:"tracking_number"`
	Status          string                      `json:"status"`
	StatusText      string                      `json:"status_text"`
	CurrentLocation string                      `json:"current_location,omitempty"`
	EstimatedDate   *time.Time                  `json:"estimated_date,omitempty"`
	DeliveredDate   *time.Time                  `json:"delivered_date,omitempty"`
	Events          []interfaces.TrackingEvent  `json:"events"`
	ProofOfDelivery *interfaces.ProofOfDelivery `json:"proof_of_delivery,omitempty"`
	LastUpdated     time.Time                   `json:"last_updated"`
}

// HandleProviderWebhook - обрабатывает webhook от провайдера доставки
func (s *Service) HandleProviderWebhook(ctx context.Context, providerCode string, payload []byte, headers map[string]string) (*interfaces.WebhookResponse, error) {
	// Получаем провайдера
	provider, err := s.GetProviderByCode(providerCode)
	if err != nil {
		return nil, fmt.Errorf("provider not found: %s", providerCode)
	}

	// Обрабатываем webhook через провайдера
	webhookResponse, err := provider.HandleWebhook(ctx, payload, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to process webhook: %w", err)
	}

	// Если webhook обработан успешно, сохраняем события отслеживания
	if webhookResponse.Processed && webhookResponse.TrackingNumber != "" {
		// Найдем отправление по трек-номеру
		shipment, err := s.storage.GetShipmentByTracking(ctx, webhookResponse.TrackingNumber)
		if err != nil {
			// Если отправление не найдено, просто логируем (возможно, это не наше отправление)
			return webhookResponse, nil
		}

		// Обновляем статус отправления
		var deliveredAt *time.Time
		if webhookResponse.Status == interfaces.StatusDelivered && webhookResponse.DeliveryDetails != nil {
			deliveredAt = &webhookResponse.DeliveryDetails.DeliveredAt
		}
		if err := s.storage.UpdateShipmentStatus(ctx, shipment.ID, webhookResponse.Status, deliveredAt); err != nil {
			return nil, fmt.Errorf("failed to update shipment status: %w", err)
		}

		// Сохраняем события отслеживания
		for _, event := range webhookResponse.Events {
			var location *string
			if event.Location != "" {
				location = &event.Location
			}
			var description *string
			if event.Description != "" {
				description = &event.Description
			}

			trackingEvent := &models.TrackingEvent{
				ShipmentID:  shipment.ID,
				ProviderID:  shipment.ProviderID,
				EventTime:   event.Timestamp,
				Status:      event.Status,
				Location:    location,
				Description: description,
				RawData:     json.RawMessage(fmt.Sprintf(`{"webhook_payload": %q}`, string(payload))),
			}

			if err := s.storage.CreateTrackingEvent(ctx, trackingEvent); err != nil {
				return nil, fmt.Errorf("failed to save tracking event: %w", err)
			}
		}

		// TODO: Отправить уведомления пользователям о изменении статуса
		// s.notificationService.SendStatusUpdate(ctx, shipment, webhookResponse.Status)
	}

	return webhookResponse, nil
}

// GetProviderByCode - получает провайдера по коду
func (s *Service) GetProviderByCode(providerCode string) (interfaces.DeliveryProvider, error) {
	// Проверяем кеш
	if provider, exists := s.cache[providerCode]; exists {
		return provider, nil
	}

	// Получаем провайдера из фабрики
	provider, err := s.providerFactory.CreateProvider(providerCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider %s: %w", providerCode, err)
	}

	// Сохраняем в кеш
	s.cache[providerCode] = provider

	return provider, nil
}

// GetShipment - получает отправление по ID
func (s *Service) GetShipment(ctx context.Context, shipmentID int) (*models.Shipment, error) {
	shipment, err := s.storage.GetShipment(ctx, shipmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipment: %w", err)
	}

	// Загружаем связанные данные
	provider, err := s.storage.GetProvider(ctx, shipment.ProviderID)
	if err == nil {
		shipment.Provider = provider
	}

	// Загружаем события отслеживания
	events, err := s.storage.GetTrackingEvents(ctx, shipmentID)
	if err == nil {
		shipment.Events = events
	}

	return shipment, nil
}

// UpdateProvider - обновляет информацию о провайдере
func (s *Service) UpdateProvider(ctx context.Context, providerID int, update *models.Provider) error {
	// Проверяем существование провайдера
	existing, err := s.storage.GetProvider(ctx, providerID)
	if err != nil {
		return fmt.Errorf("provider not found: %w", err)
	}

	// Обновляем поля
	existing.Name = update.Name
	existing.IsActive = update.IsActive
	existing.SupportsCOD = update.SupportsCOD
	existing.SupportsInsurance = update.SupportsInsurance
	existing.SupportsTracking = update.SupportsTracking
	if update.LogoURL != nil {
		existing.LogoURL = update.LogoURL
	}
	if update.APIConfig != nil {
		existing.APIConfig = update.APIConfig
	}
	if update.Capabilities != nil {
		existing.Capabilities = update.Capabilities
	}

	// Сохраняем в БД
	if err := s.storage.UpdateProvider(ctx, existing); err != nil {
		return fmt.Errorf("failed to update provider: %w", err)
	}

	// Очищаем кеш провайдера
	delete(s.cache, existing.Code)

	return nil
}

// CreatePricingRule - создает новое правило ценообразования
func (s *Service) CreatePricingRule(ctx context.Context, rule *models.PricingRule) (*models.PricingRule, error) {
	// Валидируем правило
	if rule.ProviderID <= 0 {
		return nil, fmt.Errorf("invalid provider ID")
	}
	if rule.RuleType == "" {
		return nil, fmt.Errorf("rule type is required")
	}

	// Проверяем существование провайдера
	if _, err := s.storage.GetProvider(ctx, rule.ProviderID); err != nil {
		return nil, fmt.Errorf("provider not found: %w", err)
	}

	// Создаем правило
	if err := s.storage.CreatePricingRule(ctx, rule); err != nil {
		return nil, fmt.Errorf("failed to create pricing rule: %w", err)
	}

	return rule, nil
}

// GetDeliveryAnalytics - получает аналитику по доставкам
func (s *Service) GetDeliveryAnalytics(ctx context.Context, from, to time.Time, providerID *int) (*DeliveryAnalytics, error) {
	analytics := &DeliveryAnalytics{
		Period: AnalyticsPeriod{
			From: from,
			To:   to,
		},
	}

	// Получаем статистику по отправлениям
	stats, err := s.storage.GetShipmentStatistics(ctx, from, to, providerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipment statistics: %w", err)
	}
	analytics.TotalShipments = stats.TotalShipments
	analytics.TotalCost = stats.TotalCost
	analytics.AverageCost = stats.AverageCost
	analytics.StatusBreakdown = stats.StatusBreakdown

	// Получаем статистику по провайдерам
	providerStats, err := s.storage.GetProviderStatistics(ctx, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider statistics: %w", err)
	}
	// Конвертируем типы из storage в service
	for _, stat := range providerStats {
		analytics.ProviderBreakdown = append(analytics.ProviderBreakdown, ProviderStatistics{
			ProviderID:    stat.ProviderID,
			ProviderName:  stat.ProviderName,
			ShipmentCount: stat.ShipmentCount,
			TotalCost:     stat.TotalCost,
			AverageCost:   stat.AverageCost,
			SuccessRate:   stat.SuccessRate,
		})
	}

	// Получаем топ маршруты
	topRoutes, err := s.storage.GetTopDeliveryRoutes(ctx, from, to, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to get top routes: %w", err)
	}
	// Конвертируем типы из storage в service
	for _, route := range topRoutes {
		analytics.TopRoutes = append(analytics.TopRoutes, RouteStatistics{
			FromCity:      route.FromCity,
			ToCity:        route.ToCity,
			ShipmentCount: route.ShipmentCount,
			AverageCost:   route.AverageCost,
			AverageDays:   route.AverageDays,
		})
	}

	// Получаем среднее время доставки
	avgDeliveryTimes, err := s.storage.GetAverageDeliveryTimes(ctx, from, to, providerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get delivery times: %w", err)
	}
	analytics.AverageDeliveryDays = avgDeliveryTimes

	return analytics, nil
}

// DeliveryAnalytics - структура аналитики доставки
type DeliveryAnalytics struct {
	Period              AnalyticsPeriod      `json:"period"`
	TotalShipments      int                  `json:"total_shipments"`
	TotalCost           float64              `json:"total_cost"`
	AverageCost         float64              `json:"average_cost"`
	StatusBreakdown     map[string]int       `json:"status_breakdown"`
	ProviderBreakdown   []ProviderStatistics `json:"provider_breakdown"`
	TopRoutes           []RouteStatistics    `json:"top_routes"`
	AverageDeliveryDays map[string]float64   `json:"average_delivery_days"`
}

// AnalyticsPeriod - период аналитики
type AnalyticsPeriod struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

// ProviderStatistics - статистика по провайдеру
type ProviderStatistics struct {
	ProviderID    int     `json:"provider_id"`
	ProviderName  string  `json:"provider_name"`
	ShipmentCount int     `json:"shipment_count"`
	TotalCost     float64 `json:"total_cost"`
	AverageCost   float64 `json:"average_cost"`
	SuccessRate   float64 `json:"success_rate"`
}

// RouteStatistics - статистика маршрута
type RouteStatistics struct {
	FromCity      string  `json:"from_city"`
	ToCity        string  `json:"to_city"`
	ShipmentCount int     `json:"shipment_count"`
	AverageCost   float64 `json:"average_cost"`
	AverageDays   float64 `json:"average_days"`
}
