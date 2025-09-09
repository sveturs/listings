package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Цвета для вывода в консоль
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

// Структуры для запросов и ответов
type CreateShipmentRequest struct {
	MarketplaceOrderID  int     `json:"marketplace_order_id"`
	SenderName          string  `json:"sender_name"`
	SenderAddress       string  `json:"sender_address"`
	SenderCity          string  `json:"sender_city"`
	SenderPostalCode    string  `json:"sender_postal_code"`
	SenderPhone         string  `json:"sender_phone"`
	RecipientName       string  `json:"recipient_name"`
	RecipientAddress    string  `json:"recipient_address"`
	RecipientCity       string  `json:"recipient_city"`
	RecipientPostalCode string  `json:"recipient_postal_code"`
	RecipientPhone      string  `json:"recipient_phone"`
	RecipientEmail      string  `json:"recipient_email"`
	WeightKg            float64 `json:"weight_kg"`
	LengthCm            float64 `json:"length_cm"`
	WidthCm             float64 `json:"width_cm"`
	HeightCm            float64 `json:"height_cm"`
	ServiceType         string  `json:"service_type"`
	CODAmount           float64 `json:"cod_amount"`
	InsuranceAmount     float64 `json:"insurance_amount"`
	DeliveryInstructions string `json:"delivery_instructions"`
	Notes               string  `json:"notes"`
}

type TrackingRequest struct {
	TrackingNumber string `json:"tracking_number"`
}

type CalculateRateRequest struct {
	SenderPostalCode    string  `json:"sender_postal_code"`
	RecipientPostalCode string  `json:"recipient_postal_code"`
	WeightKg            float64 `json:"weight_kg"`
	LengthCm            float64 `json:"length_cm"`
	WidthCm             float64 `json:"width_cm"`
	HeightCm            float64 `json:"height_cm"`
	ServiceType         string  `json:"service_type"`
	CODAmount           float64 `json:"cod_amount"`
	InsuranceAmount     float64 `json:"insurance_amount"`
}

func main() {
	baseURL := "http://localhost:3000/api/v1/postexpress"
	
	fmt.Printf("%s========================================%s\n", colorCyan, colorReset)
	fmt.Printf("%s  POST EXPRESS API ТЕСТИРОВАНИЕ%s\n", colorCyan, colorReset)
	fmt.Printf("%s========================================%s\n", colorCyan, colorReset)
	fmt.Printf("\nТестовый endpoint: %s%s%s\n\n", colorYellow, baseURL, colorReset)
	
	// 1. Получение JWT токена
	token := getJWTToken()
	if token == "" {
		fmt.Printf("%s❌ Не удалось получить JWT токен%s\n", colorRed, colorReset)
		return
	}
	fmt.Printf("%s✅ JWT токен получен%s\n", colorGreen, colorReset)
	
	// 2. Проверка статуса интеграции
	fmt.Printf("\n%s📋 1. ПРОВЕРКА СТАТУСА ИНТЕГРАЦИИ%s\n", colorBlue, colorReset)
	checkIntegrationStatus(baseURL, token)
	
	// 3. Расчет стоимости доставки
	fmt.Printf("\n%s💰 2. РАСЧЕТ СТОИМОСТИ ДОСТАВКИ%s\n", colorBlue, colorReset)
	testCalculateRate(baseURL, token)
	
	// 4. Получение списка локаций
	fmt.Printf("\n%s📍 3. ПОЛУЧЕНИЕ СПИСКА ЛОКАЦИЙ%s\n", colorBlue, colorReset)
	testGetLocations(baseURL, token)
	
	// 5. Создание тестовой посылки
	fmt.Printf("\n%s📦 4. СОЗДАНИЕ ТЕСТОВОЙ ПОСЫЛКИ (Транзакция 73)%s\n", colorBlue, colorReset)
	trackingNumber := testCreateShipment(baseURL, token)
	
	// 6. Отслеживание посылки
	if trackingNumber != "" {
		fmt.Printf("\n%s🔍 5. ОТСЛЕЖИВАНИЕ ПОСЫЛКИ%s\n", colorBlue, colorReset)
		testTrackShipment(baseURL, token, trackingNumber)
	}
	
	fmt.Printf("\n%s========================================%s\n", colorCyan, colorReset)
	fmt.Printf("%s  ТЕСТИРОВАНИЕ ЗАВЕРШЕНО%s\n", colorCyan, colorReset)
	fmt.Printf("%s========================================%s\n", colorCyan, colorReset)
}

