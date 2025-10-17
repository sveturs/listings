# Post Express WSP API - Critical Issues Report

**Дата составления:** 2025-10-15
**Версия API:** Post Express WSP B2B (0.2.4)
**Среда:** Production тестирование
**Контакт Post Express:** Nikola Dmitrašinović (nikola.dmitrasinovic@posta.rs)

---

## Executive Summary

В процессе интеграции Post Express WSP API были выявлены три критические проблемы на стороне Post Express, которые блокируют использование следующих транзакций:

1. **TX6 (ValidateAddress)** - Oracle database connection loss
2. **TX9 (CheckServiceAvailability)** - Missing address data error (даже после добавления всех требуемых полей)
3. **TX11 (CalculatePostage)** - Database schema error (missing column)

Все проблемы возникают на стороне Post Express API и не могут быть исправлены со стороны интегратора.

---

## 1. TX6 ValidateAddress - Oracle Connection Error

### Описание проблемы

Транзакция TX6 (валидация адреса) возвращает ошибку потери соединения с Oracle базой данных.

### Технические детали

**Endpoint:** `POST https://wsp.posta.rs/PostExpressWS3/Service.asmx`
**Transaction:** TX 6
**Дата/время:** 2025-10-15, несколько попыток между 14:00-15:00 UTC+1
**HTTP Status:** 200 OK (ошибка в теле ответа)

### Request структура

```json
{
  "TipAdrese": 0,
  "IdRukovanje": 71,
  "IdNaselje": 100001,
  "IdUlica": 1186,
  "BrojPodbroj": "2",
  "PostanskiBroj": "11000"
}
```

### Response

```json
{
  "Rezultat": 1,
  "Poruka": "Dogodila se greška prilikom izvršavanja upita!",
  "StrRezultat": "ODP error: ORA-03135: connection lost contact\nProcess ID: 19464570\nSession ID: 80 Serial number: 28891"
}
```

### Интерпретация

**ORA-03135** - критическая ошибка Oracle:
- Потеря соединения между application server и database server
- Process ID и Session ID указывают на конкретную сессию, которая была прервана
- Возможные причины:
  - Network timeout между app server и DB
  - Database server crash/restart
  - Firewall или network issues
  - Database resource exhaustion (connections, memory)

### Тестовые данные

Использовались корректные данные из TX3/TX4:
- **IdNaselje:** 100001 (Београд-Врачар, подтверждён TX3)
- **IdUlica:** 1186 (Таковска, подтверждена TX4)
- **PostanskiBroj:** 11000 (корректный индекс Белграда)
- **BrojPodbroj:** "2" (дом №2)
- **TipAdrese:** 0 (стандартный тип)
- **IdRukovanje:** 71 (корректный service ID)

### Воспроизводимость

**100% воспроизводимо** - ошибка возникает при каждом вызове TX6 с различными параметрами.

### История исправлений (наша сторона)

1. **Исправление 1:** Переименование поля `Broj` → `BrojPodbroj` (требование API)
2. **Исправление 2:** Добавление обязательных полей `TipAdrese` и `IdRukovanje`
3. **Исправление 3:** Добавление опционального поля `Datum`

После всех исправлений структура запроса полностью соответствует спецификации API, но ошибка Oracle сохраняется.

---

## 2. TX9 CheckServiceAvailability - Missing Address Data Error

### Описание проблемы

Транзакция TX9 (проверка доступности услуги) возвращает ошибку "Podaci adrese nisu prosleđeni!" (Address data not provided) даже после добавления всех требуемых полей адреса.

### Технические детали

**Endpoint:** `POST https://wsp.posta.rs/PostExpressWS3/Service.asmx`
**Transaction:** TX 9
**Дата/время:** 2025-10-15, несколько попыток между 13:00-14:00 UTC+1
**HTTP Status:** 200 OK (ошибка в теле ответа)

### Request структура (после всех исправлений)

```json
{
  "TipAdrese": 0,
  "IdRukovanje": 71,
  "IdNaseljeOdlaska": 100001,
  "IdNaseljeDolaska": 100039,
  "PostanskiBrojOdlaska": "11000",
  "PostanskiBrojDolaska": "21000"
}
```

### Response

