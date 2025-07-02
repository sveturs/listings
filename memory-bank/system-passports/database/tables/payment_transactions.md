# Паспорт таблицы: payment_transactions

## Назначение
Хранение всех платежных транзакций через внешние платежные шлюзы (AllSecure, Stripe и др.). Основная таблица для процессинга платежей на маркетплейсе.

## Структура таблицы

```sql
CREATE TABLE payment_transactions (
    id BIGSERIAL PRIMARY KEY,
    gateway_id INT REFERENCES payment_gateways(id),
    user_id INT REFERENCES users(id),
    
    -- Ссылка на покупку
    listing_id INT REFERENCES marketplace_listings(id),
    order_reference VARCHAR(255) UNIQUE NOT NULL,
    
    -- AllSecure данные
    gateway_transaction_id VARCHAR(255),
    gateway_reference_id VARCHAR(255),
    
    -- Финансовые данные
    amount DECIMAL(12,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'RSD',
    marketplace_commission DECIMAL(12,2),
    seller_amount DECIMAL(12,2),
    
    -- Статусы
    status VARCHAR(50) DEFAULT 'pending',
    gateway_status VARCHAR(50),
    
    -- Дополнительная информация
    payment_method VARCHAR(50),
    customer_email VARCHAR(255),
    description TEXT,
    
    -- Metadata
    gateway_response JSONB,
    error_details JSONB,
    
    -- Timestamps
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    authorized_at TIMESTAMP WITH TIME ZONE,
    captured_at TIMESTAMP WITH TIME ZONE,
    failed_at TIMESTAMP WITH TIME ZONE,
    
    -- Constraints
    CONSTRAINT payment_transactions_amount_positive CHECK (amount > 0),
    CONSTRAINT payment_transactions_status_valid CHECK (
        status IN ('pending', 'authorized', 'captured', 'failed', 'refunded', 'voided')
    )
);
```

## Поля таблицы

### Основные поля
- `id` - уникальный идентификатор (BIGSERIAL)
- `gateway_id` - используемый платежный шлюз (FK)
- `user_id` - покупатель (FK)
- `listing_id` - купленное объявление (FK)
- `order_reference` - уникальная ссылка на заказ (UNIQUE)

### Идентификаторы шлюза
- `gateway_transaction_id` - UUID транзакции в платежном шлюзе
- `gateway_reference_id` - Purchase ID в платежном шлюзе

### Финансовые данные
- `amount` - общая сумма платежа (CHECK > 0)
- `currency` - валюта (по умолчанию 'RSD')
- `marketplace_commission` - комиссия маркетплейса
- `seller_amount` - сумма продавцу после комиссии

### Статусы
- `status` - внутренний статус транзакции
- `gateway_status` - статус от платежного шлюза

### Информация о платеже
- `payment_method` - способ оплаты (card, bank_transfer и др.)
- `customer_email` - email покупателя
- `description` - описание платежа

### Метаданные
- `gateway_response` - полный ответ от шлюза (JSONB)
- `error_details` - детали ошибки при неудаче (JSONB)

### Временные метки
- `created_at` - создание транзакции
- `updated_at` - последнее обновление
- `authorized_at` - время авторизации
- `captured_at` - время списания средств
- `failed_at` - время неудачи

## Статусы транзакций

### Внутренние статусы (status)
- `pending` - ожидает обработки
- `authorized` - средства заблокированы
- `captured` - средства списаны
- `failed` - транзакция неудачна
- `refunded` - средства возвращены
- `voided` - транзакция отменена

### Статусы AllSecure (gateway_status)
- `CREATED` - создана
- `PENDING` - в обработке
- `SUCCESSFUL` - успешна
- `REDIRECT` - требуется 3DS
- `CANCELLED` - отменена
- `ERROR` - ошибка

## Индексы

1. **idx_payment_transactions_user_id** - по покупателю
2. **idx_payment_transactions_listing_id** - по объявлению
3. **idx_payment_transactions_status** - по статусу
4. **idx_payment_transactions_gateway_transaction_id** - по ID шлюза
5. **idx_payment_transactions_order_reference** - по ссылке заказа
6. **idx_payment_transactions_created_at** - по дате создания

## Связи с другими таблицами

### Прямые связи
- `gateway_id` → `payment_gateways.id` - платежный шлюз
- `user_id` → `users.id` - покупатель
- `listing_id` → `marketplace_listings.id` - объявление

### Обратные связи
- `escrow_payments.payment_transaction_id` - эскроу платежи

## Структура gateway_response

### AllSecure успешный ответ
```json
{
  "uuid": "e7f4b3a0-1234-5678-9abc-def012345678",
  "purchaseId": "20240115-123456",
  "status": "SUCCESSFUL",
  "result": {
    "code": "000.000.000",
    "description": "Transaction succeeded"
  },
  "returnData": {
    "cardHolder": "John Doe",
    "binCountry": "RS",
    "binBank": "RAIFFEISEN BANK"
  },
  "paymentMethod": {
    "type": "card",
    "brand": "VISA",
    "last4": "4242"
  }
}
```

