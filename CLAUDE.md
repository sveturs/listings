# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Important Workflow Rules

- **Язык общения**: Use Russian in your answers

- **Управление задачами**: Информацию о прогрессе текущей задачи храни в файле `/memory-bank/currentTask.md`. При начале новой сессии получай статус из этого файла. Прогресс сохраняй в `/memory-bank/tasks/<next_number_as_\d\d\d>-<task_name>.md`. После завершения задачи очищай `/memory-bank/current-task.md`

- **Коммиты**: Подбирай описание для коммита на основании внесенных изменений. Не используй указание авторства Claude

- **Разработка**:
   - Запускай сервер командой `yarn dev` на порту 3999
   - Останавливай сервер по порту 3999
   - Используй mcp playwright для работы с браузером Google Chrome
   - Используй mcp context7 для поиска по документации daisyui при работе с версткой

- **Качество кода**: Проверяй `yarn lint` и `yarn build` перед завершением задачи. После успешного `yarn build` выполни `yarn format`

- **Зависимости**: При добавлении ключевых зависимостей указывай их в разделе "Key Dependencies"

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

- `yarn dev -p 3040` - Start the development server with Turbopack at http://localhost:3040
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


## Backend (Go)
```bash
# Сборка
cd backend && go build -o main ./cmd/api

# Запуск с горячей перезагрузкой
cd backend && air

# Тесты
cd backend && go test ./...

# Миграции базы данных
cd backend && migrate -path ./migrations -database "postgresql://..." up
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
