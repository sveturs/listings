// backend/internal/proj/delivery/handler/test_handler.go
package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"

	pb "backend/pkg/grpc/delivery/v1"
	"backend/pkg/utils"
)

// TestShipmentRequest - запрос для создания тестового отправления через delivery микросервис
type TestShipmentRequest struct {
	// Получатель
	RecipientName    string `json:"recipient_name" example:"Petar Petrović"`
	RecipientPhone   string `json:"recipient_phone" example:"0641234567"`
	RecipientEmail   string `json:"recipient_email" example:"petar@example.com"`
	RecipientCity    string `json:"recipient_city" example:"Beograd"`
	RecipientAddress string `json:"recipient_address" example:"Takovska 2"`
	RecipientZip     string `json:"recipient_zip" example:"11000"`

	// Отправитель
	SenderName    string `json:"sender_name" example:"Sve Tu d.o.o."`
	SenderPhone   string `json:"sender_phone" example:"0641234567"`
	SenderEmail   string `json:"sender_email" example:"b2b@svetu.rs"`
	SenderCity    string `json:"sender_city" example:"Beograd"`
	SenderAddress string `json:"sender_address" example:"Bulevar kralja Aleksandra 73"`
	SenderZip     string `json:"sender_zip" example:"11000"`

	// Параметры отправления
	Weight       int    `json:"weight" example:"500"` // граммы
	Content      string `json:"content" example:"Test paket za SVETU"`
	CODAmount    int64  `json:"cod_amount" example:"0"`    // наложенный платеж (RSD)
	InsuredValue int64  `json:"insured_value" example:"0"` // объявленная ценность (RSD)

	// Дополнительные услуги (для совместимости со старым API)
	Services         string `json:"services" example:"PNA"`                  // PNA, SMS, OTK, VD
	DeliveryMethod   string `json:"delivery_method" example:"K"`             // K = Kurir, S = Šalter, PAK = Pаккетомат
	DeliveryType     string `json:"delivery_type" example:"standard"`        // standard, cod, parcel_locker
	PaymentMethod    string `json:"payment_method" example:"POF"`            // POF = gotovina
	IdRukovanje      int    `json:"id_rukovanje" example:"29"`               // ID услуги доставки (29, 30, 55, 58, 59, 71, 85)
	ParcelLockerCode string `json:"parcel_locker_code,omitempty" example:""` // Код пакketomата (для IdRukovanje = 85)
}

// TestShipmentResponse - ответ с результатом создания
type TestShipmentResponse struct {
	Success        bool        `json:"success"`
	TrackingNumber string      `json:"tracking_number,omitempty"`
	ShipmentID     string      `json:"shipment_id,omitempty"`
	Cost           string      `json:"cost,omitempty"`
	Currency       string      `json:"currency,omitempty"`
	Errors         []string    `json:"errors,omitempty"`
	RequestData    interface{} `json:"request_data,omitempty"`
	ResponseData   interface{} `json:"response_data,omitempty"`
	CreatedAt      string      `json:"created_at"`
	ProcessingTime int64       `json:"processing_time_ms"` // milliseconds
}

// CalculateTestRequest - запрос расчета стоимости
type CalculateTestRequest struct {
	FromCity  string `json:"from_city" example:"Beograd"`
	ToCity    string `json:"to_city" example:"Novi Sad"`
	Weight    int    `json:"weight" example:"1000"` // граммы
	CODAmount int64  `json:"cod_amount" example:"0"`
}

