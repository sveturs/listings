package notifications

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/smtp"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// NotificationService - сервис отправки уведомлений
type NotificationService struct {
	db             *sqlx.DB
	emailConfig    *EmailConfig
	smsConfig      *SMSConfig
	viberConfig    *ViberConfig
	telegramConfig *TelegramConfig
}

// EmailConfig - конфигурация email
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	FromAddress  string
	FromName     string
}

// SMSConfig - конфигурация SMS
type SMSConfig struct {
	Provider   string // "twilio", "infobip", "sns"
	APIKey     string
	APISecret  string
	FromNumber string
}

// ViberConfig - конфигурация Viber
type ViberConfig struct {
	BotToken string
	BotName  string
}

// TelegramConfig - конфигурация Telegram
type TelegramConfig struct {
	BotToken    string
	BotUsername string
}

// DeliveryNotification - уведомление о доставке
type DeliveryNotification struct {
	ID           int             `db:"id"`
	ShipmentID   int             `db:"shipment_id"`
	UserID       int             `db:"user_id"`
	Channel      string          `db:"channel"` // email, sms, viber, telegram, push
	Status       string          `db:"status"`  // pending, sent, failed
	Template     string          `db:"template"`
	Data         json.RawMessage `db:"data"`
	SentAt       *time.Time      `db:"sent_at"`
	ErrorMessage *string         `db:"error_message"`
	CreatedAt    time.Time       `db:"created_at"`
}

// StatusChangeEvent - событие изменения статуса
type StatusChangeEvent struct {
	ShipmentID     int
	TrackingNumber string
	OldStatus      string
	NewStatus      string
	Location       string
	Description    string
	EventTime      time.Time
	RecipientEmail string
	RecipientPhone string
	RecipientName  string
}

// NewNotificationService создает новый сервис уведомлений
func NewNotificationService(db *sqlx.DB, emailConfig *EmailConfig, smsConfig *SMSConfig) *NotificationService {
	return &NotificationService{
		db:          db,
		emailConfig: emailConfig,
		smsConfig:   smsConfig,
	}
}

// NotifyStatusChange отправляет уведомления при изменении статуса
func (s *NotificationService) NotifyStatusChange(ctx context.Context, event *StatusChangeEvent) error {
	// Определяем какие уведомления отправлять в зависимости от статуса
	shouldNotify := s.shouldNotifyForStatus(event.NewStatus)
	if !shouldNotify {
		return nil
	}

	// Создаем запись о уведомлении
	_ = &DeliveryNotification{
		ShipmentID: event.ShipmentID,
		Status:     "pending",
		Template:   s.getTemplateForStatus(event.NewStatus),
		Data:       s.prepareTemplateData(event),
		CreatedAt:  time.Now(),
	}
	// TODO: сохранить notif в БД когда будет готова таблица delivery_notifications

	// Отправляем email
	if event.RecipientEmail != "" {
		if err := s.sendEmailNotification(ctx, event); err != nil {
			log.Error().Err(err).Msg("Failed to send email notification")
			// Продолжаем с другими каналами даже если email не удался
		}
	}

	// Отправляем SMS для критических статусов
	if s.isCriticalStatus(event.NewStatus) && event.RecipientPhone != "" {
		if err := s.sendSMSNotification(ctx, event); err != nil {
			log.Error().Err(err).Msg("Failed to send SMS notification")
		}
	}

	// Отправляем push уведомление через WebSocket
	if err := s.sendPushNotification(ctx, event); err != nil {
		log.Error().Err(err).Msg("Failed to send push notification")
	}

	return nil
}

// shouldNotifyForStatus определяет, нужно ли отправлять уведомление для статуса
func (s *NotificationService) shouldNotifyForStatus(status string) bool {
	notifiableStatuses := map[string]bool{
		"confirmed":        true,
		"picked_up":        true,
		"in_transit":       true,
		"out_for_delivery": true,
		"delivered":        true,
		"failed":           true,
		"returned":         true,
		"cancelled":        true,
	}
	return notifiableStatuses[status]
}

