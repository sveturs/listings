package models

import (
	"time"
)

// RefreshToken представляет refresh токен пользователя
type RefreshToken struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	// Дополнительная информация для безопасности
	UserAgent string `json:"user_agent" db:"user_agent"`
	IP        string `json:"ip" db:"ip"`

	// Для управления сессиями
	DeviceName string     `json:"device_name,omitempty" db:"device_name"`
	IsRevoked  bool       `json:"is_revoked" db:"is_revoked"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty" db:"revoked_at"`
}

// IsValid проверяет, действителен ли refresh токен
func (rt *RefreshToken) IsValid() bool {
	if rt.IsRevoked {
		return false
	}
	return time.Now().Before(rt.ExpiresAt)
}
