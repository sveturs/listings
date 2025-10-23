package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"

	"backend/internal/proj/delivery/attributes"
	"backend/internal/proj/delivery/grpcclient"
	"backend/internal/proj/delivery/models"
	"backend/internal/proj/delivery/notifications"
	"backend/internal/proj/delivery/storage"
	notifService "backend/internal/proj/notifications/service"
	pb "backend/pkg/grpc/delivery/v1"
)

// Service - основной сервис управления доставкой
type Service struct {
	db            *sqlx.DB
	storage       *storage.Storage
	attributes    *attributes.Service
	notifications *notifications.DeliveryNotificationIntegration
	grpcClient    grpcClientInterface // gRPC клиент для делегирования к микросервису (REQUIRED)
}

// grpcClientInterface определяет интерфейс для gRPC клиента
type grpcClientInterface interface {
	CreateShipment(ctx context.Context, req *pb.CreateShipmentRequest) (*pb.CreateShipmentResponse, error)
	GetShipment(ctx context.Context, req *pb.GetShipmentRequest) (*pb.GetShipmentResponse, error)
	TrackShipment(ctx context.Context, req *pb.TrackShipmentRequest) (*pb.TrackShipmentResponse, error)
	CancelShipment(ctx context.Context, req *pb.CancelShipmentRequest) (*pb.CancelShipmentResponse, error)
	CalculateRate(ctx context.Context, req *pb.CalculateRateRequest) (*pb.CalculateRateResponse, error)
	GetSettlements(ctx context.Context, req *pb.GetSettlementsRequest) (*pb.GetSettlementsResponse, error)
	GetStreets(ctx context.Context, req *pb.GetStreetsRequest) (*pb.GetStreetsResponse, error)
	GetParcelLockers(ctx context.Context, req *pb.GetParcelLockersRequest) (*pb.GetParcelLockersResponse, error)
	Close() error
}

// NewService создает новый экземпляр сервиса доставки
// ВАЖНО: grpcClient ОБЯЗАТЕЛЕН, без него сервис не будет работать
func NewService(db *sqlx.DB, grpcClient grpcClientInterface) *Service {
	if grpcClient == nil {
		panic("delivery service requires gRPC client, but got nil")
	}

	return &Service{
		db:         db,
		storage:    storage.NewStorage(db),
		attributes: attributes.NewService(db),
		grpcClient: grpcClient,
	}
}

// SetNotificationService устанавливает сервис уведомлений
func (s *Service) SetNotificationService(notifService notifService.NotificationServiceInterface) {
	if notifService != nil {
		s.notifications = notifications.NewDeliveryNotificationIntegration(notifService)
		log.Info().Msg("Notification service integrated with delivery module")
	}
}

