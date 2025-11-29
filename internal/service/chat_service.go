// Package service provides business logic layer for the listings microservice.
package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	authservice "github.com/vondi-global/auth/pkg/service"
	"github.com/vondi-global/listings/internal/domain"
	"github.com/vondi-global/listings/internal/repository/postgres"
)

const (
	// MaxMessageLength is the maximum allowed message length
	MaxMessageLength = 10000

	// MaxImageSize is the maximum size for image attachments (10MB)
	MaxImageSize = 10 * 1024 * 1024

	// MaxVideoSize is the maximum size for video attachments (50MB)
	MaxVideoSize = 50 * 1024 * 1024

	// MaxDocumentSize is the maximum size for document attachments (20MB)
	MaxDocumentSize = 20 * 1024 * 1024
)

// isChatNotFoundError checks if an error indicates that the chat was not found.
// Repository returns fmt.Errorf("chat not found") which doesn't match errors.Is(err, ErrChatNotFound).
func isChatNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	return err == ErrChatNotFound || strings.Contains(err.Error(), "chat not found")
}

// ChatRepository defines data access operations for chats
// This will be implemented by the repository layer
type ChatRepository interface {
	// Chat CRUD
	Create(ctx context.Context, chat *domain.Chat) error
	GetByID(ctx context.Context, chatID int64) (*domain.Chat, error)
	GetByParticipantsAndListing(ctx context.Context, buyerID, sellerID, listingID int64) (*domain.Chat, error)
	GetByParticipantsAndProduct(ctx context.Context, buyerID, sellerID, productID int64) (*domain.Chat, error)
	GetByParticipantsDirect(ctx context.Context, user1ID, user2ID int64) (*domain.Chat, error)
	Update(ctx context.Context, chat *domain.Chat) error
	Delete(ctx context.Context, chatID int64) error

	// Chat listing
	GetUserChats(ctx context.Context, userID int64, status *domain.ChatStatus, archivedOnly bool, limit, offset int) ([]*domain.Chat, int, error)
	GetUnreadCount(ctx context.Context, userID int64, chatID *int64) (int32, error)

	// Archive management
	ArchiveChat(ctx context.Context, chatID int64, archived bool) error
}

// MessageRepository defines data access operations for messages
type MessageRepository interface {
	// Message CRUD
	Create(ctx context.Context, message *domain.Message) error
	GetByID(ctx context.Context, messageID int64) (*domain.Message, error)
	Delete(ctx context.Context, messageID int64) error

	// Message listing with cursor pagination
	GetMessagesByCursor(ctx context.Context, chatID int64, beforeMessageID *int64, limit int) ([]*domain.Message, bool, error)
	GetMessages(ctx context.Context, chatID int64, beforeMessageID, afterMessageID *int64, limit int) ([]*domain.Message, error)
	GetLatestMessage(ctx context.Context, chatID int64) (*domain.Message, error)

	// Read status management
	MarkMessagesAsRead(ctx context.Context, chatID, receiverID int64, messageIDs []int64) (int, error)
	MarkAllAsRead(ctx context.Context, chatID, receiverID int64) (int, error)
	GetUnreadCount(ctx context.Context, chatID, receiverID int64) (int32, error)
	GetUnreadCountByUser(ctx context.Context, receiverID int64) (int32, error)
}

// AttachmentRepository defines data access operations for attachments
type AttachmentRepository interface {
	// Attachment CRUD
	Create(ctx context.Context, attachment *domain.ChatAttachment) error
	GetByID(ctx context.Context, attachmentID int64) (*domain.ChatAttachment, error)
	GetByMessageID(ctx context.Context, messageID int64) ([]*domain.ChatAttachment, error)
	Delete(ctx context.Context, attachmentID int64) error

	// Batch operations
	CreateBatch(ctx context.Context, attachments []*domain.ChatAttachment) error
}

// ChatService defines business logic operations for chat management
type ChatService interface {
	// Chat operations
	CreateChat(ctx context.Context, req *CreateChatRequest) (*domain.Chat, error)
	GetOrCreateChat(ctx context.Context, req *GetOrCreateChatRequest) (*domain.Chat, bool, error)
	GetChat(ctx context.Context, chatID, userID int64) (*domain.Chat, error)
	GetUserChats(ctx context.Context, req *GetUserChatsRequest) ([]*domain.Chat, int, error)
	ArchiveChat(ctx context.Context, chatID, userID int64, archived bool) error
	DeleteChat(ctx context.Context, chatID int64) error

	// Message operations
	SendMessage(ctx context.Context, req *SendMessageRequest) (*domain.Message, error)
	SendSystemMessage(ctx context.Context, req *SendSystemMessageRequest) (*domain.Message, error)
	GetMessages(ctx context.Context, req *GetMessagesRequest) ([]*domain.Message, bool, error)
	MarkMessagesAsRead(ctx context.Context, req *MarkMessagesAsReadRequest) (int, error)
	GetUnreadCount(ctx context.Context, userID int64, chatID *int64) (int, error)

	// Attachment operations
	UploadAttachment(ctx context.Context, req *UploadAttachmentRequest) (*domain.ChatAttachment, error)
	GetAttachment(ctx context.Context, attachmentID, userID int64) (*domain.ChatAttachment, error)
	DeleteAttachment(ctx context.Context, attachmentID, userID int64) error

	// WebSocket hub integration
	SetHub(hub ChatHub)

	// Real-time streaming (handled at transport layer, service provides messages)
	// StreamMessages will be implemented in gRPC handler using GetMessages polling
}

