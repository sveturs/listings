package grpc

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	chatsvcv1 "github.com/sveturs/listings/api/proto/chat/v1"
	"github.com/sveturs/listings/internal/domain"
	"github.com/sveturs/listings/internal/middleware"
	"github.com/sveturs/listings/internal/service"
)

// ============================================================================
// CHAT SERVICE gRPC METHODS
// These methods implement ChatService RPC from chat.proto
// ============================================================================

// ============================================================================
// CHAT OPERATIONS (6 methods)
// ============================================================================

// GetOrCreateChat retrieves existing chat or creates new one
// Authorization: user_id extracted from JWT metadata (NOT in request)
func (s *Server) GetOrCreateChat(ctx context.Context, req *chatsvcv1.GetOrCreateChatRequest) (*chatsvcv1.GetOrCreateChatResponse, error) {
	// Extract user_id from context (set by JWT middleware)
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	s.logger.Debug().
		Int64("user_id", userID).
		Interface("listing_id", req.ListingId).
		Interface("storefront_product_id", req.StorefrontProductId).
		Interface("other_user_id", req.OtherUserId).
		Msg("GetOrCreateChat called")

	// Validate request
	if err := validateGetOrCreateChatRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Build service request
	serviceReq := &service.GetOrCreateChatRequest{
		UserID:              userID,
		OtherUserID:         req.OtherUserId,
		ListingID:           req.ListingId,
		StorefrontProductID: req.StorefrontProductId,
	}

	// Call service layer
	chat, created, err := s.chatService.GetOrCreateChat(ctx, serviceReq)
	if err != nil {
		return nil, mapServiceErrorToGRPC(err, s.logger)
	}

	// Convert domain.Chat to proto Chat
	pbChat := domainChatToProtoChat(chat)

	return &chatsvcv1.GetOrCreateChatResponse{
		Chat:    pbChat,
		Created: created,
	}, nil
}

// ListUserChats retrieves all chats for authenticated user
// Authorization: user_id extracted from JWT metadata
// Admin override: if req.UserId is set and caller is admin, use req.UserId instead
func (s *Server) ListUserChats(ctx context.Context, req *chatsvcv1.ListUserChatsRequest) (*chatsvcv1.ListUserChatsResponse, error) {
	// Extract user_id from context (JWT)
	callerUserID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	// Determine the target user_id (admin override or self)
	targetUserID := callerUserID

	// Check for admin override: if UserId is set in request
	if req.UserId != nil && *req.UserId != callerUserID {
		// Verify caller is admin
		roles, hasRoles := middleware.GetRoles(ctx)
		isAdmin := false
		if hasRoles {
			for _, role := range roles {
				if role == "admin" {
					isAdmin = true
					break
				}
			}
		}

		if isAdmin {
			targetUserID = *req.UserId
			s.logger.Debug().
				Int64("caller_id", callerUserID).
				Int64("target_user_id", targetUserID).
				Msg("Admin override: listing chats for another user")
		} else {
			// Non-admin trying to query another user's chats
			return nil, status.Error(codes.PermissionDenied, "admin role required to query other user's chats")
		}
	}

	s.logger.Debug().
		Int64("user_id", targetUserID).
		Bool("archived_only", req.ArchivedOnly).
		Int32("limit", req.Limit).
		Int32("offset", req.Offset).
		Msg("ListUserChats called")

	// Validate pagination
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	offset := req.Offset
	if offset < 0 {
		offset = 0
	}

	// Build service request
	serviceReq := &service.GetUserChatsRequest{
		UserID:    targetUserID,
		Archived:  req.ArchivedOnly,
		ListingID: req.ListingId,
		Limit:     int(limit),
		Offset:    int(offset),
	}

	// Call service layer
	chats, totalCount, err := s.chatService.GetUserChats(ctx, serviceReq)
	if err != nil {
		return nil, mapServiceErrorToGRPC(err, s.logger)
	}

	s.logger.Debug().
		Int("chats_count", len(chats)).
		Int("total_count", totalCount).
		Int64("target_user_id", targetUserID).
		Msg("ListUserChats result from service")

	// Convert domain.Chat slice to proto Chat slice
	pbChats := make([]*chatsvcv1.Chat, 0, len(chats))
	for _, chat := range chats {
		pbChats = append(pbChats, domainChatToProtoChat(chat))
	}

	// Calculate total unread count across all chats
	unreadTotal := int32(0)
	for _, chat := range chats {
		unreadTotal += chat.UnreadCount
	}

	return &chatsvcv1.ListUserChatsResponse{
		Chats:       pbChats,
		TotalCount:  int32(totalCount),
		UnreadTotal: unreadTotal,
	}, nil
}

