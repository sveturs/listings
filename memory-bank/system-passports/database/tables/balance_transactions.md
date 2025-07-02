# Паспорт таблицы: balance_transactions

## Назначение
История всех финансовых операций с балансами пользователей. Обеспечивает полный аудит и возможность восстановления баланса из истории транзакций.

## Структура таблицы

```sql
CREATE TABLE IF NOT EXISTS balance_transactions (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    type VARCHAR(20) NOT NULL,
    amount DECIMAL(12,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'RSD',
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    payment_method VARCHAR(50),
    payment_details JSONB,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP
);
```

## Поля таблицы

### Основные поля
- `id` - уникальный идентификатор транзакции (SERIAL)
- `user_id` - пользователь транзакции (FK к users, NOT NULL)

### Финансовые данные
- `type` - тип транзакции (до 20 символов)
- `amount` - сумма транзакции (до 12 цифр, 2 после запятой)
- `currency` - валюта (3 символа ISO 4217, по умолчанию 'RSD')

### Статус транзакции
- `status` - статус транзакции (по умолчанию 'pending')
- `completed_at` - время завершения транзакции

### Платежная информация
- `payment_method` - способ платежа (до 50 символов)
- `payment_details` - детали платежа в JSON формате
- `description` - описание транзакции

### Системные поля
- `created_at` - время создания записи

## Типы транзакций

### Основные типы (type)
- `deposit` - пополнение баланса
- `withdrawal` - вывод средств
- `transfer` - перевод между пользователями
- `service_payment` - оплата услуг платформы
- `purchase` - покупка товара
- `sale` - продажа товара
- `commission` - комиссия платформы
- `refund` - возврат средств

### Статусы (status)
- `pending` - ожидает обработки
- `completed` - успешно завершена
- `failed` - неудачная транзакция
- `cancelled` - отменена пользователем

### Методы платежа (payment_method)
- `bank_transfer` - банковский перевод
- `payment_slip` - платежная квитанция
- `crypto` - криптовалюта
- `card` - банковская карта
- `balance` - с баланса системы

## Индексы

1. **idx_transactions_user** - по пользователю
2. **idx_transactions_status** - по статусу
3. **idx_transactions_created** - по дате создания

## Связи с другими таблицами

### Прямые связи
- `user_id` → `users.id` - пользователь транзакции

### Обратные связи
- `user_storefronts.creation_transaction_id` - транзакция создания витрины

### Логические связи
- Связана с `user_balances` через user_id
- Может ссылаться на `payment_transactions` через payment_details

## Структура payment_details

```json
{
  // Для bank_transfer
  "bank_name": "Raiffeisen Bank",
  "account_number": "265-1234567890-12",
  "reference_number": "97-1234567890",
  
  // Для card payment
  "card_last4": "4242",
  "card_brand": "Visa",
  "gateway_transaction_id": "tx_123456",
  
  // Для transfer
  "from_user_id": 123,
  "to_user_id": 456,
  "transfer_reason": "Payment for iPhone 13",
  
  // Для purchase/sale
  "listing_id": 789,
  "order_id": "ord_123456",
  "seller_id": 321,
  "buyer_id": 654
}
```

## Примеры использования

### Создание транзакции пополнения
```sql
INSERT INTO balance_transactions (
    user_id, type, amount, currency, status, 
    payment_method, payment_details, description
) VALUES (
    :user_id, 
    'deposit', 
    1000.00, 
    'RSD', 
    'pending',
    'bank_transfer',
    jsonb_build_object(
        'bank_name', 'Raiffeisen Bank',
        'reference_number', '97-1234567890'
    ),
    'Пополнение через банковский перевод'
);
```

### Завершение транзакции
```sql
-- В транзакции БД
BEGIN;

-- Обновляем статус транзакции
UPDATE balance_transactions
SET status = 'completed',
    completed_at = CURRENT_TIMESTAMP
WHERE id = :transaction_id
  AND status = 'pending';

-- Обновляем баланс пользователя
UPDATE user_balances
SET balance = balance + (
    SELECT amount FROM balance_transactions 
    WHERE id = :transaction_id
),
updated_at = CURRENT_TIMESTAMP
WHERE user_id = (
    SELECT user_id FROM balance_transactions 
    WHERE id = :transaction_id
);

COMMIT;
```

### История транзакций пользователя
```sql
SELECT 
    bt.*,
    CASE 
        WHEN bt.type IN ('deposit', 'sale', 'refund') THEN '+'
        WHEN bt.type IN ('withdrawal', 'purchase', 'commission') THEN '-'
        ELSE '±'
    END as direction
FROM balance_transactions bt
WHERE bt.user_id = :user_id
ORDER BY bt.created_at DESC
LIMIT 50 OFFSET :offset;
```

### Расчет баланса из истории
```sql
-- Проверка консистентности
SELECT 
    SUM(CASE 
        WHEN type IN ('deposit', 'sale', 'refund') THEN amount
        WHEN type IN ('withdrawal', 'purchase', 'commission') THEN -amount
        ELSE 0
    END) as calculated_balance
FROM balance_transactions
WHERE user_id = :user_id
  AND status = 'completed';
```

### Отчет по типам транзакций
```sql
SELECT 
    type,
    COUNT(*) as count,
    SUM(amount) as total_amount,
    AVG(amount) as avg_amount
FROM balance_transactions
WHERE user_id = :user_id
  AND status = 'completed'
  AND created_at >= DATE_TRUNC('month', CURRENT_DATE)
GROUP BY type
ORDER BY total_amount DESC;
```

## Бизнес-правила

### Создание транзакций
1. **Immutable records** - транзакции не редактируются после создания
2. **Статус pending** - все транзакции начинаются как pending
3. **Обязательное описание** - для удобства пользователя

### Обработка транзакций
1. **Атомарность** - обновление транзакции и баланса в одной БД транзакции
2. **Идемпотентность** - повторное выполнение не изменяет результат
3. **Аудит** - все изменения логируются

### Валидация
1. **Положительная сумма** - amount > 0
2. **Существующий пользователь** - проверка FK
3. **Корректный тип** - из списка разрешенных

## API интеграция

### Endpoints
- `GET /api/v1/users/transactions` - история транзакций
- `GET /api/v1/users/transactions/{id}` - детали транзакции
- `POST /api/v1/users/transactions/deposit` - создать депозит
- `POST /api/v1/users/transactions/withdraw` - запрос на вывод

### Пример ответа
```json
{
  "id": 12345,
  "type": "deposit",
  "amount": "1000.00",
  "currency": "RSD",
  "status": "completed",
  "payment_method": "bank_transfer",
  "description": "Пополнение через банковский перевод",
  "created_at": "2024-01-15T10:00:00Z",
  "completed_at": "2024-01-15T10:30:00Z",
  "direction": "+"
}
```

## Известные особенности

1. **Immutable design** - записи не изменяются, только добавляются
2. **JSONB для гибкости** - payment_details хранит разные структуры
3. **Двойная запись** - синхронизация с user_balances
4. **Нет CASCADE DELETE** - для сохранения истории
5. **Временная метка завершения** - отдельное поле completed_at

## Производительность

1. **Индексы покрывают** - основные сценарии запросов
2. **Пагинация обязательна** - для истории транзакций
3. **Партиционирование** - возможно по created_at для архива
4. **Агрегация** - периодический расчет статистики

## Миграции

- **000001** - создание таблицы и индексов
- Будущие: партиционирование для масштабирования