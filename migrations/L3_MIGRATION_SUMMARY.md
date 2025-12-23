# L3 Categories Migration Summary

## Overview

Created 4 clean migrations to add **324 new L3 categories** (from 26 existing to **350 total**).

## Migration Files

| Migration | Categories Added | Total L3 After | File Size |
|-----------|-----------------|----------------|-----------|
| `20251217200001_l3_elektronika_clean` | 90 | 116 | 26KB |
| `20251217200002_l3_odeca_clean` | 90 | 206 | 28KB |
| `20251217200003_l3_dom_sport_clean` | 80 | 286 | 23KB |
| `20251217200004_l3_ostalo_clean` | 64 | **350** | 19KB |

## Categories Breakdown

### Migration 1: Elektronika (90 new L3)

**Pametni telefoni (5):**
- google-pixel-telefoni, oneplus-telefoni, realme-telefoni, oppo-telefoni, telefoni-budget-do-200e

**Laptop racunari (7):**
- laptop-2-u-1, laptop-chromebook, laptop-macbook, laptop-lenovo-thinkpad, laptop-dell-xps, laptop-hp-pavilion, laptop-asus-rog

**TV i video (15):**
- led-tv-32-43-inca, led-tv-50-55-inca, led-tv-65-75-inca, oled-tv, qled-tv
- smart-tv-android, smart-tv-tizen
- soundbar-2-1, soundbar-5-1, bluetooth-zvucnici, hi-fi-sistemi
- projektori-full-hd, projektori-4k, av-resiveri, kucni-bioskop

**Racunarske komponente (20):**
- Grafičke: graficke-rtx-4000, graficke-rtx-3000, graficke-gtx-1000, graficke-amd-rx-7000, graficke-amd-rx-6000
- Procesori Intel: procesor-intel-i9, procesor-intel-i7, procesor-intel-i5
- Procesori AMD: procesor-amd-ryzen-9, procesor-amd-ryzen-7, procesor-amd-ryzen-5
- RAM: ram-ddr5-32gb, ram-ddr5-16gb, ram-ddr4-16gb, ram-ddr4-8gb
- SSD: ssd-nvme-1tb, ssd-nvme-500gb, ssd-sata-1tb
- Matične: maticna-ploca-intel-z790, maticna-ploca-amd-x670

**Gaming oprema (15):**
- Konzole: playstation-5, playstation-5-igre, xbox-series-x, xbox-series-s, xbox-igre, nintendo-switch, nintendo-igre
- Kontroleri: gejmpad-ps5, gejmpad-xbox
- Periferija: gaming-slusalice, gaming-monitor-144hz, gaming-monitor-240hz, gaming-stolice, rgb-rasveta, streaming-mikrofon

**Foto i video kamere (12):**
- dslr-canon, dslr-nikon, mirrorless-sony, mirrorless-fujifilm
- action-kamera-gopro, action-kamera-dji, dron-sa-kamerom, 4k-video-kamera
- objektiv-za-canon, objektiv-za-nikon, stativ-foto-video, foto-rasveta

**Pametni satovi (8):**
- apple-watch-series, apple-watch-ultra, samsung-galaxy-watch
- garmin-fitness-sat, xiaomi-mi-band, xiaomi-watch, huawei-watch-gt, pametni-satovi-budget

---

### Migration 2: Odeca i Obuca (90 new L3)

**Muska odeca (12):**
- muski-dzemperi, muski-duksevi, muska-sportska-odeca, muske-trenerke, muski-sorcevi
- muski-kaputi, muski-prsluci, muska-odela, muski-smokingzi
- muske-kravate, muski-kaisevi, muski-salovi

**Zenska odeca (11):**
- zenske-pantalone-jeans, zenske-pantalone-elegantne, zenski-dzemperi, zenski-duksevi
- zenska-sportska-odeca, zenske-trenerke, zenski-sorcevi, zenski-kaputi, zenske-mantile
- zenska-vecernja-garderoba, zenska-poslovna-garderoba

**Decija odeca (15):**
- Bebe: bebe-odeca-0-3-meseca, bebe-odeca-3-6-meseci, bebe-odeca-6-12-meseci
- Dečaci: decaci-odeca-1-3-godine, decaci-odeca-4-7-godina, decaci-odeca-8-12-godina
- Devojčice: devojcice-odeca-1-3-godine, devojcice-odeca-4-7-godina, devojcice-odeca-8-12-godina
- Opšte: decije-jakne, deciji-duksevi, decije-trenerke, decije-kaputi, skolska-uniforma, svecana-decija-odeca

