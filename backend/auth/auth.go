package auth

import (
	"context"
	"sync"
"log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauth2v2 "google.golang.org/api/oauth2/v2"
)

type SessionData struct {
    Token      *oauth2.Token
    UserID     int    `json:"user_id"`
    Name       string `json:"name"`
    Email      string `json:"email"`
    GoogleID   string `json:"google_id"`
    PictureURL string `json:"picture_url"`
    Provider   string `json:"provider"`
}

type AuthManager struct {
	googleConfig *oauth2.Config
	sessions     sync.Map
}

func NewAuthManager(googleClientID, googleClientSecret, googleRedirectURL string) *AuthManager {
    log.Printf("Initializing AuthManager with: ClientID=%s, RedirectURL=%s", 
        googleClientID, googleRedirectURL)


	googleConfig := &oauth2.Config{
		ClientID:     googleClientID,
		ClientSecret: googleClientSecret,
		RedirectURL:  googleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}

	return &AuthManager{
		googleConfig: googleConfig,
	}
}


func (am *AuthManager) GetGoogleAuthURL() string {
    url := am.googleConfig.AuthCodeURL("state")
    log.Printf("Generated Google Auth URL: %s", url)
    return url
}


func (am *AuthManager) HandleGoogleCallback(ctx context.Context, code string) (*SessionData, error) {
    token, err := am.googleConfig.Exchange(ctx, code)
    if err != nil {
        return nil, err
    }

    oauth2Service, err := oauth2v2.New(am.googleConfig.Client(ctx, token))
    if err != nil {
        return nil, err
    }

    userInfo, err := oauth2Service.Userinfo.Get().Do()
    if err != nil {
        return nil, err
    }

    return &SessionData{
        Token:      token,
        Name:       userInfo.Name,
        Email:      userInfo.Email,
        GoogleID:   userInfo.Id,
        PictureURL: userInfo.Picture,
        Provider:   "google",
    }, nil
}

func (am *AuthManager) SaveSession(token string, data *SessionData) {
	am.sessions.Store(token, data)
}

func (am *AuthManager) GetSession(token string) (*SessionData, bool) {
	if value, ok := am.sessions.Load(token); ok {
		if sessionData, ok := value.(*SessionData); ok {
			return sessionData, true
		}
	}
	return nil, false
}

func (am *AuthManager) DeleteSession(token string) {
	am.sessions.Delete(token)
}