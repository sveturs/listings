// backend/internal/domain/models/user_profile.go
package models

import (
    "encoding/json"
    "time"
	"fmt"
)

// UserProfile расширяет базовую модель User дополнительными полями профиля
type UserProfile struct {
    User
    Phone            *string          `json:"phone,omitempty"`
    Bio             *string          `json:"bio,omitempty"`
    NotificationEmail bool            `json:"notification_email"`
    NotificationPush bool            `json:"notification_push"`
    Timezone        string           `json:"timezone"`
    LastSeen        *time.Time       `json:"last_seen,omitempty"`
    AccountStatus   string           `json:"account_status"`
    Settings        json.RawMessage  `json:"settings,omitempty"`
}

// UserProfileUpdate используется для частичного обновления профиля
type UserProfileUpdate struct {
    Phone            *string          `json:"phone,omitempty"`
    Bio             *string          `json:"bio,omitempty"`
    NotificationEmail *bool           `json:"notification_email,omitempty"`
    NotificationPush *bool           `json:"notification_push,omitempty"`
    Timezone        *string          `json:"timezone,omitempty"`
    Settings        json.RawMessage  `json:"settings,omitempty"`
}

// Validate проверяет корректность данных профиля
func (up *UserProfileUpdate) Validate() error {
    if up.Phone != nil && len(*up.Phone) > 20 {
        return fmt.Errorf("phone number is too long")
    }
    if up.Bio != nil && len(*up.Bio) > 1000 {
        return fmt.Errorf("bio is too long")
    }
    if up.Timezone != nil && !isValidTimezone(*up.Timezone) {
        return fmt.Errorf("invalid timezone")
    }
    return nil
}

// isValidTimezone проверяет существование часового пояса
func isValidTimezone(tz string) bool {
    _, err := time.LoadLocation(tz)
    return err == nil
}