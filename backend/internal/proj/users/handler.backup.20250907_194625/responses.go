// backend/internal/proj/users/handler/responses.go
package handler

import "backend/internal/domain/models"

// Request structures

// LoginRequest represents login request data
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=6" example:"password123"`
}

// RegisterRequest represents registration request data
type RegisterRequest struct {
	Name     string  `json:"name" validate:"required,min=2" example:"John Doe"`
	Email    string  `json:"email" validate:"required,email" example:"user@example.com"`
	Password string  `json:"password" validate:"required,min=6" example:"password123"`
	Phone    *string `json:"phone,omitempty" example:"+1234567890"`
}

// Response structures

// AuthResponse represents authentication response with access token
type AuthResponse struct {
	AccessToken string       `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	TokenType   string       `json:"token_type" example:"Bearer"`
	ExpiresIn   int          `json:"expires_in" example:"3600"`
	User        UserResponse `json:"user"`
}

// UserResponse represents user data in responses
type UserResponse struct {
	ID         int    `json:"id" example:"1"`
	Name       string `json:"name" example:"John Doe"`
	Email      string `json:"email" example:"user@example.com"`
	PictureURL string `json:"picture_url,omitempty" example:"https://example.com/avatar.jpg"`
}

// TokenResponse represents token refresh response
type TokenResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	TokenType   string `json:"token_type" example:"Bearer"`
	ExpiresIn   int    `json:"expires_in" example:"3600"`
}

// SessionResponse represents session information
type SessionResponse struct {
	Authenticated bool                 `json:"authenticated" example:"true"`
	User          *SessionUserResponse `json:"user,omitempty"`
}

// SessionUserResponse represents user data in session response
type SessionUserResponse struct {
	ID         int    `json:"id" example:"1"`
	Name       string `json:"name" example:"John Doe"`
	Email      string `json:"email" example:"user@example.com"`
	Provider   string `json:"provider" example:"password"`
	PictureURL string `json:"picture_url,omitempty" example:"https://example.com/avatar.jpg"`
	IsAdmin    bool   `json:"is_admin" example:"false"`
	City       string `json:"city,omitempty" example:"Moscow"`
	Country    string `json:"country,omitempty" example:"Russia"`
	Phone      string `json:"phone,omitempty" example:"+1234567890"`
}

// AdminUserListResponse represents paginated list of users
type AdminUserListResponse struct {
	Data  []*models.UserProfile `json:"data"`
	Total int                   `json:"total" example:"100"`
	Page  int                   `json:"page" example:"1"`
	Limit int                   `json:"limit" example:"10"`
	Pages int                   `json:"pages" example:"10"`
}

// AdminStatusUpdateRequest represents user status update request
type AdminStatusUpdateRequest struct {
	Status string `json:"status" validate:"required,oneof=active blocked pending" example:"active"`
}

// AdminMessageResponse represents success message response
type AdminMessageResponse struct {
	Message string `json:"message" example:"admin.users.success.profile_updated"`
}

// AdminBalanceResponse wraps balance data
type AdminBalanceResponse struct {
	Data interface{} `json:"data"`
}

// AdminAdminsResponse represents admin check response
type AdminAdminsResponse struct {
	Email   string `json:"email" example:"admin@example.com"`
	IsAdmin bool   `json:"is_admin" example:"true"`
}

// LoginResponse represents deprecated login response
type LoginResponse struct {
	Message string       `json:"message" example:"users.login.success.authenticated"`
	User    *models.User `json:"user"`
}

// UsersListResponse структура ответа со списком пользователей для админки
type UsersListResponse struct {
	Users      []*models.UserProfile `json:"users"`
	TotalCount int                   `json:"total_count" example:"100"`
	Page       int                   `json:"page" example:"1"`
	PageSize   int                   `json:"page_size" example:"20"`
}

// UpdateUserStatusRequest структура запроса изменения статуса пользователя
type UpdateUserStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=active inactive suspended" example:"active"`
}

// UpdateUserRoleRequest структура запроса изменения роли пользователя
type UpdateUserRoleRequest struct {
	RoleID int `json:"role_id" validate:"required,min=1" example:"2"`
}
