package minio

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioConfig содержит настройки подключения к MinIO
type MinioConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
	Location        string
	PublicURL       string // URL для публичного доступа к файлам
}

// MinioClient представляет клиент MinIO для работы с файлами
type MinioClient struct {
	client     *minio.Client
	bucketName string
	location   string
	baseURL    string
}

// NewMinioClient создает новый клиент MinIO
func NewMinioClient(config MinioConfig) (*MinioClient, error) {
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка создания клиента MinIO: %w", err)
	}

	// Проверка существования бакета
	exists, err := client.BucketExists(context.Background(), config.BucketName)
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки существования бакета: %w", err)
	}

	// Создание бакета, если он не существует
	if !exists {
		err = client.MakeBucket(context.Background(), config.BucketName, minio.MakeBucketOptions{
			Region: config.Location,
		})
		if err != nil {
			return nil, fmt.Errorf("ошибка создания бакета: %w", err)
		}
		log.Printf("Успешно создан бакет: %s", config.BucketName)

		// Устанавливаем политику доступа для публичного чтения
		policy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::` + config.BucketName + `/*"]}]}`
		err = client.SetBucketPolicy(context.Background(), config.BucketName, policy)
		if err != nil {
			return nil, fmt.Errorf("ошибка установки политики бакета: %w", err)
		}
		log.Printf("Успешно установлена политика для бакета: %s", config.BucketName)
	}

	// Также создаем bucket для файлов чата если он не существует
	chatBucket := "chat-files"
	chatExists, err := client.BucketExists(context.Background(), chatBucket)
	if err != nil {
		log.Printf("Ошибка проверки существования бакета chat-files: %v", err)
	} else if !chatExists {
		err = client.MakeBucket(context.Background(), chatBucket, minio.MakeBucketOptions{
			Region: config.Location,
		})
		if err != nil {
			log.Printf("Ошибка создания бакета chat-files: %v", err)
		} else {
			log.Printf("Успешно создан бакет: %s", chatBucket)

			// Устанавливаем политику доступа для публичного чтения
			chatPolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::` + chatBucket + `/*"]}]}`
			err = client.SetBucketPolicy(context.Background(), chatBucket, chatPolicy)
			if err != nil {
				log.Printf("Ошибка установки политики бакета chat-files: %v", err)
			} else {
				log.Printf("Успешно установлена политика для бакета: %s", chatBucket)
			}
		}
	}

	// Создаем bucket для фотографий отзывов если он не существует
	reviewPhotosBucket := "review-photos"
	reviewExists, err := client.BucketExists(context.Background(), reviewPhotosBucket)
	if err != nil {
		log.Printf("Ошибка проверки существования бакета review-photos: %v", err)
	} else if !reviewExists {
		err = client.MakeBucket(context.Background(), reviewPhotosBucket, minio.MakeBucketOptions{
			Region: config.Location,
		})
		if err != nil {
			log.Printf("Ошибка создания бакета review-photos: %v", err)
		} else {
			log.Printf("Успешно создан бакет: %s", reviewPhotosBucket)

			// Устанавливаем политику доступа для публичного чтения
			reviewPolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::` + reviewPhotosBucket + `/*"]}]}`
			err = client.SetBucketPolicy(context.Background(), reviewPhotosBucket, reviewPolicy)
			if err != nil {
				log.Printf("Ошибка установки политики бакета review-photos: %v", err)
			} else {
				log.Printf("Успешно установлена политика для бакета: %s", reviewPhotosBucket)
			}
		}
	}

	// Формируем базовый URL для файлов
	var baseURL string
	if config.PublicURL != "" {
		// Используем публичный URL если он задан (например, через proxy)
		baseURL = config.PublicURL
	} else {
		// Иначе формируем URL напрямую к MinIO
		protocol := "http"
		if config.UseSSL {
			protocol = "https"
		}
		baseURL = fmt.Sprintf("%s://%s", protocol, config.Endpoint)
	}

	return &MinioClient{
		client:     client,
		bucketName: config.BucketName,
		location:   config.Location,
		baseURL:    baseURL,
	}, nil
}

// In backend/internal/storage/minio/client.go

func (m *MinioClient) UploadFile(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	// Remove leading slash if present
	if strings.HasPrefix(objectName, "/") {
		objectName = objectName[1:]
	}

	// Upload file to MinIO
	_, err := m.client.PutObject(ctx, m.bucketName, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("error uploading file to MinIO: %w", err)
	}

	// Return public URL for the file
	// Если baseURL уже содержит путь к прокси (например, http://localhost:3000),
	// то просто добавляем путь к файлу
	publicURL := fmt.Sprintf("/%s/%s", m.bucketName, objectName)
	return publicURL, nil
}

// DeleteFile удаляет файл из MinIO
func (m *MinioClient) DeleteFile(ctx context.Context, objectName string) error {
	err := m.client.RemoveObject(ctx, m.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("ошибка удаления файла из MinIO: %w", err)
	}
	return nil
}

// GetPresignedURL создает предварительно подписанный URL для доступа к файлу
func (m *MinioClient) GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	// Удаляем начальный слеш, если он есть
	if strings.HasPrefix(objectName, "/") {
		objectName = objectName[1:]
	}

	// Создаем предварительно подписанный URL
	presignedURL, err := m.client.PresignedGetObject(ctx, m.bucketName, objectName, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("ошибка создания предварительно подписанного URL: %w", err)
	}
	return presignedURL.String(), nil
}

// GetObject возвращает файл из MinIO в виде потока
func (m *MinioClient) GetObject(ctx context.Context, objectName string) (io.ReadCloser, error) {
	// Удаляем начальный слеш, если он есть
	if strings.HasPrefix(objectName, "/") {
		objectName = objectName[1:]
	}

	log.Printf("Получение объекта из MinIO: bucket=%s, object=%s", m.bucketName, objectName)

	// Получаем объект из MinIO
	obj, err := m.client.GetObject(ctx, m.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("ошибка получения объекта из MinIO: %w", err)
	}

	return obj, nil
}

// UploadToCustomBucket загружает файл в указанный бакет
func (m *MinioClient) UploadToCustomBucket(ctx context.Context, bucketName, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	// Remove leading slash if present
	if strings.HasPrefix(objectName, "/") {
		objectName = objectName[1:]
	}

	// Upload file to MinIO
	_, err := m.client.PutObject(ctx, bucketName, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("error uploading file to MinIO bucket %s: %w", bucketName, err)
	}

	// Return public URL for the file
	publicURL := fmt.Sprintf("/%s/%s", bucketName, objectName)
	return publicURL, nil
}

// DeleteFileFromCustomBucket удаляет файл из указанного бакета
func (m *MinioClient) DeleteFileFromCustomBucket(ctx context.Context, bucketName, objectName string) error {
	err := m.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("ошибка удаления файла из MinIO bucket %s: %w", bucketName, err)
	}
	return nil
}

// GetPresignedURLFromCustomBucket создает предварительно подписанный URL для файла из указанного бакета
func (m *MinioClient) GetPresignedURLFromCustomBucket(ctx context.Context, bucketName, objectName string, expiry time.Duration) (string, error) {
	if strings.HasPrefix(objectName, "/") {
		objectName = objectName[1:]
	}

	presignedURL, err := m.client.PresignedGetObject(ctx, bucketName, objectName, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("ошибка создания предварительно подписанного URL для bucket %s: %w", bucketName, err)
	}
	return presignedURL.String(), nil
}

// GetObjectFromCustomBucket возвращает файл из указанного бакета в виде потока
func (m *MinioClient) GetObjectFromCustomBucket(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error) {
	if strings.HasPrefix(objectName, "/") {
		objectName = objectName[1:]
	}

	log.Printf("Получение объекта из MinIO: bucket=%s, object=%s", bucketName, objectName)

	obj, err := m.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("ошибка получения объекта из MinIO bucket %s: %w", bucketName, err)
	}

	return obj, nil
}