```json
{
  "Rezultat": 3,
  "Poruka": "Podaci adrese nisu prosleđeni!",
  "StrRezultat": "{\"Poruka\":\"Podaci adrese nisu prosleđeni! \",\"PorukaKorisnik\":\"Podaci adrese nisu prosleđeni! \",\"Info\":null}"
}
```

### Интерпретация

**"Podaci adrese nisu prosleđeni!"** (Address data not provided):
- API возвращает эту ошибку даже после добавления всех ожидаемых адресных полей
- Были добавлены: `TipAdrese`, `IdNaseljeOdlaska`, `IdNaseljeDolaska`, `PostanskiBrojOdlaska`, `PostanskiBrojDolaska`, `Datum`
- Возможные причины:
  - Неправильная документация API (требуются другие поля)
  - TX9 не реализована или отключена на тестовом окружении
  - Внутренняя ошибка валидации на стороне API
  - Требуются дополнительные поля, не указанные в документации

### Тестовые данные

Использовались корректные данные из TX3:
- **IdNaseljeOdlaska:** 100001 (Београд-Врачар, подтверждён TX3)
- **IdNaseljeDolaska:** 100039 (Нови Сад, подтверждён TX3)
- **PostanskiBrojOdlaska:** 11000 (Белград)
- **PostanskiBrojDolaska:** 21000 (Нови Сад)
- **TipAdrese:** 0 (стандартный тип)
- **IdRukovanje:** 71 (стандартная услуга)

### Воспроизводимость

**100% воспроизводимо** - ошибка возникает при каждом вызове TX9 с различными комбинациями полей.

### История исправлений (наша сторона)

1. **Попытка 1:** Добавили `TipAdrese` - не помогло
2. **Попытка 2:** Добавили `IdNaseljeOdlaska` и `IdNaseljeDolaska` - не помогло
3. **Попытка 3:** Добавили `Datum` - не помогло

После всех исправлений структура запроса содержит все поля из документации, но ошибка сохраняется.

---

## 3. TX11 CalculatePostage - Database Schema Error

### Описание проблемы

Транзакция TX11 (расчёт стоимости доставки) возвращает ошибку отсутствующей колонки в таблице базы данных.

### Технические детали

**Endpoint:** `POST https://wsp.posta.rs/PostExpressWS3/Service.asmx`
**Transaction:** TX 11
**Дата/время:** 2025-10-15, несколько попыток между 14:00-15:00 UTC+1
**HTTP Status:** 200 OK (ошибка в теле ответа)

### Request структура

```json
{
  "IdRukovanje": 71,
  "IdZemlja": 0,
  "PostanskiBrojOdlaska": "11000",
  "PostanskiBrojDolaska": "21000",
  "Masa": 500,
  "Otkupnina": 5000,
  "Vrednost": 5000,
  "PosebneUsluge": "PNA"
}
```

### Response

```json
{
  "Rezultat": 1,
  "Poruka": "Dogodila se greška prilikom izvršavanja upita!",
  "StrRezultat": "Column 'PREVOD_' does not belong to table Prevodi"
}
```

### Интерпретация

**"Column 'PREVOD_' does not belong to table Prevodi"** - ошибка schema:
- В коде API происходит обращение к несуществующей колонке `PREVOD_` в таблице `Prevodi`
- Возможные причины:
  - Несоответствие между кодом и актуальной schema БД
  - Неполная миграция базы данных после изменений
  - Typo в SQL запросе (пустое имя колонки после префикса)
  - Отсутствие необходимой колонки после рефакторинга

### Особенность

Название колонки `PREVOD_` с подчёркиванием в конце предполагает, что это либо:
1. Префикс для динамически генерируемого имени колонки (например, `PREVOD_EN`, `PREVOD_SR`)
2. Ошибка в коде, где имя колонки не было полностью сформировано

Учитывая предыдущую выявленную проблему с языком (TX11 возвращал данные только на сербском даже при запросе английского), это может быть связано с системой переводов/локализации.

### Тестовые данные

Использовались корректные данные:
- **IdRukovanje:** 71 (стандартная услуга)
- **IdZemlja:** 0 (внутренние отправления по Сербии)
- **PostanskiBrojOdlaska:** 11000 (Белград)
- **PostanskiBrojDolaska:** 21000 (Нови Сад)
- **Masa:** 500г (корректный вес)
- **Otkupnina:** 5000 para (наложенный платёж)
- **Vrednost:** 5000 para (объявленная ценность)
- **PosebneUsluge:** "PNA" (уведомление о вручении)