// isCriticalStatus определяет критические статусы для SMS
func (s *NotificationService) isCriticalStatus(status string) bool {
	criticalStatuses := map[string]bool{
		"out_for_delivery": true,
		"delivered":        true,
		"failed":           true,
		"returned":         true,
	}
	return criticalStatuses[status]
}

// getTemplateForStatus возвращает шаблон для статуса
func (s *NotificationService) getTemplateForStatus(status string) string {
	templates := map[string]string{
		"confirmed":        "delivery_confirmed",
		"picked_up":        "delivery_picked_up",
		"in_transit":       "delivery_in_transit",
		"out_for_delivery": "delivery_out_for_delivery",
		"delivered":        "delivery_delivered",
		"failed":           "delivery_failed",
		"returned":         "delivery_returned",
		"cancelled":        "delivery_cancelled",
	}
	return templates[status]
}

// prepareTemplateData подготавливает данные для шаблона
func (s *NotificationService) prepareTemplateData(event *StatusChangeEvent) json.RawMessage {
	data := map[string]interface{}{
		"tracking_number": event.TrackingNumber,
		"status":          event.NewStatus,
		"location":        event.Location,
		"description":     event.Description,
		"event_time":      event.EventTime.Format("02.01.2006 15:04"),
		"recipient_name":  event.RecipientName,
	}

	jsonData, _ := json.Marshal(data)
	return jsonData
}

// sendEmailNotification отправляет email уведомление
func (s *NotificationService) sendEmailNotification(ctx context.Context, event *StatusChangeEvent) error {
	if s.emailConfig == nil {
		return errors.New("email config not set")
	}

	subject := s.getEmailSubject(event.NewStatus, event.TrackingNumber)
	body, err := s.renderEmailTemplate(event)
	if err != nil {
		return errors.Wrap(err, "failed to render email template")
	}

	// Формируем email
	from := fmt.Sprintf("%s <%s>", s.emailConfig.FromName, s.emailConfig.FromAddress)
	to := event.RecipientEmail

	msg := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
			"\r\n"+
			"%s",
		from, to, subject, body,
	))

	// Отправляем через SMTP
	auth := smtp.PlainAuth("", s.emailConfig.SMTPUser, s.emailConfig.SMTPPassword, s.emailConfig.SMTPHost)
	addr := fmt.Sprintf("%s:%d", s.emailConfig.SMTPHost, s.emailConfig.SMTPPort)

	if err := smtp.SendMail(addr, auth, s.emailConfig.FromAddress, []string{to}, msg); err != nil {
		return errors.Wrap(err, "failed to send email")
	}

	// Сохраняем в БД
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO delivery_notifications (shipment_id, channel, status, template, data, sent_at)
		VALUES ($1, 'email', 'sent', $2, $3, NOW())
	`, event.ShipmentID, s.getTemplateForStatus(event.NewStatus), s.prepareTemplateData(event))

	return err
}

// getEmailSubject возвращает тему письма
func (s *NotificationService) getEmailSubject(status, trackingNumber string) string {
	subjects := map[string]string{
		"confirmed":        "Заказ #%s подтвержден",
		"picked_up":        "Заказ #%s передан в службу доставки",
		"in_transit":       "Заказ #%s в пути",
		"out_for_delivery": "Заказ #%s передан курьеру",
		"delivered":        "Заказ #%s доставлен",
		"failed":           "Проблема с доставкой заказа #%s",
		"returned":         "Заказ #%s возвращен отправителю",
		"cancelled":        "Заказ #%s отменен",
	}

	format, ok := subjects[status]
	if !ok {
		format = "Обновление статуса заказа #%s"
	}

	return fmt.Sprintf(format, trackingNumber)
}

// renderEmailTemplate рендерит email шаблон
func (s *NotificationService) renderEmailTemplate(event *StatusChangeEvent) (string, error) {
	// HTML шаблон письма
	const emailTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background: #4CAF50; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; background: #f9f9f9; }
        .tracking { background: white; padding: 15px; margin: 20px 0; border-left: 4px solid #4CAF50; }
        .status { font-size: 18px; font-weight: bold; color: #4CAF50; }
        .footer { text-align: center; padding: 20px; color: #666; font-size: 12px; }
        .button { display: inline-block; padding: 12px 24px; background: #4CAF50; color: white; text-decoration: none; border-radius: 4px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Обновление статуса доставки</h1>
        </div>
        <div class="content">
            <p>Здравствуйте, {{.RecipientName}}!</p>

            <div class="tracking">
                <p class="status">Статус: {{.StatusText}}</p>
                <p><strong>Трек-номер:</strong> {{.TrackingNumber}}</p>
                <p><strong>Местоположение:</strong> {{.Location}}</p>
                <p><strong>Время:</strong> {{.EventTime}}</p>
                {{if .Description}}<p><strong>Описание:</strong> {{.Description}}</p>{{end}}
            </div>

            <p style="text-align: center;">
                <a href="https://svetu.rs/tracking/{{.TrackingNumber}}" class="button">Отследить посылку</a>
            </p>
        </div>
        <div class="footer">
            <p>Это автоматическое уведомление. Пожалуйста, не отвечайте на это письмо.</p>
            <p>&copy; 2025 Sve Tu. Все права защищены.</p>
        </div>
    </div>
</body>
</html>
`

	// Подготавливаем данные для шаблона
	data := struct {
		RecipientName  string
		StatusText     string
		TrackingNumber string
		Location       string
		EventTime      string
		Description    string
	}{
		RecipientName:  event.RecipientName,
		StatusText:     s.getStatusText(event.NewStatus),
		TrackingNumber: event.TrackingNumber,
		Location:       event.Location,
		EventTime:      event.EventTime.Format("02.01.2006 15:04"),
		Description:    event.Description,
	}

	// Рендерим шаблон
	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// getStatusText возвращает человекочитаемый текст статуса
