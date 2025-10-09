// backend/internal/proj/c2c/service/chat.go
package service

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"backend/internal/domain/models"
	"backend/internal/logger"
	"backend/internal/storage"
	"backend/pkg/utils"

	notificationService "backend/internal/proj/notifications/service"
	userService "backend/internal/proj/users/service"
)

// ChatContextKey is a type for context keys to avoid collisions
type ChatContextKey string

const (
	// Context keys for chat service
	ChatContextKeyListingExists ChatContextKey = "listing_exists"
)

type ChatService struct {
	storage             storage.Storage
	notificationService notificationService.NotificationServiceInterface
	subscribers         sync.Map
	statusSubscribers   sync.Map
	onlineUsers         sync.Map
	userLastSeen        sync.Map
	chatTranslationSvc  *ChatTranslationService
	userService         userService.UserServiceInterface
}

func NewChatService(storage storage.Storage, notifService notificationService.NotificationServiceInterface) *ChatService {
	return &ChatService{
		storage:             storage,
		notificationService: notifService,
		subscribers:         sync.Map{},
		statusSubscribers:   sync.Map{},
		onlineUsers:         sync.Map{},
		userLastSeen:        sync.Map{},
	}
}

// SetChatTranslationService sets the chat translation service
// This is called by the global service after all dependencies are initialized
func (s *ChatService) SetChatTranslationService(translationService *ChatTranslationService) {
	s.chatTranslationSvc = translationService
}

// SetUserService sets the user service
// This is called by the global service after all dependencies are initialized
func (s *ChatService) SetUserService(usrService userService.UserServiceInterface) {
	s.userService = usrService
}

// Реализация методов для сообщений
func (s *ChatService) SendMessage(ctx context.Context, msg *models.MarketplaceMessage) error { //nolint:contextcheck
	// Санитизация контента сообщения для защиты от XSS
	msg.Content = utils.SanitizeText(msg.Content)

	// Валидация длины сообщения
	if len(msg.Content) == 0 {
		return fmt.Errorf("message content cannot be empty")
	}
	if len(msg.Content) > 10000 {
		return fmt.Errorf("message content too long (max 10000 characters)")
	}

	var listing *models.MarketplaceListing
	listingExists := false

	// Определяем тип сообщения и обрабатываем соответствующим образом
	switch {
	case msg.StorefrontProductID > 0:
		// Получаем информацию о товаре и владельце витрины
		storefrontOwnerID, err := s.storage.GetStorefrontOwnerByProductID(ctx, msg.StorefrontProductID)
		if err != nil {
			log.Printf("Error getting storefront owner for product %d: %v", msg.StorefrontProductID, err)
			return fmt.Errorf("storefront product not found: %d", msg.StorefrontProductID)
		}

		// Устанавливаем получателя как владельца витрины
		msg.ReceiverID = storefrontOwnerID
		log.Printf("Message for storefront product %d will be sent to owner %d", msg.StorefrontProductID, storefrontOwnerID)

		// Создаем виртуальный листинг для отображения
		listing = &models.MarketplaceListing{
			ID:    0,
			Title: fmt.Sprintf("Товар витрины #%d", msg.StorefrontProductID),
		}
		listingExists = false
	case msg.ListingID > 0:
		// Если есть ListingID, пытаемся найти объявление
		var err error
		listing, err = s.storage.GetListingByID(ctx, msg.ListingID)
		if err != nil {
			// Проверяем, уже существует ли чат для этого сообщения
			// Если chat_id уже есть, значит это сообщение в существующем чате
			if msg.ChatID > 0 {
				// Если чат существует, разрешаем отправку даже если листинг не найден
				listingExists = false
				// Создаем пустой листинг для подстановки информации в уведомления
				listing = &models.MarketplaceListing{
					ID:    msg.ListingID,
					Title: "__DELETED_LISTING__",
				}
			} else {
				// Если это новый чат и листинг не найден, возвращаем ошибку
				return err
			}
		} else {
			listingExists = true

			// Проверяем права доступа, только если листинг существует
			if msg.ReceiverID != listing.UserID && msg.SenderID != listing.UserID {
				return fmt.Errorf("permission denied")
			}
		}
	default:
		// Это прямое сообщение контакту без привязки к объявлению
		listingExists = false
		listing = &models.MarketplaceListing{
			ID:    0,
			Title: "Личное сообщение1",
		}
	}

	// Добавляем информацию о том, существует ли листинг в контекст
	// Это будет использовано в CreateMessage
	ctx = context.WithValue(ctx, ChatContextKeyListingExists, listingExists)

	// Определяем язык сообщения для переводов (если ещё не определён и есть сервис переводов)
	if msg.OriginalLanguage == "" && s.chatTranslationSvc != nil {
		if err := s.chatTranslationSvc.DetectAndSetLanguage(ctx, msg); err != nil {
			logger.Warn().Err(err).Msg("Failed to detect message language")
		}
	}

	if err := s.storage.CreateMessage(ctx, msg); err != nil {
		return err
	}

	// Отправляем сообщение через WebSocket с персонализированными переводами
	s.BroadcastMessageWithTranslations(ctx, msg)

	// Асинхронная отправка уведомлений
	if msg.ReceiverID != msg.SenderID {
		// Создаем копию контекста, чтобы его можно было использовать в горутине
		ctxCopy := context.Background()

		// Копируем нужные данные для формирования уведомления
		listingID := 0
		listingTitle := "Личное сообщение"
		if listing != nil {
			listingID = listing.ID
			listingTitle = listing.Title
		}
		senderName := "Пользователь"
		if msg.Sender != nil {
			senderName = msg.Sender.Name
		}
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
				logger.Error().
					Err(err).
					Int("receiverID", receiverID).
					Int("listingID", listingID).
					Msg("Error sending notification")
			}
		}()
	}

	return nil
}

