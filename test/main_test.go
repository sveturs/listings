package test

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "testing"
    "time"
)

const baseURL = "http://localhost:3000"

func TestMainEndpoints(t *testing.T) {
    // 1. Проверка корневого маршрута
    resp, err := http.Get(baseURL)
    if err != nil {
        t.Fatalf("Root endpoint error: %v", err)
    }
    if resp.StatusCode != http.StatusOK {
        t.Errorf("Expected status OK, got %v", resp.Status)
    }

    // 2. Тестирование аутентификации
    resp, err = http.Get(baseURL + "/auth/google")
    if err != nil {
        t.Fatalf("Auth endpoint error: %v", err)
    }
    if resp.StatusCode != http.StatusOK {
        t.Errorf("Expected status OK, got %v", resp.Status)
    }

    // 3. Создание комнаты
    room := map[string]interface{}{
        "name":               "Test Room",
        "capacity":          2,
        "price_per_night":   100,
        "address_city":      "Test City",
        "accommodation_type": "room",
    }
    jsonRoom, _ := json.Marshal(room)

    resp, err = http.Post(baseURL+"/rooms", "application/json", bytes.NewBuffer(jsonRoom))
    if err != nil {
        t.Fatalf("Create room error: %v", err)
    }
    if resp.StatusCode != http.StatusOK {
        t.Errorf("Expected status OK, got %v", resp.Status)
    }

    // 4. Получение списка комнат
    resp, err = http.Get(baseURL + "/rooms")
    if err != nil {
        t.Fatalf("Get rooms error: %v", err)
    }
    if resp.StatusCode != http.StatusOK {
        t.Errorf("Expected status OK, got %v", resp.Status)
    }

    // И так далее для остальных эндпоинтов...
}