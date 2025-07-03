# Паспорт Backend слоя (Go)

## Общие принципы

### Соглашения о наименованиях
- **Go структуры**: PascalCase для полей (`UserID`, `CreatedAt`)
- **JSON теги**: строго snake_case (`json:"user_id"`, `json:"created_at"`)
- **API параметры**: snake_case во всех запросах и ответах
- **База данных**: snake_case соответствует JSON тегам

### JSON сериализация
- Все модели имеют explicit JSON теги
- Приватные поля исключаются через `json:"-"`
- Опциональные поля используют `omitempty`
- Null значения обрабатываются через указатели

### Validation
- Использование тегов `validate` для проверки входных данных
- Кастомные валидаторы для сложных бизнес-правил
- Единообразная обработка ошибок валидации

## Структура паспортов Backend

### Модели данных
- [User Models](./models/user.md) - Пользователи и аутентификация
- [Marketplace Models](./models/marketplace.md) - Товары и объявления
- [Storefront Models](./models/storefront.md) - Витрины и магазины
- [Payment Models](./models/payments.md) - Платежи и балансы
- [Common Types](./models/common.md) - Общие типы и структуры

### API Handlers
- [Auth Handlers](./handlers/auth.md) - Аутентификация и авторизация
- [Marketplace Handlers](./handlers/marketplace.md) - Маркетплейс API
- [Storefront Handlers](./handlers/storefront.md) - Витрины API
- [Payment Handlers](./handlers/payments.md) - Платежи API
- [File Upload Handlers](./handlers/files.md) - Загрузка файлов

### API Documentation
- [Swagger/OpenAPI типы](./swagger-types.md) - Документация API
- [Endpoints Reference](./api-endpoints.md) - Список всех endpoints

## Ключевые особенности

### Consistency
- **100% snake_case** для JSON API
- Предсказуемое именование ID полей: `{entity}_id`
- Стандартные timestamp поля: `created_at`, `updated_at`

### Error Handling
- Единообразные error response структуры
- Локализованные сообщения об ошибках через placeholders
- Логирование реальных ошибок, возврат safe сообщений

### Security
- JWT токены для аутентификации
- Проверка владельца ресурсов
- Валидация всех входных данных
- Исключение чувствительных данных из API responses