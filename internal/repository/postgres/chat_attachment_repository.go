// Package postgres implements PostgreSQL repository layer for listings microservice.
package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"github.com/vondi-global/listings/internal/domain"
)

// ChatAttachmentRepository defines operations for chat attachment management
type ChatAttachmentRepository interface {
	// Core CRUD operations
	Create(ctx context.Context, attachment *domain.ChatAttachment) error
	GetByID(ctx context.Context, attachmentID int64) (*domain.ChatAttachment, error)
	Update(ctx context.Context, attachment *domain.ChatAttachment) error
	Delete(ctx context.Context, attachmentID int64) error

	// Query operations
	GetByMessageID(ctx context.Context, messageID int64) ([]*domain.ChatAttachment, error)
	GetByMessageIDs(ctx context.Context, messageIDs []int64) (map[int64][]*domain.ChatAttachment, error)

	// Batch operations
	CreateBatch(ctx context.Context, attachments []*domain.ChatAttachment) error

	// Count operations
	GetAttachmentCount(ctx context.Context, messageID int64) (int32, error)

	// Authorization helpers
	GetAttachmentMessageID(ctx context.Context, attachmentID int64) (int64, error)

	// Transaction support
	WithTx(tx pgx.Tx) ChatAttachmentRepository
}

// chatAttachmentRepository implements ChatAttachmentRepository using PostgreSQL
type chatAttachmentRepository struct {
	db     dbOrTx
	logger zerolog.Logger
}

// NewChatAttachmentRepository creates a new chat attachment repository
func NewChatAttachmentRepository(pool *pgxpool.Pool, logger zerolog.Logger) ChatAttachmentRepository {
	return &chatAttachmentRepository{
		db:     pool,
		logger: logger.With().Str("component", "chat_attachment_repository").Logger(),
	}
}

// WithTx returns a new repository instance using the provided transaction
func (r *chatAttachmentRepository) WithTx(tx pgx.Tx) ChatAttachmentRepository {
	return &chatAttachmentRepository{
		db:     tx,
		logger: r.logger,
	}
}

