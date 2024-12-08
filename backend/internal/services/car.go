package services

import (
    "backend/internal/domain/models"
    "backend/internal/storage"
    "context"
    "fmt"
    "mime/multipart"
    "os"
    "path/filepath"
    "time"
    "github.com/disintegration/imaging"
)

type CarService struct {
    storage storage.Storage
}

func NewCarService(storage storage.Storage) *CarService {
    return &CarService{storage: storage}
}

func (s *CarService) AddCar(ctx context.Context, car *models.Car) (int, error) {
    return s.storage.AddCar(ctx, car)
}

func (s *CarService) GetAvailableCars(ctx context.Context) ([]models.Car, error) {
    return s.storage.GetAvailableCars(ctx)
}

func (s *CarService) ProcessImage(file *multipart.FileHeader) (string, error) {
    src, err := file.Open()
    if err != nil {
        return "", err
    }
    defer src.Close()

    ext := filepath.Ext(file.Filename)
    fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
    filePath := filepath.Join("uploads", fileName)

    dst, err := os.Create(filePath)
    if err != nil {
        return "", err
    }
    defer dst.Close()

    img, err := imaging.Decode(src)
    if err != nil {
        return "", err
    }

    resized := imaging.Resize(img, 1200, 0, imaging.Lanczos)
    err = imaging.Save(resized, filePath)
    if err != nil {
        return "", err
    }

    return fileName, nil
}

func (s *CarService) AddCarImage(ctx context.Context, image *models.CarImage) (int, error) {
    return s.storage.AddCarImage(ctx, image)
}

func (s *CarService) GetCarImages(ctx context.Context, carID string) ([]models.CarImage, error) {
    return s.storage.GetCarImages(ctx, carID)
}