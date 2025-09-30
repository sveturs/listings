# Tehničke specifikacije za finansijske operatore marketplace-a Sve Tu

## 1. Opšte informacije o projektu

### 1.1 O kompaniji
- **Naziv**: Još nije određen.

### 1.2 Obim i objekti
- **Trenutni obim**: Startup, planirani lansirani Q4 2025
- **Prognoza za godinu 1**:
  - 5.000 aktivnih korisnika
  - 500 transakcija/mesečno
  - Prosečan račun: 5.000-15.000 RSD
- **Prognoza za godinu 3**:
  - 50.000 aktivnih korisnika
  - 10.000 transakcija/mesečno
  - Prosečan račun: 10.000-25.000 RSD

## 2. Poslovni model i finansijski tokovi

### 2.1 Izvori prihoda platforme

#### A. Provizije od prodaje (glavni prihod)
- **C2C transakcije**: 2-5% od iznosa transakcije
- **B2C prodavnice**: 5-7% od iznosa prodaje

#### B. Plaćene usluge za korisnike
- **Promocija oglasa**:
  - Podizanje u vrh: 200-500 RSD
  - VIP status: 1.000-3.000 RSD/mesečno
  - Označavanje bojom: 300-700 RSD
- **Pretplate za prodavce (prodavnice)**:
  - Starter: 1.500 RSD/mesečno
  - Professional: 5.000 RSD/mesečno
  - Enterprise: 15.000 RSD/mesečno
- **Dodatne usluge**:
  - Premium podrška
  - Analitika i statistika

### 2.2 Tipovi finansijskih operacija

#### Scenario 1: Plaćanje internih usluga marketplace-a
**Tok**: Kupac → Platežni gateway → Marketplace

**Primeri operacija**:
- Plaćanje pretplate na prodavnicu
- Kupovina usluga promocije
- Plaćanje VIP statusa
- Dopunjavanje internog balansa za buduća plaćanja usluga marketplace-a

**Zahtevi**:
- Trenutno pripis sredstava
- Bez holding-a
- Automatska fiskalizacija
- Podrška za rekurentna plaćanja za pretplate (mogućnost automatskog naplaćivanja novca sa kartice klijenta na redovnoj osnovi)

#### Scenario 2: B2C prodaja (prodavnice)
**Tok**: Kupac → Ekvajring → Eskrou → Split plaćanja

**Raspodela sredstava**:
- Prodavac: 85-90% (minus provizija marketplace-a)
- Marketplace: 3-5%
- Logistika: 0.01-10% (ako je primenjivo)
- PVZ: 0.1-2% (ako je primenjivo)
- Platežni sistem: 0.01-2%

**Zahtevi**:
- Zadržavanje sredstava na 7-14 dana (zaštita kupca)
- Automatski split nakon potvrde dostave
- Mogućnost parcijalnog povrata
- Integracija sa logističkim partnerima

#### Scenario 3: C2C prodaja (fizička lica)
**Tok**: Kupac → Ekvajring → Eskrou → Split plaćanja

**Raspodela sredstava**:
- Prodavac: 85-95% (minus provizija marketplace-a)
- Marketplace: 2-5%
- Logistika: 0.01-10% (ako je prodavac izabrao dostavu)
- PVZ: 0.1-2% (ako je prodavac predao robu u PVZ)
- Platežni sistem: 0.01-2%

**Karakteristike**:
- Prodavac - fizičko lice bez PIB-a
- Zadržavanje sredstava do potvrde prijema
- Isplata na bankarski račun ili karticu prodavca
- Mogućnost korišćenja PVZ-a i dostave

**Zahtevi**:
- KYC verifikacija za prodavce (iznosi > 10.000 RSD)
- Podrška za isplate na kartice fizičkih lica

## 3. Tehnički zahtevi

### 3.1 Načini plaćanja (prijem sredstava)
**Potrebni na početku**:
- Bankovne kartice (Visa, Mastercard, Dina)
- IPS QR (trenutna plaćanja preko QR koda - standard NBS Srbije)
- Bankarska doznaka

