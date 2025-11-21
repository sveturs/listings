package domain

import (
	"encoding/json"
	"fmt"
	"time"
)

// AttachmentType represents the type of attachment
type AttachmentType string

const (
	AttachmentTypeImage    AttachmentType = "image"
	AttachmentTypeVideo    AttachmentType = "video"
	AttachmentTypeDocument AttachmentType = "document"
)

// ChatAttachment represents a file attachment in a message
type ChatAttachment struct {
	ID        int64 `json:"id"`
	MessageID int64 `json:"message_id"`

	// File metadata
	FileType    AttachmentType `json:"file_type"`
	FileName    string         `json:"file_name"`
	FileSize    int64          `json:"file_size"`
	ContentType string         `json:"content_type"`

	// Storage
	StorageType   string `json:"storage_type"`
	StorageBucket string `json:"storage_bucket"`
	FilePath      string `json:"file_path"`
	PublicURL     string `json:"public_url"`
	ThumbnailURL  *string `json:"thumbnail_url,omitempty"`

	// Metadata (JSON object for dimensions, duration, etc.)
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
}

// Validate validates the attachment fields
func (a *ChatAttachment) Validate() error {
	if a.MessageID == 0 {
		return fmt.Errorf("message_id is required")
	}

	// Validate file type
	if a.FileType != AttachmentTypeImage && a.FileType != AttachmentTypeVideo && a.FileType != AttachmentTypeDocument {
		return fmt.Errorf("invalid file_type: %s", a.FileType)
	}

	// Validate file name
	if a.FileName == "" {
		return fmt.Errorf("file_name is required")
	}
	if len(a.FileName) > 255 {
		return fmt.Errorf("file_name exceeds maximum length of 255 characters")
	}

	// Validate file size
	if a.FileSize <= 0 {
		return fmt.Errorf("file_size must be positive")
	}

	// Validate file size limits per type
	const (
		maxImageSize    = 10 * 1024 * 1024  // 10MB
		maxVideoSize    = 50 * 1024 * 1024  // 50MB
		maxDocumentSize = 20 * 1024 * 1024  // 20MB
	)

	switch a.FileType {
	case AttachmentTypeImage:
		if a.FileSize > maxImageSize {
			return fmt.Errorf("image file size exceeds maximum of 10MB")
		}
	case AttachmentTypeVideo:
		if a.FileSize > maxVideoSize {
			return fmt.Errorf("video file size exceeds maximum of 50MB")
		}
	case AttachmentTypeDocument:
		if a.FileSize > maxDocumentSize {
			return fmt.Errorf("document file size exceeds maximum of 20MB")
		}
	}

	// Validate content type
	if a.ContentType == "" {
		return fmt.Errorf("content_type is required")
	}
	if len(a.ContentType) > 100 {
		return fmt.Errorf("content_type exceeds maximum length of 100 characters")
	}

	// Validate storage fields
	if a.StorageType == "" {
		return fmt.Errorf("storage_type is required")
	}
	if a.StorageBucket == "" {
		return fmt.Errorf("storage_bucket is required")
	}
	if a.FilePath == "" {
		return fmt.Errorf("file_path is required")
	}

	return nil
}

// MetadataJSON returns metadata as JSON string for database storage
func (a *ChatAttachment) MetadataJSON() (string, error) {
	if a.Metadata == nil {
		return "{}", nil
	}

	bytes, err := json.Marshal(a.Metadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal metadata: %w", err)
	}

	return string(bytes), nil
}

// ParseMetadata parses JSON metadata string into the Metadata field
func (a *ChatAttachment) ParseMetadata(jsonStr string) error {
	if jsonStr == "" || jsonStr == "{}" {
		a.Metadata = make(map[string]interface{})
		return nil
	}

	if err := json.Unmarshal([]byte(jsonStr), &a.Metadata); err != nil {
		return fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return nil
}
