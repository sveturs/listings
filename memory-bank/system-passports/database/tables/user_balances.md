# Паспорт таблицы: user_balances

## Назначение
Хранение текущего баланса пользователей в системе. Основная таблица для всех финансовых операций на платформе.

## Структура таблицы

```sql
CREATE TABLE IF NOT EXISTS user_balances (
    user_id INT PRIMARY KEY REFERENCES users(id),
    balance DECIMAL(12,2) NOT NULL DEFAULT 0,
    frozen_balance DECIMAL(12,2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'RSD',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Поля таблицы

### Ключевые поля
- `user_id` - ID пользователя (PRIMARY KEY, FK к users)

### Финансовые поля
- `balance` - доступный баланс (до 12 цифр, 2 после запятой, по умолчанию 0)
- `frozen_balance` - замороженный баланс (для pending операций)
- `currency` - валюта баланса (3 символа ISO 4217, по умолчанию 'RSD')

### Системные поля
- `updated_at` - время последнего обновления баланса

## Индексы

1. **PRIMARY KEY** на user_id (автоматический индекс)

## Связи с другими таблицами

### Прямые связи
- `user_id` → `users.id` - владелец баланса

### Обратные связи
- `balance_transactions.user_id` - транзакции пользователя
- `payment_transactions.user_id` - платежные транзакции
- `escrow_payments` - эскроу платежи
- `merchant_payouts` - выплаты продавцам

## Бизнес-правила

### Управление балансом
1. **Один баланс на пользователя** - PRIMARY KEY гарантирует уникальность
2. **Неотрицательный баланс** - проверка на уровне приложения
3. **Атомарные операции** - все изменения через транзакции

### Замороженный баланс
- Используется для резервирования средств
- При покупке: frozen_balance += amount
- При подтверждении: frozen_balance -= amount, balance -= amount
- При отмене: frozen_balance -= amount

### Валюты
- `RSD` - Сербский динар (основная)
- `EUR` - Евро
- `USD` - Доллар США

### Точность вычислений
- DECIMAL(12,2) - до 10 миллиардов с точностью до копеек
- Все расчеты в минимальных единицах валюты

## Примеры использования

### Создание баланса для нового пользователя
```sql
INSERT INTO user_balances (user_id, currency)
VALUES (:user_id, 'RSD')
ON CONFLICT (user_id) DO NOTHING;
```

### Получение баланса пользователя
```sql
SELECT 
    ub.balance,
    ub.frozen_balance,
    ub.currency,
    (ub.balance - ub.frozen_balance) as available_balance
FROM user_balances ub
WHERE ub.user_id = :user_id;
```

### Пополнение баланса
```sql
-- В транзакции
BEGIN;

-- Обновляем баланс
UPDATE user_balances 
SET balance = balance + :amount,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = :user_id;

-- Записываем транзакцию
INSERT INTO balance_transactions (
    user_id, type, amount, currency, status, description
) VALUES (
    :user_id, 'deposit', :amount, 'RSD', 'completed', :description
);

COMMIT;
```

### Заморозка средств для покупки
```sql
-- Проверяем доступные средства
SELECT (balance - frozen_balance) as available
FROM user_balances
WHERE user_id = :user_id;

-- Если достаточно, замораживаем
UPDATE user_balances
SET frozen_balance = frozen_balance + :amount,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = :user_id
  AND (balance - frozen_balance) >= :amount;
```

### Перевод между пользователями
```sql
BEGIN;

-- Проверяем и списываем у отправителя
UPDATE user_balances
SET balance = balance - :amount,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = :sender_id
  AND balance >= :amount;

-- Зачисляем получателю
UPDATE user_balances
SET balance = balance + :amount,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = :receiver_id;

-- Записываем обе транзакции
INSERT INTO balance_transactions (user_id, type, amount, status)
VALUES 
    (:sender_id, 'withdrawal', :amount, 'completed'),
    (:receiver_id, 'deposit', :amount, 'completed');

COMMIT;
```

## Известные особенности

1. **Отсутствие истории в таблице** - вся история в balance_transactions
2. **PRIMARY KEY на user_id** - автоматически создает уникальный индекс
3. **Мультивалютность** - подготовлена, но используется только RSD
4. **Нет триггеров** - все обновления через приложение для контроля
5. **Frozen balance** - для безопасных транзакций без блокировок

## API интеграция

### Endpoints
- `GET /api/v1/users/balance` - текущий баланс
- `POST /api/v1/users/balance/deposit` - пополнение
- `POST /api/v1/users/balance/withdraw` - вывод средств
- `GET /api/v1/users/balance/history` - история транзакций

### Пример ответа
```json
{
  "balance": "1500.00",
  "frozen_balance": "200.00",
  "available_balance": "1300.00",
  "currency": "RSD",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

## Безопасность

1. **Двойная запись** - баланс + транзакция для аудита
2. **Оптимистичная блокировка** - через условия WHERE
3. **Транзакционность** - все операции в транзакциях БД
4. **Проверка лимитов** - на уровне приложения
5. **Логирование** - все операции записываются

## Производительность

1. **Денормализация** - баланс хранится отдельно от транзакций
2. **Индекс по user_id** - быстрый доступ
3. **Без JOIN'ов** - для частых операций чтения
4. **Кеширование** - возможно в Redis для high-load

## Миграции

- **000001** - создание таблицы
- Будущие: добавление multi-currency support