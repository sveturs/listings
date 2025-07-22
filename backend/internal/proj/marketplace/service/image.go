package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"backend/internal/domain/models"
	"backend/internal/storage/filestorage"
	//	"io"
	//	"log"
)

// ImageService представляет сервис для работы с изображениями
type ImageService struct {
	fileStorage filestorage.FileStorageInterface
}

// NewImageService создает новый сервис для работы с изображениями
func NewImageService(fileStorage filestorage.FileStorageInterface) *ImageService {
	return &ImageService{
		fileStorage: fileStorage,
	}
}

// UploadImage загружает изображение в хранилище
func (s *ImageService) UploadImage(ctx context.Context, file *multipart.FileHeader, objectName string) (*models.MarketplaceImage, error) {
	// Открываем файл
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия файла: %w", err)
	}
	defer func() {
		if err := src.Close(); err != nil {
			// Логирование ошибки закрытия файла
		}
	}()

	// Загружаем файл в хранилище
	publicURL, err := s.fileStorage.UploadFile(ctx, objectName, src, file.Size, file.Header.Get("Content-Type"))
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки файла: %w", err)
	}

	// Создаем и возвращаем информацию об изображении
	image := &models.MarketplaceImage{
		FilePath:    objectName,
		FileName:    file.Filename,
		FileSize:    int(file.Size),
		ContentType: file.Header.Get("Content-Type"),
		StorageType: "minio", // В данной реализации всегда используем MinIO
		PublicURL:   publicURL,
		CreatedAt:   time.Now(),
	}

	return image, nil
}

// DeleteImage удаляет изображение из хранилища
func (s *ImageService) DeleteImage(ctx context.Context, image *models.MarketplaceImage) error {
	if image.FilePath == "" {
		return fmt.Errorf("пустой путь к файлу")
	}

	// Извлекаем имя файла из пути
	objectName := filepath.Base(image.FilePath)

	// Удаляем файл из хранилища
	err := s.fileStorage.DeleteFile(ctx, objectName)
	if err != nil {
		return fmt.Errorf("ошибка удаления файла: %w", err)
	}

	return nil
}

// GetImageURL возвращает URL для доступа к изображению
func (s *ImageService) GetImageURL(ctx context.Context, image *models.MarketplaceImage) (string, error) {
	if image.PublicURL != "" {
		return image.PublicURL, nil
	}

	if image.FilePath == "" {
		return "", fmt.Errorf("пустой путь к файлу")
	}

	// Извлекаем имя файла из пути
	objectName := filepath.Base(image.FilePath)

	// Получаем URL для доступа к файлу
	url, err := s.fileStorage.GetURL(ctx, objectName)
	if err != nil {
		return "", fmt.Errorf("ошибка получения URL: %w", err)
	}

	return url, nil
}
