package service

import (
	"context"
	"errors"
	"fmt"
	"os"

	"backend/internal/proj/viber/config"
	"backend/internal/proj/viber/infobip"
	"backend/internal/proj/viber/models"
	"backend/internal/storage/postgres"
)

// InfobipBotService сервис для работы с Viber через Infobip
type InfobipBotService struct {
	config         *config.ViberConfig
	infobipClient  *infobip.Client
	db             *postgres.Database
	sessionManager *SessionManager
}

// NewInfobipBotService создаёт новый сервис для работы через Infobip
func NewInfobipBotService(cfg *config.ViberConfig, db *postgres.Database) *InfobipBotService {
	infobipClient := infobip.NewClient(cfg.InfobipAPIKey, cfg.InfobipBaseURL)

	return &InfobipBotService{
		config:         cfg,
		infobipClient:  infobipClient,
		db:             db,
		sessionManager: NewSessionManager(db),
	}
}

// SendTextMessage отправляет текстовое сообщение через Infobip
func (s *InfobipBotService) SendTextMessage(ctx context.Context, viberID, text string) error {
	// Проверяем сессию для определения billable
	session, _ := s.sessionManager.GetActiveSession(ctx, viberID)

	var sessionInfo *infobip.ViberSessionInfo
	if session != nil {
		// Есть активная сессия - сообщение бесплатное
		sessionInfo = &infobip.ViberSessionInfo{
			SessionID: fmt.Sprintf("session_%d", session.ID),
			Origin:    "USER_INITIATED",
		}
	}

	resp, err := s.infobipClient.SendTextMessage(
		ctx,
		s.config.InfobipSenderID,
		viberID,
		text,
		sessionInfo,
	)
	if err != nil {
		return fmt.Errorf("failed to send text message via Infobip: %w", err)
	}

	// Сохраняем в БД
	isBillable := session == nil
	return s.saveOutgoingMessage(ctx, viberID, "text", text, nil, resp.Messages[0].MessageID, isBillable)
}

// SendImageMessage отправляет изображение через Infobip
func (s *InfobipBotService) SendImageMessage(ctx context.Context, viberID, imageURL, text string) error {
	session, _ := s.sessionManager.GetActiveSession(ctx, viberID)

	var sessionInfo *infobip.ViberSessionInfo
	if session != nil {
		sessionInfo = &infobip.ViberSessionInfo{
			SessionID: fmt.Sprintf("session_%d", session.ID),
			Origin:    "USER_INITIATED",
		}
	}

	resp, err := s.infobipClient.SendImageMessage(
		ctx,
		s.config.InfobipSenderID,
		viberID,
		imageURL,
		text,
		sessionInfo,
	)
	if err != nil {
		return fmt.Errorf("failed to send image message via Infobip: %w", err)
	}

	isBillable := session == nil
	return s.saveOutgoingMessage(ctx, viberID, "image", text, map[string]interface{}{"image_url": imageURL}, resp.Messages[0].MessageID, isBillable)
}

// SendButtonMessage отправляет сообщение с кнопкой через Infobip
func (s *InfobipBotService) SendButtonMessage(ctx context.Context, viberID, text, buttonText, buttonURL string) error {
	session, _ := s.sessionManager.GetActiveSession(ctx, viberID)

	var sessionInfo *infobip.ViberSessionInfo
	if session != nil {
		sessionInfo = &infobip.ViberSessionInfo{
			SessionID: fmt.Sprintf("session_%d", session.ID),
			Origin:    "USER_INITIATED",
		}
	}

	resp, err := s.infobipClient.SendButtonMessage(
		ctx,
		s.config.InfobipSenderID,
		viberID,
		text,
		buttonText,
		buttonURL,
		sessionInfo,
	)
	if err != nil {
		return fmt.Errorf("failed to send button message via Infobip: %w", err)
	}

	data := map[string]interface{}{
		"button_text": buttonText,
		"button_url":  buttonURL,
	}

	isBillable := session == nil
	return s.saveOutgoingMessage(ctx, viberID, "button", text, data, resp.Messages[0].MessageID, isBillable)
}