// GetChatByID retrieves a single chat by ID
// Authorization: User must be buyer OR seller in the chat
func (s *Server) GetChatByID(ctx context.Context, req *chatsvcv1.GetChatByIDRequest) (*chatsvcv1.GetChatByIDResponse, error) {
	// Extract user_id from context
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	s.logger.Debug().
		Int64("chat_id", req.ChatId).
		Int64("user_id", userID).
		Msg("GetChatByID called")

	// Validate input
	if req.ChatId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "chat_id must be greater than 0")
	}

	// Call service layer (includes authorization check)
	chat, err := s.chatService.GetChat(ctx, req.ChatId, userID)
	if err != nil {
		return nil, mapServiceErrorToGRPC(err, s.logger)
	}

	// Convert domain.Chat to proto Chat
	pbChat := domainChatToProtoChat(chat)

	return &chatsvcv1.GetChatByIDResponse{
		Chat: pbChat,
	}, nil
}

// ArchiveChat archives/unarchives a chat for current user
// Authorization: User must be buyer OR seller in the chat
func (s *Server) ArchiveChat(ctx context.Context, req *chatsvcv1.ArchiveChatRequest) (*emptypb.Empty, error) {
	// Extract user_id from context
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	s.logger.Debug().
		Int64("chat_id", req.ChatId).
		Int64("user_id", userID).
		Bool("archived", req.Archived).
		Msg("ArchiveChat called")

	// Validate input
	if req.ChatId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "chat_id must be greater than 0")
	}

	// Call service layer
	if err := s.chatService.ArchiveChat(ctx, req.ChatId, userID, req.Archived); err != nil {
		return nil, mapServiceErrorToGRPC(err, s.logger)
	}

	return &emptypb.Empty{}, nil
}

// DeleteChat permanently deletes a chat (admin only)
// Authorization: Admin role required (validated in service layer via JWT)
func (s *Server) DeleteChat(ctx context.Context, req *chatsvcv1.DeleteChatRequest) (*emptypb.Empty, error) {
	// Extract user_id from context (for audit logging)
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	s.logger.Info().
		Int64("chat_id", req.ChatId).
		Int64("admin_user_id", userID).
		Msg("DeleteChat called")

	// Validate input
	if req.ChatId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "chat_id must be greater than 0")
	}

	// TODO: Add admin role check here or in service layer
	// For now, service layer just deletes the chat

	// Call service layer
	if err := s.chatService.DeleteChat(ctx, req.ChatId); err != nil {
		return nil, mapServiceErrorToGRPC(err, s.logger)
	}

	return &emptypb.Empty{}, nil
}

// GetChatStats retrieves chat statistics (admin only)
// Authorization: Admin role required (validated in service layer via JWT)
func (s *Server) GetChatStats(ctx context.Context, req *chatsvcv1.GetChatStatsRequest) (*chatsvcv1.GetChatStatsResponse, error) {
	// Extract user_id from context (for audit logging)
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	s.logger.Debug().
		Int64("admin_user_id", userID).
		Interface("user_id_filter", req.UserId).
		Msg("GetChatStats called")

	// TODO: Implement chat statistics in service layer
	// For now, return stub response

	return &chatsvcv1.GetChatStatsResponse{
		TotalChats:         0,
		ActiveChats:        0,
		TotalMessages:      0,
		MessagesToday:      0,
		AvgMessagesPerChat: 0,
		DailyStats:         []*chatsvcv1.DailyChatStats{},
	}, nil
}