func (s *ChatService) GetMessages(ctx context.Context, listingID, userID int, offset, limit int) ([]models.MarketplaceMessage, error) {
	if limit == 0 {
		limit = 20
	}

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
	log.Printf("BroadcastMessage called: messageID=%d, senderID=%d, receiverID=%d, hasAttachments=%v, attachmentsCount=%d, attachments=%+v",
		msg.ID, msg.SenderID, msg.ReceiverID, msg.HasAttachments, msg.AttachmentsCount, msg.Attachments)

	// Отправляем сообщение только получателю и отправителю
	// Получатель должен получить сообщение
	if receiverCh, ok := s.subscribers.Load(msg.ReceiverID); ok {
		if ch, ok := receiverCh.(chan *models.MarketplaceMessage); ok {
			select {
			case ch <- msg:
				log.Printf("Message sent to receiver %d", msg.ReceiverID)
			default:
				// Канал полный или закрыт, пропускаем
				log.Printf("Failed to send message to receiver %d - channel full or closed", msg.ReceiverID)
			}
		}
	} else {
		log.Printf("No subscriber found for receiver %d", msg.ReceiverID)
	}

	// Отправитель также должен получить сообщение для обновления UI
	if senderCh, ok := s.subscribers.Load(msg.SenderID); ok {
		if ch, ok := senderCh.(chan *models.MarketplaceMessage); ok {
			select {
			case ch <- msg:
			default:
				// Канал полный или закрыт, пропускаем
			}
		}
	}
}

