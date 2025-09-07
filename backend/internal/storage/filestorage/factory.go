package filestorage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"backend/internal/config"
	"backend/internal/logger"
	"backend/internal/storage/minio"
)

// Factory создает и возвращает соответствующую реализацию FileStorageInterface
func NewFileStorage(ctx context.Context, cfg config.FileStorageConfig) (FileStorageInterface, error) {
	logger.Info().Str("provider", cfg.Provider).Msg("Инициализация хранилища файлов.")

	switch cfg.Provider {
	case "minio":
		logger.Info().Str("endpoint", cfg.MinioEndpoint).Str("bucket", cfg.MinioBucketName).
			Bool("useSSL", cfg.MinioUseSSL).Str("publicBaseURL", cfg.PublicBaseURL).
			Msgf("Настройка MinIO")

		minioClient, err := minio.NewMinioClient(ctx, minio.MinioConfig{
			Endpoint:           cfg.MinioEndpoint,
			AccessKeyID:        cfg.MinioAccessKey,
			SecretAccessKey:    cfg.MinioSecretKey,
			UseSSL:             cfg.MinioUseSSL,
			BucketName:         cfg.MinioBucketName,
			ChatBucket:         cfg.MinioChatBucket,
			StorefrontBucket:   cfg.MinioStorefrontBucket,
			ReviewPhotosBucket: cfg.MinioReviewPhotosBucket,
			Location:           cfg.MinioLocation,
			PublicURL:          cfg.PublicBaseURL,
		})
		if err != nil {
			logger.Error().Err(err).Msg("ОШИБКА при создании клиента MinIO")
			return nil, fmt.Errorf("ошибка создания клиента MinIO: %w", err)
		}
		logger.Info().Msgf("Клиент MinIO успешно создан")
		return &MinioStorage{
			client:          minioClient,
			publicBaseURL:   cfg.PublicBaseURL,
			minioBucketName: cfg.MinioBucketName,
		}, nil
	default:
		logger.Error().Str("provider", cfg.Provider).Msgf("неподдерживаемый провайдер хранилища")
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

	logger.Info().Str("objectName", objectName).Int64("size", size).Str("contentType", contentType).
		Str("minioBucketName", s.minioBucketName).Msg("Uploading file to MinIO")

	filePath, err := s.client.UploadFile(ctx, bucketPath, reader, size, contentType)
	if err != nil {
		logger.Error().Err(err).Msg("ERROR uploading file to MinIO")
		return "", err
	}

	logger.Info().Str("filePath", filePath).Msg("File successfully uploaded to MinIO")

	// Возвращаем полный URL, который уже сформирован в MinioClient
	logger.Info().Str("fileURL", filePath).Msg("URL for file")
	return filePath, nil
}

func (s *MinioStorage) DeleteFile(ctx context.Context, objectName string) error {
	// Извлекаем имя объекта из полного пути
	originalName := objectName
	if filepath.IsAbs(objectName) || objectName[0] == '/' {
		objectName = filepath.Base(objectName)
	}

	logger.Info().Str("originalName", originalName).Str("objectName", objectName).Msgf("Удаление файла из MinIO")

	err := s.client.DeleteFile(ctx, objectName)
	if err != nil {
		logger.Info().Msgf("ОШИБКА при удалении файла из MinIO: %v", err)
	} else {
		logger.Info().Msgf("Файл успешно удален из MinIO: %s", objectName)
	}

	return err
}

func (s *MinioStorage) GetURL(ctx context.Context, objectName string) (string, error) {
	// Извлекаем имя объекта из полного пути
	originalName := objectName
	if filepath.IsAbs(objectName) || objectName[0] == '/' {
		objectName = filepath.Base(objectName)
	}

	logger.Info().Msgf("Получение URL для файла: originalName=%s, objectName=%s", originalName, objectName)

	var fileURL string
	if s.publicBaseURL != "" {
		// Формируем URL для доступа через бэкенд
		// Используем имя bucket из конфигурации
		fileURL = fmt.Sprintf("%s/%s/%s", s.publicBaseURL, s.minioBucketName, objectName)
	} else {
		// Используем имя bucket из конфигурации
		fileURL = fmt.Sprintf("/%s/%s", s.minioBucketName, objectName)
	}

	logger.Info().Str("fileURL", fileURL).Msg("Сформирован URL для файла")
	return fileURL, nil
}

func (s *MinioStorage) GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	// Извлекаем имя объекта из полного пути
	originalName := objectName
	if filepath.IsAbs(objectName) || objectName[0] == '/' {
		objectName = filepath.Base(objectName)
	}

	logger.Info().Str("originalName", originalName).Str("objectName", objectName).Dur("expiry", expiry).Msg("Получение предподписанного URL для файла")

	url, err := s.client.GetPresignedURL(ctx, objectName, expiry)
	if err != nil {
		logger.Error().Err(err).Msg("ОШИБКА при получении предподписанного URL")
	} else {
		logger.Info().Str("url", url).Msg("Получен предподписанный URL")
	}

	return url, err
}

func (s *MinioStorage) GetFile(ctx context.Context, objectName string) (io.ReadCloser, error) {
	// Извлекаем имя объекта из полного пути
	originalName := objectName
	if filepath.IsAbs(objectName) || objectName[0] == '/' {
		objectName = filepath.Base(objectName)
	}

	logger.Info().Str("originalName", originalName).Str("objectName", objectName).Msg("Получение файла из MinIO")

	// Если путь содержит подкаталоги, сохраняем их
	parts := filepath.SplitList(objectName)
	if len(parts) > 1 {
		objectName = filepath.Join(parts...)
	}

	file, err := s.client.GetObject(ctx, objectName)
	if err != nil {
		logger.Error().Err(err).Msgf("ОШИБКА при получении файла из MinIO")
		return nil, err
	}

	logger.Info().Msgf("Файл успешно получен из MinIO")
	return file, nil
}
