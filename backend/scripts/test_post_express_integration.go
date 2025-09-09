package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// TransakcijaIn - основная структура запроса
type TransakcijaIn struct {
	StrKlijent         string `json:"strKlijent"`
	Servis             int    `json:"servis"`
	IdVrstaTransakcije int    `json:"idVrstaTransakcije"`
	TipSerijalizacije  int    `json:"tipSerijalizacije"`
	IdTransakcija      string `json:"idTransakcija"`
	StrIn              string `json:"strIn,omitempty"`
}

// ClientData - данные клиента для аутентификации
type ClientData struct {
	Username          string `json:"username"`
	Password          string `json:"password"`
	Jezik             string `json:"jezik"`
	IdTipUredjaja     int    `json:"idTipUredjaja"`
	VerzijaOS         string `json:"verzijaOS"`
	NazivUredjaja     string `json:"nazivUredjaja"`
	ModelUredjaja     string `json:"modelUredjaja"`
	VerzijaAplikacije string `json:"verzijaAplikacije"`
	IPAdresa          string `json:"ipAdresa"`
}

// ManifestRequest - запрос для создания манифеста
type ManifestRequest struct {
	Naziv       string `json:"naziv"`
	Napomena    string `json:"napomena"`
	IdMagacin   int    `json:"idMagacin"`
	IdTipPosiljke int  `json:"idTipPosiljke"`
}

// TrackingRequest - запрос для отслеживания
type TrackingRequest struct {
	BrojPosiljke string `json:"brojPosiljke"`
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

	fmt.Println("🚀 Post Express WSP API Test")
	fmt.Println("============================")
	fmt.Printf("Endpoint: %s\n", endpoint)
	fmt.Printf("Username: %s\n", username)
	fmt.Printf("Password: %s\n", password[:3]+"...")
	fmt.Println("")

	// Создаем HTTP клиент
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Для тестовой среды
			},
		},
	}

	// Тест 1: Простая аутентификация (транзакция 61 - получение списка магазинов)
	fmt.Println("📋 Тест 1: Получение списка магазинов (транзакция 61)")
	fmt.Println("-------------------------------------------------------")
	testGetStores(client, endpoint, username, password)

	// Тест 2: Создание манифеста (транзакция 73)
	fmt.Println("\n📦 Тест 2: Создание манифеста (транзакция 73)")
	fmt.Println("-----------------------------------------------")
	testCreateManifest(client, endpoint, username, password)

	// Тест 3: Отслеживание посылки (транзакция 63)
	fmt.Println("\n🔍 Тест 3: Отслеживание посылки (транзакция 63)")
	fmt.Println("-------------------------------------------------")
	testTracking(client, endpoint, username, password)

	// Тест 4: Получение типов посылок (транзакция 58)
	fmt.Println("\n📋 Тест 4: Получение типов посылок (транзакция 58)")
	fmt.Println("----------------------------------------------------")
	testGetShipmentTypes(client, endpoint, username, password)
}

func testGetStores(client *http.Client, endpoint, username, password string) {
	// Подготовка клиентских данных
	clientData := ClientData{
		Username:          username,
		Password:          password,
		Jezik:             "LAT",
		IdTipUredjaja:     2,
		VerzijaOS:         "Linux",
		NazivUredjaja:     "API",
		ModelUredjaja:     "SERVER",
		VerzijaAplikacije: "1.0.0",
		IPAdresa:          "127.0.0.1",
	}

	clientJSON, _ := json.Marshal(clientData)

	// Создаем запрос
	req := TransakcijaIn{
		StrKlijent:         string(clientJSON),
		Servis:             3,
		IdVrstaTransakcije: 61, // GetMagacini
		TipSerijalizacije:  2,   // JSON
		IdTransakcija:      generateGUID(),
	}

	// Отправляем запрос
	sendRequest(client, endpoint, req)
}

