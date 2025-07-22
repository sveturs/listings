// backend/internal/proj/users/service/auth.go
package service

import (
	"context"
	"log"
	"sync"
	"time"

	"backend/internal/domain/models"
	"backend/internal/logger"
	"backend/internal/storage"
	"backend/internal/types"
	"backend/pkg/jwt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	oauth2v2 "google.golang.org/api/oauth2/v2"
)

type AuthService struct {
	googleConfig *oauth2.Config
	sessions     sync.Map
	storage      storage.Storage
	jwtSecret    string
	jwtExpHours  int
}

func NewAuthService(
	googleClientID string,
	googleClientSecret string,
	googleRedirectURL string,
	storage storage.Storage,
	jwtSecret string,
	jwtExpHours int,
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
		storage:      storage,
		jwtSecret:    jwtSecret,
		jwtExpHours:  jwtExpHours,
	}
}

func (s *AuthService) GetGoogleAuthURL(origin string) string {
	// Используем origin как state для последующего редиректа
	state := origin
	if state == "" {
		state = "default"
	}
	return s.googleConfig.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("prompt", "select_account"),
		oauth2.SetAuthURLParam("access_type", "offline"),
		oauth2.SetAuthURLParam("include_granted_scopes", "true"),
	)
}

func (s *AuthService) HandleGoogleCallback(ctx context.Context, code string) (*types.SessionData, error) {
	logger.Info().
		Str("code_prefix", code[:10]+"...").
		Msg("HandleGoogleCallback: exchanging code for token")

	token, err := s.googleConfig.Exchange(ctx, code)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("HandleGoogleCallback: failed to exchange code for token")
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
	log.Printf("AuthService: Session saved - UserID: %d, Email: %s, Provider: %s",
		data.UserID, data.Email, data.Provider)
}

func (s *AuthService) GetSession(ctx context.Context, token string) (*types.SessionData, error) {
	if value, ok := s.sessions.Load(token); ok {
		if sessionData, ok := value.(*types.SessionData); ok {
			log.Printf("AuthService: Session found - UserID: %d, Email: %s, Provider: %s",
				sessionData.UserID, sessionData.Email, sessionData.Provider)
			return sessionData, nil
		}
	}

	log.Printf("AuthService: Session not found")
	return nil, nil
}

func (s *AuthService) DeleteSession(token string) {
	s.sessions.Delete(token)
}

// GenerateJWT генерирует JWT токен для пользователя
func (s *AuthService) GenerateJWT(userID int, email string) (string, error) {
	return jwt.GenerateTokenWithDuration(
		userID,
		email,
		s.jwtSecret,
		time.Duration(s.jwtExpHours)*time.Hour,
	)
}

// ValidateJWT проверяет JWT токен и возвращает claims
func (s *AuthService) ValidateJWT(tokenString string) (*jwt.Claims, error) {
	return jwt.ValidateToken(tokenString, s.jwtSecret)
}

// LoginWithEmailPassword аутентификация по email и паролю с выдачей JWT
func (s *AuthService) LoginWithEmailPassword(ctx context.Context, email, password string) (string, *models.User, error) {
	// Получаем пользователя по email
	user, err := s.storage.GetUserByEmail(ctx, email)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, types.ErrInvalidCredentials
	}

	// Проверяем пароль
	if !user.CheckPassword(password) {
		return "", nil, types.ErrInvalidCredentials
	}

	// Генерируем JWT токен
	token, err := s.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

