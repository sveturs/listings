package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

type ShipmentsRepository struct {
	db *sql.DB
}

func NewShipmentsRepository(db *sql.DB) *ShipmentsRepository {
	return &ShipmentsRepository{db: db}
}

type ShipmentStats struct {
	TotalShipments   int     `json:"total_shipments"`
	PendingShipments int     `json:"pending_shipments"`
	InTransit        int     `json:"in_transit"`
	Delivered        int     `json:"delivered"`
	Failed           int     `json:"failed"`
	TotalRevenue     float64 `json:"total_revenue"`
	AvgDeliveryTime  string  `json:"avg_delivery_time"`
}

type BEXShipment struct {
	ID                 int             `json:"id"`
	MarketplaceOrderID *int            `json:"marketplace_order_id"`
	StorefrontOrderID  *int64          `json:"storefront_order_id"`
	TrackingNumber     string          `json:"tracking_number"`
	RecipientName      string          `json:"recipient_name"`
	RecipientCity      string          `json:"recipient_city"`
	RecipientPhone     string          `json:"recipient_phone"`
	Status             int             `json:"status"`
	StatusText         string          `json:"status_text"`
	CODAmount          float64         `json:"cod_amount"`
	WeightKg           float64         `json:"weight_kg"`
	RegisteredAt       *time.Time      `json:"registered_at"`
	DeliveredAt        *time.Time      `json:"delivered_at"`
	FailedReason       *string         `json:"failed_reason"`
	StatusHistory      json.RawMessage `json:"status_history"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

type PostExpressShipment struct {
	ID                 int             `json:"id"`
	MarketplaceOrderID *int            `json:"marketplace_order_id"`
	StorefrontOrderID  *int64          `json:"storefront_order_id"`
	TrackingNumber     string          `json:"tracking_number"`
	RecipientName      string          `json:"recipient_name"`
	RecipientCity      string          `json:"recipient_city"`
	RecipientPhone     string          `json:"recipient_phone"`
	Status             string          `json:"status"`
	DeliveryStatus     *string         `json:"delivery_status"`
	CODAmount          float64         `json:"cod_amount"`
	WeightKg           float64         `json:"weight_kg"`
	TotalPrice         float64         `json:"total_price"`
	RegisteredAt       *time.Time      `json:"registered_at"`
	DeliveredAt        *time.Time      `json:"delivered_at"`
	FailedReason       *string         `json:"failed_reason"`
	StatusHistory      json.RawMessage `json:"status_history"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

func (r *ShipmentsRepository) GetDashboardStats() (*ShipmentStats, error) {
	stats := &ShipmentStats{}

	// Get BEX stats
	bexQuery := `
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN status = 1 THEN 1 END) as pending,
			COUNT(CASE WHEN status = 2 THEN 1 END) as in_transit,
			COUNT(CASE WHEN status = 3 THEN 1 END) as delivered,
			COUNT(CASE WHEN status = 4 THEN 1 END) as failed,
			COALESCE(SUM(cod_amount), 0) as revenue,
			COALESCE(AVG(
				CASE 
					WHEN delivered_at IS NOT NULL AND registered_at IS NOT NULL 
					THEN EXTRACT(EPOCH FROM (delivered_at - registered_at))/3600
				END
			), 0) as avg_delivery_hours
		FROM bex_shipments
		WHERE created_at >= NOW() - INTERVAL '30 days'
	`

	var bexStats struct {
		Total            int
		Pending          int
		InTransit        int
		Delivered        int
		Failed           int
		Revenue          float64
		AvgDeliveryHours float64
	}

	err := r.db.QueryRow(bexQuery).Scan(
		&bexStats.Total,
		&bexStats.Pending,
		&bexStats.InTransit,
		&bexStats.Delivered,
		&bexStats.Failed,
		&bexStats.Revenue,
		&bexStats.AvgDeliveryHours,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get BEX stats")
	}

	// Get Post Express stats
	postQuery := `
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending,
			COUNT(CASE WHEN status = 'in_transit' THEN 1 END) as in_transit,
			COUNT(CASE WHEN status = 'delivered' THEN 1 END) as delivered,
			COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed,
			COALESCE(SUM(cod_amount), 0) as revenue,
			COALESCE(AVG(
				CASE 
					WHEN delivered_at IS NOT NULL AND registered_at IS NOT NULL 
					THEN EXTRACT(EPOCH FROM (delivered_at - registered_at))/3600
				END
			), 0) as avg_delivery_hours
		FROM post_express_shipments
		WHERE created_at >= NOW() - INTERVAL '30 days'
	`

	var postStats struct {
		Total            int
		Pending          int
		InTransit        int
		Delivered        int
		Failed           int
		Revenue          float64
		AvgDeliveryHours float64
	}

	err = r.db.QueryRow(postQuery).Scan(
		&postStats.Total,
		&postStats.Pending,
		&postStats.InTransit,
		&postStats.Delivered,
		&postStats.Failed,
		&postStats.Revenue,
		&postStats.AvgDeliveryHours,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Post Express stats")
	}

	// Combine stats
	stats.TotalShipments = bexStats.Total + postStats.Total
	stats.PendingShipments = bexStats.Pending + postStats.Pending
	stats.InTransit = bexStats.InTransit + postStats.InTransit
	stats.Delivered = bexStats.Delivered + postStats.Delivered
	stats.Failed = bexStats.Failed + postStats.Failed
	stats.TotalRevenue = bexStats.Revenue + postStats.Revenue

	// Calculate average delivery time
	totalDelivered := bexStats.Delivered + postStats.Delivered
	if totalDelivered > 0 {
		avgHours := (bexStats.AvgDeliveryHours*float64(bexStats.Delivered) +
			postStats.AvgDeliveryHours*float64(postStats.Delivered)) /
			float64(totalDelivered)

		if avgHours < 24 {
			stats.AvgDeliveryTime = fmt.Sprintf("%.1f часов", avgHours)
		} else {
			stats.AvgDeliveryTime = fmt.Sprintf("%.1f дней", avgHours/24)
		}
	} else {
		stats.AvgDeliveryTime = "Нет данных"
	}

	return stats, nil
}

func (r *ShipmentsRepository) GetRecentShipments(limit int) ([]interface{}, error) {
	if limit <= 0 {
		limit = 10
	}

	var shipments []interface{}

	// Get recent BEX shipments
	bexQuery := `
		SELECT 
			id, marketplace_order_id, storefront_order_id, tracking_number,
			recipient_name, recipient_city, recipient_phone,
			status, status_text, cod_amount, weight_kg,
			registered_at, delivered_at, failed_reason,
			status_history, created_at, updated_at
		FROM bex_shipments
		ORDER BY created_at DESC
		LIMIT $1
	`

	rows, err := r.db.Query(bexQuery, limit/2)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get BEX shipments")
	}
	defer rows.Close()

	for rows.Next() {
		var s BEXShipment
		err := rows.Scan(
			&s.ID, &s.MarketplaceOrderID, &s.StorefrontOrderID, &s.TrackingNumber,
			&s.RecipientName, &s.RecipientCity, &s.RecipientPhone,
			&s.Status, &s.StatusText, &s.CODAmount, &s.WeightKg,
			&s.RegisteredAt, &s.DeliveredAt, &s.FailedReason,
			&s.StatusHistory, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan BEX shipment")
		}
		shipments = append(shipments, map[string]interface{}{
			"provider": "BEX",
			"shipment": s,
		})
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error iterating BEX shipments")
	}

	// Get recent Post Express shipments
	postQuery := `
		SELECT 
			id, marketplace_order_id, storefront_order_id, tracking_number,
			recipient_name, recipient_city, recipient_phone,
			status, delivery_status, cod_amount, weight_kg, total_price,
			registered_at, delivered_at, failed_reason,
			status_history, created_at, updated_at
		FROM post_express_shipments
		ORDER BY created_at DESC
		LIMIT $1
	`

	rows2, err := r.db.Query(postQuery, limit/2)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Post Express shipments")
	}
	defer rows2.Close()

	for rows2.Next() {
		var s PostExpressShipment
		err := rows2.Scan(
			&s.ID, &s.MarketplaceOrderID, &s.StorefrontOrderID, &s.TrackingNumber,
			&s.RecipientName, &s.RecipientCity, &s.RecipientPhone,
			&s.Status, &s.DeliveryStatus, &s.CODAmount, &s.WeightKg, &s.TotalPrice,
			&s.RegisteredAt, &s.DeliveredAt, &s.FailedReason,
			&s.StatusHistory, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan Post Express shipment")
		}
		shipments = append(shipments, map[string]interface{}{
			"provider": "PostExpress",
			"shipment": s,
		})
	}
	if err = rows2.Err(); err != nil {
		return nil, errors.Wrap(err, "error iterating Post Express shipments")
	}

	return shipments, nil
}

