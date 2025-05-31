package models

// ErrorResponse представляет стандартный ответ с ошибкой
// @name ErrorResponse
type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid request"`
	Message string `json:"message,omitempty" example:"Validation failed"`
	Code    int    `json:"code,omitempty" example:"400"`
} // @name ErrorResponse

// SuccessResponse представляет стандартный успешный ответ
// @name SuccessResponse
type SuccessResponse struct {
	Success bool        `json:"success" example:"true"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty" example:"Operation completed successfully"`
} // @name SuccessResponse

// User представляет пользователя для API
// @name User
type UserAPI struct {
	ID       int    `json:"id" example:"1"`
	Email    string `json:"email" example:"user@example.com"`
	Name     string `json:"name" example:"Иван Иванов"`
	Role     string `json:"role" example:"user"`
	IsActive bool   `json:"is_active" example:"true"`
} // @name User
