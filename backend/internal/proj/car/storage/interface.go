// internal/proj/car/storage/interface.go
package storage

import (
    "context"
    "backend/internal/domain/models"
)

type CarRepository interface {
    AddCar(ctx context.Context, car *models.Car) (int, error)
    GetAvailableCars(ctx context.Context, filters map[string]string) ([]models.Car, error)
    GetCarWithFeatures(ctx context.Context, carID int) (*models.Car, error)
    GetCarFeatures(ctx context.Context) ([]models.CarFeature, error)
    GetCarCategories(ctx context.Context) ([]models.CarCategory, error)
    CreateCarBooking(ctx context.Context, booking *models.CarBooking) error
}

type CarImageRepository interface {
    AddCarImage(ctx context.Context, image *models.CarImage) (int, error)
    GetCarImages(ctx context.Context, carID string) ([]models.CarImage, error)
    DeleteCarImage(ctx context.Context, imageID string) (string, error)
}

type Repository interface {
    CarRepository
    CarImageRepository
}