// SendRichMedia отправляет Rich Media сообщение через Infobip
func (s *InfobipBotService) SendRichMedia(ctx context.Context, viberID string, richMedia *models.RichMedia, text string) error {
	session, _ := s.sessionManager.GetActiveSession(ctx, viberID)

	var sessionInfo *infobip.ViberSessionInfo
	if session != nil {
		sessionInfo = &infobip.ViberSessionInfo{
			SessionID: fmt.Sprintf("session_%d", session.ID),
			Origin:    "USER_INITIATED",
		}
	}

	// Конвертируем модель RichMedia в формат Infobip
	infobipRichMedia := s.convertRichMediaToInfobip(richMedia)

	resp, err := s.infobipClient.SendRichMedia(
		ctx,
		s.config.InfobipSenderID,
		viberID,
		infobipRichMedia,
		text,
		sessionInfo,
	)
	if err != nil {
		return fmt.Errorf("failed to send rich media via Infobip: %w", err)
	}

	isBillable := session == nil
	return s.saveOutgoingMessage(ctx, viberID, "rich_media", text, s.richMediaToMap(richMedia), resp.Messages[0].MessageID, isBillable)
}

// SendTrackingNotification отправляет уведомление о трекинге через Infobip
func (s *InfobipBotService) SendTrackingNotification(ctx context.Context, viberID string, delivery *DeliveryInfo) error {
	// Генерируем статическую карту
	mapURL := s.generateStaticMapURL(delivery)

	// Создаём Rich Media с картой и кнопкой трекинга
	buttons := []infobip.ViberRichMediaButton{
		{
			Columns:    6,
			Rows:       4,
			ActionType: "none",
			Image:      mapURL,
		},
		{
			Columns:    6,
			Rows:       1,
			ActionType: "none",
			Text:       fmt.Sprintf("📍 Курьер в пути!\nОжидаемое время: %s", delivery.EstimatedTime.Format("15:04")),
			TextSize:   "medium",
			TextVAlign: "middle",
			TextHAlign: "center",
		},
		{
			Columns:    3,
			Rows:       2,
			ActionType: "open-url",
			ActionBody: fmt.Sprintf("https://svetu.rs/track/%s?viber=true&embedded=true", delivery.TrackingToken),
			Text:       "🗺️ Открыть карту",
			TextSize:   "medium",
			TextVAlign: "middle",
			TextHAlign: "center",
			BgColor:    "#1976d2",
		},
		{
			Columns:    3,
			Rows:       2,
			ActionType: "reply",
			ActionBody: fmt.Sprintf("update_track_%s", delivery.TrackingToken),
			Text:       "🔄 Обновить",
			TextSize:   "medium",
			TextVAlign: "middle",
			TextHAlign: "center",
			BgColor:    "#4caf50",
		},
	}

	infobipRichMedia := &infobip.ViberRichMedia{
		ButtonsGroupColumns: 6,
		ButtonsGroupRows:    7,
		Buttons:             buttons,
	}

	// Проверяем сессию
	session, _ := s.sessionManager.GetActiveSession(ctx, viberID)
	var sessionInfo *infobip.ViberSessionInfo
	if session != nil {
		sessionInfo = &infobip.ViberSessionInfo{
			SessionID: fmt.Sprintf("session_%d", session.ID),
			Origin:    "USER_INITIATED",
		}
	}

	resp, err := s.infobipClient.SendRichMedia(
		ctx,
		s.config.InfobipSenderID,
		viberID,
		infobipRichMedia,
		"Отслеживание доставки",
		sessionInfo,
	)
	if err != nil {
		return fmt.Errorf("failed to send tracking notification via Infobip: %w", err)
	}

	isBillable := session == nil
	return s.saveOutgoingMessage(ctx, viberID, "rich_media", "Отслеживание доставки", nil, resp.Messages[0].MessageID, isBillable)
}

// SendBulkMessages отправляет массовую рассылку через Infobip
func (s *InfobipBotService) SendBulkMessages(ctx context.Context, messages []BulkMessageRequest) error {
	var infobipMessages []infobip.ViberMessage

	for _, msg := range messages {
		infobipMsg := infobip.ViberMessage{
			From: s.config.InfobipSenderID,
			To:   msg.To,
			Content: infobip.ViberContent{
				Type: "TEXT",
				Text: msg.Text,
			},
		}

		// Добавляем метку для промо рассылки
		if msg.IsPromo {
			infobipMsg.Label = &infobip.ViberLabel{
				Type:    "PROMOTION",
				Content: "Специальное предложение",
			}
		}

		infobipMessages = append(infobipMessages, infobipMsg)
	}

	resp, err := s.infobipClient.SendBulkMessages(ctx, infobipMessages)
	if err != nil {
		return fmt.Errorf("failed to send bulk messages via Infobip: %w", err)
	}

	// Сохраняем результаты в БД
	for i, status := range resp.Messages {
		_ = s.saveOutgoingMessage(
			ctx,
			status.To,
			"text",
			messages[i].Text,
			nil,
			status.MessageID,
			true, // Bulk messages всегда billable
		)
	}

	return nil
}

