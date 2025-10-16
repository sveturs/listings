package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"

	"backend/internal/proj/postexpress"
	"backend/internal/proj/postexpress/service"
	"backend/pkg/utils"
)

// TestShipmentRequest - запрос для создания тестового отправления
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

	// Дополнительные услуги
	Services         string `json:"services" example:"PNA"`                  // PNA, SMS, OTK, VD
	DeliveryMethod   string `json:"delivery_method" example:"K"`             // K = Kurir, S = Šalter, PAK = Paкетомат
	DeliveryType     string `json:"delivery_type" example:"standard"`        // standard, cod, parcel_locker
	PaymentMethod    string `json:"payment_method" example:"POF"`            // POF = gotovina
	IdRukovanje      int    `json:"id_rukovanje" example:"29"`               // ID услуги доставки (29, 30, 55, 58, 59, 71, 85)
	ParcelLockerCode string `json:"parcel_locker_code,omitempty" example:""` // Код паккетомата (для IdRukovanje = 85)
}

// TestShipmentResponse - ответ с результатом создания
type TestShipmentResponse struct {
	Success        bool        `json:"success"`
	TrackingNumber string      `json:"tracking_number,omitempty"`
	ManifestID     int         `json:"manifest_id,omitempty"`
	ShipmentID     int         `json:"shipment_id,omitempty"`
	ExternalID     string      `json:"external_id,omitempty"`
	Cost           int         `json:"cost,omitempty"` // RSD
	Errors         []string    `json:"errors,omitempty"`
	RequestData    interface{} `json:"request_data,omitempty"`
	ResponseData   interface{} `json:"response_data,omitempty"`
	CreatedAt      string      `json:"created_at"`
	ProcessingTime int64       `json:"processing_time_ms"` // milliseconds
}

// CreateTestShipment создает тестовое отправление через реальный Post Express API
// @Summary Create test shipment
// @Description Create a test shipment using real Post Express WSP API (B2B Manifest)
// @Tags post-express-test
// @Accept json
// @Produce json
// @Param request body TestShipmentRequest true "Test shipment request"
// @Success 200 {object} utils.SuccessResponseSwag{data=TestShipmentResponse} "Test shipment result"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/v1/postexpress/test/shipment [post]
func (h *Handler) CreateTestShipment(c *fiber.Ctx) error {
	startTime := time.Now()

	var req TestShipmentRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "postexpress.invalid_request", nil)
	}

	// Валидация обязательных полей
	// Для parcel locker (IdRukovanje=85) не требуется адрес получателя
	if req.IdRukovanje == 85 {
		// Для паккетомата нужны только имя, телефон и код паккетомата
		if req.RecipientName == "" || req.RecipientPhone == "" || req.ParcelLockerCode == "" {
			return utils.SendErrorResponse(c, fiber.StatusBadRequest, "postexpress.missing_parcel_locker_data", fiber.Map{
				"missing_fields": []string{"recipient_name", "recipient_phone", "parcel_locker_code"},
			})
		}
	} else {
		// Для обычной доставки нужны полные данные получателя
		if req.RecipientName == "" || req.RecipientPhone == "" || req.RecipientCity == "" {
			return utils.SendErrorResponse(c, fiber.StatusBadRequest, "postexpress.missing_recipient_data", fiber.Map{
				"missing_fields": []string{"recipient_name", "recipient_phone", "recipient_city"},
			})
		}
	}

	if req.SenderName == "" || req.SenderCity == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "postexpress.missing_sender_data", fiber.Map{
			"missing_fields": []string{"sender_name", "sender_city"},
		})
	}

	if req.Weight <= 0 {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "postexpress.invalid_weight", fiber.Map{
			"weight": req.Weight,
		})
	}

	// Устанавливаем дефолтные значения если не указаны
	if req.Services == "" {
		req.Services = "PNA"
	}
	if req.DeliveryMethod == "" {
		req.DeliveryMethod = "K"
	}
	if req.PaymentMethod == "" {
		req.PaymentMethod = "POF"
	}

	// РЕАЛЬНОЕ создание отправления через WSP API
	result, err := h.createRealShipment(c.Context(), &req, startTime)
	if err != nil {
		h.logger.Error("Failed to create real test shipment: %v", err)
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "postexpress.create_shipment_failed", fiber.Map{
			"error": err.Error(),
		})
	}

	return utils.SendSuccessResponse(c, result, "Test shipment created successfully via real Post Express API")
}

