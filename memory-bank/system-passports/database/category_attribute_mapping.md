# Паспорт таблицы: category_attribute_mapping

## Общая информация
- **Название**: category_attribute_mapping
- **Схема**: public
- **Описание**: Таблица связи многие-ко-многим между категориями товаров и их атрибутами
- **Тип данных**: Связывающая таблица (junction table)
- **Примерный размер**: 5K-20K записей
- **Частота обновления**: Низкая (конфигурация администраторами)

## Структура таблицы

### Поля

| Поле | Тип | Обязательность | Описание | Ограничения |
|------|-----|----------------|----------|-------------|
| category_id | INT | NOT NULL | ID категории товаров | FK -> marketplace_categories(id) ON DELETE CASCADE |
| attribute_id | INT | NOT NULL | ID атрибута | FK -> category_attributes(id) ON DELETE CASCADE |
| is_enabled | BOOLEAN | DEFAULT true | Активность атрибута для категории | - |
| is_required | BOOLEAN | DEFAULT false | Обязательность для данной категории | Переопределяет глобальную настройку |

### Составной первичный ключ
- **PRIMARY KEY**: (category_id, attribute_id)
- Обеспечивает уникальность: один атрибут привязывается к категории только один раз

## Индексы

| Название | Поля | Тип | Описание |
|----------|------|-----|----------|
| category_attribute_mapping_pkey | (category_id, attribute_id) | PRIMARY KEY | Составной первичный ключ |
| category_attribute_mapping_category_id_fkey | category_id | FOREIGN KEY | Связь с категорией |
| category_attribute_mapping_attribute_id_fkey | attribute_id | FOREIGN KEY | Связь с атрибутом |

### Рекомендуемые индексы (отсутствуют)
- INDEX на attribute_id для обратного поиска категорий по атрибуту
- INDEX на (is_enabled, is_required) для фильтрации

## Связи

### Внешние ключи (Foreign Keys)

| Поле | Ссылается на | Действие при удалении | Описание |
|------|--------------|----------------------|----------|
| category_id | marketplace_categories(id) | CASCADE | При удалении категории удаляются все связи |
| attribute_id | category_attributes(id) | CASCADE | При удалении атрибута удаляются все связи |

### Связанные таблицы
- **marketplace_categories**: Категории товаров
- **category_attributes**: Определения атрибутов

## Бизнес-логика

### Основные операции
1. **Привязка атрибута**: Администратор связывает атрибут с категорией
2. **Настройка обязательности**: Для каждой категории можно указать, какие атрибуты обязательны
3. **Активация/деактивация**: Временное отключение атрибута без удаления данных
4. **Массовые операции**: Привязка набора атрибутов ко всем подкатегориям

### Логика наследования
- **Наследование от родительских категорий**: Подкатегории автоматически получают атрибуты родителей
- **Переопределение**: Подкатегория может изменить настройки наследованного атрибута
- **Специфичные атрибуты**: Каждая категория может иметь уникальные атрибуты

### Приоритеты настроек
1. **Локальная настройка**: is_required в mapping таблице
2. **Глобальная настройка**: is_required в category_attributes
3. **Значение по умолчанию**: false

## Примеры конфигурации

### Автомобили (category_id = 2000)
```sql
-- Обязательные атрибуты для всех автомобилей
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required)
VALUES 
(2000, 1, true, true),  -- make (марка)
(2000, 2, true, true),  -- model (модель) 
(2000, 3, true, true);  -- year (год)

-- Опциональные атрибуты
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required)
VALUES
(2000, 4, true, false), -- mileage (пробег)
(2000, 5, true, false); -- engine_capacity (объем двигателя)
```

### Недвижимость (category_id = 1000)
```sql
-- Обязательные для недвижимости
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required)
VALUES
(1000, 10, true, true), -- property_type (тип недвижимости)
(1000, 15, true, true); -- area (площадь)

-- Специфичные для квартир (подкатегория)
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required)
VALUES
(1100, 11, true, false), -- rooms (количество комнат)
(1100, 12, true, false); -- floor (этаж)
```

