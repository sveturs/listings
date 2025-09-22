package notifications

import (
	"context"
	"fmt"

	"backend/internal/domain/models"
	notificationService "backend/internal/proj/notifications/service"

	"github.com/rs/zerolog/log"
)

// DeliveryNotificationIntegration интегрирует систему доставки с существующим сервисом уведомлений
type DeliveryNotificationIntegration struct {
	notificationService notificationService.NotificationServiceInterface
}

// NewDeliveryNotificationIntegration создает новую интеграцию
func NewDeliveryNotificationIntegration(notifService notificationService.NotificationServiceInterface) *DeliveryNotificationIntegration {
	return &DeliveryNotificationIntegration{
		notificationService: notifService,
	}
}

// SendDeliveryStatusUpdate отправляет уведомление об изменении статуса доставки
func (i *DeliveryNotificationIntegration) SendDeliveryStatusUpdate(ctx context.Context, userID int, event *StatusChangeEvent) error {
	// Формируем сообщение в зависимости от статуса
	message := i.formatDeliveryMessage(event)

	// Отправляем через существующий сервис уведомлений
	err := i.notificationService.SendNotification(
		ctx,
		userID,
		models.NotificationTypeDeliveryStatus,
		message,
		0, // listingID не используется для доставки
	)
	if err != nil {
		log.Error().Err(err).
			Int("user_id", userID).
			Str("tracking_number", event.TrackingNumber).
			Msg("Failed to send delivery notification")
		return err
	}

	log.Info().
		Int("user_id", userID).
		Str("tracking_number", event.TrackingNumber).
		Str("status", event.NewStatus).
		Msg("Delivery notification sent successfully")

	return nil
}

// formatDeliveryMessage форматирует сообщение о доставке
func (i *DeliveryNotificationIntegration) formatDeliveryMessage(event *StatusChangeEvent) string {
	var emoji string
	var statusText string

	// Подбираем эмодзи и текст для статуса
	switch event.NewStatus {
	case "confirmed":
		emoji = "✅"
		statusText = "подтвержден"
	case "picked_up":
		emoji = "📦"
		statusText = "передан в службу доставки"
	case "in_transit":
		emoji = "🚚"
		statusText = "в пути"
	case "out_for_delivery":
		emoji = "🏃"
		statusText = "передан курьеру для доставки"
	case "delivered":
		emoji = "✨"
		statusText = "доставлен"
	case "failed":
		emoji = "❌"
		statusText = "не удалось доставить"
	case "returned":
		emoji = "↩️"
		statusText = "возвращен отправителю"
	case "canceled":
		emoji = "🚫"
		statusText = "отменен"
	default:
		emoji = "📋"
		statusText = event.NewStatus
	}

	// Формируем сообщение
	message := fmt.Sprintf(
		"%s <b>Обновление статуса доставки</b>\n\n"+
			"📦 Трек-номер: <code>%s</code>\n"+
			"📍 Статус: <b>%s</b>\n",
		emoji, event.TrackingNumber, statusText,
	)

	// Добавляем местоположение, если есть
	if event.Location != "" {
		message += fmt.Sprintf("📍 Местоположение: %s\n", event.Location)
	}

	// Добавляем описание, если есть
	if event.Description != "" {
		message += fmt.Sprintf("💬 %s\n", event.Description)
	}

	// Добавляем время события
	message += fmt.Sprintf("\n🕐 Время: %s\n", event.EventTime.Format("02.01.2006 15:04"))

	// Добавляем ссылку на отслеживание
	message += fmt.Sprintf("\n🔗 <a href=\"https://svetu.rs/tracking/%s\">Отследить посылку</a>", event.TrackingNumber)

	// Для критических статусов добавляем призыв к действию
	switch event.NewStatus {
	case "out_for_delivery":
		message += "\n\n⚡ <i>Курьер свяжется с вами в ближайшее время. Пожалуйста, будьте доступны для связи.</i>"
	case "delivered":
		message += "\n\n🎉 <i>Спасибо за покупку! Не забудьте оставить отзыв.</i>"
	case "failed":
		message += "\n\n⚠️ <i>Пожалуйста, свяжитесь с нами для решения вопроса.</i>"
	}

	return message
}

