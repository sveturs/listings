package service

import (
	"backend/internal/domain/models"
 	"backend/internal/types"
	"context"
	"sync"
    userStorage "backend/internal/proj/users/storage"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauth2v2 "google.golang.org/api/oauth2/v2"
)

type AuthService struct {
    googleConfig *oauth2.Config
    sessions     sync.Map
    storage      userStorage.UserStorage
}

func NewAuthService(
    googleClientID, 
    googleClientSecret, 
    googleRedirectURL string, 
    storage userStorage.UserStorage,
) *AuthService {
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

    return &AuthService{
        googleConfig: googleConfig,
        storage:     storage,
    }
}

func (s *AuthService) GetGoogleAuthURL() string {
	return s.googleConfig.AuthCodeURL(
		"state",
		oauth2.SetAuthURLParam("prompt", "select_account"),
		oauth2.SetAuthURLParam("access_type", "offline"),
		oauth2.SetAuthURLParam("include_granted_scopes", "true"),
	)
}

func (s *AuthService) HandleGoogleCallback(ctx context.Context, code string) (*types.SessionData, error) {
	token, err := s.googleConfig.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	oauth2Service, err := oauth2v2.New(s.googleConfig.Client(ctx, token))
	if err != nil {
		return nil, err
	}

	userInfo, err := oauth2Service.Userinfo.Get().Do()
	if err != nil {
		return nil, err
	}

	user, err := s.storage.GetOrCreateGoogleUser(ctx, &models.User{
		Name:       userInfo.Name,
		Email:      userInfo.Email,
		GoogleID:   userInfo.Id,
		PictureURL: userInfo.Picture,
	})
	if err != nil {
		return nil, err
	}

	return &types.SessionData{
		Token:      token,
		UserID:     user.ID,
		Name:       user.Name,
		Email:      user.Email,
		GoogleID:   user.GoogleID,
		PictureURL: user.PictureURL,
		Provider:   "google",
	}, nil
}

func (s *AuthService) SaveSession(token string, data *types.SessionData) {
	s.sessions.Store(token, data)
}

func (s *AuthService) GetSession(token string) (*types.SessionData, bool) {
	if value, ok := s.sessions.Load(token); ok {
		if sessionData, ok := value.(*types.SessionData); ok {
			return sessionData, true
		}
	}
	return nil, false
}

func (s *AuthService) DeleteSession(token string) {
	s.sessions.Delete(token)
}
