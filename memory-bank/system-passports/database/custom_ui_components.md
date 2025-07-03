# Паспорт таблицы: custom_ui_components

## Общая информация
- **Название**: custom_ui_components
- **Схема**: public
- **Описание**: Система хранения и управления кастомными UI компонентами для категорий и атрибутов
- **Тип данных**: Код компонентов, мета-данные UI
- **Примерный размер**: 50-500 записей
- **Частота обновления**: Низкая (разработка новых компонентов)

## Структура таблицы

### Основная таблица custom_ui_components

| Поле | Тип | Обязательность | Описание | Ограничения |
|------|-----|----------------|----------|-------------|
| id | SERIAL | PRIMARY KEY | Уникальный идентификатор | Автоинкремент |
| name | VARCHAR(255) | NOT NULL UNIQUE | Системное имя компонента | Уникальное |
| display_name | VARCHAR(255) | NOT NULL | Человекочитаемое название | - |
| description | TEXT | NULL | Описание назначения компонента | - |
| component_type | VARCHAR(50) | NOT NULL | Тип компонента | CHECK: 'category', 'attribute', 'filter' |
| component_code | TEXT | NOT NULL | JSX код компонента | React/JSX синтаксис |
| configuration | JSONB | DEFAULT '{}' | Конфигурация компонента | Пропсы и настройки |
| dependencies | JSONB | DEFAULT '[]' | Внешние зависимости | Список библиотек |
| is_active | BOOLEAN | DEFAULT true | Активность компонента | - |
| created_at | TIMESTAMP WITH TIME ZONE | DEFAULT NOW() | Дата создания | С учетом часового пояса |
| updated_at | TIMESTAMP WITH TIME ZONE | DEFAULT NOW() | Дата обновления | Обновляется триггером |
| created_by | INT | NULL | Создатель компонента | FK -> users(id) |
| updated_by | INT | NULL | Последний редактор | FK -> users(id) |

### Дополнительные поля (из последующих миграций)

| Поле | Тип | Описание | Назначение |
|------|-----|----------|------------|
| compiled_code | TEXT | Транспилированный JS код | Кеширование компиляции |
| compilation_errors | JSONB | Ошибки компиляции | Отладочная информация |
| last_compiled_at | TIMESTAMP WITH TIME ZONE | Время последней компиляции | Для инвалидации кеша |

### Связанная таблица custom_ui_component_usage

| Поле | Тип | Обязательность | Описание | Ограничения |
|------|-----|----------------|----------|-------------|
| id | SERIAL | PRIMARY KEY | Уникальный идентификатор | Автоинкремент |
| component_id | INT | NOT NULL | ID компонента | FK -> custom_ui_components(id) ON DELETE CASCADE |
| entity_type | VARCHAR(50) | NOT NULL | Тип сущности | CHECK: 'category', 'attribute' |
| entity_id | INT | NOT NULL | ID сущности | category_id или attribute_id |
| configuration | JSONB | DEFAULT '{}' | Специфичная конфигурация | Переопределения для использования |
| created_at | TIMESTAMP WITH TIME ZONE | DEFAULT NOW() | Дата создания связи | - |
| updated_at | TIMESTAMP WITH TIME ZONE | DEFAULT NOW() | Дата обновления | Обновляется триггером |

### Связанная таблица custom_ui_templates

| Поле | Тип | Обязательность | Описание | Ограничения |
|------|-----|----------------|----------|-------------|
| id | SERIAL | PRIMARY KEY | Уникальный идентификатор | Автоинкремент |
| name | VARCHAR(255) | NOT NULL UNIQUE | Имя шаблона | Уникальное |
| display_name | VARCHAR(255) | NOT NULL | Отображаемое название | - |
| description | TEXT | NULL | Описание шаблона | - |
| template_code | TEXT | NOT NULL | Код шаблона компонента | С плейсхолдерами |
| template_type | VARCHAR(50) | NOT NULL | Тип шаблона | CHECK: 'category', 'attribute', 'filter' |
| example_configuration | JSONB | DEFAULT '{}' | Пример конфигурации | Для документации |
| is_active | BOOLEAN | DEFAULT true | Активность шаблона | - |
| created_at | TIMESTAMP WITH TIME ZONE | DEFAULT NOW() | Дата создания | - |
| updated_at | TIMESTAMP WITH TIME ZONE | DEFAULT NOW() | Дата обновления | - |

## Типы компонентов (component_type)

### category
- Компоненты для отображения информации о категории
- Кастомные лейауты страниц категорий
- Специальные фильтры и сортировки

### attribute  
- Кастомные поля ввода для атрибутов
- Специализированные виджеты (карты, слайдеры, галереи)
- Условная логика отображения

### filter
- Продвинутые фильтры для поиска
- Фасетная навигация
- Интерактивные элементы фильтрации

## Индексы

| Название | Поля | Тип | Описание |
|----------|------|-----|----------|
| custom_ui_components_pkey | id | PRIMARY KEY | Первичный ключ |
| custom_ui_components_name_key | name | UNIQUE | Уникальность имени |
| idx_custom_ui_components_name | name | INDEX | Поиск по имени |
| idx_custom_ui_components_type | component_type | INDEX | Фильтрация по типу |
| idx_custom_ui_components_active | is_active | INDEX | Фильтрация активных |
| idx_component_usage_component | component_id | INDEX | Поиск использований |
| idx_component_usage_entity | (entity_type, entity_id) | INDEX | Поиск по сущности |
| idx_component_usage_unique | (component_id, entity_type, entity_id) | UNIQUE | Уникальность использования |

## Связи

### Внешние ключи