// Create creates a new attachment
// TODO: Create index: CREATE INDEX idx_attachments_message_id ON chat_attachments(message_id);
// TODO: Create index: CREATE INDEX idx_attachments_file_type ON chat_attachments(file_type);
// TODO: Create index: CREATE INDEX idx_attachments_created_at ON chat_attachments(created_at);
func (r *chatAttachmentRepository) Create(ctx context.Context, attachment *domain.ChatAttachment) error {
	// Set default storage type if not provided
	if attachment.StorageType == "" {
		attachment.StorageType = "minio"
	}
	if attachment.StorageBucket == "" {
		attachment.StorageBucket = "chat-files"
	}

	query := `
		INSERT INTO chat_attachments (
			message_id, file_type, file_name, file_size, content_type,
			storage_type, storage_bucket, file_path, public_url, thumbnail_url, metadata
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(ctx, query,
		attachment.MessageID,
		attachment.FileType,
		attachment.FileName,
		attachment.FileSize,
		attachment.ContentType,
		attachment.StorageType,
		attachment.StorageBucket,
		attachment.FilePath,
		attachment.PublicURL,
		attachment.ThumbnailURL,
		attachment.Metadata,
	).Scan(&attachment.ID, &attachment.CreatedAt)

	if err != nil {
		r.logger.Error().Err(err).
			Int64("message_id", attachment.MessageID).
			Str("file_name", attachment.FileName).
			Msg("failed to create attachment")
		return fmt.Errorf("failed to create attachment: %w", err)
	}

	r.logger.Info().
		Int64("attachment_id", attachment.ID).
		Int64("message_id", attachment.MessageID).
		Str("file_name", attachment.FileName).
		Msg("attachment created")
	return nil
}

// GetByID retrieves an attachment by its ID
func (r *chatAttachmentRepository) GetByID(ctx context.Context, attachmentID int64) (*domain.ChatAttachment, error) {
	query := `
		SELECT id, message_id, file_type, file_name, file_size, content_type,
		       storage_type, storage_bucket, file_path, public_url, thumbnail_url, metadata, created_at
		FROM chat_attachments
		WHERE id = $1
	`

	var attachment domain.ChatAttachment
	var thumbnailURL sql.NullString

	err := r.db.QueryRow(ctx, query, attachmentID).Scan(
		&attachment.ID,
		&attachment.MessageID,
		&attachment.FileType,
		&attachment.FileName,
		&attachment.FileSize,
		&attachment.ContentType,
		&attachment.StorageType,
		&attachment.StorageBucket,
		&attachment.FilePath,
		&attachment.PublicURL,
		&thumbnailURL,
		&attachment.Metadata,
		&attachment.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("attachment not found")
		}
		r.logger.Error().Err(err).Int64("attachment_id", attachmentID).Msg("failed to get attachment by ID")
		return nil, fmt.Errorf("failed to get attachment by ID: %w", err)
	}

	// Handle nullable fields
	if thumbnailURL.Valid {
		attachment.ThumbnailURL = &thumbnailURL.String
	}

	return &attachment, nil
}

// Update updates an existing attachment
func (r *chatAttachmentRepository) Update(ctx context.Context, attachment *domain.ChatAttachment) error {
	query := `
		UPDATE chat_attachments
		SET file_type = $1, file_name = $2, file_size = $3, content_type = $4,
		    storage_type = $5, storage_bucket = $6, file_path = $7,
		    public_url = $8, thumbnail_url = $9, metadata = $10
		WHERE id = $11
	`

	result, err := r.db.Exec(ctx, query,
		attachment.FileType,
		attachment.FileName,
		attachment.FileSize,
		attachment.ContentType,
		attachment.StorageType,
		attachment.StorageBucket,
		attachment.FilePath,
		attachment.PublicURL,
		attachment.ThumbnailURL,
		attachment.Metadata,
		attachment.ID,
	)

	if err != nil {
		r.logger.Error().Err(err).Int64("attachment_id", attachment.ID).Msg("failed to update attachment")
		return fmt.Errorf("failed to update attachment: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("attachment not found")
	}

	r.logger.Info().Int64("attachment_id", attachment.ID).Msg("attachment updated")
	return nil
}

// Delete deletes an attachment by ID
func (r *chatAttachmentRepository) Delete(ctx context.Context, attachmentID int64) error {
	query := `DELETE FROM chat_attachments WHERE id = $1`

	result, err := r.db.Exec(ctx, query, attachmentID)
	if err != nil {
		r.logger.Error().Err(err).Int64("attachment_id", attachmentID).Msg("failed to delete attachment")
		return fmt.Errorf("failed to delete attachment: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("attachment not found")
	}

	r.logger.Info().Int64("attachment_id", attachmentID).Msg("attachment deleted")
	return nil
}

// GetByMessageID retrieves all attachments for a message
func (r *chatAttachmentRepository) GetByMessageID(ctx context.Context, messageID int64) ([]*domain.ChatAttachment, error) {
	query := `
		SELECT id, message_id, file_type, file_name, file_size, content_type,
		       storage_type, storage_bucket, file_path, public_url, thumbnail_url, metadata, created_at
		FROM chat_attachments
		WHERE message_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, query, messageID)
	if err != nil {
		r.logger.Error().Err(err).Int64("message_id", messageID).Msg("failed to get attachments by message ID")
		return nil, fmt.Errorf("failed to get attachments by message ID: %w", err)
	}
	defer rows.Close()

	var attachments []*domain.ChatAttachment
	for rows.Next() {
		var attachment domain.ChatAttachment
		var thumbnailURL sql.NullString

		err := rows.Scan(
			&attachment.ID,
			&attachment.MessageID,
			&attachment.FileType,
			&attachment.FileName,
			&attachment.FileSize,
			&attachment.ContentType,
			&attachment.StorageType,
			&attachment.StorageBucket,
			&attachment.FilePath,
			&attachment.PublicURL,
			&thumbnailURL,
			&attachment.Metadata,
			&attachment.CreatedAt,
		)
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to scan attachment")
			return nil, fmt.Errorf("failed to scan attachment: %w", err)
		}

		// Handle nullable fields
		if thumbnailURL.Valid {
			attachment.ThumbnailURL = &thumbnailURL.String
		}

		attachments = append(attachments, &attachment)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error().Err(err).Msg("error iterating attachment rows")
		return nil, fmt.Errorf("error iterating attachment rows: %w", err)
	}

	return attachments, nil
}

