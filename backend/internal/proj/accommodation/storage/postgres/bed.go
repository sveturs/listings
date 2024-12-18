package postgres

import (
    "context"
    "fmt"
    "backend/internal/domain/models"
)

func (s *Storage) AddBed(ctx context.Context, roomID int, bedNumber string, pricePerNight float64, hasOutlet bool, hasLight bool, hasShelf bool, bedType string) (int, error) {
    var bedID int
    err := s.pool.QueryRow(ctx, `
        INSERT INTO beds (
            room_id, bed_number, price_per_night, 
            is_available, has_outlet, has_light, 
            has_shelf, bed_type
        )
        VALUES ($1, $2, $3, true, $4, $5, $6, $7)
        RETURNING id`,
        roomID, bedNumber, pricePerNight, hasOutlet, hasLight, hasShelf, bedType).Scan(&bedID)

    return bedID, err
}

func (s *Storage) GetBedByID(ctx context.Context, id int) (*models.Bed, error) {
    bed := &models.Bed{}
    err := s.pool.QueryRow(ctx, `
        SELECT 
            id, room_id, bed_number, price_per_night, 
            is_available, has_outlet, has_light, has_shelf, 
            bed_type
        FROM beds 
        WHERE id = $1`,
        id).Scan(
            &bed.ID, &bed.RoomID, &bed.BedNumber, &bed.PricePerNight,
            &bed.IsAvailable, &bed.HasOutlet, &bed.HasLight, &bed.HasShelf,
            &bed.BedType)
    if err != nil {
        return nil, fmt.Errorf("failed to get bed: %w", err)
    }
    return bed, nil
}

func (s *Storage) GetBedsByRoomID(ctx context.Context, roomID int) ([]models.Bed, error) {
    rows, err := s.pool.Query(ctx, `
        SELECT 
            id, room_id, bed_number, price_per_night,
            is_available, has_outlet, has_light, has_shelf,
            bed_type
        FROM beds
        WHERE room_id = $1
        ORDER BY bed_number`,
        roomID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var beds []models.Bed
    for rows.Next() {
        var bed models.Bed
        err := rows.Scan(
            &bed.ID, &bed.RoomID, &bed.BedNumber, &bed.PricePerNight,
            &bed.IsAvailable, &bed.HasOutlet, &bed.HasLight, &bed.HasShelf,
            &bed.BedType)
        if err != nil {
            continue
        }
        beds = append(beds, bed)
    }
    return beds, rows.Err()
}

func (s *Storage) GetAvailableBeds(ctx context.Context, roomID string, startDate string, endDate string) ([]models.Bed, error) {
    query := `
        SELECT b.id, b.room_id, b.bed_number, b.price_per_night, 
               b.has_outlet, b.has_light, b.has_shelf, b.bed_type
        FROM beds b
        WHERE b.room_id = $1
        AND b.is_available = true
        AND NOT EXISTS (
            SELECT 1
            FROM bed_bookings bb
            WHERE bb.bed_id = b.id
            AND bb.status = 'confirmed'
            AND (bb.start_date <= $3 AND bb.end_date >= $2)
        )
        ORDER BY b.bed_number
    `

    rows, err := s.pool.Query(ctx, query, roomID, startDate, endDate)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var beds []models.Bed
    for rows.Next() {
        var bed models.Bed
        if err := rows.Scan(
            &bed.ID, 
            &bed.RoomID, 
            &bed.BedNumber, 
            &bed.PricePerNight,
            &bed.HasOutlet,
            &bed.HasLight,
            &bed.HasShelf,
            &bed.BedType,
        ); err != nil {
            continue
        }
        bed.IsAvailable = true
        beds = append(beds, bed)
    }

    return beds, rows.Err()
}

func (s *Storage) UpdateBedAttributes(ctx context.Context, bedID int, bedReq *models.BedRequest) error {
    _, err := s.pool.Exec(ctx, `
        UPDATE beds 
        SET has_outlet = $1, has_light = $2, has_shelf = $3, bed_type = $4
        WHERE id = $5`,
        bedReq.HasOutlet, bedReq.HasLight, bedReq.HasShelf, bedReq.BedType, bedID)
    return err
}

func (s *Storage) UpdateBedAvailability(ctx context.Context, bedID int, isAvailable bool) error {
    _, err := s.pool.Exec(ctx, `
        UPDATE beds 
        SET is_available = $1 
        WHERE id = $2`,
        isAvailable, bedID)
    return err
}

func (s *Storage) UpdateBedPrice(ctx context.Context, bedID int, price float64) error {
    _, err := s.pool.Exec(ctx, `
        UPDATE beds 
        SET price_per_night = $1 
        WHERE id = $2`,
        price, bedID)
    return err
}