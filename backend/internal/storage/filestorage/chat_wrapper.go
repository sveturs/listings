package filestorage

import (
	"backend/internal/storage/minio"
	"context"
	"fmt"
	"io"
	"log"
	"time"
)

// ChatFilesWrapper обертка для работы с bucket chat-files
type ChatFilesWrapper struct {
	client        *minio.MinioClient
	publicBaseURL string
}

// NewChatFilesWrapper создает новую обертку для chat-files
func NewChatFilesWrapper(client *minio.MinioClient, publicBaseURL string) *ChatFilesWrapper {
	return &ChatFilesWrapper{
		client:        client,
		publicBaseURL: publicBaseURL,
	}
}

// UploadFile загружает файл в хранилище
func (w *ChatFilesWrapper) UploadFile(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	log.Printf("Uploading chat file to MinIO: objectName=%s, size=%d, contentType=%s", objectName, size, contentType)

	// Загружаем файл через MinIO клиент
	publicURL, err := w.client.UploadFile(ctx, objectName, reader, size, contentType)
	if err != nil {
		log.Printf("ERROR uploading chat file to MinIO: %v", err)
		return "", err
	}

	log.Printf("Chat file successfully uploaded to MinIO: publicURL=%s", publicURL)
	return publicURL, nil
}

// DeleteFile удаляет файл из хранилища
func (w *ChatFilesWrapper) DeleteFile(ctx context.Context, objectName string) error {
	log.Printf("Deleting chat file from MinIO: objectName=%s", objectName)
	return w.client.DeleteFile(ctx, objectName)
}

// GetURL возвращает публичный URL файла
func (w *ChatFilesWrapper) GetURL(ctx context.Context, objectName string) (string, error) {
	// Формируем URL для доступа через nginx
	fileURL := fmt.Sprintf("/chat-files/%s", objectName)
	log.Printf("Generated URL for chat file: %s", fileURL)
	return fileURL, nil
}

// GetPresignedURL создает предподписанный URL
func (w *ChatFilesWrapper) GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	return w.client.GetPresignedURL(ctx, objectName, expiry)
}

// GetFile получает файл из хранилища
func (w *ChatFilesWrapper) GetFile(ctx context.Context, objectName string) (io.ReadCloser, error) {
	return w.client.GetObject(ctx, objectName)
}