### 3.2 Načini isplata (povlačenje sredstava)
**Potrebni na početku**:
- Bankarska doznaka na račun pravnog lica (za B2C prodavce)
- Transfer na karticu fizičkog lica (za C2C prodavce)
- IPS transferi (trenutni transferi između banaka Srbije)

### 3.3 Funkcionalne mogućnosti

#### Obavezne funkcije:
1. **Eskrou (zadržavanje sredstava)**
   - Podesivi period zadržavanja (1-30 dana)
   - Automatsko oslobađanje po ispunjenim uslovima
   - Mogućnost produženja hold-a pri sporovima

2. **Split plaćanja (multisplit)**
   - Podela na 3+ primaoca
   - Podesiva pravila podele
   - Automatski proračun provizija

3. **Povrati i otkazivanja**
   - Potpuni povrat
   - Parcijalni povrat
   - Otkazivanje pre zarobljivanja sredstava

4. **Rekurentna plaćanja**
   - Za pretplate
   - Automatska obnova
   - Upravljanje pretplatama

5. **Marketplace funkcije**
   - Onboarding prodavaca (KYC/KYB)
   - Upravljanje balansima
   - Masovne isplate
   - Registri transakcija

#### Poželjne funkcije:
**Tokenizacija kartica** - čuvanje kartica kupaca u šifrovanom obliku (tokeni) za brza ponovna plaćanja bez ponovnog unosa podataka kartice

**3D Secure 2.0** - dodatna provera autentičnosti plaćanja preko SMS koda ili aplikacije banke, smanjuje rizike od prevara i prebacuje odgovornost na banku-izdavaoca

**Antifraud sistem** - automatska provera plaćanja na znakove prevare (neobična geolokacija, višestruki pokušaji, sumnjivi obrasci)

**Dinamička konverzija valuta** - mogućnost prikazivanja cena i primanja plaćanja u valuti kartice kupca sa automatskom konverzijom

**Mobilni SDK** - gotove biblioteke za integraciju plaćanja u mobilne aplikacije (iOS/Android)

**Webhook obaveštenja** - automatski HTTP zahtevi od platežnog sistema prema našem serveru pri promeni statusa plaćanja

### 3.4 Integracija

#### API zahtevi:
- REST API (poželjno) ili drugi savremeni protokol
- Webhook za statuse plaćanja
- Sandbox (test okruženje) za razvoj
- API dokumentacija (engleski ili srpski)
- Primeri integracije

**Napomena**: Spremni smo da se prilagodimo API-ju platežnog sistema. Naš stek: Go (backend), JavaScript/TypeScript (frontend)

#### Bezbednost:
**PCI DSS compliance** - sertifikat usklađenosti sa standardima bezbednosti platnih kartica. Znači da platežni sistem štiti podatke kartica po međunarodnim standardima

**TLS 1.2+ za sve konekcije** - korišćenje savremenih protokola šifrovanja za sve konekcije između našeg servera i platežnog sistema

**Tokenizacija osetljivih podataka** - zamena pravih podataka kartica jedinstvenim tokenima, koji su beskorisni za zlikovce

**Dvofaktorska autentifikacija za administratore** - dodatna zaštita administrativnih panela preko SMS-a ili autentifikator aplikacije

**IP whitelist** - ograničavanje pristupa API-ju samo sa dozvoljenih IP adresa naših servera

## 4. Pravni i compliance zahtevi

### 4.1 Licenciranje (zahtevi za platežnog operatora)

**Licenca NBS za platežne usluge** - zvanično odobrenje od Narodne banke Srbije za obavljanje platežnih operacija. Potvrđuje legalnost i pouzdanost operatora

**Pravo na obavljanje ekvajaringa** - licenca za prijem plaćanja sa bankovnih kartica u ime trgovaca (u našem slučaju - marketplace-a)

**Pravo na eskrou usluge** - dozvola za privremeno čuvanje sredstava kupca do ispunjenja uslova transakcije (kritično za zaštitu kupaca)

**Pravo na isplate fizičkim licima** - licenca za obavljanje isplata na kartice i račune fizičkih lica, koja nisu preduzetnici (neophodno za C2C)

### 4.2 Usklađenost sa regulatornim zahtevima

