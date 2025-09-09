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

// WSP API структуры
type TransakcijaIn struct {
	StrKlijent         string `json:"StrKlijent"`         // С большой буквы
	Servis             int    `json:"Servis"`             // С большой буквы
	IdVrstaTransakcije int    `json:"IdVrstaTransakcije"` // С большой буквы
	TipSerijalizacije  int    `json:"TipSerijalizacije"`  // С большой буквы
	IdTransakcija      string `json:"IdTransakcija"`      // С большой буквы
	StrIn              string `json:"StrIn"`              // С большой буквы
}

type Klijent struct {
	Username      string `json:"Username"`      // С большой буквы
	Password      string `json:"Password"`      // С большой буквы
	Jezik         string `json:"Jezik"`         // С большой буквы
	IdTipUredjaja int    `json:"IdTipUredjaja"` // С большой буквы
}

type GetNaseljeRequest struct {
	Naziv string `json:"naziv"`
}

type ManifestRequest struct {
	Posiljalac   Posiljalac `json:"Posiljalac"`
	Posiljke     []Posiljka `json:"Posiljke"`
	DatumPrijema string     `json:"DatumPrijema"`
}

type Posiljalac struct {
	Naziv         string `json:"Naziv"`
	Adresa        string `json:"Adresa"`
	Mesto         string `json:"Mesto"`
	PostanskiBroj string `json:"PostanskiBroj"`
	Telefon       string `json:"Telefon"`
}

type Posiljka struct {
	BrojPosiljke  string   `json:"BrojPosiljke"`
	IdRukovanje   int      `json:"IdRukovanje"`
	IdTipPosiljke int      `json:"IdTipPosiljke"`
	Primalac      Primalac `json:"Primalac"`
	Masa          float64  `json:"Masa"`
	Sadrzaj       string   `json:"Sadrzaj"`
}

type Primalac struct {
	TipAdrese     string `json:"TipAdrese"`
	Naziv         string `json:"Naziv"`
	Telefon       string `json:"Telefon"`
	Adresa        string `json:"Adresa"`
	Mesto         string `json:"Mesto"`
	PostanskiBroj string `json:"PostanskiBroj"`
}

func main() {
	fmt.Println("========================================")
	fmt.Println("  DIRECT WSP API TEST")
	fmt.Println("========================================")

	// Тестовый endpoint
	endpoint := "http://212.62.32.201/WspWebApi/transakcija"

	// Тестовые credentials из документации Post Express
	klijent := Klijent{
		Username:      "TEST", // Правильный username в верхнем регистре
		Password:      "t3st", // Правильный password в нижнем регистре
		Jezik:         "LAT",
		IdTipUredjaja: 2, // Server
	}

	// Тест 1: GetNaselje (транзакция 3)
	fmt.Println("\n1. TEST GetNaselje (Transaction ID: 3)")
	testGetNaselje(endpoint, klijent)

	// Тест 2: Manifest (транзакция 73)
	fmt.Println("\n2. TEST Manifest (Transaction ID: 73)")
	testManifest(endpoint, klijent)

	// Тест 3: Tracking (транзакция 63)
	fmt.Println("\n3. TEST Tracking (Transaction ID: 63)")
	testTracking(endpoint, klijent)
}

func testGetNaselje(endpoint string, klijent Klijent) {
	// Подготовка данных
	naseljeReq := GetNaseljeRequest{
		Naziv: "Novi",
	}

	inputJSON, _ := json.Marshal(naseljeReq)
	klijentJSON, _ := json.Marshal(klijent)

	// Создание транзакции
	transaction := TransakcijaIn{
		StrKlijent:         string(klijentJSON),
		Servis:             3, // Правильный servis для WSP API
		IdVrstaTransakcije: 3, // GetNaselje
		TipSerijalizacije:  2, // JSON
		IdTransakcija:      generateGUID(),
		StrIn:              string(inputJSON),
	}

	// Отправка запроса
	response := sendRequest(endpoint, transaction)
	fmt.Printf("Response: %s\n", response)
}

func testManifest(endpoint string, klijent Klijent) {
	// Подготовка манифеста
	manifest := ManifestRequest{
		Posiljalac: Posiljalac{
			Naziv:         "Sve Tu Test",
			Adresa:        "Test Adresa 1",
			Mesto:         "Novi Sad",
			PostanskiBroj: "21000",
			Telefon:       "+381621234567",
		},
		Posiljke: []Posiljka{
			{
				BrojPosiljke:  fmt.Sprintf("TEST-%d", time.Now().Unix()),
				IdRukovanje:   29, // PE_Danas_za_sutra_12
				IdTipPosiljke: 1,  // Обычная
				Primalac: Primalac{
					TipAdrese:     "S",
					Naziv:         "Test Primalac",
					Telefon:       "+381611234567",
					Adresa:        "Test Adresa 2",
					Mesto:         "Beograd",
					PostanskiBroj: "11000",
				},
				Masa:    1.5,
				Sadrzaj: "Test paket",
			},
		},
		DatumPrijema: time.Now().Format("2006-01-02"),
	}

	inputJSON, _ := json.Marshal(manifest)
	klijentJSON, _ := json.Marshal(klijent)

	// Создание транзакции
	transaction := TransakcijaIn{
		StrKlijent:         string(klijentJSON),
		Servis:             3,  // Правильный servis для WSP API
		IdVrstaTransakcije: 73, // Manifest
		TipSerijalizacije:  2,  // JSON
		IdTransakcija:      generateGUID(),
		StrIn:              string(inputJSON),
	}

	// Отправка запроса
	response := sendRequest(endpoint, transaction)
	fmt.Printf("Response: %s\n", response)

	// Попытка парсинга ответа
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(response), &result); err == nil {
		prettyJSON, _ := json.MarshalIndent(result, "", "  ")
		fmt.Printf("\nParsed Response:\n%s\n", string(prettyJSON))
	}
}

func testTracking(endpoint string, klijent Klijent) {
	// Подготовка запроса отслеживания
	trackingReq := map[string]string{
		"BrojPosiljke": "TEST123456",
	}

	inputJSON, _ := json.Marshal(trackingReq)
	klijentJSON, _ := json.Marshal(klijent)

	// Создание транзакции
	transaction := TransakcijaIn{
		StrKlijent:         string(klijentJSON),
		Servis:             3,  // Правильный servis для WSP API
		IdVrstaTransakcije: 63, // Tracking
		TipSerijalizacije:  2,  // JSON
		IdTransakcija:      generateGUID(),
		StrIn:              string(inputJSON),
	}

	// Отправка запроса
	response := sendRequest(endpoint, transaction)
	fmt.Printf("Response: %s\n", response)
}

func sendRequest(endpoint string, transaction TransakcijaIn) string {
	jsonData, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Sprintf("Error marshaling: %v", err)
	}

	fmt.Printf("Sending to: %s\n", endpoint)
	fmt.Printf("Request body: %s\n", string(jsonData))

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Sprintf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("Error reading response: %v", err)
	}

	fmt.Printf("Status: %d\n", resp.StatusCode)

	return string(body)
}

func generateGUID() string {
	return fmt.Sprintf("%d-%d", time.Now().Unix(), time.Now().Nanosecond())
}
