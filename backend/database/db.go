package database

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
	"strconv"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	pool *pgxpool.Pool
}

type Room struct {
	Name               string   `json:"name"`
	Capacity           int      `json:"capacity"`
	PricePerNight     *float64 `json:"price_per_night"`
	AddressStreet     string   `json:"address_street"`
	AddressCity       string   `json:"address_city"`
	AddressState      string   `json:"address_state"`
	AddressCountry    string   `json:"address_country"`
	AddressPostalCode string   `json:"address_postal_code"`
	AccommodationType string   `json:"accommodation_type"`
	IsShared          bool     `json:"is_shared"`
	TotalBeds         *int     `json:"total_beds"`
	AvailableBeds     *int     `json:"available_beds"`
	HasPrivateBathroom bool    `json:"has_private_bathroom"`
	Latitude          *float64 `json:"latitude"`
	Longitude         *float64 `json:"longitude"`
	FormattedAddress  string   `json:"formatted_address"`
}

type RoomImage struct {
	ID          int       `json:"id"`
	RoomID      int       `json:"room_id"`
	FilePath    string    `json:"file_path"`
	FileName    string    `json:"file_name"`
	FileSize    int       `json:"file_size"`
	ContentType string    `json:"content_type"`
	IsMain      bool      `json:"is_main"`
	CreatedAt   time.Time `json:"created_at"`
}

type Bed struct {
	ID            int     `json:"id"`
	RoomID        int     `json:"room_id"`
	BedNumber     string  `json:"bed_number"`
	IsAvailable   bool    `json:"is_available"`
	PricePerNight float64 `json:"price_per_night"`
}

type BedImage struct {
	ID          int       `json:"id"`
	BedID       int       `json:"bed_id"`
	FilePath    string    `json:"file_path"`
	FileName    string    `json:"file_name"`
	FileSize    int       `json:"file_size"`
	ContentType string    `json:"content_type"`
	CreatedAt   time.Time `json:"created_at"`
}

type BookingRequest struct {
	UserID    int    `json:"user_id"`
	RoomID    int    `json:"room_id"`
	BedID     *int   `json:"bed_id,omitempty"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type Booking struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	RoomID         int       `json:"room_id"`
	BedID          *int      `json:"bed_id"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	Status         string    `json:"status"`
	RoomName       string    `json:"room_name"`
	AccommodationType string `json:"type"`
	UserName       string    `json:"user_name"`
	UserEmail      string    `json:"user_email"`
	BookingType    string    `json:"booking_type"`
}

func New(dbURL string) (*Database, error) {
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, err
	}
	return &Database{pool: pool}, nil
}

func (db *Database) Close() {
	db.pool.Close()
}

// User methods
func (db *Database) AddUser(ctx context.Context, name, email string) error {
	_, err := db.pool.Exec(ctx, "INSERT INTO users (name, email) VALUES ($1, $2)", name, email)
	return err
}

// Room methods
func (db *Database) AddRoom(ctx context.Context, room Room) (int, error) {
	var roomID int
	totalBeds := 0
	availableBeds := 0
	
	if room.TotalBeds != nil {
		totalBeds = *room.TotalBeds
		if room.AvailableBeds != nil {
			availableBeds = *room.AvailableBeds
		} else {
			availableBeds = totalBeds
		}
	}

	err := db.pool.QueryRow(ctx, `
		INSERT INTO rooms (
			name, capacity, price_per_night,
			address_street, address_city, address_state,
			address_country, address_postal_code,
			accommodation_type, is_shared,
			total_beds, available_beds, has_private_bathroom,
			latitude, longitude, formatted_address
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16)
		RETURNING id
	`, room.Name, room.Capacity, room.PricePerNight,
		room.AddressStreet, room.AddressCity, room.AddressState,
		room.AddressCountry, room.AddressPostalCode,
		room.AccommodationType, room.IsShared,
		totalBeds, availableBeds, room.HasPrivateBathroom,
		room.Latitude, room.Longitude, room.FormattedAddress,
	).Scan(&roomID)

	return roomID, err
}

func (db *Database) DeleteRoom(ctx context.Context, id string) error {
	result, err := db.pool.Exec(ctx, "DELETE FROM rooms WHERE id=$1", id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("room not found")
	}
	return nil
}

