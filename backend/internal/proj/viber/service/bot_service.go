package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"backend/internal/proj/viber/config"
	"backend/internal/proj/viber/models"
	"backend/internal/storage/postgres"
)

// BotService управляет Viber Bot
type BotService struct {
	config   *config.ViberConfig
	db       *postgres.Database
	client   *http.Client
	sessions *SessionManager
}

// NewBotService создаёт новый сервис бота
func NewBotService(cfg *config.ViberConfig, db *postgres.Database) *BotService {
	return &BotService{
		config:   cfg,
		db:       db,
		client:   &http.Client{Timeout: 10 * time.Second},
		sessions: NewSessionManager(db),
	}
}

// SetWebhook устанавливает webhook для бота
func (s *BotService) SetWebhook() error {
	url := fmt.Sprintf("%s/set_webhook", s.config.APIEndpoint)

	payload := map[string]interface{}{
		"auth_token": s.config.AuthToken,
		"url":        s.config.WebhookURL,
		"event_types": []string{
			"delivered",
			"seen",
			"failed",
			"subscribed",
			"unsubscribed",
			"conversation_started",
			"message",
		},
		"send_name":  true,
		"send_photo": true,
	}

	resp, err := s.makeRequest("POST", url, payload)
	if err != nil {
		return fmt.Errorf("failed to set webhook: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if status, ok := result["status"].(float64); !ok || status != 0 {
		return fmt.Errorf("webhook setup failed: %v", result["status_message"])
	}

	return nil
}

// SendTextMessage отправляет текстовое сообщение
func (s *BotService) SendTextMessage(ctx context.Context, viberID, text string) error {
	// Проверяем, есть ли активная сессия (для определения billable)
	session, _ := s.sessions.GetActiveSession(ctx, viberID)
	isBillable := session == nil

	message := &models.OutgoingMessage{
		Receiver: viberID,
		Type:     "text",
		Text:     text,
		Sender: models.OutgoingSender{
			Name:   s.config.BotName,
			Avatar: s.config.BotAvatar,
		},
	}

	if err := s.sendMessage(message); err != nil {
		return err
	}

	// Сохраняем в БД
	return s.saveOutgoingMessage(ctx, viberID, "text", text, nil, isBillable)
}

// SendRichMedia отправляет Rich Media сообщение
func (s *BotService) SendRichMedia(ctx context.Context, viberID string, richMedia *models.RichMedia, text string) error {
	session, _ := s.sessions.GetActiveSession(ctx, viberID)
	isBillable := session == nil

	message := &models.OutgoingMessage{
		Receiver:  viberID,
		Type:      "rich_media",
		Text:      text,
		RichMedia: s.richMediaToMap(richMedia),
		Sender: models.OutgoingSender{
			Name:   s.config.BotName,
			Avatar: s.config.BotAvatar,
		},
	}

	if err := s.sendMessage(message); err != nil {
		return err
	}

	return s.saveOutgoingMessage(ctx, viberID, "rich_media", text, s.richMediaToMap(richMedia), isBillable)
}

// SendTrackingNotification отправляет уведомление с трекингом
func (s *BotService) SendTrackingNotification(ctx context.Context, viberID string, delivery *DeliveryInfo) error {
	// Генерируем статическую карту
	mapURL := s.generateStaticMapURL(delivery)

	// Создаём Rich Media с картой и кнопкой трекинга
	richMedia := &models.RichMedia{
		Type:                "rich_media",
		ButtonsGroupColumns: 6,
		ButtonsGroupRows:    7,
		Buttons: []models.RichButton{
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
				Columns:    6,
				Rows:       2,
				ActionType: "open-url",
				ActionBody: fmt.Sprintf("%s/track/%s?viber=true", s.config.FrontendURL, delivery.TrackingToken),
				Text:       "🗺️ Отследить курьера",
				TextSize:   "large",
				TextVAlign: "middle",
				TextHAlign: "center",
				BgColor:    "#1976d2",
			},
		},
	}

	return s.SendRichMedia(ctx, viberID, richMedia, "Отслеживание доставки")
}

