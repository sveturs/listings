// backend/internal/proj/marketplace/service/chat_interface.go
package service

import (
	"context"

	"backend/internal/domain/models"
)

type ChatServiceInterface interface {
	// Сообщения
	SendMessage(ctx context.Context, msg *models.MarketplaceMessage) error
	GetMessages(ctx context.Context, listingID, userID int, page, limit int) ([]models.MarketplaceMessage, error)
	GetMessageByID(ctx context.Context, messageID int) (*models.MarketplaceMessage, error)
	MarkMessagesAsRead(ctx context.Context, messageIDs []int, userID int) error

	// Чаты
	GetChats(ctx context.Context, userID int) ([]models.MarketplaceChat, error)
	GetChat(ctx context.Context, chatID, userID int) (*models.MarketplaceChat, error)
	ArchiveChat(ctx context.Context, chatID, userID int) error
	GetUnreadMessagesCount(ctx context.Context, userID int) (int, error)

	// WebSocket
	BroadcastMessage(msg *models.MarketplaceMessage)
	BroadcastMessageWithTranslations(ctx context.Context, msg *models.MarketplaceMessage)
	SubscribeToMessages(userID int) chan *models.MarketplaceMessage
	UnsubscribeFromMessages(userID int)

	// Online status
	SetUserOnline(userID int)
	SetUserOffline(userID int)
	GetOnlineUsers() []int
	IsUserOnline(userID int) bool
	BroadcastUserStatus(userID int, status string)
	SubscribeToStatusUpdates(userID int) chan map[string]interface{}
	UnsubscribeFromStatusUpdates(userID int)
}
