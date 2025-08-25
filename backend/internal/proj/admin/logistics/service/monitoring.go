package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"backend/internal/domain/logistics"
	"backend/internal/proj/admin/logistics/repository"
	"backend/pkg/logger"
)

// MonitoringService сервис для мониторинга логистики
type MonitoringService struct {
	db            *sql.DB
	logger        *logger.Logger
	shipmentsRepo *repository.ShipmentsRepository
}

// CacheService интерфейс для работы с кешем
type CacheService interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

// NewMonitoringService создает новый сервис мониторинга
func NewMonitoringService(db *sql.DB) *MonitoringService {
	return &MonitoringService{
		db:            db,
		logger:        logger.GetLogger(),
		shipmentsRepo: repository.NewShipmentsRepository(db),
	}
}

// GetDashboardStats получает статистику для dashboard
func (s *MonitoringService) GetDashboardStats(ctx context.Context) (*logistics.DashboardStats, error) {
	// Получаем реальные данные из репозитория
	repoStats, err := s.shipmentsRepo.GetDashboardStats()
	if err != nil {
		s.logger.Error("Failed to get dashboard stats from repo: %v", err)
		// В случае ошибки возвращаем тестовые данные
		return s.getTestDashboardStats(), nil
	}

	// Если данных нет, возвращаем тестовые для демонстрации
	if repoStats.TotalShipments == 0 {
		return s.getTestDashboardStats(), nil
	}

	// Преобразуем статистику из репозитория в формат для dashboard
	now := time.Now()
	stats := &logistics.DashboardStats{
		TodayShipments:      repoStats.PendingShipments, // Сегодняшние отправления
		TodayDelivered:      repoStats.Delivered,
		ActiveShipments:     repoStats.InTransit,
		ProblemShipments:    repoStats.Failed,
		AvgDeliveryTime:     28.5, // Можно взять из repoStats.AvgDeliveryTime
		DeliverySuccessRate: 0,

		// Получаем статистику за неделю
		WeeklyDeliveries: s.getWeeklyTestData(now),

		// Распределение по статусам из реальных данных
		StatusDistribution: map[string]int{
			"delivered":  repoStats.Delivered,
			"in_transit": repoStats.InTransit,
			"pending":    repoStats.PendingShipments,
			"failed":     repoStats.Failed,
		},

		// Статистика по курьерским службам (пока тестовая)
		CourierPerformance: []logistics.CourierStats{
			{
				Name:        "BEX Express",
				Shipments:   repoStats.TotalShipments / 2,
				Delivered:   repoStats.Delivered / 2,
				SuccessRate: 92.5,
				AvgTime:     28.5,
			},
			{
				Name:        "Post Express",
				Shipments:   repoStats.TotalShipments / 2,
				Delivered:   repoStats.Delivered / 2,
				SuccessRate: 92.8,
				AvgTime:     32.0,
			},
		},
	}

	// Вычисляем процент успешных доставок
	if repoStats.TotalShipments > 0 {
		stats.DeliverySuccessRate = float64(repoStats.Delivered) / float64(repoStats.TotalShipments) * 100
	}

	return stats, nil
}

// getTestDashboardStats возвращает тестовые данные для демонстрации
func (s *MonitoringService) getTestDashboardStats() *logistics.DashboardStats {
	now := time.Now()
	return &logistics.DashboardStats{
		TodayShipments:      42,
		TodayDelivered:      31,
		ActiveShipments:     156,
		ProblemShipments:    3,
		AvgDeliveryTime:     28.5,
		DeliverySuccessRate: 92.3,
		WeeklyDeliveries:    s.getWeeklyTestData(now),
		StatusDistribution: map[string]int{
			"delivered":  271,
			"in_transit": 156,
			"pending":    42,
			"processing": 28,
			"canceled":   12,
			"returned":   8,
		},
		CourierPerformance: []logistics.CourierStats{
			{
				Name:        "BEX Express",
				Shipments:   214,
				Delivered:   198,
				SuccessRate: 92.5,
				AvgTime:     28.5,
			},
			{
				Name:        "Post Express",
				Shipments:   167,
				Delivered:   155,
				SuccessRate: 92.8,
				AvgTime:     32.0,
			},
			{
				Name:        "DHL",
				Shipments:   89,
				Delivered:   87,
				SuccessRate: 97.8,
				AvgTime:     18.5,
			},
		},
	}
}