func testCreateManifest(client *http.Client, endpoint, username, password string) {
	// Подготовка клиентских данных
	clientData := ClientData{
		Username:          username,
		Password:          password,
		Jezik:             "LAT",
		IdTipUredjaja:     2,
		VerzijaOS:         "Linux",
		NazivUredjaja:     "API",
		ModelUredjaja:     "SERVER",
		VerzijaAplikacije: "1.0.0",
		IPAdresa:          "127.0.0.1",
	}

	clientJSON, _ := json.Marshal(clientData)

	// Подготовка данных манифеста
	manifestData := ManifestRequest{
		Naziv:         "Test Manifest SVETU",
		Napomena:      "Тестовый манифест для проверки интеграции",
		IdMagacin:     1, // Будем получать из транзакции 61
		IdTipPosiljke: 1, // Стандартная посылка
	}

	manifestJSON, _ := json.Marshal(manifestData)

	// Создаем запрос
	req := TransakcijaIn{
		StrKlijent:         string(clientJSON),
		Servis:             3,
		IdVrstaTransakcije: 73, // Manifest
		TipSerijalizacije:  2,   // JSON
		IdTransakcija:      generateGUID(),
		StrIn:              string(manifestJSON),
	}

	// Отправляем запрос
	sendRequest(client, endpoint, req)
}

func testTracking(client *http.Client, endpoint, username, password string) {
	// Подготовка клиентских данных
	clientData := ClientData{
		Username:          username,
		Password:          password,
		Jezik:             "LAT",
		IdTipUredjaja:     2,
		VerzijaOS:         "Linux",
		NazivUredjaja:     "API",
		ModelUredjaja:     "SERVER",
		VerzijaAplikacije: "1.0.0",
		IPAdresa:          "127.0.0.1",
	}

	clientJSON, _ := json.Marshal(clientData)

	// Подготовка данных для отслеживания
	trackingData := TrackingRequest{
		BrojPosiljke: "TEST123456", // Тестовый номер
	}

	trackingJSON, _ := json.Marshal(trackingData)

	// Создаем запрос
	req := TransakcijaIn{
		StrKlijent:         string(clientJSON),
		Servis:             3,
		IdVrstaTransakcije: 63, // Tracking
		TipSerijalizacije:  2,   // JSON
		IdTransakcija:      generateGUID(),
		StrIn:              string(trackingJSON),
	}

	// Отправляем запрос
	sendRequest(client, endpoint, req)
}

func testGetShipmentTypes(client *http.Client, endpoint, username, password string) {
	// Подготовка клиентских данных
	clientData := ClientData{
		Username:          username,
		Password:          password,
		Jezik:             "LAT",
		IdTipUredjaja:     2,
		VerzijaOS:         "Linux",
		NazivUredjaja:     "API",
		ModelUredjaja:     "SERVER",
		VerzijaAplikacije: "1.0.0",
		IPAdresa:          "127.0.0.1",
	}

	clientJSON, _ := json.Marshal(clientData)

	// Создаем запрос
	req := TransakcijaIn{
		StrKlijent:         string(clientJSON),
		Servis:             3,
		IdVrstaTransakcije: 58, // GetTipoviPosiljki
		TipSerijalizacije:  2,   // JSON
		IdTransakcija:      generateGUID(),
	}

	// Отправляем запрос
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
	fmt.Println(string(reqJSON))

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
	} else {
		fmt.Println("📄 Тело ответа (raw):")
		fmt.Println(string(body))
	}

	// Анализ результата
	if resp.StatusCode == 200 {
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
				if poruka, ok := errMsg["Poruka"].(string); ok {
					fmt.Printf("💬 Сообщение: %s\n", poruka)
				}
				if porukaKorisnik, ok := errMsg["PorukaKorisnik"].(string); ok {
					fmt.Printf("💬 Сообщение для пользователя: %s\n", porukaKorisnik)
				}
			} else {
				fmt.Printf("💬 Сообщение: %s\n", strRezultat)
			}
		}
	} else {
		fmt.Printf("❌ HTTP ошибка: %d\n", resp.StatusCode)
	}
}

func generateGUID() string {
	// Простая генерация GUID для тестов
	return fmt.Sprintf("%d-%d-%d-%d", 
		time.Now().Unix(), 
		time.Now().Nanosecond(),
		os.Getpid(),
		time.Now().UnixNano()%1000)
}