func getJWTToken() string {
	// Используем скрипт create_test_jwt.go для получения токена
	token := os.Getenv("TEST_JWT_TOKEN")
	if token == "" {
		// Можно добавить логику генерации токена
		token = "test_token_placeholder"
	}
	return token
}

func checkIntegrationStatus(baseURL, token string) {
	resp, err := makeRequest("GET", baseURL+"/settings", nil, token)
	if err != nil {
		fmt.Printf("%s❌ Ошибка: %v%s\n", colorRed, err, colorReset)
		return
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	
	if resp.StatusCode == 200 {
		var result map[string]interface{}
		json.Unmarshal(body, &result)
		
		if data, ok := result["data"].(map[string]interface{}); ok {
			fmt.Printf("%s✅ Интеграция активна:%s\n", colorGreen, colorReset)
			fmt.Printf("   - Enabled: %v\n", data["enabled"])
			fmt.Printf("   - Test Mode: %v\n", data["test_mode"])
			fmt.Printf("   - WSP Endpoint: %v\n", data["wsp_endpoint"])
		}
	} else {
		fmt.Printf("%s❌ Ошибка получения статуса: %s%s\n", colorRed, string(body), colorReset)
	}
}

func testCalculateRate(baseURL, token string) {
	request := CalculateRateRequest{
		SenderPostalCode:    "21000",
		RecipientPostalCode: "11000",
		WeightKg:            2.5,
		LengthCm:            30,
		WidthCm:             20,
		HeightCm:            15,
		ServiceType:         "PE_Danas_za_sutra_12",
		CODAmount:           0,
		InsuranceAmount:     1000,
	}
	
	resp, err := makeRequest("POST", baseURL+"/rates/calculate", request, token)
	if err != nil {
		fmt.Printf("%s❌ Ошибка: %v%s\n", colorRed, err, colorReset)
		return
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	
	if resp.StatusCode == 200 {
		var result map[string]interface{}
		json.Unmarshal(body, &result)
		
		if data, ok := result["data"].(map[string]interface{}); ok {
			fmt.Printf("%s✅ Стоимость рассчитана:%s\n", colorGreen, colorReset)
			fmt.Printf("   - Базовая стоимость: %.2f RSD\n", data["base_price"])
			fmt.Printf("   - Страховка: %.2f RSD\n", data["insurance_fee"])
			fmt.Printf("   - Итого: %.2f RSD\n", data["total_price"])
		}
	} else {
		fmt.Printf("%s❌ Ошибка расчета: %s%s\n", colorRed, string(body), colorReset)
	}
}

func testGetLocations(baseURL, token string) {
	resp, err := makeRequest("GET", baseURL+"/locations?query=Novi", nil, token)
	if err != nil {
		fmt.Printf("%s❌ Ошибка: %v%s\n", colorRed, err, colorReset)
		return
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	
	if resp.StatusCode == 200 {
		var result map[string]interface{}
		json.Unmarshal(body, &result)
		
		if data, ok := result["data"].([]interface{}); ok {
			fmt.Printf("%s✅ Найдено локаций: %d%s\n", colorGreen, len(data), colorReset)
			
			// Показываем первые 3 локации
			for i, loc := range data {
				if i >= 3 {
					break
				}
				if location, ok := loc.(map[string]interface{}); ok {
					fmt.Printf("   - %s (%s)\n", location["name"], location["postal_code"])
				}
			}
		}
	} else {
		fmt.Printf("%s❌ Ошибка получения локаций: %s%s\n", colorRed, string(body), colorReset)
	}
}

func testCreateShipment(baseURL, token string) string {
	request := CreateShipmentRequest{
		MarketplaceOrderID:  12345,
		SenderName:          "Sve Tu Test Sender",
		SenderAddress:       "Микија Манојловића 53",
		SenderCity:          "Нови Сад",
		SenderPostalCode:    "21000",
		SenderPhone:         "+381621234567",
		RecipientName:       "Test Recipient",
		RecipientAddress:    "Кнез Михаилова 1",
		RecipientCity:       "Београд",
		RecipientPostalCode: "11000",
		RecipientPhone:      "+381611234567",
		RecipientEmail:      "test@example.com",
		WeightKg:            1.5,
		LengthCm:            30,
		WidthCm:             20,
		HeightCm:            10,
		ServiceType:         "PE_Danas_za_sutra_12",
		CODAmount:           0,
		InsuranceAmount:     500,
		DeliveryInstructions: "Позвонить перед доставкой",
		Notes:               "Тестовая посылка через транзакцию 73",
	}
	
	fmt.Printf("%s📝 Создаем посылку через Manifest (транзакция 73)...%s\n", colorYellow, colorReset)
	
	resp, err := makeRequest("POST", baseURL+"/shipments", request, token)
	if err != nil {
		fmt.Printf("%s❌ Ошибка: %v%s\n", colorRed, err, colorReset)
		return ""
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	
	if resp.StatusCode == 201 || resp.StatusCode == 200 {
		var result map[string]interface{}
		json.Unmarshal(body, &result)
		
		if data, ok := result["data"].(map[string]interface{}); ok {
			trackingNumber := ""
			if tn, ok := data["tracking_number"].(string); ok {
				trackingNumber = tn
			}
			
			fmt.Printf("%s✅ Посылка создана успешно!%s\n", colorGreen, colorReset)
			fmt.Printf("   - ID: %.0f\n", data["id"])
			fmt.Printf("   - Tracking Number: %s\n", trackingNumber)
			fmt.Printf("   - Barcode: %s\n", data["barcode"])
			fmt.Printf("   - Status: %s\n", data["status"])
			fmt.Printf("   - Total Price: %.2f RSD\n", data["total_price"])
			
			return trackingNumber
		}
	} else {
		fmt.Printf("%s❌ Ошибка создания посылки:%s\n", colorRed, colorReset)
		fmt.Printf("   Status: %d\n", resp.StatusCode)
		fmt.Printf("   Response: %s\n", string(body))
		
		// Парсим ошибку для диагностики
		var errorResp map[string]interface{}
		if err := json.Unmarshal(body, &errorResp); err == nil {
			if msg, ok := errorResp["message"].(string); ok {
				fmt.Printf("   %sОшибка: %s%s\n", colorRed, msg, colorReset)
			}
			if details, ok := errorResp["details"].(string); ok {
				fmt.Printf("   %sДетали: %s%s\n", colorYellow, details, colorReset)
			}
		}
	}
	
	return ""
}

func testTrackShipment(baseURL, token, trackingNumber string) {
	request := TrackingRequest{
		TrackingNumber: trackingNumber,
	}
	
	resp, err := makeRequest("POST", baseURL+"/tracking", request, token)
	if err != nil {
		fmt.Printf("%s❌ Ошибка: %v%s\n", colorRed, err, colorReset)
		return
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	
	if resp.StatusCode == 200 {
		var result map[string]interface{}
		json.Unmarshal(body, &result)
		
		if data, ok := result["data"].(map[string]interface{}); ok {
			fmt.Printf("%s✅ Информация о посылке:%s\n", colorGreen, colorReset)
			fmt.Printf("   - Tracking Number: %s\n", data["tracking_number"])
			fmt.Printf("   - Status: %s\n", data["status"])
			
			if events, ok := data["events"].([]interface{}); ok {
				fmt.Printf("   - События отслеживания: %d\n", len(events))
				for _, event := range events {
					if e, ok := event.(map[string]interface{}); ok {
						fmt.Printf("     • %s - %s\n", e["date"], e["description"])
					}
				}
			}
		}
	} else {
		fmt.Printf("%s❌ Ошибка отслеживания: %s%s\n", colorRed, string(body), colorReset)
	}
}

func makeRequest(method, url string, data interface{}, token string) (*http.Response, error) {
	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonData)
	}
	
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	return client.Do(req)
}