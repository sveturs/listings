# 📮 ТЕХНИЧЕСКИЙ ОТЧЕТ ДЛЯ ПОДДЕРЖКИ POST EXPRESS

**Дата:** 2025-09-08  
**От:** Sve Tu Platforma  
**Кому:** Техническая поддержка Post Express (b2b@posta.rs)  
**Тема:** Проблема с учетными данными TEST в WSP API

## 📋 РЕЗЮМЕ ПРОБЛЕМЫ

Мы не можем успешно создать посылку через WSP API используя учетные данные **TEST/t3st** из вашей документации. API возвращает ошибку: **"Korisničko ime TEST nije registrovano!"** (Имя пользователя TEST не зарегистрировано).

## 🔧 НАША РЕАЛИЗАЦИЯ

### 1. Структура запроса (Go код)

```go
// backend/internal/proj/postexpress/models/models.go
type TransakcijaIn struct {
    StrKlijent         string `json:"StrKlijent"`
    Servis             int    `json:"Servis"`
    IdVrstaTranskacije int    `json:"IdVrstaTranskacije"` // С буквой "k"!
    TipSerijalizacije  int    `json:"TipSerijalizacije"`
    IdTransakcija      string `json:"IdTransakcija"`
    StrIn              string `json:"StrIn"`
}

type Klijent struct {
    Username      string `json:"Username"`
    Password      string `json:"Password"`
    Jezik         string `json:"Jezik"`
    IdTipUredjaja int    `json:"IdTipUredjaja"`
}
```

### 2. Функция отправки запроса

```go
// backend/internal/proj/postexpress/service/client.go
func (c *Client) SendRequest(req *TransakcijaIn) (*TransakcijaOut, error) {
    // Формируем Klijent
    klijent := Klijent{
        Username:      c.username,  // "TEST"
        Password:      c.password,  // "t3st"
        Jezik:         "LAT",
        IdTipUredjaja: 2,
    }
    
    klijentJSON, _ := json.Marshal(klijent)
    req.StrKlijent = string(klijentJSON)
    
    // Отправляем запрос
    jsonData, _ := json.Marshal(req)
    
    httpReq, _ := http.NewRequest("POST", 
        "http://212.62.32.201/WspWebApi/transakcija", 
        bytes.NewBuffer(jsonData))
    httpReq.Header.Set("Content-Type", "application/json")
    
    resp, _ := c.httpClient.Do(httpReq)
    // ... обработка ответа
}
```

## 📊 ТЕСТОВЫЕ ЗАПРОСЫ И ОТВЕТЫ

### Тест 1: Простая транзакция GetNaselje (ID=3)

**Запрос:**
```json
{
  "StrKlijent": "{\"Username\":\"TEST\",\"Password\":\"t3st\",\"Jezik\":\"LAT\",\"IdTipUredjaja\":2}",
  "Servis": 3,
  "IdVrstaTranskacije": 3,
  "TipSerijalizacije": 2,
  "IdTransakcija": "test-1736955123",
  "StrIn": "{\"Naziv\":\"Novi\"}"
}
```

**cURL команда:**
```bash
curl -X POST http://212.62.32.201/WspWebApi/transakcija \
  -H "Content-Type: application/json" \
  -d '{
    "StrKlijent": "{\"Username\":\"TEST\",\"Password\":\"t3st\",\"Jezik\":\"LAT\",\"IdTipUredjaja\":2}",
    "Servis": 3,
    "IdVrstaTranskacije": 3,
    "TipSerijalizacije": 2,
    "IdTransakcija": "test-1736955123",
    "StrIn": "{\"Naziv\":\"Novi\"}"
  }'
```

**Ответ:**
```json
{
  "Rezultat": 3,
  "StrOut": null,
  "StrRezultat": "{
    \"Poruka\": \"Korisničko ime TEST nije registrovano!\",
    \"PorukaKorisnik\": \"Korisničko ime TEST nije registrovano!\"
  }"
}
```

### Тест 2: Транзакция Manifest (ID=73) для создания посылки

**Запрос:**
```json
{
  "StrKlijent": "{\"Username\":\"TEST\",\"Password\":\"t3st\",\"Jezik\":\"LAT\",\"IdTipUredjaja\":2}",
  "Servis": 3,
  "IdVrstaTranskacije": 73,
  "TipSerijalizacije": 2,
  "IdTransakcija": "manifest-1736955456",
  "StrIn": "{
    \"Posiljalac\": {
      \"Ime\": \"Test Sender\",
      \"Adresa\": \"Bulevar kralja Aleksandra 1\",
      \"IdNaselje\": 110000,
      \"Telefon\": \"0601234567\"
    },
    \"Posiljke\": [{
      \"Primalac\": {
        \"Ime\": \"Test Receiver\",
        \"Adresa\": \"Knez Mihailova 10\",
        \"IdNaselje\": 110000,
        \"Telefon\": \"0607654321\"
      },
      \"TezinaPosiljke\": 1000,
      \"VrednostPosiljke\": 500000,
      \"BrojOtkupnice\": \"123456\",
      \"Sadrzaj\": \"Test package\"
    }],
    \"DatumPrijema\": \"2025-01-08\"
  }"
}
```

