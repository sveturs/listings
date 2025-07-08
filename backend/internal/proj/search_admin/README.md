# Search Admin Module

Модуль управления конфигурацией поиска для платформы Sve Tu.

## Описание

Этот модуль предоставляет REST API для управления различными аспектами поисковой системы:
- Веса полей поиска
- Синонимы
- Правила транслитерации
- Статистика поиска
- Общая конфигурация поиска

## API Endpoints

### Конфигурация

- `GET /api/v1/search/config` - Получить общую конфигурацию
- `PUT /api/v1/search/config` - Обновить конфигурацию (требует админ права)

### Веса полей

- `GET /api/v1/search/config/weights` - Получить все веса
- `GET /api/v1/search/config/weights/:field` - Получить вес для конкретного поля
- `POST /api/v1/search/config/weights` - Создать новый вес (требует админ права)
- `PUT /api/v1/search/config/weights/:id` - Обновить вес (требует админ права)
- `DELETE /api/v1/search/config/weights/:id` - Удалить вес (требует админ права)

### Синонимы

- `GET /api/v1/search/config/synonyms?language=ru` - Получить синонимы
- `POST /api/v1/search/config/synonyms` - Создать синоним (требует админ права)
- `PUT /api/v1/search/config/synonyms/:id` - Обновить синоним (требует админ права)
- `DELETE /api/v1/search/config/synonyms/:id` - Удалить синоним (требует админ права)

### Транслитерация

- `GET /api/v1/search/config/transliteration` - Получить правила
- `POST /api/v1/search/config/transliteration` - Создать правило (требует админ права)
- `PUT /api/v1/search/config/transliteration/:id` - Обновить правило (требует админ права)
- `DELETE /api/v1/search/config/transliteration/:id` - Удалить правило (требует админ права)

### Статистика

- `GET /api/v1/search/statistics?limit=100` - Получить статистику поиска
- `GET /api/v1/search/statistics/popular?limit=10` - Получить популярные запросы

## Примеры использования

### Обновление весов полей

```bash
curl -X PUT http://localhost:3000/api/v1/search/config/weights/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "field_name": "title",
    "weight": 5.0,
    "description": "Заголовок объявления"
  }'
```

### Добавление синонима

```bash
curl -X POST http://localhost:3000/api/v1/search/config/synonyms \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "term": "ноутбук",
    "synonyms": ["лэптоп", "notebook", "портативный компьютер"],
    "language": "ru"
  }'
```

### Обновление общей конфигурации

```bash
curl -X PUT http://localhost:3000/api/v1/search/config \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "min_search_length": 3,
    "max_suggestions": 15,
    "fuzzy_enabled": true,
    "fuzzy_max_edits": 2,
    "synonyms_enabled": true,
    "transliteration_enabled": true
  }'
```

## База данных

Модуль использует следующие таблицы:
- `search_config` - Общая конфигурация
- `search_weights` - Веса полей
- `search_synonyms_config` - Синонимы
- `transliteration_rules` - Правила транслитерации
- `search_statistics` - Статистика поиска

## Безопасность

Все модифицирующие операции (POST, PUT, DELETE) требуют админских прав и защищены middleware `AdminAuthRequired`.