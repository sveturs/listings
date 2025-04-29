package minio

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"log"
	"strings"
	"time"
)

// MinioConfig содержит настройки подключения к MinIO
type MinioConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
	Location        string
}

// MinioClient представляет клиент MinIO для работы с файлами
type MinioClient struct {
	client     *minio.Client
	bucketName string
	location   string
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

	return &MinioClient{
		client:     client,
		bucketName: config.BucketName,
		location:   config.Location,
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

	// Return path to file (relative path)
	return objectName, nil
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
	// Создаем предварительно подписанный URL
	presignedURL, err := m.client.PresignedGetObject(ctx, m.bucketName, objectName, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("ошибка создания предварительно подписанного URL: %w", err)
	}
	return presignedURL.String(), nil
}