// CreateTestShipment создает тестовое отправление через delivery gRPC микросервис
// @Summary Create test shipment via gRPC microservice
// @Description Create a test shipment using delivery gRPC microservice (supports all 5 providers)
// @Tags delivery-test
// @Accept json
// @Produce json
// @Param request body TestShipmentRequest true "Test shipment request"
// @Success 200 {object} utils.SuccessResponseSwag{data=TestShipmentResponse} "Test shipment result"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/public/delivery/test/shipment [post]
func (h *Handler) CreateTestShipment(c *fiber.Ctx) error {
	startTime := time.Now()

	var req TestShipmentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "delivery.invalid_request", fiber.Map{
			"error": err.Error(),
		})
	}

	// Валидация обязательных полей
	if req.RecipientName == "" || req.RecipientPhone == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "delivery.missing_recipient_data", fiber.Map{
			"missing_fields": []string{"recipient_name", "recipient_phone"},
		})
	}

	if req.SenderName == "" || req.SenderPhone == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "delivery.missing_sender_data", fiber.Map{
			"missing_fields": []string{"sender_name", "sender_phone"},
		})
	}

	// Преобразуем weight из граммов в кг для gRPC
	weightKg := float32(req.Weight) / 1000.0

	// Создаем gRPC запрос
	grpcReq := &pb.CreateShipmentRequest{
		Provider: pb.DeliveryProvider_DELIVERY_PROVIDER_POST_EXPRESS, // По умолчанию Post Express
		FromAddress: &pb.Address{
			ContactName:  req.SenderName,
			ContactPhone: req.SenderPhone,
			Street:       req.SenderAddress,
			City:         req.SenderCity,
			PostalCode:   req.SenderZip,
			Country:      "RS",
		},
		ToAddress: &pb.Address{
			ContactName:  req.RecipientName,
			ContactPhone: req.RecipientPhone,
			Street:       req.RecipientAddress,
			City:         req.RecipientCity,
			PostalCode:   req.RecipientZip,
			Country:      "RS",
		},
		Package: &pb.Package{
			Weight:        fmt.Sprintf("%.3f", weightKg),
			Description:   req.Content,
			DeclaredValue: fmt.Sprintf("%d", req.InsuredValue),
		},
		UserId: "test-user", // Для тестовых отправлений
	}

	// Вызываем gRPC микросервис
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := h.service.GetGRPCClient().CreateShipment(ctx, grpcReq)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("recipient_city", req.RecipientCity).
			Msg("Failed to create test shipment via gRPC")
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "delivery.shipment_creation_failed", fiber.Map{
			"error": err.Error(),
		})
	}

	processingTime := time.Since(startTime).Milliseconds()

	// Формируем ответ
	result := TestShipmentResponse{
		Success:        true,
		TrackingNumber: resp.Shipment.TrackingNumber,
		ShipmentID:     resp.Shipment.Id,
		Cost:           resp.Shipment.Cost,
		Currency:       resp.Shipment.Currency,
		CreatedAt:      time.Now().Format(time.RFC3339),
		ProcessingTime: processingTime,
		RequestData:    req,
		ResponseData:   resp.Shipment,
	}

	h.logger.Info().
		Str("tracking_number", result.TrackingNumber).
		Str("shipment_id", result.ShipmentID).
		Int64("processing_time_ms", processingTime).
		Msg("Test shipment created successfully via gRPC")

	return utils.SendSuccessResponse(c, result, "Test shipment created successfully")
}

// TrackTestShipment отслеживает тестовое отправление через gRPC
// @Summary Track test shipment
// @Description Track a test shipment using delivery gRPC microservice
// @Tags delivery-test
// @Accept json
// @Produce json
// @Param tracking_number path string true "Tracking number"
// @Success 200 {object} utils.SuccessResponseSwag "Tracking information"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid tracking number"
// @Failure 404 {object} utils.ErrorResponseSwag "Shipment not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/public/delivery/test/tracking/{tracking_number} [get]
func (h *Handler) TrackTestShipment(c *fiber.Ctx) error {
	trackingNumber := c.Params("tracking_number")
	if trackingNumber == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "delivery.missing_tracking_number", nil)
	}

	// Вызываем gRPC микросервис
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	grpcReq := &pb.TrackShipmentRequest{
		TrackingNumber: trackingNumber,
	}

	resp, err := h.service.GetGRPCClient().TrackShipment(ctx, grpcReq)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("tracking_number", trackingNumber).
			Msg("Failed to track test shipment via gRPC")
		return utils.SendErrorResponse(c, fiber.StatusNotFound, "delivery.shipment_not_found", fiber.Map{
			"tracking_number": trackingNumber,
			"error":           err.Error(),
		})
	}

	return utils.SendSuccessResponse(c, fiber.Map{
		"shipment": resp.Shipment,
		"events":   resp.Events,
	}, "Shipment tracking information retrieved successfully")
}

