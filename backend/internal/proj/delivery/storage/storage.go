package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"backend/internal/proj/delivery/models"
)

// Storage - хранилище для работы с БД
type Storage struct {
	db *sqlx.DB
}

// NewStorage создает новый экземпляр хранилища
func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{
		db: db,
	}
}

// GetProvider - получает провайдера по ID
func (s *Storage) GetProvider(ctx context.Context, id int) (*models.Provider, error) {
	var provider models.Provider
	query := `SELECT * FROM delivery_providers WHERE id = $1`

	if err := s.db.GetContext(ctx, &provider, query, id); err != nil {
		return nil, err
	}

	return &provider, nil
}

// GetProviderByCode - получает провайдера по коду
func (s *Storage) GetProviderByCode(ctx context.Context, code string) (*models.Provider, error) {
	var provider models.Provider
	query := `SELECT * FROM delivery_providers WHERE code = $1`

	if err := s.db.GetContext(ctx, &provider, query, code); err != nil {
		return nil, err
	}

	return &provider, nil
}

// GetProviders - получает список провайдеров
func (s *Storage) GetProviders(ctx context.Context, activeOnly bool) ([]models.Provider, error) {
	var providers []models.Provider
	query := `SELECT * FROM delivery_providers`

	if activeOnly {
		query += " WHERE is_active = true"
	}

	query += " ORDER BY id"

	if err := s.db.SelectContext(ctx, &providers, query); err != nil {
		return nil, err
	}

	return providers, nil
}

// CreateShipment - создает новое отправление
func (s *Storage) CreateShipment(ctx context.Context, shipment *models.Shipment) error {
	query := `
		INSERT INTO delivery_shipments (
			provider_id, order_id, external_id, tracking_number, status,
			sender_info, recipient_info, package_info,
			delivery_cost, insurance_cost, cod_amount,
			cost_breakdown, pickup_date, estimated_delivery,
			provider_response, labels
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		) RETURNING id, created_at, updated_at`

	err := s.db.GetContext(ctx, shipment, query,
		shipment.ProviderID,
		shipment.OrderID,
		shipment.ExternalID,
		shipment.TrackingNumber,
		shipment.Status,
		shipment.SenderInfo,
		shipment.RecipientInfo,
		shipment.PackageInfo,
		shipment.DeliveryCost,
		shipment.InsuranceCost,
		shipment.CODAmount,
		shipment.CostBreakdown,
		shipment.PickupDate,
		shipment.EstimatedDelivery,
		shipment.ProviderResponse,
		shipment.Labels,
	)

	return err
}

// GetShipment - получает отправление по ID
func (s *Storage) GetShipment(ctx context.Context, id int) (*models.Shipment, error) {
	var shipment models.Shipment
	query := `SELECT * FROM delivery_shipments WHERE id = $1`

	if err := s.db.GetContext(ctx, &shipment, query, id); err != nil {
		return nil, err
	}

	return &shipment, nil
}

// GetShipmentByTracking - получает отправление по трек-номеру
func (s *Storage) GetShipmentByTracking(ctx context.Context, trackingNumber string) (*models.Shipment, error) {
	var shipment models.Shipment
	query := `SELECT * FROM delivery_shipments WHERE tracking_number = $1`

	if err := s.db.GetContext(ctx, &shipment, query, trackingNumber); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("shipment not found")
		}
		return nil, err
	}

	return &shipment, nil
}

// UpdateShipmentStatus - обновляет статус отправления
func (s *Storage) UpdateShipmentStatus(ctx context.Context, id int, status string, deliveredAt *time.Time) error {
	query := `
		UPDATE delivery_shipments
		SET status = $1, actual_delivery_date = $2, updated_at = NOW()
		WHERE id = $3`

	_, err := s.db.ExecContext(ctx, query, status, deliveredAt, id)
	return err
}

// UpdateOrderShipment - связывает заказ с отправлением
func (s *Storage) UpdateOrderShipment(ctx context.Context, orderID, shipmentID int) error {
	query := `
		UPDATE marketplace_orders
		SET delivery_shipment_id = $1
		WHERE id = $2`

	_, err := s.db.ExecContext(ctx, query, shipmentID, orderID)
	return err
}

// CreateTrackingEvent - создает событие отслеживания
func (s *Storage) CreateTrackingEvent(ctx context.Context, event *models.TrackingEvent) error {
	// Проверяем, не существует ли уже такое событие
	var exists bool
	checkQuery := `
		SELECT EXISTS(
			SELECT 1 FROM delivery_tracking_events
			WHERE shipment_id = $1 AND event_time = $2 AND status = $3
		)`

	if err := s.db.GetContext(ctx, &exists, checkQuery, event.ShipmentID, event.EventTime, event.Status); err != nil {
		return err
	}

	if exists {
		// Событие уже существует, пропускаем
		return nil
	}

	query := `
		INSERT INTO delivery_tracking_events (
			shipment_id, provider_id, event_time, status,
			location, description, raw_data
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at`

	err := s.db.GetContext(ctx, event, query,
		event.ShipmentID,
		event.ProviderID,
		event.EventTime,
		event.Status,
		event.Location,
		event.Description,
		event.RawData,
	)

	return err
}

