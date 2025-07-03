# Паспорт таблицы: marketplace_favorites

## Назначение
Хранение избранных объявлений пользователей. Позволяет пользователям сохранять интересующие их товары для быстрого доступа.

## Структура таблицы

```sql
CREATE TABLE marketplace_favorites (
    user_id INT REFERENCES users(id),
    listing_id INT REFERENCES marketplace_listings(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, listing_id)
);
```

## Поля таблицы

### Основные поля
- `user_id` - идентификатор пользователя (FK к users)
- `listing_id` - идентификатор объявления (FK к marketplace_listings)

### Системные поля
- `created_at` - дата добавления в избранное

### Составной первичный ключ
- `PRIMARY KEY (user_id, listing_id)` - гарантирует уникальность пары пользователь-объявление

## Индексы

1. **PRIMARY KEY** - составной индекс по (user_id, listing_id)
2. **idx_marketplace_favorites_user_count** - индекс для подсчета избранного у пользователя
3. **idx_marketplace_favorites_listing** - индекс для подсчета добавлений объявления в избранное

## Связи с другими таблицами

### Прямые связи (эта таблица ссылается на)
- `user_id` → `users.id` - пользователь
- `listing_id` → `marketplace_listings.id` - объявление (CASCADE DELETE)

### Обратные связи
- Нет прямых ссылок из других таблиц

## Бизнес-правила

### Ограничения
1. **Уникальность** - пользователь может добавить объявление в избранное только один раз
2. **Каскадное удаление** - при удалении объявления автоматически удаляются все записи из избранного
3. **Существование пользователя** - пользователь должен существовать в системе

### Логика работы
1. Пользователь не может добавить свои объявления в избранное (проверка на уровне приложения)
2. Нет ограничений на количество избранных объявлений
3. Избранное доступно только авторизованным пользователям

## Примеры использования

### Добавление в избранное
```sql
INSERT INTO marketplace_favorites (user_id, listing_id) 
VALUES (123, 456)
ON CONFLICT (user_id, listing_id) DO NOTHING;
```

### Удаление из избранного
```sql
DELETE FROM marketplace_favorites 
WHERE user_id = 123 AND listing_id = 456;
```

### Получение избранных объявлений пользователя
```sql
SELECT l.*, mi.public_url as main_image_url
FROM marketplace_favorites f
JOIN marketplace_listings l ON f.listing_id = l.id
LEFT JOIN marketplace_images mi ON l.id = mi.listing_id AND mi.is_main = true
WHERE f.user_id = 123
  AND l.status = 'active'
ORDER BY f.created_at DESC;
```

### Проверка, добавлено ли в избранное
```sql
SELECT EXISTS(
    SELECT 1 FROM marketplace_favorites 
    WHERE user_id = 123 AND listing_id = 456
) as is_favorite;
```

### Подсчет количества добавлений в избранное для объявления
```sql
SELECT COUNT(*) as favorites_count
FROM marketplace_favorites
WHERE listing_id = 456;
```

### Популярные объявления по количеству добавлений в избранное
```sql
SELECT l.*, COUNT(f.user_id) as favorites_count
FROM marketplace_listings l
LEFT JOIN marketplace_favorites f ON l.id = f.listing_id
WHERE l.status = 'active'
GROUP BY l.id
ORDER BY favorites_count DESC
LIMIT 10;
```

## API интеграция

### Endpoints
- `GET /api/v1/marketplace/favorites` - список избранных
- `POST /api/v1/marketplace/favorites` - добавить в избранное
- `DELETE /api/v1/marketplace/favorites/{listing_id}` - удалить из избранного

### Ответ API при получении объявлений
```json
{
  "id": 456,
  "title": "iPhone 13 Pro",
  "price": 899.99,
  "is_favorite": true,
  "favorites_count": 25
}
```

## Известные особенности

1. **Составной PRIMARY KEY** - автоматически создает уникальный индекс
2. **CASCADE DELETE** - автоматическая очистка при удалении объявления
3. **ON CONFLICT DO NOTHING** - безопасное добавление без ошибок дубликата
4. **Нет обратной связи** - удаление пользователя требует ручной очистки
5. **Оптимизированные индексы** - для быстрого подсчета и выборки

## Миграции

- **000001** - создание таблицы
- **000010** - добавление CASCADE DELETE для listing_id
- **000039** - добавление оптимизационных индексов

## Интеграция с другими компонентами

### Frontend
1. **Кнопка "В избранное"** - toggle состояние
2. **Страница избранного** - список сохраненных объявлений
3. **Индикатор в листинге** - отметка избранных

### Backend
1. **Middleware проверки** - только для авторизованных
2. **Статистика** - подсчет популярности
3. **Уведомления** - об изменении цены избранных товаров

### Аналитика
1. **Популярные товары** - по количеству добавлений
2. **Конверсия** - из избранного в покупку
3. **Retention** - возвращение к избранным

## Оптимизация производительности

1. **Пагинация** - обязательна для больших списков
2. **Кеширование** - счетчики favorites_count
3. **Batch операции** - массовая проверка is_favorite
4. **Индексы** - покрывают все частые запросы