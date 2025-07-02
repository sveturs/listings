# Паспорт таблицы: user_contacts

## Назначение
Таблица для управления контактами пользователей. Позволяет пользователям добавлять друг друга в контакты, управлять запросами на добавление в контакты и блокировать нежелательных пользователей.

## Структура таблицы

```sql
CREATE TABLE IF NOT EXISTS user_contacts (
    id BIGSERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    contact_user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending, accepted, blocked
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- Дополнительные поля
    added_from_chat_id INTEGER REFERENCES marketplace_chats(id), -- Откуда добавлен контакт
    notes TEXT, -- Заметки о контакте
    
    -- Индексы и ограничения
    UNIQUE(user_id, contact_user_id),
    CHECK (user_id != contact_user_id), -- Нельзя добавить себя в контакты
    CHECK (status IN ('pending', 'accepted', 'blocked'))
);
```

## Поля таблицы

### Основные поля
- `id` - уникальный идентификатор записи контакта (BIGSERIAL)
- `user_id` - ID пользователя, который добавляет контакт
- `contact_user_id` - ID пользователя, которого добавляют в контакты
- `status` - статус контакта: 'pending', 'accepted', 'blocked'

### Контекстные поля
- `added_from_chat_id` - ID чата, из которого был добавлен контакт (опционально)
- `notes` - заметки пользователя о контакте

### Системные поля
- `created_at` - дата создания записи
- `updated_at` - дата последнего обновления

## Индексы

1. **idx_user_contacts_user_id** - индекс по user_id для быстрого поиска контактов пользователя
2. **idx_user_contacts_contact_user_id** - индекс по contact_user_id
3. **idx_user_contacts_status** - индекс по статусу контакта
4. **idx_user_contacts_created_at** - индекс по дате создания

## Триггеры

- **update_user_contacts_updated_at** - автоматически обновляет поле updated_at при изменении записи

## Ограничения

1. **UNIQUE(user_id, contact_user_id)** - предотвращает дублирование записей контактов
2. **CHECK (user_id != contact_user_id)** - пользователь не может добавить себя в контакты
3. **CHECK (status IN ('pending', 'accepted', 'blocked'))** - валидация статуса

## Связи с другими таблицами

### Прямые связи (эта таблица ссылается на)
- `users(id)` → `user_id` - пользователь, который добавляет контакт
- `users(id)` → `contact_user_id` - пользователь, которого добавляют
- `marketplace_chats(id)` → `added_from_chat_id` - чат, из которого добавлен контакт

### Обратные связи (другие таблицы ссылаются на user_contacts)
- Нет прямых связей

## Бизнес-правила

1. **Статусы контактов**:
   - `pending` - запрос на добавление в контакты отправлен, ожидает подтверждения
   - `accepted` - контакт принят, пользователи в контактах друг у друга
   - `blocked` - пользователь заблокирован

2. **Двусторонние отношения** - когда пользователь A добавляет пользователя B, создается запись только в одну сторону. При принятии запроса может создаваться обратная запись.

3. **Блокировка** - блокировка работает в одну сторону. Заблокированный пользователь не может отправлять сообщения блокирующему.

4. **Контекст добавления** - система отслеживает, из какого чата был добавлен контакт для лучшего UX.

## Примеры использования

### Отправка запроса на добавление в контакты
```sql
INSERT INTO user_contacts (user_id, contact_user_id, status, added_from_chat_id) 
VALUES (123, 456, 'pending', 789);
```

### Принятие запроса на добавление в контакты
```sql
UPDATE user_contacts 
SET status = 'accepted', updated_at = CURRENT_TIMESTAMP 
WHERE user_id = 456 AND contact_user_id = 123 AND status = 'pending';
```

### Блокировка пользователя
```sql
INSERT INTO user_contacts (user_id, contact_user_id, status) 
VALUES (123, 456, 'blocked')
ON CONFLICT (user_id, contact_user_id) 
DO UPDATE SET status = 'blocked', updated_at = CURRENT_TIMESTAMP;
```

### Получение списка контактов пользователя
```sql
SELECT u2.id, u2.name, u2.picture_url, uc.status, uc.notes
FROM user_contacts uc
JOIN users u2 ON u2.id = uc.contact_user_id
WHERE uc.user_id = 123 AND uc.status = 'accepted'
ORDER BY u2.name;
```

### Получение входящих запросов на добавление в контакты
```sql
SELECT u1.id, u1.name, u1.picture_url, uc.created_at
FROM user_contacts uc
JOIN users u1 ON u1.id = uc.user_id
WHERE uc.contact_user_id = 123 AND uc.status = 'pending'
ORDER BY uc.created_at DESC;
```

## Известные особенности

1. **Односторонние записи** - система создает запись только от инициатора запроса, что требует careful обработки при поиске взаимных контактов.

2. **Блокировка приоритетна** - если пользователь заблокирован, он не может отправлять сообщения независимо от других настроек.

3. **Связь с чатами** - система отслеживает контекст добавления контакта из чатов для улучшения UX.

4. **Заметки** - пользователи могут добавлять приватные заметки о своих контактах.

5. **BIGSERIAL ID** - используется BIGSERIAL вместо SERIAL для поддержки большого количества записей.

## Интеграция с другими компонентами

1. **Чаты** - проверяется статус контакта при отправке сообщений
2. **Настройки приватности** - взаимодействует с user_privacy_settings для контроля доступа
3. **Уведомления** - генерируются уведомления при изменении статуса контактов

## Миграции

- **000039** - создание таблицы user_contacts с индексами и триггерами