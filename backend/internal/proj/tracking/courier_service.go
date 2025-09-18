package tracking

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"backend/internal/storage/postgres"
)

// CourierService управляет курьерами
type CourierService struct {
	db *postgres.Database
}

// NewCourierService создаёт новый сервис курьеров
func NewCourierService(db *postgres.Database) *CourierService {
	return &CourierService{db: db}
}

// Courier представляет курьера
type Courier struct {
	ID                 int          `json:"id" db:"id"`
	Name               string       `json:"name" db:"name"`
	Phone              string       `json:"phone" db:"phone"`
	Email              string       `json:"email" db:"email"`
	Status             string       `json:"status" db:"status"`
	CurrentLatitude    *float64     `json:"current_latitude" db:"current_latitude"`
	CurrentLongitude   *float64     `json:"current_longitude" db:"current_longitude"`
	LastLocationUpdate sql.NullTime `json:"last_location_update" db:"last_location_update"`
	IsOnline           bool         `json:"is_online" db:"is_online"`
	ActiveDeliveries   int          `json:"active_deliveries"`
}

// LocationUpdate представляет обновление локации
type LocationUpdate struct {
	CourierID int     `json:"courier_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Speed     float64 `json:"speed"`
	Heading   int     `json:"heading"`
}

// GetCourier получает информацию о курьере
func (s *CourierService) GetCourier(ctx context.Context, courierID int) (*Courier, error) {
	courier := &Courier{}

	query := `
		SELECT
			id, name, phone, email, status,
			current_latitude, current_longitude, last_location_update,
			is_online
		FROM couriers
		WHERE id = $1`

	err := s.db.GetSQLXDB().GetContext(ctx, courier, query, courierID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("courier not found")
		}
		return nil, fmt.Errorf("failed to get courier: %w", err)
	}

	// Получаем количество активных доставок
	activeDeliveries, err := s.getActiveDeliveriesCount(ctx, courierID)
	if err == nil {
		courier.ActiveDeliveries = activeDeliveries
	}

	return courier, nil
}

// UpdateCourierLocation обновляет местоположение курьера
func (s *CourierService) UpdateCourierLocation(ctx context.Context, courierID int, latitude, longitude float64) error {
	query := `
		UPDATE couriers
		SET current_latitude = $2,
		    current_longitude = $3,
		    last_location_update = NOW(),
		    is_online = true
		WHERE id = $1`

	result, err := s.db.GetSQLXDB().ExecContext(ctx, query, courierID, latitude, longitude)
	if err != nil {
		return fmt.Errorf("failed to update courier location: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("courier not found")
	}

	return nil
}

// GetAllCouriers получает список всех курьеров
func (s *CourierService) GetAllCouriers(ctx context.Context) ([]Courier, error) {
	var couriers []Courier

	query := `
		SELECT
			id, name, phone, email, status,
			current_latitude, current_longitude, last_location_update,
			is_online
		FROM couriers
		ORDER BY name ASC`

	err := s.db.GetSQLXDB().SelectContext(ctx, &couriers, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get couriers: %w", err)
	}

	// Добавляем количество активных доставок для каждого курьера
	for i := range couriers {
		activeDeliveries, err := s.getActiveDeliveriesCount(ctx, couriers[i].ID)
		if err == nil {
			couriers[i].ActiveDeliveries = activeDeliveries
		}
	}

	return couriers, nil
}

// CreateCourier создает нового курьера
func (s *CourierService) CreateCourier(ctx context.Context, name, phone, email string) (*Courier, error) {
	courier := &Courier{}

	query := `
		INSERT INTO couriers (name, phone, email, status, is_online)
		VALUES ($1, $2, $3, 'active', false)
		RETURNING id, name, phone, email, status, current_latitude,
		          current_longitude, last_location_update, is_online`

	err := s.db.GetSQLXDB().GetContext(ctx, courier, query, name, phone, email)
	if err != nil {
		return nil, fmt.Errorf("failed to create courier: %w", err)
	}

	courier.ActiveDeliveries = 0

	return courier, nil
}

// UpdateCourierStatus обновляет статус курьера
func (s *CourierService) UpdateCourierStatus(ctx context.Context, courierID int, status string) error {
	query := `
		UPDATE couriers
		SET status = $2
		WHERE id = $1`

	result, err := s.db.GetSQLXDB().ExecContext(ctx, query, courierID, status)
	if err != nil {
		return fmt.Errorf("failed to update courier status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("courier not found")
	}

	return nil
}

// SetCourierOnlineStatus устанавливает онлайн статус курьера
func (s *CourierService) SetCourierOnlineStatus(ctx context.Context, courierID int, isOnline bool) error {
	query := `
		UPDATE couriers
		SET is_online = $2
		WHERE id = $1`

	result, err := s.db.GetSQLXDB().ExecContext(ctx, query, courierID, isOnline)
	if err != nil {
		return fmt.Errorf("failed to update courier online status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("courier not found")
	}

	return nil
}

// GetNearestCouriers получает ближайших курьеров к указанной точке
func (s *CourierService) GetNearestCouriers(ctx context.Context, latitude, longitude float64, radiusKm int) ([]Courier, error) {
	var couriers []Courier

	query := `
		SELECT
			id, name, phone, email, status,
			current_latitude, current_longitude, last_location_update,
			is_online,
			CAST(ST_Distance(
				ST_MakePoint($2, $1)::geography,
				ST_MakePoint(current_longitude, current_latitude)::geography
			) / 1000 AS INTEGER) as distance_km
		FROM couriers
		WHERE status = 'active'
		  AND is_online = true
		  AND current_latitude IS NOT NULL
		  AND current_longitude IS NOT NULL
		  AND ST_Distance(
		      ST_MakePoint($2, $1)::geography,
		      ST_MakePoint(current_longitude, current_latitude)::geography
		  ) <= $3 * 1000
		ORDER BY distance_km ASC
		LIMIT 10`

	err := s.db.GetSQLXDB().SelectContext(ctx, &couriers, query, latitude, longitude, radiusKm)
	if err != nil {
		return nil, fmt.Errorf("failed to get nearest couriers: %w", err)
	}

	return couriers, nil
}

// getActiveDeliveriesCount получает количество активных доставок курьера
func (s *CourierService) getActiveDeliveriesCount(ctx context.Context, courierID int) (int, error) {
	var count int

	query := `
		SELECT COUNT(*)
		FROM deliveries
		WHERE courier_id = $1
		  AND status IN ('pending', 'picked_up', 'in_transit')`

	err := s.db.GetSQLXDB().GetContext(ctx, &count, query, courierID)
	if err != nil {
		return 0, err
	}

	return count, nil
}
