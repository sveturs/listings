# WSP WebAPI Data Structures

## Базовые классы

### TransakcijaIn (Базовый класс запроса)
```json
{
  "TransakcijaId": int,           // Идентификатор транзакции
  "DatumVremePosiljke": datetime, // Дата и время отправки
  "Klijent": Klijent              // Объект клиента
}
```

### TransakcijaOut (Базовый класс ответа)
```json
{
  "OK": boolean,                   // Успешность операции
  "Poruka": string,                // Сообщение (обычно об ошибке)
  "TransakcijaId": int,            // Идентификатор транзакции
  "DatumVremePrijema": datetime    // Дата и время приема
}
```

### Klijent (Аутентификация)
```json
{
  "Username": string,    // Имя пользователя (обязательно)
  "Password": string     // Пароль (обязательно)
}
```

## Классы для работы с адресами

### NaseljeIn
```json
{
  "Naziv": string,    // Часть названия населенного пункта (мин. 3 символа)
  "Ptt": string       // Почтовый индекс (опционально)
}
```

### NaseljeOut
```json
{
  "OK": boolean,
  "Poruka": string,
  "Naselja": [
    {
      "Sifra": int,        // Код населенного пункта
      "Naziv": string,     // Название
      "Ptt": string,       // Почтовый индекс
      "Opstina": string    // Община
    }
  ]
}
```

### UlicaIn
```json
{
  "SifraMesta": int,    // Код населенного пункта
  "Naziv": string       // Часть названия улицы (мин. 3 символа)
}
```

### UlicaOut
```json
{
  "OK": boolean,
  "Poruka": string,
  "Ulice": [
    {
      "Sifra": int,         // Код улицы
      "Naziv": string,      // Название улицы
      "SifraMesta": int,    // Код населенного пункта
      "NazivMesta": string, // Название населенного пункта
      "PttMesta": string    // Почтовый индекс
    }
  ]
}
```

### ProveraAdreseIn / AddressCheckIn
```json
{
  "SifraMesta": int,     // Код населенного пункта
  "SifraUlice": int,     // Код улицы
  "Broj": string,        // Номер дома
  "Ulaz": string,        // Подъезд (опционально)
  "Sprat": string,       // Этаж (опционально)
  "Stan": string,        // Квартира (опционально)
  "Ptt": string          // Почтовый индекс
}
```

### ProveraAdreseOut / AddressCheckOut
```json
{
  "OK": boolean,
  "Poruka": string,
  "Adresa": {
    "DostupanKurir": boolean,       // Доступность курьерской службы
    "DostupanPostExpress": boolean, // Доступность PostExpress
    "TipKucnogBroja": string        // Тип номера дома
  }
}
```

### ProveraDostupnostiUslugeIn
```json
{
  "SifraMesta": int,     // Код населенного пункта
  "SifraUlice": int,     // Код улицы (опционально)
  "KucniBroj": string,   // Номер дома (опционально)
  "SifraUsluge": int     // Код услуги
}
```

### ProveraDostupnostiUslugeOut
```json
{
  "OK": boolean,
  "Poruka": string,
  "Dostupnost": boolean  // Доступность услуги
}
```

## Классы для расчета стоимости

### PostarinaIn
```json
{
  "SifraUsluge": int,           // Код услуги
  "Masa": float,                // Вес в граммах
  "Vrednost": float,            // Объявленная стоимость
  "Otkupnina": float,          // Сумма наложенного платежа
  "LicnoPreuzimanje": boolean, // Личное вручение
  "PovratnicaPotpis": boolean, // Уведомление о вручении
  "SMS": boolean,              // SMS уведомление
  "Potvrda": boolean,          // Подтверждение
  "SifraMestaPrimaoca": int,   // Код места получателя
  "BrojPosiljaka": int         // Количество отправлений
}
```

### PostarinaOut
```json
{
  "OK": boolean,
  "Poruka": string,
  "Postarina": {
    "UkupnaCena": float,     // Общая стоимость
    "CenaBezPDV": float,     // Стоимость без НДС
    "PDV": float,            // НДС
    "OsnovnaCena": float,    // Базовая стоимость
    "DodatneUsluge": float   // Дополнительные услуги
  }
}
```

## Классы для отслеживания

### TTKretanjeIn
```json
{
  "BrojPosiljke": string,  // Номер отправления (штрих-код)
  "Jezik": int             // Язык (1-кириллица, 2-латиница, 3-английский)
}
```

### TTKretanjeOut
```json
{
  "OK": boolean,
  "Poruka": string,
  "Kretanja": [
    {
      "DatumVreme": datetime,  // Дата и время события
      "Mesto": string,         // Место события
      "Status": string,        // Статус
      "Opis": string          // Описание события
    }
  ]
}
```

### TTStatusiPosiljakaIn
```json
{
  "BrojeviPosiljaka": [string],  // Массив номеров (до 100)
  "Jezik": int                   // Язык ответа
}
```

### TTStatusiPosiljakaOut
```json
{
  "OK": boolean,
  "Poruka": string,
  "Statusi": [
    {
      "BrojPosiljke": string,    // Номер отправления
      "Status": string,          // Текущий статус
      "DatumStatusa": datetime,  // Дата последнего статуса
      "MestoStatusa": string,    // Место последнего статуса
      "OpisStatusa": string,     // Описание статуса
      "Isporuceno": boolean      // Доставлено ли
    }
  ]
}
```

