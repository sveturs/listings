# Паспорт таблицы: listing_attribute_values

## Общая информация
- **Название**: listing_attribute_values
- **Схема**: public
- **Описание**: Таблица для хранения фактических значений атрибутов конкретных объявлений
- **Тип данных**: Пользовательские данные, характеристики товаров
- **Примерный размер**: 100K-1M записей
- **Частота обновления**: Высокая (при создании/редактировании объявлений)

## Структура таблицы

### Поля

| Поле | Тип | Обязательность | Описание | Ограничения |
|------|-----|----------------|----------|-------------|
| id | SERIAL | PRIMARY KEY | Уникальный идентификатор записи | Автоинкремент |
| listing_id | INT | NOT NULL | ID объявления | FK -> marketplace_listings(id) ON DELETE CASCADE |
| attribute_id | INT | NOT NULL | ID атрибута | FK -> category_attributes(id) ON DELETE CASCADE |
| text_value | TEXT | NULL | Текстовое значение | Для типов text, select |
| numeric_value | NUMERIC(20,5) | NULL | Числовое значение | Для типов number, range |
| boolean_value | BOOLEAN | NULL | Логическое значение | Для типа boolean |
| json_value | JSONB | NULL | JSON значение | Для сложных типов, массивов |

## Полиморфная структура значений

### Принцип хранения
Таблица использует полиморфный подход - для каждого типа атрибута заполняется только одно поле:

| Тип атрибута | Используемое поле | Пример значения |
|--------------|-------------------|----------------|
| text | text_value | "BMW X5" |
| number | numeric_value | 2.5 |
| select | text_value | "Sedan" |
| boolean | boolean_value | true |
| multiselect | json_value | ["Option1", "Option2"] |
| range | json_value | {"min": 100, "max": 200} |

### Примеры значений

```sql
-- Автомобиль BMW X5, 2020 год, пробег 50000 км
INSERT INTO listing_attribute_values (listing_id, attribute_id, text_value) VALUES
(1, 1, 'BMW'),           -- make
(1, 2, 'X5'),            -- model
(1, 3, '2020');          -- year (как text)

INSERT INTO listing_attribute_values (listing_id, attribute_id, numeric_value) VALUES
(1, 4, 50000),           -- mileage
(1, 5, 3.0);             -- engine_capacity

INSERT INTO listing_attribute_values (listing_id, attribute_id, boolean_value) VALUES
(1, 15, true);           -- has_parking

-- Квартира с множественными опциями
INSERT INTO listing_attribute_values (listing_id, attribute_id, json_value) VALUES
(2, 20, '["Balcony", "Elevator", "Parking"]'); -- amenities
```

## Индексы

| Название | Поля | Тип | Описание |
|----------|------|-----|----------|
| listing_attribute_values_pkey | id | PRIMARY KEY | Первичный ключ |
| idx_listing_attr_listing_id | listing_id | INDEX | Быстрый поиск атрибутов объявления |
| idx_listing_attr_unique | (listing_id, attribute_id) | UNIQUE | Защита от дублей |
| idx_listing_attr_text | (attribute_id, text_value) | PARTIAL INDEX | Поиск по текстовым значениям |
| idx_listing_attr_numeric | (attribute_id, numeric_value) | PARTIAL INDEX | Поиск по числовым значениям |
| idx_listing_attr_boolean | (attribute_id, boolean_value) | PARTIAL INDEX | Поиск по логическим значениям |

### Частичные индексы
Индексы созданы с условием WHERE, что экономит место:
```sql
CREATE INDEX idx_listing_attr_text ON listing_attribute_values(attribute_id, text_value) 
WHERE text_value IS NOT NULL;
```

## Связи

### Внешние ключи (Foreign Keys)

| Поле | Ссылается на | Действие при удалении | Описание |
|------|--------------|----------------------|----------|
| listing_id | marketplace_listings(id) | CASCADE | При удалении объявления удаляются все атрибуты |
| attribute_id | category_attributes(id) | CASCADE | При удалении атрибута удаляются все значения |

### Связанные таблицы
- **marketplace_listings**: Объявления товаров
- **category_attributes**: Определения атрибутов

## Бизнес-логика

### Жизненный цикл значений
1. **Создание объявления**: Пользователь заполняет форму с атрибутами
2. **Валидация**: Проверка значений согласно rules из category_attributes
3. **Сохранение**: Создание записей в listing_attribute_values
4. **Индексирование**: Отправка данных в OpenSearch для поиска
5. **Обновление**: Изменение значений при редактировании объявления

