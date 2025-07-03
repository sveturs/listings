package models

import (
	"time"
)

// ChatAttachment представляет файл, прикрепленный к сообщению чата
type ChatAttachment struct {
	ID            int                    `json:"id"`
	MessageID     int                    `json:"message_id"`
	FileType      string                 `json:"file_type"`               // image, video, document
	FilePath      string                 `json:"file_path"`               // Путь к файлу в хранилище
	FileName      string                 `json:"file_name"`               // Оригинальное имя файла
	FileSize      int64                  `json:"file_size"`               // Размер в байтах
	ContentType   string                 `json:"content_type"`            // MIME тип
	StorageType   string                 `json:"storage_type"`            // minio
	StorageBucket string                 `json:"storage_bucket"`          // chat-files
	PublicURL     string                 `json:"public_url"`              // Публичный URL для доступа
	ThumbnailURL  string                 `json:"thumbnail_url,omitempty"` // URL превью для видео
	Metadata      map[string]interface{} `json:"metadata,omitempty"`      // Дополнительные метаданные
	CreatedAt     time.Time              `json:"created_at"`
}

// FileType constants
const (
	FileTypeImage    = "image"
	FileTypeVideo    = "video"
	FileTypeDocument = "document"
)

// CreateAttachmentRequest представляет запрос на загрузку файлов
type CreateAttachmentRequest struct {
	MessageID int `json:"message_id" validate:"required"`
}
