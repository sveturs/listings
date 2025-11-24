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

// ChatRepository defines operations for chat management
type ChatRepository interface {
	// Core CRUD operations
	Create(ctx context.Context, chat *domain.Chat) error
	GetByID(ctx context.Context, chatID int64) (*domain.Chat, error)
	GetByParticipantsAndListing(ctx context.Context, buyerID, sellerID, listingID int64) (*domain.Chat, error)
	GetByParticipantsAndProduct(ctx context.Context, buyerID, sellerID, productID int64) (*domain.Chat, error)
	GetByParticipantsDirect(ctx context.Context, userID1, userID2 int64) (*domain.Chat, error)
	Update(ctx context.Context, chat *domain.Chat) error
	Delete(ctx context.Context, chatID int64) error

	// User chat operations
	GetUserChats(ctx context.Context, userID int64, status *domain.ChatStatus, archivedOnly bool, limit, offset int) ([]*domain.Chat, int, error)
	GetUserChatsCount(ctx context.Context, userID int64, status *domain.ChatStatus, archivedOnly bool) (int, error)

	// Status operations
	ArchiveChat(ctx context.Context, chatID int64, archived bool) error
	UpdateStatus(ctx context.Context, chatID int64, status domain.ChatStatus) error
	UpdateLastMessageAt(ctx context.Context, chatID int64, timestamp time.Time) error

	// Unread count operations
	GetUnreadCount(ctx context.Context, userID int64, chatID *int64) (int32, error)

	// Authorization helpers
	IsParticipant(ctx context.Context, chatID, userID int64) (bool, error)

	// Transaction support
	WithTx(tx pgx.Tx) ChatRepository
}

// chatRepository implements ChatRepository using PostgreSQL
type chatRepository struct {
	db     dbOrTx
	logger zerolog.Logger
}

// NewChatRepository creates a new chat repository
func NewChatRepository(pool *pgxpool.Pool, logger zerolog.Logger) ChatRepository {
	return &chatRepository{
		db:     pool,
		logger: logger.With().Str("component", "chat_repository").Logger(),
	}
}

// WithTx returns a new repository instance using the provided transaction
func (r *chatRepository) WithTx(tx pgx.Tx) ChatRepository {
	return &chatRepository{
		db:     tx,
		logger: r.logger,
	}
}

