// backend/internal/storage/postgres/car.go
package postgres

import (
    "backend/internal/domain/models"
    "context"
    "fmt"
	"log"
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
func (db *Database) GetAvailableCars(ctx context.Context) ([]models.Car, error) {
    query := `
        WITH car_features AS (
            SELECT cl.car_id, array_agg(f.name) as features
            FROM car_feature_links cl
            JOIN car_features f ON f.id = cl.feature_id
            GROUP BY cl.car_id
        )
        SELECT 
            c.id, c.make, c.model, c.year, c.price_per_day,
            c.location, c.latitude, c.longitude, c.description,
            c.availability, c.transmission, c.fuel_type, c.seats,
            c.daily_mileage_limit, c.insurance_included, c.created_at,
            COALESCE(cf.features, '{}'::text[]) as features
        FROM cars c
        LEFT JOIN car_features cf ON cf.car_id = c.id
        WHERE c.availability = true
        ORDER BY c.created_at DESC
    `

    rows, err := db.pool.Query(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var cars []models.Car
    for rows.Next() {
        var car models.Car
        err := rows.Scan(
            &car.ID, &car.Make, &car.Model, &car.Year, &car.PricePerDay,
            &car.Location, &car.Latitude, &car.Longitude, &car.Description,
            &car.Availability, &car.Transmission, &car.FuelType, &car.Seats,
            &car.DailyMileageLimit, &car.InsuranceIncluded, &car.CreatedAt,
            &car.Features,
        )
        if err != nil {
            continue
        }
        cars = append(cars, car)
    }

    // Загружаем изображения для каждой машины
    for i := range cars {
        images, err := db.GetCarImages(ctx, fmt.Sprintf("%d", cars[i].ID))
        if err == nil {
            cars[i].Images = images
        }
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