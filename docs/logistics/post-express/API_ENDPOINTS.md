# WSP WebAPI Endpoints

Все запросы отправляются методом **POST** на единый endpoint:
```
https://onlinepostexpress.rs/WSPWebApi/api/app/transakcija
```

## 1. Адресная информация

### GetNaselje - Поиск населенных пунктов
**TransakcijaId**: 3

Поиск населенных пунктов по части названия или почтовому индексу.

**Входные параметры (NaseljeIn)**:
- `Naziv` (string) - часть названия населенного пункта (мин. 3 символа)
- `Ptt` (string) - почтовый индекс (опционально)

**Выходные параметры (NaseljeOut)**:
- `OK` (boolean) - успешность операции
- `Poruka` (string) - сообщение об ошибке
- `Naselja` (array) - массив найденных населенных пунктов
  - `Sifra` (int) - код населенного пункта
  - `Naziv` (string) - название
  - `Ptt` (string) - почтовый индекс
  - `Opstina` (string) - община

### GetUlica - Поиск улиц
**TransakcijaId**: 4

Поиск улиц в населенном пункте.

**Входные параметры (UlicaIn)**:
- `SifraMesta` (int) - код населенного пункта
- `Naziv` (string) - часть названия улицы (мин. 3 символа)

**Выходные параметры (UlicaOut)**:
- `OK` (boolean) - успешность операции
- `Poruka` (string) - сообщение об ошибке
- `Ulice` (array) - массив найденных улиц
  - `Sifra` (int) - код улицы
  - `Naziv` (string) - название улицы
  - `SifraMesta` (int) - код населенного пункта
  - `NazivMesta` (string) - название населенного пункта
  - `PttMesta` (string) - почтовый индекс

### ProveraAdrese (AddressCheck) - Проверка адреса
**TransakcijaId**: 6

Проверка корректности и полноты адреса.

**Входные параметры (ProveraAdreseIn)**:
- `SifraMesta` (int) - код населенного пункта
- `SifraUlice` (int) - код улицы
- `Broj` (string) - номер дома
- `Ulaz` (string) - подъезд (опционально)
- `Sprat` (string) - этаж (опционально)
- `Stan` (string) - квартира (опционально)
- `Ptt` (string) - почтовый индекс

**Выходные параметры (ProveraAdreseOut)**:
- `OK` (boolean) - успешность операции
- `Poruka` (string) - сообщение об ошибке
- `Adresa` (object) - проверенная адреса
  - `DostupanKurir` (boolean) - доступность курьерской службы
  - `DostupanPostExpress` (boolean) - доступность PostExpress
  - `TipKucnogBroja` (string) - тип номера дома

### ProveraDostupnostiUsluge - Проверка доступности услуг
**TransakcijaId**: 23

Проверка доступности конкретной услуги по адресу.

**Входные параметры (ProveraDostupnostiUslugeIn)**:
- `SifraMesta` (int) - код населенного пункта
- `SifraUlice` (int) - код улицы (опционально)
- `KucniBroj` (string) - номер дома (опционально)
- `SifraUsluge` (int) - код услуги

**Выходные параметры (ProveraDostupnostiUslugeOut)**:
- `OK` (boolean) - успешность операции
- `Poruka` (string) - сообщение об ошибке
- `Dostupnost` (boolean) - доступность услуги

## 2. Расчет стоимости

### PostarinaPosiljke - Расчет почтовых расходов
**TransakcijaId**: 5

Расчет стоимости отправки.

**Входные параметры (PostarinaIn)**:
- `SifraUsluge` (int) - код услуги
- `Masa` (float) - вес в граммах
- `Vrednost` (float) - объявленная стоимость (для ценных отправлений)
- `Otkupnina` (float) - сумма наложенного платежа
- `LicnoPreuzimanje` (boolean) - личное вручение
- `PovratnicaPotpis` (boolean) - уведомление о вручении
- `SMS` (boolean) - SMS уведомление
- `Potvrda` (boolean) - подтверждение
- `SifraMestaPrimaoca` (int) - код населенного пункта получателя
- `BrojPosiljaka` (int) - количество отправлений

**Выходные параметры (PostarinaOut)**:
- `OK` (boolean) - успешность операции
- `Poruka` (string) - сообщение об ошибке
- `Postarina` (object) - расчет стоимости
  - `UkupnaCena` (float) - общая стоимость
  - `CenaBezPDV` (float) - стоимость без НДС
  - `PDV` (float) - НДС
  - `OsnovnaCena` (float) - базовая стоимость
  - `DodatneUsluge` (float) - дополнительные услуги

