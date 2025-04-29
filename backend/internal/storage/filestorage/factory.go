package filestorage

import (
	"backend/internal/config"
	"backend/internal/storage/minio"
	"context"
	"fmt"
	"io"
//	"log"
//	"net/url"
	"path/filepath"
	"time"
)

// Factory создает и возвращает соответствующую реализацию FileStorageInterface
func NewFileStorage(cfg config.FileStorageConfig) (FileStorageInterface, error) {
	switch cfg.Provider {
	case "minio":
		minioClient, err := minio.NewMinioClient(minio.MinioConfig{
			Endpoint:        cfg.MinioEndpoint,
			AccessKeyID:     cfg.MinioAccessKey,
			SecretAccessKey: cfg.MinioSecretKey,
			UseSSL:          cfg.MinioUseSSL,
			BucketName:      cfg.MinioBucketName,
			Location:        cfg.MinioLocation,
		})
		if err != nil {
			return nil, fmt.Errorf("ошибка создания клиента MinIO: %w", err)
		}
		return &MinioStorage{
			client:         minioClient,
			publicBaseURL:  cfg.PublicBaseURL,
			minioBucketName: cfg.MinioBucketName,
		}, nil
	default:
		return nil, fmt.Errorf("неподдерживаемый провайдер хранилища: %s", cfg.Provider)
	}
}

// MinioStorage реализация FileStorageInterface для MinIO
type MinioStorage struct {
	client         *minio.MinioClient
	publicBaseURL  string
	minioBucketName string
}

func (s *MinioStorage) UploadFile(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	filePath, err := s.client.UploadFile(ctx, objectName, reader, size, contentType)
	if err != nil {
		return "", err
	}

	// Строим публичный URL для файла
	if s.publicBaseURL != "" {
		return s.publicBaseURL + filePath, nil
	}
	
	// Если publicBaseURL не указан, используем путь внутри бакета
	return filePath, nil
}

func (s *MinioStorage) DeleteFile(ctx context.Context, objectName string) error {
	// Извлекаем имя объекта из полного пути
	if filepath.IsAbs(objectName) || objectName[0] == '/' {
		objectName = filepath.Base(objectName)
	}
	
	return s.client.DeleteFile(ctx, objectName)
}

func (s *MinioStorage) GetURL(ctx context.Context, objectName string) (string, error) {
	// Извлекаем имя объекта из полного пути
	if filepath.IsAbs(objectName) || objectName[0] == '/' {
		objectName = filepath.Base(objectName)
	}
	
	if s.publicBaseURL != "" {
		return fmt.Sprintf("%s/%s/%s", s.publicBaseURL, s.minioBucketName, objectName), nil
	}
	
	return fmt.Sprintf("/%s/%s", s.minioBucketName, objectName), nil
}

func (s *MinioStorage) GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	// Извлекаем имя объекта из полного пути
	if filepath.IsAbs(objectName) || objectName[0] == '/' {
		objectName = filepath.Base(objectName)
	}
	
	return s.client.GetPresignedURL(ctx, objectName, expiry)
}