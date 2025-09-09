#  TEHNIČKI IZVEŠTAJ ZA PODRŠKU POST EXPRESS

**Datum:** 2025-09-08  
**Od:** Sve Tu Platforma  
**Za:** Tehničku podršku Post Express (b2b@posta.rs)  
**Tema:** Problem sa pristupnim podacima TEST u WSP API

##  REZIME PROBLEMA

Ne možemo uspešno kreirati pošiljku kroz WSP API koristeći pristupne podatke **TEST/t3st** iz vaše dokumentacije. API vraća grešku: **"Korisničko ime TEST nije registrovano!"**

##  NAŠA IMPLEMENTACIJA

### 1. Struktura zahteva (Go kod)

```go
// backend/internal/proj/postexpress/models/models.go
type TransakcijaIn struct {
    StrKlijent         string `json:"StrKlijent"`
    Servis             int    `json:"Servis"`
    IdVrstaTranskacije int    `json:"IdVrstaTranskacije"` // Sa slovom "k"!
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

### 2. Funkcija slanja zahteva

```go
// backend/internal/proj/postexpress/service/client.go
func (c *Client) SendRequest(req *TransakcijaIn) (*TransakcijaOut, error) {
    // Formiramo Klijent
    klijent := Klijent{
        Username:      c.username,  // "TEST"
        Password:      c.password,  // "t3st"
        Jezik:         "LAT",
        IdTipUredjaja: 2,
    }
    
    klijentJSON, _ := json.Marshal(klijent)
    req.StrKlijent = string(klijentJSON)
    
    // Šaljemo zahtev
    jsonData, _ := json.Marshal(req)
    
    httpReq, _ := http.NewRequest("POST", 
        "http://212.62.32.201/WspWebApi/transakcija", 
        bytes.NewBuffer(jsonData))
    httpReq.Header.Set("Content-Type", "application/json")
    
    resp, _ := c.httpClient.Do(httpReq)
    // ... obrada odgovora
}
```

##  TEST ZAHTEVI I ODGOVORI

### Test 1: Prosta transakcija GetNaselje (ID=3)

**Zahtev:**
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

**cURL komanda:**
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

**Odgovor:**
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

### Test 2: Transakcija Manifest (ID=73) za kreiranje pošiljke

**Zahtev:**
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

**Odgovor:**
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

##  ŠTA SMO PROVERILI

###  Ispravnost polja
- Koristimo `IdVrstaTranskacije` sa slovom "k" (ne "c")
- Sva polja sa velikim slovom kao u dokumentaciji
- `Servis = 3` za B2B
- `TipSerijalizacije = 2` za JSON

###  Različite varijante autorizacije
Testirali smo nekoliko varijanti formiranja Klijent:

1. **Varijanta iz dokumentacije:**
```json
{"Username":"TEST","Password":"t3st","Jezik":"LAT","IdTipUredjaja":2}
```

2. **Sa malim slovima:**
```json
{"username":"TEST","password":"t3st","jezik":"LAT","idTipUredjaja":2}
```

3. **Različite kombinacije velikih/malih slova:**
- Username: TEST, test, Test
- Password: t3st, T3ST, T3st

**Rezultat:** Sve varijante vraćaju "Korisničko ime TEST nije registrovano!"

##  LOGOVI KOMPLETNOG TESTA

```bash
$ go run backend/scripts/test_wsp_minimal.go

========================================
  MINIMAL WSP API TEST
========================================

1. Testing IdVrstaTransakcije (with 'c')
Request: {"IdTransakcija":"test-1757341551","IdVrstaTransakcije":3,...}
Status: 200
Error message: Nepoznata vrsta transakcije (NapraviObjIn)! IdVrstaTransakcije = 0
## Failed with Rezultat=3

2. Testing IdVrstaTranskacije (with 'k')
Request: {"IdTransakcija":"test-1757341552","IdVrstaTranskacije":3,...}
Status: 200
Error message: Korisničko ime TEST nije registrovano!
## Failed with Rezultat=3

[... ostali testovi sa istim rezultatom ...]
```

##  KLJUČNO ZAPAŽANJE

Kada koristimo **pogrešno** polje `IdVrstaTransakcije` (sa "c"), dobijamo:
```
"Nepoznata vrsta transakcije (NapraviObjIn)! IdVrstaTransakcije = 0"
```

Kada koristimo **ispravno** polje `IdVrstaTranskacije` (sa "k"), dobijamo:
```
"Korisničko ime TEST nije registrovano!"
```

Ovo dokazuje da:
1.  Naš zahtev je pravilno strukturiran
2.  API ispravno parsira naš zahtev
3.  Transakcija se prepoznaje korektno
4.  Pristupni podaci TEST/t3st nisu aktivni ili ne postoje

##  POREĐENJE SA DOKUMENTACIJOM

### Iz vaše dokumentacije (https://www.posta.rs/wsp-help/uvod/uvod.aspx):
- Username: TEST
- Password: t3st
- Test URL: http://212.62.32.201/WspWebApi/transakcija

### Šta mi koristimo:
-  Username: TEST (tačno poklapanje)
-  Password: t3st (tačno poklapanje)
-  URL: http://212.62.32.201/WspWebApi/transakcija (tačno poklapanje)

##  PITANJA

1. **Da li su pristupni podaci TEST/t3st aktivni u test okruženju?**
  - Možda su deaktivirani?
  - Da li je potrebna prethodna registracija?

2. **Postoje li dodatni zahtevi za aktivaciju?**
  - IP whitelist?
  - Prethodna registracija preko emaila?
  - Specijalni header-i u zahtevu?

3. **Možete li obezbediti radni primer zahteva?**
  - Sa aktivnim test pristupnim podacima
  - Koji uspešno kreira pošiljku

---

**P.S.** Sigurni smo da drugim klijentima integracija radi uspešno. Molimo vas da nam pomognete da pronađemo šta radimo pogrešno. Možda postoji neki neočigledan korak koji propuštamo?

## 🔗 PRILOZI

### Kompletan test skript (Go)
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

### Komanda za brzi test
```bash
curl -X POST http://212.62.32.201/WspWebApi/transakcija \
  -H "Content-Type: application/json" \
  -d '{"StrKlijent":"{\"Username\":\"TEST\",\"Password\":\"t3st\",\"Jezik\":\"LAT\",\"IdTipUredjaja\":2}","Servis":3,"IdVrstaTranskacije":3,"TipSerijalizacije":2,"IdTransakcija":"test-123","StrIn":"{\"Naziv\":\"Novi\"}"}'
```
