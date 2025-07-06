# Реализация нечеткого поиска в Backend

## Обзор

Реализована поддержка нечеткого поиска (fuzzy search) в backend части системы для улучшения поисковой функциональности marketplace.

## Выполненные изменения

### 1. Миграции базы данных

Применены миграции для добавления поддержки нечеткого поиска:

- **000084_add_text_search_extensions.up.sql** - Добавлены расширения PostgreSQL:
  - `pg_trgm` - для триграммного поиска
  - `unaccent` - для поиска без учета акцентов
  - Создана функция-обертка `f_unaccent` для использования в индексах
  - Добавлены GIN индексы для триграммного поиска по полям title и description

- **000085_create_search_synonyms.up.sql** - Создана таблица синонимов:
  - Таблица `search_synonyms` для хранения синонимов
  - Функция `expand_search_query` для расширения запросов синонимами
  - Предзаполнены базовые синонимы для русского, сербского и английского языков

- **000086_add_category_search_indexes.up.sql** - Добавлены индексы для категорий:
  - Индексы для нечеткого поиска по переводам категорий
  - Функция `search_categories_fuzzy` для поиска категорий с нечетким соответствием
  - Материализованное представление для статистики категорий

### 2. Backend реализация

#### Новые файлы:

1. **internal/proj/marketplace/service/search_query_expansion.go**
   - Функция `ExpandQueryWithSynonyms` - расширение запроса синонимами
   - Функция `SearchCategoriesFuzzy` - нечеткий поиск по категориям
   - Функция `AnalyzeSearchQuery` - анализ поискового запроса

2. **internal/storage/postgres/fuzzy_search.go**
   - Реализация методов Storage для работы с нечетким поиском

3. **internal/storage/postgres/database_fuzzy_search.go**
   - Реализация методов Database для работы с нечетким поиском

4. **internal/proj/marketplace/handler/fuzzy_search.go**
   - HTTP handlers для тестирования и использования нечеткого поиска
   - Endpoints: `/test-fuzzy-search` и `/fuzzy-search`

#### Обновленные файлы:

1. **internal/domain/search/types.go**
   - Добавлено поле `UseSynonyms` в структуры SearchParams и ServiceParams

2. **internal/storage/storage.go**
   - Добавлены методы интерфейса:
     - `ExpandSearchQuery` - расширение запроса синонимами
     - `SearchCategoriesFuzzy` - нечеткий поиск по категориям

3. **internal/proj/marketplace/storage/opensearch/repository.go**
   - Обновлен метод `buildSearchQuery` для поддержки:
     - Расширения запроса синонимами
     - Поиска по n-граммам для лучшего нечеткого соответствия

4. **internal/proj/marketplace/service/marketplace.go**
   - Обновлен метод `SearchListingsAdvanced` для передачи параметров нечеткого поиска

5. **internal/proj/marketplace/handler/search.go**
   - Добавлена обработка параметра `fuzzy` (по умолчанию включен)
   - Автоматическая установка `fuzziness=AUTO` при включенном нечетком поиске

6. **internal/proj/marketplace/handler/handler.go**
   - Добавлены новые маршруты:
     - `GET /api/v1/marketplace/test-fuzzy-search` - тестирование нечеткого поиска
     - `GET /api/v1/marketplace/fuzzy-search` - поиск с настраиваемыми параметрами

## Использование

### Основной поиск

По умолчанию нечеткий поиск включен для endpoint'а `/api/v1/marketplace/search`:

```bash
GET /api/v1/marketplace/search?query=квортира&fuzzy=true
```

### Тестирование нечеткого поиска

```bash
GET /api/v1/marketplace/test-fuzzy-search?query=машына&lang=ru
```

Возвращает:
- Оригинальный запрос
- Расширенный запрос с синонимами
- Найденные похожие категории

### Расширенный нечеткий поиск

```bash
GET /api/v1/marketplace/fuzzy-search?query=телефон&fuzziness=AUTO&use_synonyms=true&minimum_should_match=30%
```

Параметры:
- `fuzziness` - уровень нечеткости (AUTO, 0, 1, 2)
- `use_synonyms` - использовать синонимы (true/false)
- `minimum_should_match` - минимальное количество совпадений

## Технические детали

### PostgreSQL функции

1. **expand_search_query(query, language)** - расширяет запрос синонимами
2. **search_categories_fuzzy(search_term, language, similarity_threshold)** - ищет категории с нечетким соответствием
3. **f_unaccent(text)** - immutable обертка для функции unaccent

### OpenSearch улучшения

- Добавлена поддержка поиска по n-граммам (если настроен соответствующий analyzer)
- Расширение запроса синонимами через multi_match запросы
- Настраиваемый параметр fuzziness для match запросов

## Дальнейшие улучшения

Для полной функциональности рекомендуется:

1. Обновить маппинг OpenSearch для поддержки n-грамм анализаторов (скрипт `update_opensearch_mapping.sh`)
2. Переиндексировать существующие данные для применения новых анализаторов
3. Настроить веса и пороги схожести для оптимальных результатов
4. Расширить базу синонимов для вашей предметной области