// GetMessageStatus получает статус сообщения через Infobip
func (s *InfobipBotService) GetMessageStatus(ctx context.Context, messageID string) (*MessageStatus, error) {
	status, err := s.infobipClient.GetMessageStatus(ctx, messageID)
	if err != nil {
		return nil, err
	}

	return &MessageStatus{
		MessageID:   status.MessageID,
		Status:      status.Status.Name,
		Description: status.Status.Description,
	}, nil
}

// convertRichMediaToInfobip конвертирует внутренний формат в формат Infobip
func (s *InfobipBotService) convertRichMediaToInfobip(rm *models.RichMedia) *infobip.ViberRichMedia {
	var buttons []infobip.ViberRichMediaButton

	for _, btn := range rm.Buttons {
		buttons = append(buttons, infobip.ViberRichMediaButton{
			Columns:    btn.Columns,
			Rows:       btn.Rows,
			ActionType: btn.ActionType,
			ActionBody: btn.ActionBody,
			Text:       btn.Text,
			Image:      btn.Image,
			TextSize:   btn.TextSize,
			TextVAlign: btn.TextVAlign,
			TextHAlign: btn.TextHAlign,
			BgColor:    btn.BgColor,
		})
	}

	return &infobip.ViberRichMedia{
		ButtonsGroupColumns: rm.ButtonsGroupColumns,
		ButtonsGroupRows:    rm.ButtonsGroupRows,
		Buttons:             buttons,
	}
}

// richMediaToMap конвертирует RichMedia в map для сохранения в БД
func (s *InfobipBotService) richMediaToMap(rm *models.RichMedia) map[string]interface{} {
	buttons := make([]map[string]interface{}, len(rm.Buttons))
	for i, btn := range rm.Buttons {
		buttons[i] = map[string]interface{}{
			"Columns":    btn.Columns,
			"Rows":       btn.Rows,
			"ActionType": btn.ActionType,
			"ActionBody": btn.ActionBody,
			"Image":      btn.Image,
			"Text":       btn.Text,
			"TextSize":   btn.TextSize,
			"TextVAlign": btn.TextVAlign,
			"TextHAlign": btn.TextHAlign,
			"BgColor":    btn.BgColor,
		}
	}

	return map[string]interface{}{
		"Type":                rm.Type,
		"ButtonsGroupColumns": rm.ButtonsGroupColumns,
		"ButtonsGroupRows":    rm.ButtonsGroupRows,
		"BgColor":             rm.BgColor,
		"Buttons":             buttons,
	}
}

// generateStaticMapURL генерирует URL для статической карты
func (s *InfobipBotService) generateStaticMapURL(delivery *DeliveryInfo) string {
	mapboxToken := os.Getenv("MAPBOX_ACCESS_TOKEN")
	if mapboxToken == "" {
		// Mapbox токен требуется для создания статичных карт
		// Установите MAPBOX_ACCESS_TOKEN в переменных окружения
		return ""
	}

	// Маркеры и путь
	courierMarker := fmt.Sprintf("pin-l-bicycle+3b82f6(%f,%f)",
		delivery.CourierLongitude, delivery.CourierLatitude)
	deliveryMarker := fmt.Sprintf("pin-l-home+ef4444(%f,%f)",
		delivery.DeliveryLongitude, delivery.DeliveryLatitude)

	return fmt.Sprintf(
		"https://api.mapbox.com/styles/v1/mapbox/streets-v11/static/%s,%s/auto/600x400@2x?access_token=%s",
		courierMarker, deliveryMarker, mapboxToken,
	)
}

