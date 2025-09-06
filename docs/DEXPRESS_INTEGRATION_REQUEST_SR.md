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
*[potrebno dodati matični broj ako postoji]*

### C3 - Adresa sedišta kompanije
- **Ulica:** Vase Stajića 18
- **Broj:** 18
- **PTT:** 21101
- **Mesto:** Novi Sad

### C4 - Kontakt osobe za praćenje procesa integracije

**Osoba 1:**
- **Ime i prezime:** Dmitrii Voroshilov
- **Email:** dima@svetu.rs

**Osoba 2:**
- **Ime i prezime:** Azamat Salakhov
- **Email:** azamat@svetu.rs

### C5 - Kontakt osobe odgovorne za razvoj u procesu integracije

**Osoba 1:**
- **Ime i prezime:** Azamat Salakhov
- **Email:** azamat@svetu.rs

**Osoba 2:**
- **Ime i prezime:** Dmitrii Voroshilov
- **Email:** dima@svetu.rs

### C6 - Mail adresa osobe ili distributivne grupe na koju će korisnik primati informacije vezane za buduće izmene i radove na servisima na D Express strani

**info@svetu.rs**

---

## D - Ukoliko je kompanija koja izdaje zahtev nema svoj razvoj unutar kompanije ili ne koristi svoj razvoj, navesti podatke integratora

*[PRESKOČENO - kompanija ima svoj razvoj]*

---

## E - Tehnologije kojom se vrši razvoj na strani korisnika

### Platforma
- **Backend:** Go (Golang)
- **Frontend:** React/Next.js
- **Baza podataka:** PostgreSQL
- **API:** RESTful API sa JSON formatom
- **Infrastruktura:** Docker, Kubernetes

---

## Pre započinjanja procesa integracije kod klijenta je utvrđeno sledeće

## F - Tip distribucije:

☑️ **Pošiljke se šalju iz sopstvenog skladišta (ili jedne sopstvene lokacije)**  
☐ Pošiljke se šalju iz zakupljenog skladišta (uslužno skladištenje i odvajanje)  
☐ Pošiljke se šalju sa sopstvenih različitih lokacija  
☐ U pitanju je drop shipping - roba se preuzima od različitih dobavljača

**Napomena:** Skladište se nalazi na adresi: Novi Sad, ul. Kisačka 27, stan 12

**Činjenice je utvrdio:**  
*[Ime i prezime D Express predstavnika - popunjava D Express]*

---

## G - Način unosa podataka:

☐ Podaci o primaocima se unose u interni sistem korisnika, unos vrše zaposleni u kompaniji  
☑️ **Podaci o primaocima se unose na portalu korisnika, od strane trećih lica**

**Napomena:** Kupci sami unose podatke za dostavu kroz naš marketplace portal

**Činjenice je utvrdio:**  
*[Ime i prezime D Express predstavnika - popunjava D Express]*

---

## H - Očekivani broj pošiljaka na mesečnom nivou:

**1,000 pošiljaka** (početak)  
**5,000 pošiljaka** (nakon 6 meseci)

---

## I - Očekivani datumi početka i kraja integracije:

### Početak integracije
**15.09.2025.**

### Početak testiranja
**25.09.2025.**

### Kraj testiranja / puštanje u produkciju
**01.10.2025.**

---

## Dodatne napomene

### Potrebni API servisi:
- Kreiranje naloga za dostavu
- Štampanje adresnica
- Praćenje pošiljaka (tracking)
- Provera poštanskih brojeva i adresa
- Kalkulacija cene dostave
- Upravljanje povratima

### Tipovi pošiljaka:
- Standardne pošiljke
- Pošiljke sa otkupninom (COD)
- Express dostava
- Dostava u D Express poslovnice



---

**Datum podnošenja:** 05.09.2025.  
**Potpis podnosioca:** _______________________
