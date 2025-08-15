package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"backend/internal/proj/postexpress/models"
	"backend/internal/proj/postexpress/storage"
)

// ErrNotFound is returned when an entity is not found
var ErrNotFound = errors.New("not found")

// Repository представляет PostgreSQL реализацию репозитория Post Express
type Repository struct {
	db *sqlx.DB
}

// NewRepository создает новый экземпляр репозитория
func NewRepository(db *sqlx.DB) storage.Repository {
	return &Repository{db: db}
}

// =============================================================================
// НАСТРОЙКИ
// =============================================================================

func (r *Repository) GetSettings(ctx context.Context) (*models.PostExpressSettings, error) {
	var settings models.PostExpressSettings
	query := `
		SELECT id, api_username, api_password, api_endpoint,
			   sender_name, sender_address, sender_city, sender_postal_code, sender_phone, sender_email,
			   enabled, test_mode, auto_print_labels, auto_track_shipments,
			   notify_on_pickup, notify_on_delivery, notify_on_failed_delivery,
			   total_shipments, successful_deliveries, failed_deliveries,
			   created_at, updated_at
		FROM post_express_settings 
		ORDER BY id DESC 
		LIMIT 1`

	err := r.db.GetContext(ctx, &settings, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound // Настройки не найдены
		}
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}

	return &settings, nil
}

func (r *Repository) UpdateSettings(ctx context.Context, settings *models.PostExpressSettings) error {
	query := `
		UPDATE post_express_settings SET
			api_username = $2, api_password = $3, api_endpoint = $4,
			sender_name = $5, sender_address = $6, sender_city = $7, 
			sender_postal_code = $8, sender_phone = $9, sender_email = $10,
			enabled = $11, test_mode = $12, auto_print_labels = $13, auto_track_shipments = $14,
			notify_on_pickup = $15, notify_on_delivery = $16, notify_on_failed_delivery = $17,
			total_shipments = $18, successful_deliveries = $19, failed_deliveries = $20,
			updated_at = NOW()
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		settings.ID, settings.APIUsername, settings.APIPassword, settings.APIEndpoint,
		settings.SenderName, settings.SenderAddress, settings.SenderCity,
		settings.SenderPostalCode, settings.SenderPhone, settings.SenderEmail,
		settings.Enabled, settings.TestMode, settings.AutoPrintLabels, settings.AutoTrackShipments,
		settings.NotifyOnPickup, settings.NotifyOnDelivery, settings.NotifyOnFailedDelivery,
		settings.TotalShipments, settings.SuccessfulDeliveries, settings.FailedDeliveries)
	if err != nil {
		return fmt.Errorf("failed to update settings: %w", err)
	}

	return nil
}

func (r *Repository) CreateSettings(ctx context.Context, settings *models.PostExpressSettings) error {
	query := `
		INSERT INTO post_express_settings (
			api_username, api_password, api_endpoint,
			sender_name, sender_address, sender_city, sender_postal_code, sender_phone, sender_email,
			enabled, test_mode, auto_print_labels, auto_track_shipments,
			notify_on_pickup, notify_on_delivery, notify_on_failed_delivery,
			total_shipments, successful_deliveries, failed_deliveries
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19
		) RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		settings.APIUsername, settings.APIPassword, settings.APIEndpoint,
		settings.SenderName, settings.SenderAddress, settings.SenderCity,
		settings.SenderPostalCode, settings.SenderPhone, settings.SenderEmail,
		settings.Enabled, settings.TestMode, settings.AutoPrintLabels, settings.AutoTrackShipments,
		settings.NotifyOnPickup, settings.NotifyOnDelivery, settings.NotifyOnFailedDelivery,
		settings.TotalShipments, settings.SuccessfulDeliveries, settings.FailedDeliveries,
	).Scan(&settings.ID, &settings.CreatedAt, &settings.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create settings: %w", err)
	}

	return nil
}

// =============================================================================
// ЛОКАЦИИ
// =============================================================================

func (r *Repository) GetLocationByID(ctx context.Context, id int) (*models.PostExpressLocation, error) {
	var location models.PostExpressLocation
	query := `
		SELECT id, post_express_id, name, name_cyrillic, postal_code, municipality,
			   latitude, longitude, region, district, delivery_zone,
			   is_active, supports_cod, supports_express, created_at, updated_at
		FROM post_express_locations 
		WHERE id = $1`

	err := r.db.GetContext(ctx, &location, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get location by ID: %w", err)
	}

	return &location, nil
}

func (r *Repository) GetLocationByPostExpressID(ctx context.Context, postExpressID int) (*models.PostExpressLocation, error) {
	var location models.PostExpressLocation
	query := `
		SELECT id, post_express_id, name, name_cyrillic, postal_code, municipality,
			   latitude, longitude, region, district, delivery_zone,
			   is_active, supports_cod, supports_express, created_at, updated_at
		FROM post_express_locations 
		WHERE post_express_id = $1`

	err := r.db.GetContext(ctx, &location, query, postExpressID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get location by PostExpress ID: %w", err)
	}

	return &location, nil
}

func (r *Repository) SearchLocations(ctx context.Context, query string, limit int) ([]*models.PostExpressLocation, error) {
	locations := []*models.PostExpressLocation{}

	sql := `
		SELECT id, post_express_id, name, name_cyrillic, postal_code, municipality,
			   latitude, longitude, region, district, delivery_zone,
			   is_active, supports_cod, supports_express, created_at, updated_at
		FROM post_express_locations 
		WHERE is_active = true 
		  AND (name ILIKE $1 OR name_cyrillic ILIKE $1 OR postal_code ILIKE $1)
		ORDER BY 
			CASE WHEN name ILIKE $2 THEN 1 ELSE 2 END,
			length(name),
			name
		LIMIT $3`

	searchPattern := "%" + query + "%"
	exactPattern := query + "%"

	err := r.db.SelectContext(ctx, &locations, sql, searchPattern, exactPattern, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search locations: %w", err)
	}

	return locations, nil
}