// GetTestConfig возвращает текущую конфигурацию Post Express для тестирования
// @Summary Get Post Express test config
// @Description Get current Post Express configuration (without password) for testing page
// @Tags post-express-test
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "Configuration"
// @Router /api/v1/postexpress/test/config [get]
func (h *Handler) GetTestConfig(c *fiber.Ctx) error {
	config := fiber.Map{
		"api_available": true,
		"test_mode":     true,
		"default_sender": fiber.Map{
			"name":    "Sve Tu d.o.o.",
			"phone":   "0641234567",
			"email":   "b2b@svetu.rs",
			"city":    "Beograd",
			"address": "Bulevar kralja Aleksandra 73",
			"zip":     "11000",
		},
		"default_recipient": fiber.Map{
			"name":    "Petar Petrović",
			"phone":   "0641234567",
			"email":   "petar@example.com",
			"city":    "Beograd",
			"address": "Takovska 2",
			"zip":     "11000",
		},
		"delivery_types": []fiber.Map{
			{"code": "standard", "name": "Обычная доставка", "description": "Стандартная доставка курьером или в отделение"},
			{"code": "cod", "name": "Откупная (COD)", "description": "Доставка с наложенным платежом"},
			{"code": "parcel_locker", "name": "Паккетомат", "description": "Доставка в паккетомат (IdRukovanje: 85)"},
		},
		"delivery_methods": []fiber.Map{
			{"code": "K", "name": "Kurir (Courier)"},
			{"code": "S", "name": "Šalter (Post Office)"},
			{"code": "PAK", "name": "Paкетомат (Parcel Locker)"},
		},
		"id_rukovanje_options": []fiber.Map{
			{"id": 29, "name": "PE_Danas_za_sutra_12", "description": "Доставка завтра до 12:00"},
			{"id": 30, "name": "PE_Danas_za_danas", "description": "Доставка сегодня"},
			{"id": 55, "name": "PE_Danas_za_odmah", "description": "Доставка немедленно"},
			{"id": 58, "name": "PE_Danas_za_sutra_19", "description": "Доставка завтра до 19:00"},
			{"id": 59, "name": "PE_Danas_za_odmah_Bg", "description": "Доставка немедленно (только Белград)"},
			{"id": 71, "name": "PE_Danas_za_sutra_isporuka", "description": "Доставка завтра (стандартная)"},
			{"id": 85, "name": "Isporuka_na_paketomatu", "description": "Доставка в паккетомат"},
		},
		"payment_methods": []fiber.Map{
			{"code": "POF", "name": "Gotovina (Cash)"},
			{"code": "K", "name": "Kartica (Card)"},
		},
		"services": []fiber.Map{
			{"code": "PNA", "name": "Prijem na adresi (Pickup at address)"},
			{"code": "SMS", "name": "SMS obaveštenje"},
			{"code": "OTK", "name": "Otkupnina (COD)"},
			{"code": "VD", "name": "Vrednost (Insured value)"},
		},
	}

	return utils.SendSuccessResponse(c, config, "Post Express test configuration")
}

// GetTestHistory возвращает историю тестовых отправлений (мок)
// @Summary Get test shipments history
// @Description Get history of test shipments (mock data for demo)
// @Tags post-express-test
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag{data=[]TestShipmentResponse} "Test shipments history"
// @Router /api/v1/postexpress/test/history [get]
func (h *Handler) GetTestHistory(c *fiber.Ctx) error {
	// Возвращаем пустую историю пока не будет реального хранилища
	history := []TestShipmentResponse{}

	return utils.SendSuccessResponse(c, history, "Test shipments history")
}

// =============================================================================
// НОВЫЕ ТЕСТОВЫЕ ENDPOINTS ДЛЯ ВСЕХ WSP API ФУНКЦИЙ
// =============================================================================

// TestTrackingRequest - запрос отслеживания
type TestTrackingRequest struct {
	TrackingNumber string `json:"tracking_number" example:"SVETU-TEST-123456"`
}

// TestTrackShipment тестирует отслеживание через WSP API (Transaction 15)
// @Summary Test tracking shipment
// @Description Track shipment using real Post Express WSP API (Transaction 15)
// @Tags post-express-test
// @Accept json
// @Produce json
// @Param request body TestTrackingRequest true "Tracking request"
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "Tracking result"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/v1/postexpress/test/track [post]
func (h *Handler) TestTrackShipment(c *fiber.Ctx) error {
	startTime := time.Now()

	var req TestTrackingRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "postexpress.invalid_request", nil)
	}

	if req.TrackingNumber == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "postexpress.tracking_number_required", nil)
	}

	// Проверяем что WSP клиент инициализирован
	if h.wspClient == nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "postexpress.wsp_client_not_initialized", nil)
	}

	h.logger.Info("Testing tracking for shipment: %s", req.TrackingNumber)

	// Вызываем реальный WSP API (Transaction 15)
	trackingResp, err := h.wspClient.GetShipmentStatus(c.Context(), req.TrackingNumber)
	if err != nil {
		h.logger.Error("GetShipmentStatus failed: %v", err)
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "postexpress.tracking_failed", fiber.Map{
			"error": err.Error(),
		})
	}

	result := fiber.Map{
		"success":         true,
		"tracking_number": req.TrackingNumber,
		"status":          trackingResp.Status,
		"events":          trackingResp.Events,
		"processing_time": time.Since(startTime).Milliseconds(),
		"response_data":   trackingResp,
	}

	h.logger.Info("Tracking successful: status=%s, events=%d", trackingResp.Status, len(trackingResp.Events))

	return utils.SendSuccessResponse(c, result, "Tracking data retrieved successfully")
}

// TestCancelRequest - запрос отмены
type TestCancelRequest struct {
	ShipmentID string `json:"shipment_id" example:"12345"`
	Reason     string `json:"reason" example:"Отмена по требованию клиента"`
}

