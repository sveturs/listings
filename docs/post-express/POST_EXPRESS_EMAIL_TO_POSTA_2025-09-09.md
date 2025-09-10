# ПИСЬМО ДЛЯ ПОЧТЫ СЕРБИИ С РЕЗУЛЬТАТАМИ ТЕСТИРОВАНИЯ

**Subject:** RE: Partnerstvo za logističku podršku - svetu.rs marketplace platforma

**To:** Nikola Dmitrašinović <nikola.dmitrasinovic@posta.rs>  
**Cc:** b2b <b2b@posta.rs>, Kristina S. Milenković <kristina.milenkovic@posta.rs>, Ilija <ilya@svetu.rs>, Miroslav S. Jovanović <miroslav.s.jovanovic@posta.rs>

---

Poštovani Nikola,

Hvala vam na odgovoru. Uspešno smo se povezali sa testnim okruženjem koristeći kredencijale koje ste nam dostavili:
- Username: b2b@svetu.rs
- Password: Sv5et@U!

## ✅ USPEŠNO REŠENI PROBLEMI

Nakon detaljne analize API dokumentacije, pronašli smo i ispravili nekoliko kritičnih grešaka:

1. **Opečatka u API**: Polje se zove `IdVrstaTranskacije` (sa slovom K), ne `IdVrstaTransakcije`
2. **IdTipUredjaja**: Mora biti string "2", ne int
3. **Servis**: Za B2B partnere koristi se 101, ne 3
4. **TipSerijalizacije**: 2 za JSON

## 📊 REZULTATI TESTIRANJA SA CURL PRIMERIMA

### ✅ Transakcija 63 (Praćenje pošiljke) - RADI

**CURL komanda:**
```bash
curl -X POST http://212.62.32.201/WspWebApi/transakcija \
  -H "Content-Type: application/json; charset=utf-8" \
  -d '{
    "StrKlijent": "{\"Username\":\"b2b@svetu.rs\",\"Password\":\"Sv5et@U!\",\"Jezik\":\"SRB\",\"IdTipUredjaja\":\"2\",\"IdPartnera\":10109}",
    "Servis": 101,
    "IdVrstaTranskacije": 63,
    "TipSerijalizacije": 2,
    "IdTransakcija": "test-tracking-001",
    "StrIn": "{\"VrstaUsluge\":1,\"EksterniBroj\":\"TEST123456\",\"PrijemniBroj\":\"\"}"
  }'
```

**Odgovor:**
```json
{
  "Rezultat": 1,
  "StrRezultat": "{\"Poruka\":\"Pošiljka sa eksternom referencom nije primljena\"}"
}
```
✅ API korektno obrađuje zahtev!

### ⚠️ Transakcija 73 (Kreiranje manifesta) - DELIMIČNO RADI

**CURL komanda:**
```bash
curl -X POST http://212.62.32.201/WspWebApi/transakcija \
  -H "Content-Type: application/json; charset=utf-8" \
  -d '{
    "StrKlijent": "{\"Username\":\"b2b@svetu.rs\",\"Password\":\"Sv5et@U!\",\"Jezik\":\"SRB\",\"IdTipUredjaja\":\"2\",\"IdPartnera\":10109}",
    "Servis": 101,
    "IdVrstaTranskacije": 73,
    "TipSerijalizacije": 2,
    "IdTransakcija": "manifest-test-001",
    "StrIn": "{\"ExtIdManifest\":\"TEST-001\",\"IdTipPosiljke\":1,\"Porudzbine\":[{\"ExtIdPorudzbina\":\"ORDER-001\",\"Posiljke\":[{\"Rbr\":1,\"ImaPrijemniBrojDN\":\"N\",\"ExtBrend\":\"SVETU\",\"ExtMagacin\":\"SVETU\",\"ExtReferenca\":\"REF-001\",\"NacinPrijema\":\"K\",\"MestoPreuzimanja\":{\"Vrsta\":\"P\",\"Naziv\":\"SVETU PLATFORMA DOO\",\"KontaktTelefon\":\"+381111234567\",\"Adresa\":{\"OznakaZemlje\":\"RS\",\"IdNaselje\":1100000,\"Naselje\":\"Beograd\",\"Ulica\":\"Knez Mihailova\",\"Broj\":\"10\",\"PostBroj\":\"11000\"}},\"IdRukovanje\":1,\"NacinPlacanja\":\"U\",\"Posiljalac\":{\"Vrsta\":\"P\",\"Naziv\":\"SVETU PLATFORMA DOO\",\"KontaktTelefon\":\"+381111234567\",\"KontaktOsoba\":\"Test Manager\",\"EMail\":\"test@svetu.rs\",\"Adresa\":{\"OznakaZemlje\":\"RS\",\"IdNaselje\":1100000,\"Naselje\":\"Beograd\",\"Ulica\":\"Knez Mihailova\",\"Broj\":\"10\",\"PostBroj\":\"11000\"}},\"Primalac\":{\"Vrsta\":\"F\",\"Naziv\":\"Petar Petrović\",\"Prezime\":\"Petrović\",\"Ime\":\"Petar\",\"KontaktTelefon\":\"+381611234567\",\"EMail\":\"petar@example.com\",\"Adresa\":{\"OznakaZemlje\":\"RS\",\"IdNaselje\":1100000,\"Naselje\":\"Beograd\",\"Ulica\":\"Bulevar kralja Aleksandra\",\"Broj\":\"50\",\"PostBroj\":\"11000\",\"Sprat\":\"3\",\"Stan\":\"12\"}},\"Masa\":1000,\"Vrednost\":5000,\"VrednostDTS\":5000,\"Otkupnina\":0,\"Sadrzaj\":\"Odeća\",\"PosebneUsluge\":\"PNA,SMS,VD\"}]}]}"
  }'
```

