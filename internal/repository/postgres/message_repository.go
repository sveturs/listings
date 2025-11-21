// Package postgres implements PostgreSQL repository layer for listings microservice.
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"github.com/sveturs/listings/internal/domain"
)

// MessageRepository defines operations for message management
type MessageRepository interface {
	// Core CRUD operations
	Create(ctx context.Context, message *domain.Message) error
	GetByID(ctx context.Context, messageID int64) (*domain.Message, error)
	Update(ctx context.Context, message *domain.Message) error
	Delete(ctx context.Context, messageID int64) error

	// Query operations with cursor-based pagination
	GetMessages(ctx context.Context, chatID int64, beforeMessageID, afterMessageID *int64, limit int) ([]*domain.Message, error)
	GetMessagesByCursor(ctx context.Context, chatID int64, beforeMessageID *int64, limit int) ([]*domain.Message, bool, error)
	GetLatestMessage(ctx context.Context, chatID int64) (*domain.Message, error)

	// Read status operations
	MarkAsRead(ctx context.Context, messageID int64) error
	MarkMessagesAsRead(ctx context.Context, chatID, receiverID int64, messageIDs []int64) (int, error)
	MarkAllAsRead(ctx context.Context, chatID, receiverID int64) (int, error)

	// Count operations
	GetUnreadCount(ctx context.Context, chatID, receiverID int64) (int32, error)
	GetUnreadCountByUser(ctx context.Context, receiverID int64) (int32, error)
	GetMessagesCount(ctx context.Context, chatID int64) (int, error)

	// Batch operations
	GetMessagesByIDs(ctx context.Context, messageIDs []int64) ([]*domain.Message, error)

	// Authorization helpers
	GetMessageSenderID(ctx context.Context, messageID int64) (int64, error)

	// Transaction support
	WithTx(tx pgx.Tx) MessageRepository
}

// messageRepository implements MessageRepository using PostgreSQL
type messageRepository struct {
	db     dbOrTx
	logger zerolog.Logger
}

// NewMessageRepository creates a new message repository
func NewMessageRepository(pool *pgxpool.Pool, logger zerolog.Logger) MessageRepository {
	return &messageRepository{
		db:     pool,
		logger: logger.With().Str("component", "message_repository").Logger(),
	}
}

// WithTx returns a new repository instance using the provided transaction
func (r *messageRepository) WithTx(tx pgx.Tx) MessageRepository {
	return &messageRepository{
		db:     tx,
		logger: r.logger,
	}
}

