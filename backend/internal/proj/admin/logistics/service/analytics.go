package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"backend/internal/domain/logistics"
	"backend/pkg/logger"
)

// AnalyticsService сервис для аналитики логистики
type AnalyticsService struct {
	db     *sql.DB
	logger *logger.Logger
}

// NewAnalyticsService создает новый сервис аналитики
func NewAnalyticsService(db *sql.DB) *AnalyticsService {
	return &AnalyticsService{
		db:     db,
		logger: logger.GetLogger(),
	}
}

// GetPerformanceMetrics получает метрики производительности
func (s *AnalyticsService) GetPerformanceMetrics(ctx context.Context, fromDate, toDate time.Time, groupBy string) (map[string]interface{}, error) {
	// Определяем формат группировки
	var dateFormat string
	switch groupBy {
	case "week":
		dateFormat = "YYYY-WW"
	case "month":
		dateFormat = "YYYY-MM"
	default:
		dateFormat = "YYYY-MM-DD"
	}

	// Запрос для получения метрик
	//nolint:gosec // dateFormat is controlled, not user input
	query := fmt.Sprintf(`
		WITH shipment_data AS (
			SELECT 
				TO_CHAR(created_at, '%s') as period,
				COUNT(*) as total_shipments,
				COUNT(*) FILTER (WHERE status = 'delivered') as delivered,
				COUNT(*) FILTER (WHERE status IN ('in_transit', 'picked_up')) as in_transit,
				COUNT(*) FILTER (WHERE status IN ('returned', 'cancelled')) as failed,
				AVG(EXTRACT(EPOCH FROM (delivered_at - created_at))/3600) FILTER (WHERE delivered_at IS NOT NULL) as avg_delivery_hours
			FROM (
				SELECT created_at, delivered_at, status 
				FROM bex_shipments 
				WHERE created_at BETWEEN $1 AND $2
				UNION ALL
				SELECT created_at, delivered_at, status 
				FROM post_express_shipments 
				WHERE created_at BETWEEN $1 AND $2
			) s
			GROUP BY TO_CHAR(created_at, '%s')
		),
		problem_data AS (
			SELECT 
				TO_CHAR(created_at, '%s') as period,
				COUNT(*) as problems,
				COUNT(*) FILTER (WHERE severity = 'critical') as critical_problems
			FROM problem_shipments
			WHERE created_at BETWEEN $1 AND $2
			GROUP BY TO_CHAR(created_at, '%s')
		)
		SELECT 
			sd.period,
			sd.total_shipments,
			sd.delivered,
			sd.in_transit,
			sd.failed,
			sd.avg_delivery_hours,
			COALESCE(pd.problems, 0) as problems,
			COALESCE(pd.critical_problems, 0) as critical_problems,
			CASE 
				WHEN sd.total_shipments > 0 
				THEN ROUND(sd.delivered * 100.0 / sd.total_shipments, 2)
				ELSE 0 
			END as success_rate
		FROM shipment_data sd
		LEFT JOIN problem_data pd ON sd.period = pd.period
		ORDER BY sd.period
	`, dateFormat, dateFormat, dateFormat, dateFormat)

	rows, err := s.db.QueryContext(ctx, query, fromDate, toDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get performance metrics: %w", err)
	}
	defer rows.Close()

	var metrics []map[string]interface{}
	var totalShipments, totalDelivered, totalProblems int
	var totalDeliveryTime float64
	deliveryCount := 0

	for rows.Next() {
		var m struct {
			Period           string
			TotalShipments   int
			Delivered        int
			InTransit        int
			Failed           int
			AvgDeliveryHours sql.NullFloat64
			Problems         int
			CriticalProblems int
			SuccessRate      float64
		}

		err := rows.Scan(
			&m.Period,
			&m.TotalShipments,
			&m.Delivered,
			&m.InTransit,
			&m.Failed,
			&m.AvgDeliveryHours,
			&m.Problems,
			&m.CriticalProblems,
			&m.SuccessRate,
		)
		if err != nil {
			continue
		}

		// Агрегируем общие показатели
		totalShipments += m.TotalShipments
		totalDelivered += m.Delivered
		totalProblems += m.Problems
		if m.AvgDeliveryHours.Valid {
			totalDeliveryTime += m.AvgDeliveryHours.Float64
			deliveryCount++
		}

		metrics = append(metrics, map[string]interface{}{
			"period":            m.Period,
			"total_shipments":   m.TotalShipments,
			"delivered":         m.Delivered,
			"in_transit":        m.InTransit,
			"failed":            m.Failed,
			"avg_delivery_time": m.AvgDeliveryHours.Float64,
			"problems":          m.Problems,
			"critical_problems": m.CriticalProblems,
			"success_rate":      m.SuccessRate,
		})
	}

	// Вычисляем общие показатели
	avgDeliveryTime := float64(0)
	if deliveryCount > 0 {
		avgDeliveryTime = totalDeliveryTime / float64(deliveryCount)
	}

	successRate := float64(0)
	if totalShipments > 0 {
		successRate = float64(totalDelivered) * 100 / float64(totalShipments)
	}

	return map[string]interface{}{
		"metrics": metrics,
		"summary": map[string]interface{}{
			"total_shipments":      totalShipments,
			"total_delivered":      totalDelivered,
			"total_problems":       totalProblems,
			"avg_delivery_time":    avgDeliveryTime,
			"overall_success_rate": successRate,
		},
	}, rows.Err()
}