// TestCancelShipment тестирует отмену через WSP API (Transaction 25)
// @Summary Test cancel shipment
// @Description Cancel shipment using real Post Express WSP API (Transaction 25)
// @Tags post-express-test
// @Accept json
// @Produce json
// @Param request body TestCancelRequest true "Cancel request"
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "Cancel result"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/v1/postexpress/test/cancel [post]
func (h *Handler) TestCancelShipment(c *fiber.Ctx) error {
	startTime := time.Now()

	var req TestCancelRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "postexpress.invalid_request", nil)
	}

	if req.ShipmentID == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "postexpress.shipment_id_required", nil)
	}

	if h.wspClient == nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "postexpress.wsp_client_not_initialized", nil)
	}

	h.logger.Info("Testing cancel for shipment: %s, reason: %s", req.ShipmentID, req.Reason)

	// Вызываем реальный WSP API (Transaction 25)
	err := h.wspClient.CancelShipment(c.Context(), req.ShipmentID)
	if err != nil {
		h.logger.Error("CancelShipment failed: %v", err)
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "postexpress.cancel_failed", fiber.Map{
			"error": err.Error(),
		})
	}

	result := fiber.Map{
		"success":         true,
		"shipment_id":     req.ShipmentID,
		"reason":          req.Reason,
		"processing_time": time.Since(startTime).Milliseconds(),
	}

	h.logger.Info("Cancel successful for shipment: %s", req.ShipmentID)

	return utils.SendSuccessResponse(c, result, "Shipment canceled successfully")
}

// TestPrintLabelRequest - запрос печати этикетки
type TestPrintLabelRequest struct {
	ShipmentID string `json:"shipment_id" example:"12345"`
}

// TestPrintLabel тестирует печать этикетки через WSP API (Transaction 20)
// @Summary Test print label
// @Description Print shipment label using real Post Express WSP API (Transaction 20)
// @Tags post-express-test
// @Accept json
// @Produce json
// @Param request body TestPrintLabelRequest true "Print label request"
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "Label data"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/v1/postexpress/test/label [post]
func (h *Handler) TestPrintLabel(c *fiber.Ctx) error {
	startTime := time.Now()

	var req TestPrintLabelRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "postexpress.invalid_request", nil)
	}

	if req.ShipmentID == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "postexpress.shipment_id_required", nil)
	}

	if h.wspClient == nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "postexpress.wsp_client_not_initialized", nil)
	}

	h.logger.Info("Testing print label for shipment: %s", req.ShipmentID)

	// Вызываем реальный WSP API (Transaction 20)
	labelData, err := h.wspClient.PrintLabel(c.Context(), req.ShipmentID)
	if err != nil {
		h.logger.Error("PrintLabel failed: %v", err)
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "postexpress.print_label_failed", fiber.Map{
			"error": err.Error(),
		})
	}

	result := fiber.Map{
		"success":         true,
		"shipment_id":     req.ShipmentID,
		"label_size":      len(labelData),
		"label_format":    "PDF",
		"processing_time": time.Since(startTime).Milliseconds(),
		// labelData содержит PDF bytes, можно вернуть как base64
		"label_data_preview": fmt.Sprintf("PDF data: %d bytes", len(labelData)),
	}

	h.logger.Info("Label printed successfully for shipment: %s, size: %d bytes", req.ShipmentID, len(labelData))

	return utils.SendSuccessResponse(c, result, "Label retrieved successfully")
}

// TestSearchLocationsRequest - запрос поиска населенных пунктов
type TestSearchLocationsRequest struct {
	Query string `json:"query" example:"Beograd"`
}

// TestSearchLocations тестирует поиск населенных пунктов через WSP API (Transaction 3)
// @Summary Test search locations
// @Description Search locations using real Post Express WSP API (Transaction 3)
// @Tags post-express-test
// @Accept json
// @Produce json
// @Param request body TestSearchLocationsRequest true "Search request"
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "Locations list"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/v1/postexpress/test/locations [post]
func (h *Handler) TestSearchLocations(c *fiber.Ctx) error {
	startTime := time.Now()

	var req TestSearchLocationsRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "postexpress.invalid_request", nil)
	}

	if req.Query == "" {
		req.Query = "Beograd" // Дефолтный поиск
	}

	if h.wspClient == nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "postexpress.wsp_client_not_initialized", nil)
	}

	h.logger.Info("Testing search locations: query=%s", req.Query)

	// Вызываем реальный WSP API (Transaction 3)
	locations, err := h.wspClient.GetLocations(c.Context(), req.Query)
	if err != nil {
		h.logger.Error("GetLocations failed: %v", err)
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "postexpress.search_locations_failed", fiber.Map{
			"error": err.Error(),
		})
	}

	result := fiber.Map{
		"success":         true,
		"query":           req.Query,
		"count":           len(locations),
		"locations":       locations,
		"processing_time": time.Since(startTime).Milliseconds(),
	}

	h.logger.Info("Search locations successful: query=%s, found=%d", req.Query, len(locations))

	return utils.SendSuccessResponse(c, result, fmt.Sprintf("Found %d locations", len(locations)))
}

// TestGetOfficesRequest - запрос получения отделений
type TestGetOfficesRequest struct {
	LocationID int `json:"location_id" example:"1"`
}