### Правила заполнения
- Одно объявление + один атрибут = одна запись (UNIQUE индекс)
- Заполняется только одно поле value в зависимости от типа атрибута
- Обязательные атрибуты должны иметь значение
- Значения валидируются согласно rules и options из category_attributes

### Типовые операции

#### Создание значений атрибутов
```sql
-- При создании объявления автомобиля
INSERT INTO listing_attribute_values (listing_id, attribute_id, text_value)
SELECT 123, ca.id, 
    CASE ca.name 
        WHEN 'make' THEN 'Toyota'
        WHEN 'model' THEN 'Camry'
        WHEN 'fuel_type' THEN 'Petrol'
    END
FROM category_attributes ca
WHERE ca.name IN ('make', 'model', 'fuel_type');
```

#### Поиск по атрибутам
```sql
-- Найти все BMW с пробегом меньше 100000
SELECT DISTINCT l.*
FROM marketplace_listings l
JOIN listing_attribute_values lav1 ON l.id = lav1.listing_id
JOIN listing_attribute_values lav2 ON l.id = lav2.listing_id
JOIN category_attributes ca1 ON lav1.attribute_id = ca1.id
JOIN category_attributes ca2 ON lav2.attribute_id = ca2.id
WHERE ca1.name = 'make' AND lav1.text_value = 'BMW'
  AND ca2.name = 'mileage' AND lav2.numeric_value < 100000;
```

## Производительность

### Проблемы производительности
- **Фрагментация данных**: Много NULL значений в value полях
- **Сложные JOIN-ы**: Поиск по нескольким атрибутам требует множественных JOIN
- **Размер таблицы**: При 10K объявлений и 10 атрибутах = 100K записей

### Оптимизации
- **Частичные индексы**: Индексируются только заполненные значения
- **Денормализация**: Популярные атрибуты дублируются в marketplace_listings
- **Материализованные представления**: Предвычисленные агрегаты
- **Партиционирование**: Разделение по категориям товаров

### Альтернативные подходы
- **JSONB в listings**: Хранение всех атрибутов в одном JSON поле
- **EAV с типизацией**: Отдельные таблицы для каждого типа значений
- **NoSQL**: Использование MongoDB для атрибутов

## Валидация и целостность

### Проверки на уровне приложения
- Соответствие типа значения типу атрибута
- Проверка обязательных атрибутов для категории
- Валидация значений согласно rules
- Уникальность значений для unique атрибутов

### Примеры валидации
```javascript
// Валидация числового атрибута
function validateNumericAttribute(attributeId, value, rules) {
    if (rules.min && value < rules.min) return false;
    if (rules.max && value > rules.max) return false;
    if (rules.step && (value % rules.step) !== 0) return false;
    return true;
}

// Валидация select атрибута
function validateSelectAttribute(attributeId, value, options) {
    return options.values.includes(value);
}
```

## Интеграции

### OpenSearch синхронизация
При изменении значений атрибутов необходимо обновлять индекс поиска:

```sql
-- Trigger для автоматической синхронизации
CREATE OR REPLACE FUNCTION sync_listing_to_opensearch()
RETURNS TRIGGER AS $$
BEGIN
    -- Отправить сигнал на переиндексацию объявления
    PERFORM pg_notify('listing_updated', NEW.listing_id::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
```

### API endpoints
- GET /api/listings/{id}/attributes - Получить все атрибуты объявления
- PUT /api/listings/{id}/attributes - Обновить атрибуты объявления
- POST /api/search/attributes - Поиск по атрибутам

## Мониторинг и аналитика

### Ключевые метрики
- Среднее количество атрибутов на объявление
- Процент заполненности обязательных атрибутов
- Топ-10 самых используемых значений атрибутов
- Скорость поиска по атрибутам

### Проблемные запросы
- Поиск по множественным атрибутам (>5 JOIN)
- Фильтрация по числовым диапазонам без индексов
- Агрегация значений атрибутов в реальном времени

## Изменения и миграции

### История изменений
1. **000011_create_category_attributes.up.sql**: Создание базовой структуры
2. **000023_attribute_indices.up.sql**: Добавление оптимизированных индексов

### Планируемые улучшения
- Сжатие пустых значений
- Кеширование популярных комбинаций атрибутов
- Автоматическое извлечение атрибутов из описаний (ML)
- Версионирование значений для истории изменений
- Bulk операции для массового импорта товаров