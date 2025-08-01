# Инструкция по добавлению категорий и переводов в базу данных

## Важные правила

1. **ВСЕ изменения в базе данных должны производиться ТОЛЬКО через миграции!**
2. **НИКОГДА не выполняйте прямые SQL запросы для изменения данных**
3. **Обязательно создавайте как up, так и down миграции**

## Структура данных

### Таблица категорий (marketplace_categories)
```sql
- id: уникальный идентификатор
- parent_id: ID родительской категории (NULL для корневых)
- name: название на сербском языке (основное)
- slug: URL-совместимое название
- level: уровень вложенности (0 - корневые, 1 - первый уровень и т.д.)
- sort_order: порядок сортировки
- is_active: активность категории
```

### Таблица переводов (translations)
```sql
- entity_type: 'category' (НЕ 'marketplace_category'!)
- entity_id: ID категории
- field_name: 'name', 'seo_title', 'seo_description'
- language: 'en', 'ru', 'sr'
- translated_text: текст перевода
```

## Пошаговая инструкция

### 1. Создание миграции для добавления категорий

```sql
-- migrations/000XXX_add_new_categories.up.sql
INSERT INTO marketplace_categories (id, parent_id, name, slug, level, sort_order, is_active) VALUES
(1304, 1303, 'Gume i točkovi', 'tires-and-wheels', 1, 1, true),
(1305, 1303, 'Motor i delovi motora', 'engine-and-parts', 1, 2, true);

-- Обновление уровней если нужно
UPDATE marketplace_categories 
SET level = 1 
WHERE parent_id = 1303;
```

### 2. Создание миграции для добавления переводов

```sql
-- migrations/000XXX_add_category_translations.up.sql
INSERT INTO translations (entity_type, entity_id, field_name, language, translated_text) VALUES
-- Английские переводы
('category', 1304, 'name', 'en', 'Tires and Wheels'),
('category', 1304, 'seo_title', 'en', 'Tires, Wheels and Rims'),
('category', 1304, 'seo_description', 'en', 'Summer, winter and all-season tires'),

-- Русские переводы
('category', 1304, 'name', 'ru', 'Шины и колеса'),
('category', 1304, 'seo_title', 'ru', 'Шины, колеса и диски'),
('category', 1304, 'seo_description', 'ru', 'Летние, зимние и всесезонные шины'),

-- Сербские переводы (латиница)
('category', 1304, 'name', 'sr', 'Gume i točkovi'),
('category', 1304, 'seo_title', 'sr', 'Gume, točkovi i felne'),
('category', 1304, 'seo_description', 'sr', 'Letnje, zimske i celogodišnje gume')

-- Используйте ON CONFLICT для обновления существующих
ON CONFLICT (entity_type, entity_id, field_name, language) 
DO UPDATE SET translated_text = EXCLUDED.translated_text;
```

### 3. Создание down миграций

```sql
-- migrations/000XXX_add_new_categories.down.sql
DELETE FROM marketplace_categories 
WHERE id IN (1304, 1305);

-- migrations/000XXX_add_category_translations.down.sql
DELETE FROM translations 
WHERE entity_type = 'category' 
AND entity_id IN (1304, 1305);
```

### 4. Применение миграций

```bash
cd /data/hostel-booking-system/backend
./migrator up
```

### 5. Очистка кэша и переиндексация

```bash
# Очистка Redis кэша
docker exec hostel_redis redis-cli FLUSHALL

# Переиндексация OpenSearch
cd /data/hostel-booking-system/backend
./reindex
```

## Важные замечания

1. **entity_type для категорий**: Всегда используйте `'category'`, а НЕ `'marketplace_category'`
2. **Языки**: 
   - `en` - английский
   - `ru` - русский
   - `sr` - сербский (используйте латиницу, не кириллицу)
3. **Поля для перевода**:
   - `name` - название категории
   - `seo_title` - SEO заголовок
   - `seo_description` - SEO описание
4. **После изменений**: Обязательно очистите кэш Redis и переиндексируйте OpenSearch

## Проверка результатов

```sql
-- Проверка категорий
SELECT id, parent_id, name, slug FROM marketplace_categories 
WHERE parent_id = 1303 ORDER BY sort_order;

-- Проверка переводов
SELECT entity_id, language, field_name, translated_text 
FROM translations 
WHERE entity_type = 'category' AND entity_id = 1304 
ORDER BY language, field_name;
```

## Отладка проблем

Если переводы не отображаются:
1. Проверьте, что используется правильный entity_type ('category')
2. Очистите кэш Redis
3. Проверьте функции GetCategories и GetAllCategories в backend коде
4. Убедитесь, что backend перезапущен после изменения кода