func (db *Database) GetRooms(ctx context.Context, filters map[string]string) ([]map[string]interface{}, error) {
	baseQuery := `
	WITH room_availability AS (
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
					COALESCE(
						(SELECT MIN(b3.price_per_night) 
						 FROM beds b3 
						 WHERE b3.room_id = r.id 
						 AND b3.is_available = true
						 AND b3.id NOT IN (
							SELECT bb.bed_id
							FROM bed_bookings bb
							WHERE bb.status = 'confirmed'
							AND bb.start_date <= $2
							AND bb.end_date >= $1
						 )),
						r.price_per_night
					)
				ELSE r.price_per_night
			END as actual_price
		FROM rooms r
	)
	SELECT 
		r.id, 
		r.name, 
		r.capacity,
		r.latitude,
		r.longitude,
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
		r.created_at
	FROM rooms r
	JOIN room_availability ra ON r.id = ra.id
	WHERE 1=1
		AND (
			CASE 
				WHEN r.accommodation_type = 'bed' 
				THEN COALESCE(ra.available_count, 0) > 0 
				ELSE COALESCE(ra.available_count, 1) = 1
			END
		)`

	var conditions []string
	args := []interface{}{filters["start_date"], filters["end_date"]}
	paramCount := 2

	if accommodationType := filters["accommodation_type"]; accommodationType != "" {
		paramCount++
		conditions = append(conditions, fmt.Sprintf("r.accommodation_type = $%d", paramCount))
		args = append(args, accommodationType)

		if accommodationType == "room" && filters["has_private_rooms"] == "true" {
			conditions = append(conditions, "r.is_shared = false")
		}
	}

	if capacity := filters["capacity"]; capacity != "" {
		paramCount++
		conditions = append(conditions, fmt.Sprintf("r.capacity >= $%d", paramCount))
		args = append(args, capacity)
	}

	if minPrice := filters["min_price"]; minPrice != "" {
		paramCount++
		conditions = append(conditions, fmt.Sprintf("ra.actual_price >= $%d", paramCount))
		args = append(args, minPrice)
	}

	if maxPrice := filters["max_price"]; maxPrice != "" {
		paramCount++
		conditions = append(conditions, fmt.Sprintf("ra.actual_price <= $%d", paramCount))
		args = append(args, maxPrice)
	}

	if city := filters["city"]; city != "" {
		paramCount++
		conditions = append(conditions, fmt.Sprintf("r.address_city ILIKE $%d", paramCount))
		args = append(args, "%"+city+"%")
	}

	if country := filters["country"]; country != "" {
		paramCount++
		conditions = append(conditions, fmt.Sprintf("r.address_country ILIKE $%d", paramCount))
		args = append(args, "%"+country+"%")
	}

	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	baseQuery += " ORDER BY ra.available_count DESC, ra.actual_price ASC, r.created_at DESC"

	rows, err := db.pool.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []map[string]interface{}
	for rows.Next() {
		var id, capacity int
		var totalBeds, availableBeds int
		var name, addressStreet, addressCity, addressState,
			addressCountry, addressPostalCode, accommodationType string
		var latitude, longitude float64
		var pricePerNight float64
		var isShared, hasPrivateBathroom bool
		var createdAt time.Time

		if err := rows.Scan(
			&id, &name, &capacity, &latitude, &longitude, &pricePerNight,
			&addressStreet, &addressCity, &addressState,
			&addressCountry, &addressPostalCode,
			&accommodationType, &isShared,
			&totalBeds, &availableBeds,
			&hasPrivateBathroom,
			&createdAt,
		); err != nil {
			log.Printf("Ошибка сканирования строки: %v", err)
			continue
		}

		rooms = append(rooms, map[string]interface{}{
			"id":                   id,
			"name":                 name,
			"capacity":             capacity,
			"latitude":             latitude,
			"longitude":            longitude,
			"price_per_night":      pricePerNight,
			"address_street":       addressStreet,
			"address_city":         addressCity,
			"address_state":        addressState,
			"address_country":      addressCountry,
			"address_postal_code":  addressPostalCode,
			"accommodation_type":   accommodationType,
			"is_shared":           isShared,
			"total_beds":          totalBeds,
			"available_beds":      availableBeds,
			"has_private_bathroom": hasPrivateBathroom,
			"created_at":          createdAt.Format("2006-01-02 15:04:05"),
		})
	}

	return rooms, nil
}

