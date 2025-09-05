# Post Express (Pošta Srbije) WSP WebAPI Documentation

## Общее описание

WSP WebAPI је REST Web API сервис Поште Србије који омогућава корисницима коришћење услуга Поште Србије путем сопствених апликација. Сервис пружа могућност размене података са информационим системом Поште Србије у реалном времену.

## Основные характеристики

- **Протокол**: REST API
- **Формат данных**: JSON
- **Метод**: POST
- **Основной endpoint**: `https://onlinepostexpress.rs/WSPWebApi/api/app/transakcija`
- **Аутентификация**: Username/Password через класс Klijent

## Регистрация пользователя

За коришћење сервиса потребна је регистрација:
1. Приступите веб страници https://onlinepostexpress.rs/registracija
2. Попуните образац за регистрацију
3. Након потврде регистрације добијате приступне параметре (Username и Password)

## Доступные функции

### 1. Адресная информация
- **GetNaselje** - Поиск населенных пунктов
- **GetUlica** - Поиск улиц в населенном пункте
- **ProveraAdrese** - Проверка корректности адреса
- **ProveraDostupnostiUsluge** - Проверка доступности услуг по адресу

### 2. Расчет стоимости
- **PostarinaPosiljke** - Расчет почтовых расходов

### 3. Отслеживание отправлений
- **TTKretanjaUsluge** - Отслеживание движения отдельного отправления
- **TTPosiljkeStatusi** - Групповое отслеживание статусов отправлений

### 4. B2B операции
- **B2BManifest** - Создание и управление манифестами для бизнес-отправлений
- Передача отправлений (Handling 58, 71, 85)

## Структура запроса

Все запросы отправляются методом POST на endpoint транзакции:

```json
{
  "TransakcijaId": 1,  // Идентификатор транзакции (зависит от операции)
  "DatumVremePosiljke": "2024-01-25T10:00:00",
  "Klijent": {
    "Username": "your_username",
    "Password": "your_password"
  },
  // Дополнительные параметры в зависимости от транзакции
}
```

## Структура ответа

```json
{
  "OK": true,  // Успешность операции
  "Poruka": "Сообщение",
  "TransakcijaId": 1,
  "DatumVremePrijema": "2024-01-25T10:00:01",
  // Результаты операции
}
```

## Важные замечания

1. **Обязательные поля**: Все поля класса Klijent обязательны для заполнения
2. **Формат даты**: ISO 8601 (YYYY-MM-DDTHH:mm:ss)
3. **Кодировка**: UTF-8
4. **Регистрация**: Обязательна для получения доступа к API

## Документация

- [API Endpoints](./API_ENDPOINTS.md) - Полное описание всех эндпоинтов
- [Data Structures](./DATA_STRUCTURES.md) - Структуры данных и классы
- [Examples](./EXAMPLES.md) - Примеры запросов и ответов
- [Integration Guide](./INTEGRATION_GUIDE.md) - Руководство по интеграции
- [Enums and Codes](./ENUMS_AND_CODES.md) - Справочники кодов и перечислений

## Поддержка

- **Телефон**: 0700 100 300
- **Email**: wsp@posta.rs
- **Сайт**: https://www.postexpress.rs

## Версии документации

Эта документация основана на:
- WSP Web Api Address information v1.0
- WSP Web Api Exchange of data v1.0
- WSP Web Api Shipment tracking v1.0