// GetByMessageIDs retrieves attachments for multiple messages
// Returns a map of messageID -> attachments
func (r *chatAttachmentRepository) GetByMessageIDs(ctx context.Context, messageIDs []int64) (map[int64][]*domain.ChatAttachment, error) {
	if len(messageIDs) == 0 {
		return make(map[int64][]*domain.ChatAttachment), nil
	}

	query := `
		SELECT id, message_id, file_type, file_name, file_size, content_type,
		       storage_type, storage_bucket, file_path, public_url, thumbnail_url, metadata, created_at
		FROM chat_attachments
		WHERE message_id = ANY($1)
		ORDER BY message_id, created_at ASC
	`

	rows, err := r.db.Query(ctx, query, messageIDs)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to get attachments by message IDs")
		return nil, fmt.Errorf("failed to get attachments by message IDs: %w", err)
	}
	defer rows.Close()

	result := make(map[int64][]*domain.ChatAttachment)
	for rows.Next() {
		var attachment domain.ChatAttachment
		var thumbnailURL sql.NullString

		err := rows.Scan(
			&attachment.ID,
			&attachment.MessageID,
			&attachment.FileType,
			&attachment.FileName,
			&attachment.FileSize,
			&attachment.ContentType,
			&attachment.StorageType,
			&attachment.StorageBucket,
			&attachment.FilePath,
			&attachment.PublicURL,
			&thumbnailURL,
			&attachment.Metadata,
			&attachment.CreatedAt,
		)
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to scan attachment")
			return nil, fmt.Errorf("failed to scan attachment: %w", err)
		}

		// Handle nullable fields
		if thumbnailURL.Valid {
			attachment.ThumbnailURL = &thumbnailURL.String
		}

		result[attachment.MessageID] = append(result[attachment.MessageID], &attachment)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error().Err(err).Msg("error iterating attachment rows")
		return nil, fmt.Errorf("error iterating attachment rows: %w", err)
	}

	return result, nil
}

// CreateBatch creates multiple attachments in a batch
func (r *chatAttachmentRepository) CreateBatch(ctx context.Context, attachments []*domain.ChatAttachment) error {
	if len(attachments) == 0 {
		return nil
	}

	query := `
		INSERT INTO chat_attachments (
			message_id, file_type, file_name, file_size, content_type,
			storage_type, storage_bucket, file_path, public_url, thumbnail_url, metadata
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at
	`

	batch := &pgx.Batch{}
	for _, attachment := range attachments {
		// Set defaults
		if attachment.StorageType == "" {
			attachment.StorageType = "minio"
		}
		if attachment.StorageBucket == "" {
			attachment.StorageBucket = "chat-files"
		}

		batch.Queue(query,
			attachment.MessageID,
			attachment.FileType,
			attachment.FileName,
			attachment.FileSize,
			attachment.ContentType,
			attachment.StorageType,
			attachment.StorageBucket,
			attachment.FilePath,
			attachment.PublicURL,
			attachment.ThumbnailURL,
			attachment.Metadata,
		)
	}

	br := r.db.SendBatch(ctx, batch)
	defer br.Close()

	for i := 0; i < len(attachments); i++ {
		err := br.QueryRow().Scan(&attachments[i].ID, &attachments[i].CreatedAt)
		if err != nil {
			r.logger.Error().Err(err).Int("attachment_index", i).Msg("failed to create attachment in batch")
			return fmt.Errorf("failed to create attachment in batch: %w", err)
		}
	}

	r.logger.Info().Int("count", len(attachments)).Msg("attachments created in batch")
	return nil
}

// GetAttachmentCount retrieves the count of attachments for a message
func (r *chatAttachmentRepository) GetAttachmentCount(ctx context.Context, messageID int64) (int32, error) {
	query := `SELECT COUNT(*) FROM chat_attachments WHERE message_id = $1`

	var count int32
	err := r.db.QueryRow(ctx, query, messageID).Scan(&count)
	if err != nil {
		r.logger.Error().Err(err).Int64("message_id", messageID).Msg("failed to count attachments")
		return 0, fmt.Errorf("failed to count attachments: %w", err)
	}

	return count, nil
}

// GetAttachmentMessageID retrieves the message ID of an attachment (for authorization checks)
func (r *chatAttachmentRepository) GetAttachmentMessageID(ctx context.Context, attachmentID int64) (int64, error) {
	query := `SELECT message_id FROM chat_attachments WHERE id = $1`

	var messageID int64
	err := r.db.QueryRow(ctx, query, attachmentID).Scan(&messageID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, fmt.Errorf("attachment not found")
		}
		r.logger.Error().Err(err).Int64("attachment_id", attachmentID).Msg("failed to get attachment message ID")
		return 0, fmt.Errorf("failed to get attachment message ID: %w", err)
	}

	return messageID, nil
}
