# Паспорт таблицы: payment_methods

## Назначение
Справочник доступных способов оплаты на платформе. Хранит конфигурацию, лимиты и комиссии для каждого платежного метода.

## Структура таблицы

```sql
CREATE TABLE IF NOT EXISTS payment_methods (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE,
    type VARCHAR(50) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    minimum_amount DECIMAL(12,2),
    maximum_amount DECIMAL(12,2),
    fee_percentage DECIMAL(5,2),
    fixed_fee DECIMAL(12,2),
    credentials JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Поля таблицы

### Основные поля
- `id` - уникальный идентификатор метода (SERIAL)
- `name` - название метода для отображения (до 100 символов)
- `code` - уникальный код метода (до 50 символов, UNIQUE)
- `type` - тип платежного метода (до 50 символов)

### Настройки
- `is_active` - активен ли метод (по умолчанию true)

### Лимиты
- `minimum_amount` - минимальная сумма транзакции
- `maximum_amount` - максимальная сумма транзакции

### Комиссии
- `fee_percentage` - процент комиссии (до 5 цифр, 2 после запятой)
- `fixed_fee` - фиксированная комиссия

### Конфигурация
- `credentials` - настройки подключения (JSONB)

### Системные поля
- `created_at` - дата добавления метода

## Типы платежных методов (type)

- `bank` - банковские переводы
- `cash` - наличные платежи
- `digital` - цифровые платежи
- `card` - карточные платежи
- `crypto` - криптовалюты
- `payment_system` - платежные системы

## Предустановленные методы

```sql
INSERT INTO payment_methods (name, code, type, is_active, minimum_amount, maximum_amount, fee_percentage, fixed_fee) 
VALUES 
    ('Bank transfer', 'bank_transfer', 'bank', true, 1000, 10000000, 0, 100),
    ('Post office', 'post_office', 'cash', true, 500, 500000, 1.5, 50),
    ('IPS QR code', 'ips_qr', 'digital', true, 100, 1000000, 0.8, 0);
```

## Структура credentials

### Для bank_transfer
```json
{
  "bank_name": "Raiffeisen Bank",
  "account_number": "265-1234567890-12",
  "swift_code": "RZBSRSBG",
  "iban": "RS35265123456789012"
}
```

### Для payment gateway
```json
{
  "merchant_id": "MERCHANT123",
  "api_key": "encrypted_key",
  "api_secret": "encrypted_secret",
  "webhook_url": "https://api.svetu.rs/webhooks/payment",
  "test_mode": false
}
```

### Для crypto
```json
{
  "wallet_address": "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
  "network": "bitcoin",
  "confirmations_required": 3
}
```

## Индексы

1. **PRIMARY KEY** на id
2. **UNIQUE** на code

## Связи с другими таблицами

### Обратные связи
- `balance_transactions.payment_method` - использование в транзакциях
- `payment_transactions.payment_method` - использование в платежах

## Примеры использования

### Получение активных методов оплаты
```sql
SELECT 
    id,
    name,
    code,
    type,
    minimum_amount,
    maximum_amount,
    fee_percentage,
    fixed_fee
FROM payment_methods
WHERE is_active = true
ORDER BY 
    CASE type 
        WHEN 'bank' THEN 1
        WHEN 'digital' THEN 2
        WHEN 'cash' THEN 3
        ELSE 4
    END,
    name;
```

### Расчет комиссии
```sql
SELECT 
    name,
    code,
    :amount as amount,
    GREATEST(
        (:amount * fee_percentage / 100),
        COALESCE(fixed_fee, 0)
    ) as commission,
    :amount + GREATEST(
        (:amount * fee_percentage / 100),
        COALESCE(fixed_fee, 0)
    ) as total_amount
FROM payment_methods
WHERE code = :payment_method_code
  AND is_active = true
  AND :amount BETWEEN COALESCE(minimum_amount, 0) 
                  AND COALESCE(maximum_amount, 999999999);
```

### Проверка лимитов
```sql
SELECT 
    CASE
        WHEN NOT is_active THEN 'Method is not active'
        WHEN :amount < COALESCE(minimum_amount, 0) THEN 
            'Amount below minimum: ' || minimum_amount
        WHEN :amount > COALESCE(maximum_amount, 999999999) THEN 
            'Amount above maximum: ' || maximum_amount
        ELSE 'OK'
    END as validation_result
FROM payment_methods
WHERE code = :payment_method_code;
```

### Добавление нового метода
```sql
INSERT INTO payment_methods (
    name, code, type, is_active,
    minimum_amount, maximum_amount,
    fee_percentage, fixed_fee,
    credentials
) VALUES (
    'AllSecure Gateway',
    'allsecure',
    'payment_system',
    true,
    100,
    1000000,
    2.5,
    0,
    jsonb_build_object(
        'merchant_id', 'SVETU_MERCHANT',
        'channel_id', 'SVETU_CHANNEL',
        'test_mode', true
    )
);
```

## Бизнес-правила

### Активация/деактивация
1. **Мягкое удаление** - используется is_active вместо DELETE
2. **Проверка использования** - нельзя деактивировать используемые методы
3. **Минимум один активный** - всегда должен быть хотя бы один метод

### Комиссии
1. **Выбор большей** - между процентной и фиксированной
2. **Прозрачность** - комиссия показывается до оплаты
3. **Включение в сумму** - может добавляться к сумме или вычитаться

### Лимиты
1. **NULL как отсутствие** - NULL в лимитах означает без ограничений
2. **Валидация на уровне БД** - через CHECK constraints
3. **Разные лимиты** - для разных типов пользователей (на уровне приложения)

## API интеграция

### Endpoints
- `GET /api/v1/payment-methods` - список доступных методов
- `GET /api/v1/payment-methods/{code}` - детали метода
- `POST /api/v1/payment-methods/{code}/calculate` - расчет комиссии

### Пример ответа
```json
{
  "id": 1,
  "name": "Bank transfer",
  "code": "bank_transfer",
  "type": "bank",
  "minimum_amount": "1000.00",
  "maximum_amount": "10000000.00",
  "fee_percentage": "0.00",
  "fixed_fee": "100.00",
  "available": true
}
```

## Известные особенности

1. **UNIQUE на code** - для идентификации в API
2. **JSONB credentials** - гибкое хранение настроек
3. **Nullable лимиты** - отсутствие означает без ограничений
4. **Два типа комиссий** - процент и фиксированная
5. **Справочная таблица** - редко изменяется

## Безопасность

1. **Шифрование credentials** - чувствительные данные шифруются
2. **Ограничение доступа** - только админы могут изменять
3. **Аудит изменений** - логирование всех операций
4. **Валидация на уровне БД** - UNIQUE и CHECK constraints

## Производительность

1. **Кеширование** - методы редко меняются
2. **Малый размер** - обычно < 20 записей
3. **Индекс на code** - быстрый поиск по коду

## Миграции

- **000001** - создание таблицы
- **000001** - вставка начальных данных
- Будущие: добавление новых методов оплаты