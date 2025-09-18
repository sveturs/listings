package infobip

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client для работы с Infobip API
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewClient создаёт новый клиент Infobip
func NewClient(apiKey, baseURL string) *Client {
	return &Client{
		apiKey:     apiKey,
		baseURL:    "https://" + baseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// ViberMessage сообщение для отправки через Viber
type ViberMessage struct {
	From         string            `json:"from"`
	To           string            `json:"to"`
	Content      ViberContent      `json:"content"`
	CallbackData string            `json:"callbackData,omitempty"`
	NotifyURL    string            `json:"notifyUrl,omitempty"`
	Label        *ViberLabel       `json:"label,omitempty"`
	SessionInfo  *ViberSessionInfo `json:"sessionInfo,omitempty"`
}

// ViberContent содержимое сообщения
type ViberContent struct {
	Type         string          `json:"type"` // TEXT, IMAGE, VIDEO, FILE, BUTTON, RICH_MEDIA
	Text         string          `json:"text,omitempty"`
	ImageURL     string          `json:"imageUrl,omitempty"`
	FileURL      string          `json:"fileUrl,omitempty"`
	ButtonText   string          `json:"buttonText,omitempty"`
	ButtonURL    string          `json:"buttonUrl,omitempty"`
	TrackingData string          `json:"trackingData,omitempty"`
	RichMedia    *ViberRichMedia `json:"richMedia,omitempty"`
}

// ViberRichMedia для rich media сообщений
type ViberRichMedia struct {
	ButtonsGroupColumns int                    `json:"buttonsGroupColumns"`
	ButtonsGroupRows    int                    `json:"buttonsGroupRows"`
	Buttons             []ViberRichMediaButton `json:"buttons"`
}

// ViberRichMediaButton кнопка в rich media
type ViberRichMediaButton struct {
	Columns    int    `json:"columns"`
	Rows       int    `json:"rows"`
	ActionType string `json:"actionType"` // reply, open-url, none
	ActionBody string `json:"actionBody,omitempty"`
	Text       string `json:"text,omitempty"`
	Image      string `json:"image,omitempty"`
	TextSize   string `json:"textSize,omitempty"`
	TextVAlign string `json:"textVAlign,omitempty"`
	TextHAlign string `json:"textHAlign,omitempty"`
	BgColor    string `json:"bgColor,omitempty"`
}

// ViberLabel метка для билируемых сообщений
type ViberLabel struct {
	Type    string `json:"type"` // PROMOTION, TRANSACTION
	Content string `json:"content,omitempty"`
}

// ViberSessionInfo информация о сессии для бесплатных сообщений
type ViberSessionInfo struct {
	SessionID string `json:"sessionId,omitempty"`
	Origin    string `json:"origin,omitempty"` // USER_INITIATED, BUSINESS_INITIATED
}

// ViberBulkMessage массовая рассылка
type ViberBulkMessage struct {
	Messages []ViberMessage `json:"messages"`
}

// ViberResponse ответ от API
type ViberResponse struct {
	Messages []ViberMessageStatus `json:"messages"`
	BulkID   string               `json:"bulkId,omitempty"`
}

// ViberMessageStatus статус отправленного сообщения
type ViberMessageStatus struct {
	To        string      `json:"to"`
	Status    ViberStatus `json:"status"`
	MessageID string      `json:"messageId"`
}

// ViberStatus статус сообщения
type ViberStatus struct {
	GroupID     int    `json:"groupId"`
	GroupName   string `json:"groupName"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ViberWebhook структура для вебхука
type ViberWebhook struct {
	MessageID      string               `json:"messageId"`
	To             string               `json:"to"`
	From           string               `json:"from"`
	SentAt         string               `json:"sentAt"`
	DoneAt         string               `json:"doneAt,omitempty"`
	Status         ViberStatus          `json:"status"`
	Price          *ViberPrice          `json:"price,omitempty"`
	Error          *ViberError          `json:"error,omitempty"`
	CallbackData   string               `json:"callbackData,omitempty"`
	InboundContent *ViberInboundContent `json:"content,omitempty"`
}

// ViberPrice стоимость сообщения
type ViberPrice struct {
	PricePerMessage float64 `json:"pricePerMessage"`
	Currency        string  `json:"currency"`
}

// ViberError ошибка доставки
type ViberError struct {
	GroupID     int    `json:"groupId"`
	GroupName   string `json:"groupName"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ViberInboundContent входящее сообщение от пользователя
type ViberInboundContent struct {
	Type         string                `json:"type"`
	Text         string                `json:"text,omitempty"`
	Media        *ViberInboundMedia    `json:"media,omitempty"`
	Location     *ViberInboundLocation `json:"location,omitempty"`
	TrackingData string                `json:"trackingData,omitempty"`
}

// ViberInboundMedia медиа во входящем сообщении
type ViberInboundMedia struct {
	URL      string `json:"url"`
	FileName string `json:"fileName,omitempty"`
	Size     int64  `json:"size,omitempty"`
}

// ViberInboundLocation локация во входящем сообщении
type ViberInboundLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// SendTextMessage отправляет текстовое сообщение
func (c *Client) SendTextMessage(ctx context.Context, from, to, text string, sessionInfo *ViberSessionInfo) (*ViberResponse, error) {
	msg := ViberMessage{
		From: from,
		To:   to,
		Content: ViberContent{
			Type: "TEXT",
			Text: text,
		},
		SessionInfo: sessionInfo,
	}

	return c.sendMessage(ctx, msg)
}

// SendInteractiveMessage отправляет сообщение с интерактивными ссылками для карты и маркетплейса
func (c *Client) SendInteractiveMessage(ctx context.Context, from, to string, sessionInfo *ViberSessionInfo) (*ViberResponse, error) {
	text := `🎯 Доступные действия в SveTu:

📍 Карта с товарами:
https://svetu.rs/map

🏪 Маркетплейс:
https://svetu.rs

🚚 Отслеживание доставки:
https://svetu.rs/track/ABC123

✨ Все ссылки открываются прямо в Viber!`

	return c.SendTextMessage(ctx, from, to, text, sessionInfo)
}

// SendTrackingLink отправляет ссылку для отслеживания заказа
func (c *Client) SendTrackingLink(ctx context.Context, from, to, orderID string, sessionInfo *ViberSessionInfo) (*ViberResponse, error) {
	text := fmt.Sprintf(`🗺️ Откройте карту трекинга!

Для просмотра карты с местоположением курьера перейдите по ссылке:

https://svetu.rs/track/%s

Ссылка откроется прямо в Viber!`, orderID)

	msg := ViberMessage{
		From: from,
		To:   to,
		Content: ViberContent{
			Type: "TEXT",
			Text: text,
		},
		CallbackData: "tracking_link_sent",
		SessionInfo:  sessionInfo,
	}

	return c.sendMessage(ctx, msg)
}

// SendImageMessage отправляет изображение
func (c *Client) SendImageMessage(ctx context.Context, from, to, imageURL, text string, sessionInfo *ViberSessionInfo) (*ViberResponse, error) {
	msg := ViberMessage{
		From: from,
		To:   to,
		Content: ViberContent{
			Type:     "IMAGE",
			ImageURL: imageURL,
			Text:     text,
		},
		SessionInfo: sessionInfo,
	}

	return c.sendMessage(ctx, msg)
}

// SendButtonMessage отправляет сообщение с кнопкой
func (c *Client) SendButtonMessage(ctx context.Context, from, to, text, buttonText, buttonURL string, sessionInfo *ViberSessionInfo) (*ViberResponse, error) {
	msg := ViberMessage{
		From: from,
		To:   to,
		Content: ViberContent{
			Type:       "BUTTON",
			Text:       text,
			ButtonText: buttonText,
			ButtonURL:  buttonURL,
		},
		SessionInfo: sessionInfo,
	}

	return c.sendMessage(ctx, msg)
}

// SendRichMedia отправляет rich media сообщение
func (c *Client) SendRichMedia(ctx context.Context, from, to string, richMedia *ViberRichMedia, text string, sessionInfo *ViberSessionInfo) (*ViberResponse, error) {
	msg := ViberMessage{
		From: from,
		To:   to,
		Content: ViberContent{
			Type:      "RICH_MEDIA",
			Text:      text,
			RichMedia: richMedia,
		},
		SessionInfo: sessionInfo,
	}

	return c.sendMessage(ctx, msg)
}

// SendBulkMessages отправляет массовую рассылку
func (c *Client) SendBulkMessages(ctx context.Context, messages []ViberMessage) (*ViberResponse, error) {
	bulk := ViberBulkMessage{Messages: messages}
	return c.sendBulk(ctx, bulk)
}

// sendMessage отправляет одно сообщение через Infobip API v2
func (c *Client) sendMessage(ctx context.Context, msg ViberMessage) (*ViberResponse, error) {
	// Используем новый универсальный endpoint API v2
	url := fmt.Sprintf("%s/viber/2/messages", c.baseURL)

	// Формируем структуру для API v2
	v2Message := map[string]interface{}{
		"messages": []map[string]interface{}{
			{
				"sender":       msg.From,
				"destinations": []map[string]string{{"to": msg.To}},
				"content":      msg.Content,
			},
		},
	}

	body, err := json.Marshal(v2Message)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "App "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var viberResp ViberResponse
	if err := json.Unmarshal(respBody, &viberResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &viberResp, nil
}

// sendBulk отправляет массовую рассылку
func (c *Client) sendBulk(ctx context.Context, bulk ViberBulkMessage) (*ViberResponse, error) {
	url := fmt.Sprintf("%s/viber/1/send/bulk", c.baseURL)

	body, err := json.Marshal(bulk)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal bulk message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "App "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusMultiStatus {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var viberResp ViberResponse
	if err := json.Unmarshal(respBody, &viberResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &viberResp, nil
}

// GetMessageStatus получает статус сообщения
func (c *Client) GetMessageStatus(ctx context.Context, messageID string) (*ViberMessageStatus, error) {
	url := fmt.Sprintf("%s/viber/1/reports/%s", c.baseURL, messageID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "App "+c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var status ViberMessageStatus
	if err := json.Unmarshal(respBody, &status); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &status, nil
}
