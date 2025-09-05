# Zahtev za kreiranje korisničkog profila za integraciju informatičkog sistema korisnika u D Express sistem podataka

## A - Datum podnošenja zahteva
**05.09.2025.**

---

## B - Razvoj funkcionalnosti vrši
(Front end za unos podataka, back end za obradu naloga i štampu nalepnica, sistem za slanje naloga ka D Express-u)

☑️ **kompanija za koju se radi integracija** - za slučaj da kompanija sama radi razvoj, preskočiti ceo blok D  
☐ spoljni saradnik (integrator) - u slučaju da kompanija ne koristi ili nema svoj razvoj, preskočiti blok C5

---

## C - Podaci kompanije za koju se vrši integracija

### C1 - Naziv kompanije
**SVE TU PLATFORMA D.O.O.**

### C2 - UK broj kompanije (ako nema ostaviti prazno)
*[ostaviti prazno ili dodati matični broj ako postoji]*

### C3 - Adresa sedišta kompanije
- **Naziv ulice:** Vase Stajića
- **Kućni broj:** 18
- **PTT broj:** 21101
- **Naziv mesta:** Novi Sad

### C4 - Kontakt osobe za praćenje procesa integracije

**Osoba 1:**
- **Ime i prezime:** Dmitrii Voroshilov
- **Email adresa:** dima@svetu.rs

**Osoba 2:**
- **Ime i prezime:** Azamat Salakhov
- **Email adresa:** azamat@svetu.rs

**Osoba 3:**
- **Ime i prezime:** [ostaviti prazno]
- **Email adresa:** [ostaviti prazno]

### C5 - Kontakt osobe odgovorne za razvoj u procesu integracije

**Osoba 1:**
- **Ime i prezime:** Azamat Salakhov
- **Email adresa:** azamat@svetu.rs

**Osoba 2:**
- **Ime i prezime:** Dmitrii Voroshilov
- **Email adresa:** dima@svetu.rs

**Osoba 3:**
- **Ime i prezime:** [ostaviti prazno]
- **Email adresa:** [ostaviti prazno]

### C6 - Mail adresa osobe ili distributivne grupe na koju će korisnik primati informacije vezane za buduće izmene i radove na servisima na D Express strani

**info@svetu.rs**

---

## D - Ukoliko je kompanija koja izdaje zahtev nema svoj razvoj unutar kompanije ili ne koristi svoj razvoj, navesti podatke integratora

*[BLOK D SE PRESKAČE - kompanija ima svoj razvoj]*

### D1 - Naziv kompanije integratora
*[preskočeno]*

### D2 - UK broj kompanije integratora
*[preskočeno]*

### D3 - Adresa sedišta kompanije integratora
*[preskočeno]*

### D4 - Kontakt osobe za praćenje procesa integracije na strani integratora
*[preskočeno]*

### D5 - Kontakt osobe odgovorne za razvoj u procesu integracije na strani integratora
*[preskočeno]*

### D6 - Mail adresa osobe ili distributivne grupe na koju će tehnička lica primati informacije
*[preskočeno]*

### D7 - Ko daje instrukcije i informacije programerskoj kući
*[preskočeno]*

---

## E - Tehnologije kojom se vrši razvoj na strani korisnika

### Opis infrastrukture - platforma, programski jezici, svoj hosting ili cloud

**Backend:**
- Programski jezik: Go (Golang) 1.21+
- Framework: Fiber v2
- Baza podataka: PostgreSQL 15
- Cache: Redis
- Search engine: OpenSearch

**Frontend:**
- Framework: Next.js 15 (React 19)
- Jezik: TypeScript
- CSS: Tailwind CSS + DaisyUI

**Infrastruktura:**
- Hosting: Vlastiti VPS serveri
- Containerization: Docker
- Orchestration: Docker Compose
- OS: Ubuntu 22.04 LTS
- Web server: Nginx (reverse proxy)

**API:**
- Protocol: REST API
- Format: JSON
- Autentifikacija: JWT tokens
- Dokumentacija: Swagger/OpenAPI

---

## Pre započinjanja procesa integracije kod klijenta je utvrđeno sledeće

## F - Tip distribucije:

☑️ **Pošiljke se šalju iz sopstvenog skladišta (ili jedne sopstvene lokacije)**  
☐ Pošiljke se šalju iz zakupljenog skladišta (uslužno skladištenje i odvajanje)  
☐ Pošiljke se šalju sa sopstvenih različitih lokacija  
☐ U pitanju je drop shipping - roba se preuzima od različitih dobavljača

**Činjenice je utvrdio:**  
*[Ime i prezime D Express predstavnika]*

---

## G - Način unosa podataka:

☐ Podaci o primaocima se unose u interni sistem korisnika, unos vrše zaposleni u kompaniji  
☑️ **Podaci o primaocima se unose na portalu korisnika, od strane trećih lica**

**Činjenice je utvrdio:**  
*[Ime i prezime D Express predstavnika]*

---

## H - Očekivani broj pošiljaka na mesečnom nivou:

**1,000 pošiljaka mesečno (početak)**  
**Projekcija rasta: 5,000 pošiljaka mesečno nakon 6 meseci**

---

## I - Očekivani datumi početka i kraja integracije:

### Početak integracije
**15.09.2025.**

### Početak testiranja
**25.09.2025.**

### Kraj testiranja / puštanje u produkciju
**01.10.2025.**

---

## Dodatne informacije

### Lokacija skladišta:
Novi Sad, ul. Kisačka 27, stan 12

### Tipovi pošiljaka koje će se koristiti:
- Standardne pošiljke
- Pošiljke sa otkupninom (COD) - većinski deo
- Express dostava (D Express Today)
- Dostava u D Express poslovnice

### Potrebne API funkcionalnosti:
- Kreiranje naloga za dostavu
- Generisanje i štampanje adresnica
- Praćenje pošiljaka u realnom vremenu
- Validacija adresa
- Kalkulacija cene dostave
- Upravljanje povratima
- Masovno kreiranje naloga (batch processing)

### Napomena o integraciji:
Već imamo iskustva sa integracijom Post Express WSP API sistema, tako da smo upoznati sa osnovnim principima rada sa kurirskim API servisima.

---

**Potpis podnosioca zahteva:** _______________________

**Datum:** 05.09.2025.