func (r *Repository) CreateLocation(ctx context.Context, location *models.PostExpressLocation) error {
	query := `
		INSERT INTO post_express_locations (
			post_express_id, name, name_cyrillic, postal_code, municipality,
			latitude, longitude, region, district, delivery_zone,
			is_active, supports_cod, supports_express
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		) RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		location.PostExpressID, location.Name, location.NameCyrillic, location.PostalCode,
		location.Municipality, location.Latitude, location.Longitude, location.Region,
		location.District, location.DeliveryZone, location.IsActive, location.SupportsCOD,
		location.SupportsExpress,
	).Scan(&location.ID, &location.CreatedAt, &location.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create location: %w", err)
	}

	return nil
}

func (r *Repository) UpdateLocation(ctx context.Context, location *models.PostExpressLocation) error {
	query := `
		UPDATE post_express_locations SET
			name = $2, name_cyrillic = $3, postal_code = $4, municipality = $5,
			latitude = $6, longitude = $7, region = $8, district = $9, delivery_zone = $10,
			is_active = $11, supports_cod = $12, supports_express = $13, updated_at = NOW()
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		location.ID, location.Name, location.NameCyrillic, location.PostalCode,
		location.Municipality, location.Latitude, location.Longitude, location.Region,
		location.District, location.DeliveryZone, location.IsActive, location.SupportsCOD,
		location.SupportsExpress)
	if err != nil {
		return fmt.Errorf("failed to update location: %w", err)
	}

	return nil
}

func (r *Repository) BulkUpsertLocations(ctx context.Context, locations []*models.PostExpressLocation) error {
	if len(locations) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO post_express_locations (
			post_express_id, name, name_cyrillic, postal_code, municipality,
			latitude, longitude, region, district, delivery_zone,
			is_active, supports_cod, supports_express
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (post_express_id) DO UPDATE SET
			name = EXCLUDED.name,
			name_cyrillic = EXCLUDED.name_cyrillic,
			postal_code = EXCLUDED.postal_code,
			municipality = EXCLUDED.municipality,
			latitude = EXCLUDED.latitude,
			longitude = EXCLUDED.longitude,
			region = EXCLUDED.region,
			district = EXCLUDED.district,
			delivery_zone = EXCLUDED.delivery_zone,
			is_active = EXCLUDED.is_active,
			supports_cod = EXCLUDED.supports_cod,
			supports_express = EXCLUDED.supports_express,
			updated_at = NOW()`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, location := range locations {
		_, err = stmt.ExecContext(ctx,
			location.PostExpressID, location.Name, location.NameCyrillic, location.PostalCode,
			location.Municipality, location.Latitude, location.Longitude, location.Region,
			location.District, location.DeliveryZone, location.IsActive, location.SupportsCOD,
			location.SupportsExpress)
		if err != nil {
			return fmt.Errorf("failed to upsert location %d: %w", location.PostExpressID, err)
		}
	}

	return tx.Commit()
}

// =============================================================================
// ОТДЕЛЕНИЯ
// =============================================================================

func (r *Repository) GetOfficeByID(ctx context.Context, id int) (*models.PostExpressOffice, error) {
	var office models.PostExpressOffice
	query := `
		SELECT id, office_code, location_id, name, address, phone, email,
			   working_hours, latitude, longitude,
			   accepts_packages, issues_packages, has_atm, has_parking, wheelchair_accessible,
			   is_active, temporary_closed, closed_until, created_at, updated_at
		FROM post_express_offices 
		WHERE id = $1`

	err := r.db.GetContext(ctx, &office, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get office by ID: %w", err)
	}

	return &office, nil
}

func (r *Repository) GetOfficeByCode(ctx context.Context, code string) (*models.PostExpressOffice, error) {
	var office models.PostExpressOffice
	query := `
		SELECT id, office_code, location_id, name, address, phone, email,
			   working_hours, latitude, longitude,
			   accepts_packages, issues_packages, has_atm, has_parking, wheelchair_accessible,
			   is_active, temporary_closed, closed_until, created_at, updated_at
		FROM post_express_offices 
		WHERE office_code = $1`

	err := r.db.GetContext(ctx, &office, query, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get office by code: %w", err)
	}

	return &office, nil
}

func (r *Repository) GetOfficesByLocationID(ctx context.Context, locationID int) ([]*models.PostExpressOffice, error) {
	offices := []*models.PostExpressOffice{}
	query := `
		SELECT id, office_code, location_id, name, address, phone, email,
			   working_hours, latitude, longitude,
			   accepts_packages, issues_packages, has_atm, has_parking, wheelchair_accessible,
			   is_active, temporary_closed, closed_until, created_at, updated_at
		FROM post_express_offices 
		WHERE location_id = $1 AND is_active = true
		ORDER BY name`

	err := r.db.SelectContext(ctx, &offices, query, locationID)
	if err != nil {
		return nil, fmt.Errorf("failed to get offices by location ID: %w", err)
	}

	return offices, nil
}

func (r *Repository) CreateOffice(ctx context.Context, office *models.PostExpressOffice) error {
	query := `
		INSERT INTO post_express_offices (
			office_code, location_id, name, address, phone, email,
			working_hours, latitude, longitude,
			accepts_packages, issues_packages, has_atm, has_parking, wheelchair_accessible,
			is_active, temporary_closed, closed_until
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
		) RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		office.OfficeCode, office.LocationID, office.Name, office.Address, office.Phone,
		office.Email, office.WorkingHours, office.Latitude, office.Longitude,
		office.AcceptsPackages, office.IssuesPackages, office.HasATM, office.HasParking,
		office.WheelchairAccessible, office.IsActive, office.TemporaryClosed, office.ClosedUntil,
	).Scan(&office.ID, &office.CreatedAt, &office.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create office: %w", err)
	}

	return nil
}

