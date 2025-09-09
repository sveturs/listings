package main

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Структуры согласно документации WSP API

// TransakcijaIn - основная структура запроса к WSP API
type TransakcijaIn struct {
	StrKlijent         string `json:"StrKlijent"`
	Servis             int    `json:"Servis"`
	IdVrstaTransakcije int    `json:"IdVrstaTranskacije"` // ВАЖНО: в API опечатка - используется Transkacije вместо Transakcije!
	TipSerijalizacije  int    `json:"TipSerijalizacije"`
	IdTransakcija      string `json:"IdTransakcija"`
	StrIn              string `json:"StrIn,omitempty"`
}

// Klijent - данные клиента для аутентификации
type Klijent struct {
	Username          string  `json:"Username"`
	Password          string  `json:"Password"`
	Jezik             string  `json:"Jezik"`
	IdTipUredjaja     string  `json:"IdTipUredjaja"` // String "2" для веб-приложения
	VerzijaOS         string  `json:"VerzijaOS"`
	NazivUredjaja     string  `json:"NazivUredjaja"`
	ModelUredjaja     string  `json:"ModelUredjaja"`
	VerzijaAplikacije string  `json:"VerzijaAplikacije"`
	IPAdresa          string  `json:"IPAdresa"`
	Geolokacija       *string `json:"Geolokacija"`
	IdPartnera        int     `json:"IdPartnera,omitempty"` // 10109 для SVETU
}

// TTKretanjeIn - структура для отслеживания посылки (транзакция 63)
type TTKretanjeIn struct {
	VrstaUsluge  int    `json:"VrstaUsluge"`
	EksterniBroj string `json:"EksterniBroj"`
	PrijemniBroj string `json:"PrijemniBroj"`
}

// ManifestIn - структура для создания манифеста (транзакция 73)
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
	Rbr                int       `json:"Rbr"`
	PrijemniBroj       string    `json:"PrijemniBroj,omitempty"`
	ImaPrijemniBrojDN  string    `json:"ImaPrijemniBrojDN"`
	ExtBrend           string    `json:"ExtBrend"`
	ExtMagacin         string    `json:"ExtMagacin"`
	ExtReferenca       string    `json:"ExtReferenca"`
	NacinPrijema       string    `json:"NacinPrijema"`
	MestoPreuzimanja   *Korisnik `json:"MestoPreuzimanja,omitempty"` // Место преузимanja как объект Korisnik
	IdRukovanje        int       `json:"IdRukovanje"`
	NacinPlacanja      string    `json:"NacinPlacanja"`
	Posiljalac         Korisnik  `json:"Posiljalac"`
	Primalac           Korisnik  `json:"Primalac"`
	Masa               int       `json:"Masa"`
	Vrednost           int64     `json:"Vrednost"`
	VrednostDTS        int64     `json:"VrednostDTS"`
	Otkupnina          int64     `json:"Otkupnina"`
	Sadrzaj            string    `json:"Sadrzaj"`
	PosebneUsluge      string    `json:"PosebneUsluge,omitempty"`
}

// Korisnik - структура пользователя
type Korisnik struct {
	Vrsta          string `json:"Vrsta"`
	Naziv          string `json:"Naziv,omitempty"`     // Обязательно для физлиц
	Prezime        string `json:"Prezime,omitempty"`
	Ime            string `json:"Ime,omitempty"`
	KontaktTelefon string `json:"KontaktTelefon"`
	KontaktOsoba   string `json:"KontaktOsoba,omitempty"`
	EMail          string `json:"EMail,omitempty"`
	Adresa         Adresa `json:"Adresa"`
}

// Adresa - структура адреса
type Adresa struct {
	OznakaZemlje string `json:"OznakaZemlje,omitempty"`
	IdNaselje    *int   `json:"IdNaselje,omitempty"`
	Naselje      string `json:"Naselje"`
	IdUlica      *int   `json:"IdUlica,omitempty"`
	Ulica        string `json:"Ulica"`
	Broj         string `json:"Broj"`
	PostBroj     string `json:"PostBroj"`
	Sprat        string `json:"Sprat,omitempty"`
	Stan         string `json:"Stan,omitempty"`
}

// generateGUID генерирует уникальный идентификатор
func generateGUID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// callWSPAPI выполняет вызов WSP API
func callWSPAPI(transakcija *TransakcijaIn) (map[string]interface{}, error) {
	// Тестовый эндпоинт (подтвержденный рабочий URL)
	url := "http://212.62.32.201/WspWebApi/transakcija"

	// Сериализация запроса
	jsonData, err := json.Marshal(transakcija)
	if err != nil {
		return nil, fmt.Errorf("ошибка сериализации: %v", err)
	}

	// Создание HTTP клиента
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Для тестового окружения
			},
		},
	}

	// Создание запроса
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %v", err)
	}

	// Установка заголовков
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json")

	// Выполнение запроса
	fmt.Printf("Отправка запроса к %s...\n", url)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	// Чтение ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	fmt.Printf("Статус ответа: %d\n", resp.StatusCode)
	fmt.Printf("Тело ответа: %s\n", string(body))

	// Декодирование ответа
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа: %v", err)
	}

	return result, nil
}

