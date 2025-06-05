package models

import "time"

// UserContact представляет контакт между пользователями
type UserContact struct {
	ID            int       `json:"id" db:"id"`
	UserID        int       `json:"user_id" db:"user_id"`
	ContactUserID int       `json:"contact_user_id" db:"contact_user_id"`
	Status        string    `json:"status" db:"status"` // pending, accepted, blocked
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	
	// Дополнительные поля
	AddedFromChatID *int   `json:"added_from_chat_id,omitempty" db:"added_from_chat_id"`
	Notes           string `json:"notes,omitempty" db:"notes"`
	
	// Связанные объекты
	User        *User `json:"user,omitempty"`
	ContactUser *User `json:"contact_user,omitempty"`
}

// UserPrivacySettings настройки приватности пользователя
type UserPrivacySettings struct {
	UserID                        int       `json:"user_id" db:"user_id"`
	AllowContactRequests          bool      `json:"allow_contact_requests" db:"allow_contact_requests"`
	AllowMessagesFromContactsOnly bool      `json:"allow_messages_from_contacts_only" db:"allow_messages_from_contacts_only"`
	CreatedAt                     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                     time.Time `json:"updated_at" db:"updated_at"`
}

// Константы для статусов контактов
const (
	ContactStatusPending  = "pending"
	ContactStatusAccepted = "accepted"
	ContactStatusBlocked  = "blocked"
)

// Структуры для запросов API
type AddContactRequest struct {
	ContactUserID   int    `json:"contact_user_id" validate:"required"`
	Notes           string `json:"notes,omitempty"`
	AddedFromChatID *int   `json:"added_from_chat_id,omitempty"`
}

type UpdateContactRequest struct {
	Status string `json:"status" validate:"required,oneof=accepted blocked"`
	Notes  string `json:"notes,omitempty"`
}

type ContactsListResponse struct {
	Contacts []UserContact `json:"contacts"`
	Total    int           `json:"total"`
	Page     int           `json:"page"`
	Limit    int           `json:"limit"`
}

type UpdatePrivacySettingsRequest struct {
	AllowContactRequests          *bool `json:"allow_contact_requests,omitempty"`
	AllowMessagesFromContactsOnly *bool `json:"allow_messages_from_contacts_only,omitempty"`
}