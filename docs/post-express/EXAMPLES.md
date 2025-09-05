# WSP WebAPI Examples

## 1. Поиск населенного пункта (GetNaselje)

### Запрос
```http
POST https://onlinepostexpress.rs/WSPWebApi/api/app/transakcija
Content-Type: application/json

{
  "TransakcijaId": 3,
  "DatumVremePosiljke": "2024-01-25T10:00:00",
  "Klijent": {
    "Username": "your_username",
    "Password": "your_password"
  },
  "NaseljeIn": {
    "Naziv": "Београд",
    "Ptt": ""
  }
}
```

### Ответ
```json
{
  "OK": true,
  "Poruka": "",
  "TransakcijaId": 3,
  "DatumVremePrijema": "2024-01-25T10:00:01",
  "NaseljeOut": {
    "OK": true,
    "Poruka": "",
    "Naselja": [
      {
        "Sifra": 1,
        "Naziv": "Београд (Стари град)",
        "Ptt": "11000",
        "Opstina": "Стари град"
      },
      {
        "Sifra": 2,
        "Naziv": "Београд (Звездара)",
        "Ptt": "11050",
        "Opstina": "Звездара"
      },
      {
        "Sifra": 3,
        "Naziv": "Београд (Нови Београд)",
        "Ptt": "11070",
        "Opstina": "Нови Београд"
      }
    ]
  }
}
```

## 2. Поиск улицы (GetUlica)

### Запрос
```json
{
  "TransakcijaId": 4,
  "DatumVremePosiljke": "2024-01-25T10:05:00",
  "Klijent": {
    "Username": "your_username",
    "Password": "your_password"
  },
  "UlicaIn": {
    "SifraMesta": 1,
    "Naziv": "Кнез Михаилова"
  }
}
```

### Ответ
```json
{
  "OK": true,
  "Poruka": "",
  "TransakcijaId": 4,
  "DatumVremePrijema": "2024-01-25T10:05:01",
  "UlicaOut": {
    "OK": true,
    "Poruka": "",
    "Ulice": [
      {
        "Sifra": 101,
        "Naziv": "Кнез Михаилова",
        "SifraMesta": 1,
        "NazivMesta": "Београд (Стари град)",
        "PttMesta": "11000"
      }
    ]
  }
}
```

## 3. Проверка адреса (ProveraAdrese/AddressCheck)

### Запрос
```json
{
  "TransakcijaId": 6,
  "DatumVremePosiljke": "2024-01-25T10:10:00",
  "Klijent": {
    "Username": "your_username",
    "Password": "your_password"
  },
  "ProveraAdreseIn": {
    "SifraMesta": 1,
    "SifraUlice": 101,
    "Broj": "10",
    "Ulaz": "1",
    "Sprat": "3",
    "Stan": "12",
    "Ptt": "11000"
  }
}
```

### Ответ
```json
{
  "OK": true,
  "Poruka": "",
  "TransakcijaId": 6,
  "DatumVremePrijema": "2024-01-25T10:10:01",
  "ProveraAdreseOut": {
    "OK": true,
    "Poruka": "",
    "Adresa": {
      "DostupanKurir": true,
      "DostupanPostExpress": true,
      "TipKucnogBroja": "REGULAR"
    }
  }
}
```

## 4. Проверка доступности услуги (ProveraDostupnostiUsluge)

### Запрос
```json
{
  "TransakcijaId": 23,
  "DatumVremePosiljke": "2024-01-25T10:15:00",
  "Klijent": {
    "Username": "your_username",
    "Password": "your_password"
  },
  "ProveraDostupnostiUslugeIn": {
    "SifraMesta": 1,
    "SifraUlice": 101,
    "KucniBroj": "10",
    "SifraUsluge": 58
  }
}
```

### Ответ
```json
{
  "OK": true,
  "Poruka": "",
  "TransakcijaId": 23,
  "DatumVremePrijema": "2024-01-25T10:15:01",
  "ProveraDostupnostiUslugeOut": {
    "OK": true,
    "Poruka": "",
    "Dostupnost": true
  }
}
```

## 5. Расчет почтовых расходов (PostarinaPosiljke)

### Запрос
```json
{
  "TransakcijaId": 5,
  "DatumVremePosiljke": "2024-01-25T10:20:00",
  "Klijent": {
    "Username": "your_username",
    "Password": "your_password"
  },
  "PostarinaIn": {
    "SifraUsluge": 58,
    "Masa": 1500,
    "Vrednost": 5000,
    "Otkupnina": 0,
    "LicnoPreuzimanje": true,
    "PovratnicaPotpis": false,
    "SMS": true,
    "Potvrda": false,
    "SifraMestaPrimaoca": 2,
    "BrojPosiljaka": 1
  }
}
```