// Create creates a new chat
// TODO: Create index: CREATE INDEX idx_chats_buyer_id ON chats(buyer_id);
// TODO: Create index: CREATE INDEX idx_chats_seller_id ON chats(seller_id);
// TODO: Create index: CREATE INDEX idx_chats_listing_id ON chats(listing_id) WHERE listing_id IS NOT NULL;
// TODO: Create index: CREATE INDEX idx_chats_storefront_product_id ON chats(storefront_product_id) WHERE storefront_product_id IS NOT NULL;
// TODO: Create index: CREATE INDEX idx_chats_last_message_at ON chats(last_message_at DESC) WHERE NOT is_archived;
func (r *chatRepository) Create(ctx context.Context, chat *domain.Chat) error {
	// Set default status if not provided
	if chat.Status == "" {
		chat.Status = domain.ChatStatusActive
	}

	query := `
		INSERT INTO chats (buyer_id, seller_id, listing_id, storefront_product_id, status, is_archived, last_message_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		chat.BuyerID,
		chat.SellerID,
		chat.ListingID,
		chat.StorefrontProductID,
		chat.Status,
		chat.IsArchived,
		chat.LastMessageAt,
	).Scan(&chat.ID, &chat.CreatedAt, &chat.UpdatedAt)

	if err != nil {
		r.logger.Error().Err(err).
			Int64("buyer_id", chat.BuyerID).
			Int64("seller_id", chat.SellerID).
			Msg("failed to create chat")
		return fmt.Errorf("failed to create chat: %w", err)
	}

	r.logger.Info().Int64("chat_id", chat.ID).
		Int64("buyer_id", chat.BuyerID).
		Int64("seller_id", chat.SellerID).
		Msg("chat created")
	return nil
}

// GetByID retrieves a chat by its ID
func (r *chatRepository) GetByID(ctx context.Context, chatID int64) (*domain.Chat, error) {
	query := `
		SELECT id, buyer_id, seller_id, listing_id, storefront_product_id,
		       status, is_archived, last_message_at, created_at, updated_at
		FROM chats
		WHERE id = $1
	`

	var chat domain.Chat
	var listingID, storefrontProductID sql.NullInt64

	err := r.db.QueryRow(ctx, query, chatID).Scan(
		&chat.ID,
		&chat.BuyerID,
		&chat.SellerID,
		&listingID,
		&storefrontProductID,
		&chat.Status,
		&chat.IsArchived,
		&chat.LastMessageAt,
		&chat.CreatedAt,
		&chat.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("chat not found")
		}
		r.logger.Error().Err(err).Int64("chat_id", chatID).Msg("failed to get chat by ID")
		return nil, fmt.Errorf("failed to get chat by ID: %w", err)
	}

	// Handle nullable fields
	if listingID.Valid {
		chat.ListingID = &listingID.Int64
	}
	if storefrontProductID.Valid {
		chat.StorefrontProductID = &storefrontProductID.Int64
	}

	return &chat, nil
}

// GetByParticipantsAndListing retrieves a chat by participants and listing
// Note: Search is symmetric - checks both (buyer_id, seller_id) and (seller_id, buyer_id) combinations
func (r *chatRepository) GetByParticipantsAndListing(ctx context.Context, buyerID, sellerID, listingID int64) (*domain.Chat, error) {
	query := `
		SELECT id, buyer_id, seller_id, listing_id, storefront_product_id,
		       status, is_archived, last_message_at, created_at, updated_at
		FROM chats
		WHERE listing_id = $3 AND (
		    (buyer_id = $1 AND seller_id = $2) OR
		    (buyer_id = $2 AND seller_id = $1)
		)
	`

	var chat domain.Chat
	var listingIDNullable, storefrontProductID sql.NullInt64

	err := r.db.QueryRow(ctx, query, buyerID, sellerID, listingID).Scan(
		&chat.ID,
		&chat.BuyerID,
		&chat.SellerID,
		&listingIDNullable,
		&storefrontProductID,
		&chat.Status,
		&chat.IsArchived,
		&chat.LastMessageAt,
		&chat.CreatedAt,
		&chat.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("chat not found")
		}
		r.logger.Error().Err(err).
			Int64("buyer_id", buyerID).
			Int64("seller_id", sellerID).
			Int64("listing_id", listingID).
			Msg("failed to get chat by participants and listing")
		return nil, fmt.Errorf("failed to get chat by participants and listing: %w", err)
	}

	// Handle nullable fields
	if listingIDNullable.Valid {
		chat.ListingID = &listingIDNullable.Int64
	}
	if storefrontProductID.Valid {
		chat.StorefrontProductID = &storefrontProductID.Int64
	}

	return &chat, nil
}

// GetByParticipantsAndProduct retrieves a chat by participants and storefront product
// Note: Search is symmetric - checks both (buyer_id, seller_id) and (seller_id, buyer_id) combinations
func (r *chatRepository) GetByParticipantsAndProduct(ctx context.Context, buyerID, sellerID, productID int64) (*domain.Chat, error) {
	query := `
		SELECT id, buyer_id, seller_id, listing_id, storefront_product_id,
		       status, is_archived, last_message_at, created_at, updated_at
		FROM chats
		WHERE storefront_product_id = $3 AND (
		    (buyer_id = $1 AND seller_id = $2) OR
		    (buyer_id = $2 AND seller_id = $1)
		)
	`

	var chat domain.Chat
	var listingID, storefrontProductIDNullable sql.NullInt64

	err := r.db.QueryRow(ctx, query, buyerID, sellerID, productID).Scan(
		&chat.ID,
		&chat.BuyerID,
		&chat.SellerID,
		&listingID,
		&storefrontProductIDNullable,
		&chat.Status,
		&chat.IsArchived,
		&chat.LastMessageAt,
		&chat.CreatedAt,
		&chat.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("chat not found")
		}
		r.logger.Error().Err(err).
			Int64("buyer_id", buyerID).
			Int64("seller_id", sellerID).
			Int64("product_id", productID).
			Msg("failed to get chat by participants and product")
		return nil, fmt.Errorf("failed to get chat by participants and product: %w", err)
	}

	// Handle nullable fields
	if listingID.Valid {
		chat.ListingID = &listingID.Int64
	}
	if storefrontProductIDNullable.Valid {
		chat.StorefrontProductID = &storefrontProductIDNullable.Int64
	}

	return &chat, nil
}

// GetByParticipantsDirect retrieves a direct message chat by participants (no listing/product context)
func (r *chatRepository) GetByParticipantsDirect(ctx context.Context, userID1, userID2 int64) (*domain.Chat, error) {
	query := `
		SELECT id, buyer_id, seller_id, listing_id, storefront_product_id,
		       status, is_archived, last_message_at, created_at, updated_at
		FROM chats
		WHERE listing_id IS NULL
		  AND storefront_product_id IS NULL
		  AND (
		      (buyer_id = $1 AND seller_id = $2) OR
		      (buyer_id = $2 AND seller_id = $1)
		  )
	`

	var chat domain.Chat
	var listingID, storefrontProductID sql.NullInt64

	err := r.db.QueryRow(ctx, query, userID1, userID2).Scan(
		&chat.ID,
		&chat.BuyerID,
		&chat.SellerID,
		&listingID,
		&storefrontProductID,
		&chat.Status,
		&chat.IsArchived,
		&chat.LastMessageAt,
		&chat.CreatedAt,
		&chat.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("chat not found")
		}
		r.logger.Error().Err(err).
			Int64("user_id1", userID1).
			Int64("user_id2", userID2).
			Msg("failed to get direct chat by participants")
		return nil, fmt.Errorf("failed to get direct chat by participants: %w", err)
	}

	// Handle nullable fields
	if listingID.Valid {
		chat.ListingID = &listingID.Int64
	}
	if storefrontProductID.Valid {
		chat.StorefrontProductID = &storefrontProductID.Int64
	}

	return &chat, nil
}

// Update updates an existing chat
func (r *chatRepository) Update(ctx context.Context, chat *domain.Chat) error {
	query := `
		UPDATE chats
		SET buyer_id = $1, seller_id = $2, listing_id = $3, storefront_product_id = $4,
		    status = $5, is_archived = $6, last_message_at = $7, updated_at = CURRENT_TIMESTAMP
		WHERE id = $8
		RETURNING updated_at
	`

	err := r.db.QueryRow(ctx, query,
		chat.BuyerID,
		chat.SellerID,
		chat.ListingID,
		chat.StorefrontProductID,
		chat.Status,
		chat.IsArchived,
		chat.LastMessageAt,
		chat.ID,
	).Scan(&chat.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("chat not found")
		}
		r.logger.Error().Err(err).Int64("chat_id", chat.ID).Msg("failed to update chat")
		return fmt.Errorf("failed to update chat: %w", err)
	}

	r.logger.Info().Int64("chat_id", chat.ID).Msg("chat updated")
	return nil
}

// Delete deletes a chat by ID
func (r *chatRepository) Delete(ctx context.Context, chatID int64) error {
	query := `DELETE FROM chats WHERE id = $1`

	result, err := r.db.Exec(ctx, query, chatID)
	if err != nil {
		r.logger.Error().Err(err).Int64("chat_id", chatID).Msg("failed to delete chat")
		return fmt.Errorf("failed to delete chat: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("chat not found")
	}

	r.logger.Info().Int64("chat_id", chatID).Msg("chat deleted")
	return nil
}

// GetUserChats retrieves all chats for a user with pagination and filters
func (r *chatRepository) GetUserChats(ctx context.Context, userID int64, status *domain.ChatStatus, archivedOnly bool, limit, offset int) ([]*domain.Chat, int, error) {
	// Build query with filters
	baseQuery := `
		FROM chats
		WHERE (buyer_id = $1 OR seller_id = $1)
	`
	args := []interface{}{userID}
	argIndex := 2

	if status != nil {
		baseQuery += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, *status)
		argIndex++
	}

	if archivedOnly {
		baseQuery += " AND is_archived = true"
	} else {
		baseQuery += " AND is_archived = false"
	}

	// Get total count
	countQuery := "SELECT COUNT(*) " + baseQuery
	var totalCount int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		r.logger.Error().Err(err).Int64("user_id", userID).Msg("failed to count user chats")
		return nil, 0, fmt.Errorf("failed to count user chats: %w", err)
	}

	// Get chats
	selectQuery := `
		SELECT id, buyer_id, seller_id, listing_id, storefront_product_id,
		       status, is_archived, last_message_at, created_at, updated_at
	` + baseQuery + fmt.Sprintf(`
		ORDER BY last_message_at DESC
		LIMIT $%d OFFSET $%d
	`, argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, selectQuery, args...)
	if err != nil {
		r.logger.Error().Err(err).Int64("user_id", userID).Msg("failed to get user chats")
		return nil, 0, fmt.Errorf("failed to get user chats: %w", err)
	}
	defer rows.Close()

	var chats []*domain.Chat
	for rows.Next() {
		var chat domain.Chat
		var listingID, storefrontProductID sql.NullInt64

		err := rows.Scan(
			&chat.ID,
			&chat.BuyerID,
			&chat.SellerID,
			&listingID,
			&storefrontProductID,
			&chat.Status,
			&chat.IsArchived,
			&chat.LastMessageAt,
			&chat.CreatedAt,
			&chat.UpdatedAt,
		)
		if err != nil {
			r.logger.Error().Err(err).Msg("failed to scan chat")
			return nil, 0, fmt.Errorf("failed to scan chat: %w", err)
		}

		// Handle nullable fields
		if listingID.Valid {
			chat.ListingID = &listingID.Int64
		}
		if storefrontProductID.Valid {
			chat.StorefrontProductID = &storefrontProductID.Int64
		}

		chats = append(chats, &chat)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error().Err(err).Msg("error iterating chat rows")
		return nil, 0, fmt.Errorf("error iterating chat rows: %w", err)
	}

	return chats, totalCount, nil
}

// GetUserChatsCount retrieves the total count of user chats with filters
func (r *chatRepository) GetUserChatsCount(ctx context.Context, userID int64, status *domain.ChatStatus, archivedOnly bool) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM chats
		WHERE (buyer_id = $1 OR seller_id = $1)
	`
	args := []interface{}{userID}
	argIndex := 2

	if status != nil {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, *status)
		argIndex++
	}

	if archivedOnly {
		query += " AND is_archived = true"
	} else {
		query += " AND is_archived = false"
	}

	var count int
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		r.logger.Error().Err(err).Int64("user_id", userID).Msg("failed to count user chats")
		return 0, fmt.Errorf("failed to count user chats: %w", err)
	}

	return count, nil
}

