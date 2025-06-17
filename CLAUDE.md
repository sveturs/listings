# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## IMPORTANT WORKFLOW RULES

- **Язык общения**: Use Russian in your answers

- **Управление задачами**: ВСЕГДА Информацию о прогрессе текущей задачи храни в файле `/memory-bank/currentTask.md`. При начале новой сессии получай статус из этого файла. Прогресс сохраняй в `/memory-bank/tasks/<next_number_as_\d\d\d>-<task_name>.md`. После завершения задачи очищай `/memory-bank/current-task.md`

- **Коммиты**: Подбирай описание для коммита на основании внесенных изменений. НЕ ИСПОЛЬЗУЙ УКАЗАНИЕ АВТОРСТВА CLAUDE

- **Разработка**:
   - Запускай сервер командой `yarn dev -p 3001`. Если надо его остановить, ищи процесс для остановке ИМЕННО ПО ПОРТУ 3001. НЕ ИСПОЛЬЗУЙ `killall node 2>/dev/null || true`
   - Используй mcp playwright для работы с браузером Google Chrome
   - ИСПОЛЬЗУЙ ДЕФОЛТНЫЕ КЛАССЫ DAISYUI. Используй mcp context7 для поиска по документации daisyui при работе с версткой
   - **Галерея изображений**: Для отображения фотографий ВСЕГДА используй компонент `ImageGallery` из `src/components/reviews/ImageGallery.tsx`. Этот компонент обеспечивает единообразный просмотр изображений с полупрозрачным фоном, навигацией и миниатюрами

- **Качество кода**: Перед завершением задачи ОБЯЗАТЕЛЬНО выполни `yarn format && yarn lint && yarn build`. Задача считается выполненной только если все команды прошли успешно

- **Зависимости**: При добавлении ключевых зависимостей указывай их в разделе "Key Dependencies"

- **Переводы**: Backend возвращает только placeholder'ы (например, "notifications.getError"), а фронт переводит их через файлы `frontend/svetu/src/messages/{en,ru}.json`

- **Переменные окружения**: При добавлении новых переменных:
   1. Обновляй `.env` и `.env.example` (с примерами и комментариями)
   2. Обновляй интерфейсы в `src/config/types.ts` (`EnvVariables` и `Config`)
   3. Добавляй обработку в `src/config/index.ts` в методе `loadConfig()`

## Проект: Sve Tu Platform - Marketplace

Состоит из
- frontend: React, NextJS, Tailwind, DaisyUI 
- backend: Go, Postgres, Minio, OpenSearch

## Frontend (NextJS)

### Development Commands

- `yarn dev -p 3001` - Start the development server with Turbopack at http://localhost:3001
- `yarn build` - Create an optimized production build
- `yarn start` - Run the production server
- `yarn lint` - Run ESLint for code quality checks
- `yarn lint:fix` - Fix ESLint errors automatically
- `yarn format` - Format all files with Prettier
- `yarn format:check` - Check formatting without changes

### Frontend Architecture

This is a Next.js 15.3.2 application using:

- React 19 with TypeScript
- Tailwind CSS v4 for styling
- App Router (located in `src/app/`)
- ESLint configured with Next.js and TypeScript rules
- Internationalization with next-intl (en/ru locales)
- Centralized configuration management in `src/config/`
- Google OAuth 2.0 authentication
- State management with Redux Toolkit (НЕ Zustand!)
  - Store: `src/store/`
  - Slices: `src/store/slices/`
  - Hooks: `src/store/hooks.ts` и `src/hooks/useChat.ts`


### Environment Variables for Frontend

Configuration is managed through environment variables.

#### Configuration Management

All environment variables are centralized in the src/config/ module:

- `src/config/types.ts` - TypeScript interfaces for configuration
- `src/config/index.ts` - Configuration manager with helper methods

### Frontend Key Dependencies

- next-intl: For internationalization support
- prettier: Code formatter with ESLint integration
- daisyui: Component library for Tailwind CSS
- @reduxjs/toolkit: State management (Redux Toolkit)
- react-redux: React bindings for Redux

### Важная информация о поиске и индексировании

**⚠️ ВАЖНО**: Главная страница маркетплейса получает данные из OpenSearch, а НЕ напрямую из PostgreSQL. 

При изменении данных в базе PostgreSQL (например, user_id объявления) изменения НЕ отобразятся на главной странице до переиндексирования OpenSearch.

Для переиндексирования используйте:
```bash
#  через команду
cd backend && ./reindex
```

### Authentication System

Authentication flow:
1. User clicks "Sign in with Google" button
2. Redirects to Google OAuth consent
3. Google redirects back to backend callback
4. Backend creates session and redirects to frontend
5. Frontend fetches user data via `/auth/session`


## Backend (Go)
```bash
# Сборка
cd backend && go build -o bin/api ./cmd/api

# Запуск с горячей перезагрузкой
cd backend && go run ./cmd/api/main.go

# Тесты
cd backend && go test ./...

# Миграции базы данных
cd backend && make migrate_up
```

## Архитектура Backend

