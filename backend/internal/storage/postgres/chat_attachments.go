package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"backend/internal/domain/models"
)

// CreateChatAttachment создает новую запись о вложении
func (db *Database) CreateChatAttachment(ctx context.Context, attachment *models.ChatAttachment) error {
	query := `
		INSERT INTO chat_attachments (
			message_id, file_type, file_path, file_name, 
			file_size, content_type, storage_type, storage_bucket, 
			public_url, thumbnail_url, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at
	`

	err := db.pool.QueryRow(ctx, query,
		attachment.MessageID,
		attachment.FileType,
		attachment.FilePath,
		attachment.FileName,
		attachment.FileSize,
		attachment.ContentType,
		attachment.StorageType,
		attachment.StorageBucket,
		attachment.PublicURL,
		attachment.ThumbnailURL,
		attachment.Metadata,
	).Scan(&attachment.ID, &attachment.CreatedAt)
	if err != nil {
		return fmt.Errorf("error creating chat attachment: %w", err)
	}

	return nil
}

// GetChatAttachment получает информацию о вложении по ID
func (db *Database) GetChatAttachment(ctx context.Context, attachmentID int) (*models.ChatAttachment, error) {
	attachment := &models.ChatAttachment{}

	query := `
		SELECT 
			id, message_id, file_type, file_path, file_name,
			file_size, content_type, storage_type, storage_bucket,
			public_url, thumbnail_url, metadata, created_at
		FROM chat_attachments
		WHERE id = $1
	`

	err := db.pool.QueryRow(ctx, query, attachmentID).Scan(
		&attachment.ID,
		&attachment.MessageID,
		&attachment.FileType,
		&attachment.FilePath,
		&attachment.FileName,
		&attachment.FileSize,
		&attachment.ContentType,
		&attachment.StorageType,
		&attachment.StorageBucket,
		&attachment.PublicURL,
		&attachment.ThumbnailURL,
		&attachment.Metadata,
		&attachment.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting chat attachment: %w", err)
	}

	return attachment, nil
}

// GetMessageAttachments получает все вложения сообщения
func (db *Database) GetMessageAttachments(ctx context.Context, messageID int) ([]*models.ChatAttachment, error) {
	query := `
		SELECT 
			id, message_id, file_type, file_path, file_name,
			file_size, content_type, storage_type, storage_bucket,
			public_url, thumbnail_url, metadata, created_at
		FROM chat_attachments
		WHERE message_id = $1
		ORDER BY created_at ASC
	`

	rows, err := db.pool.Query(ctx, query, messageID)
	if err != nil {
		return nil, fmt.Errorf("error getting message attachments: %w", err)
	}
	defer rows.Close()

	var attachments []*models.ChatAttachment

	for rows.Next() {
		attachment := &models.ChatAttachment{}
		err := rows.Scan(
			&attachment.ID,
			&attachment.MessageID,
			&attachment.FileType,
			&attachment.FilePath,
			&attachment.FileName,
			&attachment.FileSize,
			&attachment.ContentType,
			&attachment.StorageType,
			&attachment.StorageBucket,
			&attachment.PublicURL,
			&attachment.ThumbnailURL,
			&attachment.Metadata,
			&attachment.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning attachment: %w", err)
		}

		attachments = append(attachments, attachment)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating attachments: %w", err)
	}

	return attachments, nil
}

// DeleteChatAttachment удаляет вложение по ID
func (db *Database) DeleteChatAttachment(ctx context.Context, attachmentID int) error {
	query := `DELETE FROM chat_attachments WHERE id = $1`

	result, err := db.pool.Exec(ctx, query, attachmentID)
	if err != nil {
		return fmt.Errorf("error deleting chat attachment: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("attachment not found")
	}

	return nil
}

// UpdateMessageAttachmentsCount обновляет счетчик вложений в сообщении
func (db *Database) UpdateMessageAttachmentsCount(ctx context.Context, messageID int, count int) error {
	query := `
		UPDATE marketplace_messages 
		SET has_attachments = $2, attachments_count = $3
		WHERE id = $1
	`

	hasAttachments := count > 0
	_, err := db.pool.Exec(ctx, query, messageID, hasAttachments, count)
	if err != nil {
		return fmt.Errorf("error updating message attachments count: %w", err)
	}

	return nil
}

// GetMessageByID получает сообщение по ID
func (db *Database) GetMessageByID(ctx context.Context, messageID int) (*models.MarketplaceMessage, error) {
	message := &models.MarketplaceMessage{}
	var listingID sql.NullInt64

	query := `
		SELECT 
			id, chat_id, listing_id, sender_id, receiver_id,
			content, is_read, created_at, updated_at,
			original_language, has_attachments, attachments_count
		FROM marketplace_messages
		WHERE id = $1
	`

	err := db.pool.QueryRow(ctx, query, messageID).Scan(
		&message.ID,
		&message.ChatID,
		&listingID,
		&message.SenderID,
		&message.ReceiverID,
		&message.Content,
		&message.IsRead,
		&message.CreatedAt,
		&message.UpdatedAt,
		&message.OriginalLanguage,
		&message.HasAttachments,
		&message.AttachmentsCount,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting message by id: %w", err)
	}

	// Обрабатываем NULL значение для listing_id
	if listingID.Valid {
		message.ListingID = int(listingID.Int64)
	}

	return message, nil
}
