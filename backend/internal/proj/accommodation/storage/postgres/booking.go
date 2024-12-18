// backend/internal/proj/accommodation/storage/postgres/booking.go
package postgres

import (
	"backend/internal/domain/models"
	"context"
	"fmt"
)

func (s *Storage) CreateBooking(ctx context.Context, booking *models.BookingRequest) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Проверяем существование пользователя
	var userExists bool
	err = tx.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)",
		booking.UserID).Scan(&userExists)
	if err != nil || !userExists {

		return fmt.Errorf("user not found")
	}

	// Получаем тип комнаты
	var roomType string
	var isShared bool
	err = tx.QueryRow(ctx, `
			SELECT accommodation_type, is_shared 
			FROM rooms 
			WHERE id = $1`,
		booking.RoomID).Scan(&roomType, &isShared)
	if err != nil {
		return err
	}

	if booking.StartDate == booking.EndDate {
		return fmt.Errorf("check-out date must be after check-in date")
	}

	if roomType == "bed" {
		if booking.BedID == nil {
			return fmt.Errorf("bed ID is required for bed booking")
		}

		// Проверяем доступность кровати
		var isAvailable bool
		err = tx.QueryRow(ctx, `
				SELECT is_available 
				FROM beds 
				WHERE id = $1 AND room_id = $2`,
			*booking.BedID, booking.RoomID).Scan(&isAvailable)
		if err != nil || !isAvailable {
			return fmt.Errorf("bed is not available")
		}

		// Проверяем конфликты бронирования
		var conflictCount int
		err = tx.QueryRow(ctx, `
				SELECT COUNT(*) 
				FROM bed_bookings 
				WHERE bed_id = $1 
				AND status = 'confirmed'
				AND (
					(start_date <= $2 AND end_date >= $2) OR
					(start_date <= $3 AND end_date >= $3) OR
					(start_date >= $2 AND end_date <= $3)
				)`,
			*booking.BedID, booking.StartDate, booking.EndDate).Scan(&conflictCount)

		if err != nil || conflictCount > 0 {
			return fmt.Errorf("bed is already booked for these dates")
		}

		// Создаем бронирование кровати
		_, err = tx.Exec(ctx, `
				INSERT INTO bed_bookings (bed_id, user_id, start_date, end_date, status)
				VALUES ($1, $2, $3, $4, 'confirmed')`,
			*booking.BedID, booking.UserID, booking.StartDate, booking.EndDate)

		// Обновляем количество доступных кроватей
		_, err = tx.Exec(ctx, `
				UPDATE rooms 
				SET available_beds = available_beds - 1
				WHERE id = $1`,
			booking.RoomID)
	} else {
		// Проверяем доступность комнаты
		var count int
		err = tx.QueryRow(ctx, `
				SELECT COUNT(*) 
				FROM bookings 
				WHERE room_id = $1 
					AND status = 'confirmed'
					AND start_date <= $3 
					AND end_date >= $2`,
			booking.RoomID, booking.StartDate, booking.EndDate).Scan(&count)
		if err != nil || count > 0 {
			return fmt.Errorf("room is already booked for these dates")
		}

		// Создаем бронирование комнаты
		_, err = tx.Exec(ctx, `
				INSERT INTO bookings (user_id, room_id, start_date, end_date, status)
				VALUES ($1, $2, $3, $4, 'confirmed')`,
			booking.UserID, booking.RoomID, booking.StartDate, booking.EndDate)
	}

	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// GetAllBookings получает все бронирования
func (s *Storage) GetAllBookings(ctx context.Context) ([]models.Booking, error) {
	roomBookingsQuery := `
			SELECT b.id, b.user_id, b.room_id, NULL as bed_id, 
				   b.start_date, b.end_date, b.status,
				   r.name as room_name, r.accommodation_type,
				   u.name as user_name, u.email as user_email
			FROM bookings b
			JOIN rooms r ON b.room_id = r.id
			JOIN users u ON b.user_id = u.id
		`
	bedBookingsQuery := `
			SELECT bb.id, bb.user_id, b.room_id, bb.bed_id,
				   bb.start_date, bb.end_date, bb.status,
				   r.name as room_name, r.accommodation_type,
				   u.name as user_name, u.email as user_email
			FROM bed_bookings bb
			JOIN beds b ON bb.bed_id = b.id
			JOIN rooms r ON b.room_id = r.id
			JOIN users u ON bb.user_id = u.id
		`

	// Получаем бронирования комнат
	roomRows, err := s.pool.Query(ctx, roomBookingsQuery)
	if err != nil {
		return nil, err
	}
	defer roomRows.Close()

	// Получаем бронирования кроватей
	bedRows, err := s.pool.Query(ctx, bedBookingsQuery)
	if err != nil {
		return nil, err
	}
	defer bedRows.Close()

	var bookings []models.Booking

	// Обрабатываем бронирования комнат
	for roomRows.Next() {
		var booking models.Booking
		err := roomRows.Scan(
			&booking.ID, &booking.UserID, &booking.RoomID, &booking.BedID,
			&booking.StartDate, &booking.EndDate, &booking.Status,
			&booking.RoomName, &booking.AccommodationType,
			&booking.UserName, &booking.UserEmail,
		)
		if err != nil {
			continue
		}
		booking.BookingType = "room"
		bookings = append(bookings, booking)
	}

	// Обрабатываем бронирования кроватей
	for bedRows.Next() {
		var booking models.Booking
		err := bedRows.Scan(
			&booking.ID, &booking.UserID, &booking.RoomID, &booking.BedID,
			&booking.StartDate, &booking.EndDate, &booking.Status,
			&booking.RoomName, &booking.AccommodationType,
			&booking.UserName, &booking.UserEmail,
		)
		if err != nil {
			continue
		}
		booking.BookingType = "bed"
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

// DeleteBooking удаляет бронирование
func (s *Storage) DeleteBooking(ctx context.Context, bookingID string, bookingType string) error {
	var query string
	if bookingType == "room" {
		query = "DELETE FROM bookings WHERE id = $1"
	} else if bookingType == "bed" {
		query = "DELETE FROM bed_bookings WHERE id = $1"
	} else {
		return fmt.Errorf("invalid booking type")
	}

	_, err := s.pool.Exec(ctx, query, bookingID)
	return err
}
