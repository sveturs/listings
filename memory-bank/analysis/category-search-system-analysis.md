# Анализ системы категорий и поиска

## 1. Структура базы данных

### 1.1 Таблицы категорий

#### marketplace_categories
- **id**: integer (PK)
- **name**: varchar(100) - название категории
- **slug**: varchar(100) - уникальный slug
- **parent_id**: integer - родительская категория (иерархия)
- **level**: integer - уровень вложенности
- **sort_order**: integer - порядок сортировки
- **is_active**: boolean
- **seo_title**, **seo_description**, **seo_keywords**: SEO поля

**Индексы:**
- По slug (уникальный)
- По parent_id
- По external_id

#### category_attributes
- **id**: integer (PK)
- **name**: varchar(100) - системное имя атрибута
- **display_name**: varchar(255) - отображаемое имя
- **attribute_type**: varchar(50) - тип (text, select, number и т.д.)
- **options**: jsonb - опции для select
- **is_searchable**: boolean
- **is_filterable**: boolean
- **is_required**: boolean
- **sort_order**: integer

#### category_attribute_mapping
- **category_id**: integer (FK)
- **attribute_id**: integer (FK)
- **is_enabled**: boolean
- **is_required**: boolean
- **sort_order**: integer

Связывает категории с их атрибутами (many-to-many).

### 1.2 Поисковые таблицы

#### search_queries
- Хранит поисковые запросы пользователей
- **normalized_query**: нормализованный запрос
- **search_count**: счетчик использования
- **language**: язык запроса
- **results_count**: количество результатов

#### search_weights
- Веса для полей в поиске
- **field_name**: имя поля (title, description и т.д.)
- **weight**: вес от 0.0 до 1.0
- **search_type**: fulltext/fuzzy/exact
- **item_type**: marketplace/storefront/global
- **category_id**: опциональная связь с категорией

#### search_statistics
- Статистика поисковых запросов
- Время выполнения, количество результатов
- Связь с пользователем

#### user_behavior_events
- События поведения пользователей
- **event_type**: тип события (search, click, view и т.д.)
- **search_query**: поисковый запрос
- **item_id**, **item_type**: на что кликнули
- **position**: позиция в результатах

### 1.3 Переводы

#### translations
- **entity_type**: 'category' для категорий
- **entity_id**: ID категории
- **field_name**: 'name', 'description' и т.д.
- **language**: ru/en/sr
- **translated_text**: перевод

## 2. OpenSearch структура

### Индекс marketplace

#### Анализаторы:
- **serbian_analyzer**: с serbian stemmer
- **russian_analyzer**: с russian stemmer
- **english_analyzer**: с english stemmer
- **autocomplete**: edge_ngram для автодополнения
- **shingle_analyzer**: для фраз

#### Основные поля:
- **title**, **description**: с language-specific анализаторами
- **category**: объект с id, name, slug
- **category_id**: integer для фильтрации
- **category_path_ids**: массив ID для иерархии
- **attributes**: nested объект с атрибутами
- **all_attributes_text**: все атрибуты как текст для поиска

#### Поля атрибутов:
- **attribute_id**, **attribute_name**
- **text_value**: с language анализаторами
- **numeric_value**, **boolean_value**
- **display_value**: отображаемое значение
- **translations**: объект с переводами

## 3. Backend архитектура

### 3.1 Handlers

#### CategoriesHandler
- `GetCategories()`: получение списка категорий с переводами
- `GetCategoryTree()`: иерархическое дерево с кэшированием (5 минут)
- `GetCategoryAttributes()`: атрибуты конкретной категории

#### UnifiedSearchHandler
- Поиск по marketplace и storefront
- Поддержка фильтров по категориям
- Сортировка по релевантности/цене/дате/популярности
- Faceted search возможности

### 3.2 Сервисы и репозитории

- **CategoryService**: бизнес-логика категорий
- **PostgreSQL repository**: работа с БД
- **OpenSearch repository**: индексация и поиск
- **Redis**: кэширование категорий и атрибутов

### 3.3 Поисковая логика

Веса поиска из config:
- Точное совпадение title: 5.0
- Частичное совпадение title: 3.0
- Description: 2.0
- Атрибуты: 4.0-5.0
- Переводы: 1.5-4.0

## 4. Frontend

### 4.1 AI интеграция

Текущий подход:
1. AI анализирует изображение
2. Возвращает строку категории (например "automotive/auto-parts/tires-wheels")
3. Frontend функция `getCategoryData` пытается найти соответствие
4. Поиск по slug, name, переводам

### 4.2 Проблемы текущего подхода

1. **Жестко закодированные категории в промптах**
2. **Нет связи с реальными категориями из БД**
3. **Сложность поддержки при добавлении категорий**
4. **Ограничения по размеру промптов**

## 5. Существующие возможности для улучшения

### 5.1 Уже есть в системе:

1. **Таблица search_weights** - можно использовать для весов ключевых слов
2. **Таблица user_behavior_events** - для обучения на основе кликов
3. **Таблица search_queries** - популярные запросы
4. **OpenSearch с анализаторами** - для fuzzy matching
5. **Redis кэширование** - для быстрого доступа

### 5.2 Что можно добавить:

1. **Таблица category_keywords** для семантического поиска
2. **API endpoint** для умного определения категории
3. **Векторный поиск** в OpenSearch
4. **ML модель** на основе истории выборов

## 6. Выводы

Система уже имеет хорошую базу для реализации умного поиска категорий:
- Структура БД поддерживает иерархию и атрибуты
- OpenSearch настроен для мультиязычного поиска
- Есть таблицы для весов и статистики
- Backend готов к расширению

Основная задача - создать промежуточный слой между AI и категориями, который будет использовать семантический поиск вместо жестко закодированных правил.