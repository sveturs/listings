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

// Структуры из manifest.go - ПРАВИЛЬНЫЕ!
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

// Правильная структура для B2B - вложенная иерархия!
type WSPManifestRequest struct {
	ExtIdManifest  string          `json:"ExtIdManifest"`            // ОБЯЗАТЕЛЬНО: внешний ID манифеста
	IdTipPosiljke  int             `json:"IdTipPosiljke"`            // ВАЖНО: Тип посылки на ВЕРХНЕМ уровне! 1=обычная, 2=возврат
	Posiljalac     WSPPosiljalac   `json:"Posiljalac"`
	Porudzbine     []WSPPorudzbina `json:"Porudzbine"` // Заказы (верхний уровень)
	DatumPrijema   string          `json:"DatumPrijema"`
	VremePrijema   string          `json:"VremePrijema,omitempty"`
	IdPostePrijema int             `json:"IdPostePrijema,omitempty"`
	IdPartnera     int             `json:"IdPartnera,omitempty"`     // ВАЖНО: 10109 для svetu.rs
	NazivManifesta string          `json:"NazivManifesta,omitempty"` // ВАЖНО: название манифеста
}

// Поруджбина (Заказ) содержит массив посылок!
type WSPPorudzbina struct {
	ExtIdPorudzbina       string        `json:"ExtIdPorudzbina,omitempty"`       // External ID заказа
	ExtIdPorudzbinaKupca  string        `json:"ExtIdPorudzbinaKupca,omitempty"`  // External ID заказа клиента
	IndGrupnostUrucenja   *bool         `json:"IndGrupnostUrucenja,omitempty"`   // Групповая доставка
	Posiljke              []WSPPosiljka `json:"Posiljke"`                        // Посылки внутри заказа!
}

type WSPPosiljalac struct {
	Naziv         string     `json:"Naziv"`
	Adresa        *WSPAdresa `json:"Adresa"`        // ВАЖНО: Adresa - это ОБЪЕКТ!
	Mesto         string     `json:"Mesto"`
	PostanskiBroj string     `json:"PostanskiBroj"`
	Telefon       string     `json:"Telefon"`
	Email         string     `json:"Email,omitempty"`
	PIB           string     `json:"PIB,omitempty"`
	MaticniBroj   string     `json:"MaticniBroj,omitempty"`
	Kontakt       string     `json:"Kontakt,omitempty"`
	IdUgovor      int        `json:"IdUgovor,omitempty"`
	OznakaZemlje  string     `json:"OznakaZemlje,omitempty"`
}

