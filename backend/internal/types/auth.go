package types

import "golang.org/x/oauth2"

type SessionData struct {
    Token      *oauth2.Token `json:"-"`
    UserID     int          `json:"user_id"`
    Name       string       `json:"name"`
    Email      string       `json:"email"`
    GoogleID   string       `json:"google_id"`
    PictureURL string       `json:"picture_url"`
    Provider   string       `json:"provider"`
}