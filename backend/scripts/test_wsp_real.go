package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	
	"github.com/google/uuid"
)

// Структуры для WSP API согласно документации
type Klijent struct {
	Username          string  `json:"Username"`
	Password          string  `json:"Password"`
	Jezik             string  `json:"Jezik"`
	IdTipUredjaja     int     `json:"IdTipUredjaja"`
	NazivUredjaja     string  `json:"NazivUredjaja"`
	ModelUredjaja     string  `json:"ModelUredjaja"`
	VerzijaOS         string  `json:"VerzijaOS"`
	VerzijaAplikacije string  `json:"VerzijaAplikacije"`
	IPAdresa          string  `json:"IPAdresa"`
	Geolokacija       *string `json:"Geolokacija"`
}

type TransakcijaIn struct {
	StrKlijent         string `json:"StrKlijent"`
	Servis             int    `json:"Servis"`
	IdVrstaTranskacije int    `json:"IdVrstaTranskacije"`
	TipSerijalizacije  int    `json:"TipSerijalizacije"`
	IdTransakcija      string `json:"IdTransakcija"`
	StrIn              string `json:"StrIn"`
}

// Структуры для поиска населенного пункта
type NaseljeIn struct {
	Naziv string `json:"Naziv"`
	Ptt   string `json:"Ptt"`
}

// Структура для отслеживания посылки
type KretanjeIn struct {
	VrstaUsluge   string `json:"VrstaUsluge"`
	EksterniBroj  string `json:"EksterniBroj"`
	PrijemniBroj  string `json:"PrijemniBroj"`
}

func testGetNaselje() {
	fmt.Println("\n📍 Testing GetNaselje (Search Location)...")
	fmt.Println("==========================================")

	// Подготовка клиентских данных
	klijent := Klijent{
		Username:          "TEST",
		Password:          "t3st",
		Jezik:             "LAT",
		IdTipUredjaja:     11, // 11 для desktop согласно примеру
		NazivUredjaja:     "TestDevice",
		ModelUredjaja:     "API_Test",
		VerzijaOS:         "Linux",
		VerzijaAplikacije: "1.0.0.0",
		IPAdresa:          "127.0.0.1",
		Geolokacija:       nil,
	}

	klijentJSON, _ := json.Marshal(klijent)

	// Данные для поиска населенного пункта
	naseljeIn := NaseljeIn{
		Naziv: "Београд",
		Ptt:   "",
	}
	naseljeJSON, _ := json.Marshal(naseljeIn)

	// Формирование запроса
	request := TransakcijaIn{
		StrKlijent:         string(klijentJSON),
		Servis:             3,  // Всегда 3 для нашего сервиса
		IdVrstaTranskacije: 3,  // 3 = GetNaselje
		TipSerijalizacije:  1,  // 1 = JSON
		IdTransakcija:      uuid.New().String(),
		StrIn:              string(naseljeJSON),
	}

	sendRequest(request, "GetNaselje")
}

func testTrackShipment() {
	fmt.Println("\n📦 Testing Shipment Tracking...")
	fmt.Println("==========================================")

	// Подготовка клиентских данных
	klijent := Klijent{
		Username:          "TEST",
		Password:          "t3st",
		Jezik:             "LAT",
		IdTipUredjaja:     11,
		NazivUredjaja:     "TestDevice",
		ModelUredjaja:     "API_Test",
		VerzijaOS:         "Linux",
		VerzijaAplikacije: "1.0.0.0",
		IPAdresa:          "127.0.0.1",
		Geolokacija:       nil,
	}

	klijentJSON, _ := json.Marshal(klijent)

	// Данные для отслеживания (пример из документации)
	kretanjeIn := KretanjeIn{
		VrstaUsluge:  "1",
		EksterniBroj: "",
		PrijemniBroj: "PE746090324RS",
	}
	kretanjeJSON, _ := json.Marshal(kretanjeIn)

	// Формирование запроса
	request := TransakcijaIn{
		StrKlijent:         string(klijentJSON),
		Servis:             101, // 101 для отслеживания согласно примеру
		IdVrstaTranskacije: 63,  // 63 = TTKretanje
		TipSerijalizacije:  2,   // 2 = XML в примере, но пробуем JSON
		IdTransakcija:      uuid.New().String(),
		StrIn:              string(kretanjeJSON),
	}

	sendRequest(request, "Tracking")
}

func sendRequest(request TransakcijaIn, operation string) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		fmt.Printf("❌ Failed to marshal request: %v\n", err)
		return
	}

	fmt.Printf("📤 Request for %s:\n", operation)
	
	// Красиво выводим структуру запроса
	var prettyRequest map[string]interface{}
	json.Unmarshal(jsonData, &prettyRequest)
	prettyJSON, _ := json.MarshalIndent(prettyRequest, "", "  ")
	fmt.Printf("%s\n\n", string(prettyJSON))

	// Отправка запроса
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Используем правильный URL из документации
	url := "https://wsp.postexpress.rs/api/Transakcija"
	fmt.Printf("🌐 Sending to: %s\n", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ Failed to create request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	fmt.Println("⏳ Waiting for response...")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Failed to send request: %v\n", err)
		fmt.Println("   Possible reasons:")
		fmt.Println("   - Network/firewall issues")
		fmt.Println("   - SSL certificate problems")
		fmt.Println("   - API endpoint not accessible")
		return
	}
	defer resp.Body.Close()

	fmt.Printf("✅ Response Status: %s\n", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ Failed to read response: %v\n", err)
		return
	}

	// Пробуем распарсить как JSON
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		// Если не JSON, выводим как есть
		fmt.Printf("📥 Raw Response:\n%s\n", string(body))
	} else {
		// Красиво выводим JSON
		prettyResponse, _ := json.MarshalIndent(result, "", "  ")
		fmt.Printf("📥 Response:\n%s\n", string(prettyResponse))
		
		// Анализируем результат
		if ok, exists := result["OK"]; exists {
			if okBool, isBool := ok.(bool); isBool && okBool {
				fmt.Println("✅ Request successful!")
			} else {
				fmt.Println("⚠️ Request returned OK=false")
				if msg, exists := result["Poruka"]; exists {
					fmt.Printf("Message: %v\n", msg)
				}
			}
		}
	}
}

func main() {
	fmt.Println("🚀 Testing Real WSP Post Express API")
	fmt.Println("==========================================")
	fmt.Println("Using TEST credentials from documentation")
	fmt.Println("URL: https://wsp.postexpress.rs/api/Transakcija")
	
	// Тест 1: Поиск населенного пункта
	testGetNaselje()
	
	time.Sleep(2 * time.Second)
	
	// Тест 2: Отслеживание посылки
	testTrackShipment()
	
	fmt.Println("\n==========================================")
	fmt.Println("✨ Testing completed!")
	fmt.Println("\nNotes:")
	fmt.Println("- If connection failed, the API might require VPN or whitelist")
	fmt.Println("- If authentication failed, TEST account might be disabled")
	fmt.Println("- Production credentials required for real integration")
}