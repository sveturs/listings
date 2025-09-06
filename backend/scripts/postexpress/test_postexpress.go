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

	"github.com/google/uuid"
)

// WSP API structures
type TransactionRequest struct {
	StrKlijent         string `json:"StrKlijent"`
	Servis             int    `json:"Servis"`
	IdVrstaTransakcije int    `json:"IdVrstaTranskacije"` // Note: typo in API
	TipSerijalizacije  int    `json:"TipSerijalizacije"`
	IdTransakcija      string `json:"IdTransakcija"`
	StrIn              string `json:"StrIn"`
}

type TransactionResponse struct {
	Rezultat    int    `json:"Rezultat"`
	StrOut      string `json:"StrOut"`
	StrRezultat string `json:"StrRezultat"`
}

type Klijent struct {
	Username          string `json:"Username"`
	Password          string `json:"Password"`
	Jezik             string `json:"Jezik"`
	IdTipUredjaja     int    `json:"IdTipUredjaja"`
	NazivUredjaja     string `json:"NazivUredjaja"`
	ModelUredjaja     string `json:"ModelUredjaja"`
	VerzijaOS         string `json:"VerzijaOS"`
	VerzijaAplikacije string `json:"VerzijaAplikacije"`
	IPAdresa          string `json:"IPAdresa"`
	Geolokacija       string `json:"Geolokacija"`
	Referenca         string `json:"Referenca"`
}

const (
	// Test credentials from documentation
	testEndpoint = "http://212.62.32.201/WspWebApi/transakcija"
	testUsername = "TEST"
	testPassword = "t3st"
)

func createClient() *Klijent {
	return &Klijent{
		Username:          testUsername,
		Password:          testPassword,
		Jezik:             "LAT",
		IdTipUredjaja:     2, // Server
		NazivUredjaja:     "TEST_SERVER",
		ModelUredjaja:     "Go Test Client",
		VerzijaOS:         "Linux",
		VerzijaAplikacije: "1.0.0",
		IPAdresa:          "127.0.0.1",
		Geolokacija:       "",
		Referenca:         "1",
	}
}

