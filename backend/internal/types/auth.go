package types

import (
	"errors"
	"golang.org/x/oauth2"
)

type AuthProvider string

const (
	ProviderGoogle   AuthProvider = "google"
	ProviderEmail    AuthProvider = "email"
	ProviderPassword AuthProvider = "password"
	ProviderJWT      AuthProvider = "jwt"
)

type SessionData struct {
    Token      *oauth2.Token `json:"-"`
    UserID     int          `json:"user_id"`
    Name       string       `json:"name"`
    Email      string       `json:"email"`
    GoogleID   string       `json:"google_id"`
    PictureURL string       `json:"picture_url"`
    Provider   string       `json:"provider"`
}

// Ошибки аутентификации
var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrInvalidToken       = errors.New("invalid authentication token")
	ErrTokenExpired       = errors.New("authentication token has expired")
	ErrUserNotFound       = errors.New("user not found")
)