// testTracking тестирует отслеживание посылки (транзакция 63)
func testTracking() {
	fmt.Println("\n=== ТЕСТ 1: Отслеживание посылки (транзакция 63) ===")

	// Данные клиента
	klijent := Klijent{
		Username:          "b2b@svetu.rs",
		Password:          "Sv5et@U!",
		Jezik:             "SRB",
		IdTipUredjaja:     "2", // Веб-приложение
		VerzijaOS:         "Linux",
		NazivUredjaja:     "SVETU-API",
		ModelUredjaja:     "API",
		VerzijaAplikacije: "1.0.0",
		IPAdresa:          "192.168.1.1",
		IdPartnera:        10109, // ID партнера для SVETU
	}

	klijentJSON, _ := json.Marshal(klijent)

	// Данные для отслеживания
	tracking := TTKretanjeIn{
		VrstaUsluge:  1,
		EksterniBroj: "TEST123456",
		PrijemniBroj: "",
	}

	trackingJSON, _ := json.Marshal(tracking)

	// Создание транзакции
	transakcija := TransakcijaIn{
		StrKlijent:         string(klijentJSON),
		Servis:             101, // Для B2B партнеров
		IdVrstaTransakcije: 63,  // Отслеживание посылки
		TipSerijalizacije:  2,   // JSON
		IdTransakcija:      generateGUID(),
		StrIn:              string(trackingJSON),
	}

	// Вызов API
	result, err := callWSPAPI(&transakcija)
	if err != nil {
		fmt.Printf("❌ Ошибка: %v\n", err)
		return
	}

	// Анализ результата
	if rezultat, ok := result["Rezultat"]; ok {
		fmt.Printf("✅ Rezultat: %v\n", rezultat)
	}
	if strRezultat, ok := result["StrRezultat"]; ok {
		fmt.Printf("📝 StrRezultat: %v\n", strRezultat)
	}
	if greska, ok := result["Greska"]; ok {
		fmt.Printf("⚠️ Greska: %v\n", greska)
	}
}

// testManifest тестирует создание манифеста (транзакция 73)
func testManifest() {
	fmt.Println("\n=== ТЕСТ 2: Создание манифеста (транзакция 73) ===")

	// Данные клиента
	klijent := Klijent{
		Username:          "b2b@svetu.rs",
		Password:          "Sv5et@U!",
		Jezik:             "SRB",
		IdTipUredjaja:     "2",
		VerzijaOS:         "Linux",
		NazivUredjaja:     "SVETU-API",
		ModelUredjaja:     "API",
		VerzijaAplikacije: "1.0.0",
		IPAdresa:          "192.168.1.1",
		IdPartnera:        10109,
	}

	klijentJSON, _ := json.Marshal(klijent)

	// ID населенного пункта Белграда
	beogradId := 1100000

	// Создание манифеста с исправленными данными
	manifest := ManifestIn{
		ExtIdManifest: fmt.Sprintf("TEST-%d", time.Now().Unix()),
		IdTipPosiljke: 1, // Пакет
		Porudzbine: []Porudzbina{
			{
				ExtIdPorudzbina: fmt.Sprintf("ORDER-%d", time.Now().Unix()),
				Posiljke: []Posiljka{
					{
						Rbr:               1,
						ImaPrijemniBrojDN: "N",
						ExtBrend:          "SVETU",
						ExtMagacin:        "SVETU",
						ExtReferenca:      fmt.Sprintf("REF-%d", time.Now().Unix()),
						NacinPrijema:      "K", // Курьер
						MestoPreuzimanja: &Korisnik{ // Место преузимanja как объект Korisnik
							Vrsta:          "P",
							Naziv:          "SVETU PLATFORMA DOO",
							KontaktTelefon: "+381111234567",
							Adresa: Adresa{
								OznakaZemlje: "RS", // Добавлено: код страны
								IdNaselje: &beogradId,
								Naselje:   "Beograd",
								Ulica:     "Knez Mihailova",
								Broj:      "10",
								PostBroj:  "11000",
							},
						},
						IdRukovanje:       1,
						NacinPlacanja:     "U", // U - услуга уговорная (для B2B)
						Posiljalac: Korisnik{
							Vrsta:          "P", // Юридическое лицо
							Naziv:          "SVETU PLATFORMA DOO",
							KontaktTelefon: "+381111234567",
							KontaktOsoba:   "Test Manager",
							EMail:          "test@svetu.rs",
							Adresa: Adresa{
								OznakaZemlje: "RS", // Добавлено: код страны
								IdNaselje: &beogradId,
								Naselje:   "Beograd",
								Ulica:     "Knez Mihailova",
								Broj:      "10",
								PostBroj:  "11000",
							},
						},
						Primalac: Korisnik{
							Vrsta:          "F", // Физическое лицо
							Naziv:          "Petar Petrović", // Обязательно для физлиц
							Prezime:        "Petrović",
							Ime:            "Petar",
							KontaktTelefon: "+381611234567",
							EMail:          "petar@example.com",
							Adresa: Adresa{
								OznakaZemlje: "RS", // Добавлено: код страны
								IdNaselje: &beogradId,
								Naselje:   "Beograd",
								Ulica:     "Bulevar kralja Aleksandra",
								Broj:      "50",
								PostBroj:  "11000",
								Sprat:     "3",
								Stan:      "12",
							},
						},
						Masa:          1000, // 1 кг
						Vrednost:      5000, // 5000 динаров
						VrednostDTS:   5000,
						Otkupnina:     0,
						Sadrzaj:       "Odeća",
						PosebneUsluge: "PNA,SMS,VD", // PNA обязательно для курьера, VD для ценных посылок
					},
				},
			},
		},
	}

	manifestJSON, _ := json.Marshal(manifest)

	// Создание транзакции
	transakcija := TransakcijaIn{
		StrKlijent:         string(klijentJSON),
		Servis:             101,
		IdVrstaTransakcije: 73, // Создание манифеста
		TipSerijalizacije:  2,
		IdTransakcija:      generateGUID(),
		StrIn:              string(manifestJSON),
	}

	// Вызов API
	result, err := callWSPAPI(&transakcija)
	if err != nil {
		fmt.Printf("❌ Ошибка: %v\n", err)
		return
	}

	// Анализ результата
	if rezultat, ok := result["Rezultat"]; ok {
		fmt.Printf("✅ Rezultat: %v\n", rezultat)
	}
	if strOut, ok := result["StrOut"]; ok {
		fmt.Printf("📦 StrOut: %v\n", strOut)
		// Попытка декодировать StrOut
		if strOutStr, ok := strOut.(string); ok {
			var manifestOut map[string]interface{}
			if err := json.Unmarshal([]byte(strOutStr), &manifestOut); err == nil {
				if idPartner, ok := manifestOut["IdPartner"]; ok {
					fmt.Printf("🔑 IdPartner: %v\n", idPartner)
				}
			}
		}
	}
	if greske, ok := result["Greske"]; ok {
		fmt.Printf("⚠️ Greske: %v\n", greske)
	}
}