func callWSPAPI(transactionID int, strIn interface{}) (*TransactionResponse, error) {
	client := createClient()
	clientJSON, _ := json.Marshal(client)

	var strInJSON []byte
	if strIn != nil {
		strInJSON, _ = json.Marshal(strIn)
	} else {
		strInJSON = []byte("{}")
	}

	request := TransactionRequest{
		StrKlijent:         string(clientJSON),
		Servis:             101,
		IdVrstaTransakcije: transactionID,
		TipSerijalizacije:  2,
		IdTransakcija:      uuid.New().String(),
		StrIn:              string(strInJSON),
	}

	requestJSON, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	fmt.Printf("📤 Request to WSP API (Transaction ID=%d):\n", transactionID)

	resp, err := http.Post(testEndpoint, "application/json", bytes.NewBuffer(requestJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to call API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var response TransactionResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

func testGetNaselje() {
	fmt.Println("\n🏘️  TEST 1: GetNaselje (ID=3) - Search for 'Novi Sad'")
	fmt.Println("=" + string(bytes.Repeat([]byte("="), 60)))

	input := map[string]interface{}{
		"Naziv":           "Novi Sad",
		"BrojSlogova":     10,
		"NacinSortiranja": 0,
	}

	response, err := callWSPAPI(3, input)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	if response.Rezultat == 0 {
		fmt.Println("✅ Success!")

		var result map[string]interface{}
		json.Unmarshal([]byte(response.StrOut), &result)

		if naselja, ok := result["Naselja"].([]interface{}); ok {
			fmt.Printf("📍 Found %d locations:\n", len(naselja))
			for i, naselje := range naselja {
				if n, ok := naselje.(map[string]interface{}); ok {
					fmt.Printf("   %d. %s (ID: %.0f)\n", i+1, n["Naziv"], n["Id"])
				}
			}
		}
	} else {
		fmt.Printf("❌ API Error (Result=%d)\n", response.Rezultat)
		fmt.Printf("   Message: %s\n", response.StrRezultat)
	}
}

func testCalculateRate() {
	fmt.Println("\n💰 TEST 2: Poštarina pošiljke (ID=11) - Calculate shipping rate")
	fmt.Println("=" + string(bytes.Repeat([]byte("="), 60)))

	input := map[string]interface{}{
		"IdRukovanje":            58, // PE_Danas_za_sutra_19
		"IdZemlja":               "0",
		"Masa":                   1000,  // 1kg in grams
		"Vrednost":               10000, // 100 RSD in paras
		"Otkupnina":              0,
		"VrstaOtkupnogDokumenta": "",
		"PosebneUsluge":          "",
		"IdPeBoxTip":             0,
	}

	response, err := callWSPAPI(11, input)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	if response.Rezultat == 0 {
		fmt.Println("✅ Success!")

		var result map[string]interface{}
		json.Unmarshal([]byte(response.StrOut), &result)

		if iznos, ok := result["Iznos"].(float64); ok {
			fmt.Printf("💵 Total price: %.2f RSD\n", iznos/100)
		}

		if stavovi, ok := result["CenovniStavovi"].([]interface{}); ok {
			fmt.Println("📋 Price breakdown:")
			for _, stav := range stavovi {
				if s, ok := stav.(map[string]interface{}); ok {
					fmt.Printf("   - %s: %.2f RSD\n", s["Naziv"], s["Iznos"].(float64)/100)
				}
			}
		}

		if napomene, ok := result["Napomene"].([]interface{}); ok && len(napomene) > 0 {
			fmt.Println("📝 Notes:")
			for _, napomena := range napomene {
				fmt.Printf("   %s\n", napomena)
			}
		}
	} else {
		fmt.Printf("❌ API Error (Result=%d)\n", response.Rezultat)
		fmt.Printf("   Message: %s\n", response.StrRezultat)
	}
}

func testCheckAvailability() {
	fmt.Println("\n🚚 TEST 3: Provera dostupnosti usluge (ID=9) - Check service availability")
	fmt.Println("=" + string(bytes.Repeat([]byte("="), 60)))

	input := map[string]interface{}{
		"TipAdrese":   2,  // Address of recipient
		"IdRukovanje": 58, // PE_Danas_za_sutra_19
		"Adresa": map[string]interface{}{
			"IdNaselje": 0,
			"Naselje":   "Novi Beograd",
			"IdUlica":   0,
			"Ulica":     "Bulevar Mihajla Pupina",
			"Broj":      "10",
			"Posta":     "11070",
			"Pak":       "",
			"Vrsta":     "S", // Standard address
		},
		"Datum": "",
	}

	response, err := callWSPAPI(9, input)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	if response.Rezultat == 0 {
		fmt.Println("✅ Success! Service is available at this address")

		var result map[string]interface{}
		json.Unmarshal([]byte(response.StrOut), &result)

		if adresa, ok := result["Adresa"].(map[string]interface{}); ok {
			fmt.Println("📬 Validated address:")
			fmt.Printf("   Location: %s\n", adresa["Naselje"])
			fmt.Printf("   Street: %s %v\n", adresa["Ulica"], adresa["Broj"])
			fmt.Printf("   Postal: %s\n", adresa["Posta"])
			if pak, ok := adresa["Pak"].(string); ok && pak != "" {
				fmt.Printf("   PAK: %s\n", pak)
			}
		}

		if poruke, ok := result["Poruke"].(string); ok && poruke != "" {
			fmt.Printf("⚠️  Warning: %s\n", poruke)
		}
	} else {
		fmt.Printf("❌ API Error (Result=%d)\n", response.Rezultat)
		fmt.Printf("   Message: %s\n", response.StrRezultat)
	}
}

func testTracking() {
	fmt.Println("\n📦 TEST 4: Pojedinačno praćenje (ID=63) - Track shipment")
	fmt.Println("=" + string(bytes.Repeat([]byte("="), 60)))

	// Using example tracking number from documentation
	input := map[string]interface{}{
		"VrstaUsluge":  1,
		"EksterniBroj": "",
		"PrijemniBroj": "PE750151869RS", // Example from docs
	}

	response, err := callWSPAPI(63, input)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	if response.Rezultat == 0 {
		fmt.Println("✅ Success!")

		var result map[string]interface{}
		json.Unmarshal([]byte(response.StrOut), &result)

		if kretanja, ok := result["Kretanja"].([]interface{}); ok {
			fmt.Printf("📍 Found %d tracking events:\n", len(kretanja))
			for _, kretanje := range kretanja {
				if k, ok := kretanje.(map[string]interface{}); ok {
					fmt.Printf("   • %s - %s @ %s\n", k["Datum"], k["Status"], k["Mesto"])
					if potpisnik, ok := k["Potpisnik"].(string); ok && potpisnik != "" {
						fmt.Printf("     Signed by: %s\n", potpisnik)
					}
				}
			}
		} else {
			fmt.Println("ℹ️  No tracking events found (shipment might not exist in test system)")
		}
	} else {
		fmt.Printf("❌ API Error (Result=%d)\n", response.Rezultat)
		fmt.Printf("   Message: %s\n", response.StrRezultat)
	}
}

func main() {
	fmt.Println("🚀 Post Express WSP API Integration Test")
	fmt.Println("📡 Testing endpoint:", testEndpoint)
	fmt.Println("🔑 Using test credentials: Username=TEST, Password=t3st")

	// Run all tests
	testGetNaselje()
	time.Sleep(500 * time.Millisecond) // Small delay between requests

	testCalculateRate()
	time.Sleep(500 * time.Millisecond)

	testCheckAvailability()
	time.Sleep(500 * time.Millisecond)

	testTracking()

	fmt.Println("\n✨ All tests completed!")
	fmt.Println("📝 Note: These tests use the official Post Express test server")
	fmt.Println("   Some data (like tracking) might not return results in test mode")
}
