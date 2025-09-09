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

// TransakcijaIn - основная структура запроса с большими буквами
type TransakcijaIn struct {
	StrKlijent         string `json:"StrKlijent"`
	Servis             int    `json:"Servis"`
	IdVrstaTransakcije int    `json:"IdVrstaTransakcije"`
	TipSerijalizacije  int    `json:"TipSerijalizacije"`
	IdTransakcija      string `json:"IdTransakcija"`
	StrIn              string `json:"StrIn,omitempty"`
}

// ClientData - данные клиента с большими буквами
type ClientData struct {
	Username          string `json:"Username"`
	Password          string `json:"Password"`
	Jezik             string `json:"Jezik"`
	IdTipUredjaja     int    `json:"IdTipUredjaja"`
	VerzijaOS         string `json:"VerzijaOS"`
	NazivUredjaja     string `json:"NazivUredjaja"`
	ModelUredjaja     string `json:"ModelUredjaja"`
	VerzijaAplikacije string `json:"VerzijaAplikacije"`
	IPAdresa          string `json:"IPAdresa"`
}

func main() {
	// Читаем учетные данные
	username := "b2b@svetu.rs"
	password := "Sv5et@U!"
	endpoint := "http://212.62.32.201/WspWebApi/transakcija"

	fmt.Println("🚀 Post Express WSP API Test (Capital Letters)")
	fmt.Println("==============================================")
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

	// Тест с разными комбинациями
	tests := []struct {
		name        string
		transaction int
		servis      int
	}{
		{"Получение магазинов (Servis=3)", 61, 3},
		{"Получение магазинов (Servis=1)", 61, 1},
		{"Получение магазинов (Servis=101)", 61, 101},
		{"Получение типов посылок (Servis=3)", 58, 3},
		{"Простой тест подключения (transaction=1)", 1, 3},
	}

	for _, test := range tests {
		fmt.Printf("\n📋 Тест: %s\n", test.name)
		fmt.Println(strings.Repeat("-", 50))
		
		// Подготовка клиентских данных
		clientData := ClientData{
			Username:          username,
			Password:          password,
			Jezik:             "LAT",
			IdTipUredjaja:     2,
			VerzijaOS:         "Linux",
			NazivUredjaja:     "SVETU",
			ModelUredjaja:     "SERVER",
			VerzijaAplikacije: "1.0.0",
			IPAdresa:          "127.0.0.1",
		}

		clientJSON, _ := json.Marshal(clientData)

		// Создаем запрос
		req := TransakcijaIn{
			StrKlijent:         string(clientJSON),
			Servis:             test.servis,
			IdVrstaTransakcije: test.transaction,
			TipSerijalizacije:  2,
			IdTransakcija:      generateGUID(),
		}

		sendRequest(client, endpoint, req)
	}

	// Тест с альтернативной структурой
	fmt.Println("\n📋 Тест: Альтернативный формат запроса")
	fmt.Println("----------------------------------------")
	testAlternativeFormat(client, endpoint, username, password)
}

func testAlternativeFormat(client *http.Client, endpoint, username, password string) {
	// Пробуем отправить как простой объект без вложенности
	requestBody := map[string]interface{}{
		"Username":           username,
		"Password":           password,
		"Jezik":              "LAT",
		"IdVrstaTransakcije": 61,
		"Servis":             3,
		"TipSerijalizacije":  2,
	}

	reqJSON, _ := json.MarshalIndent(requestBody, "", "  ")
	fmt.Println("📤 Отправляемый запрос (альтернативный формат):")
	fmt.Println(string(reqJSON))

	httpReq, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(reqJSON))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := client.Do(httpReq)
	if err != nil {
		fmt.Printf("❌ Ошибка отправки: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("\n📥 Статус ответа: %s\n", resp.Status)
	fmt.Println("📄 Тело ответа:")
	
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err == nil {
		prettyJSON, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(prettyJSON))
	} else {
		fmt.Println(string(body))
	}
}

func sendRequest(client *http.Client, endpoint string, req TransakcijaIn) {
	reqJSON, _ := json.MarshalIndent(req, "", "  ")
	fmt.Println("📤 Отправляемый запрос:")
	fmt.Println(string(reqJSON))

	httpReq, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(reqJSON))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := client.Do(httpReq)
	if err != nil {
		fmt.Printf("❌ Ошибка отправки: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("\n📥 Статус ответа: %s\n", resp.Status)
	
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
			}
		}

		if strRezultat, ok := result["StrRezultat"].(string); ok && strRezultat != "" {
			var errMsg map[string]interface{}
			if err := json.Unmarshal([]byte(strRezultat), &errMsg); err == nil {
				if poruka, ok := errMsg["Poruka"].(string); ok {
					fmt.Printf("💬 Сообщение: %s\n", poruka)
				}
			}
		}
	} else {
		fmt.Println("📄 Тело ответа (raw):")
		fmt.Println(string(body))
	}
}

func generateGUID() string {
	return fmt.Sprintf("%d-%d-%d", 
		time.Now().Unix(), 
		time.Now().Nanosecond(),
		os.Getpid())
}

// Добавим импорт strings
var strings = struct {
	Repeat func(string, int) string
}{
	Repeat: func(s string, count int) string {
		result := ""
		for i := 0; i < count; i++ {
			result += s
		}
		return result
	},
}