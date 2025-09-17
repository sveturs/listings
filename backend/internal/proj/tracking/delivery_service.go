package tracking

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"math"
	"time"

	"backend/internal/storage/postgres"
)

// DeliveryService управляет доставками
type DeliveryService struct {
	db *postgres.Database
}

// NewDeliveryService создаёт новый сервис доставок
func NewDeliveryService(db *postgres.Database) *DeliveryService {
	return &DeliveryService{db: db}
}

// Delivery представляет доставку
type Delivery struct {
	ID                    int
	OrderID               int
	CourierID             int
	TrackingToken         string
	Status                string
	PickupAddress         string
	DeliveryAddress       string
	PickupLatitude        float64
	PickupLongitude       float64
	DeliveryLatitude      float64
	DeliveryLongitude     float64
	EstimatedDeliveryTime time.Time
	ActualDeliveryTime    *time.Time
	CourierName           string
	CourierPhone          string
	CourierLocation       *CourierLocation
	LastKnownLocation     *LocationUpdate
	Distance              int // в метрах
	Duration              int // в минутах
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// CourierLocation текущая локация курьера
type CourierLocation struct {
	Latitude  float64
	Longitude float64
	Speed     float64
	Heading   int
	UpdatedAt time.Time
}

// CreateDelivery создаёт новую доставку
func (s *DeliveryService) CreateDelivery(ctx context.Context, orderID, courierID int, pickup, delivery string, pickupLat, pickupLng, deliveryLat, deliveryLng float64) (*Delivery, error) {
	// Генерируем уникальный токен для трекинга
	token := s.generateTrackingToken()

	// Рассчитываем расстояние и примерное время
	distance := s.CalculateDistance(pickupLat, pickupLng, deliveryLat, deliveryLng)
	duration := s.estimateDuration(distance)
	estimatedTime := time.Now().Add(time.Duration(duration) * time.Minute)

	query := `
		INSERT INTO deliveries (
			order_id, courier_id, tracking_token, status,
			pickup_address, delivery_address,
			pickup_latitude, pickup_longitude,
			delivery_latitude, delivery_longitude,
			estimated_delivery_time
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`

	var d Delivery
	err := s.db.GetSQLXDB().QueryRowContext(ctx, query,
		orderID, courierID, token, "pending",
		pickup, delivery,
		pickupLat, pickupLng,
		deliveryLat, deliveryLng,
		estimatedTime,
	).Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create delivery: %w", err)
	}

	d.OrderID = orderID
	d.CourierID = courierID
	d.TrackingToken = token
	d.Status = "pending"
	d.PickupAddress = pickup
	d.DeliveryAddress = delivery
	d.PickupLatitude = pickupLat
	d.PickupLongitude = pickupLng
	d.DeliveryLatitude = deliveryLat
	d.DeliveryLongitude = deliveryLng
	d.EstimatedDeliveryTime = estimatedTime
	d.Distance = int(distance)
	d.Duration = duration

	return &d, nil
}

// GetDelivery получает информацию о доставке
func (s *DeliveryService) GetDelivery(deliveryID int) (*Delivery, error) {
	query := `
		SELECT 
			d.id, d.order_id, d.courier_id, d.tracking_token, d.status,
			d.pickup_address, d.delivery_address,
			d.pickup_latitude, d.pickup_longitude,
			d.delivery_latitude, d.delivery_longitude,
			d.estimated_delivery_time, d.actual_delivery_time,
			d.created_at, d.updated_at,
			c.name, c.phone
		FROM deliveries d
		JOIN couriers c ON c.id = d.courier_id
		WHERE d.id = $1
	`

	var d Delivery
	err := s.db.GetSQLXDB().QueryRow(query, deliveryID).Scan(
		&d.ID, &d.OrderID, &d.CourierID, &d.TrackingToken, &d.Status,
		&d.PickupAddress, &d.DeliveryAddress,
		&d.PickupLatitude, &d.PickupLongitude,
		&d.DeliveryLatitude, &d.DeliveryLongitude,
		&d.EstimatedDeliveryTime, &d.ActualDeliveryTime,
		&d.CreatedAt, &d.UpdatedAt,
		&d.CourierName, &d.CourierPhone,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("delivery not found")
		}
		return nil, err
	}

	// Получаем последнюю известную локацию
	d.CourierLocation = s.getLastCourierLocation(d.CourierID)

	return &d, nil
}

