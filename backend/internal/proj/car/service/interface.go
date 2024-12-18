package service

import (
    "context"
    "mime/multipart"
    "backend/internal/domain/models"
)

type CarServiceInterface interface {
    AddCar(ctx context.Context, car *models.Car) (int, error)
    GetAvailableCars(ctx context.Context, filters map[string]string) ([]models.Car, error)
    ProcessImage(file *multipart.FileHeader) (string, error)
    AddCarImage(ctx context.Context, image *models.CarImage) (int, error)
    GetCarImages(ctx context.Context, carID string) ([]models.CarImage, error)
    CreateBooking(ctx context.Context, booking *models.CarBooking) error
}