// GetFinancialReport получает финансовый отчет
func (s *AnalyticsService) GetFinancialReport(ctx context.Context, fromDate, toDate time.Time) (map[string]interface{}, error) {
	report := make(map[string]interface{})

	// Получаем стоимость доставок BEX Express
	var bexCosts struct {
		TotalCost     sql.NullFloat64
		TotalCOD      sql.NullFloat64
		ShipmentCount int
	}

	err := s.db.QueryRowContext(ctx, `
		SELECT 
			SUM(shipping_cost) as total_cost,
			SUM(cod_amount) as total_cod,
			COUNT(*) as shipment_count
		FROM bex_shipments
		WHERE created_at BETWEEN $1 AND $2
	`, fromDate, toDate).Scan(&bexCosts.TotalCost, &bexCosts.TotalCOD, &bexCosts.ShipmentCount)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to get BEX costs: %w", err)
	}

	// Получаем стоимость доставок Post Express
	var postCosts struct {
		TotalCost     sql.NullFloat64
		TotalCOD      sql.NullFloat64
		ShipmentCount int
	}

	err = s.db.QueryRowContext(ctx, `
		SELECT 
			SUM(price) as total_cost,
			SUM(cod_amount) as total_cod,
			COUNT(*) as shipment_count
		FROM post_express_shipments
		WHERE created_at BETWEEN $1 AND $2
	`, fromDate, toDate).Scan(&postCosts.TotalCost, &postCosts.TotalCOD, &postCosts.ShipmentCount)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to get Post Express costs: %w", err)
	}

	// Получаем стоимость возвратов
	var returnCosts sql.NullFloat64
	if err := s.db.QueryRowContext(ctx, `
		SELECT SUM(shipping_cost * 0.5) -- Предполагаем 50% от стоимости доставки за возврат
		FROM (
			SELECT shipping_cost FROM bex_shipments 
			WHERE status = 'returned' AND created_at BETWEEN $1 AND $2
			UNION ALL
			SELECT price as shipping_cost FROM post_express_shipments 
			WHERE status = 'returned' AND created_at BETWEEN $1 AND $2
		) returns
	`, fromDate, toDate).Scan(&returnCosts); err != nil {
		s.logger.Error("Failed to get return costs: %v", err)
		// Продолжаем с нулевым значением
	}

	// Формируем отчет
	totalShippingCost := bexCosts.TotalCost.Float64 + postCosts.TotalCost.Float64
	totalCODCollected := bexCosts.TotalCOD.Float64 + postCosts.TotalCOD.Float64
	totalShipments := bexCosts.ShipmentCount + postCosts.ShipmentCount

	report["summary"] = map[string]interface{}{
		"total_shipping_cost": totalShippingCost,
		"total_cod_collected": totalCODCollected,
		"total_return_cost":   returnCosts.Float64,
		"total_shipments":     totalShipments,
		"avg_cost_per_shipment": func() float64 {
			if totalShipments > 0 {
				return totalShippingCost / float64(totalShipments)
			}
			return 0
		}(),
	}

	report["by_courier"] = map[string]interface{}{
		"bex_express": map[string]interface{}{
			"total_cost": bexCosts.TotalCost.Float64,
			"total_cod":  bexCosts.TotalCOD.Float64,
			"shipments":  bexCosts.ShipmentCount,
			"avg_cost": func() float64 {
				if bexCosts.ShipmentCount > 0 {
					return bexCosts.TotalCost.Float64 / float64(bexCosts.ShipmentCount)
				}
				return 0
			}(),
		},
		"post_express": map[string]interface{}{
			"total_cost": postCosts.TotalCost.Float64,
			"total_cod":  postCosts.TotalCOD.Float64,
			"shipments":  postCosts.ShipmentCount,
			"avg_cost": func() float64 {
				if postCosts.ShipmentCount > 0 {
					return postCosts.TotalCost.Float64 / float64(postCosts.ShipmentCount)
				}
				return 0
			}(),
		},
	}

	// Получаем разбивку по месяцам
	monthlyQuery := `
		SELECT 
			TO_CHAR(created_at, 'YYYY-MM') as month,
			SUM(cost) as total_cost,
			COUNT(*) as shipments
		FROM (
			SELECT created_at, shipping_cost as cost FROM bex_shipments 
			WHERE created_at BETWEEN $1 AND $2
			UNION ALL
			SELECT created_at, price as cost FROM post_express_shipments 
			WHERE created_at BETWEEN $1 AND $2
		) s
		GROUP BY TO_CHAR(created_at, 'YYYY-MM')
		ORDER BY month
	`

	rows, err := s.db.QueryContext(ctx, monthlyQuery, fromDate, toDate)
	if err != nil {
		return report, nil // Возвращаем частичный отчет
	}
	defer rows.Close()

	var monthly []map[string]interface{}
	for rows.Next() {
		var month string
		var totalCost sql.NullFloat64
		var shipments int

		if err := rows.Scan(&month, &totalCost, &shipments); err != nil {
			continue
		}

		monthly = append(monthly, map[string]interface{}{
			"month":      month,
			"total_cost": totalCost.Float64,
			"shipments":  shipments,
		})
	}
	if err = rows.Err(); err != nil {
		s.logger.Error("error iterating monthly data rows: %v", err)
	}

	report["monthly"] = monthly

	return report, nil
}