func (r *Repository) UpdateOffice(ctx context.Context, office *models.PostExpressOffice) error {
	query := `
		UPDATE post_express_offices SET
			location_id = $2, name = $3, address = $4, phone = $5, email = $6,
			working_hours = $7, latitude = $8, longitude = $9,
			accepts_packages = $10, issues_packages = $11, has_atm = $12, has_parking = $13,
			wheelchair_accessible = $14, is_active = $15, temporary_closed = $16, closed_until = $17,
			updated_at = NOW()
		WHERE office_code = $1`

	_, err := r.db.ExecContext(ctx, query,
		office.OfficeCode, office.LocationID, office.Name, office.Address, office.Phone,
		office.Email, office.WorkingHours, office.Latitude, office.Longitude,
		office.AcceptsPackages, office.IssuesPackages, office.HasATM, office.HasParking,
		office.WheelchairAccessible, office.IsActive, office.TemporaryClosed, office.ClosedUntil)
	if err != nil {
		return fmt.Errorf("failed to update office: %w", err)
	}

	return nil
}

func (r *Repository) BulkUpsertOffices(ctx context.Context, offices []*models.PostExpressOffice) error {
	if len(offices) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO post_express_offices (
			office_code, location_id, name, address, phone, email,
			working_hours, latitude, longitude,
			accepts_packages, issues_packages, has_atm, has_parking, wheelchair_accessible,
			is_active, temporary_closed, closed_until
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
		ON CONFLICT (office_code) DO UPDATE SET
			location_id = EXCLUDED.location_id,
			name = EXCLUDED.name,
			address = EXCLUDED.address,
			phone = EXCLUDED.phone,
			email = EXCLUDED.email,
			working_hours = EXCLUDED.working_hours,
			latitude = EXCLUDED.latitude,
			longitude = EXCLUDED.longitude,
			accepts_packages = EXCLUDED.accepts_packages,
			issues_packages = EXCLUDED.issues_packages,
			has_atm = EXCLUDED.has_atm,
			has_parking = EXCLUDED.has_parking,
			wheelchair_accessible = EXCLUDED.wheelchair_accessible,
			is_active = EXCLUDED.is_active,
			temporary_closed = EXCLUDED.temporary_closed,
			closed_until = EXCLUDED.closed_until,
			updated_at = NOW()`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, office := range offices {
		_, err = stmt.ExecContext(ctx,
			office.OfficeCode, office.LocationID, office.Name, office.Address, office.Phone,
			office.Email, office.WorkingHours, office.Latitude, office.Longitude,
			office.AcceptsPackages, office.IssuesPackages, office.HasATM, office.HasParking,
			office.WheelchairAccessible, office.IsActive, office.TemporaryClosed, office.ClosedUntil)
		if err != nil {
			return fmt.Errorf("failed to upsert office %s: %w", office.OfficeCode, err)
		}
	}

	return tx.Commit()
}

// =============================================================================
// ТАРИФЫ
// =============================================================================

func (r *Repository) GetRates(ctx context.Context) ([]*models.PostExpressRate, error) {
	rates := []*models.PostExpressRate{}
	query := `
		SELECT id, weight_from, weight_to, base_price,
			   insurance_included_up_to, insurance_rate_percent, cod_fee,
			   max_length_cm, max_width_cm, max_height_cm, max_dimensions_sum_cm,
			   delivery_days_min, delivery_days_max, is_active, is_special_offer,
			   created_at, updated_at
		FROM post_express_rates 
		WHERE is_active = true
		ORDER BY weight_from`

	err := r.db.SelectContext(ctx, &rates, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get rates: %w", err)
	}

	return rates, nil
}

func (r *Repository) GetRateForWeight(ctx context.Context, weight float64) (*models.PostExpressRate, error) {
	var rate models.PostExpressRate
	query := `
		SELECT id, weight_from, weight_to, base_price,
			   insurance_included_up_to, insurance_rate_percent, cod_fee,
			   max_length_cm, max_width_cm, max_height_cm, max_dimensions_sum_cm,
			   delivery_days_min, delivery_days_max, is_active, is_special_offer,
			   created_at, updated_at
		FROM post_express_rates 
		WHERE is_active = true 
		  AND weight_from <= $1 
		  AND weight_to >= $1
		ORDER BY weight_from
		LIMIT 1`

	err := r.db.GetContext(ctx, &rate, query, weight)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get rate for weight: %w", err)
	}

	return &rate, nil
}

func (r *Repository) CreateRate(ctx context.Context, rate *models.PostExpressRate) error {
	query := `
		INSERT INTO post_express_rates (
			weight_from, weight_to, base_price,
			insurance_included_up_to, insurance_rate_percent, cod_fee,
			max_length_cm, max_width_cm, max_height_cm, max_dimensions_sum_cm,
			delivery_days_min, delivery_days_max, is_active, is_special_offer
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		) RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		rate.WeightFrom, rate.WeightTo, rate.BasePrice,
		rate.InsuranceIncludedUpTo, rate.InsuranceRatePercent, rate.CODFee,
		rate.MaxLengthCm, rate.MaxWidthCm, rate.MaxHeightCm, rate.MaxDimensionsSumCm,
		rate.DeliveryDaysMin, rate.DeliveryDaysMax, rate.IsActive, rate.IsSpecialOffer,
	).Scan(&rate.ID, &rate.CreatedAt, &rate.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create rate: %w", err)
	}

	return nil
}