// ArchiveChat archives or unarchives a chat
func (r *chatRepository) ArchiveChat(ctx context.Context, chatID int64, archived bool) error {
	query := `
		UPDATE chats
		SET is_archived = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := r.db.Exec(ctx, query, archived, chatID)
	if err != nil {
		r.logger.Error().Err(err).Int64("chat_id", chatID).Bool("archived", archived).Msg("failed to archive chat")
		return fmt.Errorf("failed to archive chat: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("chat not found")
	}

	r.logger.Info().Int64("chat_id", chatID).Bool("archived", archived).Msg("chat archive status updated")
	return nil
}

// UpdateStatus updates the status of a chat
func (r *chatRepository) UpdateStatus(ctx context.Context, chatID int64, status domain.ChatStatus) error {
	query := `
		UPDATE chats
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := r.db.Exec(ctx, query, status, chatID)
	if err != nil {
		r.logger.Error().Err(err).Int64("chat_id", chatID).Str("status", string(status)).Msg("failed to update chat status")
		return fmt.Errorf("failed to update chat status: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("chat not found")
	}

	r.logger.Info().Int64("chat_id", chatID).Str("status", string(status)).Msg("chat status updated")
	return nil
}

// UpdateLastMessageAt updates the last_message_at timestamp of a chat
func (r *chatRepository) UpdateLastMessageAt(ctx context.Context, chatID int64, timestamp time.Time) error {
	query := `
		UPDATE chats
		SET last_message_at = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	result, err := r.db.Exec(ctx, query, timestamp, chatID)
	if err != nil {
		r.logger.Error().Err(err).Int64("chat_id", chatID).Msg("failed to update last_message_at")
		return fmt.Errorf("failed to update last_message_at: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("chat not found")
	}

	r.logger.Debug().Int64("chat_id", chatID).Msg("last_message_at updated")
	return nil
}

// GetUnreadCount retrieves unread message count for a user
// If chatID is provided, count only for that chat. Otherwise, count across all chats.
// TODO: Create index: CREATE INDEX idx_messages_receiver_unread ON messages(receiver_id, is_read) WHERE NOT is_read;
// TODO: Create index: CREATE INDEX idx_messages_chat_receiver_unread ON messages(chat_id, receiver_id) WHERE NOT is_read;
func (r *chatRepository) GetUnreadCount(ctx context.Context, userID int64, chatID *int64) (int32, error) {
	var query string
	var args []interface{}

	if chatID != nil {
		query = `
			SELECT COUNT(*)
			FROM messages
			WHERE chat_id = $1 AND receiver_id = $2 AND is_read = false
		`
		args = []interface{}{*chatID, userID}
	} else {
		query = `
			SELECT COUNT(*)
			FROM messages
			WHERE receiver_id = $1 AND is_read = false
		`
		args = []interface{}{userID}
	}

	var count int32
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		r.logger.Error().Err(err).Int64("user_id", userID).Msg("failed to get unread count")
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}

	return count, nil
}

// IsParticipant checks if a user is a participant in the chat
func (r *chatRepository) IsParticipant(ctx context.Context, chatID, userID int64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM chats
			WHERE id = $1 AND (buyer_id = $2 OR seller_id = $2)
		)
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, chatID, userID).Scan(&exists)
	if err != nil {
		r.logger.Error().Err(err).Int64("chat_id", chatID).Int64("user_id", userID).Msg("failed to check if user is participant")
		return false, fmt.Errorf("failed to check if user is participant: %w", err)
	}

	return exists, nil
}
