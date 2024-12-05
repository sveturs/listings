package postgres

import (
	"backend/internal/domain/models"
	"context"
	"fmt"
	"strconv"
)

// AddBed добавляет новую кровать в комнату
func (db *Database) AddBed(ctx context.Context, roomID int, bedNumber string, pricePerNight float64, hasOutlet bool, hasLight bool, hasShelf bool, bedType string) (int, error) {
    var bedID int
    var roomExists bool
    err := db.pool.QueryRow(ctx,
        "SELECT EXISTS(SELECT 1 FROM rooms WHERE id = $1)",
        roomID).Scan(&roomExists)
    if err != nil || !roomExists {
        return 0, fmt.Errorf("room not found")
    }

    err = db.pool.QueryRow(ctx, `
        INSERT INTO beds (room_id, bed_number, price_per_night, is_available, has_outlet, has_light, has_shelf, bed_type)
        VALUES ($1, $2, $3, true, $4, $5, $6, $7)
        RETURNING id`,
        roomID, bedNumber, pricePerNight, hasOutlet, hasLight, hasShelf, bedType).Scan(&bedID)

    return bedID, err
}

func (db *Database) UpdateBedAttributes(ctx context.Context, bedID int, bedReq *models.BedRequest) error {
    _, err := db.pool.Exec(ctx, `
        UPDATE beds 
        SET has_outlet = $1, has_light = $2, has_shelf = $3, bed_type = $4
        WHERE id = $5`,
        bedReq.HasOutlet, bedReq.HasLight, bedReq.HasShelf, bedReq.BedType, bedID)
    return err
}
// GetAvailableBeds получает список доступных кроватей
func (db *Database) GetAvailableBeds(ctx context.Context, roomID string, startDate string, endDate string) ([]models.Bed, error) {
	query := `
    SELECT b.id, b.bed_number, b.price_per_night, b.has_outlet, b.has_light, b.has_shelf, b.bed_type
    FROM beds b
    WHERE b.room_id = $1
    AND b.is_available = true
    AND NOT EXISTS (
        SELECT 1
        FROM bed_bookings bb
        WHERE bb.bed_id = b.id
        AND bb.status = 'confirmed'
        AND (
            (bb.start_date <= $3 AND bb.end_date >= $2)
        )
    )
    ORDER BY b.bed_number
`

	rows, err := db.pool.Query(ctx, query, roomID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var beds []models.Bed
	for rows.Next() {
		var bed models.Bed
		if err := rows.Scan(
			&bed.ID, 
			&bed.BedNumber, 
			&bed.PricePerNight, 
			&bed.HasOutlet,
			&bed.HasLight,
			&bed.HasShelf,
			&bed.BedType,
		); err != nil {
			continue
		}
		bed.RoomID, _ = strconv.Atoi(roomID)
		bed.IsAvailable = true
		beds = append(beds, bed)
	}

	// Обновляем количество доступных кроватей в комнате
	_, err = db.pool.Exec(ctx, `
		UPDATE rooms 
		SET available_beds = $1
		WHERE id = $2
	`, len(beds), roomID)

	return beds, rows.Err()
}
