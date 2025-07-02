# Паспорт таблицы: refresh_tokens

## Назначение
Таблица для хранения refresh токенов в системе JWT аутентификации. Позволяет пользователям оставаться авторизованными и обновлять access токены без повторного ввода учетных данных. Также обеспечивает управление сессиями и безопасность токенов.

## Структура таблицы

```sql
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Информация об устройстве/браузере
    user_agent TEXT,
    ip VARCHAR(45), -- Поддержка IPv6
    device_name VARCHAR(100),
    
    -- Для отзыва токенов
    is_revoked BOOLEAN DEFAULT FALSE,
    revoked_at TIMESTAMP
);
```

## Поля таблицы

### Основные поля
- `id` - уникальный идентификатор токена (SERIAL)
- `user_id` - ID пользователя, которому принадлежит токен
- `token` - строка refresh токена (уникальная, тип TEXT после миграции 000042)
- `expires_at` - дата и время истечения токена

### Информация о контексте
- `user_agent` - User-Agent браузера/приложения
- `ip` - IP адрес с поддержкой IPv6 (до 45 символов)
- `device_name` - опциональное имя устройства для управления сессиями

### Управление токенами
- `is_revoked` - флаг принудительного отзыва токена
- `revoked_at` - дата и время отзыва токена

### Системные поля
- `created_at` - дата создания токена

## Индексы

1. **PRIMARY KEY (id)** - первичный ключ
2. **UNIQUE (token)** - уникальность токена
3. **idx_refresh_tokens_user_id** - индекс по user_id для неотозванных токенов
4. **idx_refresh_tokens_token** - индекс по token для неотозванных токенов
5. **idx_refresh_tokens_expires_at** - индекс по expires_at для неотозванных токенов

## Функции

- **cleanup_expired_refresh_tokens()** - функция для автоматической очистки истекших и отозванных токенов

## Связи с другими таблицами

### Прямые связи (эта таблица ссылается на)
- `users(id)` → `user_id` - токен принадлежит пользователю

### Обратные связи (другие таблицы ссылаются на refresh_tokens)
- Нет прямых связей

## Бизнес-правила

1. **Время жизни токенов** - каждый токен имеет определенное время истечения

2. **Уникальность токенов** - каждый токен должен быть уникальным в системе

3. **Связь с пользователем** - токен всегда привязан к конкретному пользователю

4. **Отзыв токенов** - токены могут быть принудительно отозваны (is_revoked = TRUE)

5. **Автоочистка** - истекшие и отозванные токены автоматически удаляются через 30 дней

6. **Каскадное удаление** - при удалении пользователя все его токены удаляются

## Примеры использования

### Создание нового refresh токена
```sql
INSERT INTO refresh_tokens (user_id, token, expires_at, user_agent, ip, device_name)
VALUES (
    123, 
    'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...', 
    CURRENT_TIMESTAMP + INTERVAL '30 days',
    'Mozilla/5.0 (Windows NT 10.0; Win64; x64)...',
    '192.168.1.100',
    'Windows Chrome'
);
```

### Поиск активного токена
```sql
SELECT id, user_id, expires_at, device_name
FROM refresh_tokens
WHERE token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...'
  AND NOT is_revoked
  AND expires_at > CURRENT_TIMESTAMP;
```

### Получение всех активных сессий пользователя
```sql
SELECT id, device_name, ip, user_agent, created_at, expires_at
FROM refresh_tokens
WHERE user_id = 123
  AND NOT is_revoked
  AND expires_at > CURRENT_TIMESTAMP
ORDER BY created_at DESC;
```

### Отзыв токена (выход из сессии)
```sql
UPDATE refresh_tokens
SET is_revoked = TRUE, revoked_at = CURRENT_TIMESTAMP
WHERE id = 456;
```

### Отзыв всех токенов пользователя (выход из всех устройств)
```sql
UPDATE refresh_tokens
SET is_revoked = TRUE, revoked_at = CURRENT_TIMESTAMP
WHERE user_id = 123 AND NOT is_revoked;
```

### Удаление истекших токенов
```sql
SELECT cleanup_expired_refresh_tokens();
```

### Обновление токена (rotation)
```sql
-- Отзываем старый токен
UPDATE refresh_tokens 
SET is_revoked = TRUE, revoked_at = CURRENT_TIMESTAMP
WHERE token = 'old_token';

-- Создаем новый токен
INSERT INTO refresh_tokens (user_id, token, expires_at, user_agent, ip, device_name)
VALUES (123, 'new_token', CURRENT_TIMESTAMP + INTERVAL '30 days', 'Chrome...', '192.168.1.100', 'Chrome Desktop');
```

## Известные особенности

1. **Поле token изменилось с VARCHAR(255) на TEXT** - миграция 000042 увеличила размер для поддержки более длинных токенов.

2. **Поддержка IPv6** - поле ip имеет размер 45 символов для поддержки IPv6 адресов.

3. **Частичные индексы** - индексы созданы только для неотозванных токенов (WHERE NOT is_revoked) для оптимизации.

4. **Автоочистка** - функция cleanup_expired_refresh_tokens удаляет истекшие токены и отозванные токены старше 30 дней.

5. **Опциональная информация об устройстве** - поля user_agent, ip, device_name могут быть NULL.

## Интеграция с другими компонентами

1. **Система аутентификации** - основной компонент JWT auth flow

2. **API endpoints** - используется в /auth/refresh для обновления access токенов

3. **Middleware аутентификации** - проверяет валидность refresh токенов

4. **Управление сессиями** - позволяет пользователям видеть и управлять активными сессиями

## Безопасность

1. **Хранение токенов** - токены должны быть надежно зашифрованы или захешированы

2. **Защита от replay атак** - использование expires_at и is_revoked предотвращает повторное использование

3. **Логирование доступа** - рекомендуется логировать создание и использование токенов

4. **Ротация токенов** - периодическое обновление токенов повышает безопасность

5. **Ограничение количества** - можно ограничить количество активных токенов на пользователя

## Сценарии использования

### Авторизация через JWT
1. Пользователь вводит логин/пароль
2. Система создает access и refresh токены
3. Refresh токен сохраняется в базе
4. Access токен отправляется клиенту

### Обновление access токена
1. Клиент отправляет refresh токен
2. Система проверяет токен в базе
3. Если валиден - создает новый access токен
4. Опционально - создает новый refresh токен (rotation)

### Выход из системы
1. Клиент отправляет запрос на logout
2. Система отзывает refresh токен (is_revoked = TRUE)
3. Access токен становится недействительным

## Миграции

- **000041** - создание таблицы refresh_tokens с индексами и функцией очистки
- **000042** - увеличение размера поля token с VARCHAR(255) до TEXT