func (r *Repository) UpdateRate(ctx context.Context, rate *models.PostExpressRate) error {
	query := `
		UPDATE post_express_rates SET
			weight_from = $2, weight_to = $3, base_price = $4,
			insurance_included_up_to = $5, insurance_rate_percent = $6, cod_fee = $7,
			max_length_cm = $8, max_width_cm = $9, max_height_cm = $10, max_dimensions_sum_cm = $11,
			delivery_days_min = $12, delivery_days_max = $13, is_active = $14, is_special_offer = $15,
			updated_at = NOW()
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		rate.ID, rate.WeightFrom, rate.WeightTo, rate.BasePrice,
		rate.InsuranceIncludedUpTo, rate.InsuranceRatePercent, rate.CODFee,
		rate.MaxLengthCm, rate.MaxWidthCm, rate.MaxHeightCm, rate.MaxDimensionsSumCm,
		rate.DeliveryDaysMin, rate.DeliveryDaysMax, rate.IsActive, rate.IsSpecialOffer)
	if err != nil {
		return fmt.Errorf("failed to update rate: %w", err)
	}

	return nil
}

// =============================================================================
// ОТПРАВЛЕНИЯ
// =============================================================================

func (r *Repository) CreateShipment(ctx context.Context, shipment *models.PostExpressShipment) (*models.PostExpressShipment, error) {
	query := `
		INSERT INTO post_express_shipments (
			marketplace_order_id, storefront_order_id, tracking_number, barcode, post_express_id,
			sender_name, sender_address, sender_city, sender_postal_code, sender_phone, sender_email, sender_location_id,
			recipient_name, recipient_address, recipient_city, recipient_postal_code, recipient_phone, recipient_email, recipient_location_id,
			weight_kg, length_cm, width_cm, height_cm, declared_value,
			service_type, cod_amount, cod_reference, insurance_amount,
			base_price, insurance_fee, cod_fee, total_price,
			status, delivery_status, delivery_instructions, notes
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19,
			$20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36
		) RETURNING id, status_history, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		shipment.MarketplaceOrderID, shipment.StorefrontOrderID, shipment.TrackingNumber,
		shipment.Barcode, shipment.PostExpressID,
		shipment.SenderName, shipment.SenderAddress, shipment.SenderCity, shipment.SenderPostalCode,
		shipment.SenderPhone, shipment.SenderEmail, shipment.SenderLocationID,
		shipment.RecipientName, shipment.RecipientAddress, shipment.RecipientCity, shipment.RecipientPostalCode,
		shipment.RecipientPhone, shipment.RecipientEmail, shipment.RecipientLocationID,
		shipment.WeightKg, shipment.LengthCm, shipment.WidthCm, shipment.HeightCm, shipment.DeclaredValue,
		shipment.ServiceType, shipment.CODAmount, shipment.CODReference, shipment.InsuranceAmount,
		shipment.BasePrice, shipment.InsuranceFee, shipment.CODFee, shipment.TotalPrice,
		shipment.Status, shipment.DeliveryStatus, shipment.DeliveryInstructions, shipment.Notes,
	).Scan(&shipment.ID, &shipment.StatusHistory, &shipment.CreatedAt, &shipment.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create shipment: %w", err)
	}

	return shipment, nil
}

func (r *Repository) GetShipmentByID(ctx context.Context, id int) (*models.PostExpressShipment, error) {
	var shipment models.PostExpressShipment
	query := `
		SELECT id, marketplace_order_id, storefront_order_id, tracking_number, barcode, post_express_id,
			   sender_name, sender_address, sender_city, sender_postal_code, sender_phone, sender_email, sender_location_id,
			   recipient_name, recipient_address, recipient_city, recipient_postal_code, recipient_phone, recipient_email, recipient_location_id,
			   weight_kg, length_cm, width_cm, height_cm, declared_value,
			   service_type, cod_amount, cod_reference, insurance_amount,
			   base_price, insurance_fee, cod_fee, total_price,
			   status, delivery_status, label_url, invoice_url, pod_url,
			   registered_at, picked_up_at, delivered_at, failed_at, returned_at,
			   status_history, notes, internal_notes, delivery_instructions, failed_reason,
			   created_at, updated_at
		FROM post_express_shipments 
		WHERE id = $1`

	err := r.db.GetContext(ctx, &shipment, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get shipment by ID: %w", err)
	}

	return &shipment, nil
}

func (r *Repository) GetShipmentByTrackingNumber(ctx context.Context, trackingNumber string) (*models.PostExpressShipment, error) {
	var shipment models.PostExpressShipment
	query := `
		SELECT id, marketplace_order_id, storefront_order_id, tracking_number, barcode, post_express_id,
			   sender_name, sender_address, sender_city, sender_postal_code, sender_phone, sender_email, sender_location_id,
			   recipient_name, recipient_address, recipient_city, recipient_postal_code, recipient_phone, recipient_email, recipient_location_id,
			   weight_kg, length_cm, width_cm, height_cm, declared_value,
			   service_type, cod_amount, cod_reference, insurance_amount,
			   base_price, insurance_fee, cod_fee, total_price,
			   status, delivery_status, label_url, invoice_url, pod_url,
			   registered_at, picked_up_at, delivered_at, failed_at, returned_at,
			   status_history, notes, internal_notes, delivery_instructions, failed_reason,
			   created_at, updated_at
		FROM post_express_shipments 
		WHERE tracking_number = $1`

	err := r.db.GetContext(ctx, &shipment, query, trackingNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get shipment by tracking number: %w", err)
	}

	return &shipment, nil
}

func (r *Repository) GetShipmentsByOrderID(ctx context.Context, orderID int, isStorefront bool) ([]*models.PostExpressShipment, error) {
	shipments := []*models.PostExpressShipment{}

	var query string
	if isStorefront {
		query = `
			SELECT id, marketplace_order_id, storefront_order_id, tracking_number, barcode, post_express_id,
				   sender_name, sender_address, sender_city, sender_postal_code, sender_phone, sender_email, sender_location_id,
				   recipient_name, recipient_address, recipient_city, recipient_postal_code, recipient_phone, recipient_email, recipient_location_id,
				   weight_kg, length_cm, width_cm, height_cm, declared_value,
				   service_type, cod_amount, cod_reference, insurance_amount,
				   base_price, insurance_fee, cod_fee, total_price,
				   status, delivery_status, label_url, invoice_url, pod_url,
				   registered_at, picked_up_at, delivered_at, failed_at, returned_at,
				   status_history, notes, internal_notes, delivery_instructions, failed_reason,
				   created_at, updated_at
			FROM post_express_shipments 
			WHERE storefront_order_id = $1
			ORDER BY created_at DESC`
	} else {
		query = `
			SELECT id, marketplace_order_id, storefront_order_id, tracking_number, barcode, post_express_id,
				   sender_name, sender_address, sender_city, sender_postal_code, sender_phone, sender_email, sender_location_id,
				   recipient_name, recipient_address, recipient_city, recipient_postal_code, recipient_phone, recipient_email, recipient_location_id,
				   weight_kg, length_cm, width_cm, height_cm, declared_value,
				   service_type, cod_amount, cod_reference, insurance_amount,
				   base_price, insurance_fee, cod_fee, total_price,
				   status, delivery_status, label_url, invoice_url, pod_url,
				   registered_at, picked_up_at, delivered_at, failed_at, returned_at,
				   status_history, notes, internal_notes, delivery_instructions, failed_reason,
				   created_at, updated_at
			FROM post_express_shipments 
			WHERE marketplace_order_id = $1
			ORDER BY created_at DESC`
	}

	err := r.db.SelectContext(ctx, &shipments, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipments by order ID: %w", err)
	}

	return shipments, nil
}