// SendKeyboard отправляет сообщение с клавиатурой
func (s *BotService) SendKeyboard(ctx context.Context, viberID, text string, keyboard *models.Keyboard) error {
	message := &models.OutgoingMessage{
		Receiver: viberID,
		Type:     "text",
		Text:     text,
		Keyboard: keyboard,
		Sender: models.OutgoingSender{
			Name:   s.config.BotName,
			Avatar: s.config.BotAvatar,
		},
	}

	return s.sendMessage(message)
}

// GetMainMenuKeyboard возвращает основную клавиатуру меню
func (s *BotService) GetMainMenuKeyboard() *models.Keyboard {
	return &models.Keyboard{
		Type:          "keyboard",
		DefaultHeight: true,
		Buttons: []models.Button{
			{
				Columns:    3,
				Rows:       1,
				ActionType: "reply",
				ActionBody: "search",
				Text:       "🔍 Поиск товаров",
				TextSize:   "regular",
				TextHAlign: "center",
				TextVAlign: "middle",
			},
			{
				Columns:    3,
				Rows:       1,
				ActionType: "reply",
				ActionBody: "my_orders",
				Text:       "📦 Мои заказы",
				TextSize:   "regular",
				TextHAlign: "center",
				TextVAlign: "middle",
			},
			{
				Columns:    2,
				Rows:       1,
				ActionType: "reply",
				ActionBody: "cart",
				Text:       "🛒 Корзина",
				TextSize:   "regular",
				TextHAlign: "center",
				TextVAlign: "middle",
			},
			{
				Columns:    2,
				Rows:       1,
				ActionType: "reply",
				ActionBody: "b2c_stores",
				Text:       "🏪 Витрины",
				TextSize:   "regular",
				TextHAlign: "center",
				TextVAlign: "middle",
			},
			{
				Columns:    2,
				Rows:       1,
				ActionType: "reply",
				ActionBody: "help",
				Text:       "❓ Помощь",
				TextSize:   "regular",
				TextHAlign: "center",
				TextVAlign: "middle",
			},
			{
				Columns:    6,
				Rows:       1,
				ActionType: "open-url",
				ActionBody: s.config.FrontendURL,
				Text:       "🌐 Открыть сайт",
				TextSize:   "regular",
				TextHAlign: "center",
				TextVAlign: "middle",
				BgColor:    "#e8f5e9",
			},
		},
	}
}

// sendMessage отправляет сообщение через Viber API
func (s *BotService) sendMessage(message *models.OutgoingMessage) error {
	url := fmt.Sprintf("%s/send_message", s.config.APIEndpoint)

	// Добавляем auth token
	payload := map[string]interface{}{
		"auth_token": s.config.AuthToken,
		"receiver":   message.Receiver,
		"type":       message.Type,
		"sender": map[string]interface{}{
			"name":   message.Sender.Name,
			"avatar": message.Sender.Avatar,
		},
	}

	if message.Text != "" {
		payload["text"] = message.Text
	}

	if message.RichMedia != nil {
		payload["rich_media"] = message.RichMedia
	}

	if message.Keyboard != nil {
		payload["keyboard"] = message.Keyboard
	}

	resp, err := s.makeRequest("POST", url, payload)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if status, ok := result["status"].(float64); !ok || status != 0 {
		return fmt.Errorf("message send failed: %v", result["status_message"])
	}

	return nil
}

// makeRequest выполняет HTTP запрос к Viber API
func (s *BotService) makeRequest(method, url string, payload interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// generateStaticMapURL генерирует URL для статической карты
func (s *BotService) generateStaticMapURL(delivery *DeliveryInfo) string {
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

// richMediaToMap конвертирует RichMedia в map для JSON
func (s *BotService) richMediaToMap(rm *models.RichMedia) map[string]interface{} {
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

// saveOutgoingMessage сохраняет исходящее сообщение в БД
func (s *BotService) saveOutgoingMessage(ctx context.Context, viberID, msgType, text string, richMedia map[string]interface{}, isBillable bool) error {
	// TODO: Implement database save
	return nil
}

// DeliveryInfo содержит информацию о доставке для уведомления
type DeliveryInfo struct {
	ID                int
	OrderID           int
	TrackingToken     string
	CourierName       string
	CourierLatitude   float64
	CourierLongitude  float64
	DeliveryLatitude  float64
	DeliveryLongitude float64
	EstimatedTime     time.Time
	DeliveryAddress   string
}
