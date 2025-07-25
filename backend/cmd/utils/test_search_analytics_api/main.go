package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Get JWT token
	jwtToken := os.Getenv("TEST_JWT_TOKEN")
	if jwtToken == "" {
		// Generate a test token if not provided
		fmt.Println("No TEST_JWT_TOKEN found. Using hardcoded admin token.")
		// This is the token for admin user from previous sessions
		jwtToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo5LCJlbWFpbCI6ImFkbWluQHRlc3QuY29tIiwiaXNzIjoic3ZldHUtYmFja2VuZCIsInN1YiI6InVzZXI6OSIsImV4cCI6MTc1NDUxODY0OSwibmJmIjoxNzUxOTI2NjQ5LCJpYXQiOjE3NTE5MjY2NDl9." +
			"Dh13B3RTPIYGR6-cIVznnpOX_w3vnN4MOMmgBWuMYxQ" // #nosec G101 -- test token for development
	}

	// Test the analytics endpoint
	testEndpoint("http://localhost:3000/api/v1/admin/search/analytics?range=7d", jwtToken)

	// Also test the public statistics endpoint
	fmt.Println("\n=== Testing public statistics endpoint ===")
	testEndpoint("http://localhost:3000/api/v1/search/statistics", "")

	// Test popular searches
	fmt.Println("\n=== Testing popular searches endpoint ===")
	testEndpoint("http://localhost:3000/api/v1/search/statistics/popular", "")
}

func testEndpoint(url string, token string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Content-Type", "application/json")

	fmt.Printf("Testing endpoint: %s\n", url)
	if token != "" {
		fmt.Println("Using authentication: YES")
	} else {
		fmt.Println("Using authentication: NO")
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Failed to close response body: %v", err)
		}
	}()

	fmt.Printf("Status Code: %d\n", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response: %v", err)
		return
	}

	// Pretty print JSON response
	var jsonData interface{}
	if err := json.Unmarshal(body, &jsonData); err != nil {
		fmt.Printf("Raw Response: %s\n", string(body))
	} else {
		prettyJSON, _ := json.MarshalIndent(jsonData, "", "  ")
		fmt.Printf("Response:\n%s\n", string(prettyJSON))
	}
}