func (r *Repository) UpdateShipment(ctx context.Context, shipment *models.PostExpressShipment) error {
	query := `
		UPDATE post_express_shipments SET
			tracking_number = $2, barcode = $3, post_express_id = $4,
			weight_kg = $5, length_cm = $6, width_cm = $7, height_cm = $8, declared_value = $9,
			service_type = $10, cod_amount = $11, cod_reference = $12, insurance_amount = $13,
			base_price = $14, insurance_fee = $15, cod_fee = $16, total_price = $17,
			status = $18, delivery_status = $19, label_url = $20, invoice_url = $21, pod_url = $22,
			registered_at = $23, picked_up_at = $24, delivered_at = $25, failed_at = $26, returned_at = $27,
			status_history = $28, notes = $29, internal_notes = $30, delivery_instructions = $31,
			failed_reason = $32, updated_at = NOW()
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		shipment.ID, shipment.TrackingNumber, shipment.Barcode, shipment.PostExpressID,
		shipment.WeightKg, shipment.LengthCm, shipment.WidthCm, shipment.HeightCm, shipment.DeclaredValue,
		shipment.ServiceType, shipment.CODAmount, shipment.CODReference, shipment.InsuranceAmount,
		shipment.BasePrice, shipment.InsuranceFee, shipment.CODFee, shipment.TotalPrice,
		shipment.Status, shipment.DeliveryStatus, shipment.LabelURL, shipment.InvoiceURL, shipment.PODURL,
		shipment.RegisteredAt, shipment.PickedUpAt, shipment.DeliveredAt, shipment.FailedAt, shipment.ReturnedAt,
		shipment.StatusHistory, shipment.Notes, shipment.InternalNotes, shipment.DeliveryInstructions,
		shipment.FailedReason)
	if err != nil {
		return fmt.Errorf("failed to update shipment: %w", err)
	}

	return nil
}

func (r *Repository) UpdateShipmentStatus(ctx context.Context, id int, status models.ShipmentStatus) error {
	query := `UPDATE post_express_shipments SET status = $2, updated_at = NOW() WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("failed to update shipment status: %w", err)
	}

	return nil
}

func (r *Repository) ListShipments(ctx context.Context, filters storage.ShipmentFilters) ([]*models.PostExpressShipment, int, error) {
	// Базовый запрос
	baseQuery := `
		FROM post_express_shipments s
		WHERE 1=1`

	var conditions []string
	var args []interface{}
	argIndex := 1

	// Применяем фильтры
	if filters.Status != nil {
		conditions = append(conditions, fmt.Sprintf("s.status = $%d", argIndex))
		args = append(args, *filters.Status)
		argIndex++
	}

	if filters.DateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("s.created_at >= $%d", argIndex))
		args = append(args, *filters.DateFrom)
		argIndex++
	}

	if filters.DateTo != nil {
		conditions = append(conditions, fmt.Sprintf("s.created_at <= $%d", argIndex))
		args = append(args, *filters.DateTo)
		argIndex++
	}

	if filters.OrderID != nil {
		conditions = append(conditions, fmt.Sprintf("(s.marketplace_order_id = $%d OR s.storefront_order_id = $%d)", argIndex, argIndex))
		args = append(args, *filters.OrderID)
		argIndex++
	}

	if filters.City != nil {
		conditions = append(conditions, fmt.Sprintf("s.recipient_city ILIKE $%d", argIndex))
		args = append(args, "%"+*filters.City+"%")
		argIndex++
	}

	// Добавляем условия к запросу
	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	// Получаем общее количество
	countQuery := "SELECT COUNT(*) " + baseQuery
	var totalCount int
	err := r.db.GetContext(ctx, &totalCount, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get shipments count: %w", err)
	}

	// Основной запрос с данными
	dataQuery := `
		SELECT id, marketplace_order_id, storefront_order_id, tracking_number, barcode, post_express_id,
			   sender_name, sender_address, sender_city, sender_postal_code, sender_phone, sender_email, sender_location_id,
			   recipient_name, recipient_address, recipient_city, recipient_postal_code, recipient_phone, recipient_email, recipient_location_id,
			   weight_kg, length_cm, width_cm, height_cm, declared_value,
			   service_type, cod_amount, cod_reference, insurance_amount,
			   base_price, insurance_fee, cod_fee, total_price,
			   status, delivery_status, label_url, invoice_url, pod_url,
			   registered_at, picked_up_at, delivered_at, failed_at, returned_at,
			   status_history, notes, internal_notes, delivery_instructions, failed_reason,
			   created_at, updated_at ` + baseQuery + `
		ORDER BY s.created_at DESC
		LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)

	// Добавляем параметры пагинации
	limit := filters.PageSize
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset := (filters.Page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	args = append(args, limit, offset)

	shipments := []*models.PostExpressShipment{}
	err = r.db.SelectContext(ctx, &shipments, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get shipments: %w", err)
	}

	return shipments, totalCount, nil
}

// =============================================================================
// ОТСЛЕЖИВАНИЕ
// =============================================================================

func (r *Repository) CreateTrackingEvent(ctx context.Context, event *models.TrackingEvent) error {
	query := `
		INSERT INTO post_express_tracking_events (
			shipment_id, event_code, event_description, event_location, event_timestamp, additional_info
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		event.ShipmentID, event.EventCode, event.EventDescription,
		event.EventLocation, event.EventTimestamp, event.AdditionalInfo,
	).Scan(&event.ID, &event.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create tracking event: %w", err)
	}

	return nil
}