### Структура error_details
```json
{
  "code": "PAYMENT_FAILED",
  "message": "Insufficient funds",
  "gateway_code": "100.100.101",
  "gateway_message": "Account or Bank Details Incorrect",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## Примеры использования

### Создание платежной транзакции
```sql
INSERT INTO payment_transactions (
    gateway_id, user_id, listing_id, order_reference,
    amount, currency, marketplace_commission, seller_amount,
    payment_method, customer_email, description
) VALUES (
    :gateway_id,
    :user_id,
    :listing_id,
    :order_reference,
    :amount,
    'RSD',
    :amount * 0.05, -- 5% комиссия
    :amount * 0.95,
    'card',
    :customer_email,
    'Payment for listing #' || :listing_id
) RETURNING id, order_reference;
```

### Обновление после ответа шлюза
```sql
UPDATE payment_transactions
SET 
    gateway_transaction_id = :gateway_uuid,
    gateway_reference_id = :purchase_id,
    status = CASE 
        WHEN :gateway_status = 'SUCCESSFUL' THEN 'captured'
        WHEN :gateway_status = 'PENDING' THEN 'authorized'
        ELSE 'failed'
    END,
    gateway_status = :gateway_status,
    gateway_response = :response_json,
    updated_at = NOW(),
    captured_at = CASE 
        WHEN :gateway_status = 'SUCCESSFUL' THEN NOW() 
        ELSE NULL 
    END,
    failed_at = CASE 
        WHEN :gateway_status IN ('ERROR', 'CANCELLED') THEN NOW() 
        ELSE NULL 
    END
WHERE id = :transaction_id;
```

### Получение транзакций пользователя
```sql
SELECT 
    pt.*,
    l.title as listing_title,
    pg.name as gateway_name
FROM payment_transactions pt
JOIN marketplace_listings l ON pt.listing_id = l.id
JOIN payment_gateways pg ON pt.gateway_id = pg.id
WHERE pt.user_id = :user_id
ORDER BY pt.created_at DESC
LIMIT 20;
```

### Статистика платежей
```sql
SELECT 
    DATE_TRUNC('day', created_at) as date,
    COUNT(*) as total_transactions,
    COUNT(*) FILTER (WHERE status = 'captured') as successful,
    COUNT(*) FILTER (WHERE status = 'failed') as failed,
    SUM(amount) FILTER (WHERE status = 'captured') as revenue,
    SUM(marketplace_commission) FILTER (WHERE status = 'captured') as commission
FROM payment_transactions
WHERE created_at >= CURRENT_DATE - INTERVAL '30 days'
GROUP BY DATE_TRUNC('day', created_at)
ORDER BY date DESC;
```

### Проверка дубликатов
```sql
-- Предотвращение двойного списания
SELECT COUNT(*) 
FROM payment_transactions
WHERE listing_id = :listing_id
  AND user_id = :user_id
  AND status IN ('captured', 'authorized')
  AND created_at >= NOW() - INTERVAL '5 minutes';
```

## Бизнес-правила

### Создание транзакций
1. **Уникальный order_reference** - генерируется как `{date}-{user_id}-{random}`
2. **Проверка доступности** - listing должен быть active
3. **Одна активная транзакция** - на пару user-listing

### Обработка платежей
1. **Двухэтапный платеж** - сначала authorize, потом capture
2. **Timeout авторизации** - 7 дней для capture
3. **Автоматический void** - если не captured вовремя

### Комиссии
1. **Расчет при создании** - фиксируется в момент платежа
2. **Прозрачность** - покупатель видит полную сумму
3. **Выплата продавцу** - seller_amount после успешной доставки

### Безопасность
1. **Идемпотентность** - повторные запросы не создают дубли
2. **Webhook верификация** - проверка подписи от шлюза
3. **Логирование** - все операции записываются

## API интеграция

### Endpoints
- `POST /api/v1/payments/create` - создание платежа
- `GET /api/v1/payments/{id}` - статус платежа
- `POST /api/v1/payments/{id}/capture` - подтверждение платежа
- `POST /api/v1/payments/{id}/refund` - возврат средств

### Webhook обработка
```sql
-- Обработка webhook от AllSecure
UPDATE payment_transactions
SET 
    status = 'captured',
    gateway_status = 'SUCCESSFUL',
    gateway_response = gateway_response || :webhook_data,
    captured_at = NOW()
WHERE gateway_transaction_id = :uuid
  AND status = 'authorized';
```

## Известные особенности

1. **BIGSERIAL для масштаба** - ожидается большой объем транзакций
2. **JSONB для гибкости** - разные шлюзы = разные форматы
3. **Множество timestamps** - для детального tracking
4. **CHECK constraints** - валидация на уровне БД
5. **WITH TIME ZONE** - для глобальных операций

## Производительность

1. **Индексы на все FK** - быстрые JOIN'ы
2. **Индекс на order_reference** - быстрый поиск по заказу
3. **Партиционирование** - возможно по created_at
4. **Архивация старых** - перенос в архивные таблицы

## Миграции

- **000046** - создание таблицы payment_transactions
- Будущие: добавление поддержки других шлюзов