// GetReportData получает данные для экспорта отчета
func (s *AnalyticsService) GetReportData(ctx context.Context, reportType string, fromDate, toDate time.Time) (interface{}, error) {
	switch reportType {
	case "shipments":
		return s.getShipmentsReportData(ctx, fromDate, toDate)
	case "problems":
		return s.getProblemsReportData(ctx, fromDate, toDate)
	case "performance":
		metrics, err := s.GetPerformanceMetrics(ctx, fromDate, toDate, "day")
		if err != nil {
			return nil, err
		}
		return metrics["metrics"], nil
	default:
		return nil, fmt.Errorf("unknown report type: %s", reportType)
	}
}

// getShipmentsReportData получает данные отправлений для отчета
func (s *AnalyticsService) getShipmentsReportData(ctx context.Context, fromDate, toDate time.Time) ([]map[string]interface{}, error) {
	query := `
		SELECT 
			id, type, tracking_number, status, 
			created_at, delivered_at, 
			sender_city, recipient_city
		FROM (
			SELECT 
				id, 'bex' as type, tracking_number, status,
				created_at, delivered_at,
				sender_city, recipient_city
			FROM bex_shipments
			WHERE created_at BETWEEN $1 AND $2
			UNION ALL
			SELECT 
				id, 'postexpress' as type, tracking_number, status,
				created_at, delivered_at,
				sender_city, recipient_city
			FROM post_express_shipments
			WHERE created_at BETWEEN $1 AND $2
		) s
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []map[string]interface{}
	for rows.Next() {
		var id int
		var shipmentType, trackingNumber, status string
		var createdAt time.Time
		var deliveredAt sql.NullTime
		var senderCity, recipientCity sql.NullString

		err := rows.Scan(
			&id, &shipmentType, &trackingNumber, &status,
			&createdAt, &deliveredAt,
			&senderCity, &recipientCity,
		)
		if err != nil {
			continue
		}

		item := map[string]interface{}{
			"id":              id,
			"type":            shipmentType,
			"tracking_number": trackingNumber,
			"status":          status,
			"created_at":      createdAt.Format("2006-01-02 15:04:05"),
			"sender_city":     senderCity.String,
			"recipient_city":  recipientCity.String,
		}

		if deliveredAt.Valid {
			item["delivered_at"] = deliveredAt.Time.Format("2006-01-02 15:04:05")
		} else {
			item["delivered_at"] = ""
		}

		data = append(data, item)
	}

	return data, rows.Err()
}

// getProblemsReportData получает данные проблем для отчета
func (s *AnalyticsService) getProblemsReportData(ctx context.Context, fromDate, toDate time.Time) ([]map[string]interface{}, error) {
	query := `
		SELECT 
			id, shipment_id, shipment_type, problem_type,
			severity, status, created_at, resolved_at
		FROM problem_shipments
		WHERE created_at BETWEEN $1 AND $2
		ORDER BY created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []map[string]interface{}
	for rows.Next() {
		var id, shipmentID int
		var shipmentType, problemType, severity, status string
		var createdAt time.Time
		var resolvedAt sql.NullTime

		err := rows.Scan(
			&id, &shipmentID, &shipmentType, &problemType,
			&severity, &status, &createdAt, &resolvedAt,
		)
		if err != nil {
			continue
		}

		item := map[string]interface{}{
			"id":            id,
			"shipment_id":   shipmentID,
			"shipment_type": shipmentType,
			"problem_type":  problemType,
			"severity":      severity,
			"status":        status,
			"created_at":    createdAt.Format("2006-01-02 15:04:05"),
		}

		if resolvedAt.Valid {
			item["resolved_at"] = resolvedAt.Time.Format("2006-01-02 15:04:05")
		} else {
			item["resolved_at"] = ""
		}

		data = append(data, item)
	}

	return data, rows.Err()
}

