// backend/internal/storage/filestorage/interface.go
package filestorage

import (
	"context"
	"io"
	"time"
)

// StorageProvider определяет тип провайдера хранилища
type StorageProvider string

const (
	// LocalStorage для локального хранилища
	LocalStorage StorageProvider = "local"
	// MinioStorageProvider для хранилища MinIO
	MinioStorageProvider StorageProvider = "minio"
)

// FileStorageInterface определяет интерфейс для работы с файловым хранилищем
type FileStorageInterface interface {
	// UploadFile загружает файл в хранилище
	UploadFile(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error)
	
	// DeleteFile удаляет файл из хранилища
	DeleteFile(ctx context.Context, objectName string) error
	
	// GetURL возвращает URL для доступа к файлу
	GetURL(ctx context.Context, objectName string) (string, error)
	
	// GetPresignedURL создает предварительно подписанный URL для доступа к файлу (если поддерживается)
	GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error)
}