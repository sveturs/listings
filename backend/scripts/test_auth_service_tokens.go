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
	"net/http"
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

type TestUser struct {
	UserID int
	Email  string
	Name   string
	Roles  []string
}

func main() {
	// Test users
	testUsers := []TestUser{
		{UserID: 5, Email: "voroshilovdo@gmail.com", Name: "Dmitry Voroshilov", Roles: []string{"admin"}},
		{UserID: 2, Email: "user@example.com", Name: "Test User", Roles: []string{"user"}},
		{UserID: 10, Email: "moderator@example.com", Name: "Test Moderator", Roles: []string{"moderator"}},
	}

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

	fmt.Println("Testing Auth Service JWT tokens with different users:")
	fmt.Println("============================================================")

	for _, user := range testUsers {
		// Create claims
		now := time.Now()
		claims := Claims{
			UserID: user.UserID,
			Email:  user.Email,
			Name:   user.Name,
			Roles:  user.Roles,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "https://auth.svetu.rs",
				Subject:   fmt.Sprintf("%d", user.UserID),
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
			log.Printf("Failed to sign token for %s: %v", user.Email, err)
			continue
		}

		fmt.Printf("\nUser: %s (ID: %d)\n", user.Email, user.UserID)
		fmt.Printf("Roles: %v\n", user.Roles)
		fmt.Printf("Testing token...\n")

		// Test the token
		req, err := http.NewRequest("GET", "http://localhost:3000/api/v1/user/profile", nil)
		if err != nil {
			log.Printf("Failed to create request: %v", err)
			continue
		}

		req.Header.Set("Authorization", "Bearer "+tokenString)

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Request failed: %v", err)
			continue
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)

		if resp.StatusCode == 200 {
			fmt.Printf("✅ Token validated successfully!\n")
			fmt.Printf("Response: %s\n", string(body))
		} else {
			fmt.Printf("❌ Token validation failed: %s\n", resp.Status)
			fmt.Printf("Response: %s\n", string(body))
		}

		fmt.Println("----------------------------------------")
	}

	fmt.Println("\n============================================================")
	fmt.Println("Test completed!")
}
