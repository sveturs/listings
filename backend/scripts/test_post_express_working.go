//go:build ignore
// +build ignore

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

// TransakcijaIn - структура запроса согласно документации
type TransakcijaIn struct {
	StrKlijent         string `json:"StrKlijent"`
	Servis             int    `json:"Servis"`
	IdVrstaTransakcije int    `json:"IdVrstaTranskacije"` // Обратите внимание на правильное написание!
	TipSerijalizacije  int    `json:"TipSerijalizacije"`
	IdTransakcija      string `json:"IdTransakcija"`
	StrIn              string `json:"StrIn,omitempty"`
}

// Klijent - структура клиента согласно документации
type Klijent struct {
	Username          string  `json:"Username"`
	Password          string  `json:"Password"`
	Jezik             string  `json:"Jezik"`
	IdTipUredjaja     string  `json:"IdTipUredjaja"` // String, не int!
	VerzijaOS         string  `json:"VerzijaOS"`
	NazivUredjaja     string  `json:"NazivUredjaja"`
	ModelUredjaja     string  `json:"ModelUredjaja"`
	VerzijaAplikacije string  `json:"VerzijaAplikacije"`
	IPAdresa          string  `json:"IPAdresa"`
	Geolokacija       *string `json:"Geolokacija"`
}

// ManifestIn - структура для создания манифеста
type ManifestIn struct {
	ExtIdManifest string       `json:"ExtIdManifest"`
	IdTipPosiljke int          `json:"IdTipPosiljke"`
	Porudzbine    []Porudzbina `json:"Porudzbine"`
}

// Porudzbina - структура заказа
type Porudzbina struct {
	ExtIdPorudzbinaKupca string     `json:"ExtIdPorudzbinaKupca,omitempty"`
	ExtIdPorudzbina      string     `json:"ExtIdPorudzbina"`
	Posiljke             []Posiljka `json:"Posiljke"`
}

// Posiljka - структура посылки
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

// Korisnik - структура пользователя
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

// Adresa - структура адреса
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
	fmt.Println("🚀 Post Express WSP API Test (Working Version)")
	fmt.Println("==============================================\n")

	// Загружаем .env файл из родительской директории
	if err := godotenv.Load("../.env"); err != nil {
		fmt.Println("⚠️ Warning: Failed to load .env file, using environment variables")
	}

	endpoint := os.Getenv("POST_EXPRESS_WSP_ENDPOINT")
	username := os.Getenv("POST_EXPRESS_WSP_USERNAME")
	password := os.Getenv("POST_EXPRESS_WSP_PASSWORD")

	if endpoint == "" || username == "" || password == "" {
		fmt.Println("❌ Ошибка: не заданы переменные окружения POST_EXPRESS_WSP_*")
		os.Exit(1)
	}

	fmt.Printf("Endpoint: %s\n", endpoint)
	fmt.Printf("Username: %s\n", username)
	fmt.Printf("Password: %s...\n\n", password[:4])

	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Тест создания манифеста
	testManifestFixed(client, endpoint, username, password)
}

