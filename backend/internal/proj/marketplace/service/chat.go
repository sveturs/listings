// backend/internal/proj/marketplace/service/chat.go
package service

import (
    "context"
    "fmt"
    "backend/internal/domain/models"
    "backend/internal/storage"
    "sync"
)

type ChatService struct {
    storage storage.Storage
    subscribers sync.Map // map[int]chan *models.MarketplaceMessage
}

func NewChatService(storage storage.Storage) *ChatService {
    return &ChatService{
        storage: storage,
    }
}

// Реализация методов для сообщений
func (s *ChatService) SendMessage(ctx context.Context, msg *models.MarketplaceMessage) error {
    // Проверяем, что отправитель имеет доступ к объявлению
    listing, err := s.storage.GetListingByID(ctx, msg.ListingID)
    if err != nil {
        return err
    }

    // Проверяем, что получатель либо владелец объявления, либо отправитель - владелец
    if msg.ReceiverID != listing.UserID && msg.SenderID != listing.UserID {
        return fmt.Errorf("permission denied")
    }

    if err := s.storage.CreateMessage(ctx, msg); err != nil {
        return err
    }

    // Отправляем сообщение через WebSocket
    s.BroadcastMessage(msg)

    return nil
}

func (s *ChatService) GetMessages(ctx context.Context, listingID, userID int, page, limit int) ([]models.MarketplaceMessage, error) {
    if limit == 0 {
        limit = 20
    }
    offset := (page - 1) * limit
    
    return s.storage.GetMessages(ctx, listingID, userID, offset, limit)
}

func (s *ChatService) MarkMessagesAsRead(ctx context.Context, messageIDs []int, userID int) error {
    return s.storage.MarkMessagesAsRead(ctx, messageIDs, userID)
}

// Реализация методов для чатов
func (s *ChatService) GetChats(ctx context.Context, userID int) ([]models.MarketplaceChat, error) {
    return s.storage.GetChats(ctx, userID)
}

func (s *ChatService) GetChat(ctx context.Context, chatID, userID int) (*models.MarketplaceChat, error) {
    return s.storage.GetChat(ctx, chatID, userID)
}

func (s *ChatService) ArchiveChat(ctx context.Context, chatID, userID int) error {
    return s.storage.ArchiveChat(ctx, chatID, userID)
}

// WebSocket методы
func (s *ChatService) BroadcastMessage(msg *models.MarketplaceMessage) {
    // Отправляем сообщение всем подписчикам
    s.subscribers.Range(func(key, value interface{}) bool {
        if ch, ok := value.(chan *models.MarketplaceMessage); ok {
            // Неблокирующая отправка
            select {
            case ch <- msg:
            default:
                // Канал полный или закрыт, пропускаем
            }
        }
        return true
    })
}

func (s *ChatService) SubscribeToMessages(userID int) chan *models.MarketplaceMessage {
    ch := make(chan *models.MarketplaceMessage, 100)
    s.subscribers.Store(userID, ch)
    return ch
}

func (s *ChatService) UnsubscribeFromMessages(userID int) {
    if value, loaded := s.subscribers.LoadAndDelete(userID); loaded {
        if ch, ok := value.(chan *models.MarketplaceMessage); ok {
            close(ch)
        }
    }
}