// TestGetOffices тестирует получение отделений через WSP API (Transaction 10)
// @Summary Test get offices
// @Description Get offices for location using real Post Express WSP API (Transaction 10)
// @Tags post-express-test
// @Accept json
// @Produce json
// @Param request body TestGetOfficesRequest true "Get offices request"
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "Offices list"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/v1/postexpress/test/offices [post]
func (h *Handler) TestGetOffices(c *fiber.Ctx) error {
	startTime := time.Now()

	var req TestGetOfficesRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "postexpress.invalid_request", nil)
	}

	if req.LocationID <= 0 {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "postexpress.location_id_required", nil)
	}

	if h.wspClient == nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "postexpress.wsp_client_not_initialized", nil)
	}

	h.logger.Info("Testing get offices: location_id=%d", req.LocationID)

	// Вызываем реальный WSP API (Transaction 10)
	offices, err := h.wspClient.GetOffices(c.Context(), req.LocationID)
	if err != nil {
		h.logger.Error("GetOffices failed: %v", err)
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "postexpress.get_offices_failed", fiber.Map{
			"error": err.Error(),
		})
	}

	result := fiber.Map{
		"success":         true,
		"location_id":     req.LocationID,
		"count":           len(offices),
		"offices":         offices,
		"processing_time": time.Since(startTime).Milliseconds(),
	}

	h.logger.Info("Get offices successful: location_id=%d, found=%d", req.LocationID, len(offices))

	return utils.SendSuccessResponse(c, result, fmt.Sprintf("Found %d offices", len(offices)))
}

// createRealShipment создает реальное отправление через WSP API
func (h *Handler) createRealShipment(ctx context.Context, req *TestShipmentRequest, startTime time.Time) (*TestShipmentResponse, error) {
	// Проверяем что WSP клиент инициализирован
	if h.wspClient == nil {
		return nil, fmt.Errorf("WSP client not initialized")
	}

	// Определяем ServiceType на основе IdRukovanje
	serviceType := "PE_Danas_za_sutra_isporuka" // Дефолтный сервис (IdRukovanje=71)
	switch req.IdRukovanje {
	case 29:
		serviceType = "PE_Danas_za_sutra_12"
	case 30:
		serviceType = "PE_Danas_za_danas"
	case 55:
		serviceType = "PE_Danas_za_odmah"
	case 58:
		serviceType = "PE_Danas_za_sutra_19"
	case 59:
		serviceType = "PE_Danas_za_odmah_Bg"
	case 71:
		serviceType = "PE_Danas_za_sutra_isporuka"
	case 85:
		serviceType = "PE_Klasicna" // Parcel Locker
	}

	// Подготавливаем запрос для WSP API
	wspReq := &service.WSPShipmentRequest{
		SenderName:          req.SenderName,
		SenderAddress:       req.SenderAddress,
		SenderCity:          req.SenderCity,
		SenderPostalCode:    req.SenderZip,
		SenderPhone:         req.SenderPhone,
		RecipientName:       req.RecipientName,
		RecipientAddress:    req.RecipientAddress,
		RecipientCity:       req.RecipientCity,
		RecipientPostalCode: req.RecipientZip,
		RecipientPhone:      req.RecipientPhone,
		Weight:              float64(req.Weight) / 1000.0, // Конвертируем граммы в килограммы
		CODAmount:           float64(req.CODAmount),
		InsuranceAmount:     float64(req.InsuredValue),
		ServiceType:         serviceType,
		Content:             req.Content,
		Note:                fmt.Sprintf("Test shipment from SVETU platform - %s", req.DeliveryMethod),
		ParcelLockerCode:    req.ParcelLockerCode, // Код паккетомата для IdRukovanje=85
	}

	// Вызываем реальный API
	h.logger.Info("Creating real test shipment via Post Express WSP API")
	h.logger.Debug("Request data: recipient=%s, city=%s, weight=%dg, COD=%d",
		req.RecipientName, req.RecipientCity, req.Weight, req.CODAmount)

	manifestResp, err := h.wspClient.CreateShipmentViaManifest(ctx, wspReq)
	if err != nil {
		return nil, fmt.Errorf("WSP API call failed: %w", err)
	}

	// Проверяем результат
	if manifestResp.Rezultat != 0 {
		return &TestShipmentResponse{
			Success:        false,
			Errors:         []string{manifestResp.Poruka},
			CreatedAt:      time.Now().Format(time.RFC3339),
			ProcessingTime: time.Since(startTime).Milliseconds(),
			RequestData:    buildRequestData(req),
			ResponseData:   manifestResp,
		}, nil
	}

	// Извлекаем данные посылки из ответа
	if len(manifestResp.Porudzbine) == 0 || len(manifestResp.Porudzbine[0].Posiljke) == 0 {
		return &TestShipmentResponse{
			Success:        false,
			Errors:         []string{"No shipment in manifest response"},
			CreatedAt:      time.Now().Format(time.RFC3339),
			ProcessingTime: time.Since(startTime).Milliseconds(),
			RequestData:    buildRequestData(req),
			ResponseData:   manifestResp,
		}, nil
	}

	posiljka := manifestResp.Porudzbine[0].Posiljke[0]
	if posiljka.Rezultat != 0 {
		return &TestShipmentResponse{
			Success:        false,
			Errors:         []string{posiljka.Poruka},
			CreatedAt:      time.Now().Format(time.RFC3339),
			ProcessingTime: time.Since(startTime).Milliseconds(),
			RequestData:    buildRequestData(req),
			ResponseData:   manifestResp,
		}, nil
	}

	// Извлекаем стоимость (в para, 1 RSD = 100 para)
	costPara := posiljka.Postarina
	costRSD := float64(costPara) / 100.0

	// Извлекаем трек-номер (PrijemniBroj - это настоящий tracking number)
	trackingNumber := posiljka.PrijemniBroj
	if trackingNumber == "" {
		trackingNumber = posiljka.TrackingNumber // Fallback
	}

	// Успешный ответ
	h.logger.Info("Test shipment created successfully: tracking=%s, ID=%d, cost=%.2f RSD (%d para)",
		trackingNumber, posiljka.IDPosiljke, costRSD, costPara)

	return &TestShipmentResponse{
		Success:        true,
		TrackingNumber: trackingNumber,
		ManifestID:     manifestResp.IDManifesta,
		ShipmentID:     posiljka.IDPosiljke,
		ExternalID:     posiljka.BrojPosiljke,
		Cost:           int(costRSD), // Конвертируем в RSD (целое число)
		CreatedAt:      time.Now().Format(time.RFC3339),
		ProcessingTime: time.Since(startTime).Milliseconds(),
		RequestData:    buildRequestData(req),
		ResponseData:   manifestResp,
	}, nil
}

