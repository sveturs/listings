// Package handler
// backend/internal/proj/storefront/handler/responses.go
package handler

// MessageResponse стандартный ответ с сообщением
type MessageResponse struct {
	Message string `json:"message"`
}

// CategoryMappingsUpdateResponse ответ на обновление сопоставлений категорий
type CategoryMappingsUpdateResponse struct {
	Message string `json:"message"`
}

// ApplyCategoryMappingsResponse ответ на применение сопоставлений категорий
type ApplyCategoryMappingsResponse struct {
	Message      string `json:"message"`
	UpdatedCount int    `json:"updated_count"`
}