### Backend структура
```
backend/
├── cmd/api/          # Точка входа приложения
├── internal/
│   ├── config/       # Конфигурация
│   ├── domain/       # Доменные модели
│   ├── middleware/   # Auth, CORS, Logger
│   ├── proj/         # Бизнес-логика модулей
│   │   ├── marketplace/
│   │   ├── users/
│   │   ├── payments/
│   │   └── ...
│   ├── server/       # HTTP сервер (Fiber)
│   └── storage/      # Репозитории
│       ├── postgres/
│       ├── opensearch/
│       └── minio/
└── migrations/       # SQL миграции
```

### Сервисы инфраструктуры
- **PostgreSQL** - основная база данных
- **OpenSearch** - полнотекстовый поиск
- **MinIO** - S3-совместимое хранилище изображений
- **Nginx** - обратный прокси и статика
- **Harbor** - приватный Docker registry

### External Services

- **MinIO**: Object storage for images
   - Local: http://localhost:9000
   - Production: https://svetu.rs
   - Images are served from `/listings/` path

### Конфигурация Backend

#### Переменные окружения

- Backend: `.env` в корне backend/

### Важные файлы конфигурации
- `docker-compose.yml` - локальная разработка
- `backend/categories.yaml` - структура категорий маркетплейса

## База данных

### Основные таблицы
- `users` - пользователи и аутентификация
- `marketplace_listings` - объявления
- `marketplace_categories` - категории с атрибутами
- `translations` - мультиязычные переводы
- `marketplace_images` - изображения объявлений
- `marketplace_chats` - чаты и сообщения
TODO: Нужно доработать

### Работа с миграциями

TODO: Нужно доработать


### Harbor Registry

```bash
# Авторизация
docker login -u <harbor_user> -p <harbor_user> harbor.svetu.rs

# Загрузка образа
docker push harbor.svetu.rs/svetu/backend/api:latest
```

### Управление сервером
TODO: Требуется доработать

## API Documentation (Swagger)

Backend использует **goswag** для автоматической генерации OpenAPI/Swagger документации.

### Генерация документации

```bash
cd backend && make docs
```

Эта команда запускает `swag init` со следующими параметрами:
- `--generalInfo` - путь к основному файлу приложения
- `--output` - директория для генерации документации
- `--outputTypes` - форматы вывода (go, json, yaml)
- `--parseInternal` - парсинг внутренних пакетов
- `--parseDependency` - парсинг зависимостей
- `--parseDepth` - глубина парсинга зависимостей

### Swagger аннотации

Все HTTP endpoints должны быть документированы с помощью Swagger комментариев:

```go
// GetReviews returns filtered list of reviews
// @Summary Get reviews list
// @Description Returns paginated list of reviews with filters
// @Tags reviews
// @Accept json
// @Produce json
// @Param entity_id query int false "Entity ID filter"
// @Success 200 {object} ReviewsListResponse "List of reviews"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /api/v1/reviews [get]
func (h *ReviewHandler) GetReviews(c *fiber.Ctx) error {
    // ...
}
```

### Структура Swagger аннотаций

- `@Summary` - краткое описание endpoint'а
- `@Description` - подробное описание
- `@Tags` - группировка endpoints в документации
- `@Accept` - формат входных данных
- `@Produce` - формат выходных данных
- `@Param` - описание параметров
- `@Success` - описание успешного ответа
- `@Failure` - описание ошибок
- `@Security` - требования авторизации
- `@Router` - путь и HTTP метод

### Сгенерированные файлы

После выполнения `make docs` создаются:
- `docs/docs.go` - Go код для встраивания документации
- `docs/swagger.json` - OpenAPI спецификация в JSON
- `docs/swagger.yaml` - OpenAPI спецификация в YAML

### Просмотр документации

Swagger UI доступен по адресу: http://localhost:3000/swagger/index.html (в режиме разработки)

## API Contract Management (OpenAPI v3)

Проект использует подход **Contract-First Development** на основе OpenAPI v3 схемы:

### Workflow взаимодействия Backend ↔ Frontend

1. **Backend (первичный источник)**: 
   - Разработчики пишут Swagger аннотации в Go коде
   - Аннотации описывают все endpoints, их параметры и типы ответов
   
2. **Генерация OpenAPI схемы**:
   ```bash
   cd backend && make docs
   ```
   - Создается `docs/swagger.json` - OpenAPI v3 спецификация
   
3. **Генерация типов для Frontend**:
   ```bash
   cd backend && make generate-types
   ```
   - Копирует swagger.json в frontend
   - Запускает генерацию TypeScript типов
   - Создает типизированные интерфейсы в `frontend/svetu/src/types/generated/api.ts`

### Важно при разработке

- **Backend разработчики**: Всегда документируйте endpoints с помощью Swagger аннотаций
- **Frontend разработчики**: Используйте сгенерированные типы из `@/types/generated/api`
- **При изменении API**: Обязательно перегенерируйте типы командой `make generate-types`

### Проверка качества кода

Перед завершением любой задачи на Frontend **ОБЯЗАТЕЛЬНО** выполните:

```bash
yarn format && yarn lint && yarn build
```

Эта команда:
1. `yarn format` - форматирует код согласно правилам Prettier
2. `yarn lint` - проверяет код на соответствие правилам ESLint
3. `yarn build` - создает production сборку и проверяет типы TypeScript

Задача считается выполненной только если все три команды выполнились успешно!
