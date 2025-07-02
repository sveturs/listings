# Паспорт таблицы: attribute_groups

## Общая информация
- **Название**: attribute_groups
- **Схема**: public
- **Описание**: Таблица для группировки атрибутов в логические блоки для улучшения UX
- **Тип данных**: Мета-данные, UI конфигурация
- **Примерный размер**: 20-100 записей
- **Частота обновления**: Очень низкая (настройка дизайнерами/UX)

## Структура таблицы

### Основная таблица attribute_groups

| Поле | Тип | Обязательность | Описание | Ограничения |
|------|-----|----------------|----------|-------------|
| id | SERIAL | PRIMARY KEY | Уникальный идентификатор группы | Автоинкремент |
| name | VARCHAR(100) | NOT NULL UNIQUE | Системное имя группы | snake_case, уникальное |
| display_name | VARCHAR(255) | NOT NULL | Отображаемое название группы | Локализуется |
| description | TEXT | NULL | Описание назначения группы | - |
| icon | VARCHAR(100) | NULL | Материальная иконка для группы | Material Icons name |
| sort_order | INT | DEFAULT 0 | Порядок отображения групп | Меньше = раньше |
| is_active | BOOLEAN | DEFAULT true | Активность группы | - |
| is_system | BOOLEAN | DEFAULT false | Системная группа (нельзя удалить) | - |
| created_at | TIMESTAMP | DEFAULT NOW() | Дата создания | - |
| updated_at | TIMESTAMP | DEFAULT NOW() | Дата обновления | Обновляется триггером |

### Связанная таблица attribute_group_items

| Поле | Тип | Обязательность | Описание | Ограничения |
|------|-----|----------------|----------|-------------|
| id | SERIAL | PRIMARY KEY | Уникальный идентификатор связи | Автоинкремент |
| group_id | INT | NOT NULL | ID группы | FK -> attribute_groups(id) ON DELETE CASCADE |
| attribute_id | INT | NOT NULL | ID атрибута | FK -> category_attributes(id) ON DELETE CASCADE |
| icon | VARCHAR(100) | NULL | Переопределенная иконка атрибута | Material Icons name |
| sort_order | INT | DEFAULT 0 | Порядок атрибута в группе | - |
| custom_display_name | VARCHAR(255) | NULL | Переопределенное название | Для специфичного контекста |
| visibility_condition | JSONB | NULL | Условия видимости атрибута | Логические правила |
| created_at | TIMESTAMP | DEFAULT NOW() | Дата создания связи | - |

### Связанная таблица category_attribute_groups

| Поле | Тип | Обязательность | Описание | Ограничения |
|------|-----|----------------|----------|-------------|
| id | SERIAL | PRIMARY KEY | Уникальный идентификатор | Автоинкремент |
| category_id | INT | NOT NULL | ID категории | FK -> marketplace_categories(id) ON DELETE CASCADE |
| group_id | INT | NOT NULL | ID группы | FK -> attribute_groups(id) ON DELETE CASCADE |
| component_id | INT | NULL | ID кастомного UI компонента | FK -> custom_ui_components(id) |
| sort_order | INT | DEFAULT 0 | Порядок группы в категории | - |
| is_active | BOOLEAN | DEFAULT true | Активность группы для категории | - |
| display_mode | VARCHAR(50) | DEFAULT 'list' | Режим отображения | 'list', 'grid', 'accordion', 'tabs' |
| collapsed_by_default | BOOLEAN | DEFAULT false | Свернута по умолчанию | - |
| configuration | JSONB | DEFAULT '{}' | Дополнительная конфигурация | Специфичные настройки |
| created_at | TIMESTAMP | DEFAULT NOW() | Дата создания | - |

## Предустановленные системные группы

Из миграции 000036 создаются стандартные группы:

| name | display_name | icon | sort_order | Назначение |
|------|-------------|------|------------|------------|
| basic_info | Основная информация | info | 1 | Базовые характеристики |
| technical_specs | Технические характеристики | settings | 2 | Технические параметры |
| dimensions | Размеры и вес | straighten | 3 | Физические параметры |
| appearance | Внешний вид | palette | 4 | Визуальные характеристики |
| condition_details | Состояние и комплектация | inventory_2 | 5 | Состояние товара |
| additional_info | Дополнительная информация | more_horiz | 6 | Прочие характеристики |

## Индексы

| Название | Поля | Тип | Описание |
|----------|------|-----|----------|
| attribute_groups_pkey | id | PRIMARY KEY | Первичный ключ |
| attribute_groups_name_key | name | UNIQUE | Уникальность системного имени |
| idx_attribute_groups_name | name | INDEX | Поиск по имени |
| idx_attribute_groups_active | is_active | INDEX | Фильтрация активных |
| idx_attribute_group_items_group | group_id | INDEX | Поиск элементов группы |
| idx_attribute_group_items_attribute | attribute_id | INDEX | Поиск групп атрибута |
| idx_category_attribute_groups_category | category_id | INDEX | Группы категории |
| idx_category_attribute_groups_group | group_id | INDEX | Категории группы |

## Связи

### Внешние ключи

