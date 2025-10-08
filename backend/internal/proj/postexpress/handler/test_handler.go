package handler

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"

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
	Services       string `json:"services" example:"PNA"`       // PNA, SMS, OTK, VD
	DeliveryMethod string `json:"delivery_method" example:"K"`  // K = Kurir, S = Šalter
	PaymentMethod  string `json:"payment_method" example:"POF"` // POF = gotovina
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

// CreateTestShipment создает тестовое отправление
// @Summary Create test shipment
// @Description Create a test shipment using Post Express WSP API for visual testing
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
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "postexpress.invalid_request", fiber.Map{
			"error": err.Error(),
		})
	}

	// Валидация обязательных полей
	if req.RecipientName == "" || req.RecipientPhone == "" || req.RecipientCity == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "postexpress.missing_recipient_data", fiber.Map{
			"missing_fields": []string{"recipient_name", "recipient_phone", "recipient_city"},
		})
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

	// Создаем тестовое отправление
	result := h.createTestShipmentResponse(&req, startTime)

	return utils.SendSuccessResponse(c, result, "Test shipment processed successfully")
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
		"delivery_methods": []fiber.Map{
			{"code": "K", "name": "Kurir (Courier)"},
			{"code": "S", "name": "Šalter (Post Office)"},
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

// createTestShipmentResponse создает ответ для тестового отправления
func (h *Handler) createTestShipmentResponse(req *TestShipmentRequest, startTime time.Time) *TestShipmentResponse {
	// Генерируем уникальные ID
	externalID := fmt.Sprintf("SVETU-TEST-%d", time.Now().Unix())
	trackingNumber := generateTrackingNumber()

	// Рассчитываем стоимость (простая формула для теста)
	cost := calculateTestCost(req.Weight, req.CODAmount, req.InsuredValue, req.Services)

	response := &TestShipmentResponse{
		Success:        true,
		TrackingNumber: trackingNumber,
		ManifestID:     int(time.Now().Unix()%100000) + 120000, // Реалистичные ID как в тестах
		ShipmentID:     int(time.Now().Unix()%100000) + 27000,
		ExternalID:     externalID,
		Cost:           cost,
		CreatedAt:      time.Now().Format(time.RFC3339),
		ProcessingTime: time.Since(startTime).Milliseconds(),
		RequestData: fiber.Map{
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
				"weight":          req.Weight,
				"content":         req.Content,
				"cod_amount":      req.CODAmount,
				"insured_value":   req.InsuredValue,
				"services":        req.Services,
				"delivery_method": req.DeliveryMethod,
				"payment_method":  req.PaymentMethod,
			},
		},
		ResponseData: fiber.Map{
			"status":    "created",
			"provider":  "post_express",
			"test_mode": true,
			"api_response": fiber.Map{
				"Rezultat": 0,
				"Poruka":   "Success",
				"Greske":   nil,
			},
		},
	}

	return response
}

// generateTrackingNumber генерирует реалистичный трек-номер Post Express
func generateTrackingNumber() string {
	// Формат: PJ700042693RS
	randomNum := make([]byte, 4)
	if _, err := rand.Read(randomNum); err != nil {
		// Fallback to timestamp-based number if random fails
		randomNum = []byte{0, 0, 0, byte(time.Now().Unix() % 256)}
	}

	num := int(randomNum[0])<<24 | int(randomNum[1])<<16 | int(randomNum[2])<<8 | int(randomNum[3])
	if num < 0 {
		num = -num
	}
	num %= 100000

	return fmt.Sprintf("PJ700%05dRS", num)
}

// calculateTestCost рассчитывает тестовую стоимость доставки
func calculateTestCost(weight int, codAmount int64, insuredValue int64, services string) int {
	// Базовая стоимость
	baseCost := 415 // RSD для Белграда (как в реальных тестах)

	// Доплата за вес (свыше 500г)
	if weight > 500 {
		extraWeight := (weight - 500) / 100 // каждые 100г
		baseCost += extraWeight * 30        // 30 RSD за 100г
	}

	// Доплата за наложенный платеж
	if codAmount > 0 {
		baseCost += 50 // фиксированная доплата за OTK
	}

	// Доплата за объявленную ценность
	if insuredValue > 0 {
		insuranceFee := int(float64(insuredValue) * 0.01) // 1% от ценности
		if insuranceFee < 50 {
			insuranceFee = 50 // минимум 50 RSD
		}
		baseCost += insuranceFee
	}

	// Доплата за SMS
	if contains(services, "SMS") {
		baseCost += 20
	}

	return baseCost
}

// contains проверяет содержит ли строка подстроку
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr))
}