func (r *Repository) GetTrackingEventsByShipmentID(ctx context.Context, shipmentID int) ([]*models.TrackingEvent, error) {
	events := []*models.TrackingEvent{}
	query := `
		SELECT id, shipment_id, event_code, event_description, event_location, 
			   event_timestamp, additional_info, created_at
		FROM post_express_tracking_events 
		WHERE shipment_id = $1
		ORDER BY event_timestamp DESC, created_at DESC`

	err := r.db.SelectContext(ctx, &events, query, shipmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tracking events: %w", err)
	}

	return events, nil
}

func (r *Repository) BulkCreateTrackingEvents(ctx context.Context, events []*models.TrackingEvent) error {
	if len(events) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO post_express_tracking_events (
			shipment_id, event_code, event_description, event_location, event_timestamp, additional_info
		) VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (shipment_id, event_code, event_timestamp) DO NOTHING`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, event := range events {
		_, err = stmt.ExecContext(ctx,
			event.ShipmentID, event.EventCode, event.EventDescription,
			event.EventLocation, event.EventTimestamp, event.AdditionalInfo)
		if err != nil {
			return fmt.Errorf("failed to insert tracking event: %w", err)
		}
	}

	return tx.Commit()
}

// =============================================================================
// API ЛОГИ
// =============================================================================

func (r *Repository) CreateAPILog(ctx context.Context, log *storage.APILog) error {
	query := `
		INSERT INTO post_express_api_logs (
			transaction_id, transaction_type, request_data, response_data, 
			status, error_message, shipment_id, execution_time_ms
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		log.TransactionID, log.TransactionType, log.RequestData, log.ResponseData,
		log.Status, log.ErrorMessage, log.ShipmentID, log.ExecutionTimeMs,
	).Scan(&log.ID, &log.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create API log: %w", err)
	}

	return nil
}

func (r *Repository) GetAPILogsByShipmentID(ctx context.Context, shipmentID int) ([]*storage.APILog, error) {
	logs := []*storage.APILog{}
	query := `
		SELECT id, transaction_id, transaction_type, request_data, response_data,
			   status, error_message, shipment_id, execution_time_ms, created_at
		FROM post_express_api_logs 
		WHERE shipment_id = $1
		ORDER BY created_at DESC`

	err := r.db.SelectContext(ctx, &logs, query, shipmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get API logs: %w", err)
	}

	return logs, nil
}

func (r *Repository) CleanupOldAPILogs(ctx context.Context, daysToKeep int) error {
	query := `DELETE FROM post_express_api_logs WHERE created_at < NOW() - INTERVAL '%d days'`

	_, err := r.db.ExecContext(ctx, fmt.Sprintf(query, daysToKeep))
	if err != nil {
		return fmt.Errorf("failed to cleanup old API logs: %w", err)
	}

	return nil
}

// =============================================================================
// СКЛАД
// =============================================================================

func (r *Repository) GetWarehouses(ctx context.Context) ([]*models.Warehouse, error) {
	warehouses := []*models.Warehouse{}
	query := `
		SELECT id, code, name, type, address, city, postal_code, country,
			   phone, email, manager_name, manager_phone, latitude, longitude,
			   working_hours, total_area_m2, storage_area_m2, max_capacity_m3, current_occupancy_m3,
			   supports_fbs, supports_pickup, has_refrigeration, has_loading_dock, is_active,
			   created_at, updated_at
		FROM warehouses 
		WHERE is_active = true
		ORDER BY name`

	err := r.db.SelectContext(ctx, &warehouses, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get warehouses: %w", err)
	}

	return warehouses, nil
}

func (r *Repository) GetWarehouseByID(ctx context.Context, id int) (*models.Warehouse, error) {
	var warehouse models.Warehouse
	query := `
		SELECT id, code, name, type, address, city, postal_code, country,
			   phone, email, manager_name, manager_phone, latitude, longitude,
			   working_hours, total_area_m2, storage_area_m2, max_capacity_m3, current_occupancy_m3,
			   supports_fbs, supports_pickup, has_refrigeration, has_loading_dock, is_active,
			   created_at, updated_at
		FROM warehouses 
		WHERE id = $1`

	err := r.db.GetContext(ctx, &warehouse, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get warehouse by ID: %w", err)
	}

	return &warehouse, nil
}

func (r *Repository) GetWarehouseByCode(ctx context.Context, code string) (*models.Warehouse, error) {
	var warehouse models.Warehouse
	query := `
		SELECT id, code, name, type, address, city, postal_code, country,
			   phone, email, manager_name, manager_phone, latitude, longitude,
			   working_hours, total_area_m2, storage_area_m2, max_capacity_m3, current_occupancy_m3,
			   supports_fbs, supports_pickup, has_refrigeration, has_loading_dock, is_active,
			   created_at, updated_at
		FROM warehouses 
		WHERE code = $1`

	err := r.db.GetContext(ctx, &warehouse, query, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get warehouse by code: %w", err)
	}

	return &warehouse, nil
}

// =============================================================================
// ЗАКАЗЫ НА САМОВЫВОЗ
// =============================================================================

