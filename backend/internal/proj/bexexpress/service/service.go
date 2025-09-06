package service

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"backend/internal/config"
	"backend/internal/proj/bexexpress/models"
	"backend/pkg/logger"
)

// NotificationService интерфейс для отправки уведомлений
type NotificationService interface {
	SendDeliveryStatusNotification(ctx context.Context, userID int, orderID int, deliveryProvider string, status string, statusText string, trackingNumber string) error
}

// Service представляет сервис для работы с BEX Express
type Service struct {
	db                  *sql.DB
	cfg                 *config.Config
	client              *BEXClient
	settings            *models.BEXSettings
	notificationService NotificationService
	logger              *logger.Logger
}

// NewService создает новый сервис BEX Express
func NewService(db *sql.DB, cfg *config.Config) (*Service, error) {
	s := &Service{
		db:     db,
		cfg:    cfg,
		logger: logger.GetLogger(),
	}
	return s.init()
}

// NewServiceWithNotifications создает новый сервис BEX Express с поддержкой уведомлений
func NewServiceWithNotifications(db *sql.DB, cfg *config.Config, notificationService NotificationService) (*Service, error) {
	s := &Service{
		db:                  db,
		cfg:                 cfg,
		notificationService: notificationService,
		logger:              logger.GetLogger(),
	}
	return s.init()
}

// init инициализирует сервис
func (s *Service) init() (*Service, error) {
	// Load settings from database or use defaults from config
	settings, err := s.loadSettings()
	if err != nil {
		// Use defaults from config
		settings = &models.BEXSettings{
			AuthToken:   s.cfg.BEXAuthToken,
			ClientID:    s.cfg.BEXClientID,
			APIEndpoint: s.cfg.BEXAPIURL,
			Enabled:     true,
		}
	}
	s.settings = settings

	// Create API client
	s.client = NewBEXClient(settings.AuthToken, settings.ClientID, settings.APIEndpoint)

	return s, nil
}

// loadSettings загружает настройки из базы данных
func (s *Service) loadSettings() (*models.BEXSettings, error) {
	var settings models.BEXSettings
	query := `
		SELECT id, auth_token, client_id, api_endpoint, sender_client_id,
		       sender_name, sender_address, sender_city, sender_postal_code,
		       sender_phone, sender_email, enabled, test_mode,
		       auto_print_labels, auto_track_shipments, use_address_lookup,
		       notify_on_pickup, notify_on_delivery, notify_on_failed_delivery,
		       total_shipments, successful_deliveries, failed_deliveries,
		       created_at, updated_at
		FROM bex_settings
		WHERE enabled = true
		LIMIT 1
	`

	err := s.db.QueryRow(query).Scan(
		&settings.ID, &settings.AuthToken, &settings.ClientID, &settings.APIEndpoint,
		&settings.SenderClientID, &settings.SenderName, &settings.SenderAddress,
		&settings.SenderCity, &settings.SenderPostalCode, &settings.SenderPhone,
		&settings.SenderEmail, &settings.Enabled, &settings.TestMode,
		&settings.AutoPrintLabels, &settings.AutoTrackShipments, &settings.UseAddressLookup,
		&settings.NotifyOnPickup, &settings.NotifyOnDelivery, &settings.NotifyOnFailedDelivery,
		&settings.TotalShipments, &settings.SuccessfulDeliveries, &settings.FailedDeliveries,
		&settings.CreatedAt, &settings.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no BEX settings found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to load BEX settings: %w", err)
	}

	return &settings, nil
}

