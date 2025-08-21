package handler

import (
	"encoding/csv"
	"fmt"
	"time"

	"backend/internal/proj/admin/logistics/service"
	"backend/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// AnalyticsHandler обработчик для аналитики и отчетов
type AnalyticsHandler struct {
	analyticsService *service.AnalyticsService
}

// NewAnalyticsHandler создает новый обработчик аналитики
func NewAnalyticsHandler(analyticsService *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

// GetPerformanceMetrics godoc
// @Summary Получить метрики производительности
// @Description Возвращает метрики производительности логистики за указанный период
// @Tags admin-logistics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param date_from query string false "Дата начала периода (YYYY-MM-DD)"
// @Param date_to query string false "Дата окончания периода (YYYY-MM-DD)"
// @Param group_by query string false "Группировка (day, week, month)" default(day)
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "Performance metrics"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request - invalid parameters"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/admin/logistics/analytics/performance [get]
func (h *AnalyticsHandler) GetPerformanceMetrics(c *fiber.Ctx) error {
	// Проверка прав доступа
	userID := c.Locals("user_id")
	if userID == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.unauthorized")
	}

	// Парсим параметры
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")
	periodDays := c.QueryInt("period", 30)
	groupBy := c.Query("group_by", "day")

	// Если даты не указаны, берем период в днях
	var fromDate, toDate time.Time
	var err error

	if dateFrom == "" {
		fromDate = time.Now().AddDate(0, 0, -periodDays)
	} else {
		fromDate, err = time.Parse("2006-01-02", dateFrom)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "logistics.invalid_date_from")
		}
	}

	if dateTo == "" {
		toDate = time.Now()
	} else {
		toDate, err = time.Parse("2006-01-02", dateTo)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "logistics.invalid_date_to")
		}
	}

	// Получаем метрики
	metrics, err := h.analyticsService.GetPerformanceMetrics(c.Context(), fromDate, toDate, groupBy)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "logistics.analytics.error")
	}

	return utils.SuccessResponse(c, metrics)
}

// GetFinancialReport godoc
// @Summary Получить финансовый отчет
// @Description Возвращает финансовые показатели логистики за указанный период
// @Tags admin-logistics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param date_from query string false "Дата начала периода (YYYY-MM-DD)"
// @Param date_to query string false "Дата окончания периода (YYYY-MM-DD)"
// @Success 200 {object} utils.SuccessResponseSwag{data=map[string]interface{}} "Financial report"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 403 {object} utils.ErrorResponseSwag "Forbidden - insufficient permissions"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request - invalid parameters"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/admin/logistics/analytics/financial [get]
func (h *AnalyticsHandler) GetFinancialReport(c *fiber.Ctx) error {
	// Проверка прав доступа
	userID := c.Locals("user_id")
	if userID == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.unauthorized")
	}

	// TODO: Добавить проверку прав на просмотр финансовых данных

	// Парсим параметры
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")
	periodDays := c.QueryInt("period", 30)

	// Если даты не указаны, берем период в днях
	var fromDate, toDate time.Time
	var err error

	if dateFrom == "" {
		fromDate = time.Now().AddDate(0, 0, -periodDays)
	} else {
		fromDate, err = time.Parse("2006-01-02", dateFrom)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "logistics.invalid_date_from")
		}
	}

	if dateTo == "" {
		toDate = time.Now()
	} else {
		toDate, err = time.Parse("2006-01-02", dateTo)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "logistics.invalid_date_to")
		}
	}

	// Получаем финансовый отчет
	report, err := h.analyticsService.GetFinancialReport(c.Context(), fromDate, toDate)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "logistics.financial.error")
	}

	return utils.SuccessResponse(c, report)
}

// ExportReport godoc
// @Summary Экспортировать отчет
// @Description Экспортирует отчет по логистике в указанном формате
// @Tags admin-logistics
// @Accept json
// @Produce application/octet-stream
// @Security ApiKeyAuth
// @Param format query string false "Формат экспорта (csv, xlsx)" default(csv)
// @Param report_type query string false "Тип отчета (shipments, problems, performance)" default(shipments)
// @Param date_from query string false "Дата начала периода (YYYY-MM-DD)"
// @Param date_to query string false "Дата окончания периода (YYYY-MM-DD)"
// @Success 200 {file} binary "Report file"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 400 {object} utils.ErrorResponseSwag "Bad request - invalid parameters"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/admin/logistics/analytics/export [get]
func (h *AnalyticsHandler) ExportReport(c *fiber.Ctx) error {
	// Проверка прав доступа
	userID := c.Locals("user_id")
	if userID == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.unauthorized")
	}

	// Парсим параметры
	format := c.Query("format", "csv")
	reportType := c.Query("report_type", "shipments")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	// Парсим даты
	var fromDate, toDate time.Time
	var err error

	if dateFrom == "" {
		fromDate = time.Now().AddDate(0, 0, -30)
	} else {
		fromDate, err = time.Parse("2006-01-02", dateFrom)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "logistics.invalid_date_from")
		}
	}

	if dateTo == "" {
		toDate = time.Now()
	} else {
		toDate, err = time.Parse("2006-01-02", dateTo)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "logistics.invalid_date_to")
		}
	}

	// Получаем данные для экспорта
	data, err := h.analyticsService.GetReportData(c.Context(), reportType, fromDate, toDate)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "logistics.export.error")
	}

	// Экспортируем в зависимости от формата
	switch format {
	case "csv":
		return h.exportCSV(c, data, reportType)
	case "xlsx":
		return h.exportXLSX(c, data, reportType)
	default:
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "logistics.export.invalid_format")
	}
}

