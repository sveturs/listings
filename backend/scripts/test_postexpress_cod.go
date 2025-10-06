//go:build ignore
// +build ignore

// Тест 1: Отправление с НАЛОЖЕННЫМ ПЛАТЕЖОМ (Cash on Delivery)

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
	fmt.Println("🚀 Post Express Test - COD (Наложенный платеж)")
	fmt.Println("===============================================\n")

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

	mestoPreuzimanja := &Korisnik{
		Vrsta:          "P",
		Naziv:          "Sve Tu d.o.o.",
		KontaktTelefon: "0641234567",
		KontaktOsoba:   "COD Manager",
		EMail:          "b2b@svetu.rs",
		Adresa: Adresa{
			OznakaZemlje:  "RS",
			Naselje:       "Beograd",
			Ulica:         "Bulevar kralja Aleksandra",
			Broj:          "73",
			PostanskiBroj: "11000",
		},
	}

	manifestData := ManifestIn{
		ExtIdManifest: "SVETU-COD-" + fmt.Sprintf("%d", time.Now().Unix()),
		IdTipPosiljke: 1,
		Porudzbine: []Porudzbina{
			{
				ExtIdPorudzbina: "TEST-COD-ORDER-001",
				Posiljke: []Posiljka{
					{
						Rbr:               1,
						ImaPrijemniBrojDN: "N",
						ExtBrend:          "SVETU",
						ExtMagacin:        "SVETU",
						ExtReferenca:      "TEST-COD-REF-001",
						NacinPrijema:      "K", // Kurir
						IdRukovanje:       58,
						NacinPlacanja:     "POF",
						Posiljalac: Korisnik{
							Vrsta:          "P",
							Naziv:          "Sve Tu d.o.o.",
							KontaktTelefon: "0641234567",
							KontaktOsoba:   "COD Manager",
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
							Naziv:          "Ana Anić",
							Prezime:        "Anić",
							Ime:            "Ana",
							KontaktTelefon: "0651234567",
							EMail:          "ana@example.com",
							Adresa: Adresa{
								OznakaZemlje:  "RS",
								Naselje:       "Beograd",
								Ulica:         "Makedonska",
								Broj:          "5",
								PostanskiBroj: "11000",
							},
						},
						MestoPreuzimanja: mestoPreuzimanja,
						Masa:             1000,  // 1 kg
						Vrednost:         5000,  // ОБЯЗАТЕЛЬНО для COD!
						VrednostDTS:      0,
						Otkupnina:        5000,  // 5000 RSD COD - НАЛОЖЕННЫЙ ПЛАТЕЖ!
						Sadrzaj:          "Testni paket sa otkupninom",
						PosebneUsluge:    "PNA;OTK;VD", // Prijem na adresi + Otkupnina + Vrednosna
					},
				},
			},
		},
	}

	manifestJSON, _ := json.Marshal(manifestData)

	fmt.Println("📦 ТЕСТ: Отправление с наложенным платежом 5000 RSD")
	fmt.Println("----------------------------------------------------")

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
							fmt.Printf("💰 COD сумма: 5000 RSD\n")
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
