// Package service implements business logic for the listings microservice.
package service

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/sveturs/listings/internal/domain"
)

// NOTE: Chat service tests are limited because the service depends on *postgres.Repository (concrete type)
// rather than an interface for products validation. Full testing would require:
// 1. Refactoring productsRepo to use an interface
// 2. OR using reflection/unsafe to inject mocks
// 3. OR integration tests with real database
//
// For now, we test:
// - Methods that don't require listing validation
// - Repository interactions
// - Business logic that can be isolated

// =============================================================================
// MOCK REPOSITORIES
// =============================================================================

// MockChatRepository is a mock for ChatRepository
type MockChatRepository struct {
	mock.Mock
}

func (m *MockChatRepository) Create(ctx context.Context, chat *domain.Chat) error {
	args := m.Called(ctx, chat)
	if args.Error(0) == nil {
		// Simulate DB setting ID and timestamps
		chat.ID = 1
		chat.CreatedAt = time.Now()
		chat.UpdatedAt = time.Now()
	}
	return args.Error(0)
}

func (m *MockChatRepository) GetByID(ctx context.Context, chatID int64) (*domain.Chat, error) {
	args := m.Called(ctx, chatID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Chat), args.Error(1)
}

func (m *MockChatRepository) GetByParticipantsAndListing(ctx context.Context, buyerID, sellerID, listingID int64) (*domain.Chat, error) {
	args := m.Called(ctx, buyerID, sellerID, listingID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Chat), args.Error(1)
}

func (m *MockChatRepository) GetByParticipantsAndProduct(ctx context.Context, buyerID, sellerID, productID int64) (*domain.Chat, error) {
	args := m.Called(ctx, buyerID, sellerID, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Chat), args.Error(1)
}

func (m *MockChatRepository) GetByParticipantsDirect(ctx context.Context, user1ID, user2ID int64) (*domain.Chat, error) {
	args := m.Called(ctx, user1ID, user2ID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Chat), args.Error(1)
}

func (m *MockChatRepository) Update(ctx context.Context, chat *domain.Chat) error {
	args := m.Called(ctx, chat)
	return args.Error(0)
}

func (m *MockChatRepository) Delete(ctx context.Context, chatID int64) error {
	args := m.Called(ctx, chatID)
	return args.Error(0)
}

func (m *MockChatRepository) GetUserChats(ctx context.Context, userID int64, status *domain.ChatStatus, archivedOnly bool, limit, offset int) ([]*domain.Chat, int, error) {
	args := m.Called(ctx, userID, status, archivedOnly, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int), args.Error(2)
	}
	return args.Get(0).([]*domain.Chat), args.Get(1).(int), args.Error(2)
}

func (m *MockChatRepository) GetUnreadCount(ctx context.Context, userID int64, chatID *int64) (int32, error) {
	args := m.Called(ctx, userID, chatID)
	return args.Get(0).(int32), args.Error(1)
}

func (m *MockChatRepository) ArchiveChat(ctx context.Context, chatID int64, archived bool) error {
	args := m.Called(ctx, chatID, archived)
	return args.Error(0)
}

// MockMessageRepository is a mock for MessageRepository
type MockMessageRepository struct {
	mock.Mock
}

func (m *MockMessageRepository) Create(ctx context.Context, message *domain.Message) error {
	args := m.Called(ctx, message)
	if args.Error(0) == nil {
		// Simulate DB setting ID and timestamps
		message.ID = 1
		message.CreatedAt = time.Now()
		message.UpdatedAt = time.Now()
	}
	return args.Error(0)
}

func (m *MockMessageRepository) GetByID(ctx context.Context, messageID int64) (*domain.Message, error) {
	args := m.Called(ctx, messageID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Message), args.Error(1)
}

func (m *MockMessageRepository) Delete(ctx context.Context, messageID int64) error {
	args := m.Called(ctx, messageID)
	return args.Error(0)
}

func (m *MockMessageRepository) GetMessagesByCursor(ctx context.Context, chatID int64, beforeMessageID *int64, limit int) ([]*domain.Message, bool, error) {
	args := m.Called(ctx, chatID, beforeMessageID, limit)
	if args.Get(0) == nil {
		return nil, args.Get(1).(bool), args.Error(2)
	}
	return args.Get(0).([]*domain.Message), args.Get(1).(bool), args.Error(2)
}

func (m *MockMessageRepository) GetMessages(ctx context.Context, chatID int64, beforeMessageID, afterMessageID *int64, limit int) ([]*domain.Message, error) {
	args := m.Called(ctx, chatID, beforeMessageID, afterMessageID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Message), args.Error(1)
}

func (m *MockMessageRepository) MarkMessagesAsRead(ctx context.Context, chatID, receiverID int64, messageIDs []int64) (int, error) {
	args := m.Called(ctx, chatID, receiverID, messageIDs)
	return args.Int(0), args.Error(1)
}

