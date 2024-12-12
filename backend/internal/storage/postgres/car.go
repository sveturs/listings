// backend/internal/storage/postgres/car.go

package postgres

import (
    "backend/internal/domain/models"
    "context"
    "fmt"
	"log"
    "strings"
)

func (db *Database) AddCar(ctx context.Context, car *models.Car) (int, error) {
    log.Printf("Starting AddCar with data: %+v", car)

    tx, err := db.pool.Begin(ctx)
    if err != nil {
        return 0, fmt.Errorf("error beginning transaction: %w", err)
    }
    defer tx.Rollback(ctx)

    var carID int
    // Вставляем основные данные автомобиля без features
    err = tx.QueryRow(ctx, `
        INSERT INTO cars (
            make, model, year, price_per_day, location, 
            latitude, longitude, description, availability,
            transmission, fuel_type, seats,
            daily_mileage_limit, insurance_included
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
        RETURNING id
    `,
        car.Make, car.Model, car.Year, car.PricePerDay, car.Location,
        car.Latitude, car.Longitude, car.Description, car.Availability,
        car.Transmission, car.FuelType, car.Seats,
        car.DailyMileageLimit, car.InsuranceIncluded,
    ).Scan(&carID)

    if err != nil {
        log.Printf("Error executing insert car query: %v", err)
        return 0, fmt.Errorf("error inserting car: %w", err)
    }

    // Добавляем features через связующую таблицу
    if len(car.Features) > 0 {
        // For each feature
for _, featureName := range car.Features {
    var featureID int
    // First try to get existing feature
    err = tx.QueryRow(ctx, `
        SELECT id FROM car_features WHERE name = $1
    `, featureName).Scan(&featureID)

    if err != nil {
        // If feature doesn't exist, create it
        err = tx.QueryRow(ctx, `
            INSERT INTO car_features (name)
            VALUES ($1)
            RETURNING id
        `, featureName).Scan(&featureID)
        
        if err != nil {
            log.Printf("Error creating feature %s: %v", featureName, err)
            continue
        }
    }

    // Create link between car and feature
    _, err = tx.Exec(ctx, `
        INSERT INTO car_feature_links (car_id, feature_id)
        VALUES ($1, $2)
        ON CONFLICT DO NOTHING
    `, carID, featureID)

    if err != nil {
        log.Printf("Error linking feature %s to car: %v", featureName, err)
    }
}
    }

    err = tx.Commit(ctx)
    if err != nil {
        log.Printf("Error committing transaction: %v", err)
        return 0, fmt.Errorf("error committing transaction: %w", err)
    }

    log.Printf("Successfully added car with ID: %d", carID)
    return carID, nil
}
func (db *Database) CreateCarBooking(ctx context.Context, booking *models.CarBooking) error {
    tx, err := db.pool.Begin(ctx)
    if err != nil {
        return err
    }
    defer tx.Rollback(ctx)

    // Проверяем доступность автомобиля на выбранные даты
    var count int
    err = tx.QueryRow(ctx, `
        SELECT COUNT(*)
        FROM car_bookings
        WHERE car_id = $1
        AND status = 'confirmed'
        AND (
            ($2 BETWEEN start_date AND end_date)
            OR ($3 BETWEEN start_date AND end_date)
            OR (start_date BETWEEN $2 AND $3)
        )
    `, booking.CarID, booking.StartDate, booking.EndDate).Scan(&count)

    if err != nil {
        return err
    }

    if count > 0 {
        return fmt.Errorf("car is not available for selected dates")
    }

    // Создаем бронирование
    err = tx.QueryRow(ctx, `
        INSERT INTO car_bookings (
            car_id, user_id, start_date, end_date,
            pickup_location, dropoff_location,
            status, total_price
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id
    `,
        booking.CarID,
        booking.UserID,
        booking.StartDate,
        booking.EndDate,
        booking.PickupLocation,
        booking.DropoffLocation,
        "pending",
        booking.TotalPrice,
    ).Scan(&booking.ID)

    if err != nil {
        return err
    }

    return tx.Commit(ctx)
}
func (db *Database) GetCarWithFeatures(ctx context.Context, carID int) (*models.Car, error) {
    var car models.Car
    err := db.pool.QueryRow(ctx, `
        SELECT 
            c.id, c.make, c.model, c.year, c.price_per_day,
            c.location, c.availability, c.latitude, c.longitude,
            c.description, c.seats, c.transmission, c.fuel_type,
            c.daily_mileage_limit, c.insurance_included, c.created_at,
            array_remove(array_agg(f.name), NULL) as features
        FROM cars c
        LEFT JOIN car_feature_links cfl ON c.id = cfl.car_id
        LEFT JOIN car_features f ON cfl.feature_id = f.id
        WHERE c.id = $1
        GROUP BY c.id
    `, carID).Scan(
        &car.ID,
        &car.Make,
        &car.Model,
        &car.Year,
        &car.PricePerDay,
        &car.Location,
        &car.Availability,
        &car.Latitude,
        &car.Longitude,
        &car.Description,
        &car.Seats,
        &car.Transmission,
        &car.FuelType,
        &car.DailyMileageLimit,
        &car.InsuranceIncluded,
        &car.CreatedAt,
        &car.Features,
    )
    return &car, err
}
func (db *Database) GetAvailableCars(ctx context.Context, filters map[string]string) ([]models.Car, error) {
    var count int
    err := db.pool.QueryRow(ctx, "SELECT COUNT(*) FROM cars WHERE availability = true").Scan(&count)
    if err != nil {
        log.Printf("Error counting cars: %v", err)
    } else {
        log.Printf("Total available cars in DB: %d", count)
    }

    baseQuery := `
        SELECT 
            c.id, 
            c.make, 
            c.model, 
            c.year, 
            c.price_per_day,
            c.location, 
            c.latitude, 
            c.longitude, 
            c.description,
            c.availability, 
            c.transmission, 
            c.fuel_type, 
            c.seats,
            c.daily_mileage_limit, 
            c.insurance_included, 
            c.created_at,
            ARRAY_AGG(DISTINCT f.name) as features
        FROM cars c
        LEFT JOIN car_feature_links cfl ON c.id = cfl.car_id
        LEFT JOIN car_features f ON cfl.feature_id = f.id
        WHERE c.availability = true
    `

    // Добавляем условия фильтрации
    var conditions []string
    var args []interface{}
    argCount := 1

    if v := filters["make"]; v != "" {
        conditions = append(conditions, fmt.Sprintf("LOWER(c.make) LIKE LOWER($%d)", argCount))
        args = append(args, "%"+v+"%")
        argCount++
    }

    if v := filters["transmission"]; v != "" {
        conditions = append(conditions, fmt.Sprintf("c.transmission = $%d", argCount))
        args = append(args, v)
        argCount++
    }

    if v := filters["fuel_type"]; v != "" {
        conditions = append(conditions, fmt.Sprintf("c.fuel_type = $%d", argCount))
        args = append(args, v)
        argCount++
    }

    // Добавляем условия в запрос
    if len(conditions) > 0 {
        baseQuery += " AND " + strings.Join(conditions, " AND ")
    }

    // Добавляем группировку
    baseQuery += " GROUP BY c.id"

    // Добавляем сортировку
    baseQuery += " ORDER BY c.created_at DESC"

    log.Printf("Executing cars query...")

    rows, err := db.pool.Query(ctx, baseQuery, args...)
    if err != nil {
        log.Printf("Error executing query: %v", err)
        return nil, err
    }
    defer rows.Close()

    var cars []models.Car
    for rows.Next() {
        var car models.Car
        err := rows.Scan(
            &car.ID, 
            &car.Make, 
            &car.Model, 
            &car.Year,
            &car.PricePerDay,
            &car.Location, 
            &car.Latitude,
            &car.Longitude,
            &car.Description,
            &car.Availability,
            &car.Transmission,
            &car.FuelType,
            &car.Seats,
            &car.DailyMileageLimit,
            &car.InsuranceIncluded,
            &car.CreatedAt,
            &car.Features,
        )
        if err != nil {
            log.Printf("Error scanning row: %v", err)
            continue
        }
        log.Printf("Found car: ID=%d, Make=%s, Model=%s", car.ID, car.Make, car.Model)
        cars = append(cars, car)
    }

    log.Printf("Found %d cars after scanning", len(cars))

    // Загружаем изображения для каждой машины
    for i := range cars {
        images, err := db.GetCarImages(ctx, fmt.Sprintf("%d", cars[i].ID))
        if err != nil {
            log.Printf("Error loading images for car %d: %v", cars[i].ID, err)
            continue
        }
        cars[i].Images = images
        log.Printf("Loaded %d images for car %d", len(images), cars[i].ID)
    }

    return cars, nil
}

func (db *Database) GetCarFeatures(ctx context.Context) ([]models.CarFeature, error) {
    rows, err := db.pool.Query(ctx, `
        SELECT id, name, category, description 
        FROM car_features 
        ORDER BY category, name
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var features []models.CarFeature
    for rows.Next() {
        var feature models.CarFeature
        if err := rows.Scan(&feature.ID, &feature.Name, &feature.Category, &feature.Description); err != nil {
            return nil, err
        }
        features = append(features, feature)
    }
    return features, nil
}


func (db *Database) GetCarCategories(ctx context.Context) ([]models.CarCategory, error) {
    rows, err := db.pool.Query(ctx, `
        SELECT id, name, description
        FROM car_categories
        ORDER BY name
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var categories []models.CarCategory
    for rows.Next() {
        var category models.CarCategory
        err := rows.Scan(
            &category.ID,
            &category.Name,
            &category.Description,
        )
        if err != nil {
            continue
        }
        categories = append(categories, category)
    }

    return categories, nil
}