// testAddressInfo тестирует получение информации об адресах (транзакция 91)
func testAddressInfo() {
	fmt.Println("\n=== ТЕСТ 3: Получение информации об адресах (транзакция 91) ===")

	// Данные клиента
	klijent := Klijent{
		Username:          "b2b@svetu.rs",
		Password:          "Sv5et@U!",
		Jezik:             "SRB",
		IdTipUredjaja:     "2",
		VerzijaOS:         "Linux",
		NazivUredjaja:     "SVETU-API",
		ModelUredjaja:     "API",
		VerzijaAplikacije: "1.0.0",
		IPAdresa:          "192.168.1.1",
		IdPartnera:        10109,
	}

	klijentJSON, _ := json.Marshal(klijent)

	// Запрос адресов для Белграда
	addressRequest := map[string]interface{}{
		"IdVrstaInformacije": 1, // Населенные пункты
		"PostBroj":           "11000",
	}

	addressJSON, _ := json.Marshal(addressRequest)

	// Создание транзакции
	transakcija := TransakcijaIn{
		StrKlijent:         string(klijentJSON),
		Servis:             101,
		IdVrstaTransakcije: 91, // Информация об адресах
		TipSerijalizacije:  2,
		IdTransakcija:      generateGUID(),
		StrIn:              string(addressJSON),
	}

	// Вызов API
	result, err := callWSPAPI(&transakcija)
	if err != nil {
		fmt.Printf("❌ Ошибка: %v\n", err)
		return
	}

	// Анализ результата
	if rezultat, ok := result["Rezultat"]; ok {
		fmt.Printf("✅ Rezultat: %v\n", rezultat)
	}
	if strOut, ok := result["StrOut"]; ok {
		// Ограничиваем вывод первых 500 символов
		outStr := fmt.Sprintf("%v", strOut)
		if len(outStr) > 500 {
			outStr = outStr[:500] + "..."
		}
		fmt.Printf("📍 StrOut (первые 500 символов): %s\n", outStr)
	}
	if greska, ok := result["Greska"]; ok {
		fmt.Printf("⚠️ Greska: %v\n", greska)
	}
}

func main() {
	fmt.Println("🚀 ТЕСТИРОВАНИЕ WSP API POST EXPRESS")
	fmt.Println("====================================")
	fmt.Printf("Время: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("Учетные данные: b2b@svetu.rs")
	fmt.Println("Бренд/Магазин: SVETU")
	fmt.Println("ID партнера: 10109")
	fmt.Println("====================================")

	// Запуск тестов
	testTracking()
	testManifest()
	testAddressInfo()

	fmt.Println("\n✅ Тестирование завершено!")
}