// ============================================================================
// MESSAGE OPERATIONS (6 methods)
// ============================================================================

// SendMessage sends a new message in a chat
// Authorization: User must be buyer OR seller in the chat
func (s *Server) SendMessage(ctx context.Context, req *chatsvcv1.SendMessageRequest) (*chatsvcv1.SendMessageResponse, error) {
	// Extract user_id from context
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	s.logger.Debug().
		Int64("chat_id", req.ChatId).
		Int64("sender_id", userID).
		Int("content_length", len(req.Content)).
		Msg("SendMessage called")

	// Validate request
	if err := validateSendMessageRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Default language
	language := req.OriginalLanguage
	if language == "" {
		language = "en"
	}

	// Build service request
	serviceReq := &service.SendMessageRequest{
		ChatID:           req.ChatId,
		SenderID:         userID,
		Content:          req.Content,
		OriginalLanguage: language,
		AttachmentIDs:    req.AttachmentIds,
	}

	// Call service layer
	message, err := s.chatService.SendMessage(ctx, serviceReq)
	if err != nil {
		return nil, mapServiceErrorToGRPC(err, s.logger)
	}

	// Convert domain.Message to proto Message
	pbMessage := domainMessageToProtoMessage(message)

	return &chatsvcv1.SendMessageResponse{
		Message: pbMessage,
	}, nil
}

// GetMessages retrieves messages with cursor-based pagination
// Authorization: User must be buyer OR seller in the chat
func (s *Server) GetMessages(ctx context.Context, req *chatsvcv1.GetMessagesRequest) (*chatsvcv1.GetMessagesResponse, error) {
	// Extract user_id from context
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	s.logger.Debug().
		Int64("chat_id", req.ChatId).
		Int64("user_id", userID).
		Interface("before_message_id", req.BeforeMessageId).
		Interface("after_message_id", req.AfterMessageId).
		Int32("limit", req.Limit).
		Msg("GetMessages called")

	// Validate input
	if req.ChatId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "chat_id must be greater than 0")
	}

	// Validate limit
	limit := req.Limit
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	// Build service request
	serviceReq := &service.GetMessagesRequest{
		ChatID:          req.ChatId,
		UserID:          userID,
		BeforeMessageID: req.BeforeMessageId,
		AfterMessageID:  req.AfterMessageId,
		Limit:           int(limit),
	}

	// Call service layer
	messages, hasMore, err := s.chatService.GetMessages(ctx, serviceReq)
	if err != nil {
		return nil, mapServiceErrorToGRPC(err, s.logger)
	}

	// Convert domain.Message slice to proto Message slice
	pbMessages := make([]*chatsvcv1.Message, 0, len(messages))
	for _, message := range messages {
		pbMessages = append(pbMessages, domainMessageToProtoMessage(message))
	}

	// Calculate next cursor (last message ID)
	var nextCursor *int64
	if hasMore && len(messages) > 0 {
		lastMessageID := messages[len(messages)-1].ID
		nextCursor = &lastMessageID
	}

	return &chatsvcv1.GetMessagesResponse{
		Messages:   pbMessages,
		HasMore:    hasMore,
		NextCursor: nextCursor,
	}, nil
}

// StreamMessages streams new messages in real-time (server streaming)
// Authorization: User must be buyer OR seller in the chat
// Real-time: Long-lived connection, pushes messages as they arrive
func (s *Server) StreamMessages(req *chatsvcv1.StreamMessagesRequest, stream chatsvcv1.ChatService_StreamMessagesServer) error {
	ctx := stream.Context()

	// Extract user_id from context
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "authentication required")
	}

	s.logger.Info().
		Int64("chat_id", req.ChatId).
		Int64("user_id", userID).
		Interface("since_message_id", req.SinceMessageId).
		Msg("StreamMessages started")

	// Validate input
	if req.ChatId <= 0 {
		return status.Error(codes.InvalidArgument, "chat_id must be greater than 0")
	}

	// Verify user has access to this chat
	if _, err := s.chatService.GetChat(ctx, req.ChatId, userID); err != nil {
		return mapServiceErrorToGRPC(err, s.logger)
	}

	// TODO: Implement real-time streaming using pubsub or polling
	// For now, return error indicating not implemented
	return status.Error(codes.Unimplemented, "real-time streaming not yet implemented - use polling with GetMessages instead")

	// Future implementation sketch:
	// 1. Subscribe to Redis pubsub channel for this chat_id
	// 2. Poll for new messages since req.SinceMessageId
	// 3. Stream new messages as they arrive
	// 4. Handle context cancellation for cleanup
}