// GetCourierComparison получает сравнение курьерских служб
func (s *AnalyticsService) GetCourierComparison(ctx context.Context, fromDate, toDate time.Time) ([]logistics.CourierStats, error) {
	var stats []logistics.CourierStats

	// Статистика BEX Express
	bexStats := logistics.CourierStats{Name: "BEX Express"}
	err := s.db.QueryRowContext(ctx, `
		SELECT 
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status = 'delivered') as delivered,
			COALESCE(AVG(EXTRACT(EPOCH FROM (delivered_at - created_at))/3600) FILTER (WHERE delivered_at IS NOT NULL), 0) as avg_hours
		FROM bex_shipments
		WHERE created_at BETWEEN $1 AND $2
	`, fromDate, toDate).Scan(&bexStats.Shipments, &bexStats.Delivered, &bexStats.AvgTime)

	if err == nil && bexStats.Shipments > 0 {
		bexStats.SuccessRate = float64(bexStats.Delivered) * 100 / float64(bexStats.Shipments)
		stats = append(stats, bexStats)
	}

	// Статистика Post Express
	postStats := logistics.CourierStats{Name: "Post Express"}
	err = s.db.QueryRowContext(ctx, `
		SELECT 
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status = 'delivered') as delivered,
			COALESCE(AVG(EXTRACT(EPOCH FROM (delivered_at - created_at))/3600) FILTER (WHERE delivered_at IS NOT NULL), 0) as avg_hours
		FROM post_express_shipments
		WHERE created_at BETWEEN $1 AND $2
	`, fromDate, toDate).Scan(&postStats.Shipments, &postStats.Delivered, &postStats.AvgTime)

	if err == nil && postStats.Shipments > 0 {
		postStats.SuccessRate = float64(postStats.Delivered) * 100 / float64(postStats.Shipments)
		stats = append(stats, postStats)
	}

	return stats, nil
}