**Usklađenost sa zakonom Srbije o platežnim uslugama** - rad u okviru zakona "Zakon o platnim uslugama", koji reguliše platežne operacije u Srbiji

**AML/CFT procedure** - Anti-Money Laundering (borba protiv pranja novca) i Counter Financing of Terrorism (borba protiv finansiranja terorizma). Sistem mera za sprečavanje korišćenja platforme u nezakonitе svrhe

**KYC/KYB verifikacija** - Know Your Customer (poznaj svog klijenta) za fizička lica i Know Your Business za pravna lica. Provera identiteta korisnika za sprečavanje prevara

**Fiskalizacija transakcija** - automatska predaja podataka o transakcijama poreskoj službi Srbije u skladu sa zahtevima zakona o fiskalizaciji

**GDPR compliance** - usklađenost sa evropskim pravilnikom o zaštiti ličnih podataka (General Data Protection Regulation), obaveznim u Srbiji

### 4.3 Izveštavanje
- Mesečni finansijski izveštaji
- Registri transakcija za poresku službu
- API za izvoz podataka
- Podrška za reviziju

## 5. Komercijalni uslovi

### 5.1 Uslovi saradnje
- Odsustvo zahteva po minimalnom prometu na početku
- Transparentna tarifa
- Tehnička podrška 24/7 za kritične incidente
- Dedicirani menadžer
- Pomoć u prolaženju sertifikacije kada je potrebna

## 6. Tehnološki stek platforme

### Backend
- **Jezik**: Go 1.21+
- **Baza podataka**: PostgreSQL 15
- **Keš**: Redis
- **Pretraga**: OpenSearch
- **Redovi**: RabbitMQ (planirano)

### Frontend
- **Framework**: Next.js 15 (React 19)
- **Jezik**: TypeScript
- **State**: Redux Toolkit

### Infrastruktura
- **Hosting**: VPS
- **CDN**: Cloudflare
- **Monitoring**: Prometheus + Grafana
- **CI/CD**: GitHub Actions

## 7. Mapa puta integracije

### Faza 1: MVP
- Osnovni ekvajring (kartice)
- Jednostavna direktna plaćanja
- Dopunjavanje balansa
- Osnovne isplate prodavcima

### Faza 2: Eskrou
- Zadržavanje sredstava
- Upravljanje sporovima
- Automatsko oslobađanje
- Parcijalni povrati

### Faza 3: Marketplace
- Split plaćanja
- Masovne isplate
- KYC/KYB onboarding
- Rekurentna plaćanja

## 8. Kriterijumi izbora partnera

### Obavezni zahtevi
1. ✅ Licenca NBS na teritoriji Srbije
2. ✅ Iskustvo rada sa marketplace-ima
3. ✅ Podrška za eskrou i split plaćanja (eskrou dopušten sredstvima VISA i Mastercard)
4. ✅ API dokumentacija i sandbox

### Prednosti
- Gotova rešenja za marketplace-ove
- Konkurentne tarife

## 9. Kontakt informacije

**Sajt**: https://dev.svetu.rs

## 10. Prilozi

### Prilog A: Primeri scenarija korišćenja

#### Scenario 1: Kupovina proizvoda u prodavnici
1. Kupac bira proizvod za 10.000 RSD
2. Plaća karticom preko platežnog gateway-a
3. Sredstva se drže 7 dana
4. Nakon potvrde dostave:
   - Prodavac prima: 9.500 RSD
   - Marketplace: 300 RSD (3%)
   - Platežni sistem: 200 RSD (2%)

#### Scenario 2: C2C transakcija sa dostavom
1. Kupac kupuje polovan telefon za 30.000 RSD
2. Plaćanje preko ekvajaringa
3. Zadržavanje na 14 dana
4. Podela nakon prijema:
   - Prodavac: 27.900 RSD
   - Marketplace: 1.500 RSD (5%)
   - Dostava: 600 RSD (2%)

#### Scenario 3: Pretplata na prodavnicu
1. Prodavac oformi pretplatu Professional
2. Mesečno plaćanje 5.000 RSD
3. Automatska obnova
4. Celo plaćanje ide marketplace-u