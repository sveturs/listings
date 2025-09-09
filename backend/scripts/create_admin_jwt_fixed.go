//go:build ignore
// +build ignore

package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int      `json:"user_id"`
	Email  string   `json:"email"`
	Name   string   `json:"name"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

func main() {
	// Read private key
	privateKeyData, err := ioutil.ReadFile("/data/auth_svetu/keys/private.pem")
	if err != nil {
		log.Fatalf("Failed to read private key: %v", err)
	}

	block, _ := pem.Decode(privateKeyData)
	if block == nil {
		log.Fatal("Failed to decode PEM block")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// Try PKCS8
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			log.Fatalf("Failed to parse private key: %v", err)
		}
		var ok bool
		privateKey, ok = key.(*rsa.PrivateKey)
		if !ok {
			log.Fatal("Not an RSA private key")
		}
	}

	// Create claims for admin user with CORRECT user_id from main DB
	// voroshilovdo@gmail.com has id=2 in main DB
	now := time.Now()
	claims := Claims{
		UserID: 2, // Correct ID from main DB
		Email:  "voroshilovdo@gmail.com",
		Name:   "Dmitry Voroshilov",
		Roles:  []string{"admin"},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "https://auth.svetu.rs",
			Subject:   "2", // Also fix subject to match
			Audience:  []string{"https://svetu.rs"},
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Sign with private key
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		log.Fatalf("Failed to sign token: %v", err)
	}

	fmt.Println(tokenString)
}