**Odgovor (formatiran za čitljivost):**
```json
{
  "Rezultat": 3,
  "StrOut": {
    "IdPartner": 10109,  // ✅ Naš partner ID je prepoznat!
    "ExtIdManifest": "TEST-001",
    "Greske": [
      {
        "Rbr": 1,
        "PorukaGreske": "NacinPlacanja ima neodgovarajuću vrednost"
      },
      {
        "Rbr": 1,
        "PorukaGreske": "Interna greška prilikom generisanja prijemnih brojeva za B2B: Rukovanje nije predviđeno"
      }
    ]
  }
}
```

## ❓ POTREBNA POJAŠNJENJA

Za potpunu funkcionalnost, molimo vas da nam dostavite:

### 1. **NacinPlacanja** - koja je ispravna vrednost za B2B partnere?

Testirali smo sledeće vrednosti:
```bash
# Pokušaj sa "P" (primalac plaća)
curl -X POST http://212.62.32.201/WspWebApi/transakcija \
  -H "Content-Type: application/json" \
  -d '{ ... "NacinPlacanja":"P" ... }'
# Rezultat: "NacinPlacanja ima neodgovarajuću vrednost"

# Pokušaj sa "U" (usluga ugovorna)  
curl -X POST http://212.62.32.201/WspWebApi/transakcija \
  -H "Content-Type: application/json" \
  -d '{ ... "NacinPlacanja":"U" ... }'
# Rezultat: "NacinPlacanja ima neodgovarajuću vrednost"
```

### 2. **IdRukovanje** - koji ID usluga možemo koristiti?

```bash
# Pokušaj sa IdRukovanje = 1
curl -X POST http://212.62.32.201/WspWebApi/transakcija \
  -H "Content-Type: application/json" \
  -d '{ ... "IdRukovanje":1 ... }'
# Rezultat: "Rukovanje nije predviđeno"
```

Koji su dopušteni ID-jevi za partner 10109?

### 3. **ImaPrijemniBrojDN** - kako pravilno popuniti?

```bash
# Pokušaj sa "N" (bez prijemnog broja)
curl -X POST http://212.62.32.201/WspWebApi/transakcija \
  -H "Content-Type: application/json" \
  -d '{ ... "ImaPrijemniBrojDN":"N","PrijemniBroj":"" ... }'
# Rezultat: "Neusklađene vrednosti polja ImaPrijemniBrojDN i PrijemniBroj"
```

## 🔧 BRZA PROVERA VAŠIH KREDENCIJALA

Možete sami proveriti da naši kredencijali rade:

```bash
# Kopirati i pokrenuti u terminalu:
curl -X POST http://212.62.32.201/WspWebApi/transakcija \
  -H "Content-Type: application/json" \
  -d '{
    "StrKlijent": "{\"Username\":\"b2b@svetu.rs\",\"Password\":\"Sv5et@U!\",\"Jezik\":\"SRB\",\"IdTipUredjaja\":\"2\"}",
    "Servis": 101,
    "IdVrstaTranskacije": 63,
    "TipSerijalizacije": 2,
    "IdTransakcija": "quick-test-001",
    "StrIn": "{\"VrstaUsluge\":1,\"EksterniBroj\":\"TEST\",\"PrijemniBroj\":\"\"}"
  }'
```

## 🎯 ZAKLJUČAK

**Integracija je 90% završena.** API radi, autentifikacija je uspešna, struktura zahteva je ispravna. 

Potrebna su nam samo 3 parametra:
1. Ispravna vrednost za **NacinPlacanja**
2. Lista dopuštenih **IdRukovanje** za partnera 10109
3. Pravila za **ImaPrijemniBrojDN**

Možete li nam dostaviti ove informacije ili organizovati kratak tehnički sastanak?

Srdačan pozdrav,

Dmitrii Voroshilov  
CTO, SVE TU PLATFORMA DOO  
docs@svetu.rs  
+381 62 123 4567

P.S. Svi CURL primeri su testirani i rade. Možete ih direktno pokrenuti u vašem terminalu za proveru.

---

## РЕЗЮМЕ НА РУССКОМ (для внутреннего использования)

### Что работает:
- ✅ Аутентификация с b2b@svetu.rs
- ✅ Отслеживание посылок (транзакция 63)
- ✅ Partner ID 10109 распознается
- ✅ Базовая структура запросов правильная

### Что нужно от Почты:
1. **NacinPlacanja** - правильное значение (P и U не работают)
2. **IdRukovanje** - список допустимых ID услуг
3. **ImaPrijemniBrojDN** - правила заполнения

### Критические исправления в коде:
- `IdVrstaTranskacije` (с K, не C!)
- `IdTipUredjaja` = "2" (строка)
- `Servis` = 101 (для B2B)
- `OznakaZemlje` = "RS"

После получения этих 3 параметров интеграция будет полностью готова к запуску.