### Воспроизводимость

**100% воспроизводимо** - ошибка возникает при каждом вызове TX11 с различными параметрами.

### История исправлений (наша сторона)

1. **Исправление:** Добавление обязательного поля `IdZemlja: 0` для внутренних отправлений

После исправления структура запроса полностью соответствует спецификации API, но database schema error сохраняется.

---

## 3. Работающие транзакции (для сравнения)

Для контекста - следующие транзакции работают корректно:

### TX3 - GetSettlements ✅

```bash
curl -X POST "http://localhost:3000/api/v1/postexpress/test/tx3-get-settlements" \
  -H "Content-Type: application/json" \
  -d '{"query":"Beograd"}'
```

**Результат:** Успешно находит 2 населённых пункта (Београд, Београд-Врачар)

### TX4 - GetStreets ✅

```bash
curl -X POST "http://localhost:3000/api/v1/postexpress/test/tx4-get-streets" \
  -H "Content-Type: application/json" \
  -d '{"settlement_id":100001,"query":"Takovska"}'
```

**Результат:** Успешно находит улицу Таковска (ID: 1186)

---

## 4. Воздействие на интеграцию

### Критичность

**HIGH** - все три транзакции критически важны для e-commerce интеграции:

1. **TX6** - необходима для валидации адресов доставки перед созданием заказа
2. **TX9** - необходима для проверки доступности услуги на маршруте доставки
3. **TX11** - необходима для расчёта стоимости доставки и отображения цен пользователю

### Обходные пути

**Частично доступны**:
- TX6 - единственный способ валидировать адреса (обходных путей нет)
- TX9 - можно предполагать доступность на основе почтовых кодов (неточно)
- TX11 - можно использовать фиксированные тарифы из прайс-листа Post Express (требует ручного обновления)

### Блокировка разработки

Следующие этапы разработки заблокированы:
- [ ] Создание frontend страниц для TX6, TX9, TX11
- [ ] Массовое тестирование (30+ тестовых отправлений)
- [ ] Интеграция в production checkout flow с динамическим расчётом стоимости
- [ ] End-to-end тестирование полного процесса создания заказа

---

## 5. Рекомендации для Post Express

### TX6 - Immediate Actions Required

1. **Проверить Oracle database server:**
   - Connection pool settings
   - Network connectivity между app server и DB
   - Database server logs для Process ID 19464570, Session 80
   - Resource utilization (CPU, memory, connections)

2. **Проверить application server:**
   - Connection timeout settings
   - Firewall rules
   - Network latency к database server

3. **Временное решение:**
   - Увеличить connection timeout
   - Настроить connection retry logic
   - Проверить connection pooling

### TX9 - Address Data Validation Fix Required

1. **Проверить требования к адресным полям:**
   - Документировать точный список обязательных полей для TX9
   - Сравнить с текущей реализацией
   - Обновить документацию API

2. **Проверить логику валидации:**
   - Определить, почему API считает, что адресные данные отсутствуют
   - Проверить порядок полей в запросе
   - Проверить форматирование данных (типы, регистр)

3. **Альтернативные решения:**
   - Предоставить работающий пример запроса TX9
   - Если TX9 не используется на production, документировать это
   - Рассмотреть альтернативные способы проверки доступности услуг

### TX11 - Database Schema Fix Required

1. **Проверить SQL код:**
   - Найти запросы к таблице `Prevodi`
   - Идентифицировать использование колонки `PREVOD_`
   - Проверить logic формирования динамических имён колонок

2. **Проверить database schema:**
   - Сравнить актуальную schema с кодом
   - Проверить наличие всех необходимых колонок в таблице `Prevodi`
   - Выполнить недостающие миграции если требуется

3. **Связь с языковой проблемой:**
   - Возможно это часть системы переводов (Prevodi = Translations)
   - Проверить logic выбора языка и формирования имён колонок типа `PREVOD_EN`, `PREVOD_SR` и т.д.

---

## 6. Логи и примеры для воспроизведения

### Full curl commands для воспроизведения

