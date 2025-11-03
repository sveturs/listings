# Backend - Svetu Marketplace

Бэкенд монолита Svetu Marketplace. Реализован на Go с использованием Fiber framework.

## Оглавление

- [Обязательные зависимости](#обязательные-зависимости)
- [Установка и запуск](#установка-и-запуск)
- [Архитектура](#архитектура)
- [Конфигурация](#конфигурация)

---

## ⚠️ ОБЯЗАТЕЛЬНЫЕ ЗАВИСИМОСТИ

### Listings Microservice (REQUIRED)

**Backend НЕ ЗАПУСТИТСЯ** без работающего listings микросервиса.

**Причина:** После миграции на микросервисную архитектуру (Phase 7.3), весь каталог (объявления C2C + товары B2C) обслуживается через отдельный микросервис. Монолитная реализация удалена.

#### Как запустить:

1. **Склонируйте listings репозиторий:**
   ```bash
   git clone https://github.com/sveturs/listings.git /p/github.com/sveturs/listings
   ```

2. **Запустите микросервис:**
   ```bash
   cd /p/github.com/sveturs/listings
   go run ./cmd/server/main.go
   # Слушает на localhost:50051 (gRPC)
   ```

3. **Проверьте подключение:**
   ```bash
   grpcurl -plaintext localhost:50051 list
   # Должен показать: listings.v1.ListingsService
   ```

4. **Теперь можно запускать backend:**
   ```bash
   cd /p/github.com/sveturs/svetu/backend
   go run ./cmd/api/main.go
   # Успешно подключится к микросервису
   ```

#### Troubleshooting:

```bash
# Ошибка: "failed to connect to listings microservice"
# Решение: Убедитесь что микросервис запущен на localhost:50051

# Проверка:
netstat -tlnp | grep 50051
# Должно показать: tcp ... LISTEN ... go

# Логи микросервиса:
cd /p/github.com/sveturs/listings
go run ./cmd/server/main.go
# Должно быть: "gRPC server listening on :50051"
```

---

## Установка и запуск

### Предварительные требования:

- Go 1.21+
- PostgreSQL 14+
- Redis
- MinIO (для файлового хранилища)
- OpenSearch (для поиска)
- **Listings Microservice** (см. выше)

### Шаги установки:

1. Клонируйте репозиторий
2. Скопируйте `.env.example` в `.env` и настройте переменные окружения
3. Убедитесь что listings микросервис запущен
4. Запустите миграции базы данных:
   ```bash
   ./migrator up
   ```
5. Запустите backend:
   ```bash
   go run ./cmd/api/main.go
   ```

Backend будет доступен на `http://localhost:3000`

---

## Архитектура

### Микросервисная архитектура (Phase 7.3+)

```
┌─────────────────────────────────────────────────┐
│              Backend (Monolith)                 │
│                                                 │
│  ┌───────────────────────────────────────────┐ │
│  │  Marketplace Handler                      │ │
│  │    ↓                                      │ │
│  │  Storage Layer (db_marketplace.go)        │ │
│  │    ↓                                      │ │
│  │  gRPC Client (marketplace_grpc_client.go) │ │
│  └─────────────────┬─────────────────────────┘ │
└────────────────────┼───────────────────────────┘
                     │
                     │ gRPC (port 50051)
                     ↓
       ┌─────────────────────────────┐
       │  Listings Microservice      │
       │                             │
       │  - OpenSearch indexing      │
       │  - CRUD operations          │
       │  - Search & filtering       │
       │  - Category management      │
       └─────────────────────────────┘
```

### Ключевые компоненты:

- **Backend (Monolith)**: Основной бэкенд с бизнес-логикой
- **Listings Microservice**: Независимый микросервис для управления каталогом
- **TrafficRouter**: Упрощён для dev окружения (всегда 100% на микросервис)
- **gRPC Client**: Прозрачная интеграция с микросервисом

### Структура директорий:

```
backend/
├── cmd/
│   └── api/              # Точка входа приложения
├── internal/
│   ├── config/           # Конфигурация
│   ├── domain/           # Доменные модели
│   ├── middleware/       # Middleware (auth, CORS, logger)
│   ├── proj/             # Бизнес-логика по модулям
│   │   ├── marketplace/  # Marketplace (делегирует в микросервис)
│   │   ├── users/        # Пользователи
│   │   ├── payments/     # Платежи
│   │   └── ...
│   ├── server/           # HTTP сервер (Fiber)
│   └── storage/          # Слой данных
│       ├── postgres/     # PostgreSQL репозитории
│       ├── opensearch/   # OpenSearch клиент
│       └── filestorage/  # MinIO клиент
├── migrations/           # SQL миграции
└── docs/                 # Swagger документация
```

---

## Конфигурация

Основные переменные окружения описаны в `.env.example`. Критичные настройки:

### Обязательные:
- `DATABASE_URL` - PostgreSQL connection string
- `LISTINGS_GRPC_URL` - URL микросервиса (обязательно!)
- `AUTH_SERVICE_URL` - URL auth микросервиса
- `REDIS_URL` - Redis для кеширования

### Опциональные:
- `USE_MARKETPLACE_MICROSERVICE` - Feature flag (по умолчанию true)
- `MARKETPLACE_ROLLOUT_PERCENT` - Процент трафика на микросервис (по умолчанию 100)

**Примечание:** В dev окружении микросервис используется всегда (100% трафика), feature flags игнорируются.

---

## Документация

- Swagger UI: `http://localhost:3000/swagger/`
- API спецификация: `/docs/swagger.json`
- Руководства по миграции: `/docs/migration/`

---

## Версионирование

Текущая версия хранится в `internal/version/version.go`

Для обновления версии используйте:
```bash
/home/dim/.local/bin/bump-version.sh patch  # 0.2.1 -> 0.2.2
/home/dim/.local/bin/bump-version.sh minor  # 0.2.1 -> 0.3.0
/home/dim/.local/bin/bump-version.sh major  # 0.2.1 -> 1.0.0
```

---

## Тестирование

```bash
# Запуск тестов
go test ./...

# Тесты с покрытием
go test -cover ./...

# Линтинг
make lint

# Форматирование
make format
```

---

## Дополнительная информация

- Основная документация проекта: `/p/github.com/sveturs/svetu/CLAUDE.md`
- Troubleshooting: `/p/github.com/sveturs/svetu/docs/CLAUDE_TROUBLESHOOTING.md`
- Database Guidelines: `/p/github.com/sveturs/svetu/docs/CLAUDE_DATABASE_GUIDELINES.md`