// CancelTestShipment отменяет тестовое отправление через gRPC
// @Summary Cancel test shipment
// @Description Cancel a test shipment using delivery gRPC microservice
// @Tags delivery-test
// @Accept json
// @Produce json
// @Param id path string true "Shipment ID (UUID)"
// @Success 200 {object} utils.SuccessResponseSwag "Shipment canceled"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid shipment ID"
// @Failure 404 {object} utils.ErrorResponseSwag "Shipment not found"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/public/delivery/test/cancel/{id} [post]
func (h *Handler) CancelTestShipment(c *fiber.Ctx) error {
	shipmentID := c.Params("id")
	if shipmentID == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "delivery.missing_shipment_id", nil)
	}

	// Вызываем gRPC микросервис
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	grpcReq := &pb.CancelShipmentRequest{
		Id:     shipmentID,
		Reason: "Test cancellation via delivery API",
	}

	resp, err := h.service.GetGRPCClient().CancelShipment(ctx, grpcReq)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("shipment_id", shipmentID).
			Msg("Failed to cancel test shipment via gRPC")
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "delivery.shipment_cancellation_failed", fiber.Map{
			"shipment_id": shipmentID,
			"error":       err.Error(),
		})
	}

	h.logger.Info().
		Str("shipment_id", shipmentID).
		Str("tracking_number", resp.Shipment.TrackingNumber).
		Msg("Test shipment canceled successfully via gRPC")

	return utils.SendSuccessResponse(c, fiber.Map{
		"success":  true,
		"shipment": resp.Shipment,
		"message":  "Shipment canceled successfully",
	}, "Test shipment canceled successfully")
}

// CalculateTestRate рассчитывает стоимость доставки через gRPC
// @Summary Calculate test delivery rate
// @Description Calculate delivery cost using delivery gRPC microservice
// @Tags delivery-test
// @Accept json
// @Produce json
// @Param request body CalculateTestRequest true "Calculate request"
// @Success 200 {object} utils.SuccessResponseSwag "Calculation result"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/public/delivery/test/calculate [post]
func (h *Handler) CalculateTestRate(c *fiber.Ctx) error {
	var req CalculateTestRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "delivery.invalid_request", fiber.Map{
			"error": err.Error(),
		})
	}

	// Валидация
	if req.FromCity == "" || req.ToCity == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "delivery.missing_cities", fiber.Map{
			"missing_fields": []string{"from_city", "to_city"},
		})
	}

	if req.Weight <= 0 {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "delivery.invalid_weight", fiber.Map{
			"weight": req.Weight,
		})
	}

	// Преобразуем weight из граммов в кг
	weightKg := float32(req.Weight) / 1000.0

	// Вызываем gRPC микросервис
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	grpcReq := &pb.CalculateRateRequest{
		Provider: pb.DeliveryProvider_DELIVERY_PROVIDER_POST_EXPRESS,
		FromAddress: &pb.Address{
			City:    req.FromCity,
			Country: "RS",
		},
		ToAddress: &pb.Address{
			City:    req.ToCity,
			Country: "RS",
		},
		Package: &pb.Package{
			Weight: fmt.Sprintf("%.3f", weightKg),
		},
	}

	resp, err := h.service.GetGRPCClient().CalculateRate(ctx, grpcReq)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("from_city", req.FromCity).
			Str("to_city", req.ToCity).
			Msg("Failed to calculate test rate via gRPC")
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "delivery.calculation_failed", fiber.Map{
			"error": err.Error(),
		})
	}

	return utils.SendSuccessResponse(c, fiber.Map{
		"cost":               resp.Cost,
		"currency":           resp.Currency,
		"estimated_delivery": resp.EstimatedDelivery,
		"from_city":          req.FromCity,
		"to_city":            req.ToCity,
		"weight_kg":          weightKg,
	}, "Delivery rate calculated successfully")
}

// GetTestSettlements возвращает список населенных пунктов
// @Summary Get settlements list (MOCK)
// @Description Get list of settlements (currently returns mock data, will be implemented in microservice)
// @Tags delivery-test
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag "Settlements list"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/public/delivery/test/settlements [get]
func (h *Handler) GetTestSettlements(c *fiber.Ctx) error {
	// TODO: Добавить RPC метод в микросервис для получения settlements
	// Пока возвращаем mock данные для совместимости
	h.logger.Warn().Msg("GetTestSettlements called - using mock data (implement RPC method in microservice)")

	mockSettlements := []fiber.Map{
		{"id": 1, "name": "Beograd", "zip_code": "11000"},
		{"id": 2, "name": "Novi Sad", "zip_code": "21000"},
		{"id": 3, "name": "Niš", "zip_code": "18000"},
		{"id": 4, "name": "Kragujevac", "zip_code": "34000"},
		{"id": 5, "name": "Subotica", "zip_code": "24000"},
	}

	return utils.SendSuccessResponse(c, fiber.Map{
		"settlements": mockSettlements,
		"note":        "Mock data - implement GetSettlements RPC method in microservice",
	}, "Settlements list retrieved successfully (mock data)")
}