```bash
# TX6 - Oracle error
curl -s -X POST "http://localhost:3000/api/v1/postexpress/test/tx6-validate-address" \
  -H "Content-Type: application/json" \
  -d '{
    "settlement_id": 100001,
    "street_id": 1186,
    "house_number": "2",
    "postal_code": "11000"
  }' | jq '.'

# TX9 - Missing address data error
curl -s -X POST "http://localhost:3000/api/v1/postexpress/test/tx9-service-availability" \
  -H "Content-Type: application/json" \
  -d '{
    "service_id": 71,
    "from_postal_code": "11000",
    "to_postal_code": "21000"
  }' | jq '.'

# TX11 - Schema error
curl -s -X POST "http://localhost:3000/api/v1/postexpress/test/tx11-calculate-postage" \
  -H "Content-Type: application/json" \
  -d '{
    "service_id": 71,
    "from_postal_code": "11000",
    "to_postal_code": "21000",
    "weight": 500,
    "value": 5000,
    "cod_amount": 5000,
    "services": "PNA"
  }' | jq '.'
```

### Backend logs excerpt

```
[2025-10-15 15:23:45] POST /api/v1/postexpress/test/tx6-validate-address
Request: {TipAdrese:0 IdRukovanje:71 IdNaselje:100001 IdUlica:1186 BrojPodbroj:2 PostanskiBroj:11000}
Response: {Rezultat:1 Poruka:"Dogodila se greška prilikom izvršavanja upita!" StrRezultat:"ODP error: ORA-03135: connection lost contact..."}

[2025-10-15 15:24:12] POST /api/v1/postexpress/test/tx11-calculate-postage
Request: {IdRukovanje:71 IdZemlja:0 PostanskiBrojOdlaska:11000 PostanskiBrojDolaska:21000 Masa:500...}
Response: {Rezultat:1 Poruka:"Dogodila se greška prilikom izvršavanja upita!" StrRezultat:"Column 'PREVOD_' does not belong to table Prevodi"}
```

---

## 7. Следующие шаги

### Со стороны Post Express (ТРЕБУЕТСЯ)

1. **Приоритет 1:** Исправить TX11 database schema error
2. **Приоритет 1:** Решить TX6 Oracle connection issue
3. **Приоритет 2:** Тестирование обеих транзакций после исправлений
4. **Приоритет 3:** Предоставить feedback о timeline исправлений

### Со стороны интегратора (готово к выполнению после исправлений)

1. ✅ Структуры запросов TX6 и TX11 полностью соответствуют спецификации
2. ✅ Все обязательные поля добавлены
3. ⏳ Ожидание исправлений от Post Express
4. ⏳ После исправлений - продолжение интеграции (frontend, массовое тестирование)

---

## Приложение: Полная структура запросов

### TX6 ValidateAddress

```go
type AddressValidationRequest struct {
    TipAdrese      int    `json:"TipAdrese"`      // Тип адреса (0, 1, 2)
    IdRukovanje    int    `json:"IdRukovanje"`    // ID услуги доставки
    IdNaselje      int    `json:"IdNaselje"`      // ID населённого пункта
    IdUlica        int    `json:"IdUlica,omitempty"` // ID улицы (опционально)
    BrojPodbroj    string `json:"BrojPodbroj"`    // Номер дома (например, "2" или "2a")
    PostanskiBroj  string `json:"PostanskiBroj"`  // Почтовый индекс
    Datum          string `json:"Datum,omitempty"` // Дата (опционально, формат: YYYY-MM-DD)
}
```

### TX11 CalculatePostage

```go
type PostageCalculationRequest struct {
    IdRukovanje            int    `json:"IdRukovanje"`            // ID услуги доставки
    IdZemlja               int    `json:"IdZemlja"`               // ID страны (0 = внутренние)
    PostanskiBrojOdlaska   string `json:"PostanskiBrojOdlaska"`   // Почтовый индекс отправления
    PostanskiBrojDolaska   string `json:"PostanskiBrojDolaska"`   // Почтовый индекс прибытия
    Masa                   int    `json:"Masa"`                   // Вес в граммах
    Otkupnina              int    `json:"Otkupnina,omitempty"`    // COD в para
    Vrednost               int    `json:"Vrednost,omitempty"`     // Объявленная ценность в para
    PosebneUsluge          string `json:"PosebneUsluge,omitempty"` // Дополнительные услуги
}
```

---

**Контакты для вопросов:**
- **Email:** nikola.dmitrasinovic@posta.rs
- **Тема письма:** "WSP API TX6 & TX11 Critical Issues - Integration Blocked"