## Классы для B2B операций

### Posiljalac (Отправитель)
```json
{
  "Ime": string,           // Имя
  "Prezime": string,       // Фамилия
  "Naziv": string,         // Название компании
  "SifraMesta": int,       // Код населенного пункта
  "SifraUlice": int,       // Код улицы
  "KucniBroj": string,     // Номер дома
  "Ulaz": string,          // Подъезд
  "Sprat": string,         // Этаж
  "Stan": string,          // Квартира
  "Telefon": string,       // Телефон
  "Email": string,         // Email
  "PIB": string,           // ИНН компании
  "MaticniBroj": string    // Регистрационный номер
}
```

### Primalac (Получатель)
```json
{
  "Ime": string,           // Имя
  "Prezime": string,       // Фамилия
  "Naziv": string,         // Название компании
  "SifraMesta": int,       // Код населенного пункта
  "SifraUlice": int,       // Код улицы
  "KucniBroj": string,     // Номер дома
  "Ulaz": string,          // Подъезд
  "Sprat": string,         // Этаж
  "Stan": string,          // Квартира
  "Telefon": string,       // Телефон
  "Email": string,         // Email
  "PIB": string,           // ИНН компании
  "MaticniBroj": string    // Регистрационный номер
}
```

### Posiljka (Отправление)
```json
{
  "SifraUsluge": int,      // Код услуги
  "Masa": float,           // Вес в граммах
  "Vrednost": float,       // Объявленная стоимость
  "Otkupnina": float,      // Наложенный платеж
  "Sadrzaj": string,       // Содержимое
  "BrojPaketa": int,       // Количество пакетов
  "Napomena": string,      // Примечание
  "ReferentiBroj": string, // Референтный номер клиента
  "BrojUgovora": string    // Номер договора
}
```

### DodatneUsluge (Дополнительные услуги)
```json
{
  "LicnoPreuzimanje": boolean,   // Личное вручение
  "PovratnicaPotpis": boolean,   // Уведомление о вручении
  "SMS": boolean,                // SMS уведомление
  "Email": boolean,              // Email уведомление
  "Potvrda": boolean,            // Подтверждение
  "PlaceniOdgovor": boolean,     // Оплаченный ответ
  "OsiguranjePosiljke": boolean, // Страхование
  "VracanjeDokumenata": boolean  // Возврат документов
}
```

### B2BManifestIn
```json
{
  "TipOperacije": int,       // Тип операции
  "ManifestId": string,      // ID манифеста
  "DatumManifesta": date,    // Дата манифеста
  "Posiljke": [              // Массив отправлений
    {
      "BrojPosiljke": string,
      "Posiljalac": Posiljalac,
      "Primalac": Primalac,
      "Posiljka": Posiljka,
      "DodatneUsluge": DodatneUsluge
    }
  ]
}
```

### B2BManifestOut
```json
{
  "OK": boolean,
  "Poruka": string,
  "ManifestId": string,      // ID манифеста
  "BrojPosiljaka": int,      // Количество отправлений
  "Status": string,          // Статус манифеста
  "DatumKreiranja": datetime,// Дата создания
  "DatumZatvaranja": datetime,// Дата закрытия
  "Posiljke": [              // Список отправлений
    {
      "BrojPosiljke": string,
      "Barkod": string,
      "Status": string,
      "Cena": float
    }
  ]
}
```

### HandlingIn (для операций 58, 71, 85)
```json
{
  "TransakcijaId": int,        // 58, 71 или 85
  "Klijent": Klijent,
  "Posiljalac": Posiljalac,
  "Primalac": Primalac,
  "Posiljka": Posiljka,
  "DodatneUsluge": DodatneUsluge,
  "NacinPlacanja": int,        // Способ оплаты
  "VremePrihvatanja": datetime,// Время приема
  "VremeIsporuke": datetime    // Планируемое время доставки
}
```

### HandlingOut
```json
{
  "OK": boolean,
  "Poruka": string,
  "BrojPosiljke": string,      // Присвоенный номер
  "Barkod": string,             // Штрих-код
  "Cena": float,                // Стоимость
  "CenaBezPDV": float,          // Стоимость без НДС
  "PDV": float,                 // НДС
  "VremePrihvatanja": datetime, // Время приема
  "PlaniranoVremeIsporuke": datetime // Планируемое время доставки
}
```

## Типы данных

### Форматы даты и времени
- **datetime**: ISO 8601 формат (YYYY-MM-DDTHH:mm:ss)
- **date**: ISO 8601 формат (YYYY-MM-DD)

### Числовые типы
- **int**: Целое число
- **float**: Число с плавающей точкой

### Строковые типы
- **string**: UTF-8 строка

### Логические типы
- **boolean**: true/false

## Ограничения

- Минимальная длина для поиска названий: 3 символа
- Максимальное количество отправлений для группового отслеживания: 100
- Максимальная длина референтного номера: 50 символов
- Максимальная длина примечания: 200 символов
- Формат номера телефона: международный формат (+381...)