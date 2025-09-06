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

// Структуры для WSP API
type Klijent struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

type TransakcijaIn struct {
	TransakcijaId      int        `json:"TransakcijaId"`
	DatumVremePosiljke string     `json:"DatumVremePosiljke"`
	Klijent            Klijent    `json:"Klijent"`
	NaseljeIn          *NaseljeIn `json:"NaseljeIn,omitempty"`
}

type NaseljeIn struct {
	Naziv string `json:"Naziv"`
	Ptt   string `json:"Ptt"`
}

type TransakcijaOut struct {
	OK                bool        `json:"OK"`
	Poruka            string      `json:"Poruka"`
	TransakcijaId     int         `json:"TransakcijaId"`
	DatumVremePrijema string      `json:"DatumVremePrijema"`
	NaseljeOut        *NaseljeOut `json:"NaseljeOut,omitempty"`
}

type NaseljeOut struct {
	OK      bool      `json:"OK"`
	Poruka  string    `json:"Poruka"`
	Naselja []Naselje `json:"Naselja"`
}

type Naselje struct {
	Sifra   int    `json:"Sifra"`
	Naziv   string `json:"Naziv"`
	Ptt     string `json:"Ptt"`
	Opstina string `json:"Opstina"`
}

func testWSPAPI(username, password string) {
	fmt.Println("🔍 Testing WSP API Connection...")
	fmt.Println("API Endpoint: https://onlinepostexpress.rs/WSPWebApi/api/app/transakcija")
	fmt.Printf("Username: %s\n", username)
	fmt.Println("=========================================\n")

	// Подготовка запроса для поиска населенного пункта
	request := TransakcijaIn{
		TransakcijaId:      3, // ID для GetNaselje
		DatumVremePosiljke: time.Now().Format("2006-01-02T15:04:05"),
		Klijent: Klijent{
			Username: username,
			Password: password,
		},
		NaseljeIn: &NaseljeIn{
			Naziv: "Београд",
			Ptt:   "",
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		fmt.Printf("❌ Failed to marshal request: %v\n", err)
		return
	}

	fmt.Printf("📤 Request:\n%s\n\n", string(jsonData))

	// Отправка запроса
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("POST", "https://onlinepostexpress.rs/WSPWebApi/api/app/transakcija", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ Failed to create request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	fmt.Println("⏳ Sending request to WSP API...")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("❌ Failed to send request: %v\n", err)
		fmt.Printf("   This could mean:\n")
		fmt.Printf("   - Network connection issues\n")
		fmt.Printf("   - API endpoint is not accessible\n")
		fmt.Printf("   - SSL certificate issues\n")
		return
	}
	defer resp.Body.Close()

	fmt.Printf("📥 Response Status: %s\n", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ Failed to read response: %v\n", err)
		return
	}

	fmt.Printf("📥 Raw Response:\n%s\n\n", string(body))

	// Парсинг ответа
	var response TransakcijaOut
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("⚠️ Failed to parse response as JSON: %v\n", err)
		fmt.Println("Response might be in different format (HTML error page, etc.)")
		return
	}

	// Анализ результата
	if response.OK {
		fmt.Println("✅ API Connection Successful!")
		fmt.Printf("Transaction ID: %d\n", response.TransakcijaId)
		fmt.Printf("Server Time: %s\n", response.DatumVremePrijema)

		if response.NaseljeOut != nil && response.NaseljeOut.OK {
			fmt.Printf("\n📍 Found %d locations for 'Београд':\n", len(response.NaseljeOut.Naselja))
			for i, naselje := range response.NaseljeOut.Naselja {
				if i < 5 { // Показываем первые 5
					fmt.Printf("   %d. %s (PTT: %s, Municipality: %s)\n",
						naselje.Sifra, naselje.Naziv, naselje.Ptt, naselje.Opstina)
				}
			}
		}
	} else {
		fmt.Println("❌ API Request Failed!")
		fmt.Printf("Error Message: %s\n", response.Poruka)
		if response.NaseljeOut != nil {
			fmt.Printf("Details: %s\n", response.NaseljeOut.Poruka)
		}
	}
}

func main() {
	fmt.Println("🚀 WSP API Test Tool")
	fmt.Println("=========================================")

	// Пробуем разные варианты учетных данных
	testCases := []struct {
		name     string
		username string
		password string
	}{
		{"Test Account", "test", "test"},
		{"Demo Account", "demo", "demo"},
		{"Sandbox Account", "sandbox", "sandbox"},
		{"Guest Account", "guest", "guest"},
	}

	for _, tc := range testCases {
		fmt.Printf("\n🔧 Testing with %s credentials:\n", tc.name)
		testWSPAPI(tc.username, tc.password)
		fmt.Println("\n" + "=========================================")
		time.Sleep(2 * time.Second) // Пауза между попытками
	}

	fmt.Println("\n📝 Summary:")
	fmt.Println("If all attempts failed with authentication errors,")
	fmt.Println("it means we need real production credentials from Post Express.")
	fmt.Println("\nTo get real credentials:")
	fmt.Println("1. Register at: https://onlinepostexpress.rs/registracija")
	fmt.Println("2. Wait for account approval")
	fmt.Println("3. Use provided username and password")
}
