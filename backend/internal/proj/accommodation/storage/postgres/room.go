package postgres

import (
    "context"
    "fmt"
    "log"
    "backend/internal/domain/models"
)

func (s *Storage) AddRoom(ctx context.Context, room *models.Room) (int, error) {
    var roomID int
    err := s.pool.QueryRow(ctx, `
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
        room.TotalBeds, room.AvailableBeds, room.HasPrivateBathroom,
        room.Latitude, room.Longitude, room.FormattedAddress,
    ).Scan(&roomID)

    if err != nil {
        return 0, fmt.Errorf("error adding room: %w", err)
    }
    return roomID, nil
}

func (s *Storage) GetRooms(ctx context.Context, filters map[string]string, sortBy string, sortDirection string, limit int, offset int) ([]models.Room, int64, error) {
    var conditions []string
    var args []interface{}
    paramCount := 1

    baseQuery := `
        SELECT 
            id, name, capacity, price_per_night,
            address_street, address_city, address_state,
            address_country, address_postal_code,
            accommodation_type, is_shared,
            total_beds, available_beds, has_private_bathroom,
            latitude, longitude, formatted_address,
            created_at,
            COUNT(*) OVER() as total_count
        FROM rooms
        WHERE 1=1
    `

    // Add filter conditions
    if v := filters["capacity"]; v != "" {
        conditions = append(conditions, fmt.Sprintf("capacity >= $%d", paramCount))
        args = append(args, v)
        paramCount++
    }

    // Add conditions to query
    if len(conditions) > 0 {
        baseQuery += " AND " + conditions[0]
    }

    // Add sorting
    baseQuery += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortDirection)

    // Add pagination
    baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramCount, paramCount+1)
    args = append(args, limit, offset)

    rows, err := s.pool.Query(ctx, baseQuery, args...)
    if err != nil {
        return nil, 0, err
    }
    defer rows.Close()

    var rooms []models.Room
    var totalCount int64

    for rows.Next() {
        var room models.Room
        err := rows.Scan(
            &room.ID, &room.Name, &room.Capacity, &room.PricePerNight,
            &room.AddressStreet, &room.AddressCity, &room.AddressState,
            &room.AddressCountry, &room.AddressPostalCode,
            &room.AccommodationType, &room.IsShared,
            &room.TotalBeds, &room.AvailableBeds, &room.HasPrivateBathroom,
            &room.Latitude, &room.Longitude, &room.FormattedAddress,
            &room.CreatedAt,
            &totalCount,
        )
        if err != nil {
            log.Printf("Error scanning room: %v", err)
            continue
        }
        rooms = append(rooms, room)
    }

    return rooms, totalCount, rows.Err()
}

func (s *Storage) GetRoomByID(ctx context.Context, id int) (*models.Room, error) {
    room := &models.Room{}
    err := s.pool.QueryRow(ctx, `
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
        return nil, fmt.Errorf("error getting room: %w", err)
    }
    return room, nil
}