func (r *Repository) CreatePickupOrder(ctx context.Context, order *models.WarehousePickupOrder) (*models.WarehousePickupOrder, error) {
	// Генерируем уникальный код
	var pickupCode string
	for attempts := 0; attempts < 10; attempts++ {
		pickupCode = generatePickupCode()
		// Проверяем уникальность
		var exists bool
		err := r.db.GetContext(ctx, &exists, "SELECT EXISTS(SELECT 1 FROM warehouse_pickup_orders WHERE pickup_code = $1)", pickupCode)
		if err != nil {
			return nil, fmt.Errorf("failed to check pickup code uniqueness: %w", err)
		}
		if !exists {
			break
		}
	}

	order.PickupCode = pickupCode
	order.Status = models.PickupOrderStatusPending

	// Устанавливаем срок истечения (7 дней)
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	order.ExpiresAt = &expiresAt

	query := `
		INSERT INTO warehouse_pickup_orders (
			warehouse_id, marketplace_order_id, storefront_order_id, pickup_code,
			status, expires_at, customer_name, customer_phone, customer_email, notes
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, qr_code_url, ready_at, picked_up_at, pickup_confirmed_by,
				  id_document_type, id_document_number, signature_url,
				  notification_sent_at, reminder_sent_at, pickup_photo_url,
				  created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		order.WarehouseID, order.MarketplaceOrderID, order.StorefrontOrderID, order.PickupCode,
		order.Status, order.ExpiresAt, order.CustomerName, order.CustomerPhone, order.CustomerEmail, order.Notes,
	).Scan(&order.ID, &order.QRCodeURL, &order.ReadyAt, &order.PickedUpAt, &order.PickupConfirmedBy,
		&order.IDDocumentType, &order.IDDocumentNumber, &order.SignatureURL,
		&order.NotificationSentAt, &order.ReminderSentAt, &order.PickupPhotoURL,
		&order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create pickup order: %w", err)
	}

	return order, nil
}

func (r *Repository) GetPickupOrderByID(ctx context.Context, id int) (*models.WarehousePickupOrder, error) {
	var order models.WarehousePickupOrder
	query := `
		SELECT id, warehouse_id, marketplace_order_id, storefront_order_id, pickup_code, qr_code_url,
			   status, ready_at, picked_up_at, expires_at, customer_name, customer_phone, customer_email,
			   pickup_confirmed_by, id_document_type, id_document_number, signature_url,
			   notification_sent_at, reminder_sent_at, notes, pickup_photo_url,
			   created_at, updated_at
		FROM warehouse_pickup_orders 
		WHERE id = $1`

	err := r.db.GetContext(ctx, &order, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get pickup order by ID: %w", err)
	}

	return &order, nil
}

func (r *Repository) GetPickupOrderByCode(ctx context.Context, code string) (*models.WarehousePickupOrder, error) {
	var order models.WarehousePickupOrder
	query := `
		SELECT id, warehouse_id, marketplace_order_id, storefront_order_id, pickup_code, qr_code_url,
			   status, ready_at, picked_up_at, expires_at, customer_name, customer_phone, customer_email,
			   pickup_confirmed_by, id_document_type, id_document_number, signature_url,
			   notification_sent_at, reminder_sent_at, notes, pickup_photo_url,
			   created_at, updated_at
		FROM warehouse_pickup_orders 
		WHERE pickup_code = $1`

	err := r.db.GetContext(ctx, &order, query, code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get pickup order by code: %w", err)
	}

	return &order, nil
}

func (r *Repository) UpdatePickupOrder(ctx context.Context, order *models.WarehousePickupOrder) error {
	query := `
		UPDATE warehouse_pickup_orders SET
			qr_code_url = $2, status = $3, ready_at = $4, picked_up_at = $5, expires_at = $6,
			customer_name = $7, customer_phone = $8, customer_email = $9,
			pickup_confirmed_by = $10, id_document_type = $11, id_document_number = $12, signature_url = $13,
			notification_sent_at = $14, reminder_sent_at = $15, notes = $16, pickup_photo_url = $17,
			updated_at = NOW()
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		order.ID, order.QRCodeURL, order.Status, order.ReadyAt, order.PickedUpAt, order.ExpiresAt,
		order.CustomerName, order.CustomerPhone, order.CustomerEmail,
		order.PickupConfirmedBy, order.IDDocumentType, order.IDDocumentNumber, order.SignatureURL,
		order.NotificationSentAt, order.ReminderSentAt, order.Notes, order.PickupPhotoURL)
	if err != nil {
		return fmt.Errorf("failed to update pickup order: %w", err)
	}

	return nil
}