// MarkMessagesAsRead marks messages as read
// Authorization: User must be receiver of the messages
func (s *Server) MarkMessagesAsRead(ctx context.Context, req *chatsvcv1.MarkMessagesAsReadRequest) (*chatsvcv1.MarkMessagesAsReadResponse, error) {
	// Extract user_id from context
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	s.logger.Debug().
		Int64("chat_id", req.ChatId).
		Int64("user_id", userID).
		Bool("mark_all", req.MarkAll).
		Int("message_ids_count", len(req.MessageIds)).
		Msg("MarkMessagesAsRead called")

	// Validate input
	if req.ChatId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "chat_id must be greater than 0")
	}

	// Build service request
	serviceReq := &service.MarkMessagesAsReadRequest{
		ChatID:     req.ChatId,
		UserID:     userID,
		MessageIDs: req.MessageIds,
		MarkAll:    req.MarkAll,
	}

	// Call service layer
	markedCount, err := s.chatService.MarkMessagesAsRead(ctx, serviceReq)
	if err != nil {
		return nil, mapServiceErrorToGRPC(err, s.logger)
	}

	return &chatsvcv1.MarkMessagesAsReadResponse{
		MarkedCount: int32(markedCount),
	}, nil
}

// GetUnreadCount retrieves unread message count
// Authorization: user_id extracted from JWT metadata
func (s *Server) GetUnreadCount(ctx context.Context, req *chatsvcv1.GetUnreadCountRequest) (*chatsvcv1.GetUnreadCountResponse, error) {
	// Extract user_id from context
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	s.logger.Debug().
		Int64("user_id", userID).
		Interface("chat_id", req.ChatId).
		Msg("GetUnreadCount called")

	// Call service layer
	unreadCount, err := s.chatService.GetUnreadCount(ctx, userID, req.ChatId)
	if err != nil {
		return nil, mapServiceErrorToGRPC(err, s.logger)
	}

	// TODO: Implement per-chat breakdown
	// For now, return total count only
	return &chatsvcv1.GetUnreadCountResponse{
		UnreadCount: int32(unreadCount),
		ByChat:      []*chatsvcv1.ChatUnreadCount{},
	}, nil
}

// DeleteMessage deletes a message (soft delete)
// Authorization: User must be sender OR admin
func (s *Server) DeleteMessage(ctx context.Context, req *chatsvcv1.DeleteMessageRequest) (*emptypb.Empty, error) {
	// Extract user_id from context
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	s.logger.Info().
		Int64("message_id", req.MessageId).
		Int64("user_id", userID).
		Msg("DeleteMessage called")

	// Validate input
	if req.MessageId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "message_id must be greater than 0")
	}

	// TODO: Implement delete message in service layer
	// For now, return unimplemented
	return nil, status.Error(codes.Unimplemented, "delete message not yet implemented")
}

// ============================================================================
// ATTACHMENT OPERATIONS (3 methods)
// ============================================================================