// GetTestStreets возвращает список улиц по городу
// @Summary Get streets by settlement (MOCK)
// @Description Get list of streets for a settlement (currently returns mock data, will be implemented in microservice)
// @Tags delivery-test
// @Produce json
// @Param settlement path string true "Settlement name"
// @Success 200 {object} utils.SuccessResponseSwag "Streets list"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid settlement"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/public/delivery/test/streets/{settlement} [get]
func (h *Handler) GetTestStreets(c *fiber.Ctx) error {
	settlement := c.Params("settlement")
	if settlement == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "delivery.missing_settlement", nil)
	}

	// TODO: Добавить RPC метод в микросервис для получения streets
	h.logger.Warn().
		Str("settlement", settlement).
		Msg("GetTestStreets called - using mock data (implement RPC method in microservice)")

	mockStreets := []fiber.Map{
		{"id": 1, "name": "Kneza Miloša", "settlement": settlement},
		{"id": 2, "name": "Bulevar kralja Aleksandra", "settlement": settlement},
		{"id": 3, "name": "Terazije", "settlement": settlement},
	}

	return utils.SendSuccessResponse(c, fiber.Map{
		"streets":    mockStreets,
		"settlement": settlement,
		"note":       "Mock data - implement GetStreets RPC method in microservice",
	}, "Streets list retrieved successfully (mock data)")
}

// GetTestParcelLockers возвращает список паккетоматов
// @Summary Get parcel lockers list (MOCK)
// @Description Get list of parcel lockers (currently returns mock data, will be implemented in microservice)
// @Tags delivery-test
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag "Parcel lockers list"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/public/delivery/test/parcel-lockers [get]
func (h *Handler) GetTestParcelLockers(c *fiber.Ctx) error {
	// TODO: Добавить RPC метод в микросервис для получения parcel lockers
	h.logger.Warn().Msg("GetTestParcelLockers called - using mock data (implement RPC method in microservice)")

	mockLockers := []fiber.Map{
		{"id": 1, "code": "BG001", "name": "Beograd - Terazije", "address": "Terazije 1, Beograd"},
		{"id": 2, "code": "BG002", "name": "Beograd - Savski venac", "address": "Kneza Miloša 10, Beograd"},
		{"id": 3, "code": "NS001", "name": "Novi Sad - Centar", "address": "Bulevar Oslobodjenja 1, Novi Sad"},
	}

	return utils.SendSuccessResponse(c, fiber.Map{
		"parcel_lockers": mockLockers,
		"note":           "Mock data - implement GetParcelLockers RPC method in microservice",
	}, "Parcel lockers list retrieved successfully (mock data)")
}

// GetTestDeliveryServices возвращает список услуг доставки
// @Summary Get delivery services list (MOCK)
// @Description Get list of delivery services (currently returns mock data, will be implemented in microservice)
// @Tags delivery-test
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag "Delivery services list"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/public/delivery/test/delivery-services [get]
func (h *Handler) GetTestDeliveryServices(c *fiber.Ctx) error {
	// TODO: Добавить RPC метод в микросервис для получения delivery services
	h.logger.Warn().Msg("GetTestDeliveryServices called - using mock data (implement RPC method in microservice)")

	mockServices := []fiber.Map{
		{"id": 29, "name": "Kurirska dostava - standardna", "code": "KURIR_STD"},
		{"id": 30, "name": "Kurirska dostava - ekspress", "code": "KURIR_EXP"},
		{"id": 55, "name": "Šalterska dostava", "code": "SALTER"},
		{"id": 85, "name": "Pакетомат", "code": "PARCEL_LOCKER"},
	}

	return utils.SendSuccessResponse(c, fiber.Map{
		"delivery_services": mockServices,
		"note":              "Mock data - implement GetDeliveryServices RPC method in microservice",
	}, "Delivery services list retrieved successfully (mock data)")
}

