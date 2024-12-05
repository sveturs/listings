package postgres

import (
	"backend/internal/domain/models"
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

// AddRoom добавляет новую комнату
func (db *Database) AddRoom(ctx context.Context, room *models.Room) (int, error) {
	var roomID int
	err := db.pool.QueryRow(ctx, `
		INSERT INTO rooms (
			name, capacity, price_per_night,
			address_street, address_city, address_state,
			address_country, address_postal_code,
			accommodation_type, is_shared,
			total_beds, available_beds, has_private_bathroom,
			latitude, longitude, formatted_address
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING id
	`,
		room.Name, room.Capacity, room.PricePerNight,
		room.AddressStreet, room.AddressCity, room.AddressState,
		room.AddressCountry, room.AddressPostalCode,
		room.AccommodationType, room.IsShared,
		room.TotalBeds, room.TotalBeds, // изначально available_beds = total_beds
		room.HasPrivateBathroom,
		room.Latitude, room.Longitude, room.FormattedAddress,
	).Scan(&roomID)

	return roomID, err
}
func (db *Database) GetRoomByID(ctx context.Context, id int) (*models.Room, error) {
	room := &models.Room{}
	err := db.pool.QueryRow(ctx, `
        SELECT 
            id, name, capacity, price_per_night,
            address_street, address_city, address_state,
            address_country, address_postal_code,
            accommodation_type, is_shared,
            total_beds, available_beds, has_private_bathroom,
            latitude, longitude, formatted_address,
            created_at
        FROM rooms 
        WHERE id = $1
    `, id).Scan(
		&room.ID, &room.Name, &room.Capacity, &room.PricePerNight,
		&room.AddressStreet, &room.AddressCity, &room.AddressState,
		&room.AddressCountry, &room.AddressPostalCode,
		&room.AccommodationType, &room.IsShared,
		&room.TotalBeds, &room.AvailableBeds, &room.HasPrivateBathroom,
		&room.Latitude, &room.Longitude, &room.FormattedAddress,
		&room.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get room: %w", err)
	}

	return room, nil
}

// GetRooms получает список комнат с фильтрацией
func (db *Database) GetRooms(ctx context.Context, filters map[string]string) ([]models.Room, error) {
	log.Printf("Database.GetRooms called with filters: %+v", filters)

	// Подготавливаем значения дат по умолчанию
	startDate := time.Now().Format("2006-01-02")
	if filters["start_date"] != "" {
		startDate = filters["start_date"]
	}
	endDate := time.Now().AddDate(0, 0, 1).Format("2006-01-02") // завтра по умолчанию
	if filters["end_date"] != "" {
		endDate = filters["end_date"]
	}

	// SQL запрос с CTE из предыдущего фрагмента
	baseQuery := `WITH room_availability AS (
    SELECT 
        r.id,
        COALESCE(r.total_beds, 0) as total_beds,
        CASE 
            WHEN r.accommodation_type = 'bed' THEN 
                COALESCE(
                    r.total_beds - COALESCE((
                        SELECT COUNT(DISTINCT bb.bed_id)
                        FROM beds b2
                        LEFT JOIN bed_bookings bb ON b2.id = bb.bed_id
                        WHERE b2.room_id = r.id
                        AND bb.status = 'confirmed'
                        AND bb.start_date <= $2
                        AND bb.end_date >= $1
                    ), 0),
                    r.total_beds
                )
            ELSE 
                CASE WHEN EXISTS (
                    SELECT 1 FROM bookings b
                    WHERE b.room_id = r.id
                    AND b.status = 'confirmed'
                    AND b.start_date <= $2
                    AND b.end_date >= $1
                ) THEN 0 ELSE 1 END
        END as available_count,
        CASE 
            WHEN r.accommodation_type = 'bed' THEN 
                COALESCE((
                    SELECT MIN(b2.price_per_night)
                    FROM beds b2
                    WHERE b2.room_id = r.id
                    AND b2.is_available = true
                    AND NOT EXISTS (
                        SELECT 1
                        FROM bed_bookings bb
                        WHERE bb.bed_id = b2.id
                        AND bb.status = 'confirmed'
                        AND bb.start_date <= $2
                        AND bb.end_date >= $1
                    )
                ), 0)
            ELSE r.price_per_night
        END as actual_price
    FROM rooms r
)
SELECT 
    r.id, 
    r.name, 
    r.capacity,
    ra.actual_price as price_per_night,
    r.address_street, 
    r.address_city, 
    r.address_state,
    r.address_country, 
    r.address_postal_code,
    r.accommodation_type, 
    r.is_shared,
    COALESCE(r.total_beds, 0) as total_beds,
    COALESCE(ra.available_count, 0) as available_beds,
    r.has_private_bathroom,
    COALESCE(r.latitude, 0) as latitude,
    COALESCE(r.longitude, 0) as longitude,
    r.formatted_address,
    r.created_at
FROM rooms r
JOIN room_availability ra ON r.id = ra.id
WHERE 1=1`

	var conditions []string
	args := []interface{}{startDate, endDate}
	paramCount := 2

	// Добавляем условия фильтрации
	if filters["capacity"] != "" {
		paramCount++
		conditions = append(conditions, fmt.Sprintf("r.capacity >= $%d", paramCount))
		args = append(args, filters["capacity"])
	}

	if filters["city"] != "" {
		paramCount++
		conditions = append(conditions, fmt.Sprintf("r.address_city ILIKE $%d", paramCount))
		args = append(args, "%"+filters["city"]+"%")
	}
	if filters["country"] != "" {
		paramCount++
		conditions = append(conditions, fmt.Sprintf("r.address_country ILIKE $%d", paramCount))
		args = append(args, "%"+filters["country"]+"%")
	}

	if filters["min_price"] != "" {
		paramCount++
		conditions = append(conditions, fmt.Sprintf("r.price_per_night >= $%d", paramCount))
		args = append(args, filters["min_price"])
	}

	if filters["max_price"] != "" {
		paramCount++
		conditions = append(conditions, fmt.Sprintf("r.price_per_night <= $%d", paramCount))
		args = append(args, filters["max_price"])
	}

	if filters["accommodation_type"] != "" {
		paramCount++
		conditions = append(conditions, fmt.Sprintf("r.accommodation_type = $%d", paramCount))
		args = append(args, filters["accommodation_type"])
	}

	if filters["has_private_rooms"] == "true" {
		conditions = append(conditions, "r.is_shared = false")
	}

	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	log.Printf("Executing query: %s with params: %v", baseQuery, args)
	rows, err := db.pool.Query(ctx, baseQuery, args...)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.ID,
			&room.Name,
			&room.Capacity,
			&room.PricePerNight,
			&room.AddressStreet,
			&room.AddressCity,
			&room.AddressState,
			&room.AddressCountry,
			&room.AddressPostalCode,
			&room.AccommodationType,
			&room.IsShared,
			&room.TotalBeds,
			&room.AvailableBeds,
			&room.HasPrivateBathroom,
			&room.Latitude,
			&room.Longitude,
			&room.FormattedAddress,
			&room.CreatedAt,
		)
		if err != nil {
			log.Printf("Error scanning room: %v", err)
			continue
		}
		rooms = append(rooms, room)
	}

	log.Printf("Found %d rooms in database", len(rooms))
	return rooms, nil
}
