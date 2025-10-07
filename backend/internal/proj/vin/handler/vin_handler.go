package handler

import (
	"strconv"
	"strings"

	"backend/internal/proj/vin/models"
	"backend/internal/proj/vin/service"
	"backend/internal/proj/vin/storage/postgres"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// VINHandler представляет HTTP handler для VIN операций
type VINHandler struct {
	service *service.VINService
}

// NewVINHandler создает новый VIN handler
func NewVINHandler(db *sqlx.DB) *VINHandler {
	storage := postgres.NewVINStorage(db)
	vinService := service.NewVINService(storage)

	return &VINHandler{
		service: vinService,
	}
}

// RegisterRoutes регистрирует маршруты для VIN
func (h *VINHandler) RegisterRoutes(app *fiber.App) {
	// Public endpoints
	api := app.Group("/api/v1/vin")

	// Декодирование VIN (может быть публичным с ограничениями)
	api.Post("/decode", h.DecodeVIN)
	api.Get("/validate/:vin", h.ValidateVIN)

	// Protected endpoints - требуют авторизации
	// Middleware будет применён при регистрации в модуле
	api.Get("/history", h.GetHistory)
	api.Post("/auto-fill", h.AutoFillFromVIN)
	api.Get("/stats", h.GetStats)
	api.Get("/report/:vin", h.ExportReport)
}

// DecodeVIN декодирует VIN номер
// @Summary Декодировать VIN номер
// @Description Декодирует VIN номер и возвращает информацию об автомобиле
// @Tags vin
// @Accept json
// @Produce json
// @Param request body models.VINDecodeRequest true "VIN для декодирования"
// @Success 200 {object} utils.SuccessResponseSwag "Успешное декодирование"
// @Failure 400 {object} utils.ErrorResponseSwag "Некорректный запрос"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Router /api/v1/vin/decode [post]
func (h *VINHandler) DecodeVIN(c *fiber.Ctx) error {
	var req models.VINDecodeRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidRequest", nil)
	}

	// Нормализуем VIN
	req.VIN = strings.ToUpper(strings.TrimSpace(req.VIN))

	// Получаем user ID если есть (необязательно)
	var userID *int64
	if userIDStr := c.Get("user_id"); userIDStr != "" {
		if id, err := strconv.ParseInt(userIDStr, 10, 64); err == nil {
			userID = &id
		}
	}

	// Декодируем VIN
	response, err := h.service.DecodeVIN(c.Context(), &req, userID)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "vin.decodeFailed", fiber.Map{"error": err.Error()})
	}

	if !response.Success {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "vin.invalidVIN", fiber.Map{
			"error": response.Error,
		})
	}

	return utils.SendSuccessResponse(c, response, "vin.decodeSuccess")
}

// ValidateVIN проверяет валидность VIN номера
// @Summary Проверить валидность VIN номера
// @Description Проверяет корректность VIN номера без полного декодирования
// @Tags vin
// @Accept json
// @Produce json
// @Param vin path string true "VIN номер для проверки"
// @Success 200 {object} utils.SuccessResponseSwag "VIN валиден"
// @Failure 400 {object} utils.ErrorResponseSwag "VIN невалиден"
// @Router /api/v1/vin/validate/{vin} [get]
func (h *VINHandler) ValidateVIN(c *fiber.Ctx) error {
	vin := strings.ToUpper(strings.TrimSpace(c.Params("vin")))

	decoder := service.NewVINDecoder()
	err := decoder.ValidateVIN(vin)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "vin.invalidVIN", fiber.Map{
			"error": err.Error(),
		})
	}

	// Получаем базовую информацию
	basicInfo, _ := decoder.DecodeBasicInfo(vin)

	response := map[string]interface{}{
		"valid": true,
		"vin":   vin,
	}

	if basicInfo != nil {
		response["basic_info"] = basicInfo
	}

	return utils.SendSuccessResponse(c, response, "vin.valid")
}

// GetHistory получает историю проверок VIN
// @Summary Получить историю проверок VIN
// @Description Возвращает историю проверок VIN для текущего пользователя
// @Tags vin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param vin query string false "Фильтр по конкретному VIN"
// @Param limit query int false "Лимит записей" default(20)
// @Param offset query int false "Смещение" default(0)
// @Success 200 {object} utils.SuccessResponseSwag "История проверок"
// @Failure 401 {object} utils.ErrorResponseSwag "Не авторизован"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Router /api/v1/vin/history [get]
func (h *VINHandler) GetHistory(c *fiber.Ctx) error {
	// Получаем user ID из контекста
	userIDStr := c.Get("user_id")
	if userIDStr == "" {
		return utils.SendErrorResponse(c, fiber.StatusUnauthorized, "authentication.required", nil)
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "authentication.invalidUser", nil)
	}

	// Параметры запроса
	req := models.VINHistoryRequest{
		UserID: &userID,
		VIN:    c.Query("vin"),
		Limit:  c.QueryInt("limit", 20),
		Offset: c.QueryInt("offset", 0),
	}

	history, err := h.service.GetHistory(c.Context(), &req)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "vin.historyFailed", fiber.Map{"error": err.Error()})
	}

	return utils.SendSuccessResponse(c, history, "vin.historySuccess")
}

