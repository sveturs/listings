// backend/internal/services/car.go
package service

import (
	"backend/internal/domain/models"
	"backend/internal/storage"
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"
)

type CarService struct {
    storage storage.Storage
}

func NewCarService(storage storage.Storage) CarServiceInterface {
    return &CarService{
        storage: storage,
    }
}

func (s *CarService) AddCar(ctx context.Context, car *models.Car) (int, error) {
	return s.storage.AddCar(ctx, car)
}

func (s *CarService) GetAvailableCars(ctx context.Context, filters map[string]string) ([]models.Car, error) {
    return s.storage.GetAvailableCars(ctx, filters)
}

const (
	maxWidth    = 1920 // Максимальная ширина
	maxHeight   = 1080 // Максимальная высота
	jpegQuality = 85   // Качество JPEG
)

func (s *CarService) ProcessImage(file *multipart.FileHeader) (string, error) {
	// Открываем загруженный файл
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("error opening uploaded file: %w", err)
	}
	defer src.Close()

	// Читаем содержимое в буфер
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, src); err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}

	// Декодируем изображение
	img, format, err := image.Decode(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return "", fmt.Errorf("error decoding image: %w", err)
	}

	log.Printf("Processing image: original format = %s, size = %d bytes", format, buf.Len())

	// Создаем имя файла
	fileName := fmt.Sprintf("%d.jpg", time.Now().UnixNano())
	filePath := filepath.Join("uploads", fileName)

	// Масштабируем изображение
	maxWidth := 1920
	maxHeight := 1080
	currentWidth := img.Bounds().Dx()
	currentHeight := img.Bounds().Dy()

	var resized image.Image
	if currentWidth > maxWidth || currentHeight > maxHeight {
		resized = imaging.Fit(img, maxWidth, maxHeight, imaging.Lanczos)
	} else {
		resized = img
	}

	// Оптимизируем и сохраняем как JPEG
	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("error creating output file: %w", err)
	}
	defer out.Close()

	// Настраиваем качество JPEG
	opts := jpeg.Options{
		Quality: 85, // Баланс между качеством и размером
	}

	// Конвертируем в JPEG
	if err := jpeg.Encode(out, resized, &opts); err != nil {
		return "", fmt.Errorf("error encoding JPEG: %w", err)
	}

	// Проверяем размер получившегося файла
	fi, err := out.Stat()
	if err == nil {
		log.Printf("Processed image saved: size = %d bytes, path = %s", fi.Size(), filePath)
	}

	return fileName, nil
}

func (s *CarService) AddCarImage(ctx context.Context, image *models.CarImage) (int, error) {
	return s.storage.AddCarImage(ctx, image)
}

func (s *CarService) GetCarImages(ctx context.Context, carID string) ([]models.CarImage, error) {
	return s.storage.GetCarImages(ctx, carID)
}
func (s *CarService) CreateBooking(ctx context.Context, booking *models.CarBooking) error {
	return s.storage.CreateCarBooking(ctx, booking)
}