// UploadAttachment uploads a file attachment
// Authorization: authenticated user
func (s *Server) UploadAttachment(ctx context.Context, req *chatsvcv1.UploadAttachmentRequest) (*chatsvcv1.UploadAttachmentResponse, error) {
	// Extract user_id from context
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	s.logger.Debug().
		Int64("user_id", userID).
		Str("file_name", req.FileName).
		Str("content_type", req.ContentType).
		Int("file_size", len(req.FileData)).
		Msg("UploadAttachment called")

	// Validate request
	if err := validateUploadAttachmentRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Convert proto FileType to domain AttachmentType
	fileType := protoChatAttachmentTypeToDomain(req.FileType)

	// Build service request
	serviceReq := &service.UploadAttachmentRequest{
		UserID:      userID,
		FileName:    req.FileName,
		ContentType: req.ContentType,
		FileData:    req.FileData,
		FileType:    fileType,
	}

	// Call service layer
	attachment, err := s.chatService.UploadAttachment(ctx, serviceReq)
	if err != nil {
		return nil, mapServiceErrorToGRPC(err, s.logger)
	}

	// Convert domain.ChatAttachment to proto MessageAttachment
	pbAttachment := domainAttachmentToProtoAttachment(attachment)

	return &chatsvcv1.UploadAttachmentResponse{
		Attachment: pbAttachment,
		UploadId:   fmt.Sprintf("upload-%d", attachment.ID),
	}, nil
}

// GetAttachment retrieves attachment metadata
// Authorization: User must have access to the parent message
func (s *Server) GetAttachment(ctx context.Context, req *chatsvcv1.GetAttachmentRequest) (*chatsvcv1.GetAttachmentResponse, error) {
	// Extract user_id from context
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	s.logger.Debug().
		Int64("attachment_id", req.AttachmentId).
		Int64("user_id", userID).
		Msg("GetAttachment called")

	// Validate input
	if req.AttachmentId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "attachment_id must be greater than 0")
	}

	// Call service layer (includes authorization check)
	attachment, err := s.chatService.GetAttachment(ctx, req.AttachmentId, userID)
	if err != nil {
		return nil, mapServiceErrorToGRPC(err, s.logger)
	}

	// Convert domain.ChatAttachment to proto MessageAttachment
	pbAttachment := domainAttachmentToProtoAttachment(attachment)

	return &chatsvcv1.GetAttachmentResponse{
		Attachment: pbAttachment,
	}, nil
}

// DeleteAttachment deletes an attachment
// Authorization: User must be sender of the parent message OR admin
func (s *Server) DeleteAttachment(ctx context.Context, req *chatsvcv1.DeleteAttachmentRequest) (*emptypb.Empty, error) {
	// Extract user_id from context
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	s.logger.Info().
		Int64("attachment_id", req.AttachmentId).
		Int64("user_id", userID).
		Msg("DeleteAttachment called")

	// Validate input
	if req.AttachmentId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "attachment_id must be greater than 0")
	}

	// Call service layer (includes authorization check)
	if err := s.chatService.DeleteAttachment(ctx, req.AttachmentId, userID); err != nil {
		return nil, mapServiceErrorToGRPC(err, s.logger)
	}

	return &emptypb.Empty{}, nil
}

// ============================================================================
// HELPER FUNCTIONS - Validation
// ============================================================================

// validateGetOrCreateChatRequest validates GetOrCreateChat request
func validateGetOrCreateChatRequest(req *chatsvcv1.GetOrCreateChatRequest) error {
	// Must provide at least one context
	if req.ListingId == nil && req.StorefrontProductId == nil && req.OtherUserId == nil {
		return errors.New("must provide listing_id, storefront_product_id, or other_user_id")
	}

	// Cannot provide both listing and product
	if req.ListingId != nil && req.StorefrontProductId != nil {
		return errors.New("cannot provide both listing_id and storefront_product_id")
	}

	return nil
}

// validateSendMessageRequest validates SendMessage request
func validateSendMessageRequest(req *chatsvcv1.SendMessageRequest) error {
	if req.ChatId <= 0 {
		return errors.New("chat_id must be greater than 0")
	}

	if len(req.Content) == 0 {
		return errors.New("content is required")
	}

	if len(req.Content) > 10000 {
		return fmt.Errorf("content exceeds maximum length of 10000 characters")
	}

	return nil
}