// GetTrackingEvents - получает события отслеживания
func (s *Storage) GetTrackingEvents(ctx context.Context, shipmentID int) ([]models.TrackingEvent, error) {
	var events []models.TrackingEvent
	query := `
		SELECT * FROM delivery_tracking_events
		WHERE shipment_id = $1
		ORDER BY event_time DESC`

	if err := s.db.SelectContext(ctx, &events, query, shipmentID); err != nil {
		return nil, err
	}

	return events, nil
}

// GetPricingRules - получает правила расчета стоимости для провайдера
func (s *Storage) GetPricingRules(ctx context.Context, providerID int) ([]models.PricingRule, error) {
	var rules []models.PricingRule
	query := `
		SELECT * FROM delivery_pricing_rules
		WHERE provider_id = $1 AND is_active = true
		ORDER BY priority DESC`

	if err := s.db.SelectContext(ctx, &rules, query, providerID); err != nil {
		return nil, err
	}

	return rules, nil
}

// CreatePricingRule - создает правило расчета стоимости
func (s *Storage) CreatePricingRule(ctx context.Context, rule *models.PricingRule) error {
	query := `
		INSERT INTO delivery_pricing_rules (
			provider_id, rule_type,
			weight_ranges, volume_ranges, zone_multipliers,
			fragile_surcharge, oversized_surcharge, special_handling_surcharge,
			min_price, max_price, custom_formula,
			priority, is_active
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at, updated_at`

	err := s.db.GetContext(ctx, rule, query,
		rule.ProviderID,
		rule.RuleType,
		rule.WeightRanges,
		rule.VolumeRanges,
		rule.ZoneMultipliers,
		rule.FragileSurcharge,
		rule.OversizedSurcharge,
		rule.SpecialHandlingSurcharge,
		rule.MinPrice,
		rule.MaxPrice,
		rule.CustomFormula,
		rule.Priority,
		rule.IsActive,
	)

	return err
}

// GetZones - получает зоны доставки
func (s *Storage) GetZones(ctx context.Context) ([]models.Zone, error) {
	var zones []models.Zone
	query := `SELECT id, name, type, countries, regions, cities, postal_codes, radius_km, created_at FROM delivery_zones ORDER BY type, name`

	if err := s.db.SelectContext(ctx, &zones, query); err != nil {
		return nil, err
	}

	return zones, nil
}

// GetZoneByLocation - определяет зону для местоположения
func (s *Storage) GetZoneByLocation(ctx context.Context, country, city string) (*models.Zone, error) {
	var zone models.Zone
	query := `
		SELECT id, name, type, countries, regions, cities, postal_codes, radius_km, created_at
		FROM delivery_zones
		WHERE $1 = ANY(countries)
		AND ($2 = ANY(cities) OR cities IS NULL)
		ORDER BY
			CASE type
				WHEN 'local' THEN 1
				WHEN 'regional' THEN 2
				WHEN 'national' THEN 3
				WHEN 'international' THEN 4
			END
		LIMIT 1`

	if err := s.db.GetContext(ctx, &zone, query, country, city); err != nil {
		if err == sql.ErrNoRows {
			// Возвращаем дефолтную национальную зону
			zone.Type = models.ZoneTypeNational
			zone.Name = "Default National"
			return &zone, nil
		}
		return nil, err
	}

	return &zone, nil
}

// GetShipmentsByOrder - получает отправления для заказа
func (s *Storage) GetShipmentsByOrder(ctx context.Context, orderID int) ([]models.Shipment, error) {
	var shipments []models.Shipment
	query := `SELECT * FROM delivery_shipments WHERE order_id = $1 ORDER BY created_at DESC`

	if err := s.db.SelectContext(ctx, &shipments, query, orderID); err != nil {
		return nil, err
	}

	return shipments, nil
}

// GetShipmentsByStatus - получает отправления по статусу
func (s *Storage) GetShipmentsByStatus(ctx context.Context, status string, limit int) ([]models.Shipment, error) {
	var shipments []models.Shipment
	query := `SELECT * FROM delivery_shipments WHERE status = $1 ORDER BY updated_at DESC`

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	if err := s.db.SelectContext(ctx, &shipments, query, status); err != nil {
		return nil, err
	}

	return shipments, nil
}