type WSPPosiljka struct {
	// Обязательные для B2B API
	ExtBrend           string        `json:"ExtBrend"`           // ОБЯЗАТЕЛЬНО: бренд (например "SVETU")
	ExtMagacin         string        `json:"ExtMagacin"`         // ОБЯЗАТЕЛЬНО: склад (например "WAREHOUSE1")
	ExtReferenca       string        `json:"ExtReferenca"`       // ОБЯЗАТЕЛЬНО: референс (уникальный ID в нашей системе)
	NacinPrijema       string        `json:"NacinPrijema"`       // ОБЯЗАТЕЛЬНО: способ приема (K=курьер, O=отделение)
	ImaPrijemniBrojDN  *bool         `json:"ImaPrijemniBrojDN"`  // *bool чтобы передать false (не nil!)
	NacinPlacanja      string          `json:"NacinPlacanja"`      // ОБЯЗАТЕЛЬНО: способ оплаты (N=наличные, K=карта, POF=postanska uplatnica)
	Posiljalac         WSPPosiljalac   `json:"Posiljalac"`         // ОБЯЗАТЕЛЬНО: отправитель ВНУТРИ посылки!
	MestoPreuzimanja   *WSPPosiljalac  `json:"MestoPreuzimanja,omitempty"` // Место забора (объект Korisnik)
	PosebneUsluge      string          `json:"PosebneUsluge,omitempty"`    // Особые услуги через запятую: "PNA" или "PNA,SMS"

	// Основные поля
	BrojPosiljke       string        `json:"BrojPosiljke"`  // ВАЖНО: уникальный номер
	IdRukovanje        int           `json:"IdRukovanje"`
	// IdTipPosiljke УДАЛЕНО - оно теперь на верхнем уровне WSPManifestRequest!
	Primalac           WSPPrimalac   `json:"Primalac"`
	Masa               int     `json:"Masa"` // ВАЖНО: в ГРАММАХ, integer!
	Duzina             float64 `json:"Duzina,omitempty"`
	Sirina             float64 `json:"Sirina,omitempty"`
	Visina             float64 `json:"Visina,omitempty"`
	Otkupnina          int     `json:"Otkupnina,omitempty"`          // COD в ПАРАX (1 RSD = 100 para, 5000 RSD = 500000 para)
	Vrednost           int     `json:"Vrednost,omitempty"`           // ОБЯЗАТЕЛЬНО для COD: стоимость в ПАРАX!
	SMS                bool          `json:"SMS,omitempty"`
	Povratnica         bool          `json:"Povratnica,omitempty"`
	LicnoUrucenje      bool          `json:"LicnoUrucenje,omitempty"`
	PDK                bool          `json:"PDK,omitempty"`
	VD                 bool          `json:"VD,omitempty"`
	Sadrzaj            string        `json:"Sadrzaj,omitempty"`
	Napomena           string        `json:"Napomena,omitempty"`
	ReferencaBroj      string        `json:"ReferencaBroj,omitempty"`
}

type WSPPrimalac struct {
	TipAdrese string     `json:"TipAdrese"` // ВАЖНО: S=стандарт, F=Fah, P=Post restant
	Naziv     string     `json:"Naziv"`
	Telefon   string     `json:"Telefon"`
	Email     string     `json:"Email,omitempty"`
	Adresa    *WSPAdresa `json:"Adresa,omitempty"` // ВАЖНО: Adresa - это ОБЪЕКТ!
	Fah       string     `json:"Fah,omitempty"`
	BrojFaha  string     `json:"BrojFaha,omitempty"`
	IdPoste   int        `json:"IdPoste,omitempty"`
}

// Адрес - это сложный объект!
type WSPAdresa struct {
	Ulica         string `json:"Ulica,omitempty"`         // Название улицы
	Broj          string `json:"Broj,omitempty"`          // Номер дома
	Mesto         string `json:"Mesto,omitempty"`         // Город
	PostanskiBroj string `json:"PostanskiBroj,omitempty"` // Почтовый индекс
	PAK           string `json:"PAK,omitempty"`           // Почтовый адресный код
	OznakaZemlje  string `json:"OznakaZemlje,omitempty"`  // Код страны
}

// Откупнина передается просто как число в параx (1 RSD = 100 para)!
// Не нужна структура - это ОШИБКА в старом коде!

// REAL credentials
const (
	testEndpoint = "http://212.62.32.201/WspWebApi/transakcija"
	testUsername = "b2b@svetu.rs"
	testPassword = "Sv5et@U!"
)

func createClient() Klijent {
	return Klijent{
		Username:          testUsername,
		Password:          testPassword,
		Jezik:             "LAT",
		IdTipUredjaja:     2,
		NazivUredjaja:     "SVETU_PROD",
		ModelUredjaja:     "API Client",
		VerzijaOS:         "Linux",
		VerzijaAplikacije: "1.0.0",
		IPAdresa:          "127.0.0.1",
	}
}