// CreateChatRequest contains parameters for creating a chat
type CreateChatRequest struct {
	BuyerID             int64  // User who initiates chat (authenticated)
	SellerID            int64  // User who receives chat
	ListingID           *int64 // Optional: marketplace listing context
	StorefrontProductID *int64 // Optional: B2C product context
}

// GetOrCreateChatRequest contains parameters for getting or creating a chat
type GetOrCreateChatRequest struct {
	UserID              int64  // Authenticated user
	OtherUserID         *int64 // Direct message to another user
	ListingID           *int64 // Optional: marketplace listing context
	StorefrontProductID *int64 // Optional: B2C product context
}

// GetUserChatsRequest contains parameters for listing user chats
type GetUserChatsRequest struct {
	UserID    int64  // Authenticated user
	Archived  bool   // Show archived chats
	ListingID *int64 // Filter by listing
	Limit     int    // Page size (default: 20, max: 100)
	Offset    int    // Page offset
}

// SendMessageRequest contains parameters for sending a message
type SendMessageRequest struct {
	ChatID           int64   // Chat to send message in
	SenderID         int64   // Authenticated user (must be participant)
	Content          string  // Message text (1-10000 chars)
	OriginalLanguage string  // ISO 639-1 code (en, ru, sr)
	AttachmentIDs    []int64 // Pre-uploaded attachments
}

// SendSystemMessageRequest contains parameters for sending a system message
// System messages are sent from the marketplace to notify users about events
type SendSystemMessageRequest struct {
	ReceiverID       int64  // User ID to send the message to
	Content          string // Message text
	OriginalLanguage string // ISO 639-1 code (en, ru, sr)
}

// GetMessagesRequest contains parameters for retrieving messages
type GetMessagesRequest struct {
	ChatID          int64  // Chat to retrieve messages from
	UserID          int64  // Authenticated user (must be participant)
	BeforeMessageID *int64 // Cursor: get messages before this ID (older messages)
	AfterMessageID  *int64 // Cursor: get messages after this ID (newer messages)
	Limit           int    // Max items (default: 50, max: 100)
}

// MarkMessagesAsReadRequest contains parameters for marking messages as read
type MarkMessagesAsReadRequest struct {
	ChatID     int64   // Chat containing messages
	UserID     int64   // Authenticated user (must be receiver)
	MessageIDs []int64 // Specific messages to mark (empty = mark all unread)
	MarkAll    bool    // If true, mark all unread messages in chat
}

// UploadAttachmentRequest contains parameters for uploading an attachment
type UploadAttachmentRequest struct {
	UserID      int64                 // Authenticated user
	FileName    string                // Original filename
	ContentType string                // MIME type
	FileData    []byte                // File content
	FileType    domain.AttachmentType // image, video, document
}

// chatService implements ChatService
type chatService struct {
	chatRepo       ChatRepository
	messageRepo    MessageRepository
	attachmentRepo AttachmentRepository
	productsRepo   *postgres.Repository // For validating listing/product exists
	authService    *authservice.AuthService
	pool           *pgxpool.Pool
	hub            ChatHub // WebSocket hub for real-time updates
	logger         zerolog.Logger
}

// ChatHub defines the interface for WebSocket broadcasting
// This allows chatService to broadcast events without tight coupling to websocket package
type ChatHub interface {
	BroadcastNewMessage(chatID int64, message *domain.Message)
	BroadcastMessageRead(chatID, messageID, userID int64)
	BroadcastTyping(chatID, userID int64, isTyping bool)
}

// NewChatService creates a new chat service
func NewChatService(
	chatRepo ChatRepository,
	messageRepo MessageRepository,
	attachmentRepo AttachmentRepository,
	productsRepo *postgres.Repository,
	authService *authservice.AuthService,
	pool *pgxpool.Pool,
	logger zerolog.Logger,
) ChatService {
	return &chatService{
		chatRepo:       chatRepo,
		messageRepo:    messageRepo,
		attachmentRepo: attachmentRepo,
		productsRepo:   productsRepo,
		authService:    authService,
		pool:           pool,
		hub:            nil, // Will be set via SetHub if WebSocket is enabled
		logger:         logger.With().Str("component", "chat_service").Logger(),
	}
}

// SetHub sets the WebSocket hub for real-time broadcasting
func (s *chatService) SetHub(hub ChatHub) {
	s.hub = hub
	s.logger.Info().Msg("WebSocket hub connected to chat service")
}