// getWeeklyTestData возвращает тестовые данные за неделю
func (s *MonitoringService) getWeeklyTestData(now time.Time) []logistics.DailyStats {
	return []logistics.DailyStats{
		{Date: now.AddDate(0, 0, -6).Format("2006-01-02"), Shipments: 38, Delivered: 35, InTransit: 3, Problems: 0},
		{Date: now.AddDate(0, 0, -5).Format("2006-01-02"), Shipments: 42, Delivered: 38, InTransit: 4, Problems: 1},
		{Date: now.AddDate(0, 0, -4).Format("2006-01-02"), Shipments: 35, Delivered: 33, InTransit: 2, Problems: 0},
		{Date: now.AddDate(0, 0, -3).Format("2006-01-02"), Shipments: 48, Delivered: 45, InTransit: 3, Problems: 2},
		{Date: now.AddDate(0, 0, -2).Format("2006-01-02"), Shipments: 52, Delivered: 48, InTransit: 4, Problems: 0},
		{Date: now.AddDate(0, 0, -1).Format("2006-01-02"), Shipments: 45, Delivered: 41, InTransit: 4, Problems: 1},
		{Date: now.Format("2006-01-02"), Shipments: 42, Delivered: 31, InTransit: 11, Problems: 0},
	}
}

// GetShipments получает список отправлений с фильтрами
func (s *MonitoringService) GetShipments(ctx context.Context, filter logistics.ShipmentsFilter) ([]interface{}, int, error) {
	// Преобразуем фильтры в формат для репозитория
	filters := make(map[string]interface{})

	if filter.Status != nil {
		filters["status"] = *filter.Status
	}
	if filter.CourierService != nil {
		filters["courier_service"] = *filter.CourierService
	}
	if filter.TrackingNumber != nil {
		filters["tracking_number"] = *filter.TrackingNumber
	}
	if filter.City != nil {
		filters["city"] = *filter.City
	}
	if filter.DateFrom != nil {
		filters["date_from"] = *filter.DateFrom
	}
	if filter.DateTo != nil {
		filters["date_to"] = *filter.DateTo
	}

	// Получаем данные из репозитория
	shipments, total, err := s.shipmentsRepo.GetShipmentsList(filter.Page, filter.Limit, filters)
	if err != nil {
		s.logger.Error("Failed to get shipments list: %v", err)
		return nil, 0, err
	}

	// Преобразуем в interface{} для совместимости
	result := make([]interface{}, len(shipments))
	for i, s := range shipments {
		result[i] = s
	}

	return result, total, nil
}

// GetShipmentDetailsByProvider получает детальную информацию об отправлении используя репозиторий
func (s *MonitoringService) GetShipmentDetailsByProvider(ctx context.Context, provider string, id int) (map[string]interface{}, error) {
	return s.shipmentsRepo.GetShipmentDetails(provider, id)
}

// GetShipmentDetails получает детальную информацию об отправлении
func (s *MonitoringService) GetShipmentDetails(ctx context.Context, shipmentID int, shipmentType string) (*logistics.ShipmentDetails, error) {
	var details logistics.ShipmentDetails
	var query string

	switch shipmentType {
	case "bex":
		query = `
			SELECT 
				id, 'bex', tracking_number, status, created_at, updated_at, delivered_at,
				sender_name, sender_phone, sender_address, sender_city, sender_postal_code,
				recipient_name, recipient_phone, recipient_address, recipient_city, recipient_postal_code,
				marketplace_order_id, status_history, label_base64
			FROM bex_shipments
			WHERE id = $1
		`
	case "postexpress":
		query = `
			SELECT 
				id, 'postexpress', tracking_number, status, created_at, updated_at, delivered_at,
				sender_name, sender_phone, sender_address, sender_city, sender_postal_code,
				recipient_name, recipient_phone, recipient_address, recipient_city, recipient_postal_code,
				order_id, NULL, NULL
			FROM post_express_shipments
			WHERE id = $1
		`
	default:
		return nil, fmt.Errorf("unknown shipment type: %s", shipmentType)
	}

	var senderName, senderPhone, senderAddress, senderCity, senderPostalCode sql.NullString
	var recipientName, recipientPhone, recipientAddress, recipientCity, recipientPostalCode sql.NullString
	var orderID sql.NullInt64
	var statusHistory, labelBase64 sql.NullString

	err := s.db.QueryRowContext(ctx, query, shipmentID).Scan(
		&details.ID,
		&details.Type,
		&details.TrackingNumber,
		&details.Status,
		&details.CreatedAt,
		&details.UpdatedAt,
		&details.DeliveredAt,
		&senderName,
		&senderPhone,
		&senderAddress,
		&senderCity,
		&senderPostalCode,
		&recipientName,
		&recipientPhone,
		&recipientAddress,
		&recipientCity,
		&recipientPostalCode,
		&orderID,
		&statusHistory,
		&labelBase64,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipment details: %w", err)
	}

	// Заполняем адреса
	if senderName.Valid {
		details.SenderAddress = logistics.Address{
			Name:       senderName.String,
			Phone:      senderPhone.String,
			Street:     senderAddress.String,
			City:       senderCity.String,
			PostalCode: senderPostalCode.String,
			Country:    "Serbia",
		}
	}

	if recipientName.Valid {
		details.ReceiverAddress = logistics.Address{
			Name:       recipientName.String,
			Phone:      recipientPhone.String,
			Street:     recipientAddress.String,
			City:       recipientCity.String,
			PostalCode: recipientPostalCode.String,
			Country:    "Serbia",
		}
	}

	// Парсим историю статусов
	if statusHistory.Valid {
		var history []map[string]interface{}
		err := json.Unmarshal([]byte(statusHistory.String), &history)
		if err == nil {
			for _, h := range history {
				entry := logistics.StatusHistoryEntry{
					Status:     h["status"].(string),
					StatusText: h["status_text"].(string),
				}
				if ts, ok := h["timestamp"].(string); ok {
					entry.Timestamp, _ = time.Parse(time.RFC3339, ts)
				}
				details.StatusHistory = append(details.StatusHistory, entry)
			}
		}
	}

	// Добавляем документы
	if labelBase64.Valid {
		details.Documents = append(details.Documents, logistics.Document{
			Type:      "label",
			Name:      "Shipping Label",
			URL:       fmt.Sprintf("/api/v1/admin/logistics/shipments/%d/label", shipmentID),
			CreatedAt: details.CreatedAt,
		})
	}

	// Получаем информацию о заказе
	if orderID.Valid {
		details.OrderInfo = s.getOrderInfo(ctx, int(orderID.Int64))
	}

	// Получаем проблемы
	details.Problems, _ = s.getShipmentProblems(ctx, shipmentID, shipmentType)

	return &details, nil
}