func (m *MockMessageRepository) MarkAllAsRead(ctx context.Context, chatID, receiverID int64) (int, error) {
	args := m.Called(ctx, chatID, receiverID)
	return args.Int(0), args.Error(1)
}

func (m *MockMessageRepository) GetUnreadCount(ctx context.Context, chatID, receiverID int64) (int32, error) {
	args := m.Called(ctx, chatID, receiverID)
	return args.Get(0).(int32), args.Error(1)
}

func (m *MockMessageRepository) GetUnreadCountByUser(ctx context.Context, receiverID int64) (int32, error) {
	args := m.Called(ctx, receiverID)
	return args.Get(0).(int32), args.Error(1)
}

// MockAttachmentRepository is a mock for AttachmentRepository
type MockAttachmentRepository struct {
	mock.Mock
}

func (m *MockAttachmentRepository) Create(ctx context.Context, attachment *domain.ChatAttachment) error {
	args := m.Called(ctx, attachment)
	if args.Error(0) == nil {
		// Simulate DB setting ID and timestamp
		attachment.ID = 1
		attachment.CreatedAt = time.Now()
	}
	return args.Error(0)
}

func (m *MockAttachmentRepository) GetByID(ctx context.Context, attachmentID int64) (*domain.ChatAttachment, error) {
	args := m.Called(ctx, attachmentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ChatAttachment), args.Error(1)
}

func (m *MockAttachmentRepository) GetByMessageID(ctx context.Context, messageID int64) ([]*domain.ChatAttachment, error) {
	args := m.Called(ctx, messageID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ChatAttachment), args.Error(1)
}

func (m *MockAttachmentRepository) Delete(ctx context.Context, attachmentID int64) error {
	args := m.Called(ctx, attachmentID)
	return args.Error(0)
}

func (m *MockAttachmentRepository) CreateBatch(ctx context.Context, attachments []*domain.ChatAttachment) error {
	args := m.Called(ctx, attachments)
	return args.Error(0)
}

// MockProductsRepository is a mock for products repository
type MockProductsRepository struct {
	mock.Mock
}

func (m *MockProductsRepository) GetListingByID(ctx context.Context, listingID int64) (*domain.Listing, error) {
	args := m.Called(ctx, listingID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Listing), args.Error(1)
}

// =============================================================================
// TEST: GetChat
// =============================================================================

func TestChatService_GetChat_Success(t *testing.T) {
	service, chatRepo, messageRepo, _, _ := setupTestChatService(t)
	ctx := context.Background()

	chatID := int64(1)
	userID := int64(10)

	chat := &domain.Chat{
		ID:       chatID,
		BuyerID:  userID,
		SellerID: 20,
		Status:   domain.ChatStatusActive,
	}

	chatRepo.On("GetByID", ctx, chatID).Return(chat, nil)
	messageRepo.On("GetUnreadCount", ctx, chatID, userID).Return(int32(5), nil)

	result, err := service.GetChat(ctx, chatID, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, chatID, result.ID)
	assert.Equal(t, int32(5), result.UnreadCount)

	chatRepo.AssertExpectations(t)
	messageRepo.AssertExpectations(t)
}

func TestChatService_GetChat_NotParticipant(t *testing.T) {
	service, chatRepo, _, _, _ := setupTestChatService(t)
	ctx := context.Background()

	chatID := int64(1)
	userID := int64(999) // Not a participant

	chat := &domain.Chat{
		ID:       chatID,
		BuyerID:  10,
		SellerID: 20,
		Status:   domain.ChatStatusActive,
	}

	chatRepo.On("GetByID", ctx, chatID).Return(chat, nil)

	result, err := service.GetChat(ctx, chatID, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, ErrNotParticipant, err)

	chatRepo.AssertExpectations(t)
}

func TestChatService_GetChat_NotFound(t *testing.T) {
	service, chatRepo, _, _, _ := setupTestChatService(t)
	ctx := context.Background()

	chatID := int64(999)
	userID := int64(10)

	chatRepo.On("GetByID", ctx, chatID).Return(nil, ErrChatNotFound)

	result, err := service.GetChat(ctx, chatID, userID)

	assert.Error(t, err)
	assert.Nil(t, result)

	chatRepo.AssertExpectations(t)
}

// =============================================================================
// TEST: GetUserChats
// =============================================================================