| Поле | Ссылается на | Действие при удалении | Описание |
|------|--------------|----------------------|----------|
| attribute_group_items.group_id | attribute_groups(id) | CASCADE | Удаление элементов при удалении группы |
| attribute_group_items.attribute_id | category_attributes(id) | CASCADE | Удаление связи при удалении атрибута |
| category_attribute_groups.category_id | marketplace_categories(id) | CASCADE | Удаление при удалении категории |
| category_attribute_groups.group_id | attribute_groups(id) | CASCADE | Удаление при удалении группы |
| category_attribute_groups.component_id | custom_ui_components(id) | SET NULL | Сброс компонента при удалении |

## Бизнес-логика

### Архитектура группировки

```
marketplace_categories (Автомобили)
    ↓ связь через category_attribute_groups
attribute_groups (Основная информация) 
    ↓ связь через attribute_group_items
category_attributes (марка, модель, год)
    ↓ значения в 
listing_attribute_values (BMW, X5, 2020)
```

### Режимы отображения (display_mode)

#### list (по умолчанию)
- Простой вертикальный список атрибутов
- Подходит для небольших групп (до 5 атрибутов)

#### grid
- Сетка 2x2 или 3x3 для компактного отображения
- Удобно для атрибутов с короткими значениями

#### accordion
- Сворачиваемые блоки с заголовками групп
- Экономит место на странице
- Хорошо для мобильных устройств

#### tabs
- Вкладки для переключения между группами
- Подходит для категорий с многими группами

#### custom
- Использует кастомный UI компонент (component_id)
- Полная свобода в отображении

### Условная видимость (visibility_condition)

Атрибуты могут показываться/скрываться на основе значений других атрибутов:

```json
{
  "show_when": {
    "attribute_name": "fuel_type",
    "operator": "equals",
    "value": "Electric"
  }
}

{
  "hide_when": {
    "attribute_name": "property_type", 
    "operator": "in",
    "values": ["Land", "Garage"]
  }
}
```

## Представление v_attribute_groups_with_items

Создается материализованное представление для оптимизации запросов:

```sql
CREATE OR REPLACE VIEW v_attribute_groups_with_items AS
SELECT 
    ag.id AS group_id,
    ag.name AS group_name,
    ag.display_name AS group_display_name,
    ag.icon AS group_icon,
    ag.sort_order AS group_sort_order,
    agi.attribute_id,
    ca.name AS attribute_name,
    ca.display_name AS attribute_display_name,
    COALESCE(agi.custom_display_name, ca.display_name) AS final_display_name,
    COALESCE(agi.icon, ca.icon) AS final_icon
FROM attribute_groups ag
LEFT JOIN attribute_group_items agi ON ag.id = agi.group_id
LEFT JOIN category_attributes ca ON agi.attribute_id = ca.id
WHERE ag.is_active = true
ORDER BY ag.sort_order, agi.sort_order;
```

## Триггеры

| Название | Событие | Функция | Описание |
|----------|---------|---------|----------|
| update_attribute_groups_updated_at | BEFORE UPDATE | update_attribute_groups_updated_at() | Обновление поля updated_at |

## API интеграция

### Типичные запросы

```sql
-- Получить все группы для категории с атрибутами
SELECT 
    ag.name,
    ag.display_name,
    ag.icon,
    cag.display_mode,
    cag.collapsed_by_default,
    json_agg(
        json_build_object(
            'attribute_id', ca.id,
            'name', ca.name,
            'display_name', COALESCE(agi.custom_display_name, ca.display_name),
            'type', ca.attribute_type,
            'icon', COALESCE(agi.icon, ca.icon),
            'visibility_condition', agi.visibility_condition
        ) ORDER BY agi.sort_order
    ) as attributes
FROM attribute_groups ag
JOIN category_attribute_groups cag ON ag.id = cag.group_id
JOIN attribute_group_items agi ON ag.id = agi.group_id  
JOIN category_attributes ca ON agi.attribute_id = ca.id
WHERE cag.category_id = ? AND cag.is_active = true
GROUP BY ag.id, ag.sort_order, cag.display_mode, cag.collapsed_by_default
ORDER BY ag.sort_order;
```

### Frontend использование

```javascript
// Генерация формы на основе групп
function generateAttributeForm(categoryId) {
    const groups = await fetchAttributeGroups(categoryId);
    
    return groups.map(group => ({
        title: group.display_name,
        icon: group.icon,
        display_mode: group.display_mode,
        collapsed: group.collapsed_by_default,
        attributes: group.attributes.map(attr => ({
            name: attr.name,
            label: attr.display_name,
            type: attr.type,
            icon: attr.icon,
            visible: evaluateVisibility(attr.visibility_condition, formData)
        }))
    }));
}
```

## Производительность

### Оптимизации
- Кеширование конфигурации групп на уровне приложения
- Предварительная загрузка всех групп для категории
- Минимизация JOIN-ов через денормализацию

### Рекомендации
- Не более 8 групп на категорию (UX ограничение)
- Не более 12 атрибутов в группе
- Кеширование результатов запросов на 24 часа

## Изменения и миграции

### История изменений
1. **000036_create_attribute_groups.up.sql**: Создание системы группировки атрибутов

### Планируемые улучшения
- Drag&Drop для изменения порядка атрибутов в админке
- Предпросмотр группировки перед сохранением
- A/B тестирование различных способов группировки
- Автоматические рекомендации группировки на основе аналитики
- Шаблоны группировки для импорта в новые категории