// Create creates a new message
// TODO: Create index: CREATE INDEX idx_messages_chat_id ON messages(chat_id, created_at DESC);
// TODO: Create index: CREATE INDEX idx_messages_chat_id_id ON messages(chat_id, id DESC);
// TODO: Create index: CREATE INDEX idx_messages_sender_id ON messages(sender_id);
// TODO: Create index: CREATE INDEX idx_messages_receiver_id ON messages(receiver_id);
// TODO: Create index: CREATE INDEX idx_messages_created_at ON messages(created_at);
func (r *messageRepository) Create(ctx context.Context, message *domain.Message) error {
	// Set default values if not provided
	if message.OriginalLanguage == "" {
		message.OriginalLanguage = "en"
	}
	if message.Status == "" {
		message.Status = domain.MessageStatusSent
	}

	query := `
		INSERT INTO messages (
			chat_id, sender_id, receiver_id, content, original_language,
			listing_id, storefront_product_id, status, is_read,
			has_attachments, attachments_count
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		message.ChatID,
		message.SenderID,
		message.ReceiverID,
		message.Content,
		message.OriginalLanguage,
		message.ListingID,
		message.StorefrontProductID,
		message.Status,
		message.IsRead,
		message.HasAttachments,
		message.AttachmentsCount,
	).Scan(&message.ID, &message.CreatedAt, &message.UpdatedAt)

	if err != nil {
		r.logger.Error().Err(err).
			Int64("chat_id", message.ChatID).
			Int64("sender_id", message.SenderID).
			Msg("failed to create message")
		return fmt.Errorf("failed to create message: %w", err)
	}

	r.logger.Info().
		Int64("message_id", message.ID).
		Int64("chat_id", message.ChatID).
		Int64("sender_id", message.SenderID).
		Msg("message created")
	return nil
}

// GetByID retrieves a message by its ID
func (r *messageRepository) GetByID(ctx context.Context, messageID int64) (*domain.Message, error) {
	query := `
		SELECT id, chat_id, sender_id, receiver_id, content, original_language,
		       listing_id, storefront_product_id, status, is_read,
		       has_attachments, attachments_count, created_at, updated_at, read_at
		FROM messages
		WHERE id = $1
	`

	var message domain.Message
	var listingID, storefrontProductID sql.NullInt64
	var readAt sql.NullTime

	err := r.db.QueryRow(ctx, query, messageID).Scan(
		&message.ID,
		&message.ChatID,
		&message.SenderID,
		&message.ReceiverID,
		&message.Content,
		&message.OriginalLanguage,
		&listingID,
		&storefrontProductID,
		&message.Status,
		&message.IsRead,
		&message.HasAttachments,
		&message.AttachmentsCount,
		&message.CreatedAt,
		&message.UpdatedAt,
		&readAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("message not found")
		}
		r.logger.Error().Err(err).Int64("message_id", messageID).Msg("failed to get message by ID")
		return nil, fmt.Errorf("failed to get message by ID: %w", err)
	}

	// Handle nullable fields
	if listingID.Valid {
		message.ListingID = &listingID.Int64
	}
	if storefrontProductID.Valid {
		message.StorefrontProductID = &storefrontProductID.Int64
	}
	if readAt.Valid {
		message.ReadAt = &readAt.Time
	}

	return &message, nil
}

// Update updates an existing message
func (r *messageRepository) Update(ctx context.Context, message *domain.Message) error {
	query := `
		UPDATE messages
		SET content = $1, original_language = $2, status = $3, is_read = $4,
		    has_attachments = $5, attachments_count = $6, read_at = $7,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $8
		RETURNING updated_at
	`

	err := r.db.QueryRow(ctx, query,
		message.Content,
		message.OriginalLanguage,
		message.Status,
		message.IsRead,
		message.HasAttachments,
		message.AttachmentsCount,
		message.ReadAt,
		message.ID,
	).Scan(&message.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("message not found")
		}
		r.logger.Error().Err(err).Int64("message_id", message.ID).Msg("failed to update message")
		return fmt.Errorf("failed to update message: %w", err)
	}

	r.logger.Info().Int64("message_id", message.ID).Msg("message updated")
	return nil
}

// Delete deletes a message by ID
func (r *messageRepository) Delete(ctx context.Context, messageID int64) error {
	query := `DELETE FROM messages WHERE id = $1`

	result, err := r.db.Exec(ctx, query, messageID)
	if err != nil {
		r.logger.Error().Err(err).Int64("message_id", messageID).Msg("failed to delete message")
		return fmt.Errorf("failed to delete message: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("message not found")
	}

	r.logger.Info().Int64("message_id", messageID).Msg("message deleted")
	return nil
}

// GetMessages retrieves messages with cursor-based pagination
// beforeMessageID: get messages before this ID (older messages, for scrolling up)
// afterMessageID: get messages after this ID (newer messages, for real-time updates)
// If both are nil, get most recent messages
func (r *messageRepository) GetMessages(ctx context.Context, chatID int64, beforeMessageID, afterMessageID *int64, limit int) ([]*domain.Message, error) {
	var query string
	var args []interface{}

	if beforeMessageID != nil {
		// Get messages before this ID (older)
		query = `
			SELECT id, chat_id, sender_id, receiver_id, content, original_language,
			       listing_id, storefront_product_id, status, is_read,
			       has_attachments, attachments_count, created_at, updated_at, read_at
			FROM messages
			WHERE chat_id = $1 AND id < $2
			ORDER BY id DESC
			LIMIT $3
		`
		args = []interface{}{chatID, *beforeMessageID, limit}
	} else if afterMessageID != nil {
		// Get messages after this ID (newer)
		query = `
			SELECT id, chat_id, sender_id, receiver_id, content, original_language,
			       listing_id, storefront_product_id, status, is_read,
			       has_attachments, attachments_count, created_at, updated_at, read_at
			FROM messages
			WHERE chat_id = $1 AND id > $2
			ORDER BY id ASC
			LIMIT $3
		`
		args = []interface{}{chatID, *afterMessageID, limit}
	} else {
		// Get most recent messages
		query = `
			SELECT id, chat_id, sender_id, receiver_id, content, original_language,
			       listing_id, storefront_product_id, status, is_read,
			       has_attachments, attachments_count, created_at, updated_at, read_at
			FROM messages
			WHERE chat_id = $1
			ORDER BY id DESC
			LIMIT $2
		`
		args = []interface{}{chatID, limit}
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		r.logger.Error().Err(err).Int64("chat_id", chatID).Msg("failed to get messages")
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	defer rows.Close()

	var messages []*domain.Message
	for rows.Next() {
		var message domain.Message
		var listingID, storefrontProductID sql.NullInt64
		var readAt sql.NullTime

		err := rows.Scan(
			&message.ID,
			&message.ChatID,
			&message.SenderID,
			&message.ReceiverID,
			&message.Content,
			&message.OriginalLanguage,
			&listingID,
			&storefrontProductID,
			&message.Status,
			&message.IsRead,
			&message.HasAttachments,
			&message.AttachmentsCount,
			&message.CreatedAt,
			&message.UpdatedAt,
			&readAt,
		)
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to scan message")
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

		// Handle nullable fields
		if listingID.Valid {
			message.ListingID = &listingID.Int64
		}
		if storefrontProductID.Valid {
			message.StorefrontProductID = &storefrontProductID.Int64
		}
		if readAt.Valid {
			message.ReadAt = &readAt.Time
		}

		messages = append(messages, &message)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error().Err(err).Msg("error iterating message rows")
		return nil, fmt.Errorf("error iterating message rows: %w", err)
	}

	// If we fetched older messages or most recent, reverse to chronological order
	if beforeMessageID != nil || (beforeMessageID == nil && afterMessageID == nil) {
		reverseMessages(messages)
	}

	return messages, nil
}

// GetMessagesByCursor retrieves messages with cursor pagination and returns hasMore flag
func (r *messageRepository) GetMessagesByCursor(ctx context.Context, chatID int64, beforeMessageID *int64, limit int) ([]*domain.Message, bool, error) {
	// Fetch limit + 1 to check if there are more messages
	messages, err := r.GetMessages(ctx, chatID, beforeMessageID, nil, limit+1)
	if err != nil {
		return nil, false, err
	}

	hasMore := len(messages) > limit
	if hasMore {
		messages = messages[:limit]
	}

	return messages, hasMore, nil
}

// GetLatestMessage retrieves the most recent message in a chat
func (r *messageRepository) GetLatestMessage(ctx context.Context, chatID int64) (*domain.Message, error) {
	query := `
		SELECT id, chat_id, sender_id, receiver_id, content, original_language,
		       listing_id, storefront_product_id, status, is_read,
		       has_attachments, attachments_count, created_at, updated_at, read_at
		FROM messages
		WHERE chat_id = $1
		ORDER BY id DESC
		LIMIT 1
	`

	var message domain.Message
	var listingID, storefrontProductID sql.NullInt64
	var readAt sql.NullTime

	err := r.db.QueryRow(ctx, query, chatID).Scan(
		&message.ID,
		&message.ChatID,
		&message.SenderID,
		&message.ReceiverID,
		&message.Content,
		&message.OriginalLanguage,
		&listingID,
		&storefrontProductID,
		&message.Status,
		&message.IsRead,
		&message.HasAttachments,
		&message.AttachmentsCount,
		&message.CreatedAt,
		&message.UpdatedAt,
		&readAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("no messages found in chat")
		}
		r.logger.Error().Err(err).Int64("chat_id", chatID).Msg("failed to get latest message")
		return nil, fmt.Errorf("failed to get latest message: %w", err)
	}

	// Handle nullable fields
	if listingID.Valid {
		message.ListingID = &listingID.Int64
	}
	if storefrontProductID.Valid {
		message.StorefrontProductID = &storefrontProductID.Int64
	}
	if readAt.Valid {
		message.ReadAt = &readAt.Time
	}

	return &message, nil
}

// MarkAsRead marks a single message as read
func (r *messageRepository) MarkAsRead(ctx context.Context, messageID int64) error {
	query := `
		UPDATE messages
		SET is_read = true, status = $1, read_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND is_read = false
		RETURNING read_at
	`

	var readAt time.Time
	err := r.db.QueryRow(ctx, query, domain.MessageStatusRead, messageID).Scan(&readAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Message already read or not found
			return nil
		}
		r.logger.Error().Err(err).Int64("message_id", messageID).Msg("failed to mark message as read")
		return fmt.Errorf("failed to mark message as read: %w", err)
	}

	r.logger.Debug().Int64("message_id", messageID).Msg("message marked as read")
	return nil
}

// MarkMessagesAsRead marks specific messages as read for a receiver
func (r *messageRepository) MarkMessagesAsRead(ctx context.Context, chatID, receiverID int64, messageIDs []int64) (int, error) {
	if len(messageIDs) == 0 {
		return 0, nil
	}

	query := `
		UPDATE messages
		SET is_read = true, status = $1, read_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE chat_id = $2 AND receiver_id = $3 AND id = ANY($4) AND is_read = false
	`

	result, err := r.db.Exec(ctx, query, domain.MessageStatusRead, chatID, receiverID, messageIDs)
	if err != nil {
		r.logger.Error().Err(err).
			Int64("chat_id", chatID).
			Int64("receiver_id", receiverID).
			Msg("failed to mark messages as read")
		return 0, fmt.Errorf("failed to mark messages as read: %w", err)
	}

	markedCount := int(result.RowsAffected())
	r.logger.Info().
		Int64("chat_id", chatID).
		Int64("receiver_id", receiverID).
		Int("marked_count", markedCount).
		Msg("messages marked as read")

	return markedCount, nil
}

// MarkAllAsRead marks all unread messages in a chat as read for a receiver
func (r *messageRepository) MarkAllAsRead(ctx context.Context, chatID, receiverID int64) (int, error) {
	query := `
		UPDATE messages
		SET is_read = true, status = $1, read_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE chat_id = $2 AND receiver_id = $3 AND is_read = false
	`

	result, err := r.db.Exec(ctx, query, domain.MessageStatusRead, chatID, receiverID)
	if err != nil {
		r.logger.Error().Err(err).
			Int64("chat_id", chatID).
			Int64("receiver_id", receiverID).
			Msg("failed to mark all messages as read")
		return 0, fmt.Errorf("failed to mark all messages as read: %w", err)
	}

	markedCount := int(result.RowsAffected())
	r.logger.Info().
		Int64("chat_id", chatID).
		Int64("receiver_id", receiverID).
		Int("marked_count", markedCount).
		Msg("all messages marked as read")

	return markedCount, nil
}

// GetUnreadCount retrieves unread message count for a specific chat and receiver
func (r *messageRepository) GetUnreadCount(ctx context.Context, chatID, receiverID int64) (int32, error) {
	query := `
		SELECT COUNT(*)
		FROM messages
		WHERE chat_id = $1 AND receiver_id = $2 AND is_read = false
	`

	var count int32
	err := r.db.QueryRow(ctx, query, chatID, receiverID).Scan(&count)
	if err != nil {
		r.logger.Error().Err(err).
			Int64("chat_id", chatID).
			Int64("receiver_id", receiverID).
			Msg("failed to get unread count")
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}

	return count, nil
}

// GetUnreadCountByUser retrieves total unread message count across all chats for a receiver
func (r *messageRepository) GetUnreadCountByUser(ctx context.Context, receiverID int64) (int32, error) {
	query := `
		SELECT COUNT(*)
		FROM messages
		WHERE receiver_id = $1 AND is_read = false
	`

	var count int32
	err := r.db.QueryRow(ctx, query, receiverID).Scan(&count)
	if err != nil {
		r.logger.Error().Err(err).Int64("receiver_id", receiverID).Msg("failed to get user unread count")
		return 0, fmt.Errorf("failed to get user unread count: %w", err)
	}

	return count, nil
}

// GetMessagesCount retrieves total message count in a chat
func (r *messageRepository) GetMessagesCount(ctx context.Context, chatID int64) (int, error) {
	query := `SELECT COUNT(*) FROM messages WHERE chat_id = $1`

	var count int
	err := r.db.QueryRow(ctx, query, chatID).Scan(&count)
	if err != nil {
		r.logger.Error().Err(err).Int64("chat_id", chatID).Msg("failed to count messages")
		return 0, fmt.Errorf("failed to count messages: %w", err)
	}

	return count, nil
}

// GetMessagesByIDs retrieves multiple messages by their IDs
func (r *messageRepository) GetMessagesByIDs(ctx context.Context, messageIDs []int64) ([]*domain.Message, error) {
	if len(messageIDs) == 0 {
		return []*domain.Message{}, nil
	}

	query := `
		SELECT id, chat_id, sender_id, receiver_id, content, original_language,
		       listing_id, storefront_product_id, status, is_read,
		       has_attachments, attachments_count, created_at, updated_at, read_at
		FROM messages
		WHERE id = ANY($1)
		ORDER BY id ASC
	`

	rows, err := r.db.Query(ctx, query, messageIDs)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to get messages by IDs")
		return nil, fmt.Errorf("failed to get messages by IDs: %w", err)
	}
	defer rows.Close()

	var messages []*domain.Message
	for rows.Next() {
		var message domain.Message
		var listingID, storefrontProductID sql.NullInt64
		var readAt sql.NullTime

		err := rows.Scan(
			&message.ID,
			&message.ChatID,
			&message.SenderID,
			&message.ReceiverID,
			&message.Content,
			&message.OriginalLanguage,
			&listingID,
			&storefrontProductID,
			&message.Status,
			&message.IsRead,
			&message.HasAttachments,
			&message.AttachmentsCount,
			&message.CreatedAt,
			&message.UpdatedAt,
			&readAt,
		)
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to scan message")
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

		// Handle nullable fields
		if listingID.Valid {
			message.ListingID = &listingID.Int64
		}
		if storefrontProductID.Valid {
			message.StorefrontProductID = &storefrontProductID.Int64
		}
		if readAt.Valid {
			message.ReadAt = &readAt.Time
		}

		messages = append(messages, &message)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error().Err(err).Msg("error iterating message rows")
		return nil, fmt.Errorf("error iterating message rows: %w", err)
	}

	return messages, nil
}

// GetMessageSenderID retrieves the sender ID of a message (for authorization checks)
func (r *messageRepository) GetMessageSenderID(ctx context.Context, messageID int64) (int64, error) {
	query := `SELECT sender_id FROM messages WHERE id = $1`

	var senderID int64
	err := r.db.QueryRow(ctx, query, messageID).Scan(&senderID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, fmt.Errorf("message not found")
		}
		r.logger.Error().Err(err).Int64("message_id", messageID).Msg("failed to get message sender ID")
		return 0, fmt.Errorf("failed to get message sender ID: %w", err)
	}

	return senderID, nil
}

// reverseMessages reverses a slice of messages in place
func reverseMessages(messages []*domain.Message) {
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
}