func (s *NotificationService) getStatusText(status string) string {
	statusTexts := map[string]string{
		"confirmed":        "Заказ подтвержден",
		"picked_up":        "Передан в службу доставки",
		"in_transit":       "В пути",
		"out_for_delivery": "Передан курьеру для доставки",
		"delivered":        "Доставлен",
		"failed":           "Доставка не удалась",
		"returned":         "Возвращен отправителю",
		"cancelled":        "Отменен",
	}

	if text, ok := statusTexts[status]; ok {
		return text
	}
	return status
}

// sendSMSNotification отправляет SMS уведомление
func (s *NotificationService) sendSMSNotification(ctx context.Context, event *StatusChangeEvent) error {
	if s.smsConfig == nil {
		return errors.New("SMS config not set")
	}

	// Формируем короткое SMS сообщение
	message := s.formatSMSMessage(event)

	// В зависимости от провайдера отправляем SMS
	var err error
	switch s.smsConfig.Provider {
	case "twilio":
		err = s.sendTwilioSMS(event.RecipientPhone, message)
	case "infobip":
		err = s.sendInfobipSMS(event.RecipientPhone, message)
	default:
		// Mock implementation для тестирования
		log.Info().
			Str("to", event.RecipientPhone).
			Str("message", message).
			Msg("Mock SMS sent")
	}

	if err != nil {
		return err
	}

	// Сохраняем в БД
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO delivery_notifications (shipment_id, channel, status, template, data, sent_at)
		VALUES ($1, 'sms', 'sent', $2, $3, NOW())
	`, event.ShipmentID, s.getTemplateForStatus(event.NewStatus), s.prepareTemplateData(event))

	return err
}

// formatSMSMessage форматирует SMS сообщение
func (s *NotificationService) formatSMSMessage(event *StatusChangeEvent) string {
	formats := map[string]string{
		"out_for_delivery": "Ваш заказ %s передан курьеру. Ожидайте доставку сегодня.",
		"delivered":        "Ваш заказ %s доставлен. Спасибо за покупку!",
		"failed":           "Проблема с доставкой заказа %s. Свяжитесь с нами.",
		"returned":         "Заказ %s возвращен отправителю.",
	}

	format, ok := formats[event.NewStatus]
	if !ok {
		format = "Статус заказа %s обновлен. Проверьте детали на сайте."
	}

	return fmt.Sprintf(format, event.TrackingNumber)
}

// sendTwilioSMS отправляет SMS через Twilio
func (s *NotificationService) sendTwilioSMS(to, message string) error {
	// TODO: Implement Twilio integration
	log.Info().
		Str("provider", "twilio").
		Str("to", to).
		Str("message", message).
		Msg("Twilio SMS would be sent")
	return nil
}

// sendInfobipSMS отправляет SMS через Infobip
func (s *NotificationService) sendInfobipSMS(to, message string) error {
	// TODO: Implement Infobip integration
	log.Info().
		Str("provider", "infobip").
		Str("to", to).
		Str("message", message).
		Msg("Infobip SMS would be sent")
	return nil
}

// sendPushNotification отправляет push уведомление через WebSocket
func (s *NotificationService) sendPushNotification(ctx context.Context, event *StatusChangeEvent) error {
	// Формируем push уведомление
	notification := map[string]interface{}{
		"type":            "delivery_status",
		"shipment_id":     event.ShipmentID,
		"tracking_number": event.TrackingNumber,
		"status":          event.NewStatus,
		"title":           s.getEmailSubject(event.NewStatus, event.TrackingNumber),
		"body":            s.getStatusText(event.NewStatus),
		"timestamp":       event.EventTime.Unix(),
	}

	// Отправляем через WebSocket (если подключен)
	// TODO: Integrate with WebSocket hub

	jsonData, _ := json.Marshal(notification)
	log.Info().
		RawJSON("notification", jsonData).
		Msg("Push notification would be sent via WebSocket")

	// Сохраняем в БД
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO delivery_notifications (shipment_id, channel, status, template, data, sent_at)
		VALUES ($1, 'push', 'sent', $2, $3, NOW())
	`, event.ShipmentID, s.getTemplateForStatus(event.NewStatus), jsonData)

	return err
}

