package postgres

import (
	"backend/internal/domain/models"
	"context"

)

func (db *Database) AddCar(ctx context.Context, car *models.Car) (int, error) {
    var carID int
    err := db.pool.QueryRow(ctx, `
        INSERT INTO cars (make, model, year, price_per_day, availability, location)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `, car.Make, car.Model, car.Year, car.PricePerDay, car.Availability, car.Location).Scan(&carID)
    return carID, err
}

func (db *Database) GetAvailableCars(ctx context.Context) ([]models.Car, error) {
    rows, err := db.pool.Query(ctx, `
        SELECT id, make, model, year, price_per_day, availability, location
        FROM cars WHERE availability = TRUE
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var cars []models.Car
    for rows.Next() {
        var car models.Car
        if err := rows.Scan(&car.ID, &car.Make, &car.Model, &car.Year, &car.PricePerDay, &car.Availability, &car.Location); err != nil {
            return nil, err
        }
        cars = append(cars, car)
    }
    return cars, nil
}
