//backend/internal/storage/postgres/room.go
package postgres

import (
	"backend/internal/domain/models"
	"context"
	"fmt"
	"log"
	"strings"
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
func (db *Database) GetRooms(ctx context.Context, filters map[string]string, sortBy, sortDirection string, limit, offset int) ([]models.Room, int64, error) {
    // Базовый запрос с условной логикой для дат
    baseQuery := `WITH room_availability AS (
        SELECT 
            r.id,
            COALESCE(r.total_beds, 0) as total_beds,
            CASE 
                WHEN r.accommodation_type = 'bed' THEN 
                    CASE WHEN $1 = '' OR $2 = '' THEN
                        -- Если даты не указаны, проверяем только текущие бронирования
                        COALESCE(
                            r.total_beds - COALESCE((
                                SELECT COUNT(DISTINCT bb.bed_id)
                                FROM beds b2
                                LEFT JOIN bed_bookings bb ON b2.id = bb.bed_id
                                WHERE b2.room_id = r.id
                                AND bb.status = 'confirmed'
                                AND CURRENT_DATE BETWEEN bb.start_date AND bb.end_date
                            ), 0),
                            r.total_beds
                        )
                    ELSE
                        -- Если даты указаны, проверяем доступность на конкретные даты
                        COALESCE(
                            r.total_beds - COALESCE((
                                SELECT COUNT(DISTINCT bb.bed_id)
                                FROM beds b2
                                LEFT JOIN bed_bookings bb ON b2.id = bb.bed_id
                                WHERE b2.room_id = r.id
                                AND bb.status = 'confirmed'
                                AND $1::date <= bb.end_date 
                                AND $2::date >= bb.start_date
                            ), 0),
                            r.total_beds
                        )
                    END
                ELSE 
                    CASE WHEN $1 = '' OR $2 = '' THEN
                        -- Для комнат/апартаментов без дат
                        CASE WHEN EXISTS (
                            SELECT 1 FROM bookings b
                            WHERE b.room_id = r.id
                            AND b.status = 'confirmed'
                            AND CURRENT_DATE BETWEEN b.start_date AND b.end_date
                        ) THEN 0 ELSE 1 END
                    ELSE
                        -- Для комнат/апартаментов с датами
                        CASE WHEN EXISTS (
                            SELECT 1 FROM bookings b
                            WHERE b.room_id = r.id
                            AND b.status = 'confirmed'
                            AND $1::date <= b.end_date 
                            AND $2::date >= b.start_date
                        ) THEN 0 ELSE 1 END
                    END
            END as available_count,
            CASE 
                WHEN r.accommodation_type = 'bed' THEN 
                    COALESCE((
                        SELECT MIN(b2.price_per_night)
                        FROM beds b2
                        WHERE b2.room_id = r.id
                        AND b2.is_available = true
                        AND ($1 = '' OR $2 = '' OR NOT EXISTS (
                            SELECT 1
                            FROM bed_bookings bb
                            WHERE bb.bed_id = b2.id
                            AND bb.status = 'confirmed'
                            AND $1::date <= bb.end_date 
                            AND $2::date >= bb.start_date
                        ))
                    ), 0)
                ELSE r.price_per_night
            END as actual_price
        FROM rooms r
    )`

    // Формируем условия WHERE
	
    var conditions []string
    var args []interface{}
    args = append(args, filters["start_date"], filters["end_date"])
    paramCount := 2

    if filters["capacity"] != "" {
        paramCount++
        conditions = append(conditions, fmt.Sprintf("r.capacity >= $%d", paramCount))
        args = append(args, filters["capacity"])
    }

    if filters["min_price"] != "" {
        paramCount++
        conditions = append(conditions, fmt.Sprintf("ra.actual_price >= $%d", paramCount))
        args = append(args, filters["min_price"])
    }

    if filters["max_price"] != "" {
        paramCount++
        conditions = append(conditions, fmt.Sprintf("ra.actual_price <= $%d", paramCount))
        args = append(args, filters["max_price"])
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

    if filters["accommodation_type"] != "" {
        paramCount++
        conditions = append(conditions, fmt.Sprintf("r.accommodation_type = $%d", paramCount))
        args = append(args, filters["accommodation_type"])
    }

    // Формируем ORDER BY
    var orderBy string
    switch sortBy {
    case "price_per_night":
        orderBy = fmt.Sprintf("ra.actual_price %s", sortDirection)
    case "rating":
        orderBy = fmt.Sprintf("COALESCE(avg_rating.rating, 0) %s", sortDirection)
    default:
        orderBy = fmt.Sprintf("r.created_at %s", sortDirection)
    }

    // Собираем полный запрос для получения данных
	whereClause := []string{"ra.available_count > 0"}
	if len(conditions) > 0 {
		whereClause = append(whereClause, conditions...)
	}
	
	query := fmt.Sprintf(`
		%s
		SELECT 
			r.id, 
			r.name,
		r.capacity,
		r.price_per_night,
		r.address_street,
		r.address_city,
		r.address_state, 
		r.address_country,
		r.address_postal_code,
		r.accommodation_type,
		r.is_shared,
		r.total_beds,
		r.available_beds,
		r.has_private_bathroom,
		r.latitude,
		r.longitude,
		r.formatted_address,
		r.created_at,
		COALESCE(ra.actual_price, 0) as actual_price,
		COALESCE(ra.available_count, 0) as available_count,
		COALESCE(avg_rating.rating, 0) as rating,
		COUNT(*) OVER() as total_count
    FROM rooms r
    JOIN room_availability ra ON r.id = ra.id
    LEFT JOIN (
        SELECT room_id, AVG(rating) as rating 
        FROM room_reviews 
        GROUP BY room_id
    ) avg_rating ON r.id = avg_rating.room_id
    WHERE %s
    ORDER BY %s
    LIMIT $%d OFFSET $%d
`, baseQuery, strings.Join(whereClause, " AND "), orderBy, paramCount+1, paramCount+2)

    // Добавляем параметры пагинации
    args = append(args, limit, offset)

    // Выполняем запрос
    rows, err := db.pool.Query(ctx, query, args...)
    if err != nil {
        return nil, 0, fmt.Errorf("error querying rooms: %w", err)
    }
    defer rows.Close()

    var rooms []models.Room
    var total int64

	for rows.Next() {
		var room models.Room
		var rating float64
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
			&room.ActualPrice,
			&room.AvailableCount,
			&rating,
			&total,
		)
		if err != nil {
			log.Printf("Error scanning room: %v", err)
			continue
		}
		room.Rating = rating
		rooms = append(rooms, room)
	}
    return rooms, total, rows.Err()
}