// CreateShipment создает новое отправление
func (s *Service) CreateShipment(ctx context.Context, req *models.CreateShipmentRequest) (*models.BEXShipment, error) {
	if !s.settings.Enabled {
		return nil, fmt.Errorf("BEX Express integration is disabled")
	}

	// Build shipment data for API
	shipmentData := BuildShipmentData(req, s.settings.SenderClientID)

	// Call BEX API
	response, err := s.client.CreateShipment(shipmentData)
	if err != nil {
		return nil, fmt.Errorf("failed to create shipment in BEX: %w", err)
	}

	if len(response.ShipmentsResultList) == 0 {
		return nil, fmt.Errorf("no shipment result returned from BEX")
	}

	result := response.ShipmentsResultList[0]
	if !result.State {
		return nil, fmt.Errorf("BEX rejected shipment: %s", result.Error)
	}

	// Save shipment to database
	shipment := &models.BEXShipment{
		MarketplaceOrderID: req.OrderID,
		StorefrontOrderID:  req.StorefrontOrderID,
		BexShipmentID:      &result.ShipmentID,
		TrackingNumber:     ptrStr(fmt.Sprintf("%d", result.ShipmentID)),

		// Sender info from settings
		SenderName:       s.settings.SenderName,
		SenderAddress:    s.settings.SenderAddress,
		SenderCity:       s.settings.SenderCity,
		SenderPostalCode: s.settings.SenderPostalCode,
		SenderPhone:      s.settings.SenderPhone,
		SenderEmail:      s.settings.SenderEmail,

		// Recipient info
		RecipientName:       req.RecipientName,
		RecipientAddress:    req.RecipientAddress,
		RecipientCity:       req.RecipientCity,
		RecipientPostalCode: req.RecipientPostalCode,
		RecipientPhone:      req.RecipientPhone,
		RecipientEmail:      req.RecipientEmail,

		// Shipment params
		ShipmentType:     shipmentData.ShipmentType,
		ShipmentCategory: shipmentData.ShipmentCategory,
		ShipmentContents: shipmentData.ShipmentContents,
		WeightKg:         req.WeightKg,
		TotalPackages:    req.TotalPackages,

		// Services
		PayType:                  shipmentData.PayType,
		CODAmount:                req.CODAmount,
		InsuranceAmount:          req.InsuranceAmount,
		PersonalDelivery:         req.PersonalDelivery,
		ReturnSignedInvoices:     false,
		ReturnSignedConfirmation: false,
		ReturnPackage:            false,

		// Comments
		CommentPublic:        req.Notes,
		DeliveryInstructions: req.DeliveryInstructions,

		// Status
		Status:       models.ShipmentStatusNotSentYet,
		RegisteredAt: ptrTime(time.Now()),

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insert into database
	query := `
		INSERT INTO bex_shipments (
			marketplace_order_id, storefront_order_id, bex_shipment_id, tracking_number,
			sender_name, sender_address, sender_city, sender_postal_code, sender_phone, sender_email,
			recipient_name, recipient_address, recipient_city, recipient_postal_code, recipient_phone, recipient_email,
			shipment_type, shipment_category, shipment_contents, weight_kg, total_packages,
			pay_type, cod_amount, insurance_amount, personal_delivery,
			return_signed_invoices, return_signed_confirmation, return_package,
			comment_public, comment_private, delivery_instructions,
			status, registered_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
			$21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
			$31, $32, $33, $34, $35
		) RETURNING id`

	err = s.db.QueryRowContext(ctx, query,
		shipment.MarketplaceOrderID, shipment.StorefrontOrderID, shipment.BexShipmentID, shipment.TrackingNumber,
		shipment.SenderName, shipment.SenderAddress, shipment.SenderCity, shipment.SenderPostalCode,
		shipment.SenderPhone, shipment.SenderEmail,
		shipment.RecipientName, shipment.RecipientAddress, shipment.RecipientCity, shipment.RecipientPostalCode,
		shipment.RecipientPhone, shipment.RecipientEmail,
		shipment.ShipmentType, shipment.ShipmentCategory, shipment.ShipmentContents,
		shipment.WeightKg, shipment.TotalPackages,
		shipment.PayType, shipment.CODAmount, shipment.InsuranceAmount, shipment.PersonalDelivery,
		shipment.ReturnSignedInvoices, shipment.ReturnSignedConfirmation, shipment.ReturnPackage,
		shipment.CommentPublic, shipment.CommentPrivate, shipment.DeliveryInstructions,
		shipment.Status, shipment.RegisteredAt, shipment.CreatedAt, shipment.UpdatedAt,
	).Scan(&shipment.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to save shipment to database: %w", err)
	}

	// Update statistics
	s.updateStatistics(1, 0, 0)

	// Get label if auto-print is enabled
	if s.settings.AutoPrintLabels && shipment.BexShipmentID != nil {
		labelData, err := s.client.GetShipmentLabel(*shipment.BexShipmentID, 4, 1)
		if err == nil {
			labelBase64 := base64Encode(labelData)
			shipment.LabelBase64 = &labelBase64

			// Update database with label
			updateQuery := `UPDATE bex_shipments SET label_base64 = $1 WHERE id = $2`
			if _, err := s.db.ExecContext(ctx, updateQuery, labelBase64, shipment.ID); err != nil {
				s.logger.Error("Failed to update label in database: %v", err)
			}
		}
	}

	return shipment, nil
}

// GetShipmentStatus получает статус отправления
func (s *Service) GetShipmentStatus(ctx context.Context, shipmentID int) (*models.BEXShipment, error) {
	// Get shipment from database
	shipment, err := s.getShipmentByID(ctx, shipmentID)
	if err != nil {
		return nil, err
	}

	if shipment.BexShipmentID == nil {
		return shipment, nil // No BEX ID yet
	}

	// Get status from BEX API
	statusResp, err := s.client.GetShipmentStatus(*shipment.BexShipmentID, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get status from BEX: %w", err)
	}

	// Update status in database
	newStatus := models.ShipmentStatus(statusResp.Status)
	if newStatus != shipment.Status {
		shipment.Status = newStatus
		shipment.StatusText = &statusResp.StatusText
		shipment.UpdatedAt = time.Now()

		// Update timestamps based on status
		now := time.Now()
		switch newStatus {
		case models.ShipmentStatusPickedUp:
			shipment.PickedUpAt = &now
		case models.ShipmentStatusDelivered:
			shipment.DeliveredAt = &now
			s.updateStatistics(0, 1, 0)
		case models.ShipmentStatusReturnedToSender:
			shipment.ReturnedAt = &now
			s.updateStatistics(0, 0, 1)
		case models.ShipmentStatusNotFound,
			models.ShipmentStatusDeleted,
			models.ShipmentStatusNotSentYet,
			models.ShipmentStatusAddressVerify,
			models.ShipmentStatusReturnToSender,
			models.ShipmentStatusCODPaid:
			// These statuses don't require specific timestamp updates
		}

		// Add to status history
		var history []map[string]interface{}
		if shipment.StatusHistory != nil {
			if err := json.Unmarshal(shipment.StatusHistory, &history); err != nil {
				s.logger.Error("Failed to unmarshal status history: %v", err)
			}
		}
		history = append(history, map[string]interface{}{
			"status":      newStatus,
			"status_text": statusResp.StatusText,
			"timestamp":   now,
		})
		historyJSON, _ := json.Marshal(history)
		shipment.StatusHistory = historyJSON

		// Update database
		updateQuery := `
			UPDATE bex_shipments 
			SET status = $1, status_text = $2, status_history = $3,
			    picked_up_at = $4, delivered_at = $5, returned_at = $6,
			    updated_at = $7
			WHERE id = $8`

		_, err = s.db.ExecContext(ctx, updateQuery,
			shipment.Status, shipment.StatusText, shipment.StatusHistory,
			shipment.PickedUpAt, shipment.DeliveredAt, shipment.ReturnedAt,
			shipment.UpdatedAt, shipment.ID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to update shipment status: %w", err)
		}

		// Отправляем уведомление об изменении статуса
		if s.notificationService != nil {
			// Определяем пользователя для уведомления
			var userID int
			var orderID int

			if shipment.MarketplaceOrderID != nil {
				orderID = *shipment.MarketplaceOrderID
				// Получаем user_id из marketplace_orders
				userQuery := "SELECT user_id FROM marketplace_orders WHERE id = $1"
				err = s.db.QueryRowContext(ctx, userQuery, orderID).Scan(&userID)
			} else if shipment.StorefrontOrderID != nil {
				orderID = int(*shipment.StorefrontOrderID)
				// Получаем user_id из storefront_orders
				userQuery := "SELECT user_id FROM storefront_orders WHERE id = $1"
				err = s.db.QueryRowContext(ctx, userQuery, orderID).Scan(&userID)
			}

			if err == nil && userID > 0 {
				// Отправляем уведомление
				statusText := "Неизвестно"
				if shipment.StatusText != nil {
					statusText = *shipment.StatusText
				}

				trackingNumber := ""
				if shipment.TrackingNumber != nil {
					trackingNumber = *shipment.TrackingNumber
				}

				notificationErr := s.notificationService.SendDeliveryStatusNotification(
					ctx, userID, orderID, "BEX Express",
					fmt.Sprintf("%d", newStatus), statusText, trackingNumber,
				)
				if notificationErr != nil {
					// Логируем ошибку, но не прерываем выполнение
					fmt.Printf("Failed to send delivery notification: %v", notificationErr)
				}
			}
		}
	}

	return shipment, nil
}

// GetShipmentLabel получает этикетку для печати
func (s *Service) GetShipmentLabel(ctx context.Context, shipmentID int, pageSize int) ([]byte, error) {
	// Get shipment from database
	shipment, err := s.getShipmentByID(ctx, shipmentID)
	if err != nil {
		return nil, err
	}

	if shipment.BexShipmentID == nil {
		return nil, fmt.Errorf("shipment has no BEX ID")
	}

	// Get label from BEX API
	labelData, err := s.client.GetShipmentLabel(*shipment.BexShipmentID, pageSize, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get label from BEX: %w", err)
	}

	// Save label to database
	labelBase64 := base64Encode(labelData)
	updateQuery := `UPDATE bex_shipments SET label_base64 = $1 WHERE id = $2`
	if _, err := s.db.ExecContext(ctx, updateQuery, labelBase64, shipment.ID); err != nil {
		s.logger.Error("Failed to update label in database: %v", err)
	}

	return labelData, nil
}

// CancelShipment отменяет отправление
func (s *Service) CancelShipment(ctx context.Context, shipmentID int) error {
	// Get shipment from database
	shipment, err := s.getShipmentByID(ctx, shipmentID)
	if err != nil {
		return err
	}

	if shipment.BexShipmentID == nil {
		return fmt.Errorf("shipment has no BEX ID")
	}

	// Delete from BEX
	err = s.client.DeleteShipment(*shipment.BexShipmentID)
	if err != nil {
		return fmt.Errorf("failed to delete shipment from BEX: %w", err)
	}

	// Update status in database
	updateQuery := `
		UPDATE bex_shipments 
		SET status = $1, updated_at = $2
		WHERE id = $3`

	_, err = s.db.ExecContext(ctx, updateQuery,
		models.ShipmentStatusDeleted, time.Now(), shipment.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update shipment status: %w", err)
	}

	return nil
}

// SearchAddress ищет адреса в справочниках BEX
func (s *Service) SearchAddress(ctx context.Context, req *models.SearchAddressRequest) ([]models.AddressSuggestion, error) {
	if !s.settings.UseAddressLookup {
		return []models.AddressSuggestion{}, nil
	}

	query := `
		SELECT DISTINCT
			p.bex_id as place_id,
			p.name as place_name,
			s.bex_id as street_id,
			s.name as street_name,
			p.postal_code,
			m.bex_id as municipality_id,
			m.name as municipality
		FROM bex_streets s
		JOIN bex_places p ON s.place_id = p.id
		JOIN bex_municipalities m ON p.municipality_id = m.id
		WHERE s.is_active = true AND p.is_active = true AND m.is_active = true
		  AND (s.name ILIKE $1 OR s.name_cyrillic ILIKE $1)
	`

	args := []interface{}{
		"%" + req.Query + "%",
	}
	argCount := 1

	if req.City != "" {
		argCount++
		query += fmt.Sprintf(" AND (p.name ILIKE $%d OR p.name_cyrillic ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+req.City+"%")
	}

	if req.MunicipalityID != nil {
		argCount++
		query += fmt.Sprintf(" AND m.bex_id = $%d", argCount)
		args = append(args, *req.MunicipalityID)
	}

	if req.PlaceID != nil {
		argCount++
		query += fmt.Sprintf(" AND p.id = $%d", argCount)
		args = append(args, *req.PlaceID)
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	argCount++
	query += fmt.Sprintf(" LIMIT $%d", argCount)
	args = append(args, limit)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search address: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var suggestions []models.AddressSuggestion
	for rows.Next() {
		var s models.AddressSuggestion
		err := rows.Scan(
			&s.PlaceID, &s.PlaceName,
			&s.StreetID, &s.StreetName,
			&s.PostalCode,
			&s.MunicipalityID, &s.Municipality,
		)
		if err != nil {
			continue
		}

		s.FullAddress = fmt.Sprintf("%s, %s %s, %s",
			s.StreetName, s.PostalCode, s.PlaceName, s.Municipality)

		suggestions = append(suggestions, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return suggestions, nil
}

// CalculateRate рассчитывает стоимость доставки
func (s *Service) CalculateRate(ctx context.Context, req *models.CalculateRateRequest) (*models.CalculateRateResponse, error) {
	// Basic rate calculation based on weight and category
	var basePrice float64

	switch req.ShipmentCategory {
	case 1: // Documents up to 0.5kg
		basePrice = 250
	case 2: // Package up to 1kg
		basePrice = 350
	case 3: // Package up to 2kg
		basePrice = 450
	case 31: // Package per kg
		basePrice = 450 + (req.WeightKg-2)*100
	case 32: // Pallet per kg
		basePrice = 1000 + req.WeightKg*50
	default:
		basePrice = 350
	}

	// COD fee
	codFee := 0.0
	if req.CODAmount != nil && *req.CODAmount > 0 {
		codFee = 150 + (*req.CODAmount * 0.01) // 1% of COD amount + fixed fee
	}

	// Insurance fee
	insuranceFee := 0.0
	if req.InsuranceAmount != nil && *req.InsuranceAmount > 0 {
		insuranceFee = *req.InsuranceAmount * 0.005 // 0.5% of insured amount
	}

	return &models.CalculateRateResponse{
		BasePrice:        basePrice,
		InsuranceFee:     insuranceFee,
		CODFee:           codFee,
		TotalPrice:       basePrice + insuranceFee + codFee,
		DeliveryDaysMin:  1,
		DeliveryDaysMax:  3,
		ServiceAvailable: true,
	}, nil
}

// GetShipmentByID получает отправление по ID (public method)
func (s *Service) GetShipmentByID(ctx context.Context, id int) (*models.BEXShipment, error) {
	return s.getShipmentByID(ctx, id)
}

// GetShipmentByTracking получает отправление по tracking number
func (s *Service) GetShipmentByTracking(ctx context.Context, tracking string) (*models.BEXShipment, error) {
	var shipment models.BEXShipment
	query := `
		SELECT id, marketplace_order_id, storefront_order_id, bex_shipment_id, tracking_number,
		       sender_name, sender_address, sender_city, sender_postal_code, sender_phone, sender_email,
		       recipient_name, recipient_address, recipient_city, recipient_postal_code, recipient_phone, recipient_email,
		       shipment_type, shipment_category, shipment_contents, weight_kg, total_packages,
		       pay_type, cod_amount, insurance_amount, personal_delivery,
		       return_signed_invoices, return_signed_confirmation, return_package,
		       comment_public, comment_private, delivery_instructions,
		       status, status_text, failed_reason, label_base64, label_url,
		       registered_at, picked_up_at, delivered_at, failed_at, returned_at,
		       status_history, created_at, updated_at
		FROM bex_shipments
		WHERE tracking_number = $1`

	err := s.db.QueryRowContext(ctx, query, tracking).Scan(
		&shipment.ID, &shipment.MarketplaceOrderID, &shipment.StorefrontOrderID,
		&shipment.BexShipmentID, &shipment.TrackingNumber,
		&shipment.SenderName, &shipment.SenderAddress, &shipment.SenderCity,
		&shipment.SenderPostalCode, &shipment.SenderPhone, &shipment.SenderEmail,
		&shipment.RecipientName, &shipment.RecipientAddress, &shipment.RecipientCity,
		&shipment.RecipientPostalCode, &shipment.RecipientPhone, &shipment.RecipientEmail,
		&shipment.ShipmentType, &shipment.ShipmentCategory, &shipment.ShipmentContents,
		&shipment.WeightKg, &shipment.TotalPackages,
		&shipment.PayType, &shipment.CODAmount, &shipment.InsuranceAmount,
		&shipment.PersonalDelivery, &shipment.ReturnSignedInvoices,
		&shipment.ReturnSignedConfirmation, &shipment.ReturnPackage,
		&shipment.CommentPublic, &shipment.CommentPrivate, &shipment.DeliveryInstructions,
		&shipment.Status, &shipment.StatusText, &shipment.FailedReason,
		&shipment.LabelBase64, &shipment.LabelURL,
		&shipment.RegisteredAt, &shipment.PickedUpAt, &shipment.DeliveredAt,
		&shipment.FailedAt, &shipment.ReturnedAt,
		&shipment.StatusHistory, &shipment.CreatedAt, &shipment.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("shipment not found: %w", err)
	}

	return &shipment, nil
}

// GetParcelShops получает список пунктов выдачи
func (s *Service) GetParcelShops(ctx context.Context, city string) ([]models.BEXParcelShop, error) {
	query := `
		SELECT id, bex_id, code, name, address, city, postal_code,
		       phone, working_hours, latitude, longitude, is_active
		FROM bex_parcel_shops
		WHERE is_active = true`

	args := []interface{}{}
	if city != "" {
		query += " AND city ILIKE $1"
		args = append(args, "%"+city+"%")
	}

	query += " ORDER BY name"

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var shops []models.BEXParcelShop
	for rows.Next() {
		var shop models.BEXParcelShop
		err := rows.Scan(
			&shop.ID, &shop.BexID, &shop.Code, &shop.Name,
			&shop.Address, &shop.City, &shop.PostalCode,
			&shop.Phone, &shop.WorkingHours,
			&shop.Latitude, &shop.Longitude, &shop.IsActive,
		)
		if err != nil {
			continue
		}
		shops = append(shops, shop)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return shops, nil
}

// Helper functions

func (s *Service) getShipmentByID(ctx context.Context, id int) (*models.BEXShipment, error) {
	var shipment models.BEXShipment
	query := `
		SELECT id, marketplace_order_id, storefront_order_id, bex_shipment_id, tracking_number,
		       sender_name, sender_address, sender_city, sender_postal_code, sender_phone, sender_email,
		       recipient_name, recipient_address, recipient_city, recipient_postal_code, recipient_phone, recipient_email,
		       shipment_type, shipment_category, shipment_contents, weight_kg, total_packages,
		       pay_type, cod_amount, insurance_amount, personal_delivery,
		       return_signed_invoices, return_signed_confirmation, return_package,
		       comment_public, comment_private, delivery_instructions,
		       status, status_text, failed_reason, label_base64, label_url,
		       registered_at, picked_up_at, delivered_at, failed_at, returned_at,
		       status_history, created_at, updated_at
		FROM bex_shipments
		WHERE id = $1`

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&shipment.ID, &shipment.MarketplaceOrderID, &shipment.StorefrontOrderID,
		&shipment.BexShipmentID, &shipment.TrackingNumber,
		&shipment.SenderName, &shipment.SenderAddress, &shipment.SenderCity,
		&shipment.SenderPostalCode, &shipment.SenderPhone, &shipment.SenderEmail,
		&shipment.RecipientName, &shipment.RecipientAddress, &shipment.RecipientCity,
		&shipment.RecipientPostalCode, &shipment.RecipientPhone, &shipment.RecipientEmail,
		&shipment.ShipmentType, &shipment.ShipmentCategory, &shipment.ShipmentContents,
		&shipment.WeightKg, &shipment.TotalPackages,
		&shipment.PayType, &shipment.CODAmount, &shipment.InsuranceAmount,
		&shipment.PersonalDelivery, &shipment.ReturnSignedInvoices,
		&shipment.ReturnSignedConfirmation, &shipment.ReturnPackage,
		&shipment.CommentPublic, &shipment.CommentPrivate, &shipment.DeliveryInstructions,
		&shipment.Status, &shipment.StatusText, &shipment.FailedReason,
		&shipment.LabelBase64, &shipment.LabelURL,
		&shipment.RegisteredAt, &shipment.PickedUpAt, &shipment.DeliveredAt,
		&shipment.FailedAt, &shipment.ReturnedAt,
		&shipment.StatusHistory, &shipment.CreatedAt, &shipment.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("shipment not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get shipment: %w", err)
	}

	return &shipment, nil
}

func (s *Service) updateStatistics(total, successful, failed int) {
	query := `
		UPDATE bex_settings
		SET total_shipments = total_shipments + $1,
		    successful_deliveries = successful_deliveries + $2,
		    failed_deliveries = failed_deliveries + $3,
		    updated_at = $4
		WHERE id = $5`

	if _, err := s.db.Exec(query, total, successful, failed, time.Now(), s.settings.ID); err != nil {
		s.logger.Error("Failed to update sync stats: %v", err)
	}
}

func ptrStr(s string) *string {
	return &s
}

func ptrTime(t time.Time) *time.Time {
	return &t
}

func base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
