# Zahtev za production kredencijale - Post Express WSP API

**Za:** JP "Pošta Srbije" - Post Express  
**Email:** prodaja@posta.rs  
**Predmet:** Zahtev za zaključenje ugovora i production kredencijale za WSP API - SVE TU PLATFORMA

---

Poštovani,

Obraćamo vam se u ime kompanije **Sve Tu d.o.o.** sa zahtevom za zaključenje ugovora o prenosu Post Express pošiljaka i dobijanje production kredencijala za WSP API integraciju.

## Podaci o kompaniji

**Tačan naziv preduzeća:** Sve Tu d.o.o.  
**Adresa i sedište preduzeća:** Vase Stajića 18/18, 21101 Novi Sad  
**Ime, prezime i funkcija lica ovlašćenog za potpisivanje ugovora:** [Unesite podatke]  
**Broj tekućeg računa i naziv banke:** [Unesite podatke]  
**Matični broj:** [Unesite podatke]  
**PIB:** [Unesite podatke]  
**Registrovani kao obveznik PDV-a:** DA  
**Kontakt osoba i telefon:** [Ime], [Telefon]  
**E-mail:** info@svetu.rs  
**Sajt:** www.svetu.rs  

## Opis poslovanja

**Opis sadržine pošiljaka:** Proizvodi kupljeni putem online marketplace platforme (odeća, obuća, elektronika, kućni aparati, kozmetika i drugi proizvodi)  
**Očekivani broj pošiljaka na nedeljnom nivou:** 500-1000 pošiljaka  
**Vrsta pošiljaka:**  
- ✓ Bez naznačene vrednosti
- ✓ Sa označenom vrednošću  
- ✓ Otkupne pošiljke (naloženi plaćež)

**Adrese na kojima će se vršiti preuzimanje pošiljaka:**  
Inicijalno iz centralnog skladišta na adresi sedišta kompanije, sa mogućnošću proširenja na adrese prodavaca nakon pilot perioda.

## Tehnička integracija

**Zainteresovanost za korišćenje aplikativnih rešenja Pošte:** DA

Već smo implementirali kompletnu WSP API integraciju koja podržava:
- Transakciju ID 3 (GetNaselje) - pretraga naselja
- Transakciju ID 10 (GetPostanskeJedinice) - lista poštanskih jedinica  
- Transakciju ID 15 (PracenjePosiljke) - praćenje pošiljke
- Transakciju ID 20 (StampaNalepnice) - štampa nalepnice
- Transakciju ID 25 (StorniranjePosiljke) - storniranje pošiljke
- Transakciju ID 63 (TTKretanjaUsluge) - praćenje pojedinačne pošiljke
- Transakciju ID 73 (Manifest) - prijem liste pošiljaka

## Potrebni podaci za production

Molimo vas da nam dostavite:

1. **Production kredencijale za WSP API:**
   - Username (production)
   - Password (production)
   - Production endpoint URL (ako se razlikuje od https://wsp.postexpress.rs/api/Transakcija)

2. **Dodatne tehničke informacije:**
   - ID transakcije za kreiranje pošiljke (CreatePosiljka)
   - Struktura podataka za kreiranje pošiljke
   - Proces generisanja bar kodova
   - Webhook endpoints za notifikacije o statusu (ako postoje)

3. **Komercijalni uslovi:**
   - Potvrda cena iz stimulativne ponude broj 2025-sl od 31.07.2025.
   - Ugovor o prenosu Post Express pošiljaka

## Prihvatanje uslova

Prihvatamo uslove iz vaše stimulativne ponude:
- Cene prenosa prema Tabeli 1 (340-790 dinara u zavisnosti od težine)
- Provizija za otkupninu: 45,00 dinara po nalogu
- Uslovi za pošiljke standardnih dimenzija (60x50x50cm)
- Rokovi uručenja "Danas za sutra"

## Kontakt za dalju komunikaciju

**Tehnička podrška:**  
[Ime tehničkog direktora]  
Email: tech@svetu.rs  
Telefon: [Broj]

**Komercijalni kontakt:**  
[Ime direktora]  
Email: info@svetu.rs  
Telefon: [Broj]

Spremni smo da odmah započnemo sa testiranjem nakon dobijanja production kredencijala.

Molimo vas da nam što pre dostavite tražene informacije kako bismo mogli da finalizujemo integraciju i započnemo sa korišćenjem vaših usluga.

S poštovanjem,

**Sve Tu d.o.o.**  
Novi Sad

---

*Prilog: Tehnička specifikacija implementirane integracije*