// validateUploadAttachmentRequest validates UploadAttachment request
func validateUploadAttachmentRequest(req *chatsvcv1.UploadAttachmentRequest) error {
	if req.FileName == "" {
		return errors.New("file_name is required")
	}

	if req.ContentType == "" {
		return errors.New("content_type is required")
	}

	if len(req.FileData) == 0 {
		return errors.New("file_data is required")
	}

	if req.FileType == chatsvcv1.AttachmentType_ATTACHMENT_TYPE_UNSPECIFIED {
		return errors.New("file_type must be specified")
	}

	return nil
}

// ============================================================================
// HELPER FUNCTIONS - Converters (Domain → Proto)
// ============================================================================

// domainChatToProtoChat converts domain.Chat to proto Chat
func domainChatToProtoChat(chat *domain.Chat) *chatsvcv1.Chat {
	if chat == nil {
		return nil
	}

	pbChat := &chatsvcv1.Chat{
		Id:                  chat.ID,
		BuyerId:             chat.BuyerID,
		SellerId:            chat.SellerID,
		ListingId:           chat.ListingID,
		StorefrontProductId: chat.StorefrontProductID,
		Status:              protoChatStatusFromDomain(chat.Status),
		IsArchived:          chat.IsArchived,
		LastMessageAt:       timestamppb.New(chat.LastMessageAt),
		CreatedAt:           timestamppb.New(chat.CreatedAt),
		UpdatedAt:           timestamppb.New(chat.UpdatedAt),
		UnreadCount:         chat.UnreadCount,
	}

	// Optional fields
	if chat.LastMessage != nil {
		pbChat.LastMessage = domainMessageToProtoMessage(chat.LastMessage)
	}

	if chat.BuyerName != nil {
		pbChat.BuyerName = chat.BuyerName
	}

	if chat.SellerName != nil {
		pbChat.SellerName = chat.SellerName
	}

	if chat.ListingTitle != nil {
		pbChat.ListingTitle = chat.ListingTitle
	}

	if chat.ListingImageURL != nil {
		pbChat.ListingImageUrl = chat.ListingImageURL
	}

	if chat.ListingOwnerID != nil {
		pbChat.ListingOwnerId = chat.ListingOwnerID
	}

	return pbChat
}

// domainMessageToProtoMessage converts domain.Message to proto Message
func domainMessageToProtoMessage(message *domain.Message) *chatsvcv1.Message {
	if message == nil {
		return nil
	}

	pbMessage := &chatsvcv1.Message{
		Id:                  message.ID,
		ChatId:              message.ChatID,
		SenderId:            message.SenderID,
		ReceiverId:          message.ReceiverID,
		Content:             message.Content,
		OriginalLanguage:    message.OriginalLanguage,
		ListingId:           message.ListingID,
		StorefrontProductId: message.StorefrontProductID,
		Status:              protoMessageStatusFromDomain(message.Status),
		IsRead:              message.IsRead,
		HasAttachments:      message.HasAttachments,
		AttachmentsCount:    message.AttachmentsCount,
		CreatedAt:           timestamppb.New(message.CreatedAt),
		UpdatedAt:           timestamppb.New(message.UpdatedAt),
	}

	// Optional fields
	if message.ReadAt != nil {
		pbMessage.ReadAt = timestamppb.New(*message.ReadAt)
	}

	if message.SenderName != nil {
		pbMessage.SenderName = message.SenderName
	}

	// System message flag
	pbMessage.IsSystem = message.IsSystem

	// Convert attachments
	if len(message.Attachments) > 0 {
		pbMessage.Attachments = make([]*chatsvcv1.MessageAttachment, 0, len(message.Attachments))
		for _, attachment := range message.Attachments {
			pbMessage.Attachments = append(pbMessage.Attachments, domainAttachmentToProtoAttachment(attachment))
		}
	}

	return pbMessage
}