// CreateChat creates a new chat between buyer and seller
func (s *chatService) CreateChat(ctx context.Context, req *CreateChatRequest) (*domain.Chat, error) {
	s.logger.Debug().
		Int64("buyer_id", req.BuyerID).
		Int64("seller_id", req.SellerID).
		Interface("listing_id", req.ListingID).
		Interface("storefront_product_id", req.StorefrontProductID).
		Msg("creating chat")

	// Validate input
	if err := s.validateCreateChatRequest(req); err != nil {
		return nil, err
	}

	// Prevent user from chatting with themselves
	if req.BuyerID == req.SellerID {
		return nil, ErrChatWithSelf
	}

	// Validate listing/product exists and is active (if provided)
	if req.ListingID != nil {
		// Try to get as C2C listing first
		listing, err := s.productsRepo.GetListingByID(ctx, *req.ListingID)
		if err == nil {
			// C2C listing - verify seller matches UserID
			if listing.UserID != req.SellerID {
				s.logger.Warn().
					Int64("expected_seller_id", req.SellerID).
					Int64("actual_user_id", listing.UserID).
					Msg("seller_id mismatch for C2C listing")
				return nil, ErrUnauthorized
			}
		} else {
			// Not found or error - this is expected if it's a B2C product
			// For B2C products, seller verification happens via StorefrontID
			s.logger.Debug().Err(err).Int64("listing_id", *req.ListingID).Msg("listing not found as C2C, may be B2C product")
		}
	}

	// Validate other user exists (via Auth Service)
	if err := s.validateUserExists(ctx, req.SellerID); err != nil {
		return nil, err
	}

	// Create chat
	chat := &domain.Chat{
		BuyerID:             req.BuyerID,
		SellerID:            req.SellerID,
		ListingID:           req.ListingID,
		StorefrontProductID: req.StorefrontProductID,
		Status:              domain.ChatStatusActive,
		LastMessageAt:       time.Now(),
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	if err := s.chatRepo.Create(ctx, chat); err != nil {
		s.logger.Error().Err(err).Msg("failed to create chat")
		return nil, fmt.Errorf("failed to create chat: %w", err)
	}

	s.logger.Info().Int64("chat_id", chat.ID).Msg("chat created successfully")
	return chat, nil
}

// GetOrCreateChat retrieves existing chat or creates a new one
func (s *chatService) GetOrCreateChat(ctx context.Context, req *GetOrCreateChatRequest) (*domain.Chat, bool, error) {
	s.logger.Debug().
		Int64("user_id", req.UserID).
		Interface("other_user_id", req.OtherUserID).
		Interface("listing_id", req.ListingID).
		Interface("storefront_product_id", req.StorefrontProductID).
		Msg("getting or creating chat")

	// Validate input
	if err := s.validateGetOrCreateChatRequest(req); err != nil {
		return nil, false, err
	}

	var existingChat *domain.Chat
	var err error

	// Try to find existing chat based on context
	if req.ListingID != nil {
		// Try to get as C2C listing first
		listing, listingErr := s.productsRepo.GetListingByID(ctx, *req.ListingID)

		var buyerID, sellerID int64

		if listingErr == nil {
			// C2C listing found - UserID is the seller
			if listing.UserID == req.UserID {
				// Current user is seller, other user must be buyer
				if req.OtherUserID == nil {
					return nil, false, fmt.Errorf("%w: other_user_id required when seller initiates chat", ErrInvalidInput)
				}
				sellerID = req.UserID
				buyerID = *req.OtherUserID
			} else {
				// Current user is buyer, listing owner is seller
				buyerID = req.UserID
				sellerID = listing.UserID
			}
		} else {
			// Not a C2C listing - might be B2C product or error
			// For B2C, we need other_user_id to be the seller
			if req.OtherUserID == nil {
				return nil, false, fmt.Errorf("%w: other_user_id required for B2C product chats", ErrInvalidInput)
			}
			buyerID = req.UserID
			sellerID = *req.OtherUserID
		}

		existingChat, err = s.chatRepo.GetByParticipantsAndListing(ctx, buyerID, sellerID, *req.ListingID)
	} else if req.StorefrontProductID != nil {
		// Similar logic for storefront products
		// For now, assume req.OtherUserID is the seller
		if req.OtherUserID == nil {
			return nil, false, fmt.Errorf("%w: other_user_id required for storefront product", ErrInvalidInput)
		}
		existingChat, err = s.chatRepo.GetByParticipantsAndProduct(ctx, req.UserID, *req.OtherUserID, *req.StorefrontProductID)
	} else if req.OtherUserID != nil {
		// Direct message without item context
		existingChat, err = s.chatRepo.GetByParticipantsDirect(ctx, req.UserID, *req.OtherUserID)
	}

	// Check for chat not found error - repository returns fmt.Errorf("chat not found")
	// which doesn't match ErrChatNotFound, so we need to check the error message
	if err != nil && !isChatNotFoundError(err) {
		s.logger.Error().Err(err).Msg("failed to get existing chat")
		return nil, false, fmt.Errorf("failed to get existing chat: %w", err)
	}

	// If chat exists, return it
	if existingChat != nil {
		s.logger.Debug().Int64("chat_id", existingChat.ID).Msg("existing chat found")
		return existingChat, false, nil
	}

	// Create new chat
	createReq := &CreateChatRequest{
		BuyerID:             req.UserID,
		ListingID:           req.ListingID,
		StorefrontProductID: req.StorefrontProductID,
	}

	if req.OtherUserID != nil {
		createReq.SellerID = *req.OtherUserID
	} else if req.ListingID != nil {
		// Get seller from C2C listing
		listing, _ := s.productsRepo.GetListingByID(ctx, *req.ListingID)
		if listing != nil {
			createReq.SellerID = listing.UserID
		}
	}

	newChat, err := s.CreateChat(ctx, createReq)
	if err != nil {
		return nil, false, err
	}

	s.logger.Info().Int64("chat_id", newChat.ID).Msg("new chat created")
	return newChat, true, nil
}

// GetChat retrieves a chat by ID with authorization check
func (s *chatService) GetChat(ctx context.Context, chatID, userID int64) (*domain.Chat, error) {
	s.logger.Debug().Int64("chat_id", chatID).Int64("user_id", userID).Msg("getting chat")

	// Get chat
	chat, err := s.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		if err == ErrChatNotFound {
			return nil, ErrChatNotFound
		}
		s.logger.Error().Err(err).Int64("chat_id", chatID).Msg("failed to get chat")
		return nil, fmt.Errorf("failed to get chat: %w", err)
	}

	// Authorization: user must be participant
	if !chat.IsParticipant(userID) {
		s.logger.Warn().
			Int64("chat_id", chatID).
			Int64("user_id", userID).
			Msg("user not authorized to access chat")
		return nil, ErrNotParticipant
	}

	// Load unread count for this user
	unreadCount, err := s.messageRepo.GetUnreadCount(ctx, chatID, userID)
	if err != nil {
		s.logger.Warn().Err(err).Msg("failed to count unread messages")
	} else {
		chat.UnreadCount = unreadCount
	}

	return chat, nil
}

