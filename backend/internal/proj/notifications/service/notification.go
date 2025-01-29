package service

import (
    "context"
    "backend/internal/domain/models"
    "backend/internal/storage"
	"strconv"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
    "fmt"




)

type NotificationService struct {
    storage storage.Storage
	bot     *tgbotapi.BotAPI
}

func NewNotificationService(storage storage.Storage) NotificationServiceInterface {
    bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
    if err != nil {
        log.Printf("Error initializing telegram bot: %v", err)
    }
    
    return &NotificationService{
        storage: storage,
        bot:     bot,
    }
}

func (s *NotificationService) GetNotificationSettings(ctx context.Context, userID int) ([]models.NotificationSettings, error) {
    return s.storage.GetNotificationSettings(ctx, userID)
}

func (s *NotificationService) UpdateNotificationSettings(ctx context.Context, settings *models.NotificationSettings) error {
    return s.storage.UpdateNotificationSettings(ctx, settings)
}

func (s *NotificationService) ConnectTelegram(ctx context.Context, userID int, chatID string, username string) error {
    if chatID == "" {
        return fmt.Errorf("empty chat ID")
    }
    
    // Проверяем существование пользователя перед сохранением
    _, err := s.storage.GetUserByID(ctx, userID)
    if err != nil {
        return fmt.Errorf("user not found: %w", err)
    }

    return s.storage.SaveTelegramConnection(ctx, userID, chatID, username)
}

func (s *NotificationService) GetTelegramConnection(ctx context.Context, userID int) (*models.TelegramConnection, error) {
    return s.storage.GetTelegramConnection(ctx, userID)
}

func (s *NotificationService) DisconnectTelegram(ctx context.Context, userID int) error {
    return s.storage.DeleteTelegramConnection(ctx, userID)
}

func (s *NotificationService) CreateNotification(ctx context.Context, notification *models.Notification) error {
    return s.storage.CreateNotification(ctx, notification)
}

func (s *NotificationService) GetUserNotifications(ctx context.Context, userID int, limit, offset int) ([]models.Notification, error) {
    return s.storage.GetUserNotifications(ctx, userID, limit, offset)
}

func (s *NotificationService) MarkNotificationAsRead(ctx context.Context, userID int, notificationID int) error {
    return s.storage.MarkNotificationAsRead(ctx, userID, notificationID)
}

func (s *NotificationService) SendChatNotification(ctx context.Context, userID int, message string) error {
    notification := &models.Notification{
        UserID:  userID,
        Type:    models.NotificationTypeNewMessage,
        Title:   "Новое сообщение",
        Message: message,
    }
    return s.storage.CreateNotification(ctx, notification)
}

func (s *NotificationService) SendReviewNotification(ctx context.Context, userID int, entityType string, entityID int) error {
    notification := &models.Notification{
        UserID:  userID,
        Type:    models.NotificationTypeNewReview,
        Title:   "Новый отзыв",
        Message: "Получен новый отзыв для " + entityType,
    }
    return s.storage.CreateNotification(ctx, notification)
}

func (s *NotificationService) SendListingUpdateNotification(ctx context.Context, userID int, listingID int, updateType string) error {
    notification := &models.Notification{
        UserID:  userID,
        Type:    models.NotificationTypeListingStatus,
        Title:   "Обновление объявления",
        Message: "Статус объявления изменен на: " + updateType,
    }
    return s.storage.CreateNotification(ctx, notification)
}
// Изменяем сигнатуру метода, добавляя listingID
func (s *NotificationService) SendNotification(ctx context.Context, userID int, notificationType string, message string, listingID int) error {
    // Проверяем настройки пользователя
    settings, err := s.storage.GetNotificationSettings(ctx, userID)
    if err != nil {
        return err
    }

    // Находим нужный тип уведомления
    var setting *models.NotificationSettings
    for _, s := range settings {
        if s.NotificationType == notificationType {
            setting = &s
            break
        }
    }

    // Если уведомления выключены, не отправляем
    if setting == nil || !setting.TelegramEnabled {
        return nil
    }

    // Получаем Telegram подключение
    conn, err := s.storage.GetTelegramConnection(ctx, userID)
    if err != nil {
        return err
    }

    // Формируем сообщение с полной ссылкой
    messageWithLink := fmt.Sprintf("%s\n\nПерейти к объявлению: https://landhub.rs/marketplace/listings/%d", 
        message, 
        listingID)

    chatID, _ := strconv.ParseInt(conn.TelegramChatID, 10, 64)
    msg := tgbotapi.NewMessage(chatID, messageWithLink)
    _, err = s.bot.Send(msg)
    return err
}