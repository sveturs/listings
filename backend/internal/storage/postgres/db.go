// backend/internal/storage/postgres/db.go
package postgres

import (
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/jackc/pgx/v5"
    "backend/internal/storage"
    "backend/internal/domain/models" 
    "fmt"
    //"log"
    "strconv"

    carStorage "backend/internal/proj/car/storage/postgres"

)

type Database struct {
    pool *pgxpool.Pool
    carDB *carStorage.Storage
}
func NewDatabase(dbURL string) (*Database, error) {
    pool, err := pgxpool.New(context.Background(), dbURL)
    if err != nil {
        return nil, fmt.Errorf("error creating connection pool: %w", err)
    }

    db := &Database{
        pool:  pool,
        carDB: carStorage.NewStorage(pool),
    }

    return db, nil
}
func (db *Database) AddBed(ctx context.Context, roomID int, bedNumber string, pricePerNight float64, hasOutlet bool, hasLight bool, hasShelf bool, bedType string) (int, error) {
    var bedID int
    err := db.pool.QueryRow(ctx, `
        INSERT INTO beds (room_id, bed_number, price_per_night, is_available, has_outlet, has_light, has_shelf, bed_type)
        VALUES ($1, $2, $3, true, $4, $5, $6, $7)
        RETURNING id`,
        roomID, bedNumber, pricePerNight, hasOutlet, hasLight, hasShelf, bedType).Scan(&bedID)

    return bedID, err
}

var _ storage.Storage = (*Database)(nil) 

func (db *Database) Close() {
    if db.pool != nil {
        db.pool.Close()
    }
}

func (db *Database) Ping(ctx context.Context) error {
    return db.pool.Ping(ctx)
}

type RowsWrapper struct {
    rows pgx.Rows
}

func (r *RowsWrapper) Next() bool {
    return r.rows.Next()
}

func (r *RowsWrapper) Scan(dest ...interface{}) error {
    return r.rows.Scan(dest...)
}

func (r *RowsWrapper) Close() error {
    r.rows.Close()
    return nil
}

func (db *Database) Query(ctx context.Context, sql string, args ...interface{}) (storage.Rows, error) {
    rows, err := db.pool.Query(ctx, sql, args...)
    if err != nil {
        return nil, err
    }
    return &RowsWrapper{rows: rows}, nil
}

func (db *Database) QueryRow(ctx context.Context, sql string, args ...interface{}) storage.Row {
    return db.pool.QueryRow(ctx, sql, args...)
}
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
// Car methods
func (db *Database) AddCar(ctx context.Context, car *models.Car) (int, error) {
    return db.carDB.AddCar(ctx, car)
}

func (db *Database) GetAvailableCars(ctx context.Context, filters map[string]string) ([]models.Car, error) {
    return db.carDB.GetAvailableCars(ctx, filters)
}

func (db *Database) GetCarWithFeatures(ctx context.Context, carID int) (*models.Car, error) {
    return db.carDB.GetCarWithFeatures(ctx, carID)
}

func (db *Database) GetCarFeatures(ctx context.Context) ([]models.CarFeature, error) {
    return db.carDB.GetCarFeatures(ctx)
}

func (db *Database) GetCarCategories(ctx context.Context) ([]models.CarCategory, error) {
    return db.carDB.GetCarCategories(ctx)
}

func (db *Database) CreateCarBooking(ctx context.Context, booking *models.CarBooking) error {
    return db.carDB.CreateCarBooking(ctx, booking)
}

// Car image methods
func (db *Database) AddCarImage(ctx context.Context, image *models.CarImage) (int, error) {
    return db.carDB.AddCarImage(ctx, image)
}

func (db *Database) GetCarImages(ctx context.Context, carID string) ([]models.CarImage, error) {
    return db.carDB.GetCarImages(ctx, carID)
}

func (db *Database) DeleteCarImage(ctx context.Context, imageID string) (string, error) {
    return db.carDB.DeleteCarImage(ctx, imageID)
}