package models

import (
	"encoding/json"
	"time"
)

// ViberUser представляет пользователя Viber
type ViberUser struct {
	ID                    int        `json:"id" db:"id"`
	ViberID               string     `json:"viber_id" db:"viber_id"`
	UserID                *int       `json:"user_id" db:"user_id"`
	Name                  string     `json:"name" db:"name"`
	AvatarURL             string     `json:"avatar_url" db:"avatar_url"`
	Language              string     `json:"language" db:"language"`
	CountryCode           string     `json:"country_code" db:"country_code"`
	APIVersion            int        `json:"api_version" db:"api_version"`
	Subscribed            bool       `json:"subscribed" db:"subscribed"`
	SubscribedAt          *time.Time `json:"subscribed_at" db:"subscribed_at"`
	LastSessionAt         *time.Time `json:"last_session_at" db:"last_session_at"`
	ConversationStartedAt *time.Time `json:"conversation_started_at" db:"conversation_started_at"`
	CreatedAt             time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at" db:"updated_at"`
}

// ViberSession представляет 24-часовую сессию
type ViberSession struct {
	ID            int             `json:"id" db:"id"`
	ViberUserID   int             `json:"viber_user_id" db:"viber_user_id"`
	StartedAt     time.Time       `json:"started_at" db:"started_at"`
	LastMessageAt time.Time       `json:"last_message_at" db:"last_message_at"`
	ExpiresAt     time.Time       `json:"expires_at" db:"expires_at"`
	MessageCount  int             `json:"message_count" db:"message_count"`
	Context       json.RawMessage `json:"context" db:"context"`
	Active        bool            `json:"active" db:"active"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
}

// ViberMessage представляет сообщение
type ViberMessage struct {
	ID           int                    `json:"id" db:"id"`
	ViberUserID  int                    `json:"viber_user_id" db:"viber_user_id"`
	SessionID    *int                   `json:"session_id" db:"session_id"`
	MessageToken string                 `json:"message_token" db:"message_token"`
	Direction    string                 `json:"direction" db:"direction"` // incoming, outgoing
	MessageType  string                 `json:"message_type" db:"message_type"`
	Content      string                 `json:"content" db:"content"`
	RichMedia    map[string]interface{} `json:"rich_media" db:"rich_media"`
	IsBillable   bool                   `json:"is_billable" db:"is_billable"`
	Status       string                 `json:"status" db:"status"`
	CreatedAt    time.Time              `json:"created_at" db:"created_at"`
}

// Webhook Events from Viber
type WebhookEvent struct {
	Event        string            `json:"event"`
	Timestamp    int64             `json:"timestamp"`
	ChatHostname string            `json:"chat_hostname"`
	MessageToken string            `json:"message_token"`
	Sender       *ViberSender      `json:"sender,omitempty"`
	Message      *ViberMessageData `json:"message,omitempty"`
	User         *ViberSender      `json:"user,omitempty"`
	Context      string            `json:"context,omitempty"`
	Silent       bool              `json:"silent"`
}

type ViberSender struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Avatar     string `json:"avatar,omitempty"`
	Country    string `json:"country,omitempty"`
	Language   string `json:"language,omitempty"`
	APIVersion int    `json:"api_version,omitempty"`
}

type ViberMessageData struct {
	Type         string                 `json:"type"`
	Text         string                 `json:"text,omitempty"`
	Media        string                 `json:"media,omitempty"`
	Location     *ViberLocation         `json:"location,omitempty"`
	Contact      *ViberContact          `json:"contact,omitempty"`
	TrackingData string                 `json:"tracking_data,omitempty"`
	FileName     string                 `json:"file_name,omitempty"`
	FileSize     int                    `json:"file_size,omitempty"`
	Duration     int                    `json:"duration,omitempty"`
	RichMedia    map[string]interface{} `json:"rich_media,omitempty"`
}

type ViberLocation struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

type ViberContact struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

// Outgoing message structures
type OutgoingMessage struct {
	Receiver      string                 `json:"receiver"`
	MinAPIVersion int                    `json:"min_api_version,omitempty"`
	Sender        OutgoingSender         `json:"sender"`
	TrackingData  string                 `json:"tracking_data,omitempty"`
	Type          string                 `json:"type"`
	Text          string                 `json:"text,omitempty"`
	Media         string                 `json:"media,omitempty"`
	Thumbnail     string                 `json:"thumbnail,omitempty"`
	Size          int                    `json:"size,omitempty"`
	Duration      int                    `json:"duration,omitempty"`
	RichMedia     map[string]interface{} `json:"rich_media,omitempty"`
	Keyboard      *Keyboard              `json:"keyboard,omitempty"`
}

type OutgoingSender struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar,omitempty"`
}

type Keyboard struct {
	Type          string   `json:"Type"`
	DefaultHeight bool     `json:"DefaultHeight,omitempty"`
	Buttons       []Button `json:"Buttons"`
}

type Button struct {
	Columns     int    `json:"Columns,omitempty"`
	Rows        int    `json:"Rows,omitempty"`
	BgColor     string `json:"BgColor,omitempty"`
	Silent      bool   `json:"Silent,omitempty"`
	BgMediaType string `json:"BgMediaType,omitempty"`
	BgMedia     string `json:"BgMedia,omitempty"`
	BgLoop      bool   `json:"BgLoop,omitempty"`
	ActionType  string `json:"ActionType,omitempty"`
	ActionBody  string `json:"ActionBody,omitempty"`
	Image       string `json:"Image,omitempty"`
	Text        string `json:"Text,omitempty"`
	TextVAlign  string `json:"TextVAlign,omitempty"`
	TextHAlign  string `json:"TextHAlign,omitempty"`
	TextOpacity int    `json:"TextOpacity,omitempty"`
	TextSize    string `json:"TextSize,omitempty"`
}

// Rich Media structures
type RichMedia struct {
	Type                string       `json:"Type"`
	ButtonsGroupColumns int          `json:"ButtonsGroupColumns"`
	ButtonsGroupRows    int          `json:"ButtonsGroupRows"`
	BgColor             string       `json:"BgColor,omitempty"`
	Buttons             []RichButton `json:"Buttons"`
}

type RichButton struct {
	Columns    int    `json:"Columns"`
	Rows       int    `json:"Rows"`
	ActionType string `json:"ActionType"`
	ActionBody string `json:"ActionBody,omitempty"`
	Image      string `json:"Image,omitempty"`
	Text       string `json:"Text,omitempty"`
	TextSize   string `json:"TextSize,omitempty"`
	TextVAlign string `json:"TextVAlign,omitempty"`
	TextHAlign string `json:"TextHAlign,omitempty"`
	BgColor    string `json:"BgColor,omitempty"`
}
