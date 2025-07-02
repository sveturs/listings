# Паспорт таблицы `user_storefronts`

## Назначение
Упрощенная версия витрин магазинов (deprecated/legacy). Более простая структура по сравнению с таблицей `storefronts`. Используется для обратной совместимости с существующими модулями импорта.

## Полная структура таблицы

```sql
CREATE TABLE IF NOT EXISTS user_storefronts (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    logo_path VARCHAR(255),
    slug VARCHAR(100) UNIQUE,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    creation_transaction_id INT REFERENCES balance_transactions(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Описание полей

| Поле | Тип | Обязательность | Описание |
|------|-----|---------------|----------|
| `id` | SERIAL | PRIMARY KEY | Уникальный ID витрины |
| `user_id` | INT | NOT NULL FK | ID владельца витрины |
| `name` | VARCHAR(255) | NOT NULL | Название витрины |
| `description` | TEXT | NULLABLE | Описание витрины |
| `logo_path` | VARCHAR(255) | NULLABLE | Путь к логотипу (локальный файл) |
| `slug` | VARCHAR(100) | UNIQUE NULLABLE | URL-slug для витрины |
| `status` | VARCHAR(20) | NOT NULL | Статус витрины |
| `creation_transaction_id` | INT | NULLABLE FK | ID транзакции создания |
| `created_at` | TIMESTAMP | DEFAULT NOW | Дата создания |
| `updated_at` | TIMESTAMP | DEFAULT NOW | Дата обновления |

## Индексы

```sql
CREATE INDEX IF NOT EXISTS idx_user_storefronts_user ON user_storefronts(user_id);
CREATE INDEX IF NOT EXISTS idx_user_storefronts_status ON user_storefronts(status);
```

## Ограничения

- **UNIQUE**: `slug` - уникальный URL (но может быть NULL)
- **FOREIGN KEY**: `user_id` → `users(id)`
- **FOREIGN KEY**: `creation_transaction_id` → `balance_transactions(id)`

## Связи с другими таблицами

| Связь | Тип | Описание |
|-------|-----|----------|
| `user_id` → `users.id` | Many-to-One | Владелец витрины |
| `creation_transaction_id` → `balance_transactions.id` | Many-to-One | Транзакция создания |
| `import_sources` | One-to-Many | Источники импорта товаров |

## Статусы витрины

| Статус | Описание |
|--------|----------|
| `active` | Активная витрина |
| `inactive` | Неактивная витрина |
| `suspended` | Заблокированная витрина |
| `pending` | Ожидает активации |

## Бизнес-правила

1. **Простота**: Минимальный набор полей для базовой функциональности
2. **Монетизация**: Связь с балансовыми транзакциями
3. **Миграция**: Постепенный переход на таблицу `storefronts`
4. **Обратная совместимость**: Поддержка существующих модулей импорта

## Примеры использования

### Создание простой витрины
```sql
INSERT INTO user_storefronts (user_id, name, description, slug)
VALUES (123, 'Мой магазин', 'Описание магазина', 'my-shop');
```

### Поиск витрин пользователя
```sql
SELECT * FROM user_storefronts 
WHERE user_id = 123 AND status = 'active'
ORDER BY created_at DESC;
```

### Деактивация витрины
```sql
UPDATE user_storefronts 
SET status = 'inactive', updated_at = CURRENT_TIMESTAMP
WHERE id = 42;
```

## Известные особенности

1. **Legacy система**: Старая версия витрин
2. **Используется импортом**: Модули импорта ссылаются на эту таблицу
3. **Платное создание**: Связь с балансовыми транзакциями
4. **Локальные файлы**: `logo_path` для локального хранения логотипов
5. **Простая структура**: Нет поддержки брендинга, SEO, аналитики

## Отличия от таблицы `storefronts`

| Функция | user_storefronts | storefronts |
|---------|------------------|-------------|
| Брендинг | Только логотип | Полная поддержка |
| SEO | Нет | Есть |
| Аналитика | Нет | Есть |
| Подписки | Нет | Есть |
| AI-функции | Нет | Есть |
| Геолокация | Нет | Есть |
| Статус | Простой enum | Богатые статусы |

## Использование в коде

**Backend**:
- Модули импорта: `internal/proj/import/`
- Legacy handlers: `internal/proj/user_storefronts/`
- Миграция: `internal/proj/storefront_migration/`

**Frontend**:
- Legacy компоненты: `src/components/legacy/storefront/`
- Админка: `src/components/admin/storefronts/`

## Стратегия миграции

1. **Фаза 1**: Новые витрины создаются в таблице `storefronts`
2. **Фаза 2**: Миграция существующих витрин
3. **Фаза 3**: Обновление модулей импорта
4. **Фаза 4**: Deprecation таблицы `user_storefronts`

## Связанные компоненты

- **Import sources**: Источники импорта товаров
- **Import history**: История импорта
- **Balance transactions**: Система оплаты создания витрин
- **Storefront migration**: Система миграции на новую структуру