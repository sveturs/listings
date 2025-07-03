# Паспорт таблицы `admin_users`

## Назначение
Реестр административных пользователей системы. Хранит email-адреса пользователей, которые имеют права администратора для доступа к админ-панели и выполнения привилегированных операций.

## Полная структура таблицы

```sql
CREATE TABLE IF NOT EXISTS admin_users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by INTEGER REFERENCES users(id),
    notes TEXT
);
```

## Описание полей

| Поле | Тип | Обязательность | Описание |
|------|-----|---------------|----------|
| `id` | SERIAL | PRIMARY KEY | Уникальный идентификатор записи |
| `email` | VARCHAR(255) | NOT NULL UNIQUE | Email администратора (уникальный) |
| `created_at` | TIMESTAMP WITH TIME ZONE | DEFAULT NOW() | Время добавления в систему |
| `created_by` | INTEGER | NULLABLE FK | ID пользователя, который добавил админа |
| `notes` | TEXT | NULLABLE | Дополнительные заметки об администраторе |

## Индексы

```sql
-- Быстрый поиск по email (главный use case)
CREATE INDEX IF NOT EXISTS admin_users_email_idx ON admin_users(email);
```

## Ограничения

- **UNIQUE**: `email` - один email может быть только у одного админа
- **FOREIGN KEY**: `created_by` → `users(id)` - кто добавил админа

## Связи с другими таблицами

| Связь | Тип | Описание |
|-------|-----|----------|
| `created_by` → `users.id` | Many-to-One | Кто добавил этого админа |

## Бизнес-правила

1. **Уникальность email**: Один email = один админ
2. **Проверка прав**: При аутентификации проверяется наличие email в этой таблице
3. **Логирование**: Фиксируется кто и когда добавил админа
4. **Заметки**: Можно оставлять комментарии о назначении админа

## Начальные данные

```sql
-- Системные администраторы
INSERT INTO admin_users (email, notes) 
VALUES 
    ('bevzenko.sergey@gmail.com', 'Added in migration 000028'),
    ('voroshilovdo@gmail.com', 'Added in migration 000028')
ON CONFLICT (email) DO NOTHING;
```

## Примеры использования

### Проверка прав администратора
```sql
SELECT EXISTS(
    SELECT 1 FROM admin_users 
    WHERE email = 'user@example.com'
) AS is_admin;
```

### Добавление нового администратора
```sql
INSERT INTO admin_users (email, created_by, notes)
VALUES ('newadmin@example.com', 1, 'Added by main admin');
```

### Получение списка всех админов
```sql
SELECT a.email, a.created_at, u.name as created_by_name, a.notes
FROM admin_users a
LEFT JOIN users u ON a.created_by = u.id
ORDER BY a.created_at DESC;
```

### Удаление прав администратора
```sql
DELETE FROM admin_users WHERE email = 'former_admin@example.com';
```

## Известные особенности

1. **Простая структура**: Минималистичная таблица только для проверки прав
2. **Email-based**: Авторизация основана на email, не на user_id
3. **Soft admin system**: Нет ролей или уровней - есть права или нет
4. **Audit trail**: Сохраняется информация о том, кто добавил админа

## Использование в коде

**Backend**:
- Middleware для проверки прав администратора
- Handler для управления списком админов
- Проверка при доступе к admin endpoints

```go
// Пример middleware проверки
func AdminRequired(c *fiber.Ctx) error {
    email := getUserEmail(c) // из JWT/session
    
    var exists bool
    db.Raw("SELECT EXISTS(SELECT 1 FROM admin_users WHERE email = ?)", email).Scan(&exists)
    
    if !exists {
        return c.Status(403).JSON(fiber.Map{"error": "admin.accessDenied"})
    }
    
    return c.Next()
}
```

**Frontend**:
- Проверка прав перед отображением admin UI
- Редирект неавторизованных пользователей
- Условный рендеринг admin-функций

## Security considerations

1. **RBAC**: Простейшая модель - есть права или нет
2. **Аудит**: Логируется кто добавил админа
3. **Revocation**: Удаление из таблицы = отзыв прав
4. **Email verification**: Важно проверять email при добавлении

## Связанные компоненты

- **Auth middleware**: Проверка админских прав
- **Admin handlers**: CRUD операции для управления админами  
- **Frontend admin routes**: Роуты только для админов
- **User management**: Связь с основной таблицей пользователей