// buildRequestData строит данные запроса для ответа
func buildRequestData(req *TestShipmentRequest) fiber.Map {
	return fiber.Map{
		"recipient": fiber.Map{
			"name":    req.RecipientName,
			"phone":   req.RecipientPhone,
			"email":   req.RecipientEmail,
			"city":    req.RecipientCity,
			"address": req.RecipientAddress,
			"zip":     req.RecipientZip,
		},
		"sender": fiber.Map{
			"name":    req.SenderName,
			"phone":   req.SenderPhone,
			"email":   req.SenderEmail,
			"city":    req.SenderCity,
			"address": req.SenderAddress,
			"zip":     req.SenderZip,
		},
		"shipment": fiber.Map{
			"weight":             req.Weight,
			"content":            req.Content,
			"cod_amount":         req.CODAmount,
			"insured_value":      req.InsuredValue,
			"services":           req.Services,
			"delivery_method":    req.DeliveryMethod,
			"delivery_type":      req.DeliveryType,
			"payment_method":     req.PaymentMethod,
			"id_rukovanje":       req.IdRukovanje,
			"parcel_locker_code": req.ParcelLockerCode,
		},
	}
}

// =============================================================================
// TX 3-11: Новые тестовые endpoints для полной интеграции Post Express WSP API
// =============================================================================

// TestGetSettlementsRequest - запрос поиска населённых пунктов (TX 3)
type TestGetSettlementsRequest struct {
	Query string `json:"query" example:"Beograd"`
}

// TestGetSettlements тестирует TX 3 - поиск населённых пунктов
// @Summary Test TX 3 - GetNaselje (Search Settlements)
// @Description Search settlements using real Post Express WSP API (Transaction 3)
// @Tags post-express-test
// @Accept json
// @Produce json
// @Param request body TestGetSettlementsRequest true "Search settlements request"
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "Settlements list"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/v1/postexpress/test/tx3-settlements [post]
func (h *Handler) TestGetSettlements(c *fiber.Ctx) error {
	startTime := time.Now()

	var req TestGetSettlementsRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "postexpress.invalid_request", nil)
	}

	if req.Query == "" {
		req.Query = "Beograd"
	}

	if h.wspClient == nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "postexpress.wsp_client_not_initialized", nil)
	}

	h.logger.Info("Testing TX 3 GetSettlements: query=%s", req.Query)

	resp, err := h.wspClient.GetSettlements(c.Context(), req.Query)
	if err != nil {
		h.logger.Error("GetSettlements failed: %v", err)
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "postexpress.get_settlements_failed", fiber.Map{
			"error": err.Error(),
		})
	}

	result := fiber.Map{
		"success":         resp.Rezultat == 0,
		"rezultat":        resp.Rezultat,
		"poruka":          resp.Poruka,
		"query":           req.Query,
		"count":           len(resp.Naselja),
		"naselja":         resp.Naselja,
		"processing_time": time.Since(startTime).Milliseconds(),
	}

	h.logger.Info("TX 3 GetSettlements successful: query=%s, found=%d, Rezultat=%d", req.Query, len(resp.Naselja), resp.Rezultat)

	return utils.SendSuccessResponse(c, result, fmt.Sprintf("TX 3: Found %d settlements (Rezultat: %d)", len(resp.Naselja), resp.Rezultat))
}

// TestGetStreetsRequest - запрос поиска улиц (TX 4)
type TestGetStreetsRequest struct {
	SettlementID int    `json:"settlement_id" example:"123"`
	Query        string `json:"query" example:"Takovska"`
}

