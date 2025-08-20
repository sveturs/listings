# BEX Express Integration Documentation

## Обзор

Интеграция с курьерской службой BEX Express для доставки товаров по Сербии.

## Конфигурация

### Переменные окружения

```bash
# BEX Express API
BEX_AUTH_TOKEN=d50261-18wo-8539-ee5a-67uu3tu79
BEX_CLIENT_ID=326166
BEX_API_URL=https://api.bex.rs:62502
```

## База данных

### Основные таблицы

- `bex_settings` - Настройки интеграции
- `bex_municipalities` - Справочник муниципалитетов
- `bex_places` - Справочник населенных пунктов  
- `bex_streets` - Справочник улиц
- `bex_shipments` - Отправления
- `bex_tracking_events` - События отслеживания
- `bex_rates` - Тарифы доставки

### Импорт справочников

```bash
# Импорт данных из Excel файлов BEX
go run backend/scripts/import_bex_data.go
```

Импортируется:
- 45 муниципалитетов
- 1950+ населенных пунктов
- 65500+ улиц Сербии

## API Endpoints

### Создание отправления

```bash
POST /api/v1/bex/shipments

{
  "order_id": 123,
  "recipient_name": "Петар Петровић",
  "recipient_address": "Булевар ослобођења 100",
  "recipient_city": "Нови Сад",
  "recipient_postal_code": "21000",
  "recipient_phone": "+381 21 123456",
  "recipient_email": "petar@example.com",
  "weight_kg": 2.5,
  "total_packages": 1,
  "cod_amount": 3500,
  "insurance_amount": 5000,
  "personal_delivery": false,
  "notes": "Позвонить за 30 минут до доставки"
}
```

### Получение статуса

```bash
GET /api/v1/bex/shipments/{id}/status
```

### Получение этикетки

```bash
GET /api/v1/bex/shipments/{id}/label?size=4
```

Параметры:
- `size=4` - A4 формат
- `size=6` - A6 формат

### Отмена отправления

```bash
DELETE /api/v1/bex/shipments/{id}
```

### Расчет стоимости

```bash
POST /api/v1/bex/calculate-rate

{
  "weight_kg": 2.5,
  "shipment_category": 0,
  "cod_amount": 3500,
  "insurance_amount": 5000
}
```

### Поиск адреса

```bash
POST /api/v1/bex/search-address

{
  "query": "Булевар",
  "city": "Нови Сад",
  "limit": 10
}
```

### Отслеживание посылки

```bash
GET /api/v1/bex/track/{tracking_number}
```

## Статусы отправлений

| Код | Статус | Описание |
|-----|--------|----------|
| 0 | NotSentYet | Не отправлено |
| 1 | InTransit | В пути |
| 2 | Delivered | Доставлено |
| 3 | ReturnedToSender | Возвращено отправителю |
| 4 | PickedUp | Забрано курьером |
| 5 | Deleted | Отменено |

## Категории отправлений

| ID | Категория | Вес | Базовая цена |
|----|-----------|-----|--------------|
| 1 | Документы | до 0.5кг | 250 RSD |
| 2 | Посылка | до 1кг | 350 RSD |
| 3 | Посылка | до 2кг | 450 RSD |
| 31 | Посылка за кг | 2кг+ | 450 + 100/кг |
| 32 | Паллета за кг | любой | 1000 + 50/кг |

## Содержимое отправления

| ID | Тип |
|----|-----|
| 1 | Документы |
| 2 | Товары |
| 3 | Смешанное |

## Расчет стоимости

### Базовая стоимость
- Определяется по категории и весу

### Дополнительные услуги
- Наложенный платеж (COD): 150 RSD + 1% от суммы
- Страхование: 0.5% от застрахованной суммы
- Личная доставка: бесплатно

## Интеграция с заказами

При создании заказа можно выбрать BEX как способ доставки:

```sql
-- Marketplace orders
UPDATE marketplace_orders 
SET delivery_provider = 'bex', 
    delivery_metadata = '{"tracking_number": "123456"}'
WHERE id = 123;

-- Storefront orders  
UPDATE storefront_orders
SET delivery_provider = 'bex',
    delivery_metadata = '{"tracking_number": "123456"}'
WHERE id = 456;
```

## Webhook для обновления статусов

```bash
POST /api/v1/bex/webhook/status

{
  "shipment_id": 123,
  "status": 2,
  "status_text": "Delivered",
  "timestamp": "2025-01-19T10:30:00Z"
}
```

## Тестирование

### Создание тестового отправления

```bash
# Получаем JWT токен
TOKEN=$(cd backend && go run scripts/create_test_jwt.go)

# Создаем отправление
curl -X POST http://localhost:3000/api/v1/bex/shipments \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "recipient_name": "Тест Тестовић",
    "recipient_address": "Главна 1",
    "recipient_city": "Београд",
    "recipient_postal_code": "11000",
    "recipient_phone": "+381 11 1234567",
    "weight_kg": 1,
    "total_packages": 1
  }'
```

## Примечания

1. **Адреса**: Используйте поиск адресов для валидации и получения корректных ID улиц
2. **Вес**: Автоматически определяется категория по весу если не указана
3. **COD**: При наложенном платеже изменяется тип оплаты на "получатель платит"
4. **Этикетки**: Сохраняются в base64 в БД для быстрого доступа
5. **Статусы**: Обновляются автоматически через webhook или вручную через API

## Безопасность

- Все endpoints требуют JWT авторизацию
- Credentials хранятся в таблице `bex_settings`
- Поддержка тестового режима для разработки