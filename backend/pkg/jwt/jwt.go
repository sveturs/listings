package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// RefreshClaims для refresh токенов
type RefreshClaims struct {
	UserID  int    `json:"user_id"`
	TokenID string `json:"token_id"` // Уникальный ID для отзыва
	jwt.RegisteredClaims
}

func GenerateToken(userID int, email string, secret string) (string, error) {
	return GenerateTokenWithDuration(userID, email, secret, 1*time.Hour) // По умолчанию 1 час
}

func GenerateTokenWithDuration(userID int, email string, secret string, duration time.Duration) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now), // Токен не действителен до этого времени
			Issuer:    "svetu-backend",
			Subject:   fmt.Sprintf("user:%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenString string, secret string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		// jwt.ParseWithClaims уже проверяет ExpiresAt автоматически
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Дополнительная явная проверка истечения для безопасности
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, fmt.Errorf("token expired")
	}

	// Проверка на токены из будущего (защита от clock skew attacks)
	if claims.IssuedAt != nil && claims.IssuedAt.Time.After(time.Now().Add(5*time.Minute)) {
		return nil, fmt.Errorf("token issued in the future")
	}

	return claims, nil
}

// GenerateRefreshToken генерирует refresh токен для пользователя
func GenerateRefreshToken(userID int, tokenID string, secret string) (string, error) {
	return GenerateRefreshTokenWithDuration(userID, tokenID, secret, 30*24*time.Hour) // 30 дней по умолчанию
}

// GenerateRefreshTokenWithDuration генерирует refresh токен с заданной длительностью
func GenerateRefreshTokenWithDuration(userID int, tokenID string, secret string, duration time.Duration) (string, error) {
	now := time.Now()
	claims := &RefreshClaims{
		UserID:  userID,
		TokenID: tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "svetu-backend",
			Subject:   fmt.Sprintf("refresh:%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// GenerateSecureTokenID генерирует криптографически безопасный ID для токена
func GenerateSecureTokenID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate secure token ID: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// ValidateRefreshToken валидирует refresh токен
func ValidateRefreshToken(tokenString string, secret string) (*RefreshClaims, error) {
	claims := &RefreshClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("refresh token validation failed: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Проверка истечения
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, fmt.Errorf("refresh token expired")
	}

	// Проверка на токены из будущего
	if claims.IssuedAt != nil && claims.IssuedAt.Time.After(time.Now().Add(5*time.Minute)) {
		return nil, fmt.Errorf("refresh token issued in the future")
	}

	return claims, nil
}