// GetUserChats retrieves all chats for a user with filtering
func (s *chatService) GetUserChats(ctx context.Context, req *GetUserChatsRequest) ([]*domain.Chat, int, error) {
	s.logger.Debug().
		Int64("user_id", req.UserID).
		Bool("archived", req.Archived).
		Int("limit", req.Limit).
		Int("offset", req.Offset).
		Msg("getting user chats")

	// Validate pagination
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	// Get chats
	chats, total, err := s.chatRepo.GetUserChats(ctx, req.UserID, nil, req.Archived, req.Limit, req.Offset)
	if err != nil {
		s.logger.Error().Err(err).Int64("user_id", req.UserID).Msg("failed to list user chats")
		return nil, 0, fmt.Errorf("failed to list user chats: %w", err)
	}

	// Load unread counts, last message and enrich metadata for each chat
	for _, chat := range chats {
		// Load unread count
		unreadCount, err := s.messageRepo.GetUnreadCount(ctx, chat.ID, req.UserID)
		if err != nil {
			s.logger.Warn().Err(err).Int64("chat_id", chat.ID).Msg("failed to count unread messages")
		} else {
			chat.UnreadCount = unreadCount
		}

		// Load last message (needed for notifications and chat preview)
		lastMessage, err := s.messageRepo.GetLatestMessage(ctx, chat.ID)
		if err != nil {
			s.logger.Warn().Err(err).Int64("chat_id", chat.ID).Msg("failed to get latest message")
		} else {
			chat.LastMessage = lastMessage
		}

		// Enrich with listing metadata
		s.enrichChatWithListingMetadata(ctx, chat)

		// Enrich with user names
		s.enrichChatWithUserNames(ctx, chat)
	}

	return chats, total, nil
}

// ArchiveChat archives or unarchives a chat
func (s *chatService) ArchiveChat(ctx context.Context, chatID, userID int64, archived bool) error {
	s.logger.Debug().
		Int64("chat_id", chatID).
		Int64("user_id", userID).
		Bool("archived", archived).
		Msg("archiving chat")

	// Get chat to verify ownership
	chat, err := s.GetChat(ctx, chatID, userID)
	if err != nil {
		return err
	}

	// Check if can archive
	if archived && !chat.CanArchive() {
		return fmt.Errorf("%w: chat cannot be archived in current state", ErrInvalidInput)
	}

	// Archive chat
	if err := s.chatRepo.ArchiveChat(ctx, chatID, archived); err != nil {
		s.logger.Error().Err(err).Int64("chat_id", chatID).Msg("failed to archive chat")
		return fmt.Errorf("failed to archive chat: %w", err)
	}

	s.logger.Info().Int64("chat_id", chatID).Bool("archived", archived).Msg("chat archived successfully")
	return nil
}

// DeleteChat deletes a chat (admin only)
func (s *chatService) DeleteChat(ctx context.Context, chatID int64) error {
	s.logger.Debug().Int64("chat_id", chatID).Msg("deleting chat")

	// Verify chat exists
	_, err := s.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		if err == ErrChatNotFound {
			return ErrChatNotFound
		}
		return fmt.Errorf("failed to get chat: %w", err)
	}

	// Delete chat (cascade deletes messages and attachments)
	if err := s.chatRepo.Delete(ctx, chatID); err != nil {
		s.logger.Error().Err(err).Int64("chat_id", chatID).Msg("failed to delete chat")
		return fmt.Errorf("failed to delete chat: %w", err)
	}

	s.logger.Info().Int64("chat_id", chatID).Msg("chat deleted successfully")
	return nil
}