func (r *ShipmentsRepository) GetShipmentsList(page, limit int, filters map[string]interface{}) ([]map[string]interface{}, int, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	var shipments []map[string]interface{}
	var total int

	// Получаем BEX отправления
	bexQuery := `
		SELECT 
			'BEX' as provider,
			id, tracking_number, status,
			recipient_name, recipient_city, recipient_phone,
			sender_name, sender_city,
			cod_amount, weight_kg,
			created_at, delivered_at
		FROM bex_shipments
		WHERE 1=1
	`

	// Применяем фильтры
	var bexArgs []interface{}
	argCount := 0

	if status, ok := filters["status"].(string); ok && status != "" {
		argCount++
		bexQuery += fmt.Sprintf(" AND status = $%d", argCount)
		bexArgs = append(bexArgs, getStatusCode(status))
	}

	if trackingNumber, ok := filters["tracking_number"].(string); ok && trackingNumber != "" {
		argCount++
		bexQuery += fmt.Sprintf(" AND tracking_number ILIKE $%d", argCount)
		bexArgs = append(bexArgs, "%"+trackingNumber+"%")
	}

	if city, ok := filters["city"].(string); ok && city != "" {
		argCount++
		bexQuery += fmt.Sprintf(" AND (recipient_city ILIKE $%d OR sender_city ILIKE $%d)", argCount, argCount)
		bexArgs = append(bexArgs, "%"+city+"%")
	}

	// Подсчет общего количества BEX
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s) t", bexQuery)
	var bexCount int
	err := r.db.QueryRow(countQuery, bexArgs...).Scan(&bexCount)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to count BEX shipments")
	}

	// Добавляем пагинацию
	argCount++
	bexQuery += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d", argCount)
	bexArgs = append(bexArgs, limit)
	argCount++
	bexQuery += fmt.Sprintf(" OFFSET $%d", argCount)
	bexArgs = append(bexArgs, offset)

	// Выполняем запрос
	rows, err := r.db.Query(bexQuery, bexArgs...)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to get BEX shipments")
	}
	defer rows.Close()

	for rows.Next() {
		var s struct {
			Provider       string
			ID             int
			TrackingNumber string
			Status         int
			RecipientName  string
			RecipientCity  string
			RecipientPhone string
			SenderName     string
			SenderCity     string
			CODAmount      float64
			WeightKg       float64
			CreatedAt      time.Time
			DeliveredAt    *time.Time
		}

		err := rows.Scan(
			&s.Provider, &s.ID, &s.TrackingNumber, &s.Status,
			&s.RecipientName, &s.RecipientCity, &s.RecipientPhone,
			&s.SenderName, &s.SenderCity,
			&s.CODAmount, &s.WeightKg,
			&s.CreatedAt, &s.DeliveredAt,
		)
		if err != nil {
			continue
		}

		shipments = append(shipments, map[string]interface{}{
			"id":              s.ID,
			"provider":        s.Provider,
			"tracking_number": s.TrackingNumber,
			"status":          getStatusString(s.Status),
			"recipient_name":  s.RecipientName,
			"recipient_city":  s.RecipientCity,
			"recipient_phone": s.RecipientPhone,
			"sender_name":     s.SenderName,
			"sender_city":     s.SenderCity,
			"cod_amount":      s.CODAmount,
			"weight_kg":       s.WeightKg,
			"created_at":      s.CreatedAt,
			"delivered_at":    s.DeliveredAt,
		})
	}
	if err = rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "error iterating BEX shipments")
	}

	// Аналогично для Post Express
	postQuery := `
		SELECT 
			'PostExpress' as provider,
			id, tracking_number, status,
			recipient_name, recipient_city, recipient_phone,
			sender_name, sender_city,
			cod_amount, weight_kg,
			created_at, delivered_at
		FROM post_express_shipments
		WHERE 1=1
	`

	// Фильтры для Post Express
	var postArgs []interface{}
	argCount = 0

	if status, ok := filters["status"].(string); ok && status != "" {
		argCount++
		postQuery += fmt.Sprintf(" AND status = $%d", argCount)
		postArgs = append(postArgs, status)
	}

	if trackingNumber, ok := filters["tracking_number"].(string); ok && trackingNumber != "" {
		argCount++
		postQuery += fmt.Sprintf(" AND tracking_number ILIKE $%d", argCount)
		postArgs = append(postArgs, "%"+trackingNumber+"%")
	}

	if city, ok := filters["city"].(string); ok && city != "" {
		argCount++
		postQuery += fmt.Sprintf(" AND (recipient_city ILIKE $%d OR sender_city ILIKE $%d)", argCount, argCount)
		postArgs = append(postArgs, "%"+city+"%")
	}

	// Подсчет Post Express
	countQuery = fmt.Sprintf("SELECT COUNT(*) FROM (%s) t", postQuery)
	var postCount int
	err = r.db.QueryRow(countQuery, postArgs...).Scan(&postCount)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to count Post Express shipments")
	}

	// Добавляем пагинацию
	remainingLimit := limit - len(shipments)
	if remainingLimit > 0 {
		argCount++
		postQuery += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d", argCount)
		postArgs = append(postArgs, remainingLimit)
		argCount++
		postQuery += fmt.Sprintf(" OFFSET $%d", argCount)
		postArgs = append(postArgs, 0) // offset для Post Express всегда 0 при простой пагинации

		rows2, err := r.db.Query(postQuery, postArgs...)
		if err != nil {
			return nil, 0, errors.Wrap(err, "failed to get Post Express shipments")
		}
		defer rows2.Close()

		for rows2.Next() {
			var s struct {
				Provider       string
				ID             int
				TrackingNumber string
				Status         string
				RecipientName  string
				RecipientCity  string
				RecipientPhone string
				SenderName     string
				SenderCity     string
				CODAmount      float64
				WeightKg       float64
				CreatedAt      time.Time
				DeliveredAt    *time.Time
			}

			err := rows2.Scan(
				&s.Provider, &s.ID, &s.TrackingNumber, &s.Status,
				&s.RecipientName, &s.RecipientCity, &s.RecipientPhone,
				&s.SenderName, &s.SenderCity,
				&s.CODAmount, &s.WeightKg,
				&s.CreatedAt, &s.DeliveredAt,
			)
			if err != nil {
				continue
			}

			shipments = append(shipments, map[string]interface{}{
				"id":              s.ID,
				"provider":        s.Provider,
				"tracking_number": s.TrackingNumber,
				"status":          s.Status,
				"recipient_name":  s.RecipientName,
				"recipient_city":  s.RecipientCity,
				"recipient_phone": s.RecipientPhone,
				"sender_name":     s.SenderName,
				"sender_city":     s.SenderCity,
				"cod_amount":      s.CODAmount,
				"weight_kg":       s.WeightKg,
				"created_at":      s.CreatedAt,
				"delivered_at":    s.DeliveredAt,
			})
		}
		if err = rows2.Err(); err != nil {
			return nil, 0, errors.Wrap(err, "error iterating Post Express shipments")
		}
	}

	total = bexCount + postCount
	return shipments, total, nil
}

