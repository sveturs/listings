//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// Прямой IP адрес из .env.postexpress.test
	apiURL = "http://212.62.32.201/WspWebApi/transakcija"
)

func testSearchLocation() {
	fmt.Println("\n📍 Testing Location Search (GetNaselje)...")
	fmt.Println("==========================================")

	// Формируем запрос согласно документации
	// StrKlijent должен быть JSON строкой
	klijent := map[string]interface{}{
		"Username":          "TEST",
		"Password":          "t3st",
		"Jezik":             "LAT",
		"IdTipUredjaja":     11,
		"NazivUredjaja":     "TestAPI",
		"ModelUredjaja":     "GoClient",
		"VerzijaOS":         "Linux",
		"VerzijaAplikacije": "1.0.0",
		"IPAdresa":          "127.0.0.1",
		"Geolokacija":       nil,
	}

	klijentJSON, _ := json.Marshal(klijent)

	// StrIn для поиска населенного пункта
	naseljeIn := map[string]string{
		"Naziv": "Београд",
		"Ptt":   "",
	}
	naseljeJSON, _ := json.Marshal(naseljeIn)

	// Основной запрос
	request := map[string]interface{}{
		"StrKlijent":         string(klijentJSON),
		"Servis":             3, // Всегда 3
		"IdVrstaTranskacije": 3, // 3 = GetNaselje
		"TipSerijalizacije":  1, // 1 = JSON
		"IdTransakcija":      "test-" + fmt.Sprint(time.Now().Unix()),
		"StrIn":              string(naseljeJSON),
	}

	sendRequest(request, "Location Search")
}

func testTrackShipment() {
	fmt.Println("\n📦 Testing Shipment Tracking...")
	fmt.Println("==========================================")

	klijent := map[string]interface{}{
		"Username":          "TEST",
		"Password":          "t3st",
		"Jezik":             "LAT",
		"IdTipUredjaja":     11,
		"NazivUredjaja":     "BG01022W030",
		"ModelUredjaja":     "ASUS_M11",
		"VerzijaOS":         "Microsoft Windows NT 6.2.9200.0",
		"VerzijaAplikacije": "1.0.0.0",
		"IPAdresa":          "10.200.17.21",
		"Geolokacija":       nil,
		"Referenca":         "1",
	}

	klijentJSON, _ := json.Marshal(klijent)

	// Данные для отслеживания из примера
	kretanjeIn := map[string]string{
		"VrstaUsluge":  "1",
		"EksterniBroj": "",
		"PrijemniBroj": "PE746090324RS",
	}
	kretanjeJSON, _ := json.Marshal(kretanjeIn)

	// Запрос точно как в примере из документации
	request := map[string]interface{}{
		"StrKlijent":         string(klijentJSON),
		"Servis":             101, // 101 для отслеживания
		"IdVrstaTranskacije": 63,  // 63 = TTKretanje
		"TipSerijalizacije":  2,   // 2 = XML согласно примеру
		"IdTransakcija":      "e64b381e-7b32-4629-b227-bfaa88b8660e",
		"StrIn":              string(kretanjeJSON),
	}

	sendRequest(request, "Shipment Tracking")
}

func sendRequest(request map[string]interface{}, operation string) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		fmt.Printf("❌ Failed to marshal request: %v\n", err)
		return
	}

	fmt.Printf("\n📤 Sending %s request to:\n", operation)
	fmt.Printf("   URL: %s\n", apiURL)
	fmt.Printf("   Request body:\n")
	prettyJSON, _ := json.MarshalIndent(request, "   ", "  ")
	fmt.Printf("%s\n\n", string(prettyJSON))

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ Failed to create request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	fmt.Println("⏳ Sending request...")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Failed to send request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("✅ Response received! Status: %s\n", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ Failed to read response: %v\n", err)
		return
	}

	// Пробуем распарсить ответ
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("📥 Raw response (not JSON):\n%s\n", string(body))
		return
	}

	fmt.Println("📥 Response:")
	prettyResponse, _ := json.MarshalIndent(result, "", "  ")
	fmt.Printf("%s\n", string(prettyResponse))

	// Анализируем результат
	if rezultat, exists := result["Rezultat"]; exists {
		fmt.Printf("\n📊 Result code: %v\n", rezultat)
		switch fmt.Sprint(rezultat) {
		case "0":
			fmt.Println("✅ Success!")
		case "1":
			fmt.Println("⚠️ Partial success")
		case "2":
			fmt.Println("⚠️ Warning")
		case "3":
			fmt.Println("❌ Error")
		default:
			fmt.Printf("❓ Unknown result code: %v\n", rezultat)
		}
	}

	if strOut, exists := result["StrOut"]; exists && strOut != nil {
		fmt.Printf("\n📦 Output data:\n%v\n", strOut)
	}

	if strRezultat, exists := result["StrRezultat"]; exists && strRezultat != nil {
		fmt.Printf("\n📝 Result message:\n%v\n", strRezultat)
	}
}

func main() {
	fmt.Println("🚀 Testing WSP Post Express API (Direct IP)")
	fmt.Println("==========================================")
	fmt.Printf("Using endpoint: %s\n", apiURL)
	fmt.Println("Credentials: TEST / t3st")

	// Тест 1: Поиск населенного пункта
	testSearchLocation()

	time.Sleep(2 * time.Second)

	// Тест 2: Отслеживание посылки
	testTrackShipment()

	fmt.Println("\n==========================================")
	fmt.Println("✨ Testing completed!")
}
