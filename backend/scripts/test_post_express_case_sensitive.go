package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	endpoint := "http://212.62.32.201/WspWebApi/transakcija"
	
	fmt.Println("=====================================")
	fmt.Println("Post Express Case Sensitivity Test")
	fmt.Println("=====================================")
	fmt.Printf("Endpoint: %s\n", endpoint)
	fmt.Println("=====================================\n")

	// Подготовка данных клиента в разных вариантах
	clientData := map[string]interface{}{
		"Username":          "b2b@svetu.rs",
		"Password":          "Vu3@#45$%67&*89()",
		"Jezik":             "sr",
		"IdTipUredjaja":     "2",
		"VerzijaOS":         "Linux",
		"NazivUredjaja":     "SveTu-Server",
		"ModelUredjaja":     "API",
		"VerzijaAplikacije": "1.0.0",
		"IPAdresa":          "127.0.0.1",
		"IdPartnera":        10109,
	}
	
	clientJSON, _ := json.Marshal(clientData)
	
	// Простой поисковый запрос
	searchData := map[string]interface{}{
		"Naziv":           "Novi Sad",
		"BrojSlogova":     10,
		"NacinSortiranja": 0,
	}
	
	searchJSON, _ := json.Marshal(searchData)
	
	// Тест 1: camelCase (все с маленькой буквы)
	fmt.Println("Test 1: camelCase...")
	transaction1 := map[string]interface{}{
		"strKlijent":         string(clientJSON),
		"servis":             101,
		"idVrstaTransakcije": 3,
		"tipSerijalizacije":  2,
		"idTransakcija":      fmt.Sprintf("TEST1-%d", time.Now().Unix()),
		"strIn":              string(searchJSON),
	}
	testRequest(endpoint, transaction1)
	
	// Тест 2: PascalCase (все с большой буквы)
	fmt.Println("\nTest 2: PascalCase...")
	transaction2 := map[string]interface{}{
		"StrKlijent":         string(clientJSON),
		"Servis":             101,
		"IdVrstaTransakcije": 3,
		"TipSerijalizacije":  2,
		"IdTransakcija":      fmt.Sprintf("TEST2-%d", time.Now().Unix()),
		"StrIn":              string(searchJSON),
	}
	testRequest(endpoint, transaction2)
	
	// Тест 3: lowercase (все маленькими)
	fmt.Println("\nTest 3: lowercase...")
	transaction3 := map[string]interface{}{
		"strklijent":         string(clientJSON),
		"servis":             101,
		"idvrstatransakcije": 3,
		"tipserijalizacije":  2,
		"idtransakcija":      fmt.Sprintf("TEST3-%d", time.Now().Unix()),
		"strin":              string(searchJSON),
	}
	testRequest(endpoint, transaction3)
	
	// Тест 4: UPPERCASE (все большими)
	fmt.Println("\nTest 4: UPPERCASE...")
	transaction4 := map[string]interface{}{
		"STRKLIJENT":         string(clientJSON),
		"SERVIS":             101,
		"IDVRSTATRANSAKCIJE": 3,
		"TIPSERIJALIZACIJE":  2,
		"IDTRANSAKCIJA":      fmt.Sprintf("TEST4-%d", time.Now().Unix()),
		"STRIN":              string(searchJSON),
	}
	testRequest(endpoint, transaction4)
	
	// Тест 5: Смешанный вариант (как в документации из отчета)
	fmt.Println("\nTest 5: Mixed case from documentation...")
	transaction5 := map[string]interface{}{
		"strKlijent":         string(clientJSON),
		"servis":             101,
		"idVrstaTransakcije": 3,
		"tipSerijalizacije":  2,
		"idTransakcija":      fmt.Sprintf("TEST5-%d", time.Now().Unix()),
		"strIn":              string(searchJSON),
	}
	testRequest(endpoint, transaction5)
	
	fmt.Println("\n=====================================")
	fmt.Println("Case sensitivity test completed!")
	fmt.Println("=====================================")
}

func testRequest(endpoint string, transaction map[string]interface{}) {
	jsonData, _ := json.Marshal(transaction)
	
	// Показываем только ключи для краткости
	fmt.Printf("Keys used: ")
	for key := range transaction {
		fmt.Printf("%s ", key)
	}
	fmt.Println()
	
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ HTTP Error: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("Response: %s\n", string(body))
		return
	}
	
	// Анализируем результат
	if rezultat, ok := result["Rezultat"]; ok {
		if rezultat == float64(3) {
			// Ошибка - показываем сообщение
			if strRez, ok := result["StrRezultat"].(string); ok {
				var errMsg map[string]interface{}
				if json.Unmarshal([]byte(strRez), &errMsg) == nil {
					if poruka, ok := errMsg["Poruka"].(string); ok {
						// Проверяем, содержит ли ошибка "IdVrstaTransakcije = 0"
						if bytes.Contains([]byte(poruka), []byte("IdVrstaTransakcije = 0")) {
							fmt.Printf("❌ Transaction type not recognized (IdVrstaTransakcije = 0)\n")
						} else {
							fmt.Printf("❌ Error: %s\n", poruka)
						}
					}
				}
			}
		} else if rezultat == float64(0) || rezultat == float64(1) {
			fmt.Printf("✅ Success! Rezultat = %v\n", rezultat)
			if strOut, ok := result["StrOut"]; ok && strOut != nil {
				fmt.Printf("   Response data received\n")
			}
		} else {
			fmt.Printf("⚠️ Unknown result: %v\n", rezultat)
		}
	} else {
		fmt.Printf("⚠️ Unexpected response format: %v\n", result)
	}
}