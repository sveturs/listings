# 🏗️ TEHNIČKA ARHITEKTURA PLATFORME SVE TU

## 📋 Sadržaj
1. [Opšta arhitektura sistema](#opšta-arhitektura)
2. [Komponente platforme](#komponente)
3. [Integracija sa logistikom](#logistika)
4. [Obrada plaćanja](#plaćanja)
5. [Bezbednost i skalabilnost](#bezbednost)
6. [Tehnički zahtevi za PostExpress API](#zahtevi-api)

---

## 🎯 Opšta arhitektura sistema {#opšta-arhitektura}

### Arhitektura visokog nivoa

```
┌─────────────────────────────────────────────────────────┐
│                    FRONTEND SLOJ                        │
├──────────────────┬──────────────────┬──────────────────┤
│   Web App        │   Mobile PWA     │   Admin Panel    │
│ Next.js 15       │   Responsive     │   Dashboard      │
│ React 19         │   Design         │                  │
└──────────────────┴──────────────────┴──────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                    API GATEWAY                          │
├──────────────────────┬──────────────────────────────────┤
│    Nginx             │    Auth Service                  │
│    Load Balancer     │    JWT + OAuth                   │
└──────────────────────┴──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                  BACKEND SERVISI                        │
├─────────┬─────────┬─────────┬─────────┬────────────────┤
│Main API │Market-  │Order    │Payment  │Logistics       │
│Go+Fiber │place    │Mgmt     │Service  │Adapter         │
└─────────┴─────────┴─────────┴─────────┴────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                    DATA SLOJ                            │
├──────────┬──────────┬──────────┬───────────────────────┤
│PostgreSQL│Redis     │OpenSearch│MinIO                  │
│Main DB   │Cache     │Full-text │S3 Storage             │
└──────────┴──────────┴──────────┴───────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────┐
│                 EKSTERNI SERVISI                        │
├──────────┬──────────────────┬───────────────────────────┤
│PostExpress│Payment Gateway  │SMS Provider               │
│API        │                  │                          │
└──────────┴──────────────────┴───────────────────────────┘
```

### Ključne karakteristike

| Parametar | Vrednost | Obrazloženje |
|-----------|----------|--------------|
| **Arhitektura** | Mikroservisna ready | Skalabilnost |
| **Backend jezik** | Go 1.21+ | Performanse, konkurentnost |
| **Baza podataka** | PostgreSQL 15 | ACID, pouzdanost |
| **Keširanje** | Redis Cluster | Brzina, pub/sub |
| **Pretraživanje** | OpenSearch | Full-text pretraga |
| **Fajlovi** | MinIO (S3) | Skalabilno skladištenje |
| **API** | REST + WebSocket | Standard + real-time |

---

## 🔧 Komponente platforme {#komponente}

### 1. Frontend aplikacije

```typescript
// Struktura frontend aplikacije
frontend/
├── svetu/                      // Glavna aplikacija
│   ├── src/
│   │   ├── app/               // Next.js App Router
│   │   │   ├── [locale]/      // i18n (ru, en, sr)
│   │   │   │   ├── marketplace/
│   │   │   │   ├── checkout/
│   │   │   │   ├── orders/
│   │   │   │   └── tracking/  // Praćenje dostave
│   │   ├── components/
│   │   │   ├── delivery/      // Komponente dostave
│   │   │   ├── payment/       // Plaćanje i pouzeće
│   │   │   └── tracking/      // Tracking porudžbina
│   │   └── services/
│   │       ├── api/           // API klijent
│   │       └── logistics/     // Logistika klijent
```

### 2. Backend servisi

```go
// Struktura backend-a
backend/
├── cmd/
│   └── api/                   // Ulazna tačka
├── internal/
│   ├── proj/                  // Poslovna logika
│   │   ├── marketplace/       // Marketplace
│   │   ├── orders/           // Upravljanje porudžbinama
│   │   ├── payments/         // Plaćanja i split
│   │   ├── logistics/        // Integracija sa PostExpress
│   │   │   ├── postexpress/ // PostExpress adapter
│   │   │   ├── routing/     // Rutiranje
│   │   │   └── tracking/    // Praćenje
│   │   └── notifications/   // Obaveštenja
│   └── storage/
│       ├── postgres/         // PostgreSQL repozitorijumi
│       └── redis/           // Redis keš
```

### 3. Baza podataka

```sql
-- Glavne tabele za logistiku
CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    order_number VARCHAR(32) UNIQUE,
    user_id INTEGER NOT NULL,
    seller_id INTEGER NOT NULL,
    
    -- Status
    status VARCHAR(30) NOT NULL,
    payment_method VARCHAR(20), -- 'online', 'pouzece'
    
    -- Adrese
    pickup_address JSONB,      -- Adresa prodavca
    delivery_address JSONB,     -- Adresa kupca
    
    -- Logistika
    logistics_provider VARCHAR(30) DEFAULT 'postexpress',
    tracking_number VARCHAR(100),
    estimated_delivery DATE,
    
    -- Split plaćanja
    total_amount DECIMAL(10,2),
    seller_amount DECIMAL(10,2),
    platform_fee DECIMAL(10,2),
    delivery_fee DECIMAL(10,2),
    
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE delivery_tracking (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT REFERENCES orders(id),
    tracking_number VARCHAR(100),
    
    -- Događaji od PostExpress-a
    status VARCHAR(50),
    status_description TEXT,
    location VARCHAR(255),
    
    -- Metapodaci
    event_time TIMESTAMPTZ,
    raw_event JSONB,
    
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE pouzece_transactions (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT REFERENCES orders(id),
    
    -- Sume za split
    collected_amount DECIMAL(10,2),    -- Prikupljeno od kurira
    seller_payout DECIMAL(10,2),       -- Za isplatu prodavcu
    platform_commission DECIMAL(10,2), -- Naša provizija
    delivery_cost DECIMAL(10,2),       -- PostExpress
    
    -- Status isplata
    seller_paid BOOLEAN DEFAULT FALSE,
    seller_paid_at TIMESTAMPTZ,
    commission_received BOOLEAN DEFAULT FALSE,
    
    postexpress_reference VARCHAR(100),
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

---

## 🚚 Integracija sa logistikom {#logistika}

### Tok integracije sa PostExpress

```
PROCES KREIRANJA DOSTAVE:
========================

1. Kupac kreira porudžbinu
   └─> Marketplace prima zahtev
       └─> Logistics Service poziva PostExpress API
           └─> PostExpress vraća tracking broj

2. PostExpress preuzima paket
   └─> Webhook: "picked_up"
       └─> Ažuriramo status u bazi
           └─> Šaljemo notifikaciju kupcu

3. Dostava u toku
   └─> Webhook događaji o statusu
       └─> Real-time ažuriranje
           └─> Push notifikacije

4. Isporuka kupcu
   ├─> Ako je POUZEĆE:
   │   └─> Kurir naplaćuje
   │       └─> PostExpress šalje potvrdu
   │           └─> Mi procesiramo split plaćanja
   └─> Ako je ONLINE plaćanje:
       └─> Samo potvrda dostave
```

### API integracija

```go
// logistics/postexpress/client.go
package postexpress

import (
    "context"
    "encoding/json"
    "fmt"
)

type PostExpressClient struct {
    baseURL    string
    apiKey     string
    httpClient *http.Client
}

// CreateShipment kreira zahtev za dostavu
func (c *PostExpressClient) CreateShipment(ctx context.Context, req ShipmentRequest) (*ShipmentResponse, error) {
    payload := map[string]interface{}{
        "posiljalac": map[string]interface{}{
            "ime":    req.Prodavac.Ime,
            "adresa": req.Prodavac.Adresa,
            "telefon": req.Prodavac.Telefon,
            "email":   req.Prodavac.Email,
        },
        "primalac": map[string]interface{}{
            "ime":    req.Kupac.Ime,
            "adresa": req.Kupac.Adresa,
            "telefon": req.Kupac.Telefon,
            "email":   req.Kupac.Email,
        },
        "paket": map[string]interface{}{
            "tezina":      req.Paket.Tezina,
            "dimenzije":   req.Paket.Dimenzije,
            "vrednost":    req.Paket.Vrednost,
            "opis":        req.Paket.Opis,
        },
        "usluge": map[string]interface{}{
            "pouzece":         req.JePouzece,
            "iznos_pouzeca":   req.IznosPouzeca,
            "osiguranje":      req.Osiguranje,
            "sms_obavestenje": true,
        },
        "referenca": req.BrojPorudzbine,
    }
    
    // Slanje zahteva ka PostExpress
    resp, err := c.makeRequest(ctx, "POST", "/api/v1/posiljke", payload)
    if err != nil {
        return nil, fmt.Errorf("greska pri kreiranju posiljke: %w", err)
    }
    
    var result ShipmentResponse
    if err := json.Unmarshal(resp, &result); err != nil {
        return nil, err
    }
    
    return &result, nil
}

// TrackShipment pracenje posiljke
func (c *PostExpressClient) TrackShipment(ctx context.Context, trackingNumber string) (*TrackingInfo, error) {
    endpoint := fmt.Sprintf("/api/v1/pracenje/%s", trackingNumber)
    
    resp, err := c.makeRequest(ctx, "GET", endpoint, nil)
    if err != nil {
        return nil, fmt.Errorf("greska pri pracenju posiljke: %w", err)
    }
    
    var tracking TrackingInfo
    if err := json.Unmarshal(resp, &tracking); err != nil {
        return nil, err
    }
    
    return &tracking, nil
}

// HandleWebhook obrada webhook-a od PostExpress-a
func (c *PostExpressClient) HandleWebhook(payload []byte) (*WebhookEvent, error) {
    var event WebhookEvent
    if err := json.Unmarshal(payload, &event); err != nil {
        return nil, err
    }
    
    // Verifikacija potpisa
    if !c.verifyWebhookSignature(payload, event.Potpis) {
        return nil, fmt.Errorf("nevaljan webhook potpis")
    }
    
    return &event, nil
}
```

### Obrada različitih scenarija

```go
// logistics/service.go
package logistics

type LogisticsService struct {
    postExpress *postexpress.PostExpressClient
    orderRepo   OrderRepository
    eventBus    EventBus
}

// ProcessC2COrder obrada C2C porudžbine
func (s *LogisticsService) ProcessC2COrder(ctx context.Context, order Order) error {
    // Proveravamo grad prodavca
    if s.jeBeogradIliNoviSad(order.Prodavac.Grad) {
        // Za Beograd i Novi Sad - drop-off u pošti
        return s.createDropOffShipment(ctx, order)
    }
    
    // Za ostale gradove - preuzimanje na adresi
    return s.createPickupShipment(ctx, order)
}

// ProcessB2COrder obrada B2C porudžbine
func (s *LogisticsService) ProcessB2COrder(ctx context.Context, order Order) error {
    // B2C uvek sa preuzimanjem od biznisa
    shipment, err := s.postExpress.CreateShipment(ctx, ShipmentRequest{
        Prodavac:       order.Prodavac,
        Kupac:          order.Kupac,
        Paket:          order.Paket,
        JePouzece:      order.NacinPlacanja == "pouzece",
        IznosPouzeca:   order.UkupanIznos,
        BrojPorudzbine: order.Broj,
    })
    
    if err != nil {
        return fmt.Errorf("greska pri kreiranju B2C posiljke: %w", err)
    }
    
    // Čuvamo tracking broj
    order.TrackingBroj = shipment.TrackingBroj
    order.ProcenjenaIsporuka = shipment.ProcenjenaIsporuka
    
    return s.orderRepo.UpdateOrder(ctx, order)
}

// HandleDeliveryWebhook obrada webhook-a o dostavi
func (s *LogisticsService) HandleDeliveryWebhook(ctx context.Context, event WebhookEvent) error {
    order, err := s.orderRepo.GetByTrackingNumber(ctx, event.TrackingBroj)
    if err != nil {
        return err
    }
    
    switch event.Status {
    case "preuzeto":
        order.Status = "u_transportu"
        s.eventBus.Publish(OrderPickedUpEvent{OrderID: order.ID})
        
    case "isporuceno":
        order.Status = "isporuceno"
        s.eventBus.Publish(OrderDeliveredEvent{OrderID: order.ID})
        
        // Ako je pouzeće - obrađujemo split
        if order.NacinPlacanja == "pouzece" {
            s.processCODPayment(ctx, order, event.PouzeceNaplaceno)
        }
        
    case "vraceno":
        order.Status = "vraceno"
        s.eventBus.Publish(OrderReturnedEvent{OrderID: order.ID})
    }
    
    return s.orderRepo.UpdateOrder(ctx, order)
}
```

---

## 💰 Obrada plaćanja {#plaćanja}

### Šema Split plaćanja

```
PROBLEM SA POUZEĆE PLAĆANJEM:
==============================

Trenutno stanje PostExpress:
┌──────────┐      ┌─────────┐      ┌──────────────┐
│  Kupac   │ ---> │ Kurir   │ ---> │ PostExpress  │
│          │ RSD  │         │ 100% │    Račun     │
└──────────┘      └─────────┘      └──────────────┘
                                           │
                                           ▼
                                    ┌──────────────┐
                                    │ Naš račun    │
                                    │   (100%)     │
                                    └──────────────┘

❌ PROBLEM: PostExpress NE podržava automatski split!

NAŠE REŠENJE:
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│ Escrow račun │ --> │ Dnevna       │ --> │ Automatske   │
│              │     │ reconcilacija│     │ isplate      │
└──────────────┘     └──────────────┘     └──────────────┘
        │                                          │
        ▼                                          ▼
┌──────────────┐                           ┌──────────────┐
│ Prodavac 95% │                           │ Naša provizija│
└──────────────┘                           │      5%       │
```

### Zaobilazno rešenje za Split plaćanja

```go
// payments/split_service.go
package payments

type SplitPaymentService struct {
    db           *pgxpool.Pool
    bankAPI      BankingAPI
    postExpress  PostExpressClient
}

// ProcessCODSettlement obrada pouzeće plaćanja
func (s *SplitPaymentService) ProcessCODSettlement(ctx context.Context, order Order) error {
    // 1. Dobijamo potvrdu od PostExpress-a
    potvrda, err := s.postExpress.GetCODConfirmation(ctx, order.TrackingBroj)
    if err != nil {
        return err
    }
    
    // 2. Kreiramo zapis o transakciji
    transakcija := PouzeceTransakcija{
        OrderID:           order.ID,
        NaplacenIznos:     potvrda.Iznos,
        IsplataProdavcu:   order.IznosProdavca,    // 95% od sume
        PlatformaProvizija: order.PlatformaProvizija, // 5% provizija
        TrosakDostave:     order.TrosakDostave,
        PostExpressRef:    potvrda.Referenca,
    }
    
    // 3. Čuvamo transakciju
    if err := s.saveTransaction(ctx, transakcija); err != nil {
        return err
    }
    
    // 4. Čekamo prijem sredstava od PostExpress-a
    // (obično isti dan prema njihovim uslovima)
    go s.scheduleSettlement(ctx, transakcija)
    
    return nil
}

// scheduleSettlement planira isplatu prodavcu
func (s *SplitPaymentService) scheduleSettlement(ctx context.Context, tx PouzeceTransakcija) {
    // Proveravamo prijem sredstava od PostExpress-a
    ticker := time.NewTicker(1 * time.Hour)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            // Proveravamo bankovni račun
            if s.checkFundsReceived(ctx, tx.PostExpressRef) {
                // Sredstva primljena - isplaćujemo prodavcu
                if err := s.payoutToSeller(ctx, tx); err != nil {
                    log.Printf("Greška pri isplati: %v", err)
                    continue
                }
                
                // Označavamo kao isplaćeno
                s.markAsPaid(ctx, tx.ID)
                return
            }
        case <-time.After(48 * time.Hour):
            // Eskalacija ako nismo primili sredstva
            s.escalateDelayedPayment(ctx, tx)
            return
        }
    }
}

// DailyReconciliation dnevna reconcilacija
func (s *SplitPaymentService) DailyReconciliation(ctx context.Context) error {
    // 1. Dobijamo listu COD transakcija od PostExpress-a
    transakcije, err := s.postExpress.GetDailyCODReport(ctx, time.Now())
    if err != nil {
        return err
    }
    
    // 2. Upoređujemo sa našim zapisima
    for _, peTx := range transakcije {
        nasaTx, err := s.getTransactionByRef(ctx, peTx.Referenca)
        if err != nil {
            log.Printf("Nedostaje transakcija: %s", peTx.Referenca)
            continue
        }
        
        // 3. Proveravamo sume
        if nasaTx.NaplacenIznos != peTx.Iznos {
            s.flagDiscrepancy(ctx, nasaTx.ID, "neslaganje_iznosa")
        }
        
        // 4. Ako je sve OK i sredstva primljena - isplaćujemo
        if !nasaTx.ProdavacIsplacen && s.checkFundsReceived(ctx, peTx.Referenca) {
            s.payoutToSeller(ctx, nasaTx)
        }
    }
    
    return nil
}
```

---

## 🔒 Bezbednost i skalabilnost {#bezbednost}

### Bezbednosne mere

```yaml
Bezbednosne mere:
  Autentifikacija:
    - JWT tokeni sa kratkim TTL (15 min)
    - Refresh tokeni u httpOnly kolačićima
    - OAuth 2.0 za eksterne integracije
    
  Zaštita podataka:
    - TLS 1.3 za sve konekcije
    - Šifrovanje osetljivih podataka (AES-256)
    - PCI DSS compliance za plaćanja
    
  API bezbednost:
    - Rate limiting (100 req/min po IP)
    - API ključevi za B2B partnere
    - Verifikacija webhook potpisa
    
  Infrastruktura:
    - WAF (Cloudflare)
    - DDoS zaštita
    - Redovne bezbednosne provere
```

### Skalabilnost

```
ARHITEKTURA ZA SKALIRANJE:
==========================

        ┌─────────────────┐
        │  Load Balancer  │
        │ HAProxy/Nginx   │
        └────────┬────────┘
                 │
    ┌────────────┼────────────┐
    ▼            ▼            ▼
┌─────────┐ ┌─────────┐ ┌─────────┐
│  App    │ │  App    │ │  App    │
│ Server 1│ │ Server 2│ │ Server N│
└─────────┘ └─────────┘ └─────────┘
    │            │            │
    └────────────┼────────────┘
                 │
    ┌────────────┼────────────┐
    ▼            ▼            ▼
┌─────────┐ ┌─────────┐ ┌─────────┐
│PostgreSQL│ │  Read   │ │  Read   │
│  Master │ │Replica 1│ │Replica 2│
└─────────┘ └─────────┘ └─────────┘
                 │
    ┌────────────┼────────────┐
    ▼            ▼            ▼
┌─────────┐ ┌─────────┐ ┌─────────┐
│  Redis  │ │  Redis  │ │  Redis  │
│  Node 1 │ │  Node 2 │ │  Node 3 │
└─────────┘ └─────────┘ └─────────┘
```

### Metrike performansi

| Metrika | Trenutno | Ciljano | Maksimum |
|---------|----------|---------|----------|
| **RPS (zahteva/sek)** | 500 | 2.000 | 10.000 |
| **Latencija p99** | 200ms | 100ms | 500ms |
| **Konkurentnih korisnika** | 1.000 | 5.000 | 20.000 |
| **Porudžbina/dan** | 200 | 2.000 | 10.000 |
| **Veličina baze** | 300MB | 10GB | 100GB |
| **Uptime** | 99,5% | 99,9% | 99,99% |

---

## 📝 Tehnički zahtevi za PostExpress API {#zahtevi-api}

### Minimalni skup metoda

```yaml
Potrebni API endpoint-i:
  Pošiljke:
    - POST /api/v1/posiljke         # Kreiranje pošiljke
    - GET /api/v1/posiljke/{id}     # Informacije o pošiljci
    - PUT /api/v1/posiljke/{id}     # Ažuriranje (pre preuzimanja)
    - DELETE /api/v1/posiljke/{id}  # Otkazivanje
    
  Praćenje:
    - GET /api/v1/pracenje/{broj}   # Status dostave
    - GET /api/v1/pracenje/istorija # Istorija događaja
    
  Pouzeće:
    - GET /api/v1/pouzece/transakcije # Lista pouzeće transakcija
    - GET /api/v1/pouzece/obracun     # Status isplata
    - POST /api/v1/pouzece/uskladjivanje # Reconcilacija
    
  Webhook-ovi:
    - POST /nas-endpoint/webhook    # Događaji od PostExpress-a
    
  Izveštaji:
    - GET /api/v1/izvestaji/dnevni  # Dnevni izveštaj
    - GET /api/v1/izvestaji/pouzece # Pouzeće izveštaj
```

### Webhook događaji

```json
{
  "tip_dogadjaja": "posiljka.isporucena",
  "vreme": "2025-01-08T14:30:00Z",
  "podaci": {
    "broj_pracenja": "PE123456789",
    "referenca": "ORD-2025-001234",
    "status": "isporucena",
    "vreme_isporuke": "2025-01-08T14:28:00Z",
    "primalac": "Ime Prezime",
    "pouzece": {
      "naplaceno": true,
      "iznos": 5000.00,
      "valuta": "RSD",
      "bice_uplaceno": "2025-01-08T23:59:59Z"
    },
    "potpis": "sha256=abc123..."
  }
}
```

### Zahtevi za integraciju

1. **Test okruženje**
   - Sandbox API za razvoj
   - Test tracking brojevi
   - Simulacija webhook događaja

2. **Dokumentacija**
   - OpenAPI/Swagger specifikacija
   - Primeri zahteva/odgovora
   - Kodovi grešaka i njihovi opisi

3. **SLA**
   - Dostupnost API-ja: 99,9%
   - Vreme odgovora: <500ms
   - Rate limits: minimum 1000 req/min

4. **Podrška**
   - Tehnička podrška pri integraciji
   - Posvećeni menadžer
   - Kanal za eskalaciju problema

---

## 📊 Projektovani obimi

### Rast platforme

| Period | Korisnici | Porudžbina/mes | B2C | C2C | PostExpress |
|--------|-----------|----------------|-----|-----|-------------|
| Lansiranje (Sept 2025) | 5.000 | 200 | 140 | 60 | 200 |
| 3 meseca | 20.000 | 800 | 560 | 240 | 800 |
| 6 meseci | 50.000 | 2.000 | 1.400 | 600 | 2.000 |
| 1 godina | 100.000 | 5.000 | 3.500 | 1.500 | 5.000 |
| 2 godine | 250.000 | 15.000 | 10.500 | 4.500 | 15.000 |

### Geografska pokrivenost

```
Raspodela porudžbina po gradovima:
===================================
Beograd     ████████████████████ 40%
Novi Sad    ██████████ 20%
Niš         █████ 10%
Kragujevac  ████ 8%
Subotica    ██ 5%
Ostali      ████████ 17%
```

---

## ✅ Spremnost za integraciju

### Šta imamo spremno

- ✅ Backend arhitektura sa podrškom za više provajdera
- ✅ Sistem upravljanja porudžbinama i statusima
- ✅ Mehanizam obrade webhook događaja
- ✅ UI/UX za izbor dostave i praćenje
- ✅ Sistem obaveštenja (email, SMS, push)

### Šta nam je potrebno od PostExpress-a

- 📋 API dokumentacija i pristupi za sandbox
- 🔄 Webhook endpoint-i i format događaja
- 💰 Rešenje za split plaćanja kod pouzeća
- 📦 Proces onboarding-a za naše prodavce
- 📊 Tarife i uslovi za marketplace

### Plan integracije

```
VREMENSKI PLAN INTEGRACIJE:
===========================

Priprema (15-17 januar):
├─ Dobijanje dokumentacije ─────── 3 dana
└─ Analiza API-ja ────────────────── 2 dana

Razvoj (18-29 januar):
├─ Osnovna integracija ──────────── 5 dana
├─ Testiranje sandbox ───────────── 3 dana
└─ Obrada webhook-ova ──────────── 2 dana

Pilot (30 januar - 7 februar):
├─ Pilot sa 10 prodavaca ─────────── 7 dana
└─ Ispravka problema ───────────── 3 dana

Lansiranje (8-15 februar):
├─ Production priprema ──────────── 2 dana
└─ Puno lansiranje ────────────────── 15. feb
```

---

## 📞 Kontakti tehničkog tima

**CTO / Tehnički direktor**
- Dmitrii Voroshilov
- Email: tech@svetu.rs
- Telegram: @dmitrii_tech

**Spremni smo za:**
- Tehničke konsultacije
- Zajedničku izradu rešenja
- Pilot projekte
- Dugoročno partnerstvo

---

## 📌 Napomene o dijagramima

*Dijagrami u ovom dokumentu su prikazani u ASCII art formatu za maksimalnu kompatibilnost sa svim editorima teksta. Originalni Mermaid dijagrami su zamenjeni ASCII reprezentacijom koja se ispravno prikazuje u svim text viewer-ima.*

---

*Ovaj dokument je pripremljen za sastanak sa PostExpress-om i sadrži tehničke informacije o platformi SVE TU. Otvoreni smo za prilagođavanje naše arhitekture zahtevima i mogućnostima PostExpress API-ja.*