// RegisterWithEmailPassword регистрация нового пользователя с выдачей JWT
func (s *AuthService) RegisterWithEmailPassword(ctx context.Context, name, email, password string) (string, *models.User, error) {
	// Проверяем существует ли пользователь
	existingUser, _ := s.storage.GetUserByEmail(ctx, email)
	if existingUser != nil {
		return "", nil, types.ErrUserAlreadyExists
	}

	// Создаем нового пользователя
	user := &models.User{
		Name:  name,
		Email: email,
	}

	// Хешируем пароль
	if err := user.SetPassword(password); err != nil {
		return "", nil, err
	}

	// Сохраняем в БД
	savedUser, err := s.storage.CreateUser(ctx, user)
	if err != nil {
		return "", nil, err
	}

	// Генерируем JWT токен
	token, err := s.GenerateJWT(savedUser.ID, savedUser.Email)
	if err != nil {
		return "", nil, err
	}

	return token, savedUser, nil
}

// LoginWithRefreshToken аутентифицирует пользователя и возвращает access и refresh токены
func (s *AuthService) LoginWithRefreshToken(ctx context.Context, email, password, ip, userAgent string) (accessToken, refreshToken string, user *models.User, err error) {
	// Аутентификация пользователя
	user, err = s.storage.GetUserByEmail(ctx, email)
	if err != nil || user == nil {
		return "", "", nil, types.ErrInvalidCredentials
	}

	// Проверка пароля
	if !user.CheckPassword(password) {
		return "", "", nil, types.ErrInvalidCredentials
	}

	// Генерация access токена
	accessToken, err = s.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return "", "", nil, err
	}

	// Генерация refresh токена
	tokenID, err := jwt.GenerateSecureTokenID()
	if err != nil {
		return "", "", nil, err
	}

	refreshTokenValue, err := jwt.GenerateRefreshToken(user.ID, tokenID, s.jwtSecret)
	if err != nil {
		return "", "", nil, err
	}

	// Сохранение refresh токена в базе данных
	refreshTokenModel := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenValue,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour), // 30 дней
		CreatedAt: time.Now(),
		UserAgent: userAgent,
		IP:        ip,
	}

	if err = s.storage.CreateRefreshToken(ctx, refreshTokenModel); err != nil {
		return "", "", nil, err
	}

	// Обновляем время последнего входа
	_ = s.storage.UpdateLastSeen(ctx, user.ID)

	return accessToken, refreshTokenValue, user, nil
}

// RegisterWithRefreshToken регистрирует нового пользователя и возвращает access и refresh токены
func (s *AuthService) RegisterWithRefreshToken(ctx context.Context, name, email, password, ip, userAgent string) (accessToken, refreshToken string, user *models.User, err error) {
	// Проверка существования пользователя
	existingUser, _ := s.storage.GetUserByEmail(ctx, email)
	if existingUser != nil {
		return "", "", nil, types.ErrUserAlreadyExists
	}

	// Создание пользователя
	newUser := &models.User{
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
	}

	// Установка пароля (хеширование происходит внутри метода)
	if err := newUser.SetPassword(password); err != nil {
		return "", "", nil, err
	}

	savedUser, err := s.storage.CreateUser(ctx, newUser)
	if err != nil {
		return "", "", nil, err
	}

	// Генерация access токена
	accessToken, err = s.GenerateJWT(savedUser.ID, savedUser.Email)
	if err != nil {
		return "", "", nil, err
	}

	// Генерация refresh токена
	tokenID, err := jwt.GenerateSecureTokenID()
	if err != nil {
		return "", "", nil, err
	}

	refreshTokenValue, err := jwt.GenerateRefreshToken(savedUser.ID, tokenID, s.jwtSecret)
	if err != nil {
		return "", "", nil, err
	}

	// Сохранение refresh токена в базе данных
	refreshTokenModel := &models.RefreshToken{
		UserID:    savedUser.ID,
		Token:     refreshTokenValue,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour), // 30 дней
		CreatedAt: time.Now(),
		UserAgent: userAgent,
		IP:        ip,
	}

	if err = s.storage.CreateRefreshToken(ctx, refreshTokenModel); err != nil {
		return "", "", nil, err
	}

	return accessToken, refreshTokenValue, savedUser, nil
}