// SendMessage sends a new message in a chat
func (s *chatService) SendMessage(ctx context.Context, req *SendMessageRequest) (*domain.Message, error) {
	s.logger.Debug().
		Int64("chat_id", req.ChatID).
		Int64("sender_id", req.SenderID).
		Int("content_length", len(req.Content)).
		Msg("sending message")

	// Validate input
	if err := s.validateSendMessageRequest(req); err != nil {
		return nil, err
	}

	// Get chat and verify sender is participant
	chat, err := s.GetChat(ctx, req.ChatID, req.SenderID)
	if err != nil {
		return nil, err
	}

	// Check if chat is blocked
	if chat.Status == domain.ChatStatusBlocked {
		return nil, ErrChatBlocked
	}

	// Determine receiver
	receiverID := chat.GetOtherParticipantID(req.SenderID)

	// Create message
	message := &domain.Message{
		ChatID:              req.ChatID,
		SenderID:            req.SenderID,
		ReceiverID:          receiverID,
		Content:             strings.TrimSpace(req.Content),
		OriginalLanguage:    req.OriginalLanguage,
		ListingID:           chat.ListingID,
		StorefrontProductID: chat.StorefrontProductID,
		Status:              domain.MessageStatusSent,
		IsRead:              false,
		HasAttachments:      len(req.AttachmentIDs) > 0,
		AttachmentsCount:    int32(len(req.AttachmentIDs)),
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	// Start transaction for message creation + chat update
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to begin transaction")
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Create message
	if err := s.messageRepo.Create(ctx, message); err != nil {
		s.logger.Error().Err(err).Msg("failed to create message")
		return nil, fmt.Errorf("failed to create message: %w", err)
	}

	// TODO: If attachment_ids provided, validate they exist and belong to sender
	// For now, we skip this as attachments should be uploaded separately with message_id set

	// Update chat's last_message_at
	chat.LastMessageAt = message.CreatedAt
	if err := s.chatRepo.Update(ctx, chat); err != nil {
		s.logger.Warn().Err(err).Msg("failed to update chat last_message_at")
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		s.logger.Error().Err(err).Msg("failed to commit transaction")
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.logger.Info().
		Int64("message_id", message.ID).
		Int64("chat_id", req.ChatID).
		Msg("message sent successfully")

	// Broadcast through WebSocket if hub is available
	if s.hub != nil {
		s.hub.BroadcastNewMessage(req.ChatID, message)
		s.logger.Debug().
			Int64("message_id", message.ID).
			Int64("chat_id", req.ChatID).
			Msg("message broadcasted via WebSocket")
	}

	return message, nil
}

// SendSystemMessage sends a system message from the marketplace to a user
// It creates a direct chat between the system user and the receiver if it doesn't exist
func (s *chatService) SendSystemMessage(ctx context.Context, req *SendSystemMessageRequest) (*domain.Message, error) {
	s.logger.Debug().
		Int64("receiver_id", req.ReceiverID).
		Int("content_length", len(req.Content)).
		Msg("sending system message")

	// Validate input
	if req.ReceiverID <= 0 {
		return nil, fmt.Errorf("%w: receiver_id must be greater than 0", ErrInvalidInput)
	}
	content := strings.TrimSpace(req.Content)
	if len(content) == 0 {
		return nil, ErrMessageEmpty
	}
	if len(content) > MaxMessageLength {
		return nil, &ErrMessageTooLong{Length: len(content), MaxLength: MaxMessageLength}
	}

	// Default language
	if req.OriginalLanguage == "" {
		req.OriginalLanguage = "en"
	}

	// Get or create direct chat between system user and receiver
	chat, _, err := s.GetOrCreateChat(ctx, &GetOrCreateChatRequest{
		UserID:      domain.SystemUserID,
		OtherUserID: &req.ReceiverID,
	})
	if err != nil {
		s.logger.Error().Err(err).Int64("receiver_id", req.ReceiverID).Msg("failed to get or create system chat")
		return nil, fmt.Errorf("failed to get or create system chat: %w", err)
	}

	// Create system message
	message := &domain.Message{
		ChatID:           chat.ID,
		SenderID:         domain.SystemUserID,
		ReceiverID:       req.ReceiverID,
		IsSystem:         true,
		Content:          content,
		OriginalLanguage: req.OriginalLanguage,
		Status:           domain.MessageStatusSent,
		IsRead:           false,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Start transaction
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to begin transaction")
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Create message
	if err := s.messageRepo.Create(ctx, message); err != nil {
		s.logger.Error().Err(err).Msg("failed to create system message")
		return nil, fmt.Errorf("failed to create system message: %w", err)
	}

	// Update chat's last_message_at
	chat.LastMessageAt = message.CreatedAt
	if err := s.chatRepo.Update(ctx, chat); err != nil {
		s.logger.Warn().Err(err).Msg("failed to update chat last_message_at")
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		s.logger.Error().Err(err).Msg("failed to commit transaction")
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Set sender name for UI
	senderName := "Vondi Marketplace"
	message.SenderName = &senderName

	s.logger.Info().
		Int64("message_id", message.ID).
		Int64("chat_id", chat.ID).
		Int64("receiver_id", req.ReceiverID).
		Msg("system message sent successfully")

	// Broadcast through WebSocket if hub is available
	if s.hub != nil {
		s.hub.BroadcastNewMessage(chat.ID, message)
		s.logger.Debug().
			Int64("message_id", message.ID).
			Int64("chat_id", chat.ID).
			Msg("system message broadcasted via WebSocket")
	}

	return message, nil
}

// GetMessages retrieves messages with cursor-based pagination
func (s *chatService) GetMessages(ctx context.Context, req *GetMessagesRequest) ([]*domain.Message, bool, error) {
	s.logger.Debug().
		Int64("chat_id", req.ChatID).
		Int64("user_id", req.UserID).
		Interface("before_message_id", req.BeforeMessageID).
		Interface("after_message_id", req.AfterMessageID).
		Int("limit", req.Limit).
		Msg("getting messages")

	// Validate input
	if err := s.validateGetMessagesRequest(req); err != nil {
		return nil, false, err
	}

	// Verify user is participant
	_, err := s.GetChat(ctx, req.ChatID, req.UserID)
	if err != nil {
		return nil, false, err
	}

	// Set default limit
	if req.Limit <= 0 {
		req.Limit = 50
	}
	if req.Limit > 100 {
		req.Limit = 100
	}

	// Get messages with pagination
	var messages []*domain.Message
	var hasMore bool

	if req.BeforeMessageID != nil {
		// Use cursor-based pagination for "load older" scenario
		var err error
		messages, hasMore, err = s.messageRepo.GetMessagesByCursor(ctx, req.ChatID, req.BeforeMessageID, req.Limit)
		if err != nil {
			s.logger.Error().Err(err).Int64("chat_id", req.ChatID).Msg("failed to list messages")
			return nil, false, fmt.Errorf("failed to list messages: %w", err)
		}
	} else {
		// Use regular query for initial load
		var err error
		messages, err = s.messageRepo.GetMessages(ctx, req.ChatID, req.BeforeMessageID, req.AfterMessageID, req.Limit)
		if err != nil {
			s.logger.Error().Err(err).Int64("chat_id", req.ChatID).Msg("failed to list messages")
			return nil, false, fmt.Errorf("failed to list messages: %w", err)
		}
		hasMore = len(messages) == req.Limit
	}

	return messages, hasMore, nil
}

// MarkMessagesAsRead marks messages as read
func (s *chatService) MarkMessagesAsRead(ctx context.Context, req *MarkMessagesAsReadRequest) (int, error) {
	s.logger.Debug().
		Int64("chat_id", req.ChatID).
		Int64("user_id", req.UserID).
		Bool("mark_all", req.MarkAll).
		Int("message_ids_count", len(req.MessageIDs)).
		Msg("marking messages as read")

	// Verify user is participant
	_, err := s.GetChat(ctx, req.ChatID, req.UserID)
	if err != nil {
		return 0, err
	}

	var markedCount int

	if req.MarkAll {
		// Mark all unread messages in chat
		markedCount, err = s.messageRepo.MarkAllAsRead(ctx, req.ChatID, req.UserID)
		if err != nil {
			s.logger.Error().Err(err).Int64("chat_id", req.ChatID).Msg("failed to mark all messages as read")
			return 0, fmt.Errorf("failed to mark all messages as read: %w", err)
		}
	} else if len(req.MessageIDs) > 0 {
		// Mark specific messages
		markedCount, err = s.messageRepo.MarkMessagesAsRead(ctx, req.ChatID, req.UserID, req.MessageIDs)
		if err != nil {
			s.logger.Error().Err(err).Msg("failed to mark messages as read")
			return 0, fmt.Errorf("failed to mark messages as read: %w", err)
		}
	}

	s.logger.Info().
		Int64("chat_id", req.ChatID).
		Int("marked_count", markedCount).
		Msg("messages marked as read")

	// Broadcast message read events via WebSocket if hub is available
	if s.hub != nil && markedCount > 0 {
		// If specific messages were marked, broadcast for each
		if len(req.MessageIDs) > 0 {
			for _, msgID := range req.MessageIDs {
				s.hub.BroadcastMessageRead(req.ChatID, msgID, req.UserID)
			}
		} else if req.MarkAll {
			// For mark_all, we broadcast a single event with chat_id
			// Frontend should handle this by marking all messages in that chat as read
			// We'll use messageID=0 as a special marker for "all messages"
			s.hub.BroadcastMessageRead(req.ChatID, 0, req.UserID)
		}
		s.logger.Debug().
			Int64("chat_id", req.ChatID).
			Int("count", markedCount).
			Msg("message read events broadcasted via WebSocket")
	}

	return markedCount, nil
}

// GetUnreadCount retrieves unread message count
func (s *chatService) GetUnreadCount(ctx context.Context, userID int64, chatID *int64) (int, error) {
	s.logger.Debug().
		Int64("user_id", userID).
		Interface("chat_id", chatID).
		Msg("getting unread count")

	var count32 int32
	var err error

	if chatID != nil {
		// Count for specific chat
		count32, err = s.messageRepo.GetUnreadCount(ctx, *chatID, userID)
	} else {
		// Count across all chats
		count32, err = s.messageRepo.GetUnreadCountByUser(ctx, userID)
	}

	if err != nil {
		s.logger.Error().Err(err).Int64("user_id", userID).Msg("failed to count unread messages")
		return 0, fmt.Errorf("failed to count unread messages: %w", err)
	}

	return int(count32), nil
}

// UploadAttachment uploads a file attachment
func (s *chatService) UploadAttachment(ctx context.Context, req *UploadAttachmentRequest) (*domain.ChatAttachment, error) {
	s.logger.Debug().
		Int64("user_id", req.UserID).
		Str("file_name", req.FileName).
		Str("content_type", req.ContentType).
		Int("file_size", len(req.FileData)).
		Msg("uploading attachment")

	// Validate input
	if err := s.validateUploadAttachmentRequest(req); err != nil {
		return nil, err
	}

	// Validate file size based on type
	if err := s.validateFileSize(req.FileType, int64(len(req.FileData))); err != nil {
		return nil, err
	}

	// TODO: Upload to MinIO (for now, create stub attachment)
	attachment := &domain.ChatAttachment{
		FileType:      req.FileType,
		FileName:      req.FileName,
		FileSize:      int64(len(req.FileData)),
		ContentType:   req.ContentType,
		StorageType:   "minio",
		StorageBucket: "chat-attachments",
		FilePath:      fmt.Sprintf("attachments/%d/%s", req.UserID, req.FileName),
		PublicURL:     fmt.Sprintf("https://cdn.example.com/attachments/%d/%s", req.UserID, req.FileName),
		CreatedAt:     time.Now(),
	}

	// Create attachment record
	if err := s.attachmentRepo.Create(ctx, attachment); err != nil {
		s.logger.Error().Err(err).Msg("failed to create attachment")
		return nil, fmt.Errorf("failed to create attachment: %w", err)
	}

	s.logger.Info().Int64("attachment_id", attachment.ID).Msg("attachment uploaded successfully")
	return attachment, nil
}

// GetAttachment retrieves an attachment with authorization check
func (s *chatService) GetAttachment(ctx context.Context, attachmentID, userID int64) (*domain.ChatAttachment, error) {
	s.logger.Debug().
		Int64("attachment_id", attachmentID).
		Int64("user_id", userID).
		Msg("getting attachment")

	// Get attachment
	attachment, err := s.attachmentRepo.GetByID(ctx, attachmentID)
	if err != nil {
		if err == ErrAttachmentNotFound {
			return nil, ErrAttachmentNotFound
		}
		s.logger.Error().Err(err).Int64("attachment_id", attachmentID).Msg("failed to get attachment")
		return nil, fmt.Errorf("failed to get attachment: %w", err)
	}

	// Get parent message to verify access
	message, err := s.messageRepo.GetByID(ctx, attachment.MessageID)
	if err != nil {
		return nil, fmt.Errorf("failed to get parent message: %w", err)
	}

	// Verify user has access to the chat
	_, err = s.GetChat(ctx, message.ChatID, userID)
	if err != nil {
		return nil, err
	}

	return attachment, nil
}

// DeleteAttachment deletes an attachment
func (s *chatService) DeleteAttachment(ctx context.Context, attachmentID, userID int64) error {
	s.logger.Debug().
		Int64("attachment_id", attachmentID).
		Int64("user_id", userID).
		Msg("deleting attachment")

	// Get attachment
	attachment, err := s.attachmentRepo.GetByID(ctx, attachmentID)
	if err != nil {
		if err == ErrAttachmentNotFound {
			return ErrAttachmentNotFound
		}
		return fmt.Errorf("failed to get attachment: %w", err)
	}

	// Get parent message to verify sender
	message, err := s.messageRepo.GetByID(ctx, attachment.MessageID)
	if err != nil {
		return fmt.Errorf("failed to get parent message: %w", err)
	}

	// Verify user is sender
	if message.SenderID != userID {
		return ErrUnauthorized
	}

	// TODO: Delete from MinIO storage

	// Delete attachment record
	if err := s.attachmentRepo.Delete(ctx, attachmentID); err != nil {
		s.logger.Error().Err(err).Int64("attachment_id", attachmentID).Msg("failed to delete attachment")
		return fmt.Errorf("failed to delete attachment: %w", err)
	}

	s.logger.Info().Int64("attachment_id", attachmentID).Msg("attachment deleted successfully")
	return nil
}

// Validation helpers

func (s *chatService) validateCreateChatRequest(req *CreateChatRequest) error {
	if req.BuyerID <= 0 {
		return fmt.Errorf("%w: buyer_id must be greater than 0", ErrInvalidInput)
	}
	if req.SellerID <= 0 {
		return fmt.Errorf("%w: seller_id must be greater than 0", ErrInvalidInput)
	}

	// Must provide exactly one context (listing, product, or direct message)
	contextCount := 0
	if req.ListingID != nil {
		contextCount++
	}
	if req.StorefrontProductID != nil {
		contextCount++
	}
	if contextCount > 1 {
		return fmt.Errorf("%w: can only provide one of: listing_id, storefront_product_id", ErrInvalidInput)
	}

	return nil
}

func (s *chatService) validateGetOrCreateChatRequest(req *GetOrCreateChatRequest) error {
	if req.UserID <= 0 {
		return fmt.Errorf("%w: user_id must be greater than 0", ErrInvalidInput)
	}

	// Must provide at least one context
	if req.ListingID == nil && req.StorefrontProductID == nil && req.OtherUserID == nil {
		return fmt.Errorf("%w: must provide listing_id, storefront_product_id, or other_user_id", ErrInvalidInput)
	}

	// Cannot provide both listing and product
	if req.ListingID != nil && req.StorefrontProductID != nil {
		return fmt.Errorf("%w: cannot provide both listing_id and storefront_product_id", ErrInvalidInput)
	}

	return nil
}

func (s *chatService) validateSendMessageRequest(req *SendMessageRequest) error {
	if req.ChatID <= 0 {
		return fmt.Errorf("%w: chat_id must be greater than 0", ErrInvalidInput)
	}
	if req.SenderID <= 0 {
		return fmt.Errorf("%w: sender_id must be greater than 0", ErrInvalidInput)
	}

	content := strings.TrimSpace(req.Content)
	if len(content) == 0 {
		return ErrMessageEmpty
	}
	if len(content) > MaxMessageLength {
		return &ErrMessageTooLong{
			Length:    len(content),
			MaxLength: MaxMessageLength,
		}
	}

	// Validate language code
	if req.OriginalLanguage == "" {
		req.OriginalLanguage = "en"
	}

	return nil
}

func (s *chatService) validateGetMessagesRequest(req *GetMessagesRequest) error {
	if req.ChatID <= 0 {
		return fmt.Errorf("%w: chat_id must be greater than 0", ErrInvalidInput)
	}
	if req.UserID <= 0 {
		return fmt.Errorf("%w: user_id must be greater than 0", ErrInvalidInput)
	}

	// Cannot provide both before and after cursors
	if req.BeforeMessageID != nil && req.AfterMessageID != nil {
		return fmt.Errorf("%w: cannot provide both before_message_id and after_message_id", ErrInvalidInput)
	}

	return nil
}

func (s *chatService) validateUploadAttachmentRequest(req *UploadAttachmentRequest) error {
	if req.UserID <= 0 {
		return fmt.Errorf("%w: user_id must be greater than 0", ErrInvalidInput)
	}
	if req.FileName == "" {
		return fmt.Errorf("%w: file_name is required", ErrInvalidInput)
	}
	if len(req.FileData) == 0 {
		return fmt.Errorf("%w: file_data is required", ErrInvalidInput)
	}

	// Validate content type
	if !s.isValidContentType(req.FileType, req.ContentType) {
		return &ErrInvalidFileType{ContentType: req.ContentType}
	}

	return nil
}

func (s *chatService) validateFileSize(fileType domain.AttachmentType, size int64) error {
	var maxSize int64
	var fileTypeStr AttachmentFileType

	switch fileType {
	case domain.AttachmentTypeImage:
		maxSize = MaxImageSize
		fileTypeStr = AttachmentFileTypeImage
	case domain.AttachmentTypeVideo:
		maxSize = MaxVideoSize
		fileTypeStr = AttachmentFileTypeVideo
	case domain.AttachmentTypeDocument:
		maxSize = MaxDocumentSize
		fileTypeStr = AttachmentFileTypeDocument
	default:
		return fmt.Errorf("%w: unknown file type", ErrInvalidInput)
	}

	if size > maxSize {
		return &ErrAttachmentTooLarge{
			FileType: fileTypeStr,
			Size:     size,
			MaxSize:  maxSize,
		}
	}

	return nil
}

func (s *chatService) isValidContentType(fileType domain.AttachmentType, contentType string) bool {
	validTypes := map[domain.AttachmentType][]string{
		domain.AttachmentTypeImage: {
			"image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp",
		},
		domain.AttachmentTypeVideo: {
			"video/mp4", "video/webm", "video/quicktime",
		},
		domain.AttachmentTypeDocument: {
			"application/pdf", "application/msword",
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			"text/plain",
		},
	}

	allowed, exists := validTypes[fileType]
	if !exists {
		return false
	}

	for _, valid := range allowed {
		if contentType == valid {
			return true
		}
	}

	return false
}

func (s *chatService) validateUserExists(ctx context.Context, userID int64) error {
	// TODO: Call Auth Service to verify user exists
	// For now, assume user exists
	s.logger.Debug().Int64("user_id", userID).Msg("validating user exists (stub)")
	return nil
}

// enrichChatWithListingMetadata enriches chat with listing title and image URL
func (s *chatService) enrichChatWithListingMetadata(ctx context.Context, chat *domain.Chat) {
	if chat.ListingID == nil || *chat.ListingID <= 0 {
		return
	}

	// Get listing data from listings table
	listing, err := s.productsRepo.GetListingByID(ctx, int64(*chat.ListingID))
	if err != nil {
		s.logger.Debug().Err(err).Int64("listing_id", *chat.ListingID).Msg("failed to get listing for chat enrichment")
		return
	}

	// Set listing title
	if listing.Title != "" {
		chat.ListingTitle = &listing.Title
	}

	// Set listing owner ID
	if listing.UserID > 0 {
		chat.ListingOwnerID = &listing.UserID
	}

	// Get first image URL if available
	if listing.ID > 0 {
		var imageURL string
		query := `
			SELECT url
			FROM listing_images
			WHERE listing_id = $1
			ORDER BY display_order ASC
			LIMIT 1
		`
		err := s.pool.QueryRow(ctx, query, listing.ID).Scan(&imageURL)
		if err == nil && imageURL != "" {
			chat.ListingImageURL = &imageURL
		}
	}
}

// enrichChatWithUserNames enriches chat with buyer and seller names from Auth Service
func (s *chatService) enrichChatWithUserNames(ctx context.Context, chat *domain.Chat) {
	if s.authService == nil {
		return
	}

	// Get buyer name
	if chat.BuyerID > 0 {
		buyerUser, err := s.authService.GetUser(ctx, int(chat.BuyerID))
		if err == nil && buyerUser != nil && buyerUser.Name != "" {
			chat.BuyerName = &buyerUser.Name
		}
	}

	// Get seller name
	if chat.SellerID > 0 {
		sellerUser, err := s.authService.GetUser(ctx, int(chat.SellerID))
		if err == nil && sellerUser != nil && sellerUser.Name != "" {
			chat.SellerName = &sellerUser.Name
		}
	}
}