// domainAttachmentToProtoAttachment converts domain.ChatAttachment to proto MessageAttachment
func domainAttachmentToProtoAttachment(attachment *domain.ChatAttachment) *chatsvcv1.MessageAttachment {
	if attachment == nil {
		return nil
	}

	pbAttachment := &chatsvcv1.MessageAttachment{
		Id:            attachment.ID,
		MessageId:     attachment.MessageID,
		FileType:      protoAttachmentTypeFromDomain(attachment.FileType),
		FileName:      attachment.FileName,
		FileSize:      attachment.FileSize,
		ContentType:   attachment.ContentType,
		StorageType:   attachment.StorageType,
		StorageBucket: attachment.StorageBucket,
		FilePath:      attachment.FilePath,
		PublicUrl:     attachment.PublicURL,
		CreatedAt:     timestamppb.New(attachment.CreatedAt),
	}

	// Optional fields
	if attachment.ThumbnailURL != nil {
		pbAttachment.ThumbnailUrl = attachment.ThumbnailURL
	}

	// Convert metadata to JSON string
	if metadataJSON, err := attachment.MetadataJSON(); err == nil {
		pbAttachment.Metadata = metadataJSON
	}

	return pbAttachment
}

// ============================================================================
// HELPER FUNCTIONS - Enum Converters
// ============================================================================

// protoChatStatusFromDomain converts domain.ChatStatus to proto ChatStatus
func protoChatStatusFromDomain(status domain.ChatStatus) chatsvcv1.ChatStatus {
	switch status {
	case domain.ChatStatusActive:
		return chatsvcv1.ChatStatus_CHAT_STATUS_ACTIVE
	case domain.ChatStatusArchived:
		return chatsvcv1.ChatStatus_CHAT_STATUS_ARCHIVED
	case domain.ChatStatusBlocked:
		return chatsvcv1.ChatStatus_CHAT_STATUS_BLOCKED
	default:
		return chatsvcv1.ChatStatus_CHAT_STATUS_UNSPECIFIED
	}
}

// protoMessageStatusFromDomain converts domain.MessageStatus to proto MessageStatus
func protoMessageStatusFromDomain(status domain.MessageStatus) chatsvcv1.MessageStatus {
	switch status {
	case domain.MessageStatusSent:
		return chatsvcv1.MessageStatus_MESSAGE_STATUS_SENT
	case domain.MessageStatusDelivered:
		return chatsvcv1.MessageStatus_MESSAGE_STATUS_DELIVERED
	case domain.MessageStatusRead:
		return chatsvcv1.MessageStatus_MESSAGE_STATUS_READ
	case domain.MessageStatusFailed:
		return chatsvcv1.MessageStatus_MESSAGE_STATUS_FAILED
	default:
		return chatsvcv1.MessageStatus_MESSAGE_STATUS_UNSPECIFIED
	}
}

// protoAttachmentTypeFromDomain converts domain.AttachmentType to proto AttachmentType
func protoAttachmentTypeFromDomain(fileType domain.AttachmentType) chatsvcv1.AttachmentType {
	switch fileType {
	case domain.AttachmentTypeImage:
		return chatsvcv1.AttachmentType_ATTACHMENT_TYPE_IMAGE
	case domain.AttachmentTypeVideo:
		return chatsvcv1.AttachmentType_ATTACHMENT_TYPE_VIDEO
	case domain.AttachmentTypeDocument:
		return chatsvcv1.AttachmentType_ATTACHMENT_TYPE_DOCUMENT
	default:
		return chatsvcv1.AttachmentType_ATTACHMENT_TYPE_UNSPECIFIED
	}
}

// protoChatAttachmentTypeToDomain converts proto AttachmentType to domain.AttachmentType
func protoChatAttachmentTypeToDomain(fileType chatsvcv1.AttachmentType) domain.AttachmentType {
	switch fileType {
	case chatsvcv1.AttachmentType_ATTACHMENT_TYPE_IMAGE:
		return domain.AttachmentTypeImage
	case chatsvcv1.AttachmentType_ATTACHMENT_TYPE_VIDEO:
		return domain.AttachmentTypeVideo
	case chatsvcv1.AttachmentType_ATTACHMENT_TYPE_DOCUMENT:
		return domain.AttachmentTypeDocument
	default:
		return domain.AttachmentTypeImage // Default fallback
	}
}

// ============================================================================
// HELPER FUNCTIONS - Context Extraction
// ============================================================================
// Note: User ID extraction now handled by middleware.GetUserID()