func TestChatService_GetUserChats_Success(t *testing.T) {
	service, chatRepo, messageRepo, _, _ := setupTestChatService(t)
	ctx := context.Background()

	userID := int64(1)
	req := &GetUserChatsRequest{
		UserID:   userID,
		Archived: false,
		Limit:    20,
		Offset:   0,
	}

	chats := []*domain.Chat{
		{ID: 1, BuyerID: userID, SellerID: 2},
		{ID: 2, BuyerID: userID, SellerID: 3},
	}

	chatRepo.On("GetUserChats", ctx, userID, (*domain.ChatStatus)(nil), false, 20, 0).
		Return(chats, 2, nil)

	// Mock unread counts for each chat
	messageRepo.On("GetUnreadCount", ctx, int64(1), userID).Return(int32(3), nil)
	messageRepo.On("GetUnreadCount", ctx, int64(2), userID).Return(int32(0), nil)

	result, total, err := service.GetUserChats(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, 2, total)
	assert.Len(t, result, 2)
	assert.Equal(t, int32(3), result[0].UnreadCount)
	assert.Equal(t, int32(0), result[1].UnreadCount)

	chatRepo.AssertExpectations(t)
	messageRepo.AssertExpectations(t)
}

func TestChatService_GetUserChats_WithPagination(t *testing.T) {
	service, chatRepo, _, _, _ := setupTestChatService(t)
	ctx := context.Background()

	userID := int64(1)
	req := &GetUserChatsRequest{
		UserID:   userID,
		Archived: false,
		Limit:    200, // Exceeds max
		Offset:   10,
	}

	chats := []*domain.Chat{}

	// Should limit to 100 max
	chatRepo.On("GetUserChats", ctx, userID, (*domain.ChatStatus)(nil), false, 100, 10).
		Return(chats, 0, nil)

	result, total, err := service.GetUserChats(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, 0, total)
	assert.Len(t, result, 0)

	chatRepo.AssertExpectations(t)
}

// =============================================================================
// TEST: GetMessages
// =============================================================================

func TestChatService_GetMessages_Success(t *testing.T) {
	service, chatRepo, messageRepo, _, _ := setupTestChatService(t)
	ctx := context.Background()

	chatID := int64(1)
	userID := int64(10)

	chat := &domain.Chat{
		ID:       chatID,
		BuyerID:  userID,
		SellerID: 20,
	}

	messages := []*domain.Message{
		{ID: 1, ChatID: chatID, SenderID: userID, Content: "Hello"},
		{ID: 2, ChatID: chatID, SenderID: 20, Content: "Hi"},
	}

	chatRepo.On("GetByID", ctx, chatID).Return(chat, nil)
	messageRepo.On("GetUnreadCount", ctx, chatID, userID).Return(int32(1), nil)
	messageRepo.On("GetMessages", ctx, chatID, (*int64)(nil), (*int64)(nil), 50).
		Return(messages, nil)

	req := &GetMessagesRequest{
		ChatID: chatID,
		UserID: userID,
		Limit:  0, // Will use default 50
	}

	result, hasMore, err := service.GetMessages(ctx, req)

	assert.NoError(t, err)
	assert.False(t, hasMore)
	assert.Len(t, result, 2)

	chatRepo.AssertExpectations(t)
	messageRepo.AssertExpectations(t)
}

func TestChatService_GetMessages_NotParticipant(t *testing.T) {
	service, chatRepo, _, _, _ := setupTestChatService(t)
	ctx := context.Background()

	chatID := int64(1)
	userID := int64(999) // Not a participant

	chat := &domain.Chat{
		ID:       chatID,
		BuyerID:  10,
		SellerID: 20,
	}

	chatRepo.On("GetByID", ctx, chatID).Return(chat, nil)

	req := &GetMessagesRequest{
		ChatID: chatID,
		UserID: userID,
	}

	result, hasMore, err := service.GetMessages(ctx, req)

	assert.Error(t, err)
	assert.False(t, hasMore)
	assert.Nil(t, result)
	assert.Equal(t, ErrNotParticipant, err)

	chatRepo.AssertExpectations(t)
}

// =============================================================================
// TEST: MarkMessagesAsRead
// =============================================================================

func TestChatService_MarkMessagesAsRead_SpecificMessages(t *testing.T) {
	service, chatRepo, messageRepo, _, _ := setupTestChatService(t)
	ctx := context.Background()

	chatID := int64(1)
	userID := int64(10)
	messageIDs := []int64{1, 2, 3}

	chat := &domain.Chat{
		ID:       chatID,
		BuyerID:  userID,
		SellerID: 20,
	}

	chatRepo.On("GetByID", ctx, chatID).Return(chat, nil)
	messageRepo.On("GetUnreadCount", ctx, chatID, userID).Return(int32(3), nil)
	messageRepo.On("MarkMessagesAsRead", ctx, chatID, userID, messageIDs).
		Return(3, nil)

	req := &MarkMessagesAsReadRequest{
		ChatID:     chatID,
		UserID:     userID,
		MessageIDs: messageIDs,
	}

	count, err := service.MarkMessagesAsRead(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, 3, count)

	chatRepo.AssertExpectations(t)
	messageRepo.AssertExpectations(t)
}

