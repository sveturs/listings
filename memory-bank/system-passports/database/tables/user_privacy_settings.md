# Паспорт таблицы: user_privacy_settings

## Назначение
Таблица настроек приватности пользователей. Контролирует, кто может отправлять запросы на добавление в контакты и сообщения пользователю.

## Структура таблицы

```sql
CREATE TABLE IF NOT EXISTS user_privacy_settings (
    user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    allow_contact_requests BOOLEAN NOT NULL DEFAULT TRUE,
    allow_messages_from_contacts_only BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

## Поля таблицы

### Основные поля
- `user_id` - уникальный идентификатор пользователя (первичный ключ)
- `allow_contact_requests` - разрешать ли другим пользователям отправлять запросы на добавление в контакты
- `allow_messages_from_contacts_only` - разрешать ли сообщения только от контактов

### Системные поля
- `created_at` - дата создания записи
- `updated_at` - дата последнего обновления

## Индексы

1. **PRIMARY KEY (user_id)** - первичный ключ по user_id
2. **idx_user_privacy_settings_user_id** - дополнительный индекс по user_id (технически избыточен)

## Триггеры

- **update_user_privacy_settings_updated_at** - автоматически обновляет поле updated_at при изменении записи

## Связи с другими таблицами

### Прямые связи (эта таблица ссылается на)
- `users(id)` → `user_id` - настройки принадлежат пользователю

### Обратные связи (другие таблицы ссылаются на user_privacy_settings)
- Нет прямых связей, но логически связана с user_contacts и системой чатов

## Бизнес-правила

1. **Контроль запросов на контакты**:
   - `allow_contact_requests = TRUE` - пользователь может получать запросы на добавление в контакты
   - `allow_contact_requests = FALSE` - запросы на контакты блокируются

2. **Контроль сообщений**:
   - `allow_messages_from_contacts_only = FALSE` - сообщения от всех пользователей разрешены
   - `allow_messages_from_contacts_only = TRUE` - сообщения только от добавленных в контакты

3. **Значения по умолчанию**:
   - Разрешены запросы на контакты (allow_contact_requests = TRUE)
   - Сообщения от всех пользователей (allow_messages_from_contacts_only = FALSE)

4. **Автоматическая инициализация** - при создании пользователя автоматически создаются настройки по умолчанию

## Примеры использования

### Создание настроек для нового пользователя
```sql
INSERT INTO user_privacy_settings (user_id)
VALUES (123);
```

### Изменение настроек приватности
```sql
UPDATE user_privacy_settings 
SET allow_contact_requests = FALSE,
    allow_messages_from_contacts_only = TRUE
WHERE user_id = 123;
```

### Получение настроек пользователя
```sql
SELECT allow_contact_requests, allow_messages_from_contacts_only
FROM user_privacy_settings
WHERE user_id = 123;
```

### Проверка возможности отправки запроса на контакт
```sql
SELECT ups.allow_contact_requests
FROM user_privacy_settings ups
WHERE ups.user_id = 456; -- ID получателя запроса
```

### Проверка возможности отправки сообщения
```sql
SELECT 
    ups.allow_messages_from_contacts_only,
    CASE 
        WHEN ups.allow_messages_from_contacts_only = FALSE THEN TRUE
        ELSE EXISTS (
            SELECT 1 FROM user_contacts uc 
            WHERE uc.user_id = 456 AND uc.contact_user_id = 123 AND uc.status = 'accepted'
        )
    END as can_send_message
FROM user_privacy_settings ups
WHERE ups.user_id = 456; -- ID получателя сообщения
```

### Массовая инициализация для существующих пользователей
```sql
INSERT INTO user_privacy_settings (user_id)
SELECT id FROM users
ON CONFLICT (user_id) DO NOTHING;
```

## Известные особенности

1. **Автоинициализация** - при создании пользователя в системе автоматически создается запись с настройками по умолчанию.

2. **Каскадное удаление** - при удалении пользователя настройки приватности удаляются автоматически.

3. **Значения NOT NULL** - все булевые поля имеют NOT NULL ограничения для предотвращения неопределенных состояний.

4. **Интеграция с контактами** - настройки влияют на поведение системы контактов и чатов.

5. **Триггер обновления** - автоматически отслеживает изменения для аудита.

## Интеграция с другими компонентами

1. **Система контактов** - проверяется allow_contact_requests при отправке запросов на добавление в контакты

2. **Система чатов** - проверяется allow_messages_from_contacts_only при отправке сообщений

3. **API endpoints** - настройки используются в бэкенде для авторизации действий

4. **Пользовательский интерфейс** - настройки отображаются в профиле пользователя для редактирования

## Сценарии использования

### Публичный пользователь (настройки по умолчанию)
- Принимает запросы на контакты от всех
- Принимает сообщения от всех

### Приватный пользователь
- Блокирует запросы на контакты
- Принимает сообщения только от контактов

### Полуприватный пользователь
- Принимает запросы на контакты
- Принимает сообщения только от контактов

## Миграции

- **000040** - создание таблицы user_privacy_settings с триггерами
- **000040** - автоматическая инициализация настроек для существующих пользователей