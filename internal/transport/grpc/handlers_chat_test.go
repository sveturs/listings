package grpc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	chatsvcv1 "github.com/vondi-global/listings/api/proto/chat/v1"
	"github.com/vondi-global/listings/internal/domain"
	"github.com/vondi-global/listings/internal/middleware"
	"github.com/vondi-global/listings/internal/service"
)

// =============================================================================
// MOCK CHAT SERVICE
// =============================================================================

type MockChatService struct {
	mock.Mock
}

func (m *MockChatService) CreateChat(ctx context.Context, req *service.CreateChatRequest) (*domain.Chat, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Chat), args.Error(1)
}

func (m *MockChatService) GetOrCreateChat(ctx context.Context, req *service.GetOrCreateChatRequest) (*domain.Chat, bool, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Get(1).(bool), args.Error(2)
	}
	return args.Get(0).(*domain.Chat), args.Get(1).(bool), args.Error(2)
}

func (m *MockChatService) GetChat(ctx context.Context, chatID, userID int64) (*domain.Chat, error) {
	args := m.Called(ctx, chatID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Chat), args.Error(1)
}

func (m *MockChatService) GetUserChats(ctx context.Context, req *service.GetUserChatsRequest) ([]*domain.Chat, int, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int), args.Error(2)
	}
	return args.Get(0).([]*domain.Chat), args.Get(1).(int), args.Error(2)
}

func (m *MockChatService) ArchiveChat(ctx context.Context, chatID, userID int64, archived bool) error {
	args := m.Called(ctx, chatID, userID, archived)
	return args.Error(0)
}

func (m *MockChatService) DeleteChat(ctx context.Context, chatID int64) error {
	args := m.Called(ctx, chatID)
	return args.Error(0)
}

func (m *MockChatService) SendMessage(ctx context.Context, req *service.SendMessageRequest) (*domain.Message, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Message), args.Error(1)
}

func (m *MockChatService) GetMessages(ctx context.Context, req *service.GetMessagesRequest) ([]*domain.Message, bool, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Get(1).(bool), args.Error(2)
	}
	return args.Get(0).([]*domain.Message), args.Get(1).(bool), args.Error(2)
}

func (m *MockChatService) MarkMessagesAsRead(ctx context.Context, req *service.MarkMessagesAsReadRequest) (int, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(int), args.Error(1)
}

func (m *MockChatService) GetUnreadCount(ctx context.Context, userID int64, chatID *int64) (int, error) {
	args := m.Called(ctx, userID, chatID)
	return args.Get(0).(int), args.Error(1)
}

func (m *MockChatService) UploadAttachment(ctx context.Context, req *service.UploadAttachmentRequest) (*domain.ChatAttachment, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ChatAttachment), args.Error(1)
}

func (m *MockChatService) GetAttachment(ctx context.Context, attachmentID, userID int64) (*domain.ChatAttachment, error) {
	args := m.Called(ctx, attachmentID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ChatAttachment), args.Error(1)
}

func (m *MockChatService) DeleteAttachment(ctx context.Context, attachmentID, userID int64) error {
	args := m.Called(ctx, attachmentID, userID)
	return args.Error(0)
}

func (m *MockChatService) SendSystemMessage(ctx context.Context, req *service.SendSystemMessageRequest) (*domain.Message, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Message), args.Error(1)
}

