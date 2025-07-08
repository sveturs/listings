# REST API для управления конфигурацией поиска - Отчёт о реализации

## Выполненная работа

### 1. Созданы доменные модели (domain layer)
**Файл**: `/backend/internal/domain/search_config.go`

Созданы следующие структуры:
- `SearchWeight` - веса для полей поиска
- `SearchSynonym` - синонимы для поиска
- `TransliterationRule` - правила транслитерации
- `SearchStatistics` - статистика поиска
- `SearchConfig` - общая конфигурация поиска

### 2. Созданы миграции базы данных
**Файлы**: 
- `/backend/migrations/000088_search_configuration.up.sql`
- `/backend/migrations/000088_search_configuration.down.sql`

Созданы таблицы:
- `search_weights` - веса полей поиска
- `search_synonyms_config` - синонимы (переименована чтобы избежать конфликта с существующей)
- `transliteration_rules` - правила транслитерации
- `search_statistics` - статистика поиска
- `search_config` - общая конфигурация

Добавлены начальные данные:
- Веса по умолчанию для основных полей
- Базовые синонимы на русском языке
- Правила транслитерации кириллица → латиница

### 3. Создан storage слой
**Файл**: `/backend/internal/storage/postgres/search_config_repository.go`

Реализован репозиторий `SearchConfigRepository` с методами:
- CRUD операции для весов
- CRUD операции для синонимов
- CRUD операции для правил транслитерации
- Создание статистики поиска
- Получение популярных запросов
- Управление общей конфигурацией

### 4. Создан service слой
**Файл**: `/backend/internal/proj/search_admin/service/service.go`

Реализована бизнес-логика с валидацией:
- Проверка диапазонов значений
- Проверка обязательных полей
- Проверка уникальности
- Установка значений по умолчанию

### 5. Создан handler слой
**Файлы**:
- `/backend/internal/proj/search_admin/handler/handler.go`
- `/backend/internal/proj/search_admin/handler/routes.go`

Реализованы REST endpoints:
- Управление весами полей
- Управление синонимами
- Управление транслитерацией
- Просмотр статистики
- Управление конфигурацией

Все модифицирующие операции защищены `AdminAuthRequired` middleware.

### 6. Создан модуль
**Файл**: `/backend/internal/proj/search_admin/module.go`

Модуль интегрирован в систему и реализует интерфейс `RouteRegistrar`.

### 7. Обновлён server.go
**Файл**: `/backend/internal/server/server.go`

Добавлена регистрация нового модуля `search_admin`.

### 8. Создана документация
**Файл**: `/backend/internal/proj/search_admin/README.md`

Документация включает:
- Описание модуля
- Список API endpoints
- Примеры использования
- Информацию о безопасности

## API Endpoints

### Общая конфигурация
- `GET /api/v1/search/config` - получить конфигурацию
- `PUT /api/v1/search/config` - обновить конфигурацию (админ)

### Веса полей
- `GET /api/v1/search/config/weights` - список весов
- `GET /api/v1/search/config/weights/:field` - вес конкретного поля
- `POST /api/v1/search/config/weights` - создать вес (админ)
- `PUT /api/v1/search/config/weights/:id` - обновить вес (админ)
- `DELETE /api/v1/search/config/weights/:id` - удалить вес (админ)

### Синонимы
- `GET /api/v1/search/config/synonyms?language=ru` - список синонимов
- `POST /api/v1/search/config/synonyms` - создать синоним (админ)
- `PUT /api/v1/search/config/synonyms/:id` - обновить синоним (админ)
- `DELETE /api/v1/search/config/synonyms/:id` - удалить синоним (админ)

### Транслитерация
- `GET /api/v1/search/config/transliteration` - список правил
- `POST /api/v1/search/config/transliteration` - создать правило (админ)
- `PUT /api/v1/search/config/transliteration/:id` - обновить правило (админ)
- `DELETE /api/v1/search/config/transliteration/:id` - удалить правило (админ)

### Статистика
- `GET /api/v1/search/statistics?limit=100` - статистика поиска
- `GET /api/v1/search/statistics/popular?limit=10` - популярные запросы

## Примеры использования

### Обновление веса поля
```bash
curl -X PUT http://localhost:3000/api/v1/search/config/weights/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
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
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -d '{
    "term": "ноутбук",
    "synonyms": ["лэптоп", "notebook", "портативный компьютер"],
    "language": "ru"
  }'
```

### Получение популярных запросов
```bash
curl -X GET http://localhost:3000/api/v1/search/statistics/popular?limit=20
```

## Следующие шаги

1. **Применить миграции**:
   ```bash
   cd backend
   # Использовать инструмент миграций проекта
   ```

2. **Сгенерировать Swagger документацию**:
   ```bash
   cd backend
   make generate-types
   ```

3. **Интегрировать с существующим поиском**:
   - Обновить search сервис для использования конфигурации из БД
   - Добавить логирование поисковых запросов
   - Использовать веса и синонимы при поиске

4. **Создать админ-панель во frontend**:
   - Интерфейс для управления весами
   - Интерфейс для управления синонимами
   - Дашборд со статистикой

## Замечания

1. Таблица синонимов названа `search_synonyms_config` чтобы избежать конфликта с существующей таблицей.
2. Все операции изменения требуют админских прав через middleware.
3. Модуль готов к использованию после применения миграций.
4. Есть проблема с циклической зависимостью в существующем коде (searchlogs), которая требует отдельного исправления.