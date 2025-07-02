# Паспорт таблицы: notifications

## Назначение
Централизованное хранение всех уведомлений пользователей. Поддерживает различные типы уведомлений и каналы доставки.

## Структура таблицы

```sql
CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    type VARCHAR(50) NOT NULL,
    title TEXT NOT NULL,
    message TEXT NOT NULL,
    data JSONB,
    is_read BOOLEAN DEFAULT false,
    delivered_to JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Поля таблицы

### Основные поля
- `id` - уникальный идентификатор уведомления (SERIAL)
- `user_id` - получатель уведомления (FK к users, NOT NULL)
- `type` - тип уведомления (до 50 символов)

### Содержимое
- `title` - заголовок уведомления (обязательное)
- `message` - текст уведомления (обязательное)
- `data` - дополнительные данные в JSON формате

### Статус
- `is_read` - прочитано ли уведомление (по умолчанию false)
- `delivered_to` - информация о доставке по каналам (JSON)

### Системные поля
- `created_at` - время создания уведомления

## Индексы

1. **idx_notifications_user** - поиск уведомлений пользователя
2. **idx_notifications_type** - фильтрация по типу
3. **idx_notifications_created** - сортировка по времени
4. **idx_notifications_user_unread** - непрочитанные уведомления пользователя

## Связи с другими таблицами

### Прямые связи
- `user_id` → `users.id` - получатель уведомления

### Обратные связи
- Нет прямых ссылок

### Связи через data поле
- `data.listing_id` → объявление
- `data.chat_id` → чат
- `data.order_id` → заказ
- `data.review_id` → отзыв

## Типы уведомлений

### Marketplace
- `new_message` - новое сообщение в чате
- `listing_favorited` - объявление добавлено в избранное
- `price_changed` - изменилась цена избранного товара
- `listing_sold` - товар продан

### Транзакции
- `payment_received` - получен платеж
- `payment_sent` - отправлен платеж
- `balance_low` - низкий баланс
- `payout_completed` - выплата завершена

### Отзывы
- `new_review` - новый отзыв
- `review_response` - ответ на отзыв
- `review_helpful` - отзыв отмечен полезным

### Системные
- `welcome` - приветственное уведомление
- `email_verified` - email подтвержден
- `account_suspended` - аккаунт заблокирован
- `system_maintenance` - техническое обслуживание

## Структура поля data

```json
{
  // Для new_message
  "chat_id": 123,
  "sender_id": 456,
  "sender_name": "Иван Иванов",
  "message_preview": "Здравствуйте, товар еще в наличии?",
  
  // Для listing_favorited
  "listing_id": 789,
  "listing_title": "iPhone 13 Pro",
  "user_name": "Петр Петров",
  
  // Для payment_received
  "amount": 1000.00,
  "currency": "RSD",
  "from_user_id": 321,
  "from_user_name": "Мария Иванова",
  "transaction_id": "tx_123456"
}
```

## Структура поля delivered_to

```json
{
  "in_app": {
    "delivered": true,
    "delivered_at": "2024-01-15T10:30:00Z"
  },
  "email": {
    "delivered": true,
    "delivered_at": "2024-01-15T10:30:15Z",
    "message_id": "msg_abc123"
  },
  "telegram": {
    "delivered": false,
    "error": "User not connected to Telegram"
  },
  "push": {
    "delivered": true,
    "delivered_at": "2024-01-15T10:30:05Z",
    "token": "fcm_token_xyz"
  }
}
```

## Примеры использования

### Создание уведомления о новом сообщении
```sql
INSERT INTO notifications (user_id, type, title, message, data)
VALUES (
    :receiver_id,
    'new_message',
    'Новое сообщение',
    'Вам пришло новое сообщение от ' || :sender_name,
    jsonb_build_object(
        'chat_id', :chat_id,
        'sender_id', :sender_id,
        'sender_name', :sender_name,
        'message_preview', LEFT(:message_text, 100)
    )
);
```

### Получение непрочитанных уведомлений
```sql
SELECT * FROM notifications
WHERE user_id = :user_id
  AND is_read = false
ORDER BY created_at DESC
LIMIT 20;
```

### Отметка уведомлений как прочитанных
```sql
UPDATE notifications
SET is_read = true
WHERE user_id = :user_id
  AND id = ANY(:notification_ids);
```

### Массовое уведомление
```sql
-- Уведомление всем пользователям о техработах
INSERT INTO notifications (user_id, type, title, message, data)
SELECT 
    id,
    'system_maintenance',
    'Плановые технические работы',
    'Сервис будет недоступен с 02:00 до 04:00',
    jsonb_build_object(
        'start_time', '2024-01-20T02:00:00Z',
        'end_time', '2024-01-20T04:00:00Z'
    )
FROM users
WHERE account_status = 'active';
```

### Обновление статуса доставки
```sql
UPDATE notifications
SET delivered_to = delivered_to || jsonb_build_object(
    'email', jsonb_build_object(
        'delivered', true,
        'delivered_at', CURRENT_TIMESTAMP,
        'message_id', :email_message_id
    )
)
WHERE id = :notification_id;
```

## Бизнес-правила

### Создание уведомлений
1. **Обязательные поля** - user_id, type, title, message
2. **Проверка получателя** - пользователь должен существовать
3. **Дедупликация** - избегать дублей за короткий период

### Доставка
1. **Мультиканальность** - in-app, email, telegram, push
2. **Настройки пользователя** - учитывать notification_settings
3. **Retry механизм** - повторная отправка при сбое

### Хранение
1. **Retention период** - старые уведомления удаляются через 90 дней
2. **Архивация** - важные уведомления архивируются
3. **Лимиты** - максимум 1000 непрочитанных на пользователя

## API интеграция

### Endpoints
- `GET /api/v1/notifications` - список уведомлений
- `GET /api/v1/notifications/unread/count` - счетчик непрочитанных
- `PUT /api/v1/notifications/{id}/read` - отметить прочитанным
- `PUT /api/v1/notifications/read-all` - прочитать все

### WebSocket события
- `notification:new` - новое уведомление
- `notification:read` - уведомление прочитано

## Известные особенности

1. **JSONB для гибкости** - data и delivered_to хранят произвольные данные
2. **Нет updated_at** - уведомления не редактируются после создания
3. **Каскадное удаление отсутствует** - для сохранения истории
4. **Индекс для непрочитанных** - оптимизация частого запроса
5. **Type как строка** - гибкость добавления новых типов

## Производительность

1. **Пагинация обязательна** - для списков уведомлений
2. **Очистка старых** - регулярный cron job
3. **Кеширование счетчиков** - unread count в Redis
4. **Batch insert** - для массовых уведомлений

## Миграции

- **000001** - создание таблицы
- **000039** - добавление индекса для непрочитанных