// TestGetStreets тестирует TX 4 - поиск улиц в населённом пункте
// @Summary Test TX 4 - GetUlica (Search Streets)
// @Description Search streets in settlement using real Post Express WSP API (Transaction 4)
// @Tags post-express-test
// @Accept json
// @Produce json
// @Param request body TestGetStreetsRequest true "Search streets request"
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "Streets list"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/v1/postexpress/test/tx4-streets [post]
func (h *Handler) TestGetStreets(c *fiber.Ctx) error {
	startTime := time.Now()

	var req TestGetStreetsRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "postexpress.invalid_request", nil)
	}

	if req.SettlementID <= 0 {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "postexpress.settlement_id_required", nil)
	}

	if req.Query == "" {
		req.Query = "Takovska"
	}

	if h.wspClient == nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "postexpress.wsp_client_not_initialized", nil)
	}

	h.logger.Info("Testing TX 4 GetStreets: settlement_id=%d, query=%s", req.SettlementID, req.Query)

	resp, err := h.wspClient.GetStreets(c.Context(), req.SettlementID, req.Query)
	if err != nil {
		h.logger.Error("GetStreets failed: %v", err)
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "postexpress.get_streets_failed", fiber.Map{
			"error": err.Error(),
		})
	}

	result := fiber.Map{
		"success":         resp.Rezultat == 0,
		"rezultat":        resp.Rezultat,
		"poruka":          resp.Poruka,
		"settlement_id":   req.SettlementID,
		"query":           req.Query,
		"count":           len(resp.Ulice),
		"ulice":           resp.Ulice,
		"processing_time": time.Since(startTime).Milliseconds(),
	}

	h.logger.Info("TX 4 GetStreets successful: settlement_id=%d, query=%s, found=%d, Rezultat=%d", req.SettlementID, req.Query, len(resp.Ulice), resp.Rezultat)

	return utils.SendSuccessResponse(c, result, fmt.Sprintf("TX 4: Found %d streets (Rezultat: %d)", len(resp.Ulice), resp.Rezultat))
}

// TestValidateAddressRequest - запрос валидации адреса (TX 6)
type TestValidateAddressRequest struct {
	SettlementID int    `json:"settlement_id" example:"123"`
	StreetID     int    `json:"street_id,omitempty" example:"456"`
	HouseNumber  string `json:"house_number" example:"2"`
	PostalCode   string `json:"postal_code" example:"11000"`
}

// TestValidateAddress тестирует TX 6 - валидация адреса
// @Summary Test TX 6 - ProveraAdrese (Validate Address)
// @Description Validate address using real Post Express WSP API (Transaction 6)
// @Tags post-express-test
// @Accept json
// @Produce json
// @Param request body TestValidateAddressRequest true "Validate address request"
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "Address validation result"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/v1/postexpress/test/tx6-validate-address [post]
func (h *Handler) TestValidateAddress(c *fiber.Ctx) error {
	startTime := time.Now()

	var req TestValidateAddressRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "postexpress.invalid_request", nil)
	}

	if req.SettlementID <= 0 || req.HouseNumber == "" || req.PostalCode == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "postexpress.missing_address_data", fiber.Map{
			"required_fields": []string{"settlement_id", "house_number", "postal_code"},
		})
	}

	if h.wspClient == nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "postexpress.wsp_client_not_initialized", nil)
	}

	h.logger.Info("Testing TX 6 ValidateAddress: settlement_id=%d, street_id=%d, number=%s, postal=%s",
		req.SettlementID, req.StreetID, req.HouseNumber, req.PostalCode)

	wspReq := &postexpress.AddressValidationRequest{
		TipAdrese:     0,                 // 0 = стандартный тип адреса
		IdRukovanje:   71,                // Дефолтная услуга доставки (можно любую)
		IdNaselje:     req.SettlementID,
		IdUlica:       req.StreetID,
		BrojPodbroj:   req.HouseNumber,
		PostanskiBroj: req.PostalCode,
	}

	resp, err := h.wspClient.ValidateAddress(c.Context(), wspReq)
	if err != nil {
		h.logger.Error("ValidateAddress failed: %v", err)
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "postexpress.validate_address_failed", fiber.Map{
			"error": err.Error(),
		})
	}

	result := fiber.Map{
		"success":         resp.Rezultat == 0,
		"rezultat":        resp.Rezultat,
		"poruka":          resp.Poruka,
		"postoji_adresa":  resp.PostojiAdresa,
		"id_naselje":      resp.IdNaselje,
		"naziv_naselja":   resp.NazivNaselja,
		"id_ulica":        resp.IdUlica,
		"naziv_ulice":     resp.NazivUlice,
		"broj":            resp.Broj,
		"postanski_broj":  resp.PostanskiBroj,
		"pak":             resp.PAK,
		"id_poste":        resp.IdPoste,
		"naziv_poste":     resp.NazivPoste,
		"processing_time": time.Since(startTime).Milliseconds(),
	}

	h.logger.Info("TX 6 ValidateAddress successful: postoji=%t, Rezultat=%d", resp.PostojiAdresa, resp.Rezultat)

	return utils.SendSuccessResponse(c, result, fmt.Sprintf("TX 6: Address exists: %t (Rezultat: %d)", resp.PostojiAdresa, resp.Rezultat))
}

