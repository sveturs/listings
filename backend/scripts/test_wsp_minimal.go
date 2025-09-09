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
	fmt.Println("========================================")
	fmt.Println("  MINIMAL WSP API TEST")
	fmt.Println("========================================")
	
	endpoint := "http://212.62.32.201/WspWebApi/transakcija"
	
	// Тест 1: Минимальный запрос с разными вариантами написания
	testVariants(endpoint)
}

func testVariants(endpoint string) {
	// Вариант 1: IdVrstaTransakcije (с "c")
	fmt.Println("\n1. Testing IdVrstaTransakcije (with 'c')")
	request1 := map[string]interface{}{
		"StrKlijent":         `{"Username":"TEST","Password":"t3st","Jezik":"LAT","IdTipUredjaja":2}`,
		"Servis":             3,
		"IdVrstaTransakcije": 3,  // С буквой "c"
		"TipSerijalizacije":  2,
		"IdTransakcija":      fmt.Sprintf("test-%d", time.Now().Unix()),
		"StrIn":              `{"Naziv":"Novi"}`,
	}
	sendAndPrint(endpoint, request1)
	
	// Вариант 2: IdVrstaTranskacije (с "k")
	fmt.Println("\n2. Testing IdVrstaTranskacije (with 'k')")
	request2 := map[string]interface{}{
		"StrKlijent":         `{"Username":"TEST","Password":"t3st","Jezik":"LAT","IdTipUredjaja":2}`,
		"Servis":             3,
		"IdVrstaTranskacije": 3,  // С буквой "k"
		"TipSerijalizacije":  2,
		"IdTransakcija":      fmt.Sprintf("test-%d", time.Now().Unix()+1),
		"StrIn":              `{"Naziv":"Novi"}`,
	}
	sendAndPrint(endpoint, request2)
	
	// Вариант 3: Без сериализации Klijent (прямые поля)
	fmt.Println("\n3. Testing with direct Klijent fields")
	request3 := map[string]interface{}{
		"Username":           "TEST",
		"Password":           "t3st",
		"Jezik":              "LAT",
		"IdTipUredjaja":      2,
		"Servis":             3,
		"IdVrstaTransakcije": 3,
		"TipSerijalizacije":  2,
		"IdTransakcija":      fmt.Sprintf("test-%d", time.Now().Unix()+2),
		"StrIn":              `{"Naziv":"Novi"}`,
	}
	sendAndPrint(endpoint, request3)
	
	// Вариант 4: Как в старой документации (все поля с маленькой буквы)
	fmt.Println("\n4. Testing with lowercase fields")
	request4 := map[string]interface{}{
		"strKlijent":         `{"username":"TEST","password":"t3st","jezik":"LAT","idTipUredjaja":2}`,
		"servis":             3,
		"idVrstaTransakcije": 3,
		"tipSerijalizacije":  2,
		"idTransakcija":      fmt.Sprintf("test-%d", time.Now().Unix()+3),
		"strIn":              `{"naziv":"Novi"}`,
	}
	sendAndPrint(endpoint, request4)
	
	// Вариант 5: TipSerijalizacije = 1 (возможно для JSON это 1, а не 2)
	fmt.Println("\n5. Testing with TipSerijalizacije = 1")
	request5 := map[string]interface{}{
		"StrKlijent":         `{"Username":"TEST","Password":"t3st","Jezik":"LAT","IdTipUredjaja":2}`,
		"Servis":             3,
		"IdVrstaTransakcije": 3,
		"TipSerijalizacije":  1,  // Попробуем 1 вместо 2
		"IdTransakcija":      fmt.Sprintf("test-%d", time.Now().Unix()+4),
		"StrIn":              `{"Naziv":"Novi"}`,
	}
	sendAndPrint(endpoint, request5)
	
	// Вариант 6: Servis = 101 (как в некоторых примерах)
	fmt.Println("\n6. Testing with Servis = 101")
	request6 := map[string]interface{}{
		"StrKlijent":         `{"Username":"TEST","Password":"t3st","Jezik":"LAT","IdTipUredjaja":2}`,
		"Servis":             101,  // Попробуем 101
		"IdVrstaTransakcije": 3,
		"TipSerijalizacije":  2,
		"IdTransakcija":      fmt.Sprintf("test-%d", time.Now().Unix()+5),
		"StrIn":              `{"Naziv":"Novi"}`,
	}
	sendAndPrint(endpoint, request6)
}

func sendAndPrint(endpoint string, request map[string]interface{}) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		fmt.Printf("Error marshaling: %v\n", err)
		return
	}
	
	fmt.Printf("Request: %s\n", string(jsonData))
	
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading: %v\n", err)
		return
	}
	
	fmt.Printf("Status: %d\n", resp.StatusCode)
	
	// Попробуем распарсить ответ
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err == nil {
		// Извлечем StrRezultat если есть
		if strRez, ok := result["StrRezultat"].(string); ok {
			var innerResult map[string]interface{}
			if err := json.Unmarshal([]byte(strRez), &innerResult); err == nil {
				if msg, ok := innerResult["Poruka"].(string); ok {
					fmt.Printf("Error message: %s\n", msg)
				}
			}
		}
		
		// Проверим успешность
		if rezultat, ok := result["Rezultat"].(float64); ok {
			if rezultat == 0 {
				fmt.Printf("✅ SUCCESS! Response: %s\n", string(body))
			} else {
				fmt.Printf("❌ Failed with Rezultat=%v\n", rezultat)
			}
		}
	} else {
		fmt.Printf("Raw response: %s\n", string(body))
	}
	
	fmt.Println("---")
}