func sendManifestRequest(manifest WSPManifestRequest, testName string) (*TransakcijaOut, error) {
	client := createClient()
	clientJSON, _ := json.Marshal(client)

	manifestJSON, _ := json.Marshal(manifest)

	request := TransakcijaIn{
		StrKlijent:         string(clientJSON),
		Servis:             101, // B2B
		IdVrstaTranskacije: 73,  // Manifest
		TipSerijalizacije:  2,   // JSON
		IdTransakcija:      uuid.New().String(),
		StrIn:              string(manifestJSON),
	}

	requestJSON, _ := json.Marshal(request)

	fmt.Printf("\n📤 REQUEST (%s):\n", testName)
	var prettyRequest map[string]interface{}
	json.Unmarshal(requestJSON, &prettyRequest)
	prettyJSON, _ := json.MarshalIndent(prettyRequest, "", "  ")
	fmt.Printf("%s\n", string(prettyJSON))

	httpReq, _ := http.NewRequest("POST", testEndpoint, bytes.NewBuffer(requestJSON))
	httpReq.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{Timeout: 30 * time.Second}
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	fmt.Printf("\n📥 RESPONSE (%s):\n", testName)
	fmt.Printf("HTTP Status: %s\n", resp.Status)

	var response TransakcijaOut
	json.Unmarshal(body, &response)
	prettyResponse, _ := json.MarshalIndent(response, "", "  ")
	fmt.Printf("%s\n", string(prettyResponse))

	return &response, nil
}

func testStandardShipment() {
	fmt.Println("\n" + string(bytes.Repeat([]byte("="), 80)))
	fmt.Println("🔧 TEST 1: STANDARD SHIPMENT (Правильная структура)")
	fmt.Println(string(bytes.Repeat([]byte("="), 80)))

	timestamp := time.Now().Unix()
	boolFalse := false // helper для *bool

	manifest := WSPManifestRequest{
		ExtIdManifest: fmt.Sprintf("MANIFEST-%d", timestamp), // ОБЯЗАТЕЛЬНО!
		IdTipPosiljke: 1, // ВАЖНО: Тип посылки на верхнем уровне! 1=обычная
		Posiljalac: WSPPosiljalac{
			Naziv: "SVETU Platforma d.o.o.",
			Adresa: &WSPAdresa{
				Ulica:         "Bulevar kralja Aleksandra",
				Broj:          "73",
				Mesto:         "Beograd",
				PostanskiBroj: "11000",
				OznakaZemlje:  "RS",
			},
			Mesto:         "Beograd",
			PostanskiBroj: "11000",
			Telefon:       "0641234567",
			Email:         "b2b@svetu.rs",
			OznakaZemlje:  "RS",
		},
		Porudzbine: []WSPPorudzbina{
			{
				ExtIdPorudzbina:      fmt.Sprintf("ORDER-%d", timestamp),
				ExtIdPorudzbinaKupca: fmt.Sprintf("CUSTOMER-ORDER-%d", timestamp),
				Posiljke: []WSPPosiljka{
					{
						// Обязательные B2B поля
						ExtBrend:          "SVETU",                             // бренд
						ExtMagacin:        "WAREHOUSE1",                        // склад
						ExtReferenca:      fmt.Sprintf("REF-%d", timestamp),    // референс
						NacinPrijema:      "K",                     // K=курьер, O=отделение
						ImaPrijemniBrojDN: &boolFalse,              // false как pointer
						NacinPlacanja:     "POF",                   // POF=postanska uplatnica (почтовая платежка)
						MestoPreuzimanja: &WSPPosiljalac{           // место забора для курьера
							Naziv: "SVETU Platforma d.o.o.",
							Adresa: &WSPAdresa{
								Ulica:         "Bulevar kralja Aleksandra",
								Broj:          "73",
								Mesto:         "Beograd",
								PostanskiBroj: "11000",
								OznakaZemlje:  "RS",
							},
							Mesto:         "Beograd",
							PostanskiBroj: "11000",
							Telefon:       "0641234567",
							Email:         "b2b@svetu.rs",
							OznakaZemlje:  "RS",
						},
						PosebneUsluge:     "PNA",                   // PNA=приём на адресе (для курьера обязательно!)
						Posiljalac: WSPPosiljalac{                              // отправитель внутри посылки!
							Naziv: "SVETU Platforma d.o.o.",
							Adresa: &WSPAdresa{
								Ulica:         "Bulevar kralja Aleksandra",
								Broj:          "73",
								Mesto:         "Beograd",
								PostanskiBroj: "11000",
								OznakaZemlje:  "RS",
							},
							Mesto:         "Beograd",
							PostanskiBroj: "11000",
							Telefon:       "0641234567",
							Email:         "b2b@svetu.rs",
							OznakaZemlje:  "RS",
						},

						// Основные поля
						BrojPosiljke: fmt.Sprintf("SVETU-TEST-%d", timestamp),
						IdRukovanje:  29,  // PE_Danas_za_sutra_12
						Masa:         500, // 500 грамм
						Duzina:       30,
						Sirina:       20,
						Visina:       10,
						Primalac: WSPPrimalac{
							TipAdrese: "S", // Стандартный адрес
							Naziv:     "Test Receiver 1",
							Telefon:   "0647654321",
							Email:     "test1@example.com",
							Adresa: &WSPAdresa{
								Ulica:         "Takovska",
								Broj:          "2",
								Mesto:         "Beograd",
								PostanskiBroj: "11000",
								OznakaZemlje:  "RS",
							},
						},
						Sadrzaj:       "Test package - Standard",
						ReferencaBroj: fmt.Sprintf("REF2-%d", timestamp),
					},
				},
			},
		},
		DatumPrijema:   time.Now().Format("2006-01-02"),
		VremePrijema:   time.Now().Format("15:04"),
		IdPartnera:     10109,
		NazivManifesta: fmt.Sprintf("SVETU-TEST-%s", time.Now().Format("20060102-150405")),
	}

	resp, err := sendManifestRequest(manifest, "Standard Shipment")
	if err != nil {
		fmt.Printf("❌ ERROR: %v\n", err)
		return
	}

	if resp.Rezultat == 0 {
		fmt.Println("✅ SUCCESS!")
	} else {
		fmt.Printf("⚠️  FAILED: Rezultat=%d\n", resp.Rezultat)
	}
}