// GetGRPCClient возвращает gRPC клиент для прямых вызовов (используется в test handlers)
func (s *Service) GetGRPCClient() grpcClientInterface {
	return s.grpcClient
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

// CreateShipment - создает отправление через gRPC микросервис
func (s *Service) CreateShipment(ctx context.Context, req *CreateShipmentRequest) (*models.Shipment, error) {
	if s.grpcClient == nil {
		return nil, fmt.Errorf("delivery service not configured: gRPC client is nil")
	}

	shipment, err := s.createShipmentViaGRPC(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create shipment via gRPC: %w", err)
	}

	log.Info().Int("shipment_id", shipment.ID).Msg("Shipment created via gRPC successfully")
	return shipment, nil
}

// createShipmentViaGRPC создает отправление через gRPC микросервис
func (s *Service) createShipmentViaGRPC(ctx context.Context, req *CreateShipmentRequest) (*models.Shipment, error) {
	// Конвертируем адреса
	fromAddr := &pb.Address{
		Street:       req.FromAddress.Street,
		City:         req.FromAddress.City,
		State:        "",
		PostalCode:   req.FromAddress.PostalCode,
		Country:      req.FromAddress.Country,
		ContactName:  req.FromAddress.Name,
		ContactPhone: req.FromAddress.Phone,
	}

	toAddr := &pb.Address{
		Street:       req.ToAddress.Street,
		City:         req.ToAddress.City,
		State:        "",
		PostalCode:   req.ToAddress.PostalCode,
		Country:      req.ToAddress.Country,
		ContactName:  req.ToAddress.Name,
		ContactPhone: req.ToAddress.Phone,
	}

	// Конвертируем первую посылку (новая схема поддерживает только одну)
	var pbPackage *pb.Package
	if len(req.Packages) > 0 {
		pkg := req.Packages[0]
		pbPackage = &pb.Package{
			Weight:        fmt.Sprintf("%.2f", pkg.Weight),
			Length:        fmt.Sprintf("%.2f", pkg.Dimensions.Length),
			Width:         fmt.Sprintf("%.2f", pkg.Dimensions.Width),
			Height:        fmt.Sprintf("%.2f", pkg.Dimensions.Height),
			Description:   pkg.Description,
			DeclaredValue: fmt.Sprintf("%.2f", pkg.Value),
		}
	} else {
		return nil, fmt.Errorf("at least one package is required")
	}

	// Создаем gRPC запрос (упрощенная схема)
	grpcReq := &pb.CreateShipmentRequest{
		Provider:    grpcclient.MapProviderCodeToEnum(req.ProviderCode),
		FromAddress: fromAddr,
		ToAddress:   toAddr,
		Package:     pbPackage,
		UserId:      fmt.Sprintf("%d", req.OrderID), // Используем OrderID как временный UserId
	}

	// Вызываем gRPC
	grpcResp, err := s.grpcClient.CreateShipment(ctx, grpcReq)
	if err != nil {
		return nil, fmt.Errorf("gRPC call failed: %w", err)
	}

	// Конвертируем response в модель БД
	shipment, err := grpcclient.MapShipmentFromProto(grpcResp.Shipment)
	if err != nil {
		return nil, fmt.Errorf("failed to map shipment from proto: %w", err)
	}

	// Устанавливаем поля, которые mapper не заполняет
	shipment.ProviderID = req.ProviderID
	shipment.OrderID = &req.OrderID

	// Сохраняем в локальную БД для кеширования
	if err := s.storage.CreateShipment(ctx, shipment); err != nil {
		log.Error().Err(err).Msg("Failed to cache shipment in local DB")
		// Не возвращаем ошибку, т.к. отправление создано успешно
	}

	// Обновляем заказ
	if req.OrderID > 0 {
		if err := s.storage.UpdateOrderShipment(ctx, req.OrderID, shipment.ID); err != nil {
			log.Error().Err(err).Int("order_id", req.OrderID).Msg("Failed to update order shipment")
		}
	}

	return shipment, nil
}

// TrackShipment - отслеживает отправление через gRPC микросервис
func (s *Service) TrackShipment(ctx context.Context, trackingNumber string) (*TrackingInfo, error) {
	if s.grpcClient == nil {
		return nil, fmt.Errorf("delivery service not configured: gRPC client is nil")
	}

	trackingInfo, err := s.trackShipmentViaGRPC(ctx, trackingNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to track shipment via gRPC: %w", err)
	}

	log.Info().Str("tracking_number", trackingNumber).Msg("Tracking via gRPC successful")
	return trackingInfo, nil
}

// trackShipmentViaGRPC отслеживает через gRPC
func (s *Service) trackShipmentViaGRPC(ctx context.Context, trackingNumber string) (*TrackingInfo, error) {
	// Получаем отправление из локальной БД для получения shipment_id
	shipment, err := s.storage.GetShipmentByTracking(ctx, trackingNumber)
	if err != nil {
		return nil, fmt.Errorf("shipment not found in local DB: %w", err)
	}

	// Вызываем gRPC
	grpcResp, err := s.grpcClient.TrackShipment(ctx, &pb.TrackShipmentRequest{
		TrackingNumber: trackingNumber,
	})
	if err != nil {
		return nil, fmt.Errorf("gRPC tracking call failed: %w", err)
	}

	// Конвертируем события
	events := grpcclient.MapTrackingEventsFromProto(grpcResp.Events)

	// Обновляем статус в локальной БД (берем из grpcResp.Shipment)
	if grpcResp.Shipment != nil {
		newStatus := grpcclient.MapStatusFromProto(grpcResp.Shipment.Status)
		if shipment.Status != newStatus {
			oldStatus := shipment.Status
			shipment.Status = newStatus

			var deliveredAt *time.Time
			if grpcResp.Shipment.ActualDelivery != nil {
				t := grpcResp.Shipment.ActualDelivery.AsTime()
				deliveredAt = &t
				shipment.ActualDeliveryDate = deliveredAt
			}

			if err := s.storage.UpdateShipmentStatus(ctx, shipment.ID, newStatus, deliveredAt); err != nil {
				log.Error().Err(err).Int("shipment_id", shipment.ID).Msg("Failed to update shipment status")
			}

			// Отправляем уведомление
			// Извлекаем location и description из последнего события
			var currentLocation, statusText string
			if len(grpcResp.Events) > 0 {
				lastEvent := grpcResp.Events[len(grpcResp.Events)-1]
				currentLocation = lastEvent.Location
				statusText = lastEvent.Description
			}
			s.sendStatusNotification(ctx, shipment, oldStatus, newStatus, currentLocation, statusText)
		}
	}

	// Сохраняем новые события в локальную БД
	for i := range events {
		events[i].ShipmentID = shipment.ID
		events[i].ProviderID = shipment.ProviderID
		if err := s.storage.CreateTrackingEvent(ctx, &events[i]); err != nil {
			log.Error().Err(err).Int("shipment_id", shipment.ID).Msg("Failed to cache tracking event")
		}
	}

	// Формируем ответ
	trackingInfo := &TrackingInfo{
		ShipmentID:     shipment.ID,
		TrackingNumber: trackingNumber,
		Status:         shipment.Status,
		LastUpdated:    shipment.UpdatedAt,
	}

	// Заполняем данные из grpcResp.Shipment
	if grpcResp.Shipment != nil {
		if grpcResp.Shipment.EstimatedDelivery != nil {
			t := grpcResp.Shipment.EstimatedDelivery.AsTime()
			trackingInfo.EstimatedDate = &t
		}

		if grpcResp.Shipment.ActualDelivery != nil {
			t := grpcResp.Shipment.ActualDelivery.AsTime()
			trackingInfo.DeliveredDate = &t
		}
	}

	// Извлекаем location и description из последнего события
	if len(grpcResp.Events) > 0 {
		lastEvent := grpcResp.Events[len(grpcResp.Events)-1]
		trackingInfo.CurrentLocation = lastEvent.Location
		trackingInfo.StatusText = lastEvent.Description
	}

	// Конвертируем события в TrackingEvent
	for _, e := range events {
		trackingInfo.Events = append(trackingInfo.Events, TrackingEvent{
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
		})
	}

	return trackingInfo, nil
}

// CancelShipment - отменяет отправление через gRPC микросервис
func (s *Service) CancelShipment(ctx context.Context, shipmentID int, reason string) error {
	if s.grpcClient == nil {
		return fmt.Errorf("delivery service not configured: gRPC client is nil")
	}

	err := s.cancelShipmentViaGRPC(ctx, shipmentID, reason)
	if err != nil {
		return fmt.Errorf("failed to cancel shipment via gRPC: %w", err)
	}

	log.Info().Int("shipment_id", shipmentID).Msg("Cancellation via gRPC successful")
	return nil
}

// cancelShipmentViaGRPC отменяет через gRPC
func (s *Service) cancelShipmentViaGRPC(ctx context.Context, shipmentID int, reason string) error {
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

	// Вызываем gRPC
	if shipment.ExternalID == nil {
		return fmt.Errorf("external ID not found")
	}

	_, err = s.grpcClient.CancelShipment(ctx, &pb.CancelShipmentRequest{
		Id:     *shipment.ExternalID,
		Reason: reason,
	})
	if err != nil {
		return fmt.Errorf("gRPC cancellation call failed: %w", err)
	}

	// Обновляем в локальной БД
	return s.storage.UpdateShipmentStatus(ctx, shipmentID, models.ShipmentStatusCancelled, nil)
}

// GetProviders - получает список доступных провайдеров
func (s *Service) GetProviders(ctx context.Context, activeOnly bool) ([]models.Provider, error) {
	return s.storage.GetProviders(ctx, activeOnly)
}

// ApplyCategoryDefaults - применяет дефолтные атрибуты к товарам категории
func (s *Service) ApplyCategoryDefaults(ctx context.Context, categoryID int) (int, error) {
	return s.attributes.ApplyCategoryDefaultsToProducts(ctx, categoryID)
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

// Address представляет адрес доставки
type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
}

// Dimensions представляет размеры посылки
type Dimensions struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// Package представляет посылку
type Package struct {
	Weight      float64    `json:"weight"`
	Dimensions  Dimensions `json:"dimensions"`
	Description string     `json:"description"`
	Value       float64    `json:"value"`
}

// CreateShipmentRequest - запрос создания отправления
type CreateShipmentRequest struct {
	ProviderID     int        `json:"provider_id"`
	ProviderCode   string     `json:"provider_code"`
	OrderID        int        `json:"order_id"`
	FromAddress    *Address   `json:"from_address"`
	ToAddress      *Address   `json:"to_address"`
	Packages       []Package  `json:"packages"`
	DeliveryType   string     `json:"delivery_type"`
	PickupDate     *time.Time `json:"pickup_date,omitempty"`
	InsuranceValue float64    `json:"insurance_value,omitempty"`
	CODAmount      float64    `json:"cod_amount,omitempty"`
	Services       []string   `json:"services,omitempty"`
	Reference      string     `json:"reference,omitempty"`
	Notes          string     `json:"notes,omitempty"`
}

// TrackingEvent представляет событие отслеживания
type TrackingEvent struct {
	Timestamp   time.Time `json:"timestamp"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
}

// ProofOfDelivery представляет подтверждение доставки
type ProofOfDelivery struct {
	SignatureURL string     `json:"signature_url,omitempty"`
	PhotoURL     string     `json:"photo_url,omitempty"`
	ReceivedBy   string     `json:"received_by,omitempty"`
	ReceivedAt   *time.Time `json:"received_at,omitempty"`
	Notes        string     `json:"notes,omitempty"`
}

// TrackingInfo - информация об отслеживании
type TrackingInfo struct {
	ShipmentID      int              `json:"shipment_id"`
	TrackingNumber  string           `json:"tracking_number"`
	Status          string           `json:"status"`
	StatusText      string           `json:"status_text"`
	CurrentLocation string           `json:"current_location,omitempty"`
	EstimatedDate   *time.Time       `json:"estimated_date,omitempty"`
	DeliveredDate   *time.Time       `json:"delivered_date,omitempty"`
	Events          []TrackingEvent  `json:"events"`
	ProofOfDelivery *ProofOfDelivery `json:"proof_of_delivery,omitempty"`
	LastUpdated     time.Time        `json:"last_updated"`
}

// WebhookResponse представляет ответ на webhook
type WebhookResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// GetShipment - получает отправление по ID через gRPC микросервис
func (s *Service) GetShipment(ctx context.Context, shipmentID int) (*models.Shipment, error) {
	if s.grpcClient == nil {
		return nil, fmt.Errorf("delivery service not configured: gRPC client is nil")
	}

	shipment, err := s.getShipmentViaGRPC(ctx, shipmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipment via gRPC: %w", err)
	}

	log.Info().Int("shipment_id", shipmentID).Msg("Get shipment via gRPC successful")
	return shipment, nil
}

// getShipmentViaGRPC получает отправление через gRPC
func (s *Service) getShipmentViaGRPC(ctx context.Context, shipmentID int) (*models.Shipment, error) {
	// Получаем из локальной БД для получения external_id
	shipment, err := s.storage.GetShipment(ctx, shipmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipment from local DB: %w", err)
	}

	if shipment.ExternalID == nil {
		return nil, fmt.Errorf("external ID not found")
	}

	// Вызываем gRPC
	grpcResp, err := s.grpcClient.GetShipment(ctx, &pb.GetShipmentRequest{
		Id: *shipment.ExternalID,
	})
	if err != nil {
		return nil, fmt.Errorf("gRPC get shipment call failed: %w", err)
	}

	// Конвертируем response
	grpcShipment, err := grpcclient.MapShipmentFromProto(grpcResp.Shipment)
	if err != nil {
		return nil, fmt.Errorf("failed to map shipment from proto: %w", err)
	}

	// Обновляем локальную БД
	grpcShipment.ID = shipment.ID
	grpcShipment.ProviderID = shipment.ProviderID
	if err := s.storage.UpdateShipmentStatus(ctx, shipment.ID, grpcShipment.Status, grpcShipment.ActualDeliveryDate); err != nil {
		log.Error().Err(err).Int("shipment_id", shipment.ID).Msg("Failed to update local cache")
	}

	return grpcShipment, nil
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
