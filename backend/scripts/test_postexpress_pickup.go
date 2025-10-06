//go:build ignore
// +build ignore

// Тест 3: Отправление с ПОЛУЧЕНИЕМ НА ПОЧТЕ (не курьером, а самовывоз)

package main

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type TransakcijaIn struct {
	StrKlijent         string `json:"StrKlijent"`
	Servis             int    `json:"Servis"`
	IdVrstaTransakcije int    `json:"IdVrstaTranskacije"`
	TipSerijalizacije  int    `json:"TipSerijalizacije"`
	IdTransakcija      string `json:"IdTransakcija"`
	StrIn              string `json:"StrIn,omitempty"`
}

type Klijent struct {
	Username          string  `json:"Username"`
	Password          string  `json:"Password"`
	Jezik             string  `json:"Jezik"`
	IdTipUredjaja     string  `json:"IdTipUredjaja"`
	VerzijaOS         string  `json:"VerzijaOS"`
	NazivUredjaja     string  `json:"NazivUredjaja"`
	ModelUredjaja     string  `json:"ModelUredjaja"`
	VerzijaAplikacije string  `json:"VerzijaAplikacije"`
	IPAdresa          string  `json:"IPAdresa"`
	Geolokacija       *string `json:"Geolokacija"`
}

type ManifestIn struct {
	ExtIdManifest string       `json:"ExtIdManifest"`
	IdTipPosiljke int          `json:"IdTipPosiljke"`
	Porudzbine    []Porudzbina `json:"Porudzbine"`
}

type Porudzbina struct {
	ExtIdPorudzbinaKupca string     `json:"ExtIdPorudzbinaKupca,omitempty"`
	ExtIdPorudzbina      string     `json:"ExtIdPorudzbina"`
	Posiljke             []Posiljka `json:"Posiljke"`
}

type Posiljka struct {
	Rbr               int       `json:"Rbr"`
	PrijemniBroj      string    `json:"PrijemniBroj,omitempty"`
	ImaPrijemniBrojDN string    `json:"ImaPrijemniBrojDN"`
	ExtBrend          string    `json:"ExtBrend"`
	ExtMagacin        string    `json:"ExtMagacin"`
	ExtReferenca      string    `json:"ExtReferenca"`
	NacinPrijema      string    `json:"NacinPrijema"`
	IdRukovanje       int       `json:"IdRukovanje"`
	NacinPlacanja     string    `json:"NacinPlacanja"`
	Posiljalac        Korisnik  `json:"Posiljalac"`
	Primalac          Korisnik  `json:"Primalac"`
	MestoPreuzimanja  *Korisnik `json:"MestoPreuzimanja,omitempty"`
	Masa              int       `json:"Masa"`
	Vrednost          int64     `json:"Vrednost"`
	VrednostDTS       int64     `json:"VrednostDTS"`
	Otkupnina         int64     `json:"Otkupnina"`
	Sadrzaj           string    `json:"Sadrzaj"`
	PosebneUsluge     string    `json:"PosebneUsluge,omitempty"`
}

type Korisnik struct {
	Vrsta          string `json:"Vrsta"`
	Naziv          string `json:"Naziv,omitempty"`
	Prezime        string `json:"Prezime,omitempty"`
	Ime            string `json:"Ime,omitempty"`
	KontaktTelefon string `json:"KontaktTelefon"`
	KontaktOsoba   string `json:"KontaktOsoba,omitempty"`
	EMail          string `json:"EMail,omitempty"`
	Adresa         Adresa `json:"Adresa"`
}

type Adresa struct {
	OznakaZemlje  string `json:"OznakaZemlje,omitempty"`
	IdNaselje     *int   `json:"IdNaselje,omitempty"`
	Naselje       string `json:"Naselje"`
	IdUlica       *int   `json:"IdUlica,omitempty"`
	Ulica         string `json:"Ulica,omitempty"`
	BrojPodbroj   string `json:"BrojPodbroj,omitempty"`
	Broj          string `json:"Broj,omitempty"`
	Podbroj       string `json:"Podbroj,omitempty"`
	Sprat         string `json:"Sprat,omitempty"`
	Stan          string `json:"Stan,omitempty"`
	PostanskiBroj string `json:"PostanskiBroj"`
	Pak           string `json:"Pak,omitempty"`
	Reon          string `json:"Reon,omitempty"`
	NazivPoste    string `json:"NazivPoste,omitempty"`
}