func (r *Repository) ListPickupOrders(ctx context.Context, filters storage.PickupOrderFilters) ([]*models.WarehousePickupOrder, int, error) {
	// Базовый запрос
	baseQuery := `
		FROM warehouse_pickup_orders p
		WHERE 1=1`

	var conditions []string
	var args []interface{}
	argIndex := 1

	// Применяем фильтры
	if filters.WarehouseID != nil {
		conditions = append(conditions, fmt.Sprintf("p.warehouse_id = $%d", argIndex))
		args = append(args, *filters.WarehouseID)
		argIndex++
	}

	if filters.Status != nil {
		conditions = append(conditions, fmt.Sprintf("p.status = $%d", argIndex))
		args = append(args, *filters.Status)
		argIndex++
	}

	if filters.DateFrom != nil {
		conditions = append(conditions, fmt.Sprintf("p.created_at >= $%d", argIndex))
		args = append(args, *filters.DateFrom)
		argIndex++
	}

	if filters.DateTo != nil {
		conditions = append(conditions, fmt.Sprintf("p.created_at <= $%d", argIndex))
		args = append(args, *filters.DateTo)
		argIndex++
	}

	// Добавляем условия к запросу
	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	// Получаем общее количество
	countQuery := "SELECT COUNT(*) " + baseQuery
	var totalCount int
	err := r.db.GetContext(ctx, &totalCount, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get pickup orders count: %w", err)
	}

	// Основной запрос с данными
	dataQuery := `
		SELECT id, warehouse_id, marketplace_order_id, storefront_order_id, pickup_code, qr_code_url,
			   status, ready_at, picked_up_at, expires_at, customer_name, customer_phone, customer_email,
			   pickup_confirmed_by, id_document_type, id_document_number, signature_url,
			   notification_sent_at, reminder_sent_at, notes, pickup_photo_url,
			   created_at, updated_at ` + baseQuery + `
		ORDER BY p.created_at DESC
		LIMIT $` + strconv.Itoa(argIndex) + ` OFFSET $` + strconv.Itoa(argIndex+1)

	// Добавляем параметры пагинации
	limit := filters.PageSize
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset := (filters.Page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	args = append(args, limit, offset)

	orders := []*models.WarehousePickupOrder{}
	err = r.db.SelectContext(ctx, &orders, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get pickup orders: %w", err)
	}

	return orders, totalCount, nil
}

// =============================================================================
// СТАТИСТИКА
// =============================================================================

func (r *Repository) GetShipmentStatistics(ctx context.Context, filters storage.StatisticsFilters) (*storage.ShipmentStatistics, error) {
	// Базовый WHERE для фильтров
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if filters.DateFrom != nil {
		whereClause += fmt.Sprintf(" AND created_at >= $%d", argIndex)
		args = append(args, *filters.DateFrom)
		argIndex++
	}

	if filters.DateTo != nil {
		whereClause += fmt.Sprintf(" AND created_at <= $%d", argIndex)
		args = append(args, *filters.DateTo)
	}

	// Основная статистика
	query := fmt.Sprintf(`
		SELECT 
			COUNT(*) as total_shipments,
			COUNT(CASE WHEN status = 'delivered' THEN 1 END) as delivered_shipments,
			COUNT(CASE WHEN status IN ('registered', 'picked_up', 'in_transit') THEN 1 END) as in_transit_shipments,
			COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_shipments,
			COALESCE(SUM(total_price), 0) as total_revenue,
			COALESCE(AVG(EXTRACT(EPOCH FROM (delivered_at - created_at))/3600), 0) as avg_delivery_hours
		FROM post_express_shipments 
		%s`, whereClause)

	var stats storage.ShipmentStatistics
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&stats.TotalShipments, &stats.DeliveredShipments, &stats.InTransitShipments,
		&stats.FailedShipments, &stats.TotalRevenue, &stats.AverageDeliveryTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipment statistics: %w", err)
	}

	// Расчет процента успешности
	if stats.TotalShipments > 0 {
		stats.DeliverySuccessRate = float64(stats.DeliveredShipments) / float64(stats.TotalShipments) * 100
	}

	// Статистика по статусам
	statusQuery := fmt.Sprintf(`
		SELECT status, COUNT(*)
		FROM post_express_shipments 
		%s
		GROUP BY status`, whereClause)

	rows, err := r.db.QueryContext(ctx, statusQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get status statistics: %w", err)
	}
	defer rows.Close()

	stats.ByStatus = make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, fmt.Errorf("failed to scan status row: %w", err)
		}
		stats.ByStatus[status] = count
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate status statistics: %w", err)
	}

	// Статистика по городам
	cityQuery := fmt.Sprintf(`
		SELECT recipient_city, COUNT(*)
		FROM post_express_shipments 
		%s
		GROUP BY recipient_city
		ORDER BY COUNT(*) DESC
		LIMIT 10`, whereClause)

	rows, err = r.db.QueryContext(ctx, cityQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get city statistics: %w", err)
	}
	defer rows.Close()

	stats.ByCity = make(map[string]int)
	for rows.Next() {
		var city string
		var count int
		if err := rows.Scan(&city, &count); err != nil {
			return nil, fmt.Errorf("failed to scan city row: %w", err)
		}
		stats.ByCity[city] = count
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate city statistics: %w", err)
	}

	return &stats, nil
}

func (r *Repository) GetWarehouseStatistics(ctx context.Context, warehouseID int) (*storage.WarehouseStatistics, error) {
	query := `
		SELECT 
			w.current_occupancy_m3,
			w.max_capacity_m3,
			COALESCE(pickup_stats.total_orders, 0) as total_pickup_orders,
			COALESCE(pickup_stats.pending_orders, 0) as pending_pickup_orders,
			COALESCE(pickup_stats.completed_orders, 0) as completed_pickup_orders,
			COALESCE(pickup_stats.avg_pickup_hours, 0) as avg_pickup_hours
		FROM warehouses w
		LEFT JOIN (
			SELECT 
				warehouse_id,
				COUNT(*) as total_orders,
				COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_orders,
				COUNT(CASE WHEN status = 'picked_up' THEN 1 END) as completed_orders,
				AVG(EXTRACT(EPOCH FROM (picked_up_at - ready_at))/3600) as avg_pickup_hours
			FROM warehouse_pickup_orders
			WHERE warehouse_id = $1
			GROUP BY warehouse_id
		) pickup_stats ON w.id = pickup_stats.warehouse_id
		WHERE w.id = $1`

	var stats storage.WarehouseStatistics
	var currentOccupancy, maxCapacity sql.NullFloat64

	err := r.db.QueryRowContext(ctx, query, warehouseID).Scan(
		&currentOccupancy, &maxCapacity, &stats.TotalPickupOrders,
		&stats.PendingPickupOrders, &stats.CompletedPickupOrders, &stats.AveragePickupTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get warehouse statistics: %w", err)
	}

	// Обработка nullable полей
	if currentOccupancy.Valid {
		stats.TotalVolumeM3 = currentOccupancy.Float64
	}

	if maxCapacity.Valid && maxCapacity.Float64 > 0 {
		stats.OccupancyPercent = (stats.TotalVolumeM3 / maxCapacity.Float64) * 100
	}

	// Подсчет товаров на складе (если есть таблица inventory)
	inventoryQuery := `
		SELECT COALESCE(SUM(quantity_total), 0)
		FROM warehouse_inventory 
		WHERE warehouse_id = $1`

	err = r.db.GetContext(ctx, &stats.TotalInventoryItems, inventoryQuery, warehouseID)
	if err != nil {
		// Не критичная ошибка, продолжаем без inventory данных
		stats.TotalInventoryItems = 0
	}

	return &stats, nil
}

// =============================================================================
// ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ
// =============================================================================

// generatePickupCode генерирует случайный 6-значный код для самовывоза
func generatePickupCode() string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 6)
	for i := range result {
		result[i] = chars[time.Now().UnixNano()%int64(len(chars))]
	}
	return string(result)
}