**Muska obuca (15):**
- Patike: muske-patike-casual, muske-patike-running, muske-patike-basketball, muske-patike-football
- Cipele: muske-cipele-koza, muske-cipele-oxford, muske-cipele-derby, muske-mokasine, muske-espadrile, muske-cipele-brodske
- Čizme: muske-cizme-duboke, muske-cizme-chelsea, muske-cizme-radne
- Sandale: muske-sandale, muske-papuce

**Zenska obuca (15):**
- Patike: zenske-patike-casual, zenske-patike-running, zenske-patike-fitness
- Štikle: zenske-stikle, zenske-cipele-potpetica
- Ravne: zenske-balerinke, zenske-mokasine, zenske-espadrile
- Čizme: zenske-cizme-duboke, zenske-cizme-gleznjace, zenske-cizme-preko-kolena
- Ostalo: zenske-sandale, zenske-papuce, zenske-cipele-platforme, zenske-natikace

**Decija obuca (15):**
- Bebe: bebe-obuca-0-6-meseci, bebe-obuca-6-12-meseci, decija-obuca-prve-korake
- Patike po godinama: decije-patike-1-3-godine, decije-patike-4-7-godina, decije-patike-8-12-godina
- Sport: decije-sportske-patike, decije-fudbalske-kopacke
- Čizme: decije-cizme-gumene, decije-cizme-zimske, decije-cizme-duboke
- Ostalo: decije-sandale, decije-papuce, decije-cipele-skolske, decije-baletanke

**Kupaci kostimi (7):**
- muski-kupaci-sorcevi, muski-kupaci-slip
- zenski-tankini, zenski-monokini, deciji-kupaci-kostimi
- kupaci-majice-uv-zastita, pareo-tunike

---

### Migration 3: Dom i Sport (80 new L3)

**Namestaj (15):**
- Dnevna: trosed-dvosed, ugaona-garnitura, fotelja, tv-komoda, sto-za-dnevnu-sobu
- Spavaća: bracni-krevet, samacki-krevet, sprat-krevet, orman, komode, nocni-stocic
- Trpezarija: trpezarijski-sto, trpezarijske-stolice
- Kancelarija: kancelarijski-sto, kancelarijska-stolica

**Kupatilo (12):**
- Lavabo: lavabo, ugradni-lavabo
- Kada: kada, hidromasazna-kada
- Tuš: tus-kabina, tus-set
- WC: wc-solja, bide
- Slavine: slavina-za-lavabo, slavina-za-kadu
- Nameštaj: ogledalo-sa-ormarićem, kupatilski-namestaj

**Kuhinja (13):**
- Elementi: kuhinjski-elementi, kuhinjska-radna-ploca
- Sudopera: kuhinjska-sudopera, kuhinjska-slavina
- Aparati: aspirator, stednjak, ugradna-rerna, ploca-za-kuvanje
- Hlađenje: masina-za-pranje-sudova, frizider, zamrzivac
- Ostalo: mikrotalasna-rerna, mini-bojler

**Rasveta (10):**
- Plafon: luster, plafonjera, ugradna-led-rasveta
- Zid: zidna-lampa
- Pod: stojeca-lampa
- Sto: stonalampa
- Spoljna: spoljna-rasveta
- Sijalice: led-sijalice
- Smart: pametna-rasveta
- Dekor: dekorativna-rasveta

**Fitnes oprema (10):**
- Kardio: traka-za-trcanje, bicikl-sobni, elipticni-trenazer, veslo-masina
- Tegovi: bucice, sipka-i-tegovi, klupa-za-vezbanje
- Rekviziti: podloga-za-jogu, lopta-za-pilates, trx-trake

**Bicikli (10):**
- Tipovi: brdski-bicikl, drumski-bicikl, gradski-bicikl, elektricni-bicikl, sklopivi-bicikl, deciji-bicikl, bmx-bicikl
- Oprema: kaciga-za-bicikl, biciklisticke-rukavice, biciklisticka-torba

**Kampovanje (10):**
- Šatori: sator-2-3-osobe, sator-4-6-osoba
- Spavanje: spavaca-vreca, samoduvajuci-dusek
- Ranac: kamp-ranac
- Kuhinja: kamp-rostilj, gas-reaud
- Rasveta: kamp-lampa
- Ostalo: prenosivi-frizider, kamp-stolica

---

### Migration 4: Ostalo (64 new L3)

**Kosa i frizura (10):**
- Aparati: fen-za-kosu, peglazakosu, figaro, trimer-za-kosu, masnica-za-kosu
- Nega: sampon, regenerator-za-kosu, maska-za-kosu, farba-za-kosu, sprej-za-kosu