func TestChatService_MarkMessagesAsRead_MarkAll(t *testing.T) {
	service, chatRepo, messageRepo, _, _ := setupTestChatService(t)
	ctx := context.Background()

	chatID := int64(1)
	userID := int64(10)

	chat := &domain.Chat{
		ID:       chatID,
		BuyerID:  userID,
		SellerID: 20,
	}

	chatRepo.On("GetByID", ctx, chatID).Return(chat, nil)
	messageRepo.On("GetUnreadCount", ctx, chatID, userID).Return(int32(5), nil)
	messageRepo.On("MarkAllAsRead", ctx, chatID, userID).
		Return(5, nil)

	req := &MarkMessagesAsReadRequest{
		ChatID:  chatID,
		UserID:  userID,
		MarkAll: true, // Mark all messages
	}

	count, err := service.MarkMessagesAsRead(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, 5, count)

	chatRepo.AssertExpectations(t)
	messageRepo.AssertExpectations(t)
}

// =============================================================================
// TEST: GetUnreadCount
// =============================================================================

func TestChatService_GetUnreadCount_SpecificChat(t *testing.T) {
	service, _, messageRepo, _, _ := setupTestChatService(t)
	ctx := context.Background()

	chatID := int64(1)
	userID := int64(10)

	messageRepo.On("GetUnreadCount", ctx, chatID, userID).
		Return(int32(5), nil)

	count, err := service.GetUnreadCount(ctx, userID, &chatID)

	assert.NoError(t, err)
	assert.Equal(t, 5, count)

	messageRepo.AssertExpectations(t)
}

func TestChatService_GetUnreadCount_AllChats(t *testing.T) {
	service, _, messageRepo, _, _ := setupTestChatService(t)
	ctx := context.Background()

	userID := int64(10)

	messageRepo.On("GetUnreadCountByUser", ctx, userID).
		Return(int32(12), nil)

	count, err := service.GetUnreadCount(ctx, userID, nil)

	assert.NoError(t, err)
	assert.Equal(t, 12, count)

	messageRepo.AssertExpectations(t)
}

// =============================================================================
// TEST: Archive/DeleteChat
// =============================================================================

func TestChatService_ArchiveChat_Success(t *testing.T) {
	service, chatRepo, messageRepo, _, _ := setupTestChatService(t)
	ctx := context.Background()

	chatID := int64(1)
	userID := int64(10)

	chat := &domain.Chat{
		ID:       chatID,
		BuyerID:  userID,
		SellerID: 20,
		Status:   domain.ChatStatusActive,
	}

	chatRepo.On("GetByID", ctx, chatID).Return(chat, nil)
	messageRepo.On("GetUnreadCount", ctx, chatID, userID).Return(int32(0), nil)
	chatRepo.On("ArchiveChat", ctx, chatID, true).Return(nil)

	err := service.ArchiveChat(ctx, chatID, userID, true)

	assert.NoError(t, err)

	chatRepo.AssertExpectations(t)
	messageRepo.AssertExpectations(t)
}

func TestChatService_DeleteChat_Success(t *testing.T) {
	service, chatRepo, _, _, _ := setupTestChatService(t)
	ctx := context.Background()

	chatID := int64(1)

	chat := &domain.Chat{
		ID:       chatID,
		BuyerID:  10,
		SellerID: 20,
		Status:   domain.ChatStatusActive,
	}

	chatRepo.On("GetByID", ctx, chatID).Return(chat, nil)
	chatRepo.On("Delete", ctx, chatID).Return(nil)

	err := service.DeleteChat(ctx, chatID)

	assert.NoError(t, err)

	chatRepo.AssertExpectations(t)
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

func setupTestChatService(t *testing.T) (ChatService, *MockChatRepository, *MockMessageRepository, *MockAttachmentRepository, *MockProductsRepository) {
	logger := zerolog.New(nil).Level(zerolog.Disabled)

	chatRepo := new(MockChatRepository)
	messageRepo := new(MockMessageRepository)
	attachmentRepo := new(MockAttachmentRepository)
	productsRepo := new(MockProductsRepository)

	var pool *pgxpool.Pool

	// Use real chatService with mocked repositories
	// NOTE: productsRepo will be nil since we can't easily mock the concrete type
	service := &chatService{
		chatRepo:       chatRepo,
		messageRepo:    messageRepo,
		attachmentRepo: attachmentRepo,
		productsRepo:   nil, // Cannot mock concrete type without refactoring
		authService:    nil,
		pool:           pool,
		logger:         logger,
	}

	return service, chatRepo, messageRepo, attachmentRepo, productsRepo
}