// CheckUserNotificationPreferences проверяет, хочет ли пользователь получать уведомления о доставке
func (i *DeliveryNotificationIntegration) CheckUserNotificationPreferences(ctx context.Context, userID int) (bool, error) {
	settings, err := i.notificationService.GetNotificationSettings(ctx, userID)
	if err != nil {
		return false, err
	}

	// Ищем настройки для уведомлений о доставке
	for _, setting := range settings {
		if setting.NotificationType == models.NotificationTypeDeliveryStatus {
			// Возвращаем true если включен хотя бы один канал
			return setting.TelegramEnabled || setting.EmailEnabled, nil
		}
	}

	// По умолчанию уведомления включены
	return true, nil
}

// SendBulkDeliveryUpdate отправляет массовые уведомления (например, для всех заказов с определенным статусом)
func (i *DeliveryNotificationIntegration) SendBulkDeliveryUpdate(ctx context.Context, updates map[int]*StatusChangeEvent) error {
	var errors []error

	for userID, event := range updates {
		// Проверяем предпочтения пользователя
		shouldNotify, err := i.CheckUserNotificationPreferences(ctx, userID)
		if err != nil {
			log.Warn().Err(err).Int("user_id", userID).Msg("Failed to check user preferences")
			continue
		}

		if !shouldNotify {
			log.Debug().Int("user_id", userID).Msg("User has disabled delivery notifications")
			continue
		}

		// Отправляем уведомление
		if err := i.SendDeliveryStatusUpdate(ctx, userID, event); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to send %d notifications", len(errors))
	}

	return nil
}

// FormatTrackingLink форматирует ссылку для отслеживания
func (i *DeliveryNotificationIntegration) FormatTrackingLink(trackingNumber string) string {
	return fmt.Sprintf("https://svetu.rs/tracking/%s", trackingNumber)
}

// SendDeliveryReminder отправляет напоминание о доставке
func (i *DeliveryNotificationIntegration) SendDeliveryReminder(ctx context.Context, userID int, trackingNumber string, reminderType string) error {
	var message string

	switch reminderType {
	case "pickup_ready":
		message = fmt.Sprintf(
			"📦 <b>Напоминание</b>\n\n"+
				"Ваш заказ <code>%s</code> готов к получению.\n"+
				"Пожалуйста, заберите его в ближайшее время.\n\n"+
				"🔗 <a href=\"%s\">Подробности</a>",
			trackingNumber,
			i.FormatTrackingLink(trackingNumber),
		)

	case "delivery_today":
		message = fmt.Sprintf(
			"🚚 <b>Доставка сегодня</b>\n\n"+
				"Ваш заказ <code>%s</code> будет доставлен сегодня.\n"+
				"Пожалуйста, будьте доступны для связи.\n\n"+
				"🔗 <a href=\"%s\">Отследить</a>",
			trackingNumber,
			i.FormatTrackingLink(trackingNumber),
		)

	case "feedback_request":
		message = fmt.Sprintf(
			"⭐ <b>Оцените доставку</b>\n\n"+
				"Ваш заказ <code>%s</code> был доставлен.\n"+
				"Поделитесь вашим опытом!\n\n"+
				"🔗 <a href=\"%s\">Оставить отзыв</a>",
			trackingNumber,
			i.FormatTrackingLink(trackingNumber),
		)

	default:
		return fmt.Errorf("unknown reminder type: %s", reminderType)
	}

	return i.notificationService.SendNotification(ctx, userID, models.NotificationTypeDeliveryStatus, message, 0)
}
