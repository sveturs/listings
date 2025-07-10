package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type BehaviorEvent struct {
	SessionID string         `json:"session_id"`
	UserID    *string        `json:"user_id,omitempty"`
	EventType string         `json:"event_type"`
	EventData map[string]any `json:"event_data"`
	UserAgent string         `json:"user_agent"`
	IPAddress string         `json:"ip_address"`
	PageURL   string         `json:"page_url"`
	Referrer  string         `json:"referrer"`
}

func main() {
	apiURL := "http://localhost:3000/api/v1/analytics/track"

	// 5 пар событий с разными session_id
	testCases := []struct {
		sessionID string
		userID    *string
		query     string
		hostelID  int
		position  int
	}{
		{
			sessionID: "test-session-101",
			userID:    stringPtr("user-101"),
			query:     "дешевые хостелы москва",
			hostelID:  201,
			position:  1,
		},
		{
			sessionID: "test-session-102",
			userID:    stringPtr("user-102"),
			query:     "hostel near red square",
			hostelID:  202,
			position:  2,
		},
		{
			sessionID: "test-session-103",
			userID:    nil, // анонимный пользователь
			query:     "хостел центр москвы",
			hostelID:  203,
			position:  3,
		},
		{
			sessionID: "test-session-104",
			userID:    stringPtr("user-104"),
			query:     "капсульный отель москва",
			hostelID:  204,
			position:  1,
		},
		{
			sessionID: "test-session-105",
			userID:    stringPtr("user-105"),
			query:     "hostel with kitchen",
			hostelID:  205,
			position:  4,
		},
	}

	for i, tc := range testCases {
		// 1. Отправляем search_performed
		searchEvent := BehaviorEvent{
			SessionID: tc.sessionID,
			UserID:    tc.userID,
			EventType: "search_performed",
			EventData: map[string]any{
				"query":         tc.query,
				"filters":       map[string]any{},
				"results_count": 10,
			},
			UserAgent: fmt.Sprintf("Mozilla/5.0 Test Browser %d", i+1),
			IPAddress: fmt.Sprintf("192.168.1.%d", 100+i),
			PageURL:   "http://localhost:3001/search",
			Referrer:  "http://localhost:3001/",
		}

		if err := sendEvent(apiURL, searchEvent); err != nil {
			log.Printf("Error sending search event %d: %v", i+1, err)
			continue
		}
		fmt.Printf("✓ Sent search_performed for session %s\n", tc.sessionID)

		// Небольшая задержка между событиями
		time.Sleep(500 * time.Millisecond)

		// 2. Отправляем result_clicked
		clickEvent := BehaviorEvent{
			SessionID: tc.sessionID,
			UserID:    tc.userID,
			EventType: "result_clicked",
			EventData: map[string]any{
				"query":       tc.query,
				"hostel_id":   tc.hostelID,
				"position":    tc.position,
				"search_type": "text",
			},
			UserAgent: fmt.Sprintf("Mozilla/5.0 Test Browser %d", i+1),
			IPAddress: fmt.Sprintf("192.168.1.%d", 100+i),
			PageURL:   fmt.Sprintf("http://localhost:3001/search?q=%s", tc.query),
			Referrer:  "http://localhost:3001/search",
		}

		if err := sendEvent(apiURL, clickEvent); err != nil {
			log.Printf("Error sending click event %d: %v", i+1, err)
			continue
		}
		fmt.Printf("✓ Sent result_clicked for session %s (hostel_id: %d)\n", tc.sessionID, tc.hostelID)

		// Задержка перед следующей парой
		time.Sleep(1 * time.Second)
	}

	fmt.Println("\n✅ All test events sent successfully!")
}

func sendEvent(apiURL string, event BehaviorEvent) error {
	jsonData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error marshaling event: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	// Добавляем токен авторизации (можно получить из окружения или использовать тестовый)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyLTEwMSIsImV4cCI6MTczNjQxNTgzOX0.FAKE_TOKEN")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func stringPtr(s string) *string {
	return &s
}