func testCODShipment() {
	fmt.Println("\n" + string(bytes.Repeat([]byte("="), 80)))
	fmt.Println("💰 TEST 2: COD SHIPMENT")
	fmt.Println(string(bytes.Repeat([]byte("="), 80)))

	timestamp := time.Now().Unix()
	boolFalse := false // helper для *bool

	manifest := WSPManifestRequest{
		ExtIdManifest: fmt.Sprintf("MANIFEST-COD-%d", timestamp), // ОБЯЗАТЕЛЬНО!
		IdTipPosiljke: 1, // ВАЖНО: Тип посылки на верхнем уровне! 1=обычная
		Posiljalac: WSPPosiljalac{
			Naziv: "SVETU Platforma d.o.o.",
			Adresa: &WSPAdresa{
				Ulica:         "Bulevar kralja Aleksandra",
				Broj:          "73",
				Mesto:         "Beograd",
				PostanskiBroj: "11000",
				OznakaZemlje:  "RS",
			},
			Mesto:         "Beograd",
			PostanskiBroj: "11000",
			Telefon:       "0641234567",
			Email:         "b2b@svetu.rs",
			OznakaZemlje:  "RS",
		},
		Porudzbine: []WSPPorudzbina{
			{
				ExtIdPorudzbina:      fmt.Sprintf("ORDER-COD-%d", timestamp),
				ExtIdPorudzbinaKupca: fmt.Sprintf("CUSTOMER-COD-%d", timestamp),
				Posiljke: []WSPPosiljka{
					{
						// Обязательные B2B поля
						ExtBrend:          "SVETU",                             // бренд
						ExtMagacin:        "WAREHOUSE1",                        // склад
						ExtReferenca:      fmt.Sprintf("COD-REF-%d", timestamp),    // референс
						NacinPrijema:      "K",                     // K=курьер
						ImaPrijemniBrojDN: &boolFalse,              // false как pointer
						NacinPlacanja:     "POF",                   // POF=postanska uplatnica (почтовая платежка)
						MestoPreuzimanja: &WSPPosiljalac{           // место забора для курьера
							Naziv: "SVETU Platforma d.o.o.",
							Adresa: &WSPAdresa{
								Ulica:         "Bulevar kralja Aleksandra",
								Broj:          "73",
								Mesto:         "Beograd",
								PostanskiBroj: "11000",
								OznakaZemlje:  "RS",
							},
							Mesto:         "Beograd",
							PostanskiBroj: "11000",
							Telefon:       "0641234567",
							Email:         "b2b@svetu.rs",
							OznakaZemlje:  "RS",
						},
						PosebneUsluge:     "PNA,OTK,VD",            // PNA=приём, OTK=откупнина, VD=ценная посылка (для COD обязательно!)
						Posiljalac: WSPPosiljalac{                              // отправитель внутри посылки!
							Naziv: "SVETU Platforma d.o.o.",
							Adresa: &WSPAdresa{
								Ulica:         "Bulevar kralja Aleksandra",
								Broj:          "73",
								Mesto:         "Beograd",
								PostanskiBroj: "11000",
								OznakaZemlje:  "RS",
							},
							Mesto:         "Beograd",
							PostanskiBroj: "11000",
							Telefon:       "0641234567",
							Email:         "b2b@svetu.rs",
							OznakaZemlje:  "RS",
						},

						// Основные поля
						BrojPosiljke: fmt.Sprintf("SVETU-COD-%d", timestamp),
						IdRukovanje:  29,
						Masa:         750, // 750 грамм
						Duzina:       30,
						Sirina:       20,
						Visina:       10,
						Primalac: WSPPrimalac{
							TipAdrese: "S",
							Naziv:     "Test Receiver COD",
							Telefon:   "0649876543",
							Email:     "testcod@example.com",
							Adresa: &WSPAdresa{
								Ulica:         "Knez Mihailova",
								Broj:          "10",
								Mesto:         "Beograd",
								PostanskiBroj: "11000",
								OznakaZemlje:  "RS",
							},
						},
						Otkupnina:     500000, // COD в ПАРАX: 5000 RSD = 500000 para (1 RSD = 100 para)
						Vrednost:      500000, // Стоимость в ПАРАX: 5000 RSD = 500000 para (обязательно для COD!)
						Sadrzaj:       "Test COD package",
						ReferencaBroj: fmt.Sprintf("COD2-%d", timestamp),
					},
				},
			},
		},
		DatumPrijema:   time.Now().Format("2006-01-02"),
		VremePrijema:   time.Now().Format("15:04"),
		IdPartnera:     10109,
		NazivManifesta: fmt.Sprintf("SVETU-COD-%s", time.Now().Format("20060102-150405")),
	}

	resp, err := sendManifestRequest(manifest, "COD Shipment")
	if err != nil {
		fmt.Printf("❌ ERROR: %v\n", err)
		return
	}

	if resp.Rezultat == 0 {
		fmt.Println("✅ SUCCESS!")
	} else {
		fmt.Printf("⚠️  FAILED: Rezultat=%d\n", resp.Rezultat)
	}
}

func main() {
	fmt.Println("🚀 POST EXPRESS - ПРАВИЛЬНАЯ СТРУКТУРА MANIFEST")
	fmt.Println("=" + string(bytes.Repeat([]byte("="), 79)))
	fmt.Println("📡 Endpoint:", testEndpoint)
	fmt.Println("👤 Username:", testUsername)
	fmt.Println("🔑 Password: ********")
	fmt.Println("📅 Date:", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("=" + string(bytes.Repeat([]byte("="), 79)))

	testStandardShipment()
	time.Sleep(2 * time.Second)

	testCODShipment()

	fmt.Println("\n" + string(bytes.Repeat([]byte("="), 80)))
	fmt.Println("✨ TESTS COMPLETED!")
	fmt.Println(string(bytes.Repeat([]byte("="), 80)))
}