func testManifestFixed(client *http.Client, endpoint, username, password string) {
	fmt.Println("📦 Создание манифеста (ИСПРАВЛЕННАЯ версия)")
	fmt.Println("--------------------------------------------")

	// Подготовка клиентских данных
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

	// Место забора посылки (ОБЯЗАТЕЛЬНО для курьера!)
	// Используем Белград везде - точно есть в тестовой системе
	mestoPreuzimanja := &Korisnik{
		Vrsta:          "P",
		Naziv:          "Sve Tu d.o.o.",
		KontaktTelefon: "0641234567",
		KontaktOsoba:   "Test Manager",
		EMail:          "b2b@svetu.rs",
		Adresa: Adresa{
			OznakaZemlje:  "RS",
			Naselje:       "Beograd",
			Ulica:         "Bulevar kralja Aleksandra",
			Broj:          "73",
			PostanskiBroj: "11000",
		},
	}

	// Подготовка данных манифеста
	manifestData := ManifestIn{
		ExtIdManifest: "SVETU-" + fmt.Sprintf("%d", time.Now().Unix()),
		IdTipPosiljke: 1,
		Porudzbine: []Porudzbina{
			{
				ExtIdPorudzbina: "TEST-ORDER-001",
				Posiljke: []Posiljka{
					{
						Rbr:               1,
						ImaPrijemniBrojDN: "N",
						ExtBrend:          "SVETU",
						ExtMagacin:        "SVETU",
						ExtReferenca:      "TEST-REF-001",
						NacinPrijema:      "K", // Kurir
						IdRukovanje:       58,  // B2B handling
						NacinPlacanja:     "POF",
						Posiljalac: Korisnik{
							Vrsta:          "P", // Pravno lice
							Naziv:          "Sve Tu d.o.o.",
							KontaktTelefon: "0641234567",
							KontaktOsoba:   "Test Manager",
							EMail:          "b2b@svetu.rs",
							Adresa: Adresa{
								OznakaZemlje:  "RS", // Srbija
								Naselje:       "Beograd",
								Ulica:         "Bulevar kralja Aleksandra",
								Broj:          "73",
								PostanskiBroj: "11000",
							},
						},
						Primalac: Korisnik{
							Vrsta:          "F", // Fizičko lice
							Naziv:          "Petar Petrović", // ОБЯЗАТЕЛЬНО для API!
							Prezime:        "Petrović",
							Ime:            "Petar",
							KontaktTelefon: "0641234567",
							EMail:          "test@example.com",
							Adresa: Adresa{
								OznakaZemlje:  "RS", // Srbija
								Naselje:       "Beograd",
								Ulica:         "Takovska",
								Broj:          "2",
								PostanskiBroj: "11000",
							},
						},
						MestoPreuzimanja: mestoPreuzimanja,       // ОБЯЗАТЕЛЬНО для курьера!
						Masa:             500,                    // 500 грамм
						Vrednost:         0,                      // БЕЗ объявленной ценности (убираем VD)
						VrednostDTS:      0,                      // БЕЗ DTS
						Otkupnina:        0,                      // БЕЗ наложенного платежа
						Sadrzaj:          "Test paket za SVETU", // Содержимое
						PosebneUsluge:    "PNA",                  // Прием на адресе (ОБЯЗАТЕЛЬНО для курьера!)
					},
				},
			},
		},
	}

	manifestJSON, _ := json.Marshal(manifestData)

	fmt.Println("📤 Отправляемые данные манифеста:")
	var prettyManifest bytes.Buffer
	json.Indent(&prettyManifest, manifestJSON, "  ", "  ")
	fmt.Println(prettyManifest.String())
	fmt.Println()

	// Создаем запрос
	req := TransakcijaIn{
		StrKlijent:         string(clientJSON),
		Servis:             101,
		IdVrstaTransakcije: 73, // Создание манифеста
		TipSerijalizacije:  2,  // JSON
		IdTransakcija:      generateGUID(),
		StrIn:              string(manifestJSON),
	}

	// Отправляем запрос
	sendRequest(client, endpoint, req)
}

func sendRequest(client *http.Client, endpoint string, req TransakcijaIn) {
	reqJSON, _ := json.Marshal(req)

	fmt.Println("📤 HTTP POST запрос к API:")
	var prettyReq bytes.Buffer
	json.Indent(&prettyReq, reqJSON, "  ", "  ")
	fmt.Println(prettyReq.String())
	fmt.Println()

	resp, err := client.Post(endpoint, "application/json", bytes.NewBuffer(reqJSON))
	if err != nil {
		fmt.Printf("❌ HTTP ошибка: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("📥 Статус ответа: %s\n", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ Ошибка чтения ответа: %v\n", err)
		return
	}

	fmt.Println("📄 Тело ответа:")
	var prettyResp bytes.Buffer
	if err := json.Indent(&prettyResp, body, "", "  "); err != nil {
		fmt.Println(string(body))
	} else {
		fmt.Println(prettyResp.String())
	}
	fmt.Println()

	// Парсим ответ
	var response struct {
		Rezultat    int    `json:"Rezultat"`
		StrOut      string `json:"StrOut"`
		StrRezultat string `json:"StrRezultat"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("❌ Ошибка парсинга ответа: %v\n", err)
		return
	}

	// Анализируем результат
	switch response.Rezultat {
	case 0:
		fmt.Println("✅ Успех!")
		if response.StrOut != "" {
			fmt.Println("📦 Данные результата:")
			var prettyOut bytes.Buffer
			if err := json.Indent(&prettyOut, []byte(response.StrOut), "  ", "  "); err == nil {
				fmt.Println(prettyOut.String())
			} else {
				fmt.Println(response.StrOut)
			}
		}
	case 1:
		fmt.Println("⚠️ Предупреждение")
	case 3:
		fmt.Println("❌ Критическая ошибка")
	}

	if response.StrRezultat != "" {
		var rezultat map[string]interface{}
		if err := json.Unmarshal([]byte(response.StrRezultat), &rezultat); err == nil {
			if poruka, ok := rezultat["Poruka"].(string); ok {
				fmt.Printf("💬 Сообщение: %s\n", poruka)
			}
			if porukaKorisnik, ok := rezultat["PorukaKorisnik"].(string); ok {
				fmt.Printf("💬 Сообщение для пользователя: %s\n", porukaKorisnik)
			}
		}
	}

	// Если есть ошибки в StrOut, показываем их
	if response.StrOut != "" && response.Rezultat != 0 {
		var outData map[string]interface{}
		if err := json.Unmarshal([]byte(response.StrOut), &outData); err == nil {
			if greske, ok := outData["Greske"].([]interface{}); ok && len(greske) > 0 {
				fmt.Println("\n🔴 Ошибки валидации:")
				for i, g := range greske {
					if greska, ok := g.(map[string]interface{}); ok {
						if msg, ok := greska["PorukaGreske"].(string); ok {
							fmt.Printf("  %d. %s\n", i+1, msg)
						}
					}
				}
			}
		}
	}
}

func generateGUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
