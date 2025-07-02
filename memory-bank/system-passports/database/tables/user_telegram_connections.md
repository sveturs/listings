# Паспорт таблицы: user_telegram_connections

## Назначение
Таблица для связи пользователей платформы с их Telegram аккаунтами. Позволяет отправлять уведомления через Telegram бота и интегрировать функционал с мессенджером.

## Структура таблицы

```sql
CREATE TABLE user_telegram_connections (
    user_id INT PRIMARY KEY REFERENCES users(id),
    telegram_chat_id VARCHAR(100) NOT NULL,
    telegram_username VARCHAR(100),
    connected_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Поля таблицы

### Основные поля
- `user_id` - уникальный идентификатор пользователя (первичный ключ)
- `telegram_chat_id` - идентификатор чата в Telegram (для отправки сообщений через бота)
- `telegram_username` - имя пользователя в Telegram (для отображения)

### Системные поля
- `connected_at` - дата и время подключения Telegram аккаунта

## Индексы

1. **PRIMARY KEY (user_id)** - первичный ключ по user_id, обеспечивает быстрый поиск по пользователю
2. Дополнительные индексы могут быть добавлены для telegram_chat_id при необходимости

## Триггеры

- Специальные триггеры не определены

## Связи с другими таблицами

### Прямые связи (эта таблица ссылается на)
- `users(id)` → `user_id` - Telegram подключение принадлежит пользователю

### Обратные связи (другие таблицы ссылаются на user_telegram_connections)
- Нет прямых связей, но логически используется системой уведомлений

## Бизнес-правила

1. **Уникальность подключения** - один пользователь может иметь только одно Telegram подключение

2. **Идентификация через chat_id** - telegram_chat_id используется для отправки сообщений через Telegram Bot API

3. **Опциональность username** - telegram_username может быть NULL, так как не все пользователи Telegram имеют публичные username

4. **Каскадное удаление** - при удалении пользователя его Telegram подключение удаляется автоматически

## Примеры использования

### Подключение Telegram аккаунта
```sql
INSERT INTO user_telegram_connections (user_id, telegram_chat_id, telegram_username)
VALUES (123, '456789012', 'john_doe_tg');
```

### Подключение без username
```sql
INSERT INTO user_telegram_connections (user_id, telegram_chat_id)
VALUES (124, '456789013');
```

### Получение Telegram данных пользователя
```sql
SELECT telegram_chat_id, telegram_username, connected_at
FROM user_telegram_connections
WHERE user_id = 123;
```

### Поиск пользователя по Telegram chat_id
```sql
SELECT u.id, u.name, u.email
FROM users u
JOIN user_telegram_connections utc ON u.id = utc.user_id
WHERE utc.telegram_chat_id = '456789012';
```

### Получение всех подключенных Telegram аккаунтов
```sql
SELECT u.id, u.name, utc.telegram_username, utc.connected_at
FROM users u
JOIN user_telegram_connections utc ON u.id = utc.user_id
ORDER BY utc.connected_at DESC;
```

### Отключение Telegram аккаунта
```sql
DELETE FROM user_telegram_connections
WHERE user_id = 123;
```

### Обновление Telegram username
```sql
UPDATE user_telegram_connections
SET telegram_username = 'new_username'
WHERE user_id = 123;
```

## Известные особенности

1. **Chat ID vs Username** - telegram_chat_id является основным идентификатором для отправки сообщений, telegram_username используется только для отображения.

2. **Числовой chat_id в VARCHAR** - хотя chat_id обычно числовой, он хранится как VARCHAR для поддержки различных типов чатов и будущих изменений API.

3. **Отсутствие уникальности chat_id** - теоретически один telegram_chat_id может быть привязан к разным пользователям, но это должно контролироваться логикой приложения.

4. **Timezone-naive timestamp** - connected_at использует простой TIMESTAMP без timezone информации.

## Интеграция с другими компонентами

1. **Система уведомлений** - использует telegram_chat_id для отправки push-уведомлений через Telegram бота

2. **Telegram бот** - бот использует эту таблицу для связи команд Telegram с пользователями платформы

3. **Настройки уведомлений** - интегрируется с notification_settings для контроля какие уведомления отправлять в Telegram

4. **Аутентификация** - может использоваться как дополнительный метод подтверждения личности

## Сценарии использования

### Подключение через Telegram бота
1. Пользователь открывает бота в Telegram
2. Бот генерирует код привязки
3. Пользователь вводит код на сайте
4. Система создает запись в user_telegram_connections

### Отправка уведомления
1. Система генерирует уведомление для пользователя
2. Проверяет наличие Telegram подключения
3. Отправляет сообщение через Telegram Bot API используя chat_id

### Управление подключением
1. Пользователь может отключить Telegram в настройках
2. Запись удаляется из таблицы
3. Уведомления в Telegram прекращаются

## Безопасность

1. **Защита chat_id** - telegram_chat_id должен быть защищен от утечки, так как может использоваться для спама

2. **Валидация подключения** - необходимо периодически проверять активность подключений

3. **Логирование изменений** - рекомендуется логировать подключения/отключения для безопасности

## Миграции

- **000001** - создание таблицы user_telegram_connections в составе основной миграции