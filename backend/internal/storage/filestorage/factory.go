package filestorage

import (
	"backend/internal/config"
	"backend/internal/storage/minio"
	"context"
	"fmt"
	"io"
	"log" // Разкомментируйте эту строку
	"path/filepath"
	"time"
	// "strings"
)

// Factory создает и возвращает соответствующую реализацию FileStorageInterface
func NewFileStorage(cfg config.FileStorageConfig) (FileStorageInterface, error) {
	log.Printf("Инициализация хранилища файлов. Провайдер: %s", cfg.Provider)

	switch cfg.Provider {
	case "minio":
		log.Printf("Настройка MinIO: endpoint=%s, bucket=%s, useSSL=%v, publicBaseURL=%s",
			cfg.MinioEndpoint, cfg.MinioBucketName, cfg.MinioUseSSL, cfg.PublicBaseURL)

		minioClient, err := minio.NewMinioClient(minio.MinioConfig{
			Endpoint:        cfg.MinioEndpoint,
			AccessKeyID:     cfg.MinioAccessKey,
			SecretAccessKey: cfg.MinioSecretKey,
			UseSSL:          cfg.MinioUseSSL,
			BucketName:      cfg.MinioBucketName,
			Location:        cfg.MinioLocation,
		})
		if err != nil {
			log.Printf("ОШИБКА при создании клиента MinIO: %v", err)
			return nil, fmt.Errorf("ошибка создания клиента MinIO: %w", err)
		}
		log.Printf("Клиент MinIO успешно создан")
		return &MinioStorage{
			client:          minioClient,
			publicBaseURL:   cfg.PublicBaseURL,
			minioBucketName: cfg.MinioBucketName,
		}, nil
	default:
		log.Printf("ОШИБКА: неподдерживаемый провайдер хранилища: %s", cfg.Provider)
		return nil, fmt.Errorf("неподдерживаемый провайдер хранилища: %s", cfg.Provider)
	}
}

// MinioStorage реализация FileStorageInterface для MinIO
type MinioStorage struct {
	client          *minio.MinioClient
	publicBaseURL   string
	minioBucketName string
}

func (s *MinioStorage) UploadFile(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	bucketPath := objectName

	log.Printf("Uploading file to MinIO: objectName=%s, size=%d, contentType=%s, bucket=%s",
		bucketPath, size, contentType, s.minioBucketName)

	filePath, err := s.client.UploadFile(ctx, bucketPath, reader, size, contentType)
	if err != nil {
		log.Printf("ERROR uploading file to MinIO: %v", err)
		return "", err
	}

	log.Printf("File successfully uploaded to MinIO: filePath=%s", filePath)

	// Формируем URL в зависимости от bucket
	var fileURL string
	if s.minioBucketName == "chat-files" {
		fileURL = fmt.Sprintf("/chat-files/%s", objectName)
	} else {
		fileURL = fmt.Sprintf("/listings/%s", objectName)
	}

	log.Printf("URL for file: %s", fileURL)
	return fileURL, nil
}

func (s *MinioStorage) DeleteFile(ctx context.Context, objectName string) error {
	// Извлекаем имя объекта из полного пути
	originalName := objectName
	if filepath.IsAbs(objectName) || objectName[0] == '/' {
		objectName = filepath.Base(objectName)
	}

	log.Printf("Удаление файла из MinIO: originalName=%s, objectName=%s", originalName, objectName)

	err := s.client.DeleteFile(ctx, objectName)
	if err != nil {
		log.Printf("ОШИБКА при удалении файла из MinIO: %v", err)
	} else {
		log.Printf("Файл успешно удален из MinIO: %s", objectName)
	}

	return err
}

func (s *MinioStorage) GetURL(ctx context.Context, objectName string) (string, error) {
	// Извлекаем имя объекта из полного пути
	originalName := objectName
	if filepath.IsAbs(objectName) || objectName[0] == '/' {
		objectName = filepath.Base(objectName)
	}

	log.Printf("Получение URL для файла: originalName=%s, objectName=%s", originalName, objectName)

	var fileURL string
	if s.publicBaseURL != "" {
		// Формируем URL для доступа через бэкенд
		if s.minioBucketName == "chat-files" {
			fileURL = fmt.Sprintf("%s/chat-files/%s", s.publicBaseURL, objectName)
		} else {
			fileURL = fmt.Sprintf("%s/listings/%s", s.publicBaseURL, objectName)
		}
	} else {
		if s.minioBucketName == "chat-files" {
			fileURL = fmt.Sprintf("/chat-files/%s", objectName)
		} else {
			fileURL = fmt.Sprintf("/listings/%s", objectName)
		}
	}

	log.Printf("Сформирован URL для файла: %s", fileURL)
	return fileURL, nil
}

func (s *MinioStorage) GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	// Извлекаем имя объекта из полного пути
	originalName := objectName
	if filepath.IsAbs(objectName) || objectName[0] == '/' {
		objectName = filepath.Base(objectName)
	}

	log.Printf("Получение предподписанного URL для файла: originalName=%s, objectName=%s, expiry=%v",
		originalName, objectName, expiry)

	url, err := s.client.GetPresignedURL(ctx, objectName, expiry)
	if err != nil {
		log.Printf("ОШИБКА при получении предподписанного URL: %v", err)
	} else {
		log.Printf("Получен предподписанный URL: %s", url)
	}

	return url, err
}

func (s *MinioStorage) GetFile(ctx context.Context, objectName string) (io.ReadCloser, error) {
	// Извлекаем имя объекта из полного пути
	originalName := objectName
	if filepath.IsAbs(objectName) || objectName[0] == '/' {
		objectName = filepath.Base(objectName)
	}

	log.Printf("Получение файла из MinIO: originalName=%s, objectName=%s",
		originalName, objectName)

	// Если путь содержит подкаталоги, сохраняем их
	parts := filepath.SplitList(objectName)
	if len(parts) > 1 {
		objectName = filepath.Join(parts...)
	}

	file, err := s.client.GetObject(ctx, objectName)
	if err != nil {
		log.Printf("ОШИБКА при получении файла из MinIO: %v", err)
		return nil, err
	}

	log.Printf("Файл успешно получен из MinIO")
	return file, nil
}
