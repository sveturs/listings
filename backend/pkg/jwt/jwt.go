package jwt

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AuthServiceClaims - структура claims от auth service
type AuthServiceClaims struct {
	UserID   int      `json:"user_id"`
	Email    string   `json:"email"`
	Name     string   `json:"name"`
	Roles    []string `json:"roles"`
	Provider string   `json:"provider"`
	jwt.RegisteredClaims
}

// ValidateAuthServiceToken валидирует токен от auth service (RS256)
func ValidateAuthServiceToken(tokenString string, publicKey *rsa.PublicKey) (*AuthServiceClaims, error) {
	claims := &AuthServiceClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Проверяем что используется RS256
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Проверка истечения
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("token expired")
	}

	return claims, nil
}