**Nega koze (10):**
- Čišćenje: mleko-za-ciscenje, tonik-za-lice
- Hidratacija: krema-za-lice, serum-za-lice
- Maske: maska-za-lice
- Zaštita: krema-za-suncanje, anti-age-krema
- Telo: krema-za-telo, gel-za-tusiranje, krema-za-ruke

**Bebe oprema (12):**
- Kolica: kolica-za-bebe, nosiljka-za-bebe, auto-sediste-za-bebe
- Nameštaj: krevetac-za-bebe, hranilica, ogradica-za-bebe
- Monitoring: baby-alarm
- Hranjenje: bocice-za-bebe, pumpa-za-mleko
- Kupanje: kadica-za-kupanje
- Ostalo: pelene, bebe-igracke

**Auto delovi (10):**
- Gume: gume-ljetne, gume-zimske
- Točkovi: aluminijske-felne
- Elektrika: akumulator, auto-farovi, led-auto-sijalice
- Kočnice: kocione-plocice
- Filteri: uljni-filter, vazdusni-filter
- Brisači: brisaci

**Aparati za kucu (10):**
- Pranje: masina-za-pranje-vesa, masina-za-susenje
- Čišćenje: usisivac, robotski-usisivac
- Peglanje: pegla, pegla-sa-parom
- Klima: klima-uredjaj, grejalica, preciscivac-vazduha, razvlazivac

**Knjige (6):**
- Fikcija: romani, price
- Non-fikcija: biografije, popularna-nauka
- Ostalo: decije-knjige, stripovi

**Hrana za ljubimce (6):**
- Psi: hrana-za-pse-suva, hrana-za-pse-konzerve
- Mačke: hrana-za-macke-suva, hrana-za-macke-konzerve
- Ostalo: hrana-za-ptice, hrana-za-ribe

---

## Quality Checks

### ✅ SQL Syntax
- No double-quote (") symbols in JSON fields (used "inča" instead)
- No unescaped apostrophes in English translations
- All parent_id references use subqueries with `level = 2` check
- Consistent level = 3 for all L3 categories
- All slugs lowercase with hyphens

### ✅ No Duplicates
Checked against existing 26 L3 slugs:
- apple-iphone, energetski-napici, gaming-tastature, gaming-misevi, gu-bezicne-slusalice
- huawei-telefoni, kupaci-kostimi-bikini, kupaci-kostimi-celi, laptop-gaming, laptop-profesionalni
- laptop-ultrabook, majice-decake, majice-devojcice, maske-za-telefone, muske-jakne
- muske-kosulje, muske-pantalone, pantalone-decake, punjaci-za-telefone, samsung-telefoni
- xiaomi-telefoni, zastitno-staklo, zenske-bluze, zenske-haljine, zenske-jakne, zenske-suknje

### ✅ Translations
All categories have 3 language translations (sr/en/ru) in JSONB format.

### ✅ Naming Conventions
- Sentence case: "Pametni telefoni" (not "Pametni Telefoni")
- No translated brand names: Samsung, Apple, Nike
- Metric system: cm, kg, inča (not imperial)
- Consistent technical terms: Bluetooth, LED, WiFi

---

## How to Apply

```bash
cd /p/github.com/vondi-global/listings

# Apply migrations in order
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" \
  -f migrations/20251217200001_l3_elektronika_clean.up.sql

psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" \
  -f migrations/20251217200002_l3_odeca_clean.up.sql

psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" \
  -f migrations/20251217200003_l3_dom_sport_clean.up.sql

psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" \
  -f migrations/20251217200004_l3_ostalo_clean.up.sql

# Verify count
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" \
  -c "SELECT level, COUNT(*) FROM categories GROUP BY level ORDER BY level;"
```

Expected result:
```
 level | count
-------+-------
     1 |    18
     2 |   400
     3 |   350
```

---

## Rollback

```bash
# Rollback in reverse order
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" \
  -f migrations/20251217200004_l3_ostalo_clean.down.sql

psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" \
  -f migrations/20251217200003_l3_dom_sport_clean.down.sql

psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" \
  -f migrations/20251217200002_l3_odeca_clean.down.sql

psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" \
  -f migrations/20251217200001_l3_elektronika_clean.down.sql
```

---

## Final Category Tree

```
L1: 18 categories (root)
├── L2: 400 subcategories
    └── L3: 350 leaf categories (26 existing + 324 new)
```

**Total categories:** 768

---

**Created:** 2025-12-17
**Author:** Claude Code (Content Manager)