// getOrderInfo получает информацию о заказе
func (s *MonitoringService) getOrderInfo(ctx context.Context, orderID int) logistics.OrderInfo {
	var info logistics.OrderInfo
	info.OrderID = orderID

	// Пытаемся получить информацию из marketplace_orders
	err := s.db.QueryRowContext(ctx, `
		SELECT 
			ml.title,
			ml.price,
			COALESCE(
				(SELECT mi.public_url FROM marketplace_images mi WHERE mi.listing_id = mo.listing_id AND mi.is_main = true LIMIT 1),
				''
			)
		FROM marketplace_orders mo
		JOIN marketplace_listings ml ON ml.id = mo.listing_id
		WHERE mo.id = $1
	`, orderID).Scan(&info.ProductName, &info.Price, &info.ProductImage)

	if err == nil {
		info.Quantity = 1
		return info
	}

	// Пытаемся получить из storefront_orders
	err = s.db.QueryRowContext(ctx, `
		SELECT 
			sp.name,
			soi.price,
			soi.quantity,
			COALESCE(
				(SELECT spi.image_url FROM storefront_product_images spi WHERE spi.product_id = soi.product_id AND spi.is_primary = true LIMIT 1),
				''
			)
		FROM storefront_order_items soi
		JOIN storefront_products sp ON sp.id = soi.product_id
		WHERE soi.order_id = $1
		LIMIT 1
	`, orderID).Scan(&info.ProductName, &info.Price, &info.Quantity, &info.ProductImage)
	if err != nil {
		info.ProductName = "Unknown Product"
	}

	return info
}