## Автоматизация и массовые операции

### Автоматическое создание связей
Из миграции 000011 видно массовое создание связей:

```sql
-- Привязка атрибутов автомобилей ко всем подкатегориям
INSERT INTO category_attribute_mapping (category_id, attribute_id, is_enabled, is_required)
SELECT c.id, a.id, true, false
FROM marketplace_categories c
CROSS JOIN category_attributes a
WHERE (c.id = 2000 OR c.parent_id = 2000 OR EXISTS (
    SELECT 1 FROM marketplace_categories c2 
    WHERE c2.parent_id = 2000 AND c.parent_id = c2.id
))
AND a.name IN ('make', 'model', 'year', 'mileage', ...);
```

### Обновление обязательности
```sql
-- Делаем определенные атрибуты обязательными
UPDATE category_attribute_mapping 
SET is_required = true
WHERE attribute_id IN (
    SELECT id FROM category_attributes 
    WHERE name IN ('make', 'model', 'year')
)
AND category_id IN (/* автомобильные категории */);
```

## Валидация и ограничения

### Проверки целостности
- Нельзя создать объявление без заполнения обязательных атрибутов
- Деактивированные атрибуты (is_enabled = false) не отображаются в формах
- При удалении категории все связи автоматически удаляются

### Бизнес-правила
- Один атрибут не может быть привязан к категории дважды (PRIMARY KEY)
- Нельзя сделать атрибут обязательным, если он отключен
- Системные атрибуты нельзя отключать для корневых категорий

## Производительность

### Оптимизации
- Кеширование mapping для часто используемых категорий
- Предварительная загрузка атрибутов при отображении форм
- Индексирование по составному ключу

### Проблемы производительности
- JOIN с тремя таблицами замедляет загрузку форм
- Рекурсивные запросы для наследования от родительских категорий
- Большое количество атрибутов увеличивает размер форм

## API и интеграции

### Типичные запросы
```sql
-- Получить все атрибуты для категории
SELECT ca.*, cam.is_required, cam.is_enabled
FROM category_attributes ca
JOIN category_attribute_mapping cam ON ca.id = cam.attribute_id
WHERE cam.category_id = ? AND cam.is_enabled = true
ORDER BY ca.sort_order;

-- Получить обязательные атрибуты
SELECT ca.*
FROM category_attributes ca
JOIN category_attribute_mapping cam ON ca.id = cam.attribute_id
WHERE cam.category_id = ? AND cam.is_required = true;

-- Найти категории с определенным атрибутом
SELECT mc.*
FROM marketplace_categories mc
JOIN category_attribute_mapping cam ON mc.id = cam.category_id
WHERE cam.attribute_id = ? AND cam.is_enabled = true;
```

### Frontend интеграция
- Динамическая генерация форм на основе mapping
- Условная валидация полей
- Показ/скрытие атрибутов в зависимости от выбранной категории

## Мониторинг и аналитика

### Ключевые метрики
- Количество активных атрибутов на категорию
- Процент категорий с обязательными атрибутами
- Топ-10 самых используемых атрибутов
- Категории без привязанных атрибутов

### Алерты
- Категории с очень большим количеством атрибутов (>50)
- Критические атрибуты отключены для популярных категорий
- Нарушения целостности связей

## Изменения и миграции

### История изменений
1. **000011_create_category_attributes.up.sql**: Создание таблицы и первичные связи
2. **000031_update_category_attribute_mapping.up.sql**: Обновления структуры связей

### Планируемые улучшения
- Поддержка условных атрибутов (показываются при определенных значениях других)
- Версионирование mapping для A/B тестирования
- Автоматическое наследование настроек от родительских категорий
- Аналитика использования атрибутов для оптимизации форм