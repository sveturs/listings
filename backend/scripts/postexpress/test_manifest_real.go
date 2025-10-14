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

// Структуры для WSP API согласно документации
type Klijent struct {
	Username          string `json:"Username"`
	Password          string `json:"Password"`
	Jezik             string `json:"Jezik"`
	IdTipUredjaja     int    `json:"IdTipUredjaja"`
	NazivUredjaja     string `json:"NazivUredjaja,omitempty"`
	ModelUredjaja     string `json:"ModelUredjaja,omitempty"`
	VerzijaOS         string `json:"VerzijaOS,omitempty"`
	VerzijaAplikacije string `json:"VerzijaAplikacije,omitempty"`
	IPAdresa          string `json:"IPAdresa,omitempty"`
	Geolokacija       string `json:"Geolokacija,omitempty"`
}

type TransakcijaIn struct {
	StrKlijent         string `json:"StrKlijent"`
	Servis             int    `json:"Servis"`
	IdVrstaTranskacije int    `json:"IdVrstaTranskacije"`
	TipSerijalizacije  int    `json:"TipSerijalizacije"`
	IdTransakcija      string `json:"IdTransakcija"`
	StrIn              string `json:"StrIn"`
}

type TransakcijaOut struct {
	Rezultat    int             `json:"Rezultat"`
	StrOut      json.RawMessage `json:"StrOut"`
	StrRezultat json.RawMessage `json:"StrRezultat"`
}

// Структуры для манифеста (транзакция 73)
type ManifestRequest struct {
	Posiljalac   Posiljalac `json:"Posiljalac"`
	Posiljke     []Posiljka `json:"Posiljke"`
	DatumPrijema string     `json:"DatumPrijema"`
	VremePrijema string     `json:"VremePrijema"`
}

type Posiljalac struct {
	Ime      string `json:"Ime"`
	Adresa   string `json:"Adresa"`
	IdNaselje int    `json:"IdNaselje"`
	Telefon  string `json:"Telefon"`
	Email    string `json:"Email,omitempty"`
}

type Posiljka struct {
	Primalac         Primalac      `json:"Primalac"`
	TezinaPosiljke   int           `json:"TezinaPosiljke"`   // в граммах
	VrednostPosiljke int           `json:"VrednostPosiljke"` // в параx (1 RSD = 100 para)
	BrojOtkupnice    string        `json:"BrojOtkupnice,omitempty"`
	Sadrzaj          string        `json:"Sadrzaj"`
	Otkupnina        *Otkupnina    `json:"Otkupnina,omitempty"`      // COD
	UslugePPU        *string       `json:"UslugePPU,omitempty"`      // Дополнительные услуги (PNA, SMS, etc.)
	IdRukovanje      int           `json:"IdRukovanje"`              // 29, 30, 58, 71, 85 и т.д.
	NacinIsporuke    string        `json:"NacinIsporuke,omitempty"`  // K = Kurir, S = Šalter, PAK = Paketomat
	SifraOmegranice  *string       `json:"SifraOmegranice,omitempty"` // Код паккетомата для IdRukovanje=85
}

type Primalac struct {
	Ime       string `json:"Ime"`
	Adresa    string `json:"Adresa"`
	IdNaselje int    `json:"IdNaselje"`
	Telefon   string `json:"Telefon"`
	Email     string `json:"Email,omitempty"`
}

type Otkupnina struct {
	Iznos         int    `json:"Iznos"`                   // в параx (1 RSD = 100 para)
	NacinPlacanja string `json:"NacinPlacanja,omitempty"` // POF = готівка, K = картка
	BrojRacuna    string `json:"BrojRacuna,omitempty"`
}

// Test credentials from documentation
const (
	testEndpoint = "http://212.62.32.201/WspWebApi/transakcija"
	testUsername = "TEST"
	testPassword = "t3st"
)

func createClient() Klijent {
	return Klijent{
		Username:          testUsername,
		Password:          testPassword,
		Jezik:             "LAT",
		IdTipUredjaja:     2, // Server
		NazivUredjaja:     "SVETU_TEST",
		ModelUredjaja:     "Go Test Client",
		VerzijaOS:         "Linux",
		VerzijaAplikacije: "1.0.0",
		IPAdresa:          "127.0.0.1",
		Geolokacija:       "",
	}
}

func sendManifestRequest(manifest ManifestRequest, testName string) (*TransakcijaOut, error) {
	client := createClient()
	clientJSON, _ := json.Marshal(client)

	manifestJSON, _ := json.Marshal(manifest)

	request := TransakcijaIn{
		StrKlijent:         string(clientJSON),
		Servis:             3,  // B2B сервис
		IdVrstaTranskacije: 73, // Manifest (ВАЖНО!)
		TipSerijalizacije:  2,  // JSON
		IdTransakcija:      uuid.New().String(),
		StrIn:              string(manifestJSON),
	}

	requestJSON, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	fmt.Printf("\n📤 REQUEST (%s):\n", testName)
	var prettyRequest map[string]interface{}
	json.Unmarshal(requestJSON, &prettyRequest)
	prettyJSON, _ := json.MarshalIndent(prettyRequest, "", "  ")
	fmt.Printf("%s\n", string(prettyJSON))

	httpReq, _ := http.NewRequest("POST", testEndpoint, bytes.NewBuffer(requestJSON))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	httpClient := &http.Client{Timeout: 30 * time.Second}
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	fmt.Printf("\n📥 RESPONSE (%s):\n", testName)
	fmt.Printf("HTTP Status: %s\n", resp.Status)

	var response TransakcijaOut
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("Raw response: %s\n", string(body))
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	prettyResponse, _ := json.MarshalIndent(response, "", "  ")
	fmt.Printf("%s\n", string(prettyResponse))

	return &response, nil
}