// saveOutgoingMessage сохраняет исходящее сообщение в БД
func (s *InfobipBotService) saveOutgoingMessage(ctx context.Context, viberID, msgType, text string, richMedia map[string]interface{}, messageID string, isBillable bool) error {
	query := `
		INSERT INTO viber_messages (
			viber_user_id, direction, message_type, text,
			rich_media, message_id, is_billable
		) 
		SELECT id, 'outgoing', $2, $3, $4, $5, $6
		FROM viber_users
		WHERE viber_id = $1
	`

	_, err := s.db.ExecContext(ctx, query, viberID, msgType, text, richMedia, messageID, isBillable)
	return err
}

// BulkMessageRequest запрос для массовой рассылки
type BulkMessageRequest struct {
	To      string
	Text    string
	IsPromo bool
}

// MessageStatus статус сообщения
type MessageStatus struct {
	MessageID   string
	Status      string
	Description string
}

// ProcessWebhook обрабатывает webhook от Infobip
func (s *InfobipBotService) ProcessWebhook(ctx context.Context, webhook *infobip.ViberWebhook) error {
	// Обновляем статус сообщения в БД
	query := `
		UPDATE viber_messages
		SET status = $2,
		    updated_at = CURRENT_TIMESTAMP
		WHERE message_id = $1
	`

	_, err := s.db.ExecContext(ctx, query, webhook.MessageID, webhook.Status.Name)
	if err != nil {
		return fmt.Errorf("failed to update message status: %w", err)
	}

	// Если это входящее сообщение
	if webhook.InboundContent != nil {
		return s.processInboundMessage(ctx, webhook)
	}

	// Если есть информация о цене, обновляем её
	if webhook.Price != nil {
		priceQuery := `
			UPDATE viber_messages
			SET price = $2,
			    currency = $3
			WHERE message_id = $1
		`
		_, err = s.db.ExecContext(ctx, priceQuery, webhook.MessageID, webhook.Price.PricePerMessage, webhook.Price.Currency)
	}

	return err
}

// processInboundMessage обрабатывает входящее сообщение
func (s *InfobipBotService) processInboundMessage(ctx context.Context, webhook *infobip.ViberWebhook) error {
	if webhook.InboundContent == nil {
		return nil
	}

	// Сохраняем информацию о пользователе
	user := &models.ViberSender{
		ID:   webhook.From,
		Name: webhook.From, // Infobip не предоставляет имя в webhook
	}

	if err := s.sessionManager.SaveUserInfo(ctx, user); err != nil {
		return fmt.Errorf("failed to save user info: %w", err)
	}

	// Создаём или обновляем сессию
	session, err := s.sessionManager.GetActiveSession(ctx, webhook.From)
	if err != nil && !errors.Is(err, ErrNoActiveSession) {
		return err
	}

	if errors.Is(err, ErrNoActiveSession) {
		session, err = s.sessionManager.CreateSession(ctx, webhook.From)
		if err != nil {
			return err
		}
	} else {
		if err := s.sessionManager.UpdateSession(ctx, session.ID); err != nil {
			return err
		}
	}

	// Сохраняем входящее сообщение в БД
	query := `
		INSERT INTO viber_messages (
			viber_user_id, viber_session_id, direction,
			message_type, text, tracking_data
		) 
		SELECT vu.id, $2, 'incoming', $3, $4, $5
		FROM viber_users vu
		WHERE vu.viber_id = $1
	`

	_, err = s.db.ExecContext(ctx, query,
		webhook.From,
		session.ID,
		webhook.InboundContent.Type,
		webhook.InboundContent.Text,
		webhook.InboundContent.TrackingData,
	)

	return err
}

// GetActiveSessionsCount возвращает количество активных сессий
func (s *InfobipBotService) GetActiveSessionsCount(ctx context.Context) (int, error) {
	var count int
	query := `
		SELECT COUNT(*) 
		FROM viber_sessions
		WHERE active = true AND expires_at > CURRENT_TIMESTAMP
	`

	err := s.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

// EstimateMessageCost оценивает стоимость отправки сообщения
func (s *InfobipBotService) EstimateMessageCost(ctx context.Context, viberID string, isRichMedia bool) (float64, error) {
	// Проверяем, есть ли активная сессия
	session, _ := s.sessionManager.GetActiveSession(ctx, viberID)

	if session != nil {
		// Сообщение в рамках сессии - бесплатно
		return 0, nil
	}

	// Стоимость за пределами сессии (примерные цены для Сербии)
	if isRichMedia {
		return 0.025, nil // ~2.5 цента за Rich Media
	}
	return 0.015, nil // ~1.5 цента за текстовое сообщение
}
