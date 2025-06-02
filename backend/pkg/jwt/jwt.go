package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
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