### Ответ
```json
{
  "OK": true,
  "Poruka": "",
  "TransakcijaId": 5,
  "DatumVremePrijema": "2024-01-25T10:20:01",
  "PostarinaOut": {
    "OK": true,
    "Poruka": "",
    "Postarina": {
      "UkupnaCena": 450.00,
      "CenaBezPDV": 375.00,
      "PDV": 75.00,
      "OsnovnaCena": 350.00,
      "DodatneUsluge": 25.00
    }
  }
}
```

## 6. Отслеживание отправления (TTKretanjaUsluge)

### Запрос
```json
{
  "TransakcijaId": 63,
  "DatumVremePosiljke": "2024-01-25T10:25:00",
  "Klijent": {
    "Username": "your_username",
    "Password": "your_password"
  },
  "TTKretanjeIn": {
    "BrojPosiljke": "RF123456789RS",
    "Jezik": 1
  }
}
```

### Ответ
```json
{
  "OK": true,
  "Poruka": "",
  "TransakcijaId": 63,
  "DatumVremePrijema": "2024-01-25T10:25:01",
  "TTKretanjeOut": {
    "OK": true,
    "Poruka": "",
    "Kretanja": [
      {
        "DatumVreme": "2024-01-24T09:00:00",
        "Mesto": "Београд 11000",
        "Status": "ПРИХВАЋЕНО",
        "Opis": "Пошиљка прихваћена у пошти"
      },
      {
        "DatumVreme": "2024-01-24T14:00:00",
        "Mesto": "Сортирни центар Београд",
        "Status": "У ТРАНЗИТУ",
        "Opis": "Пошиљка у сортирном центру"
      },
      {
        "DatumVreme": "2024-01-25T08:00:00",
        "Mesto": "Нови Сад 21000",
        "Status": "НА ДОСТАВУ",
        "Opis": "Пошиљка спремна за доставу"
      },
      {
        "DatumVreme": "2024-01-25T10:15:00",
        "Mesto": "Нови Сад 21000",
        "Status": "УРУЧЕНО",
        "Opis": "Пошиљка уручена примаоцу"
      }
    ]
  }
}
```

## 7. Групповое отслеживание (TTPosiljkeStatusi)

### Запрос
```json
{
  "TransakcijaId": 64,
  "DatumVremePosiljke": "2024-01-25T10:30:00",
  "Klijent": {
    "Username": "your_username",
    "Password": "your_password"
  },
  "TTStatusiPosiljakaIn": {
    "BrojeviPosiljaka": [
      "RF123456789RS",
      "RF123456790RS",
      "RF123456791RS"
    ],
    "Jezik": 1
  }
}
```

### Ответ
```json
{
  "OK": true,
  "Poruka": "",
  "TransakcijaId": 64,
  "DatumVremePrijema": "2024-01-25T10:30:01",
  "TTStatusiPosiljakaOut": {
    "OK": true,
    "Poruka": "",
    "Statusi": [
      {
        "BrojPosiljke": "RF123456789RS",
        "Status": "УРУЧЕНО",
        "DatumStatusa": "2024-01-25T10:15:00",
        "MestoStatusa": "Нови Сад 21000",
        "OpisStatusa": "Пошиљка уручена примаоцу",
        "Isporuceno": true
      },
      {
        "BrojPosiljke": "RF123456790RS",
        "Status": "НА ДОСТАВУ",
        "DatumStatusa": "2024-01-25T08:00:00",
        "MestoStatusa": "Нови Сад 21000",
        "OpisStatusa": "Пошиљка спремна за доставу",
        "Isporuceno": false
      },
      {
        "BrojPosiljke": "RF123456791RS",
        "Status": "У ТРАНЗИТУ",
        "DatumStatusa": "2024-01-24T14:00:00",
        "MestoStatusa": "Сортирни центар Београд",
        "OpisStatusa": "Пошиљка у сортирном центру",
        "Isporuceno": false
      }
    ]
  }
}
```

## 8. B2B Express Handling (Handling 58 - Express Today)