// exportCSV экспортирует данные в CSV формате
func (h *AnalyticsHandler) exportCSV(c *fiber.Ctx, data interface{}, reportType string) error {
	// Устанавливаем заголовки для скачивания файла
	filename := fmt.Sprintf("logistics_%s_%s.csv", reportType, time.Now().Format("20060102_150405"))
	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	// Создаем CSV writer
	writer := csv.NewWriter(c.Response().BodyWriter())
	defer writer.Flush()

	// В зависимости от типа отчета формируем CSV
	switch reportType {
	case "shipments":
		// Заголовки
		headers := []string{
			"ID", "Type", "Tracking Number", "Status",
			"Created At", "Delivered At", "Sender City", "Recipient City",
		}
		if err := writer.Write(headers); err != nil {
			return err
		}

		// Данные
		if shipments, ok := data.([]map[string]interface{}); ok {
			for _, shipment := range shipments {
				row := []string{
					fmt.Sprintf("%v", shipment["id"]),
					fmt.Sprintf("%v", shipment["type"]),
					fmt.Sprintf("%v", shipment["tracking_number"]),
					fmt.Sprintf("%v", shipment["status"]),
					fmt.Sprintf("%v", shipment["created_at"]),
					fmt.Sprintf("%v", shipment["delivered_at"]),
					fmt.Sprintf("%v", shipment["sender_city"]),
					fmt.Sprintf("%v", shipment["recipient_city"]),
				}
				if err := writer.Write(row); err != nil {
					return err
				}
			}
		}

	case "problems":
		// Заголовки
		headers := []string{
			"ID", "Shipment ID", "Type", "Problem Type", "Severity",
			"Status", "Created At", "Resolved At",
		}
		if err := writer.Write(headers); err != nil {
			return err
		}

		// Данные
		if problems, ok := data.([]map[string]interface{}); ok {
			for _, problem := range problems {
				row := []string{
					fmt.Sprintf("%v", problem["id"]),
					fmt.Sprintf("%v", problem["shipment_id"]),
					fmt.Sprintf("%v", problem["shipment_type"]),
					fmt.Sprintf("%v", problem["problem_type"]),
					fmt.Sprintf("%v", problem["severity"]),
					fmt.Sprintf("%v", problem["status"]),
					fmt.Sprintf("%v", problem["created_at"]),
					fmt.Sprintf("%v", problem["resolved_at"]),
				}
				if err := writer.Write(row); err != nil {
					return err
				}
			}
		}

	case "performance":
		// Заголовки
		headers := []string{
			"Date", "Total Shipments", "Delivered", "In Transit",
			"Problems", "Success Rate", "Avg Delivery Time (hours)",
		}
		if err := writer.Write(headers); err != nil {
			return err
		}

		// Данные
		if metrics, ok := data.([]map[string]interface{}); ok {
			for _, metric := range metrics {
				row := []string{
					fmt.Sprintf("%v", metric["date"]),
					fmt.Sprintf("%v", metric["total_shipments"]),
					fmt.Sprintf("%v", metric["delivered"]),
					fmt.Sprintf("%v", metric["in_transit"]),
					fmt.Sprintf("%v", metric["problems"]),
					fmt.Sprintf("%v", metric["success_rate"]),
					fmt.Sprintf("%v", metric["avg_delivery_time"]),
				}
				if err := writer.Write(row); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// exportXLSX экспортирует данные в XLSX формате (упрощенный XML)
func (h *AnalyticsHandler) exportXLSX(c *fiber.Ctx, data interface{}, reportType string) error {
	// Устанавливаем заголовки для скачивания файла
	filename := fmt.Sprintf("logistics_%s_%s.xlsx", reportType, time.Now().Format("20060102_150405"))
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	// Создаем простой XML структуру Excel файла
	xmlContent := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main">
<sheetData>`

	rowIndex := 1

	// В зависимости от типа отчета формируем Excel
	switch reportType {
	case "shipments":
		// Заголовки
		xmlContent += fmt.Sprintf(`<row r="%d">`, rowIndex)
		headers := []string{"ID", "Type", "Tracking Number", "Status", "Created At", "Delivered At", "Sender City", "Recipient City"}
		for colIndex, header := range headers {
			xmlContent += fmt.Sprintf(`<c r="%s%d" t="inlineStr"><is><t>%s</t></is></c>`,
				string(rune('A'+colIndex)), rowIndex, header)
		}
		xmlContent += `</row>`
		rowIndex++

		// Данные
		if shipments, ok := data.([]map[string]interface{}); ok {
			for _, shipment := range shipments {
				xmlContent += fmt.Sprintf(`<row r="%d">`, rowIndex)
				values := []interface{}{
					shipment["id"], shipment["type"], shipment["tracking_number"],
					shipment["status"], shipment["created_at"], shipment["delivered_at"],
					shipment["sender_city"], shipment["recipient_city"],
				}
				for colIndex, value := range values {
					xmlContent += fmt.Sprintf(`<c r="%s%d" t="inlineStr"><is><t>%v</t></is></c>`,
						string(rune('A'+colIndex)), rowIndex, value)
				}
				xmlContent += `</row>`
				rowIndex++
			}
		}

	case "problems":
		// Заголовки
		xmlContent += fmt.Sprintf(`<row r="%d">`, rowIndex)
		headers := []string{"ID", "Shipment ID", "Type", "Problem Type", "Severity", "Status", "Created At", "Resolved At"}
		for colIndex, header := range headers {
			xmlContent += fmt.Sprintf(`<c r="%s%d" t="inlineStr"><is><t>%s</t></is></c>`,
				string(rune('A'+colIndex)), rowIndex, header)
		}
		xmlContent += `</row>`
		rowIndex++

		// Данные
		if problems, ok := data.([]map[string]interface{}); ok {
			for _, problem := range problems {
				xmlContent += fmt.Sprintf(`<row r="%d">`, rowIndex)
				values := []interface{}{
					problem["id"], problem["shipment_id"], problem["shipment_type"],
					problem["problem_type"], problem["severity"], problem["status"],
					problem["created_at"], problem["resolved_at"],
				}
				for colIndex, value := range values {
					xmlContent += fmt.Sprintf(`<c r="%s%d" t="inlineStr"><is><t>%v</t></is></c>`,
						string(rune('A'+colIndex)), rowIndex, value)
				}
				xmlContent += `</row>`
				rowIndex++
			}
		}

	case "performance":
		// Заголовки
		xmlContent += fmt.Sprintf(`<row r="%d">`, rowIndex)
		headers := []string{"Date", "Total Shipments", "Delivered", "In Transit", "Problems", "Success Rate", "Avg Delivery Time (hours)"}
		for colIndex, header := range headers {
			xmlContent += fmt.Sprintf(`<c r="%s%d" t="inlineStr"><is><t>%s</t></is></c>`,
				string(rune('A'+colIndex)), rowIndex, header)
		}
		xmlContent += `</row>`
		rowIndex++

		// Данные
		if metrics, ok := data.([]map[string]interface{}); ok {
			for _, metric := range metrics {
				xmlContent += fmt.Sprintf(`<row r="%d">`, rowIndex)
				values := []interface{}{
					metric["date"], metric["total_shipments"], metric["delivered"],
					metric["in_transit"], metric["problems"], metric["success_rate"],
					metric["avg_delivery_time"],
				}
				for colIndex, value := range values {
					xmlContent += fmt.Sprintf(`<c r="%s%d" t="inlineStr"><is><t>%v</t></is></c>`,
						string(rune('A'+colIndex)), rowIndex, value)
				}
				xmlContent += `</row>`
				rowIndex++
			}
		}
	}

	xmlContent += `</sheetData>
</worksheet>`

	// Записываем в ответ
	if _, err := c.Write([]byte(xmlContent)); err != nil {
		return err
	}
	return nil
}

// GetCourierComparison godoc
// @Summary Сравнение курьерских служб
// @Description Возвращает сравнительную аналитику по курьерским службам
// @Tags admin-logistics
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param date_from query string false "Дата начала периода (YYYY-MM-DD)"
// @Param date_to query string false "Дата окончания периода (YYYY-MM-DD)"
// @Success 200 {object} utils.SuccessResponseSwag{data=[]backend_internal_domain_logistics.CourierStats} "Courier comparison"
// @Failure 401 {object} utils.ErrorResponseSwag "Unauthorized"
// @Failure 500 {object} utils.ErrorResponseSwag "Internal server error"
// @Router /api/v1/admin/logistics/analytics/couriers [get]
func (h *AnalyticsHandler) GetCourierComparison(c *fiber.Ctx) error {
	// Проверка прав доступа
	userID := c.Locals("user_id")
	if userID == nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "auth.unauthorized")
	}

	// Парсим параметры
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")
	periodDays := c.QueryInt("period", 30)

	// Если даты не указаны, берем период в днях
	var fromDate, toDate time.Time
	var err error

	if dateFrom == "" {
		fromDate = time.Now().AddDate(0, 0, -periodDays)
	} else {
		fromDate, err = time.Parse("2006-01-02", dateFrom)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "logistics.invalid_date_from")
		}
	}

	if dateTo == "" {
		toDate = time.Now()
	} else {
		toDate, err = time.Parse("2006-01-02", dateTo)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "logistics.invalid_date_to")
		}
	}

	// Получаем сравнение курьеров
	comparison, err := h.analyticsService.GetCourierComparison(c.Context(), fromDate, toDate)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "logistics.couriers.error")
	}

	return utils.SuccessResponse(c, comparison)
}