// TestCheckServiceAvailabilityRequest - запрос проверки доступности услуги (TX 9)
type TestCheckServiceAvailabilityRequest struct {
	ServiceID      int    `json:"service_id" example:"71"`
	FromPostalCode string `json:"from_postal_code" example:"11000"`
	ToPostalCode   string `json:"to_postal_code" example:"21000"`
}

// TestCheckServiceAvailability тестирует TX 9 - проверка доступности услуги
// @Summary Test TX 9 - ProveraDostupnostiUsluge (Check Service Availability)
// @Description Check service availability using real Post Express WSP API (Transaction 9)
// @Tags post-express-test
// @Accept json
// @Produce json
// @Param request body TestCheckServiceAvailabilityRequest true "Check service availability request"
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "Service availability result"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/v1/postexpress/test/tx9-service-availability [post]
func (h *Handler) TestCheckServiceAvailability(c *fiber.Ctx) error {
	startTime := time.Now()

	var req TestCheckServiceAvailabilityRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "postexpress.invalid_request", nil)
	}

	if req.ServiceID <= 0 || req.FromPostalCode == "" || req.ToPostalCode == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "postexpress.missing_service_data", fiber.Map{
			"required_fields": []string{"service_id", "from_postal_code", "to_postal_code"},
		})
	}

	if h.wspClient == nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "postexpress.wsp_client_not_initialized", nil)
	}

	h.logger.Info("Testing TX 9 CheckServiceAvailability: service_id=%d, from=%s, to=%s",
		req.ServiceID, req.FromPostalCode, req.ToPostalCode)

	// Получаем IdNaselje для отправления и прибытия через TX3
	idNaseljeOdlaska, err := h.getSettlementIDByPostalCode(c.Context(), req.FromPostalCode)
	if err != nil {
		h.logger.Error("Failed to get settlement ID for from postal code %s: %v", req.FromPostalCode, err)
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "postexpress.invalid_from_postal_code", fiber.Map{
			"postal_code": req.FromPostalCode,
			"error":       err.Error(),
		})
	}

	idNaseljeDolaska, err := h.getSettlementIDByPostalCode(c.Context(), req.ToPostalCode)
	if err != nil {
		h.logger.Error("Failed to get settlement ID for to postal code %s: %v", req.ToPostalCode, err)
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "postexpress.invalid_to_postal_code", fiber.Map{
			"postal_code": req.ToPostalCode,
			"error":       err.Error(),
		})
	}

	h.logger.Debug("Resolved settlement IDs: from=%d (postal=%s), to=%d (postal=%s)",
		idNaseljeOdlaska, req.FromPostalCode, idNaseljeDolaska, req.ToPostalCode)

	wspReq := &postexpress.ServiceAvailabilityRequest{
		TipAdrese:            0, // 0 = стандартный тип адреса
		IdRukovanje:          req.ServiceID,
		IdNaseljeOdlaska:     idNaseljeOdlaska, // Получено через TX3
		IdNaseljeDolaska:     idNaseljeDolaska, // Получено через TX3
		PostanskiBrojOdlaska: req.FromPostalCode,
		PostanskiBrojDolaska: req.ToPostalCode,
		Datum:                time.Now().Format("02.01.2006"), // Формат DD.MM.YYYY (Post Express требование)
	}

	resp, err := h.wspClient.CheckServiceAvailability(c.Context(), wspReq)
	if err != nil {
		h.logger.Error("CheckServiceAvailability failed: %v", err)
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "postexpress.check_service_failed", fiber.Map{
			"error": err.Error(),
		})
	}

	result := fiber.Map{
		"success":         resp.Rezultat == 0,
		"rezultat":        resp.Rezultat,
		"poruka":          resp.Poruka,
		"dostupna":        resp.Dostupna,
		"id_rukovanje":    resp.IdRukovanje,
		"naziv_usluge":    resp.NazivUsluge,
		"ocekivano_dana":  resp.OcekivanoDana,
		"napomena":        resp.Napomena,
		"processing_time": time.Since(startTime).Milliseconds(),
	}

	h.logger.Info("TX 9 CheckServiceAvailability successful: dostupna=%t, days=%d, Rezultat=%d",
		resp.Dostupna, resp.OcekivanoDana, resp.Rezultat)

	return utils.SendSuccessResponse(c, result, fmt.Sprintf("TX 9: Service available: %t, Expected days: %d (Rezultat: %d)", resp.Dostupna, resp.OcekivanoDana, resp.Rezultat))
}

// TestCalculatePostageRequest - запрос расчёта стоимости доставки (TX 11)
type TestCalculatePostageRequest struct {
	ServiceID      int    `json:"service_id" example:"71"`
	FromPostalCode string `json:"from_postal_code" example:"11000"`
	ToPostalCode   string `json:"to_postal_code" example:"21000"`
	Weight         int    `json:"weight" example:"500"` // граммы
	CODAmount      int    `json:"cod_amount,omitempty" example:"0"` // para (1 RSD = 100 para)
	InsuredValue   int    `json:"insured_value,omitempty" example:"0"` // para
	Services       string `json:"services,omitempty" example:"PNA"` // дополнительные услуги
}

