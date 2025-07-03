# Паспорт таблицы: notification_settings

## Назначение
Персональные настройки уведомлений для каждого пользователя. Определяет, какие типы уведомлений и через какие каналы пользователь хочет получать.

## Структура таблицы

```sql
CREATE TABLE notification_settings (
    user_id INT NOT NULL REFERENCES users(id),
    notification_type VARCHAR(50) NOT NULL,
    telegram_enabled BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, notification_type)
);
```

## Поля таблицы

### Ключевые поля
- `user_id` - ID пользователя (FK к users)
- `notification_type` - тип уведомления (до 50 символов)
- **PRIMARY KEY** - составной ключ (user_id, notification_type)

### Настройки каналов
- `telegram_enabled` - включены ли Telegram уведомления (по умолчанию false)

### Системные поля
- `created_at` - дата создания настройки
- `updated_at` - дата последнего изменения

## Индексы

1. **PRIMARY KEY** - составной индекс по (user_id, notification_type)

## Триггеры

- **update_notification_settings_timestamp** - автоматическое обновление updated_at

## Связи с другими таблицами

### Прямые связи
- `user_id` → `users.id` - пользователь

### Логические связи
- Связан с `notifications.type` - определяет настройки для типов уведомлений
- Связан с `user_telegram_connections` - для проверки возможности отправки

## Типы уведомлений

Те же типы, что и в таблице `notifications`:

### Сообщения и чаты
- `new_message` - новое сообщение
- `chat_started` - начат новый чат

### Маркетплейс
- `listing_favorited` - товар добавлен в избранное
- `price_changed` - изменение цены
- `listing_sold` - товар продан
- `listing_expired` - объявление истекло

### Финансы
- `payment_received` - получен платеж
- `payment_sent` - отправлен платеж
- `balance_low` - низкий баланс
- `payout_completed` - выплата завершена

### Отзывы
- `new_review` - новый отзыв
- `review_response` - ответ на отзыв

### Системные
- `system_maintenance` - техработы
- `security_alert` - предупреждение безопасности
- `feature_announcement` - анонс новых функций

## Настройки по умолчанию

При регистрации пользователя создаются настройки по умолчанию:

```sql
-- Критические уведомления включены
('new_message', telegram_enabled: true)
('payment_received', telegram_enabled: true)
('security_alert', telegram_enabled: true)

-- Информационные выключены
('listing_favorited', telegram_enabled: false)
('feature_announcement', telegram_enabled: false)
```

## Примеры использования

### Создание настроек для нового пользователя
```sql
INSERT INTO notification_settings (user_id, notification_type, telegram_enabled)
VALUES 
    (:user_id, 'new_message', true),
    (:user_id, 'payment_received', true),
    (:user_id, 'security_alert', true),
    (:user_id, 'listing_favorited', false),
    (:user_id, 'price_changed', false)
ON CONFLICT (user_id, notification_type) DO NOTHING;
```

### Получение настроек пользователя
```sql
SELECT 
    nt.type,
    nt.description,
    COALESCE(ns.telegram_enabled, false) as telegram_enabled
FROM (
    VALUES 
        ('new_message', 'Новые сообщения'),
        ('payment_received', 'Получение платежей'),
        ('listing_favorited', 'Добавление в избранное')
) as nt(type, description)
LEFT JOIN notification_settings ns 
    ON ns.user_id = :user_id 
    AND ns.notification_type = nt.type
ORDER BY nt.type;
```

### Обновление настроек
```sql
INSERT INTO notification_settings (user_id, notification_type, telegram_enabled)
VALUES (:user_id, :notification_type, :telegram_enabled)
ON CONFLICT (user_id, notification_type) 
DO UPDATE SET 
    telegram_enabled = EXCLUDED.telegram_enabled,
    updated_at = CURRENT_TIMESTAMP;
```

### Проверка перед отправкой уведомления
```sql
-- Проверяем, нужно ли отправлять Telegram уведомление
SELECT 
    ns.telegram_enabled,
    utc.telegram_chat_id IS NOT NULL as has_telegram
FROM notification_settings ns
LEFT JOIN user_telegram_connections utc ON ns.user_id = utc.user_id
WHERE ns.user_id = :user_id 
  AND ns.notification_type = :notification_type;
```

### Массовое включение/выключение
```sql
-- Выключить все Telegram уведомления для пользователя
UPDATE notification_settings
SET telegram_enabled = false,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = :user_id;
```

## Бизнес-правила

### Управление настройками
1. **Инициализация** - при регистрации создаются базовые настройки
2. **Отсутствие записи** - означает использование дефолтных значений
3. **Новые типы** - автоматически выключены для существующих пользователей

### Каналы доставки
1. **In-app** - всегда включен, не настраивается
2. **Email** - управляется через users.notification_email
3. **Telegram** - управляется через эту таблицу
4. **Push** - будет добавлен в будущем

### Приоритеты
- Критические уведомления (платежи, безопасность) включены по умолчанию
- Маркетинговые уведомления выключены по умолчанию
- Пользователь может переопределить любые настройки

## API интеграция

### Endpoints
- `GET /api/v1/users/notification-settings` - получить настройки
- `PUT /api/v1/users/notification-settings` - обновить настройки
- `POST /api/v1/users/notification-settings/reset` - сброс к дефолтным

### Пример запроса
```json
PUT /api/v1/users/notification-settings
{
  "settings": [
    {
      "notification_type": "new_message",
      "telegram_enabled": true
    },
    {
      "notification_type": "price_changed",
      "telegram_enabled": false
    }
  ]
}
```

## Известные особенности

1. **Составной PRIMARY KEY** - одна запись на тип уведомления
2. **Расширяемость** - легко добавить новые каналы (email_enabled, push_enabled)
3. **ON CONFLICT** - безопасное обновление настроек
4. **Дефолтные значения** - обрабатываются на уровне приложения
5. **Связь с Telegram** - требует проверки user_telegram_connections

## Будущие расширения

```sql
-- Планируемые дополнительные поля
ALTER TABLE notification_settings
ADD COLUMN email_enabled BOOLEAN DEFAULT true,
ADD COLUMN push_enabled BOOLEAN DEFAULT false,
ADD COLUMN quiet_hours_start TIME,
ADD COLUMN quiet_hours_end TIME,
ADD COLUMN frequency VARCHAR(20) DEFAULT 'instant';
-- frequency: instant, hourly, daily, weekly
```

## Миграции

- **000021** - создание таблицы
- Будущие миграции для добавления новых каналов