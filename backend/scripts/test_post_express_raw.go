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

func main() {
	endpoint := "http://212.62.32.201/WspWebApi/transakcija"

	fmt.Println("=====================================")
	fmt.Println("Post Express RAW API Test")
	fmt.Println("=====================================")
	fmt.Printf("Endpoint: %s\n", endpoint)
	fmt.Println("=====================================\n")

	// Тест 1: Простой запрос с корректными данными
	fmt.Println("Test 1: Basic transaction with correct structure...")

	// Подготовка данных клиента
	clientData := map[string]interface{}{
		"Username":          "b2b@svetu.rs",
		"Password":          "Vu3@#45$%67&*89()",
		"Jezik":             "sr",
		"IdTipUredjaja":     "2", // Строка!
		"VerzijaOS":         "Linux",
		"NazivUredjaja":     "SveTu-Server",
		"ModelUredjaja":     "API",
		"VerzijaAplikacije": "1.0.0",
		"IPAdresa":          "127.0.0.1",
		"IdPartnera":        10109, // Partner ID
	}

	clientJSON, _ := json.Marshal(clientData)

	// Запрос для поиска населенных пунктов
	searchData := map[string]interface{}{
		"Naziv":           "Novi Sad",
		"BrojSlogova":     10,
		"NacinSortiranja": 0,
	}

	searchJSON, _ := json.Marshal(searchData)

	// Формирование транзакции
	transaction := map[string]interface{}{
		"StrKlijent":         string(clientJSON),
		"Servis":             101, // 101 для B2B!
		"IdVrstaTransakcije": 3,   // Поиск населенных пунктов
		"TipSerijalizacije":  2,   // JSON
		"IdTransakcija":      fmt.Sprintf("TEST-%d", time.Now().Unix()),
		"StrIn":              string(searchJSON),
	}

	// Отправка запроса
	response, err := sendRequest(endpoint, transaction)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Response received:\n")
		prettyPrint(response)
	}

	fmt.Println("\n=====================================")

	// Тест 2: Создание манифеста (транзакция 73)
	fmt.Println("Test 2: Create manifest (transaction 73)...")

	// Данные манифеста
	manifest := map[string]interface{}{
		"Posiljalac": map[string]interface{}{
			"Naziv":         "Sve Tu Platform",
			"Adresa":        "Улица Микија Манојловића 53",
			"Mesto":         "Нови Сад",
			"PostanskiBroj": "21000",
			"Telefon":       "+381 60 1234567",
			"Email":         "b2b@svetu.rs",
			"IdUgovor":      0, // Если есть договор
		},
		"Posiljke": []map[string]interface{}{
			{
				"BrojPosiljke":  fmt.Sprintf("SVT-%d", time.Now().Unix()),
				"IdRukovanje":   29, // PE_Danas_za_sutra_12
				"IdTipPosiljke": 1,  // Обычная посылка
				"Masa":          1.5,
				"Duzina":        30,
				"Sirina":        20,
				"Visina":        10,
				"Primalac": map[string]interface{}{
					"TipAdrese":     "S",
					"Naziv":         "Тест Получатель",
					"Telefon":       "+381 60 7654321",
					"Adresa":        "Булевар Ослобођења 100",
					"Mesto":         "Београд",
					"PostanskiBroj": "11000",
				},
				"Sadrzaj":       "Тестовая посылка",
				"Napomena":      "Осторожно!",
				"ReferencaBroj": fmt.Sprintf("REF-%d", time.Now().Unix()),
			},
		},
		"DatumPrijema":   time.Now().Format("2006-01-02"),
		"VremePrijema":   time.Now().Format("15:04"),
		"IdPartnera":     10109,
		"NazivManifesta": fmt.Sprintf("SVETU-%s", time.Now().Format("20060102-150405")),
	}

	manifestJSON, _ := json.Marshal(manifest)

	// Транзакция для манифеста
	manifestTx := map[string]interface{}{
		"StrKlijent":         string(clientJSON),
		"Servis":             101, // 101 для B2B!
		"IdVrstaTransakcije": 73,  // Создание манифеста
		"TipSerijalizacije":  2,   // JSON
		"IdTransakcija":      fmt.Sprintf("MAN-%d", time.Now().Unix()),
		"StrIn":              string(manifestJSON),
	}

	response, err = sendRequest(endpoint, manifestTx)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Response received:\n")
		prettyPrint(response)
	}

	fmt.Println("\n=====================================")

	// Тест 3: Отслеживание (транзакция 15)
	fmt.Println("Test 3: Track shipment (transaction 15)...")

	trackData := map[string]interface{}{
		"BrojPosiljke": "TEST123456",
	}

	trackJSON, _ := json.Marshal(trackData)

	trackTx := map[string]interface{}{
		"StrKlijent":         string(clientJSON),
		"Servis":             101, // 101 для B2B!
		"IdVrstaTransakcije": 15,  // Отслеживание
		"TipSerijalizacije":  2,   // JSON
		"IdTransakcija":      fmt.Sprintf("TRK-%d", time.Now().Unix()),
		"StrIn":              string(trackJSON),
	}

	response, err = sendRequest(endpoint, trackTx)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
	} else {
		fmt.Printf("✅ Response received:\n")
		prettyPrint(response)
	}

	fmt.Println("\n=====================================")
	fmt.Println("Raw API test completed!")
	fmt.Println("=====================================")
}

func sendRequest(endpoint string, data interface{}) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	fmt.Printf("Request:\n")
	prettyPrint(data)
	fmt.Println()

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	fmt.Printf("Raw response: %s\n", string(body))

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		// Если не удается распарсить как JSON, возвращаем как строку
		return map[string]interface{}{
			"raw": string(body),
		}, nil
	}

	return result, nil
}

func prettyPrint(data interface{}) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("%v\n", data)
		return
	}
	fmt.Println(string(b))
}
