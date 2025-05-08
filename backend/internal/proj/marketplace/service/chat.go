// backend/internal/proj/marketplace/service/chat.go
package service

import (
	"backend/internal/domain/models"
	"backend/internal/storage"
	"context"
	"fmt"
	"sync"

	"backend/internal/proj/notifications/service"
)

type ChatService struct {
	storage             storage.Storage
	notificationService service.NotificationServiceInterface
	subscribers         sync.Map
}

func NewChatService(storage storage.Storage, notificationService service.NotificationServiceInterface) *ChatService {
	return &ChatService{
		storage:             storage,
		notificationService: notificationService,
		subscribers:         sync.Map{},
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

	// Отправляем сообщение через WebSocket сразу, не дожидаясь отправки уведомлений
	s.BroadcastMessage(msg)

	// Асинхронная отправка уведомлений
	if msg.ReceiverID != msg.SenderID {
		// Создаем копию контекста, чтобы его можно было использовать в горутине
		ctxCopy := context.Background()

		// Копируем нужные данные для формирования уведомления
		listingID := listing.ID
		listingTitle := listing.Title
		senderName := msg.Sender.Name
		messageContent := msg.Content
		receiverID := msg.ReceiverID

		// Запускаем отправку уведомлений в отдельной горутине
		go func() {
			notificationText := fmt.Sprintf(
				"Новое сообщение от %s\nТовар: %s\n\n%s",
				senderName,
				listingTitle,
				messageContent,
			)

			// Игнорируем ошибки при отправке уведомлений, они не должны влиять на основной поток
			err := s.notificationService.SendNotification(
				ctxCopy,
				receiverID,
				models.NotificationTypeNewMessage,
				notificationText,
				listingID,
			)

			if err != nil {
				// Просто логируем ошибку, не возвращаем ее в основной поток
				fmt.Printf("Error sending notification: %v\n", err)
			}
		}()
	}

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

func (s *ChatService) GetUnreadMessagesCount(ctx context.Context, userID int) (int, error) {
	// Используем storage для получения количества непрочитанных сообщений
	return s.storage.GetUnreadMessagesCount(ctx, userID)
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
