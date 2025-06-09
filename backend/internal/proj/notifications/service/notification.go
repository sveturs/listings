package service

import (
	"backend/internal/domain/models"
	"backend/internal/storage"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type NotificationService struct {
	storage storage.Storage
	bot     *tgbotapi.BotAPI
	email   *EmailService
}

func NewNotificationService(storage storage.Storage) NotificationServiceInterface {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Printf("Error initializing telegram bot: %v", err)
	}

	// Инициализируем сервис для отправки email
	emailService := NewEmailService(
		"mailserver",                // Используем имя контейнера вместо домена
		"25",                        // SMTP-порт без шифрования
		"info@svetu.rs",             // Адрес отправителя
		"SveTu.rs",                  // Имя отправителя
		"info@svetu.rs",             // Имя пользователя
		os.Getenv("EMAIL_PASSWORD"), // Пароль
	)

	return &NotificationService{
		storage: storage,
		bot:     bot,
		email:   emailService,
	}
}

func (s *NotificationService) GetNotificationSettings(ctx context.Context, userID int) ([]models.NotificationSettings, error) {
	return s.storage.GetNotificationSettings(ctx, userID)
}

func (s *NotificationService) UpdateNotificationSettings(ctx context.Context, settings *models.NotificationSettings) error {
	// Логирование для отладки
	log.Printf("Updating notification settings: %+v", settings)

	// Проверка на валидность типа уведомления
	validTypes := map[string]bool{
		models.NotificationTypeNewMessage:     true,
		models.NotificationTypeNewReview:      true,
		models.NotificationTypeReviewVote:     true,
		models.NotificationTypeReviewResponse: true,
		models.NotificationTypeListingStatus:  true,
		models.NotificationTypeFavoritePrice:  true,
	}

	if !validTypes[settings.NotificationType] {
		log.Printf("Invalid notification type: %s", settings.NotificationType)
		return fmt.Errorf("invalid notification type: %s", settings.NotificationType)
	}

	// Проверим текущие настройки, чтобы не затереть случайно другие значения
	existingSettings, err := s.storage.GetNotificationSettings(ctx, settings.UserID)
	if err != nil {
		log.Printf("Error getting existing settings: %v", err)
		// Продолжаем выполнение, будем использовать то, что пришло
	} else {
		// Найдем текущие настройки для этого типа
		var existingSetting *models.NotificationSettings
		for i := range existingSettings {
			if existingSettings[i].NotificationType == settings.NotificationType {
				existingSetting = &existingSettings[i]
				break
			}
		}

		// Если нашли существующие настройки, проверим, не затираем ли мы случайно поля
		if existingSetting != nil {
			// Если обновление касается только telegram_enabled, сохраним текущее значение email_enabled
			if settings.EmailEnabled == false && existingSetting.EmailEnabled == true {
				if c := ctx.Value("request"); c != nil {
					if req, ok := c.(map[string]interface{}); ok {
						if _, exists := req["email_enabled"]; !exists {
							log.Printf("Preserving email_enabled=true for notification type %s", settings.NotificationType)
							settings.EmailEnabled = true
						}
					}
				}
			}

			// И наоборот - если обновление касается только email_enabled
			if settings.TelegramEnabled == false && existingSetting.TelegramEnabled == true {
				if c := ctx.Value("request"); c != nil {
					if req, ok := c.(map[string]interface{}); ok {
						if _, exists := req["telegram_enabled"]; !exists {
							log.Printf("Preserving telegram_enabled=true for notification type %s", settings.NotificationType)
							settings.TelegramEnabled = true
						}
					}
				}
			}
		}
	}

	err = s.storage.UpdateNotificationSettings(ctx, settings)
	if err != nil {
		log.Printf("Error updating settings in DB: %v", err)
		return err
	}

	return nil
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

	// Получаем информацию о пользователе для email
	user, err := s.storage.GetUserByID(ctx, userID)
	if err != nil {
		log.Printf("Error getting user info for notifications: %v", err)
		// Продолжаем выполнение, возможно, отправим только Telegram
	}

	// Создаем базовое уведомление для БД
	notification := &models.Notification{
		UserID:    userID,
		Type:      notificationType,
		Title:     s.getNotificationTitle(notificationType),
		Message:   message,
		ListingID: listingID,
	}

	// Сохраняем уведомление в БД
	if err := s.storage.CreateNotification(ctx, notification); err != nil {
		log.Printf("Error creating notification in DB: %v", err)
		// Продолжаем выполнение, попробуем отправить уведомления
	}

	// Отправляем в Telegram, если включено
	if setting != nil && setting.TelegramEnabled {
		if err := s.sendTelegramNotification(ctx, userID, message, listingID); err != nil {
			log.Printf("Error sending Telegram notification: %v", err)
		}
	}

	// Отправляем на email, если включено
	if setting != nil && setting.EmailEnabled && user != nil && user.Email != "" {
		title := s.getNotificationTitle(notificationType)
		htmlBody := s.email.FormatNotificationEmail(title, message, strconv.Itoa(listingID))

		if err := s.email.SendEmail(user.Email, title, htmlBody); err != nil {
			log.Printf("Error sending email notification: %v", err)
		}
	}

	return nil
}
func (s *NotificationService) sendTelegramNotification(ctx context.Context, userID int, message string, listingID int) error {
	// Получаем Telegram подключение
	conn, err := s.storage.GetTelegramConnection(ctx, userID)
	if err != nil {
		return err
	}

	// Формируем сообщение с полной ссылкой
	messageWithLink := fmt.Sprintf("%s\n\nПерейти к объявлению: https://SveTu.rs/marketplace/listings/%d",
		message,
		listingID)

	chatID, _ := strconv.ParseInt(conn.TelegramChatID, 10, 64)
	msg := tgbotapi.NewMessage(chatID, messageWithLink)
	_, err = s.bot.Send(msg)
	return err
}

// Получение заголовка уведомления в зависимости от типа
func (s *NotificationService) getNotificationTitle(notificationType string) string {
	switch notificationType {
	case models.NotificationTypeNewMessage:
		return "Новое сообщение"
	case models.NotificationTypeNewReview:
		return "Новый отзыв"
	case models.NotificationTypeReviewVote:
		return "Оценка отзыва"
	case models.NotificationTypeReviewResponse:
		return "Ответ на отзыв"
	case models.NotificationTypeListingStatus:
		return "Обновление объявления"
	case models.NotificationTypeFavoritePrice:
		return "Изменение цены"
	default:
		return "Уведомление SveTu.rs"
	}
}