// AutoFillFromVIN заполняет данные из VIN
// @Summary Автозаполнение данных из VIN
// @Description Декодирует VIN и возвращает данные для автозаполнения формы создания объявления
// @Tags vin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body map[string]string true "VIN номер"
// @Success 200 {object} utils.SuccessResponseSwag "Данные для автозаполнения"
// @Failure 400 {object} utils.ErrorResponseSwag "Некорректный запрос"
// @Failure 401 {object} utils.ErrorResponseSwag "Не авторизован"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Router /api/v1/vin/auto-fill [post]
func (h *VINHandler) AutoFillFromVIN(c *fiber.Ctx) error {
	// Получаем user ID из контекста
	userIDStr := c.Get("user_id")
	if userIDStr == "" {
		return utils.SendErrorResponse(c, fiber.StatusUnauthorized, "authentication.required", nil)
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "authentication.invalidUser", nil)
	}

	// Парсим запрос
	var req map[string]string
	if err := c.BodyParser(&req); err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "marketplace.invalidRequest", nil)
	}

	vin, ok := req["vin"]
	if !ok || vin == "" {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "vin.required", nil)
	}

	vin = strings.ToUpper(strings.TrimSpace(vin))

	// Получаем данные для автозаполнения
	autoFillData, err := h.service.AutoFillFromVIN(c.Context(), vin, &userID)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "vin.autoFillFailed", fiber.Map{
			"error": err.Error(),
		})
	}

	return utils.SendSuccessResponse(c, autoFillData, "vin.autoFillSuccess")
}

// GetStats получает статистику использования VIN декодера
// @Summary Получить статистику использования VIN декодера
// @Description Возвращает статистику проверок VIN для текущего пользователя
// @Tags vin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponseSwag "Статистика"
// @Failure 401 {object} utils.ErrorResponseSwag "Не авторизован"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Router /api/v1/vin/stats [get]
func (h *VINHandler) GetStats(c *fiber.Ctx) error {
	// Получаем user ID из контекста
	userIDStr := c.Get("user_id")
	if userIDStr == "" {
		return utils.SendErrorResponse(c, fiber.StatusUnauthorized, "authentication.required", nil)
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "authentication.invalidUser", nil)
	}

	stats, err := h.service.GetVINStats(c.Context(), &userID)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "vin.statsFailed", fiber.Map{"error": err.Error()})
	}

	return utils.SendSuccessResponse(c, stats, "vin.statsSuccess")
}

// ExportReport экспортирует отчет о VIN
// @Summary Экспортировать отчет о VIN
// @Description Генерирует и возвращает детальный отчет о VIN в JSON формате
// @Tags vin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param vin path string true "VIN номер"
// @Success 200 {object} map[string]interface{} "Отчет о VIN"
// @Failure 400 {object} utils.ErrorResponseSwag "Некорректный VIN"
// @Failure 401 {object} utils.ErrorResponseSwag "Не авторизован"
// @Failure 500 {object} utils.ErrorResponseSwag "Внутренняя ошибка сервера"
// @Router /api/v1/vin/report/{vin} [get]
func (h *VINHandler) ExportReport(c *fiber.Ctx) error {
	// Получаем user ID из контекста
	userIDStr := c.Get("user_id")
	if userIDStr == "" {
		return utils.SendErrorResponse(c, fiber.StatusUnauthorized, "authentication.required", nil)
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "authentication.invalidUser", nil)
	}

	vin := strings.ToUpper(strings.TrimSpace(c.Params("vin")))

	// Генерируем отчет
	reportData, err := h.service.ExportVINReport(c.Context(), vin, &userID)
	if err != nil {
		return utils.SendErrorResponse(c, fiber.StatusBadRequest, "vin.reportFailed", fiber.Map{
			"error": err.Error(),
		})
	}

	// Устанавливаем заголовки для скачивания
	c.Set("Content-Type", "application/json")
	c.Set("Content-Disposition", "attachment; filename=\"vin_report_"+vin+".json\"")

	return c.Send(reportData)
}