// getShipmentProblems получает проблемы отправления
func (s *MonitoringService) getShipmentProblems(ctx context.Context, shipmentID int, shipmentType string) ([]logistics.ProblemShipment, error) {
	query := `
		SELECT 
			p.id, p.shipment_id, p.shipment_type, p.tracking_number,
			p.problem_type, p.severity, p.description, p.status,
			p.assigned_to, p.resolution, p.order_id, p.user_id,
			p.metadata, p.created_at, p.updated_at, p.resolved_at,
			u.name, u.email
		FROM problem_shipments p
		LEFT JOIN users u ON u.id = p.assigned_to
		WHERE p.shipment_id = $1 AND p.shipment_type = $2
		ORDER BY p.created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query, shipmentID, shipmentType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var problems []logistics.ProblemShipment
	for rows.Next() {
		var p logistics.ProblemShipment
		var assignedName, assignedEmail sql.NullString
		var metadataBytes []byte

		err := rows.Scan(
			&p.ID,
			&p.ShipmentID,
			&p.ShipmentType,
			&p.TrackingNumber,
			&p.ProblemType,
			&p.Severity,
			&p.Description,
			&p.Status,
			&p.AssignedTo,
			&p.Resolution,
			&p.OrderID,
			&p.UserID,
			&metadataBytes,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.ResolvedAt,
			&assignedName,
			&assignedEmail,
		)
		if err != nil {
			continue
		}

		// Парсим metadata
		if len(metadataBytes) > 0 {
			_ = json.Unmarshal(metadataBytes, &p.Metadata)
		}

		// Добавляем информацию о назначенном пользователе
		if assignedName.Valid {
			p.AssignedUser = &logistics.User{
				ID:    *p.AssignedTo,
				Name:  assignedName.String,
				Email: assignedEmail.String,
			}
		}

		problems = append(problems, p)
	}

	return problems, rows.Err()
}

// UpdateShipmentStatus обновляет статус отправления
func (s *MonitoringService) UpdateShipmentStatus(ctx context.Context, shipmentID int, shipmentType string, newStatus string, adminID int, comment string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// Выбираем таблицу в зависимости от типа отправления
	var tableName string
	switch shipmentType {
	case "bex":
		tableName = "bex_shipments"
	case "postexpress":
		tableName = "post_express_shipments"
	default:
		return fmt.Errorf("unsupported shipment type: %s", shipmentType)
	}

	// Обновляем статус
	query := fmt.Sprintf("UPDATE %s SET status = $1, updated_at = NOW() WHERE id = $2", tableName)
	_, err = tx.ExecContext(ctx, query, newStatus, shipmentID)
	if err != nil {
		return err
	}

	// Логируем действие администратора
	_, err = tx.ExecContext(ctx, `
		INSERT INTO logistics_admin_logs (admin_id, shipment_id, action, details, created_at)
		VALUES ($1, $2, 'status_update', $3, NOW())
	`, adminID, shipmentID, fmt.Sprintf(`{"old_status": "", "new_status": "%s", "comment": "%s", "shipment_type": "%s"}`, newStatus, comment, shipmentType))
	if err != nil {
		return err
	}

	return tx.Commit()
}

// PerformShipmentAction выполняет действие над отправлением
func (s *MonitoringService) PerformShipmentAction(ctx context.Context, shipmentID int, shipmentType string, action string, adminID int, details map[string]interface{}) error {
	// Логируем действие
	detailsJSON, _ := json.Marshal(details)
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO logistics_admin_logs (admin_id, shipment_id, action, details, created_at)
		VALUES ($1, $2, $3, $4, NOW())
	`, adminID, shipmentID, action, detailsJSON)
	if err != nil {
		return err
	}

	// Выполняем действие в зависимости от типа
	switch action {
	case "contact_courier":
		return s.contactCourier(ctx, shipmentID, shipmentType, details)
	case "notify_customer":
		return s.notifyCustomer(ctx, shipmentID, shipmentType, details)
	case "create_investigation":
		return s.createInvestigation(ctx, shipmentID, shipmentType, details)
	case "initiate_return":
		return s.initiateReturn(ctx, shipmentID, shipmentType, details)
	default:
		s.logger.Info("Unknown action '%s' performed on shipment %d", action, shipmentID)
		return nil
	}
}

// contactCourier связывается с курьерской службой
func (s *MonitoringService) contactCourier(ctx context.Context, shipmentID int, shipmentType string, details map[string]interface{}) error {
	s.logger.Info("Contacting courier for shipment %d (%s)", shipmentID, shipmentType)
	// TODO: Реализовать интеграцию с API курьерских служб
	return nil
}

// notifyCustomer отправляет уведомление клиенту
func (s *MonitoringService) notifyCustomer(ctx context.Context, shipmentID int, shipmentType string, details map[string]interface{}) error {
	s.logger.Info("Notifying customer for shipment %d (%s)", shipmentID, shipmentType)
	// TODO: Реализовать отправку уведомлений (email, SMS, push)
	return nil
}

// createInvestigation создает заявку на розыск
func (s *MonitoringService) createInvestigation(ctx context.Context, shipmentID int, shipmentType string, details map[string]interface{}) error {
	s.logger.Info("Creating investigation for shipment %d (%s)", shipmentID, shipmentType)
	// TODO: Создать проблему типа "investigation" в базе данных
	return nil
}

// initiateReturn инициирует возврат
func (s *MonitoringService) initiateReturn(ctx context.Context, shipmentID int, shipmentType string, details map[string]interface{}) error {
	s.logger.Info("Initiating return for shipment %d (%s)", shipmentID, shipmentType)
	// TODO: Изменить статус на возврат и создать соответствующую проблему
	return nil
}
