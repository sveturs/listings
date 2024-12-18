// backend/internal/storage/postgres/bed.go
package postgres

import (
    "backend/internal/domain/models"
    "context"
    "fmt"
)

func (db *Database) GetBedByID(ctx context.Context, id int) (*models.Bed, error) {
    bed := &models.Bed{}
    err := db.pool.QueryRow(ctx, `
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

func (db *Database) GetBedsByRoomID(ctx context.Context, roomID int) ([]models.Bed, error) {
    rows, err := db.pool.Query(ctx, `
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
        if err := rows.Scan(
            &bed.ID, &bed.RoomID, &bed.BedNumber, &bed.PricePerNight,
            &bed.IsAvailable, &bed.HasOutlet, &bed.HasLight, &bed.HasShelf,
            &bed.BedType); err != nil {
            continue
        }
        beds = append(beds, bed)
    }

    return beds, rows.Err()
}

func (db *Database) UpdateBedAttributes(ctx context.Context, bedID int, bedReq *models.BedRequest) error {
    _, err := db.pool.Exec(ctx, `
        UPDATE beds 
        SET has_outlet = $1, has_light = $2, has_shelf = $3, bed_type = $4
        WHERE id = $5`,
        bedReq.HasOutlet, bedReq.HasLight, bedReq.HasShelf, bedReq.BedType, bedID)
    return err
}

func (db *Database) UpdateBedAvailability(ctx context.Context, bedID int, isAvailable bool) error {
    _, err := db.pool.Exec(ctx, `
        UPDATE beds 
        SET is_available = $1 
        WHERE id = $2`,
        isAvailable, bedID)
    return err
}

func (db *Database) UpdateBedPrice(ctx context.Context, bedID int, price float64) error {
    _, err := db.pool.Exec(ctx, `
        UPDATE beds 
        SET price_per_night = $1 
        WHERE id = $2`,
        price, bedID)
    return err
}