// Room Image methods
func (db *Database) AddRoomImage(ctx context.Context, roomID int, image RoomImage) (int, error) {
	var imageID int
	err := db.pool.QueryRow(ctx, `
		INSERT INTO room_images (room_id, file_path, file_name, file_size, content_type, is_main)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, roomID, image.FilePath, image.FileName, image.FileSize, image.ContentType, image.IsMain).Scan(&imageID)
	return imageID, err
}

func (db *Database) GetRoomImages(ctx context.Context, roomID string) ([]RoomImage, error) {
	rows, err := db.pool.Query(ctx, `
		SELECT id, room_id, file_path, file_name, file_size, content_type, is_main, created_at
		FROM room_images
		WHERE room_id = $1
		ORDER BY is_main DESC, created_at DESC
	`, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []RoomImage
	for rows.Next() {
		var img RoomImage
		err := rows.Scan(
			&img.ID,
			&img.RoomID,
			&img.FilePath,
			&img.FileName,
			&img.FileSize,
			&img.ContentType,
			&img.IsMain,
			&img.CreatedAt,
		)
		if err != nil {
			log.Printf("Ошибка сканирования изображения: %v", err)
			continue
		}
		images = append(images, img)
	}
	return images, rows.Err()
}

func (db *Database) DeleteRoomImage(ctx context.Context, imageID string) (string, error) {
	var filePath string
	err := db.pool.QueryRow(ctx, "SELECT file_path FROM room_images WHERE id = $1", imageID).Scan(&filePath)
	if err != nil {
		return "", err
	}

	_, err = db.pool.Exec(ctx, "DELETE FROM room_images WHERE id = $1", imageID)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

// Bed methods
func (db *Database) AddBed(ctx context.Context, roomID int, bedNumber string, pricePerNight float64) (int, error) {
	var bedID int
	// Проверяем, существует ли комната
	var roomExists bool
	err := db.pool.QueryRow(ctx, 
		"SELECT EXISTS(SELECT 1 FROM rooms WHERE id = $1)", 
		roomID).Scan(&roomExists)
	if err != nil || !roomExists {
		return 0, fmt.Errorf("room not found")
	}

	// Добавляем кровать
	err = db.pool.QueryRow(ctx, `
		INSERT INTO beds (room_id, bed_number, price_per_night, is_available) 
		VALUES ($1, $2, $3, true)
		RETURNING id`,
		roomID, bedNumber, pricePerNight).Scan(&bedID)

	return bedID, err
}

func (db *Database) GetAvailableBeds(ctx context.Context, roomID string, startDate, endDate string) ([]Bed, error) {
	query := `
		SELECT b.id, b.bed_number, b.price_per_night
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

	var beds []Bed
	for rows.Next() {
		var bed Bed
		if err := rows.Scan(&bed.ID, &bed.BedNumber, &bed.PricePerNight); err != nil {
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

// Bed Image methods
func (db *Database) AddBedImage(ctx context.Context, bedID int, image BedImage) (int, error) {
	var imageID int
	err := db.pool.QueryRow(ctx, `
		INSERT INTO bed_images (bed_id, file_path, file_name, file_size, content_type)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, bedID, image.FilePath, image.FileName, image.FileSize, image.ContentType).Scan(&imageID)
	return imageID, err
}

func (db *Database) GetBedImages(ctx context.Context, bedID string) ([]BedImage, error) {
	rows, err := db.pool.Query(ctx, `
		SELECT id, bed_id, file_path, file_name, file_size, content_type, created_at
		FROM bed_images
		WHERE bed_id = $1
		ORDER BY created_at DESC
	`, bedID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []BedImage
	for rows.Next() {
		var img BedImage
		err := rows.Scan(
			&img.ID,
			&img.BedID,
			&img.FilePath,
			&img.FileName,
			&img.FileSize,
			&img.ContentType,
			&img.CreatedAt,
		)
		if err != nil {
			continue
		}
		images = append(images, img)
	}
	return images, rows.Err()
}

// Booking methods
func (db *Database) CreateBooking(ctx context.Context, booking BookingRequest) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var userExists bool
	err = tx.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)",
		booking.UserID).Scan(&userExists)
	if err != nil || !userExists {
		return fmt.Errorf("user not found")
	}

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

		var isAvailable bool
		err = tx.QueryRow(ctx, `
			SELECT is_available 
			FROM beds 
			WHERE id = $1 AND room_id = $2`,
			*booking.BedID, booking.RoomID).Scan(&isAvailable)
		if err != nil || !isAvailable {
			return fmt.Errorf("bed is not available")
		}

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

		_, err = tx.Exec(ctx, `
			INSERT INTO bed_bookings (bed_id, user_id, start_date, end_date, status)
			VALUES ($1, $2, $3, $4, 'confirmed')`,
			*booking.BedID, booking.UserID, booking.StartDate, booking.EndDate)

		_, err = tx.Exec(ctx, `
			UPDATE rooms 
			SET available_beds = available_beds - 1
			WHERE id = $1`,
			booking.RoomID)
	} else {
		var count int
		err = tx.QueryRow(ctx, `
			SELECT COUNT(*) 
			FROM bookings 
			WHERE room_id = $1 
				AND start_date <= $3 
				AND end_date >= $2
				AND status = 'confirmed'`,
			booking.RoomID, booking.StartDate, booking.EndDate).Scan(&count)
		if err != nil || count > 0 {
			return fmt.Errorf("room is already booked for these dates")
		}

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

func (db *Database) GetAllBookings(ctx context.Context) ([]Booking, error) {
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

	roomRows, err := db.pool.Query(ctx, roomBookingsQuery)
	if err != nil {
		return nil, err
	}
	defer roomRows.Close()

	bedRows, err := db.pool.Query(ctx, bedBookingsQuery)
	if err != nil {
		return nil, err
	}
	defer bedRows.Close()

	var bookings []Booking

	// Process room bookings
	for roomRows.Next() {
		var booking Booking
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

	// Process bed bookings
	for bedRows.Next() {
		var booking Booking
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

func (db *Database) DeleteBooking(ctx context.Context, bookingID, bookingType string) error {
	var query string
	if bookingType == "room" {
		query = "DELETE FROM bookings WHERE id = $1"
	} else if bookingType == "bed" {
		query = "DELETE FROM bed_bookings WHERE id = $1"
	} else {
		return fmt.Errorf("invalid booking type")
	}

	_, err := db.pool.Exec(ctx, query, bookingID)
	return err
}