// BroadcastMessageWithTranslations отправляет персонализированные сообщения с переводами
// каждому участнику чата согласно их языковым настройкам
func (s *ChatService) BroadcastMessageWithTranslations(ctx context.Context, msg *models.MarketplaceMessage) {
	log.Printf("BroadcastMessageWithTranslations called: messageID=%d, senderID=%d, receiverID=%d, originalLang=%s",
		msg.ID, msg.SenderID, msg.ReceiverID, msg.OriginalLanguage)

	// Если зависимости не установлены - используем старый механизм
	if s.chatTranslationSvc == nil || s.userService == nil {
		logger.Warn().
			Bool("hasChatTranslationSvc", s.chatTranslationSvc != nil).
			Bool("hasUserService", s.userService != nil).
			Msg("BroadcastMessageWithTranslations: dependencies not set, falling back to BroadcastMessage")
		s.BroadcastMessage(msg)
		return
	}

	logger.Debug().
		Int("messageId", msg.ID).
		Str("originalLang", msg.OriginalLanguage).
		Msg("BroadcastMessageWithTranslations: starting personalized broadcast")

	// Список участников чата
	participants := []int{msg.SenderID, msg.ReceiverID}

	for _, participantID := range participants {
		// Клонируем сообщение для каждого участника
		msgCopy := *msg

		// Получаем настройки участника
		settings, err := s.userService.GetChatSettings(ctx, participantID)
		if err != nil {
			logger.Warn().Err(err).Int("userId", participantID).Msg("Failed to get chat settings")
			// Используем defaults если не удалось загрузить настройки
			settings = &models.ChatUserSettings{
				AutoTranslate:     false,
				PreferredLanguage: "en",
				ShowLanguageBadge: true,
				ModerateTone:      true,
			}
		}

		// Если нужен перевод (автоперевод включен И язык отличается от оригинала)
		if settings.AutoTranslate &&
			msgCopy.OriginalLanguage != "" &&
			msgCopy.OriginalLanguage != settings.PreferredLanguage {

			logger.Debug().
				Int("userId", participantID).
				Str("originalLang", msgCopy.OriginalLanguage).
				Str("preferredLang", settings.PreferredLanguage).
				Msg("Translating WebSocket message")

			// Переводим сообщение (используя Redis кеш!)
			err = s.chatTranslationSvc.TranslateMessage(
				ctx,
				&msgCopy,
				settings.PreferredLanguage,
				settings.ModerateTone,
			)
			if err != nil {
				logger.Warn().Err(err).Msg("WebSocket translation failed")
				// Продолжаем с оригинальным сообщением при ошибке перевода
			}
		}

		// Отправляем персонализированное сообщение через существующий механизм
		if participantCh, ok := s.subscribers.Load(participantID); ok {
			if ch, ok := participantCh.(chan *models.MarketplaceMessage); ok {
				select {
				case ch <- &msgCopy:
					logger.Debug().
						Int("participantId", participantID).
						Str("preferredLang", settings.PreferredLanguage).
						Msg("Personalized message sent to participant")
				default:
					// Канал полный или закрыт, пропускаем
					logger.Warn().Int("participantId", participantID).Msg("Failed to send message - channel full or closed")
				}
			}
		} else {
			logger.Debug().Int("participantId", participantID).Msg("No subscriber found for participant")
		}
	}
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

// GetMessageByID возвращает сообщение по ID
func (s *ChatService) GetMessageByID(ctx context.Context, messageID int) (*models.MarketplaceMessage, error) {
	return s.storage.GetMessageByID(ctx, messageID)
}

// Online status methods
func (s *ChatService) SetUserOnline(userID int) {
	s.onlineUsers.Store(userID, true)
	s.userLastSeen.Delete(userID) // Удаляем время последнего визита для онлайн пользователей
	log.Printf("User %d is now online", userID)
	s.BroadcastUserStatus(userID, "online")
}

func (s *ChatService) SetUserOffline(userID int) {
	s.onlineUsers.Delete(userID)
	s.userLastSeen.Store(userID, time.Now().Format(time.RFC3339))
	log.Printf("User %d is now offline", userID)
	s.BroadcastUserStatus(userID, "offline")
}

func (s *ChatService) GetOnlineUsers() []int {
	var users []int
	s.onlineUsers.Range(func(key, value interface{}) bool {
		if userID, ok := key.(int); ok {
			users = append(users, userID)
		}
		return true
	})
	return users
}

func (s *ChatService) IsUserOnline(userID int) bool {
	_, ok := s.onlineUsers.Load(userID)
	return ok
}

func (s *ChatService) BroadcastUserStatus(userID int, status string) {
	statusMsg := map[string]interface{}{
		"type": "user_" + status,
		"payload": map[string]interface{}{
			"user_id": userID,
			"status":  status,
		},
	}

	// Добавляем last_seen для offline статуса
	if status == "offline" {
		if lastSeen, ok := s.userLastSeen.Load(userID); ok {
			statusMsg["payload"].(map[string]interface{})["last_seen"] = lastSeen
		}
	}

	// Отправляем всем подписчикам
	s.statusSubscribers.Range(func(key, value interface{}) bool {
		if ch, ok := value.(chan map[string]interface{}); ok {
			select {
			case ch <- statusMsg:
			default:
				// Канал полный, пропускаем
			}
		}
		return true
	})
}

func (s *ChatService) SubscribeToStatusUpdates(userID int) chan map[string]interface{} {
	ch := make(chan map[string]interface{}, 100)
	s.statusSubscribers.Store(userID, ch)

	// Отправляем текущий список онлайн пользователей
	go func() {
		onlineUsers := s.GetOnlineUsers()
		ch <- map[string]interface{}{
			"type": "online_users_list",
			"payload": map[string]interface{}{
				"users": onlineUsers,
			},
		}
	}()

	return ch
}

func (s *ChatService) UnsubscribeFromStatusUpdates(userID int) {
	if value, loaded := s.statusSubscribers.LoadAndDelete(userID); loaded {
		if ch, ok := value.(chan map[string]interface{}); ok {
			close(ch)
		}
	}
}