// ValidateTrackingToken проверяет токен трекинга
func (s *DeliveryService) ValidateTrackingToken(token string) (*Delivery, error) {
	query := `
		SELECT 
			d.id, d.order_id, d.courier_id, d.status,
			d.pickup_address, d.delivery_address,
			d.pickup_latitude, d.pickup_longitude,
			d.delivery_latitude, d.delivery_longitude,
			d.estimated_delivery_time,
			c.name
		FROM deliveries d
		JOIN couriers c ON c.id = d.courier_id
		WHERE d.tracking_token = $1
	`

	var d Delivery
	err := s.db.GetSQLXDB().QueryRow(query, token).Scan(
		&d.ID, &d.OrderID, &d.CourierID, &d.Status,
		&d.PickupAddress, &d.DeliveryAddress,
		&d.PickupLatitude, &d.PickupLongitude,
		&d.DeliveryLatitude, &d.DeliveryLongitude,
		&d.EstimatedDeliveryTime,
		&d.CourierName,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invalid tracking token")
		}
		return nil, err
	}

	d.TrackingToken = token
	return &d, nil
}

// UpdateDeliveryStatus обновляет статус доставки
func (s *DeliveryService) UpdateDeliveryStatus(ctx context.Context, deliveryID int, status string) error {
	query := `UPDATE deliveries SET status = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	
	if status == "delivered" {
		query = `UPDATE deliveries SET status = $2, actual_delivery_time = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	}

	_, err := s.db.GetSQLXDB().ExecContext(ctx, query, deliveryID, status)
	return err
}

