// backend/internal/domain/models/admin.go
package models

import "time"

// AdminUser представляет запись администратора
type AdminUser struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy *int      `json:"created_by,omitempty"`
	Notes     *string   `json:"notes,omitempty"`
}
