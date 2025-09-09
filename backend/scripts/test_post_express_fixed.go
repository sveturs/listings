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

// TTKretanjeIn - структура для отслеживания посылки
type TTKretanjeIn struct {
	VrstaUsluge  int    `json:"VrstaUsluge"`
	EksterniBroj string `json:"EksterniBroj"`
	PrijemniBroj string `json:"PrijemniBroj"`
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
	Rbr               int      `json:"Rbr"`
	PrijemniBroj      string   `json:"PrijemniBroj,omitempty"`
	ImaPrijemniBrojDN string   `json:"ImaPrijemniBrojDN"`
	ExtBrend          string   `json:"ExtBrend"`
	ExtMagacin        string   `json:"ExtMagacin"`
	ExtReferenca      string   `json:"ExtReferenca"`
	NacinPrijema      string   `json:"NacinPrijema"`
	IdRukovanje       int      `json:"IdRukovanje"`
	NacinPlacanja     string   `json:"NacinPlacanja"`
	Posiljalac        Korisnik `json:"Posiljalac"`
	Primalac          Korisnik `json:"Primalac"`
	Masa              int      `json:"Masa"`
	Vrednost          int64    `json:"Vrednost"`
	VrednostDTS       int64    `json:"VrednostDTS"`
	Otkupnina         int64    `json:"Otkupnina"`
	Sadrzaj           string   `json:"Sadrzaj"`
	PosebneUsluge     string   `json:"PosebneUsluge,omitempty"`
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
	Ulica         string `json:"Ulica"`
	Broj          string `json:"Broj"`
	Podbroj       string `json:"Podbroj,omitempty"`
	Sprat         string `json:"Sprat,omitempty"`
	Stan          string `json:"Stan,omitempty"`
	PostanskiBroj string `json:"PostanskiBroj"`
	Pak           string `json:"Pak,omitempty"`
}

func main() {
	// Читаем учетные данные из переменных окружения
	username := os.Getenv("POST_EXPRESS_WSP_USERNAME")
	password := os.Getenv("POST_EXPRESS_WSP_PASSWORD")
	endpoint := os.Getenv("POST_EXPRESS_WSP_ENDPOINT")

	if username == "" {
		username = "b2b@svetu.rs"
	}
	if password == "" {
		password = "Sv5et@U!"
	}
	if endpoint == "" {
		endpoint = "http://212.62.32.201/WspWebApi/transakcija"
	}

	fmt.Println("🚀 Post Express WSP API Test (Fixed)")
	fmt.Println("=====================================")
	fmt.Printf("Endpoint: %s\n", endpoint)
	fmt.Printf("Username: %s\n", username)
	fmt.Printf("Password: %s\n", password[:3]+"...")
	fmt.Println("")

	// Создаем HTTP клиент
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	// Тест 1: Отслеживание посылки (транзакция 63)
	fmt.Println("📋 Тест 1: Отслеживание посылки (транзакция 63)")
	fmt.Println("------------------------------------------------")
	testTracking(client, endpoint, username, password)

	// Тест 2: Создание манифеста (транзакция 73)
	fmt.Println("\n📦 Тест 2: Создание манифеста (транзакция 73)")
	fmt.Println("----------------------------------------------")
	testManifest(client, endpoint, username, password)

	// Тест 3: Альтернативные учетные данные TEST
	fmt.Println("\n🔑 Тест 3: Использование тестовых учетных данных")
	fmt.Println("-------------------------------------------------")
	testWithTestCredentials(client, endpoint)
}

func testTracking(client *http.Client, endpoint, username, password string) {
	// Подготовка клиентских данных согласно документации
	clientData := Klijent{
		Username:          username,
		Password:          password,
		Jezik:             "LAT",
		IdTipUredjaja:     "2", // String!
		VerzijaOS:         "Linux",
		NazivUredjaja:     "SVETU",
		ModelUredjaja:     "SERVER",
		VerzijaAplikacije: "1.0.0",
		IPAdresa:          "127.0.0.1",
		Geolokacija:       nil,
	}

	clientJSON, _ := json.Marshal(clientData)

	// Подготовка данных для отслеживания
	trackingData := TTKretanjeIn{
		VrstaUsluge:  1,
		EksterniBroj: "",
		PrijemniBroj: "PE123456785",
	}

	trackingJSON, _ := json.Marshal(trackingData)

	// Создаем запрос
	req := TransakcijaIn{
		StrKlijent:         string(clientJSON),
		Servis:             101, // Как в примере из документации
		IdVrstaTransakcije: 63,
		TipSerijalizacije:  2, // 1 = JSON, 2 = XML
		IdTransakcija:      generateGUID(),
		StrIn:              string(trackingJSON),
	}

	// Отправляем запрос
	sendRequest(client, endpoint, req)
}