## 3. Отслеживание отправлений

### TTKretanjaUsluge - Отслеживание отдельного отправления
**TransakcijaId**: 63

Получение информации о движении отправления.

**Входные параметры (TTKretanjeIn)**:
- `BrojPosiljke` (string) - номер отправления (штрих-код)
- `Jezik` (int) - язык ответа (1-сербский кириллица, 2-сербский латиница, 3-английский)

**Выходные параметры (TTKretanjeOut)**:
- `OK` (boolean) - успешность операции
- `Poruka` (string) - сообщение об ошибке
- `Kretanja` (array) - массив событий движения
  - `DatumVreme` (datetime) - дата и время события
  - `Mesto` (string) - место события
  - `Status` (string) - статус
  - `Opis` (string) - описание события

### TTPosiljkeStatusi - Групповое отслеживание
**TransakcijaId**: 64

Получение статусов для группы отправлений.

**Входные параметры (TTStatusiPosiljakaIn)**:
- `BrojeviPosiljaka` (array of string) - массив номеров отправлений (до 100)
- `Jezik` (int) - язык ответа

**Выходные параметры (TTStatusiPosiljakaOut)**:
- `OK` (boolean) - успешность операции
- `Poruka` (string) - сообщение об ошибке
- `Statusi` (array) - массив статусов
  - `BrojPosiljke` (string) - номер отправления
  - `Status` (string) - текущий статус
  - `DatumStatusa` (datetime) - дата последнего статуса
  - `MestoStatusa` (string) - место последнего статуса
  - `OpisStatusa` (string) - описание статуса
  - `Isporuceno` (boolean) - доставлено ли отправление

## 4. B2B операции

### B2BManifest - Управление манифестами
**TransakcijaId**: Зависит от операции

Создание и управление манифестами для бизнес-отправлений.

#### Операции:
- **Создание манифеста** (TransakcijaId: варьируется)
- **Добавление отправлений** (TransakcijaId: варьируется)
- **Закрытие манифеста** (TransakcijaId: варьируется)
- **Получение манифеста** (TransakcijaId: варьируется)

**Входные параметры (B2BManifestIn)**:
- `TipOperacije` (int) - тип операции
- `ManifestId` (string) - идентификатор манифеста
- `Posiljke` (array) - массив отправлений

**Выходные параметры (B2BManifestOut)**:
- `OK` (boolean) - успешность операции
- `Poruka` (string) - сообщение об ошибке
- `ManifestId` (string) - идентификатор манифеста
- `BrojPosiljaka` (int) - количество отправлений
- `Status` (string) - статус манифеста

### B2B Handling Operations

#### Handling 58 - Express Today
**TransakcijaId**: 58

Быстрая доставка в тот же день в пределах города.

#### Handling 71 - City Express
**TransakcijaId**: 71

Экспресс-доставка в пределах города (до 90 минут).

#### Handling 85 - Dedicated Delivery
**TransakcijaId**: 85

Выделенная доставка с персональным курьером.

**Общие входные параметры для Handling**:
- `Posiljalac` (object) - данные отправителя
- `Primalac` (object) - данные получателя
- `Posiljka` (object) - данные отправления
- `DodatneUsluge` (object) - дополнительные услуги
- `NacinPlacanja` (int) - способ оплаты

**Общие выходные параметры для Handling**:
- `OK` (boolean) - успешность операции
- `Poruka` (string) - сообщение об ошибке
- `BrojPosiljke` (string) - присвоенный номер отправления
- `Barkod` (string) - штрих-код
- `Cena` (float) - стоимость услуги

## Общие параметры для всех запросов

### Класс Klijent (обязателен для всех запросов)
- `Username` (string) - имя пользователя
- `Password` (string) - пароль

### Базовые поля запроса (TransakcijaIn)
- `TransakcijaId` (int) - идентификатор транзакции
- `DatumVremePosiljke` (datetime) - дата и время отправки запроса
- `Klijent` (object) - данные клиента для аутентификации

### Базовые поля ответа (TransakcijaOut)
- `OK` (boolean) - успешность выполнения
- `Poruka` (string) - сообщение (обычно об ошибке)
- `TransakcijaId` (int) - идентификатор транзакции
- `DatumVremePrijema` (datetime) - дата и время приема запроса

## Коды ошибок

- `AUTH_ERROR` - ошибка аутентификации
- `INVALID_PARAMS` - некорректные параметры запроса
- `NOT_FOUND` - запрашиваемые данные не найдены
- `SERVICE_UNAVAILABLE` - сервис временно недоступен
- `QUOTA_EXCEEDED` - превышен лимит запросов