// TestCalculatePostage тестирует TX 11 - расчёт стоимости доставки
// @Summary Test TX 11 - PostarinaPosiljke (Calculate Postage)
// @Description Calculate postage using real Post Express WSP API (Transaction 11)
// @Tags post-express-test
// @Accept json
// @Produce json
// @Param request body TestCalculatePostageRequest true "Calculate postage request"
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "Postage calculation result"
// @Failure 400 {object} utils.ErrorResponseSwag "Invalid request"
// @Failure 500 {object} utils.ErrorResponseSwag "Server error"
// @Router /api/v1/postexpress/test/tx11-calculate-postage [post]
func (h *Handler) TestCalculatePostage(c *fiber.Ctx) error {
	startTime := time.Now()

	var req TestCalculatePostageRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "postexpress.invalid_request", nil)
	}

	if req.ServiceID <= 0 || req.FromPostalCode == "" || req.ToPostalCode == "" || req.Weight <= 0 {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "postexpress.missing_postage_data", fiber.Map{
			"required_fields": []string{"service_id", "from_postal_code", "to_postal_code", "weight"},
		})
	}

	if req.Services == "" {
		req.Services = "PNA"
	}

	if h.wspClient == nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "postexpress.wsp_client_not_initialized", nil)
	}

	h.logger.Info("Testing TX 11 CalculatePostage: service_id=%d, from=%s, to=%s, weight=%dg",
		req.ServiceID, req.FromPostalCode, req.ToPostalCode, req.Weight)

	wspReq := &postexpress.PostageCalculationRequest{
		IdRukovanje:          req.ServiceID,
		IdZemlja:             0, // 0 = domestic shipments
		PostanskiBrojOdlaska: req.FromPostalCode,
		PostanskiBrojDolaska: req.ToPostalCode,
		Masa:                 req.Weight,
		Otkupnina:            req.CODAmount,
		Vrednost:             req.InsuredValue,
		PosebneUsluge:        req.Services,
	}

	resp, err := h.wspClient.CalculatePostage(c.Context(), wspReq)
	if err != nil {
		h.logger.Error("CalculatePostage failed: %v", err)
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "postexpress.calculate_postage_failed", fiber.Map{
			"error": err.Error(),
		})
	}

	result := fiber.Map{
		"success":                resp.Rezultat == 0,
		"rezultat":               resp.Rezultat,
		"poruka":                 resp.Poruka,
		"postarina_para":         resp.Postarina,
		"postarina_rsd":          float64(resp.Postarina) / 100.0,
		"id_rukovanje":           resp.IdRukovanje,
		"naziv_usluge":           resp.NazivUsluge,
		"postanski_broj_odlaska": resp.PostanskiBrojOdlaska,
		"postanski_broj_dolaska": resp.PostanskiBrojDolaska,
		"masa":                   resp.Masa,
		"otkupnina":              resp.Otkupnina,
		"vrednost":               resp.Vrednost,
		"posebne_usluge":         resp.PosebneUsluge,
		"napomena":               resp.Napomena,
		"processing_time":        time.Since(startTime).Milliseconds(),
	}

	h.logger.Info("TX 11 CalculatePostage successful: postage=%d para (%.2f RSD), Rezultat=%d",
		resp.Postarina, float64(resp.Postarina)/100.0, resp.Rezultat)

	return utils.SendSuccessResponse(c, result, fmt.Sprintf("TX 11: Postage: %.2f RSD (Rezultat: %d)", float64(resp.Postarina)/100.0, resp.Rezultat))
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

// getSettlementIDByPostalCode получает IdNaselje по почтовому индексу через TX3
func (h *Handler) getSettlementIDByPostalCode(ctx context.Context, postalCode string) (int, error) {
	// Mapping популярных почтовых кодов (оптимизация для часто используемых)
	knownPostalCodes := map[string]int{
		"11000": 100001, // Beograd
		"21000": 100039, // Novi Sad
		"18000": 100040, // Niš
		"34000": 100041, // Kragujevac
		"31000": 100042, // Užice
	}

	// Проверяем известные коды
	if id, ok := knownPostalCodes[postalCode]; ok {
		h.logger.Debug("Using known postal code mapping: %s -> %d", postalCode, id)
		return id, nil
	}

	// Для неизвестных кодов пробуем поиск через API
	// TX3 требует минимум 2 символа, используем широкий поиск
	searchQueries := []string{
		"Be", "No", "Ni", "Kr", "Su", "Pa", "Le", "Sm", "Va", "Za",
	}

	for _, query := range searchQueries {
		resp, err := h.wspClient.GetSettlements(ctx, query)
		if err != nil {
			continue
		}

		if resp.Rezultat != 0 {
			continue
		}

		// Ищем населенный пункт с нужным почтовым индексом
		for _, naselje := range resp.Naselja {
			if naselje.PostanskiBroj == postalCode {
				id := naselje.Id
				if id == 0 {
					id = naselje.IdNaselje
				}
				h.logger.Info("Found settlement for postal code %s: Id=%d, Naziv=%s", postalCode, id, naselje.Naziv)
				return id, nil
			}
		}
	}

	return 0, fmt.Errorf("settlement not found for postal code: %s (tried multiple searches)", postalCode)
}