**Ответ:**
```json
{
  "Rezultat": 3,
  "StrOut": null,
  "StrRezultat": "{
    \"Poruka\": \"Korisničko ime TEST nije registrovano!\",
    \"PorukaKorisnik\": \"Korisničko ime TEST nije registrovano!\"
  }"
}
```

## 🔍 ЧТО МЫ ПРОВЕРИЛИ

### ✅ Правильность полей
- Используем `IdVrstaTranskacije` с буквой "k" (не "c")
- Все поля с большой буквы как в документации
- `Servis = 3` для B2B
- `TipSerijalizacije = 2` для JSON

### ✅ Различные варианты авторизации
Протестировали несколько вариантов формирования Klijent:

1. **Вариант из документации:**
```json
{"Username":"TEST","Password":"t3st","Jezik":"LAT","IdTipUredjaja":2}
```

2. **С маленькими буквами:**
```json
{"username":"TEST","password":"t3st","jezik":"LAT","idTipUredjaja":2}
```

3. **Различные комбинации регистра:**
- Username: TEST, test, Test
- Password: t3st, T3ST, T3st

**Результат:** Все варианты возвращают "Korisničko ime TEST nije registrovano!"

## 📝 ЛОГИ ПОЛНОГО ТЕСТА

```bash
$ go run backend/scripts/test_wsp_minimal.go

========================================
  MINIMAL WSP API TEST
========================================

1. Testing IdVrstaTransakcije (with 'c')
Request: {"IdTransakcija":"test-1757341551","IdVrstaTransakcije":3,...}
Status: 200
Error message: Nepoznata vrsta transakcije (NapraviObjIn)! IdVrstaTransakcije = 0
❌ Failed with Rezultat=3

2. Testing IdVrstaTranskacije (with 'k')
Request: {"IdTransakcija":"test-1757341552","IdVrstaTranskacije":3,...}
Status: 200
Error message: Korisničko ime TEST nije registrovano!
❌ Failed with Rezultat=3

[... остальные тесты с тем же результатом ...]
```

## 🎯 КЛЮЧЕВОЕ НАБЛЮДЕНИЕ

Когда мы используем **неправильное** поле `IdVrstaTransakcije` (с "c"), получаем:
```
"Nepoznata vrsta transakcije (NapraviObjIn)! IdVrstaTransakcije = 0"
```

Когда используем **правильное** поле `IdVrstaTranskacije` (с "k"), получаем:
```
"Korisničko ime TEST nije registrovano!"
```

Это доказывает, что:
1. ✅ Наш запрос структурирован правильно
2. ✅ API правильно парсит наш запрос
3. ✅ Транзакция распознается корректно
4. ❌ Учетные данные TEST/t3st не активны или не существуют

## 📊 СРАВНЕНИЕ С ДОКУМЕНТАЦИЕЙ

### Из вашей документации (https://www.posta.rs/wsp-help/uvod/uvod.aspx):
- Username: TEST
- Password: t3st
- Test URL: http://212.62.32.201/WspWebApi/transakcija

### Что мы используем:
- ✅ Username: TEST (точное совпадение)
- ✅ Password: t3st (точное совпадение)  
- ✅ URL: http://212.62.32.201/WspWebApi/transakcija (точное совпадение)

## ❓ ВОПРОСЫ

1. **Активны ли учетные данные TEST/t3st в тестовой среде?**
   - Возможно, они были деактивированы?
   - Требуется ли предварительная регистрация?

2. **Есть ли дополнительные требования для активации?**
   - IP whitelist?
   - Предварительная регистрация через email?
   - Специальные headers в запросе?

3. **Можете предоставить рабочий пример запроса?**
   - С активными тестовыми credentials
   - Который успешно создает посылку

---

**P.S.** Мы уверены, что у других клиентов интеграция работает успешно. Просим помочь нам найти, что мы делаем не так. Возможно, есть какой-то неочевидный шаг, который мы пропускаем?

## 🔗 ПРИЛОЖЕНИЯ

### Полный тестовый скрипт (Go)
```go
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
    
    request := map[string]interface{}{
        "StrKlijent":         `{"Username":"TEST","Password":"t3st","Jezik":"LAT","IdTipUredjaja":2}`,
        "Servis":             3,
        "IdVrstaTranskacije": 3,
        "TipSerijalizacije":  2,
        "IdTransakcija":      fmt.Sprintf("test-%d", time.Now().Unix()),
        "StrIn":              `{"Naziv":"Novi"}`,
    }
    
    jsonData, _ := json.Marshal(request)
    fmt.Printf("Request: %s\n", string(jsonData))
    
    req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    
    client := &http.Client{Timeout: 10 * time.Second}
    resp, _ := client.Do(req)
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    fmt.Printf("Response: %s\n", string(body))
}
```

### Команда для быстрого теста
```bash
curl -X POST http://212.62.32.201/WspWebApi/transakcija \
  -H "Content-Type: application/json" \
  -d '{"StrKlijent":"{\"Username\":\"TEST\",\"Password\":\"t3st\",\"Jezik\":\"LAT\",\"IdTipUredjaja\":2}","Servis":3,"IdVrstaTranskacije":3,"TipSerijalizacije":2,"IdTransakcija":"test-123","StrIn":"{\"Naziv\":\"Novi\"}"}'
```