// UpdateCourierLocation обновляет локацию курьера
func (s *DeliveryService) UpdateCourierLocation(deliveryID int, lat, lng, speed float64, heading int) error {
	// Сохраняем в историю локаций
	query := `
		INSERT INTO delivery_location_history (
			delivery_id, latitude, longitude, speed, heading
		) VALUES ($1, $2, $3, $4, $5)
	`

	_, err := s.db.GetSQLXDB().Exec(query, deliveryID, lat, lng, speed, heading)
	if err != nil {
		return fmt.Errorf("failed to save location history: %w", err)
	}

	// Обновляем текущую локацию курьера
	var courierID int
	err = s.db.GetSQLXDB().QueryRow("SELECT courier_id FROM deliveries WHERE id = $1", deliveryID).Scan(&courierID)
	if err != nil {
		return err
	}

	updateQuery := `
		UPDATE couriers 
		SET current_latitude = $2, 
		    current_longitude = $3,
		    last_location_update = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	_, err = s.db.GetSQLXDB().Exec(updateQuery, courierID, lat, lng)
	return err
}

// GetNearbyItems получает товары рядом с указанной точкой
func (s *DeliveryService) GetNearbyItems(lat, lng float64, radiusMeters int) ([]map[string]interface{}, error) {
	// Конвертируем радиус в градусы (приблизительно)
	radiusDegrees := float64(radiusMeters) / 111000.0

	query := `
		SELECT 
			l.id, l.title, l.price, l.currency,
			l.location_latitude, l.location_longitude,
			CAST(ST_Distance(
				ST_MakePoint($1, $2)::geography,
				ST_MakePoint(l.location_longitude, l.location_latitude)::geography
			) AS INTEGER) as distance_meters,
			s.name as store_name,
			s.logo_url as store_logo
		FROM marketplace_listings l
		LEFT JOIN storefronts s ON s.id = l.storefront_id
		WHERE l.status = 'active'
		  AND l.location_latitude IS NOT NULL
		  AND l.location_longitude IS NOT NULL
		  AND l.location_latitude BETWEEN $2 - $4 AND $2 + $4
		  AND l.location_longitude BETWEEN $1 - $4 AND $1 + $4
		ORDER BY distance_meters ASC
		LIMIT 10
	`

	rows, err := s.db.GetSQLXDB().Query(query, lng, lat, radiusDegrees)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []map[string]interface{}
	for rows.Next() {
		var item struct {
			ID         int
			Title      string
			Price      float64
			Currency   string
			Latitude   float64
			Longitude  float64
			Distance   int
			StoreName  sql.NullString
			StoreLogo  sql.NullString
		}

		err := rows.Scan(
			&item.ID, &item.Title, &item.Price, &item.Currency,
			&item.Latitude, &item.Longitude, &item.Distance,
			&item.StoreName, &item.StoreLogo,
		)
		if err != nil {
			continue
		}

		itemMap := map[string]interface{}{
			"id":       item.ID,
			"title":    item.Title,
			"price":    item.Price,
			"currency": item.Currency,
			"location": map[string]float64{
				"latitude":  item.Latitude,
				"longitude": item.Longitude,
			},
			"distance_meters": item.Distance,
		}

		if item.StoreName.Valid {
			itemMap["store_name"] = item.StoreName.String
		}
		if item.StoreLogo.Valid {
			itemMap["store_logo"] = item.StoreLogo.String
		}

		items = append(items, itemMap)
	}

	return items, nil
}

// CalculateDistance вычисляет расстояние между двумя точками (формула Haversine)
func (s *DeliveryService) CalculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadius = 6371000 // метры

	phi1 := lat1 * math.Pi / 180
	phi2 := lat2 * math.Pi / 180
	deltaPhi := (lat2 - lat1) * math.Pi / 180
	deltaLambda := (lng2 - lng1) * math.Pi / 180

	a := math.Sin(deltaPhi/2)*math.Sin(deltaPhi/2) +
		math.Cos(phi1)*math.Cos(phi2)*
			math.Sin(deltaLambda/2)*math.Sin(deltaLambda/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// CalculateETA рассчитывает примерное время прибытия
func (s *DeliveryService) CalculateETA(delivery *Delivery) time.Time {
	if delivery.CourierLocation == nil {
		return delivery.EstimatedDeliveryTime
	}

	// Расстояние до точки доставки
	distance := s.CalculateDistance(
		delivery.CourierLocation.Latitude,
		delivery.CourierLocation.Longitude,
		delivery.DeliveryLatitude,
		delivery.DeliveryLongitude,
	)

	// Средняя скорость в м/с (если текущая скорость 0, используем среднюю)
	speed := delivery.CourierLocation.Speed
	if speed == 0 {
		speed = 10 // 36 км/ч средняя скорость курьера
	}

	// Время в секундах
	seconds := distance / speed
	return time.Now().Add(time.Duration(seconds) * time.Second)
}

// generateTrackingToken генерирует уникальный токен для трекинга
func (s *DeliveryService) generateTrackingToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// estimateDuration оценивает время доставки в минутах
func (s *DeliveryService) estimateDuration(distanceMeters float64) int {
	// Средняя скорость курьера 15 км/ч + время на получение/передачу
	minutes := (distanceMeters/1000)/15*60 + 10
	return int(minutes)
}

// getLastCourierLocation получает последнюю локацию курьера
func (s *DeliveryService) getLastCourierLocation(courierID int) *CourierLocation {
	var loc CourierLocation
	query := `
		SELECT current_latitude, current_longitude, last_location_update
		FROM couriers
		WHERE id = $1 AND current_latitude IS NOT NULL
	`

	err := s.db.GetSQLXDB().QueryRow(query, courierID).Scan(
		&loc.Latitude, &loc.Longitude, &loc.UpdatedAt,
	)

	if err != nil {
		return nil
	}

	return &loc
}

// GetActiveDeliveries получает активные доставки курьера
func (s *DeliveryService) GetActiveDeliveries(courierID int) ([]*Delivery, error) {
	query := `
		SELECT 
			d.id, d.order_id, d.tracking_token, d.status,
			d.pickup_address, d.delivery_address,
			d.pickup_latitude, d.pickup_longitude,
			d.delivery_latitude, d.delivery_longitude,
			d.estimated_delivery_time
		FROM deliveries d
		WHERE d.courier_id = $1 
		  AND d.status IN ('pending', 'picked_up', 'in_transit')
		ORDER BY d.created_at DESC
	`

	rows, err := s.db.GetSQLXDB().Query(query, courierID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deliveries []*Delivery
	for rows.Next() {
		var d Delivery
		err := rows.Scan(
			&d.ID, &d.OrderID, &d.TrackingToken, &d.Status,
			&d.PickupAddress, &d.DeliveryAddress,
			&d.PickupLatitude, &d.PickupLongitude,
			&d.DeliveryLatitude, &d.DeliveryLongitude,
			&d.EstimatedDeliveryTime,
		)
		if err != nil {
			continue
		}
		d.CourierID = courierID
		deliveries = append(deliveries, &d)
	}

	return deliveries, nil
}