func (m *MockChatService) SetHub(hub service.ChatHub) {
	m.Called(hub)
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

func setupTestChatServer(t *testing.T) (*Server, *MockChatService) {
	mockService := new(MockChatService)
	logger := zerolog.Nop()

	server := &Server{
		chatService: mockService,
		logger:      logger,
	}

	return server, mockService
}

// contextWithUserID creates a context with user_id set (simulating JWT middleware)
func contextWithUserID(userID int64) context.Context {
	ctx := context.Background()
	return context.WithValue(ctx, middleware.UserIDKey{}, userID)
}

// =============================================================================
// TEST: GetOrCreateChat
// =============================================================================

func TestGetOrCreateChat_Success_NewChat(t *testing.T) {
	server, mockService := setupTestChatServer(t)
	userID := int64(10)
	listingID := int64(100)

	ctx := contextWithUserID(userID)
	req := &chatsvcv1.GetOrCreateChatRequest{
		ListingId: &listingID,
	}

	expectedChat := &domain.Chat{
		ID:        1,
		BuyerID:   userID,
		SellerID:  20,
		ListingID: &listingID,
		Status:    domain.ChatStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService.On("GetOrCreateChat", ctx, mock.MatchedBy(func(r *service.GetOrCreateChatRequest) bool {
		return r.UserID == userID && r.ListingID != nil && *r.ListingID == listingID
	})).Return(expectedChat, true, nil)

	resp, err := server.GetOrCreateChat(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Created)
	assert.NotNil(t, resp.Chat)
	assert.Equal(t, int64(1), resp.Chat.Id)
	assert.Equal(t, userID, resp.Chat.BuyerId)

	mockService.AssertExpectations(t)
}

func TestGetOrCreateChat_Success_ExistingChat(t *testing.T) {
	server, mockService := setupTestChatServer(t)
	userID := int64(10)
	otherUserID := int64(20)

	ctx := contextWithUserID(userID)
	req := &chatsvcv1.GetOrCreateChatRequest{
		OtherUserId: &otherUserID,
	}

	existingChat := &domain.Chat{
		ID:        5,
		BuyerID:   userID,
		SellerID:  otherUserID,
		Status:    domain.ChatStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService.On("GetOrCreateChat", ctx, mock.Anything).
		Return(existingChat, false, nil)

	resp, err := server.GetOrCreateChat(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.Created)
	assert.Equal(t, int64(5), resp.Chat.Id)

	mockService.AssertExpectations(t)
}

func TestGetOrCreateChat_NoAuth(t *testing.T) {
	server, _ := setupTestChatServer(t)

	ctx := context.Background() // No user_id
	req := &chatsvcv1.GetOrCreateChatRequest{
		OtherUserId: ptrInt64(20),
	}

	resp, err := server.GetOrCreateChat(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
}

func TestGetOrCreateChat_InvalidRequest(t *testing.T) {
	server, _ := setupTestChatServer(t)

	ctx := contextWithUserID(10)
	req := &chatsvcv1.GetOrCreateChatRequest{
		// No listing_id, storefront_product_id, or other_user_id
	}

	resp, err := server.GetOrCreateChat(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestGetOrCreateChat_BothListingAndProduct(t *testing.T) {
	server, _ := setupTestChatServer(t)

	ctx := contextWithUserID(10)
	req := &chatsvcv1.GetOrCreateChatRequest{
		ListingId:           ptrInt64(100),
		StorefrontProductId: ptrInt64(200),
	}

	resp, err := server.GetOrCreateChat(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

// =============================================================================
// TEST: ListUserChats
// =============================================================================

func TestListUserChats_Success(t *testing.T) {
	server, mockService := setupTestChatServer(t)
	userID := int64(10)

	ctx := contextWithUserID(userID)
	req := &chatsvcv1.ListUserChatsRequest{
		ArchivedOnly: false,
		Limit:        20,
		Offset:       0,
	}

	chats := []*domain.Chat{
		{
			ID:          1,
			BuyerID:     userID,
			SellerID:    20,
			Status:      domain.ChatStatusActive,
			UnreadCount: 3,
		},
		{
			ID:          2,
			BuyerID:     30,
			SellerID:    userID,
			Status:      domain.ChatStatusActive,
			UnreadCount: 0,
		},
	}

	mockService.On("GetUserChats", ctx, mock.MatchedBy(func(r *service.GetUserChatsRequest) bool {
		return r.UserID == userID && r.Limit == 20 && r.Offset == 0
	})).Return(chats, 2, nil)

	resp, err := server.ListUserChats(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Chats, 2)
	assert.Equal(t, int32(2), resp.TotalCount)
	assert.Equal(t, int32(3), resp.UnreadTotal) // 3 + 0

	mockService.AssertExpectations(t)
}

func TestListUserChats_PaginationLimits(t *testing.T) {
	server, mockService := setupTestChatServer(t)
	userID := int64(10)

	ctx := contextWithUserID(userID)
	req := &chatsvcv1.ListUserChatsRequest{
		Limit:  200, // Exceeds max of 100
		Offset: -5,  // Negative, should be corrected to 0
	}

	mockService.On("GetUserChats", ctx, mock.MatchedBy(func(r *service.GetUserChatsRequest) bool {
		return r.Limit == 100 && r.Offset == 0
	})).Return([]*domain.Chat{}, 0, nil)

	resp, err := server.ListUserChats(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	mockService.AssertExpectations(t)
}

// =============================================================================
// TEST: GetChatByID
// =============================================================================

func TestGetChatByID_Success(t *testing.T) {
	server, mockService := setupTestChatServer(t)
	userID := int64(10)
	chatID := int64(1)

	ctx := contextWithUserID(userID)
	req := &chatsvcv1.GetChatByIDRequest{
		ChatId: chatID,
	}

	chat := &domain.Chat{
		ID:       chatID,
		BuyerID:  userID,
		SellerID: 20,
		Status:   domain.ChatStatusActive,
	}

	mockService.On("GetChat", ctx, chatID, userID).
		Return(chat, nil)

	resp, err := server.GetChatByID(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, chatID, resp.Chat.Id)

	mockService.AssertExpectations(t)
}

func TestGetChatByID_InvalidID(t *testing.T) {
	server, _ := setupTestChatServer(t)

	ctx := contextWithUserID(10)
	req := &chatsvcv1.GetChatByIDRequest{
		ChatId: 0,
	}

	resp, err := server.GetChatByID(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestGetChatByID_NotFound(t *testing.T) {
	server, mockService := setupTestChatServer(t)
	userID := int64(10)
	chatID := int64(999)

	ctx := contextWithUserID(userID)
	req := &chatsvcv1.GetChatByIDRequest{
		ChatId: chatID,
	}

	mockService.On("GetChat", ctx, chatID, userID).
		Return(nil, service.ErrChatNotFound)

	resp, err := server.GetChatByID(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())

	mockService.AssertExpectations(t)
}

func TestGetChatByID_NotParticipant(t *testing.T) {
	server, mockService := setupTestChatServer(t)
	userID := int64(10)
	chatID := int64(1)

	ctx := contextWithUserID(userID)
	req := &chatsvcv1.GetChatByIDRequest{
		ChatId: chatID,
	}

	mockService.On("GetChat", ctx, chatID, userID).
		Return(nil, service.ErrNotParticipant)

	resp, err := server.GetChatByID(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.PermissionDenied, st.Code())

	mockService.AssertExpectations(t)
}

// =============================================================================
// TEST: SendMessage
// =============================================================================

func TestSendMessage_Success(t *testing.T) {
	server, mockService := setupTestChatServer(t)
	userID := int64(10)
	chatID := int64(1)

	ctx := contextWithUserID(userID)
	req := &chatsvcv1.SendMessageRequest{
		ChatId:           chatID,
		Content:          "Hello, world!",
		OriginalLanguage: "en",
	}

	message := &domain.Message{
		ID:               1,
		ChatID:           chatID,
		SenderID:         userID,
		ReceiverID:       20,
		Content:          "Hello, world!",
		OriginalLanguage: "en",
		Status:           domain.MessageStatusSent,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	mockService.On("SendMessage", ctx, mock.MatchedBy(func(r *service.SendMessageRequest) bool {
		return r.ChatID == chatID && r.SenderID == userID && r.Content == "Hello, world!"
	})).Return(message, nil)

	resp, err := server.SendMessage(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Message)
	assert.Equal(t, int64(1), resp.Message.Id)
	assert.Equal(t, "Hello, world!", resp.Message.Content)

	mockService.AssertExpectations(t)
}

func TestSendMessage_EmptyContent(t *testing.T) {
	server, _ := setupTestChatServer(t)

	ctx := contextWithUserID(10)
	req := &chatsvcv1.SendMessageRequest{
		ChatId:  1,
		Content: "",
	}

	resp, err := server.SendMessage(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestSendMessage_ContentTooLong(t *testing.T) {
	server, _ := setupTestChatServer(t)

	ctx := contextWithUserID(10)
	req := &chatsvcv1.SendMessageRequest{
		ChatId:  1,
		Content: string(make([]byte, 10001)), // Exceeds 10000 limit
	}

	resp, err := server.SendMessage(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestSendMessage_ChatBlocked(t *testing.T) {
	server, mockService := setupTestChatServer(t)
	userID := int64(10)
	chatID := int64(1)

	ctx := contextWithUserID(userID)
	req := &chatsvcv1.SendMessageRequest{
		ChatId:  chatID,
		Content: "Hello",
	}

	mockService.On("SendMessage", ctx, mock.Anything).
		Return(nil, service.ErrChatBlocked)

	resp, err := server.SendMessage(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.FailedPrecondition, st.Code())

	mockService.AssertExpectations(t)
}

// =============================================================================
// TEST: GetMessages
// =============================================================================

func TestGetMessages_Success(t *testing.T) {
	server, mockService := setupTestChatServer(t)
	userID := int64(10)
	chatID := int64(1)

	ctx := contextWithUserID(userID)
	req := &chatsvcv1.GetMessagesRequest{
		ChatId: chatID,
		Limit:  50,
	}

	messages := []*domain.Message{
		{ID: 1, ChatID: chatID, Content: "Message 1"},
		{ID: 2, ChatID: chatID, Content: "Message 2"},
	}

	mockService.On("GetMessages", ctx, mock.MatchedBy(func(r *service.GetMessagesRequest) bool {
		return r.ChatID == chatID && r.UserID == userID && r.Limit == 50
	})).Return(messages, true, nil)

	resp, err := server.GetMessages(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Messages, 2)
	assert.True(t, resp.HasMore)
	assert.NotNil(t, resp.NextCursor)
	assert.Equal(t, int64(2), *resp.NextCursor)

	mockService.AssertExpectations(t)
}

func TestGetMessages_CursorPagination(t *testing.T) {
	server, mockService := setupTestChatServer(t)
	userID := int64(10)
	chatID := int64(1)
	beforeID := int64(100)

	ctx := contextWithUserID(userID)
	req := &chatsvcv1.GetMessagesRequest{
		ChatId:          chatID,
		BeforeMessageId: &beforeID,
		Limit:           50,
	}

	messages := []*domain.Message{
		{ID: 99, ChatID: chatID, Content: "Message 99"},
	}

	mockService.On("GetMessages", ctx, mock.MatchedBy(func(r *service.GetMessagesRequest) bool {
		return r.BeforeMessageID != nil && *r.BeforeMessageID == beforeID
	})).Return(messages, false, nil)

	resp, err := server.GetMessages(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Messages, 1)
	assert.False(t, resp.HasMore)

	mockService.AssertExpectations(t)
}

func TestGetMessages_LimitValidation(t *testing.T) {
	server, mockService := setupTestChatServer(t)
	userID := int64(10)

	ctx := contextWithUserID(userID)
	req := &chatsvcv1.GetMessagesRequest{
		ChatId: 1,
		Limit:  200, // Exceeds max of 100
	}

	mockService.On("GetMessages", ctx, mock.MatchedBy(func(r *service.GetMessagesRequest) bool {
		return r.Limit == 100
	})).Return([]*domain.Message{}, false, nil)

	resp, err := server.GetMessages(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	mockService.AssertExpectations(t)
}

// =============================================================================
// TEST: MarkMessagesAsRead
// =============================================================================

func TestMarkMessagesAsRead_SpecificMessages(t *testing.T) {
	server, mockService := setupTestChatServer(t)
	userID := int64(10)
	chatID := int64(1)
	messageIDs := []int64{1, 2, 3}

	ctx := contextWithUserID(userID)
	req := &chatsvcv1.MarkMessagesAsReadRequest{
		ChatId:     chatID,
		MessageIds: messageIDs,
		MarkAll:    false,
	}

	mockService.On("MarkMessagesAsRead", ctx, mock.MatchedBy(func(r *service.MarkMessagesAsReadRequest) bool {
		return r.ChatID == chatID && r.UserID == userID && len(r.MessageIDs) == 3
	})).Return(3, nil)

	resp, err := server.MarkMessagesAsRead(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int32(3), resp.MarkedCount)

	mockService.AssertExpectations(t)
}

func TestMarkMessagesAsRead_MarkAll(t *testing.T) {
	server, mockService := setupTestChatServer(t)
	userID := int64(10)
	chatID := int64(1)

	ctx := contextWithUserID(userID)
	req := &chatsvcv1.MarkMessagesAsReadRequest{
		ChatId:  chatID,
		MarkAll: true,
	}

	mockService.On("MarkMessagesAsRead", ctx, mock.MatchedBy(func(r *service.MarkMessagesAsReadRequest) bool {
		return r.MarkAll == true
	})).Return(10, nil)

	resp, err := server.MarkMessagesAsRead(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int32(10), resp.MarkedCount)

	mockService.AssertExpectations(t)
}

func TestMarkMessagesAsRead_InvalidChatID(t *testing.T) {
	server, _ := setupTestChatServer(t)

	ctx := contextWithUserID(10)
	req := &chatsvcv1.MarkMessagesAsReadRequest{
		ChatId: 0,
	}

	resp, err := server.MarkMessagesAsRead(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

// =============================================================================
// TEST: GetUnreadCount
// =============================================================================

func TestGetUnreadCount_SpecificChat(t *testing.T) {
	server, mockService := setupTestChatServer(t)
	userID := int64(10)
	chatID := int64(1)

	ctx := contextWithUserID(userID)
	req := &chatsvcv1.GetUnreadCountRequest{
		ChatId: &chatID,
	}

	mockService.On("GetUnreadCount", ctx, userID, &chatID).
		Return(5, nil)

	resp, err := server.GetUnreadCount(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int32(5), resp.UnreadCount)

	mockService.AssertExpectations(t)
}

func TestGetUnreadCount_AllChats(t *testing.T) {
	server, mockService := setupTestChatServer(t)
	userID := int64(10)

	ctx := contextWithUserID(userID)
	req := &chatsvcv1.GetUnreadCountRequest{
		ChatId: nil,
	}

	mockService.On("GetUnreadCount", ctx, userID, (*int64)(nil)).
		Return(15, nil)

	resp, err := server.GetUnreadCount(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int32(15), resp.UnreadCount)

	mockService.AssertExpectations(t)
}

// =============================================================================
// TEST: UploadAttachment
// =============================================================================

func TestUploadAttachment_Success(t *testing.T) {
	server, mockService := setupTestChatServer(t)
	userID := int64(10)

	ctx := contextWithUserID(userID)
	req := &chatsvcv1.UploadAttachmentRequest{
		FileName:    "test.jpg",
		ContentType: "image/jpeg",
		FileData:    []byte("fake image data"),
		FileType:    chatsvcv1.AttachmentType_ATTACHMENT_TYPE_IMAGE,
	}

	attachment := &domain.ChatAttachment{
		ID:          1,
		MessageID:   0, // Not yet attached to message
		FileName:    "test.jpg",
		FileType:    domain.AttachmentTypeImage,
		FileSize:    15,
		ContentType: "image/jpeg",
		CreatedAt:   time.Now(),
	}

	mockService.On("UploadAttachment", ctx, mock.MatchedBy(func(r *service.UploadAttachmentRequest) bool {
		return r.UserID == userID && r.FileName == "test.jpg"
	})).Return(attachment, nil)

	resp, err := server.UploadAttachment(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Attachment)
	assert.Equal(t, "test.jpg", resp.Attachment.FileName)

	mockService.AssertExpectations(t)
}

func TestUploadAttachment_EmptyFileName(t *testing.T) {
	server, _ := setupTestChatServer(t)

	ctx := contextWithUserID(10)
	req := &chatsvcv1.UploadAttachmentRequest{
		FileName:    "",
		ContentType: "image/jpeg",
		FileData:    []byte("data"),
		FileType:    chatsvcv1.AttachmentType_ATTACHMENT_TYPE_IMAGE,
	}

	resp, err := server.UploadAttachment(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestUploadAttachment_UnspecifiedFileType(t *testing.T) {
	server, _ := setupTestChatServer(t)

	ctx := contextWithUserID(10)
	req := &chatsvcv1.UploadAttachmentRequest{
		FileName:    "test.jpg",
		ContentType: "image/jpeg",
		FileData:    []byte("data"),
		FileType:    chatsvcv1.AttachmentType_ATTACHMENT_TYPE_UNSPECIFIED,
	}

	resp, err := server.UploadAttachment(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestUploadAttachment_FileTooLarge(t *testing.T) {
	server, mockService := setupTestChatServer(t)
	userID := int64(10)

	ctx := contextWithUserID(userID)
	req := &chatsvcv1.UploadAttachmentRequest{
		FileName:    "large.jpg",
		ContentType: "image/jpeg",
		FileData:    make([]byte, 11*1024*1024), // 11MB, exceeds 10MB limit
		FileType:    chatsvcv1.AttachmentType_ATTACHMENT_TYPE_IMAGE,
	}

	mockService.On("UploadAttachment", ctx, mock.Anything).
		Return(nil, &service.ErrAttachmentTooLarge{
			FileType: service.AttachmentFileTypeImage,
			Size:     11 * 1024 * 1024,
			MaxSize:  10 * 1024 * 1024,
		})

	resp, err := server.UploadAttachment(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())

	mockService.AssertExpectations(t)
}

// =============================================================================
// TEST: DeleteAttachment
// =============================================================================

func TestDeleteAttachment_Success(t *testing.T) {
	server, mockService := setupTestChatServer(t)
	userID := int64(10)
	attachmentID := int64(1)

	ctx := contextWithUserID(userID)
	req := &chatsvcv1.DeleteAttachmentRequest{
		AttachmentId: attachmentID,
	}

	mockService.On("DeleteAttachment", ctx, attachmentID, userID).
		Return(nil)

	resp, err := server.DeleteAttachment(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	mockService.AssertExpectations(t)
}

func TestDeleteAttachment_InvalidID(t *testing.T) {
	server, _ := setupTestChatServer(t)

	ctx := contextWithUserID(10)
	req := &chatsvcv1.DeleteAttachmentRequest{
		AttachmentId: 0,
	}

	resp, err := server.DeleteAttachment(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestDeleteAttachment_Unauthorized(t *testing.T) {
	server, mockService := setupTestChatServer(t)
	userID := int64(10)
	attachmentID := int64(1)

	ctx := contextWithUserID(userID)
	req := &chatsvcv1.DeleteAttachmentRequest{
		AttachmentId: attachmentID,
	}

	mockService.On("DeleteAttachment", ctx, attachmentID, userID).
		Return(service.ErrUnauthorized)

	resp, err := server.DeleteAttachment(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.PermissionDenied, st.Code())

	mockService.AssertExpectations(t)
}

// =============================================================================
// TEST: ArchiveChat
// =============================================================================

func TestArchiveChat_Success(t *testing.T) {
	server, mockService := setupTestChatServer(t)
	userID := int64(10)
	chatID := int64(1)

	ctx := contextWithUserID(userID)
	req := &chatsvcv1.ArchiveChatRequest{
		ChatId:   chatID,
		Archived: true,
	}

	mockService.On("ArchiveChat", ctx, chatID, userID, true).
		Return(nil)

	resp, err := server.ArchiveChat(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	mockService.AssertExpectations(t)
}

func TestArchiveChat_InvalidID(t *testing.T) {
	server, _ := setupTestChatServer(t)

	ctx := contextWithUserID(10)
	req := &chatsvcv1.ArchiveChatRequest{
		ChatId:   0,
		Archived: true,
	}

	resp, err := server.ArchiveChat(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

// =============================================================================
// TEST: DeleteChat
// =============================================================================

func TestDeleteChat_Success(t *testing.T) {
	server, mockService := setupTestChatServer(t)
	userID := int64(10)
	chatID := int64(1)

	ctx := contextWithUserID(userID)
	req := &chatsvcv1.DeleteChatRequest{
		ChatId: chatID,
	}

	mockService.On("DeleteChat", ctx, chatID).
		Return(nil)

	resp, err := server.DeleteChat(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	mockService.AssertExpectations(t)
}

func TestDeleteChat_InvalidID(t *testing.T) {
	server, _ := setupTestChatServer(t)

	ctx := contextWithUserID(10)
	req := &chatsvcv1.DeleteChatRequest{
		ChatId: 0,
	}

	resp, err := server.DeleteChat(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

// =============================================================================
// TEST: Error Mapping
// =============================================================================

func TestErrorMapping(t *testing.T) {
	logger := zerolog.Nop()

	tests := []struct {
		name         string
		serviceError error
		expectedCode codes.Code
	}{
		{
			name:         "ErrChatNotFound maps to NotFound",
			serviceError: service.ErrChatNotFound,
			expectedCode: codes.NotFound,
		},
		{
			name:         "ErrNotParticipant maps to PermissionDenied",
			serviceError: service.ErrNotParticipant,
			expectedCode: codes.PermissionDenied,
		},
		{
			name:         "ErrUnauthorized maps to PermissionDenied",
			serviceError: service.ErrUnauthorized,
			expectedCode: codes.PermissionDenied,
		},
		{
			name:         "ErrInvalidInput maps to InvalidArgument",
			serviceError: service.ErrInvalidInput,
			expectedCode: codes.InvalidArgument,
		},
		{
			name:         "ErrChatBlocked maps to FailedPrecondition",
			serviceError: service.ErrChatBlocked,
			expectedCode: codes.FailedPrecondition,
		},
		{
			name:         "Generic error maps to Internal",
			serviceError: errors.New("some error"),
			expectedCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grpcErr := mapServiceErrorToGRPC(tt.serviceError, logger)
			st, ok := status.FromError(grpcErr)
			assert.True(t, ok)
			assert.Equal(t, tt.expectedCode, st.Code())
		})
	}
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

// Note: Helper functions (ptrInt64, etc.) are defined in converters_search_test.go
// to avoid duplication across test files
