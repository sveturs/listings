# Backend Handlers Index

Полный список всех обработчиков в системе Sve Tu Platform.

## Основные обработчики

### 1. Users Handler
- **Путь**: `/internal/proj/users/handler/`
- **Префикс**: `/api/v1/users`
- **Описание**: Управление пользователями, аутентификация, профили
- **Документация**: [users_handler.md](users_handler.md)

### 2. Marketplace Handler
- **Путь**: `/internal/proj/marketplace/handler/`
- **Префикс**: `/api/v1/marketplace`
- **Описание**: Объявления, категории, изображения, чат
- **Документация**: [marketplace_handler.md](marketplace_handler.md)

### 3. Payments Handler
- **Путь**: `/internal/proj/payments/handler/`
- **Префикс**: `/api/v1/payments`
- **Описание**: Платежи, интеграция с AllSecure и Stripe
- **Документация**: [payments_handler.md](payments_handler.md)

### 4. Notifications Handler
- **Путь**: `/internal/proj/notifications/handler/`
- **Префикс**: `/api/v1/notifications`
- **Описание**: Уведомления пользователей
- **Документация**: [notifications_handler.md](notifications_handler.md)

### 5. Reviews Handler
- **Путь**: `/internal/proj/reviews/handler/`
- **Префикс**: `/api/v1/reviews`
- **Описание**: Отзывы и рейтинги
- **Документация**: [reviews_handler.md](reviews_handler.md)

### 6. Analytics Handler
- **Путь**: `/internal/proj/analytics/handler/`
- **Префикс**: `/api/v1/analytics`
- **Описание**: Аналитика и статистика
- **Документация**: [analytics_handler.md](analytics_handler.md)

### 7. Storefronts Handler
- **Путь**: `/internal/proj/storefronts/handler/`
- **Префикс**: `/api/v1/storefronts`
- **Описание**: Витрины магазинов, импорт товаров
- **Документация**: [storefronts_handler.md](storefronts_handler.md)

### 8. Balance Handler
- **Путь**: `/internal/proj/balance/handler/`
- **Префикс**: `/api/v1/balance`
- **Описание**: Балансы пользователей, транзакции
- **Документация**: [balance_handler.md](balance_handler.md)

### 9. Contacts Handler
- **Путь**: `/internal/proj/contacts/handler/`
- **Префикс**: `/api/v1/contacts`
- **Описание**: Управление контактами пользователей
- **Документация**: [contacts_handler.md](contacts_handler.md)

### 10. Chat Handler
- **Путь**: `/internal/proj/marketplace/handler/chat.go`
- **Префикс**: `/api/v1/marketplace/chats`
- **Описание**: Чат между покупателями и продавцами
- **Документация**: [chat_handler.md](chat_handler.md)

### 11. Admin Handler
- **Путь**: `/internal/proj/users/handler/admin.go`
- **Префикс**: `/api/v1/admin`
- **Описание**: Административные функции
- **Документация**: [admin_handler.md](admin_handler.md)

## Дополнительные обработчики

### 12. DocServer Handler
- **Путь**: `/internal/proj/docserver/handler/`
- **Префикс**: `/api/v1/docs`
- **Описание**: API для доступа к документации (только для админов)
- **Документация**: [miscellaneous_handlers.md](miscellaneous_handlers.md#1-docserver-handler)

### 13. Geocode Handler
- **Путь**: `/internal/proj/geocode/handler/`
- **Префикс**: `/api/v1/geocode`
- **Описание**: Геокодирование и поиск городов
- **Документация**: [miscellaneous_handlers.md](miscellaneous_handlers.md#2-geocode-handler)

### 14. Global Handler (Unified Search)
- **Путь**: `/internal/proj/global/handler/`
- **Префикс**: `/api/v1/search`
- **Описание**: Унифицированный поиск по всем типам товаров
- **Документация**: [miscellaneous_handlers.md](miscellaneous_handlers.md#3-global-handler-unified-search)

## Архитектура

### Структура модулей
```
internal/proj/{module}/
├── handler/      # HTTP обработчики
├── service/      # Бизнес-логика
├── storage/      # Работа с БД
└── routes/       # Регистрация маршрутов
```

### Общие принципы
1. **Модульность**: Каждый функциональный блок - отдельный модуль
2. **Dependency Injection**: Сервисы передаются через конструкторы
3. **Middleware**: Централизованная обработка auth, CORS, rate limiting
4. **Swagger**: Все эндпоинты документированы через swagger комментарии
5. **Безопасность**: Валидация входных данных, защита от атак

### Регистрация маршрутов
Все обработчики реализуют интерфейс `RouteRegistrar`:
```go
type RouteRegistrar interface {
    RegisterRoutes(app *fiber.App, middleware *middleware.Middleware) error
    GetPrefix() string
}
```

## Статистика

- **Всего обработчиков**: 14
- **Основных модулей**: 11
- **Дополнительных модулей**: 3
- **Публичных эндпоинтов**: ~40%
- **Защищенных эндпоинтов**: ~60%

## См. также
- [API Documentation (Swagger)](http://localhost:3000/swagger/index.html)
- [Middleware Documentation](../middleware/index.md)
- [Services Documentation](../services/index.md)