// SendViberNotification отправляет уведомление через Viber
func (s *NotificationService) SendViberNotification(ctx context.Context, event *StatusChangeEvent) error {
	if s.viberConfig == nil {
		return errors.New("Viber config not set")
	}

	// TODO: Integrate with Viber API
	log.Info().
		Str("tracking_number", event.TrackingNumber).
		Str("status", event.NewStatus).
		Msg("Viber notification would be sent")

	return nil
}

// SendTelegramNotification отправляет уведомление через Telegram
func (s *NotificationService) SendTelegramNotification(ctx context.Context, event *StatusChangeEvent) error {
	if s.telegramConfig == nil {
		return errors.New("Telegram config not set")
	}

	// TODO: Integrate with Telegram Bot API
	log.Info().
		Str("tracking_number", event.TrackingNumber).
		Str("status", event.NewStatus).
		Msg("Telegram notification would be sent")

	return nil
}

// GetNotificationHistory получает историю уведомлений
func (s *NotificationService) GetNotificationHistory(ctx context.Context, shipmentID int) ([]*DeliveryNotification, error) {
	var notifications []*DeliveryNotification

	err := s.db.SelectContext(ctx, &notifications, `
		SELECT id, shipment_id, channel, status, template, data, sent_at, error_message, created_at
		FROM delivery_notifications
		WHERE shipment_id = $1
		ORDER BY created_at DESC
	`, shipmentID)

	return notifications, err
}

// ResendNotification повторно отправляет уведомление
func (s *NotificationService) ResendNotification(ctx context.Context, notificationID int) error {
	var notif DeliveryNotification
	err := s.db.GetContext(ctx, &notif, `
		SELECT * FROM delivery_notifications WHERE id = $1
	`, notificationID)
	if err != nil {
		return errors.Wrap(err, "failed to get notification")
	}

	// TODO: Implement resend logic based on channel

	return nil
}