### Запрос
```json
{
  "TransakcijaId": 58,
  "DatumVremePosiljke": "2024-01-25T11:00:00",
  "Klijent": {
    "Username": "company_username",
    "Password": "company_password"
  },
  "HandlingIn": {
    "Posiljalac": {
      "Naziv": "Компанија АБЦ д.о.о.",
      "SifraMesta": 1,
      "SifraUlice": 101,
      "KucniBroj": "25",
      "Telefon": "+381601234567",
      "Email": "office@kompanija.rs",
      "PIB": "100123456",
      "MaticniBroj": "12345678"
    },
    "Primalac": {
      "Ime": "Петар",
      "Prezime": "Петровић",
      "SifraMesta": 2,
      "SifraUlice": 202,
      "KucniBroj": "30",
      "Stan": "5",
      "Telefon": "+381607654321",
      "Email": "petar@example.com"
    },
    "Posiljka": {
      "SifraUsluge": 58,
      "Masa": 2500,
      "Vrednost": 10000,
      "Otkupnina": 0,
      "Sadrzaj": "Документи",
      "BrojPaketa": 1,
      "Napomena": "Хитна достава",
      "ReferentiBroj": "ORD-2024-001"
    },
    "DodatneUsluge": {
      "LicnoPreuzimanje": true,
      "SMS": true,
      "Email": true,
      "PovratnicaPotpis": true
    },
    "NacinPlacanja": 1,
    "VremePrihvatanja": "2024-01-25T11:00:00",
    "VremeIsporuke": "2024-01-25T16:00:00"
  }
}
```

### Ответ
```json
{
  "OK": true,
  "Poruka": "",
  "TransakcijaId": 58,
  "DatumVremePrijema": "2024-01-25T11:00:01",
  "HandlingOut": {
    "OK": true,
    "Poruka": "",
    "BrojPosiljke": "EE123456789RS",
    "Barkod": "EE123456789RS",
    "Cena": 850.00,
    "CenaBezPDV": 708.33,
    "PDV": 141.67,
    "VremePrihvatanja": "2024-01-25T11:00:00",
    "PlaniranoVremeIsporuke": "2024-01-25T16:00:00"
  }
}
```

## 9. B2B Manifest - Создание манифеста

### Запрос
```json
{
  "TransakcijaId": 100,
  "DatumVremePosiljke": "2024-01-25T12:00:00",
  "Klijent": {
    "Username": "company_username",
    "Password": "company_password"
  },
  "B2BManifestIn": {
    "TipOperacije": 1,
    "DatumManifesta": "2024-01-25",
    "Posiljke": [
      {
        "Posiljalac": {
          "Naziv": "Компанија XYZ",
          "SifraMesta": 1,
          "SifraUlice": 101,
          "KucniBroj": "10"
        },
        "Primalac": {
          "Ime": "Марко",
          "Prezime": "Марковић",
          "SifraMesta": 3,
          "SifraUlice": 303,
          "KucniBroj": "15"
        },
        "Posiljka": {
          "SifraUsluge": 71,
          "Masa": 1000,
          "Sadrzaj": "Пакет"
        }
      }
    ]
  }
}
```

### Ответ
```json
{
  "OK": true,
  "Poruka": "",
  "TransakcijaId": 100,
  "DatumVremePrijema": "2024-01-25T12:00:01",
  "B2BManifestOut": {
    "OK": true,
    "Poruka": "",
    "ManifestId": "MAN-2024-0125-001",
    "BrojPosiljaka": 1,
    "Status": "ОТВОРЕН",
    "DatumKreiranja": "2024-01-25T12:00:01",
    "Posiljke": [
      {
        "BrojPosiljke": "CE123456789RS",
        "Barkod": "CE123456789RS",
        "Status": "КРЕИРАНО",
        "Cena": 450.00
      }
    ]
  }
}
```

## Обработка ошибок

### Ошибка аутентификации
```json
{
  "OK": false,
  "Poruka": "Неисправно корисничко име или лозинка",
  "TransakcijaId": 3,
  "DatumVremePrijema": "2024-01-25T10:00:01"
}
```

### Ошибка валидации
```json
{
  "OK": false,
  "Poruka": "Недовољна дужина назива за претрагу (минимум 3 карактера)",
  "TransakcijaId": 3,
  "DatumVremePrijema": "2024-01-25T10:00:01"
}
```

### Данные не найдены
```json
{
  "OK": false,
  "Poruka": "Нису пронађена насеља по задатим критеријумима",
  "TransakcijaId": 3,
  "DatumVremePrijema": "2024-01-25T10:00:01"
}
```

### Сервис недоступен
```json
{
  "OK": false,
  "Poruka": "Сервис је тренутно недоступан. Покушајте поново касније.",
  "TransakcijaId": 3,
  "DatumVremePrijema": "2024-01-25T10:00:01"
}
```

## Примечания

1. **Аутентификация**: Все запросы требуют валидные Username и Password в объекте Klijent
2. **Формат даты**: Используйте ISO 8601 формат (YYYY-MM-DDTHH:mm:ss)
3. **Кодировка**: Все запросы и ответы используют UTF-8
4. **Content-Type**: Всегда используйте `application/json`
5. **Язык ответов**: Параметр Jezik (1-кириллица, 2-латиница, 3-английский)
6. **Массовые операции**: Максимум 100 отправлений для группового отслеживания
7. **Коды населенных пунктов**: Используйте метод GetNaselje для получения кодов
8. **Коды улиц**: Используйте метод GetUlica для получения кодов улиц