// GenerateTokensForOAuth генерирует пару токенов для OAuth авторизации
func (s *AuthService) GenerateTokensForOAuth(ctx context.Context, userID int, email, ip, userAgent string) (accessToken, refreshToken string, err error) {
	log.Printf("GenerateTokensForOAuth called for user %d (%s)", userID, email)

	// Генерация access токена
	accessToken, err = s.GenerateJWT(userID, email)
	if err != nil {
		return "", "", err
	}

	// Генерация уникального ID для токена
	tokenID, err := jwt.GenerateSecureTokenID()
	if err != nil {
		return "", "", err
	}

	// Генерация refresh токена
	refreshToken, err = jwt.GenerateRefreshToken(userID, tokenID, s.jwtSecret)
	if err != nil {
		return "", "", err
	}

	// Сохранение refresh токена в базе данных
	refreshTokenModel := &models.RefreshToken{
		UserID:    userID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour), // 30 дней
		CreatedAt: time.Now(),
		UserAgent: userAgent,
		IP:        ip,
	}

	if err = s.storage.CreateRefreshToken(ctx, refreshTokenModel); err != nil {
		return "", "", err
	}

	log.Printf("OAuth tokens generated successfully for user %d", userID)
	return accessToken, refreshToken, nil
}

// RefreshTokens обновляет access и refresh токены
func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken, ip, userAgent string) (newAccessToken, newRefreshToken string, err error) {
	// Валидация refresh токена
	claims, err := jwt.ValidateRefreshToken(refreshToken, s.jwtSecret)
	if err != nil {
		log.Printf("RefreshTokens: JWT validation failed: %v", err)
		return "", "", err
	}

	// Проверка токена в базе данных
	storedToken, err := s.storage.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		log.Printf("RefreshTokens: Failed to get token from DB: %v", err)
		return "", "", types.ErrInvalidToken
	}
	if storedToken == nil {
		log.Printf("RefreshTokens: Token not found in DB")
		return "", "", types.ErrInvalidToken
	}

	// Проверка валидности токена
	if !storedToken.IsValid() {
		log.Printf("RefreshTokens: Token is not valid (isRevoked=%v, expired=%v)",
			storedToken.IsRevoked, time.Now().After(storedToken.ExpiresAt))
		return "", "", types.ErrInvalidToken
	}

	// Получение пользователя
	user, err := s.storage.GetUserByID(ctx, claims.UserID)
	if err != nil || user == nil {
		return "", "", types.ErrUserNotFound
	}

	// Генерация нового access токена
	newAccessToken, err = s.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return "", "", err
	}

	// Ротация refresh токена - отзываем старый и создаем новый
	if err = s.storage.RevokeRefreshTokenByValue(ctx, refreshToken); err != nil {
		log.Printf("Failed to revoke old refresh token: %v", err)
	}

	// Генерация нового refresh токена
	newTokenID, err := jwt.GenerateSecureTokenID()
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err = jwt.GenerateRefreshToken(user.ID, newTokenID, s.jwtSecret)
	if err != nil {
		return "", "", err
	}

	// Сохранение нового refresh токена
	newRefreshTokenModel := &models.RefreshToken{
		UserID:    user.ID,
		Token:     newRefreshToken,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour), // 30 дней
		CreatedAt: time.Now(),
		UserAgent: userAgent,
		IP:        ip,
	}

	if err = s.storage.CreateRefreshToken(ctx, newRefreshTokenModel); err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

// RevokeRefreshToken отзывает refresh токен
func (s *AuthService) RevokeRefreshToken(ctx context.Context, refreshToken string) error {
	return s.storage.RevokeRefreshTokenByValue(ctx, refreshToken)
}

// RevokeUserRefreshTokens отзывает все refresh токены пользователя
func (s *AuthService) RevokeUserRefreshTokens(ctx context.Context, userID int) error {
	return s.storage.RevokeUserRefreshTokens(ctx, userID)
}