// GetRecentShipments - получает последние отправления
func (s *Storage) GetRecentShipments(ctx context.Context, limit int) ([]models.Shipment, error) {
	var shipments []models.Shipment
	query := `SELECT * FROM delivery_shipments ORDER BY created_at DESC LIMIT $1`

	if err := s.db.SelectContext(ctx, &shipments, query, limit); err != nil {
		return nil, err
	}

	return shipments, nil
}

// GetShipmentStats - получает статистику по отправлениям
func (s *Storage) GetShipmentStats(ctx context.Context, providerID *int, from, to time.Time) (*ShipmentStats, error) {
	var stats ShipmentStats

	whereClause := "WHERE created_at BETWEEN $1 AND $2"
	args := []interface{}{from, to}

	if providerID != nil {
		whereClause += " AND provider_id = $3"
		args = append(args, *providerID)
	}

	// Общее количество
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM delivery_shipments %s`, whereClause)
	if err := s.db.GetContext(ctx, &stats.TotalCount, countQuery, args...); err != nil {
		return nil, err
	}

	// По статусам
	statusQuery := fmt.Sprintf(`
		SELECT status, COUNT(*) as count
		FROM delivery_shipments %s
		GROUP BY status`, whereClause)

	rows, err := s.db.QueryContext(ctx, statusQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats.ByStatus = make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		stats.ByStatus[status] = count
	}

	// Общая стоимость
	costQuery := fmt.Sprintf(`
		SELECT
			COALESCE(SUM(delivery_cost), 0) as total_delivery,
			COALESCE(SUM(insurance_cost), 0) as total_insurance,
			COALESCE(SUM(cod_amount), 0) as total_cod
		FROM delivery_shipments %s`, whereClause)

	if err := s.db.GetContext(ctx, &stats, costQuery, args...); err != nil {
		return nil, err
	}

	// Среднее время доставки
	avgQuery := fmt.Sprintf(`
		SELECT AVG(EXTRACT(DAY FROM (actual_delivery_date - created_at)))
		FROM delivery_shipments %s
		AND actual_delivery_date IS NOT NULL`, whereClause)

	var avgDays sql.NullFloat64
	if err := s.db.GetContext(ctx, &avgDays, avgQuery, args...); err != nil {
		return nil, err
	}
	if avgDays.Valid {
		stats.AvgDeliveryDays = avgDays.Float64
	}

	return &stats, nil
}

// ShipmentStats - статистика по отправлениям
type ShipmentStats struct {
	TotalCount       int            `json:"total_count"`
	ByStatus         map[string]int `json:"by_status"`
	TotalDelivery    float64        `json:"total_delivery" db:"total_delivery"`
	TotalInsurance   float64        `json:"total_insurance" db:"total_insurance"`
	TotalCOD         float64        `json:"total_cod" db:"total_cod"`
	AvgDeliveryDays  float64        `json:"avg_delivery_days"`
}

// UpdateProvider - обновляет провайдера
func (s *Storage) UpdateProvider(ctx context.Context, provider *models.Provider) error {
	query := `
		UPDATE delivery_providers
		SET name = $1, logo_url = $2, is_active = $3,
		    supports_cod = $4, supports_insurance = $5, supports_tracking = $6,
		    api_config = $7, capabilities = $8, updated_at = NOW()
		WHERE id = $9`

	_, err := s.db.ExecContext(ctx, query,
		provider.Name,
		provider.LogoURL,
		provider.IsActive,
		provider.SupportsCOD,
		provider.SupportsInsurance,
		provider.SupportsTracking,
		provider.APIConfig,
		provider.Capabilities,
		provider.ID,
	)
	return err
}

// GetShipmentStatistics - получает статистику отправлений
func (s *Storage) GetShipmentStatistics(ctx context.Context, from, to time.Time, providerID *int) (*ShipmentStatistics, error) {
	stats := &ShipmentStatistics{
		StatusBreakdown: make(map[string]int),
	}

	whereClause := "WHERE created_at BETWEEN $1 AND $2"
	args := []interface{}{from, to}

	if providerID != nil {
		whereClause += fmt.Sprintf(" AND provider_id = $%d", len(args)+1)
		args = append(args, *providerID)
	}

	// Общая статистика
	query := fmt.Sprintf(`
		SELECT
			COUNT(*) as total_shipments,
			COALESCE(SUM(delivery_cost), 0) as total_cost,
			COALESCE(AVG(delivery_cost), 0) as average_cost
		FROM delivery_shipments %s`, whereClause)

	if err := s.db.QueryRowContext(ctx, query, args...).Scan(
		&stats.TotalShipments,
		&stats.TotalCost,
		&stats.AverageCost,
	); err != nil {
		return nil, err
	}

	// Статистика по статусам
	statusQuery := fmt.Sprintf(`
		SELECT status, COUNT(*)
		FROM delivery_shipments %s
		GROUP BY status`, whereClause)

	rows, err := s.db.QueryContext(ctx, statusQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		stats.StatusBreakdown[status] = count
	}

	return stats, nil
}

// ShipmentStatistics - статистика отправлений
type ShipmentStatistics struct {
	TotalShipments  int            `json:"total_shipments"`
	TotalCost       float64        `json:"total_cost"`
	AverageCost     float64        `json:"average_cost"`
	StatusBreakdown map[string]int `json:"status_breakdown"`
}

// GetProviderStatistics - получает статистику по провайдерам
func (s *Storage) GetProviderStatistics(ctx context.Context, from, to time.Time) ([]ProviderStatistics, error) {
	query := `
		SELECT
			p.id as provider_id,
			p.name as provider_name,
			COUNT(s.id) as shipment_count,
			COALESCE(SUM(s.delivery_cost), 0) as total_cost,
			COALESCE(AVG(s.delivery_cost), 0) as average_cost,
			COALESCE(
				CAST(COUNT(CASE WHEN s.status = 'delivered' THEN 1 END) AS FLOAT) /
				NULLIF(COUNT(s.id), 0) * 100, 0
			) as success_rate
		FROM delivery_providers p
		LEFT JOIN delivery_shipments s ON p.id = s.provider_id
			AND s.created_at BETWEEN $1 AND $2
		GROUP BY p.id, p.name
		ORDER BY shipment_count DESC`

	var stats []ProviderStatistics
	if err := s.db.SelectContext(ctx, &stats, query, from, to); err != nil {
		return nil, err
	}

	return stats, nil
}

// ProviderStatistics - статистика по провайдеру
type ProviderStatistics struct {
	ProviderID    int     `json:"provider_id" db:"provider_id"`
	ProviderName  string  `json:"provider_name" db:"provider_name"`
	ShipmentCount int     `json:"shipment_count" db:"shipment_count"`
	TotalCost     float64 `json:"total_cost" db:"total_cost"`
	AverageCost   float64 `json:"average_cost" db:"average_cost"`
	SuccessRate   float64 `json:"success_rate" db:"success_rate"`
}

// GetTopDeliveryRoutes - получает топ маршрутов доставки
func (s *Storage) GetTopDeliveryRoutes(ctx context.Context, from, to time.Time, limit int) ([]RouteStatistics, error) {
	query := `
		SELECT
			sender_info->>'city' as from_city,
			recipient_info->>'city' as to_city,
			COUNT(*) as shipment_count,
			COALESCE(AVG(delivery_cost), 0) as average_cost,
			COALESCE(AVG(
				EXTRACT(EPOCH FROM (actual_delivery_date - created_at)) / 86400
			), 0) as average_days
		FROM delivery_shipments
		WHERE created_at BETWEEN $1 AND $2
			AND sender_info->>'city' IS NOT NULL
			AND recipient_info->>'city' IS NOT NULL
		GROUP BY sender_info->>'city', recipient_info->>'city'
		ORDER BY shipment_count DESC
		LIMIT $3`

	var routes []RouteStatistics
	if err := s.db.SelectContext(ctx, &routes, query, from, to, limit); err != nil {
		return nil, err
	}

	return routes, nil
}

// RouteStatistics - статистика маршрута
type RouteStatistics struct {
	FromCity      string  `json:"from_city" db:"from_city"`
	ToCity        string  `json:"to_city" db:"to_city"`
	ShipmentCount int     `json:"shipment_count" db:"shipment_count"`
	AverageCost   float64 `json:"average_cost" db:"average_cost"`
	AverageDays   float64 `json:"average_days" db:"average_days"`
}

// GetAverageDeliveryTimes - получает среднее время доставки
func (s *Storage) GetAverageDeliveryTimes(ctx context.Context, from, to time.Time, providerID *int) (map[string]float64, error) {
	whereClause := "WHERE created_at BETWEEN $1 AND $2 AND actual_delivery_date IS NOT NULL"
	args := []interface{}{from, to}

	if providerID != nil {
		whereClause += fmt.Sprintf(" AND provider_id = $%d", len(args)+1)
		args = append(args, *providerID)
	}

	query := fmt.Sprintf(`
		SELECT
			status,
			AVG(EXTRACT(EPOCH FROM (actual_delivery_date - created_at)) / 86400) as avg_days
		FROM delivery_shipments %s
		GROUP BY status`, whereClause)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]float64)
	for rows.Next() {
		var status string
		var avgDays float64
		if err := rows.Scan(&status, &avgDays); err != nil {
			return nil, err
		}
		result[status] = avgDays
	}

	return result, nil
}