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
	Endpoint           string
	AccessKeyID        string
	SecretAccessKey    string
	UseSSL             bool
	BucketName         string // Основной bucket для объявлений
	ChatBucket         string // Bucket для файлов чата
	StorefrontBucket   string // Bucket для товаров витрин
	ReviewPhotosBucket string // Bucket для фотографий отзывов
	Location           string
	PublicURL          string // URL для публичного доступа к файлам
}

// MinioClient представляет клиент MinIO для работы с файлами
type MinioClient struct {
	client     *minio.Client
	bucketName string
	location   string
	baseURL    string
}

// NewMinioClient создает новый клиент MinIO
func NewMinioClient(ctx context.Context, config MinioConfig) (*MinioClient, error) {
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка создания клиента MinIO: %w", err)
	}

	// Проверка существования бакета
	exists, err := client.BucketExists(ctx, config.BucketName)
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки существования бакета: %w", err)
	}

	// Создание бакета, если он не существует
	if !exists {
		err = client.MakeBucket(ctx, config.BucketName, minio.MakeBucketOptions{
			Region: config.Location,
		})
		if err != nil {
			return nil, fmt.Errorf("ошибка создания бакета: %w", err)
		}
		log.Printf("Успешно создан бакет: %s", config.BucketName)

		// Устанавливаем политику доступа для публичного чтения
		policy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::` + config.BucketName + `/*"]}]}`
		err = client.SetBucketPolicy(ctx, config.BucketName, policy)
		if err != nil {
			return nil, fmt.Errorf("ошибка установки политики бакета: %w", err)
		}
		log.Printf("Успешно установлена политика для бакета: %s", config.BucketName)
	}

	// Также создаем bucket для файлов чата если он не существует
	if config.ChatBucket != "" {
		chatExists, err := client.BucketExists(ctx, config.ChatBucket)
		if err != nil {
			log.Printf("Ошибка проверки существования бакета %s: %v", config.ChatBucket, err)
		} else if !chatExists {
			err = client.MakeBucket(ctx, config.ChatBucket, minio.MakeBucketOptions{
				Region: config.Location,
			})
			if err != nil {
				log.Printf("Ошибка создания бакета %s: %v", config.ChatBucket, err)
			} else {
				log.Printf("Успешно создан бакет: %s", config.ChatBucket)

				// Устанавливаем политику доступа для публичного чтения
				chatPolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::` + config.ChatBucket + `/*"]}]}`
				err = client.SetBucketPolicy(ctx, config.ChatBucket, chatPolicy)
				if err != nil {
					log.Printf("Ошибка установки политики бакета %s: %v", config.ChatBucket, err)
				} else {
					log.Printf("Успешно установлена политика для бакета: %s", config.ChatBucket)
				}
			}
		}
	}

	// Создаем bucket для фотографий отзывов если он не существует
	if config.ReviewPhotosBucket != "" {
		reviewExists, err := client.BucketExists(ctx, config.ReviewPhotosBucket)
		if err != nil {
			log.Printf("Ошибка проверки существования бакета %s: %v", config.ReviewPhotosBucket, err)
		} else if !reviewExists {
			err = client.MakeBucket(ctx, config.ReviewPhotosBucket, minio.MakeBucketOptions{
				Region: config.Location,
			})
			if err != nil {
				log.Printf("Ошибка создания бакета %s: %v", config.ReviewPhotosBucket, err)
			} else {
				log.Printf("Успешно создан бакет: %s", config.ReviewPhotosBucket)

				// Устанавливаем политику доступа для публичного чтения
				reviewPolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::` + config.ReviewPhotosBucket + `/*"]}]}`
				err = client.SetBucketPolicy(ctx, config.ReviewPhotosBucket, reviewPolicy)
				if err != nil {
					log.Printf("Ошибка установки политики бакета %s: %v", config.ReviewPhotosBucket, err)
				} else {
					log.Printf("Успешно установлена политика для бакета: %s", config.ReviewPhotosBucket)
				}
			}
		}
	}

	// Создаем bucket для товаров витрин если он не существует
	if config.StorefrontBucket != "" {
		storefrontExists, err := client.BucketExists(ctx, config.StorefrontBucket)
		if err != nil {
			log.Printf("Ошибка проверки существования бакета %s: %v", config.StorefrontBucket, err)
		} else if !storefrontExists {
			err = client.MakeBucket(ctx, config.StorefrontBucket, minio.MakeBucketOptions{
				Region: config.Location,
			})
			if err != nil {
				log.Printf("Ошибка создания бакета %s: %v", config.StorefrontBucket, err)
			} else {
				log.Printf("Успешно создан бакет: %s", config.StorefrontBucket)

				// Устанавливаем политику доступа для публичного чтения
				storefrontPolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::` + config.StorefrontBucket + `/*"]}]}`
				err = client.SetBucketPolicy(ctx, config.StorefrontBucket, storefrontPolicy)
				if err != nil {
					log.Printf("Ошибка установки политики бакета %s: %v", config.StorefrontBucket, err)
				} else {
					log.Printf("Успешно установлена политика для бакета: %s", config.StorefrontBucket)
				}
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
	objectName = strings.TrimPrefix(objectName, "/")

	// Upload file to MinIO
	_, err := m.client.PutObject(ctx, m.bucketName, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("error uploading file to MinIO: %w", err)
	}

	// Return public URL for the file
	// Формируем полный URL с базовым адресом и именем бакета
	var publicURL string
	if m.baseURL != "" {
		// Если есть базовый URL, используем его
		publicURL = fmt.Sprintf("%s/%s/%s", m.baseURL, m.bucketName, objectName)
	} else {
		// Иначе возвращаем относительный путь (для обратной совместимости)
		publicURL = fmt.Sprintf("/%s/%s", m.bucketName, objectName)
	}
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
	objectName = strings.TrimPrefix(objectName, "/")

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
	objectName = strings.TrimPrefix(objectName, "/")

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
	objectName = strings.TrimPrefix(objectName, "/")

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
	objectName = strings.TrimPrefix(objectName, "/")

	presignedURL, err := m.client.PresignedGetObject(ctx, bucketName, objectName, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("ошибка создания предварительно подписанного URL для bucket %s: %w", bucketName, err)
	}
	return presignedURL.String(), nil
}

// GetObjectFromCustomBucket возвращает файл из указанного бакета в виде потока
func (m *MinioClient) GetObjectFromCustomBucket(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error) {
	objectName = strings.TrimPrefix(objectName, "/")

	log.Printf("Получение объекта из MinIO: bucket=%s, object=%s", bucketName, objectName)

	obj, err := m.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("ошибка получения объекта из MinIO bucket %s: %w", bucketName, err)
	}

	return obj, nil
}