func testStandardShipment() {
	fmt.Println("\n" + string(bytes.Repeat([]byte("="), 80)))
	fmt.Println("🔧 TEST 1: STANDARD SHIPMENT")
	fmt.Println(string(bytes.Repeat([]byte("="), 80)))

	services := "PNA" // Prijem na adresi
	manifest := ManifestRequest{
		Posiljalac: Posiljalac{
			Ime:       "SVETU Platforma d.o.o.",
			Adresa:    "Bulevar kralja Aleksandra 73",
			IdNaselje: 110000, // Beograd
			Telefon:   "0641234567",
			Email:     "b2b@svetu.rs",
		},
		Posiljke: []Posiljka{
			{
				Primalac: Primalac{
					Ime:       "Test Receiver 1",
					Adresa:    "Takovska 2",
					IdNaselje: 110000, // Beograd
					Telefon:   "0647654321",
					Email:     "test1@example.com",
				},
				TezinaPosiljke:   500,   // 500g
				VrednostPosiljke: 50000, // 500 RSD (в параx)
				Sadrzaj:          "Test package - Standard kurir",
				IdRukovanje:      29,          // PE_Danas_za_sutra_12
				NacinIsporuke:    "K",         // Kurir
				UslugePPU:        &services,
			},
		},
		DatumPrijema: time.Now().Format("2006-01-02"),
		VremePrijema: time.Now().Format("15:04"),
	}

	resp, err := sendManifestRequest(manifest, "Standard Shipment")
	if err != nil {
		fmt.Printf("❌ ERROR: %v\n", err)
		return
	}

	if resp.Rezultat == 0 {
		fmt.Println("✅ SUCCESS: Shipment created!")
	} else {
		fmt.Printf("⚠️  FAILED: Rezultat=%d\n", resp.Rezultat)
		fmt.Printf("Error details: %s\n", string(resp.StrRezultat))
	}
}

func testCODShipment() {
	fmt.Println("\n" + string(bytes.Repeat([]byte("="), 80)))
	fmt.Println("💰 TEST 2: COD SHIPMENT (Otkupna pošiljka)")
	fmt.Println(string(bytes.Repeat([]byte("="), 80)))

	services := "PNA,OTK" // Prijem na adresi + Otkupnina
	codAmount := 500000  // 5000 RSD в параx (1 RSD = 100 para)

	manifest := ManifestRequest{
		Posiljalac: Posiljalac{
			Ime:       "SVETU Platforma d.o.o.",
			Adresa:    "Bulevar kralja Aleksandra 73",
			IdNaselje: 110000,
			Telefon:   "0641234567",
			Email:     "b2b@svetu.rs",
		},
		Posiljke: []Posiljka{
			{
				Primalac: Primalac{
					Ime:       "Test Receiver COD",
					Adresa:    "Knez Mihailova 10",
					IdNaselje: 110000,
					Telefon:   "0649876543",
					Email:     "testcod@example.com",
				},
				TezinaPosiljke:   750,
				VrednostPosiljke: 120000, // 1200 RSD
				BrojOtkupnice:    fmt.Sprintf("OTK-%d", time.Now().Unix()),
				Sadrzaj:          "Test COD package - 5000 RSD",
				Otkupnina: &Otkupnina{
					Iznos:         codAmount,
					NacinPlacanja: "POF", // Gotovina (Cash on Delivery)
				},
				IdRukovanje:   29, // PE_Danas_za_sutra_12
				NacinIsporuke: "K",
				UslugePPU:     &services,
			},
		},
		DatumPrijema: time.Now().Format("2006-01-02"),
		VremePrijema: time.Now().Format("15:04"),
	}

	resp, err := sendManifestRequest(manifest, "COD Shipment")
	if err != nil {
		fmt.Printf("❌ ERROR: %v\n", err)
		return
	}

	if resp.Rezultat == 0 {
		fmt.Println("✅ SUCCESS: COD Shipment created!")
	} else {
		fmt.Printf("⚠️  FAILED: Rezultat=%d\n", resp.Rezultat)
		fmt.Printf("Error details: %s\n", string(resp.StrRezultat))
	}
}