func getStatusCode(status string) int {
	switch status {
	case "pending":
		return 1
	case "in_transit":
		return 2
	case "delivered":
		return 3
	case "failed":
		return 4
	default:
		return 0
	}
}

func getStatusString(code int) string {
	switch code {
	case 1:
		return "pending"
	case 2:
		return "in_transit"
	case 3:
		return "delivered"
	case 4:
		return "failed"
	default:
		return "unknown"
	}
}

func (r *ShipmentsRepository) GetShipmentDetails(provider string, id int) (map[string]interface{}, error) {
	var result map[string]interface{}

	switch provider {
	case "BEX":
		query := `
			SELECT 
				id, marketplace_order_id, storefront_order_id, tracking_number,
				recipient_name, recipient_address, recipient_city, recipient_postal_code,
				recipient_phone, recipient_email,
				sender_name, sender_address, sender_city, sender_postal_code,
				sender_phone, sender_email,
				status, status_text, cod_amount, weight_kg,
				shipment_contents, comment_public,
				registered_at, delivered_at, failed_reason,
				status_history, created_at, updated_at
			FROM bex_shipments
			WHERE id = $1
		`

		var s struct {
			ID                 int
			MarketplaceOrderID *int
			StorefrontOrderID  *int64
			TrackingNumber     string
			RecipientName      string
			RecipientAddress   string
			RecipientCity      string
			RecipientZip       string
			RecipientPhone     string
			RecipientEmail     *string
			SenderName         string
			SenderAddress      string
			SenderCity         string
			SenderZip          string
			SenderPhone        string
			SenderEmail        *string
			Status             int
			StatusText         string
			CODAmount          float64
			WeightKg           float64
			PackageContents    *int
			CommentPublic      *string
			RegisteredAt       *time.Time
			DeliveredAt        *time.Time
			FailedReason       *string
			StatusHistory      json.RawMessage
			CreatedAt          time.Time
			UpdatedAt          time.Time
		}

		err := r.db.QueryRow(query, id).Scan(
			&s.ID, &s.MarketplaceOrderID, &s.StorefrontOrderID, &s.TrackingNumber,
			&s.RecipientName, &s.RecipientAddress, &s.RecipientCity, &s.RecipientZip,
			&s.RecipientPhone, &s.RecipientEmail,
			&s.SenderName, &s.SenderAddress, &s.SenderCity, &s.SenderZip,
			&s.SenderPhone, &s.SenderEmail,
			&s.Status, &s.StatusText, &s.CODAmount, &s.WeightKg,
			&s.PackageContents, &s.CommentPublic,
			&s.RegisteredAt, &s.DeliveredAt, &s.FailedReason,
			&s.StatusHistory, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("shipment not found")
			}
			return nil, errors.Wrap(err, "failed to get BEX shipment details")
		}

		// Parse status history
		var statusHistory []map[string]interface{}
		if s.StatusHistory != nil {
			_ = json.Unmarshal(s.StatusHistory, &statusHistory)
		}

		result = map[string]interface{}{
			"provider":             "BEX",
			"id":                   s.ID,
			"marketplace_order_id": s.MarketplaceOrderID,
			"storefront_order_id":  s.StorefrontOrderID,
			"tracking_number":      s.TrackingNumber,
			"recipient": map[string]interface{}{
				"name":    s.RecipientName,
				"address": s.RecipientAddress,
				"city":    s.RecipientCity,
				"zip":     s.RecipientZip,
				"phone":   s.RecipientPhone,
				"email":   s.RecipientEmail,
			},
			"sender": map[string]interface{}{
				"name":    s.SenderName,
				"address": s.SenderAddress,
				"city":    s.SenderCity,
				"zip":     s.SenderZip,
				"phone":   s.SenderPhone,
				"email":   s.SenderEmail,
			},
			"status":           getStatusString(s.Status),
			"status_text":      s.StatusText,
			"cod_amount":       s.CODAmount,
			"weight_kg":        s.WeightKg,
			"package_contents": s.PackageContents,
			"reference_number": s.CommentPublic,
			"registered_at":    s.RegisteredAt,
			"delivered_at":     s.DeliveredAt,
			"failed_reason":    s.FailedReason,
			"status_history":   statusHistory,
			"created_at":       s.CreatedAt,
			"updated_at":       s.UpdatedAt,
		}

	case "PostExpress":
		query := `
			SELECT 
				id, marketplace_order_id, storefront_order_id, tracking_number,
				recipient_name, recipient_address, recipient_city, recipient_postal_code,
				recipient_phone, recipient_email,
				sender_name, sender_address, sender_city, sender_postal_code,
				sender_phone, sender_email,
				status, delivery_status, cod_amount, weight_kg, total_price,
				notes, delivery_instructions,
				registered_at, delivered_at, failed_reason,
				status_history, created_at, updated_at
			FROM post_express_shipments
			WHERE id = $1
		`

		var s struct {
			ID                   int
			MarketplaceOrderID   *int
			StorefrontOrderID    *int64
			TrackingNumber       string
			RecipientName        string
			RecipientAddress     string
			RecipientCity        string
			RecipientZip         string
			RecipientPhone       string
			RecipientEmail       *string
			SenderName           string
			SenderAddress        string
			SenderCity           string
			SenderZip            string
			SenderPhone          string
			SenderEmail          *string
			Status               string
			DeliveryStatus       *string
			CODAmount            float64
			WeightKg             float64
			TotalPrice           float64
			Notes                *string
			DeliveryInstructions *string
			RegisteredAt         *time.Time
			DeliveredAt          *time.Time
			FailedReason         *string
			StatusHistory        json.RawMessage
			CreatedAt            time.Time
			UpdatedAt            time.Time
		}

		err := r.db.QueryRow(query, id).Scan(
			&s.ID, &s.MarketplaceOrderID, &s.StorefrontOrderID, &s.TrackingNumber,
			&s.RecipientName, &s.RecipientAddress, &s.RecipientCity, &s.RecipientZip,
			&s.RecipientPhone, &s.RecipientEmail,
			&s.SenderName, &s.SenderAddress, &s.SenderCity, &s.SenderZip,
			&s.SenderPhone, &s.SenderEmail,
			&s.Status, &s.DeliveryStatus, &s.CODAmount, &s.WeightKg, &s.TotalPrice,
			&s.Notes, &s.DeliveryInstructions,
			&s.RegisteredAt, &s.DeliveredAt, &s.FailedReason,
			&s.StatusHistory, &s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("shipment not found")
			}
			return nil, errors.Wrap(err, "failed to get Post Express shipment details")
		}

		// Parse status history
		var statusHistory []map[string]interface{}
		if s.StatusHistory != nil {
			_ = json.Unmarshal(s.StatusHistory, &statusHistory)
		}

		result = map[string]interface{}{
			"provider":             "PostExpress",
			"id":                   s.ID,
			"marketplace_order_id": s.MarketplaceOrderID,
			"storefront_order_id":  s.StorefrontOrderID,
			"tracking_number":      s.TrackingNumber,
			"recipient": map[string]interface{}{
				"name":    s.RecipientName,
				"address": s.RecipientAddress,
				"city":    s.RecipientCity,
				"zip":     s.RecipientZip,
				"phone":   s.RecipientPhone,
				"email":   s.RecipientEmail,
			},
			"sender": map[string]interface{}{
				"name":    s.SenderName,
				"address": s.SenderAddress,
				"city":    s.SenderCity,
				"zip":     s.SenderZip,
				"phone":   s.SenderPhone,
				"email":   s.SenderEmail,
			},
			"status":                s.Status,
			"delivery_status":       s.DeliveryStatus,
			"cod_amount":            s.CODAmount,
			"weight_kg":             s.WeightKg,
			"total_price":           s.TotalPrice,
			"notes":                 s.Notes,
			"delivery_instructions": s.DeliveryInstructions,
			"registered_at":         s.RegisteredAt,
			"delivered_at":          s.DeliveredAt,
			"failed_reason":         s.FailedReason,
			"status_history":        statusHistory,
			"created_at":            s.CreatedAt,
			"updated_at":            s.UpdatedAt,
		}
	default:
		return nil, errors.New("invalid provider")
	}

	return result, nil
}