| Поле | Ссылается на | Действие при удалении | Описание |
|------|--------------|----------------------|----------|
| created_by | users(id) | SET NULL | Автор компонента |
| updated_by | users(id) | SET NULL | Редактор компонента |
| custom_ui_component_usage.component_id | custom_ui_components(id) | CASCADE | Удаление использований |

### Обратные связи

| Таблица | Поле | Описание |
|---------|------|----------|
| category_attribute_groups | component_id | Использование в группах атрибутов |
| custom_ui_component_usage | component_id | Отслеживание использования |

## Структура JSX компонентов

### Пример компонента для атрибута
```jsx
// component_code для кастомного слайдера цены
import React, { useState } from 'react';
import { Slider } from '@mui/material';

const PriceRangeSlider = ({ 
    attribute, 
    value, 
    onChange, 
    configuration 
}) => {
    const { min = 0, max = 1000000, step = 1000 } = configuration;
    const [range, setRange] = useState(value || [min, max]);
    
    const handleChange = (event, newValue) => {
        setRange(newValue);
        onChange(newValue);
    };
    
    return (
        <div className="price-range-slider">
            <h3>{attribute.display_name}</h3>
            <Slider
                value={range}
                onChange={handleChange}
                valueLabelDisplay="auto"
                min={min}
                max={max}
                step={step}
                marks={[
                    { value: min, label: `${min}₽` },
                    { value: max, label: `${max}₽` }
                ]}
            />
            <div className="range-display">
                {range[0]}₽ - {range[1]}₽
            </div>
        </div>
    );
};

export default PriceRangeSlider;
```

### Пример конфигурации
```json
{
  "min": 0,
  "max": 5000000,
  "step": 50000,
  "currency": "RUB",
  "showMarks": true,
  "valueLabelFormat": "currency"
}
```

## Система компиляции

### Процесс компиляции
1. **Валидация синтаксиса**: Проверка JSX на корректность
2. **Транспиляция**: Babel преобразует JSX в обычный JS
3. **Минификация**: Сжатие кода для продакшена
4. **Кеширование**: Сохранение в поле compiled_code

### Обработка ошибок
```json
{
  "compilation_errors": [
    {
      "line": 15,
      "column": 23,
      "message": "Unexpected token '<'",
      "type": "SyntaxError"
    }
  ]
}
```

## Система шаблонов

### Шаблон для атрибута-селектора
```jsx
// template_code с плейсхолдерами
import React from 'react';

const {{COMPONENT_NAME}} = ({ 
    attribute, 
    value, 
    onChange, 
    options = {{DEFAULT_OPTIONS}} 
}) => {
    return (
        <select 
            value={value} 
            onChange={(e) => onChange(e.target.value)}
            className="{{CSS_CLASSES}}"
        >
            <option value="">{{PLACEHOLDER_TEXT}}</option>
            {options.map(option => (
                <option key={option.value} value={option.value}>
                    {option.label}
                </option>
            ))}
        </select>
    );
};

export default {{COMPONENT_NAME}};
```

### Генерация из шаблона
```javascript
function generateComponentFromTemplate(template, config) {
    return template.template_code
        .replace(/\{\{COMPONENT_NAME\}\}/g, config.componentName)
        .replace(/\{\{DEFAULT_OPTIONS\}\}/g, JSON.stringify(config.options))
        .replace(/\{\{CSS_CLASSES\}\}/g, config.cssClasses)
        .replace(/\{\{PLACEHOLDER_TEXT\}\}/g, config.placeholder);
}
```

## Триггеры

| Название | Событие | Функция | Описание |
|----------|---------|---------|----------|
| update_custom_ui_components_updated_at | BEFORE UPDATE | update_updated_at_column() | Обновление поля updated_at |
| update_custom_ui_component_usage_updated_at | BEFORE UPDATE | update_updated_at_column() | Обновление поля updated_at |
| update_custom_ui_templates_updated_at | BEFORE UPDATE | update_updated_at_column() | Обновление поля updated_at |

## Безопасность

### Ограничения безопасности
- Запрет на выполнение eval() в component_code
- Санитизация JSX кода перед компиляцией
- Ограничение доступных API в runtime
- Изоляция компонентов в sandbox

### Права доступа
- **Создание**: Только разработчики и администраторы
- **Редактирование**: Автор компонента + администраторы
- **Использование**: Все пользователи (если is_active = true)
- **Компиляция**: Автоматически при сохранении

## Производительность

### Оптимизации
- Lazy loading компонентов на frontend
- Кеширование скомпилированного кода
- CDN для статических зависимостей
- Code splitting по типам компонентов

### Мониторинг
- Время компиляции компонентов
- Частота использования компонентов
- Ошибки рендеринга в production
- Размер bundle после добавления компонентов

## API endpoints

### Основные операции
```
GET    /api/admin/ui-components              - Список всех компонентов
POST   /api/admin/ui-components              - Создание компонента
GET    /api/admin/ui-components/{id}         - Получение компонента
PUT    /api/admin/ui-components/{id}         - Обновление компонента
DELETE /api/admin/ui-components/{id}         - Удаление компонента
POST   /api/admin/ui-components/{id}/compile - Принудительная компиляция
GET    /api/ui-components/render/{name}      - Рендеринг компонента
```

### Публичные endpoints
```
GET /api/public/ui-components?type=attribute&entity_id=123
GET /api/public/ui-templates?type=category
```

## Изменения и миграции

### История изменений
1. **000033_create_custom_ui_components.up.sql**: Создание базовой системы
2. **000034_update_custom_ui_components.up.sql**: Добавление полей компиляции
3. **000037_update_custom_ui_components.up.sql**: Обновление структуры

### Планируемые улучшения
- Visual компонент-билдер с drag&drop
- Версионирование компонентов
- A/B тестирование UI компонентов
- Marketplace компонентов от сообщества
- Hot reload в development режиме
- TypeScript поддержка для компонентов