func testParcelLocker() {
	fmt.Println("\n" + string(bytes.Repeat([]byte("="), 80)))
	fmt.Println("📦 TEST 3: PARCEL LOCKER (Paketomat)")
	fmt.Println(string(bytes.Repeat([]byte("="), 80)))

	services := "PNA"
	parcelLockerCode := "BEO-001-TEST"

	manifest := ManifestRequest{
		Posiljalac: Posiljalac{
			Ime:       "SVETU Platforma d.o.o.",
			Adresa:    "Bulevar kralja Aleksandra 73",
			IdNaselje: 110000,
			Telefon:   "0641234567",
			Email:     "b2b@svetu.rs",
		},
		Posiljke: []Posiljka{
			{
				Primalac: Primalac{
					Ime:       "Test Receiver Paketomat",
					Adresa:    "Paketomat BEO-001-TEST",
					IdNaselje: 110000,
					Telefon:   "0645555555",
					Email:     "testpak@example.com",
				},
				TezinaPosiljke:   600,
				VrednostPosiljke: 80000, // 800 RSD
				Sadrzaj:          "Test Parcel Locker package",
				IdRukovanje:      85,    // Isporuka_na_paketomatu
				NacinIsporuke:    "PAK", // Paketomat
				SifraOmegranice:  &parcelLockerCode,
				UslugePPU:        &services,
			},
		},
		DatumPrijema: time.Now().Format("2006-01-02"),
		VremePrijema: time.Now().Format("15:04"),
	}

	resp, err := sendManifestRequest(manifest, "Parcel Locker")
	if err != nil {
		fmt.Printf("❌ ERROR: %v\n", err)
		return
	}

	if resp.Rezultat == 0 {
		fmt.Println("✅ SUCCESS: Parcel Locker shipment created!")
	} else {
		fmt.Printf("⚠️  FAILED: Rezultat=%d\n", resp.Rezultat)
		fmt.Printf("Error details: %s\n", string(resp.StrRezultat))
	}
}

func testCODWithSMS() {
	fmt.Println("\n" + string(bytes.Repeat([]byte("="), 80)))
	fmt.Println("📱 TEST 4: COD + SMS NOTIFICATION")
	fmt.Println(string(bytes.Repeat([]byte("="), 80)))

	services := "PNA,OTK,SMS"
	codAmount := 1200000 // 12000 RSD

	manifest := ManifestRequest{
		Posiljalac: Posiljalac{
			Ime:       "SVETU Platforma d.o.o.",
			Adresa:    "Bulevar kralja Aleksandra 73",
			IdNaselje: 110000,
			Telefon:   "0641234567",
			Email:     "b2b@svetu.rs",
		},
		Posiljke: []Posiljka{
			{
				Primalac: Primalac{
					Ime:       "Test Receiver COD+SMS",
					Adresa:    "Takovska 5",
					IdNaselje: 214000, // Kragujevac
					Telefon:   "0643333333",
					Email:     "testcodsms@example.com",
				},
				TezinaPosiljke:   900,
				VrednostPosiljke: 150000, // 1500 RSD
				BrojOtkupnice:    fmt.Sprintf("OTKSMS-%d", time.Now().Unix()),
				Sadrzaj:          "Test COD+SMS package - 12000 RSD",
				Otkupnina: &Otkupnina{
					Iznos:         codAmount,
					NacinPlacanja: "POF",
				},
				IdRukovanje:   29,
				NacinIsporuke: "K",
				UslugePPU:     &services,
			},
		},
		DatumPrijema: time.Now().Format("2006-01-02"),
		VremePrijema: time.Now().Format("15:04"),
	}

	resp, err := sendManifestRequest(manifest, "COD+SMS")
	if err != nil {
		fmt.Printf("❌ ERROR: %v\n", err)
		return
	}

	if resp.Rezultat == 0 {
		fmt.Println("✅ SUCCESS: COD+SMS Shipment created!")
	} else {
		fmt.Printf("⚠️  FAILED: Rezultat=%d\n", resp.Rezultat)
		fmt.Printf("Error details: %s\n", string(resp.StrRezultat))
	}
}

func main() {
	fmt.Println("🚀 POST EXPRESS WSP API - REAL MANIFEST TEST (Transaction 73)")
	fmt.Println("=" + string(bytes.Repeat([]byte("="), 79)))
	fmt.Println("📡 Endpoint:", testEndpoint)
	fmt.Println("👤 Username:", testUsername)
	fmt.Println("🔑 Password: t3st")
	fmt.Println("📅 Date:", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("=" + string(bytes.Repeat([]byte("="), 79)))

	// Run all tests
	testStandardShipment()
	time.Sleep(1 * time.Second)

	testCODShipment()
	time.Sleep(1 * time.Second)

	testParcelLocker()
	time.Sleep(1 * time.Second)

	testCODWithSMS()

	fmt.Println("\n" + string(bytes.Repeat([]byte("="), 80)))
	fmt.Println("✨ ALL TESTS COMPLETED!")
	fmt.Println(string(bytes.Repeat([]byte("="), 80)))
	fmt.Println("\n📝 NOTES:")
	fmt.Println("- These tests use Transaction 73 (B2B Manifest) - the CORRECT way to create shipments")
	fmt.Println("- All tests send REAL requests to Post Express test API")
	fmt.Println("- Tracking numbers will be visible in Post Express system")
	fmt.Println("- Test credentials: TEST / t3st")
}