func testManifest(client *http.Client, endpoint, username, password string) {
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
						NacinPrijema:      "K",
						IdRukovanje:       58, // B2B handling
						NacinPlacanja:     "POF",
						Posiljalac: Korisnik{
							Vrsta:          "P",
							Naziv:          "Sve Tu d.o.o.",
							KontaktTelefon: "0641234567",
							KontaktOsoba:   "Test",
							EMail:          "b2b@svetu.rs",
							Adresa: Adresa{
								Naselje:       "Novi Sad",
								Ulica:         "Mikija Manojlovića",
								Broj:          "53",
								PostanskiBroj: "21000",
							},
						},
						Primalac: Korisnik{
							Vrsta:          "F",
							Prezime:        "Petrović",
							Ime:            "Petar",
							KontaktTelefon: "0641234567",
							EMail:          "test@example.com",
							Adresa: Adresa{
								Naselje:       "Beograd",
								Ulica:         "Kneza Miloša",
								Broj:          "10",
								PostanskiBroj: "11000",
							},
						},
						Masa:          500,
						Vrednost:      1000,
						VrednostDTS:   0,
						Otkupnina:     0,
						Sadrzaj:       "Test paket",
						PosebneUsluge: "SMS",
					},
				},
			},
		},
	}

	manifestJSON, _ := json.Marshal(manifestData)

	// Создаем запрос
	req := TransakcijaIn{
		StrKlijent:         string(clientJSON),
		Servis:             101,
		IdVrstaTransakcije: 73,
		TipSerijalizacije:  2, // JSON
		IdTransakcija:      generateGUID(),
		StrIn:              string(manifestJSON),
	}

	// Отправляем запрос
	sendRequest(client, endpoint, req)
}

func testWithTestCredentials(client *http.Client, endpoint string) {
	// Используем учетные данные из примера в документации
	clientData := Klijent{
		Username:          "TEST",
		Password:          "t3st",
		Jezik:             "LAT",
		IdTipUredjaja:     "11",
		VerzijaOS:         "Microsoft Windows NT 6.2.9200.0",
		NazivUredjaja:     "BG01022W030",
		ModelUredjaja:     "ASUS_M11",
		VerzijaAplikacije: "1.0.0.0",
		IPAdresa:          "10.200.17.21",
		Geolokacija:       nil,
	}

	clientJSON, _ := json.Marshal(clientData)

	// Простой тест отслеживания
	trackingData := TTKretanjeIn{
		VrstaUsluge:  1,
		EksterniBroj: "",
		PrijemniBroj: "PE123456785",
	}

	trackingJSON, _ := json.Marshal(trackingData)

	req := TransakcijaIn{
		StrKlijent:         string(clientJSON),
		Servis:             101,
		IdVrstaTransakcije: 63,
		TipSerijalizacije:  2,
		IdTransakcija:      generateGUID(),
		StrIn:              string(trackingJSON),
	}

	sendRequest(client, endpoint, req)
}

func sendRequest(client *http.Client, endpoint string, req TransakcijaIn) {
	// Маршалим запрос
	reqJSON, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		fmt.Printf("❌ Ошибка маршалинга: %v\n", err)
		return
	}

	fmt.Println("📤 Отправляемый запрос:")
	// Показываем только структуру, без полного JSON для читаемости
	fmt.Printf("  Servis: %d\n", req.Servis)
	fmt.Printf("  IdVrstaTransakcije: %d\n", req.IdVrstaTransakcije)
	fmt.Printf("  TipSerijalizacije: %d\n", req.TipSerijalizacije)
	fmt.Printf("  IdTransakcija: %s\n", req.IdTransakcija)

	// Парсим клиента для показа
	var klijent Klijent
	json.Unmarshal([]byte(req.StrKlijent), &klijent)
	fmt.Printf("  Username: %s\n", klijent.Username)
	fmt.Printf("  IdTipUredjaja: %s\n", klijent.IdTipUredjaja)

	// Создаем HTTP запрос
	httpReq, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(reqJSON))
	if err != nil {
		fmt.Printf("❌ Ошибка создания запроса: %v\n", err)
		return
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	// Отправляем запрос
	resp, err := client.Do(httpReq)
	if err != nil {
		fmt.Printf("❌ Ошибка отправки запроса: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ Ошибка чтения ответа: %v\n", err)
		return
	}

	fmt.Printf("\n📥 Статус ответа: %s\n", resp.Status)

	// Пытаемся распарсить как JSON для красивого вывода
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err == nil {
		prettyJSON, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println("📄 Тело ответа:")
		fmt.Println(string(prettyJSON))

		// Анализ результата
		if rezultat, ok := result["Rezultat"].(float64); ok {
			switch int(rezultat) {
			case 0:
				fmt.Println("✅ Успешно!")
			case 1:
				fmt.Println("⚠️ Предупреждение")
			case 2:
				fmt.Println("❌ Ошибка")
			case 3:
				fmt.Println("❌ Критическая ошибка")
			default:
				fmt.Printf("❓ Неизвестный код результата: %d\n", int(rezultat))
			}
		}

		// Выводим сообщения об ошибках, если есть
		if strRezultat, ok := result["StrRezultat"].(string); ok && strRezultat != "" {
			var errMsg map[string]interface{}
			if err := json.Unmarshal([]byte(strRezultat), &errMsg); err == nil {
				if poruka, ok := errMsg["Poruka"].(string); ok && poruka != "" {
					fmt.Printf("💬 Сообщение: %s\n", poruka)
				}
				if porukaKorisnik, ok := errMsg["PorukaKorisnik"].(string); ok && porukaKorisnik != "" {
					fmt.Printf("💬 Сообщение для пользователя: %s\n", porukaKorisnik)
				}
			}
		}
	} else {
		fmt.Println("📄 Тело ответа (raw):")
		fmt.Println(string(body))
	}
}

func generateGUID() string {
	// Генерируем настоящий GUID
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
