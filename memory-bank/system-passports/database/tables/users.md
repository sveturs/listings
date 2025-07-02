# Паспорт таблицы: users

## Назначение
Основная таблица пользователей системы. Хранит данные аутентификации, профиля и настроек всех пользователей платформы Sve Tu.

## Структура таблицы

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) NOT NULL UNIQUE,
    google_id VARCHAR(255),
    picture_url TEXT,
    phone VARCHAR(20),
    bio TEXT,
    notification_email BOOLEAN DEFAULT true,
    timezone VARCHAR(50) DEFAULT 'UTC',
    last_seen TIMESTAMP,
    account_status VARCHAR(20) DEFAULT 'active',
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    password VARCHAR(255),
    provider VARCHAR(50) DEFAULT 'email',
    city VARCHAR(100),
    country VARCHAR(100),
    
    CONSTRAINT users_account_status_check CHECK (account_status IN ('active', 'inactive', 'suspended'))
);
```

## Поля таблицы

### Основные поля
- `id` - уникальный идентификатор пользователя (SERIAL)
- `name` - имя пользователя, обязательное поле (до 100 символов)
- `email` - email пользователя, уникальный, обязательный (до 150 символов)

### Аутентификация
- `google_id` - идентификатор Google для OAuth авторизации
- `password` - хеш пароля для email авторизации
- `provider` - способ регистрации: 'email', 'google' (по умолчанию 'email')

### Профиль
- `picture_url` - URL аватара пользователя
- `phone` - телефон пользователя (до 20 символов)
- `bio` - описание/о себе
- `city` - город пользователя
- `country` - страна пользователя

### Настройки и статус
- `notification_email` - получать ли email уведомления (по умолчанию true)
- `timezone` - часовой пояс пользователя (по умолчанию 'UTC')
- `account_status` - статус аккаунта: 'active', 'inactive', 'suspended'
- `settings` - JSONB поле для дополнительных настроек
- `last_seen` - время последней активности

### Системные поля
- `created_at` - дата создания записи
- `updated_at` - дата последнего обновления

## Индексы

1. **idx_users_phone** - индекс по телефону для быстрого поиска
2. **idx_users_status** - индекс по статусу аккаунта
3. **idx_users_email** - индекс по email
4. **idx_users_provider** - индекс по провайдеру авторизации
5. **idx_users_email_lower** - функциональный индекс по lowercase email
6. **idx_users_active** - частичный индекс для активных пользователей по last_seen

## Триггеры

- **update_users_updated_at** - автоматически обновляет поле updated_at при изменении записи

## Связи с другими таблицами

### Прямые связи (эта таблица ссылается на)
- Нет

### Обратные связи (другие таблицы ссылаются на users)
- `marketplace_listings.user_id` - объявления пользователя
- `marketplace_favorites.user_id` - избранные объявления
- `marketplace_chats.buyer_id` - чаты где пользователь покупатель
- `marketplace_chats.seller_id` - чаты где пользователь продавец
- `marketplace_messages.sender_id` - отправленные сообщения
- `reviews.reviewer_id` - написанные отзывы
- `reviews.reviewed_user_id` - отзывы о пользователе
- `user_balances.user_id` - баланс пользователя
- `user_contacts.user_id` - контакты пользователя
- `notifications.user_id` - уведомления пользователя
- `user_storefronts.user_id` - витрины пользователя

## Бизнес-правила

1. **Email уникальность** - email должен быть уникальным в системе
2. **Обязательные поля** - name и email обязательны для заполнения
3. **Статусы аккаунта**:
   - `active` - активный пользователь
   - `inactive` - неактивный (сам деактивировал)
   - `suspended` - заблокирован администрацией
4. **Провайдеры**:
   - `email` - регистрация через email/пароль
   - `google` - регистрация через Google OAuth

## Примеры использования

### Создание пользователя через email
```sql
INSERT INTO users (name, email, password, provider) 
VALUES ('Иван Иванов', 'ivan@example.com', '$2a$10$...hash...', 'email');
```

### Создание пользователя через Google OAuth
```sql
INSERT INTO users (name, email, google_id, picture_url, provider) 
VALUES ('Петр Петров', 'petr@gmail.com', '123456789', 'https://...', 'google');
```

### Поиск активных пользователей
```sql
SELECT * FROM users 
WHERE account_status = 'active' 
ORDER BY last_seen DESC;
```

### Обновление настроек
```sql
UPDATE users 
SET settings = settings || '{"theme": "dark", "language": "ru"}'::jsonb
WHERE id = 123;
```

## Известные особенности

1. **Google ID не уникален** - в миграции 000038 было удалено уникальное ограничение, что позволяет одному Google аккаунту иметь несколько пользователей
2. **Case-insensitive email** - есть индекс по LOWER(email) для поиска без учета регистра
3. **JSONB настройки** - поле settings позволяет хранить произвольные настройки без изменения схемы
4. **Два способа авторизации** - система поддерживает как email/password, так и Google OAuth
5. **Автообновление updated_at** - триггер автоматически обновляет это поле

## Миграции

- **000001** - создание таблицы
- **000014** - добавление city, country
- **000038** - добавление password, удаление UNIQUE с google_id
- **000039** - добавление оптимизационных индексов
- **000043** - добавление provider