// ValidateTestAddress валидирует адрес через gRPC (MOCK)
// @Summary Validate address (MOCK)
// @Description Validate address (currently returns mock data, will be implemented in microservice)
// @Tags delivery-test
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag "Address validation result"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/public/delivery/test/validate-address [post]
func (h *Handler) ValidateTestAddress(c *fiber.Ctx) error {
	var req struct {
		City    string `json:"city"`
		Street  string `json:"street"`
		ZipCode string `json:"zip_code"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "delivery.invalid_request", fiber.Map{
			"error": err.Error(),
		})
	}

	// TODO: Добавить RPC метод в микросервис для валидации адреса
	h.logger.Warn().
		Str("city", req.City).
		Str("street", req.Street).
		Msg("ValidateTestAddress called - using mock validation (implement RPC method in microservice)")

	// Mock validation - считаем адрес валидным если указан город
	isValid := req.City != ""

	return utils.SendSuccessResponse(c, fiber.Map{
		"valid":   isValid,
		"city":    req.City,
		"street":  req.Street,
		"zipcode": req.ZipCode,
		"note":    "Mock validation - implement ValidateAddress RPC method in microservice",
	}, "Address validation completed successfully (mock data)")
}

// GetTestProviders возвращает список провайдеров доставки
// @Summary Get delivery providers list (MOCK)
// @Description Get list of available delivery providers (currently returns mock data, will be implemented in microservice)
// @Tags delivery-test
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag "Providers list"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/public/delivery/test/providers [get]
func (h *Handler) GetTestProviders(c *fiber.Ctx) error {
	// TODO: Добавить RPC метод в микросервис для получения providers
	h.logger.Warn().Msg("GetTestProviders called - using mock data (implement RPC method in microservice)")

	mockProviders := []fiber.Map{
		{"code": "post_express", "name": "Post Express", "enabled": true, "supports_cod": true, "supports_tracking": true},
		{"code": "bex_express", "name": "BEX Express", "enabled": true, "supports_cod": true, "supports_tracking": true},
		{"code": "dhl", "name": "DHL", "enabled": false, "supports_cod": false, "supports_tracking": true},
		{"code": "ups", "name": "UPS", "enabled": false, "supports_cod": false, "supports_tracking": true},
		{"code": "fedex", "name": "FedEx", "enabled": false, "supports_cod": false, "supports_tracking": true},
	}

	return utils.SendSuccessResponse(c, fiber.Map{
		"providers": mockProviders,
		"note":      "Mock data - implement GetProviders RPC method in microservice",
	}, "Providers list retrieved successfully (mock data)")
}

// GetTestConfig возвращает конфигурацию доставки
// @Summary Get delivery config (MOCK)
// @Description Get delivery configuration (currently returns mock data, will be implemented in microservice)
// @Tags delivery-test
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag "Delivery config"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/public/delivery/test/config [get]
func (h *Handler) GetTestConfig(c *fiber.Ctx) error {
	// TODO: Добавить RPC метод в микросервис для получения config
	h.logger.Warn().Msg("GetTestConfig called - using mock data (implement RPC method in microservice)")

	mockConfig := fiber.Map{
		"default_provider":    "post_express",
		"supported_countries": []string{"RS", "BA", "HR", "ME", "MK"},
		"currency":            "RSD",
		"max_weight_kg":       30,
		"min_weight_kg":       0.1,
		"max_declared_value":  1000000,
		"tracking_enabled":    true,
		"cod_enabled":         true,
		"insurance_enabled":   true,
	}

	return utils.SendSuccessResponse(c, fiber.Map{
		"config": mockConfig,
		"note":   "Mock data - implement GetConfig RPC method in microservice",
	}, "Delivery config retrieved successfully (mock data)")
}

// GetTestHistory возвращает историю отправлений
// @Summary Get shipments history (MOCK)
// @Description Get history of shipments (currently returns mock data, will be implemented in microservice)
// @Tags delivery-test
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} utils.SuccessResponseSwag "Shipments history"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/public/delivery/test/history [get]
func (h *Handler) GetTestHistory(c *fiber.Ctx) error {
	// TODO: Добавить RPC метод в микросервис для получения history
	h.logger.Warn().Msg("GetTestHistory called - using mock data (implement RPC method in microservice)")

	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	mockShipments := []fiber.Map{
		{
			"id":              "test-shipment-1",
			"tracking_number": "PE20250101001",
			"status":          "delivered",
			"provider":        "post_express",
			"created_at":      "2025-01-15T10:30:00Z",
			"delivered_at":    "2025-01-17T14:20:00Z",
			"cost":            "450.00",
			"currency":        "RSD",
		},
		{
			"id":              "test-shipment-2",
			"tracking_number": "PE20250101002",
			"status":          "in_transit",
			"provider":        "post_express",
			"created_at":      "2025-01-20T09:15:00Z",
			"cost":            "320.00",
			"currency":        "RSD",
		},
	}

	return utils.SendSuccessResponse(c, fiber.Map{
		"shipments": mockShipments,
		"limit":     limit,
		"offset":    offset,
		"total":     2,
		"note":      "Mock data - implement GetHistory RPC method in microservice",
	}, "Shipments history retrieved successfully (mock data)")
}
