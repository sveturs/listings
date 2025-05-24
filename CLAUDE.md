# Инструкция по работе с проектом Sve Tu Platform для AI агента

## Общая информация
- Проект называется Sve Tu Platform
- В разработке находится маркетплейс
- Основная цель: площадка для размещения и поиска объявлений

## Рабочие окружения

1. Локальная разработка: `/data/hostel-booking-system/`
2. Боевой сервер: ssh root@svetu.rs, проект в `/opt/hostel-booking-system/`

## Команды сборки и запуска

- Backend: `go build -o main ./cmd/api`
- Frontend: `cd frontend/hostel-frontend && npm run build`
- Запуск: `docker-compose up -d`

## Деплой

- На сервер: `./harbor-scripts/blue_green_deploy_on_svetu.rs.sh [backend|frontend|all]`
- Синий-зеленый деплой для нулевого простоя

## Структура проекта

### Backend
- Язык: Go
- Основная точка входа: `backend/cmd/api/main.go`
- Внутренняя структура:
  - `internal/config` - конфигурация
  - `internal/domain` - бизнес-логика
  - `internal/middleware` - промежуточные обработчики
  - `internal/server` - HTTP сервер
  - `internal/storage` - работа с хранилищами данных
- Миграции БД: `backend/migrations`

### Frontend
- Framework: React 18 с TypeScript
- Структура:
  - `frontend/hostel-frontend/src/components` - переиспользуемые компоненты
  - `frontend/hostel-frontend/src/pages` - страницы приложения
  - `frontend/hostel-frontend/src/contexts` - контексты для состояния
  - `frontend/hostel-frontend/src/api` - работа с API
  - `frontend/hostel-frontend/src/hooks` - пользовательские хуки
  - `frontend/hostel-frontend/public/env.js` - конфигурация окружения

### Docker-инфраструктура
- `docker-compose.yml` - для локальной разработки
- `docker-compose.prod.yml` - для продакшена
- Основные сервисы:
  - backend - Go API
  - frontend - React приложение
  - db - PostgreSQL
  - opensearch - поисковый движок
  - minio - хранилище файлов
  - nginx - веб-сервер и прокси

## Технологии

### Backend
- Go
- PostgreSQL - основная БД
- OpenSearch - поисковый движок для объявлений
- MinIO - хранилище файлов и изображений

### Frontend
- React 18
- TypeScript
- Material UI (MUI) - компоненты интерфейса
- React Router - маршрутизация
- React Query - управление запросами
- i18next - интернационализация
- Leaflet - карты и геолокация
- Zustand - управление состоянием

### Инфраструктура
- Docker и Docker Compose
- Nginx
- Harbor - реестр Docker-образов
- Certbot - SSL-сертификаты

## Ключевые функции маркетплейса
- Публикация и поиск объявлений
- Фильтрация и категоризация
- Профили пользователей
- Чат между пользователями
- Система отзывов
- Геолокация и отображение на карте
- Многоязычность (русский/английский/сербский)
- Административная панель

## Конфигурация

- Frontend-настройки: `frontend/hostel-frontend/public/env.js`
- Nginx: `/opt/hostel-booking-system/nginx.conf` (на сервере)
- Docker: `docker-compose.yml`, `docker-compose.prod.yml`

## Текущие задачи разработки
- Улучшение маркетплейса
- Рассмотрение возможности миграции фронтенда на Next.js
- Оптимизация производительности и SEO

Пожалуйста, отвечай мне на русском языке. в