func main() {
	fmt.Println("🚀 Post Express Test - Post Office Pickup (Самовывоз с почты)")
	fmt.Println("=============================================================\n")

	godotenv.Load("../.env")

	endpoint := os.Getenv("POST_EXPRESS_WSP_ENDPOINT")
	username := os.Getenv("POST_EXPRESS_WSP_USERNAME")
	password := os.Getenv("POST_EXPRESS_WSP_PASSWORD")

	if endpoint == "" || username == "" || password == "" {
		fmt.Println("❌ Ошибка: не заданы переменные окружения")
		os.Exit(1)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	clientData := Klijent{
		Username:          username,
		Password:          password,
		Jezik:             "LAT",
		IdTipUredjaja:     "2",
		VerzijaOS:         "Linux",
		NazivUredjaja:     "SVETU",
		ModelUredjaja:     "SERVER",
		VerzijaAplikacije: "1.0.0",
		IPAdresa:          "127.0.0.1",
		Geolokacija:       nil,
	}

	clientJSON, _ := json.Marshal(clientData)

	// Для самовывоза НЕ нужно MestoPreuzimanja
	manifestData := ManifestIn{
		ExtIdManifest: "SVETU-PICKUP-" + fmt.Sprintf("%d", time.Now().Unix()),
		IdTipPosiljke: 1,
		Porudzbine: []Porudzbina{
			{
				ExtIdPorudzbina: "TEST-PICKUP-ORDER-001",
				Posiljke: []Posiljka{
					{
						Rbr:               1,
						ImaPrijemniBrojDN: "N",
						ExtBrend:          "SVETU",
						ExtMagacin:        "SVETU",
						ExtReferenca:      "TEST-PICKUP-REF-001",
						NacinPrijema:      "S", // S = Šalter (получение на почте)
						IdRukovanje:       58,
						NacinPlacanja:     "POF",
						Posiljalac: Korisnik{
							Vrsta:          "P",
							Naziv:          "Sve Tu d.o.o.",
							KontaktTelefon: "0641234567",
							KontaktOsoba:   "Pickup Manager",
							EMail:          "b2b@svetu.rs",
							Adresa: Adresa{
								OznakaZemlje:  "RS",
								Naselje:       "Beograd",
								Ulica:         "Bulevar kralja Aleksandra",
								Broj:          "73",
								PostanskiBroj: "11000",
							},
						},
						Primalac: Korisnik{
							Vrsta:          "F",
							Naziv:          "Jovan Jovanović",
							Prezime:        "Jovanović",
							Ime:            "Jovan",
							KontaktTelefon: "0631234567",
							EMail:          "jovan@example.com",
							Adresa: Adresa{
								OznakaZemlje:  "RS",
								Naselje:       "Beograd",
								Ulica:         "Kneza Miloša",
								Broj:          "20",
								PostanskiBroj: "11000",
							},
						},
						MestoPreuzimanja: nil,   // НЕ нужно для самовывоза
						Masa:             2000,  // 2 kg - больший вес
						Vrednost:         0,
						VrednostDTS:      0,
						Otkupnina:        0,
						Sadrzaj:          "Testni paket za preuzimanje na šalteru",
						PosebneUsluge:    "",    // Без дополнительных услуг для самовывоза
					},
				},
			},
		},
	}

	manifestJSON, _ := json.Marshal(manifestData)

	fmt.Println("📦 ТЕСТ: Отправление с получением на почте (самовывоз)")
	fmt.Println("-------------------------------------------------------")
	fmt.Println("📍 Способ получения: Šalter (не курьер)")
	fmt.Println("⚖️  Вес: 2 kg")
	fmt.Println()

	req := TransakcijaIn{
		StrKlijent:         string(clientJSON),
		Servis:             101,
		IdVrstaTransakcije: 73,
		TipSerijalizacije:  2,
		IdTransakcija:      generateGUID(),
		StrIn:              string(manifestJSON),
	}

	sendRequest(client, endpoint, req)
}

func sendRequest(client *http.Client, endpoint string, req TransakcijaIn) {
	reqJSON, _ := json.Marshal(req)

	resp, err := client.Post(endpoint, "application/json", bytes.NewBuffer(reqJSON))
	if err != nil {
		fmt.Printf("❌ HTTP ошибка: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var response struct {
		Rezultat    int    `json:"Rezultat"`
		StrOut      string `json:"StrOut"`
		StrRezultat string `json:"StrRezultat"`
	}
	json.Unmarshal(body, &response)

	if response.Rezultat == 0 {
		fmt.Println("✅ УСПЕХ!")

		var outData map[string]interface{}
		if json.Unmarshal([]byte(response.StrOut), &outData) == nil {
			if porudzbine, ok := outData["Porudzbine"].([]interface{}); ok && len(porudzbine) > 0 {
				if porudzbina, ok := porudzbine[0].(map[string]interface{}); ok {
					if posiljke, ok := porudzbina["Posiljke"].([]interface{}); ok && len(posiljke) > 0 {
						if posiljka, ok := posiljke[0].(map[string]interface{}); ok {
							fmt.Printf("📦 Tracking Number: %v\n", posiljka["PrijemniBroj"])
							fmt.Printf("📍 Получение: На почтовом отделении\n")
							fmt.Printf("⚖️  Вес: 2 kg\n")
							fmt.Printf("📍 ID Посылки: %v\n", posiljka["IdPosiljka"])
						}
					}
				}
			}
			fmt.Printf("🆔 ID Манифеста: %v\n", outData["IdManifest"])
		}
	} else {
		fmt.Println("❌ ОШИБКА")
		var prettyResp bytes.Buffer
		json.Indent(&prettyResp, body, "", "  ")
		fmt.Println(prettyResp.String())
	}
}

func generateGUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