func (r *ShipmentsRepository) GetProblemShipments() ([]interface{}, error) {
	var problems []interface{}

	// Get BEX problems (failed or stuck in transit > 7 days)
	bexQuery := `
		SELECT 
			id, tracking_number, recipient_name, recipient_city,
			status, status_text, failed_reason, registered_at,
			CASE 
				WHEN status = 4 THEN 'failed'
				WHEN status = 2 AND registered_at < NOW() - INTERVAL '7 days' THEN 'delayed'
				ELSE 'unknown'
			END as problem_type
		FROM bex_shipments
		WHERE status = 4 
		   OR (status = 2 AND registered_at < NOW() - INTERVAL '7 days')
		ORDER BY created_at DESC
		LIMIT 50
	`

	rows, err := r.db.Query(bexQuery)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get BEX problems")
	}
	defer rows.Close()

	for rows.Next() {
		var p struct {
			ID             int
			TrackingNumber string
			RecipientName  string
			RecipientCity  string
			Status         int
			StatusText     string
			FailedReason   *string
			RegisteredAt   *time.Time
			ProblemType    string
		}

		err := rows.Scan(
			&p.ID, &p.TrackingNumber, &p.RecipientName, &p.RecipientCity,
			&p.Status, &p.StatusText, &p.FailedReason, &p.RegisteredAt,
			&p.ProblemType,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan BEX problem")
		}

		problems = append(problems, map[string]interface{}{
			"provider":        "BEX",
			"id":              p.ID,
			"tracking_number": p.TrackingNumber,
			"recipient_name":  p.RecipientName,
			"recipient_city":  p.RecipientCity,
			"problem_type":    p.ProblemType,
			"failed_reason":   p.FailedReason,
			"registered_at":   p.RegisteredAt,
		})
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error iterating BEX problems")
	}

	// Get Post Express problems
	postQuery := `
		SELECT 
			id, tracking_number, recipient_name, recipient_city,
			status, failed_reason, registered_at,
			CASE 
				WHEN status = 'failed' THEN 'failed'
				WHEN status = 'in_transit' AND registered_at < NOW() - INTERVAL '7 days' THEN 'delayed'
				ELSE 'unknown'
			END as problem_type
		FROM post_express_shipments
		WHERE status = 'failed' 
		   OR (status = 'in_transit' AND registered_at < NOW() - INTERVAL '7 days')
		ORDER BY created_at DESC
		LIMIT 50
	`

	rows2, err := r.db.Query(postQuery)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Post Express problems")
	}
	defer rows2.Close()

	for rows2.Next() {
		var p struct {
			ID             int
			TrackingNumber string
			RecipientName  string
			RecipientCity  string
			Status         string
			FailedReason   *string
			RegisteredAt   *time.Time
			ProblemType    string
		}

		err := rows2.Scan(
			&p.ID, &p.TrackingNumber, &p.RecipientName, &p.RecipientCity,
			&p.Status, &p.FailedReason, &p.RegisteredAt,
			&p.ProblemType,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan Post Express problem")
		}

		problems = append(problems, map[string]interface{}{
			"provider":        "PostExpress",
			"id":              p.ID,
			"tracking_number": p.TrackingNumber,
			"recipient_name":  p.RecipientName,
			"recipient_city":  p.RecipientCity,
			"problem_type":    p.ProblemType,
			"failed_reason":   p.FailedReason,
			"registered_at":   p.RegisteredAt,
		})
	}
	if err = rows2.Err(); err != nil {
		return nil, errors.Wrap(err, "error iterating Post Express problems")
	}

	return problems, nil
}
